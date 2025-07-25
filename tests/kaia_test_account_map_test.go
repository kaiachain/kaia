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
	"github.com/kaiachain/kaia/work/builder"
)

// //////////////////////////////////////////////////////////////////////////////
// AddressBalanceMap
// //////////////////////////////////////////////////////////////////////////////
type AccountInfo struct {
	balance *big.Int
	nonce   uint64
}

type AccountMap struct {
	m          map[common.Address]*AccountInfo
	rewardbase common.Address
}

func NewAccountMap() *AccountMap {
	return &AccountMap{
		m: make(map[common.Address]*AccountInfo),
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

	return nil
}

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
		if ok {
			feeByFeePayer, feeBySender := types.CalcFeeWithRatio(feeRatio, fee)
			a.SubBalance(feePayer, feeByFeePayer)
			a.SubBalance(from, feeBySender)
		} else {
			a.SubBalance(feePayer, fee)
		}

		a.AddBalance(a.rewardbase, fee)

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
