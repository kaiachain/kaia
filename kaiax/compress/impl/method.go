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
	"bytes"
	"io"

	"github.com/klauspost/compress/s2"
	"github.com/klauspost/compress/zstd"
)

type (
	CompressMethod   func(src []byte) (int, []byte)
	DecompressMethod func(compressedBuffer []byte) ([]byte, error)
)

var (
	DefaultCompressor   = CompressZstd
	DefaultDecompressor = DecompressZstd
	Compress            = DefaultCompressor
	Decompress          = DefaultDecompressor
)

func CompressZstd(src []byte) (int, []byte) {
	var (
		srcReader = bytes.NewReader(src)
		dst       bytes.Buffer
	)
	compressor, err := zstd.NewWriter(&dst)
	if err != nil {
		logger.Warn("[Compression] Failed to create compressor", "err", err)
		return 0, nil
	}
	_, err = io.Copy(compressor, srcReader)
	if err != nil {
		logger.Warn("[Compression] Failed to compress", "err", err)
		return 0, nil
	}
	err = compressor.Close()
	if err != nil {
		logger.Warn("[Compression] Failed to close compressor", "err", err)
		return 0, nil
	}
	return dst.Len(), dst.Bytes()
}

func DecompressZstd(compressedBuffer []byte) ([]byte, error) {
	compressedReader := bytes.NewReader(compressedBuffer)
	decompressor, err := zstd.NewReader(compressedReader, zstd.WithDecoderConcurrency(0))
	if err != nil {
		logger.Warn("[Compression] Failed to create zstd decompressor", "err", err)
		return nil, err
	}
	defer decompressor.Close()

	var decompressedBuffer bytes.Buffer
	_, err = io.Copy(&decompressedBuffer, decompressor)
	if err != nil {
		logger.Warn("[Compression] Failed to decompress", "err", err)
		return nil, err
	}
	return decompressedBuffer.Bytes(), nil
}

func CompressS2(src []byte) (int, []byte) {
	var (
		dst        bytes.Buffer
		compressor = s2.NewWriter(&dst)
	)
	defer compressor.Close()
	err := compressor.EncodeBuffer(src)
	if err != nil {
		logger.Warn("[Compression] Failed to compress", "err", err)
		return 0, nil
	}
	return dst.Len(), dst.Bytes()
}

func DecompressS2(compressedBuffer []byte) ([]byte, error) {
	compressedReader := bytes.NewReader(compressedBuffer)
	decompressor, err := zstd.NewReader(compressedReader, zstd.WithDecoderConcurrency(0))
	if err != nil {
		logger.Warn("[Compression] Failed to create zstd decompressor", "err", err)
		return nil, err
	}
	defer decompressor.Close()

	var decompressedBuffer bytes.Buffer
	_, err = io.Copy(&decompressedBuffer, decompressor)
	if err != nil {
		logger.Warn("[Compression] Failed to decompress", "err", err)
		return nil, err
	}
	return decompressedBuffer.Bytes(), nil
}
