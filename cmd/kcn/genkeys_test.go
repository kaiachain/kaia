// Copyright 2026 The Kaia Authors
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

package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/keystore"
	"github.com/kaiachain/kaia/cmd/utils"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestMain(m *testing.M) {
	keystoreScryptN, keystoreScryptP = keystore.LightScryptN, keystore.LightScryptP
	os.Exit(m.Run())
}

func runGenerateKeys(datadir string) error {
	app := &cli.App{Name: "kcn", Commands: []*cli.Command{GenerateKeysCommand}}
	return app.Run([]string{"kcn", "generate-keys", "--datadir", datadir})
}

type fileSnapshot struct {
	content string
	modTime time.Time
}

// snapshotDir records the content and mtime of every regular file under dir.
func snapshotDir(t *testing.T, dir string) map[string]fileSnapshot {
	snap := make(map[string]fileSnapshot)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		snap[path] = fileSnapshot{string(b), info.ModTime()}
		return nil
	})
	require.NoError(t, err)
	return snap
}

func allTargets(datadir string) []string {
	klayDir := filepath.Join(datadir, "klay")
	opDir := filepath.Join(klayDir, "operator-keys")
	targets := []string{
		filepath.Join(klayDir, "nodekey"), filepath.Join(klayDir, "bls-nodekey"),
		filepath.Join(opDir, "bls-pub"), filepath.Join(opDir, "bls-pop"),
		filepath.Join(opDir, "bls-nodekey.json"), filepath.Join(opDir, "bls-nodekey.pass"),
	}
	for _, n := range ecdsaKeyNames {
		targets = append(targets, filepath.Join(opDir, n+".json"), filepath.Join(opDir, n+".pass"))
	}
	return targets
}

// checkKeySet verifies that every target exists and that the derived files
// match the raw keys. It returns the raw keys.
func checkKeySet(t *testing.T, datadir string) (*ecdsa.PrivateKey, bls.SecretKey) {
	for _, p := range allTargets(datadir) {
		assert.FileExists(t, p)
	}
	klayDir := filepath.Join(datadir, "klay")
	opDir := filepath.Join(klayDir, "operator-keys")

	nodeKey, err := crypto.LoadECDSA(filepath.Join(klayDir, "nodekey"))
	require.NoError(t, err)
	blsKey, err := bls.LoadKey(filepath.Join(klayDir, "bls-nodekey"))
	require.NoError(t, err)

	ksKey, err := decryptV3Keystore(opDir, "nodekey")
	require.NoError(t, err)
	assert.Equal(t, crypto.PubkeyToAddress(nodeKey.PublicKey), crypto.PubkeyToAddress(ksKey.PublicKey))

	// Keystores decrypt with their .pass and hold the raw keys.
	js, err := os.ReadFile(filepath.Join(opDir, "nodekey.json"))
	require.NoError(t, err)
	pw, err := os.ReadFile(filepath.Join(opDir, "nodekey.pass"))
	require.NoError(t, err)
	k, err := keystore.DecryptKey(js, string(pw))
	require.NoError(t, err)
	assert.Equal(t, crypto.FromECDSA(nodeKey), crypto.FromECDSA(k.GetPrivateKey()))

	js, err = os.ReadFile(filepath.Join(opDir, "bls-nodekey.json"))
	require.NoError(t, err)
	pw, err = os.ReadFile(filepath.Join(opDir, "bls-nodekey.pass"))
	require.NoError(t, err)
	bk, err := keystore.DecryptKeyEIP2335(js, string(pw))
	require.NoError(t, err)
	assert.Equal(t, blsKey.Marshal(), bk.SecretKey.Marshal())

	pub, err := readHexFile(filepath.Join(opDir, "bls-pub"))
	require.NoError(t, err)
	assert.Equal(t, blsKey.PublicKey().Marshal(), pub)
	pop, err := readHexFile(filepath.Join(opDir, "bls-pop"))
	require.NoError(t, err)
	sig, err := bls.SignatureFromBytes(pop)
	require.NoError(t, err)
	assert.True(t, bls.PopVerify(blsKey.PublicKey(), sig))

	// Keystore address and file formats: lowercase without 0x.
	for _, n := range ecdsaKeyNames {
		b, err := os.ReadFile(filepath.Join(opDir, n+".json"))
		require.NoError(t, err)
		assert.NotContains(t, string(b), `"address":"0x`)
	}
	b, err := os.ReadFile(filepath.Join(opDir, "bls-pub"))
	require.NoError(t, err)
	assert.False(t, strings.HasPrefix(string(b), "0x"))
	return nodeKey, blsKey
}

func TestGenerateKeys_Fresh(t *testing.T) {
	datadir := t.TempDir()
	require.NoError(t, runGenerateKeys(datadir))
	checkKeySet(t, datadir)
}

func TestGenerateKeys_ReuseRawKeys(t *testing.T) {
	testcases := []struct {
		name   string
		hasBls bool
	}{
		{"nodekey only", false},
		{"nodekey and node-derived bls-nodekey", true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			datadir := t.TempDir()
			klayDir := filepath.Join(datadir, "klay")
			require.NoError(t, os.MkdirAll(klayDir, 0o700))

			// Same files as a node started once on this datadir (node.Config.NodeKey/BlsNodeKey).
			nodeKey, err := crypto.GenerateKey()
			require.NoError(t, err)
			require.NoError(t, crypto.SaveECDSA(filepath.Join(klayDir, "nodekey"), nodeKey))
			wantBls, err := bls.GenerateKey(crypto.FromECDSA(nodeKey))
			require.NoError(t, err)
			if tc.hasBls {
				require.NoError(t, bls.SaveKey(filepath.Join(klayDir, "bls-nodekey"), wantBls))
			}

			require.NoError(t, runGenerateKeys(datadir))
			gotNodeKey, gotBls := checkKeySet(t, datadir)
			assert.Equal(t, crypto.PubkeyToAddress(nodeKey.PublicKey), crypto.PubkeyToAddress(gotNodeKey.PublicKey))
			assert.Equal(t, wantBls.PublicKey().Marshal(), gotBls.PublicKey().Marshal())
		})
	}
}

// TestGenerateKeys_LegacyPath checks that generate-keys uses the raw keys the
// node loads (node.Config.NodeKey/BlsNodeKey), which prefer <datadir>/<file>
// over <datadir>/klay/<file>.
func TestGenerateKeys_LegacyPath(t *testing.T) {
	genKeys := func(t *testing.T) (*ecdsa.PrivateKey, bls.SecretKey) {
		k, err := crypto.GenerateKey()
		require.NoError(t, err)
		b, err := bls.RandKey()
		require.NoError(t, err)
		return k, b
	}
	testcases := []struct {
		name          string
		legacyBls     bool // bls-nodekey at <datadir>/bls-nodekey
		klayKeys      bool // different nodekey and bls-nodekey at <datadir>/klay
		wantBlsDerive bool // bls-nodekey is derived from the legacy nodekey
	}{
		{"legacy only", true, false, false},
		{"legacy and klay differ", true, true, false},
		{"legacy nodekey only", false, false, true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			datadir := t.TempDir()
			klayDir := filepath.Join(datadir, "klay")
			opDir := filepath.Join(klayDir, "operator-keys")

			legacyKey, legacyBls := genKeys(t)
			require.NoError(t, crypto.SaveECDSA(filepath.Join(datadir, "nodekey"), legacyKey))
			if tc.legacyBls {
				require.NoError(t, bls.SaveKey(filepath.Join(datadir, "bls-nodekey"), legacyBls))
			}
			if tc.klayKeys {
				klayKey, klayBls := genKeys(t)
				require.NoError(t, os.MkdirAll(klayDir, 0o700))
				require.NoError(t, crypto.SaveECDSA(filepath.Join(klayDir, "nodekey"), klayKey))
				require.NoError(t, bls.SaveKey(filepath.Join(klayDir, "bls-nodekey"), klayBls))
			}

			// The keys the node itself would load from this datadir.
			cfg := node.Config{DataDir: datadir, Name: utils.ClientIdentifier}
			nodeKey := cfg.NodeKey()
			wantBls := legacyBls
			if tc.wantBlsDerive {
				var err error
				wantBls, err = bls.GenerateKey(crypto.FromECDSA(legacyKey))
				require.NoError(t, err)
			}
			wantID := crypto.PubkeyToAddress(legacyKey.PublicKey)
			require.Equal(t, wantID, crypto.PubkeyToAddress(nodeKey.PublicKey))

			out := captureStdout(t, func() { require.NoError(t, runGenerateKeys(datadir)) })
			assert.Contains(t, out, wantID.Hex())
			assert.Contains(t, out, filepath.Join(datadir, "nodekey")+" ["+statusReused+"]")
			assert.NotContains(t, out, "WARNING")
			if !tc.klayKeys {
				assert.NoFileExists(t, filepath.Join(klayDir, "nodekey"))
			}

			ksKey, err := decryptV3Keystore(opDir, "nodekey")
			require.NoError(t, err)
			assert.Equal(t, wantID, crypto.PubkeyToAddress(ksKey.PublicKey))
			pub, err := readHexFile(filepath.Join(opDir, "bls-pub"))
			require.NoError(t, err)
			assert.Equal(t, wantBls.PublicKey().Marshal(), pub)

			// The node, started after generate-keys, loads the same BLS key.
			assert.Equal(t, wantBls.PublicKey().Marshal(), cfg.BlsNodeKey().PublicKey().Marshal())
		})
	}
}

func TestGenerateKeys_SecondRunKeepsFiles(t *testing.T) {
	datadir := t.TempDir()
	require.NoError(t, runGenerateKeys(datadir))
	before := snapshotDir(t, datadir)

	out := captureStdout(t, func() { require.NoError(t, runGenerateKeys(datadir)) })
	assert.Equal(t, before, snapshotDir(t, datadir))
	assert.NotContains(t, out, statusCreated+"]")
	assert.Equal(t, 2, strings.Count(out, "["+statusReused+"]"))  // nodekey, bls-nodekey
	assert.Equal(t, 7, strings.Count(out, "["+statusSkipped+"]")) // keystores
	assert.Equal(t, 2, strings.Count(out, " "+statusSkipped+"]")) // bls-pub, bls-pop
}

// captureStdout returns what fn prints to stdout.
func captureStdout(t *testing.T, fn func()) string {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()
	fn()
	require.NoError(t, w.Close())
	b, err := io.ReadAll(r)
	require.NoError(t, err)
	return string(b)
}

func TestGenerateKeys_PartialOperatorKeys(t *testing.T) {
	datadir := t.TempDir()
	require.NoError(t, runGenerateKeys(datadir))
	opDir := filepath.Join(datadir, "klay", "operator-keys")
	managerJSON, err := os.ReadFile(filepath.Join(opDir, "manager.json"))
	require.NoError(t, err)
	require.NoError(t, os.Remove(filepath.Join(opDir, "voter.json")))
	require.NoError(t, os.Remove(filepath.Join(opDir, "voter.pass")))

	require.NoError(t, runGenerateKeys(datadir))
	checkKeySet(t, datadir)
	got, err := os.ReadFile(filepath.Join(opDir, "manager.json"))
	require.NoError(t, err)
	assert.Equal(t, managerJSON, got)
}

func TestGenerateKeys_Mismatch(t *testing.T) {
	newDatadir := func(t *testing.T) (string, string) {
		datadir := t.TempDir()
		klayDir := filepath.Join(datadir, "klay")
		opDir := filepath.Join(klayDir, "operator-keys")
		require.NoError(t, os.MkdirAll(opDir, 0o700))
		nodeKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		require.NoError(t, crypto.SaveECDSA(filepath.Join(klayDir, "nodekey"), nodeKey))
		return datadir, opDir
	}
	writeFiles := func(t *testing.T, w []fileWrite) {
		for _, f := range w {
			require.NoError(t, os.WriteFile(f.path, []byte(f.content), f.perm))
		}
	}
	writeKeystore := func(t *testing.T, opDir, name string) {
		other, err := crypto.GenerateKey()
		require.NoError(t, err)
		w, err := encryptV3(opDir, name, other)
		require.NoError(t, err)
		writeFiles(t, w)
	}
	writeBlsKeystore := func(t *testing.T, opDir string) {
		other, err := bls.RandKey()
		require.NoError(t, err)
		w, err := encryptEIP2335(opDir, other)
		require.NoError(t, err)
		writeFiles(t, w)
	}
	// editJSON rewrites one field of a keystore JSON; path is the nested key path.
	editJSON := func(t *testing.T, file string, value string, path ...string) {
		b, err := os.ReadFile(file)
		require.NoError(t, err)
		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		obj := m
		for _, k := range path[:len(path)-1] {
			obj = obj[k].(map[string]any)
		}
		obj[path[len(path)-1]] = value
		b, err = json.Marshal(m)
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(file, b, 0o600))
	}
	wrongPass := func(t *testing.T, opDir, name string) {
		require.NoError(t, os.WriteFile(filepath.Join(opDir, name+".pass"), []byte("wrong"), 0o600))
	}

	testcases := []struct {
		name    string
		prepare func(t *testing.T, opDir string)
		errStr  string
	}{
		{
			name:    "nodekey.json of another key",
			prepare: func(t *testing.T, opDir string) { writeKeystore(t, opDir, "nodekey") },
			errStr:  "nodekey.json does not match",
		},
		{
			name: "bls-pub of another key",
			prepare: func(t *testing.T, opDir string) {
				other, err := bls.RandKey()
				require.NoError(t, err)
				pub := hex.EncodeToString(other.PublicKey().Marshal())
				require.NoError(t, os.WriteFile(filepath.Join(opDir, "bls-pub"), []byte(pub), 0o644))
			},
			errStr: "bls-pub does not match",
		},
		{
			name: "bls-pop of another key",
			prepare: func(t *testing.T, opDir string) {
				other, err := bls.RandKey()
				require.NoError(t, err)
				pop := hex.EncodeToString(bls.PopProve(other).Marshal())
				require.NoError(t, os.WriteFile(filepath.Join(opDir, "bls-pop"), []byte(pop), 0o644))
			},
			errStr: "bls-pop is not a valid proof-of-possession",
		},
		{
			name: "bls-nodekey.json without pubkey",
			prepare: func(t *testing.T, opDir string) {
				ks := `{"crypto":{},"pubkey":"","uuid":"00000000-0000-0000-0000-000000000000","version":4}`
				require.NoError(t, os.WriteFile(filepath.Join(opDir, "bls-nodekey.json"), []byte(ks), 0o600))
				require.NoError(t, os.WriteFile(filepath.Join(opDir, "bls-nodekey.pass"), []byte("x"), 0o600))
			},
			errStr: "missing pubkey",
		},
		{
			name: "nodekey.json with a wrong password",
			prepare: func(t *testing.T, opDir string) {
				writeKeystore(t, opDir, "nodekey")
				wrongPass(t, opDir, "nodekey")
			},
			errStr: "decrypt",
		},
		{
			name: "operator keystore with a wrong password",
			prepare: func(t *testing.T, opDir string) {
				writeKeystore(t, opDir, "manager")
				wrongPass(t, opDir, "manager")
			},
			errStr: "decrypt",
		},
		{
			name: "operator keystore with a corrupted ciphertext",
			prepare: func(t *testing.T, opDir string) {
				writeKeystore(t, opDir, "voter")
				editJSON(t, filepath.Join(opDir, "voter.json"), strings.Repeat("00", 32), "crypto", "ciphertext")
			},
			errStr: "decrypt",
		},
		{
			name: "operator keystore with another address field",
			prepare: func(t *testing.T, opDir string) {
				writeKeystore(t, opDir, "reward")
				editJSON(t, filepath.Join(opDir, "reward.json"), strings.Repeat("11", 20), "address")
			},
			errStr: "is not the address of the key",
		},
		{
			name: "bls-nodekey.json with a wrong password",
			prepare: func(t *testing.T, opDir string) {
				writeBlsKeystore(t, opDir)
				wrongPass(t, opDir, "bls-nodekey")
			},
			errStr: "decrypt",
		},
		{
			name: "bls-nodekey.json with a corrupted ciphertext",
			prepare: func(t *testing.T, opDir string) {
				writeBlsKeystore(t, opDir)
				editJSON(t, filepath.Join(opDir, "bls-nodekey.json"), strings.Repeat("00", 32), "crypto", "cipher", "message")
			},
			errStr: "decrypt",
		},
		{
			name: "bls-nodekey.json with another pubkey field",
			prepare: func(t *testing.T, opDir string) {
				writeBlsKeystore(t, opDir)
				other, err := bls.RandKey()
				require.NoError(t, err)
				editJSON(t, filepath.Join(opDir, "bls-nodekey.json"), hex.EncodeToString(other.PublicKey().Marshal()), "pubkey")
			},
			errStr: "is not the public key of the key",
		},
		{
			name:    "bls-nodekey.json of another key",
			prepare: func(t *testing.T, opDir string) { writeBlsKeystore(t, opDir) },
			errStr:  "bls-nodekey.json does not match",
		},
		{
			name: "operator keystore without .pass",
			prepare: func(t *testing.T, opDir string) {
				writeKeystore(t, opDir, "voter")
				require.NoError(t, os.Remove(filepath.Join(opDir, "voter.pass")))
			},
			errStr: "inconsistent key files",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			datadir, opDir := newDatadir(t)
			tc.prepare(t, opDir)
			before := snapshotDir(t, datadir)

			err := runGenerateKeys(datadir)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.errStr)
			assert.Equal(t, before, snapshotDir(t, datadir), "no file must be written")
		})
	}
}

func TestLoadKey_Datadir(t *testing.T) {
	datadir := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(datadir, "klay"), 0o700))
	fileKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	require.NoError(t, crypto.SaveECDSA(filepath.Join(datadir, "klay", "nodekey"), fileKey))
	flagKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	load := func(args ...string) (*ecdsa.PrivateKey, string, error) {
		var (
			key     *ecdsa.PrivateKey
			datadir string
			err     error
		)
		cmd := &cli.Command{
			Name:  "load",
			Flags: abv2Flags(),
			Action: func(ctx *cli.Context) error {
				datadir = ctx.String("datadir")
				key, err = loadKey(ctx)
				return nil
			},
		}
		app := &cli.App{Name: "kcn", Commands: []*cli.Command{cmd}}
		require.NoError(t, app.Run(append([]string{"kcn", "load"}, args...)))
		return key, datadir, err
	}

	// --datadir: key comes from <datadir>/klay/nodekey.
	key, _, err := load("--datadir", datadir)
	require.NoError(t, err)
	assert.Equal(t, crypto.FromECDSA(fileKey), crypto.FromECDSA(key))

	// --private-key wins over --datadir.
	key, _, err = load("--datadir", datadir, "--private-key", "0x"+hex.EncodeToString(crypto.FromECDSA(flagKey)))
	require.NoError(t, err)
	assert.Equal(t, crypto.FromECDSA(flagKey), crypto.FromECDSA(key))

	// Without flags, the key is read from the rpm data directory.
	_, got, _ := load()
	assert.Equal(t, defaultDatadir, got)
	_, _, err = load("--datadir", filepath.Join(datadir, "missing"))
	assert.ErrorContains(t, err, "check --datadir")

	// <datadir>/nodekey wins over <datadir>/klay/nodekey, as in the node.
	legacyKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	require.NoError(t, crypto.SaveECDSA(filepath.Join(datadir, "nodekey"), legacyKey))
	key, _, err = load("--datadir", datadir)
	require.NoError(t, err)
	assert.Equal(t, crypto.FromECDSA(legacyKey), crypto.FromECDSA(key))

	// Every command that loads the node key accepts --datadir.
	for _, cmd := range append(ValOpsCommand.Subcommands, SignCreateNodeCommand) {
		// cli adds a help subcommand once ValOpsCommand has been run in another test.
		if cmd == GenerateKeysCommand || cmd.Name == "help" {
			continue
		}
		var found bool
		for _, f := range cmd.Flags {
			if f == abv2DatadirFlag {
				found = true
			}
		}
		assert.True(t, found, "%s has no --datadir flag", cmd.Name)
	}
}

func TestWriteNewFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nodekey")
	require.NoError(t, writeNewFile(path, "first", 0o600))

	// An existing file is never overwritten.
	err := writeNewFile(path, "second", 0o600)
	assert.ErrorContains(t, err, "refusing to overwrite")
	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "first", string(got))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}
