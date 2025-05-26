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
		name  gov.ParamName
		value any
	}{
		{name: gov.GovernanceDeriveShaImpl, value: uint64(2)},
		{name: gov.GovernanceGoverningNode, value: "000000000000000000000000000abcd000000000"},
		{name: gov.GovernanceGoverningNode, value: "0x0000000000000000000000000000000000000000"},
		{name: gov.GovernanceGoverningNode, value: "0x000000000000000000000000000abcd000000000"},
		{name: gov.GovernanceGoverningNode, value: "0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456"},
		{name: gov.GovernanceGoverningNode, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{name: gov.GovernanceGoverningNode, value: common.HexToAddress("0xc0cbe1c770fbce1eb7786bfba1ac2115d5c0a456")},
		{name: gov.GovernanceGovParamContract, value: "000000000000000000000000000abcd000000000"},
		{name: gov.GovernanceGovParamContract, value: "0x0000000000000000000000000000000000000000"},
		{name: gov.GovernanceGovParamContract, value: "0x000000000000000000000000000abcd000000000"},
		{name: gov.GovernanceGovParamContract, value: common.HexToAddress("000000000000000000000000000abcd000000000")},
		{name: gov.GovernanceUnitPrice, value: float64(0.0)},
		{name: gov.GovernanceUnitPrice, value: float64(25e9)},
		{name: gov.GovernanceUnitPrice, value: uint64(25e9)},
		{name: gov.IstanbulCommitteeSize, value: float64(7.0)},
		{name: gov.IstanbulCommitteeSize, value: uint64(7)},
		{name: gov.Kip71BaseFeeDenominator, value: uint64(64)},
		{name: gov.Kip71GasTarget, value: uint64(15000000)},
		{name: gov.Kip71GasTarget, value: uint64(30000000)},
		{name: gov.Kip71LowerBoundBaseFee, value: uint64(25000000000)},
		{name: gov.Kip71MaxBlockGasUsedForBaseFee, value: uint64(84000000)},
		{name: gov.Kip71UpperBoundBaseFee, value: uint64(750000000000)},
		{name: gov.RewardKip82Ratio, value: "10/90"},
		{name: gov.RewardKip82Ratio, value: "20/80"},
		{name: gov.RewardMintingAmount, value: "0"},
		{name: gov.RewardMintingAmount, value: "9600000000000000000"},
		{name: gov.RewardMintingAmount, value: new(big.Int).SetUint64(9.6e18)},
		{name: gov.RewardRatio, value: "0/0/100"},
		{name: gov.RewardRatio, value: "0/100/0"},
		{name: gov.RewardRatio, value: "10/10/80"},
		{name: gov.RewardRatio, value: "100/0/0"},
		{name: gov.RewardRatio, value: "30/40/30"},
		{name: gov.RewardRatio, value: "50/25/25"},
		{name: gov.AddValidator, value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		{name: gov.AddValidator, value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,0x70997970C51812dc3A010C7d01b50e0d17dc79C8"},
		{name: gov.AddValidator, value: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")},
		{name: gov.AddValidator, value: []common.Address{common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")}},
		{name: gov.RemoveValidator, value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		{name: gov.RemoveValidator, value: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266,0x70997970C51812dc3A010C7d01b50e0d17dc79C8"},
		{name: gov.RemoveValidator, value: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")},
		{name: gov.RemoveValidator, value: []common.Address{common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")}},
	}

	for _, tc := range goodVotes {
		t.Run("goodVote/"+string(tc.name), func(t *testing.T) {
			assert.NotNil(t, NewVoteData(common.Address{}, string(tc.name), tc.value))
		})
	}

	badVotes := []struct {
		name  gov.ParamName
		value any
	}{
		{name: gov.GovernanceDeriveShaImpl, value: "2"},
		{name: gov.GovernanceDeriveShaImpl, value: false},
		{name: gov.GovernanceDeriveShaImpl, value: float64(-1)},
		{name: gov.GovernanceDeriveShaImpl, value: float64(0.1)},
		{name: gov.GovernanceGovernanceMode, value: "ballot"},
		{name: gov.GovernanceGovernanceMode, value: "none"},
		{name: gov.GovernanceGovernanceMode, value: "single"},
		{name: gov.GovernanceGovernanceMode, value: "unexpected"},
		{name: gov.GovernanceGovernanceMode, value: 0},
		{name: gov.GovernanceGovernanceMode, value: 1},
		{name: gov.GovernanceGovernanceMode, value: 2},
		{name: gov.GovernanceGoverningNode, value: "0x00000000000000000000"},
		{name: gov.GovernanceGoverningNode, value: "0x000000000000000000000000000xxxx000000000"},
		{name: gov.GovernanceGoverningNode, value: "address"},
		{name: gov.GovernanceGoverningNode, value: 0},
		{name: gov.GovernanceGoverningNode, value: []byte{0}},
		{name: gov.GovernanceGoverningNode, value: []byte{}},
		{name: gov.GovernanceGoverningNode, value: false},
		{name: gov.GovernanceGovParamContract, value: "0x00000000000000000000"},
		{name: gov.GovernanceGovParamContract, value: "0x000000000000000000000000000xxxx000000000"},
		{name: gov.GovernanceGovParamContract, value: "address"},
		{name: gov.GovernanceGovParamContract, value: 0},
		{name: gov.GovernanceGovParamContract, value: []byte{0}},
		{name: gov.GovernanceGovParamContract, value: []byte{}},
		{name: gov.GovernanceGovParamContract, value: false},
		{name: gov.GovernanceUnitPrice, value: "25000000000"},
		{name: gov.GovernanceUnitPrice, value: false},
		{name: gov.GovernanceUnitPrice, value: float64(-10)},
		{name: gov.GovernanceUnitPrice, value: float64(0.1)},
		{name: gov.IstanbulEpoch, value: float64(30000.10)},
		{name: gov.IstanbulCommitteeSize, value: "7"},
		{name: gov.IstanbulCommitteeSize, value: false},
		{name: gov.IstanbulCommitteeSize, value: float64(-7)},
		{name: gov.IstanbulCommitteeSize, value: float64(7.1)},
		{name: gov.IstanbulCommitteeSize, value: uint64(0)},
		{name: gov.IstanbulEpoch, value: "bad"},
		{name: gov.IstanbulEpoch, value: false},
		{name: gov.IstanbulEpoch, value: float64(30000.00)},
		{name: gov.IstanbulEpoch, value: uint64(30000)},
		{name: gov.IstanbulPolicy, value: "RoundRobin"},
		{name: gov.IstanbulPolicy, value: "WeightedRandom"},
		{name: gov.IstanbulPolicy, value: "roundrobin"},
		{name: gov.IstanbulPolicy, value: "sticky"},
		{name: gov.IstanbulPolicy, value: "weightedrandom"},
		{name: gov.IstanbulPolicy, value: false},
		{name: gov.IstanbulPolicy, value: float64(1.0)},
		{name: gov.IstanbulPolicy, value: float64(1.2)},
		{name: gov.IstanbulPolicy, value: uint64(0)},
		{name: gov.IstanbulPolicy, value: uint64(1)},
		{name: gov.IstanbulPolicy, value: uint64(2)},
		{name: gov.Kip71BaseFeeDenominator, value: "64"},
		{name: gov.Kip71BaseFeeDenominator, value: "sixtyfour"},
		{name: gov.Kip71BaseFeeDenominator, value: 64},
		{name: gov.Kip71BaseFeeDenominator, value: false},
		{name: gov.Kip71GasTarget, value: "30000"},
		{name: gov.Kip71GasTarget, value: 3000},
		{name: gov.Kip71GasTarget, value: false},
		{name: gov.Kip71GasTarget, value: true},
		{name: gov.Kip71LowerBoundBaseFee, value: "250000000"},
		{name: gov.Kip71LowerBoundBaseFee, value: "test"},
		{name: gov.Kip71LowerBoundBaseFee, value: 25000000},
		{name: gov.Kip71LowerBoundBaseFee, value: false},
		{name: gov.Kip71MaxBlockGasUsedForBaseFee, value: "84000"},
		{name: gov.Kip71MaxBlockGasUsedForBaseFee, value: 0},
		{name: gov.Kip71MaxBlockGasUsedForBaseFee, value: 840000},
		{name: gov.Kip71MaxBlockGasUsedForBaseFee, value: false},
		{name: gov.Kip71UpperBoundBaseFee, value: "750000"},
		{name: gov.Kip71UpperBoundBaseFee, value: 7500000},
		{name: gov.Kip71UpperBoundBaseFee, value: false},
		{name: gov.Kip71UpperBoundBaseFee, value: true},
		{name: gov.RewardDeferredTxFee, value: "false"},
		{name: gov.RewardDeferredTxFee, value: 0},
		{name: gov.RewardDeferredTxFee, value: 1},
		{name: gov.RewardDeferredTxFee, value: false},
		{name: gov.RewardDeferredTxFee, value: true},
		{name: gov.RewardKip82Ratio, value: "30/30/40"},
		{name: gov.RewardKip82Ratio, value: "30/80"},
		{name: gov.RewardKip82Ratio, value: "49.5/50.5"},
		{name: gov.RewardKip82Ratio, value: "50.5/50.5"},
		{name: gov.RewardMinimumStake, value: "-1"},
		{name: gov.RewardMinimumStake, value: "0"},
		{name: gov.RewardMinimumStake, value: "2000000000000000000000000"},
		{name: gov.RewardMinimumStake, value: 0},
		{name: gov.RewardMinimumStake, value: 1.1},
		{name: gov.RewardMinimumStake, value: 200000000000000},
		{name: gov.RewardMintingAmount, value: "many"},
		{name: gov.RewardMintingAmount, value: 96000},
		{name: gov.RewardMintingAmount, value: false},
		{name: gov.RewardProposerUpdateInterval, value: "20"},
		{name: gov.RewardProposerUpdateInterval, value: float64(20.0)},
		{name: gov.RewardProposerUpdateInterval, value: float64(20.2)},
		{name: gov.RewardProposerUpdateInterval, value: uint64(20)},
		{name: gov.RewardRatio, value: "0/0/0"},
		{name: gov.RewardRatio, value: "30.5/40/29.5"},
		{name: gov.RewardRatio, value: "30.5/40/30.5"},
		{name: gov.RewardRatio, value: "30/40/29"},
		{name: gov.RewardRatio, value: "30/40/31"},
		{name: gov.RewardRatio, value: "30/70"},
		{name: gov.RewardRatio, value: 30 / 40 / 30},
		{name: gov.RewardStakingUpdateInterval, value: "20"},
		{name: gov.RewardStakingUpdateInterval, value: float64(20.0)},
		{name: gov.RewardStakingUpdateInterval, value: float64(20.2)},
		{name: gov.RewardStakingUpdateInterval, value: uint64(20)},
		{name: gov.RewardUseGiniCoeff, value: "false"},
		{name: gov.RewardUseGiniCoeff, value: 0},
		{name: gov.RewardUseGiniCoeff, value: 1},
		{name: gov.RewardUseGiniCoeff, value: false},
		{name: gov.RewardUseGiniCoeff, value: true},
		{name: gov.AddValidator, value: "0x"},
		{name: gov.AddValidator, value: "0x1"},
		{name: gov.AddValidator, value: "0x1,0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
		{name: gov.RemoveValidator, value: "0x"},
		{name: gov.RemoveValidator, value: "0x1"},
		{name: gov.RemoveValidator, value: "0x1,0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"},
	}

	for _, tc := range badVotes {
		t.Run("badVote/"+string(tc.name), func(t *testing.T) {
			assert.Nil(t, NewVoteData(common.Address{}, string(tc.name), tc.value))
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
	v3 := common.HexToAddress("0x419bC3D20E71b711314D8E65F66161f97e1a73F0")

	tcs := []struct {
		serializedVoteData string
		blockNum           uint64
		voteData           VoteData
		noSerialize        bool
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

		// Real mainnet validator vote data.
		{serializedVoteData: "0xf84594c0cbe1c770fbce1eb7786bfba1ac2115d5c0a4569a676f7665726e616e63652e72656d6f766576616c696461746f729452c0f3654e9ac47ba5e64ffcb398be485718a74b", blockNum: 127293857, voteData: NewVoteData(v2, "governance.removevalidator", "0x52c0f3654e9ac47ba5e64ffcb398be485718a74b")},
		{serializedVoteData: "0xf8459452d41ca72af615a1ac3301b0a93efa222ecc75419a676f7665726e616e63652e72656d6f766576616c696461746f7294a2ba8f7798649a778a1fd66d3035904949fec555", blockNum: 5505383, voteData: NewVoteData(v1, "governance.removevalidator", "0xa2ba8f7798649a778a1fd66d3035904949fec555")},

		// Dummy validator vote data (multiple addresses).
		{serializedVoteData: "0xf85694419bc3d20e71b711314d8e65f66161f97e1a73f097676f7665726e616e63652e61646476616c696461746f72a852c0f3654e9ac47ba5e64ffcb398be485718a74b52c0f3654e9ac47ba5e64ffcb398be485718a74c", blockNum: 127293858, voteData: NewVoteData(v3, "governance.addvalidator", "0x52c0f3654e9ac47ba5e64ffcb398be485718a74b,0x52c0f3654e9ac47ba5e64ffcb398be485718a74c")},
		{serializedVoteData: "0xf85994f39fd6e51aad88f6f4ce6ab8827279cfffb922669a676f7665726e616e63652e72656d6f766576616c696461746f72a83c44cdddb6a900fa2b585dd299e03d12fa4293bc90f79bf6eb2c4f870365e785982e1f101e93b906", blockNum: 123123, voteData: NewVoteData(common.HexToAddress("0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"), "governance.removevalidator", "0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc,0x90f79bf6eb2c4f870365e785982e1f101e93b906")},

		// Real mainnet validator vote data with hex-encoded address (for backward compatibility). Deserialize test only.
		{serializedVoteData: "0xf85b9452d41ca72af615a1ac3301b0a93efa222ecc75419a676f7665726e616e63652e72656d6f766576616c696461746f72aa307831366331393235383561306162323462353532373833623462663764386463396636383535633335", blockNum: 90915008, voteData: NewVoteData(v1, "governance.removevalidator", "0x3833623462663764386463396636383535633335"), noSerialize: true},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("TestCase_block_%d", tc.blockNum), func(t *testing.T) {
			// Test deserialization
			var vb VoteBytes = hexutil.MustDecode(tc.serializedVoteData)
			actual, err := vb.ToVoteData()
			assert.NoError(t, err)
			assert.Equal(t, tc.voteData, actual, "ToVoteData() failed")

			if !tc.noSerialize {
				// Test serialization
				serialized, err := tc.voteData.ToVoteBytes()
				assert.NoError(t, err)
				assert.Equal(t, tc.serializedVoteData, hexutil.Encode(serialized), "voteData.Serialize() failed")
			}
		})
	}
}
