package impl

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/kaiax/randao"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test low-level computation components
func TestCalcRandaoPrimitives(t *testing.T) {
	var (
		skhex = hexutil.MustDecode("0x6c605527c8e4f31c959478801d51384d690a22dfc6438604646f7709032c893a")
		sk, _ = bls.SecretKeyFromBytes(skhex)
		pk    = sk.PublicKey()

		// block_num_to_bytes() = num.to_bytes(32, byteorder="big")
		num = big.NewInt(31337)
		msg = common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000007a69")

		// mix2 = xor(mix1, keccak256(sig))
		sig  = hexutil.MustDecode("0xadfe25ced45819332cbf088f01cdd2807686dd6309b11d7440237dd623624f401d4753747f5fb92374235e997edcd18318bae2806a1675b1e685e792abd1fbdf5c50ec1e148cc7fe861984d8bc3204c1b2136725b176902bc52eeb595919df3b")
		mix1 = hexutil.MustDecode("0x8019df1a2a9f833dc7f400a15b33e54a5c80295165c5953dc23891aab9203810")
		mix2 = hexutil.MustDecode("0x8772d58248bdf34e81ecbf36f28299cfa758b61ccf3f64e1dc0646687a55892f")
	)

	// Calculate RandomReveal and MixHash
	assert.Equal(t, msg, calcRandaoMsg(num))
	assert.Equal(t, sig, bls.Sign(sk, msg[:]).Marshal())
	assert.Equal(t, mix2, calcMixHash(sig, mix1))

	// Verify signature
	ok, err := bls.VerifySignature(sig, msg, pk)
	assert.Nil(t, err)
	assert.True(t, ok)
}

func setupConsensusTestModule(t *testing.T, forkNum *big.Int) (*RandaoModule, backends.BlockChainForCaller, *types.Header) {
	t.Helper()

	skhex := hexutil.MustDecode("0x6c605527c8e4f31c959478801d51384d690a22dfc6438604646f7709032c893a")
	blsKey, err := bls.SecretKeyFromBytes(skhex)
	require.NoError(t, err)

	config := testRandaoForkChainConfig(forkNum)
	config.KaiaCompatibleBlock = new(big.Int).Set(forkNum)
	db := database.NewMemoryDBManager()

	registryStorage := system.AllocRegistry(&params.RegistryConfig{
		Records: map[string]common.Address{
			system.Kip113Name: system.Kip113LogicAddrMock,
		},
		Owner: common.HexToAddress("0xffff"),
	})
	infos := map[common.Address]system.BlsPublicKeyInfo{
		params.AuthorAddressForTesting: {
			PublicKey: blsKey.PublicKey().Marshal(),
			Pop:       bls.PopProve(blsKey).Marshal(),
		},
	}
	kip113Storage := system.AllocKip113Proxy(system.AllocKip113Init{
		Infos: infos,
		Owner: common.HexToAddress("0xffff"),
	})
	alloc := blockchain.GenesisAlloc{
		system.RegistryAddr: {
			Code:    system.RegistryMockCode,
			Balance: big.NewInt(0),
			Storage: registryStorage,
		},
		system.Kip113LogicAddrMock: {
			Code:    system.Kip113MockCode,
			Balance: big.NewInt(0),
			Storage: kip113Storage,
		},
	}

	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)
	chain := backend.BlockChain()
	parent := chain.CurrentHeader()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     new(big.Int).Add(parent.Number, common.Big1),
	}

	m := NewRandaoModule()
	err = m.Init(&InitOpts{
		ChainConfig:  config,
		Chain:        chain,
		Downloader:   &downloader.FakeDownloader{},
		BlsSecretKey: blsKey,
	})
	require.NoError(t, err)

	return m, chain, header
}

func TestPrepareAndVerifyHeader(t *testing.T) {
	m, _, header := setupConsensusTestModule(t, big.NewInt(1))

	require.NoError(t, m.PrepareHeader(header))
	assert.Len(t, header.RandomReveal, 96)
	assert.Len(t, header.MixHash, 32)

	require.NoError(t, m.VerifyHeader(header))
}

func TestVerifyHeaderRejectsUnexpectedRandao(t *testing.T) {
	m, _, header := setupConsensusTestModule(t, big.NewInt(2))
	header.RandomReveal = make([]byte, 96)
	header.MixHash = make([]byte, 32)

	err := m.VerifyHeader(header)
	require.Error(t, err)
	assert.ErrorIs(t, err, randao.ErrUnexpectedRandao)
}

func TestVerifyHeaderRejectsInvalidRandaoFields(t *testing.T) {
	m, _, header := setupConsensusTestModule(t, big.NewInt(1))
	require.NoError(t, m.PrepareHeader(header))

	header.MixHash[0] ^= 0xff

	err := m.VerifyHeader(header)
	require.Error(t, err)
	assert.ErrorIs(t, err, randao.ErrInvalidRandaoFields)
}
