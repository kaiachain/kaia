// Modifications Copyright 2024 The Kaia Authors
// Copyright 2023 The klaytn Authors
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

package blst

import (
	"crypto/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/kaiachain/kaia/v2/crypto/bls/types"
)

var benchAggregateLen = 100

func BenchmarkPublicKeyFromBytes(b *testing.B) {
	sk, _ := RandKey()
	pkb := sk.PublicKey().Marshal()

	fn := func() {
		if _, err := PublicKeyFromBytes(pkb); err != nil {
			b.Fatal(err)
		}
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

func BenchmarkMultiplePublicKeysFromBytes(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)

	fn := func() {
		if _, err := MultiplePublicKeysFromBytes(tc.pkbs); err != nil {
			b.Fatal(err)
		}
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

func BenchmarkAggregatePublicKeys(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)

	for i := 0; i < b.N; i++ {
		if _, err := AggregatePublicKeys(tc.pks); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAggregatePublicKeysFromBytes(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)

	fn := func() {
		if _, err := AggregatePublicKeysFromBytes(tc.pkbs); err != nil {
			b.Fatal(err)
		}
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

func BenchmarkSignatureFromBytes(b *testing.B) {
	sigb := testSignatureBytes

	fn := func() {
		if _, err := SignatureFromBytes(sigb); err != nil {
			b.Fatal(err)
		}
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

func BenchmarkMultipleSignaturesFromBytes(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)

	fn := func() {
		if _, err := MultipleSignaturesFromBytes(tc.sigbs); err != nil {
			b.Fatal(err)
		}
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

func BenchmarkAggregateSignatures(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)

	for i := 0; i < b.N; i++ {
		if _, err := AggregateSignatures(tc.sigs); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAggregateSignaturesFromBytes(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)

	fn := func() {
		if _, err := AggregateSignaturesFromBytes(tc.sigbs); err != nil {
			b.Fatal(err)
		}
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

func BenchmarkSign(b *testing.B) {
	sk, _ := RandKey()
	msg := []byte("test message")

	for i := 0; i < b.N; i++ {
		Sign(sk, msg)
	}
}

func BenchmarkVerify(b *testing.B) {
	sk, _ := RandKey()
	pk := sk.PublicKey()
	msg := []byte("test message")
	sig := Sign(sk, msg)

	for i := 0; i < b.N; i++ {
		Verify(sig, msg, pk)
	}
}

func BenchmarkParallelVerify(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterialMulti(L)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		threads := runtime.NumCPU()
		var wg sync.WaitGroup
		wg.Add(threads)

		jobs := make(chan int, len(tc.sigbs))
		for i := 0; i < L; i++ {
			jobs <- i
		}

		for i := 0; i < threads; i++ {
			go func() {
				for i := range jobs {
					sig, _ := SignatureFromBytes(tc.sigbs[i])
					pk, _ := PublicKeyFromBytes(tc.pkbs[i])
					Verify(sig, tc.msgs[i][:], pk)
				}
				wg.Done()
			}()
		}

		close(jobs)
		wg.Wait()
	}
}

// End-to-end benchmark where all inputs are []byte
// Aggregate-Then-SingleVerify
func BenchmarkAggregateAndVerify(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)
	asig, _ := AggregateSignatures(tc.sigs)
	asigb := asig.Marshal()

	fn := func() {
		apk, _ := AggregatePublicKeysFromBytes(tc.pkbs)
		sig, _ := SignatureFromBytes(asigb)
		Verify(sig, tc.msg, apk)
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

// End-to-end benchmark where all inputs are []byte
// FastAggregateVerify
func BenchmarkFastAggregateVerify(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterial(L)
	asig, _ := AggregateSignatures(tc.sigs)
	asigb := asig.Marshal()

	fn := func() {
		pks, _ := MultiplePublicKeysFromBytes(tc.pkbs)
		sig, _ := SignatureFromBytes(asigb)
		FastAggregateVerify(sig, tc.msg, pks)
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

// End-to-end benchmark where all inputs are []byte
// Distinct messages, VerifyMultipleSignatures
func BenchmarkVerifyMultipleSignatures(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterialMulti(L)

	fn := func() {
		pks, _ := MultiplePublicKeysFromBytes(tc.pkbs)
		VerifyMultipleSignatures(tc.sigbs, tc.msgs, pks)
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

// End-to-end benchmark where all inputs are []byte
// Distinct messages, verify each signature one by one
func BenchmarkVerifyMultipleNaive(b *testing.B) {
	L := benchAggregateLen
	tc := generateBenchmarkMaterialMulti(L)

	fn := func() {
		pks, _ := MultiplePublicKeysFromBytes(tc.pkbs)
		for i := 0; i < L; i++ {
			VerifySignature(tc.sigbs[i], tc.msgs[i], pks[i])
		}
	}
	runUncached(b, "Uncached", fn)
	runCached(b, "Cached", fn)
}

type benchmarkTestCase struct {
	sks   []types.SecretKey
	pks   []types.PublicKey
	pkbs  [][]byte
	msg   []byte
	sigs  []types.Signature
	sigbs [][]byte
	asig  types.Signature
}

func generateBenchmarkMaterial(L int) *benchmarkTestCase {
	tc := &benchmarkTestCase{}
	tc.sks = make([]types.SecretKey, L)
	tc.pks = make([]types.PublicKey, L)
	tc.pkbs = make([][]byte, L)
	tc.msg = make([]byte, 32)
	tc.sigs = make([]types.Signature, L)
	tc.sigbs = make([][]byte, L)

	rand.Read(tc.msg)

	for i := 0; i < L; i++ {
		sk, _ := RandKey()
		pk := sk.PublicKey()
		sig := Sign(sk, tc.msg)
		tc.sks[i] = sk
		tc.pks[i] = pk
		tc.pkbs[i] = pk.Marshal()
		tc.sigs[i] = sig
		tc.sigbs[i] = sig.Marshal()
	}
	return tc
}

type benchmarkTestCaseMulti struct {
	sks   []types.SecretKey
	pks   []types.PublicKey
	pkbs  [][]byte
	msgs  [][32]byte
	sigs  []types.Signature
	sigbs [][]byte
	asig  types.Signature
}

func generateBenchmarkMaterialMulti(L int) *benchmarkTestCaseMulti {
	tc := &benchmarkTestCaseMulti{}
	tc.sks = make([]types.SecretKey, L)
	tc.pks = make([]types.PublicKey, L)
	tc.pkbs = make([][]byte, L)
	tc.msgs = make([][32]byte, L)
	tc.sigs = make([]types.Signature, L)
	tc.sigbs = make([][]byte, L)

	for i := 0; i < L; i++ {
		sk, _ := RandKey()
		pk := sk.PublicKey()
		rand.Read(tc.msgs[i][:])
		sig := Sign(sk, tc.msgs[i][:])
		tc.sks[i] = sk
		tc.pks[i] = pk
		tc.pkbs[i] = pk.Marshal()
		tc.sigs[i] = sig
		tc.sigbs[i] = sig.Marshal()
	}
	return tc
}

func runUncached(b *testing.B, name string, fn func()) {
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			resetCache()
			b.StartTimer()
			fn()
		}
	})
}

func runCached(b *testing.B, name string, fn func()) {
	for i := 0; i < b.N; i++ {
		fn() // populate cache
	}
	b.Run(name, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			fn()
		}
	})
}

func resetCache() {
	publicKeyCache.Purge()
	signatureCache.Purge()
}
