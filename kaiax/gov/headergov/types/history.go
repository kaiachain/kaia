package types

import (
	"sort"

	gov_types "github.com/kaiachain/kaia/kaiax/gov/types"
)

type History map[uint64]gov_types.ParamSet

func GetHistory(govs map[uint64]GovData) History {
	gh := make(map[uint64]gov_types.ParamSet)

	var sortedNums []uint64
	for num := range govs {
		sortedNums = append(sortedNums, num)
	}
	sort.Slice(sortedNums, func(i, j int) bool {
		return sortedNums[i] < sortedNums[j]
	})

	gp := gov_types.ParamSet{}
	for _, num := range sortedNums {
		govData := govs[num]
		gp.SetFromEnumMap(govData.Items())
		gh[num] = gp
	}
	return gh
}

func (g *History) Search(blockNum uint64) (gov_types.ParamSet, error) {
	idx := uint64(0)
	for num := range *g {
		if idx < num && num <= blockNum {
			idx = num
		}
	}
	if ret, ok := (*g)[idx]; ok {
		return ret, nil
	} else {
		return gov_types.ParamSet{}, ErrNotFoundInHistory
	}
}
