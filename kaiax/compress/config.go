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

import "github.com/urfave/cli/v2"

var (
	DataCompressFlag = &cli.BoolFlag{
		Name:     "data.compress",
		Usage:    "Enable data compression",
		Aliases:  []string{},
		EnvVars:  []string{"KLAYTN_DATA_COMPRESS", "KAIA_DATA_COMPRESS"},
		Category: "DATA",
	}
	DataCompressRetentionFlag = &cli.Uint64Flag{
		Name:    "data.compress.retention",
		Usage:   "Number of blocks from the latest block whose data should not be compressed",
		Value:   DefaultRetention,
		Aliases: []string{},
	}
	DataCompressChunkItemCapFlag = &cli.IntFlag{
		Name:    "data.compress.item-cap",
		Usage:   "Maximum number of items in a compressed chunk",
		Value:   DefaultChunkItemCap,
		Aliases: []string{},
	}
	DataCompressChunkByteCapFlag = &cli.IntFlag{
		Name:    "data.compress.byte-cap",
		Usage:   "Maximum number of bytes in a compressed chunk",
		Value:   DefaultChunkByteCap,
		Aliases: []string{},
	}
)

const (
	// Bounds for sanity checks.
	MinRetention     = 128
	DefaultRetention = 172800 // 48 hours

	MinChunkItemCap     = 100
	DefaultChunkItemCap = 10000
	MaxChunkItemCap     = 1000000

	MinChunkByteCap     = 1024               // 1KB
	DefaultChunkByteCap = 1024 * 1024        // 1MB
	MaxChunkByteCap     = 1024 * 1024 * 1024 // 1GB
)

type CompressConfig struct {
	// True to enable compression.
	// False to disable compression, but sill this module support reading from the compressed database,
	// and uncompression upon block rewind.
	Enabled bool

	Retention    uint64 // number of blocks to keep in the uncompressed database
	ChunkItemCap int    // maximum number of items in a chunk
	ChunkByteCap int    // maximum size of uncompressed data in a chunk
}

func GetDefaultCompressConfig() *CompressConfig {
	return &CompressConfig{
		Enabled:      false,
		Retention:    DefaultRetention,
		ChunkItemCap: DefaultChunkItemCap,
		ChunkByteCap: DefaultChunkByteCap,
	}
}

func GetCompressConfig(ctx *cli.Context) *CompressConfig {
	return &CompressConfig{
		Enabled:      ctx.Bool(DataCompressFlag.Name),
		Retention:    ctx.Uint64(DataCompressRetentionFlag.Name),
		ChunkItemCap: ctx.Int(DataCompressChunkItemCapFlag.Name),
		ChunkByteCap: ctx.Int(DataCompressChunkByteCapFlag.Name),
	}
}
