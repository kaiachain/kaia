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

package compress

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=mock/compress.go -package=mock github.com/kaiachain/kaia/kaiax/compress CompressModule
type CompressModule interface {
	kaiax.BaseModule
	kaiax.RewindableModule

	// FindFromCompressed functions return the uncompressed data from compressed databases.
	FindFromCompressedHeader(num uint64, hash common.Hash) ([]byte, bool)
	FindFromCompressedBody(num uint64, hash common.Hash) ([]byte, bool)
	FindFromCompressedReceipts(num uint64, hash common.Hash) ([]byte, bool)
}
