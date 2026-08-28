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

package impl

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	mock_randao "github.com/kaiachain/kaia/kaiax/randao/mock"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCheckpointInterval = params.DefaultVRankEpoch / 8

// vrankOpt customizes the VRankModule built by newVRank. Unset fields get sensible defaults.
type vrankOpt func(*vrankOpts)

type vrankOpts struct {
	valset                 *mock_valset.MockValsetModule
	randao                 *mock_randao.MockRandaoModule
	db                     database.Database
	chain                  vrank.Chain // present during Init (catchUp walks it)
	chainAfterInit         vrank.Chain // swapped in AFTER Init (avoids catchUp touching this chain)
	chainAfterInitConflict bool        // set if multiple options try to write chainAfterInit
	hardfork               string      // e.g. "permissionless" (default) or "osaka" (pre-fork)
	start                  bool        // if true, Start + t.Cleanup(Stop); defaults to true

	// Pre-set valset expectations applied AFTER all options are processed, so they target
	// the final opts.valset (which may have been replaced by withValset).
	// has* flags distinguish "option not called" from "option called with nil/zero value".
	validators    []common.Address // GetCommittee(any, any) returns these
	hasValidators bool
	candidates    []common.Address // GetCandTesting(any) returns these
	hasCandidates bool
	proposer      common.Address // GetProposer(any, any) returns this
	hasProposer   bool
}

func withValset(v *mock_valset.MockValsetModule) vrankOpt {
	return func(o *vrankOpts) { o.valset = v }
}

func withRandao(r *mock_randao.MockRandaoModule) vrankOpt {
	return func(o *vrankOpts) { o.randao = r }
}

func withDB(db database.Database) vrankOpt {
	return func(o *vrankOpts) { o.db = db }
}

// withChainForCatchup installs the given chain BEFORE Init. For testing catchUpScoreCaches.
func withChainForCatchup(c vrank.Chain) vrankOpt {
	return func(o *vrankOpts) { o.chain = c }
}

// withHeaders installs the given headers as the chain AFTER Init. Skips catchUpScoreCaches.
func withHeaders(h map[uint64]*types.Header) vrankOpt {
	return func(o *vrankOpts) {
		if o.chainAfterInit != nil {
			o.chainAfterInitConflict = true
		}
		o.chainAfterInit = &testChain{headers: h}
	}
}

// withGenesis installs a minimal chain (block 0, round 0) AFTER Init, so HandleVRankCandidate's
// CurrentHeader() call has a non-nil head. Init runs against an empty chain, so catchUp exits
// early and doesn't consume any GetCandTesting expectations from the caller's valset mock.
func withGenesis() vrankOpt {
	return func(o *vrankOpts) {
		if o.chainAfterInit != nil {
			o.chainAfterInitConflict = true
		}
		o.chainAfterInit = &testChain{headers: map[uint64]*types.Header{0: makeHeaderWithRound(0, 0)}}
	}
}

func withHardfork(hf string) vrankOpt {
	return func(o *vrankOpts) { o.hardfork = hf }
}

func withoutStart() vrankOpt {
	return func(o *vrankOpts) { o.start = false }
}

// withCommittee registers GetCommittee(any, any).Return(addrs, nil).AnyTimes() on the valset
// mock. Convenience for tests that want the same committee across all blocks.
func withCommittee(addrs []common.Address) vrankOpt {
	return func(o *vrankOpts) { o.validators = addrs; o.hasValidators = true }
}

// withCandidates registers GetCandTesting(any).Return(addrs, nil).AnyTimes() on the valset
// mock. Convenience for tests that want the same candidate set across all blocks.
func withCandidates(addrs []common.Address) vrankOpt {
	return func(o *vrankOpts) { o.candidates = addrs; o.hasCandidates = true }
}

// withProposer registers GetProposer(any, any).Return(addr, nil).AnyTimes() on the valset
// mock. Convenience for tests that want the same proposer across all blocks/rounds.
func withProposer(addr common.Address) vrankOpt {
	return func(o *vrankOpts) { o.proposer = addr; o.hasProposer = true }
}

func numToAddr(n int) common.Address {
	return common.BigToAddress(big.NewInt(int64(n)))
}

type CN struct {
	Key         *ecdsa.PrivateKey
	BlsKey      bls.SecretKey
	Addr        common.Address
	Valset      *mock_valset.MockValsetModule
	Randao      *mock_randao.MockRandaoModule
	VRankModule *VRankModule
	DB          database.Database
	sub         chan *vrank.VRankBroadcastEvent
}

func newCN(t *testing.T, options ...vrankOpt) *CN {
	t.Helper()

	ctrl := gomock.NewController(t)
	opts := &vrankOpts{
		valset:   mock_valset.NewMockValsetModule(ctrl),
		randao:   mock_randao.NewMockRandaoModule(ctrl),
		db:       database.NewMemDB(),
		chain:    &testChain{headers: map[uint64]*types.Header{}},
		hardfork: "permissionless",
		start:    true,
	}
	for _, opt := range options {
		opt(opts)
	}
	require.False(t, opts.chainAfterInitConflict,
		"withHeaders and withGenesis are mutually exclusive — only one may be passed to newCN. "+
			"If you need a head plus extra headers, include block 0 in your withHeaders map.")

	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	blsKey, err := bls.DeriveFromECDSA(key)
	require.NoError(t, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)

	// Register on whichever mocks are in use after options (withValset/withRandao may have replaced
	// the defaults, so this must run after the for-loop above).
	opts.randao.EXPECT().GetBlsPubkey(addr, gomock.Any()).Return(blsKey.PublicKey(), nil).AnyTimes()
	if opts.hasValidators {
		opts.valset.EXPECT().GetCommittee(gomock.Any(), gomock.Any()).Return(opts.validators, nil).AnyTimes()
	}
	if opts.hasCandidates {
		opts.valset.EXPECT().GetCandTesting(gomock.Any()).Return(opts.candidates, nil).AnyTimes()
	}
	if opts.hasProposer {
		opts.valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(opts.proposer, nil).AnyTimes()
	}

	module := NewVRankModule()
	require.NoError(t, module.Init(&InitOpts{
		NodeKey:     key,
		BlsKey:      blsKey,
		Randao:      opts.randao,
		Valset:      opts.valset,
		Sealer:      testIstanbulSealer{},
		ChainConfig: params.TestKaiaConfig(opts.hardfork),
		Chain:       opts.chain,
		ChainKv:     opts.db,
	}))

	// Post-Init chain swap: catchUpScoreCaches already ran against the (empty) opts.chain, so
	// assigning a new Chain here won't trigger additional GetCandTesting/GetPFS calls.
	if opts.chainAfterInit != nil {
		module.Chain = opts.chainAfterInit
	}

	sub := make(chan *vrank.VRankBroadcastEvent)
	module.broadcastFeed.Subscribe(sub)

	if opts.start {
		require.NoError(t, module.Start())
		t.Cleanup(module.Stop)
	}

	return &CN{
		Key:         key,
		BlsKey:      blsKey,
		Addr:        addr,
		Valset:      opts.valset,
		Randao:      opts.randao,
		VRankModule: module,
		DB:          opts.db,
		sub:         sub,
	}
}

type testChain struct {
	headers map[uint64]*types.Header
	engine  consensus.Engine
}

func (c *testChain) CurrentHeader() *types.Header {
	var result *types.Header
	for num, h := range c.headers {
		if result == nil || num > result.Number.Uint64() {
			result = h
		}
	}
	return result
}

func (c *testChain) GetHeaderByNumber(number uint64) *types.Header {
	return c.headers[number]
}

var testSealerKey, _ = crypto.GenerateKey()

// newTestSealer returns a real sealer so seal writes and recovery behave as in production.
func newTestSealer() *istanbul.IstanbulSealer {
	return istanbul.NewSealerImpl(testSealerKey)
}

// testIstanbulSealer implements the vrank.Sealer interface for tests by reading the
// round byte that makeHeaderWithRound writes into header.Extra. Seal reads and recovery
// delegate to a real sealer, which is what the parent certificate is verified against.
type testIstanbulSealer struct{}

func (testIstanbulSealer) Round(h *types.Header) (byte, error) {
	if len(h.Extra) < istanbul.IstanbulExtraVanity {
		return 0, errors.New("header extra is too short")
	}
	return h.Extra[istanbul.IstanbulExtraVanity-1], nil
}

func (testIstanbulSealer) RawSeals(h *types.Header) ([]byte, [][]byte, error) {
	return newTestSealer().RawSeals(h)
}

func (testIstanbulSealer) RecoverCommitters(hash common.Hash, round byte, seals [][]byte) ([]common.Address, error) {
	return newTestSealer().RecoverCommitters(hash, round, seals)
}

func (testIstanbulSealer) Quorum(blockNum uint64, qualifiedLen, committeeSize int) int {
	return newTestSealer().Quorum(blockNum, qualifiedLen, committeeSize)
}

// Author returns the address withAuthor stashed in Rewardbase, so tests can make the author
// differ from the round's proposer the way a hash-locked block does.
func (testIstanbulSealer) Author(h *types.Header) (common.Address, error) {
	if len(h.Extra) < istanbul.IstanbulExtraVanity {
		return common.Address{}, errors.New("header extra is too short")
	}
	return h.Rewardbase, nil
}

// withAuthor sets the header's author for testIstanbulSealer.Author.
func withAuthor(h *types.Header, author common.Address) *types.Header {
	h.Rewardbase = author
	return h
}

func makeHeaderWithRound(number uint64, round int64) *types.Header {
	h := &types.Header{
		Number:     big.NewInt(int64(number)),
		Time:       big.NewInt(0),
		BlockScore: big.NewInt(1),
	}
	sealer := newTestSealer()
	// A well-formed Istanbul extra, so the sealer can parse this header.
	if err := sealer.WriteValidators(h, nil); err != nil {
		panic(err)
	}
	sealer.WriteRound(h, round)
	return h
}

// withParentAt seeds the round-0 parent that PrepareHeader(num) reads for the certificate.
func withParentAt(num uint64) vrankOpt {
	return withHeaders(map[uint64]*types.Header{num - 1: makeHeaderWithRound(num-1, 0)})
}

// makeHeaderWithParentRound declares the parent's round without a certificate: pfReport reads
// the round, only VerifyHeader checks the seals.
func makeHeaderWithParentRound(number uint64, parentRound uint8) *types.Header {
	h := makeHeaderWithRound(number, 0)
	encoded, err := vrank.EncodeVRank(vrank.VRankPayload{ParentRound: parentRound})
	if err != nil {
		panic(err)
	}
	h.VRank = encoded
	return h
}

func mustPop(t *testing.T, sub chan *vrank.VRankBroadcastEvent) *vrank.VRankBroadcastEvent {
	select {
	case req := <-sub:
		return req
	case <-time.After(2 * time.Second):
		t.Fatal("should broadcast")
	}
	return nil
}

func mustNotPop(t *testing.T, sub chan *vrank.VRankBroadcastEvent) *vrank.VRankBroadcastEvent {
	select {
	case <-sub:
		t.Fatal("should not broadcast")
	default:
	}
	return nil
}

// TestInit_CatchUpFromCheckpoint verifies that Init replays only the tail blocks after the
// stored DB checkpoint.
func TestInit_CatchUpFromCheckpoint(t *testing.T) {
	checkpoint := testCheckpointInterval // must be a testCheckpointInterval multiple
	P1, P2, C1 := numToAddr(1), numToAddr(2), numToAddr(10)

	// Tail block (checkpoint+1): round=1, author=P2, parent round=1, cfReport=[C1].
	//   - PFS:      pfReport=[P1], the parent's round-0 proposer
	//   - cpMatrix: reporter=P2, the block's author
	tail := withAuthor(makeHeaderWithRound(checkpoint+1, 1), P2)
	tailVRank, err := vrank.EncodeVRank(vrank.VRankPayload{Report: []common.Address{C1}, ParentRound: 1})
	require.NoError(t, err)
	tail.VRank = tailVRank
	headers := map[uint64]*types.Header{
		checkpoint:     makeHeaderWithRound(checkpoint, 0),
		checkpoint + 1: tail,
	}

	db := database.NewMemDB()
	WriteCheckpoint(db, checkpoint, map[common.Address]uint64{P1: 1}, vrank.CPMatrix{C1: {P1: 1}})
	WriteLastCheckpoint(db, checkpoint)

	valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
	// round=0 → P1 (pfReport proposer); round=1 → P2 (cfReport reporter).
	// Per-round variation, so this stays as direct EXPECT (helper takes a single addr).
	valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ uint64, round uint64) (common.Address, error) {
			if round == 1 {
				return P2, nil
			}
			return P1, nil
		},
	).AnyTimes()

	// withChain at Init time so catchUpScoreCaches actually walks these headers — that's
	// what the test is verifying. (Regular tests use withHeaders, which installs after Init.)
	// withCandidates(nil) covers catchUp's GetCandTesting call at head.
	module := newCN(t,
		withValset(valset),
		withCandidates(nil),
		withChainForCatchup(&testChain{headers: headers}),
		withDB(db),
		withoutStart(),
	).VRankModule

	// pfsCache: checkpoint seeded P1=1; tail block adds 1 more (round=1 → pfReport=[P1]).
	pfsCached, ok := module.pfsCache.Get(checkpoint + 1)
	require.True(t, ok, "pfsCache should be populated at head")
	pfs := pfsCached.(map[common.Address]uint64)
	assert.Equal(t, uint64(2), pfs[P1], "P1 score: 1 from checkpoint + 1 from tail block")

	// cpMatrixCache: checkpoint seeded C1:{P1:1}; tail block adds C1:{P2:1} (reporter=P2 at round=1).
	cpCached, ok := module.cpMatrixCache.Get(checkpoint + 1)
	require.True(t, ok, "cpMatrixCache should be populated at head")
	cpMatrix := cpCached.(vrank.CPMatrix)
	assert.Equal(t, uint64(1), cpMatrix[C1][P1], "P1 contribution should carry over from checkpoint")
	assert.Equal(t, uint64(1), cpMatrix[C1][P2], "P2 should be credited as reporter in tail block")
}

func TestCheckpointRoundTrip_PreservesZeroScores(t *testing.T) {
	checkpoint := testCheckpointInterval
	P1, P2, C1, C2 := numToAddr(1), numToAddr(2), numToAddr(10), numToAddr(11)

	db := database.NewMemDB()
	pfsIn := map[common.Address]uint64{
		P1: 3,
		P2: 0,
	}
	cpMatrixIn := vrank.CPMatrix{
		C1: {P1: 1, P2: 0},
		C2: {P1: 0, P2: 2},
	}
	WriteCheckpoint(db, checkpoint,
		pfsIn,
		cpMatrixIn,
	)

	pfsOut := ReadCheckpointPFS(db, checkpoint)
	cpMatrixOut := ReadCheckpointCPMatrix(db, checkpoint)
	require.NotNil(t, pfsOut)
	require.NotNil(t, cpMatrixOut)
	assert.Equal(t, pfsIn, pfsOut)
	assert.Equal(t, cpMatrixIn, cpMatrixOut)
}

func TestVRankModule_RestartAfterStop(t *testing.T) {
	module := newCN(t, withoutStart()).VRankModule // withoutStart so the test drives Start/Stop/Start manually.

	sink := make(chan *vrank.VRankBroadcastEvent, 1)
	sub := module.SubscribeVRank(sink)
	defer sub.Unsubscribe()

	require.NoError(t, module.Start())
	module.Stop()
	module.Stop() // Stop must be idempotent

	require.NoError(t, module.Start())
	module.broadcast([]common.Address{numToAddr(0)}, "test")

	select {
	case req := <-sink:
		assert.Equal(t, []common.Address{numToAddr(0)}, req.Targets)
		assert.Equal(t, "test", req.Msg)
	case <-time.After(time.Second):
		t.Fatal("broadcast did not resume after restart")
	}
}

// strictSealer rejects a header with no author stashed, the way IstanbulSealer.Author rejects a
// header with no author seal. Genesis is the case that matters.
type strictSealer struct{ testIstanbulSealer }

func (strictSealer) Author(h *types.Header) (common.Address, error) {
	if h.Rewardbase == (common.Address{}) {
		return common.Address{}, errors.New("no author seal")
	}
	return h.Rewardbase, nil
}
