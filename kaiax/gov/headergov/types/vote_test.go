package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	govtypes "github.com/kaiachain/kaia/kaiax/gov/types"
	"github.com/stretchr/testify/assert"
)

func TestNewVoteData(t *testing.T) {
	goodVotes := []struct {
		ty    govtypes.ParamEnum
		value interface{}
	}{
		{ty: govtypes.GovernanceDeriveShaImpl, value: uint64(2)},
		{ty: govtypes.GovernanceGoverningNode, value: "000000000000000000000000000abcd000000000"},
		{ty: govtypes.GovernanceGoverningNode, value: "0x0000000000000000000000000000000000000000"},
		{ty: govtypes.GovernanceGoverningNode, value: "0x000000000000000000000000000abcd000000000"},
		{ty: govtypes.GovernanceGoverningNode, value: "0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456"},
		{ty: govtypes.GovernanceGoverningNode, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{ty: govtypes.GovernanceGoverningNode, value: common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")},
		{ty: govtypes.GovernanceGovParamContract, value: "000000000000000000000000000abcd000000000"},
		{ty: govtypes.GovernanceGovParamContract, value: "0x0000000000000000000000000000000000000000"},
		{ty: govtypes.GovernanceGovParamContract, value: "0x000000000000000000000000000abcd000000000"},
		{ty: govtypes.GovernanceGovParamContract, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{ty: govtypes.GovernanceUnitPrice, value: float64(0.0)},
		{ty: govtypes.GovernanceUnitPrice, value: float64(25e9)},
		{ty: govtypes.GovernanceUnitPrice, value: uint64(25e9)},
		{ty: govtypes.IstanbulCommitteeSize, value: float64(7.0)},
		{ty: govtypes.IstanbulCommitteeSize, value: uint64(7)},
		{ty: govtypes.Kip71BaseFeeDenominator, value: uint64(64)},
		{ty: govtypes.Kip71GasTarget, value: uint64(15000000)},
		{ty: govtypes.Kip71GasTarget, value: uint64(30000000)},
		{ty: govtypes.Kip71LowerBoundBaseFee, value: uint64(25000000000)},
		{ty: govtypes.Kip71MaxBlockGasUsedForBaseFee, value: uint64(84000000)},
		{ty: govtypes.Kip71UpperBoundBaseFee, value: uint64(750000000000)},
		{ty: govtypes.RewardKip82Ratio, value: "10/90"},
		{ty: govtypes.RewardKip82Ratio, value: "20/80"},
		{ty: govtypes.RewardMintingAmount, value: "0"},
		{ty: govtypes.RewardMintingAmount, value: "9600000000000000000"},
		{ty: govtypes.RewardMintingAmount, value: new(big.Int).SetUint64(9.6e18)},
		{ty: govtypes.RewardRatio, value: "0/0/100"},
		{ty: govtypes.RewardRatio, value: "0/100/0"},
		{ty: govtypes.RewardRatio, value: "10/10/80"},
		{ty: govtypes.RewardRatio, value: "100/0/0"},
		{ty: govtypes.RewardRatio, value: "30/40/30"},
		{ty: govtypes.RewardRatio, value: "50/25/25"},
	}

	for _, tc := range goodVotes {
		param := govtypes.Params[tc.ty]
		t.Run("goodVote/"+param.Name, func(t *testing.T) {
			assert.NotNil(t, NewVoteData(common.Address{}, param.Name, tc.value))
		})
	}

	t.Run("goodVote/validators", func(t *testing.T) {
		assert.NotNil(t, NewVoteData(common.Address{}, "governance.addvalidator", common.Address{}))
		assert.NotNil(t, NewVoteData(common.Address{}, "governance.removevalidator", common.Address{}))
	})

	badVotes := []struct {
		ty    govtypes.ParamEnum
		value interface{}
	}{
		{ty: govtypes.GovernanceDeriveShaImpl, value: "2"},
		{ty: govtypes.GovernanceDeriveShaImpl, value: false},
		{ty: govtypes.GovernanceDeriveShaImpl, value: float64(-1)},
		{ty: govtypes.GovernanceDeriveShaImpl, value: float64(0.1)},
		{ty: govtypes.GovernanceGovernanceMode, value: "ballot"},
		{ty: govtypes.GovernanceGovernanceMode, value: "none"},
		{ty: govtypes.GovernanceGovernanceMode, value: "single"},
		{ty: govtypes.GovernanceGovernanceMode, value: "unexpected"},
		{ty: govtypes.GovernanceGovernanceMode, value: 0},
		{ty: govtypes.GovernanceGovernanceMode, value: 1},
		{ty: govtypes.GovernanceGovernanceMode, value: 2},
		{ty: govtypes.GovernanceGoverningNode, value: "0x00000000000000000000"},
		{ty: govtypes.GovernanceGoverningNode, value: "0x000000000000000000000000000xxxx000000000"},
		{ty: govtypes.GovernanceGoverningNode, value: "address"},
		{ty: govtypes.GovernanceGoverningNode, value: 0},
		{ty: govtypes.GovernanceGoverningNode, value: []byte{0}},
		{ty: govtypes.GovernanceGoverningNode, value: []byte{}},
		{ty: govtypes.GovernanceGoverningNode, value: false},
		{ty: govtypes.GovernanceGovParamContract, value: "0x00000000000000000000"},
		{ty: govtypes.GovernanceGovParamContract, value: "0x000000000000000000000000000xxxx000000000"},
		{ty: govtypes.GovernanceGovParamContract, value: "address"},
		{ty: govtypes.GovernanceGovParamContract, value: 0},
		{ty: govtypes.GovernanceGovParamContract, value: []byte{0}},
		{ty: govtypes.GovernanceGovParamContract, value: []byte{}},
		{ty: govtypes.GovernanceGovParamContract, value: false},
		{ty: govtypes.GovernanceUnitPrice, value: "25000000000"},
		{ty: govtypes.GovernanceUnitPrice, value: false},
		{ty: govtypes.GovernanceUnitPrice, value: float64(-10)},
		{ty: govtypes.GovernanceUnitPrice, value: float64(0.1)},
		{ty: govtypes.IstanbulEpoch, value: float64(30000.10)},
		{ty: govtypes.IstanbulCommitteeSize, value: "7"},
		{ty: govtypes.IstanbulCommitteeSize, value: false},
		{ty: govtypes.IstanbulCommitteeSize, value: float64(-7)},
		{ty: govtypes.IstanbulCommitteeSize, value: float64(7.1)},
		{ty: govtypes.IstanbulCommitteeSize, value: uint64(0)},
		{ty: govtypes.IstanbulEpoch, value: "bad"},
		{ty: govtypes.IstanbulEpoch, value: false},
		{ty: govtypes.IstanbulEpoch, value: float64(30000.00)},
		{ty: govtypes.IstanbulEpoch, value: uint64(30000)},
		{ty: govtypes.IstanbulPolicy, value: "RoundRobin"},
		{ty: govtypes.IstanbulPolicy, value: "WeightedRandom"},
		{ty: govtypes.IstanbulPolicy, value: "roundrobin"},
		{ty: govtypes.IstanbulPolicy, value: "sticky"},
		{ty: govtypes.IstanbulPolicy, value: "weightedrandom"},
		{ty: govtypes.IstanbulPolicy, value: false},
		{ty: govtypes.IstanbulPolicy, value: float64(1.0)},
		{ty: govtypes.IstanbulPolicy, value: float64(1.2)},
		{ty: govtypes.IstanbulPolicy, value: uint64(0)},
		{ty: govtypes.IstanbulPolicy, value: uint64(1)},
		{ty: govtypes.IstanbulPolicy, value: uint64(2)},
		{ty: govtypes.Kip71BaseFeeDenominator, value: "64"},
		{ty: govtypes.Kip71BaseFeeDenominator, value: "sixtyfour"},
		{ty: govtypes.Kip71BaseFeeDenominator, value: 64},
		{ty: govtypes.Kip71BaseFeeDenominator, value: false},
		{ty: govtypes.Kip71GasTarget, value: "30000"},
		{ty: govtypes.Kip71GasTarget, value: 3000},
		{ty: govtypes.Kip71GasTarget, value: false},
		{ty: govtypes.Kip71GasTarget, value: true},
		{ty: govtypes.Kip71LowerBoundBaseFee, value: "250000000"},
		{ty: govtypes.Kip71LowerBoundBaseFee, value: "test"},
		{ty: govtypes.Kip71LowerBoundBaseFee, value: 25000000},
		{ty: govtypes.Kip71LowerBoundBaseFee, value: false},
		{ty: govtypes.Kip71MaxBlockGasUsedForBaseFee, value: "84000"},
		{ty: govtypes.Kip71MaxBlockGasUsedForBaseFee, value: 0},
		{ty: govtypes.Kip71MaxBlockGasUsedForBaseFee, value: 840000},
		{ty: govtypes.Kip71MaxBlockGasUsedForBaseFee, value: false},
		{ty: govtypes.Kip71UpperBoundBaseFee, value: "750000"},
		{ty: govtypes.Kip71UpperBoundBaseFee, value: 7500000},
		{ty: govtypes.Kip71UpperBoundBaseFee, value: false},
		{ty: govtypes.Kip71UpperBoundBaseFee, value: true},
		{ty: govtypes.RewardDeferredTxFee, value: "false"},
		{ty: govtypes.RewardDeferredTxFee, value: 0},
		{ty: govtypes.RewardDeferredTxFee, value: 1},
		{ty: govtypes.RewardDeferredTxFee, value: false},
		{ty: govtypes.RewardDeferredTxFee, value: true},
		{ty: govtypes.RewardKip82Ratio, value: "30/30/40"},
		{ty: govtypes.RewardKip82Ratio, value: "30/80"},
		{ty: govtypes.RewardKip82Ratio, value: "49.5/50.5"},
		{ty: govtypes.RewardKip82Ratio, value: "50.5/50.5"},
		{ty: govtypes.RewardMinimumStake, value: "-1"},
		{ty: govtypes.RewardMinimumStake, value: "0"},
		{ty: govtypes.RewardMinimumStake, value: "2000000000000000000000000"},
		{ty: govtypes.RewardMinimumStake, value: 0},
		{ty: govtypes.RewardMinimumStake, value: 1.1},
		{ty: govtypes.RewardMinimumStake, value: 200000000000000},
		{ty: govtypes.RewardMintingAmount, value: "many"},
		{ty: govtypes.RewardMintingAmount, value: 96000},
		{ty: govtypes.RewardMintingAmount, value: false},
		{ty: govtypes.RewardProposerUpdateInterval, value: "20"},
		{ty: govtypes.RewardProposerUpdateInterval, value: float64(20.0)},
		{ty: govtypes.RewardProposerUpdateInterval, value: float64(20.2)},
		{ty: govtypes.RewardProposerUpdateInterval, value: uint64(20)},
		{ty: govtypes.RewardRatio, value: "0/0/0"},
		{ty: govtypes.RewardRatio, value: "30.5/40/29.5"},
		{ty: govtypes.RewardRatio, value: "30.5/40/30.5"},
		{ty: govtypes.RewardRatio, value: "30/40/29"},
		{ty: govtypes.RewardRatio, value: "30/40/31"},
		{ty: govtypes.RewardRatio, value: "30/70"},
		{ty: govtypes.RewardRatio, value: 30 / 40 / 30},
		{ty: govtypes.RewardStakingUpdateInterval, value: "20"},
		{ty: govtypes.RewardStakingUpdateInterval, value: float64(20.0)},
		{ty: govtypes.RewardStakingUpdateInterval, value: float64(20.2)},
		{ty: govtypes.RewardStakingUpdateInterval, value: uint64(20)},
		{ty: govtypes.RewardUseGiniCoeff, value: "false"},
		{ty: govtypes.RewardUseGiniCoeff, value: 0},
		{ty: govtypes.RewardUseGiniCoeff, value: 1},
		{ty: govtypes.RewardUseGiniCoeff, value: false},
		{ty: govtypes.RewardUseGiniCoeff, value: true},
	}

	for _, tc := range badVotes {
		param := govtypes.Params[tc.ty]
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
			actual, err := DeserializeHeaderVote(hexutil.MustDecode(tc.serializedVoteData), tc.blockNum)
			assert.NoError(t, err)
			assert.Equal(t, tc.voteData, actual, "DeserializeHeaderVote() failed")

			// Test serialization
			serialized, err := tc.voteData.Serialize()
			assert.NoError(t, err)
			assert.Equal(t, tc.serializedVoteData, hexutil.Encode(serialized), "voteData.Serialize() failed")
		})
	}
}
