package impl

import (
	"bytes"
	"errors"
	"reflect"
	"strings"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/rlp"
)

type voteData struct {
	voter common.Address
	name  string
	value []common.Address
}

// NewVoteData returns valset vote
func newVoteData(voter common.Address, name string, value any) (*voteData, error) {
	_, ok := gov.ValSetVoteKeyMap[strings.Trim(strings.ToLower(name), " ")]
	if !ok {
		return nil, errInvalidVoteKey
	}

	v, ok := value.([]common.Address)
	if !ok || len(v) == 0 {
		return nil, errInvalidVoteValue
	}

	return &voteData{
		voter: voter,
		name:  name,
		value: value.([]common.Address),
	}, nil
}

func newVoteDataFromBytes(vb []byte) (string, *voteData, error) {
	var v struct {
		Validator common.Address
		Key       string
		Value     []common.Address
	}
	if err := rlp.DecodeBytes(vb, &v); err != nil {
		return "", nil, err
	}

	vote, err := newVoteData(v.Validator, v.Key, v.Value)
	if err != nil {
		return v.Key, nil, err
	}
	return v.Key, vote, err
}

// Name method
func (v *voteData) Name() string {
	return v.name
}

// Value method
func (v *voteData) Value() []common.Address {
	return v.value
}

func (v *voteData) Equal(v2 *voteData) bool {
	var (
		isVoterEqual = bytes.Equal(v.voter.Bytes(), v2.voter.Bytes())
		isNameEqual  = v.name == v2.name
		isValueEqual = reflect.DeepEqual(v.value, v2.value)
	)
	return isVoterEqual && isNameEqual && isValueEqual
}

// ToVoteBytes method
func (v *voteData) ToVoteBytes() ([]byte, error) {
	vote := &struct {
		Validator common.Address
		Key       string
		Value     []common.Address
	}{
		Validator: v.voter,
		Key:       v.name,
		Value:     v.value,
	}
	return rlp.EncodeToBytes(vote)
}

func (v *ValsetModule) Vote(blockNumber uint64, voter common.Address, name string, value any) (string, error) {
	vote, err := newVoteData(voter, name, value)
	if errors.Is(err, errInvalidVoteKey) {
		return "not valSet vote", nil
	} else if err != nil {
		return "", err
	}

	if err = v.checkConsistency(blockNumber, vote); err != nil {
		return "", err
	}

	v.myVotes = append(v.myVotes, vote)
	return "(kaiax) Your vote has been successfully put into the vote queue. \n" +
		"Your node will proposer the block with this vote. \n" +
		"The new validators will take effect from the next block following the proposed block.", nil
}

func (v *ValsetModule) checkConsistency(blockNumber uint64, vote *voteData) error {
	var duplicateCheckMap map[common.Address]bool

	gp := v.headerGov.EffectiveParamSet(blockNumber + 1)
	if vote.voter != gp.GoverningNode {
		return errInvalidVoter
	}
	for _, address := range vote.Value() {
		// checkConsistency: if it's a single mode, do not include GoverningNode
		if gp.GovernanceMode == "single" && address == gp.GoverningNode {
			return errInvalidVoteValue
		}
		// checkConsistency: do not have duplicated addresses
		if ok, _ := duplicateCheckMap[address]; ok {
			return errInvalidVoteValue
		}
	}
	return nil
}
