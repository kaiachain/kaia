// Copyright 2026 The Kaia Authors
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

package kzg4844

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"math/bits"
	"sync"
	"testing"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	gokzg4844 "github.com/crate-crypto/go-eth-kzg"
)

// go-eth-kzg v1.3.0 derived the random linear combination of VerifyCellKZGProofBatch
// from a Fiat-Shamir transcript that omitted the proof points. An attacker could
// therefore compute the challenge r before choosing the proofs, and correlate two
// invalid proof points so that the aggregate pairing equation still held. Invalid cell
// proofs consequently passed batch verification, which let BlobTxSidecar.ValidateWithBlobHashes
// accept a sidecar whose blobs do not belong to the committed commitments: the same
// signed transaction could be paired with conflicting sidecars without changing its
// signature or transaction hash.
//
// In v1.5.0 the Go verifier removes that transcript and samples r randomly per
// verification, so a forgery aimed at a precomputed challenge no longer verifies. The
// tests below rebuild the v1.3.0 transcript, forge proofs against it, and assert that
// both backends reject them. The forged vector is accepted by the v1.3.0 Go verifier
// and rejected by the v1.5.0 Go verifier; the C backend independently rejects it when
// available.
//
// The defect was identified in the Ethereum Foundation PeerDAS KZG audit and fixed
// upstream in crate-crypto/go-eth-kzg#111, first released in v1.4.0. The forgery
// construction below is derived from the proof of concept attached to the security
// report that flagged it for Kaia.

// transcriptUint64 encodes v the way the v1.3.0 transcript did: big endian in 16 bytes.
func transcriptUint64(v uint64) []byte {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[8:], v)
	return b
}

// vulnerableBatchChallenge reproduces go-eth-kzg v1.3.0's transcript for a single
// commitment covering all of its cells. The proofs are absent from the transcript,
// which is what allows them to be chosen after r.
func vulnerableBatchChallenge(commitment Commitment, cells [gokzg4844.CellsPerExtBlob]*gokzg4844.Cell) fr.Element {
	const domainSepProtocol = "RCKZGCBATCH__V1_"

	point, err := gokzg4844.DeserializeKZGCommitment(gokzg4844.KZGCommitment(commitment))
	if err != nil {
		panic(err)
	}
	h := sha256.New()
	h.Write([]byte(domainSepProtocol))
	h.Write(transcriptUint64(gokzg4844.ScalarsPerBlob))
	h.Write(transcriptUint64(gokzg4844.BytesPerCell / gokzg4844.SerializedScalarSize))
	h.Write(transcriptUint64(1))
	h.Write(transcriptUint64(gokzg4844.CellsPerExtBlob))
	h.Write(point.Marshal())
	for i, cell := range cells {
		h.Write(transcriptUint64(0))
		h.Write(transcriptUint64(uint64(i)))
		h.Write(cell[:])
	}
	var challenge fr.Element
	challenge.SetBytes(h.Sum(nil))
	return challenge
}

// batchCosetWeight returns the weight the batch verifier applies to the proof of the
// cell at index: the coset shift raised to the coset size.
func batchCosetWeight(index uint64) fr.Element {
	// rootOfUnity is the 2^32-th root of unity of the BLS12-381 scalar field.
	const rootOfUnityDec = "10238227357739495823651030575849232062558860180284477541189508159991286009131"
	const extendedDomainSize = uint64(gokzg4844.ScalarsPerBlob * 2)

	var rootOfUnity fr.Element
	if _, err := rootOfUnity.SetString(rootOfUnityDec); err != nil {
		panic(err)
	}
	rootOfUnity.Exp(rootOfUnity, new(big.Int).SetUint64(1<<(32-bits.TrailingZeros64(extendedDomainSize))))

	// NewOpeningKey bit-reverses the extended domain roots and selects root[index*64].
	rootIndex := bits.Reverse64(index*64) >> (64 - bits.TrailingZeros64(extendedDomainSize))
	var shift fr.Element
	shift.Exp(rootOfUnity, new(big.Int).SetUint64(rootIndex))
	return *new(fr.Element).Exp(shift, big.NewInt(64))
}

// multiplyG1 returns scalar * point.
func multiplyG1(point *bls12381.G1Affine, scalar fr.Element) bls12381.G1Affine {
	var n big.Int
	scalar.BigInt(&n)
	var result bls12381.G1Affine
	result.ScalarMultiplication(point, &n)
	return result
}

// forgeCellProofs computes the cell proofs of blob and then perturbs the first two of
// them so that the batch verifier's aggregate pairing equation still holds when the
// batch is verified against commitment, which does not belong to blob.
//
// The two perturbations delta0 and delta1 satisfy
//
//	delta0 + r*delta1 = 0
//	q0*delta0 + r*q1*delta1 = -(sum of r^i) * (commitment - commitment(blob))
//
// where q0 and q1 are the public coset weights. The first equation leaves the weighted
// proof sum unchanged, the second cancels the commitment difference on the other side of
// the aggregate equation.
func forgeCellProofs(commitment Commitment, blob *Blob) ([]Proof, [gokzg4844.CellsPerExtBlob]*gokzg4844.Cell) {
	blobCommitment, err := BlobToCommitment(blob)
	if err != nil {
		panic(err)
	}
	proofs, err := ComputeCellProofs(blob)
	if err != nil {
		panic(err)
	}
	gokzgIniter.Do(gokzgInit)
	cells, err := context.ComputeCells((*gokzg4844.Blob)(blob), 0)
	if err != nil {
		panic(err)
	}
	r := vulnerableBatchChallenge(commitment, cells)

	// delta is the r-weighted difference between the commitment the batch is verified
	// against and the commitment that actually belongs to the blob.
	var sumOfPowers, power fr.Element
	power.SetOne()
	for range gokzg4844.CellsPerExtBlob {
		sumOfPowers.Add(&sumOfPowers, &power)
		power.Mul(&power, &r)
	}
	point, err := gokzg4844.DeserializeKZGCommitment(gokzg4844.KZGCommitment(commitment))
	if err != nil {
		panic(err)
	}
	blobPoint, err := gokzg4844.DeserializeKZGCommitment(gokzg4844.KZGCommitment(blobCommitment))
	if err != nil {
		panic(err)
	}
	var difference bls12381.G1Affine
	difference.Sub(&point, &blobPoint)
	delta := multiplyG1(&difference, sumOfPowers)

	q0, q1 := batchCosetWeight(0), batchCosetWeight(1)
	var denominator, delta0Scalar fr.Element
	denominator.Sub(&q0, &q1)
	delta0Scalar.Inverse(&denominator).Neg(&delta0Scalar)
	delta0 := multiplyG1(&delta, delta0Scalar)

	var inverseR, delta1Scalar fr.Element
	inverseR.Inverse(&r)
	delta1Scalar.Mul(&delta0Scalar, &inverseR).Neg(&delta1Scalar)
	delta1 := multiplyG1(&delta, delta1Scalar)

	forged := append([]Proof(nil), proofs...)
	for i, d := range []bls12381.G1Affine{delta0, delta1} {
		proofPoint, err := gokzg4844.DeserializeKZGProof(gokzg4844.KZGProof(forged[i]))
		if err != nil {
			panic(err)
		}
		var forgedPoint bls12381.G1Affine
		forgedPoint.Add(&proofPoint, &d)
		forged[i] = Proof(gokzg4844.SerializeG1Point(forgedPoint))
	}
	return forged, cells
}

// forgery holds the fixture shared by the tests below. The blobs are fixed so that the
// forged proofs are reproducible, and it is built with the Go backend so that the
// forgery always targets the same input regardless of which backend verifies it.
type forgery struct {
	// committedBlob is the blob the transaction commits to, different is the blob a
	// forged sidecar would carry instead.
	committedBlob, different         Blob
	commitment, differentCommitment  Commitment
	committedProofs, differentProofs []Proof

	// forgedProofs are different's cell proofs with the first two points perturbed so
	// that the batch verifies against commitment instead of differentCommitment.
	forgedProofs []Proof
	// cells are different's cells, in the order the batch verifier receives them.
	cells [gokzg4844.CellsPerExtBlob]*gokzg4844.Cell
}

// deterministicBlob returns a reproducible blob whose scalars are all distinct, so that
// the cells and the cell proofs of the blob are distinct as well. Leaving the top byte of
// every scalar zero keeps them canonical.
func deterministicBlob(seed uint64) Blob {
	var blob Blob
	for i := range gokzg4844.ScalarsPerBlob {
		offset := i * gokzg4844.SerializedScalarSize
		binary.BigEndian.PutUint64(blob[offset+24:offset+32], seed+uint64(i))
	}
	return blob
}

var newForgery = sync.OnceValue(func() *forgery {
	defer func(old bool) { useCKZG.Store(old) }(useCKZG.Load())
	useCKZG.Store(false)

	f := new(forgery)
	f.committedBlob = deterministicBlob(1)
	f.different = deterministicBlob(1 << 32)

	var err error
	if f.commitment, err = BlobToCommitment(&f.committedBlob); err != nil {
		panic(err)
	}
	if f.differentCommitment, err = BlobToCommitment(&f.different); err != nil {
		panic(err)
	}
	if f.commitment == f.differentCommitment {
		panic("the two blobs must not share a commitment, or the forgery below is vacuous")
	}
	if f.committedProofs, err = ComputeCellProofs(&f.committedBlob); err != nil {
		panic(err)
	}
	if f.differentProofs, err = ComputeCellProofs(&f.different); err != nil {
		panic(err)
	}
	f.forgedProofs, f.cells = forgeCellProofs(f.commitment, &f.different)
	return f
})

// cellProofVector is one accept/reject vector for VerifyCellProofs.
type cellProofVector struct {
	name    string
	blobs   []Blob
	commits []Commitment
	proofs  []Proof
	valid   bool
}

func cellProofVectors() []cellProofVector {
	f := newForgery()

	// A valid proof taken from the wrong cell index. Unlike a randomly mangled proof
	// this still deserializes, so it reaches the pairing check rather than failing
	// during input validation.
	misplaced := append([]Proof(nil), f.committedProofs...)
	misplaced[0], misplaced[1] = misplaced[1], misplaced[0]

	bothBlobs := []Blob{f.committedBlob, f.different}
	bothProofs := append(append([]Proof(nil), f.committedProofs...), f.differentProofs...)

	return []cellProofVector{
		{
			// Valid batches must keep passing.
			name:    "valid",
			blobs:   []Blob{f.committedBlob},
			commits: []Commitment{f.commitment},
			proofs:  f.committedProofs,
			valid:   true,
		},
		{
			// The multi-commitment path is what Kaia uses for multi-blob sidecars.
			name:    "valid batch of two blobs",
			blobs:   bothBlobs,
			commits: []Commitment{f.commitment, f.differentCommitment},
			proofs:  bothProofs,
			valid:   true,
		},
		{
			// A blob paired with another blob's commitment must fail.
			name:    "blob paired with another blob's commitment",
			blobs:   []Blob{f.different},
			commits: []Commitment{f.commitment},
			proofs:  f.differentProofs,
		},
		{
			// The same, in the multi-commitment path.
			name:    "swapped commitments",
			blobs:   bothBlobs,
			commits: []Commitment{f.differentCommitment, f.commitment},
			proofs:  bothProofs,
		},
		{
			// A single misplaced proof point must fail.
			name:    "proofs of two cells swapped",
			blobs:   []Blob{f.committedBlob},
			commits: []Commitment{f.commitment},
			proofs:  misplaced,
		},
		{
			// The regression: correlated modifications to two proof points must not
			// make the batch pass. This vector verified with go-eth-kzg v1.3.0.
			name:    "correlated proof forgery against the v1.3.0 challenge",
			blobs:   []Blob{f.different},
			commits: []Commitment{f.commitment},
			proofs:  f.forgedProofs,
		},
	}
}

func TestCKZGCellProofVectors(t *testing.T)  { testCellProofVectors(t, true) }
func TestGoKZGCellProofVectors(t *testing.T) { testCellProofVectors(t, false) }

// testCellProofVectors asserts that both backends agree on the accepted and rejected
// vectors, in particular that the correlated proof forgery is rejected.
func testCellProofVectors(t *testing.T, ckzg bool) {
	if ckzg && !ckzgAvailable {
		t.Skip("CKZG unavailable in this test build")
	}
	vectors := cellProofVectors()

	defer func(old bool) { useCKZG.Store(old) }(useCKZG.Load())
	useCKZG.Store(ckzg)

	for _, vector := range vectors {
		t.Run(vector.name, func(t *testing.T) {
			err := VerifyCellProofs(vector.blobs, vector.commits, vector.proofs)
			if vector.valid && err != nil {
				t.Fatalf("valid cell proofs rejected: %v", err)
			}
			if !vector.valid && err == nil {
				t.Fatal("invalid cell proofs accepted")
			}
		})
	}
}

// The recorded forgery. forgeCellProofs derives the perturbations from library internals
// (the v1.3.0 transcript and the coset shifts of the opening key). If an upstream change
// silently invalidates that derivation, the forged batch would be rejected for the wrong
// reason and the regression vector above would pass without testing anything, so the
// derived proofs are pinned here.
const (
	wantForgedProof0 = "8267a551529005b61a6cd823fd4425891b563356a165b4f2398f40476356b19a3775ab1d39bc5754968e4c91662e30c3"
	wantForgedProof1 = "a4105300becface860a360071a1cf14bcf5f15a4cfbe6464f9fc082be6076483527b227e8fcfcdfd3606c4170ccf648f"
)

// TestForgedCellProofsAreOnlyValidInAggregate asserts that the regression vector is not
// vacuous: the forgery is the one that was derived against the v1.3.0 transcript, and
// neither of its two perturbed proofs is a valid opening on its own. The batch as a
// whole only ever verified because the aggregate pairing equation was satisfied.
func TestForgedCellProofsAreOnlyValidInAggregate(t *testing.T) {
	f := newForgery()

	for i, want := range []string{wantForgedProof0, wantForgedProof1} {
		if got := hex.EncodeToString(f.forgedProofs[i][:]); got != want {
			t.Fatalf("forged proof %d drifted from the recorded vector:\ngot  %s\nwant %s", i, got, want)
		}
	}

	defer func(old bool) { useCKZG.Store(old) }(useCKZG.Load())
	useCKZG.Store(false) // the raw batch API below is only exposed by the Go backend

	for i := range 2 {
		err := context.VerifyCellKZGProofBatch(
			[]gokzg4844.KZGCommitment{gokzg4844.KZGCommitment(f.commitment)},
			[]uint64{uint64(i)},
			[]*gokzg4844.Cell{f.cells[i]},
			[]gokzg4844.KZGProof{gokzg4844.KZGProof(f.forgedProofs[i])},
		)
		if err == nil {
			t.Fatalf("forged proof %d unexpectedly passed individual verification", i)
		}
	}
}
