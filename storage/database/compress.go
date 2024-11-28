package database

import (
	"io"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
	"github.com/klauspost/compress/zstd"
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
	number          uint64
	blkHash         common.Hash
	storageReceipts []*types.ReceiptForStorage
}

type ReceiptCompressionRLP struct {
	Number          uint64
	BlkHash         common.Hash
	StorageReceipts []*types.ReceiptForStorage
}

func (r *ReceiptCompression) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &ReceiptCompressionRLP{Number: r.number, BlkHash: r.blkHash, StorageReceipts: r.storageReceipts})
}

func (r *ReceiptCompression) DecodeRLP(s *rlp.Stream) error {
	var dec ReceiptCompressionRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	r.number = dec.Number
	r.blkHash = dec.BlkHash
	r.storageReceipts = dec.StorageReceipts
	return nil
}

// TODO-hyunsooda: Implement header compress RLP
// TODO-hyunsooda: Implement body compress RLP
