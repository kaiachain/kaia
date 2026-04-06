// Modifications Copyright 2024 The Kaia Authors
// Copyright 2018 The klaytn Authors
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
	"errors"
	"fmt"
	"math/big"
	"slices"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/work/builder"
)

// //////////////////////////////////////////////////////////////////////////////
// AddressBalanceMap
// //////////////////////////////////////////////////////////////////////////////
type AccountInfo struct {
	balance *big.Int
	nonce   uint64
}

// txFeeInfo caches per-transaction fee data computed in Update(), so that
// AdjustFeesByReceipts can reconcile with actual gas used from receipts.
type txFeeInfo struct {
	estimatedFee *big.Int
	from         common.Address
	feePayer     common.Address
	feeRatio     types.FeeRatio
	hasFeeRatio  bool
}

type AccountMap struct {
	m           map[common.Address]*AccountInfo
	rewardbase  common.Address
	chainConfig *params.ChainConfig
	txFeeCache  map[common.Hash]*txFeeInfo // tx hash → fee info from Update()
}

func NewAccountMap() *AccountMap {
	return &AccountMap{
		m:          make(map[common.Address]*AccountInfo),
		txFeeCache: make(map[common.Hash]*txFeeInfo),
	}
}

func (a *AccountMap) AddBalance(addr common.Address, v *big.Int) {
	if acc, ok := a.m[addr]; ok {
		acc.balance.Add(acc.balance, v)
	} else {
		// create an account
		a.Set(addr, v, 0)
	}
}

func (a *AccountMap) SubBalance(addr common.Address, v *big.Int) error {
	if acc, ok := a.m[addr]; ok {
		acc.balance.Sub(acc.balance, v)
	} else {
		return fmt.Errorf("trying to subtract balance from an uninitiailzed address (%s)", addr.Hex())
	}

	return nil
}

func (a *AccountMap) GetNonce(addr common.Address) uint64 {
	if acc, ok := a.m[addr]; ok {
		return acc.nonce
	}
	// 'StateDB.GetNonce' returns 0 when the address doesn't exist
	return 0
}

func (a *AccountMap) IncNonce(addr common.Address) {
	if acc, ok := a.m[addr]; ok {
		acc.nonce++
	}
}

func (a *AccountMap) Set(addr common.Address, v *big.Int, nonce uint64) {
	a.m[addr] = &AccountInfo{new(big.Int).Set(v), nonce}
}

func (a *AccountMap) Initialize(bcdata *BCData) error {
	statedb, err := bcdata.bc.State()
	if err != nil {
		return err
	}

	// NOTE-Kaia-Issue973 Developing Kaia token economy
	// Add predefined accounts related to reward mechanism
	rewardContractAddr := common.HexToAddress("0x0000000000000000000000000000000000000441")
	kefContractAddr := common.HexToAddress("0x0000000000000000000000000000000000000442")
	kifContractAddr := common.HexToAddress("0x0000000000000000000000000000000000000443")
	addrs := append(bcdata.addrs, &rewardContractAddr, &kefContractAddr, &kifContractAddr)

	for _, addr := range addrs {
		a.Set(*addr, statedb.GetBalance(*addr), statedb.GetNonce(*addr))
	}

	a.rewardbase = *bcdata.addrs[0]
	a.chainConfig = bcdata.bc.Config()

	return nil
}

// Update simulates the balance and nonce changes for txs as if they were applied in the block
// at currentBlockNumber. currentBlockNumber must be the mined block number, not the parent.
func (a *AccountMap) Update(txs types.Transactions, txHashesExpectedFail []common.Hash, txBundlingModules []kaiax.TxBundlingModule, signer types.Signer, picker types.AccountKeyPicker, currentBlockNumber uint64) error {
	modules := make([]builder.TxBundlingModule, len(txBundlingModules)) // TODO-Kaia: Remove this cast.
	for i, module := range txBundlingModules {
		modules[i] = module.(builder.TxBundlingModule)
	}
	incorporatedTxs, _ := builder.ExtractBundlesAndIncorporate(txs, modules)
	for _, txOrGen := range incorporatedTxs {
		// To simulate tx, the nonce given to generate is set to zero.
		// This does not affect subsequent operations on the AccountMap state.
		tx, err := txOrGen.GetTx(0)
		if err != nil {
			return err
		}
		if slices.Contains(txHashesExpectedFail, tx.Hash()) {
			continue
		}
		to := tx.To()
		v := tx.Value()

		gasFrom, err := tx.ValidateSender(signer, picker, currentBlockNumber)
		if err != nil {
			return err
		}
		from := tx.ValidatedSender()

		gasFeePayer := uint64(0)
		feePayer := from
		if tx.IsFeeDelegatedTransaction() {
			gasFeePayer, err = tx.ValidateFeePayer(signer, picker, currentBlockNumber)
			if err != nil {
				return err
			}
			feePayer = tx.ValidatedFeePayer()
		}
		if to == nil {
			nonce := a.GetNonce(from)
			addr := crypto.CreateAddress(from, nonce)
			to = &addr
		}

		a.AddBalance(*to, v)
		a.SubBalance(from, v)

		// TODO-Kaia: This gas fee calculation is correct only if the transaction is a value transfer transaction.
		// Calculate the correct transaction fee by checking the corresponding receipt.
		intrinsicGas, err := tx.IntrinsicGas(currentBlockNumber)
		if err != nil {
			return err
		}

		intrinsicGas += gasFrom + gasFeePayer

		fee := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(intrinsicGas))
		feeRatio, ok := tx.FeeRatio()
		a.txFeeCache[tx.Hash()] = &txFeeInfo{
			estimatedFee: new(big.Int).Set(fee),
			from:         from,
			feePayer:     feePayer,
			feeRatio:     feeRatio,
			hasFeeRatio:  ok,
		}
		if ok {
			feeByFeePayer, feeBySender := types.CalcFeeWithRatio(feeRatio, fee)
			a.SubBalance(feePayer, feeByFeePayer)
			a.SubBalance(from, feeBySender)
		} else {
			a.SubBalance(feePayer, fee)
		}

		// In non-deferred Magma mode, only half the fee is distributed to the rewardbase;
		// the other half is burned. Mirror that logic here.
		// currentBlockNumber is the mined block number (invariant upheld by all callers).
		rewardFee := new(big.Int).Set(fee)
		if a.chainConfig != nil {
			rules := a.chainConfig.Rules(new(big.Int).SetUint64(currentBlockNumber))
			isNonDeferred := a.chainConfig.Governance == nil || !a.chainConfig.Governance.DeferredTxFee()
			if rules.IsMagma && isNonDeferred {
				rewardFee.Div(rewardFee, big.NewInt(2))
			}
		}
		a.AddBalance(a.rewardbase, rewardFee)

		a.IncNonce(from)

		if tx.IsEthereumTransaction() && tx.To() == nil {
			a.IncNonce(*to)
		}

		if tx.Type() == types.TxTypeSmartContractDeploy || tx.Type() == types.TxTypeFeeDelegatedSmartContractDeploy || tx.Type() == types.TxTypeFeeDelegatedSmartContractDeployWithRatio {
			a.IncNonce(*to)
		}
	}

	return nil
}

// AdjustFeesByReceipts corrects the accountMap after mining by reconciling the estimated
// intrinsic gas fee (used in Update) against the actual gas used recorded in receipts.
// This is necessary because Update() uses IntrinsicGas as an approximation, which excludes
// execution overhead (e.g., constructor gas, code storage cost for contract deployments).
// txs must be b.Transactions() and receipts must be the corresponding block receipts.
func (a *AccountMap) AdjustFeesByReceipts(txs types.Transactions, receipts types.Receipts, blockNum uint64) error {
	if a.chainConfig == nil || len(a.txFeeCache) == 0 {
		return nil
	}

	blockNumBig := new(big.Int).SetUint64(blockNum)
	rules := a.chainConfig.Rules(blockNumBig)
	isNonDeferred := a.chainConfig.Governance == nil || !a.chainConfig.Governance.DeferredTxFee()

	for i, tx := range txs {
		cached, ok := a.txFeeCache[tx.Hash()]
		if !ok {
			continue // tx was skipped in Update (txHashesExpectedFail) or not tracked
		}

		actualFee := new(big.Int).Mul(tx.GasPrice(), new(big.Int).SetUint64(receipts[i].GasUsed))
		delta := new(big.Int).Sub(actualFee, cached.estimatedFee)
		if delta.Sign() == 0 {
			continue
		}

		// delta > 0: actual > estimated (e.g. constructor or execution gas not in IntrinsicGas)
		//   → deduct delta from payer, add delta/2 to rewardbase
		// delta < 0: actual < estimated (e.g. EIP-7702 authorization refund > execution gas)
		//   → refund |delta| to payer, subtract |delta|/2 from rewardbase
		// SubBalance with a negative value adds to balance; AddBalance with negative subtracts.
		if cached.hasFeeRatio {
			deltaByFeePayer, deltaBySender := types.CalcFeeWithRatio(cached.feeRatio, delta)
			a.SubBalance(cached.feePayer, deltaByFeePayer)
			a.SubBalance(cached.from, deltaBySender)
		} else {
			a.SubBalance(cached.feePayer, delta)
		}

		rewardDelta := new(big.Int).Set(delta)
		if rules.IsMagma && isNonDeferred {
			rewardDelta.Div(rewardDelta, big.NewInt(2))
		}
		a.AddBalance(a.rewardbase, rewardDelta)
	}

	return nil
}

func (a *AccountMap) Verify(statedb *state.StateDB) error {
	for addr, acc := range a.m {
		if acc.nonce != statedb.GetNonce(addr) {
			return errors.New(fmt.Sprintf("[%s] nonce is different!! statedb(%d) != accountMap(%d).\n",
				addr.Hex(), statedb.GetNonce(addr), acc.nonce))
		}

		if acc.balance.Cmp(statedb.GetBalance(addr)) != 0 {
			return errors.New(fmt.Sprintf("[%s] balance is different!! statedb(%s) != accountMap(%s).\n",
				addr.Hex(), statedb.GetBalance(addr).String(), acc.balance.String()))
		}
	}

	return nil
}
