// Modifications Copyright 2023 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package nodecmd

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/kaiachain/kaia/accounts/keystore"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/urfave/cli/v2"
)

const (
	DECODE_EXTRA = "decode-extra"
	DECODE_VOTE  = "decode-vote"
	DECODE_GOV   = "decode-gov"
	DECRYPT_KEY  = "decrypt-keystore"
)

var ErrInvalidCmd = errors.New("Invalid command. Check usage through --help command")

var UtilCommand = &cli.Command{
	Name:     "util",
	Usage:    "offline utility",
	Category: "MISCELLANEOUS COMMANDS",
	Subcommands: []*cli.Command{
		{
			Name:        DECODE_EXTRA,
			Usage:       "<header file (json format)>",
			Action:      action,
			Description: "Decode header extra field",
		},
		{
			Name:        DECODE_VOTE,
			Usage:       "<hex bytes>",
			Action:      action,
			Description: "Decode header vote field",
		},
		{
			Name:        DECODE_GOV,
			Usage:       "<hex bytes>",
			Action:      action,
			Description: "Decode header governance field",
		},
		{
			Name:        DECRYPT_KEY,
			Usage:       "<keystore path> <password>",
			Action:      action,
			Description: "Decrypt keystore",
		},
	},
}

func action(ctx *cli.Context) error {
	var (
		m   map[string]interface{}
		err error
	)
	switch ctx.Command.Name {
	case DECODE_EXTRA:
		if ctx.Args().Len() != 1 {
			return ErrInvalidCmd
		}
		m, err = decodeExtra(ctx.Args().Get(0))
	case DECODE_VOTE:
		if ctx.Args().Len() != 1 {
			return ErrInvalidCmd
		}
		m, err = decodeVote(hex2Bytes(ctx.Args().Get(0)))
	case DECODE_GOV:
		if ctx.Args().Len() != 1 {
			return ErrInvalidCmd
		}
		m, err = decodeGov(hex2Bytes(ctx.Args().Get(0)))
	case DECRYPT_KEY:
		if ctx.Args().Len() != 2 {
			return ErrInvalidCmd
		}
		keystorePath, passwd := ctx.Args().Get(0), ctx.Args().Get(1)
		m, err = extractKeypair(keystorePath, passwd)
	default:
		return ErrInvalidCmd
	}
	if err == nil {
		prettyPrint(m)
	}
	return err
}

func hex2Bytes(s string) []byte {
	if data, err := hexutil.Decode(s); err == nil {
		return data
	} else {
		panic(err)
	}
}

func prettyPrint(m map[string]interface{}) {
	if b, err := json.MarshalIndent(m, "", "  "); err == nil {
		fmt.Println(string(b))
	} else {
		panic(err)
	}
}

func extractKeypair(keystorePath, passwd string) (map[string]interface{}, error) {
	keyjson, err := os.ReadFile(keystorePath)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(keyjson, passwd)
	if err != nil {
		return nil, err
	}
	addr := key.GetAddress().String()
	pubkey := key.GetPrivateKey().PublicKey
	privkey := key.GetPrivateKey()
	m := make(map[string]interface{})
	m["addr"] = addr
	m["privkey"] = hex.EncodeToString(crypto.FromECDSA(privkey))
	m["pubkey"] = hex.EncodeToString(crypto.FromECDSAPub(&pubkey))
	return m, nil
}

func decodeGov(bytes []byte) (map[string]interface{}, error) {
	var b []byte
	m := make(map[string]interface{})
	if err := rlp.DecodeBytes(bytes, &b); err == nil {
		if err := json.Unmarshal(b, &m); err == nil {
			return m, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func parseHeaderFile(headerFile string) (*types.Header, error) {
	header := new(types.Header)
	bytes, err := os.ReadFile(headerFile)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &header); err != nil {
		return nil, err
	}
	return header, nil
}

func decodeExtra(headerFile string) (map[string]interface{}, error) {
	header, err := parseHeaderFile(headerFile)
	if err != nil {
		return nil, err
	}
	// decode-extra can use a fixed sealer config because only one effective
	// scheme implementation is available. If a new scheme is introduced later,
	// this command should inspect the chainConfig from DB and reintroduce a
	// scheme option so DecodeExtra uses the correct implementation.
	s := bft.NewSealer(params.TestChainConfig.Copy(), nil)
	validators, err := s.Validators(header)
	if err != nil {
		return nil, err
	}
	authorSeal, committedSeals, err := s.RawSeals(header)
	if err != nil {
		return nil, err
	}
	round, err := s.Round(header)
	if err != nil {
		return nil, err
	}
	author, err := s.Author(header)
	if err != nil {
		author = common.Address{}
	}
	committers, err := s.Committers(header)
	if err != nil {
		committers = []common.Address{}
	}

	headerHash := s.HeaderHash(header)
	sigHash := s.SigHash(header)

	result := make(map[string]interface{})
	result["hash"] = headerHash.Hex()
	result["sigHash"] = sigHash.Hex()
	result["validators"] = validators
	result["seal"] = hexutil.Encode(authorSeal)
	result["committedSeal"] = committedSeals
	result["committers"] = committers
	result["validatorSize"] = len(validators)
	result["committedSealSize"] = len(committedSeals)
	result["proposer"] = author.String()
	result["round"] = round
	return result, nil
}

func decodeVote(bytes []byte) (map[string]interface{}, error) {
	var vb headergov.VoteBytes = bytes
	vote, err := vb.ToVoteData()
	if err != nil {
		return nil, err
	}

	m := make(map[string]interface{})
	m["validator"] = vote.Voter()
	m["key"] = vote.Name()
	m["value"] = vote.Value()
	return m, nil
}
