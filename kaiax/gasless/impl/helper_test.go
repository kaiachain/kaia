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

package impl

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/builder"
	gasless_cfg "github.com/kaiachain/kaia/kaiax/gasless/config"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

var (
	testChainConfig = &params.ChainConfig{
		ChainID: big.NewInt(1),
		Governance: &params.GovernanceConfig{
			Reward: &params.RewardConfig{
				UseGiniCoeff:          true,
				StakingUpdateInterval: 86400,
			},
			KIP71: params.GetDefaultKIP71Config(),
		},
		Gasless: &gasless_cfg.ChainConfig{
			AllowedTokens: nil,
			IsDisabled:    false,
		},
		IstanbulCompatibleBlock:  big.NewInt(0),
		LondonCompatibleBlock:    big.NewInt(0),
		EthTxTypeCompatibleBlock: big.NewInt(0),
		MagmaCompatibleBlock:     big.NewInt(0),
		KoreCompatibleBlock:      big.NewInt(0),
		ShanghaiCompatibleBlock:  big.NewInt(0),
		CancunCompatibleBlock:    big.NewInt(0),
		KaiaCompatibleBlock:      big.NewInt(0),
		PragueCompatibleBlock:    big.NewInt(0),
		Kip103CompatibleBlock:    big.NewInt(0),
		Kip160CompatibleBlock:    big.NewInt(0),
		RandaoCompatibleBlock:    big.NewInt(0),
	}

	// interface DummyGSR {
	// 	function addToken(address token) public
	// 	function removeToken(address token) public
	// 	function getSupportedTokens() external view returns (address[] memory)
	// }
	dummyGSRCode    = "0x608060405234801561000f575f80fd5b506004361061003f575f3560e01c80635fa7b58414610043578063d3c7c2c71461005f578063d48bfca71461007d575b5f80fd5b61005d6004803603810190610058919061035a565b610099565b005b610067610210565b604051610074919061043c565b60405180910390f35b6100976004803603810190610092919061035a565b61029a565b005b5f5b5f8054905081101561020c578173ffffffffffffffffffffffffffffffffffffffff165f82815481106100d1576100d061045c565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036101ff575f60015f8054905061012691906104bf565b815481106101375761013661045c565b5b905f5260205f20015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff165f82815481106101725761017161045c565b5b905f5260205f20015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505f8054806101c8576101c76104f2565b5b600190038181905f5260205f20015f6101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055905561020c565b808060010191505061009b565b5050565b60605f80548060200260200160405190810160405280929190818152602001828054801561029057602002820191905f5260205f20905b815f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610247575b5050505050905090565b5f81908060018154018082558091505060019003905f5260205f20015f9091909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b5f80fd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61032982610300565b9050919050565b6103398161031f565b8114610343575f80fd5b50565b5f8135905061035481610330565b92915050565b5f6020828403121561036f5761036e6102fc565b5b5f61037c84828501610346565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b6103b78161031f565b82525050565b5f6103c883836103ae565b60208301905092915050565b5f602082019050919050565b5f6103ea82610385565b6103f4818561038f565b93506103ff8361039f565b805f5b8381101561042f57815161041688826103bd565b9750610421836103d4565b925050600181019050610402565b5085935050505092915050565b5f6020820190508181035f83015261045481846103e0565b905092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f819050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6104c982610489565b91506104d483610489565b92508282039050818111156104ec576104eb610492565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffdfea264697066735822122000038ab02a3505bfb41836b5d6e96fe339af9309293b26d4edefd7b4a1e90bca64736f6c634300081a0033"
	dummyGSRAddress = common.HexToAddress("0x1234")

	dummyTokenAddress1 = common.HexToAddress("0xabcd")
	dummyTokenAddress2 = common.HexToAddress("0xbcde")
	dummyTokenAddress3 = common.HexToAddress("0xcdef")
)

type testBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{Number: big.NewInt(0)}, nil, nil)
}

func (bc *testBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) State() (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

type dummyGovModule struct {
	chainConfig *params.ChainConfig
}

func (d *dummyGovModule) GetParamSet(blockNum uint64) gov.ParamSet {
	return gov.ParamSet{UnitPrice: d.chainConfig.UnitPrice}
}

type AccountKeyPickerForTest struct {
	AddrKeyMap map[common.Address]accountkey.AccountKey
}

func (a *AccountKeyPickerForTest) GetKey(addr common.Address) accountkey.AccountKey {
	return a.AddrKeyMap[addr]
}

func (a *AccountKeyPickerForTest) SetKey(addr common.Address, key accountkey.AccountKey) {
	a.AddrKeyMap[addr] = key
}

func (a *AccountKeyPickerForTest) Exist(addr common.Address) bool {
	return a.AddrKeyMap[addr] != nil
}

func makeTx(t *testing.T, privKey *ecdsa.PrivateKey, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	if privKey == nil {
		var err error
		privKey, err = crypto.GenerateKey()
		require.NoError(t, err)
	}
	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	p := &AccountKeyPickerForTest{
		AddrKeyMap: make(map[common.Address]accountkey.AccountKey),
	}
	p.SetKey(addr, accountkey.NewAccountKeyLegacy())

	signer := types.LatestSignerForChainID(big.NewInt(1))
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
	tx, err := types.SignTx(tx, signer, privKey)
	require.NoError(t, err)

	return tx
}

func makeApproveTx(t *testing.T, privKey *ecdsa.PrivateKey, nonce uint64, approveArgs ApproveArgs) *types.Transaction {
	var err error
	if privKey == nil {
		privKey, err = crypto.GenerateKey()
		require.NoError(t, err)
	}

	data := append([]byte{}, common.Hex2Bytes("095ea7b3")...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(approveArgs.Spender.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(approveArgs.Amount.Bytes()), 32)...)
	approveTx := makeTx(t, privKey, nonce, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(1), data)

	return approveTx
}

func makeSwapTx(t *testing.T, privKey *ecdsa.PrivateKey, nonce uint64, swapArgs SwapArgs) *types.Transaction {
	var err error
	if privKey == nil {
		privKey, err = crypto.GenerateKey()
		require.NoError(t, err)
	}

	data := append([]byte{}, common.Hex2Bytes("43bab9f7")...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.Token.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.AmountIn.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.MinAmountOut.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.AmountRepay.Bytes()), 32)...)
	swapTx := makeTx(t, privKey, nonce, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1), data)

	return swapTx
}

func flattenPoolTxs(structured map[common.Address]types.Transactions) map[common.Hash]bool {
	flattened := map[common.Hash]bool{}
	for _, txs := range structured {
		for _, tx := range txs {
			flattened[tx.Hash()] = true
		}
	}
	return flattened
}

func flattenBundleTxs(txs []interface{}) ([]common.Hash, error) {
	nodeNonce := uint64(0)
	hashes := []common.Hash{}
	for _, txi := range txs {
		var tx *types.Transaction
		var err error
		if genLendTx, ok := txi.(builder.TxGenerator); ok {
			tx, err = genLendTx.Generate(nodeNonce)
			if err != nil {
				return nil, err
			}
			nodeNonce += 1
		} else if tx, ok = txi.(*types.Transaction); ok {
		} else {
			err = errors.New("unsupported bundle tx")
		}
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, tx.Hash())
	}
	return hashes, nil
}

type testTxPool struct {
	statedb *state.StateDB
}

func (pool *testTxPool) GetCurrentState() *state.StateDB {
	return pool.statedb
}

func testAllocStorage() blockchain.GenesisAlloc {
	allocStorage := system.AllocRegistry(&params.RegistryConfig{
		Records: map[string]common.Address{
			GaslessSwapRouterName: dummyGSRAddress,
		},
		Owner: common.HexToAddress("0xffff"),
	})
	alloc := blockchain.GenesisAlloc{
		system.RegistryAddr: {
			Code:    system.RegistryMockCode,
			Balance: big.NewInt(0),
			Storage: allocStorage,
		},
		dummyGSRAddress: {
			Code:    hexutil.MustDecode(dummyGSRCode),
			Balance: big.NewInt(0),
			Storage: map[common.Hash]common.Hash{
				// key: slot of _supportedTokens, value: length of _supportedTokens
				common.HexToHash("0x0"): common.HexToHash("0x3"),
				// key: slot of _supportedTokens elements (keccak(0x0) + n), value: element values
				common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"): dummyTokenAddress1.Hash(),
				common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e564"): dummyTokenAddress2.Hash(),
				common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e565"): dummyTokenAddress3.Hash(),
			},
		},
	}
	return alloc
}
