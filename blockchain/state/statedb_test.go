// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from core/state/statedb_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package state

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strings"
	"sync"
	"testing"
	"testing/quick"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"
)

// Updating a state statedb without commit must not affect persistent DB.
func TestUpdateLeaks(t *testing.T) {
	// Create an empty state database
	memDBManager := database.NewMemoryDBManager()
	db := memDBManager.GetMemDB()
	state, _ := New(common.Hash{}, NewDatabase(memDBManager), nil, nil)

	// Update it with some accounts
	for i := byte(0); i < 255; i++ {
		addr := common.BytesToAddress([]byte{i})
		if i%2 == 0 {
			state.SetState(addr, common.BytesToHash([]byte{i, i, i}), common.BytesToHash([]byte{i, i, i, i}))
		}
		if i%3 == 0 {
			state.SetCode(addr, []byte{i, i, i, i, i})
		}
		state.AddBalance(addr, big.NewInt(int64(11*i)))
		state.SetNonce(addr, uint64(42*i))
		state.IntermediateRoot(false)
	}
	// Ensure that no data was leaked into the database.
	// DB should be empty.
	for _, key := range db.Keys() {
		value, _ := db.Get(key)
		t.Errorf("State leaked into database: %x -> %x", key, value)
	}
}

// Tests that no intermediate state of an object is stored into the database,
// only the one right before the commit.
func TestIntermediateLeaks(t *testing.T) {
	// Create two state databases, one transitioning to the final state, the other final from the beginning
	transDBManager := database.NewMemoryDBManager()
	finalDBManager := database.NewMemoryDBManager()

	transDb := transDBManager.GetMemDB()
	finalDb := finalDBManager.GetMemDB()

	transState, _ := New(common.Hash{}, NewDatabase(transDBManager), nil, nil)
	finalState, _ := New(common.Hash{}, NewDatabase(finalDBManager), nil, nil)

	modify := func(state *StateDB, addr common.Address, i, tweak byte) {
		if i%2 == 0 {
			state.SetState(addr, common.Hash{i, i, i, 0}, common.Hash{})
			state.SetState(addr, common.Hash{i, i, i, tweak}, common.Hash{i, i, i, i, tweak})
		}
		if i%3 == 0 {
			state.SetCode(addr, []byte{i, i, i, i, i, tweak})
		}
		state.SetBalance(addr, big.NewInt(int64(11*i)+int64(tweak)))
		state.SetNonce(addr, uint64(42*i+tweak))
	}

	// Modify the transient state.
	for i := byte(0); i < 255; i++ {
		modify(transState, common.Address{byte(i)}, i, 0)
	}
	// Write modifications to trie.
	transState.IntermediateRoot(false)

	// Overwrite all the data with new values in the transient database.
	for i := byte(0); i < 255; i++ {
		modify(transState, common.Address{byte(i)}, i, 99)
		modify(finalState, common.Address{byte(i)}, i, 99)
	}

	// Commit and cross check the databases.
	if _, err := transState.Commit(false); err != nil {
		t.Fatalf("failed to commit transition state: %v", err)
	}
	if _, err := finalState.Commit(false); err != nil {
		t.Fatalf("failed to commit final state: %v", err)
	}
	for _, key := range finalDb.Keys() {
		if _, err := transDb.Get(key); err != nil {
			val, _ := finalDb.Get(key)
			t.Errorf("entry missing from the transition database: %x -> %x", key, val)
		}
	}
	for _, key := range transDb.Keys() {
		if _, err := finalDb.Get(key); err != nil {
			val, _ := transDb.Get(key)
			t.Errorf("extra entry in the transition database: %x -> %x", key, val)
		}
	}
}

// TestCopy tests that copying a statedb object indeed makes the original and
// the copy independent of each other. This test is a regression test against
// https://github.com/ethereum/go-ethereum/pull/15549.
func TestCopy(t *testing.T) {
	// Create a random state test to copy and modify "independently"
	orig, _ := New(common.Hash{}, NewDatabase(database.NewMemoryDBManager()), nil, nil)

	for i := byte(0); i < 255; i++ {
		obj := orig.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		obj.AddBalance(big.NewInt(int64(i)))
		orig.updateStateObject(obj)
	}
	orig.Finalise(false, true)

	// Copy the state, modify both in-memory
	copy := orig.Copy()

	// Copy the copy state
	ccopy := copy.Copy()

	for i := byte(0); i < 255; i++ {
		origObj := orig.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		copyObj := copy.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		ccopyObj := ccopy.GetOrNewStateObject(common.BytesToAddress([]byte{i}))

		origObj.AddBalance(big.NewInt(2 * int64(i)))
		copyObj.AddBalance(big.NewInt(3 * int64(i)))
		ccopyObj.AddBalance(big.NewInt(4 * int64(i)))

		orig.updateStateObject(origObj)
		copy.updateStateObject(copyObj)
		ccopy.updateStateObject(ccopyObj)
	}

	// Finalise the changes on all concurrently
	finalise := func(wg *sync.WaitGroup, db *StateDB) {
		defer wg.Done()
		db.Finalise(true, true)
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go finalise(&wg, orig)
	go finalise(&wg, copy)
	go finalise(&wg, ccopy)
	wg.Wait()

	// Verify that the two states have been updated independently
	for i := byte(0); i < 255; i++ {
		origObj := orig.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		copyObj := copy.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		ccopyObj := ccopy.GetOrNewStateObject(common.BytesToAddress([]byte{i}))

		if want := big.NewInt(3 * int64(i)); origObj.Balance().Cmp(want) != 0 {
			t.Errorf("orig obj %d: balance mismatch: have %v, want %v", i, origObj.Balance(), want)
		}
		if want := big.NewInt(4 * int64(i)); copyObj.Balance().Cmp(want) != 0 {
			t.Errorf("copy obj %d: balance mismatch: have %v, want %v", i, copyObj.Balance(), want)
		}
		if want := big.NewInt(5 * int64(i)); ccopyObj.Balance().Cmp(want) != 0 {
			t.Errorf("ccopy obj %d: balance mismatch: have %v, want %v", i, ccopyObj.Balance(), want)
		}
	}
}

func TestSnapshotRandom(t *testing.T) {
	config := &quick.Config{MaxCount: 1000}
	err := quick.Check((*snapshotTest).run, config)
	if cerr, ok := err.(*quick.CheckError); ok {
		test := cerr.In[0].(*snapshotTest)
		t.Errorf("%v:\n%s", test.err, test)
	} else if err != nil {
		t.Error(err)
	}
}

// TestStateObjects tests basic functional operations of StateObjects.
// It will be updated by StateDB.Commit() with state objects in StateDB.stateObjects.
func TestStateObjects(t *testing.T) {
	stateDB, _ := New(common.Hash{}, NewDatabase(database.NewMemoryDBManager()), nil, nil)

	// Update each account, it will update StateDB.stateObjects.
	for i := byte(0); i < 128; i++ {
		addr := common.BytesToAddress([]byte{i})
		stateObj := stateDB.GetOrNewStateObject(addr)

		stateObj.AddBalance(big.NewInt(int64(i)))
		stateDB.updateStateObject(stateObj)
	}

	assert.Equal(t, 128, len(stateDB.stateObjects))
}

// TestCopiedEIP7702 tests that copied EOA has the same code related fields as the original EOA.
// This test has been introduced since the implementation of EIP-7702.
func TestCopiedEIP7702(t *testing.T) {
	stateDB, _ := New(common.Hash{}, NewDatabase(database.NewMemoryDBManager()), nil, nil)

	testCode := common.Hex2Bytes("0xef0100")
	testCodeHash := crypto.Keccak256Hash(testCode)

	addr := common.BytesToAddress([]byte{5})
	stateDB.SetCodeToEOA(addr, testCode, params.Rules{})

	assert.Equal(t, stateDB.GetCodeHash(addr), testCodeHash)
	pa := account.GetProgramAccount(stateDB.GetAccount(addr))
	assert.Equal(t, pa.GetStorageRoot(), types.EmptyRootHash.ExtendZero())

	copy := stateDB.Copy()

	assert.Equal(t, copy.GetCodeHash(addr), testCodeHash)
	pa = account.GetProgramAccount(copy.GetAccount(addr))
	assert.Equal(t, pa.GetStorageRoot(), types.EmptyRootHash.ExtendZero())
}

// Test that invalid pruning options are prohibited.
func TestPruningOptions(t *testing.T) {
	opens := func(pruning bool, pruningNum bool) bool {
		dbm := database.NewMemoryDBManager()
		opts := &statedb.TrieOpts{}
		if pruning {
			dbm.WritePruningEnabled()
		}
		if pruningNum {
			opts.PruningBlockNumber = 1
		}
		_, err := New(common.Hash{}, NewDatabase(dbm), nil, opts)
		return err == nil
	}

	// DB pruning disabled & not request pruning. Normal non-pruning setup.
	assert.True(t, opens(false, false))
	// DB pruning disabled & request pruning must fail.
	assert.False(t, opens(false, true))

	// DB pruning enabled & not request pruning. Temporary trie,
	// such as in debug_traceTransaction or eth_call.
	assert.True(t, opens(true, false))
	// DB pruning enabled & request pruning. Normal pruning setup,
	// such as in InsertChain.
	assert.True(t, opens(true, true))
}

// Test that the storage root (ExtHash) has correct extensions
// under different pruning options.
func TestPruningRoot(t *testing.T) {
	addr := common.HexToAddress("0xaaaa")

	makeState := func(db Database) common.Hash {
		stateDB, _ := New(common.Hash{}, db, nil, nil)
		stateDB.CreateSmartContractAccount(addr, params.CodeFormatEVM, params.Rules{})
		stateDB.SetState(addr, common.HexToHash("1"), common.HexToHash("2"))
		root, _ := stateDB.Commit(false)
		return root
	}

	// When pruning is disabled, storage root is zero-extended.
	dbm := database.NewMemoryDBManager()

	db := NewDatabase(dbm)
	root := makeState(db)
	stateDB, _ := New(root, db, nil, nil)
	storageRoot, _ := stateDB.GetContractStorageRoot(addr)
	assert.True(t, storageRoot.IsZeroExtended())

	// When pruning is enabled, storage root is nonzero-extended.
	dbm = database.NewMemoryDBManager()
	dbm.WritePruningEnabled()

	db = NewDatabase(dbm)
	root = makeState(db)
	stateDB, _ = New(root, db, nil, nil) // Reopen trie to check the account stored in disk.
	storageRoot, _ = stateDB.GetContractStorageRoot(addr)
	assert.False(t, storageRoot.IsZeroExtended())
}

// A snapshotTest checks that reverting StateDB snapshots properly undoes all changes
// captured by the snapshot. Instances of this test with pseudorandom content are created
// by Generate.
//
// The test works as follows:
//
// A new state is created and all actions are applied to it. Several snapshots are taken
// in between actions. The test then reverts each snapshot. For each snapshot the actions
// leading up to it are replayed on a fresh, empty state. The behaviour of all public
// accessor methods on the reverted state must match the return value of the equivalent
// methods on the replayed state.
type snapshotTest struct {
	addrs     []common.Address // all account addresses
	actions   []testAction     // modifications to the state
	snapshots []int            // actions indexes at which snapshot is taken
	err       error            // failure details are reported through this field
}

type testAction struct {
	name   string
	fn     func(testAction, *StateDB)
	args   []int64
	noAddr bool
}

// newTestAction creates a random action that changes state.
func newTestAction(addr common.Address, r *rand.Rand) testAction {
	actions := []testAction{
		{
			name: "SetBalance",
			fn: func(a testAction, s *StateDB) {
				s.SetBalance(addr, big.NewInt(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "AddBalance",
			fn: func(a testAction, s *StateDB) {
				s.AddBalance(addr, big.NewInt(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "SetNonce",
			fn: func(a testAction, s *StateDB) {
				s.SetNonce(addr, uint64(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "SetState",
			fn: func(a testAction, s *StateDB) {
				var key, val common.Hash
				binary.BigEndian.PutUint16(key[:], uint16(a.args[0]))
				binary.BigEndian.PutUint16(val[:], uint16(a.args[1]))
				s.SetState(addr, key, val)
			},
			args: make([]int64, 2),
		},
		{
			name: "SetCode",
			fn: func(a testAction, s *StateDB) {
				code := make([]byte, 16)
				binary.BigEndian.PutUint64(code, uint64(a.args[0]))
				binary.BigEndian.PutUint64(code[8:], uint64(a.args[1]))
				s.SetCode(addr, code)
			},
			args: make([]int64, 2),
		},
		{
			name: "CreateAccount",
			fn: func(a testAction, s *StateDB) {
				s.CreateSmartContractAccount(addr, params.CodeFormatEVM, params.Rules{IsIstanbul: true})
			},
		},
		{
			name: "SelfDestruct",
			fn: func(a testAction, s *StateDB) {
				s.SelfDestruct(addr)
			},
		},
		{
			name: "AddRefund",
			fn: func(a testAction, s *StateDB) {
				s.AddRefund(uint64(a.args[0]))
			},
			args:   make([]int64, 1),
			noAddr: true,
		},
		{
			name: "AddLog",
			fn: func(a testAction, s *StateDB) {
				data := make([]byte, 2)
				binary.BigEndian.PutUint16(data, uint16(a.args[0]))
				s.AddLog(&types.Log{Address: addr, Data: data})
			},
			args: make([]int64, 1),
		},
		{
			name: "AddPreimage",
			fn: func(a testAction, s *StateDB) {
				preimage := []byte{1}
				hash := common.BytesToHash(preimage)
				s.AddPreimage(hash, preimage)
			},
			args: make([]int64, 1),
		},
		{
			name: "AddAddressToAccessList",
			fn: func(a testAction, s *StateDB) {
				s.AddAddressToAccessList(addr)
			},
		},
		{
			name: "AddSlotToAccessList",
			fn: func(a testAction, s *StateDB) {
				s.AddSlotToAccessList(addr,
					common.Hash{byte(a.args[0])})
			},
			args: make([]int64, 1),
		},
		{
			name: "SetTransientState",
			fn: func(a testAction, s *StateDB) {
				var key, val common.Hash
				binary.BigEndian.PutUint16(key[:], uint16(a.args[0]))
				binary.BigEndian.PutUint16(val[:], uint16(a.args[1]))
				s.SetTransientState(addr, key, val)
			},
			args: make([]int64, 2),
		},
	}
	action := actions[r.Intn(len(actions))]
	var nameargs []string
	if !action.noAddr {
		nameargs = append(nameargs, addr.Hex())
	}
	for _, i := range action.args {
		action.args[i] = rand.Int63n(100)
		nameargs = append(nameargs, fmt.Sprint(action.args[i]))
	}
	action.name += strings.Join(nameargs, ", ")
	return action
}

// Generate returns a new snapshot test of the given size. All randomness is
// derived from r.
func (*snapshotTest) Generate(r *rand.Rand, size int) reflect.Value {
	// Generate random actions.
	addrs := make([]common.Address, 50)
	for i := range addrs {
		addrs[i][0] = byte(i)
	}
	actions := make([]testAction, size)
	for i := range actions {
		addr := addrs[r.Intn(len(addrs))]
		actions[i] = newTestAction(addr, r)
	}
	// Generate snapshot indexes.
	nsnapshots := int(math.Sqrt(float64(size)))
	if size > 0 && nsnapshots == 0 {
		nsnapshots = 1
	}
	snapshots := make([]int, nsnapshots)
	snaplen := len(actions) / nsnapshots
	for i := range snapshots {
		// Try to place the snapshots some number of actions apart from each other.
		snapshots[i] = (i * snaplen) + r.Intn(snaplen)
	}
	return reflect.ValueOf(&snapshotTest{addrs, actions, snapshots, nil})
}

func (test *snapshotTest) String() string {
	out := new(bytes.Buffer)
	sindex := 0
	for i, action := range test.actions {
		if len(test.snapshots) > sindex && i == test.snapshots[sindex] {
			fmt.Fprintf(out, "---- snapshot %d ----\n", sindex)
			sindex++
		}
		fmt.Fprintf(out, "%4d: %s\n", i, action.name)
	}
	return out.String()
}

func (test *snapshotTest) run() bool {
	// Run all actions and create snapshots.
	var (
		state, _     = New(common.Hash{}, NewDatabase(database.NewMemoryDBManager()), nil, nil)
		snapshotRevs = make([]int, len(test.snapshots))
		sindex       = 0
	)
	for i, action := range test.actions {
		if len(test.snapshots) > sindex && i == test.snapshots[sindex] {
			snapshotRevs[sindex] = state.Snapshot()
			sindex++
		}
		action.fn(action, state)
	}
	// Revert all snapshots in reverse order. Each revert must yield a state
	// that is equivalent to fresh state with all actions up the snapshot applied.
	for sindex--; sindex >= 0; sindex-- {
		checkstate, _ := New(common.Hash{}, state.Database(), nil, nil)
		for _, action := range test.actions[:test.snapshots[sindex]] {
			action.fn(action, checkstate)
		}
		state.RevertToSnapshot(snapshotRevs[sindex])
		if err := test.checkEqual(state, checkstate); err != nil {
			test.err = fmt.Errorf("state mismatch after revert to snapshot %d\n%v", sindex, err)
			return false
		}
	}
	return true
}

// checkEqual checks that methods of state and checkstate return the same values.
func (test *snapshotTest) checkEqual(state, checkstate *StateDB) error {
	for _, addr := range test.addrs {
		var err error
		checkeq := func(op string, a, b interface{}) bool {
			if err == nil && !reflect.DeepEqual(a, b) {
				err = fmt.Errorf("got %s(%s) == %v, want %v", op, addr.Hex(), a, b)
				return false
			}
			return true
		}
		// Check basic accessor methods.
		checkeq("Exist", state.Exist(addr), checkstate.Exist(addr))
		checkeq("HasSelfDestructed", state.HasSelfDestructed(addr), checkstate.HasSelfDestructed(addr))
		checkeq("GetBalance", state.GetBalance(addr), checkstate.GetBalance(addr))
		checkeq("GetNonce", state.GetNonce(addr), checkstate.GetNonce(addr))
		checkeq("GetCode", state.GetCode(addr), checkstate.GetCode(addr))
		checkeq("GetCodeHash", state.GetCodeHash(addr), checkstate.GetCodeHash(addr))
		checkeq("GetCodeSize", state.GetCodeSize(addr), checkstate.GetCodeSize(addr))
		// Check storage.
		if obj := state.getStateObject(addr); obj != nil {
			state.ForEachStorage(addr, func(key, value common.Hash) bool {
				return checkeq("GetState("+key.Hex()+")", checkstate.GetState(addr, key), value)
			})
			checkstate.ForEachStorage(addr, func(key, value common.Hash) bool {
				return checkeq("GetState("+key.Hex()+")", checkstate.GetState(addr, key), value)
			})
		}
		if err != nil {
			return err
		}
	}

	if state.GetRefund() != checkstate.GetRefund() {
		return fmt.Errorf("got GetRefund() == %d, want GetRefund() == %d",
			state.GetRefund(), checkstate.GetRefund())
	}
	if !reflect.DeepEqual(state.GetLogs(common.Hash{}), checkstate.GetLogs(common.Hash{})) {
		return fmt.Errorf("got GetLogs(common.Hash{}) == %v, want GetLogs(common.Hash{}) == %v",
			state.GetLogs(common.Hash{}), checkstate.GetLogs(common.Hash{}))
	}
	return nil
}

// This test is to check the functionality of restoring with dirties of journal.
// Snapshot must remember the exact number of dirties.
func (s *StateSuite) TestSnapshotWithJournalDirties(c *check.C) {
	s.state.GetOrNewStateObject(common.Address{})
	root, _ := s.state.Commit(false)
	s.state.Reset(root)

	snapshot := s.state.Snapshot()
	s.state.AddBalance(common.Address{}, new(big.Int))

	if len(s.state.journal.dirties) != 1 {
		c.Fatal("expected one dirty state object")
	}
	s.state.RevertToSnapshot(snapshot)
	if len(s.state.journal.dirties) != 0 {
		c.Fatal("expected no dirty state object")
	}
}

// TestCopyOfCopy tests that modified objects are carried over to the copy, and the copy of the copy.
// See https://github.com/ethereum/go-ethereum/pull/15225#issuecomment-380191512
func TestCopyOfCopy(t *testing.T) {
	sdb, _ := New(common.Hash{}, NewDatabase(database.NewMemoryDBManager()), nil, nil)
	addr := common.HexToAddress("aaaa")
	sdb.SetBalance(addr, big.NewInt(42))

	if got := sdb.Copy().GetBalance(addr).Uint64(); got != 42 {
		t.Fatalf("1st copy fail, expected 42, got %v", got)
	}
	if got := sdb.Copy().Copy().GetBalance(addr).Uint64(); got != 42 {
		t.Fatalf("2nd copy fail, expected 42, got %v", got)
	}
}

// TestZeroHashNode checks returning values of `(db *Database) Node` function.
// The function should return (nil, ErrZeroHashNode) for default common.Hash{} value.
func TestZeroHashNode(t *testing.T) {
	zeroHash := common.ExtHash{}

	db := database.NewMemoryDBManager()
	sdb := NewDatabase(db)
	node, err := sdb.TrieDB().Node(zeroHash)
	if err != nil {
		assert.Equal(t, statedb.ErrZeroHashNode, err)
	}
	if node != nil {
		t.Fatalf("node should return nil value for zero hash")
	}
}

// TestMissingTrieNodes tests that if the statedb fails to load parts of the trie,
// the Commit operation fails with an error
// If we are missing trie nodes, we should not continue writing to the trie
func TestMissingTrieNodes(t *testing.T) {
	// Create an initial state with a few accounts
	memDb := database.NewMemoryDBManager()
	db := NewDatabase(memDb)
	var root common.Hash
	state, _ := New(common.Hash{}, db, nil, nil)
	addr := toAddr([]byte("so"))
	{
		state.SetBalance(addr, big.NewInt(1))
		state.SetCode(addr, []byte{1, 2, 3})
		a2 := toAddr([]byte("another"))
		state.SetBalance(a2, big.NewInt(100))
		state.SetCode(a2, []byte{1, 2, 4})
		root, _ = state.Commit(false)
		t.Logf("root: %x", root)
		// force-flush
		state.Database().TrieDB().Cap(0)
	}
	// Create a new state on the old root
	state, _ = New(root, db, nil, nil)
	// Now we clear out the memdb
	it := memDb.GetMemDB().NewIterator(nil, nil)
	for it.Next() {
		k := it.Key()
		// Leave the root intact
		if !bytes.Equal(k, root[:]) {
			t.Logf("key: %x", k)
			memDb.GetMemDB().Delete(k)
		}
	}
	balance := state.GetBalance(addr)
	// The removed elem should lead to it returning zero balance
	if exp, got := uint64(0), balance.Uint64(); got != exp {
		t.Errorf("expected %d, got %d", exp, got)
	}
	// Modify the state
	state.SetBalance(addr, big.NewInt(2))
	root, err := state.Commit(false)
	if err == nil {
		t.Fatalf("expected error, got root :%x", root)
	}
}

func TestStateDBAccessList(t *testing.T) {
	// Some helpers
	addr := func(a string) common.Address {
		return common.HexToAddress(a)
	}
	slot := func(a string) common.Hash {
		return common.HexToHash(a)
	}

	memDb := database.NewMemoryDBManager()
	db := NewDatabase(memDb)
	state, _ := New(common.Hash{}, db, nil, nil)
	state.accessList = newAccessList()

	verifyAddrs := func(astrings ...string) {
		t.Helper()
		// convert to common.Address form
		var addresses []common.Address
		addressMap := make(map[common.Address]struct{})
		for _, astring := range astrings {
			address := addr(astring)
			addresses = append(addresses, address)
			addressMap[address] = struct{}{}
		}
		// Check that the given addresses are in the access list
		for _, address := range addresses {
			if !state.AddressInAccessList(address) {
				t.Fatalf("expected %x to be in access list", address)
			}
		}
		// Check that only the expected addresses are present in the acesslist
		for address := range state.accessList.addresses {
			if _, exist := addressMap[address]; !exist {
				t.Fatalf("extra address %x in access list", address)
			}
		}
	}
	verifySlots := func(addrString string, slotStrings ...string) {
		if !state.AddressInAccessList(addr(addrString)) {
			t.Fatalf("scope missing address/slots %v", addrString)
		}
		address := addr(addrString)
		// convert to common.Hash form
		var slots []common.Hash
		slotMap := make(map[common.Hash]struct{})
		for _, slotString := range slotStrings {
			s := slot(slotString)
			slots = append(slots, s)
			slotMap[s] = struct{}{}
		}
		// Check that the expected items are in the access list
		for i, s := range slots {
			if _, slotPresent := state.SlotInAccessList(address, s); !slotPresent {
				t.Fatalf("input %d: scope missing slot %v (address %v)", i, s, addrString)
			}
		}
		// Check that no extra elements are in the access list
		index := state.accessList.addresses[address]
		if index >= 0 {
			stateSlots := state.accessList.slots[index]
			for s := range stateSlots {
				if _, slotPresent := slotMap[s]; !slotPresent {
					t.Fatalf("scope has extra slot %v (address %v)", s, addrString)
				}
			}
		}
	}

	state.AddAddressToAccessList(addr("aa"))          // 1
	state.AddSlotToAccessList(addr("bb"), slot("01")) // 2,3
	state.AddSlotToAccessList(addr("bb"), slot("02")) // 4
	verifyAddrs("aa", "bb")
	verifySlots("bb", "01", "02")

	// Make a copy
	stateCopy1 := state.Copy()
	if exp, got := 4, state.journal.length(); exp != got {
		t.Fatalf("journal length mismatch: have %d, want %d", got, exp)
	}

	// same again, should cause no journal entries
	state.AddSlotToAccessList(addr("bb"), slot("01"))
	state.AddSlotToAccessList(addr("bb"), slot("02"))
	state.AddAddressToAccessList(addr("aa"))
	if exp, got := 4, state.journal.length(); exp != got {
		t.Fatalf("journal length mismatch: have %d, want %d", got, exp)
	}
	// some new ones
	state.AddSlotToAccessList(addr("bb"), slot("03")) // 5
	state.AddSlotToAccessList(addr("aa"), slot("01")) // 6
	state.AddSlotToAccessList(addr("cc"), slot("01")) // 7,8
	state.AddAddressToAccessList(addr("cc"))
	if exp, got := 8, state.journal.length(); exp != got {
		t.Fatalf("journal length mismatch: have %d, want %d", got, exp)
	}

	verifyAddrs("aa", "bb", "cc")
	verifySlots("aa", "01")
	verifySlots("bb", "01", "02", "03")
	verifySlots("cc", "01")

	// now start rolling back changes
	state.journal.revert(state, 7)
	if _, ok := state.SlotInAccessList(addr("cc"), slot("01")); ok {
		t.Fatalf("slot present, expected missing")
	}
	verifyAddrs("aa", "bb", "cc")
	verifySlots("aa", "01")
	verifySlots("bb", "01", "02", "03")

	state.journal.revert(state, 6)
	if state.AddressInAccessList(addr("cc")) {
		t.Fatalf("addr present, expected missing")
	}
	verifyAddrs("aa", "bb")
	verifySlots("aa", "01")
	verifySlots("bb", "01", "02", "03")

	state.journal.revert(state, 5)
	if _, ok := state.SlotInAccessList(addr("aa"), slot("01")); ok {
		t.Fatalf("slot present, expected missing")
	}
	verifyAddrs("aa", "bb")
	verifySlots("bb", "01", "02", "03")

	state.journal.revert(state, 4)
	if _, ok := state.SlotInAccessList(addr("bb"), slot("03")); ok {
		t.Fatalf("slot present, expected missing")
	}
	verifyAddrs("aa", "bb")
	verifySlots("bb", "01", "02")

	state.journal.revert(state, 3)
	if _, ok := state.SlotInAccessList(addr("bb"), slot("02")); ok {
		t.Fatalf("slot present, expected missing")
	}
	verifyAddrs("aa", "bb")
	verifySlots("bb", "01")

	state.journal.revert(state, 2)
	if _, ok := state.SlotInAccessList(addr("bb"), slot("01")); ok {
		t.Fatalf("slot present, expected missing")
	}
	verifyAddrs("aa", "bb")

	state.journal.revert(state, 1)
	if state.AddressInAccessList(addr("bb")) {
		t.Fatalf("addr present, expected missing")
	}
	verifyAddrs("aa")

	state.journal.revert(state, 0)
	if state.AddressInAccessList(addr("aa")) {
		t.Fatalf("addr present, expected missing")
	}
	if got, exp := len(state.accessList.addresses), 0; got != exp {
		t.Fatalf("expected empty, got %d", got)
	}
	if got, exp := len(state.accessList.slots), 0; got != exp {
		t.Fatalf("expected empty, got %d", got)
	}
	// Check the copy
	// Make a copy
	state = stateCopy1
	verifyAddrs("aa", "bb")
	verifySlots("bb", "01", "02")
	if got, exp := len(state.accessList.addresses), 2; got != exp {
		t.Fatalf("expected empty, got %d", got)
	}
	if got, exp := len(state.accessList.slots), 1; got != exp {
		t.Fatalf("expected empty, got %d", got)
	}
}

func TestStateDBTransientStorage(t *testing.T) {
	memDb := database.NewMemoryDBManager()
	db := NewDatabase(memDb)
	state, _ := New(common.Hash{}, db, nil, nil)

	key := common.Hash{0x01}
	value := common.Hash{0x02}
	addr := common.Address{}

	state.SetTransientState(addr, key, value)
	if exp, got := 1, state.journal.length(); exp != got {
		t.Fatalf("journal length mismatch: have %d, want %d", got, exp)
	}
	// the retrieved value should equal what was set
	if got := state.GetTransientState(addr, key); got != value {
		t.Fatalf("transient storage mismatch: have %x, want %x", got, value)
	}

	// revert the transient state being set and then check that the
	// value is now the empty hash
	state.journal.revert(state, 0)
	if got, exp := state.GetTransientState(addr, key), (common.Hash{}); exp != got {
		t.Fatalf("transient storage mismatch: have %x, want %x", got, exp)
	}

	// set transient state and then copy the statedb and ensure that
	// the transient state is copied
	state.SetTransientState(addr, key, value)
	cpy := state.Copy()
	if got := cpy.GetTransientState(addr, key); got != value {
		t.Fatalf("transient storage mismatch: have %x, want %x", got, value)
	}
}
