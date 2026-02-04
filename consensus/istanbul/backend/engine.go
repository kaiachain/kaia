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
// This file is derived from quorum/consensus/istanbul/backend/engine.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"bytes"
	"encoding/hex"
	"errors"
	"math/big"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	consensuscommon "github.com/kaiachain/kaia/consensus/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulCore "github.com/kaiachain/kaia/consensus/istanbul/core"
	"github.com/kaiachain/kaia/consensus/misc"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/randao"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

const (
	inmemoryPeers    = 200
	inmemoryMessages = 4096
)

var (
	inmemoryBlocks             = 2048 // Number of blocks to precompute validators' addresses
	inmemoryValidatorsPerBlock = 30   // Approximate number of validators' addresses from ecrecover
	signatureAddresses, _      = lru.NewARC(inmemoryBlocks * inmemoryValidatorsPerBlock)
)

// cacheSignatureAddresses extracts the address from the given data and signature and cache them for later usage.
func cacheSignatureAddresses(data []byte, sig []byte) (common.Address, error) {
	sigStr := hex.EncodeToString(sig)
	if addr, ok := signatureAddresses.Get(sigStr); ok {
		return addr.(common.Address), nil
	}
	addr, err := istanbul.GetSignatureAddress(data, sig)
	if err != nil {
		return common.Address{}, err
	}
	signatureAddresses.Add(sigStr, addr)
	return addr, err
}

// Author retrieves the Kaia address of the account that minted the given block.
func (sb *backend) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header)
}

// CanVerifyHeadersConcurrently returns true if concurrent header verification possible, otherwise returns false.
func (sb *backend) CanVerifyHeadersConcurrently() bool {
	return false
}

// PreprocessHeaderVerification prepares header verification for heavy computation before synchronous header verification such as ecrecover.
func (sb *backend) PreprocessHeaderVerification(headers []*types.Header) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, inmemoryBlocks)
	go func() {
		errored := false
		for _, header := range headers {
			var err error
			if errored { // If errored once in the batch, skip the rest
				err = consensus.ErrUnknownAncestor
			} else {
				err = sb.computeSignatureAddrs(header)
			}

			if err != nil {
				errored = true
			}

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// computeSignatureAddrs computes the addresses of signer and validators and caches them.
func (sb *backend) computeSignatureAddrs(header *types.Header) error {
	_, err := ecrecover(header)
	if err != nil {
		return err
	}

	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return err
	}

	proposalSeal := istanbulCore.PrepareCommittedSeal(header.Hash())
	for _, seal := range istanbulExtra.CommittedSeal {
		_, err := cacheSignatureAddresses(proposalSeal, seal)
		if err != nil {
			return istanbul.ErrInvalidSignature
		}
	}
	return nil
}

func (sb *backend) VerifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return consensus.ErrUnknownBlock
	}

	// Header verify before/after magma fork
	if chain.Config().IsMagmaForkEnabled(header.Number) {
		if len(parents) > 0 {
			// the kip71Config used when creating the block number is a previous block config.
			blockNum := header.Number.Uint64()
			pset := sb.govModule.GetParamSet(blockNum)
			kip71 := pset.ToKip71Config()
			if err := misc.VerifyMagmaHeader(parents[len(parents)-1], header, kip71); err != nil {
				return err
			}
		}
		// For Magma fork, BaseFee is allowed even without parents (first header)
	} else if header.BaseFee != nil {
		return consensus.ErrInvalidBaseFee
	}

	// Ensure that the extra data format is satisfied
	if _, err := types.ExtractIstanbulExtra(header); err != nil {
		return istanbul.ErrInvalidExtraDataFormat
	}

	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Get parent header for consensus-dependent checks
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	if parent.Time.Uint64()+sb.config.BlockPeriod > header.Time.Uint64() {
		return istanbul.ErrInvalidTimestamp
	}
	if err := sb.verifySigner(chain, header, parents); err != nil {
		return err
	}

	// VerifyRandao must be after verifySigner because it needs the signer (proposer) address
	if chain.Config().IsRandaoForkEnabled(header.Number) {
		prevMixHash := headerMixHash(chain, parent)
		if err := sb.VerifyRandao(chain, header, prevMixHash); err != nil {
			return err
		}
	} else if header.RandomReveal != nil || header.MixHash != nil {
		return istanbul.ErrUnexpectedRandao
	}

	if err := sb.verifyCommittedSeals(chain, header, nil); err != nil {
		return err
	}

	for _, module := range sb.consensusModules {
		if err := module.VerifyHeader(header); err != nil {
			return err
		}
	}

	return nil
}

// verifySigner checks whether the signer is in parent's validator set
func (sb *backend) verifySigner(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return consensus.ErrUnknownBlock
	}

	// Retrieve the snapshot needed to verify this header and cache it
	if sb.valsetModule == nil {
		return istanbul.ErrNoEssentialModule
	}
	qualified, err := sb.valsetModule.GetQualifiedValidators(number)
	if err != nil {
		return err
	}

	// resolve the authorization key and check against signers
	signer, err := ecrecover(header)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.
	if !valset.NewAddressSet(qualified).Contains(signer) {
		return istanbul.ErrUnauthorized
	}
	return nil
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (sb *backend) verifyCommittedSeals(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	number := header.Number.Uint64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}
	if sb.valsetModule == nil || sb.govModule == nil {
		return istanbul.ErrNoEssentialModule
	}

	// Retrieve the snapshot needed to verify this header and cache it
	council, err := sb.valsetModule.GetCouncil(number)
	if err != nil {
		return err
	}
	committeeSize := sb.govModule.GetParamSet(number).CommitteeSize
	qualified, err := sb.valsetModule.GetQualifiedValidators(number)
	if err != nil {
		return err
	}
	f := consensuscommon.CalcFaultTolerance(len(qualified), committeeSize)

	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return err
	}
	// The length of Committed seals should be larger than 0
	if len(extra.CommittedSeal) == 0 {
		return istanbul.ErrEmptyCommittedSeals
	}

	// Check whether the committed seals are generated by parent's validators
	councilSet := valset.NewAddressSet(council).Copy()
	validSeal := 0
	proposalSeal := istanbulCore.PrepareCommittedSeal(header.Hash())
	// 1. Get committed seals from current header
	for _, seal := range extra.CommittedSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := cacheSignatureAddresses(proposalSeal, seal)
		if err != nil {
			return istanbul.ErrInvalidSignature
		}
		// Every validator can have only one seal. If more than one seals are signed by a
		// validator, the validator cannot be found and errInvalidCommittedSeals is returned.
		if councilSet.Remove(addr) {
			validSeal++
		} else {
			return istanbul.ErrInvalidCommittedSeals
		}
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= 2*f {
		return istanbul.ErrInvalidCommittedSeals
	}
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (sb *backend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// get parent header and ensure the signer is in parent's validator set
	number := header.Number.Uint64()
	if number == 0 {
		return consensus.ErrUnknownBlock
	}

	// ensure that the blockscore equals to defaultBlockScore
	if header.BlockScore.Cmp(istanbul.DefaultBlockScore) != 0 {
		return consensus.ErrInvalidBlockScore
	}
	return sb.verifySigner(chain, header, nil)
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// unused fields, force to set to empty
	header.Rewardbase = sb.rewardbase
	// use the same blockscore for all blocks
	header.BlockScore = istanbul.DefaultBlockScore

	if chain.Config().IsRandaoForkEnabled(header.Number) {
		prevMixHash := headerMixHash(chain, parent)
		randomReveal, mixHash, err := sb.CalcRandao(header.Number, prevMixHash)
		if err != nil {
			return err
		}
		header.RandomReveal = randomReveal
		header.MixHash = mixHash
	}

	// add qualified validators to extraData's validators section
	if sb.valsetModule == nil {
		return istanbul.ErrNoEssentialModule
	}
	qualified, err := sb.valsetModule.GetQualifiedValidators(number)
	if err != nil {
		return err
	}
	extra, err := prepareExtra(header, qualified)
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp
	header.Time = new(big.Int).Add(parent.Time, new(big.Int).SetUint64(sb.config.BlockPeriod))
	header.TimeFoS = parent.TimeFoS
	if header.Time.Int64() < time.Now().Unix() {
		t := time.Now()
		header.Time = big.NewInt(t.Unix())
		header.TimeFoS = uint8((t.UnixNano() / 1000 / 1000 / 10) % 100)
	}

	for _, module := range sb.consensusModules {
		if err = module.PrepareHeader(header); err != nil {
			return err
		}
	}

	return nil
}

func (sb *backend) Initialize(chain consensus.ChainReader, header *types.Header, state *state.StateDB) {
	// [EIP-2935] stores the parent block hash in the history storage contract
	if chain.Config().IsPragueForkEnabled(header.Number) {
		context := blockchain.NewEVMBlockContext(header, chain, nil)
		vmenv := vm.NewEVM(context, vm.TxContext{}, state, chain.Config(), &vm.Config{})
		blockchain.ProcessParentBlockHash(header, vmenv, state, chain.Config().Rules(header.Number))
	}
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	receipts []*types.Receipt,
) (*types.Block, error) {
	// We can assure that if the magma hard forked block should have the field of base fee
	if chain.Config().IsMagmaForkEnabled(header.Number) {
		if header.BaseFee == nil {
			logger.Error("Magma hard forked block should have baseFee", "blockNum", header.Number.Uint64())
			return nil, errors.New("Invalid Magma block without baseFee")
		}
	} else if header.BaseFee != nil {
		logger.Error("A block before Magma hardfork shouldn't have baseFee", "blockNum", header.Number.Uint64())
		return nil, consensus.ErrInvalidBaseFee
	}

	// TODO-kaiax: Simplify header.Rewardbase setting using kaiax/valset
	// If sb.chain is nil, it means backend is not initialized yet.
	if sb.chain != nil && sb.chain.Config().Istanbul.ProposerPolicy == uint64(istanbul.WeightedRandom) {
		// Determine and update Rewardbase when mining. When mining, state root is not yet determined and will be determined at the end of this Finalize below.
		if common.EmptyHash(header.Root) {
			// TODO-Kaia Let's redesign below logic and remove dependency between block reward and istanbul consensus.
			var (
				blockNum      = header.Number.Uint64()
				rewardAddress = sb.GetRewardAddress(blockNum, sb.address)
			)

			if sb.valsetModule == nil {
				return nil, istanbul.ErrNoEssentialModule
			}
			qualified, err := sb.valsetModule.GetQualifiedValidators(blockNum)
			if err != nil {
				return nil, err
			}

			var logMsg string
			if !valset.NewAddressSet(qualified).Contains(sb.address) || (rewardAddress == common.Address{}) {
				logMsg = "No reward address for nodeValidator. Use node's rewardbase."
			} else {
				// use reward address of current node.
				// only a block made by proposer will be accepted. However, due to round change any node can be the proposer of a block.
				// so need to write reward address of current node to receive reward when it becomes proposer.
				// if current node does not become proposer, the block will be abandoned
				header.Rewardbase = rewardAddress
				logMsg = "Use reward address for nodeValidator."
			}
			logger.Trace(logMsg, "header.Number", header.Number.Uint64(), "node address", sb.address, "rewardbase", header.Rewardbase)
		}
	}

	// TODO-kaiax: Reward distribution must be before KIP160,103. When we moved KIP103,160,Randao,Credit to kaiax modules,
	// this module.FinalizeHeader loop should be at the end of this function, like VerifyHeader and PrepareHeader does.
	for _, module := range sb.consensusModules {
		if err := module.FinalizeHeader(header, state, txs, receipts); err != nil {
			return nil, err
		}
	}

	// RebalanceTreasury can modify the global state (state),
	// so the existing state db should be used to apply the rebalancing result.
	// Only on the KIP-103 or KIP-160 hardfork block, the following logic should be executed
	if chain.Config().IsKIP160ForkBlock(header.Number) || chain.Config().IsKIP103ForkBlock(header.Number) {
		rebalanceResult, err := system.RebalanceTreasury(state, chain, header)
		if err != nil {
			logger.Error("failed to execute treasury rebalancing. State not changed", "err", err)
		} else {
			// Leave the memo in the log for later contract finalization
			isKIP103 := chain.Config().IsKIP103ForkBlock(header.Number) // because memo format differs between KIP-103 and KIP-160
			logger.Info("successfully executed treasury rebalancing", "memo", string(rebalanceResult.Memo(isKIP103)))
		}
	}

	// The Registry contract are installed at RandaoCompatibleBlock with a KIP113 record
	if chain.Config().IsRandaoForkBlock(header.Number) {
		err := system.InstallRegistry(state, chain.Config().RandaoRegistry)
		if err != nil {
			return nil, err
		}
	}

	// Replace the Mainnet credit contract
	if chain.Config().IsKaiaForkBlockParent(header.Number) {
		if chain.Config().ChainID.Uint64() == params.MainnetNetworkId && state.GetCode(system.MainnetCreditAddr) != nil {
			if err := state.SetCode(system.MainnetCreditAddr, system.MainnetCreditV2Code); err != nil {
				return nil, err
			}
			logger.Info("Replaced CypressCredit with CypressCreditV2", "blockNum", header.Number.Uint64())
		}
	}

	header.Root = state.IntermediateRoot(true)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, receipts), nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block) (*types.Block, error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()

	// Bail out if we're unauthorized to sign a block
	if sb.valsetModule == nil {
		return nil, istanbul.ErrNoEssentialModule
	}
	qualified, err := sb.valsetModule.GetQualifiedValidators(number)
	if err != nil {
		return nil, err
	}
	if !valset.NewAddressSet(qualified).Contains(sb.address) {
		return nil, istanbul.ErrUnauthorized
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return nil, consensus.ErrUnknownAncestor
	}
	block, err = sb.updateBlock(block)
	if err != nil {
		return nil, err
	}

	// Initialize seal state; returns nil if block was already committed before Seal started
	commitCh := sb.initSealState(number, block.Hash())
	if commitCh == nil {
		return nil, nil
	}
	defer sb.cleanupSealState()

	// post block into Istanbul engine
	go sb.EventMux().Post(istanbul.RequestEvent{
		Proposal: block,
	})

	logger.Debug("[Seal] Waiting for commitCh", "blockNum", number)

	for {
		select {
		case result := <-commitCh:
			logger.Debug("[Seal] Received from commitCh", "blockNum", number, "resultNil", result == nil)
			if result == nil {
				return nil, nil
			}
			// if the block hash and the hash from channel are the same,
			// return the result. Otherwise, keep waiting the next hash.
			block = types.SetRoundToBlock(block, result.Round)
			if block.Hash() == result.Block.Hash() {
				return result.Block, nil
			}
		}
	}
}

// update timestamp and signature of the block based on its number of transactions
func (sb *backend) updateBlock(block *types.Block) (*types.Block, error) {
	header := block.Header()
	// sign the hash
	seal, err := sb.Sign(sigHash(header).Bytes())
	if err != nil {
		return nil, err
	}

	err = writeSeal(header, seal)
	if err != nil {
		return nil, err
	}

	return block.WithSeal(header), nil
}

func (sb *backend) CalcBlockScore(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(0)
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{
		{
			Namespace: "istanbul",
			Version:   "1.0",
			Service:   &API{chain: chain, istanbul: sb},
			Public:    true,
		},
	}
}

// SetChain sets chain of the Istanbul backend
func (sb *backend) SetChain(chain consensus.ChainReader) {
	sb.chain = chain
}

// RegisterKaiaxModules sets kaiax modules of the Istanbul backend
func (sb *backend) RegisterKaiaxModules(mGov gov.GovModule, mStaking staking.StakingModule, mValset valset.ValsetModule, mRandao randao.RandaoModule) {
	sb.govModule = mGov
	sb.RegisterStakingModule(mStaking)
	sb.valsetModule = mValset
	sb.randaoModule = mRandao

	sb.core.RegisterKaiaxModules(mValset, mGov)
}

func (sb *backend) RegisterStakingModule(module staking.StakingModule) {
	sb.stakingModule = module
}

func (sb *backend) RegisterConsensusModule(modules ...kaiax.ConsensusModule) {
	sb.consensusModules = append(sb.consensusModules, modules...)
}

func (sb *backend) UnregisterConsensusModule(module kaiax.ConsensusModule) {
	for i, m := range sb.consensusModules {
		if m == module {
			sb.consensusModules = append(sb.consensusModules[:i], sb.consensusModules[i+1:]...)
			break
		}
	}
}

func (sb *backend) RegisterTxBundlingModule(modules ...kaiax.TxBundlingModule) {
	sb.txBundlingModules = append(sb.txBundlingModules, modules...)
}

// Start implements consensus.Istanbul.Start
func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return istanbul.ErrStartedEngine
	}

	// clear previous data
	sb.sealMu.Lock()
	sb.proposedBlockHash = common.Hash{}
	sb.sealSkippedNum = 0
	if sb.commitCh != nil {
		close(sb.commitCh)
		sb.commitCh = nil
	}
	sb.sealMu.Unlock()

	sb.SetChain(chain)
	sb.currentBlock = currentBlock
	sb.hasBadBlock = hasBadBlock

	if err := sb.core.Start(); err != nil {
		return err
	}

	sb.coreStarted = true
	return nil
}

// Stop implements consensus.Istanbul.Stop
func (sb *backend) Stop() error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if !sb.coreStarted {
		return istanbul.ErrStoppedEngine
	}
	// Close commitCh to stop any pending Seal() calls
	sb.sealMu.Lock()
	if sb.commitCh != nil {
		close(sb.commitCh)
		sb.commitCh = nil
	}
	sb.sealMu.Unlock()
	if err := sb.core.Stop(); err != nil {
		return err
	}
	sb.coreStarted = false
	return nil
}

// GetConsensusInfo returns consensus information regarding the given block number.
func (sb *backend) GetConsensusInfo(block *types.Block) (consensus.ConsensusInfo, error) {
	blockNumber := block.NumberU64()
	if blockNumber == 0 {
		return consensus.ConsensusInfo{}, nil
	}

	if sb.chain == nil {
		return consensus.ConsensusInfo{}, errNoChainReader
	}

	// get the committers of this block from committed seals
	extra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		return consensus.ConsensusInfo{}, err
	}
	committers, err := RecoverCommittedSeals(extra, block.Hash())
	if err != nil {
		return consensus.ConsensusInfo{}, err
	}

	round := block.Header().Round()
	// get the committee list of this block (blockNumber, round)
	currentProposer, err := sb.valsetModule.GetProposer(blockNumber, uint64(round))
	if err != nil {
		logger.Error("Failed to get proposer.", "blockNum", blockNumber, "round", uint64(round), "err", err)
		return consensus.ConsensusInfo{}, istanbul.ErrInternalError
	}

	var currentCommittee []common.Address
	if sb.valsetModule != nil {
		currentCommittee, _ = sb.valsetModule.GetCommittee(block.NumberU64(), uint64(round))
	}

	// Uncomment to validate if committers are in the committee
	// for _, recovered := range committers {
	// 	found := false
	// 	for _, calculated := range currentCommittee {
	// 		if recovered == calculated {
	// 			found = true
	// 		}
	// 	}
	// 	if !found {
	// 		return consensus.ConsensusInfo{}, errInvalidCommittedSeals
	// 	}
	// }

	// get origin proposer at 0 round.
	var roundZeroProposer *common.Address
	if sb.valsetModule != nil {
		addr, err := sb.valsetModule.GetProposer(blockNumber, 0)
		if err == nil {
			roundZeroProposer = &addr
		}
	}

	cInfo := consensus.ConsensusInfo{
		SigHash:        sigHash(block.Header()),
		Proposer:       currentProposer,
		OriginProposer: roundZeroProposer,
		Committee:      currentCommittee,
		Committers:     committers,
		Round:          round,
	}

	return cInfo, nil
}

func (sb *backend) PurgeCache() {
	// TODO-kaiax: Implement this
}

// FIXME: Need to update this for Istanbul
// sigHash returns the hash which is used as input for the Istanbul
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

// ecrecover extracts the Kaia account address from a signed header.
func ecrecover(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return common.Address{}, err
	}
	addr, err := cacheSignatureAddresses(sigHash(header).Bytes(), istanbulExtra.Seal)
	if err != nil {
		return addr, err
	}

	return addr, nil
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(header *types.Header, vals []common.Address) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < types.IstanbulExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.IstanbulExtraVanity])

	ist := &types.IstanbulExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

// writeSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func writeSeal(h *types.Header, seal []byte) error {
	if len(seal)%types.IstanbulExtraSeal != 0 {
		return istanbul.ErrInvalidSignature
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}

// writeCommittedSeals writes the extra-data field of a block header with given committed seals.
func writeCommittedSeals(h *types.Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return istanbul.ErrInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != types.IstanbulExtraSeal {
			return istanbul.ErrInvalidCommittedSeals
		}
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(istanbulExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}
