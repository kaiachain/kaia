package valset

import (
	"bytes"
	"sort"
	"strings"

	"github.com/kaiachain/kaia/common"
)

type AddressList []common.Address

func (sc AddressList) Len() int {
	return len(sc)
}

func (sc AddressList) Less(i, j int) bool {
	return strings.Compare(sc[i].String(), sc[j].String()) < 0
}

func (sc AddressList) Swap(i, j int) {
	sc[i], sc[j] = sc[j], sc[i]
}

func (sc AddressList) AddressStringList() []string {
	stringAddrs := make([]string, len(sc))
	for _, val := range sc {
		stringAddrs = append(stringAddrs, val.String())
	}
	return stringAddrs
}

func (sc AddressList) GetIdxByAddress(addr common.Address) int {
	for i, val := range sc {
		if addr == val {
			return i
		}
	}
	// TODO-Kaia-Istanbul: Enable this log when non-committee nodes don't call `core.startNewRound()`
	// logger.Warn("failed to find an address in the validator list",
	// 	"address", addr, "validatorAddrs", valSet.validators.AddressStringList())
	return -1
}

// SortedAddressList retrieves the sorted address list of ValidatorSet in "ascending order".
// if public is false, sort it using bytes.Compare. It's for public purpose.
// - public-false usage: (getValidators/getDemotedValidators, defaultSet snap store, prepareExtra.validators)
// if public is true, sort it using strings.Compare. It's used for internal consensus purpose, especially for the source of committee.
// - public-true usage: (snap read/store/apply except defaultSet snap store, vrank log)
// TODO-kaia-valset: unify sorting.
func (sc AddressList) SortedAddressList(public bool) []common.Address {
	copiedSlice := make(AddressList, len(sc))
	copy(copiedSlice, sc)

	if public {
		// want reverse-sort: ascending order - bytes.Compare(ValidatorSet[i][:], ValidatorSet[j][:]) > 0
		sort.Slice(copiedSlice, func(i, j int) bool {
			return bytes.Compare(copiedSlice[i].Bytes(), copiedSlice[j].Bytes()) >= 0
		})
		sort.Sort(sort.Reverse(copiedSlice))
	} else {
		// want sort: descending order - strings.Compare(ValidatorSet[i].String(), ValidatorSet[j].String()) < 0
		sort.Sort(copiedSlice)
	}
	return copiedSlice
}
