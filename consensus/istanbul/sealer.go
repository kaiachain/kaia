package istanbul

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"io"
	"math"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/rlp"
)

var _ consensus.Sealer = (*IstanbulSealer)(nil)

var (
	inmemoryBlocks             = 2048
	inmemoryValidatorsPerBlock = 30
	signatureAddresses, _      = lru.NewARC(inmemoryBlocks * inmemoryValidatorsPerBlock)

	// IstanbulDigest represents a hash of "Istanbul practical byzantine fault tolerance"
	// to identify whether the block is from Istanbul consensus engine.
	IstanbulDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	IstanbulExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	IstanbulExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	// ErrInvalidIstanbulHeaderExtra is returned if the length of extra-data is less than 32 bytes.
	ErrInvalidIstanbulHeaderExtra = errors.New("invalid istanbul header extra-data")
)

type IstanbulSealer struct {
	privateKey *ecdsa.PrivateKey

	parsedExtraCache *lru.ARCCache
}

func NewSealerImpl(privateKey *ecdsa.PrivateKey) *IstanbulSealer {
	parsedExtraCache, _ := lru.NewARC(inmemoryBlocks)
	return &IstanbulSealer{privateKey: privateKey, parsedExtraCache: parsedExtraCache}
}

type IstanbulExtra struct {
	Validators    []common.Address
	Seal          []byte
	CommittedSeal [][]byte
}

// EncodeRLP serializes the istanbul fields into the Kaia RLP format.
func (ist *IstanbulExtra) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{
		ist.Validators,
		ist.Seal,
		ist.CommittedSeal,
	})
}

// DecodeRLP implements rlp.Decoder and loads the istanbul fields from a RLP stream.
func (ist *IstanbulExtra) DecodeRLP(s *rlp.Stream) error {
	var istanbulExtra struct {
		Validators    []common.Address
		Seal          []byte
		CommittedSeal [][]byte
	}
	if err := s.Decode(&istanbulExtra); err != nil {
		return err
	}
	ist.Validators, ist.Seal, ist.CommittedSeal = istanbulExtra.Validators, istanbulExtra.Seal, istanbulExtra.CommittedSeal
	return nil
}

func (ist *IstanbulExtra) Copy() *IstanbulExtra {
	out := &IstanbulExtra{
		Validators: make([]common.Address, len(ist.Validators)),
		Seal:       append([]byte(nil), ist.Seal...),
	}
	copy(out.Validators, ist.Validators)
	out.CommittedSeal = make([][]byte, len(ist.CommittedSeal))
	for i, seal := range ist.CommittedSeal {
		out.CommittedSeal[i] = append([]byte(nil), seal...)
	}
	return out
}

func ensureIstanbulVanityBytes(extra []byte) []byte {
	if len(extra) >= IstanbulExtraVanity {
		return append([]byte(nil), extra...)
	}
	out := append([]byte(nil), extra...)
	out = append(out, bytes.Repeat([]byte{0x00}, IstanbulExtraVanity-len(out))...)
	return out
}

func (m *IstanbulSealer) parseIstanbulExtra(extraBytes []byte) (*IstanbulExtra, error) {
	if len(extraBytes) < IstanbulExtraVanity {
		return nil, ErrInvalidIstanbulHeaderExtra
	}

	if m.parsedExtraCache != nil {
		key := crypto.Keccak256Hash(extraBytes)
		if cached, ok := m.parsedExtraCache.Get(key); ok {
			return cached.(*IstanbulExtra).Copy(), nil
		}
	}

	var extra *IstanbulExtra
	if err := rlp.DecodeBytes(extraBytes[IstanbulExtraVanity:], &extra); err != nil {
		return nil, err
	}
	if m.parsedExtraCache != nil {
		m.parsedExtraCache.Add(crypto.Keccak256Hash(extraBytes), extra.Copy())
	}
	return extra, nil
}

func istanbulFilteredHeader(h *types.Header, keepSeal bool) *types.Header {
	newHeader := types.CopyHeader(h)
	if len(newHeader.Extra) < IstanbulExtraVanity {
		return nil
	}
	var istanbulExtra *IstanbulExtra
	err := rlp.DecodeBytes(newHeader.Extra[IstanbulExtraVanity:], &istanbulExtra)
	if err != nil {
		return nil
	}

	if !keepSeal {
		istanbulExtra.Seal = []byte{}
	}
	istanbulExtra.CommittedSeal = [][]byte{}

	payload, err := rlp.EncodeToBytes(istanbulExtra)
	if err != nil {
		return nil
	}
	newHeader.Extra[IstanbulExtraVanity-1] = 0
	newHeader.Extra = append(newHeader.Extra[:IstanbulExtraVanity], payload...)
	return newHeader
}

func prepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(2)}) // keep in sync with bft.MsgCommit
	return buf.Bytes()
}

func (m *IstanbulSealer) Author(header *types.Header) (common.Address, error) {
	extra, err := m.parseIstanbulExtra(header.Extra)
	if err != nil {
		return common.Address{}, err
	}
	if len(extra.Seal) == 0 {
		return common.Address{}, ErrInvalidSignature
	}
	addr, err := cacheSignatureAddress(m.SigHash(header).Bytes(), extra.Seal)
	if err != nil {
		return common.Address{}, ErrInvalidSignature
	}
	return addr, nil
}

func (m *IstanbulSealer) Committers(header *types.Header) ([]common.Address, error) {
	extra, err := m.parseIstanbulExtra(header.Extra)
	if err != nil {
		return nil, err
	}
	if len(extra.CommittedSeal) == 0 {
		return []common.Address{}, nil
	}

	proposalSeal := prepareCommittedSeal(m.HeaderHash(header)) // bft.MsgCommit
	committers := make([]common.Address, 0, len(extra.CommittedSeal))
	for _, seal := range extra.CommittedSeal {
		addr, err := cacheSignatureAddress(proposalSeal, seal)
		if err != nil {
			return nil, ErrInvalidSignature
		}
		committers = append(committers, addr)
	}
	return committers, nil
}

func (m *IstanbulSealer) Vanity(header *types.Header) ([]byte, error) {
	if _, err := m.parseIstanbulExtra(header.Extra); err != nil {
		return nil, err
	}
	vanity := make([]byte, IstanbulExtraVanity)
	copy(vanity, header.Extra[:IstanbulExtraVanity])
	return vanity, nil
}

func (m *IstanbulSealer) Validators(header *types.Header) ([]common.Address, error) {
	extra, err := m.parseIstanbulExtra(header.Extra)
	if err != nil {
		return nil, err
	}
	validators := make([]common.Address, len(extra.Validators))
	copy(validators, extra.Validators)
	return validators, nil
}

func (m *IstanbulSealer) Round(header *types.Header) (byte, error) {
	if _, err := m.parseIstanbulExtra(header.Extra); err != nil {
		return 0, err
	}
	return header.Extra[IstanbulExtraVanity-1], nil
}

func (m *IstanbulSealer) RawSeals(header *types.Header) ([]byte, [][]byte, error) {
	extra, err := m.parseIstanbulExtra(header.Extra)
	if err != nil {
		return nil, nil, err
	}
	authorSeal := append([]byte(nil), extra.Seal...)
	committedSeals := make([][]byte, len(extra.CommittedSeal))
	for i, seal := range extra.CommittedSeal {
		committedSeals[i] = append([]byte(nil), seal...)
	}
	return authorSeal, committedSeals, nil
}

func (m *IstanbulSealer) writeExtra(header *types.Header, mutator func(*IstanbulExtra) error) error {
	header.Extra = ensureIstanbulVanityBytes(header.Extra)

	extra, err := m.parseIstanbulExtra(header.Extra)
	if err != nil {
		extra = &IstanbulExtra{}
	}

	if err := mutator(extra); err != nil {
		return err
	}

	payload, err := rlp.EncodeToBytes(extra)
	if err != nil {
		return err
	}
	header.Extra = append(header.Extra[:IstanbulExtraVanity], payload...)
	return nil
}

func (m *IstanbulSealer) WriteRound(header *types.Header, round int64) {
	header.Extra = ensureIstanbulVanityBytes(header.Extra)
	header.Extra[IstanbulExtraVanity-1] = byte(round)
}

func (m *IstanbulSealer) WriteValidators(header *types.Header, validators []common.Address) error {
	return m.writeExtra(header, func(extra *IstanbulExtra) error {
		extra.Validators = append([]common.Address(nil), validators...)
		// Validators are written when preparing a new proposal header.
		// Any inherited proposer/committed seals from a parent header must be cleared.
		extra.Seal = nil
		extra.CommittedSeal = nil
		return nil
	})
}

func (m *IstanbulSealer) WriteAuthorSeal(header *types.Header, seal []byte) error {
	if len(seal) != IstanbulExtraSeal {
		return ErrInvalidSignature
	}
	return m.writeExtra(header, func(extra *IstanbulExtra) error {
		extra.Seal = append([]byte(nil), seal...)
		return nil
	})
}

func (m *IstanbulSealer) WriteCommittedSeals(header *types.Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return ErrInvalidCommittedSeals
	}
	return m.writeExtra(header, func(extra *IstanbulExtra) error {
		extra.CommittedSeal = make([][]byte, len(committedSeals))
		for i, seal := range committedSeals {
			if len(seal) != IstanbulExtraSeal {
				return ErrInvalidCommittedSeals
			}
			extra.CommittedSeal[i] = append([]byte(nil), seal...)
		}
		return nil
	})
}

func (m *IstanbulSealer) HeaderHash(header *types.Header) common.Hash {
	if filtered := istanbulFilteredHeader(header, true); filtered != nil {
		return RLPHash(filtered)
	}
	return RLPHash(header)
}

func (m *IstanbulSealer) SigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()
	rlp.Encode(hasher, istanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

func (m *IstanbulSealer) F(_ uint64, qualifiedLen, committeeSize int) int {
	if qualifiedLen > int(committeeSize) {
		return int(math.Ceil(float64(committeeSize)/3)) - 1
	}
	return int(math.Ceil(float64(qualifiedLen)/3)) - 1
}

func (m *IstanbulSealer) Quorum(blockNum uint64, qualifiedlen, committeeSize int) int {
	return 2*m.F(blockNum, qualifiedlen, committeeSize) + 1
}

func (m *IstanbulSealer) MakeCommittedSeal(header *types.Header) ([]byte, error) {
	return crypto.Sign(crypto.Keccak256(prepareCommittedSeal(m.HeaderHash(header))), m.privateKey)
}

func (m *IstanbulSealer) MakeAuthorSeal(header *types.Header) ([]byte, error) {
	if header.Number.Sign() == 0 {
		return []byte{}, nil
	}
	canonicalHeader := types.CopyHeader(header)
	if err := m.writeExtra(canonicalHeader, func(*IstanbulExtra) error { return nil }); err != nil {
		return nil, err
	}
	return crypto.Sign(crypto.Keccak256(m.SigHash(canonicalHeader).Bytes()), m.privateKey)
}

func cacheSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	sigHex := hex.EncodeToString(sig)
	if addr, ok := signatureAddresses.Get(sigHex); ok {
		return addr.(common.Address), nil
	}
	addr, err := GetSignatureAddress(data, sig)
	if err != nil {
		return common.Address{}, err
	}
	signatureAddresses.Add(sigHex, addr)
	return addr, nil
}
