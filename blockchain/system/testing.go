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

package system

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	addressbookv2contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/AddressBookV2"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

// MakeTestPermissionlessConfig creates a test AllocPermissionlessConfig with n validators.
func MakeTestPermissionlessConfig(t *testing.T, n int) *AllocPermissionlessConfig {
	t.Helper()

	ownerKey, _ := crypto.GenerateKey()
	owner := crypto.PubkeyToAddress(ownerKey.PublicKey)

	nodeIds := make([]common.Address, n)
	nodeInfos := make([]addressbookv2contract.NodeInfo, n)
	stakeAmts := make([]*big.Int, n)
	stakeAmt := new(big.Int).Mul(big.NewInt(5_000_000), big.NewInt(params.KAIA))

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		_, pub, pop := MakeTestBlsKey()

		nodeIds[i] = addr
		stakeAmts[i] = new(big.Int).Set(stakeAmt)
		nodeInfos[i] = addressbookv2contract.NodeInfo{
			Manager:       addr,
			RewardAddress: common.BytesToAddress(crypto.Keccak256(addr.Bytes())),
			VoterAddress:  addr,
			TimeoutAt:     new(big.Int),
			GcId:          big.NewInt(int64(i + 1)),
			BlsInfo: addressbookv2contract.BlsPublicKeyInfo{
				PublicKey: pub,
				Pop:       pop,
			},
			State: valset.ValActive.ToUint8(),
		}
	}

	return &AllocPermissionlessConfig{
		Owner:     owner,
		NodeIds:   nodeIds,
		NodeInfos: nodeInfos,
		StakeAmts: stakeAmts,
		DataConfig: addressbookv2contract.IABv2DataContractInitData{
			InitialOwner:           owner,
			ExitThreshold:          big.NewInt(2),
			PauseTimeout:           big.NewInt(8 * 3600),   // 8h
			IdleTimeout:            big.NewInt(30 * 86400), // 30d
			MaxValidatorCount:      big.NewInt(50),
			MaxReadyCandidateCount: big.NewInt(3),
			KefAddress:             common.HexToAddress("0x1111"),
			KifAddress:             common.HexToAddress("0x2222"),
			KpfAddress:             common.HexToAddress("0x3333"),
		},
	}
}

// MakeTestBlsKey generates a random BLS key pair for testing.
func MakeTestBlsKey() (priv, pub, pop []byte) {
	ikm := make([]byte, 32)
	rand.Read(ikm)

	sk, _ := bls.GenerateKey(ikm)
	pk := sk.PublicKey()
	sig := bls.PopProve(sk)

	priv = sk.Marshal()
	pub = pk.Marshal()
	pop = sig.Marshal()
	if len(priv) != 32 || len(pub) != 48 || len(pop) != 96 {
		panic("bad bls key")
	}
	return
}
