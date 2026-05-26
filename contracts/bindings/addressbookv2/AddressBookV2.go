// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package addressbookv2

import (
	"errors"
	"math/big"
	"strings"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = kaia.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BlsPublicKeyInfo is an auto generated low-level Go binding around an user-defined struct.
type BlsPublicKeyInfo struct {
	PublicKey []byte
	Pop       []byte
}

// GovernanceInfo is an auto generated low-level Go binding around an user-defined struct.
type GovernanceInfo struct {
	NodeId          common.Address
	StakingContract common.Address
	VoterAddress    common.Address
	GcId            *big.Int
}

// NodeInfo is an auto generated low-level Go binding around an user-defined struct.
type NodeInfo struct {
	Manager         common.Address
	StakingContract common.Address
	RewardAddress   common.Address
	VoterAddress    common.Address
	TimeoutAt       *big.Int
	GcId            *big.Int
	BlsInfo         BlsPublicKeyInfo
	Name            string
	Metadata        string
	State           uint8
}

// Profile is an auto generated low-level Go binding around an user-defined struct.
type Profile struct {
	NodeId          common.Address
	StakingContract common.Address
	RewardAddress   common.Address
	TimeoutAt       *big.Int
	State           uint8
}

// AddressBookV2MetaData contains all meta data concerning the AddressBookV2 contract.
var AddressBookV2MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_epochBlockInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CONTRACT_TYPE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_METADATA_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_NODE_BALANCE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_STAKE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SYSTEM_SENDER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"assignGcId\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blsInfo\",\"type\":\"tuple\",\"internalType\":\"structBlsPublicKeyInfo\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadata\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deleteNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"epochBlockInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exit\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"typeList\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"},{\"name\":\"addressList\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllAddressInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllBlsInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeIdList\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"pubkeyList\",\"type\":\"tuple[]\",\"internalType\":\"structBlsPublicKeyInfo[]\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllGovernanceInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"infos\",\"type\":\"tuple[]\",\"internalType\":\"structGovernanceInfo[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesLength\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllProfiles\",\"inputs\":[],\"outputs\":[{\"name\":\"profiles\",\"type\":\"tuple[]\",\"internalType\":\"structProfile[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timeoutAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCfsThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCnInfo\",\"inputs\":[{\"name\":\"_cnNodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigurator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochVACount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFundAddresses\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxCounts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxValActivePausedCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNodeInfo\",\"components\":[{\"name\":\"manager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timeoutAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blsInfo\",\"type\":\"tuple\",\"internalType\":\"structBlsPublicKeyInfo\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadata\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeInfos\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"outputs\":[{\"name\":\"infos\",\"type\":\"tuple[]\",\"internalType\":\"structNodeInfo[]\",\"components\":[{\"name\":\"manager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timeoutAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blsInfo\",\"type\":\"tuple\",\"internalType\":\"structBlsPublicKeyInfo\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadata\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeState\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPfsThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegisteredNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSlotLimits\",\"inputs\":[],\"outputs\":[{\"name\":\"maxSlotAvailable\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minActiveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSlotLimitsFor\",\"inputs\":[{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"maxSlotAvailable\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minActiveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStateCount\",\"inputs\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSuspendedValidators\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSuspender\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimeouts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isActivated\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isConstructed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isUsedAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kirContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"offboard\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pocContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processSystemTransition\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"newStates\",\"type\":\"uint8[]\",\"internalType\":\"enumState[]\"},{\"name\":\"timeoutAts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"epochVACount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"readyCandidate\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"readyValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requirement\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"resume\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeGcId\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"spareContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"suspendValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unreadyCandidate\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unreadyValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unsuspendValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCfsThreshold\",\"inputs\":[{\"name\":\"newCfsThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConfigurator\",\"inputs\":[{\"name\":\"newConfigurator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateIdleTimeout\",\"inputs\":[{\"name\":\"newIdleTimeout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKefAddress\",\"inputs\":[{\"name\":\"newKefAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKifAddress\",\"inputs\":[{\"name\":\"newKifAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKpfAddress\",\"inputs\":[{\"name\":\"newKpfAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateManager\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxCandReadyCount\",\"inputs\":[{\"name\":\"newMaxCandReadyCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxNodeCount\",\"inputs\":[{\"name\":\"newMaxNodeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxValActivePausedCount\",\"inputs\":[{\"name\":\"newMaxValActivePausedCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMetadata\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newMetadata\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauseTimeout\",\"inputs\":[{\"name\":\"newPauseTimeout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePfsThreshold\",\"inputs\":[{\"name\":\"newPfsThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRewardAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newRewardAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateSuspender\",\"inputs\":[{\"name\":\"newSuspender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateVoterAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newVoterAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AddressConfigUpdated\",\"inputs\":[{\"name\":\"configId\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"oldValue\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newValue\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CandidateReadied\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CandidateUnreadied\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochTransitionProcessed\",\"inputs\":[{\"name\":\"epochVACount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GcIdAssigned\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ManagerUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newMetadata\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeCreated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeDeleted\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldRewardAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRewardAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StateChanged\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromState\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumState\"},{\"name\":\"toState\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SystemTransitionProcessed\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"newStates\",\"type\":\"uint8[]\",\"indexed\":false,\"internalType\":\"enumState[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UintConfigUpdated\",\"inputs\":[{\"name\":\"configId\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorSuspended\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorUnsuspended\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorsInitialized\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VoterAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldVoterAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newVoterAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AlreadySuspended\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CnNodeNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FactoryNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GcIdAlreadyAssigned\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GcIdNotAssigned\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientNodeBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LegacyFunctionDeprecated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSuspended\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyConfigurator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyNodeId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySuspender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySystemTx\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PDEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlotsFull\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakingDeployerMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakingTooLow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimeoutExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60c060405230608052348015610013575f80fd5b5060405161604c38038061604c833981016040819052610032916100fb565b60a08190528080610041610049565b505050610112565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100995760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100f85780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b5f6020828403121561010b575f80fd5b5051919050565b60805160a051615efd61014f5f395f818161083e01528181613ed5015261438d01525f818161408f015281816140b8015261428e0152615efd5ff3fe60806040526004361061045c575f3560e01c80637df40c621161023e578063b9f96f4011610138578063d9abb38b116100b5578063e70c38f111610079578063e70c38f114610d8b578063e8868e9f14610d9f578063f0a92ba814610db4578063f2fde38b14610dc8578063ffa1ad7414610de75761045c565b8063d9abb38b14610cf0578063da38d49814610d0f578063e1789d8114610d2e578063e4f0d37c14610d4d578063e59d7a8414610d6c5761045c565b8063cb1c2b5c116100fc578063cb1c2b5c14610c6c578063cf8c6f5214610c8a578063d18c07ab14610c9e578063d267eda514610cbd578063d3b5490714610cdc5761045c565b8063b9f96f4014610bd5578063ba70d01814610bf4578063be535f8b14610c13578063c732e08514610c32578063c9a86af214610c4d5761045c565b80639d8cf08f116101c6578063ad3cb1cc1161018a578063ad3cb1cc14610b35578063b42652e914610b65578063b57873a514610b84578063b756393014610ba3578063b858dd9514610bb65761045c565b80639d8cf08f14610a9a5780639f9e3cba14610ab9578063a41b600014610ad8578063a4c98ada14610af7578063a9ee547214610b165761045c565b80638da5cb5b1161020d5780638da5cb5b14610a2b5780638fabf38914610a3f5780639b7ae5ec14610a535780639d0e234d14610a675780639d0f5ef114610a865761045c565b80637df40c62146109b25780638129fc1c146109d157806387b7b8fd146109e55780638beeb439146109ff5761045c565b8063453e962e1161035a5780635b27b6c9116102d7578063715b208b1161029b578063715b208b1461091f578063766718081461094157806376a67a511461095557806378b84a5c14610974578063793c1946146109935761045c565b80635b27b6c91461088c578063656f5869146108ab5780636968b53f146108ca5780636abd623d146108ec578063715018a61461090b5761045c565b806350a5bb691161031e57806350a5bb69146107dc57806350de2fb3146107fa57806352d1902d14610819578063567b0b6c1461082d578063582115fb146108605761045c565b8063453e962e1461071f578063468e3a7e1461073e5780634a8c1fb41461076d5780634b6a94cc146107865780634f1ef286146107c95761045c565b80631b1a478b116103e857806325cf0943116103ac57806325cf094314610698578063291937f5146106ac5780632aca5091146106c05780632d4ede93146106e1578063394f8899146107005761045c565b80631b1a478b146105e15780631b8f34ca146106005780631ba3fd581461061f57806321d2320014610640578063229bb8231461066c5761045c565b80630a4ff2391161042f5780630a4ff2391461050e5780630b1fe7841461053057806315575d5a14610551578063160370b81461059a5780631865c57d146105bf5761045c565b806303e6689d14610481578063058529fb146104af57806306bb8471146104ce57806307ecec3e146104ef575b348015610467575f80fd5b50604051632053d6b560e11b815260040160405180910390fd5b34801561048c575f80fd5b50610495610dfb565b604080519283526020830191909152015b60405180910390f35b3480156104ba575f80fd5b506104956104c9366004614d38565b610e1b565b3480156104d9575f80fd5b506104ed6104e8366004614d73565b610e38565b005b3480156104fa575f80fd5b506104ed610509366004614d8e565b610f41565b348015610519575f80fd5b50610522610fde565b6040519081526020016104a6565b34801561053b575f80fd5b50610544610ff7565b6040516104a69190614df9565b34801561055c575f80fd5b5061057061056b366004614d73565b611141565b604080516001600160a01b03948516815292841660208401529216918101919091526060016104a6565b3480156105a5575f80fd5b506105ae611185565b6040516104a6959493929190614ec3565b3480156105ca575f80fd5b506105d36111af565b6040516104a6929190614f21565b3480156105ec575f80fd5b506105226105fb366004614f4e565b611210565b34801561060b575f80fd5b506104ed61061a366004614fb0565b611255565b34801561062a575f80fd5b506106336113a2565b6040516104a6919061504a565b34801561064b575f80fd5b506106546113b7565b6040516001600160a01b0390911681526020016104a6565b348015610677575f80fd5b5061068b610686366004614d73565b6113d2565b6040516104a6919061505c565b3480156106a3575f80fd5b506105706113ff565b3480156106b7575f80fd5b50610522611434565b3480156106cb575f80fd5b506106d4611446565b6040516104a6919061506a565b3480156106ec575f80fd5b506104ed6106fb366004614d73565b61155c565b34801561070b575f80fd5b506104ed61071a366004614d8e565b6116f2565b34801561072a575f80fd5b506104ed610739366004614d73565b611893565b348015610749575f80fd5b5061075d610758366004614d73565b6119c4565b60405190151581526020016104a6565b348015610778575f80fd5b50600c5461075d9060ff1681565b348015610791575f80fd5b506107bc6040518060400160405280600b81526020016a41646472657373426f6f6b60a81b81525081565b6040516104a691906150fc565b6104ed6107d7366004615238565b6119f1565b3480156107e7575f80fd5b50600c5461075d90610100900460ff1681565b348015610805575f80fd5b506104ed610814366004614d73565b611a10565b348015610824575f80fd5b50610522611a59565b348015610838575f80fd5b506105227f000000000000000000000000000000000000000000000000000000000000000081565b34801561086b575f80fd5b5061087f61087a366004614d73565b611a74565b6040516104a6919061539d565b348015610897575f80fd5b506104ed6108a6366004614d38565b611db8565b3480156108b6575f80fd5b506104ed6108c5366004614d73565b611def565b3480156108d5575f80fd5b506108de611e5b565b6040516104a69291906153af565b3480156108f7575f80fd5b50600754610654906001600160a01b031681565b348015610916575f80fd5b506104ed6120fb565b34801561092a575f80fd5b5061093361210e565b6040516104a692919061541f565b34801561094c575f80fd5b5061052261245a565b348015610960575f80fd5b506104ed61096f366004614d73565b612463565b34801561097f575f80fd5b506104ed61098e366004614d73565b61252f565b34801561099e575f80fd5b506104ed6109ad366004614d73565b6125a3565b3480156109bd575f80fd5b506104ed6109cc366004614d73565b6125ea565b3480156109dc575f80fd5b506104ed61260d565b3480156109f0575f80fd5b506106546002600160a01b0381565b348015610a0a575f80fd5b50610a1e610a19366004615479565b612b2e565b6040516104a691906154b7565b348015610a36575f80fd5b50610654612ee3565b348015610a4a575f80fd5b50610522612f11565b348015610a5e575f80fd5b50610654612f23565b348015610a72575f80fd5b506104ed610a81366004614d38565b612f3e565b348015610a91575f80fd5b50610495612f60565b348015610aa5575f80fd5b506104ed610ab4366004614d73565b612f7e565b348015610ac4575f80fd5b506104ed610ad3366004614d8e565b612fa0565b348015610ae3575f80fd5b506104ed610af2366004614d73565b613128565b348015610b02575f80fd5b506104ed610b11366004614d38565b61319c565b348015610b21575f80fd5b506104ed610b30366004614d38565b6131bf565b348015610b40575f80fd5b506107bc604051806040016040528060058152602001640352e302e360dc1b81525081565b348015610b70575f80fd5b506104ed610b7f366004614d73565b6131e2565b348015610b8f575f80fd5b506104ed610b9e366004614d73565b6132eb565b348015610bae575f80fd5b506001610522565b348015610bc1575f80fd5b50600654610654906001600160a01b031681565b348015610be0575f80fd5b506104ed610bef366004614d73565b61330e565b348015610bff575f80fd5b506104ed610c0e366004614d38565b61334c565b348015610c1e575f80fd5b506104ed610c2d366004614d73565b61336f565b348015610c3d575f80fd5b50610522678ac7230489e8000081565b348015610c58575f80fd5b506104ed610c67366004614d73565b6133de565b348015610c77575f80fd5b506105226a0422ca8b0a00a42500000081565b348015610c95575f80fd5b50610633613401565b348015610ca9575f80fd5b506104ed610cb8366004614d38565b613416565b348015610cc8575f80fd5b50600554610654906001600160a01b031681565b348015610ce7575f80fd5b50610522613439565b348015610cfb575f80fd5b506104ed610d0a366004614d73565b61344b565b348015610d1a575f80fd5b506104ed610d29366004615519565b61348c565b348015610d39575f80fd5b506104ed610d48366004615600565b613534565b348015610d58575f80fd5b506104ed610d67366004614d73565b6137c8565b348015610d77575f80fd5b506104ed610d86366004614d38565b61383d565b348015610d96575f80fd5b50610495613860565b348015610daa575f80fd5b5061052261080081565b348015610dbf575f80fd5b50610522613880565b348015610dd3575f80fd5b506104ed610de2366004614d73565b613892565b348015610df2575f80fd5b50610522600281565b5f805f610e066138d4565b905080600e015481600f015492509250509091565b5f80610e26836138f8565b610e2f84613931565b91509150915091565b610e40613957565b5f610e496138d4565b90505f6001600160a01b0383165f908152602083905260409020600a015460ff166008811115610e7b57610e7b614dc5565b03610e9957604051634825e09360e01b815260040160405180910390fd5b6001600160a01b0382165f9081526020829052604090206005015415610ed257604051637be80ce960e11b815260040160405180910390fd5b5f816007015f8154610ee3906156da565b91829055506001600160a01b0384165f8181526020858152604091829020600501849055905183815292935090917fe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb910160405180910390a2505050565b81610f4b8161398b565b6001600160a01b038216610f725760405163b4fa3fb360e01b815260040160405180910390fd5b5f610f7b6138d4565b6001600160a01b038086165f8181526020849052604080822080548986166001600160a01b03198216811790925591519596509316938492917f8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda8791a45050505050565b5f610ff2610fea6138d4565b6001016139ce565b905090565b60605f6110026138d4565b90505f611011826001016139ce565b9050806001600160401b0381111561102b5761102b61510e565b60405190808252806020026020018201604052801561108a57816020015b6110776040805160a0810182525f808252602082018190529181018290526060810182905290608082015290565b8152602001906001900390816110495790505b5092505f5b8181101561113b575f6110a560018501836139d7565b6001600160a01b038082165f8181526020888152604091829020825160a081018452938452600181015485169184019190915260028101549093169082015260048201546060820152600a8201549293509091608082019060ff16600881111561111157611111614dc5565b815250868481518110611126576111266156f2565b6020908102919091010152505060010161108f565b50505090565b5f805f805f80611150876139e9565b92509250925080611174576040516342dc2dc560e01b815260040160405180910390fd5b5085945090925090505b9193909250565b60608060605f805f805f805f611199613a61565b939e929d50909b50995090975095505050505050565b6040805160018082528183019092526060915f918291602080830190803683370190505090506111dd612ee3565b815f815181106111ef576111ef6156f2565b6001600160a01b039092166020928302919091019091015292600192509050565b5f6112196138d4565b6008015f83600881111561122f5761122f614dc5565b600881111561124057611240614dc5565b81526020019081526020015f20549050919050565b61125d613c45565b85848114158061126d5750808314155b1561128b5760405163b4fa3fb360e01b815260040160405180910390fd5b5f5b8181101561130c576113048989838181106112aa576112aa6156f2565b90506020020160208101906112bf9190614d73565b8888848181106112d1576112d16156f2565b90506020020160208101906112e69190614f4e565b8787858181106112f8576112f86156f2565b90506020020135613c6c565b60010161128d565b50611315613ecf565b1561135b57816113236138d4565b601001556040518281527fd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e848160529060200160405180910390a15b7fab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e888888886040516113909493929190615706565b60405180910390a15050505050505050565b6060610ff26113af6138d4565b600501613f00565b5f6113c06138d4565b601501546001600160a01b0316919050565b5f6113db6138d4565b6001600160a01b039092165f9081526020929092525060409020600a015460ff1690565b5f805f8061140b6138d4565b601281015460138201546014909201546001600160a01b03918216979282169650169350915050565b5f61143d6138d4565b600a0154905090565b60605f6114516138d4565b90505f611460826001016139ce565b9050806001600160401b0381111561147a5761147a61510e565b6040519080825280602002602001820160405280156114ca57816020015b604080516080810182525f8082526020808301829052928201819052606082015282525f199092019101816114985790505b5092505f5b8181101561113b575f6114e560018501836139d7565b6001600160a01b038082165f8181526020888152604091829020825160808101845293845260018101548516918401919091526003810154909316908201526005820154606082015287519293509091879085908110611547576115476156f2565b602090810291909101015250506001016114cf565b806115668161398b565b611571826001613f0c565b61158e5760405163baf3f0f760e01b815260040160405180910390fd5b5f6115976138d4565b6001600160a01b038481165f908152602083815260408083206001808201546002830154600989018652848720805460ff1990811690915591881687528487208054831690559096168552828520805490961690955593835260088501909152812080549394509192919061160b8361579b565b9091555061161e90506003830185613f41565b506001600160a01b0384165f90815260208390526040812080546001600160a01b031990811682556001820180548216905560028201805482169055600382018054909116905560048101829055600581018290559060068201816116838282614c5f565b611690600183015f614c5f565b506116a09050600883015f614c5f565b6116ad600983015f614c5f565b50600a01805460ff191690556040516001600160a01b038516907f1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad905f90a250505050565b816116fc8161398b565b6001600160a01b0382166117235760405163b4fa3fb360e01b815260040160405180910390fd5b5f61172c6138d4565b6001600160a01b038086165f908152602083815260408083206001810154825163e1a12d3560e01b8152925196975090959394169263e1a12d35926004808401939192918290030181865afa158015611787573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906117ab91906157bb565b6001600160a01b0316146117d257604051638ed87ef960e01b815260040160405180910390fd5b6001600160a01b0384165f90815260098301602052604090205460ff161561180d576040516316a163b960e11b815260040160405180910390fd5b6002810180546001600160a01b039081165f818152600986016020526040808220805460ff19908116909155898516808452828420805490921660011790915585546001600160a01b0319168117909555519193928492908a16917f270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc6970191a4505050505050565b8061189d81613f55565b6118a8826001613f0c565b6118c55760405163baf3f0f760e01b815260040160405180910390fd5b6118ce82613f7e565b6118eb5760405163bf74735560e01b815260040160405180910390fd5b678ac7230489e80000826001600160a01b031631101561191d5760405162b8ec7b60e61b815260040160405180910390fd5b5f6119266138d4565b905080600f01546119376002611210565b106119555760405163848084dd60e01b815260040160405180910390fd5b80600e0154611962610fde565b106119805760405163848084dd60e01b815260040160405180910390fd5b61198c8360025f613c6c565b6040516001600160a01b038416907fb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98905f90a2505050565b5f6119cd6138d4565b6001600160a01b039092165f90815260099290920160205250604090205460ff1690565b6119f9614084565b611a0282614128565b611a0c8282614130565b5050565b611a186141f1565b60035f80516020615ed1833981519152611a33601584614223565b604080516001600160a01b0392831681529185166020830152015b60405180910390a250565b5f611a62614283565b505f80516020615eb183398151915290565b611a7c614c96565b5f611a856138d4565b6001600160a01b0384165f908152602091909152604081209150600a82015460ff166008811115611ab857611ab8614dc5565b03611ad657604051634825e09360e01b815260040160405180910390fd5b604080516101408101825282546001600160a01b0390811682526001840154811660208301526002840154811682840152600384015416606082015260048301546080820152600583015460a082015281518083019092526006830180549192849260c08501929082908290611b4b906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611b77906157d6565b8015611bc25780601f10611b9957610100808354040283529160200191611bc2565b820191905f5260205f20905b815481529060010190602001808311611ba557829003601f168201915b50505050508152602001600182018054611bdb906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611c07906157d6565b8015611c525780601f10611c2957610100808354040283529160200191611c52565b820191905f5260205f20905b815481529060010190602001808311611c3557829003601f168201915b5050505050815250508152602001600882018054611c6f906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611c9b906157d6565b8015611ce65780601f10611cbd57610100808354040283529160200191611ce6565b820191905f5260205f20905b815481529060010190602001808311611cc957829003601f168201915b50505050508152602001600982018054611cff906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611d2b906157d6565b8015611d765780601f10611d4d57610100808354040283529160200191611d76565b820191905f5260205f20905b815481529060010190602001808311611d5957829003601f168201915b5050509183525050600a82015460209091019060ff166008811115611d9d57611d9d614dc5565b6008811115611dae57611dae614dc5565b9052509392505050565b611dc0613957565b60055f80516020615e91833981519152611ddb6011846142cc565b604080519182526020820185905201611a4e565b80611df981613f55565b611e04826004613f0c565b611e215760405163baf3f0f760e01b815260040160405180910390fd5b611e2a82613f7e565b611e475760405163bf74735560e01b815260040160405180910390fd5b611a0c826005611e56856142ed565b613c6c565b6060805f611e676138d4565b90505f611e76826001016139ce565b9050806001600160401b03811115611e9057611e9061510e565b604051908082528060200260200182016040528015611eb9578160200160208202803683370190505b509350806001600160401b03811115611ed457611ed461510e565b604051908082528060200260200182016040528015611f1957816020015b6040805180820190915260608082526020820152815260200190600190039081611ef25790505b5092505f5b818110156120f457611f3360018401826139d7565b858281518110611f4557611f456156f2565b60200260200101906001600160a01b031690816001600160a01b031681525050825f015f868381518110611f7b57611f7b6156f2565b60200260200101516001600160a01b03166001600160a01b031681526020019081526020015f206006016040518060400160405290815f82018054611fbf906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611feb906157d6565b80156120365780601f1061200d57610100808354040283529160200191612036565b820191905f5260205f20905b81548152906001019060200180831161201957829003601f168201915b5050505050815260200160018201805461204f906157d6565b80601f016020809104026020016040519081016040528092919081815260200182805461207b906157d6565b80156120c65780601f1061209d576101008083540402835291602001916120c6565b820191905f5260205f20905b8154815290600101906020018083116120a957829003601f168201915b5050505050815250508482815181106120e1576120e16156f2565b6020908102919091010152600101611f1e565b5050509091565b6121036141f1565b61210c5f614317565b565b600c54606090819060ff16612137575050604080515f8082526020820190815281830190925291565b5f805f805f612144613a61565b84519499509297509095509350915061215e81600361580e565b612169906002615825565b6001600160401b038111156121805761218061510e565b6040519080825280602002602001820160405280156121a9578160200160208202803683370190505b5097506121b781600361580e565b6121c2906002615825565b6001600160401b038111156121d9576121d961510e565b604051908082528060200260200182016040528015612202578160200160208202803683370190505b5096505f805b8281101561238f575f8a8381518110612223576122236156f2565b602002602001019060ff16908160ff1681525050878181518110612249576122496156f2565b602002602001015189838061225d906156da565b94508151811061226f5761226f6156f2565b60200260200101906001600160a01b031690816001600160a01b03168152505060018a83815181106122a3576122a36156f2565b602002602001019060ff16908160ff16815250508681815181106122c9576122c96156f2565b60200260200101518983806122dd906156da565b9450815181106122ef576122ef6156f2565b60200260200101906001600160a01b031690816001600160a01b03168152505060028a8381518110612323576123236156f2565b602002602001019060ff16908160ff1681525050858181518110612349576123496156f2565b602002602001015189838061235d906156da565b94508151811061236f5761236f6156f2565b6001600160a01b0390921660209283029190910190910152600101612208565b5060038982815181106123a4576123a46156f2565b60ff909216602092830291909101909101528388826123c2816156da565b9350815181106123d4576123d46156f2565b60200260200101906001600160a01b031690816001600160a01b0316815250506004898281518110612408576124086156f2565b602002602001019060ff16908160ff16815250508288828151811061242f5761242f6156f2565b60200260200101906001600160a01b031690816001600160a01b031681525050505050505050509091565b5f610ff2614387565b8061246d81613f55565b612478826006613f0c565b6124955760405163baf3f0f760e01b815260040160405180910390fd5b5f61249e6138d4565b90506124ad81601001546138f8565b6124b76007611210565b106124d55760405163848084dd60e01b815260040160405180910390fd5b6124e28160100154613931565b6124ec6006611210565b1161250a5760405163848084dd60e01b815260040160405180910390fd5b5f81600c01544261251b9190615825565b905061252984600783613c6c565b50505050565b6125376143b2565b5f6125406138d4565b905061254f6005820183613f41565b61256c5760405163d33ff8c160e01b815260040160405180910390fd5b6040516001600160a01b038316907f814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e905f90a25050565b806125ad81613f55565b6125b8826007613f0c565b6125d55760405163baf3f0f760e01b815260040160405180910390fd5b6125de826143e6565b611a0c8260065f613c6c565b6125f2613957565b60015f80516020615ed1833981519152611a33601384614223565b5f61261661441f565b805490915060ff600160401b82041615906001600160401b03165f8115801561263c5750825b90505f826001600160401b031660011480156126575750303b155b905081158015612665575080155b156126835760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156126ad57845460ff60401b1916600160401b1785555b60405163e2693e3f60e01b815260206004820152601060248201526f10509d8c91185d1850dbdb9d1c9858dd60821b60448201525f906104019063e2693e3f90606401602060405180830381865afa15801561270b573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061272f91906157bb565b90506001600160a01b0381166127585760405163aed5959560e01b815260040160405180910390fd5b5f816001600160a01b031663ebe58ed76040518163ffffffff1660e01b81526004015f60405180830381865afa158015612794573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526127bb9190810190615aee565b90505f6127c66138d4565b90506127d4825f0151614447565b6060820151600a8201556080820151600b82015560a0820151600c82015560c0820151600d82015560e0820151600e8201556101008201516011820155610120820151600f8201556101408201516012820180546001600160a01b03199081166001600160a01b03938416179091556101608401516013840180548316918416919091179055610180840151601484018054831691841691909117905560208401516015840180548316918416919091179055604084015160168401805490921692169190911790556101a0820151515f5b81811015612a7d575f846101a0015182815181106128c6576128c66156f2565b602002602001015190505f856101c0015183815181106128e8576128e86156f2565b6020908102919091018101516001600160a01b038481165f90815260098901845260408082208054600160ff199182168117909255958501518416835281832080548716821790558185015190931682529020805490931617909155905060066101208201819052505f608082018181526001600160a01b0380851683526020888152604093849020855181549084166001600160a01b03199182161782559186015160018201805491851691841691909117905593850151600285018054918416918316919091179055606085015160038501805491909316911617905551600482015560a0820151600582015560c0820151805183929190600683019081906129f39082615c76565b5060208201516001820190612a089082615c76565b50505060e08201516008820190612a1f9082615c76565b506101008201516009820190612a359082615c76565b50610120820151600a8201805460ff19166001836008811115612a5a57612a5a614dc5565b0217905550612a6f9150506001860183614458565b5050508060010190506128a6565b5060065f9081526008830160205260409081902082905560108301829055606460078401556101a084015190517f820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d95291612ad59161504a565b60405180910390a1505050508315612b2757845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b60605f612b396138d4565b905082806001600160401b03811115612b5457612b5461510e565b604051908082528060200260200182016040528015612b8d57816020015b612b7a614c96565b815260200190600190039081612b725790505b5092505f5b81811015612eda57825f878784818110612bae57612bae6156f2565b9050602002016020810190612bc39190614d73565b6001600160a01b03908116825260208083019390935260409182015f20825161014081018452815483168152600182015483169481019490945260028101548216848401526003810154909116606084015260048101546080840152600581015460a08401528151808301909252600681018054919260c085019290919082908290612c4e906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054612c7a906157d6565b8015612cc55780601f10612c9c57610100808354040283529160200191612cc5565b820191905f5260205f20905b815481529060010190602001808311612ca857829003601f168201915b50505050508152602001600182018054612cde906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054612d0a906157d6565b8015612d555780601f10612d2c57610100808354040283529160200191612d55565b820191905f5260205f20905b815481529060010190602001808311612d3857829003601f168201915b5050505050815250508152602001600882018054612d72906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054612d9e906157d6565b8015612de95780601f10612dc057610100808354040283529160200191612de9565b820191905f5260205f20905b815481529060010190602001808311612dcc57829003601f168201915b50505050508152602001600982018054612e02906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054612e2e906157d6565b8015612e795780601f10612e5057610100808354040283529160200191612e79565b820191905f5260205f20905b815481529060010190602001808311612e5c57829003601f168201915b5050509183525050600a82015460209091019060ff166008811115612ea057612ea0614dc5565b6008811115612eb157612eb1614dc5565b81525050848281518110612ec757612ec76156f2565b6020908102919091010152600101612b92565b50505092915050565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b5f612f1a6138d4565b60110154905090565b5f612f2c6138d4565b601601546001600160a01b0316919050565b612f46613957565b5f5f80516020615e91833981519152611ddb600c846142cc565b5f80612f76612f6d6138d4565b60100154610e1b565b915091509091565b612f86613957565b5f5f80516020615ed1833981519152611a33601284614223565b81612faa8161398b565b5f612fb36138d4565b6001600160a01b038086165f9081526020839052604080822060030180548885166001600160a01b0319821617909155905163e2693e3f60e01b8152939450909116916104019063e2693e3f9061302f906004016020808252600e908201526d29ba30b5b4b733aa3930b1b5b2b960911b604082015260600190565b602060405180830381865afa15801561304a573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061306e91906157bb565b90506001600160a01b038116156130d65760405163aad8cb3f60e01b81526001600160a01b03878116600483015282169063aad8cb3f906024015f604051808303815f87803b1580156130bf575f80fd5b505af11580156130d1573d5f803e3d5ffd5b505050505b846001600160a01b0316826001600160a01b0316876001600160a01b03167f23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d9660405160405180910390a4505050505050565b6131306143b2565b5f6131396138d4565b90506131486005820183614458565b61316557604051633ad2b1bb60e11b815260040160405180910390fd5b6040516001600160a01b038316907fb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd905f90a25050565b6131a4613957565b60045f80516020615e91833981519152611ddb600e846142cc565b6131c7613957565b60065f80516020615e91833981519152611ddb600f846142cc565b806131ec81613f55565b5f6131f6836113d2565b9050600781600881111561320c5761320c614dc5565b0361321f5761321a836143e6565b613251565b600681600881111561323357613233614dc5565b146132515760405163baf3f0f760e01b815260040160405180910390fd5b5f61325a6138d4565b905061326981601001546138f8565b6132736008611210565b106132915760405163848084dd60e01b815260040160405180910390fd5b60068260088111156132a5576132a5614dc5565b036132df576132b78160100154613931565b6132c16006611210565b116132df5760405163848084dd60e01b815260040160405180910390fd5b6125298460085f613c6c565b6132f36141f1565b60045f80516020615ed1833981519152611a33601684614223565b8061331881613f55565b613323826004613f0c565b6133405760405163baf3f0f760e01b815260040160405180910390fd5b611a0c8260015f613c6c565b613354613957565b60025f80516020615e91833981519152611ddb600a846142cc565b613377613957565b5f6133806138d4565b6001600160a01b0383165f908152602082905260408120600501549192508190036133be576040516357024f6d60e11b815260040160405180910390fd5b506001600160a01b039091165f9081526020919091526040812060050155565b6133e6613957565b60025f80516020615ed1833981519152611a33601484614223565b6060610ff261340e6138d4565b600301613f00565b61341e613957565b60035f80516020615e91833981519152611ddb600b846142cc565b5f6134426138d4565b60100154905090565b8061345581613f55565b613460826005613f0c565b61347d5760405163baf3f0f760e01b815260040160405180910390fd5b611a0c826004611e56856142ed565b826134968161398b565b6108008211156134b95760405163b4fa3fb360e01b815260040160405180910390fd5b82826134c36138d4565b6001600160a01b0387165f90815260209190915260409020600901916134ea919083615d31565b50836001600160a01b03167f2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf8484604051613526929190615dea565b60405180910390a250505050565b5f61353e886113d2565b600881111561354f5761354f614dc5565b1461356d5760405163731918fb60e11b815260040160405180910390fd5b81515f0361358e5760405163b4fa3fb360e01b815260040160405180910390fd5b610800815111156135b25760405163b4fa3fb360e01b815260040160405180910390fd5b5f6135bb6138d4565b90506135cd816009018989898861446c565b5f604051806101400160405280336001600160a01b03168152602001896001600160a01b03168152602001886001600160a01b03168152602001876001600160a01b031681526020015f81526020015f81526020018681526020018581526020018481526020016001600881111561364757613647614dc5565b90526001600160a01b03808b165f9081526020858152604091829020845181549085166001600160a01b031991821617825591850151600182018054918616918416919091179055918401516002830180549185169183169190911790556060840151600383018054919094169116179091556080820151600482015560a0820151600582015560c082015180519293508392600683019081906136eb9082615c76565b50602082015160018201906137009082615c76565b50505060e082015160088201906137179082615c76565b50610100820151600982019061372d9082615c76565b50610120820151600a8201805460ff1916600183600881111561375257613752614dc5565b0217905550613767915050600383018a614458565b5060015f9081526008830160205260408120805491613785836156da565b90915550506040516001600160a01b038a16907f55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792905f90a2505050505050505050565b806137d281613f55565b6137dd826002613f0c565b6137fa5760405163baf3f0f760e01b815260040160405180910390fd5b6138068260015f613c6c565b6040516001600160a01b038316907f8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2905f90a25050565b613845613957565b60015f80516020615e91833981519152611ddb600d846142cc565b5f805f61386b6138d4565b905080600c015481600d015492509250509091565b5f6138896138d4565b600b0154905090565b61389a6141f1565b6001600160a01b0381166138c857604051631e4fbdf760e01b81525f60048201526024015b60405180910390fd5b6138d181614317565b50565b7f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff0090565b5f600482101561390957505f919050565b6002613916600384615e2c565b613921906001615825565b61392b9190615e2c565b92915050565b5f600482101561393f575090565b600361394c83600261580e565b613921906002615825565b61395f6138d4565b601601546001600160a01b0316331461210c5760405163033b71e160e41b815260040160405180910390fd5b6139936138d4565b6001600160a01b038281165f9081526020929092526040909120541633146138d15760405163605919ad60e11b815260040160405180910390fd5b5f61392b825490565b5f6139e2838361455a565b9392505050565b5f805f806139f56138d4565b6001600160a01b0386165f908152602082905260408120919250600a82015460ff166008811115613a2857613a28614dc5565b03613a3d575f805f945094509450505061117e565b6001818101546002909201546001600160a01b039283169892169650945092505050565b60608060605f805f613a716138d4565b90505f613a80826001016139ce565b9050806001600160401b03811115613a9a57613a9a61510e565b604051908082528060200260200182016040528015613ac3578160200160208202803683370190505b509650806001600160401b03811115613ade57613ade61510e565b604051908082528060200260200182016040528015613b07578160200160208202803683370190505b509550806001600160401b03811115613b2257613b2261510e565b604051908082528060200260200182016040528015613b4b578160200160208202803683370190505b5094505f5b81811015613c1d575f613b6660018501836139d7565b6001600160a01b0381165f9081526020869052604090208a519192509082908b9085908110613b9757613b976156f2565b6001600160a01b03928316602091820292909201015260018201548a519116908a9085908110613bc957613bc96156f2565b6001600160a01b03928316602091820292909201015260028201548951911690899085908110613bfb57613bfb6156f2565b6001600160a01b03909216602092830291909101909101525050600101613b50565b505060138101546012909101549596949593946001600160a01b039182169490911692509050565b336002600160a01b031461210c576040516354d325c360e01b815260040160405180910390fd5b5f613c756138d4565b6001600160a01b0385165f908152602082905260408120600a8101549293509160ff1690816008811115613cab57613cab614dc5565b1480613cc757505f856008811115613cc557613cc5614dc5565b145b80613cf35750846008811115613cdf57613cdf614dc5565b816008811115613cf157613cf1614dc5565b145b15613d0057505050505050565b826008015f826008811115613d1757613d17614dc5565b6008811115613d2857613d28614dc5565b81526020019081526020015f205f815480929190613d459061579b565b9190505550826008015f866008811115613d6157613d61614dc5565b6008811115613d7257613d72614dc5565b81526020019081526020015f205f815480929190613d8f906156da565b9091555060019050816008811115613da957613da9614dc5565b148015613dc857506001856008811115613dc557613dc5614dc5565b14155b15613dee57613dda6001840187614458565b50613de86003840187613f41565b50613e43565b6001816008811115613e0257613e02614dc5565b14158015613e2157506001856008811115613e1f57613e1f614dc5565b145b15613e4357613e336001840187613f41565b50613e416003840187614458565b505b600a8201805486919060ff19166001836008811115613e6457613e64614dc5565b021790555060048201849055846008811115613e8257613e82614dc5565b816008811115613e9457613e94614dc5565b6040516001600160a01b038916907fcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee905f90a4505050505050565b5f613efa7f000000000000000000000000000000000000000000000000000000000000000043615e3f565b15919050565b60605f6139e283614580565b5f816008811115613f1f57613f1f614dc5565b613f28846113d2565b6008811115613f3957613f39614dc5565b149392505050565b5f6139e2836001600160a01b0384166145d9565b336001600160a01b038216146138d1576040516335f1334d60e11b815260040160405180910390fd5b5f80613f886138d4565b6001600160a01b038085165f9081526020928352604080822060010154815163318588a360e11b81529151931694509092849263630b11469260048082019392918290030181865afa158015613fe0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906140049190615e52565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015614040573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906140649190615e52565b61406e9190615e69565b6a0422ca8b0a00a4250000001115949350505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061410a57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166140fe5f80516020615eb1833981519152546001600160a01b031690565b6001600160a01b031614155b1561210c5760405163703e46dd60e11b815260040160405180910390fd5b6138d16141f1565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561418a575060408051601f3d908101601f1916820190925261418791810190615e52565b60015b6141b257604051634c9c8ce360e01b81526001600160a01b03831660048201526024016138bf565b5f80516020615eb183398151915281146141e257604051632a87526960e21b8152600481018290526024016138bf565b6141ec83836146c3565b505050565b336141fa612ee3565b6001600160a01b03161461210c5760405163118cdaa760e01b81523360048201526024016138bf565b5f6001600160a01b03821661424b5760405163b4fa3fb360e01b815260040160405180910390fd5b5f614276847f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff00615825565b8054939055509092915050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461210c5760405163703e46dd60e11b815260040160405180910390fd5b5f815f0361424b5760405163b4fa3fb360e01b815260040160405180910390fd5b5f6142f66138d4565b6001600160a01b039092165f90815260209290925250604090206004015490565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f610ff27f000000000000000000000000000000000000000000000000000000000000000043615e2c565b6143ba6138d4565b601501546001600160a01b0316331461210c5760405163333f4e6560e01b815260040160405180910390fd5b5f6143f0826142ed565b905080158015906144015750804210155b15611a0c5760405163b48d5fc760e01b815260040160405180910390fd5b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061392b565b61444f614718565b6138d18161473d565b5f6139e2836001600160a01b038416614745565b61447884848484614791565b5f6144816148fd565b90506001600160a01b0381166144aa5760405163cdded31d60e01b815260040160405180910390fd5b60405163669d8d4560e01b81526001600160a01b03858116600483015233919083169063669d8d4590602401602060405180830381865afa1580156144f1573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061451591906157bb565b6001600160a01b03161461453c57604051632281776f60e01b815260040160405180910390fd5b614546848461497f565b61455286868686614a2a565b505050505050565b5f825f01828154811061456f5761456f6156f2565b905f5260205f200154905092915050565b6060815f018054806020026020016040519081016040528092919081815260200182805480156145cd57602002820191905f5260205f20905b8154815260200190600101908083116145b9575b50505050509050919050565b5f81815260018301602052604081205480156146b3575f6145fb600183615e69565b85549091505f9061460e90600190615e69565b905080821461466d575f865f01828154811061462c5761462c6156f2565b905f5260205f200154905080875f01848154811061464c5761464c6156f2565b5f918252602080832090910192909255918252600188019052604090208390555b855486908061467e5761467e615e7c565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f90556001935050505061392b565b5f91505061392b565b5092915050565b6146cc82614af4565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115614710576141ec8282614b57565b611a0c614bf0565b614720614c0f565b61210c57604051631afcd79f60e31b815260040160405180910390fd5b61389a614718565b5f81815260018301602052604081205461478a57508154600181810184555f84815260208082209093018490558454848252828601909352604090209190915561392b565b505f61392b565b6001600160a01b03841615806147ae57506001600160a01b038316155b806147c057506001600160a01b038216155b156147de5760405163b4fa3fb360e01b815260040160405180910390fd5b826001600160a01b0316846001600160a01b0316148061480f5750816001600160a01b0316846001600160a01b0316145b8061482b5750816001600160a01b0316836001600160a01b0316145b156148495760405163b4fa3fb360e01b815260040160405180910390fd5b80515160301415806148615750806020015151606014155b1561487f5760405163b4fa3fb360e01b815260040160405180910390fd5b805180516020909101207fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29314806148df575060208082015180519101207f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21145b156125295760405163b4fa3fb360e01b815260040160405180910390fd5b60405163e2693e3f60e01b815260206004820152601060248201526f436e5374616b696e67466163746f727960801b60448201525f906104019063e2693e3f90606401602060405180830381865afa15801561495b573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ff291906157bb565b5f826001600160a01b031663e1a12d356040518163ffffffff1660e01b8152600401602060405180830381865afa1580156149bc573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906149e091906157bb565b90506001600160a01b03811615801590614a0c5750816001600160a01b0316816001600160a01b031614155b156141ec5760405163b4fa3fb360e01b815260040160405180910390fd5b6001600160a01b0383165f9081526020859052604090205460ff1680614a6757506001600160a01b0382165f9081526020859052604090205460ff165b80614a8957506001600160a01b0381165f9081526020859052604090205460ff165b15614aa7576040516316a163b960e11b815260040160405180910390fd5b6001600160a01b039283165f90815260209490945260408085208054600160ff19918216811790925593851686528186208054851682179055919093168452919092208054909216179055565b806001600160a01b03163b5f03614b2957604051634c9c8ce360e01b81526001600160a01b03821660048201526024016138bf565b5f80516020615eb183398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b60605f614b648484614c28565b9050808015614b8557505f3d1180614b8557505f846001600160a01b03163b115b15614b9a57614b92614c3b565b91505061392b565b8015614bc457604051639996b31560e01b81526001600160a01b03851660048201526024016138bf565b3d15614bd757614bd2614c54565b6146bc565b60405163d6bda27560e01b815260040160405180910390fd5b341561210c5760405163b398979f60e01b815260040160405180910390fd5b5f614c1861441f565b54600160401b900460ff16919050565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b508054614c6b906157d6565b5f825580601f10614c7a575050565b601f0160209004905f5260205f20908101906138d19190614d20565b6040518061014001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f81526020015f8152602001614d07604051806040016040528060608152602001606081525090565b815260606020820181905260408201819052015f905290565b5b80821115614d34575f8155600101614d21565b5090565b5f60208284031215614d48575f80fd5b5035919050565b6001600160a01b03811681146138d1575f80fd5b8035614d6e81614d4f565b919050565b5f60208284031215614d83575f80fd5b81356139e281614d4f565b5f8060408385031215614d9f575f80fd5b8235614daa81614d4f565b91506020830135614dba81614d4f565b809150509250929050565b634e487b7160e01b5f52602160045260245ffd5b60098110614df557634e487b7160e01b5f52602160045260245ffd5b9052565b602080825282518282018190525f919060409081850190868401855b82811015614e7357815180516001600160a01b039081168652878201518116888701528682015116868601526060808201519086015260809081015190614e5e81870183614dd9565b505060a0939093019290850190600101614e15565b5091979650505050505050565b5f815180845260208085019450602084015f5b83811015614eb85781516001600160a01b031687529582019590820190600101614e93565b509495945050505050565b60a081525f614ed560a0830188614e80565b8281036020840152614ee78188614e80565b90508281036040840152614efb8187614e80565b6001600160a01b0395861660608501529390941660809092019190915250949350505050565b604081525f614f336040830185614e80565b90508260208301529392505050565b600981106138d1575f80fd5b5f60208284031215614f5e575f80fd5b81356139e281614f42565b5f8083601f840112614f79575f80fd5b5081356001600160401b03811115614f8f575f80fd5b6020830191508360208260051b8501011115614fa9575f80fd5b9250929050565b5f805f805f805f6080888a031215614fc6575f80fd5b87356001600160401b0380821115614fdc575f80fd5b614fe88b838c01614f69565b909950975060208a0135915080821115615000575f80fd5b61500c8b838c01614f69565b909750955060408a0135915080821115615024575f80fd5b506150318a828b01614f69565b989b979a50959894979596606090950135949350505050565b602081525f6139e26020830184614e80565b6020810161392b8284614dd9565b602080825282518282018190525f919060409081850190868401855b82811015614e7357815180516001600160a01b039081168652878201518116888701528682015116868601526060908101519085015260809093019290850190600101615086565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f6139e260208301846150ce565b634e487b7160e01b5f52604160045260245ffd5b604080519081016001600160401b03811182821017156151445761514461510e565b60405290565b60405161014081016001600160401b03811182821017156151445761514461510e565b6040516101e081016001600160401b03811182821017156151445761514461510e565b604051601f8201601f191681016001600160401b03811182821017156151b8576151b861510e565b604052919050565b5f6001600160401b038211156151d8576151d861510e565b50601f01601f191660200190565b5f82601f8301126151f5575f80fd5b8135615208615203826151c0565b615190565b81815284602083860101111561521c575f80fd5b816020850160208301375f918101602001919091529392505050565b5f8060408385031215615249575f80fd5b823561525481614d4f565b915060208301356001600160401b0381111561526e575f80fd5b61527a858286016151e6565b9150509250929050565b5f81516040845261529860408501826150ce565b9050602083015184820360208601526152b182826150ce565b95945050505050565b80516001600160a01b031682525f61014060208301516152e560208601826001600160a01b03169052565b50604083015161530060408601826001600160a01b03169052565b50606083015161531b60608601826001600160a01b03169052565b506080830151608085015260a083015160a085015260c08301518160c086015261534782860182615284565b91505060e083015184820360e086015261536182826150ce565b915050610100808401518583038287015261537c83826150ce565b925050506101208084015161539382870182614dd9565b5090949350505050565b602081525f6139e260208301846152ba565b604081525f6153c16040830185614e80565b6020838203818501528185518084528284019150828160051b8501018388015f5b8381101561541057601f198784030185526153fe838351615284565b948601949250908501906001016153e2565b50909998505050505050505050565b604080825283519082018190525f906020906060840190828701845b8281101561545a57815160ff168452928401929084019060010161543b565b505050838103602085015261546f8186614e80565b9695505050505050565b5f806020838503121561548a575f80fd5b82356001600160401b0381111561549f575f80fd5b6154ab85828601614f69565b90969095509350505050565b5f60208083016020845280855180835260408601915060408160051b8701019250602087015f5b8281101561550c57603f198886030184526154fa8583516152ba565b945092850192908501906001016154de565b5092979650505050505050565b5f805f6040848603121561552b575f80fd5b833561553681614d4f565b925060208401356001600160401b0380821115615551575f80fd5b818601915086601f830112615564575f80fd5b813581811115615572575f80fd5b876020828501011115615583575f80fd5b6020830194508093505050509250925092565b5f604082840312156155a6575f80fd5b6155ae615122565b905081356001600160401b03808211156155c6575f80fd5b6155d2858386016151e6565b835260208401359150808211156155e7575f80fd5b506155f4848285016151e6565b60208301525092915050565b5f805f805f805f60e0888a031215615616575f80fd5b873561562181614d4f565b9650602088013561563181614d4f565b955061563f60408901614d63565b945061564d60608901614d63565b935060808801356001600160401b0380821115615668575f80fd5b6156748b838c01615596565b945060a08a0135915080821115615689575f80fd5b6156958b838c016151e6565b935060c08a01359150808211156156aa575f80fd5b506156b78a828b016151e6565b91505092959891949750929550565b634e487b7160e01b5f52601160045260245ffd5b5f600182016156eb576156eb6156c6565b5060010190565b634e487b7160e01b5f52603260045260245ffd5b604080825281018490525f8560608301825b8781101561574857823561572b81614d4f565b6001600160a01b0316825260209283019290910190600101615718565b508381036020858101919091528582529150859082015f5b8681101561578e57823561577381614f42565b61577d8382614dd9565b509183019190830190600101615760565b5098975050505050505050565b5f816157a9576157a96156c6565b505f190190565b8051614d6e81614d4f565b5f602082840312156157cb575f80fd5b81516139e281614d4f565b600181811c908216806157ea57607f821691505b60208210810361580857634e487b7160e01b5f52602260045260245ffd5b50919050565b808202811582820484141761392b5761392b6156c6565b8082018082111561392b5761392b6156c6565b5f6001600160401b038211156158505761585061510e565b5060051b60200190565b5f82601f830112615869575f80fd5b8151602061587961520383615838565b8083825260208201915060208460051b87010193508684111561589a575f80fd5b602086015b848110156158bf5780516158b281614d4f565b835291830191830161589f565b509695505050505050565b5f82601f8301126158d9575f80fd5b81516158e7615203826151c0565b8181528460208386010111156158fb575f80fd5b8160208501602083015e5f918101602001919091529392505050565b5f60408284031215615927575f80fd5b61592f615122565b905081516001600160401b0380821115615947575f80fd5b615953858386016158ca565b83526020840151915080821115615968575f80fd5b506155f4848285016158ca565b8051614d6e81614f42565b5f82601f83011261598f575f80fd5b8151602061599f61520383615838565b82815260059290921b840181019181810190868411156159bd575f80fd5b8286015b848110156158bf5780516001600160401b03808211156159df575f80fd5b90880190610140828b03601f19018113156159f8575f80fd5b615a0061514a565b615a0b8885016157b0565b81526040615a1a8186016157b0565b898301526060615a2b8187016157b0565b8284015260809150615a3e8287016157b0565b818401525060a0808601518284015260c0915081860151818401525060e08086015185811115615a6c575f80fd5b615a7a8f8c838a0101615917565b838501525061010091508186015185811115615a94575f80fd5b615aa28f8c838a01016158ca565b8285015250506101208086015185811115615abb575f80fd5b615ac98f8c838a01016158ca565b8385015250615ad9848701615975565b908301525086525050509183019183016159c1565b5f60208284031215615afe575f80fd5b81516001600160401b0380821115615b14575f80fd5b908301906101e08286031215615b28575f80fd5b615b3061516d565b615b39836157b0565b8152615b47602084016157b0565b6020820152615b58604084016157b0565b6040820152606083015160608201526080830151608082015260a083015160a082015260c083015160c082015260e083015160e0820152610100808401518183015250610120808401518183015250610140615bb58185016157b0565b90820152610160615bc78482016157b0565b90820152610180615bd98482016157b0565b908201526101a08381015183811115615bf0575f80fd5b615bfc8882870161585a565b8284015250506101c08084015183811115615c15575f80fd5b615c2188828701615980565b918301919091525095945050505050565b601f8211156141ec57805f5260205f20601f840160051c81016020851015615c575750805b601f840160051c820191505b81811015612b27575f8155600101615c63565b81516001600160401b03811115615c8f57615c8f61510e565b615ca381615c9d84546157d6565b84615c32565b602080601f831160018114615cd6575f8415615cbf5750858301515b5f19600386901b1c1916600185901b178555614552565b5f85815260208120601f198616915b82811015615d0457888601518255948401946001909101908401615ce5565b5085821015615d2157878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160401b03831115615d4857615d4861510e565b615d5c83615d5683546157d6565b83615c32565b5f601f841160018114615d8d575f8515615d765750838201355b5f19600387901b1c1916600186901b178355612b27565b5f83815260208120601f198716915b82811015615dbc5786850135825560209485019460019092019101615d9c565b5086821015615dd8575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b60208152816020820152818360408301375f818301604090810191909152601f909201601f19160101919050565b634e487b7160e01b5f52601260045260245ffd5b5f82615e3a57615e3a615e18565b500490565b5f82615e4d57615e4d615e18565b500690565b5f60208284031215615e62575f80fd5b5051919050565b8181038181111561392b5761392b6156c6565b634e487b7160e01b5f52603160045260245ffdfe34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbcd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4ba164736f6c6343000819000a",
}

// AddressBookV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use AddressBookV2MetaData.ABI instead.
var AddressBookV2ABI = AddressBookV2MetaData.ABI

// AddressBookV2BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressBookV2BinRuntime = `60806040526004361061045c575f3560e01c80637df40c621161023e578063b9f96f4011610138578063d9abb38b116100b5578063e70c38f111610079578063e70c38f114610d8b578063e8868e9f14610d9f578063f0a92ba814610db4578063f2fde38b14610dc8578063ffa1ad7414610de75761045c565b8063d9abb38b14610cf0578063da38d49814610d0f578063e1789d8114610d2e578063e4f0d37c14610d4d578063e59d7a8414610d6c5761045c565b8063cb1c2b5c116100fc578063cb1c2b5c14610c6c578063cf8c6f5214610c8a578063d18c07ab14610c9e578063d267eda514610cbd578063d3b5490714610cdc5761045c565b8063b9f96f4014610bd5578063ba70d01814610bf4578063be535f8b14610c13578063c732e08514610c32578063c9a86af214610c4d5761045c565b80639d8cf08f116101c6578063ad3cb1cc1161018a578063ad3cb1cc14610b35578063b42652e914610b65578063b57873a514610b84578063b756393014610ba3578063b858dd9514610bb65761045c565b80639d8cf08f14610a9a5780639f9e3cba14610ab9578063a41b600014610ad8578063a4c98ada14610af7578063a9ee547214610b165761045c565b80638da5cb5b1161020d5780638da5cb5b14610a2b5780638fabf38914610a3f5780639b7ae5ec14610a535780639d0e234d14610a675780639d0f5ef114610a865761045c565b80637df40c62146109b25780638129fc1c146109d157806387b7b8fd146109e55780638beeb439146109ff5761045c565b8063453e962e1161035a5780635b27b6c9116102d7578063715b208b1161029b578063715b208b1461091f578063766718081461094157806376a67a511461095557806378b84a5c14610974578063793c1946146109935761045c565b80635b27b6c91461088c578063656f5869146108ab5780636968b53f146108ca5780636abd623d146108ec578063715018a61461090b5761045c565b806350a5bb691161031e57806350a5bb69146107dc57806350de2fb3146107fa57806352d1902d14610819578063567b0b6c1461082d578063582115fb146108605761045c565b8063453e962e1461071f578063468e3a7e1461073e5780634a8c1fb41461076d5780634b6a94cc146107865780634f1ef286146107c95761045c565b80631b1a478b116103e857806325cf0943116103ac57806325cf094314610698578063291937f5146106ac5780632aca5091146106c05780632d4ede93146106e1578063394f8899146107005761045c565b80631b1a478b146105e15780631b8f34ca146106005780631ba3fd581461061f57806321d2320014610640578063229bb8231461066c5761045c565b80630a4ff2391161042f5780630a4ff2391461050e5780630b1fe7841461053057806315575d5a14610551578063160370b81461059a5780631865c57d146105bf5761045c565b806303e6689d14610481578063058529fb146104af57806306bb8471146104ce57806307ecec3e146104ef575b348015610467575f80fd5b50604051632053d6b560e11b815260040160405180910390fd5b34801561048c575f80fd5b50610495610dfb565b604080519283526020830191909152015b60405180910390f35b3480156104ba575f80fd5b506104956104c9366004614d38565b610e1b565b3480156104d9575f80fd5b506104ed6104e8366004614d73565b610e38565b005b3480156104fa575f80fd5b506104ed610509366004614d8e565b610f41565b348015610519575f80fd5b50610522610fde565b6040519081526020016104a6565b34801561053b575f80fd5b50610544610ff7565b6040516104a69190614df9565b34801561055c575f80fd5b5061057061056b366004614d73565b611141565b604080516001600160a01b03948516815292841660208401529216918101919091526060016104a6565b3480156105a5575f80fd5b506105ae611185565b6040516104a6959493929190614ec3565b3480156105ca575f80fd5b506105d36111af565b6040516104a6929190614f21565b3480156105ec575f80fd5b506105226105fb366004614f4e565b611210565b34801561060b575f80fd5b506104ed61061a366004614fb0565b611255565b34801561062a575f80fd5b506106336113a2565b6040516104a6919061504a565b34801561064b575f80fd5b506106546113b7565b6040516001600160a01b0390911681526020016104a6565b348015610677575f80fd5b5061068b610686366004614d73565b6113d2565b6040516104a6919061505c565b3480156106a3575f80fd5b506105706113ff565b3480156106b7575f80fd5b50610522611434565b3480156106cb575f80fd5b506106d4611446565b6040516104a6919061506a565b3480156106ec575f80fd5b506104ed6106fb366004614d73565b61155c565b34801561070b575f80fd5b506104ed61071a366004614d8e565b6116f2565b34801561072a575f80fd5b506104ed610739366004614d73565b611893565b348015610749575f80fd5b5061075d610758366004614d73565b6119c4565b60405190151581526020016104a6565b348015610778575f80fd5b50600c5461075d9060ff1681565b348015610791575f80fd5b506107bc6040518060400160405280600b81526020016a41646472657373426f6f6b60a81b81525081565b6040516104a691906150fc565b6104ed6107d7366004615238565b6119f1565b3480156107e7575f80fd5b50600c5461075d90610100900460ff1681565b348015610805575f80fd5b506104ed610814366004614d73565b611a10565b348015610824575f80fd5b50610522611a59565b348015610838575f80fd5b506105227f000000000000000000000000000000000000000000000000000000000000000081565b34801561086b575f80fd5b5061087f61087a366004614d73565b611a74565b6040516104a6919061539d565b348015610897575f80fd5b506104ed6108a6366004614d38565b611db8565b3480156108b6575f80fd5b506104ed6108c5366004614d73565b611def565b3480156108d5575f80fd5b506108de611e5b565b6040516104a69291906153af565b3480156108f7575f80fd5b50600754610654906001600160a01b031681565b348015610916575f80fd5b506104ed6120fb565b34801561092a575f80fd5b5061093361210e565b6040516104a692919061541f565b34801561094c575f80fd5b5061052261245a565b348015610960575f80fd5b506104ed61096f366004614d73565b612463565b34801561097f575f80fd5b506104ed61098e366004614d73565b61252f565b34801561099e575f80fd5b506104ed6109ad366004614d73565b6125a3565b3480156109bd575f80fd5b506104ed6109cc366004614d73565b6125ea565b3480156109dc575f80fd5b506104ed61260d565b3480156109f0575f80fd5b506106546002600160a01b0381565b348015610a0a575f80fd5b50610a1e610a19366004615479565b612b2e565b6040516104a691906154b7565b348015610a36575f80fd5b50610654612ee3565b348015610a4a575f80fd5b50610522612f11565b348015610a5e575f80fd5b50610654612f23565b348015610a72575f80fd5b506104ed610a81366004614d38565b612f3e565b348015610a91575f80fd5b50610495612f60565b348015610aa5575f80fd5b506104ed610ab4366004614d73565b612f7e565b348015610ac4575f80fd5b506104ed610ad3366004614d8e565b612fa0565b348015610ae3575f80fd5b506104ed610af2366004614d73565b613128565b348015610b02575f80fd5b506104ed610b11366004614d38565b61319c565b348015610b21575f80fd5b506104ed610b30366004614d38565b6131bf565b348015610b40575f80fd5b506107bc604051806040016040528060058152602001640352e302e360dc1b81525081565b348015610b70575f80fd5b506104ed610b7f366004614d73565b6131e2565b348015610b8f575f80fd5b506104ed610b9e366004614d73565b6132eb565b348015610bae575f80fd5b506001610522565b348015610bc1575f80fd5b50600654610654906001600160a01b031681565b348015610be0575f80fd5b506104ed610bef366004614d73565b61330e565b348015610bff575f80fd5b506104ed610c0e366004614d38565b61334c565b348015610c1e575f80fd5b506104ed610c2d366004614d73565b61336f565b348015610c3d575f80fd5b50610522678ac7230489e8000081565b348015610c58575f80fd5b506104ed610c67366004614d73565b6133de565b348015610c77575f80fd5b506105226a0422ca8b0a00a42500000081565b348015610c95575f80fd5b50610633613401565b348015610ca9575f80fd5b506104ed610cb8366004614d38565b613416565b348015610cc8575f80fd5b50600554610654906001600160a01b031681565b348015610ce7575f80fd5b50610522613439565b348015610cfb575f80fd5b506104ed610d0a366004614d73565b61344b565b348015610d1a575f80fd5b506104ed610d29366004615519565b61348c565b348015610d39575f80fd5b506104ed610d48366004615600565b613534565b348015610d58575f80fd5b506104ed610d67366004614d73565b6137c8565b348015610d77575f80fd5b506104ed610d86366004614d38565b61383d565b348015610d96575f80fd5b50610495613860565b348015610daa575f80fd5b5061052261080081565b348015610dbf575f80fd5b50610522613880565b348015610dd3575f80fd5b506104ed610de2366004614d73565b613892565b348015610df2575f80fd5b50610522600281565b5f805f610e066138d4565b905080600e015481600f015492509250509091565b5f80610e26836138f8565b610e2f84613931565b91509150915091565b610e40613957565b5f610e496138d4565b90505f6001600160a01b0383165f908152602083905260409020600a015460ff166008811115610e7b57610e7b614dc5565b03610e9957604051634825e09360e01b815260040160405180910390fd5b6001600160a01b0382165f9081526020829052604090206005015415610ed257604051637be80ce960e11b815260040160405180910390fd5b5f816007015f8154610ee3906156da565b91829055506001600160a01b0384165f8181526020858152604091829020600501849055905183815292935090917fe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb910160405180910390a2505050565b81610f4b8161398b565b6001600160a01b038216610f725760405163b4fa3fb360e01b815260040160405180910390fd5b5f610f7b6138d4565b6001600160a01b038086165f8181526020849052604080822080548986166001600160a01b03198216811790925591519596509316938492917f8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda8791a45050505050565b5f610ff2610fea6138d4565b6001016139ce565b905090565b60605f6110026138d4565b90505f611011826001016139ce565b9050806001600160401b0381111561102b5761102b61510e565b60405190808252806020026020018201604052801561108a57816020015b6110776040805160a0810182525f808252602082018190529181018290526060810182905290608082015290565b8152602001906001900390816110495790505b5092505f5b8181101561113b575f6110a560018501836139d7565b6001600160a01b038082165f8181526020888152604091829020825160a081018452938452600181015485169184019190915260028101549093169082015260048201546060820152600a8201549293509091608082019060ff16600881111561111157611111614dc5565b815250868481518110611126576111266156f2565b6020908102919091010152505060010161108f565b50505090565b5f805f805f80611150876139e9565b92509250925080611174576040516342dc2dc560e01b815260040160405180910390fd5b5085945090925090505b9193909250565b60608060605f805f805f805f611199613a61565b939e929d50909b50995090975095505050505050565b6040805160018082528183019092526060915f918291602080830190803683370190505090506111dd612ee3565b815f815181106111ef576111ef6156f2565b6001600160a01b039092166020928302919091019091015292600192509050565b5f6112196138d4565b6008015f83600881111561122f5761122f614dc5565b600881111561124057611240614dc5565b81526020019081526020015f20549050919050565b61125d613c45565b85848114158061126d5750808314155b1561128b5760405163b4fa3fb360e01b815260040160405180910390fd5b5f5b8181101561130c576113048989838181106112aa576112aa6156f2565b90506020020160208101906112bf9190614d73565b8888848181106112d1576112d16156f2565b90506020020160208101906112e69190614f4e565b8787858181106112f8576112f86156f2565b90506020020135613c6c565b60010161128d565b50611315613ecf565b1561135b57816113236138d4565b601001556040518281527fd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e848160529060200160405180910390a15b7fab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e888888886040516113909493929190615706565b60405180910390a15050505050505050565b6060610ff26113af6138d4565b600501613f00565b5f6113c06138d4565b601501546001600160a01b0316919050565b5f6113db6138d4565b6001600160a01b039092165f9081526020929092525060409020600a015460ff1690565b5f805f8061140b6138d4565b601281015460138201546014909201546001600160a01b03918216979282169650169350915050565b5f61143d6138d4565b600a0154905090565b60605f6114516138d4565b90505f611460826001016139ce565b9050806001600160401b0381111561147a5761147a61510e565b6040519080825280602002602001820160405280156114ca57816020015b604080516080810182525f8082526020808301829052928201819052606082015282525f199092019101816114985790505b5092505f5b8181101561113b575f6114e560018501836139d7565b6001600160a01b038082165f8181526020888152604091829020825160808101845293845260018101548516918401919091526003810154909316908201526005820154606082015287519293509091879085908110611547576115476156f2565b602090810291909101015250506001016114cf565b806115668161398b565b611571826001613f0c565b61158e5760405163baf3f0f760e01b815260040160405180910390fd5b5f6115976138d4565b6001600160a01b038481165f908152602083815260408083206001808201546002830154600989018652848720805460ff1990811690915591881687528487208054831690559096168552828520805490961690955593835260088501909152812080549394509192919061160b8361579b565b9091555061161e90506003830185613f41565b506001600160a01b0384165f90815260208390526040812080546001600160a01b031990811682556001820180548216905560028201805482169055600382018054909116905560048101829055600581018290559060068201816116838282614c5f565b611690600183015f614c5f565b506116a09050600883015f614c5f565b6116ad600983015f614c5f565b50600a01805460ff191690556040516001600160a01b038516907f1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad905f90a250505050565b816116fc8161398b565b6001600160a01b0382166117235760405163b4fa3fb360e01b815260040160405180910390fd5b5f61172c6138d4565b6001600160a01b038086165f908152602083815260408083206001810154825163e1a12d3560e01b8152925196975090959394169263e1a12d35926004808401939192918290030181865afa158015611787573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906117ab91906157bb565b6001600160a01b0316146117d257604051638ed87ef960e01b815260040160405180910390fd5b6001600160a01b0384165f90815260098301602052604090205460ff161561180d576040516316a163b960e11b815260040160405180910390fd5b6002810180546001600160a01b039081165f818152600986016020526040808220805460ff19908116909155898516808452828420805490921660011790915585546001600160a01b0319168117909555519193928492908a16917f270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc6970191a4505050505050565b8061189d81613f55565b6118a8826001613f0c565b6118c55760405163baf3f0f760e01b815260040160405180910390fd5b6118ce82613f7e565b6118eb5760405163bf74735560e01b815260040160405180910390fd5b678ac7230489e80000826001600160a01b031631101561191d5760405162b8ec7b60e61b815260040160405180910390fd5b5f6119266138d4565b905080600f01546119376002611210565b106119555760405163848084dd60e01b815260040160405180910390fd5b80600e0154611962610fde565b106119805760405163848084dd60e01b815260040160405180910390fd5b61198c8360025f613c6c565b6040516001600160a01b038416907fb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98905f90a2505050565b5f6119cd6138d4565b6001600160a01b039092165f90815260099290920160205250604090205460ff1690565b6119f9614084565b611a0282614128565b611a0c8282614130565b5050565b611a186141f1565b60035f80516020615ed1833981519152611a33601584614223565b604080516001600160a01b0392831681529185166020830152015b60405180910390a250565b5f611a62614283565b505f80516020615eb183398151915290565b611a7c614c96565b5f611a856138d4565b6001600160a01b0384165f908152602091909152604081209150600a82015460ff166008811115611ab857611ab8614dc5565b03611ad657604051634825e09360e01b815260040160405180910390fd5b604080516101408101825282546001600160a01b0390811682526001840154811660208301526002840154811682840152600384015416606082015260048301546080820152600583015460a082015281518083019092526006830180549192849260c08501929082908290611b4b906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611b77906157d6565b8015611bc25780601f10611b9957610100808354040283529160200191611bc2565b820191905f5260205f20905b815481529060010190602001808311611ba557829003601f168201915b50505050508152602001600182018054611bdb906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611c07906157d6565b8015611c525780601f10611c2957610100808354040283529160200191611c52565b820191905f5260205f20905b815481529060010190602001808311611c3557829003601f168201915b5050505050815250508152602001600882018054611c6f906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611c9b906157d6565b8015611ce65780601f10611cbd57610100808354040283529160200191611ce6565b820191905f5260205f20905b815481529060010190602001808311611cc957829003601f168201915b50505050508152602001600982018054611cff906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611d2b906157d6565b8015611d765780601f10611d4d57610100808354040283529160200191611d76565b820191905f5260205f20905b815481529060010190602001808311611d5957829003601f168201915b5050509183525050600a82015460209091019060ff166008811115611d9d57611d9d614dc5565b6008811115611dae57611dae614dc5565b9052509392505050565b611dc0613957565b60055f80516020615e91833981519152611ddb6011846142cc565b604080519182526020820185905201611a4e565b80611df981613f55565b611e04826004613f0c565b611e215760405163baf3f0f760e01b815260040160405180910390fd5b611e2a82613f7e565b611e475760405163bf74735560e01b815260040160405180910390fd5b611a0c826005611e56856142ed565b613c6c565b6060805f611e676138d4565b90505f611e76826001016139ce565b9050806001600160401b03811115611e9057611e9061510e565b604051908082528060200260200182016040528015611eb9578160200160208202803683370190505b509350806001600160401b03811115611ed457611ed461510e565b604051908082528060200260200182016040528015611f1957816020015b6040805180820190915260608082526020820152815260200190600190039081611ef25790505b5092505f5b818110156120f457611f3360018401826139d7565b858281518110611f4557611f456156f2565b60200260200101906001600160a01b031690816001600160a01b031681525050825f015f868381518110611f7b57611f7b6156f2565b60200260200101516001600160a01b03166001600160a01b031681526020019081526020015f206006016040518060400160405290815f82018054611fbf906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054611feb906157d6565b80156120365780601f1061200d57610100808354040283529160200191612036565b820191905f5260205f20905b81548152906001019060200180831161201957829003601f168201915b5050505050815260200160018201805461204f906157d6565b80601f016020809104026020016040519081016040528092919081815260200182805461207b906157d6565b80156120c65780601f1061209d576101008083540402835291602001916120c6565b820191905f5260205f20905b8154815290600101906020018083116120a957829003601f168201915b5050505050815250508482815181106120e1576120e16156f2565b6020908102919091010152600101611f1e565b5050509091565b6121036141f1565b61210c5f614317565b565b600c54606090819060ff16612137575050604080515f8082526020820190815281830190925291565b5f805f805f612144613a61565b84519499509297509095509350915061215e81600361580e565b612169906002615825565b6001600160401b038111156121805761218061510e565b6040519080825280602002602001820160405280156121a9578160200160208202803683370190505b5097506121b781600361580e565b6121c2906002615825565b6001600160401b038111156121d9576121d961510e565b604051908082528060200260200182016040528015612202578160200160208202803683370190505b5096505f805b8281101561238f575f8a8381518110612223576122236156f2565b602002602001019060ff16908160ff1681525050878181518110612249576122496156f2565b602002602001015189838061225d906156da565b94508151811061226f5761226f6156f2565b60200260200101906001600160a01b031690816001600160a01b03168152505060018a83815181106122a3576122a36156f2565b602002602001019060ff16908160ff16815250508681815181106122c9576122c96156f2565b60200260200101518983806122dd906156da565b9450815181106122ef576122ef6156f2565b60200260200101906001600160a01b031690816001600160a01b03168152505060028a8381518110612323576123236156f2565b602002602001019060ff16908160ff1681525050858181518110612349576123496156f2565b602002602001015189838061235d906156da565b94508151811061236f5761236f6156f2565b6001600160a01b0390921660209283029190910190910152600101612208565b5060038982815181106123a4576123a46156f2565b60ff909216602092830291909101909101528388826123c2816156da565b9350815181106123d4576123d46156f2565b60200260200101906001600160a01b031690816001600160a01b0316815250506004898281518110612408576124086156f2565b602002602001019060ff16908160ff16815250508288828151811061242f5761242f6156f2565b60200260200101906001600160a01b031690816001600160a01b031681525050505050505050509091565b5f610ff2614387565b8061246d81613f55565b612478826006613f0c565b6124955760405163baf3f0f760e01b815260040160405180910390fd5b5f61249e6138d4565b90506124ad81601001546138f8565b6124b76007611210565b106124d55760405163848084dd60e01b815260040160405180910390fd5b6124e28160100154613931565b6124ec6006611210565b1161250a5760405163848084dd60e01b815260040160405180910390fd5b5f81600c01544261251b9190615825565b905061252984600783613c6c565b50505050565b6125376143b2565b5f6125406138d4565b905061254f6005820183613f41565b61256c5760405163d33ff8c160e01b815260040160405180910390fd5b6040516001600160a01b038316907f814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e905f90a25050565b806125ad81613f55565b6125b8826007613f0c565b6125d55760405163baf3f0f760e01b815260040160405180910390fd5b6125de826143e6565b611a0c8260065f613c6c565b6125f2613957565b60015f80516020615ed1833981519152611a33601384614223565b5f61261661441f565b805490915060ff600160401b82041615906001600160401b03165f8115801561263c5750825b90505f826001600160401b031660011480156126575750303b155b905081158015612665575080155b156126835760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff1916600117855583156126ad57845460ff60401b1916600160401b1785555b60405163e2693e3f60e01b815260206004820152601060248201526f10509d8c91185d1850dbdb9d1c9858dd60821b60448201525f906104019063e2693e3f90606401602060405180830381865afa15801561270b573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061272f91906157bb565b90506001600160a01b0381166127585760405163aed5959560e01b815260040160405180910390fd5b5f816001600160a01b031663ebe58ed76040518163ffffffff1660e01b81526004015f60405180830381865afa158015612794573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526127bb9190810190615aee565b90505f6127c66138d4565b90506127d4825f0151614447565b6060820151600a8201556080820151600b82015560a0820151600c82015560c0820151600d82015560e0820151600e8201556101008201516011820155610120820151600f8201556101408201516012820180546001600160a01b03199081166001600160a01b03938416179091556101608401516013840180548316918416919091179055610180840151601484018054831691841691909117905560208401516015840180548316918416919091179055604084015160168401805490921692169190911790556101a0820151515f5b81811015612a7d575f846101a0015182815181106128c6576128c66156f2565b602002602001015190505f856101c0015183815181106128e8576128e86156f2565b6020908102919091018101516001600160a01b038481165f90815260098901845260408082208054600160ff199182168117909255958501518416835281832080548716821790558185015190931682529020805490931617909155905060066101208201819052505f608082018181526001600160a01b0380851683526020888152604093849020855181549084166001600160a01b03199182161782559186015160018201805491851691841691909117905593850151600285018054918416918316919091179055606085015160038501805491909316911617905551600482015560a0820151600582015560c0820151805183929190600683019081906129f39082615c76565b5060208201516001820190612a089082615c76565b50505060e08201516008820190612a1f9082615c76565b506101008201516009820190612a359082615c76565b50610120820151600a8201805460ff19166001836008811115612a5a57612a5a614dc5565b0217905550612a6f9150506001860183614458565b5050508060010190506128a6565b5060065f9081526008830160205260409081902082905560108301829055606460078401556101a084015190517f820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d95291612ad59161504a565b60405180910390a1505050508315612b2757845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b60605f612b396138d4565b905082806001600160401b03811115612b5457612b5461510e565b604051908082528060200260200182016040528015612b8d57816020015b612b7a614c96565b815260200190600190039081612b725790505b5092505f5b81811015612eda57825f878784818110612bae57612bae6156f2565b9050602002016020810190612bc39190614d73565b6001600160a01b03908116825260208083019390935260409182015f20825161014081018452815483168152600182015483169481019490945260028101548216848401526003810154909116606084015260048101546080840152600581015460a08401528151808301909252600681018054919260c085019290919082908290612c4e906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054612c7a906157d6565b8015612cc55780601f10612c9c57610100808354040283529160200191612cc5565b820191905f5260205f20905b815481529060010190602001808311612ca857829003601f168201915b50505050508152602001600182018054612cde906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054612d0a906157d6565b8015612d555780601f10612d2c57610100808354040283529160200191612d55565b820191905f5260205f20905b815481529060010190602001808311612d3857829003601f168201915b5050505050815250508152602001600882018054612d72906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054612d9e906157d6565b8015612de95780601f10612dc057610100808354040283529160200191612de9565b820191905f5260205f20905b815481529060010190602001808311612dcc57829003601f168201915b50505050508152602001600982018054612e02906157d6565b80601f0160208091040260200160405190810160405280929190818152602001828054612e2e906157d6565b8015612e795780601f10612e5057610100808354040283529160200191612e79565b820191905f5260205f20905b815481529060010190602001808311612e5c57829003601f168201915b5050509183525050600a82015460209091019060ff166008811115612ea057612ea0614dc5565b6008811115612eb157612eb1614dc5565b81525050848281518110612ec757612ec76156f2565b6020908102919091010152600101612b92565b50505092915050565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b5f612f1a6138d4565b60110154905090565b5f612f2c6138d4565b601601546001600160a01b0316919050565b612f46613957565b5f5f80516020615e91833981519152611ddb600c846142cc565b5f80612f76612f6d6138d4565b60100154610e1b565b915091509091565b612f86613957565b5f5f80516020615ed1833981519152611a33601284614223565b81612faa8161398b565b5f612fb36138d4565b6001600160a01b038086165f9081526020839052604080822060030180548885166001600160a01b0319821617909155905163e2693e3f60e01b8152939450909116916104019063e2693e3f9061302f906004016020808252600e908201526d29ba30b5b4b733aa3930b1b5b2b960911b604082015260600190565b602060405180830381865afa15801561304a573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061306e91906157bb565b90506001600160a01b038116156130d65760405163aad8cb3f60e01b81526001600160a01b03878116600483015282169063aad8cb3f906024015f604051808303815f87803b1580156130bf575f80fd5b505af11580156130d1573d5f803e3d5ffd5b505050505b846001600160a01b0316826001600160a01b0316876001600160a01b03167f23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d9660405160405180910390a4505050505050565b6131306143b2565b5f6131396138d4565b90506131486005820183614458565b61316557604051633ad2b1bb60e11b815260040160405180910390fd5b6040516001600160a01b038316907fb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd905f90a25050565b6131a4613957565b60045f80516020615e91833981519152611ddb600e846142cc565b6131c7613957565b60065f80516020615e91833981519152611ddb600f846142cc565b806131ec81613f55565b5f6131f6836113d2565b9050600781600881111561320c5761320c614dc5565b0361321f5761321a836143e6565b613251565b600681600881111561323357613233614dc5565b146132515760405163baf3f0f760e01b815260040160405180910390fd5b5f61325a6138d4565b905061326981601001546138f8565b6132736008611210565b106132915760405163848084dd60e01b815260040160405180910390fd5b60068260088111156132a5576132a5614dc5565b036132df576132b78160100154613931565b6132c16006611210565b116132df5760405163848084dd60e01b815260040160405180910390fd5b6125298460085f613c6c565b6132f36141f1565b60045f80516020615ed1833981519152611a33601684614223565b8061331881613f55565b613323826004613f0c565b6133405760405163baf3f0f760e01b815260040160405180910390fd5b611a0c8260015f613c6c565b613354613957565b60025f80516020615e91833981519152611ddb600a846142cc565b613377613957565b5f6133806138d4565b6001600160a01b0383165f908152602082905260408120600501549192508190036133be576040516357024f6d60e11b815260040160405180910390fd5b506001600160a01b039091165f9081526020919091526040812060050155565b6133e6613957565b60025f80516020615ed1833981519152611a33601484614223565b6060610ff261340e6138d4565b600301613f00565b61341e613957565b60035f80516020615e91833981519152611ddb600b846142cc565b5f6134426138d4565b60100154905090565b8061345581613f55565b613460826005613f0c565b61347d5760405163baf3f0f760e01b815260040160405180910390fd5b611a0c826004611e56856142ed565b826134968161398b565b6108008211156134b95760405163b4fa3fb360e01b815260040160405180910390fd5b82826134c36138d4565b6001600160a01b0387165f90815260209190915260409020600901916134ea919083615d31565b50836001600160a01b03167f2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf8484604051613526929190615dea565b60405180910390a250505050565b5f61353e886113d2565b600881111561354f5761354f614dc5565b1461356d5760405163731918fb60e11b815260040160405180910390fd5b81515f0361358e5760405163b4fa3fb360e01b815260040160405180910390fd5b610800815111156135b25760405163b4fa3fb360e01b815260040160405180910390fd5b5f6135bb6138d4565b90506135cd816009018989898861446c565b5f604051806101400160405280336001600160a01b03168152602001896001600160a01b03168152602001886001600160a01b03168152602001876001600160a01b031681526020015f81526020015f81526020018681526020018581526020018481526020016001600881111561364757613647614dc5565b90526001600160a01b03808b165f9081526020858152604091829020845181549085166001600160a01b031991821617825591850151600182018054918616918416919091179055918401516002830180549185169183169190911790556060840151600383018054919094169116179091556080820151600482015560a0820151600582015560c082015180519293508392600683019081906136eb9082615c76565b50602082015160018201906137009082615c76565b50505060e082015160088201906137179082615c76565b50610100820151600982019061372d9082615c76565b50610120820151600a8201805460ff1916600183600881111561375257613752614dc5565b0217905550613767915050600383018a614458565b5060015f9081526008830160205260408120805491613785836156da565b90915550506040516001600160a01b038a16907f55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792905f90a2505050505050505050565b806137d281613f55565b6137dd826002613f0c565b6137fa5760405163baf3f0f760e01b815260040160405180910390fd5b6138068260015f613c6c565b6040516001600160a01b038316907f8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2905f90a25050565b613845613957565b60015f80516020615e91833981519152611ddb600d846142cc565b5f805f61386b6138d4565b905080600c015481600d015492509250509091565b5f6138896138d4565b600b0154905090565b61389a6141f1565b6001600160a01b0381166138c857604051631e4fbdf760e01b81525f60048201526024015b60405180910390fd5b6138d181614317565b50565b7f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff0090565b5f600482101561390957505f919050565b6002613916600384615e2c565b613921906001615825565b61392b9190615e2c565b92915050565b5f600482101561393f575090565b600361394c83600261580e565b613921906002615825565b61395f6138d4565b601601546001600160a01b0316331461210c5760405163033b71e160e41b815260040160405180910390fd5b6139936138d4565b6001600160a01b038281165f9081526020929092526040909120541633146138d15760405163605919ad60e11b815260040160405180910390fd5b5f61392b825490565b5f6139e2838361455a565b9392505050565b5f805f806139f56138d4565b6001600160a01b0386165f908152602082905260408120919250600a82015460ff166008811115613a2857613a28614dc5565b03613a3d575f805f945094509450505061117e565b6001818101546002909201546001600160a01b039283169892169650945092505050565b60608060605f805f613a716138d4565b90505f613a80826001016139ce565b9050806001600160401b03811115613a9a57613a9a61510e565b604051908082528060200260200182016040528015613ac3578160200160208202803683370190505b509650806001600160401b03811115613ade57613ade61510e565b604051908082528060200260200182016040528015613b07578160200160208202803683370190505b509550806001600160401b03811115613b2257613b2261510e565b604051908082528060200260200182016040528015613b4b578160200160208202803683370190505b5094505f5b81811015613c1d575f613b6660018501836139d7565b6001600160a01b0381165f9081526020869052604090208a519192509082908b9085908110613b9757613b976156f2565b6001600160a01b03928316602091820292909201015260018201548a519116908a9085908110613bc957613bc96156f2565b6001600160a01b03928316602091820292909201015260028201548951911690899085908110613bfb57613bfb6156f2565b6001600160a01b03909216602092830291909101909101525050600101613b50565b505060138101546012909101549596949593946001600160a01b039182169490911692509050565b336002600160a01b031461210c576040516354d325c360e01b815260040160405180910390fd5b5f613c756138d4565b6001600160a01b0385165f908152602082905260408120600a8101549293509160ff1690816008811115613cab57613cab614dc5565b1480613cc757505f856008811115613cc557613cc5614dc5565b145b80613cf35750846008811115613cdf57613cdf614dc5565b816008811115613cf157613cf1614dc5565b145b15613d0057505050505050565b826008015f826008811115613d1757613d17614dc5565b6008811115613d2857613d28614dc5565b81526020019081526020015f205f815480929190613d459061579b565b9190505550826008015f866008811115613d6157613d61614dc5565b6008811115613d7257613d72614dc5565b81526020019081526020015f205f815480929190613d8f906156da565b9091555060019050816008811115613da957613da9614dc5565b148015613dc857506001856008811115613dc557613dc5614dc5565b14155b15613dee57613dda6001840187614458565b50613de86003840187613f41565b50613e43565b6001816008811115613e0257613e02614dc5565b14158015613e2157506001856008811115613e1f57613e1f614dc5565b145b15613e4357613e336001840187613f41565b50613e416003840187614458565b505b600a8201805486919060ff19166001836008811115613e6457613e64614dc5565b021790555060048201849055846008811115613e8257613e82614dc5565b816008811115613e9457613e94614dc5565b6040516001600160a01b038916907fcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee905f90a4505050505050565b5f613efa7f000000000000000000000000000000000000000000000000000000000000000043615e3f565b15919050565b60605f6139e283614580565b5f816008811115613f1f57613f1f614dc5565b613f28846113d2565b6008811115613f3957613f39614dc5565b149392505050565b5f6139e2836001600160a01b0384166145d9565b336001600160a01b038216146138d1576040516335f1334d60e11b815260040160405180910390fd5b5f80613f886138d4565b6001600160a01b038085165f9081526020928352604080822060010154815163318588a360e11b81529151931694509092849263630b11469260048082019392918290030181865afa158015613fe0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906140049190615e52565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015614040573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906140649190615e52565b61406e9190615e69565b6a0422ca8b0a00a4250000001115949350505050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016148061410a57507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166140fe5f80516020615eb1833981519152546001600160a01b031690565b6001600160a01b031614155b1561210c5760405163703e46dd60e11b815260040160405180910390fd5b6138d16141f1565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561418a575060408051601f3d908101601f1916820190925261418791810190615e52565b60015b6141b257604051634c9c8ce360e01b81526001600160a01b03831660048201526024016138bf565b5f80516020615eb183398151915281146141e257604051632a87526960e21b8152600481018290526024016138bf565b6141ec83836146c3565b505050565b336141fa612ee3565b6001600160a01b03161461210c5760405163118cdaa760e01b81523360048201526024016138bf565b5f6001600160a01b03821661424b5760405163b4fa3fb360e01b815260040160405180910390fd5b5f614276847f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff00615825565b8054939055509092915050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461210c5760405163703e46dd60e11b815260040160405180910390fd5b5f815f0361424b5760405163b4fa3fb360e01b815260040160405180910390fd5b5f6142f66138d4565b6001600160a01b039092165f90815260209290925250604090206004015490565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f610ff27f000000000000000000000000000000000000000000000000000000000000000043615e2c565b6143ba6138d4565b601501546001600160a01b0316331461210c5760405163333f4e6560e01b815260040160405180910390fd5b5f6143f0826142ed565b905080158015906144015750804210155b15611a0c5760405163b48d5fc760e01b815260040160405180910390fd5b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a0061392b565b61444f614718565b6138d18161473d565b5f6139e2836001600160a01b038416614745565b61447884848484614791565b5f6144816148fd565b90506001600160a01b0381166144aa5760405163cdded31d60e01b815260040160405180910390fd5b60405163669d8d4560e01b81526001600160a01b03858116600483015233919083169063669d8d4590602401602060405180830381865afa1580156144f1573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061451591906157bb565b6001600160a01b03161461453c57604051632281776f60e01b815260040160405180910390fd5b614546848461497f565b61455286868686614a2a565b505050505050565b5f825f01828154811061456f5761456f6156f2565b905f5260205f200154905092915050565b6060815f018054806020026020016040519081016040528092919081815260200182805480156145cd57602002820191905f5260205f20905b8154815260200190600101908083116145b9575b50505050509050919050565b5f81815260018301602052604081205480156146b3575f6145fb600183615e69565b85549091505f9061460e90600190615e69565b905080821461466d575f865f01828154811061462c5761462c6156f2565b905f5260205f200154905080875f01848154811061464c5761464c6156f2565b5f918252602080832090910192909255918252600188019052604090208390555b855486908061467e5761467e615e7c565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f90556001935050505061392b565b5f91505061392b565b5092915050565b6146cc82614af4565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115614710576141ec8282614b57565b611a0c614bf0565b614720614c0f565b61210c57604051631afcd79f60e31b815260040160405180910390fd5b61389a614718565b5f81815260018301602052604081205461478a57508154600181810184555f84815260208082209093018490558454848252828601909352604090209190915561392b565b505f61392b565b6001600160a01b03841615806147ae57506001600160a01b038316155b806147c057506001600160a01b038216155b156147de5760405163b4fa3fb360e01b815260040160405180910390fd5b826001600160a01b0316846001600160a01b0316148061480f5750816001600160a01b0316846001600160a01b0316145b8061482b5750816001600160a01b0316836001600160a01b0316145b156148495760405163b4fa3fb360e01b815260040160405180910390fd5b80515160301415806148615750806020015151606014155b1561487f5760405163b4fa3fb360e01b815260040160405180910390fd5b805180516020909101207fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29314806148df575060208082015180519101207f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21145b156125295760405163b4fa3fb360e01b815260040160405180910390fd5b60405163e2693e3f60e01b815260206004820152601060248201526f436e5374616b696e67466163746f727960801b60448201525f906104019063e2693e3f90606401602060405180830381865afa15801561495b573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ff291906157bb565b5f826001600160a01b031663e1a12d356040518163ffffffff1660e01b8152600401602060405180830381865afa1580156149bc573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906149e091906157bb565b90506001600160a01b03811615801590614a0c5750816001600160a01b0316816001600160a01b031614155b156141ec5760405163b4fa3fb360e01b815260040160405180910390fd5b6001600160a01b0383165f9081526020859052604090205460ff1680614a6757506001600160a01b0382165f9081526020859052604090205460ff165b80614a8957506001600160a01b0381165f9081526020859052604090205460ff165b15614aa7576040516316a163b960e11b815260040160405180910390fd5b6001600160a01b039283165f90815260209490945260408085208054600160ff19918216811790925593851686528186208054851682179055919093168452919092208054909216179055565b806001600160a01b03163b5f03614b2957604051634c9c8ce360e01b81526001600160a01b03821660048201526024016138bf565b5f80516020615eb183398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b60605f614b648484614c28565b9050808015614b8557505f3d1180614b8557505f846001600160a01b03163b115b15614b9a57614b92614c3b565b91505061392b565b8015614bc457604051639996b31560e01b81526001600160a01b03851660048201526024016138bf565b3d15614bd757614bd2614c54565b6146bc565b60405163d6bda27560e01b815260040160405180910390fd5b341561210c5760405163b398979f60e01b815260040160405180910390fd5b5f614c1861441f565b54600160401b900460ff16919050565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b508054614c6b906157d6565b5f825580601f10614c7a575050565b601f0160209004905f5260205f20908101906138d19190614d20565b6040518061014001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f81526020015f8152602001614d07604051806040016040528060608152602001606081525090565b815260606020820181905260408201819052015f905290565b5b80821115614d34575f8155600101614d21565b5090565b5f60208284031215614d48575f80fd5b5035919050565b6001600160a01b03811681146138d1575f80fd5b8035614d6e81614d4f565b919050565b5f60208284031215614d83575f80fd5b81356139e281614d4f565b5f8060408385031215614d9f575f80fd5b8235614daa81614d4f565b91506020830135614dba81614d4f565b809150509250929050565b634e487b7160e01b5f52602160045260245ffd5b60098110614df557634e487b7160e01b5f52602160045260245ffd5b9052565b602080825282518282018190525f919060409081850190868401855b82811015614e7357815180516001600160a01b039081168652878201518116888701528682015116868601526060808201519086015260809081015190614e5e81870183614dd9565b505060a0939093019290850190600101614e15565b5091979650505050505050565b5f815180845260208085019450602084015f5b83811015614eb85781516001600160a01b031687529582019590820190600101614e93565b509495945050505050565b60a081525f614ed560a0830188614e80565b8281036020840152614ee78188614e80565b90508281036040840152614efb8187614e80565b6001600160a01b0395861660608501529390941660809092019190915250949350505050565b604081525f614f336040830185614e80565b90508260208301529392505050565b600981106138d1575f80fd5b5f60208284031215614f5e575f80fd5b81356139e281614f42565b5f8083601f840112614f79575f80fd5b5081356001600160401b03811115614f8f575f80fd5b6020830191508360208260051b8501011115614fa9575f80fd5b9250929050565b5f805f805f805f6080888a031215614fc6575f80fd5b87356001600160401b0380821115614fdc575f80fd5b614fe88b838c01614f69565b909950975060208a0135915080821115615000575f80fd5b61500c8b838c01614f69565b909750955060408a0135915080821115615024575f80fd5b506150318a828b01614f69565b989b979a50959894979596606090950135949350505050565b602081525f6139e26020830184614e80565b6020810161392b8284614dd9565b602080825282518282018190525f919060409081850190868401855b82811015614e7357815180516001600160a01b039081168652878201518116888701528682015116868601526060908101519085015260809093019290850190600101615086565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f6139e260208301846150ce565b634e487b7160e01b5f52604160045260245ffd5b604080519081016001600160401b03811182821017156151445761514461510e565b60405290565b60405161014081016001600160401b03811182821017156151445761514461510e565b6040516101e081016001600160401b03811182821017156151445761514461510e565b604051601f8201601f191681016001600160401b03811182821017156151b8576151b861510e565b604052919050565b5f6001600160401b038211156151d8576151d861510e565b50601f01601f191660200190565b5f82601f8301126151f5575f80fd5b8135615208615203826151c0565b615190565b81815284602083860101111561521c575f80fd5b816020850160208301375f918101602001919091529392505050565b5f8060408385031215615249575f80fd5b823561525481614d4f565b915060208301356001600160401b0381111561526e575f80fd5b61527a858286016151e6565b9150509250929050565b5f81516040845261529860408501826150ce565b9050602083015184820360208601526152b182826150ce565b95945050505050565b80516001600160a01b031682525f61014060208301516152e560208601826001600160a01b03169052565b50604083015161530060408601826001600160a01b03169052565b50606083015161531b60608601826001600160a01b03169052565b506080830151608085015260a083015160a085015260c08301518160c086015261534782860182615284565b91505060e083015184820360e086015261536182826150ce565b915050610100808401518583038287015261537c83826150ce565b925050506101208084015161539382870182614dd9565b5090949350505050565b602081525f6139e260208301846152ba565b604081525f6153c16040830185614e80565b6020838203818501528185518084528284019150828160051b8501018388015f5b8381101561541057601f198784030185526153fe838351615284565b948601949250908501906001016153e2565b50909998505050505050505050565b604080825283519082018190525f906020906060840190828701845b8281101561545a57815160ff168452928401929084019060010161543b565b505050838103602085015261546f8186614e80565b9695505050505050565b5f806020838503121561548a575f80fd5b82356001600160401b0381111561549f575f80fd5b6154ab85828601614f69565b90969095509350505050565b5f60208083016020845280855180835260408601915060408160051b8701019250602087015f5b8281101561550c57603f198886030184526154fa8583516152ba565b945092850192908501906001016154de565b5092979650505050505050565b5f805f6040848603121561552b575f80fd5b833561553681614d4f565b925060208401356001600160401b0380821115615551575f80fd5b818601915086601f830112615564575f80fd5b813581811115615572575f80fd5b876020828501011115615583575f80fd5b6020830194508093505050509250925092565b5f604082840312156155a6575f80fd5b6155ae615122565b905081356001600160401b03808211156155c6575f80fd5b6155d2858386016151e6565b835260208401359150808211156155e7575f80fd5b506155f4848285016151e6565b60208301525092915050565b5f805f805f805f60e0888a031215615616575f80fd5b873561562181614d4f565b9650602088013561563181614d4f565b955061563f60408901614d63565b945061564d60608901614d63565b935060808801356001600160401b0380821115615668575f80fd5b6156748b838c01615596565b945060a08a0135915080821115615689575f80fd5b6156958b838c016151e6565b935060c08a01359150808211156156aa575f80fd5b506156b78a828b016151e6565b91505092959891949750929550565b634e487b7160e01b5f52601160045260245ffd5b5f600182016156eb576156eb6156c6565b5060010190565b634e487b7160e01b5f52603260045260245ffd5b604080825281018490525f8560608301825b8781101561574857823561572b81614d4f565b6001600160a01b0316825260209283019290910190600101615718565b508381036020858101919091528582529150859082015f5b8681101561578e57823561577381614f42565b61577d8382614dd9565b509183019190830190600101615760565b5098975050505050505050565b5f816157a9576157a96156c6565b505f190190565b8051614d6e81614d4f565b5f602082840312156157cb575f80fd5b81516139e281614d4f565b600181811c908216806157ea57607f821691505b60208210810361580857634e487b7160e01b5f52602260045260245ffd5b50919050565b808202811582820484141761392b5761392b6156c6565b8082018082111561392b5761392b6156c6565b5f6001600160401b038211156158505761585061510e565b5060051b60200190565b5f82601f830112615869575f80fd5b8151602061587961520383615838565b8083825260208201915060208460051b87010193508684111561589a575f80fd5b602086015b848110156158bf5780516158b281614d4f565b835291830191830161589f565b509695505050505050565b5f82601f8301126158d9575f80fd5b81516158e7615203826151c0565b8181528460208386010111156158fb575f80fd5b8160208501602083015e5f918101602001919091529392505050565b5f60408284031215615927575f80fd5b61592f615122565b905081516001600160401b0380821115615947575f80fd5b615953858386016158ca565b83526020840151915080821115615968575f80fd5b506155f4848285016158ca565b8051614d6e81614f42565b5f82601f83011261598f575f80fd5b8151602061599f61520383615838565b82815260059290921b840181019181810190868411156159bd575f80fd5b8286015b848110156158bf5780516001600160401b03808211156159df575f80fd5b90880190610140828b03601f19018113156159f8575f80fd5b615a0061514a565b615a0b8885016157b0565b81526040615a1a8186016157b0565b898301526060615a2b8187016157b0565b8284015260809150615a3e8287016157b0565b818401525060a0808601518284015260c0915081860151818401525060e08086015185811115615a6c575f80fd5b615a7a8f8c838a0101615917565b838501525061010091508186015185811115615a94575f80fd5b615aa28f8c838a01016158ca565b8285015250506101208086015185811115615abb575f80fd5b615ac98f8c838a01016158ca565b8385015250615ad9848701615975565b908301525086525050509183019183016159c1565b5f60208284031215615afe575f80fd5b81516001600160401b0380821115615b14575f80fd5b908301906101e08286031215615b28575f80fd5b615b3061516d565b615b39836157b0565b8152615b47602084016157b0565b6020820152615b58604084016157b0565b6040820152606083015160608201526080830151608082015260a083015160a082015260c083015160c082015260e083015160e0820152610100808401518183015250610120808401518183015250610140615bb58185016157b0565b90820152610160615bc78482016157b0565b90820152610180615bd98482016157b0565b908201526101a08381015183811115615bf0575f80fd5b615bfc8882870161585a565b8284015250506101c08084015183811115615c15575f80fd5b615c2188828701615980565b918301919091525095945050505050565b601f8211156141ec57805f5260205f20601f840160051c81016020851015615c575750805b601f840160051c820191505b81811015612b27575f8155600101615c63565b81516001600160401b03811115615c8f57615c8f61510e565b615ca381615c9d84546157d6565b84615c32565b602080601f831160018114615cd6575f8415615cbf5750858301515b5f19600386901b1c1916600185901b178555614552565b5f85815260208120601f198616915b82811015615d0457888601518255948401946001909101908401615ce5565b5085821015615d2157878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b6001600160401b03831115615d4857615d4861510e565b615d5c83615d5683546157d6565b83615c32565b5f601f841160018114615d8d575f8515615d765750838201355b5f19600387901b1c1916600186901b178355612b27565b5f83815260208120601f198716915b82811015615dbc5786850135825560209485019460019092019101615d9c565b5086821015615dd8575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b60208152816020820152818360408301375f818301604090810191909152601f909201601f19160101919050565b634e487b7160e01b5f52601260045260245ffd5b5f82615e3a57615e3a615e18565b500490565b5f82615e4d57615e4d615e18565b500690565b5f60208284031215615e62575f80fd5b5051919050565b8181038181111561392b5761392b6156c6565b634e487b7160e01b5f52603160045260245ffdfe34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbcd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4ba164736f6c6343000819000a`

// AddressBookV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AddressBookV2MetaData.Bin instead.
var AddressBookV2Bin = AddressBookV2MetaData.Bin

// DeployAddressBookV2 deploys a new Kaia contract, binding an instance of AddressBookV2 to it.
func DeployAddressBookV2(auth *bind.TransactOpts, backend bind.ContractBackend, _epochBlockInterval *big.Int) (common.Address, *types.Transaction, *AddressBookV2, error) {
	parsed, err := AddressBookV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AddressBookV2Bin), backend, _epochBlockInterval)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AddressBookV2{AddressBookV2Caller: AddressBookV2Caller{contract: contract}, AddressBookV2Transactor: AddressBookV2Transactor{contract: contract}, AddressBookV2Filterer: AddressBookV2Filterer{contract: contract}}, nil
}

// AddressBookV2 is an auto generated Go binding around a Kaia contract.
type AddressBookV2 struct {
	AddressBookV2Caller     // Read-only binding to the contract
	AddressBookV2Transactor // Write-only binding to the contract
	AddressBookV2Filterer   // Log filterer for contract events
}

// AddressBookV2Caller is an auto generated read-only Go binding around a Kaia contract.
type AddressBookV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressBookV2Transactor is an auto generated write-only Go binding around a Kaia contract.
type AddressBookV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressBookV2Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type AddressBookV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressBookV2Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type AddressBookV2Session struct {
	Contract     *AddressBookV2    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AddressBookV2CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type AddressBookV2CallerSession struct {
	Contract *AddressBookV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// AddressBookV2TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type AddressBookV2TransactorSession struct {
	Contract     *AddressBookV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// AddressBookV2Raw is an auto generated low-level Go binding around a Kaia contract.
type AddressBookV2Raw struct {
	Contract *AddressBookV2 // Generic contract binding to access the raw methods on
}

// AddressBookV2CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type AddressBookV2CallerRaw struct {
	Contract *AddressBookV2Caller // Generic read-only contract binding to access the raw methods on
}

// AddressBookV2TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type AddressBookV2TransactorRaw struct {
	Contract *AddressBookV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewAddressBookV2 creates a new instance of AddressBookV2, bound to a specific deployed contract.
func NewAddressBookV2(address common.Address, backend bind.ContractBackend) (*AddressBookV2, error) {
	contract, err := bindAddressBookV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2{AddressBookV2Caller: AddressBookV2Caller{contract: contract}, AddressBookV2Transactor: AddressBookV2Transactor{contract: contract}, AddressBookV2Filterer: AddressBookV2Filterer{contract: contract}}, nil
}

// NewAddressBookV2Caller creates a new read-only instance of AddressBookV2, bound to a specific deployed contract.
func NewAddressBookV2Caller(address common.Address, caller bind.ContractCaller) (*AddressBookV2Caller, error) {
	contract, err := bindAddressBookV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2Caller{contract: contract}, nil
}

// NewAddressBookV2Transactor creates a new write-only instance of AddressBookV2, bound to a specific deployed contract.
func NewAddressBookV2Transactor(address common.Address, transactor bind.ContractTransactor) (*AddressBookV2Transactor, error) {
	contract, err := bindAddressBookV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2Transactor{contract: contract}, nil
}

// NewAddressBookV2Filterer creates a new log filterer instance of AddressBookV2, bound to a specific deployed contract.
func NewAddressBookV2Filterer(address common.Address, filterer bind.ContractFilterer) (*AddressBookV2Filterer, error) {
	contract, err := bindAddressBookV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2Filterer{contract: contract}, nil
}

// bindAddressBookV2 binds a generic wrapper to an already deployed contract.
func bindAddressBookV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AddressBookV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AddressBookV2 *AddressBookV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AddressBookV2.Contract.AddressBookV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AddressBookV2 *AddressBookV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AddressBookV2.Contract.AddressBookV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AddressBookV2 *AddressBookV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AddressBookV2.Contract.AddressBookV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AddressBookV2 *AddressBookV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AddressBookV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AddressBookV2 *AddressBookV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AddressBookV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AddressBookV2 *AddressBookV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AddressBookV2.Contract.contract.Transact(opts, method, params...)
}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() view returns(string)
func (_AddressBookV2 *AddressBookV2Caller) CONTRACTTYPE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "CONTRACT_TYPE")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() view returns(string)
func (_AddressBookV2 *AddressBookV2Session) CONTRACTTYPE() (string, error) {
	return _AddressBookV2.Contract.CONTRACTTYPE(&_AddressBookV2.CallOpts)
}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() view returns(string)
func (_AddressBookV2 *AddressBookV2CallerSession) CONTRACTTYPE() (string, error) {
	return _AddressBookV2.Contract.CONTRACTTYPE(&_AddressBookV2.CallOpts)
}

// MAXMETADATALENGTH is a free data retrieval call binding the contract method 0xe8868e9f.
//
// Solidity: function MAX_METADATA_LENGTH() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) MAXMETADATALENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "MAX_METADATA_LENGTH")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// MAXMETADATALENGTH is a free data retrieval call binding the contract method 0xe8868e9f.
//
// Solidity: function MAX_METADATA_LENGTH() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) MAXMETADATALENGTH() (*big.Int, error) {
	return _AddressBookV2.Contract.MAXMETADATALENGTH(&_AddressBookV2.CallOpts)
}

// MAXMETADATALENGTH is a free data retrieval call binding the contract method 0xe8868e9f.
//
// Solidity: function MAX_METADATA_LENGTH() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) MAXMETADATALENGTH() (*big.Int, error) {
	return _AddressBookV2.Contract.MAXMETADATALENGTH(&_AddressBookV2.CallOpts)
}

// MINNODEBALANCE is a free data retrieval call binding the contract method 0xc732e085.
//
// Solidity: function MIN_NODE_BALANCE() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) MINNODEBALANCE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "MIN_NODE_BALANCE")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// MINNODEBALANCE is a free data retrieval call binding the contract method 0xc732e085.
//
// Solidity: function MIN_NODE_BALANCE() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) MINNODEBALANCE() (*big.Int, error) {
	return _AddressBookV2.Contract.MINNODEBALANCE(&_AddressBookV2.CallOpts)
}

// MINNODEBALANCE is a free data retrieval call binding the contract method 0xc732e085.
//
// Solidity: function MIN_NODE_BALANCE() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) MINNODEBALANCE() (*big.Int, error) {
	return _AddressBookV2.Contract.MINNODEBALANCE(&_AddressBookV2.CallOpts)
}

// MINSTAKE is a free data retrieval call binding the contract method 0xcb1c2b5c.
//
// Solidity: function MIN_STAKE() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) MINSTAKE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "MIN_STAKE")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// MINSTAKE is a free data retrieval call binding the contract method 0xcb1c2b5c.
//
// Solidity: function MIN_STAKE() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) MINSTAKE() (*big.Int, error) {
	return _AddressBookV2.Contract.MINSTAKE(&_AddressBookV2.CallOpts)
}

// MINSTAKE is a free data retrieval call binding the contract method 0xcb1c2b5c.
//
// Solidity: function MIN_STAKE() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) MINSTAKE() (*big.Int, error) {
	return _AddressBookV2.Contract.MINSTAKE(&_AddressBookV2.CallOpts)
}

// SYSTEMSENDER is a free data retrieval call binding the contract method 0x87b7b8fd.
//
// Solidity: function SYSTEM_SENDER() view returns(address)
func (_AddressBookV2 *AddressBookV2Caller) SYSTEMSENDER(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "SYSTEM_SENDER")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// SYSTEMSENDER is a free data retrieval call binding the contract method 0x87b7b8fd.
//
// Solidity: function SYSTEM_SENDER() view returns(address)
func (_AddressBookV2 *AddressBookV2Session) SYSTEMSENDER() (common.Address, error) {
	return _AddressBookV2.Contract.SYSTEMSENDER(&_AddressBookV2.CallOpts)
}

// SYSTEMSENDER is a free data retrieval call binding the contract method 0x87b7b8fd.
//
// Solidity: function SYSTEM_SENDER() view returns(address)
func (_AddressBookV2 *AddressBookV2CallerSession) SYSTEMSENDER() (common.Address, error) {
	return _AddressBookV2.Contract.SYSTEMSENDER(&_AddressBookV2.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AddressBookV2 *AddressBookV2Caller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AddressBookV2 *AddressBookV2Session) UPGRADEINTERFACEVERSION() (string, error) {
	return _AddressBookV2.Contract.UPGRADEINTERFACEVERSION(&_AddressBookV2.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_AddressBookV2 *AddressBookV2CallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _AddressBookV2.Contract.UPGRADEINTERFACEVERSION(&_AddressBookV2.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) VERSION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "VERSION")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) VERSION() (*big.Int, error) {
	return _AddressBookV2.Contract.VERSION(&_AddressBookV2.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) VERSION() (*big.Int, error) {
	return _AddressBookV2.Contract.VERSION(&_AddressBookV2.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "currentEpoch")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) CurrentEpoch() (*big.Int, error) {
	return _AddressBookV2.Contract.CurrentEpoch(&_AddressBookV2.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) CurrentEpoch() (*big.Int, error) {
	return _AddressBookV2.Contract.CurrentEpoch(&_AddressBookV2.CallOpts)
}

// EpochBlockInterval is a free data retrieval call binding the contract method 0x567b0b6c.
//
// Solidity: function epochBlockInterval() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) EpochBlockInterval(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "epochBlockInterval")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// EpochBlockInterval is a free data retrieval call binding the contract method 0x567b0b6c.
//
// Solidity: function epochBlockInterval() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) EpochBlockInterval() (*big.Int, error) {
	return _AddressBookV2.Contract.EpochBlockInterval(&_AddressBookV2.CallOpts)
}

// EpochBlockInterval is a free data retrieval call binding the contract method 0x567b0b6c.
//
// Solidity: function epochBlockInterval() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) EpochBlockInterval() (*big.Int, error) {
	return _AddressBookV2.Contract.EpochBlockInterval(&_AddressBookV2.CallOpts)
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_AddressBookV2 *AddressBookV2Caller) GetAllAddress(opts *bind.CallOpts) (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error,
) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getAllAddress")

	outstruct := new(struct {
		TypeList    []uint8
		AddressList []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TypeList = *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)
	outstruct.AddressList = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_AddressBookV2 *AddressBookV2Session) GetAllAddress() (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error,
) {
	return _AddressBookV2.Contract.GetAllAddress(&_AddressBookV2.CallOpts)
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_AddressBookV2 *AddressBookV2CallerSession) GetAllAddress() (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error,
) {
	return _AddressBookV2.Contract.GetAllAddress(&_AddressBookV2.CallOpts)
}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[], address[], address[], address, address)
func (_AddressBookV2 *AddressBookV2Caller) GetAllAddressInfo(opts *bind.CallOpts) ([]common.Address, []common.Address, []common.Address, common.Address, common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getAllAddressInfo")
	if err != nil {
		return *new([]common.Address), *new([]common.Address), *new([]common.Address), *new(common.Address), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	out2 := *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	out3 := *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	out4 := *abi.ConvertType(out[4], new(common.Address)).(*common.Address)

	return out0, out1, out2, out3, out4, err
}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[], address[], address[], address, address)
func (_AddressBookV2 *AddressBookV2Session) GetAllAddressInfo() ([]common.Address, []common.Address, []common.Address, common.Address, common.Address, error) {
	return _AddressBookV2.Contract.GetAllAddressInfo(&_AddressBookV2.CallOpts)
}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[], address[], address[], address, address)
func (_AddressBookV2 *AddressBookV2CallerSession) GetAllAddressInfo() ([]common.Address, []common.Address, []common.Address, common.Address, common.Address, error) {
	return _AddressBookV2.Contract.GetAllAddressInfo(&_AddressBookV2.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_AddressBookV2 *AddressBookV2Caller) GetAllBlsInfo(opts *bind.CallOpts) (struct {
	NodeIdList []common.Address
	PubkeyList []BlsPublicKeyInfo
}, error,
) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getAllBlsInfo")

	outstruct := new(struct {
		NodeIdList []common.Address
		PubkeyList []BlsPublicKeyInfo
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NodeIdList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.PubkeyList = *abi.ConvertType(out[1], new([]BlsPublicKeyInfo)).(*[]BlsPublicKeyInfo)

	return *outstruct, err
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_AddressBookV2 *AddressBookV2Session) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []BlsPublicKeyInfo
}, error,
) {
	return _AddressBookV2.Contract.GetAllBlsInfo(&_AddressBookV2.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_AddressBookV2 *AddressBookV2CallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []BlsPublicKeyInfo
}, error,
) {
	return _AddressBookV2.Contract.GetAllBlsInfo(&_AddressBookV2.CallOpts)
}

// GetAllGovernanceInfo is a free data retrieval call binding the contract method 0x2aca5091.
//
// Solidity: function getAllGovernanceInfo() view returns((address,address,address,uint256)[] infos)
func (_AddressBookV2 *AddressBookV2Caller) GetAllGovernanceInfo(opts *bind.CallOpts) ([]GovernanceInfo, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getAllGovernanceInfo")
	if err != nil {
		return *new([]GovernanceInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]GovernanceInfo)).(*[]GovernanceInfo)

	return out0, err
}

// GetAllGovernanceInfo is a free data retrieval call binding the contract method 0x2aca5091.
//
// Solidity: function getAllGovernanceInfo() view returns((address,address,address,uint256)[] infos)
func (_AddressBookV2 *AddressBookV2Session) GetAllGovernanceInfo() ([]GovernanceInfo, error) {
	return _AddressBookV2.Contract.GetAllGovernanceInfo(&_AddressBookV2.CallOpts)
}

// GetAllGovernanceInfo is a free data retrieval call binding the contract method 0x2aca5091.
//
// Solidity: function getAllGovernanceInfo() view returns((address,address,address,uint256)[] infos)
func (_AddressBookV2 *AddressBookV2CallerSession) GetAllGovernanceInfo() ([]GovernanceInfo, error) {
	return _AddressBookV2.Contract.GetAllGovernanceInfo(&_AddressBookV2.CallOpts)
}

// GetAllNodesLength is a free data retrieval call binding the contract method 0x0a4ff239.
//
// Solidity: function getAllNodesLength() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetAllNodesLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getAllNodesLength")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetAllNodesLength is a free data retrieval call binding the contract method 0x0a4ff239.
//
// Solidity: function getAllNodesLength() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) GetAllNodesLength() (*big.Int, error) {
	return _AddressBookV2.Contract.GetAllNodesLength(&_AddressBookV2.CallOpts)
}

// GetAllNodesLength is a free data retrieval call binding the contract method 0x0a4ff239.
//
// Solidity: function getAllNodesLength() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetAllNodesLength() (*big.Int, error) {
	return _AddressBookV2.Contract.GetAllNodesLength(&_AddressBookV2.CallOpts)
}

// GetAllProfiles is a free data retrieval call binding the contract method 0x0b1fe784.
//
// Solidity: function getAllProfiles() view returns((address,address,address,uint256,uint8)[] profiles)
func (_AddressBookV2 *AddressBookV2Caller) GetAllProfiles(opts *bind.CallOpts) ([]Profile, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getAllProfiles")
	if err != nil {
		return *new([]Profile), err
	}

	out0 := *abi.ConvertType(out[0], new([]Profile)).(*[]Profile)

	return out0, err
}

// GetAllProfiles is a free data retrieval call binding the contract method 0x0b1fe784.
//
// Solidity: function getAllProfiles() view returns((address,address,address,uint256,uint8)[] profiles)
func (_AddressBookV2 *AddressBookV2Session) GetAllProfiles() ([]Profile, error) {
	return _AddressBookV2.Contract.GetAllProfiles(&_AddressBookV2.CallOpts)
}

// GetAllProfiles is a free data retrieval call binding the contract method 0x0b1fe784.
//
// Solidity: function getAllProfiles() view returns((address,address,address,uint256,uint8)[] profiles)
func (_AddressBookV2 *AddressBookV2CallerSession) GetAllProfiles() ([]Profile, error) {
	return _AddressBookV2.Contract.GetAllProfiles(&_AddressBookV2.CallOpts)
}

// GetCfsThreshold is a free data retrieval call binding the contract method 0xf0a92ba8.
//
// Solidity: function getCfsThreshold() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetCfsThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getCfsThreshold")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetCfsThreshold is a free data retrieval call binding the contract method 0xf0a92ba8.
//
// Solidity: function getCfsThreshold() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) GetCfsThreshold() (*big.Int, error) {
	return _AddressBookV2.Contract.GetCfsThreshold(&_AddressBookV2.CallOpts)
}

// GetCfsThreshold is a free data retrieval call binding the contract method 0xf0a92ba8.
//
// Solidity: function getCfsThreshold() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetCfsThreshold() (*big.Int, error) {
	return _AddressBookV2.Contract.GetCfsThreshold(&_AddressBookV2.CallOpts)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address, address, address)
func (_AddressBookV2 *AddressBookV2Caller) GetCnInfo(opts *bind.CallOpts, _cnNodeId common.Address) (common.Address, common.Address, common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getCnInfo", _cnNodeId)
	if err != nil {
		return *new(common.Address), *new(common.Address), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return out0, out1, out2, err
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address, address, address)
func (_AddressBookV2 *AddressBookV2Session) GetCnInfo(_cnNodeId common.Address) (common.Address, common.Address, common.Address, error) {
	return _AddressBookV2.Contract.GetCnInfo(&_AddressBookV2.CallOpts, _cnNodeId)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address, address, address)
func (_AddressBookV2 *AddressBookV2CallerSession) GetCnInfo(_cnNodeId common.Address) (common.Address, common.Address, common.Address, error) {
	return _AddressBookV2.Contract.GetCnInfo(&_AddressBookV2.CallOpts, _cnNodeId)
}

// GetConfigurator is a free data retrieval call binding the contract method 0x9b7ae5ec.
//
// Solidity: function getConfigurator() view returns(address)
func (_AddressBookV2 *AddressBookV2Caller) GetConfigurator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getConfigurator")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// GetConfigurator is a free data retrieval call binding the contract method 0x9b7ae5ec.
//
// Solidity: function getConfigurator() view returns(address)
func (_AddressBookV2 *AddressBookV2Session) GetConfigurator() (common.Address, error) {
	return _AddressBookV2.Contract.GetConfigurator(&_AddressBookV2.CallOpts)
}

// GetConfigurator is a free data retrieval call binding the contract method 0x9b7ae5ec.
//
// Solidity: function getConfigurator() view returns(address)
func (_AddressBookV2 *AddressBookV2CallerSession) GetConfigurator() (common.Address, error) {
	return _AddressBookV2.Contract.GetConfigurator(&_AddressBookV2.CallOpts)
}

// GetEpochVACount is a free data retrieval call binding the contract method 0xd3b54907.
//
// Solidity: function getEpochVACount() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetEpochVACount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getEpochVACount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetEpochVACount is a free data retrieval call binding the contract method 0xd3b54907.
//
// Solidity: function getEpochVACount() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) GetEpochVACount() (*big.Int, error) {
	return _AddressBookV2.Contract.GetEpochVACount(&_AddressBookV2.CallOpts)
}

// GetEpochVACount is a free data retrieval call binding the contract method 0xd3b54907.
//
// Solidity: function getEpochVACount() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetEpochVACount() (*big.Int, error) {
	return _AddressBookV2.Contract.GetEpochVACount(&_AddressBookV2.CallOpts)
}

// GetFundAddresses is a free data retrieval call binding the contract method 0x25cf0943.
//
// Solidity: function getFundAddresses() view returns(address, address, address)
func (_AddressBookV2 *AddressBookV2Caller) GetFundAddresses(opts *bind.CallOpts) (common.Address, common.Address, common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getFundAddresses")
	if err != nil {
		return *new(common.Address), *new(common.Address), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return out0, out1, out2, err
}

// GetFundAddresses is a free data retrieval call binding the contract method 0x25cf0943.
//
// Solidity: function getFundAddresses() view returns(address, address, address)
func (_AddressBookV2 *AddressBookV2Session) GetFundAddresses() (common.Address, common.Address, common.Address, error) {
	return _AddressBookV2.Contract.GetFundAddresses(&_AddressBookV2.CallOpts)
}

// GetFundAddresses is a free data retrieval call binding the contract method 0x25cf0943.
//
// Solidity: function getFundAddresses() view returns(address, address, address)
func (_AddressBookV2 *AddressBookV2CallerSession) GetFundAddresses() (common.Address, common.Address, common.Address, error) {
	return _AddressBookV2.Contract.GetFundAddresses(&_AddressBookV2.CallOpts)
}

// GetMaxCounts is a free data retrieval call binding the contract method 0x03e6689d.
//
// Solidity: function getMaxCounts() view returns(uint256, uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetMaxCounts(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getMaxCounts")
	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err
}

// GetMaxCounts is a free data retrieval call binding the contract method 0x03e6689d.
//
// Solidity: function getMaxCounts() view returns(uint256, uint256)
func (_AddressBookV2 *AddressBookV2Session) GetMaxCounts() (*big.Int, *big.Int, error) {
	return _AddressBookV2.Contract.GetMaxCounts(&_AddressBookV2.CallOpts)
}

// GetMaxCounts is a free data retrieval call binding the contract method 0x03e6689d.
//
// Solidity: function getMaxCounts() view returns(uint256, uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetMaxCounts() (*big.Int, *big.Int, error) {
	return _AddressBookV2.Contract.GetMaxCounts(&_AddressBookV2.CallOpts)
}

// GetMaxValActivePausedCount is a free data retrieval call binding the contract method 0x8fabf389.
//
// Solidity: function getMaxValActivePausedCount() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetMaxValActivePausedCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getMaxValActivePausedCount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetMaxValActivePausedCount is a free data retrieval call binding the contract method 0x8fabf389.
//
// Solidity: function getMaxValActivePausedCount() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) GetMaxValActivePausedCount() (*big.Int, error) {
	return _AddressBookV2.Contract.GetMaxValActivePausedCount(&_AddressBookV2.CallOpts)
}

// GetMaxValActivePausedCount is a free data retrieval call binding the contract method 0x8fabf389.
//
// Solidity: function getMaxValActivePausedCount() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetMaxValActivePausedCount() (*big.Int, error) {
	return _AddressBookV2.Contract.GetMaxValActivePausedCount(&_AddressBookV2.CallOpts)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0x582115fb.
//
// Solidity: function getNodeInfo(address nodeId) view returns((address,address,address,address,uint256,uint256,(bytes,bytes),string,string,uint8))
func (_AddressBookV2 *AddressBookV2Caller) GetNodeInfo(opts *bind.CallOpts, nodeId common.Address) (NodeInfo, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getNodeInfo", nodeId)
	if err != nil {
		return *new(NodeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(NodeInfo)).(*NodeInfo)

	return out0, err
}

// GetNodeInfo is a free data retrieval call binding the contract method 0x582115fb.
//
// Solidity: function getNodeInfo(address nodeId) view returns((address,address,address,address,uint256,uint256,(bytes,bytes),string,string,uint8))
func (_AddressBookV2 *AddressBookV2Session) GetNodeInfo(nodeId common.Address) (NodeInfo, error) {
	return _AddressBookV2.Contract.GetNodeInfo(&_AddressBookV2.CallOpts, nodeId)
}

// GetNodeInfo is a free data retrieval call binding the contract method 0x582115fb.
//
// Solidity: function getNodeInfo(address nodeId) view returns((address,address,address,address,uint256,uint256,(bytes,bytes),string,string,uint8))
func (_AddressBookV2 *AddressBookV2CallerSession) GetNodeInfo(nodeId common.Address) (NodeInfo, error) {
	return _AddressBookV2.Contract.GetNodeInfo(&_AddressBookV2.CallOpts, nodeId)
}

// GetNodeInfos is a free data retrieval call binding the contract method 0x8beeb439.
//
// Solidity: function getNodeInfos(address[] nodeIds) view returns((address,address,address,address,uint256,uint256,(bytes,bytes),string,string,uint8)[] infos)
func (_AddressBookV2 *AddressBookV2Caller) GetNodeInfos(opts *bind.CallOpts, nodeIds []common.Address) ([]NodeInfo, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getNodeInfos", nodeIds)
	if err != nil {
		return *new([]NodeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new([]NodeInfo)).(*[]NodeInfo)

	return out0, err
}

// GetNodeInfos is a free data retrieval call binding the contract method 0x8beeb439.
//
// Solidity: function getNodeInfos(address[] nodeIds) view returns((address,address,address,address,uint256,uint256,(bytes,bytes),string,string,uint8)[] infos)
func (_AddressBookV2 *AddressBookV2Session) GetNodeInfos(nodeIds []common.Address) ([]NodeInfo, error) {
	return _AddressBookV2.Contract.GetNodeInfos(&_AddressBookV2.CallOpts, nodeIds)
}

// GetNodeInfos is a free data retrieval call binding the contract method 0x8beeb439.
//
// Solidity: function getNodeInfos(address[] nodeIds) view returns((address,address,address,address,uint256,uint256,(bytes,bytes),string,string,uint8)[] infos)
func (_AddressBookV2 *AddressBookV2CallerSession) GetNodeInfos(nodeIds []common.Address) ([]NodeInfo, error) {
	return _AddressBookV2.Contract.GetNodeInfos(&_AddressBookV2.CallOpts, nodeIds)
}

// GetNodeState is a free data retrieval call binding the contract method 0x229bb823.
//
// Solidity: function getNodeState(address nodeId) view returns(uint8)
func (_AddressBookV2 *AddressBookV2Caller) GetNodeState(opts *bind.CallOpts, nodeId common.Address) (uint8, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getNodeState", nodeId)
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// GetNodeState is a free data retrieval call binding the contract method 0x229bb823.
//
// Solidity: function getNodeState(address nodeId) view returns(uint8)
func (_AddressBookV2 *AddressBookV2Session) GetNodeState(nodeId common.Address) (uint8, error) {
	return _AddressBookV2.Contract.GetNodeState(&_AddressBookV2.CallOpts, nodeId)
}

// GetNodeState is a free data retrieval call binding the contract method 0x229bb823.
//
// Solidity: function getNodeState(address nodeId) view returns(uint8)
func (_AddressBookV2 *AddressBookV2CallerSession) GetNodeState(nodeId common.Address) (uint8, error) {
	return _AddressBookV2.Contract.GetNodeState(&_AddressBookV2.CallOpts, nodeId)
}

// GetPfsThreshold is a free data retrieval call binding the contract method 0x291937f5.
//
// Solidity: function getPfsThreshold() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetPfsThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getPfsThreshold")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetPfsThreshold is a free data retrieval call binding the contract method 0x291937f5.
//
// Solidity: function getPfsThreshold() view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) GetPfsThreshold() (*big.Int, error) {
	return _AddressBookV2.Contract.GetPfsThreshold(&_AddressBookV2.CallOpts)
}

// GetPfsThreshold is a free data retrieval call binding the contract method 0x291937f5.
//
// Solidity: function getPfsThreshold() view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetPfsThreshold() (*big.Int, error) {
	return _AddressBookV2.Contract.GetPfsThreshold(&_AddressBookV2.CallOpts)
}

// GetRegisteredNodes is a free data retrieval call binding the contract method 0xcf8c6f52.
//
// Solidity: function getRegisteredNodes() view returns(address[])
func (_AddressBookV2 *AddressBookV2Caller) GetRegisteredNodes(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getRegisteredNodes")
	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err
}

// GetRegisteredNodes is a free data retrieval call binding the contract method 0xcf8c6f52.
//
// Solidity: function getRegisteredNodes() view returns(address[])
func (_AddressBookV2 *AddressBookV2Session) GetRegisteredNodes() ([]common.Address, error) {
	return _AddressBookV2.Contract.GetRegisteredNodes(&_AddressBookV2.CallOpts)
}

// GetRegisteredNodes is a free data retrieval call binding the contract method 0xcf8c6f52.
//
// Solidity: function getRegisteredNodes() view returns(address[])
func (_AddressBookV2 *AddressBookV2CallerSession) GetRegisteredNodes() ([]common.Address, error) {
	return _AddressBookV2.Contract.GetRegisteredNodes(&_AddressBookV2.CallOpts)
}

// GetSlotLimits is a free data retrieval call binding the contract method 0x9d0f5ef1.
//
// Solidity: function getSlotLimits() view returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2Caller) GetSlotLimits(opts *bind.CallOpts) (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error,
) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getSlotLimits")

	outstruct := new(struct {
		MaxSlotAvailable *big.Int
		MinActiveCount   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxSlotAvailable = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.MinActiveCount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err
}

// GetSlotLimits is a free data retrieval call binding the contract method 0x9d0f5ef1.
//
// Solidity: function getSlotLimits() view returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2Session) GetSlotLimits() (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error,
) {
	return _AddressBookV2.Contract.GetSlotLimits(&_AddressBookV2.CallOpts)
}

// GetSlotLimits is a free data retrieval call binding the contract method 0x9d0f5ef1.
//
// Solidity: function getSlotLimits() view returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2CallerSession) GetSlotLimits() (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error,
) {
	return _AddressBookV2.Contract.GetSlotLimits(&_AddressBookV2.CallOpts)
}

// GetSlotLimitsFor is a free data retrieval call binding the contract method 0x058529fb.
//
// Solidity: function getSlotLimitsFor(uint256 n) pure returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2Caller) GetSlotLimitsFor(opts *bind.CallOpts, n *big.Int) (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error,
) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getSlotLimitsFor", n)

	outstruct := new(struct {
		MaxSlotAvailable *big.Int
		MinActiveCount   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxSlotAvailable = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.MinActiveCount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err
}

// GetSlotLimitsFor is a free data retrieval call binding the contract method 0x058529fb.
//
// Solidity: function getSlotLimitsFor(uint256 n) pure returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2Session) GetSlotLimitsFor(n *big.Int) (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error,
) {
	return _AddressBookV2.Contract.GetSlotLimitsFor(&_AddressBookV2.CallOpts, n)
}

// GetSlotLimitsFor is a free data retrieval call binding the contract method 0x058529fb.
//
// Solidity: function getSlotLimitsFor(uint256 n) pure returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2CallerSession) GetSlotLimitsFor(n *big.Int) (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error,
) {
	return _AddressBookV2.Contract.GetSlotLimitsFor(&_AddressBookV2.CallOpts, n)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[], uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetState(opts *bind.CallOpts) ([]common.Address, *big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getState")
	if err != nil {
		return *new([]common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[], uint256)
func (_AddressBookV2 *AddressBookV2Session) GetState() ([]common.Address, *big.Int, error) {
	return _AddressBookV2.Contract.GetState(&_AddressBookV2.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[], uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetState() ([]common.Address, *big.Int, error) {
	return _AddressBookV2.Contract.GetState(&_AddressBookV2.CallOpts)
}

// GetStateCount is a free data retrieval call binding the contract method 0x1b1a478b.
//
// Solidity: function getStateCount(uint8 state) view returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetStateCount(opts *bind.CallOpts, state uint8) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getStateCount", state)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetStateCount is a free data retrieval call binding the contract method 0x1b1a478b.
//
// Solidity: function getStateCount(uint8 state) view returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) GetStateCount(state uint8) (*big.Int, error) {
	return _AddressBookV2.Contract.GetStateCount(&_AddressBookV2.CallOpts, state)
}

// GetStateCount is a free data retrieval call binding the contract method 0x1b1a478b.
//
// Solidity: function getStateCount(uint8 state) view returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetStateCount(state uint8) (*big.Int, error) {
	return _AddressBookV2.Contract.GetStateCount(&_AddressBookV2.CallOpts, state)
}

// GetSuspendedValidators is a free data retrieval call binding the contract method 0x1ba3fd58.
//
// Solidity: function getSuspendedValidators() view returns(address[])
func (_AddressBookV2 *AddressBookV2Caller) GetSuspendedValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getSuspendedValidators")
	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err
}

// GetSuspendedValidators is a free data retrieval call binding the contract method 0x1ba3fd58.
//
// Solidity: function getSuspendedValidators() view returns(address[])
func (_AddressBookV2 *AddressBookV2Session) GetSuspendedValidators() ([]common.Address, error) {
	return _AddressBookV2.Contract.GetSuspendedValidators(&_AddressBookV2.CallOpts)
}

// GetSuspendedValidators is a free data retrieval call binding the contract method 0x1ba3fd58.
//
// Solidity: function getSuspendedValidators() view returns(address[])
func (_AddressBookV2 *AddressBookV2CallerSession) GetSuspendedValidators() ([]common.Address, error) {
	return _AddressBookV2.Contract.GetSuspendedValidators(&_AddressBookV2.CallOpts)
}

// GetSuspender is a free data retrieval call binding the contract method 0x21d23200.
//
// Solidity: function getSuspender() view returns(address)
func (_AddressBookV2 *AddressBookV2Caller) GetSuspender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getSuspender")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// GetSuspender is a free data retrieval call binding the contract method 0x21d23200.
//
// Solidity: function getSuspender() view returns(address)
func (_AddressBookV2 *AddressBookV2Session) GetSuspender() (common.Address, error) {
	return _AddressBookV2.Contract.GetSuspender(&_AddressBookV2.CallOpts)
}

// GetSuspender is a free data retrieval call binding the contract method 0x21d23200.
//
// Solidity: function getSuspender() view returns(address)
func (_AddressBookV2 *AddressBookV2CallerSession) GetSuspender() (common.Address, error) {
	return _AddressBookV2.Contract.GetSuspender(&_AddressBookV2.CallOpts)
}

// GetTimeouts is a free data retrieval call binding the contract method 0xe70c38f1.
//
// Solidity: function getTimeouts() view returns(uint256, uint256)
func (_AddressBookV2 *AddressBookV2Caller) GetTimeouts(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "getTimeouts")
	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err
}

// GetTimeouts is a free data retrieval call binding the contract method 0xe70c38f1.
//
// Solidity: function getTimeouts() view returns(uint256, uint256)
func (_AddressBookV2 *AddressBookV2Session) GetTimeouts() (*big.Int, *big.Int, error) {
	return _AddressBookV2.Contract.GetTimeouts(&_AddressBookV2.CallOpts)
}

// GetTimeouts is a free data retrieval call binding the contract method 0xe70c38f1.
//
// Solidity: function getTimeouts() view returns(uint256, uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) GetTimeouts() (*big.Int, *big.Int, error) {
	return _AddressBookV2.Contract.GetTimeouts(&_AddressBookV2.CallOpts)
}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_AddressBookV2 *AddressBookV2Caller) IsActivated(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "isActivated")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_AddressBookV2 *AddressBookV2Session) IsActivated() (bool, error) {
	return _AddressBookV2.Contract.IsActivated(&_AddressBookV2.CallOpts)
}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_AddressBookV2 *AddressBookV2CallerSession) IsActivated() (bool, error) {
	return _AddressBookV2.Contract.IsActivated(&_AddressBookV2.CallOpts)
}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_AddressBookV2 *AddressBookV2Caller) IsConstructed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "isConstructed")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_AddressBookV2 *AddressBookV2Session) IsConstructed() (bool, error) {
	return _AddressBookV2.Contract.IsConstructed(&_AddressBookV2.CallOpts)
}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_AddressBookV2 *AddressBookV2CallerSession) IsConstructed() (bool, error) {
	return _AddressBookV2.Contract.IsConstructed(&_AddressBookV2.CallOpts)
}

// IsUsedAddress is a free data retrieval call binding the contract method 0x468e3a7e.
//
// Solidity: function isUsedAddress(address addr) view returns(bool)
func (_AddressBookV2 *AddressBookV2Caller) IsUsedAddress(opts *bind.CallOpts, addr common.Address) (bool, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "isUsedAddress", addr)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsUsedAddress is a free data retrieval call binding the contract method 0x468e3a7e.
//
// Solidity: function isUsedAddress(address addr) view returns(bool)
func (_AddressBookV2 *AddressBookV2Session) IsUsedAddress(addr common.Address) (bool, error) {
	return _AddressBookV2.Contract.IsUsedAddress(&_AddressBookV2.CallOpts, addr)
}

// IsUsedAddress is a free data retrieval call binding the contract method 0x468e3a7e.
//
// Solidity: function isUsedAddress(address addr) view returns(bool)
func (_AddressBookV2 *AddressBookV2CallerSession) IsUsedAddress(addr common.Address) (bool, error) {
	return _AddressBookV2.Contract.IsUsedAddress(&_AddressBookV2.CallOpts, addr)
}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2Caller) KirContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "kirContractAddress")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2Session) KirContractAddress() (common.Address, error) {
	return _AddressBookV2.Contract.KirContractAddress(&_AddressBookV2.CallOpts)
}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2CallerSession) KirContractAddress() (common.Address, error) {
	return _AddressBookV2.Contract.KirContractAddress(&_AddressBookV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AddressBookV2 *AddressBookV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AddressBookV2 *AddressBookV2Session) Owner() (common.Address, error) {
	return _AddressBookV2.Contract.Owner(&_AddressBookV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AddressBookV2 *AddressBookV2CallerSession) Owner() (common.Address, error) {
	return _AddressBookV2.Contract.Owner(&_AddressBookV2.CallOpts)
}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2Caller) PocContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "pocContractAddress")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2Session) PocContractAddress() (common.Address, error) {
	return _AddressBookV2.Contract.PocContractAddress(&_AddressBookV2.CallOpts)
}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2CallerSession) PocContractAddress() (common.Address, error) {
	return _AddressBookV2.Contract.PocContractAddress(&_AddressBookV2.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AddressBookV2 *AddressBookV2Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AddressBookV2 *AddressBookV2Session) ProxiableUUID() ([32]byte, error) {
	return _AddressBookV2.Contract.ProxiableUUID(&_AddressBookV2.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_AddressBookV2 *AddressBookV2CallerSession) ProxiableUUID() ([32]byte, error) {
	return _AddressBookV2.Contract.ProxiableUUID(&_AddressBookV2.CallOpts)
}

// Requirement is a free data retrieval call binding the contract method 0xb7563930.
//
// Solidity: function requirement() pure returns(uint256)
func (_AddressBookV2 *AddressBookV2Caller) Requirement(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "requirement")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Requirement is a free data retrieval call binding the contract method 0xb7563930.
//
// Solidity: function requirement() pure returns(uint256)
func (_AddressBookV2 *AddressBookV2Session) Requirement() (*big.Int, error) {
	return _AddressBookV2.Contract.Requirement(&_AddressBookV2.CallOpts)
}

// Requirement is a free data retrieval call binding the contract method 0xb7563930.
//
// Solidity: function requirement() pure returns(uint256)
func (_AddressBookV2 *AddressBookV2CallerSession) Requirement() (*big.Int, error) {
	return _AddressBookV2.Contract.Requirement(&_AddressBookV2.CallOpts)
}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2Caller) SpareContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AddressBookV2.contract.Call(opts, &out, "spareContractAddress")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2Session) SpareContractAddress() (common.Address, error) {
	return _AddressBookV2.Contract.SpareContractAddress(&_AddressBookV2.CallOpts)
}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_AddressBookV2 *AddressBookV2CallerSession) SpareContractAddress() (common.Address, error) {
	return _AddressBookV2.Contract.SpareContractAddress(&_AddressBookV2.CallOpts)
}

// AssignGcId is a paid mutator transaction binding the contract method 0x06bb8471.
//
// Solidity: function assignGcId(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) AssignGcId(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "assignGcId", nodeId)
}

// AssignGcId is a paid mutator transaction binding the contract method 0x06bb8471.
//
// Solidity: function assignGcId(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) AssignGcId(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.AssignGcId(&_AddressBookV2.TransactOpts, nodeId)
}

// AssignGcId is a paid mutator transaction binding the contract method 0x06bb8471.
//
// Solidity: function assignGcId(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) AssignGcId(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.AssignGcId(&_AddressBookV2.TransactOpts, nodeId)
}

// CreateNode is a paid mutator transaction binding the contract method 0xe1789d81.
//
// Solidity: function createNode(address nodeId, address stakingContract, address rewardAddress, address voterAddress, (bytes,bytes) blsInfo, string name, string metadata) returns()
func (_AddressBookV2 *AddressBookV2Transactor) CreateNode(opts *bind.TransactOpts, nodeId common.Address, stakingContract common.Address, rewardAddress common.Address, voterAddress common.Address, blsInfo BlsPublicKeyInfo, name string, metadata string) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "createNode", nodeId, stakingContract, rewardAddress, voterAddress, blsInfo, name, metadata)
}

// CreateNode is a paid mutator transaction binding the contract method 0xe1789d81.
//
// Solidity: function createNode(address nodeId, address stakingContract, address rewardAddress, address voterAddress, (bytes,bytes) blsInfo, string name, string metadata) returns()
func (_AddressBookV2 *AddressBookV2Session) CreateNode(nodeId common.Address, stakingContract common.Address, rewardAddress common.Address, voterAddress common.Address, blsInfo BlsPublicKeyInfo, name string, metadata string) (*types.Transaction, error) {
	return _AddressBookV2.Contract.CreateNode(&_AddressBookV2.TransactOpts, nodeId, stakingContract, rewardAddress, voterAddress, blsInfo, name, metadata)
}

// CreateNode is a paid mutator transaction binding the contract method 0xe1789d81.
//
// Solidity: function createNode(address nodeId, address stakingContract, address rewardAddress, address voterAddress, (bytes,bytes) blsInfo, string name, string metadata) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) CreateNode(nodeId common.Address, stakingContract common.Address, rewardAddress common.Address, voterAddress common.Address, blsInfo BlsPublicKeyInfo, name string, metadata string) (*types.Transaction, error) {
	return _AddressBookV2.Contract.CreateNode(&_AddressBookV2.TransactOpts, nodeId, stakingContract, rewardAddress, voterAddress, blsInfo, name, metadata)
}

// DeleteNode is a paid mutator transaction binding the contract method 0x2d4ede93.
//
// Solidity: function deleteNode(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) DeleteNode(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "deleteNode", nodeId)
}

// DeleteNode is a paid mutator transaction binding the contract method 0x2d4ede93.
//
// Solidity: function deleteNode(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) DeleteNode(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.DeleteNode(&_AddressBookV2.TransactOpts, nodeId)
}

// DeleteNode is a paid mutator transaction binding the contract method 0x2d4ede93.
//
// Solidity: function deleteNode(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) DeleteNode(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.DeleteNode(&_AddressBookV2.TransactOpts, nodeId)
}

// Exit is a paid mutator transaction binding the contract method 0xb42652e9.
//
// Solidity: function exit(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) Exit(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "exit", nodeId)
}

// Exit is a paid mutator transaction binding the contract method 0xb42652e9.
//
// Solidity: function exit(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) Exit(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Exit(&_AddressBookV2.TransactOpts, nodeId)
}

// Exit is a paid mutator transaction binding the contract method 0xb42652e9.
//
// Solidity: function exit(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) Exit(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Exit(&_AddressBookV2.TransactOpts, nodeId)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AddressBookV2 *AddressBookV2Transactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AddressBookV2 *AddressBookV2Session) Initialize() (*types.Transaction, error) {
	return _AddressBookV2.Contract.Initialize(&_AddressBookV2.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) Initialize() (*types.Transaction, error) {
	return _AddressBookV2.Contract.Initialize(&_AddressBookV2.TransactOpts)
}

// Offboard is a paid mutator transaction binding the contract method 0xb9f96f40.
//
// Solidity: function offboard(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) Offboard(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "offboard", nodeId)
}

// Offboard is a paid mutator transaction binding the contract method 0xb9f96f40.
//
// Solidity: function offboard(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) Offboard(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Offboard(&_AddressBookV2.TransactOpts, nodeId)
}

// Offboard is a paid mutator transaction binding the contract method 0xb9f96f40.
//
// Solidity: function offboard(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) Offboard(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Offboard(&_AddressBookV2.TransactOpts, nodeId)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) Pause(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "pause", nodeId)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) Pause(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Pause(&_AddressBookV2.TransactOpts, nodeId)
}

// Pause is a paid mutator transaction binding the contract method 0x76a67a51.
//
// Solidity: function pause(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) Pause(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Pause(&_AddressBookV2.TransactOpts, nodeId)
}

// ProcessSystemTransition is a paid mutator transaction binding the contract method 0x1b8f34ca.
//
// Solidity: function processSystemTransition(address[] nodeIds, uint8[] newStates, uint256[] timeoutAts, uint256 epochVACount) returns()
func (_AddressBookV2 *AddressBookV2Transactor) ProcessSystemTransition(opts *bind.TransactOpts, nodeIds []common.Address, newStates []uint8, timeoutAts []*big.Int, epochVACount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "processSystemTransition", nodeIds, newStates, timeoutAts, epochVACount)
}

// ProcessSystemTransition is a paid mutator transaction binding the contract method 0x1b8f34ca.
//
// Solidity: function processSystemTransition(address[] nodeIds, uint8[] newStates, uint256[] timeoutAts, uint256 epochVACount) returns()
func (_AddressBookV2 *AddressBookV2Session) ProcessSystemTransition(nodeIds []common.Address, newStates []uint8, timeoutAts []*big.Int, epochVACount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.ProcessSystemTransition(&_AddressBookV2.TransactOpts, nodeIds, newStates, timeoutAts, epochVACount)
}

// ProcessSystemTransition is a paid mutator transaction binding the contract method 0x1b8f34ca.
//
// Solidity: function processSystemTransition(address[] nodeIds, uint8[] newStates, uint256[] timeoutAts, uint256 epochVACount) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) ProcessSystemTransition(nodeIds []common.Address, newStates []uint8, timeoutAts []*big.Int, epochVACount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.ProcessSystemTransition(&_AddressBookV2.TransactOpts, nodeIds, newStates, timeoutAts, epochVACount)
}

// ReadyCandidate is a paid mutator transaction binding the contract method 0x453e962e.
//
// Solidity: function readyCandidate(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) ReadyCandidate(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "readyCandidate", nodeId)
}

// ReadyCandidate is a paid mutator transaction binding the contract method 0x453e962e.
//
// Solidity: function readyCandidate(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) ReadyCandidate(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.ReadyCandidate(&_AddressBookV2.TransactOpts, nodeId)
}

// ReadyCandidate is a paid mutator transaction binding the contract method 0x453e962e.
//
// Solidity: function readyCandidate(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) ReadyCandidate(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.ReadyCandidate(&_AddressBookV2.TransactOpts, nodeId)
}

// ReadyValidator is a paid mutator transaction binding the contract method 0x656f5869.
//
// Solidity: function readyValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) ReadyValidator(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "readyValidator", nodeId)
}

// ReadyValidator is a paid mutator transaction binding the contract method 0x656f5869.
//
// Solidity: function readyValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) ReadyValidator(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.ReadyValidator(&_AddressBookV2.TransactOpts, nodeId)
}

// ReadyValidator is a paid mutator transaction binding the contract method 0x656f5869.
//
// Solidity: function readyValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) ReadyValidator(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.ReadyValidator(&_AddressBookV2.TransactOpts, nodeId)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AddressBookV2 *AddressBookV2Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AddressBookV2 *AddressBookV2Session) RenounceOwnership() (*types.Transaction, error) {
	return _AddressBookV2.Contract.RenounceOwnership(&_AddressBookV2.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AddressBookV2.Contract.RenounceOwnership(&_AddressBookV2.TransactOpts)
}

// Resume is a paid mutator transaction binding the contract method 0x793c1946.
//
// Solidity: function resume(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) Resume(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "resume", nodeId)
}

// Resume is a paid mutator transaction binding the contract method 0x793c1946.
//
// Solidity: function resume(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) Resume(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Resume(&_AddressBookV2.TransactOpts, nodeId)
}

// Resume is a paid mutator transaction binding the contract method 0x793c1946.
//
// Solidity: function resume(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) Resume(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Resume(&_AddressBookV2.TransactOpts, nodeId)
}

// RevokeGcId is a paid mutator transaction binding the contract method 0xbe535f8b.
//
// Solidity: function revokeGcId(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) RevokeGcId(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "revokeGcId", nodeId)
}

// RevokeGcId is a paid mutator transaction binding the contract method 0xbe535f8b.
//
// Solidity: function revokeGcId(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) RevokeGcId(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.RevokeGcId(&_AddressBookV2.TransactOpts, nodeId)
}

// RevokeGcId is a paid mutator transaction binding the contract method 0xbe535f8b.
//
// Solidity: function revokeGcId(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) RevokeGcId(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.RevokeGcId(&_AddressBookV2.TransactOpts, nodeId)
}

// SuspendValidator is a paid mutator transaction binding the contract method 0xa41b6000.
//
// Solidity: function suspendValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) SuspendValidator(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "suspendValidator", nodeId)
}

// SuspendValidator is a paid mutator transaction binding the contract method 0xa41b6000.
//
// Solidity: function suspendValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) SuspendValidator(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.SuspendValidator(&_AddressBookV2.TransactOpts, nodeId)
}

// SuspendValidator is a paid mutator transaction binding the contract method 0xa41b6000.
//
// Solidity: function suspendValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) SuspendValidator(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.SuspendValidator(&_AddressBookV2.TransactOpts, nodeId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AddressBookV2 *AddressBookV2Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AddressBookV2 *AddressBookV2Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.TransferOwnership(&_AddressBookV2.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.TransferOwnership(&_AddressBookV2.TransactOpts, newOwner)
}

// UnreadyCandidate is a paid mutator transaction binding the contract method 0xe4f0d37c.
//
// Solidity: function unreadyCandidate(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UnreadyCandidate(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "unreadyCandidate", nodeId)
}

// UnreadyCandidate is a paid mutator transaction binding the contract method 0xe4f0d37c.
//
// Solidity: function unreadyCandidate(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) UnreadyCandidate(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UnreadyCandidate(&_AddressBookV2.TransactOpts, nodeId)
}

// UnreadyCandidate is a paid mutator transaction binding the contract method 0xe4f0d37c.
//
// Solidity: function unreadyCandidate(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UnreadyCandidate(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UnreadyCandidate(&_AddressBookV2.TransactOpts, nodeId)
}

// UnreadyValidator is a paid mutator transaction binding the contract method 0xd9abb38b.
//
// Solidity: function unreadyValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UnreadyValidator(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "unreadyValidator", nodeId)
}

// UnreadyValidator is a paid mutator transaction binding the contract method 0xd9abb38b.
//
// Solidity: function unreadyValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) UnreadyValidator(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UnreadyValidator(&_AddressBookV2.TransactOpts, nodeId)
}

// UnreadyValidator is a paid mutator transaction binding the contract method 0xd9abb38b.
//
// Solidity: function unreadyValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UnreadyValidator(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UnreadyValidator(&_AddressBookV2.TransactOpts, nodeId)
}

// UnsuspendValidator is a paid mutator transaction binding the contract method 0x78b84a5c.
//
// Solidity: function unsuspendValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UnsuspendValidator(opts *bind.TransactOpts, nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "unsuspendValidator", nodeId)
}

// UnsuspendValidator is a paid mutator transaction binding the contract method 0x78b84a5c.
//
// Solidity: function unsuspendValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2Session) UnsuspendValidator(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UnsuspendValidator(&_AddressBookV2.TransactOpts, nodeId)
}

// UnsuspendValidator is a paid mutator transaction binding the contract method 0x78b84a5c.
//
// Solidity: function unsuspendValidator(address nodeId) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UnsuspendValidator(nodeId common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UnsuspendValidator(&_AddressBookV2.TransactOpts, nodeId)
}

// UpdateCfsThreshold is a paid mutator transaction binding the contract method 0xd18c07ab.
//
// Solidity: function updateCfsThreshold(uint256 newCfsThreshold) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateCfsThreshold(opts *bind.TransactOpts, newCfsThreshold *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateCfsThreshold", newCfsThreshold)
}

// UpdateCfsThreshold is a paid mutator transaction binding the contract method 0xd18c07ab.
//
// Solidity: function updateCfsThreshold(uint256 newCfsThreshold) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateCfsThreshold(newCfsThreshold *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateCfsThreshold(&_AddressBookV2.TransactOpts, newCfsThreshold)
}

// UpdateCfsThreshold is a paid mutator transaction binding the contract method 0xd18c07ab.
//
// Solidity: function updateCfsThreshold(uint256 newCfsThreshold) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateCfsThreshold(newCfsThreshold *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateCfsThreshold(&_AddressBookV2.TransactOpts, newCfsThreshold)
}

// UpdateConfigurator is a paid mutator transaction binding the contract method 0xb57873a5.
//
// Solidity: function updateConfigurator(address newConfigurator) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateConfigurator(opts *bind.TransactOpts, newConfigurator common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateConfigurator", newConfigurator)
}

// UpdateConfigurator is a paid mutator transaction binding the contract method 0xb57873a5.
//
// Solidity: function updateConfigurator(address newConfigurator) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateConfigurator(newConfigurator common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateConfigurator(&_AddressBookV2.TransactOpts, newConfigurator)
}

// UpdateConfigurator is a paid mutator transaction binding the contract method 0xb57873a5.
//
// Solidity: function updateConfigurator(address newConfigurator) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateConfigurator(newConfigurator common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateConfigurator(&_AddressBookV2.TransactOpts, newConfigurator)
}

// UpdateIdleTimeout is a paid mutator transaction binding the contract method 0xe59d7a84.
//
// Solidity: function updateIdleTimeout(uint256 newIdleTimeout) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateIdleTimeout(opts *bind.TransactOpts, newIdleTimeout *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateIdleTimeout", newIdleTimeout)
}

// UpdateIdleTimeout is a paid mutator transaction binding the contract method 0xe59d7a84.
//
// Solidity: function updateIdleTimeout(uint256 newIdleTimeout) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateIdleTimeout(newIdleTimeout *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateIdleTimeout(&_AddressBookV2.TransactOpts, newIdleTimeout)
}

// UpdateIdleTimeout is a paid mutator transaction binding the contract method 0xe59d7a84.
//
// Solidity: function updateIdleTimeout(uint256 newIdleTimeout) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateIdleTimeout(newIdleTimeout *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateIdleTimeout(&_AddressBookV2.TransactOpts, newIdleTimeout)
}

// UpdateKefAddress is a paid mutator transaction binding the contract method 0x9d8cf08f.
//
// Solidity: function updateKefAddress(address newKefAddress) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateKefAddress(opts *bind.TransactOpts, newKefAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateKefAddress", newKefAddress)
}

// UpdateKefAddress is a paid mutator transaction binding the contract method 0x9d8cf08f.
//
// Solidity: function updateKefAddress(address newKefAddress) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateKefAddress(newKefAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateKefAddress(&_AddressBookV2.TransactOpts, newKefAddress)
}

// UpdateKefAddress is a paid mutator transaction binding the contract method 0x9d8cf08f.
//
// Solidity: function updateKefAddress(address newKefAddress) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateKefAddress(newKefAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateKefAddress(&_AddressBookV2.TransactOpts, newKefAddress)
}

// UpdateKifAddress is a paid mutator transaction binding the contract method 0x7df40c62.
//
// Solidity: function updateKifAddress(address newKifAddress) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateKifAddress(opts *bind.TransactOpts, newKifAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateKifAddress", newKifAddress)
}

// UpdateKifAddress is a paid mutator transaction binding the contract method 0x7df40c62.
//
// Solidity: function updateKifAddress(address newKifAddress) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateKifAddress(newKifAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateKifAddress(&_AddressBookV2.TransactOpts, newKifAddress)
}

// UpdateKifAddress is a paid mutator transaction binding the contract method 0x7df40c62.
//
// Solidity: function updateKifAddress(address newKifAddress) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateKifAddress(newKifAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateKifAddress(&_AddressBookV2.TransactOpts, newKifAddress)
}

// UpdateKpfAddress is a paid mutator transaction binding the contract method 0xc9a86af2.
//
// Solidity: function updateKpfAddress(address newKpfAddress) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateKpfAddress(opts *bind.TransactOpts, newKpfAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateKpfAddress", newKpfAddress)
}

// UpdateKpfAddress is a paid mutator transaction binding the contract method 0xc9a86af2.
//
// Solidity: function updateKpfAddress(address newKpfAddress) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateKpfAddress(newKpfAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateKpfAddress(&_AddressBookV2.TransactOpts, newKpfAddress)
}

// UpdateKpfAddress is a paid mutator transaction binding the contract method 0xc9a86af2.
//
// Solidity: function updateKpfAddress(address newKpfAddress) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateKpfAddress(newKpfAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateKpfAddress(&_AddressBookV2.TransactOpts, newKpfAddress)
}

// UpdateManager is a paid mutator transaction binding the contract method 0x07ecec3e.
//
// Solidity: function updateManager(address nodeId, address newManager) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateManager(opts *bind.TransactOpts, nodeId common.Address, newManager common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateManager", nodeId, newManager)
}

// UpdateManager is a paid mutator transaction binding the contract method 0x07ecec3e.
//
// Solidity: function updateManager(address nodeId, address newManager) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateManager(nodeId common.Address, newManager common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateManager(&_AddressBookV2.TransactOpts, nodeId, newManager)
}

// UpdateManager is a paid mutator transaction binding the contract method 0x07ecec3e.
//
// Solidity: function updateManager(address nodeId, address newManager) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateManager(nodeId common.Address, newManager common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateManager(&_AddressBookV2.TransactOpts, nodeId, newManager)
}

// UpdateMaxCandReadyCount is a paid mutator transaction binding the contract method 0xa9ee5472.
//
// Solidity: function updateMaxCandReadyCount(uint256 newMaxCandReadyCount) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateMaxCandReadyCount(opts *bind.TransactOpts, newMaxCandReadyCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateMaxCandReadyCount", newMaxCandReadyCount)
}

// UpdateMaxCandReadyCount is a paid mutator transaction binding the contract method 0xa9ee5472.
//
// Solidity: function updateMaxCandReadyCount(uint256 newMaxCandReadyCount) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateMaxCandReadyCount(newMaxCandReadyCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateMaxCandReadyCount(&_AddressBookV2.TransactOpts, newMaxCandReadyCount)
}

// UpdateMaxCandReadyCount is a paid mutator transaction binding the contract method 0xa9ee5472.
//
// Solidity: function updateMaxCandReadyCount(uint256 newMaxCandReadyCount) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateMaxCandReadyCount(newMaxCandReadyCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateMaxCandReadyCount(&_AddressBookV2.TransactOpts, newMaxCandReadyCount)
}

// UpdateMaxNodeCount is a paid mutator transaction binding the contract method 0xa4c98ada.
//
// Solidity: function updateMaxNodeCount(uint256 newMaxNodeCount) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateMaxNodeCount(opts *bind.TransactOpts, newMaxNodeCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateMaxNodeCount", newMaxNodeCount)
}

// UpdateMaxNodeCount is a paid mutator transaction binding the contract method 0xa4c98ada.
//
// Solidity: function updateMaxNodeCount(uint256 newMaxNodeCount) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateMaxNodeCount(newMaxNodeCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateMaxNodeCount(&_AddressBookV2.TransactOpts, newMaxNodeCount)
}

// UpdateMaxNodeCount is a paid mutator transaction binding the contract method 0xa4c98ada.
//
// Solidity: function updateMaxNodeCount(uint256 newMaxNodeCount) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateMaxNodeCount(newMaxNodeCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateMaxNodeCount(&_AddressBookV2.TransactOpts, newMaxNodeCount)
}

// UpdateMaxValActivePausedCount is a paid mutator transaction binding the contract method 0x5b27b6c9.
//
// Solidity: function updateMaxValActivePausedCount(uint256 newMaxValActivePausedCount) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateMaxValActivePausedCount(opts *bind.TransactOpts, newMaxValActivePausedCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateMaxValActivePausedCount", newMaxValActivePausedCount)
}

// UpdateMaxValActivePausedCount is a paid mutator transaction binding the contract method 0x5b27b6c9.
//
// Solidity: function updateMaxValActivePausedCount(uint256 newMaxValActivePausedCount) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateMaxValActivePausedCount(newMaxValActivePausedCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateMaxValActivePausedCount(&_AddressBookV2.TransactOpts, newMaxValActivePausedCount)
}

// UpdateMaxValActivePausedCount is a paid mutator transaction binding the contract method 0x5b27b6c9.
//
// Solidity: function updateMaxValActivePausedCount(uint256 newMaxValActivePausedCount) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateMaxValActivePausedCount(newMaxValActivePausedCount *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateMaxValActivePausedCount(&_AddressBookV2.TransactOpts, newMaxValActivePausedCount)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0xda38d498.
//
// Solidity: function updateMetadata(address nodeId, string newMetadata) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateMetadata(opts *bind.TransactOpts, nodeId common.Address, newMetadata string) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateMetadata", nodeId, newMetadata)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0xda38d498.
//
// Solidity: function updateMetadata(address nodeId, string newMetadata) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateMetadata(nodeId common.Address, newMetadata string) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateMetadata(&_AddressBookV2.TransactOpts, nodeId, newMetadata)
}

// UpdateMetadata is a paid mutator transaction binding the contract method 0xda38d498.
//
// Solidity: function updateMetadata(address nodeId, string newMetadata) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateMetadata(nodeId common.Address, newMetadata string) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateMetadata(&_AddressBookV2.TransactOpts, nodeId, newMetadata)
}

// UpdatePauseTimeout is a paid mutator transaction binding the contract method 0x9d0e234d.
//
// Solidity: function updatePauseTimeout(uint256 newPauseTimeout) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdatePauseTimeout(opts *bind.TransactOpts, newPauseTimeout *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updatePauseTimeout", newPauseTimeout)
}

// UpdatePauseTimeout is a paid mutator transaction binding the contract method 0x9d0e234d.
//
// Solidity: function updatePauseTimeout(uint256 newPauseTimeout) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdatePauseTimeout(newPauseTimeout *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdatePauseTimeout(&_AddressBookV2.TransactOpts, newPauseTimeout)
}

// UpdatePauseTimeout is a paid mutator transaction binding the contract method 0x9d0e234d.
//
// Solidity: function updatePauseTimeout(uint256 newPauseTimeout) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdatePauseTimeout(newPauseTimeout *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdatePauseTimeout(&_AddressBookV2.TransactOpts, newPauseTimeout)
}

// UpdatePfsThreshold is a paid mutator transaction binding the contract method 0xba70d018.
//
// Solidity: function updatePfsThreshold(uint256 newPfsThreshold) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdatePfsThreshold(opts *bind.TransactOpts, newPfsThreshold *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updatePfsThreshold", newPfsThreshold)
}

// UpdatePfsThreshold is a paid mutator transaction binding the contract method 0xba70d018.
//
// Solidity: function updatePfsThreshold(uint256 newPfsThreshold) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdatePfsThreshold(newPfsThreshold *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdatePfsThreshold(&_AddressBookV2.TransactOpts, newPfsThreshold)
}

// UpdatePfsThreshold is a paid mutator transaction binding the contract method 0xba70d018.
//
// Solidity: function updatePfsThreshold(uint256 newPfsThreshold) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdatePfsThreshold(newPfsThreshold *big.Int) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdatePfsThreshold(&_AddressBookV2.TransactOpts, newPfsThreshold)
}

// UpdateRewardAddress is a paid mutator transaction binding the contract method 0x394f8899.
//
// Solidity: function updateRewardAddress(address nodeId, address newRewardAddress) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateRewardAddress(opts *bind.TransactOpts, nodeId common.Address, newRewardAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateRewardAddress", nodeId, newRewardAddress)
}

// UpdateRewardAddress is a paid mutator transaction binding the contract method 0x394f8899.
//
// Solidity: function updateRewardAddress(address nodeId, address newRewardAddress) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateRewardAddress(nodeId common.Address, newRewardAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateRewardAddress(&_AddressBookV2.TransactOpts, nodeId, newRewardAddress)
}

// UpdateRewardAddress is a paid mutator transaction binding the contract method 0x394f8899.
//
// Solidity: function updateRewardAddress(address nodeId, address newRewardAddress) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateRewardAddress(nodeId common.Address, newRewardAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateRewardAddress(&_AddressBookV2.TransactOpts, nodeId, newRewardAddress)
}

// UpdateSuspender is a paid mutator transaction binding the contract method 0x50de2fb3.
//
// Solidity: function updateSuspender(address newSuspender) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateSuspender(opts *bind.TransactOpts, newSuspender common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateSuspender", newSuspender)
}

// UpdateSuspender is a paid mutator transaction binding the contract method 0x50de2fb3.
//
// Solidity: function updateSuspender(address newSuspender) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateSuspender(newSuspender common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateSuspender(&_AddressBookV2.TransactOpts, newSuspender)
}

// UpdateSuspender is a paid mutator transaction binding the contract method 0x50de2fb3.
//
// Solidity: function updateSuspender(address newSuspender) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateSuspender(newSuspender common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateSuspender(&_AddressBookV2.TransactOpts, newSuspender)
}

// UpdateVoterAddress is a paid mutator transaction binding the contract method 0x9f9e3cba.
//
// Solidity: function updateVoterAddress(address nodeId, address newVoterAddress) returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpdateVoterAddress(opts *bind.TransactOpts, nodeId common.Address, newVoterAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "updateVoterAddress", nodeId, newVoterAddress)
}

// UpdateVoterAddress is a paid mutator transaction binding the contract method 0x9f9e3cba.
//
// Solidity: function updateVoterAddress(address nodeId, address newVoterAddress) returns()
func (_AddressBookV2 *AddressBookV2Session) UpdateVoterAddress(nodeId common.Address, newVoterAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateVoterAddress(&_AddressBookV2.TransactOpts, nodeId, newVoterAddress)
}

// UpdateVoterAddress is a paid mutator transaction binding the contract method 0x9f9e3cba.
//
// Solidity: function updateVoterAddress(address nodeId, address newVoterAddress) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpdateVoterAddress(nodeId common.Address, newVoterAddress common.Address) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpdateVoterAddress(&_AddressBookV2.TransactOpts, nodeId, newVoterAddress)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AddressBookV2 *AddressBookV2Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AddressBookV2 *AddressBookV2Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpgradeToAndCall(&_AddressBookV2.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _AddressBookV2.Contract.UpgradeToAndCall(&_AddressBookV2.TransactOpts, newImplementation, data)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_AddressBookV2 *AddressBookV2Transactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _AddressBookV2.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_AddressBookV2 *AddressBookV2Session) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Fallback(&_AddressBookV2.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _AddressBookV2.Contract.Fallback(&_AddressBookV2.TransactOpts, calldata)
}

// AddressBookV2AddressConfigUpdatedIterator is returned from FilterAddressConfigUpdated and is used to iterate over the raw logs and unpacked data for AddressConfigUpdated events raised by the AddressBookV2 contract.
type AddressBookV2AddressConfigUpdatedIterator struct {
	Event *AddressBookV2AddressConfigUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2AddressConfigUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2AddressConfigUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2AddressConfigUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2AddressConfigUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2AddressConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2AddressConfigUpdated represents a AddressConfigUpdated event raised by the AddressBookV2 contract.
type AddressBookV2AddressConfigUpdated struct {
	ConfigId uint8
	OldValue common.Address
	NewValue common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAddressConfigUpdated is a free log retrieval operation binding the contract event 0xd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4b.
//
// Solidity: event AddressConfigUpdated(uint8 indexed configId, address oldValue, address newValue)
func (_AddressBookV2 *AddressBookV2Filterer) FilterAddressConfigUpdated(opts *bind.FilterOpts, configId []uint8) (*AddressBookV2AddressConfigUpdatedIterator, error) {
	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "AddressConfigUpdated", configIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2AddressConfigUpdatedIterator{contract: _AddressBookV2.contract, event: "AddressConfigUpdated", logs: logs, sub: sub}, nil
}

// WatchAddressConfigUpdated is a free log subscription operation binding the contract event 0xd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4b.
//
// Solidity: event AddressConfigUpdated(uint8 indexed configId, address oldValue, address newValue)
func (_AddressBookV2 *AddressBookV2Filterer) WatchAddressConfigUpdated(opts *bind.WatchOpts, sink chan<- *AddressBookV2AddressConfigUpdated, configId []uint8) (event.Subscription, error) {
	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "AddressConfigUpdated", configIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2AddressConfigUpdated)
				if err := _AddressBookV2.contract.UnpackLog(event, "AddressConfigUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAddressConfigUpdated is a log parse operation binding the contract event 0xd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4b.
//
// Solidity: event AddressConfigUpdated(uint8 indexed configId, address oldValue, address newValue)
func (_AddressBookV2 *AddressBookV2Filterer) ParseAddressConfigUpdated(log types.Log) (*AddressBookV2AddressConfigUpdated, error) {
	event := new(AddressBookV2AddressConfigUpdated)
	if err := _AddressBookV2.contract.UnpackLog(event, "AddressConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2CandidateReadiedIterator is returned from FilterCandidateReadied and is used to iterate over the raw logs and unpacked data for CandidateReadied events raised by the AddressBookV2 contract.
type AddressBookV2CandidateReadiedIterator struct {
	Event *AddressBookV2CandidateReadied // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2CandidateReadiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2CandidateReadied)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2CandidateReadied)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2CandidateReadiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2CandidateReadiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2CandidateReadied represents a CandidateReadied event raised by the AddressBookV2 contract.
type AddressBookV2CandidateReadied struct {
	NodeId common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCandidateReadied is a free log retrieval operation binding the contract event 0xb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98.
//
// Solidity: event CandidateReadied(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) FilterCandidateReadied(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2CandidateReadiedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "CandidateReadied", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2CandidateReadiedIterator{contract: _AddressBookV2.contract, event: "CandidateReadied", logs: logs, sub: sub}, nil
}

// WatchCandidateReadied is a free log subscription operation binding the contract event 0xb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98.
//
// Solidity: event CandidateReadied(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) WatchCandidateReadied(opts *bind.WatchOpts, sink chan<- *AddressBookV2CandidateReadied, nodeId []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "CandidateReadied", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2CandidateReadied)
				if err := _AddressBookV2.contract.UnpackLog(event, "CandidateReadied", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCandidateReadied is a log parse operation binding the contract event 0xb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98.
//
// Solidity: event CandidateReadied(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) ParseCandidateReadied(log types.Log) (*AddressBookV2CandidateReadied, error) {
	event := new(AddressBookV2CandidateReadied)
	if err := _AddressBookV2.contract.UnpackLog(event, "CandidateReadied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2CandidateUnreadiedIterator is returned from FilterCandidateUnreadied and is used to iterate over the raw logs and unpacked data for CandidateUnreadied events raised by the AddressBookV2 contract.
type AddressBookV2CandidateUnreadiedIterator struct {
	Event *AddressBookV2CandidateUnreadied // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2CandidateUnreadiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2CandidateUnreadied)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2CandidateUnreadied)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2CandidateUnreadiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2CandidateUnreadiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2CandidateUnreadied represents a CandidateUnreadied event raised by the AddressBookV2 contract.
type AddressBookV2CandidateUnreadied struct {
	NodeId common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCandidateUnreadied is a free log retrieval operation binding the contract event 0x8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2.
//
// Solidity: event CandidateUnreadied(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) FilterCandidateUnreadied(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2CandidateUnreadiedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "CandidateUnreadied", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2CandidateUnreadiedIterator{contract: _AddressBookV2.contract, event: "CandidateUnreadied", logs: logs, sub: sub}, nil
}

// WatchCandidateUnreadied is a free log subscription operation binding the contract event 0x8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2.
//
// Solidity: event CandidateUnreadied(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) WatchCandidateUnreadied(opts *bind.WatchOpts, sink chan<- *AddressBookV2CandidateUnreadied, nodeId []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "CandidateUnreadied", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2CandidateUnreadied)
				if err := _AddressBookV2.contract.UnpackLog(event, "CandidateUnreadied", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCandidateUnreadied is a log parse operation binding the contract event 0x8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2.
//
// Solidity: event CandidateUnreadied(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) ParseCandidateUnreadied(log types.Log) (*AddressBookV2CandidateUnreadied, error) {
	event := new(AddressBookV2CandidateUnreadied)
	if err := _AddressBookV2.contract.UnpackLog(event, "CandidateUnreadied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2EpochTransitionProcessedIterator is returned from FilterEpochTransitionProcessed and is used to iterate over the raw logs and unpacked data for EpochTransitionProcessed events raised by the AddressBookV2 contract.
type AddressBookV2EpochTransitionProcessedIterator struct {
	Event *AddressBookV2EpochTransitionProcessed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2EpochTransitionProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2EpochTransitionProcessed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2EpochTransitionProcessed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2EpochTransitionProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2EpochTransitionProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2EpochTransitionProcessed represents a EpochTransitionProcessed event raised by the AddressBookV2 contract.
type AddressBookV2EpochTransitionProcessed struct {
	EpochVACount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterEpochTransitionProcessed is a free log retrieval operation binding the contract event 0xd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e84816052.
//
// Solidity: event EpochTransitionProcessed(uint256 epochVACount)
func (_AddressBookV2 *AddressBookV2Filterer) FilterEpochTransitionProcessed(opts *bind.FilterOpts) (*AddressBookV2EpochTransitionProcessedIterator, error) {
	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "EpochTransitionProcessed")
	if err != nil {
		return nil, err
	}
	return &AddressBookV2EpochTransitionProcessedIterator{contract: _AddressBookV2.contract, event: "EpochTransitionProcessed", logs: logs, sub: sub}, nil
}

// WatchEpochTransitionProcessed is a free log subscription operation binding the contract event 0xd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e84816052.
//
// Solidity: event EpochTransitionProcessed(uint256 epochVACount)
func (_AddressBookV2 *AddressBookV2Filterer) WatchEpochTransitionProcessed(opts *bind.WatchOpts, sink chan<- *AddressBookV2EpochTransitionProcessed) (event.Subscription, error) {
	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "EpochTransitionProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2EpochTransitionProcessed)
				if err := _AddressBookV2.contract.UnpackLog(event, "EpochTransitionProcessed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseEpochTransitionProcessed is a log parse operation binding the contract event 0xd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e84816052.
//
// Solidity: event EpochTransitionProcessed(uint256 epochVACount)
func (_AddressBookV2 *AddressBookV2Filterer) ParseEpochTransitionProcessed(log types.Log) (*AddressBookV2EpochTransitionProcessed, error) {
	event := new(AddressBookV2EpochTransitionProcessed)
	if err := _AddressBookV2.contract.UnpackLog(event, "EpochTransitionProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2GcIdAssignedIterator is returned from FilterGcIdAssigned and is used to iterate over the raw logs and unpacked data for GcIdAssigned events raised by the AddressBookV2 contract.
type AddressBookV2GcIdAssignedIterator struct {
	Event *AddressBookV2GcIdAssigned // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2GcIdAssignedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2GcIdAssigned)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2GcIdAssigned)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2GcIdAssignedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2GcIdAssignedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2GcIdAssigned represents a GcIdAssigned event raised by the AddressBookV2 contract.
type AddressBookV2GcIdAssigned struct {
	NodeId common.Address
	GcId   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterGcIdAssigned is a free log retrieval operation binding the contract event 0xe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb.
//
// Solidity: event GcIdAssigned(address indexed nodeId, uint256 gcId)
func (_AddressBookV2 *AddressBookV2Filterer) FilterGcIdAssigned(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2GcIdAssignedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "GcIdAssigned", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2GcIdAssignedIterator{contract: _AddressBookV2.contract, event: "GcIdAssigned", logs: logs, sub: sub}, nil
}

// WatchGcIdAssigned is a free log subscription operation binding the contract event 0xe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb.
//
// Solidity: event GcIdAssigned(address indexed nodeId, uint256 gcId)
func (_AddressBookV2 *AddressBookV2Filterer) WatchGcIdAssigned(opts *bind.WatchOpts, sink chan<- *AddressBookV2GcIdAssigned, nodeId []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "GcIdAssigned", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2GcIdAssigned)
				if err := _AddressBookV2.contract.UnpackLog(event, "GcIdAssigned", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGcIdAssigned is a log parse operation binding the contract event 0xe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb.
//
// Solidity: event GcIdAssigned(address indexed nodeId, uint256 gcId)
func (_AddressBookV2 *AddressBookV2Filterer) ParseGcIdAssigned(log types.Log) (*AddressBookV2GcIdAssigned, error) {
	event := new(AddressBookV2GcIdAssigned)
	if err := _AddressBookV2.contract.UnpackLog(event, "GcIdAssigned", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AddressBookV2 contract.
type AddressBookV2InitializedIterator struct {
	Event *AddressBookV2Initialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2Initialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2Initialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2Initialized represents a Initialized event raised by the AddressBookV2 contract.
type AddressBookV2Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AddressBookV2 *AddressBookV2Filterer) FilterInitialized(opts *bind.FilterOpts) (*AddressBookV2InitializedIterator, error) {
	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AddressBookV2InitializedIterator{contract: _AddressBookV2.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AddressBookV2 *AddressBookV2Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AddressBookV2Initialized) (event.Subscription, error) {
	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2Initialized)
				if err := _AddressBookV2.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AddressBookV2 *AddressBookV2Filterer) ParseInitialized(log types.Log) (*AddressBookV2Initialized, error) {
	event := new(AddressBookV2Initialized)
	if err := _AddressBookV2.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2ManagerUpdatedIterator is returned from FilterManagerUpdated and is used to iterate over the raw logs and unpacked data for ManagerUpdated events raised by the AddressBookV2 contract.
type AddressBookV2ManagerUpdatedIterator struct {
	Event *AddressBookV2ManagerUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2ManagerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2ManagerUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2ManagerUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2ManagerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2ManagerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2ManagerUpdated represents a ManagerUpdated event raised by the AddressBookV2 contract.
type AddressBookV2ManagerUpdated struct {
	NodeId     common.Address
	OldManager common.Address
	NewManager common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterManagerUpdated is a free log retrieval operation binding the contract event 0x8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda87.
//
// Solidity: event ManagerUpdated(address indexed nodeId, address indexed oldManager, address indexed newManager)
func (_AddressBookV2 *AddressBookV2Filterer) FilterManagerUpdated(opts *bind.FilterOpts, nodeId []common.Address, oldManager []common.Address, newManager []common.Address) (*AddressBookV2ManagerUpdatedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var oldManagerRule []interface{}
	for _, oldManagerItem := range oldManager {
		oldManagerRule = append(oldManagerRule, oldManagerItem)
	}
	var newManagerRule []interface{}
	for _, newManagerItem := range newManager {
		newManagerRule = append(newManagerRule, newManagerItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "ManagerUpdated", nodeIdRule, oldManagerRule, newManagerRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2ManagerUpdatedIterator{contract: _AddressBookV2.contract, event: "ManagerUpdated", logs: logs, sub: sub}, nil
}

// WatchManagerUpdated is a free log subscription operation binding the contract event 0x8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda87.
//
// Solidity: event ManagerUpdated(address indexed nodeId, address indexed oldManager, address indexed newManager)
func (_AddressBookV2 *AddressBookV2Filterer) WatchManagerUpdated(opts *bind.WatchOpts, sink chan<- *AddressBookV2ManagerUpdated, nodeId []common.Address, oldManager []common.Address, newManager []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var oldManagerRule []interface{}
	for _, oldManagerItem := range oldManager {
		oldManagerRule = append(oldManagerRule, oldManagerItem)
	}
	var newManagerRule []interface{}
	for _, newManagerItem := range newManager {
		newManagerRule = append(newManagerRule, newManagerItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "ManagerUpdated", nodeIdRule, oldManagerRule, newManagerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2ManagerUpdated)
				if err := _AddressBookV2.contract.UnpackLog(event, "ManagerUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseManagerUpdated is a log parse operation binding the contract event 0x8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda87.
//
// Solidity: event ManagerUpdated(address indexed nodeId, address indexed oldManager, address indexed newManager)
func (_AddressBookV2 *AddressBookV2Filterer) ParseManagerUpdated(log types.Log) (*AddressBookV2ManagerUpdated, error) {
	event := new(AddressBookV2ManagerUpdated)
	if err := _AddressBookV2.contract.UnpackLog(event, "ManagerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2MetadataUpdatedIterator is returned from FilterMetadataUpdated and is used to iterate over the raw logs and unpacked data for MetadataUpdated events raised by the AddressBookV2 contract.
type AddressBookV2MetadataUpdatedIterator struct {
	Event *AddressBookV2MetadataUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2MetadataUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2MetadataUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2MetadataUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2MetadataUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2MetadataUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2MetadataUpdated represents a MetadataUpdated event raised by the AddressBookV2 contract.
type AddressBookV2MetadataUpdated struct {
	NodeId      common.Address
	NewMetadata string
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdated is a free log retrieval operation binding the contract event 0x2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf.
//
// Solidity: event MetadataUpdated(address indexed nodeId, string newMetadata)
func (_AddressBookV2 *AddressBookV2Filterer) FilterMetadataUpdated(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2MetadataUpdatedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "MetadataUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2MetadataUpdatedIterator{contract: _AddressBookV2.contract, event: "MetadataUpdated", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdated is a free log subscription operation binding the contract event 0x2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf.
//
// Solidity: event MetadataUpdated(address indexed nodeId, string newMetadata)
func (_AddressBookV2 *AddressBookV2Filterer) WatchMetadataUpdated(opts *bind.WatchOpts, sink chan<- *AddressBookV2MetadataUpdated, nodeId []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "MetadataUpdated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2MetadataUpdated)
				if err := _AddressBookV2.contract.UnpackLog(event, "MetadataUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMetadataUpdated is a log parse operation binding the contract event 0x2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf.
//
// Solidity: event MetadataUpdated(address indexed nodeId, string newMetadata)
func (_AddressBookV2 *AddressBookV2Filterer) ParseMetadataUpdated(log types.Log) (*AddressBookV2MetadataUpdated, error) {
	event := new(AddressBookV2MetadataUpdated)
	if err := _AddressBookV2.contract.UnpackLog(event, "MetadataUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2NodeCreatedIterator is returned from FilterNodeCreated and is used to iterate over the raw logs and unpacked data for NodeCreated events raised by the AddressBookV2 contract.
type AddressBookV2NodeCreatedIterator struct {
	Event *AddressBookV2NodeCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2NodeCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2NodeCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2NodeCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2NodeCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2NodeCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2NodeCreated represents a NodeCreated event raised by the AddressBookV2 contract.
type AddressBookV2NodeCreated struct {
	NodeId common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeCreated is a free log retrieval operation binding the contract event 0x55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792.
//
// Solidity: event NodeCreated(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) FilterNodeCreated(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2NodeCreatedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "NodeCreated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2NodeCreatedIterator{contract: _AddressBookV2.contract, event: "NodeCreated", logs: logs, sub: sub}, nil
}

// WatchNodeCreated is a free log subscription operation binding the contract event 0x55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792.
//
// Solidity: event NodeCreated(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) WatchNodeCreated(opts *bind.WatchOpts, sink chan<- *AddressBookV2NodeCreated, nodeId []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "NodeCreated", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2NodeCreated)
				if err := _AddressBookV2.contract.UnpackLog(event, "NodeCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNodeCreated is a log parse operation binding the contract event 0x55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792.
//
// Solidity: event NodeCreated(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) ParseNodeCreated(log types.Log) (*AddressBookV2NodeCreated, error) {
	event := new(AddressBookV2NodeCreated)
	if err := _AddressBookV2.contract.UnpackLog(event, "NodeCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2NodeDeletedIterator is returned from FilterNodeDeleted and is used to iterate over the raw logs and unpacked data for NodeDeleted events raised by the AddressBookV2 contract.
type AddressBookV2NodeDeletedIterator struct {
	Event *AddressBookV2NodeDeleted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2NodeDeletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2NodeDeleted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2NodeDeleted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2NodeDeletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2NodeDeletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2NodeDeleted represents a NodeDeleted event raised by the AddressBookV2 contract.
type AddressBookV2NodeDeleted struct {
	NodeId common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNodeDeleted is a free log retrieval operation binding the contract event 0x1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad.
//
// Solidity: event NodeDeleted(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) FilterNodeDeleted(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2NodeDeletedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "NodeDeleted", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2NodeDeletedIterator{contract: _AddressBookV2.contract, event: "NodeDeleted", logs: logs, sub: sub}, nil
}

// WatchNodeDeleted is a free log subscription operation binding the contract event 0x1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad.
//
// Solidity: event NodeDeleted(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) WatchNodeDeleted(opts *bind.WatchOpts, sink chan<- *AddressBookV2NodeDeleted, nodeId []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "NodeDeleted", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2NodeDeleted)
				if err := _AddressBookV2.contract.UnpackLog(event, "NodeDeleted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNodeDeleted is a log parse operation binding the contract event 0x1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad.
//
// Solidity: event NodeDeleted(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) ParseNodeDeleted(log types.Log) (*AddressBookV2NodeDeleted, error) {
	event := new(AddressBookV2NodeDeleted)
	if err := _AddressBookV2.contract.UnpackLog(event, "NodeDeleted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AddressBookV2 contract.
type AddressBookV2OwnershipTransferredIterator struct {
	Event *AddressBookV2OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2OwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2OwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2OwnershipTransferred represents a OwnershipTransferred event raised by the AddressBookV2 contract.
type AddressBookV2OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AddressBookV2 *AddressBookV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AddressBookV2OwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2OwnershipTransferredIterator{contract: _AddressBookV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AddressBookV2 *AddressBookV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AddressBookV2OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2OwnershipTransferred)
				if err := _AddressBookV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AddressBookV2 *AddressBookV2Filterer) ParseOwnershipTransferred(log types.Log) (*AddressBookV2OwnershipTransferred, error) {
	event := new(AddressBookV2OwnershipTransferred)
	if err := _AddressBookV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2RewardAddressUpdatedIterator is returned from FilterRewardAddressUpdated and is used to iterate over the raw logs and unpacked data for RewardAddressUpdated events raised by the AddressBookV2 contract.
type AddressBookV2RewardAddressUpdatedIterator struct {
	Event *AddressBookV2RewardAddressUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2RewardAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2RewardAddressUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2RewardAddressUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2RewardAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2RewardAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2RewardAddressUpdated represents a RewardAddressUpdated event raised by the AddressBookV2 contract.
type AddressBookV2RewardAddressUpdated struct {
	NodeId           common.Address
	OldRewardAddress common.Address
	NewRewardAddress common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterRewardAddressUpdated is a free log retrieval operation binding the contract event 0x270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc69701.
//
// Solidity: event RewardAddressUpdated(address indexed nodeId, address indexed oldRewardAddress, address indexed newRewardAddress)
func (_AddressBookV2 *AddressBookV2Filterer) FilterRewardAddressUpdated(opts *bind.FilterOpts, nodeId []common.Address, oldRewardAddress []common.Address, newRewardAddress []common.Address) (*AddressBookV2RewardAddressUpdatedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var oldRewardAddressRule []interface{}
	for _, oldRewardAddressItem := range oldRewardAddress {
		oldRewardAddressRule = append(oldRewardAddressRule, oldRewardAddressItem)
	}
	var newRewardAddressRule []interface{}
	for _, newRewardAddressItem := range newRewardAddress {
		newRewardAddressRule = append(newRewardAddressRule, newRewardAddressItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "RewardAddressUpdated", nodeIdRule, oldRewardAddressRule, newRewardAddressRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2RewardAddressUpdatedIterator{contract: _AddressBookV2.contract, event: "RewardAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchRewardAddressUpdated is a free log subscription operation binding the contract event 0x270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc69701.
//
// Solidity: event RewardAddressUpdated(address indexed nodeId, address indexed oldRewardAddress, address indexed newRewardAddress)
func (_AddressBookV2 *AddressBookV2Filterer) WatchRewardAddressUpdated(opts *bind.WatchOpts, sink chan<- *AddressBookV2RewardAddressUpdated, nodeId []common.Address, oldRewardAddress []common.Address, newRewardAddress []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var oldRewardAddressRule []interface{}
	for _, oldRewardAddressItem := range oldRewardAddress {
		oldRewardAddressRule = append(oldRewardAddressRule, oldRewardAddressItem)
	}
	var newRewardAddressRule []interface{}
	for _, newRewardAddressItem := range newRewardAddress {
		newRewardAddressRule = append(newRewardAddressRule, newRewardAddressItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "RewardAddressUpdated", nodeIdRule, oldRewardAddressRule, newRewardAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2RewardAddressUpdated)
				if err := _AddressBookV2.contract.UnpackLog(event, "RewardAddressUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRewardAddressUpdated is a log parse operation binding the contract event 0x270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc69701.
//
// Solidity: event RewardAddressUpdated(address indexed nodeId, address indexed oldRewardAddress, address indexed newRewardAddress)
func (_AddressBookV2 *AddressBookV2Filterer) ParseRewardAddressUpdated(log types.Log) (*AddressBookV2RewardAddressUpdated, error) {
	event := new(AddressBookV2RewardAddressUpdated)
	if err := _AddressBookV2.contract.UnpackLog(event, "RewardAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2StateChangedIterator is returned from FilterStateChanged and is used to iterate over the raw logs and unpacked data for StateChanged events raised by the AddressBookV2 contract.
type AddressBookV2StateChangedIterator struct {
	Event *AddressBookV2StateChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2StateChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2StateChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2StateChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2StateChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2StateChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2StateChanged represents a StateChanged event raised by the AddressBookV2 contract.
type AddressBookV2StateChanged struct {
	NodeId    common.Address
	FromState uint8
	ToState   uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStateChanged is a free log retrieval operation binding the contract event 0xcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee.
//
// Solidity: event StateChanged(address indexed nodeId, uint8 indexed fromState, uint8 indexed toState)
func (_AddressBookV2 *AddressBookV2Filterer) FilterStateChanged(opts *bind.FilterOpts, nodeId []common.Address, fromState []uint8, toState []uint8) (*AddressBookV2StateChangedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var fromStateRule []interface{}
	for _, fromStateItem := range fromState {
		fromStateRule = append(fromStateRule, fromStateItem)
	}
	var toStateRule []interface{}
	for _, toStateItem := range toState {
		toStateRule = append(toStateRule, toStateItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "StateChanged", nodeIdRule, fromStateRule, toStateRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2StateChangedIterator{contract: _AddressBookV2.contract, event: "StateChanged", logs: logs, sub: sub}, nil
}

// WatchStateChanged is a free log subscription operation binding the contract event 0xcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee.
//
// Solidity: event StateChanged(address indexed nodeId, uint8 indexed fromState, uint8 indexed toState)
func (_AddressBookV2 *AddressBookV2Filterer) WatchStateChanged(opts *bind.WatchOpts, sink chan<- *AddressBookV2StateChanged, nodeId []common.Address, fromState []uint8, toState []uint8) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var fromStateRule []interface{}
	for _, fromStateItem := range fromState {
		fromStateRule = append(fromStateRule, fromStateItem)
	}
	var toStateRule []interface{}
	for _, toStateItem := range toState {
		toStateRule = append(toStateRule, toStateItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "StateChanged", nodeIdRule, fromStateRule, toStateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2StateChanged)
				if err := _AddressBookV2.contract.UnpackLog(event, "StateChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseStateChanged is a log parse operation binding the contract event 0xcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee.
//
// Solidity: event StateChanged(address indexed nodeId, uint8 indexed fromState, uint8 indexed toState)
func (_AddressBookV2 *AddressBookV2Filterer) ParseStateChanged(log types.Log) (*AddressBookV2StateChanged, error) {
	event := new(AddressBookV2StateChanged)
	if err := _AddressBookV2.contract.UnpackLog(event, "StateChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2SystemTransitionProcessedIterator is returned from FilterSystemTransitionProcessed and is used to iterate over the raw logs and unpacked data for SystemTransitionProcessed events raised by the AddressBookV2 contract.
type AddressBookV2SystemTransitionProcessedIterator struct {
	Event *AddressBookV2SystemTransitionProcessed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2SystemTransitionProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2SystemTransitionProcessed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2SystemTransitionProcessed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2SystemTransitionProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2SystemTransitionProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2SystemTransitionProcessed represents a SystemTransitionProcessed event raised by the AddressBookV2 contract.
type AddressBookV2SystemTransitionProcessed struct {
	NodeIds   []common.Address
	NewStates []uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSystemTransitionProcessed is a free log retrieval operation binding the contract event 0xab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e.
//
// Solidity: event SystemTransitionProcessed(address[] nodeIds, uint8[] newStates)
func (_AddressBookV2 *AddressBookV2Filterer) FilterSystemTransitionProcessed(opts *bind.FilterOpts) (*AddressBookV2SystemTransitionProcessedIterator, error) {
	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "SystemTransitionProcessed")
	if err != nil {
		return nil, err
	}
	return &AddressBookV2SystemTransitionProcessedIterator{contract: _AddressBookV2.contract, event: "SystemTransitionProcessed", logs: logs, sub: sub}, nil
}

// WatchSystemTransitionProcessed is a free log subscription operation binding the contract event 0xab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e.
//
// Solidity: event SystemTransitionProcessed(address[] nodeIds, uint8[] newStates)
func (_AddressBookV2 *AddressBookV2Filterer) WatchSystemTransitionProcessed(opts *bind.WatchOpts, sink chan<- *AddressBookV2SystemTransitionProcessed) (event.Subscription, error) {
	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "SystemTransitionProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2SystemTransitionProcessed)
				if err := _AddressBookV2.contract.UnpackLog(event, "SystemTransitionProcessed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSystemTransitionProcessed is a log parse operation binding the contract event 0xab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e.
//
// Solidity: event SystemTransitionProcessed(address[] nodeIds, uint8[] newStates)
func (_AddressBookV2 *AddressBookV2Filterer) ParseSystemTransitionProcessed(log types.Log) (*AddressBookV2SystemTransitionProcessed, error) {
	event := new(AddressBookV2SystemTransitionProcessed)
	if err := _AddressBookV2.contract.UnpackLog(event, "SystemTransitionProcessed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2UintConfigUpdatedIterator is returned from FilterUintConfigUpdated and is used to iterate over the raw logs and unpacked data for UintConfigUpdated events raised by the AddressBookV2 contract.
type AddressBookV2UintConfigUpdatedIterator struct {
	Event *AddressBookV2UintConfigUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2UintConfigUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2UintConfigUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2UintConfigUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2UintConfigUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2UintConfigUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2UintConfigUpdated represents a UintConfigUpdated event raised by the AddressBookV2 contract.
type AddressBookV2UintConfigUpdated struct {
	ConfigId uint8
	OldValue *big.Int
	NewValue *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUintConfigUpdated is a free log retrieval operation binding the contract event 0x34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd.
//
// Solidity: event UintConfigUpdated(uint8 indexed configId, uint256 oldValue, uint256 newValue)
func (_AddressBookV2 *AddressBookV2Filterer) FilterUintConfigUpdated(opts *bind.FilterOpts, configId []uint8) (*AddressBookV2UintConfigUpdatedIterator, error) {
	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "UintConfigUpdated", configIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2UintConfigUpdatedIterator{contract: _AddressBookV2.contract, event: "UintConfigUpdated", logs: logs, sub: sub}, nil
}

// WatchUintConfigUpdated is a free log subscription operation binding the contract event 0x34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd.
//
// Solidity: event UintConfigUpdated(uint8 indexed configId, uint256 oldValue, uint256 newValue)
func (_AddressBookV2 *AddressBookV2Filterer) WatchUintConfigUpdated(opts *bind.WatchOpts, sink chan<- *AddressBookV2UintConfigUpdated, configId []uint8) (event.Subscription, error) {
	var configIdRule []interface{}
	for _, configIdItem := range configId {
		configIdRule = append(configIdRule, configIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "UintConfigUpdated", configIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2UintConfigUpdated)
				if err := _AddressBookV2.contract.UnpackLog(event, "UintConfigUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUintConfigUpdated is a log parse operation binding the contract event 0x34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd.
//
// Solidity: event UintConfigUpdated(uint8 indexed configId, uint256 oldValue, uint256 newValue)
func (_AddressBookV2 *AddressBookV2Filterer) ParseUintConfigUpdated(log types.Log) (*AddressBookV2UintConfigUpdated, error) {
	event := new(AddressBookV2UintConfigUpdated)
	if err := _AddressBookV2.contract.UnpackLog(event, "UintConfigUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the AddressBookV2 contract.
type AddressBookV2UpgradedIterator struct {
	Event *AddressBookV2Upgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2Upgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2Upgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2Upgraded represents a Upgraded event raised by the AddressBookV2 contract.
type AddressBookV2Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AddressBookV2 *AddressBookV2Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*AddressBookV2UpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2UpgradedIterator{contract: _AddressBookV2.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AddressBookV2 *AddressBookV2Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *AddressBookV2Upgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2Upgraded)
				if err := _AddressBookV2.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_AddressBookV2 *AddressBookV2Filterer) ParseUpgraded(log types.Log) (*AddressBookV2Upgraded, error) {
	event := new(AddressBookV2Upgraded)
	if err := _AddressBookV2.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2ValidatorSuspendedIterator is returned from FilterValidatorSuspended and is used to iterate over the raw logs and unpacked data for ValidatorSuspended events raised by the AddressBookV2 contract.
type AddressBookV2ValidatorSuspendedIterator struct {
	Event *AddressBookV2ValidatorSuspended // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2ValidatorSuspendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2ValidatorSuspended)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2ValidatorSuspended)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2ValidatorSuspendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2ValidatorSuspendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2ValidatorSuspended represents a ValidatorSuspended event raised by the AddressBookV2 contract.
type AddressBookV2ValidatorSuspended struct {
	NodeId common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterValidatorSuspended is a free log retrieval operation binding the contract event 0xb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd.
//
// Solidity: event ValidatorSuspended(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) FilterValidatorSuspended(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2ValidatorSuspendedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "ValidatorSuspended", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2ValidatorSuspendedIterator{contract: _AddressBookV2.contract, event: "ValidatorSuspended", logs: logs, sub: sub}, nil
}

// WatchValidatorSuspended is a free log subscription operation binding the contract event 0xb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd.
//
// Solidity: event ValidatorSuspended(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) WatchValidatorSuspended(opts *bind.WatchOpts, sink chan<- *AddressBookV2ValidatorSuspended, nodeId []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "ValidatorSuspended", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2ValidatorSuspended)
				if err := _AddressBookV2.contract.UnpackLog(event, "ValidatorSuspended", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorSuspended is a log parse operation binding the contract event 0xb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd.
//
// Solidity: event ValidatorSuspended(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) ParseValidatorSuspended(log types.Log) (*AddressBookV2ValidatorSuspended, error) {
	event := new(AddressBookV2ValidatorSuspended)
	if err := _AddressBookV2.contract.UnpackLog(event, "ValidatorSuspended", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2ValidatorUnsuspendedIterator is returned from FilterValidatorUnsuspended and is used to iterate over the raw logs and unpacked data for ValidatorUnsuspended events raised by the AddressBookV2 contract.
type AddressBookV2ValidatorUnsuspendedIterator struct {
	Event *AddressBookV2ValidatorUnsuspended // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2ValidatorUnsuspendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2ValidatorUnsuspended)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2ValidatorUnsuspended)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2ValidatorUnsuspendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2ValidatorUnsuspendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2ValidatorUnsuspended represents a ValidatorUnsuspended event raised by the AddressBookV2 contract.
type AddressBookV2ValidatorUnsuspended struct {
	NodeId common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterValidatorUnsuspended is a free log retrieval operation binding the contract event 0x814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e.
//
// Solidity: event ValidatorUnsuspended(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) FilterValidatorUnsuspended(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2ValidatorUnsuspendedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "ValidatorUnsuspended", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2ValidatorUnsuspendedIterator{contract: _AddressBookV2.contract, event: "ValidatorUnsuspended", logs: logs, sub: sub}, nil
}

// WatchValidatorUnsuspended is a free log subscription operation binding the contract event 0x814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e.
//
// Solidity: event ValidatorUnsuspended(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) WatchValidatorUnsuspended(opts *bind.WatchOpts, sink chan<- *AddressBookV2ValidatorUnsuspended, nodeId []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "ValidatorUnsuspended", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2ValidatorUnsuspended)
				if err := _AddressBookV2.contract.UnpackLog(event, "ValidatorUnsuspended", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorUnsuspended is a log parse operation binding the contract event 0x814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e.
//
// Solidity: event ValidatorUnsuspended(address indexed nodeId)
func (_AddressBookV2 *AddressBookV2Filterer) ParseValidatorUnsuspended(log types.Log) (*AddressBookV2ValidatorUnsuspended, error) {
	event := new(AddressBookV2ValidatorUnsuspended)
	if err := _AddressBookV2.contract.UnpackLog(event, "ValidatorUnsuspended", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2ValidatorsInitializedIterator is returned from FilterValidatorsInitialized and is used to iterate over the raw logs and unpacked data for ValidatorsInitialized events raised by the AddressBookV2 contract.
type AddressBookV2ValidatorsInitializedIterator struct {
	Event *AddressBookV2ValidatorsInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2ValidatorsInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2ValidatorsInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2ValidatorsInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2ValidatorsInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2ValidatorsInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2ValidatorsInitialized represents a ValidatorsInitialized event raised by the AddressBookV2 contract.
type AddressBookV2ValidatorsInitialized struct {
	NodeIds []common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterValidatorsInitialized is a free log retrieval operation binding the contract event 0x820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d952.
//
// Solidity: event ValidatorsInitialized(address[] nodeIds)
func (_AddressBookV2 *AddressBookV2Filterer) FilterValidatorsInitialized(opts *bind.FilterOpts) (*AddressBookV2ValidatorsInitializedIterator, error) {
	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "ValidatorsInitialized")
	if err != nil {
		return nil, err
	}
	return &AddressBookV2ValidatorsInitializedIterator{contract: _AddressBookV2.contract, event: "ValidatorsInitialized", logs: logs, sub: sub}, nil
}

// WatchValidatorsInitialized is a free log subscription operation binding the contract event 0x820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d952.
//
// Solidity: event ValidatorsInitialized(address[] nodeIds)
func (_AddressBookV2 *AddressBookV2Filterer) WatchValidatorsInitialized(opts *bind.WatchOpts, sink chan<- *AddressBookV2ValidatorsInitialized) (event.Subscription, error) {
	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "ValidatorsInitialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2ValidatorsInitialized)
				if err := _AddressBookV2.contract.UnpackLog(event, "ValidatorsInitialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorsInitialized is a log parse operation binding the contract event 0x820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d952.
//
// Solidity: event ValidatorsInitialized(address[] nodeIds)
func (_AddressBookV2 *AddressBookV2Filterer) ParseValidatorsInitialized(log types.Log) (*AddressBookV2ValidatorsInitialized, error) {
	event := new(AddressBookV2ValidatorsInitialized)
	if err := _AddressBookV2.contract.UnpackLog(event, "ValidatorsInitialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AddressBookV2VoterAddressUpdatedIterator is returned from FilterVoterAddressUpdated and is used to iterate over the raw logs and unpacked data for VoterAddressUpdated events raised by the AddressBookV2 contract.
type AddressBookV2VoterAddressUpdatedIterator struct {
	Event *AddressBookV2VoterAddressUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AddressBookV2VoterAddressUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2VoterAddressUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AddressBookV2VoterAddressUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AddressBookV2VoterAddressUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2VoterAddressUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2VoterAddressUpdated represents a VoterAddressUpdated event raised by the AddressBookV2 contract.
type AddressBookV2VoterAddressUpdated struct {
	NodeId          common.Address
	OldVoterAddress common.Address
	NewVoterAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterVoterAddressUpdated is a free log retrieval operation binding the contract event 0x23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d96.
//
// Solidity: event VoterAddressUpdated(address indexed nodeId, address indexed oldVoterAddress, address indexed newVoterAddress)
func (_AddressBookV2 *AddressBookV2Filterer) FilterVoterAddressUpdated(opts *bind.FilterOpts, nodeId []common.Address, oldVoterAddress []common.Address, newVoterAddress []common.Address) (*AddressBookV2VoterAddressUpdatedIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var oldVoterAddressRule []interface{}
	for _, oldVoterAddressItem := range oldVoterAddress {
		oldVoterAddressRule = append(oldVoterAddressRule, oldVoterAddressItem)
	}
	var newVoterAddressRule []interface{}
	for _, newVoterAddressItem := range newVoterAddress {
		newVoterAddressRule = append(newVoterAddressRule, newVoterAddressItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "VoterAddressUpdated", nodeIdRule, oldVoterAddressRule, newVoterAddressRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2VoterAddressUpdatedIterator{contract: _AddressBookV2.contract, event: "VoterAddressUpdated", logs: logs, sub: sub}, nil
}

// WatchVoterAddressUpdated is a free log subscription operation binding the contract event 0x23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d96.
//
// Solidity: event VoterAddressUpdated(address indexed nodeId, address indexed oldVoterAddress, address indexed newVoterAddress)
func (_AddressBookV2 *AddressBookV2Filterer) WatchVoterAddressUpdated(opts *bind.WatchOpts, sink chan<- *AddressBookV2VoterAddressUpdated, nodeId []common.Address, oldVoterAddress []common.Address, newVoterAddress []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var oldVoterAddressRule []interface{}
	for _, oldVoterAddressItem := range oldVoterAddress {
		oldVoterAddressRule = append(oldVoterAddressRule, oldVoterAddressItem)
	}
	var newVoterAddressRule []interface{}
	for _, newVoterAddressItem := range newVoterAddress {
		newVoterAddressRule = append(newVoterAddressRule, newVoterAddressItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "VoterAddressUpdated", nodeIdRule, oldVoterAddressRule, newVoterAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2VoterAddressUpdated)
				if err := _AddressBookV2.contract.UnpackLog(event, "VoterAddressUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoterAddressUpdated is a log parse operation binding the contract event 0x23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d96.
//
// Solidity: event VoterAddressUpdated(address indexed nodeId, address indexed oldVoterAddress, address indexed newVoterAddress)
func (_AddressBookV2 *AddressBookV2Filterer) ParseVoterAddressUpdated(log types.Log) (*AddressBookV2VoterAddressUpdated, error) {
	event := new(AddressBookV2VoterAddressUpdated)
	if err := _AddressBookV2.contract.UnpackLog(event, "VoterAddressUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
