package valset

import (
	"encoding/json"
	"fmt"
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
	EqualState(ValidatorStateMap) bool
	GetDemoted(CommonAddressSet, map[common.Address]float64, uint64) *AddressSet
}

func NewCommonAddressSet(addrs []common.Address) CommonAddressSet {
	return NewAddressSet(addrs)
}

type State uint8

const (
	Unknown State = iota
	CandInactive
	CandReady
	CandTesting
	ValInactive
	ValPaused
	ValExiting
	ValReady
	ValActive
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

type ValidatorStateMap map[common.Address]*ValidatorState

func (v ValidatorStateMap) String() string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("ValidatorStateMap error: %v", err)
	}
	return string(b)
}

func (v ValidatorStateMap) Copy() ValidatorStateMap {
	if v == nil {
		return nil
	}

	cp := make(ValidatorStateMap, len(v))
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

func (v ValidatorStateMap) EqualState(other ValidatorStateMap) bool {
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
