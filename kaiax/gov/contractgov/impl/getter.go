package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	govcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/gov"
	"github.com/kaiachain/kaia/kaiax/gov"
)

// GetParamSet returns default parameter set in case of the following errors:
// (1) contractgov is disabled (i.e., pre-Kore or GovParam address is zero)
// (2) GovParam address is not set
// (3) Contract call to GovParam failed
// Invalid parameters in the contract (i.e., invalid parameter name or non-canonical value) are ignored.
func (c *contractGovModule) GetParamSet(blockNum uint64) gov.ParamSet {
	m, err := c.contractGetAllParamsAt(blockNum)
	if err != nil {
		return *gov.GetDefaultGovernanceParamSet()
	}

	ret := *gov.GetDefaultGovernanceParamSet()
	for k, v := range m {
		err = ret.Set(k, v)
		if err != nil {
			return *gov.GetDefaultGovernanceParamSet()
		}
	}

	return ret
}

func (c *contractGovModule) GetPartialParamSet(blockNum uint64) gov.PartialParamSet {
	m, err := c.contractGetAllParamsAt(blockNum)
	if err != nil {
		return nil
	}
	return m
}

func (c *contractGovModule) contractGetAllParamsAt(blockNum uint64) (gov.PartialParamSet, error) {
	addr, err := c.contractAddrAt(blockNum)
	if err != nil {
		return nil, err
	}
	if common.EmptyAddress(addr) {
		logger.Trace("ContractEngine disabled: GovParamContract address not set")
		return nil, nil
	}

	return c.contractGetAllParamsAtFromAddr(blockNum, addr)
}

func (c *contractGovModule) contractGetAllParamsAtFromAddr(blockNum uint64, addr common.Address) (gov.PartialParamSet, error) {
	chain := c.Chain
	if chain == nil {
		return nil, ErrNotReady
	}

	config := c.ChainConfig
	if !config.IsKoreForkEnabled(new(big.Int).SetUint64(blockNum)) {
		return nil, ErrNotReady
	}

	var (
		// Get storage root for this contract at the parent block
		storageRoot = c.getParentStorageRootHash(blockNum, addr)

		// Create composite cache key: contract address + storage root (64 bytes)
		cacheKey [64]byte
	)

	addrHash := addr.Hash()
	copy(cacheKey[:32], addrHash.Bytes())
	copy(cacheKey[32:], storageRoot.Bytes())

	if !common.EmptyHash(storageRoot) {
		// Try to get from cache first
		c.cacheMutex.RLock()
		if cached, exists := c.paramSetCache.Get(cacheKey); exists {
			c.cacheMutex.RUnlock()
			// Cache hit
			cacheHits.Inc(1)
			return cached, nil
		}
		c.cacheMutex.RUnlock()

		// Cache miss
		cacheMisses.Inc(1)
	}

	// compute the result
	result, err := c.computeContractParams(blockNum, addr)
	if err != nil {
		return nil, err
	}

	if !common.EmptyHash(storageRoot) {
		// Store in cache
		c.cacheMutex.Lock()
		c.paramSetCache.Add(cacheKey, result)
		c.cacheMutex.Unlock()
	}

	return result, nil
}

// computeContractParams computes the contract parameters without caching
func (c *contractGovModule) computeContractParams(blockNum uint64, addr common.Address) (gov.PartialParamSet, error) {
	caller := backends.NewBlockchainContractBackend(c.Chain, nil, nil)
	contract, err := govcontract.NewGovParamCaller(addr, caller)
	if err != nil {
		return nil, err
	}

	names, values, err := contract.GetAllParamsAt(nil, new(big.Int).SetUint64(blockNum))
	if err != nil {
		logger.Warn("ContractEngine disabled: getAllParams call failed", "err", err)
		return nil, nil
	}

	if len(names) != len(values) {
		logger.Warn("ContractEngine disabled: getAllParams result invalid", "len(names)", len(names), "len(values)", len(values))
		return nil, nil
	}

	ret := ParseContractCall(names, values)
	return ret, nil
}

func (c *contractGovModule) contractAddrAt(blockNum uint64) (common.Address, error) {
	headerParams := c.Hgm.GetParamSet(blockNum)
	return headerParams.GovParamContract, nil
}

// getParentStorageRootHash computes the storage root for the ContractGov contract at the parent block
func (c *contractGovModule) getParentStorageRootHash(blockNum uint64, contractAddr common.Address) common.Hash {
	if blockNum == 0 {
		return common.Hash{}
	}

	parent := c.Chain.GetHeaderByNumber(blockNum - 1)
	if parent == nil {
		return common.Hash{}
	}

	state, err := c.Chain.StateAt(parent.Root)
	if err != nil {
		return common.Hash{}
	}

	// Get the storage root for the specific contract
	return state.GetStorageRoot(contractAddr)
}

func ParseContractCall(names []string, values [][]byte) gov.PartialParamSet {
	ret := make(gov.PartialParamSet)
	for i := 0; i < len(names); i++ {
		ret.Add(names[i], values[i])
	}

	return ret
}
