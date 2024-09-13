package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/rlp"
)

type GovData interface {
	Items() map[ParamEnum]interface{}
	Serialize() ([]byte, error)
}

var _ GovData = (*govData)(nil)

type govData struct {
	items map[ParamEnum]interface{}
}

// NewGovData returns a canonical & formatted gov data. VoteForbidden flag and consistency is NOT checked.
// In genesis, forbidden-vote params can exist. Thus, unlike NewVoteData, here we must not check VoteForbidden flag.
func NewGovData(m map[ParamEnum]interface{}) GovData {
	items := make(map[ParamEnum]interface{})
	for ty, value := range m {
		param, ok := Params[ty]
		if !ok {
			return nil
		}

		cv, err := param.Canonicalizer(value)
		if err != nil {
			return nil
		}

		if !param.FormatChecker(cv) {
			return nil
		}

		items[ty] = cv
	}
	return &govData{
		items: items,
	}
}

func (g *govData) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for ty, value := range g.items {
		if bigInt, ok := value.(*big.Int); ok {
			tmp[Params[ty].Name] = bigInt.String()
		} else {
			tmp[Params[ty].Name] = value
		}
	}

	return json.Marshal(tmp)
}

func (g *govData) Items() map[ParamEnum]interface{} {
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

	strMap := make(map[string]interface{})
	err = json.Unmarshal(rlpDecoded, &strMap)
	if err != nil {
		return nil, err
	}

	for name, value := range strMap {
		param, err := GetParamByName(name)
		if err != nil {
			return nil, err
		}

		cv, err := param.Canonicalizer(value)
		if err != nil {
			return nil, err
		}
		strMap[name] = cv
	}

	gov := NewGovData(strMapToEnumMap(strMap))
	if gov == nil {
		return nil, errors.New("failed to create gov data")
	}

	return gov, nil
}

func strMapToEnumMap(strMap map[string]interface{}) map[ParamEnum]interface{} {
	ret := make(map[ParamEnum]interface{})
	for name, value := range strMap {
		ret[ParamNameToEnum[name]] = value
	}
	return ret
}
