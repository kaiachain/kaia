package headergov

import (
	"encoding/json"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/rlp"
)

type VoteBytes []byte
type VotesInEpoch map[uint64]VoteData
type GroupedVotesMap map[uint64]VotesInEpoch

type voteData struct {
	voter common.Address
	name  gov.ParamName
	value any // canonicalized value
}

// NewVoteData returns a valid, canonical vote data.
// If return is not nil, the name and the value is valid.
// The format of the value is checked, but consistency is NOT checked.
func NewVoteData(voter common.Address, name string, value any) VoteData {
	param, ok := gov.Params[gov.ParamName(name)]
	if !ok {
		if name == "governance.addvalidator" || name == "governance.removevalidator" {
			return &voteData{
				voter: voter,
				name:  gov.ParamName(name),
				value: []common.Address{},
			}
		}

		return nil
	}

	if param.VoteForbidden {
		return nil
	}

	cv, err := param.Canonicalizer(value)
	if err != nil {
		return nil
	}

	if !param.FormatChecker(cv) {
		return nil
	}

	return &voteData{
		voter: voter,
		name:  gov.ParamName(name),
		value: cv,
	}
}

func (vote *voteData) Voter() common.Address {
	return vote.voter
}

func (vote *voteData) Name() gov.ParamName {
	return vote.name
}

func (vote *voteData) Value() any {
	return vote.value
}

func (vote *voteData) ToVoteBytes() (VoteBytes, error) {
	v := &struct {
		Validator common.Address
		Key       string
		Value     any
	}{
		Validator: vote.voter,
		Key:       string(vote.name),
		Value:     vote.value,
	}

	if cv, ok := vote.value.(*big.Int); ok {
		v.Value = cv.String()
	}

	return rlp.EncodeToBytes(v)
}

func (vote *voteData) MarshalJSON() ([]byte, error) {
	v := &struct {
		Voter common.Address
		Name  string
		Value any
	}{
		Voter: vote.voter,
		Name:  string(vote.name),
		Value: vote.value,
	}

	return json.Marshal(v)
}

func (vb VoteBytes) ToVoteData() (VoteData, error) {
	var v struct {
		Validator common.Address
		Key       string
		Value     []byte
	}

	err := rlp.DecodeBytes(vb, &v)
	if err != nil {
		return nil, ErrInvalidRlp
	}

	vote := NewVoteData(v.Validator, v.Key, v.Value)
	if vote == nil {
		return nil, ErrInvalidVoteData
	}

	return vote, nil
}

func (vb VoteBytes) String() string {
	return hexutil.Encode(vb)
}
