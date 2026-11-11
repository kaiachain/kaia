package engine

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockHash(t *testing.T) {
	cases := []struct {
		name         string
		encoded      string
		expectedHash common.Hash
	}{
		{
			name:         "kairos-144654443", // Kairos 144654443
			encoded:      "f904f1f903c1a046f5775da65a7c9e9a12d2ee4d010285df656ebd64c97cad2669c42c4aa39d8b94f90675a56a03f836204d66c0f923e00500ddc90aa0cec442b90b1d21e465474f7e5c233c3461e950bf383b39227917ad8a608653dfa004de211f466c3543bb2e3db8c976e9d046c214bc33c982d664b7e494d0c2c88aa0390d7c4c79e1e594bcc3236742112d9a7358652e66b1c97dea1c490e69276617b90100000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000002000000000000000000800000000000000000000010000020000000000000010000000000000000000000000000000000000006000000000000000000000004000000000000001240004000000000000000000000000000000000000004800000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000002000000000000000000000000000000000000000400000a00000000000000000000000000000000000000000000184089f406b8302c7158465b315d945b90187d883010c00846b6c617988676f312e32302e36856c696e757800000000000000f90164f85494571e53df607be97431a5bbefca1dffe5aef56f4d945cb1a7dccbd0dc446e3640898ede8820368554c89499fb17d324fa0e07f23b49d09028ac0919414db694b74ff9dea397fe9e231df545eb53fe2adf776cb2b841dc7dca1dc15d06e83a0a39217052ad7915d3b0f0a68f49d4b53498130a1212327b9dd5514362eb3e5f0ebb277e7d38baecad9cd045d3627d3a1f322422eeee8a01f8c9b8410a3fb9227632261c8520f1ee859c503f017701fdb92c2e39f532e190779888a63cf112cf2e530f656153ab64103558c26fb9ca43d5cf31425e62569f891e1eca00b8414a32097b87c481f145ae75ded4e8ec2ea4230c7b2feb3f3afcd2d22ea0b7d4596f4354e185d4f1f93b4d758ddb97e9203e1d306413ad0db1f4a779156879baa300b841fde66a2abb916bf314c82a2dfcd82ed356af829a98573d7cdb4a74cccb6ceecc22358750332ad31dccd9a890dbde7e0d4e21f173515011656f6c222169ba94620080808505d21dba00b860a027d1591d784686d063234bc000a8b9f829da534c41b555c867329b2b4abc5ef8f36245546adcc35c31dbc7c16b084a0d48487817b0eb191ac5b9aa8f3c908d413f19a1e784c3920fc0c0ef9180582902dc30ee8b70859e493af467bc227403a047f8b901fe5f4a3387cb6ccd5c14a00f7649247d2d2d315a652c378f328c1970f9012a31f90126830a596a850ba43b740083061a8094a317038414a275365ed4a085b786e83e761d20a580944dd5fa43054709155d5e68fadd3aeb3f853f34b0b844202ee0ed00000000000000000000000000000000000000000000000000000000000a596a0000000000000000000000000000000000000000000000000000000000000377f847f8458207f5a085d89c94a1e010f90b8e825e5321b1fff91ea748f50424448e17d5f1658498bda053936941c152ce9721447711a80c39b07fb36d42c8a8b6a258d8b071ab3b6c96945e6b99bca5a21818d40d12c56194674989146fc8f847f8458207f5a0c40d29e71cdd4f07051398fa35d41a25cc2f5db165d3ad8aa8b5eff20cc81996a01ea154d91320586165b139870e688a0ab4837137c5208a2da60cbdc60d7c7d81",
			expectedHash: common.HexToHash("0x4d0a7bafd8a88c983431031270df6bfc0e9a0909b126f2a3f1ac5b3c03392eb7"),
		},
	}

	s := NewSealer(params.KairosChainConfig, nil)
	defer types.SetHeaderHashFn(nil)
	for _, tc := range cases {
		var block types.Block
		require.NoError(t, rlp.DecodeBytes(common.FromHex(tc.encoded), &block), tc.name)

		assert.Equal(t, tc.expectedHash, block.Hash(), tc.name)
		assert.Equal(t, tc.expectedHash, block.Header().Hash(), tc.name)
		assert.Equal(t, tc.expectedHash, s.HeaderHash(block.Header()), tc.name)
	}
}

func TestBlockHash_PreIstanbulCompatibleBlockUsesIstanbulHeaderHash(t *testing.T) {
	header := &types.Header{
		ParentHash:  common.HexToHash("0x1"),
		Rewardbase:  common.HexToAddress("0x2"),
		Root:        common.HexToHash("0x3"),
		TxHash:      common.HexToHash("0x4"),
		ReceiptHash: common.HexToHash("0x5"),
		BlockScore:  big.NewInt(1),
		Number:      big.NewInt(6145),
		GasUsed:     21000,
		Time:        big.NewInt(1711959900),
		Governance:  []byte{0x01},
		Vote:        []byte{0x02},
	}

	istanbulSealer := istanbul.NewSealerImpl(nil)
	require.NoError(t, istanbulSealer.WriteValidators(header, []common.Address{
		common.HexToAddress("0x44add0ec310f115a0e603b2d7db9f067778eaf8a"),
		common.HexToAddress("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212"),
	}))
	istanbulSealer.WriteRound(header, 3)
	require.NoError(t, istanbulSealer.WriteAuthorSeal(header, bytes.Repeat([]byte{0x07}, istanbul.IstanbulExtraSeal)))
	require.NoError(t, istanbulSealer.WriteCommittedSeals(header, [][]byte{bytes.Repeat([]byte{0x09}, istanbul.IstanbulExtraSeal)}))

	plainHash := istanbul.RLPHash(header)
	expectedHash := istanbulSealer.HeaderHash(header)
	require.NotEqual(t, plainHash, expectedHash)

	s := NewSealer(params.KairosChainConfig, nil)
	defer types.SetHeaderHashFn(nil)

	assert.Equal(t, expectedHash, s.HeaderHash(header))
	assert.Equal(t, expectedHash, header.Hash())
	assert.Equal(t, uint64(75373312), params.KairosChainConfig.IstanbulCompatibleBlock.Uint64())
	assert.Less(t, header.Number.Uint64(), params.KairosChainConfig.IstanbulCompatibleBlock.Uint64())
}

// TestCommittersForkGating: the engine sealer recovers committed seals with the round-bound
// preimage post-permissionless and the legacy preimage before it.
func TestCommittersForkGating(t *testing.T) {
	defer types.SetHeaderHashFn(nil)

	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	addr := crypto.PubkeyToAddress(key.PublicKey)

	// header carrying a round-bound committed seal
	impl := istanbul.NewSealerImpl(key)
	header := &types.Header{Number: big.NewInt(1)}
	require.NoError(t, impl.WriteValidators(header, []common.Address{addr}))
	impl.WriteRound(header, 0)
	seal, err := impl.MakeCommittedSealFromHashWithRound(impl.HeaderHash(header), 0)
	require.NoError(t, err)
	require.NoError(t, impl.WriteCommittedSeals(header, [][]byte{seal}))

	// post-fork routes to round-bound recovery (finds the signer); pre-fork stays legacy (misses it)
	got, err := NewSealer(&params.ChainConfig{PermissionlessCompatibleBlock: big.NewInt(0)}, nil).Committers(header)
	require.NoError(t, err)
	assert.Equal(t, []common.Address{addr}, got)

	got, err = NewSealer(&params.ChainConfig{}, nil).Committers(header)
	require.NoError(t, err)
	assert.NotEqual(t, []common.Address{addr}, got)
}
