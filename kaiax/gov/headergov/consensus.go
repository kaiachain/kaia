package headergov

import (
	"bytes"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
)

func (h *headerGovModule) VerifyHeader(header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil
	}

	// 1. Check Vote
	if len(header.Vote) > 0 {
		vote, err := DeserializeHeaderVote(header.Vote)
		if err != nil {
			logger.Error("Failed to parse vote", "num", header.Number.Uint64(), "err", err)
			return err
		}

		err = h.VerifyVote(header.Number.Uint64(), vote)
		if err != nil {
			logger.Error("Failed to verify vote", "num", header.Number.Uint64(), "err", err)
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

	actual, err := DeserializeHeaderGov(header.Governance)
	if err != nil {
		logger.Error("Failed to parse governance", "num", header.Number.Uint64(), "err", err)
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
		header.Vote, _ = h.myVotes[0].Serialize()
	}

	if header.Number.Uint64()%h.epoch == 0 {
		gov := h.getExpectedGovernance(header.Number.Uint64())
		if len(gov.Items()) > 0 {
			header.Governance, _ = gov.Serialize()
		}
	}

	return nil
}

func (h *headerGovModule) FinalizeHeader(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	return nil
}

// VerifyVote takes canonical VoteData and performs the semantic check.
func (h *headerGovModule) VerifyVote(blockNum uint64, vote VoteData) error {
	if vote == nil {
		return ErrNilVote
	}

	// consistency check
	switch vote.Enum() {
	case GovernanceGoverningNode:
		// TODO: check in valset
		break
	case GovernanceGovParamContract:
		state, err := h.Chain.State()
		if err != nil {
			return err
		}

		acc := state.GetAccount(vote.Value().(common.Address))
		if acc == nil {
			return ErrGovParamNotAccount
		}

		pa := account.GetProgramAccount(acc)
		emptyCodeHash := crypto.Keccak256(nil)
		if pa != nil && !bytes.Equal(pa.GetCodeHash(), emptyCodeHash) {
			return ErrGovParamNotContract
		}
	case Kip71LowerBoundBaseFee:
		params, err := h.EffectiveParamSet(blockNum)
		if err != nil {
			return err
		}
		if vote.Value().(uint64) > params.UpperBoundBaseFee {
			return ErrLowerBoundBaseFee
		}
	case Kip71UpperBoundBaseFee:
		params, err := h.EffectiveParamSet(blockNum)
		if err != nil {
			return err
		}
		if vote.Value().(uint64) < params.LowerBoundBaseFee {
			return ErrUpperBoundBaseFee
		}
	}

	return nil
}

// blockNum must be greater than epoch.
func (h *headerGovModule) getExpectedGovernance(blockNum uint64) GovData {
	prevEpochIdx := calcEpochIdx(blockNum, h.epoch) - 1
	prevEpochVotes := h.getVotesInEpoch(prevEpochIdx)
	govs := make(map[ParamEnum]interface{})

	for _, vote := range prevEpochVotes {
		govs[vote.Enum()] = vote.Value()
	}

	return NewGovData(govs)
}

func (h *headerGovModule) getVotesInEpoch(epochIdx uint64) map[uint64]VoteData {
	votes := make(map[uint64]VoteData)
	for blockNum, vote := range h.cache.GroupedVotes()[epochIdx] {
		votes[blockNum] = vote
	}
	return votes
}
