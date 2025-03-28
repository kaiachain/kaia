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
			TargetTxHash:  common.HexToHash("0xc7f1b27b0c69006738b17567a7127c4d163fac7b575d046c6cbc90e62e6355e8"),
			BlockNumber:   1,
			Sender:        common.HexToAddress("0x14791697260E4c9A71f18484C9f997B308e59325"),
			To:            common.HexToAddress("0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"),
			Nonce:         4,
			Bid:           big.NewInt(3),
			CallGasLimit:  2,
			Data:          common.Hex2Bytes("1234"),
			SearcherSig:   common.Hex2Bytes("2cd97ec3eb8230a8cac9169146ea6ca406d908edd488e5fda30811ebf56647d94740d582c592e3476481b3fbab38a100623d2f4b0615da8b8dfd0f99128879901b"),
			AuctioneerSig: common.Hex2Bytes("d87718806c267dd6de19f4ed1111742750ee8040fdb3d18b1bd0dc1020ad8ca84262dfb4a3449f53b2cef8e2142796a96cca9ff8d08302f07db1d53a7b792e8d1c"),
		}
		invalidTargetTx            = baseBid
		invalidSearcherSigLenBid   = baseBid
		unexpectedSEarcherSigBid   = baseBid
		invalidAuctioneerSigLenBid = baseBid
		unexpectedAuctioneerSigBid = baseBid
		validBid                   = baseBid
		anotherBid                 = baseBid
	)
	invalidTargetTx.TargetTxRaw = common.Hex2Bytes("1234")

	invalidSearcherSigLenBid.SearcherSig = common.Hex2Bytes("1234")
	unexpectedSEarcherSigBid.SearcherSig = common.Hex2Bytes("2cd97ec3eb8230a8cac9169146ea6ca406d908edd488e5fda30811ebf56647d94740d582c592e3476481b3fbab38a100623d2f4b0615da8b8dfd0f99128879901c")

	invalidAuctioneerSigLenBid.AuctioneerSig = common.Hex2Bytes("1234")
	unexpectedAuctioneerSigBid.AuctioneerSig = common.Hex2Bytes("d87718806c267dd6de19f4ed1111742750ee8040fdb3d18b1bd0dc1020ad8ca84262dfb4a3449f53b2cef8e2142796a96cca9ff8d08302f07dc1d53a7b792e8d1c")

	validBid.TargetTxRaw = nil

	anotherBid.Bid = big.NewInt(10)
	anotherBid.SearcherSig = common.Hex2Bytes("a80aef9b383c2d947a4fdbedd6c211e83946a475a8bb9ac47afb494fe1bb87bc492898e25cf867a62f63f4e667208a172064abc86d9f854eb01905dc5aad02ea1b")
	anotherBid.AuctioneerSig = common.Hex2Bytes("cb14cd6c0016da027ebf680988584fc9332a57234c5a3702baa488894b1e50c939cf19ba405e1f6c9692a41939e0594f664d9030e98a889519e5ade73b46a57c1c")

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
			makeRPCOutput(common.Hash{}, nil, nil, errors.New("invalid searcher sig: expected 0x14791697260E4c9A71f18484C9f997B308e59325, calculated 0xBAc7570F225089fE23C6cF96e4D37fB94BDAd222")),
		},
		{
			"invalid auctioneer signature length",
			invalidAuctioneerSigLenBid,
			makeRPCOutput(common.Hash{}, nil, nil, auction.ErrInvalidAuctioneerSig),
		},
		{
			"unexpected auctioneer signature length",
			unexpectedAuctioneerSigBid,
			makeRPCOutput(common.Hash{}, nil, nil, errors.New("invalid auctioneer sig: expected 0x96Bd8E216c0D894C0486341288Bf486d5686C5b6, calculated 0x913c15715cAdC50aAD43F11C88f7a6Ee4964925B")),
		},
		{
			"valid bid and target tx can be empty",
			validBid,
			makeRPCOutput(common.HexToHash("0x60dda343662263ebcf704871a94420e2c21968662f3026f28730f9dfbf1edae7"), nil, nil, nil),
		},
		{
			"another bid with same target tx",
			anotherBid,
			makeRPCOutput(common.HexToHash("0x42afaab759d44e9f5d6cacc714435e653aba87b674a3ebe8bf2edf6333cff5e7"), nil, nil, nil),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := api.SubmitBid(context.Background(), tc.bidInput)
			assert.Equal(t, err, tc.expected)
		})
	}
}
