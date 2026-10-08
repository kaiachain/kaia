// Copyright 2024 The Kaia Authors
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
package core

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	mock_istanbul "github.com/kaiachain/kaia/consensus/istanbul/mocks"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/kaiax/gov"
	mock_gov "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMockBackend create a mock-backend and mock valset/gov modules initialized with default values.
// Caller must call istCore.RegisterKaiaxModules(mockValset, mockGov) after New().
func newMockBackend(t *testing.T, validatorAddrs []common.Address, permissionless bool) (*mock_istanbul.MockBackend, *gomock.Controller, *valset_mock.MockValsetModule, *mock_gov.MockGovModule) {
	committeeSize := uint64(len(validatorAddrs) / 3)

	istExtra := &istanbul.IstanbulExtra{
		Validators:    validatorAddrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}
	extra, err := rlp.EncodeToBytes(istExtra)
	if err != nil {
		t.Fatal(err)
	}

	initBlock := types.NewBlockWithHeader(&types.Header{
		ParentHash: common.Hash{},
		Number:     common.Big0,
		GasUsed:    0,
		Extra:      append(make([]byte, istanbul.IstanbulExtraVanity), extra...),
		Time:       new(big.Int).SetUint64(1234),
		BlockScore: common.Big0,
	})

	eventMux := new(event.TypeMux)

	mockCtrl := gomock.NewController(t)
	mockValsetModule := valset_mock.NewMockValsetModule(mockCtrl)
	mockGovModule := mock_gov.NewMockGovModule(mockCtrl)
	mockBackend := mock_istanbul.NewMockBackend(mockCtrl)

	// Valset module: used by core startRound (no BlockValSet)
	mockValsetModule.EXPECT().GetCouncil(gomock.Any()).Return(validatorAddrs, nil).AnyTimes()
	mockValsetModule.EXPECT().GetDemotedValidators(gomock.Any()).Return([]common.Address{}, nil).AnyTimes()
	mockValsetModule.EXPECT().GetCommittee(gomock.Any(), gomock.Any()).DoAndReturn(
		func(num uint64, round uint64) ([]common.Address, error) {
			if round == 2 {
				return validatorAddrs[2 : committeeSize+2], nil
			}
			return validatorAddrs[0:committeeSize], nil
		},
	).AnyTimes()
	mockValsetModule.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(validatorAddrs[0], nil).AnyTimes()

	// Gov module: required by core registration checks.
	mockGovModule.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{CommitteeSize: committeeSize}).AnyTimes()

	// Consider the last proposal is "initBlock" and the owner of mockBackend is validatorAddrs[0]
	mockBackend.EXPECT().Address().Return(validatorAddrs[0]).AnyTimes()
	sealerKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	mockBackend.EXPECT().Sealer().Return(istanbul.NewSealerImpl(sealerKey)).AnyTimes()
	mockBackend.EXPECT().LastProposal().Return(initBlock, validatorAddrs[0]).AnyTimes()
	mockBackend.EXPECT().NodeType().Return(common.CONSENSUSNODE).AnyTimes()
	mockBackend.EXPECT().IsPermissionlessAt(gomock.Any()).Return(permissionless).AnyTimes()

	// Set an eventMux in which istanbul core will subscribe istanbul events
	mockBackend.EXPECT().EventMux().Return(eventMux).AnyTimes()

	// Just for bypassing an unused function
	mockBackend.EXPECT().SetCurrentView(gomock.Any()).Return().AnyTimes()

	// Always return nil for broadcasting related functions
	mockBackend.EXPECT().Sign(gomock.Any()).Return(nil, nil).AnyTimes()
	mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	mockBackend.EXPECT().GossipSubPeer(gomock.Any(), gomock.Any()).Return().AnyTimes()

	// Verify checks whether the proposal of the preprepare message is a valid block. Consider it valid.
	mockBackend.EXPECT().Verify(gomock.Any()).Return(time.Duration(0), nil).AnyTimes()

	return mockBackend, mockCtrl, mockValsetModule, mockGovModule
}

func TestGetRoundCommitteeStateUsesReturnedCommitteeSize(t *testing.T) {
	validatorAddrs := []common.Address{
		common.HexToAddress("0x0001"),
		common.HexToAddress("0x0002"),
		common.HexToAddress("0x0003"),
		common.HexToAddress("0x0004"),
		common.HexToAddress("0x0005"),
	}
	const (
		round           = uint64(0)
		govCommitteeLen = uint64(2)
	)
	proposer := validatorAddrs[0]

	for _, tc := range []struct {
		name              string
		committee         []common.Address
		expectedCommittee uint64
	}{
		{name: "subset committee keeps old effective size", committee: validatorAddrs[:int(govCommitteeLen)], expectedCommittee: govCommitteeLen},
		{name: "full committee uses qualified length", committee: validatorAddrs, expectedCommittee: uint64(len(validatorAddrs))},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const seq = uint64(7)
			mockCtrl := gomock.NewController(t)
			t.Cleanup(mockCtrl.Finish)

			mockValsetModule := valset_mock.NewMockValsetModule(mockCtrl)
			mockGovModule := mock_gov.NewMockGovModule(mockCtrl)
			mockValsetModule.EXPECT().GetCouncil(seq).Return(validatorAddrs, nil)
			mockValsetModule.EXPECT().GetDemotedValidators(seq).Return([]common.Address{}, nil)
			mockValsetModule.EXPECT().GetCommittee(seq, round).Return(tc.committee, nil)
			mockValsetModule.EXPECT().GetProposer(seq, round).Return(proposer, nil)

			c := &core{valsetModule: mockValsetModule, govModule: mockGovModule}
			qualified, committee, gotProposer, committeeSize, requiredMsgCnt, fNum, err := getRoundCommitteeState(c, seq, round)
			require.NoError(t, err)

			assert.Equal(t, validatorAddrs, qualified.List())
			assert.Equal(t, tc.committee, committee.List())
			assert.Equal(t, proposer, gotProposer)
			assert.Equal(t, tc.expectedCommittee, committeeSize)
			assert.Equal(t, calcQuorumSize(len(validatorAddrs), tc.expectedCommittee), requiredMsgCnt)
			assert.Equal(t, calcFaultTolerance(len(validatorAddrs), tc.expectedCommittee), fNum)
		})
	}
}

func TestStartNewRoundUsesViewRound(t *testing.T) {
	validators := []common.Address{
		common.HexToAddress("0x1"),
		common.HexToAddress("0x2"),
		common.HexToAddress("0x3"),
		common.HexToAddress("0x4"),
	}
	ctrl := gomock.NewController(t)
	backend := mock_istanbul.NewMockBackend(ctrl)
	backend.EXPECT().Address().Return(validators[0]).AnyTimes()
	backend.EXPECT().LastProposal().Return(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(10)}), common.Address{})
	backend.EXPECT().SetCurrentView(gomock.Any())
	backend.EXPECT().EventMux().Return(new(event.TypeMux))

	mValset := valset_mock.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCouncil(uint64(11)).Return(validators, nil)
	mValset.EXPECT().GetDemotedValidators(uint64(11)).Return(nil, nil)
	mValset.EXPECT().GetCommittee(uint64(11), uint64(0)).Return(validators, nil)
	mValset.EXPECT().GetProposer(uint64(11), uint64(0)).Return(validators[1], nil)

	c := New(backend, istanbul.DefaultConfig).(*core)
	c.RegisterKaiaxModules(mValset, mock_gov.NewMockGovModule(ctrl))
	c.current = newRoundState(
		&bft.View{Sequence: big.NewInt(10), Round: big.NewInt(0)},
		valset.NewAddressSet(validators), common.Hash{}, nil, nil, nil, backend.HasBadProposal,
	)
	c.startNewRound(big.NewInt(3))
	t.Cleanup(c.stopTimer)

	assert.Equal(t, uint64(0), c.current.Round().Uint64())
	assert.Equal(t, validators[1], c.current.proposer)
}

// getTestCommitteeState returns qualified, committee, proposer, nonCommittee for tests using the same logic as newMockBackend
func getTestCommitteeState(validatorAddrs []common.Address, committeeSize uint64, seq, round uint64) (qualified, committee *valset.AddressSet, proposer common.Address, nonCommittee *valset.AddressSet) {
	qualified = valset.NewAddressSet(validatorAddrs).Subtract(valset.NewAddressSet([]common.Address{}))
	var committeeAddrs []common.Address
	if round == 2 {
		committeeAddrs = validatorAddrs[2 : committeeSize+2]
	} else {
		committeeAddrs = validatorAddrs[0:committeeSize]
	}
	committee = valset.NewAddressSet(committeeAddrs)
	proposer = validatorAddrs[0]
	nonCommittee = qualified.Subtract(committee)
	return qualified, committee, proposer, nonCommittee
}

func TestGetRoundCommitteeStateMissingModules(t *testing.T) {
	_, _, _, _, _, _, err := getRoundCommitteeState(&core{}, 1, 0)
	assert.ErrorIs(t, err, istanbul.ErrNoEssentialModule)
}

// genValidators returns a set of addresses and corresponding keys used for generating a validator set
func genValidators(n int) ([]common.Address, map[common.Address]*ecdsa.PrivateKey) {
	addrs := make([]common.Address, n)
	keyMap := make(map[common.Address]*ecdsa.PrivateKey, n)

	for i := range n {
		key, _ := crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(key.PublicKey)
		keyMap[addrs[i]] = key
	}
	return addrs, keyMap
}

// signBlock signs the given block with the given private key
func signBlock(block *types.Block, privateKey *ecdsa.PrivateKey) (*types.Block, error) {
	header := block.Header()

	testSealer := istanbul.NewSealerImpl(privateKey)
	authorSeal, err := testSealer.MakeAuthorSeal(header)
	if err != nil {
		return nil, err
	}
	if err := testSealer.WriteAuthorSeal(header, authorSeal); err != nil {
		return nil, err
	}
	return block.WithSeal(header), nil
}

// genBlock generates a signed block indicating prevBlock with ParentHash
func genBlock(prevBlock *types.Block, signerKey *ecdsa.PrivateKey) (*types.Block, error) {
	block := types.NewBlockWithHeader(&types.Header{
		ParentHash: prevBlock.Hash(),
		Number:     new(big.Int).Add(prevBlock.Number(), common.Big1),
		GasUsed:    0,
		Extra:      prevBlock.Extra(),
		Time:       new(big.Int).Add(prevBlock.Time(), common.Big1),
		BlockScore: new(big.Int).Add(prevBlock.BlockScore(), common.Big1),
	})
	return signBlock(block, signerKey)
}

// genIstanbulMsg generates an istanbul message with given values, in the wire
// format selected by permissionless.
func genIstanbulMsg(msgType uint64, prevHash common.Hash, proposal *types.Block, signerAddr common.Address, signerKey *ecdsa.PrivateKey, permissionless bool) (istanbul.MessageEvent, error) {
	return genIstanbulMsgWithSealKey(msgType, prevHash, proposal, signerAddr, signerKey, signerKey, permissionless)
}

// genIstanbulMsgWithSealKey builds an istanbul message signed (outer signature) by
// signerKey, but whose COMMIT CommittedSeal is signed by sealKey over the legacy
// (non-round-bound) preimage. A sealKey different from signerKey forges the
// committed seal.
func genIstanbulMsgWithSealKey(msgType uint64, prevHash common.Hash, proposal *types.Block, signerAddr common.Address, signerKey, sealKey *ecdsa.PrivateKey, permissionless bool) (istanbul.MessageEvent, error) {
	view := &bft.View{Round: big.NewInt(0), Sequence: proposal.Number()}
	var payload any
	switch msgType {
	case bft.MsgPreprepare:
		payload = &bft.Preprepare{View: view, Proposal: proposal}
	case bft.MsgPrepare:
		payload = &bft.Prepare{View: view, Digest: proposal.Hash()}
	case bft.MsgCommit:
		// A COMMIT carries a CommittedSeal: the sender's signature over the proposal's
		// committed-seal preimage, which handleCommit verifies.
		seal, err := crypto.Sign(crypto.Keccak256(istanbul.PrepareCommittedSeal(proposal.Hash())), sealKey)
		if err != nil {
			return istanbul.MessageEvent{}, err
		}
		payload = &bft.Commit{View: view, Digest: proposal.Hash(), CommittedSeal: seal}
	case bft.MsgRoundChange:
		payload = &bft.RoundChange{View: view}
	default:
		return istanbul.MessageEvent{}, bft.ErrInvalidMessage
	}
	return signIstanbulMsg(prevHash, msgType, payload, signerAddr, signerKey, permissionless)
}

// genIstanbulCommitWithRoundBoundSeal builds a COMMIT with a round-bound committed seal (post-permissionless).
func genIstanbulCommitWithRoundBoundSeal(prevHash common.Hash, proposal *types.Block, signerAddr common.Address, signerKey *ecdsa.PrivateKey, round byte) (istanbul.MessageEvent, error) {
	seal, err := crypto.Sign(crypto.Keccak256(istanbul.PrepareCommittedSealWithRound(proposal.Hash(), round)), signerKey)
	if err != nil {
		return istanbul.MessageEvent{}, err
	}
	commit := &bft.Commit{
		View:          &bft.View{Round: big.NewInt(int64(round)), Sequence: proposal.Number()},
		Digest:        proposal.Hash(),
		CommittedSeal: seal,
	}
	return signIstanbulMsg(prevHash, bft.MsgCommit, commit, signerAddr, signerKey, true)
}

// signIstanbulMsg assembles a bft.Message carrying payload and signs it with
// signerKey in the wire format selected by permissionless.
func signIstanbulMsg(prevHash common.Hash, code uint64, payload any, signerAddr common.Address, signerKey *ecdsa.PrivateKey, permissionless bool) (istanbul.MessageEvent, error) {
	encodedPayload, err := bft.Encode(payload)
	if err != nil {
		return istanbul.MessageEvent{}, err
	}
	msg := &bft.Message{
		PrevHash: prevHash,
		Code:     code,
		Msg:      encodedPayload,
		Address:  signerAddr,
	}
	isPermissionlessAt := func(uint64) bool { return permissionless }
	data, err := msg.PayloadNoSigForFork(isPermissionlessAt)
	if err != nil {
		return istanbul.MessageEvent{}, err
	}
	if msg.Signature, err = crypto.Sign(crypto.Keccak256(data), signerKey); err != nil {
		return istanbul.MessageEvent{}, err
	}
	wire, err := msg.PayloadForFork(isPermissionlessAt)
	if err != nil {
		return istanbul.MessageEvent{}, err
	}
	return istanbul.MessageEvent{Hash: msg.PrevHash, Payload: wire}, nil
}

func neverPermissionless(uint64) bool  { return false }
func alwaysPermissionless(uint64) bool { return true }

// startCoreAtPreprepare starts a core, drives a preprepare, and returns the core plus a committee
// sender (≠ proposer) with its key, ready to send COMMITs. permissionless sets the backend fork status.
func startCoreAtPreprepare(t *testing.T, permissionless bool) (c *core, prevHash common.Hash, proposal *types.Block, sender common.Address, senderKey *ecdsa.PrivateKey) {
	t.Helper()
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	t.Cleanup(fork.ClearHardForkBlockNumberConfig)

	addrs, keys := genValidators(30)
	mockBackend, mockCtrl, mockValset, mockGov := newMockBackend(t, addrs, permissionless)
	t.Cleanup(mockCtrl.Finish)

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom
	c = New(mockBackend, istConfig).(*core)
	c.RegisterKaiaxModules(mockValset, mockGov)
	require.NoError(t, c.Start())
	t.Cleanup(func() { c.Stop() })

	lastProposal, _ := mockBackend.LastProposal()
	lastBlock := lastProposal.(*types.Block)
	_, committee, proposer, _ := getTestCommitteeState(addrs, uint64(len(addrs)/3), c.currentView().Sequence.Uint64(), c.currentView().Round.Uint64())
	cand, err := genBlock(lastBlock, keys[proposer])
	require.NoError(t, err)
	pp, err := genIstanbulMsg(bft.MsgPreprepare, lastBlock.Hash(), cand, proposer, keys[proposer], permissionless)
	require.NoError(t, err)
	require.NoError(t, mockBackend.EventMux().Post(pp))
	time.Sleep(time.Second)
	require.NotNil(t, c.current.Preprepare)

	for i := 0; i < committee.Len(); i++ {
		if committee.At(i) != proposer {
			sender = committee.At(i)
			break
		}
	}
	return c, lastBlock.Hash(), c.current.Preprepare.Proposal.(*types.Block), sender, keys[sender]
}

// TestCore_handleCommit_RejectsForgedCommittedSeal checks that a COMMIT from a committee member
// whose CommittedSeal is signed by a different key (a forged seal) is rejected, while a
// correctly-signed COMMIT is accepted.
func TestCore_handleCommit_RejectsForgedCommittedSeal(t *testing.T) {
	c, prevHash, proposal, sender, senderKey := startCoreAtPreprepare(t, false)

	// Seal signed by a key other than the sender's: recovered signer != sender, so it is rejected.
	forgeKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	forged, err := genIstanbulMsgWithSealKey(bft.MsgCommit, prevHash, proposal, sender, senderKey, forgeKey, false)
	require.NoError(t, err)
	require.NoError(t, c.backend.EventMux().Post(forged))
	time.Sleep(time.Second)
	assert.Equal(t, 0, len(c.current.Commits.messages), "forged committed seal must be rejected")

	// A correctly-signed seal from the same sender is accepted (so it was the seal, not the sender).
	valid, err := genIstanbulMsg(bft.MsgCommit, prevHash, proposal, sender, senderKey, false)
	require.NoError(t, err)
	require.NoError(t, c.backend.EventMux().Post(valid))
	time.Sleep(time.Second)
	assert.Equal(t, 1, len(c.current.Commits.messages), "valid committed seal must be accepted")
}

// TestCore_handleCommit_PermissionlessRoundBoundSeal checks that post-permissionless handleCommit
// rejects a legacy committed seal and accepts a round-bound one from the same committee member.
func TestCore_handleCommit_PermissionlessRoundBoundSeal(t *testing.T) {
	c, prevHash, proposal, sender, senderKey := startCoreAtPreprepare(t, true)

	// A legacy (non-round-bound) committed seal must be rejected post-permissionless.
	legacy, err := genIstanbulMsg(bft.MsgCommit, prevHash, proposal, sender, senderKey, true)
	require.NoError(t, err)
	require.NoError(t, c.backend.EventMux().Post(legacy))
	time.Sleep(time.Second)
	assert.Equal(t, 0, len(c.current.Commits.messages), "legacy committed seal must be rejected post-permissionless")

	// A round-bound committed seal from the same sender is accepted.
	roundBound, err := genIstanbulCommitWithRoundBoundSeal(prevHash, proposal, sender, senderKey, 0)
	require.NoError(t, err)
	require.NoError(t, c.backend.EventMux().Post(roundBound))
	time.Sleep(time.Second)
	assert.Equal(t, 1, len(c.current.Commits.messages), "round-bound committed seal must be accepted post-permissionless")
}

func TestCore_handlerMsg(t *testing.T) {
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	validatorAddrs, validatorKeyMap := genValidators(10)
	mockBackend, mockCtrl, mockValset, mockGov := newMockBackend(t, validatorAddrs, false)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig.Copy()
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	istCore.RegisterKaiaxModules(mockValset, mockGov)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}
	defer istCore.Stop()

	lastProposal, _ := mockBackend.LastProposal()
	lastBlock := lastProposal.(*types.Block)
	committeeSize := uint64(len(validatorAddrs) / 3)
	_, _, proposer, _ := getTestCommitteeState(validatorAddrs, committeeSize, lastBlock.NumberU64()+1, 0)

	// invalid format
	{
		invalidMsg := []byte{0x1, 0x2, 0x3, 0x4}
		err := istCore.handleMsg(invalidMsg)
		assert.NotNil(t, err)
	}

	// invali sender (non-validator)
	{
		newAddr, keyMap := genValidators(1)
		nonValidatorAddr := newAddr[0]
		nonValidatorKey := keyMap[nonValidatorAddr]

		newProposal, err := genBlock(lastBlock, nonValidatorKey)
		if err != nil {
			t.Fatal(err)
		}

		istanbulMsg, err := genIstanbulMsg(bft.MsgPreprepare, lastBlock.Hash(), newProposal, nonValidatorAddr, nonValidatorKey, false)
		if err != nil {
			t.Fatal(err)
		}

		err = istCore.handleMsg(istanbulMsg.Payload)
		assert.NotNil(t, err)
	}

	// valid message
	{
		msgSender := proposer
		msgSenderKey := validatorKeyMap[msgSender]

		newProposal, err := genBlock(lastBlock, msgSenderKey)
		if err != nil {
			t.Fatal(err)
		}

		istanbulMsg, err := genIstanbulMsg(bft.MsgPreprepare, lastBlock.Hash(), newProposal, msgSender, msgSenderKey, false)
		if err != nil {
			t.Fatal(err)
		}

		err = istCore.handleMsg(istanbulMsg.Payload)
		assert.Nil(t, err)
	}
}

// TestCore_postPrepreparedEvent checks that the event loop keeps processing after
// an invalid sender and emits PrepreparedEvent for the subsequent valid proposal.
func TestCore_postPrepreparedEvent(t *testing.T) {
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	validatorAddrs, validatorKeyMap := genValidators(10)
	mockBackend, mockCtrl, mockValset, mockGov := newMockBackend(t, validatorAddrs, false)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig.Copy()
	istConfig.ProposerPolicy = istanbul.WeightedRandom
	istCore := New(mockBackend, istConfig).(*core)
	istCore.RegisterKaiaxModules(mockValset, mockGov)
	require.NoError(t, istCore.Start())
	defer istCore.Stop()

	lastProposal, _ := mockBackend.LastProposal()
	lastBlock := lastProposal.(*types.Block)
	wantSeq := lastBlock.NumberU64() + 1
	committeeSize := uint64(len(validatorAddrs) / 3)
	_, _, proposer, nonCommittee := getTestCommitteeState(validatorAddrs, committeeSize, wantSeq, 0)

	invalidSender := nonCommittee.At(0)
	invalidProposal, err := genBlock(lastBlock, validatorKeyMap[invalidSender])
	require.NoError(t, err)
	invalidMsg, err := genIstanbulMsg(bft.MsgPreprepare, lastBlock.Hash(), invalidProposal,
		invalidSender, validatorKeyMap[invalidSender], false)
	require.NoError(t, err)

	proposal, err := genBlock(lastBlock, validatorKeyMap[proposer])
	require.NoError(t, err)
	validMsg, err := genIstanbulMsg(bft.MsgPreprepare, lastBlock.Hash(), proposal,
		proposer, validatorKeyMap[proposer], false)
	require.NoError(t, err)
	require.NotEqual(t, invalidProposal.Hash(), proposal.Hash())

	mux := mockBackend.EventMux()
	sub := mux.Subscribe(istanbul.PrepreparedEvent{})
	defer sub.Unsubscribe()

	// Post and receive concurrently: TypeMux delivery is unbuffered.
	postErrors := make(chan error, 1)
	go func() {
		for _, msg := range []istanbul.MessageEvent{invalidMsg, validMsg} {
			if err := mux.Post(msg); err != nil {
				postErrors <- err
				return
			}
		}
	}()

	select {
	case ev := <-sub.Chan():
		require.NotNil(t, ev)
		pe, ok := ev.Data.(istanbul.PrepreparedEvent)
		require.True(t, ok)
		require.NotNil(t, pe.Block)
		require.Equal(t, proposal.Hash(), pe.Block.Hash(), "only the valid proposal may produce the event")
		require.Equal(t, wantSeq, pe.View.Sequence.Uint64())
		require.Equal(t, uint64(0), pe.View.Round.Uint64())
	case err := <-postErrors:
		t.Fatalf("failed to post consensus message: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("expected PrepreparedEvent after rejecting the invalid sender")
	}
}

// TestCore_handleTimeoutMsg_race tests a race condition between round change triggers.
// There should be no race condition when round change message and timeout event are handled simultaneously.
func TestCore_handleTimeoutMsg_race(t *testing.T) {
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	// important variables to construct test cases
	const sleepTime = 200 * time.Millisecond
	const processingTime = 400 * time.Millisecond

	type testCase struct {
		name          string
		timeoutTime   time.Duration
		messageRound  int64
		expectedRound int64
	}
	testCases := []testCase{
		{
			// if timeoutTime < sleepTime,
			// timeout event will be posted and then round change message will be processed
			name:          "timeout before processing the (2f+1)th round change message",
			timeoutTime:   50 * time.Millisecond,
			messageRound:  10,
			expectedRound: 10,
		},
		{
			// if timeoutTime > sleepTime && timeoutTime < (processingTime + sleepTime),
			// timeout event will be posted during the processing of (2f+1)th round change message
			name:          "timeout during processing the (2f+1)th round change message",
			timeoutTime:   300 * time.Millisecond,
			messageRound:  20,
			expectedRound: 20,
		},
	}

	validatorAddrs, validatorKeys := genValidators(10)
	mockBackend, mockCtrl, mockValset, mockGov := newMockBackend(t, validatorAddrs, false)
	defer mockCtrl.Finish()

	istConfig := istanbul.DefaultConfig
	istConfig.ProposerPolicy = istanbul.WeightedRandom

	istCore := New(mockBackend, istConfig).(*core)
	istCore.RegisterKaiaxModules(mockValset, mockGov)
	if err := istCore.Start(); err != nil {
		t.Fatal(err)
	}
	defer istCore.Stop()

	eventMux := mockBackend.EventMux()
	lastProposal, _ := mockBackend.LastProposal()
	sequence := istCore.current.sequence.Int64()

	for _, tc := range testCases {
		handler := func(t *testing.T) {
			roundChangeTimer := istCore.roundChangeTimer.Load().(*time.Timer)

			// reset timeout timer of this round and wait some time
			roundChangeTimer.Reset(tc.timeoutTime)
			time.Sleep(sleepTime)

			// `istCore.validateFn` will be executed on processing a istanbul message
			istCore.validateFn = func(arg1 []byte, arg2 []byte) (common.Address, error) {
				// delays the processing of a istanbul message
				time.Sleep(processingTime)
				return istCore.checkValidatorSignature(arg1, arg2)
			}

			// prepare a round change message payload
			payload := makeRCMsgPayload(t, tc.messageRound, sequence, lastProposal.Hash(), validatorAddrs[0], validatorKeys[validatorAddrs[0]])
			if payload == nil {
				t.Fatal("failed to make a round change message payload")
			}

			// one round change message changes the round because the committee size of mockBackend is 3
			err := eventMux.Post(istanbul.MessageEvent{
				Hash:    lastProposal.Hash(),
				Payload: payload,
			})
			if err != nil {
				t.Fatal(err)
			}

			// wait until the istanbul message have processed
			time.Sleep(processingTime + sleepTime)
			roundChangeTimer.Stop()

			// check the result
			assert.Equal(t, tc.expectedRound, istCore.current.round.Int64())
		}
		t.Run(tc.name, handler)
	}
}

// makeRCMsgPayload makes a payload of round change message.
func makeRCMsgPayload(t *testing.T, round int64, sequence int64, prevHash common.Hash, senderAddr common.Address, signerKey *ecdsa.PrivateKey) []byte {
	encoded, err := bft.Encode(&bft.RoundChange{
		View: &bft.View{
			Round:    big.NewInt(round),
			Sequence: big.NewInt(sequence),
		},
	})
	require.Nil(t, err)

	msg := &bft.Message{
		PrevHash: prevHash,
		Code:     bft.MsgRoundChange,
		Msg:      encoded,
		Address:  senderAddr,
	}

	// The mock backend in this test is pre-Permissionless.
	data, err := msg.PayloadNoSigForFork(neverPermissionless)
	require.Nil(t, err)

	msg.Signature, err = crypto.Sign(crypto.Keccak256([]byte(data)), signerKey)
	require.Nil(t, err)

	payload, err := msg.PayloadForFork(neverPermissionless)
	require.Nil(t, err)

	return payload
}

// An oversized PREPARE, COMMIT or ROUND CHANGE must be rejected before any of
// the retention paths (backlog, roundChangeSet, messageSet) can keep it. The
// signed payload can be inflated through a view field, since rlp does not cap
// big.Int on decode, or through any variable-length payload field.
func TestHandleCheckedMsgRejectsOversizedSubjectMessage(t *testing.T) {
	src := common.HexToAddress("0x1")
	oversizedView := &bft.View{
		Sequence: big.NewInt(1),
		Round:    new(big.Int).Lsh(big.NewInt(1), 8*maxSubjectMessageBytes),
	}
	encode := func(val any) []byte {
		encoded, err := bft.Encode(val)
		require.NoError(t, err)
		return encoded
	}
	view := &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)}
	payloads := map[uint64]map[string][]byte{
		bft.MsgPrepare: {"view field": encode(&bft.Prepare{View: oversizedView})},
		bft.MsgCommit: {
			"view field":     encode(&bft.Commit{View: oversizedView}),
			"committed seal": encode(&bft.Commit{View: view, CommittedSeal: make([]byte, maxSubjectMessageBytes+1)}),
		},
		bft.MsgRoundChange: {"view field": encode(&bft.RoundChange{View: oversizedView})},
	}

	for code, byName := range payloads {
		for name, payload := range byName {
			msg := &bft.Message{Code: code, Msg: payload}
			t.Run(fmt.Sprintf("code %d/%s", code, name), func(t *testing.T) {
				c := newTestBacklogCore()
				c.roundChangeSet = newRoundChangeSet(valset.NewAddressSet([]common.Address{src}), 1)

				err := c.handleCheckedMsg(msg, src)

				require.ErrorIs(t, err, errMessageTooLarge)
				assert.Empty(t, c.backlogs)
				assert.Zero(t, c.current.Prepares.Size())
				assert.Zero(t, c.current.Commits.Size())
				assert.Empty(t, c.roundChangeSet.roundChanges)
			})
		}
	}
}

// A well-formed message of every code must stay inside the size limit, even at
// the largest view. PREPREPARE carries a block and is bounded by the block size
// rather than by this check.
func TestCheckMessageSizeFitsWellFormedMessages(t *testing.T) {
	view := &bft.View{
		Sequence: new(big.Int).SetUint64(^uint64(0)),
		Round:    new(big.Int).SetUint64(^uint64(0)),
	}
	maxRound := new(big.Int).SetUint64(^uint64(0))
	for code, payload := range map[uint64]any{
		bft.MsgPrepare: &bft.Prepare{View: view, Digest: common.HexToHash("0x01")},
		bft.MsgCommit:  &bft.Commit{View: view, Digest: common.HexToHash("0x01"), CommittedSeal: make([]byte, crypto.SignatureLength)},
		bft.MsgRoundChange: &bft.RoundChange{View: view, Prepared: &bft.PreparedClaim{
			Round: maxRound, Digest: common.HexToHash("0x01"),
		}},
	} {
		encoded, err := bft.Encode(payload)
		require.NoError(t, err)
		require.NoError(t, checkMessageSize(&bft.Message{
			PrevHash:  common.HexToHash("0x02"),
			Code:      code,
			Msg:       encoded,
			Signature: make([]byte, crypto.SignatureLength),
		}, alwaysPermissionless))
	}
	require.NoError(t, checkMessageSize(&bft.Message{
		Code: bft.MsgPreprepare,
		Msg:  make([]byte, maxSubjectMessageBytes+1),
	}, alwaysPermissionless))
}

func TestCheckMessageSizeIncludesConsensusP2PWrapper(t *testing.T) {
	encoded, err := bft.Encode(&bft.RoundChange{View: &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)}})
	require.NoError(t, err)
	msg := &bft.Message{
		Code:      bft.MsgRoundChange,
		Msg:       encoded,
		Signature: make([]byte, crypto.SignatureLength),
		Evidence:  make([]byte, maxConsensusP2PMessageBytes-64),
	}
	// The attachment alone is below 12 MiB, but the signed message and its
	// ConsensusMsg wrapper push the wire message over the protocol cap.
	require.LessOrEqual(t, len(msg.Evidence), maxConsensusP2PMessageBytes)
	payload, err := msg.Payload()
	require.NoError(t, err)
	wire, err := bft.Encode(&bft.ConsensusMsg{Payload: payload})
	require.NoError(t, err)
	require.Equal(t, uint64(len(wire)), consensusP2PMessageSize(payload))
	require.Greater(t, consensusP2PMessageSize(payload), uint64(maxConsensusP2PMessageBytes))
	require.ErrorIs(t, checkMessageSize(msg, alwaysPermissionless), errMessageTooLarge)

	msg.Evidence = msg.Evidence[:len(msg.Evidence)-2048]
	require.NoError(t, checkMessageSize(msg, alwaysPermissionless))
}
