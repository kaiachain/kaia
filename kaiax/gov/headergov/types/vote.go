package types

import (
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
)

type VoteData interface {
	Voter() common.Address
	Name() string
	Type() ParamEnum
	Value() interface{}

	Serialize() ([]byte, error)
}

var _ VoteData = (*voteData)(nil)

type voteData struct {
	voter common.Address
	name  string
	ty    ParamEnum
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
		ty:    paramNameToEnum[name],
		value: cv,
	}
}

func (vote *voteData) Voter() common.Address {
	return vote.voter
}

func (vote *voteData) Name() string {
	return vote.name
}

func (vote *voteData) Type() ParamEnum {
	return vote.ty
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

func DeserializeHeaderVote(b []byte, blockNum uint64) (VoteData, error) {
	var v struct {
		Validator common.Address
		Key       string
		Value     []byte
	}

	err := rlp.DecodeBytes(b, &v)
	if err != nil {
		return nil, err
	}

	vote := NewVoteData(v.Validator, v.Key, v.Value)
	if vote == nil {
		return nil, ErrInvalidVoteData
	}

	return vote, nil
}
