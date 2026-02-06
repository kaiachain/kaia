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

// HandleIstanbulPreprepare is executed by validators for timer start, andthe proposer for broadcast
func (v *VRankModule) HandleIstanbulPreprepare(block *types.Block, view *istanbul.View) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return
	}

	blockNum := block.NumberU64()
	if v.isValidator(blockNum) {
		v.candResponses.Clear()
		v.prepreparedTime = time.Now()
	}
	if v.isProposer(blockNum, view.Round.Uint64()) {
		vrankPreprepare := &vrank.VRankPreprepare{Block: block}
		v.BroadcastVRankPreprepare(vrankPreprepare)
	}
}

// HandleVRankPreprepare is executed by candidates
func (v *VRankModule) HandleVRankPreprepare(preprepare *vrank.VRankPreprepare) error {
	block := preprepare.Block
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return nil
	}
	if preprepare == nil || preprepare.Block == nil {
		logger.Error("VRankPreprepare is nil")
		return vrank.ErrVRankPreprepareNil
	}

	if v.isCandidate(block.NumberU64()) {
		sig, err := crypto.Sign(crypto.Keccak256(block.Hash().Bytes()), v.NodeKey)
		if err != nil {
			logger.Error("Sign failed", "blockNum", block.NumberU64())
			return err
		}
		msg := &vrank.VRankCandidate{
			BlockNumber: block.NumberU64(),
			Round:       0,
			BlockHash:   block.Hash(),
			Sig:         sig,
		}
		v.BroadcastVRankCandidate(msg)
	}
	return nil
}

// HandleVRankCandidate is executed by validators
func (v *VRankModule) HandleVRankCandidate(msg *vrank.VRankCandidate) error {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(msg.BlockNumber)) {
		return nil
	}

	if msg == nil {
		logger.Error("Unexpected nil")
		return vrank.ErrVRankCandidateNil
	}

	elapsed := time.Since(v.prepreparedTime)
	if v.isValidator(msg.BlockNumber) {
		cand, err := istanbul.GetSignatureAddress(msg.BlockHash.Bytes(), msg.Sig)
		if err != nil {
			logger.Error("GetSignatureAddress failed", "blockNum", msg.BlockNumber, "blockHash", msg.BlockHash, "sig", msg.Sig)
			return err
		}
		logger.Trace("HandleVRankCandidate", "cand", cand, "elapsed", elapsed, "blockHash", msg.BlockHash.Hex())
		v.candResponses.Store(cand, elapsed)
	}
	return nil
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
	validators, err := v.Valset.GetCouncil(vrankCandidate.BlockNumber)
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
		logger.Error("GetCandidates failed", "blockNum", blockNum)
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
	for {
		select {
		case req, ok := <-v.broadcastCh:
			if !ok {
				return
			}
			v.broadcastFeed.Send(req)
		}
	}
}
