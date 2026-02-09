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
	"math/big"
	"slices"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// HandleIstanbulPreprepare starts timer and broadcasts VRankPreprepare to candidates
func (v *VRankModule) HandleIstanbulPreprepare(block *types.Block, view *istanbul.View) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return
	}

	blockNum := block.NumberU64()
	// if I'm a validator for the next block, then I need to collect VRankCandidate
	if v.isValidator(blockNum + 1) {
		v.startNewCollection(block, view)
	}
	// if I'm the proposer that broadcasted IstsanbulPreprepare to other validators,
	// then I need to broadcast VRankPreprepare as well
	if v.isProposer(blockNum, view.Round.Uint64()) {
		v.BroadcastVRankPreprepare(&vrank.VRankPreprepare{Block: block, View: view})
	}
}

// HandleVRankPreprepare broadcasts VRankCandidate to validators
func (v *VRankModule) HandleVRankPreprepare(msg *vrank.VRankPreprepare) error {
	block := msg.Block
	view := msg.View
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return nil
	}

	if v.isCandidate(block.NumberU64()) {
		sig, err := crypto.Sign(crypto.Keccak256(block.Hash().Bytes()), v.NodeKey)
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

// HandleVRankCandidate collects VRankCandidate from candidates
func (v *VRankModule) HandleVRankCandidate(msg *vrank.VRankCandidate) error {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(msg.BlockNumber)) {
		return nil
	}

	if msg == nil {
		logger.Error("Unexpected nil")
		return vrank.ErrVRankCandidateNil
	}

	elapsed := time.Since(v.prepreparedTime)
	// if I'm a validator for the next block, then I need to collect VRankCandidate
	if v.prepreparedView.Sequence != nil && v.isValidator(v.prepreparedView.Sequence.Uint64()+1) {
		sender, err := v.verifyVRankCandidate(msg)
		if err != nil {
			return err
		}

		if _, ok := v.candResponses.Load(sender); !ok {
			logger.Trace("HandleVRankCandidate", "sender", sender.Hex(), "elapsed", elapsed, "blockHash", msg.BlockHash.Hex())
			v.candResponses.Store(sender, elapsed)
		}
	}
	return nil
}

func (v *VRankModule) verifyVRankCandidate(msg *vrank.VRankCandidate) (common.Address, error) {
	if msg.BlockNumber != v.prepreparedView.Sequence.Uint64() {
		logger.Debug("sequence mismatch", "expected", v.prepreparedView.Sequence.Uint64(), "got", msg.BlockNumber)
		return common.Address{}, vrank.ErrViewMismatch
	}
	if msg.Round != uint8(v.prepreparedView.Round.Uint64()) {
		logger.Debug("round mismatch", "expected", v.prepreparedView.Round.Uint64(), "got", msg.Round)
		return common.Address{}, vrank.ErrViewMismatch
	}
	if msg.BlockHash != v.prepreparedBlockHash {
		logger.Debug("BlockHash mismatch", "expected", v.prepreparedBlockHash.Hex(), "got", msg.BlockHash.Hex())
		return common.Address{}, vrank.ErrBlockHashMismatch
	}

	sender, err := istanbul.GetSignatureAddress(msg.BlockHash.Bytes(), msg.Sig)
	if err != nil {
		logger.Debug("GetSignatureAddress failed", "err", err, "blockNum", msg.BlockNumber, "blockHash", msg.BlockHash, "sig", msg.Sig)
		return common.Address{}, err
	}
	candidates, err := v.Valset.GetCandidates(msg.BlockNumber)
	if err != nil || candidates == nil {
		logger.Debug("GetCandidates failed", "err", err, "blockNum", msg.BlockNumber)
		return common.Address{}, err
	}
	if !slices.Contains(candidates, sender) {
		logger.Debug("Sender is not a candidate", "sender", sender.Hex(), "blockNum", msg.BlockNumber)
		return common.Address{}, vrank.ErrMsgFromNonCandidate
	}
	return sender, nil
}

func (v *VRankModule) startNewCollection(block *types.Block, view *istanbul.View) {
	v.prepreparedTime = time.Now()
	v.prepreparedView = *view
	v.prepreparedBlockHash = block.Hash()
	v.candResponses.Clear()
}

// BroadcastVRankPreprepare is called by the proposer
func (v *VRankModule) BroadcastVRankPreprepare(vrankPreprepare *vrank.VRankPreprepare) {
	block := vrankPreprepare.Block
	candidates, err := v.Valset.GetCandidates(block.NumberU64())
	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", block.NumberU64())
		return
	}
	v.broadcast(candidates, VRankPreprepareMsg, vrankPreprepare)
}

// BroadcastVRankPreprepare is called by candidates
func (v *VRankModule) BroadcastVRankCandidate(vrankCandidate *vrank.VRankCandidate) {
	validators, err := v.Valset.GetCouncil(vrankCandidate.BlockNumber + 1)
	if err != nil || validators == nil {
		logger.Error("GetCouncil failed", "blockNum", vrankCandidate.BlockNumber)
		return
	}

	v.broadcast(validators, VRankCandidateMsg, vrankCandidate)
}

func (v *VRankModule) broadcast(targets []common.Address, code int, msg any) {
	req := &vrank.BroadcastRequest{
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

	return proposer == v.nodeId
}

func (v *VRankModule) isCandidate(blockNum uint64) bool {
	candidates, err := v.Valset.GetCandidates(blockNum)
	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", blockNum)
		return false
	}

	return slices.Contains(candidates, v.nodeId)
}

func (v *VRankModule) isValidator(blockNum uint64) bool {
	validators, err := v.Valset.GetCouncil(blockNum)
	if err != nil || validators == nil {
		logger.Error("GetCouncil failed", "blockNum", blockNum)
		return false
	}

	return slices.Contains(validators, v.nodeId)
}

func (v *VRankModule) GetCfReport(blockNum uint64) (vrank.CfReport, error) {
	if blockNum == 0 {
		return vrank.CfReport{}, nil
	}

	candidates, err := v.Valset.GetCandidates(blockNum)

	if err != nil || candidates == nil {
		logger.Error("GetCandidates failed", "blockNum", blockNum)
		return nil, vrank.ErrGetCandidateFailed
	}

	var cfReport vrank.CfReport
	for _, addr := range candidates {
		elapsed, ok := v.candResponses.Load(addr)
		if !ok || elapsed.(time.Duration) > candidatePrepareDeadlineMs {
			cfReport = append(cfReport, addr)
		}
	}

	return cfReport, nil
}

func (v *VRankModule) handleBroadcastLoop() {
	for req := range v.broadcastCh {
		v.broadcastFeed.Send(req)
	}
}
