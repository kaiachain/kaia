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

package valset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/kaiachain/kaia/common"
)

// Genesis defaults for the ABv2 parameters. Chains may configure other values; these are
// what homi writes and what tests build their genesis from.
const (
	DefaultValPausedTimeout        = time.Hour * 8
	DefaultValIdleTimeout          = 30 * 24 * time.Hour
	DefaultMaxNodeCount            = 100
	DefaultMaxValActivePausedCount = 50
	DefaultMaxCandReadyCount       = 3
	DefaultPfsThreshold            = 2
	DefaultCfsThreshold            = 300
)

type NodeState uint8

const (
	Unknown     NodeState = iota // 0
	Registered                   // 1
	CandReady                    // 2
	CandTesting                  // 3
	ValInactive                  // 4
	ValReady                     // 5
	ValActive                    // 6
	ValPaused                    // 7
	ValExiting                   // 8
)

const (
	RegisteredStr  = "Registered"
	CandReadyStr   = "CandReady"
	CandTestingStr = "CandTesting"
	ValInactiveStr = "ValInactive"
	ValPausedStr   = "ValPaused"
	ValExitingStr  = "ValExiting"
	ValReadyStr    = "ValReady"
	ValActiveStr   = "ValActive"
)

// ToUint8 converts the State to its uint8 representation for ABI encoding.
func (s NodeState) ToUint8() uint8 {
	return uint8(s)
}

// IsRewardEligible returns true if the state is eligible for block rewards.
func (s NodeState) IsRewardEligible() bool {
	return s == ValActive
}

// String returns the human-readable name of the state.
func (s NodeState) String() string {
	switch s {
	case Registered:
		return RegisteredStr
	case CandReady:
		return CandReadyStr
	case CandTesting:
		return CandTestingStr
	case ValInactive:
		return ValInactiveStr
	case ValPaused:
		return ValPausedStr
	case ValExiting:
		return ValExitingStr
	case ValReady:
		return ValReadyStr
	case ValActive:
		return ValActiveStr
	default:
		return fmt.Sprintf("UnknownState(%d)", s)
	}
}

// ParseState converts a string to a State value.
func ParseState(s string) (NodeState, error) {
	switch s {
	case RegisteredStr:
		return Registered, nil
	case CandReadyStr:
		return CandReady, nil
	case CandTestingStr:
		return CandTesting, nil
	case ValInactiveStr:
		return ValInactive, nil
	case ValPausedStr:
		return ValPaused, nil
	case ValExitingStr:
		return ValExiting, nil
	case ValReadyStr:
		return ValReady, nil
	case ValActiveStr:
		return ValActive, nil
	default:
		return 0, fmt.Errorf("invalid state string: %s", s)
	}
}

// UnmarshalJSON deserializes a JSON string into a State value.
func (s *NodeState) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	val, err := ParseState(str)
	if err != nil {
		return err
	}
	*s = val
	return nil
}

// MarshalJSON serializes the State as a JSON string.
func (s NodeState) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, "%q", s.String()), nil
}

// Node holds a node's lifecycle state and the metadata that drives transitions.
// Lifecycle: read from ABv2 -> state transition -> written to ABv2
type Node struct {
	State         NodeState `json:"state"`
	StakingAmount uint64    `json:"stakingAmount"` // in KAIA unit
	IdleTimeout   time.Time `json:"idleTimeout"`
	PausedTimeout time.Time `json:"pausedTimeout"`
	Suspended     bool      `json:"suspended"`
}

func (n *Node) Copy() *Node {
	if n == nil {
		return nil
	}
	copied := *n
	return &copied
}

// NodeMap is the keyed collection of per-node records. Used for permissionless state transitions.
type NodeMap map[common.Address]*Node

func (v NodeMap) String() string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("NodeMap error: %v", err)
	}
	return string(b)
}

// Copy returns a deep copy of the NodeMap. Node values are
// duplicated so that mutations through the returned map do not affect the
// receiver. Pure transition methods rely on this property.
func (v NodeMap) Copy() NodeMap {
	if v == nil {
		return nil
	}

	cp := make(NodeMap, len(v))
	for key, value := range v {
		cp[key] = value.Copy()
	}
	return cp
}

// EqualState returns true if both maps have the same addresses with matching State values.
// Only compares State field, ignoring StakingAmount and timeouts.
func (v NodeMap) EqualState(other NodeMap) bool {
	if len(v) != len(other) {
		return false
	}

	for addr, val := range v {
		otherVal, exists := other[addr]
		if !exists {
			return false
		}
		if val == nil || otherVal == nil {
			if val != otherVal {
				return false
			}
			continue
		}
		if val.State != otherVal.State {
			return false
		}
	}
	return true
}

// Addresses returns a deterministically sorted list of addresses in the map.
func (v NodeMap) Addresses() []common.Address {
	addrs := slices.Collect(maps.Keys(v))
	slices.SortFunc(addrs, func(a, b common.Address) int {
		return bytes.Compare(a[:], b[:])
	})
	return addrs
}

// CountByState counts the number of nodes in a given state.
func (v NodeMap) CountByState(state NodeState) uint64 {
	var count uint64
	for _, val := range v {
		if val.State == state {
			count++
		}
	}
	return count
}

// MarkSuspended sets the Suspended flag on nodes whose address is in the given list.
// Mutates the receiver in place; the only in-place mutator on NodeMap.
func (v NodeMap) MarkSuspended(suspended []common.Address) {
	set := make(map[common.Address]struct{}, len(suspended))
	for _, addr := range suspended {
		set[addr] = struct{}{}
	}
	for addr, vs := range v {
		_, vs.Suspended = set[addr]
	}
}

// ExcludeSuspended returns a new NodeMap without suspended nodes.
func (v NodeMap) ExcludeSuspended() NodeMap {
	filtered := make(NodeMap, len(v))
	for addr, val := range v {
		if !val.Suspended {
			filtered[addr] = val
		}
	}
	return filtered
}

// FilterByState returns a new NodeMap containing only nodes in one of the given states.
func (v NodeMap) FilterByState(states ...NodeState) NodeMap {
	filtered := make(NodeMap)
	for addr, val := range v {
		if slices.Contains(states, val.State) {
			filtered[addr] = val
		}
	}
	return filtered
}

// Council returns nodes committed to the Istanbul consensus cycle.
// {ValActive, ValPaused}. ValReady excluded — voluntary standby, not consensus commitment.
func (v NodeMap) Council() NodeMap {
	return v.FilterByState(ValActive, ValPaused)
}

// Committee returns nodes eligible for consensus signing and proposing.
// {ValActive} excluding suspended nodes.
func (v NodeMap) Committee() NodeMap {
	return v.FilterByState(ValActive).ExcludeSuspended()
}

// HeaderGovVoters returns nodes eligible for governance header votes.
// {ValActive} excluding suspended nodes.
func (v NodeMap) HeaderGovVoters() NodeMap {
	return v.FilterByState(ValActive).ExcludeSuspended()
}

// CNPeers returns nodes that should maintain CN-CN P2P connections.
// Includes CandReady — P2P must be established before CandTesting promotion.
func (v NodeMap) CNPeers() NodeMap {
	return v.FilterByState(ValActive, ValReady, ValPaused, CandReady, CandTesting)
}
