// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package publicdelegation

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

// IPublicDelegationPDConstructorArgs is an auto generated low-level Go binding around an user-defined struct.
type IPublicDelegationPDConstructorArgs struct {
	Owner          common.Address
	CommissionTo   common.Address
	CommissionRate *big.Int
	GcName         string
}

// AddressMetaData contains all meta data concerning the Address contract.
var AddressMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"}]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea26469706673582212206d20bff17dc7a186180f241e0975b1633a8914ef132d046516045de0e07e8d0c64736f6c63430008190033",
}

// AddressABI is the input ABI used to generate the binding from.
// Deprecated: Use AddressMetaData.ABI instead.
var AddressABI = AddressMetaData.ABI

// AddressBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea26469706673582212206d20bff17dc7a186180f241e0975b1633a8914ef132d046516045de0e07e8d0c64736f6c63430008190033`

// AddressBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AddressMetaData.Bin instead.
var AddressBin = AddressMetaData.Bin

// DeployAddress deploys a new Kaia contract, binding an instance of Address to it.
func DeployAddress(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Address, error) {
	parsed, err := AddressMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AddressBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Address{AddressCaller: AddressCaller{contract: contract}, AddressTransactor: AddressTransactor{contract: contract}, AddressFilterer: AddressFilterer{contract: contract}}, nil
}

// Address is an auto generated Go binding around a Kaia contract.
type Address struct {
	AddressCaller     // Read-only binding to the contract
	AddressTransactor // Write-only binding to the contract
	AddressFilterer   // Log filterer for contract events
}

// AddressCaller is an auto generated read-only Go binding around a Kaia contract.
type AddressCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressTransactor is an auto generated write-only Go binding around a Kaia contract.
type AddressTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type AddressFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type AddressSession struct {
	Contract     *Address          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AddressCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type AddressCallerSession struct {
	Contract *AddressCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// AddressTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type AddressTransactorSession struct {
	Contract     *AddressTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AddressRaw is an auto generated low-level Go binding around a Kaia contract.
type AddressRaw struct {
	Contract *Address // Generic contract binding to access the raw methods on
}

// AddressCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type AddressCallerRaw struct {
	Contract *AddressCaller // Generic read-only contract binding to access the raw methods on
}

// AddressTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type AddressTransactorRaw struct {
	Contract *AddressTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAddress creates a new instance of Address, bound to a specific deployed contract.
func NewAddress(address common.Address, backend bind.ContractBackend) (*Address, error) {
	contract, err := bindAddress(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Address{AddressCaller: AddressCaller{contract: contract}, AddressTransactor: AddressTransactor{contract: contract}, AddressFilterer: AddressFilterer{contract: contract}}, nil
}

// NewAddressCaller creates a new read-only instance of Address, bound to a specific deployed contract.
func NewAddressCaller(address common.Address, caller bind.ContractCaller) (*AddressCaller, error) {
	contract, err := bindAddress(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AddressCaller{contract: contract}, nil
}

// NewAddressTransactor creates a new write-only instance of Address, bound to a specific deployed contract.
func NewAddressTransactor(address common.Address, transactor bind.ContractTransactor) (*AddressTransactor, error) {
	contract, err := bindAddress(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AddressTransactor{contract: contract}, nil
}

// NewAddressFilterer creates a new log filterer instance of Address, bound to a specific deployed contract.
func NewAddressFilterer(address common.Address, filterer bind.ContractFilterer) (*AddressFilterer, error) {
	contract, err := bindAddress(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AddressFilterer{contract: contract}, nil
}

// bindAddress binds a generic wrapper to an already deployed contract.
func bindAddress(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AddressMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Address *AddressRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Address.Contract.AddressCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Address *AddressRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Address.Contract.AddressTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Address *AddressRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Address.Contract.AddressTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Address *AddressCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Address.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Address *AddressTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Address.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Address *AddressTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Address.Contract.contract.Transact(opts, method, params...)
}

// ContextUpgradeableMetaData contains all meta data concerning the ContextUpgradeable contract.
var ContextUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"}]",
}

// ContextUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use ContextUpgradeableMetaData.ABI instead.
var ContextUpgradeableABI = ContextUpgradeableMetaData.ABI

// ContextUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ContextUpgradeableBinRuntime = ``

// ContextUpgradeable is an auto generated Go binding around a Kaia contract.
type ContextUpgradeable struct {
	ContextUpgradeableCaller     // Read-only binding to the contract
	ContextUpgradeableTransactor // Write-only binding to the contract
	ContextUpgradeableFilterer   // Log filterer for contract events
}

// ContextUpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type ContextUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextUpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type ContextUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextUpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ContextUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextUpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ContextUpgradeableSession struct {
	Contract     *ContextUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ContextUpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ContextUpgradeableCallerSession struct {
	Contract *ContextUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// ContextUpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ContextUpgradeableTransactorSession struct {
	Contract     *ContextUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// ContextUpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type ContextUpgradeableRaw struct {
	Contract *ContextUpgradeable // Generic contract binding to access the raw methods on
}

// ContextUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ContextUpgradeableCallerRaw struct {
	Contract *ContextUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// ContextUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ContextUpgradeableTransactorRaw struct {
	Contract *ContextUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContextUpgradeable creates a new instance of ContextUpgradeable, bound to a specific deployed contract.
func NewContextUpgradeable(address common.Address, backend bind.ContractBackend) (*ContextUpgradeable, error) {
	contract, err := bindContextUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeable{ContextUpgradeableCaller: ContextUpgradeableCaller{contract: contract}, ContextUpgradeableTransactor: ContextUpgradeableTransactor{contract: contract}, ContextUpgradeableFilterer: ContextUpgradeableFilterer{contract: contract}}, nil
}

// NewContextUpgradeableCaller creates a new read-only instance of ContextUpgradeable, bound to a specific deployed contract.
func NewContextUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*ContextUpgradeableCaller, error) {
	contract, err := bindContextUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableCaller{contract: contract}, nil
}

// NewContextUpgradeableTransactor creates a new write-only instance of ContextUpgradeable, bound to a specific deployed contract.
func NewContextUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextUpgradeableTransactor, error) {
	contract, err := bindContextUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableTransactor{contract: contract}, nil
}

// NewContextUpgradeableFilterer creates a new log filterer instance of ContextUpgradeable, bound to a specific deployed contract.
func NewContextUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextUpgradeableFilterer, error) {
	contract, err := bindContextUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableFilterer{contract: contract}, nil
}

// bindContextUpgradeable binds a generic wrapper to an already deployed contract.
func bindContextUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContextUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContextUpgradeable *ContextUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContextUpgradeable.Contract.ContextUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContextUpgradeable *ContextUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContextUpgradeable.Contract.ContextUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContextUpgradeable *ContextUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContextUpgradeable.Contract.ContextUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContextUpgradeable *ContextUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContextUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContextUpgradeable *ContextUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContextUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContextUpgradeable *ContextUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContextUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ContextUpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ContextUpgradeable contract.
type ContextUpgradeableInitializedIterator struct {
	Event *ContextUpgradeableInitialized // Event containing the contract specifics and raw log

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
func (it *ContextUpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContextUpgradeableInitialized)
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
		it.Event = new(ContextUpgradeableInitialized)
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
func (it *ContextUpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContextUpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContextUpgradeableInitialized represents a Initialized event raised by the ContextUpgradeable contract.
type ContextUpgradeableInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ContextUpgradeable *ContextUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContextUpgradeableInitializedIterator, error) {

	logs, sub, err := _ContextUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableInitializedIterator{contract: _ContextUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ContextUpgradeable *ContextUpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ContextUpgradeableInitialized) (event.Subscription, error) {

	logs, sub, err := _ContextUpgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContextUpgradeableInitialized)
				if err := _ContextUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ContextUpgradeable *ContextUpgradeableFilterer) ParseInitialized(log types.Log) (*ContextUpgradeableInitialized, error) {
	event := new(ContextUpgradeableInitialized)
	if err := _ContextUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradeableMetaData contains all meta data concerning the ERC20Upgradeable contract.
var ERC20UpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// ERC20UpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20UpgradeableMetaData.ABI instead.
var ERC20UpgradeableABI = ERC20UpgradeableMetaData.ABI

// ERC20UpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC20UpgradeableBinRuntime = ``

// Deprecated: Use ERC20UpgradeableMetaData.Sigs instead.
// ERC20UpgradeableFuncSigs maps the 4-byte function signature to its string representation.
var ERC20UpgradeableFuncSigs = ERC20UpgradeableMetaData.Sigs

// ERC20Upgradeable is an auto generated Go binding around a Kaia contract.
type ERC20Upgradeable struct {
	ERC20UpgradeableCaller     // Read-only binding to the contract
	ERC20UpgradeableTransactor // Write-only binding to the contract
	ERC20UpgradeableFilterer   // Log filterer for contract events
}

// ERC20UpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type ERC20UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20UpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type ERC20UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20UpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ERC20UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20UpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ERC20UpgradeableSession struct {
	Contract     *ERC20Upgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20UpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ERC20UpgradeableCallerSession struct {
	Contract *ERC20UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ERC20UpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ERC20UpgradeableTransactorSession struct {
	Contract     *ERC20UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ERC20UpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type ERC20UpgradeableRaw struct {
	Contract *ERC20Upgradeable // Generic contract binding to access the raw methods on
}

// ERC20UpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ERC20UpgradeableCallerRaw struct {
	Contract *ERC20UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ERC20UpgradeableTransactorRaw struct {
	Contract *ERC20UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Upgradeable creates a new instance of ERC20Upgradeable, bound to a specific deployed contract.
func NewERC20Upgradeable(address common.Address, backend bind.ContractBackend) (*ERC20Upgradeable, error) {
	contract, err := bindERC20Upgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Upgradeable{ERC20UpgradeableCaller: ERC20UpgradeableCaller{contract: contract}, ERC20UpgradeableTransactor: ERC20UpgradeableTransactor{contract: contract}, ERC20UpgradeableFilterer: ERC20UpgradeableFilterer{contract: contract}}, nil
}

// NewERC20UpgradeableCaller creates a new read-only instance of ERC20Upgradeable, bound to a specific deployed contract.
func NewERC20UpgradeableCaller(address common.Address, caller bind.ContractCaller) (*ERC20UpgradeableCaller, error) {
	contract, err := bindERC20Upgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradeableCaller{contract: contract}, nil
}

// NewERC20UpgradeableTransactor creates a new write-only instance of ERC20Upgradeable, bound to a specific deployed contract.
func NewERC20UpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20UpgradeableTransactor, error) {
	contract, err := bindERC20Upgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradeableTransactor{contract: contract}, nil
}

// NewERC20UpgradeableFilterer creates a new log filterer instance of ERC20Upgradeable, bound to a specific deployed contract.
func NewERC20UpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20UpgradeableFilterer, error) {
	contract, err := bindERC20Upgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradeableFilterer{contract: contract}, nil
}

// bindERC20Upgradeable binds a generic wrapper to an already deployed contract.
func bindERC20Upgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC20UpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Upgradeable *ERC20UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Upgradeable.Contract.ERC20UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Upgradeable *ERC20UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.ERC20UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Upgradeable *ERC20UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.ERC20UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Upgradeable *ERC20UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Upgradeable *ERC20UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Upgradeable *ERC20UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Upgradeable.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Upgradeable.Contract.Allowance(&_ERC20Upgradeable.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ERC20Upgradeable.Contract.Allowance(&_ERC20Upgradeable.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Upgradeable.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20Upgradeable.Contract.BalanceOf(&_ERC20Upgradeable.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ERC20Upgradeable.Contract.BalanceOf(&_ERC20Upgradeable.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Upgradeable *ERC20UpgradeableCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ERC20Upgradeable.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Upgradeable *ERC20UpgradeableSession) Decimals() (uint8, error) {
	return _ERC20Upgradeable.Contract.Decimals(&_ERC20Upgradeable.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC20Upgradeable *ERC20UpgradeableCallerSession) Decimals() (uint8, error) {
	return _ERC20Upgradeable.Contract.Decimals(&_ERC20Upgradeable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Upgradeable *ERC20UpgradeableCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC20Upgradeable.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Upgradeable *ERC20UpgradeableSession) Name() (string, error) {
	return _ERC20Upgradeable.Contract.Name(&_ERC20Upgradeable.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ERC20Upgradeable *ERC20UpgradeableCallerSession) Name() (string, error) {
	return _ERC20Upgradeable.Contract.Name(&_ERC20Upgradeable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Upgradeable *ERC20UpgradeableCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ERC20Upgradeable.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Upgradeable *ERC20UpgradeableSession) Symbol() (string, error) {
	return _ERC20Upgradeable.Contract.Symbol(&_ERC20Upgradeable.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ERC20Upgradeable *ERC20UpgradeableCallerSession) Symbol() (string, error) {
	return _ERC20Upgradeable.Contract.Symbol(&_ERC20Upgradeable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Upgradeable.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableSession) TotalSupply() (*big.Int, error) {
	return _ERC20Upgradeable.Contract.TotalSupply(&_ERC20Upgradeable.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ERC20Upgradeable *ERC20UpgradeableCallerSession) TotalSupply() (*big.Int, error) {
	return _ERC20Upgradeable.Contract.TotalSupply(&_ERC20Upgradeable.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.Approve(&_ERC20Upgradeable.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.Approve(&_ERC20Upgradeable.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.Transfer(&_ERC20Upgradeable.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.Transfer(&_ERC20Upgradeable.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.TransferFrom(&_ERC20Upgradeable.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_ERC20Upgradeable *ERC20UpgradeableTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _ERC20Upgradeable.Contract.TransferFrom(&_ERC20Upgradeable.TransactOpts, from, to, value)
}

// ERC20UpgradeableApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ERC20Upgradeable contract.
type ERC20UpgradeableApprovalIterator struct {
	Event *ERC20UpgradeableApproval // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradeableApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradeableApproval)
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
		it.Event = new(ERC20UpgradeableApproval)
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
func (it *ERC20UpgradeableApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradeableApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradeableApproval represents a Approval event raised by the ERC20Upgradeable contract.
type ERC20UpgradeableApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ERC20UpgradeableApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Upgradeable.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradeableApprovalIterator{contract: _ERC20Upgradeable.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ERC20UpgradeableApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ERC20Upgradeable.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradeableApproval)
				if err := _ERC20Upgradeable.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) ParseApproval(log types.Log) (*ERC20UpgradeableApproval, error) {
	event := new(ERC20UpgradeableApproval)
	if err := _ERC20Upgradeable.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ERC20Upgradeable contract.
type ERC20UpgradeableInitializedIterator struct {
	Event *ERC20UpgradeableInitialized // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradeableInitialized)
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
		it.Event = new(ERC20UpgradeableInitialized)
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
func (it *ERC20UpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradeableInitialized represents a Initialized event raised by the ERC20Upgradeable contract.
type ERC20UpgradeableInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*ERC20UpgradeableInitializedIterator, error) {

	logs, sub, err := _ERC20Upgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradeableInitializedIterator{contract: _ERC20Upgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ERC20UpgradeableInitialized) (event.Subscription, error) {

	logs, sub, err := _ERC20Upgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradeableInitialized)
				if err := _ERC20Upgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) ParseInitialized(log types.Log) (*ERC20UpgradeableInitialized, error) {
	event := new(ERC20UpgradeableInitialized)
	if err := _ERC20Upgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC20UpgradeableTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ERC20Upgradeable contract.
type ERC20UpgradeableTransferIterator struct {
	Event *ERC20UpgradeableTransfer // Event containing the contract specifics and raw log

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
func (it *ERC20UpgradeableTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC20UpgradeableTransfer)
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
		it.Event = new(ERC20UpgradeableTransfer)
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
func (it *ERC20UpgradeableTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC20UpgradeableTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC20UpgradeableTransfer represents a Transfer event raised by the ERC20Upgradeable contract.
type ERC20UpgradeableTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ERC20UpgradeableTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Upgradeable.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ERC20UpgradeableTransferIterator{contract: _ERC20Upgradeable.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ERC20UpgradeableTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ERC20Upgradeable.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC20UpgradeableTransfer)
				if err := _ERC20Upgradeable.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ERC20Upgradeable *ERC20UpgradeableFilterer) ParseTransfer(log types.Log) (*ERC20UpgradeableTransfer, error) {
	event := new(ERC20UpgradeableTransfer)
	if err := _ERC20Upgradeable.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ErrorsMetaData contains all meta data concerning the Errors contract.
var ErrorsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedDeployment\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"MissingPrecompile\",\"type\":\"error\"}]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea264697066735822122045bd5ccaa497192dc3cf7cf9c42324c6699541b8e1c440fb85f637fcb9bb93ba64736f6c63430008190033",
}

// ErrorsABI is the input ABI used to generate the binding from.
// Deprecated: Use ErrorsMetaData.ABI instead.
var ErrorsABI = ErrorsMetaData.ABI

// ErrorsBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ErrorsBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea264697066735822122045bd5ccaa497192dc3cf7cf9c42324c6699541b8e1c440fb85f637fcb9bb93ba64736f6c63430008190033`

// ErrorsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ErrorsMetaData.Bin instead.
var ErrorsBin = ErrorsMetaData.Bin

// DeployErrors deploys a new Kaia contract, binding an instance of Errors to it.
func DeployErrors(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Errors, error) {
	parsed, err := ErrorsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ErrorsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Errors{ErrorsCaller: ErrorsCaller{contract: contract}, ErrorsTransactor: ErrorsTransactor{contract: contract}, ErrorsFilterer: ErrorsFilterer{contract: contract}}, nil
}

// Errors is an auto generated Go binding around a Kaia contract.
type Errors struct {
	ErrorsCaller     // Read-only binding to the contract
	ErrorsTransactor // Write-only binding to the contract
	ErrorsFilterer   // Log filterer for contract events
}

// ErrorsCaller is an auto generated read-only Go binding around a Kaia contract.
type ErrorsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErrorsTransactor is an auto generated write-only Go binding around a Kaia contract.
type ErrorsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErrorsFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ErrorsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ErrorsSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ErrorsSession struct {
	Contract     *Errors           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ErrorsCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ErrorsCallerSession struct {
	Contract *ErrorsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ErrorsTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ErrorsTransactorSession struct {
	Contract     *ErrorsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ErrorsRaw is an auto generated low-level Go binding around a Kaia contract.
type ErrorsRaw struct {
	Contract *Errors // Generic contract binding to access the raw methods on
}

// ErrorsCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ErrorsCallerRaw struct {
	Contract *ErrorsCaller // Generic read-only contract binding to access the raw methods on
}

// ErrorsTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ErrorsTransactorRaw struct {
	Contract *ErrorsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErrors creates a new instance of Errors, bound to a specific deployed contract.
func NewErrors(address common.Address, backend bind.ContractBackend) (*Errors, error) {
	contract, err := bindErrors(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Errors{ErrorsCaller: ErrorsCaller{contract: contract}, ErrorsTransactor: ErrorsTransactor{contract: contract}, ErrorsFilterer: ErrorsFilterer{contract: contract}}, nil
}

// NewErrorsCaller creates a new read-only instance of Errors, bound to a specific deployed contract.
func NewErrorsCaller(address common.Address, caller bind.ContractCaller) (*ErrorsCaller, error) {
	contract, err := bindErrors(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ErrorsCaller{contract: contract}, nil
}

// NewErrorsTransactor creates a new write-only instance of Errors, bound to a specific deployed contract.
func NewErrorsTransactor(address common.Address, transactor bind.ContractTransactor) (*ErrorsTransactor, error) {
	contract, err := bindErrors(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ErrorsTransactor{contract: contract}, nil
}

// NewErrorsFilterer creates a new log filterer instance of Errors, bound to a specific deployed contract.
func NewErrorsFilterer(address common.Address, filterer bind.ContractFilterer) (*ErrorsFilterer, error) {
	contract, err := bindErrors(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ErrorsFilterer{contract: contract}, nil
}

// bindErrors binds a generic wrapper to an already deployed contract.
func bindErrors(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ErrorsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Errors *ErrorsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Errors.Contract.ErrorsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Errors *ErrorsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Errors.Contract.ErrorsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Errors *ErrorsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Errors.Contract.ErrorsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Errors *ErrorsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Errors.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Errors *ErrorsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Errors.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Errors *ErrorsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Errors.Contract.contract.Transact(opts, method, params...)
}

// ICnStakingMetaData contains all meta data concerning the ICnStaking contract.
var ICnStakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"BaseCnStakingMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidStakeFor\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidWithdrawalState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotStaker\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotUnstakingManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotWithdrawableYet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RedelegationCooldown\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RedelegationDisabled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroValue\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"approvedWithdrawalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"withdrawableFrom\",\"type\":\"uint256\"}],\"name\":\"ApproveStakingWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"approvedWithdrawalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"CancelApprovedStakingWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"DelegateKaia\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"contractType\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"publicDelegation\",\"type\":\"address\"}],\"name\":\"DeployCnStakingV4\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"prevCnStaking\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"targetCnStaking\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"HandleRedelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"targetCnStaking\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Redelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"approvedWithdrawalId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"WithdrawApprovedStaking\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"CONTRACT_TYPE\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"STAKE_LOCKUP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approveStakingWithdrawal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"cancelApprovedStakingWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"delegate\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_from\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_to\",\"type\":\"uint256\"},{\"internalType\":\"enumICnStaking.WithdrawalStakingState\",\"name\":\"_state\",\"type\":\"uint8\"}],\"name\":\"getApprovedStakingWithdrawalIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getApprovedStakingWithdrawalInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"withdrawableFrom\",\"type\":\"uint256\"},{\"internalType\":\"enumICnStaking.WithdrawalStakingState\",\"name\":\"state\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"handleRedelegation\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_publicDelegation\",\"type\":\"address\"}],\"name\":\"initializeWithPD\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"}],\"name\":\"lastRedelegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"publicDelegation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_targetCnStaking\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"redelegate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"staking\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstaking\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"withdrawApprovedStaking\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawalRequestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"4b6a94cc": "CONTRACT_TYPE()",
		"96106ae4": "STAKE_LOCKUP()",
		"ffa1ad74": "VERSION()",
		"5df8b09a": "approveStakingWithdrawal(address,uint256)",
		"c804b115": "cancelApprovedStakingWithdrawal(uint256)",
		"c89e4361": "delegate()",
		"d2569eb9": "getApprovedStakingWithdrawalIds(uint256,uint256,uint8)",
		"725c0503": "getApprovedStakingWithdrawalInfo(uint256)",
		"a006e90c": "handleRedelegation(address)",
		"c4d66de8": "initialize(address)",
		"3e8f5335": "initializeWithPD(address,address)",
		"14d3ce10": "lastRedelegation(address)",
		"e1a12d35": "publicDelegation()",
		"6bd8f804": "redelegate(address,address,uint256)",
		"4cf088d9": "staking()",
		"630b1146": "unstaking()",
		"6e93df0d": "withdrawApprovedStaking(uint256)",
		"19e44e32": "withdrawalRequestCount()",
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

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() pure returns(string)
func (_ICnStaking *ICnStakingCaller) CONTRACTTYPE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "CONTRACT_TYPE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() pure returns(string)
func (_ICnStaking *ICnStakingSession) CONTRACTTYPE() (string, error) {
	return _ICnStaking.Contract.CONTRACTTYPE(&_ICnStaking.CallOpts)
}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() pure returns(string)
func (_ICnStaking *ICnStakingCallerSession) CONTRACTTYPE() (string, error) {
	return _ICnStaking.Contract.CONTRACTTYPE(&_ICnStaking.CallOpts)
}

// STAKELOCKUP is a free data retrieval call binding the contract method 0x96106ae4.
//
// Solidity: function STAKE_LOCKUP() pure returns(uint256)
func (_ICnStaking *ICnStakingCaller) STAKELOCKUP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "STAKE_LOCKUP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// STAKELOCKUP is a free data retrieval call binding the contract method 0x96106ae4.
//
// Solidity: function STAKE_LOCKUP() pure returns(uint256)
func (_ICnStaking *ICnStakingSession) STAKELOCKUP() (*big.Int, error) {
	return _ICnStaking.Contract.STAKELOCKUP(&_ICnStaking.CallOpts)
}

// STAKELOCKUP is a free data retrieval call binding the contract method 0x96106ae4.
//
// Solidity: function STAKE_LOCKUP() pure returns(uint256)
func (_ICnStaking *ICnStakingCallerSession) STAKELOCKUP() (*big.Int, error) {
	return _ICnStaking.Contract.STAKELOCKUP(&_ICnStaking.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() pure returns(uint256)
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
// Solidity: function VERSION() pure returns(uint256)
func (_ICnStaking *ICnStakingSession) VERSION() (*big.Int, error) {
	return _ICnStaking.Contract.VERSION(&_ICnStaking.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() pure returns(uint256)
func (_ICnStaking *ICnStakingCallerSession) VERSION() (*big.Int, error) {
	return _ICnStaking.Contract.VERSION(&_ICnStaking.CallOpts)
}

// GetApprovedStakingWithdrawalIds is a free data retrieval call binding the contract method 0xd2569eb9.
//
// Solidity: function getApprovedStakingWithdrawalIds(uint256 _from, uint256 _to, uint8 _state) view returns(uint256[] ids)
func (_ICnStaking *ICnStakingCaller) GetApprovedStakingWithdrawalIds(opts *bind.CallOpts, _from *big.Int, _to *big.Int, _state uint8) ([]*big.Int, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "getApprovedStakingWithdrawalIds", _from, _to, _state)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetApprovedStakingWithdrawalIds is a free data retrieval call binding the contract method 0xd2569eb9.
//
// Solidity: function getApprovedStakingWithdrawalIds(uint256 _from, uint256 _to, uint8 _state) view returns(uint256[] ids)
func (_ICnStaking *ICnStakingSession) GetApprovedStakingWithdrawalIds(_from *big.Int, _to *big.Int, _state uint8) ([]*big.Int, error) {
	return _ICnStaking.Contract.GetApprovedStakingWithdrawalIds(&_ICnStaking.CallOpts, _from, _to, _state)
}

// GetApprovedStakingWithdrawalIds is a free data retrieval call binding the contract method 0xd2569eb9.
//
// Solidity: function getApprovedStakingWithdrawalIds(uint256 _from, uint256 _to, uint8 _state) view returns(uint256[] ids)
func (_ICnStaking *ICnStakingCallerSession) GetApprovedStakingWithdrawalIds(_from *big.Int, _to *big.Int, _state uint8) ([]*big.Int, error) {
	return _ICnStaking.Contract.GetApprovedStakingWithdrawalIds(&_ICnStaking.CallOpts, _from, _to, _state)
}

// GetApprovedStakingWithdrawalInfo is a free data retrieval call binding the contract method 0x725c0503.
//
// Solidity: function getApprovedStakingWithdrawalInfo(uint256 _index) view returns(address to, uint256 value, uint256 withdrawableFrom, uint8 state)
func (_ICnStaking *ICnStakingCaller) GetApprovedStakingWithdrawalInfo(opts *bind.CallOpts, _index *big.Int) (struct {
	To               common.Address
	Value            *big.Int
	WithdrawableFrom *big.Int
	State            uint8
}, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "getApprovedStakingWithdrawalInfo", _index)

	outstruct := new(struct {
		To               common.Address
		Value            *big.Int
		WithdrawableFrom *big.Int
		State            uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.To = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Value = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.WithdrawableFrom = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.State = *abi.ConvertType(out[3], new(uint8)).(*uint8)

	return *outstruct, err

}

// GetApprovedStakingWithdrawalInfo is a free data retrieval call binding the contract method 0x725c0503.
//
// Solidity: function getApprovedStakingWithdrawalInfo(uint256 _index) view returns(address to, uint256 value, uint256 withdrawableFrom, uint8 state)
func (_ICnStaking *ICnStakingSession) GetApprovedStakingWithdrawalInfo(_index *big.Int) (struct {
	To               common.Address
	Value            *big.Int
	WithdrawableFrom *big.Int
	State            uint8
}, error) {
	return _ICnStaking.Contract.GetApprovedStakingWithdrawalInfo(&_ICnStaking.CallOpts, _index)
}

// GetApprovedStakingWithdrawalInfo is a free data retrieval call binding the contract method 0x725c0503.
//
// Solidity: function getApprovedStakingWithdrawalInfo(uint256 _index) view returns(address to, uint256 value, uint256 withdrawableFrom, uint8 state)
func (_ICnStaking *ICnStakingCallerSession) GetApprovedStakingWithdrawalInfo(_index *big.Int) (struct {
	To               common.Address
	Value            *big.Int
	WithdrawableFrom *big.Int
	State            uint8
}, error) {
	return _ICnStaking.Contract.GetApprovedStakingWithdrawalInfo(&_ICnStaking.CallOpts, _index)
}

// LastRedelegation is a free data retrieval call binding the contract method 0x14d3ce10.
//
// Solidity: function lastRedelegation(address _account) view returns(uint256)
func (_ICnStaking *ICnStakingCaller) LastRedelegation(opts *bind.CallOpts, _account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "lastRedelegation", _account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastRedelegation is a free data retrieval call binding the contract method 0x14d3ce10.
//
// Solidity: function lastRedelegation(address _account) view returns(uint256)
func (_ICnStaking *ICnStakingSession) LastRedelegation(_account common.Address) (*big.Int, error) {
	return _ICnStaking.Contract.LastRedelegation(&_ICnStaking.CallOpts, _account)
}

// LastRedelegation is a free data retrieval call binding the contract method 0x14d3ce10.
//
// Solidity: function lastRedelegation(address _account) view returns(uint256)
func (_ICnStaking *ICnStakingCallerSession) LastRedelegation(_account common.Address) (*big.Int, error) {
	return _ICnStaking.Contract.LastRedelegation(&_ICnStaking.CallOpts, _account)
}

// PublicDelegation is a free data retrieval call binding the contract method 0xe1a12d35.
//
// Solidity: function publicDelegation() view returns(address)
func (_ICnStaking *ICnStakingCaller) PublicDelegation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "publicDelegation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PublicDelegation is a free data retrieval call binding the contract method 0xe1a12d35.
//
// Solidity: function publicDelegation() view returns(address)
func (_ICnStaking *ICnStakingSession) PublicDelegation() (common.Address, error) {
	return _ICnStaking.Contract.PublicDelegation(&_ICnStaking.CallOpts)
}

// PublicDelegation is a free data retrieval call binding the contract method 0xe1a12d35.
//
// Solidity: function publicDelegation() view returns(address)
func (_ICnStaking *ICnStakingCallerSession) PublicDelegation() (common.Address, error) {
	return _ICnStaking.Contract.PublicDelegation(&_ICnStaking.CallOpts)
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

// WithdrawalRequestCount is a free data retrieval call binding the contract method 0x19e44e32.
//
// Solidity: function withdrawalRequestCount() view returns(uint256)
func (_ICnStaking *ICnStakingCaller) WithdrawalRequestCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICnStaking.contract.Call(opts, &out, "withdrawalRequestCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalRequestCount is a free data retrieval call binding the contract method 0x19e44e32.
//
// Solidity: function withdrawalRequestCount() view returns(uint256)
func (_ICnStaking *ICnStakingSession) WithdrawalRequestCount() (*big.Int, error) {
	return _ICnStaking.Contract.WithdrawalRequestCount(&_ICnStaking.CallOpts)
}

// WithdrawalRequestCount is a free data retrieval call binding the contract method 0x19e44e32.
//
// Solidity: function withdrawalRequestCount() view returns(uint256)
func (_ICnStaking *ICnStakingCallerSession) WithdrawalRequestCount() (*big.Int, error) {
	return _ICnStaking.Contract.WithdrawalRequestCount(&_ICnStaking.CallOpts)
}

// ApproveStakingWithdrawal is a paid mutator transaction binding the contract method 0x5df8b09a.
//
// Solidity: function approveStakingWithdrawal(address _to, uint256 _value) returns(uint256 id)
func (_ICnStaking *ICnStakingTransactor) ApproveStakingWithdrawal(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ICnStaking.contract.Transact(opts, "approveStakingWithdrawal", _to, _value)
}

// ApproveStakingWithdrawal is a paid mutator transaction binding the contract method 0x5df8b09a.
//
// Solidity: function approveStakingWithdrawal(address _to, uint256 _value) returns(uint256 id)
func (_ICnStaking *ICnStakingSession) ApproveStakingWithdrawal(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ICnStaking.Contract.ApproveStakingWithdrawal(&_ICnStaking.TransactOpts, _to, _value)
}

// ApproveStakingWithdrawal is a paid mutator transaction binding the contract method 0x5df8b09a.
//
// Solidity: function approveStakingWithdrawal(address _to, uint256 _value) returns(uint256 id)
func (_ICnStaking *ICnStakingTransactorSession) ApproveStakingWithdrawal(_to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ICnStaking.Contract.ApproveStakingWithdrawal(&_ICnStaking.TransactOpts, _to, _value)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _id) returns()
func (_ICnStaking *ICnStakingTransactor) CancelApprovedStakingWithdrawal(opts *bind.TransactOpts, _id *big.Int) (*types.Transaction, error) {
	return _ICnStaking.contract.Transact(opts, "cancelApprovedStakingWithdrawal", _id)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _id) returns()
func (_ICnStaking *ICnStakingSession) CancelApprovedStakingWithdrawal(_id *big.Int) (*types.Transaction, error) {
	return _ICnStaking.Contract.CancelApprovedStakingWithdrawal(&_ICnStaking.TransactOpts, _id)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _id) returns()
func (_ICnStaking *ICnStakingTransactorSession) CancelApprovedStakingWithdrawal(_id *big.Int) (*types.Transaction, error) {
	return _ICnStaking.Contract.CancelApprovedStakingWithdrawal(&_ICnStaking.TransactOpts, _id)
}

// Delegate is a paid mutator transaction binding the contract method 0xc89e4361.
//
// Solidity: function delegate() payable returns()
func (_ICnStaking *ICnStakingTransactor) Delegate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICnStaking.contract.Transact(opts, "delegate")
}

// Delegate is a paid mutator transaction binding the contract method 0xc89e4361.
//
// Solidity: function delegate() payable returns()
func (_ICnStaking *ICnStakingSession) Delegate() (*types.Transaction, error) {
	return _ICnStaking.Contract.Delegate(&_ICnStaking.TransactOpts)
}

// Delegate is a paid mutator transaction binding the contract method 0xc89e4361.
//
// Solidity: function delegate() payable returns()
func (_ICnStaking *ICnStakingTransactorSession) Delegate() (*types.Transaction, error) {
	return _ICnStaking.Contract.Delegate(&_ICnStaking.TransactOpts)
}

// HandleRedelegation is a paid mutator transaction binding the contract method 0xa006e90c.
//
// Solidity: function handleRedelegation(address _user) payable returns()
func (_ICnStaking *ICnStakingTransactor) HandleRedelegation(opts *bind.TransactOpts, _user common.Address) (*types.Transaction, error) {
	return _ICnStaking.contract.Transact(opts, "handleRedelegation", _user)
}

// HandleRedelegation is a paid mutator transaction binding the contract method 0xa006e90c.
//
// Solidity: function handleRedelegation(address _user) payable returns()
func (_ICnStaking *ICnStakingSession) HandleRedelegation(_user common.Address) (*types.Transaction, error) {
	return _ICnStaking.Contract.HandleRedelegation(&_ICnStaking.TransactOpts, _user)
}

// HandleRedelegation is a paid mutator transaction binding the contract method 0xa006e90c.
//
// Solidity: function handleRedelegation(address _user) payable returns()
func (_ICnStaking *ICnStakingTransactorSession) HandleRedelegation(_user common.Address) (*types.Transaction, error) {
	return _ICnStaking.Contract.HandleRedelegation(&_ICnStaking.TransactOpts, _user)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_ICnStaking *ICnStakingTransactor) Initialize(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _ICnStaking.contract.Transact(opts, "initialize", _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_ICnStaking *ICnStakingSession) Initialize(_owner common.Address) (*types.Transaction, error) {
	return _ICnStaking.Contract.Initialize(&_ICnStaking.TransactOpts, _owner)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address _owner) returns()
func (_ICnStaking *ICnStakingTransactorSession) Initialize(_owner common.Address) (*types.Transaction, error) {
	return _ICnStaking.Contract.Initialize(&_ICnStaking.TransactOpts, _owner)
}

// InitializeWithPD is a paid mutator transaction binding the contract method 0x3e8f5335.
//
// Solidity: function initializeWithPD(address _owner, address _publicDelegation) returns()
func (_ICnStaking *ICnStakingTransactor) InitializeWithPD(opts *bind.TransactOpts, _owner common.Address, _publicDelegation common.Address) (*types.Transaction, error) {
	return _ICnStaking.contract.Transact(opts, "initializeWithPD", _owner, _publicDelegation)
}

// InitializeWithPD is a paid mutator transaction binding the contract method 0x3e8f5335.
//
// Solidity: function initializeWithPD(address _owner, address _publicDelegation) returns()
func (_ICnStaking *ICnStakingSession) InitializeWithPD(_owner common.Address, _publicDelegation common.Address) (*types.Transaction, error) {
	return _ICnStaking.Contract.InitializeWithPD(&_ICnStaking.TransactOpts, _owner, _publicDelegation)
}

// InitializeWithPD is a paid mutator transaction binding the contract method 0x3e8f5335.
//
// Solidity: function initializeWithPD(address _owner, address _publicDelegation) returns()
func (_ICnStaking *ICnStakingTransactorSession) InitializeWithPD(_owner common.Address, _publicDelegation common.Address) (*types.Transaction, error) {
	return _ICnStaking.Contract.InitializeWithPD(&_ICnStaking.TransactOpts, _owner, _publicDelegation)
}

// Redelegate is a paid mutator transaction binding the contract method 0x6bd8f804.
//
// Solidity: function redelegate(address _user, address _targetCnStaking, uint256 _value) returns()
func (_ICnStaking *ICnStakingTransactor) Redelegate(opts *bind.TransactOpts, _user common.Address, _targetCnStaking common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ICnStaking.contract.Transact(opts, "redelegate", _user, _targetCnStaking, _value)
}

// Redelegate is a paid mutator transaction binding the contract method 0x6bd8f804.
//
// Solidity: function redelegate(address _user, address _targetCnStaking, uint256 _value) returns()
func (_ICnStaking *ICnStakingSession) Redelegate(_user common.Address, _targetCnStaking common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ICnStaking.Contract.Redelegate(&_ICnStaking.TransactOpts, _user, _targetCnStaking, _value)
}

// Redelegate is a paid mutator transaction binding the contract method 0x6bd8f804.
//
// Solidity: function redelegate(address _user, address _targetCnStaking, uint256 _value) returns()
func (_ICnStaking *ICnStakingTransactorSession) Redelegate(_user common.Address, _targetCnStaking common.Address, _value *big.Int) (*types.Transaction, error) {
	return _ICnStaking.Contract.Redelegate(&_ICnStaking.TransactOpts, _user, _targetCnStaking, _value)
}

// WithdrawApprovedStaking is a paid mutator transaction binding the contract method 0x6e93df0d.
//
// Solidity: function withdrawApprovedStaking(uint256 _id) returns()
func (_ICnStaking *ICnStakingTransactor) WithdrawApprovedStaking(opts *bind.TransactOpts, _id *big.Int) (*types.Transaction, error) {
	return _ICnStaking.contract.Transact(opts, "withdrawApprovedStaking", _id)
}

// WithdrawApprovedStaking is a paid mutator transaction binding the contract method 0x6e93df0d.
//
// Solidity: function withdrawApprovedStaking(uint256 _id) returns()
func (_ICnStaking *ICnStakingSession) WithdrawApprovedStaking(_id *big.Int) (*types.Transaction, error) {
	return _ICnStaking.Contract.WithdrawApprovedStaking(&_ICnStaking.TransactOpts, _id)
}

// WithdrawApprovedStaking is a paid mutator transaction binding the contract method 0x6e93df0d.
//
// Solidity: function withdrawApprovedStaking(uint256 _id) returns()
func (_ICnStaking *ICnStakingTransactorSession) WithdrawApprovedStaking(_id *big.Int) (*types.Transaction, error) {
	return _ICnStaking.Contract.WithdrawApprovedStaking(&_ICnStaking.TransactOpts, _id)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ICnStaking *ICnStakingTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICnStaking.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ICnStaking *ICnStakingSession) Receive() (*types.Transaction, error) {
	return _ICnStaking.Contract.Receive(&_ICnStaking.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ICnStaking *ICnStakingTransactorSession) Receive() (*types.Transaction, error) {
	return _ICnStaking.Contract.Receive(&_ICnStaking.TransactOpts)
}

// ICnStakingApproveStakingWithdrawalIterator is returned from FilterApproveStakingWithdrawal and is used to iterate over the raw logs and unpacked data for ApproveStakingWithdrawal events raised by the ICnStaking contract.
type ICnStakingApproveStakingWithdrawalIterator struct {
	Event *ICnStakingApproveStakingWithdrawal // Event containing the contract specifics and raw log

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
func (it *ICnStakingApproveStakingWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingApproveStakingWithdrawal)
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
		it.Event = new(ICnStakingApproveStakingWithdrawal)
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
func (it *ICnStakingApproveStakingWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingApproveStakingWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingApproveStakingWithdrawal represents a ApproveStakingWithdrawal event raised by the ICnStaking contract.
type ICnStakingApproveStakingWithdrawal struct {
	ApprovedWithdrawalId *big.Int
	To                   common.Address
	Value                *big.Int
	WithdrawableFrom     *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterApproveStakingWithdrawal is a free log retrieval operation binding the contract event 0xdd0988b8c11a867814c87be652f93a86d35c6a235b92977e99d5394bd8580ced.
//
// Solidity: event ApproveStakingWithdrawal(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value, uint256 withdrawableFrom)
func (_ICnStaking *ICnStakingFilterer) FilterApproveStakingWithdrawal(opts *bind.FilterOpts, approvedWithdrawalId []*big.Int, to []common.Address) (*ICnStakingApproveStakingWithdrawalIterator, error) {

	var approvedWithdrawalIdRule []interface{}
	for _, approvedWithdrawalIdItem := range approvedWithdrawalId {
		approvedWithdrawalIdRule = append(approvedWithdrawalIdRule, approvedWithdrawalIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ICnStaking.contract.FilterLogs(opts, "ApproveStakingWithdrawal", approvedWithdrawalIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingApproveStakingWithdrawalIterator{contract: _ICnStaking.contract, event: "ApproveStakingWithdrawal", logs: logs, sub: sub}, nil
}

// WatchApproveStakingWithdrawal is a free log subscription operation binding the contract event 0xdd0988b8c11a867814c87be652f93a86d35c6a235b92977e99d5394bd8580ced.
//
// Solidity: event ApproveStakingWithdrawal(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value, uint256 withdrawableFrom)
func (_ICnStaking *ICnStakingFilterer) WatchApproveStakingWithdrawal(opts *bind.WatchOpts, sink chan<- *ICnStakingApproveStakingWithdrawal, approvedWithdrawalId []*big.Int, to []common.Address) (event.Subscription, error) {

	var approvedWithdrawalIdRule []interface{}
	for _, approvedWithdrawalIdItem := range approvedWithdrawalId {
		approvedWithdrawalIdRule = append(approvedWithdrawalIdRule, approvedWithdrawalIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ICnStaking.contract.WatchLogs(opts, "ApproveStakingWithdrawal", approvedWithdrawalIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingApproveStakingWithdrawal)
				if err := _ICnStaking.contract.UnpackLog(event, "ApproveStakingWithdrawal", log); err != nil {
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

// ParseApproveStakingWithdrawal is a log parse operation binding the contract event 0xdd0988b8c11a867814c87be652f93a86d35c6a235b92977e99d5394bd8580ced.
//
// Solidity: event ApproveStakingWithdrawal(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value, uint256 withdrawableFrom)
func (_ICnStaking *ICnStakingFilterer) ParseApproveStakingWithdrawal(log types.Log) (*ICnStakingApproveStakingWithdrawal, error) {
	event := new(ICnStakingApproveStakingWithdrawal)
	if err := _ICnStaking.contract.UnpackLog(event, "ApproveStakingWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICnStakingCancelApprovedStakingWithdrawalIterator is returned from FilterCancelApprovedStakingWithdrawal and is used to iterate over the raw logs and unpacked data for CancelApprovedStakingWithdrawal events raised by the ICnStaking contract.
type ICnStakingCancelApprovedStakingWithdrawalIterator struct {
	Event *ICnStakingCancelApprovedStakingWithdrawal // Event containing the contract specifics and raw log

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
func (it *ICnStakingCancelApprovedStakingWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingCancelApprovedStakingWithdrawal)
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
		it.Event = new(ICnStakingCancelApprovedStakingWithdrawal)
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
func (it *ICnStakingCancelApprovedStakingWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingCancelApprovedStakingWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingCancelApprovedStakingWithdrawal represents a CancelApprovedStakingWithdrawal event raised by the ICnStaking contract.
type ICnStakingCancelApprovedStakingWithdrawal struct {
	ApprovedWithdrawalId *big.Int
	To                   common.Address
	Value                *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterCancelApprovedStakingWithdrawal is a free log retrieval operation binding the contract event 0xcc847b5f283b573ff21408ad42cb442f358bbce95269470b05097215598173df.
//
// Solidity: event CancelApprovedStakingWithdrawal(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value)
func (_ICnStaking *ICnStakingFilterer) FilterCancelApprovedStakingWithdrawal(opts *bind.FilterOpts, approvedWithdrawalId []*big.Int, to []common.Address) (*ICnStakingCancelApprovedStakingWithdrawalIterator, error) {

	var approvedWithdrawalIdRule []interface{}
	for _, approvedWithdrawalIdItem := range approvedWithdrawalId {
		approvedWithdrawalIdRule = append(approvedWithdrawalIdRule, approvedWithdrawalIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ICnStaking.contract.FilterLogs(opts, "CancelApprovedStakingWithdrawal", approvedWithdrawalIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingCancelApprovedStakingWithdrawalIterator{contract: _ICnStaking.contract, event: "CancelApprovedStakingWithdrawal", logs: logs, sub: sub}, nil
}

// WatchCancelApprovedStakingWithdrawal is a free log subscription operation binding the contract event 0xcc847b5f283b573ff21408ad42cb442f358bbce95269470b05097215598173df.
//
// Solidity: event CancelApprovedStakingWithdrawal(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value)
func (_ICnStaking *ICnStakingFilterer) WatchCancelApprovedStakingWithdrawal(opts *bind.WatchOpts, sink chan<- *ICnStakingCancelApprovedStakingWithdrawal, approvedWithdrawalId []*big.Int, to []common.Address) (event.Subscription, error) {

	var approvedWithdrawalIdRule []interface{}
	for _, approvedWithdrawalIdItem := range approvedWithdrawalId {
		approvedWithdrawalIdRule = append(approvedWithdrawalIdRule, approvedWithdrawalIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ICnStaking.contract.WatchLogs(opts, "CancelApprovedStakingWithdrawal", approvedWithdrawalIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingCancelApprovedStakingWithdrawal)
				if err := _ICnStaking.contract.UnpackLog(event, "CancelApprovedStakingWithdrawal", log); err != nil {
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

// ParseCancelApprovedStakingWithdrawal is a log parse operation binding the contract event 0xcc847b5f283b573ff21408ad42cb442f358bbce95269470b05097215598173df.
//
// Solidity: event CancelApprovedStakingWithdrawal(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value)
func (_ICnStaking *ICnStakingFilterer) ParseCancelApprovedStakingWithdrawal(log types.Log) (*ICnStakingCancelApprovedStakingWithdrawal, error) {
	event := new(ICnStakingCancelApprovedStakingWithdrawal)
	if err := _ICnStaking.contract.UnpackLog(event, "CancelApprovedStakingWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICnStakingDelegateKaiaIterator is returned from FilterDelegateKaia and is used to iterate over the raw logs and unpacked data for DelegateKaia events raised by the ICnStaking contract.
type ICnStakingDelegateKaiaIterator struct {
	Event *ICnStakingDelegateKaia // Event containing the contract specifics and raw log

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
func (it *ICnStakingDelegateKaiaIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingDelegateKaia)
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
		it.Event = new(ICnStakingDelegateKaia)
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
func (it *ICnStakingDelegateKaiaIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingDelegateKaiaIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingDelegateKaia represents a DelegateKaia event raised by the ICnStaking contract.
type ICnStakingDelegateKaia struct {
	From  common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDelegateKaia is a free log retrieval operation binding the contract event 0x8ecbcbd560048c1f474847748189e8c43f1fdee737fbd34523b04b5a9bbdffaa.
//
// Solidity: event DelegateKaia(address indexed from, uint256 value)
func (_ICnStaking *ICnStakingFilterer) FilterDelegateKaia(opts *bind.FilterOpts, from []common.Address) (*ICnStakingDelegateKaiaIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _ICnStaking.contract.FilterLogs(opts, "DelegateKaia", fromRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingDelegateKaiaIterator{contract: _ICnStaking.contract, event: "DelegateKaia", logs: logs, sub: sub}, nil
}

// WatchDelegateKaia is a free log subscription operation binding the contract event 0x8ecbcbd560048c1f474847748189e8c43f1fdee737fbd34523b04b5a9bbdffaa.
//
// Solidity: event DelegateKaia(address indexed from, uint256 value)
func (_ICnStaking *ICnStakingFilterer) WatchDelegateKaia(opts *bind.WatchOpts, sink chan<- *ICnStakingDelegateKaia, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _ICnStaking.contract.WatchLogs(opts, "DelegateKaia", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingDelegateKaia)
				if err := _ICnStaking.contract.UnpackLog(event, "DelegateKaia", log); err != nil {
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

// ParseDelegateKaia is a log parse operation binding the contract event 0x8ecbcbd560048c1f474847748189e8c43f1fdee737fbd34523b04b5a9bbdffaa.
//
// Solidity: event DelegateKaia(address indexed from, uint256 value)
func (_ICnStaking *ICnStakingFilterer) ParseDelegateKaia(log types.Log) (*ICnStakingDelegateKaia, error) {
	event := new(ICnStakingDelegateKaia)
	if err := _ICnStaking.contract.UnpackLog(event, "DelegateKaia", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICnStakingDeployCnStakingV4Iterator is returned from FilterDeployCnStakingV4 and is used to iterate over the raw logs and unpacked data for DeployCnStakingV4 events raised by the ICnStaking contract.
type ICnStakingDeployCnStakingV4Iterator struct {
	Event *ICnStakingDeployCnStakingV4 // Event containing the contract specifics and raw log

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
func (it *ICnStakingDeployCnStakingV4Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingDeployCnStakingV4)
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
		it.Event = new(ICnStakingDeployCnStakingV4)
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
func (it *ICnStakingDeployCnStakingV4Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingDeployCnStakingV4Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingDeployCnStakingV4 represents a DeployCnStakingV4 event raised by the ICnStaking contract.
type ICnStakingDeployCnStakingV4 struct {
	ContractType     string
	PublicDelegation common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDeployCnStakingV4 is a free log retrieval operation binding the contract event 0x7ecd2f575e249099716283fbe0b54d644cc3aaec0465b36cb9e541aa0ed1dd94.
//
// Solidity: event DeployCnStakingV4(string contractType, address indexed publicDelegation)
func (_ICnStaking *ICnStakingFilterer) FilterDeployCnStakingV4(opts *bind.FilterOpts, publicDelegation []common.Address) (*ICnStakingDeployCnStakingV4Iterator, error) {

	var publicDelegationRule []interface{}
	for _, publicDelegationItem := range publicDelegation {
		publicDelegationRule = append(publicDelegationRule, publicDelegationItem)
	}

	logs, sub, err := _ICnStaking.contract.FilterLogs(opts, "DeployCnStakingV4", publicDelegationRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingDeployCnStakingV4Iterator{contract: _ICnStaking.contract, event: "DeployCnStakingV4", logs: logs, sub: sub}, nil
}

// WatchDeployCnStakingV4 is a free log subscription operation binding the contract event 0x7ecd2f575e249099716283fbe0b54d644cc3aaec0465b36cb9e541aa0ed1dd94.
//
// Solidity: event DeployCnStakingV4(string contractType, address indexed publicDelegation)
func (_ICnStaking *ICnStakingFilterer) WatchDeployCnStakingV4(opts *bind.WatchOpts, sink chan<- *ICnStakingDeployCnStakingV4, publicDelegation []common.Address) (event.Subscription, error) {

	var publicDelegationRule []interface{}
	for _, publicDelegationItem := range publicDelegation {
		publicDelegationRule = append(publicDelegationRule, publicDelegationItem)
	}

	logs, sub, err := _ICnStaking.contract.WatchLogs(opts, "DeployCnStakingV4", publicDelegationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingDeployCnStakingV4)
				if err := _ICnStaking.contract.UnpackLog(event, "DeployCnStakingV4", log); err != nil {
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

// ParseDeployCnStakingV4 is a log parse operation binding the contract event 0x7ecd2f575e249099716283fbe0b54d644cc3aaec0465b36cb9e541aa0ed1dd94.
//
// Solidity: event DeployCnStakingV4(string contractType, address indexed publicDelegation)
func (_ICnStaking *ICnStakingFilterer) ParseDeployCnStakingV4(log types.Log) (*ICnStakingDeployCnStakingV4, error) {
	event := new(ICnStakingDeployCnStakingV4)
	if err := _ICnStaking.contract.UnpackLog(event, "DeployCnStakingV4", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICnStakingHandleRedelegationIterator is returned from FilterHandleRedelegation and is used to iterate over the raw logs and unpacked data for HandleRedelegation events raised by the ICnStaking contract.
type ICnStakingHandleRedelegationIterator struct {
	Event *ICnStakingHandleRedelegation // Event containing the contract specifics and raw log

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
func (it *ICnStakingHandleRedelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingHandleRedelegation)
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
		it.Event = new(ICnStakingHandleRedelegation)
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
func (it *ICnStakingHandleRedelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingHandleRedelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingHandleRedelegation represents a HandleRedelegation event raised by the ICnStaking contract.
type ICnStakingHandleRedelegation struct {
	User            common.Address
	PrevCnStaking   common.Address
	TargetCnStaking common.Address
	Value           *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterHandleRedelegation is a free log retrieval operation binding the contract event 0xcbbd772ac9792b90fb77aa3194b0279ec4dbea1880ecc3674ecd2f4bf39b0c90.
//
// Solidity: event HandleRedelegation(address indexed user, address indexed prevCnStaking, address indexed targetCnStaking, uint256 value)
func (_ICnStaking *ICnStakingFilterer) FilterHandleRedelegation(opts *bind.FilterOpts, user []common.Address, prevCnStaking []common.Address, targetCnStaking []common.Address) (*ICnStakingHandleRedelegationIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var prevCnStakingRule []interface{}
	for _, prevCnStakingItem := range prevCnStaking {
		prevCnStakingRule = append(prevCnStakingRule, prevCnStakingItem)
	}
	var targetCnStakingRule []interface{}
	for _, targetCnStakingItem := range targetCnStaking {
		targetCnStakingRule = append(targetCnStakingRule, targetCnStakingItem)
	}

	logs, sub, err := _ICnStaking.contract.FilterLogs(opts, "HandleRedelegation", userRule, prevCnStakingRule, targetCnStakingRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingHandleRedelegationIterator{contract: _ICnStaking.contract, event: "HandleRedelegation", logs: logs, sub: sub}, nil
}

// WatchHandleRedelegation is a free log subscription operation binding the contract event 0xcbbd772ac9792b90fb77aa3194b0279ec4dbea1880ecc3674ecd2f4bf39b0c90.
//
// Solidity: event HandleRedelegation(address indexed user, address indexed prevCnStaking, address indexed targetCnStaking, uint256 value)
func (_ICnStaking *ICnStakingFilterer) WatchHandleRedelegation(opts *bind.WatchOpts, sink chan<- *ICnStakingHandleRedelegation, user []common.Address, prevCnStaking []common.Address, targetCnStaking []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var prevCnStakingRule []interface{}
	for _, prevCnStakingItem := range prevCnStaking {
		prevCnStakingRule = append(prevCnStakingRule, prevCnStakingItem)
	}
	var targetCnStakingRule []interface{}
	for _, targetCnStakingItem := range targetCnStaking {
		targetCnStakingRule = append(targetCnStakingRule, targetCnStakingItem)
	}

	logs, sub, err := _ICnStaking.contract.WatchLogs(opts, "HandleRedelegation", userRule, prevCnStakingRule, targetCnStakingRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingHandleRedelegation)
				if err := _ICnStaking.contract.UnpackLog(event, "HandleRedelegation", log); err != nil {
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

// ParseHandleRedelegation is a log parse operation binding the contract event 0xcbbd772ac9792b90fb77aa3194b0279ec4dbea1880ecc3674ecd2f4bf39b0c90.
//
// Solidity: event HandleRedelegation(address indexed user, address indexed prevCnStaking, address indexed targetCnStaking, uint256 value)
func (_ICnStaking *ICnStakingFilterer) ParseHandleRedelegation(log types.Log) (*ICnStakingHandleRedelegation, error) {
	event := new(ICnStakingHandleRedelegation)
	if err := _ICnStaking.contract.UnpackLog(event, "HandleRedelegation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICnStakingRedelegationIterator is returned from FilterRedelegation and is used to iterate over the raw logs and unpacked data for Redelegation events raised by the ICnStaking contract.
type ICnStakingRedelegationIterator struct {
	Event *ICnStakingRedelegation // Event containing the contract specifics and raw log

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
func (it *ICnStakingRedelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingRedelegation)
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
		it.Event = new(ICnStakingRedelegation)
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
func (it *ICnStakingRedelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingRedelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingRedelegation represents a Redelegation event raised by the ICnStaking contract.
type ICnStakingRedelegation struct {
	User            common.Address
	TargetCnStaking common.Address
	Value           *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRedelegation is a free log retrieval operation binding the contract event 0x4830c21a6eb4b2b3af6c4d24cfbd78ce437d5b690aac56ceec54df1d4d24e1a0.
//
// Solidity: event Redelegation(address indexed user, address indexed targetCnStaking, uint256 value)
func (_ICnStaking *ICnStakingFilterer) FilterRedelegation(opts *bind.FilterOpts, user []common.Address, targetCnStaking []common.Address) (*ICnStakingRedelegationIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var targetCnStakingRule []interface{}
	for _, targetCnStakingItem := range targetCnStaking {
		targetCnStakingRule = append(targetCnStakingRule, targetCnStakingItem)
	}

	logs, sub, err := _ICnStaking.contract.FilterLogs(opts, "Redelegation", userRule, targetCnStakingRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingRedelegationIterator{contract: _ICnStaking.contract, event: "Redelegation", logs: logs, sub: sub}, nil
}

// WatchRedelegation is a free log subscription operation binding the contract event 0x4830c21a6eb4b2b3af6c4d24cfbd78ce437d5b690aac56ceec54df1d4d24e1a0.
//
// Solidity: event Redelegation(address indexed user, address indexed targetCnStaking, uint256 value)
func (_ICnStaking *ICnStakingFilterer) WatchRedelegation(opts *bind.WatchOpts, sink chan<- *ICnStakingRedelegation, user []common.Address, targetCnStaking []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var targetCnStakingRule []interface{}
	for _, targetCnStakingItem := range targetCnStaking {
		targetCnStakingRule = append(targetCnStakingRule, targetCnStakingItem)
	}

	logs, sub, err := _ICnStaking.contract.WatchLogs(opts, "Redelegation", userRule, targetCnStakingRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingRedelegation)
				if err := _ICnStaking.contract.UnpackLog(event, "Redelegation", log); err != nil {
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

// ParseRedelegation is a log parse operation binding the contract event 0x4830c21a6eb4b2b3af6c4d24cfbd78ce437d5b690aac56ceec54df1d4d24e1a0.
//
// Solidity: event Redelegation(address indexed user, address indexed targetCnStaking, uint256 value)
func (_ICnStaking *ICnStakingFilterer) ParseRedelegation(log types.Log) (*ICnStakingRedelegation, error) {
	event := new(ICnStakingRedelegation)
	if err := _ICnStaking.contract.UnpackLog(event, "Redelegation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICnStakingWithdrawApprovedStakingIterator is returned from FilterWithdrawApprovedStaking and is used to iterate over the raw logs and unpacked data for WithdrawApprovedStaking events raised by the ICnStaking contract.
type ICnStakingWithdrawApprovedStakingIterator struct {
	Event *ICnStakingWithdrawApprovedStaking // Event containing the contract specifics and raw log

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
func (it *ICnStakingWithdrawApprovedStakingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingWithdrawApprovedStaking)
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
		it.Event = new(ICnStakingWithdrawApprovedStaking)
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
func (it *ICnStakingWithdrawApprovedStakingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingWithdrawApprovedStakingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingWithdrawApprovedStaking represents a WithdrawApprovedStaking event raised by the ICnStaking contract.
type ICnStakingWithdrawApprovedStaking struct {
	ApprovedWithdrawalId *big.Int
	To                   common.Address
	Value                *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterWithdrawApprovedStaking is a free log retrieval operation binding the contract event 0x7a2d48d1df1249429730a253e5713a7d7a2024913de2fbccafdf36efdd32bf05.
//
// Solidity: event WithdrawApprovedStaking(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value)
func (_ICnStaking *ICnStakingFilterer) FilterWithdrawApprovedStaking(opts *bind.FilterOpts, approvedWithdrawalId []*big.Int, to []common.Address) (*ICnStakingWithdrawApprovedStakingIterator, error) {

	var approvedWithdrawalIdRule []interface{}
	for _, approvedWithdrawalIdItem := range approvedWithdrawalId {
		approvedWithdrawalIdRule = append(approvedWithdrawalIdRule, approvedWithdrawalIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ICnStaking.contract.FilterLogs(opts, "WithdrawApprovedStaking", approvedWithdrawalIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingWithdrawApprovedStakingIterator{contract: _ICnStaking.contract, event: "WithdrawApprovedStaking", logs: logs, sub: sub}, nil
}

// WatchWithdrawApprovedStaking is a free log subscription operation binding the contract event 0x7a2d48d1df1249429730a253e5713a7d7a2024913de2fbccafdf36efdd32bf05.
//
// Solidity: event WithdrawApprovedStaking(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value)
func (_ICnStaking *ICnStakingFilterer) WatchWithdrawApprovedStaking(opts *bind.WatchOpts, sink chan<- *ICnStakingWithdrawApprovedStaking, approvedWithdrawalId []*big.Int, to []common.Address) (event.Subscription, error) {

	var approvedWithdrawalIdRule []interface{}
	for _, approvedWithdrawalIdItem := range approvedWithdrawalId {
		approvedWithdrawalIdRule = append(approvedWithdrawalIdRule, approvedWithdrawalIdItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ICnStaking.contract.WatchLogs(opts, "WithdrawApprovedStaking", approvedWithdrawalIdRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingWithdrawApprovedStaking)
				if err := _ICnStaking.contract.UnpackLog(event, "WithdrawApprovedStaking", log); err != nil {
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

// ParseWithdrawApprovedStaking is a log parse operation binding the contract event 0x7a2d48d1df1249429730a253e5713a7d7a2024913de2fbccafdf36efdd32bf05.
//
// Solidity: event WithdrawApprovedStaking(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value)
func (_ICnStaking *ICnStakingFilterer) ParseWithdrawApprovedStaking(log types.Log) (*ICnStakingWithdrawApprovedStaking, error) {
	event := new(ICnStakingWithdrawApprovedStaking)
	if err := _ICnStaking.contract.UnpackLog(event, "WithdrawApprovedStaking", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC1155ErrorsMetaData contains all meta data concerning the IERC1155Errors contract.
var IERC1155ErrorsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC1155InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"idsLength\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valuesLength\",\"type\":\"uint256\"}],\"name\":\"ERC1155InvalidArrayLength\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC1155InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC1155MissingApprovalForAll\",\"type\":\"error\"}]",
}

// IERC1155ErrorsABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1155ErrorsMetaData.ABI instead.
var IERC1155ErrorsABI = IERC1155ErrorsMetaData.ABI

// IERC1155ErrorsBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1155ErrorsBinRuntime = ``

// IERC1155Errors is an auto generated Go binding around a Kaia contract.
type IERC1155Errors struct {
	IERC1155ErrorsCaller     // Read-only binding to the contract
	IERC1155ErrorsTransactor // Write-only binding to the contract
	IERC1155ErrorsFilterer   // Log filterer for contract events
}

// IERC1155ErrorsCaller is an auto generated read-only Go binding around a Kaia contract.
type IERC1155ErrorsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1155ErrorsTransactor is an auto generated write-only Go binding around a Kaia contract.
type IERC1155ErrorsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1155ErrorsFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IERC1155ErrorsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1155ErrorsSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IERC1155ErrorsSession struct {
	Contract     *IERC1155Errors   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC1155ErrorsCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IERC1155ErrorsCallerSession struct {
	Contract *IERC1155ErrorsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IERC1155ErrorsTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IERC1155ErrorsTransactorSession struct {
	Contract     *IERC1155ErrorsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IERC1155ErrorsRaw is an auto generated low-level Go binding around a Kaia contract.
type IERC1155ErrorsRaw struct {
	Contract *IERC1155Errors // Generic contract binding to access the raw methods on
}

// IERC1155ErrorsCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IERC1155ErrorsCallerRaw struct {
	Contract *IERC1155ErrorsCaller // Generic read-only contract binding to access the raw methods on
}

// IERC1155ErrorsTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IERC1155ErrorsTransactorRaw struct {
	Contract *IERC1155ErrorsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1155Errors creates a new instance of IERC1155Errors, bound to a specific deployed contract.
func NewIERC1155Errors(address common.Address, backend bind.ContractBackend) (*IERC1155Errors, error) {
	contract, err := bindIERC1155Errors(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1155Errors{IERC1155ErrorsCaller: IERC1155ErrorsCaller{contract: contract}, IERC1155ErrorsTransactor: IERC1155ErrorsTransactor{contract: contract}, IERC1155ErrorsFilterer: IERC1155ErrorsFilterer{contract: contract}}, nil
}

// NewIERC1155ErrorsCaller creates a new read-only instance of IERC1155Errors, bound to a specific deployed contract.
func NewIERC1155ErrorsCaller(address common.Address, caller bind.ContractCaller) (*IERC1155ErrorsCaller, error) {
	contract, err := bindIERC1155Errors(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1155ErrorsCaller{contract: contract}, nil
}

// NewIERC1155ErrorsTransactor creates a new write-only instance of IERC1155Errors, bound to a specific deployed contract.
func NewIERC1155ErrorsTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC1155ErrorsTransactor, error) {
	contract, err := bindIERC1155Errors(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1155ErrorsTransactor{contract: contract}, nil
}

// NewIERC1155ErrorsFilterer creates a new log filterer instance of IERC1155Errors, bound to a specific deployed contract.
func NewIERC1155ErrorsFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC1155ErrorsFilterer, error) {
	contract, err := bindIERC1155Errors(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1155ErrorsFilterer{contract: contract}, nil
}

// bindIERC1155Errors binds a generic wrapper to an already deployed contract.
func bindIERC1155Errors(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1155ErrorsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1155Errors *IERC1155ErrorsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1155Errors.Contract.IERC1155ErrorsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1155Errors *IERC1155ErrorsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1155Errors.Contract.IERC1155ErrorsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1155Errors *IERC1155ErrorsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1155Errors.Contract.IERC1155ErrorsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1155Errors *IERC1155ErrorsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1155Errors.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1155Errors *IERC1155ErrorsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1155Errors.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1155Errors *IERC1155ErrorsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1155Errors.Contract.contract.Transact(opts, method, params...)
}

// IERC20MetaData contains all meta data concerning the IERC20 contract.
var IERC20MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
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

// IERC20ErrorsMetaData contains all meta data concerning the IERC20Errors contract.
var IERC20ErrorsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"}]",
}

// IERC20ErrorsABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC20ErrorsMetaData.ABI instead.
var IERC20ErrorsABI = IERC20ErrorsMetaData.ABI

// IERC20ErrorsBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC20ErrorsBinRuntime = ``

// IERC20Errors is an auto generated Go binding around a Kaia contract.
type IERC20Errors struct {
	IERC20ErrorsCaller     // Read-only binding to the contract
	IERC20ErrorsTransactor // Write-only binding to the contract
	IERC20ErrorsFilterer   // Log filterer for contract events
}

// IERC20ErrorsCaller is an auto generated read-only Go binding around a Kaia contract.
type IERC20ErrorsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20ErrorsTransactor is an auto generated write-only Go binding around a Kaia contract.
type IERC20ErrorsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20ErrorsFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IERC20ErrorsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20ErrorsSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IERC20ErrorsSession struct {
	Contract     *IERC20Errors     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20ErrorsCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IERC20ErrorsCallerSession struct {
	Contract *IERC20ErrorsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// IERC20ErrorsTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IERC20ErrorsTransactorSession struct {
	Contract     *IERC20ErrorsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// IERC20ErrorsRaw is an auto generated low-level Go binding around a Kaia contract.
type IERC20ErrorsRaw struct {
	Contract *IERC20Errors // Generic contract binding to access the raw methods on
}

// IERC20ErrorsCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IERC20ErrorsCallerRaw struct {
	Contract *IERC20ErrorsCaller // Generic read-only contract binding to access the raw methods on
}

// IERC20ErrorsTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IERC20ErrorsTransactorRaw struct {
	Contract *IERC20ErrorsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC20Errors creates a new instance of IERC20Errors, bound to a specific deployed contract.
func NewIERC20Errors(address common.Address, backend bind.ContractBackend) (*IERC20Errors, error) {
	contract, err := bindIERC20Errors(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20Errors{IERC20ErrorsCaller: IERC20ErrorsCaller{contract: contract}, IERC20ErrorsTransactor: IERC20ErrorsTransactor{contract: contract}, IERC20ErrorsFilterer: IERC20ErrorsFilterer{contract: contract}}, nil
}

// NewIERC20ErrorsCaller creates a new read-only instance of IERC20Errors, bound to a specific deployed contract.
func NewIERC20ErrorsCaller(address common.Address, caller bind.ContractCaller) (*IERC20ErrorsCaller, error) {
	contract, err := bindIERC20Errors(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20ErrorsCaller{contract: contract}, nil
}

// NewIERC20ErrorsTransactor creates a new write-only instance of IERC20Errors, bound to a specific deployed contract.
func NewIERC20ErrorsTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC20ErrorsTransactor, error) {
	contract, err := bindIERC20Errors(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20ErrorsTransactor{contract: contract}, nil
}

// NewIERC20ErrorsFilterer creates a new log filterer instance of IERC20Errors, bound to a specific deployed contract.
func NewIERC20ErrorsFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC20ErrorsFilterer, error) {
	contract, err := bindIERC20Errors(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC20ErrorsFilterer{contract: contract}, nil
}

// bindIERC20Errors binds a generic wrapper to an already deployed contract.
func bindIERC20Errors(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC20ErrorsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20Errors *IERC20ErrorsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20Errors.Contract.IERC20ErrorsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20Errors *IERC20ErrorsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20Errors.Contract.IERC20ErrorsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20Errors *IERC20ErrorsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20Errors.Contract.IERC20ErrorsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20Errors *IERC20ErrorsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20Errors.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20Errors *IERC20ErrorsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20Errors.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20Errors *IERC20ErrorsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20Errors.Contract.contract.Transact(opts, method, params...)
}

// IERC20MetadataMetaData contains all meta data concerning the IERC20Metadata contract.
var IERC20MetadataMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// IERC20MetadataABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC20MetadataMetaData.ABI instead.
var IERC20MetadataABI = IERC20MetadataMetaData.ABI

// IERC20MetadataBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC20MetadataBinRuntime = ``

// Deprecated: Use IERC20MetadataMetaData.Sigs instead.
// IERC20MetadataFuncSigs maps the 4-byte function signature to its string representation.
var IERC20MetadataFuncSigs = IERC20MetadataMetaData.Sigs

// IERC20Metadata is an auto generated Go binding around a Kaia contract.
type IERC20Metadata struct {
	IERC20MetadataCaller     // Read-only binding to the contract
	IERC20MetadataTransactor // Write-only binding to the contract
	IERC20MetadataFilterer   // Log filterer for contract events
}

// IERC20MetadataCaller is an auto generated read-only Go binding around a Kaia contract.
type IERC20MetadataCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20MetadataTransactor is an auto generated write-only Go binding around a Kaia contract.
type IERC20MetadataTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20MetadataFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IERC20MetadataFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC20MetadataSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IERC20MetadataSession struct {
	Contract     *IERC20Metadata   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC20MetadataCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IERC20MetadataCallerSession struct {
	Contract *IERC20MetadataCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// IERC20MetadataTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IERC20MetadataTransactorSession struct {
	Contract     *IERC20MetadataTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// IERC20MetadataRaw is an auto generated low-level Go binding around a Kaia contract.
type IERC20MetadataRaw struct {
	Contract *IERC20Metadata // Generic contract binding to access the raw methods on
}

// IERC20MetadataCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IERC20MetadataCallerRaw struct {
	Contract *IERC20MetadataCaller // Generic read-only contract binding to access the raw methods on
}

// IERC20MetadataTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IERC20MetadataTransactorRaw struct {
	Contract *IERC20MetadataTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC20Metadata creates a new instance of IERC20Metadata, bound to a specific deployed contract.
func NewIERC20Metadata(address common.Address, backend bind.ContractBackend) (*IERC20Metadata, error) {
	contract, err := bindIERC20Metadata(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC20Metadata{IERC20MetadataCaller: IERC20MetadataCaller{contract: contract}, IERC20MetadataTransactor: IERC20MetadataTransactor{contract: contract}, IERC20MetadataFilterer: IERC20MetadataFilterer{contract: contract}}, nil
}

// NewIERC20MetadataCaller creates a new read-only instance of IERC20Metadata, bound to a specific deployed contract.
func NewIERC20MetadataCaller(address common.Address, caller bind.ContractCaller) (*IERC20MetadataCaller, error) {
	contract, err := bindIERC20Metadata(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20MetadataCaller{contract: contract}, nil
}

// NewIERC20MetadataTransactor creates a new write-only instance of IERC20Metadata, bound to a specific deployed contract.
func NewIERC20MetadataTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC20MetadataTransactor, error) {
	contract, err := bindIERC20Metadata(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC20MetadataTransactor{contract: contract}, nil
}

// NewIERC20MetadataFilterer creates a new log filterer instance of IERC20Metadata, bound to a specific deployed contract.
func NewIERC20MetadataFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC20MetadataFilterer, error) {
	contract, err := bindIERC20Metadata(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC20MetadataFilterer{contract: contract}, nil
}

// bindIERC20Metadata binds a generic wrapper to an already deployed contract.
func bindIERC20Metadata(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC20MetadataMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20Metadata *IERC20MetadataRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20Metadata.Contract.IERC20MetadataCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20Metadata *IERC20MetadataRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.IERC20MetadataTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20Metadata *IERC20MetadataRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.IERC20MetadataTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC20Metadata *IERC20MetadataCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC20Metadata.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC20Metadata *IERC20MetadataTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC20Metadata *IERC20MetadataTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20Metadata *IERC20MetadataCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20Metadata.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20Metadata *IERC20MetadataSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20Metadata.Contract.Allowance(&_IERC20Metadata.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_IERC20Metadata *IERC20MetadataCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _IERC20Metadata.Contract.Allowance(&_IERC20Metadata.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20Metadata *IERC20MetadataCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IERC20Metadata.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20Metadata *IERC20MetadataSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20Metadata.Contract.BalanceOf(&_IERC20Metadata.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_IERC20Metadata *IERC20MetadataCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _IERC20Metadata.Contract.BalanceOf(&_IERC20Metadata.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_IERC20Metadata *IERC20MetadataCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _IERC20Metadata.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_IERC20Metadata *IERC20MetadataSession) Decimals() (uint8, error) {
	return _IERC20Metadata.Contract.Decimals(&_IERC20Metadata.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_IERC20Metadata *IERC20MetadataCallerSession) Decimals() (uint8, error) {
	return _IERC20Metadata.Contract.Decimals(&_IERC20Metadata.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IERC20Metadata *IERC20MetadataCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IERC20Metadata.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IERC20Metadata *IERC20MetadataSession) Name() (string, error) {
	return _IERC20Metadata.Contract.Name(&_IERC20Metadata.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_IERC20Metadata *IERC20MetadataCallerSession) Name() (string, error) {
	return _IERC20Metadata.Contract.Name(&_IERC20Metadata.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_IERC20Metadata *IERC20MetadataCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IERC20Metadata.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_IERC20Metadata *IERC20MetadataSession) Symbol() (string, error) {
	return _IERC20Metadata.Contract.Symbol(&_IERC20Metadata.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_IERC20Metadata *IERC20MetadataCallerSession) Symbol() (string, error) {
	return _IERC20Metadata.Contract.Symbol(&_IERC20Metadata.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20Metadata *IERC20MetadataCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IERC20Metadata.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20Metadata *IERC20MetadataSession) TotalSupply() (*big.Int, error) {
	return _IERC20Metadata.Contract.TotalSupply(&_IERC20Metadata.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_IERC20Metadata *IERC20MetadataCallerSession) TotalSupply() (*big.Int, error) {
	return _IERC20Metadata.Contract.TotalSupply(&_IERC20Metadata.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.Approve(&_IERC20Metadata.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.Approve(&_IERC20Metadata.TransactOpts, spender, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.Transfer(&_IERC20Metadata.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.Transfer(&_IERC20Metadata.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.TransferFrom(&_IERC20Metadata.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_IERC20Metadata *IERC20MetadataTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _IERC20Metadata.Contract.TransferFrom(&_IERC20Metadata.TransactOpts, from, to, value)
}

// IERC20MetadataApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the IERC20Metadata contract.
type IERC20MetadataApprovalIterator struct {
	Event *IERC20MetadataApproval // Event containing the contract specifics and raw log

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
func (it *IERC20MetadataApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20MetadataApproval)
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
		it.Event = new(IERC20MetadataApproval)
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
func (it *IERC20MetadataApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20MetadataApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20MetadataApproval represents a Approval event raised by the IERC20Metadata contract.
type IERC20MetadataApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20Metadata *IERC20MetadataFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*IERC20MetadataApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20Metadata.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &IERC20MetadataApprovalIterator{contract: _IERC20Metadata.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_IERC20Metadata *IERC20MetadataFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *IERC20MetadataApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _IERC20Metadata.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20MetadataApproval)
				if err := _IERC20Metadata.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_IERC20Metadata *IERC20MetadataFilterer) ParseApproval(log types.Log) (*IERC20MetadataApproval, error) {
	event := new(IERC20MetadataApproval)
	if err := _IERC20Metadata.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC20MetadataTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the IERC20Metadata contract.
type IERC20MetadataTransferIterator struct {
	Event *IERC20MetadataTransfer // Event containing the contract specifics and raw log

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
func (it *IERC20MetadataTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC20MetadataTransfer)
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
		it.Event = new(IERC20MetadataTransfer)
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
func (it *IERC20MetadataTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC20MetadataTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC20MetadataTransfer represents a Transfer event raised by the IERC20Metadata contract.
type IERC20MetadataTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20Metadata *IERC20MetadataFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*IERC20MetadataTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20Metadata.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &IERC20MetadataTransferIterator{contract: _IERC20Metadata.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_IERC20Metadata *IERC20MetadataFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *IERC20MetadataTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _IERC20Metadata.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC20MetadataTransfer)
				if err := _IERC20Metadata.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_IERC20Metadata *IERC20MetadataFilterer) ParseTransfer(log types.Log) (*IERC20MetadataTransfer, error) {
	event := new(IERC20MetadataTransfer)
	if err := _IERC20Metadata.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC721ErrorsMetaData contains all meta data concerning the IERC721Errors contract.
var IERC721ErrorsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"}]",
}

// IERC721ErrorsABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC721ErrorsMetaData.ABI instead.
var IERC721ErrorsABI = IERC721ErrorsMetaData.ABI

// IERC721ErrorsBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC721ErrorsBinRuntime = ``

// IERC721Errors is an auto generated Go binding around a Kaia contract.
type IERC721Errors struct {
	IERC721ErrorsCaller     // Read-only binding to the contract
	IERC721ErrorsTransactor // Write-only binding to the contract
	IERC721ErrorsFilterer   // Log filterer for contract events
}

// IERC721ErrorsCaller is an auto generated read-only Go binding around a Kaia contract.
type IERC721ErrorsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721ErrorsTransactor is an auto generated write-only Go binding around a Kaia contract.
type IERC721ErrorsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721ErrorsFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IERC721ErrorsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC721ErrorsSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IERC721ErrorsSession struct {
	Contract     *IERC721Errors    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC721ErrorsCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IERC721ErrorsCallerSession struct {
	Contract *IERC721ErrorsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// IERC721ErrorsTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IERC721ErrorsTransactorSession struct {
	Contract     *IERC721ErrorsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// IERC721ErrorsRaw is an auto generated low-level Go binding around a Kaia contract.
type IERC721ErrorsRaw struct {
	Contract *IERC721Errors // Generic contract binding to access the raw methods on
}

// IERC721ErrorsCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IERC721ErrorsCallerRaw struct {
	Contract *IERC721ErrorsCaller // Generic read-only contract binding to access the raw methods on
}

// IERC721ErrorsTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IERC721ErrorsTransactorRaw struct {
	Contract *IERC721ErrorsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC721Errors creates a new instance of IERC721Errors, bound to a specific deployed contract.
func NewIERC721Errors(address common.Address, backend bind.ContractBackend) (*IERC721Errors, error) {
	contract, err := bindIERC721Errors(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC721Errors{IERC721ErrorsCaller: IERC721ErrorsCaller{contract: contract}, IERC721ErrorsTransactor: IERC721ErrorsTransactor{contract: contract}, IERC721ErrorsFilterer: IERC721ErrorsFilterer{contract: contract}}, nil
}

// NewIERC721ErrorsCaller creates a new read-only instance of IERC721Errors, bound to a specific deployed contract.
func NewIERC721ErrorsCaller(address common.Address, caller bind.ContractCaller) (*IERC721ErrorsCaller, error) {
	contract, err := bindIERC721Errors(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721ErrorsCaller{contract: contract}, nil
}

// NewIERC721ErrorsTransactor creates a new write-only instance of IERC721Errors, bound to a specific deployed contract.
func NewIERC721ErrorsTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC721ErrorsTransactor, error) {
	contract, err := bindIERC721Errors(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC721ErrorsTransactor{contract: contract}, nil
}

// NewIERC721ErrorsFilterer creates a new log filterer instance of IERC721Errors, bound to a specific deployed contract.
func NewIERC721ErrorsFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC721ErrorsFilterer, error) {
	contract, err := bindIERC721Errors(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC721ErrorsFilterer{contract: contract}, nil
}

// bindIERC721Errors binds a generic wrapper to an already deployed contract.
func bindIERC721Errors(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC721ErrorsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721Errors *IERC721ErrorsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC721Errors.Contract.IERC721ErrorsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721Errors *IERC721ErrorsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721Errors.Contract.IERC721ErrorsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721Errors *IERC721ErrorsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721Errors.Contract.IERC721ErrorsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC721Errors *IERC721ErrorsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC721Errors.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC721Errors *IERC721ErrorsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC721Errors.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC721Errors *IERC721ErrorsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC721Errors.Contract.contract.Transact(opts, method, params...)
}

// IKIP163MetaData contains all meta data concerning the IKIP163 contract.
var IKIP163MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"stakeFor\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"228cb733": "reward()",
		"4bf69206": "stakeFor(address)",
	},
}

// IKIP163ABI is the input ABI used to generate the binding from.
// Deprecated: Use IKIP163MetaData.ABI instead.
var IKIP163ABI = IKIP163MetaData.ABI

// IKIP163BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IKIP163BinRuntime = ``

// Deprecated: Use IKIP163MetaData.Sigs instead.
// IKIP163FuncSigs maps the 4-byte function signature to its string representation.
var IKIP163FuncSigs = IKIP163MetaData.Sigs

// IKIP163 is an auto generated Go binding around a Kaia contract.
type IKIP163 struct {
	IKIP163Caller     // Read-only binding to the contract
	IKIP163Transactor // Write-only binding to the contract
	IKIP163Filterer   // Log filterer for contract events
}

// IKIP163Caller is an auto generated read-only Go binding around a Kaia contract.
type IKIP163Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP163Transactor is an auto generated write-only Go binding around a Kaia contract.
type IKIP163Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP163Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type IKIP163Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP163Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IKIP163Session struct {
	Contract     *IKIP163          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IKIP163CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IKIP163CallerSession struct {
	Contract *IKIP163Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IKIP163TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IKIP163TransactorSession struct {
	Contract     *IKIP163Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IKIP163Raw is an auto generated low-level Go binding around a Kaia contract.
type IKIP163Raw struct {
	Contract *IKIP163 // Generic contract binding to access the raw methods on
}

// IKIP163CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IKIP163CallerRaw struct {
	Contract *IKIP163Caller // Generic read-only contract binding to access the raw methods on
}

// IKIP163TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IKIP163TransactorRaw struct {
	Contract *IKIP163Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIKIP163 creates a new instance of IKIP163, bound to a specific deployed contract.
func NewIKIP163(address common.Address, backend bind.ContractBackend) (*IKIP163, error) {
	contract, err := bindIKIP163(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IKIP163{IKIP163Caller: IKIP163Caller{contract: contract}, IKIP163Transactor: IKIP163Transactor{contract: contract}, IKIP163Filterer: IKIP163Filterer{contract: contract}}, nil
}

// NewIKIP163Caller creates a new read-only instance of IKIP163, bound to a specific deployed contract.
func NewIKIP163Caller(address common.Address, caller bind.ContractCaller) (*IKIP163Caller, error) {
	contract, err := bindIKIP163(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IKIP163Caller{contract: contract}, nil
}

// NewIKIP163Transactor creates a new write-only instance of IKIP163, bound to a specific deployed contract.
func NewIKIP163Transactor(address common.Address, transactor bind.ContractTransactor) (*IKIP163Transactor, error) {
	contract, err := bindIKIP163(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IKIP163Transactor{contract: contract}, nil
}

// NewIKIP163Filterer creates a new log filterer instance of IKIP163, bound to a specific deployed contract.
func NewIKIP163Filterer(address common.Address, filterer bind.ContractFilterer) (*IKIP163Filterer, error) {
	contract, err := bindIKIP163(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IKIP163Filterer{contract: contract}, nil
}

// bindIKIP163 binds a generic wrapper to an already deployed contract.
func bindIKIP163(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IKIP163MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IKIP163 *IKIP163Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKIP163.Contract.IKIP163Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IKIP163 *IKIP163Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKIP163.Contract.IKIP163Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IKIP163 *IKIP163Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKIP163.Contract.IKIP163Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IKIP163 *IKIP163CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKIP163.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IKIP163 *IKIP163TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKIP163.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IKIP163 *IKIP163TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKIP163.Contract.contract.Transact(opts, method, params...)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_IKIP163 *IKIP163Caller) Reward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IKIP163.contract.Call(opts, &out, "reward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_IKIP163 *IKIP163Session) Reward() (*big.Int, error) {
	return _IKIP163.Contract.Reward(&_IKIP163.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_IKIP163 *IKIP163CallerSession) Reward() (*big.Int, error) {
	return _IKIP163.Contract.Reward(&_IKIP163.CallOpts)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_IKIP163 *IKIP163Transactor) StakeFor(opts *bind.TransactOpts, _recipient common.Address) (*types.Transaction, error) {
	return _IKIP163.contract.Transact(opts, "stakeFor", _recipient)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_IKIP163 *IKIP163Session) StakeFor(_recipient common.Address) (*types.Transaction, error) {
	return _IKIP163.Contract.StakeFor(&_IKIP163.TransactOpts, _recipient)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_IKIP163 *IKIP163TransactorSession) StakeFor(_recipient common.Address) (*types.Transaction, error) {
	return _IKIP163.Contract.StakeFor(&_IKIP163.TransactOpts, _recipient)
}

// IPublicDelegationMetaData contains all meta data concerning the IPublicDelegation contract.
var IPublicDelegationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CommissionRateTooHigh\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotRequestOwner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RedelegateAmountTooLow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakeAmountTooLow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalAmountTooLow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"Claimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"contractType\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"baseCnStaking\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gcName\",\"type\":\"string\"}],\"indexed\":false,\"internalType\":\"structIPublicDelegation.PDConstructorArgs\",\"name\":\"pdArgs\",\"type\":\"tuple\"}],\"name\":\"DeployPublicDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"targetCnStaking\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"}],\"name\":\"Redelegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"RequestCancelWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"}],\"name\":\"RequestWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"commission\",\"type\":\"uint256\"}],\"name\":\"SendCommission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"prevCommissionRate\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"}],\"name\":\"UpdateCommissionRate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"prevCommissionTo\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"}],\"name\":\"UpdateCommissionTo\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMMISSION_DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CONTRACT_TYPE\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_COMMISSION_RATE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseCnStaking\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"cancelApprovedStakingWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commissionTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"convertToAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"convertToShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getCurrentWithdrawalRequestState\",\"outputs\":[{\"internalType\":\"enumIPublicDelegation.WithdrawalRequestState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"getUserRequestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"getUserRequestIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"enumIPublicDelegation.WithdrawalRequestState\",\"name\":\"_state\",\"type\":\"uint8\"}],\"name\":\"getUserRequestIdsWithState\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_baseCnStaking\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gcName\",\"type\":\"string\"}],\"internalType\":\"structIPublicDelegation.PDConstructorArgs\",\"name\":\"_args\",\"type\":\"tuple\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"maxRedeem\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"maxWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"previewDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"previewRedeem\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"previewWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_targetCnStaking\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"redelegateByAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_targetCnStaking\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"redelegateByShares\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"requestIdToOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"stakeFor\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sweep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_commissionRate\",\"type\":\"uint256\"}],\"name\":\"updateCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_commissionTo\",\"type\":\"address\"}],\"name\":\"updateCommissionTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"userRequestIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"3b1dbfcc": "COMMISSION_DENOMINATOR()",
		"4b6a94cc": "CONTRACT_TYPE()",
		"207239c0": "MAX_COMMISSION_RATE()",
		"ffa1ad74": "VERSION()",
		"0c11487f": "baseCnStaking()",
		"c804b115": "cancelApprovedStakingWithdrawal(uint256)",
		"379607f5": "claim(uint256)",
		"5ea1d6f8": "commissionRate()",
		"2f9ac83a": "commissionTo()",
		"07a2d13a": "convertToAssets(uint256)",
		"c6e6f592": "convertToShares(uint256)",
		"04ddc9d1": "getCurrentWithdrawalRequestState(uint256)",
		"c166c458": "getUserRequestCount(address)",
		"60df7c6c": "getUserRequestIds(address)",
		"93b89a84": "getUserRequestIdsWithState(address,uint8)",
		"26cf277a": "initialize(address,(address,address,uint256,string))",
		"d905777e": "maxRedeem(address)",
		"ce96cb77": "maxWithdraw(address)",
		"ef8b30f7": "previewDeposit(uint256)",
		"4cdad506": "previewRedeem(uint256)",
		"0a28a477": "previewWithdraw(uint256)",
		"1e9a6950": "redeem(address,uint256)",
		"e659d7d7": "redelegateByAssets(address,uint256)",
		"e15fc350": "redelegateByShares(address,uint256)",
		"f29177c3": "requestIdToOwner(uint256)",
		"228cb733": "reward()",
		"3a4b66f1": "stake()",
		"4bf69206": "stakeFor(address)",
		"35faa416": "sweep()",
		"01e1d114": "totalAssets()",
		"00fa3d50": "updateCommissionRate(uint256)",
		"052028d0": "updateCommissionTo(address)",
		"97feb23c": "userRequestIds(address,uint256)",
		"f3fef3a3": "withdraw(address,uint256)",
	},
}

// IPublicDelegationABI is the input ABI used to generate the binding from.
// Deprecated: Use IPublicDelegationMetaData.ABI instead.
var IPublicDelegationABI = IPublicDelegationMetaData.ABI

// IPublicDelegationBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IPublicDelegationBinRuntime = ``

// Deprecated: Use IPublicDelegationMetaData.Sigs instead.
// IPublicDelegationFuncSigs maps the 4-byte function signature to its string representation.
var IPublicDelegationFuncSigs = IPublicDelegationMetaData.Sigs

// IPublicDelegation is an auto generated Go binding around a Kaia contract.
type IPublicDelegation struct {
	IPublicDelegationCaller     // Read-only binding to the contract
	IPublicDelegationTransactor // Write-only binding to the contract
	IPublicDelegationFilterer   // Log filterer for contract events
}

// IPublicDelegationCaller is an auto generated read-only Go binding around a Kaia contract.
type IPublicDelegationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPublicDelegationTransactor is an auto generated write-only Go binding around a Kaia contract.
type IPublicDelegationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPublicDelegationFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IPublicDelegationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IPublicDelegationSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IPublicDelegationSession struct {
	Contract     *IPublicDelegation // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IPublicDelegationCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IPublicDelegationCallerSession struct {
	Contract *IPublicDelegationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IPublicDelegationTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IPublicDelegationTransactorSession struct {
	Contract     *IPublicDelegationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IPublicDelegationRaw is an auto generated low-level Go binding around a Kaia contract.
type IPublicDelegationRaw struct {
	Contract *IPublicDelegation // Generic contract binding to access the raw methods on
}

// IPublicDelegationCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IPublicDelegationCallerRaw struct {
	Contract *IPublicDelegationCaller // Generic read-only contract binding to access the raw methods on
}

// IPublicDelegationTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IPublicDelegationTransactorRaw struct {
	Contract *IPublicDelegationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIPublicDelegation creates a new instance of IPublicDelegation, bound to a specific deployed contract.
func NewIPublicDelegation(address common.Address, backend bind.ContractBackend) (*IPublicDelegation, error) {
	contract, err := bindIPublicDelegation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegation{IPublicDelegationCaller: IPublicDelegationCaller{contract: contract}, IPublicDelegationTransactor: IPublicDelegationTransactor{contract: contract}, IPublicDelegationFilterer: IPublicDelegationFilterer{contract: contract}}, nil
}

// NewIPublicDelegationCaller creates a new read-only instance of IPublicDelegation, bound to a specific deployed contract.
func NewIPublicDelegationCaller(address common.Address, caller bind.ContractCaller) (*IPublicDelegationCaller, error) {
	contract, err := bindIPublicDelegation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationCaller{contract: contract}, nil
}

// NewIPublicDelegationTransactor creates a new write-only instance of IPublicDelegation, bound to a specific deployed contract.
func NewIPublicDelegationTransactor(address common.Address, transactor bind.ContractTransactor) (*IPublicDelegationTransactor, error) {
	contract, err := bindIPublicDelegation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationTransactor{contract: contract}, nil
}

// NewIPublicDelegationFilterer creates a new log filterer instance of IPublicDelegation, bound to a specific deployed contract.
func NewIPublicDelegationFilterer(address common.Address, filterer bind.ContractFilterer) (*IPublicDelegationFilterer, error) {
	contract, err := bindIPublicDelegation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationFilterer{contract: contract}, nil
}

// bindIPublicDelegation binds a generic wrapper to an already deployed contract.
func bindIPublicDelegation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IPublicDelegationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPublicDelegation *IPublicDelegationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPublicDelegation.Contract.IPublicDelegationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPublicDelegation *IPublicDelegationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.IPublicDelegationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPublicDelegation *IPublicDelegationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.IPublicDelegationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IPublicDelegation *IPublicDelegationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IPublicDelegation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IPublicDelegation *IPublicDelegationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IPublicDelegation *IPublicDelegationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.contract.Transact(opts, method, params...)
}

// COMMISSIONDENOMINATOR is a free data retrieval call binding the contract method 0x3b1dbfcc.
//
// Solidity: function COMMISSION_DENOMINATOR() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) COMMISSIONDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "COMMISSION_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// COMMISSIONDENOMINATOR is a free data retrieval call binding the contract method 0x3b1dbfcc.
//
// Solidity: function COMMISSION_DENOMINATOR() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) COMMISSIONDENOMINATOR() (*big.Int, error) {
	return _IPublicDelegation.Contract.COMMISSIONDENOMINATOR(&_IPublicDelegation.CallOpts)
}

// COMMISSIONDENOMINATOR is a free data retrieval call binding the contract method 0x3b1dbfcc.
//
// Solidity: function COMMISSION_DENOMINATOR() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) COMMISSIONDENOMINATOR() (*big.Int, error) {
	return _IPublicDelegation.Contract.COMMISSIONDENOMINATOR(&_IPublicDelegation.CallOpts)
}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() pure returns(string)
func (_IPublicDelegation *IPublicDelegationCaller) CONTRACTTYPE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "CONTRACT_TYPE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() pure returns(string)
func (_IPublicDelegation *IPublicDelegationSession) CONTRACTTYPE() (string, error) {
	return _IPublicDelegation.Contract.CONTRACTTYPE(&_IPublicDelegation.CallOpts)
}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() pure returns(string)
func (_IPublicDelegation *IPublicDelegationCallerSession) CONTRACTTYPE() (string, error) {
	return _IPublicDelegation.Contract.CONTRACTTYPE(&_IPublicDelegation.CallOpts)
}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) MAXCOMMISSIONRATE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "MAX_COMMISSION_RATE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) MAXCOMMISSIONRATE() (*big.Int, error) {
	return _IPublicDelegation.Contract.MAXCOMMISSIONRATE(&_IPublicDelegation.CallOpts)
}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) MAXCOMMISSIONRATE() (*big.Int, error) {
	return _IPublicDelegation.Contract.MAXCOMMISSIONRATE(&_IPublicDelegation.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) VERSION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) VERSION() (*big.Int, error) {
	return _IPublicDelegation.Contract.VERSION(&_IPublicDelegation.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() pure returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) VERSION() (*big.Int, error) {
	return _IPublicDelegation.Contract.VERSION(&_IPublicDelegation.CallOpts)
}

// BaseCnStaking is a free data retrieval call binding the contract method 0x0c11487f.
//
// Solidity: function baseCnStaking() view returns(address)
func (_IPublicDelegation *IPublicDelegationCaller) BaseCnStaking(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "baseCnStaking")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BaseCnStaking is a free data retrieval call binding the contract method 0x0c11487f.
//
// Solidity: function baseCnStaking() view returns(address)
func (_IPublicDelegation *IPublicDelegationSession) BaseCnStaking() (common.Address, error) {
	return _IPublicDelegation.Contract.BaseCnStaking(&_IPublicDelegation.CallOpts)
}

// BaseCnStaking is a free data retrieval call binding the contract method 0x0c11487f.
//
// Solidity: function baseCnStaking() view returns(address)
func (_IPublicDelegation *IPublicDelegationCallerSession) BaseCnStaking() (common.Address, error) {
	return _IPublicDelegation.Contract.BaseCnStaking(&_IPublicDelegation.CallOpts)
}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) CommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "commissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) CommissionRate() (*big.Int, error) {
	return _IPublicDelegation.Contract.CommissionRate(&_IPublicDelegation.CallOpts)
}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) CommissionRate() (*big.Int, error) {
	return _IPublicDelegation.Contract.CommissionRate(&_IPublicDelegation.CallOpts)
}

// CommissionTo is a free data retrieval call binding the contract method 0x2f9ac83a.
//
// Solidity: function commissionTo() view returns(address)
func (_IPublicDelegation *IPublicDelegationCaller) CommissionTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "commissionTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CommissionTo is a free data retrieval call binding the contract method 0x2f9ac83a.
//
// Solidity: function commissionTo() view returns(address)
func (_IPublicDelegation *IPublicDelegationSession) CommissionTo() (common.Address, error) {
	return _IPublicDelegation.Contract.CommissionTo(&_IPublicDelegation.CallOpts)
}

// CommissionTo is a free data retrieval call binding the contract method 0x2f9ac83a.
//
// Solidity: function commissionTo() view returns(address)
func (_IPublicDelegation *IPublicDelegationCallerSession) CommissionTo() (common.Address, error) {
	return _IPublicDelegation.Contract.CommissionTo(&_IPublicDelegation.CallOpts)
}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 _shares) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) ConvertToAssets(opts *bind.CallOpts, _shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "convertToAssets", _shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 _shares) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) ConvertToAssets(_shares *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.ConvertToAssets(&_IPublicDelegation.CallOpts, _shares)
}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 _shares) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) ConvertToAssets(_shares *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.ConvertToAssets(&_IPublicDelegation.CallOpts, _shares)
}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) ConvertToShares(opts *bind.CallOpts, _assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "convertToShares", _assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) ConvertToShares(_assets *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.ConvertToShares(&_IPublicDelegation.CallOpts, _assets)
}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) ConvertToShares(_assets *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.ConvertToShares(&_IPublicDelegation.CallOpts, _assets)
}

// GetCurrentWithdrawalRequestState is a free data retrieval call binding the contract method 0x04ddc9d1.
//
// Solidity: function getCurrentWithdrawalRequestState(uint256 _requestId) view returns(uint8)
func (_IPublicDelegation *IPublicDelegationCaller) GetCurrentWithdrawalRequestState(opts *bind.CallOpts, _requestId *big.Int) (uint8, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "getCurrentWithdrawalRequestState", _requestId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetCurrentWithdrawalRequestState is a free data retrieval call binding the contract method 0x04ddc9d1.
//
// Solidity: function getCurrentWithdrawalRequestState(uint256 _requestId) view returns(uint8)
func (_IPublicDelegation *IPublicDelegationSession) GetCurrentWithdrawalRequestState(_requestId *big.Int) (uint8, error) {
	return _IPublicDelegation.Contract.GetCurrentWithdrawalRequestState(&_IPublicDelegation.CallOpts, _requestId)
}

// GetCurrentWithdrawalRequestState is a free data retrieval call binding the contract method 0x04ddc9d1.
//
// Solidity: function getCurrentWithdrawalRequestState(uint256 _requestId) view returns(uint8)
func (_IPublicDelegation *IPublicDelegationCallerSession) GetCurrentWithdrawalRequestState(_requestId *big.Int) (uint8, error) {
	return _IPublicDelegation.Contract.GetCurrentWithdrawalRequestState(&_IPublicDelegation.CallOpts, _requestId)
}

// GetUserRequestCount is a free data retrieval call binding the contract method 0xc166c458.
//
// Solidity: function getUserRequestCount(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) GetUserRequestCount(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "getUserRequestCount", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserRequestCount is a free data retrieval call binding the contract method 0xc166c458.
//
// Solidity: function getUserRequestCount(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) GetUserRequestCount(_owner common.Address) (*big.Int, error) {
	return _IPublicDelegation.Contract.GetUserRequestCount(&_IPublicDelegation.CallOpts, _owner)
}

// GetUserRequestCount is a free data retrieval call binding the contract method 0xc166c458.
//
// Solidity: function getUserRequestCount(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) GetUserRequestCount(_owner common.Address) (*big.Int, error) {
	return _IPublicDelegation.Contract.GetUserRequestCount(&_IPublicDelegation.CallOpts, _owner)
}

// GetUserRequestIds is a free data retrieval call binding the contract method 0x60df7c6c.
//
// Solidity: function getUserRequestIds(address _owner) view returns(uint256[])
func (_IPublicDelegation *IPublicDelegationCaller) GetUserRequestIds(opts *bind.CallOpts, _owner common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "getUserRequestIds", _owner)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUserRequestIds is a free data retrieval call binding the contract method 0x60df7c6c.
//
// Solidity: function getUserRequestIds(address _owner) view returns(uint256[])
func (_IPublicDelegation *IPublicDelegationSession) GetUserRequestIds(_owner common.Address) ([]*big.Int, error) {
	return _IPublicDelegation.Contract.GetUserRequestIds(&_IPublicDelegation.CallOpts, _owner)
}

// GetUserRequestIds is a free data retrieval call binding the contract method 0x60df7c6c.
//
// Solidity: function getUserRequestIds(address _owner) view returns(uint256[])
func (_IPublicDelegation *IPublicDelegationCallerSession) GetUserRequestIds(_owner common.Address) ([]*big.Int, error) {
	return _IPublicDelegation.Contract.GetUserRequestIds(&_IPublicDelegation.CallOpts, _owner)
}

// GetUserRequestIdsWithState is a free data retrieval call binding the contract method 0x93b89a84.
//
// Solidity: function getUserRequestIdsWithState(address _owner, uint8 _state) view returns(uint256[])
func (_IPublicDelegation *IPublicDelegationCaller) GetUserRequestIdsWithState(opts *bind.CallOpts, _owner common.Address, _state uint8) ([]*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "getUserRequestIdsWithState", _owner, _state)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUserRequestIdsWithState is a free data retrieval call binding the contract method 0x93b89a84.
//
// Solidity: function getUserRequestIdsWithState(address _owner, uint8 _state) view returns(uint256[])
func (_IPublicDelegation *IPublicDelegationSession) GetUserRequestIdsWithState(_owner common.Address, _state uint8) ([]*big.Int, error) {
	return _IPublicDelegation.Contract.GetUserRequestIdsWithState(&_IPublicDelegation.CallOpts, _owner, _state)
}

// GetUserRequestIdsWithState is a free data retrieval call binding the contract method 0x93b89a84.
//
// Solidity: function getUserRequestIdsWithState(address _owner, uint8 _state) view returns(uint256[])
func (_IPublicDelegation *IPublicDelegationCallerSession) GetUserRequestIdsWithState(_owner common.Address, _state uint8) ([]*big.Int, error) {
	return _IPublicDelegation.Contract.GetUserRequestIdsWithState(&_IPublicDelegation.CallOpts, _owner, _state)
}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) MaxRedeem(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "maxRedeem", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) MaxRedeem(_owner common.Address) (*big.Int, error) {
	return _IPublicDelegation.Contract.MaxRedeem(&_IPublicDelegation.CallOpts, _owner)
}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) MaxRedeem(_owner common.Address) (*big.Int, error) {
	return _IPublicDelegation.Contract.MaxRedeem(&_IPublicDelegation.CallOpts, _owner)
}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) MaxWithdraw(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "maxWithdraw", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) MaxWithdraw(_owner common.Address) (*big.Int, error) {
	return _IPublicDelegation.Contract.MaxWithdraw(&_IPublicDelegation.CallOpts, _owner)
}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address _owner) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) MaxWithdraw(_owner common.Address) (*big.Int, error) {
	return _IPublicDelegation.Contract.MaxWithdraw(&_IPublicDelegation.CallOpts, _owner)
}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) PreviewDeposit(opts *bind.CallOpts, _assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "previewDeposit", _assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) PreviewDeposit(_assets *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.PreviewDeposit(&_IPublicDelegation.CallOpts, _assets)
}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) PreviewDeposit(_assets *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.PreviewDeposit(&_IPublicDelegation.CallOpts, _assets)
}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 _shares) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) PreviewRedeem(opts *bind.CallOpts, _shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "previewRedeem", _shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 _shares) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) PreviewRedeem(_shares *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.PreviewRedeem(&_IPublicDelegation.CallOpts, _shares)
}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 _shares) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) PreviewRedeem(_shares *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.PreviewRedeem(&_IPublicDelegation.CallOpts, _shares)
}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) PreviewWithdraw(opts *bind.CallOpts, _assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "previewWithdraw", _assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) PreviewWithdraw(_assets *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.PreviewWithdraw(&_IPublicDelegation.CallOpts, _assets)
}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 _assets) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) PreviewWithdraw(_assets *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.PreviewWithdraw(&_IPublicDelegation.CallOpts, _assets)
}

// RequestIdToOwner is a free data retrieval call binding the contract method 0xf29177c3.
//
// Solidity: function requestIdToOwner(uint256 _requestId) view returns(address)
func (_IPublicDelegation *IPublicDelegationCaller) RequestIdToOwner(opts *bind.CallOpts, _requestId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "requestIdToOwner", _requestId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RequestIdToOwner is a free data retrieval call binding the contract method 0xf29177c3.
//
// Solidity: function requestIdToOwner(uint256 _requestId) view returns(address)
func (_IPublicDelegation *IPublicDelegationSession) RequestIdToOwner(_requestId *big.Int) (common.Address, error) {
	return _IPublicDelegation.Contract.RequestIdToOwner(&_IPublicDelegation.CallOpts, _requestId)
}

// RequestIdToOwner is a free data retrieval call binding the contract method 0xf29177c3.
//
// Solidity: function requestIdToOwner(uint256 _requestId) view returns(address)
func (_IPublicDelegation *IPublicDelegationCallerSession) RequestIdToOwner(_requestId *big.Int) (common.Address, error) {
	return _IPublicDelegation.Contract.RequestIdToOwner(&_IPublicDelegation.CallOpts, _requestId)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) Reward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "reward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) Reward() (*big.Int, error) {
	return _IPublicDelegation.Contract.Reward(&_IPublicDelegation.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) Reward() (*big.Int, error) {
	return _IPublicDelegation.Contract.Reward(&_IPublicDelegation.CallOpts)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) TotalAssets(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "totalAssets")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) TotalAssets() (*big.Int, error) {
	return _IPublicDelegation.Contract.TotalAssets(&_IPublicDelegation.CallOpts)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) TotalAssets() (*big.Int, error) {
	return _IPublicDelegation.Contract.TotalAssets(&_IPublicDelegation.CallOpts)
}

// UserRequestIds is a free data retrieval call binding the contract method 0x97feb23c.
//
// Solidity: function userRequestIds(address _owner, uint256 _index) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCaller) UserRequestIds(opts *bind.CallOpts, _owner common.Address, _index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _IPublicDelegation.contract.Call(opts, &out, "userRequestIds", _owner, _index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserRequestIds is a free data retrieval call binding the contract method 0x97feb23c.
//
// Solidity: function userRequestIds(address _owner, uint256 _index) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationSession) UserRequestIds(_owner common.Address, _index *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.UserRequestIds(&_IPublicDelegation.CallOpts, _owner, _index)
}

// UserRequestIds is a free data retrieval call binding the contract method 0x97feb23c.
//
// Solidity: function userRequestIds(address _owner, uint256 _index) view returns(uint256)
func (_IPublicDelegation *IPublicDelegationCallerSession) UserRequestIds(_owner common.Address, _index *big.Int) (*big.Int, error) {
	return _IPublicDelegation.Contract.UserRequestIds(&_IPublicDelegation.CallOpts, _owner, _index)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _requestId) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) CancelApprovedStakingWithdrawal(opts *bind.TransactOpts, _requestId *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "cancelApprovedStakingWithdrawal", _requestId)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _requestId) returns()
func (_IPublicDelegation *IPublicDelegationSession) CancelApprovedStakingWithdrawal(_requestId *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.CancelApprovedStakingWithdrawal(&_IPublicDelegation.TransactOpts, _requestId)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _requestId) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) CancelApprovedStakingWithdrawal(_requestId *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.CancelApprovedStakingWithdrawal(&_IPublicDelegation.TransactOpts, _requestId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 _requestId) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) Claim(opts *bind.TransactOpts, _requestId *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "claim", _requestId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 _requestId) returns()
func (_IPublicDelegation *IPublicDelegationSession) Claim(_requestId *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Claim(&_IPublicDelegation.TransactOpts, _requestId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 _requestId) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) Claim(_requestId *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Claim(&_IPublicDelegation.TransactOpts, _requestId)
}

// Initialize is a paid mutator transaction binding the contract method 0x26cf277a.
//
// Solidity: function initialize(address _baseCnStaking, (address,address,uint256,string) _args) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) Initialize(opts *bind.TransactOpts, _baseCnStaking common.Address, _args IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "initialize", _baseCnStaking, _args)
}

// Initialize is a paid mutator transaction binding the contract method 0x26cf277a.
//
// Solidity: function initialize(address _baseCnStaking, (address,address,uint256,string) _args) returns()
func (_IPublicDelegation *IPublicDelegationSession) Initialize(_baseCnStaking common.Address, _args IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Initialize(&_IPublicDelegation.TransactOpts, _baseCnStaking, _args)
}

// Initialize is a paid mutator transaction binding the contract method 0x26cf277a.
//
// Solidity: function initialize(address _baseCnStaking, (address,address,uint256,string) _args) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) Initialize(_baseCnStaking common.Address, _args IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Initialize(&_IPublicDelegation.TransactOpts, _baseCnStaking, _args)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address _recipient, uint256 _shares) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) Redeem(opts *bind.TransactOpts, _recipient common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "redeem", _recipient, _shares)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address _recipient, uint256 _shares) returns()
func (_IPublicDelegation *IPublicDelegationSession) Redeem(_recipient common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Redeem(&_IPublicDelegation.TransactOpts, _recipient, _shares)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address _recipient, uint256 _shares) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) Redeem(_recipient common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Redeem(&_IPublicDelegation.TransactOpts, _recipient, _shares)
}

// RedelegateByAssets is a paid mutator transaction binding the contract method 0xe659d7d7.
//
// Solidity: function redelegateByAssets(address _targetCnStaking, uint256 _assets) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) RedelegateByAssets(opts *bind.TransactOpts, _targetCnStaking common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "redelegateByAssets", _targetCnStaking, _assets)
}

// RedelegateByAssets is a paid mutator transaction binding the contract method 0xe659d7d7.
//
// Solidity: function redelegateByAssets(address _targetCnStaking, uint256 _assets) returns()
func (_IPublicDelegation *IPublicDelegationSession) RedelegateByAssets(_targetCnStaking common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.RedelegateByAssets(&_IPublicDelegation.TransactOpts, _targetCnStaking, _assets)
}

// RedelegateByAssets is a paid mutator transaction binding the contract method 0xe659d7d7.
//
// Solidity: function redelegateByAssets(address _targetCnStaking, uint256 _assets) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) RedelegateByAssets(_targetCnStaking common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.RedelegateByAssets(&_IPublicDelegation.TransactOpts, _targetCnStaking, _assets)
}

// RedelegateByShares is a paid mutator transaction binding the contract method 0xe15fc350.
//
// Solidity: function redelegateByShares(address _targetCnStaking, uint256 _shares) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) RedelegateByShares(opts *bind.TransactOpts, _targetCnStaking common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "redelegateByShares", _targetCnStaking, _shares)
}

// RedelegateByShares is a paid mutator transaction binding the contract method 0xe15fc350.
//
// Solidity: function redelegateByShares(address _targetCnStaking, uint256 _shares) returns()
func (_IPublicDelegation *IPublicDelegationSession) RedelegateByShares(_targetCnStaking common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.RedelegateByShares(&_IPublicDelegation.TransactOpts, _targetCnStaking, _shares)
}

// RedelegateByShares is a paid mutator transaction binding the contract method 0xe15fc350.
//
// Solidity: function redelegateByShares(address _targetCnStaking, uint256 _shares) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) RedelegateByShares(_targetCnStaking common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.RedelegateByShares(&_IPublicDelegation.TransactOpts, _targetCnStaking, _shares)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_IPublicDelegation *IPublicDelegationTransactor) Stake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "stake")
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_IPublicDelegation *IPublicDelegationSession) Stake() (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Stake(&_IPublicDelegation.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) Stake() (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Stake(&_IPublicDelegation.TransactOpts)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_IPublicDelegation *IPublicDelegationTransactor) StakeFor(opts *bind.TransactOpts, _recipient common.Address) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "stakeFor", _recipient)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_IPublicDelegation *IPublicDelegationSession) StakeFor(_recipient common.Address) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.StakeFor(&_IPublicDelegation.TransactOpts, _recipient)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) StakeFor(_recipient common.Address) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.StakeFor(&_IPublicDelegation.TransactOpts, _recipient)
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_IPublicDelegation *IPublicDelegationTransactor) Sweep(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "sweep")
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_IPublicDelegation *IPublicDelegationSession) Sweep() (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Sweep(&_IPublicDelegation.TransactOpts)
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) Sweep() (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Sweep(&_IPublicDelegation.TransactOpts)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0x00fa3d50.
//
// Solidity: function updateCommissionRate(uint256 _commissionRate) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) UpdateCommissionRate(opts *bind.TransactOpts, _commissionRate *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "updateCommissionRate", _commissionRate)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0x00fa3d50.
//
// Solidity: function updateCommissionRate(uint256 _commissionRate) returns()
func (_IPublicDelegation *IPublicDelegationSession) UpdateCommissionRate(_commissionRate *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.UpdateCommissionRate(&_IPublicDelegation.TransactOpts, _commissionRate)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0x00fa3d50.
//
// Solidity: function updateCommissionRate(uint256 _commissionRate) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) UpdateCommissionRate(_commissionRate *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.UpdateCommissionRate(&_IPublicDelegation.TransactOpts, _commissionRate)
}

// UpdateCommissionTo is a paid mutator transaction binding the contract method 0x052028d0.
//
// Solidity: function updateCommissionTo(address _commissionTo) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) UpdateCommissionTo(opts *bind.TransactOpts, _commissionTo common.Address) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "updateCommissionTo", _commissionTo)
}

// UpdateCommissionTo is a paid mutator transaction binding the contract method 0x052028d0.
//
// Solidity: function updateCommissionTo(address _commissionTo) returns()
func (_IPublicDelegation *IPublicDelegationSession) UpdateCommissionTo(_commissionTo common.Address) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.UpdateCommissionTo(&_IPublicDelegation.TransactOpts, _commissionTo)
}

// UpdateCommissionTo is a paid mutator transaction binding the contract method 0x052028d0.
//
// Solidity: function updateCommissionTo(address _commissionTo) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) UpdateCommissionTo(_commissionTo common.Address) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.UpdateCommissionTo(&_IPublicDelegation.TransactOpts, _commissionTo)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _assets) returns()
func (_IPublicDelegation *IPublicDelegationTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.contract.Transact(opts, "withdraw", _recipient, _assets)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _assets) returns()
func (_IPublicDelegation *IPublicDelegationSession) Withdraw(_recipient common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Withdraw(&_IPublicDelegation.TransactOpts, _recipient, _assets)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _assets) returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) Withdraw(_recipient common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Withdraw(&_IPublicDelegation.TransactOpts, _recipient, _assets)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IPublicDelegation *IPublicDelegationTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IPublicDelegation.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IPublicDelegation *IPublicDelegationSession) Receive() (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Receive(&_IPublicDelegation.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_IPublicDelegation *IPublicDelegationTransactorSession) Receive() (*types.Transaction, error) {
	return _IPublicDelegation.Contract.Receive(&_IPublicDelegation.TransactOpts)
}

// IPublicDelegationClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the IPublicDelegation contract.
type IPublicDelegationClaimedIterator struct {
	Event *IPublicDelegationClaimed // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationClaimed)
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
		it.Event = new(IPublicDelegationClaimed)
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
func (it *IPublicDelegationClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationClaimed represents a Claimed event raised by the IPublicDelegation contract.
type IPublicDelegationClaimed struct {
	User      common.Address
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0xd8138f8a3f377c5259ca548e70e4c2de94f129f5a11036a15b69513cba2b426a.
//
// Solidity: event Claimed(address indexed user, uint256 indexed requestId)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterClaimed(opts *bind.FilterOpts, user []common.Address, requestId []*big.Int) (*IPublicDelegationClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "Claimed", userRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationClaimedIterator{contract: _IPublicDelegation.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0xd8138f8a3f377c5259ca548e70e4c2de94f129f5a11036a15b69513cba2b426a.
//
// Solidity: event Claimed(address indexed user, uint256 indexed requestId)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *IPublicDelegationClaimed, user []common.Address, requestId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "Claimed", userRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationClaimed)
				if err := _IPublicDelegation.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0xd8138f8a3f377c5259ca548e70e4c2de94f129f5a11036a15b69513cba2b426a.
//
// Solidity: event Claimed(address indexed user, uint256 indexed requestId)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseClaimed(log types.Log) (*IPublicDelegationClaimed, error) {
	event := new(IPublicDelegationClaimed)
	if err := _IPublicDelegation.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationDeployPublicDelegationIterator is returned from FilterDeployPublicDelegation and is used to iterate over the raw logs and unpacked data for DeployPublicDelegation events raised by the IPublicDelegation contract.
type IPublicDelegationDeployPublicDelegationIterator struct {
	Event *IPublicDelegationDeployPublicDelegation // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationDeployPublicDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationDeployPublicDelegation)
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
		it.Event = new(IPublicDelegationDeployPublicDelegation)
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
func (it *IPublicDelegationDeployPublicDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationDeployPublicDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationDeployPublicDelegation represents a DeployPublicDelegation event raised by the IPublicDelegation contract.
type IPublicDelegationDeployPublicDelegation struct {
	ContractType  string
	BaseCnStaking common.Address
	PdArgs        IPublicDelegationPDConstructorArgs
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDeployPublicDelegation is a free log retrieval operation binding the contract event 0xae1ee3fdcc521a751b838d314164fcdbc7c3dcd4f0a910e153506bf5d89ca1ec.
//
// Solidity: event DeployPublicDelegation(string contractType, address indexed baseCnStaking, (address,address,uint256,string) pdArgs)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterDeployPublicDelegation(opts *bind.FilterOpts, baseCnStaking []common.Address) (*IPublicDelegationDeployPublicDelegationIterator, error) {

	var baseCnStakingRule []interface{}
	for _, baseCnStakingItem := range baseCnStaking {
		baseCnStakingRule = append(baseCnStakingRule, baseCnStakingItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "DeployPublicDelegation", baseCnStakingRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationDeployPublicDelegationIterator{contract: _IPublicDelegation.contract, event: "DeployPublicDelegation", logs: logs, sub: sub}, nil
}

// WatchDeployPublicDelegation is a free log subscription operation binding the contract event 0xae1ee3fdcc521a751b838d314164fcdbc7c3dcd4f0a910e153506bf5d89ca1ec.
//
// Solidity: event DeployPublicDelegation(string contractType, address indexed baseCnStaking, (address,address,uint256,string) pdArgs)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchDeployPublicDelegation(opts *bind.WatchOpts, sink chan<- *IPublicDelegationDeployPublicDelegation, baseCnStaking []common.Address) (event.Subscription, error) {

	var baseCnStakingRule []interface{}
	for _, baseCnStakingItem := range baseCnStaking {
		baseCnStakingRule = append(baseCnStakingRule, baseCnStakingItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "DeployPublicDelegation", baseCnStakingRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationDeployPublicDelegation)
				if err := _IPublicDelegation.contract.UnpackLog(event, "DeployPublicDelegation", log); err != nil {
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

// ParseDeployPublicDelegation is a log parse operation binding the contract event 0xae1ee3fdcc521a751b838d314164fcdbc7c3dcd4f0a910e153506bf5d89ca1ec.
//
// Solidity: event DeployPublicDelegation(string contractType, address indexed baseCnStaking, (address,address,uint256,string) pdArgs)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseDeployPublicDelegation(log types.Log) (*IPublicDelegationDeployPublicDelegation, error) {
	event := new(IPublicDelegationDeployPublicDelegation)
	if err := _IPublicDelegation.contract.UnpackLog(event, "DeployPublicDelegation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationRedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the IPublicDelegation contract.
type IPublicDelegationRedeemedIterator struct {
	Event *IPublicDelegationRedeemed // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationRedeemed)
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
		it.Event = new(IPublicDelegationRedeemed)
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
func (it *IPublicDelegationRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationRedeemed represents a Redeemed event raised by the IPublicDelegation contract.
type IPublicDelegationRedeemed struct {
	User      common.Address
	Recipient common.Address
	Assets    *big.Int
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea83480.
//
// Solidity: event Redeemed(address indexed user, address indexed recipient, uint256 assets, uint256 shares)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterRedeemed(opts *bind.FilterOpts, user []common.Address, recipient []common.Address) (*IPublicDelegationRedeemedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "Redeemed", userRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationRedeemedIterator{contract: _IPublicDelegation.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea83480.
//
// Solidity: event Redeemed(address indexed user, address indexed recipient, uint256 assets, uint256 shares)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *IPublicDelegationRedeemed, user []common.Address, recipient []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "Redeemed", userRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationRedeemed)
				if err := _IPublicDelegation.contract.UnpackLog(event, "Redeemed", log); err != nil {
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

// ParseRedeemed is a log parse operation binding the contract event 0x5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea83480.
//
// Solidity: event Redeemed(address indexed user, address indexed recipient, uint256 assets, uint256 shares)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseRedeemed(log types.Log) (*IPublicDelegationRedeemed, error) {
	event := new(IPublicDelegationRedeemed)
	if err := _IPublicDelegation.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationRedelegatedIterator is returned from FilterRedelegated and is used to iterate over the raw logs and unpacked data for Redelegated events raised by the IPublicDelegation contract.
type IPublicDelegationRedelegatedIterator struct {
	Event *IPublicDelegationRedelegated // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationRedelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationRedelegated)
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
		it.Event = new(IPublicDelegationRedelegated)
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
func (it *IPublicDelegationRedelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationRedelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationRedelegated represents a Redelegated event raised by the IPublicDelegation contract.
type IPublicDelegationRedelegated struct {
	User            common.Address
	TargetCnStaking common.Address
	Assets          *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRedelegated is a free log retrieval operation binding the contract event 0x78d93753d5a8009153a294711a82a3d1ba802938d66537c1b9142a053782d693.
//
// Solidity: event Redelegated(address indexed user, address indexed targetCnStaking, uint256 assets)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterRedelegated(opts *bind.FilterOpts, user []common.Address, targetCnStaking []common.Address) (*IPublicDelegationRedelegatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var targetCnStakingRule []interface{}
	for _, targetCnStakingItem := range targetCnStaking {
		targetCnStakingRule = append(targetCnStakingRule, targetCnStakingItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "Redelegated", userRule, targetCnStakingRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationRedelegatedIterator{contract: _IPublicDelegation.contract, event: "Redelegated", logs: logs, sub: sub}, nil
}

// WatchRedelegated is a free log subscription operation binding the contract event 0x78d93753d5a8009153a294711a82a3d1ba802938d66537c1b9142a053782d693.
//
// Solidity: event Redelegated(address indexed user, address indexed targetCnStaking, uint256 assets)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchRedelegated(opts *bind.WatchOpts, sink chan<- *IPublicDelegationRedelegated, user []common.Address, targetCnStaking []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var targetCnStakingRule []interface{}
	for _, targetCnStakingItem := range targetCnStaking {
		targetCnStakingRule = append(targetCnStakingRule, targetCnStakingItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "Redelegated", userRule, targetCnStakingRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationRedelegated)
				if err := _IPublicDelegation.contract.UnpackLog(event, "Redelegated", log); err != nil {
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

// ParseRedelegated is a log parse operation binding the contract event 0x78d93753d5a8009153a294711a82a3d1ba802938d66537c1b9142a053782d693.
//
// Solidity: event Redelegated(address indexed user, address indexed targetCnStaking, uint256 assets)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseRedelegated(log types.Log) (*IPublicDelegationRedelegated, error) {
	event := new(IPublicDelegationRedelegated)
	if err := _IPublicDelegation.contract.UnpackLog(event, "Redelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationRequestCancelWithdrawalIterator is returned from FilterRequestCancelWithdrawal and is used to iterate over the raw logs and unpacked data for RequestCancelWithdrawal events raised by the IPublicDelegation contract.
type IPublicDelegationRequestCancelWithdrawalIterator struct {
	Event *IPublicDelegationRequestCancelWithdrawal // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationRequestCancelWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationRequestCancelWithdrawal)
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
		it.Event = new(IPublicDelegationRequestCancelWithdrawal)
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
func (it *IPublicDelegationRequestCancelWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationRequestCancelWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationRequestCancelWithdrawal represents a RequestCancelWithdrawal event raised by the IPublicDelegation contract.
type IPublicDelegationRequestCancelWithdrawal struct {
	User      common.Address
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestCancelWithdrawal is a free log retrieval operation binding the contract event 0x853eba293fc3859716cf2e1a34b0c266f8265db181ec04748d0f25b8a19fc80e.
//
// Solidity: event RequestCancelWithdrawal(address indexed user, uint256 indexed requestId)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterRequestCancelWithdrawal(opts *bind.FilterOpts, user []common.Address, requestId []*big.Int) (*IPublicDelegationRequestCancelWithdrawalIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "RequestCancelWithdrawal", userRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationRequestCancelWithdrawalIterator{contract: _IPublicDelegation.contract, event: "RequestCancelWithdrawal", logs: logs, sub: sub}, nil
}

// WatchRequestCancelWithdrawal is a free log subscription operation binding the contract event 0x853eba293fc3859716cf2e1a34b0c266f8265db181ec04748d0f25b8a19fc80e.
//
// Solidity: event RequestCancelWithdrawal(address indexed user, uint256 indexed requestId)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchRequestCancelWithdrawal(opts *bind.WatchOpts, sink chan<- *IPublicDelegationRequestCancelWithdrawal, user []common.Address, requestId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "RequestCancelWithdrawal", userRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationRequestCancelWithdrawal)
				if err := _IPublicDelegation.contract.UnpackLog(event, "RequestCancelWithdrawal", log); err != nil {
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

// ParseRequestCancelWithdrawal is a log parse operation binding the contract event 0x853eba293fc3859716cf2e1a34b0c266f8265db181ec04748d0f25b8a19fc80e.
//
// Solidity: event RequestCancelWithdrawal(address indexed user, uint256 indexed requestId)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseRequestCancelWithdrawal(log types.Log) (*IPublicDelegationRequestCancelWithdrawal, error) {
	event := new(IPublicDelegationRequestCancelWithdrawal)
	if err := _IPublicDelegation.contract.UnpackLog(event, "RequestCancelWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationRequestWithdrawalIterator is returned from FilterRequestWithdrawal and is used to iterate over the raw logs and unpacked data for RequestWithdrawal events raised by the IPublicDelegation contract.
type IPublicDelegationRequestWithdrawalIterator struct {
	Event *IPublicDelegationRequestWithdrawal // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationRequestWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationRequestWithdrawal)
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
		it.Event = new(IPublicDelegationRequestWithdrawal)
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
func (it *IPublicDelegationRequestWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationRequestWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationRequestWithdrawal represents a RequestWithdrawal event raised by the IPublicDelegation contract.
type IPublicDelegationRequestWithdrawal struct {
	User      common.Address
	Recipient common.Address
	RequestId *big.Int
	Assets    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestWithdrawal is a free log retrieval operation binding the contract event 0xd71e6ec4eed83207b08d7ee4a0773c0ff8f8a1ab94b8ce85737fc0c5ea2b5f0c.
//
// Solidity: event RequestWithdrawal(address indexed user, address indexed recipient, uint256 indexed requestId, uint256 assets)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterRequestWithdrawal(opts *bind.FilterOpts, user []common.Address, recipient []common.Address, requestId []*big.Int) (*IPublicDelegationRequestWithdrawalIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "RequestWithdrawal", userRule, recipientRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationRequestWithdrawalIterator{contract: _IPublicDelegation.contract, event: "RequestWithdrawal", logs: logs, sub: sub}, nil
}

// WatchRequestWithdrawal is a free log subscription operation binding the contract event 0xd71e6ec4eed83207b08d7ee4a0773c0ff8f8a1ab94b8ce85737fc0c5ea2b5f0c.
//
// Solidity: event RequestWithdrawal(address indexed user, address indexed recipient, uint256 indexed requestId, uint256 assets)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchRequestWithdrawal(opts *bind.WatchOpts, sink chan<- *IPublicDelegationRequestWithdrawal, user []common.Address, recipient []common.Address, requestId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "RequestWithdrawal", userRule, recipientRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationRequestWithdrawal)
				if err := _IPublicDelegation.contract.UnpackLog(event, "RequestWithdrawal", log); err != nil {
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

// ParseRequestWithdrawal is a log parse operation binding the contract event 0xd71e6ec4eed83207b08d7ee4a0773c0ff8f8a1ab94b8ce85737fc0c5ea2b5f0c.
//
// Solidity: event RequestWithdrawal(address indexed user, address indexed recipient, uint256 indexed requestId, uint256 assets)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseRequestWithdrawal(log types.Log) (*IPublicDelegationRequestWithdrawal, error) {
	event := new(IPublicDelegationRequestWithdrawal)
	if err := _IPublicDelegation.contract.UnpackLog(event, "RequestWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationSendCommissionIterator is returned from FilterSendCommission and is used to iterate over the raw logs and unpacked data for SendCommission events raised by the IPublicDelegation contract.
type IPublicDelegationSendCommissionIterator struct {
	Event *IPublicDelegationSendCommission // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationSendCommissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationSendCommission)
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
		it.Event = new(IPublicDelegationSendCommission)
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
func (it *IPublicDelegationSendCommissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationSendCommissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationSendCommission represents a SendCommission event raised by the IPublicDelegation contract.
type IPublicDelegationSendCommission struct {
	CommissionTo common.Address
	Commission   *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSendCommission is a free log retrieval operation binding the contract event 0x6c3b15f4a619d5331e4708792bb87f858ffb7a08f1a87aabca7cd15e51e04a63.
//
// Solidity: event SendCommission(address indexed commissionTo, uint256 commission)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterSendCommission(opts *bind.FilterOpts, commissionTo []common.Address) (*IPublicDelegationSendCommissionIterator, error) {

	var commissionToRule []interface{}
	for _, commissionToItem := range commissionTo {
		commissionToRule = append(commissionToRule, commissionToItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "SendCommission", commissionToRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationSendCommissionIterator{contract: _IPublicDelegation.contract, event: "SendCommission", logs: logs, sub: sub}, nil
}

// WatchSendCommission is a free log subscription operation binding the contract event 0x6c3b15f4a619d5331e4708792bb87f858ffb7a08f1a87aabca7cd15e51e04a63.
//
// Solidity: event SendCommission(address indexed commissionTo, uint256 commission)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchSendCommission(opts *bind.WatchOpts, sink chan<- *IPublicDelegationSendCommission, commissionTo []common.Address) (event.Subscription, error) {

	var commissionToRule []interface{}
	for _, commissionToItem := range commissionTo {
		commissionToRule = append(commissionToRule, commissionToItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "SendCommission", commissionToRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationSendCommission)
				if err := _IPublicDelegation.contract.UnpackLog(event, "SendCommission", log); err != nil {
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

// ParseSendCommission is a log parse operation binding the contract event 0x6c3b15f4a619d5331e4708792bb87f858ffb7a08f1a87aabca7cd15e51e04a63.
//
// Solidity: event SendCommission(address indexed commissionTo, uint256 commission)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseSendCommission(log types.Log) (*IPublicDelegationSendCommission, error) {
	event := new(IPublicDelegationSendCommission)
	if err := _IPublicDelegation.contract.UnpackLog(event, "SendCommission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the IPublicDelegation contract.
type IPublicDelegationStakedIterator struct {
	Event *IPublicDelegationStaked // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationStaked)
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
		it.Event = new(IPublicDelegationStaked)
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
func (it *IPublicDelegationStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationStaked represents a Staked event raised by the IPublicDelegation contract.
type IPublicDelegationStaked struct {
	User   common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address indexed user, uint256 assets, uint256 shares)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address) (*IPublicDelegationStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "Staked", userRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationStakedIterator{contract: _IPublicDelegation.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address indexed user, uint256 assets, uint256 shares)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *IPublicDelegationStaked, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "Staked", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationStaked)
				if err := _IPublicDelegation.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address indexed user, uint256 assets, uint256 shares)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseStaked(log types.Log) (*IPublicDelegationStaked, error) {
	event := new(IPublicDelegationStaked)
	if err := _IPublicDelegation.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationUpdateCommissionRateIterator is returned from FilterUpdateCommissionRate and is used to iterate over the raw logs and unpacked data for UpdateCommissionRate events raised by the IPublicDelegation contract.
type IPublicDelegationUpdateCommissionRateIterator struct {
	Event *IPublicDelegationUpdateCommissionRate // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationUpdateCommissionRateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationUpdateCommissionRate)
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
		it.Event = new(IPublicDelegationUpdateCommissionRate)
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
func (it *IPublicDelegationUpdateCommissionRateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationUpdateCommissionRateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationUpdateCommissionRate represents a UpdateCommissionRate event raised by the IPublicDelegation contract.
type IPublicDelegationUpdateCommissionRate struct {
	PrevCommissionRate *big.Int
	CommissionRate     *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpdateCommissionRate is a free log retrieval operation binding the contract event 0x67fb2216d844c3553cf557bffa85f0fde0294999f808e61dcae1773d50d5e150.
//
// Solidity: event UpdateCommissionRate(uint256 indexed prevCommissionRate, uint256 indexed commissionRate)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterUpdateCommissionRate(opts *bind.FilterOpts, prevCommissionRate []*big.Int, commissionRate []*big.Int) (*IPublicDelegationUpdateCommissionRateIterator, error) {

	var prevCommissionRateRule []interface{}
	for _, prevCommissionRateItem := range prevCommissionRate {
		prevCommissionRateRule = append(prevCommissionRateRule, prevCommissionRateItem)
	}
	var commissionRateRule []interface{}
	for _, commissionRateItem := range commissionRate {
		commissionRateRule = append(commissionRateRule, commissionRateItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "UpdateCommissionRate", prevCommissionRateRule, commissionRateRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationUpdateCommissionRateIterator{contract: _IPublicDelegation.contract, event: "UpdateCommissionRate", logs: logs, sub: sub}, nil
}

// WatchUpdateCommissionRate is a free log subscription operation binding the contract event 0x67fb2216d844c3553cf557bffa85f0fde0294999f808e61dcae1773d50d5e150.
//
// Solidity: event UpdateCommissionRate(uint256 indexed prevCommissionRate, uint256 indexed commissionRate)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchUpdateCommissionRate(opts *bind.WatchOpts, sink chan<- *IPublicDelegationUpdateCommissionRate, prevCommissionRate []*big.Int, commissionRate []*big.Int) (event.Subscription, error) {

	var prevCommissionRateRule []interface{}
	for _, prevCommissionRateItem := range prevCommissionRate {
		prevCommissionRateRule = append(prevCommissionRateRule, prevCommissionRateItem)
	}
	var commissionRateRule []interface{}
	for _, commissionRateItem := range commissionRate {
		commissionRateRule = append(commissionRateRule, commissionRateItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "UpdateCommissionRate", prevCommissionRateRule, commissionRateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationUpdateCommissionRate)
				if err := _IPublicDelegation.contract.UnpackLog(event, "UpdateCommissionRate", log); err != nil {
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

// ParseUpdateCommissionRate is a log parse operation binding the contract event 0x67fb2216d844c3553cf557bffa85f0fde0294999f808e61dcae1773d50d5e150.
//
// Solidity: event UpdateCommissionRate(uint256 indexed prevCommissionRate, uint256 indexed commissionRate)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseUpdateCommissionRate(log types.Log) (*IPublicDelegationUpdateCommissionRate, error) {
	event := new(IPublicDelegationUpdateCommissionRate)
	if err := _IPublicDelegation.contract.UnpackLog(event, "UpdateCommissionRate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IPublicDelegationUpdateCommissionToIterator is returned from FilterUpdateCommissionTo and is used to iterate over the raw logs and unpacked data for UpdateCommissionTo events raised by the IPublicDelegation contract.
type IPublicDelegationUpdateCommissionToIterator struct {
	Event *IPublicDelegationUpdateCommissionTo // Event containing the contract specifics and raw log

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
func (it *IPublicDelegationUpdateCommissionToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IPublicDelegationUpdateCommissionTo)
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
		it.Event = new(IPublicDelegationUpdateCommissionTo)
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
func (it *IPublicDelegationUpdateCommissionToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IPublicDelegationUpdateCommissionToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IPublicDelegationUpdateCommissionTo represents a UpdateCommissionTo event raised by the IPublicDelegation contract.
type IPublicDelegationUpdateCommissionTo struct {
	PrevCommissionTo common.Address
	CommissionTo     common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpdateCommissionTo is a free log retrieval operation binding the contract event 0xa466bc069316af563b6a90c9b6f29a97752ac7be7c6bbf3767bc9b16c2fd90eb.
//
// Solidity: event UpdateCommissionTo(address indexed prevCommissionTo, address indexed commissionTo)
func (_IPublicDelegation *IPublicDelegationFilterer) FilterUpdateCommissionTo(opts *bind.FilterOpts, prevCommissionTo []common.Address, commissionTo []common.Address) (*IPublicDelegationUpdateCommissionToIterator, error) {

	var prevCommissionToRule []interface{}
	for _, prevCommissionToItem := range prevCommissionTo {
		prevCommissionToRule = append(prevCommissionToRule, prevCommissionToItem)
	}
	var commissionToRule []interface{}
	for _, commissionToItem := range commissionTo {
		commissionToRule = append(commissionToRule, commissionToItem)
	}

	logs, sub, err := _IPublicDelegation.contract.FilterLogs(opts, "UpdateCommissionTo", prevCommissionToRule, commissionToRule)
	if err != nil {
		return nil, err
	}
	return &IPublicDelegationUpdateCommissionToIterator{contract: _IPublicDelegation.contract, event: "UpdateCommissionTo", logs: logs, sub: sub}, nil
}

// WatchUpdateCommissionTo is a free log subscription operation binding the contract event 0xa466bc069316af563b6a90c9b6f29a97752ac7be7c6bbf3767bc9b16c2fd90eb.
//
// Solidity: event UpdateCommissionTo(address indexed prevCommissionTo, address indexed commissionTo)
func (_IPublicDelegation *IPublicDelegationFilterer) WatchUpdateCommissionTo(opts *bind.WatchOpts, sink chan<- *IPublicDelegationUpdateCommissionTo, prevCommissionTo []common.Address, commissionTo []common.Address) (event.Subscription, error) {

	var prevCommissionToRule []interface{}
	for _, prevCommissionToItem := range prevCommissionTo {
		prevCommissionToRule = append(prevCommissionToRule, prevCommissionToItem)
	}
	var commissionToRule []interface{}
	for _, commissionToItem := range commissionTo {
		commissionToRule = append(commissionToRule, commissionToItem)
	}

	logs, sub, err := _IPublicDelegation.contract.WatchLogs(opts, "UpdateCommissionTo", prevCommissionToRule, commissionToRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IPublicDelegationUpdateCommissionTo)
				if err := _IPublicDelegation.contract.UnpackLog(event, "UpdateCommissionTo", log); err != nil {
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

// ParseUpdateCommissionTo is a log parse operation binding the contract event 0xa466bc069316af563b6a90c9b6f29a97752ac7be7c6bbf3767bc9b16c2fd90eb.
//
// Solidity: event UpdateCommissionTo(address indexed prevCommissionTo, address indexed commissionTo)
func (_IPublicDelegation *IPublicDelegationFilterer) ParseUpdateCommissionTo(log types.Log) (*IPublicDelegationUpdateCommissionTo, error) {
	event := new(IPublicDelegationUpdateCommissionTo)
	if err := _IPublicDelegation.contract.UnpackLog(event, "UpdateCommissionTo", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// InitializableMetaData contains all meta data concerning the Initializable contract.
var InitializableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"}]",
}

// InitializableABI is the input ABI used to generate the binding from.
// Deprecated: Use InitializableMetaData.ABI instead.
var InitializableABI = InitializableMetaData.ABI

// InitializableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const InitializableBinRuntime = ``

// Initializable is an auto generated Go binding around a Kaia contract.
type Initializable struct {
	InitializableCaller     // Read-only binding to the contract
	InitializableTransactor // Write-only binding to the contract
	InitializableFilterer   // Log filterer for contract events
}

// InitializableCaller is an auto generated read-only Go binding around a Kaia contract.
type InitializableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitializableTransactor is an auto generated write-only Go binding around a Kaia contract.
type InitializableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitializableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type InitializableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InitializableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type InitializableSession struct {
	Contract     *Initializable    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// InitializableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type InitializableCallerSession struct {
	Contract *InitializableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// InitializableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type InitializableTransactorSession struct {
	Contract     *InitializableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// InitializableRaw is an auto generated low-level Go binding around a Kaia contract.
type InitializableRaw struct {
	Contract *Initializable // Generic contract binding to access the raw methods on
}

// InitializableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type InitializableCallerRaw struct {
	Contract *InitializableCaller // Generic read-only contract binding to access the raw methods on
}

// InitializableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type InitializableTransactorRaw struct {
	Contract *InitializableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInitializable creates a new instance of Initializable, bound to a specific deployed contract.
func NewInitializable(address common.Address, backend bind.ContractBackend) (*Initializable, error) {
	contract, err := bindInitializable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Initializable{InitializableCaller: InitializableCaller{contract: contract}, InitializableTransactor: InitializableTransactor{contract: contract}, InitializableFilterer: InitializableFilterer{contract: contract}}, nil
}

// NewInitializableCaller creates a new read-only instance of Initializable, bound to a specific deployed contract.
func NewInitializableCaller(address common.Address, caller bind.ContractCaller) (*InitializableCaller, error) {
	contract, err := bindInitializable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InitializableCaller{contract: contract}, nil
}

// NewInitializableTransactor creates a new write-only instance of Initializable, bound to a specific deployed contract.
func NewInitializableTransactor(address common.Address, transactor bind.ContractTransactor) (*InitializableTransactor, error) {
	contract, err := bindInitializable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InitializableTransactor{contract: contract}, nil
}

// NewInitializableFilterer creates a new log filterer instance of Initializable, bound to a specific deployed contract.
func NewInitializableFilterer(address common.Address, filterer bind.ContractFilterer) (*InitializableFilterer, error) {
	contract, err := bindInitializable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InitializableFilterer{contract: contract}, nil
}

// bindInitializable binds a generic wrapper to an already deployed contract.
func bindInitializable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InitializableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Initializable *InitializableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Initializable.Contract.InitializableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Initializable *InitializableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Initializable.Contract.InitializableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Initializable *InitializableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Initializable.Contract.InitializableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Initializable *InitializableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Initializable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Initializable *InitializableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Initializable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Initializable *InitializableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Initializable.Contract.contract.Transact(opts, method, params...)
}

// InitializableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Initializable contract.
type InitializableInitializedIterator struct {
	Event *InitializableInitialized // Event containing the contract specifics and raw log

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
func (it *InitializableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(InitializableInitialized)
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
		it.Event = new(InitializableInitialized)
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
func (it *InitializableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *InitializableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// InitializableInitialized represents a Initialized event raised by the Initializable contract.
type InitializableInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Initializable *InitializableFilterer) FilterInitialized(opts *bind.FilterOpts) (*InitializableInitializedIterator, error) {

	logs, sub, err := _Initializable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &InitializableInitializedIterator{contract: _Initializable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Initializable *InitializableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *InitializableInitialized) (event.Subscription, error) {

	logs, sub, err := _Initializable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(InitializableInitialized)
				if err := _Initializable.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Initializable *InitializableFilterer) ParseInitialized(log types.Log) (*InitializableInitialized, error) {
	event := new(InitializableInitialized)
	if err := _Initializable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LowLevelCallMetaData contains all meta data concerning the LowLevelCall contract.
var LowLevelCallMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220a9c01c846acbeb4f949683453874c8f22300476e0b9b9d6e7875a98789c479a464736f6c63430008190033",
}

// LowLevelCallABI is the input ABI used to generate the binding from.
// Deprecated: Use LowLevelCallMetaData.ABI instead.
var LowLevelCallABI = LowLevelCallMetaData.ABI

// LowLevelCallBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const LowLevelCallBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220a9c01c846acbeb4f949683453874c8f22300476e0b9b9d6e7875a98789c479a464736f6c63430008190033`

// LowLevelCallBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LowLevelCallMetaData.Bin instead.
var LowLevelCallBin = LowLevelCallMetaData.Bin

// DeployLowLevelCall deploys a new Kaia contract, binding an instance of LowLevelCall to it.
func DeployLowLevelCall(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *LowLevelCall, error) {
	parsed, err := LowLevelCallMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LowLevelCallBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LowLevelCall{LowLevelCallCaller: LowLevelCallCaller{contract: contract}, LowLevelCallTransactor: LowLevelCallTransactor{contract: contract}, LowLevelCallFilterer: LowLevelCallFilterer{contract: contract}}, nil
}

// LowLevelCall is an auto generated Go binding around a Kaia contract.
type LowLevelCall struct {
	LowLevelCallCaller     // Read-only binding to the contract
	LowLevelCallTransactor // Write-only binding to the contract
	LowLevelCallFilterer   // Log filterer for contract events
}

// LowLevelCallCaller is an auto generated read-only Go binding around a Kaia contract.
type LowLevelCallCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LowLevelCallTransactor is an auto generated write-only Go binding around a Kaia contract.
type LowLevelCallTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LowLevelCallFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type LowLevelCallFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LowLevelCallSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type LowLevelCallSession struct {
	Contract     *LowLevelCall     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LowLevelCallCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type LowLevelCallCallerSession struct {
	Contract *LowLevelCallCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// LowLevelCallTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type LowLevelCallTransactorSession struct {
	Contract     *LowLevelCallTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// LowLevelCallRaw is an auto generated low-level Go binding around a Kaia contract.
type LowLevelCallRaw struct {
	Contract *LowLevelCall // Generic contract binding to access the raw methods on
}

// LowLevelCallCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type LowLevelCallCallerRaw struct {
	Contract *LowLevelCallCaller // Generic read-only contract binding to access the raw methods on
}

// LowLevelCallTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type LowLevelCallTransactorRaw struct {
	Contract *LowLevelCallTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLowLevelCall creates a new instance of LowLevelCall, bound to a specific deployed contract.
func NewLowLevelCall(address common.Address, backend bind.ContractBackend) (*LowLevelCall, error) {
	contract, err := bindLowLevelCall(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LowLevelCall{LowLevelCallCaller: LowLevelCallCaller{contract: contract}, LowLevelCallTransactor: LowLevelCallTransactor{contract: contract}, LowLevelCallFilterer: LowLevelCallFilterer{contract: contract}}, nil
}

// NewLowLevelCallCaller creates a new read-only instance of LowLevelCall, bound to a specific deployed contract.
func NewLowLevelCallCaller(address common.Address, caller bind.ContractCaller) (*LowLevelCallCaller, error) {
	contract, err := bindLowLevelCall(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LowLevelCallCaller{contract: contract}, nil
}

// NewLowLevelCallTransactor creates a new write-only instance of LowLevelCall, bound to a specific deployed contract.
func NewLowLevelCallTransactor(address common.Address, transactor bind.ContractTransactor) (*LowLevelCallTransactor, error) {
	contract, err := bindLowLevelCall(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LowLevelCallTransactor{contract: contract}, nil
}

// NewLowLevelCallFilterer creates a new log filterer instance of LowLevelCall, bound to a specific deployed contract.
func NewLowLevelCallFilterer(address common.Address, filterer bind.ContractFilterer) (*LowLevelCallFilterer, error) {
	contract, err := bindLowLevelCall(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LowLevelCallFilterer{contract: contract}, nil
}

// bindLowLevelCall binds a generic wrapper to an already deployed contract.
func bindLowLevelCall(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LowLevelCallMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LowLevelCall *LowLevelCallRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LowLevelCall.Contract.LowLevelCallCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LowLevelCall *LowLevelCallRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LowLevelCall.Contract.LowLevelCallTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LowLevelCall *LowLevelCallRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LowLevelCall.Contract.LowLevelCallTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LowLevelCall *LowLevelCallCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LowLevelCall.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LowLevelCall *LowLevelCallTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LowLevelCall.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LowLevelCall *LowLevelCallTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LowLevelCall.Contract.contract.Transact(opts, method, params...)
}

// MathMetaData contains all meta data concerning the Math contract.
var MathMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea26469706673582212200d6b6b8f3aef955b67b64e96216116caf99bc52accccffdcda9b915871ae129b64736f6c63430008190033",
}

// MathABI is the input ABI used to generate the binding from.
// Deprecated: Use MathMetaData.ABI instead.
var MathABI = MathMetaData.ABI

// MathBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const MathBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea26469706673582212200d6b6b8f3aef955b67b64e96216116caf99bc52accccffdcda9b915871ae129b64736f6c63430008190033`

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

// OwnableUpgradeableMetaData contains all meta data concerning the OwnableUpgradeable contract.
var OwnableUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"8da5cb5b": "owner()",
		"715018a6": "renounceOwnership()",
		"f2fde38b": "transferOwnership(address)",
	},
}

// OwnableUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableUpgradeableMetaData.ABI instead.
var OwnableUpgradeableABI = OwnableUpgradeableMetaData.ABI

// OwnableUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const OwnableUpgradeableBinRuntime = ``

// Deprecated: Use OwnableUpgradeableMetaData.Sigs instead.
// OwnableUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
var OwnableUpgradeableFuncSigs = OwnableUpgradeableMetaData.Sigs

// OwnableUpgradeable is an auto generated Go binding around a Kaia contract.
type OwnableUpgradeable struct {
	OwnableUpgradeableCaller     // Read-only binding to the contract
	OwnableUpgradeableTransactor // Write-only binding to the contract
	OwnableUpgradeableFilterer   // Log filterer for contract events
}

// OwnableUpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type OwnableUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableUpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type OwnableUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableUpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type OwnableUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableUpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type OwnableUpgradeableSession struct {
	Contract     *OwnableUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OwnableUpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type OwnableUpgradeableCallerSession struct {
	Contract *OwnableUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// OwnableUpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type OwnableUpgradeableTransactorSession struct {
	Contract     *OwnableUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OwnableUpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type OwnableUpgradeableRaw struct {
	Contract *OwnableUpgradeable // Generic contract binding to access the raw methods on
}

// OwnableUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type OwnableUpgradeableCallerRaw struct {
	Contract *OwnableUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type OwnableUpgradeableTransactorRaw struct {
	Contract *OwnableUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnableUpgradeable creates a new instance of OwnableUpgradeable, bound to a specific deployed contract.
func NewOwnableUpgradeable(address common.Address, backend bind.ContractBackend) (*OwnableUpgradeable, error) {
	contract, err := bindOwnableUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeable{OwnableUpgradeableCaller: OwnableUpgradeableCaller{contract: contract}, OwnableUpgradeableTransactor: OwnableUpgradeableTransactor{contract: contract}, OwnableUpgradeableFilterer: OwnableUpgradeableFilterer{contract: contract}}, nil
}

// NewOwnableUpgradeableCaller creates a new read-only instance of OwnableUpgradeable, bound to a specific deployed contract.
func NewOwnableUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*OwnableUpgradeableCaller, error) {
	contract, err := bindOwnableUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableCaller{contract: contract}, nil
}

// NewOwnableUpgradeableTransactor creates a new write-only instance of OwnableUpgradeable, bound to a specific deployed contract.
func NewOwnableUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableUpgradeableTransactor, error) {
	contract, err := bindOwnableUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableTransactor{contract: contract}, nil
}

// NewOwnableUpgradeableFilterer creates a new log filterer instance of OwnableUpgradeable, bound to a specific deployed contract.
func NewOwnableUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableUpgradeableFilterer, error) {
	contract, err := bindOwnableUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableFilterer{contract: contract}, nil
}

// bindOwnableUpgradeable binds a generic wrapper to an already deployed contract.
func bindOwnableUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OwnableUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableUpgradeable *OwnableUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableUpgradeable.Contract.OwnableUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableUpgradeable *OwnableUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.OwnableUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableUpgradeable *OwnableUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.OwnableUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OwnableUpgradeable *OwnableUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OwnableUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OwnableUpgradeable *OwnableUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OwnableUpgradeable *OwnableUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableUpgradeable *OwnableUpgradeableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _OwnableUpgradeable.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableUpgradeable *OwnableUpgradeableSession) Owner() (common.Address, error) {
	return _OwnableUpgradeable.Contract.Owner(&_OwnableUpgradeable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_OwnableUpgradeable *OwnableUpgradeableCallerSession) Owner() (common.Address, error) {
	return _OwnableUpgradeable.Contract.Owner(&_OwnableUpgradeable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableUpgradeable *OwnableUpgradeableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OwnableUpgradeable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableUpgradeable *OwnableUpgradeableSession) RenounceOwnership() (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.RenounceOwnership(&_OwnableUpgradeable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_OwnableUpgradeable *OwnableUpgradeableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.RenounceOwnership(&_OwnableUpgradeable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableUpgradeable *OwnableUpgradeableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _OwnableUpgradeable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableUpgradeable *OwnableUpgradeableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.TransferOwnership(&_OwnableUpgradeable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_OwnableUpgradeable *OwnableUpgradeableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _OwnableUpgradeable.Contract.TransferOwnership(&_OwnableUpgradeable.TransactOpts, newOwner)
}

// OwnableUpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the OwnableUpgradeable contract.
type OwnableUpgradeableInitializedIterator struct {
	Event *OwnableUpgradeableInitialized // Event containing the contract specifics and raw log

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
func (it *OwnableUpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableUpgradeableInitialized)
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
		it.Event = new(OwnableUpgradeableInitialized)
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
func (it *OwnableUpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableUpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableUpgradeableInitialized represents a Initialized event raised by the OwnableUpgradeable contract.
type OwnableUpgradeableInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*OwnableUpgradeableInitializedIterator, error) {

	logs, sub, err := _OwnableUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableInitializedIterator{contract: _OwnableUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *OwnableUpgradeableInitialized) (event.Subscription, error) {

	logs, sub, err := _OwnableUpgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableUpgradeableInitialized)
				if err := _OwnableUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) ParseInitialized(log types.Log) (*OwnableUpgradeableInitialized, error) {
	event := new(OwnableUpgradeableInitialized)
	if err := _OwnableUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnableUpgradeableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OwnableUpgradeable contract.
type OwnableUpgradeableOwnershipTransferredIterator struct {
	Event *OwnableUpgradeableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OwnableUpgradeableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableUpgradeableOwnershipTransferred)
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
		it.Event = new(OwnableUpgradeableOwnershipTransferred)
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
func (it *OwnableUpgradeableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableUpgradeableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableUpgradeableOwnershipTransferred represents a OwnershipTransferred event raised by the OwnableUpgradeable contract.
type OwnableUpgradeableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableUpgradeableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnableUpgradeable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableOwnershipTransferredIterator{contract: _OwnableUpgradeable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableUpgradeableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _OwnableUpgradeable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableUpgradeableOwnershipTransferred)
				if err := _OwnableUpgradeable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) ParseOwnershipTransferred(log types.Log) (*OwnableUpgradeableOwnershipTransferred, error) {
	event := new(OwnableUpgradeableOwnershipTransferred)
	if err := _OwnableUpgradeable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PanicMetaData contains all meta data concerning the Panic contract.
var PanicMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220061e90e634f79a5e2f0f24b3e4fbf25eadda2c47def87cfa4a1cebcb86c1526864736f6c63430008190033",
}

// PanicABI is the input ABI used to generate the binding from.
// Deprecated: Use PanicMetaData.ABI instead.
var PanicABI = PanicMetaData.ABI

// PanicBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const PanicBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220061e90e634f79a5e2f0f24b3e4fbf25eadda2c47def87cfa4a1cebcb86c1526864736f6c63430008190033`

// PanicBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PanicMetaData.Bin instead.
var PanicBin = PanicMetaData.Bin

// DeployPanic deploys a new Kaia contract, binding an instance of Panic to it.
func DeployPanic(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Panic, error) {
	parsed, err := PanicMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PanicBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Panic{PanicCaller: PanicCaller{contract: contract}, PanicTransactor: PanicTransactor{contract: contract}, PanicFilterer: PanicFilterer{contract: contract}}, nil
}

// Panic is an auto generated Go binding around a Kaia contract.
type Panic struct {
	PanicCaller     // Read-only binding to the contract
	PanicTransactor // Write-only binding to the contract
	PanicFilterer   // Log filterer for contract events
}

// PanicCaller is an auto generated read-only Go binding around a Kaia contract.
type PanicCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PanicTransactor is an auto generated write-only Go binding around a Kaia contract.
type PanicTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PanicFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type PanicFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PanicSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type PanicSession struct {
	Contract     *Panic            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PanicCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type PanicCallerSession struct {
	Contract *PanicCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PanicTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type PanicTransactorSession struct {
	Contract     *PanicTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PanicRaw is an auto generated low-level Go binding around a Kaia contract.
type PanicRaw struct {
	Contract *Panic // Generic contract binding to access the raw methods on
}

// PanicCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type PanicCallerRaw struct {
	Contract *PanicCaller // Generic read-only contract binding to access the raw methods on
}

// PanicTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type PanicTransactorRaw struct {
	Contract *PanicTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPanic creates a new instance of Panic, bound to a specific deployed contract.
func NewPanic(address common.Address, backend bind.ContractBackend) (*Panic, error) {
	contract, err := bindPanic(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Panic{PanicCaller: PanicCaller{contract: contract}, PanicTransactor: PanicTransactor{contract: contract}, PanicFilterer: PanicFilterer{contract: contract}}, nil
}

// NewPanicCaller creates a new read-only instance of Panic, bound to a specific deployed contract.
func NewPanicCaller(address common.Address, caller bind.ContractCaller) (*PanicCaller, error) {
	contract, err := bindPanic(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PanicCaller{contract: contract}, nil
}

// NewPanicTransactor creates a new write-only instance of Panic, bound to a specific deployed contract.
func NewPanicTransactor(address common.Address, transactor bind.ContractTransactor) (*PanicTransactor, error) {
	contract, err := bindPanic(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PanicTransactor{contract: contract}, nil
}

// NewPanicFilterer creates a new log filterer instance of Panic, bound to a specific deployed contract.
func NewPanicFilterer(address common.Address, filterer bind.ContractFilterer) (*PanicFilterer, error) {
	contract, err := bindPanic(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PanicFilterer{contract: contract}, nil
}

// bindPanic binds a generic wrapper to an already deployed contract.
func bindPanic(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PanicMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Panic *PanicRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Panic.Contract.PanicCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Panic *PanicRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Panic.Contract.PanicTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Panic *PanicRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Panic.Contract.PanicTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Panic *PanicCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Panic.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Panic *PanicTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Panic.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Panic *PanicTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Panic.Contract.contract.Transact(opts, method, params...)
}

// PublicDelegationMetaData contains all meta data concerning the PublicDelegation contract.
var PublicDelegationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"CommissionRateTooHigh\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotRequestOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RedelegateAmountTooLow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakeAmountTooLow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WithdrawalAmountTooLow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"Claimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"contractType\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"baseCnStaking\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gcName\",\"type\":\"string\"}],\"indexed\":false,\"internalType\":\"structIPublicDelegation.PDConstructorArgs\",\"name\":\"pdArgs\",\"type\":\"tuple\"}],\"name\":\"DeployPublicDelegation\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"targetCnStaking\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"}],\"name\":\"Redelegated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"}],\"name\":\"RequestCancelWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"requestId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"}],\"name\":\"RequestWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"commission\",\"type\":\"uint256\"}],\"name\":\"SendCommission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assets\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"shares\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"prevCommissionRate\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"}],\"name\":\"UpdateCommissionRate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"prevCommissionTo\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"}],\"name\":\"UpdateCommissionTo\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"COMMISSION_DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"CONTRACT_TYPE\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_COMMISSION_RATE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"VERSION\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"baseCnStaking\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"cancelApprovedStakingWithdrawal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commissionRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"commissionTo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"convertToAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"convertToShares\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"getCurrentWithdrawalRequestState\",\"outputs\":[{\"internalType\":\"enumIPublicDelegation.WithdrawalRequestState\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"getUserRequestCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"getUserRequestIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"enumIPublicDelegation.WithdrawalRequestState\",\"name\":\"_state\",\"type\":\"uint8\"}],\"name\":\"getUserRequestIdsWithState\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_baseCnStaking\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gcName\",\"type\":\"string\"}],\"internalType\":\"structIPublicDelegation.PDConstructorArgs\",\"name\":\"_args\",\"type\":\"tuple\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"maxRedeem\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"maxWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"previewDeposit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"previewRedeem\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"previewWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_targetCnStaking\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"redelegateByAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_targetCnStaking\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_shares\",\"type\":\"uint256\"}],\"name\":\"redelegateByShares\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requestId\",\"type\":\"uint256\"}],\"name\":\"requestIdToOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"stakeFor\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sweep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalAssets\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_commissionRate\",\"type\":\"uint256\"}],\"name\":\"updateCommissionRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_commissionTo\",\"type\":\"address\"}],\"name\":\"updateCommissionTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"userRequestIds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_assets\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Sigs: map[string]string{
		"3b1dbfcc": "COMMISSION_DENOMINATOR()",
		"4b6a94cc": "CONTRACT_TYPE()",
		"207239c0": "MAX_COMMISSION_RATE()",
		"ffa1ad74": "VERSION()",
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"0c11487f": "baseCnStaking()",
		"c804b115": "cancelApprovedStakingWithdrawal(uint256)",
		"379607f5": "claim(uint256)",
		"5ea1d6f8": "commissionRate()",
		"2f9ac83a": "commissionTo()",
		"07a2d13a": "convertToAssets(uint256)",
		"c6e6f592": "convertToShares(uint256)",
		"313ce567": "decimals()",
		"04ddc9d1": "getCurrentWithdrawalRequestState(uint256)",
		"c166c458": "getUserRequestCount(address)",
		"60df7c6c": "getUserRequestIds(address)",
		"93b89a84": "getUserRequestIdsWithState(address,uint8)",
		"26cf277a": "initialize(address,(address,address,uint256,string))",
		"d905777e": "maxRedeem(address)",
		"ce96cb77": "maxWithdraw(address)",
		"06fdde03": "name()",
		"8da5cb5b": "owner()",
		"ef8b30f7": "previewDeposit(uint256)",
		"4cdad506": "previewRedeem(uint256)",
		"0a28a477": "previewWithdraw(uint256)",
		"1e9a6950": "redeem(address,uint256)",
		"e659d7d7": "redelegateByAssets(address,uint256)",
		"e15fc350": "redelegateByShares(address,uint256)",
		"715018a6": "renounceOwnership()",
		"f29177c3": "requestIdToOwner(uint256)",
		"228cb733": "reward()",
		"3a4b66f1": "stake()",
		"4bf69206": "stakeFor(address)",
		"35faa416": "sweep()",
		"95d89b41": "symbol()",
		"01e1d114": "totalAssets()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"f2fde38b": "transferOwnership(address)",
		"00fa3d50": "updateCommissionRate(uint256)",
		"052028d0": "updateCommissionTo(address)",
		"97feb23c": "userRequestIds(address,uint256)",
		"f3fef3a3": "withdraw(address,uint256)",
	},
	Bin: "0x6080604052348015600e575f80fd5b5060156019565b60c9565b7ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00805468010000000000000000900460ff161560685760405163f92ee8a960e01b815260040160405180910390fd5b80546001600160401b039081161460c65780546001600160401b0319166001600160401b0390811782556040519081527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a15b50565b612c16806100d65f395ff3fe608060405260043610610277575f3560e01c80634cdad5061161014a578063c6e6f592116100be578063e659d7d711610078578063e659d7d7146107c5578063ef8b30f7146107e4578063f29177c314610803578063f2fde38b14610856578063f3fef3a314610875578063ffa1ad7414610894575f80fd5b8063c6e6f5921461070b578063c804b1151461072a578063ce96cb7714610749578063d905777e14610768578063dd62ed3e14610787578063e15fc350146107a6575f80fd5b80638da5cb5b1161010f5780638da5cb5b1461061e57806393b89a841461065a57806395d89b411461067957806397feb23c1461068d578063a9059cbb146106ac578063c166c458146106cb575f80fd5b80634cdad5061461056d5780635ea1d6f81461058c57806360df7c6c146105bf57806370a08231146105eb578063715018a61461060a575f80fd5b8063207239c0116101ec57806335faa416116101a657806335faa416146104e4578063379607f5146104f85780633a4b66f1146105175780633b1dbfcc146104265780634b6a94cc1461051f5780634bf692061461055a575f80fd5b8063207239c014610426578063228cb7331461043b57806323b872dd1461044f57806326cf277a1461046e5780632f9ac83a1461048d578063313ce567146104c9575f80fd5b806307a2d13a1161023d57806307a2d13a1461034e578063095ea7b31461036d5780630a28a4771461039c5780630c11487f146103bb57806318160ddd146103e75780631e9a695014610407575f80fd5b8062fa3d501461029c57806301e1d114146102bb57806304ddc9d1146102e2578063052028d01461030e57806306fdde031461032d575f80fd5b36610298576102846108a8565b61028e33346108df565b610296610b43565b005b5f80fd5b3480156102a7575f80fd5b506102966102b636600461252a565b610b6d565b3480156102c6575f80fd5b506102cf610c1a565b6040519081526020015b60405180910390f35b3480156102ed575f80fd5b506103016102fc36600461252a565b610c28565b6040516102d99190612555565b348015610319575f80fd5b5061029661032836600461258f565b610d9b565b348015610338575f80fd5b50610341610e49565b6040516102d991906125d8565b348015610359575f80fd5b506102cf61036836600461252a565b610f09565b348015610378575f80fd5b5061038c6103873660046125ea565b610f4a565b60405190151581526020016102d9565b3480156103a7575f80fd5b506102cf6103b636600461252a565b610f63565b3480156103c6575f80fd5b506103cf610f98565b6040516001600160a01b0390911681526020016102d9565b3480156103f2575f80fd5b505f80516020612ba1833981519152546102cf565b348015610412575f80fd5b506102966104213660046125ea565b610fb3565b348015610431575f80fd5b506102cf61271081565b348015610446575f80fd5b506102cf611045565b34801561045a575f80fd5b5061038c610469366004612614565b61104e565b348015610479575f80fd5b506102966104883660046126c0565b611071565b348015610498575f80fd5b507f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c62701546001600160a01b03166103cf565b3480156104d4575f80fd5b50604051601281526020016102d9565b3480156104ef575f80fd5b506102966112e3565b348015610503575f80fd5b5061029661051236600461252a565b6112fd565b610296611492565b34801561052a575f80fd5b506103416040518060400160405280601081526020016f283ab13634b1a232b632b3b0ba34b7b760811b81525081565b61029661056836600461258f565b6114a4565b348015610578575f80fd5b506102cf61058736600461252a565b6114c8565b348015610597575f80fd5b507f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c62702546102cf565b3480156105ca575f80fd5b506105de6105d936600461258f565b6114d2565b6040516102d991906127be565b3480156105f6575f80fd5b506102cf61060536600461258f565b611547565b348015610615575f80fd5b5061029661156d565b348015610629575f80fd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166103cf565b348015610665575f80fd5b506105de610674366004612801565b61157e565b348015610684575f80fd5b506103416116a7565b348015610698575f80fd5b506102cf6106a73660046125ea565b6116e5565b3480156106b7575f80fd5b5061038c6106c63660046125ea565b61172b565b3480156106d6575f80fd5b506102cf6106e536600461258f565b6001600160a01b03165f9081525f80516020612b81833981519152602052604090205490565b348015610716575f80fd5b506102cf61072536600461252a565b611738565b348015610735575f80fd5b5061029661074436600461252a565b61176c565b348015610754575f80fd5b506102cf61076336600461258f565b6118a3565b348015610773575f80fd5b506102cf61078236600461258f565b6118b0565b348015610792575f80fd5b506102cf6107a136600461283b565b6118ba565b3480156107b1575f80fd5b506102966107c03660046125ea565b611903565b3480156107d0575f80fd5b506102966107df3660046125ea565b61193e565b3480156107ef575f80fd5b506102cf6107fe36600461252a565b611970565b34801561080e575f80fd5b506103cf61081d36600461252a565b5f9081527f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270460205260409020546001600160a01b031690565b348015610861575f80fd5b5061029661087036600461258f565b61197a565b348015610880575f80fd5b5061029661088f3660046125ea565b6119b9565b34801561089f575f80fd5b506102cf600281565b6108b0611a39565b6108dd60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005b90611a79565b565b5f80516020612bc18339815191525f6108f6610f98565b90505f83470390505f61090d828560020154611a80565b90506001600160a01b03861615610a79575f8183856001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa15801561095d573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109819190612867565b866001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109bd573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109e19190612867565b03010390505f610a0687610a005f80516020612ba18339815191525490565b84611a8f565b9050805f03610a2857604051631fe91a3f60e11b815260040160405180910390fd5b610a328882611ab1565b60408051888152602081018390526001600160a01b038a16917f1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90910160405180910390a250505b818501818103908214610ad757836001600160a01b031663c89e4361826040518263ffffffff1660e01b81526004015f604051808303818588803b158015610abf575f80fd5b505af1158015610ad1573d5f803e3d5ffd5b50505050505b8115610b3a576001850154610af5906001600160a01b031683611ae5565b60018501546040518381526001600160a01b03909116907f6c3b15f4a619d5331e4708792bb87f858ffb7a08f1a87aabca7cd15e51e04a639060200160405180910390a25b50505050505050565b6108dd5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f006108d7565b610b75611b58565b610b7d6108a8565b612710811115610ba057604051636bdaa48360e11b815260040160405180910390fd5b610baa5f806108df565b7f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c627028054908290556040515f80516020612bc18339815191529190839082907f67fb2216d844c3553cf557bffa85f0fde0294999f808e61dcae1773d50d5e150905f90a35050610c17610b43565b50565b5f610c23611bb3565b905090565b5f80610c32610f98565b90505f80826001600160a01b031663725c0503866040518263ffffffff1660e01b8152600401610c6491815260200190565b608060405180830381865afa158015610c7f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ca3919061287e565b935093505050815f03610cba57505f949350505050565b6002816002811115610cce57610cce612541565b03610cde57506005949350505050565b6001816002811115610cf257610cf2612541565b03610d0257506003949350505050565b5f836001600160a01b03166396106ae46040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d3f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d639190612867565b8301905082421015610d7b5750600195945050505050565b80421015610d8f5750600295945050505050565b50600495945050505050565b610da3611b58565b80610dad81611bca565b610db56108a8565b610dbf5f806108df565b7f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270180546001600160a01b031981166001600160a01b038581169182179093556040515f80516020612bc1833981519152939092169182907fa466bc069316af563b6a90c9b6f29a97752ac7be7c6bbf3767bc9b16c2fd90eb905f90a35050610e45610b43565b5050565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0380546060915f80516020612b6183398151915291610e87906128ca565b80601f0160208091040260200160405190810160405280929190818152602001828054610eb3906128ca565b8015610efe5780601f10610ed557610100808354040283529160200191610efe565b820191905f5260205f20905b815481529060010190602001808311610ee157829003601f168201915b505050505091505090565b5f80610f205f80516020612ba18339815191525490565b90508015610f4157610f3c610f33610c1a565b8490835f611bf1565b610f43565b825b9392505050565b5f33610f57818585611c3c565b60019150505b92915050565b5f80610f7a5f80516020612ba18339815191525490565b90508015610f4157610f3c81610f8e610c1a565b8591906001611bf1565b5f80516020612bc1833981519152546001600160a01b031690565b81610fbd81611bca565b610fc56108a8565b610fcf5f806108df565b5f610fd9836114c8565b9050610fe53384611c49565b610ff0338583611c7d565b60408051828152602081018590526001600160a01b0386169133917f5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea8348091015b60405180910390a350611040610b43565b505050565b5f610c23611db2565b5f3361105b858285611dd9565b611066858585611e3d565b506001949350505050565b5f61107a611e9a565b805490915060ff600160401b820416159067ffffffffffffffff165f811580156110a15750825b90505f8267ffffffffffffffff1660011480156110bd5750303b155b9050811580156110cb575080155b156110e95760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561111357845460ff60401b1916600160401b1785555b8661111d81611bca565b865161112881611bca565b6127108860400151111561114f57604051636bdaa48360e11b815260040160405180910390fd5b61119f88606001516040516020016111679190612919565b604051602081830303815290604052896060015160405160200161118b919061294a565b604051602081830303815290604052611ec2565b87516111aa90611ed4565b5f80516020612bc183398151915280546001600160a01b03199081166001600160a01b038c8116918217845560208c8101517f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270180549095169216919091179092556040808c01517f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c627025580518082018252601081526f283ab13634b1a232b632b3b0ba34b7b760811b938101939093525190917fae1ee3fdcc521a751b838d314164fcdbc7c3dcd4f0a910e153506bf5d89ca1ec9161128a91908d9061296c565b60405180910390a25050508315610b3a57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a150505050505050565b6112eb6108a8565b6112f55f806108df565b6108dd610b43565b6113056108a8565b61130f5f806108df565b61131881611ee5565b5f611321610f98565b604051636e93df0d60e01b8152600481018490529091506001600160a01b03821690636e93df0d906024015f604051808303815f87803b158015611363575f80fd5b505af1158015611375573d5f803e3d5ffd5b505060405163725c050360e01b8152600481018590525f92508291506001600160a01b0384169063725c050390602401608060405180830381865afa1580156113c0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906113e4919061287e565b9350509250506002808111156113fc576113fc612541565b81600281111561140e5761140e612541565b0361145a575f6114458361142d5f80516020612ba18339815191525490565b85611436611bb3565b61144091906129e1565b611a8f565b90506114513382611ab1565b5050505061148a565b604051849033907fd8138f8a3f377c5259ca548e70e4c2de94f129f5a11036a15b69513cba2b426a905f90a35050505b610c17610b43565b61149a6108a8565b6112f533346108df565b806114ae81611bca565b6114b66108a8565b6114c082346108df565b610e45610b43565b5f610f5d82610f09565b6001600160a01b0381165f9081525f80516020612b81833981519152602090815260409182902080548351818402810184019094528084526060939283018282801561153b57602002820191905f5260205f20905b815481526020019060010190808311611527575b50505050509050919050565b6001600160a01b03165f9081525f80516020612b61833981519152602052604090205490565b611575611b58565b6108dd5f611f3a565b6001600160a01b0382165f9081525f80516020612b818339815191526020526040812080546060925f80516020612bc18339815191529291908167ffffffffffffffff8111156115d0576115d0612652565b6040519080825280602002602001820160405280156115f9578160200160208202803683370190505b5090505f805b8381101561169a5787600581111561161957611619612541565b61163c86838154811061162e5761162e6129f4565b905f5260205f200154610c28565b600581111561164d5761164d612541565b0361169257848181548110611664576116646129f4565b905f5260205f200154838380600101945081518110611685576116856129f4565b6020026020010181815250505b6001016115ff565b5081529695505050505050565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0480546060915f80516020612b6183398151915291610e87906128ca565b6001600160a01b0382165f9081525f80516020612b818339815191526020526040812080548390811061171a5761171a6129f4565b905f5260205f200154905092915050565b5f33610f57818585611e3d565b5f8061174f5f80516020612ba18339815191525490565b90508015610f4157610f3c81611763610c1a565b8591905f611bf1565b6117746108a8565b61177e5f806108df565b61178781611ee5565b5f611790610f98565b60405163725c050360e01b8152600481018490529091505f906001600160a01b0383169063725c050390602401608060405180830381865afa1580156117d8573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906117fc919061287e565b50509150505f61180b82611970565b90506118173382611ab1565b60405163c804b11560e01b8152600481018590526001600160a01b0384169063c804b115906024015f604051808303815f87803b158015611856575f80fd5b505af1158015611868573d5f803e3d5ffd5b50506040518692503391507f853eba293fc3859716cf2e1a34b0c266f8265db181ec04748d0f25b8a19fc80e905f90a3505050610c17610b43565b5f610f5d61058783611547565b5f610f5d82611547565b6001600160a01b039182165f9081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace016020908152604080832093909416825291909152205490565b61190b6108a8565b6119155f806108df565b5f61191f826114c8565b905061192b3383611c49565b6119358382611faa565b50610e45610b43565b6119466108a8565b6119505f806108df565b5f61195a82610f63565b90506119663382611c49565b6119358383611faa565b5f610f5d82611738565b611982611b58565b6001600160a01b0381166119b057604051631e4fbdf760e01b81525f60048201526024015b60405180910390fd5b610c1781611f3a565b816119c381611bca565b6119cb6108a8565b6119d55f806108df565b5f6119df83610f63565b90506119eb3382611c49565b6119f6338585611c7d565b60408051848152602081018390526001600160a01b0386169133917f5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea83480910161102f565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c156108dd57604051633ee5aeb560e01b815260040160405180910390fd5b80825d5050565b5f610f43838361271084611bf1565b5f8215611aa757611aa28484845f611bf1565b611aa9565b835b949350505050565b6001600160a01b038216611ada5760405163ec442f0560e01b81525f60048201526024016119a7565b610e455f838361207b565b80471015611b0f5760405163cf47918160e01b8152476004820152602481018290526044016119a7565b611b28828260405180602001604052805f8152506121b4565b15611b31575050565b3d15611b3f57610e456121c9565b60405163d6bda27560e01b815260040160405180910390fd5b33611b8a7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146108dd5760405163118cdaa760e01b81523360048201526024016119a7565b5f611bbc611db2565b611bc46121d4565b01905090565b6001600160a01b038116610c175760405163d92e233d60e01b815260040160405180910390fd5b5f611c1e611bfe836122a7565b8015611c1957505f8480611c1457611c14612a08565b868809115b151590565b611c298686866122d3565b611c339190612a1c565b95945050505050565b6110408383836001612383565b6001600160a01b038216611c7257604051634b637e8f60e11b81525f60048201526024016119a7565b610e45825f8361207b565b805f03611c9d576040516374c51ccb60e11b815260040160405180910390fd5b5f80516020612bc18339815191525f611cb4610f98565b604051632efc584d60e11b81526001600160a01b038681166004830152602482018690529190911690635df8b09a906044016020604051808303815f875af1158015611d02573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611d269190612867565b6001600160a01b038681165f81815260038601602090815260408083208054600181018255908452828420018690558583526004880182529182902080546001600160a01b0319168417905590518781529394508493928816927fd71e6ec4eed83207b08d7ee4a0773c0ff8f8a1ab94b8ce85737fc0c5ea2b5f0c910160405180910390a45050505050565b5f4781611dd0825f80516020612bc183398151915260020154611a80565b90910392915050565b5f611de484846118ba565b90505f19811015611e375781811015611e2957604051637dc7a0d960e11b81526001600160a01b038416600482015260248101829052604481018390526064016119a7565b611e3784848484035f612383565b50505050565b6001600160a01b038316611e6657604051634b637e8f60e11b81525f60048201526024016119a7565b6001600160a01b038216611e8f5760405163ec442f0560e01b81525f60048201526024016119a7565b61104083838361207b565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610f5d565b611eca612467565b610e45828261248c565b611edc612467565b610c17816124dc565b5f8181527f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270460205260409020546001600160a01b03163314610c175760405163517907dd60e01b815260040160405180910390fd5b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b805f03611fca5760405163b7897d3760e01b815260040160405180910390fd5b611fd2610f98565b604051631af63e0160e21b81523360048201526001600160a01b038481166024830152604482018490529190911690636bd8f804906064015f604051808303815f87803b158015612021575f80fd5b505af1158015612033573d5f803e3d5ffd5b50506040518381526001600160a01b03851692503391507f78d93753d5a8009153a294711a82a3d1ba802938d66537c1b9142a053782d6939060200160405180910390a35050565b5f80516020612b618339815191526001600160a01b0384166120b55781816002015f8282546120aa9190612a1c565b909155506121259050565b6001600160a01b0384165f90815260208290526040902054828110156121075760405163391434e360e21b81526001600160a01b038616600482015260248101829052604481018490526064016119a7565b6001600160a01b0385165f9081526020839052604090209083900390555b6001600160a01b038316612143576002810180548390039055612161565b6001600160a01b0383165f9081526020829052604090208054830190555b826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516121a691815260200190565b60405180910390a350505050565b5f805f83516020850186885af1949350505050565b6040513d5f823e3d81fd5b5f806121de610f98565b9050806001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa15801561221c573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906122409190612867565b816001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa15801561227c573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906122a09190612867565b0391505090565b5f60028260038111156122bc576122bc612541565b6122c69190612a2f565b60ff166001149050919050565b5f805f6122e086866124e4565b91509150815f03612304578381816122fa576122fa612a08565b0492505050610f43565b81841161231b5761231b6003851502601118612500565b5f848688095f868103871696879004966002600389028118808a02820302808a02820302808a02820302808a02820302808a02820302808a02909103029181900381900460010185841190960395909502919093039390930492909217029150509392505050565b5f80516020612b618339815191526001600160a01b0385166123ba5760405163e602df0560e01b81525f60048201526024016119a7565b6001600160a01b0384166123e357604051634a1406b160e11b81525f60048201526024016119a7565b6001600160a01b038086165f9081526001830160209081526040808320938816835292905220839055811561246057836001600160a01b0316856001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258560405161245791815260200190565b60405180910390a35b5050505050565b61246f612511565b6108dd57604051631afcd79f60e31b815260040160405180910390fd5b612494612467565b5f80516020612b618339815191527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace036124cd8482612aa0565b5060048101611e378382612aa0565b611982612467565b5f805f1983850993909202808410938190039390930393915050565b634e487b715f52806020526024601cfd5b5f61251a611e9a565b54600160401b900460ff16919050565b5f6020828403121561253a575f80fd5b5035919050565b634e487b7160e01b5f52602160045260245ffd5b602081016006831061257557634e487b7160e01b5f52602160045260245ffd5b91905290565b6001600160a01b0381168114610c17575f80fd5b5f6020828403121561259f575f80fd5b8135610f438161257b565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f610f4360208301846125aa565b5f80604083850312156125fb575f80fd5b82356126068161257b565b946020939093013593505050565b5f805f60608486031215612626575f80fd5b83356126318161257b565b925060208401356126418161257b565b929592945050506040919091013590565b634e487b7160e01b5f52604160045260245ffd5b6040516080810167ffffffffffffffff8111828210171561268957612689612652565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156126b8576126b8612652565b604052919050565b5f80604083850312156126d1575f80fd5b82356126dc8161257b565b915060208381013567ffffffffffffffff808211156126f9575f80fd5b908501906080828803121561270c575f80fd5b612714612666565b823561271f8161257b565b81528284013561272e8161257b565b818501526040838101359082015260608301358281111561274d575f80fd5b80840193505087601f840112612761575f80fd5b82358281111561277357612773612652565b612785601f8201601f1916860161268f565b9250808352888582860101111561279a575f80fd5b80858501868501375f85828501015250816060820152809450505050509250929050565b602080825282518282018190525f9190848201906040850190845b818110156127f5578351835292840192918401916001016127d9565b50909695505050505050565b5f8060408385031215612812575f80fd5b823561281d8161257b565b9150602083013560068110612830575f80fd5b809150509250929050565b5f806040838503121561284c575f80fd5b82356128578161257b565b915060208301356128308161257b565b5f60208284031215612877575f80fd5b5051919050565b5f805f8060808587031215612891575f80fd5b845161289c8161257b565b8094505060208501519250604085015191506060850151600381106128bf575f80fd5b939692955090935050565b600181811c908216806128de57607f821691505b6020821081036128fc57634e487b7160e01b5f52602260045260245ffd5b50919050565b5f81518060208401855e5f93019283525090919050565b5f6129248284612902565b75205075626c69632044656c656761746564204b41494160501b81526016019392505050565b5f6129558284612902565b662d70644b41494160c81b81526007019392505050565b604081525f61297e60408301856125aa565b828103602084015260018060a01b0380855116825280602086015116602083015250604084015160408201526060840151608060608301526129c360808301826125aa565b9695505050505050565b634e487b7160e01b5f52601160045260245ffd5b81810381811115610f5d57610f5d6129cd565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601260045260245ffd5b80820180821115610f5d57610f5d6129cd565b5f60ff831680612a4d57634e487b7160e01b5f52601260045260245ffd5b8060ff84160691505092915050565b601f82111561104057805f5260205f20601f840160051c81016020851015612a815750805b601f840160051c820191505b81811015612460575f8155600101612a8d565b815167ffffffffffffffff811115612aba57612aba612652565b612ace81612ac884546128ca565b84612a5c565b602080601f831160018114612b01575f8415612aea5750858301515b5f19600386901b1c1916600185901b178555612b58565b5f85815260208120601f198616915b82811015612b2f57888601518255948401946001909101908401612b10565b5085821015612b4c57878501515f19600388901b60f8161c191681555b505060018460011b0185555b50505050505056fe52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace005bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270352c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace025bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c62700a26469706673582212202873104daeb97dafd272614940fdff2a727d0929b354a554e0431d591dab7f1d64736f6c63430008190033",
}

// PublicDelegationABI is the input ABI used to generate the binding from.
// Deprecated: Use PublicDelegationMetaData.ABI instead.
var PublicDelegationABI = PublicDelegationMetaData.ABI

// PublicDelegationBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const PublicDelegationBinRuntime = `608060405260043610610277575f3560e01c80634cdad5061161014a578063c6e6f592116100be578063e659d7d711610078578063e659d7d7146107c5578063ef8b30f7146107e4578063f29177c314610803578063f2fde38b14610856578063f3fef3a314610875578063ffa1ad7414610894575f80fd5b8063c6e6f5921461070b578063c804b1151461072a578063ce96cb7714610749578063d905777e14610768578063dd62ed3e14610787578063e15fc350146107a6575f80fd5b80638da5cb5b1161010f5780638da5cb5b1461061e57806393b89a841461065a57806395d89b411461067957806397feb23c1461068d578063a9059cbb146106ac578063c166c458146106cb575f80fd5b80634cdad5061461056d5780635ea1d6f81461058c57806360df7c6c146105bf57806370a08231146105eb578063715018a61461060a575f80fd5b8063207239c0116101ec57806335faa416116101a657806335faa416146104e4578063379607f5146104f85780633a4b66f1146105175780633b1dbfcc146104265780634b6a94cc1461051f5780634bf692061461055a575f80fd5b8063207239c014610426578063228cb7331461043b57806323b872dd1461044f57806326cf277a1461046e5780632f9ac83a1461048d578063313ce567146104c9575f80fd5b806307a2d13a1161023d57806307a2d13a1461034e578063095ea7b31461036d5780630a28a4771461039c5780630c11487f146103bb57806318160ddd146103e75780631e9a695014610407575f80fd5b8062fa3d501461029c57806301e1d114146102bb57806304ddc9d1146102e2578063052028d01461030e57806306fdde031461032d575f80fd5b36610298576102846108a8565b61028e33346108df565b610296610b43565b005b5f80fd5b3480156102a7575f80fd5b506102966102b636600461252a565b610b6d565b3480156102c6575f80fd5b506102cf610c1a565b6040519081526020015b60405180910390f35b3480156102ed575f80fd5b506103016102fc36600461252a565b610c28565b6040516102d99190612555565b348015610319575f80fd5b5061029661032836600461258f565b610d9b565b348015610338575f80fd5b50610341610e49565b6040516102d991906125d8565b348015610359575f80fd5b506102cf61036836600461252a565b610f09565b348015610378575f80fd5b5061038c6103873660046125ea565b610f4a565b60405190151581526020016102d9565b3480156103a7575f80fd5b506102cf6103b636600461252a565b610f63565b3480156103c6575f80fd5b506103cf610f98565b6040516001600160a01b0390911681526020016102d9565b3480156103f2575f80fd5b505f80516020612ba1833981519152546102cf565b348015610412575f80fd5b506102966104213660046125ea565b610fb3565b348015610431575f80fd5b506102cf61271081565b348015610446575f80fd5b506102cf611045565b34801561045a575f80fd5b5061038c610469366004612614565b61104e565b348015610479575f80fd5b506102966104883660046126c0565b611071565b348015610498575f80fd5b507f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c62701546001600160a01b03166103cf565b3480156104d4575f80fd5b50604051601281526020016102d9565b3480156104ef575f80fd5b506102966112e3565b348015610503575f80fd5b5061029661051236600461252a565b6112fd565b610296611492565b34801561052a575f80fd5b506103416040518060400160405280601081526020016f283ab13634b1a232b632b3b0ba34b7b760811b81525081565b61029661056836600461258f565b6114a4565b348015610578575f80fd5b506102cf61058736600461252a565b6114c8565b348015610597575f80fd5b507f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c62702546102cf565b3480156105ca575f80fd5b506105de6105d936600461258f565b6114d2565b6040516102d991906127be565b3480156105f6575f80fd5b506102cf61060536600461258f565b611547565b348015610615575f80fd5b5061029661156d565b348015610629575f80fd5b507f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b03166103cf565b348015610665575f80fd5b506105de610674366004612801565b61157e565b348015610684575f80fd5b506103416116a7565b348015610698575f80fd5b506102cf6106a73660046125ea565b6116e5565b3480156106b7575f80fd5b5061038c6106c63660046125ea565b61172b565b3480156106d6575f80fd5b506102cf6106e536600461258f565b6001600160a01b03165f9081525f80516020612b81833981519152602052604090205490565b348015610716575f80fd5b506102cf61072536600461252a565b611738565b348015610735575f80fd5b5061029661074436600461252a565b61176c565b348015610754575f80fd5b506102cf61076336600461258f565b6118a3565b348015610773575f80fd5b506102cf61078236600461258f565b6118b0565b348015610792575f80fd5b506102cf6107a136600461283b565b6118ba565b3480156107b1575f80fd5b506102966107c03660046125ea565b611903565b3480156107d0575f80fd5b506102966107df3660046125ea565b61193e565b3480156107ef575f80fd5b506102cf6107fe36600461252a565b611970565b34801561080e575f80fd5b506103cf61081d36600461252a565b5f9081527f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270460205260409020546001600160a01b031690565b348015610861575f80fd5b5061029661087036600461258f565b61197a565b348015610880575f80fd5b5061029661088f3660046125ea565b6119b9565b34801561089f575f80fd5b506102cf600281565b6108b0611a39565b6108dd60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005b90611a79565b565b5f80516020612bc18339815191525f6108f6610f98565b90505f83470390505f61090d828560020154611a80565b90506001600160a01b03861615610a79575f8183856001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa15801561095d573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109819190612867565b866001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa1580156109bd573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109e19190612867565b03010390505f610a0687610a005f80516020612ba18339815191525490565b84611a8f565b9050805f03610a2857604051631fe91a3f60e11b815260040160405180910390fd5b610a328882611ab1565b60408051888152602081018390526001600160a01b038a16917f1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90910160405180910390a250505b818501818103908214610ad757836001600160a01b031663c89e4361826040518263ffffffff1660e01b81526004015f604051808303818588803b158015610abf575f80fd5b505af1158015610ad1573d5f803e3d5ffd5b50505050505b8115610b3a576001850154610af5906001600160a01b031683611ae5565b60018501546040518381526001600160a01b03909116907f6c3b15f4a619d5331e4708792bb87f858ffb7a08f1a87aabca7cd15e51e04a639060200160405180910390a25b50505050505050565b6108dd5f7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f006108d7565b610b75611b58565b610b7d6108a8565b612710811115610ba057604051636bdaa48360e11b815260040160405180910390fd5b610baa5f806108df565b7f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c627028054908290556040515f80516020612bc18339815191529190839082907f67fb2216d844c3553cf557bffa85f0fde0294999f808e61dcae1773d50d5e150905f90a35050610c17610b43565b50565b5f610c23611bb3565b905090565b5f80610c32610f98565b90505f80826001600160a01b031663725c0503866040518263ffffffff1660e01b8152600401610c6491815260200190565b608060405180830381865afa158015610c7f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610ca3919061287e565b935093505050815f03610cba57505f949350505050565b6002816002811115610cce57610cce612541565b03610cde57506005949350505050565b6001816002811115610cf257610cf2612541565b03610d0257506003949350505050565b5f836001600160a01b03166396106ae46040518163ffffffff1660e01b8152600401602060405180830381865afa158015610d3f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d639190612867565b8301905082421015610d7b5750600195945050505050565b80421015610d8f5750600295945050505050565b50600495945050505050565b610da3611b58565b80610dad81611bca565b610db56108a8565b610dbf5f806108df565b7f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270180546001600160a01b031981166001600160a01b038581169182179093556040515f80516020612bc1833981519152939092169182907fa466bc069316af563b6a90c9b6f29a97752ac7be7c6bbf3767bc9b16c2fd90eb905f90a35050610e45610b43565b5050565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0380546060915f80516020612b6183398151915291610e87906128ca565b80601f0160208091040260200160405190810160405280929190818152602001828054610eb3906128ca565b8015610efe5780601f10610ed557610100808354040283529160200191610efe565b820191905f5260205f20905b815481529060010190602001808311610ee157829003601f168201915b505050505091505090565b5f80610f205f80516020612ba18339815191525490565b90508015610f4157610f3c610f33610c1a565b8490835f611bf1565b610f43565b825b9392505050565b5f33610f57818585611c3c565b60019150505b92915050565b5f80610f7a5f80516020612ba18339815191525490565b90508015610f4157610f3c81610f8e610c1a565b8591906001611bf1565b5f80516020612bc1833981519152546001600160a01b031690565b81610fbd81611bca565b610fc56108a8565b610fcf5f806108df565b5f610fd9836114c8565b9050610fe53384611c49565b610ff0338583611c7d565b60408051828152602081018590526001600160a01b0386169133917f5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea8348091015b60405180910390a350611040610b43565b505050565b5f610c23611db2565b5f3361105b858285611dd9565b611066858585611e3d565b506001949350505050565b5f61107a611e9a565b805490915060ff600160401b820416159067ffffffffffffffff165f811580156110a15750825b90505f8267ffffffffffffffff1660011480156110bd5750303b155b9050811580156110cb575080155b156110e95760405163f92ee8a960e01b815260040160405180910390fd5b845467ffffffffffffffff19166001178555831561111357845460ff60401b1916600160401b1785555b8661111d81611bca565b865161112881611bca565b6127108860400151111561114f57604051636bdaa48360e11b815260040160405180910390fd5b61119f88606001516040516020016111679190612919565b604051602081830303815290604052896060015160405160200161118b919061294a565b604051602081830303815290604052611ec2565b87516111aa90611ed4565b5f80516020612bc183398151915280546001600160a01b03199081166001600160a01b038c8116918217845560208c8101517f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270180549095169216919091179092556040808c01517f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c627025580518082018252601081526f283ab13634b1a232b632b3b0ba34b7b760811b938101939093525190917fae1ee3fdcc521a751b838d314164fcdbc7c3dcd4f0a910e153506bf5d89ca1ec9161128a91908d9061296c565b60405180910390a25050508315610b3a57845460ff60401b19168555604051600181527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d29060200160405180910390a150505050505050565b6112eb6108a8565b6112f55f806108df565b6108dd610b43565b6113056108a8565b61130f5f806108df565b61131881611ee5565b5f611321610f98565b604051636e93df0d60e01b8152600481018490529091506001600160a01b03821690636e93df0d906024015f604051808303815f87803b158015611363575f80fd5b505af1158015611375573d5f803e3d5ffd5b505060405163725c050360e01b8152600481018590525f92508291506001600160a01b0384169063725c050390602401608060405180830381865afa1580156113c0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906113e4919061287e565b9350509250506002808111156113fc576113fc612541565b81600281111561140e5761140e612541565b0361145a575f6114458361142d5f80516020612ba18339815191525490565b85611436611bb3565b61144091906129e1565b611a8f565b90506114513382611ab1565b5050505061148a565b604051849033907fd8138f8a3f377c5259ca548e70e4c2de94f129f5a11036a15b69513cba2b426a905f90a35050505b610c17610b43565b61149a6108a8565b6112f533346108df565b806114ae81611bca565b6114b66108a8565b6114c082346108df565b610e45610b43565b5f610f5d82610f09565b6001600160a01b0381165f9081525f80516020612b81833981519152602090815260409182902080548351818402810184019094528084526060939283018282801561153b57602002820191905f5260205f20905b815481526020019060010190808311611527575b50505050509050919050565b6001600160a01b03165f9081525f80516020612b61833981519152602052604090205490565b611575611b58565b6108dd5f611f3a565b6001600160a01b0382165f9081525f80516020612b818339815191526020526040812080546060925f80516020612bc18339815191529291908167ffffffffffffffff8111156115d0576115d0612652565b6040519080825280602002602001820160405280156115f9578160200160208202803683370190505b5090505f805b8381101561169a5787600581111561161957611619612541565b61163c86838154811061162e5761162e6129f4565b905f5260205f200154610c28565b600581111561164d5761164d612541565b0361169257848181548110611664576116646129f4565b905f5260205f200154838380600101945081518110611685576116856129f4565b6020026020010181815250505b6001016115ff565b5081529695505050505050565b7f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace0480546060915f80516020612b6183398151915291610e87906128ca565b6001600160a01b0382165f9081525f80516020612b818339815191526020526040812080548390811061171a5761171a6129f4565b905f5260205f200154905092915050565b5f33610f57818585611e3d565b5f8061174f5f80516020612ba18339815191525490565b90508015610f4157610f3c81611763610c1a565b8591905f611bf1565b6117746108a8565b61177e5f806108df565b61178781611ee5565b5f611790610f98565b60405163725c050360e01b8152600481018490529091505f906001600160a01b0383169063725c050390602401608060405180830381865afa1580156117d8573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906117fc919061287e565b50509150505f61180b82611970565b90506118173382611ab1565b60405163c804b11560e01b8152600481018590526001600160a01b0384169063c804b115906024015f604051808303815f87803b158015611856575f80fd5b505af1158015611868573d5f803e3d5ffd5b50506040518692503391507f853eba293fc3859716cf2e1a34b0c266f8265db181ec04748d0f25b8a19fc80e905f90a3505050610c17610b43565b5f610f5d61058783611547565b5f610f5d82611547565b6001600160a01b039182165f9081527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace016020908152604080832093909416825291909152205490565b61190b6108a8565b6119155f806108df565b5f61191f826114c8565b905061192b3383611c49565b6119358382611faa565b50610e45610b43565b6119466108a8565b6119505f806108df565b5f61195a82610f63565b90506119663382611c49565b6119358383611faa565b5f610f5d82611738565b611982611b58565b6001600160a01b0381166119b057604051631e4fbdf760e01b81525f60048201526024015b60405180910390fd5b610c1781611f3a565b816119c381611bca565b6119cb6108a8565b6119d55f806108df565b5f6119df83610f63565b90506119eb3382611c49565b6119f6338585611c7d565b60408051848152602081018390526001600160a01b0386169133917f5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea83480910161102f565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005c156108dd57604051633ee5aeb560e01b815260040160405180910390fd5b80825d5050565b5f610f43838361271084611bf1565b5f8215611aa757611aa28484845f611bf1565b611aa9565b835b949350505050565b6001600160a01b038216611ada5760405163ec442f0560e01b81525f60048201526024016119a7565b610e455f838361207b565b80471015611b0f5760405163cf47918160e01b8152476004820152602481018290526044016119a7565b611b28828260405180602001604052805f8152506121b4565b15611b31575050565b3d15611b3f57610e456121c9565b60405163d6bda27560e01b815260040160405180910390fd5b33611b8a7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c199300546001600160a01b031690565b6001600160a01b0316146108dd5760405163118cdaa760e01b81523360048201526024016119a7565b5f611bbc611db2565b611bc46121d4565b01905090565b6001600160a01b038116610c175760405163d92e233d60e01b815260040160405180910390fd5b5f611c1e611bfe836122a7565b8015611c1957505f8480611c1457611c14612a08565b868809115b151590565b611c298686866122d3565b611c339190612a1c565b95945050505050565b6110408383836001612383565b6001600160a01b038216611c7257604051634b637e8f60e11b81525f60048201526024016119a7565b610e45825f8361207b565b805f03611c9d576040516374c51ccb60e11b815260040160405180910390fd5b5f80516020612bc18339815191525f611cb4610f98565b604051632efc584d60e11b81526001600160a01b038681166004830152602482018690529190911690635df8b09a906044016020604051808303815f875af1158015611d02573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611d269190612867565b6001600160a01b038681165f81815260038601602090815260408083208054600181018255908452828420018690558583526004880182529182902080546001600160a01b0319168417905590518781529394508493928816927fd71e6ec4eed83207b08d7ee4a0773c0ff8f8a1ab94b8ce85737fc0c5ea2b5f0c910160405180910390a45050505050565b5f4781611dd0825f80516020612bc183398151915260020154611a80565b90910392915050565b5f611de484846118ba565b90505f19811015611e375781811015611e2957604051637dc7a0d960e11b81526001600160a01b038416600482015260248101829052604481018390526064016119a7565b611e3784848484035f612383565b50505050565b6001600160a01b038316611e6657604051634b637e8f60e11b81525f60048201526024016119a7565b6001600160a01b038216611e8f5760405163ec442f0560e01b81525f60048201526024016119a7565b61104083838361207b565b5f807ff0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00610f5d565b611eca612467565b610e45828261248c565b611edc612467565b610c17816124dc565b5f8181527f5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270460205260409020546001600160a01b03163314610c175760405163517907dd60e01b815260040160405180910390fd5b7f9016d09d72d40fdae2fd8ceac6b6234c7706214fd39c1cd1e609a0528c19930080546001600160a01b031981166001600160a01b03848116918217845560405192169182907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a3505050565b805f03611fca5760405163b7897d3760e01b815260040160405180910390fd5b611fd2610f98565b604051631af63e0160e21b81523360048201526001600160a01b038481166024830152604482018490529190911690636bd8f804906064015f604051808303815f87803b158015612021575f80fd5b505af1158015612033573d5f803e3d5ffd5b50506040518381526001600160a01b03851692503391507f78d93753d5a8009153a294711a82a3d1ba802938d66537c1b9142a053782d6939060200160405180910390a35050565b5f80516020612b618339815191526001600160a01b0384166120b55781816002015f8282546120aa9190612a1c565b909155506121259050565b6001600160a01b0384165f90815260208290526040902054828110156121075760405163391434e360e21b81526001600160a01b038616600482015260248101829052604481018490526064016119a7565b6001600160a01b0385165f9081526020839052604090209083900390555b6001600160a01b038316612143576002810180548390039055612161565b6001600160a01b0383165f9081526020829052604090208054830190555b826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef846040516121a691815260200190565b60405180910390a350505050565b5f805f83516020850186885af1949350505050565b6040513d5f823e3d81fd5b5f806121de610f98565b9050806001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa15801561221c573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906122409190612867565b816001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa15801561227c573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906122a09190612867565b0391505090565b5f60028260038111156122bc576122bc612541565b6122c69190612a2f565b60ff166001149050919050565b5f805f6122e086866124e4565b91509150815f03612304578381816122fa576122fa612a08565b0492505050610f43565b81841161231b5761231b6003851502601118612500565b5f848688095f868103871696879004966002600389028118808a02820302808a02820302808a02820302808a02820302808a02820302808a02909103029181900381900460010185841190960395909502919093039390930492909217029150509392505050565b5f80516020612b618339815191526001600160a01b0385166123ba5760405163e602df0560e01b81525f60048201526024016119a7565b6001600160a01b0384166123e357604051634a1406b160e11b81525f60048201526024016119a7565b6001600160a01b038086165f9081526001830160209081526040808320938816835292905220839055811561246057836001600160a01b0316856001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258560405161245791815260200190565b60405180910390a35b5050505050565b61246f612511565b6108dd57604051631afcd79f60e31b815260040160405180910390fd5b612494612467565b5f80516020612b618339815191527f52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace036124cd8482612aa0565b5060048101611e378382612aa0565b611982612467565b5f805f1983850993909202808410938190039390930393915050565b634e487b715f52806020526024601cfd5b5f61251a611e9a565b54600160401b900460ff16919050565b5f6020828403121561253a575f80fd5b5035919050565b634e487b7160e01b5f52602160045260245ffd5b602081016006831061257557634e487b7160e01b5f52602160045260245ffd5b91905290565b6001600160a01b0381168114610c17575f80fd5b5f6020828403121561259f575f80fd5b8135610f438161257b565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f610f4360208301846125aa565b5f80604083850312156125fb575f80fd5b82356126068161257b565b946020939093013593505050565b5f805f60608486031215612626575f80fd5b83356126318161257b565b925060208401356126418161257b565b929592945050506040919091013590565b634e487b7160e01b5f52604160045260245ffd5b6040516080810167ffffffffffffffff8111828210171561268957612689612652565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156126b8576126b8612652565b604052919050565b5f80604083850312156126d1575f80fd5b82356126dc8161257b565b915060208381013567ffffffffffffffff808211156126f9575f80fd5b908501906080828803121561270c575f80fd5b612714612666565b823561271f8161257b565b81528284013561272e8161257b565b818501526040838101359082015260608301358281111561274d575f80fd5b80840193505087601f840112612761575f80fd5b82358281111561277357612773612652565b612785601f8201601f1916860161268f565b9250808352888582860101111561279a575f80fd5b80858501868501375f85828501015250816060820152809450505050509250929050565b602080825282518282018190525f9190848201906040850190845b818110156127f5578351835292840192918401916001016127d9565b50909695505050505050565b5f8060408385031215612812575f80fd5b823561281d8161257b565b9150602083013560068110612830575f80fd5b809150509250929050565b5f806040838503121561284c575f80fd5b82356128578161257b565b915060208301356128308161257b565b5f60208284031215612877575f80fd5b5051919050565b5f805f8060808587031215612891575f80fd5b845161289c8161257b565b8094505060208501519250604085015191506060850151600381106128bf575f80fd5b939692955090935050565b600181811c908216806128de57607f821691505b6020821081036128fc57634e487b7160e01b5f52602260045260245ffd5b50919050565b5f81518060208401855e5f93019283525090919050565b5f6129248284612902565b75205075626c69632044656c656761746564204b41494160501b81526016019392505050565b5f6129558284612902565b662d70644b41494160c81b81526007019392505050565b604081525f61297e60408301856125aa565b828103602084015260018060a01b0380855116825280602086015116602083015250604084015160408201526060840151608060608301526129c360808301826125aa565b9695505050505050565b634e487b7160e01b5f52601160045260245ffd5b81810381811115610f5d57610f5d6129cd565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52601260045260245ffd5b80820180821115610f5d57610f5d6129cd565b5f60ff831680612a4d57634e487b7160e01b5f52601260045260245ffd5b8060ff84160691505092915050565b601f82111561104057805f5260205f20601f840160051c81016020851015612a815750805b601f840160051c820191505b81811015612460575f8155600101612a8d565b815167ffffffffffffffff811115612aba57612aba612652565b612ace81612ac884546128ca565b84612a5c565b602080601f831160018114612b01575f8415612aea5750858301515b5f19600386901b1c1916600185901b178555612b58565b5f85815260208120601f198616915b82811015612b2f57888601518255948401946001909101908401612b10565b5085821015612b4c57878501515f19600388901b60f8161c191681555b505060018460011b0185555b50505050505056fe52c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace005bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c6270352c63247e1f47db19d5ce0460030c497f067ca4cebf71ba98eeadabe20bace025bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c62700a26469706673582212202873104daeb97dafd272614940fdff2a727d0929b354a554e0431d591dab7f1d64736f6c63430008190033`

// Deprecated: Use PublicDelegationMetaData.Sigs instead.
// PublicDelegationFuncSigs maps the 4-byte function signature to its string representation.
var PublicDelegationFuncSigs = PublicDelegationMetaData.Sigs

// PublicDelegationBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PublicDelegationMetaData.Bin instead.
var PublicDelegationBin = PublicDelegationMetaData.Bin

// DeployPublicDelegation deploys a new Kaia contract, binding an instance of PublicDelegation to it.
func DeployPublicDelegation(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *PublicDelegation, error) {
	parsed, err := PublicDelegationMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PublicDelegationBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PublicDelegation{PublicDelegationCaller: PublicDelegationCaller{contract: contract}, PublicDelegationTransactor: PublicDelegationTransactor{contract: contract}, PublicDelegationFilterer: PublicDelegationFilterer{contract: contract}}, nil
}

// PublicDelegation is an auto generated Go binding around a Kaia contract.
type PublicDelegation struct {
	PublicDelegationCaller     // Read-only binding to the contract
	PublicDelegationTransactor // Write-only binding to the contract
	PublicDelegationFilterer   // Log filterer for contract events
}

// PublicDelegationCaller is an auto generated read-only Go binding around a Kaia contract.
type PublicDelegationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicDelegationTransactor is an auto generated write-only Go binding around a Kaia contract.
type PublicDelegationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicDelegationFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type PublicDelegationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicDelegationSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type PublicDelegationSession struct {
	Contract     *PublicDelegation // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PublicDelegationCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type PublicDelegationCallerSession struct {
	Contract *PublicDelegationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// PublicDelegationTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type PublicDelegationTransactorSession struct {
	Contract     *PublicDelegationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// PublicDelegationRaw is an auto generated low-level Go binding around a Kaia contract.
type PublicDelegationRaw struct {
	Contract *PublicDelegation // Generic contract binding to access the raw methods on
}

// PublicDelegationCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type PublicDelegationCallerRaw struct {
	Contract *PublicDelegationCaller // Generic read-only contract binding to access the raw methods on
}

// PublicDelegationTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type PublicDelegationTransactorRaw struct {
	Contract *PublicDelegationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPublicDelegation creates a new instance of PublicDelegation, bound to a specific deployed contract.
func NewPublicDelegation(address common.Address, backend bind.ContractBackend) (*PublicDelegation, error) {
	contract, err := bindPublicDelegation(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PublicDelegation{PublicDelegationCaller: PublicDelegationCaller{contract: contract}, PublicDelegationTransactor: PublicDelegationTransactor{contract: contract}, PublicDelegationFilterer: PublicDelegationFilterer{contract: contract}}, nil
}

// NewPublicDelegationCaller creates a new read-only instance of PublicDelegation, bound to a specific deployed contract.
func NewPublicDelegationCaller(address common.Address, caller bind.ContractCaller) (*PublicDelegationCaller, error) {
	contract, err := bindPublicDelegation(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationCaller{contract: contract}, nil
}

// NewPublicDelegationTransactor creates a new write-only instance of PublicDelegation, bound to a specific deployed contract.
func NewPublicDelegationTransactor(address common.Address, transactor bind.ContractTransactor) (*PublicDelegationTransactor, error) {
	contract, err := bindPublicDelegation(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationTransactor{contract: contract}, nil
}

// NewPublicDelegationFilterer creates a new log filterer instance of PublicDelegation, bound to a specific deployed contract.
func NewPublicDelegationFilterer(address common.Address, filterer bind.ContractFilterer) (*PublicDelegationFilterer, error) {
	contract, err := bindPublicDelegation(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationFilterer{contract: contract}, nil
}

// bindPublicDelegation binds a generic wrapper to an already deployed contract.
func bindPublicDelegation(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PublicDelegationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicDelegation *PublicDelegationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicDelegation.Contract.PublicDelegationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicDelegation *PublicDelegationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicDelegation.Contract.PublicDelegationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicDelegation *PublicDelegationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicDelegation.Contract.PublicDelegationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicDelegation *PublicDelegationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PublicDelegation.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicDelegation *PublicDelegationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicDelegation.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicDelegation *PublicDelegationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicDelegation.Contract.contract.Transact(opts, method, params...)
}

// COMMISSIONDENOMINATOR is a free data retrieval call binding the contract method 0x3b1dbfcc.
//
// Solidity: function COMMISSION_DENOMINATOR() view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) COMMISSIONDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "COMMISSION_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// COMMISSIONDENOMINATOR is a free data retrieval call binding the contract method 0x3b1dbfcc.
//
// Solidity: function COMMISSION_DENOMINATOR() view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) COMMISSIONDENOMINATOR() (*big.Int, error) {
	return _PublicDelegation.Contract.COMMISSIONDENOMINATOR(&_PublicDelegation.CallOpts)
}

// COMMISSIONDENOMINATOR is a free data retrieval call binding the contract method 0x3b1dbfcc.
//
// Solidity: function COMMISSION_DENOMINATOR() view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) COMMISSIONDENOMINATOR() (*big.Int, error) {
	return _PublicDelegation.Contract.COMMISSIONDENOMINATOR(&_PublicDelegation.CallOpts)
}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() view returns(string)
func (_PublicDelegation *PublicDelegationCaller) CONTRACTTYPE(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "CONTRACT_TYPE")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() view returns(string)
func (_PublicDelegation *PublicDelegationSession) CONTRACTTYPE() (string, error) {
	return _PublicDelegation.Contract.CONTRACTTYPE(&_PublicDelegation.CallOpts)
}

// CONTRACTTYPE is a free data retrieval call binding the contract method 0x4b6a94cc.
//
// Solidity: function CONTRACT_TYPE() view returns(string)
func (_PublicDelegation *PublicDelegationCallerSession) CONTRACTTYPE() (string, error) {
	return _PublicDelegation.Contract.CONTRACTTYPE(&_PublicDelegation.CallOpts)
}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) MAXCOMMISSIONRATE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "MAX_COMMISSION_RATE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) MAXCOMMISSIONRATE() (*big.Int, error) {
	return _PublicDelegation.Contract.MAXCOMMISSIONRATE(&_PublicDelegation.CallOpts)
}

// MAXCOMMISSIONRATE is a free data retrieval call binding the contract method 0x207239c0.
//
// Solidity: function MAX_COMMISSION_RATE() view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) MAXCOMMISSIONRATE() (*big.Int, error) {
	return _PublicDelegation.Contract.MAXCOMMISSIONRATE(&_PublicDelegation.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) VERSION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "VERSION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) VERSION() (*big.Int, error) {
	return _PublicDelegation.Contract.VERSION(&_PublicDelegation.CallOpts)
}

// VERSION is a free data retrieval call binding the contract method 0xffa1ad74.
//
// Solidity: function VERSION() view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) VERSION() (*big.Int, error) {
	return _PublicDelegation.Contract.VERSION(&_PublicDelegation.CallOpts)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.Allowance(&_PublicDelegation.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.Allowance(&_PublicDelegation.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.BalanceOf(&_PublicDelegation.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.BalanceOf(&_PublicDelegation.CallOpts, account)
}

// BaseCnStaking is a free data retrieval call binding the contract method 0x0c11487f.
//
// Solidity: function baseCnStaking() view returns(address)
func (_PublicDelegation *PublicDelegationCaller) BaseCnStaking(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "baseCnStaking")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BaseCnStaking is a free data retrieval call binding the contract method 0x0c11487f.
//
// Solidity: function baseCnStaking() view returns(address)
func (_PublicDelegation *PublicDelegationSession) BaseCnStaking() (common.Address, error) {
	return _PublicDelegation.Contract.BaseCnStaking(&_PublicDelegation.CallOpts)
}

// BaseCnStaking is a free data retrieval call binding the contract method 0x0c11487f.
//
// Solidity: function baseCnStaking() view returns(address)
func (_PublicDelegation *PublicDelegationCallerSession) BaseCnStaking() (common.Address, error) {
	return _PublicDelegation.Contract.BaseCnStaking(&_PublicDelegation.CallOpts)
}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) CommissionRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "commissionRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) CommissionRate() (*big.Int, error) {
	return _PublicDelegation.Contract.CommissionRate(&_PublicDelegation.CallOpts)
}

// CommissionRate is a free data retrieval call binding the contract method 0x5ea1d6f8.
//
// Solidity: function commissionRate() view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) CommissionRate() (*big.Int, error) {
	return _PublicDelegation.Contract.CommissionRate(&_PublicDelegation.CallOpts)
}

// CommissionTo is a free data retrieval call binding the contract method 0x2f9ac83a.
//
// Solidity: function commissionTo() view returns(address)
func (_PublicDelegation *PublicDelegationCaller) CommissionTo(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "commissionTo")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CommissionTo is a free data retrieval call binding the contract method 0x2f9ac83a.
//
// Solidity: function commissionTo() view returns(address)
func (_PublicDelegation *PublicDelegationSession) CommissionTo() (common.Address, error) {
	return _PublicDelegation.Contract.CommissionTo(&_PublicDelegation.CallOpts)
}

// CommissionTo is a free data retrieval call binding the contract method 0x2f9ac83a.
//
// Solidity: function commissionTo() view returns(address)
func (_PublicDelegation *PublicDelegationCallerSession) CommissionTo() (common.Address, error) {
	return _PublicDelegation.Contract.CommissionTo(&_PublicDelegation.CallOpts)
}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 _shares) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) ConvertToAssets(opts *bind.CallOpts, _shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "convertToAssets", _shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 _shares) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) ConvertToAssets(_shares *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.ConvertToAssets(&_PublicDelegation.CallOpts, _shares)
}

// ConvertToAssets is a free data retrieval call binding the contract method 0x07a2d13a.
//
// Solidity: function convertToAssets(uint256 _shares) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) ConvertToAssets(_shares *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.ConvertToAssets(&_PublicDelegation.CallOpts, _shares)
}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) ConvertToShares(opts *bind.CallOpts, _assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "convertToShares", _assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) ConvertToShares(_assets *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.ConvertToShares(&_PublicDelegation.CallOpts, _assets)
}

// ConvertToShares is a free data retrieval call binding the contract method 0xc6e6f592.
//
// Solidity: function convertToShares(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) ConvertToShares(_assets *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.ConvertToShares(&_PublicDelegation.CallOpts, _assets)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PublicDelegation *PublicDelegationCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PublicDelegation *PublicDelegationSession) Decimals() (uint8, error) {
	return _PublicDelegation.Contract.Decimals(&_PublicDelegation.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_PublicDelegation *PublicDelegationCallerSession) Decimals() (uint8, error) {
	return _PublicDelegation.Contract.Decimals(&_PublicDelegation.CallOpts)
}

// GetCurrentWithdrawalRequestState is a free data retrieval call binding the contract method 0x04ddc9d1.
//
// Solidity: function getCurrentWithdrawalRequestState(uint256 _requestId) view returns(uint8)
func (_PublicDelegation *PublicDelegationCaller) GetCurrentWithdrawalRequestState(opts *bind.CallOpts, _requestId *big.Int) (uint8, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "getCurrentWithdrawalRequestState", _requestId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetCurrentWithdrawalRequestState is a free data retrieval call binding the contract method 0x04ddc9d1.
//
// Solidity: function getCurrentWithdrawalRequestState(uint256 _requestId) view returns(uint8)
func (_PublicDelegation *PublicDelegationSession) GetCurrentWithdrawalRequestState(_requestId *big.Int) (uint8, error) {
	return _PublicDelegation.Contract.GetCurrentWithdrawalRequestState(&_PublicDelegation.CallOpts, _requestId)
}

// GetCurrentWithdrawalRequestState is a free data retrieval call binding the contract method 0x04ddc9d1.
//
// Solidity: function getCurrentWithdrawalRequestState(uint256 _requestId) view returns(uint8)
func (_PublicDelegation *PublicDelegationCallerSession) GetCurrentWithdrawalRequestState(_requestId *big.Int) (uint8, error) {
	return _PublicDelegation.Contract.GetCurrentWithdrawalRequestState(&_PublicDelegation.CallOpts, _requestId)
}

// GetUserRequestCount is a free data retrieval call binding the contract method 0xc166c458.
//
// Solidity: function getUserRequestCount(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) GetUserRequestCount(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "getUserRequestCount", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserRequestCount is a free data retrieval call binding the contract method 0xc166c458.
//
// Solidity: function getUserRequestCount(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) GetUserRequestCount(_owner common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.GetUserRequestCount(&_PublicDelegation.CallOpts, _owner)
}

// GetUserRequestCount is a free data retrieval call binding the contract method 0xc166c458.
//
// Solidity: function getUserRequestCount(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) GetUserRequestCount(_owner common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.GetUserRequestCount(&_PublicDelegation.CallOpts, _owner)
}

// GetUserRequestIds is a free data retrieval call binding the contract method 0x60df7c6c.
//
// Solidity: function getUserRequestIds(address _owner) view returns(uint256[])
func (_PublicDelegation *PublicDelegationCaller) GetUserRequestIds(opts *bind.CallOpts, _owner common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "getUserRequestIds", _owner)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUserRequestIds is a free data retrieval call binding the contract method 0x60df7c6c.
//
// Solidity: function getUserRequestIds(address _owner) view returns(uint256[])
func (_PublicDelegation *PublicDelegationSession) GetUserRequestIds(_owner common.Address) ([]*big.Int, error) {
	return _PublicDelegation.Contract.GetUserRequestIds(&_PublicDelegation.CallOpts, _owner)
}

// GetUserRequestIds is a free data retrieval call binding the contract method 0x60df7c6c.
//
// Solidity: function getUserRequestIds(address _owner) view returns(uint256[])
func (_PublicDelegation *PublicDelegationCallerSession) GetUserRequestIds(_owner common.Address) ([]*big.Int, error) {
	return _PublicDelegation.Contract.GetUserRequestIds(&_PublicDelegation.CallOpts, _owner)
}

// GetUserRequestIdsWithState is a free data retrieval call binding the contract method 0x93b89a84.
//
// Solidity: function getUserRequestIdsWithState(address _owner, uint8 _state) view returns(uint256[])
func (_PublicDelegation *PublicDelegationCaller) GetUserRequestIdsWithState(opts *bind.CallOpts, _owner common.Address, _state uint8) ([]*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "getUserRequestIdsWithState", _owner, _state)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUserRequestIdsWithState is a free data retrieval call binding the contract method 0x93b89a84.
//
// Solidity: function getUserRequestIdsWithState(address _owner, uint8 _state) view returns(uint256[])
func (_PublicDelegation *PublicDelegationSession) GetUserRequestIdsWithState(_owner common.Address, _state uint8) ([]*big.Int, error) {
	return _PublicDelegation.Contract.GetUserRequestIdsWithState(&_PublicDelegation.CallOpts, _owner, _state)
}

// GetUserRequestIdsWithState is a free data retrieval call binding the contract method 0x93b89a84.
//
// Solidity: function getUserRequestIdsWithState(address _owner, uint8 _state) view returns(uint256[])
func (_PublicDelegation *PublicDelegationCallerSession) GetUserRequestIdsWithState(_owner common.Address, _state uint8) ([]*big.Int, error) {
	return _PublicDelegation.Contract.GetUserRequestIdsWithState(&_PublicDelegation.CallOpts, _owner, _state)
}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) MaxRedeem(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "maxRedeem", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) MaxRedeem(_owner common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.MaxRedeem(&_PublicDelegation.CallOpts, _owner)
}

// MaxRedeem is a free data retrieval call binding the contract method 0xd905777e.
//
// Solidity: function maxRedeem(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) MaxRedeem(_owner common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.MaxRedeem(&_PublicDelegation.CallOpts, _owner)
}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) MaxWithdraw(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "maxWithdraw", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) MaxWithdraw(_owner common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.MaxWithdraw(&_PublicDelegation.CallOpts, _owner)
}

// MaxWithdraw is a free data retrieval call binding the contract method 0xce96cb77.
//
// Solidity: function maxWithdraw(address _owner) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) MaxWithdraw(_owner common.Address) (*big.Int, error) {
	return _PublicDelegation.Contract.MaxWithdraw(&_PublicDelegation.CallOpts, _owner)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_PublicDelegation *PublicDelegationCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_PublicDelegation *PublicDelegationSession) Name() (string, error) {
	return _PublicDelegation.Contract.Name(&_PublicDelegation.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_PublicDelegation *PublicDelegationCallerSession) Name() (string, error) {
	return _PublicDelegation.Contract.Name(&_PublicDelegation.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PublicDelegation *PublicDelegationCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PublicDelegation *PublicDelegationSession) Owner() (common.Address, error) {
	return _PublicDelegation.Contract.Owner(&_PublicDelegation.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PublicDelegation *PublicDelegationCallerSession) Owner() (common.Address, error) {
	return _PublicDelegation.Contract.Owner(&_PublicDelegation.CallOpts)
}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) PreviewDeposit(opts *bind.CallOpts, _assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "previewDeposit", _assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) PreviewDeposit(_assets *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.PreviewDeposit(&_PublicDelegation.CallOpts, _assets)
}

// PreviewDeposit is a free data retrieval call binding the contract method 0xef8b30f7.
//
// Solidity: function previewDeposit(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) PreviewDeposit(_assets *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.PreviewDeposit(&_PublicDelegation.CallOpts, _assets)
}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 _shares) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) PreviewRedeem(opts *bind.CallOpts, _shares *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "previewRedeem", _shares)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 _shares) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) PreviewRedeem(_shares *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.PreviewRedeem(&_PublicDelegation.CallOpts, _shares)
}

// PreviewRedeem is a free data retrieval call binding the contract method 0x4cdad506.
//
// Solidity: function previewRedeem(uint256 _shares) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) PreviewRedeem(_shares *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.PreviewRedeem(&_PublicDelegation.CallOpts, _shares)
}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) PreviewWithdraw(opts *bind.CallOpts, _assets *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "previewWithdraw", _assets)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) PreviewWithdraw(_assets *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.PreviewWithdraw(&_PublicDelegation.CallOpts, _assets)
}

// PreviewWithdraw is a free data retrieval call binding the contract method 0x0a28a477.
//
// Solidity: function previewWithdraw(uint256 _assets) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) PreviewWithdraw(_assets *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.PreviewWithdraw(&_PublicDelegation.CallOpts, _assets)
}

// RequestIdToOwner is a free data retrieval call binding the contract method 0xf29177c3.
//
// Solidity: function requestIdToOwner(uint256 _requestId) view returns(address)
func (_PublicDelegation *PublicDelegationCaller) RequestIdToOwner(opts *bind.CallOpts, _requestId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "requestIdToOwner", _requestId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RequestIdToOwner is a free data retrieval call binding the contract method 0xf29177c3.
//
// Solidity: function requestIdToOwner(uint256 _requestId) view returns(address)
func (_PublicDelegation *PublicDelegationSession) RequestIdToOwner(_requestId *big.Int) (common.Address, error) {
	return _PublicDelegation.Contract.RequestIdToOwner(&_PublicDelegation.CallOpts, _requestId)
}

// RequestIdToOwner is a free data retrieval call binding the contract method 0xf29177c3.
//
// Solidity: function requestIdToOwner(uint256 _requestId) view returns(address)
func (_PublicDelegation *PublicDelegationCallerSession) RequestIdToOwner(_requestId *big.Int) (common.Address, error) {
	return _PublicDelegation.Contract.RequestIdToOwner(&_PublicDelegation.CallOpts, _requestId)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) Reward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "reward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) Reward() (*big.Int, error) {
	return _PublicDelegation.Contract.Reward(&_PublicDelegation.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0x228cb733.
//
// Solidity: function reward() view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) Reward() (*big.Int, error) {
	return _PublicDelegation.Contract.Reward(&_PublicDelegation.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_PublicDelegation *PublicDelegationCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_PublicDelegation *PublicDelegationSession) Symbol() (string, error) {
	return _PublicDelegation.Contract.Symbol(&_PublicDelegation.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_PublicDelegation *PublicDelegationCallerSession) Symbol() (string, error) {
	return _PublicDelegation.Contract.Symbol(&_PublicDelegation.CallOpts)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) TotalAssets(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "totalAssets")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) TotalAssets() (*big.Int, error) {
	return _PublicDelegation.Contract.TotalAssets(&_PublicDelegation.CallOpts)
}

// TotalAssets is a free data retrieval call binding the contract method 0x01e1d114.
//
// Solidity: function totalAssets() view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) TotalAssets() (*big.Int, error) {
	return _PublicDelegation.Contract.TotalAssets(&_PublicDelegation.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) TotalSupply() (*big.Int, error) {
	return _PublicDelegation.Contract.TotalSupply(&_PublicDelegation.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) TotalSupply() (*big.Int, error) {
	return _PublicDelegation.Contract.TotalSupply(&_PublicDelegation.CallOpts)
}

// UserRequestIds is a free data retrieval call binding the contract method 0x97feb23c.
//
// Solidity: function userRequestIds(address _owner, uint256 _index) view returns(uint256)
func (_PublicDelegation *PublicDelegationCaller) UserRequestIds(opts *bind.CallOpts, _owner common.Address, _index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PublicDelegation.contract.Call(opts, &out, "userRequestIds", _owner, _index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserRequestIds is a free data retrieval call binding the contract method 0x97feb23c.
//
// Solidity: function userRequestIds(address _owner, uint256 _index) view returns(uint256)
func (_PublicDelegation *PublicDelegationSession) UserRequestIds(_owner common.Address, _index *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.UserRequestIds(&_PublicDelegation.CallOpts, _owner, _index)
}

// UserRequestIds is a free data retrieval call binding the contract method 0x97feb23c.
//
// Solidity: function userRequestIds(address _owner, uint256 _index) view returns(uint256)
func (_PublicDelegation *PublicDelegationCallerSession) UserRequestIds(_owner common.Address, _index *big.Int) (*big.Int, error) {
	return _PublicDelegation.Contract.UserRequestIds(&_PublicDelegation.CallOpts, _owner, _index)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Approve(&_PublicDelegation.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Approve(&_PublicDelegation.TransactOpts, spender, value)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _requestId) returns()
func (_PublicDelegation *PublicDelegationTransactor) CancelApprovedStakingWithdrawal(opts *bind.TransactOpts, _requestId *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "cancelApprovedStakingWithdrawal", _requestId)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _requestId) returns()
func (_PublicDelegation *PublicDelegationSession) CancelApprovedStakingWithdrawal(_requestId *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.CancelApprovedStakingWithdrawal(&_PublicDelegation.TransactOpts, _requestId)
}

// CancelApprovedStakingWithdrawal is a paid mutator transaction binding the contract method 0xc804b115.
//
// Solidity: function cancelApprovedStakingWithdrawal(uint256 _requestId) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) CancelApprovedStakingWithdrawal(_requestId *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.CancelApprovedStakingWithdrawal(&_PublicDelegation.TransactOpts, _requestId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 _requestId) returns()
func (_PublicDelegation *PublicDelegationTransactor) Claim(opts *bind.TransactOpts, _requestId *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "claim", _requestId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 _requestId) returns()
func (_PublicDelegation *PublicDelegationSession) Claim(_requestId *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Claim(&_PublicDelegation.TransactOpts, _requestId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 _requestId) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) Claim(_requestId *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Claim(&_PublicDelegation.TransactOpts, _requestId)
}

// Initialize is a paid mutator transaction binding the contract method 0x26cf277a.
//
// Solidity: function initialize(address _baseCnStaking, (address,address,uint256,string) _args) returns()
func (_PublicDelegation *PublicDelegationTransactor) Initialize(opts *bind.TransactOpts, _baseCnStaking common.Address, _args IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "initialize", _baseCnStaking, _args)
}

// Initialize is a paid mutator transaction binding the contract method 0x26cf277a.
//
// Solidity: function initialize(address _baseCnStaking, (address,address,uint256,string) _args) returns()
func (_PublicDelegation *PublicDelegationSession) Initialize(_baseCnStaking common.Address, _args IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Initialize(&_PublicDelegation.TransactOpts, _baseCnStaking, _args)
}

// Initialize is a paid mutator transaction binding the contract method 0x26cf277a.
//
// Solidity: function initialize(address _baseCnStaking, (address,address,uint256,string) _args) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) Initialize(_baseCnStaking common.Address, _args IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Initialize(&_PublicDelegation.TransactOpts, _baseCnStaking, _args)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address _recipient, uint256 _shares) returns()
func (_PublicDelegation *PublicDelegationTransactor) Redeem(opts *bind.TransactOpts, _recipient common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "redeem", _recipient, _shares)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address _recipient, uint256 _shares) returns()
func (_PublicDelegation *PublicDelegationSession) Redeem(_recipient common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Redeem(&_PublicDelegation.TransactOpts, _recipient, _shares)
}

// Redeem is a paid mutator transaction binding the contract method 0x1e9a6950.
//
// Solidity: function redeem(address _recipient, uint256 _shares) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) Redeem(_recipient common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Redeem(&_PublicDelegation.TransactOpts, _recipient, _shares)
}

// RedelegateByAssets is a paid mutator transaction binding the contract method 0xe659d7d7.
//
// Solidity: function redelegateByAssets(address _targetCnStaking, uint256 _assets) returns()
func (_PublicDelegation *PublicDelegationTransactor) RedelegateByAssets(opts *bind.TransactOpts, _targetCnStaking common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "redelegateByAssets", _targetCnStaking, _assets)
}

// RedelegateByAssets is a paid mutator transaction binding the contract method 0xe659d7d7.
//
// Solidity: function redelegateByAssets(address _targetCnStaking, uint256 _assets) returns()
func (_PublicDelegation *PublicDelegationSession) RedelegateByAssets(_targetCnStaking common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.RedelegateByAssets(&_PublicDelegation.TransactOpts, _targetCnStaking, _assets)
}

// RedelegateByAssets is a paid mutator transaction binding the contract method 0xe659d7d7.
//
// Solidity: function redelegateByAssets(address _targetCnStaking, uint256 _assets) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) RedelegateByAssets(_targetCnStaking common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.RedelegateByAssets(&_PublicDelegation.TransactOpts, _targetCnStaking, _assets)
}

// RedelegateByShares is a paid mutator transaction binding the contract method 0xe15fc350.
//
// Solidity: function redelegateByShares(address _targetCnStaking, uint256 _shares) returns()
func (_PublicDelegation *PublicDelegationTransactor) RedelegateByShares(opts *bind.TransactOpts, _targetCnStaking common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "redelegateByShares", _targetCnStaking, _shares)
}

// RedelegateByShares is a paid mutator transaction binding the contract method 0xe15fc350.
//
// Solidity: function redelegateByShares(address _targetCnStaking, uint256 _shares) returns()
func (_PublicDelegation *PublicDelegationSession) RedelegateByShares(_targetCnStaking common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.RedelegateByShares(&_PublicDelegation.TransactOpts, _targetCnStaking, _shares)
}

// RedelegateByShares is a paid mutator transaction binding the contract method 0xe15fc350.
//
// Solidity: function redelegateByShares(address _targetCnStaking, uint256 _shares) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) RedelegateByShares(_targetCnStaking common.Address, _shares *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.RedelegateByShares(&_PublicDelegation.TransactOpts, _targetCnStaking, _shares)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PublicDelegation *PublicDelegationTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PublicDelegation *PublicDelegationSession) RenounceOwnership() (*types.Transaction, error) {
	return _PublicDelegation.Contract.RenounceOwnership(&_PublicDelegation.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_PublicDelegation *PublicDelegationTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _PublicDelegation.Contract.RenounceOwnership(&_PublicDelegation.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_PublicDelegation *PublicDelegationTransactor) Stake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "stake")
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_PublicDelegation *PublicDelegationSession) Stake() (*types.Transaction, error) {
	return _PublicDelegation.Contract.Stake(&_PublicDelegation.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x3a4b66f1.
//
// Solidity: function stake() payable returns()
func (_PublicDelegation *PublicDelegationTransactorSession) Stake() (*types.Transaction, error) {
	return _PublicDelegation.Contract.Stake(&_PublicDelegation.TransactOpts)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_PublicDelegation *PublicDelegationTransactor) StakeFor(opts *bind.TransactOpts, _recipient common.Address) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "stakeFor", _recipient)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_PublicDelegation *PublicDelegationSession) StakeFor(_recipient common.Address) (*types.Transaction, error) {
	return _PublicDelegation.Contract.StakeFor(&_PublicDelegation.TransactOpts, _recipient)
}

// StakeFor is a paid mutator transaction binding the contract method 0x4bf69206.
//
// Solidity: function stakeFor(address _recipient) payable returns()
func (_PublicDelegation *PublicDelegationTransactorSession) StakeFor(_recipient common.Address) (*types.Transaction, error) {
	return _PublicDelegation.Contract.StakeFor(&_PublicDelegation.TransactOpts, _recipient)
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_PublicDelegation *PublicDelegationTransactor) Sweep(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "sweep")
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_PublicDelegation *PublicDelegationSession) Sweep() (*types.Transaction, error) {
	return _PublicDelegation.Contract.Sweep(&_PublicDelegation.TransactOpts)
}

// Sweep is a paid mutator transaction binding the contract method 0x35faa416.
//
// Solidity: function sweep() returns()
func (_PublicDelegation *PublicDelegationTransactorSession) Sweep() (*types.Transaction, error) {
	return _PublicDelegation.Contract.Sweep(&_PublicDelegation.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Transfer(&_PublicDelegation.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Transfer(&_PublicDelegation.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.TransferFrom(&_PublicDelegation.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_PublicDelegation *PublicDelegationTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.TransferFrom(&_PublicDelegation.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PublicDelegation *PublicDelegationTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PublicDelegation *PublicDelegationSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PublicDelegation.Contract.TransferOwnership(&_PublicDelegation.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PublicDelegation.Contract.TransferOwnership(&_PublicDelegation.TransactOpts, newOwner)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0x00fa3d50.
//
// Solidity: function updateCommissionRate(uint256 _commissionRate) returns()
func (_PublicDelegation *PublicDelegationTransactor) UpdateCommissionRate(opts *bind.TransactOpts, _commissionRate *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "updateCommissionRate", _commissionRate)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0x00fa3d50.
//
// Solidity: function updateCommissionRate(uint256 _commissionRate) returns()
func (_PublicDelegation *PublicDelegationSession) UpdateCommissionRate(_commissionRate *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.UpdateCommissionRate(&_PublicDelegation.TransactOpts, _commissionRate)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0x00fa3d50.
//
// Solidity: function updateCommissionRate(uint256 _commissionRate) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) UpdateCommissionRate(_commissionRate *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.UpdateCommissionRate(&_PublicDelegation.TransactOpts, _commissionRate)
}

// UpdateCommissionTo is a paid mutator transaction binding the contract method 0x052028d0.
//
// Solidity: function updateCommissionTo(address _commissionTo) returns()
func (_PublicDelegation *PublicDelegationTransactor) UpdateCommissionTo(opts *bind.TransactOpts, _commissionTo common.Address) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "updateCommissionTo", _commissionTo)
}

// UpdateCommissionTo is a paid mutator transaction binding the contract method 0x052028d0.
//
// Solidity: function updateCommissionTo(address _commissionTo) returns()
func (_PublicDelegation *PublicDelegationSession) UpdateCommissionTo(_commissionTo common.Address) (*types.Transaction, error) {
	return _PublicDelegation.Contract.UpdateCommissionTo(&_PublicDelegation.TransactOpts, _commissionTo)
}

// UpdateCommissionTo is a paid mutator transaction binding the contract method 0x052028d0.
//
// Solidity: function updateCommissionTo(address _commissionTo) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) UpdateCommissionTo(_commissionTo common.Address) (*types.Transaction, error) {
	return _PublicDelegation.Contract.UpdateCommissionTo(&_PublicDelegation.TransactOpts, _commissionTo)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _assets) returns()
func (_PublicDelegation *PublicDelegationTransactor) Withdraw(opts *bind.TransactOpts, _recipient common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.contract.Transact(opts, "withdraw", _recipient, _assets)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _assets) returns()
func (_PublicDelegation *PublicDelegationSession) Withdraw(_recipient common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Withdraw(&_PublicDelegation.TransactOpts, _recipient, _assets)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address _recipient, uint256 _assets) returns()
func (_PublicDelegation *PublicDelegationTransactorSession) Withdraw(_recipient common.Address, _assets *big.Int) (*types.Transaction, error) {
	return _PublicDelegation.Contract.Withdraw(&_PublicDelegation.TransactOpts, _recipient, _assets)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PublicDelegation *PublicDelegationTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicDelegation.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PublicDelegation *PublicDelegationSession) Receive() (*types.Transaction, error) {
	return _PublicDelegation.Contract.Receive(&_PublicDelegation.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_PublicDelegation *PublicDelegationTransactorSession) Receive() (*types.Transaction, error) {
	return _PublicDelegation.Contract.Receive(&_PublicDelegation.TransactOpts)
}

// PublicDelegationApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the PublicDelegation contract.
type PublicDelegationApprovalIterator struct {
	Event *PublicDelegationApproval // Event containing the contract specifics and raw log

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
func (it *PublicDelegationApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationApproval)
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
		it.Event = new(PublicDelegationApproval)
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
func (it *PublicDelegationApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationApproval represents a Approval event raised by the PublicDelegation contract.
type PublicDelegationApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_PublicDelegation *PublicDelegationFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*PublicDelegationApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationApprovalIterator{contract: _PublicDelegation.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_PublicDelegation *PublicDelegationFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *PublicDelegationApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationApproval)
				if err := _PublicDelegation.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_PublicDelegation *PublicDelegationFilterer) ParseApproval(log types.Log) (*PublicDelegationApproval, error) {
	event := new(PublicDelegationApproval)
	if err := _PublicDelegation.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the PublicDelegation contract.
type PublicDelegationClaimedIterator struct {
	Event *PublicDelegationClaimed // Event containing the contract specifics and raw log

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
func (it *PublicDelegationClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationClaimed)
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
		it.Event = new(PublicDelegationClaimed)
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
func (it *PublicDelegationClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationClaimed represents a Claimed event raised by the PublicDelegation contract.
type PublicDelegationClaimed struct {
	User      common.Address
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0xd8138f8a3f377c5259ca548e70e4c2de94f129f5a11036a15b69513cba2b426a.
//
// Solidity: event Claimed(address indexed user, uint256 indexed requestId)
func (_PublicDelegation *PublicDelegationFilterer) FilterClaimed(opts *bind.FilterOpts, user []common.Address, requestId []*big.Int) (*PublicDelegationClaimedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "Claimed", userRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationClaimedIterator{contract: _PublicDelegation.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0xd8138f8a3f377c5259ca548e70e4c2de94f129f5a11036a15b69513cba2b426a.
//
// Solidity: event Claimed(address indexed user, uint256 indexed requestId)
func (_PublicDelegation *PublicDelegationFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *PublicDelegationClaimed, user []common.Address, requestId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "Claimed", userRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationClaimed)
				if err := _PublicDelegation.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0xd8138f8a3f377c5259ca548e70e4c2de94f129f5a11036a15b69513cba2b426a.
//
// Solidity: event Claimed(address indexed user, uint256 indexed requestId)
func (_PublicDelegation *PublicDelegationFilterer) ParseClaimed(log types.Log) (*PublicDelegationClaimed, error) {
	event := new(PublicDelegationClaimed)
	if err := _PublicDelegation.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationDeployPublicDelegationIterator is returned from FilterDeployPublicDelegation and is used to iterate over the raw logs and unpacked data for DeployPublicDelegation events raised by the PublicDelegation contract.
type PublicDelegationDeployPublicDelegationIterator struct {
	Event *PublicDelegationDeployPublicDelegation // Event containing the contract specifics and raw log

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
func (it *PublicDelegationDeployPublicDelegationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationDeployPublicDelegation)
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
		it.Event = new(PublicDelegationDeployPublicDelegation)
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
func (it *PublicDelegationDeployPublicDelegationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationDeployPublicDelegationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationDeployPublicDelegation represents a DeployPublicDelegation event raised by the PublicDelegation contract.
type PublicDelegationDeployPublicDelegation struct {
	ContractType  string
	BaseCnStaking common.Address
	PdArgs        IPublicDelegationPDConstructorArgs
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterDeployPublicDelegation is a free log retrieval operation binding the contract event 0xae1ee3fdcc521a751b838d314164fcdbc7c3dcd4f0a910e153506bf5d89ca1ec.
//
// Solidity: event DeployPublicDelegation(string contractType, address indexed baseCnStaking, (address,address,uint256,string) pdArgs)
func (_PublicDelegation *PublicDelegationFilterer) FilterDeployPublicDelegation(opts *bind.FilterOpts, baseCnStaking []common.Address) (*PublicDelegationDeployPublicDelegationIterator, error) {

	var baseCnStakingRule []interface{}
	for _, baseCnStakingItem := range baseCnStaking {
		baseCnStakingRule = append(baseCnStakingRule, baseCnStakingItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "DeployPublicDelegation", baseCnStakingRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationDeployPublicDelegationIterator{contract: _PublicDelegation.contract, event: "DeployPublicDelegation", logs: logs, sub: sub}, nil
}

// WatchDeployPublicDelegation is a free log subscription operation binding the contract event 0xae1ee3fdcc521a751b838d314164fcdbc7c3dcd4f0a910e153506bf5d89ca1ec.
//
// Solidity: event DeployPublicDelegation(string contractType, address indexed baseCnStaking, (address,address,uint256,string) pdArgs)
func (_PublicDelegation *PublicDelegationFilterer) WatchDeployPublicDelegation(opts *bind.WatchOpts, sink chan<- *PublicDelegationDeployPublicDelegation, baseCnStaking []common.Address) (event.Subscription, error) {

	var baseCnStakingRule []interface{}
	for _, baseCnStakingItem := range baseCnStaking {
		baseCnStakingRule = append(baseCnStakingRule, baseCnStakingItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "DeployPublicDelegation", baseCnStakingRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationDeployPublicDelegation)
				if err := _PublicDelegation.contract.UnpackLog(event, "DeployPublicDelegation", log); err != nil {
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

// ParseDeployPublicDelegation is a log parse operation binding the contract event 0xae1ee3fdcc521a751b838d314164fcdbc7c3dcd4f0a910e153506bf5d89ca1ec.
//
// Solidity: event DeployPublicDelegation(string contractType, address indexed baseCnStaking, (address,address,uint256,string) pdArgs)
func (_PublicDelegation *PublicDelegationFilterer) ParseDeployPublicDelegation(log types.Log) (*PublicDelegationDeployPublicDelegation, error) {
	event := new(PublicDelegationDeployPublicDelegation)
	if err := _PublicDelegation.contract.UnpackLog(event, "DeployPublicDelegation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the PublicDelegation contract.
type PublicDelegationInitializedIterator struct {
	Event *PublicDelegationInitialized // Event containing the contract specifics and raw log

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
func (it *PublicDelegationInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationInitialized)
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
		it.Event = new(PublicDelegationInitialized)
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
func (it *PublicDelegationInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationInitialized represents a Initialized event raised by the PublicDelegation contract.
type PublicDelegationInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicDelegation *PublicDelegationFilterer) FilterInitialized(opts *bind.FilterOpts) (*PublicDelegationInitializedIterator, error) {

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PublicDelegationInitializedIterator{contract: _PublicDelegation.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_PublicDelegation *PublicDelegationFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PublicDelegationInitialized) (event.Subscription, error) {

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationInitialized)
				if err := _PublicDelegation.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_PublicDelegation *PublicDelegationFilterer) ParseInitialized(log types.Log) (*PublicDelegationInitialized, error) {
	event := new(PublicDelegationInitialized)
	if err := _PublicDelegation.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PublicDelegation contract.
type PublicDelegationOwnershipTransferredIterator struct {
	Event *PublicDelegationOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PublicDelegationOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationOwnershipTransferred)
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
		it.Event = new(PublicDelegationOwnershipTransferred)
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
func (it *PublicDelegationOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationOwnershipTransferred represents a OwnershipTransferred event raised by the PublicDelegation contract.
type PublicDelegationOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PublicDelegation *PublicDelegationFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PublicDelegationOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationOwnershipTransferredIterator{contract: _PublicDelegation.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_PublicDelegation *PublicDelegationFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PublicDelegationOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationOwnershipTransferred)
				if err := _PublicDelegation.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_PublicDelegation *PublicDelegationFilterer) ParseOwnershipTransferred(log types.Log) (*PublicDelegationOwnershipTransferred, error) {
	event := new(PublicDelegationOwnershipTransferred)
	if err := _PublicDelegation.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationRedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the PublicDelegation contract.
type PublicDelegationRedeemedIterator struct {
	Event *PublicDelegationRedeemed // Event containing the contract specifics and raw log

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
func (it *PublicDelegationRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationRedeemed)
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
		it.Event = new(PublicDelegationRedeemed)
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
func (it *PublicDelegationRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationRedeemed represents a Redeemed event raised by the PublicDelegation contract.
type PublicDelegationRedeemed struct {
	User      common.Address
	Recipient common.Address
	Assets    *big.Int
	Shares    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea83480.
//
// Solidity: event Redeemed(address indexed user, address indexed recipient, uint256 assets, uint256 shares)
func (_PublicDelegation *PublicDelegationFilterer) FilterRedeemed(opts *bind.FilterOpts, user []common.Address, recipient []common.Address) (*PublicDelegationRedeemedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "Redeemed", userRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationRedeemedIterator{contract: _PublicDelegation.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea83480.
//
// Solidity: event Redeemed(address indexed user, address indexed recipient, uint256 assets, uint256 shares)
func (_PublicDelegation *PublicDelegationFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *PublicDelegationRedeemed, user []common.Address, recipient []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "Redeemed", userRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationRedeemed)
				if err := _PublicDelegation.contract.UnpackLog(event, "Redeemed", log); err != nil {
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

// ParseRedeemed is a log parse operation binding the contract event 0x5cdf07ad0fc222442720b108e3ed4c4640f0fadc2ab2253e66f259a0fea83480.
//
// Solidity: event Redeemed(address indexed user, address indexed recipient, uint256 assets, uint256 shares)
func (_PublicDelegation *PublicDelegationFilterer) ParseRedeemed(log types.Log) (*PublicDelegationRedeemed, error) {
	event := new(PublicDelegationRedeemed)
	if err := _PublicDelegation.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationRedelegatedIterator is returned from FilterRedelegated and is used to iterate over the raw logs and unpacked data for Redelegated events raised by the PublicDelegation contract.
type PublicDelegationRedelegatedIterator struct {
	Event *PublicDelegationRedelegated // Event containing the contract specifics and raw log

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
func (it *PublicDelegationRedelegatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationRedelegated)
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
		it.Event = new(PublicDelegationRedelegated)
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
func (it *PublicDelegationRedelegatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationRedelegatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationRedelegated represents a Redelegated event raised by the PublicDelegation contract.
type PublicDelegationRedelegated struct {
	User            common.Address
	TargetCnStaking common.Address
	Assets          *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRedelegated is a free log retrieval operation binding the contract event 0x78d93753d5a8009153a294711a82a3d1ba802938d66537c1b9142a053782d693.
//
// Solidity: event Redelegated(address indexed user, address indexed targetCnStaking, uint256 assets)
func (_PublicDelegation *PublicDelegationFilterer) FilterRedelegated(opts *bind.FilterOpts, user []common.Address, targetCnStaking []common.Address) (*PublicDelegationRedelegatedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var targetCnStakingRule []interface{}
	for _, targetCnStakingItem := range targetCnStaking {
		targetCnStakingRule = append(targetCnStakingRule, targetCnStakingItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "Redelegated", userRule, targetCnStakingRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationRedelegatedIterator{contract: _PublicDelegation.contract, event: "Redelegated", logs: logs, sub: sub}, nil
}

// WatchRedelegated is a free log subscription operation binding the contract event 0x78d93753d5a8009153a294711a82a3d1ba802938d66537c1b9142a053782d693.
//
// Solidity: event Redelegated(address indexed user, address indexed targetCnStaking, uint256 assets)
func (_PublicDelegation *PublicDelegationFilterer) WatchRedelegated(opts *bind.WatchOpts, sink chan<- *PublicDelegationRedelegated, user []common.Address, targetCnStaking []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var targetCnStakingRule []interface{}
	for _, targetCnStakingItem := range targetCnStaking {
		targetCnStakingRule = append(targetCnStakingRule, targetCnStakingItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "Redelegated", userRule, targetCnStakingRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationRedelegated)
				if err := _PublicDelegation.contract.UnpackLog(event, "Redelegated", log); err != nil {
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

// ParseRedelegated is a log parse operation binding the contract event 0x78d93753d5a8009153a294711a82a3d1ba802938d66537c1b9142a053782d693.
//
// Solidity: event Redelegated(address indexed user, address indexed targetCnStaking, uint256 assets)
func (_PublicDelegation *PublicDelegationFilterer) ParseRedelegated(log types.Log) (*PublicDelegationRedelegated, error) {
	event := new(PublicDelegationRedelegated)
	if err := _PublicDelegation.contract.UnpackLog(event, "Redelegated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationRequestCancelWithdrawalIterator is returned from FilterRequestCancelWithdrawal and is used to iterate over the raw logs and unpacked data for RequestCancelWithdrawal events raised by the PublicDelegation contract.
type PublicDelegationRequestCancelWithdrawalIterator struct {
	Event *PublicDelegationRequestCancelWithdrawal // Event containing the contract specifics and raw log

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
func (it *PublicDelegationRequestCancelWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationRequestCancelWithdrawal)
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
		it.Event = new(PublicDelegationRequestCancelWithdrawal)
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
func (it *PublicDelegationRequestCancelWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationRequestCancelWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationRequestCancelWithdrawal represents a RequestCancelWithdrawal event raised by the PublicDelegation contract.
type PublicDelegationRequestCancelWithdrawal struct {
	User      common.Address
	RequestId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestCancelWithdrawal is a free log retrieval operation binding the contract event 0x853eba293fc3859716cf2e1a34b0c266f8265db181ec04748d0f25b8a19fc80e.
//
// Solidity: event RequestCancelWithdrawal(address indexed user, uint256 indexed requestId)
func (_PublicDelegation *PublicDelegationFilterer) FilterRequestCancelWithdrawal(opts *bind.FilterOpts, user []common.Address, requestId []*big.Int) (*PublicDelegationRequestCancelWithdrawalIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "RequestCancelWithdrawal", userRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationRequestCancelWithdrawalIterator{contract: _PublicDelegation.contract, event: "RequestCancelWithdrawal", logs: logs, sub: sub}, nil
}

// WatchRequestCancelWithdrawal is a free log subscription operation binding the contract event 0x853eba293fc3859716cf2e1a34b0c266f8265db181ec04748d0f25b8a19fc80e.
//
// Solidity: event RequestCancelWithdrawal(address indexed user, uint256 indexed requestId)
func (_PublicDelegation *PublicDelegationFilterer) WatchRequestCancelWithdrawal(opts *bind.WatchOpts, sink chan<- *PublicDelegationRequestCancelWithdrawal, user []common.Address, requestId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "RequestCancelWithdrawal", userRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationRequestCancelWithdrawal)
				if err := _PublicDelegation.contract.UnpackLog(event, "RequestCancelWithdrawal", log); err != nil {
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

// ParseRequestCancelWithdrawal is a log parse operation binding the contract event 0x853eba293fc3859716cf2e1a34b0c266f8265db181ec04748d0f25b8a19fc80e.
//
// Solidity: event RequestCancelWithdrawal(address indexed user, uint256 indexed requestId)
func (_PublicDelegation *PublicDelegationFilterer) ParseRequestCancelWithdrawal(log types.Log) (*PublicDelegationRequestCancelWithdrawal, error) {
	event := new(PublicDelegationRequestCancelWithdrawal)
	if err := _PublicDelegation.contract.UnpackLog(event, "RequestCancelWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationRequestWithdrawalIterator is returned from FilterRequestWithdrawal and is used to iterate over the raw logs and unpacked data for RequestWithdrawal events raised by the PublicDelegation contract.
type PublicDelegationRequestWithdrawalIterator struct {
	Event *PublicDelegationRequestWithdrawal // Event containing the contract specifics and raw log

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
func (it *PublicDelegationRequestWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationRequestWithdrawal)
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
		it.Event = new(PublicDelegationRequestWithdrawal)
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
func (it *PublicDelegationRequestWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationRequestWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationRequestWithdrawal represents a RequestWithdrawal event raised by the PublicDelegation contract.
type PublicDelegationRequestWithdrawal struct {
	User      common.Address
	Recipient common.Address
	RequestId *big.Int
	Assets    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRequestWithdrawal is a free log retrieval operation binding the contract event 0xd71e6ec4eed83207b08d7ee4a0773c0ff8f8a1ab94b8ce85737fc0c5ea2b5f0c.
//
// Solidity: event RequestWithdrawal(address indexed user, address indexed recipient, uint256 indexed requestId, uint256 assets)
func (_PublicDelegation *PublicDelegationFilterer) FilterRequestWithdrawal(opts *bind.FilterOpts, user []common.Address, recipient []common.Address, requestId []*big.Int) (*PublicDelegationRequestWithdrawalIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "RequestWithdrawal", userRule, recipientRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationRequestWithdrawalIterator{contract: _PublicDelegation.contract, event: "RequestWithdrawal", logs: logs, sub: sub}, nil
}

// WatchRequestWithdrawal is a free log subscription operation binding the contract event 0xd71e6ec4eed83207b08d7ee4a0773c0ff8f8a1ab94b8ce85737fc0c5ea2b5f0c.
//
// Solidity: event RequestWithdrawal(address indexed user, address indexed recipient, uint256 indexed requestId, uint256 assets)
func (_PublicDelegation *PublicDelegationFilterer) WatchRequestWithdrawal(opts *bind.WatchOpts, sink chan<- *PublicDelegationRequestWithdrawal, user []common.Address, recipient []common.Address, requestId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var requestIdRule []interface{}
	for _, requestIdItem := range requestId {
		requestIdRule = append(requestIdRule, requestIdItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "RequestWithdrawal", userRule, recipientRule, requestIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationRequestWithdrawal)
				if err := _PublicDelegation.contract.UnpackLog(event, "RequestWithdrawal", log); err != nil {
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

// ParseRequestWithdrawal is a log parse operation binding the contract event 0xd71e6ec4eed83207b08d7ee4a0773c0ff8f8a1ab94b8ce85737fc0c5ea2b5f0c.
//
// Solidity: event RequestWithdrawal(address indexed user, address indexed recipient, uint256 indexed requestId, uint256 assets)
func (_PublicDelegation *PublicDelegationFilterer) ParseRequestWithdrawal(log types.Log) (*PublicDelegationRequestWithdrawal, error) {
	event := new(PublicDelegationRequestWithdrawal)
	if err := _PublicDelegation.contract.UnpackLog(event, "RequestWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationSendCommissionIterator is returned from FilterSendCommission and is used to iterate over the raw logs and unpacked data for SendCommission events raised by the PublicDelegation contract.
type PublicDelegationSendCommissionIterator struct {
	Event *PublicDelegationSendCommission // Event containing the contract specifics and raw log

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
func (it *PublicDelegationSendCommissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationSendCommission)
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
		it.Event = new(PublicDelegationSendCommission)
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
func (it *PublicDelegationSendCommissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationSendCommissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationSendCommission represents a SendCommission event raised by the PublicDelegation contract.
type PublicDelegationSendCommission struct {
	CommissionTo common.Address
	Commission   *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSendCommission is a free log retrieval operation binding the contract event 0x6c3b15f4a619d5331e4708792bb87f858ffb7a08f1a87aabca7cd15e51e04a63.
//
// Solidity: event SendCommission(address indexed commissionTo, uint256 commission)
func (_PublicDelegation *PublicDelegationFilterer) FilterSendCommission(opts *bind.FilterOpts, commissionTo []common.Address) (*PublicDelegationSendCommissionIterator, error) {

	var commissionToRule []interface{}
	for _, commissionToItem := range commissionTo {
		commissionToRule = append(commissionToRule, commissionToItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "SendCommission", commissionToRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationSendCommissionIterator{contract: _PublicDelegation.contract, event: "SendCommission", logs: logs, sub: sub}, nil
}

// WatchSendCommission is a free log subscription operation binding the contract event 0x6c3b15f4a619d5331e4708792bb87f858ffb7a08f1a87aabca7cd15e51e04a63.
//
// Solidity: event SendCommission(address indexed commissionTo, uint256 commission)
func (_PublicDelegation *PublicDelegationFilterer) WatchSendCommission(opts *bind.WatchOpts, sink chan<- *PublicDelegationSendCommission, commissionTo []common.Address) (event.Subscription, error) {

	var commissionToRule []interface{}
	for _, commissionToItem := range commissionTo {
		commissionToRule = append(commissionToRule, commissionToItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "SendCommission", commissionToRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationSendCommission)
				if err := _PublicDelegation.contract.UnpackLog(event, "SendCommission", log); err != nil {
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

// ParseSendCommission is a log parse operation binding the contract event 0x6c3b15f4a619d5331e4708792bb87f858ffb7a08f1a87aabca7cd15e51e04a63.
//
// Solidity: event SendCommission(address indexed commissionTo, uint256 commission)
func (_PublicDelegation *PublicDelegationFilterer) ParseSendCommission(log types.Log) (*PublicDelegationSendCommission, error) {
	event := new(PublicDelegationSendCommission)
	if err := _PublicDelegation.contract.UnpackLog(event, "SendCommission", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the PublicDelegation contract.
type PublicDelegationStakedIterator struct {
	Event *PublicDelegationStaked // Event containing the contract specifics and raw log

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
func (it *PublicDelegationStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationStaked)
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
		it.Event = new(PublicDelegationStaked)
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
func (it *PublicDelegationStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationStaked represents a Staked event raised by the PublicDelegation contract.
type PublicDelegationStaked struct {
	User   common.Address
	Assets *big.Int
	Shares *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address indexed user, uint256 assets, uint256 shares)
func (_PublicDelegation *PublicDelegationFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address) (*PublicDelegationStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "Staked", userRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationStakedIterator{contract: _PublicDelegation.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address indexed user, uint256 assets, uint256 shares)
func (_PublicDelegation *PublicDelegationFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *PublicDelegationStaked, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "Staked", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationStaked)
				if err := _PublicDelegation.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x1449c6dd7851abc30abf37f57715f492010519147cc2652fbc38202c18a6ee90.
//
// Solidity: event Staked(address indexed user, uint256 assets, uint256 shares)
func (_PublicDelegation *PublicDelegationFilterer) ParseStaked(log types.Log) (*PublicDelegationStaked, error) {
	event := new(PublicDelegationStaked)
	if err := _PublicDelegation.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the PublicDelegation contract.
type PublicDelegationTransferIterator struct {
	Event *PublicDelegationTransfer // Event containing the contract specifics and raw log

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
func (it *PublicDelegationTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationTransfer)
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
		it.Event = new(PublicDelegationTransfer)
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
func (it *PublicDelegationTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationTransfer represents a Transfer event raised by the PublicDelegation contract.
type PublicDelegationTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_PublicDelegation *PublicDelegationFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*PublicDelegationTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationTransferIterator{contract: _PublicDelegation.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_PublicDelegation *PublicDelegationFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *PublicDelegationTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationTransfer)
				if err := _PublicDelegation.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_PublicDelegation *PublicDelegationFilterer) ParseTransfer(log types.Log) (*PublicDelegationTransfer, error) {
	event := new(PublicDelegationTransfer)
	if err := _PublicDelegation.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationUpdateCommissionRateIterator is returned from FilterUpdateCommissionRate and is used to iterate over the raw logs and unpacked data for UpdateCommissionRate events raised by the PublicDelegation contract.
type PublicDelegationUpdateCommissionRateIterator struct {
	Event *PublicDelegationUpdateCommissionRate // Event containing the contract specifics and raw log

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
func (it *PublicDelegationUpdateCommissionRateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationUpdateCommissionRate)
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
		it.Event = new(PublicDelegationUpdateCommissionRate)
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
func (it *PublicDelegationUpdateCommissionRateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationUpdateCommissionRateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationUpdateCommissionRate represents a UpdateCommissionRate event raised by the PublicDelegation contract.
type PublicDelegationUpdateCommissionRate struct {
	PrevCommissionRate *big.Int
	CommissionRate     *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpdateCommissionRate is a free log retrieval operation binding the contract event 0x67fb2216d844c3553cf557bffa85f0fde0294999f808e61dcae1773d50d5e150.
//
// Solidity: event UpdateCommissionRate(uint256 indexed prevCommissionRate, uint256 indexed commissionRate)
func (_PublicDelegation *PublicDelegationFilterer) FilterUpdateCommissionRate(opts *bind.FilterOpts, prevCommissionRate []*big.Int, commissionRate []*big.Int) (*PublicDelegationUpdateCommissionRateIterator, error) {

	var prevCommissionRateRule []interface{}
	for _, prevCommissionRateItem := range prevCommissionRate {
		prevCommissionRateRule = append(prevCommissionRateRule, prevCommissionRateItem)
	}
	var commissionRateRule []interface{}
	for _, commissionRateItem := range commissionRate {
		commissionRateRule = append(commissionRateRule, commissionRateItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "UpdateCommissionRate", prevCommissionRateRule, commissionRateRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationUpdateCommissionRateIterator{contract: _PublicDelegation.contract, event: "UpdateCommissionRate", logs: logs, sub: sub}, nil
}

// WatchUpdateCommissionRate is a free log subscription operation binding the contract event 0x67fb2216d844c3553cf557bffa85f0fde0294999f808e61dcae1773d50d5e150.
//
// Solidity: event UpdateCommissionRate(uint256 indexed prevCommissionRate, uint256 indexed commissionRate)
func (_PublicDelegation *PublicDelegationFilterer) WatchUpdateCommissionRate(opts *bind.WatchOpts, sink chan<- *PublicDelegationUpdateCommissionRate, prevCommissionRate []*big.Int, commissionRate []*big.Int) (event.Subscription, error) {

	var prevCommissionRateRule []interface{}
	for _, prevCommissionRateItem := range prevCommissionRate {
		prevCommissionRateRule = append(prevCommissionRateRule, prevCommissionRateItem)
	}
	var commissionRateRule []interface{}
	for _, commissionRateItem := range commissionRate {
		commissionRateRule = append(commissionRateRule, commissionRateItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "UpdateCommissionRate", prevCommissionRateRule, commissionRateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationUpdateCommissionRate)
				if err := _PublicDelegation.contract.UnpackLog(event, "UpdateCommissionRate", log); err != nil {
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

// ParseUpdateCommissionRate is a log parse operation binding the contract event 0x67fb2216d844c3553cf557bffa85f0fde0294999f808e61dcae1773d50d5e150.
//
// Solidity: event UpdateCommissionRate(uint256 indexed prevCommissionRate, uint256 indexed commissionRate)
func (_PublicDelegation *PublicDelegationFilterer) ParseUpdateCommissionRate(log types.Log) (*PublicDelegationUpdateCommissionRate, error) {
	event := new(PublicDelegationUpdateCommissionRate)
	if err := _PublicDelegation.contract.UnpackLog(event, "UpdateCommissionRate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PublicDelegationUpdateCommissionToIterator is returned from FilterUpdateCommissionTo and is used to iterate over the raw logs and unpacked data for UpdateCommissionTo events raised by the PublicDelegation contract.
type PublicDelegationUpdateCommissionToIterator struct {
	Event *PublicDelegationUpdateCommissionTo // Event containing the contract specifics and raw log

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
func (it *PublicDelegationUpdateCommissionToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicDelegationUpdateCommissionTo)
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
		it.Event = new(PublicDelegationUpdateCommissionTo)
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
func (it *PublicDelegationUpdateCommissionToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicDelegationUpdateCommissionToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicDelegationUpdateCommissionTo represents a UpdateCommissionTo event raised by the PublicDelegation contract.
type PublicDelegationUpdateCommissionTo struct {
	PrevCommissionTo common.Address
	CommissionTo     common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterUpdateCommissionTo is a free log retrieval operation binding the contract event 0xa466bc069316af563b6a90c9b6f29a97752ac7be7c6bbf3767bc9b16c2fd90eb.
//
// Solidity: event UpdateCommissionTo(address indexed prevCommissionTo, address indexed commissionTo)
func (_PublicDelegation *PublicDelegationFilterer) FilterUpdateCommissionTo(opts *bind.FilterOpts, prevCommissionTo []common.Address, commissionTo []common.Address) (*PublicDelegationUpdateCommissionToIterator, error) {

	var prevCommissionToRule []interface{}
	for _, prevCommissionToItem := range prevCommissionTo {
		prevCommissionToRule = append(prevCommissionToRule, prevCommissionToItem)
	}
	var commissionToRule []interface{}
	for _, commissionToItem := range commissionTo {
		commissionToRule = append(commissionToRule, commissionToItem)
	}

	logs, sub, err := _PublicDelegation.contract.FilterLogs(opts, "UpdateCommissionTo", prevCommissionToRule, commissionToRule)
	if err != nil {
		return nil, err
	}
	return &PublicDelegationUpdateCommissionToIterator{contract: _PublicDelegation.contract, event: "UpdateCommissionTo", logs: logs, sub: sub}, nil
}

// WatchUpdateCommissionTo is a free log subscription operation binding the contract event 0xa466bc069316af563b6a90c9b6f29a97752ac7be7c6bbf3767bc9b16c2fd90eb.
//
// Solidity: event UpdateCommissionTo(address indexed prevCommissionTo, address indexed commissionTo)
func (_PublicDelegation *PublicDelegationFilterer) WatchUpdateCommissionTo(opts *bind.WatchOpts, sink chan<- *PublicDelegationUpdateCommissionTo, prevCommissionTo []common.Address, commissionTo []common.Address) (event.Subscription, error) {

	var prevCommissionToRule []interface{}
	for _, prevCommissionToItem := range prevCommissionTo {
		prevCommissionToRule = append(prevCommissionToRule, prevCommissionToItem)
	}
	var commissionToRule []interface{}
	for _, commissionToItem := range commissionTo {
		commissionToRule = append(commissionToRule, commissionToItem)
	}

	logs, sub, err := _PublicDelegation.contract.WatchLogs(opts, "UpdateCommissionTo", prevCommissionToRule, commissionToRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicDelegationUpdateCommissionTo)
				if err := _PublicDelegation.contract.UnpackLog(event, "UpdateCommissionTo", log); err != nil {
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

// ParseUpdateCommissionTo is a log parse operation binding the contract event 0xa466bc069316af563b6a90c9b6f29a97752ac7be7c6bbf3767bc9b16c2fd90eb.
//
// Solidity: event UpdateCommissionTo(address indexed prevCommissionTo, address indexed commissionTo)
func (_PublicDelegation *PublicDelegationFilterer) ParseUpdateCommissionTo(log types.Log) (*PublicDelegationUpdateCommissionTo, error) {
	event := new(PublicDelegationUpdateCommissionTo)
	if err := _PublicDelegation.contract.UnpackLog(event, "UpdateCommissionTo", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ReentrancyGuardTransientMetaData contains all meta data concerning the ReentrancyGuardTransient contract.
var ReentrancyGuardTransientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"}]",
}

// ReentrancyGuardTransientABI is the input ABI used to generate the binding from.
// Deprecated: Use ReentrancyGuardTransientMetaData.ABI instead.
var ReentrancyGuardTransientABI = ReentrancyGuardTransientMetaData.ABI

// ReentrancyGuardTransientBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ReentrancyGuardTransientBinRuntime = ``

// ReentrancyGuardTransient is an auto generated Go binding around a Kaia contract.
type ReentrancyGuardTransient struct {
	ReentrancyGuardTransientCaller     // Read-only binding to the contract
	ReentrancyGuardTransientTransactor // Write-only binding to the contract
	ReentrancyGuardTransientFilterer   // Log filterer for contract events
}

// ReentrancyGuardTransientCaller is an auto generated read-only Go binding around a Kaia contract.
type ReentrancyGuardTransientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentrancyGuardTransientTransactor is an auto generated write-only Go binding around a Kaia contract.
type ReentrancyGuardTransientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentrancyGuardTransientFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ReentrancyGuardTransientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ReentrancyGuardTransientSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ReentrancyGuardTransientSession struct {
	Contract     *ReentrancyGuardTransient // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ReentrancyGuardTransientCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ReentrancyGuardTransientCallerSession struct {
	Contract *ReentrancyGuardTransientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// ReentrancyGuardTransientTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ReentrancyGuardTransientTransactorSession struct {
	Contract     *ReentrancyGuardTransientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// ReentrancyGuardTransientRaw is an auto generated low-level Go binding around a Kaia contract.
type ReentrancyGuardTransientRaw struct {
	Contract *ReentrancyGuardTransient // Generic contract binding to access the raw methods on
}

// ReentrancyGuardTransientCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ReentrancyGuardTransientCallerRaw struct {
	Contract *ReentrancyGuardTransientCaller // Generic read-only contract binding to access the raw methods on
}

// ReentrancyGuardTransientTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ReentrancyGuardTransientTransactorRaw struct {
	Contract *ReentrancyGuardTransientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReentrancyGuardTransient creates a new instance of ReentrancyGuardTransient, bound to a specific deployed contract.
func NewReentrancyGuardTransient(address common.Address, backend bind.ContractBackend) (*ReentrancyGuardTransient, error) {
	contract, err := bindReentrancyGuardTransient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ReentrancyGuardTransient{ReentrancyGuardTransientCaller: ReentrancyGuardTransientCaller{contract: contract}, ReentrancyGuardTransientTransactor: ReentrancyGuardTransientTransactor{contract: contract}, ReentrancyGuardTransientFilterer: ReentrancyGuardTransientFilterer{contract: contract}}, nil
}

// NewReentrancyGuardTransientCaller creates a new read-only instance of ReentrancyGuardTransient, bound to a specific deployed contract.
func NewReentrancyGuardTransientCaller(address common.Address, caller bind.ContractCaller) (*ReentrancyGuardTransientCaller, error) {
	contract, err := bindReentrancyGuardTransient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ReentrancyGuardTransientCaller{contract: contract}, nil
}

// NewReentrancyGuardTransientTransactor creates a new write-only instance of ReentrancyGuardTransient, bound to a specific deployed contract.
func NewReentrancyGuardTransientTransactor(address common.Address, transactor bind.ContractTransactor) (*ReentrancyGuardTransientTransactor, error) {
	contract, err := bindReentrancyGuardTransient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ReentrancyGuardTransientTransactor{contract: contract}, nil
}

// NewReentrancyGuardTransientFilterer creates a new log filterer instance of ReentrancyGuardTransient, bound to a specific deployed contract.
func NewReentrancyGuardTransientFilterer(address common.Address, filterer bind.ContractFilterer) (*ReentrancyGuardTransientFilterer, error) {
	contract, err := bindReentrancyGuardTransient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ReentrancyGuardTransientFilterer{contract: contract}, nil
}

// bindReentrancyGuardTransient binds a generic wrapper to an already deployed contract.
func bindReentrancyGuardTransient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ReentrancyGuardTransientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReentrancyGuardTransient *ReentrancyGuardTransientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReentrancyGuardTransient.Contract.ReentrancyGuardTransientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReentrancyGuardTransient *ReentrancyGuardTransientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReentrancyGuardTransient.Contract.ReentrancyGuardTransientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReentrancyGuardTransient *ReentrancyGuardTransientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReentrancyGuardTransient.Contract.ReentrancyGuardTransientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ReentrancyGuardTransient *ReentrancyGuardTransientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ReentrancyGuardTransient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ReentrancyGuardTransient *ReentrancyGuardTransientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ReentrancyGuardTransient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ReentrancyGuardTransient *ReentrancyGuardTransientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ReentrancyGuardTransient.Contract.contract.Transact(opts, method, params...)
}

// SafeCastMetaData contains all meta data concerning the SafeCast contract.
var SafeCastMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"bits\",\"type\":\"uint8\"},{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"SafeCastOverflowedIntDowncast\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int256\",\"name\":\"value\",\"type\":\"int256\"}],\"name\":\"SafeCastOverflowedIntToUint\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"bits\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SafeCastOverflowedUintDowncast\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"SafeCastOverflowedUintToInt\",\"type\":\"error\"}]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea264697066735822122004223638d8d1e86abf400699f30c266341b4680b04bb0c8679d0cbf532d1c6d664736f6c63430008190033",
}

// SafeCastABI is the input ABI used to generate the binding from.
// Deprecated: Use SafeCastMetaData.ABI instead.
var SafeCastABI = SafeCastMetaData.ABI

// SafeCastBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const SafeCastBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea264697066735822122004223638d8d1e86abf400699f30c266341b4680b04bb0c8679d0cbf532d1c6d664736f6c63430008190033`

// SafeCastBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SafeCastMetaData.Bin instead.
var SafeCastBin = SafeCastMetaData.Bin

// DeploySafeCast deploys a new Kaia contract, binding an instance of SafeCast to it.
func DeploySafeCast(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeCast, error) {
	parsed, err := SafeCastMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SafeCastBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeCast{SafeCastCaller: SafeCastCaller{contract: contract}, SafeCastTransactor: SafeCastTransactor{contract: contract}, SafeCastFilterer: SafeCastFilterer{contract: contract}}, nil
}

// SafeCast is an auto generated Go binding around a Kaia contract.
type SafeCast struct {
	SafeCastCaller     // Read-only binding to the contract
	SafeCastTransactor // Write-only binding to the contract
	SafeCastFilterer   // Log filterer for contract events
}

// SafeCastCaller is an auto generated read-only Go binding around a Kaia contract.
type SafeCastCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeCastTransactor is an auto generated write-only Go binding around a Kaia contract.
type SafeCastTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeCastFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type SafeCastFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeCastSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type SafeCastSession struct {
	Contract     *SafeCast         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeCastCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type SafeCastCallerSession struct {
	Contract *SafeCastCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeCastTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type SafeCastTransactorSession struct {
	Contract     *SafeCastTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeCastRaw is an auto generated low-level Go binding around a Kaia contract.
type SafeCastRaw struct {
	Contract *SafeCast // Generic contract binding to access the raw methods on
}

// SafeCastCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type SafeCastCallerRaw struct {
	Contract *SafeCastCaller // Generic read-only contract binding to access the raw methods on
}

// SafeCastTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type SafeCastTransactorRaw struct {
	Contract *SafeCastTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeCast creates a new instance of SafeCast, bound to a specific deployed contract.
func NewSafeCast(address common.Address, backend bind.ContractBackend) (*SafeCast, error) {
	contract, err := bindSafeCast(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeCast{SafeCastCaller: SafeCastCaller{contract: contract}, SafeCastTransactor: SafeCastTransactor{contract: contract}, SafeCastFilterer: SafeCastFilterer{contract: contract}}, nil
}

// NewSafeCastCaller creates a new read-only instance of SafeCast, bound to a specific deployed contract.
func NewSafeCastCaller(address common.Address, caller bind.ContractCaller) (*SafeCastCaller, error) {
	contract, err := bindSafeCast(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeCastCaller{contract: contract}, nil
}

// NewSafeCastTransactor creates a new write-only instance of SafeCast, bound to a specific deployed contract.
func NewSafeCastTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeCastTransactor, error) {
	contract, err := bindSafeCast(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeCastTransactor{contract: contract}, nil
}

// NewSafeCastFilterer creates a new log filterer instance of SafeCast, bound to a specific deployed contract.
func NewSafeCastFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeCastFilterer, error) {
	contract, err := bindSafeCast(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeCastFilterer{contract: contract}, nil
}

// bindSafeCast binds a generic wrapper to an already deployed contract.
func bindSafeCast(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SafeCastMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeCast *SafeCastRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeCast.Contract.SafeCastCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeCast *SafeCastRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeCast.Contract.SafeCastTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeCast *SafeCastRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeCast.Contract.SafeCastTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeCast *SafeCastCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SafeCast.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeCast *SafeCastTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeCast.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeCast *SafeCastTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeCast.Contract.contract.Transact(opts, method, params...)
}

// TransientSlotMetaData contains all meta data concerning the TransientSlot contract.
var TransientSlotMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220f3766cd635456275dc53c8d16aa6229230804ad1f4ee4f00852a5b42220a365464736f6c63430008190033",
}

// TransientSlotABI is the input ABI used to generate the binding from.
// Deprecated: Use TransientSlotMetaData.ABI instead.
var TransientSlotABI = TransientSlotMetaData.ABI

// TransientSlotBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const TransientSlotBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220f3766cd635456275dc53c8d16aa6229230804ad1f4ee4f00852a5b42220a365464736f6c63430008190033`

// TransientSlotBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TransientSlotMetaData.Bin instead.
var TransientSlotBin = TransientSlotMetaData.Bin

// DeployTransientSlot deploys a new Kaia contract, binding an instance of TransientSlot to it.
func DeployTransientSlot(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TransientSlot, error) {
	parsed, err := TransientSlotMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TransientSlotBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TransientSlot{TransientSlotCaller: TransientSlotCaller{contract: contract}, TransientSlotTransactor: TransientSlotTransactor{contract: contract}, TransientSlotFilterer: TransientSlotFilterer{contract: contract}}, nil
}

// TransientSlot is an auto generated Go binding around a Kaia contract.
type TransientSlot struct {
	TransientSlotCaller     // Read-only binding to the contract
	TransientSlotTransactor // Write-only binding to the contract
	TransientSlotFilterer   // Log filterer for contract events
}

// TransientSlotCaller is an auto generated read-only Go binding around a Kaia contract.
type TransientSlotCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransientSlotTransactor is an auto generated write-only Go binding around a Kaia contract.
type TransientSlotTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransientSlotFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type TransientSlotFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransientSlotSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type TransientSlotSession struct {
	Contract     *TransientSlot    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TransientSlotCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type TransientSlotCallerSession struct {
	Contract *TransientSlotCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TransientSlotTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type TransientSlotTransactorSession struct {
	Contract     *TransientSlotTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TransientSlotRaw is an auto generated low-level Go binding around a Kaia contract.
type TransientSlotRaw struct {
	Contract *TransientSlot // Generic contract binding to access the raw methods on
}

// TransientSlotCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type TransientSlotCallerRaw struct {
	Contract *TransientSlotCaller // Generic read-only contract binding to access the raw methods on
}

// TransientSlotTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type TransientSlotTransactorRaw struct {
	Contract *TransientSlotTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransientSlot creates a new instance of TransientSlot, bound to a specific deployed contract.
func NewTransientSlot(address common.Address, backend bind.ContractBackend) (*TransientSlot, error) {
	contract, err := bindTransientSlot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TransientSlot{TransientSlotCaller: TransientSlotCaller{contract: contract}, TransientSlotTransactor: TransientSlotTransactor{contract: contract}, TransientSlotFilterer: TransientSlotFilterer{contract: contract}}, nil
}

// NewTransientSlotCaller creates a new read-only instance of TransientSlot, bound to a specific deployed contract.
func NewTransientSlotCaller(address common.Address, caller bind.ContractCaller) (*TransientSlotCaller, error) {
	contract, err := bindTransientSlot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransientSlotCaller{contract: contract}, nil
}

// NewTransientSlotTransactor creates a new write-only instance of TransientSlot, bound to a specific deployed contract.
func NewTransientSlotTransactor(address common.Address, transactor bind.ContractTransactor) (*TransientSlotTransactor, error) {
	contract, err := bindTransientSlot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransientSlotTransactor{contract: contract}, nil
}

// NewTransientSlotFilterer creates a new log filterer instance of TransientSlot, bound to a specific deployed contract.
func NewTransientSlotFilterer(address common.Address, filterer bind.ContractFilterer) (*TransientSlotFilterer, error) {
	contract, err := bindTransientSlot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransientSlotFilterer{contract: contract}, nil
}

// bindTransientSlot binds a generic wrapper to an already deployed contract.
func bindTransientSlot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TransientSlotMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransientSlot *TransientSlotRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransientSlot.Contract.TransientSlotCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransientSlot *TransientSlotRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransientSlot.Contract.TransientSlotTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransientSlot *TransientSlotRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransientSlot.Contract.TransientSlotTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TransientSlot *TransientSlotCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TransientSlot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TransientSlot *TransientSlotTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TransientSlot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TransientSlot *TransientSlotTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TransientSlot.Contract.contract.Transact(opts, method, params...)
}
