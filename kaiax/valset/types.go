package valset

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/kaiachain/kaia/common"
)

type CommonAddressSet interface {
	String() string
	Copy() CommonAddressSet
	Council() []common.Address
	Len() int
	Contains(addr common.Address) bool
	Subtract(other *AddressSet) *AddressSet
	Add(addr common.Address)
	Remove(addr common.Address) bool
	Marshal() ([]byte, error)
	EqualState(NodeStateMap) bool
}

func NewCommonAddressSet(addrs []common.Address) CommonAddressSet {
	return NewAddressSet(addrs)
}

type State uint8

const (
	Unknown      State = iota // 0
	CandInactive              // 1
	CandReady                 // 2
	CandTesting               // 3
	ValInactive               // 4
	ValReady                  // 5
	ValActive                 // 6
	ValPaused                 // 7
	ValExiting                // 8
)

const (
	CandInactiveStr = "CandInactive"
	CandReadyStr    = "CandReady"
	CandTestingStr  = "CandTesting"
	ValInactiveStr  = "ValInactive"
	ValPausedStr    = "ValPaused"
	ValExitingStr   = "ValExiting"
	ValReadyStr     = "ValReady"
	ValActiveStr    = "ValActive"
)

func (s State) ToUint8() uint8 {
	return uint8(s)
}

func (s State) String() string {
	switch s {
	case CandInactive:
		return CandInactiveStr
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

func ParseState(s string) (State, error) {
	switch s {
	case CandInactiveStr:
		return CandInactive, nil
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

func (s State) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", s.String())), nil
}

type ValidatorState struct {
	State         State     `json:"state"`
	StakingAmount uint64    `json:"stakingAmount"` // in KAIA unit
	IdleTimeout   time.Time `json:"idleTimeout"`
	PausedTimeout time.Time `json:"pausedTimeout"`
}

type NodeStateMap map[common.Address]*ValidatorState

func (v NodeStateMap) String() string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("NodeStateMap error: %v", err)
	}
	return string(b)
}

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

func (v NodeStateMap) Addresses() []common.Address {
	return slices.Collect(maps.Keys(v))
}
