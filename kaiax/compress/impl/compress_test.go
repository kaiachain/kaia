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
	"math/big"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
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
		nBlocks                           = 100
		_, dbm, headers, bodies, receipts = runCompress(t, nBlocks)
	)
	checkCompressedIntegrity(t, dbm, 0, nBlocks-1, headers, bodies, receipts, false)
}

func TestRewind(t *testing.T) {
	var (
		nBlocks                                   = 100
		setHeadTo                                 = 55
		mCompress, dbm, headers, bodies, receipts = runCompress(t, nBlocks)
	)

	// Scenario Description:
	// 1. `setHead` invoked targeting block number 55
	//    - Removed origin data from 100 to 55
	//    - Removed compressed data from 100 to 50 (because chunk size is set to 10)
	// 2. Generate new blocks from 55 to 100
	// 3. Restart the compression module
	// 	  - start compression from block number 50 to 100

	for i := nBlocks - 1; i >= setHeadTo; i-- {
		num := uint64(i)
		hash := dbm.ReadCanonicalHash(num)
		// delete origin data
		dbm.DeleteBody(hash, num)
		dbm.DeleteReceipts(hash, num)
		dbm.DeleteHeader(hash, num)
		// delete compression data
		mCompress.RewindDelete(hash, num)
	}
	checkCompressedIntegrity(t, dbm, setHeadTo, nBlocks-1, headers, bodies, receipts, true)

	var (
		newHeaders  []*types.Header
		newBodies   []*types.Body
		newReceipts []types.Receipts
	)
	for i := setHeadTo; i < nBlocks; i++ {
		h := genHeader()
		h.Number = big.NewInt(int64(i))
		hn, hh := h.Number.Uint64(), h.Hash()
		dbm.WriteCanonicalHash(hh, hn)
		dbm.WriteHeader(h)

		b, r := genBody(100), genReceipts(100)
		dbm.WriteBody(hh, hn, b)
		dbm.WriteReceipts(hh, hn, *r)

		newHeaders = append(newHeaders, h)
		newBodies = append(newBodies, b)
		newReceipts = append(newReceipts, *r)
	}
	canonicalHeaders := append(headers[:setHeadTo], newHeaders[:]...)
	canonicalBodies := append(bodies[:setHeadTo], newBodies[:]...)
	canonicalReceipts := append(receipts[:setHeadTo], newReceipts[:]...)

	go mCompress.Compress()
	time.Sleep(time.Second)
	waitCompression(mCompress)
	checkCompressedIntegrity(t, dbm, 0, nBlocks-1, canonicalHeaders, canonicalBodies, canonicalReceipts, false)
}
