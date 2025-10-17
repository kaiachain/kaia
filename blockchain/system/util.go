package system

import (
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
)

// ActiveSystemContracts returns the currently active system contracts at the
// given block number.
func ActiveSystemContracts(c *params.ChainConfig, genesis common.Hash, head *big.Int) map[string]common.Address {
	active := make(map[string]common.Address)

	if head.Cmp(c.OsakaCompatibleBlock) >= 0 {
	}
	if head.Cmp(c.PragueCompatibleBlock) >= 0 {
		active["HISTORY_STORAGE_ADDRESS"] = params.HistoryStorageAddress
	}
	if head.Cmp(c.RandaoCompatibleBlock) >= 0 {
		active["REGISTRY"] = RegistryAddr
	}
	if head.Cmp(c.Kip103CompatibleBlock) >= 0 {
		active["KIP103"] = c.Kip103ContractAddress
	}
	if head.Cmp(c.Kip160CompatibleBlock) >= 0 {
		active["KIP160"] = c.Kip160ContractAddress
	}

	// These contracts are active from genesis.
	if genesis == params.MainnetGenesisHash {
		active["MAINNET_CREDIT"] = MainnetCreditAddr
	}
	active["ADDRESS_BOOK"] = AddressBookAddr

	return active
}
