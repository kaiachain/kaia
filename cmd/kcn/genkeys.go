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
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/kaiachain/kaia/accounts/keystore"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/node"
	"github.com/urfave/cli/v2"
)

var genkeysDatadirFlag = &cli.StringFlag{
	Name:  "datadir",
	Usage: "Node data directory (rpm installs use the default; for a tar install, DATA_DIR in conf/kcnd.conf). Keys are written under <datadir>/klay",
	Value: defaultDatadir,
}

// ecdsaKeyNames are the secp256k1 keys written to klay/operator-keys as v3
// keystores. nodekey is also kept as raw hex in klay/ for the node to load.
// reward is unused under PublicDelegation; mev-reward is for the auction
// contract, not createNode.
var ecdsaKeyNames = []string{"nodekey", "manager", "voter", "reward", "cnstaking-owner", "mev-reward"}

// Scrypt parameters for the generated keystores.
var (
	keystoreScryptN = keystore.StandardScryptN
	keystoreScryptP = keystore.StandardScryptP
)

// Per-key status reported by generate-keys.
const (
	statusCreated = "created" // generated or derived, and written
	statusReused  = "reused"  // existing raw key loaded from klay/
	statusSkipped = "skipped" // existing file verified and kept
)

// GenerateKeysCommand generates the validator onboarding key set offline. It
// reuses the raw keys a started node has already written to <datadir>/klay and
// creates only the missing files.
var GenerateKeysCommand = &cli.Command{
	Name:   "generate-keys",
	Usage:  "Generate missing validator onboarding keys into <datadir>/klay, reusing existing ones (offline)",
	Flags:  []cli.Flag{genkeysDatadirFlag},
	Action: generateKeysAction,
}

// fileWrite is a file write deferred until every check has passed.
type fileWrite struct {
	path    string
	content string
	perm    os.FileMode
}

// writeNewFile creates path with content and fails if the file already exists,
// so a key file is never overwritten, even one created after the checks ran.
func writeNewFile(path, content string, perm os.FileMode) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, perm)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%s already exists; refusing to overwrite it", path)
		}
		return err
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat %s: %w", path, err)
}

// keystorePairExists reports whether <name>.json and <name>.pass both exist.
// Only one of the two existing is an error.
func keystorePairExists(opDir, name string) (bool, error) {
	jsonPath := filepath.Join(opDir, name+".json")
	passPath := filepath.Join(opDir, name+".pass")
	hasJSON, err := fileExists(jsonPath)
	if err != nil {
		return false, err
	}
	hasPass, err := fileExists(passPath)
	if err != nil {
		return false, err
	}
	switch {
	case hasJSON && !hasPass:
		return false, fmt.Errorf("inconsistent key files: %s exists but %s is missing", jsonPath, passPath)
	case !hasJSON && hasPass:
		return false, fmt.Errorf("inconsistent key files: %s exists but %s is missing", passPath, jsonPath)
	}
	return hasJSON, nil
}

// readKeystorePair reads <name>.json and <name>.pass under opDir.
func readKeystorePair(opDir, name string) (js []byte, pw string, err error) {
	jsonPath := filepath.Join(opDir, name+".json")
	passPath := filepath.Join(opDir, name+".pass")
	if js, err = os.ReadFile(jsonPath); err != nil {
		return nil, "", fmt.Errorf("read %s: %w", jsonPath, err)
	}
	b, err := os.ReadFile(passPath)
	if err != nil {
		return nil, "", fmt.Errorf("read %s: %w", passPath, err)
	}
	return js, string(b), nil
}

// decryptV3Keystore decrypts <name>.json with <name>.pass and checks that its
// address field is the address of the decrypted key.
func decryptV3Keystore(opDir, name string) (*ecdsa.PrivateKey, error) {
	jsonPath := filepath.Join(opDir, name+".json")
	js, pw, err := readKeystorePair(opDir, name)
	if err != nil {
		return nil, err
	}
	k, err := keystore.DecryptKey(js, pw)
	if err != nil {
		return nil, fmt.Errorf("decrypt %s with %s.pass: %w (nothing written)", jsonPath, name, err)
	}
	priv := k.GetPrivateKey()
	if priv == nil {
		return nil, fmt.Errorf("decrypt %s: no key (nothing written)", jsonPath)
	}
	if addr := crypto.PubkeyToAddress(priv.PublicKey); k.GetAddress() != addr {
		return nil, fmt.Errorf("%s: address field %s is not the address of the key %s (nothing written)", jsonPath, k.GetAddress().Hex(), addr.Hex())
	}
	return priv, nil
}

// decryptEIP2335Keystore decrypts bls-nodekey.json with bls-nodekey.pass and
// checks that its pubkey field is the public key of the decrypted key.
// EIP-2335 requires the pubkey field for a private-key keystore.
func decryptEIP2335Keystore(opDir string) (bls.SecretKey, error) {
	jsonPath := filepath.Join(opDir, "bls-nodekey.json")
	js, pw, err := readKeystorePair(opDir, "bls-nodekey")
	if err != nil {
		return nil, err
	}
	var field struct {
		Pubkey string `json:"pubkey"`
	}
	if err := json.Unmarshal(js, &field); err != nil {
		return nil, fmt.Errorf("parse %s: %w (nothing written)", jsonPath, err)
	}
	if field.Pubkey == "" {
		return nil, fmt.Errorf("parse %s: missing pubkey (nothing written)", jsonPath)
	}
	k, err := keystore.DecryptKeyEIP2335(js, pw)
	if err != nil {
		return nil, fmt.Errorf("decrypt %s with bls-nodekey.pass: %w (nothing written)", jsonPath, err)
	}
	pub := hex.EncodeToString(k.SecretKey.PublicKey().Marshal())
	if strings.ToLower(strings.TrimPrefix(field.Pubkey, "0x")) != pub {
		return nil, fmt.Errorf("%s: pubkey field %s is not the public key of the key %s (nothing written)", jsonPath, field.Pubkey, pub)
	}
	return k.SecretKey, nil
}

// readHexFile decodes a file that holds one hex value, with or without 0x.
func readHexFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	v, err := hex.DecodeString(strings.TrimPrefix(strings.TrimSpace(string(b)), "0x"))
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return v, nil
}

// randomPassword returns a fresh high-entropy password for one keystore.
func randomPassword() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// encryptV3 encrypts an ECDSA key as a v3 keystore and returns the writes for
// the keystore JSON and its random password file.
func encryptV3(opDir, name string, priv *ecdsa.PrivateKey) ([]fileWrite, error) {
	pw, err := randomPassword()
	if err != nil {
		return nil, err
	}
	k := &keystore.KeyV3{Id: uuid.New(), Address: crypto.PubkeyToAddress(priv.PublicKey), PrivateKey: priv}
	js, err := keystore.EncryptKeyV3(k, pw, keystoreScryptN, keystoreScryptP)
	if err != nil {
		return nil, fmt.Errorf("encrypt %s: %w", name, err)
	}
	return []fileWrite{
		{filepath.Join(opDir, name+".json"), string(js), 0o600},
		{filepath.Join(opDir, name+".pass"), pw, 0o600},
	}, nil
}

// encryptEIP2335 encrypts the BLS key as an EIP-2335 keystore and returns the
// writes for the keystore JSON and its random password file.
func encryptEIP2335(opDir string, sk bls.SecretKey) ([]fileWrite, error) {
	pw, err := randomPassword()
	if err != nil {
		return nil, err
	}
	js, err := keystore.EncryptKeyEIP2335(keystore.NewKeyEIP2335(sk), pw, keystoreScryptN, keystoreScryptP)
	if err != nil {
		return nil, fmt.Errorf("encrypt bls-nodekey: %w", err)
	}
	return []fileWrite{
		{filepath.Join(opDir, "bls-nodekey.json"), string(js), 0o600},
		{filepath.Join(opDir, "bls-nodekey.pass"), pw, 0o600},
	}, nil
}

func generateKeysAction(ctx *cli.Context) error {
	datadir := ctx.String("datadir")
	if datadir == "" {
		return errors.New("--datadir must not be empty")
	}
	fmt.Printf("datadir: %s (%s)\n", datadir, datadirSource(ctx))
	klayDir := filepath.Join(datadir, "klay")
	opDir := filepath.Join(klayDir, "operator-keys")

	nodekeyHexPath := nodeFilePath(datadir, node.DatadirPrivateKey)
	blsNodekeyHexPath := nodeFilePath(datadir, node.DatadirBlsSecretKey)
	nodekeyJSONPath := filepath.Join(opDir, "nodekey.json")
	blsJSONPath := filepath.Join(opDir, "bls-nodekey.json")
	blsPubPath := filepath.Join(opDir, "bls-pub")
	blsPopPath := filepath.Join(opDir, "bls-pop")

	// Check every file and prepare every write first, so a failed check writes nothing.
	var writes []fileWrite
	status := make(map[string]string)

	// Raw nodekey: reuse the one the node loads, else generate.
	var nodeKey *ecdsa.PrivateKey
	if ok, err := fileExists(nodekeyHexPath); err != nil {
		return err
	} else if ok {
		if nodeKey, err = crypto.LoadECDSA(nodekeyHexPath); err != nil {
			return fmt.Errorf("load %s: %w", nodekeyHexPath, err)
		}
		status["nodekey"] = statusReused
	} else {
		if nodeKey, err = crypto.GenerateKey(); err != nil {
			return fmt.Errorf("generate nodekey: %w", err)
		}
		writes = append(writes, fileWrite{nodekeyHexPath, hex.EncodeToString(crypto.FromECDSA(nodeKey)), 0o600})
		status["nodekey"] = statusCreated
	}
	nodeID := crypto.PubkeyToAddress(nodeKey.PublicKey)

	// Raw bls-nodekey: reuse; else derive from an existing nodekey the same way
	// the node does (node.Config.BlsNodeKey), so both end up with the same key;
	// else generate.
	var blsKey bls.SecretKey
	if ok, err := fileExists(blsNodekeyHexPath); err != nil {
		return err
	} else if ok {
		if blsKey, err = bls.LoadKey(blsNodekeyHexPath); err != nil {
			return fmt.Errorf("load %s: %w", blsNodekeyHexPath, err)
		}
		status["bls-nodekey"] = statusReused
	} else {
		if status["nodekey"] == statusReused {
			blsKey, err = bls.GenerateKey(crypto.FromECDSA(nodeKey))
		} else {
			blsKey, err = bls.RandKey()
		}
		if err != nil {
			return fmt.Errorf("generate bls key: %w", err)
		}
		writes = append(writes, fileWrite{blsNodekeyHexPath, hex.EncodeToString(blsKey.Marshal()), 0o600})
		status["bls-nodekey"] = statusCreated
	}
	blsPubBytes := blsKey.PublicKey().Marshal()
	blsPub := hex.EncodeToString(blsPubBytes)

	// nodekey keystore: keep it if it decrypts to the raw nodekey.
	if ok, err := keystorePairExists(opDir, "nodekey"); err != nil {
		return err
	} else if ok {
		k, err := decryptV3Keystore(opDir, "nodekey")
		if err != nil {
			return err
		}
		if addr := crypto.PubkeyToAddress(k.PublicKey); addr != nodeID {
			return fmt.Errorf("%s does not match %s: address %s, want %s (nothing written)", nodekeyJSONPath, nodekeyHexPath, addr.Hex(), nodeID.Hex())
		}
		status["nodekey.json"] = statusSkipped
	} else {
		w, err := encryptV3(opDir, "nodekey", nodeKey)
		if err != nil {
			return err
		}
		writes = append(writes, w...)
		status["nodekey.json"] = statusCreated
	}

	// bls-nodekey keystore: keep it if it decrypts to the raw bls-nodekey.
	if ok, err := keystorePairExists(opDir, "bls-nodekey"); err != nil {
		return err
	} else if ok {
		k, err := decryptEIP2335Keystore(opDir)
		if err != nil {
			return err
		}
		if pub := k.PublicKey().Marshal(); !bytes.Equal(pub, blsPubBytes) {
			return fmt.Errorf("%s does not match %s: pubkey %x, want %s (nothing written)", blsJSONPath, blsNodekeyHexPath, pub, blsPub)
		}
		status["bls-nodekey.json"] = statusSkipped
	} else {
		w, err := encryptEIP2335(opDir, blsKey)
		if err != nil {
			return err
		}
		writes = append(writes, w...)
		status["bls-nodekey.json"] = statusCreated
	}

	// bls-pub: keep it if it equals the BLS public key.
	if ok, err := fileExists(blsPubPath); err != nil {
		return err
	} else if ok {
		pub, err := readHexFile(blsPubPath)
		if err != nil {
			return err
		}
		if !bytes.Equal(pub, blsPubBytes) {
			return fmt.Errorf("%s does not match %s: %x, want %s (nothing written)", blsPubPath, blsNodekeyHexPath, pub, blsPub)
		}
		status["bls-pub"] = statusSkipped
	} else {
		writes = append(writes, fileWrite{blsPubPath, blsPub, 0o644})
		status["bls-pub"] = statusCreated
	}

	// bls-pop: keep it if it is a valid proof-of-possession for the BLS public key.
	var blsPop string
	if ok, err := fileExists(blsPopPath); err != nil {
		return err
	} else if ok {
		pop, err := readHexFile(blsPopPath)
		if err != nil {
			return err
		}
		sig, err := bls.SignatureFromBytes(pop)
		if err != nil || !bls.PopVerify(blsKey.PublicKey(), sig) {
			return fmt.Errorf("%s is not a valid proof-of-possession for %s (nothing written)", blsPopPath, blsNodekeyHexPath)
		}
		blsPop = hex.EncodeToString(pop)
		status["bls-pop"] = statusSkipped
	} else {
		blsPop = hex.EncodeToString(bls.PopProve(blsKey).Marshal())
		writes = append(writes, fileWrite{blsPopPath, blsPop, 0o644})
		status["bls-pop"] = statusCreated
	}

	// Operator keys: keep existing keystores, generate the missing ones.
	addrs := make(map[string]common.Address, len(ecdsaKeyNames))
	for _, n := range ecdsaKeyNames {
		if n == "nodekey" {
			continue
		}
		if ok, err := keystorePairExists(opDir, n); err != nil {
			return err
		} else if ok {
			k, err := decryptV3Keystore(opDir, n)
			if err != nil {
				return err
			}
			addrs[n] = crypto.PubkeyToAddress(k.PublicKey)
			status[n] = statusSkipped
			continue
		}
		k, err := crypto.GenerateKey()
		if err != nil {
			return fmt.Errorf("generate %s key: %w", n, err)
		}
		w, err := encryptV3(opDir, n, k)
		if err != nil {
			return err
		}
		writes = append(writes, w...)
		addrs[n] = crypto.PubkeyToAddress(k.PublicKey)
		status[n] = statusCreated
	}

	if len(writes) > 0 {
		if err := os.MkdirAll(opDir, 0o700); err != nil {
			return fmt.Errorf("create %s: %w", opDir, err)
		}
	}
	for _, w := range writes {
		if err := writeNewFile(w.path, w.content, w.perm); err != nil {
			return fmt.Errorf("write %s: %w", w.path, err)
		}
	}

	// Print public values needed for createNode / onboarding; secrets stay in files.
	fmt.Printf("Validator keys under %s\n", klayDir)
	fmt.Println("(created: generated and written, reused: existing raw key loaded, skipped: existing file verified and kept)")
	fmt.Println()
	fmt.Println("node:")
	fmt.Printf("  hex        : %s [%s]\n", nodekeyHexPath, status["nodekey"])
	fmt.Printf("  keystore   : %s (+ .pass) [%s]\n", nodekeyJSONPath, status["nodekey.json"])
	fmt.Printf("  nodeId     : %s\n", nodeID.Hex())
	fmt.Println("bls:")
	fmt.Printf("  hex        : %s [%s]\n", blsNodekeyHexPath, status["bls-nodekey"])
	fmt.Printf("  keystore   : %s (+ .pass, EIP-2335) [%s]\n", blsJSONPath, status["bls-nodekey.json"])
	fmt.Printf("  publicKey  : 0x%s [bls-pub %s]\n", blsPub, status["bls-pub"])
	fmt.Printf("  pop        : 0x%s [bls-pop %s]\n", blsPop, status["bls-pop"])
	for _, n := range ecdsaKeyNames {
		if n == "nodekey" {
			continue
		}
		fmt.Printf("%s:\n", n)
		fmt.Printf("  keystore   : %s (+ .pass) [%s]\n", filepath.Join(opDir, n+".json"), status[n])
		fmt.Printf("  address    : %s\n", addrs[n].Hex())
	}
	fmt.Println("\nNote: reward is unused when PublicDelegation is enabled. Each .pass holds that keystore's password.")
	if status["nodekey"] == statusCreated {
		fmt.Printf("\nWARNING: no nodekey was found at %s or %s, so a new node ID was created. If this node has already run, check the data directory.\n", filepath.Join(datadir, node.DatadirPrivateKey), filepath.Join(klayDir, node.DatadirPrivateKey))
	}
	return nil
}
