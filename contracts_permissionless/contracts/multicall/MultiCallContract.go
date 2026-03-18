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
	ABI: "[{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAddress\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isActivated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"ffa1ad74": "VERSION()",
		"715b208b": "getAllAddress()",
		"4a8c1fb4": "isActivated()",
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

// IAddressBookV2MetaData contains all meta data concerning the IAddressBookV2 contract.
var IAddressBookV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllProfiles\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structProfile[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFundAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"0b1fe784": "getAllProfiles()",
		"25cf0943": "getFundAddresses()",
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
	ABI: "[{\"inputs\":[],\"name\":\"multiCallDPStakingInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIds\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"clPools\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallGaslessInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"gsr\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallStakingInfo\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallStakingInfoPermissionless\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structProfile[]\",\"name\":\"profiles\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"address\",\"name\":\"kefAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kifAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kpfAddr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"6082579d": "multiCallDPStakingInfo()",
		"bfe8e683": "multiCallGaslessInfo()",
		"adde19c6": "multiCallStakingInfo()",
		"b04ba218": "multiCallStakingInfoPermissionless()",
	},
	Bin: "0x6080604052348015600e575f80fd5b50610f648061001c5f395ff3fe608060405234801561000f575f80fd5b506004361061004a575f3560e01c80636082579d1461004e578063adde19c61461006e578063b04ba21814610085578063bfe8e6831461009e575b5f80fd5b6100566100b4565b604051610065939291906108a8565b60405180910390f35b61007661035b565b604051610065939291906108ea565b61008d610373565b60405161006595949392919061094f565b6100a66104f3565b604051610065929190610a2e565b60405163e2693e3f60e01b815260206004820152600a602482015269434c526567697374727960b01b6044820152606090819081905f906104019063e2693e3f90606401602060405180830381865afa158015610113573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101379190610a74565b60405163e2693e3f60e01b815260206004820152600b60248201526a577261707065644b61696160a81b60448201529091505f906104019063e2693e3f90606401602060405180830381865afa158015610193573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101b79190610a74565b90506001600160a01b0382166101ce575050909192565b816001600160a01b03166390599c076040518163ffffffff1660e01b81526004015f60405180830381865afa158015610209573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526102309190810190610b98565b8051929750955081905067ffffffffffffffff81111561025257610252610a94565b60405190808252806020026020018201604052801561027b578160200160208202803683370190505b5093506001600160a01b0382161561035357815f5b8281101561035057816001600160a01b03166370a082318883815181106102b9576102b9610c6d565b60200260200101516040518263ffffffff1660e01b81526004016102ec91906001600160a01b0391909116815260200190565b602060405180830381865afa158015610307573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061032b9190610c81565b86828151811061033d5761033d610c6d565b6020908102919091010152600101610290565b50505b505050909192565b60608060606103686105f5565b925092509250909192565b6060805f805f806104009050806001600160a01b0316630b1fe7846040518163ffffffff1660e01b81526004015f60405180830381865afa1580156103ba573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526103e19190810190610c98565b80519096508067ffffffffffffffff8111156103ff576103ff610a94565b604051908082528060200260200182016040528015610428578160200160208202803683370190505b5095505f5b818110156104805761045b88828151811061044a5761044a610c6d565b602002602001015160200151610765565b87828151811061046d5761046d610c6d565b602090810291909101015260010161042d565b50816001600160a01b03166325cf09436040518163ffffffff1660e01b8152600401606060405180830381865afa1580156104bd573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104e19190610d84565b98999798919790965090945092505050565b60405163e2693e3f60e01b815260206004820152601160248201527023b0b9b632b9b9a9bbb0b82937baba32b960791b60448201525f906060906104019063e2693e3f90606401602060405180830381865afa158015610555573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105799190610a74565b91506001600160a01b03821661058d579091565b816001600160a01b031663d3c7c2c76040518163ffffffff1660e01b81526004015f60405180830381865afa1580156105c8573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526105ef9190810190610dc4565b90509091565b60608060605f6104009050806001600160a01b031663715b208b6040518163ffffffff1660e01b81526004015f60405180830381865afa15801561063b573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526106629190810190610df6565b80519195509350600511156106775750909192565b5f600284516106869190610ece565b9050610693600382610ef5565b1561069f575050909192565b6106aa600382610f08565b67ffffffffffffffff8111156106c2576106c2610a94565b6040519080825280602002602001820160405280156106eb578160200160208202803683370190505b5092505f5b818110156103535761072c85610707836001610f1b565b8151811061071757610717610c6d565b60200260200101516001600160a01b03163190565b84610738600384610f08565b8151811061074857610748610c6d565b602090810291909101015261075e600382610f1b565b90506106f0565b5f816001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa1580156107a2573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107c69190610c81565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015610802573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906108269190610c81565b6108309190610ece565b92915050565b5f815180845260208085019450602084015f5b8381101561086e5781516001600160a01b031687529582019590820190600101610849565b509495945050505050565b5f815180845260208085019450602084015f5b8381101561086e5781518752958201959082019060010161088c565b606081525f6108ba6060830186610836565b82810360208401526108cc8186610836565b905082810360408401526108e08185610879565b9695505050505050565b606080825284519082018190525f906020906080840190828801845b8281101561092557815160ff1684529284019290840190600101610906565b505050838103602085015261093a8187610836565b91505082810360408401526108e08185610879565b60a080825286518282018190525f919060209060c0850190828b01855b828110156109dc57815180516001600160a01b03908116865286820151811687870152604080830151909116908601526060808201519086015260809081015190600982106109c957634e487b7160e01b5f52602160045260245ffd5b850152928501929084019060010161096c565b50505084810360208601526109f1818a610879565b9350505050610a0b60408301866001600160a01b03169052565b6001600160a01b03841660608301526001600160a01b03831660808301526108e0565b6001600160a01b03831681526040602082018190525f90610a5190830184610836565b949350505050565b80516001600160a01b0381168114610a6f575f80fd5b919050565b5f60208284031215610a84575f80fd5b610a8d82610a59565b9392505050565b634e487b7160e01b5f52604160045260245ffd5b60405160a0810167ffffffffffffffff81118282101715610acb57610acb610a94565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610afa57610afa610a94565b604052919050565b5f67ffffffffffffffff821115610b1b57610b1b610a94565b5060051b60200190565b5f82601f830112610b34575f80fd5b81516020610b49610b4483610b02565b610ad1565b8083825260208201915060208460051b870101935086841115610b6a575f80fd5b602086015b84811015610b8d57610b8081610a59565b8352918301918301610b6f565b509695505050505050565b5f805f60608486031215610baa575f80fd5b835167ffffffffffffffff80821115610bc1575f80fd5b610bcd87838801610b25565b9450602091508186015181811115610be3575f80fd5b8601601f81018813610bf3575f80fd5b8051610c01610b4482610b02565b81815260059190911b8201840190848101908a831115610c1f575f80fd5b928501925b82841015610c3d57835182529285019290850190610c24565b60408a0151909750945050505080821115610c56575f80fd5b50610c6386828701610b25565b9150509250925092565b634e487b7160e01b5f52603260045260245ffd5b5f60208284031215610c91575f80fd5b5051919050565b5f6020808385031215610ca9575f80fd5b825167ffffffffffffffff811115610cbf575f80fd5b8301601f81018513610ccf575f80fd5b8051610cdd610b4482610b02565b81815260a09182028301840191848201919088841115610cfb575f80fd5b938501935b83851015610d785780858a031215610d16575f80fd5b610d1e610aa8565b610d2786610a59565b8152610d34878701610a59565b878201526040610d45818801610a59565b908201526060868101519082015260808087015160098110610d65575f80fd5b9082015283529384019391850191610d00565b50979650505050505050565b5f805f60608486031215610d96575f80fd5b610d9f84610a59565b9250610dad60208501610a59565b9150610dbb60408501610a59565b90509250925092565b5f60208284031215610dd4575f80fd5b815167ffffffffffffffff811115610dea575f80fd5b610a5184828501610b25565b5f8060408385031215610e07575f80fd5b825167ffffffffffffffff80821115610e1e575f80fd5b818501915085601f830112610e31575f80fd5b81516020610e41610b4483610b02565b82815260059290921b84018101918181019089841115610e5f575f80fd5b948201945b83861015610e8b57855160ff81168114610e7c575f80fd5b82529482019490820190610e64565b91880151919650909350505080821115610ea3575f80fd5b50610eb085828601610b25565b9150509250929050565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561083057610830610eba565b634e487b7160e01b5f52601260045260245ffd5b5f82610f0357610f03610ee1565b500690565b5f82610f1657610f16610ee1565b500490565b8082018082111561083057610830610eba56fea264697066735822122098b8ba88ec7d160ac45ce16b520cae2e4a7df0f275f4526306a9ffcd7020c6b664736f6c63430008190033",
}

// MultiCallContractABI is the input ABI used to generate the binding from.
// Deprecated: Use MultiCallContractMetaData.ABI instead.
var MultiCallContractABI = MultiCallContractMetaData.ABI

// MultiCallContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MultiCallContractBinRuntime = `608060405234801561000f575f80fd5b506004361061004a575f3560e01c80636082579d1461004e578063adde19c61461006e578063b04ba21814610085578063bfe8e6831461009e575b5f80fd5b6100566100b4565b604051610065939291906108a8565b60405180910390f35b61007661035b565b604051610065939291906108ea565b61008d610373565b60405161006595949392919061094f565b6100a66104f3565b604051610065929190610a2e565b60405163e2693e3f60e01b815260206004820152600a602482015269434c526567697374727960b01b6044820152606090819081905f906104019063e2693e3f90606401602060405180830381865afa158015610113573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101379190610a74565b60405163e2693e3f60e01b815260206004820152600b60248201526a577261707065644b61696160a81b60448201529091505f906104019063e2693e3f90606401602060405180830381865afa158015610193573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101b79190610a74565b90506001600160a01b0382166101ce575050909192565b816001600160a01b03166390599c076040518163ffffffff1660e01b81526004015f60405180830381865afa158015610209573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526102309190810190610b98565b8051929750955081905067ffffffffffffffff81111561025257610252610a94565b60405190808252806020026020018201604052801561027b578160200160208202803683370190505b5093506001600160a01b0382161561035357815f5b8281101561035057816001600160a01b03166370a082318883815181106102b9576102b9610c6d565b60200260200101516040518263ffffffff1660e01b81526004016102ec91906001600160a01b0391909116815260200190565b602060405180830381865afa158015610307573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061032b9190610c81565b86828151811061033d5761033d610c6d565b6020908102919091010152600101610290565b50505b505050909192565b60608060606103686105f5565b925092509250909192565b6060805f805f806104009050806001600160a01b0316630b1fe7846040518163ffffffff1660e01b81526004015f60405180830381865afa1580156103ba573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526103e19190810190610c98565b80519096508067ffffffffffffffff8111156103ff576103ff610a94565b604051908082528060200260200182016040528015610428578160200160208202803683370190505b5095505f5b818110156104805761045b88828151811061044a5761044a610c6d565b602002602001015160200151610765565b87828151811061046d5761046d610c6d565b602090810291909101015260010161042d565b50816001600160a01b03166325cf09436040518163ffffffff1660e01b8152600401606060405180830381865afa1580156104bd573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104e19190610d84565b98999798919790965090945092505050565b60405163e2693e3f60e01b815260206004820152601160248201527023b0b9b632b9b9a9bbb0b82937baba32b960791b60448201525f906060906104019063e2693e3f90606401602060405180830381865afa158015610555573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105799190610a74565b91506001600160a01b03821661058d579091565b816001600160a01b031663d3c7c2c76040518163ffffffff1660e01b81526004015f60405180830381865afa1580156105c8573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526105ef9190810190610dc4565b90509091565b60608060605f6104009050806001600160a01b031663715b208b6040518163ffffffff1660e01b81526004015f60405180830381865afa15801561063b573d5f803e3d5ffd5b505050506040513d5f823e601f3d908101601f191682016040526106629190810190610df6565b80519195509350600511156106775750909192565b5f600284516106869190610ece565b9050610693600382610ef5565b1561069f575050909192565b6106aa600382610f08565b67ffffffffffffffff8111156106c2576106c2610a94565b6040519080825280602002602001820160405280156106eb578160200160208202803683370190505b5092505f5b818110156103535761072c85610707836001610f1b565b8151811061071757610717610c6d565b60200260200101516001600160a01b03163190565b84610738600384610f08565b8151811061074857610748610c6d565b602090810291909101015261075e600382610f1b565b90506106f0565b5f816001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa1580156107a2573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107c69190610c81565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa158015610802573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906108269190610c81565b6108309190610ece565b92915050565b5f815180845260208085019450602084015f5b8381101561086e5781516001600160a01b031687529582019590820190600101610849565b509495945050505050565b5f815180845260208085019450602084015f5b8381101561086e5781518752958201959082019060010161088c565b606081525f6108ba6060830186610836565b82810360208401526108cc8186610836565b905082810360408401526108e08185610879565b9695505050505050565b606080825284519082018190525f906020906080840190828801845b8281101561092557815160ff1684529284019290840190600101610906565b505050838103602085015261093a8187610836565b91505082810360408401526108e08185610879565b60a080825286518282018190525f919060209060c0850190828b01855b828110156109dc57815180516001600160a01b03908116865286820151811687870152604080830151909116908601526060808201519086015260809081015190600982106109c957634e487b7160e01b5f52602160045260245ffd5b850152928501929084019060010161096c565b50505084810360208601526109f1818a610879565b9350505050610a0b60408301866001600160a01b03169052565b6001600160a01b03841660608301526001600160a01b03831660808301526108e0565b6001600160a01b03831681526040602082018190525f90610a5190830184610836565b949350505050565b80516001600160a01b0381168114610a6f575f80fd5b919050565b5f60208284031215610a84575f80fd5b610a8d82610a59565b9392505050565b634e487b7160e01b5f52604160045260245ffd5b60405160a0810167ffffffffffffffff81118282101715610acb57610acb610a94565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610afa57610afa610a94565b604052919050565b5f67ffffffffffffffff821115610b1b57610b1b610a94565b5060051b60200190565b5f82601f830112610b34575f80fd5b81516020610b49610b4483610b02565b610ad1565b8083825260208201915060208460051b870101935086841115610b6a575f80fd5b602086015b84811015610b8d57610b8081610a59565b8352918301918301610b6f565b509695505050505050565b5f805f60608486031215610baa575f80fd5b835167ffffffffffffffff80821115610bc1575f80fd5b610bcd87838801610b25565b9450602091508186015181811115610be3575f80fd5b8601601f81018813610bf3575f80fd5b8051610c01610b4482610b02565b81815260059190911b8201840190848101908a831115610c1f575f80fd5b928501925b82841015610c3d57835182529285019290850190610c24565b60408a0151909750945050505080821115610c56575f80fd5b50610c6386828701610b25565b9150509250925092565b634e487b7160e01b5f52603260045260245ffd5b5f60208284031215610c91575f80fd5b5051919050565b5f6020808385031215610ca9575f80fd5b825167ffffffffffffffff811115610cbf575f80fd5b8301601f81018513610ccf575f80fd5b8051610cdd610b4482610b02565b81815260a09182028301840191848201919088841115610cfb575f80fd5b938501935b83851015610d785780858a031215610d16575f80fd5b610d1e610aa8565b610d2786610a59565b8152610d34878701610a59565b878201526040610d45818801610a59565b908201526060868101519082015260808087015160098110610d65575f80fd5b9082015283529384019391850191610d00565b50979650505050505050565b5f805f60608486031215610d96575f80fd5b610d9f84610a59565b9250610dad60208501610a59565b9150610dbb60408501610a59565b90509250925092565b5f60208284031215610dd4575f80fd5b815167ffffffffffffffff811115610dea575f80fd5b610a5184828501610b25565b5f8060408385031215610e07575f80fd5b825167ffffffffffffffff80821115610e1e575f80fd5b818501915085601f830112610e31575f80fd5b81516020610e41610b4483610b02565b82815260059290921b84018101918181019089841115610e5f575f80fd5b948201945b83861015610e8b57855160ff81168114610e7c575f80fd5b82529482019490820190610e64565b91880151919650909350505080821115610ea3575f80fd5b50610eb085828601610b25565b9150509250929050565b634e487b7160e01b5f52601160045260245ffd5b8181038181111561083057610830610eba565b634e487b7160e01b5f52601260045260245ffd5b5f82610f0357610f03610ee1565b500690565b5f82610f1657610f16610ee1565b500490565b8082018082111561083057610830610eba56fea264697066735822122098b8ba88ec7d160ac45ce16b520cae2e4a7df0f275f4526306a9ffcd7020c6b664736f6c63430008190033`

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

// MultiCallStakingInfo is a free data retrieval call binding the contract method 0xadde19c6.
//
// Solidity: function multiCallStakingInfo() view returns(uint8[] typeList, address[] addressList, uint256[] stakingAmounts)
func (_MultiCallContract *MultiCallContractCaller) MultiCallStakingInfo(opts *bind.CallOpts) (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
}, error) {
	var out []interface{}
	err := _MultiCallContract.contract.Call(opts, &out, "multiCallStakingInfo")

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
func (_MultiCallContract *MultiCallContractSession) MultiCallStakingInfo() (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
}, error) {
	return _MultiCallContract.Contract.MultiCallStakingInfo(&_MultiCallContract.CallOpts)
}

// MultiCallStakingInfo is a free data retrieval call binding the contract method 0xadde19c6.
//
// Solidity: function multiCallStakingInfo() view returns(uint8[] typeList, address[] addressList, uint256[] stakingAmounts)
func (_MultiCallContract *MultiCallContractCallerSession) MultiCallStakingInfo() (struct {
	TypeList       []uint8
	AddressList    []common.Address
	StakingAmounts []*big.Int
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
