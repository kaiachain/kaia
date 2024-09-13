package types

import (
	"encoding/json"
	"math/big"

	"github.com/kaiachain/kaia/rlp"
)

type GovData interface {
	Items() map[string]interface{}
	Serialize() ([]byte, error)
}

var _ GovData = (*govData)(nil)

type govData struct {
	items map[string]interface{}
}

// NewGovData returns a canonical & formatted gov data. VoteForbidden flag and consistency is NOT checked.
// In genesis, forbidden-vote params can exist. Thus, unlike NewVoteData, here we must not check VoteForbidden flag.
func NewGovData(m map[string]interface{}) GovData {
	items := make(map[string]interface{})
	for name, value := range m {
		param, err := GetParamByName(name)
		if err != nil {
			return nil
		}

		cv, err := param.Canonicalizer(value)
		if err != nil {
			return nil
		}

		if !param.FormatChecker(cv) {
			return nil
		}

		items[name] = cv
	}
	return &govData{
		items: items,
	}
}

func (g *govData) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for name, value := range g.items {
		if bigInt, ok := value.(*big.Int); ok {
			tmp[name] = bigInt.String()
		} else {
			tmp[name] = value
		}
	}

	return json.Marshal(tmp)
}

func (g *govData) Items() map[string]interface{} {
	return g.items
}

func (g *govData) Serialize() ([]byte, error) {
	j, err := g.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(j)
}

func DeserializeHeaderGov(b []byte, blockNum uint64) (GovData, error) {
	rlpDecoded := []byte("")
	err := rlp.DecodeBytes(b, &rlpDecoded)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	err = json.Unmarshal(rlpDecoded, &ret)
	if err != nil {
		return nil, err
	}

	for name, value := range ret {
		param, err := GetParamByName(name)
		if err != nil {
			return nil, err
		}

		cv, err := param.Canonicalizer(value)
		if err != nil {
			return nil, err
		}
		ret[name] = cv
	}

	return &govData{
		items: ret,
	}, nil
}
