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
	"encoding/binary"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	blstypes "github.com/kaiachain/kaia/crypto/bls/types"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// HandleIstanbulPreprepare records the view when this node is a validator for the next block.
// When proposer, it broadcasts VRankPreprepare to candidates.
func (v *VRankModule) HandleIstanbulPreprepare(block *types.Block, view *bft.View) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return
	}

	prepreparedAt := time.Now()
	blockNum := block.NumberU64()
	// if I'm a committee member (ValActive), then I need to collect VRankCandidate
	// ideally isCommitteeMember(blockNum + 1, round), but committee is not finalized during `blockNum` consensus, thus (blockNum, round).
	if v.isCommitteeMember(blockNum, view.Round.Uint64()) {
		copiedView := bft.View{
			Sequence: new(big.Int).Set(view.Sequence),
			Round:    new(big.Int).Set(view.Round),
		}
		v.prepreparedViewMu.Lock()
		v.prepreparedView = copiedView
		v.prepreparedViewMu.Unlock()
		v.collector.AddPrepreparedTime(vrank.ViewKey{N: blockNum, R: uint8(view.Round.Uint64())}, prepreparedAt, block.Hash())
	}
	// if I'm the proposer that broadcast IstanbulPreprepare to other validators,
	// then I need to broadcast VRankPreprepare as well
	if v.isProposer(blockNum, view.Round.Uint64()) {
		v.recordOwnProposal(blockNum)
		v.BroadcastVRankPreprepare(&vrank.VRankPreprepare{Block: block, View: view})
	}

	if blockNum > maxWindow {
		v.collector.RemoveOldViews(vrank.ViewKey{N: blockNum - maxWindow, R: maxRound}, v.ownProposalSnapshot())
	}
}

// recordOwnProposal marks blockNum as proposed by this node so its collector view survives until
// the next proposal reports it. Prior-epoch entries are dropped (CFS is epoch-local).
func (v *VRankModule) recordOwnProposal(blockNum uint64) {
	epochStart := calcEpochStart(blockNum, v.vrankEpoch())
	v.ownProposalsMu.Lock()
	defer v.ownProposalsMu.Unlock()
	v.ownProposals[blockNum] = struct{}{}
	for n := range v.ownProposals {
		if n < epochStart {
			delete(v.ownProposals, n)
		}
	}
}

// ownProposalSnapshot copies the retained own-proposal block numbers (protected from pruning).
func (v *VRankModule) ownProposalSnapshot() map[uint64]struct{} {
	v.ownProposalsMu.Lock()
	defer v.ownProposalsMu.Unlock()
	snap := make(map[uint64]struct{}, len(v.ownProposals))
	for n := range v.ownProposals {
		snap[n] = struct{}{}
	}
	return snap
}

// pruneReportedProposals drops entries strictly below upto. The selected block is kept (strict <)
// so a round change or a failed commit still re-reports it; it goes only once a later block supersedes it.
func (v *VRankModule) pruneReportedProposals(upto uint64) {
	v.ownProposalsMu.Lock()
	defer v.ownProposalsMu.Unlock()
	for n := range v.ownProposals {
		if n < upto {
			delete(v.ownProposals, n)
		}
	}
}

// selectReportTarget returns the most recent block this node produced before number in the same
// epoch. Rounds it proposed but another validator committed are skipped (committed-header proposer
// check). ok=false when none exists (first proposal, or a restart cleared the set).
func (v *VRankModule) selectReportTarget(number uint64) (targetNum, round uint64, ok bool) {
	epochStart := calcEpochStart(number, v.vrankEpoch())
	v.ownProposalsMu.Lock()
	cands := make([]uint64, 0, len(v.ownProposals))
	for n := range v.ownProposals {
		if n < number && n >= epochStart {
			cands = append(cands, n)
		}
	}
	v.ownProposalsMu.Unlock()

	slices.Sort(cands)
	for i := len(cands) - 1; i >= 0; i-- {
		n := cands[i]
		proposer, r, err := v.proposerOf(n)
		if err != nil {
			continue
		}
		if proposer == v.nodeID {
			return n, r, true
		}
	}
	return 0, 0, false
}

// proposerOf returns the proposer and final round of a committed block.
func (v *VRankModule) proposerOf(number uint64) (common.Address, uint64, error) {
	header := v.Chain.GetHeaderByNumber(number)
	if header == nil {
		return common.Address{}, 0, vrank.ErrHeaderNotFound
	}
	roundByte, err := v.RoundReader.Round(header)
	if err != nil {
		return common.Address{}, 0, err
	}
	round := uint64(roundByte)
	proposer, err := v.Valset.GetProposer(number, round)
	if err != nil {
		return common.Address{}, 0, err
	}
	return proposer, round, nil
}

// HandleVRankPreprepare processes VRankPreprepare; if this node is a candidate, it verifies the
// proposer's signature and broadcasts VRankCandidate.
func (v *VRankModule) HandleVRankPreprepare(msg *vrank.VRankPreprepare) error {
	block := msg.Block
	view := msg.View
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return nil
	}

	if v.isCandidate(block.NumberU64()) {
		sender, err := v.recoverVRankPreprepareSender(msg)
		if err != nil {
			return err
		}
		if err := v.verifyVRankPreprepareSender(msg, sender); err != nil {
			return err
		}
		v.pruneSeenPreprepare(block.NumberU64())
		if exactReplay, conflictingView := v.markSeenPreprepare(vrank.ViewKey{N: block.NumberU64(), R: uint8(view.Round.Uint64())}, block.Hash()); exactReplay {
			// ignore seen preprepare
			return nil
		} else if conflictingView {
			logger.Warn("Conflicting VRankPreprepare ignored", "blockNum", block.NumberU64(), "round", view.Round.Uint64(), "blockHash", block.Hash().Hex())
			return nil
		}

		sigHash := v.vrankCandidateSigHash(block.NumberU64(), uint8(view.Round.Uint64()), block.Hash())
		sig, err := crypto.Sign(sigHash.Bytes(), v.NodeKey)
		if err != nil {
			logger.Error("Sign failed", "blockNum", block.NumberU64(), "blockHash", block.Hash().Hex())
			return err
		}
		blsSig := bls.Sign(v.BlsKey, sigHash.Bytes()).Marshal()
		// TODO-Permissionless: Testing only. Remove before production release.
		if v.skipCandidate.Load() {
			logger.Warn("SkipCandidate is enabled, skipping VRankCandidate broadcast")
			return nil
		}
		v.BroadcastVRankCandidate(&vrank.VRankCandidate{
			BlockNumber: block.NumberU64(),
			Round:       uint8(view.Round.Uint64()),
			BlockHash:   block.Hash(),
			Sig:         [crypto.SignatureLength]byte(sig),
			BlsSig:      [blstypes.SignatureLength]byte(blsSig),
		})
	}
	return nil
}

// HandleVRankCandidate stores VRankCandidate from candidates. Verification is performed at EvaluateCandidates.
func (v *VRankModule) HandleVRankCandidate(msg *vrank.VRankCandidate) error {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(msg.BlockNumber)) {
		return nil
	}

	receivedAt := time.Now()
	v.prepreparedViewMu.RLock()
	prepreparedSeqNum, prepreparedRound := uint64(0), uint64(0)
	hasPrepreparedView := v.prepreparedView.Sequence != nil && v.prepreparedView.Round != nil
	if hasPrepreparedView {
		prepreparedSeqNum = v.prepreparedView.Sequence.Uint64()
		prepreparedRound = v.prepreparedView.Round.Uint64()
	}
	v.prepreparedViewMu.RUnlock()
	if !hasPrepreparedView {
		return vrank.ErrPrepreparedViewNotSet
	}
	if msg.BlockNumber > prepreparedSeqNum+maxWindow {
		return vrank.ErrTooFar
	}
	if msg.Round > maxRound {
		return vrank.ErrRoundOutOfRange
	}
	if isStaleVRankCandidate(msg, prepreparedSeqNum, prepreparedRound) {
		return nil
	}

	sigHash := v.vrankCandidateSigHash(msg.BlockNumber, msg.Round, msg.BlockHash)
	sender, err := v.recoverVRankCandidateSender(sigHash, msg.Sig[:])
	if err != nil {
		return err
	}
	blsNum := big.NewInt(0).Add(v.Chain.CurrentHeader().Number, big.NewInt(1)) // head + 1
	blsPub, err := v.Randao.GetBlsPubkey(sender, blsNum)
	if err != nil {
		return fmt.Errorf("%w: %v", vrank.ErrInvalidCandidateBlsSig, err)
	}
	ok, err := bls.VerifySignature(msg.BlsSig[:], sigHash, blsPub)
	if err != nil || !ok {
		return vrank.ErrInvalidCandidateBlsSig
	}
	vk := vrank.ViewKey{N: msg.BlockNumber, R: msg.Round}
	if v.collector.HasCandMsg(vk, sender) {
		return nil
	}
	v.collector.AddCandMsg(vk, sender, receivedAt, msg)
	return nil
}

func isStaleVRankCandidate(msg *vrank.VRankCandidate, prepreparedSeqNum, prepreparedRound uint64) bool {
	if msg.BlockNumber < prepreparedSeqNum {
		return true
	}
	return msg.BlockNumber == prepreparedSeqNum && uint64(msg.Round) < prepreparedRound
}

func (v *VRankModule) pruneSeenPreprepare(currentBlockNum uint64) {
	if currentBlockNum <= maxWindow {
		return
	}
	threshold := currentBlockNum - maxWindow

	v.seenPreprepareMu.Lock()
	defer v.seenPreprepareMu.Unlock()
	for vk := range v.seenPreprepare {
		if vk.N < threshold {
			delete(v.seenPreprepare, vk)
		}
	}
}

// markSeenPreprepare records a candidate response for the given view.
// It returns (true, false) for an exact replay, (false, true) for a conflicting
// block hash in the same view, and (false, false) for a new view/hash pair.
func (v *VRankModule) markSeenPreprepare(vk vrank.ViewKey, blockHash common.Hash) (bool, bool) {
	v.seenPreprepareMu.Lock()
	defer v.seenPreprepareMu.Unlock()

	if seenHash, ok := v.seenPreprepare[vk]; ok {
		if seenHash == blockHash {
			return true, false
		}
		return false, true
	}
	v.seenPreprepare[vk] = blockHash
	return false, false
}

func (v *VRankModule) vrankPreprepareSigHash(blockNum uint64, round uint8, blockHash common.Hash) common.Hash {
	chainID := v.ChainConfig.ChainID.Uint64()

	// Canonical encoding:
	// domain separator || chain_id(uint64 BE) || block_number(uint64 BE) || round(uint8) || block_hash(32 bytes)
	payload := make([]byte, 0, len(vrankPreprepareSigDomain)+8+8+1+len(blockHash))
	payload = append(payload, []byte(vrankPreprepareSigDomain)...)
	payload = binary.BigEndian.AppendUint64(payload, chainID)
	payload = binary.BigEndian.AppendUint64(payload, blockNum)
	payload = append(payload, round)
	payload = append(payload, blockHash[:]...)
	return crypto.Keccak256Hash(payload)
}

func (v *VRankModule) recoverVRankPreprepareSender(msg *vrank.VRankPreprepare) (common.Address, error) {
	sigHash := v.vrankPreprepareSigHash(msg.Block.NumberU64(), uint8(msg.View.Round.Uint64()), msg.Block.Hash())
	pubkey, err := crypto.SigToPub(sigHash.Bytes(), msg.Sig[:])
	if err != nil {
		logger.Debug("SigToPub failed for VRankPreprepare", "err", err, "blockNum", msg.Block.NumberU64())
		return common.Address{}, fmt.Errorf("%w: %v", vrank.ErrInvalidProposerSig, err)
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

func (v *VRankModule) verifyVRankPreprepareSender(msg *vrank.VRankPreprepare, sender common.Address) error {
	blockNum := msg.Block.NumberU64()
	round := msg.View.Round.Uint64()
	proposer, err := v.Valset.GetProposer(blockNum, round)
	if err != nil {
		logger.Debug("GetProposer failed", "err", err, "blockNum", blockNum)
		return err
	}
	if sender != proposer {
		logger.Debug("VRankPreprepare from non-proposer", "sender", sender.Hex(), "proposer", proposer.Hex(), "blockNum", blockNum)
		return vrank.ErrMsgFromNonProposer
	}
	return nil
}

func (v *VRankModule) recoverVRankCandidateSender(sigHash common.Hash, signature []byte) (common.Address, error) {
	pubkey, err := crypto.SigToPub(sigHash.Bytes(), signature)
	if err != nil {
		logger.Debug("SigToPub failed", "err", err, "sigHash", sigHash, "sig", signature)
		return common.Address{}, fmt.Errorf("%w: %v", vrank.ErrInvalidCandidateSig, err)
	}
	sender := crypto.PubkeyToAddress(*pubkey)
	return sender, nil
}

func (v *VRankModule) vrankCandidateSigHash(blockNum uint64, round uint8, blockHash common.Hash) common.Hash {
	chainID := v.ChainConfig.ChainID.Uint64()

	// Canonical encoding:
	// domain separator || chain_id(uint64 BE) || block_number(uint64 BE) || round(uint8) || block_hash(32 bytes)
	payload := make([]byte, 0, len(vrankCandidateSigDomain)+8+8+1+len(blockHash))
	payload = append(payload, []byte(vrankCandidateSigDomain)...)
	payload = binary.BigEndian.AppendUint64(payload, chainID)
	payload = binary.BigEndian.AppendUint64(payload, blockNum)
	payload = append(payload, round)
	payload = append(payload, blockHash[:]...)
	return crypto.Keccak256Hash(payload)
}

// BroadcastVRankPreprepare is called by the proposer. It signs the message before broadcasting.
func (v *VRankModule) BroadcastVRankPreprepare(vrankPreprepare *vrank.VRankPreprepare) {
	block := vrankPreprepare.Block
	candidates, err := v.Valset.GetCandTesting(block.NumberU64())
	if err != nil {
		logger.Error("GetCandTesting failed", "blockNum", block.NumberU64(), "err", err)
		return
	}
	if len(candidates) == 0 {
		return
	}
	sigHash := v.vrankPreprepareSigHash(block.NumberU64(), uint8(vrankPreprepare.View.Round.Uint64()), block.Hash())
	sig, err := crypto.Sign(sigHash.Bytes(), v.NodeKey)
	if err != nil {
		logger.Error("Sign VRankPreprepare failed", "blockNum", block.NumberU64())
		return
	}
	vrankPreprepare.Sig = [crypto.SignatureLength]byte(sig)
	v.broadcast(candidates, vrankPreprepare)
}

// BroadcastVRankCandidate is called by candidates.
func (v *VRankModule) BroadcastVRankCandidate(vrankCandidate *vrank.VRankCandidate) {
	// ideally GetCommittee(blockNum + 1, round), but committee is not finalized during `blockNum` consensus, thus (blockNum, round).
	validators, err := v.Valset.GetCommittee(vrankCandidate.BlockNumber, uint64(vrankCandidate.Round))
	if err != nil || validators == nil {
		logger.Error("GetCommittee failed", "blockNum", vrankCandidate.BlockNumber)
		return
	}

	v.broadcast(validators, vrankCandidate)
}

func (v *VRankModule) broadcast(targets []common.Address, msg any) {
	req := &vrank.VRankBroadcastEvent{
		Targets: targets,
		Msg:     msg,
	}
	v.broadcastCh <- req
}

func (v *VRankModule) isProposer(blockNum, round uint64) bool {
	proposer, err := v.Valset.GetProposer(blockNum, round)
	if err != nil {
		logger.Error("GetProposer failed", "blockNum", blockNum, "round", round)
		return false
	}

	return proposer == v.nodeID
}

func (v *VRankModule) isCandidate(blockNum uint64) bool {
	candidates, err := v.Valset.GetCandTesting(blockNum)
	if err != nil {
		logger.Error("GetCandTesting failed", "blockNum", blockNum, "err", err)
		return false
	}

	return slices.Contains(candidates, v.nodeID)
}

func (v *VRankModule) isCommitteeMember(blockNum, round uint64) bool {
	committee, err := v.Valset.GetCommittee(blockNum, round)
	if err != nil || committee == nil {
		logger.Error("GetCommittee failed", "blockNum", blockNum)
		return false
	}

	return slices.Contains(committee, v.nodeID)
}

func (v *VRankModule) handleBroadcastLoop(stopCh <-chan struct{}) {
	for {
		select {
		case req := <-v.broadcastCh:
			v.broadcastFeed.Send(req)
		case <-stopCh:
			return
		}
	}
}
