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
	"testing"
)

func TestCompressStorage(t *testing.T) {
	for _, compressTyp := range []CompressionType{HeaderCompressType, BodyCompressType, ReceiptCompressType} {
		switch compressTyp {
		case HeaderCompressType:
			testHeaderCompress(t)
		case BodyCompressType:
			testBodyCompress(t)
		case ReceiptCompressType:
			testReceiptsCompress(t)
		}
	}
}

func TestDecompressStorage(t *testing.T) {
	for _, compressTyp := range []CompressionType{HeaderCompressType, BodyCompressType, ReceiptCompressType} {
		switch compressTyp {
		case HeaderCompressType:
			testHeaderDecompress(t)
		case BodyCompressType:
			testBodyDecompress(t)
		case ReceiptCompressType:
			testReceiptsDecompress(t)
		}
	}
}

func TestDeleteStorage(t *testing.T) {
	for _, compressTyp := range []CompressionType{HeaderCompressType, BodyCompressType, ReceiptCompressType} {
		switch compressTyp {
		case HeaderCompressType:
			testCompressedHeaderDelete(t)
		case BodyCompressType:
			testCompressedBodyDelete(t)
		case ReceiptCompressType:
			testCompressedReceiptsDelete(t)
		}
	}
}

func TestCompressModule(t *testing.T) {
	var (
		nBlocks = 100
		_, dbm  = runCompress(t, nBlocks)
	)
	checkCompressedIntegrity(t, dbm, nBlocks-1, false)
}

func TestRewind(t *testing.T) {
	var (
		nBlocks        = 100
		mCompress, dbm = runCompress(t, nBlocks)
	)
	for i := range nBlocks - 1 {
		num := uint64(i)
		hash := dbm.ReadCanonicalHash(num)
		mCompress.RewindDelete(hash, num)
	}
	checkCompressedIntegrity(t, dbm, nBlocks-1, true)
}
