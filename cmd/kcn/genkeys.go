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
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/kaiachain/kaia/accounts/keystore"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/urfave/cli/v2"
)

var genkeysDatadirFlag = &cli.StringFlag{
	Name:     "datadir",
	Usage:    "Node data directory; keys are written under <datadir>/klay",
	Required: true,
}

// ecdsaKeyNames are the secp256k1 keys written to klay/operator-keys as v3
// keystores. nodekey is also kept as raw hex in klay/ for the node to load.
// reward is unused under PublicDelegation; mev-reward is for the auction
// contract, not createNode.
var ecdsaKeyNames = []string{"nodekey", "manager", "voter", "reward", "cnstaking-owner", "mev-reward"}

// GenerateKeysCommand generates the full validator onboarding key set offline.
var GenerateKeysCommand = &cli.Command{
	Name:   "generate-keys",
	Usage:  "Generate all validator onboarding keys into <datadir>/klay (offline)",
	Flags:  []cli.Flag{genkeysDatadirFlag},
	Action: generateKeysAction,
}

func writeFile(path, content string, perm os.FileMode) error {
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// randomPassword returns a fresh high-entropy password for one keystore.
func randomPassword() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// writeV3 writes an ECDSA key as a v3 keystore JSON plus its random password file.
func writeV3(opDir, name string, priv *ecdsa.PrivateKey) error {
	pw, err := randomPassword()
	if err != nil {
		return err
	}
	k := &keystore.KeyV3{Id: uuid.New(), Address: crypto.PubkeyToAddress(priv.PublicKey), PrivateKey: priv}
	js, err := keystore.EncryptKeyV3(k, pw, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return fmt.Errorf("encrypt %s: %w", name, err)
	}
	if err := writeFile(filepath.Join(opDir, name+".json"), string(js), 0o600); err != nil {
		return err
	}
	return writeFile(filepath.Join(opDir, name+".pass"), pw, 0o600)
}

func generateKeysAction(ctx *cli.Context) error {
	datadir := ctx.String("datadir")
	klayDir := filepath.Join(datadir, "klay")
	opDir := filepath.Join(klayDir, "operator-keys")

	nodekeyHexPath := filepath.Join(klayDir, "nodekey")
	blsNodekeyHexPath := filepath.Join(klayDir, "bls-nodekey")
	blsPubPath := filepath.Join(opDir, "bls-pub")
	blsPopPath := filepath.Join(opDir, "bls-pop")

	// No-overwrite precheck (atomic): if any target exists, write nothing.
	targets := []string{
		nodekeyHexPath, blsNodekeyHexPath, blsPubPath, blsPopPath,
		filepath.Join(opDir, "bls-nodekey.json"), filepath.Join(opDir, "bls-nodekey.pass"),
	}
	for _, n := range ecdsaKeyNames {
		targets = append(targets, filepath.Join(opDir, n+".json"), filepath.Join(opDir, n+".pass"))
	}
	for _, p := range targets {
		if _, err := os.Stat(p); err == nil {
			return fmt.Errorf("refusing to overwrite existing file: %s (use a fresh --datadir or remove existing keys)", p)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s: %w", p, err)
		}
	}

	// Generate all keys before writing, so a generation failure leaves nothing.
	ecdsaKeys := make(map[string]*ecdsa.PrivateKey, len(ecdsaKeyNames))
	for _, n := range ecdsaKeyNames {
		k, err := crypto.GenerateKey()
		if err != nil {
			return fmt.Errorf("generate %s key: %w", n, err)
		}
		ecdsaKeys[n] = k
	}
	blsKey, err := bls.RandKey()
	if err != nil {
		return fmt.Errorf("generate bls key: %w", err)
	}

	if err := os.MkdirAll(opDir, 0o700); err != nil {
		return fmt.Errorf("create %s: %w", opDir, err)
	}

	// ECDSA keys -> v3 keystore (+ .pass) in operator-keys; nodekey also raw hex in klay/.
	for _, n := range ecdsaKeyNames {
		if err := writeV3(opDir, n, ecdsaKeys[n]); err != nil {
			return err
		}
	}
	nodeKey := ecdsaKeys["nodekey"]
	if err := writeFile(nodekeyHexPath, hex.EncodeToString(crypto.FromECDSA(nodeKey)), 0o600); err != nil {
		return err
	}

	// BLS -> raw hex in klay/, EIP-2335 keystore (+ .pass), and public pub/pop hex.
	if err := writeFile(blsNodekeyHexPath, hex.EncodeToString(blsKey.Marshal()), 0o600); err != nil {
		return err
	}
	blsPw, err := randomPassword()
	if err != nil {
		return err
	}
	blsJSON, err := keystore.EncryptKeyEIP2335(keystore.NewKeyEIP2335(blsKey), blsPw, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return fmt.Errorf("encrypt bls-nodekey: %w", err)
	}
	if err := writeFile(filepath.Join(opDir, "bls-nodekey.json"), string(blsJSON), 0o600); err != nil {
		return err
	}
	if err := writeFile(filepath.Join(opDir, "bls-nodekey.pass"), blsPw, 0o600); err != nil {
		return err
	}
	blsPub := hex.EncodeToString(blsKey.PublicKey().Marshal())
	blsPop := hex.EncodeToString(bls.PopProve(blsKey).Marshal())
	if err := writeFile(blsPubPath, blsPub, 0o644); err != nil {
		return err
	}
	if err := writeFile(blsPopPath, blsPop, 0o644); err != nil {
		return err
	}

	// Print public values needed for createNode / onboarding; secrets stay in files.
	nodeID := crypto.PubkeyToAddress(nodeKey.PublicKey).Hex()
	fmt.Printf("Generated validator keys under %s\n\n", klayDir)
	fmt.Println("node:")
	fmt.Printf("  hex        : %s\n", nodekeyHexPath)
	fmt.Printf("  keystore   : %s (+ .pass)\n", filepath.Join(opDir, "nodekey.json"))
	fmt.Printf("  nodeId     : %s\n", nodeID)
	fmt.Println("bls:")
	fmt.Printf("  hex        : %s\n", blsNodekeyHexPath)
	fmt.Printf("  keystore   : %s (+ .pass, EIP-2335)\n", filepath.Join(opDir, "bls-nodekey.json"))
	fmt.Printf("  publicKey  : 0x%s\n", blsPub)
	fmt.Printf("  pop        : 0x%s\n", blsPop)
	for _, n := range ecdsaKeyNames {
		if n == "nodekey" {
			continue
		}
		fmt.Printf("%s:\n", n)
		fmt.Printf("  keystore   : %s (+ .pass)\n", filepath.Join(opDir, n+".json"))
		fmt.Printf("  address    : %s\n", crypto.PubkeyToAddress(ecdsaKeys[n].PublicKey).Hex())
	}
	fmt.Println("\nNote: reward is unused when PublicDelegation is enabled. Each .pass holds that keystore's password.")
	return nil
}
