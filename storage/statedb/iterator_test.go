// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from trie/iterator_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package statedb

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/storage/database"
)

func TestIterator(t *testing.T) {
	trie := newEmptyTrie()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"klaytn", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"dog", "puppy"},
		{"somethingveryoddindeedthis is", "myothernodedata"},
	}
	all := make(map[string]string)
	for _, val := range vals {
		all[val.k] = val.v
		trie.Update([]byte(val.k), []byte(val.v))
	}
	trie.Commit(nil)

	found := make(map[string]string)
	it := NewIterator(trie.NodeIterator(nil))
	for it.Next() {
		found[string(it.Key)] = string(it.Value)
	}

	for k, v := range all {
		if found[k] != v {
			t.Errorf("iterator value mismatch for %s: got %q want %q", k, found[k], v)
		}
	}
}

type kv struct {
	k, v []byte
	t    bool
}

func TestIteratorLargeData(t *testing.T) {
	trie := newEmptyTrie()
	vals := make(map[string]*kv)

	for i := byte(0); i < 255; i++ {
		value := &kv{common.LeftPadBytes([]byte{i}, 32), []byte{i}, false}
		value2 := &kv{common.LeftPadBytes([]byte{10, i}, 32), []byte{i}, false}
		trie.Update(value.k, value.v)
		trie.Update(value2.k, value2.v)
		vals[string(value.k)] = value
		vals[string(value2.k)] = value2
	}

	it := NewIterator(trie.NodeIterator(nil))
	for it.Next() {
		vals[string(it.Key)].t = true
	}

	var untouched []*kv
	for _, value := range vals {
		if !value.t {
			untouched = append(untouched, value)
		}
	}

	if len(untouched) > 0 {
		t.Errorf("Missed %d nodes", len(untouched))
		for _, value := range untouched {
			t.Error(value)
		}
	}
}

// Tests that the node iterator indeed walks over the entire database contents.
func TestNodeIteratorCoverage(t *testing.T) {
	// Create some arbitrary test trie to iterate
	db, trie, _ := makeTestTrie()

	// Gather all the node hashes found by the iterator
	iterated := make(map[common.Hash]struct{})
	for it := trie.NodeIterator(nil); it.Next(true); {
		if it.Hash() != (common.Hash{}) {
			iterated[it.Hash()] = struct{}{}
		}
	}

	// Cross check the hashes and the database itself
	// db.Node contains all iterated hashes
	for itHash := range iterated {
		if _, err := db.Node(itHash.ExtendZero()); err != nil {
			t.Errorf("failed to retrieve reported node %x: %v", itHash, err)
		}
	}
	// iterated hashes contains all from db.nodes
	for exthash, obj := range db.nodes {
		hash := exthash.Unextend()
		if obj == nil || common.EmptyExtHash(exthash) {
			continue // skip empty entry
		}
		if _, ok := iterated[hash]; !ok {
			t.Errorf("state entry not reported %x", hash)
		}
	}
	// iterated hashes contains all diskDB keys
	db.Cap(0) // flush to diskDB
	for _, key := range db.diskDB.GetMemDB().Keys() {
		hash := common.BytesToExtHash(key).Unextend()
		if _, ok := iterated[hash]; !ok {
			t.Errorf("state entry not reported %x", hash)
		}
	}
}

// NodeIterator yields exact same result for Trie and StroageTrie
func TestNodeIteratorStorageTrie(t *testing.T) {
	dbm := database.NewMemoryDBManager()
	dbm.WritePruningEnabled()
	triedb := NewDatabase(dbm)

	trie1, _ := NewTrie(common.Hash{}, triedb, nil)
	hashes1 := make(map[common.Hash]struct{})
	for it := trie1.NodeIterator(nil); it.Next(true); {
		hashes1[it.Hash()] = struct{}{}
	}

	trie2, _ := NewStorageTrie(common.ExtHash{}, triedb, nil)
	hashes2 := make(map[common.Hash]struct{})
	for it := trie2.NodeIterator(nil); it.Next(true); {
		hashes2[it.Hash()] = struct{}{}
	}

	for hash := range hashes1 {
		if _, ok := hashes2[hash]; !ok {
			t.Errorf("state entry not reported %x", hash)
		}
	}
	for hash := range hashes2 {
		if _, ok := hashes1[hash]; !ok {
			t.Errorf("state entry not reported %x", hash)
		}
	}
}

type kvs struct{ k, v string }

var testdata1 = []kvs{
	{"barb", "ba"},
	{"bard", "bc"},
	{"bars", "bb"},
	{"bar", "b"},
	{"fab", "z"},
	{"food", "ab"},
	{"foos", "aa"},
	{"foo", "a"},
}

var testdata2 = []kvs{
	{"aardvark", "c"},
	{"bar", "b"},
	{"barb", "bd"},
	{"bars", "be"},
	{"fab", "z"},
	{"foo", "a"},
	{"foos", "aa"},
	{"food", "ab"},
	{"jars", "d"},
}

func TestIteratorSeek(t *testing.T) {
	trie := newEmptyTrie()
	for _, val := range testdata1 {
		trie.Update([]byte(val.k), []byte(val.v))
	}

	// Seek to the middle.
	it := NewIterator(trie.NodeIterator([]byte("fab")))
	if err := checkIteratorOrder(testdata1[4:], it); err != nil {
		t.Fatal(err)
	}

	// Seek to a non-existent key.
	it = NewIterator(trie.NodeIterator([]byte("barc")))
	if err := checkIteratorOrder(testdata1[1:], it); err != nil {
		t.Fatal(err)
	}

	// Seek beyond the end.
	it = NewIterator(trie.NodeIterator([]byte("z")))
	if err := checkIteratorOrder(nil, it); err != nil {
		t.Fatal(err)
	}
}

func checkIteratorOrder(want []kvs, it *Iterator) error {
	for it.Next() {
		if len(want) == 0 {
			return fmt.Errorf("didn't expect any more values, got key %q", it.Key)
		}
		if !bytes.Equal(it.Key, []byte(want[0].k)) {
			return fmt.Errorf("wrong key: got %q, want %q", it.Key, want[0].k)
		}
		want = want[1:]
	}
	if len(want) > 0 {
		return fmt.Errorf("iterator ended early, want key %q", want[0])
	}
	return nil
}

func TestDifferenceIterator(t *testing.T) {
	triea := newEmptyTrie()
	for _, val := range testdata1 {
		triea.Update([]byte(val.k), []byte(val.v))
	}
	triea.Commit(nil)

	trieb := newEmptyTrie()
	for _, val := range testdata2 {
		trieb.Update([]byte(val.k), []byte(val.v))
	}
	trieb.Commit(nil)

	found := make(map[string]string)
	di, _ := NewDifferenceIterator(triea.NodeIterator(nil), trieb.NodeIterator(nil))
	it := NewIterator(di)
	for it.Next() {
		found[string(it.Key)] = string(it.Value)
	}

	all := []struct{ k, v string }{
		{"aardvark", "c"},
		{"barb", "bd"},
		{"bars", "be"},
		{"jars", "d"},
	}
	for _, item := range all {
		if found[item.k] != item.v {
			t.Errorf("iterator value mismatch for %s: got %v want %v", item.k, found[item.k], item.v)
		}
	}
	if len(found) != len(all) {
		t.Errorf("iterator count mismatch: got %d values, want %d", len(found), len(all))
	}
}

func TestUnionIterator(t *testing.T) {
	triea := newEmptyTrie()
	for _, val := range testdata1 {
		triea.Update([]byte(val.k), []byte(val.v))
	}
	triea.Commit(nil)

	trieb := newEmptyTrie()
	for _, val := range testdata2 {
		trieb.Update([]byte(val.k), []byte(val.v))
	}
	trieb.Commit(nil)

	di, _ := NewUnionIterator([]NodeIterator{triea.NodeIterator(nil), trieb.NodeIterator(nil)})
	it := NewIterator(di)

	all := []struct{ k, v string }{
		{"aardvark", "c"},
		{"barb", "ba"},
		{"barb", "bd"},
		{"bard", "bc"},
		{"bars", "bb"},
		{"bars", "be"},
		{"bar", "b"},
		{"fab", "z"},
		{"food", "ab"},
		{"foos", "aa"},
		{"foo", "a"},
		{"jars", "d"},
	}

	for i, kv := range all {
		if !it.Next() {
			t.Errorf("Iterator ends prematurely at element %d", i)
		}
		if kv.k != string(it.Key) {
			t.Errorf("iterator value mismatch for element %d: got key %s want %s", i, it.Key, kv.k)
		}
		if kv.v != string(it.Value) {
			t.Errorf("iterator value mismatch for element %d: got value %s want %s", i, it.Value, kv.v)
		}
	}
	if it.Next() {
		t.Errorf("Iterator returned extra values.")
	}
}

func TestIteratorNoDups(t *testing.T) {
	var tr Trie
	for _, val := range testdata1 {
		tr.Update([]byte(val.k), []byte(val.v))
	}
	checkIteratorNoDups(t, tr.NodeIterator(nil), nil)
}

// This test checks that nodeIterator.Next can be retried after inserting missing trie nodes.
func TestIteratorContinueAfterErrorDisk(t *testing.T)    { testIteratorContinueAfterError(t, false) }
func TestIteratorContinueAfterErrorMemonly(t *testing.T) { testIteratorContinueAfterError(t, true) }

func testIteratorContinueAfterError(t *testing.T, memonly bool) {
	dbm := database.NewMemoryDBManager()
	diskdb := dbm.GetMemDB()
	triedb := NewDatabase(dbm)

	tr, _ := NewTrie(common.Hash{}, triedb, nil)
	for _, val := range testdata1 {
		tr.Update([]byte(val.k), []byte(val.v))
	}
	tr.Commit(nil)
	if !memonly {
		triedb.Commit(tr.Hash(), true, 0)
	}
	wantNodeCount := checkIteratorNoDups(t, tr.NodeIterator(nil), nil)

	var (
		diskKeys [][]byte
		memKeys  []common.ExtHash
	)
	if memonly {
		memKeys = triedb.Nodes()
	} else {
		diskKeys = diskdb.Keys()
	}
	for i := 0; i < 20; i++ {
		// Create trie that will load all nodes from DB.
		tr, _ := NewTrie(tr.Hash(), triedb, nil)

		// Remove a random node from the database. It can't be the root node
		// because that one is already loaded.
		var (
			rkey     common.Hash
			nodehash common.ExtHash
			rval     []byte
			robj     *cachedNode
		)
		for {
			if memonly {
				idx := rand.Intn(len(memKeys))
				nodehash = memKeys[idx]
			} else {
				idx := rand.Intn(len(diskKeys))
				nodehash = common.BytesToExtHash(diskKeys[idx])
			}
			rkey = nodehash.Unextend()
			if rkey != tr.Hash() {
				break
			}
		}

		if memonly {
			robj = triedb.nodes[nodehash]
			delete(triedb.nodes, nodehash)
		} else {
			rval, _ = dbm.ReadTrieNode(nodehash)
			dbm.DeleteTrieNode(nodehash)
		}
		// Iterate until the error is hit.
		seen := make(map[string]bool)
		it := tr.NodeIterator(nil)
		checkIteratorNoDups(t, it, seen)
		missing, ok := it.Error().(*MissingNodeError)
		if !ok || missing.NodeHash != rkey {
			t.Fatal("didn't hit missing node, got", it.Error())
		}

		// Add the node back and continue iteration.
		if memonly {
			triedb.nodes[nodehash] = robj
		} else {
			dbm.WriteTrieNode(nodehash, rval)
		}
		checkIteratorNoDups(t, it, seen)
		if it.Error() != nil {
			t.Fatal("unexpected error", it.Error())
		}
		if len(seen) != wantNodeCount {
			t.Fatal("wrong node iteration count, got", len(seen), "want", wantNodeCount)
		}
	}
}

// Similar to the test above, this one checks that failure to create nodeIterator at a
// certain key prefix behaves correctly when Next is called. The expectation is that Next
// should retry seeking before returning true for the first time.
func TestIteratorContinueAfterSeekErrorDisk(t *testing.T) {
	testIteratorContinueAfterSeekError(t, false)
}

func TestIteratorContinueAfterSeekErrorMemonly(t *testing.T) {
	testIteratorContinueAfterSeekError(t, true)
}

func testIteratorContinueAfterSeekError(t *testing.T, memonly bool) {
	// Commit test trie to db, then remove the node containing "bars".
	dbm := database.NewMemoryDBManager()
	triedb := NewDatabase(dbm)

	ctr, _ := NewTrie(common.Hash{}, triedb, nil)
	for _, val := range testdata1 {
		ctr.Update([]byte(val.k), []byte(val.v))
	}
	root, _ := ctr.Commit(nil)
	if !memonly {
		triedb.Commit(root, true, 0)
	}
	barNodeHash := common.HexToHash("05041990364eb72fcb1127652ce40d8bab765f2bfe53225b1170d276cc101c2e")
	nodehash := barNodeHash.ExtendZero()
	var (
		barNodeBlob []byte
		barNodeObj  *cachedNode
	)
	if memonly {
		barNodeObj = triedb.nodes[nodehash]
		delete(triedb.nodes, nodehash)
	} else {
		barNodeBlob, _ = dbm.ReadTrieNode(nodehash)
		dbm.DeleteTrieNode(nodehash)
	}
	// Create a new iterator that seeks to "bars". Seeking can't proceed because
	// the node is missing.
	tr, _ := NewTrie(root, triedb, nil)
	it := tr.NodeIterator([]byte("bars"))
	missing, ok := it.Error().(*MissingNodeError)
	if !ok {
		t.Fatal("want MissingNodeError, got", it.Error())
	} else if missing.NodeHash != barNodeHash {
		t.Fatal("wrong node missing")
	}
	// Reinsert the missing node.
	if memonly {
		triedb.nodes[nodehash] = barNodeObj
	} else {
		dbm.WriteTrieNode(nodehash, barNodeBlob)
	}
	// Check that iteration produces the right set of values.
	if err := checkIteratorOrder(testdata1[2:], NewIterator(it)); err != nil {
		t.Fatal(err)
	}
}

func checkIteratorNoDups(t *testing.T, it NodeIterator, seen map[string]bool) int {
	if seen == nil {
		seen = make(map[string]bool)
	}
	for it.Next(true) {
		if seen[string(it.Path())] {
			t.Fatalf("iterator visited node path %x twice", it.Path())
		}
		seen[string(it.Path())] = true
	}
	return len(seen)
}
