// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package factory

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

// IERC20MetaData contains all meta data concerning the IERC20 contract.
var IERC20MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"313ce567": "decimals()",
		"06fdde03": "name()",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
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

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "allowance", owner, spender)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20 *IERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20.Contract.Allowance(&_IERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IERC20 *IERC20Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "balanceOf", owner)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IERC20 *IERC20Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IERC20 *IERC20CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IERC20.Contract.BalanceOf(&_IERC20.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_IERC20 *IERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "decimals")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_IERC20 *IERC20Session) Decimals() (uint8, error) {
	return _IERC20.Contract.Decimals(&_IERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_IERC20 *IERC20CallerSession) Decimals() (uint8, error) {
	return _IERC20.Contract.Decimals(&_IERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IERC20 *IERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "name")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IERC20 *IERC20Session) Name() (string, error) {
	return _IERC20.Contract.Name(&_IERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IERC20 *IERC20CallerSession) Name() (string, error) {
	return _IERC20.Contract.Name(&_IERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_IERC20 *IERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "symbol")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_IERC20 *IERC20Session) Symbol() (string, error) {
	return _IERC20.Contract.Symbol(&_IERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_IERC20 *IERC20CallerSession) Symbol() (string, error) {
	return _IERC20.Contract.Symbol(&_IERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IERC20.contract.Call(opts, &out, "totalSupply")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20Session) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20 *IERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _IERC20.Contract.TotalSupply(&_IERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IERC20 *IERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IERC20 *IERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IERC20 *IERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Approve(&_IERC20.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IERC20 *IERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IERC20 *IERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IERC20 *IERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.Transfer(&_IERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IERC20 *IERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IERC20 *IERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IERC20 *IERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20.Contract.TransferFrom(&_IERC20.TransactOpts, from, to, value)
}

// IERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IERC20 contract.
type IERC20ApprovalIterator struct {
	Event *IERC20Approval // Event containing the contract specifics and raw log

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
func (it *IERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20Approval)
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
		it.Event = new(IERC20Approval)
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
func (it *IERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20Approval represents a Approval event raised by the IERC20 contract.
type IERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IERC20ApprovalIterator, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IERC20ApprovalIterator{contract: _IERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20Approval)
				if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20 *IERC20Filterer) ParseApproval(log types.Log) (*IERC20Approval, error) {
	event := new(IERC20Approval)
	if err := _IERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IERC20 contract.
type IERC20TransferIterator struct {
	Event *IERC20Transfer // Event containing the contract specifics and raw log

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
func (it *IERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20Transfer)
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
		it.Event = new(IERC20Transfer)
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
func (it *IERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20Transfer represents a Transfer event raised by the IERC20 contract.
type IERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IERC20TransferIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IERC20TransferIterator{contract: _IERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20Transfer)
				if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20 *IERC20Filterer) ParseTransfer(log types.Log) (*IERC20Transfer, error) {
	event := new(IERC20Transfer)
	if err := _IERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2CalleeMetaData contains all meta data concerning the IUniswapV2Callee contract.
var IUniswapV2CalleeMetaData = &bind.MetaData{
	ABI: "[{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"uniswapV2Call\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"10d1e85c": "uniswapV2Call(address,uint256,uint256,bytes)",
	},
}

// IUniswapV2CalleeABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV2CalleeMetaData.ABI instead.
var IUniswapV2CalleeABI = IUniswapV2CalleeMetaData.ABI

// IUniswapV2CalleeBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IUniswapV2CalleeBinRuntime = ``

// Deprecated: Use IUniswapV2CalleeMetaData.Sigs instead.
// IUniswapV2CalleeFuncSigs maps the 4-byte function signature to its string representation.
var IUniswapV2CalleeFuncSigs = IUniswapV2CalleeMetaData.Sigs

// IUniswapV2Callee is an auto generated Go binding around a Kaia contract.
type IUniswapV2Callee struct {
	IUniswapV2CalleeCaller     // Read-only binding to the contract
	IUniswapV2CalleeTransactor // Write-only binding to the contract
	IUniswapV2CalleeFilterer   // Log filterer for contract events
}

// IUniswapV2CalleeCaller is an auto generated read-only Go binding around a Kaia contract.
type IUniswapV2CalleeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2CalleeTransactor is an auto generated write-only Go binding around a Kaia contract.
type IUniswapV2CalleeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2CalleeFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IUniswapV2CalleeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2CalleeSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IUniswapV2CalleeSession struct {
	Contract     *IUniswapV2Callee // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IUniswapV2CalleeCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IUniswapV2CalleeCallerSession struct {
	Contract *IUniswapV2CalleeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IUniswapV2CalleeTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IUniswapV2CalleeTransactorSession struct {
	Contract     *IUniswapV2CalleeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IUniswapV2CalleeRaw is an auto generated low-level Go binding around a Kaia contract.
type IUniswapV2CalleeRaw struct {
	Contract *IUniswapV2Callee // Generic contract binding to access the raw methods on
}

// IUniswapV2CalleeCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IUniswapV2CalleeCallerRaw struct {
	Contract *IUniswapV2CalleeCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV2CalleeTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IUniswapV2CalleeTransactorRaw struct {
	Contract *IUniswapV2CalleeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV2Callee creates a new instance of IUniswapV2Callee, bound to a specific deployed contract.
func NewIUniswapV2Callee(address common.Address, backend bind.ContractBackend) (*IUniswapV2Callee, error) {
	contract, err := bindIUniswapV2Callee(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2Callee{IUniswapV2CalleeCaller: IUniswapV2CalleeCaller{contract: contract}, IUniswapV2CalleeTransactor: IUniswapV2CalleeTransactor{contract: contract}, IUniswapV2CalleeFilterer: IUniswapV2CalleeFilterer{contract: contract}}, nil
}

// NewIUniswapV2CalleeCaller creates a new read-only instance of IUniswapV2Callee, bound to a specific deployed contract.
func NewIUniswapV2CalleeCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV2CalleeCaller, error) {
	contract, err := bindIUniswapV2Callee(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2CalleeCaller{contract: contract}, nil
}

// NewIUniswapV2CalleeTransactor creates a new write-only instance of IUniswapV2Callee, bound to a specific deployed contract.
func NewIUniswapV2CalleeTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV2CalleeTransactor, error) {
	contract, err := bindIUniswapV2Callee(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2CalleeTransactor{contract: contract}, nil
}

// NewIUniswapV2CalleeFilterer creates a new log filterer instance of IUniswapV2Callee, bound to a specific deployed contract.
func NewIUniswapV2CalleeFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV2CalleeFilterer, error) {
	contract, err := bindIUniswapV2Callee(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2CalleeFilterer{contract: contract}, nil
}

// bindIUniswapV2Callee binds a generic wrapper to an already deployed contract.
func bindIUniswapV2Callee(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV2CalleeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Callee *IUniswapV2CalleeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Callee.Contract.IUniswapV2CalleeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Callee *IUniswapV2CalleeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.IUniswapV2CalleeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Callee *IUniswapV2CalleeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.IUniswapV2CalleeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Callee *IUniswapV2CalleeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Callee.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Callee *IUniswapV2CalleeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Callee *IUniswapV2CalleeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.contract.Transact(opts, method, params...)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV2Callee *IUniswapV2CalleeTransactor) UniswapV2Call(opts *bind.TransactOpts, sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Callee.contract.Transact(opts, "uniswapV2Call", sender, amount0, amount1, data)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV2Callee *IUniswapV2CalleeSession) UniswapV2Call(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.UniswapV2Call(&_IUniswapV2Callee.TransactOpts, sender, amount0, amount1, data)
}

// UniswapV2Call is a paid mutator transaction binding the contract method 0x10d1e85c.
//
// Solidity: function uniswapV2Call(address sender, uint256 amount0, uint256 amount1, bytes data) returns()
func (_IUniswapV2Callee *IUniswapV2CalleeTransactorSession) UniswapV2Call(sender common.Address, amount0 *big.Int, amount1 *big.Int, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Callee.Contract.UniswapV2Call(&_IUniswapV2Callee.TransactOpts, sender, amount0, amount1, data)
}

// IUniswapV2ERC20MetaData contains all meta data concerning the IUniswapV2ERC20 contract.
var IUniswapV2ERC20MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"3644e515": "DOMAIN_SEPARATOR()",
		"30adf81f": "PERMIT_TYPEHASH()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"313ce567": "decimals()",
		"06fdde03": "name()",
		"7ecebe00": "nonces(address)",
		"d505accf": "permit(address,address,uint256,uint256,uint8,bytes32,bytes32)",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
	},
}

// IUniswapV2ERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV2ERC20MetaData.ABI instead.
var IUniswapV2ERC20ABI = IUniswapV2ERC20MetaData.ABI

// IUniswapV2ERC20BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IUniswapV2ERC20BinRuntime = ``

// Deprecated: Use IUniswapV2ERC20MetaData.Sigs instead.
// IUniswapV2ERC20FuncSigs maps the 4-byte function signature to its string representation.
var IUniswapV2ERC20FuncSigs = IUniswapV2ERC20MetaData.Sigs

// IUniswapV2ERC20 is an auto generated Go binding around a Kaia contract.
type IUniswapV2ERC20 struct {
	IUniswapV2ERC20Caller     // Read-only binding to the contract
	IUniswapV2ERC20Transactor // Write-only binding to the contract
	IUniswapV2ERC20Filterer   // Log filterer for contract events
}

// IUniswapV2ERC20Caller is an auto generated read-only Go binding around a Kaia contract.
type IUniswapV2ERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2ERC20Transactor is an auto generated write-only Go binding around a Kaia contract.
type IUniswapV2ERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2ERC20Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type IUniswapV2ERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2ERC20Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IUniswapV2ERC20Session struct {
	Contract     *IUniswapV2ERC20  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IUniswapV2ERC20CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IUniswapV2ERC20CallerSession struct {
	Contract *IUniswapV2ERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IUniswapV2ERC20TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IUniswapV2ERC20TransactorSession struct {
	Contract     *IUniswapV2ERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IUniswapV2ERC20Raw is an auto generated low-level Go binding around a Kaia contract.
type IUniswapV2ERC20Raw struct {
	Contract *IUniswapV2ERC20 // Generic contract binding to access the raw methods on
}

// IUniswapV2ERC20CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IUniswapV2ERC20CallerRaw struct {
	Contract *IUniswapV2ERC20Caller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV2ERC20TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IUniswapV2ERC20TransactorRaw struct {
	Contract *IUniswapV2ERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV2ERC20 creates a new instance of IUniswapV2ERC20, bound to a specific deployed contract.
func NewIUniswapV2ERC20(address common.Address, backend bind.ContractBackend) (*IUniswapV2ERC20, error) {
	contract, err := bindIUniswapV2ERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2ERC20{IUniswapV2ERC20Caller: IUniswapV2ERC20Caller{contract: contract}, IUniswapV2ERC20Transactor: IUniswapV2ERC20Transactor{contract: contract}, IUniswapV2ERC20Filterer: IUniswapV2ERC20Filterer{contract: contract}}, nil
}

// NewIUniswapV2ERC20Caller creates a new read-only instance of IUniswapV2ERC20, bound to a specific deployed contract.
func NewIUniswapV2ERC20Caller(address common.Address, caller bind.ContractCaller) (*IUniswapV2ERC20Caller, error) {
	contract, err := bindIUniswapV2ERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2ERC20Caller{contract: contract}, nil
}

// NewIUniswapV2ERC20Transactor creates a new write-only instance of IUniswapV2ERC20, bound to a specific deployed contract.
func NewIUniswapV2ERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV2ERC20Transactor, error) {
	contract, err := bindIUniswapV2ERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2ERC20Transactor{contract: contract}, nil
}

// NewIUniswapV2ERC20Filterer creates a new log filterer instance of IUniswapV2ERC20, bound to a specific deployed contract.
func NewIUniswapV2ERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV2ERC20Filterer, error) {
	contract, err := bindIUniswapV2ERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2ERC20Filterer{contract: contract}, nil
}

// bindIUniswapV2ERC20 binds a generic wrapper to an already deployed contract.
func bindIUniswapV2ERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV2ERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2ERC20 *IUniswapV2ERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2ERC20.Contract.IUniswapV2ERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2ERC20 *IUniswapV2ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.IUniswapV2ERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2ERC20 *IUniswapV2ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.IUniswapV2ERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2ERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2ERC20 *IUniswapV2ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2ERC20 *IUniswapV2ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "DOMAIN_SEPARATOR")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) DOMAINSEPARATOR() ([32]byte, error) {
	return _IUniswapV2ERC20.Contract.DOMAINSEPARATOR(&_IUniswapV2ERC20.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _IUniswapV2ERC20.Contract.DOMAINSEPARATOR(&_IUniswapV2ERC20.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() pure returns(bytes32)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "PERMIT_TYPEHASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() pure returns(bytes32)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) PERMITTYPEHASH() ([32]byte, error) {
	return _IUniswapV2ERC20.Contract.PERMITTYPEHASH(&_IUniswapV2ERC20.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() pure returns(bytes32)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _IUniswapV2ERC20.Contract.PERMITTYPEHASH(&_IUniswapV2ERC20.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "allowance", owner, spender)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IUniswapV2ERC20.Contract.Allowance(&_IUniswapV2ERC20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IUniswapV2ERC20.Contract.Allowance(&_IUniswapV2ERC20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "balanceOf", owner)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IUniswapV2ERC20.Contract.BalanceOf(&_IUniswapV2ERC20.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IUniswapV2ERC20.Contract.BalanceOf(&_IUniswapV2ERC20.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "decimals")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) Decimals() (uint8, error) {
	return _IUniswapV2ERC20.Contract.Decimals(&_IUniswapV2ERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) Decimals() (uint8, error) {
	return _IUniswapV2ERC20.Contract.Decimals(&_IUniswapV2ERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "name")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) Name() (string, error) {
	return _IUniswapV2ERC20.Contract.Name(&_IUniswapV2ERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) Name() (string, error) {
	return _IUniswapV2ERC20.Contract.Name(&_IUniswapV2ERC20.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "nonces", owner)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) Nonces(owner common.Address) (*big.Int, error) {
	return _IUniswapV2ERC20.Contract.Nonces(&_IUniswapV2ERC20.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _IUniswapV2ERC20.Contract.Nonces(&_IUniswapV2ERC20.CallOpts, owner)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "symbol")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) Symbol() (string, error) {
	return _IUniswapV2ERC20.Contract.Symbol(&_IUniswapV2ERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) Symbol() (string, error) {
	return _IUniswapV2ERC20.Contract.Symbol(&_IUniswapV2ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2ERC20.contract.Call(opts, &out, "totalSupply")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) TotalSupply() (*big.Int, error) {
	return _IUniswapV2ERC20.Contract.TotalSupply(&_IUniswapV2ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IUniswapV2ERC20 *IUniswapV2ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _IUniswapV2ERC20.Contract.TotalSupply(&_IUniswapV2ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.Approve(&_IUniswapV2ERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.Approve(&_IUniswapV2ERC20.TransactOpts, spender, value)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_IUniswapV2ERC20 *IUniswapV2ERC20Transactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _IUniswapV2ERC20.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.Permit(&_IUniswapV2ERC20.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_IUniswapV2ERC20 *IUniswapV2ERC20TransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.Permit(&_IUniswapV2ERC20.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.Transfer(&_IUniswapV2ERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.Transfer(&_IUniswapV2ERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.TransferFrom(&_IUniswapV2ERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IUniswapV2ERC20 *IUniswapV2ERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2ERC20.Contract.TransferFrom(&_IUniswapV2ERC20.TransactOpts, from, to, value)
}

// IUniswapV2ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IUniswapV2ERC20 contract.
type IUniswapV2ERC20ApprovalIterator struct {
	Event *IUniswapV2ERC20Approval // Event containing the contract specifics and raw log

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
func (it *IUniswapV2ERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2ERC20Approval)
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
		it.Event = new(IUniswapV2ERC20Approval)
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
func (it *IUniswapV2ERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2ERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2ERC20Approval represents a Approval event raised by the IUniswapV2ERC20 contract.
type IUniswapV2ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IUniswapV2ERC20ApprovalIterator, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IUniswapV2ERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2ERC20ApprovalIterator{contract: _IUniswapV2ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IUniswapV2ERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IUniswapV2ERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2ERC20Approval)
				if err := _IUniswapV2ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Filterer) ParseApproval(log types.Log) (*IUniswapV2ERC20Approval, error) {
	event := new(IUniswapV2ERC20Approval)
	if err := _IUniswapV2ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IUniswapV2ERC20 contract.
type IUniswapV2ERC20TransferIterator struct {
	Event *IUniswapV2ERC20Transfer // Event containing the contract specifics and raw log

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
func (it *IUniswapV2ERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2ERC20Transfer)
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
		it.Event = new(IUniswapV2ERC20Transfer)
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
func (it *IUniswapV2ERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2ERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2ERC20Transfer represents a Transfer event raised by the IUniswapV2ERC20 contract.
type IUniswapV2ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IUniswapV2ERC20TransferIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IUniswapV2ERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2ERC20TransferIterator{contract: _IUniswapV2ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IUniswapV2ERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IUniswapV2ERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2ERC20Transfer)
				if err := _IUniswapV2ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IUniswapV2ERC20 *IUniswapV2ERC20Filterer) ParseTransfer(log types.Log) (*IUniswapV2ERC20Transfer, error) {
	event := new(IUniswapV2ERC20Transfer)
	if err := _IUniswapV2ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2FactoryMetaData contains all meta data concerning the IUniswapV2Factory contract.
var IUniswapV2FactoryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"PairCreated\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allPairs\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"allPairsLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"createPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeToSetter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"getPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"setFeeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"setFeeToSetter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1e3dd18b": "allPairs(uint256)",
		"574f2ba3": "allPairsLength()",
		"c9c65396": "createPair(address,address)",
		"017e7e58": "feeTo()",
		"094b7415": "feeToSetter()",
		"e6a43905": "getPair(address,address)",
		"f46901ed": "setFeeTo(address)",
		"a2e74af6": "setFeeToSetter(address)",
	},
}

// IUniswapV2FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV2FactoryMetaData.ABI instead.
var IUniswapV2FactoryABI = IUniswapV2FactoryMetaData.ABI

// IUniswapV2FactoryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IUniswapV2FactoryBinRuntime = ``

// Deprecated: Use IUniswapV2FactoryMetaData.Sigs instead.
// IUniswapV2FactoryFuncSigs maps the 4-byte function signature to its string representation.
var IUniswapV2FactoryFuncSigs = IUniswapV2FactoryMetaData.Sigs

// IUniswapV2Factory is an auto generated Go binding around a Kaia contract.
type IUniswapV2Factory struct {
	IUniswapV2FactoryCaller     // Read-only binding to the contract
	IUniswapV2FactoryTransactor // Write-only binding to the contract
	IUniswapV2FactoryFilterer   // Log filterer for contract events
}

// IUniswapV2FactoryCaller is an auto generated read-only Go binding around a Kaia contract.
type IUniswapV2FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2FactoryTransactor is an auto generated write-only Go binding around a Kaia contract.
type IUniswapV2FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2FactoryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IUniswapV2FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2FactorySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IUniswapV2FactorySession struct {
	Contract     *IUniswapV2Factory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IUniswapV2FactoryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IUniswapV2FactoryCallerSession struct {
	Contract *IUniswapV2FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IUniswapV2FactoryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IUniswapV2FactoryTransactorSession struct {
	Contract     *IUniswapV2FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IUniswapV2FactoryRaw is an auto generated low-level Go binding around a Kaia contract.
type IUniswapV2FactoryRaw struct {
	Contract *IUniswapV2Factory // Generic contract binding to access the raw methods on
}

// IUniswapV2FactoryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IUniswapV2FactoryCallerRaw struct {
	Contract *IUniswapV2FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV2FactoryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IUniswapV2FactoryTransactorRaw struct {
	Contract *IUniswapV2FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV2Factory creates a new instance of IUniswapV2Factory, bound to a specific deployed contract.
func NewIUniswapV2Factory(address common.Address, backend bind.ContractBackend) (*IUniswapV2Factory, error) {
	contract, err := bindIUniswapV2Factory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2Factory{IUniswapV2FactoryCaller: IUniswapV2FactoryCaller{contract: contract}, IUniswapV2FactoryTransactor: IUniswapV2FactoryTransactor{contract: contract}, IUniswapV2FactoryFilterer: IUniswapV2FactoryFilterer{contract: contract}}, nil
}

// NewIUniswapV2FactoryCaller creates a new read-only instance of IUniswapV2Factory, bound to a specific deployed contract.
func NewIUniswapV2FactoryCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV2FactoryCaller, error) {
	contract, err := bindIUniswapV2Factory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2FactoryCaller{contract: contract}, nil
}

// NewIUniswapV2FactoryTransactor creates a new write-only instance of IUniswapV2Factory, bound to a specific deployed contract.
func NewIUniswapV2FactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV2FactoryTransactor, error) {
	contract, err := bindIUniswapV2Factory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2FactoryTransactor{contract: contract}, nil
}

// NewIUniswapV2FactoryFilterer creates a new log filterer instance of IUniswapV2Factory, bound to a specific deployed contract.
func NewIUniswapV2FactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV2FactoryFilterer, error) {
	contract, err := bindIUniswapV2Factory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2FactoryFilterer{contract: contract}, nil
}

// bindIUniswapV2Factory binds a generic wrapper to an already deployed contract.
func bindIUniswapV2Factory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV2FactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Factory *IUniswapV2FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Factory.Contract.IUniswapV2FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Factory *IUniswapV2FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.IUniswapV2FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Factory *IUniswapV2FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.IUniswapV2FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Factory *IUniswapV2FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.contract.Transact(opts, method, params...)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) AllPairs(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "allPairs", arg0)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactorySession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _IUniswapV2Factory.Contract.AllPairs(&_IUniswapV2Factory.CallOpts, arg0)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _IUniswapV2Factory.Contract.AllPairs(&_IUniswapV2Factory.CallOpts, arg0)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) AllPairsLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "allPairsLength")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_IUniswapV2Factory *IUniswapV2FactorySession) AllPairsLength() (*big.Int, error) {
	return _IUniswapV2Factory.Contract.AllPairsLength(&_IUniswapV2Factory.CallOpts)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) AllPairsLength() (*big.Int, error) {
	return _IUniswapV2Factory.Contract.AllPairsLength(&_IUniswapV2Factory.CallOpts)
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) FeeTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "feeTo")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactorySession) FeeTo() (common.Address, error) {
	return _IUniswapV2Factory.Contract.FeeTo(&_IUniswapV2Factory.CallOpts)
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) FeeTo() (common.Address, error) {
	return _IUniswapV2Factory.Contract.FeeTo(&_IUniswapV2Factory.CallOpts)
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) FeeToSetter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "feeToSetter")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactorySession) FeeToSetter() (common.Address, error) {
	return _IUniswapV2Factory.Contract.FeeToSetter(&_IUniswapV2Factory.CallOpts)
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) FeeToSetter() (common.Address, error) {
	return _IUniswapV2Factory.Contract.FeeToSetter(&_IUniswapV2Factory.CallOpts)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address tokenA, address tokenB) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryCaller) GetPair(opts *bind.CallOpts, tokenA common.Address, tokenB common.Address) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Factory.contract.Call(opts, &out, "getPair", tokenA, tokenB)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address tokenA, address tokenB) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactorySession) GetPair(tokenA common.Address, tokenB common.Address) (common.Address, error) {
	return _IUniswapV2Factory.Contract.GetPair(&_IUniswapV2Factory.CallOpts, tokenA, tokenB)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address tokenA, address tokenB) view returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryCallerSession) GetPair(tokenA common.Address, tokenB common.Address) (common.Address, error) {
	return _IUniswapV2Factory.Contract.GetPair(&_IUniswapV2Factory.CallOpts, tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryTransactor) CreatePair(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.contract.Transact(opts, "createPair", tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactorySession) CreatePair(tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.CreatePair(&_IUniswapV2Factory.TransactOpts, tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorSession) CreatePair(tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.CreatePair(&_IUniswapV2Factory.TransactOpts, tokenA, tokenB)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactoryTransactor) SetFeeTo(opts *bind.TransactOpts, arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.contract.Transact(opts, "setFeeTo", arg0)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactorySession) SetFeeTo(arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.SetFeeTo(&_IUniswapV2Factory.TransactOpts, arg0)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorSession) SetFeeTo(arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.SetFeeTo(&_IUniswapV2Factory.TransactOpts, arg0)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactoryTransactor) SetFeeToSetter(opts *bind.TransactOpts, arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.contract.Transact(opts, "setFeeToSetter", arg0)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactorySession) SetFeeToSetter(arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.SetFeeToSetter(&_IUniswapV2Factory.TransactOpts, arg0)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address ) returns()
func (_IUniswapV2Factory *IUniswapV2FactoryTransactorSession) SetFeeToSetter(arg0 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Factory.Contract.SetFeeToSetter(&_IUniswapV2Factory.TransactOpts, arg0)
}

// IUniswapV2FactoryPairCreatedIterator is returned from FilterPairCreated and is used to iterate over the raw logs and unpacked data for PairCreated events raised by the IUniswapV2Factory contract.
type IUniswapV2FactoryPairCreatedIterator struct {
	Event *IUniswapV2FactoryPairCreated // Event containing the contract specifics and raw log

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
func (it *IUniswapV2FactoryPairCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2FactoryPairCreated)
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
		it.Event = new(IUniswapV2FactoryPairCreated)
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
func (it *IUniswapV2FactoryPairCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2FactoryPairCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2FactoryPairCreated represents a PairCreated event raised by the IUniswapV2Factory contract.
type IUniswapV2FactoryPairCreated struct {
	Token0 common.Address
	Token1 common.Address
	Pair   common.Address
	Arg3   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPairCreated is a free log retrieval operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_IUniswapV2Factory *IUniswapV2FactoryFilterer) FilterPairCreated(opts *bind.FilterOpts, token0 []common.Address, token1 []common.Address) (*IUniswapV2FactoryPairCreatedIterator, error) {
	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _IUniswapV2Factory.contract.FilterLogs(opts, "PairCreated", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2FactoryPairCreatedIterator{contract: _IUniswapV2Factory.contract, event: "PairCreated", logs: logs, sub: sub}, nil
}

// WatchPairCreated is a free log subscription operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_IUniswapV2Factory *IUniswapV2FactoryFilterer) WatchPairCreated(opts *bind.WatchOpts, sink chan<- *IUniswapV2FactoryPairCreated, token0 []common.Address, token1 []common.Address) (event.Subscription, error) {
	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _IUniswapV2Factory.contract.WatchLogs(opts, "PairCreated", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2FactoryPairCreated)
				if err := _IUniswapV2Factory.contract.UnpackLog(event, "PairCreated", log); err != nil {
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

// ParsePairCreated is a log parse operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_IUniswapV2Factory *IUniswapV2FactoryFilterer) ParsePairCreated(log types.Log) (*IUniswapV2FactoryPairCreated, error) {
	event := new(IUniswapV2FactoryPairCreated)
	if err := _IUniswapV2Factory.contract.UnpackLog(event, "PairCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2PairMetaData contains all meta data concerning the IUniswapV2Pair contract.
var IUniswapV2PairMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint112\",\"name\":\"reserve0\",\"type\":\"uint112\"},{\"indexed\":false,\"internalType\":\"uint112\",\"name\":\"reserve1\",\"type\":\"uint112\"}],\"name\":\"Sync\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MINIMUM_LIQUIDITY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint112\",\"name\":\"reserve0\",\"type\":\"uint112\"},{\"internalType\":\"uint112\",\"name\":\"reserve1\",\"type\":\"uint112\"},{\"internalType\":\"uint32\",\"name\":\"blockTimestampLast\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"kLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"price0CumulativeLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"price1CumulativeLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"skim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"sync\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token0\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token1\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"3644e515": "DOMAIN_SEPARATOR()",
		"ba9a7a56": "MINIMUM_LIQUIDITY()",
		"30adf81f": "PERMIT_TYPEHASH()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"89afcb44": "burn(address)",
		"313ce567": "decimals()",
		"c45a0155": "factory()",
		"0902f1ac": "getReserves()",
		"485cc955": "initialize(address,address)",
		"7464fc3d": "kLast()",
		"6a627842": "mint(address)",
		"06fdde03": "name()",
		"7ecebe00": "nonces(address)",
		"d505accf": "permit(address,address,uint256,uint256,uint8,bytes32,bytes32)",
		"5909c0d5": "price0CumulativeLast()",
		"5a3d5493": "price1CumulativeLast()",
		"bc25cf77": "skim(address)",
		"022c0d9f": "swap(uint256,uint256,address,bytes)",
		"95d89b41": "symbol()",
		"fff6cae9": "sync()",
		"0dfe1681": "token0()",
		"d21220a7": "token1()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
	},
}

// IUniswapV2PairABI is the input ABI used to generate the binding from.
// Deprecated: Use IUniswapV2PairMetaData.ABI instead.
var IUniswapV2PairABI = IUniswapV2PairMetaData.ABI

// IUniswapV2PairBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IUniswapV2PairBinRuntime = ``

// Deprecated: Use IUniswapV2PairMetaData.Sigs instead.
// IUniswapV2PairFuncSigs maps the 4-byte function signature to its string representation.
var IUniswapV2PairFuncSigs = IUniswapV2PairMetaData.Sigs

// IUniswapV2Pair is an auto generated Go binding around a Kaia contract.
type IUniswapV2Pair struct {
	IUniswapV2PairCaller     // Read-only binding to the contract
	IUniswapV2PairTransactor // Write-only binding to the contract
	IUniswapV2PairFilterer   // Log filterer for contract events
}

// IUniswapV2PairCaller is an auto generated read-only Go binding around a Kaia contract.
type IUniswapV2PairCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2PairTransactor is an auto generated write-only Go binding around a Kaia contract.
type IUniswapV2PairTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2PairFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IUniswapV2PairFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IUniswapV2PairSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IUniswapV2PairSession struct {
	Contract     *IUniswapV2Pair   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IUniswapV2PairCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IUniswapV2PairCallerSession struct {
	Contract *IUniswapV2PairCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IUniswapV2PairTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IUniswapV2PairTransactorSession struct {
	Contract     *IUniswapV2PairTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IUniswapV2PairRaw is an auto generated low-level Go binding around a Kaia contract.
type IUniswapV2PairRaw struct {
	Contract *IUniswapV2Pair // Generic contract binding to access the raw methods on
}

// IUniswapV2PairCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IUniswapV2PairCallerRaw struct {
	Contract *IUniswapV2PairCaller // Generic read-only contract binding to access the raw methods on
}

// IUniswapV2PairTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IUniswapV2PairTransactorRaw struct {
	Contract *IUniswapV2PairTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIUniswapV2Pair creates a new instance of IUniswapV2Pair, bound to a specific deployed contract.
func NewIUniswapV2Pair(address common.Address, backend bind.ContractBackend) (*IUniswapV2Pair, error) {
	contract, err := bindIUniswapV2Pair(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2Pair{IUniswapV2PairCaller: IUniswapV2PairCaller{contract: contract}, IUniswapV2PairTransactor: IUniswapV2PairTransactor{contract: contract}, IUniswapV2PairFilterer: IUniswapV2PairFilterer{contract: contract}}, nil
}

// NewIUniswapV2PairCaller creates a new read-only instance of IUniswapV2Pair, bound to a specific deployed contract.
func NewIUniswapV2PairCaller(address common.Address, caller bind.ContractCaller) (*IUniswapV2PairCaller, error) {
	contract, err := bindIUniswapV2Pair(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairCaller{contract: contract}, nil
}

// NewIUniswapV2PairTransactor creates a new write-only instance of IUniswapV2Pair, bound to a specific deployed contract.
func NewIUniswapV2PairTransactor(address common.Address, transactor bind.ContractTransactor) (*IUniswapV2PairTransactor, error) {
	contract, err := bindIUniswapV2Pair(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairTransactor{contract: contract}, nil
}

// NewIUniswapV2PairFilterer creates a new log filterer instance of IUniswapV2Pair, bound to a specific deployed contract.
func NewIUniswapV2PairFilterer(address common.Address, filterer bind.ContractFilterer) (*IUniswapV2PairFilterer, error) {
	contract, err := bindIUniswapV2Pair(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairFilterer{contract: contract}, nil
}

// bindIUniswapV2Pair binds a generic wrapper to an already deployed contract.
func bindIUniswapV2Pair(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IUniswapV2PairMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Pair *IUniswapV2PairRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Pair.Contract.IUniswapV2PairCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Pair *IUniswapV2PairRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.IUniswapV2PairTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Pair *IUniswapV2PairRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.IUniswapV2PairTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IUniswapV2Pair *IUniswapV2PairCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IUniswapV2Pair.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IUniswapV2Pair *IUniswapV2PairTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IUniswapV2Pair *IUniswapV2PairTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_IUniswapV2Pair *IUniswapV2PairCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "DOMAIN_SEPARATOR")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_IUniswapV2Pair *IUniswapV2PairSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _IUniswapV2Pair.Contract.DOMAINSEPARATOR(&_IUniswapV2Pair.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _IUniswapV2Pair.Contract.DOMAINSEPARATOR(&_IUniswapV2Pair.CallOpts)
}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() pure returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCaller) MINIMUMLIQUIDITY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "MINIMUM_LIQUIDITY")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() pure returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairSession) MINIMUMLIQUIDITY() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.MINIMUMLIQUIDITY(&_IUniswapV2Pair.CallOpts)
}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() pure returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) MINIMUMLIQUIDITY() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.MINIMUMLIQUIDITY(&_IUniswapV2Pair.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() pure returns(bytes32)
func (_IUniswapV2Pair *IUniswapV2PairCaller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "PERMIT_TYPEHASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() pure returns(bytes32)
func (_IUniswapV2Pair *IUniswapV2PairSession) PERMITTYPEHASH() ([32]byte, error) {
	return _IUniswapV2Pair.Contract.PERMITTYPEHASH(&_IUniswapV2Pair.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() pure returns(bytes32)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _IUniswapV2Pair.Contract.PERMITTYPEHASH(&_IUniswapV2Pair.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "allowance", owner, spender)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IUniswapV2Pair.Contract.Allowance(&_IUniswapV2Pair.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IUniswapV2Pair.Contract.Allowance(&_IUniswapV2Pair.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "balanceOf", owner)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IUniswapV2Pair.Contract.BalanceOf(&_IUniswapV2Pair.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _IUniswapV2Pair.Contract.BalanceOf(&_IUniswapV2Pair.CallOpts, owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "decimals")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_IUniswapV2Pair *IUniswapV2PairSession) Decimals() (uint8, error) {
	return _IUniswapV2Pair.Contract.Decimals(&_IUniswapV2Pair.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Decimals() (uint8, error) {
	return _IUniswapV2Pair.Contract.Decimals(&_IUniswapV2Pair.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "factory")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairSession) Factory() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Factory(&_IUniswapV2Pair.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Factory() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Factory(&_IUniswapV2Pair.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_IUniswapV2Pair *IUniswapV2PairCaller) GetReserves(opts *bind.CallOpts) (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error,
) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "getReserves")

	outstruct := new(struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reserve0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestampLast = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_IUniswapV2Pair *IUniswapV2PairSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error,
) {
	return _IUniswapV2Pair.Contract.GetReserves(&_IUniswapV2Pair.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error,
) {
	return _IUniswapV2Pair.Contract.GetReserves(&_IUniswapV2Pair.CallOpts)
}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCaller) KLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "kLast")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairSession) KLast() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.KLast(&_IUniswapV2Pair.CallOpts)
}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) KLast() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.KLast(&_IUniswapV2Pair.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "name")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_IUniswapV2Pair *IUniswapV2PairSession) Name() (string, error) {
	return _IUniswapV2Pair.Contract.Name(&_IUniswapV2Pair.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() pure returns(string)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Name() (string, error) {
	return _IUniswapV2Pair.Contract.Name(&_IUniswapV2Pair.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Nonces(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "nonces", owner)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairSession) Nonces(owner common.Address) (*big.Int, error) {
	return _IUniswapV2Pair.Contract.Nonces(&_IUniswapV2Pair.CallOpts, owner)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address owner) view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Nonces(owner common.Address) (*big.Int, error) {
	return _IUniswapV2Pair.Contract.Nonces(&_IUniswapV2Pair.CallOpts, owner)
}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Price0CumulativeLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "price0CumulativeLast")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairSession) Price0CumulativeLast() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.Price0CumulativeLast(&_IUniswapV2Pair.CallOpts)
}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Price0CumulativeLast() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.Price0CumulativeLast(&_IUniswapV2Pair.CallOpts)
}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Price1CumulativeLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "price1CumulativeLast")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairSession) Price1CumulativeLast() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.Price1CumulativeLast(&_IUniswapV2Pair.CallOpts)
}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Price1CumulativeLast() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.Price1CumulativeLast(&_IUniswapV2Pair.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "symbol")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_IUniswapV2Pair *IUniswapV2PairSession) Symbol() (string, error) {
	return _IUniswapV2Pair.Contract.Symbol(&_IUniswapV2Pair.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() pure returns(string)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Symbol() (string, error) {
	return _IUniswapV2Pair.Contract.Symbol(&_IUniswapV2Pair.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Token0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "token0")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairSession) Token0() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Token0(&_IUniswapV2Pair.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Token0() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Token0(&_IUniswapV2Pair.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCaller) Token1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "token1")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairSession) Token1() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Token1(&_IUniswapV2Pair.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) Token1() (common.Address, error) {
	return _IUniswapV2Pair.Contract.Token1(&_IUniswapV2Pair.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IUniswapV2Pair.contract.Call(opts, &out, "totalSupply")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairSession) TotalSupply() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.TotalSupply(&_IUniswapV2Pair.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IUniswapV2Pair *IUniswapV2PairCallerSession) TotalSupply() (*big.Int, error) {
	return _IUniswapV2Pair.Contract.TotalSupply(&_IUniswapV2Pair.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Approve(&_IUniswapV2Pair.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Approve(&_IUniswapV2Pair.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Burn(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "burn", to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV2Pair *IUniswapV2PairSession) Burn(to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Burn(&_IUniswapV2Pair.TransactOpts, to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Burn(to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Burn(&_IUniswapV2Pair.TransactOpts, to)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address , address ) returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Initialize(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "initialize", arg0, arg1)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address , address ) returns()
func (_IUniswapV2Pair *IUniswapV2PairSession) Initialize(arg0 common.Address, arg1 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Initialize(&_IUniswapV2Pair.TransactOpts, arg0, arg1)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address , address ) returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Initialize(arg0 common.Address, arg1 common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Initialize(&_IUniswapV2Pair.TransactOpts, arg0, arg1)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Mint(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "mint", to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_IUniswapV2Pair *IUniswapV2PairSession) Mint(to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Mint(&_IUniswapV2Pair.TransactOpts, to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Mint(to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Mint(&_IUniswapV2Pair.TransactOpts, to)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_IUniswapV2Pair *IUniswapV2PairSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Permit(&_IUniswapV2Pair.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Permit(&_IUniswapV2Pair.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Skim(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "skim", to)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_IUniswapV2Pair *IUniswapV2PairSession) Skim(to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Skim(&_IUniswapV2Pair.TransactOpts, to)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Skim(to common.Address) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Skim(&_IUniswapV2Pair.TransactOpts, to)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Swap(opts *bind.TransactOpts, amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "swap", amount0Out, amount1Out, to, data)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_IUniswapV2Pair *IUniswapV2PairSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Swap(&_IUniswapV2Pair.TransactOpts, amount0Out, amount1Out, to, data)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Swap(&_IUniswapV2Pair.TransactOpts, amount0Out, amount1Out, to, data)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Sync(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "sync")
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_IUniswapV2Pair *IUniswapV2PairSession) Sync() (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Sync(&_IUniswapV2Pair.TransactOpts)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Sync() (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Sync(&_IUniswapV2Pair.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Transfer(&_IUniswapV2Pair.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.Transfer(&_IUniswapV2Pair.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.TransferFrom(&_IUniswapV2Pair.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IUniswapV2Pair *IUniswapV2PairTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IUniswapV2Pair.Contract.TransferFrom(&_IUniswapV2Pair.TransactOpts, from, to, value)
}

// IUniswapV2PairApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IUniswapV2Pair contract.
type IUniswapV2PairApprovalIterator struct {
	Event *IUniswapV2PairApproval // Event containing the contract specifics and raw log

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
func (it *IUniswapV2PairApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2PairApproval)
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
		it.Event = new(IUniswapV2PairApproval)
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
func (it *IUniswapV2PairApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2PairApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2PairApproval represents a Approval event raised by the IUniswapV2Pair contract.
type IUniswapV2PairApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IUniswapV2PairApprovalIterator, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairApprovalIterator{contract: _IUniswapV2Pair.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IUniswapV2PairApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2PairApproval)
				if err := _IUniswapV2Pair.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) ParseApproval(log types.Log) (*IUniswapV2PairApproval, error) {
	event := new(IUniswapV2PairApproval)
	if err := _IUniswapV2Pair.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2PairBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the IUniswapV2Pair contract.
type IUniswapV2PairBurnIterator struct {
	Event *IUniswapV2PairBurn // Event containing the contract specifics and raw log

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
func (it *IUniswapV2PairBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2PairBurn)
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
		it.Event = new(IUniswapV2PairBurn)
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
func (it *IUniswapV2PairBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2PairBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2PairBurn represents a Burn event raised by the IUniswapV2Pair contract.
type IUniswapV2PairBurn struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	To      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) FilterBurn(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*IUniswapV2PairBurnIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.FilterLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairBurnIterator{contract: _IUniswapV2Pair.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *IUniswapV2PairBurn, sender []common.Address, to []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.WatchLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2PairBurn)
				if err := _IUniswapV2Pair.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) ParseBurn(log types.Log) (*IUniswapV2PairBurn, error) {
	event := new(IUniswapV2PairBurn)
	if err := _IUniswapV2Pair.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2PairMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the IUniswapV2Pair contract.
type IUniswapV2PairMintIterator struct {
	Event *IUniswapV2PairMint // Event containing the contract specifics and raw log

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
func (it *IUniswapV2PairMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2PairMint)
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
		it.Event = new(IUniswapV2PairMint)
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
func (it *IUniswapV2PairMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2PairMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2PairMint represents a Mint event raised by the IUniswapV2Pair contract.
type IUniswapV2PairMint struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) FilterMint(opts *bind.FilterOpts, sender []common.Address) (*IUniswapV2PairMintIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.FilterLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairMintIterator{contract: _IUniswapV2Pair.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *IUniswapV2PairMint, sender []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.WatchLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2PairMint)
				if err := _IUniswapV2Pair.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) ParseMint(log types.Log) (*IUniswapV2PairMint, error) {
	event := new(IUniswapV2PairMint)
	if err := _IUniswapV2Pair.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2PairSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the IUniswapV2Pair contract.
type IUniswapV2PairSwapIterator struct {
	Event *IUniswapV2PairSwap // Event containing the contract specifics and raw log

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
func (it *IUniswapV2PairSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2PairSwap)
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
		it.Event = new(IUniswapV2PairSwap)
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
func (it *IUniswapV2PairSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2PairSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2PairSwap represents a Swap event raised by the IUniswapV2Pair contract.
type IUniswapV2PairSwap struct {
	Sender     common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) FilterSwap(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*IUniswapV2PairSwapIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.FilterLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairSwapIterator{contract: _IUniswapV2Pair.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *IUniswapV2PairSwap, sender []common.Address, to []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.WatchLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2PairSwap)
				if err := _IUniswapV2Pair.contract.UnpackLog(event, "Swap", log); err != nil {
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

// ParseSwap is a log parse operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) ParseSwap(log types.Log) (*IUniswapV2PairSwap, error) {
	event := new(IUniswapV2PairSwap)
	if err := _IUniswapV2Pair.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2PairSyncIterator is returned from FilterSync and is used to iterate over the raw logs and unpacked data for Sync events raised by the IUniswapV2Pair contract.
type IUniswapV2PairSyncIterator struct {
	Event *IUniswapV2PairSync // Event containing the contract specifics and raw log

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
func (it *IUniswapV2PairSyncIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2PairSync)
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
		it.Event = new(IUniswapV2PairSync)
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
func (it *IUniswapV2PairSyncIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2PairSyncIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2PairSync represents a Sync event raised by the IUniswapV2Pair contract.
type IUniswapV2PairSync struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSync is a free log retrieval operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) FilterSync(opts *bind.FilterOpts) (*IUniswapV2PairSyncIterator, error) {
	logs, sub, err := _IUniswapV2Pair.contract.FilterLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairSyncIterator{contract: _IUniswapV2Pair.contract, event: "Sync", logs: logs, sub: sub}, nil
}

// WatchSync is a free log subscription operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) WatchSync(opts *bind.WatchOpts, sink chan<- *IUniswapV2PairSync) (event.Subscription, error) {
	logs, sub, err := _IUniswapV2Pair.contract.WatchLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2PairSync)
				if err := _IUniswapV2Pair.contract.UnpackLog(event, "Sync", log); err != nil {
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

// ParseSync is a log parse operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) ParseSync(log types.Log) (*IUniswapV2PairSync, error) {
	event := new(IUniswapV2PairSync)
	if err := _IUniswapV2Pair.contract.UnpackLog(event, "Sync", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IUniswapV2PairTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IUniswapV2Pair contract.
type IUniswapV2PairTransferIterator struct {
	Event *IUniswapV2PairTransfer // Event containing the contract specifics and raw log

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
func (it *IUniswapV2PairTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IUniswapV2PairTransfer)
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
		it.Event = new(IUniswapV2PairTransfer)
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
func (it *IUniswapV2PairTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IUniswapV2PairTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IUniswapV2PairTransfer represents a Transfer event raised by the IUniswapV2Pair contract.
type IUniswapV2PairTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IUniswapV2PairTransferIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IUniswapV2PairTransferIterator{contract: _IUniswapV2Pair.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IUniswapV2PairTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IUniswapV2Pair.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IUniswapV2PairTransfer)
				if err := _IUniswapV2Pair.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IUniswapV2Pair *IUniswapV2PairFilterer) ParseTransfer(log types.Log) (*IUniswapV2PairTransfer, error) {
	event := new(IUniswapV2PairTransfer)
	if err := _IUniswapV2Pair.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MathMetaData contains all meta data concerning the Math contract.
var MathMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a7231582062d1272f0f40b69c2c03cbf7be07334fc5d55569598ee16efb9c83efad76f62764736f6c63430005100032",
}

// MathABI is the input ABI used to generate the binding from.
// Deprecated: Use MathMetaData.ABI instead.
var MathABI = MathMetaData.ABI

// MathBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MathBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a7231582062d1272f0f40b69c2c03cbf7be07334fc5d55569598ee16efb9c83efad76f62764736f6c63430005100032`

// MathBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MathMetaData.Bin instead.
var MathBin = MathMetaData.Bin

// DeployMath deploys a new Kaia contract, binding an instance of Math to it.
func DeployMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Math, error) {
	parsed, err := MathMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Math{MathCaller: MathCaller{contract: contract}, MathTransactor: MathTransactor{contract: contract}, MathFilterer: MathFilterer{contract: contract}}, nil
}

// Math is an auto generated Go binding around a Kaia contract.
type Math struct {
	MathCaller     // Read-only binding to the contract
	MathTransactor // Write-only binding to the contract
	MathFilterer   // Log filterer for contract events
}

// MathCaller is an auto generated read-only Go binding around a Kaia contract.
type MathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathTransactor is an auto generated write-only Go binding around a Kaia contract.
type MathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type MathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MathSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type MathSession struct {
	Contract     *Math             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MathCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type MathCallerSession struct {
	Contract *MathCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MathTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type MathTransactorSession struct {
	Contract     *MathTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MathRaw is an auto generated low-level Go binding around a Kaia contract.
type MathRaw struct {
	Contract *Math // Generic contract binding to access the raw methods on
}

// MathCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type MathCallerRaw struct {
	Contract *MathCaller // Generic read-only contract binding to access the raw methods on
}

// MathTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type MathTransactorRaw struct {
	Contract *MathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMath creates a new instance of Math, bound to a specific deployed contract.
func NewMath(address common.Address, backend bind.ContractBackend) (*Math, error) {
	contract, err := bindMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Math{MathCaller: MathCaller{contract: contract}, MathTransactor: MathTransactor{contract: contract}, MathFilterer: MathFilterer{contract: contract}}, nil
}

// NewMathCaller creates a new read-only instance of Math, bound to a specific deployed contract.
func NewMathCaller(address common.Address, caller bind.ContractCaller) (*MathCaller, error) {
	contract, err := bindMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MathCaller{contract: contract}, nil
}

// NewMathTransactor creates a new write-only instance of Math, bound to a specific deployed contract.
func NewMathTransactor(address common.Address, transactor bind.ContractTransactor) (*MathTransactor, error) {
	contract, err := bindMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MathTransactor{contract: contract}, nil
}

// NewMathFilterer creates a new log filterer instance of Math, bound to a specific deployed contract.
func NewMathFilterer(address common.Address, filterer bind.ContractFilterer) (*MathFilterer, error) {
	contract, err := bindMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MathFilterer{contract: contract}, nil
}

// bindMath binds a generic wrapper to an already deployed contract.
func bindMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MathMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Math *MathRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Math.Contract.MathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Math *MathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Math.Contract.MathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Math *MathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Math.Contract.MathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Math *MathCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Math.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Math *MathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Math.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Math *MathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Math.Contract.contract.Transact(opts, method, params...)
}

// SafeMathMetaData contains all meta data concerning the SafeMath contract.
var SafeMathMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a7231582000424e399d0e68f5b92db35ac72a2d44c433cfb7e62585f42d9dc9ed7971755c64736f6c63430005100032",
}

// SafeMathABI is the input ABI used to generate the binding from.
// Deprecated: Use SafeMathMetaData.ABI instead.
var SafeMathABI = SafeMathMetaData.ABI

// SafeMathBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const SafeMathBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a7231582000424e399d0e68f5b92db35ac72a2d44c433cfb7e62585f42d9dc9ed7971755c64736f6c63430005100032`

// SafeMathBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SafeMathMetaData.Bin instead.
var SafeMathBin = SafeMathMetaData.Bin

// DeploySafeMath deploys a new Kaia contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := SafeMathMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around a Kaia contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around a Kaia contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around a Kaia contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around a Kaia contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SafeMathMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// UQ112x112MetaData contains all meta data concerning the UQ112x112 contract.
var UQ112x112MetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556023600b82828239805160001a607314601657fe5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a723158206db49779fdee9ff00451921ee4875e5d07d94b5a3c6e3bbcee3c3d2c725b7b1264736f6c63430005100032",
}

// UQ112x112ABI is the input ABI used to generate the binding from.
// Deprecated: Use UQ112x112MetaData.ABI instead.
var UQ112x112ABI = UQ112x112MetaData.ABI

// UQ112x112BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const UQ112x112BinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea265627a7a723158206db49779fdee9ff00451921ee4875e5d07d94b5a3c6e3bbcee3c3d2c725b7b1264736f6c63430005100032`

// UQ112x112Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UQ112x112MetaData.Bin instead.
var UQ112x112Bin = UQ112x112MetaData.Bin

// DeployUQ112x112 deploys a new Kaia contract, binding an instance of UQ112x112 to it.
func DeployUQ112x112(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UQ112x112, error) {
	parsed, err := UQ112x112MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UQ112x112Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UQ112x112{UQ112x112Caller: UQ112x112Caller{contract: contract}, UQ112x112Transactor: UQ112x112Transactor{contract: contract}, UQ112x112Filterer: UQ112x112Filterer{contract: contract}}, nil
}

// UQ112x112 is an auto generated Go binding around a Kaia contract.
type UQ112x112 struct {
	UQ112x112Caller     // Read-only binding to the contract
	UQ112x112Transactor // Write-only binding to the contract
	UQ112x112Filterer   // Log filterer for contract events
}

// UQ112x112Caller is an auto generated read-only Go binding around a Kaia contract.
type UQ112x112Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UQ112x112Transactor is an auto generated write-only Go binding around a Kaia contract.
type UQ112x112Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UQ112x112Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type UQ112x112Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UQ112x112Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type UQ112x112Session struct {
	Contract     *UQ112x112        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UQ112x112CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type UQ112x112CallerSession struct {
	Contract *UQ112x112Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// UQ112x112TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type UQ112x112TransactorSession struct {
	Contract     *UQ112x112Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// UQ112x112Raw is an auto generated low-level Go binding around a Kaia contract.
type UQ112x112Raw struct {
	Contract *UQ112x112 // Generic contract binding to access the raw methods on
}

// UQ112x112CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type UQ112x112CallerRaw struct {
	Contract *UQ112x112Caller // Generic read-only contract binding to access the raw methods on
}

// UQ112x112TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type UQ112x112TransactorRaw struct {
	Contract *UQ112x112Transactor // Generic write-only contract binding to access the raw methods on
}

// NewUQ112x112 creates a new instance of UQ112x112, bound to a specific deployed contract.
func NewUQ112x112(address common.Address, backend bind.ContractBackend) (*UQ112x112, error) {
	contract, err := bindUQ112x112(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UQ112x112{UQ112x112Caller: UQ112x112Caller{contract: contract}, UQ112x112Transactor: UQ112x112Transactor{contract: contract}, UQ112x112Filterer: UQ112x112Filterer{contract: contract}}, nil
}

// NewUQ112x112Caller creates a new read-only instance of UQ112x112, bound to a specific deployed contract.
func NewUQ112x112Caller(address common.Address, caller bind.ContractCaller) (*UQ112x112Caller, error) {
	contract, err := bindUQ112x112(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UQ112x112Caller{contract: contract}, nil
}

// NewUQ112x112Transactor creates a new write-only instance of UQ112x112, bound to a specific deployed contract.
func NewUQ112x112Transactor(address common.Address, transactor bind.ContractTransactor) (*UQ112x112Transactor, error) {
	contract, err := bindUQ112x112(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UQ112x112Transactor{contract: contract}, nil
}

// NewUQ112x112Filterer creates a new log filterer instance of UQ112x112, bound to a specific deployed contract.
func NewUQ112x112Filterer(address common.Address, filterer bind.ContractFilterer) (*UQ112x112Filterer, error) {
	contract, err := bindUQ112x112(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UQ112x112Filterer{contract: contract}, nil
}

// bindUQ112x112 binds a generic wrapper to an already deployed contract.
func bindUQ112x112(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UQ112x112MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UQ112x112 *UQ112x112Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UQ112x112.Contract.UQ112x112Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UQ112x112 *UQ112x112Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UQ112x112.Contract.UQ112x112Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UQ112x112 *UQ112x112Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UQ112x112.Contract.UQ112x112Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UQ112x112 *UQ112x112CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UQ112x112.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UQ112x112 *UQ112x112TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UQ112x112.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UQ112x112 *UQ112x112TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UQ112x112.Contract.contract.Transact(opts, method, params...)
}

// UniswapV2ERC20MetaData contains all meta data concerning the UniswapV2ERC20 contract.
var UniswapV2ERC20MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"3644e515": "DOMAIN_SEPARATOR()",
		"30adf81f": "PERMIT_TYPEHASH()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"313ce567": "decimals()",
		"06fdde03": "name()",
		"7ecebe00": "nonces(address)",
		"d505accf": "permit(address,address,uint256,uint256,uint8,bytes32,bytes32)",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
	},
	Bin: "0x608060405234801561001057600080fd5b50604051469080605261096c8239604080519182900360520182208282018252600a8352692ab734b9bbb0b8102b1960b11b6020938401528151808301835260018152603160f81b908401528151808401919091527fbfcc8ef98ffbf7b6c3fec7bf5185b566b9863e35a9d83acd49ad6824b5969738818301527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6606082015260808101949094523060a0808601919091528151808603909101815260c0909401905282519201919091206003555061087e806100ee6000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c80633644e5151161008c57806395d89b411161006657806395d89b411461025b578063a9059cbb14610263578063d505accf1461028f578063dd62ed3e146102e2576100cf565b80633644e5151461020757806370a082311461020f5780637ecebe0014610235576100cf565b806306fdde03146100d4578063095ea7b31461015157806318160ddd1461019157806323b872dd146101ab57806330adf81f146101e1578063313ce567146101e9575b600080fd5b6100dc610310565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101165781810151838201526020016100fe565b50505050905090810190601f1680156101435780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61017d6004803603604081101561016757600080fd5b506001600160a01b038135169060200135610336565b604080519115158252519081900360200190f35b61019961034d565b60408051918252519081900360200190f35b61017d600480360360608110156101c157600080fd5b506001600160a01b03813581169160208101359091169060400135610353565b6101996103ed565b6101f1610411565b6040805160ff9092168252519081900360200190f35b610199610416565b6101996004803603602081101561022557600080fd5b50356001600160a01b031661041c565b6101996004803603602081101561024b57600080fd5b50356001600160a01b031661042e565b6100dc610440565b61017d6004803603604081101561027957600080fd5b506001600160a01b038135169060200135610462565b6102e0600480360360e08110156102a557600080fd5b506001600160a01b03813581169160208101359091169060408101359060608101359060ff6080820135169060a08101359060c0013561046f565b005b610199600480360360408110156102f857600080fd5b506001600160a01b0381358116916020013516610671565b6040518060400160405280600a8152602001692ab734b9bbb0b8102b1960b11b81525081565b600061034333848461068e565b5060015b92915050565b60005481565b6001600160a01b0383166000908152600260209081526040808320338452909152812054600019146103d8576001600160a01b03841660009081526002602090815260408083203384529091529020546103b3908363ffffffff6106f016565b6001600160a01b03851660009081526002602090815260408083203384529091529020555b6103e3848484610740565b5060019392505050565b7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c981565b601281565b60035481565b60016020526000908152604090205481565b60046020526000908152604090205481565b604051806040016040528060068152602001652aa72496ab1960d11b81525081565b6000610343338484610740565b428410156104b9576040805162461bcd60e51b8152602060048201526012602482015271155b9a5cddd85c158c8e881156141254915160721b604482015290519081900360640190fd5b6003546001600160a01b0380891660008181526004602090815260408083208054600180820190925582517f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c98186015280840196909652958d166060860152608085018c905260a085019590955260c08085018b90528151808603909101815260e08501825280519083012061190160f01b6101008601526101028501969096526101228085019690965280518085039096018652610142840180825286519683019690962095839052610162840180825286905260ff89166101828501526101a284018890526101c28401879052519193926101e280820193601f1981019281900390910190855afa1580156105d4573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381161580159061060a5750886001600160a01b0316816001600160a01b0316145b61065b576040805162461bcd60e51b815260206004820152601c60248201527f556e697377617056323a20494e56414c49445f5349474e415455524500000000604482015290519081900360640190fd5b61066689898961068e565b505050505050505050565b600260209081526000928352604080842090915290825290205481565b6001600160a01b03808416600081815260026020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b80820382811115610347576040805162461bcd60e51b815260206004820152601560248201527464732d6d6174682d7375622d756e646572666c6f7760581b604482015290519081900360640190fd5b6001600160a01b038316600090815260016020526040902054610769908263ffffffff6106f016565b6001600160a01b03808516600090815260016020526040808220939093559084168152205461079e908263ffffffff6107fa16565b6001600160a01b0380841660008181526001602090815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b80820182811015610347576040805162461bcd60e51b815260206004820152601460248201527364732d6d6174682d6164642d6f766572666c6f7760601b604482015290519081900360640190fdfea265627a7a72315820c13018ae2ca64151e6ef76adb64d370e8f1acff8736bb21511f6adb4c73ddbb764736f6c63430005100032454950373132446f6d61696e28737472696e67206e616d652c737472696e672076657273696f6e2c75696e7432353620636861696e49642c6164647265737320766572696679696e67436f6e747261637429",
}

// UniswapV2ERC20ABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV2ERC20MetaData.ABI instead.
var UniswapV2ERC20ABI = UniswapV2ERC20MetaData.ABI

// UniswapV2ERC20BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const UniswapV2ERC20BinRuntime = `608060405234801561001057600080fd5b50600436106100cf5760003560e01c80633644e5151161008c57806395d89b411161006657806395d89b411461025b578063a9059cbb14610263578063d505accf1461028f578063dd62ed3e146102e2576100cf565b80633644e5151461020757806370a082311461020f5780637ecebe0014610235576100cf565b806306fdde03146100d4578063095ea7b31461015157806318160ddd1461019157806323b872dd146101ab57806330adf81f146101e1578063313ce567146101e9575b600080fd5b6100dc610310565b6040805160208082528351818301528351919283929083019185019080838360005b838110156101165781810151838201526020016100fe565b50505050905090810190601f1680156101435780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61017d6004803603604081101561016757600080fd5b506001600160a01b038135169060200135610336565b604080519115158252519081900360200190f35b61019961034d565b60408051918252519081900360200190f35b61017d600480360360608110156101c157600080fd5b506001600160a01b03813581169160208101359091169060400135610353565b6101996103ed565b6101f1610411565b6040805160ff9092168252519081900360200190f35b610199610416565b6101996004803603602081101561022557600080fd5b50356001600160a01b031661041c565b6101996004803603602081101561024b57600080fd5b50356001600160a01b031661042e565b6100dc610440565b61017d6004803603604081101561027957600080fd5b506001600160a01b038135169060200135610462565b6102e0600480360360e08110156102a557600080fd5b506001600160a01b03813581169160208101359091169060408101359060608101359060ff6080820135169060a08101359060c0013561046f565b005b610199600480360360408110156102f857600080fd5b506001600160a01b0381358116916020013516610671565b6040518060400160405280600a8152602001692ab734b9bbb0b8102b1960b11b81525081565b600061034333848461068e565b5060015b92915050565b60005481565b6001600160a01b0383166000908152600260209081526040808320338452909152812054600019146103d8576001600160a01b03841660009081526002602090815260408083203384529091529020546103b3908363ffffffff6106f016565b6001600160a01b03851660009081526002602090815260408083203384529091529020555b6103e3848484610740565b5060019392505050565b7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c981565b601281565b60035481565b60016020526000908152604090205481565b60046020526000908152604090205481565b604051806040016040528060068152602001652aa72496ab1960d11b81525081565b6000610343338484610740565b428410156104b9576040805162461bcd60e51b8152602060048201526012602482015271155b9a5cddd85c158c8e881156141254915160721b604482015290519081900360640190fd5b6003546001600160a01b0380891660008181526004602090815260408083208054600180820190925582517f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c98186015280840196909652958d166060860152608085018c905260a085019590955260c08085018b90528151808603909101815260e08501825280519083012061190160f01b6101008601526101028501969096526101228085019690965280518085039096018652610142840180825286519683019690962095839052610162840180825286905260ff89166101828501526101a284018890526101c28401879052519193926101e280820193601f1981019281900390910190855afa1580156105d4573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381161580159061060a5750886001600160a01b0316816001600160a01b0316145b61065b576040805162461bcd60e51b815260206004820152601c60248201527f556e697377617056323a20494e56414c49445f5349474e415455524500000000604482015290519081900360640190fd5b61066689898961068e565b505050505050505050565b600260209081526000928352604080842090915290825290205481565b6001600160a01b03808416600081815260026020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b80820382811115610347576040805162461bcd60e51b815260206004820152601560248201527464732d6d6174682d7375622d756e646572666c6f7760581b604482015290519081900360640190fd5b6001600160a01b038316600090815260016020526040902054610769908263ffffffff6106f016565b6001600160a01b03808516600090815260016020526040808220939093559084168152205461079e908263ffffffff6107fa16565b6001600160a01b0380841660008181526001602090815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b80820182811015610347576040805162461bcd60e51b815260206004820152601460248201527364732d6d6174682d6164642d6f766572666c6f7760601b604482015290519081900360640190fdfea265627a7a72315820c13018ae2ca64151e6ef76adb64d370e8f1acff8736bb21511f6adb4c73ddbb764736f6c63430005100032`

// Deprecated: Use UniswapV2ERC20MetaData.Sigs instead.
// UniswapV2ERC20FuncSigs maps the 4-byte function signature to its string representation.
var UniswapV2ERC20FuncSigs = UniswapV2ERC20MetaData.Sigs

// UniswapV2ERC20Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniswapV2ERC20MetaData.Bin instead.
var UniswapV2ERC20Bin = UniswapV2ERC20MetaData.Bin

// DeployUniswapV2ERC20 deploys a new Kaia contract, binding an instance of UniswapV2ERC20 to it.
func DeployUniswapV2ERC20(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UniswapV2ERC20, error) {
	parsed, err := UniswapV2ERC20MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniswapV2ERC20Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniswapV2ERC20{UniswapV2ERC20Caller: UniswapV2ERC20Caller{contract: contract}, UniswapV2ERC20Transactor: UniswapV2ERC20Transactor{contract: contract}, UniswapV2ERC20Filterer: UniswapV2ERC20Filterer{contract: contract}}, nil
}

// UniswapV2ERC20 is an auto generated Go binding around a Kaia contract.
type UniswapV2ERC20 struct {
	UniswapV2ERC20Caller     // Read-only binding to the contract
	UniswapV2ERC20Transactor // Write-only binding to the contract
	UniswapV2ERC20Filterer   // Log filterer for contract events
}

// UniswapV2ERC20Caller is an auto generated read-only Go binding around a Kaia contract.
type UniswapV2ERC20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2ERC20Transactor is an auto generated write-only Go binding around a Kaia contract.
type UniswapV2ERC20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2ERC20Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type UniswapV2ERC20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2ERC20Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type UniswapV2ERC20Session struct {
	Contract     *UniswapV2ERC20   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapV2ERC20CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type UniswapV2ERC20CallerSession struct {
	Contract *UniswapV2ERC20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// UniswapV2ERC20TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type UniswapV2ERC20TransactorSession struct {
	Contract     *UniswapV2ERC20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// UniswapV2ERC20Raw is an auto generated low-level Go binding around a Kaia contract.
type UniswapV2ERC20Raw struct {
	Contract *UniswapV2ERC20 // Generic contract binding to access the raw methods on
}

// UniswapV2ERC20CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type UniswapV2ERC20CallerRaw struct {
	Contract *UniswapV2ERC20Caller // Generic read-only contract binding to access the raw methods on
}

// UniswapV2ERC20TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type UniswapV2ERC20TransactorRaw struct {
	Contract *UniswapV2ERC20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV2ERC20 creates a new instance of UniswapV2ERC20, bound to a specific deployed contract.
func NewUniswapV2ERC20(address common.Address, backend bind.ContractBackend) (*UniswapV2ERC20, error) {
	contract, err := bindUniswapV2ERC20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV2ERC20{UniswapV2ERC20Caller: UniswapV2ERC20Caller{contract: contract}, UniswapV2ERC20Transactor: UniswapV2ERC20Transactor{contract: contract}, UniswapV2ERC20Filterer: UniswapV2ERC20Filterer{contract: contract}}, nil
}

// NewUniswapV2ERC20Caller creates a new read-only instance of UniswapV2ERC20, bound to a specific deployed contract.
func NewUniswapV2ERC20Caller(address common.Address, caller bind.ContractCaller) (*UniswapV2ERC20Caller, error) {
	contract, err := bindUniswapV2ERC20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2ERC20Caller{contract: contract}, nil
}

// NewUniswapV2ERC20Transactor creates a new write-only instance of UniswapV2ERC20, bound to a specific deployed contract.
func NewUniswapV2ERC20Transactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV2ERC20Transactor, error) {
	contract, err := bindUniswapV2ERC20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2ERC20Transactor{contract: contract}, nil
}

// NewUniswapV2ERC20Filterer creates a new log filterer instance of UniswapV2ERC20, bound to a specific deployed contract.
func NewUniswapV2ERC20Filterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV2ERC20Filterer, error) {
	contract, err := bindUniswapV2ERC20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV2ERC20Filterer{contract: contract}, nil
}

// bindUniswapV2ERC20 binds a generic wrapper to an already deployed contract.
func bindUniswapV2ERC20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV2ERC20MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2ERC20 *UniswapV2ERC20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2ERC20.Contract.UniswapV2ERC20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2ERC20 *UniswapV2ERC20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.UniswapV2ERC20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2ERC20 *UniswapV2ERC20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.UniswapV2ERC20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2ERC20 *UniswapV2ERC20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2ERC20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2ERC20 *UniswapV2ERC20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2ERC20 *UniswapV2ERC20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "DOMAIN_SEPARATOR")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) DOMAINSEPARATOR() ([32]byte, error) {
	return _UniswapV2ERC20.Contract.DOMAINSEPARATOR(&_UniswapV2ERC20.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _UniswapV2ERC20.Contract.DOMAINSEPARATOR(&_UniswapV2ERC20.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "PERMIT_TYPEHASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) PERMITTYPEHASH() ([32]byte, error) {
	return _UniswapV2ERC20.Contract.PERMITTYPEHASH(&_UniswapV2ERC20.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _UniswapV2ERC20.Contract.PERMITTYPEHASH(&_UniswapV2ERC20.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "allowance", arg0, arg1)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _UniswapV2ERC20.Contract.Allowance(&_UniswapV2ERC20.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _UniswapV2ERC20.Contract.Allowance(&_UniswapV2ERC20.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "balanceOf", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _UniswapV2ERC20.Contract.BalanceOf(&_UniswapV2ERC20.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _UniswapV2ERC20.Contract.BalanceOf(&_UniswapV2ERC20.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "decimals")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) Decimals() (uint8, error) {
	return _UniswapV2ERC20.Contract.Decimals(&_UniswapV2ERC20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) Decimals() (uint8, error) {
	return _UniswapV2ERC20.Contract.Decimals(&_UniswapV2ERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "name")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) Name() (string, error) {
	return _UniswapV2ERC20.Contract.Name(&_UniswapV2ERC20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) Name() (string, error) {
	return _UniswapV2ERC20.Contract.Name(&_UniswapV2ERC20.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "nonces", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) Nonces(arg0 common.Address) (*big.Int, error) {
	return _UniswapV2ERC20.Contract.Nonces(&_UniswapV2ERC20.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _UniswapV2ERC20.Contract.Nonces(&_UniswapV2ERC20.CallOpts, arg0)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "symbol")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) Symbol() (string, error) {
	return _UniswapV2ERC20.Contract.Symbol(&_UniswapV2ERC20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) Symbol() (string, error) {
	return _UniswapV2ERC20.Contract.Symbol(&_UniswapV2ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2ERC20.contract.Call(opts, &out, "totalSupply")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) TotalSupply() (*big.Int, error) {
	return _UniswapV2ERC20.Contract.TotalSupply(&_UniswapV2ERC20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_UniswapV2ERC20 *UniswapV2ERC20CallerSession) TotalSupply() (*big.Int, error) {
	return _UniswapV2ERC20.Contract.TotalSupply(&_UniswapV2ERC20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.Approve(&_UniswapV2ERC20.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20TransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.Approve(&_UniswapV2ERC20.TransactOpts, spender, value)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_UniswapV2ERC20 *UniswapV2ERC20Transactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _UniswapV2ERC20.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_UniswapV2ERC20 *UniswapV2ERC20Session) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.Permit(&_UniswapV2ERC20.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_UniswapV2ERC20 *UniswapV2ERC20TransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.Permit(&_UniswapV2ERC20.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20Transactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.Transfer(&_UniswapV2ERC20.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20TransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.Transfer(&_UniswapV2ERC20.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20Transactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20Session) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.TransferFrom(&_UniswapV2ERC20.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_UniswapV2ERC20 *UniswapV2ERC20TransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2ERC20.Contract.TransferFrom(&_UniswapV2ERC20.TransactOpts, from, to, value)
}

// UniswapV2ERC20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the UniswapV2ERC20 contract.
type UniswapV2ERC20ApprovalIterator struct {
	Event *UniswapV2ERC20Approval // Event containing the contract specifics and raw log

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
func (it *UniswapV2ERC20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2ERC20Approval)
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
		it.Event = new(UniswapV2ERC20Approval)
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
func (it *UniswapV2ERC20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2ERC20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2ERC20Approval represents a Approval event raised by the UniswapV2ERC20 contract.
type UniswapV2ERC20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UniswapV2ERC20 *UniswapV2ERC20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*UniswapV2ERC20ApprovalIterator, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _UniswapV2ERC20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV2ERC20ApprovalIterator{contract: _UniswapV2ERC20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UniswapV2ERC20 *UniswapV2ERC20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *UniswapV2ERC20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _UniswapV2ERC20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2ERC20Approval)
				if err := _UniswapV2ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UniswapV2ERC20 *UniswapV2ERC20Filterer) ParseApproval(log types.Log) (*UniswapV2ERC20Approval, error) {
	event := new(UniswapV2ERC20Approval)
	if err := _UniswapV2ERC20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV2ERC20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the UniswapV2ERC20 contract.
type UniswapV2ERC20TransferIterator struct {
	Event *UniswapV2ERC20Transfer // Event containing the contract specifics and raw log

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
func (it *UniswapV2ERC20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2ERC20Transfer)
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
		it.Event = new(UniswapV2ERC20Transfer)
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
func (it *UniswapV2ERC20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2ERC20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2ERC20Transfer represents a Transfer event raised by the UniswapV2ERC20 contract.
type UniswapV2ERC20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UniswapV2ERC20 *UniswapV2ERC20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UniswapV2ERC20TransferIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV2ERC20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV2ERC20TransferIterator{contract: _UniswapV2ERC20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UniswapV2ERC20 *UniswapV2ERC20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *UniswapV2ERC20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV2ERC20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2ERC20Transfer)
				if err := _UniswapV2ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UniswapV2ERC20 *UniswapV2ERC20Filterer) ParseTransfer(log types.Log) (*UniswapV2ERC20Transfer, error) {
	event := new(UniswapV2ERC20Transfer)
	if err := _UniswapV2ERC20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV2FactoryMetaData contains all meta data concerning the UniswapV2Factory contract.
var UniswapV2FactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeToSetter\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"PairCreated\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allPairs\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"allPairsLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"}],\"name\":\"createPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"pair\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"feeToSetter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"getPair\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeTo\",\"type\":\"address\"}],\"name\":\"setFeeTo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeToSetter\",\"type\":\"address\"}],\"name\":\"setFeeToSetter\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"1e3dd18b": "allPairs(uint256)",
		"574f2ba3": "allPairsLength()",
		"c9c65396": "createPair(address,address)",
		"017e7e58": "feeTo()",
		"094b7415": "feeToSetter()",
		"e6a43905": "getPair(address,address)",
		"f46901ed": "setFeeTo(address)",
		"a2e74af6": "setFeeToSetter(address)",
	},
	Bin: "0x608060405234801561001057600080fd5b506040516136863803806136868339818101604052602081101561003357600080fd5b5051600180546001600160a01b0319166001600160a01b03909216919091179055613623806100636000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c8063a2e74af61161005b578063a2e74af6146100fd578063c9c6539614610132578063e6a439051461016d578063f46901ed146101a857610088565b8063017e7e581461008d578063094b7415146100be5780631e3dd18b146100c6578063574f2ba3146100e3575b600080fd5b6100956101db565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100956101f7565b610095600480360360208110156100dc57600080fd5b5035610213565b6100eb610247565b60408051918252519081900360200190f35b6101306004803603602081101561011357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661024d565b005b6100956004803603604081101561014857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602001351661031a565b6100956004803603604081101561018357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602001351661076d565b610130600480360360208110156101be57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166107a0565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60015473ffffffffffffffffffffffffffffffffffffffff1681565b6003818154811061022057fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16905081565b60035490565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102d357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f556e697377617056323a20464f5242494444454e000000000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60008173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156103b757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f556e697377617056323a204944454e544943414c5f4144445245535345530000604482015290519081900360640190fd5b6000808373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16106103f45783856103f7565b84845b909250905073ffffffffffffffffffffffffffffffffffffffff821661047e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f556e697377617056323a205a45524f5f41444452455353000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff82811660009081526002602090815260408083208585168452909152902054161561051f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f556e697377617056323a20504149525f45584953545300000000000000000000604482015290519081900360640190fd5b6060604051806020016105319061086d565b6020820181038252601f19601f82011660405250905060008383604051602001808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b81526014018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b815260140192505050604051602081830303815290604052805190602001209050808251602084016000f5604080517f485cc95500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152868116602483015291519297509087169163485cc9559160448082019260009290919082900301818387803b15801561065e57600080fd5b505af1158015610672573d6000803e3d6000fd5b5050505073ffffffffffffffffffffffffffffffffffffffff84811660008181526002602081815260408084208987168086529083528185208054978d167fffffffffffffffffffffffff000000000000000000000000000000000000000098891681179091559383528185208686528352818520805488168517905560038054600181018255958190527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b90950180549097168417909655925483519283529082015281517f0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9929181900390910190a35050505092915050565b600260209081526000928352604080842090915290825290205473ffffffffffffffffffffffffffffffffffffffff1681565b60015473ffffffffffffffffffffffffffffffffffffffff16331461082657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f556e697377617056323a20464f5242494444454e000000000000000000000000604482015290519081900360640190fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b612d748061087b8339019056fe60806040526001600c5534801561001557600080fd5b506040514690806052612d228239604080519182900360520182208282018252600a8352692ab734b9bbb0b8102b1960b11b6020938401528151808301835260018152603160f81b908401528151808401919091527fbfcc8ef98ffbf7b6c3fec7bf5185b566b9863e35a9d83acd49ad6824b5969738818301527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6606082015260808101949094523060a0808601919091528151808603909101815260c09094019052825192019190912060035550600580546001600160a01b03191633179055612c1d806101056000396000f3fe608060405234801561001057600080fd5b50600436106101b95760003560e01c80636a627842116100f9578063ba9a7a5611610097578063d21220a711610071578063d21220a7146105da578063d505accf146105e2578063dd62ed3e14610640578063fff6cae91461067b576101b9565b8063ba9a7a5614610597578063bc25cf771461059f578063c45a0155146105d2576101b9565b80637ecebe00116100d35780637ecebe00146104d757806389afcb441461050a57806395d89b4114610556578063a9059cbb1461055e576101b9565b80636a6278421461046957806370a082311461049c5780637464fc3d146104cf576101b9565b806323b872dd116101665780633644e515116101405780633644e51514610416578063485cc9551461041e5780635909c0d5146104595780635a3d549314610461576101b9565b806323b872dd146103ad57806330adf81f146103f0578063313ce567146103f8576101b9565b8063095ea7b311610197578063095ea7b3146103155780630dfe16811461036257806318160ddd14610393576101b9565b8063022c0d9f146101be57806306fdde03146102595780630902f1ac146102d6575b600080fd5b610257600480360360808110156101d457600080fd5b81359160208101359173ffffffffffffffffffffffffffffffffffffffff604083013516919081019060808101606082013564010000000081111561021857600080fd5b82018360208201111561022a57600080fd5b8035906020019184600183028401116401000000008311171561024c57600080fd5b509092509050610683565b005b610261610d57565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561029b578181015183820152602001610283565b50505050905090810190601f1680156102c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102de610d90565b604080516dffffffffffffffffffffffffffff948516815292909316602083015263ffffffff168183015290519081900360600190f35b61034e6004803603604081101561032b57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610de5565b604080519115158252519081900360200190f35b61036a610dfc565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61039b610e18565b60408051918252519081900360200190f35b61034e600480360360608110156103c357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060400135610e1e565b61039b610efd565b610400610f21565b6040805160ff9092168252519081900360200190f35b61039b610f26565b6102576004803603604081101561043457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516610f2c565b61039b611005565b61039b61100b565b61039b6004803603602081101561047f57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611011565b61039b600480360360208110156104b257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113cb565b61039b6113dd565b61039b600480360360208110156104ed57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113e3565b61053d6004803603602081101561052057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113f5565b6040805192835260208301919091528051918290030190f35b610261611892565b61034e6004803603604081101561057457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356118cb565b61039b6118d8565b610257600480360360208110156105b557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166118de565b61036a611ad4565b61036a611af0565b610257600480360360e08110156105f857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060408101359060608101359060ff6080820135169060a08101359060c00135611b0c565b61039b6004803603604081101561065657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516611dd8565b610257611df5565b600c546001146106f457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55841515806107075750600084115b61075c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180612b2f6025913960400191505060405180910390fd5b600080610767610d90565b5091509150816dffffffffffffffffffffffffffff168710801561079a5750806dffffffffffffffffffffffffffff1686105b6107ef576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180612b786021913960400191505060405180910390fd5b600654600754600091829173ffffffffffffffffffffffffffffffffffffffff91821691908116908916821480159061085457508073ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff1614155b6108bf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f556e697377617056323a20494e56414c49445f544f0000000000000000000000604482015290519081900360640190fd5b8a156108d0576108d0828a8d611fdb565b89156108e1576108e1818a8c611fdb565b86156109c3578873ffffffffffffffffffffffffffffffffffffffff166310d1e85c338d8d8c8c6040518663ffffffff1660e01b8152600401808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001858152602001848152602001806020018281038252848482818152602001925080828437600081840152601f19601f8201169050808301925050509650505050505050600060405180830381600087803b1580156109aa57600080fd5b505af11580156109be573d6000803e3d6000fd5b505050505b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff8416916370a08231916024808301926020929190829003018186803b158015610a2f57600080fd5b505afa158015610a43573d6000803e3d6000fd5b505050506040513d6020811015610a5957600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191955073ffffffffffffffffffffffffffffffffffffffff8316916370a0823191602480820192602092909190829003018186803b158015610acb57600080fd5b505afa158015610adf573d6000803e3d6000fd5b505050506040513d6020811015610af557600080fd5b5051925060009150506dffffffffffffffffffffffffffff85168a90038311610b1f576000610b35565b89856dffffffffffffffffffffffffffff160383035b9050600089856dffffffffffffffffffffffffffff16038311610b59576000610b6f565b89856dffffffffffffffffffffffffffff160383035b90506000821180610b805750600081115b610bd5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526024815260200180612b546024913960400191505060405180910390fd5b6000610c09610beb84600363ffffffff6121e816565b610bfd876103e863ffffffff6121e816565b9063ffffffff61226e16565b90506000610c21610beb84600363ffffffff6121e816565b9050610c59620f4240610c4d6dffffffffffffffffffffffffffff8b8116908b1663ffffffff6121e816565b9063ffffffff6121e816565b610c69838363ffffffff6121e816565b1015610cd657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f556e697377617056323a204b0000000000000000000000000000000000000000604482015290519081900360640190fd5b5050610ce4848488886122e0565b60408051838152602081018390528082018d9052606081018c9052905173ffffffffffffffffffffffffffffffffffffffff8b169133917fd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d8229181900360800190a350506001600c55505050505050505050565b6040518060400160405280600a81526020017f556e69737761702056320000000000000000000000000000000000000000000081525081565b6008546dffffffffffffffffffffffffffff808216926e0100000000000000000000000000008304909116917c0100000000000000000000000000000000000000000000000000000000900463ffffffff1690565b6000610df233848461259c565b5060015b92915050565b60065473ffffffffffffffffffffffffffffffffffffffff1681565b60005481565b73ffffffffffffffffffffffffffffffffffffffff831660009081526002602090815260408083203384529091528120547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff14610ee85773ffffffffffffffffffffffffffffffffffffffff84166000908152600260209081526040808320338452909152902054610eb6908363ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff851660009081526002602090815260408083203384529091529020555b610ef384848461260b565b5060019392505050565b7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c981565b601281565b60035481565b60055473ffffffffffffffffffffffffffffffffffffffff163314610fb257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f556e697377617056323a20464f5242494444454e000000000000000000000000604482015290519081900360640190fd5b6006805473ffffffffffffffffffffffffffffffffffffffff9384167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560078054929093169116179055565b60095481565b600a5481565b6000600c5460011461108457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c81905580611094610d90565b50600654604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905193955091935060009273ffffffffffffffffffffffffffffffffffffffff909116916370a08231916024808301926020929190829003018186803b15801561110e57600080fd5b505afa158015611122573d6000803e3d6000fd5b505050506040513d602081101561113857600080fd5b5051600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905192935060009273ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b1580156111b157600080fd5b505afa1580156111c5573d6000803e3d6000fd5b505050506040513d60208110156111db57600080fd5b505190506000611201836dffffffffffffffffffffffffffff871663ffffffff61226e16565b90506000611225836dffffffffffffffffffffffffffff871663ffffffff61226e16565b9050600061123387876126ec565b600054909150806112705761125c6103e8610bfd611257878763ffffffff6121e816565b612878565b985061126b60006103e86128ca565b6112cd565b6112ca6dffffffffffffffffffffffffffff8916611294868463ffffffff6121e816565b8161129b57fe5b046dffffffffffffffffffffffffffff89166112bd868563ffffffff6121e816565b816112c457fe5b0461297a565b98505b60008911611326576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180612bc16028913960400191505060405180910390fd5b6113308a8a6128ca565b61133c86868a8a6122e0565b811561137e5760085461137a906dffffffffffffffffffffffffffff808216916e01000000000000000000000000000090041663ffffffff6121e816565b600b555b6040805185815260208101859052815133927f4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f928290030190a250506001600c5550949695505050505050565b60016020526000908152604090205481565b600b5481565b60046020526000908152604090205481565b600080600c5460011461146957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c81905580611479610d90565b50600654600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905194965092945073ffffffffffffffffffffffffffffffffffffffff9182169391169160009184916370a08231916024808301926020929190829003018186803b1580156114fb57600080fd5b505afa15801561150f573d6000803e3d6000fd5b505050506040513d602081101561152557600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191925060009173ffffffffffffffffffffffffffffffffffffffff8516916370a08231916024808301926020929190829003018186803b15801561159957600080fd5b505afa1580156115ad573d6000803e3d6000fd5b505050506040513d60208110156115c357600080fd5b5051306000908152600160205260408120549192506115e288886126ec565b600054909150806115f9848763ffffffff6121e816565b8161160057fe5b049a5080611614848663ffffffff6121e816565b8161161b57fe5b04995060008b11801561162e575060008a115b611683576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180612b996028913960400191505060405180910390fd5b61168d3084612992565b611698878d8d611fdb565b6116a3868d8c611fdb565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff8916916370a08231916024808301926020929190829003018186803b15801561170f57600080fd5b505afa158015611723573d6000803e3d6000fd5b505050506040513d602081101561173957600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191965073ffffffffffffffffffffffffffffffffffffffff8816916370a0823191602480820192602092909190829003018186803b1580156117ab57600080fd5b505afa1580156117bf573d6000803e3d6000fd5b505050506040513d60208110156117d557600080fd5b505193506117e585858b8b6122e0565b811561182757600854611823906dffffffffffffffffffffffffffff808216916e01000000000000000000000000000090041663ffffffff6121e816565b600b555b604080518c8152602081018c9052815173ffffffffffffffffffffffffffffffffffffffff8f169233927fdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496929081900390910190a35050505050505050506001600c81905550915091565b6040518060400160405280600681526020017f554e492d5632000000000000000000000000000000000000000000000000000081525081565b6000610df233848461260b565b6103e881565b600c5460011461194f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55600654600754600854604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff9485169490931692611a2b9285928792611a26926dffffffffffffffffffffffffffff169185916370a0823191602480820192602092909190829003018186803b1580156119ee57600080fd5b505afa158015611a02573d6000803e3d6000fd5b505050506040513d6020811015611a1857600080fd5b50519063ffffffff61226e16565b611fdb565b600854604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529051611aca9284928792611a26926e01000000000000000000000000000090046dffffffffffffffffffffffffffff169173ffffffffffffffffffffffffffffffffffffffff8616916370a0823191602480820192602092909190829003018186803b1580156119ee57600080fd5b50506001600c5550565b60055473ffffffffffffffffffffffffffffffffffffffff1681565b60075473ffffffffffffffffffffffffffffffffffffffff1681565b42841015611b7b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e697377617056323a20455850495245440000000000000000000000000000604482015290519081900360640190fd5b60035473ffffffffffffffffffffffffffffffffffffffff80891660008181526004602090815260408083208054600180820190925582517f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c98186015280840196909652958d166060860152608085018c905260a085019590955260c08085018b90528151808603909101815260e0850182528051908301207f19010000000000000000000000000000000000000000000000000000000000006101008601526101028501969096526101228085019690965280518085039096018652610142840180825286519683019690962095839052610162840180825286905260ff89166101828501526101a284018890526101c28401879052519193926101e2808201937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081019281900390910190855afa158015611cdc573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff811615801590611d5757508873ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16145b611dc257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f556e697377617056323a20494e56414c49445f5349474e415455524500000000604482015290519081900360640190fd5b611dcd89898961259c565b505050505050505050565b600260209081526000928352604080842090915290825290205481565b600c54600114611e6657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55600654604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529051611fd49273ffffffffffffffffffffffffffffffffffffffff16916370a08231916024808301926020929190829003018186803b158015611edd57600080fd5b505afa158015611ef1573d6000803e3d6000fd5b505050506040513d6020811015611f0757600080fd5b5051600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b158015611f7a57600080fd5b505afa158015611f8e573d6000803e3d6000fd5b505050506040513d6020811015611fa457600080fd5b50516008546dffffffffffffffffffffffffffff808216916e0100000000000000000000000000009004166122e0565b6001600c55565b604080518082018252601981527f7472616e7366657228616464726573732c75696e743235362900000000000000602091820152815173ffffffffffffffffffffffffffffffffffffffff85811660248301526044808301869052845180840390910181526064909201845291810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001781529251815160009460609489169392918291908083835b602083106120e157805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016120a4565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114612143576040519150601f19603f3d011682016040523d82523d6000602084013e612148565b606091505b5091509150818015612176575080511580612176575080806020019051602081101561217357600080fd5b50515b6121e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f556e697377617056323a205452414e534645525f4641494c4544000000000000604482015290519081900360640190fd5b5050505050565b60008115806122035750508082028282828161220057fe5b04145b610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f64732d6d6174682d6d756c2d6f766572666c6f77000000000000000000000000604482015290519081900360640190fd5b80820382811115610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f64732d6d6174682d7375622d756e646572666c6f770000000000000000000000604482015290519081900360640190fd5b6dffffffffffffffffffffffffffff841180159061230c57506dffffffffffffffffffffffffffff8311155b61237757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f556e697377617056323a204f564552464c4f5700000000000000000000000000604482015290519081900360640190fd5b60085463ffffffff428116917c0100000000000000000000000000000000000000000000000000000000900481168203908116158015906123c757506dffffffffffffffffffffffffffff841615155b80156123e257506dffffffffffffffffffffffffffff831615155b15612492578063ffffffff16612425856123fb86612a57565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169063ffffffff612a7b16565b600980547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092169290920201905563ffffffff8116612465846123fb87612a57565b600a80547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216929092020190555b600880547fffffffffffffffffffffffffffffffffffff0000000000000000000000000000166dffffffffffffffffffffffffffff888116919091177fffffffff0000000000000000000000000000ffffffffffffffffffffffffffff166e0100000000000000000000000000008883168102919091177bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff871602179283905560408051848416815291909304909116602082015281517f1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1929181900390910190a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff808416600081815260026020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260016020526040902054612641908263ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600160205260408082209390935590841681522054612683908263ffffffff612abc16565b73ffffffffffffffffffffffffffffffffffffffff80841660008181526001602090815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b600080600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663017e7e586040518163ffffffff1660e01b815260040160206040518083038186803b15801561275757600080fd5b505afa15801561276b573d6000803e3d6000fd5b505050506040513d602081101561278157600080fd5b5051600b5473ffffffffffffffffffffffffffffffffffffffff821615801594509192509061286457801561285f5760006127d86112576dffffffffffffffffffffffffffff88811690881663ffffffff6121e816565b905060006127e583612878565b90508082111561285c576000612813612804848463ffffffff61226e16565b6000549063ffffffff6121e816565b905060006128388361282c86600563ffffffff6121e816565b9063ffffffff612abc16565b9050600081838161284557fe5b04905080156128585761285887826128ca565b5050505b50505b612870565b8015612870576000600b555b505092915050565b600060038211156128bb575080600160028204015b818110156128b5578091506002818285816128a457fe5b0401816128ad57fe5b04905061288d565b506128c5565b81156128c5575060015b919050565b6000546128dd908263ffffffff612abc16565b600090815573ffffffffffffffffffffffffffffffffffffffff8316815260016020526040902054612915908263ffffffff612abc16565b73ffffffffffffffffffffffffffffffffffffffff831660008181526001602090815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6000818310612989578161298b565b825b9392505050565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600160205260409020546129c8908263ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff831660009081526001602052604081209190915554612a02908263ffffffff61226e16565b600090815560408051838152905173ffffffffffffffffffffffffffffffffffffffff8516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef919081900360200190a35050565b6dffffffffffffffffffffffffffff166e0100000000000000000000000000000290565b60006dffffffffffffffffffffffffffff82167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff841681612ab457fe5b049392505050565b80820182811015610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f64732d6d6174682d6164642d6f766572666c6f77000000000000000000000000604482015290519081900360640190fdfe556e697377617056323a20494e53554646494349454e545f4f55545055545f414d4f554e54556e697377617056323a20494e53554646494349454e545f494e5055545f414d4f554e54556e697377617056323a20494e53554646494349454e545f4c4951554944495459556e697377617056323a20494e53554646494349454e545f4c49515549444954595f4255524e4544556e697377617056323a20494e53554646494349454e545f4c49515549444954595f4d494e544544a265627a7a723158207dca18479e58487606bf70c79e44d8dee62353c9ee6d01f9a9d70885b8765f2264736f6c63430005100032454950373132446f6d61696e28737472696e67206e616d652c737472696e672076657273696f6e2c75696e7432353620636861696e49642c6164647265737320766572696679696e67436f6e747261637429a265627a7a723158202760f92d7fa1db6f5aa16307bad65df4ebcc8550c4b1f03755ab8dfd830c178f64736f6c63430005100032",
}

// UniswapV2FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV2FactoryMetaData.ABI instead.
var UniswapV2FactoryABI = UniswapV2FactoryMetaData.ABI

// UniswapV2FactoryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const UniswapV2FactoryBinRuntime = `608060405234801561001057600080fd5b50600436106100885760003560e01c8063a2e74af61161005b578063a2e74af6146100fd578063c9c6539614610132578063e6a439051461016d578063f46901ed146101a857610088565b8063017e7e581461008d578063094b7415146100be5780631e3dd18b146100c6578063574f2ba3146100e3575b600080fd5b6100956101db565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b6100956101f7565b610095600480360360208110156100dc57600080fd5b5035610213565b6100eb610247565b60408051918252519081900360200190f35b6101306004803603602081101561011357600080fd5b503573ffffffffffffffffffffffffffffffffffffffff1661024d565b005b6100956004803603604081101561014857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602001351661031a565b6100956004803603604081101561018357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135811691602001351661076d565b610130600480360360208110156101be57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166107a0565b60005473ffffffffffffffffffffffffffffffffffffffff1681565b60015473ffffffffffffffffffffffffffffffffffffffff1681565b6003818154811061022057fe5b60009182526020909120015473ffffffffffffffffffffffffffffffffffffffff16905081565b60035490565b60015473ffffffffffffffffffffffffffffffffffffffff1633146102d357604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f556e697377617056323a20464f5242494444454e000000000000000000000000604482015290519081900360640190fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60008173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156103b757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f556e697377617056323a204944454e544943414c5f4144445245535345530000604482015290519081900360640190fd5b6000808373ffffffffffffffffffffffffffffffffffffffff168573ffffffffffffffffffffffffffffffffffffffff16106103f45783856103f7565b84845b909250905073ffffffffffffffffffffffffffffffffffffffff821661047e57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601760248201527f556e697377617056323a205a45524f5f41444452455353000000000000000000604482015290519081900360640190fd5b73ffffffffffffffffffffffffffffffffffffffff82811660009081526002602090815260408083208585168452909152902054161561051f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601660248201527f556e697377617056323a20504149525f45584953545300000000000000000000604482015290519081900360640190fd5b6060604051806020016105319061086d565b6020820181038252601f19601f82011660405250905060008383604051602001808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b81526014018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b815260140192505050604051602081830303815290604052805190602001209050808251602084016000f5604080517f485cc95500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8781166004830152868116602483015291519297509087169163485cc9559160448082019260009290919082900301818387803b15801561065e57600080fd5b505af1158015610672573d6000803e3d6000fd5b5050505073ffffffffffffffffffffffffffffffffffffffff84811660008181526002602081815260408084208987168086529083528185208054978d167fffffffffffffffffffffffff000000000000000000000000000000000000000098891681179091559383528185208686528352818520805488168517905560038054600181018255958190527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b90950180549097168417909655925483519283529082015281517f0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9929181900390910190a35050505092915050565b600260209081526000928352604080842090915290825290205473ffffffffffffffffffffffffffffffffffffffff1681565b60015473ffffffffffffffffffffffffffffffffffffffff16331461082657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f556e697377617056323a20464f5242494444454e000000000000000000000000604482015290519081900360640190fd5b600080547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b612d748061087b8339019056fe60806040526001600c5534801561001557600080fd5b506040514690806052612d228239604080519182900360520182208282018252600a8352692ab734b9bbb0b8102b1960b11b6020938401528151808301835260018152603160f81b908401528151808401919091527fbfcc8ef98ffbf7b6c3fec7bf5185b566b9863e35a9d83acd49ad6824b5969738818301527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6606082015260808101949094523060a0808601919091528151808603909101815260c09094019052825192019190912060035550600580546001600160a01b03191633179055612c1d806101056000396000f3fe608060405234801561001057600080fd5b50600436106101b95760003560e01c80636a627842116100f9578063ba9a7a5611610097578063d21220a711610071578063d21220a7146105da578063d505accf146105e2578063dd62ed3e14610640578063fff6cae91461067b576101b9565b8063ba9a7a5614610597578063bc25cf771461059f578063c45a0155146105d2576101b9565b80637ecebe00116100d35780637ecebe00146104d757806389afcb441461050a57806395d89b4114610556578063a9059cbb1461055e576101b9565b80636a6278421461046957806370a082311461049c5780637464fc3d146104cf576101b9565b806323b872dd116101665780633644e515116101405780633644e51514610416578063485cc9551461041e5780635909c0d5146104595780635a3d549314610461576101b9565b806323b872dd146103ad57806330adf81f146103f0578063313ce567146103f8576101b9565b8063095ea7b311610197578063095ea7b3146103155780630dfe16811461036257806318160ddd14610393576101b9565b8063022c0d9f146101be57806306fdde03146102595780630902f1ac146102d6575b600080fd5b610257600480360360808110156101d457600080fd5b81359160208101359173ffffffffffffffffffffffffffffffffffffffff604083013516919081019060808101606082013564010000000081111561021857600080fd5b82018360208201111561022a57600080fd5b8035906020019184600183028401116401000000008311171561024c57600080fd5b509092509050610683565b005b610261610d57565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561029b578181015183820152602001610283565b50505050905090810190601f1680156102c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102de610d90565b604080516dffffffffffffffffffffffffffff948516815292909316602083015263ffffffff168183015290519081900360600190f35b61034e6004803603604081101561032b57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610de5565b604080519115158252519081900360200190f35b61036a610dfc565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61039b610e18565b60408051918252519081900360200190f35b61034e600480360360608110156103c357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060400135610e1e565b61039b610efd565b610400610f21565b6040805160ff9092168252519081900360200190f35b61039b610f26565b6102576004803603604081101561043457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516610f2c565b61039b611005565b61039b61100b565b61039b6004803603602081101561047f57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611011565b61039b600480360360208110156104b257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113cb565b61039b6113dd565b61039b600480360360208110156104ed57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113e3565b61053d6004803603602081101561052057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113f5565b6040805192835260208301919091528051918290030190f35b610261611892565b61034e6004803603604081101561057457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356118cb565b61039b6118d8565b610257600480360360208110156105b557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166118de565b61036a611ad4565b61036a611af0565b610257600480360360e08110156105f857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060408101359060608101359060ff6080820135169060a08101359060c00135611b0c565b61039b6004803603604081101561065657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516611dd8565b610257611df5565b600c546001146106f457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55841515806107075750600084115b61075c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180612b2f6025913960400191505060405180910390fd5b600080610767610d90565b5091509150816dffffffffffffffffffffffffffff168710801561079a5750806dffffffffffffffffffffffffffff1686105b6107ef576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180612b786021913960400191505060405180910390fd5b600654600754600091829173ffffffffffffffffffffffffffffffffffffffff91821691908116908916821480159061085457508073ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff1614155b6108bf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f556e697377617056323a20494e56414c49445f544f0000000000000000000000604482015290519081900360640190fd5b8a156108d0576108d0828a8d611fdb565b89156108e1576108e1818a8c611fdb565b86156109c3578873ffffffffffffffffffffffffffffffffffffffff166310d1e85c338d8d8c8c6040518663ffffffff1660e01b8152600401808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001858152602001848152602001806020018281038252848482818152602001925080828437600081840152601f19601f8201169050808301925050509650505050505050600060405180830381600087803b1580156109aa57600080fd5b505af11580156109be573d6000803e3d6000fd5b505050505b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff8416916370a08231916024808301926020929190829003018186803b158015610a2f57600080fd5b505afa158015610a43573d6000803e3d6000fd5b505050506040513d6020811015610a5957600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191955073ffffffffffffffffffffffffffffffffffffffff8316916370a0823191602480820192602092909190829003018186803b158015610acb57600080fd5b505afa158015610adf573d6000803e3d6000fd5b505050506040513d6020811015610af557600080fd5b5051925060009150506dffffffffffffffffffffffffffff85168a90038311610b1f576000610b35565b89856dffffffffffffffffffffffffffff160383035b9050600089856dffffffffffffffffffffffffffff16038311610b59576000610b6f565b89856dffffffffffffffffffffffffffff160383035b90506000821180610b805750600081115b610bd5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526024815260200180612b546024913960400191505060405180910390fd5b6000610c09610beb84600363ffffffff6121e816565b610bfd876103e863ffffffff6121e816565b9063ffffffff61226e16565b90506000610c21610beb84600363ffffffff6121e816565b9050610c59620f4240610c4d6dffffffffffffffffffffffffffff8b8116908b1663ffffffff6121e816565b9063ffffffff6121e816565b610c69838363ffffffff6121e816565b1015610cd657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f556e697377617056323a204b0000000000000000000000000000000000000000604482015290519081900360640190fd5b5050610ce4848488886122e0565b60408051838152602081018390528082018d9052606081018c9052905173ffffffffffffffffffffffffffffffffffffffff8b169133917fd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d8229181900360800190a350506001600c55505050505050505050565b6040518060400160405280600a81526020017f556e69737761702056320000000000000000000000000000000000000000000081525081565b6008546dffffffffffffffffffffffffffff808216926e0100000000000000000000000000008304909116917c0100000000000000000000000000000000000000000000000000000000900463ffffffff1690565b6000610df233848461259c565b5060015b92915050565b60065473ffffffffffffffffffffffffffffffffffffffff1681565b60005481565b73ffffffffffffffffffffffffffffffffffffffff831660009081526002602090815260408083203384529091528120547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff14610ee85773ffffffffffffffffffffffffffffffffffffffff84166000908152600260209081526040808320338452909152902054610eb6908363ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff851660009081526002602090815260408083203384529091529020555b610ef384848461260b565b5060019392505050565b7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c981565b601281565b60035481565b60055473ffffffffffffffffffffffffffffffffffffffff163314610fb257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f556e697377617056323a20464f5242494444454e000000000000000000000000604482015290519081900360640190fd5b6006805473ffffffffffffffffffffffffffffffffffffffff9384167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560078054929093169116179055565b60095481565b600a5481565b6000600c5460011461108457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c81905580611094610d90565b50600654604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905193955091935060009273ffffffffffffffffffffffffffffffffffffffff909116916370a08231916024808301926020929190829003018186803b15801561110e57600080fd5b505afa158015611122573d6000803e3d6000fd5b505050506040513d602081101561113857600080fd5b5051600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905192935060009273ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b1580156111b157600080fd5b505afa1580156111c5573d6000803e3d6000fd5b505050506040513d60208110156111db57600080fd5b505190506000611201836dffffffffffffffffffffffffffff871663ffffffff61226e16565b90506000611225836dffffffffffffffffffffffffffff871663ffffffff61226e16565b9050600061123387876126ec565b600054909150806112705761125c6103e8610bfd611257878763ffffffff6121e816565b612878565b985061126b60006103e86128ca565b6112cd565b6112ca6dffffffffffffffffffffffffffff8916611294868463ffffffff6121e816565b8161129b57fe5b046dffffffffffffffffffffffffffff89166112bd868563ffffffff6121e816565b816112c457fe5b0461297a565b98505b60008911611326576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180612bc16028913960400191505060405180910390fd5b6113308a8a6128ca565b61133c86868a8a6122e0565b811561137e5760085461137a906dffffffffffffffffffffffffffff808216916e01000000000000000000000000000090041663ffffffff6121e816565b600b555b6040805185815260208101859052815133927f4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f928290030190a250506001600c5550949695505050505050565b60016020526000908152604090205481565b600b5481565b60046020526000908152604090205481565b600080600c5460011461146957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c81905580611479610d90565b50600654600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905194965092945073ffffffffffffffffffffffffffffffffffffffff9182169391169160009184916370a08231916024808301926020929190829003018186803b1580156114fb57600080fd5b505afa15801561150f573d6000803e3d6000fd5b505050506040513d602081101561152557600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191925060009173ffffffffffffffffffffffffffffffffffffffff8516916370a08231916024808301926020929190829003018186803b15801561159957600080fd5b505afa1580156115ad573d6000803e3d6000fd5b505050506040513d60208110156115c357600080fd5b5051306000908152600160205260408120549192506115e288886126ec565b600054909150806115f9848763ffffffff6121e816565b8161160057fe5b049a5080611614848663ffffffff6121e816565b8161161b57fe5b04995060008b11801561162e575060008a115b611683576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180612b996028913960400191505060405180910390fd5b61168d3084612992565b611698878d8d611fdb565b6116a3868d8c611fdb565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff8916916370a08231916024808301926020929190829003018186803b15801561170f57600080fd5b505afa158015611723573d6000803e3d6000fd5b505050506040513d602081101561173957600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191965073ffffffffffffffffffffffffffffffffffffffff8816916370a0823191602480820192602092909190829003018186803b1580156117ab57600080fd5b505afa1580156117bf573d6000803e3d6000fd5b505050506040513d60208110156117d557600080fd5b505193506117e585858b8b6122e0565b811561182757600854611823906dffffffffffffffffffffffffffff808216916e01000000000000000000000000000090041663ffffffff6121e816565b600b555b604080518c8152602081018c9052815173ffffffffffffffffffffffffffffffffffffffff8f169233927fdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496929081900390910190a35050505050505050506001600c81905550915091565b6040518060400160405280600681526020017f554e492d5632000000000000000000000000000000000000000000000000000081525081565b6000610df233848461260b565b6103e881565b600c5460011461194f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55600654600754600854604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff9485169490931692611a2b9285928792611a26926dffffffffffffffffffffffffffff169185916370a0823191602480820192602092909190829003018186803b1580156119ee57600080fd5b505afa158015611a02573d6000803e3d6000fd5b505050506040513d6020811015611a1857600080fd5b50519063ffffffff61226e16565b611fdb565b600854604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529051611aca9284928792611a26926e01000000000000000000000000000090046dffffffffffffffffffffffffffff169173ffffffffffffffffffffffffffffffffffffffff8616916370a0823191602480820192602092909190829003018186803b1580156119ee57600080fd5b50506001600c5550565b60055473ffffffffffffffffffffffffffffffffffffffff1681565b60075473ffffffffffffffffffffffffffffffffffffffff1681565b42841015611b7b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e697377617056323a20455850495245440000000000000000000000000000604482015290519081900360640190fd5b60035473ffffffffffffffffffffffffffffffffffffffff80891660008181526004602090815260408083208054600180820190925582517f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c98186015280840196909652958d166060860152608085018c905260a085019590955260c08085018b90528151808603909101815260e0850182528051908301207f19010000000000000000000000000000000000000000000000000000000000006101008601526101028501969096526101228085019690965280518085039096018652610142840180825286519683019690962095839052610162840180825286905260ff89166101828501526101a284018890526101c28401879052519193926101e2808201937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081019281900390910190855afa158015611cdc573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff811615801590611d5757508873ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16145b611dc257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f556e697377617056323a20494e56414c49445f5349474e415455524500000000604482015290519081900360640190fd5b611dcd89898961259c565b505050505050505050565b600260209081526000928352604080842090915290825290205481565b600c54600114611e6657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55600654604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529051611fd49273ffffffffffffffffffffffffffffffffffffffff16916370a08231916024808301926020929190829003018186803b158015611edd57600080fd5b505afa158015611ef1573d6000803e3d6000fd5b505050506040513d6020811015611f0757600080fd5b5051600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b158015611f7a57600080fd5b505afa158015611f8e573d6000803e3d6000fd5b505050506040513d6020811015611fa457600080fd5b50516008546dffffffffffffffffffffffffffff808216916e0100000000000000000000000000009004166122e0565b6001600c55565b604080518082018252601981527f7472616e7366657228616464726573732c75696e743235362900000000000000602091820152815173ffffffffffffffffffffffffffffffffffffffff85811660248301526044808301869052845180840390910181526064909201845291810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001781529251815160009460609489169392918291908083835b602083106120e157805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016120a4565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114612143576040519150601f19603f3d011682016040523d82523d6000602084013e612148565b606091505b5091509150818015612176575080511580612176575080806020019051602081101561217357600080fd5b50515b6121e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f556e697377617056323a205452414e534645525f4641494c4544000000000000604482015290519081900360640190fd5b5050505050565b60008115806122035750508082028282828161220057fe5b04145b610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f64732d6d6174682d6d756c2d6f766572666c6f77000000000000000000000000604482015290519081900360640190fd5b80820382811115610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f64732d6d6174682d7375622d756e646572666c6f770000000000000000000000604482015290519081900360640190fd5b6dffffffffffffffffffffffffffff841180159061230c57506dffffffffffffffffffffffffffff8311155b61237757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f556e697377617056323a204f564552464c4f5700000000000000000000000000604482015290519081900360640190fd5b60085463ffffffff428116917c0100000000000000000000000000000000000000000000000000000000900481168203908116158015906123c757506dffffffffffffffffffffffffffff841615155b80156123e257506dffffffffffffffffffffffffffff831615155b15612492578063ffffffff16612425856123fb86612a57565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169063ffffffff612a7b16565b600980547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092169290920201905563ffffffff8116612465846123fb87612a57565b600a80547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216929092020190555b600880547fffffffffffffffffffffffffffffffffffff0000000000000000000000000000166dffffffffffffffffffffffffffff888116919091177fffffffff0000000000000000000000000000ffffffffffffffffffffffffffff166e0100000000000000000000000000008883168102919091177bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff871602179283905560408051848416815291909304909116602082015281517f1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1929181900390910190a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff808416600081815260026020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260016020526040902054612641908263ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600160205260408082209390935590841681522054612683908263ffffffff612abc16565b73ffffffffffffffffffffffffffffffffffffffff80841660008181526001602090815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b600080600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663017e7e586040518163ffffffff1660e01b815260040160206040518083038186803b15801561275757600080fd5b505afa15801561276b573d6000803e3d6000fd5b505050506040513d602081101561278157600080fd5b5051600b5473ffffffffffffffffffffffffffffffffffffffff821615801594509192509061286457801561285f5760006127d86112576dffffffffffffffffffffffffffff88811690881663ffffffff6121e816565b905060006127e583612878565b90508082111561285c576000612813612804848463ffffffff61226e16565b6000549063ffffffff6121e816565b905060006128388361282c86600563ffffffff6121e816565b9063ffffffff612abc16565b9050600081838161284557fe5b04905080156128585761285887826128ca565b5050505b50505b612870565b8015612870576000600b555b505092915050565b600060038211156128bb575080600160028204015b818110156128b5578091506002818285816128a457fe5b0401816128ad57fe5b04905061288d565b506128c5565b81156128c5575060015b919050565b6000546128dd908263ffffffff612abc16565b600090815573ffffffffffffffffffffffffffffffffffffffff8316815260016020526040902054612915908263ffffffff612abc16565b73ffffffffffffffffffffffffffffffffffffffff831660008181526001602090815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6000818310612989578161298b565b825b9392505050565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600160205260409020546129c8908263ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff831660009081526001602052604081209190915554612a02908263ffffffff61226e16565b600090815560408051838152905173ffffffffffffffffffffffffffffffffffffffff8516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef919081900360200190a35050565b6dffffffffffffffffffffffffffff166e0100000000000000000000000000000290565b60006dffffffffffffffffffffffffffff82167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff841681612ab457fe5b049392505050565b80820182811015610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f64732d6d6174682d6164642d6f766572666c6f77000000000000000000000000604482015290519081900360640190fdfe556e697377617056323a20494e53554646494349454e545f4f55545055545f414d4f554e54556e697377617056323a20494e53554646494349454e545f494e5055545f414d4f554e54556e697377617056323a20494e53554646494349454e545f4c4951554944495459556e697377617056323a20494e53554646494349454e545f4c49515549444954595f4255524e4544556e697377617056323a20494e53554646494349454e545f4c49515549444954595f4d494e544544a265627a7a723158207dca18479e58487606bf70c79e44d8dee62353c9ee6d01f9a9d70885b8765f2264736f6c63430005100032454950373132446f6d61696e28737472696e67206e616d652c737472696e672076657273696f6e2c75696e7432353620636861696e49642c6164647265737320766572696679696e67436f6e747261637429a265627a7a723158202760f92d7fa1db6f5aa16307bad65df4ebcc8550c4b1f03755ab8dfd830c178f64736f6c63430005100032`

// Deprecated: Use UniswapV2FactoryMetaData.Sigs instead.
// UniswapV2FactoryFuncSigs maps the 4-byte function signature to its string representation.
var UniswapV2FactoryFuncSigs = UniswapV2FactoryMetaData.Sigs

// UniswapV2FactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniswapV2FactoryMetaData.Bin instead.
var UniswapV2FactoryBin = UniswapV2FactoryMetaData.Bin

// DeployUniswapV2Factory deploys a new Kaia contract, binding an instance of UniswapV2Factory to it.
func DeployUniswapV2Factory(auth *bind.TransactOpts, backend bind.ContractBackend, _feeToSetter common.Address) (common.Address, *types.Transaction, *UniswapV2Factory, error) {
	parsed, err := UniswapV2FactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniswapV2FactoryBin), backend, _feeToSetter)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniswapV2Factory{UniswapV2FactoryCaller: UniswapV2FactoryCaller{contract: contract}, UniswapV2FactoryTransactor: UniswapV2FactoryTransactor{contract: contract}, UniswapV2FactoryFilterer: UniswapV2FactoryFilterer{contract: contract}}, nil
}

// UniswapV2Factory is an auto generated Go binding around a Kaia contract.
type UniswapV2Factory struct {
	UniswapV2FactoryCaller     // Read-only binding to the contract
	UniswapV2FactoryTransactor // Write-only binding to the contract
	UniswapV2FactoryFilterer   // Log filterer for contract events
}

// UniswapV2FactoryCaller is an auto generated read-only Go binding around a Kaia contract.
type UniswapV2FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2FactoryTransactor is an auto generated write-only Go binding around a Kaia contract.
type UniswapV2FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2FactoryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type UniswapV2FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2FactorySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type UniswapV2FactorySession struct {
	Contract     *UniswapV2Factory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapV2FactoryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type UniswapV2FactoryCallerSession struct {
	Contract *UniswapV2FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// UniswapV2FactoryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type UniswapV2FactoryTransactorSession struct {
	Contract     *UniswapV2FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// UniswapV2FactoryRaw is an auto generated low-level Go binding around a Kaia contract.
type UniswapV2FactoryRaw struct {
	Contract *UniswapV2Factory // Generic contract binding to access the raw methods on
}

// UniswapV2FactoryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type UniswapV2FactoryCallerRaw struct {
	Contract *UniswapV2FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV2FactoryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type UniswapV2FactoryTransactorRaw struct {
	Contract *UniswapV2FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV2Factory creates a new instance of UniswapV2Factory, bound to a specific deployed contract.
func NewUniswapV2Factory(address common.Address, backend bind.ContractBackend) (*UniswapV2Factory, error) {
	contract, err := bindUniswapV2Factory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV2Factory{UniswapV2FactoryCaller: UniswapV2FactoryCaller{contract: contract}, UniswapV2FactoryTransactor: UniswapV2FactoryTransactor{contract: contract}, UniswapV2FactoryFilterer: UniswapV2FactoryFilterer{contract: contract}}, nil
}

// NewUniswapV2FactoryCaller creates a new read-only instance of UniswapV2Factory, bound to a specific deployed contract.
func NewUniswapV2FactoryCaller(address common.Address, caller bind.ContractCaller) (*UniswapV2FactoryCaller, error) {
	contract, err := bindUniswapV2Factory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2FactoryCaller{contract: contract}, nil
}

// NewUniswapV2FactoryTransactor creates a new write-only instance of UniswapV2Factory, bound to a specific deployed contract.
func NewUniswapV2FactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV2FactoryTransactor, error) {
	contract, err := bindUniswapV2Factory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2FactoryTransactor{contract: contract}, nil
}

// NewUniswapV2FactoryFilterer creates a new log filterer instance of UniswapV2Factory, bound to a specific deployed contract.
func NewUniswapV2FactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV2FactoryFilterer, error) {
	contract, err := bindUniswapV2Factory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV2FactoryFilterer{contract: contract}, nil
}

// bindUniswapV2Factory binds a generic wrapper to an already deployed contract.
func bindUniswapV2Factory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV2FactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Factory *UniswapV2FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Factory.Contract.UniswapV2FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Factory *UniswapV2FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.UniswapV2FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Factory *UniswapV2FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.UniswapV2FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Factory *UniswapV2FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Factory *UniswapV2FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Factory *UniswapV2FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.contract.Transact(opts, method, params...)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCaller) AllPairs(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Factory.contract.Call(opts, &out, "allPairs", arg0)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactorySession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _UniswapV2Factory.Contract.AllPairs(&_UniswapV2Factory.CallOpts, arg0)
}

// AllPairs is a free data retrieval call binding the contract method 0x1e3dd18b.
//
// Solidity: function allPairs(uint256 ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCallerSession) AllPairs(arg0 *big.Int) (common.Address, error) {
	return _UniswapV2Factory.Contract.AllPairs(&_UniswapV2Factory.CallOpts, arg0)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_UniswapV2Factory *UniswapV2FactoryCaller) AllPairsLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Factory.contract.Call(opts, &out, "allPairsLength")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_UniswapV2Factory *UniswapV2FactorySession) AllPairsLength() (*big.Int, error) {
	return _UniswapV2Factory.Contract.AllPairsLength(&_UniswapV2Factory.CallOpts)
}

// AllPairsLength is a free data retrieval call binding the contract method 0x574f2ba3.
//
// Solidity: function allPairsLength() view returns(uint256)
func (_UniswapV2Factory *UniswapV2FactoryCallerSession) AllPairsLength() (*big.Int, error) {
	return _UniswapV2Factory.Contract.AllPairsLength(&_UniswapV2Factory.CallOpts)
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCaller) FeeTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Factory.contract.Call(opts, &out, "feeTo")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_UniswapV2Factory *UniswapV2FactorySession) FeeTo() (common.Address, error) {
	return _UniswapV2Factory.Contract.FeeTo(&_UniswapV2Factory.CallOpts)
}

// FeeTo is a free data retrieval call binding the contract method 0x017e7e58.
//
// Solidity: function feeTo() view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCallerSession) FeeTo() (common.Address, error) {
	return _UniswapV2Factory.Contract.FeeTo(&_UniswapV2Factory.CallOpts)
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCaller) FeeToSetter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Factory.contract.Call(opts, &out, "feeToSetter")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_UniswapV2Factory *UniswapV2FactorySession) FeeToSetter() (common.Address, error) {
	return _UniswapV2Factory.Contract.FeeToSetter(&_UniswapV2Factory.CallOpts)
}

// FeeToSetter is a free data retrieval call binding the contract method 0x094b7415.
//
// Solidity: function feeToSetter() view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCallerSession) FeeToSetter() (common.Address, error) {
	return _UniswapV2Factory.Contract.FeeToSetter(&_UniswapV2Factory.CallOpts)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCaller) GetPair(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Factory.contract.Call(opts, &out, "getPair", arg0, arg1)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactorySession) GetPair(arg0 common.Address, arg1 common.Address) (common.Address, error) {
	return _UniswapV2Factory.Contract.GetPair(&_UniswapV2Factory.CallOpts, arg0, arg1)
}

// GetPair is a free data retrieval call binding the contract method 0xe6a43905.
//
// Solidity: function getPair(address , address ) view returns(address)
func (_UniswapV2Factory *UniswapV2FactoryCallerSession) GetPair(arg0 common.Address, arg1 common.Address) (common.Address, error) {
	return _UniswapV2Factory.Contract.GetPair(&_UniswapV2Factory.CallOpts, arg0, arg1)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_UniswapV2Factory *UniswapV2FactoryTransactor) CreatePair(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.contract.Transact(opts, "createPair", tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_UniswapV2Factory *UniswapV2FactorySession) CreatePair(tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.CreatePair(&_UniswapV2Factory.TransactOpts, tokenA, tokenB)
}

// CreatePair is a paid mutator transaction binding the contract method 0xc9c65396.
//
// Solidity: function createPair(address tokenA, address tokenB) returns(address pair)
func (_UniswapV2Factory *UniswapV2FactoryTransactorSession) CreatePair(tokenA common.Address, tokenB common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.CreatePair(&_UniswapV2Factory.TransactOpts, tokenA, tokenB)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address _feeTo) returns()
func (_UniswapV2Factory *UniswapV2FactoryTransactor) SetFeeTo(opts *bind.TransactOpts, _feeTo common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.contract.Transact(opts, "setFeeTo", _feeTo)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address _feeTo) returns()
func (_UniswapV2Factory *UniswapV2FactorySession) SetFeeTo(_feeTo common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.SetFeeTo(&_UniswapV2Factory.TransactOpts, _feeTo)
}

// SetFeeTo is a paid mutator transaction binding the contract method 0xf46901ed.
//
// Solidity: function setFeeTo(address _feeTo) returns()
func (_UniswapV2Factory *UniswapV2FactoryTransactorSession) SetFeeTo(_feeTo common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.SetFeeTo(&_UniswapV2Factory.TransactOpts, _feeTo)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address _feeToSetter) returns()
func (_UniswapV2Factory *UniswapV2FactoryTransactor) SetFeeToSetter(opts *bind.TransactOpts, _feeToSetter common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.contract.Transact(opts, "setFeeToSetter", _feeToSetter)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address _feeToSetter) returns()
func (_UniswapV2Factory *UniswapV2FactorySession) SetFeeToSetter(_feeToSetter common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.SetFeeToSetter(&_UniswapV2Factory.TransactOpts, _feeToSetter)
}

// SetFeeToSetter is a paid mutator transaction binding the contract method 0xa2e74af6.
//
// Solidity: function setFeeToSetter(address _feeToSetter) returns()
func (_UniswapV2Factory *UniswapV2FactoryTransactorSession) SetFeeToSetter(_feeToSetter common.Address) (*types.Transaction, error) {
	return _UniswapV2Factory.Contract.SetFeeToSetter(&_UniswapV2Factory.TransactOpts, _feeToSetter)
}

// UniswapV2FactoryPairCreatedIterator is returned from FilterPairCreated and is used to iterate over the raw logs and unpacked data for PairCreated events raised by the UniswapV2Factory contract.
type UniswapV2FactoryPairCreatedIterator struct {
	Event *UniswapV2FactoryPairCreated // Event containing the contract specifics and raw log

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
func (it *UniswapV2FactoryPairCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2FactoryPairCreated)
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
		it.Event = new(UniswapV2FactoryPairCreated)
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
func (it *UniswapV2FactoryPairCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2FactoryPairCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2FactoryPairCreated represents a PairCreated event raised by the UniswapV2Factory contract.
type UniswapV2FactoryPairCreated struct {
	Token0 common.Address
	Token1 common.Address
	Pair   common.Address
	Arg3   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPairCreated is a free log retrieval operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_UniswapV2Factory *UniswapV2FactoryFilterer) FilterPairCreated(opts *bind.FilterOpts, token0 []common.Address, token1 []common.Address) (*UniswapV2FactoryPairCreatedIterator, error) {
	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _UniswapV2Factory.contract.FilterLogs(opts, "PairCreated", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return &UniswapV2FactoryPairCreatedIterator{contract: _UniswapV2Factory.contract, event: "PairCreated", logs: logs, sub: sub}, nil
}

// WatchPairCreated is a free log subscription operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_UniswapV2Factory *UniswapV2FactoryFilterer) WatchPairCreated(opts *bind.WatchOpts, sink chan<- *UniswapV2FactoryPairCreated, token0 []common.Address, token1 []common.Address) (event.Subscription, error) {
	var token0Rule []interface{}
	for _, token0Item := range token0 {
		token0Rule = append(token0Rule, token0Item)
	}
	var token1Rule []interface{}
	for _, token1Item := range token1 {
		token1Rule = append(token1Rule, token1Item)
	}

	logs, sub, err := _UniswapV2Factory.contract.WatchLogs(opts, "PairCreated", token0Rule, token1Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2FactoryPairCreated)
				if err := _UniswapV2Factory.contract.UnpackLog(event, "PairCreated", log); err != nil {
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

// ParsePairCreated is a log parse operation binding the contract event 0x0d3648bd0f6ba80134a33ba9275ac585d9d315f0ad8355cddefde31afa28d0e9.
//
// Solidity: event PairCreated(address indexed token0, address indexed token1, address pair, uint256 arg3)
func (_UniswapV2Factory *UniswapV2FactoryFilterer) ParsePairCreated(log types.Log) (*UniswapV2FactoryPairCreated, error) {
	event := new(UniswapV2FactoryPairCreated)
	if err := _UniswapV2Factory.contract.UnpackLog(event, "PairCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV2PairMetaData contains all meta data concerning the UniswapV2Pair contract.
var UniswapV2PairMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1In\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint112\",\"name\":\"reserve0\",\"type\":\"uint112\"},{\"indexed\":false,\"internalType\":\"uint112\",\"name\":\"reserve1\",\"type\":\"uint112\"}],\"name\":\"Sync\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"DOMAIN_SEPARATOR\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MINIMUM_LIQUIDITY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"PERMIT_TYPEHASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getReserves\",\"outputs\":[{\"internalType\":\"uint112\",\"name\":\"_reserve0\",\"type\":\"uint112\"},{\"internalType\":\"uint112\",\"name\":\"_reserve1\",\"type\":\"uint112\"},{\"internalType\":\"uint32\",\"name\":\"_blockTimestampLast\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token1\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"kLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"permit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"price0CumulativeLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"price1CumulativeLast\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"skim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount0Out\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1Out\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"sync\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token0\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token1\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"3644e515": "DOMAIN_SEPARATOR()",
		"ba9a7a56": "MINIMUM_LIQUIDITY()",
		"30adf81f": "PERMIT_TYPEHASH()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"89afcb44": "burn(address)",
		"313ce567": "decimals()",
		"c45a0155": "factory()",
		"0902f1ac": "getReserves()",
		"485cc955": "initialize(address,address)",
		"7464fc3d": "kLast()",
		"6a627842": "mint(address)",
		"06fdde03": "name()",
		"7ecebe00": "nonces(address)",
		"d505accf": "permit(address,address,uint256,uint256,uint8,bytes32,bytes32)",
		"5909c0d5": "price0CumulativeLast()",
		"5a3d5493": "price1CumulativeLast()",
		"bc25cf77": "skim(address)",
		"022c0d9f": "swap(uint256,uint256,address,bytes)",
		"95d89b41": "symbol()",
		"fff6cae9": "sync()",
		"0dfe1681": "token0()",
		"d21220a7": "token1()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
	},
	Bin: "0x60806040526001600c5534801561001557600080fd5b506040514690806052612d228239604080519182900360520182208282018252600a8352692ab734b9bbb0b8102b1960b11b6020938401528151808301835260018152603160f81b908401528151808401919091527fbfcc8ef98ffbf7b6c3fec7bf5185b566b9863e35a9d83acd49ad6824b5969738818301527fc89efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6606082015260808101949094523060a0808601919091528151808603909101815260c09094019052825192019190912060035550600580546001600160a01b03191633179055612c1d806101056000396000f3fe608060405234801561001057600080fd5b50600436106101b95760003560e01c80636a627842116100f9578063ba9a7a5611610097578063d21220a711610071578063d21220a7146105da578063d505accf146105e2578063dd62ed3e14610640578063fff6cae91461067b576101b9565b8063ba9a7a5614610597578063bc25cf771461059f578063c45a0155146105d2576101b9565b80637ecebe00116100d35780637ecebe00146104d757806389afcb441461050a57806395d89b4114610556578063a9059cbb1461055e576101b9565b80636a6278421461046957806370a082311461049c5780637464fc3d146104cf576101b9565b806323b872dd116101665780633644e515116101405780633644e51514610416578063485cc9551461041e5780635909c0d5146104595780635a3d549314610461576101b9565b806323b872dd146103ad57806330adf81f146103f0578063313ce567146103f8576101b9565b8063095ea7b311610197578063095ea7b3146103155780630dfe16811461036257806318160ddd14610393576101b9565b8063022c0d9f146101be57806306fdde03146102595780630902f1ac146102d6575b600080fd5b610257600480360360808110156101d457600080fd5b81359160208101359173ffffffffffffffffffffffffffffffffffffffff604083013516919081019060808101606082013564010000000081111561021857600080fd5b82018360208201111561022a57600080fd5b8035906020019184600183028401116401000000008311171561024c57600080fd5b509092509050610683565b005b610261610d57565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561029b578181015183820152602001610283565b50505050905090810190601f1680156102c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102de610d90565b604080516dffffffffffffffffffffffffffff948516815292909316602083015263ffffffff168183015290519081900360600190f35b61034e6004803603604081101561032b57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610de5565b604080519115158252519081900360200190f35b61036a610dfc565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61039b610e18565b60408051918252519081900360200190f35b61034e600480360360608110156103c357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060400135610e1e565b61039b610efd565b610400610f21565b6040805160ff9092168252519081900360200190f35b61039b610f26565b6102576004803603604081101561043457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516610f2c565b61039b611005565b61039b61100b565b61039b6004803603602081101561047f57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611011565b61039b600480360360208110156104b257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113cb565b61039b6113dd565b61039b600480360360208110156104ed57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113e3565b61053d6004803603602081101561052057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113f5565b6040805192835260208301919091528051918290030190f35b610261611892565b61034e6004803603604081101561057457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356118cb565b61039b6118d8565b610257600480360360208110156105b557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166118de565b61036a611ad4565b61036a611af0565b610257600480360360e08110156105f857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060408101359060608101359060ff6080820135169060a08101359060c00135611b0c565b61039b6004803603604081101561065657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516611dd8565b610257611df5565b600c546001146106f457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55841515806107075750600084115b61075c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180612b2f6025913960400191505060405180910390fd5b600080610767610d90565b5091509150816dffffffffffffffffffffffffffff168710801561079a5750806dffffffffffffffffffffffffffff1686105b6107ef576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180612b786021913960400191505060405180910390fd5b600654600754600091829173ffffffffffffffffffffffffffffffffffffffff91821691908116908916821480159061085457508073ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff1614155b6108bf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f556e697377617056323a20494e56414c49445f544f0000000000000000000000604482015290519081900360640190fd5b8a156108d0576108d0828a8d611fdb565b89156108e1576108e1818a8c611fdb565b86156109c3578873ffffffffffffffffffffffffffffffffffffffff166310d1e85c338d8d8c8c6040518663ffffffff1660e01b8152600401808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001858152602001848152602001806020018281038252848482818152602001925080828437600081840152601f19601f8201169050808301925050509650505050505050600060405180830381600087803b1580156109aa57600080fd5b505af11580156109be573d6000803e3d6000fd5b505050505b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff8416916370a08231916024808301926020929190829003018186803b158015610a2f57600080fd5b505afa158015610a43573d6000803e3d6000fd5b505050506040513d6020811015610a5957600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191955073ffffffffffffffffffffffffffffffffffffffff8316916370a0823191602480820192602092909190829003018186803b158015610acb57600080fd5b505afa158015610adf573d6000803e3d6000fd5b505050506040513d6020811015610af557600080fd5b5051925060009150506dffffffffffffffffffffffffffff85168a90038311610b1f576000610b35565b89856dffffffffffffffffffffffffffff160383035b9050600089856dffffffffffffffffffffffffffff16038311610b59576000610b6f565b89856dffffffffffffffffffffffffffff160383035b90506000821180610b805750600081115b610bd5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526024815260200180612b546024913960400191505060405180910390fd5b6000610c09610beb84600363ffffffff6121e816565b610bfd876103e863ffffffff6121e816565b9063ffffffff61226e16565b90506000610c21610beb84600363ffffffff6121e816565b9050610c59620f4240610c4d6dffffffffffffffffffffffffffff8b8116908b1663ffffffff6121e816565b9063ffffffff6121e816565b610c69838363ffffffff6121e816565b1015610cd657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f556e697377617056323a204b0000000000000000000000000000000000000000604482015290519081900360640190fd5b5050610ce4848488886122e0565b60408051838152602081018390528082018d9052606081018c9052905173ffffffffffffffffffffffffffffffffffffffff8b169133917fd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d8229181900360800190a350506001600c55505050505050505050565b6040518060400160405280600a81526020017f556e69737761702056320000000000000000000000000000000000000000000081525081565b6008546dffffffffffffffffffffffffffff808216926e0100000000000000000000000000008304909116917c0100000000000000000000000000000000000000000000000000000000900463ffffffff1690565b6000610df233848461259c565b5060015b92915050565b60065473ffffffffffffffffffffffffffffffffffffffff1681565b60005481565b73ffffffffffffffffffffffffffffffffffffffff831660009081526002602090815260408083203384529091528120547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff14610ee85773ffffffffffffffffffffffffffffffffffffffff84166000908152600260209081526040808320338452909152902054610eb6908363ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff851660009081526002602090815260408083203384529091529020555b610ef384848461260b565b5060019392505050565b7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c981565b601281565b60035481565b60055473ffffffffffffffffffffffffffffffffffffffff163314610fb257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f556e697377617056323a20464f5242494444454e000000000000000000000000604482015290519081900360640190fd5b6006805473ffffffffffffffffffffffffffffffffffffffff9384167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560078054929093169116179055565b60095481565b600a5481565b6000600c5460011461108457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c81905580611094610d90565b50600654604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905193955091935060009273ffffffffffffffffffffffffffffffffffffffff909116916370a08231916024808301926020929190829003018186803b15801561110e57600080fd5b505afa158015611122573d6000803e3d6000fd5b505050506040513d602081101561113857600080fd5b5051600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905192935060009273ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b1580156111b157600080fd5b505afa1580156111c5573d6000803e3d6000fd5b505050506040513d60208110156111db57600080fd5b505190506000611201836dffffffffffffffffffffffffffff871663ffffffff61226e16565b90506000611225836dffffffffffffffffffffffffffff871663ffffffff61226e16565b9050600061123387876126ec565b600054909150806112705761125c6103e8610bfd611257878763ffffffff6121e816565b612878565b985061126b60006103e86128ca565b6112cd565b6112ca6dffffffffffffffffffffffffffff8916611294868463ffffffff6121e816565b8161129b57fe5b046dffffffffffffffffffffffffffff89166112bd868563ffffffff6121e816565b816112c457fe5b0461297a565b98505b60008911611326576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180612bc16028913960400191505060405180910390fd5b6113308a8a6128ca565b61133c86868a8a6122e0565b811561137e5760085461137a906dffffffffffffffffffffffffffff808216916e01000000000000000000000000000090041663ffffffff6121e816565b600b555b6040805185815260208101859052815133927f4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f928290030190a250506001600c5550949695505050505050565b60016020526000908152604090205481565b600b5481565b60046020526000908152604090205481565b600080600c5460011461146957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c81905580611479610d90565b50600654600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905194965092945073ffffffffffffffffffffffffffffffffffffffff9182169391169160009184916370a08231916024808301926020929190829003018186803b1580156114fb57600080fd5b505afa15801561150f573d6000803e3d6000fd5b505050506040513d602081101561152557600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191925060009173ffffffffffffffffffffffffffffffffffffffff8516916370a08231916024808301926020929190829003018186803b15801561159957600080fd5b505afa1580156115ad573d6000803e3d6000fd5b505050506040513d60208110156115c357600080fd5b5051306000908152600160205260408120549192506115e288886126ec565b600054909150806115f9848763ffffffff6121e816565b8161160057fe5b049a5080611614848663ffffffff6121e816565b8161161b57fe5b04995060008b11801561162e575060008a115b611683576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180612b996028913960400191505060405180910390fd5b61168d3084612992565b611698878d8d611fdb565b6116a3868d8c611fdb565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff8916916370a08231916024808301926020929190829003018186803b15801561170f57600080fd5b505afa158015611723573d6000803e3d6000fd5b505050506040513d602081101561173957600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191965073ffffffffffffffffffffffffffffffffffffffff8816916370a0823191602480820192602092909190829003018186803b1580156117ab57600080fd5b505afa1580156117bf573d6000803e3d6000fd5b505050506040513d60208110156117d557600080fd5b505193506117e585858b8b6122e0565b811561182757600854611823906dffffffffffffffffffffffffffff808216916e01000000000000000000000000000090041663ffffffff6121e816565b600b555b604080518c8152602081018c9052815173ffffffffffffffffffffffffffffffffffffffff8f169233927fdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496929081900390910190a35050505050505050506001600c81905550915091565b6040518060400160405280600681526020017f554e492d5632000000000000000000000000000000000000000000000000000081525081565b6000610df233848461260b565b6103e881565b600c5460011461194f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55600654600754600854604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff9485169490931692611a2b9285928792611a26926dffffffffffffffffffffffffffff169185916370a0823191602480820192602092909190829003018186803b1580156119ee57600080fd5b505afa158015611a02573d6000803e3d6000fd5b505050506040513d6020811015611a1857600080fd5b50519063ffffffff61226e16565b611fdb565b600854604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529051611aca9284928792611a26926e01000000000000000000000000000090046dffffffffffffffffffffffffffff169173ffffffffffffffffffffffffffffffffffffffff8616916370a0823191602480820192602092909190829003018186803b1580156119ee57600080fd5b50506001600c5550565b60055473ffffffffffffffffffffffffffffffffffffffff1681565b60075473ffffffffffffffffffffffffffffffffffffffff1681565b42841015611b7b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e697377617056323a20455850495245440000000000000000000000000000604482015290519081900360640190fd5b60035473ffffffffffffffffffffffffffffffffffffffff80891660008181526004602090815260408083208054600180820190925582517f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c98186015280840196909652958d166060860152608085018c905260a085019590955260c08085018b90528151808603909101815260e0850182528051908301207f19010000000000000000000000000000000000000000000000000000000000006101008601526101028501969096526101228085019690965280518085039096018652610142840180825286519683019690962095839052610162840180825286905260ff89166101828501526101a284018890526101c28401879052519193926101e2808201937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081019281900390910190855afa158015611cdc573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff811615801590611d5757508873ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16145b611dc257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f556e697377617056323a20494e56414c49445f5349474e415455524500000000604482015290519081900360640190fd5b611dcd89898961259c565b505050505050505050565b600260209081526000928352604080842090915290825290205481565b600c54600114611e6657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55600654604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529051611fd49273ffffffffffffffffffffffffffffffffffffffff16916370a08231916024808301926020929190829003018186803b158015611edd57600080fd5b505afa158015611ef1573d6000803e3d6000fd5b505050506040513d6020811015611f0757600080fd5b5051600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b158015611f7a57600080fd5b505afa158015611f8e573d6000803e3d6000fd5b505050506040513d6020811015611fa457600080fd5b50516008546dffffffffffffffffffffffffffff808216916e0100000000000000000000000000009004166122e0565b6001600c55565b604080518082018252601981527f7472616e7366657228616464726573732c75696e743235362900000000000000602091820152815173ffffffffffffffffffffffffffffffffffffffff85811660248301526044808301869052845180840390910181526064909201845291810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001781529251815160009460609489169392918291908083835b602083106120e157805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016120a4565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114612143576040519150601f19603f3d011682016040523d82523d6000602084013e612148565b606091505b5091509150818015612176575080511580612176575080806020019051602081101561217357600080fd5b50515b6121e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f556e697377617056323a205452414e534645525f4641494c4544000000000000604482015290519081900360640190fd5b5050505050565b60008115806122035750508082028282828161220057fe5b04145b610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f64732d6d6174682d6d756c2d6f766572666c6f77000000000000000000000000604482015290519081900360640190fd5b80820382811115610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f64732d6d6174682d7375622d756e646572666c6f770000000000000000000000604482015290519081900360640190fd5b6dffffffffffffffffffffffffffff841180159061230c57506dffffffffffffffffffffffffffff8311155b61237757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f556e697377617056323a204f564552464c4f5700000000000000000000000000604482015290519081900360640190fd5b60085463ffffffff428116917c0100000000000000000000000000000000000000000000000000000000900481168203908116158015906123c757506dffffffffffffffffffffffffffff841615155b80156123e257506dffffffffffffffffffffffffffff831615155b15612492578063ffffffff16612425856123fb86612a57565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169063ffffffff612a7b16565b600980547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092169290920201905563ffffffff8116612465846123fb87612a57565b600a80547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216929092020190555b600880547fffffffffffffffffffffffffffffffffffff0000000000000000000000000000166dffffffffffffffffffffffffffff888116919091177fffffffff0000000000000000000000000000ffffffffffffffffffffffffffff166e0100000000000000000000000000008883168102919091177bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff871602179283905560408051848416815291909304909116602082015281517f1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1929181900390910190a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff808416600081815260026020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260016020526040902054612641908263ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600160205260408082209390935590841681522054612683908263ffffffff612abc16565b73ffffffffffffffffffffffffffffffffffffffff80841660008181526001602090815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b600080600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663017e7e586040518163ffffffff1660e01b815260040160206040518083038186803b15801561275757600080fd5b505afa15801561276b573d6000803e3d6000fd5b505050506040513d602081101561278157600080fd5b5051600b5473ffffffffffffffffffffffffffffffffffffffff821615801594509192509061286457801561285f5760006127d86112576dffffffffffffffffffffffffffff88811690881663ffffffff6121e816565b905060006127e583612878565b90508082111561285c576000612813612804848463ffffffff61226e16565b6000549063ffffffff6121e816565b905060006128388361282c86600563ffffffff6121e816565b9063ffffffff612abc16565b9050600081838161284557fe5b04905080156128585761285887826128ca565b5050505b50505b612870565b8015612870576000600b555b505092915050565b600060038211156128bb575080600160028204015b818110156128b5578091506002818285816128a457fe5b0401816128ad57fe5b04905061288d565b506128c5565b81156128c5575060015b919050565b6000546128dd908263ffffffff612abc16565b600090815573ffffffffffffffffffffffffffffffffffffffff8316815260016020526040902054612915908263ffffffff612abc16565b73ffffffffffffffffffffffffffffffffffffffff831660008181526001602090815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6000818310612989578161298b565b825b9392505050565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600160205260409020546129c8908263ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff831660009081526001602052604081209190915554612a02908263ffffffff61226e16565b600090815560408051838152905173ffffffffffffffffffffffffffffffffffffffff8516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef919081900360200190a35050565b6dffffffffffffffffffffffffffff166e0100000000000000000000000000000290565b60006dffffffffffffffffffffffffffff82167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff841681612ab457fe5b049392505050565b80820182811015610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f64732d6d6174682d6164642d6f766572666c6f77000000000000000000000000604482015290519081900360640190fdfe556e697377617056323a20494e53554646494349454e545f4f55545055545f414d4f554e54556e697377617056323a20494e53554646494349454e545f494e5055545f414d4f554e54556e697377617056323a20494e53554646494349454e545f4c4951554944495459556e697377617056323a20494e53554646494349454e545f4c49515549444954595f4255524e4544556e697377617056323a20494e53554646494349454e545f4c49515549444954595f4d494e544544a265627a7a723158207dca18479e58487606bf70c79e44d8dee62353c9ee6d01f9a9d70885b8765f2264736f6c63430005100032454950373132446f6d61696e28737472696e67206e616d652c737472696e672076657273696f6e2c75696e7432353620636861696e49642c6164647265737320766572696679696e67436f6e747261637429",
}

// UniswapV2PairABI is the input ABI used to generate the binding from.
// Deprecated: Use UniswapV2PairMetaData.ABI instead.
var UniswapV2PairABI = UniswapV2PairMetaData.ABI

// UniswapV2PairBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const UniswapV2PairBinRuntime = `608060405234801561001057600080fd5b50600436106101b95760003560e01c80636a627842116100f9578063ba9a7a5611610097578063d21220a711610071578063d21220a7146105da578063d505accf146105e2578063dd62ed3e14610640578063fff6cae91461067b576101b9565b8063ba9a7a5614610597578063bc25cf771461059f578063c45a0155146105d2576101b9565b80637ecebe00116100d35780637ecebe00146104d757806389afcb441461050a57806395d89b4114610556578063a9059cbb1461055e576101b9565b80636a6278421461046957806370a082311461049c5780637464fc3d146104cf576101b9565b806323b872dd116101665780633644e515116101405780633644e51514610416578063485cc9551461041e5780635909c0d5146104595780635a3d549314610461576101b9565b806323b872dd146103ad57806330adf81f146103f0578063313ce567146103f8576101b9565b8063095ea7b311610197578063095ea7b3146103155780630dfe16811461036257806318160ddd14610393576101b9565b8063022c0d9f146101be57806306fdde03146102595780630902f1ac146102d6575b600080fd5b610257600480360360808110156101d457600080fd5b81359160208101359173ffffffffffffffffffffffffffffffffffffffff604083013516919081019060808101606082013564010000000081111561021857600080fd5b82018360208201111561022a57600080fd5b8035906020019184600183028401116401000000008311171561024c57600080fd5b509092509050610683565b005b610261610d57565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561029b578181015183820152602001610283565b50505050905090810190601f1680156102c85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6102de610d90565b604080516dffffffffffffffffffffffffffff948516815292909316602083015263ffffffff168183015290519081900360600190f35b61034e6004803603604081101561032b57600080fd5b5073ffffffffffffffffffffffffffffffffffffffff8135169060200135610de5565b604080519115158252519081900360200190f35b61036a610dfc565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b61039b610e18565b60408051918252519081900360200190f35b61034e600480360360608110156103c357600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060400135610e1e565b61039b610efd565b610400610f21565b6040805160ff9092168252519081900360200190f35b61039b610f26565b6102576004803603604081101561043457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516610f2c565b61039b611005565b61039b61100b565b61039b6004803603602081101561047f57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff16611011565b61039b600480360360208110156104b257600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113cb565b61039b6113dd565b61039b600480360360208110156104ed57600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113e3565b61053d6004803603602081101561052057600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166113f5565b6040805192835260208301919091528051918290030190f35b610261611892565b61034e6004803603604081101561057457600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81351690602001356118cb565b61039b6118d8565b610257600480360360208110156105b557600080fd5b503573ffffffffffffffffffffffffffffffffffffffff166118de565b61036a611ad4565b61036a611af0565b610257600480360360e08110156105f857600080fd5b5073ffffffffffffffffffffffffffffffffffffffff813581169160208101359091169060408101359060608101359060ff6080820135169060a08101359060c00135611b0c565b61039b6004803603604081101561065657600080fd5b5073ffffffffffffffffffffffffffffffffffffffff81358116916020013516611dd8565b610257611df5565b600c546001146106f457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55841515806107075750600084115b61075c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526025815260200180612b2f6025913960400191505060405180910390fd5b600080610767610d90565b5091509150816dffffffffffffffffffffffffffff168710801561079a5750806dffffffffffffffffffffffffffff1686105b6107ef576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526021815260200180612b786021913960400191505060405180910390fd5b600654600754600091829173ffffffffffffffffffffffffffffffffffffffff91821691908116908916821480159061085457508073ffffffffffffffffffffffffffffffffffffffff168973ffffffffffffffffffffffffffffffffffffffff1614155b6108bf57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f556e697377617056323a20494e56414c49445f544f0000000000000000000000604482015290519081900360640190fd5b8a156108d0576108d0828a8d611fdb565b89156108e1576108e1818a8c611fdb565b86156109c3578873ffffffffffffffffffffffffffffffffffffffff166310d1e85c338d8d8c8c6040518663ffffffff1660e01b8152600401808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001858152602001848152602001806020018281038252848482818152602001925080828437600081840152601f19601f8201169050808301925050509650505050505050600060405180830381600087803b1580156109aa57600080fd5b505af11580156109be573d6000803e3d6000fd5b505050505b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff8416916370a08231916024808301926020929190829003018186803b158015610a2f57600080fd5b505afa158015610a43573d6000803e3d6000fd5b505050506040513d6020811015610a5957600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191955073ffffffffffffffffffffffffffffffffffffffff8316916370a0823191602480820192602092909190829003018186803b158015610acb57600080fd5b505afa158015610adf573d6000803e3d6000fd5b505050506040513d6020811015610af557600080fd5b5051925060009150506dffffffffffffffffffffffffffff85168a90038311610b1f576000610b35565b89856dffffffffffffffffffffffffffff160383035b9050600089856dffffffffffffffffffffffffffff16038311610b59576000610b6f565b89856dffffffffffffffffffffffffffff160383035b90506000821180610b805750600081115b610bd5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526024815260200180612b546024913960400191505060405180910390fd5b6000610c09610beb84600363ffffffff6121e816565b610bfd876103e863ffffffff6121e816565b9063ffffffff61226e16565b90506000610c21610beb84600363ffffffff6121e816565b9050610c59620f4240610c4d6dffffffffffffffffffffffffffff8b8116908b1663ffffffff6121e816565b9063ffffffff6121e816565b610c69838363ffffffff6121e816565b1015610cd657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152600c60248201527f556e697377617056323a204b0000000000000000000000000000000000000000604482015290519081900360640190fd5b5050610ce4848488886122e0565b60408051838152602081018390528082018d9052606081018c9052905173ffffffffffffffffffffffffffffffffffffffff8b169133917fd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d8229181900360800190a350506001600c55505050505050505050565b6040518060400160405280600a81526020017f556e69737761702056320000000000000000000000000000000000000000000081525081565b6008546dffffffffffffffffffffffffffff808216926e0100000000000000000000000000008304909116917c0100000000000000000000000000000000000000000000000000000000900463ffffffff1690565b6000610df233848461259c565b5060015b92915050565b60065473ffffffffffffffffffffffffffffffffffffffff1681565b60005481565b73ffffffffffffffffffffffffffffffffffffffff831660009081526002602090815260408083203384529091528120547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff14610ee85773ffffffffffffffffffffffffffffffffffffffff84166000908152600260209081526040808320338452909152902054610eb6908363ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff851660009081526002602090815260408083203384529091529020555b610ef384848461260b565b5060019392505050565b7f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c981565b601281565b60035481565b60055473ffffffffffffffffffffffffffffffffffffffff163314610fb257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f556e697377617056323a20464f5242494444454e000000000000000000000000604482015290519081900360640190fd5b6006805473ffffffffffffffffffffffffffffffffffffffff9384167fffffffffffffffffffffffff00000000000000000000000000000000000000009182161790915560078054929093169116179055565b60095481565b600a5481565b6000600c5460011461108457604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c81905580611094610d90565b50600654604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905193955091935060009273ffffffffffffffffffffffffffffffffffffffff909116916370a08231916024808301926020929190829003018186803b15801561110e57600080fd5b505afa158015611122573d6000803e3d6000fd5b505050506040513d602081101561113857600080fd5b5051600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905192935060009273ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b1580156111b157600080fd5b505afa1580156111c5573d6000803e3d6000fd5b505050506040513d60208110156111db57600080fd5b505190506000611201836dffffffffffffffffffffffffffff871663ffffffff61226e16565b90506000611225836dffffffffffffffffffffffffffff871663ffffffff61226e16565b9050600061123387876126ec565b600054909150806112705761125c6103e8610bfd611257878763ffffffff6121e816565b612878565b985061126b60006103e86128ca565b6112cd565b6112ca6dffffffffffffffffffffffffffff8916611294868463ffffffff6121e816565b8161129b57fe5b046dffffffffffffffffffffffffffff89166112bd868563ffffffff6121e816565b816112c457fe5b0461297a565b98505b60008911611326576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180612bc16028913960400191505060405180910390fd5b6113308a8a6128ca565b61133c86868a8a6122e0565b811561137e5760085461137a906dffffffffffffffffffffffffffff808216916e01000000000000000000000000000090041663ffffffff6121e816565b600b555b6040805185815260208101859052815133927f4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f928290030190a250506001600c5550949695505050505050565b60016020526000908152604090205481565b600b5481565b60046020526000908152604090205481565b600080600c5460011461146957604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c81905580611479610d90565b50600654600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905194965092945073ffffffffffffffffffffffffffffffffffffffff9182169391169160009184916370a08231916024808301926020929190829003018186803b1580156114fb57600080fd5b505afa15801561150f573d6000803e3d6000fd5b505050506040513d602081101561152557600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191925060009173ffffffffffffffffffffffffffffffffffffffff8516916370a08231916024808301926020929190829003018186803b15801561159957600080fd5b505afa1580156115ad573d6000803e3d6000fd5b505050506040513d60208110156115c357600080fd5b5051306000908152600160205260408120549192506115e288886126ec565b600054909150806115f9848763ffffffff6121e816565b8161160057fe5b049a5080611614848663ffffffff6121e816565b8161161b57fe5b04995060008b11801561162e575060008a115b611683576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526028815260200180612b996028913960400191505060405180910390fd5b61168d3084612992565b611698878d8d611fdb565b6116a3868d8c611fdb565b604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff8916916370a08231916024808301926020929190829003018186803b15801561170f57600080fd5b505afa158015611723573d6000803e3d6000fd5b505050506040513d602081101561173957600080fd5b5051604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905191965073ffffffffffffffffffffffffffffffffffffffff8816916370a0823191602480820192602092909190829003018186803b1580156117ab57600080fd5b505afa1580156117bf573d6000803e3d6000fd5b505050506040513d60208110156117d557600080fd5b505193506117e585858b8b6122e0565b811561182757600854611823906dffffffffffffffffffffffffffff808216916e01000000000000000000000000000090041663ffffffff6121e816565b600b555b604080518c8152602081018c9052815173ffffffffffffffffffffffffffffffffffffffff8f169233927fdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496929081900390910190a35050505050505050506001600c81905550915091565b6040518060400160405280600681526020017f554e492d5632000000000000000000000000000000000000000000000000000081525081565b6000610df233848461260b565b6103e881565b600c5460011461194f57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55600654600754600854604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff9485169490931692611a2b9285928792611a26926dffffffffffffffffffffffffffff169185916370a0823191602480820192602092909190829003018186803b1580156119ee57600080fd5b505afa158015611a02573d6000803e3d6000fd5b505050506040513d6020811015611a1857600080fd5b50519063ffffffff61226e16565b611fdb565b600854604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529051611aca9284928792611a26926e01000000000000000000000000000090046dffffffffffffffffffffffffffff169173ffffffffffffffffffffffffffffffffffffffff8616916370a0823191602480820192602092909190829003018186803b1580156119ee57600080fd5b50506001600c5550565b60055473ffffffffffffffffffffffffffffffffffffffff1681565b60075473ffffffffffffffffffffffffffffffffffffffff1681565b42841015611b7b57604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601260248201527f556e697377617056323a20455850495245440000000000000000000000000000604482015290519081900360640190fd5b60035473ffffffffffffffffffffffffffffffffffffffff80891660008181526004602090815260408083208054600180820190925582517f6e71edae12b1b97f4d1f60370fef10105fa2faae0126114a169c64845d6126c98186015280840196909652958d166060860152608085018c905260a085019590955260c08085018b90528151808603909101815260e0850182528051908301207f19010000000000000000000000000000000000000000000000000000000000006101008601526101028501969096526101228085019690965280518085039096018652610142840180825286519683019690962095839052610162840180825286905260ff89166101828501526101a284018890526101c28401879052519193926101e2808201937fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081019281900390910190855afa158015611cdc573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff811615801590611d5757508873ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16145b611dc257604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f556e697377617056323a20494e56414c49445f5349474e415455524500000000604482015290519081900360640190fd5b611dcd89898961259c565b505050505050505050565b600260209081526000928352604080842090915290825290205481565b600c54600114611e6657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f556e697377617056323a204c4f434b4544000000000000000000000000000000604482015290519081900360640190fd5b6000600c55600654604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201529051611fd49273ffffffffffffffffffffffffffffffffffffffff16916370a08231916024808301926020929190829003018186803b158015611edd57600080fd5b505afa158015611ef1573d6000803e3d6000fd5b505050506040513d6020811015611f0757600080fd5b5051600754604080517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152905173ffffffffffffffffffffffffffffffffffffffff909216916370a0823191602480820192602092909190829003018186803b158015611f7a57600080fd5b505afa158015611f8e573d6000803e3d6000fd5b505050506040513d6020811015611fa457600080fd5b50516008546dffffffffffffffffffffffffffff808216916e0100000000000000000000000000009004166122e0565b6001600c55565b604080518082018252601981527f7472616e7366657228616464726573732c75696e743235362900000000000000602091820152815173ffffffffffffffffffffffffffffffffffffffff85811660248301526044808301869052845180840390910181526064909201845291810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb000000000000000000000000000000000000000000000000000000001781529251815160009460609489169392918291908083835b602083106120e157805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016120a4565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d8060008114612143576040519150601f19603f3d011682016040523d82523d6000602084013e612148565b606091505b5091509150818015612176575080511580612176575080806020019051602081101561217357600080fd5b50515b6121e157604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f556e697377617056323a205452414e534645525f4641494c4544000000000000604482015290519081900360640190fd5b5050505050565b60008115806122035750508082028282828161220057fe5b04145b610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f64732d6d6174682d6d756c2d6f766572666c6f77000000000000000000000000604482015290519081900360640190fd5b80820382811115610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f64732d6d6174682d7375622d756e646572666c6f770000000000000000000000604482015290519081900360640190fd5b6dffffffffffffffffffffffffffff841180159061230c57506dffffffffffffffffffffffffffff8311155b61237757604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f556e697377617056323a204f564552464c4f5700000000000000000000000000604482015290519081900360640190fd5b60085463ffffffff428116917c0100000000000000000000000000000000000000000000000000000000900481168203908116158015906123c757506dffffffffffffffffffffffffffff841615155b80156123e257506dffffffffffffffffffffffffffff831615155b15612492578063ffffffff16612425856123fb86612a57565b7bffffffffffffffffffffffffffffffffffffffffffffffffffffffff169063ffffffff612a7b16565b600980547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff929092169290920201905563ffffffff8116612465846123fb87612a57565b600a80547bffffffffffffffffffffffffffffffffffffffffffffffffffffffff92909216929092020190555b600880547fffffffffffffffffffffffffffffffffffff0000000000000000000000000000166dffffffffffffffffffffffffffff888116919091177fffffffff0000000000000000000000000000ffffffffffffffffffffffffffff166e0100000000000000000000000000008883168102919091177bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167c010000000000000000000000000000000000000000000000000000000063ffffffff871602179283905560408051848416815291909304909116602082015281517f1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1929181900390910190a1505050505050565b73ffffffffffffffffffffffffffffffffffffffff808416600081815260026020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b73ffffffffffffffffffffffffffffffffffffffff8316600090815260016020526040902054612641908263ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff8085166000908152600160205260408082209390935590841681522054612683908263ffffffff612abc16565b73ffffffffffffffffffffffffffffffffffffffff80841660008181526001602090815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b600080600560009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663017e7e586040518163ffffffff1660e01b815260040160206040518083038186803b15801561275757600080fd5b505afa15801561276b573d6000803e3d6000fd5b505050506040513d602081101561278157600080fd5b5051600b5473ffffffffffffffffffffffffffffffffffffffff821615801594509192509061286457801561285f5760006127d86112576dffffffffffffffffffffffffffff88811690881663ffffffff6121e816565b905060006127e583612878565b90508082111561285c576000612813612804848463ffffffff61226e16565b6000549063ffffffff6121e816565b905060006128388361282c86600563ffffffff6121e816565b9063ffffffff612abc16565b9050600081838161284557fe5b04905080156128585761285887826128ca565b5050505b50505b612870565b8015612870576000600b555b505092915050565b600060038211156128bb575080600160028204015b818110156128b5578091506002818285816128a457fe5b0401816128ad57fe5b04905061288d565b506128c5565b81156128c5575060015b919050565b6000546128dd908263ffffffff612abc16565b600090815573ffffffffffffffffffffffffffffffffffffffff8316815260016020526040902054612915908263ffffffff612abc16565b73ffffffffffffffffffffffffffffffffffffffff831660008181526001602090815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b6000818310612989578161298b565b825b9392505050565b73ffffffffffffffffffffffffffffffffffffffff82166000908152600160205260409020546129c8908263ffffffff61226e16565b73ffffffffffffffffffffffffffffffffffffffff831660009081526001602052604081209190915554612a02908263ffffffff61226e16565b600090815560408051838152905173ffffffffffffffffffffffffffffffffffffffff8516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef919081900360200190a35050565b6dffffffffffffffffffffffffffff166e0100000000000000000000000000000290565b60006dffffffffffffffffffffffffffff82167bffffffffffffffffffffffffffffffffffffffffffffffffffffffff841681612ab457fe5b049392505050565b80820182811015610df657604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601460248201527f64732d6d6174682d6164642d6f766572666c6f77000000000000000000000000604482015290519081900360640190fdfe556e697377617056323a20494e53554646494349454e545f4f55545055545f414d4f554e54556e697377617056323a20494e53554646494349454e545f494e5055545f414d4f554e54556e697377617056323a20494e53554646494349454e545f4c4951554944495459556e697377617056323a20494e53554646494349454e545f4c49515549444954595f4255524e4544556e697377617056323a20494e53554646494349454e545f4c49515549444954595f4d494e544544a265627a7a723158207dca18479e58487606bf70c79e44d8dee62353c9ee6d01f9a9d70885b8765f2264736f6c63430005100032`

// Deprecated: Use UniswapV2PairMetaData.Sigs instead.
// UniswapV2PairFuncSigs maps the 4-byte function signature to its string representation.
var UniswapV2PairFuncSigs = UniswapV2PairMetaData.Sigs

// UniswapV2PairBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use UniswapV2PairMetaData.Bin instead.
var UniswapV2PairBin = UniswapV2PairMetaData.Bin

// DeployUniswapV2Pair deploys a new Kaia contract, binding an instance of UniswapV2Pair to it.
func DeployUniswapV2Pair(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *UniswapV2Pair, error) {
	parsed, err := UniswapV2PairMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(UniswapV2PairBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &UniswapV2Pair{UniswapV2PairCaller: UniswapV2PairCaller{contract: contract}, UniswapV2PairTransactor: UniswapV2PairTransactor{contract: contract}, UniswapV2PairFilterer: UniswapV2PairFilterer{contract: contract}}, nil
}

// UniswapV2Pair is an auto generated Go binding around a Kaia contract.
type UniswapV2Pair struct {
	UniswapV2PairCaller     // Read-only binding to the contract
	UniswapV2PairTransactor // Write-only binding to the contract
	UniswapV2PairFilterer   // Log filterer for contract events
}

// UniswapV2PairCaller is an auto generated read-only Go binding around a Kaia contract.
type UniswapV2PairCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2PairTransactor is an auto generated write-only Go binding around a Kaia contract.
type UniswapV2PairTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2PairFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type UniswapV2PairFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UniswapV2PairSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type UniswapV2PairSession struct {
	Contract     *UniswapV2Pair    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UniswapV2PairCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type UniswapV2PairCallerSession struct {
	Contract *UniswapV2PairCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// UniswapV2PairTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type UniswapV2PairTransactorSession struct {
	Contract     *UniswapV2PairTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// UniswapV2PairRaw is an auto generated low-level Go binding around a Kaia contract.
type UniswapV2PairRaw struct {
	Contract *UniswapV2Pair // Generic contract binding to access the raw methods on
}

// UniswapV2PairCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type UniswapV2PairCallerRaw struct {
	Contract *UniswapV2PairCaller // Generic read-only contract binding to access the raw methods on
}

// UniswapV2PairTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type UniswapV2PairTransactorRaw struct {
	Contract *UniswapV2PairTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapV2Pair creates a new instance of UniswapV2Pair, bound to a specific deployed contract.
func NewUniswapV2Pair(address common.Address, backend bind.ContractBackend) (*UniswapV2Pair, error) {
	contract, err := bindUniswapV2Pair(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UniswapV2Pair{UniswapV2PairCaller: UniswapV2PairCaller{contract: contract}, UniswapV2PairTransactor: UniswapV2PairTransactor{contract: contract}, UniswapV2PairFilterer: UniswapV2PairFilterer{contract: contract}}, nil
}

// NewUniswapV2PairCaller creates a new read-only instance of UniswapV2Pair, bound to a specific deployed contract.
func NewUniswapV2PairCaller(address common.Address, caller bind.ContractCaller) (*UniswapV2PairCaller, error) {
	contract, err := bindUniswapV2Pair(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairCaller{contract: contract}, nil
}

// NewUniswapV2PairTransactor creates a new write-only instance of UniswapV2Pair, bound to a specific deployed contract.
func NewUniswapV2PairTransactor(address common.Address, transactor bind.ContractTransactor) (*UniswapV2PairTransactor, error) {
	contract, err := bindUniswapV2Pair(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairTransactor{contract: contract}, nil
}

// NewUniswapV2PairFilterer creates a new log filterer instance of UniswapV2Pair, bound to a specific deployed contract.
func NewUniswapV2PairFilterer(address common.Address, filterer bind.ContractFilterer) (*UniswapV2PairFilterer, error) {
	contract, err := bindUniswapV2Pair(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairFilterer{contract: contract}, nil
}

// bindUniswapV2Pair binds a generic wrapper to an already deployed contract.
func bindUniswapV2Pair(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UniswapV2PairMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Pair *UniswapV2PairRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Pair.Contract.UniswapV2PairCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Pair *UniswapV2PairRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.UniswapV2PairTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Pair *UniswapV2PairRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.UniswapV2PairTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UniswapV2Pair *UniswapV2PairCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UniswapV2Pair.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UniswapV2Pair *UniswapV2PairTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UniswapV2Pair *UniswapV2PairTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.contract.Transact(opts, method, params...)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_UniswapV2Pair *UniswapV2PairCaller) DOMAINSEPARATOR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "DOMAIN_SEPARATOR")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_UniswapV2Pair *UniswapV2PairSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _UniswapV2Pair.Contract.DOMAINSEPARATOR(&_UniswapV2Pair.CallOpts)
}

// DOMAINSEPARATOR is a free data retrieval call binding the contract method 0x3644e515.
//
// Solidity: function DOMAIN_SEPARATOR() view returns(bytes32)
func (_UniswapV2Pair *UniswapV2PairCallerSession) DOMAINSEPARATOR() ([32]byte, error) {
	return _UniswapV2Pair.Contract.DOMAINSEPARATOR(&_UniswapV2Pair.CallOpts)
}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCaller) MINIMUMLIQUIDITY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "MINIMUM_LIQUIDITY")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairSession) MINIMUMLIQUIDITY() (*big.Int, error) {
	return _UniswapV2Pair.Contract.MINIMUMLIQUIDITY(&_UniswapV2Pair.CallOpts)
}

// MINIMUMLIQUIDITY is a free data retrieval call binding the contract method 0xba9a7a56.
//
// Solidity: function MINIMUM_LIQUIDITY() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCallerSession) MINIMUMLIQUIDITY() (*big.Int, error) {
	return _UniswapV2Pair.Contract.MINIMUMLIQUIDITY(&_UniswapV2Pair.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_UniswapV2Pair *UniswapV2PairCaller) PERMITTYPEHASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "PERMIT_TYPEHASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_UniswapV2Pair *UniswapV2PairSession) PERMITTYPEHASH() ([32]byte, error) {
	return _UniswapV2Pair.Contract.PERMITTYPEHASH(&_UniswapV2Pair.CallOpts)
}

// PERMITTYPEHASH is a free data retrieval call binding the contract method 0x30adf81f.
//
// Solidity: function PERMIT_TYPEHASH() view returns(bytes32)
func (_UniswapV2Pair *UniswapV2PairCallerSession) PERMITTYPEHASH() ([32]byte, error) {
	return _UniswapV2Pair.Contract.PERMITTYPEHASH(&_UniswapV2Pair.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "allowance", arg0, arg1)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _UniswapV2Pair.Contract.Allowance(&_UniswapV2Pair.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _UniswapV2Pair.Contract.Allowance(&_UniswapV2Pair.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "balanceOf", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _UniswapV2Pair.Contract.BalanceOf(&_UniswapV2Pair.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _UniswapV2Pair.Contract.BalanceOf(&_UniswapV2Pair.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_UniswapV2Pair *UniswapV2PairCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "decimals")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_UniswapV2Pair *UniswapV2PairSession) Decimals() (uint8, error) {
	return _UniswapV2Pair.Contract.Decimals(&_UniswapV2Pair.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Decimals() (uint8, error) {
	return _UniswapV2Pair.Contract.Decimals(&_UniswapV2Pair.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV2Pair *UniswapV2PairCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "factory")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV2Pair *UniswapV2PairSession) Factory() (common.Address, error) {
	return _UniswapV2Pair.Contract.Factory(&_UniswapV2Pair.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() view returns(address)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Factory() (common.Address, error) {
	return _UniswapV2Pair.Contract.Factory(&_UniswapV2Pair.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 _reserve0, uint112 _reserve1, uint32 _blockTimestampLast)
func (_UniswapV2Pair *UniswapV2PairCaller) GetReserves(opts *bind.CallOpts) (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error,
) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "getReserves")

	outstruct := new(struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Reserve0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Reserve1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.BlockTimestampLast = *abi.ConvertType(out[2], new(uint32)).(*uint32)

	return *outstruct, err
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 _reserve0, uint112 _reserve1, uint32 _blockTimestampLast)
func (_UniswapV2Pair *UniswapV2PairSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error,
) {
	return _UniswapV2Pair.Contract.GetReserves(&_UniswapV2Pair.CallOpts)
}

// GetReserves is a free data retrieval call binding the contract method 0x0902f1ac.
//
// Solidity: function getReserves() view returns(uint112 _reserve0, uint112 _reserve1, uint32 _blockTimestampLast)
func (_UniswapV2Pair *UniswapV2PairCallerSession) GetReserves() (struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}, error,
) {
	return _UniswapV2Pair.Contract.GetReserves(&_UniswapV2Pair.CallOpts)
}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCaller) KLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "kLast")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairSession) KLast() (*big.Int, error) {
	return _UniswapV2Pair.Contract.KLast(&_UniswapV2Pair.CallOpts)
}

// KLast is a free data retrieval call binding the contract method 0x7464fc3d.
//
// Solidity: function kLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCallerSession) KLast() (*big.Int, error) {
	return _UniswapV2Pair.Contract.KLast(&_UniswapV2Pair.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_UniswapV2Pair *UniswapV2PairCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "name")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_UniswapV2Pair *UniswapV2PairSession) Name() (string, error) {
	return _UniswapV2Pair.Contract.Name(&_UniswapV2Pair.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Name() (string, error) {
	return _UniswapV2Pair.Contract.Name(&_UniswapV2Pair.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "nonces", arg0)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _UniswapV2Pair.Contract.Nonces(&_UniswapV2Pair.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _UniswapV2Pair.Contract.Nonces(&_UniswapV2Pair.CallOpts, arg0)
}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCaller) Price0CumulativeLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "price0CumulativeLast")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairSession) Price0CumulativeLast() (*big.Int, error) {
	return _UniswapV2Pair.Contract.Price0CumulativeLast(&_UniswapV2Pair.CallOpts)
}

// Price0CumulativeLast is a free data retrieval call binding the contract method 0x5909c0d5.
//
// Solidity: function price0CumulativeLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Price0CumulativeLast() (*big.Int, error) {
	return _UniswapV2Pair.Contract.Price0CumulativeLast(&_UniswapV2Pair.CallOpts)
}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCaller) Price1CumulativeLast(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "price1CumulativeLast")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairSession) Price1CumulativeLast() (*big.Int, error) {
	return _UniswapV2Pair.Contract.Price1CumulativeLast(&_UniswapV2Pair.CallOpts)
}

// Price1CumulativeLast is a free data retrieval call binding the contract method 0x5a3d5493.
//
// Solidity: function price1CumulativeLast() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Price1CumulativeLast() (*big.Int, error) {
	return _UniswapV2Pair.Contract.Price1CumulativeLast(&_UniswapV2Pair.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_UniswapV2Pair *UniswapV2PairCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "symbol")
	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_UniswapV2Pair *UniswapV2PairSession) Symbol() (string, error) {
	return _UniswapV2Pair.Contract.Symbol(&_UniswapV2Pair.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Symbol() (string, error) {
	return _UniswapV2Pair.Contract.Symbol(&_UniswapV2Pair.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_UniswapV2Pair *UniswapV2PairCaller) Token0(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "token0")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_UniswapV2Pair *UniswapV2PairSession) Token0() (common.Address, error) {
	return _UniswapV2Pair.Contract.Token0(&_UniswapV2Pair.CallOpts)
}

// Token0 is a free data retrieval call binding the contract method 0x0dfe1681.
//
// Solidity: function token0() view returns(address)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Token0() (common.Address, error) {
	return _UniswapV2Pair.Contract.Token0(&_UniswapV2Pair.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_UniswapV2Pair *UniswapV2PairCaller) Token1(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "token1")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_UniswapV2Pair *UniswapV2PairSession) Token1() (common.Address, error) {
	return _UniswapV2Pair.Contract.Token1(&_UniswapV2Pair.CallOpts)
}

// Token1 is a free data retrieval call binding the contract method 0xd21220a7.
//
// Solidity: function token1() view returns(address)
func (_UniswapV2Pair *UniswapV2PairCallerSession) Token1() (common.Address, error) {
	return _UniswapV2Pair.Contract.Token1(&_UniswapV2Pair.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _UniswapV2Pair.contract.Call(opts, &out, "totalSupply")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairSession) TotalSupply() (*big.Int, error) {
	return _UniswapV2Pair.Contract.TotalSupply(&_UniswapV2Pair.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_UniswapV2Pair *UniswapV2PairCallerSession) TotalSupply() (*big.Int, error) {
	return _UniswapV2Pair.Contract.TotalSupply(&_UniswapV2Pair.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Approve(&_UniswapV2Pair.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Approve(&_UniswapV2Pair.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_UniswapV2Pair *UniswapV2PairTransactor) Burn(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "burn", to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_UniswapV2Pair *UniswapV2PairSession) Burn(to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Burn(&_UniswapV2Pair.TransactOpts, to)
}

// Burn is a paid mutator transaction binding the contract method 0x89afcb44.
//
// Solidity: function burn(address to) returns(uint256 amount0, uint256 amount1)
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Burn(to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Burn(&_UniswapV2Pair.TransactOpts, to)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _token0, address _token1) returns()
func (_UniswapV2Pair *UniswapV2PairTransactor) Initialize(opts *bind.TransactOpts, _token0 common.Address, _token1 common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "initialize", _token0, _token1)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _token0, address _token1) returns()
func (_UniswapV2Pair *UniswapV2PairSession) Initialize(_token0 common.Address, _token1 common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Initialize(&_UniswapV2Pair.TransactOpts, _token0, _token1)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _token0, address _token1) returns()
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Initialize(_token0 common.Address, _token1 common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Initialize(&_UniswapV2Pair.TransactOpts, _token0, _token1)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_UniswapV2Pair *UniswapV2PairTransactor) Mint(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "mint", to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_UniswapV2Pair *UniswapV2PairSession) Mint(to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Mint(&_UniswapV2Pair.TransactOpts, to)
}

// Mint is a paid mutator transaction binding the contract method 0x6a627842.
//
// Solidity: function mint(address to) returns(uint256 liquidity)
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Mint(to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Mint(&_UniswapV2Pair.TransactOpts, to)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_UniswapV2Pair *UniswapV2PairTransactor) Permit(opts *bind.TransactOpts, owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "permit", owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_UniswapV2Pair *UniswapV2PairSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Permit(&_UniswapV2Pair.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Permit is a paid mutator transaction binding the contract method 0xd505accf.
//
// Solidity: function permit(address owner, address spender, uint256 value, uint256 deadline, uint8 v, bytes32 r, bytes32 s) returns()
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Permit(owner common.Address, spender common.Address, value *big.Int, deadline *big.Int, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Permit(&_UniswapV2Pair.TransactOpts, owner, spender, value, deadline, v, r, s)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_UniswapV2Pair *UniswapV2PairTransactor) Skim(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "skim", to)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_UniswapV2Pair *UniswapV2PairSession) Skim(to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Skim(&_UniswapV2Pair.TransactOpts, to)
}

// Skim is a paid mutator transaction binding the contract method 0xbc25cf77.
//
// Solidity: function skim(address to) returns()
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Skim(to common.Address) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Skim(&_UniswapV2Pair.TransactOpts, to)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_UniswapV2Pair *UniswapV2PairTransactor) Swap(opts *bind.TransactOpts, amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "swap", amount0Out, amount1Out, to, data)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_UniswapV2Pair *UniswapV2PairSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Swap(&_UniswapV2Pair.TransactOpts, amount0Out, amount1Out, to, data)
}

// Swap is a paid mutator transaction binding the contract method 0x022c0d9f.
//
// Solidity: function swap(uint256 amount0Out, uint256 amount1Out, address to, bytes data) returns()
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Swap(amount0Out *big.Int, amount1Out *big.Int, to common.Address, data []byte) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Swap(&_UniswapV2Pair.TransactOpts, amount0Out, amount1Out, to, data)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_UniswapV2Pair *UniswapV2PairTransactor) Sync(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "sync")
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_UniswapV2Pair *UniswapV2PairSession) Sync() (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Sync(&_UniswapV2Pair.TransactOpts)
}

// Sync is a paid mutator transaction binding the contract method 0xfff6cae9.
//
// Solidity: function sync() returns()
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Sync() (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Sync(&_UniswapV2Pair.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Transfer(&_UniswapV2Pair.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.Transfer(&_UniswapV2Pair.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.TransferFrom(&_UniswapV2Pair.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_UniswapV2Pair *UniswapV2PairTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _UniswapV2Pair.Contract.TransferFrom(&_UniswapV2Pair.TransactOpts, from, to, value)
}

// UniswapV2PairApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the UniswapV2Pair contract.
type UniswapV2PairApprovalIterator struct {
	Event *UniswapV2PairApproval // Event containing the contract specifics and raw log

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
func (it *UniswapV2PairApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2PairApproval)
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
		it.Event = new(UniswapV2PairApproval)
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
func (it *UniswapV2PairApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2PairApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2PairApproval represents a Approval event raised by the UniswapV2Pair contract.
type UniswapV2PairApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UniswapV2Pair *UniswapV2PairFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*UniswapV2PairApprovalIterator, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairApprovalIterator{contract: _UniswapV2Pair.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UniswapV2Pair *UniswapV2PairFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *UniswapV2PairApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2PairApproval)
				if err := _UniswapV2Pair.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_UniswapV2Pair *UniswapV2PairFilterer) ParseApproval(log types.Log) (*UniswapV2PairApproval, error) {
	event := new(UniswapV2PairApproval)
	if err := _UniswapV2Pair.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV2PairBurnIterator is returned from FilterBurn and is used to iterate over the raw logs and unpacked data for Burn events raised by the UniswapV2Pair contract.
type UniswapV2PairBurnIterator struct {
	Event *UniswapV2PairBurn // Event containing the contract specifics and raw log

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
func (it *UniswapV2PairBurnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2PairBurn)
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
		it.Event = new(UniswapV2PairBurn)
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
func (it *UniswapV2PairBurnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2PairBurnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2PairBurn represents a Burn event raised by the UniswapV2Pair contract.
type UniswapV2PairBurn struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	To      common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBurn is a free log retrieval operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_UniswapV2Pair *UniswapV2PairFilterer) FilterBurn(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*UniswapV2PairBurnIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.FilterLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairBurnIterator{contract: _UniswapV2Pair.contract, event: "Burn", logs: logs, sub: sub}, nil
}

// WatchBurn is a free log subscription operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_UniswapV2Pair *UniswapV2PairFilterer) WatchBurn(opts *bind.WatchOpts, sink chan<- *UniswapV2PairBurn, sender []common.Address, to []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.WatchLogs(opts, "Burn", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2PairBurn)
				if err := _UniswapV2Pair.contract.UnpackLog(event, "Burn", log); err != nil {
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

// ParseBurn is a log parse operation binding the contract event 0xdccd412f0b1252819cb1fd330b93224ca42612892bb3f4f789976e6d81936496.
//
// Solidity: event Burn(address indexed sender, uint256 amount0, uint256 amount1, address indexed to)
func (_UniswapV2Pair *UniswapV2PairFilterer) ParseBurn(log types.Log) (*UniswapV2PairBurn, error) {
	event := new(UniswapV2PairBurn)
	if err := _UniswapV2Pair.contract.UnpackLog(event, "Burn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV2PairMintIterator is returned from FilterMint and is used to iterate over the raw logs and unpacked data for Mint events raised by the UniswapV2Pair contract.
type UniswapV2PairMintIterator struct {
	Event *UniswapV2PairMint // Event containing the contract specifics and raw log

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
func (it *UniswapV2PairMintIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2PairMint)
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
		it.Event = new(UniswapV2PairMint)
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
func (it *UniswapV2PairMintIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2PairMintIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2PairMint represents a Mint event raised by the UniswapV2Pair contract.
type UniswapV2PairMint struct {
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMint is a free log retrieval operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_UniswapV2Pair *UniswapV2PairFilterer) FilterMint(opts *bind.FilterOpts, sender []common.Address) (*UniswapV2PairMintIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.FilterLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairMintIterator{contract: _UniswapV2Pair.contract, event: "Mint", logs: logs, sub: sub}, nil
}

// WatchMint is a free log subscription operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_UniswapV2Pair *UniswapV2PairFilterer) WatchMint(opts *bind.WatchOpts, sink chan<- *UniswapV2PairMint, sender []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.WatchLogs(opts, "Mint", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2PairMint)
				if err := _UniswapV2Pair.contract.UnpackLog(event, "Mint", log); err != nil {
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

// ParseMint is a log parse operation binding the contract event 0x4c209b5fc8ad50758f13e2e1088ba56a560dff690a1c6fef26394f4c03821c4f.
//
// Solidity: event Mint(address indexed sender, uint256 amount0, uint256 amount1)
func (_UniswapV2Pair *UniswapV2PairFilterer) ParseMint(log types.Log) (*UniswapV2PairMint, error) {
	event := new(UniswapV2PairMint)
	if err := _UniswapV2Pair.contract.UnpackLog(event, "Mint", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV2PairSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the UniswapV2Pair contract.
type UniswapV2PairSwapIterator struct {
	Event *UniswapV2PairSwap // Event containing the contract specifics and raw log

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
func (it *UniswapV2PairSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2PairSwap)
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
		it.Event = new(UniswapV2PairSwap)
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
func (it *UniswapV2PairSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2PairSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2PairSwap represents a Swap event raised by the UniswapV2Pair contract.
type UniswapV2PairSwap struct {
	Sender     common.Address
	Amount0In  *big.Int
	Amount1In  *big.Int
	Amount0Out *big.Int
	Amount1Out *big.Int
	To         common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_UniswapV2Pair *UniswapV2PairFilterer) FilterSwap(opts *bind.FilterOpts, sender []common.Address, to []common.Address) (*UniswapV2PairSwapIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.FilterLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairSwapIterator{contract: _UniswapV2Pair.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_UniswapV2Pair *UniswapV2PairFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *UniswapV2PairSwap, sender []common.Address, to []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.WatchLogs(opts, "Swap", senderRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2PairSwap)
				if err := _UniswapV2Pair.contract.UnpackLog(event, "Swap", log); err != nil {
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

// ParseSwap is a log parse operation binding the contract event 0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822.
//
// Solidity: event Swap(address indexed sender, uint256 amount0In, uint256 amount1In, uint256 amount0Out, uint256 amount1Out, address indexed to)
func (_UniswapV2Pair *UniswapV2PairFilterer) ParseSwap(log types.Log) (*UniswapV2PairSwap, error) {
	event := new(UniswapV2PairSwap)
	if err := _UniswapV2Pair.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV2PairSyncIterator is returned from FilterSync and is used to iterate over the raw logs and unpacked data for Sync events raised by the UniswapV2Pair contract.
type UniswapV2PairSyncIterator struct {
	Event *UniswapV2PairSync // Event containing the contract specifics and raw log

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
func (it *UniswapV2PairSyncIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2PairSync)
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
		it.Event = new(UniswapV2PairSync)
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
func (it *UniswapV2PairSyncIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2PairSyncIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2PairSync represents a Sync event raised by the UniswapV2Pair contract.
type UniswapV2PairSync struct {
	Reserve0 *big.Int
	Reserve1 *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSync is a free log retrieval operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_UniswapV2Pair *UniswapV2PairFilterer) FilterSync(opts *bind.FilterOpts) (*UniswapV2PairSyncIterator, error) {
	logs, sub, err := _UniswapV2Pair.contract.FilterLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairSyncIterator{contract: _UniswapV2Pair.contract, event: "Sync", logs: logs, sub: sub}, nil
}

// WatchSync is a free log subscription operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_UniswapV2Pair *UniswapV2PairFilterer) WatchSync(opts *bind.WatchOpts, sink chan<- *UniswapV2PairSync) (event.Subscription, error) {
	logs, sub, err := _UniswapV2Pair.contract.WatchLogs(opts, "Sync")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2PairSync)
				if err := _UniswapV2Pair.contract.UnpackLog(event, "Sync", log); err != nil {
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

// ParseSync is a log parse operation binding the contract event 0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1.
//
// Solidity: event Sync(uint112 reserve0, uint112 reserve1)
func (_UniswapV2Pair *UniswapV2PairFilterer) ParseSync(log types.Log) (*UniswapV2PairSync, error) {
	event := new(UniswapV2PairSync)
	if err := _UniswapV2Pair.contract.UnpackLog(event, "Sync", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UniswapV2PairTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the UniswapV2Pair contract.
type UniswapV2PairTransferIterator struct {
	Event *UniswapV2PairTransfer // Event containing the contract specifics and raw log

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
func (it *UniswapV2PairTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UniswapV2PairTransfer)
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
		it.Event = new(UniswapV2PairTransfer)
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
func (it *UniswapV2PairTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UniswapV2PairTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UniswapV2PairTransfer represents a Transfer event raised by the UniswapV2Pair contract.
type UniswapV2PairTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UniswapV2Pair *UniswapV2PairFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*UniswapV2PairTransferIterator, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &UniswapV2PairTransferIterator{contract: _UniswapV2Pair.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UniswapV2Pair *UniswapV2PairFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *UniswapV2PairTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _UniswapV2Pair.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UniswapV2PairTransfer)
				if err := _UniswapV2Pair.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_UniswapV2Pair *UniswapV2PairFilterer) ParseTransfer(log types.Log) (*UniswapV2PairTransfer, error) {
	event := new(UniswapV2PairTransfer)
	if err := _UniswapV2Pair.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
