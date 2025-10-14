// Copyright 2025 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package statedb

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"

	"github.com/erigontech/erigon-lib/commitment"
	"github.com/erigontech/erigon-lib/kaiatrie"
	"github.com/erigontech/erigon-lib/state"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FlatTrie_Import(t *testing.T) {
	t.Log(commitment.ModeDirect)
	t.Log(len(state.Schema))
	t.Log(kaiatrie.ModeRawBytes)
}

func Test_FlatTrie_Random(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	r := rand.New(rand.NewSource(42)) // for determinism
	accounts, storages := randTrie(t, r, 6, 10)

	// To use in other test files
	// fmt.Printf("accounts = [][2]string{\n")
	// for _, a := range accounts {
	// 	fmt.Printf("  {\"%s\", \"%s\"},\n", a[0], a[1])
	// }
	// fmt.Printf("}\n")
	// fmt.Printf("storages = [][3]string{\n")
	// for _, s := range storages {
	// 	fmt.Printf("  {\"%s\", \"%s\", \"%s\"},\n", s[0], s[1], s[2])
	// }
	// fmt.Printf("}\n")

	// Construct the tries in both SecureTrie and FlatTrie
	secureTrie := newEmptySecureTrie()
	fnNewSecureTrie := func(common.Address) trieInterface { return newEmptySecureTrie() }
	accountTrie1, storageTries1 := constructTrie(t, secureTrie, fnNewSecureTrie, accounts, storages)

	dm, err := kaiatrie.NewTemporaryDomainsManager(t.TempDir())
	require.NoError(t, err)
	defer dm.Close()

	flatAccountTrie, err := NewFlatAccountTrie(dm, common.Hash{}, &TrieOpts{
		BaseBlockNumber: 0,
		CommitGenesis:   true,
	})
	require.NoError(t, err)
	fnNewFlatStorageTrie := func(addr common.Address) trieInterface {
		flatStorageTrie, err := NewFlatStorageTrie(dm, addr, common.Hash{}, &TrieOpts{
			BaseBlockNumber: 0,
			CommitGenesis:   true,
			AccountTrie:     flatAccountTrie,
		})
		require.NoError(t, err)
		return flatStorageTrie
	}
	accountTrie2, storageTries2 := constructTrie(t, flatAccountTrie, fnNewFlatStorageTrie, accounts, storages)

	// Compare the trie roots
	stateRoot1, storageRoots1 := calcTrieRoots(t, accountTrie1, storageTries1)
	stateRoot2, storageRoots2 := calcTrieRoots(t, accountTrie2, storageTries2)
	assert.Equal(t, stateRoot1, stateRoot2, "stateRoot")
	for addr, root := range storageRoots1 {
		assert.Equal(t, root, storageRoots2[addr], "addr = %s storageRoot = %s", addr, root)
	}

	// Check the Get() function
	checkTrieGet(t, accountTrie1, storageTries1, accounts, storages)
	checkTrieGet(t, accountTrie2, storageTries2, accounts, storages)

	// Check the NodeIterator() function, only at the leaf nodes
	require.NoError(t, dm.CommitWrites())
	checkNodeIterator(t, accountTrie1, storageTries1, accounts, storages)
	checkNodeIterator(t, accountTrie2, storageTries2, accounts, storages)
}

// Trie checking helpers
type trieInterface interface {
	GetKey(key []byte) []byte
	TryUpdate(key, value []byte) error
	TryGet(key []byte) ([]byte, error)
	Hash() common.Hash
	Commit(onleaf LeafCallback) (common.Hash, error)
	NodeIterator(start []byte) NodeIterator
}

func constructTrie(t *testing.T, accountTrie trieInterface, fnNewStorageTrie func(common.Address) trieInterface, accounts [][2]string, storages [][3]string) (trieInterface, map[string]trieInterface) {
	for _, a := range accounts {
		k, v := hexutil.MustDecode(a[0]), hexutil.MustDecode(a[1])
		require.NoError(t, accountTrie.TryUpdate(k, v))
	}

	storageTries := make(map[string]trieInterface)
	for _, s := range storages {
		addrS, k, v := s[0], hexutil.MustDecode(s[1]), hexutil.MustDecode(s[2])
		if _, ok := storageTries[addrS]; !ok {
			storageTries[addrS] = fnNewStorageTrie(common.HexToAddress(addrS))
		}
		require.NoError(t, storageTries[addrS].TryUpdate(k, v))
	}

	return accountTrie, storageTries
}

func calcTrieRoots(t *testing.T, accountTrie trieInterface, storageTries map[string]trieInterface) (string, map[string]string) {
	storageRoots := make(map[string]string)
	for addr, trie := range storageTries {
		root, err := trie.Commit(nil)
		require.NoError(t, err)
		storageRoots[addr] = root.Hex()
		t.Logf("storageRoot[%s] = %s", addr, storageRoots[addr])
	}

	root, err := accountTrie.Commit(nil)
	require.NoError(t, err)
	stateRoot := root.Hex()
	t.Logf("stateRoot = %s", stateRoot)

	return stateRoot, storageRoots
}

func checkTrieGet(t *testing.T, accountTrie trieInterface, storageTries map[string]trieInterface, accounts [][2]string, storages [][3]string) {
	for _, a := range accounts {
		k, v := hexutil.MustDecode(a[0]), hexutil.MustDecode(a[1])
		value, err := accountTrie.TryGet(k)
		require.NoError(t, err)
		assert.Equal(t, v, value)
	}
	for _, s := range storages {
		addrS, k, v := s[0], hexutil.MustDecode(s[1]), hexutil.MustDecode(s[2])
		storageTrie := storageTries[addrS]
		require.NotNil(t, storageTrie)

		value, err := storageTrie.TryGet(k)
		require.NoError(t, err)
		assert.Equal(t, v, value)
	}
}

func checkNodeIterator(t *testing.T, accountTrie trieInterface, storageTries map[string]trieInterface, accounts [][2]string, storages [][3]string) {
	{
		// iterate over the whole account trie
		iterated := make(map[string]string)
		it := accountTrie.NodeIterator(nil)
		for it.Next(true) {
			if !it.Leaf() {
				continue
			}
			k, v := it.LeafKey(), it.LeafBlob()
			k = accountTrie.GetKey(k)
			iterated[common.BytesToAddress(k).Hex()] = hexutil.Encode(v)
		}
		// Compare with the correct answer
		assert.Equal(t, len(iterated), len(accounts)) // |iterated| = |accounts|
		for _, a := range accounts {                  // {accounts} in {iterated}
			addrS, accS := a[0], a[1]
			assert.Equal(t, accS, iterated[addrS], "addr = %s, acc = %s", addrS, accS)
		}
	}
	{
		iterated := make(map[string]map[string]string)
		its := make(map[string]NodeIterator)
		// batch create iterators
		for _, s := range storages {
			addrS := s[0]
			if _, ok := its[addrS]; !ok {
				iterated[addrS] = make(map[string]string)
				storageTrie := storageTries[addrS]
				require.NotNil(t, storageTrie)
				its[addrS] = storageTrie.NodeIterator(nil)
			}
		}
		// run one iterator at a time
		count := 0
		for addrS, it := range its {
			for it.Next(true) {
				if !it.Leaf() {
					continue
				}
				k, v := it.LeafKey(), it.LeafBlob()
				k = storageTries[addrS].GetKey(k)
				iterated[addrS][hexutil.Encode(k)] = hexutil.Encode(v)
				count++
			}
		}
		// Compare with the correct answer
		assert.Equal(t, len(storages), count)
		for _, s := range storages {
			addrS, kS, vS := s[0], s[1], s[2]
			assert.Equal(t, vS, iterated[addrS][kS], "addr = %s, k = %s, v = %s", addrS, kS, vS)
		}
	}
}

// Random accounts generator

func randTrie(t *testing.T, r *rand.Rand, numAccounts, maxSlotsPerAccount int) ([][2]string, [][3]string) {
	accounts := make([][2]string, numAccounts)
	storages := make([][3]string, 0)
	for i := 0; i < len(accounts); i++ {
		if i%2 == 0 { // EOA
			accounts[i] = [2]string{randAddr(r).Hex(), hexutil.Encode(randEOA(t, r))}
			continue
		} else { // SCA
			addr := randAddr(r).Hex()
			storage := randStorage(r, addr, maxSlotsPerAccount)
			storageRoot := correctStorageRoot(storage)
			accounts[i] = [2]string{addr, hexutil.Encode(randSCA(t, r, storageRoot))}
			storages = append(storages, storage...)
		}
	}
	return accounts, storages
}

func randStorage(r *rand.Rand, addr string, maxSlotsPerAccount int) [][3]string {
	storage := make([][3]string, r.Intn(maxSlotsPerAccount+1)) // [0, maxSlotsPerAccount)
	for j := 0; j < len(storage); j++ {
		// value with some leading zeros
		value := randHash(r).Bytes()
		value = value[:r.Intn(32)]
		if common.EmptyHash(common.BytesToHash(value)) {
			// Ensure a nonzero slot; See state_object.go:updateStorageTrie where the slot is rather subject to trie.TryDelete() when (value == common.Hash{}).
			// Here we are generating nonzero slot to be trie.TryUpdate().
			value = []byte{0x42}
		} else {
			value = common.BytesToHash(value).Bytes()
		}
		// as in state_object.go:updateStorageTrie
		value, _ = rlp.EncodeToBytes(bytes.TrimLeft(value[:], "\x00"))
		storage[j] = [3]string{addr, randHash(r).Hex(), hexutil.Encode(value)}
	}
	return storage
}

func correctStorageRoot(storage [][3]string) common.Hash {
	trie := newEmptySecureTrie()
	for i := 0; i < len(storage); i++ {
		k, v := []byte(storage[i][1]), []byte(storage[i][2])
		trie.TryUpdate(k, v)
	}
	return trie.Hash()
}

func randEOA(t *testing.T, r *rand.Rand) []byte {
	acc, err := account.NewAccountWithMap(account.ExternallyOwnedAccountType, map[account.AccountValueKeyType]interface{}{
		account.AccountValueKeyNonce:         uint64(r.Int()),
		account.AccountValueKeyBalance:       big.NewInt(r.Int63()),
		account.AccountValueKeyHumanReadable: false,
		account.AccountValueKeyAccountKey:    randAccountKey(r),
	})
	require.NoError(t, err)

	ser := account.NewAccountSerializerWithAccount(acc)
	data, err := rlp.EncodeToBytes(ser)
	require.NoError(t, err)
	return data
}

func randSCA(t *testing.T, r *rand.Rand, storageRoot common.Hash) []byte {
	acc, err := account.NewAccountWithMap(account.SmartContractAccountType, map[account.AccountValueKeyType]interface{}{
		account.AccountValueKeyNonce:         uint64(r.Int()),
		account.AccountValueKeyBalance:       big.NewInt(r.Int63()),
		account.AccountValueKeyHumanReadable: false,
		account.AccountValueKeyAccountKey:    accountkey.NewAccountKeyFail(),
		account.AccountValueKeyStorageRoot:   storageRoot,
		account.AccountValueKeyCodeHash:      randHash(r),
		account.AccountValueKeyCodeInfo:      params.CodeInfo(0x10),
	})
	require.NoError(t, err)

	ser := account.NewAccountSerializerWithAccount(acc)
	data, err := rlp.EncodeToBytes(ser)
	require.NoError(t, err)
	return data
}

func randAccountKey(r *rand.Rand) accountkey.AccountKey {
	ty := r.Intn(5)
	switch ty {
	case 0:
		return accountkey.NewAccountKeyLegacy()
	case 1:
		return accountkey.NewAccountKeyPublicWithValue(randPub(r))
	case 2:
		return accountkey.NewAccountKeyFail()
	case 3:
		n := r.Intn(9) + 1 // [1, MaxNumKeysForMultiSig]
		m := 1
		if n > 1 {
			m = r.Intn(n-1) + 1 // [1, n]
		}
		keys := make(accountkey.WeightedPublicKeys, n)
		for i := 0; i < n; i++ {
			keys[i] = accountkey.NewWeightedPublicKey(uint(r.Intn(10)), (*accountkey.PublicKeySerializable)(randPub(r)))
		}
		return accountkey.NewAccountKeyWeightedMultiSigWithValues(uint(m), keys)
	default:
		return accountkey.NewAccountKeyRoleBasedWithValues([]accountkey.AccountKey{
			randAccountKey(r),
			randAccountKey(r),
			randAccountKey(r),
		})
	}
}

func randPub(r *rand.Rand) *ecdsa.PublicKey {
	privB := make([]byte, 32)
	r.Read(privB)
	priv := crypto.ToECDSAUnsafe(privB)
	return &priv.PublicKey
}

func randAddr(r *rand.Rand) common.Address {
	h := make([]byte, 20)
	r.Read(h)
	return common.BytesToAddress(h)
}

func randHash(r *rand.Rand) common.Hash {
	h := make([]byte, 32)
	r.Read(h)
	return common.BytesToHash(h)
}
