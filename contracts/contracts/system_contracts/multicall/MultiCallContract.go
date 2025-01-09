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

// IAddressBookMetaData contains all meta data concerning the IAddressBook contract.
var IAddressBookMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllAddress\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isActivated\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
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

// ICLRegistryMetaData contains all meta data concerning the ICLRegistry contract.
var ICLRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllCLs\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// Solidity: function getAllCLs() view returns(address[], uint256[], address[], address[])
func (_ICLRegistry *ICLRegistryCaller) GetAllCLs(opts *bind.CallOpts) ([]common.Address, []*big.Int, []common.Address, []common.Address, error) {
	var out []interface{}
	err := _ICLRegistry.contract.Call(opts, &out, "getAllCLs")

	if err != nil {
		return *new([]common.Address), *new([]*big.Int), *new([]common.Address), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	out2 := *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	out3 := *abi.ConvertType(out[3], new([]common.Address)).(*[]common.Address)

	return out0, out1, out2, out3, err

}

// GetAllCLs is a free data retrieval call binding the contract method 0x90599c07.
//
// Solidity: function getAllCLs() view returns(address[], uint256[], address[], address[])
func (_ICLRegistry *ICLRegistrySession) GetAllCLs() ([]common.Address, []*big.Int, []common.Address, []common.Address, error) {
	return _ICLRegistry.Contract.GetAllCLs(&_ICLRegistry.CallOpts)
}

// GetAllCLs is a free data retrieval call binding the contract method 0x90599c07.
//
// Solidity: function getAllCLs() view returns(address[], uint256[], address[], address[])
func (_ICLRegistry *ICLRegistryCallerSession) GetAllCLs() ([]common.Address, []*big.Int, []common.Address, []common.Address, error) {
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
	ABI: "[{\"inputs\":[],\"name\":\"multiCallDPStakingInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIds\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"clPools\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"clStakings\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiCallStakingInfo\",\"outputs\":[{\"internalType\":\"uint8[]\",\"name\":\"typeList\",\"type\":\"uint8[]\"},{\"internalType\":\"address[]\",\"name\":\"addressList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"stakingAmounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"6082579d": "multiCallDPStakingInfo()",
		"adde19c6": "multiCallStakingInfo()",
	},
	Bin: "0x608060405234801561001057600080fd5b506109d9806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80636082579d1461003b578063adde19c61461005c575b600080fd5b610043610073565b604051610053949392919061052b565b60405180910390f35b610064610338565b60405161005393929190610583565b60405163e2693e3f60e01b815260206004820152600a602482015269434c526567697374727960b01b60448201526060908190819081906000906104019063e2693e3f90606401602060405180830381865afa1580156100d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100fb919061060e565b60405163e2693e3f60e01b815260206004820152600b60248201526a577261707065644b61696160a81b60448201529091506000906104019063e2693e3f90606401602060405180830381865afa15801561015a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061017e919061060e565b90506001600160a01b038216610195575050610332565b816001600160a01b03166390599c076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156101d3573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526101fb919081019061070d565b8151939950909750955081905067ffffffffffffffff81111561022057610220610630565b604051908082528060200260200182016040528015610249578160200160208202803683370190505b5093506001600160a01b0382161561032e578160005b8281101561032b57816001600160a01b03166370a082318983815181106102885761028861080f565b60200260200101516040518263ffffffff1660e01b81526004016102bb91906001600160a01b0391909116815260200190565b602060405180830381865afa1580156102d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102fc9190610825565b86828151811061030e5761030e61080f565b60209081029190910101528061032381610854565b91505061025f565b50505b5050505b90919293565b606080606060006104009050806001600160a01b031663715b208b6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610382573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103aa919081019061086d565b80519195509350600511156103bf5750909192565b6000600284516103cf9190610939565b90506103dc600382610968565b156103e8575050909192565b6103f360038261097c565b67ffffffffffffffff81111561040b5761040b610630565b604051908082528060200260200182016040528015610434578160200160208202803683370190505b50925060005b818110156104af5761047685610451836001610990565b815181106104615761046161080f565b60200260200101516001600160a01b03163190565b8461048260038461097c565b815181106104925761049261080f565b60209081029190910101526104a8600382610990565b905061043a565b505050909192565b600081518084526020808501945080840160005b838110156104f05781516001600160a01b0316875295820195908201906001016104cb565b509495945050505050565b600081518084526020808501945080840160005b838110156104f05781518752958201959082019060010161050f565b60808152600061053e60808301876104b7565b828103602084015261055081876104b7565b9050828103604084015261056481866104b7565b9050828103606084015261057881856104fb565b979650505050505050565b606080825284519082018190526000906020906080840190828801845b828110156105bf57815160ff16845292840192908401906001016105a0565b505050838103828501526105d381876104b7565b91505082810360408401526105e881856104fb565b9695505050505050565b80516001600160a01b038116811461060957600080fd5b919050565b60006020828403121561062057600080fd5b610629826105f2565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561066f5761066f610630565b604052919050565b600067ffffffffffffffff82111561069157610691610630565b5060051b60200190565b600082601f8301126106ac57600080fd5b815160206106c16106bc83610677565b610646565b82815260059290921b840181019181810190868411156106e057600080fd5b8286015b84811015610702576106f5816105f2565b83529183019183016106e4565b509695505050505050565b6000806000806080858703121561072357600080fd5b845167ffffffffffffffff8082111561073b57600080fd5b6107478883890161069b565b955060209150818701518181111561075e57600080fd5b8701601f8101891361076f57600080fd5b805161077d6106bc82610677565b81815260059190911b8201840190848101908b83111561079c57600080fd5b928501925b828410156107ba578351825292850192908501906107a1565b60408b01519098509450505050808211156107d457600080fd5b6107e08883890161069b565b935060608701519150808211156107f657600080fd5b506108038782880161069b565b91505092959194509250565b634e487b7160e01b600052603260045260246000fd5b60006020828403121561083757600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b6000600182016108665761086661083e565b5060010190565b6000806040838503121561088057600080fd5b825167ffffffffffffffff8082111561089857600080fd5b818501915085601f8301126108ac57600080fd5b815160206108bc6106bc83610677565b82815260059290921b840181019181810190898411156108db57600080fd5b948201945b8386101561090957855160ff811681146108fa5760008081fd5b825294820194908201906108e0565b9188015191965090935050508082111561092257600080fd5b5061092f8582860161069b565b9150509250929050565b8181038181111561094c5761094c61083e565b92915050565b634e487b7160e01b600052601260045260246000fd5b60008261097757610977610952565b500690565b60008261098b5761098b610952565b500490565b8082018082111561094c5761094c61083e56fea2646970667358221220f2777d2f144f3b75748af5d8c1431536882b8844df9f7077c25c6a7720f1e49a64736f6c63430008130033",
}

// MultiCallContractABI is the input ABI used to generate the binding from.
// Deprecated: Use MultiCallContractMetaData.ABI instead.
var MultiCallContractABI = MultiCallContractMetaData.ABI

// MultiCallContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MultiCallContractBinRuntime = `608060405234801561001057600080fd5b50600436106100365760003560e01c80636082579d1461003b578063adde19c61461005c575b600080fd5b610043610073565b604051610053949392919061052b565b60405180910390f35b610064610338565b60405161005393929190610583565b60405163e2693e3f60e01b815260206004820152600a602482015269434c526567697374727960b01b60448201526060908190819081906000906104019063e2693e3f90606401602060405180830381865afa1580156100d7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906100fb919061060e565b60405163e2693e3f60e01b815260206004820152600b60248201526a577261707065644b61696160a81b60448201529091506000906104019063e2693e3f90606401602060405180830381865afa15801561015a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061017e919061060e565b90506001600160a01b038216610195575050610332565b816001600160a01b03166390599c076040518163ffffffff1660e01b8152600401600060405180830381865afa1580156101d3573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526101fb919081019061070d565b8151939950909750955081905067ffffffffffffffff81111561022057610220610630565b604051908082528060200260200182016040528015610249578160200160208202803683370190505b5093506001600160a01b0382161561032e578160005b8281101561032b57816001600160a01b03166370a082318983815181106102885761028861080f565b60200260200101516040518263ffffffff1660e01b81526004016102bb91906001600160a01b0391909116815260200190565b602060405180830381865afa1580156102d8573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102fc9190610825565b86828151811061030e5761030e61080f565b60209081029190910101528061032381610854565b91505061025f565b50505b5050505b90919293565b606080606060006104009050806001600160a01b031663715b208b6040518163ffffffff1660e01b8152600401600060405180830381865afa158015610382573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103aa919081019061086d565b80519195509350600511156103bf5750909192565b6000600284516103cf9190610939565b90506103dc600382610968565b156103e8575050909192565b6103f360038261097c565b67ffffffffffffffff81111561040b5761040b610630565b604051908082528060200260200182016040528015610434578160200160208202803683370190505b50925060005b818110156104af5761047685610451836001610990565b815181106104615761046161080f565b60200260200101516001600160a01b03163190565b8461048260038461097c565b815181106104925761049261080f565b60209081029190910101526104a8600382610990565b905061043a565b505050909192565b600081518084526020808501945080840160005b838110156104f05781516001600160a01b0316875295820195908201906001016104cb565b509495945050505050565b600081518084526020808501945080840160005b838110156104f05781518752958201959082019060010161050f565b60808152600061053e60808301876104b7565b828103602084015261055081876104b7565b9050828103604084015261056481866104b7565b9050828103606084015261057881856104fb565b979650505050505050565b606080825284519082018190526000906020906080840190828801845b828110156105bf57815160ff16845292840192908401906001016105a0565b505050838103828501526105d381876104b7565b91505082810360408401526105e881856104fb565b9695505050505050565b80516001600160a01b038116811461060957600080fd5b919050565b60006020828403121561062057600080fd5b610629826105f2565b9392505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561066f5761066f610630565b604052919050565b600067ffffffffffffffff82111561069157610691610630565b5060051b60200190565b600082601f8301126106ac57600080fd5b815160206106c16106bc83610677565b610646565b82815260059290921b840181019181810190868411156106e057600080fd5b8286015b84811015610702576106f5816105f2565b83529183019183016106e4565b509695505050505050565b6000806000806080858703121561072357600080fd5b845167ffffffffffffffff8082111561073b57600080fd5b6107478883890161069b565b955060209150818701518181111561075e57600080fd5b8701601f8101891361076f57600080fd5b805161077d6106bc82610677565b81815260059190911b8201840190848101908b83111561079c57600080fd5b928501925b828410156107ba578351825292850192908501906107a1565b60408b01519098509450505050808211156107d457600080fd5b6107e08883890161069b565b935060608701519150808211156107f657600080fd5b506108038782880161069b565b91505092959194509250565b634e487b7160e01b600052603260045260246000fd5b60006020828403121561083757600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b6000600182016108665761086661083e565b5060010190565b6000806040838503121561088057600080fd5b825167ffffffffffffffff8082111561089857600080fd5b818501915085601f8301126108ac57600080fd5b815160206108bc6106bc83610677565b82815260059290921b840181019181810190898411156108db57600080fd5b948201945b8386101561090957855160ff811681146108fa5760008081fd5b825294820194908201906108e0565b9188015191965090935050508082111561092257600080fd5b5061092f8582860161069b565b9150509250929050565b8181038181111561094c5761094c61083e565b92915050565b634e487b7160e01b600052601260045260246000fd5b60008261097757610977610952565b500690565b60008261098b5761098b610952565b500490565b8082018082111561094c5761094c61083e56fea2646970667358221220f2777d2f144f3b75748af5d8c1431536882b8844df9f7077c25c6a7720f1e49a64736f6c63430008130033`

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
// Solidity: function multiCallDPStakingInfo() view returns(address[] nodeIds, address[] clPools, address[] clStakings, uint256[] stakingAmounts)
func (_MultiCallContract *MultiCallContractCaller) MultiCallDPStakingInfo(opts *bind.CallOpts) (struct {
	NodeIds        []common.Address
	ClPools        []common.Address
	ClStakings     []common.Address
	StakingAmounts []*big.Int
}, error) {
	var out []interface{}
	err := _MultiCallContract.contract.Call(opts, &out, "multiCallDPStakingInfo")

	outstruct := new(struct {
		NodeIds        []common.Address
		ClPools        []common.Address
		ClStakings     []common.Address
		StakingAmounts []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NodeIds = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.ClPools = *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)
	outstruct.ClStakings = *abi.ConvertType(out[2], new([]common.Address)).(*[]common.Address)
	outstruct.StakingAmounts = *abi.ConvertType(out[3], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// MultiCallDPStakingInfo is a free data retrieval call binding the contract method 0x6082579d.
//
// Solidity: function multiCallDPStakingInfo() view returns(address[] nodeIds, address[] clPools, address[] clStakings, uint256[] stakingAmounts)
func (_MultiCallContract *MultiCallContractSession) MultiCallDPStakingInfo() (struct {
	NodeIds        []common.Address
	ClPools        []common.Address
	ClStakings     []common.Address
	StakingAmounts []*big.Int
}, error) {
	return _MultiCallContract.Contract.MultiCallDPStakingInfo(&_MultiCallContract.CallOpts)
}

// MultiCallDPStakingInfo is a free data retrieval call binding the contract method 0x6082579d.
//
// Solidity: function multiCallDPStakingInfo() view returns(address[] nodeIds, address[] clPools, address[] clStakings, uint256[] stakingAmounts)
func (_MultiCallContract *MultiCallContractCallerSession) MultiCallDPStakingInfo() (struct {
	NodeIds        []common.Address
	ClPools        []common.Address
	ClStakings     []common.Address
	StakingAmounts []*big.Int
}, error) {
	return _MultiCallContract.Contract.MultiCallDPStakingInfo(&_MultiCallContract.CallOpts)
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
