package headergov

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/stretchr/testify/assert"
)

func TestNewVoteData(t *testing.T) {
	goodVotes := []struct {
		enum  gov.ParamEnum
		value any
	}{
		{enum: gov.GovernanceDeriveShaImpl, value: uint64(2)},
		{enum: gov.GovernanceGoverningNode, value: "000000000000000000000000000abcd000000000"},
		{enum: gov.GovernanceGoverningNode, value: "0x0000000000000000000000000000000000000000"},
		{enum: gov.GovernanceGoverningNode, value: "0x000000000000000000000000000abcd000000000"},
		{enum: gov.GovernanceGoverningNode, value: "0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456"},
		{enum: gov.GovernanceGoverningNode, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{enum: gov.GovernanceGoverningNode, value: common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")},
		{enum: gov.GovernanceGovParamContract, value: "000000000000000000000000000abcd000000000"},
		{enum: gov.GovernanceGovParamContract, value: "0x0000000000000000000000000000000000000000"},
		{enum: gov.GovernanceGovParamContract, value: "0x000000000000000000000000000abcd000000000"},
		{enum: gov.GovernanceGovParamContract, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{enum: gov.GovernanceUnitPrice, value: float64(0.0)},
		{enum: gov.GovernanceUnitPrice, value: float64(25e9)},
		{enum: gov.GovernanceUnitPrice, value: uint64(25e9)},
		{enum: gov.IstanbulCommitteeSize, value: float64(7.0)},
		{enum: gov.IstanbulCommitteeSize, value: uint64(7)},
		{enum: gov.Kip71BaseFeeDenominator, value: uint64(64)},
		{enum: gov.Kip71GasTarget, value: uint64(15000000)},
		{enum: gov.Kip71GasTarget, value: uint64(30000000)},
		{enum: gov.Kip71LowerBoundBaseFee, value: uint64(25000000000)},
		{enum: gov.Kip71MaxBlockGasUsedForBaseFee, value: uint64(84000000)},
		{enum: gov.Kip71UpperBoundBaseFee, value: uint64(750000000000)},
		{enum: gov.RewardKip82Ratio, value: "10/90"},
		{enum: gov.RewardKip82Ratio, value: "20/80"},
		{enum: gov.RewardMintingAmount, value: "0"},
		{enum: gov.RewardMintingAmount, value: "9600000000000000000"},
		{enum: gov.RewardMintingAmount, value: new(big.Int).SetUint64(9.6e18)},
		{enum: gov.RewardRatio, value: "0/0/100"},
		{enum: gov.RewardRatio, value: "0/100/0"},
		{enum: gov.RewardRatio, value: "10/10/80"},
		{enum: gov.RewardRatio, value: "100/0/0"},
		{enum: gov.RewardRatio, value: "30/40/30"},
		{enum: gov.RewardRatio, value: "50/25/25"},
	}

	for _, tc := range goodVotes {
		param := gov.Params[tc.enum]
		t.Run("goodVote/"+param.Name, func(t *testing.T) {
			assert.NotNil(t, NewVoteData(common.Address{}, param.Name, tc.value))
		})
	}

	t.Run("goodVote/validators", func(t *testing.T) {
		assert.NotNil(t, NewVoteData(common.Address{}, "governance.addvalidator", common.Address{}))
		assert.NotNil(t, NewVoteData(common.Address{}, "governance.removevalidator", common.Address{}))
	})

	badVotes := []struct {
		enum  gov.ParamEnum
		value any
	}{
		{enum: gov.GovernanceDeriveShaImpl, value: "2"},
		{enum: gov.GovernanceDeriveShaImpl, value: false},
		{enum: gov.GovernanceDeriveShaImpl, value: float64(-1)},
		{enum: gov.GovernanceDeriveShaImpl, value: float64(0.1)},
		{enum: gov.GovernanceGovernanceMode, value: "ballot"},
		{enum: gov.GovernanceGovernanceMode, value: "none"},
		{enum: gov.GovernanceGovernanceMode, value: "single"},
		{enum: gov.GovernanceGovernanceMode, value: "unexpected"},
		{enum: gov.GovernanceGovernanceMode, value: 0},
		{enum: gov.GovernanceGovernanceMode, value: 1},
		{enum: gov.GovernanceGovernanceMode, value: 2},
		{enum: gov.GovernanceGoverningNode, value: "0x00000000000000000000"},
		{enum: gov.GovernanceGoverningNode, value: "0x000000000000000000000000000xxxx000000000"},
		{enum: gov.GovernanceGoverningNode, value: "address"},
		{enum: gov.GovernanceGoverningNode, value: 0},
		{enum: gov.GovernanceGoverningNode, value: []byte{0}},
		{enum: gov.GovernanceGoverningNode, value: []byte{}},
		{enum: gov.GovernanceGoverningNode, value: false},
		{enum: gov.GovernanceGovParamContract, value: "0x00000000000000000000"},
		{enum: gov.GovernanceGovParamContract, value: "0x000000000000000000000000000xxxx000000000"},
		{enum: gov.GovernanceGovParamContract, value: "address"},
		{enum: gov.GovernanceGovParamContract, value: 0},
		{enum: gov.GovernanceGovParamContract, value: []byte{0}},
		{enum: gov.GovernanceGovParamContract, value: []byte{}},
		{enum: gov.GovernanceGovParamContract, value: false},
		{enum: gov.GovernanceUnitPrice, value: "25000000000"},
		{enum: gov.GovernanceUnitPrice, value: false},
		{enum: gov.GovernanceUnitPrice, value: float64(-10)},
		{enum: gov.GovernanceUnitPrice, value: float64(0.1)},
		{enum: gov.IstanbulEpoch, value: float64(30000.10)},
		{enum: gov.IstanbulCommitteeSize, value: "7"},
		{enum: gov.IstanbulCommitteeSize, value: false},
		{enum: gov.IstanbulCommitteeSize, value: float64(-7)},
		{enum: gov.IstanbulCommitteeSize, value: float64(7.1)},
		{enum: gov.IstanbulCommitteeSize, value: uint64(0)},
		{enum: gov.IstanbulEpoch, value: "bad"},
		{enum: gov.IstanbulEpoch, value: false},
		{enum: gov.IstanbulEpoch, value: float64(30000.00)},
		{enum: gov.IstanbulEpoch, value: uint64(30000)},
		{enum: gov.IstanbulPolicy, value: "RoundRobin"},
		{enum: gov.IstanbulPolicy, value: "WeightedRandom"},
		{enum: gov.IstanbulPolicy, value: "roundrobin"},
		{enum: gov.IstanbulPolicy, value: "sticky"},
		{enum: gov.IstanbulPolicy, value: "weightedrandom"},
		{enum: gov.IstanbulPolicy, value: false},
		{enum: gov.IstanbulPolicy, value: float64(1.0)},
		{enum: gov.IstanbulPolicy, value: float64(1.2)},
		{enum: gov.IstanbulPolicy, value: uint64(0)},
		{enum: gov.IstanbulPolicy, value: uint64(1)},
		{enum: gov.IstanbulPolicy, value: uint64(2)},
		{enum: gov.Kip71BaseFeeDenominator, value: "64"},
		{enum: gov.Kip71BaseFeeDenominator, value: "sixtyfour"},
		{enum: gov.Kip71BaseFeeDenominator, value: 64},
		{enum: gov.Kip71BaseFeeDenominator, value: false},
		{enum: gov.Kip71GasTarget, value: "30000"},
		{enum: gov.Kip71GasTarget, value: 3000},
		{enum: gov.Kip71GasTarget, value: false},
		{enum: gov.Kip71GasTarget, value: true},
		{enum: gov.Kip71LowerBoundBaseFee, value: "250000000"},
		{enum: gov.Kip71LowerBoundBaseFee, value: "test"},
		{enum: gov.Kip71LowerBoundBaseFee, value: 25000000},
		{enum: gov.Kip71LowerBoundBaseFee, value: false},
		{enum: gov.Kip71MaxBlockGasUsedForBaseFee, value: "84000"},
		{enum: gov.Kip71MaxBlockGasUsedForBaseFee, value: 0},
		{enum: gov.Kip71MaxBlockGasUsedForBaseFee, value: 840000},
		{enum: gov.Kip71MaxBlockGasUsedForBaseFee, value: false},
		{enum: gov.Kip71UpperBoundBaseFee, value: "750000"},
		{enum: gov.Kip71UpperBoundBaseFee, value: 7500000},
		{enum: gov.Kip71UpperBoundBaseFee, value: false},
		{enum: gov.Kip71UpperBoundBaseFee, value: true},
		{enum: gov.RewardDeferredTxFee, value: "false"},
		{enum: gov.RewardDeferredTxFee, value: 0},
		{enum: gov.RewardDeferredTxFee, value: 1},
		{enum: gov.RewardDeferredTxFee, value: false},
		{enum: gov.RewardDeferredTxFee, value: true},
		{enum: gov.RewardKip82Ratio, value: "30/30/40"},
		{enum: gov.RewardKip82Ratio, value: "30/80"},
		{enum: gov.RewardKip82Ratio, value: "49.5/50.5"},
		{enum: gov.RewardKip82Ratio, value: "50.5/50.5"},
		{enum: gov.RewardMinimumStake, value: "-1"},
		{enum: gov.RewardMinimumStake, value: "0"},
		{enum: gov.RewardMinimumStake, value: "2000000000000000000000000"},
		{enum: gov.RewardMinimumStake, value: 0},
		{enum: gov.RewardMinimumStake, value: 1.1},
		{enum: gov.RewardMinimumStake, value: 200000000000000},
		{enum: gov.RewardMintingAmount, value: "many"},
		{enum: gov.RewardMintingAmount, value: 96000},
		{enum: gov.RewardMintingAmount, value: false},
		{enum: gov.RewardProposerUpdateInterval, value: "20"},
		{enum: gov.RewardProposerUpdateInterval, value: float64(20.0)},
		{enum: gov.RewardProposerUpdateInterval, value: float64(20.2)},
		{enum: gov.RewardProposerUpdateInterval, value: uint64(20)},
		{enum: gov.RewardRatio, value: "0/0/0"},
		{enum: gov.RewardRatio, value: "30.5/40/29.5"},
		{enum: gov.RewardRatio, value: "30.5/40/30.5"},
		{enum: gov.RewardRatio, value: "30/40/29"},
		{enum: gov.RewardRatio, value: "30/40/31"},
		{enum: gov.RewardRatio, value: "30/70"},
		{enum: gov.RewardRatio, value: 30 / 40 / 30},
		{enum: gov.RewardStakingUpdateInterval, value: "20"},
		{enum: gov.RewardStakingUpdateInterval, value: float64(20.0)},
		{enum: gov.RewardStakingUpdateInterval, value: float64(20.2)},
		{enum: gov.RewardStakingUpdateInterval, value: uint64(20)},
		{enum: gov.RewardUseGiniCoeff, value: "false"},
		{enum: gov.RewardUseGiniCoeff, value: 0},
		{enum: gov.RewardUseGiniCoeff, value: 1},
		{enum: gov.RewardUseGiniCoeff, value: false},
		{enum: gov.RewardUseGiniCoeff, value: true},
	}

	for _, tc := range badVotes {
		param := gov.Params[tc.enum]
		t.Run("badVote/"+param.Name, func(t *testing.T) {
			assert.Nil(t, NewVoteData(common.Address{}, param.Name, tc.value))
		})
	}

	t.Run("badVote/invalidParam", func(t *testing.T) {
		assert.Nil(t, NewVoteData(common.Address{}, "nonexistent.param", "2"))
		assert.Nil(t, NewVoteData(common.Address{}, "nonexistent.param", uint64(100)))
		assert.Nil(t, NewVoteData(common.Address{}, "governance.unitprice", "100"))
	})
}

func TestVoteSerialization(t *testing.T) {
	v1 := common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541")
	v2 := common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")

	tcs := []struct {
		serializedVoteData string
		blockNum           uint64
		voteData           VoteData
	}{
		///// all vote datas.
		{serializedVoteData: "0xf8439452d41ca72af615a1ac3301b0a93efa222ecc754198676f7665726e616e63652e676f7665726e696e676e6f64659452d41ca72af615a1ac3301b0a93efa222ecc7541", blockNum: 1, voteData: NewVoteData(v1, "governance.governingnode", v1)},
		{serializedVoteData: "0xed9452d41ca72af615a1ac3301b0a93efa222ecc7541917265776172642e6b69703832726174696f8533332f3637", blockNum: 2, voteData: NewVoteData(v1, "reward.kip82ratio", "33/67")},
		{serializedVoteData: "0xf39452d41ca72af615a1ac3301b0a93efa222ecc7541976b697037312e6c6f776572626f756e64626173656665658505d21dba00", blockNum: 3, voteData: NewVoteData(v1, "kip71.lowerboundbasefee", uint64(25e9))},
		{serializedVoteData: "0xf39452d41ca72af615a1ac3301b0a93efa222ecc7541976b697037312e7570706572626f756e646261736566656585ae9f7bcc00", blockNum: 4, voteData: NewVoteData(v1, "kip71.upperboundbasefee", uint64(750e9))},
		{serializedVoteData: "0xef9452d41ca72af615a1ac3301b0a93efa222ecc7541986b697037312e6261736566656564656e6f6d696e61746f7264", blockNum: 5, voteData: NewVoteData(v1, "kip71.basefeedenominator", uint64(100))},
		{serializedVoteData: "0xf83e9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e749331303030303030303030303030303030303030", blockNum: 6, voteData: NewVoteData(v1, "reward.mintingamount", big.NewInt(1000000000000000000))},
		// TODO: add govparamcontract from baobab

		///// Real mainnet vote data.
		{serializedVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e6974707269636585ae9f7bcc00", blockNum: 86119166, voteData: NewVoteData(v1, "governance.unitprice", uint64(750000000000))},
		{serializedVoteData: "0xf09452d41ca72af615a1ac3301b0a93efa222ecc754194676f7665726e616e63652e756e69747072696365853a35294400", blockNum: 90355962, voteData: NewVoteData(v1, "governance.unitprice", uint64(250000000000))},
		{serializedVoteData: "0xed9452d41ca72af615a1ac3301b0a93efa222ecc754196697374616e62756c2e636f6d6d697474656573697a651f", blockNum: 95352567, voteData: NewVoteData(v1, "istanbul.committeesize", uint64(31))},
		{serializedVoteData: "0xf83e9452d41ca72af615a1ac3301b0a93efa222ecc7541947265776172642e6d696e74696e67616d6f756e749336343030303030303030303030303030303030", blockNum: 105629058, voteData: NewVoteData(v1, "reward.mintingamount", big.NewInt(6400000000000000000))},
		{serializedVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f34302f3130", blockNum: 105629111, voteData: NewVoteData(v1, "reward.ratio", "50/40/10")},
		{serializedVoteData: "0xeb9452d41ca72af615a1ac3301b0a93efa222ecc75418c7265776172642e726174696f8835302f32302f3330", blockNum: 118753908, voteData: NewVoteData(v1, "reward.ratio", "50/20/30")},
		{serializedVoteData: "0xf8439452d41ca72af615a1ac3301b0a93efa222ecc754198676f7665726e616e63652e676f7665726e696e676e6f646594c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456", blockNum: 126061533, voteData: NewVoteData(v1, "governance.governingnode", v2)},
		{serializedVoteData: "0xef94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45698676f7665726e616e63652e646572697665736861696d706c80", blockNum: 127692621, voteData: NewVoteData(v2, "governance.deriveshaimpl", uint64(0))},
		{serializedVoteData: "0xe994c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568f6b697037312e67617374617267657483e4e1c0", blockNum: 140916059, voteData: NewVoteData(v2, "kip71.gastarget", uint64(15000000))},
		{serializedVoteData: "0xf83a94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4569f6b697037312e6d6178626c6f636b67617375736564666f72626173656665658401c9c380", blockNum: 140916152, voteData: NewVoteData(v2, "kip71.maxblockgasusedforbasefee", uint64(30000000))},
		{serializedVoteData: "0xed94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a45696697374616e62756c2e636f6d6d697474656573697a6532", blockNum: 161809335, voteData: NewVoteData(v2, "istanbul.committeesize", uint64(50))},
		{serializedVoteData: "0xf83e94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a456947265776172642e6d696e74696e67616d6f756e749339363030303030303030303030303030303030", blockNum: 161809370, voteData: NewVoteData(v2, "reward.mintingamount", new(big.Int).SetUint64(9.6e18))},
		{serializedVoteData: "0xeb94c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4568c7265776172642e726174696f8835302f32352f3235", blockNum: 161809416, voteData: NewVoteData(v2, "reward.ratio", "50/25/25")},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("TestCase_block_%d", tc.blockNum), func(t *testing.T) {
			// Test deserialization
			actual, err := DeserializeHeaderVote(hexutil.MustDecode(tc.serializedVoteData))
			assert.NoError(t, err)
			assert.Equal(t, tc.voteData, actual, "DeserializeHeaderVote() failed")

			// Test serialization
			serialized, err := tc.voteData.Serialize()
			assert.NoError(t, err)
			assert.Equal(t, tc.serializedVoteData, hexutil.Encode(serialized), "voteData.Serialize() failed")
		})
	}
}
