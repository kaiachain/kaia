// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package kip113

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

// IKIP113BlsPublicKeyInfo is an auto generated low-level Go binding around an user-defined struct.
type IKIP113BlsPublicKeyInfo struct {
	PublicKey []byte
	Pop       []byte
}

// AddressUpgradeableMetaData contains all meta data concerning the AddressUpgradeable contract.
var AddressUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220c2b4e09038a93c465ca14ebc76b734be33d74eebb5af4ca36db46e6cba52808d64736f6c63430008130033",
}

// AddressUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use AddressUpgradeableMetaData.ABI instead.
var AddressUpgradeableABI = AddressUpgradeableMetaData.ABI

// AddressUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressUpgradeableBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea2646970667358221220c2b4e09038a93c465ca14ebc76b734be33d74eebb5af4ca36db46e6cba52808d64736f6c63430008130033`

// AddressUpgradeableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AddressUpgradeableMetaData.Bin instead.
var AddressUpgradeableBin = AddressUpgradeableMetaData.Bin

// DeployAddressUpgradeable deploys a new Kaia contract, binding an instance of AddressUpgradeable to it.
func DeployAddressUpgradeable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AddressUpgradeable, error) {
	parsed, err := AddressUpgradeableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AddressUpgradeableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AddressUpgradeable{AddressUpgradeableCaller: AddressUpgradeableCaller{contract: contract}, AddressUpgradeableTransactor: AddressUpgradeableTransactor{contract: contract}, AddressUpgradeableFilterer: AddressUpgradeableFilterer{contract: contract}}, nil
}

// AddressUpgradeable is an auto generated Go binding around a Kaia contract.
type AddressUpgradeable struct {
	AddressUpgradeableCaller     // Read-only binding to the contract
	AddressUpgradeableTransactor // Write-only binding to the contract
	AddressUpgradeableFilterer   // Log filterer for contract events
}

// AddressUpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type AddressUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressUpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type AddressUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressUpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type AddressUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressUpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type AddressUpgradeableSession struct {
	Contract     *AddressUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AddressUpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type AddressUpgradeableCallerSession struct {
	Contract *AddressUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// AddressUpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type AddressUpgradeableTransactorSession struct {
	Contract     *AddressUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// AddressUpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type AddressUpgradeableRaw struct {
	Contract *AddressUpgradeable // Generic contract binding to access the raw methods on
}

// AddressUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type AddressUpgradeableCallerRaw struct {
	Contract *AddressUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// AddressUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type AddressUpgradeableTransactorRaw struct {
	Contract *AddressUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAddressUpgradeable creates a new instance of AddressUpgradeable, bound to a specific deployed contract.
func NewAddressUpgradeable(address common.Address, backend bind.ContractBackend) (*AddressUpgradeable, error) {
	contract, err := bindAddressUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AddressUpgradeable{AddressUpgradeableCaller: AddressUpgradeableCaller{contract: contract}, AddressUpgradeableTransactor: AddressUpgradeableTransactor{contract: contract}, AddressUpgradeableFilterer: AddressUpgradeableFilterer{contract: contract}}, nil
}

// NewAddressUpgradeableCaller creates a new read-only instance of AddressUpgradeable, bound to a specific deployed contract.
func NewAddressUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*AddressUpgradeableCaller, error) {
	contract, err := bindAddressUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AddressUpgradeableCaller{contract: contract}, nil
}

// NewAddressUpgradeableTransactor creates a new write-only instance of AddressUpgradeable, bound to a specific deployed contract.
func NewAddressUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*AddressUpgradeableTransactor, error) {
	contract, err := bindAddressUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AddressUpgradeableTransactor{contract: contract}, nil
}

// NewAddressUpgradeableFilterer creates a new log filterer instance of AddressUpgradeable, bound to a specific deployed contract.
func NewAddressUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*AddressUpgradeableFilterer, error) {
	contract, err := bindAddressUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AddressUpgradeableFilterer{contract: contract}, nil
}

// bindAddressUpgradeable binds a generic wrapper to an already deployed contract.
func bindAddressUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AddressUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AddressUpgradeable *AddressUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AddressUpgradeable.Contract.AddressUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AddressUpgradeable *AddressUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AddressUpgradeable.Contract.AddressUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AddressUpgradeable *AddressUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AddressUpgradeable.Contract.AddressUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AddressUpgradeable *AddressUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AddressUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AddressUpgradeable *AddressUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AddressUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AddressUpgradeable *AddressUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AddressUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ContextUpgradeableMetaData contains all meta data concerning the ContextUpgradeable contract.
var ContextUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"}]",
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

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
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
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContextUpgradeable *ContextUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*ContextUpgradeableInitializedIterator, error) {
	logs, sub, err := _ContextUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ContextUpgradeableInitializedIterator{contract: _ContextUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ContextUpgradeable *ContextUpgradeableFilterer) ParseInitialized(log types.Log) (*ContextUpgradeableInitialized, error) {
	event := new(ContextUpgradeableInitialized)
	if err := _ContextUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradeableMetaData contains all meta data concerning the ERC1967UpgradeUpgradeable contract.
var ERC1967UpgradeUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"}]",
}

// ERC1967UpgradeUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1967UpgradeUpgradeableMetaData.ABI instead.
var ERC1967UpgradeUpgradeableABI = ERC1967UpgradeUpgradeableMetaData.ABI

// ERC1967UpgradeUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC1967UpgradeUpgradeableBinRuntime = ``

// ERC1967UpgradeUpgradeable is an auto generated Go binding around a Kaia contract.
type ERC1967UpgradeUpgradeable struct {
	ERC1967UpgradeUpgradeableCaller     // Read-only binding to the contract
	ERC1967UpgradeUpgradeableTransactor // Write-only binding to the contract
	ERC1967UpgradeUpgradeableFilterer   // Log filterer for contract events
}

// ERC1967UpgradeUpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type ERC1967UpgradeUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeUpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type ERC1967UpgradeUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeUpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ERC1967UpgradeUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeUpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ERC1967UpgradeUpgradeableSession struct {
	Contract     *ERC1967UpgradeUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts              // Call options to use throughout this session
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// ERC1967UpgradeUpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ERC1967UpgradeUpgradeableCallerSession struct {
	Contract *ERC1967UpgradeUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                    // Call options to use throughout this session
}

// ERC1967UpgradeUpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ERC1967UpgradeUpgradeableTransactorSession struct {
	Contract     *ERC1967UpgradeUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                    // Transaction auth options to use throughout this session
}

// ERC1967UpgradeUpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type ERC1967UpgradeUpgradeableRaw struct {
	Contract *ERC1967UpgradeUpgradeable // Generic contract binding to access the raw methods on
}

// ERC1967UpgradeUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ERC1967UpgradeUpgradeableCallerRaw struct {
	Contract *ERC1967UpgradeUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1967UpgradeUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ERC1967UpgradeUpgradeableTransactorRaw struct {
	Contract *ERC1967UpgradeUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1967UpgradeUpgradeable creates a new instance of ERC1967UpgradeUpgradeable, bound to a specific deployed contract.
func NewERC1967UpgradeUpgradeable(address common.Address, backend bind.ContractBackend) (*ERC1967UpgradeUpgradeable, error) {
	contract, err := bindERC1967UpgradeUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeable{ERC1967UpgradeUpgradeableCaller: ERC1967UpgradeUpgradeableCaller{contract: contract}, ERC1967UpgradeUpgradeableTransactor: ERC1967UpgradeUpgradeableTransactor{contract: contract}, ERC1967UpgradeUpgradeableFilterer: ERC1967UpgradeUpgradeableFilterer{contract: contract}}, nil
}

// NewERC1967UpgradeUpgradeableCaller creates a new read-only instance of ERC1967UpgradeUpgradeable, bound to a specific deployed contract.
func NewERC1967UpgradeUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*ERC1967UpgradeUpgradeableCaller, error) {
	contract, err := bindERC1967UpgradeUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableCaller{contract: contract}, nil
}

// NewERC1967UpgradeUpgradeableTransactor creates a new write-only instance of ERC1967UpgradeUpgradeable, bound to a specific deployed contract.
func NewERC1967UpgradeUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1967UpgradeUpgradeableTransactor, error) {
	contract, err := bindERC1967UpgradeUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableTransactor{contract: contract}, nil
}

// NewERC1967UpgradeUpgradeableFilterer creates a new log filterer instance of ERC1967UpgradeUpgradeable, bound to a specific deployed contract.
func NewERC1967UpgradeUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1967UpgradeUpgradeableFilterer, error) {
	contract, err := bindERC1967UpgradeUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableFilterer{contract: contract}, nil
}

// bindERC1967UpgradeUpgradeable binds a generic wrapper to an already deployed contract.
func bindERC1967UpgradeUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1967UpgradeUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967UpgradeUpgradeable.Contract.ERC1967UpgradeUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967UpgradeUpgradeable.Contract.ERC1967UpgradeUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967UpgradeUpgradeable.Contract.ERC1967UpgradeUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967UpgradeUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967UpgradeUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967UpgradeUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ERC1967UpgradeUpgradeableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableAdminChangedIterator struct {
	Event *ERC1967UpgradeUpgradeableAdminChanged // Event containing the contract specifics and raw log

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
func (it *ERC1967UpgradeUpgradeableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgradeableAdminChanged)
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
		it.Event = new(ERC1967UpgradeUpgradeableAdminChanged)
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
func (it *ERC1967UpgradeUpgradeableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradeableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgradeableAdminChanged represents a AdminChanged event raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ERC1967UpgradeUpgradeableAdminChangedIterator, error) {
	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableAdminChangedIterator{contract: _ERC1967UpgradeUpgradeable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgradeableAdminChanged) (event.Subscription, error) {
	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgradeableAdminChanged)
				if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) ParseAdminChanged(log types.Log) (*ERC1967UpgradeUpgradeableAdminChanged, error) {
	event := new(ERC1967UpgradeUpgradeableAdminChanged)
	if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradeableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableBeaconUpgradedIterator struct {
	Event *ERC1967UpgradeUpgradeableBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC1967UpgradeUpgradeableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgradeableBeaconUpgraded)
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
		it.Event = new(ERC1967UpgradeUpgradeableBeaconUpgraded)
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
func (it *ERC1967UpgradeUpgradeableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradeableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgradeableBeaconUpgraded represents a BeaconUpgraded event raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ERC1967UpgradeUpgradeableBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableBeaconUpgradedIterator{contract: _ERC1967UpgradeUpgradeable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgradeableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgradeableBeaconUpgraded)
				if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) ParseBeaconUpgraded(log types.Log) (*ERC1967UpgradeUpgradeableBeaconUpgraded, error) {
	event := new(ERC1967UpgradeUpgradeableBeaconUpgraded)
	if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableInitializedIterator struct {
	Event *ERC1967UpgradeUpgradeableInitialized // Event containing the contract specifics and raw log

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
func (it *ERC1967UpgradeUpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgradeableInitialized)
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
		it.Event = new(ERC1967UpgradeUpgradeableInitialized)
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
func (it *ERC1967UpgradeUpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgradeableInitialized represents a Initialized event raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*ERC1967UpgradeUpgradeableInitializedIterator, error) {
	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableInitializedIterator{contract: _ERC1967UpgradeUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgradeableInitialized) (event.Subscription, error) {
	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgradeableInitialized)
				if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) ParseInitialized(log types.Log) (*ERC1967UpgradeUpgradeableInitialized, error) {
	event := new(ERC1967UpgradeUpgradeableInitialized)
	if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradeableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableUpgradedIterator struct {
	Event *ERC1967UpgradeUpgradeableUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC1967UpgradeUpgradeableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgradeableUpgraded)
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
		it.Event = new(ERC1967UpgradeUpgradeableUpgraded)
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
func (it *ERC1967UpgradeUpgradeableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradeableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgradeableUpgraded represents a Upgraded event raised by the ERC1967UpgradeUpgradeable contract.
type ERC1967UpgradeUpgradeableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC1967UpgradeUpgradeableUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradeableUpgradedIterator{contract: _ERC1967UpgradeUpgradeable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgradeableUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967UpgradeUpgradeable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgradeableUpgraded)
				if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ERC1967UpgradeUpgradeable *ERC1967UpgradeUpgradeableFilterer) ParseUpgraded(log types.Log) (*ERC1967UpgradeUpgradeableUpgraded, error) {
	event := new(ERC1967UpgradeUpgradeableUpgraded)
	if err := _ERC1967UpgradeUpgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
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

// IAddressBookFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IAddressBookMetaData.Sigs instead.
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
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getAllAddress")

	outstruct := new(struct {
		TypeList    []uint8
		AddressList []common.Address
	})

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
}, error,
) {
	return _IAddressBook.Contract.GetAllAddress(&_IAddressBook.CallOpts)
}

// GetAllAddress is a free data retrieval call binding the contract method 0x715b208b.
//
// Solidity: function getAllAddress() view returns(uint8[] typeList, address[] addressList)
func (_IAddressBook *IAddressBookCallerSession) GetAllAddress() (struct {
	TypeList    []uint8
	AddressList []common.Address
}, error,
) {
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
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getAllAddressInfo")

	outstruct := new(struct {
		CnNodeIdList          []common.Address
		CnStakingContractList []common.Address
		CnRewardAddressList   []common.Address
		PocContractAddress    common.Address
		KirContractAddress    common.Address
	})

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
}, error,
) {
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
}, error,
) {
	return _IAddressBook.Contract.GetAllAddressInfo(&_IAddressBook.CallOpts)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookCaller) GetCnInfo(opts *bind.CallOpts, _cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getCnInfo", _cnNodeId)

	outstruct := new(struct {
		CnNodeId          common.Address
		CnStakingcontract common.Address
		CnRewardAddress   common.Address
	})

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
}, error,
) {
	return _IAddressBook.Contract.GetCnInfo(&_IAddressBook.CallOpts, _cnNodeId)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookCallerSession) GetCnInfo(_cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error,
) {
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
}, error,
) {
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
}, error,
) {
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
}, error,
) {
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
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getRequestInfoByArgs", _functionId, _firstArg, _secondArg, _thirdArg)

	outstruct := new(struct {
		Id                  [32]byte
		Confirmers          []common.Address
		InitialProposedTime *big.Int
		State               uint8
	})

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
}, error,
) {
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
}, error,
) {
	return _IAddressBook.Contract.GetRequestInfoByArgs(&_IAddressBook.CallOpts, _functionId, _firstArg, _secondArg, _thirdArg)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 requirement)
func (_IAddressBook *IAddressBookCaller) GetState(opts *bind.CallOpts) (struct {
	AdminList   []common.Address
	Requirement *big.Int
}, error,
) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getState")

	outstruct := new(struct {
		AdminList   []common.Address
		Requirement *big.Int
	})

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
}, error,
) {
	return _IAddressBook.Contract.GetState(&_IAddressBook.CallOpts)
}

// GetState is a free data retrieval call binding the contract method 0x1865c57d.
//
// Solidity: function getState() view returns(address[] adminList, uint256 requirement)
func (_IAddressBook *IAddressBookCallerSession) GetState() (struct {
	AdminList   []common.Address
	Requirement *big.Int
}, error,
) {
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

// IBeaconUpgradeableMetaData contains all meta data concerning the IBeaconUpgradeable contract.
var IBeaconUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"5c60da1b": "implementation()",
	},
}

// IBeaconUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use IBeaconUpgradeableMetaData.ABI instead.
var IBeaconUpgradeableABI = IBeaconUpgradeableMetaData.ABI

// IBeaconUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IBeaconUpgradeableBinRuntime = ``

// IBeaconUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IBeaconUpgradeableMetaData.Sigs instead.
var IBeaconUpgradeableFuncSigs = IBeaconUpgradeableMetaData.Sigs

// IBeaconUpgradeable is an auto generated Go binding around a Kaia contract.
type IBeaconUpgradeable struct {
	IBeaconUpgradeableCaller     // Read-only binding to the contract
	IBeaconUpgradeableTransactor // Write-only binding to the contract
	IBeaconUpgradeableFilterer   // Log filterer for contract events
}

// IBeaconUpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type IBeaconUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconUpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type IBeaconUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconUpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IBeaconUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconUpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IBeaconUpgradeableSession struct {
	Contract     *IBeaconUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IBeaconUpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IBeaconUpgradeableCallerSession struct {
	Contract *IBeaconUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// IBeaconUpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IBeaconUpgradeableTransactorSession struct {
	Contract     *IBeaconUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IBeaconUpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type IBeaconUpgradeableRaw struct {
	Contract *IBeaconUpgradeable // Generic contract binding to access the raw methods on
}

// IBeaconUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IBeaconUpgradeableCallerRaw struct {
	Contract *IBeaconUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// IBeaconUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IBeaconUpgradeableTransactorRaw struct {
	Contract *IBeaconUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBeaconUpgradeable creates a new instance of IBeaconUpgradeable, bound to a specific deployed contract.
func NewIBeaconUpgradeable(address common.Address, backend bind.ContractBackend) (*IBeaconUpgradeable, error) {
	contract, err := bindIBeaconUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBeaconUpgradeable{IBeaconUpgradeableCaller: IBeaconUpgradeableCaller{contract: contract}, IBeaconUpgradeableTransactor: IBeaconUpgradeableTransactor{contract: contract}, IBeaconUpgradeableFilterer: IBeaconUpgradeableFilterer{contract: contract}}, nil
}

// NewIBeaconUpgradeableCaller creates a new read-only instance of IBeaconUpgradeable, bound to a specific deployed contract.
func NewIBeaconUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*IBeaconUpgradeableCaller, error) {
	contract, err := bindIBeaconUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBeaconUpgradeableCaller{contract: contract}, nil
}

// NewIBeaconUpgradeableTransactor creates a new write-only instance of IBeaconUpgradeable, bound to a specific deployed contract.
func NewIBeaconUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*IBeaconUpgradeableTransactor, error) {
	contract, err := bindIBeaconUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBeaconUpgradeableTransactor{contract: contract}, nil
}

// NewIBeaconUpgradeableFilterer creates a new log filterer instance of IBeaconUpgradeable, bound to a specific deployed contract.
func NewIBeaconUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*IBeaconUpgradeableFilterer, error) {
	contract, err := bindIBeaconUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBeaconUpgradeableFilterer{contract: contract}, nil
}

// bindIBeaconUpgradeable binds a generic wrapper to an already deployed contract.
func bindIBeaconUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBeaconUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBeaconUpgradeable *IBeaconUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBeaconUpgradeable.Contract.IBeaconUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBeaconUpgradeable *IBeaconUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBeaconUpgradeable.Contract.IBeaconUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBeaconUpgradeable *IBeaconUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBeaconUpgradeable.Contract.IBeaconUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBeaconUpgradeable *IBeaconUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBeaconUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBeaconUpgradeable *IBeaconUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBeaconUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBeaconUpgradeable *IBeaconUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBeaconUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeaconUpgradeable *IBeaconUpgradeableCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IBeaconUpgradeable.contract.Call(opts, &out, "implementation")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeaconUpgradeable *IBeaconUpgradeableSession) Implementation() (common.Address, error) {
	return _IBeaconUpgradeable.Contract.Implementation(&_IBeaconUpgradeable.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeaconUpgradeable *IBeaconUpgradeableCallerSession) Implementation() (common.Address, error) {
	return _IBeaconUpgradeable.Contract.Implementation(&_IBeaconUpgradeable.CallOpts)
}

// IERC1822ProxiableUpgradeableMetaData contains all meta data concerning the IERC1822ProxiableUpgradeable contract.
var IERC1822ProxiableUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"52d1902d": "proxiableUUID()",
	},
}

// IERC1822ProxiableUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1822ProxiableUpgradeableMetaData.ABI instead.
var IERC1822ProxiableUpgradeableABI = IERC1822ProxiableUpgradeableMetaData.ABI

// IERC1822ProxiableUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1822ProxiableUpgradeableBinRuntime = ``

// IERC1822ProxiableUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IERC1822ProxiableUpgradeableMetaData.Sigs instead.
var IERC1822ProxiableUpgradeableFuncSigs = IERC1822ProxiableUpgradeableMetaData.Sigs

// IERC1822ProxiableUpgradeable is an auto generated Go binding around a Kaia contract.
type IERC1822ProxiableUpgradeable struct {
	IERC1822ProxiableUpgradeableCaller     // Read-only binding to the contract
	IERC1822ProxiableUpgradeableTransactor // Write-only binding to the contract
	IERC1822ProxiableUpgradeableFilterer   // Log filterer for contract events
}

// IERC1822ProxiableUpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type IERC1822ProxiableUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableUpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type IERC1822ProxiableUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableUpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IERC1822ProxiableUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableUpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IERC1822ProxiableUpgradeableSession struct {
	Contract     *IERC1822ProxiableUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IERC1822ProxiableUpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IERC1822ProxiableUpgradeableCallerSession struct {
	Contract *IERC1822ProxiableUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// IERC1822ProxiableUpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IERC1822ProxiableUpgradeableTransactorSession struct {
	Contract     *IERC1822ProxiableUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// IERC1822ProxiableUpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type IERC1822ProxiableUpgradeableRaw struct {
	Contract *IERC1822ProxiableUpgradeable // Generic contract binding to access the raw methods on
}

// IERC1822ProxiableUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IERC1822ProxiableUpgradeableCallerRaw struct {
	Contract *IERC1822ProxiableUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// IERC1822ProxiableUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IERC1822ProxiableUpgradeableTransactorRaw struct {
	Contract *IERC1822ProxiableUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1822ProxiableUpgradeable creates a new instance of IERC1822ProxiableUpgradeable, bound to a specific deployed contract.
func NewIERC1822ProxiableUpgradeable(address common.Address, backend bind.ContractBackend) (*IERC1822ProxiableUpgradeable, error) {
	contract, err := bindIERC1822ProxiableUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableUpgradeable{IERC1822ProxiableUpgradeableCaller: IERC1822ProxiableUpgradeableCaller{contract: contract}, IERC1822ProxiableUpgradeableTransactor: IERC1822ProxiableUpgradeableTransactor{contract: contract}, IERC1822ProxiableUpgradeableFilterer: IERC1822ProxiableUpgradeableFilterer{contract: contract}}, nil
}

// NewIERC1822ProxiableUpgradeableCaller creates a new read-only instance of IERC1822ProxiableUpgradeable, bound to a specific deployed contract.
func NewIERC1822ProxiableUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*IERC1822ProxiableUpgradeableCaller, error) {
	contract, err := bindIERC1822ProxiableUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableUpgradeableCaller{contract: contract}, nil
}

// NewIERC1822ProxiableUpgradeableTransactor creates a new write-only instance of IERC1822ProxiableUpgradeable, bound to a specific deployed contract.
func NewIERC1822ProxiableUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC1822ProxiableUpgradeableTransactor, error) {
	contract, err := bindIERC1822ProxiableUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableUpgradeableTransactor{contract: contract}, nil
}

// NewIERC1822ProxiableUpgradeableFilterer creates a new log filterer instance of IERC1822ProxiableUpgradeable, bound to a specific deployed contract.
func NewIERC1822ProxiableUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC1822ProxiableUpgradeableFilterer, error) {
	contract, err := bindIERC1822ProxiableUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableUpgradeableFilterer{contract: contract}, nil
}

// bindIERC1822ProxiableUpgradeable binds a generic wrapper to an already deployed contract.
func bindIERC1822ProxiableUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1822ProxiableUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1822ProxiableUpgradeable.Contract.IERC1822ProxiableUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1822ProxiableUpgradeable.Contract.IERC1822ProxiableUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1822ProxiableUpgradeable.Contract.IERC1822ProxiableUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1822ProxiableUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1822ProxiableUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1822ProxiableUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IERC1822ProxiableUpgradeable.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableSession) ProxiableUUID() ([32]byte, error) {
	return _IERC1822ProxiableUpgradeable.Contract.ProxiableUUID(&_IERC1822ProxiableUpgradeable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822ProxiableUpgradeable *IERC1822ProxiableUpgradeableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _IERC1822ProxiableUpgradeable.Contract.ProxiableUUID(&_IERC1822ProxiableUpgradeable.CallOpts)
}

// IERC1967UpgradeableMetaData contains all meta data concerning the IERC1967Upgradeable contract.
var IERC1967UpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"}]",
}

// IERC1967UpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1967UpgradeableMetaData.ABI instead.
var IERC1967UpgradeableABI = IERC1967UpgradeableMetaData.ABI

// IERC1967UpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1967UpgradeableBinRuntime = ``

// IERC1967Upgradeable is an auto generated Go binding around a Kaia contract.
type IERC1967Upgradeable struct {
	IERC1967UpgradeableCaller     // Read-only binding to the contract
	IERC1967UpgradeableTransactor // Write-only binding to the contract
	IERC1967UpgradeableFilterer   // Log filterer for contract events
}

// IERC1967UpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type IERC1967UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967UpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type IERC1967UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967UpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IERC1967UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967UpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IERC1967UpgradeableSession struct {
	Contract     *IERC1967Upgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IERC1967UpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IERC1967UpgradeableCallerSession struct {
	Contract *IERC1967UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// IERC1967UpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IERC1967UpgradeableTransactorSession struct {
	Contract     *IERC1967UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// IERC1967UpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type IERC1967UpgradeableRaw struct {
	Contract *IERC1967Upgradeable // Generic contract binding to access the raw methods on
}

// IERC1967UpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IERC1967UpgradeableCallerRaw struct {
	Contract *IERC1967UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// IERC1967UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IERC1967UpgradeableTransactorRaw struct {
	Contract *IERC1967UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1967Upgradeable creates a new instance of IERC1967Upgradeable, bound to a specific deployed contract.
func NewIERC1967Upgradeable(address common.Address, backend bind.ContractBackend) (*IERC1967Upgradeable, error) {
	contract, err := bindIERC1967Upgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1967Upgradeable{IERC1967UpgradeableCaller: IERC1967UpgradeableCaller{contract: contract}, IERC1967UpgradeableTransactor: IERC1967UpgradeableTransactor{contract: contract}, IERC1967UpgradeableFilterer: IERC1967UpgradeableFilterer{contract: contract}}, nil
}

// NewIERC1967UpgradeableCaller creates a new read-only instance of IERC1967Upgradeable, bound to a specific deployed contract.
func NewIERC1967UpgradeableCaller(address common.Address, caller bind.ContractCaller) (*IERC1967UpgradeableCaller, error) {
	contract, err := bindIERC1967Upgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableCaller{contract: contract}, nil
}

// NewIERC1967UpgradeableTransactor creates a new write-only instance of IERC1967Upgradeable, bound to a specific deployed contract.
func NewIERC1967UpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC1967UpgradeableTransactor, error) {
	contract, err := bindIERC1967Upgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableTransactor{contract: contract}, nil
}

// NewIERC1967UpgradeableFilterer creates a new log filterer instance of IERC1967Upgradeable, bound to a specific deployed contract.
func NewIERC1967UpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC1967UpgradeableFilterer, error) {
	contract, err := bindIERC1967Upgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableFilterer{contract: contract}, nil
}

// bindIERC1967Upgradeable binds a generic wrapper to an already deployed contract.
func bindIERC1967Upgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1967UpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1967Upgradeable *IERC1967UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1967Upgradeable.Contract.IERC1967UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1967Upgradeable *IERC1967UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1967Upgradeable.Contract.IERC1967UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1967Upgradeable *IERC1967UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1967Upgradeable.Contract.IERC1967UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1967Upgradeable *IERC1967UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1967Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1967Upgradeable *IERC1967UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1967Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1967Upgradeable *IERC1967UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1967Upgradeable.Contract.contract.Transact(opts, method, params...)
}

// IERC1967UpgradeableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableAdminChangedIterator struct {
	Event *IERC1967UpgradeableAdminChanged // Event containing the contract specifics and raw log

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
func (it *IERC1967UpgradeableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967UpgradeableAdminChanged)
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
		it.Event = new(IERC1967UpgradeableAdminChanged)
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
func (it *IERC1967UpgradeableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967UpgradeableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967UpgradeableAdminChanged represents a AdminChanged event raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*IERC1967UpgradeableAdminChangedIterator, error) {
	logs, sub, err := _IERC1967Upgradeable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableAdminChangedIterator{contract: _IERC1967Upgradeable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *IERC1967UpgradeableAdminChanged) (event.Subscription, error) {
	logs, sub, err := _IERC1967Upgradeable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967UpgradeableAdminChanged)
				if err := _IERC1967Upgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) ParseAdminChanged(log types.Log) (*IERC1967UpgradeableAdminChanged, error) {
	event := new(IERC1967UpgradeableAdminChanged)
	if err := _IERC1967Upgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC1967UpgradeableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableBeaconUpgradedIterator struct {
	Event *IERC1967UpgradeableBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *IERC1967UpgradeableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967UpgradeableBeaconUpgraded)
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
		it.Event = new(IERC1967UpgradeableBeaconUpgraded)
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
func (it *IERC1967UpgradeableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967UpgradeableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967UpgradeableBeaconUpgraded represents a BeaconUpgraded event raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*IERC1967UpgradeableBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _IERC1967Upgradeable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableBeaconUpgradedIterator{contract: _IERC1967Upgradeable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *IERC1967UpgradeableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _IERC1967Upgradeable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967UpgradeableBeaconUpgraded)
				if err := _IERC1967Upgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) ParseBeaconUpgraded(log types.Log) (*IERC1967UpgradeableBeaconUpgraded, error) {
	event := new(IERC1967UpgradeableBeaconUpgraded)
	if err := _IERC1967Upgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IERC1967UpgradeableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableUpgradedIterator struct {
	Event *IERC1967UpgradeableUpgraded // Event containing the contract specifics and raw log

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
func (it *IERC1967UpgradeableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967UpgradeableUpgraded)
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
		it.Event = new(IERC1967UpgradeableUpgraded)
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
func (it *IERC1967UpgradeableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967UpgradeableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967UpgradeableUpgraded represents a Upgraded event raised by the IERC1967Upgradeable contract.
type IERC1967UpgradeableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*IERC1967UpgradeableUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IERC1967Upgradeable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradeableUpgradedIterator{contract: _IERC1967Upgradeable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *IERC1967UpgradeableUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IERC1967Upgradeable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967UpgradeableUpgraded)
				if err := _IERC1967Upgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_IERC1967Upgradeable *IERC1967UpgradeableFilterer) ParseUpgraded(log types.Log) (*IERC1967UpgradeableUpgraded, error) {
	event := new(IERC1967UpgradeableUpgraded)
	if err := _IERC1967Upgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// IKIP113MetaData contains all meta data concerning the IKIP113 contract.
var IKIP113MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllBlsInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIdList\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structIKIP113.BlsPublicKeyInfo[]\",\"name\":\"pubkeyList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"6968b53f": "getAllBlsInfo()",
	},
}

// IKIP113ABI is the input ABI used to generate the binding from.
// Deprecated: Use IKIP113MetaData.ABI instead.
var IKIP113ABI = IKIP113MetaData.ABI

// IKIP113BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IKIP113BinRuntime = ``

// IKIP113FuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IKIP113MetaData.Sigs instead.
var IKIP113FuncSigs = IKIP113MetaData.Sigs

// IKIP113 is an auto generated Go binding around a Kaia contract.
type IKIP113 struct {
	IKIP113Caller     // Read-only binding to the contract
	IKIP113Transactor // Write-only binding to the contract
	IKIP113Filterer   // Log filterer for contract events
}

// IKIP113Caller is an auto generated read-only Go binding around a Kaia contract.
type IKIP113Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP113Transactor is an auto generated write-only Go binding around a Kaia contract.
type IKIP113Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP113Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type IKIP113Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IKIP113Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IKIP113Session struct {
	Contract     *IKIP113          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IKIP113CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IKIP113CallerSession struct {
	Contract *IKIP113Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IKIP113TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IKIP113TransactorSession struct {
	Contract     *IKIP113Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IKIP113Raw is an auto generated low-level Go binding around a Kaia contract.
type IKIP113Raw struct {
	Contract *IKIP113 // Generic contract binding to access the raw methods on
}

// IKIP113CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IKIP113CallerRaw struct {
	Contract *IKIP113Caller // Generic read-only contract binding to access the raw methods on
}

// IKIP113TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IKIP113TransactorRaw struct {
	Contract *IKIP113Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIKIP113 creates a new instance of IKIP113, bound to a specific deployed contract.
func NewIKIP113(address common.Address, backend bind.ContractBackend) (*IKIP113, error) {
	contract, err := bindIKIP113(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IKIP113{IKIP113Caller: IKIP113Caller{contract: contract}, IKIP113Transactor: IKIP113Transactor{contract: contract}, IKIP113Filterer: IKIP113Filterer{contract: contract}}, nil
}

// NewIKIP113Caller creates a new read-only instance of IKIP113, bound to a specific deployed contract.
func NewIKIP113Caller(address common.Address, caller bind.ContractCaller) (*IKIP113Caller, error) {
	contract, err := bindIKIP113(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IKIP113Caller{contract: contract}, nil
}

// NewIKIP113Transactor creates a new write-only instance of IKIP113, bound to a specific deployed contract.
func NewIKIP113Transactor(address common.Address, transactor bind.ContractTransactor) (*IKIP113Transactor, error) {
	contract, err := bindIKIP113(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IKIP113Transactor{contract: contract}, nil
}

// NewIKIP113Filterer creates a new log filterer instance of IKIP113, bound to a specific deployed contract.
func NewIKIP113Filterer(address common.Address, filterer bind.ContractFilterer) (*IKIP113Filterer, error) {
	contract, err := bindIKIP113(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IKIP113Filterer{contract: contract}, nil
}

// bindIKIP113 binds a generic wrapper to an already deployed contract.
func bindIKIP113(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IKIP113MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IKIP113 *IKIP113Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKIP113.Contract.IKIP113Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IKIP113 *IKIP113Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKIP113.Contract.IKIP113Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IKIP113 *IKIP113Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKIP113.Contract.IKIP113Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IKIP113 *IKIP113CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IKIP113.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IKIP113 *IKIP113TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IKIP113.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IKIP113 *IKIP113TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IKIP113.Contract.contract.Transact(opts, method, params...)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_IKIP113 *IKIP113Caller) GetAllBlsInfo(opts *bind.CallOpts) (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	var out []interface{}
	err := _IKIP113.contract.Call(opts, &out, "getAllBlsInfo")

	outstruct := new(struct {
		NodeIdList []common.Address
		PubkeyList []IKIP113BlsPublicKeyInfo
	})

	outstruct.NodeIdList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.PubkeyList = *abi.ConvertType(out[1], new([]IKIP113BlsPublicKeyInfo)).(*[]IKIP113BlsPublicKeyInfo)
	return *outstruct, err
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_IKIP113 *IKIP113Session) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _IKIP113.Contract.GetAllBlsInfo(&_IKIP113.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_IKIP113 *IKIP113CallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _IKIP113.Contract.GetAllBlsInfo(&_IKIP113.CallOpts)
}

// InitializableMetaData contains all meta data concerning the Initializable contract.
var InitializableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"}]",
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

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
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
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Initializable *InitializableFilterer) FilterInitialized(opts *bind.FilterOpts) (*InitializableInitializedIterator, error) {
	logs, sub, err := _Initializable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &InitializableInitializedIterator{contract: _Initializable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Initializable *InitializableFilterer) ParseInitialized(log types.Log) (*InitializableInitialized, error) {
	event := new(InitializableInitialized)
	if err := _Initializable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OwnableUpgradeableMetaData contains all meta data concerning the OwnableUpgradeable contract.
var OwnableUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// OwnableUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use OwnableUpgradeableMetaData.Sigs instead.
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

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
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
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*OwnableUpgradeableInitializedIterator, error) {
	logs, sub, err := _OwnableUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &OwnableUpgradeableInitializedIterator{contract: _OwnableUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_OwnableUpgradeable *OwnableUpgradeableFilterer) ParseInitialized(log types.Log) (*OwnableUpgradeableInitialized, error) {
	event := new(OwnableUpgradeableInitialized)
	if err := _OwnableUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OwnableUpgradeableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the OwnableUpgradeable contract.
type OwnableUpgradeableOwnershipTransferredIterator struct {
	Event *OwnableUpgradeableOwnershipTransferred // Event containing the contract specifics and raw log

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
	return event, nil
}

// SimpleBlsRegistryMetaData contains all meta data concerning the SimpleBlsRegistry contract.
var SimpleBlsRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"Unregistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ZERO48HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ZERO96HASH\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"abook\",\"outputs\":[{\"internalType\":\"contractIAddressBook\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allNodeIds\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllBlsInfo\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"nodeIdList\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structIKIP113.BlsPublicKeyInfo[]\",\"name\":\"pubkeyList\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"record\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"}],\"name\":\"unregister\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"6fc522c6": "ZERO48HASH()",
		"20abd458": "ZERO96HASH()",
		"829d639d": "abook()",
		"a5834971": "allNodeIds(uint256)",
		"6968b53f": "getAllBlsInfo()",
		"8129fc1c": "initialize()",
		"8da5cb5b": "owner()",
		"52d1902d": "proxiableUUID()",
		"3465d6d5": "record(address)",
		"786cd4d7": "register(address,bytes,bytes)",
		"715018a6": "renounceOwnership()",
		"f2fde38b": "transferOwnership(address)",
		"2ec2c246": "unregister(address)",
		"3659cfe6": "upgradeTo(address)",
		"4f1ef286": "upgradeToAndCall(address,bytes)",
	},
	Bin: "0x60a06040523060805234801561001457600080fd5b5061001d610022565b6100e1565b600054610100900460ff161561008e5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b60005460ff908116146100df576000805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b608051611ecc61011860003960008181610593015281816105d301528181610672015281816106b201526107450152611ecc6000f3fe6080604052600436106100e85760003560e01c80636fc522c61161008a578063829d639d11610059578063829d639d1461026d5780638da5cb5b1461029b578063a5834971146102b9578063f2fde38b146102d957600080fd5b80636fc522c6146101ef578063715018a614610223578063786cd4d7146102385780638129fc1c1461025857600080fd5b80633659cfe6116100c65780633659cfe6146101845780634f1ef286146101a457806352d1902d146101b75780636968b53f146101cc57600080fd5b806320abd458146100ed5780632ec2c246146101345780633465d6d514610156575b600080fd5b3480156100f957600080fd5b506101217f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2181565b6040519081526020015b60405180910390f35b34801561014057600080fd5b5061015461014f3660046116cb565b6102f9565b005b34801561016257600080fd5b506101766101713660046116cb565b61045d565b60405161012b92919061173f565b34801561019057600080fd5b5061015461019f3660046116cb565b610589565b6101546101b2366004611783565b610668565b3480156101c357600080fd5b50610121610738565b3480156101d857600080fd5b506101e16107eb565b60405161012b929190611847565b3480156101fb57600080fd5b506101217fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29381565b34801561022f57600080fd5b50610154610aa8565b34801561024457600080fd5b50610154610253366004611955565b610abc565b34801561026457600080fd5b50610154610e30565b34801561027957600080fd5b5061028361040081565b6040516001600160a01b03909116815260200161012b565b3480156102a757600080fd5b506097546001600160a01b0316610283565b3480156102c557600080fd5b506102836102d43660046119d8565b610f48565b3480156102e557600080fd5b506101546102f43660046116cb565b610f72565b610301610fe8565b61030a81611042565b1561035c5760405162461bcd60e51b815260206004820152601a60248201527f434e206973207374696c6c20696e2041646472657373426f6f6b00000000000060448201526064015b60405180910390fd5b6001600160a01b038116600090815260ca60205260409020805461037f906119f1565b90506000036103c75760405162461bcd60e51b815260206004820152601460248201527310d3881a5cc81b9bdd081c9959da5cdd195c995960621b6044820152606401610353565b6103d0816110be565b6001600160a01b038116600090815260ca60205260409081902090517fb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d99161041f918491906001820190611aa8565b60405180910390a16001600160a01b038116600090815260ca602052604081209061044a8282611668565b610458600183016000611668565b505050565b60ca60205260009081526040902080548190610478906119f1565b80601f01602080910402602001604051908101604052809291908181526020018280546104a4906119f1565b80156104f15780601f106104c6576101008083540402835291602001916104f1565b820191906000526020600020905b8154815290600101906020018083116104d457829003601f168201915b505050505090806001018054610506906119f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610532906119f1565b801561057f5780601f106105545761010080835404028352916020019161057f565b820191906000526020600020905b81548152906001019060200180831161056257829003601f168201915b5050505050905082565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036105d15760405162461bcd60e51b815260040161035390611ade565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661061a600080516020611e50833981519152546001600160a01b031690565b6001600160a01b0316146106405760405162461bcd60e51b815260040161035390611b2a565b610649816111c5565b60408051600080825260208201909252610665918391906111cd565b50565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036106b05760405162461bcd60e51b815260040161035390611ade565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106f9600080516020611e50833981519152546001600160a01b031690565b6001600160a01b03161461071f5760405162461bcd60e51b815260040161035390611b2a565b610728826111c5565b610734828260016111cd565b5050565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107d85760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610353565b50600080516020611e5083398151915290565b60c954606090819067ffffffffffffffff81111561080b5761080b61176d565b604051908082528060200260200182016040528015610834578160200160208202803683370190505b5060c95490925067ffffffffffffffff8111156108535761085361176d565b60405190808252806020026020018201604052801561089857816020015b60408051808201909152606080825260208201528152602001906001900390816108715790505b50905060005b8251811015610aa35760c981815481106108ba576108ba611b76565b9060005260206000200160009054906101000a90046001600160a01b03168382815181106108ea576108ea611b76565b60200260200101906001600160a01b031690816001600160a01b03168152505060ca600060c9838154811061092157610921611b76565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120815180830190925280548290829061095f906119f1565b80601f016020809104026020016040519081016040528092919081815260200182805461098b906119f1565b80156109d85780601f106109ad576101008083540402835291602001916109d8565b820191906000526020600020905b8154815290600101906020018083116109bb57829003601f168201915b505050505081526020016001820180546109f1906119f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610a1d906119f1565b8015610a6a5780601f10610a3f57610100808354040283529160200191610a6a565b820191906000526020600020905b815481529060010190602001808311610a4d57829003601f168201915b505050505081525050828281518110610a8557610a85611b76565b60200260200101819052508080610a9b90611ba2565b91505061089e565b509091565b610ab0610fe8565b610aba6000611338565b565b610ac4610fe8565b838360308114610b165760405162461bcd60e51b815260206004820152601b60248201527f5075626c6963206b6579206d75737420626520343820627974657300000000006044820152606401610353565b6040517fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29390610b489084908490611bbb565b604051809103902003610b9d5760405162461bcd60e51b815260206004820152601960248201527f5075626c6963206b65792063616e6e6f74206265207a65726f000000000000006044820152606401610353565b838360608114610be65760405162461bcd60e51b8152602060048201526014602482015273506f70206d75737420626520393620627974657360601b6044820152606401610353565b6040517f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2190610c189084908490611bbb565b604051809103902003610c625760405162461bcd60e51b8152602060048201526012602482015271506f702063616e6e6f74206265207a65726f60701b6044820152606401610353565b610c6b89611042565b610cb75760405162461bcd60e51b815260206004820152601e60248201527f636e4e6f64654964206973206e6f7420696e2041646472657373426f6f6b00006044820152606401610353565b6001600160a01b038916600090815260ca602052604090208054610cda906119f1565b9050600003610d2f5760c980546001810182556000919091527f66be4f155c5ef2ebd3772b228f2f00681e4ed5826cdb3b1943cc11ad15ad1d280180546001600160a01b0319166001600160a01b038b161790555b6040805160606020601f8b018190040282018101835291810189815290918291908b908b9081908501838280828437600092019190915250505090825250604080516020601f8a018190048102820181019092528881529181019190899089908190840183828082843760009201829052509390945250506001600160a01b038c16815260ca6020526040902082519091508190610dcd9082611c19565b5060208201516001820190610de29082611c19565b509050507f79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f8989898989604051610e1d959493929190611d02565b60405180910390a1505050505050505050565b600054610100900460ff1615808015610e505750600054600160ff909116105b80610e6a5750303b158015610e6a575060005460ff166001145b610ecd5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610353565b6000805460ff191660011790558015610ef0576000805461ff0019166101001790555b610ef861138a565b610f006113b9565b8015610665576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b60c98181548110610f5857600080fd5b6000918252602090912001546001600160a01b0316905081565b610f7a610fe8565b6001600160a01b038116610fdf5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610353565b61066581611338565b6097546001600160a01b03163314610aba5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610353565b604051630aabaead60e11b81526001600160a01b0382166004820152600090610400906315575d5a90602401606060405180830381865afa9250505080156110a7575060408051601f3d908101601f191682019092526110a491810190611d46565b60015b6110b357506000919050565b506001949350505050565b60005b60c95481101561073457816001600160a01b031660c982815481106110e8576110e8611b76565b6000918252602090912001546001600160a01b0316036111b35760c9805461111290600190611d93565b8154811061112257611122611b76565b60009182526020909120015460c980546001600160a01b03909216918390811061114e5761114e611b76565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060c980548061118d5761118d611da6565b600082815260209020810160001990810180546001600160a01b03191690550190555050565b806111bd81611ba2565b9150506110c1565b610665610fe8565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff161561120057610458836113e0565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561125a575060408051601f3d908101601f1916820190925261125791810190611dbc565b60015b6112bd5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610353565b600080516020611e50833981519152811461132c5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610353565b5061045883838361147c565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166113b15760405162461bcd60e51b815260040161035390611dd5565b610aba6114a7565b600054610100900460ff16610aba5760405162461bcd60e51b815260040161035390611dd5565b6001600160a01b0381163b61144d5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610353565b600080516020611e5083398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611485836114d7565b6000825111806114925750805b15610458576114a18383611517565b50505050565b600054610100900460ff166114ce5760405162461bcd60e51b815260040161035390611dd5565b610aba33611338565b6114e0816113e0565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061153c8383604051806060016040528060278152602001611e7060279139611545565b90505b92915050565b6060600080856001600160a01b0316856040516115629190611e20565b600060405180830381855af49150503d806000811461159d576040519150601f19603f3d011682016040523d82523d6000602084013e6115a2565b606091505b50915091506115b3868383876115bd565b9695505050505050565b6060831561162c578251600003611625576001600160a01b0385163b6116255760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610353565b5081611636565b611636838361163e565b949350505050565b81511561164e5781518083602001fd5b8060405162461bcd60e51b81526004016103539190611e3c565b508054611674906119f1565b6000825580601f10611684575050565b601f01602090049060005260206000209081019061066591905b808211156116b2576000815560010161169e565b5090565b6001600160a01b038116811461066557600080fd5b6000602082840312156116dd57600080fd5b81356116e8816116b6565b9392505050565b60005b8381101561170a5781810151838201526020016116f2565b50506000910152565b6000815180845261172b8160208601602086016116ef565b601f01601f19169290920160200192915050565b6040815260006117526040830185611713565b82810360208401526117648185611713565b95945050505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561179657600080fd5b82356117a1816116b6565b9150602083013567ffffffffffffffff808211156117be57600080fd5b818501915085601f8301126117d257600080fd5b8135818111156117e4576117e461176d565b604051601f8201601f19908116603f0116810190838211818310171561180c5761180c61176d565b8160405282815288602084870101111561182557600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60408082528351828201819052600091906020906060850190828801855b8281101561188a5781516001600160a01b031684529284019290840190600101611865565b50505084810382860152855180825282820190600581901b8301840188850160005b838110156118fc57858303601f19018552815180518985526118d08a860182611713565b91890151858303868b01529190506118e88183611713565b9689019694505050908601906001016118ac565b50909a9950505050505050505050565b60008083601f84011261191e57600080fd5b50813567ffffffffffffffff81111561193657600080fd5b60208301915083602082850101111561194e57600080fd5b9250929050565b60008060008060006060868803121561196d57600080fd5b8535611978816116b6565b9450602086013567ffffffffffffffff8082111561199557600080fd5b6119a189838a0161190c565b909650945060408801359150808211156119ba57600080fd5b506119c78882890161190c565b969995985093965092949392505050565b6000602082840312156119ea57600080fd5b5035919050565b600181811c90821680611a0557607f821691505b602082108103611a2557634e487b7160e01b600052602260045260246000fd5b50919050565b60008154611a38816119f1565b808552602060018381168015611a555760018114611a6f57611a9d565b60ff1985168884015283151560051b880183019550611a9d565b866000528260002060005b85811015611a955781548a8201860152908301908401611a7a565b890184019650505b505050505092915050565b6001600160a01b0384168152606060208201819052600090611acc90830185611a2b565b82810360408401526115b38185611a2b565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201611bb457611bb4611b8c565b5060010190565b8183823760009101908152919050565b601f82111561045857600081815260208120601f850160051c81016020861015611bf25750805b601f850160051c820191505b81811015611c1157828155600101611bfe565b505050505050565b815167ffffffffffffffff811115611c3357611c3361176d565b611c4781611c4184546119f1565b84611bcb565b602080601f831160018114611c7c5760008415611c645750858301515b600019600386901b1c1916600185901b178555611c11565b600085815260208120601f198616915b82811015611cab57888601518255948401946001909101908401611c8c565b5085821015611cc95787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0386168152606060208201819052600090611d279083018688611cd9565b8281036040840152611d3a818587611cd9565b98975050505050505050565b600080600060608486031215611d5b57600080fd5b8351611d66816116b6565b6020850151909350611d77816116b6565b6040850151909250611d88816116b6565b809150509250925092565b8181038181111561153f5761153f611b8c565b634e487b7160e01b600052603160045260246000fd5b600060208284031215611dce57600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008251611e328184602087016116ef565b9190910192915050565b60208152600061153c602083018461171356fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220cf3c282151123924c9c8275c323310bbf7c513b7905cf4ab928cb0d42f59f3a664736f6c63430008130033",
}

// SimpleBlsRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use SimpleBlsRegistryMetaData.ABI instead.
var SimpleBlsRegistryABI = SimpleBlsRegistryMetaData.ABI

// SimpleBlsRegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const SimpleBlsRegistryBinRuntime = `6080604052600436106100e85760003560e01c80636fc522c61161008a578063829d639d11610059578063829d639d1461026d5780638da5cb5b1461029b578063a5834971146102b9578063f2fde38b146102d957600080fd5b80636fc522c6146101ef578063715018a614610223578063786cd4d7146102385780638129fc1c1461025857600080fd5b80633659cfe6116100c65780633659cfe6146101845780634f1ef286146101a457806352d1902d146101b75780636968b53f146101cc57600080fd5b806320abd458146100ed5780632ec2c246146101345780633465d6d514610156575b600080fd5b3480156100f957600080fd5b506101217f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2181565b6040519081526020015b60405180910390f35b34801561014057600080fd5b5061015461014f3660046116cb565b6102f9565b005b34801561016257600080fd5b506101766101713660046116cb565b61045d565b60405161012b92919061173f565b34801561019057600080fd5b5061015461019f3660046116cb565b610589565b6101546101b2366004611783565b610668565b3480156101c357600080fd5b50610121610738565b3480156101d857600080fd5b506101e16107eb565b60405161012b929190611847565b3480156101fb57600080fd5b506101217fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29381565b34801561022f57600080fd5b50610154610aa8565b34801561024457600080fd5b50610154610253366004611955565b610abc565b34801561026457600080fd5b50610154610e30565b34801561027957600080fd5b5061028361040081565b6040516001600160a01b03909116815260200161012b565b3480156102a757600080fd5b506097546001600160a01b0316610283565b3480156102c557600080fd5b506102836102d43660046119d8565b610f48565b3480156102e557600080fd5b506101546102f43660046116cb565b610f72565b610301610fe8565b61030a81611042565b1561035c5760405162461bcd60e51b815260206004820152601a60248201527f434e206973207374696c6c20696e2041646472657373426f6f6b00000000000060448201526064015b60405180910390fd5b6001600160a01b038116600090815260ca60205260409020805461037f906119f1565b90506000036103c75760405162461bcd60e51b815260206004820152601460248201527310d3881a5cc81b9bdd081c9959da5cdd195c995960621b6044820152606401610353565b6103d0816110be565b6001600160a01b038116600090815260ca60205260409081902090517fb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d99161041f918491906001820190611aa8565b60405180910390a16001600160a01b038116600090815260ca602052604081209061044a8282611668565b610458600183016000611668565b505050565b60ca60205260009081526040902080548190610478906119f1565b80601f01602080910402602001604051908101604052809291908181526020018280546104a4906119f1565b80156104f15780601f106104c6576101008083540402835291602001916104f1565b820191906000526020600020905b8154815290600101906020018083116104d457829003601f168201915b505050505090806001018054610506906119f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610532906119f1565b801561057f5780601f106105545761010080835404028352916020019161057f565b820191906000526020600020905b81548152906001019060200180831161056257829003601f168201915b5050505050905082565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036105d15760405162461bcd60e51b815260040161035390611ade565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031661061a600080516020611e50833981519152546001600160a01b031690565b6001600160a01b0316146106405760405162461bcd60e51b815260040161035390611b2a565b610649816111c5565b60408051600080825260208201909252610665918391906111cd565b50565b6001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001630036106b05760405162461bcd60e51b815260040161035390611ade565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166106f9600080516020611e50833981519152546001600160a01b031690565b6001600160a01b03161461071f5760405162461bcd60e51b815260040161035390611b2a565b610728826111c5565b610734828260016111cd565b5050565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146107d85760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c00000000000000006064820152608401610353565b50600080516020611e5083398151915290565b60c954606090819067ffffffffffffffff81111561080b5761080b61176d565b604051908082528060200260200182016040528015610834578160200160208202803683370190505b5060c95490925067ffffffffffffffff8111156108535761085361176d565b60405190808252806020026020018201604052801561089857816020015b60408051808201909152606080825260208201528152602001906001900390816108715790505b50905060005b8251811015610aa35760c981815481106108ba576108ba611b76565b9060005260206000200160009054906101000a90046001600160a01b03168382815181106108ea576108ea611b76565b60200260200101906001600160a01b031690816001600160a01b03168152505060ca600060c9838154811061092157610921611b76565b60009182526020808320909101546001600160a01b031683528201929092526040908101909120815180830190925280548290829061095f906119f1565b80601f016020809104026020016040519081016040528092919081815260200182805461098b906119f1565b80156109d85780601f106109ad576101008083540402835291602001916109d8565b820191906000526020600020905b8154815290600101906020018083116109bb57829003601f168201915b505050505081526020016001820180546109f1906119f1565b80601f0160208091040260200160405190810160405280929190818152602001828054610a1d906119f1565b8015610a6a5780601f10610a3f57610100808354040283529160200191610a6a565b820191906000526020600020905b815481529060010190602001808311610a4d57829003601f168201915b505050505081525050828281518110610a8557610a85611b76565b60200260200101819052508080610a9b90611ba2565b91505061089e565b509091565b610ab0610fe8565b610aba6000611338565b565b610ac4610fe8565b838360308114610b165760405162461bcd60e51b815260206004820152601b60248201527f5075626c6963206b6579206d75737420626520343820627974657300000000006044820152606401610353565b6040517fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd29390610b489084908490611bbb565b604051809103902003610b9d5760405162461bcd60e51b815260206004820152601960248201527f5075626c6963206b65792063616e6e6f74206265207a65726f000000000000006044820152606401610353565b838360608114610be65760405162461bcd60e51b8152602060048201526014602482015273506f70206d75737420626520393620627974657360601b6044820152606401610353565b6040517f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c2190610c189084908490611bbb565b604051809103902003610c625760405162461bcd60e51b8152602060048201526012602482015271506f702063616e6e6f74206265207a65726f60701b6044820152606401610353565b610c6b89611042565b610cb75760405162461bcd60e51b815260206004820152601e60248201527f636e4e6f64654964206973206e6f7420696e2041646472657373426f6f6b00006044820152606401610353565b6001600160a01b038916600090815260ca602052604090208054610cda906119f1565b9050600003610d2f5760c980546001810182556000919091527f66be4f155c5ef2ebd3772b228f2f00681e4ed5826cdb3b1943cc11ad15ad1d280180546001600160a01b0319166001600160a01b038b161790555b6040805160606020601f8b018190040282018101835291810189815290918291908b908b9081908501838280828437600092019190915250505090825250604080516020601f8a018190048102820181019092528881529181019190899089908190840183828082843760009201829052509390945250506001600160a01b038c16815260ca6020526040902082519091508190610dcd9082611c19565b5060208201516001820190610de29082611c19565b509050507f79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f8989898989604051610e1d959493929190611d02565b60405180910390a1505050505050505050565b600054610100900460ff1615808015610e505750600054600160ff909116105b80610e6a5750303b158015610e6a575060005460ff166001145b610ecd5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b6064820152608401610353565b6000805460ff191660011790558015610ef0576000805461ff0019166101001790555b610ef861138a565b610f006113b9565b8015610665576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b60c98181548110610f5857600080fd5b6000918252602090912001546001600160a01b0316905081565b610f7a610fe8565b6001600160a01b038116610fdf5760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152608401610353565b61066581611338565b6097546001600160a01b03163314610aba5760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610353565b604051630aabaead60e11b81526001600160a01b0382166004820152600090610400906315575d5a90602401606060405180830381865afa9250505080156110a7575060408051601f3d908101601f191682019092526110a491810190611d46565b60015b6110b357506000919050565b506001949350505050565b60005b60c95481101561073457816001600160a01b031660c982815481106110e8576110e8611b76565b6000918252602090912001546001600160a01b0316036111b35760c9805461111290600190611d93565b8154811061112257611122611b76565b60009182526020909120015460c980546001600160a01b03909216918390811061114e5761114e611b76565b9060005260206000200160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060c980548061118d5761118d611da6565b600082815260209020810160001990810180546001600160a01b03191690550190555050565b806111bd81611ba2565b9150506110c1565b610665610fe8565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff161561120057610458836113e0565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b8152600401602060405180830381865afa92505050801561125a575060408051601f3d908101601f1916820190925261125791810190611dbc565b60015b6112bd5760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b6064820152608401610353565b600080516020611e50833981519152811461132c5760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b6064820152608401610353565b5061045883838361147c565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff166113b15760405162461bcd60e51b815260040161035390611dd5565b610aba6114a7565b600054610100900460ff16610aba5760405162461bcd60e51b815260040161035390611dd5565b6001600160a01b0381163b61144d5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b6064820152608401610353565b600080516020611e5083398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b611485836114d7565b6000825111806114925750805b15610458576114a18383611517565b50505050565b600054610100900460ff166114ce5760405162461bcd60e51b815260040161035390611dd5565b610aba33611338565b6114e0816113e0565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061153c8383604051806060016040528060278152602001611e7060279139611545565b90505b92915050565b6060600080856001600160a01b0316856040516115629190611e20565b600060405180830381855af49150503d806000811461159d576040519150601f19603f3d011682016040523d82523d6000602084013e6115a2565b606091505b50915091506115b3868383876115bd565b9695505050505050565b6060831561162c578251600003611625576001600160a01b0385163b6116255760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610353565b5081611636565b611636838361163e565b949350505050565b81511561164e5781518083602001fd5b8060405162461bcd60e51b81526004016103539190611e3c565b508054611674906119f1565b6000825580601f10611684575050565b601f01602090049060005260206000209081019061066591905b808211156116b2576000815560010161169e565b5090565b6001600160a01b038116811461066557600080fd5b6000602082840312156116dd57600080fd5b81356116e8816116b6565b9392505050565b60005b8381101561170a5781810151838201526020016116f2565b50506000910152565b6000815180845261172b8160208601602086016116ef565b601f01601f19169290920160200192915050565b6040815260006117526040830185611713565b82810360208401526117648185611713565b95945050505050565b634e487b7160e01b600052604160045260246000fd5b6000806040838503121561179657600080fd5b82356117a1816116b6565b9150602083013567ffffffffffffffff808211156117be57600080fd5b818501915085601f8301126117d257600080fd5b8135818111156117e4576117e461176d565b604051601f8201601f19908116603f0116810190838211818310171561180c5761180c61176d565b8160405282815288602084870101111561182557600080fd5b8260208601602083013760006020848301015280955050505050509250929050565b60408082528351828201819052600091906020906060850190828801855b8281101561188a5781516001600160a01b031684529284019290840190600101611865565b50505084810382860152855180825282820190600581901b8301840188850160005b838110156118fc57858303601f19018552815180518985526118d08a860182611713565b91890151858303868b01529190506118e88183611713565b9689019694505050908601906001016118ac565b50909a9950505050505050505050565b60008083601f84011261191e57600080fd5b50813567ffffffffffffffff81111561193657600080fd5b60208301915083602082850101111561194e57600080fd5b9250929050565b60008060008060006060868803121561196d57600080fd5b8535611978816116b6565b9450602086013567ffffffffffffffff8082111561199557600080fd5b6119a189838a0161190c565b909650945060408801359150808211156119ba57600080fd5b506119c78882890161190c565b969995985093965092949392505050565b6000602082840312156119ea57600080fd5b5035919050565b600181811c90821680611a0557607f821691505b602082108103611a2557634e487b7160e01b600052602260045260246000fd5b50919050565b60008154611a38816119f1565b808552602060018381168015611a555760018114611a6f57611a9d565b60ff1985168884015283151560051b880183019550611a9d565b866000528260002060005b85811015611a955781548a8201860152908301908401611a7a565b890184019650505b505050505092915050565b6001600160a01b0384168152606060208201819052600090611acc90830185611a2b565b82810360408401526115b38185611a2b565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b600060018201611bb457611bb4611b8c565b5060010190565b8183823760009101908152919050565b601f82111561045857600081815260208120601f850160051c81016020861015611bf25750805b601f850160051c820191505b81811015611c1157828155600101611bfe565b505050505050565b815167ffffffffffffffff811115611c3357611c3361176d565b611c4781611c4184546119f1565b84611bcb565b602080601f831160018114611c7c5760008415611c645750858301515b600019600386901b1c1916600185901b178555611c11565b600085815260208120601f198616915b82811015611cab57888601518255948401946001909101908401611c8c565b5085821015611cc95787850151600019600388901b60f8161c191681555b5050505050600190811b01905550565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0386168152606060208201819052600090611d279083018688611cd9565b8281036040840152611d3a818587611cd9565b98975050505050505050565b600080600060608486031215611d5b57600080fd5b8351611d66816116b6565b6020850151909350611d77816116b6565b6040850151909250611d88816116b6565b809150509250925092565b8181038181111561153f5761153f611b8c565b634e487b7160e01b600052603160045260246000fd5b600060208284031215611dce57600080fd5b5051919050565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b60008251611e328184602087016116ef565b9190910192915050565b60208152600061153c602083018461171356fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220cf3c282151123924c9c8275c323310bbf7c513b7905cf4ab928cb0d42f59f3a664736f6c63430008130033`

// SimpleBlsRegistryFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use SimpleBlsRegistryMetaData.Sigs instead.
var SimpleBlsRegistryFuncSigs = SimpleBlsRegistryMetaData.Sigs

// SimpleBlsRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SimpleBlsRegistryMetaData.Bin instead.
var SimpleBlsRegistryBin = SimpleBlsRegistryMetaData.Bin

// DeploySimpleBlsRegistry deploys a new Kaia contract, binding an instance of SimpleBlsRegistry to it.
func DeploySimpleBlsRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SimpleBlsRegistry, error) {
	parsed, err := SimpleBlsRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SimpleBlsRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SimpleBlsRegistry{SimpleBlsRegistryCaller: SimpleBlsRegistryCaller{contract: contract}, SimpleBlsRegistryTransactor: SimpleBlsRegistryTransactor{contract: contract}, SimpleBlsRegistryFilterer: SimpleBlsRegistryFilterer{contract: contract}}, nil
}

// SimpleBlsRegistry is an auto generated Go binding around a Kaia contract.
type SimpleBlsRegistry struct {
	SimpleBlsRegistryCaller     // Read-only binding to the contract
	SimpleBlsRegistryTransactor // Write-only binding to the contract
	SimpleBlsRegistryFilterer   // Log filterer for contract events
}

// SimpleBlsRegistryCaller is an auto generated read-only Go binding around a Kaia contract.
type SimpleBlsRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleBlsRegistryTransactor is an auto generated write-only Go binding around a Kaia contract.
type SimpleBlsRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleBlsRegistryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type SimpleBlsRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SimpleBlsRegistrySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type SimpleBlsRegistrySession struct {
	Contract     *SimpleBlsRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// SimpleBlsRegistryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type SimpleBlsRegistryCallerSession struct {
	Contract *SimpleBlsRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// SimpleBlsRegistryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type SimpleBlsRegistryTransactorSession struct {
	Contract     *SimpleBlsRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// SimpleBlsRegistryRaw is an auto generated low-level Go binding around a Kaia contract.
type SimpleBlsRegistryRaw struct {
	Contract *SimpleBlsRegistry // Generic contract binding to access the raw methods on
}

// SimpleBlsRegistryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type SimpleBlsRegistryCallerRaw struct {
	Contract *SimpleBlsRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SimpleBlsRegistryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type SimpleBlsRegistryTransactorRaw struct {
	Contract *SimpleBlsRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSimpleBlsRegistry creates a new instance of SimpleBlsRegistry, bound to a specific deployed contract.
func NewSimpleBlsRegistry(address common.Address, backend bind.ContractBackend) (*SimpleBlsRegistry, error) {
	contract, err := bindSimpleBlsRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistry{SimpleBlsRegistryCaller: SimpleBlsRegistryCaller{contract: contract}, SimpleBlsRegistryTransactor: SimpleBlsRegistryTransactor{contract: contract}, SimpleBlsRegistryFilterer: SimpleBlsRegistryFilterer{contract: contract}}, nil
}

// NewSimpleBlsRegistryCaller creates a new read-only instance of SimpleBlsRegistry, bound to a specific deployed contract.
func NewSimpleBlsRegistryCaller(address common.Address, caller bind.ContractCaller) (*SimpleBlsRegistryCaller, error) {
	contract, err := bindSimpleBlsRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryCaller{contract: contract}, nil
}

// NewSimpleBlsRegistryTransactor creates a new write-only instance of SimpleBlsRegistry, bound to a specific deployed contract.
func NewSimpleBlsRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SimpleBlsRegistryTransactor, error) {
	contract, err := bindSimpleBlsRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryTransactor{contract: contract}, nil
}

// NewSimpleBlsRegistryFilterer creates a new log filterer instance of SimpleBlsRegistry, bound to a specific deployed contract.
func NewSimpleBlsRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SimpleBlsRegistryFilterer, error) {
	contract, err := bindSimpleBlsRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryFilterer{contract: contract}, nil
}

// bindSimpleBlsRegistry binds a generic wrapper to an already deployed contract.
func bindSimpleBlsRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SimpleBlsRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleBlsRegistry *SimpleBlsRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleBlsRegistry.Contract.SimpleBlsRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleBlsRegistry *SimpleBlsRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.SimpleBlsRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleBlsRegistry *SimpleBlsRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.SimpleBlsRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SimpleBlsRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.contract.Transact(opts, method, params...)
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistryCaller) ZERO48HASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SimpleBlsRegistry.contract.Call(opts, &out, "ZERO48HASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) ZERO48HASH() ([32]byte, error) {
	return _SimpleBlsRegistry.Contract.ZERO48HASH(&_SimpleBlsRegistry.CallOpts)
}

// ZERO48HASH is a free data retrieval call binding the contract method 0x6fc522c6.
//
// Solidity: function ZERO48HASH() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerSession) ZERO48HASH() ([32]byte, error) {
	return _SimpleBlsRegistry.Contract.ZERO48HASH(&_SimpleBlsRegistry.CallOpts)
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistryCaller) ZERO96HASH(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SimpleBlsRegistry.contract.Call(opts, &out, "ZERO96HASH")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) ZERO96HASH() ([32]byte, error) {
	return _SimpleBlsRegistry.Contract.ZERO96HASH(&_SimpleBlsRegistry.CallOpts)
}

// ZERO96HASH is a free data retrieval call binding the contract method 0x20abd458.
//
// Solidity: function ZERO96HASH() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerSession) ZERO96HASH() ([32]byte, error) {
	return _SimpleBlsRegistry.Contract.ZERO96HASH(&_SimpleBlsRegistry.CallOpts)
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistryCaller) Abook(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleBlsRegistry.contract.Call(opts, &out, "abook")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) Abook() (common.Address, error) {
	return _SimpleBlsRegistry.Contract.Abook(&_SimpleBlsRegistry.CallOpts)
}

// Abook is a free data retrieval call binding the contract method 0x829d639d.
//
// Solidity: function abook() view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerSession) Abook() (common.Address, error) {
	return _SimpleBlsRegistry.Contract.Abook(&_SimpleBlsRegistry.CallOpts)
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistryCaller) AllNodeIds(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SimpleBlsRegistry.contract.Call(opts, &out, "allNodeIds", arg0)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) AllNodeIds(arg0 *big.Int) (common.Address, error) {
	return _SimpleBlsRegistry.Contract.AllNodeIds(&_SimpleBlsRegistry.CallOpts, arg0)
}

// AllNodeIds is a free data retrieval call binding the contract method 0xa5834971.
//
// Solidity: function allNodeIds(uint256 ) view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerSession) AllNodeIds(arg0 *big.Int) (common.Address, error) {
	return _SimpleBlsRegistry.Contract.AllNodeIds(&_SimpleBlsRegistry.CallOpts, arg0)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_SimpleBlsRegistry *SimpleBlsRegistryCaller) GetAllBlsInfo(opts *bind.CallOpts) (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	var out []interface{}
	err := _SimpleBlsRegistry.contract.Call(opts, &out, "getAllBlsInfo")

	outstruct := new(struct {
		NodeIdList []common.Address
		PubkeyList []IKIP113BlsPublicKeyInfo
	})

	outstruct.NodeIdList = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.PubkeyList = *abi.ConvertType(out[1], new([]IKIP113BlsPublicKeyInfo)).(*[]IKIP113BlsPublicKeyInfo)
	return *outstruct, err
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _SimpleBlsRegistry.Contract.GetAllBlsInfo(&_SimpleBlsRegistry.CallOpts)
}

// GetAllBlsInfo is a free data retrieval call binding the contract method 0x6968b53f.
//
// Solidity: function getAllBlsInfo() view returns(address[] nodeIdList, (bytes,bytes)[] pubkeyList)
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerSession) GetAllBlsInfo() (struct {
	NodeIdList []common.Address
	PubkeyList []IKIP113BlsPublicKeyInfo
}, error,
) {
	return _SimpleBlsRegistry.Contract.GetAllBlsInfo(&_SimpleBlsRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SimpleBlsRegistry.contract.Call(opts, &out, "owner")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) Owner() (common.Address, error) {
	return _SimpleBlsRegistry.Contract.Owner(&_SimpleBlsRegistry.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerSession) Owner() (common.Address, error) {
	return _SimpleBlsRegistry.Contract.Owner(&_SimpleBlsRegistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistryCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _SimpleBlsRegistry.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) ProxiableUUID() ([32]byte, error) {
	return _SimpleBlsRegistry.Contract.ProxiableUUID(&_SimpleBlsRegistry.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerSession) ProxiableUUID() ([32]byte, error) {
	return _SimpleBlsRegistry.Contract.ProxiableUUID(&_SimpleBlsRegistry.CallOpts)
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistryCaller) Record(opts *bind.CallOpts, arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	var out []interface{}
	err := _SimpleBlsRegistry.contract.Call(opts, &out, "record", arg0)

	outstruct := new(struct {
		PublicKey []byte
		Pop       []byte
	})

	outstruct.PublicKey = *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	outstruct.Pop = *abi.ConvertType(out[1], new([]byte)).(*[]byte)
	return *outstruct, err
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) Record(arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	return _SimpleBlsRegistry.Contract.Record(&_SimpleBlsRegistry.CallOpts, arg0)
}

// Record is a free data retrieval call binding the contract method 0x3465d6d5.
//
// Solidity: function record(address ) view returns(bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistryCallerSession) Record(arg0 common.Address) (struct {
	PublicKey []byte
	Pop       []byte
}, error,
) {
	return _SimpleBlsRegistry.Contract.Record(&_SimpleBlsRegistry.CallOpts, arg0)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleBlsRegistry.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) Initialize() (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.Initialize(&_SimpleBlsRegistry.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorSession) Initialize() (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.Initialize(&_SimpleBlsRegistry.TransactOpts)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address cnNodeId, bytes publicKey, bytes pop) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactor) Register(opts *bind.TransactOpts, cnNodeId common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _SimpleBlsRegistry.contract.Transact(opts, "register", cnNodeId, publicKey, pop)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address cnNodeId, bytes publicKey, bytes pop) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) Register(cnNodeId common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.Register(&_SimpleBlsRegistry.TransactOpts, cnNodeId, publicKey, pop)
}

// Register is a paid mutator transaction binding the contract method 0x786cd4d7.
//
// Solidity: function register(address cnNodeId, bytes publicKey, bytes pop) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorSession) Register(cnNodeId common.Address, publicKey []byte, pop []byte) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.Register(&_SimpleBlsRegistry.TransactOpts, cnNodeId, publicKey, pop)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SimpleBlsRegistry.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) RenounceOwnership() (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.RenounceOwnership(&_SimpleBlsRegistry.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.RenounceOwnership(&_SimpleBlsRegistry.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.TransferOwnership(&_SimpleBlsRegistry.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.TransferOwnership(&_SimpleBlsRegistry.TransactOpts, newOwner)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactor) Unregister(opts *bind.TransactOpts, cnNodeId common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.contract.Transact(opts, "unregister", cnNodeId)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) Unregister(cnNodeId common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.Unregister(&_SimpleBlsRegistry.TransactOpts, cnNodeId)
}

// Unregister is a paid mutator transaction binding the contract method 0x2ec2c246.
//
// Solidity: function unregister(address cnNodeId) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorSession) Unregister(cnNodeId common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.Unregister(&_SimpleBlsRegistry.TransactOpts, cnNodeId)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.UpgradeTo(&_SimpleBlsRegistry.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.UpgradeTo(&_SimpleBlsRegistry.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SimpleBlsRegistry.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SimpleBlsRegistry *SimpleBlsRegistrySession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.UpgradeToAndCall(&_SimpleBlsRegistry.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_SimpleBlsRegistry *SimpleBlsRegistryTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _SimpleBlsRegistry.Contract.UpgradeToAndCall(&_SimpleBlsRegistry.TransactOpts, newImplementation, data)
}

// SimpleBlsRegistryAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryAdminChangedIterator struct {
	Event *SimpleBlsRegistryAdminChanged // Event containing the contract specifics and raw log

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
func (it *SimpleBlsRegistryAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleBlsRegistryAdminChanged)
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
		it.Event = new(SimpleBlsRegistryAdminChanged)
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
func (it *SimpleBlsRegistryAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleBlsRegistryAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleBlsRegistryAdminChanged represents a AdminChanged event raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*SimpleBlsRegistryAdminChangedIterator, error) {
	logs, sub, err := _SimpleBlsRegistry.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryAdminChangedIterator{contract: _SimpleBlsRegistry.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *SimpleBlsRegistryAdminChanged) (event.Subscription, error) {
	logs, sub, err := _SimpleBlsRegistry.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleBlsRegistryAdminChanged)
				if err := _SimpleBlsRegistry.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) ParseAdminChanged(log types.Log) (*SimpleBlsRegistryAdminChanged, error) {
	event := new(SimpleBlsRegistryAdminChanged)
	if err := _SimpleBlsRegistry.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SimpleBlsRegistryBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryBeaconUpgradedIterator struct {
	Event *SimpleBlsRegistryBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *SimpleBlsRegistryBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleBlsRegistryBeaconUpgraded)
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
		it.Event = new(SimpleBlsRegistryBeaconUpgraded)
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
func (it *SimpleBlsRegistryBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleBlsRegistryBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleBlsRegistryBeaconUpgraded represents a BeaconUpgraded event raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*SimpleBlsRegistryBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SimpleBlsRegistry.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryBeaconUpgradedIterator{contract: _SimpleBlsRegistry.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *SimpleBlsRegistryBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _SimpleBlsRegistry.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleBlsRegistryBeaconUpgraded)
				if err := _SimpleBlsRegistry.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) ParseBeaconUpgraded(log types.Log) (*SimpleBlsRegistryBeaconUpgraded, error) {
	event := new(SimpleBlsRegistryBeaconUpgraded)
	if err := _SimpleBlsRegistry.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SimpleBlsRegistryInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryInitializedIterator struct {
	Event *SimpleBlsRegistryInitialized // Event containing the contract specifics and raw log

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
func (it *SimpleBlsRegistryInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleBlsRegistryInitialized)
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
		it.Event = new(SimpleBlsRegistryInitialized)
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
func (it *SimpleBlsRegistryInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleBlsRegistryInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleBlsRegistryInitialized represents a Initialized event raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) FilterInitialized(opts *bind.FilterOpts) (*SimpleBlsRegistryInitializedIterator, error) {
	logs, sub, err := _SimpleBlsRegistry.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryInitializedIterator{contract: _SimpleBlsRegistry.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *SimpleBlsRegistryInitialized) (event.Subscription, error) {
	logs, sub, err := _SimpleBlsRegistry.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleBlsRegistryInitialized)
				if err := _SimpleBlsRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) ParseInitialized(log types.Log) (*SimpleBlsRegistryInitialized, error) {
	event := new(SimpleBlsRegistryInitialized)
	if err := _SimpleBlsRegistry.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SimpleBlsRegistryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryOwnershipTransferredIterator struct {
	Event *SimpleBlsRegistryOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SimpleBlsRegistryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleBlsRegistryOwnershipTransferred)
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
		it.Event = new(SimpleBlsRegistryOwnershipTransferred)
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
func (it *SimpleBlsRegistryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleBlsRegistryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleBlsRegistryOwnershipTransferred represents a OwnershipTransferred event raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SimpleBlsRegistryOwnershipTransferredIterator, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimpleBlsRegistry.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryOwnershipTransferredIterator{contract: _SimpleBlsRegistry.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SimpleBlsRegistryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SimpleBlsRegistry.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleBlsRegistryOwnershipTransferred)
				if err := _SimpleBlsRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) ParseOwnershipTransferred(log types.Log) (*SimpleBlsRegistryOwnershipTransferred, error) {
	event := new(SimpleBlsRegistryOwnershipTransferred)
	if err := _SimpleBlsRegistry.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SimpleBlsRegistryRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryRegisteredIterator struct {
	Event *SimpleBlsRegistryRegistered // Event containing the contract specifics and raw log

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
func (it *SimpleBlsRegistryRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleBlsRegistryRegistered)
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
		it.Event = new(SimpleBlsRegistryRegistered)
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
func (it *SimpleBlsRegistryRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleBlsRegistryRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleBlsRegistryRegistered represents a Registered event raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryRegistered struct {
	CnNodeId  common.Address
	PublicKey []byte
	Pop       []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) FilterRegistered(opts *bind.FilterOpts) (*SimpleBlsRegistryRegisteredIterator, error) {
	logs, sub, err := _SimpleBlsRegistry.contract.FilterLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryRegisteredIterator{contract: _SimpleBlsRegistry.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *SimpleBlsRegistryRegistered) (event.Subscription, error) {
	logs, sub, err := _SimpleBlsRegistry.contract.WatchLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleBlsRegistryRegistered)
				if err := _SimpleBlsRegistry.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x79c75399e89a1f580d9a6252cb8bdcf4cd80f73b3597c94d845eb52174ad930f.
//
// Solidity: event Registered(address cnNodeId, bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) ParseRegistered(log types.Log) (*SimpleBlsRegistryRegistered, error) {
	event := new(SimpleBlsRegistryRegistered)
	if err := _SimpleBlsRegistry.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SimpleBlsRegistryUnregisteredIterator is returned from FilterUnregistered and is used to iterate over the raw logs and unpacked data for Unregistered events raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryUnregisteredIterator struct {
	Event *SimpleBlsRegistryUnregistered // Event containing the contract specifics and raw log

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
func (it *SimpleBlsRegistryUnregisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleBlsRegistryUnregistered)
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
		it.Event = new(SimpleBlsRegistryUnregistered)
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
func (it *SimpleBlsRegistryUnregisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleBlsRegistryUnregisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleBlsRegistryUnregistered represents a Unregistered event raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryUnregistered struct {
	CnNodeId  common.Address
	PublicKey []byte
	Pop       []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterUnregistered is a free log retrieval operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) FilterUnregistered(opts *bind.FilterOpts) (*SimpleBlsRegistryUnregisteredIterator, error) {
	logs, sub, err := _SimpleBlsRegistry.contract.FilterLogs(opts, "Unregistered")
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryUnregisteredIterator{contract: _SimpleBlsRegistry.contract, event: "Unregistered", logs: logs, sub: sub}, nil
}

// WatchUnregistered is a free log subscription operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) WatchUnregistered(opts *bind.WatchOpts, sink chan<- *SimpleBlsRegistryUnregistered) (event.Subscription, error) {
	logs, sub, err := _SimpleBlsRegistry.contract.WatchLogs(opts, "Unregistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleBlsRegistryUnregistered)
				if err := _SimpleBlsRegistry.contract.UnpackLog(event, "Unregistered", log); err != nil {
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

// ParseUnregistered is a log parse operation binding the contract event 0xb98b07c4d52e8d65fa5416810f2746a810eb074b1ac7784e1b07e315c0dfd2d9.
//
// Solidity: event Unregistered(address cnNodeId, bytes publicKey, bytes pop)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) ParseUnregistered(log types.Log) (*SimpleBlsRegistryUnregistered, error) {
	event := new(SimpleBlsRegistryUnregistered)
	if err := _SimpleBlsRegistry.contract.UnpackLog(event, "Unregistered", log); err != nil {
		return nil, err
	}
	return event, nil
}

// SimpleBlsRegistryUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryUpgradedIterator struct {
	Event *SimpleBlsRegistryUpgraded // Event containing the contract specifics and raw log

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
func (it *SimpleBlsRegistryUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SimpleBlsRegistryUpgraded)
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
		it.Event = new(SimpleBlsRegistryUpgraded)
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
func (it *SimpleBlsRegistryUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SimpleBlsRegistryUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SimpleBlsRegistryUpgraded represents a Upgraded event raised by the SimpleBlsRegistry contract.
type SimpleBlsRegistryUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*SimpleBlsRegistryUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SimpleBlsRegistry.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &SimpleBlsRegistryUpgradedIterator{contract: _SimpleBlsRegistry.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *SimpleBlsRegistryUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _SimpleBlsRegistry.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SimpleBlsRegistryUpgraded)
				if err := _SimpleBlsRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_SimpleBlsRegistry *SimpleBlsRegistryFilterer) ParseUpgraded(log types.Log) (*SimpleBlsRegistryUpgraded, error) {
	event := new(SimpleBlsRegistryUpgraded)
	if err := _SimpleBlsRegistry.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StorageSlotUpgradeableMetaData contains all meta data concerning the StorageSlotUpgradeable contract.
var StorageSlotUpgradeableMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212206d15c78c51d0895d5280fdb75a0c35fc5f82d2d01a0b996cd84838d0d7b5f77964736f6c63430008130033",
}

// StorageSlotUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageSlotUpgradeableMetaData.ABI instead.
var StorageSlotUpgradeableABI = StorageSlotUpgradeableMetaData.ABI

// StorageSlotUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const StorageSlotUpgradeableBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212206d15c78c51d0895d5280fdb75a0c35fc5f82d2d01a0b996cd84838d0d7b5f77964736f6c63430008130033`

// StorageSlotUpgradeableBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageSlotUpgradeableMetaData.Bin instead.
var StorageSlotUpgradeableBin = StorageSlotUpgradeableMetaData.Bin

// DeployStorageSlotUpgradeable deploys a new Kaia contract, binding an instance of StorageSlotUpgradeable to it.
func DeployStorageSlotUpgradeable(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StorageSlotUpgradeable, error) {
	parsed, err := StorageSlotUpgradeableMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StorageSlotUpgradeableBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageSlotUpgradeable{StorageSlotUpgradeableCaller: StorageSlotUpgradeableCaller{contract: contract}, StorageSlotUpgradeableTransactor: StorageSlotUpgradeableTransactor{contract: contract}, StorageSlotUpgradeableFilterer: StorageSlotUpgradeableFilterer{contract: contract}}, nil
}

// StorageSlotUpgradeable is an auto generated Go binding around a Kaia contract.
type StorageSlotUpgradeable struct {
	StorageSlotUpgradeableCaller     // Read-only binding to the contract
	StorageSlotUpgradeableTransactor // Write-only binding to the contract
	StorageSlotUpgradeableFilterer   // Log filterer for contract events
}

// StorageSlotUpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type StorageSlotUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotUpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type StorageSlotUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotUpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type StorageSlotUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotUpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type StorageSlotUpgradeableSession struct {
	Contract     *StorageSlotUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// StorageSlotUpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type StorageSlotUpgradeableCallerSession struct {
	Contract *StorageSlotUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// StorageSlotUpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type StorageSlotUpgradeableTransactorSession struct {
	Contract     *StorageSlotUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// StorageSlotUpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type StorageSlotUpgradeableRaw struct {
	Contract *StorageSlotUpgradeable // Generic contract binding to access the raw methods on
}

// StorageSlotUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type StorageSlotUpgradeableCallerRaw struct {
	Contract *StorageSlotUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// StorageSlotUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type StorageSlotUpgradeableTransactorRaw struct {
	Contract *StorageSlotUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageSlotUpgradeable creates a new instance of StorageSlotUpgradeable, bound to a specific deployed contract.
func NewStorageSlotUpgradeable(address common.Address, backend bind.ContractBackend) (*StorageSlotUpgradeable, error) {
	contract, err := bindStorageSlotUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageSlotUpgradeable{StorageSlotUpgradeableCaller: StorageSlotUpgradeableCaller{contract: contract}, StorageSlotUpgradeableTransactor: StorageSlotUpgradeableTransactor{contract: contract}, StorageSlotUpgradeableFilterer: StorageSlotUpgradeableFilterer{contract: contract}}, nil
}

// NewStorageSlotUpgradeableCaller creates a new read-only instance of StorageSlotUpgradeable, bound to a specific deployed contract.
func NewStorageSlotUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*StorageSlotUpgradeableCaller, error) {
	contract, err := bindStorageSlotUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSlotUpgradeableCaller{contract: contract}, nil
}

// NewStorageSlotUpgradeableTransactor creates a new write-only instance of StorageSlotUpgradeable, bound to a specific deployed contract.
func NewStorageSlotUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageSlotUpgradeableTransactor, error) {
	contract, err := bindStorageSlotUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSlotUpgradeableTransactor{contract: contract}, nil
}

// NewStorageSlotUpgradeableFilterer creates a new log filterer instance of StorageSlotUpgradeable, bound to a specific deployed contract.
func NewStorageSlotUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageSlotUpgradeableFilterer, error) {
	contract, err := bindStorageSlotUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageSlotUpgradeableFilterer{contract: contract}, nil
}

// bindStorageSlotUpgradeable binds a generic wrapper to an already deployed contract.
func bindStorageSlotUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StorageSlotUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSlotUpgradeable.Contract.StorageSlotUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSlotUpgradeable.Contract.StorageSlotUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSlotUpgradeable.Contract.StorageSlotUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSlotUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSlotUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSlotUpgradeable *StorageSlotUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSlotUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// UUPSUpgradeableMetaData contains all meta data concerning the UUPSUpgradeable contract.
var UUPSUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"}],\"name\":\"upgradeTo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"52d1902d": "proxiableUUID()",
		"3659cfe6": "upgradeTo(address)",
		"4f1ef286": "upgradeToAndCall(address,bytes)",
	},
}

// UUPSUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use UUPSUpgradeableMetaData.ABI instead.
var UUPSUpgradeableABI = UUPSUpgradeableMetaData.ABI

// UUPSUpgradeableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const UUPSUpgradeableBinRuntime = ``

// UUPSUpgradeableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use UUPSUpgradeableMetaData.Sigs instead.
var UUPSUpgradeableFuncSigs = UUPSUpgradeableMetaData.Sigs

// UUPSUpgradeable is an auto generated Go binding around a Kaia contract.
type UUPSUpgradeable struct {
	UUPSUpgradeableCaller     // Read-only binding to the contract
	UUPSUpgradeableTransactor // Write-only binding to the contract
	UUPSUpgradeableFilterer   // Log filterer for contract events
}

// UUPSUpgradeableCaller is an auto generated read-only Go binding around a Kaia contract.
type UUPSUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSUpgradeableTransactor is an auto generated write-only Go binding around a Kaia contract.
type UUPSUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSUpgradeableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type UUPSUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UUPSUpgradeableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type UUPSUpgradeableSession struct {
	Contract     *UUPSUpgradeable  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UUPSUpgradeableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type UUPSUpgradeableCallerSession struct {
	Contract *UUPSUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// UUPSUpgradeableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type UUPSUpgradeableTransactorSession struct {
	Contract     *UUPSUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// UUPSUpgradeableRaw is an auto generated low-level Go binding around a Kaia contract.
type UUPSUpgradeableRaw struct {
	Contract *UUPSUpgradeable // Generic contract binding to access the raw methods on
}

// UUPSUpgradeableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type UUPSUpgradeableCallerRaw struct {
	Contract *UUPSUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// UUPSUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type UUPSUpgradeableTransactorRaw struct {
	Contract *UUPSUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUUPSUpgradeable creates a new instance of UUPSUpgradeable, bound to a specific deployed contract.
func NewUUPSUpgradeable(address common.Address, backend bind.ContractBackend) (*UUPSUpgradeable, error) {
	contract, err := bindUUPSUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeable{UUPSUpgradeableCaller: UUPSUpgradeableCaller{contract: contract}, UUPSUpgradeableTransactor: UUPSUpgradeableTransactor{contract: contract}, UUPSUpgradeableFilterer: UUPSUpgradeableFilterer{contract: contract}}, nil
}

// NewUUPSUpgradeableCaller creates a new read-only instance of UUPSUpgradeable, bound to a specific deployed contract.
func NewUUPSUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*UUPSUpgradeableCaller, error) {
	contract, err := bindUUPSUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableCaller{contract: contract}, nil
}

// NewUUPSUpgradeableTransactor creates a new write-only instance of UUPSUpgradeable, bound to a specific deployed contract.
func NewUUPSUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*UUPSUpgradeableTransactor, error) {
	contract, err := bindUUPSUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableTransactor{contract: contract}, nil
}

// NewUUPSUpgradeableFilterer creates a new log filterer instance of UUPSUpgradeable, bound to a specific deployed contract.
func NewUUPSUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*UUPSUpgradeableFilterer, error) {
	contract, err := bindUUPSUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableFilterer{contract: contract}, nil
}

// bindUUPSUpgradeable binds a generic wrapper to an already deployed contract.
func bindUUPSUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := UUPSUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UUPSUpgradeable *UUPSUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UUPSUpgradeable.Contract.UUPSUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UUPSUpgradeable *UUPSUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UUPSUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UUPSUpgradeable *UUPSUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UUPSUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_UUPSUpgradeable *UUPSUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _UUPSUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_UUPSUpgradeable *UUPSUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_UUPSUpgradeable *UUPSUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSUpgradeable *UUPSUpgradeableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _UUPSUpgradeable.contract.Call(opts, &out, "proxiableUUID")
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSUpgradeable *UUPSUpgradeableSession) ProxiableUUID() ([32]byte, error) {
	return _UUPSUpgradeable.Contract.ProxiableUUID(&_UUPSUpgradeable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_UUPSUpgradeable *UUPSUpgradeableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _UUPSUpgradeable.Contract.ProxiableUUID(&_UUPSUpgradeable.CallOpts)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSUpgradeable *UUPSUpgradeableTransactor) UpgradeTo(opts *bind.TransactOpts, newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSUpgradeable.contract.Transact(opts, "upgradeTo", newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSUpgradeable *UUPSUpgradeableSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UpgradeTo(&_UUPSUpgradeable.TransactOpts, newImplementation)
}

// UpgradeTo is a paid mutator transaction binding the contract method 0x3659cfe6.
//
// Solidity: function upgradeTo(address newImplementation) returns()
func (_UUPSUpgradeable *UUPSUpgradeableTransactorSession) UpgradeTo(newImplementation common.Address) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UpgradeTo(&_UUPSUpgradeable.TransactOpts, newImplementation)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSUpgradeable *UUPSUpgradeableTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSUpgradeable.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSUpgradeable *UUPSUpgradeableSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UpgradeToAndCall(&_UUPSUpgradeable.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_UUPSUpgradeable *UUPSUpgradeableTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _UUPSUpgradeable.Contract.UpgradeToAndCall(&_UUPSUpgradeable.TransactOpts, newImplementation, data)
}

// UUPSUpgradeableAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the UUPSUpgradeable contract.
type UUPSUpgradeableAdminChangedIterator struct {
	Event *UUPSUpgradeableAdminChanged // Event containing the contract specifics and raw log

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
func (it *UUPSUpgradeableAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSUpgradeableAdminChanged)
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
		it.Event = new(UUPSUpgradeableAdminChanged)
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
func (it *UUPSUpgradeableAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSUpgradeableAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSUpgradeableAdminChanged represents a AdminChanged event raised by the UUPSUpgradeable contract.
type UUPSUpgradeableAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*UUPSUpgradeableAdminChangedIterator, error) {
	logs, sub, err := _UUPSUpgradeable.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableAdminChangedIterator{contract: _UUPSUpgradeable.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *UUPSUpgradeableAdminChanged) (event.Subscription, error) {
	logs, sub, err := _UUPSUpgradeable.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSUpgradeableAdminChanged)
				if err := _UUPSUpgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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

// ParseAdminChanged is a log parse operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) ParseAdminChanged(log types.Log) (*UUPSUpgradeableAdminChanged, error) {
	event := new(UUPSUpgradeableAdminChanged)
	if err := _UUPSUpgradeable.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// UUPSUpgradeableBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the UUPSUpgradeable contract.
type UUPSUpgradeableBeaconUpgradedIterator struct {
	Event *UUPSUpgradeableBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *UUPSUpgradeableBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSUpgradeableBeaconUpgraded)
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
		it.Event = new(UUPSUpgradeableBeaconUpgraded)
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
func (it *UUPSUpgradeableBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSUpgradeableBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSUpgradeableBeaconUpgraded represents a BeaconUpgraded event raised by the UUPSUpgradeable contract.
type UUPSUpgradeableBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*UUPSUpgradeableBeaconUpgradedIterator, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _UUPSUpgradeable.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableBeaconUpgradedIterator{contract: _UUPSUpgradeable.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *UUPSUpgradeableBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {
	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _UUPSUpgradeable.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSUpgradeableBeaconUpgraded)
				if err := _UUPSUpgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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

// ParseBeaconUpgraded is a log parse operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) ParseBeaconUpgraded(log types.Log) (*UUPSUpgradeableBeaconUpgraded, error) {
	event := new(UUPSUpgradeableBeaconUpgraded)
	if err := _UUPSUpgradeable.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// UUPSUpgradeableInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the UUPSUpgradeable contract.
type UUPSUpgradeableInitializedIterator struct {
	Event *UUPSUpgradeableInitialized // Event containing the contract specifics and raw log

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
func (it *UUPSUpgradeableInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSUpgradeableInitialized)
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
		it.Event = new(UUPSUpgradeableInitialized)
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
func (it *UUPSUpgradeableInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSUpgradeableInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSUpgradeableInitialized represents a Initialized event raised by the UUPSUpgradeable contract.
type UUPSUpgradeableInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) FilterInitialized(opts *bind.FilterOpts) (*UUPSUpgradeableInitializedIterator, error) {
	logs, sub, err := _UUPSUpgradeable.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableInitializedIterator{contract: _UUPSUpgradeable.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *UUPSUpgradeableInitialized) (event.Subscription, error) {
	logs, sub, err := _UUPSUpgradeable.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSUpgradeableInitialized)
				if err := _UUPSUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) ParseInitialized(log types.Log) (*UUPSUpgradeableInitialized, error) {
	event := new(UUPSUpgradeableInitialized)
	if err := _UUPSUpgradeable.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	return event, nil
}

// UUPSUpgradeableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the UUPSUpgradeable contract.
type UUPSUpgradeableUpgradedIterator struct {
	Event *UUPSUpgradeableUpgraded // Event containing the contract specifics and raw log

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
func (it *UUPSUpgradeableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UUPSUpgradeableUpgraded)
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
		it.Event = new(UUPSUpgradeableUpgraded)
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
func (it *UUPSUpgradeableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UUPSUpgradeableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UUPSUpgradeableUpgraded represents a Upgraded event raised by the UUPSUpgradeable contract.
type UUPSUpgradeableUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*UUPSUpgradeableUpgradedIterator, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UUPSUpgradeable.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &UUPSUpgradeableUpgradedIterator{contract: _UUPSUpgradeable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *UUPSUpgradeableUpgraded, implementation []common.Address) (event.Subscription, error) {
	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _UUPSUpgradeable.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UUPSUpgradeableUpgraded)
				if err := _UUPSUpgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_UUPSUpgradeable *UUPSUpgradeableFilterer) ParseUpgraded(log types.Log) (*UUPSUpgradeableUpgraded, error) {
	event := new(UUPSUpgradeableUpgraded)
	if err := _UUPSUpgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}
