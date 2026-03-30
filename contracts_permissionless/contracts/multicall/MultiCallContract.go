// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package multicall

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
	ABI: "[{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAddress\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isActivated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"spareContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"ffa1ad74": "VERSION()",
		"715b208b": "getAllAddress()",
		"4a8c1fb4": "isActivated()",
		"6abd623d": "spareContractAddress()",
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

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_IAddressBook *IAddressBookCaller) VERSION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_IAddressBook *IAddressBookSession) VERSION() (*big.Int, error) {
	return _IAddressBook.Contract.VERSION(&_IAddressBook.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_IAddressBook *IAddressBookCallerSession) VERSION() (*big.Int, error) {
	return _IAddressBook.Contract.VERSION(&_IAddressBook.CallOpts)
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

// IAddressBookV2MetaData contains all meta data concerning the IAddressBookV2 contract.
var IAddressBookV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllProfiles\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structProfile[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExitThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFundAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMaxCounts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"maxValidatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxReadyCandidateCount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTimeouts\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"pauseTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idleTimeout\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"0b1fe784": "getAllProfiles()",
		"3aea9faf": "getExitThreshold()",
		"25cf0943": "getFundAddresses()",
		"03e6689d": "getMaxCounts()",
		"e70c38f1": "getTimeouts()",
	},
}

// IAddressBookV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use IAddressBookV2MetaData.ABI instead.
var IAddressBookV2ABI = IAddressBookV2MetaData.ABI

// IAddressBookV2BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IAddressBookV2BinRuntime = ``

// Deprecated: Use IAddressBookV2MetaData.Sigs instead.
// IAddressBookV2FuncSigs maps the 4-byte function signature to its string representation.
var IAddressBookV2FuncSigs = IAddressBookV2MetaData.Sigs

// IAddressBookV2 is an auto generated Go binding around a Kaia contract.
type IAddressBookV2 struct {
	IAddressBookV2Caller     // Read-only binding to the contract
	IAddressBookV2Transactor // Write-only binding to the contract
	IAddressBookV2Filterer   // Log filterer for contract events
}

// IAddressBookV2Caller is an auto generated read-only Go binding around a Kaia contract.
type IAddressBookV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookV2Transactor is an auto generated write-only Go binding around a Kaia contract.
type IAddressBookV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookV2Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type IAddressBookV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAddressBookV2Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IAddressBookV2Session struct {
	Contract     *IAddressBookV2   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAddressBookV2CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IAddressBookV2CallerSession struct {
	Contract *IAddressBookV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IAddressBookV2TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IAddressBookV2TransactorSession struct {
	Contract     *IAddressBookV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IAddressBookV2Raw is an auto generated low-level Go binding around a Kaia contract.
type IAddressBookV2Raw struct {
	Contract *IAddressBookV2 // Generic contract binding to access the raw methods on
}

// IAddressBookV2CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IAddressBookV2CallerRaw struct {
	Contract *IAddressBookV2Caller // Generic read-only contract binding to access the raw methods on
}

// IAddressBookV2TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IAddressBookV2TransactorRaw struct {
	Contract *IAddressBookV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIAddressBookV2 creates a new instance of IAddressBookV2, bound to a specific deployed contract.
func NewIAddressBookV2(address common.Address, backend bind.ContractBackend) (*IAddressBookV2, error) {
	contract, err := bindIAddressBookV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAddressBookV2{IAddressBookV2Caller: IAddressBookV2Caller{contract: contract}, IAddressBookV2Transactor: IAddressBookV2Transactor{contract: contract}, IAddressBookV2Filterer: IAddressBookV2Filterer{contract: contract}}, nil
}

// NewIAddressBookV2Caller creates a new read-only instance of IAddressBookV2, bound to a specific deployed contract.
func NewIAddressBookV2Caller(address common.Address, caller bind.ContractCaller) (*IAddressBookV2Caller, error) {
	contract, err := bindIAddressBookV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAddressBookV2Caller{contract: contract}, nil
}

// NewIAddressBookV2Transactor creates a new write-only instance of IAddressBookV2, bound to a specific deployed contract.
func NewIAddressBookV2Transactor(address common.Address, transactor bind.ContractTransactor) (*IAddressBookV2Transactor, error) {
	contract, err := bindIAddressBookV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAddressBookV2Transactor{contract: contract}, nil
}

// NewIAddressBookV2Filterer creates a new log filterer instance of IAddressBookV2, bound to a specific deployed contract.
func NewIAddressBookV2Filterer(address common.Address, filterer bind.ContractFilterer) (*IAddressBookV2Filterer, error) {
	contract, err := bindIAddressBookV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAddressBookV2Filterer{contract: contract}, nil
}

// bindIAddressBookV2 binds a generic wrapper to an already deployed contract.
func bindIAddressBookV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAddressBookV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAddressBookV2 *IAddressBookV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAddressBookV2.Contract.IAddressBookV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAddressBookV2 *IAddressBookV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBookV2.Contract.IAddressBookV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAddressBookV2 *IAddressBookV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAddressBookV2.Contract.IAddressBookV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAddressBookV2 *IAddressBookV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAddressBookV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAddressBookV2 *IAddressBookV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAddressBookV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAddressBookV2 *IAddressBookV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAddressBookV2.Contract.contract.Transact(opts, method, params...)
}

// GetAllProfiles is a free data retrieval call binding the contract method 0x0b1fe784.
//
// Solidity: function getAllProfiles() view returns((address,address,address,uint256,uint8)[])
func (_IAddressBookV2 *IAddressBookV2Caller) GetAllProfiles(opts *bind.CallOpts) ([]Profile, error) {
	var out []interface{}
	err := _IAddressBookV2.contract.Call(opts, &out, "getAllProfiles")

	if err != nil {
		return *new([]Profile), err
	}

	out0 := *abi.ConvertType(out[0], new([]Profile)).(*[]Profile)

	return out0, err

}

// GetAllProfiles is a free data retrieval call binding the contract method 0x0b1fe784.
//
// Solidity: function getAllProfiles() view returns((address,address,address,uint256,uint8)[])
func (_IAddressBookV2 *IAddressBookV2Session) GetAllProfiles() ([]Profile, error) {
	return _IAddressBookV2.Contract.GetAllProfiles(&_IAddressBookV2.CallOpts)
}

// GetAllProfiles is a free data retrieval call binding the contract method 0x0b1fe784.
//
// Solidity: function getAllProfiles() view returns((address,address,address,uint256,uint8)[])
func (_IAddressBookV2 *IAddressBookV2CallerSession) GetAllProfiles() ([]Profile, error) {
	return _IAddressBookV2.Contract.GetAllProfiles(&_IAddressBookV2.CallOpts)
}

// GetExitThreshold is a free data retrieval call binding the contract method 0x3aea9faf.
//
// Solidity: function getExitThreshold() view returns(uint256)
func (_IAddressBookV2 *IAddressBookV2Caller) GetExitThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAddressBookV2.contract.Call(opts, &out, "getExitThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExitThreshold is a free data retrieval call binding the contract method 0x3aea9faf.
//
// Solidity: function getExitThreshold() view returns(uint256)
func (_IAddressBookV2 *IAddressBookV2Session) GetExitThreshold() (*big.Int, error) {
	return _IAddressBookV2.Contract.GetExitThreshold(&_IAddressBookV2.CallOpts)
}

// GetExitThreshold is a free data retrieval call binding the contract method 0x3aea9faf.
//
// Solidity: function getExitThreshold() view returns(uint256)
func (_IAddressBookV2 *IAddressBookV2CallerSession) GetExitThreshold() (*big.Int, error) {
	return _IAddressBookV2.Contract.GetExitThreshold(&_IAddressBookV2.CallOpts)
}

// GetFundAddresses is a free data retrieval call binding the contract method 0x25cf0943.
//
// Solidity: function getFundAddresses() view returns(address, address, address)
func (_IAddressBookV2 *IAddressBookV2Caller) GetFundAddresses(opts *bind.CallOpts) (common.Address, common.Address, common.Address, error) {
	var out []interface{}
	err := _IAddressBookV2.contract.Call(opts, &out, "getFundAddresses")

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
func (_IAddressBookV2 *IAddressBookV2Session) GetFundAddresses() (common.Address, common.Address, common.Address, error) {
	return _IAddressBookV2.Contract.GetFundAddresses(&_IAddressBookV2.CallOpts)
}

// GetFundAddresses is a free data retrieval call binding the contract method 0x25cf0943.
//
// Solidity: function getFundAddresses() view returns(address, address, address)
func (_IAddressBookV2 *IAddressBookV2CallerSession) GetFundAddresses() (common.Address, common.Address, common.Address, error) {
	return _IAddressBookV2.Contract.GetFundAddresses(&_IAddressBookV2.CallOpts)
}

// GetMaxCounts is a free data retrieval call binding the contract method 0x03e6689d.
//
// Solidity: function getMaxCounts() view returns(uint256 maxValidatorCount, uint256 maxReadyCandidateCount)
func (_IAddressBookV2 *IAddressBookV2Caller) GetMaxCounts(opts *bind.CallOpts) (struct {
	MaxValidatorCount      *big.Int
	MaxReadyCandidateCount *big.Int
}, error) {
	var out []interface{}
	err := _IAddressBookV2.contract.Call(opts, &out, "getMaxCounts")

	outstruct := new(struct {
		MaxValidatorCount      *big.Int
		MaxReadyCandidateCount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.MaxValidatorCount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.MaxReadyCandidateCount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetMaxCounts is a free data retrieval call binding the contract method 0x03e6689d.
//
// Solidity: function getMaxCounts() view returns(uint256 maxValidatorCount, uint256 maxReadyCandidateCount)
func (_IAddressBookV2 *IAddressBookV2Session) GetMaxCounts() (struct {
	MaxValidatorCount      *big.Int
	MaxReadyCandidateCount *big.Int
}, error) {
	return _IAddressBookV2.Contract.GetMaxCounts(&_IAddressBookV2.CallOpts)
}

// GetMaxCounts is a free data retrieval call binding the contract method 0x03e6689d.
//
// Solidity: function getMaxCounts() view returns(uint256 maxValidatorCount, uint256 maxReadyCandidateCount)
func (_IAddressBookV2 *IAddressBookV2CallerSession) GetMaxCounts() (struct {
	MaxValidatorCount      *big.Int
	MaxReadyCandidateCount *big.Int
}, error) {
	return _IAddressBookV2.Contract.GetMaxCounts(&_IAddressBookV2.CallOpts)
}

// GetTimeouts is a free data retrieval call binding the contract method 0xe70c38f1.
//
// Solidity: function getTimeouts() view returns(uint256 pauseTimeout, uint256 idleTimeout)
func (_IAddressBookV2 *IAddressBookV2Caller) GetTimeouts(opts *bind.CallOpts) (struct {
	PauseTimeout *big.Int
	IdleTimeout  *big.Int
}, error) {
	var out []interface{}
	err := _IAddressBookV2.contract.Call(opts, &out, "getTimeouts")

	outstruct := new(struct {
		PauseTimeout *big.Int
		IdleTimeout  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PauseTimeout = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.IdleTimeout = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetTimeouts is a free data retrieval call binding the contract method 0xe70c38f1.
//
// Solidity: function getTimeouts() view returns(uint256 pauseTimeout, uint256 idleTimeout)
func (_IAddressBookV2 *IAddressBookV2Session) GetTimeouts() (struct {
	PauseTimeout *big.Int
	IdleTimeout  *big.Int
}, error) {
	return _IAddressBookV2.Contract.GetTimeouts(&_IAddressBookV2.CallOpts)
}

// GetTimeouts is a free data retrieval call binding the contract method 0xe70c38f1.
//
// Solidity: function getTimeouts() view returns(uint256 pauseTimeout, uint256 idleTimeout)
func (_IAddressBookV2 *IAddressBookV2CallerSession) GetTimeouts() (struct {
	PauseTimeout *big.Int
	IdleTimeout  *big.Int
}, error) {
	return _IAddressBookV2.Contract.GetTimeouts(&_IAddressBookV2.CallOpts)
}

// ICLRegistryMetaData contains all meta data concerning the ICLRegistry contract.
var ICLRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllCLs\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"90599c07": "getAllCLs()",
	},
}

// ICLRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use ICLRegistryMetaData.ABI instead.
var ICLRegistryABI = ICLRegistryMetaData.ABI

// ICLRegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ICLRegistryBinRuntime = ``

// Deprecated: Use ICLRegistryMetaData.Sigs instead.
// ICLRegistryFuncSigs maps the 4-byte function signature to its string representation.
var ICLRegistryFuncSigs = ICLRegistryMetaData.Sigs

// ICLRegistry is an auto generated Go binding around a Kaia contract.
type ICLRegistry struct {
	ICLRegistryCaller     // Read-only binding to the contract
	ICLRegistryTransactor // Write-only binding to the contract
	ICLRegistryFilterer   // Log filterer for contract events
}

// ICLRegistryCaller is an auto generated read-only Go binding around a Kaia contract.
type ICLRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICLRegistryTransactor is an auto generated write-only Go binding around a Kaia contract.
type ICLRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICLRegistryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ICLRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICLRegistrySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ICLRegistrySession struct {
	Contract     *ICLRegistry      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ICLRegistryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ICLRegistryCallerSession struct {
	Contract *ICLRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ICLRegistryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ICLRegistryTransactorSession struct {
	Contract     *ICLRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ICLRegistryRaw is an auto generated low-level Go binding around a Kaia contract.
type ICLRegistryRaw struct {
	Contract *ICLRegistry // Generic contract binding to access the raw methods on
}

// ICLRegistryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ICLRegistryCallerRaw struct {
	Contract *ICLRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// ICLRegistryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ICLRegistryTransactorRaw struct {
	Contract *ICLRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewICLRegistry creates a new instance of ICLRegistry, bound to a specific deployed contract.
func NewICLRegistry(address common.Address, backend bind.ContractBackend) (*ICLRegistry, error) {
	contract, err := bindICLRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ICLRegistry{ICLRegistryCaller: ICLRegistryCaller{contract: contract}, ICLRegistryTransactor: ICLRegistryTransactor{contract: contract}, ICLRegistryFilterer: ICLRegistryFilterer{contract: contract}}, nil
}

// NewICLRegistryCaller creates a new read-only instance of ICLRegistry, bound to a specific deployed contract.
func NewICLRegistryCaller(address common.Address, caller bind.ContractCaller) (*ICLRegistryCaller, error) {
	contract, err := bindICLRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ICLRegistryCaller{contract: contract}, nil
}

// NewICLRegistryTransactor creates a new write-only instance of ICLRegistry, bound to a specific deployed contract.
func NewICLRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*ICLRegistryTransactor, error) {
	contract, err := bindICLRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ICLRegistryTransactor{contract: contract}, nil
}

// NewICLRegistryFilterer creates a new log filterer instance of ICLRegistry, bound to a specific deployed contract.
func NewICLRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*ICLRegistryFilterer, error) {
	contract, err := bindICLRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ICLRegistryFilterer{contract: contract}, nil
}

// bindICLRegistry binds a generic wrapper to an already deployed contract.
func bindICLRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ICLRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICLRegistry *ICLRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICLRegistry.Contract.ICLRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICLRegistry *ICLRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICLRegistry.Contract.ICLRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICLRegistry *ICLRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICLRegistry.Contract.ICLRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICLRegistry *ICLRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICLRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICLRegistry *ICLRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICLRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICLRegistry *ICLRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICLRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetAllCLs is a free data retrieval call binding the contract method 0x90599c07.
//
// Solidity: function getAllCLs() view returns(address[], uint256[], address[])
func (_ICLRegistry *ICLRegistryCaller) GetAllCLs(opts *bind.CallOpts) ([]common.Address, []*big.Int, []common.Address, error) {
	var out []interface{}
	err := _ICLRegistry.contract.Call(opts, &out, "getAllCLs")

	if err != nil {
		return *new([]common.Address), *new([]*big.Int), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	out2 := *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)

	return out0, out1, out2, err

}

// GetAllCLs is a free data retrieval call binding the contract method 0x90599c07.
//
// Solidity: function getAllCLs() view returns(address[], uint256[], address[])
func (_ICLRegistry *ICLRegistrySession) GetAllCLs() ([]common.Address, []*big.Int, []common.Address, error) {
	return _ICLRegistry.Contract.GetAllCLs(&_ICLRegistry.CallOpts)
}

// GetAllCLs is a free data retrieval call binding the contract method 0x90599c07.
//
// Solidity: function getAllCLs() view returns(address[], uint256[], address[])
func (_ICLRegistry *ICLRegistryCallerSession) GetAllCLs() ([]common.Address, []*big.Int, []common.Address, error) {
	return _ICLRegistry.Contract.GetAllCLs(&_ICLRegistry.CallOpts)
}

// ICnStakingMetaData contains all meta data concerning the ICnStaking contract.
var ICnStakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staking\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstaking\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"ffa1ad74": "VERSION()",
		"4cf088d9": "staking()",
		"630b1146": "unstaking()",
	},
}

// ICnStakingABI is the input ABI used to generate the binding from.
// Deprecated: Use ICnStakingMetaData.ABI instead.
var ICnStakingABI = ICnStakingMetaData.ABI

// ICnStakingBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ICnStakingBinRuntime = ``

// Deprecated: Use ICnStakingMetaData.Sigs instead.
// ICnStakingFuncSigs maps the 4-byte function signature to its string representation.
var ICnStakingFuncSigs = ICnStakingMetaData.Sigs

// ICnStaking is an auto generated Go binding around a Kaia contract.
type ICnStaking struct {
	ICnStakingCaller     // Read-only binding to the contract
	ICnStakingTransactor // Write-only binding to the contract
	ICnStakingFilterer   // Log filterer for contract events
}

// ICnStakingCaller is an auto generated read-only Go binding around a Kaia contract.
type ICnStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICnStakingTransactor is an auto generated write-only Go binding around a Kaia contract.
type ICnStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICnStakingFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ICnStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICnStakingSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ICnStakingSession struct {
	Contract     *ICnStaking       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ICnStakingCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ICnStakingCallerSession struct {
	Contract *ICnStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// ICnStakingTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ICnStakingTransactorSession struct {
	Contract     *ICnStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// ICnStakingRaw is an auto generated low-level Go binding around a Kaia contract.
type ICnStakingRaw struct {
	Contract *ICnStaking // Generic contract binding to access the raw methods on
}

// ICnStakingCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ICnStakingCallerRaw struct {
	Contract *ICnStakingCaller // Generic read-only contract binding to access the raw methods on
}

// ICnStakingTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ICnStakingTransactorRaw struct {
	Contract *ICnStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewICnStaking creates a new instance of ICnStaking, bound to a specific deployed contract.
func NewICnStaking(address common.Address, backend bind.ContractBackend) (*ICnStaking, error) {
	contract, err := bindICnStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ICnStaking{ICnStakingCaller: ICnStakingCaller{contract: contract}, ICnStakingTransactor: ICnStakingTransactor{contract: contract}, ICnStakingFilterer: ICnStakingFilterer{contract: contract}}, nil
}

// NewICnStakingCaller creates a new read-only instance of ICnStaking, bound to a specific deployed contract.
func NewICnStakingCaller(address common.Address, caller bind.ContractCaller) (*ICnStakingCaller, error) {
	contract, err := bindICnStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ICnStakingCaller{contract: contract}, nil
}

// NewICnStakingTransactor creates a new write-only instance of ICnStaking, bound to a specific deployed contract.
func NewICnStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*ICnStakingTransactor, error) {
	contract, err := bindICnStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ICnStakingTransactor{contract: contract}, nil
}

// NewICnStakingFilterer creates a new log filterer instance of ICnStaking, bound to a specific deployed contract.
func NewICnStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*ICnStakingFilterer, error) {
	contract, err := bindICnStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ICnStakingFilterer{contract: contract}, nil
}

// bindICnStaking binds a generic wrapper to an already deployed contract.
func bindICnStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ICnStakingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICnStaking *ICnStakingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICnStaking.Contract.ICnStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICnStaking *ICnStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICnStaking.Contract.ICnStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICnStaking *ICnStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICnStaking.Contract.ICnStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICnStaking *ICnStakingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICnStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICnStaking *ICnStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICnStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICnStaking *ICnStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICnStaking.Contract.contract.Transact(opts, method, params...)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_ICnStaking *ICnStakingCaller) VERSION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_ICnStaking *ICnStakingSession) VERSION() (*big.Int, error) {
	return _ICnStaking.Contract.VERSION(&_ICnStaking.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_ICnStaking *ICnStakingCallerSession) VERSION() (*big.Int, error) {
	return _ICnStaking.Contract.VERSION(&_ICnStaking.CallOpts)
}

// Staking is a free data retrieval call binding the contract method 0x4cf088d9.
//
// Solidity: function staking() view returns(uint256)
func (_ICnStaking *ICnStakingCaller) Staking(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "staking")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Staking is a free data retrieval call binding the contract method 0x4cf088d9.
//
// Solidity: function staking() view returns(uint256)
func (_ICnStaking *ICnStakingSession) Staking() (*big.Int, error) {
	return _ICnStaking.Contract.Staking(&_ICnStaking.CallOpts)
}

// Staking is a free data retrieval call binding the contract method 0x4cf088d9.
//
// Solidity: function staking() view returns(uint256)
func (_ICnStaking *ICnStakingCallerSession) Staking() (*big.Int, error) {
	return _ICnStaking.Contract.Staking(&_ICnStaking.CallOpts)
}

// Unstaking is a free data retrieval call binding the contract method 0x630b1146.
//
// Solidity: function unstaking() view returns(uint256)
func (_ICnStaking *ICnStakingCaller) Unstaking(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "unstaking")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Unstaking is a free data retrieval call binding the contract method 0x630b1146.
//
// Solidity: function unstaking() view returns(uint256)
func (_ICnStaking *ICnStakingSession) Unstaking() (*big.Int, error) {
	return _ICnStaking.Contract.Unstaking(&_ICnStaking.CallOpts)
}

// Unstaking is a free data retrieval call binding the contract method 0x630b1146.
//
// Solidity: function unstaking() view returns(uint256)
func (_ICnStaking *ICnStakingCallerSession) Unstaking() (*big.Int, error) {
	return _ICnStaking.Contract.Unstaking(&_ICnStaking.CallOpts)
}

// IERC20MetaData contains all meta data concerning the IERC20 contract.
var IERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"70a08231": "balanceOf(address)",
	},
}

// IERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC20MetaData.ABI instead.
var IERC20ABI = IERC20MetaData.ABI

// IERC20BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC20BinRuntime = ``

// Deprecated: Use IERC20MetaData.Sigs instead.
// IERC20FuncSigs maps the 4-byte function signature to its string representation.
var IERC20FuncSigs = IERC20MetaData.Sigs

// IERC20 is an auto generated Go binding around a Kaia contract.
type IERC20 struct {
	IERC20Caller     // Read-only binding to the contract
	IERC20Transactor // Write-only binding to the contract
	IERC20Filterer   // Log filterer for contract events
}

// IERC20Caller is an auto generated read-only Go binding around a Kaia contract.
type IERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Transactor is an auto generated write-only Go binding around a Kaia contract.
type IERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type IERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IERC20Session struct {
	Contract     *IERC20           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IERC20CallerSession struct {
	Contract *IERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// IERC20TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IERC20TransactorSession struct {
	Contract     *IERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20Raw is an auto generated low-level Go binding around a Kaia contract.
type IERC20Raw struct {
	Contract *IERC20 // Generic contract binding to access the raw methods on
}

// IERC20CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IERC20CallerRaw struct {
	Contract *IERC20Caller // Generic read-only contract binding to access the raw methods on
}

// IERC20TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IERC20TransactorRaw struct {
	Contract *IERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC20 creates a new instance of IERC20, bound to a specific deployed contract.
func NewIERC20(address common.Address, backend bind.ContractBackend) (*IERC20, error) {
	contract, err := bindIERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20{IERC20Caller: IERC20Caller{contract: contract}, IERC20Transactor: IERC20Transactor{contract: contract}, IERC20Filterer: IERC20Filterer{contract: contract}}, nil
}

// NewIERC20Caller creates a new read-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Caller(address common.Address, caller bind.ContractCaller) (*IERC20Caller, error) {
	contract, err := bindIERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Caller{contract: contract}, nil
}

// NewIERC20Transactor creates a new write-only instance of IERC20, bound to a specific deployed contract.
func NewIERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC20Transactor, error) {
	contract, err := bindIERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20Transactor{contract: contract}, nil
}

// NewIERC20Filterer creates a new log filterer instance of IERC20, bound to a specific deployed contract.
func NewIERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC20Filterer, error) {
	contract, err := bindIERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC20Filterer{contract: contract}, nil
}

// bindIERC20 binds a generic wrapper to an already deployed contract.
func bindIERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20.Contract.IERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.IERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20 *IERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20 *IERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20 *IERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20 *IERC20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, account)
}

// IGaslessSwapRouterMetaData contains all meta data concerning the IGaslessSwapRouter contract.
var IGaslessSwapRouterMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getSupportedTokens\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"d3c7c2c7": "getSupportedTokens()",
	},
}

// IGaslessSwapRouterABI is the input ABI used to generate the binding from.
// Deprecated: Use IGaslessSwapRouterMetaData.ABI instead.
var IGaslessSwapRouterABI = IGaslessSwapRouterMetaData.ABI

// IGaslessSwapRouterBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IGaslessSwapRouterBinRuntime = ``

// Deprecated: Use IGaslessSwapRouterMetaData.Sigs instead.
// IGaslessSwapRouterFuncSigs maps the 4-byte function signature to its string representation.
var IGaslessSwapRouterFuncSigs = IGaslessSwapRouterMetaData.Sigs

// IGaslessSwapRouter is an auto generated Go binding around a Kaia contract.
type IGaslessSwapRouter struct {
	IGaslessSwapRouterCaller     // Read-only binding to the contract
	IGaslessSwapRouterTransactor // Write-only binding to the contract
	IGaslessSwapRouterFilterer   // Log filterer for contract events
}

// IGaslessSwapRouterCaller is an auto generated read-only Go binding around a Kaia contract.
type IGaslessSwapRouterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGaslessSwapRouterTransactor is an auto generated write-only Go binding around a Kaia contract.
type IGaslessSwapRouterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGaslessSwapRouterFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IGaslessSwapRouterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IGaslessSwapRouterSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IGaslessSwapRouterSession struct {
	Contract     *IGaslessSwapRouter // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IGaslessSwapRouterCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IGaslessSwapRouterCallerSession struct {
	Contract *IGaslessSwapRouterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// IGaslessSwapRouterTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IGaslessSwapRouterTransactorSession struct {
	Contract     *IGaslessSwapRouterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IGaslessSwapRouterRaw is an auto generated low-level Go binding around a Kaia contract.
type IGaslessSwapRouterRaw struct {
	Contract *IGaslessSwapRouter // Generic contract binding to access the raw methods on
}

// IGaslessSwapRouterCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IGaslessSwapRouterCallerRaw struct {
	Contract *IGaslessSwapRouterCaller // Generic read-only contract binding to access the raw methods on
}

// IGaslessSwapRouterTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IGaslessSwapRouterTransactorRaw struct {
	Contract *IGaslessSwapRouterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIGaslessSwapRouter creates a new instance of IGaslessSwapRouter, bound to a specific deployed contract.
func NewIGaslessSwapRouter(address common.Address, backend bind.ContractBackend) (*IGaslessSwapRouter, error) {
	contract, err := bindIGaslessSwapRouter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IGaslessSwapRouter{IGaslessSwapRouterCaller: IGaslessSwapRouterCaller{contract: contract}, IGaslessSwapRouterTransactor: IGaslessSwapRouterTransactor{contract: contract}, IGaslessSwapRouterFilterer: IGaslessSwapRouterFilterer{contract: contract}}, nil
}

// NewIGaslessSwapRouterCaller creates a new read-only instance of IGaslessSwapRouter, bound to a specific deployed contract.
func NewIGaslessSwapRouterCaller(address common.Address, caller bind.ContractCaller) (*IGaslessSwapRouterCaller, error) {
	contract, err := bindIGaslessSwapRouter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IGaslessSwapRouterCaller{contract: contract}, nil
}

// NewIGaslessSwapRouterTransactor creates a new write-only instance of IGaslessSwapRouter, bound to a specific deployed contract.
func NewIGaslessSwapRouterTransactor(address common.Address, transactor bind.ContractTransactor) (*IGaslessSwapRouterTransactor, error) {
	contract, err := bindIGaslessSwapRouter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IGaslessSwapRouterTransactor{contract: contract}, nil
}

// NewIGaslessSwapRouterFilterer creates a new log filterer instance of IGaslessSwapRouter, bound to a specific deployed contract.
func NewIGaslessSwapRouterFilterer(address common.Address, filterer bind.ContractFilterer) (*IGaslessSwapRouterFilterer, error) {
	contract, err := bindIGaslessSwapRouter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IGaslessSwapRouterFilterer{contract: contract}, nil
}

// bindIGaslessSwapRouter binds a generic wrapper to an already deployed contract.
func bindIGaslessSwapRouter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IGaslessSwapRouterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IGaslessSwapRouter *IGaslessSwapRouterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IGaslessSwapRouter.Contract.IGaslessSwapRouterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IGaslessSwapRouter *IGaslessSwapRouterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGaslessSwapRouter.Contract.IGaslessSwapRouterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IGaslessSwapRouter *IGaslessSwapRouterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IGaslessSwapRouter.Contract.IGaslessSwapRouterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IGaslessSwapRouter *IGaslessSwapRouterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IGaslessSwapRouter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IGaslessSwapRouter *IGaslessSwapRouterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IGaslessSwapRouter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IGaslessSwapRouter *IGaslessSwapRouterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IGaslessSwapRouter.Contract.contract.Transact(opts, method, params...)
}

// GetSupportedTokens is a free data retrieval call binding the contract method 0xd3c7c2c7.
//
// Solidity: function getSupportedTokens() view returns(address[])
func (_IGaslessSwapRouter *IGaslessSwapRouterCaller) GetSupportedTokens(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _IGaslessSwapRouter.contract.Call(opts, &out, "getSupportedTokens")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetSupportedTokens is a free data retrieval call binding the contract method 0xd3c7c2c7.
//
// Solidity: function getSupportedTokens() view returns(address[])
func (_IGaslessSwapRouter *IGaslessSwapRouterSession) GetSupportedTokens() ([]common.Address, error) {
	return _IGaslessSwapRouter.Contract.GetSupportedTokens(&_IGaslessSwapRouter.CallOpts)
}

// GetSupportedTokens is a free data retrieval call binding the contract method 0xd3c7c2c7.
//
// Solidity: function getSupportedTokens() view returns(address[])
func (_IGaslessSwapRouter *IGaslessSwapRouterCallerSession) GetSupportedTokens() ([]common.Address, error) {
	return _IGaslessSwapRouter.Contract.GetSupportedTokens(&_IGaslessSwapRouter.CallOpts)
}

// IRegistryMetaData contains all meta data concerning the IRegistry contract.
var IRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e2693e3f": "getActiveAddr(string)",
	},
}

// IRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use IRegistryMetaData.ABI instead.
var IRegistryABI = IRegistryMetaData.ABI

// IRegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IRegistryBinRuntime = ``

// Deprecated: Use IRegistryMetaData.Sigs instead.
// IRegistryFuncSigs maps the 4-byte function signature to its string representation.
var IRegistryFuncSigs = IRegistryMetaData.Sigs

// IRegistry is an auto generated Go binding around a Kaia contract.
type IRegistry struct {
	IRegistryCaller     // Read-only binding to the contract
	IRegistryTransactor // Write-only binding to the contract
	IRegistryFilterer   // Log filterer for contract events
}

// IRegistryCaller is an auto generated read-only Go binding around a Kaia contract.
type IRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistryTransactor is an auto generated write-only Go binding around a Kaia contract.
type IRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistrySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IRegistrySession struct {
	Contract     *IRegistry        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IRegistryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IRegistryCallerSession struct {
	Contract *IRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IRegistryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IRegistryTransactorSession struct {
	Contract     *IRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IRegistryRaw is an auto generated low-level Go binding around a Kaia contract.
type IRegistryRaw struct {
	Contract *IRegistry // Generic contract binding to access the raw methods on
}

// IRegistryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IRegistryCallerRaw struct {
	Contract *IRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// IRegistryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IRegistryTransactorRaw struct {
	Contract *IRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIRegistry creates a new instance of IRegistry, bound to a specific deployed contract.
func NewIRegistry(address common.Address, backend bind.ContractBackend) (*IRegistry, error) {
	contract, err := bindIRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IRegistry{IRegistryCaller: IRegistryCaller{contract: contract}, IRegistryTransactor: IRegistryTransactor{contract: contract}, IRegistryFilterer: IRegistryFilterer{contract: contract}}, nil
}

// NewIRegistryCaller creates a new read-only instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryCaller(address common.Address, caller bind.ContractCaller) (*IRegistryCaller, error) {
	contract, err := bindIRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IRegistryCaller{contract: contract}, nil
}

// NewIRegistryTransactor creates a new write-only instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*IRegistryTransactor, error) {
	contract, err := bindIRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IRegistryTransactor{contract: contract}, nil
}

// NewIRegistryFilterer creates a new log filterer instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*IRegistryFilterer, error) {
	contract, err := bindIRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IRegistryFilterer{contract: contract}, nil
}

// bindIRegistry binds a generic wrapper to an already deployed contract.
func bindIRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRegistry *IRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IRegistry.Contract.IRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRegistry *IRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRegistry.Contract.IRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRegistry *IRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRegistry.Contract.IRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRegistry *IRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRegistry *IRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRegistry *IRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistryCaller) GetActiveAddr(opts *bind.CallOpts, name string) (common.Address, error) {
	var out []interface{}
	err := _IRegistry.contract.Call(opts, &out, "getActiveAddr", name)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistrySession) GetActiveAddr(name string) (common.Address, error) {
	return _IRegistry.Contract.GetActiveAddr(&_IRegistry.CallOpts, name)
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistryCallerSession) GetActiveAddr(name string) (common.Address, error) {
	return _IRegistry.Contract.GetActiveAddr(&_IRegistry.CallOpts, name)
}

// MultiCallContractMetaData contains all meta data concerning the MultiCallContract contract.
var MultiCallContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"multiCallDPStakingInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIds\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"clPools\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallGaslessInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"gsr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallNodeStatesPermissionless\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structProfile[]\",\"name\":\"profiles\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"pauseTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idleTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxValidatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxReadyCandidateCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"exitThreshold\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallStakingInfo\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"spareAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallStakingInfoPermissionless\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structProfile[]\",\"name\":\"profiles\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"kefAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kifAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kpfAddr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"6082579d": "multiCallDPStakingInfo()",
		"bfe8e683": "multiCallGaslessInfo()",
		"2ada6c5c": "multiCallNodeStatesPermissionless()",
		"adde19c6": "multiCallStakingInfo()",
		"b04ba218": "multiCallStakingInfoPermissionless()",
	},
	Bin: "0x608060405234801561001057600080fd5b50611298806100206000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80632ada6c5c1461005c5780636082579d14610080578063adde19c614610097578063b04ba218146100af578063bfe8e683146100c8575b600080fd5b6100646100de565b6040516100779796959493929190610b5b565b60405180910390f35b610088610237565b60405161007793929190610be4565b61009f6104f5565b6040516100779493929190610c27565b6100b7610513565b604051610077959493929190610ca6565b6100d06105a5565b604051610077929190610cf6565b60608060008060008060006100f16106ae565b6040805163e70c38f160e01b81528151939a5091985061040092839263e70c38f19260048083019391928290030181865afa158015610134573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101589190610d22565b8096508197505050806001600160a01b03166303e6689d6040518163ffffffff1660e01b81526004016040805180830381865afa15801561019d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101c19190610d22565b8094508195505050806001600160a01b0316633aea9faf6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610207573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061022b9190610d46565b91505090919293949596565b60405163e2693e3f60e01b815260206004820152600a602482015269434c526567697374727960b01b6044820152606090819081906000906104019063e2693e3f90606401602060405180830381865afa158015610299573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102bd9190610d7b565b60405163e2693e3f60e01b815260206004820152600b60248201526a577261707065644b61696160a81b60448201529091506000906104019063e2693e3f90606401602060405180830381865afa15801561031c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103409190610d7b565b90506001600160a01b038216610357575050909192565b816001600160a01b03166390599c076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610395573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103bd9190810190610ea3565b8051929750955081905067ffffffffffffffff8111156103df576103df610d9d565b604051908082528060200260200182016040528015610408578160200160208202803683370190505b5093506001600160a01b038216156104ed578160005b828110156104ea57816001600160a01b03166370a0823188838151811061044757610447610f80565b60200260200101516040518263ffffffff1660e01b815260040161047a91906001600160a01b0391909116815260200190565b602060405180830381865afa158015610497573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104bb9190610d46565b8682815181106104cd576104cd610f80565b6020908102919091010152806104e281610fac565b91505061041e565b50505b505050909192565b606080606060006105046107cf565b93509350935093505b90919293565b60608060008060006105236106ae565b809550819650505060006104009050806001600160a01b03166325cf09436040518163ffffffff1660e01b8152600401606060405180830381865afa158015610570573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105949190610fc5565b979896979196909550909350915050565b60405163e2693e3f60e01b815260206004820152601160248201527023b0b9b632b9b9a9bbb0b82937baba32b960791b60448201526000906060906104019063e2693e3f90606401602060405180830381865afa15801561060a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061062e9190610d7b565b91506001600160a01b038216610642579091565b816001600160a01b031663d3c7c2c76040518163ffffffff1660e01b8152600401600060405180830381865afa158015610680573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526106a89190810190611008565b90509091565b60608060006104009050806001600160a01b0316630b1fe7846040518163ffffffff1660e01b8152600401600060405180830381865afa1580156106f6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261071e919081019061103d565b80519093508067ffffffffffffffff81111561073c5761073c610d9d565b604051908082528060200260200182016040528015610765578160200160208202803683370190505b50925060005b818110156107c85761079985828151811061078857610788610f80565b6020026020010151602001516109b7565b8482815181106107ab576107ab610f80565b6020908102919091010152806107c081610fac565b91505061076b565b5050509091565b60608060606000806104009050806001600160a01b031663715b208b6040518163ffffffff1660e01b8152600401600060405180830381865afa15801561081a573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526108429190810190611132565b8095508196505050806001600160a01b0316636abd623d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610888573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108ac9190610d7b565b91506005845110156108be575061050d565b6000600285516108ce91906111fe565b90506108db600382611227565b156108e757505061050d565b6108f260038261123b565b67ffffffffffffffff81111561090a5761090a610d9d565b604051908082528060200260200182016040528015610933578160200160208202803683370190505b50935060005b818110156109ae576109758661095083600161124f565b8151811061096057610960610f80565b60200260200101516001600160a01b03163190565b8561098160038461123b565b8151811061099157610991610f80565b60209081029190910101526109a760038261124f565b9050610939565b50505090919293565b6000816001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109f7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a1b9190610d46565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a59573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a7d9190610d46565b610a8791906111fe565b92915050565b60008151808452602080850194508084016000805b84811015610b1457825180516001600160a01b039081168a52858201518116868b0152604080830151909116908a0152606080820151908a01526080908101519060098210610aff57634e487b7160e01b84526021600452602484fd5b89015260a09097019691830191600101610aa2565b50959695505050505050565b600081518084526020808501945080840160005b83811015610b5057815187529582019590820190600101610b34565b509495945050505050565b60e081526000610b6e60e083018a610a8d565b8281036020840152610b80818a610b20565b604084019890985250506060810194909452608084019290925260a083015260c09091015292915050565b600081518084526020808501945080840160005b83811015610b505781516001600160a01b031687529582019590820190600101610bbf565b606081526000610bf76060830186610bab565b8281036020840152610c098186610bab565b90508281036040840152610c1d8185610b20565b9695505050505050565b6080808252855190820181905260009060209060a0840190828901845b82811015610c6357815160ff1684529284019290840190600101610c44565b50505083810382850152610c778188610bab565b9150508281036040840152610c8c8186610b20565b91505060018060a01b038316606083015295945050505050565b60a081526000610cb960a0830188610a8d565b8281036020840152610ccb8188610b20565b6001600160a01b03968716604085015294861660608401525050921660809092019190915292915050565b6001600160a01b0383168152604060208201819052600090610d1a90830184610bab565b949350505050565b60008060408385031215610d3557600080fd5b505080516020909101519092909150565b600060208284031215610d5857600080fd5b5051919050565b80516001600160a01b0381168114610d7657600080fd5b919050565b600060208284031215610d8d57600080fd5b610d9682610d5f565b9392505050565b634e487b7160e01b600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715610dd657610dd6610d9d565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610e0557610e05610d9d565b604052919050565b600067ffffffffffffffff821115610e2757610e27610d9d565b5060051b60200190565b600082601f830112610e4257600080fd5b81516020610e57610e5283610e0d565b610ddc565b82815260059290921b84018101918181019086841115610e7657600080fd5b8286015b84811015610e9857610e8b81610d5f565b8352918301918301610e7a565b509695505050505050565b600080600060608486031215610eb857600080fd5b835167ffffffffffffffff80821115610ed057600080fd5b610edc87838801610e31565b9450602091508186015181811115610ef357600080fd5b8601601f81018813610f0457600080fd5b8051610f12610e5282610e0d565b81815260059190911b8201840190848101908a831115610f3157600080fd5b928501925b82841015610f4f57835182529285019290850190610f36565b60408a0151909750945050505080821115610f6957600080fd5b50610f7686828701610e31565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201610fbe57610fbe610f96565b5060010190565b600080600060608486031215610fda57600080fd5b610fe384610d5f565b9250610ff160208501610d5f565b9150610fff60408501610d5f565b90509250925092565b60006020828403121561101a57600080fd5b815167ffffffffffffffff81111561103157600080fd5b610d1a84828501610e31565b6000602080838503121561105057600080fd5b825167ffffffffffffffff81111561106757600080fd5b8301601f8101851361107857600080fd5b8051611086610e5282610e0d565b81815260a091820283018401918482019190888411156110a557600080fd5b938501935b838510156111265780858a0312156110c25760008081fd5b6110ca610db3565b6110d386610d5f565b81526110e0878701610d5f565b8782015260406110f1818801610d5f565b9082015260608681015190820152608080870151600981106111135760008081fd5b90820152835293840193918501916110aa565b50979650505050505050565b6000806040838503121561114557600080fd5b825167ffffffffffffffff8082111561115d57600080fd5b818501915085601f83011261117157600080fd5b81516020611181610e5283610e0d565b82815260059290921b840181019181810190898411156111a057600080fd5b948201945b838610156111ce57855160ff811681146111bf5760008081fd5b825294820194908201906111a5565b918801519196509093505050808211156111e757600080fd5b506111f485828601610e31565b9150509250929050565b81810381811115610a8757610a87610f96565b634e487b7160e01b600052601260045260246000fd5b60008261123657611236611211565b500690565b60008261124a5761124a611211565b500490565b80820180821115610a8757610a87610f9656fea2646970667358221220a6c9591189a194fd90e74aa45842d42e00b630c242b5f57ad514723ce7a8884c64736f6c63430008130033",
}

// MultiCallContractABI is the input ABI used to generate the binding from.
// Deprecated: Use MultiCallContractMetaData.ABI instead.
var MultiCallContractABI = MultiCallContractMetaData.ABI

// MultiCallContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MultiCallContractBinRuntime = `608060405234801561001057600080fd5b50600436106100575760003560e01c80632ada6c5c1461005c5780636082579d14610080578063adde19c614610097578063b04ba218146100af578063bfe8e683146100c8575b600080fd5b6100646100de565b6040516100779796959493929190610b5b565b60405180910390f35b610088610237565b60405161007793929190610be4565b61009f6104f5565b6040516100779493929190610c27565b6100b7610513565b604051610077959493929190610ca6565b6100d06105a5565b604051610077929190610cf6565b60608060008060008060006100f16106ae565b6040805163e70c38f160e01b81528151939a5091985061040092839263e70c38f19260048083019391928290030181865afa158015610134573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101589190610d22565b8096508197505050806001600160a01b03166303e6689d6040518163ffffffff1660e01b81526004016040805180830381865afa15801561019d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101c19190610d22565b8094508195505050806001600160a01b0316633aea9faf6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610207573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061022b9190610d46565b91505090919293949596565b60405163e2693e3f60e01b815260206004820152600a602482015269434c526567697374727960b01b6044820152606090819081906000906104019063e2693e3f90606401602060405180830381865afa158015610299573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102bd9190610d7b565b60405163e2693e3f60e01b815260206004820152600b60248201526a577261707065644b61696160a81b60448201529091506000906104019063e2693e3f90606401602060405180830381865afa15801561031c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103409190610d7b565b90506001600160a01b038216610357575050909192565b816001600160a01b03166390599c076040518163ffffffff1660e01b8152600401600060405180830381865afa158015610395573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103bd9190810190610ea3565b8051929750955081905067ffffffffffffffff8111156103df576103df610d9d565b604051908082528060200260200182016040528015610408578160200160208202803683370190505b5093506001600160a01b038216156104ed578160005b828110156104ea57816001600160a01b03166370a0823188838151811061044757610447610f80565b60200260200101516040518263ffffffff1660e01b815260040161047a91906001600160a01b0391909116815260200190565b602060405180830381865afa158015610497573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104bb9190610d46565b8682815181106104cd576104cd610f80565b6020908102919091010152806104e281610fac565b91505061041e565b50505b505050909192565b606080606060006105046107cf565b93509350935093505b90919293565b60608060008060006105236106ae565b809550819650505060006104009050806001600160a01b03166325cf09436040518163ffffffff1660e01b8152600401606060405180830381865afa158015610570573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105949190610fc5565b979896979196909550909350915050565b60405163e2693e3f60e01b815260206004820152601160248201527023b0b9b632b9b9a9bbb0b82937baba32b960791b60448201526000906060906104019063e2693e3f90606401602060405180830381865afa15801561060a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061062e9190610d7b565b91506001600160a01b038216610642579091565b816001600160a01b031663d3c7c2c76040518163ffffffff1660e01b8152600401600060405180830381865afa158015610680573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526106a89190810190611008565b90509091565b60608060006104009050806001600160a01b0316630b1fe7846040518163ffffffff1660e01b8152600401600060405180830381865afa1580156106f6573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f1916820160405261071e919081019061103d565b80519093508067ffffffffffffffff81111561073c5761073c610d9d565b604051908082528060200260200182016040528015610765578160200160208202803683370190505b50925060005b818110156107c85761079985828151811061078857610788610f80565b6020026020010151602001516109b7565b8482815181106107ab576107ab610f80565b6020908102919091010152806107c081610fac565b91505061076b565b5050509091565b60608060606000806104009050806001600160a01b031663715b208b6040518163ffffffff1660e01b8152600401600060405180830381865afa15801561081a573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526108429190810190611132565b8095508196505050806001600160a01b0316636abd623d6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610888573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108ac9190610d7b565b91506005845110156108be575061050d565b6000600285516108ce91906111fe565b90506108db600382611227565b156108e757505061050d565b6108f260038261123b565b67ffffffffffffffff81111561090a5761090a610d9d565b604051908082528060200260200182016040528015610933578160200160208202803683370190505b50935060005b818110156109ae576109758661095083600161124f565b8151811061096057610960610f80565b60200260200101516001600160a01b03163190565b8561098160038461123b565b8151811061099157610991610f80565b60209081029190910101526109a760038261124f565b9050610939565b50505090919293565b6000816001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109f7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a1b9190610d46565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015610a59573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a7d9190610d46565b610a8791906111fe565b92915050565b60008151808452602080850194508084016000805b84811015610b1457825180516001600160a01b039081168a52858201518116868b0152604080830151909116908a0152606080820151908a01526080908101519060098210610aff57634e487b7160e01b84526021600452602484fd5b89015260a09097019691830191600101610aa2565b50959695505050505050565b600081518084526020808501945080840160005b83811015610b5057815187529582019590820190600101610b34565b509495945050505050565b60e081526000610b6e60e083018a610a8d565b8281036020840152610b80818a610b20565b604084019890985250506060810194909452608084019290925260a083015260c09091015292915050565b600081518084526020808501945080840160005b83811015610b505781516001600160a01b031687529582019590820190600101610bbf565b606081526000610bf76060830186610bab565b8281036020840152610c098186610bab565b90508281036040840152610c1d8185610b20565b9695505050505050565b6080808252855190820181905260009060209060a0840190828901845b82811015610c6357815160ff1684529284019290840190600101610c44565b50505083810382850152610c778188610bab565b9150508281036040840152610c8c8186610b20565b91505060018060a01b038316606083015295945050505050565b60a081526000610cb960a0830188610a8d565b8281036020840152610ccb8188610b20565b6001600160a01b03968716604085015294861660608401525050921660809092019190915292915050565b6001600160a01b0383168152604060208201819052600090610d1a90830184610bab565b949350505050565b60008060408385031215610d3557600080fd5b505080516020909101519092909150565b600060208284031215610d5857600080fd5b5051919050565b80516001600160a01b0381168114610d7657600080fd5b919050565b600060208284031215610d8d57600080fd5b610d9682610d5f565b9392505050565b634e487b7160e01b600052604160045260246000fd5b60405160a0810167ffffffffffffffff81118282101715610dd657610dd6610d9d565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610e0557610e05610d9d565b604052919050565b600067ffffffffffffffff821115610e2757610e27610d9d565b5060051b60200190565b600082601f830112610e4257600080fd5b81516020610e57610e5283610e0d565b610ddc565b82815260059290921b84018101918181019086841115610e7657600080fd5b8286015b84811015610e9857610e8b81610d5f565b8352918301918301610e7a565b509695505050505050565b600080600060608486031215610eb857600080fd5b835167ffffffffffffffff80821115610ed057600080fd5b610edc87838801610e31565b9450602091508186015181811115610ef357600080fd5b8601601f81018813610f0457600080fd5b8051610f12610e5282610e0d565b81815260059190911b8201840190848101908a831115610f3157600080fd5b928501925b82841015610f4f57835182529285019290850190610f36565b60408a0151909750945050505080821115610f6957600080fd5b50610f7686828701610e31565b9150509250925092565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201610fbe57610fbe610f96565b5060010190565b600080600060608486031215610fda57600080fd5b610fe384610d5f565b9250610ff160208501610d5f565b9150610fff60408501610d5f565b90509250925092565b60006020828403121561101a57600080fd5b815167ffffffffffffffff81111561103157600080fd5b610d1a84828501610e31565b6000602080838503121561105057600080fd5b825167ffffffffffffffff81111561106757600080fd5b8301601f8101851361107857600080fd5b8051611086610e5282610e0d565b81815260a091820283018401918482019190888411156110a557600080fd5b938501935b838510156111265780858a0312156110c25760008081fd5b6110ca610db3565b6110d386610d5f565b81526110e0878701610d5f565b8782015260406110f1818801610d5f565b9082015260608681015190820152608080870151600981106111135760008081fd5b90820152835293840193918501916110aa565b50979650505050505050565b6000806040838503121561114557600080fd5b825167ffffffffffffffff8082111561115d57600080fd5b818501915085601f83011261117157600080fd5b81516020611181610e5283610e0d565b82815260059290921b840181019181810190898411156111a057600080fd5b948201945b838610156111ce57855160ff811681146111bf5760008081fd5b825294820194908201906111a5565b918801519196509093505050808211156111e757600080fd5b506111f485828601610e31565b9150509250929050565b81810381811115610a8757610a87610f96565b634e487b7160e01b600052601260045260246000fd5b60008261123657611236611211565b500690565b60008261124a5761124a611211565b500490565b80820180821115610a8757610a87610f9656fea2646970667358221220a6c9591189a194fd90e74aa45842d42e00b630c242b5f57ad514723ce7a8884c64736f6c63430008130033`

// Deprecated: Use MultiCallContractMetaData.Sigs instead.
// MultiCallContractFuncSigs maps the 4-byte function signature to its string representation.
var MultiCallContractFuncSigs = MultiCallContractMetaData.Sigs

// MultiCallContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MultiCallContractMetaData.Bin instead.
var MultiCallContractBin = MultiCallContractMetaData.Bin

// DeployMultiCallContract deploys a new Kaia contract, binding an instance of MultiCallContract to it.
func DeployMultiCallContract(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *MultiCallContract, error) {
	parsed, err := MultiCallContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MultiCallContractBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &MultiCallContract{MultiCallContractCaller: MultiCallContractCaller{contract: contract}, MultiCallContractTransactor: MultiCallContractTransactor{contract: contract}, MultiCallContractFilterer: MultiCallContractFilterer{contract: contract}}, nil
}

// MultiCallContract is an auto generated Go binding around a Kaia contract.
type MultiCallContract struct {
	MultiCallContractCaller     // Read-only binding to the contract
	MultiCallContractTransactor // Write-only binding to the contract
	MultiCallContractFilterer   // Log filterer for contract events
}

// MultiCallContractCaller is an auto generated read-only Go binding around a Kaia contract.
type MultiCallContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiCallContractTransactor is an auto generated write-only Go binding around a Kaia contract.
type MultiCallContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiCallContractFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type MultiCallContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MultiCallContractSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type MultiCallContractSession struct {
	Contract     *MultiCallContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MultiCallContractCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type MultiCallContractCallerSession struct {
	Contract *MultiCallContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MultiCallContractTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type MultiCallContractTransactorSession struct {
	Contract     *MultiCallContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MultiCallContractRaw is an auto generated low-level Go binding around a Kaia contract.
type MultiCallContractRaw struct {
	Contract *MultiCallContract // Generic contract binding to access the raw methods on
}

// MultiCallContractCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type MultiCallContractCallerRaw struct {
	Contract *MultiCallContractCaller // Generic read-only contract binding to access the raw methods on
}

// MultiCallContractTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type MultiCallContractTransactorRaw struct {
	Contract *MultiCallContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMultiCallContract creates a new instance of MultiCallContract, bound to a specific deployed contract.
func NewMultiCallContract(address common.Address, backend bind.ContractBackend) (*MultiCallContract, error) {
	contract, err := bindMultiCallContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MultiCallContract{MultiCallContractCaller: MultiCallContractCaller{contract: contract}, MultiCallContractTransactor: MultiCallContractTransactor{contract: contract}, MultiCallContractFilterer: MultiCallContractFilterer{contract: contract}}, nil
}

// NewMultiCallContractCaller creates a new read-only instance of MultiCallContract, bound to a specific deployed contract.
func NewMultiCallContractCaller(address common.Address, caller bind.ContractCaller) (*MultiCallContractCaller, error) {
	contract, err := bindMultiCallContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MultiCallContractCaller{contract: contract}, nil
}

// NewMultiCallContractTransactor creates a new write-only instance of MultiCallContract, bound to a specific deployed contract.
func NewMultiCallContractTransactor(address common.Address, transactor bind.ContractTransactor) (*MultiCallContractTransactor, error) {
	contract, err := bindMultiCallContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MultiCallContractTransactor{contract: contract}, nil
}

// NewMultiCallContractFilterer creates a new log filterer instance of MultiCallContract, bound to a specific deployed contract.
func NewMultiCallContractFilterer(address common.Address, filterer bind.ContractFilterer) (*MultiCallContractFilterer, error) {
	contract, err := bindMultiCallContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MultiCallContractFilterer{contract: contract}, nil
}

// bindMultiCallContract binds a generic wrapper to an already deployed contract.
func bindMultiCallContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MultiCallContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiCallContract *MultiCallContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiCallContract.Contract.MultiCallContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiCallContract *MultiCallContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiCallContract.Contract.MultiCallContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiCallContract *MultiCallContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiCallContract.Contract.MultiCallContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MultiCallContract *MultiCallContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MultiCallContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MultiCallContract *MultiCallContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MultiCallContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MultiCallContract *MultiCallContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MultiCallContract.Contract.contract.Transact(opts, method, params...)
}

// MultiCallDPStakingInfo is a free data retrieval call binding the contract method 0x6082579d.
//
// Solidity: function multiCallDPStakingInfo() view returns(address[] nodeIds, address[] clPools, uint256[] stakingAmounts)
func (_MultiCallContract *MultiCallContractCaller) MultiCallDPStakingInfo(opts *bind.CallOpts) (struct {
	NodeIds        []common.Address
	ClPools        []common.Address
	StakingAmounts []*big.Int
}, error) {
	var out []interface{}
	err := _MultiCallContract.contract.Call(opts, &out, "multiCallDPStakingInfo")

	outstruct := new(struct {
		NodeIds        []common.Address
		ClPools        []common.Address
		StakingAmounts []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NodeIds = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.ClPools = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.StakingAmounts = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// MultiCallDPStakingInfo is a free data retrieval call binding the contract method 0x6082579d.
//
// Solidity: function multiCallDPStakingInfo() view returns(address[] nodeIds, address[] clPools, uint256[] stakingAmounts)
func (_MultiCallContract *MultiCallContractSession) MultiCallDPStakingInfo() (struct {
	NodeIds        []common.Address
	ClPools        []common.Address
	StakingAmounts []*big.Int
}, error) {
	return _MultiCallContract.Contract.MultiCallDPStakingInfo(&_MultiCallContract.CallOpts)
}

// MultiCallDPStakingInfo is a free data retrieval call binding the contract method 0x6082579d.
//
// Solidity: function multiCallDPStakingInfo() view returns(address[] nodeIds, address[] clPools, uint256[] stakingAmounts)
func (_MultiCallContract *MultiCallContractCallerSession) MultiCallDPStakingInfo() (struct {
	NodeIds        []common.Address
	ClPools        []common.Address
	StakingAmounts []*big.Int
}, error) {
	return _MultiCallContract.Contract.MultiCallDPStakingInfo(&_MultiCallContract.CallOpts)
}

// MultiCallGaslessInfo is a free data retrieval call binding the contract method 0xbfe8e683.
//
// Solidity: function multiCallGaslessInfo() view returns(address gsr, address[] tokens)
func (_MultiCallContract *MultiCallContractCaller) MultiCallGaslessInfo(opts *bind.CallOpts) (struct {
	Gsr    common.Address
	Tokens []common.Address
}, error) {
	var out []interface{}
	err := _MultiCallContract.contract.Call(opts, &out, "multiCallGaslessInfo")

	outstruct := new(struct {
		Gsr    common.Address
		Tokens []common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Gsr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Tokens = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return *outstruct, err

}

// MultiCallGaslessInfo is a free data retrieval call binding the contract method 0xbfe8e683.
//
// Solidity: function multiCallGaslessInfo() view returns(address gsr, address[] tokens)
func (_MultiCallContract *MultiCallContractSession) MultiCallGaslessInfo() (struct {
	Gsr    common.Address
	Tokens []common.Address
}, error) {
	return _MultiCallContract.Contract.MultiCallGaslessInfo(&_MultiCallContract.CallOpts)
}

// MultiCallGaslessInfo is a free data retrieval call binding the contract method 0xbfe8e683.
//
// Solidity: function multiCallGaslessInfo() view returns(address gsr, address[] tokens)
func (_MultiCallContract *MultiCallContractCallerSession) MultiCallGaslessInfo() (struct {
	Gsr    common.Address
	Tokens []common.Address
}, error) {
	return _MultiCallContract.Contract.MultiCallGaslessInfo(&_MultiCallContract.CallOpts)
}

// MultiCallNodeStatesPermissionless is a free data retrieval call binding the contract method 0x2ada6c5c.
//
// Solidity: function multiCallNodeStatesPermissionless() view returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, uint256 pauseTimeout, uint256 idleTimeout, uint256 maxValidatorCount, uint256 maxReadyCandidateCount, uint256 exitThreshold)
func (_MultiCallContract *MultiCallContractCaller) MultiCallNodeStatesPermissionless(opts *bind.CallOpts) (struct {
	Profiles               []Profile
	StakingAmounts         []*big.Int
	PauseTimeout           *big.Int
	IdleTimeout            *big.Int
	MaxValidatorCount      *big.Int
	MaxReadyCandidateCount *big.Int
	ExitThreshold          *big.Int
}, error) {
	var out []interface{}
	err := _MultiCallContract.contract.Call(opts, &out, "multiCallNodeStatesPermissionless")

	outstruct := new(struct {
		Profiles               []Profile
		StakingAmounts         []*big.Int
		PauseTimeout           *big.Int
		IdleTimeout            *big.Int
		MaxValidatorCount      *big.Int
		MaxReadyCandidateCount *big.Int
		ExitThreshold          *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Profiles = *abi.ConvertType(out[0], new([]Profile)).(*[]Profile)
	outstruct.StakingAmounts = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.PauseTimeout = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.IdleTimeout = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.MaxValidatorCount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.MaxReadyCandidateCount = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.ExitThreshold = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// MultiCallNodeStatesPermissionless is a free data retrieval call binding the contract method 0x2ada6c5c.
//
// Solidity: function multiCallNodeStatesPermissionless() view returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, uint256 pauseTimeout, uint256 idleTimeout, uint256 maxValidatorCount, uint256 maxReadyCandidateCount, uint256 exitThreshold)
func (_MultiCallContract *MultiCallContractSession) MultiCallNodeStatesPermissionless() (struct {
	Profiles               []Profile
	StakingAmounts         []*big.Int
	PauseTimeout           *big.Int
	IdleTimeout            *big.Int
	MaxValidatorCount      *big.Int
	MaxReadyCandidateCount *big.Int
	ExitThreshold          *big.Int
}, error) {
	return _MultiCallContract.Contract.MultiCallNodeStatesPermissionless(&_MultiCallContract.CallOpts)
}

// MultiCallNodeStatesPermissionless is a free data retrieval call binding the contract method 0x2ada6c5c.
//
// Solidity: function multiCallNodeStatesPermissionless() view returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, uint256 pauseTimeout, uint256 idleTimeout, uint256 maxValidatorCount, uint256 maxReadyCandidateCount, uint256 exitThreshold)
func (_MultiCallContract *MultiCallContractCallerSession) MultiCallNodeStatesPermissionless() (struct {
	Profiles               []Profile
	StakingAmounts         []*big.Int
	PauseTimeout           *big.Int
	IdleTimeout            *big.Int
	MaxValidatorCount      *big.Int
	MaxReadyCandidateCount *big.Int
	ExitThreshold          *big.Int
}, error) {
	return _MultiCallContract.Contract.MultiCallNodeStatesPermissionless(&_MultiCallContract.CallOpts)
}

// MultiCallStakingInfo is a free data retrieval call binding the contract method 0xadde19c6.
//
// Solidity: function multiCallStakingInfo() view returns(uint8[] typeList, address[] addressList, uint256[] stakingAmounts, address spareAddress)
func (_MultiCallContract *MultiCallContractCaller) MultiCallStakingInfo(opts *bind.CallOpts) (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
	SpareAddress   common.Address
}, error) {
	var out []interface{}
	err := _MultiCallContract.contract.Call(opts, &out, "multiCallStakingInfo")

	outstruct := new(struct {
		TypeList       []uint8
		AddressList    []common.Address
		StakingAmounts []*big.Int
		SpareAddress   common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TypeList = *abi.ConvertType(out[0], new([]uint8)).(*[]uint8)
	outstruct.AddressList = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.StakingAmounts = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)
	outstruct.SpareAddress = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// MultiCallStakingInfo is a free data retrieval call binding the contract method 0xadde19c6.
//
// Solidity: function multiCallStakingInfo() view returns(uint8[] typeList, address[] addressList, uint256[] stakingAmounts, address spareAddress)
func (_MultiCallContract *MultiCallContractSession) MultiCallStakingInfo() (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
	SpareAddress   common.Address
}, error) {
	return _MultiCallContract.Contract.MultiCallStakingInfo(&_MultiCallContract.CallOpts)
}

// MultiCallStakingInfo is a free data retrieval call binding the contract method 0xadde19c6.
//
// Solidity: function multiCallStakingInfo() view returns(uint8[] typeList, address[] addressList, uint256[] stakingAmounts, address spareAddress)
func (_MultiCallContract *MultiCallContractCallerSession) MultiCallStakingInfo() (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
	SpareAddress   common.Address
}, error) {
	return _MultiCallContract.Contract.MultiCallStakingInfo(&_MultiCallContract.CallOpts)
}

// MultiCallStakingInfoPermissionless is a free data retrieval call binding the contract method 0xb04ba218.
//
// Solidity: function multiCallStakingInfoPermissionless() view returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, address kefAddr, address kifAddr, address kpfAddr)
func (_MultiCallContract *MultiCallContractCaller) MultiCallStakingInfoPermissionless(opts *bind.CallOpts) (struct {
	Profiles       []Profile
	StakingAmounts []*big.Int
	KefAddr        common.Address
	KifAddr        common.Address
	KpfAddr        common.Address
}, error) {
	var out []interface{}
	err := _MultiCallContract.contract.Call(opts, &out, "multiCallStakingInfoPermissionless")

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
// Solidity: function multiCallStakingInfoPermissionless() view returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, address kefAddr, address kifAddr, address kpfAddr)
func (_MultiCallContract *MultiCallContractSession) MultiCallStakingInfoPermissionless() (struct {
	Profiles       []Profile
	StakingAmounts []*big.Int
	KefAddr        common.Address
	KifAddr        common.Address
	KpfAddr        common.Address
}, error) {
	return _MultiCallContract.Contract.MultiCallStakingInfoPermissionless(&_MultiCallContract.CallOpts)
}

// MultiCallStakingInfoPermissionless is a free data retrieval call binding the contract method 0xb04ba218.
//
// Solidity: function multiCallStakingInfoPermissionless() view returns((address,address,address,uint256,uint8)[] profiles, uint256[] stakingAmounts, address kefAddr, address kifAddr, address kpfAddr)
func (_MultiCallContract *MultiCallContractCallerSession) MultiCallStakingInfoPermissionless() (struct {
	Profiles       []Profile
	StakingAmounts []*big.Int
	KefAddr        common.Address
	KifAddr        common.Address
	KpfAddr        common.Address
}, error) {
	return _MultiCallContract.Contract.MultiCallStakingInfoPermissionless(&_MultiCallContract.CallOpts)
}
