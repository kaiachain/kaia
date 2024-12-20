// Copyright 2024 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package valset

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"

	"github.com/kaiachain/kaia/common"
)

type sortableAddressList []common.Address

func (sa sortableAddressList) Len() int {
	return len(sa)
}

func (sa sortableAddressList) Less(i, j int) bool {
	// Sort by the EIP-155 mixed-case checksummed representation. It differs from the lower-case sorted order.
	// This order is somewhat counterintuitive, but keeping it for backward compatibility.
	return strings.Compare(sa[i].String(), sa[j].String()) < 0
}

func (sa sortableAddressList) Swap(i, j int) {
	sa[i], sa[j] = sa[j], sa[i]
}

// AddressSet is an ordered set of addresses. Its internal list is always sorted.
type AddressSet struct {
	list sortableAddressList
	mu   sync.RWMutex
}

func NewAddressSet(addrs []common.Address) *AddressSet {
	list := make(sortableAddressList, len(addrs))
	copy(list, addrs)
	sort.Sort(list)
	return &AddressSet{
		list: list,
	}
}

func (as *AddressSet) String() string {
	as.mu.RLock()
	defer as.mu.RUnlock()

	addrs := make([]string, len(as.list))
	for i, addr := range as.list {
		addrs[i] = addr.Hex()
	}
	return fmt.Sprintf("[%s]", strings.Join(addrs, ","))
}

func (as *AddressSet) Copy() *AddressSet {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return NewAddressSet(as.list)
}

func (as *AddressSet) List() []common.Address {
	as.mu.RLock()
	defer as.mu.RUnlock()
	result := make([]common.Address, len(as.list))
	copy(result, as.list)
	return result
}

func (as *AddressSet) Len() int {
	as.mu.RLock()
	defer as.mu.RUnlock()
	return len(as.list)
}

func (as *AddressSet) At(i int) common.Address {
	as.mu.RLock()
	defer as.mu.RUnlock()
	if i < 0 {
		return common.Address{}
	}
	return as.list[i%len(as.list)]
}

func (as *AddressSet) IndexOf(addr common.Address) int {
	as.mu.RLock()
	defer as.mu.RUnlock()
	for i, a := range as.list {
		if a == addr {
			return i
		}
	}
	return -1
}

func (as *AddressSet) Contains(addr common.Address) bool {
	return as.IndexOf(addr) != -1
}

func (as *AddressSet) Add(addr common.Address) {
	as.mu.Lock()
	defer as.mu.Unlock()

	for _, a := range as.list {
		if a == addr {
			return
		}
	}
	as.list = append(as.list, addr)
	sort.Sort(as.list)
}

func (as *AddressSet) Remove(addr common.Address) bool {
	as.mu.Lock()
	defer as.mu.Unlock()
	for i, a := range as.list {
		if a == addr {
			as.list = append(as.list[:i], as.list[i+1:]...)
			return true
		}
	}
	return false
}

func (as *AddressSet) Subtract(other *AddressSet) *AddressSet {
	as.mu.RLock()
	defer as.mu.RUnlock()

	m := make(map[common.Address]bool)
	for _, addr := range as.list {
		m[addr] = true
	}
	for _, addr := range other.list {
		delete(m, addr)
	}

	result := make([]common.Address, 0, len(m))
	for addr := range m {
		result = append(result, addr)
	}
	return NewAddressSet(result)
}

// ShuffledListLegacy returns a shuffled list using the legacy shuffling algorithm.
// This is used for backward compatibility with the old committee & proposer selection.
func (as *AddressSet) ShuffledListLegacy(seed int64) []common.Address {
	as.mu.RLock()
	defer as.mu.RUnlock()

	shuffled := make([]common.Address, len(as.list))
	copy(shuffled, as.list)

	r := rand.New(rand.NewSource(seed))
	for i := range shuffled {
		j := r.Intn(len(shuffled))
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

// ShuffledList returns a shuffled list using the Fisher-Yates algorithm.
func (as *AddressSet) ShuffledList(seed int64) []common.Address {
	as.mu.RLock()
	defer as.mu.RUnlock()

	shuffled := make([]common.Address, len(as.list))
	copy(shuffled, as.list)

	r := rand.New(rand.NewSource(seed))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

// HashToSeedLegacy returns the first 15 nibbles (7.5 bytes) of the hash as a seed.
// This is used for backward compatibility with the old committee & proposer selection.
func HashToSeedLegacy(hash common.Hash) int64 {
	n8 := binary.BigEndian.Uint64(hash[:8])
	return int64(n8 >> 4)
}

// HashToSeed returns the first 8 bytes of the hash as a seed.
func HashToSeed(hash []byte) int64 {
	n8 := binary.BigEndian.Uint64(hash[:8])
	return int64(n8)
}
