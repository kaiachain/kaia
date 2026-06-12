// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

// Package kaiabft implements a BFT consensus engine with speculative execution.
//
// KaiaBFT is wire-compatible with Istanbul BFT — they use the same message
// format, same extra-data layout, and can reach consensus together in a mixed
// network. The key enhancement is that non-proposer validators speculatively
// execute the proposed block upon receiving a pre-prepare, populating a cache
// that InsertChain can use to skip re-execution.
package kaiabft

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/networks/rpc"
)

const (
	inmemoryPeers    = 200
	inmemoryMessages = 4096
	fetcherID        = "kaiabft"
	chainInitTimeout = 30 * time.Second
)

var (
	logger = log.NewModuleLogger(log.ConsensusIstanbulBackend)

	errStoppedEngine = errors.New("kaiabft: stopped engine")
	errStartedEngine = errors.New("kaiabft: started engine")
	errDecodeFailed  = errors.New("kaiabft: fail to decode message")
	errNoChainReader = errors.New("kaiabft: chain is nil, --mine option might be missing")
	errNoModule      = errors.New("kaiabft: essential kaiax module not registered")
	errInvalidPeer   = errors.New("kaiabft: invalid peer address")
)

// Opts bundles the inputs needed to instantiate a kaiabft consensus engine.
type Opts struct {
	Timeout    uint64            // Round-change timeout in milliseconds
	PrivateKey *ecdsa.PrivateKey // Consensus message signing key
	NodeType   common.ConnType   // CN, PN, or EN
	Sealer     consensus.Sealer  // Shared BFT seal operations
}

// New creates a new kaiabft consensus engine.
func New(opts *Opts) consensus.Engine {
	recentMessages, _ := lru.NewARC(inmemoryPeers)
	knownMessages, _ := lru.NewARC(inmemoryMessages)
	b := &backend{
		timeout:        opts.Timeout,
		privateKey:     opts.PrivateKey,
		address:        crypto.PubkeyToAddress(opts.PrivateKey.PublicKey),
		sealer:         opts.Sealer,
		nodetype:       opts.NodeType,
		eventMux:       new(event.TypeMux),
		logger:         logger.NewWith(),
		commitCh:       make(chan *types.Result, 1),
		recentMessages: recentMessages,
		knownMessages:  knownMessages,
		chainInitCh:    make(chan struct{}),
	}
	b.currentView.Store(&bft.View{Sequence: big.NewInt(0), Round: big.NewInt(0)})
	return b
}

// backend implements consensus.Engine and consensus.Handler.
type backend struct {
	timeout    uint64
	privateKey *ecdsa.PrivateKey
	address    common.Address
	sealer     consensus.Sealer
	nodetype   common.ConnType
	logger     log.Logger

	chain    consensus.ChainReader
	executor consensus.Executor
	eventMux *event.TypeMux

	// Core state machine
	machine     *machine
	coreStarted bool
	coreMu      sync.RWMutex

	// Seal synchronization
	commitCh          chan *types.Result
	proposedBlockHash common.Hash
	sealMu            sync.Mutex
	sealSkippedNum    uint64

	// P2P message caching
	broadcaster    consensus.Broadcaster
	recentMessages *lru.ARCCache
	knownMessages  *lru.ARCCache
	currentView    atomic.Value // *bft.View

	// Kaiax modules
	valsetModule valset.ValsetModule
	govModule    gov.GovModule

	// Peer registration gating
	chainInitCh   chan struct{}
	chainInitOnce sync.Once

	// Speculative execution (cache set in Start from chain type assertion;
	// cancel/wg track the in-flight goroutine launched by the machine).
	specCache  *blockchain.SpeculativeResultCache
	specCancel context.CancelFunc
	specWg     sync.WaitGroup
	specMu     sync.Mutex
}

// ---------------------------------------------------------------------------
// consensus.Engine implementation
// ---------------------------------------------------------------------------

func (b *backend) RegisterKaiaxModules(mGov gov.GovModule, mValset valset.ValsetModule) {
	b.valsetModule = mValset
	b.govModule = mGov
}

func (b *backend) Start(chain consensus.ChainReader, executor consensus.Executor) error {
	b.coreMu.Lock()
	defer b.coreMu.Unlock()
	if b.coreStarted {
		return errStartedEngine
	}

	// Clear previous seal state.
	b.sealMu.Lock()
	b.proposedBlockHash = common.Hash{}
	b.sealSkippedNum = 0
	if b.commitCh != nil {
		close(b.commitCh)
	}
	b.commitCh = nil
	b.sealMu.Unlock()

	defer b.signalPeerRegistrable()

	b.chain = chain
	b.executor = executor

	// Grab the speculative cache from the concrete BlockChain if available.
	if bc, ok := chain.(*blockchain.BlockChain); ok {
		b.specCache = bc.SpeculativeCache()
	}

	b.machine = newMachine(b)
	if err := b.machine.start(); err != nil {
		return err
	}

	b.coreStarted = true
	return nil
}

func (b *backend) Stop() error {
	b.coreMu.Lock()
	defer b.coreMu.Unlock()
	b.signalPeerRegistrable()
	if !b.coreStarted {
		return errStoppedEngine
	}
	b.sealMu.Lock()
	if b.commitCh != nil {
		close(b.commitCh)
		b.commitCh = nil
	}
	b.sealMu.Unlock()
	b.machine.stop()
	b.cancelSpeculativeExecution()
	b.specWg.Wait()
	b.coreStarted = false
	return nil
}

func (b *backend) SubmitTransactions(txs *types.TransactionsByPriceAndNonce, statedb *state.StateDB, header *types.Header, mux *event.TypeMux) <-chan *consensus.ExecutionResult {
	resultCh := make(chan *consensus.ExecutionResult, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("SubmitTransactions panic", "err", r)
				resultCh <- nil
			}
		}()

		if cv, ok := b.currentView.Load().(*bft.View); ok && cv != nil &&
			cv.Sequence.Cmp(header.Number) == 0 {
			proposer, err := b.valsetModule.GetProposer(cv.Sequence.Uint64(), cv.Round.Uint64())
			if err != nil {
				logger.Warn("Failed to resolve proposer for local block execution",
					"number", header.Number.Uint64(), "viewSeq", cv.Sequence.Uint64(),
					"viewRound", cv.Round.Uint64(), "self", b.address, "err", err)
			} else if proposer != b.address {
				logger.Debug("Skipping local block execution on non-proposer",
					"number", header.Number.Uint64(), "proposer", proposer)
				resultCh <- nil
				return
			}
		}
		// If currentView is stale or unavailable, fall through to the legacy
		// full execution path as a defensive fallback.

		validators, err := b.valsetModule.GetQualifiedValidators(header.Number.Uint64())
		if err != nil {
			resultCh <- nil
			return
		}
		if err := b.sealer.WriteValidators(header, validators); err != nil {
			resultCh <- nil
			return
		}

		if err := b.executor.ResetWithState(statedb, header); err != nil {
			resultCh <- nil
			return
		}

		executeStart := time.Now()
		result, err := b.executor.Execute(txs, mux)
		if err != nil {
			resultCh <- nil
			return
		}
		result.ExecuteTime = time.Since(executeStart)

		finalizeStart := time.Now()
		block, err := b.executor.FinalizeState(result)
		if err != nil {
			resultCh <- nil
			return
		}
		result.FinalizeTime = time.Since(finalizeStart)
		result.Block = block

		// Log block preparation completion (all validators log this, before seal)
		logger.Info("Prepared new block",
			"number", result.Block.Number(),
			"hash", result.Block.Hash(),
			"txs", len(result.Txs),
			"elapsed", common.PrettyDuration(result.ExecuteTime+result.FinalizeTime),
			"executeTime", common.PrettyDuration(result.ExecuteTime),
			"finalizeTime", common.PrettyDuration(result.FinalizeTime))

		sealStart := time.Now()
		sealedBlock, err := b.seal(block)
		result.SealTime = time.Since(sealStart)

		if err != nil {
			logger.Error("Failed to seal block", "err", err, "sealTime", result.SealTime)
			resultCh <- nil
			return
		}
		if sealedBlock == nil {
			logger.Debug("Seal skipped, not the proposer", "number", block.Number(), "sealTime", result.SealTime)
			resultCh <- nil
			return
		}

		logger.Info("Successfully sealed new block", "number", sealedBlock.Number(), "hash", sealedBlock.Hash(), "sealTime", result.SealTime)
		result.Block = sealedBlock
		resultCh <- result
	}()

	return resultCh
}

func (b *backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{
		{
			Namespace: "istanbul",
			Version:   "1.0",
			Service:   &api{chain: chain, backend: b},
			Public:    true,
		},
	}
}

func (b *backend) PurgeCache() {}

func (b *backend) SubscribeNewSequence() *event.TypeMuxSubscription {
	return b.eventMux.Subscribe(consensus.NewSequenceEvent{})
}

// ---------------------------------------------------------------------------
// consensus.Handler implementation
// ---------------------------------------------------------------------------

func (b *backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	b.coreMu.Lock()
	defer b.coreMu.Unlock()

	if msg.Code != consensus.ConsensusMsgCode {
		return false, nil
	}

	if !b.coreStarted {
		return true, errStoppedEngine
	}

	var cmsg bft.ConsensusMsg
	if err := msg.Decode(&cmsg); err != nil {
		return true, errDecodeFailed
	}
	data := cmsg.Payload
	hash := bft.RLPHash(data)

	// Mark peer's message.
	var m *lru.ARCCache
	ms, ok := b.recentMessages.Get(addr)
	if ok {
		m, _ = ms.(*lru.ARCCache)
	} else {
		m, _ = lru.NewARC(inmemoryMessages)
		b.recentMessages.Add(addr, m)
	}
	m.Add(hash, true)

	// Deduplicate.
	if _, ok := b.knownMessages.Get(hash); ok {
		return true, nil
	}
	b.knownMessages.Add(hash, true)

	go b.eventMux.Post(messageEvent{
		Payload: data,
		Hash:    cmsg.PrevHash,
	})

	return true, nil
}

func (b *backend) NewChainHead() error {
	b.coreMu.RLock()
	defer b.coreMu.RUnlock()
	if !b.coreStarted {
		return errStoppedEngine
	}
	go b.eventMux.Post(chainHeadEvent{})
	return nil
}

func (b *backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	b.broadcaster = broadcaster
	if b.nodetype == common.CONSENSUSNODE {
		b.broadcaster.RegisterValidator(common.CONSENSUSNODE, b)
	}
}

// ---------------------------------------------------------------------------
// Peer validation (p2p.PeerTypeValidator)
// ---------------------------------------------------------------------------

func (b *backend) ValidatePeerType(addr common.Address) error {
	select {
	case <-b.chainInitCh:
	case <-time.After(chainInitTimeout):
		return errNoChainReader
	}
	if b.chain == nil {
		return errNoChainReader
	}
	if b.valsetModule == nil {
		return errNoModule
	}
	council, err := b.valsetModule.GetCouncil(b.chain.CurrentHeader().Number.Uint64() + 1)
	if err != nil {
		return err
	}
	if valset.NewAddressSet(council).Contains(addr) {
		return nil
	}
	return errInvalidPeer
}

func (b *backend) signalPeerRegistrable() {
	b.chainInitOnce.Do(func() { close(b.chainInitCh) })
}

// ---------------------------------------------------------------------------
// Seal — triggers BFT consensus on a proposed block
// ---------------------------------------------------------------------------

func (b *backend) seal(block *types.Block) (*types.Block, error) {
	header := block.Header()
	number := header.Number.Uint64()

	block, err := b.updateBlock(block)
	if err != nil {
		return nil, err
	}

	commitCh := b.initSealState(number, block.Hash())
	if commitCh == nil {
		return nil, nil
	}
	defer b.cleanupSealState(commitCh)

	go b.eventMux.Post(requestEvent{Proposal: block})

	for result := range commitCh {
		if result == nil {
			return nil, nil
		}
		header := block.Header()
		b.sealer.WriteRound(header, result.Round)
		block = block.WithSeal(header)
		if block.Hash() == result.Block.Hash() {
			return result.Block, nil
		}
	}
	return nil, nil
}

func (b *backend) updateBlock(block *types.Block) (*types.Block, error) {
	header := block.Header()
	if header != nil && header.Number != nil && header.Number.Uint64() > 0 {
		parent := b.chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
		if parent == nil || parent.Hash() != header.ParentHash || parent.Number.Uint64()+1 != header.Number.Uint64() {
			return nil, consensus.ErrUnknownAncestor
		}
	}
	if b.valsetModule == nil {
		return nil, errNoModule
	}
	qualified, err := b.valsetModule.GetQualifiedValidators(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	if !valset.NewAddressSet(qualified).Contains(b.address) {
		return nil, consensus.ErrUnauthorized
	}
	authorSeal, err := b.sealer.MakeAuthorSeal(header)
	if err != nil {
		return nil, err
	}
	if err := b.sealer.WriteAuthorSeal(header, authorSeal); err != nil {
		return nil, err
	}
	return block.WithSeal(header), nil
}

func (b *backend) initSealState(number uint64, blockHash common.Hash) chan *types.Result {
	b.sealMu.Lock()
	defer b.sealMu.Unlock()
	if b.sealSkippedNum == number {
		b.sealSkippedNum = 0
		return nil
	}
	// A prior in-flight seal() for this sequence (e.g. a round-change rebuild
	// superseding an earlier build) is still blocked on the old channel. Release
	// it with a nil result so its goroutine returns instead of leaking. The
	// channel is buffered (cap 1) and the old sealer is its only reader, so this
	// never blocks: it either hands nil to the waiting sealer or fills the buffer
	// of a sealer that has already returned.
	if b.commitCh != nil {
		b.commitCh <- nil
	}
	b.proposedBlockHash = blockHash
	b.commitCh = make(chan *types.Result, 1)
	return b.commitCh
}

func (b *backend) cleanupSealState(ch chan *types.Result) {
	b.sealMu.Lock()
	defer b.sealMu.Unlock()
	// Only clear if we still own the current channel: a superseding seal()
	// (see initSealState) may have already replaced it, and clearing then would
	// wipe the live seal's state.
	if b.commitCh != ch {
		return
	}
	b.proposedBlockHash = common.Hash{}
	b.commitCh = nil
}

// ---------------------------------------------------------------------------
// Verify, Commit, Broadcast — called from the state machine
// ---------------------------------------------------------------------------

func (b *backend) verify(proposal bft.Proposal) (time.Duration, error) {
	block, ok := proposal.(*types.Block)
	if !ok {
		return 0, errors.New("kaiabft: invalid proposal type")
	}
	if b.chain.HasBadBlock(block.Hash()) {
		return 0, blockchain.ErrBlacklistedHash
	}
	txnHash := types.DeriveTransactionsRoot(block.Transactions(), block.Number())
	if txnHash != block.Header().TxHash {
		return 0, errors.New("kaiabft: mismatch transaction hashes")
	}
	for _, tx := range block.Transactions() {
		if tx.Type() == types.TxTypeEthereumBlob {
			sidecar := tx.BlobTxSidecar()
			if sidecar == nil {
				return 0, errors.New("kaiabft: no blob sidecar for blob tx")
			}
			if err := sidecar.ValidateWithBlobHashes(tx.BlobHashes()); err != nil {
				return 0, err
			}
		}
	}
	err := b.chain.ValidateHeader(block.Header())
	if err == nil || err.Error() == "zero committed seals" {
		return 0, nil
	}
	if err == consensus.ErrFutureBlock {
		return time.Until(time.Unix(block.Header().Time.Int64(), 0)), consensus.ErrFutureBlock
	}
	return 0, err
}

func (b *backend) commit(proposal bft.Proposal, seals [][]byte) error {
	block, ok := proposal.(*types.Block)
	if !ok {
		return errors.New("kaiabft: invalid proposal type")
	}
	proposalHash := proposal.Hash()
	h := block.Header()
	round := b.currentView.Load().(*bft.View).Round.Int64()
	b.sealer.WriteRound(h, round)
	if err := b.sealer.WriteCommittedSeals(h, seals); err != nil {
		return err
	}
	block = block.WithSeal(h)

	b.logger.Info("Committed", "number", proposal.Number().Uint64(), "hash", proposalHash, "address", b.address)

	b.sealMu.Lock()
	if b.proposedBlockHash == proposalHash {
		if b.commitCh != nil {
			b.commitCh <- &types.Result{Block: block, Round: round}
		}
		b.sealMu.Unlock()
		return nil
	}
	if b.commitCh != nil {
		b.commitCh <- nil
	} else {
		b.sealSkippedNum = proposal.Number().Uint64()
	}
	b.sealMu.Unlock()

	if b.broadcaster != nil {
		b.broadcaster.Enqueue(fetcherID, block)
	}
	return nil
}

func (b *backend) broadcast(prevHash common.Hash, payload []byte) error {
	go b.eventMux.Post(messageEvent{Hash: prevHash, Payload: payload})
	return nil
}

func (b *backend) gossipSubPeer(prevHash common.Hash, payload []byte) {
	targets := b.getTargetReceivers()
	if targets == nil {
		return
	}
	hash := bft.RLPHash(payload)
	b.knownMessages.Add(hash, true)

	if b.broadcaster == nil || len(targets) == 0 {
		return
	}
	ps := b.broadcaster.FindCNPeers(targets)
	for addr, p := range ps {
		ms, ok := b.recentMessages.Get(addr)
		var m *lru.ARCCache
		if ok {
			m, _ = ms.(*lru.ARCCache)
			if _, k := m.Get(hash); k {
				continue
			}
		} else {
			m, _ = lru.NewARC(inmemoryMessages)
		}
		m.Add(hash, true)
		b.recentMessages.Add(addr, m)
		go p.Send(consensus.ConsensusMsgCode, &bft.ConsensusMsg{PrevHash: prevHash, Payload: payload})
	}
}

func (b *backend) getTargetReceivers() map[common.Address]bool {
	if b.valsetModule == nil {
		return nil
	}
	cv, ok := b.currentView.Load().(*bft.View)
	if !ok {
		return nil
	}
	num := cv.Sequence.Uint64()
	targets := make(map[common.Address]bool)
	for i := range 2 {
		round := cv.Round.Uint64() + uint64(i)
		committee, err := b.valsetModule.GetCommittee(num, round)
		if err != nil {
			return nil
		}
		committeeSet := valset.NewAddressSet(committee)
		if i == 0 && !committeeSet.Contains(b.address) {
			return nil
		}
		for _, val := range committee {
			if val != b.address {
				targets[val] = true
			}
		}
	}
	return targets
}

func (b *backend) sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256([]byte(data))
	return crypto.Sign(hashData, b.privateKey)
}

func (b *backend) lastProposal() (bft.Proposal, common.Address) {
	block := b.chain.CurrentBlock()
	var proposer common.Address
	if block.Number().Cmp(common.Big0) > 0 {
		var err error
		proposer, err = b.sealer.Author(block.Header())
		if err != nil {
			b.logger.Error("Failed to get block proposer", "err", err)
			return nil, common.Address{}
		}
	}
	return block, proposer
}

func (b *backend) hasBadProposal(hash common.Hash) bool {
	if b.chain == nil {
		return false
	}
	return b.chain.HasBadBlock(hash)
}

func (b *backend) hasProposal(hash common.Hash, number *big.Int) bool {
	return b.chain.GetHeader(hash, number.Uint64()) != nil
}

// ---------------------------------------------------------------------------
// Speculative execution
// ---------------------------------------------------------------------------

// startSpeculativeExecution validates the proposed block in a background
// goroutine and populates the speculative cache so InsertChain can skip
// re-execution. Cancels any prior in-flight execution. Caller (the machine)
// owns the decision of when to trigger; this method owns the execution.
func (b *backend) startSpeculativeExecution(proposal bft.Proposal) {
	block, ok := proposal.(*types.Block)
	if !ok || b.specCache == nil || b.executor == nil {
		return
	}

	blockHash := block.Hash()
	// Re-entry for the same proposal (e.g. future-block retry) keeps the
	// in-flight or completed execution instead of restarting it.
	if b.specCache.HasUsable(blockHash) {
		return
	}

	b.specMu.Lock()
	if b.specCancel != nil {
		b.specCancel()
	}
	ctx, cancel := context.WithCancel(context.Background())
	b.specCancel = cancel
	b.specMu.Unlock()

	entry := b.specCache.Reserve(blockHash)

	signer := types.MakeSigner(b.chain.Config(), block.Number())
	executor := b.executor.Clone()

	parentHeader := b.chain.GetHeader(block.ParentHash(), block.NumberU64()-1)
	if parentHeader == nil {
		entry.Complete(nil, errors.New("parent header not found"))
		cancel()
		return
	}

	// Adopt pool-known senders; kick async ecrecover for the rest. Runs
	// synchronously so execution and prefetch always see warmed senders.
	blockchain.WarmSenders(signer, block, b.chain.TxLookup())

	// Warm trie-node cache for spec-exec; ctx ties prefetch to this round.
	b.specWg.Add(1)
	go func() {
		defer b.specWg.Done()
		blockchain.PrefetchBlockState(ctx, b.chain, parentHeader.Root, block.NumberU64(), block.Transactions(), signer)
	}()

	b.specWg.Add(1)
	go func() {
		defer b.specWg.Done()
		defer cancel()

		// PrunableStateAt marks obsolete trie nodes for live pruning the same
		// way InsertChain's normal path does; falls back to StateAt when
		// pruning is disabled.
		parentState, err := b.chain.PrunableStateAt(parentHeader.Root, parentHeader.Number.Uint64())
		if err != nil {
			entry.Complete(nil, err)
			return
		}

		if err := executor.ResetWithState(parentState, block.Header()); err != nil {
			entry.Complete(nil, err)
			return
		}

		select {
		case <-ctx.Done():
			entry.Complete(nil, ctx.Err())
			return
		default:
		}

		// ProcessBlock delegates to StateProcessor.Process() which already
		// applies rewards and computes the state root — do NOT call
		// FinalizeState again.
		result, err := executor.ProcessBlock(block.Transactions())
		if err != nil {
			entry.Complete(nil, err)
			return
		}

		// Pre-compute bloom + receipt hash so InsertChain can validate the
		// header without re-deriving them; both are O(n) over receipts and
		// dominate ValidateState time.
		bloom := types.CreateBloom(result.Receipts)
		receiptHash := types.DeriveReceiptsRoot(result.Receipts, block.Number())

		entry.Complete(&blockchain.SpeculativeResult{
			State:            result.State,
			Receipts:         result.Receipts,
			Logs:             result.Logs,
			UsedGas:          result.UsedGas,
			InternalTxTraces: result.InternalTxTraces,
			ProcessStats: blockchain.ProcessStats{
				BeforeApplyTxs: result.BeforeApplyTxs,
				AfterApplyTxs:  result.AfterApplyTxs,
				AfterFinalize:  result.AfterFinalize,
			},
			Bloom:       bloom,
			ReceiptHash: receiptHash,
		}, nil)

		logger.Info("Speculative execution completed",
			"number", block.NumberU64(), "hash", blockHash, "txs", len(block.Transactions()))
	}()
}

// cancelSpeculativeExecution signals any in-flight speculative execution to
// abort. Safe to call when nothing is running.
func (b *backend) cancelSpeculativeExecution() {
	b.specMu.Lock()
	defer b.specMu.Unlock()
	if b.specCancel != nil {
		b.specCancel()
		b.specCancel = nil
	}
}

// ---------------------------------------------------------------------------
// Internal event types used by the state machine
// ---------------------------------------------------------------------------

type requestEvent struct {
	Proposal bft.Proposal
}

type messageEvent struct {
	Hash    common.Hash
	Payload []byte
}

type chainHeadEvent struct{}

type backlogEvent struct {
	src  common.Address
	msg  *bft.Message
	Hash common.Hash
}

type timeoutEvent struct {
	nextView *bft.View
}
