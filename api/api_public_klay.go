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

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/rlp"
)

// PublicKaiaAPI provides an API to access Kaia related information.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicKaiaAPI struct {
	b Backend
}

// NewPublicKaiaAPI creates a new Kaia protocol API.
func NewPublicKaiaAPI(b Backend) *PublicKaiaAPI {
	return &PublicKaiaAPI{b}
}

// GasPrice returns a suggestion for a gas price (baseFee * 2).
func (s *PublicKaiaAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	price, err := s.b.SuggestPrice(ctx)
	return (*hexutil.Big)(price), err
}

func (s *PublicKaiaAPI) UpperBoundGasPrice(ctx context.Context) *hexutil.Big {
	return (*hexutil.Big)(s.b.UpperBoundGasPrice(ctx))
}

func (s *PublicKaiaAPI) LowerBoundGasPrice(ctx context.Context) *hexutil.Big {
	return (*hexutil.Big)(s.b.LowerBoundGasPrice(ctx))
}

// ProtocolVersion returns the current Kaia protocol version this node supports.
func (s *PublicKaiaAPI) ProtocolVersion() hexutil.Uint {
	return hexutil.Uint(s.b.ProtocolVersion())
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (s *PublicKaiaAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
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
func (s *PublicKaiaAPI) FeeHistory(ctx context.Context, blockCount DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	oldest, reward, baseFee, gasUsed, err := s.b.FeeHistory(ctx, int(blockCount), lastBlock, rewardPercentiles)
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

type TotalSupplyResult struct {
	Number      *hexutil.Big `json:"number"`          // Block number in which the total supply was calculated.
	Error       *string      `json:"error,omitempty"` // Errors that occurred while fetching the components, thus failed to deliver the total supply.
	TotalSupply *hexutil.Big `json:"totalSupply"`     // The total supply of the native token. i.e. Minted - Burnt.
	TotalMinted *hexutil.Big `json:"totalMinted"`     // Total minted amount.
	TotalBurnt  *hexutil.Big `json:"totalBurnt"`      // Total burnt amount. Sum of all burnt amounts below.
	BurntFee    *hexutil.Big `json:"burntFee"`        // from tx fee burn. ReadAccReward(num).BurntFee.
	ZeroBurn    *hexutil.Big `json:"zeroBurn"`        // balance of 0x0 (zero) address.
	DeadBurn    *hexutil.Big `json:"deadBurn"`        // balance of 0xdead (dead) address.
	Kip103Burn  *hexutil.Big `json:"kip103Burn"`      // by KIP103 fork. Read from its memo.
	Kip160Burn  *hexutil.Big `json:"kip160Burn"`      // by KIP160 fork. Read from its memo.
}

func (s *PublicKaiaAPI) GetTotalSupply(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*TotalSupplyResult, error) {
	block, err := s.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}

	// Case 1. Failed to fetch essential components. The API fails.
	ts, err := s.b.GetTotalSupply(ctx, blockNrOrHash)
	if ts == nil {
		return nil, err
	}

	// Case 2. Failed to fetch some components. The API delivers the partial result with the 'error' field set.
	// Case 3. Succeeded to fetch all components. The API delivers the full result.
	res := &TotalSupplyResult{
		Number:      (*hexutil.Big)(block.Number()),
		Error:       nil,
		TotalSupply: (*hexutil.Big)(ts.TotalSupply),
		TotalMinted: (*hexutil.Big)(ts.TotalMinted),
		TotalBurnt:  (*hexutil.Big)(ts.TotalBurnt),
		BurntFee:    (*hexutil.Big)(ts.BurntFee),
		ZeroBurn:    (*hexutil.Big)(ts.ZeroBurn),
		DeadBurn:    (*hexutil.Big)(ts.DeadBurn),
		Kip103Burn:  (*hexutil.Big)(ts.Kip103Burn),
		Kip160Burn:  (*hexutil.Big)(ts.Kip160Burn),
	}
	if err != nil {
		errStr := err.Error()
		res.Error = &errStr
	}
	return res, nil
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (s *PublicKaiaAPI) Syncing() (interface{}, error) {
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
func (s *PublicKaiaAPI) EncodeAccountKey(accKey accountkey.AccountKeyJSON) (hexutil.Bytes, error) {
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
func (s *PublicKaiaAPI) DecodeAccountKey(encodedAccKey hexutil.Bytes) (*accountkey.AccountKeySerializer, error) {
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
