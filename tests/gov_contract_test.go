// Modifications Copyright 2024 The Kaia Authors
// Copyright 2022 The klaytn Authors
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
	"strings"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	govcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/gov"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	gov_impl "github.com/kaiachain/kaia/kaiax/gov/impl"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/node/cn"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGovernance_GovModule tests GovModule and ContractGovModule
func TestGovernance_GovModule(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)

	config := params.MainnetChainConfig.Copy()
	config.IstanbulCompatibleBlock = new(big.Int).SetUint64(0)
	config.LondonCompatibleBlock = new(big.Int).SetUint64(0)
	config.EthTxTypeCompatibleBlock = new(big.Int).SetUint64(0)
	config.MagmaCompatibleBlock = new(big.Int).SetUint64(0)
	config.KoreCompatibleBlock = new(big.Int).SetUint64(0)

	config.Istanbul.Epoch = 2
	config.Istanbul.SubGroupSize = 1
	config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)
	config.Governance.Reward.MintingAmount = new(big.Int).Mul(big.NewInt(9000000000000000000), big.NewInt(params.KAIA))
	config.Governance.Reward.Kip82Ratio = params.DefaultKip82Ratio

	config.Governance.GovParamContract = common.Address{}
	config.Governance.GovernanceMode = "none"

	fullNode, node, validator, chainId, workspace := newBlockchain(t, config, nil)
	defer os.RemoveAll(workspace)
	defer fullNode.Stop()

	var (
		chain        = node.BlockChain().(*blockchain.BlockChain)
		owner        = validator
		contractAddr = crypto.CreateAddress(owner.Addr, owner.Nonce)

		paramName  = "istanbul.committeesize"
		oldVal     = config.Istanbul.SubGroupSize
		newVal     = uint64(22)
		paramBytes = []byte{22}

		govBlock  uint64 // Before vote: 0, After vote: the governance block
		stopBlock uint64 // Before govBlock is set: 0, After: the block to stop receiving new blocks
	)

	// Here we are running (tx sender) and (param reader) in parallel.
	// This is to check that param reader (mixed engine) works in such situations:
	// (a) contract engine disabled
	// (b) contract engine enabled (via vote)

	govModule := node.GovModule().(*gov_impl.GovModule)
	chain.Engine().(consensus.Istanbul).RegisterConsensusModule(govModule)
	chain.RegisterExecutionModule(govModule)
	node.Miner().RegisterExecutionModule(govModule)

	// Run tx sender thread
	go func() {
		deployGovParamTx_constructor(t, node, owner, chainId)

		// Give some time for txpool to recognize the contract, because otherwise
		// the txpool may reject the setParam tx with 'not a program account'
		time.Sleep(2 * time.Second)

		deployGovParamTx_setParamIn(t, node, owner, chainId, contractAddr, paramName, paramBytes)

		vote := headergov.NewVoteData(validator.Addr, string(gov.GovernanceGovParamContract), contractAddr)
		require.NotNil(t, vote)
		govModule.Hgm.PushMyVotes(vote)
	}()

	// Validate current params from mixedEngine.CurrentParams() & contractEngine.CurrentParams(),
	// alongside block processing.
	// At block #N, CurrentParams() returns the parameters to be used when building
	// block #N+1 (i.e. pending block).
	chainEventCh := make(chan blockchain.ChainEvent)
	subscription := chain.SubscribeChainEvent(chainEventCh)
	defer subscription.Unsubscribe()

	// 1. wait until header.Governance block + 5
	for {
		ev := <-chainEventCh
		time.Sleep(100 * time.Millisecond) // wait for tx sender thread to set deployBlock, etc.

		num := ev.Block.Number().Uint64()
		if govBytes := ev.Block.Header().Governance; len(govBytes) > 0 {
			govBlock = num
			// stopBlock is the epoch block, so we stop when receiving it
			// otherwise, GetParamSet(stopBlock) may fail
			stopBlock = govBlock + 5
			stopBlock = stopBlock - (stopBlock % config.Istanbul.Epoch)
			t.Logf("Governance at block=%2d, stopBlock=%2d, gov=%v", num, stopBlock, hexutil.Encode(govBytes))
		}

		if stopBlock != 0 && num >= stopBlock {
			break
		}

		if num >= 60 {
			t.Fatal("test taking too long; something must be wrong")
		}
	}

	// 2. test GetParamSet():  Validate historic params from both Engines
	for num := uint64(0); num < stopBlock; num++ {
		govVal := govModule.GetParamSet(num).CommitteeSize
		contractVal := govModule.Cgm.GetParamSet(num).CommitteeSize

		if num <= govBlock+1 { // ContractEngine disabled
			assert.Equal(t, oldVal, govVal)
			assert.Equal(t, gov.Params[gov.IstanbulCommitteeSize].DefaultValue, contractVal)
		} else { // ContractEngine enabled
			assert.Equal(t, newVal, govVal)
			assert.Equal(t, newVal, contractVal)
		}
	}
}

func deployGovParamTx_constructor(t *testing.T, node *cn.CN, owner *TestAccountType, chainId *big.Int,
) (uint64, common.Address, *types.Transaction) {
	var (
		// Deploy contract: constructor(address _owner)
		contractAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))
		contractBin    = govcontract.GovParamBin
		ctorArgs, _    = contractAbi.Pack("")
		code           = contractBin + hex.EncodeToString(ctorArgs)
	)

	// Deploy contract
	tx, addr := deployContractDeployTx(t, node.TxPool(), chainId, owner, code)

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("GovParam deployed at block=%2d, addr=%s", num, addr.Hex())

	return num, addr, tx
}

func deployGovParamTx_setParamIn(t *testing.T, node *cn.CN, owner *TestAccountType, chainId *big.Int,
	contractAddr common.Address, name string, value []byte,
) (uint64, *types.Transaction) {
	var (
		contractAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))
		callArgs, _    = contractAbi.Pack("setParamIn", name, true, value, big.NewInt(1))
		data           = common.ToHex(callArgs)
	)

	tx := deployContractExecutionTx(t, node.TxPool(), chainId, owner, contractAddr, data)

	chain := node.BlockChain().(*blockchain.BlockChain)
	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status, "setParamIn failed")

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("GovParam.setParamIn executed at block=%2d", num)
	return num, tx
}

func deployGovParamTx_batchSetParamIn(t *testing.T, node *cn.CN, owner *TestAccountType, chainId *big.Int,
	contractAddr common.Address, bytesMap map[string][]byte,
) []*types.Transaction {
	var (
		chain          = node.BlockChain().(*blockchain.BlockChain)
		beginBlock     = chain.CurrentHeader().Number.Uint64()
		contractAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))
		txs            = []*types.Transaction{}
	)

	// Send all setParamIn() calls at once
	for name, value := range bytesMap {
		callArgs, _ := contractAbi.Pack("setParamIn", name, true, value, big.NewInt(1))
		data := common.ToHex(callArgs)
		tx := deployContractExecutionTx(t, node.TxPool(), chainId, owner, contractAddr, data)
		txs = append(txs, tx)
	}

	// Wait for all txs
	for _, tx := range txs {
		receipt := waitReceipt(chain, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status, "batchSetParamIn failed")
	}
	num := chain.CurrentHeader().Number.Uint64()
	t.Logf("GovParam.setParamIn executed %d times between blocks=%2d,%2d", len(txs), beginBlock, num)
	return txs
}

// Kaia node only decodes the byte-array param values (refer to params/governance_paramset.go).
// Encoding is the job of transaction senders (i.e. clients and dApps).
// This is a reference implementation of such encoder.
func chainConfigToBytesMap(t *testing.T, config *params.ChainConfig) map[string][]byte {
	pset, err := params.NewGovParamSetChainConfig(config)
	require.Nil(t, err)
	strMap := pset.StrMap()

	bytesMap := map[string][]byte{}
	for name, value := range strMap {
		switch value.(type) {
		case string:
			bytesMap[name] = []byte(value.(string))
		case common.Address:
			bytesMap[name] = value.(common.Address).Bytes()
		case uint64:
			bytesMap[name] = new(big.Int).SetUint64(value.(uint64)).Bytes()
		case bool:
			if value.(bool) == true {
				bytesMap[name] = []byte{0x01}
			} else {
				bytesMap[name] = []byte{0x00}
			}
		}
	}

	// Check that bytesMap is correct just in case
	qset, err := params.NewGovParamSetBytesMap(bytesMap)
	require.Nil(t, err)
	require.Equal(t, pset.StrMap(), qset.StrMap())
	return bytesMap
}
