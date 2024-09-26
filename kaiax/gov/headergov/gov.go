package headergov

import (
	"encoding/json"
	"math/big"

	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/rlp"
)

type GovBytes []byte

type govData struct {
	items map[gov.ParamEnum]any
}

// NewGovData returns a canonical & formatted gov data. VoteForbidden flag and consistency is NOT checked.
// In genesis, forbidden-vote params can exist. Thus, unlike NewVoteData, here we must not check VoteForbidden flag.
func NewGovData(m map[gov.ParamEnum]any) GovData {
	items := make(map[gov.ParamEnum]any)
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
	tmp := make(map[string]any)
	for enum, value := range g.items {
		if bigInt, ok := value.(*big.Int); ok {
			tmp[gov.Params[enum].Name] = bigInt.String()
		} else {
			tmp[gov.Params[enum].Name] = value
		}
	}

	return json.Marshal(tmp)
}

func (g *govData) Items() map[gov.ParamEnum]any {
	return g.items
}

func (g *govData) ToGovBytes() (GovBytes, error) {
	j, err := g.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(j)
}

func (gb GovBytes) ToGovData() (GovData, error) {
	rlpDecoded := []byte("")
	err := rlp.DecodeBytes(gb, &rlpDecoded)
	if err != nil {
		return nil, ErrInvalidRlp
	}

	strMap := make(map[string]any)
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

func (gb GovBytes) String() string {
	return hexutil.Encode(gb)
}
