// Copyright 2026 The Kaia Authors
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

package vm

import (
	"bytes"
	"sort"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
)

// legacyActivePrecompiles mirrors the pre-prebuilt implementation: collect the
// active contract map's keys, then append the vmversion0 compat addresses for
// istanbul+ rules. Guards the prebuilt slices against fork-behavior drift.
func legacyActivePrecompiles(rules params.Rules) []common.Address {
	contracts := ActivePrecompiledContracts(rules)
	addrs := make([]common.Address, 0, len(contracts))
	for addr := range contracts {
		addrs = append(addrs, addr)
	}
	if rules.IsIstanbul {
		return append(addrs,
			[]common.Address{common.BytesToAddress([]byte{10}), common.BytesToAddress([]byte{11})}...)
	}
	return addrs
}

func sortedAddrs(addrs []common.Address) []common.Address {
	out := append([]common.Address{}, addrs...)
	sort.Slice(out, func(i, j int) bool { return bytes.Compare(out[i][:], out[j][:]) < 0 })
	return out
}

func TestActivePrecompilesMatchesLegacy(t *testing.T) {
	forks := []struct {
		name  string
		rules params.Rules
	}{
		{"byzantium", params.Rules{}},
		{"istanbul", params.Rules{IsIstanbul: true}},
		{"kore", params.Rules{IsIstanbul: true, IsKore: true}},
		{"shanghai", params.Rules{IsIstanbul: true, IsKore: true, IsShanghai: true}},
		{"cancun", params.Rules{IsIstanbul: true, IsKore: true, IsShanghai: true, IsCancun: true}},
		{"prague", params.Rules{IsIstanbul: true, IsKore: true, IsShanghai: true, IsCancun: true, IsPrague: true}},
		{"osaka", params.Rules{IsIstanbul: true, IsKore: true, IsShanghai: true, IsCancun: true, IsPrague: true, IsOsaka: true}},
	}
	for _, f := range forks {
		got := sortedAddrs(ActivePrecompiles(f.rules))
		want := sortedAddrs(legacyActivePrecompiles(f.rules))
		if len(got) != len(want) {
			t.Fatalf("%s: address count mismatch: got %d, want %d", f.name, len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Fatalf("%s: address[%d] mismatch: got %v, want %v", f.name, i, got[i], want[i])
			}
		}
	}
}
