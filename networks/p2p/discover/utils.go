// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package discover

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
)

// #region Constants

const (
	alpha           = 3  // Kademlia concurrency factor
	bucketSize      = 16 // Kademlia bucket size
	maxReplacements = 10 // Size of per-bucket replacement list

	maxBondingPingPongs = 16 // Limit on the number of concurrent ping/pong interactions
	maxFindnodeFailures = 5  // Nodes exceeding this limit are dropped

	refreshInterval    = 30 * time.Minute
	revalidateInterval = 10 * time.Second

	seedCount  = 30
	seedMaxAge = 5 * 24 * time.Hour
)

const (
	tableIPLimit, tableSubnet   = 10, 24
	bucketIPLimit, bucketSubnet = 2, 24 // at most 2 addresses from the same /24
	// We keep buckets for the upper 1/15 of distances because
	// it's very unlikely we'll ever encounter a node that's closer.
	hashBits          = len(common.Hash{}) * 8
	nBuckets          = hashBits / 15       // Number of buckets
	bucketMinDistance = hashBits - nBuckets // Log distance of closest bucket
)

// #region Types

// bucket contains nodes, ordered by their last activity. the entry
// that was most recently active is the first element in entries.
type bucket struct {
	entries      []*Node // live entries, sorted by time of last contact
	replacements []*Node // recently seen nodes to be used if revalidation fails
	ips          netutil.DistinctNetSet
}

// bump moves the given node to the front of the bucket entry list
// if it is contained in that list.
func (b *bucket) bump(n *Node) bool {
	for i := range b.entries {
		if b.entries[i].ID == n.ID {
			// move it to the front
			copy(b.entries[1:], b.entries[:i])
			b.entries[0] = n
			return true
		}
	}
	return false
}

// nodesByDistance is a list of nodes, ordered by
// distance to target.
type nodesByDistance struct {
	entries []*Node
	target  common.Hash
}

// push adds the given node to the list, keeping the total size below maxElems.
func (h *nodesByDistance) push(n *Node, maxElems int) {
	ix := sort.Search(len(h.entries), func(i int) bool {
		return distcmp(h.target, h.entries[i].sha, n.sha) > 0
	})
	if len(h.entries) < maxElems {
		h.entries = append(h.entries, n)
	}
	if ix == len(h.entries) {
		// farther away than all nodes we already have.
		// if there was room for it, the node is now the last element.
	} else {
		// slide existing entries down to make room
		// this will overwrite the entry we just appended.
		copy(h.entries[ix+1:], h.entries[ix:])
		h.entries[ix] = n
	}
}

func (h *nodesByDistance) String() string {
	return fmt.Sprintf("nodeByDistance target: %s, entries: %s", h.target, h.entries)
}

// #region Node list helpers

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

// pushNode adds n to the front of list, keeping at most max items.
func pushNode(list []*Node, n *Node, max int) ([]*Node, *Node) {
	if len(list) < max {
		list = append(list, nil)
	}
	removed := list[len(list)-1]
	copy(list[1:], list)
	list[0] = n
	return list, removed
}

// deleteNode removes n from list.
func deleteNode(list []*Node, n *Node) []*Node {
	for i := range list {
		if list[i].ID == n.ID {
			return append(list[:i], list[i+1:]...)
		}
	}
	return list
}

func removeBn(nodes []*Node) []*Node {
	tmp := nodes[:0]
	for _, n := range nodes {
		if n.NType != NodeTypeBN {
			tmp = append(tmp, n)
		}
	}
	return tmp
}

// #region sharedRand

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

// #region NodeType helpers

// nodeTypeName converts NodeType to string for logging.
func nodeTypeName(nt NodeType) string { // TODO-Kaia-Node Consolidate p2p.NodeType and common.ConnType
	switch nt {
	case NodeTypeCN:
		return "CN"
	case NodeTypePN:
		return "PN"
	case NodeTypeEN:
		return "EN"
	case NodeTypeBN:
		return "BN"
	default:
		return "Unknown Node Type"
	}
}

// EffectiveNodeType returns the KIP-311 role bucket for a wire-level node type.
func EffectiveNodeType(nt NodeType) NodeType {
	if nt == NodeTypePN {
		return NodeTypeEN
	}
	return nt
}

func ParseNodeType(nt string) NodeType {
	switch nt {
	case "cn":
		return NodeTypeCN
	case "pn":
		return NodeTypePN
	case "en":
		return NodeTypeEN
	case "bn":
		return NodeTypeBN
	default:
		return NodeTypeUnknown
	}
}

// StringNodeType converts NodeType to string
func StringNodeType(nType NodeType) string { // TODO-Kaia-Node Consolidate p2p.NodeType and common.ConnType
	switch nType {
	case NodeTypeCN:
		return "cn"
	case NodeTypePN:
		return "pn"
	case NodeTypeEN:
		return "en"
	case NodeTypeBN:
		return "bn"
	default:
		return "unknown"
	}
}
