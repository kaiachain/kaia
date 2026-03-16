// Modifications Copyright 2024 The Kaia Authors
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
//
// This file is derived from p2p/discover/table_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package discover

import (
	"fmt"
	"math/rand"
	"net"
	"reflect"
	"sync"
	"testing"
	"testing/quick"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
)

func TestTable_pingReplace(t *testing.T) {
	run := func(newNodeResponding, lastInBucketResponding bool) {
		name := fmt.Sprintf("newNodeResponding=%t/lastInBucketResponding=%t", newNodeResponding, lastInBucketResponding)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testPingReplace(t, newNodeResponding, lastInBucketResponding)
		})
	}

	run(true, true)
	run(false, true)
	run(true, false)
	run(false, false)
}

func testPingReplace(t *testing.T, newNodeIsResponding, lastInBucketIsResponding bool) {
	transport := newPingRecorder()
	conf := Config{
		udp:        transport,
		Id:         NodeID{},
		Addr:       &net.UDPAddr{},
		Bootnodes:  nil,
		NodeDBPath: "",
	}
	discv, _ := newTable(&conf)
	tab := discv.(*Table)
	tab.addStorage(NodeTypeUnknown, &KademliaStorage{targetType: NodeTypeUnknown})
	go tab.loop()
	defer tab.Close()

	// Wait for init so bond is accepted.
	<-tab.initDone

	// fill up the sender's bucket.
	pingSender := NewNode(MustHexID("a502af0f59b2aab7746995408c79e9ca312d2793cc997e44fc55eda62f0150bbb8c59a6f9269ba3a081518b62699ee807c7c19c20125ddfccca872608af9e370"),
		net.IP{}, 99, 99, nil, NodeTypeUnknown)
	last := fillOldBucket(tab, pingSender)

	// this call to bond should replace the last node
	// in its bucket if the node is not responding.
	transport.dead[last.ID] = !lastInBucketIsResponding
	transport.dead[pingSender.ID] = !newNodeIsResponding
	tab.Bond(true, pingSender.ID, &net.UDPAddr{}, 0, NodeTypeUnknown)
	tab.doRevalidate(make(chan struct{}, 1))

	// first ping goes to sender (bonding pingback)
	if !transport.pinged[pingSender.ID] {
		t.Error("table did not ping back sender")
	}
	if !transport.pinged[last.ID] {
		// second ping goes to oldest node in bucket
		// to see whether it is still alive.
		t.Error("table did not ping last node in bucket")
	}

	wantSize := bucketSize
	if !lastInBucketIsResponding && !newNodeIsResponding {
		wantSize--
	}
	if l := len(tab.bucket(pingSender.sha, NodeTypeUnknown).entries); l != wantSize {
		t.Errorf("wrong bucket size after bond: got %d, want %d", l, wantSize)
	}
	if found := contains(tab.bucket(pingSender.sha, NodeTypeUnknown).entries, last.ID); found != lastInBucketIsResponding {
		t.Errorf("last entry found: %t, want: %t", found, lastInBucketIsResponding)
	}
	wantNewEntry := newNodeIsResponding && !lastInBucketIsResponding
	if found := contains(tab.bucket(pingSender.sha, NodeTypeUnknown).entries, pingSender.ID); found != wantNewEntry {
		t.Errorf("new entry found: %t, want: %t", found, wantNewEntry)
	}
}

func TestBucket_bumpNoDuplicates(t *testing.T) {
	t.Parallel()
	cfg := &quick.Config{
		MaxCount: 1000,
		Rand:     rand.New(rand.NewSource(time.Now().Unix())),
		Values: func(args []reflect.Value, rand *rand.Rand) {
			// generate a random list of nodes. this will be the content of the bucket.
			n := rand.Intn(bucketSize-1) + 1
			nodes := make([]*Node, n)
			for i := range nodes {
				nodes[i] = nodeAtDistance(common.Hash{}, 200, NodeTypeUnknown)
			}
			args[0] = reflect.ValueOf(nodes)
			// generate random bump positions.
			bumps := make([]int, rand.Intn(100))
			for i := range bumps {
				bumps[i] = rand.Intn(len(nodes))
			}
			args[1] = reflect.ValueOf(bumps)
		},
	}

	prop := func(nodes []*Node, bumps []int) (ok bool) {
		b := &bucket{entries: make([]*Node, len(nodes))}
		copy(b.entries, nodes)
		for i, pos := range bumps {
			b.bump(b.entries[pos])
			if hasDuplicates(b.entries) {
				t.Logf("bucket has duplicates after %d/%d bumps:", i+1, len(bumps))
				for _, n := range b.entries {
					t.Logf("  %v", n.ID)
				}
				return false
			}
		}
		return true
	}
	if err := quick.Check(prop, cfg); err != nil {
		t.Error(err)
	}
}

// This checks that the table-wide IP limit is applied correctly.
func TestTable_IPLimit(t *testing.T) {
	transport := newPingRecorder()
	conf := Config{
		udp:        transport,
		Id:         NodeID{},
		Addr:       &net.UDPAddr{},
		Bootnodes:  nil,
		NodeDBPath: "",
	}
	discv, _ := newTable(&conf)
	tab := discv.(*Table)
	tab.addStorage(NodeTypeUnknown, &KademliaStorage{targetType: NodeTypeUnknown})
	go tab.loop()
	defer tab.Close()

	for i := range tableIPLimit + 1 {
		n := nodeAtDistance(tab.self.sha, i, NodeTypeUnknown)
		n.IP = net.IP{172, 0, 1, byte(i)}
		tab.add(n)
	}
	if tab.len() > tableIPLimit {
		t.Errorf("too many nodes in table")
	}
}

// This checks that the table-wide IP limit is applied correctly.
func TestTable_BucketIPLimit(t *testing.T) {
	transport := newPingRecorder()
	conf := Config{
		udp:        transport,
		Id:         NodeID{},
		Addr:       &net.UDPAddr{},
		Bootnodes:  nil,
		NodeDBPath: "",
	}
	discv, _ := newTable(&conf)
	tab := discv.(*Table)
	tab.addStorage(NodeTypeUnknown, &KademliaStorage{targetType: NodeTypeUnknown})
	go tab.loop()
	defer tab.Close()

	d := 3
	for i := range bucketIPLimit + 1 {
		n := nodeAtDistance(tab.self.sha, d, NodeTypeUnknown)
		n.NType = NodeTypeUnknown
		n.IP = net.IP{172, 0, 1, byte(i)}
		tab.add(n)
	}
	if tab.len() > bucketIPLimit {
		t.Errorf("too many nodes in table")
	}
}

// fillBucket inserts nodes into the given bucket until
// it is full. The node's IDs dont correspond to their
// hashes.
func fillOldBucket(tab *Table, n *Node) (last *Node) {
	ld := logdist(tab.self.sha, n.sha)
	b := tab.bucket(n.sha, n.NType)
	for len(b.entries) < bucketSize {
		b.entries = append(b.entries, nodeAtDistance(tab.self.sha, ld, n.NType))
	}
	return b.entries[bucketSize-1]
}

type pingRecorder struct {
	mu           sync.Mutex
	dead, pinged map[NodeID]bool
}

func newPingRecorder() *pingRecorder {
	return &pingRecorder{
		dead:   make(map[NodeID]bool),
		pinged: make(map[NodeID]bool),
	}
}

func (t *pingRecorder) findnode(toid NodeID, toaddr *net.UDPAddr, target NodeID, nType NodeType, max int) ([]*Node, error) {
	return nil, nil
}
func (t *pingRecorder) close() {}
func (t *pingRecorder) waitping(from NodeID, fromIP net.IP) error {
	return nil // remote always pings
}

func (t *pingRecorder) ping(toid NodeID, toaddr *net.UDPAddr) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.pinged[toid] = true
	if t.dead[toid] {
		return errTimeout
	} else {
		return nil
	}
}

func TestTable_closest(t *testing.T) {
	t.Parallel()

	test := func(test *closeTest) bool {
		// for any node table, Target and N
		transport := newPingRecorder()
		conf := Config{
			udp:        transport,
			Id:         test.Self,
			Addr:       &net.UDPAddr{},
			Bootnodes:  nil,
			NodeDBPath: "",
		}
		discv, _ := newTable(&conf)
		tab := discv.(*Table)
		tab.addStorage(NodeTypeUnknown, &KademliaStorage{targetType: NodeTypeUnknown})
		go tab.loop()
		defer tab.Close()
		tab.stuff(test.All, NodeTypeUnknown)

		// check that doClosest(Target, N) returns nodes
		result := tab.closest(test.Target, NodeTypeUnknown, test.N).entries
		if hasDuplicates(result) {
			t.Errorf("result contains duplicates")
			return false
		}
		if !sortedByDistanceTo(test.Target, result) {
			t.Errorf("result is not sorted by distance to target")
			return false
		}

		// check that the number of results is min(N, tablen)
		wantN := test.N
		if tlen := tab.len(); tlen < test.N {
			wantN = tlen
		}
		if len(result) != wantN {
			t.Errorf("wrong number of nodes: got %d, want %d", len(result), wantN)
			return false
		} else if len(result) == 0 {
			return true // no need to check distance
		}

		// check that the result nodes have minimum distance to target.
		for _, n := range tab.nodes() {
			if contains(result, n.ID) {
				continue // don't run the check below for nodes in result
			}
			farthestResult := result[len(result)-1].sha
			if distcmp(test.Target, n.sha, farthestResult) < 0 {
				t.Errorf("table contains node that is closer to target but it's not in result")
				t.Logf("  Target:          %v", test.Target)
				t.Logf("  Farthest Result: %v", farthestResult)
				t.Logf("  ID:              %v", n.ID)
				return false
			}
		}
		return true
	}
	if err := quick.Check(test, quickcfg()); err != nil {
		t.Error(err)
	}
}

func TestTable_ReadRandomNodesGetAll(t *testing.T) {
	cfg := &quick.Config{
		MaxCount: 200,
		Rand:     rand.New(rand.NewSource(time.Now().Unix())),
		Values: func(args []reflect.Value, rand *rand.Rand) {
			args[0] = reflect.ValueOf(make([]*Node, rand.Intn(1000)))
		},
	}
	test := func(buf []*Node) bool {
		transport := newPingRecorder()
		conf := Config{
			udp:        transport,
			Id:         NodeID{},
			Addr:       &net.UDPAddr{},
			Bootnodes:  nil,
			NodeDBPath: "",
		}
		discv, _ := newTable(&conf)
		tab := discv.(*Table)
		tab.addStorage(NodeTypeUnknown, &KademliaStorage{targetType: NodeTypeUnknown})
		go tab.loop()
		defer tab.Close()
		<-tab.initDone

		for range buf {
			ld := cfg.Rand.Intn(len(tab.storages[NodeTypeUnknown].(*KademliaStorage).buckets))
			tab.stuff([]*Node{nodeAtDistance(tab.self.sha, ld, NodeTypeUnknown)}, NodeTypeUnknown)
		}
		gotN := tab.ReadRandomNodes(buf, NodeTypeUnknown)
		if gotN != tab.len() {
			t.Errorf("wrong number of nodes, got %d, want %d", gotN, tab.len())
			return false
		}
		if hasDuplicates(buf[:gotN]) {
			t.Errorf("result contains duplicates")
			return false
		}
		return true
	}
	if err := quick.Check(test, cfg); err != nil {
		t.Error(err)
	}
}

type closeTest struct {
	Self   NodeID
	Target common.Hash
	All    []*Node
	N      int
}

func (*closeTest) Generate(rand *rand.Rand, size int) reflect.Value {
	t := &closeTest{
		Self:   gen(NodeID{}, rand).(NodeID),
		Target: gen(common.Hash{}, rand).(common.Hash),
		N:      rand.Intn(bucketSize),
	}
	for _, id := range gen([]NodeID{}, rand).([]NodeID) {
		t.All = append(t.All, &Node{ID: id})
	}
	return reflect.ValueOf(t)
}

func TestTable_Lookup(t *testing.T) {
	self := nodeAtDistance(common.Hash{}, 0, NodeTypeUnknown)
	conf := Config{
		udp:        lookupTestnet,
		Id:         self.ID,
		Addr:       &net.UDPAddr{},
		Bootnodes:  nil,
		NodeDBPath: "",
	}
	discv, _ := newTable(&conf)
	tab := discv.(*Table)
	tab.addStorage(NodeTypeUnknown, &KademliaStorage{targetType: NodeTypeUnknown})
	go tab.loop()
	defer tab.Close()

	// lookup on empty table returns no nodes
	if results := tab.Lookup(lookupTestnet.target, NodeTypeUnknown); len(results) > 0 {
		t.Fatalf("lookup on empty table returned %d results: %#v", len(results), results)
	}
	// seed table with initial node (otherwise lookup will terminate immediately)
	seed := NewNode(lookupTestnet.dists[256][0], net.IP{}, 256, 0, nil, NodeTypeUnknown)
	tab.stuff([]*Node{seed}, NodeTypeUnknown)

	results := tab.Lookup(lookupTestnet.target, NodeTypeUnknown)
	t.Logf("results:")
	for _, e := range results {
		t.Logf("  ld=%d, %x", logdist(lookupTestnet.targetSha, e.sha), e.sha[:])
	}
	if len(results) != bucketSize {
		t.Errorf("wrong number of results: got %d, want %d", len(results), bucketSize)
	}
	if hasDuplicates(results) {
		t.Errorf("result set contains duplicate entries")
	}
	if !sortedByDistanceTo(lookupTestnet.targetSha, results) {
		t.Errorf("result set not sorted by distance to target")
	}
	// TODO: check result nodes are actually closest
}

// mine generates a testnet struct literal with nodes at
// various distances to the given target.
func (tn *preminedTestnet) mine(target NodeID) {
	tn.target = target
	tn.targetSha = crypto.Keccak256Hash(tn.target[:])
	found := 0
	for found < bucketSize*10 {
		k := newkey()
		id := PubkeyID(&k.PublicKey)
		sha := crypto.Keccak256Hash(id[:])
		ld := logdist(tn.targetSha, sha)
		if len(tn.dists[ld]) < bucketSize {
			tn.dists[ld] = append(tn.dists[ld], id)
			fmt.Println("found ID with ld", ld)
			found++
		}
	}
	fmt.Println("&preminedTestnet{")
	fmt.Printf("	target: %#v,\n", tn.target)
	fmt.Printf("	targetSha: %#v,\n", tn.targetSha)
	fmt.Printf("	dists: [%d][]NodeID{\n", len(tn.dists))
	for ld, ns := range tn.dists {
		if len(ns) == 0 {
			continue
		}
		fmt.Printf("		%d: []NodeID{\n", ld)
		for _, n := range ns {
			fmt.Printf("			MustHexID(\"%x\"),\n", n[:])
		}
		fmt.Println("		},")
	}
	fmt.Println("	},")
	fmt.Println("}")
}

func contains(ns []*Node, id NodeID) bool {
	for _, n := range ns {
		if n.ID == id {
			return true
		}
	}
	return false
}
