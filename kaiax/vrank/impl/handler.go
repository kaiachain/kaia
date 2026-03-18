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
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// HandleIstanbulPreprepare records the view when this node is a validator for the next block.
// When proposer, it broadcasts VRankPreprepare to candidates.
func (v *VRankModule) HandleIstanbulPreprepare(block *types.Block, view *istanbul.View) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return
	}

	prepreparedAt := time.Now()
	blockNum := block.NumberU64()
	// if I'm a committee member (ValActive), then I need to collect VRankCandidate
	// should be isCommitteeMember(blockNum + 1), but committee is not finalized during `blockNum` consensus.
	if v.isCommitteeMember(blockNum, view.Round.Uint64()) {
		copiedView := istanbul.View{
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
		v.BroadcastVRankPreprepare(&vrank.VRankPreprepare{Block: block, View: view})
	}

	if blockNum > maxCollectorWindow {
		v.collector.RemoveOldViews(vrank.ViewKey{N: blockNum - maxCollectorWindow, R: maxRound})
	}
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

		sigHash := v.vrankCandidateSigHash(block.NumberU64(), uint8(view.Round.Uint64()), block.Hash())
		sig, err := crypto.Sign(sigHash.Bytes(), v.NodeKey)
		if err != nil {
			logger.Error("Sign failed", "blockNum", block.NumberU64(), "blockHash", block.Hash().Hex())
			return err
		}
		v.BroadcastVRankCandidate(&vrank.VRankCandidate{
			BlockNumber: block.NumberU64(),
			Round:       uint8(view.Round.Uint64()),
			BlockHash:   block.Hash(),
			Sig:         sig,
		})
	}
	return nil
}

// HandleVRankCandidate stores VRankCandidate from candidates. Verification is performed at TallyCfReport.
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
	if msg.BlockNumber > prepreparedSeqNum+maxCollectorWindow {
		return vrank.ErrTooFar
	}
	if msg.Round > maxRound {
		return vrank.ErrRoundOutOfRange
	}
	if isStaleVRankCandidate(msg, prepreparedSeqNum, prepreparedRound) {
		return nil
	}

	sender, err := v.recoverVRankCandidateSender(msg)
	if err != nil {
		return err
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
	pubkey, err := crypto.SigToPub(sigHash.Bytes(), msg.Sig)
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

func (v *VRankModule) recoverVRankCandidateSender(msg *vrank.VRankCandidate) (common.Address, error) {
	sigHash := v.vrankCandidateSigHash(msg.BlockNumber, msg.Round, msg.BlockHash)
	pubkey, err := crypto.SigToPub(sigHash.Bytes(), msg.Sig)
	if err != nil {
		logger.Debug("SigToPub failed", "err", err, "blockNum", msg.BlockNumber, "blockHash", msg.BlockHash, "sig", msg.Sig)
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
	candidates, err := v.Valset.GetCandidates(block.NumberU64())
	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", block.NumberU64())
		return
	}
	sigHash := v.vrankPreprepareSigHash(block.NumberU64(), uint8(vrankPreprepare.View.Round.Uint64()), block.Hash())
	sig, err := crypto.Sign(sigHash.Bytes(), v.NodeKey)
	if err != nil {
		logger.Error("Sign VRankPreprepare failed", "blockNum", block.NumberU64())
		return
	}
	vrankPreprepare.Sig = sig
	v.broadcast(candidates, vrank.VRankPreprepareMsg, vrankPreprepare)
}

// BroadcastVRankCandidate is called by candidates.
func (v *VRankModule) BroadcastVRankCandidate(vrankCandidate *vrank.VRankCandidate) {
	// should be GetCommittee(blockNum + 1, 0), but committee is not finalized during `blockNum` consensus.
	validators, err := v.Valset.GetCommittee(vrankCandidate.BlockNumber, uint64(vrankCandidate.Round))
	if err != nil || validators == nil {
		logger.Error("GetCommittee failed", "blockNum", vrankCandidate.BlockNumber)
		return
	}

	v.broadcast(validators, vrank.VRankCandidateMsg, vrankCandidate)
}

func (v *VRankModule) broadcast(targets []common.Address, code int, msg any) {
	req := &vrank.VRankBroadcastEvent{
		Targets: targets,
		Code:    code,
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
	candidates, err := v.Valset.GetCandidates(blockNum)
	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", blockNum)
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
