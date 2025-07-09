package impl

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

func prep() *AuctionModule {
	var (
		db     = database.NewMemoryDBManager()
		alloc  = testAllocStorage()
		config = testRandaoForkChainConfig(big.NewInt(0))
	)

	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	mAuction := NewAuctionModule()
	apiBackend := &MockBackend{}
	fakeDownloader := &downloader.FakeDownloader{}
	mAuction.Init(&InitOpts{
		ChainConfig: config,
		Chain:       backend.BlockChain(),
		Backend:     apiBackend,
		Downloader:  fakeDownloader,
	})
	mAuction.bidPool.running = 1
	mAuction.bidPool.auctioneer = common.HexToAddress("0x96Bd8E216c0D894C0486341288Bf486d5686C5b6")
	mAuction.bidPool.ChainConfig.ChainID = big.NewInt(1000)
	mAuction.bidPool.auctionEntryPoint = common.HexToAddress("0x6869431f189dCd7C2B92002aA61FCD4c1c0C1A33")
	return mAuction
}

func TestSubmitBid(t *testing.T) {
	var (
		mAuction = prep()
		api      = newAuctionAPI(mAuction)
		baseBid  = BidInput{
			TargetTxRaw:   common.Hex2Bytes("f8674785066720b30083015f909496bd8e216c0d894c0486341288bf486d5686c5b601808207f4a0a97fa83b989a6d66acc942d1cbd70f548c21e24eefea12e72f8c27ba4369a434a01900811315ba3c64055e9778470f438128b54a46712cc032f25a1487e2144578"),
			TargetTxHash:  common.HexToHash("0xacb81e7c775471be3e286a461701436f74b7bf7b951096f979b8557d870f246e"),
			BlockNumber:   1,
			Sender:        common.HexToAddress("0x14791697260E4c9A71f18484C9f997B308e59325"),
			To:            common.HexToAddress("0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"),
			Nonce:         4,
			Bid:           big.NewInt(3),
			CallGasLimit:  2,
			Data:          common.Hex2Bytes("1234"),
			SearcherSig:   common.Hex2Bytes("4f17bd3304ab18e8fd19938b6b3898c491134ecdd6a104244b458dc339ce2bf043f3b8d0a6a96d34cb27180146fe265e3213bb9ddcbafc0778cc39cde4d388d31b"),
			AuctioneerSig: common.Hex2Bytes("a9cfe35e9352818d7062b9a2ecfff939f46781ca352f35f56e790d4eaeb261e03564b4113517bda854eb530d642fdbc082085ead664e31014e902dbf4061fb841c"),
		}
		invalidTargetTx            = baseBid
		invalidSearcherSigLenBid   = baseBid
		unexpectedSEarcherSigBid   = baseBid
		invalidAuctioneerSigLenBid = baseBid
		unexpectedAuctioneerSigBid = baseBid
		diffTargetTx               = baseBid
		validBid                   = baseBid
		anotherBid                 = baseBid
	)
	invalidTargetTx.TargetTxRaw = common.Hex2Bytes("1234")

	invalidSearcherSigLenBid.SearcherSig = common.Hex2Bytes("1234")
	unexpectedSEarcherSigBid.SearcherSig = common.Hex2Bytes("2cd97ec3eb8230a8cac9169146ea6ca406d908edd488e5fda30811ebf56647d94740d582c592e3476481b3fbab38a100623d2f4b0615da8b8dfd0f99128879901c")

	invalidAuctioneerSigLenBid.AuctioneerSig = common.Hex2Bytes("1234")
	unexpectedAuctioneerSigBid.AuctioneerSig = common.Hex2Bytes("d87718806c267dd6de19f4ed1111742750ee8040fdb3d18b1bd0dc1020ad8ca84262dfb4a3449f53b2cef8e2142796a96cca9ff8d08302f07dc1d53a7b792e8d1c")

	diffTargetTx.TargetTxRaw = common.Hex2Bytes("f8674785066720b30083015f909496bd8e216c0d894c0486341288bf486d5686c5b601808207f4a0a97fa83b989a6d66acc942d1cbd70f548c21e24eefea12e72f8c27ba4369a434a01900811315ba3c64055e9778470f438128b54a46712cc032f25a1487e2144579")

	validBid.TargetTxRaw = nil

	anotherBid.Bid = big.NewInt(10)
	anotherBid.SearcherSig = common.Hex2Bytes("6439652673f1544bcd95d25c1dad31944321bdc0e6720f6c59a582aa0c40cc403ef4b5d1865eb3fa0e26fc49d7ef88f77f42d1559131a83a2326445eab3649741b")
	anotherBid.AuctioneerSig = common.Hex2Bytes("640a09994942d99bb751db3347ea3e909752b363a90dc3ed9c0b4d8ad512ae3d44ec820239f9875dcdbffc28fafafd3ca7d48a4ff6f4a2d8969a8f5d309460361b")

	tcs := []struct {
		name     string
		bidInput BidInput
		expected RPCOutput
	}{
		{
			"target tx decoding error",
			invalidTargetTx,
			makeRPCOutput(common.Hash{}, errors.New("undefined tx type"), nil, nil),
		},
		{
			"invalid seacher signature length",
			invalidSearcherSigLenBid,
			makeRPCOutput(common.Hash{}, nil, nil, auction.ErrInvalidSearcherSig),
		},
		{
			"unexpected seacher signature",
			unexpectedSEarcherSigBid,
			makeRPCOutput(common.Hash{}, nil, nil, errors.New("invalid searcher sig: expected 0x14791697260E4c9A71f18484C9f997B308e59325, calculated 0x5CD48323C0ebc334437ae933E782F2761F8196cA")),
		},
		{
			"invalid auctioneer signature length",
			invalidAuctioneerSigLenBid,
			makeRPCOutput(common.Hash{}, nil, nil, auction.ErrInvalidAuctioneerSig),
		},
		{
			"unexpected auctioneer signature length",
			unexpectedAuctioneerSigBid,
			makeRPCOutput(common.Hash{}, nil, nil, errors.New("invalid auctioneer sig: expected 0x96Bd8E216c0D894C0486341288Bf486d5686C5b6, calculated 0xd9094A8A697677ab51AA715F6449253Eb9c9885A")),
		},
		{
			"if target tx is not empty, its hash must be the same with bid's target tx hash",
			diffTargetTx,
			makeRPCOutput(common.Hash{}, auction.ErrInvalidTargetTxHash, nil, nil),
		},
		{
			"valid bid and target tx can be empty",
			validBid,
			makeRPCOutput(common.HexToHash("0xec633e59d7237fc6cce22dc3ca2dacf5ce9230276644d8e134b5d307ca7981bd"), nil, nil, nil),
		},
		{
			"another bid with same target tx",
			anotherBid,
			makeRPCOutput(common.HexToHash("0x26688d0fc660b6fed98b7f96ab5602e4c4dbe133e278fc08cc6bc51131d1bdd2"), nil, nil, nil),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := api.SubmitBid(context.Background(), tc.bidInput)
			assert.Equal(t, err, tc.expected)
		})
	}
}
