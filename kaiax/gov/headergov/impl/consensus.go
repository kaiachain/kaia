package impl

import (
	"reflect"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
)

func (h *headerGovModule) VerifyHeader(header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil
	}

	// 1. Check Vote
	if len(header.Vote) > 0 {
		var vb headergov.VoteBytes = header.Vote
		vote, err := vb.ToVoteData()
		if err != nil {
			logger.Error("ToVoteData error", "num", header.Number.Uint64(), "vote", vb, "err", err)
			return err
		}

		err = h.VerifyVote(header.Number.Uint64(), vote)
		if err != nil {
			logger.Error("VerifyVote error", "num", header.Number.Uint64(), "vote", vb, "err", err)
			return err
		}
	}

	// 2. Check Governance
	if header.Number.Uint64()%h.epoch != 0 {
		if len(header.Governance) > 0 {
			logger.Error("governance is not allowed in non-epoch block", "num", header.Number.Uint64())
			return ErrGovInNonEpochBlock
		} else {
			return nil
		}
	}

	expected := h.getExpectedGovernance(header.Number.Uint64())
	if len(header.Governance) == 0 {
		if len(expected.Items()) != 0 {
			return ErrGovVerification
		}

		return nil
	}

	var gb headergov.GovBytes = header.Governance
	actual, err := gb.ToGovData()
	if err != nil {
		logger.Error("DeserializeHeaderGov error", "num", header.Number.Uint64(), "governance", gb, "err", err)
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		logger.Error("Governance mismatch", "expected", expected, "actual", actual)
		return ErrGovVerification
	}

	return nil
}

func (h *headerGovModule) PrepareHeader(header *types.Header) error {
	// if epoch block & vote exists in the last epoch, put Governance to header.
	if len(h.myVotes) > 0 {
		header.Vote, _ = h.myVotes[0].ToVoteBytes()
	}

	if header.Number.Uint64()%h.epoch == 0 {
		gov := h.getExpectedGovernance(header.Number.Uint64())
		if len(gov.Items()) > 0 {
			header.Governance, _ = gov.ToGovBytes()
		}
	}

	return nil
}

func (h *headerGovModule) FinalizeHeader(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	return nil
}

// VerifyVote checks the followings:
// (1) if voter is in valset,
// (2) integrity of the voter (ensures that voter is the block proposer),
// (3) consistency check of the vote value.
func (h *headerGovModule) VerifyVote(blockNum uint64, vote headergov.VoteData) error {
	if vote == nil {
		return ErrNilVote
	}

	// TODO: check if if voter is in valset.
	// TODO: check if Voter is the block proposer.

	// consistency check
	return h.checkConsistency(blockNum, vote)
}

func (h *headerGovModule) checkConsistency(blockNum uint64, vote headergov.VoteData) error {
	switch vote.Name() {
	case gov.GovernanceGoverningNode:
		// TODO: check in valset
		break
	case gov.GovernanceGovParamContract:
		state, err := h.Chain.State()
		if err != nil {
			return err
		}

		acc := state.GetAccount(vote.Value().(common.Address))
		if acc == nil {
			return ErrGovParamNotAccount
		}

		pa := account.GetProgramAccount(acc)
		if pa == nil || pa.Empty() {
			return ErrGovParamNotContract
		}
	case gov.Kip71LowerBoundBaseFee:
		params := h.EffectiveParamSet(blockNum)
		if vote.Value().(uint64) > params.UpperBoundBaseFee {
			return ErrLowerBoundBaseFee
		}
	case gov.Kip71UpperBoundBaseFee:
		params := h.EffectiveParamSet(blockNum)
		if vote.Value().(uint64) < params.LowerBoundBaseFee {
			return ErrUpperBoundBaseFee
		}
	}

	return nil
}

// blockNum must be greater than epoch.
func (h *headerGovModule) getExpectedGovernance(blockNum uint64) headergov.GovData {
	prevEpochIdx := calcEpochIdx(blockNum, h.epoch) - 1
	prevEpochVotes := h.getVotesInEpoch(prevEpochIdx)
	govs := make(map[gov.ParamName]any)

	for _, vote := range prevEpochVotes {
		govs[vote.Name()] = vote.Value()
	}

	return headergov.NewGovData(govs)
}

func (h *headerGovModule) getVotesInEpoch(epochIdx uint64) map[uint64]headergov.VoteData {
	votes := make(map[uint64]headergov.VoteData)
	for blockNum, vote := range h.cache.GroupedVotes()[epochIdx] {
		votes[blockNum] = vote
	}
	return votes
}
