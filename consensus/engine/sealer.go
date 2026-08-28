package engine

import (
	"crypto/ecdsa"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
)

var _ consensus.Sealer = (*dynamicSealer)(nil)

func newSealerImpl(privateKey *ecdsa.PrivateKey) consensus.Sealer {
	return istanbul.NewSealerImpl(privateKey)
}

type dynamicSealer struct {
	chainConfig *params.ChainConfig
	address     common.Address

	sealer consensus.Sealer
}

func NewSealer(chainConfig *params.ChainConfig, privateKey *ecdsa.PrivateKey) consensus.Sealer {
	s := &dynamicSealer{
		chainConfig: chainConfig,
	}
	if privateKey != nil {
		s.address = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
	s.sealer = newSealerImpl(privateKey)
	types.SetHeaderHashFn(s.HeaderHash)
	return s
}

func (s *dynamicSealer) Address() common.Address {
	return s.address
}

func (s *dynamicSealer) MakeCommittedSeal(header *types.Header) ([]byte, error) {
	if header == nil || header.Number == nil {
		return nil, consensus.ErrUnknownBlock
	}
	// If it's genesis block, return empty seal
	if header.Number.Sign() == 0 {
		return []byte{}, nil
	}
	return s.implAt(header.Number.Uint64()).MakeCommittedSeal(header)
}

func (s *dynamicSealer) MakeAuthorSeal(header *types.Header) ([]byte, error) {
	if header == nil || header.Number == nil {
		return nil, consensus.ErrUnknownBlock
	}
	// If it's genesis block, return empty seal
	if header.Number.Sign() == 0 {
		return []byte{}, nil
	}
	return s.implAt(header.Number.Uint64()).MakeAuthorSeal(header)
}

func (s *dynamicSealer) Author(header *types.Header) (common.Address, error) {
	if header == nil || header.Number == nil {
		return common.Address{}, consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).Author(header)
}

func (s *dynamicSealer) Committers(header *types.Header) ([]common.Address, error) {
	if header == nil || header.Number == nil {
		return nil, consensus.ErrUnknownBlock
	}
	// Post-permissionless committed seals bind the round; recover with the matching preimage.
	if s.chainConfig.IsPermissionlessForkEnabled(header.Number) {
		return s.implAt(header.Number.Uint64()).CommittersWithRound(header)
	}
	return s.implAt(header.Number.Uint64()).Committers(header)
}

func (s *dynamicSealer) CommittersWithRound(header *types.Header) ([]common.Address, error) {
	if header == nil || header.Number == nil {
		return nil, consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).CommittersWithRound(header)
}

func (s *dynamicSealer) RecoverCommitters(hash common.Hash, round byte, seals [][]byte) ([]common.Address, error) {
	return s.implAt(0).RecoverCommitters(hash, round, seals) // implAt ignores the number
}

func (s *dynamicSealer) Vanity(header *types.Header) ([]byte, error) {
	if header == nil || header.Number == nil {
		return nil, consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).Vanity(header)
}

func (s *dynamicSealer) Validators(header *types.Header) ([]common.Address, error) {
	if header == nil || header.Number == nil {
		return nil, consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).Validators(header)
}

func (s *dynamicSealer) Round(header *types.Header) (byte, error) {
	if header == nil || header.Number == nil {
		return 0, consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).Round(header)
}

func (s *dynamicSealer) RawSeals(header *types.Header) ([]byte, [][]byte, error) {
	if header == nil || header.Number == nil {
		return nil, nil, consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).RawSeals(header)
}

func (s *dynamicSealer) WriteRound(header *types.Header, round int64) {
	if header == nil || header.Number == nil {
		return
	}
	s.implAt(header.Number.Uint64()).WriteRound(header, round)
}

func (s *dynamicSealer) WriteValidators(header *types.Header, validators []common.Address) error {
	if header == nil || header.Number == nil {
		return consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).WriteValidators(header, validators)
}

func (s *dynamicSealer) WriteAuthorSeal(header *types.Header, seal []byte) error {
	if header == nil || header.Number == nil {
		return consensus.ErrUnknownBlock
	}
	// It shoudn't support writing author seal for the genesis block
	if header.Number.Sign() == 0 {
		return consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).WriteAuthorSeal(header, seal)
}

func (s *dynamicSealer) WriteCommittedSeals(header *types.Header, committedSeals [][]byte) error {
	if header == nil || header.Number == nil {
		return consensus.ErrUnknownBlock
	}
	// It shouldn't support writing committed seals for the genesis block
	if header.Number.Sign() == 0 {
		return consensus.ErrUnknownBlock
	}
	return s.implAt(header.Number.Uint64()).WriteCommittedSeals(header, committedSeals)
}

func (s *dynamicSealer) F(blockNum uint64, qualifiedlen, committeesize int) int {
	return s.implAt(blockNum).F(blockNum, qualifiedlen, committeesize)
}

func (s *dynamicSealer) Quorum(blockNum uint64, qualifiedlen, committeesize int) int {
	return s.implAt(blockNum).Quorum(blockNum, qualifiedlen, committeesize)
}

func (s *dynamicSealer) SigHash(header *types.Header) common.Hash {
	if header == nil || header.Number == nil {
		return common.Hash{}
	}
	return s.implAt(header.Number.Uint64()).SigHash(header)
}

func (s *dynamicSealer) HeaderHash(header *types.Header) common.Hash {
	if header == nil || header.Number == nil {
		return istanbul.RLPHash(header)
	}
	return s.implAt(header.Number.Uint64()).HeaderHash(header)
}

func (s *dynamicSealer) implAt(number uint64) consensus.Sealer {
	_ = number
	return s.sealer
}
