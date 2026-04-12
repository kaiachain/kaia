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

package dbsyncer

import (
	"strings"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
)

const (
	// mysql driver (go) has query parameter max size (65535)
	// transaction record has 15 parameters
	BULK_INSERT_SIZE = 3000
	TX_KEY_FACTOR    = 100000
)

func getProposerAndValidatorsFromBlock(s consensus.Sealer, block *types.Block) (proposer string, validators string, err error) {
	blockNumber := block.NumberU64()
	if blockNumber == 0 {
		return "", "", nil
	}
	proposerAddr, err := s.Author(block.Header())
	if err != nil {
		return "", "", err
	}
	validatorsList, err := s.Validators(block.Header())
	if err != nil {
		return "", "", err
	}
	var strValidators []string
	for _, validator := range validatorsList {
		strValidators = append(strValidators, validator.Hex())
	}

	return proposerAddr.Hex(), strings.Join(strValidators, ","), nil
}

// ecrecover extracts the Kaia account address from a signed header.
func ecrecover(s consensus.Sealer, header *types.Header) (common.Address, error) {
	author, err := s.Author(header)
	return author, err
}
