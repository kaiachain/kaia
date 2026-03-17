package impl

import (
	"bytes"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/kaiax/randao"
	"github.com/kaiachain/kaia/params"
)

// Calculate KIP-114 Randao header fields
// https://github.com/klaytn/kips/blob/kip114/KIPs/kip-114.md
func (r *RandaoModule) CalcRandao(number *big.Int, prevMixHash []byte) ([]byte, []byte, error) {
	if r.BlsSecretKey == nil {
		return nil, nil, randao.ErrNoBlsKey
	}
	if len(prevMixHash) != 32 {
		logger.Error("invalid prevMixHash", "number", number.Uint64(), "prevMixHash", hexutil.Encode(prevMixHash))
		return nil, nil, randao.ErrInvalidRandaoFields
	}

	// block_num_to_bytes() = num.to_bytes(32, byteorder="big")
	msg := calcRandaoMsg(number)

	// calc_random_reveal() = sign(privateKey, headerNumber)
	randomReveal := bls.Sign(r.BlsSecretKey, msg[:]).Marshal()

	// calc_mix_hash() = xor(prevMixHash, keccak256(randomReveal))
	mixHash := calcMixHash(randomReveal, prevMixHash)

	return randomReveal, mixHash, nil
}

func (r *RandaoModule) VerifyHeader(header *types.Header, parent *types.Header) error {
	if header.Number.Sign() == 0 {
		return nil // Do not verify genesis block
	}

	if !r.ChainConfig.IsRandaoForkEnabled(header.Number) {
		if header.RandomReveal != nil || header.MixHash != nil {
			return randao.ErrUnexpectedRandao
		}
		return nil
	}

	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	proposer, err := r.Chain.Engine().Author(header)
	if err != nil {
		return err
	}

	// [proposerPubkey, proposerPop] = get_proposer_pubkey_pop()
	// if not pop_verify(proposerPubkey, proposerPop): return False
	proposerPub, err := r.GetBlsPubkey(proposer, header.Number)
	if err != nil {
		return err
	}

	// if not verify(proposerPubkey, newHeader.number, newHeader.randomReveal): return False
	msg := calcRandaoMsg(header.Number)
	ok, err := bls.VerifySignature(header.RandomReveal, msg, proposerPub)
	if err != nil {
		return err
	}
	if !ok {
		return randao.ErrInvalidRandaoFields
	}

	prevMixHash := headerMixHash(r.ChainConfig, parent)
	if len(prevMixHash) != 32 {
		return randao.ErrInvalidRandaoFields
	}
	mixHash := calcMixHash(header.RandomReveal, prevMixHash)
	if !bytes.Equal(header.MixHash, mixHash) {
		return randao.ErrInvalidRandaoFields
	}

	return nil
}

func (r *RandaoModule) PrepareHeader(header *types.Header) error {
	if !r.ChainConfig.IsRandaoForkEnabled(header.Number) {
		return nil
	}

	parent := r.Chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	prevMixHash := headerMixHash(r.ChainConfig, parent)

	randomReveal, mixHash, err := r.CalcRandao(header.Number, prevMixHash)
	if err != nil {
		return err
	}

	header.RandomReveal = randomReveal
	header.MixHash = mixHash
	return nil
}

// block_num_to_bytes() = num.to_bytes(32, byteorder="big")
func calcRandaoMsg(number *big.Int) common.Hash {
	return common.BytesToHash(number.Bytes())
}

// calc_mix_hash() = xor(prevMixHash, keccak256(randomReveal))
func calcMixHash(randomReveal, prevMixHash []byte) []byte {
	mixHash := make([]byte, 32)
	revealHash := crypto.Keccak256(randomReveal)
	for i := 0; i < 32; i++ {
		mixHash[i] = prevMixHash[i] ^ revealHash[i]
	}
	return mixHash
}

// At the fork block's parent, pretend that prevMixHash is ZeroMixHash.
func headerMixHash(chainConfig *params.ChainConfig, header *types.Header) []byte {
	if chainConfig.IsRandaoForkBlockParent(header.Number) {
		return params.ZeroMixHash
	}
	return header.MixHash
}
