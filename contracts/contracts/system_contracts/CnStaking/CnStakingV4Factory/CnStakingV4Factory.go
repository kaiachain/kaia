// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cnstakingv4factory

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

// BeaconProxyMetaData contains all meta data concerning the BeaconProxy contract.
var BeaconProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidBeacon\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
	Bin: "0x60a060405260405161053f38038061053f83398101604081905261002291610331565b61002c828261003e565b506001600160a01b031660805261040d565b610047826100fb565b6040516001600160a01b038316907f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e905f90a28051156100ef576100ea826001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156100c0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906100e491906103ed565b82610209565b505050565b6100f76102aa565b5050565b806001600160a01b03163b5f0361013557604051631933b43b60e21b81526001600160a01b03821660048201526024015b60405180910390fd5b807fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b0392831617905560408051635c60da1b60e01b815290515f92841691635c60da1b9160048083019260209291908290030181865afa1580156101ae573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101d291906103ed565b9050806001600160a01b03163b5f036100f757604051634c9c8ce360e01b81526001600160a01b038216600482015260240161012c565b60605f61021684846102cb565b905080801561023757505f3d118061023757505f846001600160a01b03163b115b1561024c576102446102de565b9150506102a4565b801561027657604051639996b31560e01b81526001600160a01b038516600482015260240161012c565b3d15610289576102846102f7565b6102a2565b60405163d6bda27560e01b815260040160405180910390fd5b505b92915050565b34156102c95760405163b398979f60e01b815260040160405180910390fd5b565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b80516001600160a01b0381168114610318575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f8060408385031215610342575f80fd5b61034b83610302565b60208401519092506001600160401b0380821115610367575f80fd5b818501915085601f83011261037a575f80fd5b81518181111561038c5761038c61031d565b604051601f8201601f19908116603f011681019083821181831017156103b4576103b461031d565b816040528281528860208487010111156103cc575f80fd5b8260208601602083015e5f6020848301015280955050505050509250929050565b5f602082840312156103fd575f80fd5b61040682610302565b9392505050565b60805161011b6104245f395f601d015261011b5ff3fe6080604052600a600c565b005b60186014601a565b609d565b565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156076573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906098919060ba565b905090565b365f80375f80365f845af43d5f803e80801560b6573d5ff35b3d5ffd5b5f6020828403121560c9575f80fd5b81516001600160a01b038116811460de575f80fd5b939250505056fea2646970667358221220b4d171c91104b11935acfb0e21f3489de043d1f81010125fab26bfe1470a144e64736f6c63430008190033",
}

// BeaconProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use BeaconProxyMetaData.ABI instead.
var BeaconProxyABI = BeaconProxyMetaData.ABI

// BeaconProxyBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const BeaconProxyBinRuntime = `6080604052600a600c565b005b60186014601a565b609d565b565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156076573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906098919060ba565b905090565b365f80375f80365f845af43d5f803e80801560b6573d5ff35b3d5ffd5b5f6020828403121560c9575f80fd5b81516001600160a01b038116811460de575f80fd5b939250505056fea2646970667358221220b4d171c91104b11935acfb0e21f3489de043d1f81010125fab26bfe1470a144e64736f6c63430008190033`

// BeaconProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BeaconProxyMetaData.Bin instead.
var BeaconProxyBin = BeaconProxyMetaData.Bin

// DeployBeaconProxy deploys a new Kaia contract, binding an instance of BeaconProxy to it.
func DeployBeaconProxy(auth *bind.TransactOpts, backend bind.ContractBackend, beacon common.Address, data []byte) (common.Address, *types.Transaction, *BeaconProxy, error) {
	parsed, err := BeaconProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BeaconProxyBin), backend, beacon, data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BeaconProxy{BeaconProxyCaller: BeaconProxyCaller{contract: contract}, BeaconProxyTransactor: BeaconProxyTransactor{contract: contract}, BeaconProxyFilterer: BeaconProxyFilterer{contract: contract}}, nil
}

// BeaconProxy is an auto generated Go binding around a Kaia contract.
type BeaconProxy struct {
	BeaconProxyCaller     // Read-only binding to the contract
	BeaconProxyTransactor // Write-only binding to the contract
	BeaconProxyFilterer   // Log filterer for contract events
}

// BeaconProxyCaller is an auto generated read-only Go binding around a Kaia contract.
type BeaconProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconProxyTransactor is an auto generated write-only Go binding around a Kaia contract.
type BeaconProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconProxyFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type BeaconProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BeaconProxySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type BeaconProxySession struct {
	Contract     *BeaconProxy      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BeaconProxyCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type BeaconProxyCallerSession struct {
	Contract *BeaconProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BeaconProxyTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type BeaconProxyTransactorSession struct {
	Contract     *BeaconProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BeaconProxyRaw is an auto generated low-level Go binding around a Kaia contract.
type BeaconProxyRaw struct {
	Contract *BeaconProxy // Generic contract binding to access the raw methods on
}

// BeaconProxyCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type BeaconProxyCallerRaw struct {
	Contract *BeaconProxyCaller // Generic read-only contract binding to access the raw methods on
}

// BeaconProxyTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type BeaconProxyTransactorRaw struct {
	Contract *BeaconProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBeaconProxy creates a new instance of BeaconProxy, bound to a specific deployed contract.
func NewBeaconProxy(address common.Address, backend bind.ContractBackend) (*BeaconProxy, error) {
	contract, err := bindBeaconProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BeaconProxy{BeaconProxyCaller: BeaconProxyCaller{contract: contract}, BeaconProxyTransactor: BeaconProxyTransactor{contract: contract}, BeaconProxyFilterer: BeaconProxyFilterer{contract: contract}}, nil
}

// NewBeaconProxyCaller creates a new read-only instance of BeaconProxy, bound to a specific deployed contract.
func NewBeaconProxyCaller(address common.Address, caller bind.ContractCaller) (*BeaconProxyCaller, error) {
	contract, err := bindBeaconProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconProxyCaller{contract: contract}, nil
}

// NewBeaconProxyTransactor creates a new write-only instance of BeaconProxy, bound to a specific deployed contract.
func NewBeaconProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*BeaconProxyTransactor, error) {
	contract, err := bindBeaconProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BeaconProxyTransactor{contract: contract}, nil
}

// NewBeaconProxyFilterer creates a new log filterer instance of BeaconProxy, bound to a specific deployed contract.
func NewBeaconProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*BeaconProxyFilterer, error) {
	contract, err := bindBeaconProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BeaconProxyFilterer{contract: contract}, nil
}

// bindBeaconProxy binds a generic wrapper to an already deployed contract.
func bindBeaconProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BeaconProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconProxy *BeaconProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconProxy.Contract.BeaconProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconProxy *BeaconProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconProxy.Contract.BeaconProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconProxy *BeaconProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconProxy.Contract.BeaconProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BeaconProxy *BeaconProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BeaconProxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BeaconProxy *BeaconProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BeaconProxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BeaconProxy *BeaconProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BeaconProxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BeaconProxy *BeaconProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _BeaconProxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BeaconProxy *BeaconProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BeaconProxy.Contract.Fallback(&_BeaconProxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_BeaconProxy *BeaconProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _BeaconProxy.Contract.Fallback(&_BeaconProxy.TransactOpts, calldata)
}

// BeaconProxyBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the BeaconProxy contract.
type BeaconProxyBeaconUpgradedIterator struct {
	Event *BeaconProxyBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *BeaconProxyBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BeaconProxyBeaconUpgraded)
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
		it.Event = new(BeaconProxyBeaconUpgraded)
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
func (it *BeaconProxyBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BeaconProxyBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BeaconProxyBeaconUpgraded represents a BeaconUpgraded event raised by the BeaconProxy contract.
type BeaconProxyBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_BeaconProxy *BeaconProxyFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*BeaconProxyBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _BeaconProxy.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &BeaconProxyBeaconUpgradedIterator{contract: _BeaconProxy.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_BeaconProxy *BeaconProxyFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *BeaconProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _BeaconProxy.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BeaconProxyBeaconUpgraded)
				if err := _BeaconProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_BeaconProxy *BeaconProxyFilterer) ParseBeaconUpgraded(log types.Log) (*BeaconProxyBeaconUpgraded, error) {
	event := new(BeaconProxyBeaconUpgraded)
	if err := _BeaconProxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CnStakingV4FactoryMetaData contains all meta data concerning the CnStakingV4Factory contract.
var CnStakingV4FactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_cnStakingBeacon\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_pdBeacon\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"InsufficientInitialStake\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"DeployCnStakingV4\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"publicDelegation\",\"type\":\"address\"}],\"name\":\"DeployCnStakingV4WithPD\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEAD_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIAL_LOCKUP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cnStakingBeacon\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"deployCnStaking\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gcName\",\"type\":\"string\"}],\"internalType\":\"structIPublicDelegation.PDConstructorArgs\",\"name\":\"_pdArgs\",\"type\":\"tuple\"}],\"name\":\"deployCnStakingWithPD\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicDelegation\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getDeployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isDeployedCnStaking\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isDeployedPublicDelegation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pdBeacon\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"4e6fd6c4": "DEAD_ADDRESS()",
		"0cf74c5c": "INITIAL_LOCKUP()",
		"7b0e7fdd": "cnStakingBeacon()",
		"4ed7a764": "deployCnStaking(address)",
		"33f5db27": "deployCnStakingWithPD(address,(address,address,uint256,string))",
		"669d8d45": "getDeployer(address)",
		"9a429925": "isDeployedCnStaking(address)",
		"7dfe297c": "isDeployedPublicDelegation(address)",
		"aa777f4c": "pdBeacon()",
	},
	Bin: "0x60c060405234801561000f575f80fd5b50604051610fc5380380610fc583398101604081905261002e916100ae565b6001600160a01b0382166100555760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b03811661007c5760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b039182166080521660a0526100df565b80516001600160a01b03811681146100a9575f80fd5b919050565b5f80604083850312156100bf575f80fd5b6100c883610093565b91506100d660208401610093565b90509250929050565b60805160a051610eb06101155f395f818161022a015261034101525f8181610179015281816102d501526105a00152610eb05ff3fe608060405260043610610084575f3560e01c8063669d8d4511610057578063669d8d45146101315780637b0e7fdd146101685780637dfe297c1461019b5780639a429925146101e2578063aa777f4c14610219575f80fd5b80630cf74c5c1461008857806333f5db27146100b25780634e6fd6c4146100e55780634ed7a76414610112575b5f80fd5b348015610093575f80fd5b5061009f633b9aca0081565b6040519081526020015b60405180910390f35b6100c56100c0366004610762565b61024c565b604080516001600160a01b039384168152929091166020830152016100a9565b3480156100f0575f80fd5b506100fa61dead81565b6040516001600160a01b0390911681526020016100a9565b34801561011d575f80fd5b506100fa61012c36600461085a565b610567565b34801561013c575f80fd5b506100fa61014b36600461085a565b6001600160a01b039081165f908152600260205260409020541690565b348015610173575f80fd5b506100fa7f000000000000000000000000000000000000000000000000000000000000000081565b3480156101a6575f80fd5b506101d26101b536600461085a565b6001600160a01b03165f9081526001602052604090205460ff1690565b60405190151581526020016100a9565b3480156101ed575f80fd5b506101d26101fc36600461085a565b6001600160a01b03165f9081526020819052604090205460ff1690565b348015610224575f80fd5b506100fa7f000000000000000000000000000000000000000000000000000000000000000081565b5f80633b9aca003410156102735760405163176c085f60e31b815260040160405180910390fd5b604080513360208083018290526001600160a01b03881683850152835180840385018152606084019094528351930192909220915f916102ba9190889088906080016108dc565b604051602081830303815290604052805190602001209050817f0000000000000000000000000000000000000000000000000000000000000000604051610300906106cc565b6001600160a01b0390911681526040602082018190525f908201526060018190604051809103905ff590508015801561033b573d5f803e3d5ffd5b509350807f000000000000000000000000000000000000000000000000000000000000000060405161036c906106cc565b6001600160a01b0390911681526040602082018190525f908201526060018190604051809103905ff59050801580156103a7573d5f803e3d5ffd5b5060405163136793bd60e11b81529093506001600160a01b038416906326cf277a906103d99087908990600401610910565b5f604051808303815f87803b1580156103f0575f80fd5b505af1158015610402573d5f803e3d5ffd5b5050604051633e8f533560e01b81526001600160a01b038981166004830152868116602483015287169250633e8f533591506044015f604051808303815f87803b15801561044e575f80fd5b505af1158015610460573d5f803e3d5ffd5b50506040516325fb490360e11b815261dead60048201526001600160a01b0386169250634bf69206915034906024015f604051808303818588803b1580156104a6575f80fd5b505af11580156104b8573d5f803e3d5ffd5b505050506001600160a01b038581165f818152602081815260408083208054600160ff1991821681179092558a871680865282855283862080549092169092179055848452600283528184208054336001600160a01b031991821681179092558286529483902080549095161790935551918252928a16935090917f87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f910160405180910390a350509250929050565b60408051336020808301919091526001600160a01b0384168284015282518083038401815260609092019283905281519101205f9181907f0000000000000000000000000000000000000000000000000000000000000000906105c9906106cc565b6001600160a01b0390911681526040602082018190525f908201526060018190604051809103905ff5905080158015610604573d5f803e3d5ffd5b5060405163189acdbd60e31b81526001600160a01b0385811660048301529193509083169063c4d66de8906024015f604051808303815f87803b158015610649575f80fd5b505af115801561065b573d5f803e3d5ffd5b5050506001600160a01b038084165f81815260208181526040808320805460ff19166001179055600290915280822080546001600160a01b0319163317905551928716935090917f1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef9190a350919050565b61053f8061093c83390190565b80356001600160a01b03811681146106ef575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b6040516080810167ffffffffffffffff8111828210171561072b5761072b6106f4565b60405290565b604051601f8201601f1916810167ffffffffffffffff8111828210171561075a5761075a6106f4565b604052919050565b5f8060408385031215610773575f80fd5b61077c836106d9565b915060208084013567ffffffffffffffff80821115610799575f80fd5b90850190608082880312156107ac575f80fd5b6107b4610708565b6107bd836106d9565b81526107ca8484016106d9565b84820152604083013560408201526060830135828111156107e9575f80fd5b80840193505087601f8401126107fd575f80fd5b82358281111561080f5761080f6106f4565b610821601f8201601f19168601610731565b92508083528885828601011115610836575f80fd5b80858501868501375f85828501015250816060820152809450505050509250929050565b5f6020828403121561086a575f80fd5b610873826106d9565b9392505050565b5f60018060a01b0380835116845280602084015116602085015250604082015160408401526060820151608060608501528051806080860152806020830160a087015e5f60a0828701015260a0601f19601f8301168601019250505092915050565b6001600160a01b038481168252831660208201526060604082018190525f906109079083018461087a565b95945050505050565b6001600160a01b03831681526040602082018190525f906109339083018461087a565b94935050505056fe60a060405260405161053f38038061053f83398101604081905261002291610331565b61002c828261003e565b506001600160a01b031660805261040d565b610047826100fb565b6040516001600160a01b038316907f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e905f90a28051156100ef576100ea826001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156100c0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906100e491906103ed565b82610209565b505050565b6100f76102aa565b5050565b806001600160a01b03163b5f0361013557604051631933b43b60e21b81526001600160a01b03821660048201526024015b60405180910390fd5b807fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b0392831617905560408051635c60da1b60e01b815290515f92841691635c60da1b9160048083019260209291908290030181865afa1580156101ae573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101d291906103ed565b9050806001600160a01b03163b5f036100f757604051634c9c8ce360e01b81526001600160a01b038216600482015260240161012c565b60605f61021684846102cb565b905080801561023757505f3d118061023757505f846001600160a01b03163b115b1561024c576102446102de565b9150506102a4565b801561027657604051639996b31560e01b81526001600160a01b038516600482015260240161012c565b3d15610289576102846102f7565b6102a2565b60405163d6bda27560e01b815260040160405180910390fd5b505b92915050565b34156102c95760405163b398979f60e01b815260040160405180910390fd5b565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b80516001600160a01b0381168114610318575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f8060408385031215610342575f80fd5b61034b83610302565b60208401519092506001600160401b0380821115610367575f80fd5b818501915085601f83011261037a575f80fd5b81518181111561038c5761038c61031d565b604051601f8201601f19908116603f011681019083821181831017156103b4576103b461031d565b816040528281528860208487010111156103cc575f80fd5b8260208601602083015e5f6020848301015280955050505050509250929050565b5f602082840312156103fd575f80fd5b61040682610302565b9392505050565b60805161011b6104245f395f601d015261011b5ff3fe6080604052600a600c565b005b60186014601a565b609d565b565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156076573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906098919060ba565b905090565b365f80375f80365f845af43d5f803e80801560b6573d5ff35b3d5ffd5b5f6020828403121560c9575f80fd5b81516001600160a01b038116811460de575f80fd5b939250505056fea2646970667358221220b4d171c91104b11935acfb0e21f3489de043d1f81010125fab26bfe1470a144e64736f6c63430008190033a26469706673582212205c1ad77ef1ff82a1619bd43dc20c936d5e5c1fe63ca204f53981f15c6e496d7164736f6c63430008190033",
}

// CnStakingV4FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use CnStakingV4FactoryMetaData.ABI instead.
var CnStakingV4FactoryABI = CnStakingV4FactoryMetaData.ABI

// CnStakingV4FactoryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const CnStakingV4FactoryBinRuntime = `608060405260043610610084575f3560e01c8063669d8d4511610057578063669d8d45146101315780637b0e7fdd146101685780637dfe297c1461019b5780639a429925146101e2578063aa777f4c14610219575f80fd5b80630cf74c5c1461008857806333f5db27146100b25780634e6fd6c4146100e55780634ed7a76414610112575b5f80fd5b348015610093575f80fd5b5061009f633b9aca0081565b6040519081526020015b60405180910390f35b6100c56100c0366004610762565b61024c565b604080516001600160a01b039384168152929091166020830152016100a9565b3480156100f0575f80fd5b506100fa61dead81565b6040516001600160a01b0390911681526020016100a9565b34801561011d575f80fd5b506100fa61012c36600461085a565b610567565b34801561013c575f80fd5b506100fa61014b36600461085a565b6001600160a01b039081165f908152600260205260409020541690565b348015610173575f80fd5b506100fa7f000000000000000000000000000000000000000000000000000000000000000081565b3480156101a6575f80fd5b506101d26101b536600461085a565b6001600160a01b03165f9081526001602052604090205460ff1690565b60405190151581526020016100a9565b3480156101ed575f80fd5b506101d26101fc36600461085a565b6001600160a01b03165f9081526020819052604090205460ff1690565b348015610224575f80fd5b506100fa7f000000000000000000000000000000000000000000000000000000000000000081565b5f80633b9aca003410156102735760405163176c085f60e31b815260040160405180910390fd5b604080513360208083018290526001600160a01b03881683850152835180840385018152606084019094528351930192909220915f916102ba9190889088906080016108dc565b604051602081830303815290604052805190602001209050817f0000000000000000000000000000000000000000000000000000000000000000604051610300906106cc565b6001600160a01b0390911681526040602082018190525f908201526060018190604051809103905ff590508015801561033b573d5f803e3d5ffd5b509350807f000000000000000000000000000000000000000000000000000000000000000060405161036c906106cc565b6001600160a01b0390911681526040602082018190525f908201526060018190604051809103905ff59050801580156103a7573d5f803e3d5ffd5b5060405163136793bd60e11b81529093506001600160a01b038416906326cf277a906103d99087908990600401610910565b5f604051808303815f87803b1580156103f0575f80fd5b505af1158015610402573d5f803e3d5ffd5b5050604051633e8f533560e01b81526001600160a01b038981166004830152868116602483015287169250633e8f533591506044015f604051808303815f87803b15801561044e575f80fd5b505af1158015610460573d5f803e3d5ffd5b50506040516325fb490360e11b815261dead60048201526001600160a01b0386169250634bf69206915034906024015f604051808303818588803b1580156104a6575f80fd5b505af11580156104b8573d5f803e3d5ffd5b505050506001600160a01b038581165f818152602081815260408083208054600160ff1991821681179092558a871680865282855283862080549092169092179055848452600283528184208054336001600160a01b031991821681179092558286529483902080549095161790935551918252928a16935090917f87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f910160405180910390a350509250929050565b60408051336020808301919091526001600160a01b0384168284015282518083038401815260609092019283905281519101205f9181907f0000000000000000000000000000000000000000000000000000000000000000906105c9906106cc565b6001600160a01b0390911681526040602082018190525f908201526060018190604051809103905ff5905080158015610604573d5f803e3d5ffd5b5060405163189acdbd60e31b81526001600160a01b0385811660048301529193509083169063c4d66de8906024015f604051808303815f87803b158015610649575f80fd5b505af115801561065b573d5f803e3d5ffd5b5050506001600160a01b038084165f81815260208181526040808320805460ff19166001179055600290915280822080546001600160a01b0319163317905551928716935090917f1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef9190a350919050565b61053f8061093c83390190565b80356001600160a01b03811681146106ef575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b6040516080810167ffffffffffffffff8111828210171561072b5761072b6106f4565b60405290565b604051601f8201601f1916810167ffffffffffffffff8111828210171561075a5761075a6106f4565b604052919050565b5f8060408385031215610773575f80fd5b61077c836106d9565b915060208084013567ffffffffffffffff80821115610799575f80fd5b90850190608082880312156107ac575f80fd5b6107b4610708565b6107bd836106d9565b81526107ca8484016106d9565b84820152604083013560408201526060830135828111156107e9575f80fd5b80840193505087601f8401126107fd575f80fd5b82358281111561080f5761080f6106f4565b610821601f8201601f19168601610731565b92508083528885828601011115610836575f80fd5b80858501868501375f85828501015250816060820152809450505050509250929050565b5f6020828403121561086a575f80fd5b610873826106d9565b9392505050565b5f60018060a01b0380835116845280602084015116602085015250604082015160408401526060820151608060608501528051806080860152806020830160a087015e5f60a0828701015260a0601f19601f8301168601019250505092915050565b6001600160a01b038481168252831660208201526060604082018190525f906109079083018461087a565b95945050505050565b6001600160a01b03831681526040602082018190525f906109339083018461087a565b94935050505056fe60a060405260405161053f38038061053f83398101604081905261002291610331565b61002c828261003e565b506001600160a01b031660805261040d565b610047826100fb565b6040516001600160a01b038316907f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e905f90a28051156100ef576100ea826001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156100c0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906100e491906103ed565b82610209565b505050565b6100f76102aa565b5050565b806001600160a01b03163b5f0361013557604051631933b43b60e21b81526001600160a01b03821660048201526024015b60405180910390fd5b807fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b0392831617905560408051635c60da1b60e01b815290515f92841691635c60da1b9160048083019260209291908290030181865afa1580156101ae573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101d291906103ed565b9050806001600160a01b03163b5f036100f757604051634c9c8ce360e01b81526001600160a01b038216600482015260240161012c565b60605f61021684846102cb565b905080801561023757505f3d118061023757505f846001600160a01b03163b115b1561024c576102446102de565b9150506102a4565b801561027657604051639996b31560e01b81526001600160a01b038516600482015260240161012c565b3d15610289576102846102f7565b6102a2565b60405163d6bda27560e01b815260040160405180910390fd5b505b92915050565b34156102c95760405163b398979f60e01b815260040160405180910390fd5b565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b80516001600160a01b0381168114610318575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f8060408385031215610342575f80fd5b61034b83610302565b60208401519092506001600160401b0380821115610367575f80fd5b818501915085601f83011261037a575f80fd5b81518181111561038c5761038c61031d565b604051601f8201601f19908116603f011681019083821181831017156103b4576103b461031d565b816040528281528860208487010111156103cc575f80fd5b8260208601602083015e5f6020848301015280955050505050509250929050565b5f602082840312156103fd575f80fd5b61040682610302565b9392505050565b60805161011b6104245f395f601d015261011b5ff3fe6080604052600a600c565b005b60186014601a565b609d565b565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156076573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906098919060ba565b905090565b365f80375f80365f845af43d5f803e80801560b6573d5ff35b3d5ffd5b5f6020828403121560c9575f80fd5b81516001600160a01b038116811460de575f80fd5b939250505056fea2646970667358221220b4d171c91104b11935acfb0e21f3489de043d1f81010125fab26bfe1470a144e64736f6c63430008190033a26469706673582212205c1ad77ef1ff82a1619bd43dc20c936d5e5c1fe63ca204f53981f15c6e496d7164736f6c63430008190033`

// Deprecated: Use CnStakingV4FactoryMetaData.Sigs instead.
// CnStakingV4FactoryFuncSigs maps the 4-byte function signature to its string representation.
var CnStakingV4FactoryFuncSigs = CnStakingV4FactoryMetaData.Sigs

// CnStakingV4FactoryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CnStakingV4FactoryMetaData.Bin instead.
var CnStakingV4FactoryBin = CnStakingV4FactoryMetaData.Bin

// DeployCnStakingV4Factory deploys a new Kaia contract, binding an instance of CnStakingV4Factory to it.
func DeployCnStakingV4Factory(auth *bind.TransactOpts, backend bind.ContractBackend, _cnStakingBeacon common.Address, _pdBeacon common.Address) (common.Address, *types.Transaction, *CnStakingV4Factory, error) {
	parsed, err := CnStakingV4FactoryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CnStakingV4FactoryBin), backend, _cnStakingBeacon, _pdBeacon)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CnStakingV4Factory{CnStakingV4FactoryCaller: CnStakingV4FactoryCaller{contract: contract}, CnStakingV4FactoryTransactor: CnStakingV4FactoryTransactor{contract: contract}, CnStakingV4FactoryFilterer: CnStakingV4FactoryFilterer{contract: contract}}, nil
}

// CnStakingV4Factory is an auto generated Go binding around a Kaia contract.
type CnStakingV4Factory struct {
	CnStakingV4FactoryCaller     // Read-only binding to the contract
	CnStakingV4FactoryTransactor // Write-only binding to the contract
	CnStakingV4FactoryFilterer   // Log filterer for contract events
}

// CnStakingV4FactoryCaller is an auto generated read-only Go binding around a Kaia contract.
type CnStakingV4FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CnStakingV4FactoryTransactor is an auto generated write-only Go binding around a Kaia contract.
type CnStakingV4FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CnStakingV4FactoryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type CnStakingV4FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CnStakingV4FactorySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type CnStakingV4FactorySession struct {
	Contract     *CnStakingV4Factory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// CnStakingV4FactoryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type CnStakingV4FactoryCallerSession struct {
	Contract *CnStakingV4FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// CnStakingV4FactoryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type CnStakingV4FactoryTransactorSession struct {
	Contract     *CnStakingV4FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// CnStakingV4FactoryRaw is an auto generated low-level Go binding around a Kaia contract.
type CnStakingV4FactoryRaw struct {
	Contract *CnStakingV4Factory // Generic contract binding to access the raw methods on
}

// CnStakingV4FactoryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type CnStakingV4FactoryCallerRaw struct {
	Contract *CnStakingV4FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// CnStakingV4FactoryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type CnStakingV4FactoryTransactorRaw struct {
	Contract *CnStakingV4FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCnStakingV4Factory creates a new instance of CnStakingV4Factory, bound to a specific deployed contract.
func NewCnStakingV4Factory(address common.Address, backend bind.ContractBackend) (*CnStakingV4Factory, error) {
	contract, err := bindCnStakingV4Factory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CnStakingV4Factory{CnStakingV4FactoryCaller: CnStakingV4FactoryCaller{contract: contract}, CnStakingV4FactoryTransactor: CnStakingV4FactoryTransactor{contract: contract}, CnStakingV4FactoryFilterer: CnStakingV4FactoryFilterer{contract: contract}}, nil
}

// NewCnStakingV4FactoryCaller creates a new read-only instance of CnStakingV4Factory, bound to a specific deployed contract.
func NewCnStakingV4FactoryCaller(address common.Address, caller bind.ContractCaller) (*CnStakingV4FactoryCaller, error) {
	contract, err := bindCnStakingV4Factory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CnStakingV4FactoryCaller{contract: contract}, nil
}

// NewCnStakingV4FactoryTransactor creates a new write-only instance of CnStakingV4Factory, bound to a specific deployed contract.
func NewCnStakingV4FactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*CnStakingV4FactoryTransactor, error) {
	contract, err := bindCnStakingV4Factory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CnStakingV4FactoryTransactor{contract: contract}, nil
}

// NewCnStakingV4FactoryFilterer creates a new log filterer instance of CnStakingV4Factory, bound to a specific deployed contract.
func NewCnStakingV4FactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*CnStakingV4FactoryFilterer, error) {
	contract, err := bindCnStakingV4Factory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CnStakingV4FactoryFilterer{contract: contract}, nil
}

// bindCnStakingV4Factory binds a generic wrapper to an already deployed contract.
func bindCnStakingV4Factory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CnStakingV4FactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CnStakingV4Factory *CnStakingV4FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CnStakingV4Factory.Contract.CnStakingV4FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CnStakingV4Factory *CnStakingV4FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.CnStakingV4FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CnStakingV4Factory *CnStakingV4FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.CnStakingV4FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CnStakingV4Factory *CnStakingV4FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CnStakingV4Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CnStakingV4Factory *CnStakingV4FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CnStakingV4Factory *CnStakingV4FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.contract.Transact(opts, method, params...)
}

// DEADADDRESS is a free data retrieval call binding the contract method 0x4e6fd6c4.
//
// Solidity: function DEAD_ADDRESS() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactoryCaller) DEADADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CnStakingV4Factory.contract.Call(opts, &out, "DEAD_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DEADADDRESS is a free data retrieval call binding the contract method 0x4e6fd6c4.
//
// Solidity: function DEAD_ADDRESS() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactorySession) DEADADDRESS() (common.Address, error) {
	return _CnStakingV4Factory.Contract.DEADADDRESS(&_CnStakingV4Factory.CallOpts)
}

// DEADADDRESS is a free data retrieval call binding the contract method 0x4e6fd6c4.
//
// Solidity: function DEAD_ADDRESS() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactoryCallerSession) DEADADDRESS() (common.Address, error) {
	return _CnStakingV4Factory.Contract.DEADADDRESS(&_CnStakingV4Factory.CallOpts)
}

// INITIALLOCKUP is a free data retrieval call binding the contract method 0x0cf74c5c.
//
// Solidity: function INITIAL_LOCKUP() view returns(uint256)
func (_CnStakingV4Factory *CnStakingV4FactoryCaller) INITIALLOCKUP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CnStakingV4Factory.contract.Call(opts, &out, "INITIAL_LOCKUP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITIALLOCKUP is a free data retrieval call binding the contract method 0x0cf74c5c.
//
// Solidity: function INITIAL_LOCKUP() view returns(uint256)
func (_CnStakingV4Factory *CnStakingV4FactorySession) INITIALLOCKUP() (*big.Int, error) {
	return _CnStakingV4Factory.Contract.INITIALLOCKUP(&_CnStakingV4Factory.CallOpts)
}

// INITIALLOCKUP is a free data retrieval call binding the contract method 0x0cf74c5c.
//
// Solidity: function INITIAL_LOCKUP() view returns(uint256)
func (_CnStakingV4Factory *CnStakingV4FactoryCallerSession) INITIALLOCKUP() (*big.Int, error) {
	return _CnStakingV4Factory.Contract.INITIALLOCKUP(&_CnStakingV4Factory.CallOpts)
}

// CnStakingBeacon is a free data retrieval call binding the contract method 0x7b0e7fdd.
//
// Solidity: function cnStakingBeacon() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactoryCaller) CnStakingBeacon(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CnStakingV4Factory.contract.Call(opts, &out, "cnStakingBeacon")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CnStakingBeacon is a free data retrieval call binding the contract method 0x7b0e7fdd.
//
// Solidity: function cnStakingBeacon() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactorySession) CnStakingBeacon() (common.Address, error) {
	return _CnStakingV4Factory.Contract.CnStakingBeacon(&_CnStakingV4Factory.CallOpts)
}

// CnStakingBeacon is a free data retrieval call binding the contract method 0x7b0e7fdd.
//
// Solidity: function cnStakingBeacon() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactoryCallerSession) CnStakingBeacon() (common.Address, error) {
	return _CnStakingV4Factory.Contract.CnStakingBeacon(&_CnStakingV4Factory.CallOpts)
}

// GetDeployer is a free data retrieval call binding the contract method 0x669d8d45.
//
// Solidity: function getDeployer(address _addr) view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactoryCaller) GetDeployer(opts *bind.CallOpts, _addr common.Address) (common.Address, error) {
	var out []interface{}
	err := _CnStakingV4Factory.contract.Call(opts, &out, "getDeployer", _addr)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDeployer is a free data retrieval call binding the contract method 0x669d8d45.
//
// Solidity: function getDeployer(address _addr) view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactorySession) GetDeployer(_addr common.Address) (common.Address, error) {
	return _CnStakingV4Factory.Contract.GetDeployer(&_CnStakingV4Factory.CallOpts, _addr)
}

// GetDeployer is a free data retrieval call binding the contract method 0x669d8d45.
//
// Solidity: function getDeployer(address _addr) view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactoryCallerSession) GetDeployer(_addr common.Address) (common.Address, error) {
	return _CnStakingV4Factory.Contract.GetDeployer(&_CnStakingV4Factory.CallOpts, _addr)
}

// IsDeployedCnStaking is a free data retrieval call binding the contract method 0x9a429925.
//
// Solidity: function isDeployedCnStaking(address _addr) view returns(bool)
func (_CnStakingV4Factory *CnStakingV4FactoryCaller) IsDeployedCnStaking(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var out []interface{}
	err := _CnStakingV4Factory.contract.Call(opts, &out, "isDeployedCnStaking", _addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDeployedCnStaking is a free data retrieval call binding the contract method 0x9a429925.
//
// Solidity: function isDeployedCnStaking(address _addr) view returns(bool)
func (_CnStakingV4Factory *CnStakingV4FactorySession) IsDeployedCnStaking(_addr common.Address) (bool, error) {
	return _CnStakingV4Factory.Contract.IsDeployedCnStaking(&_CnStakingV4Factory.CallOpts, _addr)
}

// IsDeployedCnStaking is a free data retrieval call binding the contract method 0x9a429925.
//
// Solidity: function isDeployedCnStaking(address _addr) view returns(bool)
func (_CnStakingV4Factory *CnStakingV4FactoryCallerSession) IsDeployedCnStaking(_addr common.Address) (bool, error) {
	return _CnStakingV4Factory.Contract.IsDeployedCnStaking(&_CnStakingV4Factory.CallOpts, _addr)
}

// IsDeployedPublicDelegation is a free data retrieval call binding the contract method 0x7dfe297c.
//
// Solidity: function isDeployedPublicDelegation(address _addr) view returns(bool)
func (_CnStakingV4Factory *CnStakingV4FactoryCaller) IsDeployedPublicDelegation(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var out []interface{}
	err := _CnStakingV4Factory.contract.Call(opts, &out, "isDeployedPublicDelegation", _addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDeployedPublicDelegation is a free data retrieval call binding the contract method 0x7dfe297c.
//
// Solidity: function isDeployedPublicDelegation(address _addr) view returns(bool)
func (_CnStakingV4Factory *CnStakingV4FactorySession) IsDeployedPublicDelegation(_addr common.Address) (bool, error) {
	return _CnStakingV4Factory.Contract.IsDeployedPublicDelegation(&_CnStakingV4Factory.CallOpts, _addr)
}

// IsDeployedPublicDelegation is a free data retrieval call binding the contract method 0x7dfe297c.
//
// Solidity: function isDeployedPublicDelegation(address _addr) view returns(bool)
func (_CnStakingV4Factory *CnStakingV4FactoryCallerSession) IsDeployedPublicDelegation(_addr common.Address) (bool, error) {
	return _CnStakingV4Factory.Contract.IsDeployedPublicDelegation(&_CnStakingV4Factory.CallOpts, _addr)
}

// PdBeacon is a free data retrieval call binding the contract method 0xaa777f4c.
//
// Solidity: function pdBeacon() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactoryCaller) PdBeacon(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CnStakingV4Factory.contract.Call(opts, &out, "pdBeacon")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PdBeacon is a free data retrieval call binding the contract method 0xaa777f4c.
//
// Solidity: function pdBeacon() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactorySession) PdBeacon() (common.Address, error) {
	return _CnStakingV4Factory.Contract.PdBeacon(&_CnStakingV4Factory.CallOpts)
}

// PdBeacon is a free data retrieval call binding the contract method 0xaa777f4c.
//
// Solidity: function pdBeacon() view returns(address)
func (_CnStakingV4Factory *CnStakingV4FactoryCallerSession) PdBeacon() (common.Address, error) {
	return _CnStakingV4Factory.Contract.PdBeacon(&_CnStakingV4Factory.CallOpts)
}

// DeployCnStaking is a paid mutator transaction binding the contract method 0x4ed7a764.
//
// Solidity: function deployCnStaking(address _owner) returns(address proxy)
func (_CnStakingV4Factory *CnStakingV4FactoryTransactor) DeployCnStaking(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _CnStakingV4Factory.contract.Transact(opts, "deployCnStaking", _owner)
}

// DeployCnStaking is a paid mutator transaction binding the contract method 0x4ed7a764.
//
// Solidity: function deployCnStaking(address _owner) returns(address proxy)
func (_CnStakingV4Factory *CnStakingV4FactorySession) DeployCnStaking(_owner common.Address) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.DeployCnStaking(&_CnStakingV4Factory.TransactOpts, _owner)
}

// DeployCnStaking is a paid mutator transaction binding the contract method 0x4ed7a764.
//
// Solidity: function deployCnStaking(address _owner) returns(address proxy)
func (_CnStakingV4Factory *CnStakingV4FactoryTransactorSession) DeployCnStaking(_owner common.Address) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.DeployCnStaking(&_CnStakingV4Factory.TransactOpts, _owner)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address proxy, address publicDelegation)
func (_CnStakingV4Factory *CnStakingV4FactoryTransactor) DeployCnStakingWithPD(opts *bind.TransactOpts, _owner common.Address, _pdArgs IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _CnStakingV4Factory.contract.Transact(opts, "deployCnStakingWithPD", _owner, _pdArgs)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address proxy, address publicDelegation)
func (_CnStakingV4Factory *CnStakingV4FactorySession) DeployCnStakingWithPD(_owner common.Address, _pdArgs IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.DeployCnStakingWithPD(&_CnStakingV4Factory.TransactOpts, _owner, _pdArgs)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address proxy, address publicDelegation)
func (_CnStakingV4Factory *CnStakingV4FactoryTransactorSession) DeployCnStakingWithPD(_owner common.Address, _pdArgs IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.DeployCnStakingWithPD(&_CnStakingV4Factory.TransactOpts, _owner, _pdArgs)
}

// CnStakingV4FactoryDeployCnStakingV4Iterator is returned from FilterDeployCnStakingV4 and is used to iterate over the raw logs and unpacked data for DeployCnStakingV4 events raised by the CnStakingV4Factory contract.
type CnStakingV4FactoryDeployCnStakingV4Iterator struct {
	Event *CnStakingV4FactoryDeployCnStakingV4 // Event containing the contract specifics and raw log

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
func (it *CnStakingV4FactoryDeployCnStakingV4Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CnStakingV4FactoryDeployCnStakingV4)
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
		it.Event = new(CnStakingV4FactoryDeployCnStakingV4)
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
func (it *CnStakingV4FactoryDeployCnStakingV4Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CnStakingV4FactoryDeployCnStakingV4Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CnStakingV4FactoryDeployCnStakingV4 represents a DeployCnStakingV4 event raised by the CnStakingV4Factory contract.
type CnStakingV4FactoryDeployCnStakingV4 struct {
	Proxy common.Address
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDeployCnStakingV4 is a free log retrieval operation binding the contract event 0x1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef.
//
// Solidity: event DeployCnStakingV4(address indexed proxy, address indexed owner)
func (_CnStakingV4Factory *CnStakingV4FactoryFilterer) FilterDeployCnStakingV4(opts *bind.FilterOpts, proxy []common.Address, owner []common.Address) (*CnStakingV4FactoryDeployCnStakingV4Iterator, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CnStakingV4Factory.contract.FilterLogs(opts, "DeployCnStakingV4", proxyRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &CnStakingV4FactoryDeployCnStakingV4Iterator{contract: _CnStakingV4Factory.contract, event: "DeployCnStakingV4", logs: logs, sub: sub}, nil
}

// WatchDeployCnStakingV4 is a free log subscription operation binding the contract event 0x1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef.
//
// Solidity: event DeployCnStakingV4(address indexed proxy, address indexed owner)
func (_CnStakingV4Factory *CnStakingV4FactoryFilterer) WatchDeployCnStakingV4(opts *bind.WatchOpts, sink chan<- *CnStakingV4FactoryDeployCnStakingV4, proxy []common.Address, owner []common.Address) (event.Subscription, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CnStakingV4Factory.contract.WatchLogs(opts, "DeployCnStakingV4", proxyRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CnStakingV4FactoryDeployCnStakingV4)
				if err := _CnStakingV4Factory.contract.UnpackLog(event, "DeployCnStakingV4", log); err != nil {
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

// ParseDeployCnStakingV4 is a log parse operation binding the contract event 0x1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef.
//
// Solidity: event DeployCnStakingV4(address indexed proxy, address indexed owner)
func (_CnStakingV4Factory *CnStakingV4FactoryFilterer) ParseDeployCnStakingV4(log types.Log) (*CnStakingV4FactoryDeployCnStakingV4, error) {
	event := new(CnStakingV4FactoryDeployCnStakingV4)
	if err := _CnStakingV4Factory.contract.UnpackLog(event, "DeployCnStakingV4", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CnStakingV4FactoryDeployCnStakingV4WithPDIterator is returned from FilterDeployCnStakingV4WithPD and is used to iterate over the raw logs and unpacked data for DeployCnStakingV4WithPD events raised by the CnStakingV4Factory contract.
type CnStakingV4FactoryDeployCnStakingV4WithPDIterator struct {
	Event *CnStakingV4FactoryDeployCnStakingV4WithPD // Event containing the contract specifics and raw log

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
func (it *CnStakingV4FactoryDeployCnStakingV4WithPDIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CnStakingV4FactoryDeployCnStakingV4WithPD)
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
		it.Event = new(CnStakingV4FactoryDeployCnStakingV4WithPD)
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
func (it *CnStakingV4FactoryDeployCnStakingV4WithPDIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CnStakingV4FactoryDeployCnStakingV4WithPDIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CnStakingV4FactoryDeployCnStakingV4WithPD represents a DeployCnStakingV4WithPD event raised by the CnStakingV4Factory contract.
type CnStakingV4FactoryDeployCnStakingV4WithPD struct {
	Proxy            common.Address
	Owner            common.Address
	PublicDelegation common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDeployCnStakingV4WithPD is a free log retrieval operation binding the contract event 0x87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f.
//
// Solidity: event DeployCnStakingV4WithPD(address indexed proxy, address indexed owner, address publicDelegation)
func (_CnStakingV4Factory *CnStakingV4FactoryFilterer) FilterDeployCnStakingV4WithPD(opts *bind.FilterOpts, proxy []common.Address, owner []common.Address) (*CnStakingV4FactoryDeployCnStakingV4WithPDIterator, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CnStakingV4Factory.contract.FilterLogs(opts, "DeployCnStakingV4WithPD", proxyRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &CnStakingV4FactoryDeployCnStakingV4WithPDIterator{contract: _CnStakingV4Factory.contract, event: "DeployCnStakingV4WithPD", logs: logs, sub: sub}, nil
}

// WatchDeployCnStakingV4WithPD is a free log subscription operation binding the contract event 0x87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f.
//
// Solidity: event DeployCnStakingV4WithPD(address indexed proxy, address indexed owner, address publicDelegation)
func (_CnStakingV4Factory *CnStakingV4FactoryFilterer) WatchDeployCnStakingV4WithPD(opts *bind.WatchOpts, sink chan<- *CnStakingV4FactoryDeployCnStakingV4WithPD, proxy []common.Address, owner []common.Address) (event.Subscription, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CnStakingV4Factory.contract.WatchLogs(opts, "DeployCnStakingV4WithPD", proxyRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CnStakingV4FactoryDeployCnStakingV4WithPD)
				if err := _CnStakingV4Factory.contract.UnpackLog(event, "DeployCnStakingV4WithPD", log); err != nil {
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

// ParseDeployCnStakingV4WithPD is a log parse operation binding the contract event 0x87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f.
//
// Solidity: event DeployCnStakingV4WithPD(address indexed proxy, address indexed owner, address publicDelegation)
func (_CnStakingV4Factory *CnStakingV4FactoryFilterer) ParseDeployCnStakingV4WithPD(log types.Log) (*CnStakingV4FactoryDeployCnStakingV4WithPD, error) {
	event := new(CnStakingV4FactoryDeployCnStakingV4WithPD)
	if err := _CnStakingV4Factory.contract.UnpackLog(event, "DeployCnStakingV4WithPD", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC1967UtilsMetaData contains all meta data concerning the ERC1967Utils contract.
var ERC1967UtilsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidAdmin\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidBeacon\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"}]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220dd033e6377b33b7a96b145f6c106e3ab52adadbda2ef4b30b7d04a45d2523b7464736f6c63430008190033",
}

// ERC1967UtilsABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1967UtilsMetaData.ABI instead.
var ERC1967UtilsABI = ERC1967UtilsMetaData.ABI

// ERC1967UtilsBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC1967UtilsBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220dd033e6377b33b7a96b145f6c106e3ab52adadbda2ef4b30b7d04a45d2523b7464736f6c63430008190033`

// ERC1967UtilsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC1967UtilsMetaData.Bin instead.
var ERC1967UtilsBin = ERC1967UtilsMetaData.Bin

// DeployERC1967Utils deploys a new Kaia contract, binding an instance of ERC1967Utils to it.
func DeployERC1967Utils(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ERC1967Utils, error) {
	parsed, err := ERC1967UtilsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC1967UtilsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC1967Utils{ERC1967UtilsCaller: ERC1967UtilsCaller{contract: contract}, ERC1967UtilsTransactor: ERC1967UtilsTransactor{contract: contract}, ERC1967UtilsFilterer: ERC1967UtilsFilterer{contract: contract}}, nil
}

// ERC1967Utils is an auto generated Go binding around a Kaia contract.
type ERC1967Utils struct {
	ERC1967UtilsCaller     // Read-only binding to the contract
	ERC1967UtilsTransactor // Write-only binding to the contract
	ERC1967UtilsFilterer   // Log filterer for contract events
}

// ERC1967UtilsCaller is an auto generated read-only Go binding around a Kaia contract.
type ERC1967UtilsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UtilsTransactor is an auto generated write-only Go binding around a Kaia contract.
type ERC1967UtilsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UtilsFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ERC1967UtilsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UtilsSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ERC1967UtilsSession struct {
	Contract     *ERC1967Utils     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC1967UtilsCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ERC1967UtilsCallerSession struct {
	Contract *ERC1967UtilsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ERC1967UtilsTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ERC1967UtilsTransactorSession struct {
	Contract     *ERC1967UtilsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ERC1967UtilsRaw is an auto generated low-level Go binding around a Kaia contract.
type ERC1967UtilsRaw struct {
	Contract *ERC1967Utils // Generic contract binding to access the raw methods on
}

// ERC1967UtilsCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ERC1967UtilsCallerRaw struct {
	Contract *ERC1967UtilsCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1967UtilsTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ERC1967UtilsTransactorRaw struct {
	Contract *ERC1967UtilsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1967Utils creates a new instance of ERC1967Utils, bound to a specific deployed contract.
func NewERC1967Utils(address common.Address, backend bind.ContractBackend) (*ERC1967Utils, error) {
	contract, err := bindERC1967Utils(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1967Utils{ERC1967UtilsCaller: ERC1967UtilsCaller{contract: contract}, ERC1967UtilsTransactor: ERC1967UtilsTransactor{contract: contract}, ERC1967UtilsFilterer: ERC1967UtilsFilterer{contract: contract}}, nil
}

// NewERC1967UtilsCaller creates a new read-only instance of ERC1967Utils, bound to a specific deployed contract.
func NewERC1967UtilsCaller(address common.Address, caller bind.ContractCaller) (*ERC1967UtilsCaller, error) {
	contract, err := bindERC1967Utils(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UtilsCaller{contract: contract}, nil
}

// NewERC1967UtilsTransactor creates a new write-only instance of ERC1967Utils, bound to a specific deployed contract.
func NewERC1967UtilsTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1967UtilsTransactor, error) {
	contract, err := bindERC1967Utils(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UtilsTransactor{contract: contract}, nil
}

// NewERC1967UtilsFilterer creates a new log filterer instance of ERC1967Utils, bound to a specific deployed contract.
func NewERC1967UtilsFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1967UtilsFilterer, error) {
	contract, err := bindERC1967Utils(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1967UtilsFilterer{contract: contract}, nil
}

// bindERC1967Utils binds a generic wrapper to an already deployed contract.
func bindERC1967Utils(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1967UtilsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Utils *ERC1967UtilsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Utils.Contract.ERC1967UtilsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Utils *ERC1967UtilsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Utils.Contract.ERC1967UtilsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Utils *ERC1967UtilsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Utils.Contract.ERC1967UtilsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Utils *ERC1967UtilsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Utils.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Utils *ERC1967UtilsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Utils.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Utils *ERC1967UtilsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Utils.Contract.contract.Transact(opts, method, params...)
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

// IBeaconMetaData contains all meta data concerning the IBeacon contract.
var IBeaconMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"5c60da1b": "implementation()",
	},
}

// IBeaconABI is the input ABI used to generate the binding from.
// Deprecated: Use IBeaconMetaData.ABI instead.
var IBeaconABI = IBeaconMetaData.ABI

// IBeaconBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IBeaconBinRuntime = ``

// Deprecated: Use IBeaconMetaData.Sigs instead.
// IBeaconFuncSigs maps the 4-byte function signature to its string representation.
var IBeaconFuncSigs = IBeaconMetaData.Sigs

// IBeacon is an auto generated Go binding around a Kaia contract.
type IBeacon struct {
	IBeaconCaller     // Read-only binding to the contract
	IBeaconTransactor // Write-only binding to the contract
	IBeaconFilterer   // Log filterer for contract events
}

// IBeaconCaller is an auto generated read-only Go binding around a Kaia contract.
type IBeaconCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconTransactor is an auto generated write-only Go binding around a Kaia contract.
type IBeaconTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IBeaconFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IBeaconSession struct {
	Contract     *IBeacon          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBeaconCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IBeaconCallerSession struct {
	Contract *IBeaconCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IBeaconTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IBeaconTransactorSession struct {
	Contract     *IBeaconTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IBeaconRaw is an auto generated low-level Go binding around a Kaia contract.
type IBeaconRaw struct {
	Contract *IBeacon // Generic contract binding to access the raw methods on
}

// IBeaconCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IBeaconCallerRaw struct {
	Contract *IBeaconCaller // Generic read-only contract binding to access the raw methods on
}

// IBeaconTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IBeaconTransactorRaw struct {
	Contract *IBeaconTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIBeacon creates a new instance of IBeacon, bound to a specific deployed contract.
func NewIBeacon(address common.Address, backend bind.ContractBackend) (*IBeacon, error) {
	contract, err := bindIBeacon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IBeacon{IBeaconCaller: IBeaconCaller{contract: contract}, IBeaconTransactor: IBeaconTransactor{contract: contract}, IBeaconFilterer: IBeaconFilterer{contract: contract}}, nil
}

// NewIBeaconCaller creates a new read-only instance of IBeacon, bound to a specific deployed contract.
func NewIBeaconCaller(address common.Address, caller bind.ContractCaller) (*IBeaconCaller, error) {
	contract, err := bindIBeacon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IBeaconCaller{contract: contract}, nil
}

// NewIBeaconTransactor creates a new write-only instance of IBeacon, bound to a specific deployed contract.
func NewIBeaconTransactor(address common.Address, transactor bind.ContractTransactor) (*IBeaconTransactor, error) {
	contract, err := bindIBeacon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IBeaconTransactor{contract: contract}, nil
}

// NewIBeaconFilterer creates a new log filterer instance of IBeacon, bound to a specific deployed contract.
func NewIBeaconFilterer(address common.Address, filterer bind.ContractFilterer) (*IBeaconFilterer, error) {
	contract, err := bindIBeacon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IBeaconFilterer{contract: contract}, nil
}

// bindIBeacon binds a generic wrapper to an already deployed contract.
func bindIBeacon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IBeaconMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBeacon *IBeaconRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBeacon.Contract.IBeaconCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBeacon *IBeaconRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBeacon.Contract.IBeaconTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBeacon *IBeaconRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBeacon.Contract.IBeaconTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IBeacon *IBeaconCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IBeacon.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IBeacon *IBeaconTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IBeacon.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IBeacon *IBeaconTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IBeacon.Contract.contract.Transact(opts, method, params...)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeacon *IBeaconCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IBeacon.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeacon *IBeaconSession) Implementation() (common.Address, error) {
	return _IBeacon.Contract.Implementation(&_IBeacon.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IBeacon *IBeaconCallerSession) Implementation() (common.Address, error) {
	return _IBeacon.Contract.Implementation(&_IBeacon.CallOpts)
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

// ICnStakingV4FactoryMetaData contains all meta data concerning the ICnStakingV4Factory contract.
var ICnStakingV4FactoryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"InsufficientInitialStake\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"DeployCnStakingV4\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"publicDelegation\",\"type\":\"address\"}],\"name\":\"DeployCnStakingV4WithPD\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEAD_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIAL_LOCKUP\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cnStakingBeacon\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"deployCnStaking\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"commissionTo\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"gcName\",\"type\":\"string\"}],\"internalType\":\"structIPublicDelegation.PDConstructorArgs\",\"name\":\"_pdArgs\",\"type\":\"tuple\"}],\"name\":\"deployCnStakingWithPD\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proxy\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"publicDelegation\",\"type\":\"address\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"getDeployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isDeployedCnStaking\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isDeployedPublicDelegation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pdBeacon\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"4e6fd6c4": "DEAD_ADDRESS()",
		"0cf74c5c": "INITIAL_LOCKUP()",
		"7b0e7fdd": "cnStakingBeacon()",
		"4ed7a764": "deployCnStaking(address)",
		"33f5db27": "deployCnStakingWithPD(address,(address,address,uint256,string))",
		"669d8d45": "getDeployer(address)",
		"9a429925": "isDeployedCnStaking(address)",
		"7dfe297c": "isDeployedPublicDelegation(address)",
		"aa777f4c": "pdBeacon()",
	},
}

// ICnStakingV4FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use ICnStakingV4FactoryMetaData.ABI instead.
var ICnStakingV4FactoryABI = ICnStakingV4FactoryMetaData.ABI

// ICnStakingV4FactoryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ICnStakingV4FactoryBinRuntime = ``

// Deprecated: Use ICnStakingV4FactoryMetaData.Sigs instead.
// ICnStakingV4FactoryFuncSigs maps the 4-byte function signature to its string representation.
var ICnStakingV4FactoryFuncSigs = ICnStakingV4FactoryMetaData.Sigs

// ICnStakingV4Factory is an auto generated Go binding around a Kaia contract.
type ICnStakingV4Factory struct {
	ICnStakingV4FactoryCaller     // Read-only binding to the contract
	ICnStakingV4FactoryTransactor // Write-only binding to the contract
	ICnStakingV4FactoryFilterer   // Log filterer for contract events
}

// ICnStakingV4FactoryCaller is an auto generated read-only Go binding around a Kaia contract.
type ICnStakingV4FactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICnStakingV4FactoryTransactor is an auto generated write-only Go binding around a Kaia contract.
type ICnStakingV4FactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICnStakingV4FactoryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ICnStakingV4FactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ICnStakingV4FactorySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ICnStakingV4FactorySession struct {
	Contract     *ICnStakingV4Factory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ICnStakingV4FactoryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ICnStakingV4FactoryCallerSession struct {
	Contract *ICnStakingV4FactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// ICnStakingV4FactoryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ICnStakingV4FactoryTransactorSession struct {
	Contract     *ICnStakingV4FactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// ICnStakingV4FactoryRaw is an auto generated low-level Go binding around a Kaia contract.
type ICnStakingV4FactoryRaw struct {
	Contract *ICnStakingV4Factory // Generic contract binding to access the raw methods on
}

// ICnStakingV4FactoryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ICnStakingV4FactoryCallerRaw struct {
	Contract *ICnStakingV4FactoryCaller // Generic read-only contract binding to access the raw methods on
}

// ICnStakingV4FactoryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ICnStakingV4FactoryTransactorRaw struct {
	Contract *ICnStakingV4FactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewICnStakingV4Factory creates a new instance of ICnStakingV4Factory, bound to a specific deployed contract.
func NewICnStakingV4Factory(address common.Address, backend bind.ContractBackend) (*ICnStakingV4Factory, error) {
	contract, err := bindICnStakingV4Factory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ICnStakingV4Factory{ICnStakingV4FactoryCaller: ICnStakingV4FactoryCaller{contract: contract}, ICnStakingV4FactoryTransactor: ICnStakingV4FactoryTransactor{contract: contract}, ICnStakingV4FactoryFilterer: ICnStakingV4FactoryFilterer{contract: contract}}, nil
}

// NewICnStakingV4FactoryCaller creates a new read-only instance of ICnStakingV4Factory, bound to a specific deployed contract.
func NewICnStakingV4FactoryCaller(address common.Address, caller bind.ContractCaller) (*ICnStakingV4FactoryCaller, error) {
	contract, err := bindICnStakingV4Factory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ICnStakingV4FactoryCaller{contract: contract}, nil
}

// NewICnStakingV4FactoryTransactor creates a new write-only instance of ICnStakingV4Factory, bound to a specific deployed contract.
func NewICnStakingV4FactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*ICnStakingV4FactoryTransactor, error) {
	contract, err := bindICnStakingV4Factory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ICnStakingV4FactoryTransactor{contract: contract}, nil
}

// NewICnStakingV4FactoryFilterer creates a new log filterer instance of ICnStakingV4Factory, bound to a specific deployed contract.
func NewICnStakingV4FactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*ICnStakingV4FactoryFilterer, error) {
	contract, err := bindICnStakingV4Factory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ICnStakingV4FactoryFilterer{contract: contract}, nil
}

// bindICnStakingV4Factory binds a generic wrapper to an already deployed contract.
func bindICnStakingV4Factory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ICnStakingV4FactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICnStakingV4Factory *ICnStakingV4FactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICnStakingV4Factory.Contract.ICnStakingV4FactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICnStakingV4Factory *ICnStakingV4FactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICnStakingV4Factory.Contract.ICnStakingV4FactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICnStakingV4Factory *ICnStakingV4FactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICnStakingV4Factory.Contract.ICnStakingV4FactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ICnStakingV4Factory *ICnStakingV4FactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ICnStakingV4Factory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ICnStakingV4Factory *ICnStakingV4FactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ICnStakingV4Factory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ICnStakingV4Factory *ICnStakingV4FactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ICnStakingV4Factory.Contract.contract.Transact(opts, method, params...)
}

// DEADADDRESS is a free data retrieval call binding the contract method 0x4e6fd6c4.
//
// Solidity: function DEAD_ADDRESS() pure returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCaller) DEADADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ICnStakingV4Factory.contract.Call(opts, &out, "DEAD_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DEADADDRESS is a free data retrieval call binding the contract method 0x4e6fd6c4.
//
// Solidity: function DEAD_ADDRESS() pure returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) DEADADDRESS() (common.Address, error) {
	return _ICnStakingV4Factory.Contract.DEADADDRESS(&_ICnStakingV4Factory.CallOpts)
}

// DEADADDRESS is a free data retrieval call binding the contract method 0x4e6fd6c4.
//
// Solidity: function DEAD_ADDRESS() pure returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCallerSession) DEADADDRESS() (common.Address, error) {
	return _ICnStakingV4Factory.Contract.DEADADDRESS(&_ICnStakingV4Factory.CallOpts)
}

// INITIALLOCKUP is a free data retrieval call binding the contract method 0x0cf74c5c.
//
// Solidity: function INITIAL_LOCKUP() pure returns(uint256)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCaller) INITIALLOCKUP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ICnStakingV4Factory.contract.Call(opts, &out, "INITIAL_LOCKUP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// INITIALLOCKUP is a free data retrieval call binding the contract method 0x0cf74c5c.
//
// Solidity: function INITIAL_LOCKUP() pure returns(uint256)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) INITIALLOCKUP() (*big.Int, error) {
	return _ICnStakingV4Factory.Contract.INITIALLOCKUP(&_ICnStakingV4Factory.CallOpts)
}

// INITIALLOCKUP is a free data retrieval call binding the contract method 0x0cf74c5c.
//
// Solidity: function INITIAL_LOCKUP() pure returns(uint256)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCallerSession) INITIALLOCKUP() (*big.Int, error) {
	return _ICnStakingV4Factory.Contract.INITIALLOCKUP(&_ICnStakingV4Factory.CallOpts)
}

// CnStakingBeacon is a free data retrieval call binding the contract method 0x7b0e7fdd.
//
// Solidity: function cnStakingBeacon() view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCaller) CnStakingBeacon(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ICnStakingV4Factory.contract.Call(opts, &out, "cnStakingBeacon")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CnStakingBeacon is a free data retrieval call binding the contract method 0x7b0e7fdd.
//
// Solidity: function cnStakingBeacon() view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) CnStakingBeacon() (common.Address, error) {
	return _ICnStakingV4Factory.Contract.CnStakingBeacon(&_ICnStakingV4Factory.CallOpts)
}

// CnStakingBeacon is a free data retrieval call binding the contract method 0x7b0e7fdd.
//
// Solidity: function cnStakingBeacon() view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCallerSession) CnStakingBeacon() (common.Address, error) {
	return _ICnStakingV4Factory.Contract.CnStakingBeacon(&_ICnStakingV4Factory.CallOpts)
}

// GetDeployer is a free data retrieval call binding the contract method 0x669d8d45.
//
// Solidity: function getDeployer(address _addr) view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCaller) GetDeployer(opts *bind.CallOpts, _addr common.Address) (common.Address, error) {
	var out []interface{}
	err := _ICnStakingV4Factory.contract.Call(opts, &out, "getDeployer", _addr)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDeployer is a free data retrieval call binding the contract method 0x669d8d45.
//
// Solidity: function getDeployer(address _addr) view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) GetDeployer(_addr common.Address) (common.Address, error) {
	return _ICnStakingV4Factory.Contract.GetDeployer(&_ICnStakingV4Factory.CallOpts, _addr)
}

// GetDeployer is a free data retrieval call binding the contract method 0x669d8d45.
//
// Solidity: function getDeployer(address _addr) view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCallerSession) GetDeployer(_addr common.Address) (common.Address, error) {
	return _ICnStakingV4Factory.Contract.GetDeployer(&_ICnStakingV4Factory.CallOpts, _addr)
}

// IsDeployedCnStaking is a free data retrieval call binding the contract method 0x9a429925.
//
// Solidity: function isDeployedCnStaking(address _addr) view returns(bool)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCaller) IsDeployedCnStaking(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var out []interface{}
	err := _ICnStakingV4Factory.contract.Call(opts, &out, "isDeployedCnStaking", _addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDeployedCnStaking is a free data retrieval call binding the contract method 0x9a429925.
//
// Solidity: function isDeployedCnStaking(address _addr) view returns(bool)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) IsDeployedCnStaking(_addr common.Address) (bool, error) {
	return _ICnStakingV4Factory.Contract.IsDeployedCnStaking(&_ICnStakingV4Factory.CallOpts, _addr)
}

// IsDeployedCnStaking is a free data retrieval call binding the contract method 0x9a429925.
//
// Solidity: function isDeployedCnStaking(address _addr) view returns(bool)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCallerSession) IsDeployedCnStaking(_addr common.Address) (bool, error) {
	return _ICnStakingV4Factory.Contract.IsDeployedCnStaking(&_ICnStakingV4Factory.CallOpts, _addr)
}

// IsDeployedPublicDelegation is a free data retrieval call binding the contract method 0x7dfe297c.
//
// Solidity: function isDeployedPublicDelegation(address _addr) view returns(bool)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCaller) IsDeployedPublicDelegation(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var out []interface{}
	err := _ICnStakingV4Factory.contract.Call(opts, &out, "isDeployedPublicDelegation", _addr)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsDeployedPublicDelegation is a free data retrieval call binding the contract method 0x7dfe297c.
//
// Solidity: function isDeployedPublicDelegation(address _addr) view returns(bool)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) IsDeployedPublicDelegation(_addr common.Address) (bool, error) {
	return _ICnStakingV4Factory.Contract.IsDeployedPublicDelegation(&_ICnStakingV4Factory.CallOpts, _addr)
}

// IsDeployedPublicDelegation is a free data retrieval call binding the contract method 0x7dfe297c.
//
// Solidity: function isDeployedPublicDelegation(address _addr) view returns(bool)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCallerSession) IsDeployedPublicDelegation(_addr common.Address) (bool, error) {
	return _ICnStakingV4Factory.Contract.IsDeployedPublicDelegation(&_ICnStakingV4Factory.CallOpts, _addr)
}

// PdBeacon is a free data retrieval call binding the contract method 0xaa777f4c.
//
// Solidity: function pdBeacon() view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCaller) PdBeacon(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ICnStakingV4Factory.contract.Call(opts, &out, "pdBeacon")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PdBeacon is a free data retrieval call binding the contract method 0xaa777f4c.
//
// Solidity: function pdBeacon() view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) PdBeacon() (common.Address, error) {
	return _ICnStakingV4Factory.Contract.PdBeacon(&_ICnStakingV4Factory.CallOpts)
}

// PdBeacon is a free data retrieval call binding the contract method 0xaa777f4c.
//
// Solidity: function pdBeacon() view returns(address)
func (_ICnStakingV4Factory *ICnStakingV4FactoryCallerSession) PdBeacon() (common.Address, error) {
	return _ICnStakingV4Factory.Contract.PdBeacon(&_ICnStakingV4Factory.CallOpts)
}

// DeployCnStaking is a paid mutator transaction binding the contract method 0x4ed7a764.
//
// Solidity: function deployCnStaking(address _owner) returns(address proxy)
func (_ICnStakingV4Factory *ICnStakingV4FactoryTransactor) DeployCnStaking(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _ICnStakingV4Factory.contract.Transact(opts, "deployCnStaking", _owner)
}

// DeployCnStaking is a paid mutator transaction binding the contract method 0x4ed7a764.
//
// Solidity: function deployCnStaking(address _owner) returns(address proxy)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) DeployCnStaking(_owner common.Address) (*types.Transaction, error) {
	return _ICnStakingV4Factory.Contract.DeployCnStaking(&_ICnStakingV4Factory.TransactOpts, _owner)
}

// DeployCnStaking is a paid mutator transaction binding the contract method 0x4ed7a764.
//
// Solidity: function deployCnStaking(address _owner) returns(address proxy)
func (_ICnStakingV4Factory *ICnStakingV4FactoryTransactorSession) DeployCnStaking(_owner common.Address) (*types.Transaction, error) {
	return _ICnStakingV4Factory.Contract.DeployCnStaking(&_ICnStakingV4Factory.TransactOpts, _owner)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address proxy, address publicDelegation)
func (_ICnStakingV4Factory *ICnStakingV4FactoryTransactor) DeployCnStakingWithPD(opts *bind.TransactOpts, _owner common.Address, _pdArgs IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _ICnStakingV4Factory.contract.Transact(opts, "deployCnStakingWithPD", _owner, _pdArgs)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address proxy, address publicDelegation)
func (_ICnStakingV4Factory *ICnStakingV4FactorySession) DeployCnStakingWithPD(_owner common.Address, _pdArgs IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _ICnStakingV4Factory.Contract.DeployCnStakingWithPD(&_ICnStakingV4Factory.TransactOpts, _owner, _pdArgs)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address proxy, address publicDelegation)
func (_ICnStakingV4Factory *ICnStakingV4FactoryTransactorSession) DeployCnStakingWithPD(_owner common.Address, _pdArgs IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _ICnStakingV4Factory.Contract.DeployCnStakingWithPD(&_ICnStakingV4Factory.TransactOpts, _owner, _pdArgs)
}

// ICnStakingV4FactoryDeployCnStakingV4Iterator is returned from FilterDeployCnStakingV4 and is used to iterate over the raw logs and unpacked data for DeployCnStakingV4 events raised by the ICnStakingV4Factory contract.
type ICnStakingV4FactoryDeployCnStakingV4Iterator struct {
	Event *ICnStakingV4FactoryDeployCnStakingV4 // Event containing the contract specifics and raw log

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
func (it *ICnStakingV4FactoryDeployCnStakingV4Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingV4FactoryDeployCnStakingV4)
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
		it.Event = new(ICnStakingV4FactoryDeployCnStakingV4)
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
func (it *ICnStakingV4FactoryDeployCnStakingV4Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingV4FactoryDeployCnStakingV4Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingV4FactoryDeployCnStakingV4 represents a DeployCnStakingV4 event raised by the ICnStakingV4Factory contract.
type ICnStakingV4FactoryDeployCnStakingV4 struct {
	Proxy common.Address
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterDeployCnStakingV4 is a free log retrieval operation binding the contract event 0x1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef.
//
// Solidity: event DeployCnStakingV4(address indexed proxy, address indexed owner)
func (_ICnStakingV4Factory *ICnStakingV4FactoryFilterer) FilterDeployCnStakingV4(opts *bind.FilterOpts, proxy []common.Address, owner []common.Address) (*ICnStakingV4FactoryDeployCnStakingV4Iterator, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ICnStakingV4Factory.contract.FilterLogs(opts, "DeployCnStakingV4", proxyRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingV4FactoryDeployCnStakingV4Iterator{contract: _ICnStakingV4Factory.contract, event: "DeployCnStakingV4", logs: logs, sub: sub}, nil
}

// WatchDeployCnStakingV4 is a free log subscription operation binding the contract event 0x1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef.
//
// Solidity: event DeployCnStakingV4(address indexed proxy, address indexed owner)
func (_ICnStakingV4Factory *ICnStakingV4FactoryFilterer) WatchDeployCnStakingV4(opts *bind.WatchOpts, sink chan<- *ICnStakingV4FactoryDeployCnStakingV4, proxy []common.Address, owner []common.Address) (event.Subscription, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ICnStakingV4Factory.contract.WatchLogs(opts, "DeployCnStakingV4", proxyRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingV4FactoryDeployCnStakingV4)
				if err := _ICnStakingV4Factory.contract.UnpackLog(event, "DeployCnStakingV4", log); err != nil {
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

// ParseDeployCnStakingV4 is a log parse operation binding the contract event 0x1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef.
//
// Solidity: event DeployCnStakingV4(address indexed proxy, address indexed owner)
func (_ICnStakingV4Factory *ICnStakingV4FactoryFilterer) ParseDeployCnStakingV4(log types.Log) (*ICnStakingV4FactoryDeployCnStakingV4, error) {
	event := new(ICnStakingV4FactoryDeployCnStakingV4)
	if err := _ICnStakingV4Factory.contract.UnpackLog(event, "DeployCnStakingV4", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ICnStakingV4FactoryDeployCnStakingV4WithPDIterator is returned from FilterDeployCnStakingV4WithPD and is used to iterate over the raw logs and unpacked data for DeployCnStakingV4WithPD events raised by the ICnStakingV4Factory contract.
type ICnStakingV4FactoryDeployCnStakingV4WithPDIterator struct {
	Event *ICnStakingV4FactoryDeployCnStakingV4WithPD // Event containing the contract specifics and raw log

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
func (it *ICnStakingV4FactoryDeployCnStakingV4WithPDIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ICnStakingV4FactoryDeployCnStakingV4WithPD)
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
		it.Event = new(ICnStakingV4FactoryDeployCnStakingV4WithPD)
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
func (it *ICnStakingV4FactoryDeployCnStakingV4WithPDIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ICnStakingV4FactoryDeployCnStakingV4WithPDIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ICnStakingV4FactoryDeployCnStakingV4WithPD represents a DeployCnStakingV4WithPD event raised by the ICnStakingV4Factory contract.
type ICnStakingV4FactoryDeployCnStakingV4WithPD struct {
	Proxy            common.Address
	Owner            common.Address
	PublicDelegation common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterDeployCnStakingV4WithPD is a free log retrieval operation binding the contract event 0x87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f.
//
// Solidity: event DeployCnStakingV4WithPD(address indexed proxy, address indexed owner, address publicDelegation)
func (_ICnStakingV4Factory *ICnStakingV4FactoryFilterer) FilterDeployCnStakingV4WithPD(opts *bind.FilterOpts, proxy []common.Address, owner []common.Address) (*ICnStakingV4FactoryDeployCnStakingV4WithPDIterator, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ICnStakingV4Factory.contract.FilterLogs(opts, "DeployCnStakingV4WithPD", proxyRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &ICnStakingV4FactoryDeployCnStakingV4WithPDIterator{contract: _ICnStakingV4Factory.contract, event: "DeployCnStakingV4WithPD", logs: logs, sub: sub}, nil
}

// WatchDeployCnStakingV4WithPD is a free log subscription operation binding the contract event 0x87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f.
//
// Solidity: event DeployCnStakingV4WithPD(address indexed proxy, address indexed owner, address publicDelegation)
func (_ICnStakingV4Factory *ICnStakingV4FactoryFilterer) WatchDeployCnStakingV4WithPD(opts *bind.WatchOpts, sink chan<- *ICnStakingV4FactoryDeployCnStakingV4WithPD, proxy []common.Address, owner []common.Address) (event.Subscription, error) {

	var proxyRule []interface{}
	for _, proxyItem := range proxy {
		proxyRule = append(proxyRule, proxyItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _ICnStakingV4Factory.contract.WatchLogs(opts, "DeployCnStakingV4WithPD", proxyRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ICnStakingV4FactoryDeployCnStakingV4WithPD)
				if err := _ICnStakingV4Factory.contract.UnpackLog(event, "DeployCnStakingV4WithPD", log); err != nil {
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

// ParseDeployCnStakingV4WithPD is a log parse operation binding the contract event 0x87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f.
//
// Solidity: event DeployCnStakingV4WithPD(address indexed proxy, address indexed owner, address publicDelegation)
func (_ICnStakingV4Factory *ICnStakingV4FactoryFilterer) ParseDeployCnStakingV4WithPD(log types.Log) (*ICnStakingV4FactoryDeployCnStakingV4WithPD, error) {
	event := new(ICnStakingV4FactoryDeployCnStakingV4WithPD)
	if err := _ICnStakingV4Factory.contract.UnpackLog(event, "DeployCnStakingV4WithPD", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC1967MetaData contains all meta data concerning the IERC1967 contract.
var IERC1967MetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"}]",
}

// IERC1967ABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1967MetaData.ABI instead.
var IERC1967ABI = IERC1967MetaData.ABI

// IERC1967BinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1967BinRuntime = ``

// IERC1967 is an auto generated Go binding around a Kaia contract.
type IERC1967 struct {
	IERC1967Caller     // Read-only binding to the contract
	IERC1967Transactor // Write-only binding to the contract
	IERC1967Filterer   // Log filterer for contract events
}

// IERC1967Caller is an auto generated read-only Go binding around a Kaia contract.
type IERC1967Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Transactor is an auto generated write-only Go binding around a Kaia contract.
type IERC1967Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Filterer is an auto generated log filtering Go binding around a Kaia contract events.
type IERC1967Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Session is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IERC1967Session struct {
	Contract     *IERC1967         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC1967CallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IERC1967CallerSession struct {
	Contract *IERC1967Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IERC1967TransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IERC1967TransactorSession struct {
	Contract     *IERC1967Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IERC1967Raw is an auto generated low-level Go binding around a Kaia contract.
type IERC1967Raw struct {
	Contract *IERC1967 // Generic contract binding to access the raw methods on
}

// IERC1967CallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IERC1967CallerRaw struct {
	Contract *IERC1967Caller // Generic read-only contract binding to access the raw methods on
}

// IERC1967TransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IERC1967TransactorRaw struct {
	Contract *IERC1967Transactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1967 creates a new instance of IERC1967, bound to a specific deployed contract.
func NewIERC1967(address common.Address, backend bind.ContractBackend) (*IERC1967, error) {
	contract, err := bindIERC1967(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1967{IERC1967Caller: IERC1967Caller{contract: contract}, IERC1967Transactor: IERC1967Transactor{contract: contract}, IERC1967Filterer: IERC1967Filterer{contract: contract}}, nil
}

// NewIERC1967Caller creates a new read-only instance of IERC1967, bound to a specific deployed contract.
func NewIERC1967Caller(address common.Address, caller bind.ContractCaller) (*IERC1967Caller, error) {
	contract, err := bindIERC1967(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1967Caller{contract: contract}, nil
}

// NewIERC1967Transactor creates a new write-only instance of IERC1967, bound to a specific deployed contract.
func NewIERC1967Transactor(address common.Address, transactor bind.ContractTransactor) (*IERC1967Transactor, error) {
	contract, err := bindIERC1967(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1967Transactor{contract: contract}, nil
}

// NewIERC1967Filterer creates a new log filterer instance of IERC1967, bound to a specific deployed contract.
func NewIERC1967Filterer(address common.Address, filterer bind.ContractFilterer) (*IERC1967Filterer, error) {
	contract, err := bindIERC1967(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1967Filterer{contract: contract}, nil
}

// bindIERC1967 binds a generic wrapper to an already deployed contract.
func bindIERC1967(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1967MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1967 *IERC1967Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1967.Contract.IERC1967Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1967 *IERC1967Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1967.Contract.IERC1967Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1967 *IERC1967Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1967.Contract.IERC1967Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1967 *IERC1967CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1967.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1967 *IERC1967TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1967.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1967 *IERC1967TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1967.Contract.contract.Transact(opts, method, params...)
}

// IERC1967AdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the IERC1967 contract.
type IERC1967AdminChangedIterator struct {
	Event *IERC1967AdminChanged // Event containing the contract specifics and raw log

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
func (it *IERC1967AdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967AdminChanged)
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
		it.Event = new(IERC1967AdminChanged)
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
func (it *IERC1967AdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967AdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967AdminChanged represents a AdminChanged event raised by the IERC1967 contract.
type IERC1967AdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967 *IERC1967Filterer) FilterAdminChanged(opts *bind.FilterOpts) (*IERC1967AdminChangedIterator, error) {

	logs, sub, err := _IERC1967.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &IERC1967AdminChangedIterator{contract: _IERC1967.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_IERC1967 *IERC1967Filterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *IERC1967AdminChanged) (event.Subscription, error) {

	logs, sub, err := _IERC1967.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967AdminChanged)
				if err := _IERC1967.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_IERC1967 *IERC1967Filterer) ParseAdminChanged(log types.Log) (*IERC1967AdminChanged, error) {
	event := new(IERC1967AdminChanged)
	if err := _IERC1967.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC1967BeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the IERC1967 contract.
type IERC1967BeaconUpgradedIterator struct {
	Event *IERC1967BeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *IERC1967BeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967BeaconUpgraded)
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
		it.Event = new(IERC1967BeaconUpgraded)
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
func (it *IERC1967BeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967BeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967BeaconUpgraded represents a BeaconUpgraded event raised by the IERC1967 contract.
type IERC1967BeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967 *IERC1967Filterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*IERC1967BeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _IERC1967.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &IERC1967BeaconUpgradedIterator{contract: _IERC1967.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_IERC1967 *IERC1967Filterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *IERC1967BeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _IERC1967.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967BeaconUpgraded)
				if err := _IERC1967.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_IERC1967 *IERC1967Filterer) ParseBeaconUpgraded(log types.Log) (*IERC1967BeaconUpgraded, error) {
	event := new(IERC1967BeaconUpgraded)
	if err := _IERC1967.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IERC1967UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the IERC1967 contract.
type IERC1967UpgradedIterator struct {
	Event *IERC1967Upgraded // Event containing the contract specifics and raw log

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
func (it *IERC1967UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IERC1967Upgraded)
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
		it.Event = new(IERC1967Upgraded)
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
func (it *IERC1967UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IERC1967UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IERC1967Upgraded represents a Upgraded event raised by the IERC1967 contract.
type IERC1967Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967 *IERC1967Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*IERC1967UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IERC1967.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &IERC1967UpgradedIterator{contract: _IERC1967.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_IERC1967 *IERC1967Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *IERC1967Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _IERC1967.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IERC1967Upgraded)
				if err := _IERC1967.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_IERC1967 *IERC1967Filterer) ParseUpgraded(log types.Log) (*IERC1967Upgraded, error) {
	event := new(IERC1967Upgraded)
	if err := _IERC1967.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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

// ProxyMetaData contains all meta data concerning the Proxy contract.
var ProxyMetaData = &bind.MetaData{
	ABI: "[{\"stateMutability\":\"payable\",\"type\":\"fallback\"}]",
}

// ProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use ProxyMetaData.ABI instead.
var ProxyABI = ProxyMetaData.ABI

// ProxyBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ProxyBinRuntime = ``

// Proxy is an auto generated Go binding around a Kaia contract.
type Proxy struct {
	ProxyCaller     // Read-only binding to the contract
	ProxyTransactor // Write-only binding to the contract
	ProxyFilterer   // Log filterer for contract events
}

// ProxyCaller is an auto generated read-only Go binding around a Kaia contract.
type ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyTransactor is an auto generated write-only Go binding around a Kaia contract.
type ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ProxySession struct {
	Contract     *Proxy            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ProxyCallerSession struct {
	Contract *ProxyCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ProxyTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ProxyTransactorSession struct {
	Contract     *ProxyTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyRaw is an auto generated low-level Go binding around a Kaia contract.
type ProxyRaw struct {
	Contract *Proxy // Generic contract binding to access the raw methods on
}

// ProxyCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ProxyCallerRaw struct {
	Contract *ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ProxyTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ProxyTransactorRaw struct {
	Contract *ProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProxy creates a new instance of Proxy, bound to a specific deployed contract.
func NewProxy(address common.Address, backend bind.ContractBackend) (*Proxy, error) {
	contract, err := bindProxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Proxy{ProxyCaller: ProxyCaller{contract: contract}, ProxyTransactor: ProxyTransactor{contract: contract}, ProxyFilterer: ProxyFilterer{contract: contract}}, nil
}

// NewProxyCaller creates a new read-only instance of Proxy, bound to a specific deployed contract.
func NewProxyCaller(address common.Address, caller bind.ContractCaller) (*ProxyCaller, error) {
	contract, err := bindProxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyCaller{contract: contract}, nil
}

// NewProxyTransactor creates a new write-only instance of Proxy, bound to a specific deployed contract.
func NewProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ProxyTransactor, error) {
	contract, err := bindProxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProxyTransactor{contract: contract}, nil
}

// NewProxyFilterer creates a new log filterer instance of Proxy, bound to a specific deployed contract.
func NewProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ProxyFilterer, error) {
	contract, err := bindProxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProxyFilterer{contract: contract}, nil
}

// bindProxy binds a generic wrapper to an already deployed contract.
func bindProxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proxy *ProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Proxy.Contract.ProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proxy *ProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.Contract.ProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proxy *ProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proxy.Contract.ProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proxy *ProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Proxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proxy *ProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proxy *ProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Proxy *ProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Proxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Proxy *ProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Proxy.Contract.Fallback(&_Proxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Proxy *ProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Proxy.Contract.Fallback(&_Proxy.TransactOpts, calldata)
}

// StorageSlotMetaData contains all meta data concerning the StorageSlot contract.
var StorageSlotMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220f2ed8da391a97d25a7979584f89d4134e64cc7d5147560095d1bcd8063085ba264736f6c63430008190033",
}

// StorageSlotABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageSlotMetaData.ABI instead.
var StorageSlotABI = StorageSlotMetaData.ABI

// StorageSlotBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const StorageSlotBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220f2ed8da391a97d25a7979584f89d4134e64cc7d5147560095d1bcd8063085ba264736f6c63430008190033`

// StorageSlotBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageSlotMetaData.Bin instead.
var StorageSlotBin = StorageSlotMetaData.Bin

// DeployStorageSlot deploys a new Kaia contract, binding an instance of StorageSlot to it.
func DeployStorageSlot(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *StorageSlot, error) {
	parsed, err := StorageSlotMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StorageSlotBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StorageSlot{StorageSlotCaller: StorageSlotCaller{contract: contract}, StorageSlotTransactor: StorageSlotTransactor{contract: contract}, StorageSlotFilterer: StorageSlotFilterer{contract: contract}}, nil
}

// StorageSlot is an auto generated Go binding around a Kaia contract.
type StorageSlot struct {
	StorageSlotCaller     // Read-only binding to the contract
	StorageSlotTransactor // Write-only binding to the contract
	StorageSlotFilterer   // Log filterer for contract events
}

// StorageSlotCaller is an auto generated read-only Go binding around a Kaia contract.
type StorageSlotCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotTransactor is an auto generated write-only Go binding around a Kaia contract.
type StorageSlotTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type StorageSlotFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type StorageSlotSession struct {
	Contract     *StorageSlot      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageSlotCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type StorageSlotCallerSession struct {
	Contract *StorageSlotCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StorageSlotTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type StorageSlotTransactorSession struct {
	Contract     *StorageSlotTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StorageSlotRaw is an auto generated low-level Go binding around a Kaia contract.
type StorageSlotRaw struct {
	Contract *StorageSlot // Generic contract binding to access the raw methods on
}

// StorageSlotCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type StorageSlotCallerRaw struct {
	Contract *StorageSlotCaller // Generic read-only contract binding to access the raw methods on
}

// StorageSlotTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type StorageSlotTransactorRaw struct {
	Contract *StorageSlotTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStorageSlot creates a new instance of StorageSlot, bound to a specific deployed contract.
func NewStorageSlot(address common.Address, backend bind.ContractBackend) (*StorageSlot, error) {
	contract, err := bindStorageSlot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StorageSlot{StorageSlotCaller: StorageSlotCaller{contract: contract}, StorageSlotTransactor: StorageSlotTransactor{contract: contract}, StorageSlotFilterer: StorageSlotFilterer{contract: contract}}, nil
}

// NewStorageSlotCaller creates a new read-only instance of StorageSlot, bound to a specific deployed contract.
func NewStorageSlotCaller(address common.Address, caller bind.ContractCaller) (*StorageSlotCaller, error) {
	contract, err := bindStorageSlot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSlotCaller{contract: contract}, nil
}

// NewStorageSlotTransactor creates a new write-only instance of StorageSlot, bound to a specific deployed contract.
func NewStorageSlotTransactor(address common.Address, transactor bind.ContractTransactor) (*StorageSlotTransactor, error) {
	contract, err := bindStorageSlot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StorageSlotTransactor{contract: contract}, nil
}

// NewStorageSlotFilterer creates a new log filterer instance of StorageSlot, bound to a specific deployed contract.
func NewStorageSlotFilterer(address common.Address, filterer bind.ContractFilterer) (*StorageSlotFilterer, error) {
	contract, err := bindStorageSlot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StorageSlotFilterer{contract: contract}, nil
}

// bindStorageSlot binds a generic wrapper to an already deployed contract.
func bindStorageSlot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StorageSlotMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSlot *StorageSlotRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSlot.Contract.StorageSlotCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSlot *StorageSlotRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSlot.Contract.StorageSlotTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSlot *StorageSlotRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSlot.Contract.StorageSlotTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StorageSlot *StorageSlotCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StorageSlot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StorageSlot *StorageSlotTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StorageSlot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StorageSlot *StorageSlotTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StorageSlot.Contract.contract.Transact(opts, method, params...)
}
