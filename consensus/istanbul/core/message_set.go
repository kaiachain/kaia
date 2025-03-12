// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from quorum/consensus/istanbul/core/message_set.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package core

import (
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
)

// Construct a new message set to accumulate messages for given sequence/view number.
func newMessageSet(valSet *istanbul.BlockValSet) *messageSet {
	return &messageSet{
		view: &istanbul.View{
			Round:    new(big.Int),
			Sequence: new(big.Int),
		},
		messagesMu: new(sync.Mutex),
		messages:   make(map[common.Address]*message),
		valSet:     valSet,
	}
}

// ----------------------------------------------------------------------------

type messageSet struct {
	view       *istanbul.View
	valSet     *istanbul.BlockValSet
	messagesMu *sync.Mutex
	messages   map[common.Address]*message
}

func (ms *messageSet) View() *istanbul.View {
	return ms.view
}

func (ms *messageSet) GetMessages() string {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	ret := "\n"
	for k := range ms.messages {
		ret = ret + "\t" + k.String() + "\n"
	}
	return ret
}

func (ms *messageSet) Add(msg *message) error {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	if err := ms.verify(msg); err != nil {
		return err
	}

	return ms.addVerifiedMessage(msg)
}

func (ms *messageSet) Values() (result []*message) {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()

	for _, v := range ms.messages {
		result = append(result, v)
	}

	return result
}

func (ms *messageSet) Size() int {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	return len(ms.messages)
}

func (ms *messageSet) Get(addr common.Address) *message {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	return ms.messages[addr]
}

// ----------------------------------------------------------------------------

func (ms *messageSet) verify(msg *message) error {
	// verify if the message comes from one of the validators
	if !ms.valSet.Qualified().Contains(msg.Address) {
		return istanbul.ErrUnauthorizedAddress
	}

	// TODO: check view number and sequence number

	return nil
}

func (ms *messageSet) addVerifiedMessage(msg *message) error {
	ms.messages[msg.Address] = msg
	return nil
}

func (ms *messageSet) String() string {
	ms.messagesMu.Lock()
	defer ms.messagesMu.Unlock()
	addresses := make([]string, 0, len(ms.messages))
	for _, v := range ms.messages {
		addresses = append(addresses, v.Address.String())
	}
	return fmt.Sprintf("[%v]", strings.Join(addresses, ", "))
}
