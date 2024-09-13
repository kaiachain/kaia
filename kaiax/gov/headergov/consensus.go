package headergov

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
)

var (
	errGovInNonEpochBlock = errors.New("governance is not allowed in non-epoch block")
	errNilVote            = errors.New("vote is nil")
	errInvalidGov         = errors.New("header.Governance does not match the vote in previous epoch")

	errGovParamNotAccount  = errors.New("govparamcontract is not an account")
	errGovParamNotContract = errors.New("govparamcontract is not an contract account")
	errLowerBoundBaseFee   = errors.New("lowerboundbasefee is greater than upperboundbasefee")
	errUpperBoundBaseFee   = errors.New("upperboundbasefee is less than lowerboundbasefee")
)

func (h *headerGovModule) VerifyHeader(header *types.Header) error {
	if header.Number.Uint64() == 0 {
		return nil
	}

	// 1. Check Vote
	if len(header.Vote) > 0 {
		vote, err := DeserializeHeaderVote(header.Vote, header.Number.Uint64())
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
			return errGovInNonEpochBlock
		} else {
			return nil
		}
	}

	gov, err := DeserializeHeaderGov(header.Governance, header.Number.Uint64())
	if err != nil {
		logger.Error("Failed to parse governance", "num", header.Number.Uint64(), "err", err)
		return err
	}
	return h.VerifyGov(header.Number.Uint64(), gov)
}

func (h *headerGovModule) PrepareHeader(header *types.Header) (*types.Header, error) {
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

	return header, nil
}

func (h *headerGovModule) FinalizeBlock(b *types.Block) (*types.Block, error) {
	return b, nil
}

// VerifyVote takes canonical VoteData and performs the semantic check.
func (h *headerGovModule) VerifyVote(blockNum uint64, vote VoteData) error {
	if vote == nil {
		return errNilVote
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
			return errGovParamNotAccount
		}

		pa := account.GetProgramAccount(acc)
		emptyCodeHash := crypto.Keccak256(nil)
		if pa != nil && !bytes.Equal(pa.GetCodeHash(), emptyCodeHash) {
			return errGovParamNotContract
		}
	case Kip71LowerBoundBaseFee:
		params, err := h.EffectiveParamSet(blockNum)
		if err != nil {
			return err
		}
		if vote.Value().(uint64) > params.UpperBoundBaseFee {
			return errLowerBoundBaseFee
		}
	case Kip71UpperBoundBaseFee:
		params, err := h.EffectiveParamSet(blockNum)
		if err != nil {
			return err
		}
		if vote.Value().(uint64) < params.LowerBoundBaseFee {
			return errUpperBoundBaseFee
		}
	}

	return nil
}

func (h *headerGovModule) VerifyGov(blockNum uint64, gov GovData) error {
	expected := h.getExpectedGovernance(blockNum)
	if !reflect.DeepEqual(expected, gov) {
		return errInvalidGov
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
