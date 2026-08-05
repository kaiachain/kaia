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

// HandleIstanbulPreprepare records the view and broadcasts VRankPreprepare to the candidates when
// this node is the proposer of the next block. Other nodes do nothing.
func (v *VRankModule) HandleIstanbulPreprepare(block *types.Block, view *bft.View) {
	if !v.ChainConfig.IsPermissionlessForkEnabled(block.Number()) {
		return
	}

	prepreparedAt := time.Now()
	blockNum := block.NumberU64()
	// Only the proposer keeps VRank state for this view: it sends VRankPreprepare to CandTesting so
	// the candidates reply with VRankCandidate (distinct from the Istanbul Preprepare consensus
	// already sent to validators), and it records the preprepared time to measure those replies.
	// It reports the result from its own next proposal, so no other node needs to collect.
	if v.isProposer(blockNum, view.Round.Uint64()) {
		v.collector.AddPrepreparedTime(vrank.ViewKey{N: blockNum, R: uint8(view.Round.Uint64())}, prepreparedAt, block.Hash())
		// Prior-epoch views can never be reported. The report-time prune would drop them too, but a
		// node that keeps losing its rounds has nothing to report and never gets there.
		v.collector.PruneReported(calcEpochStart(blockNum, v.vrankEpoch()))
		v.BroadcastVRankPreprepare(&vrank.VRankPreprepare{Block: block, View: view})
	}
}

// selectReportTarget returns the most recent block this node produced before number in the same
// epoch. Rounds it proposed but another validator committed are skipped (committed-header proposer
// check). ok=false when none exists (first proposal, or a restart cleared the collector).
func (v *VRankModule) selectReportTarget(number uint64) (targetNum, round uint64, ok bool) {
	cands := v.collector.PendingEvaluations(calcEpochStart(number, v.vrankEpoch()))
	slices.Sort(cands)
	for i := len(cands) - 1; i >= 0; i-- {
		n := cands[i]
		// The collector may already hold the view of the block being built; report only on prior ones.
		if n >= number {
			continue
		}
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
		}, sender)
	}
	return nil
}

// HandleVRankCandidate stores a VRankCandidate reply for a view this node proposed. Candidates
// reply to the proposer directly, and only the proposer reports on the view (from its own next
// proposal), so a reply for any other view is dropped. Full validation is at EvaluateCandidates.
func (v *VRankModule) HandleVRankCandidate(msg *vrank.VRankCandidate) error {
	if !v.ChainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(msg.BlockNumber)) {
		return nil
	}
	receivedAt := time.Now()
	if msg.Round > maxRound {
		return vrank.ErrRoundOutOfRange
	}
	vk := vrank.ViewKey{N: msg.BlockNumber, R: msg.Round}
	// Accept only a view this node proposed (its preprepared time is recorded). This drops forged
	// or misdirected replies cheaply, before signature recovery.
	if !v.collector.HasPreprepared(vk) {
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
	if v.collector.HasCandMsg(vk, sender) {
		return nil
	}
	v.collector.AddCandMsg(vk, sender, receivedAt, msg)
	return nil
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

// BroadcastVRankCandidate is called by a candidate to reply to the block's proposer — the only node
// that reports on this view, since it collects the replies to its own VRankPreprepare.
func (v *VRankModule) BroadcastVRankCandidate(vrankCandidate *vrank.VRankCandidate, proposer common.Address) {
	v.broadcast([]common.Address{proposer}, vrankCandidate)
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
