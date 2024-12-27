package headergov

import (
	"encoding/json"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/rlp"
)

type (
	VoteBytes       []byte
	VotesInEpoch    map[uint64]VoteData
	GroupedVotesMap map[uint64]VotesInEpoch
)

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
		param, ok = gov.ValidatorParams[gov.ParamName(name)]
		if !ok {
			logger.Error("Invalid vote name", "name", name)
			return nil
		}
	}

	if param.VoteForbidden {
		logger.Error("Vote is forbidden", "name", name)
		return nil
	}

	cv, err := param.Canonicalizer(value)
	if err != nil {
		logger.Error("Canonicalize error", "name", name, "value", value, "err", err)
		return nil
	}

	if !param.FormatChecker(cv) {
		logger.Error("Format check error", "name", name, "value", value)
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
	} else if cv, ok := vote.value.([]common.Address); ok {
		// concat all addresses into []byte
		concatBytes := make([]byte, 0, len(cv)*common.AddressLength)
		for _, addr := range cv {
			concatBytes = append(concatBytes, addr.Bytes()...)
		}
		v.Value = concatBytes
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
