package headergov

import (
	"encoding/json"
	"math/big"

	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/rlp"
)

type govData struct {
	items map[gov.ParamEnum]interface{}
}

// NewGovData returns a canonical & formatted gov data. VoteForbidden flag and consistency is NOT checked.
// In genesis, forbidden-vote params can exist. Thus, unlike NewVoteData, here we must not check VoteForbidden flag.
func NewGovData(m map[gov.ParamEnum]interface{}) GovData {
	items := make(map[gov.ParamEnum]interface{})
	for enum, value := range m {
		param, ok := gov.Params[enum]
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

		items[enum] = cv
	}
	return &govData{
		items: items,
	}
}

func (g *govData) MarshalJSON() ([]byte, error) {
	tmp := make(map[string]interface{})
	for enum, value := range g.items {
		if bigInt, ok := value.(*big.Int); ok {
			tmp[gov.Params[enum].Name] = bigInt.String()
		} else {
			tmp[gov.Params[enum].Name] = value
		}
	}

	return json.Marshal(tmp)
}

func (g *govData) Items() map[gov.ParamEnum]interface{} {
	return g.items
}

func (g *govData) Serialize() ([]byte, error) {
	j, err := g.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(j)
}

func DeserializeHeaderGov(b []byte) (GovData, error) {
	rlpDecoded := []byte("")
	err := rlp.DecodeBytes(b, &rlpDecoded)
	if err != nil {
		return nil, ErrInvalidRlp
	}

	strMap := make(map[string]interface{})
	err = json.Unmarshal(rlpDecoded, &strMap)
	if err != nil {
		return nil, ErrInvalidJson
	}

	for name, value := range strMap {
		param, err := gov.GetParamByName(name)
		if err != nil {
			return nil, err
		}

		cv, err := param.Canonicalizer(value)
		if err != nil {
			return nil, err
		}
		strMap[name] = cv
	}

	gov := NewGovData(gov.StrMapToEnumMap(strMap))
	if gov == nil {
		return nil, ErrInvalidGovData
	}

	return gov, nil
}
