// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package testing

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

// Profile is an auto generated low-level Go binding around an user-defined struct.
type Profile struct {
	NodeId          common.Address
	StakingContract common.Address
	RewardAddress   common.Address
	TimeoutAt       *big.Int
	State           uint8
}

// IAddressBookMetaData contains all meta data concerning the IAddressBook contract.
var IAddressBookMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_adminList\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"_requirement\",\"type\":\"uint256\"}],\"name\":\"constructContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAddress\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAddressInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"cnNodeIdList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"cnStakingContractList\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"cnRewardAddressList\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"pocContractAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kirContractAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_cnNodeId\",\"type\":\"address\"}],\"name\":\"getCnInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"cnStakingcontract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"cnRewardAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingRequestList\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"pendingRequestList\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_id\",\"type\":\"bytes32\"}],\"name\":\"getRequestInfo\",\"outputs\":[{\"internalType\":\"enumIAddressBook.Functions\",\"name\":\"functionId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"firstArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"secondArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"thirdArg\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"confirmers\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"initialProposedTime\",\"type\":\"uint256\"},{\"internalType\":\"enumIAddressBook.RequestState\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIAddressBook.Functions\",\"name\":\"_functionId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_firstArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_secondArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_thirdArg\",\"type\":\"bytes32\"}],\"name\":\"getRequestInfoByArgs\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"confirmers\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"initialProposedTime\",\"type\":\"uint256\"},{\"internalType\":\"enumIAddressBook.RequestState\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"adminList\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"requirement\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isActivated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isConstructed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kirContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pocContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rewardAddress\",\"type\":\"address\"}],\"name\":\"reviseRewardAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"enumIAddressBook.Functions\",\"name\":\"_functionId\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_firstArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_secondArg\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_thirdArg\",\"type\":\"bytes32\"}],\"name\":\"revokeRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"spareContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"submitActivateAddressBook\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"submitAddAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"submitClearRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"}],\"name\":\"submitDeleteAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_cnNodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_cnStakingContractAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_cnRewardAddress\",\"type\":\"address\"}],\"name\":\"submitRegisterCnStakingContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_cnNodeId\",\"type\":\"address\"}],\"name\":\"submitUnregisterCnStakingContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_kirContractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_version\",\"type\":\"uint256\"}],\"name\":\"submitUpdateKirContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_pocContractAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_version\",\"type\":\"uint256\"}],\"name\":\"submitUpdatePocContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requirement\",\"type\":\"uint256\"}],\"name\":\"submitUpdateRequirement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_spareContractAddress\",\"type\":\"address\"}],\"name\":\"submitUpdateSpareContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"7894c366": "constructContract(address[],uint256)",
		"715b208b": "getAllAddress()",
		"160370b8": "getAllAddressInfo()",
		"15575d5a": "getCnInfo(address)",
		"da34a0bd": "getPendingRequestList()",
		"82d67e5a": "getRequestInfo(bytes32)",
		"407091eb": "getRequestInfoByArgs(uint8,bytes32,bytes32,bytes32)",
		"1865c57d": "getState()",
		"4a8c1fb4": "isActivated()",
		"50a5bb69": "isConstructed()",
		"b858dd95": "kirContractAddress()",
		"d267eda5": "pocContractAddress()",
		"832a2aad": "reviseRewardAddress(address)",
		"3f0628b1": "revokeRequest(uint8,bytes32,bytes32,bytes32)",
		"6abd623d": "spareContractAddress()",
		"feb15ca1": "submitActivateAddressBook()",
		"863f5c0a": "submitAddAdmin(address)",
		"87cd9feb": "submitClearRequest()",
		"791b5123": "submitDeleteAdmin(address)",
		"cc11efc0": "submitRegisterCnStakingContract(address,address,address)",
		"b5067706": "submitUnregisterCnStakingContract(address)",
		"9258d768": "submitUpdateKirContract(address,uint256)",
		"21ac4ad4": "submitUpdatePocContract(address,uint256)",
		"e748357b": "submitUpdateRequirement(uint256)",
		"394a144a": "submitUpdateSpareContract(address)",
	},
}

// IAddressBookABI is the input ABI used to generate the binding from.
// Deprecated: Use IAddressBookMetaData.ABI instead.
var IAddressBookABI = IAddressBookMetaData.ABI

// IAddressBookBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IAddressBookBinRuntime = ``

// Deprecated: Use IAddressBookMetaData.Sigs instead.
// IAddressBookFuncSigs maps the 4-byte function signature to its string representation.
var IAddressBookFuncSigs = IAddressBookMetaData.Sigs

// IAddressBook is an auto generated Go binding around a Kaia contract.
type IAddressBook struct {
	IAddressBookCaller     // Read-only binding to the contract
	IAddressBookTransactor // Write-only binding to the contract
	IAddressBookFilterer   // Log filterer for contract events
}

// IAddressBookCaller is an auto generated read-only Go binding around a Kaia contract.
type IAddressBookCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookTransactor is an auto generated write-only Go binding around a Kaia contract.
type IAddressBookTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IAddressBookFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IAddressBookSession struct {
	Contract     *IAddressBook     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAddressBookCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IAddressBookCallerSession struct {
	Contract *IAddressBookCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// IAddressBookTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IAddressBookTransactorSession struct {
	Contract     *IAddressBookTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// IAddressBookRaw is an auto generated low-level Go binding around a Kaia contract.
type IAddressBookRaw struct {
	Contract *IAddressBook // Generic contract binding to access the raw methods on
}

// IAddressBookCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IAddressBookCallerRaw struct {
	Contract *IAddressBookCaller // Generic read-only contract binding to access the raw methods on
}

// IAddressBookTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IAddressBookTransactorRaw struct {
	Contract *IAddressBookTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAddressBook creates a new instance of IAddressBook, bound to a specific deployed contract.
func NewIAddressBook(address common.Address, backend bind.ContractBackend) (*IAddressBook, error) {
	contract, err := bindIAddressBook(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAddressBook{IAddressBookCaller: IAddressBookCaller{contract: contract}, IAddressBookTransactor: IAddressBookTransactor{contract: contract}, IAddressBookFilterer: IAddressBookFilterer{contract: contract}}, nil
}

// NewIAddressBookCaller creates a new read-only instance of IAddressBook, bound to a specific deployed contract.
func NewIAddressBookCaller(address common.Address, caller bind.ContractCaller) (*IAddressBookCaller, error) {
	contract, err := bindIAddressBook(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAddressBookCaller{contract: contract}, nil
}

// NewIAddressBookTransactor creates a new write-only instance of IAddressBook, bound to a specific deployed contract.
func NewIAddressBookTransactor(address common.Address, transactor bind.ContractTransactor) (*IAddressBookTransactor, error) {
	contract, err := bindIAddressBook(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAddressBookTransactor{contract: contract}, nil
}

// NewIAddressBookFilterer creates a new log filterer instance of IAddressBook, bound to a specific deployed contract.
func NewIAddressBookFilterer(address common.Address, filterer bind.ContractFilterer) (*IAddressBookFilterer, error) {
	contract, err := bindIAddressBook(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAddressBookFilterer{contract: contract}, nil
}

// bindIAddressBook binds a generic wrapper to an already deployed contract.
func bindIAddressBook(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAddressBookMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAddressBook *IAddressBookRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAddressBook.Contract.IAddressBookCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAddressBook *IAddressBookRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBook.Contract.IAddressBookTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAddressBook *IAddressBookRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAddressBook.Contract.IAddressBookTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAddressBook *IAddressBookCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAddressBook.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAddressBook *IAddressBookTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBook.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAddressBook *IAddressBookTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAddressBook.Contract.contract.Transact(opts, method, params...)
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_IAddressBook *IAddressBookCaller) GetAllAddress(opts *bind.CallOpts) (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getAllAddress")

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
func (_IAddressBook *IAddressBookSession) GetAllAddress() (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error) {
	return _IAddressBook.Contract.GetAllAddress(&_IAddressBook.CallOpts)
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_IAddressBook *IAddressBookCallerSession) GetAllAddress() (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error) {
	return _IAddressBook.Contract.GetAllAddress(&_IAddressBook.CallOpts)
}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[] cnNodeIdList, address[] cnStakingContractList, address[] cnRewardAddressList, address pocContractAddress, address kirContractAddress)
func (_IAddressBook *IAddressBookCaller) GetAllAddressInfo(opts *bind.CallOpts) (struct {
	CnNodeIdList          []common.Address
	CnStakingContractList []common.Address
	CnRewardAddressList   []common.Address
	PocContractAddress    common.Address
	KirContractAddress    common.Address
}, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getAllAddressInfo")

	outstruct := new(struct {
		CnNodeIdList          []common.Address
		CnStakingContractList []common.Address
		CnRewardAddressList   []common.Address
		PocContractAddress    common.Address
		KirContractAddress    common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CnNodeIdList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.CnStakingContractList = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.CnRewardAddressList = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.PocContractAddress = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.KirContractAddress = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[] cnNodeIdList, address[] cnStakingContractList, address[] cnRewardAddressList, address pocContractAddress, address kirContractAddress)
func (_IAddressBook *IAddressBookSession) GetAllAddressInfo() (struct {
	CnNodeIdList          []common.Address
	CnStakingContractList []common.Address
	CnRewardAddressList   []common.Address
	PocContractAddress    common.Address
	KirContractAddress    common.Address
}, error) {
	return _IAddressBook.Contract.GetAllAddressInfo(&_IAddressBook.CallOpts)
}

// GetAllAddressInfo is a free data retrieval call binding the contract method 0x160370b8.
//
// Solidity: function getAllAddressInfo() view returns(address[] cnNodeIdList, address[] cnStakingContractList, address[] cnRewardAddressList, address pocContractAddress, address kirContractAddress)
func (_IAddressBook *IAddressBookCallerSession) GetAllAddressInfo() (struct {
	CnNodeIdList          []common.Address
	CnStakingContractList []common.Address
	CnRewardAddressList   []common.Address
	PocContractAddress    common.Address
	KirContractAddress    common.Address
}, error) {
	return _IAddressBook.Contract.GetAllAddressInfo(&_IAddressBook.CallOpts)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookCaller) GetCnInfo(opts *bind.CallOpts, _cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getCnInfo", _cnNodeId)

	outstruct := new(struct {
		CnNodeId          common.Address
		CnStakingcontract common.Address
		CnRewardAddress   common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CnNodeId = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.CnStakingcontract = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.CnRewardAddress = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookSession) GetCnInfo(_cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error) {
	return _IAddressBook.Contract.GetCnInfo(&_IAddressBook.CallOpts, _cnNodeId)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookCallerSession) GetCnInfo(_cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error) {
	return _IAddressBook.Contract.GetCnInfo(&_IAddressBook.CallOpts, _cnNodeId)
}

// GetPendingRequestList is a free data retrieval call binding the contract method 0xda34a0bd.
//
// Solidity: function getPendingRequestList() view returns(bytes32[] pendingRequestList)
func (_IAddressBook *IAddressBookCaller) GetPendingRequestList(opts *bind.CallOpts) ([][32]byte, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getPendingRequestList")

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// GetPendingRequestList is a free data retrieval call binding the contract method 0xda34a0bd.
//
// Solidity: function getPendingRequestList() view returns(bytes32[] pendingRequestList)
func (_IAddressBook *IAddressBookSession) GetPendingRequestList() ([][32]byte, error) {
	return _IAddressBook.Contract.GetPendingRequestList(&_IAddressBook.CallOpts)
}

// GetPendingRequestList is a free data retrieval call binding the contract method 0xda34a0bd.
//
// Solidity: function getPendingRequestList() view returns(bytes32[] pendingRequestList)
func (_IAddressBook *IAddressBookCallerSession) GetPendingRequestList() ([][32]byte, error) {
	return _IAddressBook.Contract.GetPendingRequestList(&_IAddressBook.CallOpts)
}

// GetRequestInfo is a free data retrieval call binding the contract method 0x82d67e5a.
//
// Solidity: function getRequestInfo(bytes32 _id) view returns(uint8 functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookCaller) GetRequestInfo(opts *bind.CallOpts, _id [32]byte) (struct {
	FunctionId          uint8
	FirstArg            [32]byte
	SecondArg           [32]byte
	ThirdArg            [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getRequestInfo", _id)

	outstruct := new(struct {
		FunctionId          uint8
		FirstArg            [32]byte
		SecondArg           [32]byte
		ThirdArg            [32]byte
		Confirmers          []common.Address
		InitialProposedTime *big.Int
		State               uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FunctionId = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.FirstArg = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.SecondArg = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.ThirdArg = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)
	outstruct.Confirmers = *abi.ConvertType(out[4], new([]common.Address)).(*[]common.Address)
	outstruct.InitialProposedTime = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.State = *abi.ConvertType(out[6], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetRequestInfo is a free data retrieval call binding the contract method 0x82d67e5a.
//
// Solidity: function getRequestInfo(bytes32 _id) view returns(uint8 functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookSession) GetRequestInfo(_id [32]byte) (struct {
	FunctionId          uint8
	FirstArg            [32]byte
	SecondArg           [32]byte
	ThirdArg            [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error) {
	return _IAddressBook.Contract.GetRequestInfo(&_IAddressBook.CallOpts, _id)
}

// GetRequestInfo is a free data retrieval call binding the contract method 0x82d67e5a.
//
// Solidity: function getRequestInfo(bytes32 _id) view returns(uint8 functionId, bytes32 firstArg, bytes32 secondArg, bytes32 thirdArg, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookCallerSession) GetRequestInfo(_id [32]byte) (struct {
	FunctionId          uint8
	FirstArg            [32]byte
	SecondArg           [32]byte
	ThirdArg            [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error) {
	return _IAddressBook.Contract.GetRequestInfo(&_IAddressBook.CallOpts, _id)
}

// GetRequestInfoByArgs is a free data retrieval call binding the contract method 0x407091eb.
//
// Solidity: function getRequestInfoByArgs(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) view returns(bytes32 id, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookCaller) GetRequestInfoByArgs(opts *bind.CallOpts, _functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (struct {
	Id                  [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getRequestInfoByArgs", _functionId, _firstArg, _secondArg, _thirdArg)

	outstruct := new(struct {
		Id                  [32]byte
		Confirmers          []common.Address
		InitialProposedTime *big.Int
		State               uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Confirmers = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.InitialProposedTime = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.State = *abi.ConvertType(out[3], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetRequestInfoByArgs is a free data retrieval call binding the contract method 0x407091eb.
//
// Solidity: function getRequestInfoByArgs(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) view returns(bytes32 id, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookSession) GetRequestInfoByArgs(_functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (struct {
	Id                  [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error) {
	return _IAddressBook.Contract.GetRequestInfoByArgs(&_IAddressBook.CallOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// GetRequestInfoByArgs is a free data retrieval call binding the contract method 0x407091eb.
//
// Solidity: function getRequestInfoByArgs(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) view returns(bytes32 id, address[] confirmers, uint256 initialProposedTime, uint8 state)
func (_IAddressBook *IAddressBookCallerSession) GetRequestInfoByArgs(_functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (struct {
	Id                  [32]byte
	Confirmers          []common.Address
	InitialProposedTime *big.Int
	State               uint8
}, error) {
	return _IAddressBook.Contract.GetRequestInfoByArgs(&_IAddressBook.CallOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 requirement)
func (_IAddressBook *IAddressBookCaller) GetState(opts *bind.CallOpts) (struct {
	AdminList   []common.Address
	Requirement *big.Int
}, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		AdminList   []common.Address
		Requirement *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AdminList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Requirement = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 requirement)
func (_IAddressBook *IAddressBookSession) GetState() (struct {
	AdminList   []common.Address
	Requirement *big.Int
}, error) {
	return _IAddressBook.Contract.GetState(&_IAddressBook.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 requirement)
func (_IAddressBook *IAddressBookCallerSession) GetState() (struct {
	AdminList   []common.Address
	Requirement *big.Int
}, error) {
	return _IAddressBook.Contract.GetState(&_IAddressBook.CallOpts)
}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_IAddressBook *IAddressBookCaller) IsActivated(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "isActivated")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_IAddressBook *IAddressBookSession) IsActivated() (bool, error) {
	return _IAddressBook.Contract.IsActivated(&_IAddressBook.CallOpts)
}

// IsActivated is a free data retrieval call binding the contract method 0x4a8c1fb4.
//
// Solidity: function isActivated() view returns(bool)
func (_IAddressBook *IAddressBookCallerSession) IsActivated() (bool, error) {
	return _IAddressBook.Contract.IsActivated(&_IAddressBook.CallOpts)
}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_IAddressBook *IAddressBookCaller) IsConstructed(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "isConstructed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_IAddressBook *IAddressBookSession) IsConstructed() (bool, error) {
	return _IAddressBook.Contract.IsConstructed(&_IAddressBook.CallOpts)
}

// IsConstructed is a free data retrieval call binding the contract method 0x50a5bb69.
//
// Solidity: function isConstructed() view returns(bool)
func (_IAddressBook *IAddressBookCallerSession) IsConstructed() (bool, error) {
	return _IAddressBook.Contract.IsConstructed(&_IAddressBook.CallOpts)
}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCaller) KirContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "kirContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_IAddressBook *IAddressBookSession) KirContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.KirContractAddress(&_IAddressBook.CallOpts)
}

// KirContractAddress is a free data retrieval call binding the contract method 0xb858dd95.
//
// Solidity: function kirContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCallerSession) KirContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.KirContractAddress(&_IAddressBook.CallOpts)
}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCaller) PocContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "pocContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_IAddressBook *IAddressBookSession) PocContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.PocContractAddress(&_IAddressBook.CallOpts)
}

// PocContractAddress is a free data retrieval call binding the contract method 0xd267eda5.
//
// Solidity: function pocContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCallerSession) PocContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.PocContractAddress(&_IAddressBook.CallOpts)
}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCaller) SpareContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "spareContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_IAddressBook *IAddressBookSession) SpareContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.SpareContractAddress(&_IAddressBook.CallOpts)
}

// SpareContractAddress is a free data retrieval call binding the contract method 0x6abd623d.
//
// Solidity: function spareContractAddress() view returns(address)
func (_IAddressBook *IAddressBookCallerSession) SpareContractAddress() (common.Address, error) {
	return _IAddressBook.Contract.SpareContractAddress(&_IAddressBook.CallOpts)
}

// ConstructContract is a paid mutator transaction binding the contract method 0x7894c366.
//
// Solidity: function constructContract(address[] _adminList, uint256 _requirement) returns()
func (_IAddressBook *IAddressBookTransactor) ConstructContract(opts *bind.TransactOpts, _adminList []common.Address, _requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "constructContract", _adminList, _requirement)
}

// ConstructContract is a paid mutator transaction binding the contract method 0x7894c366.
//
// Solidity: function constructContract(address[] _adminList, uint256 _requirement) returns()
func (_IAddressBook *IAddressBookSession) ConstructContract(_adminList []common.Address, _requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.ConstructContract(&_IAddressBook.TransactOpts, _adminList, _requirement)
}

// ConstructContract is a paid mutator transaction binding the contract method 0x7894c366.
//
// Solidity: function constructContract(address[] _adminList, uint256 _requirement) returns()
func (_IAddressBook *IAddressBookTransactorSession) ConstructContract(_adminList []common.Address, _requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.ConstructContract(&_IAddressBook.TransactOpts, _adminList, _requirement)
}

// ReviseRewardAddress is a paid mutator transaction binding the contract method 0x832a2aad.
//
// Solidity: function reviseRewardAddress(address _rewardAddress) returns()
func (_IAddressBook *IAddressBookTransactor) ReviseRewardAddress(opts *bind.TransactOpts, _rewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "reviseRewardAddress", _rewardAddress)
}

// ReviseRewardAddress is a paid mutator transaction binding the contract method 0x832a2aad.
//
// Solidity: function reviseRewardAddress(address _rewardAddress) returns()
func (_IAddressBook *IAddressBookSession) ReviseRewardAddress(_rewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.ReviseRewardAddress(&_IAddressBook.TransactOpts, _rewardAddress)
}

// ReviseRewardAddress is a paid mutator transaction binding the contract method 0x832a2aad.
//
// Solidity: function reviseRewardAddress(address _rewardAddress) returns()
func (_IAddressBook *IAddressBookTransactorSession) ReviseRewardAddress(_rewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.ReviseRewardAddress(&_IAddressBook.TransactOpts, _rewardAddress)
}

// RevokeRequest is a paid mutator transaction binding the contract method 0x3f0628b1.
//
// Solidity: function revokeRequest(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) returns()
func (_IAddressBook *IAddressBookTransactor) RevokeRequest(opts *bind.TransactOpts, _functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "revokeRequest", _functionId, _firstArg, _secondArg, _thirdArg)
}

// RevokeRequest is a paid mutator transaction binding the contract method 0x3f0628b1.
//
// Solidity: function revokeRequest(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) returns()
func (_IAddressBook *IAddressBookSession) RevokeRequest(_functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (*types.Transaction, error) {
	return _IAddressBook.Contract.RevokeRequest(&_IAddressBook.TransactOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// RevokeRequest is a paid mutator transaction binding the contract method 0x3f0628b1.
//
// Solidity: function revokeRequest(uint8 _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) returns()
func (_IAddressBook *IAddressBookTransactorSession) RevokeRequest(_functionId uint8, _firstArg [32]byte, _secondArg [32]byte, _thirdArg [32]byte) (*types.Transaction, error) {
	return _IAddressBook.Contract.RevokeRequest(&_IAddressBook.TransactOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// SubmitActivateAddressBook is a paid mutator transaction binding the contract method 0xfeb15ca1.
//
// Solidity: function submitActivateAddressBook() returns()
func (_IAddressBook *IAddressBookTransactor) SubmitActivateAddressBook(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitActivateAddressBook")
}

// SubmitActivateAddressBook is a paid mutator transaction binding the contract method 0xfeb15ca1.
//
// Solidity: function submitActivateAddressBook() returns()
func (_IAddressBook *IAddressBookSession) SubmitActivateAddressBook() (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitActivateAddressBook(&_IAddressBook.TransactOpts)
}

// SubmitActivateAddressBook is a paid mutator transaction binding the contract method 0xfeb15ca1.
//
// Solidity: function submitActivateAddressBook() returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitActivateAddressBook() (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitActivateAddressBook(&_IAddressBook.TransactOpts)
}

// SubmitAddAdmin is a paid mutator transaction binding the contract method 0x863f5c0a.
//
// Solidity: function submitAddAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitAddAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitAddAdmin", _admin)
}

// SubmitAddAdmin is a paid mutator transaction binding the contract method 0x863f5c0a.
//
// Solidity: function submitAddAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookSession) SubmitAddAdmin(_admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitAddAdmin(&_IAddressBook.TransactOpts, _admin)
}

// SubmitAddAdmin is a paid mutator transaction binding the contract method 0x863f5c0a.
//
// Solidity: function submitAddAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitAddAdmin(_admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitAddAdmin(&_IAddressBook.TransactOpts, _admin)
}

// SubmitClearRequest is a paid mutator transaction binding the contract method 0x87cd9feb.
//
// Solidity: function submitClearRequest() returns()
func (_IAddressBook *IAddressBookTransactor) SubmitClearRequest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitClearRequest")
}

// SubmitClearRequest is a paid mutator transaction binding the contract method 0x87cd9feb.
//
// Solidity: function submitClearRequest() returns()
func (_IAddressBook *IAddressBookSession) SubmitClearRequest() (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitClearRequest(&_IAddressBook.TransactOpts)
}

// SubmitClearRequest is a paid mutator transaction binding the contract method 0x87cd9feb.
//
// Solidity: function submitClearRequest() returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitClearRequest() (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitClearRequest(&_IAddressBook.TransactOpts)
}

// SubmitDeleteAdmin is a paid mutator transaction binding the contract method 0x791b5123.
//
// Solidity: function submitDeleteAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitDeleteAdmin(opts *bind.TransactOpts, _admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitDeleteAdmin", _admin)
}

// SubmitDeleteAdmin is a paid mutator transaction binding the contract method 0x791b5123.
//
// Solidity: function submitDeleteAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookSession) SubmitDeleteAdmin(_admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitDeleteAdmin(&_IAddressBook.TransactOpts, _admin)
}

// SubmitDeleteAdmin is a paid mutator transaction binding the contract method 0x791b5123.
//
// Solidity: function submitDeleteAdmin(address _admin) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitDeleteAdmin(_admin common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitDeleteAdmin(&_IAddressBook.TransactOpts, _admin)
}

// SubmitRegisterCnStakingContract is a paid mutator transaction binding the contract method 0xcc11efc0.
//
// Solidity: function submitRegisterCnStakingContract(address _cnNodeId, address _cnStakingContractAddress, address _cnRewardAddress) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitRegisterCnStakingContract(opts *bind.TransactOpts, _cnNodeId common.Address, _cnStakingContractAddress common.Address, _cnRewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitRegisterCnStakingContract", _cnNodeId, _cnStakingContractAddress, _cnRewardAddress)
}

// SubmitRegisterCnStakingContract is a paid mutator transaction binding the contract method 0xcc11efc0.
//
// Solidity: function submitRegisterCnStakingContract(address _cnNodeId, address _cnStakingContractAddress, address _cnRewardAddress) returns()
func (_IAddressBook *IAddressBookSession) SubmitRegisterCnStakingContract(_cnNodeId common.Address, _cnStakingContractAddress common.Address, _cnRewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitRegisterCnStakingContract(&_IAddressBook.TransactOpts, _cnNodeId, _cnStakingContractAddress, _cnRewardAddress)
}

// SubmitRegisterCnStakingContract is a paid mutator transaction binding the contract method 0xcc11efc0.
//
// Solidity: function submitRegisterCnStakingContract(address _cnNodeId, address _cnStakingContractAddress, address _cnRewardAddress) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitRegisterCnStakingContract(_cnNodeId common.Address, _cnStakingContractAddress common.Address, _cnRewardAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitRegisterCnStakingContract(&_IAddressBook.TransactOpts, _cnNodeId, _cnStakingContractAddress, _cnRewardAddress)
}

// SubmitUnregisterCnStakingContract is a paid mutator transaction binding the contract method 0xb5067706.
//
// Solidity: function submitUnregisterCnStakingContract(address _cnNodeId) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUnregisterCnStakingContract(opts *bind.TransactOpts, _cnNodeId common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUnregisterCnStakingContract", _cnNodeId)
}

// SubmitUnregisterCnStakingContract is a paid mutator transaction binding the contract method 0xb5067706.
//
// Solidity: function submitUnregisterCnStakingContract(address _cnNodeId) returns()
func (_IAddressBook *IAddressBookSession) SubmitUnregisterCnStakingContract(_cnNodeId common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUnregisterCnStakingContract(&_IAddressBook.TransactOpts, _cnNodeId)
}

// SubmitUnregisterCnStakingContract is a paid mutator transaction binding the contract method 0xb5067706.
//
// Solidity: function submitUnregisterCnStakingContract(address _cnNodeId) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUnregisterCnStakingContract(_cnNodeId common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUnregisterCnStakingContract(&_IAddressBook.TransactOpts, _cnNodeId)
}

// SubmitUpdateKirContract is a paid mutator transaction binding the contract method 0x9258d768.
//
// Solidity: function submitUpdateKirContract(address _kirContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUpdateKirContract(opts *bind.TransactOpts, _kirContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUpdateKirContract", _kirContractAddress, _version)
}

// SubmitUpdateKirContract is a paid mutator transaction binding the contract method 0x9258d768.
//
// Solidity: function submitUpdateKirContract(address _kirContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookSession) SubmitUpdateKirContract(_kirContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateKirContract(&_IAddressBook.TransactOpts, _kirContractAddress, _version)
}

// SubmitUpdateKirContract is a paid mutator transaction binding the contract method 0x9258d768.
//
// Solidity: function submitUpdateKirContract(address _kirContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUpdateKirContract(_kirContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateKirContract(&_IAddressBook.TransactOpts, _kirContractAddress, _version)
}

// SubmitUpdatePocContract is a paid mutator transaction binding the contract method 0x21ac4ad4.
//
// Solidity: function submitUpdatePocContract(address _pocContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUpdatePocContract(opts *bind.TransactOpts, _pocContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUpdatePocContract", _pocContractAddress, _version)
}

// SubmitUpdatePocContract is a paid mutator transaction binding the contract method 0x21ac4ad4.
//
// Solidity: function submitUpdatePocContract(address _pocContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookSession) SubmitUpdatePocContract(_pocContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdatePocContract(&_IAddressBook.TransactOpts, _pocContractAddress, _version)
}

// SubmitUpdatePocContract is a paid mutator transaction binding the contract method 0x21ac4ad4.
//
// Solidity: function submitUpdatePocContract(address _pocContractAddress, uint256 _version) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUpdatePocContract(_pocContractAddress common.Address, _version *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdatePocContract(&_IAddressBook.TransactOpts, _pocContractAddress, _version)
}

// SubmitUpdateRequirement is a paid mutator transaction binding the contract method 0xe748357b.
//
// Solidity: function submitUpdateRequirement(uint256 _requirement) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUpdateRequirement(opts *bind.TransactOpts, _requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUpdateRequirement", _requirement)
}

// SubmitUpdateRequirement is a paid mutator transaction binding the contract method 0xe748357b.
//
// Solidity: function submitUpdateRequirement(uint256 _requirement) returns()
func (_IAddressBook *IAddressBookSession) SubmitUpdateRequirement(_requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateRequirement(&_IAddressBook.TransactOpts, _requirement)
}

// SubmitUpdateRequirement is a paid mutator transaction binding the contract method 0xe748357b.
//
// Solidity: function submitUpdateRequirement(uint256 _requirement) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUpdateRequirement(_requirement *big.Int) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateRequirement(&_IAddressBook.TransactOpts, _requirement)
}

// SubmitUpdateSpareContract is a paid mutator transaction binding the contract method 0x394a144a.
//
// Solidity: function submitUpdateSpareContract(address _spareContractAddress) returns()
func (_IAddressBook *IAddressBookTransactor) SubmitUpdateSpareContract(opts *bind.TransactOpts, _spareContractAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.contract.Transact(opts, "submitUpdateSpareContract", _spareContractAddress)
}

// SubmitUpdateSpareContract is a paid mutator transaction binding the contract method 0x394a144a.
//
// Solidity: function submitUpdateSpareContract(address _spareContractAddress) returns()
func (_IAddressBook *IAddressBookSession) SubmitUpdateSpareContract(_spareContractAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateSpareContract(&_IAddressBook.TransactOpts, _spareContractAddress)
}

// SubmitUpdateSpareContract is a paid mutator transaction binding the contract method 0x394a144a.
//
// Solidity: function submitUpdateSpareContract(address _spareContractAddress) returns()
func (_IAddressBook *IAddressBookTransactorSession) SubmitUpdateSpareContract(_spareContractAddress common.Address) (*types.Transaction, error) {
	return _IAddressBook.Contract.SubmitUpdateSpareContract(&_IAddressBook.TransactOpts, _spareContractAddress)
}

// MultiCallContractMockPermissionlessMetaData contains all meta data concerning the MultiCallContractMockPermissionless contract.
var MultiCallContractMockPermissionlessMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"multiCallStakingInfo\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallStakingInfoPermissionless\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structProfile[]\",\"name\":\"profiles\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"kefAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kifAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kpfAddr\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"adde19c6": "multiCallStakingInfo()",
		"b04ba218": "multiCallStakingInfoPermissionless()",
	},
	Bin: "0x608060405234801561001057600080fd5b50610c60806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063adde19c61461003b578063b04ba2181461005b575b600080fd5b610043610074565b60405161005293929190610845565b60405180910390f35b610063610212565b6040516100529594939291906108df565b606080806104003b156101ff5760006104009050806001600160a01b031663715b208b6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156100c6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526100ee9190810190610ab1565b805191955093506005111561011157610105610597565b93509350935050909192565b6000600284516101219190610b93565b905061012e600382610bc2565b1561013a575050909192565b610145600382610bd6565b67ffffffffffffffff81111561015d5761015d6109c5565b604051908082528060200260200182016040528015610186578160200160208202803683370190505b50925060005b818110156101f7576101a9816a0422ca8b0a00a425000000610bea565b6101be906a0422ca8b0a00a425000000610c01565b846101ca600384610bd6565b815181106101da576101da610c14565b60209081029190910101526101f0600382610c01565b905061018c565b505050909192565b610207610597565b925092509250909192565b60408051600680825260e08201909252606091829160009182918291816020015b6102626040805160a0810182526000808252602082018190529181018290526060810182905290608082015290565b81526020019060019003908161023357505060408051600680825260e082019092529196506020820160c080368337019050506040805160a081018252610f008152610f016020820152610f029181019190915260006060820152909450608081016006815250856000815181106102dc576102dc610c14565b6020908102919091018101919091526040805160a081018252610f038152610f0492810192909252610f0590820152600060608201526080810160078152508560018151811061032e5761032e610c14565b6020908102919091018101919091526040805160a08101825261100081526110019281019290925261100290820152600060608201526080810160058152508560028151811061038057610380610c14565b6020908102919091018101919091526040805160a0810182526120008152612001928101929092526120029082015260006060820152608081016008815250856003815181106103d2576103d2610c14565b6020908102919091018101919091526040805160a08101825261300081526130019281019290925261300290820152600060608201526080810160028152508560048151811061042457610424610c14565b6020908102919091018101919091526040805160a08101825261400081526140019281019290925261400290820152600060608201526080810160068152508560058151811061047657610476610c14565b60200260200101819052506a0422ca8b0a00a425000000846000815181106104a0576104a0610c14565b6020026020010181815250506a084595161401484a000000846001815181106104cb576104cb610c14565b6020026020010181815250506a069e10de76676d08000000846002815181106104f6576104f6610c14565b6020026020010181815250506a05ca4ec2a79a7f670000008460038151811061052157610521610c14565b6020026020010181815250506a04f68ca6d8cd91c60000008460048151811061054c5761054c610c14565b6020026020010181815250506a0771d2fa45345aa90000008460058151811061057757610577610c14565b6020908102919091010152509293919250610a0191610a029150610a0390565b60408051600580825260c08201909252606091829182916020820160a080368337505060408051600580825260c0820190925292955090506020820160a080368337505060408051600180825281830190925292945090506020808301908036833701905050905060008360008151811061061457610614610c14565b602002602001019060ff16908160ff168152505060018360018151811061063d5761063d610c14565b602002602001019060ff16908160ff168152505060028360028151811061066657610666610c14565b602002602001019060ff16908160ff168152505060038360038151811061068f5761068f610c14565b602002602001019060ff16908160ff16815250506004836004815181106106b8576106b8610c14565b602002602001019060ff16908160ff1681525050610f00826000815181106106e2576106e2610c14565b60200260200101906001600160a01b031690816001600160a01b031681525050610f018260018151811061071857610718610c14565b60200260200101906001600160a01b031690816001600160a01b031681525050610f028260028151811061074e5761074e610c14565b60200260200101906001600160a01b031690816001600160a01b031681525050610f038260038151811061078457610784610c14565b60200260200101906001600160a01b031690816001600160a01b031681525050610f04826004815181106107ba576107ba610c14565b60200260200101906001600160a01b031690816001600160a01b0316815250506a05ca4ec2a79a7f67000000816000815181106107f9576107f9610c14565b602002602001018181525050909192565b600081518084526020808501945080840160005b8381101561083a5781518752958201959082019060010161081e565b509495945050505050565b606080825284519082018190526000906020906080840190828801845b8281101561088157815160ff1684529284019290840190600101610862565b5050508381038285015285518082528683019183019060005b818110156108bf5783516001600160a01b03168352928401929184019160010161089a565b505084810360408601526108d3818761080a565b98975050505050505050565b60a080825286518282018190526000919060209060c0850190828b0185805b8381101561096e57825180516001600160a01b039081168752878201518116888801526040808301519091169087015260608082015190870152608090810151906009821061095b57634e487b7160e01b84526021600452602484fd5b86015293860193918501916001016108fe565b5050505084810382860152610983818a61080a565b935050505061099d60408301866001600160a01b03169052565b6001600160a01b03841660608301526001600160a01b03831660808301529695505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610a0457610a046109c5565b604052919050565b600067ffffffffffffffff821115610a2657610a266109c5565b5060051b60200190565b600082601f830112610a4157600080fd5b81516020610a56610a5183610a0c565b6109db565b82815260059290921b84018101918181019086841115610a7557600080fd5b8286015b84811015610aa65780516001600160a01b0381168114610a995760008081fd5b8352918301918301610a79565b509695505050505050565b60008060408385031215610ac457600080fd5b825167ffffffffffffffff80821115610adc57600080fd5b818501915085601f830112610af057600080fd5b81516020610b00610a5183610a0c565b82815260059290921b84018101918181019089841115610b1f57600080fd5b948201945b83861015610b4d57855160ff81168114610b3e5760008081fd5b82529482019490820190610b24565b91880151919650909350505080821115610b6657600080fd5b50610b7385828601610a30565b9150509250929050565b634e487b7160e01b600052601160045260246000fd5b81810381811115610ba657610ba6610b7d565b92915050565b634e487b7160e01b600052601260045260246000fd5b600082610bd157610bd1610bac565b500690565b600082610be557610be5610bac565b500490565b8082028115828204841417610ba657610ba6610b7d565b80820180821115610ba657610ba6610b7d565b634e487b7160e01b600052603260045260246000fdfea26469706673582212204e29eb3d6f915a4170232fad30819e79eb67f3d1ce17d415a192d6672c5a7d6264736f6c63430008130033",
}

// MultiCallContractMockPermissionlessABI is the input ABI used to generate the binding from.
// Deprecated: Use MultiCallContractMockPermissionlessMetaData.ABI instead.
var MultiCallContractMockPermissionlessABI = MultiCallContractMockPermissionlessMetaData.ABI

// MultiCallContractMockPermissionlessBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MultiCallContractMockPermissionlessBinRuntime = `608060405234801561001057600080fd5b50600436106100365760003560e01c8063adde19c61461003b578063b04ba2181461005b575b600080fd5b610043610074565b60405161005293929190610845565b60405180910390f35b610063610212565b6040516100529594939291906108df565b606080806104003b156101ff5760006104009050806001600160a01b031663715b208b6040518163ffffffff1660e01b8152600401600060405180830381865afa1580156100c6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526100ee9190810190610ab1565b805191955093506005111561011157610105610597565b93509350935050909192565b6000600284516101219190610b93565b905061012e600382610bc2565b1561013a575050909192565b610145600382610bd6565b67ffffffffffffffff81111561015d5761015d6109c5565b604051908082528060200260200182016040528015610186578160200160208202803683370190505b50925060005b818110156101f7576101a9816a0422ca8b0a00a425000000610bea565b6101be906a0422ca8b0a00a425000000610c01565b846101ca600384610bd6565b815181106101da576101da610c14565b60209081029190910101526101f0600382610c01565b905061018c565b505050909192565b610207610597565b925092509250909192565b60408051600680825260e08201909252606091829160009182918291816020015b6102626040805160a0810182526000808252602082018190529181018290526060810182905290608082015290565b81526020019060019003908161023357505060408051600680825260e082019092529196506020820160c080368337019050506040805160a081018252610f008152610f016020820152610f029181019190915260006060820152909450608081016006815250856000815181106102dc576102dc610c14565b6020908102919091018101919091526040805160a081018252610f038152610f0492810192909252610f0590820152600060608201526080810160078152508560018151811061032e5761032e610c14565b6020908102919091018101919091526040805160a08101825261100081526110019281019290925261100290820152600060608201526080810160058152508560028151811061038057610380610c14565b6020908102919091018101919091526040805160a0810182526120008152612001928101929092526120029082015260006060820152608081016008815250856003815181106103d2576103d2610c14565b6020908102919091018101919091526040805160a08101825261300081526130019281019290925261300290820152600060608201526080810160028152508560048151811061042457610424610c14565b6020908102919091018101919091526040805160a08101825261400081526140019281019290925261400290820152600060608201526080810160068152508560058151811061047657610476610c14565b60200260200101819052506a0422ca8b0a00a425000000846000815181106104a0576104a0610c14565b6020026020010181815250506a084595161401484a000000846001815181106104cb576104cb610c14565b6020026020010181815250506a069e10de76676d08000000846002815181106104f6576104f6610c14565b6020026020010181815250506a05ca4ec2a79a7f670000008460038151811061052157610521610c14565b6020026020010181815250506a04f68ca6d8cd91c60000008460048151811061054c5761054c610c14565b6020026020010181815250506a0771d2fa45345aa90000008460058151811061057757610577610c14565b6020908102919091010152509293919250610a0191610a029150610a0390565b60408051600580825260c08201909252606091829182916020820160a080368337505060408051600580825260c0820190925292955090506020820160a080368337505060408051600180825281830190925292945090506020808301908036833701905050905060008360008151811061061457610614610c14565b602002602001019060ff16908160ff168152505060018360018151811061063d5761063d610c14565b602002602001019060ff16908160ff168152505060028360028151811061066657610666610c14565b602002602001019060ff16908160ff168152505060038360038151811061068f5761068f610c14565b602002602001019060ff16908160ff16815250506004836004815181106106b8576106b8610c14565b602002602001019060ff16908160ff1681525050610f00826000815181106106e2576106e2610c14565b60200260200101906001600160a01b031690816001600160a01b031681525050610f018260018151811061071857610718610c14565b60200260200101906001600160a01b031690816001600160a01b031681525050610f028260028151811061074e5761074e610c14565b60200260200101906001600160a01b031690816001600160a01b031681525050610f038260038151811061078457610784610c14565b60200260200101906001600160a01b031690816001600160a01b031681525050610f04826004815181106107ba576107ba610c14565b60200260200101906001600160a01b031690816001600160a01b0316815250506a05ca4ec2a79a7f67000000816000815181106107f9576107f9610c14565b602002602001018181525050909192565b600081518084526020808501945080840160005b8381101561083a5781518752958201959082019060010161081e565b509495945050505050565b606080825284519082018190526000906020906080840190828801845b8281101561088157815160ff1684529284019290840190600101610862565b5050508381038285015285518082528683019183019060005b818110156108bf5783516001600160a01b03168352928401929184019160010161089a565b505084810360408601526108d3818761080a565b98975050505050505050565b60a080825286518282018190526000919060209060c0850190828b0185805b8381101561096e57825180516001600160a01b039081168752878201518116888801526040808301519091169087015260608082015190870152608090810151906009821061095b57634e487b7160e01b84526021600452602484fd5b86015293860193918501916001016108fe565b5050505084810382860152610983818a61080a565b935050505061099d60408301866001600160a01b03169052565b6001600160a01b03841660608301526001600160a01b03831660808301529695505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff81118282101715610a0457610a046109c5565b604052919050565b600067ffffffffffffffff821115610a2657610a266109c5565b5060051b60200190565b600082601f830112610a4157600080fd5b81516020610a56610a5183610a0c565b6109db565b82815260059290921b84018101918181019086841115610a7557600080fd5b8286015b84811015610aa65780516001600160a01b0381168114610a995760008081fd5b8352918301918301610a79565b509695505050505050565b60008060408385031215610ac457600080fd5b825167ffffffffffffffff80821115610adc57600080fd5b818501915085601f830112610af057600080fd5b81516020610b00610a5183610a0c565b82815260059290921b84018101918181019089841115610b1f57600080fd5b948201945b83861015610b4d57855160ff81168114610b3e5760008081fd5b82529482019490820190610b24565b91880151919650909350505080821115610b6657600080fd5b50610b7385828601610a30565b9150509250929050565b634e487b7160e01b600052601160045260246000fd5b81810381811115610ba657610ba6610b7d565b92915050565b634e487b7160e01b600052601260045260246000fd5b600082610bd157610bd1610bac565b500690565b600082610be557610be5610bac565b500490565b8082028115828204841417610ba657610ba6610b7d565b80820180821115610ba657610ba6610b7d565b634e487b7160e01b600052603260045260246000fdfea26469706673582212204e29eb3d6f915a4170232fad30819e79eb67f3d1ce17d415a192d6672c5a7d6264736f6c63430008130033`

// Deprecated: Use MultiCallContractMockPermissionlessMetaData.Sigs instead.
// MultiCallContractMockPermissionlessFuncSigs maps the 4-byte function signature to its string representation.
var MultiCallContractMockPermissionlessFuncSigs = MultiCallContractMockPermissionlessMetaData.Sigs

// MultiCallContractMockPermissionlessBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MultiCallContractMockPermissionlessMetaData.Bin instead.
var MultiCallContractMockPermissionlessBin = MultiCallContractMockPermissionlessMetaData.Bin

// DeployMultiCallContractMockPermissionless deploys a new Kaia contract, binding an instance of MultiCallContractMockPermissionless to it.
func DeployMultiCallContractMockPermissionless(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MultiCallContractMockPermissionless, error) {
	parsed, err := MultiCallContractMockPermissionlessMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MultiCallContractMockPermissionlessBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultiCallContractMockPermissionless{MultiCallContractMockPermissionlessCaller: MultiCallContractMockPermissionlessCaller{contract: contract}, MultiCallContractMockPermissionlessTransactor: MultiCallContractMockPermissionlessTransactor{contract: contract}, MultiCallContractMockPermissionlessFilterer: MultiCallContractMockPermissionlessFilterer{contract: contract}}, nil
}

// MultiCallContractMockPermissionless is an auto generated Go binding around a Kaia contract.
type MultiCallContractMockPermissionless struct {
	MultiCallContractMockPermissionlessCaller     // Read-only binding to the contract
	MultiCallContractMockPermissionlessTransactor // Write-only binding to the contract
	MultiCallContractMockPermissionlessFilterer   // Log filterer for contract events
}

// MultiCallContractMockPermissionlessCaller is an auto generated read-only Go binding around a Kaia contract.
type MultiCallContractMockPermissionlessCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiCallContractMockPermissionlessTransactor is an auto generated write-only Go binding around a Kaia contract.
type MultiCallContractMockPermissionlessTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiCallContractMockPermissionlessFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type MultiCallContractMockPermissionlessFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiCallContractMockPermissionlessSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type MultiCallContractMockPermissionlessSession struct {
	Contract     *MultiCallContractMockPermissionless // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                        // Call options to use throughout this session
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// MultiCallContractMockPermissionlessCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type MultiCallContractMockPermissionlessCallerSession struct {
	Contract *MultiCallContractMockPermissionlessCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                              // Call options to use throughout this session
}

// MultiCallContractMockPermissionlessTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type MultiCallContractMockPermissionlessTransactorSession struct {
	Contract     *MultiCallContractMockPermissionlessTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                              // Transaction auth options to use throughout this session
}

// MultiCallContractMockPermissionlessRaw is an auto generated low-level Go binding around a Kaia contract.
type MultiCallContractMockPermissionlessRaw struct {
	Contract *MultiCallContractMockPermissionless // Generic contract binding to access the raw methods on
}

// MultiCallContractMockPermissionlessCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type MultiCallContractMockPermissionlessCallerRaw struct {
	Contract *MultiCallContractMockPermissionlessCaller // Generic read-only contract binding to access the raw methods on
}

// MultiCallContractMockPermissionlessTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type MultiCallContractMockPermissionlessTransactorRaw struct {
	Contract *MultiCallContractMockPermissionlessTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiCallContractMockPermissionless creates a new instance of MultiCallContractMockPermissionless, bound to a specific deployed contract.
func NewMultiCallContractMockPermissionless(address common.Address, backend bind.ContractBackend) (*MultiCallContractMockPermissionless, error) {
	contract, err := bindMultiCallContractMockPermissionless(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiCallContractMockPermissionless{MultiCallContractMockPermissionlessCaller: MultiCallContractMockPermissionlessCaller{contract: contract}, MultiCallContractMockPermissionlessTransactor: MultiCallContractMockPermissionlessTransactor{contract: contract}, MultiCallContractMockPermissionlessFilterer: MultiCallContractMockPermissionlessFilterer{contract: contract}}, nil
}

// NewMultiCallContractMockPermissionlessCaller creates a new read-only instance of MultiCallContractMockPermissionless, bound to a specific deployed contract.
func NewMultiCallContractMockPermissionlessCaller(address common.Address, caller bind.ContractCaller) (*MultiCallContractMockPermissionlessCaller, error) {
	contract, err := bindMultiCallContractMockPermissionless(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiCallContractMockPermissionlessCaller{contract: contract}, nil
}

// NewMultiCallContractMockPermissionlessTransactor creates a new write-only instance of MultiCallContractMockPermissionless, bound to a specific deployed contract.
func NewMultiCallContractMockPermissionlessTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiCallContractMockPermissionlessTransactor, error) {
	contract, err := bindMultiCallContractMockPermissionless(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiCallContractMockPermissionlessTransactor{contract: contract}, nil
}

// NewMultiCallContractMockPermissionlessFilterer creates a new log filterer instance of MultiCallContractMockPermissionless, bound to a specific deployed contract.
func NewMultiCallContractMockPermissionlessFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiCallContractMockPermissionlessFilterer, error) {
	contract, err := bindMultiCallContractMockPermissionless(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiCallContractMockPermissionlessFilterer{contract: contract}, nil
}

// bindMultiCallContractMockPermissionless binds a generic wrapper to an already deployed contract.
func bindMultiCallContractMockPermissionless(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultiCallContractMockPermissionlessMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiCallContractMockPermissionless.Contract.MultiCallContractMockPermissionlessCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiCallContractMockPermissionless.Contract.MultiCallContractMockPermissionlessTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiCallContractMockPermissionless.Contract.MultiCallContractMockPermissionlessTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiCallContractMockPermissionless.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiCallContractMockPermissionless.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiCallContractMockPermissionless.Contract.contract.Transact(opts, method, params...)
}

// MultiCallStakingInfo is a free data retrieval call binding the contract method 0xadde19c6.
//
// Solidity: function multiCallStakingInfo() view returns(uint8[] typeList, address[] addressList, uint256[] stakingAmounts)
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessCaller) MultiCallStakingInfo(opts *bind.CallOpts) (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
}, error) {
	var out []interface{}
	err := _MultiCallContractMockPermissionless.contract.Call(opts, &out, "multiCallStakingInfo")

	outstruct := new(struct {
		TypeList       []uint8
		AddressList    []common.Address
		StakingAmounts []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TypeList = *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)
	outstruct.AddressList = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.StakingAmounts = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// MultiCallStakingInfo is a free data retrieval call binding the contract method 0xadde19c6.
//
// Solidity: function multiCallStakingInfo() view returns(uint8[] typeList, address[] addressList, uint256[] stakingAmounts)
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessSession) MultiCallStakingInfo() (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
}, error) {
	return _MultiCallContractMockPermissionless.Contract.MultiCallStakingInfo(&_MultiCallContractMockPermissionless.CallOpts)
}

// MultiCallStakingInfo is a free data retrieval call binding the contract method 0xadde19c6.
//
// Solidity: function multiCallStakingInfo() view returns(uint8[] typeList, address[] addressList, uint256[] stakingAmounts)
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessCallerSession) MultiCallStakingInfo() (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
}, error) {
	return _MultiCallContractMockPermissionless.Contract.MultiCallStakingInfo(&_MultiCallContractMockPermissionless.CallOpts)
}

// MultiCallStakingInfoPermissionless is a free data retrieval call binding the contract method 0xb04ba218.
//
// Solidity: function multiCallStakingInfoPermissionless() pure returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, address kefAddr, address kifAddr, address kpfAddr)
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessCaller) MultiCallStakingInfoPermissionless(opts *bind.CallOpts) (struct {
	Profiles       []Profile
	StakingAmounts []*big.Int
	KefAddr        common.Address
	KifAddr        common.Address
	KpfAddr        common.Address
}, error) {
	var out []interface{}
	err := _MultiCallContractMockPermissionless.contract.Call(opts, &out, "multiCallStakingInfoPermissionless")

	outstruct := new(struct {
		Profiles       []Profile
		StakingAmounts []*big.Int
		KefAddr        common.Address
		KifAddr        common.Address
		KpfAddr        common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Profiles = *abi.ConvertType(out[0], new([]Profile)).(*[]Profile)
	outstruct.StakingAmounts = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.KefAddr = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.KifAddr = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.KpfAddr = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// MultiCallStakingInfoPermissionless is a free data retrieval call binding the contract method 0xb04ba218.
//
// Solidity: function multiCallStakingInfoPermissionless() pure returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, address kefAddr, address kifAddr, address kpfAddr)
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessSession) MultiCallStakingInfoPermissionless() (struct {
	Profiles       []Profile
	StakingAmounts []*big.Int
	KefAddr        common.Address
	KifAddr        common.Address
	KpfAddr        common.Address
}, error) {
	return _MultiCallContractMockPermissionless.Contract.MultiCallStakingInfoPermissionless(&_MultiCallContractMockPermissionless.CallOpts)
}

// MultiCallStakingInfoPermissionless is a free data retrieval call binding the contract method 0xb04ba218.
//
// Solidity: function multiCallStakingInfoPermissionless() pure returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, address kefAddr, address kifAddr, address kpfAddr)
func (_MultiCallContractMockPermissionless *MultiCallContractMockPermissionlessCallerSession) MultiCallStakingInfoPermissionless() (struct {
	Profiles       []Profile
	StakingAmounts []*big.Int
	KefAddr        common.Address
	KifAddr        common.Address
	KpfAddr        common.Address
}, error) {
	return _MultiCallContractMockPermissionless.Contract.MultiCallStakingInfoPermissionless(&_MultiCallContractMockPermissionless.CallOpts)
}
