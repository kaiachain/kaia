package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts"
	mock_accounts "github.com/kaiachain/kaia/accounts/mocks"
	mock_api "github.com/kaiachain/kaia/api/mocks"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/gxhash"
	"github.com/kaiachain/kaia/consensus/mocks"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dummyChainConfigForEthAPITest = &params.ChainConfig{
	ChainID:                  new(big.Int).SetUint64(111111),
	IstanbulCompatibleBlock:  new(big.Int).SetUint64(0),
	LondonCompatibleBlock:    new(big.Int).SetUint64(0),
	EthTxTypeCompatibleBlock: new(big.Int).SetUint64(0),
	UnitPrice:                25000000000, // 25 gkei
}

var (
	testLondonConfig = &params.ChainConfig{
		ChainID:                 new(big.Int).SetUint64(111111),
		IstanbulCompatibleBlock: common.Big0,
		LondonCompatibleBlock:   common.Big0,
		UnitPrice:               25000000000,
	}
	testEthTxTypeConfig = &params.ChainConfig{
		ChainID:                  new(big.Int).SetUint64(111111),
		IstanbulCompatibleBlock:  common.Big0,
		LondonCompatibleBlock:    common.Big0,
		EthTxTypeCompatibleBlock: common.Big0,
		UnitPrice:                25000000000, // 25 gkei
	}
	testRandaoConfig = &params.ChainConfig{
		ChainID:                  new(big.Int).SetUint64(111111),
		IstanbulCompatibleBlock:  common.Big0,
		LondonCompatibleBlock:    common.Big0,
		EthTxTypeCompatibleBlock: common.Big0,
		MagmaCompatibleBlock:     common.Big0,
		KoreCompatibleBlock:      common.Big0,
		ShanghaiCompatibleBlock:  common.Big0,
		CancunCompatibleBlock:    common.Big0,
		RandaoCompatibleBlock:    common.Big0,
		UnitPrice:                25000000000, // 25 gkei
	}
)

// For floor data gas error case (EIP-7623)
var (
	// // SPDX-License-Identifier: GPL-3.0

	// pragma solidity >=0.8.2 <0.9.0;
	//
	// /**
	//  * @title Storage
	//  * @dev Store & retrieve value in a variable
	//  * @custom:dev-run-script ./scripts/deploy_with_ethers.ts
	//  */
	// contract Storage {
	//
	//     uint256 number;
	//
	//     /**
	//      * @dev Return value
	//      * @return value of 'number'
	//      */
	//     function retrieve() public view returns (uint256){
	//         return number;
	//     }
	// }
	floorDataGasTestCode = hexutil.Bytes(common.Hex2Bytes("6080604052348015600e575f5ffd5b50600436106026575f3560e01c80632e64cec114602a575b5f5ffd5b60306044565b604051603b91906062565b60405180910390f35b5f5f54905090565b5f819050919050565b605c81604c565b82525050565b5f60208201905060735f8301846055565b9291505056fea26469706673582212206aeab8d313a899d42a212113167e622ff770e746a3c3d0596d15fe2551d2c97464736f6c634300081e0033"))
	// retrieve() with long junk to increase floor data gas. 104 nonzero tokens = 4160 floor data gas = 25160 intrinsic gas
	floorDataGasTestData = hexutil.Bytes(append(common.Hex2Bytes("2e64cec1"), bytes.Repeat([]byte{0xff}, 100)...))
)

// TestEthereumAPI_Etherbase tests Etherbase.
func TestEthAPI_Etherbase(t *testing.T) {
	testNodeAddress(t, "Etherbase")
}

// TestEthAPI_Coinbase tests Coinbase.
func TestEthAPI_Coinbase(t *testing.T) {
	testNodeAddress(t, "Coinbase")
}

// testNodeAddress generates nodeAddress and tests Etherbase and Coinbase.
func testNodeAddress(t *testing.T, testAPIName string) {
	key, _ := crypto.GenerateKey()
	nodeAddress := crypto.PubkeyToAddress(key.PublicKey)

	api := EthAPI{nodeAddress: nodeAddress}
	results := reflect.ValueOf(&api).MethodByName(testAPIName).Call([]reflect.Value{})
	result, ok := results[0].Interface().(common.Address)
	assert.True(t, ok)
	assert.Equal(t, nodeAddress, result)
}

// TestEthAPI_Hashrate tests Hasharate.
func TestEthAPI_Hashrate(t *testing.T) {
	api := &EthAPI{}
	assert.Equal(t, hexutil.Uint64(ZeroHashrate), api.Hashrate())
}

// TestEthAPI_Mining tests Mining.
func TestEthAPI_Mining(t *testing.T) {
	api := &EthAPI{}
	assert.Equal(t, false, api.Mining())
}

// TestEthAPI_GetWork tests GetWork.
func TestEthAPI_GetWork(t *testing.T) {
	api := &EthAPI{}
	_, err := api.GetWork()
	assert.Equal(t, errNoMiningWork, err)
}

// TestEthAPI_SubmitWork tests SubmitWork.
func TestEthAPI_SubmitWork(t *testing.T) {
	api := &EthAPI{}
	assert.Equal(t, false, api.SubmitWork(BlockNonce{}, common.Hash{}, common.Hash{}))
}

// TestEthAPI_SubmitHashrate tests SubmitHashrate.
func TestEthAPI_SubmitHashrate(t *testing.T) {
	api := &EthAPI{}
	assert.Equal(t, false, api.SubmitHashrate(hexutil.Uint64(0), common.Hash{}))
}

// TestEthAPI_GetHashrate tests GetHashrate.
func TestEthAPI_GetHashrate(t *testing.T) {
	api := &EthAPI{}
	assert.Equal(t, ZeroHashrate, api.GetHashrate())
}

// TestEthAPI_GetUncleByBlockNumberAndIndex tests GetUncleByBlockNumberAndIndex.
func TestEthAPI_GetUncleByBlockNumberAndIndex(t *testing.T) {
	api := &EthAPI{}
	uncleBlock, err := api.GetUncleByBlockNumberAndIndex(context.Background(), rpc.BlockNumber(0), hexutil.Uint(0))
	assert.NoError(t, err)
	assert.Nil(t, uncleBlock)
}

// TestEthAPI_GetUncleByBlockHashAndIndex tests GetUncleByBlockHashAndIndex.
func TestEthAPI_GetUncleByBlockHashAndIndex(t *testing.T) {
	api := &EthAPI{}
	uncleBlock, err := api.GetUncleByBlockHashAndIndex(context.Background(), common.Hash{}, hexutil.Uint(0))
	assert.NoError(t, err)
	assert.Nil(t, uncleBlock)
}

// TestTestEthAPI_GetUncleCountByBlockNumber tests GetUncleCountByBlockNumber.
func TestTestEthAPI_GetUncleCountByBlockNumber(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, _, _, _, _ := createTestData(t, nil)

	// For existing block number, it must return 0.
	mockBackend.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil)
	existingBlockNumber := rpc.BlockNumber(block.Number().Int64())
	assert.Equal(t, hexutil.Uint(ZeroUncleCount), *api.GetUncleCountByBlockNumber(context.Background(), existingBlockNumber))

	// For non-existing block number, it must return nil.
	mockBackend.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(nil, nil)
	nonExistingBlockNumber := rpc.BlockNumber(5)
	uncleCount := api.GetUncleCountByBlockNumber(context.Background(), nonExistingBlockNumber)
	uintNil := hexutil.Uint(uint(0))
	expectedResult := &uintNil
	expectedResult = nil
	assert.Equal(t, expectedResult, uncleCount)

	mockCtrl.Finish()
}

// TestTestEthAPI_GetUncleCountByBlockHash tests GetUncleCountByBlockHash.
func TestTestEthAPI_GetUncleCountByBlockHash(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, _, _, _, _ := createTestData(t, nil)

	// For existing block hash, it must return 0.
	mockBackend.EXPECT().BlockByHash(gomock.Any(), gomock.Any()).Return(block, nil)
	existingHash := block.Hash()
	assert.Equal(t, hexutil.Uint(ZeroUncleCount), *api.GetUncleCountByBlockHash(context.Background(), existingHash))

	// For non-existing block hash, it must return nil.
	mockBackend.EXPECT().BlockByHash(gomock.Any(), gomock.Any()).Return(nil, nil)
	nonExistingHash := block.Hash()
	uncleCount := api.GetUncleCountByBlockHash(context.Background(), nonExistingHash)
	uintNil := hexutil.Uint(uint(0))
	expectedResult := &uintNil
	expectedResult = nil
	assert.Equal(t, expectedResult, uncleCount)

	mockCtrl.Finish()
}

// TestEthAPI_GetHeaderByNumber tests GetHeaderByNumber.
func TestEthAPI_GetHeaderByNumber(t *testing.T) {
	testGetHeader(t, "GetHeaderByNumber", testLondonConfig)
	testGetHeader(t, "GetHeaderByNumber", testEthTxTypeConfig)
	testGetHeader(t, "GetHeaderByNumber", testRandaoConfig)
}

// TestEthAPI_GetHeaderByHash tests GetHeaderByNumber.
func TestEthAPI_GetHeaderByHash(t *testing.T) {
	testGetHeader(t, "GetHeaderByHash", testLondonConfig)
	testGetHeader(t, "GetHeaderByHash", testEthTxTypeConfig)
	testGetHeader(t, "GetHeaderByHash", testRandaoConfig)
}

// testGetHeader generates data to test GetHeader related functions in EthAPI
// and actually tests the API function passed as a parameter.
func testGetHeader(t *testing.T, testAPIName string, config *params.ChainConfig) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)

	// Creates a MockEngine.
	mockEngine := mocks.NewMockEngine(mockCtrl)
	// GetHeader APIs calls internally below methods.
	mockBackend.EXPECT().Engine().Return(mockEngine)
	mockBackend.EXPECT().ChainConfig().Return(config).AnyTimes()

	// Author is called when calculates miner field of Header.
	dummyMiner := common.HexToAddress("0x9712f943b296758aaae79944ec975884188d3a96")
	mockEngine.EXPECT().Author(gomock.Any()).Return(dummyMiner, nil)

	// Create dummy header
	header := types.CopyHeader(&types.Header{
		ParentHash:  common.HexToHash("0xc8036293065bacdfce87debec0094a71dbbe40345b078d21dcc47adb4513f348"),
		Rewardbase:  common.Address{},
		TxHash:      types.EmptyTxRootOriginal,
		Root:        common.HexToHash("0xad31c32942fa033166e4ef588ab973dbe26657c594de4ba98192108becf0fec9"),
		ReceiptHash: types.EmptyTxRootOriginal,
		Bloom:       types.Bloom{},
		BlockScore:  new(big.Int).SetUint64(1),
		Number:      new(big.Int).SetUint64(4),
		GasUsed:     uint64(10000),
		Time:        new(big.Int).SetUint64(1641363540),
		TimeFoS:     uint8(85),
		Extra:       common.Hex2Bytes("0xd983010701846b6c617988676f312e31362e338664617277696e000000000000f89ed5949712f943b296758aaae79944ec975884188d3a96b8415a0614be7fd5ea40f11ce558e02993bd55f11ae72a3cfbc861875a57483ec5ec3adda3e5845fd7ab271d670c755480f9ef5b8dd731f4e1f032fff5d165b763ac01f843b8418867d3733167a0c737fa5b62dcc59ec3b0af5748bcc894e7990a0b5a642da4546713c9127b3358cdfe7894df1ca1db5a97560599986d7f1399003cd63660b98200"),
		Governance:  []byte{},
		Vote:        []byte{},
	})
	if config.IsRandaoForkEnabled(common.Big0) {
		header.RandomReveal = hexutil.MustDecode("0x94516a8bc695b5bf43aa077cd682d9475a3a6bed39a633395b78ed8f276e7c5bb00bb26a77825013c6718579f1b3ee2275b158801705ea77989e3acc849ee9c524bd1822bde3cba7be2aae04347f0d91508b7b7ce2f11ec36cbf763173421ae7")
		header.MixHash = hexutil.MustDecode("0xdf117d1245dceaae0a47f05371b23cd0d0db963ff9d5c8ba768dc989f4c31883")
	}

	var blockParam interface{}
	switch testAPIName {
	case "GetHeaderByNumber":
		blockParam = rpc.BlockNumber(header.Number.Uint64())
		mockBackend.EXPECT().HeaderByNumber(gomock.Any(), gomock.Any()).Return(header, nil)
	case "GetHeaderByHash":
		blockParam = header.Hash()
		mockBackend.EXPECT().HeaderByHash(gomock.Any(), gomock.Any()).Return(header, nil)
	}

	results := reflect.ValueOf(&api).MethodByName(testAPIName).Call(
		[]reflect.Value{
			reflect.ValueOf(context.Background()),
			reflect.ValueOf(blockParam),
		},
	)
	ethHeader, ok := results[0].Interface().(map[string]interface{})
	assert.Equal(t, true, ok)
	assert.NotEqual(t, ethHeader, nil)

	// We can get a real mashaled data by using real backend instance, not mock
	// Mock just return a header instance, not rlp decoded json data
	expected := make(map[string]interface{})
	assert.NoError(t, json.Unmarshal([]byte(`
	{
		"difficulty": "0x1",
		"extraData": "0x",
		"gasLimit": "0xe8d4a50fff",
		"gasUsed": "0x2710",
		"hash": "0x852754129164bc6f3cdf4beaab557f2b058634be42e3470d5ef56a8e4ff01685",
		"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"miner": "0x9712f943b296758aaae79944ec975884188d3a96",
		"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"nonce": "0x0000000000000000",
		"number": "0x4",
		"parentHash": "0xc8036293065bacdfce87debec0094a71dbbe40345b078d21dcc47adb4513f348",
		"receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size": "0x244",
		"stateRoot": "0xad31c32942fa033166e4ef588ab973dbe26657c594de4ba98192108becf0fec9",
		"timestamp": "0x61d53854",
		"transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421"
	}`), &expected))

	if config.IsEthTxTypeForkEnabled(common.Big0) {
		expected["baseFeePerGas"] = "0x0"
	}
	if config.IsRandaoForkEnabled(common.Big0) {
		expected["randomReveal"] = "0x94516a8bc695b5bf43aa077cd682d9475a3a6bed39a633395b78ed8f276e7c5bb00bb26a77825013c6718579f1b3ee2275b158801705ea77989e3acc849ee9c524bd1822bde3cba7be2aae04347f0d91508b7b7ce2f11ec36cbf763173421ae7"
		expected["mixHash"] = "0xdf117d1245dceaae0a47f05371b23cd0d0db963ff9d5c8ba768dc989f4c31883"
		expected["hash"] = "0x36f1c36d1723049abf1202a1cda828eec6399edd654dae12b72a1642097a29e4"
		expected["size"] = "0x2c4"
	}
	assert.Equal(t, stringifyMap(expected), stringifyMap(ethHeader))
}

// TestEthAPI_GetBlockByNumber tests GetBlockByNumber.
func TestEthAPI_GetBlockByNumber(t *testing.T) {
	testGetBlock(t, "GetBlockByNumber", false)
	testGetBlock(t, "GetBlockByNumber", true)
}

// TestEthAPI_GetBlockByHash tests GetBlockByHash.
func TestEthAPI_GetBlockByHash(t *testing.T) {
	testGetBlock(t, "GetBlockByHash", false)
	testGetBlock(t, "GetBlockByHash", true)
}

// testGetBlock generates data to test GetBlock related functions in EthAPI
// and actually tests the API function passed as a parameter.
func testGetBlock(t *testing.T, testAPIName string, fullTxs bool) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)

	// Creates a MockEngine.
	mockEngine := mocks.NewMockEngine(mockCtrl)
	// GetHeader APIs calls internally below methods.
	mockBackend.EXPECT().Engine().Return(mockEngine)
	mockBackend.EXPECT().ChainConfig().Return(dummyChainConfigForEthAPITest).AnyTimes()
	// Author is called when calculates miner field of Header.
	dummyMiner := common.HexToAddress("0x9712f943b296758aaae79944ec975884188d3a96")
	mockEngine.EXPECT().Author(gomock.Any()).Return(dummyMiner, nil)

	// Create dummy header
	header := types.CopyHeader(&types.Header{
		ParentHash: common.HexToHash("0xc8036293065bacdfce87debec0094a71dbbe40345b078d21dcc47adb4513f348"), Rewardbase: common.Address{}, TxHash: types.EmptyTxRootOriginal,
		Root:        common.HexToHash("0xad31c32942fa033166e4ef588ab973dbe26657c594de4ba98192108becf0fec9"),
		ReceiptHash: types.EmptyTxRootOriginal,
		Bloom:       types.Bloom{},
		BlockScore:  new(big.Int).SetUint64(1),
		Number:      new(big.Int).SetUint64(4),
		GasUsed:     uint64(10000),
		Time:        new(big.Int).SetUint64(1641363540),
		TimeFoS:     uint8(85),
		Extra:       common.Hex2Bytes("0xd983010701846b6c617988676f312e31362e338664617277696e000000000000f89ed5949712f943b296758aaae79944ec975884188d3a96b8415a0614be7fd5ea40f11ce558e02993bd55f11ae72a3cfbc861875a57483ec5ec3adda3e5845fd7ab271d670c755480f9ef5b8dd731f4e1f032fff5d165b763ac01f843b8418867d3733167a0c737fa5b62dcc59ec3b0af5748bcc894e7990a0b5a642da4546713c9127b3358cdfe7894df1ca1db5a97560599986d7f1399003cd63660b98200"),
		Governance:  []byte{},
		Vote:        []byte{},
	})
	block, _, _, _, _ := createTestData(t, header)
	var blockParam interface{}
	switch testAPIName {
	case "GetBlockByNumber":
		blockParam = rpc.BlockNumber(block.NumberU64())
		mockBackend.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil)
	case "GetBlockByHash":
		blockParam = block.Hash()
		mockBackend.EXPECT().BlockByHash(gomock.Any(), gomock.Any()).Return(block, nil)
	}

	results := reflect.ValueOf(&api).MethodByName(testAPIName).Call(
		[]reflect.Value{
			reflect.ValueOf(context.Background()),
			reflect.ValueOf(blockParam),
			reflect.ValueOf(fullTxs),
		},
	)
	ethBlock, ok := results[0].Interface().(map[string]interface{})
	assert.Equal(t, true, ok)
	assert.NotEqual(t, ethBlock, nil)

	expected := make(map[string]interface{})
	if fullTxs {
		assert.NoError(t, json.Unmarshal([]byte(`
    {
        "baseFeePerGas": "0x0",
        "difficulty": "0x1",
        "extraData": "0x",
        "gasLimit": "0xe8d4a50fff",
        "gasUsed": "0x2710",
        "hash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
        "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "miner": "0x9712f943b296758aaae79944ec975884188d3a96",
        "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "nonce": "0x0000000000000000",
        "number": "0x4",
        "parentHash": "0xc8036293065bacdfce87debec0094a71dbbe40345b078d21dcc47adb4513f348",
        "receiptsRoot": "0xf6278dd71ffc1637f78dc2ee54f6f9e64d4b1633c1179dfdbc8c3b482efbdbec",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0xe44",
        "stateRoot": "0xad31c32942fa033166e4ef588ab973dbe26657c594de4ba98192108becf0fec9",
        "timestamp": "0x61d53854",
        "transactions": [
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x0000000000000000000000000000000000000000",
              "gas": "0x1c9c380",
              "gasPrice": "0x5d21dba00",
              "hash": "0x6231f24f79d28bb5b8425ce577b3b77cd9c1ab766fcfc5233358a2b1c2f4ff70",
              "input": "0x3078653331393765386630303030303030303030303030303030303030303030303065306265663939623461323232383665323736333062343835643036633561313437636565393331303030303030303030303030303030303030303030303030313538626566663863386364656264363436353461646435663661316439393337653733353336633030303030303030303030303030303030303030303030303030303030303030303030303030303030303030323962623565376662366265616533326366383030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030306530303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303138303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030316236306662343631346132326530303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303031303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030343030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303031353862656666386338636465626436343635346164643566366131643939333765373335333663303030303030303030303030303030303030303030303030373462613033313938666564326231356135316166323432623963363366616633633866346433343030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303033303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030",
              "nonce": "0x0",
              "to": "0x3736346135356338333362313038373730343930",
              "transactionIndex": "0x0",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3036656164333031646165616636376537376538",
              "gas": "0x989680",
              "gasPrice": "0x5d21dba00",
              "hash": "0xf146858415c060eae65a389cbeea8aeadc79461038fbee331ffd97b41279dd63",
              "input": "0x",
              "nonce": "0x1",
              "to": "0x3364613566326466626334613262333837316462",
              "transactionIndex": "0x1",
              "value": "0x5",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3730323366383135666136613633663761613063",
              "gas": "0x1312d00",
              "gasPrice": "0x5d21dba00",
              "hash": "0x0a01fc67bb4c15c32fa43563c0fcf05cd5bf2fdcd4ec78122b5d0295993bca24",
              "input": "0x68656c6c6f",
              "nonce": "0x2",
              "to": "0x3336623562313539333066323466653862616538",
              "transactionIndex": "0x2",
              "value": "0x3",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x1312d00",
              "gasPrice": "0x5d21dba00",
              "hash": "0x486f7561375c38f1627264f8676f92ec0dd1c4a7c52002ba8714e61fcc6bb649",
              "input": "0x",
              "nonce": "0x3",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0x3",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x5f5e100",
              "gasPrice": "0x5d21dba00",
              "hash": "0xbd3e57cd31dd3d6679326f7a949f0de312e9ae53bec5ef3c23b43a5319c220a4",
              "input": "0x",
              "nonce": "0x4",
              "to": null,
              "transactionIndex": "0x4",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0xff666129a0c7227b17681d668ecdef5d6681fc93dbd58856eea1374880c598b0",
              "input": "0x",
              "nonce": "0x5",
              "to": "0x3632323232656162393565396564323963346266",
              "transactionIndex": "0x5",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0xa8ad4f295f2acff9ef56b476b1c52ecb74fb3fd95a789d768c2edb3376dbeacf",
              "input": "0x",
              "nonce": "0x6",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0x6",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0x47dbfd201fc1dd4188fd2003c6328a09bf49414be607867ca3a5d63573aede93",
              "input": "0xf8ad80b8aaf8a8a0072409b14b96f9d7dbf4788dbc68c5d30bd5fac1431c299e0ab55c92e70a28a4a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a00000000000000000000000000000000000000000000000000000000000000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a00000000000000000000000000000000000000000000000000000000000000000808080",
              "nonce": "0x7",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0x7",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3036656164333031646165616636376537376538",
              "gas": "0x989680",
              "gasPrice": "0x5d21dba00",
              "hash": "0x2283294e89b41df2df4dd37c375a3f51c3ad11877aa0a4b59d0f68cf5cfd865a",
              "input": "0x",
              "nonce": "0x8",
              "to": "0x3364613566326466626334613262333837316462",
              "transactionIndex": "0x8",
              "value": "0x5",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3730323366383135666136613633663761613063",
              "gas": "0x1312d00",
              "gasPrice": "0x5d21dba00",
              "hash": "0x80e05750d02d22d73926179a0611b431ae7658846406f836e903d76191423716",
              "input": "0x68656c6c6f",
              "nonce": "0x9",
              "to": "0x3336623562313539333066323466653862616538",
              "transactionIndex": "0x9",
              "value": "0x3",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x1312d00",
              "gasPrice": "0x5d21dba00",
              "hash": "0xe8abdee5e8fef72fe4d98f7dbef36000407e97874e8c880df4d85646958dd2c1",
              "input": "0x",
              "nonce": "0xa",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0xa",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x5f5e100",
              "gasPrice": "0x5d21dba00",
              "hash": "0x4c970be1815e58e6f69321202ce38b2e5c5e5ecb70205634848afdbc57224811",
              "input": "0x",
              "nonce": "0xb",
              "to": null,
              "transactionIndex": "0xb",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0x7ff0a809387d0a4cab77624d467f4d65ffc1ac95f4cc46c2246daab0407a7d83",
              "input": "0x",
              "nonce": "0xc",
              "to": "0x3632323232656162393565396564323963346266",
              "transactionIndex": "0xc",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0xb510b11415b39d18a972a00e3b43adae1e0f583ea0481a4296e169561ff4d916",
              "input": "0x",
              "nonce": "0xd",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0xd",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0xab466145fb71a2d24d6f6af3bddf3bcfa43c20a5937905dd01963eaf9fc5e382",
              "input": "0xf8ad80b8aaf8a8a0072409b14b96f9d7dbf4788dbc68c5d30bd5fac1431c299e0ab55c92e70a28a4a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a00000000000000000000000000000000000000000000000000000000000000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a00000000000000000000000000000000000000000000000000000000000000000808080",
              "nonce": "0xe",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0xe",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3036656164333031646165616636376537376538",
              "gas": "0x989680",
              "gasPrice": "0x5d21dba00",
              "hash": "0xec714ab0875768f482daeabf7eb7be804e3c94bc1f1b687359da506c7f3a66b2",
              "input": "0x",
              "nonce": "0xf",
              "to": "0x3364613566326466626334613262333837316462",
              "transactionIndex": "0xf",
              "value": "0x5",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3730323366383135666136613633663761613063",
              "gas": "0x1312d00",
              "gasPrice": "0x5d21dba00",
              "hash": "0x069af125fe88784e46f90ace9960a09e5d23e6ace20350062be75964a7ece8e6",
              "input": "0x68656c6c6f",
              "nonce": "0x10",
              "to": "0x3336623562313539333066323466653862616538",
              "transactionIndex": "0x10",
              "value": "0x3",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x1312d00",
              "gasPrice": "0x5d21dba00",
              "hash": "0x4a6bb7b2cd68265eb6a693aa270daffa3cc297765267f92be293b12e64948c82",
              "input": "0x",
              "nonce": "0x11",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0x11",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x5f5e100",
              "gasPrice": "0x5d21dba00",
              "hash": "0xa354fe3fdde6292e85545e6327c314827a20e0d7a1525398b38526fe28fd36e1",
              "input": "0x",
              "nonce": "0x12",
              "to": null,
              "transactionIndex": "0x12",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0x5bb64e885f196f7b515e62e3b90496864d960e2f5e0d7ad88550fa1c875ca691",
              "input": "0x",
              "nonce": "0x13",
              "to": "0x3632323232656162393565396564323963346266",
              "transactionIndex": "0x13",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0x6f4308b3c98db2db215d02c0df24472a215df7aa283261fcb06a6c9f796df9af",
              "input": "0x",
              "nonce": "0x14",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0x14",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            },
            {
              "blockHash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
              "blockNumber": "0x4",
              "from": "0x3936663364636533666637396132333733653330",
              "gas": "0x2faf080",
              "gasPrice": "0x5d21dba00",
              "hash": "0x1df88d113f0c5833c1f7264687cd6ac43888c232600ffba8d3a7d89bb5013e71",
              "input": "0xf8ad80b8aaf8a8a0072409b14b96f9d7dbf4788dbc68c5d30bd5fac1431c299e0ab55c92e70a28a4a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a00000000000000000000000000000000000000000000000000000000000000000a056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421a00000000000000000000000000000000000000000000000000000000000000000808080",
              "nonce": "0x15",
              "to": "0x3936663364636533666637396132333733653330",
              "transactionIndex": "0x15",
              "value": "0x0",
              "type": "0x0",
              "v": "0x1",
              "r": "0x2",
              "s": "0x3"
            }
        ],
        "transactionsRoot": "0x0a83e34ab7302f42f4a9203e8295f545517645989da6555d8cbdc1e9599df85b",
        "uncles": []
    }
    `,
		), &expected))
	} else {
		assert.NoError(t, json.Unmarshal([]byte(`
    {
        "baseFeePerGas": "0x0",
        "difficulty": "0x1",
        "extraData": "0x",
        "gasLimit": "0xe8d4a50fff",
        "gasUsed": "0x2710",
        "hash": "0xc74d8c04d4d2f2e4ed9cd1731387248367cea7f149731b7a015371b220ffa0fb",
        "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
        "miner": "0x9712f943b296758aaae79944ec975884188d3a96",
        "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "nonce": "0x0000000000000000",
        "number": "0x4",
        "parentHash": "0xc8036293065bacdfce87debec0094a71dbbe40345b078d21dcc47adb4513f348",
        "receiptsRoot": "0xf6278dd71ffc1637f78dc2ee54f6f9e64d4b1633c1179dfdbc8c3b482efbdbec",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0xe44",
        "stateRoot": "0xad31c32942fa033166e4ef588ab973dbe26657c594de4ba98192108becf0fec9",
        "timestamp": "0x61d53854",
        "transactions": [
            "0x6231f24f79d28bb5b8425ce577b3b77cd9c1ab766fcfc5233358a2b1c2f4ff70",
            "0xf146858415c060eae65a389cbeea8aeadc79461038fbee331ffd97b41279dd63",
            "0x0a01fc67bb4c15c32fa43563c0fcf05cd5bf2fdcd4ec78122b5d0295993bca24",
            "0x486f7561375c38f1627264f8676f92ec0dd1c4a7c52002ba8714e61fcc6bb649",
            "0xbd3e57cd31dd3d6679326f7a949f0de312e9ae53bec5ef3c23b43a5319c220a4",
            "0xff666129a0c7227b17681d668ecdef5d6681fc93dbd58856eea1374880c598b0",
            "0xa8ad4f295f2acff9ef56b476b1c52ecb74fb3fd95a789d768c2edb3376dbeacf",
            "0x47dbfd201fc1dd4188fd2003c6328a09bf49414be607867ca3a5d63573aede93",
            "0x2283294e89b41df2df4dd37c375a3f51c3ad11877aa0a4b59d0f68cf5cfd865a",
            "0x80e05750d02d22d73926179a0611b431ae7658846406f836e903d76191423716",
            "0xe8abdee5e8fef72fe4d98f7dbef36000407e97874e8c880df4d85646958dd2c1",
            "0x4c970be1815e58e6f69321202ce38b2e5c5e5ecb70205634848afdbc57224811",
            "0x7ff0a809387d0a4cab77624d467f4d65ffc1ac95f4cc46c2246daab0407a7d83",
            "0xb510b11415b39d18a972a00e3b43adae1e0f583ea0481a4296e169561ff4d916",
            "0xab466145fb71a2d24d6f6af3bddf3bcfa43c20a5937905dd01963eaf9fc5e382",
            "0xec714ab0875768f482daeabf7eb7be804e3c94bc1f1b687359da506c7f3a66b2",
            "0x069af125fe88784e46f90ace9960a09e5d23e6ace20350062be75964a7ece8e6",
            "0x4a6bb7b2cd68265eb6a693aa270daffa3cc297765267f92be293b12e64948c82",
            "0xa354fe3fdde6292e85545e6327c314827a20e0d7a1525398b38526fe28fd36e1",
            "0x5bb64e885f196f7b515e62e3b90496864d960e2f5e0d7ad88550fa1c875ca691",
            "0x6f4308b3c98db2db215d02c0df24472a215df7aa283261fcb06a6c9f796df9af",
            "0x1df88d113f0c5833c1f7264687cd6ac43888c232600ffba8d3a7d89bb5013e71"
        ],
        "transactionsRoot": "0x0a83e34ab7302f42f4a9203e8295f545517645989da6555d8cbdc1e9599df85b",
        "uncles": []
    }
    `,
		), &expected))
	}
	assert.Equal(t, stringifyMap(expected), stringifyMap(ethBlock))
}

// marshal and unmarshal to stringify map fields.
func stringifyMap(m map[string]interface{}) map[string]interface{} {
	marshaled, _ := json.Marshal(m)
	var unmarshaled map[string]interface{}
	json.Unmarshal(marshaled, &unmarshaled)
	return unmarshaled
}

// TestEthAPI_GetTransactionByBlockNumberAndIndex tests GetTransactionByBlockNumberAndIndex.
func TestEthAPI_GetTransactionByBlockNumberAndIndex(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, _, _, _ := createTestData(t, nil)

	// Mock Backend functions.
	mockBackend.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil).Times(txs.Len())
	mockBackend.EXPECT().ChainConfig().Return(dummyChainConfigForEthAPITest).AnyTimes()
	// Get transaction by block number and index for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		ethTx := api.GetTransactionByBlockNumberAndIndex(context.Background(), rpc.BlockNumber(block.NumberU64()), hexutil.Uint(i))
		checkEthRPCTransactionFormat(t, block, ethTx, txs[i], hexutil.Uint64(i))
	}

	mockCtrl.Finish()
}

// TestEthAPI_GetTransactionByBlockHashAndIndex tests GetTransactionByBlockHashAndIndex.
func TestEthAPI_GetTransactionByBlockHashAndIndex(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, _, _, _ := createTestData(t, nil)

	// Mock Backend functions.
	mockBackend.EXPECT().BlockByHash(gomock.Any(), gomock.Any()).Return(block, nil).Times(txs.Len())
	mockBackend.EXPECT().ChainConfig().Return(dummyChainConfigForEthAPITest).AnyTimes()

	// Get transaction by block hash and index for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		ethTx := api.GetTransactionByBlockHashAndIndex(context.Background(), block.Hash(), hexutil.Uint(i))
		checkEthRPCTransactionFormat(t, block, ethTx, txs[i], hexutil.Uint64(i))
	}

	mockCtrl.Finish()
}

// TestEthAPI_GetTransactionByHash tests GetTransactionByHash.
func TestEthAPI_GetTransactionByHash(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, txHashMap, _, _ := createTestData(t, nil)

	// Define queryFromPool for ReadTxAndLookupInfo function return tx from hash map.
	// MockDatabaseManager will initiate data with txHashMap, block and queryFromPool.
	// If queryFromPool is true, MockDatabaseManager will return nil to query transactions from transaction pool,
	// otherwise return a transaction from txHashMap.
	mockDBManager := &MockDatabaseManager{txHashMap: txHashMap, blockData: block, queryFromPool: false}

	// Mock Backend functions.
	mockBackend.EXPECT().ChainDB().Return(mockDBManager).Times(txs.Len())
	mockBackend.EXPECT().BlockByHash(gomock.Any(), block.Hash()).Return(block, nil).Times(txs.Len())
	mockBackend.EXPECT().ChainConfig().Return(dummyChainConfigForEthAPITest).AnyTimes()

	// Get transaction by hash for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		ethTx, err := api.GetTransactionByHash(context.Background(), txs[i].Hash())
		if err != nil {
			t.Fatal(err)
		}
		checkEthRPCTransactionFormat(t, block, ethTx, txs[i], hexutil.Uint64(i))
	}

	mockCtrl.Finish()
}

// TestEthAPI_GetTransactionByHash tests GetTransactionByHash from transaction pool.
func TestEthAPI_GetTransactionByHashFromPool(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, txHashMap, _, _ := createTestData(t, nil)

	// Define queryFromPool for ReadTxAndLookupInfo function return nil.
	// MockDatabaseManager will initiate data with txHashMap, block and queryFromPool.
	// If queryFromPool is true, MockDatabaseManager will return nil to query transactions from transaction pool,
	// otherwise return a transaction from txHashMap.
	mockDBManager := &MockDatabaseManager{txHashMap: txHashMap, blockData: block, queryFromPool: true}

	// Mock Backend functions.
	mockBackend.EXPECT().ChainDB().Return(mockDBManager).Times(txs.Len())
	mockBackend.EXPECT().ChainConfig().Return(dummyChainConfigForEthAPITest).AnyTimes()
	mockBackend.EXPECT().GetPoolTransaction(gomock.Any()).DoAndReturn(
		func(hash common.Hash) *types.Transaction {
			return txHashMap[hash]
		},
	).Times(txs.Len())

	//  Get transaction by hash from the transaction pool for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		ethTx, err := api.GetTransactionByHash(context.Background(), txs[i].Hash())
		if err != nil {
			t.Fatal(err)
		}
		checkEthRPCTransactionFormat(t, nil, ethTx, txs[i], 0)
	}

	mockCtrl.Finish()
}

// TestEthAPI_PendingTransactions tests PendingTransactions.
func TestEthAPI_PendingTransactions(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	_, txs, txHashMap, _, _ := createTestData(t, nil)

	mockAccountManager := mock_accounts.NewMockAccountManager(mockCtrl)
	mockBackend.EXPECT().AccountManager().Return(mockAccountManager)
	mockBackend.EXPECT().ChainConfig().Return(dummyChainConfigForEthAPITest).AnyTimes()
	mockBackend.EXPECT().GetPoolTransactions().Return(txs, nil)

	wallets := make([]accounts.Wallet, 1)
	wallets[0] = NewMockWallet(txs)
	mockAccountManager.EXPECT().Wallets().Return(wallets)

	pendingTxs, err := api.PendingTransactions()
	if err != nil {
		t.Fatal(err)
	}

	for _, pt := range pendingTxs {
		checkEthRPCTransactionFormat(t, nil, pt, txHashMap[pt.Hash], 0)
	}

	mockCtrl.Finish()
}

// TestEthAPI_GetTransactionReceipt tests GetTransactionReceipt.
func TestEthAPI_GetTransactionReceipt(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, txHashMap, receiptMap, receipts := createTestData(t, nil)

	// Mock Backend functions.
	mockBackend.EXPECT().GetTxLookupInfoAndReceipt(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, *types.Receipt) {
			txLookupInfo := txHashMap[hash]
			idx := txLookupInfo.Nonce() // Assume idx of the transaction is nonce
			return txLookupInfo, block.Hash(), block.NumberU64(), idx, receiptMap[hash]
		},
	).Times(txs.Len())
	mockBackend.EXPECT().GetBlockReceipts(gomock.Any(), gomock.Any()).Return(receipts).Times(txs.Len())
	mockBackend.EXPECT().HeaderByHash(gomock.Any(), block.Hash()).Return(block.Header(), nil).Times(txs.Len())
	mockBackend.EXPECT().ChainConfig().Return(testRandaoConfig).AnyTimes()

	// Get receipt for each transaction types.
	for i := 0; i < txs.Len(); i++ {
		receipt, err := api.GetTransactionReceipt(context.Background(), txs[i].Hash())
		if err != nil {
			t.Fatal(err)
		}
		txIdx := uint64(i)
		checkEthTransactionReceiptFormat(t, block, receipts, receipt, RpcOutputReceipt(block.Header(), txs[i], block.Hash(), block.NumberU64(), txIdx, receiptMap[txs[i].Hash()], params.TestChainConfig), txIdx)
	}

	mockCtrl.Finish()
}

func testInitForEthApi(t *testing.T) (*gomock.Controller, *mock_api.MockBackend, EthAPI) {
	mockCtrl := gomock.NewController(t)
	mockBackend := mock_api.NewMockBackend(mockCtrl)

	blockchain.InitDeriveSha(dummyChainConfigForEthAPITest)

	api := EthAPI{
		kaiaTransactionAPI: NewKaiaTransactionAPI(mockBackend, new(AddrLocker)),
		kaiaAPI:            NewKaiaAPI(mockBackend),
		kaiaBlockChainAPI:  NewKaiaBlockChainAPI(mockBackend),
	}
	return mockCtrl, mockBackend, api
}

func checkEthRPCTransactionFormat(t *testing.T, block *types.Block, ethTx *EthRPCTransaction, tx *types.Transaction, expectedIndex hexutil.Uint64) {
	// All Kaia transaction types must be returned as TxTypeLegacyTransaction types.
	assert.Equal(t, types.TxType(ethTx.Type), types.TxTypeLegacyTransaction)

	// Check the data of common fields of the transaction.
	from := getFrom(tx)
	assert.Equal(t, from, ethTx.From)
	assert.Equal(t, hexutil.Uint64(tx.Gas()), ethTx.Gas)
	assert.Equal(t, tx.GasPrice(), ethTx.GasPrice.ToInt())
	assert.Equal(t, tx.Hash(), ethTx.Hash)
	assert.Equal(t, tx.GetTxInternalData().RawSignatureValues()[0].V, ethTx.V.ToInt())
	assert.Equal(t, tx.GetTxInternalData().RawSignatureValues()[0].R, ethTx.R.ToInt())
	assert.Equal(t, tx.GetTxInternalData().RawSignatureValues()[0].S, ethTx.S.ToInt())
	assert.Equal(t, hexutil.Uint64(tx.Nonce()), ethTx.Nonce)

	// Check the optional field of Kaia transactions.
	assert.Equal(t, 0, bytes.Compare(ethTx.Input, tx.Data()))

	to := tx.To()
	switch tx.Type() {
	case types.TxTypeAccountUpdate, types.TxTypeFeeDelegatedAccountUpdate, types.TxTypeFeeDelegatedAccountUpdateWithRatio,
		types.TxTypeCancel, types.TxTypeFeeDelegatedCancel, types.TxTypeFeeDelegatedCancelWithRatio,
		types.TxTypeChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		assert.Equal(t, &from, ethTx.To)
	default:
		assert.Equal(t, to, ethTx.To)
	}
	value := tx.Value()
	assert.Equal(t, value, ethTx.Value.ToInt())

	// If it is not a pending transaction and has already been processed and added into a block,
	// the following fields should be returned.
	if block != nil {
		assert.Equal(t, block.Hash().String(), ethTx.BlockHash.String())
		assert.Equal(t, block.NumberU64(), ethTx.BlockNumber.ToInt().Uint64())
		assert.Equal(t, expectedIndex, *ethTx.TransactionIndex)
	}

	// Fields additionally used for Ethereum transaction types are not used
	// when returning Kaia transactions.
	assert.Equal(t, true, reflect.ValueOf(ethTx.Accesses).IsNil())
	assert.Equal(t, true, reflect.ValueOf(ethTx.ChainID).IsNil())
	assert.Equal(t, true, reflect.ValueOf(ethTx.GasFeeCap).IsNil())
	assert.Equal(t, true, reflect.ValueOf(ethTx.GasTipCap).IsNil())
}

func checkEthTransactionReceiptFormat(t *testing.T, block *types.Block, receipts []*types.Receipt, ethReceipt map[string]interface{}, kReceipt map[string]interface{}, idx uint64) {
	tx := block.Transactions()[idx]

	// Check the common receipt fields.
	blockHash, ok := ethReceipt["blockHash"]
	if !ok {
		t.Fatal("blockHash is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, blockHash, kReceipt["blockHash"])

	blockNumber, ok := ethReceipt["blockNumber"]
	if !ok {
		t.Fatal("blockNumber is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, blockNumber.(hexutil.Uint64), hexutil.Uint64(kReceipt["blockNumber"].(*hexutil.Big).ToInt().Uint64()))

	transactionHash, ok := ethReceipt["transactionHash"]
	if !ok {
		t.Fatal("transactionHash is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, transactionHash, kReceipt["transactionHash"])

	transactionIndex, ok := ethReceipt["transactionIndex"]
	if !ok {
		t.Fatal("transactionIndex is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, transactionIndex, hexutil.Uint64(kReceipt["transactionIndex"].(hexutil.Uint)))

	from, ok := ethReceipt["from"]
	if !ok {
		t.Fatal("from is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, from, kReceipt["from"])

	// Kaia transactions that do not use the 'To' field
	// fill in 'To' with from during converting format.
	toInTx := tx.To()
	fromAddress := getFrom(tx)
	to, ok := ethReceipt["to"]
	if !ok {
		t.Fatal("to is not defined in Ethereum transaction receipt format.")
	}
	switch tx.Type() {
	case types.TxTypeAccountUpdate, types.TxTypeFeeDelegatedAccountUpdate, types.TxTypeFeeDelegatedAccountUpdateWithRatio,
		types.TxTypeCancel, types.TxTypeFeeDelegatedCancel, types.TxTypeFeeDelegatedCancelWithRatio,
		types.TxTypeChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		assert.Equal(t, &fromAddress, to)
	default:
		assert.Equal(t, toInTx, to)
	}

	gasUsed, ok := ethReceipt["gasUsed"]
	if !ok {
		t.Fatal("gasUsed is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, gasUsed, kReceipt["gasUsed"])

	// Compare with the calculated cumulative gas used value
	// to check whether the cumulativeGasUsed value is calculated properly.
	cumulativeGasUsed, ok := ethReceipt["cumulativeGasUsed"]
	if !ok {
		t.Fatal("cumulativeGasUsed is not defined in Ethereum transaction receipt format.")
	}
	calculatedCumulativeGas := uint64(0)
	for i := 0; i <= int(idx); i++ {
		calculatedCumulativeGas += receipts[i].GasUsed
	}
	assert.Equal(t, cumulativeGasUsed, hexutil.Uint64(calculatedCumulativeGas))

	contractAddress, ok := ethReceipt["contractAddress"]
	if !ok {
		t.Fatal("contractAddress is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, contractAddress, kReceipt["contractAddress"])

	logs, ok := ethReceipt["logs"]
	if !ok {
		t.Fatal("logs is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, logs, kReceipt["logs"])

	logsBloom, ok := ethReceipt["logsBloom"]
	if !ok {
		t.Fatal("logsBloom is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, logsBloom, kReceipt["logsBloom"])

	typeInt, ok := ethReceipt["type"]
	if !ok {
		t.Fatal("type is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, types.TxType(typeInt.(hexutil.Uint)), types.TxTypeLegacyTransaction)

	effectiveGasPrice, ok := ethReceipt["effectiveGasPrice"]
	if !ok {
		t.Fatal("effectiveGasPrice is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, effectiveGasPrice, hexutil.Uint64(kReceipt["gasPrice"].(*hexutil.Big).ToInt().Uint64()))

	status, ok := ethReceipt["status"]
	if !ok {
		t.Fatal("status is not defined in Ethereum transaction receipt format.")
	}
	assert.Equal(t, status, kReceipt["status"])

	// Check the receipt fields that should be removed.
	var shouldNotExisted []string
	shouldNotExisted = append(shouldNotExisted, "gas", "gasPrice", "senderTxHash", "signatures", "txError", "typeInt", "feePayer", "feePayerSignatures", "feeRatio", "input", "value", "codeFormat", "humanReadable", "key", "inputJSON")
	for i := 0; i < len(shouldNotExisted); i++ {
		k := shouldNotExisted[i]
		_, ok = ethReceipt[k]
		if ok {
			t.Fatal(k, " should not be defined in the Ethereum transaction receipt format.")
		}
	}
}

func createTestData(t *testing.T, header *types.Header) (*types.Block, types.Transactions, map[common.Hash]*types.Transaction, map[common.Hash]*types.Receipt, []*types.Receipt) {
	var txs types.Transactions

	gasPrice := big.NewInt(25 * params.Gkei)
	deployData := "0x60806040526000805534801561001457600080fd5b506101ea806100246000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306661abd1461007257806342cbb15c1461009d578063767800de146100c8578063b22636271461011f578063d14e62b814610150575b600080fd5b34801561007e57600080fd5b5061008761017d565b6040518082815260200191505060405180910390f35b3480156100a957600080fd5b506100b2610183565b6040518082815260200191505060405180910390f35b3480156100d457600080fd5b506100dd61018b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561012b57600080fd5b5061014e60048036038101908080356000191690602001909291905050506101b1565b005b34801561015c57600080fd5b5061017b600480360381019080803590602001909291905050506101b4565b005b60005481565b600043905090565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b50565b80600081905550505600a165627a7a7230582053c65686a3571c517e2cf4f741d842e5ee6aa665c96ce70f46f9a594794f11eb0029"
	executeData := "0xa9059cbb0000000000000000000000008a4c9c443bb0645df646a2d5bb55def0ed1e885a0000000000000000000000000000000000000000000000000000000000003039"
	var anchorData []byte

	txHashMap := make(map[common.Hash]*types.Transaction)
	receiptMap := make(map[common.Hash]*types.Receipt)
	var receipts []*types.Receipt

	// Create test data for chainDataAnchoring tx
	{
		dummyBlock := types.NewBlock(&types.Header{}, nil, nil)
		scData, err := types.NewAnchoringDataType0(dummyBlock, 0, uint64(dummyBlock.Transactions().Len()))
		if err != nil {
			t.Fatal(err)
		}
		anchorData, _ = rlp.EncodeToBytes(scData)
	}

	// Make test transactions data
	{
		// TxTypeLegacyTransaction
		values := map[types.TxValueKeyType]interface{}{
			// Simply set the nonce to txs.Len() to have a different nonce for each transaction type.
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyTo:       common.StringToAddress("0xe0680cfce04f80a386f1764a55c833b108770490"),
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(30000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     []byte("0xe3197e8f000000000000000000000000e0bef99b4a22286e27630b485d06c5a147cee931000000000000000000000000158beff8c8cdebd64654add5f6a1d9937e73536c0000000000000000000000000000000000000000000029bb5e7fb6beae32cf8000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000180000000000000000000000000000000000000000000000000001b60fb4614a22e000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000158beff8c8cdebd64654add5f6a1d9937e73536c00000000000000000000000074ba03198fed2b15a51af242b9c63faf3c8f4d3400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeLegacyTransaction, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])

	}
	{
		// TxTypeValueTransfer
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x520af902892196a3449b06ead301daeaf67e77e8"),
			types.TxValueKeyTo:       common.StringToAddress("0xa06fa690d92788cac4953da5f2dfbc4a2b3871db"),
			types.TxValueKeyAmount:   big.NewInt(5),
			types.TxValueKeyGasLimit: uint64(10000000),
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeValueTransferMemo
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0xc05e11f9075d453b4fc87023f815fa6a63f7aa0c"),
			types.TxValueKeyTo:       common.StringToAddress("0xb5a2d79e9228f3d278cb36b5b15930f24fe8bae8"),
			types.TxValueKeyAmount:   big.NewInt(3),
			types.TxValueKeyGasLimit: uint64(20000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     []byte(string("hello")),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeValueTransferMemo, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeAccountUpdate
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      uint64(txs.Len()),
			types.TxValueKeyFrom:       common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:   uint64(20000000),
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: accountkey.NewAccountKeyLegacy(),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeSmartContractDeploy
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         uint64(txs.Len()),
			types.TxValueKeyFrom:          common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        big.NewInt(0),
			types.TxValueKeyGasLimit:      uint64(100000000),
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyData:          common.Hex2Bytes(deployData),
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		r := createReceipt(t, tx, tx.Gas())
		fromAddress, err := tx.From()
		if err != nil {
			t.Fatal(err)
		}
		tx.FillContractAddress(fromAddress, r)
		receiptMap[tx.Hash()] = r
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeSmartContractExecution
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:       common.StringToAddress("0x00ca1eee49a4d2b04e6562222eab95e9ed29c4bf"),
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(50000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     common.Hex2Bytes(executeData),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeSmartContractExecution, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeCancel
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit: uint64(50000000),
			types.TxValueKeyGasPrice: gasPrice,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeCancel, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeChainDataAnchoring
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:        uint64(txs.Len()),
			types.TxValueKeyFrom:         common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:     uint64(50000000),
			types.TxValueKeyGasPrice:     gasPrice,
			types.TxValueKeyAnchoredData: anchorData,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeChainDataAnchoring, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedValueTransfer
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x520af902892196a3449b06ead301daeaf67e77e8"),
			types.TxValueKeyTo:       common.StringToAddress("0xa06fa690d92788cac4953da5f2dfbc4a2b3871db"),
			types.TxValueKeyAmount:   big.NewInt(5),
			types.TxValueKeyGasLimit: uint64(10000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransfer, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedValueTransferMemo
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0xc05e11f9075d453b4fc87023f815fa6a63f7aa0c"),
			types.TxValueKeyTo:       common.StringToAddress("0xb5a2d79e9228f3d278cb36b5b15930f24fe8bae8"),
			types.TxValueKeyAmount:   big.NewInt(3),
			types.TxValueKeyGasLimit: uint64(20000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     []byte(string("hello")),
			types.TxValueKeyFeePayer: common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemo, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx

		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedAccountUpdate
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      uint64(txs.Len()),
			types.TxValueKeyFrom:       common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:   uint64(20000000),
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccountKey: accountkey.NewAccountKeyLegacy(),
			types.TxValueKeyFeePayer:   common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdate, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedSmartContractDeploy
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:         uint64(txs.Len()),
			types.TxValueKeyFrom:          common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:            (*common.Address)(nil),
			types.TxValueKeyAmount:        big.NewInt(0),
			types.TxValueKeyGasLimit:      uint64(100000000),
			types.TxValueKeyGasPrice:      gasPrice,
			types.TxValueKeyData:          common.Hex2Bytes(deployData),
			types.TxValueKeyHumanReadable: false,
			types.TxValueKeyCodeFormat:    params.CodeFormatEVM,
			types.TxValueKeyFeePayer:      common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeploy, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		r := createReceipt(t, tx, tx.Gas())
		fromAddress, err := tx.From()
		if err != nil {
			t.Fatal(err)
		}
		tx.FillContractAddress(fromAddress, r)
		receiptMap[tx.Hash()] = r
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedSmartContractExecution
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:       common.StringToAddress("0x00ca1eee49a4d2b04e6562222eab95e9ed29c4bf"),
			types.TxValueKeyAmount:   big.NewInt(0),
			types.TxValueKeyGasLimit: uint64(50000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyData:     common.Hex2Bytes(executeData),
			types.TxValueKeyFeePayer: common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecution, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedCancel
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:    uint64(txs.Len()),
			types.TxValueKeyFrom:     common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit: uint64(50000000),
			types.TxValueKeyGasPrice: gasPrice,
			types.TxValueKeyFeePayer: common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancel, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedChainDataAnchoring
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:        uint64(txs.Len()),
			types.TxValueKeyFrom:         common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:     uint64(50000000),
			types.TxValueKeyGasPrice:     gasPrice,
			types.TxValueKeyAnchoredData: anchorData,
			types.TxValueKeyFeePayer:     common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedChainDataAnchoring, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx

		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedValueTransferWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x520af902892196a3449b06ead301daeaf67e77e8"),
			types.TxValueKeyTo:                 common.StringToAddress("0xa06fa690d92788cac4953da5f2dfbc4a2b3871db"),
			types.TxValueKeyAmount:             big.NewInt(5),
			types.TxValueKeyGasLimit:           uint64(10000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedValueTransferMemoWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0xc05e11f9075d453b4fc87023f815fa6a63f7aa0c"),
			types.TxValueKeyTo:                 common.StringToAddress("0xb5a2d79e9228f3d278cb36b5b15930f24fe8bae8"),
			types.TxValueKeyAmount:             big.NewInt(3),
			types.TxValueKeyGasLimit:           uint64(20000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               []byte(string("hello")),
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedValueTransferMemoWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedAccountUpdateWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:           uint64(20000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyAccountKey:         accountkey.NewAccountKeyLegacy(),
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedAccountUpdateWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedSmartContractDeployWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:                 (*common.Address)(nil),
			types.TxValueKeyAmount:             big.NewInt(0),
			types.TxValueKeyGasLimit:           uint64(100000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               common.Hex2Bytes(deployData),
			types.TxValueKeyHumanReadable:      false,
			types.TxValueKeyCodeFormat:         params.CodeFormatEVM,
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractDeployWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		r := createReceipt(t, tx, tx.Gas())
		fromAddress, err := tx.From()
		if err != nil {
			t.Fatal(err)
		}
		tx.FillContractAddress(fromAddress, r)
		receiptMap[tx.Hash()] = r
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedSmartContractExecutionWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyTo:                 common.StringToAddress("0x00ca1eee49a4d2b04e6562222eab95e9ed29c4bf"),
			types.TxValueKeyAmount:             big.NewInt(0),
			types.TxValueKeyGasLimit:           uint64(50000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyData:               common.Hex2Bytes(executeData),
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedSmartContractExecutionWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedCancelWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:           uint64(50000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedCancelWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeFeeDelegatedChainDataAnchoringWithRatio
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:              uint64(txs.Len()),
			types.TxValueKeyFrom:               common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			types.TxValueKeyGasLimit:           uint64(50000000),
			types.TxValueKeyGasPrice:           gasPrice,
			types.TxValueKeyAnchoredData:       anchorData,
			types.TxValueKeyFeePayer:           common.StringToAddress("0xa142f7b24a618778165c9b06e15a61f100c51400"),
			types.TxValueKeyFeeRatioOfFeePayer: types.FeeRatio(20),
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeFeeDelegatedChainDataAnchoringWithRatio, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		feePayerSignatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(3), R: big.NewInt(4), S: big.NewInt(5)},
			&types.TxSignature{V: big.NewInt(4), R: big.NewInt(5), S: big.NewInt(6)},
		}
		tx.SetFeePayerSignatures(feePayerSignatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}

	// Create a block which includes all transaction data.
	var block *types.Block
	if header != nil {
		block = types.NewBlock(header, txs, receipts)
	} else {
		block = types.NewBlock(&types.Header{Number: big.NewInt(1)}, txs, nil)
	}

	return block, txs, txHashMap, receiptMap, receipts
}

func createEthereumTypedTestData(t *testing.T, header *types.Header) (*types.Block, types.Transactions, map[common.Hash]*types.Transaction, map[common.Hash]*types.Receipt, []*types.Receipt) {
	var txs types.Transactions

	gasPrice := big.NewInt(25 * params.Gkei)
	deployData := "0x60806040526000805534801561001457600080fd5b506101ea806100246000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306661abd1461007257806342cbb15c1461009d578063767800de146100c8578063b22636271461011f578063d14e62b814610150575b600080fd5b34801561007e57600080fd5b5061008761017d565b6040518082815260200191505060405180910390f35b3480156100a957600080fd5b506100b2610183565b6040518082815260200191505060405180910390f35b3480156100d457600080fd5b506100dd61018b565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561012b57600080fd5b5061014e60048036038101908080356000191690602001909291905050506101b1565b005b34801561015c57600080fd5b5061017b600480360381019080803590602001909291905050506101b4565b005b60005481565b600043905090565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b50565b80600081905550505600a165627a7a7230582053c65686a3571c517e2cf4f741d842e5ee6aa665c96ce70f46f9a594794f11eb0029"
	accessList := types.AccessList{
		types.AccessTuple{
			Address: common.StringToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			StorageKeys: []common.Hash{
				common.HexToHash("0xa145cd642157a5df01f5bc3837a1bb59b3dcefbbfad5ec435919780aebeaba2b"),
				common.HexToHash("0x12e2c26dca2fb2b8879f54a5ea1604924edf0e37965c2be8aa6133b75818da40"),
			},
		},
	}
	chainId := new(big.Int).SetUint64(2019)

	txHashMap := make(map[common.Hash]*types.Transaction)
	receiptMap := make(map[common.Hash]*types.Receipt)
	var receipts []*types.Receipt

	// Make test transactions data
	{
		// TxTypeEthereumAccessList
		to := common.StringToAddress("0xb5a2d79e9228f3d278cb36b5b15930f24fe8bae8")
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      uint64(txs.Len()),
			types.TxValueKeyTo:         &to,
			types.TxValueKeyAmount:     big.NewInt(10),
			types.TxValueKeyGasLimit:   uint64(50000000),
			types.TxValueKeyData:       common.Hex2Bytes(deployData),
			types.TxValueKeyGasPrice:   gasPrice,
			types.TxValueKeyAccessList: accessList,
			types.TxValueKeyChainID:    chainId,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeEthereumAccessList, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(1), R: big.NewInt(2), S: big.NewInt(3)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}
	{
		// TxTypeEthereumDynamicFee
		to := common.StringToAddress("0xb5a2d79e9228f3d278cb36b5b15930f24fe8bae8")
		values := map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      uint64(txs.Len()),
			types.TxValueKeyTo:         &to,
			types.TxValueKeyAmount:     big.NewInt(3),
			types.TxValueKeyGasLimit:   uint64(50000000),
			types.TxValueKeyData:       common.Hex2Bytes(deployData),
			types.TxValueKeyGasTipCap:  gasPrice,
			types.TxValueKeyGasFeeCap:  gasPrice,
			types.TxValueKeyAccessList: accessList,
			types.TxValueKeyChainID:    chainId,
		}
		tx, err := types.NewTransactionWithMap(types.TxTypeEthereumDynamicFee, values)
		assert.Equal(t, nil, err)

		signatures := types.TxSignatures{
			&types.TxSignature{V: big.NewInt(2), R: big.NewInt(3), S: big.NewInt(4)},
		}
		tx.SetSignature(signatures)

		txs = append(txs, tx)
		txHashMap[tx.Hash()] = tx
		// For testing, set GasUsed with tx.Gas()
		receiptMap[tx.Hash()] = createReceipt(t, tx, tx.Gas())
		receipts = append(receipts, receiptMap[tx.Hash()])
	}

	// Create a block which includes all transaction data.
	var block *types.Block
	if header != nil {
		block = types.NewBlock(header, txs, receipts)
	} else {
		block = types.NewBlock(&types.Header{Number: big.NewInt(1)}, txs, nil)
	}

	return block, txs, txHashMap, receiptMap, receipts
}

func createReceipt(t *testing.T, tx *types.Transaction, gasUsed uint64) *types.Receipt {
	rct := types.NewReceipt(uint(0), tx.Hash(), gasUsed)
	rct.Logs = []*types.Log{}
	rct.Bloom = types.Bloom{}
	return rct
}

// MockDatabaseManager is a mock of database.DBManager interface for overriding the ReadTxAndLookupInfo function.
type MockDatabaseManager struct {
	database.DBManager

	txHashMap     map[common.Hash]*types.Transaction
	blockData     *types.Block
	queryFromPool bool
}

// GetTxLookupInfoAndReceipt retrieves a tx and lookup info and receipt for a given transaction hash.
func (dbm *MockDatabaseManager) ReadTxAndLookupInfo(hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	// If queryFromPool, return nil to query from pool after this function
	if dbm.queryFromPool {
		return nil, common.Hash{}, 0, 0
	}

	txFromHashMap := dbm.txHashMap[hash]
	if txFromHashMap == nil {
		return nil, common.Hash{}, 0, 0
	}
	return txFromHashMap, dbm.blockData.Hash(), dbm.blockData.NumberU64(), txFromHashMap.Nonce()
}

// MockWallet is a mock of accounts.Wallet interface for overriding the Accounts function.
type MockWallet struct {
	accounts.Wallet

	accounts []accounts.Account
}

// NewMockWallet prepares accounts based on tx from.
func NewMockWallet(txs types.Transactions) *MockWallet {
	mw := &MockWallet{}

	for _, t := range txs {
		mw.accounts = append(mw.accounts, accounts.Account{Address: getFrom(t)})
	}
	return mw
}

// Accounts implements accounts.Wallet, returning an account list.
func (mw *MockWallet) Accounts() []accounts.Account {
	return mw.accounts
}

// TestEthTransactionArgs_setDefaults tests setDefaults method of EthTransactionArgs.
func TestEthTransactionArgs_setDefaults(t *testing.T) {
	_, mockBackend, _ := testInitForEthApi(t)
	// To clarify the exact scope of this test, it is assumed that the user must fill in the gas.
	// Because when user does not specify gas, it calls estimateGas internally and it requires
	// many backend calls which are not directly related with this test.
	gas := hexutil.Uint64(1000000)
	from := common.HexToAddress("0x2eaad2bf70a070aaa2e007beee99c6148f47718e")
	poolNonce := uint64(1)
	accountNonce := uint64(5)
	to := common.HexToAddress("0x9712f943b296758aaae79944ec975884188d3a96")
	byteCode := common.Hex2Bytes("6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632e64cec114604e5780636057361d146076575b600080fd5b348015605957600080fd5b50606060a0565b6040518082815260200191505060405180910390f35b348015608157600080fd5b50609e6004803603810190808035906020019092919050505060a9565b005b60008054905090565b80600081905550505600a165627a7a723058207783dba41884f73679e167576362b7277f88458815141651f48ca38c25b498f80029")
	unitPrice := new(big.Int).SetUint64(dummyChainConfigForEthAPITest.UnitPrice)
	value := new(big.Int).SetUint64(500)
	testSet := []struct {
		txArgs              EthTransactionArgs
		expectedResult      EthTransactionArgs
		dynamicFeeParamsSet bool
		nonceSet            bool
		chainIdSet          bool
		expectedError       error
	}{
		{
			txArgs: EthTransactionArgs{
				From:                 nil,
				To:                   nil,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         nil,
				MaxPriorityFeePerGas: nil,
				Value:                nil,
				Nonce:                nil,
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult: EthTransactionArgs{
				From:                 nil,
				To:                   nil,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(new(big.Int)),
				Nonce:                (*hexutil.Uint64)(&poolNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(dummyChainConfigForEthAPITest.ChainID),
			},
			dynamicFeeParamsSet: false,
			nonceSet:            false,
			chainIdSet:          false,
			expectedError:       nil,
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             (*hexutil.Big)(unitPrice),
				MaxFeePerGas:         nil,
				MaxPriorityFeePerGas: nil,
				Value:                (*hexutil.Big)(value),
				Nonce:                nil,
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             (*hexutil.Big)(unitPrice),
				MaxFeePerGas:         nil,
				MaxPriorityFeePerGas: nil,
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&poolNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(dummyChainConfigForEthAPITest.ChainID),
			},
			dynamicFeeParamsSet: false,
			nonceSet:            false,
			chainIdSet:          false,
			expectedError:       nil,
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(new(big.Int).SetUint64(1)),
				MaxPriorityFeePerGas: nil,
				Value:                (*hexutil.Big)(value),
				Nonce:                nil,
				Data:                 nil,
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult:      EthTransactionArgs{},
			dynamicFeeParamsSet: false,
			nonceSet:            false,
			chainIdSet:          false,
			expectedError:       fmt.Errorf("only %s is allowed to be used as maxFeePerGas and maxPriorityPerGas", unitPrice.Text(16)),
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         nil,
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                nil,
				Data:                 nil,
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&poolNonce),
				Data:                 nil,
				Input:                nil,
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(dummyChainConfigForEthAPITest.ChainID),
			},
			dynamicFeeParamsSet: false,
			nonceSet:            false,
			chainIdSet:          false,
			expectedError:       nil,
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         nil,
				MaxPriorityFeePerGas: (*hexutil.Big)(new(big.Int).SetUint64(1)),
				Value:                (*hexutil.Big)(value),
				Nonce:                nil,
				Data:                 nil,
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult:      EthTransactionArgs{},
			dynamicFeeParamsSet: false,
			nonceSet:            false,
			chainIdSet:          false,
			expectedError:       fmt.Errorf("only %s is allowed to be used as maxFeePerGas and maxPriorityPerGas", unitPrice.Text(16)),
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             (*hexutil.Big)(unitPrice),
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                nil,
				Data:                 nil,
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult:      EthTransactionArgs{},
			dynamicFeeParamsSet: false,
			nonceSet:            false,
			chainIdSet:          false,
			expectedError:       errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified"),
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                nil,
				Data:                 nil,
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&poolNonce),
				Data:                 nil,
				Input:                nil,
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(dummyChainConfigForEthAPITest.ChainID),
			},
			dynamicFeeParamsSet: true,
			nonceSet:            false,
			chainIdSet:          false,
			expectedError:       nil,
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         nil,
				MaxPriorityFeePerGas: nil,
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&accountNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&accountNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(dummyChainConfigForEthAPITest.ChainID),
			},
			dynamicFeeParamsSet: false,
			nonceSet:            true,
			chainIdSet:          false,
			expectedError:       nil,
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&accountNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              nil,
			},
			expectedResult: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&accountNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(dummyChainConfigForEthAPITest.ChainID),
			},
			dynamicFeeParamsSet: true,
			nonceSet:            true,
			chainIdSet:          false,
			expectedError:       nil,
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&accountNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(new(big.Int).SetUint64(1234)),
			},
			expectedResult: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&accountNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                nil,
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(new(big.Int).SetUint64(1234)),
			},
			dynamicFeeParamsSet: true,
			nonceSet:            true,
			chainIdSet:          true,
			expectedError:       nil,
		},
		{
			txArgs: EthTransactionArgs{
				From:                 &from,
				To:                   &to,
				Gas:                  &gas,
				GasPrice:             nil,
				MaxFeePerGas:         (*hexutil.Big)(unitPrice),
				MaxPriorityFeePerGas: (*hexutil.Big)(unitPrice),
				Value:                (*hexutil.Big)(value),
				Nonce:                (*hexutil.Uint64)(&accountNonce),
				Data:                 (*hexutil.Bytes)(&byteCode),
				Input:                (*hexutil.Bytes)(&[]byte{0x1}),
				AccessList:           nil,
				ChainID:              (*hexutil.Big)(new(big.Int).SetUint64(1234)),
			},
			expectedResult:      EthTransactionArgs{},
			dynamicFeeParamsSet: true,
			nonceSet:            true,
			chainIdSet:          true,
			expectedError:       errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`),
		},
	}
	for _, test := range testSet {
		mockBackend.EXPECT().CurrentBlock().Return(
			types.NewBlockWithHeader(&types.Header{Number: new(big.Int).SetUint64(0)}),
		)
		mockBackend.EXPECT().SuggestTipCap(gomock.Any()).Return(unitPrice, nil)
		mockBackend.EXPECT().SuggestPrice(gomock.Any()).Return(unitPrice, nil)
		if !test.dynamicFeeParamsSet {
			mockBackend.EXPECT().ChainConfig().Return(dummyChainConfigForEthAPITest)
		}
		if !test.nonceSet {
			mockBackend.EXPECT().GetPoolNonce(context.Background(), gomock.Any()).Return(poolNonce)
		}
		if !test.chainIdSet {
			mockBackend.EXPECT().ChainConfig().Return(dummyChainConfigForEthAPITest)
		}
		mockBackend.EXPECT().RPCGasCap().Return(nil)
		txArgs := test.txArgs
		err := txArgs.setDefaults(context.Background(), mockBackend)
		require.Equal(t, test.expectedError, err)
		if err == nil {
			require.Equal(t, test.expectedResult, txArgs)
		}
	}
}

func TestEthAPI_GetRawTransactionByHash(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, txHashMap, _, _ := createEthereumTypedTestData(t, nil)

	// Define queryFromPool for ReadTxAndLookupInfo function return tx from hash map.
	// MockDatabaseManager will initiate data with txHashMap, block and queryFromPool.
	// If queryFromPool is true, MockDatabaseManager will return nil to query transactions from transaction pool,
	// otherwise return a transaction from txHashMap.
	mockDBManager := &MockDatabaseManager{txHashMap: txHashMap, blockData: block, queryFromPool: false}

	// Mock Backend functions.
	mockBackend.EXPECT().ChainDB().Return(mockDBManager).Times(txs.Len())

	for i := 0; i < txs.Len(); i++ {
		rawTx, err := api.GetRawTransactionByHash(context.Background(), txs[i].Hash())
		if err != nil {
			t.Fatal(err)
		}
		prefix := types.TxType(rawTx[0])
		// When get raw transaction by eth namespace API, EthereumTxTypeEnvelope must not be included.
		require.NotEqual(t, types.EthereumTxTypeEnvelope, prefix)
	}

	mockCtrl.Finish()
}

func TestEthAPI_GetRawTransactionByBlockNumberAndIndex(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	block, txs, _, _, _ := createEthereumTypedTestData(t, nil)

	// Mock Backend functions.
	mockBackend.EXPECT().BlockByNumber(gomock.Any(), gomock.Any()).Return(block, nil).Times(txs.Len())

	for i := 0; i < txs.Len(); i++ {
		rawTx := api.GetRawTransactionByBlockNumberAndIndex(context.Background(), rpc.BlockNumber(block.NumberU64()), hexutil.Uint(i))
		prefix := types.TxType(rawTx[0])
		// When get raw transaction by eth namespace API, EthereumTxTypeEnvelope must not be included.
		require.NotEqual(t, types.EthereumTxTypeEnvelope, prefix)
	}

	mockCtrl.Finish()
}

type testChainContext struct {
	header *types.Header
}

func (mc *testChainContext) Engine() consensus.Engine {
	return gxhash.NewFaker()
}

func (mc *testChainContext) GetHeader(common.Hash, uint64) *types.Header {
	return mc.header
}

// Contract C { constructor() { revert("hello"); } }
var codeRevertHello = "0x6080604052348015600f57600080fd5b5060405162461bcd60e51b815260206004820152600560248201526468656c6c6f60d81b604482015260640160405180910390fdfe"

func testEstimateGas(t *testing.T, mockBackend *mock_api.MockBackend, fnEstimateGas func(EthTransactionArgs, *EthStateOverride) (hexutil.Uint64, error)) {
	chainConfig := &params.ChainConfig{}
	chainConfig.IstanbulCompatibleBlock = common.Big0
	chainConfig.LondonCompatibleBlock = common.Big0
	chainConfig.EthTxTypeCompatibleBlock = common.Big0
	chainConfig.MagmaCompatibleBlock = common.Big0
	chainConfig.KoreCompatibleBlock = common.Big0
	chainConfig.ShanghaiCompatibleBlock = common.Big0
	chainConfig.CancunCompatibleBlock = common.Big0
	chainConfig.KaiaCompatibleBlock = common.Big0
	chainConfig.PragueCompatibleBlock = common.Big0
	var (
		// genesis
		account1 = common.HexToAddress("0xaaaa")
		account2 = common.HexToAddress("0xbbbb")
		account3 = common.HexToAddress("0xcccc")
		account4 = common.HexToAddress("0xdddd")
		account5 = common.HexToAddress("0xeeee")
		gspec    = &blockchain.Genesis{Alloc: blockchain.GenesisAlloc{
			account1: {Balance: big.NewInt(params.KAIA * 2)},
			account2: {Balance: common.Big0},
			account3: {Balance: common.Big0, Code: hexutil.MustDecode(codeRevertHello)},
			account4: {Balance: big.NewInt(params.KAIA * 2), Code: append(types.DelegationPrefix, account5.Bytes()...)},
		}, Config: chainConfig}

		// blockchain
		dbm    = database.NewMemoryDBManager()
		db     = state.NewDatabase(dbm)
		block  = gspec.MustCommit(dbm)
		header = block.Header()
		chain  = &testChainContext{header: header}

		// tx arguments
		KAIA     = hexutil.Big(*big.NewInt(params.KAIA))
		mKAIA    = hexutil.Big(*big.NewInt(params.KAIA / 1000))
		KAIA2_1  = hexutil.Big(*big.NewInt(params.KAIA*2 + 1))
		gas1000  = hexutil.Uint64(1000)
		gas40000 = hexutil.Uint64(40000)
		baddata  = hexutil.Bytes(hexutil.MustDecode("0xdeadbeef"))
	)

	any := gomock.Any()
	getStateAndHeader := func(...interface{}) (*state.StateDB, *types.Header, error) {
		// Return a new state for each call because the state is modified by EstimateGas.
		state, err := state.New(block.Root(), db, nil, nil)
		return state, header, err
	}
	getEVM := func(_ context.Context, msg blockchain.Message, state *state.StateDB, header *types.Header, vmConfig vm.Config) (*vm.EVM, func() error, error) {
		// Taken from node/cn/api_backend.go
		vmError := func() error { return nil }
		txContext := blockchain.NewEVMTxContext(msg, header, chainConfig)
		blockContext := blockchain.NewEVMBlockContext(header, chain, nil)
		return vm.NewEVM(blockContext, txContext, state, chainConfig, &vmConfig), vmError, nil
	}
	mockBackend.EXPECT().ChainConfig().Return(chainConfig).AnyTimes()
	mockBackend.EXPECT().RPCGasCap().Return(common.Big0).AnyTimes()
	mockBackend.EXPECT().RPCEVMTimeout().Return(5 * time.Second).AnyTimes()
	mockBackend.EXPECT().StateAndHeaderByNumber(any, any).DoAndReturn(getStateAndHeader).AnyTimes()
	mockBackend.EXPECT().StateAndHeaderByNumberOrHash(any, any).DoAndReturn(getStateAndHeader).AnyTimes()
	mockBackend.EXPECT().GetEVM(any, any, any, any, any).DoAndReturn(getEVM).AnyTimes()
	mockBackend.EXPECT().IsConsoleLogEnabled().Return(false).AnyTimes()

	testcases := []struct {
		args      EthTransactionArgs
		expectErr string
		expectGas uint64
		overrides EthStateOverride
	}{
		{ // simple transfer
			args: EthTransactionArgs{
				From:  &account1,
				To:    &account2,
				Value: &KAIA,
			},
			expectGas: 21000,
		},
		{ // simple transfer with insufficient funds with zero gasPrice
			args: EthTransactionArgs{
				From:  &account2, // sender has 0 KAIA
				To:    &account1,
				Value: &KAIA, // transfer 1 KAIA
			},
			expectErr: "insufficient balance for transfer",
		},
		{ // simple transfer with slightly insufficient funds with zero gasPrice
			// this testcase is to check whether the gas prefunded in EthDoCall is not too much
			args: EthTransactionArgs{
				From:  &account1, // sender has 2 KAIA
				To:    &account2,
				Value: &KAIA2_1, // transfer 2.0000...1 KAIA
			},
			expectErr: "insufficient balance for transfer",
		},
		{ // simple transfer with insufficient funds with nonzero gasPrice
			args: EthTransactionArgs{
				From:     &account2, // sender has 0 KAIA
				To:       &account1,
				Value:    &KAIA, // transfer 1 KAIA
				GasPrice: &mKAIA,
			},
			expectErr: "insufficient funds for transfer",
		},
		{ // simple transfer too high gasPrice
			args: EthTransactionArgs{
				From:     &account1, // sender has 2 KAIA
				To:       &account2,
				Value:    &KAIA,  // transfer 1 KAIA
				GasPrice: &mKAIA, // allowance = (2 - 1) / 0.001 = 1000 gas
			},
			expectErr: "gas required exceeds allowance",
		},
		{ // empty create
			args:      EthTransactionArgs{},
			expectGas: 53000,
		},
		{ // ignore too small gasLimit
			args: EthTransactionArgs{
				Gas: &gas1000,
			},
			expectGas: 53000,
		},
		{ // capped by gasLimit
			args: EthTransactionArgs{
				Gas: &gas40000,
			},
			expectErr: "gas required exceeds allowance",
		},
		{ // fails with VM error
			args: EthTransactionArgs{
				From: &account1,
				Data: &baddata,
			},
			expectErr: "VM error occurs while running smart contract",
		},
		{ // fails with contract revert
			args: EthTransactionArgs{
				From: &account1,
				To:   &account3,
			},
			expectErr: "execution reverted: hello",
		},
		{ // Should be able to send to an EIP-7702 delegated account.
			args: EthTransactionArgs{
				From:  &account1,
				To:    &account4,
				Value: (*hexutil.Big)(big.NewInt(1)),
			},
			expectGas: 21000,
		},
		{ // Should be able to send as EIP-7702 delegated account.
			args: EthTransactionArgs{
				From:  &account4,
				To:    &account2,
				Value: (*hexutil.Big)(big.NewInt(1)),
			},
			expectGas: 21000,
		},
		{ // Should be able to handle EIP-7623.
			args: EthTransactionArgs{
				From: &account1,
				To:   &account2,
				Data: &floorDataGasTestData,
			},
			overrides: EthStateOverride{
				account2: EthOverrideAccount{
					Code: &floorDataGasTestCode,
				},
			},
			expectGas: 25160,
			// We return ErrIntrinsicGas instead of ErrFloorDataGas when floor data gas hits the cap.
			// So EstimateGas will be able to return expected gas instead of this error.
			// expectErr: "insufficient gas for floor data gas cost",
		},
	}

	for i, tc := range testcases {
		gas, err := fnEstimateGas(tc.args, &tc.overrides)
		t.Logf("tc[%02d] = %d %v", i, gas, err)
		if len(tc.expectErr) > 0 {
			require.NotNil(t, err)
			assert.Contains(t, err.Error(), tc.expectErr, i)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.expectGas, uint64(gas), i)
		}
	}
}

func TestEthAPI_EstimateGas(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForEthApi(t)
	defer mockCtrl.Finish()

	testEstimateGas(t, mockBackend, func(args EthTransactionArgs, overrides *EthStateOverride) (hexutil.Uint64, error) {
		return api.EstimateGas(context.Background(), args, nil, overrides)
	})
}
