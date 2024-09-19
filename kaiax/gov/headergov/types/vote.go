package types

import (
	"encoding/json"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
)

type voteData struct {
	voter common.Address
	name  string
	enum  ParamEnum
	value interface{} // canonicalized value
}

// NewVoteData returns a valid, canonical vote data.
// If return is not nil, the name and the value is valid.
// The format of the value is checked, but consistency is NOT checked.
func NewVoteData(voter common.Address, name string, value interface{}) VoteData {
	param, err := GetParamByName(name)
	if err != nil {
		if name == "governance.addvalidator" || name == "governance.removevalidator" {
			return &voteData{
				voter: voter,
				name:  name,
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
		name:  name,
		enum:  ParamNameToEnum[name],
		value: cv,
	}
}

func (vote *voteData) Voter() common.Address {
	return vote.voter
}

func (vote *voteData) Name() string {
	return vote.name
}

func (vote *voteData) Enum() ParamEnum {
	return vote.enum
}

func (vote *voteData) Value() interface{} {
	return vote.value
}

func (vote *voteData) Serialize() ([]byte, error) {
	v := &struct {
		Validator common.Address
		Key       string
		Value     interface{}
	}{
		Validator: vote.voter,
		Key:       vote.name,
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
		Enum  ParamEnum
		Value interface{}
	}{
		Voter: vote.voter,
		Name:  vote.name,
		Enum:  vote.enum,
		Value: vote.value,
	}

	return json.Marshal(v)
}

func DeserializeHeaderVote(b []byte) (VoteData, error) {
	var v struct {
		Validator common.Address
		Key       string
		Value     []byte
	}

	err := rlp.DecodeBytes(b, &v)
	if err != nil {
		return nil, ErrInvalidRlp
	}

	vote := NewVoteData(v.Validator, v.Key, v.Value)
	if vote == nil {
		return nil, ErrInvalidVoteData
	}

	return vote, nil
}
