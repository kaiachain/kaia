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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_epochBlockInterval\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"fallback\",\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"CONTRACT_TYPE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_METADATA_LENGTH\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_NODE_BALANCE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_STAKE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SYSTEM_SENDER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"assignGcId\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"createNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"blsInfo\",\"type\":\"tuple\",\"internalType\":\"structBlsPublicKeyInfo\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadata\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"nodeIdSig\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currentEpoch\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deleteNode\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"epochBlockInterval\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"exit\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getAllAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"typeList\",\"type\":\"uint8[]\",\"internalType\":\"uint8[]\"},{\"name\":\"addressList\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllAddressInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllBlsInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"nodeIdList\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"pubkeyList\",\"type\":\"tuple[]\",\"internalType\":\"structBlsPublicKeyInfo[]\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllGovernanceInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"infos\",\"type\":\"tuple[]\",\"internalType\":\"structGovernanceInfo[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllNodesLength\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAllProfiles\",\"inputs\":[],\"outputs\":[{\"name\":\"profiles\",\"type\":\"tuple[]\",\"internalType\":\"structProfile[]\",\"components\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timeoutAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCfsThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCnInfo\",\"inputs\":[{\"name\":\"_cnNodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getConfigurator\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getEpochVACount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getFundAddresses\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxCounts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxValActivePausedCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeInfo\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structNodeInfo\",\"components\":[{\"name\":\"manager\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"stakingContract\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewardAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"voterAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"timeoutAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blsInfo\",\"type\":\"tuple\",\"internalType\":\"structBlsPublicKeyInfo\",\"components\":[{\"name\":\"publicKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"pop\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"metadata\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getNodeState\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPfsThreshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRegisteredNodes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSlotLimits\",\"inputs\":[],\"outputs\":[{\"name\":\"maxSlotAvailable\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minActiveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSlotLimitsFor\",\"inputs\":[{\"name\":\"n\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"maxSlotAvailable\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"minActiveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getState\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getStateCount\",\"inputs\":[{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumState\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSuspendedValidators\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getSuspender\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTimeouts\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"isActivated\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isConstructed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isUsedAddress\",\"inputs\":[{\"name\":\"addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"kirContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"offboard\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"pocContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"processSystemTransition\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"newStates\",\"type\":\"uint8[]\",\"internalType\":\"enumState[]\"},{\"name\":\"timeoutAts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"epochVACount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"readyCandidate\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"readyValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"requirement\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"resume\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeGcId\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"spareContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"suspendValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unreadyCandidate\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unreadyValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unsuspendValidator\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateCfsThreshold\",\"inputs\":[{\"name\":\"newCfsThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateConfigurator\",\"inputs\":[{\"name\":\"newConfigurator\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateIdleTimeout\",\"inputs\":[{\"name\":\"newIdleTimeout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKefAddress\",\"inputs\":[{\"name\":\"newKefAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKifAddress\",\"inputs\":[{\"name\":\"newKifAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateKpfAddress\",\"inputs\":[{\"name\":\"newKpfAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateManager\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newManager\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxCandReadyCount\",\"inputs\":[{\"name\":\"newMaxCandReadyCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxNodeCount\",\"inputs\":[{\"name\":\"newMaxNodeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMaxValActivePausedCount\",\"inputs\":[{\"name\":\"newMaxValActivePausedCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateMetadata\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newMetadata\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePauseTimeout\",\"inputs\":[{\"name\":\"newPauseTimeout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updatePfsThreshold\",\"inputs\":[{\"name\":\"newPfsThreshold\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateRewardAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newRewardAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateSuspender\",\"inputs\":[{\"name\":\"newSuspender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"updateVoterAddress\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"newVoterAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"AddressConfigUpdated\",\"inputs\":[{\"name\":\"configId\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"oldValue\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"newValue\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CandidateReadied\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CandidateUnreadied\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"EpochTransitionProcessed\",\"inputs\":[{\"name\":\"epochVACount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GcIdAssigned\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GcIdRevoked\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gcId\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ManagerUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newManager\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MetadataUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newMetadata\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeCreated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"NodeDeleted\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldRewardAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newRewardAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StateChanged\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"fromState\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumState\"},{\"name\":\"toState\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"enumState\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SystemTransitionProcessed\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"},{\"name\":\"newStates\",\"type\":\"uint8[]\",\"indexed\":false,\"internalType\":\"enumState[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UintConfigUpdated\",\"inputs\":[{\"name\":\"configId\",\"type\":\"uint8\",\"indexed\":true,\"internalType\":\"uint8\"},{\"name\":\"oldValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorSuspended\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorUnsuspended\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ValidatorsInitialized\",\"inputs\":[{\"name\":\"nodeIds\",\"type\":\"address[]\",\"indexed\":false,\"internalType\":\"address[]\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"VoterAddressUpdated\",\"inputs\":[{\"name\":\"nodeId\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"oldVoterAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newVoterAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressAlreadyRegistered\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"AlreadySuspended\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CnNodeNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FactoryNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GcIdAlreadyAssigned\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GcIdNotAssigned\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientNodeBalance\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInput\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"LegacyFunctionDeprecated\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeIdProofInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NodeNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSuspended\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyConfigurator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyNodeId\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySuspender\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlySystemTx\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"PDEnabled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SlotsFull\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakingDeployerMismatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StakingTooLow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimeoutExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60c060405230608052348015610013575f80fd5b50604051616149380380616149833981016040819052610032916100fb565b60a08190528080610041610049565b505050610112565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff16156100995760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b03908116146100f85780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b5f6020828403121561010b575f80fd5b5051919050565b60805160a051615ffa61014f5f395f818161078701528181613cc001526143f801525f8181613f3501528181613f5e015261412b0152615ffa5ff3fe608060405260043610610391575f3560e01c806378b84a5c116101de578063b858dd9511610108578063b858dd9514610ad3578063b9f96f4014610af2578063ba70d01814610b11578063be535f8b14610b30578063c732e08514610b4f578063c9a86af214610b6a578063cb1c2b5c14610b89578063cf8c6f5214610ba7578063d18c07ab14610bbb578063d267eda514610bda578063d3b5490714610bf9578063d9abb38b14610c0d578063da38d49814610c2c578063e4f0d37c14610c4b578063e59d7a8414610c6a578063e70c38f114610c89578063e8868e9f14610c9d578063f0a92ba814610cb2578063f2fde38b14610cc6578063ffa1ad7414610ce557610391565b806378b84a5c146108bd578063793c1946146108dc5780637df40c62146108fb5780638129fc1c1461091a57806387b7b8fd1461092e5780638da5cb5b146109485780638fabf3891461095c5780639b7ae5ec146109705780639d0e234d146109845780639d0f5ef1146109a35780639d8cf08f146109b75780639f9e3cba146109d6578063a41b6000146109f5578063a4c98ada14610a14578063a9ee547214610a33578063ad3cb1cc14610a52578063b42652e914610a82578063b57873a514610aa1578063b756393014610ac057610391565b8063453e962e116102bf578063453e962e14610649578063468e3a7e146106685780634a8c1fb4146106975780634b6a94cc146106b05780634f1ef286146106f357806350a5bb691461070657806350de2fb31461072457806352d1902d1461074357806353d39bfb14610757578063567b0b6c14610776578063582115fb146107a95780635b27b6c9146107d5578063656f5869146107f45780636968b53f146108135780636abd623d14610835578063715018a614610854578063715b208b14610868578063766718081461088a57806376a67a511461089e57610391565b806303e6689d146103b6578063058529fb146103e457806306bb84711461040357806307ecec3e146104245780630a4ff239146104435780630b1fe7841461046557806315575d5a14610486578063160370b8146104cf5780631865c57d146104f45780631b1a478b146105165780631b8f34ca146105355780631ba3fd581461055457806321d2320014610575578063229bb8231461059657806325cf0943146105c2578063291937f5146105d65780632aca5091146105ea5780632d4ede931461060b578063394f88991461062a575b34801561039c575f80fd5b50604051632053d6b560e11b815260040160405180910390fd5b3480156103c1575f80fd5b506103ca610cf9565b604080519283526020830191909152015b60405180910390f35b3480156103ef575f80fd5b506103ca6103fe366004614e5d565b610d19565b34801561040e575f80fd5b5061042261041d366004614e98565b610d36565b005b34801561042f575f80fd5b5061042261043e366004614eb3565b610f1c565b34801561044e575f80fd5b50610457610fb9565b6040519081526020016103db565b348015610470575f80fd5b50610479610fd2565b6040516103db9190614f2b565b348015610491575f80fd5b506104a56104a0366004614e98565b61111c565b604080516001600160a01b03948516815292841660208401529216918101919091526060016103db565b3480156104da575f80fd5b506104e3611160565b6040516103db959493929190614ff5565b3480156104ff575f80fd5b5061050861118a565b6040516103db929190615053565b348015610521575f80fd5b50610457610530366004615080565b6111eb565b348015610540575f80fd5b5061042261054f3660046150e2565b611230565b34801561055f575f80fd5b5061056861137d565b6040516103db919061517c565b348015610580575f80fd5b50610589611392565b6040516103db919061518e565b3480156105a1575f80fd5b506105b56105b0366004614e98565b6113ad565b6040516103db91906151a2565b3480156105cd575f80fd5b506104a56113da565b3480156105e1575f80fd5b5061045761140f565b3480156105f5575f80fd5b506105fe611421565b6040516103db91906151b0565b348015610616575f80fd5b50610422610625366004614e98565b611537565b348015610635575f80fd5b50610422610644366004614eb3565b6117c9565b348015610654575f80fd5b50610422610663366004614e98565b611973565b348015610673575f80fd5b50610687610682366004614e98565b611aa4565b60405190151581526020016103db565b3480156106a2575f80fd5b50600c546106879060ff1681565b3480156106bb575f80fd5b506106e66040518060400160405280600b81526020016a41646472657373426f6f6b60a81b81525081565b6040516103db9190615242565b61042261070136600461537e565b611ad1565b348015610711575f80fd5b50600c5461068790610100900460ff1681565b34801561072f575f80fd5b5061042261073e366004614e98565b611af0565b34801561074e575f80fd5b50610457611b39565b348015610762575f80fd5b50610422610771366004615434565b611b54565b348015610781575f80fd5b506104577f000000000000000000000000000000000000000000000000000000000000000081565b3480156107b4575f80fd5b506107c86107c3366004614e98565b611dea565b6040516103db9190615550565b3480156107e0575f80fd5b506104226107ef366004614e5d565b61212e565b3480156107ff575f80fd5b5061042261080e366004614e98565b612165565b34801561081e575f80fd5b506108276121da565b6040516103db92919061562b565b348015610840575f80fd5b50600754610589906001600160a01b031681565b34801561085f575f80fd5b5061042261247a565b348015610873575f80fd5b5061087c61248d565b6040516103db92919061569b565b348015610895575f80fd5b506104576127d9565b3480156108a9575f80fd5b506104226108b8366004614e98565b6127e2565b3480156108c8575f80fd5b506104226108d7366004614e98565b6128ae565b3480156108e7575f80fd5b506104226108f6366004614e98565b612922565b348015610906575f80fd5b50610422610915366004614e98565b612969565b348015610925575f80fd5b5061042261298c565b348015610939575f80fd5b506105896002600160a01b0381565b348015610953575f80fd5b50610589612ead565b348015610967575f80fd5b50610457612ec7565b34801561097b575f80fd5b50610589612ed9565b34801561098f575f80fd5b5061042261099e366004614e5d565b612ef4565b3480156109ae575f80fd5b506103ca612f16565b3480156109c2575f80fd5b506104226109d1366004614e98565b612f34565b3480156109e1575f80fd5b506104226109f0366004614eb3565b612f56565b348015610a00575f80fd5b50610422610a0f366004614e98565b6130be565b348015610a1f575f80fd5b50610422610a2e366004614e5d565b613132565b348015610a3e575f80fd5b50610422610a4d366004614e5d565b613155565b348015610a5d575f80fd5b506106e6604051806040016040528060058152602001640352e302e360dc1b81525081565b348015610a8d575f80fd5b50610422610a9c366004614e98565b613178565b348015610aac575f80fd5b50610422610abb366004614e98565b613281565b348015610acb575f80fd5b506001610457565b348015610ade575f80fd5b50600654610589906001600160a01b031681565b348015610afd575f80fd5b50610422610b0c366004614e98565b6132a4565b348015610b1c575f80fd5b50610422610b2b366004614e5d565b6132e2565b348015610b3b575f80fd5b50610422610b4a366004614e98565b613305565b348015610b5a575f80fd5b50610457678ac7230489e8000081565b348015610b75575f80fd5b50610422610b84366004614e98565b613467565b348015610b94575f80fd5b506104576a0422ca8b0a00a42500000081565b348015610bb2575f80fd5b5061056861348a565b348015610bc6575f80fd5b50610422610bd5366004614e5d565b61349f565b348015610be5575f80fd5b50600554610589906001600160a01b031681565b348015610c04575f80fd5b506104576134c2565b348015610c18575f80fd5b50610422610c27366004614e98565b6134d4565b348015610c37575f80fd5b50610422610c463660046156f5565b613515565b348015610c56575f80fd5b50610422610c65366004614e98565b6135af565b348015610c75575f80fd5b50610422610c84366004614e5d565b613624565b348015610c94575f80fd5b506103ca613647565b348015610ca8575f80fd5b5061045761080081565b348015610cbd575f80fd5b50610457613667565b348015610cd1575f80fd5b50610422610ce0366004614e98565b613679565b348015610cf0575f80fd5b50610457600281565b5f805f610d046136bf565b905080600e015481600f015492509250509091565b5f80610d24836136e3565b610d2d8461371c565b91509150915091565b610d3e613742565b5f610d476136bf565b90505f6001600160a01b0383165f908152602083905260409020600a015460ff166008811115610d7957610d79614ef7565b03610d9757604051634825e09360e01b815260040160405180910390fd5b6001600160a01b0382165f9081526020829052604090206005015415610dd057604051637be80ce960e11b815260040160405180910390fd5b5f816007015f8154610de190615786565b91829055506001600160a01b0384165f908152602084905260408082206005018390555163e2693e3f60e01b8152919250906104019063e2693e3f90610e299060040161579e565b602060405180830381865afa158015610e44573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e6891906157d1565b90506001600160a01b03811615610ed35760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f90610ea590879060040161518e565b5f604051808303815f87803b158015610ebc575f80fd5b505af1158015610ece573d5f803e3d5ffd5b505050505b836001600160a01b03167fe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb83604051610f0e91815260200190565b60405180910390a250505050565b81610f2681613776565b6001600160a01b038216610f4d5760405163b4fa3fb360e01b815260040160405180910390fd5b5f610f566136bf565b6001600160a01b038086165f8181526020849052604080822080548986166001600160a01b03198216811790925591519596509316938492917f8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda8791a45050505050565b5f610fcd610fc56136bf565b6001016137b9565b905090565b60605f610fdd6136bf565b90505f610fec826001016137b9565b9050806001600160401b0381111561100657611006615254565b60405190808252806020026020018201604052801561106557816020015b6110526040805160a0810182525f808252602082018190529181018290526060810182905290608082015290565b8152602001906001900390816110245790505b5092505f5b81811015611116575f61108060018501836137c2565b6001600160a01b038082165f8181526020888152604091829020825160a081018452938452600181015485169184019190915260028101549093169082015260048201546060820152600a8201549293509091608082019060ff1660088111156110ec576110ec614ef7565b815250868481518110611101576111016157ec565b6020908102919091010152505060010161106a565b50505090565b5f805f805f8061112b876137d4565b9250925092508061114f576040516342dc2dc560e01b815260040160405180910390fd5b5085945090925090505b9193909250565b60608060605f805f805f805f61117461384c565b939e929d50909b50995090975095505050505050565b6040805160018082528183019092526060915f918291602080830190803683370190505090506111b8612ead565b815f815181106111ca576111ca6157ec565b6001600160a01b039092166020928302919091019091015292600192509050565b5f6111f46136bf565b6008015f83600881111561120a5761120a614ef7565b600881111561121b5761121b614ef7565b81526020019081526020015f20549050919050565b611238613a30565b8584811415806112485750808314155b156112665760405163b4fa3fb360e01b815260040160405180910390fd5b5f5b818110156112e7576112df898983818110611285576112856157ec565b905060200201602081019061129a9190614e98565b8888848181106112ac576112ac6157ec565b90506020020160208101906112c19190615080565b8787858181106112d3576112d36157ec565b90506020020135613a57565b600101611268565b506112f0613cba565b1561133657816112fe6136bf565b601001556040518281527fd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e848160529060200160405180910390a15b7fab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e8888888860405161136b9493929190615800565b60405180910390a15050505050505050565b6060610fcd61138a6136bf565b600501613ceb565b5f61139b6136bf565b601501546001600160a01b0316919050565b5f6113b66136bf565b6001600160a01b039092165f9081526020929092525060409020600a015460ff1690565b5f805f806113e66136bf565b601281015460138201546014909201546001600160a01b03918216979282169650169350915050565b5f6114186136bf565b600a0154905090565b60605f61142c6136bf565b90505f61143b826001016137b9565b9050806001600160401b0381111561145557611455615254565b6040519080825280602002602001820160405280156114a557816020015b604080516080810182525f8082526020808301829052928201819052606082015282525f199092019101816114735790505b5092505f5b81811015611116575f6114c060018501836137c2565b6001600160a01b038082165f8181526020888152604091829020825160808101845293845260018101548516918401919091526003810154909316908201526005820154606082015287519293509091879085908110611522576115226157ec565b602090810291909101015250506001016114aa565b8061154181613776565b61154c826001613cf7565b6115695760405163baf3f0f760e01b815260040160405180910390fd5b5f6115726136bf565b6001600160a01b038085165f9081526020839052604090206003810154929350911615611678576003810180546001600160a01b031916905560405163e2693e3f60e01b81525f906104019063e2693e3f906115d09060040161579e565b602060405180830381865afa1580156115eb573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061160f91906157d1565b90506001600160a01b038116156116765760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f9061164c90889060040161518e565b5f604051808303815f87803b158015611663575f80fd5b505af1925050508015611674575060015b505b505b60018181015460028301546001600160a01b038781165f9081526009870160209081526040808320805460ff19908116909155958416835280832080548716905592909316815281812080549094169093559282526008850190529081208054916116e283615895565b909155506116f590506003830185613d2c565b506001600160a01b0384165f90815260208390526040812080546001600160a01b0319908116825560018201805482169055600282018054821690556003820180549091169055600481018290556005810182905590600682018161175a8282614d84565b611767600183015f614d84565b506117779050600883015f614d84565b611784600983015f614d84565b50600a01805460ff191690556040516001600160a01b038516907f1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad905f90a250505050565b816117d381613776565b6001600160a01b0382166117fa5760405163b4fa3fb360e01b815260040160405180910390fd5b5f6118036136bf565b6001600160a01b038086165f908152602083815260408083206001810154825163e1a12d3560e01b8152925196975090959394169263e1a12d35926004808401939192918290030181865afa15801561185e573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061188291906157d1565b6001600160a01b0316146118a957604051638ed87ef960e01b815260040160405180910390fd5b6118b284613d40565b6001600160a01b0384165f90815260098301602052604090205460ff16156118ed576040516316a163b960e11b815260040160405180910390fd5b6002810180546001600160a01b039081165f818152600986016020526040808220805460ff19908116909155898516808452828420805490921660011790915585546001600160a01b0319168117909555519193928492908a16917f270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc6970191a4505050505050565b8061197d81613dfb565b611988826001613cf7565b6119a55760405163baf3f0f760e01b815260040160405180910390fd5b6119ae82613e24565b6119cb5760405163bf74735560e01b815260040160405180910390fd5b678ac7230489e80000826001600160a01b03163110156119fd5760405162b8ec7b60e61b815260040160405180910390fd5b5f611a066136bf565b905080600f0154611a1760026111eb565b10611a355760405163848084dd60e01b815260040160405180910390fd5b80600e0154611a42610fb9565b10611a605760405163848084dd60e01b815260040160405180910390fd5b611a6c8360025f613a57565b6040516001600160a01b038416907fb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98905f90a2505050565b5f611aad6136bf565b6001600160a01b039092165f90815260099290920160205250604090205460ff1690565b611ad9613f2a565b611ae282613fce565b611aec8282613fd6565b5050565b611af861408e565b60035f80516020615fce833981519152611b136015846140c0565b604080516001600160a01b0392831681529185166020830152015b60405180910390a250565b5f611b42614120565b505f80516020615fae83398151915290565b5f611b5e896113ad565b6008811115611b6f57611b6f614ef7565b14611b8d5760405163731918fb60e11b815260040160405180910390fd5b82515f03611bae5760405163b4fa3fb360e01b815260040160405180910390fd5b61080082511115611bd25760405163b4fa3fb360e01b815260040160405180910390fd5b5f611bdb6136bf565b9050611bee816009018a8a8a8987614169565b5f604051806101400160405280336001600160a01b031681526020018a6001600160a01b03168152602001896001600160a01b03168152602001886001600160a01b031681526020015f81526020015f815260200187815260200186815260200185815260200160016008811115611c6857611c68614ef7565b90526001600160a01b03808c165f9081526020858152604091829020845181549085166001600160a01b031991821617825591850151600182018054918616918416919091179055918401516002830180549185169183169190911790556060840151600383018054919094169116179091556080820151600482015560a0820151600582015560c08201518051929350839260068301908190611d0c908261593a565b5060208201516001820190611d21908261593a565b50505060e08201516008820190611d38908261593a565b506101008201516009820190611d4e908261593a565b50610120820151600a8201805460ff19166001836008811115611d7357611d73614ef7565b0217905550611d88915050600383018b614300565b5060015f9081526008830160205260408120805491611da683615786565b90915550506040516001600160a01b038b16907f55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792905f90a250505050505050505050565b611df2614dbb565b5f611dfb6136bf565b6001600160a01b0384165f908152602091909152604081209150600a82015460ff166008811115611e2e57611e2e614ef7565b03611e4c57604051634825e09360e01b815260040160405180910390fd5b604080516101408101825282546001600160a01b0390811682526001840154811660208301526002840154811682840152600384015416606082015260048301546080820152600583015460a082015281518083019092526006830180549192849260c08501929082908290611ec1906158aa565b80601f0160208091040260200160405190810160405280929190818152602001828054611eed906158aa565b8015611f385780601f10611f0f57610100808354040283529160200191611f38565b820191905f5260205f20905b815481529060010190602001808311611f1b57829003601f168201915b50505050508152602001600182018054611f51906158aa565b80601f0160208091040260200160405190810160405280929190818152602001828054611f7d906158aa565b8015611fc85780601f10611f9f57610100808354040283529160200191611fc8565b820191905f5260205f20905b815481529060010190602001808311611fab57829003601f168201915b5050505050815250508152602001600882018054611fe5906158aa565b80601f0160208091040260200160405190810160405280929190818152602001828054612011906158aa565b801561205c5780601f106120335761010080835404028352916020019161205c565b820191905f5260205f20905b81548152906001019060200180831161203f57829003601f168201915b50505050508152602001600982018054612075906158aa565b80601f01602080910402602001604051908101604052809291908181526020018280546120a1906158aa565b80156120ec5780601f106120c3576101008083540402835291602001916120ec565b820191905f5260205f20905b8154815290600101906020018083116120cf57829003601f168201915b5050509183525050600a82015460209091019060ff16600881111561211357612113614ef7565b600881111561212457612124614ef7565b9052509392505050565b612136613742565b60055f80516020615f8e833981519152612151601184614314565b604080519182526020820185905201611b2e565b8061216f81613dfb565b61217a826004613cf7565b6121975760405163baf3f0f760e01b815260040160405180910390fd5b6121a082613e24565b6121bd5760405163bf74735560e01b815260040160405180910390fd5b6121c682614335565b611aec8260056121d58561436e565b613a57565b6060805f6121e66136bf565b90505f6121f5826001016137b9565b9050806001600160401b0381111561220f5761220f615254565b604051908082528060200260200182016040528015612238578160200160208202803683370190505b509350806001600160401b0381111561225357612253615254565b60405190808252806020026020018201604052801561229857816020015b60408051808201909152606080825260208201528152602001906001900390816122715790505b5092505f5b81811015612473576122b260018401826137c2565b8582815181106122c4576122c46157ec565b60200260200101906001600160a01b031690816001600160a01b031681525050825f015f8683815181106122fa576122fa6157ec565b60200260200101516001600160a01b03166001600160a01b031681526020019081526020015f206006016040518060400160405290815f8201805461233e906158aa565b80601f016020809104026020016040519081016040528092919081815260200182805461236a906158aa565b80156123b55780601f1061238c576101008083540402835291602001916123b5565b820191905f5260205f20905b81548152906001019060200180831161239857829003601f168201915b505050505081526020016001820180546123ce906158aa565b80601f01602080910402602001604051908101604052809291908181526020018280546123fa906158aa565b80156124455780601f1061241c57610100808354040283529160200191612445565b820191905f5260205f20905b81548152906001019060200180831161242857829003601f168201915b505050505081525050848281518110612460576124606157ec565b602090810291909101015260010161229d565b5050509091565b61248261408e565b61248b5f614398565b565b600c54606090819060ff166124b6575050604080515f8082526020820190815281830190925291565b5f805f805f6124c361384c565b8451949950929750909550935091506124dd8160036159f0565b6124e8906002615a07565b6001600160401b038111156124ff576124ff615254565b604051908082528060200260200182016040528015612528578160200160208202803683370190505b5097506125368160036159f0565b612541906002615a07565b6001600160401b0381111561255857612558615254565b604051908082528060200260200182016040528015612581578160200160208202803683370190505b5096505f805b8281101561270e575f8a83815181106125a2576125a26157ec565b602002602001019060ff16908160ff16815250508781815181106125c8576125c86157ec565b60200260200101518983806125dc90615786565b9450815181106125ee576125ee6157ec565b60200260200101906001600160a01b031690816001600160a01b03168152505060018a8381518110612622576126226157ec565b602002602001019060ff16908160ff1681525050868181518110612648576126486157ec565b602002602001015189838061265c90615786565b94508151811061266e5761266e6157ec565b60200260200101906001600160a01b031690816001600160a01b03168152505060028a83815181106126a2576126a26157ec565b602002602001019060ff16908160ff16815250508581815181106126c8576126c86157ec565b60200260200101518983806126dc90615786565b9450815181106126ee576126ee6157ec565b6001600160a01b0390921660209283029190910190910152600101612587565b506003898281518110612723576127236157ec565b60ff9092166020928302919091019091015283888261274181615786565b935081518110612753576127536157ec565b60200260200101906001600160a01b031690816001600160a01b0316815250506004898281518110612787576127876157ec565b602002602001019060ff16908160ff1681525050828882815181106127ae576127ae6157ec565b60200260200101906001600160a01b031690816001600160a01b031681525050505050505050509091565b5f610fcd6143f2565b806127ec81613dfb565b6127f7826006613cf7565b6128145760405163baf3f0f760e01b815260040160405180910390fd5b5f61281d6136bf565b905061282c81601001546136e3565b61283660076111eb565b106128545760405163848084dd60e01b815260040160405180910390fd5b612861816010015461371c565b61286b60066111eb565b116128895760405163848084dd60e01b815260040160405180910390fd5b5f81600c01544261289a9190615a07565b90506128a884600783613a57565b50505050565b6128b661441d565b5f6128bf6136bf565b90506128ce6005820183613d2c565b6128eb5760405163d33ff8c160e01b815260040160405180910390fd5b6040516001600160a01b038316907f814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e905f90a25050565b8061292c81613dfb565b612937826007613cf7565b6129545760405163baf3f0f760e01b815260040160405180910390fd5b61295d82614335565b611aec8260065f613a57565b612971613742565b60015f80516020615fce833981519152611b136013846140c0565b5f612995614451565b805490915060ff600160401b82041615906001600160401b03165f811580156129bb5750825b90505f826001600160401b031660011480156129d65750303b155b9050811580156129e4575080155b15612a025760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff191660011785558315612a2c57845460ff60401b1916600160401b1785555b60405163e2693e3f60e01b815260206004820152601060248201526f10509d8c91185d1850dbdb9d1c9858dd60821b60448201525f906104019063e2693e3f90606401602060405180830381865afa158015612a8a573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190612aae91906157d1565b90506001600160a01b038116612ad75760405163aed5959560e01b815260040160405180910390fd5b5f816001600160a01b031663ebe58ed76040518163ffffffff1660e01b81526004015f60405180830381865afa158015612b13573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f19168201604052612b3a9190810190615cd0565b90505f612b456136bf565b9050612b53825f0151614479565b6060820151600a8201556080820151600b82015560a0820151600c82015560c0820151600d82015560e0820151600e8201556101008201516011820155610120820151600f8201556101408201516012820180546001600160a01b03199081166001600160a01b03938416179091556101608401516013840180548316918416919091179055610180840151601484018054831691841691909117905560208401516015840180548316918416919091179055604084015160168401805490921692169190911790556101a0820151515f5b81811015612dfc575f846101a001518281518110612c4557612c456157ec565b602002602001015190505f856101c001518381518110612c6757612c676157ec565b6020908102919091018101516001600160a01b038481165f90815260098901845260408082208054600160ff199182168117909255958501518416835281832080548716821790558185015190931682529020805490931617909155905060066101208201819052505f608082018181526001600160a01b0380851683526020888152604093849020855181549084166001600160a01b03199182161782559186015160018201805491851691841691909117905593850151600285018054918416918316919091179055606085015160038501805491909316911617905551600482015560a0820151600582015560c082015180518392919060068301908190612d72908261593a565b5060208201516001820190612d87908261593a565b50505060e08201516008820190612d9e908261593a565b506101008201516009820190612db4908261593a565b50610120820151600a8201805460ff19166001836008811115612dd957612dd9614ef7565b0217905550612dee9150506001860183614300565b505050806001019050612c25565b5060065f9081526008830160205260409081902082905560108301829055606460078401556101a084015190517f820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d95291612e549161517c565b60405180910390a1505050508315612ea657845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b5f80612eb761448a565b546001600160a01b031692915050565b5f612ed06136bf565b60110154905090565b5f612ee26136bf565b601601546001600160a01b0316919050565b612efc613742565b5f5f80516020615f8e833981519152612151600c84614314565b5f80612f2c612f236136bf565b60100154610d19565b915091509091565b612f3c613742565b5f5f80516020615fce833981519152611b136012846140c0565b81612f6081613776565b5f612f696136bf565b6001600160a01b038086165f9081526020839052604080822060030180548885166001600160a01b0319821617909155905163e2693e3f60e01b8152939450909116916104019063e2693e3f90612fc29060040161579e565b602060405180830381865afa158015612fdd573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061300191906157d1565b90506001600160a01b0381161561306c5760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f9061303e90899060040161518e565b5f604051808303815f87803b158015613055575f80fd5b505af1158015613067573d5f803e3d5ffd5b505050505b846001600160a01b0316826001600160a01b0316876001600160a01b03167f23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d9660405160405180910390a4505050505050565b6130c661441d565b5f6130cf6136bf565b90506130de6005820183614300565b6130fb57604051633ad2b1bb60e11b815260040160405180910390fd5b6040516001600160a01b038316907fb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd905f90a25050565b61313a613742565b60045f80516020615f8e833981519152612151600e84614314565b61315d613742565b60065f80516020615f8e833981519152612151600f84614314565b8061318281613dfb565b5f61318c836113ad565b905060078160088111156131a2576131a2614ef7565b036131b5576131b083614335565b6131e7565b60068160088111156131c9576131c9614ef7565b146131e75760405163baf3f0f760e01b815260040160405180910390fd5b5f6131f06136bf565b90506131ff81601001546136e3565b61320960086111eb565b106132275760405163848084dd60e01b815260040160405180910390fd5b600682600881111561323b5761323b614ef7565b036132755761324d816010015461371c565b61325760066111eb565b116132755760405163848084dd60e01b815260040160405180910390fd5b6128a88460085f613a57565b61328961408e565b60045f80516020615fce833981519152611b136016846140c0565b806132ae81613dfb565b6132b9826004613cf7565b6132d65760405163baf3f0f760e01b815260040160405180910390fd5b611aec8260015f613a57565b6132ea613742565b60025f80516020615f8e833981519152612151600a84614314565b61330d613742565b5f6133166136bf565b6001600160a01b0383165f90815260208290526040812060050154919250819003613354576040516357024f6d60e11b815260040160405180910390fd5b60405163e2693e3f60e01b81525f906104019063e2693e3f906133799060040161579e565b602060405180830381865afa158015613394573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906133b891906157d1565b90506001600160a01b0381161561341b5760405163f575c5a760e01b8152600481018390526001600160a01b0382169063f575c5a7906024015f604051808303815f87803b158015613408575f80fd5b505af1925050508015613419575060015b505b6001600160a01b0384165f818152602085815260408083206005019290925590518481527f04078f2e4bb3259952b29491fa528a9a33d9496afca06a26bee6b7b1df26241f9101610f0e565b61346f613742565b60025f80516020615fce833981519152611b136014846140c0565b6060610fcd6134976136bf565b600301613ceb565b6134a7613742565b60035f80516020615f8e833981519152612151600b84614314565b5f6134cb6136bf565b60100154905090565b806134de81613dfb565b6134e9826005613cf7565b6135065760405163baf3f0f760e01b815260040160405180910390fd5b611aec8260046121d58561436e565b8261351f81613776565b6108008211156135425760405163b4fa3fb360e01b815260040160405180910390fd5b828261354c6136bf565b6001600160a01b0387165f9081526020919091526040902060090191613573919083615e14565b50836001600160a01b03167f2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf8484604051610f0e929190615ec8565b806135b981613dfb565b6135c4826002613cf7565b6135e15760405163baf3f0f760e01b815260040160405180910390fd5b6135ed8260015f613a57565b6040516001600160a01b038316907f8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2905f90a25050565b61362c613742565b60015f80516020615f8e833981519152612151600d84614314565b5f805f6136526136bf565b905080600c015481600d015492509250509091565b5f6136706136bf565b600b0154905090565b61368161408e565b6001600160a01b0381166136b3575f604051631e4fbdf760e01b81526004016136aa919061518e565b60405180910390fd5b6136bc81614398565b50565b7f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff0090565b5f60048210156136f457505f919050565b6002613701600384615f0a565b61370c906001615a07565b6137169190615f0a565b92915050565b5f600482101561372a575090565b60036137378360026159f0565b61370c906002615a07565b61374a6136bf565b601601546001600160a01b0316331461248b5760405163033b71e160e41b815260040160405180910390fd5b61377e6136bf565b6001600160a01b038281165f9081526020929092526040909120541633146136bc5760405163605919ad60e11b815260040160405180910390fd5b5f613716825490565b5f6137cd83836144ae565b9392505050565b5f805f806137e06136bf565b6001600160a01b0386165f908152602082905260408120919250600a82015460ff16600881111561381357613813614ef7565b03613828575f805f9450945094505050611159565b6001818101546002909201546001600160a01b039283169892169650945092505050565b60608060605f805f61385c6136bf565b90505f61386b826001016137b9565b9050806001600160401b0381111561388557613885615254565b6040519080825280602002602001820160405280156138ae578160200160208202803683370190505b509650806001600160401b038111156138c9576138c9615254565b6040519080825280602002602001820160405280156138f2578160200160208202803683370190505b509550806001600160401b0381111561390d5761390d615254565b604051908082528060200260200182016040528015613936578160200160208202803683370190505b5094505f5b81811015613a08575f61395160018501836137c2565b6001600160a01b0381165f9081526020869052604090208a519192509082908b9085908110613982576139826157ec565b6001600160a01b03928316602091820292909201015260018201548a519116908a90859081106139b4576139b46157ec565b6001600160a01b039283166020918202929092010152600282015489519116908990859081106139e6576139e66157ec565b6001600160a01b0390921660209283029190910190910152505060010161393b565b505060138101546012909101549596949593946001600160a01b039182169490911692509050565b336002600160a01b031461248b576040516354d325c360e01b815260040160405180910390fd5b5f613a606136bf565b6001600160a01b0385165f908152602082905260408120600a8101549293509160ff1690816008811115613a9657613a96614ef7565b1480613ab257505f856008811115613ab057613ab0614ef7565b145b80613ade5750846008811115613aca57613aca614ef7565b816008811115613adc57613adc614ef7565b145b15613aeb57505050505050565b826008015f826008811115613b0257613b02614ef7565b6008811115613b1357613b13614ef7565b81526020019081526020015f205f815480929190613b3090615895565b9190505550826008015f866008811115613b4c57613b4c614ef7565b6008811115613b5d57613b5d614ef7565b81526020019081526020015f205f815480929190613b7a90615786565b9091555060019050816008811115613b9457613b94614ef7565b148015613bb357506001856008811115613bb057613bb0614ef7565b14155b15613bd957613bc56001840187614300565b50613bd36003840187613d2c565b50613c2e565b6001816008811115613bed57613bed614ef7565b14158015613c0c57506001856008811115613c0a57613c0a614ef7565b145b15613c2e57613c1e6001840187613d2c565b50613c2c6003840187614300565b505b600a8201805486919060ff19166001836008811115613c4f57613c4f614ef7565b021790555060048201849055846008811115613c6d57613c6d614ef7565b816008811115613c7f57613c7f614ef7565b6040516001600160a01b038916907fcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee905f90a4505050505050565b5f613ce57f000000000000000000000000000000000000000000000000000000000000000043615f1d565b15919050565b60605f6137cd836144d4565b5f816008811115613d0a57613d0a614ef7565b613d13846113ad565b6008811115613d2457613d24614ef7565b149392505050565b5f6137cd836001600160a01b03841661452d565b5f613d49614617565b90506001600160a01b038116613d725760405163cdded31d60e01b815260040160405180910390fd5b604051631f7f8a5f60e21b81526001600160a01b03821690637dfe297c90613d9e90859060040161518e565b602060405180830381865afa158015613db9573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613ddd9190615f30565b15611aec5760405163b4fa3fb360e01b815260040160405180910390fd5b336001600160a01b038216146136bc576040516335f1334d60e11b815260040160405180910390fd5b5f80613e2e6136bf565b6001600160a01b038085165f9081526020928352604080822060010154815163318588a360e11b81529151931694509092849263630b11469260048082019392918290030181865afa158015613e86573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613eaa9190615f4f565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015613ee6573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613f0a9190615f4f565b613f149190615f66565b6a0422ca8b0a00a4250000001115949350505050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480613fb057507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316613fa45f80516020615fae833981519152546001600160a01b031690565b6001600160a01b031614155b1561248b5760405163703e46dd60e11b815260040160405180910390fd5b6136bc61408e565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015614030575060408051601f3d908101601f1916820190925261402d91810190615f4f565b60015b61404f5781604051634c9c8ce360e01b81526004016136aa919061518e565b5f80516020615fae833981519152811461407f57604051632a87526960e21b8152600481018290526024016136aa565b6140898383614699565b505050565b33614097612ead565b6001600160a01b03161461248b573360405163118cdaa760e01b81526004016136aa919061518e565b5f6001600160a01b0382166140e85760405163b4fa3fb360e01b815260040160405180910390fd5b5f614113847f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff00615a07565b8054939055509092915050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461248b5760405163703e46dd60e11b815260040160405180910390fd5b614175858585856146ee565b5f61417e614617565b90506001600160a01b0381166141a75760405163cdded31d60e01b815260040160405180910390fd5b60405163669d8d4560e01b815233906001600160a01b0383169063669d8d45906141d590899060040161518e565b602060405180830381865afa1580156141f0573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061421491906157d1565b6001600160a01b03161461423b57604051632281776f60e01b815260040160405180910390fd5b61424686868461485a565b5f6142518686614936565b6001600160a01b03161480156142cd5750604051631f7f8a5f60e21b81526001600160a01b03821690637dfe297c9061428e90879060040161518e565b602060405180830381865afa1580156142a9573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906142cd9190615f30565b156142eb5760405163b4fa3fb360e01b815260040160405180910390fd5b6142f7878787876149e1565b50505050505050565b5f6137cd836001600160a01b038416614aab565b5f815f036140e85760405163b4fa3fb360e01b815260040160405180910390fd5b5f61433f8261436e565b905080158015906143505750804210155b15611aec5760405163b48d5fc760e01b815260040160405180910390fd5b5f6143776136bf565b6001600160a01b039092165f90815260209290925250604090206004015490565b5f6143a161448a565b80546001600160a01b038481166001600160a01b031983168117845560405193945091169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f610fcd7f000000000000000000000000000000000000000000000000000000000000000043615f0a565b6144256136bf565b601501546001600160a01b0316331461248b5760405163333f4e6560e01b815260040160405180910390fd5b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00613716565b614481614af7565b6136bc81614b1c565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930090565b5f825f0182815481106144c3576144c36157ec565b905f5260205f200154905092915050565b6060815f0180548060200260200160405190810160405280929190818152602001828054801561452157602002820191905f5260205f20905b81548152602001906001019080831161450d575b50505050509050919050565b5f8181526001830160205260408120548015614607575f61454f600183615f66565b85549091505f9061456290600190615f66565b90508082146145c1575f865f018281548110614580576145806157ec565b905f5260205f200154905080875f0184815481106145a0576145a06157ec565b5f918252602080832090910192909255918252600188019052604090208390555b85548690806145d2576145d2615f79565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f905560019350505050613716565b5f915050613716565b5092915050565b60405163e2693e3f60e01b815260206004820152601060248201526f436e5374616b696e67466163746f727960801b60448201525f906104019063e2693e3f90606401602060405180830381865afa158015614675573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610fcd91906157d1565b6146a282614b24565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156146e6576140898282614b7e565b611aec614c0e565b6001600160a01b038416158061470b57506001600160a01b038316155b8061471d57506001600160a01b038216155b1561473b5760405163b4fa3fb360e01b815260040160405180910390fd5b826001600160a01b0316846001600160a01b0316148061476c5750816001600160a01b0316846001600160a01b0316145b806147885750816001600160a01b0316836001600160a01b0316145b156147a65760405163b4fa3fb360e01b815260040160405180910390fd5b80515160301415806147be5750806020015151606014155b156147dc5760405163b4fa3fb360e01b815260040160405180910390fd5b805180516020909101207fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd293148061483c575060208082015180519101207f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21145b156128a85760405163b4fa3fb360e01b815260040160405180910390fd5b604080517f23ae25c387ef8bd2c14b622e10202a494f464e31b796f40e83e9aecdf9cb42fb602082015246918101919091523060608201523360808201526001600160a01b0380851660a0830152831660c08201525f9060e0016040516020818303038152906040528051906020012090505f806148d88385614c2d565b5090925090505f8160038111156148f1576148f1614ef7565b1415806149105750856001600160a01b0316826001600160a01b031614155b1561492e57604051631ea9ff4d60e21b815260040160405180910390fd5b505050505050565b5f826001600160a01b031663e1a12d356040518163ffffffff1660e01b8152600401602060405180830381865afa158015614973573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061499791906157d1565b90506001600160a01b038116158015906149c35750816001600160a01b0316816001600160a01b031614155b156137165760405163b4fa3fb360e01b815260040160405180910390fd5b6001600160a01b0383165f9081526020859052604090205460ff1680614a1e57506001600160a01b0382165f9081526020859052604090205460ff165b80614a4057506001600160a01b0381165f9081526020859052604090205460ff165b15614a5e576040516316a163b960e11b815260040160405180910390fd5b6001600160a01b039283165f90815260209490945260408085208054600160ff19918216811790925593851686528186208054851682179055919093168452919092208054909216179055565b5f818152600183016020526040812054614af057508154600181810184555f848152602080822090930184905584548482528286019093526040902091909155613716565b505f613716565b614aff614c76565b61248b57604051631afcd79f60e31b815260040160405180910390fd5b613681614af7565b806001600160a01b03163b5f03614b505780604051634c9c8ce360e01b81526004016136aa919061518e565b5f80516020615fae83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b60605f614b8b8484614c8f565b9050808015614bac57505f3d1180614bac57505f846001600160a01b03163b115b15614bc157614bb9614ca2565b915050613716565b8015614be25783604051639996b31560e01b81526004016136aa919061518e565b3d15614bf557614bf0614cbb565b614610565b60405163d6bda27560e01b815260040160405180910390fd5b341561248b5760405163b398979f60e01b815260040160405180910390fd5b5f805f8351604103614c64576020840151604085015160608601515f1a614c5688828585614cc6565b955095509550505050614c6f565b505081515f91506002905b9250925092565b5f614c7f614451565b54600160401b900460ff16919050565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b5f80806fa2a8918ca85bafe22016d0b997e4df60600160ff1b03841115614cf557505f91506003905082614d7a565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015614d46573d5f803e3d5ffd5b5050604051601f1901519150506001600160a01b038116614d7157505f925060019150829050614d7a565b92505f91508190505b9450945094915050565b508054614d90906158aa565b5f825580601f10614d9f575050565b601f0160209004905f5260205f20908101906136bc9190614e45565b6040518061014001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f81526020015f8152602001614e2c604051806040016040528060608152602001606081525090565b815260606020820181905260408201819052015f905290565b5b80821115614e59575f8155600101614e46565b5090565b5f60208284031215614e6d575f80fd5b5035919050565b6001600160a01b03811681146136bc575f80fd5b8035614e9381614e74565b919050565b5f60208284031215614ea8575f80fd5b81356137cd81614e74565b5f8060408385031215614ec4575f80fd5b8235614ecf81614e74565b91506020830135614edf81614e74565b809150509250929050565b6001600160a01b03169052565b634e487b7160e01b5f52602160045260245ffd5b60098110614f2757634e487b7160e01b5f52602160045260245ffd5b9052565b602080825282518282018190525f919060409081850190868401855b82811015614fa557815180516001600160a01b039081168652878201518116888701528682015116868601526060808201519086015260809081015190614f9081870183614f0b565b505060a0939093019290850190600101614f47565b5091979650505050505050565b5f815180845260208085019450602084015f5b83811015614fea5781516001600160a01b031687529582019590820190600101614fc5565b509495945050505050565b60a081525f61500760a0830188614fb2565b82810360208401526150198188614fb2565b9050828103604084015261502d8187614fb2565b6001600160a01b0395861660608501529390941660809092019190915250949350505050565b604081525f6150656040830185614fb2565b90508260208301529392505050565b600981106136bc575f80fd5b5f60208284031215615090575f80fd5b81356137cd81615074565b5f8083601f8401126150ab575f80fd5b5081356001600160401b038111156150c1575f80fd5b6020830191508360208260051b85010111156150db575f80fd5b9250929050565b5f805f805f805f6080888a0312156150f8575f80fd5b87356001600160401b038082111561510e575f80fd5b61511a8b838c0161509b565b909950975060208a0135915080821115615132575f80fd5b61513e8b838c0161509b565b909750955060408a0135915080821115615156575f80fd5b506151638a828b0161509b565b989b979a50959894979596606090950135949350505050565b602081525f6137cd6020830184614fb2565b6001600160a01b0391909116815260200190565b602081016137168284614f0b565b602080825282518282018190525f919060409081850190868401855b82811015614fa557815180516001600160a01b0390811686528782015181168887015286820151168686015260609081015190850152608090930192908501906001016151cc565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f6137cd6020830184615214565b634e487b7160e01b5f52604160045260245ffd5b604080519081016001600160401b038111828210171561528a5761528a615254565b60405290565b60405161014081016001600160401b038111828210171561528a5761528a615254565b6040516101e081016001600160401b038111828210171561528a5761528a615254565b604051601f8201601f191681016001600160401b03811182821017156152fe576152fe615254565b604052919050565b5f6001600160401b0382111561531e5761531e615254565b50601f01601f191660200190565b5f82601f83011261533b575f80fd5b813561534e61534982615306565b6152d6565b818152846020838601011115615362575f80fd5b816020850160208301375f918101602001919091529392505050565b5f806040838503121561538f575f80fd5b823561539a81614e74565b915060208301356001600160401b038111156153b4575f80fd5b6153c08582860161532c565b9150509250929050565b5f604082840312156153da575f80fd5b6153e2615268565b905081356001600160401b03808211156153fa575f80fd5b6154068583860161532c565b8352602084013591508082111561541b575f80fd5b506154288482850161532c565b60208301525092915050565b5f805f805f805f80610100898b03121561544c575f80fd5b61545589614e88565b975061546360208a01614e88565b965061547160408a01614e88565b955061547f60608a01614e88565b945060808901356001600160401b038082111561549a575f80fd5b6154a68c838d016153ca565b955060a08b01359150808211156154bb575f80fd5b6154c78c838d0161532c565b945060c08b01359150808211156154dc575f80fd5b6154e88c838d0161532c565b935060e08b01359150808211156154fd575f80fd5b5061550a8b828c0161532c565b9150509295985092959890939650565b5f81516040845261552e6040850182615214565b9050602083015184820360208601526155478282615214565b95945050505050565b60208152615562602082018351614eea565b5f60208301516155756040840182614eea565b5060408301516155886060840182614eea565b50606083015161559b6080840182614eea565b50608083015160a083015260a083015160c083015260c08301516101408060e08501526155cc61016085018361551a565b915060e0850151601f196101008187860301818801526155ec8584615214565b94508088015192505061012081878603018188015261560b8584615214565b9450808801519250505061562182860182614f0b565b5090949350505050565b604081525f61563d6040830185614fb2565b6020838203818501528185518084528284019150828160051b8501018388015f5b8381101561568c57601f1987840301855261567a83835161551a565b9486019492509085019060010161565e565b50909998505050505050505050565b604080825283519082018190525f906020906060840190828701845b828110156156d657815160ff16845292840192908401906001016156b7565b50505083810360208501526156eb8186614fb2565b9695505050505050565b5f805f60408486031215615707575f80fd5b833561571281614e74565b925060208401356001600160401b038082111561572d575f80fd5b818601915086601f830112615740575f80fd5b81358181111561574e575f80fd5b87602082850101111561575f575f80fd5b6020830194508093505050509250925092565b634e487b7160e01b5f52601160045260245ffd5b5f6001820161579757615797615772565b5060010190565b6020808252600e908201526d29ba30b5b4b733aa3930b1b5b2b960911b604082015260600190565b8051614e9381614e74565b5f602082840312156157e1575f80fd5b81516137cd81614e74565b634e487b7160e01b5f52603260045260245ffd5b604080825281018490525f8560608301825b8781101561584257823561582581614e74565b6001600160a01b0316825260209283019290910190600101615812565b508381036020858101919091528582529150859082015f5b8681101561588857823561586d81615074565b6158778382614f0b565b50918301919083019060010161585a565b5098975050505050505050565b5f816158a3576158a3615772565b505f190190565b600181811c908216806158be57607f821691505b6020821081036158dc57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561408957805f5260205f20601f840160051c810160208510156159075750805b601f840160051c820191505b81811015612ea6575f8155600101615913565b5f19600383901b1c191660019190911b1790565b81516001600160401b0381111561595357615953615254565b6159678161596184546158aa565b846158e2565b602080601f831160018114615995575f84156159835750858301515b61598d8582615926565b86555061492e565b5f85815260208120601f198616915b828110156159c3578886015182559484019460019091019084016159a4565b50858210156159e057878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b808202811582820484141761371657613716615772565b8082018082111561371657613716615772565b5f6001600160401b03821115615a3257615a32615254565b5060051b60200190565b5f82601f830112615a4b575f80fd5b81516020615a5b61534983615a1a565b8083825260208201915060208460051b870101935086841115615a7c575f80fd5b602086015b84811015615aa1578051615a9481614e74565b8352918301918301615a81565b509695505050505050565b5f82601f830112615abb575f80fd5b8151615ac961534982615306565b818152846020838601011115615add575f80fd5b8160208501602083015e5f918101602001919091529392505050565b5f60408284031215615b09575f80fd5b615b11615268565b905081516001600160401b0380821115615b29575f80fd5b615b3585838601615aac565b83526020840151915080821115615b4a575f80fd5b5061542884828501615aac565b8051614e9381615074565b5f82601f830112615b71575f80fd5b81516020615b8161534983615a1a565b82815260059290921b84018101918181019086841115615b9f575f80fd5b8286015b84811015615aa15780516001600160401b0380821115615bc1575f80fd5b90880190610140828b03601f1901811315615bda575f80fd5b615be2615290565b615bed8885016157c6565b81526040615bfc8186016157c6565b898301526060615c0d8187016157c6565b8284015260809150615c208287016157c6565b818401525060a0808601518284015260c0915081860151818401525060e08086015185811115615c4e575f80fd5b615c5c8f8c838a0101615af9565b838501525061010091508186015185811115615c76575f80fd5b615c848f8c838a0101615aac565b8285015250506101208086015185811115615c9d575f80fd5b615cab8f8c838a0101615aac565b8385015250615cbb848701615b57565b90830152508652505050918301918301615ba3565b5f60208284031215615ce0575f80fd5b81516001600160401b0380821115615cf6575f80fd5b908301906101e08286031215615d0a575f80fd5b615d126152b3565b615d1b836157c6565b8152615d29602084016157c6565b6020820152615d3a604084016157c6565b6040820152606083015160608201526080830151608082015260a083015160a082015260c083015160c082015260e083015160e0820152610100808401518183015250610120808401518183015250610140615d978185016157c6565b90820152610160615da98482016157c6565b90820152610180615dbb8482016157c6565b908201526101a08381015183811115615dd2575f80fd5b615dde88828701615a3c565b8284015250506101c08084015183811115615df7575f80fd5b615e0388828701615b62565b918301919091525095945050505050565b6001600160401b03831115615e2b57615e2b615254565b615e3f83615e3983546158aa565b836158e2565b5f601f841160018114615e6b575f8515615e595750838201355b615e638682615926565b845550612ea6565b5f83815260208120601f198716915b82811015615e9a5786850135825560209485019460019092019101615e7a565b5086821015615eb6575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b60208152816020820152818360408301375f818301604090810191909152601f909201601f19160101919050565b634e487b7160e01b5f52601260045260245ffd5b5f82615f1857615f18615ef6565b500490565b5f82615f2b57615f2b615ef6565b500690565b5f60208284031215615f40575f80fd5b815180151581146137cd575f80fd5b5f60208284031215615f5f575f80fd5b5051919050565b8181038181111561371657613716615772565b634e487b7160e01b5f52603160045260245ffdfe34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbcd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4ba164736f6c6343000819000a",
}

// AddressBookV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use AddressBookV2MetaData.ABI instead.
var AddressBookV2ABI = AddressBookV2MetaData.ABI

// AddressBookV2BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressBookV2BinRuntime = `608060405260043610610391575f3560e01c806378b84a5c116101de578063b858dd9511610108578063b858dd9514610ad3578063b9f96f4014610af2578063ba70d01814610b11578063be535f8b14610b30578063c732e08514610b4f578063c9a86af214610b6a578063cb1c2b5c14610b89578063cf8c6f5214610ba7578063d18c07ab14610bbb578063d267eda514610bda578063d3b5490714610bf9578063d9abb38b14610c0d578063da38d49814610c2c578063e4f0d37c14610c4b578063e59d7a8414610c6a578063e70c38f114610c89578063e8868e9f14610c9d578063f0a92ba814610cb2578063f2fde38b14610cc6578063ffa1ad7414610ce557610391565b806378b84a5c146108bd578063793c1946146108dc5780637df40c62146108fb5780638129fc1c1461091a57806387b7b8fd1461092e5780638da5cb5b146109485780638fabf3891461095c5780639b7ae5ec146109705780639d0e234d146109845780639d0f5ef1146109a35780639d8cf08f146109b75780639f9e3cba146109d6578063a41b6000146109f5578063a4c98ada14610a14578063a9ee547214610a33578063ad3cb1cc14610a52578063b42652e914610a82578063b57873a514610aa1578063b756393014610ac057610391565b8063453e962e116102bf578063453e962e14610649578063468e3a7e146106685780634a8c1fb4146106975780634b6a94cc146106b05780634f1ef286146106f357806350a5bb691461070657806350de2fb31461072457806352d1902d1461074357806353d39bfb14610757578063567b0b6c14610776578063582115fb146107a95780635b27b6c9146107d5578063656f5869146107f45780636968b53f146108135780636abd623d14610835578063715018a614610854578063715b208b14610868578063766718081461088a57806376a67a511461089e57610391565b806303e6689d146103b6578063058529fb146103e457806306bb84711461040357806307ecec3e146104245780630a4ff239146104435780630b1fe7841461046557806315575d5a14610486578063160370b8146104cf5780631865c57d146104f45780631b1a478b146105165780631b8f34ca146105355780631ba3fd581461055457806321d2320014610575578063229bb8231461059657806325cf0943146105c2578063291937f5146105d65780632aca5091146105ea5780632d4ede931461060b578063394f88991461062a575b34801561039c575f80fd5b50604051632053d6b560e11b815260040160405180910390fd5b3480156103c1575f80fd5b506103ca610cf9565b604080519283526020830191909152015b60405180910390f35b3480156103ef575f80fd5b506103ca6103fe366004614e5d565b610d19565b34801561040e575f80fd5b5061042261041d366004614e98565b610d36565b005b34801561042f575f80fd5b5061042261043e366004614eb3565b610f1c565b34801561044e575f80fd5b50610457610fb9565b6040519081526020016103db565b348015610470575f80fd5b50610479610fd2565b6040516103db9190614f2b565b348015610491575f80fd5b506104a56104a0366004614e98565b61111c565b604080516001600160a01b03948516815292841660208401529216918101919091526060016103db565b3480156104da575f80fd5b506104e3611160565b6040516103db959493929190614ff5565b3480156104ff575f80fd5b5061050861118a565b6040516103db929190615053565b348015610521575f80fd5b50610457610530366004615080565b6111eb565b348015610540575f80fd5b5061042261054f3660046150e2565b611230565b34801561055f575f80fd5b5061056861137d565b6040516103db919061517c565b348015610580575f80fd5b50610589611392565b6040516103db919061518e565b3480156105a1575f80fd5b506105b56105b0366004614e98565b6113ad565b6040516103db91906151a2565b3480156105cd575f80fd5b506104a56113da565b3480156105e1575f80fd5b5061045761140f565b3480156105f5575f80fd5b506105fe611421565b6040516103db91906151b0565b348015610616575f80fd5b50610422610625366004614e98565b611537565b348015610635575f80fd5b50610422610644366004614eb3565b6117c9565b348015610654575f80fd5b50610422610663366004614e98565b611973565b348015610673575f80fd5b50610687610682366004614e98565b611aa4565b60405190151581526020016103db565b3480156106a2575f80fd5b50600c546106879060ff1681565b3480156106bb575f80fd5b506106e66040518060400160405280600b81526020016a41646472657373426f6f6b60a81b81525081565b6040516103db9190615242565b61042261070136600461537e565b611ad1565b348015610711575f80fd5b50600c5461068790610100900460ff1681565b34801561072f575f80fd5b5061042261073e366004614e98565b611af0565b34801561074e575f80fd5b50610457611b39565b348015610762575f80fd5b50610422610771366004615434565b611b54565b348015610781575f80fd5b506104577f000000000000000000000000000000000000000000000000000000000000000081565b3480156107b4575f80fd5b506107c86107c3366004614e98565b611dea565b6040516103db9190615550565b3480156107e0575f80fd5b506104226107ef366004614e5d565b61212e565b3480156107ff575f80fd5b5061042261080e366004614e98565b612165565b34801561081e575f80fd5b506108276121da565b6040516103db92919061562b565b348015610840575f80fd5b50600754610589906001600160a01b031681565b34801561085f575f80fd5b5061042261247a565b348015610873575f80fd5b5061087c61248d565b6040516103db92919061569b565b348015610895575f80fd5b506104576127d9565b3480156108a9575f80fd5b506104226108b8366004614e98565b6127e2565b3480156108c8575f80fd5b506104226108d7366004614e98565b6128ae565b3480156108e7575f80fd5b506104226108f6366004614e98565b612922565b348015610906575f80fd5b50610422610915366004614e98565b612969565b348015610925575f80fd5b5061042261298c565b348015610939575f80fd5b506105896002600160a01b0381565b348015610953575f80fd5b50610589612ead565b348015610967575f80fd5b50610457612ec7565b34801561097b575f80fd5b50610589612ed9565b34801561098f575f80fd5b5061042261099e366004614e5d565b612ef4565b3480156109ae575f80fd5b506103ca612f16565b3480156109c2575f80fd5b506104226109d1366004614e98565b612f34565b3480156109e1575f80fd5b506104226109f0366004614eb3565b612f56565b348015610a00575f80fd5b50610422610a0f366004614e98565b6130be565b348015610a1f575f80fd5b50610422610a2e366004614e5d565b613132565b348015610a3e575f80fd5b50610422610a4d366004614e5d565b613155565b348015610a5d575f80fd5b506106e6604051806040016040528060058152602001640352e302e360dc1b81525081565b348015610a8d575f80fd5b50610422610a9c366004614e98565b613178565b348015610aac575f80fd5b50610422610abb366004614e98565b613281565b348015610acb575f80fd5b506001610457565b348015610ade575f80fd5b50600654610589906001600160a01b031681565b348015610afd575f80fd5b50610422610b0c366004614e98565b6132a4565b348015610b1c575f80fd5b50610422610b2b366004614e5d565b6132e2565b348015610b3b575f80fd5b50610422610b4a366004614e98565b613305565b348015610b5a575f80fd5b50610457678ac7230489e8000081565b348015610b75575f80fd5b50610422610b84366004614e98565b613467565b348015610b94575f80fd5b506104576a0422ca8b0a00a42500000081565b348015610bb2575f80fd5b5061056861348a565b348015610bc6575f80fd5b50610422610bd5366004614e5d565b61349f565b348015610be5575f80fd5b50600554610589906001600160a01b031681565b348015610c04575f80fd5b506104576134c2565b348015610c18575f80fd5b50610422610c27366004614e98565b6134d4565b348015610c37575f80fd5b50610422610c463660046156f5565b613515565b348015610c56575f80fd5b50610422610c65366004614e98565b6135af565b348015610c75575f80fd5b50610422610c84366004614e5d565b613624565b348015610c94575f80fd5b506103ca613647565b348015610ca8575f80fd5b5061045761080081565b348015610cbd575f80fd5b50610457613667565b348015610cd1575f80fd5b50610422610ce0366004614e98565b613679565b348015610cf0575f80fd5b50610457600281565b5f805f610d046136bf565b905080600e015481600f015492509250509091565b5f80610d24836136e3565b610d2d8461371c565b91509150915091565b610d3e613742565b5f610d476136bf565b90505f6001600160a01b0383165f908152602083905260409020600a015460ff166008811115610d7957610d79614ef7565b03610d9757604051634825e09360e01b815260040160405180910390fd5b6001600160a01b0382165f9081526020829052604090206005015415610dd057604051637be80ce960e11b815260040160405180910390fd5b5f816007015f8154610de190615786565b91829055506001600160a01b0384165f908152602084905260408082206005018390555163e2693e3f60e01b8152919250906104019063e2693e3f90610e299060040161579e565b602060405180830381865afa158015610e44573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610e6891906157d1565b90506001600160a01b03811615610ed35760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f90610ea590879060040161518e565b5f604051808303815f87803b158015610ebc575f80fd5b505af1158015610ece573d5f803e3d5ffd5b505050505b836001600160a01b03167fe1fbe15fca2fbb149763b54900ac143ffd56dbbc787c6bfd0e3d45fae47e01eb83604051610f0e91815260200190565b60405180910390a250505050565b81610f2681613776565b6001600160a01b038216610f4d5760405163b4fa3fb360e01b815260040160405180910390fd5b5f610f566136bf565b6001600160a01b038086165f8181526020849052604080822080548986166001600160a01b03198216811790925591519596509316938492917f8df26d30992ecfde135bbe59c1f267d82e2aae9d32fdae41551a38fe8b7bda8791a45050505050565b5f610fcd610fc56136bf565b6001016137b9565b905090565b60605f610fdd6136bf565b90505f610fec826001016137b9565b9050806001600160401b0381111561100657611006615254565b60405190808252806020026020018201604052801561106557816020015b6110526040805160a0810182525f808252602082018190529181018290526060810182905290608082015290565b8152602001906001900390816110245790505b5092505f5b81811015611116575f61108060018501836137c2565b6001600160a01b038082165f8181526020888152604091829020825160a081018452938452600181015485169184019190915260028101549093169082015260048201546060820152600a8201549293509091608082019060ff1660088111156110ec576110ec614ef7565b815250868481518110611101576111016157ec565b6020908102919091010152505060010161106a565b50505090565b5f805f805f8061112b876137d4565b9250925092508061114f576040516342dc2dc560e01b815260040160405180910390fd5b5085945090925090505b9193909250565b60608060605f805f805f805f61117461384c565b939e929d50909b50995090975095505050505050565b6040805160018082528183019092526060915f918291602080830190803683370190505090506111b8612ead565b815f815181106111ca576111ca6157ec565b6001600160a01b039092166020928302919091019091015292600192509050565b5f6111f46136bf565b6008015f83600881111561120a5761120a614ef7565b600881111561121b5761121b614ef7565b81526020019081526020015f20549050919050565b611238613a30565b8584811415806112485750808314155b156112665760405163b4fa3fb360e01b815260040160405180910390fd5b5f5b818110156112e7576112df898983818110611285576112856157ec565b905060200201602081019061129a9190614e98565b8888848181106112ac576112ac6157ec565b90506020020160208101906112c19190615080565b8787858181106112d3576112d36157ec565b90506020020135613a57565b600101611268565b506112f0613cba565b1561133657816112fe6136bf565b601001556040518281527fd45be950fd3aceb65c6059b131cc8e06ab2390da6780d464b82c153e848160529060200160405180910390a15b7fab95e7867bd336dde387ba31a71307c75dcc78b0344b873a5e993eb4470eb37e8888888860405161136b9493929190615800565b60405180910390a15050505050505050565b6060610fcd61138a6136bf565b600501613ceb565b5f61139b6136bf565b601501546001600160a01b0316919050565b5f6113b66136bf565b6001600160a01b039092165f9081526020929092525060409020600a015460ff1690565b5f805f806113e66136bf565b601281015460138201546014909201546001600160a01b03918216979282169650169350915050565b5f6114186136bf565b600a0154905090565b60605f61142c6136bf565b90505f61143b826001016137b9565b9050806001600160401b0381111561145557611455615254565b6040519080825280602002602001820160405280156114a557816020015b604080516080810182525f8082526020808301829052928201819052606082015282525f199092019101816114735790505b5092505f5b81811015611116575f6114c060018501836137c2565b6001600160a01b038082165f8181526020888152604091829020825160808101845293845260018101548516918401919091526003810154909316908201526005820154606082015287519293509091879085908110611522576115226157ec565b602090810291909101015250506001016114aa565b8061154181613776565b61154c826001613cf7565b6115695760405163baf3f0f760e01b815260040160405180910390fd5b5f6115726136bf565b6001600160a01b038085165f9081526020839052604090206003810154929350911615611678576003810180546001600160a01b031916905560405163e2693e3f60e01b81525f906104019063e2693e3f906115d09060040161579e565b602060405180830381865afa1580156115eb573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061160f91906157d1565b90506001600160a01b038116156116765760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f9061164c90889060040161518e565b5f604051808303815f87803b158015611663575f80fd5b505af1925050508015611674575060015b505b505b60018181015460028301546001600160a01b038781165f9081526009870160209081526040808320805460ff19908116909155958416835280832080548716905592909316815281812080549094169093559282526008850190529081208054916116e283615895565b909155506116f590506003830185613d2c565b506001600160a01b0384165f90815260208390526040812080546001600160a01b0319908116825560018201805482169055600282018054821690556003820180549091169055600481018290556005810182905590600682018161175a8282614d84565b611767600183015f614d84565b506117779050600883015f614d84565b611784600983015f614d84565b50600a01805460ff191690556040516001600160a01b038516907f1629bfc36423a1b4749d3fe1d6970b9d32d42bbee47dd5540670696ab6b9a4ad905f90a250505050565b816117d381613776565b6001600160a01b0382166117fa5760405163b4fa3fb360e01b815260040160405180910390fd5b5f6118036136bf565b6001600160a01b038086165f908152602083815260408083206001810154825163e1a12d3560e01b8152925196975090959394169263e1a12d35926004808401939192918290030181865afa15801561185e573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061188291906157d1565b6001600160a01b0316146118a957604051638ed87ef960e01b815260040160405180910390fd5b6118b284613d40565b6001600160a01b0384165f90815260098301602052604090205460ff16156118ed576040516316a163b960e11b815260040160405180910390fd5b6002810180546001600160a01b039081165f818152600986016020526040808220805460ff19908116909155898516808452828420805490921660011790915585546001600160a01b0319168117909555519193928492908a16917f270e800343b82239558a49df43a4ab4ec495dbfd29f864df4fbd9b927dc6970191a4505050505050565b8061197d81613dfb565b611988826001613cf7565b6119a55760405163baf3f0f760e01b815260040160405180910390fd5b6119ae82613e24565b6119cb5760405163bf74735560e01b815260040160405180910390fd5b678ac7230489e80000826001600160a01b03163110156119fd5760405162b8ec7b60e61b815260040160405180910390fd5b5f611a066136bf565b905080600f0154611a1760026111eb565b10611a355760405163848084dd60e01b815260040160405180910390fd5b80600e0154611a42610fb9565b10611a605760405163848084dd60e01b815260040160405180910390fd5b611a6c8360025f613a57565b6040516001600160a01b038416907fb6cfd7c953a120707430bb9a474b9062b3dd92baab50f0c69ea822b324a31b98905f90a2505050565b5f611aad6136bf565b6001600160a01b039092165f90815260099290920160205250604090205460ff1690565b611ad9613f2a565b611ae282613fce565b611aec8282613fd6565b5050565b611af861408e565b60035f80516020615fce833981519152611b136015846140c0565b604080516001600160a01b0392831681529185166020830152015b60405180910390a250565b5f611b42614120565b505f80516020615fae83398151915290565b5f611b5e896113ad565b6008811115611b6f57611b6f614ef7565b14611b8d5760405163731918fb60e11b815260040160405180910390fd5b82515f03611bae5760405163b4fa3fb360e01b815260040160405180910390fd5b61080082511115611bd25760405163b4fa3fb360e01b815260040160405180910390fd5b5f611bdb6136bf565b9050611bee816009018a8a8a8987614169565b5f604051806101400160405280336001600160a01b031681526020018a6001600160a01b03168152602001896001600160a01b03168152602001886001600160a01b031681526020015f81526020015f815260200187815260200186815260200185815260200160016008811115611c6857611c68614ef7565b90526001600160a01b03808c165f9081526020858152604091829020845181549085166001600160a01b031991821617825591850151600182018054918616918416919091179055918401516002830180549185169183169190911790556060840151600383018054919094169116179091556080820151600482015560a0820151600582015560c08201518051929350839260068301908190611d0c908261593a565b5060208201516001820190611d21908261593a565b50505060e08201516008820190611d38908261593a565b506101008201516009820190611d4e908261593a565b50610120820151600a8201805460ff19166001836008811115611d7357611d73614ef7565b0217905550611d88915050600383018b614300565b5060015f9081526008830160205260408120805491611da683615786565b90915550506040516001600160a01b038b16907f55fdf3ae96916cdb0bf329ba2d19e0618b01d8f4d6cfe27ec8bbb79c62be7792905f90a250505050505050505050565b611df2614dbb565b5f611dfb6136bf565b6001600160a01b0384165f908152602091909152604081209150600a82015460ff166008811115611e2e57611e2e614ef7565b03611e4c57604051634825e09360e01b815260040160405180910390fd5b604080516101408101825282546001600160a01b0390811682526001840154811660208301526002840154811682840152600384015416606082015260048301546080820152600583015460a082015281518083019092526006830180549192849260c08501929082908290611ec1906158aa565b80601f0160208091040260200160405190810160405280929190818152602001828054611eed906158aa565b8015611f385780601f10611f0f57610100808354040283529160200191611f38565b820191905f5260205f20905b815481529060010190602001808311611f1b57829003601f168201915b50505050508152602001600182018054611f51906158aa565b80601f0160208091040260200160405190810160405280929190818152602001828054611f7d906158aa565b8015611fc85780601f10611f9f57610100808354040283529160200191611fc8565b820191905f5260205f20905b815481529060010190602001808311611fab57829003601f168201915b5050505050815250508152602001600882018054611fe5906158aa565b80601f0160208091040260200160405190810160405280929190818152602001828054612011906158aa565b801561205c5780601f106120335761010080835404028352916020019161205c565b820191905f5260205f20905b81548152906001019060200180831161203f57829003601f168201915b50505050508152602001600982018054612075906158aa565b80601f01602080910402602001604051908101604052809291908181526020018280546120a1906158aa565b80156120ec5780601f106120c3576101008083540402835291602001916120ec565b820191905f5260205f20905b8154815290600101906020018083116120cf57829003601f168201915b5050509183525050600a82015460209091019060ff16600881111561211357612113614ef7565b600881111561212457612124614ef7565b9052509392505050565b612136613742565b60055f80516020615f8e833981519152612151601184614314565b604080519182526020820185905201611b2e565b8061216f81613dfb565b61217a826004613cf7565b6121975760405163baf3f0f760e01b815260040160405180910390fd5b6121a082613e24565b6121bd5760405163bf74735560e01b815260040160405180910390fd5b6121c682614335565b611aec8260056121d58561436e565b613a57565b6060805f6121e66136bf565b90505f6121f5826001016137b9565b9050806001600160401b0381111561220f5761220f615254565b604051908082528060200260200182016040528015612238578160200160208202803683370190505b509350806001600160401b0381111561225357612253615254565b60405190808252806020026020018201604052801561229857816020015b60408051808201909152606080825260208201528152602001906001900390816122715790505b5092505f5b81811015612473576122b260018401826137c2565b8582815181106122c4576122c46157ec565b60200260200101906001600160a01b031690816001600160a01b031681525050825f015f8683815181106122fa576122fa6157ec565b60200260200101516001600160a01b03166001600160a01b031681526020019081526020015f206006016040518060400160405290815f8201805461233e906158aa565b80601f016020809104026020016040519081016040528092919081815260200182805461236a906158aa565b80156123b55780601f1061238c576101008083540402835291602001916123b5565b820191905f5260205f20905b81548152906001019060200180831161239857829003601f168201915b505050505081526020016001820180546123ce906158aa565b80601f01602080910402602001604051908101604052809291908181526020018280546123fa906158aa565b80156124455780601f1061241c57610100808354040283529160200191612445565b820191905f5260205f20905b81548152906001019060200180831161242857829003601f168201915b505050505081525050848281518110612460576124606157ec565b602090810291909101015260010161229d565b5050509091565b61248261408e565b61248b5f614398565b565b600c54606090819060ff166124b6575050604080515f8082526020820190815281830190925291565b5f805f805f6124c361384c565b8451949950929750909550935091506124dd8160036159f0565b6124e8906002615a07565b6001600160401b038111156124ff576124ff615254565b604051908082528060200260200182016040528015612528578160200160208202803683370190505b5097506125368160036159f0565b612541906002615a07565b6001600160401b0381111561255857612558615254565b604051908082528060200260200182016040528015612581578160200160208202803683370190505b5096505f805b8281101561270e575f8a83815181106125a2576125a26157ec565b602002602001019060ff16908160ff16815250508781815181106125c8576125c86157ec565b60200260200101518983806125dc90615786565b9450815181106125ee576125ee6157ec565b60200260200101906001600160a01b031690816001600160a01b03168152505060018a8381518110612622576126226157ec565b602002602001019060ff16908160ff1681525050868181518110612648576126486157ec565b602002602001015189838061265c90615786565b94508151811061266e5761266e6157ec565b60200260200101906001600160a01b031690816001600160a01b03168152505060028a83815181106126a2576126a26157ec565b602002602001019060ff16908160ff16815250508581815181106126c8576126c86157ec565b60200260200101518983806126dc90615786565b9450815181106126ee576126ee6157ec565b6001600160a01b0390921660209283029190910190910152600101612587565b506003898281518110612723576127236157ec565b60ff9092166020928302919091019091015283888261274181615786565b935081518110612753576127536157ec565b60200260200101906001600160a01b031690816001600160a01b0316815250506004898281518110612787576127876157ec565b602002602001019060ff16908160ff1681525050828882815181106127ae576127ae6157ec565b60200260200101906001600160a01b031690816001600160a01b031681525050505050505050509091565b5f610fcd6143f2565b806127ec81613dfb565b6127f7826006613cf7565b6128145760405163baf3f0f760e01b815260040160405180910390fd5b5f61281d6136bf565b905061282c81601001546136e3565b61283660076111eb565b106128545760405163848084dd60e01b815260040160405180910390fd5b612861816010015461371c565b61286b60066111eb565b116128895760405163848084dd60e01b815260040160405180910390fd5b5f81600c01544261289a9190615a07565b90506128a884600783613a57565b50505050565b6128b661441d565b5f6128bf6136bf565b90506128ce6005820183613d2c565b6128eb5760405163d33ff8c160e01b815260040160405180910390fd5b6040516001600160a01b038316907f814c4b6f6fc147ebb6fbe4ffcd3554d0309170fd0a70e66cc4e4c0784f4aa32e905f90a25050565b8061292c81613dfb565b612937826007613cf7565b6129545760405163baf3f0f760e01b815260040160405180910390fd5b61295d82614335565b611aec8260065f613a57565b612971613742565b60015f80516020615fce833981519152611b136013846140c0565b5f612995614451565b805490915060ff600160401b82041615906001600160401b03165f811580156129bb5750825b90505f826001600160401b031660011480156129d65750303b155b9050811580156129e4575080155b15612a025760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff191660011785558315612a2c57845460ff60401b1916600160401b1785555b60405163e2693e3f60e01b815260206004820152601060248201526f10509d8c91185d1850dbdb9d1c9858dd60821b60448201525f906104019063e2693e3f90606401602060405180830381865afa158015612a8a573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190612aae91906157d1565b90506001600160a01b038116612ad75760405163aed5959560e01b815260040160405180910390fd5b5f816001600160a01b031663ebe58ed76040518163ffffffff1660e01b81526004015f60405180830381865afa158015612b13573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f19168201604052612b3a9190810190615cd0565b90505f612b456136bf565b9050612b53825f0151614479565b6060820151600a8201556080820151600b82015560a0820151600c82015560c0820151600d82015560e0820151600e8201556101008201516011820155610120820151600f8201556101408201516012820180546001600160a01b03199081166001600160a01b03938416179091556101608401516013840180548316918416919091179055610180840151601484018054831691841691909117905560208401516015840180548316918416919091179055604084015160168401805490921692169190911790556101a0820151515f5b81811015612dfc575f846101a001518281518110612c4557612c456157ec565b602002602001015190505f856101c001518381518110612c6757612c676157ec565b6020908102919091018101516001600160a01b038481165f90815260098901845260408082208054600160ff199182168117909255958501518416835281832080548716821790558185015190931682529020805490931617909155905060066101208201819052505f608082018181526001600160a01b0380851683526020888152604093849020855181549084166001600160a01b03199182161782559186015160018201805491851691841691909117905593850151600285018054918416918316919091179055606085015160038501805491909316911617905551600482015560a0820151600582015560c082015180518392919060068301908190612d72908261593a565b5060208201516001820190612d87908261593a565b50505060e08201516008820190612d9e908261593a565b506101008201516009820190612db4908261593a565b50610120820151600a8201805460ff19166001836008811115612dd957612dd9614ef7565b0217905550612dee9150506001860183614300565b505050806001019050612c25565b5060065f9081526008830160205260409081902082905560108301829055606460078401556101a084015190517f820f68b9d060f5d911b3243881ada086c3768ea90e97a10f7f5023d84b94d95291612e549161517c565b60405180910390a1505050508315612ea657845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b5050505050565b5f80612eb761448a565b546001600160a01b031692915050565b5f612ed06136bf565b60110154905090565b5f612ee26136bf565b601601546001600160a01b0316919050565b612efc613742565b5f5f80516020615f8e833981519152612151600c84614314565b5f80612f2c612f236136bf565b60100154610d19565b915091509091565b612f3c613742565b5f5f80516020615fce833981519152611b136012846140c0565b81612f6081613776565b5f612f696136bf565b6001600160a01b038086165f9081526020839052604080822060030180548885166001600160a01b0319821617909155905163e2693e3f60e01b8152939450909116916104019063e2693e3f90612fc29060040161579e565b602060405180830381865afa158015612fdd573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061300191906157d1565b90506001600160a01b0381161561306c5760405163aad8cb3f60e01b81526001600160a01b0382169063aad8cb3f9061303e90899060040161518e565b5f604051808303815f87803b158015613055575f80fd5b505af1158015613067573d5f803e3d5ffd5b505050505b846001600160a01b0316826001600160a01b0316876001600160a01b03167f23ac4832f230e9863286feaebf415429c2049d9f5098e1c1f8743a48773c6d9660405160405180910390a4505050505050565b6130c661441d565b5f6130cf6136bf565b90506130de6005820183614300565b6130fb57604051633ad2b1bb60e11b815260040160405180910390fd5b6040516001600160a01b038316907fb102f7913267c344ac15011acd7185602a74269c32e7783833f5311450fb43dd905f90a25050565b61313a613742565b60045f80516020615f8e833981519152612151600e84614314565b61315d613742565b60065f80516020615f8e833981519152612151600f84614314565b8061318281613dfb565b5f61318c836113ad565b905060078160088111156131a2576131a2614ef7565b036131b5576131b083614335565b6131e7565b60068160088111156131c9576131c9614ef7565b146131e75760405163baf3f0f760e01b815260040160405180910390fd5b5f6131f06136bf565b90506131ff81601001546136e3565b61320960086111eb565b106132275760405163848084dd60e01b815260040160405180910390fd5b600682600881111561323b5761323b614ef7565b036132755761324d816010015461371c565b61325760066111eb565b116132755760405163848084dd60e01b815260040160405180910390fd5b6128a88460085f613a57565b61328961408e565b60045f80516020615fce833981519152611b136016846140c0565b806132ae81613dfb565b6132b9826004613cf7565b6132d65760405163baf3f0f760e01b815260040160405180910390fd5b611aec8260015f613a57565b6132ea613742565b60025f80516020615f8e833981519152612151600a84614314565b61330d613742565b5f6133166136bf565b6001600160a01b0383165f90815260208290526040812060050154919250819003613354576040516357024f6d60e11b815260040160405180910390fd5b60405163e2693e3f60e01b81525f906104019063e2693e3f906133799060040161579e565b602060405180830381865afa158015613394573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906133b891906157d1565b90506001600160a01b0381161561341b5760405163f575c5a760e01b8152600481018390526001600160a01b0382169063f575c5a7906024015f604051808303815f87803b158015613408575f80fd5b505af1925050508015613419575060015b505b6001600160a01b0384165f818152602085815260408083206005019290925590518481527f04078f2e4bb3259952b29491fa528a9a33d9496afca06a26bee6b7b1df26241f9101610f0e565b61346f613742565b60025f80516020615fce833981519152611b136014846140c0565b6060610fcd6134976136bf565b600301613ceb565b6134a7613742565b60035f80516020615f8e833981519152612151600b84614314565b5f6134cb6136bf565b60100154905090565b806134de81613dfb565b6134e9826005613cf7565b6135065760405163baf3f0f760e01b815260040160405180910390fd5b611aec8260046121d58561436e565b8261351f81613776565b6108008211156135425760405163b4fa3fb360e01b815260040160405180910390fd5b828261354c6136bf565b6001600160a01b0387165f9081526020919091526040902060090191613573919083615e14565b50836001600160a01b03167f2013570c343af8ab14a9778150e381a0fda34ed6368127a95fd5e7210cbec5bf8484604051610f0e929190615ec8565b806135b981613dfb565b6135c4826002613cf7565b6135e15760405163baf3f0f760e01b815260040160405180910390fd5b6135ed8260015f613a57565b6040516001600160a01b038316907f8f87baa66b5a3109ebbdf710997ed35a0939f537a70e7ffd6b937a3867e718e2905f90a25050565b61362c613742565b60015f80516020615f8e833981519152612151600d84614314565b5f805f6136526136bf565b905080600c015481600d015492509250509091565b5f6136706136bf565b600b0154905090565b61368161408e565b6001600160a01b0381166136b3575f604051631e4fbdf760e01b81526004016136aa919061518e565b60405180910390fd5b6136bc81614398565b50565b7f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff0090565b5f60048210156136f457505f919050565b6002613701600384615f0a565b61370c906001615a07565b6137169190615f0a565b92915050565b5f600482101561372a575090565b60036137378360026159f0565b61370c906002615a07565b61374a6136bf565b601601546001600160a01b0316331461248b5760405163033b71e160e41b815260040160405180910390fd5b61377e6136bf565b6001600160a01b038281165f9081526020929092526040909120541633146136bc5760405163605919ad60e11b815260040160405180910390fd5b5f613716825490565b5f6137cd83836144ae565b9392505050565b5f805f806137e06136bf565b6001600160a01b0386165f908152602082905260408120919250600a82015460ff16600881111561381357613813614ef7565b03613828575f805f9450945094505050611159565b6001818101546002909201546001600160a01b039283169892169650945092505050565b60608060605f805f61385c6136bf565b90505f61386b826001016137b9565b9050806001600160401b0381111561388557613885615254565b6040519080825280602002602001820160405280156138ae578160200160208202803683370190505b509650806001600160401b038111156138c9576138c9615254565b6040519080825280602002602001820160405280156138f2578160200160208202803683370190505b509550806001600160401b0381111561390d5761390d615254565b604051908082528060200260200182016040528015613936578160200160208202803683370190505b5094505f5b81811015613a08575f61395160018501836137c2565b6001600160a01b0381165f9081526020869052604090208a519192509082908b9085908110613982576139826157ec565b6001600160a01b03928316602091820292909201015260018201548a519116908a90859081106139b4576139b46157ec565b6001600160a01b039283166020918202929092010152600282015489519116908990859081106139e6576139e66157ec565b6001600160a01b0390921660209283029190910190910152505060010161393b565b505060138101546012909101549596949593946001600160a01b039182169490911692509050565b336002600160a01b031461248b576040516354d325c360e01b815260040160405180910390fd5b5f613a606136bf565b6001600160a01b0385165f908152602082905260408120600a8101549293509160ff1690816008811115613a9657613a96614ef7565b1480613ab257505f856008811115613ab057613ab0614ef7565b145b80613ade5750846008811115613aca57613aca614ef7565b816008811115613adc57613adc614ef7565b145b15613aeb57505050505050565b826008015f826008811115613b0257613b02614ef7565b6008811115613b1357613b13614ef7565b81526020019081526020015f205f815480929190613b3090615895565b9190505550826008015f866008811115613b4c57613b4c614ef7565b6008811115613b5d57613b5d614ef7565b81526020019081526020015f205f815480929190613b7a90615786565b9091555060019050816008811115613b9457613b94614ef7565b148015613bb357506001856008811115613bb057613bb0614ef7565b14155b15613bd957613bc56001840187614300565b50613bd36003840187613d2c565b50613c2e565b6001816008811115613bed57613bed614ef7565b14158015613c0c57506001856008811115613c0a57613c0a614ef7565b145b15613c2e57613c1e6001840187613d2c565b50613c2c6003840187614300565b505b600a8201805486919060ff19166001836008811115613c4f57613c4f614ef7565b021790555060048201849055846008811115613c6d57613c6d614ef7565b816008811115613c7f57613c7f614ef7565b6040516001600160a01b038916907fcfb25346bbf2c2f19e20af8b4b4d54cbc6c83057934c1f28539760e8f8065dee905f90a4505050505050565b5f613ce57f000000000000000000000000000000000000000000000000000000000000000043615f1d565b15919050565b60605f6137cd836144d4565b5f816008811115613d0a57613d0a614ef7565b613d13846113ad565b6008811115613d2457613d24614ef7565b149392505050565b5f6137cd836001600160a01b03841661452d565b5f613d49614617565b90506001600160a01b038116613d725760405163cdded31d60e01b815260040160405180910390fd5b604051631f7f8a5f60e21b81526001600160a01b03821690637dfe297c90613d9e90859060040161518e565b602060405180830381865afa158015613db9573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613ddd9190615f30565b15611aec5760405163b4fa3fb360e01b815260040160405180910390fd5b336001600160a01b038216146136bc576040516335f1334d60e11b815260040160405180910390fd5b5f80613e2e6136bf565b6001600160a01b038085165f9081526020928352604080822060010154815163318588a360e11b81529151931694509092849263630b11469260048082019392918290030181865afa158015613e86573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613eaa9190615f4f565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015613ee6573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190613f0a9190615f4f565b613f149190615f66565b6a0422ca8b0a00a4250000001115949350505050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161480613fb057507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316613fa45f80516020615fae833981519152546001600160a01b031690565b6001600160a01b031614155b1561248b5760405163703e46dd60e11b815260040160405180910390fd5b6136bc61408e565b816001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa925050508015614030575060408051601f3d908101601f1916820190925261402d91810190615f4f565b60015b61404f5781604051634c9c8ce360e01b81526004016136aa919061518e565b5f80516020615fae833981519152811461407f57604051632a87526960e21b8152600481018290526024016136aa565b6140898383614699565b505050565b33614097612ead565b6001600160a01b03161461248b573360405163118cdaa760e01b81526004016136aa919061518e565b5f6001600160a01b0382166140e85760405163b4fa3fb360e01b815260040160405180910390fd5b5f614113847f1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff00615a07565b8054939055509092915050565b306001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000161461248b5760405163703e46dd60e11b815260040160405180910390fd5b614175858585856146ee565b5f61417e614617565b90506001600160a01b0381166141a75760405163cdded31d60e01b815260040160405180910390fd5b60405163669d8d4560e01b815233906001600160a01b0383169063669d8d45906141d590899060040161518e565b602060405180830381865afa1580156141f0573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061421491906157d1565b6001600160a01b03161461423b57604051632281776f60e01b815260040160405180910390fd5b61424686868461485a565b5f6142518686614936565b6001600160a01b03161480156142cd5750604051631f7f8a5f60e21b81526001600160a01b03821690637dfe297c9061428e90879060040161518e565b602060405180830381865afa1580156142a9573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906142cd9190615f30565b156142eb5760405163b4fa3fb360e01b815260040160405180910390fd5b6142f7878787876149e1565b50505050505050565b5f6137cd836001600160a01b038416614aab565b5f815f036140e85760405163b4fa3fb360e01b815260040160405180910390fd5b5f61433f8261436e565b905080158015906143505750804210155b15611aec5760405163b48d5fc760e01b815260040160405180910390fd5b5f6143776136bf565b6001600160a01b039092165f90815260209290925250604090206004015490565b5f6143a161448a565b80546001600160a01b038481166001600160a01b031983168117845560405193945091169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b5f610fcd7f000000000000000000000000000000000000000000000000000000000000000043615f0a565b6144256136bf565b601501546001600160a01b0316331461248b5760405163333f4e6560e01b815260040160405180910390fd5b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00613716565b614481614af7565b6136bc81614b1c565b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930090565b5f825f0182815481106144c3576144c36157ec565b905f5260205f200154905092915050565b6060815f0180548060200260200160405190810160405280929190818152602001828054801561452157602002820191905f5260205f20905b81548152602001906001019080831161450d575b50505050509050919050565b5f8181526001830160205260408120548015614607575f61454f600183615f66565b85549091505f9061456290600190615f66565b90508082146145c1575f865f018281548110614580576145806157ec565b905f5260205f200154905080875f0184815481106145a0576145a06157ec565b5f918252602080832090910192909255918252600188019052604090208390555b85548690806145d2576145d2615f79565b600190038181905f5260205f20015f90559055856001015f8681526020019081526020015f205f905560019350505050613716565b5f915050613716565b5092915050565b60405163e2693e3f60e01b815260206004820152601060248201526f436e5374616b696e67466163746f727960801b60448201525f906104019063e2693e3f90606401602060405180830381865afa158015614675573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610fcd91906157d1565b6146a282614b24565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156146e6576140898282614b7e565b611aec614c0e565b6001600160a01b038416158061470b57506001600160a01b038316155b8061471d57506001600160a01b038216155b1561473b5760405163b4fa3fb360e01b815260040160405180910390fd5b826001600160a01b0316846001600160a01b0316148061476c5750816001600160a01b0316846001600160a01b0316145b806147885750816001600160a01b0316836001600160a01b0316145b156147a65760405163b4fa3fb360e01b815260040160405180910390fd5b80515160301415806147be5750806020015151606014155b156147dc5760405163b4fa3fb360e01b815260040160405180910390fd5b805180516020909101207fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd293148061483c575060208082015180519101207f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21145b156128a85760405163b4fa3fb360e01b815260040160405180910390fd5b604080517f23ae25c387ef8bd2c14b622e10202a494f464e31b796f40e83e9aecdf9cb42fb602082015246918101919091523060608201523360808201526001600160a01b0380851660a0830152831660c08201525f9060e0016040516020818303038152906040528051906020012090505f806148d88385614c2d565b5090925090505f8160038111156148f1576148f1614ef7565b1415806149105750856001600160a01b0316826001600160a01b031614155b1561492e57604051631ea9ff4d60e21b815260040160405180910390fd5b505050505050565b5f826001600160a01b031663e1a12d356040518163ffffffff1660e01b8152600401602060405180830381865afa158015614973573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061499791906157d1565b90506001600160a01b038116158015906149c35750816001600160a01b0316816001600160a01b031614155b156137165760405163b4fa3fb360e01b815260040160405180910390fd5b6001600160a01b0383165f9081526020859052604090205460ff1680614a1e57506001600160a01b0382165f9081526020859052604090205460ff165b80614a4057506001600160a01b0381165f9081526020859052604090205460ff165b15614a5e576040516316a163b960e11b815260040160405180910390fd5b6001600160a01b039283165f90815260209490945260408085208054600160ff19918216811790925593851686528186208054851682179055919093168452919092208054909216179055565b5f818152600183016020526040812054614af057508154600181810184555f848152602080822090930184905584548482528286019093526040902091909155613716565b505f613716565b614aff614c76565b61248b57604051631afcd79f60e31b815260040160405180910390fd5b613681614af7565b806001600160a01b03163b5f03614b505780604051634c9c8ce360e01b81526004016136aa919061518e565b5f80516020615fae83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b60605f614b8b8484614c8f565b9050808015614bac57505f3d1180614bac57505f846001600160a01b03163b115b15614bc157614bb9614ca2565b915050613716565b8015614be25783604051639996b31560e01b81526004016136aa919061518e565b3d15614bf557614bf0614cbb565b614610565b60405163d6bda27560e01b815260040160405180910390fd5b341561248b5760405163b398979f60e01b815260040160405180910390fd5b5f805f8351604103614c64576020840151604085015160608601515f1a614c5688828585614cc6565b955095509550505050614c6f565b505081515f91506002905b9250925092565b5f614c7f614451565b54600160401b900460ff16919050565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b5f80806fa2a8918ca85bafe22016d0b997e4df60600160ff1b03841115614cf557505f91506003905082614d7a565b604080515f808252602082018084528a905260ff891692820192909252606081018790526080810186905260019060a0016020604051602081039080840390855afa158015614d46573d5f803e3d5ffd5b5050604051601f1901519150506001600160a01b038116614d7157505f925060019150829050614d7a565b92505f91508190505b9450945094915050565b508054614d90906158aa565b5f825580601f10614d9f575050565b601f0160209004905f5260205f20908101906136bc9190614e45565b6040518061014001604052805f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f81526020015f8152602001614e2c604051806040016040528060608152602001606081525090565b815260606020820181905260408201819052015f905290565b5b80821115614e59575f8155600101614e46565b5090565b5f60208284031215614e6d575f80fd5b5035919050565b6001600160a01b03811681146136bc575f80fd5b8035614e9381614e74565b919050565b5f60208284031215614ea8575f80fd5b81356137cd81614e74565b5f8060408385031215614ec4575f80fd5b8235614ecf81614e74565b91506020830135614edf81614e74565b809150509250929050565b6001600160a01b03169052565b634e487b7160e01b5f52602160045260245ffd5b60098110614f2757634e487b7160e01b5f52602160045260245ffd5b9052565b602080825282518282018190525f919060409081850190868401855b82811015614fa557815180516001600160a01b039081168652878201518116888701528682015116868601526060808201519086015260809081015190614f9081870183614f0b565b505060a0939093019290850190600101614f47565b5091979650505050505050565b5f815180845260208085019450602084015f5b83811015614fea5781516001600160a01b031687529582019590820190600101614fc5565b509495945050505050565b60a081525f61500760a0830188614fb2565b82810360208401526150198188614fb2565b9050828103604084015261502d8187614fb2565b6001600160a01b0395861660608501529390941660809092019190915250949350505050565b604081525f6150656040830185614fb2565b90508260208301529392505050565b600981106136bc575f80fd5b5f60208284031215615090575f80fd5b81356137cd81615074565b5f8083601f8401126150ab575f80fd5b5081356001600160401b038111156150c1575f80fd5b6020830191508360208260051b85010111156150db575f80fd5b9250929050565b5f805f805f805f6080888a0312156150f8575f80fd5b87356001600160401b038082111561510e575f80fd5b61511a8b838c0161509b565b909950975060208a0135915080821115615132575f80fd5b61513e8b838c0161509b565b909750955060408a0135915080821115615156575f80fd5b506151638a828b0161509b565b989b979a50959894979596606090950135949350505050565b602081525f6137cd6020830184614fb2565b6001600160a01b0391909116815260200190565b602081016137168284614f0b565b602080825282518282018190525f919060409081850190868401855b82811015614fa557815180516001600160a01b0390811686528782015181168887015286820151168686015260609081015190850152608090930192908501906001016151cc565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f6137cd6020830184615214565b634e487b7160e01b5f52604160045260245ffd5b604080519081016001600160401b038111828210171561528a5761528a615254565b60405290565b60405161014081016001600160401b038111828210171561528a5761528a615254565b6040516101e081016001600160401b038111828210171561528a5761528a615254565b604051601f8201601f191681016001600160401b03811182821017156152fe576152fe615254565b604052919050565b5f6001600160401b0382111561531e5761531e615254565b50601f01601f191660200190565b5f82601f83011261533b575f80fd5b813561534e61534982615306565b6152d6565b818152846020838601011115615362575f80fd5b816020850160208301375f918101602001919091529392505050565b5f806040838503121561538f575f80fd5b823561539a81614e74565b915060208301356001600160401b038111156153b4575f80fd5b6153c08582860161532c565b9150509250929050565b5f604082840312156153da575f80fd5b6153e2615268565b905081356001600160401b03808211156153fa575f80fd5b6154068583860161532c565b8352602084013591508082111561541b575f80fd5b506154288482850161532c565b60208301525092915050565b5f805f805f805f80610100898b03121561544c575f80fd5b61545589614e88565b975061546360208a01614e88565b965061547160408a01614e88565b955061547f60608a01614e88565b945060808901356001600160401b038082111561549a575f80fd5b6154a68c838d016153ca565b955060a08b01359150808211156154bb575f80fd5b6154c78c838d0161532c565b945060c08b01359150808211156154dc575f80fd5b6154e88c838d0161532c565b935060e08b01359150808211156154fd575f80fd5b5061550a8b828c0161532c565b9150509295985092959890939650565b5f81516040845261552e6040850182615214565b9050602083015184820360208601526155478282615214565b95945050505050565b60208152615562602082018351614eea565b5f60208301516155756040840182614eea565b5060408301516155886060840182614eea565b50606083015161559b6080840182614eea565b50608083015160a083015260a083015160c083015260c08301516101408060e08501526155cc61016085018361551a565b915060e0850151601f196101008187860301818801526155ec8584615214565b94508088015192505061012081878603018188015261560b8584615214565b9450808801519250505061562182860182614f0b565b5090949350505050565b604081525f61563d6040830185614fb2565b6020838203818501528185518084528284019150828160051b8501018388015f5b8381101561568c57601f1987840301855261567a83835161551a565b9486019492509085019060010161565e565b50909998505050505050505050565b604080825283519082018190525f906020906060840190828701845b828110156156d657815160ff16845292840192908401906001016156b7565b50505083810360208501526156eb8186614fb2565b9695505050505050565b5f805f60408486031215615707575f80fd5b833561571281614e74565b925060208401356001600160401b038082111561572d575f80fd5b818601915086601f830112615740575f80fd5b81358181111561574e575f80fd5b87602082850101111561575f575f80fd5b6020830194508093505050509250925092565b634e487b7160e01b5f52601160045260245ffd5b5f6001820161579757615797615772565b5060010190565b6020808252600e908201526d29ba30b5b4b733aa3930b1b5b2b960911b604082015260600190565b8051614e9381614e74565b5f602082840312156157e1575f80fd5b81516137cd81614e74565b634e487b7160e01b5f52603260045260245ffd5b604080825281018490525f8560608301825b8781101561584257823561582581614e74565b6001600160a01b0316825260209283019290910190600101615812565b508381036020858101919091528582529150859082015f5b8681101561588857823561586d81615074565b6158778382614f0b565b50918301919083019060010161585a565b5098975050505050505050565b5f816158a3576158a3615772565b505f190190565b600181811c908216806158be57607f821691505b6020821081036158dc57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561408957805f5260205f20601f840160051c810160208510156159075750805b601f840160051c820191505b81811015612ea6575f8155600101615913565b5f19600383901b1c191660019190911b1790565b81516001600160401b0381111561595357615953615254565b6159678161596184546158aa565b846158e2565b602080601f831160018114615995575f84156159835750858301515b61598d8582615926565b86555061492e565b5f85815260208120601f198616915b828110156159c3578886015182559484019460019091019084016159a4565b50858210156159e057878501515f19600388901b60f8161c191681555b5050505050600190811b01905550565b808202811582820484141761371657613716615772565b8082018082111561371657613716615772565b5f6001600160401b03821115615a3257615a32615254565b5060051b60200190565b5f82601f830112615a4b575f80fd5b81516020615a5b61534983615a1a565b8083825260208201915060208460051b870101935086841115615a7c575f80fd5b602086015b84811015615aa1578051615a9481614e74565b8352918301918301615a81565b509695505050505050565b5f82601f830112615abb575f80fd5b8151615ac961534982615306565b818152846020838601011115615add575f80fd5b8160208501602083015e5f918101602001919091529392505050565b5f60408284031215615b09575f80fd5b615b11615268565b905081516001600160401b0380821115615b29575f80fd5b615b3585838601615aac565b83526020840151915080821115615b4a575f80fd5b5061542884828501615aac565b8051614e9381615074565b5f82601f830112615b71575f80fd5b81516020615b8161534983615a1a565b82815260059290921b84018101918181019086841115615b9f575f80fd5b8286015b84811015615aa15780516001600160401b0380821115615bc1575f80fd5b90880190610140828b03601f1901811315615bda575f80fd5b615be2615290565b615bed8885016157c6565b81526040615bfc8186016157c6565b898301526060615c0d8187016157c6565b8284015260809150615c208287016157c6565b818401525060a0808601518284015260c0915081860151818401525060e08086015185811115615c4e575f80fd5b615c5c8f8c838a0101615af9565b838501525061010091508186015185811115615c76575f80fd5b615c848f8c838a0101615aac565b8285015250506101208086015185811115615c9d575f80fd5b615cab8f8c838a0101615aac565b8385015250615cbb848701615b57565b90830152508652505050918301918301615ba3565b5f60208284031215615ce0575f80fd5b81516001600160401b0380821115615cf6575f80fd5b908301906101e08286031215615d0a575f80fd5b615d126152b3565b615d1b836157c6565b8152615d29602084016157c6565b6020820152615d3a604084016157c6565b6040820152606083015160608201526080830151608082015260a083015160a082015260c083015160c082015260e083015160e0820152610100808401518183015250610120808401518183015250610140615d978185016157c6565b90820152610160615da98482016157c6565b90820152610180615dbb8482016157c6565b908201526101a08381015183811115615dd2575f80fd5b615dde88828701615a3c565b8284015250506101c08084015183811115615df7575f80fd5b615e0388828701615b62565b918301919091525095945050505050565b6001600160401b03831115615e2b57615e2b615254565b615e3f83615e3983546158aa565b836158e2565b5f601f841160018114615e6b575f8515615e595750838201355b615e638682615926565b845550612ea6565b5f83815260208120601f198716915b82811015615e9a5786850135825560209485019460019092019101615e7a565b5086821015615eb6575f1960f88860031b161c19848701351681555b505060018560011b0183555050505050565b60208152816020820152818360408301375f818301604090810191909152601f909201601f19160101919050565b634e487b7160e01b5f52601260045260245ffd5b5f82615f1857615f18615ef6565b500490565b5f82615f2b57615f2b615ef6565b500690565b5f60208284031215615f40575f80fd5b815180151581146137cd575f80fd5b5f60208284031215615f5f575f80fd5b5051919050565b8181038181111561371657613716615772565b634e487b7160e01b5f52603160045260245ffdfe34e70e79c69eb46175bef4bfa16c239443c0aaf0bf701d30389fc9144da5e6bd360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbcd07d74a393c991393c31f5d832e6c292f2557e27ae0daffecbc1dd50f89cbe4ba164736f6c6343000819000a`

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

// AddressBookV2GcIdRevokedIterator is returned from FilterGcIdRevoked and is used to iterate over the raw logs and unpacked data for GcIdRevoked events raised by the AddressBookV2 contract.
type AddressBookV2GcIdRevokedIterator struct {
	Event *AddressBookV2GcIdRevoked // Event containing the contract specifics and raw log

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
func (it *AddressBookV2GcIdRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AddressBookV2GcIdRevoked)
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
		it.Event = new(AddressBookV2GcIdRevoked)
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
func (it *AddressBookV2GcIdRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AddressBookV2GcIdRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AddressBookV2GcIdRevoked represents a GcIdRevoked event raised by the AddressBookV2 contract.
type AddressBookV2GcIdRevoked struct {
	NodeId common.Address
	GcId   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterGcIdRevoked is a free log retrieval operation binding the contract event 0x04078f2e4bb3259952b29491fa528a9a33d9496afca06a26bee6b7b1df26241f.
//
// Solidity: event GcIdRevoked(address indexed nodeId, uint256 gcId)
func (_AddressBookV2 *AddressBookV2Filterer) FilterGcIdRevoked(opts *bind.FilterOpts, nodeId []common.Address) (*AddressBookV2GcIdRevokedIterator, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.FilterLogs(opts, "GcIdRevoked", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return &AddressBookV2GcIdRevokedIterator{contract: _AddressBookV2.contract, event: "GcIdRevoked", logs: logs, sub: sub}, nil
}

// WatchGcIdRevoked is a free log subscription operation binding the contract event 0x04078f2e4bb3259952b29491fa528a9a33d9496afca06a26bee6b7b1df26241f.
//
// Solidity: event GcIdRevoked(address indexed nodeId, uint256 gcId)
func (_AddressBookV2 *AddressBookV2Filterer) WatchGcIdRevoked(opts *bind.WatchOpts, sink chan<- *AddressBookV2GcIdRevoked, nodeId []common.Address) (event.Subscription, error) {

	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}

	logs, sub, err := _AddressBookV2.contract.WatchLogs(opts, "GcIdRevoked", nodeIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AddressBookV2GcIdRevoked)
				if err := _AddressBookV2.contract.UnpackLog(event, "GcIdRevoked", log); err != nil {
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

// ParseGcIdRevoked is a log parse operation binding the contract event 0x04078f2e4bb3259952b29491fa528a9a33d9496afca06a26bee6b7b1df26241f.
//
// Solidity: event GcIdRevoked(address indexed nodeId, uint256 gcId)
func (_AddressBookV2 *AddressBookV2Filterer) ParseGcIdRevoked(log types.Log) (*AddressBookV2GcIdRevoked, error) {
	event := new(AddressBookV2GcIdRevoked)
	if err := _AddressBookV2.contract.UnpackLog(event, "GcIdRevoked", log); err != nil {
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
