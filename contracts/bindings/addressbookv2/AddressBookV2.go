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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_epochBlockInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CONTRACT_TYPE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_METADATA_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_NODE_BALANCE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_STAKE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SYSTEM_SENDER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"assignGcId\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blsInfo\",\"type\":\"tuple\",\"internalType\":\"structBlsPublicKeyInfo\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadata\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"nodeIdSig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deleteNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"epochBlockInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exit\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"typeList\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"},{\"name\":\"addressList\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllAddressInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllBlsInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeIdList\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"pubkeyList\",\"type\":\"tuple[]\",\"internalType\":\"structBlsPublicKeyInfo[]\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllGovernanceInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"infos\",\"type\":\"tuple[]\",\"internalType\":\"structGovernanceInfo[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesLength\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllProfiles\",\"inputs\":[],\"outputs\":[{\"name\":\"profiles\",\"type\":\"tuple[]\",\"internalType\":\"structProfile[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timeoutAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCfsThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCnInfo\",\"inputs\":[{\"name\":\"_cnNodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigurator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochVACount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFundAddresses\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxCounts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxValActivePausedCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNodeInfo\",\"components\":[{\"name\":\"manager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timeoutAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blsInfo\",\"type\":\"tuple\",\"internalType\":\"structBlsPublicKeyInfo\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadata\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeState\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPfsThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegisteredNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSlotLimits\",\"inputs\":[],\"outputs\":[{\"name\":\"maxSlotAvailable\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minActiveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSlotLimitsFor\",\"inputs\":[{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"maxSlotAvailable\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minActiveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStateCount\",\"inputs\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSuspendedValidators\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSuspender\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimeouts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isActivated\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isConstructed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isUsedAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kirContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"offboard\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pocContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processSystemTransition\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"newStates\",\"type\":\"uint8[]\",\"internalType\":\"enumState[]\"},{\"name\":\"timeoutAts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"epochVACount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"readyCandidate\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"readyValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requirement\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"resume\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeGcId\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"spareContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"suspendValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unreadyCandidate\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unreadyValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unsuspendValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCfsThreshold\",\"inputs\":[{\"name\":\"newCfsThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConfigurator\",\"inputs\":[{\"name\":\"newConfigurator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateIdleTimeout\",\"inputs\":[{\"name\":\"newIdleTimeout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKefAddress\",\"inputs\":[{\"name\":\"newKefAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKifAddress\",\"inputs\":[{\"name\":\"newKifAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKpfAddress\",\"inputs\":[{\"name\":\"newKpfAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateManager\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxCandReadyCount\",\"inputs\":[{\"name\":\"newMaxCandReadyCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxNodeCount\",\"inputs\":[{\"name\":\"newMaxNodeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxValActivePausedCount\",\"inputs\":[{\"name\":\"newMaxValActivePausedCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMetadata\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newMetadata\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauseTimeout\",\"inputs\":[{\"name\":\"newPauseTimeout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePfsThreshold\",\"inputs\":[{\"name\":\"newPfsThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRewardAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newRewardAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateSuspender\",\"inputs\":[{\"name\":\"newSuspender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateVoterAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newVoterAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AddressConfigUpdated\",\"inputs\":[{\"name\":\"configId\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"oldValue\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newValue\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CandidateReadied\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CandidateUnreadied\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochTransitionProcessed\",\"inputs\":[{\"name\":\"epochVACount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GcIdAssigned\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ManagerUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newMetadata\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeCreated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeDeleted\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldRewardAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRewardAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StateChanged\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromState\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumState\"},{\"name\":\"toState\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SystemTransitionProcessed\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"newStates\",\"type\":\"uint8[]\",\"indexed\":false,\"internalType\":\"enumState[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UintConfigUpdated\",\"inputs\":[{\"name\":\"configId\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorSuspended\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorUnsuspended\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorsInitialized\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VoterAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldVoterAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newVoterAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AlreadySuspended\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CnNodeNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FactoryNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GcIdAlreadyAssigned\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GcIdNotAssigned\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientNodeBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LegacyFunctionDeprecated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeIdProofInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSuspended\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyConfigurator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyNodeId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySuspender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySystemTx\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PDEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlotsFull\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakingDeployerMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakingTooLow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimeoutExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60c060405230608052348015610013575f80fd5b506040516160ea3803806160ea833981016040819052610032916100fb565b60a08190528080610041610049565b505050610112565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100995760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100f85780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b5f6020828403121561010b575f80fd5b5051919050565b60805160a051615f9b61014f5f395f818161081701528181613d23015261437d01525f8181613edd01528181613f0601526140d30152615f9b5ff3fe608060405260043610610421575f3560e01c806378b84a5c1161022e578063b858dd9511610138578063d3b54907116100b5578063e70c38f111610079578063e70c38f114610d19578063e8868e9f14610d2d578063f0a92ba814610d42578063f2fde38b14610d56578063ffa1ad7414610d7557610421565b8063d3b5490714610c89578063d9abb38b14610c9d578063da38d49814610cbc578063e4f0d37c14610cdb578063e59d7a8414610cfa57610421565b8063c9a86af2116100fc578063c9a86af214610bfa578063cb1c2b5c14610c19578063cf8c6f5214610c37578063d18c07ab14610c4b578063d267eda514610c6a57610421565b8063b858dd9514610b63578063b9f96f4014610b82578063ba70d01814610ba1578063be535f8b14610bc0578063c732e08514610bdf57610421565b80639d0f5ef1116101c6578063a9ee54721161018a578063a9ee547214610ac3578063ad3cb1cc14610ae2578063b42652e914610b12578063b57873a514610b31578063b756393014610b5057610421565b80639d0f5ef114610a335780639d8cf08f14610a475780639f9e3cba14610a66578063a41b600014610a85578063a4c98ada14610aa457610421565b806378b84a5c1461094d578063793c19461461096c5780637df40c621461098b5780638129fc1c146109aa57806387b7b8fd146109be5780638da5cb5b146109d85780638fabf389146109ec5780639b7ae5ec14610a005780639d0e234d14610a1457610421565b8063453e962e1161032f578063567b0b6c116102c75780636abd623d1161028b5780636abd623d146108c5578063715018a6146108e4578063715b208b146108f8578063766718081461091a57806376a67a511461092e57610421565b8063567b0b6c14610806578063582115fb146108395780635b27b6c914610865578063656f5869146108845780636968b53f146108a357610421565b8063453e962e146106d9578063468e3a7e146106f85780634a8c1fb4146107275780634b6a94cc146107405780634f1ef2861461078357806350a5bb691461079657806350de2fb3146107b457806352d1902d146107d357806353d39bfb146107e757610421565b80631b1a478b116103bd57806325cf09431161038157806325cf094314610652578063291937f5146106665780632aca50911461067a5780632d4ede931461069b578063394f8899146106ba57610421565b80631b1a478b146105a65780631b8f34ca146105c55780631ba3fd58146105e457806321d2320014610605578063229bb8231461062657610421565b806303e6689d14610446578063058529fb1461047457806306bb84711461049357806307ecec3e146104b45780630a4ff239146104d35780630b1fe784146104f557806315575d5a14610516578063160370b81461055f5780631865c57d14610584575b34801561042c575f80fd5b50604051632053d6b560e11b815260040160405180910390fd5b348015610451575f80fd5b5061045a610d89565b604080519283526020830191909152015b60405180910390f35b34801561047f575f80fd5b5061045a61048e366004614df7565b610da9565b34801561049e575f80fd5b506104b26104ad366004614e32565b610dc6565b005b3480156104bf575f80fd5b506104b26104ce366004614e4d565b610fa8565b3480156104de575f80fd5b506104e7611045565b60405190815260200161046b565b348015610500575f80fd5b5061050961105e565b60405161046b9190614eb8565b348015610521575f80fd5b50610535610530366004614e32565b6111a8565b604080516001600160a01b039485168152928416602084015292169181019190915260600161046b565b34801561056a575f80fd5b506105736111ec565b60405161046b959493929190614f82565b34801561058f575f80fd5b50610598611216565b60405161046b929190614fe0565b3480156105b1575f80fd5b506104e76105c036600461500d565b611277565b3480156105d0575f80fd5b506104b26105df36600461506f565b6112bc565b3480156105ef575f80fd5b506105f8611409565b60405161046b9190615109565b348015610610575f80fd5b5061061961141e565b60405161046b919061511b565b348015610631575f80fd5b50610645610640366004614e32565b611439565b60405161046b919061512f565b34801561065d575f80fd5b50610535611466565b348015610671575f80fd5b506104e761149b565b348015610685575f80fd5b5061068e6114ad565b60405161046b919061513d565b3480156106a6575f80fd5b506104b26106b5366004614e32565b6115c3565b3480156106c5575f80fd5b506104b26106d4366004614e4d565b611855565b3480156106e4575f80fd5b506104b26106f3366004614e32565b6119f6565b348015610703575f80fd5b50610717610712366004614e32565b611b27565b604051901515815260200161046b565b348015610732575f80fd5b50600c546107179060ff1681565b34801561074b575f80fd5b506107766040518060400160405280600b81526020016a41646472657373426f6f6b60a81b81525081565b60405161046b91906151cf565b6104b261079136600461530b565b611b54565b3480156107a1575f80fd5b50600c5461071790610100900460ff1681565b3480156107bf575f80fd5b506104b26107ce366004614e32565b611b73565b3480156107de575f80fd5b506104e7611bbc565b3480156107f2575f80fd5b506104b26108013660046153c1565b611bd7565b348015610811575f80fd5b506104e77f000000000000000000000000000000000000000000000000000000000000000081565b348015610844575f80fd5b50610858610853366004614e32565b611e6d565b60405161046b91906154dd565b348015610870575f80fd5b506104b261087f366004614df7565b6121b1565b34801561088f575f80fd5b506104b261089e366004614e32565b6121e8565b3480156108ae575f80fd5b506108b7612254565b60405161046b9291906155cc565b3480156108d0575f80fd5b50600754610619906001600160a01b031681565b3480156108ef575f80fd5b506104b26124f4565b348015610903575f80fd5b5061090c612507565b60405161046b92919061563c565b348015610925575f80fd5b506104e7612853565b348015610939575f80fd5b506104b2610948366004614e32565b61285c565b348015610958575f80fd5b506104b2610967366004614e32565b612928565b348015610977575f80fd5b506104b2610986366004614e32565b61299c565b348015610996575f80fd5b506104b26109a5366004614e32565b6129e3565b3480156109b5575f80fd5b506104b2612a06565b3480156109c9575f80fd5b506106196002600160a01b0381565b3480156109e3575f80fd5b50610619612f27565b3480156109f7575f80fd5b506104e7612f55565b348015610a0b575f80fd5b50610619612f67565b348015610a1f575f80fd5b506104b2610a2e366004614df7565b612f82565b348015610a3e575f80fd5b5061045a612fa4565b348015610a52575f80fd5b506104b2610a61366004614e32565b612fc2565b348015610a71575f80fd5b506104b2610a80366004614e4d565b612fe4565b348015610a90575f80fd5b506104b2610a9f366004614e32565b61314c565b348015610aaf575f80fd5b506104b2610abe366004614df7565b6131c0565b348015610ace575f80fd5b506104b2610add366004614df7565b6131e3565b348015610aed575f80fd5b50610776604051806040016040528060058152602001640352e302e360dc1b81525081565b348015610b1d575f80fd5b506104b2610b2c366004614e32565b613206565b348015610b3c575f80fd5b506104b2610b4b366004614e32565b61330f565b348015610b5b575f80fd5b5060016104e7565b348015610b6e575f80fd5b50600654610619906001600160a01b031681565b348015610b8d575f80fd5b506104b2610b9c366004614e32565b613332565b348015610bac575f80fd5b506104b2610bbb366004614df7565b613370565b348015610bcb575f80fd5b506104b2610bda366004614e32565b613393565b348015610bea575f80fd5b506104e7678ac7230489e8000081565b348015610c05575f80fd5b506104b2610c14366004614e32565b6134ca565b348015610c24575f80fd5b506104e76a0422ca8b0a00a42500000081565b348015610c42575f80fd5b506105f86134ed565b348015610c56575f80fd5b506104b2610c65366004614df7565b613502565b348015610c75575f80fd5b50600554610619906001600160a01b031681565b348015610c94575f80fd5b506104e7613525565b348015610ca8575f80fd5b506104b2610cb7366004614e32565b613537565b348015610cc7575f80fd5b506104b2610cd6366004615696565b613578565b348015610ce6575f80fd5b506104b2610cf5366004614e32565b613612565b348015610d05575f80fd5b506104b2610d14366004614df7565b613687565b348015610d24575f80fd5b5061045a6136aa565b348015610d38575f80fd5b506104e761080081565b348015610d4d575f80fd5b506104e76136ca565b348015610d61575f80fd5b506104b2610d70366004614e32565b6136dc565b348015610d80575f80fd5b506104e7600281565b5f805f610d94613722565b905080600e015481600f015492509250509091565b5f80610db483613746565b610dbd8461377f565b91509150915091565b610dce6137a5565b5f610dd7613722565b90505f6001600160a01b0383165f908152602083905260409020600a015460ff166008811115610e0957610e09614e84565b03610e2757604051634825e09360e01b815260040160405180910390fd5b6001600160a01b0382165f9081526020829052604090206005015415610e6057604051637be80ce960e11b815260040160405180910390fd5b5f816007015f8154610e7190615727565b91829055506001600160a01b0384165f908152602084905260408082206005018390555163e2693e3f60e01b8152919250906104019063e2693e3f90610eb99060040161573f565b602060405180830381865afa158015610ed4573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ef89190615772565b90506001600160a01b03811615610f5f5760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f90610f3590879060040161511b565b5f604051808303815f87803b158015610f4c575f80fd5b505af1925050508015610f5d575060015b505b836001600160a01b03167fe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb83604051610f9a91815260200190565b60405180910390a250505050565b81610fb2816137d9565b6001600160a01b038216610fd95760405163b4fa3fb360e01b815260040160405180910390fd5b5f610fe2613722565b6001600160a01b038086165f8181526020849052604080822080548986166001600160a01b03198216811790925591519596509316938492917f8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda8791a45050505050565b5f611059611051613722565b60010161381c565b905090565b60605f611069613722565b90505f6110788260010161381c565b9050806001600160401b03811115611092576110926151e1565b6040519080825280602002602001820160405280156110f157816020015b6110de6040805160a0810182525f808252602082018190529181018290526060810182905290608082015290565b8152602001906001900390816110b05790505b5092505f5b818110156111a2575f61110c6001850183613825565b6001600160a01b038082165f8181526020888152604091829020825160a081018452938452600181015485169184019190915260028101549093169082015260048201546060820152600a8201549293509091608082019060ff16600881111561117857611178614e84565b81525086848151811061118d5761118d61578d565b602090810291909101015250506001016110f6565b50505090565b5f805f805f806111b787613837565b925092509250806111db576040516342dc2dc560e01b815260040160405180910390fd5b5085945090925090505b9193909250565b60608060605f805f805f805f6112006138af565b939e929d50909b50995090975095505050505050565b6040805160018082528183019092526060915f91829160208083019080368337019050509050611244612f27565b815f815181106112565761125661578d565b6001600160a01b039092166020928302919091019091015292600192509050565b5f611280613722565b6008015f83600881111561129657611296614e84565b60088111156112a7576112a7614e84565b81526020019081526020015f20549050919050565b6112c4613a93565b8584811415806112d45750808314155b156112f25760405163b4fa3fb360e01b815260040160405180910390fd5b5f5b818110156113735761136b8989838181106113115761131161578d565b90506020020160208101906113269190614e32565b8888848181106113385761133861578d565b905060200201602081019061134d919061500d565b87878581811061135f5761135f61578d565b90506020020135613aba565b6001016112f4565b5061137c613d1d565b156113c2578161138a613722565b601001556040518281527fd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e848160529060200160405180910390a15b7fab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e888888886040516113f794939291906157a1565b60405180910390a15050505050505050565b6060611059611416613722565b600501613d4e565b5f611427613722565b601501546001600160a01b0316919050565b5f611442613722565b6001600160a01b039092165f9081526020929092525060409020600a015460ff1690565b5f805f80611472613722565b601281015460138201546014909201546001600160a01b03918216979282169650169350915050565b5f6114a4613722565b600a0154905090565b60605f6114b8613722565b90505f6114c78260010161381c565b9050806001600160401b038111156114e1576114e16151e1565b60405190808252806020026020018201604052801561153157816020015b604080516080810182525f8082526020808301829052928201819052606082015282525f199092019101816114ff5790505b5092505f5b818110156111a2575f61154c6001850183613825565b6001600160a01b038082165f81815260208881526040918290208251608081018452938452600181015485169184019190915260038101549093169082015260058201546060820152875192935090918790859081106115ae576115ae61578d565b60209081029190910101525050600101611536565b806115cd816137d9565b6115d8826001613d5a565b6115f55760405163baf3f0f760e01b815260040160405180910390fd5b5f6115fe613722565b6001600160a01b038085165f9081526020839052604090206003810154929350911615611704576003810180546001600160a01b031916905560405163e2693e3f60e01b81525f906104019063e2693e3f9061165c9060040161573f565b602060405180830381865afa158015611677573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061169b9190615772565b90506001600160a01b038116156117025760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f906116d890889060040161511b565b5f604051808303815f87803b1580156116ef575f80fd5b505af1925050508015611700575060015b505b505b60018181015460028301546001600160a01b038781165f9081526009870160209081526040808320805460ff199081169091559584168352808320805487169055929093168152818120805490941690935592825260088501905290812080549161176e83615836565b9091555061178190506003830185613d8f565b506001600160a01b0384165f90815260208390526040812080546001600160a01b031990811682556001820180548216905560028201805482169055600382018054909116905560048101829055600581018290559060068201816117e68282614d1e565b6117f3600183015f614d1e565b506118039050600883015f614d1e565b611810600983015f614d1e565b50600a01805460ff191690556040516001600160a01b038516907f1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad905f90a250505050565b8161185f816137d9565b6001600160a01b0382166118865760405163b4fa3fb360e01b815260040160405180910390fd5b5f61188f613722565b6001600160a01b038086165f908152602083815260408083206001810154825163e1a12d3560e01b8152925196975090959394169263e1a12d35926004808401939192918290030181865afa1580156118ea573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061190e9190615772565b6001600160a01b03161461193557604051638ed87ef960e01b815260040160405180910390fd5b6001600160a01b0384165f90815260098301602052604090205460ff1615611970576040516316a163b960e11b815260040160405180910390fd5b6002810180546001600160a01b039081165f818152600986016020526040808220805460ff19908116909155898516808452828420805490921660011790915585546001600160a01b0319168117909555519193928492908a16917f270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc6970191a4505050505050565b80611a0081613da3565b611a0b826001613d5a565b611a285760405163baf3f0f760e01b815260040160405180910390fd5b611a3182613dcc565b611a4e5760405163bf74735560e01b815260040160405180910390fd5b678ac7230489e80000826001600160a01b0316311015611a805760405162b8ec7b60e61b815260040160405180910390fd5b5f611a89613722565b905080600f0154611a9a6002611277565b10611ab85760405163848084dd60e01b815260040160405180910390fd5b80600e0154611ac5611045565b10611ae35760405163848084dd60e01b815260040160405180910390fd5b611aef8360025f613aba565b6040516001600160a01b038416907fb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98905f90a2505050565b5f611b30613722565b6001600160a01b039092165f90815260099290920160205250604090205460ff1690565b611b5c613ed2565b611b6582613f76565b611b6f8282613f7e565b5050565b611b7b614036565b60035f80516020615f6f833981519152611b96601584614068565b604080516001600160a01b0392831681529185166020830152015b60405180910390a250565b5f611bc56140c8565b505f80516020615f4f83398151915290565b5f611be189611439565b6008811115611bf257611bf2614e84565b14611c105760405163731918fb60e11b815260040160405180910390fd5b82515f03611c315760405163b4fa3fb360e01b815260040160405180910390fd5b61080082511115611c555760405163b4fa3fb360e01b815260040160405180910390fd5b5f611c5e613722565b9050611c71816009018a8a8a8987614111565b5f604051806101400160405280336001600160a01b031681526020018a6001600160a01b03168152602001896001600160a01b03168152602001886001600160a01b031681526020015f81526020015f815260200187815260200186815260200185815260200160016008811115611ceb57611ceb614e84565b90526001600160a01b03808c165f9081526020858152604091829020845181549085166001600160a01b031991821617825591850151600182018054918616918416919091179055918401516002830180549185169183169190911790556060840151600383018054919094169116179091556080820151600482015560a0820151600582015560c08201518051929350839260068301908190611d8f90826158db565b5060208201516001820190611da490826158db565b50505060e08201516008820190611dbb90826158db565b506101008201516009820190611dd190826158db565b50610120820151600a8201805460ff19166001836008811115611df657611df6614e84565b0217905550611e0b915050600383018b6142a8565b5060015f9081526008830160205260408120805491611e2983615727565b90915550506040516001600160a01b038b16907f55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792905f90a250505050505050505050565b611e75614d55565b5f611e7e613722565b6001600160a01b0384165f908152602091909152604081209150600a82015460ff166008811115611eb157611eb1614e84565b03611ecf57604051634825e09360e01b815260040160405180910390fd5b604080516101408101825282546001600160a01b0390811682526001840154811660208301526002840154811682840152600384015416606082015260048301546080820152600583015460a082015281518083019092526006830180549192849260c08501929082908290611f449061584b565b80601f0160208091040260200160405190810160405280929190818152602001828054611f709061584b565b8015611fbb5780601f10611f9257610100808354040283529160200191611fbb565b820191905f5260205f20905b815481529060010190602001808311611f9e57829003601f168201915b50505050508152602001600182018054611fd49061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546120009061584b565b801561204b5780601f106120225761010080835404028352916020019161204b565b820191905f5260205f20905b81548152906001019060200180831161202e57829003601f168201915b50505050508152505081526020016008820180546120689061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546120949061584b565b80156120df5780601f106120b6576101008083540402835291602001916120df565b820191905f5260205f20905b8154815290600101906020018083116120c257829003601f168201915b505050505081526020016009820180546120f89061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546121249061584b565b801561216f5780601f106121465761010080835404028352916020019161216f565b820191905f5260205f20905b81548152906001019060200180831161215257829003601f168201915b5050509183525050600a82015460209091019060ff16600881111561219657612196614e84565b60088111156121a7576121a7614e84565b9052509392505050565b6121b96137a5565b60055f80516020615f2f8339815191526121d46011846142bc565b604080519182526020820185905201611bb1565b806121f281613da3565b6121fd826004613d5a565b61221a5760405163baf3f0f760e01b815260040160405180910390fd5b61222382613dcc565b6122405760405163bf74735560e01b815260040160405180910390fd5b611b6f82600561224f856142dd565b613aba565b6060805f612260613722565b90505f61226f8260010161381c565b9050806001600160401b03811115612289576122896151e1565b6040519080825280602002602001820160405280156122b2578160200160208202803683370190505b509350806001600160401b038111156122cd576122cd6151e1565b60405190808252806020026020018201604052801561231257816020015b60408051808201909152606080825260208201528152602001906001900390816122eb5790505b5092505f5b818110156124ed5761232c6001840182613825565b85828151811061233e5761233e61578d565b60200260200101906001600160a01b031690816001600160a01b031681525050825f015f8683815181106123745761237461578d565b60200260200101516001600160a01b03166001600160a01b031681526020019081526020015f206006016040518060400160405290815f820180546123b89061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546123e49061584b565b801561242f5780601f106124065761010080835404028352916020019161242f565b820191905f5260205f20905b81548152906001019060200180831161241257829003601f168201915b505050505081526020016001820180546124489061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546124749061584b565b80156124bf5780601f10612496576101008083540402835291602001916124bf565b820191905f5260205f20905b8154815290600101906020018083116124a257829003601f168201915b5050505050815250508482815181106124da576124da61578d565b6020908102919091010152600101612317565b5050509091565b6124fc614036565b6125055f614307565b565b600c54606090819060ff16612530575050604080515f8082526020820190815281830190925291565b5f805f805f61253d6138af565b845194995092975090955093509150612557816003615991565b6125629060026159a8565b6001600160401b03811115612579576125796151e1565b6040519080825280602002602001820160405280156125a2578160200160208202803683370190505b5097506125b0816003615991565b6125bb9060026159a8565b6001600160401b038111156125d2576125d26151e1565b6040519080825280602002602001820160405280156125fb578160200160208202803683370190505b5096505f805b82811015612788575f8a838151811061261c5761261c61578d565b602002602001019060ff16908160ff16815250508781815181106126425761264261578d565b602002602001015189838061265690615727565b9450815181106126685761266861578d565b60200260200101906001600160a01b031690816001600160a01b03168152505060018a838151811061269c5761269c61578d565b602002602001019060ff16908160ff16815250508681815181106126c2576126c261578d565b60200260200101518983806126d690615727565b9450815181106126e8576126e861578d565b60200260200101906001600160a01b031690816001600160a01b03168152505060028a838151811061271c5761271c61578d565b602002602001019060ff16908160ff16815250508581815181106127425761274261578d565b602002602001015189838061275690615727565b9450815181106127685761276861578d565b6001600160a01b0390921660209283029190910190910152600101612601565b50600389828151811061279d5761279d61578d565b60ff909216602092830291909101909101528388826127bb81615727565b9350815181106127cd576127cd61578d565b60200260200101906001600160a01b031690816001600160a01b03168152505060048982815181106128015761280161578d565b602002602001019060ff16908160ff1681525050828882815181106128285761282861578d565b60200260200101906001600160a01b031690816001600160a01b031681525050505050505050509091565b5f611059614377565b8061286681613da3565b612871826006613d5a565b61288e5760405163baf3f0f760e01b815260040160405180910390fd5b5f612897613722565b90506128a68160100154613746565b6128b06007611277565b106128ce5760405163848084dd60e01b815260040160405180910390fd5b6128db816010015461377f565b6128e56006611277565b116129035760405163848084dd60e01b815260040160405180910390fd5b5f81600c01544261291491906159a8565b905061292284600783613aba565b50505050565b6129306143a2565b5f612939613722565b90506129486005820183613d8f565b6129655760405163d33ff8c160e01b815260040160405180910390fd5b6040516001600160a01b038316907f814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e905f90a25050565b806129a681613da3565b6129b1826007613d5a565b6129ce5760405163baf3f0f760e01b815260040160405180910390fd5b6129d7826143d6565b611b6f8260065f613aba565b6129eb6137a5565b60015f80516020615f6f833981519152611b96601384614068565b5f612a0f61440f565b805490915060ff600160401b82041615906001600160401b03165f81158015612a355750825b90505f826001600160401b03166001148015612a505750303b155b905081158015612a5e575080155b15612a7c5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff191660011785558315612aa657845460ff60401b1916600160401b1785555b60405163e2693e3f60e01b815260206004820152601060248201526f10509d8c91185d1850dbdb9d1c9858dd60821b60448201525f906104019063e2693e3f90606401602060405180830381865afa158015612b04573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190612b289190615772565b90506001600160a01b038116612b515760405163aed5959560e01b815260040160405180910390fd5b5f816001600160a01b031663ebe58ed76040518163ffffffff1660e01b81526004015f60405180830381865afa158015612b8d573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f19168201604052612bb49190810190615c71565b90505f612bbf613722565b9050612bcd825f0151614437565b6060820151600a8201556080820151600b82015560a0820151600c82015560c0820151600d82015560e0820151600e8201556101008201516011820155610120820151600f8201556101408201516012820180546001600160a01b03199081166001600160a01b03938416179091556101608401516013840180548316918416919091179055610180840151601484018054831691841691909117905560208401516015840180548316918416919091179055604084015160168401805490921692169190911790556101a0820151515f5b81811015612e76575f846101a001518281518110612cbf57612cbf61578d565b602002602001015190505f856101c001518381518110612ce157612ce161578d565b6020908102919091018101516001600160a01b038481165f90815260098901845260408082208054600160ff199182168117909255958501518416835281832080548716821790558185015190931682529020805490931617909155905060066101208201819052505f608082018181526001600160a01b0380851683526020888152604093849020855181549084166001600160a01b03199182161782559186015160018201805491851691841691909117905593850151600285018054918416918316919091179055606085015160038501805491909316911617905551600482015560a0820151600582015560c082015180518392919060068301908190612dec90826158db565b5060208201516001820190612e0190826158db565b50505060e08201516008820190612e1890826158db565b506101008201516009820190612e2e90826158db565b50610120820151600a8201805460ff19166001836008811115612e5357612e53614e84565b0217905550612e6891505060018601836142a8565b505050806001019050612c9f565b5060065f9081526008830160205260409081902082905560108301829055606460078401556101a084015190517f820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d95291612ece91615109565b60405180910390a1505050508315612f2057845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b5f612f5e613722565b60110154905090565b5f612f70613722565b601601546001600160a01b0316919050565b612f8a6137a5565b5f5f80516020615f2f8339815191526121d4600c846142bc565b5f80612fba612fb1613722565b60100154610da9565b915091509091565b612fca6137a5565b5f5f80516020615f6f833981519152611b96601284614068565b81612fee816137d9565b5f612ff7613722565b6001600160a01b038086165f9081526020839052604080822060030180548885166001600160a01b0319821617909155905163e2693e3f60e01b8152939450909116916104019063e2693e3f906130509060040161573f565b602060405180830381865afa15801561306b573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061308f9190615772565b90506001600160a01b038116156130fa5760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f906130cc90899060040161511b565b5f604051808303815f87803b1580156130e3575f80fd5b505af11580156130f5573d5f803e3d5ffd5b505050505b846001600160a01b0316826001600160a01b0316876001600160a01b03167f23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d9660405160405180910390a4505050505050565b6131546143a2565b5f61315d613722565b905061316c60058201836142a8565b61318957604051633ad2b1bb60e11b815260040160405180910390fd5b6040516001600160a01b038316907fb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd905f90a25050565b6131c86137a5565b60045f80516020615f2f8339815191526121d4600e846142bc565b6131eb6137a5565b60065f80516020615f2f8339815191526121d4600f846142bc565b8061321081613da3565b5f61321a83611439565b9050600781600881111561323057613230614e84565b036132435761323e836143d6565b613275565b600681600881111561325757613257614e84565b146132755760405163baf3f0f760e01b815260040160405180910390fd5b5f61327e613722565b905061328d8160100154613746565b6132976008611277565b106132b55760405163848084dd60e01b815260040160405180910390fd5b60068260088111156132c9576132c9614e84565b03613303576132db816010015461377f565b6132e56006611277565b116133035760405163848084dd60e01b815260040160405180910390fd5b6129228460085f613aba565b613317614036565b60045f80516020615f6f833981519152611b96601684614068565b8061333c81613da3565b613347826004613d5a565b6133645760405163baf3f0f760e01b815260040160405180910390fd5b611b6f8260015f613aba565b6133786137a5565b60025f80516020615f2f8339815191526121d4600a846142bc565b61339b6137a5565b5f6133a4613722565b6001600160a01b0383165f908152602082905260408120600501549192508190036133e2576040516357024f6d60e11b815260040160405180910390fd5b60405163e2693e3f60e01b81525f906104019063e2693e3f906134079060040161573f565b602060405180830381865afa158015613422573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906134469190615772565b90506001600160a01b038116156134a95760405163f575c5a760e01b8152600481018390526001600160a01b0382169063f575c5a7906024015f604051808303815f87803b158015613496575f80fd5b505af19250505080156134a7575060015b505b50506001600160a01b039091165f9081526020919091526040812060050155565b6134d26137a5565b60025f80516020615f6f833981519152611b96601484614068565b60606110596134fa613722565b600301613d4e565b61350a6137a5565b60035f80516020615f2f8339815191526121d4600b846142bc565b5f61352e613722565b60100154905090565b8061354181613da3565b61354c826005613d5a565b6135695760405163baf3f0f760e01b815260040160405180910390fd5b611b6f82600461224f856142dd565b82613582816137d9565b6108008211156135a55760405163b4fa3fb360e01b815260040160405180910390fd5b82826135af613722565b6001600160a01b0387165f90815260209190915260409020600901916135d6919083615db5565b50836001600160a01b03167f2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf8484604051610f9a929190615e69565b8061361c81613da3565b613627826002613d5a565b6136445760405163baf3f0f760e01b815260040160405180910390fd5b6136508260015f613aba565b6040516001600160a01b038316907f8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2905f90a25050565b61368f6137a5565b60015f80516020615f2f8339815191526121d4600d846142bc565b5f805f6136b5613722565b905080600c015481600d015492509250509091565b5f6136d3613722565b600b0154905090565b6136e4614036565b6001600160a01b038116613716575f604051631e4fbdf760e01b815260040161370d919061511b565b60405180910390fd5b61371f81614307565b50565b7f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff0090565b5f600482101561375757505f919050565b6002613764600384615eab565b61376f9060016159a8565b6137799190615eab565b92915050565b5f600482101561378d575090565b600361379a836002615991565b61376f9060026159a8565b6137ad613722565b601601546001600160a01b031633146125055760405163033b71e160e41b815260040160405180910390fd5b6137e1613722565b6001600160a01b038281165f90815260209290925260409091205416331461371f5760405163605919ad60e11b815260040160405180910390fd5b5f613779825490565b5f6138308383614448565b9392505050565b5f805f80613843613722565b6001600160a01b0386165f908152602082905260408120919250600a82015460ff16600881111561387657613876614e84565b0361388b575f805f94509450945050506111e5565b6001818101546002909201546001600160a01b039283169892169650945092505050565b60608060605f805f6138bf613722565b90505f6138ce8260010161381c565b9050806001600160401b038111156138e8576138e86151e1565b604051908082528060200260200182016040528015613911578160200160208202803683370190505b509650806001600160401b0381111561392c5761392c6151e1565b604051908082528060200260200182016040528015613955578160200160208202803683370190505b509550806001600160401b03811115613970576139706151e1565b604051908082528060200260200182016040528015613999578160200160208202803683370190505b5094505f5b81811015613a6b575f6139b46001850183613825565b6001600160a01b0381165f9081526020869052604090208a519192509082908b90859081106139e5576139e561578d565b6001600160a01b03928316602091820292909201015260018201548a519116908a9085908110613a1757613a1761578d565b6001600160a01b03928316602091820292909201015260028201548951911690899085908110613a4957613a4961578d565b6001600160a01b0390921660209283029190910190910152505060010161399e565b505060138101546012909101549596949593946001600160a01b039182169490911692509050565b336002600160a01b0314612505576040516354d325c360e01b815260040160405180910390fd5b5f613ac3613722565b6001600160a01b0385165f908152602082905260408120600a8101549293509160ff1690816008811115613af957613af9614e84565b1480613b1557505f856008811115613b1357613b13614e84565b145b80613b415750846008811115613b2d57613b2d614e84565b816008811115613b3f57613b3f614e84565b145b15613b4e57505050505050565b826008015f826008811115613b6557613b65614e84565b6008811115613b7657613b76614e84565b81526020019081526020015f205f815480929190613b9390615836565b9190505550826008015f866008811115613baf57613baf614e84565b6008811115613bc057613bc0614e84565b81526020019081526020015f205f815480929190613bdd90615727565b9091555060019050816008811115613bf757613bf7614e84565b148015613c1657506001856008811115613c1357613c13614e84565b14155b15613c3c57613c2860018401876142a8565b50613c366003840187613d8f565b50613c91565b6001816008811115613c5057613c50614e84565b14158015613c6f57506001856008811115613c6d57613c6d614e84565b145b15613c9157613c816001840187613d8f565b50613c8f60038401876142a8565b505b600a8201805486919060ff19166001836008811115613cb257613cb2614e84565b021790555060048201849055846008811115613cd057613cd0614e84565b816008811115613ce257613ce2614e84565b6040516001600160a01b038916907fcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee905f90a4505050505050565b5f613d487f000000000000000000000000000000000000000000000000000000000000000043615ebe565b15919050565b60605f6138308361446e565b5f816008811115613d6d57613d6d614e84565b613d7684611439565b6008811115613d8757613d87614e84565b149392505050565b5f613830836001600160a01b0384166144c7565b336001600160a01b0382161461371f576040516335f1334d60e11b815260040160405180910390fd5b5f80613dd6613722565b6001600160a01b038085165f9081526020928352604080822060010154815163318588a360e11b81529151931694509092849263630b11469260048082019392918290030181865afa158015613e2e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613e529190615ed1565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015613e8e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613eb29190615ed1565b613ebc9190615ee8565b6a0422ca8b0a00a4250000001115949350505050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480613f5857507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316613f4c5f80516020615f4f833981519152546001600160a01b031690565b6001600160a01b031614155b156125055760405163703e46dd60e11b815260040160405180910390fd5b61371f614036565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015613fd8575060408051601f3d908101601f19168201909252613fd591810190615ed1565b60015b613ff75781604051634c9c8ce360e01b815260040161370d919061511b565b5f80516020615f4f833981519152811461402757604051632a87526960e21b81526004810182905260240161370d565b61403183836145b1565b505050565b3361403f612f27565b6001600160a01b031614612505573360405163118cdaa760e01b815260040161370d919061511b565b5f6001600160a01b0382166140905760405163b4fa3fb360e01b815260040160405180910390fd5b5f6140bb847f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff006159a8565b8054939055509092915050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146125055760405163703e46dd60e11b815260040160405180910390fd5b61411d85858585614606565b5f614126614772565b90506001600160a01b03811661414f5760405163cdded31d60e01b815260040160405180910390fd5b60405163669d8d4560e01b815233906001600160a01b0383169063669d8d459061417d90899060040161511b565b602060405180830381865afa158015614198573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906141bc9190615772565b6001600160a01b0316146141e357604051632281776f60e01b815260040160405180910390fd5b6141ee8686846147f4565b5f6141f986866148d0565b6001600160a01b03161480156142755750604051631f7f8a5f60e21b81526001600160a01b03821690637dfe297c9061423690879060040161511b565b602060405180830381865afa158015614251573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906142759190615efb565b156142935760405163b4fa3fb360e01b815260040160405180910390fd5b61429f8787878761497b565b50505050505050565b5f613830836001600160a01b038416614a45565b5f815f036140905760405163b4fa3fb360e01b815260040160405180910390fd5b5f6142e6613722565b6001600160a01b039092165f90815260209290925250604090206004015490565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f6110597f000000000000000000000000000000000000000000000000000000000000000043615eab565b6143aa613722565b601501546001600160a01b031633146125055760405163333f4e6560e01b815260040160405180910390fd5b5f6143e0826142dd565b905080158015906143f15750804210155b15611b6f5760405163b48d5fc760e01b815260040160405180910390fd5b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00613779565b61443f614a91565b61371f81614ab6565b5f825f01828154811061445d5761445d61578d565b905f5260205f200154905092915050565b6060815f018054806020026020016040519081016040528092919081815260200182805480156144bb57602002820191905f5260205f20905b8154815260200190600101908083116144a7575b50505050509050919050565b5f81815260018301602052604081205480156145a1575f6144e9600183615ee8565b85549091505f906144fc90600190615ee8565b905080821461455b575f865f01828154811061451a5761451a61578d565b905f5260205f200154905080875f01848154811061453a5761453a61578d565b5f918252602080832090910192909255918252600188019052604090208390555b855486908061456c5761456c615f1a565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f905560019350505050613779565b5f915050613779565b5092915050565b6145ba82614abe565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156145fe576140318282614b18565b611b6f614ba8565b6001600160a01b038416158061462357506001600160a01b038316155b8061463557506001600160a01b038216155b156146535760405163b4fa3fb360e01b815260040160405180910390fd5b826001600160a01b0316846001600160a01b031614806146845750816001600160a01b0316846001600160a01b0316145b806146a05750816001600160a01b0316836001600160a01b0316145b156146be5760405163b4fa3fb360e01b815260040160405180910390fd5b80515160301415806146d65750806020015151606014155b156146f45760405163b4fa3fb360e01b815260040160405180910390fd5b805180516020909101207fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd2931480614754575060208082015180519101207f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21145b156129225760405163b4fa3fb360e01b815260040160405180910390fd5b60405163e2693e3f60e01b815260206004820152601060248201526f436e5374616b696e67466163746f727960801b60448201525f906104019063e2693e3f90606401602060405180830381865afa1580156147d0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906110599190615772565b604080517f23ae25c387ef8bd2c14b622e10202a494f464e31b796f40e83e9aecdf9cb42fb602082015246918101919091523060608201523360808201526001600160a01b0380851660a0830152831660c08201525f9060e0016040516020818303038152906040528051906020012090505f806148728385614bc7565b5090925090505f81600381111561488b5761488b614e84565b1415806148aa5750856001600160a01b0316826001600160a01b031614155b156148c857604051631ea9ff4d60e21b815260040160405180910390fd5b505050505050565b5f826001600160a01b031663e1a12d356040518163ffffffff1660e01b8152600401602060405180830381865afa15801561490d573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906149319190615772565b90506001600160a01b0381161580159061495d5750816001600160a01b0316816001600160a01b031614155b156137795760405163b4fa3fb360e01b815260040160405180910390fd5b6001600160a01b0383165f9081526020859052604090205460ff16806149b857506001600160a01b0382165f9081526020859052604090205460ff165b806149da57506001600160a01b0381165f9081526020859052604090205460ff165b156149f8576040516316a163b960e11b815260040160405180910390fd5b6001600160a01b039283165f90815260209490945260408085208054600160ff19918216811790925593851686528186208054851682179055919093168452919092208054909216179055565b5f818152600183016020526040812054614a8a57508154600181810184555f848152602080822090930184905584548482528286019093526040902091909155613779565b505f613779565b614a99614c10565b61250557604051631afcd79f60e31b815260040160405180910390fd5b6136e4614a91565b806001600160a01b03163b5f03614aea5780604051634c9c8ce360e01b815260040161370d919061511b565b5f80516020615f4f83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b60605f614b258484614c29565b9050808015614b4657505f3d1180614b4657505f846001600160a01b03163b115b15614b5b57614b53614c3c565b915050613779565b8015614b7c5783604051639996b31560e01b815260040161370d919061511b565b3d15614b8f57614b8a614c55565b6145aa565b60405163d6bda27560e01b815260040160405180910390fd5b34156125055760405163b398979f60e01b815260040160405180910390fd5b5f805f8351604103614bfe576020840151604085015160608601515f1a614bf088828585614c60565b955095509550505050614c09565b505081515f91506002905b9250925092565b5f614c1961440f565b54600160401b900460ff16919050565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b5f80806fa2a8918ca85bafe22016d0b997e4df60600160ff1b03841115614c8f57505f91506003905082614d14565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015614ce0573d5f803e3d5ffd5b5050604051601f1901519150506001600160a01b038116614d0b57505f925060019150829050614d14565b92505f91508190505b9450945094915050565b508054614d2a9061584b565b5f825580601f10614d39575050565b601f0160209004905f5260205f209081019061371f9190614ddf565b6040518061014001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f81526020015f8152602001614dc6604051806040016040528060608152602001606081525090565b815260606020820181905260408201819052015f905290565b5b80821115614df3575f8155600101614de0565b5090565b5f60208284031215614e07575f80fd5b5035919050565b6001600160a01b038116811461371f575f80fd5b8035614e2d81614e0e565b919050565b5f60208284031215614e42575f80fd5b813561383081614e0e565b5f8060408385031215614e5e575f80fd5b8235614e6981614e0e565b91506020830135614e7981614e0e565b809150509250929050565b634e487b7160e01b5f52602160045260245ffd5b60098110614eb457634e487b7160e01b5f52602160045260245ffd5b9052565b602080825282518282018190525f919060409081850190868401855b82811015614f3257815180516001600160a01b039081168652878201518116888701528682015116868601526060808201519086015260809081015190614f1d81870183614e98565b505060a0939093019290850190600101614ed4565b5091979650505050505050565b5f815180845260208085019450602084015f5b83811015614f775781516001600160a01b031687529582019590820190600101614f52565b509495945050505050565b60a081525f614f9460a0830188614f3f565b8281036020840152614fa68188614f3f565b90508281036040840152614fba8187614f3f565b6001600160a01b0395861660608501529390941660809092019190915250949350505050565b604081525f614ff26040830185614f3f565b90508260208301529392505050565b6009811061371f575f80fd5b5f6020828403121561501d575f80fd5b813561383081615001565b5f8083601f840112615038575f80fd5b5081356001600160401b0381111561504e575f80fd5b6020830191508360208260051b8501011115615068575f80fd5b9250929050565b5f805f805f805f6080888a031215615085575f80fd5b87356001600160401b038082111561509b575f80fd5b6150a78b838c01615028565b909950975060208a01359150808211156150bf575f80fd5b6150cb8b838c01615028565b909750955060408a01359150808211156150e3575f80fd5b506150f08a828b01615028565b989b979a50959894979596606090950135949350505050565b602081525f6138306020830184614f3f565b6001600160a01b0391909116815260200190565b602081016137798284614e98565b602080825282518282018190525f919060409081850190868401855b82811015614f3257815180516001600160a01b039081168652878201518116888701528682015116868601526060908101519085015260809093019290850190600101615159565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f61383060208301846151a1565b634e487b7160e01b5f52604160045260245ffd5b604080519081016001600160401b0381118282101715615217576152176151e1565b60405290565b60405161014081016001600160401b0381118282101715615217576152176151e1565b6040516101e081016001600160401b0381118282101715615217576152176151e1565b604051601f8201601f191681016001600160401b038111828210171561528b5761528b6151e1565b604052919050565b5f6001600160401b038211156152ab576152ab6151e1565b50601f01601f191660200190565b5f82601f8301126152c8575f80fd5b81356152db6152d682615293565b615263565b8181528460208386010111156152ef575f80fd5b816020850160208301375f918101602001919091529392505050565b5f806040838503121561531c575f80fd5b823561532781614e0e565b915060208301356001600160401b03811115615341575f80fd5b61534d858286016152b9565b9150509250929050565b5f60408284031215615367575f80fd5b61536f6151f5565b905081356001600160401b0380821115615387575f80fd5b615393858386016152b9565b835260208401359150808211156153a8575f80fd5b506153b5848285016152b9565b60208301525092915050565b5f805f805f805f80610100898b0312156153d9575f80fd5b6153e289614e22565b97506153f060208a01614e22565b96506153fe60408a01614e22565b955061540c60608a01614e22565b945060808901356001600160401b0380821115615427575f80fd5b6154338c838d01615357565b955060a08b0135915080821115615448575f80fd5b6154548c838d016152b9565b945060c08b0135915080821115615469575f80fd5b6154758c838d016152b9565b935060e08b013591508082111561548a575f80fd5b506154978b828c016152b9565b9150509295985092959890939650565b5f8151604084526154bb60408501826151a1565b9050602083015184820360208601526154d482826151a1565b95945050505050565b602081526154f76020820183516001600160a01b03169052565b5f602083015161551260408401826001600160a01b03169052565b5060408301516001600160a01b03811660608401525060608301516001600160a01b038116608084015250608083015160a083015260a083015160c083015260c08301516101408060e085015261556d6101608501836154a7565b915060e0850151601f1961010081878603018188015261558d85846151a1565b9450808801519250506101208187860301818801526155ac85846151a1565b945080880151925050506155c282860182614e98565b5090949350505050565b604081525f6155de6040830185614f3f565b6020838203818501528185518084528284019150828160051b8501018388015f5b8381101561562d57601f1987840301855261561b8383516154a7565b948601949250908501906001016155ff565b50909998505050505050505050565b604080825283519082018190525f906020906060840190828701845b8281101561567757815160ff1684529284019290840190600101615658565b505050838103602085015261568c8186614f3f565b9695505050505050565b5f805f604084860312156156a8575f80fd5b83356156b381614e0e565b925060208401356001600160401b03808211156156ce575f80fd5b818601915086601f8301126156e1575f80fd5b8135818111156156ef575f80fd5b876020828501011115615700575f80fd5b6020830194508093505050509250925092565b634e487b7160e01b5f52601160045260245ffd5b5f6001820161573857615738615713565b5060010190565b6020808252600e908201526d29ba30b5b4b733aa3930b1b5b2b960911b604082015260600190565b8051614e2d81614e0e565b5f60208284031215615782575f80fd5b815161383081614e0e565b634e487b7160e01b5f52603260045260245ffd5b604080825281018490525f8560608301825b878110156157e35782356157c681614e0e565b6001600160a01b03168252602092830192909101906001016157b3565b508381036020858101919091528582529150859082015f5b8681101561582957823561580e81615001565b6158188382614e98565b5091830191908301906001016157fb565b5098975050505050505050565b5f8161584457615844615713565b505f190190565b600181811c9082168061585f57607f821691505b60208210810361587d57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561403157805f5260205f20601f840160051c810160208510156158a85750805b601f840160051c820191505b81811015612f20575f81556001016158b4565b5f19600383901b1c191660019190911b1790565b81516001600160401b038111156158f4576158f46151e1565b61590881615902845461584b565b84615883565b602080601f831160018114615936575f84156159245750858301515b61592e85826158c7565b8655506148c8565b5f85815260208120601f198616915b8281101561596457888601518255948401946001909101908401615945565b508582101561598157878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b808202811582820484141761377957613779615713565b8082018082111561377957613779615713565b5f6001600160401b038211156159d3576159d36151e1565b5060051b60200190565b5f82601f8301126159ec575f80fd5b815160206159fc6152d6836159bb565b8083825260208201915060208460051b870101935086841115615a1d575f80fd5b602086015b84811015615a42578051615a3581614e0e565b8352918301918301615a22565b509695505050505050565b5f82601f830112615a5c575f80fd5b8151615a6a6152d682615293565b818152846020838601011115615a7e575f80fd5b8160208501602083015e5f918101602001919091529392505050565b5f60408284031215615aaa575f80fd5b615ab26151f5565b905081516001600160401b0380821115615aca575f80fd5b615ad685838601615a4d565b83526020840151915080821115615aeb575f80fd5b506153b584828501615a4d565b8051614e2d81615001565b5f82601f830112615b12575f80fd5b81516020615b226152d6836159bb565b82815260059290921b84018101918181019086841115615b40575f80fd5b8286015b84811015615a425780516001600160401b0380821115615b62575f80fd5b90880190610140828b03601f1901811315615b7b575f80fd5b615b8361521d565b615b8e888501615767565b81526040615b9d818601615767565b898301526060615bae818701615767565b8284015260809150615bc1828701615767565b818401525060a0808601518284015260c0915081860151818401525060e08086015185811115615bef575f80fd5b615bfd8f8c838a0101615a9a565b838501525061010091508186015185811115615c17575f80fd5b615c258f8c838a0101615a4d565b8285015250506101208086015185811115615c3e575f80fd5b615c4c8f8c838a0101615a4d565b8385015250615c5c848701615af8565b90830152508652505050918301918301615b44565b5f60208284031215615c81575f80fd5b81516001600160401b0380821115615c97575f80fd5b908301906101e08286031215615cab575f80fd5b615cb3615240565b615cbc83615767565b8152615cca60208401615767565b6020820152615cdb60408401615767565b6040820152606083015160608201526080830151608082015260a083015160a082015260c083015160c082015260e083015160e0820152610100808401518183015250610120808401518183015250610140615d38818501615767565b90820152610160615d4a848201615767565b90820152610180615d5c848201615767565b908201526101a08381015183811115615d73575f80fd5b615d7f888287016159dd565b8284015250506101c08084015183811115615d98575f80fd5b615da488828701615b03565b918301919091525095945050505050565b6001600160401b03831115615dcc57615dcc6151e1565b615de083615dda835461584b565b83615883565b5f601f841160018114615e0c575f8515615dfa5750838201355b615e0486826158c7565b845550612f20565b5f83815260208120601f198716915b82811015615e3b5786850135825560209485019460019092019101615e1b565b5086821015615e57575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b60208152816020820152818360408301375f818301604090810191909152601f909201601f19160101919050565b634e487b7160e01b5f52601260045260245ffd5b5f82615eb957615eb9615e97565b500490565b5f82615ecc57615ecc615e97565b500690565b5f60208284031215615ee1575f80fd5b5051919050565b8181038181111561377957613779615713565b5f60208284031215615f0b575f80fd5b81518015158114613830575f80fd5b634e487b7160e01b5f52603160045260245ffdfe34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbcd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4ba164736f6c6343000819000a",
}

// AddressBookV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use AddressBookV2MetaData.ABI instead.
var AddressBookV2ABI = AddressBookV2MetaData.ABI

// AddressBookV2BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressBookV2BinRuntime = `608060405260043610610421575f3560e01c806378b84a5c1161022e578063b858dd9511610138578063d3b54907116100b5578063e70c38f111610079578063e70c38f114610d19578063e8868e9f14610d2d578063f0a92ba814610d42578063f2fde38b14610d56578063ffa1ad7414610d7557610421565b8063d3b5490714610c89578063d9abb38b14610c9d578063da38d49814610cbc578063e4f0d37c14610cdb578063e59d7a8414610cfa57610421565b8063c9a86af2116100fc578063c9a86af214610bfa578063cb1c2b5c14610c19578063cf8c6f5214610c37578063d18c07ab14610c4b578063d267eda514610c6a57610421565b8063b858dd9514610b63578063b9f96f4014610b82578063ba70d01814610ba1578063be535f8b14610bc0578063c732e08514610bdf57610421565b80639d0f5ef1116101c6578063a9ee54721161018a578063a9ee547214610ac3578063ad3cb1cc14610ae2578063b42652e914610b12578063b57873a514610b31578063b756393014610b5057610421565b80639d0f5ef114610a335780639d8cf08f14610a475780639f9e3cba14610a66578063a41b600014610a85578063a4c98ada14610aa457610421565b806378b84a5c1461094d578063793c19461461096c5780637df40c621461098b5780638129fc1c146109aa57806387b7b8fd146109be5780638da5cb5b146109d85780638fabf389146109ec5780639b7ae5ec14610a005780639d0e234d14610a1457610421565b8063453e962e1161032f578063567b0b6c116102c75780636abd623d1161028b5780636abd623d146108c5578063715018a6146108e4578063715b208b146108f8578063766718081461091a57806376a67a511461092e57610421565b8063567b0b6c14610806578063582115fb146108395780635b27b6c914610865578063656f5869146108845780636968b53f146108a357610421565b8063453e962e146106d9578063468e3a7e146106f85780634a8c1fb4146107275780634b6a94cc146107405780634f1ef2861461078357806350a5bb691461079657806350de2fb3146107b457806352d1902d146107d357806353d39bfb146107e757610421565b80631b1a478b116103bd57806325cf09431161038157806325cf094314610652578063291937f5146106665780632aca50911461067a5780632d4ede931461069b578063394f8899146106ba57610421565b80631b1a478b146105a65780631b8f34ca146105c55780631ba3fd58146105e457806321d2320014610605578063229bb8231461062657610421565b806303e6689d14610446578063058529fb1461047457806306bb84711461049357806307ecec3e146104b45780630a4ff239146104d35780630b1fe784146104f557806315575d5a14610516578063160370b81461055f5780631865c57d14610584575b34801561042c575f80fd5b50604051632053d6b560e11b815260040160405180910390fd5b348015610451575f80fd5b5061045a610d89565b604080519283526020830191909152015b60405180910390f35b34801561047f575f80fd5b5061045a61048e366004614df7565b610da9565b34801561049e575f80fd5b506104b26104ad366004614e32565b610dc6565b005b3480156104bf575f80fd5b506104b26104ce366004614e4d565b610fa8565b3480156104de575f80fd5b506104e7611045565b60405190815260200161046b565b348015610500575f80fd5b5061050961105e565b60405161046b9190614eb8565b348015610521575f80fd5b50610535610530366004614e32565b6111a8565b604080516001600160a01b039485168152928416602084015292169181019190915260600161046b565b34801561056a575f80fd5b506105736111ec565b60405161046b959493929190614f82565b34801561058f575f80fd5b50610598611216565b60405161046b929190614fe0565b3480156105b1575f80fd5b506104e76105c036600461500d565b611277565b3480156105d0575f80fd5b506104b26105df36600461506f565b6112bc565b3480156105ef575f80fd5b506105f8611409565b60405161046b9190615109565b348015610610575f80fd5b5061061961141e565b60405161046b919061511b565b348015610631575f80fd5b50610645610640366004614e32565b611439565b60405161046b919061512f565b34801561065d575f80fd5b50610535611466565b348015610671575f80fd5b506104e761149b565b348015610685575f80fd5b5061068e6114ad565b60405161046b919061513d565b3480156106a6575f80fd5b506104b26106b5366004614e32565b6115c3565b3480156106c5575f80fd5b506104b26106d4366004614e4d565b611855565b3480156106e4575f80fd5b506104b26106f3366004614e32565b6119f6565b348015610703575f80fd5b50610717610712366004614e32565b611b27565b604051901515815260200161046b565b348015610732575f80fd5b50600c546107179060ff1681565b34801561074b575f80fd5b506107766040518060400160405280600b81526020016a41646472657373426f6f6b60a81b81525081565b60405161046b91906151cf565b6104b261079136600461530b565b611b54565b3480156107a1575f80fd5b50600c5461071790610100900460ff1681565b3480156107bf575f80fd5b506104b26107ce366004614e32565b611b73565b3480156107de575f80fd5b506104e7611bbc565b3480156107f2575f80fd5b506104b26108013660046153c1565b611bd7565b348015610811575f80fd5b506104e77f000000000000000000000000000000000000000000000000000000000000000081565b348015610844575f80fd5b50610858610853366004614e32565b611e6d565b60405161046b91906154dd565b348015610870575f80fd5b506104b261087f366004614df7565b6121b1565b34801561088f575f80fd5b506104b261089e366004614e32565b6121e8565b3480156108ae575f80fd5b506108b7612254565b60405161046b9291906155cc565b3480156108d0575f80fd5b50600754610619906001600160a01b031681565b3480156108ef575f80fd5b506104b26124f4565b348015610903575f80fd5b5061090c612507565b60405161046b92919061563c565b348015610925575f80fd5b506104e7612853565b348015610939575f80fd5b506104b2610948366004614e32565b61285c565b348015610958575f80fd5b506104b2610967366004614e32565b612928565b348015610977575f80fd5b506104b2610986366004614e32565b61299c565b348015610996575f80fd5b506104b26109a5366004614e32565b6129e3565b3480156109b5575f80fd5b506104b2612a06565b3480156109c9575f80fd5b506106196002600160a01b0381565b3480156109e3575f80fd5b50610619612f27565b3480156109f7575f80fd5b506104e7612f55565b348015610a0b575f80fd5b50610619612f67565b348015610a1f575f80fd5b506104b2610a2e366004614df7565b612f82565b348015610a3e575f80fd5b5061045a612fa4565b348015610a52575f80fd5b506104b2610a61366004614e32565b612fc2565b348015610a71575f80fd5b506104b2610a80366004614e4d565b612fe4565b348015610a90575f80fd5b506104b2610a9f366004614e32565b61314c565b348015610aaf575f80fd5b506104b2610abe366004614df7565b6131c0565b348015610ace575f80fd5b506104b2610add366004614df7565b6131e3565b348015610aed575f80fd5b50610776604051806040016040528060058152602001640352e302e360dc1b81525081565b348015610b1d575f80fd5b506104b2610b2c366004614e32565b613206565b348015610b3c575f80fd5b506104b2610b4b366004614e32565b61330f565b348015610b5b575f80fd5b5060016104e7565b348015610b6e575f80fd5b50600654610619906001600160a01b031681565b348015610b8d575f80fd5b506104b2610b9c366004614e32565b613332565b348015610bac575f80fd5b506104b2610bbb366004614df7565b613370565b348015610bcb575f80fd5b506104b2610bda366004614e32565b613393565b348015610bea575f80fd5b506104e7678ac7230489e8000081565b348015610c05575f80fd5b506104b2610c14366004614e32565b6134ca565b348015610c24575f80fd5b506104e76a0422ca8b0a00a42500000081565b348015610c42575f80fd5b506105f86134ed565b348015610c56575f80fd5b506104b2610c65366004614df7565b613502565b348015610c75575f80fd5b50600554610619906001600160a01b031681565b348015610c94575f80fd5b506104e7613525565b348015610ca8575f80fd5b506104b2610cb7366004614e32565b613537565b348015610cc7575f80fd5b506104b2610cd6366004615696565b613578565b348015610ce6575f80fd5b506104b2610cf5366004614e32565b613612565b348015610d05575f80fd5b506104b2610d14366004614df7565b613687565b348015610d24575f80fd5b5061045a6136aa565b348015610d38575f80fd5b506104e761080081565b348015610d4d575f80fd5b506104e76136ca565b348015610d61575f80fd5b506104b2610d70366004614e32565b6136dc565b348015610d80575f80fd5b506104e7600281565b5f805f610d94613722565b905080600e015481600f015492509250509091565b5f80610db483613746565b610dbd8461377f565b91509150915091565b610dce6137a5565b5f610dd7613722565b90505f6001600160a01b0383165f908152602083905260409020600a015460ff166008811115610e0957610e09614e84565b03610e2757604051634825e09360e01b815260040160405180910390fd5b6001600160a01b0382165f9081526020829052604090206005015415610e6057604051637be80ce960e11b815260040160405180910390fd5b5f816007015f8154610e7190615727565b91829055506001600160a01b0384165f908152602084905260408082206005018390555163e2693e3f60e01b8152919250906104019063e2693e3f90610eb99060040161573f565b602060405180830381865afa158015610ed4573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ef89190615772565b90506001600160a01b03811615610f5f5760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f90610f3590879060040161511b565b5f604051808303815f87803b158015610f4c575f80fd5b505af1925050508015610f5d575060015b505b836001600160a01b03167fe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb83604051610f9a91815260200190565b60405180910390a250505050565b81610fb2816137d9565b6001600160a01b038216610fd95760405163b4fa3fb360e01b815260040160405180910390fd5b5f610fe2613722565b6001600160a01b038086165f8181526020849052604080822080548986166001600160a01b03198216811790925591519596509316938492917f8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda8791a45050505050565b5f611059611051613722565b60010161381c565b905090565b60605f611069613722565b90505f6110788260010161381c565b9050806001600160401b03811115611092576110926151e1565b6040519080825280602002602001820160405280156110f157816020015b6110de6040805160a0810182525f808252602082018190529181018290526060810182905290608082015290565b8152602001906001900390816110b05790505b5092505f5b818110156111a2575f61110c6001850183613825565b6001600160a01b038082165f8181526020888152604091829020825160a081018452938452600181015485169184019190915260028101549093169082015260048201546060820152600a8201549293509091608082019060ff16600881111561117857611178614e84565b81525086848151811061118d5761118d61578d565b602090810291909101015250506001016110f6565b50505090565b5f805f805f806111b787613837565b925092509250806111db576040516342dc2dc560e01b815260040160405180910390fd5b5085945090925090505b9193909250565b60608060605f805f805f805f6112006138af565b939e929d50909b50995090975095505050505050565b6040805160018082528183019092526060915f91829160208083019080368337019050509050611244612f27565b815f815181106112565761125661578d565b6001600160a01b039092166020928302919091019091015292600192509050565b5f611280613722565b6008015f83600881111561129657611296614e84565b60088111156112a7576112a7614e84565b81526020019081526020015f20549050919050565b6112c4613a93565b8584811415806112d45750808314155b156112f25760405163b4fa3fb360e01b815260040160405180910390fd5b5f5b818110156113735761136b8989838181106113115761131161578d565b90506020020160208101906113269190614e32565b8888848181106113385761133861578d565b905060200201602081019061134d919061500d565b87878581811061135f5761135f61578d565b90506020020135613aba565b6001016112f4565b5061137c613d1d565b156113c2578161138a613722565b601001556040518281527fd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e848160529060200160405180910390a15b7fab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e888888886040516113f794939291906157a1565b60405180910390a15050505050505050565b6060611059611416613722565b600501613d4e565b5f611427613722565b601501546001600160a01b0316919050565b5f611442613722565b6001600160a01b039092165f9081526020929092525060409020600a015460ff1690565b5f805f80611472613722565b601281015460138201546014909201546001600160a01b03918216979282169650169350915050565b5f6114a4613722565b600a0154905090565b60605f6114b8613722565b90505f6114c78260010161381c565b9050806001600160401b038111156114e1576114e16151e1565b60405190808252806020026020018201604052801561153157816020015b604080516080810182525f8082526020808301829052928201819052606082015282525f199092019101816114ff5790505b5092505f5b818110156111a2575f61154c6001850183613825565b6001600160a01b038082165f81815260208881526040918290208251608081018452938452600181015485169184019190915260038101549093169082015260058201546060820152875192935090918790859081106115ae576115ae61578d565b60209081029190910101525050600101611536565b806115cd816137d9565b6115d8826001613d5a565b6115f55760405163baf3f0f760e01b815260040160405180910390fd5b5f6115fe613722565b6001600160a01b038085165f9081526020839052604090206003810154929350911615611704576003810180546001600160a01b031916905560405163e2693e3f60e01b81525f906104019063e2693e3f9061165c9060040161573f565b602060405180830381865afa158015611677573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061169b9190615772565b90506001600160a01b038116156117025760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f906116d890889060040161511b565b5f604051808303815f87803b1580156116ef575f80fd5b505af1925050508015611700575060015b505b505b60018181015460028301546001600160a01b038781165f9081526009870160209081526040808320805460ff199081169091559584168352808320805487169055929093168152818120805490941690935592825260088501905290812080549161176e83615836565b9091555061178190506003830185613d8f565b506001600160a01b0384165f90815260208390526040812080546001600160a01b031990811682556001820180548216905560028201805482169055600382018054909116905560048101829055600581018290559060068201816117e68282614d1e565b6117f3600183015f614d1e565b506118039050600883015f614d1e565b611810600983015f614d1e565b50600a01805460ff191690556040516001600160a01b038516907f1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad905f90a250505050565b8161185f816137d9565b6001600160a01b0382166118865760405163b4fa3fb360e01b815260040160405180910390fd5b5f61188f613722565b6001600160a01b038086165f908152602083815260408083206001810154825163e1a12d3560e01b8152925196975090959394169263e1a12d35926004808401939192918290030181865afa1580156118ea573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061190e9190615772565b6001600160a01b03161461193557604051638ed87ef960e01b815260040160405180910390fd5b6001600160a01b0384165f90815260098301602052604090205460ff1615611970576040516316a163b960e11b815260040160405180910390fd5b6002810180546001600160a01b039081165f818152600986016020526040808220805460ff19908116909155898516808452828420805490921660011790915585546001600160a01b0319168117909555519193928492908a16917f270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc6970191a4505050505050565b80611a0081613da3565b611a0b826001613d5a565b611a285760405163baf3f0f760e01b815260040160405180910390fd5b611a3182613dcc565b611a4e5760405163bf74735560e01b815260040160405180910390fd5b678ac7230489e80000826001600160a01b0316311015611a805760405162b8ec7b60e61b815260040160405180910390fd5b5f611a89613722565b905080600f0154611a9a6002611277565b10611ab85760405163848084dd60e01b815260040160405180910390fd5b80600e0154611ac5611045565b10611ae35760405163848084dd60e01b815260040160405180910390fd5b611aef8360025f613aba565b6040516001600160a01b038416907fb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98905f90a2505050565b5f611b30613722565b6001600160a01b039092165f90815260099290920160205250604090205460ff1690565b611b5c613ed2565b611b6582613f76565b611b6f8282613f7e565b5050565b611b7b614036565b60035f80516020615f6f833981519152611b96601584614068565b604080516001600160a01b0392831681529185166020830152015b60405180910390a250565b5f611bc56140c8565b505f80516020615f4f83398151915290565b5f611be189611439565b6008811115611bf257611bf2614e84565b14611c105760405163731918fb60e11b815260040160405180910390fd5b82515f03611c315760405163b4fa3fb360e01b815260040160405180910390fd5b61080082511115611c555760405163b4fa3fb360e01b815260040160405180910390fd5b5f611c5e613722565b9050611c71816009018a8a8a8987614111565b5f604051806101400160405280336001600160a01b031681526020018a6001600160a01b03168152602001896001600160a01b03168152602001886001600160a01b031681526020015f81526020015f815260200187815260200186815260200185815260200160016008811115611ceb57611ceb614e84565b90526001600160a01b03808c165f9081526020858152604091829020845181549085166001600160a01b031991821617825591850151600182018054918616918416919091179055918401516002830180549185169183169190911790556060840151600383018054919094169116179091556080820151600482015560a0820151600582015560c08201518051929350839260068301908190611d8f90826158db565b5060208201516001820190611da490826158db565b50505060e08201516008820190611dbb90826158db565b506101008201516009820190611dd190826158db565b50610120820151600a8201805460ff19166001836008811115611df657611df6614e84565b0217905550611e0b915050600383018b6142a8565b5060015f9081526008830160205260408120805491611e2983615727565b90915550506040516001600160a01b038b16907f55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792905f90a250505050505050505050565b611e75614d55565b5f611e7e613722565b6001600160a01b0384165f908152602091909152604081209150600a82015460ff166008811115611eb157611eb1614e84565b03611ecf57604051634825e09360e01b815260040160405180910390fd5b604080516101408101825282546001600160a01b0390811682526001840154811660208301526002840154811682840152600384015416606082015260048301546080820152600583015460a082015281518083019092526006830180549192849260c08501929082908290611f449061584b565b80601f0160208091040260200160405190810160405280929190818152602001828054611f709061584b565b8015611fbb5780601f10611f9257610100808354040283529160200191611fbb565b820191905f5260205f20905b815481529060010190602001808311611f9e57829003601f168201915b50505050508152602001600182018054611fd49061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546120009061584b565b801561204b5780601f106120225761010080835404028352916020019161204b565b820191905f5260205f20905b81548152906001019060200180831161202e57829003601f168201915b50505050508152505081526020016008820180546120689061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546120949061584b565b80156120df5780601f106120b6576101008083540402835291602001916120df565b820191905f5260205f20905b8154815290600101906020018083116120c257829003601f168201915b505050505081526020016009820180546120f89061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546121249061584b565b801561216f5780601f106121465761010080835404028352916020019161216f565b820191905f5260205f20905b81548152906001019060200180831161215257829003601f168201915b5050509183525050600a82015460209091019060ff16600881111561219657612196614e84565b60088111156121a7576121a7614e84565b9052509392505050565b6121b96137a5565b60055f80516020615f2f8339815191526121d46011846142bc565b604080519182526020820185905201611bb1565b806121f281613da3565b6121fd826004613d5a565b61221a5760405163baf3f0f760e01b815260040160405180910390fd5b61222382613dcc565b6122405760405163bf74735560e01b815260040160405180910390fd5b611b6f82600561224f856142dd565b613aba565b6060805f612260613722565b90505f61226f8260010161381c565b9050806001600160401b03811115612289576122896151e1565b6040519080825280602002602001820160405280156122b2578160200160208202803683370190505b509350806001600160401b038111156122cd576122cd6151e1565b60405190808252806020026020018201604052801561231257816020015b60408051808201909152606080825260208201528152602001906001900390816122eb5790505b5092505f5b818110156124ed5761232c6001840182613825565b85828151811061233e5761233e61578d565b60200260200101906001600160a01b031690816001600160a01b031681525050825f015f8683815181106123745761237461578d565b60200260200101516001600160a01b03166001600160a01b031681526020019081526020015f206006016040518060400160405290815f820180546123b89061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546123e49061584b565b801561242f5780601f106124065761010080835404028352916020019161242f565b820191905f5260205f20905b81548152906001019060200180831161241257829003601f168201915b505050505081526020016001820180546124489061584b565b80601f01602080910402602001604051908101604052809291908181526020018280546124749061584b565b80156124bf5780601f10612496576101008083540402835291602001916124bf565b820191905f5260205f20905b8154815290600101906020018083116124a257829003601f168201915b5050505050815250508482815181106124da576124da61578d565b6020908102919091010152600101612317565b5050509091565b6124fc614036565b6125055f614307565b565b600c54606090819060ff16612530575050604080515f8082526020820190815281830190925291565b5f805f805f61253d6138af565b845194995092975090955093509150612557816003615991565b6125629060026159a8565b6001600160401b03811115612579576125796151e1565b6040519080825280602002602001820160405280156125a2578160200160208202803683370190505b5097506125b0816003615991565b6125bb9060026159a8565b6001600160401b038111156125d2576125d26151e1565b6040519080825280602002602001820160405280156125fb578160200160208202803683370190505b5096505f805b82811015612788575f8a838151811061261c5761261c61578d565b602002602001019060ff16908160ff16815250508781815181106126425761264261578d565b602002602001015189838061265690615727565b9450815181106126685761266861578d565b60200260200101906001600160a01b031690816001600160a01b03168152505060018a838151811061269c5761269c61578d565b602002602001019060ff16908160ff16815250508681815181106126c2576126c261578d565b60200260200101518983806126d690615727565b9450815181106126e8576126e861578d565b60200260200101906001600160a01b031690816001600160a01b03168152505060028a838151811061271c5761271c61578d565b602002602001019060ff16908160ff16815250508581815181106127425761274261578d565b602002602001015189838061275690615727565b9450815181106127685761276861578d565b6001600160a01b0390921660209283029190910190910152600101612601565b50600389828151811061279d5761279d61578d565b60ff909216602092830291909101909101528388826127bb81615727565b9350815181106127cd576127cd61578d565b60200260200101906001600160a01b031690816001600160a01b03168152505060048982815181106128015761280161578d565b602002602001019060ff16908160ff1681525050828882815181106128285761282861578d565b60200260200101906001600160a01b031690816001600160a01b031681525050505050505050509091565b5f611059614377565b8061286681613da3565b612871826006613d5a565b61288e5760405163baf3f0f760e01b815260040160405180910390fd5b5f612897613722565b90506128a68160100154613746565b6128b06007611277565b106128ce5760405163848084dd60e01b815260040160405180910390fd5b6128db816010015461377f565b6128e56006611277565b116129035760405163848084dd60e01b815260040160405180910390fd5b5f81600c01544261291491906159a8565b905061292284600783613aba565b50505050565b6129306143a2565b5f612939613722565b90506129486005820183613d8f565b6129655760405163d33ff8c160e01b815260040160405180910390fd5b6040516001600160a01b038316907f814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e905f90a25050565b806129a681613da3565b6129b1826007613d5a565b6129ce5760405163baf3f0f760e01b815260040160405180910390fd5b6129d7826143d6565b611b6f8260065f613aba565b6129eb6137a5565b60015f80516020615f6f833981519152611b96601384614068565b5f612a0f61440f565b805490915060ff600160401b82041615906001600160401b03165f81158015612a355750825b90505f826001600160401b03166001148015612a505750303b155b905081158015612a5e575080155b15612a7c5760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff191660011785558315612aa657845460ff60401b1916600160401b1785555b60405163e2693e3f60e01b815260206004820152601060248201526f10509d8c91185d1850dbdb9d1c9858dd60821b60448201525f906104019063e2693e3f90606401602060405180830381865afa158015612b04573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190612b289190615772565b90506001600160a01b038116612b515760405163aed5959560e01b815260040160405180910390fd5b5f816001600160a01b031663ebe58ed76040518163ffffffff1660e01b81526004015f60405180830381865afa158015612b8d573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f19168201604052612bb49190810190615c71565b90505f612bbf613722565b9050612bcd825f0151614437565b6060820151600a8201556080820151600b82015560a0820151600c82015560c0820151600d82015560e0820151600e8201556101008201516011820155610120820151600f8201556101408201516012820180546001600160a01b03199081166001600160a01b03938416179091556101608401516013840180548316918416919091179055610180840151601484018054831691841691909117905560208401516015840180548316918416919091179055604084015160168401805490921692169190911790556101a0820151515f5b81811015612e76575f846101a001518281518110612cbf57612cbf61578d565b602002602001015190505f856101c001518381518110612ce157612ce161578d565b6020908102919091018101516001600160a01b038481165f90815260098901845260408082208054600160ff199182168117909255958501518416835281832080548716821790558185015190931682529020805490931617909155905060066101208201819052505f608082018181526001600160a01b0380851683526020888152604093849020855181549084166001600160a01b03199182161782559186015160018201805491851691841691909117905593850151600285018054918416918316919091179055606085015160038501805491909316911617905551600482015560a0820151600582015560c082015180518392919060068301908190612dec90826158db565b5060208201516001820190612e0190826158db565b50505060e08201516008820190612e1890826158db565b506101008201516009820190612e2e90826158db565b50610120820151600a8201805460ff19166001836008811115612e5357612e53614e84565b0217905550612e6891505060018601836142a8565b505050806001019050612c9f565b5060065f9081526008830160205260409081902082905560108301829055606460078401556101a084015190517f820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d95291612ece91615109565b60405180910390a1505050508315612f2057845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b5f612f5e613722565b60110154905090565b5f612f70613722565b601601546001600160a01b0316919050565b612f8a6137a5565b5f5f80516020615f2f8339815191526121d4600c846142bc565b5f80612fba612fb1613722565b60100154610da9565b915091509091565b612fca6137a5565b5f5f80516020615f6f833981519152611b96601284614068565b81612fee816137d9565b5f612ff7613722565b6001600160a01b038086165f9081526020839052604080822060030180548885166001600160a01b0319821617909155905163e2693e3f60e01b8152939450909116916104019063e2693e3f906130509060040161573f565b602060405180830381865afa15801561306b573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061308f9190615772565b90506001600160a01b038116156130fa5760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f906130cc90899060040161511b565b5f604051808303815f87803b1580156130e3575f80fd5b505af11580156130f5573d5f803e3d5ffd5b505050505b846001600160a01b0316826001600160a01b0316876001600160a01b03167f23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d9660405160405180910390a4505050505050565b6131546143a2565b5f61315d613722565b905061316c60058201836142a8565b61318957604051633ad2b1bb60e11b815260040160405180910390fd5b6040516001600160a01b038316907fb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd905f90a25050565b6131c86137a5565b60045f80516020615f2f8339815191526121d4600e846142bc565b6131eb6137a5565b60065f80516020615f2f8339815191526121d4600f846142bc565b8061321081613da3565b5f61321a83611439565b9050600781600881111561323057613230614e84565b036132435761323e836143d6565b613275565b600681600881111561325757613257614e84565b146132755760405163baf3f0f760e01b815260040160405180910390fd5b5f61327e613722565b905061328d8160100154613746565b6132976008611277565b106132b55760405163848084dd60e01b815260040160405180910390fd5b60068260088111156132c9576132c9614e84565b03613303576132db816010015461377f565b6132e56006611277565b116133035760405163848084dd60e01b815260040160405180910390fd5b6129228460085f613aba565b613317614036565b60045f80516020615f6f833981519152611b96601684614068565b8061333c81613da3565b613347826004613d5a565b6133645760405163baf3f0f760e01b815260040160405180910390fd5b611b6f8260015f613aba565b6133786137a5565b60025f80516020615f2f8339815191526121d4600a846142bc565b61339b6137a5565b5f6133a4613722565b6001600160a01b0383165f908152602082905260408120600501549192508190036133e2576040516357024f6d60e11b815260040160405180910390fd5b60405163e2693e3f60e01b81525f906104019063e2693e3f906134079060040161573f565b602060405180830381865afa158015613422573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906134469190615772565b90506001600160a01b038116156134a95760405163f575c5a760e01b8152600481018390526001600160a01b0382169063f575c5a7906024015f604051808303815f87803b158015613496575f80fd5b505af19250505080156134a7575060015b505b50506001600160a01b039091165f9081526020919091526040812060050155565b6134d26137a5565b60025f80516020615f6f833981519152611b96601484614068565b60606110596134fa613722565b600301613d4e565b61350a6137a5565b60035f80516020615f2f8339815191526121d4600b846142bc565b5f61352e613722565b60100154905090565b8061354181613da3565b61354c826005613d5a565b6135695760405163baf3f0f760e01b815260040160405180910390fd5b611b6f82600461224f856142dd565b82613582816137d9565b6108008211156135a55760405163b4fa3fb360e01b815260040160405180910390fd5b82826135af613722565b6001600160a01b0387165f90815260209190915260409020600901916135d6919083615db5565b50836001600160a01b03167f2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf8484604051610f9a929190615e69565b8061361c81613da3565b613627826002613d5a565b6136445760405163baf3f0f760e01b815260040160405180910390fd5b6136508260015f613aba565b6040516001600160a01b038316907f8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2905f90a25050565b61368f6137a5565b60015f80516020615f2f8339815191526121d4600d846142bc565b5f805f6136b5613722565b905080600c015481600d015492509250509091565b5f6136d3613722565b600b0154905090565b6136e4614036565b6001600160a01b038116613716575f604051631e4fbdf760e01b815260040161370d919061511b565b60405180910390fd5b61371f81614307565b50565b7f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff0090565b5f600482101561375757505f919050565b6002613764600384615eab565b61376f9060016159a8565b6137799190615eab565b92915050565b5f600482101561378d575090565b600361379a836002615991565b61376f9060026159a8565b6137ad613722565b601601546001600160a01b031633146125055760405163033b71e160e41b815260040160405180910390fd5b6137e1613722565b6001600160a01b038281165f90815260209290925260409091205416331461371f5760405163605919ad60e11b815260040160405180910390fd5b5f613779825490565b5f6138308383614448565b9392505050565b5f805f80613843613722565b6001600160a01b0386165f908152602082905260408120919250600a82015460ff16600881111561387657613876614e84565b0361388b575f805f94509450945050506111e5565b6001818101546002909201546001600160a01b039283169892169650945092505050565b60608060605f805f6138bf613722565b90505f6138ce8260010161381c565b9050806001600160401b038111156138e8576138e86151e1565b604051908082528060200260200182016040528015613911578160200160208202803683370190505b509650806001600160401b0381111561392c5761392c6151e1565b604051908082528060200260200182016040528015613955578160200160208202803683370190505b509550806001600160401b03811115613970576139706151e1565b604051908082528060200260200182016040528015613999578160200160208202803683370190505b5094505f5b81811015613a6b575f6139b46001850183613825565b6001600160a01b0381165f9081526020869052604090208a519192509082908b90859081106139e5576139e561578d565b6001600160a01b03928316602091820292909201015260018201548a519116908a9085908110613a1757613a1761578d565b6001600160a01b03928316602091820292909201015260028201548951911690899085908110613a4957613a4961578d565b6001600160a01b0390921660209283029190910190910152505060010161399e565b505060138101546012909101549596949593946001600160a01b039182169490911692509050565b336002600160a01b0314612505576040516354d325c360e01b815260040160405180910390fd5b5f613ac3613722565b6001600160a01b0385165f908152602082905260408120600a8101549293509160ff1690816008811115613af957613af9614e84565b1480613b1557505f856008811115613b1357613b13614e84565b145b80613b415750846008811115613b2d57613b2d614e84565b816008811115613b3f57613b3f614e84565b145b15613b4e57505050505050565b826008015f826008811115613b6557613b65614e84565b6008811115613b7657613b76614e84565b81526020019081526020015f205f815480929190613b9390615836565b9190505550826008015f866008811115613baf57613baf614e84565b6008811115613bc057613bc0614e84565b81526020019081526020015f205f815480929190613bdd90615727565b9091555060019050816008811115613bf757613bf7614e84565b148015613c1657506001856008811115613c1357613c13614e84565b14155b15613c3c57613c2860018401876142a8565b50613c366003840187613d8f565b50613c91565b6001816008811115613c5057613c50614e84565b14158015613c6f57506001856008811115613c6d57613c6d614e84565b145b15613c9157613c816001840187613d8f565b50613c8f60038401876142a8565b505b600a8201805486919060ff19166001836008811115613cb257613cb2614e84565b021790555060048201849055846008811115613cd057613cd0614e84565b816008811115613ce257613ce2614e84565b6040516001600160a01b038916907fcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee905f90a4505050505050565b5f613d487f000000000000000000000000000000000000000000000000000000000000000043615ebe565b15919050565b60605f6138308361446e565b5f816008811115613d6d57613d6d614e84565b613d7684611439565b6008811115613d8757613d87614e84565b149392505050565b5f613830836001600160a01b0384166144c7565b336001600160a01b0382161461371f576040516335f1334d60e11b815260040160405180910390fd5b5f80613dd6613722565b6001600160a01b038085165f9081526020928352604080822060010154815163318588a360e11b81529151931694509092849263630b11469260048082019392918290030181865afa158015613e2e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613e529190615ed1565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015613e8e573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613eb29190615ed1565b613ebc9190615ee8565b6a0422ca8b0a00a4250000001115949350505050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480613f5857507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316613f4c5f80516020615f4f833981519152546001600160a01b031690565b6001600160a01b031614155b156125055760405163703e46dd60e11b815260040160405180910390fd5b61371f614036565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015613fd8575060408051601f3d908101601f19168201909252613fd591810190615ed1565b60015b613ff75781604051634c9c8ce360e01b815260040161370d919061511b565b5f80516020615f4f833981519152811461402757604051632a87526960e21b81526004810182905260240161370d565b61403183836145b1565b505050565b3361403f612f27565b6001600160a01b031614612505573360405163118cdaa760e01b815260040161370d919061511b565b5f6001600160a01b0382166140905760405163b4fa3fb360e01b815260040160405180910390fd5b5f6140bb847f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff006159a8565b8054939055509092915050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146125055760405163703e46dd60e11b815260040160405180910390fd5b61411d85858585614606565b5f614126614772565b90506001600160a01b03811661414f5760405163cdded31d60e01b815260040160405180910390fd5b60405163669d8d4560e01b815233906001600160a01b0383169063669d8d459061417d90899060040161511b565b602060405180830381865afa158015614198573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906141bc9190615772565b6001600160a01b0316146141e357604051632281776f60e01b815260040160405180910390fd5b6141ee8686846147f4565b5f6141f986866148d0565b6001600160a01b03161480156142755750604051631f7f8a5f60e21b81526001600160a01b03821690637dfe297c9061423690879060040161511b565b602060405180830381865afa158015614251573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906142759190615efb565b156142935760405163b4fa3fb360e01b815260040160405180910390fd5b61429f8787878761497b565b50505050505050565b5f613830836001600160a01b038416614a45565b5f815f036140905760405163b4fa3fb360e01b815260040160405180910390fd5b5f6142e6613722565b6001600160a01b039092165f90815260209290925250604090206004015490565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f6110597f000000000000000000000000000000000000000000000000000000000000000043615eab565b6143aa613722565b601501546001600160a01b031633146125055760405163333f4e6560e01b815260040160405180910390fd5b5f6143e0826142dd565b905080158015906143f15750804210155b15611b6f5760405163b48d5fc760e01b815260040160405180910390fd5b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00613779565b61443f614a91565b61371f81614ab6565b5f825f01828154811061445d5761445d61578d565b905f5260205f200154905092915050565b6060815f018054806020026020016040519081016040528092919081815260200182805480156144bb57602002820191905f5260205f20905b8154815260200190600101908083116144a7575b50505050509050919050565b5f81815260018301602052604081205480156145a1575f6144e9600183615ee8565b85549091505f906144fc90600190615ee8565b905080821461455b575f865f01828154811061451a5761451a61578d565b905f5260205f200154905080875f01848154811061453a5761453a61578d565b5f918252602080832090910192909255918252600188019052604090208390555b855486908061456c5761456c615f1a565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f905560019350505050613779565b5f915050613779565b5092915050565b6145ba82614abe565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156145fe576140318282614b18565b611b6f614ba8565b6001600160a01b038416158061462357506001600160a01b038316155b8061463557506001600160a01b038216155b156146535760405163b4fa3fb360e01b815260040160405180910390fd5b826001600160a01b0316846001600160a01b031614806146845750816001600160a01b0316846001600160a01b0316145b806146a05750816001600160a01b0316836001600160a01b0316145b156146be5760405163b4fa3fb360e01b815260040160405180910390fd5b80515160301415806146d65750806020015151606014155b156146f45760405163b4fa3fb360e01b815260040160405180910390fd5b805180516020909101207fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd2931480614754575060208082015180519101207f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21145b156129225760405163b4fa3fb360e01b815260040160405180910390fd5b60405163e2693e3f60e01b815260206004820152601060248201526f436e5374616b696e67466163746f727960801b60448201525f906104019063e2693e3f90606401602060405180830381865afa1580156147d0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906110599190615772565b604080517f23ae25c387ef8bd2c14b622e10202a494f464e31b796f40e83e9aecdf9cb42fb602082015246918101919091523060608201523360808201526001600160a01b0380851660a0830152831660c08201525f9060e0016040516020818303038152906040528051906020012090505f806148728385614bc7565b5090925090505f81600381111561488b5761488b614e84565b1415806148aa5750856001600160a01b0316826001600160a01b031614155b156148c857604051631ea9ff4d60e21b815260040160405180910390fd5b505050505050565b5f826001600160a01b031663e1a12d356040518163ffffffff1660e01b8152600401602060405180830381865afa15801561490d573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906149319190615772565b90506001600160a01b0381161580159061495d5750816001600160a01b0316816001600160a01b031614155b156137795760405163b4fa3fb360e01b815260040160405180910390fd5b6001600160a01b0383165f9081526020859052604090205460ff16806149b857506001600160a01b0382165f9081526020859052604090205460ff165b806149da57506001600160a01b0381165f9081526020859052604090205460ff165b156149f8576040516316a163b960e11b815260040160405180910390fd5b6001600160a01b039283165f90815260209490945260408085208054600160ff19918216811790925593851686528186208054851682179055919093168452919092208054909216179055565b5f818152600183016020526040812054614a8a57508154600181810184555f848152602080822090930184905584548482528286019093526040902091909155613779565b505f613779565b614a99614c10565b61250557604051631afcd79f60e31b815260040160405180910390fd5b6136e4614a91565b806001600160a01b03163b5f03614aea5780604051634c9c8ce360e01b815260040161370d919061511b565b5f80516020615f4f83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b60605f614b258484614c29565b9050808015614b4657505f3d1180614b4657505f846001600160a01b03163b115b15614b5b57614b53614c3c565b915050613779565b8015614b7c5783604051639996b31560e01b815260040161370d919061511b565b3d15614b8f57614b8a614c55565b6145aa565b60405163d6bda27560e01b815260040160405180910390fd5b34156125055760405163b398979f60e01b815260040160405180910390fd5b5f805f8351604103614bfe576020840151604085015160608601515f1a614bf088828585614c60565b955095509550505050614c09565b505081515f91506002905b9250925092565b5f614c1961440f565b54600160401b900460ff16919050565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b5f80806fa2a8918ca85bafe22016d0b997e4df60600160ff1b03841115614c8f57505f91506003905082614d14565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015614ce0573d5f803e3d5ffd5b5050604051601f1901519150506001600160a01b038116614d0b57505f925060019150829050614d14565b92505f91508190505b9450945094915050565b508054614d2a9061584b565b5f825580601f10614d39575050565b601f0160209004905f5260205f209081019061371f9190614ddf565b6040518061014001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f81526020015f8152602001614dc6604051806040016040528060608152602001606081525090565b815260606020820181905260408201819052015f905290565b5b80821115614df3575f8155600101614de0565b5090565b5f60208284031215614e07575f80fd5b5035919050565b6001600160a01b038116811461371f575f80fd5b8035614e2d81614e0e565b919050565b5f60208284031215614e42575f80fd5b813561383081614e0e565b5f8060408385031215614e5e575f80fd5b8235614e6981614e0e565b91506020830135614e7981614e0e565b809150509250929050565b634e487b7160e01b5f52602160045260245ffd5b60098110614eb457634e487b7160e01b5f52602160045260245ffd5b9052565b602080825282518282018190525f919060409081850190868401855b82811015614f3257815180516001600160a01b039081168652878201518116888701528682015116868601526060808201519086015260809081015190614f1d81870183614e98565b505060a0939093019290850190600101614ed4565b5091979650505050505050565b5f815180845260208085019450602084015f5b83811015614f775781516001600160a01b031687529582019590820190600101614f52565b509495945050505050565b60a081525f614f9460a0830188614f3f565b8281036020840152614fa68188614f3f565b90508281036040840152614fba8187614f3f565b6001600160a01b0395861660608501529390941660809092019190915250949350505050565b604081525f614ff26040830185614f3f565b90508260208301529392505050565b6009811061371f575f80fd5b5f6020828403121561501d575f80fd5b813561383081615001565b5f8083601f840112615038575f80fd5b5081356001600160401b0381111561504e575f80fd5b6020830191508360208260051b8501011115615068575f80fd5b9250929050565b5f805f805f805f6080888a031215615085575f80fd5b87356001600160401b038082111561509b575f80fd5b6150a78b838c01615028565b909950975060208a01359150808211156150bf575f80fd5b6150cb8b838c01615028565b909750955060408a01359150808211156150e3575f80fd5b506150f08a828b01615028565b989b979a50959894979596606090950135949350505050565b602081525f6138306020830184614f3f565b6001600160a01b0391909116815260200190565b602081016137798284614e98565b602080825282518282018190525f919060409081850190868401855b82811015614f3257815180516001600160a01b039081168652878201518116888701528682015116868601526060908101519085015260809093019290850190600101615159565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f61383060208301846151a1565b634e487b7160e01b5f52604160045260245ffd5b604080519081016001600160401b0381118282101715615217576152176151e1565b60405290565b60405161014081016001600160401b0381118282101715615217576152176151e1565b6040516101e081016001600160401b0381118282101715615217576152176151e1565b604051601f8201601f191681016001600160401b038111828210171561528b5761528b6151e1565b604052919050565b5f6001600160401b038211156152ab576152ab6151e1565b50601f01601f191660200190565b5f82601f8301126152c8575f80fd5b81356152db6152d682615293565b615263565b8181528460208386010111156152ef575f80fd5b816020850160208301375f918101602001919091529392505050565b5f806040838503121561531c575f80fd5b823561532781614e0e565b915060208301356001600160401b03811115615341575f80fd5b61534d858286016152b9565b9150509250929050565b5f60408284031215615367575f80fd5b61536f6151f5565b905081356001600160401b0380821115615387575f80fd5b615393858386016152b9565b835260208401359150808211156153a8575f80fd5b506153b5848285016152b9565b60208301525092915050565b5f805f805f805f80610100898b0312156153d9575f80fd5b6153e289614e22565b97506153f060208a01614e22565b96506153fe60408a01614e22565b955061540c60608a01614e22565b945060808901356001600160401b0380821115615427575f80fd5b6154338c838d01615357565b955060a08b0135915080821115615448575f80fd5b6154548c838d016152b9565b945060c08b0135915080821115615469575f80fd5b6154758c838d016152b9565b935060e08b013591508082111561548a575f80fd5b506154978b828c016152b9565b9150509295985092959890939650565b5f8151604084526154bb60408501826151a1565b9050602083015184820360208601526154d482826151a1565b95945050505050565b602081526154f76020820183516001600160a01b03169052565b5f602083015161551260408401826001600160a01b03169052565b5060408301516001600160a01b03811660608401525060608301516001600160a01b038116608084015250608083015160a083015260a083015160c083015260c08301516101408060e085015261556d6101608501836154a7565b915060e0850151601f1961010081878603018188015261558d85846151a1565b9450808801519250506101208187860301818801526155ac85846151a1565b945080880151925050506155c282860182614e98565b5090949350505050565b604081525f6155de6040830185614f3f565b6020838203818501528185518084528284019150828160051b8501018388015f5b8381101561562d57601f1987840301855261561b8383516154a7565b948601949250908501906001016155ff565b50909998505050505050505050565b604080825283519082018190525f906020906060840190828701845b8281101561567757815160ff1684529284019290840190600101615658565b505050838103602085015261568c8186614f3f565b9695505050505050565b5f805f604084860312156156a8575f80fd5b83356156b381614e0e565b925060208401356001600160401b03808211156156ce575f80fd5b818601915086601f8301126156e1575f80fd5b8135818111156156ef575f80fd5b876020828501011115615700575f80fd5b6020830194508093505050509250925092565b634e487b7160e01b5f52601160045260245ffd5b5f6001820161573857615738615713565b5060010190565b6020808252600e908201526d29ba30b5b4b733aa3930b1b5b2b960911b604082015260600190565b8051614e2d81614e0e565b5f60208284031215615782575f80fd5b815161383081614e0e565b634e487b7160e01b5f52603260045260245ffd5b604080825281018490525f8560608301825b878110156157e35782356157c681614e0e565b6001600160a01b03168252602092830192909101906001016157b3565b508381036020858101919091528582529150859082015f5b8681101561582957823561580e81615001565b6158188382614e98565b5091830191908301906001016157fb565b5098975050505050505050565b5f8161584457615844615713565b505f190190565b600181811c9082168061585f57607f821691505b60208210810361587d57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561403157805f5260205f20601f840160051c810160208510156158a85750805b601f840160051c820191505b81811015612f20575f81556001016158b4565b5f19600383901b1c191660019190911b1790565b81516001600160401b038111156158f4576158f46151e1565b61590881615902845461584b565b84615883565b602080601f831160018114615936575f84156159245750858301515b61592e85826158c7565b8655506148c8565b5f85815260208120601f198616915b8281101561596457888601518255948401946001909101908401615945565b508582101561598157878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b808202811582820484141761377957613779615713565b8082018082111561377957613779615713565b5f6001600160401b038211156159d3576159d36151e1565b5060051b60200190565b5f82601f8301126159ec575f80fd5b815160206159fc6152d6836159bb565b8083825260208201915060208460051b870101935086841115615a1d575f80fd5b602086015b84811015615a42578051615a3581614e0e565b8352918301918301615a22565b509695505050505050565b5f82601f830112615a5c575f80fd5b8151615a6a6152d682615293565b818152846020838601011115615a7e575f80fd5b8160208501602083015e5f918101602001919091529392505050565b5f60408284031215615aaa575f80fd5b615ab26151f5565b905081516001600160401b0380821115615aca575f80fd5b615ad685838601615a4d565b83526020840151915080821115615aeb575f80fd5b506153b584828501615a4d565b8051614e2d81615001565b5f82601f830112615b12575f80fd5b81516020615b226152d6836159bb565b82815260059290921b84018101918181019086841115615b40575f80fd5b8286015b84811015615a425780516001600160401b0380821115615b62575f80fd5b90880190610140828b03601f1901811315615b7b575f80fd5b615b8361521d565b615b8e888501615767565b81526040615b9d818601615767565b898301526060615bae818701615767565b8284015260809150615bc1828701615767565b818401525060a0808601518284015260c0915081860151818401525060e08086015185811115615bef575f80fd5b615bfd8f8c838a0101615a9a565b838501525061010091508186015185811115615c17575f80fd5b615c258f8c838a0101615a4d565b8285015250506101208086015185811115615c3e575f80fd5b615c4c8f8c838a0101615a4d565b8385015250615c5c848701615af8565b90830152508652505050918301918301615b44565b5f60208284031215615c81575f80fd5b81516001600160401b0380821115615c97575f80fd5b908301906101e08286031215615cab575f80fd5b615cb3615240565b615cbc83615767565b8152615cca60208401615767565b6020820152615cdb60408401615767565b6040820152606083015160608201526080830151608082015260a083015160a082015260c083015160c082015260e083015160e0820152610100808401518183015250610120808401518183015250610140615d38818501615767565b90820152610160615d4a848201615767565b90820152610180615d5c848201615767565b908201526101a08381015183811115615d73575f80fd5b615d7f888287016159dd565b8284015250506101c08084015183811115615d98575f80fd5b615da488828701615b03565b918301919091525095945050505050565b6001600160401b03831115615dcc57615dcc6151e1565b615de083615dda835461584b565b83615883565b5f601f841160018114615e0c575f8515615dfa5750838201355b615e0486826158c7565b845550612f20565b5f83815260208120601f198716915b82811015615e3b5786850135825560209485019460019092019101615e1b565b5086821015615e57575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b60208152816020820152818360408301375f818301604090810191909152601f909201601f19160101919050565b634e487b7160e01b5f52601260045260245ffd5b5f82615eb957615eb9615e97565b500490565b5f82615ecc57615ecc615e97565b500690565b5f60208284031215615ee1575f80fd5b5051919050565b8181038181111561377957613779615713565b5f60208284031215615f0b575f80fd5b81518015158114613830575f80fd5b634e487b7160e01b5f52603160045260245ffdfe34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbcd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4ba164736f6c6343000819000a`

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
}, error) {
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
}, error) {
	return _AddressBookV2.Contract.GetAllAddress(&_AddressBookV2.CallOpts)
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_AddressBookV2 *AddressBookV2CallerSession) GetAllAddress() (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error) {
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
}, error) {
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
}, error) {
	return _AddressBookV2.Contract.GetAllBlsInfo(&_AddressBookV2.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_AddressBookV2 *AddressBookV2CallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []BlsPublicKeyInfo
}, error) {
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
}, error) {
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
}, error) {
	return _AddressBookV2.Contract.GetSlotLimits(&_AddressBookV2.CallOpts)
}

// GetSlotLimits is a free data retrieval call binding the contract method 0x9d0f5ef1.
//
// Solidity: function getSlotLimits() view returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2CallerSession) GetSlotLimits() (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error) {
	return _AddressBookV2.Contract.GetSlotLimits(&_AddressBookV2.CallOpts)
}

// GetSlotLimitsFor is a free data retrieval call binding the contract method 0x058529fb.
//
// Solidity: function getSlotLimitsFor(uint256 n) pure returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2Caller) GetSlotLimitsFor(opts *bind.CallOpts, n *big.Int) (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error) {
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
}, error) {
	return _AddressBookV2.Contract.GetSlotLimitsFor(&_AddressBookV2.CallOpts, n)
}

// GetSlotLimitsFor is a free data retrieval call binding the contract method 0x058529fb.
//
// Solidity: function getSlotLimitsFor(uint256 n) pure returns(uint256 maxSlotAvailable, uint256 minActiveCount)
func (_AddressBookV2 *AddressBookV2CallerSession) GetSlotLimitsFor(n *big.Int) (struct {
	MaxSlotAvailable *big.Int
	MinActiveCount   *big.Int
}, error) {
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

// CreateNode is a paid mutator transaction binding the contract method 0x53d39bfb.
//
// Solidity: function createNode(address nodeId, address stakingContract, address rewardAddress, address voterAddress, (bytes,bytes) blsInfo, string name, string metadata, bytes nodeIdSig) returns()
func (_AddressBookV2 *AddressBookV2Transactor) CreateNode(opts *bind.TransactOpts, nodeId common.Address, stakingContract common.Address, rewardAddress common.Address, voterAddress common.Address, blsInfo BlsPublicKeyInfo, name string, metadata string, nodeIdSig []byte) (*types.Transaction, error) {
	return _AddressBookV2.contract.Transact(opts, "createNode", nodeId, stakingContract, rewardAddress, voterAddress, blsInfo, name, metadata, nodeIdSig)
}

// CreateNode is a paid mutator transaction binding the contract method 0x53d39bfb.
//
// Solidity: function createNode(address nodeId, address stakingContract, address rewardAddress, address voterAddress, (bytes,bytes) blsInfo, string name, string metadata, bytes nodeIdSig) returns()
func (_AddressBookV2 *AddressBookV2Session) CreateNode(nodeId common.Address, stakingContract common.Address, rewardAddress common.Address, voterAddress common.Address, blsInfo BlsPublicKeyInfo, name string, metadata string, nodeIdSig []byte) (*types.Transaction, error) {
	return _AddressBookV2.Contract.CreateNode(&_AddressBookV2.TransactOpts, nodeId, stakingContract, rewardAddress, voterAddress, blsInfo, name, metadata, nodeIdSig)
}

// CreateNode is a paid mutator transaction binding the contract method 0x53d39bfb.
//
// Solidity: function createNode(address nodeId, address stakingContract, address rewardAddress, address voterAddress, (bytes,bytes) blsInfo, string name, string metadata, bytes nodeIdSig) returns()
func (_AddressBookV2 *AddressBookV2TransactorSession) CreateNode(nodeId common.Address, stakingContract common.Address, rewardAddress common.Address, voterAddress common.Address, blsInfo BlsPublicKeyInfo, name string, metadata string, nodeIdSig []byte) (*types.Transaction, error) {
	return _AddressBookV2.Contract.CreateNode(&_AddressBookV2.TransactOpts, nodeId, stakingContract, rewardAddress, voterAddress, blsInfo, name, metadata, nodeIdSig)
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
