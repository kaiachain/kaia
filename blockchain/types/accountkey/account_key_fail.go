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

package accountkey

import (
	"crypto/ecdsa"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kerrors"
)

// AccountKeyFail is used to prevent smart contract accounts from withdrawing tokens
// from themselves with a public key recovery mechanism.
// Kaia assumes that the only way to take tokens from smart contract account is using
// `transfer()` in the smart contract code.
type AccountKeyFail struct{}

var globalFailKey = &AccountKeyFail{}

// NewAccountKeyFail creates a new AccountKeyFail object.
// Since AccountKeyFail has no attributes, use one global variable for all allocations.
func NewAccountKeyFail() *AccountKeyFail { return globalFailKey }

func (a *AccountKeyFail) Type() AccountKeyType {
	return AccountKeyTypeFail
}

func (a *AccountKeyFail) IsCompositeType() bool {
	return false
}

// This type of account key always return false
func (a *AccountKeyFail) ValidateMember(recoveredKey *ecdsa.PublicKey, from common.Address) bool {
	return false
}

func (a *AccountKeyFail) Equal(b AccountKey) bool {
	// This type of account key always returns false.
	return false
}

func (a *AccountKeyFail) Validate(currentBlockNumber uint64, r RoleType, recoveredKeys []*ecdsa.PublicKey, from common.Address) bool {
	// This type of account key always fails to validate.
	return false
}

func (a *AccountKeyFail) String() string {
	return "AccountKeyFail"
}

func (a *AccountKeyFail) DeepCopy() AccountKey {
	return NewAccountKeyFail()
}

func (a *AccountKeyFail) AccountCreationGas(currentBlockNumber uint64) (uint64, error) {
	// No gas required to make an account with a failed key.
	return 0, nil
}

func (a *AccountKeyFail) SigValidationGas(currentBlockNumber uint64, r RoleType, numSigs int) (uint64, error) {
	// No gas required to make an account with a failed key.
	return 0, nil
}

func (a *AccountKeyFail) CheckInstallable(currentBlockNumber uint64) error {
	// AccountKeyFail can be assigned to an account. Since it does not have any value, it returns always nil.
	return nil
}

func (a *AccountKeyFail) CheckUpdatable(newKey AccountKey, currentBlockNumber uint64) error {
	// AccountKeyFail cannot be updated with any key, hence it returns always an error.
	return kerrors.ErrAccountKeyFailNotUpdatable
}

func (a *AccountKeyFail) Update(newKey AccountKey, currentBlockNumber uint64) error {
	// AccountKeyFail cannot be updated with any key, hence it returns always an error.
	return kerrors.ErrAccountKeyFailNotUpdatable
}
