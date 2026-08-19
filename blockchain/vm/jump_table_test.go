package vm

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

func TestJumpTableCopy(t *testing.T) {
	tbl := newConstantinopleInstructionSet()
	require.Equal(t, params.SloadGasEIP150, tbl[SLOAD].constantGas)

	deepCopy := copyJumpTable(&tbl)
	deepCopy[SLOAD].constantGas = params.SloadGasEIP2200

	require.Equal(t, params.SloadGasEIP2200, deepCopy[SLOAD].constantGas)
	require.Equal(t, params.SloadGasEIP150, tbl[SLOAD].constantGas)
}

func newTestEVM(t *testing.T, config *params.ChainConfig, vmConfig *Config) *EVM {
	t.Helper()
	statedb, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	require.NoError(t, err)
	blockCtx := BlockContext{
		BlockNumber: big.NewInt(1),
		Time:        big.NewInt(0),
		CanTransfer: func(StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(StateDB, common.Address, common.Address, *big.Int) {},
	}
	return NewEVM(blockCtx, TxContext{}, statedb, config, vmConfig)
}

// preIstanbulConfig pins the interpreter to ConstantinopleInstructionSet.
func preIstanbulConfig() *params.ChainConfig {
	config := params.TestChainConfig.Copy()
	far := big.NewInt(1000000)
	config.IstanbulCompatibleBlock = far
	config.LondonCompatibleBlock = far
	config.KoreCompatibleBlock = far
	config.ShanghaiCompatibleBlock = far
	config.CancunCompatibleBlock = far
	config.PragueCompatibleBlock = far
	config.OsakaCompatibleBlock = far
	config.PermissionlessCompatibleBlock = far
	return config
}

// ExtraEips must apply to the interpreter's own table only. EnableEIP mutates
// operations in place, so without a deep copy it rewrites the package-level table
// and every later interpreter on the same fork inherits the extra EIP.
func TestExtraEipsKeepSharedTableIntact(t *testing.T) {
	config := preIstanbulConfig()

	evm := newTestEVM(t, config, &Config{ExtraEips: []int{2200}})
	require.Equal(t, params.SloadGasEIP2200, evm.Config.JumpTable[SLOAD].constantGas, "ExtraEips was not applied")

	require.Equal(t, params.SloadGasEIP150, ConstantinopleInstructionSet[SLOAD].constantGas)

	victim := newTestEVM(t, config, &Config{})
	require.Equal(t, params.SloadGasEIP150, victim.Config.JumpTable[SLOAD].constantGas)
}

// ChangeGasCostForTest is handed evm.Config.JumpTable, which shares its operations
// with a package-level table. It runs for every Cancun-or-later fixture in the
// tests package, so the leak needs no ExtraEips at all.
func TestChangeGasCostForTestKeepsSharedTableIntact(t *testing.T) {
	config := params.TestChainConfig
	require.True(t, config.Rules(big.NewInt(1)).IsOsaka, "fixture must select OsakaInstructionSet")
	before := OsakaInstructionSet[EXTCODEHASH].constantGas
	require.NotEqual(t, params.WarmStorageReadCostEIP2929, before)

	evm := newTestEVM(t, config, &Config{})
	ChangeGasCostForTest(&evm.Config.JumpTable)

	require.Equal(t, params.WarmStorageReadCostEIP2929, evm.Config.JumpTable[EXTCODEHASH].constantGas)
	require.Equal(t, before, OsakaInstructionSet[EXTCODEHASH].constantGas)
}
