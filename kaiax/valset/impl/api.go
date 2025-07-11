package impl

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/rpc"
)

// APIs returns the RPC APIs this valset module provides.
func (v *ValsetModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   NewValsetAPI(v),
			Public:    true,
		},
	}
}

type ValsetAPI struct {
	vs *ValsetModule
}

func NewValsetAPI(vs *ValsetModule) *ValsetAPI {
	return &ValsetAPI{vs: vs}
}

// GetCouncil retrieves the list of authorized validators at the specified block.
func (api *ValsetAPI) GetCouncil(number *rpc.BlockNumber) ([]common.Address, error) {
	num, err := api.resolveRpcNumber(number, true)
	if err != nil {
		return nil, err
	}

	return api.vs.GetCouncil(num)
}

func (api *ValsetAPI) GetCouncilSize(number *rpc.BlockNumber) (int, error) {
	council, err := api.GetCouncil(number)
	if err != nil {
		return -1, err
	}
	return len(council), nil
}

func (api *ValsetAPI) GetCommittee(number *rpc.BlockNumber) ([]common.Address, error) {
	// cannot determine the committee of not-yet finalized block because it depends on the round.
	num, err := api.resolveRpcNumber(number, false)
	if err != nil {
		return nil, err
	}
	header := api.vs.Chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, errUnknownBlock
	}
	round := uint64(header.Round())

	return api.vs.GetCommittee(num, round)
}

func (api *ValsetAPI) GetCommitteeSize(number *rpc.BlockNumber) (int, error) {
	committee, err := api.GetCommittee(number)
	if err != nil {
		return -1, err
	}
	return len(committee), nil
}

func (api *ValsetAPI) GetAllRecordsFromRegistry(name string, number rpc.BlockNumber) ([]interface{}, error) {
	bn := big.NewInt(number.Int64())
	if number == rpc.LatestBlockNumber || number == rpc.PendingBlockNumber {
		bn = big.NewInt(api.vs.Chain.CurrentBlock().Number().Int64())
	}

	if api.vs.Chain.Config().IsRandaoForkEnabled(bn) {
		backend := api.newBlockchainContractBackend()
		records, err := system.ReadAllRecordsFromRegistry(backend, name, bn)
		if err != nil {
			return nil, err
		}

		if len(records) == 0 {
			return nil, fmt.Errorf("%s has not been registered", name)
		}

		recordsList := make([]interface{}, len(records))
		for i, record := range records {
			recordsList[i] = map[string]interface{}{"addr": record.Addr, "activation": record.Activation}
		}
		return recordsList, nil
	} else {
		return nil, errors.New("Randao fork is not enabled")
	}
}

func (api *ValsetAPI) GetActiveAddressFromRegistry(name string, number rpc.BlockNumber) (common.Address, error) {
	bn := big.NewInt(number.Int64())
	if number == rpc.LatestBlockNumber || number == rpc.PendingBlockNumber {
		bn = big.NewInt(api.vs.Chain.CurrentBlock().Number().Int64())
	}

	if api.vs.Chain.Config().IsRandaoForkEnabled(bn) {
		backend := api.newBlockchainContractBackend()
		addr, err := system.ReadActiveAddressFromRegistry(backend, name, bn)
		if err != nil {
			return common.Address{}, err
		}

		if addr == (common.Address{}) {
			return common.Address{}, errors.New("no active address for " + name)
		}
		return addr, nil
	} else {
		return common.Address{}, errors.New("Randao fork is not enabled")
	}
}

// newBlockchainContractBackend creates a new blockchain contract backend.
func (api *ValsetAPI) newBlockchainContractBackend() *backends.BlockchainContractBackend {
	if chain, ok := api.vs.Chain.(backends.BlockChainForCaller); ok {
		return backends.NewBlockchainContractBackend(chain, nil, nil)
	}
	return nil
}

// resolveRpcNumber resolves the RPC block number to a uint64.
func (api *ValsetAPI) resolveRpcNumber(number *rpc.BlockNumber, allowPending bool) (uint64, error) {
	headNum := api.vs.Chain.CurrentBlock().NumberU64()
	var num uint64
	if number == nil || *number == rpc.LatestBlockNumber {
		num = headNum
	} else if *number == rpc.PendingBlockNumber {
		num = headNum + 1
	} else {
		num = uint64(number.Int64())
	}

	if num > headNum+1 { // May allow up to head + 1 to query the pending block.
		return 0, errUnknownBlock
	} else if num == headNum+1 && !allowPending {
		return 0, errPendingNotAllowed
	} else {
		return num, nil
	}
}
