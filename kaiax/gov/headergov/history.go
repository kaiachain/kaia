package headergov

import (
	"sort"

	"github.com/kaiachain/kaia/kaiax/gov"
)

type History map[uint64]gov.ParamSet

// GetHistory generates history based on sorted gov blocks.
func GetHistory(govs map[uint64]GovData) History {
	gh := make(map[uint64]gov.ParamSet)

	// we must ensure that gov history is not empty
	gh[0] = *gov.GetDefaultGovernanceParamSet()

	sortedNums := make([]uint64, 0, len(govs))
	for num := range govs {
		sortedNums = append(sortedNums, num)
	}
	sort.Slice(sortedNums, func(i, j int) bool {
		return sortedNums[i] < sortedNums[j]
	})

	gp := *gov.GetDefaultGovernanceParamSet()
	for _, num := range sortedNums {
		govData := govs[num]
		if err := gp.SetFromEnumMap(govData.Items()); err != nil {
			continue
		}
		gh[num] = gp
	}

	return gh
}

// Search finds the maximum gov block number that is less than or equal to the given block number.
func (g *History) Search(blockNum uint64) (gov.ParamSet, error) {
	idx := uint64(0)
	for num := range *g {
		if idx < num && num <= blockNum {
			idx = num
		}
	}
	if ret, ok := (*g)[idx]; ok {
		return ret, nil
	}

	panic("must not happen")
}
