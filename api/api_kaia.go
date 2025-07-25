// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
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
// This file is derived from internal/ethapi/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package api

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
)

// KaiaAPI provides an API to access Kaia related information.
// It offers only methods that operate on public data that is freely available to anyone.
type KaiaAPI struct {
	b Backend
}

// NewKaiaAPI creates a new Kaia protocol API.
func NewKaiaAPI(b Backend) *KaiaAPI {
	return &KaiaAPI{b}
}

// GasPrice returns a suggestion for a gas price (baseFee * 2).
func (s *KaiaAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	price, err := s.b.SuggestPrice(ctx)
	return (*hexutil.Big)(price), err
}

func (s *KaiaAPI) UpperBoundGasPrice(ctx context.Context) *hexutil.Big {
	return (*hexutil.Big)(s.b.UpperBoundGasPrice(ctx))
}

func (s *KaiaAPI) LowerBoundGasPrice(ctx context.Context) *hexutil.Big {
	return (*hexutil.Big)(s.b.LowerBoundGasPrice(ctx))
}

// ProtocolVersion returns the current Kaia protocol version this node supports.
func (s *KaiaAPI) ProtocolVersion() hexutil.Uint {
	return hexutil.Uint(s.b.ProtocolVersion())
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (s *KaiaAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
	price, err := s.b.SuggestTipCap(ctx)
	return (*hexutil.Big)(price), err
}

type FeeHistoryResult struct {
	OldestBlock  *hexutil.Big     `json:"oldestBlock"`
	Reward       [][]*hexutil.Big `json:"reward,omitempty"`
	BaseFee      []*hexutil.Big   `json:"baseFeePerGas,omitempty"`
	GasUsedRatio []float64        `json:"gasUsedRatio"`
}

// FeeHistory returns data relevant for fee estimation based on the specified range of blocks.
func (s *KaiaAPI) FeeHistory(ctx context.Context, blockCount DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	oldest, reward, baseFee, gasUsed, err := s.b.FeeHistory(ctx, uint64(blockCount), lastBlock, rewardPercentiles)
	if err != nil {
		return nil, err
	}
	results := &FeeHistoryResult{
		OldestBlock:  (*hexutil.Big)(oldest),
		GasUsedRatio: gasUsed,
	}
	if reward != nil {
		results.Reward = make([][]*hexutil.Big, len(reward))
		for i, w := range reward {
			results.Reward[i] = make([]*hexutil.Big, len(w))
			for j, v := range w {
				results.Reward[i][j] = (*hexutil.Big)(v)
			}
		}
	}
	if baseFee != nil {
		results.BaseFee = make([]*hexutil.Big, len(baseFee))
		for i, v := range baseFee {
			results.BaseFee[i] = (*hexutil.Big)(v)
		}
	}
	return results, nil
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (s *KaiaAPI) Syncing() (interface{}, error) {
	progress := s.b.Progress()

	// Return not syncing if the synchronisation already completed
	if progress.CurrentBlock >= progress.HighestBlock {
		return false, nil
	}
	// Otherwise gather the block sync stats
	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(progress.StartingBlock),
		"currentBlock":  hexutil.Uint64(progress.CurrentBlock),
		"highestBlock":  hexutil.Uint64(progress.HighestBlock),
		"pulledStates":  hexutil.Uint64(progress.PulledStates),
		"knownStates":   hexutil.Uint64(progress.KnownStates),
	}, nil
}

// EncodeAccountKey gets an account key of JSON format and returns RLP encoded bytes of the key.
func (s *KaiaAPI) EncodeAccountKey(accKey accountkey.AccountKeyJSON) (hexutil.Bytes, error) {
	if accKey.KeyType == nil {
		return nil, errors.New("key type is not specified")
	}
	key, err := accountkey.NewAccountKey(*accKey.KeyType)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(accKey.Key, key); err != nil {
		return nil, err
	}
	// Invalidate zero values of threshold and weight to prevent users' mistake
	// JSON unmarshalling sets zero for those values if they are not exist on JSON input
	if err := checkAccountKeyZeroValues(key, false); err != nil {
		return nil, err
	}
	accKeySerializer := accountkey.NewAccountKeySerializerWithAccountKey(key)
	encodedKey, err := rlp.EncodeToBytes(accKeySerializer)
	if err != nil {
		return nil, errors.New("the key probably contains an invalid public key: " + err.Error())
	}
	return (hexutil.Bytes)(encodedKey), nil
}

// DecodeAccountKey gets an RLP encoded bytes of an account key and returns the decoded account key.
func (s *KaiaAPI) DecodeAccountKey(encodedAccKey hexutil.Bytes) (*accountkey.AccountKeySerializer, error) {
	dec := accountkey.NewAccountKeySerializer()
	if err := rlp.DecodeBytes(encodedAccKey, &dec); err != nil {
		return nil, err
	}
	return dec, nil
}

// checkAccountKeyZeroValues returns errors if the input account key contains zero values of threshold or weight.
func checkAccountKeyZeroValues(key accountkey.AccountKey, isNested bool) error {
	switch key.Type() {
	case accountkey.AccountKeyTypeWeightedMultiSig:
		multiSigKey, _ := key.(*accountkey.AccountKeyWeightedMultiSig)
		if multiSigKey.Threshold == 0 {
			return errors.New("invalid threshold of the multiSigKey")
		}
		for _, weightedKey := range multiSigKey.Keys {
			if weightedKey.Weight == 0 {
				return errors.New("invalid weight of the multiSigKey")
			}
		}
	case accountkey.AccountKeyTypeRoleBased:
		if isNested {
			return errors.New("roleBasedKey cannot contains a roleBasedKey as a role key")
		}
		roleBasedKey, _ := key.(*accountkey.AccountKeyRoleBased)
		for _, roleKey := range *roleBasedKey {
			if err := checkAccountKeyZeroValues(roleKey, true); err != nil {
				return err
			}
		}
	}
	return nil
}
