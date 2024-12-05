package database

import (
	"io"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
	"github.com/klauspost/compress/zstd"
)

type CompressStructTyp interface {
	GetBlkNumber() uint64
	GetBlkHash() common.Hash
}

type (
	Finder       func(uint64, uint64, uint64, common.Hash) (any, error)
	CompressFn   func(from, to, headNumber uint64, migrationMode bool) (uint64, error)
	DecompressFn func(compressTyp CompressionType, from, to uint64) ([]CompressStructTyp, error)
)

// Create a writer that caches compressors.
// For this operation type we supply a nil Reader.
var (
	encoder, _ = zstd.NewWriter(nil)
	decoder, _ = zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
)

// Compress a buffer.
// If you have a destination buffer, the allocation in the call can also be eliminated.
func Compress(src []byte) []byte {
	return encoder.EncodeAll(src, make([]byte, 0, len(src)))
}

// Decompress a buffer. We don't supply a destination buffer,
// so it will be allocated by the decoder.
func Decompress(src []byte) ([]byte, error) {
	return decoder.DecodeAll(src, nil)
}

type ReceiptCompression struct {
	BlkNumber       uint64
	BlkHash         common.Hash
	StorageReceipts []*types.ReceiptForStorage
}

type ReceiptCompressionRLP struct {
	BlkNumber       uint64
	BlkHash         common.Hash
	StorageReceipts []*types.ReceiptForStorage
}

func (r *ReceiptCompression) GetBlkNumber() uint64 {
	return r.BlkNumber
}

func (r *ReceiptCompression) GetBlkHash() common.Hash {
	return r.BlkHash
}

func (r *ReceiptCompression) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &ReceiptCompressionRLP{BlkNumber: r.BlkNumber, BlkHash: r.BlkHash, StorageReceipts: r.StorageReceipts})
}

func (r *ReceiptCompression) DecodeRLP(s *rlp.Stream) error {
	var dec ReceiptCompressionRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	r.BlkNumber = dec.BlkNumber
	r.BlkHash = dec.BlkHash
	r.StorageReceipts = dec.StorageReceipts
	return nil
}

type BodyCompression struct {
	BlkNumber uint64
	BlkHash   common.Hash
	Body      *types.Body
}

type BodyCompressionRLP struct {
	BlkNumber uint64
	BlkHash   common.Hash
	Body      *types.Body
}

func (r *BodyCompression) GetBlkNumber() uint64 {
	return r.BlkNumber
}

func (r *BodyCompression) GetBlkHash() common.Hash {
	return r.BlkHash
}

func (r *BodyCompression) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &BodyCompressionRLP{BlkNumber: r.BlkNumber, BlkHash: r.BlkHash, Body: r.Body})
}

func (r *BodyCompression) DecodeRLP(s *rlp.Stream) error {
	var dec BodyCompressionRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	r.BlkNumber = dec.BlkNumber
	r.BlkHash = dec.BlkHash
	r.Body = dec.Body
	return nil
}

// TODO-hyunsooda: Implement header compress RLP
