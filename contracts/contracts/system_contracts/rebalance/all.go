// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rebalance

import (
	"errors"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = klaytn.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IRetiredContractMetaData contains all meta data concerning the IRetiredContract contract.
var IRetiredContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"adminList\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"quorom\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1865c57d": "getState()",
	},
}

// IRetiredContractABI is the input ABI used to generate the binding from.
// Deprecated: Use IRetiredContractMetaData.ABI instead.
var IRetiredContractABI = IRetiredContractMetaData.ABI

// IRetiredContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IRetiredContractBinRuntime = ``

// IRetiredContractFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IRetiredContractMetaData.Sigs instead.
var IRetiredContractFuncSigs = IRetiredContractMetaData.Sigs

// IRetiredContract is an auto generated Go binding around a Klaytn contract.
type IRetiredContract struct {
	IRetiredContractCaller     // Read-only binding to the contract
	IRetiredContractTransactor // Write-only binding to the contract
	IRetiredContractFilterer   // Log filterer for contract events
}

// IRetiredContractCaller is an auto generated read-only Go binding around a Klaytn contract.
type IRetiredContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRetiredContractTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IRetiredContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRetiredContractFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IRetiredContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRetiredContractSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IRetiredContractSession struct {
	Contract     *IRetiredContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IRetiredContractCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IRetiredContractCallerSession struct {
	Contract *IRetiredContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IRetiredContractTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IRetiredContractTransactorSession struct {
	Contract     *IRetiredContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IRetiredContractRaw is an auto generated low-level Go binding around a Klaytn contract.
type IRetiredContractRaw struct {
	Contract *IRetiredContract // Generic contract binding to access the raw methods on
}

// IRetiredContractCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IRetiredContractCallerRaw struct {
	Contract *IRetiredContractCaller // Generic read-only contract binding to access the raw methods on
}

// IRetiredContractTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IRetiredContractTransactorRaw struct {
	Contract *IRetiredContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIRetiredContract creates a new instance of IRetiredContract, bound to a specific deployed contract.
func NewIRetiredContract(address common.Address, backend bind.ContractBackend) (*IRetiredContract, error) {
	contract, err := bindIRetiredContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IRetiredContract{IRetiredContractCaller: IRetiredContractCaller{contract: contract}, IRetiredContractTransactor: IRetiredContractTransactor{contract: contract}, IRetiredContractFilterer: IRetiredContractFilterer{contract: contract}}, nil
}

// NewIRetiredContractCaller creates a new read-only instance of IRetiredContract, bound to a specific deployed contract.
func NewIRetiredContractCaller(address common.Address, caller bind.ContractCaller) (*IRetiredContractCaller, error) {
	contract, err := bindIRetiredContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IRetiredContractCaller{contract: contract}, nil
}

// NewIRetiredContractTransactor creates a new write-only instance of IRetiredContract, bound to a specific deployed contract.
func NewIRetiredContractTransactor(address common.Address, transactor bind.ContractTransactor) (*IRetiredContractTransactor, error) {
	contract, err := bindIRetiredContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IRetiredContractTransactor{contract: contract}, nil
}

// NewIRetiredContractFilterer creates a new log filterer instance of IRetiredContract, bound to a specific deployed contract.
func NewIRetiredContractFilterer(address common.Address, filterer bind.ContractFilterer) (*IRetiredContractFilterer, error) {
	contract, err := bindIRetiredContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IRetiredContractFilterer{contract: contract}, nil
}

// bindIRetiredContract binds a generic wrapper to an already deployed contract.
func bindIRetiredContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IRetiredContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRetiredContract *IRetiredContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IRetiredContract.Contract.IRetiredContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRetiredContract *IRetiredContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRetiredContract.Contract.IRetiredContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRetiredContract *IRetiredContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRetiredContract.Contract.IRetiredContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRetiredContract *IRetiredContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IRetiredContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRetiredContract *IRetiredContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRetiredContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRetiredContract *IRetiredContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRetiredContract.Contract.contract.Transact(opts, method, params...)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 quorom)
func (_IRetiredContract *IRetiredContractCaller) GetState(opts *bind.CallOpts) (struct {
	AdminList []common.Address
	Quorom    *big.Int
}, error,
) {
	var out []interface{}
	err := _IRetiredContract.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		AdminList []common.Address
		Quorom    *big.Int
	})

	outstruct.AdminList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Quorom = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, err
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 quorom)
func (_IRetiredContract *IRetiredContractSession) GetState() (struct {
	AdminList []common.Address
	Quorom    *big.Int
}, error,
) {
	return _IRetiredContract.Contract.GetState(&_IRetiredContract.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 quorom)
func (_IRetiredContract *IRetiredContractCallerSession) GetState() (struct {
	AdminList []common.Address
	Quorom    *big.Int
}, error,
) {
	return _IRetiredContract.Contract.GetState(&_IRetiredContract.CallOpts)
}

// ITreasuryRebalanceMetaData contains all meta data concerning the ITreasuryRebalance contract.
var ITreasuryRebalanceMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"retired\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"approversCount\",\"type\":\"uint256\"}],\"name\":\"Approved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumITreasuryRebalance.Status\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rebalanceBlockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deployedBlockNumber\",\"type\":\"uint256\"}],\"name\":\"ContractDeployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"enumITreasuryRebalance.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"Finalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newbie\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundAllocation\",\"type\":\"uint256\"}],\"name\":\"NewbieRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newbie\",\"type\":\"address\"}],\"name\":\"NewbieRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"retired\",\"type\":\"address\"}],\"name\":\"RetiredRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"retired\",\"type\":\"address\"}],\"name\":\"RetiredRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumITreasuryRebalance.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"StatusChanged\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"retiredAddress\",\"type\":\"address\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkRetiredsApproved\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"}],\"name\":\"finalizeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeRegistration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newbieAddress\",\"type\":\"address\"}],\"name\":\"getNewbie\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewbieCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"retiredAddress\",\"type\":\"address\"}],\"name\":\"getRetired\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiredCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"treasuryAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"memo\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rebalanceBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newbieAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"registerNewbie\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"retiredAddress\",\"type\":\"address\"}],\"name\":\"registerRetired\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newbieAddress\",\"type\":\"address\"}],\"name\":\"removeNewbie\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"retiredAddress\",\"type\":\"address\"}],\"name\":\"removeRetired\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"status\",\"outputs\":[{\"internalType\":\"enumITreasuryRebalance.Status\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sumOfRetiredBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"retireesBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"daea85c5": "approve(address)",
		"966e0794": "checkRetiredsApproved()",
		"faaf9ca6": "finalizeApproval()",
		"ea6d4a9b": "finalizeContract(string)",
		"48409096": "finalizeRegistration()",
		"eb5a8e55": "getNewbie(address)",
		"91734d86": "getNewbieCount()",
		"bf680590": "getRetired(address)",
		"d1ed33fc": "getRetiredCount()",
		"e20fcf00": "getTreasuryAmount()",
		"58c3b870": "memo()",
		"49a3fb45": "rebalanceBlockNumber()",
		"652e27e0": "registerNewbie(address,uint256)",
		"1f8c1798": "registerRetired(address)",
		"6864b95b": "removeNewbie(address)",
		"1c1dac59": "removeRetired(address)",
		"d826f88f": "reset()",
		"200d2ed2": "status()",
		"45205a6b": "sumOfRetiredBalance()",
	},
}

// ITreasuryRebalanceABI is the input ABI used to generate the binding from.
// Deprecated: Use ITreasuryRebalanceMetaData.ABI instead.
var ITreasuryRebalanceABI = ITreasuryRebalanceMetaData.ABI

// ITreasuryRebalanceBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ITreasuryRebalanceBinRuntime = ``

// ITreasuryRebalanceFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use ITreasuryRebalanceMetaData.Sigs instead.
var ITreasuryRebalanceFuncSigs = ITreasuryRebalanceMetaData.Sigs

// ITreasuryRebalance is an auto generated Go binding around a Klaytn contract.
type ITreasuryRebalance struct {
	ITreasuryRebalanceCaller     // Read-only binding to the contract
	ITreasuryRebalanceTransactor // Write-only binding to the contract
	ITreasuryRebalanceFilterer   // Log filterer for contract events
}

// ITreasuryRebalanceCaller is an auto generated read-only Go binding around a Klaytn contract.
type ITreasuryRebalanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITreasuryRebalanceTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ITreasuryRebalanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITreasuryRebalanceFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ITreasuryRebalanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ITreasuryRebalanceSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ITreasuryRebalanceSession struct {
	Contract     *ITreasuryRebalance // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ITreasuryRebalanceCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ITreasuryRebalanceCallerSession struct {
	Contract *ITreasuryRebalanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ITreasuryRebalanceTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ITreasuryRebalanceTransactorSession struct {
	Contract     *ITreasuryRebalanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ITreasuryRebalanceRaw is an auto generated low-level Go binding around a Klaytn contract.
type ITreasuryRebalanceRaw struct {
	Contract *ITreasuryRebalance // Generic contract binding to access the raw methods on
}

// ITreasuryRebalanceCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ITreasuryRebalanceCallerRaw struct {
	Contract *ITreasuryRebalanceCaller // Generic read-only contract binding to access the raw methods on
}

// ITreasuryRebalanceTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ITreasuryRebalanceTransactorRaw struct {
	Contract *ITreasuryRebalanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewITreasuryRebalance creates a new instance of ITreasuryRebalance, bound to a specific deployed contract.
func NewITreasuryRebalance(address common.Address, backend bind.ContractBackend) (*ITreasuryRebalance, error) {
	contract, err := bindITreasuryRebalance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalance{ITreasuryRebalanceCaller: ITreasuryRebalanceCaller{contract: contract}, ITreasuryRebalanceTransactor: ITreasuryRebalanceTransactor{contract: contract}, ITreasuryRebalanceFilterer: ITreasuryRebalanceFilterer{contract: contract}}, nil
}

// NewITreasuryRebalanceCaller creates a new read-only instance of ITreasuryRebalance, bound to a specific deployed contract.
func NewITreasuryRebalanceCaller(address common.Address, caller bind.ContractCaller) (*ITreasuryRebalanceCaller, error) {
	contract, err := bindITreasuryRebalance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceCaller{contract: contract}, nil
}

// NewITreasuryRebalanceTransactor creates a new write-only instance of ITreasuryRebalance, bound to a specific deployed contract.
func NewITreasuryRebalanceTransactor(address common.Address, transactor bind.ContractTransactor) (*ITreasuryRebalanceTransactor, error) {
	contract, err := bindITreasuryRebalance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceTransactor{contract: contract}, nil
}

// NewITreasuryRebalanceFilterer creates a new log filterer instance of ITreasuryRebalance, bound to a specific deployed contract.
func NewITreasuryRebalanceFilterer(address common.Address, filterer bind.ContractFilterer) (*ITreasuryRebalanceFilterer, error) {
	contract, err := bindITreasuryRebalance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceFilterer{contract: contract}, nil
}

// bindITreasuryRebalance binds a generic wrapper to an already deployed contract.
func bindITreasuryRebalance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ITreasuryRebalanceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ITreasuryRebalance *ITreasuryRebalanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ITreasuryRebalance.Contract.ITreasuryRebalanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ITreasuryRebalance *ITreasuryRebalanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.ITreasuryRebalanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ITreasuryRebalance *ITreasuryRebalanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.ITreasuryRebalanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ITreasuryRebalance *ITreasuryRebalanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ITreasuryRebalance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.contract.Transact(opts, method, params...)
}

// CheckRetiredsApproved is a free data retrieval call binding the contract method 0x966e0794.
//
// Solidity: function checkRetiredsApproved() view returns()
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) CheckRetiredsApproved(opts *bind.CallOpts) error {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "checkRetiredsApproved")
	if err != nil {
		return err
	}

	return err
}

// CheckRetiredsApproved is a free data retrieval call binding the contract method 0x966e0794.
//
// Solidity: function checkRetiredsApproved() view returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) CheckRetiredsApproved() error {
	return _ITreasuryRebalance.Contract.CheckRetiredsApproved(&_ITreasuryRebalance.CallOpts)
}

// CheckRetiredsApproved is a free data retrieval call binding the contract method 0x966e0794.
//
// Solidity: function checkRetiredsApproved() view returns()
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) CheckRetiredsApproved() error {
	return _ITreasuryRebalance.Contract.CheckRetiredsApproved(&_ITreasuryRebalance.CallOpts)
}

// GetNewbie is a free data retrieval call binding the contract method 0xeb5a8e55.
//
// Solidity: function getNewbie(address newbieAddress) view returns(address, uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) GetNewbie(opts *bind.CallOpts, newbieAddress common.Address) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "getNewbie", newbieAddress)
	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err
}

// GetNewbie is a free data retrieval call binding the contract method 0xeb5a8e55.
//
// Solidity: function getNewbie(address newbieAddress) view returns(address, uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceSession) GetNewbie(newbieAddress common.Address) (common.Address, *big.Int, error) {
	return _ITreasuryRebalance.Contract.GetNewbie(&_ITreasuryRebalance.CallOpts, newbieAddress)
}

// GetNewbie is a free data retrieval call binding the contract method 0xeb5a8e55.
//
// Solidity: function getNewbie(address newbieAddress) view returns(address, uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) GetNewbie(newbieAddress common.Address) (common.Address, *big.Int, error) {
	return _ITreasuryRebalance.Contract.GetNewbie(&_ITreasuryRebalance.CallOpts, newbieAddress)
}

// GetNewbieCount is a free data retrieval call binding the contract method 0x91734d86.
//
// Solidity: function getNewbieCount() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) GetNewbieCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "getNewbieCount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetNewbieCount is a free data retrieval call binding the contract method 0x91734d86.
//
// Solidity: function getNewbieCount() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceSession) GetNewbieCount() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.GetNewbieCount(&_ITreasuryRebalance.CallOpts)
}

// GetNewbieCount is a free data retrieval call binding the contract method 0x91734d86.
//
// Solidity: function getNewbieCount() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) GetNewbieCount() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.GetNewbieCount(&_ITreasuryRebalance.CallOpts)
}

// GetRetired is a free data retrieval call binding the contract method 0xbf680590.
//
// Solidity: function getRetired(address retiredAddress) view returns(address, address[])
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) GetRetired(opts *bind.CallOpts, retiredAddress common.Address) (common.Address, []common.Address, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "getRetired", retiredAddress)
	if err != nil {
		return *new(common.Address), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err
}

// GetRetired is a free data retrieval call binding the contract method 0xbf680590.
//
// Solidity: function getRetired(address retiredAddress) view returns(address, address[])
func (_ITreasuryRebalance *ITreasuryRebalanceSession) GetRetired(retiredAddress common.Address) (common.Address, []common.Address, error) {
	return _ITreasuryRebalance.Contract.GetRetired(&_ITreasuryRebalance.CallOpts, retiredAddress)
}

// GetRetired is a free data retrieval call binding the contract method 0xbf680590.
//
// Solidity: function getRetired(address retiredAddress) view returns(address, address[])
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) GetRetired(retiredAddress common.Address) (common.Address, []common.Address, error) {
	return _ITreasuryRebalance.Contract.GetRetired(&_ITreasuryRebalance.CallOpts, retiredAddress)
}

// GetRetiredCount is a free data retrieval call binding the contract method 0xd1ed33fc.
//
// Solidity: function getRetiredCount() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) GetRetiredCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "getRetiredCount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetRetiredCount is a free data retrieval call binding the contract method 0xd1ed33fc.
//
// Solidity: function getRetiredCount() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceSession) GetRetiredCount() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.GetRetiredCount(&_ITreasuryRebalance.CallOpts)
}

// GetRetiredCount is a free data retrieval call binding the contract method 0xd1ed33fc.
//
// Solidity: function getRetiredCount() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) GetRetiredCount() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.GetRetiredCount(&_ITreasuryRebalance.CallOpts)
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) GetTreasuryAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "getTreasuryAmount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_ITreasuryRebalance *ITreasuryRebalanceSession) GetTreasuryAmount() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.GetTreasuryAmount(&_ITreasuryRebalance.CallOpts)
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) GetTreasuryAmount() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.GetTreasuryAmount(&_ITreasuryRebalance.CallOpts)
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) Memo(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "memo")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_ITreasuryRebalance *ITreasuryRebalanceSession) Memo() (string, error) {
	return _ITreasuryRebalance.Contract.Memo(&_ITreasuryRebalance.CallOpts)
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) Memo() (string, error) {
	return _ITreasuryRebalance.Contract.Memo(&_ITreasuryRebalance.CallOpts)
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) RebalanceBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "rebalanceBlockNumber")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceSession) RebalanceBlockNumber() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.RebalanceBlockNumber(&_ITreasuryRebalance.CallOpts)
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) RebalanceBlockNumber() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.RebalanceBlockNumber(&_ITreasuryRebalance.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) Status(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "status")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_ITreasuryRebalance *ITreasuryRebalanceSession) Status() (uint8, error) {
	return _ITreasuryRebalance.Contract.Status(&_ITreasuryRebalance.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) Status() (uint8, error) {
	return _ITreasuryRebalance.Contract.Status(&_ITreasuryRebalance.CallOpts)
}

// SumOfRetiredBalance is a free data retrieval call binding the contract method 0x45205a6b.
//
// Solidity: function sumOfRetiredBalance() view returns(uint256 retireesBalance)
func (_ITreasuryRebalance *ITreasuryRebalanceCaller) SumOfRetiredBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ITreasuryRebalance.contract.Call(opts, &out, "sumOfRetiredBalance")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// SumOfRetiredBalance is a free data retrieval call binding the contract method 0x45205a6b.
//
// Solidity: function sumOfRetiredBalance() view returns(uint256 retireesBalance)
func (_ITreasuryRebalance *ITreasuryRebalanceSession) SumOfRetiredBalance() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.SumOfRetiredBalance(&_ITreasuryRebalance.CallOpts)
}

// SumOfRetiredBalance is a free data retrieval call binding the contract method 0x45205a6b.
//
// Solidity: function sumOfRetiredBalance() view returns(uint256 retireesBalance)
func (_ITreasuryRebalance *ITreasuryRebalanceCallerSession) SumOfRetiredBalance() (*big.Int, error) {
	return _ITreasuryRebalance.Contract.SumOfRetiredBalance(&_ITreasuryRebalance.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) Approve(opts *bind.TransactOpts, retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "approve", retiredAddress)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) Approve(retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.Approve(&_ITreasuryRebalance.TransactOpts, retiredAddress)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) Approve(retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.Approve(&_ITreasuryRebalance.TransactOpts, retiredAddress)
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) FinalizeApproval(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "finalizeApproval")
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) FinalizeApproval() (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.FinalizeApproval(&_ITreasuryRebalance.TransactOpts)
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) FinalizeApproval() (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.FinalizeApproval(&_ITreasuryRebalance.TransactOpts)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string memo) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) FinalizeContract(opts *bind.TransactOpts, memo string) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "finalizeContract", memo)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string memo) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) FinalizeContract(memo string) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.FinalizeContract(&_ITreasuryRebalance.TransactOpts, memo)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string memo) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) FinalizeContract(memo string) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.FinalizeContract(&_ITreasuryRebalance.TransactOpts, memo)
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) FinalizeRegistration(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "finalizeRegistration")
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) FinalizeRegistration() (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.FinalizeRegistration(&_ITreasuryRebalance.TransactOpts)
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) FinalizeRegistration() (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.FinalizeRegistration(&_ITreasuryRebalance.TransactOpts)
}

// RegisterNewbie is a paid mutator transaction binding the contract method 0x652e27e0.
//
// Solidity: function registerNewbie(address newbieAddress, uint256 amount) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) RegisterNewbie(opts *bind.TransactOpts, newbieAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "registerNewbie", newbieAddress, amount)
}

// RegisterNewbie is a paid mutator transaction binding the contract method 0x652e27e0.
//
// Solidity: function registerNewbie(address newbieAddress, uint256 amount) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) RegisterNewbie(newbieAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.RegisterNewbie(&_ITreasuryRebalance.TransactOpts, newbieAddress, amount)
}

// RegisterNewbie is a paid mutator transaction binding the contract method 0x652e27e0.
//
// Solidity: function registerNewbie(address newbieAddress, uint256 amount) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) RegisterNewbie(newbieAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.RegisterNewbie(&_ITreasuryRebalance.TransactOpts, newbieAddress, amount)
}

// RegisterRetired is a paid mutator transaction binding the contract method 0x1f8c1798.
//
// Solidity: function registerRetired(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) RegisterRetired(opts *bind.TransactOpts, retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "registerRetired", retiredAddress)
}

// RegisterRetired is a paid mutator transaction binding the contract method 0x1f8c1798.
//
// Solidity: function registerRetired(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) RegisterRetired(retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.RegisterRetired(&_ITreasuryRebalance.TransactOpts, retiredAddress)
}

// RegisterRetired is a paid mutator transaction binding the contract method 0x1f8c1798.
//
// Solidity: function registerRetired(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) RegisterRetired(retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.RegisterRetired(&_ITreasuryRebalance.TransactOpts, retiredAddress)
}

// RemoveNewbie is a paid mutator transaction binding the contract method 0x6864b95b.
//
// Solidity: function removeNewbie(address newbieAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) RemoveNewbie(opts *bind.TransactOpts, newbieAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "removeNewbie", newbieAddress)
}

// RemoveNewbie is a paid mutator transaction binding the contract method 0x6864b95b.
//
// Solidity: function removeNewbie(address newbieAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) RemoveNewbie(newbieAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.RemoveNewbie(&_ITreasuryRebalance.TransactOpts, newbieAddress)
}

// RemoveNewbie is a paid mutator transaction binding the contract method 0x6864b95b.
//
// Solidity: function removeNewbie(address newbieAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) RemoveNewbie(newbieAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.RemoveNewbie(&_ITreasuryRebalance.TransactOpts, newbieAddress)
}

// RemoveRetired is a paid mutator transaction binding the contract method 0x1c1dac59.
//
// Solidity: function removeRetired(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) RemoveRetired(opts *bind.TransactOpts, retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "removeRetired", retiredAddress)
}

// RemoveRetired is a paid mutator transaction binding the contract method 0x1c1dac59.
//
// Solidity: function removeRetired(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) RemoveRetired(retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.RemoveRetired(&_ITreasuryRebalance.TransactOpts, retiredAddress)
}

// RemoveRetired is a paid mutator transaction binding the contract method 0x1c1dac59.
//
// Solidity: function removeRetired(address retiredAddress) returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) RemoveRetired(retiredAddress common.Address) (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.RemoveRetired(&_ITreasuryRebalance.TransactOpts, retiredAddress)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ITreasuryRebalance.contract.Transact(opts, "reset")
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceSession) Reset() (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.Reset(&_ITreasuryRebalance.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_ITreasuryRebalance *ITreasuryRebalanceTransactorSession) Reset() (*types.Transaction, error) {
	return _ITreasuryRebalance.Contract.Reset(&_ITreasuryRebalance.TransactOpts)
}

// ITreasuryRebalanceApprovedIterator is returned from FilterApproved and is used to iterate over the raw logs and unpacked data for Approved events raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceApprovedIterator struct {
	Event *ITreasuryRebalanceApproved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITreasuryRebalanceApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITreasuryRebalanceApproved)
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
		it.Event = new(ITreasuryRebalanceApproved)
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
func (it *ITreasuryRebalanceApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITreasuryRebalanceApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITreasuryRebalanceApproved represents a Approved event raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceApproved struct {
	Retired        common.Address
	Approver       common.Address
	ApproversCount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterApproved is a free log retrieval operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address retired, address approver, uint256 approversCount)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) FilterApproved(opts *bind.FilterOpts) (*ITreasuryRebalanceApprovedIterator, error) {
	logs, sub, err := _ITreasuryRebalance.contract.FilterLogs(opts, "Approved")
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceApprovedIterator{contract: _ITreasuryRebalance.contract, event: "Approved", logs: logs, sub: sub}, nil
}

// WatchApproved is a free log subscription operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address retired, address approver, uint256 approversCount)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) WatchApproved(opts *bind.WatchOpts, sink chan<- *ITreasuryRebalanceApproved) (event.Subscription, error) {
	logs, sub, err := _ITreasuryRebalance.contract.WatchLogs(opts, "Approved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITreasuryRebalanceApproved)
				if err := _ITreasuryRebalance.contract.UnpackLog(event, "Approved", log); err != nil {
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

// ParseApproved is a log parse operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address retired, address approver, uint256 approversCount)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) ParseApproved(log types.Log) (*ITreasuryRebalanceApproved, error) {
	event := new(ITreasuryRebalanceApproved)
	if err := _ITreasuryRebalance.contract.UnpackLog(event, "Approved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ITreasuryRebalanceContractDeployedIterator is returned from FilterContractDeployed and is used to iterate over the raw logs and unpacked data for ContractDeployed events raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceContractDeployedIterator struct {
	Event *ITreasuryRebalanceContractDeployed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITreasuryRebalanceContractDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITreasuryRebalanceContractDeployed)
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
		it.Event = new(ITreasuryRebalanceContractDeployed)
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
func (it *ITreasuryRebalanceContractDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITreasuryRebalanceContractDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITreasuryRebalanceContractDeployed represents a ContractDeployed event raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceContractDeployed struct {
	Status               uint8
	RebalanceBlockNumber *big.Int
	DeployedBlockNumber  *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterContractDeployed is a free log retrieval operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) FilterContractDeployed(opts *bind.FilterOpts) (*ITreasuryRebalanceContractDeployedIterator, error) {
	logs, sub, err := _ITreasuryRebalance.contract.FilterLogs(opts, "ContractDeployed")
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceContractDeployedIterator{contract: _ITreasuryRebalance.contract, event: "ContractDeployed", logs: logs, sub: sub}, nil
}

// WatchContractDeployed is a free log subscription operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) WatchContractDeployed(opts *bind.WatchOpts, sink chan<- *ITreasuryRebalanceContractDeployed) (event.Subscription, error) {
	logs, sub, err := _ITreasuryRebalance.contract.WatchLogs(opts, "ContractDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITreasuryRebalanceContractDeployed)
				if err := _ITreasuryRebalance.contract.UnpackLog(event, "ContractDeployed", log); err != nil {
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

// ParseContractDeployed is a log parse operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) ParseContractDeployed(log types.Log) (*ITreasuryRebalanceContractDeployed, error) {
	event := new(ITreasuryRebalanceContractDeployed)
	if err := _ITreasuryRebalance.contract.UnpackLog(event, "ContractDeployed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ITreasuryRebalanceFinalizedIterator is returned from FilterFinalized and is used to iterate over the raw logs and unpacked data for Finalized events raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceFinalizedIterator struct {
	Event *ITreasuryRebalanceFinalized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITreasuryRebalanceFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITreasuryRebalanceFinalized)
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
		it.Event = new(ITreasuryRebalanceFinalized)
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
func (it *ITreasuryRebalanceFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITreasuryRebalanceFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITreasuryRebalanceFinalized represents a Finalized event raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceFinalized struct {
	Memo   string
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFinalized is a free log retrieval operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) FilterFinalized(opts *bind.FilterOpts) (*ITreasuryRebalanceFinalizedIterator, error) {
	logs, sub, err := _ITreasuryRebalance.contract.FilterLogs(opts, "Finalized")
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceFinalizedIterator{contract: _ITreasuryRebalance.contract, event: "Finalized", logs: logs, sub: sub}, nil
}

// WatchFinalized is a free log subscription operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) WatchFinalized(opts *bind.WatchOpts, sink chan<- *ITreasuryRebalanceFinalized) (event.Subscription, error) {
	logs, sub, err := _ITreasuryRebalance.contract.WatchLogs(opts, "Finalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITreasuryRebalanceFinalized)
				if err := _ITreasuryRebalance.contract.UnpackLog(event, "Finalized", log); err != nil {
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

// ParseFinalized is a log parse operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) ParseFinalized(log types.Log) (*ITreasuryRebalanceFinalized, error) {
	event := new(ITreasuryRebalanceFinalized)
	if err := _ITreasuryRebalance.contract.UnpackLog(event, "Finalized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ITreasuryRebalanceNewbieRegisteredIterator is returned from FilterNewbieRegistered and is used to iterate over the raw logs and unpacked data for NewbieRegistered events raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceNewbieRegisteredIterator struct {
	Event *ITreasuryRebalanceNewbieRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITreasuryRebalanceNewbieRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITreasuryRebalanceNewbieRegistered)
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
		it.Event = new(ITreasuryRebalanceNewbieRegistered)
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
func (it *ITreasuryRebalanceNewbieRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITreasuryRebalanceNewbieRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITreasuryRebalanceNewbieRegistered represents a NewbieRegistered event raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceNewbieRegistered struct {
	Newbie         common.Address
	FundAllocation *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterNewbieRegistered is a free log retrieval operation binding the contract event 0xd261b37cd56b21cd1af841dca6331a133e5d8b9d55c2c6fe0ec822e2a303ef74.
//
// Solidity: event NewbieRegistered(address newbie, uint256 fundAllocation)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) FilterNewbieRegistered(opts *bind.FilterOpts) (*ITreasuryRebalanceNewbieRegisteredIterator, error) {
	logs, sub, err := _ITreasuryRebalance.contract.FilterLogs(opts, "NewbieRegistered")
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceNewbieRegisteredIterator{contract: _ITreasuryRebalance.contract, event: "NewbieRegistered", logs: logs, sub: sub}, nil
}

// WatchNewbieRegistered is a free log subscription operation binding the contract event 0xd261b37cd56b21cd1af841dca6331a133e5d8b9d55c2c6fe0ec822e2a303ef74.
//
// Solidity: event NewbieRegistered(address newbie, uint256 fundAllocation)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) WatchNewbieRegistered(opts *bind.WatchOpts, sink chan<- *ITreasuryRebalanceNewbieRegistered) (event.Subscription, error) {
	logs, sub, err := _ITreasuryRebalance.contract.WatchLogs(opts, "NewbieRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITreasuryRebalanceNewbieRegistered)
				if err := _ITreasuryRebalance.contract.UnpackLog(event, "NewbieRegistered", log); err != nil {
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

// ParseNewbieRegistered is a log parse operation binding the contract event 0xd261b37cd56b21cd1af841dca6331a133e5d8b9d55c2c6fe0ec822e2a303ef74.
//
// Solidity: event NewbieRegistered(address newbie, uint256 fundAllocation)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) ParseNewbieRegistered(log types.Log) (*ITreasuryRebalanceNewbieRegistered, error) {
	event := new(ITreasuryRebalanceNewbieRegistered)
	if err := _ITreasuryRebalance.contract.UnpackLog(event, "NewbieRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ITreasuryRebalanceNewbieRemovedIterator is returned from FilterNewbieRemoved and is used to iterate over the raw logs and unpacked data for NewbieRemoved events raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceNewbieRemovedIterator struct {
	Event *ITreasuryRebalanceNewbieRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITreasuryRebalanceNewbieRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITreasuryRebalanceNewbieRemoved)
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
		it.Event = new(ITreasuryRebalanceNewbieRemoved)
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
func (it *ITreasuryRebalanceNewbieRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITreasuryRebalanceNewbieRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITreasuryRebalanceNewbieRemoved represents a NewbieRemoved event raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceNewbieRemoved struct {
	Newbie common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNewbieRemoved is a free log retrieval operation binding the contract event 0xe630072edaed8f0fccf534c7eaa063290db8f775b0824c7261d01e6619da4b38.
//
// Solidity: event NewbieRemoved(address newbie)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) FilterNewbieRemoved(opts *bind.FilterOpts) (*ITreasuryRebalanceNewbieRemovedIterator, error) {
	logs, sub, err := _ITreasuryRebalance.contract.FilterLogs(opts, "NewbieRemoved")
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceNewbieRemovedIterator{contract: _ITreasuryRebalance.contract, event: "NewbieRemoved", logs: logs, sub: sub}, nil
}

// WatchNewbieRemoved is a free log subscription operation binding the contract event 0xe630072edaed8f0fccf534c7eaa063290db8f775b0824c7261d01e6619da4b38.
//
// Solidity: event NewbieRemoved(address newbie)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) WatchNewbieRemoved(opts *bind.WatchOpts, sink chan<- *ITreasuryRebalanceNewbieRemoved) (event.Subscription, error) {
	logs, sub, err := _ITreasuryRebalance.contract.WatchLogs(opts, "NewbieRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITreasuryRebalanceNewbieRemoved)
				if err := _ITreasuryRebalance.contract.UnpackLog(event, "NewbieRemoved", log); err != nil {
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

// ParseNewbieRemoved is a log parse operation binding the contract event 0xe630072edaed8f0fccf534c7eaa063290db8f775b0824c7261d01e6619da4b38.
//
// Solidity: event NewbieRemoved(address newbie)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) ParseNewbieRemoved(log types.Log) (*ITreasuryRebalanceNewbieRemoved, error) {
	event := new(ITreasuryRebalanceNewbieRemoved)
	if err := _ITreasuryRebalance.contract.UnpackLog(event, "NewbieRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ITreasuryRebalanceRetiredRegisteredIterator is returned from FilterRetiredRegistered and is used to iterate over the raw logs and unpacked data for RetiredRegistered events raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceRetiredRegisteredIterator struct {
	Event *ITreasuryRebalanceRetiredRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITreasuryRebalanceRetiredRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITreasuryRebalanceRetiredRegistered)
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
		it.Event = new(ITreasuryRebalanceRetiredRegistered)
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
func (it *ITreasuryRebalanceRetiredRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITreasuryRebalanceRetiredRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITreasuryRebalanceRetiredRegistered represents a RetiredRegistered event raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceRetiredRegistered struct {
	Retired common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRetiredRegistered is a free log retrieval operation binding the contract event 0x7da2e87d0b02df1162d5736cc40dfcfffd17198aaf093ddff4a8f4eb26002fde.
//
// Solidity: event RetiredRegistered(address retired)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) FilterRetiredRegistered(opts *bind.FilterOpts) (*ITreasuryRebalanceRetiredRegisteredIterator, error) {
	logs, sub, err := _ITreasuryRebalance.contract.FilterLogs(opts, "RetiredRegistered")
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceRetiredRegisteredIterator{contract: _ITreasuryRebalance.contract, event: "RetiredRegistered", logs: logs, sub: sub}, nil
}

// WatchRetiredRegistered is a free log subscription operation binding the contract event 0x7da2e87d0b02df1162d5736cc40dfcfffd17198aaf093ddff4a8f4eb26002fde.
//
// Solidity: event RetiredRegistered(address retired)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) WatchRetiredRegistered(opts *bind.WatchOpts, sink chan<- *ITreasuryRebalanceRetiredRegistered) (event.Subscription, error) {
	logs, sub, err := _ITreasuryRebalance.contract.WatchLogs(opts, "RetiredRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITreasuryRebalanceRetiredRegistered)
				if err := _ITreasuryRebalance.contract.UnpackLog(event, "RetiredRegistered", log); err != nil {
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

// ParseRetiredRegistered is a log parse operation binding the contract event 0x7da2e87d0b02df1162d5736cc40dfcfffd17198aaf093ddff4a8f4eb26002fde.
//
// Solidity: event RetiredRegistered(address retired)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) ParseRetiredRegistered(log types.Log) (*ITreasuryRebalanceRetiredRegistered, error) {
	event := new(ITreasuryRebalanceRetiredRegistered)
	if err := _ITreasuryRebalance.contract.UnpackLog(event, "RetiredRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ITreasuryRebalanceRetiredRemovedIterator is returned from FilterRetiredRemoved and is used to iterate over the raw logs and unpacked data for RetiredRemoved events raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceRetiredRemovedIterator struct {
	Event *ITreasuryRebalanceRetiredRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITreasuryRebalanceRetiredRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITreasuryRebalanceRetiredRemoved)
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
		it.Event = new(ITreasuryRebalanceRetiredRemoved)
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
func (it *ITreasuryRebalanceRetiredRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITreasuryRebalanceRetiredRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITreasuryRebalanceRetiredRemoved represents a RetiredRemoved event raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceRetiredRemoved struct {
	Retired common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRetiredRemoved is a free log retrieval operation binding the contract event 0x1f46b11b62ae5cc6363d0d5c2e597c4cb8849543d9126353adb73c5d7215e237.
//
// Solidity: event RetiredRemoved(address retired)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) FilterRetiredRemoved(opts *bind.FilterOpts) (*ITreasuryRebalanceRetiredRemovedIterator, error) {
	logs, sub, err := _ITreasuryRebalance.contract.FilterLogs(opts, "RetiredRemoved")
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceRetiredRemovedIterator{contract: _ITreasuryRebalance.contract, event: "RetiredRemoved", logs: logs, sub: sub}, nil
}

// WatchRetiredRemoved is a free log subscription operation binding the contract event 0x1f46b11b62ae5cc6363d0d5c2e597c4cb8849543d9126353adb73c5d7215e237.
//
// Solidity: event RetiredRemoved(address retired)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) WatchRetiredRemoved(opts *bind.WatchOpts, sink chan<- *ITreasuryRebalanceRetiredRemoved) (event.Subscription, error) {
	logs, sub, err := _ITreasuryRebalance.contract.WatchLogs(opts, "RetiredRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITreasuryRebalanceRetiredRemoved)
				if err := _ITreasuryRebalance.contract.UnpackLog(event, "RetiredRemoved", log); err != nil {
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

// ParseRetiredRemoved is a log parse operation binding the contract event 0x1f46b11b62ae5cc6363d0d5c2e597c4cb8849543d9126353adb73c5d7215e237.
//
// Solidity: event RetiredRemoved(address retired)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) ParseRetiredRemoved(log types.Log) (*ITreasuryRebalanceRetiredRemoved, error) {
	event := new(ITreasuryRebalanceRetiredRemoved)
	if err := _ITreasuryRebalance.contract.UnpackLog(event, "RetiredRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ITreasuryRebalanceStatusChangedIterator is returned from FilterStatusChanged and is used to iterate over the raw logs and unpacked data for StatusChanged events raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceStatusChangedIterator struct {
	Event *ITreasuryRebalanceStatusChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ITreasuryRebalanceStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ITreasuryRebalanceStatusChanged)
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
		it.Event = new(ITreasuryRebalanceStatusChanged)
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
func (it *ITreasuryRebalanceStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ITreasuryRebalanceStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ITreasuryRebalanceStatusChanged represents a StatusChanged event raised by the ITreasuryRebalance contract.
type ITreasuryRebalanceStatusChanged struct {
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStatusChanged is a free log retrieval operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) FilterStatusChanged(opts *bind.FilterOpts) (*ITreasuryRebalanceStatusChangedIterator, error) {
	logs, sub, err := _ITreasuryRebalance.contract.FilterLogs(opts, "StatusChanged")
	if err != nil {
		return nil, err
	}
	return &ITreasuryRebalanceStatusChangedIterator{contract: _ITreasuryRebalance.contract, event: "StatusChanged", logs: logs, sub: sub}, nil
}

// WatchStatusChanged is a free log subscription operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) WatchStatusChanged(opts *bind.WatchOpts, sink chan<- *ITreasuryRebalanceStatusChanged) (event.Subscription, error) {
	logs, sub, err := _ITreasuryRebalance.contract.WatchLogs(opts, "StatusChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ITreasuryRebalanceStatusChanged)
				if err := _ITreasuryRebalance.contract.UnpackLog(event, "StatusChanged", log); err != nil {
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

// ParseStatusChanged is a log parse operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_ITreasuryRebalance *ITreasuryRebalanceFilterer) ParseStatusChanged(log types.Log) (*ITreasuryRebalanceStatusChanged, error) {
	event := new(ITreasuryRebalanceStatusChanged)
	if err := _ITreasuryRebalance.contract.UnpackLog(event, "StatusChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IZeroedContractMetaData contains all meta data concerning the IZeroedContract contract.
var IZeroedContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getState\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"adminList\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"quorom\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1865c57d": "getState()",
	},
}

// IZeroedContractABI is the input ABI used to generate the binding from.
// Deprecated: Use IZeroedContractMetaData.ABI instead.
var IZeroedContractABI = IZeroedContractMetaData.ABI

// IZeroedContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IZeroedContractBinRuntime = ``

// IZeroedContractFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IZeroedContractMetaData.Sigs instead.
var IZeroedContractFuncSigs = IZeroedContractMetaData.Sigs

// IZeroedContract is an auto generated Go binding around a Klaytn contract.
type IZeroedContract struct {
	IZeroedContractCaller     // Read-only binding to the contract
	IZeroedContractTransactor // Write-only binding to the contract
	IZeroedContractFilterer   // Log filterer for contract events
}

// IZeroedContractCaller is an auto generated read-only Go binding around a Klaytn contract.
type IZeroedContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IZeroedContractTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IZeroedContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IZeroedContractFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IZeroedContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IZeroedContractSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IZeroedContractSession struct {
	Contract     *IZeroedContract  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IZeroedContractCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IZeroedContractCallerSession struct {
	Contract *IZeroedContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IZeroedContractTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IZeroedContractTransactorSession struct {
	Contract     *IZeroedContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IZeroedContractRaw is an auto generated low-level Go binding around a Klaytn contract.
type IZeroedContractRaw struct {
	Contract *IZeroedContract // Generic contract binding to access the raw methods on
}

// IZeroedContractCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IZeroedContractCallerRaw struct {
	Contract *IZeroedContractCaller // Generic read-only contract binding to access the raw methods on
}

// IZeroedContractTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IZeroedContractTransactorRaw struct {
	Contract *IZeroedContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIZeroedContract creates a new instance of IZeroedContract, bound to a specific deployed contract.
func NewIZeroedContract(address common.Address, backend bind.ContractBackend) (*IZeroedContract, error) {
	contract, err := bindIZeroedContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IZeroedContract{IZeroedContractCaller: IZeroedContractCaller{contract: contract}, IZeroedContractTransactor: IZeroedContractTransactor{contract: contract}, IZeroedContractFilterer: IZeroedContractFilterer{contract: contract}}, nil
}

// NewIZeroedContractCaller creates a new read-only instance of IZeroedContract, bound to a specific deployed contract.
func NewIZeroedContractCaller(address common.Address, caller bind.ContractCaller) (*IZeroedContractCaller, error) {
	contract, err := bindIZeroedContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IZeroedContractCaller{contract: contract}, nil
}

// NewIZeroedContractTransactor creates a new write-only instance of IZeroedContract, bound to a specific deployed contract.
func NewIZeroedContractTransactor(address common.Address, transactor bind.ContractTransactor) (*IZeroedContractTransactor, error) {
	contract, err := bindIZeroedContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IZeroedContractTransactor{contract: contract}, nil
}

// NewIZeroedContractFilterer creates a new log filterer instance of IZeroedContract, bound to a specific deployed contract.
func NewIZeroedContractFilterer(address common.Address, filterer bind.ContractFilterer) (*IZeroedContractFilterer, error) {
	contract, err := bindIZeroedContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IZeroedContractFilterer{contract: contract}, nil
}

// bindIZeroedContract binds a generic wrapper to an already deployed contract.
func bindIZeroedContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IZeroedContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IZeroedContract *IZeroedContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IZeroedContract.Contract.IZeroedContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IZeroedContract *IZeroedContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IZeroedContract.Contract.IZeroedContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IZeroedContract *IZeroedContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IZeroedContract.Contract.IZeroedContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IZeroedContract *IZeroedContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IZeroedContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IZeroedContract *IZeroedContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IZeroedContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IZeroedContract *IZeroedContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IZeroedContract.Contract.contract.Transact(opts, method, params...)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 quorom)
func (_IZeroedContract *IZeroedContractCaller) GetState(opts *bind.CallOpts) (struct {
	AdminList []common.Address
	Quorom    *big.Int
}, error,
) {
	var out []interface{}
	err := _IZeroedContract.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		AdminList []common.Address
		Quorom    *big.Int
	})

	outstruct.AdminList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Quorom = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, err
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 quorom)
func (_IZeroedContract *IZeroedContractSession) GetState() (struct {
	AdminList []common.Address
	Quorom    *big.Int
}, error,
) {
	return _IZeroedContract.Contract.GetState(&_IZeroedContract.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 quorom)
func (_IZeroedContract *IZeroedContractCallerSession) GetState() (struct {
	AdminList []common.Address
	Quorom    *big.Int
}, error,
) {
	return _IZeroedContract.Contract.GetState(&_IZeroedContract.CallOpts)
}

// OwnableMetaData contains all meta data concerning the Ownable contract.
var OwnableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"8f32d59b": "isOwner()",
		"8da5cb5b": "owner()",
		"715018a6": "renounceOwnership()",
		"f2fde38b": "transferOwnership(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b50600080546001600160a01b0319163390811782556040519091907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a36102e18061005f6000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c8063715018a6146100515780638da5cb5b1461005b5780638f32d59b1461007b578063f2fde38b14610099575b600080fd5b6100596100ac565b005b6000546040516001600160a01b0390911681526020015b60405180910390f35b6000546001600160a01b031633146040519015158152602001610072565b6100596100a736600461027b565b610155565b6000546001600160a01b0316331461010b5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031633146101af5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610102565b6101b8816101bb565b50565b6001600160a01b0381166102205760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610102565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b60006020828403121561028d57600080fd5b81356001600160a01b03811681146102a457600080fd5b939250505056fea2646970667358221220f69393ded8e7101091799304ce90e6a51b8cf5e03b617676571c199bd79c95d964736f6c63430008130033",
}

// OwnableABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableMetaData.ABI instead.
var OwnableABI = OwnableMetaData.ABI

// OwnableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const OwnableBinRuntime = `608060405234801561001057600080fd5b506004361061004c5760003560e01c8063715018a6146100515780638da5cb5b1461005b5780638f32d59b1461007b578063f2fde38b14610099575b600080fd5b6100596100ac565b005b6000546040516001600160a01b0390911681526020015b60405180910390f35b6000546001600160a01b031633146040519015158152602001610072565b6100596100a736600461027b565b610155565b6000546001600160a01b0316331461010b5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064015b60405180910390fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031633146101af5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610102565b6101b8816101bb565b50565b6001600160a01b0381166102205760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610102565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b60006020828403121561028d57600080fd5b81356001600160a01b03811681146102a457600080fd5b939250505056fea2646970667358221220f69393ded8e7101091799304ce90e6a51b8cf5e03b617676571c199bd79c95d964736f6c63430008130033`

// OwnableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use OwnableMetaData.Sigs instead.
var OwnableFuncSigs = OwnableMetaData.Sigs

// OwnableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OwnableMetaData.Bin instead.
var OwnableBin = OwnableMetaData.Bin

// DeployOwnable deploys a new Klaytn contract, binding an instance of Ownable to it.
func DeployOwnable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Ownable, error) {
	parsed, err := OwnableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OwnableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// Ownable is an auto generated Go binding around a Klaytn contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around a Klaytn contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around a Klaytn contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type OwnableTransactorRaw struct {
	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
	contract, err := bindOwnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
	contract, err := bindOwnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableCaller{contract: contract}, nil
}

// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
	contract, err := bindOwnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableTransactor{contract: contract}, nil
}

// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
	contract, err := bindOwnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableFilterer{contract: contract}, nil
}

// bindOwnable binds a generic wrapper to an already deployed contract.
func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OwnableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Ownable *OwnableCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Ownable.contract.Call(opts, &out, "isOwner")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Ownable *OwnableSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Ownable *OwnableCallerSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ownable.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
type OwnableOwnershipTransferredIterator struct {
	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OwnableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipTransferred)
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
		it.Event = new(OwnableOwnershipTransferred)
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
func (it *OwnableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipTransferred)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Ownable *OwnableFilterer) ParseOwnershipTransferred(log types.Log) (*OwnableOwnershipTransferred, error) {
	event := new(OwnableOwnershipTransferred)
	if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceMetaData contains all meta data concerning the TreasuryRebalance contract.
var TreasuryRebalanceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rebalanceBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"retired\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"approversCount\",\"type\":\"uint256\"}],\"name\":\"Approved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumITreasuryRebalance.Status\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rebalanceBlockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deployedBlockNumber\",\"type\":\"uint256\"}],\"name\":\"ContractDeployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"enumITreasuryRebalance.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"Finalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newbie\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundAllocation\",\"type\":\"uint256\"}],\"name\":\"NewbieRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newbie\",\"type\":\"address\"}],\"name\":\"NewbieRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"retired\",\"type\":\"address\"}],\"name\":\"RetiredRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"retired\",\"type\":\"address\"}],\"name\":\"RetiredRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumITreasuryRebalance.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"StatusChanged\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_retiredAddress\",\"type\":\"address\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkRetiredsApproved\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"finalizeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeRegistration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newbieAddress\",\"type\":\"address\"}],\"name\":\"getNewbie\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNewbieCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newbieAddress\",\"type\":\"address\"}],\"name\":\"getNewbieIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_retiredAddress\",\"type\":\"address\"}],\"name\":\"getRetired\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRetiredCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_retiredAddress\",\"type\":\"address\"}],\"name\":\"getRetiredIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"treasuryAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isContractAddr\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"memo\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newbieAddress\",\"type\":\"address\"}],\"name\":\"newbieExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"newbies\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"newbie\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rebalanceBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newbieAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"registerNewbie\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_retiredAddress\",\"type\":\"address\"}],\"name\":\"registerRetired\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newbieAddress\",\"type\":\"address\"}],\"name\":\"removeNewbie\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_retiredAddress\",\"type\":\"address\"}],\"name\":\"removeRetired\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_retiredAddress\",\"type\":\"address\"}],\"name\":\"retiredExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"retirees\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"retired\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"status\",\"outputs\":[{\"internalType\":\"enumITreasuryRebalance.Status\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sumOfRetiredBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"retireesBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"daea85c5": "approve(address)",
		"966e0794": "checkRetiredsApproved()",
		"faaf9ca6": "finalizeApproval()",
		"ea6d4a9b": "finalizeContract(string)",
		"48409096": "finalizeRegistration()",
		"eb5a8e55": "getNewbie(address)",
		"91734d86": "getNewbieCount()",
		"11f5c466": "getNewbieIndex(address)",
		"bf680590": "getRetired(address)",
		"d1ed33fc": "getRetiredCount()",
		"681f6e7c": "getRetiredIndex(address)",
		"e20fcf00": "getTreasuryAmount()",
		"e2384cb3": "isContractAddr(address)",
		"8f32d59b": "isOwner()",
		"58c3b870": "memo()",
		"683e13cb": "newbieExists(address)",
		"94393e11": "newbies(uint256)",
		"8da5cb5b": "owner()",
		"49a3fb45": "rebalanceBlockNumber()",
		"652e27e0": "registerNewbie(address,uint256)",
		"1f8c1798": "registerRetired(address)",
		"6864b95b": "removeNewbie(address)",
		"1c1dac59": "removeRetired(address)",
		"715018a6": "renounceOwnership()",
		"d826f88f": "reset()",
		"01784e05": "retiredExists(address)",
		"5a12667b": "retirees(uint256)",
		"200d2ed2": "status()",
		"45205a6b": "sumOfRetiredBalance()",
		"f2fde38b": "transferOwnership(address)",
	},
	Bin: "0x60806040523480156200001157600080fd5b5060405162002696380380620026968339810160408190526200003491620000c8565b600080546001600160a01b0319163390811782556040519091907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a360048190556003805460ff191690556040517f6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a90620000b99060009084904290620000e2565b60405180910390a15062000119565b600060208284031215620000db57600080fd5b5051919050565b60608101600485106200010557634e487b7160e01b600052602160045260246000fd5b938152602081019290925260409091015290565b61256d80620001296000396000f3fe6080604052600436106101cd5760003560e01c80638da5cb5b116100f7578063d826f88f11610095578063ea6d4a9b11610064578063ea6d4a9b1461057d578063eb5a8e551461059d578063f2fde38b146105bd578063faaf9ca6146105dd576101cd565b8063d826f88f14610512578063daea85c514610527578063e20fcf0014610547578063e2384cb31461055c576101cd565b806394393e11116100d157806394393e111461047b578063966e0794146104ba578063bf680590146104cf578063d1ed33fc146104fd576101cd565b80638da5cb5b146104285780638f32d59b1461044657806391734d8614610466576101cd565b806349a3fb451161016f578063681f6e7c1161013e578063681f6e7c146103b3578063683e13cb146103d35780636864b95b146103f3578063715018a614610413576101cd565b806349a3fb451461032357806358c3b870146103395780635a12667b1461035b578063652e27e014610393576101cd565b80631f8c1798116101ab5780631f8c1798146102b2578063200d2ed2146102d257806345205a6b146102f9578063484090961461030e576101cd565b806301784e051461022d57806311f5c466146102625780631c1dac5914610290575b60405162461bcd60e51b815260206004820152602a60248201527f5468697320636f6e747261637420646f6573206e6f742061636365707420616e60448201526979207061796d656e747360b01b60648201526084015b60405180910390fd5b34801561023957600080fd5b5061024d610248366004611f0c565b6105f2565b60405190151581526020015b60405180910390f35b34801561026e57600080fd5b5061028261027d366004611f0c565b6106a6565b604051908152602001610259565b34801561029c57600080fd5b506102b06102ab366004611f0c565b610712565b005b3480156102be57600080fd5b506102b06102cd366004611f0c565b6108b0565b3480156102de57600080fd5b506003546102ec9060ff1681565b6040516102599190611f68565b34801561030557600080fd5b506102826109f5565b34801561031a57600080fd5b506102b0610a53565b34801561032f57600080fd5b5061028260045481565b34801561034557600080fd5b5061034e610b0a565b6040516102599190611f7c565b34801561036757600080fd5b5061037b610376366004611fca565b610b98565b6040516001600160a01b039091168152602001610259565b34801561039f57600080fd5b506102b06103ae366004611fe3565b610bc7565b3480156103bf57600080fd5b506102826103ce366004611f0c565b610db0565b3480156103df57600080fd5b5061024d6103ee366004611f0c565b610e12565b3480156103ff57600080fd5b506102b061040e366004611f0c565b610ec0565b34801561041f57600080fd5b506102b0611069565b34801561043457600080fd5b506000546001600160a01b031661037b565b34801561045257600080fd5b506000546001600160a01b0316331461024d565b34801561047257600080fd5b50600254610282565b34801561048757600080fd5b5061049b610496366004611fca565b6110dd565b604080516001600160a01b039093168352602083019190915201610259565b3480156104c657600080fd5b506102b0611115565b3480156104db57600080fd5b506104ef6104ea366004611f0c565b6112f9565b60405161025992919061200f565b34801561050957600080fd5b50600154610282565b34801561051e57600080fd5b506102b06113e0565b34801561053357600080fd5b506102b0610542366004611f0c565b6114bf565b34801561055357600080fd5b506102826116a3565b34801561056857600080fd5b5061024d610577366004611f0c565b3b151590565b34801561058957600080fd5b506102b06105983660046120b2565b6116f5565b3480156105a957600080fd5b5061049b6105b8366004611f0c565b61181d565b3480156105c957600080fd5b506102b06105d8366004611f0c565b6118cd565b3480156105e957600080fd5b506102b0611900565b60006001600160a01b03821661063c5760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964206164647265737360881b6044820152606401610224565b60005b6001548110156106a057826001600160a01b03166001828154811061066657610666612147565b60009182526020909120600290910201546001600160a01b03160361068e5750600192915050565b8061069881612173565b91505061063f565b50919050565b6000805b60025481101561070857826001600160a01b0316600282815481106106d1576106d1612147565b60009182526020909120600290910201546001600160a01b0316036106f65792915050565b8061070081612173565b9150506106aa565b5060001992915050565b6000546001600160a01b0316331461073c5760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff169081111561075657610756611f30565b146107735760405162461bcd60e51b8152600401610224906121c1565b600061077e83610db0565b905060001981036107a15760405162461bcd60e51b8152600401610224906121f8565b600180546107b0908290612228565b815481106107c0576107c0612147565b9060005260206000209060020201600182815481106107e1576107e1612147565b60009182526020909120825460029092020180546001600160a01b0319166001600160a01b03909216919091178155600180830180546108249284019190611dac565b5090505060018054806108395761083961223b565b60008281526020812060026000199093019283020180546001600160a01b03191681559061086a6001830182611df8565b505090556040516001600160a01b03841681527f1f46b11b62ae5cc6363d0d5c2e597c4cb8849543d9126353adb73c5d7215e237906020015b60405180910390a1505050565b6000546001600160a01b031633146108da5760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff16908111156108f4576108f4611f30565b146109115760405162461bcd60e51b8152600401610224906121c1565b61091a826105f2565b156109755760405162461bcd60e51b815260206004820152602560248201527f52657469726564206164647265737320697320616c72656164792072656769736044820152641d195c995960da1b6064820152608401610224565b6001805480820182556000919091526002027fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180546001600160a01b0384166001600160a01b0319909116811782556040519081527f7da2e87d0b02df1162d5736cc40dfcfffd17198aaf093ddff4a8f4eb26002fde906020016108a3565b6000805b600154811015610a4f5760018181548110610a1657610a16612147565b6000918252602090912060029091020154610a3b906001600160a01b03163183612251565b915080610a4781612173565b9150506109f9565b5090565b6000546001600160a01b03163314610a7d5760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff1690811115610a9757610a97611f30565b14610ab45760405162461bcd60e51b8152600401610224906121c1565b600380546001919060ff191682805b02179055506003546040517fafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e91610aff9160ff90911690611f68565b60405180910390a150565b60058054610b1790612264565b80601f0160208091040260200160405190810160405280929190818152602001828054610b4390612264565b8015610b905780601f10610b6557610100808354040283529160200191610b90565b820191906000526020600020905b815481529060010190602001808311610b7357829003601f168201915b505050505081565b60018181548110610ba857600080fd5b60009182526020909120600290910201546001600160a01b0316905081565b6000546001600160a01b03163314610bf15760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff1690811115610c0b57610c0b611f30565b14610c285760405162461bcd60e51b8152600401610224906121c1565b610c3183610e12565b15610c8a5760405162461bcd60e51b8152602060048201526024808201527f4e6577626965206164647265737320697320616c726561647920726567697374604482015263195c995960e21b6064820152608401610224565b81600003610cda5760405162461bcd60e51b815260206004820152601960248201527f416d6f756e742063616e6e6f742062652073657420746f2030000000000000006044820152606401610224565b6040805180820182526001600160a01b038581168083526020808401878152600280546001810182556000829052865191027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace81018054929096166001600160a01b031990921691909117909455517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf90930192909255835190815290810185905290917fd261b37cd56b21cd1af841dca6331a133e5d8b9d55c2c6fe0ec822e2a303ef7491015b60405180910390a150505050565b6000805b60015481101561070857826001600160a01b031660018281548110610ddb57610ddb612147565b60009182526020909120600290910201546001600160a01b031603610e005792915050565b80610e0a81612173565b915050610db4565b60006001600160a01b038216610e5c5760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964206164647265737360881b6044820152606401610224565b60005b6002548110156106a057826001600160a01b031660028281548110610e8657610e86612147565b60009182526020909120600290910201546001600160a01b031603610eae5750600192915050565b80610eb881612173565b915050610e5f565b6000546001600160a01b03163314610eea5760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff1690811115610f0457610f04611f30565b14610f215760405162461bcd60e51b8152600401610224906121c1565b6000610f2c836106a6565b90506000198103610f775760405162461bcd60e51b815260206004820152601560248201527413995dd89a59481b9bdd081c9959da5cdd195c9959605a1b6044820152606401610224565b60028054610f8790600190612228565b81548110610f9757610f97612147565b906000526020600020906002020160028281548110610fb857610fb8612147565b600091825260209091208254600292830290910180546001600160a01b0319166001600160a01b039092169190911781556001928301549201919091558054806110045761100461223b565b600082815260208082206002600019949094019384020180546001600160a01b03191681556001019190915591556040516001600160a01b03851681527fe630072edaed8f0fccf534c7eaa063290db8f775b0824c7261d01e6619da4b3891016108a3565b6000546001600160a01b031633146110935760405162461bcd60e51b81526004016102249061218c565b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b600281815481106110ed57600080fd5b6000918252602090912060029091020180546001909101546001600160a01b03909116915082565b60005b6001548110156112f65760006001828154811061113757611137612147565b6000918252602091829020604080518082018252600290930290910180546001600160a01b031683526001810180548351818702810187019094528084529394919385830193928301828280156111b757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611199575b505050505081525050905060006111d282600001513b151590565b90508015611297576000806111ea8460000151611a14565b915091508084602001515110156112135760405162461bcd60e51b815260040161022490612298565b60208401516000805b825181101561126d5761124883828151811061123a5761123a612147565b602002602001015186611a8d565b1561125b578161125781612173565b9250505b8061126581612173565b91505061121c565b508281101561128e5760405162461bcd60e51b815260040161022490612298565b505050506112e1565b8160200151516001146112e15760405162461bcd60e51b8152602060048201526012602482015271454f412073686f756c6420617070726f766560701b6044820152606401610224565b505080806112ee90612173565b915050611118565b50565b60006060600061130884610db0565b9050600019810361132b5760405162461bcd60e51b8152600401610224906121f8565b60006001828154811061134057611340612147565b6000918252602091829020604080518082018252600290930290910180546001600160a01b031683526001810180548351818702810187019094528084529394919385830193928301828280156113c057602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116113a2575b505050505081525050905080600001518160200151935093505050915091565b6000546001600160a01b0316331461140a5760405162461bcd60e51b81526004016102249061218c565b6003805460ff168181111561142157611421611f30565b14158015611430575060045443105b61148f5760405162461bcd60e51b815260206004820152602a60248201527f436f6e74726163742069732066696e616c697a65642c2063616e6e6f742072656044820152697365742076616c75657360b01b6064820152608401610224565b61149b60016000611e16565b6114a760026000611e37565b6114b360056000611e58565b6003805460ff19169055565b6001806003805460ff16908111156114d9576114d9611f30565b146114f65760405162461bcd60e51b8152600401610224906121c1565b6114ff826105f2565b6115625760405162461bcd60e51b815260206004820152602e60248201527f72657469726564206e6565647320746f2062652072656769737465726564206260448201526d19599bdc9948185c1c1c9bdd985b60921b6064820152608401610224565b813b1515806115de57336001600160a01b038416146115cf5760405162461bcd60e51b8152602060048201526024808201527f7265746972656441646472657373206973206e6f7420746865206d73672e7365604482015263373232b960e11b6064820152608401610224565b6115d98333611aea565b505050565b60006115e984611a14565b509050805160000361163d5760405162461bcd60e51b815260206004820152601a60248201527f61646d696e206c6973742063616e6e6f7420626520656d7074790000000000006044820152606401610224565b6116473382611a8d565b6116935760405162461bcd60e51b815260206004820152601b60248201527f6d73672e73656e646572206973206e6f74207468652061646d696e00000000006044820152606401610224565b61169d8433611aea565b50505050565b6000805b600254811015610a4f57600281815481106116c4576116c4612147565b906000526020600020906002020160010154826116e19190612251565b9150806116ed81612173565b9150506116a7565b6000546001600160a01b0316331461171f5760405162461bcd60e51b81526004016102249061218c565b6002806003805460ff169081111561173957611739611f30565b146117565760405162461bcd60e51b8152600401610224906121c1565b60056117628382612328565b506003805460ff1916811781556040517f8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404916117a191600591906123e8565b60405180910390a160045443116118195760405162461bcd60e51b815260206004820152603660248201527f436f6e74726163742063616e206f6e6c792066696e616c697a6520616674657260448201527520657865637574696e6720726562616c616e63696e6760501b6064820152608401610224565b5050565b600080600061182b846106a6565b905060001981036118765760405162461bcd60e51b815260206004820152601560248201527413995dd89a59481b9bdd081c9959da5cdd195c9959605a1b6044820152606401610224565b60006002828154811061188b5761188b612147565b60009182526020918290206040805180820190915260029092020180546001600160a01b03168083526001909101549190920181905290969095509350505050565b6000546001600160a01b031633146118f75760405162461bcd60e51b81526004016102249061218c565b6112f681611cec565b6000546001600160a01b0316331461192a5760405162461bcd60e51b81526004016102249061218c565b6001806003805460ff169081111561194457611944611f30565b146119615760405162461bcd60e51b8152600401610224906121c1565b6119696109f5565b6119716116a3565b106119f85760405162461bcd60e51b815260206004820152604b60248201527f747265617375727920616d6f756e742073686f756c64206265206c657373207460448201527f68616e207468652073756d206f6620616c6c207265746972656420616464726560648201526a73732062616c616e63657360a81b608482015260a401610224565b611a00611115565b600380546002919060ff1916600183610ac3565b6060600080839050806001600160a01b0316631865c57d6040518163ffffffff1660e01b8152600401600060405180830381865afa158015611a5a573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611a82919081019061247d565b909590945092505050565b6000805b8251811015611ae357828181518110611aac57611aac612147565b60200260200101516001600160a01b0316846001600160a01b031603611ad157600191505b80611adb81612173565b915050611a91565b5092915050565b6000611af583610db0565b90506000198103611b185760405162461bcd60e51b8152600401610224906121f8565b600060018281548110611b2d57611b2d612147565b9060005260206000209060020201600101805480602002602001604051908101604052809291908181526020018280548015611b9257602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611b74575b5050505050905060005b8151811015611c2457836001600160a01b0316828281518110611bc157611bc1612147565b60200260200101516001600160a01b031603611c125760405162461bcd60e51b815260206004820152601060248201526f105b1c9958591e48185c1c1c9bdd995960821b6044820152606401610224565b80611c1c81612173565b915050611b9c565b5060018281548110611c3857611c38612147565b600091825260208083206001600290930201820180548084018255908452922090910180546001600160a01b0386166001600160a01b031990911617905580547f80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f9091869186919086908110611caf57611caf612147565b600091825260209182902060016002909202010154604080516001600160a01b039586168152949093169184019190915290820152606001610da2565b6001600160a01b038116611d515760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610224565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b828054828255906000526020600020908101928215611dec5760005260206000209182015b82811115611dec578254825591600101919060010190611dd1565b50610a4f929150611e8e565b50805460008255906000526020600020908101906112f69190611e8e565b50805460008255600202906000526020600020908101906112f69190611ea3565b50805460008255600202906000526020600020908101906112f69190611ed1565b508054611e6490612264565b6000825580601f10611e74575050565b601f0160209004906000526020600020908101906112f691905b5b80821115610a4f5760008155600101611e8f565b80821115610a4f5780546001600160a01b03191681556000611ec86001830182611df8565b50600201611ea3565b5b80821115610a4f5780546001600160a01b031916815560006001820155600201611ed2565b6001600160a01b03811681146112f657600080fd5b600060208284031215611f1e57600080fd5b8135611f2981611ef7565b9392505050565b634e487b7160e01b600052602160045260246000fd5b60048110611f6457634e487b7160e01b600052602160045260246000fd5b9052565b60208101611f768284611f46565b92915050565b600060208083528351808285015260005b81811015611fa957858101830151858201604001528201611f8d565b506000604082860101526040601f19601f8301168501019250505092915050565b600060208284031215611fdc57600080fd5b5035919050565b60008060408385031215611ff657600080fd5b823561200181611ef7565b946020939093013593505050565b6001600160a01b038381168252604060208084018290528451918401829052600092858201929091906060860190855b8181101561205d57855185168352948301949183019160010161203f565b509098975050505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156120aa576120aa61206b565b604052919050565b600060208083850312156120c557600080fd5b823567ffffffffffffffff808211156120dd57600080fd5b818501915085601f8301126120f157600080fd5b8135818111156121035761210361206b565b612115601f8201601f19168501612081565b9150808252868482850101111561212b57600080fd5b8084840185840137600090820190930192909252509392505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016121855761218561215d565b5060010190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252601c908201527f4e6f7420696e207468652064657369676e617465642073746174757300000000604082015260600190565b60208082526016908201527514995d1a5c9959081b9bdd081c9959da5cdd195c995960521b604082015260600190565b81810381811115611f7657611f7661215d565b634e487b7160e01b600052603160045260246000fd5b80820180821115611f7657611f7661215d565b600181811c9082168061227857607f821691505b6020821081036106a057634e487b7160e01b600052602260045260246000fd5b60208082526022908201527f6d696e2072657175697265642061646d696e732073686f756c6420617070726f604082015261766560f01b606082015260800190565b601f8211156115d957600081815260208120601f850160051c810160208610156123015750805b601f850160051c820191505b818110156123205782815560010161230d565b505050505050565b815167ffffffffffffffff8111156123425761234261206b565b612356816123508454612264565b846122da565b602080601f83116001811461238b57600084156123735750858301515b600019600386901b1c1916600185901b178555612320565b600085815260208120601f198616915b828110156123ba5788860151825594840194600190910190840161239b565b50858210156123d85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6040815260008084546123fa81612264565b806040860152606060018084166000811461241c576001811461243657612467565b60ff1985168884015283151560051b880183019550612467565b8960005260208060002060005b8681101561245e5781548b8201870152908401908201612443565b8a018501975050505b505050505080915050611f296020830184611f46565b6000806040838503121561249057600080fd5b825167ffffffffffffffff808211156124a857600080fd5b818501915085601f8301126124bc57600080fd5b81516020828211156124d0576124d061206b565b8160051b92506124e1818401612081565b82815292840181019281810190898511156124fb57600080fd5b948201945b84861015612525578551935061251584611ef7565b8382529482019490820190612500565b9790910151969896975050505050505056fea26469706673582212204680134776f9249e58cbc4909a9b899ff12f8578cb7f31b1e2b3a2b1d44f65a064736f6c63430008130033",
}

// TreasuryRebalanceABI is the input ABI used to generate the binding from.
// Deprecated: Use TreasuryRebalanceMetaData.ABI instead.
var TreasuryRebalanceABI = TreasuryRebalanceMetaData.ABI

// TreasuryRebalanceBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const TreasuryRebalanceBinRuntime = `6080604052600436106101cd5760003560e01c80638da5cb5b116100f7578063d826f88f11610095578063ea6d4a9b11610064578063ea6d4a9b1461057d578063eb5a8e551461059d578063f2fde38b146105bd578063faaf9ca6146105dd576101cd565b8063d826f88f14610512578063daea85c514610527578063e20fcf0014610547578063e2384cb31461055c576101cd565b806394393e11116100d157806394393e111461047b578063966e0794146104ba578063bf680590146104cf578063d1ed33fc146104fd576101cd565b80638da5cb5b146104285780638f32d59b1461044657806391734d8614610466576101cd565b806349a3fb451161016f578063681f6e7c1161013e578063681f6e7c146103b3578063683e13cb146103d35780636864b95b146103f3578063715018a614610413576101cd565b806349a3fb451461032357806358c3b870146103395780635a12667b1461035b578063652e27e014610393576101cd565b80631f8c1798116101ab5780631f8c1798146102b2578063200d2ed2146102d257806345205a6b146102f9578063484090961461030e576101cd565b806301784e051461022d57806311f5c466146102625780631c1dac5914610290575b60405162461bcd60e51b815260206004820152602a60248201527f5468697320636f6e747261637420646f6573206e6f742061636365707420616e60448201526979207061796d656e747360b01b60648201526084015b60405180910390fd5b34801561023957600080fd5b5061024d610248366004611f0c565b6105f2565b60405190151581526020015b60405180910390f35b34801561026e57600080fd5b5061028261027d366004611f0c565b6106a6565b604051908152602001610259565b34801561029c57600080fd5b506102b06102ab366004611f0c565b610712565b005b3480156102be57600080fd5b506102b06102cd366004611f0c565b6108b0565b3480156102de57600080fd5b506003546102ec9060ff1681565b6040516102599190611f68565b34801561030557600080fd5b506102826109f5565b34801561031a57600080fd5b506102b0610a53565b34801561032f57600080fd5b5061028260045481565b34801561034557600080fd5b5061034e610b0a565b6040516102599190611f7c565b34801561036757600080fd5b5061037b610376366004611fca565b610b98565b6040516001600160a01b039091168152602001610259565b34801561039f57600080fd5b506102b06103ae366004611fe3565b610bc7565b3480156103bf57600080fd5b506102826103ce366004611f0c565b610db0565b3480156103df57600080fd5b5061024d6103ee366004611f0c565b610e12565b3480156103ff57600080fd5b506102b061040e366004611f0c565b610ec0565b34801561041f57600080fd5b506102b0611069565b34801561043457600080fd5b506000546001600160a01b031661037b565b34801561045257600080fd5b506000546001600160a01b0316331461024d565b34801561047257600080fd5b50600254610282565b34801561048757600080fd5b5061049b610496366004611fca565b6110dd565b604080516001600160a01b039093168352602083019190915201610259565b3480156104c657600080fd5b506102b0611115565b3480156104db57600080fd5b506104ef6104ea366004611f0c565b6112f9565b60405161025992919061200f565b34801561050957600080fd5b50600154610282565b34801561051e57600080fd5b506102b06113e0565b34801561053357600080fd5b506102b0610542366004611f0c565b6114bf565b34801561055357600080fd5b506102826116a3565b34801561056857600080fd5b5061024d610577366004611f0c565b3b151590565b34801561058957600080fd5b506102b06105983660046120b2565b6116f5565b3480156105a957600080fd5b5061049b6105b8366004611f0c565b61181d565b3480156105c957600080fd5b506102b06105d8366004611f0c565b6118cd565b3480156105e957600080fd5b506102b0611900565b60006001600160a01b03821661063c5760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964206164647265737360881b6044820152606401610224565b60005b6001548110156106a057826001600160a01b03166001828154811061066657610666612147565b60009182526020909120600290910201546001600160a01b03160361068e5750600192915050565b8061069881612173565b91505061063f565b50919050565b6000805b60025481101561070857826001600160a01b0316600282815481106106d1576106d1612147565b60009182526020909120600290910201546001600160a01b0316036106f65792915050565b8061070081612173565b9150506106aa565b5060001992915050565b6000546001600160a01b0316331461073c5760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff169081111561075657610756611f30565b146107735760405162461bcd60e51b8152600401610224906121c1565b600061077e83610db0565b905060001981036107a15760405162461bcd60e51b8152600401610224906121f8565b600180546107b0908290612228565b815481106107c0576107c0612147565b9060005260206000209060020201600182815481106107e1576107e1612147565b60009182526020909120825460029092020180546001600160a01b0319166001600160a01b03909216919091178155600180830180546108249284019190611dac565b5090505060018054806108395761083961223b565b60008281526020812060026000199093019283020180546001600160a01b03191681559061086a6001830182611df8565b505090556040516001600160a01b03841681527f1f46b11b62ae5cc6363d0d5c2e597c4cb8849543d9126353adb73c5d7215e237906020015b60405180910390a1505050565b6000546001600160a01b031633146108da5760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff16908111156108f4576108f4611f30565b146109115760405162461bcd60e51b8152600401610224906121c1565b61091a826105f2565b156109755760405162461bcd60e51b815260206004820152602560248201527f52657469726564206164647265737320697320616c72656164792072656769736044820152641d195c995960da1b6064820152608401610224565b6001805480820182556000919091526002027fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180546001600160a01b0384166001600160a01b0319909116811782556040519081527f7da2e87d0b02df1162d5736cc40dfcfffd17198aaf093ddff4a8f4eb26002fde906020016108a3565b6000805b600154811015610a4f5760018181548110610a1657610a16612147565b6000918252602090912060029091020154610a3b906001600160a01b03163183612251565b915080610a4781612173565b9150506109f9565b5090565b6000546001600160a01b03163314610a7d5760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff1690811115610a9757610a97611f30565b14610ab45760405162461bcd60e51b8152600401610224906121c1565b600380546001919060ff191682805b02179055506003546040517fafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e91610aff9160ff90911690611f68565b60405180910390a150565b60058054610b1790612264565b80601f0160208091040260200160405190810160405280929190818152602001828054610b4390612264565b8015610b905780601f10610b6557610100808354040283529160200191610b90565b820191906000526020600020905b815481529060010190602001808311610b7357829003601f168201915b505050505081565b60018181548110610ba857600080fd5b60009182526020909120600290910201546001600160a01b0316905081565b6000546001600160a01b03163314610bf15760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff1690811115610c0b57610c0b611f30565b14610c285760405162461bcd60e51b8152600401610224906121c1565b610c3183610e12565b15610c8a5760405162461bcd60e51b8152602060048201526024808201527f4e6577626965206164647265737320697320616c726561647920726567697374604482015263195c995960e21b6064820152608401610224565b81600003610cda5760405162461bcd60e51b815260206004820152601960248201527f416d6f756e742063616e6e6f742062652073657420746f2030000000000000006044820152606401610224565b6040805180820182526001600160a01b038581168083526020808401878152600280546001810182556000829052865191027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace81018054929096166001600160a01b031990921691909117909455517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf90930192909255835190815290810185905290917fd261b37cd56b21cd1af841dca6331a133e5d8b9d55c2c6fe0ec822e2a303ef7491015b60405180910390a150505050565b6000805b60015481101561070857826001600160a01b031660018281548110610ddb57610ddb612147565b60009182526020909120600290910201546001600160a01b031603610e005792915050565b80610e0a81612173565b915050610db4565b60006001600160a01b038216610e5c5760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964206164647265737360881b6044820152606401610224565b60005b6002548110156106a057826001600160a01b031660028281548110610e8657610e86612147565b60009182526020909120600290910201546001600160a01b031603610eae5750600192915050565b80610eb881612173565b915050610e5f565b6000546001600160a01b03163314610eea5760405162461bcd60e51b81526004016102249061218c565b6000806003805460ff1690811115610f0457610f04611f30565b14610f215760405162461bcd60e51b8152600401610224906121c1565b6000610f2c836106a6565b90506000198103610f775760405162461bcd60e51b815260206004820152601560248201527413995dd89a59481b9bdd081c9959da5cdd195c9959605a1b6044820152606401610224565b60028054610f8790600190612228565b81548110610f9757610f97612147565b906000526020600020906002020160028281548110610fb857610fb8612147565b600091825260209091208254600292830290910180546001600160a01b0319166001600160a01b039092169190911781556001928301549201919091558054806110045761100461223b565b600082815260208082206002600019949094019384020180546001600160a01b03191681556001019190915591556040516001600160a01b03851681527fe630072edaed8f0fccf534c7eaa063290db8f775b0824c7261d01e6619da4b3891016108a3565b6000546001600160a01b031633146110935760405162461bcd60e51b81526004016102249061218c565b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b600281815481106110ed57600080fd5b6000918252602090912060029091020180546001909101546001600160a01b03909116915082565b60005b6001548110156112f65760006001828154811061113757611137612147565b6000918252602091829020604080518082018252600290930290910180546001600160a01b031683526001810180548351818702810187019094528084529394919385830193928301828280156111b757602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611199575b505050505081525050905060006111d282600001513b151590565b90508015611297576000806111ea8460000151611a14565b915091508084602001515110156112135760405162461bcd60e51b815260040161022490612298565b60208401516000805b825181101561126d5761124883828151811061123a5761123a612147565b602002602001015186611a8d565b1561125b578161125781612173565b9250505b8061126581612173565b91505061121c565b508281101561128e5760405162461bcd60e51b815260040161022490612298565b505050506112e1565b8160200151516001146112e15760405162461bcd60e51b8152602060048201526012602482015271454f412073686f756c6420617070726f766560701b6044820152606401610224565b505080806112ee90612173565b915050611118565b50565b60006060600061130884610db0565b9050600019810361132b5760405162461bcd60e51b8152600401610224906121f8565b60006001828154811061134057611340612147565b6000918252602091829020604080518082018252600290930290910180546001600160a01b031683526001810180548351818702810187019094528084529394919385830193928301828280156113c057602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116113a2575b505050505081525050905080600001518160200151935093505050915091565b6000546001600160a01b0316331461140a5760405162461bcd60e51b81526004016102249061218c565b6003805460ff168181111561142157611421611f30565b14158015611430575060045443105b61148f5760405162461bcd60e51b815260206004820152602a60248201527f436f6e74726163742069732066696e616c697a65642c2063616e6e6f742072656044820152697365742076616c75657360b01b6064820152608401610224565b61149b60016000611e16565b6114a760026000611e37565b6114b360056000611e58565b6003805460ff19169055565b6001806003805460ff16908111156114d9576114d9611f30565b146114f65760405162461bcd60e51b8152600401610224906121c1565b6114ff826105f2565b6115625760405162461bcd60e51b815260206004820152602e60248201527f72657469726564206e6565647320746f2062652072656769737465726564206260448201526d19599bdc9948185c1c1c9bdd985b60921b6064820152608401610224565b813b1515806115de57336001600160a01b038416146115cf5760405162461bcd60e51b8152602060048201526024808201527f7265746972656441646472657373206973206e6f7420746865206d73672e7365604482015263373232b960e11b6064820152608401610224565b6115d98333611aea565b505050565b60006115e984611a14565b509050805160000361163d5760405162461bcd60e51b815260206004820152601a60248201527f61646d696e206c6973742063616e6e6f7420626520656d7074790000000000006044820152606401610224565b6116473382611a8d565b6116935760405162461bcd60e51b815260206004820152601b60248201527f6d73672e73656e646572206973206e6f74207468652061646d696e00000000006044820152606401610224565b61169d8433611aea565b50505050565b6000805b600254811015610a4f57600281815481106116c4576116c4612147565b906000526020600020906002020160010154826116e19190612251565b9150806116ed81612173565b9150506116a7565b6000546001600160a01b0316331461171f5760405162461bcd60e51b81526004016102249061218c565b6002806003805460ff169081111561173957611739611f30565b146117565760405162461bcd60e51b8152600401610224906121c1565b60056117628382612328565b506003805460ff1916811781556040517f8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404916117a191600591906123e8565b60405180910390a160045443116118195760405162461bcd60e51b815260206004820152603660248201527f436f6e74726163742063616e206f6e6c792066696e616c697a6520616674657260448201527520657865637574696e6720726562616c616e63696e6760501b6064820152608401610224565b5050565b600080600061182b846106a6565b905060001981036118765760405162461bcd60e51b815260206004820152601560248201527413995dd89a59481b9bdd081c9959da5cdd195c9959605a1b6044820152606401610224565b60006002828154811061188b5761188b612147565b60009182526020918290206040805180820190915260029092020180546001600160a01b03168083526001909101549190920181905290969095509350505050565b6000546001600160a01b031633146118f75760405162461bcd60e51b81526004016102249061218c565b6112f681611cec565b6000546001600160a01b0316331461192a5760405162461bcd60e51b81526004016102249061218c565b6001806003805460ff169081111561194457611944611f30565b146119615760405162461bcd60e51b8152600401610224906121c1565b6119696109f5565b6119716116a3565b106119f85760405162461bcd60e51b815260206004820152604b60248201527f747265617375727920616d6f756e742073686f756c64206265206c657373207460448201527f68616e207468652073756d206f6620616c6c207265746972656420616464726560648201526a73732062616c616e63657360a81b608482015260a401610224565b611a00611115565b600380546002919060ff1916600183610ac3565b6060600080839050806001600160a01b0316631865c57d6040518163ffffffff1660e01b8152600401600060405180830381865afa158015611a5a573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611a82919081019061247d565b909590945092505050565b6000805b8251811015611ae357828181518110611aac57611aac612147565b60200260200101516001600160a01b0316846001600160a01b031603611ad157600191505b80611adb81612173565b915050611a91565b5092915050565b6000611af583610db0565b90506000198103611b185760405162461bcd60e51b8152600401610224906121f8565b600060018281548110611b2d57611b2d612147565b9060005260206000209060020201600101805480602002602001604051908101604052809291908181526020018280548015611b9257602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611b74575b5050505050905060005b8151811015611c2457836001600160a01b0316828281518110611bc157611bc1612147565b60200260200101516001600160a01b031603611c125760405162461bcd60e51b815260206004820152601060248201526f105b1c9958591e48185c1c1c9bdd995960821b6044820152606401610224565b80611c1c81612173565b915050611b9c565b5060018281548110611c3857611c38612147565b600091825260208083206001600290930201820180548084018255908452922090910180546001600160a01b0386166001600160a01b031990911617905580547f80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f9091869186919086908110611caf57611caf612147565b600091825260209182902060016002909202010154604080516001600160a01b039586168152949093169184019190915290820152606001610da2565b6001600160a01b038116611d515760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610224565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b828054828255906000526020600020908101928215611dec5760005260206000209182015b82811115611dec578254825591600101919060010190611dd1565b50610a4f929150611e8e565b50805460008255906000526020600020908101906112f69190611e8e565b50805460008255600202906000526020600020908101906112f69190611ea3565b50805460008255600202906000526020600020908101906112f69190611ed1565b508054611e6490612264565b6000825580601f10611e74575050565b601f0160209004906000526020600020908101906112f691905b5b80821115610a4f5760008155600101611e8f565b80821115610a4f5780546001600160a01b03191681556000611ec86001830182611df8565b50600201611ea3565b5b80821115610a4f5780546001600160a01b031916815560006001820155600201611ed2565b6001600160a01b03811681146112f657600080fd5b600060208284031215611f1e57600080fd5b8135611f2981611ef7565b9392505050565b634e487b7160e01b600052602160045260246000fd5b60048110611f6457634e487b7160e01b600052602160045260246000fd5b9052565b60208101611f768284611f46565b92915050565b600060208083528351808285015260005b81811015611fa957858101830151858201604001528201611f8d565b506000604082860101526040601f19601f8301168501019250505092915050565b600060208284031215611fdc57600080fd5b5035919050565b60008060408385031215611ff657600080fd5b823561200181611ef7565b946020939093013593505050565b6001600160a01b038381168252604060208084018290528451918401829052600092858201929091906060860190855b8181101561205d57855185168352948301949183019160010161203f565b509098975050505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff811182821017156120aa576120aa61206b565b604052919050565b600060208083850312156120c557600080fd5b823567ffffffffffffffff808211156120dd57600080fd5b818501915085601f8301126120f157600080fd5b8135818111156121035761210361206b565b612115601f8201601f19168501612081565b9150808252868482850101111561212b57600080fd5b8084840185840137600090820190930192909252509392505050565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000600182016121855761218561215d565b5060010190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252601c908201527f4e6f7420696e207468652064657369676e617465642073746174757300000000604082015260600190565b60208082526016908201527514995d1a5c9959081b9bdd081c9959da5cdd195c995960521b604082015260600190565b81810381811115611f7657611f7661215d565b634e487b7160e01b600052603160045260246000fd5b80820180821115611f7657611f7661215d565b600181811c9082168061227857607f821691505b6020821081036106a057634e487b7160e01b600052602260045260246000fd5b60208082526022908201527f6d696e2072657175697265642061646d696e732073686f756c6420617070726f604082015261766560f01b606082015260800190565b601f8211156115d957600081815260208120601f850160051c810160208610156123015750805b601f850160051c820191505b818110156123205782815560010161230d565b505050505050565b815167ffffffffffffffff8111156123425761234261206b565b612356816123508454612264565b846122da565b602080601f83116001811461238b57600084156123735750858301515b600019600386901b1c1916600185901b178555612320565b600085815260208120601f198616915b828110156123ba5788860151825594840194600190910190840161239b565b50858210156123d85787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6040815260008084546123fa81612264565b806040860152606060018084166000811461241c576001811461243657612467565b60ff1985168884015283151560051b880183019550612467565b8960005260208060002060005b8681101561245e5781548b8201870152908401908201612443565b8a018501975050505b505050505080915050611f296020830184611f46565b6000806040838503121561249057600080fd5b825167ffffffffffffffff808211156124a857600080fd5b818501915085601f8301126124bc57600080fd5b81516020828211156124d0576124d061206b565b8160051b92506124e1818401612081565b82815292840181019281810190898511156124fb57600080fd5b948201945b84861015612525578551935061251584611ef7565b8382529482019490820190612500565b9790910151969896975050505050505056fea26469706673582212204680134776f9249e58cbc4909a9b899ff12f8578cb7f31b1e2b3a2b1d44f65a064736f6c63430008130033`

// TreasuryRebalanceFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use TreasuryRebalanceMetaData.Sigs instead.
var TreasuryRebalanceFuncSigs = TreasuryRebalanceMetaData.Sigs

// TreasuryRebalanceBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TreasuryRebalanceMetaData.Bin instead.
var TreasuryRebalanceBin = TreasuryRebalanceMetaData.Bin

// DeployTreasuryRebalance deploys a new Klaytn contract, binding an instance of TreasuryRebalance to it.
func DeployTreasuryRebalance(auth *bind.TransactOpts, backend bind.ContractBackend, _rebalanceBlockNumber *big.Int) (common.Address, *types.Transaction, *TreasuryRebalance, error) {
	parsed, err := TreasuryRebalanceMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TreasuryRebalanceBin), backend, _rebalanceBlockNumber)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TreasuryRebalance{TreasuryRebalanceCaller: TreasuryRebalanceCaller{contract: contract}, TreasuryRebalanceTransactor: TreasuryRebalanceTransactor{contract: contract}, TreasuryRebalanceFilterer: TreasuryRebalanceFilterer{contract: contract}}, nil
}

// TreasuryRebalance is an auto generated Go binding around a Klaytn contract.
type TreasuryRebalance struct {
	TreasuryRebalanceCaller     // Read-only binding to the contract
	TreasuryRebalanceTransactor // Write-only binding to the contract
	TreasuryRebalanceFilterer   // Log filterer for contract events
}

// TreasuryRebalanceCaller is an auto generated read-only Go binding around a Klaytn contract.
type TreasuryRebalanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreasuryRebalanceTransactor is an auto generated write-only Go binding around a Klaytn contract.
type TreasuryRebalanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreasuryRebalanceFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type TreasuryRebalanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreasuryRebalanceSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type TreasuryRebalanceSession struct {
	Contract     *TreasuryRebalance // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// TreasuryRebalanceCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type TreasuryRebalanceCallerSession struct {
	Contract *TreasuryRebalanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// TreasuryRebalanceTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type TreasuryRebalanceTransactorSession struct {
	Contract     *TreasuryRebalanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// TreasuryRebalanceRaw is an auto generated low-level Go binding around a Klaytn contract.
type TreasuryRebalanceRaw struct {
	Contract *TreasuryRebalance // Generic contract binding to access the raw methods on
}

// TreasuryRebalanceCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type TreasuryRebalanceCallerRaw struct {
	Contract *TreasuryRebalanceCaller // Generic read-only contract binding to access the raw methods on
}

// TreasuryRebalanceTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type TreasuryRebalanceTransactorRaw struct {
	Contract *TreasuryRebalanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTreasuryRebalance creates a new instance of TreasuryRebalance, bound to a specific deployed contract.
func NewTreasuryRebalance(address common.Address, backend bind.ContractBackend) (*TreasuryRebalance, error) {
	contract, err := bindTreasuryRebalance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalance{TreasuryRebalanceCaller: TreasuryRebalanceCaller{contract: contract}, TreasuryRebalanceTransactor: TreasuryRebalanceTransactor{contract: contract}, TreasuryRebalanceFilterer: TreasuryRebalanceFilterer{contract: contract}}, nil
}

// NewTreasuryRebalanceCaller creates a new read-only instance of TreasuryRebalance, bound to a specific deployed contract.
func NewTreasuryRebalanceCaller(address common.Address, caller bind.ContractCaller) (*TreasuryRebalanceCaller, error) {
	contract, err := bindTreasuryRebalance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceCaller{contract: contract}, nil
}

// NewTreasuryRebalanceTransactor creates a new write-only instance of TreasuryRebalance, bound to a specific deployed contract.
func NewTreasuryRebalanceTransactor(address common.Address, transactor bind.ContractTransactor) (*TreasuryRebalanceTransactor, error) {
	contract, err := bindTreasuryRebalance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceTransactor{contract: contract}, nil
}

// NewTreasuryRebalanceFilterer creates a new log filterer instance of TreasuryRebalance, bound to a specific deployed contract.
func NewTreasuryRebalanceFilterer(address common.Address, filterer bind.ContractFilterer) (*TreasuryRebalanceFilterer, error) {
	contract, err := bindTreasuryRebalance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceFilterer{contract: contract}, nil
}

// bindTreasuryRebalance binds a generic wrapper to an already deployed contract.
func bindTreasuryRebalance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TreasuryRebalanceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TreasuryRebalance *TreasuryRebalanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TreasuryRebalance.Contract.TreasuryRebalanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TreasuryRebalance *TreasuryRebalanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.TreasuryRebalanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TreasuryRebalance *TreasuryRebalanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.TreasuryRebalanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TreasuryRebalance *TreasuryRebalanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TreasuryRebalance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TreasuryRebalance *TreasuryRebalanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TreasuryRebalance *TreasuryRebalanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.contract.Transact(opts, method, params...)
}

// CheckRetiredsApproved is a free data retrieval call binding the contract method 0x966e0794.
//
// Solidity: function checkRetiredsApproved() view returns()
func (_TreasuryRebalance *TreasuryRebalanceCaller) CheckRetiredsApproved(opts *bind.CallOpts) error {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "checkRetiredsApproved")
	if err != nil {
		return err
	}

	return err
}

// CheckRetiredsApproved is a free data retrieval call binding the contract method 0x966e0794.
//
// Solidity: function checkRetiredsApproved() view returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) CheckRetiredsApproved() error {
	return _TreasuryRebalance.Contract.CheckRetiredsApproved(&_TreasuryRebalance.CallOpts)
}

// CheckRetiredsApproved is a free data retrieval call binding the contract method 0x966e0794.
//
// Solidity: function checkRetiredsApproved() view returns()
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) CheckRetiredsApproved() error {
	return _TreasuryRebalance.Contract.CheckRetiredsApproved(&_TreasuryRebalance.CallOpts)
}

// GetNewbie is a free data retrieval call binding the contract method 0xeb5a8e55.
//
// Solidity: function getNewbie(address _newbieAddress) view returns(address, uint256)
func (_TreasuryRebalance *TreasuryRebalanceCaller) GetNewbie(opts *bind.CallOpts, _newbieAddress common.Address) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "getNewbie", _newbieAddress)
	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err
}

// GetNewbie is a free data retrieval call binding the contract method 0xeb5a8e55.
//
// Solidity: function getNewbie(address _newbieAddress) view returns(address, uint256)
func (_TreasuryRebalance *TreasuryRebalanceSession) GetNewbie(_newbieAddress common.Address) (common.Address, *big.Int, error) {
	return _TreasuryRebalance.Contract.GetNewbie(&_TreasuryRebalance.CallOpts, _newbieAddress)
}

// GetNewbie is a free data retrieval call binding the contract method 0xeb5a8e55.
//
// Solidity: function getNewbie(address _newbieAddress) view returns(address, uint256)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) GetNewbie(_newbieAddress common.Address) (common.Address, *big.Int, error) {
	return _TreasuryRebalance.Contract.GetNewbie(&_TreasuryRebalance.CallOpts, _newbieAddress)
}

// GetNewbieCount is a free data retrieval call binding the contract method 0x91734d86.
//
// Solidity: function getNewbieCount() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCaller) GetNewbieCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "getNewbieCount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetNewbieCount is a free data retrieval call binding the contract method 0x91734d86.
//
// Solidity: function getNewbieCount() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceSession) GetNewbieCount() (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetNewbieCount(&_TreasuryRebalance.CallOpts)
}

// GetNewbieCount is a free data retrieval call binding the contract method 0x91734d86.
//
// Solidity: function getNewbieCount() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) GetNewbieCount() (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetNewbieCount(&_TreasuryRebalance.CallOpts)
}

// GetNewbieIndex is a free data retrieval call binding the contract method 0x11f5c466.
//
// Solidity: function getNewbieIndex(address _newbieAddress) view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCaller) GetNewbieIndex(opts *bind.CallOpts, _newbieAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "getNewbieIndex", _newbieAddress)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetNewbieIndex is a free data retrieval call binding the contract method 0x11f5c466.
//
// Solidity: function getNewbieIndex(address _newbieAddress) view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceSession) GetNewbieIndex(_newbieAddress common.Address) (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetNewbieIndex(&_TreasuryRebalance.CallOpts, _newbieAddress)
}

// GetNewbieIndex is a free data retrieval call binding the contract method 0x11f5c466.
//
// Solidity: function getNewbieIndex(address _newbieAddress) view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) GetNewbieIndex(_newbieAddress common.Address) (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetNewbieIndex(&_TreasuryRebalance.CallOpts, _newbieAddress)
}

// GetRetired is a free data retrieval call binding the contract method 0xbf680590.
//
// Solidity: function getRetired(address _retiredAddress) view returns(address, address[])
func (_TreasuryRebalance *TreasuryRebalanceCaller) GetRetired(opts *bind.CallOpts, _retiredAddress common.Address) (common.Address, []common.Address, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "getRetired", _retiredAddress)
	if err != nil {
		return *new(common.Address), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err
}

// GetRetired is a free data retrieval call binding the contract method 0xbf680590.
//
// Solidity: function getRetired(address _retiredAddress) view returns(address, address[])
func (_TreasuryRebalance *TreasuryRebalanceSession) GetRetired(_retiredAddress common.Address) (common.Address, []common.Address, error) {
	return _TreasuryRebalance.Contract.GetRetired(&_TreasuryRebalance.CallOpts, _retiredAddress)
}

// GetRetired is a free data retrieval call binding the contract method 0xbf680590.
//
// Solidity: function getRetired(address _retiredAddress) view returns(address, address[])
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) GetRetired(_retiredAddress common.Address) (common.Address, []common.Address, error) {
	return _TreasuryRebalance.Contract.GetRetired(&_TreasuryRebalance.CallOpts, _retiredAddress)
}

// GetRetiredCount is a free data retrieval call binding the contract method 0xd1ed33fc.
//
// Solidity: function getRetiredCount() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCaller) GetRetiredCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "getRetiredCount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetRetiredCount is a free data retrieval call binding the contract method 0xd1ed33fc.
//
// Solidity: function getRetiredCount() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceSession) GetRetiredCount() (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetRetiredCount(&_TreasuryRebalance.CallOpts)
}

// GetRetiredCount is a free data retrieval call binding the contract method 0xd1ed33fc.
//
// Solidity: function getRetiredCount() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) GetRetiredCount() (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetRetiredCount(&_TreasuryRebalance.CallOpts)
}

// GetRetiredIndex is a free data retrieval call binding the contract method 0x681f6e7c.
//
// Solidity: function getRetiredIndex(address _retiredAddress) view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCaller) GetRetiredIndex(opts *bind.CallOpts, _retiredAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "getRetiredIndex", _retiredAddress)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetRetiredIndex is a free data retrieval call binding the contract method 0x681f6e7c.
//
// Solidity: function getRetiredIndex(address _retiredAddress) view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceSession) GetRetiredIndex(_retiredAddress common.Address) (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetRetiredIndex(&_TreasuryRebalance.CallOpts, _retiredAddress)
}

// GetRetiredIndex is a free data retrieval call binding the contract method 0x681f6e7c.
//
// Solidity: function getRetiredIndex(address _retiredAddress) view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) GetRetiredIndex(_retiredAddress common.Address) (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetRetiredIndex(&_TreasuryRebalance.CallOpts, _retiredAddress)
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_TreasuryRebalance *TreasuryRebalanceCaller) GetTreasuryAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "getTreasuryAmount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_TreasuryRebalance *TreasuryRebalanceSession) GetTreasuryAmount() (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetTreasuryAmount(&_TreasuryRebalance.CallOpts)
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) GetTreasuryAmount() (*big.Int, error) {
	return _TreasuryRebalance.Contract.GetTreasuryAmount(&_TreasuryRebalance.CallOpts)
}

// IsContractAddr is a free data retrieval call binding the contract method 0xe2384cb3.
//
// Solidity: function isContractAddr(address _addr) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceCaller) IsContractAddr(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "isContractAddr", _addr)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsContractAddr is a free data retrieval call binding the contract method 0xe2384cb3.
//
// Solidity: function isContractAddr(address _addr) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceSession) IsContractAddr(_addr common.Address) (bool, error) {
	return _TreasuryRebalance.Contract.IsContractAddr(&_TreasuryRebalance.CallOpts, _addr)
}

// IsContractAddr is a free data retrieval call binding the contract method 0xe2384cb3.
//
// Solidity: function isContractAddr(address _addr) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) IsContractAddr(_addr common.Address) (bool, error) {
	return _TreasuryRebalance.Contract.IsContractAddr(&_TreasuryRebalance.CallOpts, _addr)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "isOwner")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceSession) IsOwner() (bool, error) {
	return _TreasuryRebalance.Contract.IsOwner(&_TreasuryRebalance.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) IsOwner() (bool, error) {
	return _TreasuryRebalance.Contract.IsOwner(&_TreasuryRebalance.CallOpts)
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_TreasuryRebalance *TreasuryRebalanceCaller) Memo(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "memo")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_TreasuryRebalance *TreasuryRebalanceSession) Memo() (string, error) {
	return _TreasuryRebalance.Contract.Memo(&_TreasuryRebalance.CallOpts)
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) Memo() (string, error) {
	return _TreasuryRebalance.Contract.Memo(&_TreasuryRebalance.CallOpts)
}

// NewbieExists is a free data retrieval call binding the contract method 0x683e13cb.
//
// Solidity: function newbieExists(address _newbieAddress) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceCaller) NewbieExists(opts *bind.CallOpts, _newbieAddress common.Address) (bool, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "newbieExists", _newbieAddress)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// NewbieExists is a free data retrieval call binding the contract method 0x683e13cb.
//
// Solidity: function newbieExists(address _newbieAddress) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceSession) NewbieExists(_newbieAddress common.Address) (bool, error) {
	return _TreasuryRebalance.Contract.NewbieExists(&_TreasuryRebalance.CallOpts, _newbieAddress)
}

// NewbieExists is a free data retrieval call binding the contract method 0x683e13cb.
//
// Solidity: function newbieExists(address _newbieAddress) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) NewbieExists(_newbieAddress common.Address) (bool, error) {
	return _TreasuryRebalance.Contract.NewbieExists(&_TreasuryRebalance.CallOpts, _newbieAddress)
}

// Newbies is a free data retrieval call binding the contract method 0x94393e11.
//
// Solidity: function newbies(uint256 ) view returns(address newbie, uint256 amount)
func (_TreasuryRebalance *TreasuryRebalanceCaller) Newbies(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Newbie common.Address
	Amount *big.Int
}, error,
) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "newbies", arg0)

	outstruct := new(struct {
		Newbie common.Address
		Amount *big.Int
	})

	outstruct.Newbie = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, err
}

// Newbies is a free data retrieval call binding the contract method 0x94393e11.
//
// Solidity: function newbies(uint256 ) view returns(address newbie, uint256 amount)
func (_TreasuryRebalance *TreasuryRebalanceSession) Newbies(arg0 *big.Int) (struct {
	Newbie common.Address
	Amount *big.Int
}, error,
) {
	return _TreasuryRebalance.Contract.Newbies(&_TreasuryRebalance.CallOpts, arg0)
}

// Newbies is a free data retrieval call binding the contract method 0x94393e11.
//
// Solidity: function newbies(uint256 ) view returns(address newbie, uint256 amount)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) Newbies(arg0 *big.Int) (struct {
	Newbie common.Address
	Amount *big.Int
}, error,
) {
	return _TreasuryRebalance.Contract.Newbies(&_TreasuryRebalance.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TreasuryRebalance *TreasuryRebalanceCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TreasuryRebalance *TreasuryRebalanceSession) Owner() (common.Address, error) {
	return _TreasuryRebalance.Contract.Owner(&_TreasuryRebalance.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) Owner() (common.Address, error) {
	return _TreasuryRebalance.Contract.Owner(&_TreasuryRebalance.CallOpts)
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCaller) RebalanceBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "rebalanceBlockNumber")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceSession) RebalanceBlockNumber() (*big.Int, error) {
	return _TreasuryRebalance.Contract.RebalanceBlockNumber(&_TreasuryRebalance.CallOpts)
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) RebalanceBlockNumber() (*big.Int, error) {
	return _TreasuryRebalance.Contract.RebalanceBlockNumber(&_TreasuryRebalance.CallOpts)
}

// RetiredExists is a free data retrieval call binding the contract method 0x01784e05.
//
// Solidity: function retiredExists(address _retiredAddress) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceCaller) RetiredExists(opts *bind.CallOpts, _retiredAddress common.Address) (bool, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "retiredExists", _retiredAddress)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// RetiredExists is a free data retrieval call binding the contract method 0x01784e05.
//
// Solidity: function retiredExists(address _retiredAddress) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceSession) RetiredExists(_retiredAddress common.Address) (bool, error) {
	return _TreasuryRebalance.Contract.RetiredExists(&_TreasuryRebalance.CallOpts, _retiredAddress)
}

// RetiredExists is a free data retrieval call binding the contract method 0x01784e05.
//
// Solidity: function retiredExists(address _retiredAddress) view returns(bool)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) RetiredExists(_retiredAddress common.Address) (bool, error) {
	return _TreasuryRebalance.Contract.RetiredExists(&_TreasuryRebalance.CallOpts, _retiredAddress)
}

// Retirees is a free data retrieval call binding the contract method 0x5a12667b.
//
// Solidity: function retirees(uint256 ) view returns(address retired)
func (_TreasuryRebalance *TreasuryRebalanceCaller) Retirees(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "retirees", arg0)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Retirees is a free data retrieval call binding the contract method 0x5a12667b.
//
// Solidity: function retirees(uint256 ) view returns(address retired)
func (_TreasuryRebalance *TreasuryRebalanceSession) Retirees(arg0 *big.Int) (common.Address, error) {
	return _TreasuryRebalance.Contract.Retirees(&_TreasuryRebalance.CallOpts, arg0)
}

// Retirees is a free data retrieval call binding the contract method 0x5a12667b.
//
// Solidity: function retirees(uint256 ) view returns(address retired)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) Retirees(arg0 *big.Int) (common.Address, error) {
	return _TreasuryRebalance.Contract.Retirees(&_TreasuryRebalance.CallOpts, arg0)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_TreasuryRebalance *TreasuryRebalanceCaller) Status(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "status")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_TreasuryRebalance *TreasuryRebalanceSession) Status() (uint8, error) {
	return _TreasuryRebalance.Contract.Status(&_TreasuryRebalance.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) Status() (uint8, error) {
	return _TreasuryRebalance.Contract.Status(&_TreasuryRebalance.CallOpts)
}

// SumOfRetiredBalance is a free data retrieval call binding the contract method 0x45205a6b.
//
// Solidity: function sumOfRetiredBalance() view returns(uint256 retireesBalance)
func (_TreasuryRebalance *TreasuryRebalanceCaller) SumOfRetiredBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalance.contract.Call(opts, &out, "sumOfRetiredBalance")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// SumOfRetiredBalance is a free data retrieval call binding the contract method 0x45205a6b.
//
// Solidity: function sumOfRetiredBalance() view returns(uint256 retireesBalance)
func (_TreasuryRebalance *TreasuryRebalanceSession) SumOfRetiredBalance() (*big.Int, error) {
	return _TreasuryRebalance.Contract.SumOfRetiredBalance(&_TreasuryRebalance.CallOpts)
}

// SumOfRetiredBalance is a free data retrieval call binding the contract method 0x45205a6b.
//
// Solidity: function sumOfRetiredBalance() view returns(uint256 retireesBalance)
func (_TreasuryRebalance *TreasuryRebalanceCallerSession) SumOfRetiredBalance() (*big.Int, error) {
	return _TreasuryRebalance.Contract.SumOfRetiredBalance(&_TreasuryRebalance.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) Approve(opts *bind.TransactOpts, _retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "approve", _retiredAddress)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) Approve(_retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.Approve(&_TreasuryRebalance.TransactOpts, _retiredAddress)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) Approve(_retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.Approve(&_TreasuryRebalance.TransactOpts, _retiredAddress)
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) FinalizeApproval(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "finalizeApproval")
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) FinalizeApproval() (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.FinalizeApproval(&_TreasuryRebalance.TransactOpts)
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) FinalizeApproval() (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.FinalizeApproval(&_TreasuryRebalance.TransactOpts)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string _memo) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) FinalizeContract(opts *bind.TransactOpts, _memo string) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "finalizeContract", _memo)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string _memo) returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) FinalizeContract(_memo string) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.FinalizeContract(&_TreasuryRebalance.TransactOpts, _memo)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string _memo) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) FinalizeContract(_memo string) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.FinalizeContract(&_TreasuryRebalance.TransactOpts, _memo)
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) FinalizeRegistration(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "finalizeRegistration")
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) FinalizeRegistration() (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.FinalizeRegistration(&_TreasuryRebalance.TransactOpts)
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) FinalizeRegistration() (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.FinalizeRegistration(&_TreasuryRebalance.TransactOpts)
}

// RegisterNewbie is a paid mutator transaction binding the contract method 0x652e27e0.
//
// Solidity: function registerNewbie(address _newbieAddress, uint256 _amount) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) RegisterNewbie(opts *bind.TransactOpts, _newbieAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "registerNewbie", _newbieAddress, _amount)
}

// RegisterNewbie is a paid mutator transaction binding the contract method 0x652e27e0.
//
// Solidity: function registerNewbie(address _newbieAddress, uint256 _amount) returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) RegisterNewbie(_newbieAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RegisterNewbie(&_TreasuryRebalance.TransactOpts, _newbieAddress, _amount)
}

// RegisterNewbie is a paid mutator transaction binding the contract method 0x652e27e0.
//
// Solidity: function registerNewbie(address _newbieAddress, uint256 _amount) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) RegisterNewbie(_newbieAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RegisterNewbie(&_TreasuryRebalance.TransactOpts, _newbieAddress, _amount)
}

// RegisterRetired is a paid mutator transaction binding the contract method 0x1f8c1798.
//
// Solidity: function registerRetired(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) RegisterRetired(opts *bind.TransactOpts, _retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "registerRetired", _retiredAddress)
}

// RegisterRetired is a paid mutator transaction binding the contract method 0x1f8c1798.
//
// Solidity: function registerRetired(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) RegisterRetired(_retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RegisterRetired(&_TreasuryRebalance.TransactOpts, _retiredAddress)
}

// RegisterRetired is a paid mutator transaction binding the contract method 0x1f8c1798.
//
// Solidity: function registerRetired(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) RegisterRetired(_retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RegisterRetired(&_TreasuryRebalance.TransactOpts, _retiredAddress)
}

// RemoveNewbie is a paid mutator transaction binding the contract method 0x6864b95b.
//
// Solidity: function removeNewbie(address _newbieAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) RemoveNewbie(opts *bind.TransactOpts, _newbieAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "removeNewbie", _newbieAddress)
}

// RemoveNewbie is a paid mutator transaction binding the contract method 0x6864b95b.
//
// Solidity: function removeNewbie(address _newbieAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) RemoveNewbie(_newbieAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RemoveNewbie(&_TreasuryRebalance.TransactOpts, _newbieAddress)
}

// RemoveNewbie is a paid mutator transaction binding the contract method 0x6864b95b.
//
// Solidity: function removeNewbie(address _newbieAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) RemoveNewbie(_newbieAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RemoveNewbie(&_TreasuryRebalance.TransactOpts, _newbieAddress)
}

// RemoveRetired is a paid mutator transaction binding the contract method 0x1c1dac59.
//
// Solidity: function removeRetired(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) RemoveRetired(opts *bind.TransactOpts, _retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "removeRetired", _retiredAddress)
}

// RemoveRetired is a paid mutator transaction binding the contract method 0x1c1dac59.
//
// Solidity: function removeRetired(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) RemoveRetired(_retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RemoveRetired(&_TreasuryRebalance.TransactOpts, _retiredAddress)
}

// RemoveRetired is a paid mutator transaction binding the contract method 0x1c1dac59.
//
// Solidity: function removeRetired(address _retiredAddress) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) RemoveRetired(_retiredAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RemoveRetired(&_TreasuryRebalance.TransactOpts, _retiredAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) RenounceOwnership() (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RenounceOwnership(&_TreasuryRebalance.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.RenounceOwnership(&_TreasuryRebalance.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "reset")
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) Reset() (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.Reset(&_TreasuryRebalance.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) Reset() (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.Reset(&_TreasuryRebalance.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.TransferOwnership(&_TreasuryRebalance.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.TransferOwnership(&_TreasuryRebalance.TransactOpts, newOwner)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _TreasuryRebalance.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TreasuryRebalance *TreasuryRebalanceSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.Fallback(&_TreasuryRebalance.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_TreasuryRebalance *TreasuryRebalanceTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _TreasuryRebalance.Contract.Fallback(&_TreasuryRebalance.TransactOpts, calldata)
}

// TreasuryRebalanceApprovedIterator is returned from FilterApproved and is used to iterate over the raw logs and unpacked data for Approved events raised by the TreasuryRebalance contract.
type TreasuryRebalanceApprovedIterator struct {
	Event *TreasuryRebalanceApproved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceApproved)
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
		it.Event = new(TreasuryRebalanceApproved)
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
func (it *TreasuryRebalanceApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceApproved represents a Approved event raised by the TreasuryRebalance contract.
type TreasuryRebalanceApproved struct {
	Retired        common.Address
	Approver       common.Address
	ApproversCount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterApproved is a free log retrieval operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address retired, address approver, uint256 approversCount)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterApproved(opts *bind.FilterOpts) (*TreasuryRebalanceApprovedIterator, error) {
	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "Approved")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceApprovedIterator{contract: _TreasuryRebalance.contract, event: "Approved", logs: logs, sub: sub}, nil
}

// WatchApproved is a free log subscription operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address retired, address approver, uint256 approversCount)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchApproved(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceApproved) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "Approved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceApproved)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "Approved", log); err != nil {
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

// ParseApproved is a log parse operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address retired, address approver, uint256 approversCount)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseApproved(log types.Log) (*TreasuryRebalanceApproved, error) {
	event := new(TreasuryRebalanceApproved)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "Approved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceContractDeployedIterator is returned from FilterContractDeployed and is used to iterate over the raw logs and unpacked data for ContractDeployed events raised by the TreasuryRebalance contract.
type TreasuryRebalanceContractDeployedIterator struct {
	Event *TreasuryRebalanceContractDeployed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceContractDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceContractDeployed)
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
		it.Event = new(TreasuryRebalanceContractDeployed)
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
func (it *TreasuryRebalanceContractDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceContractDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceContractDeployed represents a ContractDeployed event raised by the TreasuryRebalance contract.
type TreasuryRebalanceContractDeployed struct {
	Status               uint8
	RebalanceBlockNumber *big.Int
	DeployedBlockNumber  *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterContractDeployed is a free log retrieval operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterContractDeployed(opts *bind.FilterOpts) (*TreasuryRebalanceContractDeployedIterator, error) {
	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "ContractDeployed")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceContractDeployedIterator{contract: _TreasuryRebalance.contract, event: "ContractDeployed", logs: logs, sub: sub}, nil
}

// WatchContractDeployed is a free log subscription operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchContractDeployed(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceContractDeployed) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "ContractDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceContractDeployed)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "ContractDeployed", log); err != nil {
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

// ParseContractDeployed is a log parse operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseContractDeployed(log types.Log) (*TreasuryRebalanceContractDeployed, error) {
	event := new(TreasuryRebalanceContractDeployed)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "ContractDeployed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceFinalizedIterator is returned from FilterFinalized and is used to iterate over the raw logs and unpacked data for Finalized events raised by the TreasuryRebalance contract.
type TreasuryRebalanceFinalizedIterator struct {
	Event *TreasuryRebalanceFinalized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceFinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceFinalized)
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
		it.Event = new(TreasuryRebalanceFinalized)
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
func (it *TreasuryRebalanceFinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceFinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceFinalized represents a Finalized event raised by the TreasuryRebalance contract.
type TreasuryRebalanceFinalized struct {
	Memo   string
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFinalized is a free log retrieval operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterFinalized(opts *bind.FilterOpts) (*TreasuryRebalanceFinalizedIterator, error) {
	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "Finalized")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceFinalizedIterator{contract: _TreasuryRebalance.contract, event: "Finalized", logs: logs, sub: sub}, nil
}

// WatchFinalized is a free log subscription operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchFinalized(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceFinalized) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "Finalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceFinalized)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "Finalized", log); err != nil {
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

// ParseFinalized is a log parse operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseFinalized(log types.Log) (*TreasuryRebalanceFinalized, error) {
	event := new(TreasuryRebalanceFinalized)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "Finalized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceNewbieRegisteredIterator is returned from FilterNewbieRegistered and is used to iterate over the raw logs and unpacked data for NewbieRegistered events raised by the TreasuryRebalance contract.
type TreasuryRebalanceNewbieRegisteredIterator struct {
	Event *TreasuryRebalanceNewbieRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceNewbieRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceNewbieRegistered)
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
		it.Event = new(TreasuryRebalanceNewbieRegistered)
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
func (it *TreasuryRebalanceNewbieRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceNewbieRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceNewbieRegistered represents a NewbieRegistered event raised by the TreasuryRebalance contract.
type TreasuryRebalanceNewbieRegistered struct {
	Newbie         common.Address
	FundAllocation *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterNewbieRegistered is a free log retrieval operation binding the contract event 0xd261b37cd56b21cd1af841dca6331a133e5d8b9d55c2c6fe0ec822e2a303ef74.
//
// Solidity: event NewbieRegistered(address newbie, uint256 fundAllocation)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterNewbieRegistered(opts *bind.FilterOpts) (*TreasuryRebalanceNewbieRegisteredIterator, error) {
	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "NewbieRegistered")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceNewbieRegisteredIterator{contract: _TreasuryRebalance.contract, event: "NewbieRegistered", logs: logs, sub: sub}, nil
}

// WatchNewbieRegistered is a free log subscription operation binding the contract event 0xd261b37cd56b21cd1af841dca6331a133e5d8b9d55c2c6fe0ec822e2a303ef74.
//
// Solidity: event NewbieRegistered(address newbie, uint256 fundAllocation)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchNewbieRegistered(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceNewbieRegistered) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "NewbieRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceNewbieRegistered)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "NewbieRegistered", log); err != nil {
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

// ParseNewbieRegistered is a log parse operation binding the contract event 0xd261b37cd56b21cd1af841dca6331a133e5d8b9d55c2c6fe0ec822e2a303ef74.
//
// Solidity: event NewbieRegistered(address newbie, uint256 fundAllocation)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseNewbieRegistered(log types.Log) (*TreasuryRebalanceNewbieRegistered, error) {
	event := new(TreasuryRebalanceNewbieRegistered)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "NewbieRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceNewbieRemovedIterator is returned from FilterNewbieRemoved and is used to iterate over the raw logs and unpacked data for NewbieRemoved events raised by the TreasuryRebalance contract.
type TreasuryRebalanceNewbieRemovedIterator struct {
	Event *TreasuryRebalanceNewbieRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceNewbieRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceNewbieRemoved)
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
		it.Event = new(TreasuryRebalanceNewbieRemoved)
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
func (it *TreasuryRebalanceNewbieRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceNewbieRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceNewbieRemoved represents a NewbieRemoved event raised by the TreasuryRebalance contract.
type TreasuryRebalanceNewbieRemoved struct {
	Newbie common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterNewbieRemoved is a free log retrieval operation binding the contract event 0xe630072edaed8f0fccf534c7eaa063290db8f775b0824c7261d01e6619da4b38.
//
// Solidity: event NewbieRemoved(address newbie)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterNewbieRemoved(opts *bind.FilterOpts) (*TreasuryRebalanceNewbieRemovedIterator, error) {
	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "NewbieRemoved")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceNewbieRemovedIterator{contract: _TreasuryRebalance.contract, event: "NewbieRemoved", logs: logs, sub: sub}, nil
}

// WatchNewbieRemoved is a free log subscription operation binding the contract event 0xe630072edaed8f0fccf534c7eaa063290db8f775b0824c7261d01e6619da4b38.
//
// Solidity: event NewbieRemoved(address newbie)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchNewbieRemoved(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceNewbieRemoved) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "NewbieRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceNewbieRemoved)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "NewbieRemoved", log); err != nil {
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

// ParseNewbieRemoved is a log parse operation binding the contract event 0xe630072edaed8f0fccf534c7eaa063290db8f775b0824c7261d01e6619da4b38.
//
// Solidity: event NewbieRemoved(address newbie)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseNewbieRemoved(log types.Log) (*TreasuryRebalanceNewbieRemoved, error) {
	event := new(TreasuryRebalanceNewbieRemoved)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "NewbieRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TreasuryRebalance contract.
type TreasuryRebalanceOwnershipTransferredIterator struct {
	Event *TreasuryRebalanceOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceOwnershipTransferred)
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
		it.Event = new(TreasuryRebalanceOwnershipTransferred)
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
func (it *TreasuryRebalanceOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceOwnershipTransferred represents a OwnershipTransferred event raised by the TreasuryRebalance contract.
type TreasuryRebalanceOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TreasuryRebalanceOwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceOwnershipTransferredIterator{contract: _TreasuryRebalance.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceOwnershipTransferred)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseOwnershipTransferred(log types.Log) (*TreasuryRebalanceOwnershipTransferred, error) {
	event := new(TreasuryRebalanceOwnershipTransferred)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceRetiredRegisteredIterator is returned from FilterRetiredRegistered and is used to iterate over the raw logs and unpacked data for RetiredRegistered events raised by the TreasuryRebalance contract.
type TreasuryRebalanceRetiredRegisteredIterator struct {
	Event *TreasuryRebalanceRetiredRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceRetiredRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceRetiredRegistered)
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
		it.Event = new(TreasuryRebalanceRetiredRegistered)
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
func (it *TreasuryRebalanceRetiredRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceRetiredRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceRetiredRegistered represents a RetiredRegistered event raised by the TreasuryRebalance contract.
type TreasuryRebalanceRetiredRegistered struct {
	Retired common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRetiredRegistered is a free log retrieval operation binding the contract event 0x7da2e87d0b02df1162d5736cc40dfcfffd17198aaf093ddff4a8f4eb26002fde.
//
// Solidity: event RetiredRegistered(address retired)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterRetiredRegistered(opts *bind.FilterOpts) (*TreasuryRebalanceRetiredRegisteredIterator, error) {
	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "RetiredRegistered")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceRetiredRegisteredIterator{contract: _TreasuryRebalance.contract, event: "RetiredRegistered", logs: logs, sub: sub}, nil
}

// WatchRetiredRegistered is a free log subscription operation binding the contract event 0x7da2e87d0b02df1162d5736cc40dfcfffd17198aaf093ddff4a8f4eb26002fde.
//
// Solidity: event RetiredRegistered(address retired)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchRetiredRegistered(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceRetiredRegistered) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "RetiredRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceRetiredRegistered)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "RetiredRegistered", log); err != nil {
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

// ParseRetiredRegistered is a log parse operation binding the contract event 0x7da2e87d0b02df1162d5736cc40dfcfffd17198aaf093ddff4a8f4eb26002fde.
//
// Solidity: event RetiredRegistered(address retired)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseRetiredRegistered(log types.Log) (*TreasuryRebalanceRetiredRegistered, error) {
	event := new(TreasuryRebalanceRetiredRegistered)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "RetiredRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceRetiredRemovedIterator is returned from FilterRetiredRemoved and is used to iterate over the raw logs and unpacked data for RetiredRemoved events raised by the TreasuryRebalance contract.
type TreasuryRebalanceRetiredRemovedIterator struct {
	Event *TreasuryRebalanceRetiredRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceRetiredRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceRetiredRemoved)
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
		it.Event = new(TreasuryRebalanceRetiredRemoved)
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
func (it *TreasuryRebalanceRetiredRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceRetiredRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceRetiredRemoved represents a RetiredRemoved event raised by the TreasuryRebalance contract.
type TreasuryRebalanceRetiredRemoved struct {
	Retired common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRetiredRemoved is a free log retrieval operation binding the contract event 0x1f46b11b62ae5cc6363d0d5c2e597c4cb8849543d9126353adb73c5d7215e237.
//
// Solidity: event RetiredRemoved(address retired)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterRetiredRemoved(opts *bind.FilterOpts) (*TreasuryRebalanceRetiredRemovedIterator, error) {
	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "RetiredRemoved")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceRetiredRemovedIterator{contract: _TreasuryRebalance.contract, event: "RetiredRemoved", logs: logs, sub: sub}, nil
}

// WatchRetiredRemoved is a free log subscription operation binding the contract event 0x1f46b11b62ae5cc6363d0d5c2e597c4cb8849543d9126353adb73c5d7215e237.
//
// Solidity: event RetiredRemoved(address retired)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchRetiredRemoved(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceRetiredRemoved) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "RetiredRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceRetiredRemoved)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "RetiredRemoved", log); err != nil {
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

// ParseRetiredRemoved is a log parse operation binding the contract event 0x1f46b11b62ae5cc6363d0d5c2e597c4cb8849543d9126353adb73c5d7215e237.
//
// Solidity: event RetiredRemoved(address retired)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseRetiredRemoved(log types.Log) (*TreasuryRebalanceRetiredRemoved, error) {
	event := new(TreasuryRebalanceRetiredRemoved)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "RetiredRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceStatusChangedIterator is returned from FilterStatusChanged and is used to iterate over the raw logs and unpacked data for StatusChanged events raised by the TreasuryRebalance contract.
type TreasuryRebalanceStatusChangedIterator struct {
	Event *TreasuryRebalanceStatusChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceStatusChanged)
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
		it.Event = new(TreasuryRebalanceStatusChanged)
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
func (it *TreasuryRebalanceStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceStatusChanged represents a StatusChanged event raised by the TreasuryRebalance contract.
type TreasuryRebalanceStatusChanged struct {
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStatusChanged is a free log retrieval operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) FilterStatusChanged(opts *bind.FilterOpts) (*TreasuryRebalanceStatusChangedIterator, error) {
	logs, sub, err := _TreasuryRebalance.contract.FilterLogs(opts, "StatusChanged")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceStatusChangedIterator{contract: _TreasuryRebalance.contract, event: "StatusChanged", logs: logs, sub: sub}, nil
}

// WatchStatusChanged is a free log subscription operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) WatchStatusChanged(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceStatusChanged) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalance.contract.WatchLogs(opts, "StatusChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceStatusChanged)
				if err := _TreasuryRebalance.contract.UnpackLog(event, "StatusChanged", log); err != nil {
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

// ParseStatusChanged is a log parse operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_TreasuryRebalance *TreasuryRebalanceFilterer) ParseStatusChanged(log types.Log) (*TreasuryRebalanceStatusChanged, error) {
	event := new(TreasuryRebalanceStatusChanged)
	if err := _TreasuryRebalance.contract.UnpackLog(event, "StatusChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2MetaData contains all meta data concerning the TreasuryRebalanceV2 contract.
var TreasuryRebalanceV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rebalanceBlockNumber\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"allocated\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"fundAllocation\",\"type\":\"uint256\"}],\"name\":\"AllocatedRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"allocated\",\"type\":\"address\"}],\"name\":\"AllocatedRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"zeroed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"approversCount\",\"type\":\"uint256\"}],\"name\":\"Approved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumTreasuryRebalanceV2.Status\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rebalanceBlockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deployedBlockNumber\",\"type\":\"uint256\"}],\"name\":\"ContractDeployed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"memo\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"enumTreasuryRebalanceV2.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"Finalized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumTreasuryRebalanceV2.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"StatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"zeroed\",\"type\":\"address\"}],\"name\":\"ZeroedRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"zeroed\",\"type\":\"address\"}],\"name\":\"ZeroedRemoved\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_allocatedAddress\",\"type\":\"address\"}],\"name\":\"allocatedExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allocateds\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_zeroedAddress\",\"type\":\"address\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"checkZeroedsApproved\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_memo\",\"type\":\"string\"}],\"name\":\"finalizeContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeRegistration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_allocatedAddress\",\"type\":\"address\"}],\"name\":\"getAllocated\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllocatedCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_allocatedAddress\",\"type\":\"address\"}],\"name\":\"getAllocatedIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTreasuryAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"treasuryAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_zeroedAddress\",\"type\":\"address\"}],\"name\":\"getZeroed\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getZeroedCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_zeroedAddress\",\"type\":\"address\"}],\"name\":\"getZeroedIndex\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isContractAddr\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"memo\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rebalanceBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_allocatedAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"registerAllocated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_zeroedAddress\",\"type\":\"address\"}],\"name\":\"registerZeroed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_allocatedAddress\",\"type\":\"address\"}],\"name\":\"removeAllocated\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_zeroedAddress\",\"type\":\"address\"}],\"name\":\"removeZeroed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"status\",\"outputs\":[{\"internalType\":\"enumTreasuryRebalanceV2.Status\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sumOfZeroedBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"zeroedsBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rebalanceBlockNumber\",\"type\":\"uint256\"}],\"name\":\"updateRebalanceBlocknumber\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_zeroedAddress\",\"type\":\"address\"}],\"name\":\"zeroedExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"zeroeds\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"bd786f57": "allocatedExists(address)",
		"343e2c85": "allocateds(uint256)",
		"daea85c5": "approve(address)",
		"0287d126": "checkZeroedsApproved()",
		"faaf9ca6": "finalizeApproval()",
		"ea6d4a9b": "finalizeContract(string)",
		"48409096": "finalizeRegistration()",
		"9e59eb14": "getAllocated(address)",
		"ed355529": "getAllocatedCount()",
		"7bfaf7b7": "getAllocatedIndex(address)",
		"e20fcf00": "getTreasuryAmount()",
		"cea1c338": "getZeroed(address)",
		"9dc954ba": "getZeroedCount()",
		"518592da": "getZeroedIndex(address)",
		"e2384cb3": "isContractAddr(address)",
		"8f32d59b": "isOwner()",
		"58c3b870": "memo()",
		"8da5cb5b": "owner()",
		"49a3fb45": "rebalanceBlockNumber()",
		"ecd86778": "registerAllocated(address,uint256)",
		"5f9b0df7": "registerZeroed(address)",
		"27704cb5": "removeAllocated(address)",
		"db27b50b": "removeZeroed(address)",
		"715018a6": "renounceOwnership()",
		"d826f88f": "reset()",
		"200d2ed2": "status()",
		"9ab29b70": "sumOfZeroedBalance()",
		"f2fde38b": "transferOwnership(address)",
		"1804692f": "updateRebalanceBlocknumber(uint256)",
		"5f8798c0": "zeroedExists(address)",
		"62aa3e91": "zeroeds(uint256)",
	},
	Bin: "0x60806040523480156200001157600080fd5b50604051620026e5380380620026e5833981016040819052620000349162000142565b600080546001600160a01b0319163390811782556040519091907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a3438111620000ed5760405162461bcd60e51b815260206004820152603a60248201527f726562616c616e636520626c6f636b4e756d6265722073686f756c642062652060448201527f67726561746572207468616e2063757272656e7420626c6f636b000000000000606482015260840160405180910390fd5b60048190556003805460ff191690556040517f6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a906200013390600090849042906200015c565b60405180910390a15062000193565b6000602082840312156200015557600080fd5b5051919050565b60608101600485106200017f57634e487b7160e01b600052602160045260246000fd5b938152602081019290925260409091015290565b61254280620001a36000396000f3fe608060405234801561001057600080fd5b50600436106101e55760003560e01c80638f32d59b1161010f578063db27b50b116100a2578063ecd8677811610071578063ecd867781461040a578063ed3555291461041d578063f2fde38b14610425578063faaf9ca61461043857600080fd5b8063db27b50b146103c8578063e20fcf00146103db578063e2384cb3146103e3578063ea6d4a9b146103f757600080fd5b8063bd786f57116100de578063bd786f5714610379578063cea1c3381461038c578063d826f88f146103ad578063daea85c5146103b557600080fd5b80638f32d59b146103435780639ab29b70146103565780639dc954ba1461035e5780639e59eb141461036657600080fd5b8063518592da1161018757806362aa3e911161015657806362aa3e91146102ec578063715018a6146103175780637bfaf7b71461031f5780638da5cb5b1461033257600080fd5b8063518592da1461028e57806358c3b870146102a15780635f8798c0146102b65780635f9b0df7146102d957600080fd5b806327704cb5116101c357806327704cb51461022a578063343e2c851461023d578063484090961461026f57806349a3fb451461027757600080fd5b80630287d126146101ea5780631804692f146101f4578063200d2ed214610207575b600080fd5b6101f2610440565b005b6101f2610202366004611ed3565b61062d565b6003546102149060ff1681565b6040516102219190611f24565b60405180910390f35b6101f2610238366004611f47565b610748565b61025061024b366004611ed3565b6108fd565b604080516001600160a01b039093168352602083019190915201610221565b6101f2610935565b61028060045481565b604051908152602001610221565b61028061029c366004611f47565b6109ec565b6102a9610a58565b6040516102219190611f6b565b6102c96102c4366004611f47565b610ae6565b6040519015158152602001610221565b6101f26102e7366004611f47565b610b9a565b6102ff6102fa366004611ed3565b610cdd565b6040516001600160a01b039091168152602001610221565b6101f2610d0c565b61028061032d366004611f47565b610d80565b6000546001600160a01b03166102ff565b6000546001600160a01b031633146102c9565b610280610de2565b600154610280565b610250610374366004611f47565b610e40565b6102c9610387366004611f47565b610ef3565b61039f61039a366004611f47565b610fa1565b604051610221929190611fb9565b6101f2611088565b6101f26103c3366004611f47565b6111b6565b6101f26103d6366004611f47565b61139c565b610280611531565b6102c96103f1366004611f47565b3b151590565b6101f261040536600461205c565b611583565b6101f26104183660046120f1565b6116ab565b600254610280565b6101f2610433366004611f47565b611898565b6101f26118cb565b60005b60015481101561062a576000600182815481106104625761046261211d565b6000918252602091829020604080518082018252600290930290910180546001600160a01b031683526001810180548351818702810187019094528084529394919385830193928301828280156104e257602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116104c4575b505050505081525050905060006104fd82600001513b151590565b905080156105cb5760008061051584600001516119de565b915091508084602001515110156105475760405162461bcd60e51b815260040161053e90612133565b60405180910390fd5b60208401516000805b82518110156105a15761057c83828151811061056e5761056e61211d565b602002602001015186611a57565b1561058f578161058b8161218b565b9250505b806105998161218b565b915050610550565b50828110156105c25760405162461bcd60e51b815260040161053e90612133565b50505050610615565b8160200151516001146106155760405162461bcd60e51b8152602060048201526012602482015271454f412073686f756c6420617070726f766560701b604482015260640161053e565b505080806106229061218b565b915050610443565b50565b6000546001600160a01b031633146106575760405162461bcd60e51b815260040161053e906121a4565b60045443106106ce5760405162461bcd60e51b815260206004820152603e60248201527f63757272656e7420626c6f636b2073686f756c646e277420626520706173742060448201527f7468652063757272656e746c792073657420626c6f636b206e756d6265720000606482015260840161053e565b8043106107435760405162461bcd60e51b815260206004820152603a60248201527f726562616c616e636520626c6f636b4e756d6265722073686f756c642062652060448201527f67726561746572207468616e2063757272656e7420626c6f636b000000000000606482015260840161053e565b600455565b6000546001600160a01b031633146107725760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff169081111561078c5761078c611eec565b146107a95760405162461bcd60e51b815260040161053e906121d9565b60006107b483610d80565b905060001981036108025760405162461bcd60e51b8152602060048201526018602482015277105b1b1bd8d85d1959081b9bdd081c9959da5cdd195c995960421b604482015260640161053e565b6002805461081290600190612210565b815481106108225761082261211d565b9060005260206000209060020201600282815481106108435761084361211d565b600091825260209091208254600292830290910180546001600160a01b0319166001600160a01b0390921691909117815560019283015492019190915580548061088f5761088f612223565b600082815260208082206002600019949094019384020180546001600160a01b03191681556001019190915591556040516001600160a01b03851681527ff8f67464bea52432645435be9c46c427173a75aefaa1001272e08a4b8572f06e91015b60405180910390a1505050565b6002818154811061090d57600080fd5b6000918252602090912060029091020180546001909101546001600160a01b03909116915082565b6000546001600160a01b0316331461095f5760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff169081111561097957610979611eec565b146109965760405162461bcd60e51b815260040161053e906121d9565b600380546001919060ff191682805b02179055506003546040517fafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e916109e19160ff90911690611f24565b60405180910390a150565b6000805b600154811015610a4e57826001600160a01b031660018281548110610a1757610a1761211d565b60009182526020909120600290910201546001600160a01b031603610a3c5792915050565b80610a468161218b565b9150506109f0565b5060001992915050565b60058054610a6590612239565b80601f0160208091040260200160405190810160405280929190818152602001828054610a9190612239565b8015610ade5780601f10610ab357610100808354040283529160200191610ade565b820191906000526020600020905b815481529060010190602001808311610ac157829003601f168201915b505050505081565b60006001600160a01b038216610b305760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964206164647265737360881b604482015260640161053e565b60005b600154811015610b9457826001600160a01b031660018281548110610b5a57610b5a61211d565b60009182526020909120600290910201546001600160a01b031603610b825750600192915050565b80610b8c8161218b565b915050610b33565b50919050565b6000546001600160a01b03163314610bc45760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff1690811115610bde57610bde611eec565b14610bfb5760405162461bcd60e51b815260040161053e906121d9565b610c0482610ae6565b15610c5d5760405162461bcd60e51b8152602060048201526024808201527f5a65726f6564206164647265737320697320616c726561647920726567697374604482015263195c995960e21b606482015260840161053e565b6001805480820182556000919091526002027fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180546001600160a01b0384166001600160a01b0319909116811782556040519081527fa9a4f3b74b03e48e76814dbc308d3f20104d608c67a42a7ae678d0945daa8e92906020016108f0565b60018181548110610ced57600080fd5b60009182526020909120600290910201546001600160a01b0316905081565b6000546001600160a01b03163314610d365760405162461bcd60e51b815260040161053e906121a4565b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000805b600254811015610a4e57826001600160a01b031660028281548110610dab57610dab61211d565b60009182526020909120600290910201546001600160a01b031603610dd05792915050565b80610dda8161218b565b915050610d84565b6000805b600154811015610e3c5760018181548110610e0357610e0361211d565b6000918252602090912060029091020154610e28906001600160a01b0316318361226d565b915080610e348161218b565b915050610de6565b5090565b6000806000610e4e84610d80565b90506000198103610e9c5760405162461bcd60e51b8152602060048201526018602482015277105b1b1bd8d85d1959081b9bdd081c9959da5cdd195c995960421b604482015260640161053e565b600060028281548110610eb157610eb161211d565b60009182526020918290206040805180820190915260029092020180546001600160a01b03168083526001909101549190920181905290969095509350505050565b60006001600160a01b038216610f3d5760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964206164647265737360881b604482015260640161053e565b60005b600254811015610b9457826001600160a01b031660028281548110610f6757610f6761211d565b60009182526020909120600290910201546001600160a01b031603610f8f5750600192915050565b80610f998161218b565b915050610f40565b600060606000610fb0846109ec565b90506000198103610fd35760405162461bcd60e51b815260040161053e90612280565b600060018281548110610fe857610fe861211d565b6000918252602091829020604080518082018252600290930290910180546001600160a01b0316835260018101805483518187028101870190945280845293949193858301939283018282801561106857602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161104a575b505050505081525050905080600001518160200151935093505050915091565b6000546001600160a01b031633146110b25760405162461bcd60e51b815260040161053e906121a4565b6003805460ff16818111156110c9576110c9611eec565b036111295760405162461bcd60e51b815260206004820152602a60248201527f436f6e74726163742069732066696e616c697a65642c2063616e6e6f742072656044820152697365742076616c75657360b01b606482015260840161053e565b60045443106111865760405162461bcd60e51b8152602060048201526024808201527f526562616c616e636520626c6f636b6e756d62657220616c72656164792070616044820152631cdcd95960e21b606482015260840161053e565b61119260016000611d84565b61119e60026000611da5565b6111aa60056000611dc6565b6003805460ff19169055565b6001806003805460ff16908111156111d0576111d0611eec565b146111ed5760405162461bcd60e51b815260040161053e906121d9565b6111f682610ae6565b6112585760405162461bcd60e51b815260206004820152602d60248201527f7a65726f6564206e6565647320746f206265207265676973746572656420626560448201526c199bdc9948185c1c1c9bdd985b609a1b606482015260840161053e565b813b1515806112d457336001600160a01b038416146112c55760405162461bcd60e51b815260206004820152602360248201527f7a65726f656441646472657373206973206e6f7420746865206d73672e73656e6044820152623232b960e91b606482015260840161053e565b6112cf8333611ac2565b505050565b6000806112e0856119de565b9150915081516000036113355760405162461bcd60e51b815260206004820152601a60248201527f61646d696e206c6973742063616e6e6f7420626520656d707479000000000000604482015260640161053e565b61133f3383611a57565b61138b5760405162461bcd60e51b815260206004820152601b60248201527f6d73672e73656e646572206973206e6f74207468652061646d696e0000000000604482015260640161053e565b6113958533611ac2565b5050505050565b6000546001600160a01b031633146113c65760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff16908111156113e0576113e0611eec565b146113fd5760405162461bcd60e51b815260040161053e906121d9565b6000611408836109ec565b9050600019810361142b5760405162461bcd60e51b815260040161053e90612280565b6001805461143a908290612210565b8154811061144a5761144a61211d565b90600052602060002090600202016001828154811061146b5761146b61211d565b60009182526020909120825460029092020180546001600160a01b0319166001600160a01b03909216919091178155600180830180546114ae9284019190611e00565b5090505060018054806114c3576114c3612223565b60008281526020812060026000199093019283020180546001600160a01b0319168155906114f46001830182611e4c565b505090556040516001600160a01b03841681527f8a654c98d0a7856a8d216c621bb8073316efcaa2b65774d2050c4c1fc7a85a0c906020016108f0565b6000805b600254811015610e3c57600281815481106115525761155261211d565b9060005260206000209060020201600101548261156f919061226d565b91508061157b8161218b565b915050611535565b6000546001600160a01b031633146115ad5760405162461bcd60e51b815260040161053e906121a4565b6002806003805460ff16908111156115c7576115c7611eec565b146115e45760405162461bcd60e51b815260040161053e906121d9565b60056115f083826122fd565b506003805460ff1916811781556040517f8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca64049161162f91600591906123bd565b60405180910390a160045443116116a75760405162461bcd60e51b815260206004820152603660248201527f436f6e74726163742063616e206f6e6c792066696e616c697a6520616674657260448201527520657865637574696e6720726562616c616e63696e6760501b606482015260840161053e565b5050565b6000546001600160a01b031633146116d55760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff16908111156116ef576116ef611eec565b1461170c5760405162461bcd60e51b815260040161053e906121d9565b61171583610ef3565b156117725760405162461bcd60e51b815260206004820152602760248201527f416c6c6f6361746564206164647265737320697320616c7265616479207265676044820152661a5cdd195c995960ca1b606482015260840161053e565b816000036117c25760405162461bcd60e51b815260206004820152601960248201527f416d6f756e742063616e6e6f742062652073657420746f203000000000000000604482015260640161053e565b6040805180820182526001600160a01b038581168083526020808401878152600280546001810182556000829052865191027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace81018054929096166001600160a01b031990921691909117909455517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf90930192909255835190815290810185905290917fab5b2126f71ee7e0b39eadc53fb5d08a8f6c68dc61795fa05ed7d176cd2665ed91015b60405180910390a150505050565b6000546001600160a01b031633146118c25760405162461bcd60e51b815260040161053e906121a4565b61062a81611cc4565b6000546001600160a01b031633146118f55760405162461bcd60e51b815260040161053e906121a4565b6001806003805460ff169081111561190f5761190f611eec565b1461192c5760405162461bcd60e51b815260040161053e906121d9565b611934610de2565b61193c611531565b106119c25760405162461bcd60e51b815260206004820152604a60248201527f747265617375727920616d6f756e742073686f756c64206265206c657373207460448201527f68616e207468652073756d206f6620616c6c207a65726f656420616464726573606482015269732062616c616e63657360b01b608482015260a40161053e565b6119ca610440565b600380546002919060ff19166001836109a5565b6060600080839050806001600160a01b0316631865c57d6040518163ffffffff1660e01b8152600401600060405180830381865afa158015611a24573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611a4c9190810190612452565b909590945092505050565b8051600090815b81811015611ab557838181518110611a7857611a7861211d565b60200260200101516001600160a01b0316856001600160a01b031603611aa357600192505050611abc565b80611aad8161218b565b915050611a5e565b5060009150505b92915050565b6000611acd836109ec565b90506000198103611af05760405162461bcd60e51b815260040161053e90612280565b600060018281548110611b0557611b0561211d565b9060005260206000209060020201600101805480602002602001604051908101604052809291908181526020018280548015611b6a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611b4c575b5050505050905060005b8151811015611bfc57836001600160a01b0316828281518110611b9957611b9961211d565b60200260200101516001600160a01b031603611bea5760405162461bcd60e51b815260206004820152601060248201526f105b1c9958591e48185c1c1c9bdd995960821b604482015260640161053e565b80611bf48161218b565b915050611b74565b5060018281548110611c1057611c1061211d565b600091825260208083206001600290930201820180548084018255908452922090910180546001600160a01b0386166001600160a01b031990911617905580547f80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f9091869186919086908110611c8757611c8761211d565b600091825260209182902060016002909202010154604080516001600160a01b03958616815294909316918401919091529082015260600161188a565b6001600160a01b038116611d295760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161053e565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b508054600082556002029060005260206000209081019061062a9190611e6a565b508054600082556002029060005260206000209081019061062a9190611e98565b508054611dd290612239565b6000825580601f10611de2575050565b601f01602090049060005260206000209081019061062a9190611ebe565b828054828255906000526020600020908101928215611e405760005260206000209182015b82811115611e40578254825591600101919060010190611e25565b50610e3c929150611ebe565b508054600082559060005260206000209081019061062a9190611ebe565b80821115610e3c5780546001600160a01b03191681556000611e8f6001830182611e4c565b50600201611e6a565b5b80821115610e3c5780546001600160a01b031916815560006001820155600201611e99565b5b80821115610e3c5760008155600101611ebf565b600060208284031215611ee557600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b60048110611f2057634e487b7160e01b600052602160045260246000fd5b9052565b60208101611abc8284611f02565b6001600160a01b038116811461062a57600080fd5b600060208284031215611f5957600080fd5b8135611f6481611f32565b9392505050565b600060208083528351808285015260005b81811015611f9857858101830151858201604001528201611f7c565b506000604082860101526040601f19601f8301168501019250505092915050565b6001600160a01b038381168252604060208084018290528451918401829052600092858201929091906060860190855b81811015612007578551851683529483019491830191600101611fe9565b509098975050505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561205457612054612015565b604052919050565b6000602080838503121561206f57600080fd5b823567ffffffffffffffff8082111561208757600080fd5b818501915085601f83011261209b57600080fd5b8135818111156120ad576120ad612015565b6120bf601f8201601f1916850161202b565b915080825286848285010111156120d557600080fd5b8084840185840137600090820190930192909252509392505050565b6000806040838503121561210457600080fd5b823561210f81611f32565b946020939093013593505050565b634e487b7160e01b600052603260045260246000fd5b60208082526022908201527f6d696e2072657175697265642061646d696e732073686f756c6420617070726f604082015261766560f01b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b60006001820161219d5761219d612175565b5060010190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252601c908201527f4e6f7420696e207468652064657369676e617465642073746174757300000000604082015260600190565b81810381811115611abc57611abc612175565b634e487b7160e01b600052603160045260246000fd5b600181811c9082168061224d57607f821691505b602082108103610b9457634e487b7160e01b600052602260045260246000fd5b80820180821115611abc57611abc612175565b60208082526015908201527416995c9bd959081b9bdd081c9959da5cdd195c9959605a1b604082015260600190565b601f8211156112cf57600081815260208120601f850160051c810160208610156122d65750805b601f850160051c820191505b818110156122f5578281556001016122e2565b505050505050565b815167ffffffffffffffff81111561231757612317612015565b61232b816123258454612239565b846122af565b602080601f83116001811461236057600084156123485750858301515b600019600386901b1c1916600185901b1785556122f5565b600085815260208120601f198616915b8281101561238f57888601518255948401946001909101908401612370565b50858210156123ad5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6040815260008084546123cf81612239565b80604086015260606001808416600081146123f1576001811461240b5761243c565b60ff1985168884015283151560051b88018301955061243c565b8960005260208060002060005b868110156124335781548b8201870152908401908201612418565b8a018501975050505b505050505080915050611f646020830184611f02565b6000806040838503121561246557600080fd5b825167ffffffffffffffff8082111561247d57600080fd5b818501915085601f83011261249157600080fd5b81516020828211156124a5576124a5612015565b8160051b92506124b681840161202b565b82815292840181019281810190898511156124d057600080fd5b948201945b848610156124fa57855193506124ea84611f32565b83825294820194908201906124d5565b9790910151969896975050505050505056fea2646970667358221220d937546b8b082e09d72e451ac1c9bc7ee3016e2906fdce592ba18709cb53040764736f6c63430008130033",
}

// TreasuryRebalanceV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use TreasuryRebalanceV2MetaData.ABI instead.
var TreasuryRebalanceV2ABI = TreasuryRebalanceV2MetaData.ABI

// TreasuryRebalanceV2BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const TreasuryRebalanceV2BinRuntime = `608060405234801561001057600080fd5b50600436106101e55760003560e01c80638f32d59b1161010f578063db27b50b116100a2578063ecd8677811610071578063ecd867781461040a578063ed3555291461041d578063f2fde38b14610425578063faaf9ca61461043857600080fd5b8063db27b50b146103c8578063e20fcf00146103db578063e2384cb3146103e3578063ea6d4a9b146103f757600080fd5b8063bd786f57116100de578063bd786f5714610379578063cea1c3381461038c578063d826f88f146103ad578063daea85c5146103b557600080fd5b80638f32d59b146103435780639ab29b70146103565780639dc954ba1461035e5780639e59eb141461036657600080fd5b8063518592da1161018757806362aa3e911161015657806362aa3e91146102ec578063715018a6146103175780637bfaf7b71461031f5780638da5cb5b1461033257600080fd5b8063518592da1461028e57806358c3b870146102a15780635f8798c0146102b65780635f9b0df7146102d957600080fd5b806327704cb5116101c357806327704cb51461022a578063343e2c851461023d578063484090961461026f57806349a3fb451461027757600080fd5b80630287d126146101ea5780631804692f146101f4578063200d2ed214610207575b600080fd5b6101f2610440565b005b6101f2610202366004611ed3565b61062d565b6003546102149060ff1681565b6040516102219190611f24565b60405180910390f35b6101f2610238366004611f47565b610748565b61025061024b366004611ed3565b6108fd565b604080516001600160a01b039093168352602083019190915201610221565b6101f2610935565b61028060045481565b604051908152602001610221565b61028061029c366004611f47565b6109ec565b6102a9610a58565b6040516102219190611f6b565b6102c96102c4366004611f47565b610ae6565b6040519015158152602001610221565b6101f26102e7366004611f47565b610b9a565b6102ff6102fa366004611ed3565b610cdd565b6040516001600160a01b039091168152602001610221565b6101f2610d0c565b61028061032d366004611f47565b610d80565b6000546001600160a01b03166102ff565b6000546001600160a01b031633146102c9565b610280610de2565b600154610280565b610250610374366004611f47565b610e40565b6102c9610387366004611f47565b610ef3565b61039f61039a366004611f47565b610fa1565b604051610221929190611fb9565b6101f2611088565b6101f26103c3366004611f47565b6111b6565b6101f26103d6366004611f47565b61139c565b610280611531565b6102c96103f1366004611f47565b3b151590565b6101f261040536600461205c565b611583565b6101f26104183660046120f1565b6116ab565b600254610280565b6101f2610433366004611f47565b611898565b6101f26118cb565b60005b60015481101561062a576000600182815481106104625761046261211d565b6000918252602091829020604080518082018252600290930290910180546001600160a01b031683526001810180548351818702810187019094528084529394919385830193928301828280156104e257602002820191906000526020600020905b81546001600160a01b031681526001909101906020018083116104c4575b505050505081525050905060006104fd82600001513b151590565b905080156105cb5760008061051584600001516119de565b915091508084602001515110156105475760405162461bcd60e51b815260040161053e90612133565b60405180910390fd5b60208401516000805b82518110156105a15761057c83828151811061056e5761056e61211d565b602002602001015186611a57565b1561058f578161058b8161218b565b9250505b806105998161218b565b915050610550565b50828110156105c25760405162461bcd60e51b815260040161053e90612133565b50505050610615565b8160200151516001146106155760405162461bcd60e51b8152602060048201526012602482015271454f412073686f756c6420617070726f766560701b604482015260640161053e565b505080806106229061218b565b915050610443565b50565b6000546001600160a01b031633146106575760405162461bcd60e51b815260040161053e906121a4565b60045443106106ce5760405162461bcd60e51b815260206004820152603e60248201527f63757272656e7420626c6f636b2073686f756c646e277420626520706173742060448201527f7468652063757272656e746c792073657420626c6f636b206e756d6265720000606482015260840161053e565b8043106107435760405162461bcd60e51b815260206004820152603a60248201527f726562616c616e636520626c6f636b4e756d6265722073686f756c642062652060448201527f67726561746572207468616e2063757272656e7420626c6f636b000000000000606482015260840161053e565b600455565b6000546001600160a01b031633146107725760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff169081111561078c5761078c611eec565b146107a95760405162461bcd60e51b815260040161053e906121d9565b60006107b483610d80565b905060001981036108025760405162461bcd60e51b8152602060048201526018602482015277105b1b1bd8d85d1959081b9bdd081c9959da5cdd195c995960421b604482015260640161053e565b6002805461081290600190612210565b815481106108225761082261211d565b9060005260206000209060020201600282815481106108435761084361211d565b600091825260209091208254600292830290910180546001600160a01b0319166001600160a01b0390921691909117815560019283015492019190915580548061088f5761088f612223565b600082815260208082206002600019949094019384020180546001600160a01b03191681556001019190915591556040516001600160a01b03851681527ff8f67464bea52432645435be9c46c427173a75aefaa1001272e08a4b8572f06e91015b60405180910390a1505050565b6002818154811061090d57600080fd5b6000918252602090912060029091020180546001909101546001600160a01b03909116915082565b6000546001600160a01b0316331461095f5760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff169081111561097957610979611eec565b146109965760405162461bcd60e51b815260040161053e906121d9565b600380546001919060ff191682805b02179055506003546040517fafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e916109e19160ff90911690611f24565b60405180910390a150565b6000805b600154811015610a4e57826001600160a01b031660018281548110610a1757610a1761211d565b60009182526020909120600290910201546001600160a01b031603610a3c5792915050565b80610a468161218b565b9150506109f0565b5060001992915050565b60058054610a6590612239565b80601f0160208091040260200160405190810160405280929190818152602001828054610a9190612239565b8015610ade5780601f10610ab357610100808354040283529160200191610ade565b820191906000526020600020905b815481529060010190602001808311610ac157829003601f168201915b505050505081565b60006001600160a01b038216610b305760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964206164647265737360881b604482015260640161053e565b60005b600154811015610b9457826001600160a01b031660018281548110610b5a57610b5a61211d565b60009182526020909120600290910201546001600160a01b031603610b825750600192915050565b80610b8c8161218b565b915050610b33565b50919050565b6000546001600160a01b03163314610bc45760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff1690811115610bde57610bde611eec565b14610bfb5760405162461bcd60e51b815260040161053e906121d9565b610c0482610ae6565b15610c5d5760405162461bcd60e51b8152602060048201526024808201527f5a65726f6564206164647265737320697320616c726561647920726567697374604482015263195c995960e21b606482015260840161053e565b6001805480820182556000919091526002027fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180546001600160a01b0384166001600160a01b0319909116811782556040519081527fa9a4f3b74b03e48e76814dbc308d3f20104d608c67a42a7ae678d0945daa8e92906020016108f0565b60018181548110610ced57600080fd5b60009182526020909120600290910201546001600160a01b0316905081565b6000546001600160a01b03163314610d365760405162461bcd60e51b815260040161053e906121a4565b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000805b600254811015610a4e57826001600160a01b031660028281548110610dab57610dab61211d565b60009182526020909120600290910201546001600160a01b031603610dd05792915050565b80610dda8161218b565b915050610d84565b6000805b600154811015610e3c5760018181548110610e0357610e0361211d565b6000918252602090912060029091020154610e28906001600160a01b0316318361226d565b915080610e348161218b565b915050610de6565b5090565b6000806000610e4e84610d80565b90506000198103610e9c5760405162461bcd60e51b8152602060048201526018602482015277105b1b1bd8d85d1959081b9bdd081c9959da5cdd195c995960421b604482015260640161053e565b600060028281548110610eb157610eb161211d565b60009182526020918290206040805180820190915260029092020180546001600160a01b03168083526001909101549190920181905290969095509350505050565b60006001600160a01b038216610f3d5760405162461bcd60e51b815260206004820152600f60248201526e496e76616c6964206164647265737360881b604482015260640161053e565b60005b600254811015610b9457826001600160a01b031660028281548110610f6757610f6761211d565b60009182526020909120600290910201546001600160a01b031603610f8f5750600192915050565b80610f998161218b565b915050610f40565b600060606000610fb0846109ec565b90506000198103610fd35760405162461bcd60e51b815260040161053e90612280565b600060018281548110610fe857610fe861211d565b6000918252602091829020604080518082018252600290930290910180546001600160a01b0316835260018101805483518187028101870190945280845293949193858301939283018282801561106857602002820191906000526020600020905b81546001600160a01b0316815260019091019060200180831161104a575b505050505081525050905080600001518160200151935093505050915091565b6000546001600160a01b031633146110b25760405162461bcd60e51b815260040161053e906121a4565b6003805460ff16818111156110c9576110c9611eec565b036111295760405162461bcd60e51b815260206004820152602a60248201527f436f6e74726163742069732066696e616c697a65642c2063616e6e6f742072656044820152697365742076616c75657360b01b606482015260840161053e565b60045443106111865760405162461bcd60e51b8152602060048201526024808201527f526562616c616e636520626c6f636b6e756d62657220616c72656164792070616044820152631cdcd95960e21b606482015260840161053e565b61119260016000611d84565b61119e60026000611da5565b6111aa60056000611dc6565b6003805460ff19169055565b6001806003805460ff16908111156111d0576111d0611eec565b146111ed5760405162461bcd60e51b815260040161053e906121d9565b6111f682610ae6565b6112585760405162461bcd60e51b815260206004820152602d60248201527f7a65726f6564206e6565647320746f206265207265676973746572656420626560448201526c199bdc9948185c1c1c9bdd985b609a1b606482015260840161053e565b813b1515806112d457336001600160a01b038416146112c55760405162461bcd60e51b815260206004820152602360248201527f7a65726f656441646472657373206973206e6f7420746865206d73672e73656e6044820152623232b960e91b606482015260840161053e565b6112cf8333611ac2565b505050565b6000806112e0856119de565b9150915081516000036113355760405162461bcd60e51b815260206004820152601a60248201527f61646d696e206c6973742063616e6e6f7420626520656d707479000000000000604482015260640161053e565b61133f3383611a57565b61138b5760405162461bcd60e51b815260206004820152601b60248201527f6d73672e73656e646572206973206e6f74207468652061646d696e0000000000604482015260640161053e565b6113958533611ac2565b5050505050565b6000546001600160a01b031633146113c65760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff16908111156113e0576113e0611eec565b146113fd5760405162461bcd60e51b815260040161053e906121d9565b6000611408836109ec565b9050600019810361142b5760405162461bcd60e51b815260040161053e90612280565b6001805461143a908290612210565b8154811061144a5761144a61211d565b90600052602060002090600202016001828154811061146b5761146b61211d565b60009182526020909120825460029092020180546001600160a01b0319166001600160a01b03909216919091178155600180830180546114ae9284019190611e00565b5090505060018054806114c3576114c3612223565b60008281526020812060026000199093019283020180546001600160a01b0319168155906114f46001830182611e4c565b505090556040516001600160a01b03841681527f8a654c98d0a7856a8d216c621bb8073316efcaa2b65774d2050c4c1fc7a85a0c906020016108f0565b6000805b600254811015610e3c57600281815481106115525761155261211d565b9060005260206000209060020201600101548261156f919061226d565b91508061157b8161218b565b915050611535565b6000546001600160a01b031633146115ad5760405162461bcd60e51b815260040161053e906121a4565b6002806003805460ff16908111156115c7576115c7611eec565b146115e45760405162461bcd60e51b815260040161053e906121d9565b60056115f083826122fd565b506003805460ff1916811781556040517f8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca64049161162f91600591906123bd565b60405180910390a160045443116116a75760405162461bcd60e51b815260206004820152603660248201527f436f6e74726163742063616e206f6e6c792066696e616c697a6520616674657260448201527520657865637574696e6720726562616c616e63696e6760501b606482015260840161053e565b5050565b6000546001600160a01b031633146116d55760405162461bcd60e51b815260040161053e906121a4565b6000806003805460ff16908111156116ef576116ef611eec565b1461170c5760405162461bcd60e51b815260040161053e906121d9565b61171583610ef3565b156117725760405162461bcd60e51b815260206004820152602760248201527f416c6c6f6361746564206164647265737320697320616c7265616479207265676044820152661a5cdd195c995960ca1b606482015260840161053e565b816000036117c25760405162461bcd60e51b815260206004820152601960248201527f416d6f756e742063616e6e6f742062652073657420746f203000000000000000604482015260640161053e565b6040805180820182526001600160a01b038581168083526020808401878152600280546001810182556000829052865191027f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace81018054929096166001600160a01b031990921691909117909455517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf90930192909255835190815290810185905290917fab5b2126f71ee7e0b39eadc53fb5d08a8f6c68dc61795fa05ed7d176cd2665ed91015b60405180910390a150505050565b6000546001600160a01b031633146118c25760405162461bcd60e51b815260040161053e906121a4565b61062a81611cc4565b6000546001600160a01b031633146118f55760405162461bcd60e51b815260040161053e906121a4565b6001806003805460ff169081111561190f5761190f611eec565b1461192c5760405162461bcd60e51b815260040161053e906121d9565b611934610de2565b61193c611531565b106119c25760405162461bcd60e51b815260206004820152604a60248201527f747265617375727920616d6f756e742073686f756c64206265206c657373207460448201527f68616e207468652073756d206f6620616c6c207a65726f656420616464726573606482015269732062616c616e63657360b01b608482015260a40161053e565b6119ca610440565b600380546002919060ff19166001836109a5565b6060600080839050806001600160a01b0316631865c57d6040518163ffffffff1660e01b8152600401600060405180830381865afa158015611a24573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f19168201604052611a4c9190810190612452565b909590945092505050565b8051600090815b81811015611ab557838181518110611a7857611a7861211d565b60200260200101516001600160a01b0316856001600160a01b031603611aa357600192505050611abc565b80611aad8161218b565b915050611a5e565b5060009150505b92915050565b6000611acd836109ec565b90506000198103611af05760405162461bcd60e51b815260040161053e90612280565b600060018281548110611b0557611b0561211d565b9060005260206000209060020201600101805480602002602001604051908101604052809291908181526020018280548015611b6a57602002820191906000526020600020905b81546001600160a01b03168152600190910190602001808311611b4c575b5050505050905060005b8151811015611bfc57836001600160a01b0316828281518110611b9957611b9961211d565b60200260200101516001600160a01b031603611bea5760405162461bcd60e51b815260206004820152601060248201526f105b1c9958591e48185c1c1c9bdd995960821b604482015260640161053e565b80611bf48161218b565b915050611b74565b5060018281548110611c1057611c1061211d565b600091825260208083206001600290930201820180548084018255908452922090910180546001600160a01b0386166001600160a01b031990911617905580547f80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f9091869186919086908110611c8757611c8761211d565b600091825260209182902060016002909202010154604080516001600160a01b03958616815294909316918401919091529082015260600161188a565b6001600160a01b038116611d295760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161053e565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b508054600082556002029060005260206000209081019061062a9190611e6a565b508054600082556002029060005260206000209081019061062a9190611e98565b508054611dd290612239565b6000825580601f10611de2575050565b601f01602090049060005260206000209081019061062a9190611ebe565b828054828255906000526020600020908101928215611e405760005260206000209182015b82811115611e40578254825591600101919060010190611e25565b50610e3c929150611ebe565b508054600082559060005260206000209081019061062a9190611ebe565b80821115610e3c5780546001600160a01b03191681556000611e8f6001830182611e4c565b50600201611e6a565b5b80821115610e3c5780546001600160a01b031916815560006001820155600201611e99565b5b80821115610e3c5760008155600101611ebf565b600060208284031215611ee557600080fd5b5035919050565b634e487b7160e01b600052602160045260246000fd5b60048110611f2057634e487b7160e01b600052602160045260246000fd5b9052565b60208101611abc8284611f02565b6001600160a01b038116811461062a57600080fd5b600060208284031215611f5957600080fd5b8135611f6481611f32565b9392505050565b600060208083528351808285015260005b81811015611f9857858101830151858201604001528201611f7c565b506000604082860101526040601f19601f8301168501019250505092915050565b6001600160a01b038381168252604060208084018290528451918401829052600092858201929091906060860190855b81811015612007578551851683529483019491830191600101611fe9565b509098975050505050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f1916810167ffffffffffffffff8111828210171561205457612054612015565b604052919050565b6000602080838503121561206f57600080fd5b823567ffffffffffffffff8082111561208757600080fd5b818501915085601f83011261209b57600080fd5b8135818111156120ad576120ad612015565b6120bf601f8201601f1916850161202b565b915080825286848285010111156120d557600080fd5b8084840185840137600090820190930192909252509392505050565b6000806040838503121561210457600080fd5b823561210f81611f32565b946020939093013593505050565b634e487b7160e01b600052603260045260246000fd5b60208082526022908201527f6d696e2072657175697265642061646d696e732073686f756c6420617070726f604082015261766560f01b606082015260800190565b634e487b7160e01b600052601160045260246000fd5b60006001820161219d5761219d612175565b5060010190565b6020808252818101527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604082015260600190565b6020808252601c908201527f4e6f7420696e207468652064657369676e617465642073746174757300000000604082015260600190565b81810381811115611abc57611abc612175565b634e487b7160e01b600052603160045260246000fd5b600181811c9082168061224d57607f821691505b602082108103610b9457634e487b7160e01b600052602260045260246000fd5b80820180821115611abc57611abc612175565b60208082526015908201527416995c9bd959081b9bdd081c9959da5cdd195c9959605a1b604082015260600190565b601f8211156112cf57600081815260208120601f850160051c810160208610156122d65750805b601f850160051c820191505b818110156122f5578281556001016122e2565b505050505050565b815167ffffffffffffffff81111561231757612317612015565b61232b816123258454612239565b846122af565b602080601f83116001811461236057600084156123485750858301515b600019600386901b1c1916600185901b1785556122f5565b600085815260208120601f198616915b8281101561238f57888601518255948401946001909101908401612370565b50858210156123ad5787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b6040815260008084546123cf81612239565b80604086015260606001808416600081146123f1576001811461240b5761243c565b60ff1985168884015283151560051b88018301955061243c565b8960005260208060002060005b868110156124335781548b8201870152908401908201612418565b8a018501975050505b505050505080915050611f646020830184611f02565b6000806040838503121561246557600080fd5b825167ffffffffffffffff8082111561247d57600080fd5b818501915085601f83011261249157600080fd5b81516020828211156124a5576124a5612015565b8160051b92506124b681840161202b565b82815292840181019281810190898511156124d057600080fd5b948201945b848610156124fa57855193506124ea84611f32565b83825294820194908201906124d5565b9790910151969896975050505050505056fea2646970667358221220d937546b8b082e09d72e451ac1c9bc7ee3016e2906fdce592ba18709cb53040764736f6c63430008130033`

// TreasuryRebalanceV2FuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use TreasuryRebalanceV2MetaData.Sigs instead.
var TreasuryRebalanceV2FuncSigs = TreasuryRebalanceV2MetaData.Sigs

// TreasuryRebalanceV2Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TreasuryRebalanceV2MetaData.Bin instead.
var TreasuryRebalanceV2Bin = TreasuryRebalanceV2MetaData.Bin

// DeployTreasuryRebalanceV2 deploys a new Klaytn contract, binding an instance of TreasuryRebalanceV2 to it.
func DeployTreasuryRebalanceV2(auth *bind.TransactOpts, backend bind.ContractBackend, _rebalanceBlockNumber *big.Int) (common.Address, *types.Transaction, *TreasuryRebalanceV2, error) {
	parsed, err := TreasuryRebalanceV2MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TreasuryRebalanceV2Bin), backend, _rebalanceBlockNumber)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TreasuryRebalanceV2{TreasuryRebalanceV2Caller: TreasuryRebalanceV2Caller{contract: contract}, TreasuryRebalanceV2Transactor: TreasuryRebalanceV2Transactor{contract: contract}, TreasuryRebalanceV2Filterer: TreasuryRebalanceV2Filterer{contract: contract}}, nil
}

// TreasuryRebalanceV2 is an auto generated Go binding around a Klaytn contract.
type TreasuryRebalanceV2 struct {
	TreasuryRebalanceV2Caller     // Read-only binding to the contract
	TreasuryRebalanceV2Transactor // Write-only binding to the contract
	TreasuryRebalanceV2Filterer   // Log filterer for contract events
}

// TreasuryRebalanceV2Caller is an auto generated read-only Go binding around a Klaytn contract.
type TreasuryRebalanceV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreasuryRebalanceV2Transactor is an auto generated write-only Go binding around a Klaytn contract.
type TreasuryRebalanceV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreasuryRebalanceV2Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type TreasuryRebalanceV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TreasuryRebalanceV2Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type TreasuryRebalanceV2Session struct {
	Contract     *TreasuryRebalanceV2 // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TreasuryRebalanceV2CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type TreasuryRebalanceV2CallerSession struct {
	Contract *TreasuryRebalanceV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// TreasuryRebalanceV2TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type TreasuryRebalanceV2TransactorSession struct {
	Contract     *TreasuryRebalanceV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// TreasuryRebalanceV2Raw is an auto generated low-level Go binding around a Klaytn contract.
type TreasuryRebalanceV2Raw struct {
	Contract *TreasuryRebalanceV2 // Generic contract binding to access the raw methods on
}

// TreasuryRebalanceV2CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type TreasuryRebalanceV2CallerRaw struct {
	Contract *TreasuryRebalanceV2Caller // Generic read-only contract binding to access the raw methods on
}

// TreasuryRebalanceV2TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type TreasuryRebalanceV2TransactorRaw struct {
	Contract *TreasuryRebalanceV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewTreasuryRebalanceV2 creates a new instance of TreasuryRebalanceV2, bound to a specific deployed contract.
func NewTreasuryRebalanceV2(address common.Address, backend bind.ContractBackend) (*TreasuryRebalanceV2, error) {
	contract, err := bindTreasuryRebalanceV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2{TreasuryRebalanceV2Caller: TreasuryRebalanceV2Caller{contract: contract}, TreasuryRebalanceV2Transactor: TreasuryRebalanceV2Transactor{contract: contract}, TreasuryRebalanceV2Filterer: TreasuryRebalanceV2Filterer{contract: contract}}, nil
}

// NewTreasuryRebalanceV2Caller creates a new read-only instance of TreasuryRebalanceV2, bound to a specific deployed contract.
func NewTreasuryRebalanceV2Caller(address common.Address, caller bind.ContractCaller) (*TreasuryRebalanceV2Caller, error) {
	contract, err := bindTreasuryRebalanceV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2Caller{contract: contract}, nil
}

// NewTreasuryRebalanceV2Transactor creates a new write-only instance of TreasuryRebalanceV2, bound to a specific deployed contract.
func NewTreasuryRebalanceV2Transactor(address common.Address, transactor bind.ContractTransactor) (*TreasuryRebalanceV2Transactor, error) {
	contract, err := bindTreasuryRebalanceV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2Transactor{contract: contract}, nil
}

// NewTreasuryRebalanceV2Filterer creates a new log filterer instance of TreasuryRebalanceV2, bound to a specific deployed contract.
func NewTreasuryRebalanceV2Filterer(address common.Address, filterer bind.ContractFilterer) (*TreasuryRebalanceV2Filterer, error) {
	contract, err := bindTreasuryRebalanceV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2Filterer{contract: contract}, nil
}

// bindTreasuryRebalanceV2 binds a generic wrapper to an already deployed contract.
func bindTreasuryRebalanceV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TreasuryRebalanceV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TreasuryRebalanceV2.Contract.TreasuryRebalanceV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.TreasuryRebalanceV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.TreasuryRebalanceV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TreasuryRebalanceV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.contract.Transact(opts, method, params...)
}

// AllocatedExists is a free data retrieval call binding the contract method 0xbd786f57.
//
// Solidity: function allocatedExists(address _allocatedAddress) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) AllocatedExists(opts *bind.CallOpts, _allocatedAddress common.Address) (bool, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "allocatedExists", _allocatedAddress)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// AllocatedExists is a free data retrieval call binding the contract method 0xbd786f57.
//
// Solidity: function allocatedExists(address _allocatedAddress) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) AllocatedExists(_allocatedAddress common.Address) (bool, error) {
	return _TreasuryRebalanceV2.Contract.AllocatedExists(&_TreasuryRebalanceV2.CallOpts, _allocatedAddress)
}

// AllocatedExists is a free data retrieval call binding the contract method 0xbd786f57.
//
// Solidity: function allocatedExists(address _allocatedAddress) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) AllocatedExists(_allocatedAddress common.Address) (bool, error) {
	return _TreasuryRebalanceV2.Contract.AllocatedExists(&_TreasuryRebalanceV2.CallOpts, _allocatedAddress)
}

// Allocateds is a free data retrieval call binding the contract method 0x343e2c85.
//
// Solidity: function allocateds(uint256 ) view returns(address addr, uint256 amount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) Allocateds(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Addr   common.Address
	Amount *big.Int
}, error,
) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "allocateds", arg0)

	outstruct := new(struct {
		Addr   common.Address
		Amount *big.Int
	})

	outstruct.Addr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Amount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	return *outstruct, err
}

// Allocateds is a free data retrieval call binding the contract method 0x343e2c85.
//
// Solidity: function allocateds(uint256 ) view returns(address addr, uint256 amount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) Allocateds(arg0 *big.Int) (struct {
	Addr   common.Address
	Amount *big.Int
}, error,
) {
	return _TreasuryRebalanceV2.Contract.Allocateds(&_TreasuryRebalanceV2.CallOpts, arg0)
}

// Allocateds is a free data retrieval call binding the contract method 0x343e2c85.
//
// Solidity: function allocateds(uint256 ) view returns(address addr, uint256 amount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) Allocateds(arg0 *big.Int) (struct {
	Addr   common.Address
	Amount *big.Int
}, error,
) {
	return _TreasuryRebalanceV2.Contract.Allocateds(&_TreasuryRebalanceV2.CallOpts, arg0)
}

// CheckZeroedsApproved is a free data retrieval call binding the contract method 0x0287d126.
//
// Solidity: function checkZeroedsApproved() view returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) CheckZeroedsApproved(opts *bind.CallOpts) error {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "checkZeroedsApproved")
	if err != nil {
		return err
	}

	return err
}

// CheckZeroedsApproved is a free data retrieval call binding the contract method 0x0287d126.
//
// Solidity: function checkZeroedsApproved() view returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) CheckZeroedsApproved() error {
	return _TreasuryRebalanceV2.Contract.CheckZeroedsApproved(&_TreasuryRebalanceV2.CallOpts)
}

// CheckZeroedsApproved is a free data retrieval call binding the contract method 0x0287d126.
//
// Solidity: function checkZeroedsApproved() view returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) CheckZeroedsApproved() error {
	return _TreasuryRebalanceV2.Contract.CheckZeroedsApproved(&_TreasuryRebalanceV2.CallOpts)
}

// GetAllocated is a free data retrieval call binding the contract method 0x9e59eb14.
//
// Solidity: function getAllocated(address _allocatedAddress) view returns(address, uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) GetAllocated(opts *bind.CallOpts, _allocatedAddress common.Address) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "getAllocated", _allocatedAddress)
	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err
}

// GetAllocated is a free data retrieval call binding the contract method 0x9e59eb14.
//
// Solidity: function getAllocated(address _allocatedAddress) view returns(address, uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) GetAllocated(_allocatedAddress common.Address) (common.Address, *big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetAllocated(&_TreasuryRebalanceV2.CallOpts, _allocatedAddress)
}

// GetAllocated is a free data retrieval call binding the contract method 0x9e59eb14.
//
// Solidity: function getAllocated(address _allocatedAddress) view returns(address, uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) GetAllocated(_allocatedAddress common.Address) (common.Address, *big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetAllocated(&_TreasuryRebalanceV2.CallOpts, _allocatedAddress)
}

// GetAllocatedCount is a free data retrieval call binding the contract method 0xed355529.
//
// Solidity: function getAllocatedCount() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) GetAllocatedCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "getAllocatedCount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetAllocatedCount is a free data retrieval call binding the contract method 0xed355529.
//
// Solidity: function getAllocatedCount() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) GetAllocatedCount() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetAllocatedCount(&_TreasuryRebalanceV2.CallOpts)
}

// GetAllocatedCount is a free data retrieval call binding the contract method 0xed355529.
//
// Solidity: function getAllocatedCount() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) GetAllocatedCount() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetAllocatedCount(&_TreasuryRebalanceV2.CallOpts)
}

// GetAllocatedIndex is a free data retrieval call binding the contract method 0x7bfaf7b7.
//
// Solidity: function getAllocatedIndex(address _allocatedAddress) view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) GetAllocatedIndex(opts *bind.CallOpts, _allocatedAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "getAllocatedIndex", _allocatedAddress)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetAllocatedIndex is a free data retrieval call binding the contract method 0x7bfaf7b7.
//
// Solidity: function getAllocatedIndex(address _allocatedAddress) view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) GetAllocatedIndex(_allocatedAddress common.Address) (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetAllocatedIndex(&_TreasuryRebalanceV2.CallOpts, _allocatedAddress)
}

// GetAllocatedIndex is a free data retrieval call binding the contract method 0x7bfaf7b7.
//
// Solidity: function getAllocatedIndex(address _allocatedAddress) view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) GetAllocatedIndex(_allocatedAddress common.Address) (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetAllocatedIndex(&_TreasuryRebalanceV2.CallOpts, _allocatedAddress)
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) GetTreasuryAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "getTreasuryAmount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) GetTreasuryAmount() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetTreasuryAmount(&_TreasuryRebalanceV2.CallOpts)
}

// GetTreasuryAmount is a free data retrieval call binding the contract method 0xe20fcf00.
//
// Solidity: function getTreasuryAmount() view returns(uint256 treasuryAmount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) GetTreasuryAmount() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetTreasuryAmount(&_TreasuryRebalanceV2.CallOpts)
}

// GetZeroed is a free data retrieval call binding the contract method 0xcea1c338.
//
// Solidity: function getZeroed(address _zeroedAddress) view returns(address, address[])
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) GetZeroed(opts *bind.CallOpts, _zeroedAddress common.Address) (common.Address, []common.Address, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "getZeroed", _zeroedAddress)
	if err != nil {
		return *new(common.Address), *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new([]common.Address)).(*[]common.Address)

	return out0, out1, err
}

// GetZeroed is a free data retrieval call binding the contract method 0xcea1c338.
//
// Solidity: function getZeroed(address _zeroedAddress) view returns(address, address[])
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) GetZeroed(_zeroedAddress common.Address) (common.Address, []common.Address, error) {
	return _TreasuryRebalanceV2.Contract.GetZeroed(&_TreasuryRebalanceV2.CallOpts, _zeroedAddress)
}

// GetZeroed is a free data retrieval call binding the contract method 0xcea1c338.
//
// Solidity: function getZeroed(address _zeroedAddress) view returns(address, address[])
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) GetZeroed(_zeroedAddress common.Address) (common.Address, []common.Address, error) {
	return _TreasuryRebalanceV2.Contract.GetZeroed(&_TreasuryRebalanceV2.CallOpts, _zeroedAddress)
}

// GetZeroedCount is a free data retrieval call binding the contract method 0x9dc954ba.
//
// Solidity: function getZeroedCount() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) GetZeroedCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "getZeroedCount")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetZeroedCount is a free data retrieval call binding the contract method 0x9dc954ba.
//
// Solidity: function getZeroedCount() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) GetZeroedCount() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetZeroedCount(&_TreasuryRebalanceV2.CallOpts)
}

// GetZeroedCount is a free data retrieval call binding the contract method 0x9dc954ba.
//
// Solidity: function getZeroedCount() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) GetZeroedCount() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetZeroedCount(&_TreasuryRebalanceV2.CallOpts)
}

// GetZeroedIndex is a free data retrieval call binding the contract method 0x518592da.
//
// Solidity: function getZeroedIndex(address _zeroedAddress) view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) GetZeroedIndex(opts *bind.CallOpts, _zeroedAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "getZeroedIndex", _zeroedAddress)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetZeroedIndex is a free data retrieval call binding the contract method 0x518592da.
//
// Solidity: function getZeroedIndex(address _zeroedAddress) view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) GetZeroedIndex(_zeroedAddress common.Address) (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetZeroedIndex(&_TreasuryRebalanceV2.CallOpts, _zeroedAddress)
}

// GetZeroedIndex is a free data retrieval call binding the contract method 0x518592da.
//
// Solidity: function getZeroedIndex(address _zeroedAddress) view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) GetZeroedIndex(_zeroedAddress common.Address) (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.GetZeroedIndex(&_TreasuryRebalanceV2.CallOpts, _zeroedAddress)
}

// IsContractAddr is a free data retrieval call binding the contract method 0xe2384cb3.
//
// Solidity: function isContractAddr(address _addr) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) IsContractAddr(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "isContractAddr", _addr)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsContractAddr is a free data retrieval call binding the contract method 0xe2384cb3.
//
// Solidity: function isContractAddr(address _addr) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) IsContractAddr(_addr common.Address) (bool, error) {
	return _TreasuryRebalanceV2.Contract.IsContractAddr(&_TreasuryRebalanceV2.CallOpts, _addr)
}

// IsContractAddr is a free data retrieval call binding the contract method 0xe2384cb3.
//
// Solidity: function isContractAddr(address _addr) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) IsContractAddr(_addr common.Address) (bool, error) {
	return _TreasuryRebalanceV2.Contract.IsContractAddr(&_TreasuryRebalanceV2.CallOpts, _addr)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "isOwner")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) IsOwner() (bool, error) {
	return _TreasuryRebalanceV2.Contract.IsOwner(&_TreasuryRebalanceV2.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) IsOwner() (bool, error) {
	return _TreasuryRebalanceV2.Contract.IsOwner(&_TreasuryRebalanceV2.CallOpts)
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) Memo(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "memo")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) Memo() (string, error) {
	return _TreasuryRebalanceV2.Contract.Memo(&_TreasuryRebalanceV2.CallOpts)
}

// Memo is a free data retrieval call binding the contract method 0x58c3b870.
//
// Solidity: function memo() view returns(string)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) Memo() (string, error) {
	return _TreasuryRebalanceV2.Contract.Memo(&_TreasuryRebalanceV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) Owner() (common.Address, error) {
	return _TreasuryRebalanceV2.Contract.Owner(&_TreasuryRebalanceV2.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) Owner() (common.Address, error) {
	return _TreasuryRebalanceV2.Contract.Owner(&_TreasuryRebalanceV2.CallOpts)
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) RebalanceBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "rebalanceBlockNumber")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) RebalanceBlockNumber() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.RebalanceBlockNumber(&_TreasuryRebalanceV2.CallOpts)
}

// RebalanceBlockNumber is a free data retrieval call binding the contract method 0x49a3fb45.
//
// Solidity: function rebalanceBlockNumber() view returns(uint256)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) RebalanceBlockNumber() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.RebalanceBlockNumber(&_TreasuryRebalanceV2.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) Status(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "status")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) Status() (uint8, error) {
	return _TreasuryRebalanceV2.Contract.Status(&_TreasuryRebalanceV2.CallOpts)
}

// Status is a free data retrieval call binding the contract method 0x200d2ed2.
//
// Solidity: function status() view returns(uint8)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) Status() (uint8, error) {
	return _TreasuryRebalanceV2.Contract.Status(&_TreasuryRebalanceV2.CallOpts)
}

// SumOfZeroedBalance is a free data retrieval call binding the contract method 0x9ab29b70.
//
// Solidity: function sumOfZeroedBalance() view returns(uint256 zeroedsBalance)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) SumOfZeroedBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "sumOfZeroedBalance")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// SumOfZeroedBalance is a free data retrieval call binding the contract method 0x9ab29b70.
//
// Solidity: function sumOfZeroedBalance() view returns(uint256 zeroedsBalance)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) SumOfZeroedBalance() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.SumOfZeroedBalance(&_TreasuryRebalanceV2.CallOpts)
}

// SumOfZeroedBalance is a free data retrieval call binding the contract method 0x9ab29b70.
//
// Solidity: function sumOfZeroedBalance() view returns(uint256 zeroedsBalance)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) SumOfZeroedBalance() (*big.Int, error) {
	return _TreasuryRebalanceV2.Contract.SumOfZeroedBalance(&_TreasuryRebalanceV2.CallOpts)
}

// ZeroedExists is a free data retrieval call binding the contract method 0x5f8798c0.
//
// Solidity: function zeroedExists(address _zeroedAddress) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) ZeroedExists(opts *bind.CallOpts, _zeroedAddress common.Address) (bool, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "zeroedExists", _zeroedAddress)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// ZeroedExists is a free data retrieval call binding the contract method 0x5f8798c0.
//
// Solidity: function zeroedExists(address _zeroedAddress) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) ZeroedExists(_zeroedAddress common.Address) (bool, error) {
	return _TreasuryRebalanceV2.Contract.ZeroedExists(&_TreasuryRebalanceV2.CallOpts, _zeroedAddress)
}

// ZeroedExists is a free data retrieval call binding the contract method 0x5f8798c0.
//
// Solidity: function zeroedExists(address _zeroedAddress) view returns(bool)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) ZeroedExists(_zeroedAddress common.Address) (bool, error) {
	return _TreasuryRebalanceV2.Contract.ZeroedExists(&_TreasuryRebalanceV2.CallOpts, _zeroedAddress)
}

// Zeroeds is a free data retrieval call binding the contract method 0x62aa3e91.
//
// Solidity: function zeroeds(uint256 ) view returns(address addr)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Caller) Zeroeds(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TreasuryRebalanceV2.contract.Call(opts, &out, "zeroeds", arg0)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Zeroeds is a free data retrieval call binding the contract method 0x62aa3e91.
//
// Solidity: function zeroeds(uint256 ) view returns(address addr)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) Zeroeds(arg0 *big.Int) (common.Address, error) {
	return _TreasuryRebalanceV2.Contract.Zeroeds(&_TreasuryRebalanceV2.CallOpts, arg0)
}

// Zeroeds is a free data retrieval call binding the contract method 0x62aa3e91.
//
// Solidity: function zeroeds(uint256 ) view returns(address addr)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2CallerSession) Zeroeds(arg0 *big.Int) (common.Address, error) {
	return _TreasuryRebalanceV2.Contract.Zeroeds(&_TreasuryRebalanceV2.CallOpts, arg0)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) Approve(opts *bind.TransactOpts, _zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "approve", _zeroedAddress)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) Approve(_zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.Approve(&_TreasuryRebalanceV2.TransactOpts, _zeroedAddress)
}

// Approve is a paid mutator transaction binding the contract method 0xdaea85c5.
//
// Solidity: function approve(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) Approve(_zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.Approve(&_TreasuryRebalanceV2.TransactOpts, _zeroedAddress)
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) FinalizeApproval(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "finalizeApproval")
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) FinalizeApproval() (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.FinalizeApproval(&_TreasuryRebalanceV2.TransactOpts)
}

// FinalizeApproval is a paid mutator transaction binding the contract method 0xfaaf9ca6.
//
// Solidity: function finalizeApproval() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) FinalizeApproval() (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.FinalizeApproval(&_TreasuryRebalanceV2.TransactOpts)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string _memo) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) FinalizeContract(opts *bind.TransactOpts, _memo string) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "finalizeContract", _memo)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string _memo) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) FinalizeContract(_memo string) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.FinalizeContract(&_TreasuryRebalanceV2.TransactOpts, _memo)
}

// FinalizeContract is a paid mutator transaction binding the contract method 0xea6d4a9b.
//
// Solidity: function finalizeContract(string _memo) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) FinalizeContract(_memo string) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.FinalizeContract(&_TreasuryRebalanceV2.TransactOpts, _memo)
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) FinalizeRegistration(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "finalizeRegistration")
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) FinalizeRegistration() (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.FinalizeRegistration(&_TreasuryRebalanceV2.TransactOpts)
}

// FinalizeRegistration is a paid mutator transaction binding the contract method 0x48409096.
//
// Solidity: function finalizeRegistration() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) FinalizeRegistration() (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.FinalizeRegistration(&_TreasuryRebalanceV2.TransactOpts)
}

// RegisterAllocated is a paid mutator transaction binding the contract method 0xecd86778.
//
// Solidity: function registerAllocated(address _allocatedAddress, uint256 _amount) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) RegisterAllocated(opts *bind.TransactOpts, _allocatedAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "registerAllocated", _allocatedAddress, _amount)
}

// RegisterAllocated is a paid mutator transaction binding the contract method 0xecd86778.
//
// Solidity: function registerAllocated(address _allocatedAddress, uint256 _amount) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) RegisterAllocated(_allocatedAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RegisterAllocated(&_TreasuryRebalanceV2.TransactOpts, _allocatedAddress, _amount)
}

// RegisterAllocated is a paid mutator transaction binding the contract method 0xecd86778.
//
// Solidity: function registerAllocated(address _allocatedAddress, uint256 _amount) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) RegisterAllocated(_allocatedAddress common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RegisterAllocated(&_TreasuryRebalanceV2.TransactOpts, _allocatedAddress, _amount)
}

// RegisterZeroed is a paid mutator transaction binding the contract method 0x5f9b0df7.
//
// Solidity: function registerZeroed(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) RegisterZeroed(opts *bind.TransactOpts, _zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "registerZeroed", _zeroedAddress)
}

// RegisterZeroed is a paid mutator transaction binding the contract method 0x5f9b0df7.
//
// Solidity: function registerZeroed(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) RegisterZeroed(_zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RegisterZeroed(&_TreasuryRebalanceV2.TransactOpts, _zeroedAddress)
}

// RegisterZeroed is a paid mutator transaction binding the contract method 0x5f9b0df7.
//
// Solidity: function registerZeroed(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) RegisterZeroed(_zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RegisterZeroed(&_TreasuryRebalanceV2.TransactOpts, _zeroedAddress)
}

// RemoveAllocated is a paid mutator transaction binding the contract method 0x27704cb5.
//
// Solidity: function removeAllocated(address _allocatedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) RemoveAllocated(opts *bind.TransactOpts, _allocatedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "removeAllocated", _allocatedAddress)
}

// RemoveAllocated is a paid mutator transaction binding the contract method 0x27704cb5.
//
// Solidity: function removeAllocated(address _allocatedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) RemoveAllocated(_allocatedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RemoveAllocated(&_TreasuryRebalanceV2.TransactOpts, _allocatedAddress)
}

// RemoveAllocated is a paid mutator transaction binding the contract method 0x27704cb5.
//
// Solidity: function removeAllocated(address _allocatedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) RemoveAllocated(_allocatedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RemoveAllocated(&_TreasuryRebalanceV2.TransactOpts, _allocatedAddress)
}

// RemoveZeroed is a paid mutator transaction binding the contract method 0xdb27b50b.
//
// Solidity: function removeZeroed(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) RemoveZeroed(opts *bind.TransactOpts, _zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "removeZeroed", _zeroedAddress)
}

// RemoveZeroed is a paid mutator transaction binding the contract method 0xdb27b50b.
//
// Solidity: function removeZeroed(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) RemoveZeroed(_zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RemoveZeroed(&_TreasuryRebalanceV2.TransactOpts, _zeroedAddress)
}

// RemoveZeroed is a paid mutator transaction binding the contract method 0xdb27b50b.
//
// Solidity: function removeZeroed(address _zeroedAddress) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) RemoveZeroed(_zeroedAddress common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RemoveZeroed(&_TreasuryRebalanceV2.TransactOpts, _zeroedAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) RenounceOwnership() (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RenounceOwnership(&_TreasuryRebalanceV2.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.RenounceOwnership(&_TreasuryRebalanceV2.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "reset")
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) Reset() (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.Reset(&_TreasuryRebalanceV2.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) Reset() (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.Reset(&_TreasuryRebalanceV2.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.TransferOwnership(&_TreasuryRebalanceV2.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.TransferOwnership(&_TreasuryRebalanceV2.TransactOpts, newOwner)
}

// UpdateRebalanceBlocknumber is a paid mutator transaction binding the contract method 0x1804692f.
//
// Solidity: function updateRebalanceBlocknumber(uint256 _rebalanceBlockNumber) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Transactor) UpdateRebalanceBlocknumber(opts *bind.TransactOpts, _rebalanceBlockNumber *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.contract.Transact(opts, "updateRebalanceBlocknumber", _rebalanceBlockNumber)
}

// UpdateRebalanceBlocknumber is a paid mutator transaction binding the contract method 0x1804692f.
//
// Solidity: function updateRebalanceBlocknumber(uint256 _rebalanceBlockNumber) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Session) UpdateRebalanceBlocknumber(_rebalanceBlockNumber *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.UpdateRebalanceBlocknumber(&_TreasuryRebalanceV2.TransactOpts, _rebalanceBlockNumber)
}

// UpdateRebalanceBlocknumber is a paid mutator transaction binding the contract method 0x1804692f.
//
// Solidity: function updateRebalanceBlocknumber(uint256 _rebalanceBlockNumber) returns()
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2TransactorSession) UpdateRebalanceBlocknumber(_rebalanceBlockNumber *big.Int) (*types.Transaction, error) {
	return _TreasuryRebalanceV2.Contract.UpdateRebalanceBlocknumber(&_TreasuryRebalanceV2.TransactOpts, _rebalanceBlockNumber)
}

// TreasuryRebalanceV2AllocatedRegisteredIterator is returned from FilterAllocatedRegistered and is used to iterate over the raw logs and unpacked data for AllocatedRegistered events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2AllocatedRegisteredIterator struct {
	Event *TreasuryRebalanceV2AllocatedRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2AllocatedRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2AllocatedRegistered)
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
		it.Event = new(TreasuryRebalanceV2AllocatedRegistered)
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
func (it *TreasuryRebalanceV2AllocatedRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2AllocatedRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2AllocatedRegistered represents a AllocatedRegistered event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2AllocatedRegistered struct {
	Allocated      common.Address
	FundAllocation *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterAllocatedRegistered is a free log retrieval operation binding the contract event 0xab5b2126f71ee7e0b39eadc53fb5d08a8f6c68dc61795fa05ed7d176cd2665ed.
//
// Solidity: event AllocatedRegistered(address allocated, uint256 fundAllocation)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterAllocatedRegistered(opts *bind.FilterOpts) (*TreasuryRebalanceV2AllocatedRegisteredIterator, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "AllocatedRegistered")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2AllocatedRegisteredIterator{contract: _TreasuryRebalanceV2.contract, event: "AllocatedRegistered", logs: logs, sub: sub}, nil
}

// WatchAllocatedRegistered is a free log subscription operation binding the contract event 0xab5b2126f71ee7e0b39eadc53fb5d08a8f6c68dc61795fa05ed7d176cd2665ed.
//
// Solidity: event AllocatedRegistered(address allocated, uint256 fundAllocation)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchAllocatedRegistered(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2AllocatedRegistered) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "AllocatedRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2AllocatedRegistered)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "AllocatedRegistered", log); err != nil {
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

// ParseAllocatedRegistered is a log parse operation binding the contract event 0xab5b2126f71ee7e0b39eadc53fb5d08a8f6c68dc61795fa05ed7d176cd2665ed.
//
// Solidity: event AllocatedRegistered(address allocated, uint256 fundAllocation)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseAllocatedRegistered(log types.Log) (*TreasuryRebalanceV2AllocatedRegistered, error) {
	event := new(TreasuryRebalanceV2AllocatedRegistered)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "AllocatedRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2AllocatedRemovedIterator is returned from FilterAllocatedRemoved and is used to iterate over the raw logs and unpacked data for AllocatedRemoved events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2AllocatedRemovedIterator struct {
	Event *TreasuryRebalanceV2AllocatedRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2AllocatedRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2AllocatedRemoved)
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
		it.Event = new(TreasuryRebalanceV2AllocatedRemoved)
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
func (it *TreasuryRebalanceV2AllocatedRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2AllocatedRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2AllocatedRemoved represents a AllocatedRemoved event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2AllocatedRemoved struct {
	Allocated common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAllocatedRemoved is a free log retrieval operation binding the contract event 0xf8f67464bea52432645435be9c46c427173a75aefaa1001272e08a4b8572f06e.
//
// Solidity: event AllocatedRemoved(address allocated)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterAllocatedRemoved(opts *bind.FilterOpts) (*TreasuryRebalanceV2AllocatedRemovedIterator, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "AllocatedRemoved")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2AllocatedRemovedIterator{contract: _TreasuryRebalanceV2.contract, event: "AllocatedRemoved", logs: logs, sub: sub}, nil
}

// WatchAllocatedRemoved is a free log subscription operation binding the contract event 0xf8f67464bea52432645435be9c46c427173a75aefaa1001272e08a4b8572f06e.
//
// Solidity: event AllocatedRemoved(address allocated)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchAllocatedRemoved(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2AllocatedRemoved) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "AllocatedRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2AllocatedRemoved)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "AllocatedRemoved", log); err != nil {
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

// ParseAllocatedRemoved is a log parse operation binding the contract event 0xf8f67464bea52432645435be9c46c427173a75aefaa1001272e08a4b8572f06e.
//
// Solidity: event AllocatedRemoved(address allocated)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseAllocatedRemoved(log types.Log) (*TreasuryRebalanceV2AllocatedRemoved, error) {
	event := new(TreasuryRebalanceV2AllocatedRemoved)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "AllocatedRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2ApprovedIterator is returned from FilterApproved and is used to iterate over the raw logs and unpacked data for Approved events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2ApprovedIterator struct {
	Event *TreasuryRebalanceV2Approved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2ApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2Approved)
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
		it.Event = new(TreasuryRebalanceV2Approved)
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
func (it *TreasuryRebalanceV2ApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2ApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2Approved represents a Approved event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2Approved struct {
	Zeroed         common.Address
	Approver       common.Address
	ApproversCount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterApproved is a free log retrieval operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address zeroed, address approver, uint256 approversCount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterApproved(opts *bind.FilterOpts) (*TreasuryRebalanceV2ApprovedIterator, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "Approved")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2ApprovedIterator{contract: _TreasuryRebalanceV2.contract, event: "Approved", logs: logs, sub: sub}, nil
}

// WatchApproved is a free log subscription operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address zeroed, address approver, uint256 approversCount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchApproved(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2Approved) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "Approved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2Approved)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "Approved", log); err != nil {
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

// ParseApproved is a log parse operation binding the contract event 0x80da462ebfbe41cfc9bc015e7a9a3c7a2a73dbccede72d8ceb583606c27f8f90.
//
// Solidity: event Approved(address zeroed, address approver, uint256 approversCount)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseApproved(log types.Log) (*TreasuryRebalanceV2Approved, error) {
	event := new(TreasuryRebalanceV2Approved)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "Approved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2ContractDeployedIterator is returned from FilterContractDeployed and is used to iterate over the raw logs and unpacked data for ContractDeployed events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2ContractDeployedIterator struct {
	Event *TreasuryRebalanceV2ContractDeployed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2ContractDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2ContractDeployed)
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
		it.Event = new(TreasuryRebalanceV2ContractDeployed)
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
func (it *TreasuryRebalanceV2ContractDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2ContractDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2ContractDeployed represents a ContractDeployed event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2ContractDeployed struct {
	Status               uint8
	RebalanceBlockNumber *big.Int
	DeployedBlockNumber  *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterContractDeployed is a free log retrieval operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterContractDeployed(opts *bind.FilterOpts) (*TreasuryRebalanceV2ContractDeployedIterator, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "ContractDeployed")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2ContractDeployedIterator{contract: _TreasuryRebalanceV2.contract, event: "ContractDeployed", logs: logs, sub: sub}, nil
}

// WatchContractDeployed is a free log subscription operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchContractDeployed(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2ContractDeployed) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "ContractDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2ContractDeployed)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "ContractDeployed", log); err != nil {
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

// ParseContractDeployed is a log parse operation binding the contract event 0x6f182006c5a12fe70c0728eedb2d1b0628c41483ca6721c606707d778d22ed0a.
//
// Solidity: event ContractDeployed(uint8 status, uint256 rebalanceBlockNumber, uint256 deployedBlockNumber)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseContractDeployed(log types.Log) (*TreasuryRebalanceV2ContractDeployed, error) {
	event := new(TreasuryRebalanceV2ContractDeployed)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "ContractDeployed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2FinalizedIterator is returned from FilterFinalized and is used to iterate over the raw logs and unpacked data for Finalized events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2FinalizedIterator struct {
	Event *TreasuryRebalanceV2Finalized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2FinalizedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2Finalized)
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
		it.Event = new(TreasuryRebalanceV2Finalized)
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
func (it *TreasuryRebalanceV2FinalizedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2FinalizedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2Finalized represents a Finalized event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2Finalized struct {
	Memo   string
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFinalized is a free log retrieval operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterFinalized(opts *bind.FilterOpts) (*TreasuryRebalanceV2FinalizedIterator, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "Finalized")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2FinalizedIterator{contract: _TreasuryRebalanceV2.contract, event: "Finalized", logs: logs, sub: sub}, nil
}

// WatchFinalized is a free log subscription operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchFinalized(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2Finalized) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "Finalized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2Finalized)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "Finalized", log); err != nil {
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

// ParseFinalized is a log parse operation binding the contract event 0x8f8636c7757ca9b7d154e1d44ca90d8e8c885b9eac417c59bbf8eb7779ca6404.
//
// Solidity: event Finalized(string memo, uint8 status)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseFinalized(log types.Log) (*TreasuryRebalanceV2Finalized, error) {
	event := new(TreasuryRebalanceV2Finalized)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "Finalized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2OwnershipTransferredIterator struct {
	Event *TreasuryRebalanceV2OwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2OwnershipTransferred)
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
		it.Event = new(TreasuryRebalanceV2OwnershipTransferred)
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
func (it *TreasuryRebalanceV2OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2OwnershipTransferred represents a OwnershipTransferred event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TreasuryRebalanceV2OwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2OwnershipTransferredIterator{contract: _TreasuryRebalanceV2.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2OwnershipTransferred)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseOwnershipTransferred(log types.Log) (*TreasuryRebalanceV2OwnershipTransferred, error) {
	event := new(TreasuryRebalanceV2OwnershipTransferred)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2StatusChangedIterator is returned from FilterStatusChanged and is used to iterate over the raw logs and unpacked data for StatusChanged events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2StatusChangedIterator struct {
	Event *TreasuryRebalanceV2StatusChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2StatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2StatusChanged)
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
		it.Event = new(TreasuryRebalanceV2StatusChanged)
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
func (it *TreasuryRebalanceV2StatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2StatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2StatusChanged represents a StatusChanged event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2StatusChanged struct {
	Status uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStatusChanged is a free log retrieval operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterStatusChanged(opts *bind.FilterOpts) (*TreasuryRebalanceV2StatusChangedIterator, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "StatusChanged")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2StatusChangedIterator{contract: _TreasuryRebalanceV2.contract, event: "StatusChanged", logs: logs, sub: sub}, nil
}

// WatchStatusChanged is a free log subscription operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchStatusChanged(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2StatusChanged) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "StatusChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2StatusChanged)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "StatusChanged", log); err != nil {
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

// ParseStatusChanged is a log parse operation binding the contract event 0xafa725e7f44cadb687a7043853fa1a7e7b8f0da74ce87ec546e9420f04da8c1e.
//
// Solidity: event StatusChanged(uint8 status)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseStatusChanged(log types.Log) (*TreasuryRebalanceV2StatusChanged, error) {
	event := new(TreasuryRebalanceV2StatusChanged)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "StatusChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2ZeroedRegisteredIterator is returned from FilterZeroedRegistered and is used to iterate over the raw logs and unpacked data for ZeroedRegistered events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2ZeroedRegisteredIterator struct {
	Event *TreasuryRebalanceV2ZeroedRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2ZeroedRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2ZeroedRegistered)
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
		it.Event = new(TreasuryRebalanceV2ZeroedRegistered)
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
func (it *TreasuryRebalanceV2ZeroedRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2ZeroedRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2ZeroedRegistered represents a ZeroedRegistered event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2ZeroedRegistered struct {
	Zeroed common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterZeroedRegistered is a free log retrieval operation binding the contract event 0xa9a4f3b74b03e48e76814dbc308d3f20104d608c67a42a7ae678d0945daa8e92.
//
// Solidity: event ZeroedRegistered(address zeroed)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterZeroedRegistered(opts *bind.FilterOpts) (*TreasuryRebalanceV2ZeroedRegisteredIterator, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "ZeroedRegistered")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2ZeroedRegisteredIterator{contract: _TreasuryRebalanceV2.contract, event: "ZeroedRegistered", logs: logs, sub: sub}, nil
}

// WatchZeroedRegistered is a free log subscription operation binding the contract event 0xa9a4f3b74b03e48e76814dbc308d3f20104d608c67a42a7ae678d0945daa8e92.
//
// Solidity: event ZeroedRegistered(address zeroed)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchZeroedRegistered(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2ZeroedRegistered) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "ZeroedRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2ZeroedRegistered)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "ZeroedRegistered", log); err != nil {
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

// ParseZeroedRegistered is a log parse operation binding the contract event 0xa9a4f3b74b03e48e76814dbc308d3f20104d608c67a42a7ae678d0945daa8e92.
//
// Solidity: event ZeroedRegistered(address zeroed)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseZeroedRegistered(log types.Log) (*TreasuryRebalanceV2ZeroedRegistered, error) {
	event := new(TreasuryRebalanceV2ZeroedRegistered)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "ZeroedRegistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// TreasuryRebalanceV2ZeroedRemovedIterator is returned from FilterZeroedRemoved and is used to iterate over the raw logs and unpacked data for ZeroedRemoved events raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2ZeroedRemovedIterator struct {
	Event *TreasuryRebalanceV2ZeroedRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TreasuryRebalanceV2ZeroedRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TreasuryRebalanceV2ZeroedRemoved)
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
		it.Event = new(TreasuryRebalanceV2ZeroedRemoved)
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
func (it *TreasuryRebalanceV2ZeroedRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TreasuryRebalanceV2ZeroedRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TreasuryRebalanceV2ZeroedRemoved represents a ZeroedRemoved event raised by the TreasuryRebalanceV2 contract.
type TreasuryRebalanceV2ZeroedRemoved struct {
	Zeroed common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterZeroedRemoved is a free log retrieval operation binding the contract event 0x8a654c98d0a7856a8d216c621bb8073316efcaa2b65774d2050c4c1fc7a85a0c.
//
// Solidity: event ZeroedRemoved(address zeroed)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) FilterZeroedRemoved(opts *bind.FilterOpts) (*TreasuryRebalanceV2ZeroedRemovedIterator, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.FilterLogs(opts, "ZeroedRemoved")
	if err != nil {
		return nil, err
	}
	return &TreasuryRebalanceV2ZeroedRemovedIterator{contract: _TreasuryRebalanceV2.contract, event: "ZeroedRemoved", logs: logs, sub: sub}, nil
}

// WatchZeroedRemoved is a free log subscription operation binding the contract event 0x8a654c98d0a7856a8d216c621bb8073316efcaa2b65774d2050c4c1fc7a85a0c.
//
// Solidity: event ZeroedRemoved(address zeroed)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) WatchZeroedRemoved(opts *bind.WatchOpts, sink chan<- *TreasuryRebalanceV2ZeroedRemoved) (event.Subscription, error) {
	logs, sub, err := _TreasuryRebalanceV2.contract.WatchLogs(opts, "ZeroedRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TreasuryRebalanceV2ZeroedRemoved)
				if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "ZeroedRemoved", log); err != nil {
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

// ParseZeroedRemoved is a log parse operation binding the contract event 0x8a654c98d0a7856a8d216c621bb8073316efcaa2b65774d2050c4c1fc7a85a0c.
//
// Solidity: event ZeroedRemoved(address zeroed)
func (_TreasuryRebalanceV2 *TreasuryRebalanceV2Filterer) ParseZeroedRemoved(log types.Log) (*TreasuryRebalanceV2ZeroedRemoved, error) {
	event := new(TreasuryRebalanceV2ZeroedRemoved)
	if err := _TreasuryRebalanceV2.contract.UnpackLog(event, "ZeroedRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}
