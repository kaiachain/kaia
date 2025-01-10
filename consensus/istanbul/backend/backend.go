// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from quorum/consensus/istanbul/backend/backend.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"crypto/ecdsa"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulCore "github.com/kaiachain/kaia/consensus/istanbul/core"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/storage/database"
)

const (
	// fetcherID is the ID indicates the block is from Istanbul engine
	fetcherID = "istanbul"
)

var logger = log.NewModuleLogger(log.ConsensusIstanbulBackend)

type BackendOpts struct {
	IstanbulConfig    *istanbul.Config // Istanbul consensus core config
	Rewardbase        common.Address
	PrivateKey        *ecdsa.PrivateKey // Consensus message signing key
	BlsSecretKey      bls.SecretKey     // Randao signing key. Required since Randao fork
	DB                database.DBManager
	GovModule         gov.GovModule
	BlsPubkeyProvider BlsPubkeyProvider // If not nil, override the default BLS public key provider
	NodeType          common.ConnType
}

func New(opts *BackendOpts) consensus.Istanbul {
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)
	backend := &backend{
		config:            opts.IstanbulConfig,
		istanbulEventMux:  new(event.TypeMux),
		privateKey:        opts.PrivateKey,
		address:           crypto.PubkeyToAddress(opts.PrivateKey.PublicKey),
		blsSecretKey:      opts.BlsSecretKey,
		logger:            logger.NewWith(),
		db:                opts.DB,
		commitCh:          make(chan *types.Result, 1),
		candidates:        make(map[common.Address]bool),
		coreStarted:       false,
		recentMessages:    recentMessages,
		knownMessages:     knownMessages,
		rewardbase:        opts.Rewardbase,
		govModule:         opts.GovModule,
		blsPubkeyProvider: opts.BlsPubkeyProvider,
		nodetype:          opts.NodeType,
	}
	if backend.blsPubkeyProvider == nil {
		backend.blsPubkeyProvider = newChainBlsPubkeyProvider()
	}

	backend.currentView.Store(&istanbul.View{Sequence: big.NewInt(0), Round: big.NewInt(0)})
	backend.core = istanbulCore.New(backend, backend.config)
	return backend
}

// ----------------------------------------------------------------------------

type backend struct {
	config           *istanbul.Config
	istanbulEventMux *event.TypeMux
	privateKey       *ecdsa.PrivateKey
	address          common.Address
	blsSecretKey     bls.SecretKey
	core             istanbulCore.Engine
	logger           log.Logger
	db               database.DBManager
	chain            consensus.ChainReader
	stakingModule    staking.StakingModule
	valsetModule     valset.ValsetModule
	consensusModules []kaiax.ConsensusModule
	currentBlock     func() *types.Block
	hasBadBlock      func(hash common.Hash) bool

	// the channels for istanbul engine notifications
	commitCh          chan *types.Result
	proposedBlockHash common.Hash
	sealMu            sync.Mutex
	coreStarted       bool
	coreMu            sync.RWMutex

	// Current list of candidates we are pushing
	candidates map[common.Address]bool
	// Protects the signer fields
	candidatesLock sync.RWMutex

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	recentMessages *lru.ARCCache // the cache of peer's messages
	knownMessages  *lru.ARCCache // the cache of self messages

	rewardbase  common.Address
	currentView atomic.Value //*istanbul.View

	// Reference to the governance.Engine
	govModule gov.GovModule

	// Reference to BlsPubkeyProvider
	blsPubkeyProvider BlsPubkeyProvider

	// Node type
	nodetype common.ConnType

	isRestoringSnapshots atomic.Bool
}

func (sb *backend) NodeType() common.ConnType {
	return sb.nodetype
}

func (sb *backend) GetRewardBase() common.Address {
	return sb.rewardbase
}

func (sb *backend) SetCurrentView(view *istanbul.View) {
	sb.currentView.Store(view)
}

// Address implements istanbul.Backend.Address
func (sb *backend) Address() common.Address {
	return sb.address
}

// Broadcast implements istanbul.Backend.Broadcast
func (sb *backend) Broadcast(prevHash common.Hash, payload []byte) error {
	// send to others
	// TODO Check gossip again in event handle
	// sb.Gossip(valSet, payload)
	// send to self
	msg := istanbul.MessageEvent{
		Hash:    prevHash,
		Payload: payload,
	}
	go sb.istanbulEventMux.Post(msg)
	return nil
}

// Broadcast implements istanbul.Backend.Gossip
func (sb *backend) Gossip(payload []byte) error {
	hash := istanbul.RLPHash(payload)
	sb.knownMessages.Add(hash, true)

	if sb.broadcaster != nil {
		ps := sb.broadcaster.GetCNPeers()
		for addr, p := range ps {
			ms, ok := sb.recentMessages.Get(addr)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					// This peer had this event, skip it
					continue
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}

			m.Add(hash, true)
			sb.recentMessages.Add(addr, m)

			cmsg := &istanbul.ConsensusMsg{
				PrevHash: common.Hash{},
				Payload:  payload,
			}

			// go p.Send(IstanbulMsg, payload)
			go p.Send(IstanbulMsg, cmsg)
		}
	}
	return nil
}

// getTargetReceivers returns a map of nodes which need to receive a message
func (sb *backend) getTargetReceivers() map[common.Address]bool {
	cv, ok := sb.currentView.Load().(*istanbul.View)
	if !ok {
		logger.Error("Failed to assert type from sb.currentView!!", "cv", cv)
		return nil
	}

	// calculates a map of target nodes which need to receive a message
	// committee[currentView].Union(committee[nextView]) => targets
	targets := make(map[common.Address]bool)
	for i := 0; i < 2; i++ {
		roundCommittee, err := sb.GetCommitteeStateByRound(cv.Sequence.Uint64(), cv.Round.Uint64()+uint64(i))
		if err != nil {
			return nil
		}
		// i == 0: get current round's committee. additionally, check if the node is in the current view's committee
		if i == 0 && !roundCommittee.Committee().Contains(sb.Address()) {
			return nil
		}
		for _, val := range roundCommittee.Committee().List() {
			if val != sb.Address() {
				targets[val] = true
			}
		}
	}
	return targets
}

// GossipSubPeer implements istanbul.Backend.Gossip
func (sb *backend) GossipSubPeer(prevHash common.Hash, payload []byte) {
	targets := sb.getTargetReceivers()
	if targets == nil {
		return
	}
	hash := istanbul.RLPHash(payload)
	sb.knownMessages.Add(hash, true)

	if sb.broadcaster != nil && len(targets) > 0 {
		ps := sb.broadcaster.FindCNPeers(targets)
		for addr, p := range ps {
			ms, ok := sb.recentMessages.Get(addr)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, k := m.Get(hash); k {
					// This peer had this event, skip it
					continue
				}
			} else {
				m, _ = lru.NewARC(inmemoryMessages)
			}

			m.Add(hash, true)
			sb.recentMessages.Add(addr, m)

			cmsg := &istanbul.ConsensusMsg{
				PrevHash: prevHash,
				Payload:  payload,
			}

			go p.Send(IstanbulMsg, cmsg)
		}
	}
	return
}

// Commit implements istanbul.Backend.Commit
func (sb *backend) Commit(proposal istanbul.Proposal, seals [][]byte) error {
	// Check if the proposal is a valid block
	block, ok := proposal.(*types.Block)
	if !ok {
		sb.logger.Error("Invalid proposal, %v", proposal)
		return errInvalidProposal
	}
	h := block.Header()
	round := sb.currentView.Load().(*istanbul.View).Round.Int64()
	h = types.SetRoundToHeader(h, round)
	// Append seals into extra-data
	err := writeCommittedSeals(h, seals)
	if err != nil {
		return err
	}
	// update block's header
	block = block.WithSeal(h)

	sb.logger.Info("Committed", "number", proposal.Number().Uint64(), "hash", proposal.Hash(), "address", sb.Address())
	// - if the proposed and committed blocks are the same, send the proposed hash
	//   to commit channel, which is being watched inside the engine.Seal() function.
	// - otherwise, we try to insert the block.
	// -- if success, the ChainHeadEvent event will be broadcasted, try to build
	//    the next block and the previous Seal() will be stopped.
	// -- otherwise, a error will be returned and a round change event will be fired.
	if sb.proposedBlockHash == block.Hash() {
		// feed block hash to Seal() and wait the Seal() result
		sb.commitCh <- &types.Result{Block: block, Round: round}
		return nil
	}

	if sb.broadcaster != nil {
		sb.broadcaster.Enqueue(fetcherID, block)
	}
	return nil
}

// EventMux implements istanbul.Backend.EventMux
func (sb *backend) EventMux() *event.TypeMux {
	return sb.istanbulEventMux
}

// Verify implements istanbul.Backend.Verify
func (sb *backend) Verify(proposal istanbul.Proposal) (time.Duration, error) {
	// Check if the proposal is a valid block
	block, ok := proposal.(*types.Block)
	if !ok {
		sb.logger.Error("Invalid proposal, %v", proposal)
		return 0, errInvalidProposal
	}

	// check bad block
	if sb.HasBadProposal(block.Hash()) {
		return 0, blockchain.ErrBlacklistedHash
	}

	// check block body
	txnHash := types.DeriveSha(block.Transactions(), block.Number())
	if txnHash != block.Header().TxHash {
		return 0, errMismatchTxhashes
	}

	// verify the header of proposed block
	err := sb.VerifyHeader(sb.chain, block.Header(), false)
	// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
	if err == nil || err == errEmptyCommittedSeals {
		return 0, nil
	} else if err == consensus.ErrFutureBlock {
		return time.Unix(block.Header().Time.Int64(), 0).Sub(now()), consensus.ErrFutureBlock
	}
	return 0, err
}

// Sign implements istanbul.Backend.Sign
func (sb *backend) Sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256([]byte(data))
	return crypto.Sign(hashData, sb.privateKey)
}

// CheckSignature implements istanbul.Backend.CheckSignature
func (sb *backend) CheckSignature(data []byte, address common.Address, sig []byte) error {
	signer, err := cacheSignatureAddresses(data, sig)
	if err != nil {
		logger.Error("Failed to get signer address", "err", err)
		return err
	}
	// Compare derived addresses
	if signer != address {
		return errInvalidSignature
	}
	return nil
}

// HasPropsal implements istanbul.Backend.HashBlock
func (sb *backend) HasPropsal(hash common.Hash, number *big.Int) bool {
	return sb.chain.GetHeader(hash, number.Uint64()) != nil
}

func (sb *backend) LastProposal() (istanbul.Proposal, common.Address) {
	block := sb.currentBlock()

	var proposer common.Address
	if block.Number().Cmp(common.Big0) > 0 {
		var err error
		proposer, err = sb.Author(block.Header())
		if err != nil {
			sb.logger.Error("Failed to get block proposer", "err", err)
			return nil, common.Address{}
		}
	}

	// Return header only block here since we don't need block body
	return block, proposer
}

func (sb *backend) HasBadProposal(hash common.Hash) bool {
	if sb.hasBadBlock == nil {
		return false
	}
	return sb.hasBadBlock(hash)
}

func (sb *backend) GetValidatorSet(num uint64) (*istanbul.BlockValSet, error) {
	council, err := sb.valsetModule.GetCouncil(num)
	if err != nil {
		return nil, err
	}

	demoted, err := sb.valsetModule.GetDemotedValidators(num)
	if err != nil {
		return nil, err
	}

	return istanbul.NewBlockValSet(council, demoted), nil
}

func (sb *backend) GetCommitteeState(num uint64) (*istanbul.RoundCommitteeState, error) {
	header := sb.chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, errUnknownBlock
	}

	return sb.GetCommitteeStateByRound(num, uint64(header.Round()))
}

func (sb *backend) GetCommitteeStateByRound(num uint64, round uint64) (*istanbul.RoundCommitteeState, error) {
	blockValSet, err := sb.GetValidatorSet(num)
	if err != nil {
		return nil, err
	}

	committee, err := sb.valsetModule.GetCommittee(num, round)
	if err != nil {
		return nil, err
	}

	proposer, err := sb.valsetModule.GetProposer(num, round)
	if err != nil {
		return nil, err
	}

	committeeSize := sb.govModule.GetParamSet(num).CommitteeSize
	return istanbul.NewRoundCommitteeState(blockValSet, committeeSize, committee, proposer), nil
}

// GetProposer implements istanbul.Backend.GetProposer
func (sb *backend) GetProposer(number uint64) common.Address {
	if h := sb.chain.GetHeaderByNumber(number); h != nil {
		a, _ := sb.Author(h)
		return a
	}
	return common.Address{}
}

func (sb *backend) GetRewardAddress(num uint64, nodeId common.Address) common.Address {
	sInfo, err := sb.stakingModule.GetStakingInfo(num)
	if err != nil {
		return common.Address{}
	}

	for idx, id := range sInfo.NodeIds {
		if id == nodeId {
			return sInfo.RewardAddrs[idx]
		}
	}
	return common.Address{}
}
