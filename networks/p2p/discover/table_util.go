// Modifications Copyright 2026 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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

package discover

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
	"time"
)

// TODO: salvage pushNode and deleteNode to this file.

func hasNode(list []*Node, n *Node) bool {
	for _, e := range list {
		if e.ID == n.ID {
			return true
		}
	}
	return false
}

func bumpNode(list []*Node, n *Node) bool {
	for i := range list {
		if list[i].ID == n.ID {
			// move it to the front
			copy(list[1:], list[:i])
			list[0] = n
			return true
		}
	}
	return false
}

// A shared pseudorandom source that can be shared by multiple storages.
// This RNG is shared so that the main table loop can seed it periodically.
type sharedRand struct {
	mu   sync.Mutex
	rand *rand.Rand
}

func newSharedRand() *sharedRand {
	return &sharedRand{
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *sharedRand) Seed() {
	s.mu.Lock()
	defer s.mu.Unlock()
	var b [8]byte
	if _, err := crand.Read(b[:]); err != nil {
		// If crypto/rand seed fails, fall back to time seed
		s.rand.Seed(time.Now().UnixNano())
		return
	}
	s.rand.Seed(int64(binary.BigEndian.Uint64(b[:])))
}

func (s *sharedRand) Intn(n int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rand.Intn(n)
}

func (s *sharedRand) Int63n(n int64) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rand.Int63n(n)
}

func (s *sharedRand) Perm(n int) []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rand.Perm(n)
}

func (s *sharedRand) Shuffle(n int, swap func(i, j int)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rand.Shuffle(n, swap)
}

func (s *sharedRand) Read(p []byte) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rand.Read(p)
}
