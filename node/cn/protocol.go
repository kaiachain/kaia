// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of go-ethereum.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/protocol.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package cn

import (
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/datasync/fetcher"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/node/cn/snap"
	"github.com/kaiachain/kaia/rlp"
)

// `kaia` is the default fallback protocol for all consensus engine types.
// The name, versions, and lengths are defined in consensus/protocol.go.
const (
	kaia63 = 63
	kaia65 = 65
	kaia66 = 66
)

const ProtocolMaxMsgSize = 12 * 1024 * 1024 // Maximum cap on the size of a protocol message

// Kaia protocol message codes
// TODO-Klaytn-Issue751 Protocol message should be refactored. Present code is not used.
const (
	// Protocol messages belonging to kaia/62
	StatusMsg                   = 0x00
	NewBlockHashesMsg           = 0x01
	BlockHeaderFetchRequestMsg  = 0x02
	BlockHeaderFetchResponseMsg = 0x03
	BlockBodiesFetchRequestMsg  = 0x04
	BlockBodiesFetchResponseMsg = 0x05
	TxMsg                       = 0x06
	BlockHeadersRequestMsg      = 0x07
	BlockHeadersMsg             = 0x08
	BlockBodiesRequestMsg       = 0x09
	BlockBodiesMsg              = 0x0a
	NewBlockMsg                 = 0x0b

	// Protocol messages belonging to kaia/63
	NodeDataRequestMsg = 0x0c
	NodeDataMsg        = 0x0d
	ReceiptsRequestMsg = 0x0e
	ReceiptsMsg        = 0x0f

	// Protocol messages belonging to kaia/64
	Unused10 = 0x10 // Skipped a number because 0x11 is already taken
	Unused11 = 0x11 // Already used by consensus (IstanbulMsg)

	// Protocol messages belonging to kaia/65
	StakingInfoRequestMsg = 0x12
	StakingInfoMsg        = 0x13

	// Protocol messages belonging to kaia/66
	BidMsg = 0x14

	MsgCodeEnd = 0x15
)

type errCode int

const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrProtocolVersionMismatch
	ErrNetworkIdMismatch
	ErrGenesisBlockMismatch
	ErrChainIDMismatch
	ErrNoStatusMsg
	ErrExtraStatusMsg
	ErrSuspendedPeer
	ErrUnexpectedTxType
	ErrFailedToGetStateDB
	ErrUnsupportedEnginePolicy
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

// XXX change once legacy code is out
var errorToString = map[int]string{
	ErrMsgTooLarge:             "Message too long",
	ErrDecode:                  "Invalid message",
	ErrInvalidMsgCode:          "Invalid message code",
	ErrProtocolVersionMismatch: "Protocol version mismatch",
	ErrNetworkIdMismatch:       "NetworkId mismatch",
	ErrGenesisBlockMismatch:    "Genesis block mismatch",
	ErrChainIDMismatch:         "ChainID mismatch",
	ErrNoStatusMsg:             "No status message",
	ErrExtraStatusMsg:          "Extra status message",
	ErrSuspendedPeer:           "Suspended peer",
	ErrUnexpectedTxType:        "Unexpected tx type",
	ErrFailedToGetStateDB:      "Failed to get stateDB",
	ErrUnsupportedEnginePolicy: "Unsupported engine or policy",
}

// ProtocolManagerDownloader is an interface of downloader.Downloader used by ProtocolManager.
//
//go:generate mockgen -destination=./mocks/downloader_mock.go -package=mocks github.com/kaiachain/kaia/node/cn ProtocolManagerDownloader
type ProtocolManagerDownloader interface {
	RegisterPeer(id string, version int, peer downloader.Peer) error
	UnregisterPeer(id string) error

	DeliverBodies(id string, transactions [][]*types.Transaction) error
	DeliverHeaders(id string, headers []*types.Header) error
	DeliverNodeData(id string, data [][]byte) error
	DeliverReceipts(id string, receipts [][]*types.Receipt) error
	DeliverStakingInfos(id string, stakingInfos []*staking.P2PStakingInfo) error
	DeliverSnapPacket(peer *snap.Peer, packet snap.Packet) error

	Terminate()
	Synchronise(id string, head common.Hash, td *big.Int, mode downloader.SyncMode) error
	Synchronising() bool
	Progress() kaia.SyncProgress
	Cancel()

	GetSnapSyncer() *snap.Syncer
	SyncStakingInfo(id string, from, to uint64) error
	SyncStakingInfoStatus() *downloader.SyncingStatus
}

// ProtocolManagerFetcher is an interface of fetcher.Fetcher used by ProtocolManager.
//
//go:generate mockgen -destination=./mocks/fetcher_mock.go -package=mocks github.com/kaiachain/kaia/node/cn ProtocolManagerFetcher
type ProtocolManagerFetcher interface {
	Enqueue(peer string, block *types.Block) error
	FilterBodies(peer string, transactions [][]*types.Transaction, time time.Time) [][]*types.Transaction
	FilterHeaders(peer string, headers []*types.Header, time time.Time) []*types.Header
	Notify(peer string, hash common.Hash, number uint64, time time.Time, headerFetcher fetcher.HeaderRequesterFn, bodyFetcher fetcher.BodyRequesterFn) error
	Start()
	Stop()
}

// statusData is the network packet for the status message.
type statusData struct {
	ProtocolVersion uint32
	NetworkId       uint64
	TD              *big.Int
	CurrentBlock    common.Hash
	GenesisBlock    common.Hash
	ChainID         *big.Int // ChainID to sign a transaction.
}

// newBlockHashesData is the network packet for the block announcements.
type newBlockHashesData []struct {
	Hash   common.Hash // Hash of one particular block being announced
	Number uint64      // Number of one particular block being announced
}

// getBlockHeadersData represents a block header query.
type getBlockHeadersData struct {
	Origin  hashOrNumber // Block from which to retrieve headers
	Amount  uint64       // Maximum number of headers to retrieve
	Skip    uint64       // Blocks to skip between consecutive headers
	Reverse bool         // Query direction (false = rising towards latest, true = falling towards genesis)
}

// hashOrNumber is a combined field for specifying an origin block.
type hashOrNumber struct {
	Hash   common.Hash // Block hash from which to retrieve headers (excludes Number)
	Number uint64      // Block hash from which to retrieve headers (excludes Hash)
}

// EncodeRLP is a specialized encoder for hashOrNumber to encode only one of the
// two contained union fields.
func (hn *hashOrNumber) EncodeRLP(w io.Writer) error {
	if hn.Hash == (common.Hash{}) {
		return rlp.Encode(w, hn.Number)
	}
	if hn.Number != 0 {
		return fmt.Errorf("both origin hash (%x) and number (%d) provided", hn.Hash, hn.Number)
	}
	return rlp.Encode(w, hn.Hash)
}

// DecodeRLP is a specialized decoder for hashOrNumber to decode the contents
// into either a block hash or a block number.
func (hn *hashOrNumber) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	origin, err := s.Raw()
	if err == nil {
		switch {
		case size == 32:
			err = rlp.DecodeBytes(origin, &hn.Hash)
		case size <= 8:
			err = rlp.DecodeBytes(origin, &hn.Number)
		default:
			err = fmt.Errorf("invalid input size %d for origin", size)
		}
	}
	return err
}

// newBlockData is the network packet for the block propagation message.
type newBlockData struct {
	Block *types.Block
	TD    *big.Int
}

// blockBody represents the data content of a single block.
type blockBody struct {
	Transactions []*types.Transaction // Transactions contained within a block
}

// blockBodiesData is the network packet for block content distribution.
type blockBodiesData []*blockBody
