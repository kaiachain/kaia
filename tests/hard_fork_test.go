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
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/common/profile"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulBackend "github.com/kaiachain/kaia/consensus/istanbul/backend"
	"github.com/kaiachain/kaia/crypto"
	gov_impl "github.com/kaiachain/kaia/kaiax/gov/impl"
	reward_impl "github.com/kaiachain/kaia/kaiax/reward/impl"
	staking_impl "github.com/kaiachain/kaia/kaiax/staking/impl"
	valset_impl "github.com/kaiachain/kaia/kaiax/valset/impl"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHardForkBlock tests whether the change incurs a hard fork or not.
// genesis.json, b1.rlp, and b2.rlp has raw data of genesis, and consecutive two blocks after the genesis block.
// If anything is failed, it can be considered that a hard fork occurs.
func TestHardForkBlock(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	var genesis blockchain.Genesis

	// If you uncomment the below, you can find this test failed with an error "!!!!!HARD FORK DETECTED!!!!!"
	//fork.UpdateHardForkConfig(&fork.HardForkConfig{
	//})

	// If you print out b1.rlp and b2.rlp, uncomment below.
	// `genBlocks` could be failed sometimes depends on the order of transaction in a block. Just try again.
	// genBlocks(t)
	// return

	// load raw data from files.
	genesisJson, err := os.ReadFile("genesis.json")
	require.Equal(t, nil, err)
	rawb1, err := os.ReadFile("b1.rlp")
	require.Equal(t, nil, err)
	rawb2, err := os.ReadFile("b2.rlp")
	require.Equal(t, nil, err)

	err = json.Unmarshal([]byte(genesisJson), &genesis)
	require.Equal(t, nil, err)

	genesisKey, err := crypto.HexToECDSA("42eb1412d77987043716f425964b1c8d4c27ce9fb3e9a5b9ab243bc9882fe731")
	require.Equal(t, nil, err)

	var genesisAddr common.Address
	for addr := range genesis.Alloc {
		genesisAddr = addr
		break
	}

	dir := "chaindata-hardfork"
	os.RemoveAll(dir)

	chainDb := NewDatabase(dir, database.LevelDB)
	defer func() {
		os.RemoveAll(dir)
	}()

	chainConfig, _, err := blockchain.SetupGenesisBlock(chainDb, &genesis, params.UnusedNetworkId, false, false)
	chainConfig.Istanbul = &params.IstanbulConfig{
		Epoch: 100,
	}
	require.Nil(t, err)

	govModule := gov_impl.NewGovModule()
	engine := istanbulBackend.New(&istanbulBackend.BackendOpts{
		IstanbulConfig: istanbul.DefaultConfig,
		Rewardbase:     genesisAddr,
		PrivateKey:     genesisKey,
		DB:             chainDb,
		GovModule:      govModule,
		NodeType:       common.CONSENSUSNODE,
	})
	chain, err := blockchain.NewBlockChain(chainDb, nil, chainConfig, engine, vm.Config{})
	require.NoError(t, err)

	mStaking := staking_impl.NewStakingModule()
	mReward := reward_impl.NewRewardModule()
	mValset := valset_impl.NewValsetModule()
	err = errors.Join(
		govModule.Init(&gov_impl.InitOpts{
			ChainConfig: chainConfig,
			ChainKv:     chainDb.GetMiscDB(),
			Chain:       chain,
			Valset:      mValset,
		}),
		mReward.Init(&reward_impl.InitOpts{
			ChainConfig:   chainConfig,
			Chain:         chain,
			GovModule:     govModule,
			StakingModule: mStaking, // Not used in "Simple" istanbul policy
		}),
		mValset.Init(&valset_impl.InitOpts{
			ChainKv:       chainDb.GetMiscDB(),
			Chain:         chain,
			GovModule:     govModule,
			StakingModule: mStaking,
		}),
	)
	require.NoError(t, err)
	engine.RegisterConsensusModule(mReward)
	engine.RegisterKaiaxModules(govModule, mStaking, mValset)
	mValset.Start()

	r1, err := hexutil.Decode(string(rawb1))
	require.NoError(t, err)
	r2, err := hexutil.Decode(string(rawb2))
	require.NoError(t, err)
	rawBlocks := [...][]byte{r1, r2}

	var blocks types.Blocks
	for _, raw := range rawBlocks {
		var blk types.Block

		err := rlp.DecodeBytes(raw, &blk)
		require.Equal(t, nil, err)

		blocks = append(blocks, &blk)
	}

	idx, err := chain.InsertChain(blocks)
	require.Equalf(t, nil, err, "!!!!!HARD FORK DETECTED!!!!!")
	require.Equal(t, 0, idx)
}

// genBlock generates two blocks including transactions utilizing all transaction types and account types.
func genBlocks(t *testing.T) {
	testFunctions := []struct {
		Name  string
		genTx genTransaction
	}{
		{"LegacyTransaction", genLegacyTransaction},
		{"ValueTransfer", genValueTransfer},
		{"ValueTransferWithMemo", genValueTransferWithMemo},
		{"AccountUpdate", genAccountUpdateIdem},
		{"SmartContractExecution", genSmartContractExecution},
		{"Cancel", genCancel},
		{"ChainDataAnchoring", genChainDataAnchoring},
		{"FeeDelegatedValueTransfer", genFeeDelegatedValueTransfer},
		{"FeeDelegatedValueTransferWithMemo", genFeeDelegatedValueTransferWithMemo},
		{"FeeDelegatedAccountUpdate", genFeeDelegatedAccountUpdateIdem},
		{"FeeDelegatedSmartContractExecution", genFeeDelegatedSmartContractExecution},
		{"FeeDelegatedCancel", genFeeDelegatedCancel},
		{"FeeDelegatedWithRatioValueTransfer", genFeeDelegatedWithRatioValueTransfer},
		{"FeeDelegatedWithRatioValueTransferWithMemo", genFeeDelegatedWithRatioValueTransferWithMemo},
		{"FeeDelegatedWithRatioAccountUpdate", genFeeDelegatedWithRatioAccountUpdateIdem},
		{"FeeDelegatedWithRatioSmartContractExecution", genFeeDelegatedWithRatioSmartContractExecution},
		{"FeeDelegatedWithRatioCancel", genFeeDelegatedWithRatioCancel},
	}

	accountTypes := []struct {
		Type    string
		account TestAccount
	}{
		{"KaiaLegacy", genKaiaLegacyAccount(t)},
		{"Public", genPublicAccount(t)},
		{"MultiSig", genMultiSigAccount(t)},
		{"RoleBasedWithPublic", genRoleBasedWithPublicAccount(t)},
		{"RoleBasedWithMultiSig", genRoleBasedWithMultiSigAccount(t)},
	}

	log.EnableLogForTest(log.LvlCrit, log.LvlTrace)
	prof := profile.NewProfiler()

	// Initialize blockchain
	start := time.Now()
	bcdata, err := NewBCDataWithForkConfig(6, 4, Forks["Prague"])
	assert.NoError(t, err)
	prof.Profile("main_init_blockchain", time.Now().Sub(start))

	b, err := json.Marshal(bcdata.genesis)
	os.WriteFile("genesis.json", b, 0o755)

	defer bcdata.Shutdown()

	// Initialize address-balance map for verification
	start = time.Now()
	accountMap := NewAccountMap()
	if err := accountMap.Initialize(bcdata); err != nil {
		t.Fatal(err)
	}
	prof.Profile("main_init_accountMap", time.Now().Sub(start))

	// reservoir account
	var reservoir TestAccount
	reservoir = &TestAccountType{
		Addr:  *bcdata.addrs[0],
		Keys:  []*ecdsa.PrivateKey{bcdata.privKeys[0]},
		Nonce: uint64(0),
	}

	signer := types.LatestSignerForChainID(bcdata.bc.Config().ChainID)
	gasPrice := new(big.Int).SetUint64(bcdata.bc.Config().UnitPrice)

	// For smart contract
	contract, err := createAnonymousAccount("ed34b0cf47a0021e9897760f0a904a69260c2f638e0bcc805facb745ec3ff9ab")
	assert.Equal(t, nil, err)

	// Preparing step
	{
		var txs types.Transactions
		// Preparing step. Send KAIA to LegacyAccount.
		{
			amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KAIA))
			tx := types.NewTransaction(reservoir.GetNonce(),
				accountTypes[0].account.GetAddr(), amount, gasLimit, gasPrice, []byte{})

			err := tx.SignWithKeys(signer, reservoir.GetTxKeys())
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			reservoir.AddNonce()
		}

		// Preparing step. Send KAIA to KaiaAcounts.
		for i := 1; i < len(accountTypes); i++ {
			// create an account which account key will be replaced to one of account key types.
			anon, err := createAnonymousAccount(getRandomPrivateKeyString(t))
			assert.Equal(t, nil, err)

			accountTypes[i].account.SetAddr(anon.Addr)

			{
				amount := new(big.Int).Mul(big.NewInt(3000), new(big.Int).SetUint64(params.KAIA))
				tx := types.NewTransaction(reservoir.GetNonce(),
					accountTypes[i].account.GetAddr(), amount, gasLimit, gasPrice, []byte{})

				err := tx.SignWithKeys(signer, reservoir.GetTxKeys())
				assert.Equal(t, nil, err)
				txs = append(txs, tx)

				reservoir.AddNonce()
			}

			// update the account's key
			{
				values := map[types.TxValueKeyType]interface{}{
					types.TxValueKeyNonce:      accountTypes[i].account.GetNonce(),
					types.TxValueKeyFrom:       accountTypes[i].account.GetAddr(),
					types.TxValueKeyGasLimit:   gasLimit,
					types.TxValueKeyGasPrice:   gasPrice,
					types.TxValueKeyAccountKey: accountTypes[i].account.GetAccKey(),
				}
				tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
				assert.Equal(t, nil, err)

				err = tx.SignWithKeys(signer, anon.Keys)
				assert.Equal(t, nil, err)

				txs = append(txs, tx)

				accountTypes[i].account.AddNonce()
			}
		}

		{
			amount := new(big.Int).SetUint64(0)

			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:         reservoir.GetNonce(),
				types.TxValueKeyFrom:          reservoir.GetAddr(),
				types.TxValueKeyTo:            (*common.Address)(nil),
				types.TxValueKeyAmount:        amount,
				types.TxValueKeyGasLimit:      gasLimit,
				types.TxValueKeyGasPrice:      gasPrice,
				types.TxValueKeyHumanReadable: false,
				types.TxValueKeyData:          common.FromHex(code),
				types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.GetTxKeys())
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			contract.Addr = crypto.CreateAddress(reservoir.GetAddr(), reservoir.GetNonce())

			reservoir.AddNonce()
		}
		{
			amount := new(big.Int).SetUint64(0)

			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:         reservoir.GetNonce(),
				types.TxValueKeyFrom:          reservoir.GetAddr(),
				types.TxValueKeyTo:            (*common.Address)(nil),
				types.TxValueKeyAmount:        amount,
				types.TxValueKeyGasLimit:      gasLimit,
				types.TxValueKeyGasPrice:      gasPrice,
				types.TxValueKeyHumanReadable: false,
				types.TxValueKeyData:          common.FromHex(code),
				types.TxValueKeyFeePayer:      reservoir.GetAddr(),
				types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.GetTxKeys())
			assert.Equal(t, nil, err)

			err = tx.SignFeePayerWithKeys(signer, reservoir.GetFeeKeys())
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			reservoir.AddNonce()
		}
		{
			amount := new(big.Int).SetUint64(0)

			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:              reservoir.GetNonce(),
				types.TxValueKeyFrom:               reservoir.GetAddr(),
				types.TxValueKeyTo:                 (*common.Address)(nil),
				types.TxValueKeyAmount:             amount,
				types.TxValueKeyGasLimit:           gasLimit,
				types.TxValueKeyGasPrice:           gasPrice,
				types.TxValueKeyHumanReadable:      false,
				types.TxValueKeyData:               common.FromHex(code),
				types.TxValueKeyFeePayer:           reservoir.GetAddr(),
				types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
				types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.GetTxKeys())
			assert.Equal(t, nil, err)

			err = tx.SignFeePayerWithKeys(signer, reservoir.GetFeeKeys())
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			reservoir.AddNonce()
		}

		// SmartContractDeploy with Nil Recipient.
		{
			amount := new(big.Int).SetUint64(0)

			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:         reservoir.GetNonce(),
				types.TxValueKeyFrom:          reservoir.GetAddr(),
				types.TxValueKeyTo:            (*common.Address)(nil),
				types.TxValueKeyAmount:        amount,
				types.TxValueKeyGasLimit:      gasLimit,
				types.TxValueKeyGasPrice:      gasPrice,
				types.TxValueKeyHumanReadable: false,
				types.TxValueKeyData:          common.FromHex(code),
				types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.GetTxKeys())
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			reservoir.AddNonce()
		}
		{
			amount := new(big.Int).SetUint64(0)

			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:         reservoir.GetNonce(),
				types.TxValueKeyFrom:          reservoir.GetAddr(),
				types.TxValueKeyTo:            (*common.Address)(nil),
				types.TxValueKeyAmount:        amount,
				types.TxValueKeyGasLimit:      gasLimit,
				types.TxValueKeyGasPrice:      gasPrice,
				types.TxValueKeyHumanReadable: false,
				types.TxValueKeyData:          common.FromHex(code),
				types.TxValueKeyFeePayer:      reservoir.GetAddr(),
				types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.GetTxKeys())
			assert.Equal(t, nil, err)

			err = tx.SignFeePayerWithKeys(signer, reservoir.GetFeeKeys())
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			reservoir.AddNonce()
		}
		{
			amount := new(big.Int).SetUint64(0)

			values := map[types.TxValueKeyType]interface{}{
				types.TxValueKeyNonce:              reservoir.GetNonce(),
				types.TxValueKeyFrom:               reservoir.GetAddr(),
				types.TxValueKeyTo:                 (*common.Address)(nil),
				types.TxValueKeyAmount:             amount,
				types.TxValueKeyGasLimit:           gasLimit,
				types.TxValueKeyGasPrice:           gasPrice,
				types.TxValueKeyHumanReadable:      false,
				types.TxValueKeyData:               common.FromHex(code),
				types.TxValueKeyFeePayer:           reservoir.GetAddr(),
				types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
				types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
			}
			tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
			assert.Equal(t, nil, err)

			err = tx.SignWithKeys(signer, reservoir.GetTxKeys())
			assert.Equal(t, nil, err)

			err = tx.SignFeePayerWithKeys(signer, reservoir.GetFeeKeys())
			assert.Equal(t, nil, err)

			txs = append(txs, tx)

			reservoir.AddNonce()
		}

		if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
			t.Fatal(err)
		}
	}

	var txs types.Transactions

	for _, f := range testFunctions {
		for _, sender := range accountTypes {
			toAccount := reservoir

			// LegacyTransaction can be used only with LegacyAccount and KaiaAccount with AccountKeyLegacy.
			if !strings.Contains(sender.Type, "Legacy") && strings.Contains(f.Name, "Legacy") {
				continue
			}

			// Sender can't be a LegacyAccount with AccountUpdate
			if sender.Type == "Legacy" && strings.Contains(f.Name, "AccountUpdate") {
				continue
			}

			gasPriceLocal := gasPrice
			// Set contract's address with SmartContractExecution
			if strings.Contains(f.Name, "SmartContractExecution") {
				toAccount = contract
				gasPriceLocal = big.NewInt(0)
			}

			if !strings.Contains(f.Name, "FeeDelegated") {
				// For NonFeeDelegated Transactions
				tx, _ := f.genTx(t, signer, sender.account, toAccount, nil, gasPriceLocal)
				txs = append(txs, tx)
				sender.account.AddNonce()
			} else {
				// For FeeDelegated(WithRatio) Transactions
				for _, payer := range accountTypes {
					tx, _ := f.genTx(t, signer, sender.account, toAccount, payer.account, gasPriceLocal)
					txs = append(txs, tx)
					sender.account.AddNonce()
				}
			}
		}
	}

	if err := bcdata.GenABlockWithTransactions(accountMap, txs, prof); err != nil {
		t.Fatal(err)
	}

	lastBlock := bcdata.bc.CurrentBlock().NumberU64()
	for i := uint64(0); i <= lastBlock; i++ {
		blk := bcdata.bc.GetBlockByNumber(i)
		b, err := rlp.EncodeToBytes(blk)
		require.Equal(t, nil, err)

		// fmt.Println(blk.String())
		// fmt.Println("encoded===")
		// fmt.Println((hexutil.Bytes)(b))

		filename := fmt.Sprintf("b%d.rlp", i)
		f, err := os.Create(filename)
		require.Equal(t, nil, err)

		_, err = f.WriteString(hexutil.Encode(b))
		require.Equal(t, nil, err)

		err = f.Close()
		require.Equal(t, nil, err)
	}
}

func genAccountUpdateIdem(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForUpdate(from, to, gasPrice, from.GetAccKey(), types.TxTypeAccountUpdate)

	tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetUpdateKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genFeeDelegatedAccountUpdateIdem(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForUpdate(from, to, gasPrice, from.GetAccKey(), types.TxTypeFeeDelegatedAccountUpdate)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetUpdateKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}

func genFeeDelegatedWithRatioAccountUpdateIdem(t *testing.T, signer types.Signer, from TestAccount, to TestAccount, payer TestAccount, gasPrice *big.Int) (*types.Transaction, uint64) {
	values, intrinsic := genMapForUpdate(from, to, gasPrice, from.GetAccKey(), types.TxTypeFeeDelegatedAccountUpdateWithRatio)
	values[types.TxValueKeyFeePayer] = payer.GetAddr()
	values[types.TxValueKeyFeeRatioOfFeePayer] = types.FeeRatio(30)

	tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, values)
	assert.Equal(t, nil, err)

	err = tx.SignWithKeys(signer, from.GetUpdateKeys())
	assert.Equal(t, nil, err)

	err = tx.SignFeePayerWithKeys(signer, payer.GetFeeKeys())
	assert.Equal(t, nil, err)

	return tx, intrinsic
}
