package impl

import (
	"reflect"
	"slices"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"golang.org/x/exp/maps"
)

func (h *headerGovModule) VerifyHeader(header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil
	}
	err := h.VerifyVote(header)
	if err != nil {
		logger.Error("VerifyVote error", "num", header.Number.Uint64(), "vote", header.Vote, "err", err)
		return err
	}

	return h.VerifyGov(header)
}

func (h *headerGovModule) PrepareHeader(header *types.Header) error {
	// if this node has a vote waiting to be casted, put Vote field.
	if len(h.myVotes) > 0 {
		header.Vote, _ = h.myVotes[0].ToVoteBytes()
		logger.Debug("Prepare header with vote", "num", header.Number.Uint64(), "vote", hexutil.Encode(header.Vote))
	}

	// if epoch block & vote exists in the last epoch, put Governance field.
	if header.Number.Uint64()%h.epoch == 0 {
		gov := h.getExpectedGovernance(header.Number.Uint64())
		if len(gov.Items()) > 0 {
			header.Governance, _ = gov.ToGovBytes()
			logger.Debug("Prepare header with governance", "num", header.Number.Uint64(), "governance", hexutil.Encode(header.Governance))
		}
	}

	return nil
}

func (h *headerGovModule) FinalizeHeader(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	return nil
}

// VerifyVote checks the followings:
// (1) voter must be in valset,
// (2) integrity of the voter (the voter must be the block proposer),
// (3) the vote value must be consistent compared to the latest ParamSet.
func (h *headerGovModule) VerifyVote(header *types.Header) error {
	if len(header.Vote) == 0 {
		return nil
	}
	if header.Vote == nil {
		return ErrNilVote
	}

	var (
		vb       headergov.VoteBytes = header.Vote
		blockNum                     = header.Number.Uint64()
	)

	vote, err := vb.ToVoteData()
	if err != nil {
		logger.Error("ToVoteData error", "num", blockNum, "vote", vb, "err", err)
		return err
	}

	council, err := h.ValSet.GetCouncil(blockNum)
	if err != nil {
		return err
	}

	// check if the voter is in council
	found := false
	for _, addr := range council {
		if addr == vote.Voter() {
			found = true
			break
		}
	}
	if !found {
		return ErrInvalidKeyValue
	}

	// check if Voter is the block proposer.
	// TODO-kaia: derive author from header
	//author, err := h.Chain.Engine().Author(header)
	//if err != nil {
	//	return err
	//}
	//if author != vote.Voter() {
	//	return ErrInvalidVoter
	//}
	return h.checkConsistency(blockNum, vote)
}

// VerifyGov checks the followings:
// (1) governance must be empty in non-epoch block,
// (2) if there are no votes in the previous epoch, governance must be empty,
// (3) if any vote exists in the previous epoch, governance must not be empty,
// (4) the json must not contain unknown fields,
// (5) the parsed json must exactly match the map derived locally from the previous epoch's votes.
func (h *headerGovModule) VerifyGov(header *types.Header) error {
	// (1)
	if header.Number.Uint64()%h.epoch != 0 {
		if len(header.Governance) > 0 {
			logger.Error("governance is not allowed in non-epoch block", "num", header.Number.Uint64())
			return ErrGovInNonEpochBlock
		} else {
			return nil
		}
	}

	// (2), (3)
	expected := h.getExpectedGovernance(header.Number.Uint64())
	if len(header.Governance) == 0 {
		if len(expected.Items()) != 0 {
			return ErrGovVerification
		}

		return nil
	}

	// (4)
	var gb headergov.GovBytes = header.Governance
	actual, err := gb.ToGovData()
	if err != nil {
		logger.Error("DeserializeHeaderGov error", "num", header.Number.Uint64(), "governance", gb, "err", err)
		return err
	}

	// (5)
	if !reflect.DeepEqual(expected, actual) {
		logger.Error("Governance mismatch", "expected", expected, "actual", actual)
		return ErrGovVerification
	}

	return nil
}

func (h *headerGovModule) checkConsistency(blockNum uint64, vote headergov.VoteData) error {
	//nolint:exhaustive
	switch vote.Name() {
	case gov.GovernanceGoverningNode:
		params := h.GetParamSet(blockNum)

		// skip the consistency check if governingMode is non-single.
		if params.GovernanceMode != "single" {
			return nil
		}

		// we'll use blockNum-1 for the blocknumber of GetCouncil since blockNum cannot be available(eg. vote)
		// it's definite that the valSet vote is not included in this block
		// so the council(blockNum - 1) and council(blockNum) should be same
		council, err := h.ValSet.GetCouncil(blockNum - 1)
		if err != nil {
			return err
		}

		for _, addr := range council {
			if addr == params.GoverningNode {
				return nil
			}
		}
		return ErrGovNodeNotInValSetList
	case gov.GovernanceGovParamContract:
		state, err := h.Chain.State()
		if err != nil {
			return err
		}

		acc := state.GetAccount(vote.Value().(common.Address))
		if acc == nil {
			return ErrGovParamNotAccount
		}

		if acc.Type() != account.SmartContractAccountType || acc.Empty() {
			return ErrGovParamNotContract
		}
	case gov.Kip71LowerBoundBaseFee:
		params := h.GetParamSet(blockNum)
		if vote.Value().(uint64) > params.UpperBoundBaseFee {
			return ErrLowerBoundBaseFee
		}
	case gov.Kip71UpperBoundBaseFee:
		params := h.GetParamSet(blockNum)
		if vote.Value().(uint64) < params.LowerBoundBaseFee {
			return ErrUpperBoundBaseFee
		}
	case gov.AddValidator, gov.RemoveValidator:
		params := h.GetParamSet(blockNum)

		if params.GovernanceMode != "single" {
			return nil
		}
		for _, address := range vote.Value().([]common.Address) {
			if address == params.GoverningNode {
				return ErrGovNodeInValSetVoteValue
			}
		}
	}

	return nil
}

// The blockNum's epoch index must be greater than 0. That is, it must be blockNum >= epoch.
func (h *headerGovModule) getExpectedGovernance(blockNum uint64) headergov.GovData {
	prevEpochIdx := calcEpochIdx(blockNum, h.epoch) - 1
	prevEpochVotes := h.getVotesInEpoch(prevEpochIdx)
	govs := make(gov.PartialParamSet)

	sortedVoteBlocks := maps.Keys(prevEpochVotes)
	slices.Sort(sortedVoteBlocks)

	for _, voteBlock := range sortedVoteBlocks {
		vote := prevEpochVotes[voteBlock]
		govs.Add(string(vote.Name()), vote.Value())
	}

	// assert(len(headergov.NewGovData(govs).Items()) == len(govs))
	return headergov.NewGovData(govs)
}

func (h *headerGovModule) getVotesInEpoch(epochIdx uint64) map[uint64]headergov.VoteData {
	pBorder := ReadLowestVoteScannedEpochIdx(h.ChainKv)
	if pBorder == nil {
		panic("lowest vote scanned epoch index must exist")
	}
	border := *pBorder

	if border <= epochIdx {
		logger.Debug("Scanning votes fastpath")
		votes := make(map[uint64]headergov.VoteData)

		h.mu.RLock()
		defer h.mu.RUnlock()
		for blockNum, vote := range h.groupedVotes[epochIdx] {
			votes[blockNum] = vote
		}
		return votes
	} else {
		logger.Debug("Scanning votes slowpath")
		return h.scanAllVotesInEpoch(epochIdx)
	}
}
