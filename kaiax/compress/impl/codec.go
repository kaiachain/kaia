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

package impl

import (
	"bytes"
	"io"

	"github.com/kaiachain/kaia/rlp"
	"github.com/klauspost/compress/zstd"
)

var _ Codec = (*ZstdCodec)(nil)

type Codec interface {
	compress(src []byte) ([]byte, error)
	decompress(src []byte) ([]byte, error)
}

func compressChunk(codec Codec, chunk []ChunkItem) ([]byte, error) {
	bytes, err := rlp.EncodeToBytes(chunk)
	if err != nil {
		return nil, err
	}
	return codec.compress(bytes)
}

func decompressChunk(codec Codec, compressed []byte) ([]ChunkItem, error) {
	bytes, err := codec.decompress(compressed)
	if err != nil {
		return nil, err
	}
	chunk := []ChunkItem{}
	err = rlp.DecodeBytes(bytes, &chunk)
	return chunk, err
}

type ZstdCodec struct{}

func NewZstdCodec() *ZstdCodec {
	return &ZstdCodec{}
}

func (c *ZstdCodec) compress(src []byte) ([]byte, error) {
	var dst bytes.Buffer
	reader := bytes.NewReader(src)
	writer, err := zstd.NewWriter(&dst)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(writer, reader); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}

func (c *ZstdCodec) decompress(src []byte) ([]byte, error) {
	var dst bytes.Buffer
	reader, err := zstd.NewReader(bytes.NewReader(src), zstd.WithDecoderConcurrency(0))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	if _, err = io.Copy(&dst, reader); err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
}
