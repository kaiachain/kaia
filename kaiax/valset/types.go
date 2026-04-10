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

type State uint8

const (
	Unknown     State = iota // 0
	Registered               // 1
	CandReady                // 2
	CandTesting              // 3
	ValInactive              // 4
	ValReady                 // 5
	ValActive                // 6
	ValPaused                // 7
	ValExiting               // 8
)

// State groups for permissionless consensus.
// Validator = ValActive + ValReady + ValPaused + ValExiting + ValInactive
// Council = ValActive + ValReady + ValPaused
// Committee = ValActive
// Candidate = CandTesting
var (
	ValidatorStates = []State{ValActive, ValReady, ValPaused, ValExiting, ValInactive}
	CouncilStates   = []State{ValActive, ValReady, ValPaused}
	CommitteeStates = []State{ValActive}
	CandidateStates = []State{CandTesting}
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
func (s State) ToUint8() uint8 {
	return uint8(s)
}

// String returns the human-readable name of the state.
func (s State) String() string {
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
func ParseState(s string) (State, error) {
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
func (s *State) UnmarshalJSON(b []byte) error {
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
func (s State) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", s.String())), nil
}

// ValidatorState represents the current state and metadata of a validator node.
type ValidatorState struct {
	State         State     `json:"state"`
	StakingAmount uint64    `json:"stakingAmount"` // in KAIA unit
	IdleTimeout   time.Time `json:"idleTimeout"`
	PausedTimeout time.Time `json:"pausedTimeout"`
	Suspended     bool      `json:"suspended"`
}

// NodeStateMap maps validator addresses to their state. Used for permissionless state transitions.
type NodeStateMap map[common.Address]*ValidatorState

func (v NodeStateMap) String() string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("NodeStateMap error: %v", err)
	}
	return string(b)
}

// Copy returns a deep copy of the NodeStateMap.
func (v NodeStateMap) Copy() NodeStateMap {
	if v == nil {
		return nil
	}

	cp := make(NodeStateMap, len(v))
	for key, value := range v {
		if value == nil {
			cp[key] = nil
			continue
		}
		newValue := *value
		cp[key] = &newValue
	}
	return cp
}

// EqualState returns true if both maps have the same addresses with matching State values.
// Only compares State field, ignoring StakingAmount and timeouts.
func (v NodeStateMap) EqualState(other NodeStateMap) bool {
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
func (v NodeStateMap) Addresses() []common.Address {
	addrs := slices.Collect(maps.Keys(v))
	slices.SortFunc(addrs, func(a, b common.Address) int {
		return bytes.Compare(a[:], b[:])
	})
	return addrs
}

// CountByState counts the number of validators in a given state.
func (v NodeStateMap) CountByState(state State) uint64 {
	var count uint64
	for _, val := range v {
		if val.State == state {
			count++
		}
	}
	return count
}

// MarkSuspended sets the Suspended flag on validators whose address is in the given list.
func (v NodeStateMap) MarkSuspended(suspended []common.Address) {
	set := make(map[common.Address]struct{}, len(suspended))
	for _, addr := range suspended {
		set[addr] = struct{}{}
	}
	for addr, vs := range v {
		_, vs.Suspended = set[addr]
	}
}

// ExcludeSuspended returns a new NodeStateMap without suspended validators.
func (v NodeStateMap) ExcludeSuspended() NodeStateMap {
	filtered := make(NodeStateMap, len(v))
	for addr, val := range v {
		if !val.Suspended {
			filtered[addr] = val
		}
	}
	return filtered
}
