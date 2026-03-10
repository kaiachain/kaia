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
	// if I'm a validator, then I need to collect VRankCandidate
	// should be isValidator(blockNum + 1), but validators are not finalized during `blockNum` consensus.
	if v.isValidator(blockNum) {
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

// HandleVRankPreprepare processes VRankPreprepare; if this node is a candidate, it broadcasts VRankCandidate.
func (v *VRankModule) HandleVRankPreprepare(msg *vrank.VRankPreprepare) error {
	block := msg.Block
	view := msg.View
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return nil
	}

	if v.isCandidate(block.NumberU64()) {
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
	prepreparedSeqNum := uint64(0)
	hasPrepreparedSeq := v.prepreparedView.Sequence != nil
	if hasPrepreparedSeq {
		prepreparedSeqNum = v.prepreparedView.Sequence.Uint64()
	}
	v.prepreparedViewMu.RUnlock()
	if !hasPrepreparedSeq {
		return vrank.ErrPrepreparedViewNotSet
	}
	// should be isValidator(v.prepreparedView.Sequence.Uint64() + 1), but validators are not finalized during `seq` consensus.
	if v.isValidator(prepreparedSeqNum) {
		if msg.BlockNumber > prepreparedSeqNum+maxCollectorWindow {
			return vrank.ErrTooFar
		}
		if msg.Round > maxRound {
			return vrank.ErrRoundOutOfRange
		}

		sender, err := v.recoverVRankCandidateSender(msg)
		if err != nil {
			return err
		}
		vk := vrank.ViewKey{N: msg.BlockNumber, R: msg.Round}
		if v.collector.HasCandMsg(vk, sender) {
			return nil
		}
		if err := v.verifyVRankCandidateSender(msg, sender, prepreparedSeqNum); err != nil {
			return err
		}
		v.collector.AddCandMsg(vk, sender, receivedAt, msg)
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

func (v *VRankModule) verifyVRankCandidateSender(msg *vrank.VRankCandidate, sender common.Address, prepreparedSeqNum uint64) error {
	// should be GetCandidates(msg.BlockNumber), but candidates are not finalized during `Sequence` consensus.
	candidates, err := v.Valset.GetCandidates(prepreparedSeqNum)
	if err != nil {
		logger.Debug("GetCandidates failed", "err", err, "blockNum", msg.BlockNumber)
		return err
	}
	if !slices.Contains(candidates, sender) {
		logger.Debug("Sender is not a candidate", "sender", sender.Hex(), "blockNum", msg.BlockNumber)
		return vrank.ErrMsgFromNonCandidate
	}
	return nil
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

// BroadcastVRankPreprepare is called by the proposer
func (v *VRankModule) BroadcastVRankPreprepare(vrankPreprepare *vrank.VRankPreprepare) {
	block := vrankPreprepare.Block
	candidates, err := v.Valset.GetCandidates(block.NumberU64())
	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", block.NumberU64())
		return
	}
	v.broadcast(candidates, vrank.VRankPreprepareMsg, vrankPreprepare)
}

// BroadcastVRankCandidate is called by candidates.
func (v *VRankModule) BroadcastVRankCandidate(vrankCandidate *vrank.VRankCandidate) {
	// should be GetCouncil(blockNum + 1), but validators are not finalized during `blockNum` consensus.
	validators, err := v.Valset.GetCouncil(vrankCandidate.BlockNumber)
	if err != nil || validators == nil {
		logger.Error("GetCouncil failed", "blockNum", vrankCandidate.BlockNumber)
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

func (v *VRankModule) isValidator(blockNum uint64) bool {
	validators, err := v.Valset.GetCouncil(blockNum)
	if err != nil || validators == nil {
		logger.Error("GetCouncil failed", "blockNum", blockNum)
		return false
	}

	return slices.Contains(validators, v.nodeID)
}

func (v *VRankModule) handleBroadcastLoop() {
	for {
		select {
		case req := <-v.broadcastCh:
			v.broadcastFeed.Send(req)
		case <-v.stopCh:
			return
		}
	}
}
