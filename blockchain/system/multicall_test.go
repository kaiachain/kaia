package system

import (
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func TestContractCallerForMultiCall(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	backend := backends.NewSimulatedBackend(nil)
	originCode := MultiCallCode
	defer func() {
		MultiCallCode = originCode
	}()

	// Temporary code injection
	MultiCallCode = MultiCallMockCode

	header := backend.BlockChain().CurrentHeader()
	state, _ := backend.BlockChain().StateAt(header.Root)
	chain := backend.BlockChain()

	caller, _ := NewMultiCallContractCaller(state, chain, header)
	ret, err := caller.MultiCallStakingInfo(&bind.CallOpts{BlockNumber: header.Number})
	assert.Nil(t, err)

	// MultiCall Code injected
	assert.Equal(t, state.GetCode(MultiCallAddr), MultiCallMockCode)

	// Does not affect the original state
	state, _ = backend.BlockChain().StateAt(header.Root)
	assert.Equal(t, []byte(nil), state.GetCode(MultiCallAddr))

	// Mock data
	assert.Equal(t, 5, len(ret.TypeList))
	assert.Equal(t, 5, len(ret.AddressList))
	assert.Equal(t, 1, len(ret.StakingAmounts))

	expectedAddress := []common.Address{
		common.HexToAddress("0x0000000000000000000000000000000000000F00"),
		common.HexToAddress("0x0000000000000000000000000000000000000F01"),
		common.HexToAddress("0x0000000000000000000000000000000000000F02"),
		common.HexToAddress("0x0000000000000000000000000000000000000F03"),
		common.HexToAddress("0x0000000000000000000000000000000000000F04"),
	}
	for i := 0; i < 5; i++ {
		assert.Equal(t, uint8(i), ret.TypeList[i])
		assert.Equal(t, expectedAddress[i], ret.AddressList[i])
	}
	assert.Equal(t, new(big.Int).Mul(big.NewInt(7_000_000), big.NewInt(params.KLAY)), ret.StakingAmounts[0])
}
