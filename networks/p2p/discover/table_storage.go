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
	"net"
	"sync"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/networks/p2p/netutil"
)

var (
	_ tableStorage = (*kademliaStorage)(nil)
	_ tableStorage = (*flatStorage)(nil)
)

type tableStorage interface {
	// add adds a node to the table.
	add(n *Node)
	// delete deletes a node from the table.
	delete(n *Node)
	// random fills the given slice with random nodes from the table. Returns the number of nodes added to the buffer.
	random(buf []*Node) int
	// closest returns the n nodes in the table that are closest to the given id.
	closest(target common.Hash, max int) *nodesByDistance
	// oldest returns the oldest node, subject to revalidation.
	oldest() *Node
	// len returns the number of nodes in the table.
	len() int
	// numReplacements returns the number of replacements in the KademliaStorage. 0 for FlatStorage.
	numReplacements() int
}

type kademliaStorage struct {
	mu      sync.RWMutex
	buckets [nBuckets]*bucket
	ips     netutil.DistinctNetSet

	rand    *sharedRand
	selfID  NodeID
	selfSha common.Hash
}

func newKademliaStorage(self *Node, rand *sharedRand) *kademliaStorage {
	buckets := [nBuckets]*bucket{}
	for i := range buckets {
		buckets[i] = &bucket{
			ips: netutil.DistinctNetSet{Subnet: bucketSubnet, Limit: bucketIPLimit},
		}
	}
	return &kademliaStorage{
		buckets: buckets,
		ips:     netutil.DistinctNetSet{Subnet: tableSubnet, Limit: tableIPLimit},
		rand:    rand,
		selfID:  self.ID,
		selfSha: self.sha,
	}
}

func (s *kademliaStorage) add(n *Node) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b := s.bucket(n.sha)
	if b.bump(n) {
		// n was already in the bucket.
		return
	}
	if len(b.entries) >= bucketSize {
		// Bucket full, maybe add as replacement.
		s.addReplacement(b, n)
		return
	}
	if !s.addIP(b, n.IP) {
		// IP limit reached.
		return
	}

	// We have the space. Safe to add.
	b.entries, _ = pushNode(b.entries, n, bucketSize) // Add to front
	b.replacements = deleteNode(b.replacements, n)    // Remove from replacements (if present)
	n.addedAt = time.Now()
}

func (s *kademliaStorage) delete(n *Node) {
	s.mu.Lock()
	defer s.mu.Unlock()

	b := s.bucket(n.sha)

	// Remove the node if exists.
	if !hasNode(b.entries, n) {
		return
	}
	b.entries = deleteNode(b.entries, n)
	s.removeIP(b, n.IP)

	// Insert a random replacement instead, if available.
	r := s.popReplacement(b)
	if r != nil {
		b.entries = append(b.entries, r)
	}
}

func (s *kademliaStorage) random(buf []*Node) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Collect all non-empty buckets.
	var buckets [][]*Node
	for _, b := range s.buckets {
		if len(b.entries) > 0 {
			buckets = append(buckets, b.entries[:])
		}
	}
	if len(buckets) == 0 {
		return 0
	}
	// Shuffle the buckets.
	s.rand.Shuffle(len(buckets), func(i, j int) {
		buckets[i], buckets[j] = buckets[j], buckets[i]
	})

	// Revolve through the buckets, pick the first entry of each bucket.
	outIdx := 0
	bucketIdx := 0
	for outIdx < len(buf) {
		buf[outIdx] = buckets[bucketIdx][0] // yield the first entry of the bucket
		outIdx++

		if len(buckets[bucketIdx]) == 1 {
			// it was the last entry of the bucket. remove the bucket itself.
			buckets = append(buckets[:bucketIdx], buckets[bucketIdx+1:]...)
		} else {
			// remove the first entry of the bucket.
			buckets[bucketIdx] = buckets[bucketIdx][1:]
		}
		if len(buckets) == 0 {
			// no more buckets left.
			break
		}

		bucketIdx = (bucketIdx + 1) % len(buckets)
	}

	return outIdx
}

func (s *kademliaStorage) closest(target common.Hash, max int) *nodesByDistance {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := &nodesByDistance{target: target}
	for _, b := range &s.buckets {
		for _, n := range b.entries {
			out.push(n, max)
		}
	}
	return out
}

func (s *kademliaStorage) oldest() *Node {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// The last node in a random bucket.
	for _, bucketIdx := range s.rand.Perm(len(s.buckets)) {
		b := s.buckets[bucketIdx]
		if len(b.entries) > 0 {
			return b.entries[len(b.entries)-1]
		}
	}
	return nil
}

func (s *kademliaStorage) len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, b := range s.buckets {
		n += len(b.entries)
	}
	return n
}

func (s *kademliaStorage) numReplacements() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, b := range s.buckets {
		n += len(b.replacements)
	}
	return n
}

// Choose the bucket for the given node based on the distance from the self.
func (s *kademliaStorage) bucket(sha common.Hash) *bucket {
	d := logdist(s.selfSha, sha)
	if d <= bucketMinDistance {
		return s.buckets[0]
	}
	return s.buckets[d-bucketMinDistance-1]
}

func (s *kademliaStorage) addReplacement(b *bucket, n *Node) {
	for _, e := range b.replacements {
		if e.ID == n.ID {
			return // already in list
		}
	}
	if !s.addIP(b, n.IP) {
		return
	}
	var removed *Node
	b.replacements, removed = pushNode(b.replacements, n, maxReplacements)
	if removed != nil {
		s.removeIP(b, removed.IP)
	}
}

func (s *kademliaStorage) popReplacement(b *bucket) *Node {
	if len(b.replacements) == 0 {
		return nil
	}
	idx := s.rand.Intn(len(b.replacements))
	r := b.replacements[idx]
	b.replacements = deleteNode(b.replacements, r)
	return r
}

// addIP adds the given IP to both the table-level and bucket-level IP sets.
func (s *kademliaStorage) addIP(b *bucket, ip net.IP) bool {
	if ip == nil || ip.IsUnspecified() {
		return false
	}
	if netutil.IsLAN(ip) {
		return true
	}
	if !s.ips.Add(ip) {
		return false
	}
	if !b.ips.Add(ip) {
		s.ips.Remove(ip)
		return false
	}
	return true
}

// removeIP removes the given IP from both the table-level and bucket-level IP sets.
func (s *kademliaStorage) removeIP(b *bucket, ip net.IP) {
	if netutil.IsLAN(ip) {
		return
	}
	s.ips.Remove(ip)
	b.ips.Remove(ip)
}

type flatStorage struct {
	mu    sync.RWMutex
	nodes []*Node // A flat list of nodes. Similar to Kademlia buckets, recent nodes are at the front.
	rand  *sharedRand
}

func newFlatStorage(rand *sharedRand) *flatStorage {
	return &flatStorage{
		nodes: []*Node{},
		rand:  rand,
	}
}

func (s *flatStorage) add(n *Node) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if bumpNode(s.nodes, n) {
		return
	}
	s.nodes, _ = pushNode(s.nodes, n, math.MaxInt64)
	n.addedAt = time.Now()
}

func (s *flatStorage) delete(n *Node) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes = deleteNode(s.nodes, n)
}

func (s *flatStorage) random(buf []*Node) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := min(len(buf), len(s.nodes))
	indices := s.rand.Perm(count)
	for i := 0; i < count; i++ {
		buf[i] = s.nodes[indices[i]]
	}
	return count
}

func (s *flatStorage) closest(target common.Hash, max int) *nodesByDistance {
	// There is no concept of distance in the FlatStorage.
	// Return random nodes instead.
	buf := make([]*Node, max)
	n := s.random(buf)

	out := &nodesByDistance{target: target}
	out.entries = buf[:n]
	return out
}

func (s *flatStorage) oldest() *Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.nodes) == 0 {
		return nil
	} else {
		return s.nodes[len(s.nodes)-1]
	}
}

func (s *flatStorage) len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.nodes)
}

func (s *flatStorage) numReplacements() int {
	return 0
}
