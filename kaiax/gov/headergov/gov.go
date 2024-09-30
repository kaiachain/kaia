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
	items gov.PartialParamSet
}

// NewGovData returns a canonical & formatted gov data. VoteForbidden flag and consistency is NOT checked.
// In genesis, forbidden-vote params can exist. Thus, unlike NewVoteData, here we must not check VoteForbidden flag.
func NewGovData(m gov.PartialParamSet) GovData {
	items := make(gov.PartialParamSet)
	for name, value := range m {
		items.Add(string(name), value)
	}
	return &govData{
		items: items,
	}
}

func (g *govData) MarshalJSON() ([]byte, error) {
	tmp := make(gov.PartialParamSet)
	for name, value := range g.items {
		if bigInt, ok := value.(*big.Int); ok {
			tmp[name] = bigInt.String()
		} else {
			tmp[name] = value
		}
	}

	return json.Marshal(tmp)
}

func (g *govData) Items() gov.PartialParamSet {
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

	m := make(gov.PartialParamSet)
	err = json.Unmarshal(rlpDecoded, &m)
	if err != nil {
		return nil, ErrInvalidJson
	}

	for name, value := range m {
		m.Add(string(name), value)
	}

	gov := NewGovData(m)
	if gov == nil {
		return nil, ErrInvalidGovData
	}

	return gov, nil
}

func (gb GovBytes) String() string {
	return hexutil.Encode(gb)
}
