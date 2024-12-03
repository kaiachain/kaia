// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package tests

import (
	"encoding/hex"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func testEIP3541ChainConfig(forkNum *big.Int) *params.ChainConfig {
	config := params.MainnetChainConfig.Copy()
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.ShanghaiCompatibleBlock = big.NewInt(0)
	config.CancunCompatibleBlock = big.NewInt(0)
	config.RandaoCompatibleBlock = nil
	config.KaiaCompatibleBlock = big.NewInt(0)
	config.PragueCompatibleBlock = forkNum

	config.Istanbul.SubGroupSize = 1
	config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)

	return config
}

func TestEIP3541(t *testing.T) {
    log.EnableLogForTest(log.LvlError, log.LvlInfo)

    testCases := []struct {
        name           string
        pragueBlock    *big.Int
        expectDeployed bool
    }{
        {"BeforePrague", nil, true},
        {"AfterPrague", big.NewInt(0), false},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Configure chain and create node
            config := testEIP3541ChainConfig(tc.pragueBlock)
            fullNode, node, validator, chainId, workspace := newBlockchain(t, config, nil)
            defer func() {
                fullNode.Stop()
                os.RemoveAll(workspace)
            }()

            var (
				// Sample creation code
                inputBytes, _  = hex.DecodeString("610017600081600b8239f3ef01000000000000000000000000000000000000001000")
				// Expected code hash if deployed
                expectedCode, _ = hex.DecodeString("ef01000000000000000000000000000000000000001000")
                optsOwner      = bind.NewKeyedTransactor(validator.Keys[0])
                transactor     = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
                nonce, _       = transactor.PendingNonceAt(optsOwner.Context, optsOwner.From)
                rawTx         = types.NewContractCreation(nonce, big.NewInt(0), 1000000, big.NewInt(25000000000), inputBytes)
                signer        = types.LatestSignerForChainID(chainId)
            )

            // Deploy contract
            tx, err := optsOwner.Signer(signer, optsOwner.From, rawTx)
            if err != nil {
                t.Fatal(err)
            }
            err = transactor.SendTransaction(optsOwner.Context, tx)
            if err != nil {
                t.Fatal(err)
            }
            receipt := waitReceipt(node.BlockChain().(*blockchain.BlockChain), tx.Hash())
            if receipt == nil {
                t.Fatal("timeout")
            }
            time.Sleep(1 * time.Second)

            deployedAddress := crypto.CreateAddress(optsOwner.From, nonce)

            // Check receipt status
            assert.Equal(t, receipt.Status, types.ReceiptStatusErrDefault)

            // Verify state based on whether Prague is active
            state, _ := node.BlockChain().State()
            assert.Equal(t, tc.expectDeployed, state.IsContractAvailable(deployedAddress))
            
            if tc.expectDeployed {
                assert.Equal(t, expectedCode, state.GetCode(deployedAddress))
                assert.Equal(t, crypto.Keccak256Hash(expectedCode), state.GetCodeHash(deployedAddress))
            } else {
                assert.Equal(t, []byte(nil), state.GetCode(deployedAddress))
                assert.Equal(t, crypto.Keccak256Hash(nil), state.GetCodeHash(deployedAddress))
            }
        })
    }
}
