// Copyright 2024 The Kaia Authors
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

package valset

import (
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax"
)

//go:generate mockgen -destination=./mock/module.go -package=mock github.com/kaiachain/kaia/v2/kaiax/valset ValsetModule
type ValsetModule interface {
	kaiax.BaseModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	GetCouncil(num uint64) ([]common.Address, error)
	GetCommittee(num uint64, round uint64) ([]common.Address, error)
	GetDemotedValidators(num uint64) ([]common.Address, error)
	GetProposer(num uint64, round uint64) (common.Address, error)
}
