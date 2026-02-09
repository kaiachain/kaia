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
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testSelf = NewNode(NodeID{}, net.IP{127, 0, 0, 1}, 30303, 30303, nil, NodeTypeUnknown)
)

func newTestKademliaStorage() *kademliaStorage {
	return newKademliaStorage(testSelf, newSharedRand())
}

// Create a unique node at the given distance.
func nodeAtDist(ld int) *Node {
	n := new(Node)
	n.sha = hashAtDistance(testSelf.sha, ld)
	n.IP = net.IP{172, byte(ld), 1, 0}
	n.NType = NodeTypeUnknown
	copy(n.ID[:], n.sha[:]) // ensure the node has a unique random ID
	return n
}

// Note that this function bypasses the IP limit, makes it easier to test.
func fillBucket2(b *bucket, ld int) {
	for len(b.entries) < bucketSize {
		b.entries = append(b.entries, nodeAtDist(ld))
	}
}

func TestKademliaStorage_AddBump(t *testing.T) {
	s := newTestKademliaStorage()
	n1 := nodeAtDist(10)

	// Add node first time
	s.add(n1)
	b := s.bucket(n1.sha)
	require.Equal(t, 1, len(b.entries), "expected 1 entry after first add")

	// Add second node to same bucket
	n2 := nodeAtDist(10)
	s.add(n2)
	require.Equal(t, 2, len(b.entries), "expected 2 entries after second add")

	// Bump first node - it should move to front
	s.add(n1)
	assert.Equal(t, 2, len(b.entries), "expected 2 entries after bump")
	assert.Equal(t, n1.ID, b.entries[0].ID, "bumped node should be at front")
}

func TestKademliaStorage_AddBucketFull(t *testing.T) {
	s := newTestKademliaStorage()
	targetHash := hashAtDistance(s.selfSha, 10)

	// Fill the bucket.
	b := s.bucket(targetHash)
	fillBucket2(b, 10)
	require.Equal(t, bucketSize, len(b.entries), "expected bucket to be full")

	// Add one more - should go to replacements since bucket is full
	n := nodeAtDist(10)
	s.add(n)
	assert.Equal(t, bucketSize, len(b.entries), "expected bucket to be full")
	assert.Equal(t, 1, len(b.replacements), "expected at least 1 replacement when bucket is full")

	// Manually drop a node to make room for the new node.
	b.entries = b.entries[1:]

	// Add the new node again - should go to the bucket, and removed from the replacements list.
	s.add(n)
	assert.Equal(t, bucketSize, len(b.entries), "expected bucket to be full")
	assert.Equal(t, 0, len(b.replacements), "expected no replacements when bucket is full")
}

func TestKademliaStorage_AddTableIPLimit(t *testing.T) {
	s := newTestKademliaStorage()

	// Add nodes from same /24 subnet until table limit is reached
	for i := 0; i < tableIPLimit+5; i++ {
		n := nodeAtDist(i)
		n.IP = net.IP{172, 0, 1, byte(i + 1)}
		s.add(n)
	}

	assert.LessOrEqual(t, s.len(), tableIPLimit, "table IP limit exceeded")
}

func TestKademliaStorage_AddBucketIPLimit(t *testing.T) {
	s := newTestKademliaStorage()
	b := s.bucket(hashAtDistance(s.selfSha, 10))

	// Add nodes from same /24 subnet to same bucket
	for i := 0; i < bucketIPLimit+5; i++ {
		n := nodeAtDist(10)
		n.IP = net.IP{172, 0, 1, byte(i + 1)}
		s.add(n)
	}

	assert.LessOrEqual(t, len(b.entries), bucketIPLimit, "bucket IP limit exceeded")
}

func TestKademliaStorage_Delete(t *testing.T) {
	s := newTestKademliaStorage()
	n := nodeAtDist(10)
	n.IP = net.IP{172, 0, 1, 1}

	// Add node
	s.add(n)
	require.Equal(t, 1, s.len(), "expected 1 node after add")

	// Delete node
	s.delete(n)
	assert.Equal(t, 0, s.len(), "expected 0 nodes after delete")

	// Delete again - nothing to delete.
	s.delete(n)
	assert.Equal(t, 0, s.len(), "expected 0 nodes after delete")
}

func TestKademliaStorage_DeletePromoteReplacements(t *testing.T) {
	s := newTestKademliaStorage()
	targetHash := hashAtDistance(s.selfSha, 10)

	// Fill bucket.
	b := s.bucket(targetHash)
	fillBucket2(b, 10)
	require.Equal(t, bucketSize, len(b.entries), "expected bucket to be full")

	// Add one more - should go to replacements since bucket is full
	n := nodeAtDist(10)
	s.add(n)
	assert.Equal(t, bucketSize, len(b.entries), "expected bucket to be full")
	assert.Equal(t, 1, len(b.replacements), "expected at least 1 replacement when bucket is full")

	// Delete a node from the bucket - refilled from replacements.
	x := b.entries[0]
	s.delete(x)
	assert.Equal(t, 0, len(b.replacements), "popped a replacement")
	assert.Equal(t, bucketSize, len(b.entries), "entries refilled from replacements")

	// Delete another node from the bucket - nothing to refill because replacements are empty.
	y := b.entries[0]
	s.delete(y)
	assert.Equal(t, bucketSize-1, len(b.entries), "entries decreased")
}

func TestKademliaStorage_Random(t *testing.T) {
	s := newTestKademliaStorage()

	// 3 different buckets. Note: bucketMinDistance is 239.
	fillBucket2(s.bucket(hashAtDistance(s.selfSha, 235)), 235)
	fillBucket2(s.bucket(hashAtDistance(s.selfSha, 245)), 245)
	fillBucket2(s.bucket(hashAtDistance(s.selfSha, 255)), 255)
	assert.Equal(t, 3*bucketSize, s.len(), "added all nodes over 3 different distance buckets")

	buf := make([]*Node, 100)
	n := s.random(buf)
	assert.Equal(t, 3*bucketSize, n, "returned all nodes")

	buf = make([]*Node, 7)
	n = s.random(buf)
	assert.Equal(t, 7, n, "return size capped at output buf size")

	// First 3 entrires should be from the distinct buckets.
	d1 := logdist(s.selfSha, buf[0].sha)
	d2 := logdist(s.selfSha, buf[1].sha)
	d3 := logdist(s.selfSha, buf[2].sha)
	assert.NotEqual(t, d1, d2)
	assert.NotEqual(t, d1, d3)
	assert.NotEqual(t, d2, d3)
}

func TestKademliaStorage_Closest(t *testing.T) {
	s := newTestKademliaStorage()
	targetHash := hashAtDistance(s.selfSha, 10)

	// Only two nodes in the table.
	s.add(nodeAtDist(235))
	s.add(nodeAtDist(245))
	s.add(nodeAtDist(255))
	require.Equal(t, 3, s.len())

	// Request 7 nodes, should return 3 nodes, which is all we have.
	out := s.closest(targetHash, 7)
	assert.Equal(t, 3, len(out.entries))

	// Add enough nodes to the table.
	fillBucket2(s.bucket(hashAtDistance(s.selfSha, 235)), 235)
	fillBucket2(s.bucket(hashAtDistance(s.selfSha, 245)), 245)
	fillBucket2(s.bucket(hashAtDistance(s.selfSha, 255)), 255)
	require.Equal(t, 3*bucketSize, s.len())

	// Return as many nodes as requested.
	out = s.closest(targetHash, 7)
	assert.Equal(t, 7, len(out.entries))
}

func TestKademliaStorage_Oldest(t *testing.T) {
	s := newTestKademliaStorage()

	// Returns nothing when empty.
	oldest := s.oldest()
	assert.Nil(t, oldest)

	s.add(nodeAtDist(10))
	oldest = s.oldest()
	assert.NotNil(t, oldest)
}

func TestSimpleStorage2_AddDelete(t *testing.T) {
	s := newSimpleStorage2(newSharedRand())
	n1 := nodeAtDist(10)
	n2 := nodeAtDist(10)

	// Add node first time
	s.add(n1)
	require.Equal(t, 1, len(s.nodes))

	// Add second node to same bucket
	s.add(n2)
	require.Equal(t, 2, len(s.nodes))

	// Bump first node - it should move to front
	s.add(n1)
	require.Equal(t, 2, len(s.nodes))
	assert.Equal(t, n1.ID, s.nodes[0].ID)

	// Delete a node
	s.delete(n2)
	require.Equal(t, 1, len(s.nodes))

	// Delete the same node again - nothing to delete.
	s.delete(n2)
	assert.Equal(t, 1, len(s.nodes))
}

func TestSimpleStorage2_Random(t *testing.T) {
	s := newSimpleStorage2(newSharedRand())
	n1 := nodeAtDist(10)
	n2 := nodeAtDist(10)
	n3 := nodeAtDist(10)

	// Add nodes
	s.add(n1)
	s.add(n2)
	s.add(n3)
	require.Equal(t, 3, len(s.nodes))

	buf := make([]*Node, 100)
	n := s.random(buf)
	assert.Equal(t, 3, n, "returned all nodes")

	buf = make([]*Node, 2)
	n = s.random(buf)
	assert.Equal(t, 2, n, "return size capped at output buf size")
}

func TestSimpleStorage2_Closest(t *testing.T) {
	s := newSimpleStorage2(newSharedRand())
	targetHash := common.Hash{} // irrelevant for SimpleStorage
	n1 := nodeAtDist(10)
	n2 := nodeAtDist(10)
	n3 := nodeAtDist(10)

	// Add nodes
	s.add(n1)
	s.add(n2)
	s.add(n3)
	require.Equal(t, 3, len(s.nodes))

	out := s.closest(targetHash, 100)
	assert.Equal(t, 3, len(out.entries), "returned all nodes")

	out = s.closest(targetHash, 2)
	assert.Equal(t, 2, len(out.entries), "return size capped at max requested nodes")
}

func TestSimpleStorage2_Oldest(t *testing.T) {
	s := newSimpleStorage2(newSharedRand())

	// Returns nothing when empty.
	oldest := s.oldest()
	assert.Nil(t, oldest)

	s.add(nodeAtDist(10))
	oldest = s.oldest()
	assert.NotNil(t, oldest)
}
