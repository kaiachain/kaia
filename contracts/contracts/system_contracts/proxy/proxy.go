// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package proxy

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

// AddressMetaData contains all meta data concerning the Address contract.
var AddressMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212206a5ffe3918ffd67bc8201cd7c53cdf8cf8012872ab683529e81e227957132c4364736f6c63430008130033",
}

// AddressABI is the input ABI used to generate the binding from.
// Deprecated: Use AddressMetaData.ABI instead.
var AddressABI = AddressMetaData.ABI

// AddressBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const AddressBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212206a5ffe3918ffd67bc8201cd7c53cdf8cf8012872ab683529e81e227957132c4364736f6c63430008130033`

// AddressBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AddressMetaData.Bin instead.
var AddressBin = AddressMetaData.Bin

// DeployAddress deploys a new Klaytn contract, binding an instance of Address to it.
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

// Address is an auto generated Go binding around a Klaytn contract.
type Address struct {
	AddressCaller     // Read-only binding to the contract
	AddressTransactor // Write-only binding to the contract
	AddressFilterer   // Log filterer for contract events
}

// AddressCaller is an auto generated read-only Go binding around a Klaytn contract.
type AddressCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressTransactor is an auto generated write-only Go binding around a Klaytn contract.
type AddressTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type AddressFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AddressSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type AddressSession struct {
	Contract     *Address          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AddressCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type AddressCallerSession struct {
	Contract *AddressCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// AddressTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type AddressTransactorSession struct {
	Contract     *AddressTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AddressRaw is an auto generated low-level Go binding around a Klaytn contract.
type AddressRaw struct {
	Contract *Address // Generic contract binding to access the raw methods on
}

// AddressCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type AddressCallerRaw struct {
	Contract *AddressCaller // Generic read-only contract binding to access the raw methods on
}

// AddressTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
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

// ERC1967ProxyMetaData contains all meta data concerning the ERC1967Proxy contract.
var ERC1967ProxyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x60806040526040516104e13803806104e1833981016040819052610022916102de565b61002e82826000610035565b50506103fb565b61003e83610061565b60008251118061004b5750805b1561005c5761005a83836100a1565b505b505050565b61006a816100cd565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606100c683836040518060600160405280602781526020016104ba60279139610180565b9392505050565b6001600160a01b0381163b61013f5760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b0319166001600160a01b0392909216919091179055565b6060600080856001600160a01b03168560405161019d91906103ac565b600060405180830381855af49150503d80600081146101d8576040519150601f19603f3d011682016040523d82523d6000602084013e6101dd565b606091505b5090925090506101ef868383876101f9565b9695505050505050565b60608315610268578251600003610261576001600160a01b0385163b6102615760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610136565b5081610272565b610272838361027a565b949350505050565b81511561028a5781518083602001fd5b8060405162461bcd60e51b815260040161013691906103c8565b634e487b7160e01b600052604160045260246000fd5b60005b838110156102d55781810151838201526020016102bd565b50506000910152565b600080604083850312156102f157600080fd5b82516001600160a01b038116811461030857600080fd5b60208401519092506001600160401b038082111561032557600080fd5b818501915085601f83011261033957600080fd5b81518181111561034b5761034b6102a4565b604051601f8201601f19908116603f01168101908382118183101715610373576103736102a4565b8160405282815288602084870101111561038c57600080fd5b61039d8360208301602088016102ba565b80955050505050509250929050565b600082516103be8184602087016102ba565b9190910192915050565b60208152600082518060208401526103e78160408501602087016102ba565b601f01601f19169190910160400192915050565b60b1806104096000396000f3fe608060405236601057600e6013565b005b600e5b601f601b6021565b6058565b565b600060537f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e8080156076573d6000f35b3d6000fdfea26469706673582212205d6a8985a07624eb030e02ca1492718ebe851f42cfdc725208d8f797b3b0befa64736f6c63430008130033416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564",
}

// ERC1967ProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1967ProxyMetaData.ABI instead.
var ERC1967ProxyABI = ERC1967ProxyMetaData.ABI

// ERC1967ProxyBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC1967ProxyBinRuntime = `608060405236601057600e6013565b005b600e5b601f601b6021565b6058565b565b600060537f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b3660008037600080366000845af43d6000803e8080156076573d6000f35b3d6000fdfea26469706673582212205d6a8985a07624eb030e02ca1492718ebe851f42cfdc725208d8f797b3b0befa64736f6c63430008130033`

// ERC1967ProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ERC1967ProxyMetaData.Bin instead.
var ERC1967ProxyBin = ERC1967ProxyMetaData.Bin

// DeployERC1967Proxy deploys a new Klaytn contract, binding an instance of ERC1967Proxy to it.
func DeployERC1967Proxy(auth *bind.TransactOpts, backend bind.ContractBackend, _logic common.Address, _data []byte) (common.Address, *types.Transaction, *ERC1967Proxy, error) {
	parsed, err := ERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ERC1967ProxyBin), backend, _logic, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ERC1967Proxy{ERC1967ProxyCaller: ERC1967ProxyCaller{contract: contract}, ERC1967ProxyTransactor: ERC1967ProxyTransactor{contract: contract}, ERC1967ProxyFilterer: ERC1967ProxyFilterer{contract: contract}}, nil
}

// ERC1967Proxy is an auto generated Go binding around a Klaytn contract.
type ERC1967Proxy struct {
	ERC1967ProxyCaller     // Read-only binding to the contract
	ERC1967ProxyTransactor // Write-only binding to the contract
	ERC1967ProxyFilterer   // Log filterer for contract events
}

// ERC1967ProxyCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC1967ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxyTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC1967ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxyFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC1967ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967ProxySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC1967ProxySession struct {
	Contract     *ERC1967Proxy     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC1967ProxyCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC1967ProxyCallerSession struct {
	Contract *ERC1967ProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ERC1967ProxyTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC1967ProxyTransactorSession struct {
	Contract     *ERC1967ProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ERC1967ProxyRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC1967ProxyRaw struct {
	Contract *ERC1967Proxy // Generic contract binding to access the raw methods on
}

// ERC1967ProxyCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC1967ProxyCallerRaw struct {
	Contract *ERC1967ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1967ProxyTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC1967ProxyTransactorRaw struct {
	Contract *ERC1967ProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1967Proxy creates a new instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967Proxy(address common.Address, backend bind.ContractBackend) (*ERC1967Proxy, error) {
	contract, err := bindERC1967Proxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1967Proxy{ERC1967ProxyCaller: ERC1967ProxyCaller{contract: contract}, ERC1967ProxyTransactor: ERC1967ProxyTransactor{contract: contract}, ERC1967ProxyFilterer: ERC1967ProxyFilterer{contract: contract}}, nil
}

// NewERC1967ProxyCaller creates a new read-only instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyCaller(address common.Address, caller bind.ContractCaller) (*ERC1967ProxyCaller, error) {
	contract, err := bindERC1967Proxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyCaller{contract: contract}, nil
}

// NewERC1967ProxyTransactor creates a new write-only instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1967ProxyTransactor, error) {
	contract, err := bindERC1967Proxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyTransactor{contract: contract}, nil
}

// NewERC1967ProxyFilterer creates a new log filterer instance of ERC1967Proxy, bound to a specific deployed contract.
func NewERC1967ProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1967ProxyFilterer, error) {
	contract, err := bindERC1967Proxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyFilterer{contract: contract}, nil
}

// bindERC1967Proxy binds a generic wrapper to an already deployed contract.
func bindERC1967Proxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Proxy *ERC1967ProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Proxy.Contract.ERC1967ProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Proxy *ERC1967ProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.ERC1967ProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Proxy *ERC1967ProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.ERC1967ProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Proxy *ERC1967ProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Proxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Proxy *ERC1967ProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Proxy *ERC1967ProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Fallback(&_ERC1967Proxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Fallback(&_ERC1967Proxy.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Proxy.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC1967Proxy *ERC1967ProxySession) Receive() (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Receive(&_ERC1967Proxy.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_ERC1967Proxy *ERC1967ProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _ERC1967Proxy.Contract.Receive(&_ERC1967Proxy.TransactOpts)
}

// ERC1967ProxyAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ERC1967Proxy contract.
type ERC1967ProxyAdminChangedIterator struct {
	Event *ERC1967ProxyAdminChanged // Event containing the contract specifics and raw log

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
func (it *ERC1967ProxyAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967ProxyAdminChanged)
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
		it.Event = new(ERC1967ProxyAdminChanged)
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
func (it *ERC1967ProxyAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967ProxyAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967ProxyAdminChanged represents a AdminChanged event raised by the ERC1967Proxy contract.
type ERC1967ProxyAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Proxy *ERC1967ProxyFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ERC1967ProxyAdminChangedIterator, error) {

	logs, sub, err := _ERC1967Proxy.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyAdminChangedIterator{contract: _ERC1967Proxy.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Proxy *ERC1967ProxyFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ERC1967ProxyAdminChanged) (event.Subscription, error) {

	logs, sub, err := _ERC1967Proxy.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967ProxyAdminChanged)
				if err := _ERC1967Proxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_ERC1967Proxy *ERC1967ProxyFilterer) ParseAdminChanged(log types.Log) (*ERC1967ProxyAdminChanged, error) {
	event := new(ERC1967ProxyAdminChanged)
	if err := _ERC1967Proxy.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967ProxyBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ERC1967Proxy contract.
type ERC1967ProxyBeaconUpgradedIterator struct {
	Event *ERC1967ProxyBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC1967ProxyBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967ProxyBeaconUpgraded)
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
		it.Event = new(ERC1967ProxyBeaconUpgraded)
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
func (it *ERC1967ProxyBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967ProxyBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967ProxyBeaconUpgraded represents a BeaconUpgraded event raised by the ERC1967Proxy contract.
type ERC1967ProxyBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Proxy *ERC1967ProxyFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ERC1967ProxyBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyBeaconUpgradedIterator{contract: _ERC1967Proxy.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Proxy *ERC1967ProxyFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967ProxyBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967ProxyBeaconUpgraded)
				if err := _ERC1967Proxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_ERC1967Proxy *ERC1967ProxyFilterer) ParseBeaconUpgraded(log types.Log) (*ERC1967ProxyBeaconUpgraded, error) {
	event := new(ERC1967ProxyBeaconUpgraded)
	if err := _ERC1967Proxy.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967ProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC1967Proxy contract.
type ERC1967ProxyUpgradedIterator struct {
	Event *ERC1967ProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC1967ProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967ProxyUpgraded)
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
		it.Event = new(ERC1967ProxyUpgraded)
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
func (it *ERC1967ProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967ProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967ProxyUpgraded represents a Upgraded event raised by the ERC1967Proxy contract.
type ERC1967ProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Proxy *ERC1967ProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC1967ProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967ProxyUpgradedIterator{contract: _ERC1967Proxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Proxy *ERC1967ProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967ProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Proxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967ProxyUpgraded)
				if err := _ERC1967Proxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ERC1967Proxy *ERC1967ProxyFilterer) ParseUpgraded(log types.Log) (*ERC1967ProxyUpgraded, error) {
	event := new(ERC1967ProxyUpgraded)
	if err := _ERC1967Proxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeMetaData contains all meta data concerning the ERC1967Upgrade contract.
var ERC1967UpgradeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"}]",
}

// ERC1967UpgradeABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC1967UpgradeMetaData.ABI instead.
var ERC1967UpgradeABI = ERC1967UpgradeMetaData.ABI

// ERC1967UpgradeBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ERC1967UpgradeBinRuntime = ``

// ERC1967Upgrade is an auto generated Go binding around a Klaytn contract.
type ERC1967Upgrade struct {
	ERC1967UpgradeCaller     // Read-only binding to the contract
	ERC1967UpgradeTransactor // Write-only binding to the contract
	ERC1967UpgradeFilterer   // Log filterer for contract events
}

// ERC1967UpgradeCaller is an auto generated read-only Go binding around a Klaytn contract.
type ERC1967UpgradeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ERC1967UpgradeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ERC1967UpgradeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC1967UpgradeSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ERC1967UpgradeSession struct {
	Contract     *ERC1967Upgrade   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC1967UpgradeCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ERC1967UpgradeCallerSession struct {
	Contract *ERC1967UpgradeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ERC1967UpgradeTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ERC1967UpgradeTransactorSession struct {
	Contract     *ERC1967UpgradeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ERC1967UpgradeRaw is an auto generated low-level Go binding around a Klaytn contract.
type ERC1967UpgradeRaw struct {
	Contract *ERC1967Upgrade // Generic contract binding to access the raw methods on
}

// ERC1967UpgradeCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ERC1967UpgradeCallerRaw struct {
	Contract *ERC1967UpgradeCaller // Generic read-only contract binding to access the raw methods on
}

// ERC1967UpgradeTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type ERC1967UpgradeTransactorRaw struct {
	Contract *ERC1967UpgradeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC1967Upgrade creates a new instance of ERC1967Upgrade, bound to a specific deployed contract.
func NewERC1967Upgrade(address common.Address, backend bind.ContractBackend) (*ERC1967Upgrade, error) {
	contract, err := bindERC1967Upgrade(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC1967Upgrade{ERC1967UpgradeCaller: ERC1967UpgradeCaller{contract: contract}, ERC1967UpgradeTransactor: ERC1967UpgradeTransactor{contract: contract}, ERC1967UpgradeFilterer: ERC1967UpgradeFilterer{contract: contract}}, nil
}

// NewERC1967UpgradeCaller creates a new read-only instance of ERC1967Upgrade, bound to a specific deployed contract.
func NewERC1967UpgradeCaller(address common.Address, caller bind.ContractCaller) (*ERC1967UpgradeCaller, error) {
	contract, err := bindERC1967Upgrade(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeCaller{contract: contract}, nil
}

// NewERC1967UpgradeTransactor creates a new write-only instance of ERC1967Upgrade, bound to a specific deployed contract.
func NewERC1967UpgradeTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC1967UpgradeTransactor, error) {
	contract, err := bindERC1967Upgrade(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeTransactor{contract: contract}, nil
}

// NewERC1967UpgradeFilterer creates a new log filterer instance of ERC1967Upgrade, bound to a specific deployed contract.
func NewERC1967UpgradeFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC1967UpgradeFilterer, error) {
	contract, err := bindERC1967Upgrade(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeFilterer{contract: contract}, nil
}

// bindERC1967Upgrade binds a generic wrapper to an already deployed contract.
func bindERC1967Upgrade(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ERC1967UpgradeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Upgrade *ERC1967UpgradeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Upgrade.Contract.ERC1967UpgradeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Upgrade *ERC1967UpgradeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Upgrade.Contract.ERC1967UpgradeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Upgrade *ERC1967UpgradeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Upgrade.Contract.ERC1967UpgradeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC1967Upgrade *ERC1967UpgradeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC1967Upgrade.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC1967Upgrade *ERC1967UpgradeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC1967Upgrade.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC1967Upgrade *ERC1967UpgradeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC1967Upgrade.Contract.contract.Transact(opts, method, params...)
}

// ERC1967UpgradeAdminChangedIterator is returned from FilterAdminChanged and is used to iterate over the raw logs and unpacked data for AdminChanged events raised by the ERC1967Upgrade contract.
type ERC1967UpgradeAdminChangedIterator struct {
	Event *ERC1967UpgradeAdminChanged // Event containing the contract specifics and raw log

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
func (it *ERC1967UpgradeAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeAdminChanged)
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
		it.Event = new(ERC1967UpgradeAdminChanged)
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
func (it *ERC1967UpgradeAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeAdminChanged represents a AdminChanged event raised by the ERC1967Upgrade contract.
type ERC1967UpgradeAdminChanged struct {
	PreviousAdmin common.Address
	NewAdmin      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAdminChanged is a free log retrieval operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) FilterAdminChanged(opts *bind.FilterOpts) (*ERC1967UpgradeAdminChangedIterator, error) {

	logs, sub, err := _ERC1967Upgrade.contract.FilterLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeAdminChangedIterator{contract: _ERC1967Upgrade.contract, event: "AdminChanged", logs: logs, sub: sub}, nil
}

// WatchAdminChanged is a free log subscription operation binding the contract event 0x7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f.
//
// Solidity: event AdminChanged(address previousAdmin, address newAdmin)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) WatchAdminChanged(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeAdminChanged) (event.Subscription, error) {

	logs, sub, err := _ERC1967Upgrade.contract.WatchLogs(opts, "AdminChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeAdminChanged)
				if err := _ERC1967Upgrade.contract.UnpackLog(event, "AdminChanged", log); err != nil {
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
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) ParseAdminChanged(log types.Log) (*ERC1967UpgradeAdminChanged, error) {
	event := new(ERC1967UpgradeAdminChanged)
	if err := _ERC1967Upgrade.contract.UnpackLog(event, "AdminChanged", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeBeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the ERC1967Upgrade contract.
type ERC1967UpgradeBeaconUpgradedIterator struct {
	Event *ERC1967UpgradeBeaconUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC1967UpgradeBeaconUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeBeaconUpgraded)
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
		it.Event = new(ERC1967UpgradeBeaconUpgraded)
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
func (it *ERC1967UpgradeBeaconUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeBeaconUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeBeaconUpgraded represents a BeaconUpgraded event raised by the ERC1967Upgrade contract.
type ERC1967UpgradeBeaconUpgraded struct {
	Beacon common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBeaconUpgraded is a free log retrieval operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) FilterBeaconUpgraded(opts *bind.FilterOpts, beacon []common.Address) (*ERC1967UpgradeBeaconUpgradedIterator, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967Upgrade.contract.FilterLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeBeaconUpgradedIterator{contract: _ERC1967Upgrade.contract, event: "BeaconUpgraded", logs: logs, sub: sub}, nil
}

// WatchBeaconUpgraded is a free log subscription operation binding the contract event 0x1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e.
//
// Solidity: event BeaconUpgraded(address indexed beacon)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) WatchBeaconUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeBeaconUpgraded, beacon []common.Address) (event.Subscription, error) {

	var beaconRule []interface{}
	for _, beaconItem := range beacon {
		beaconRule = append(beaconRule, beaconItem)
	}

	logs, sub, err := _ERC1967Upgrade.contract.WatchLogs(opts, "BeaconUpgraded", beaconRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeBeaconUpgraded)
				if err := _ERC1967Upgrade.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
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
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) ParseBeaconUpgraded(log types.Log) (*ERC1967UpgradeBeaconUpgraded, error) {
	event := new(ERC1967UpgradeBeaconUpgraded)
	if err := _ERC1967Upgrade.contract.UnpackLog(event, "BeaconUpgraded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ERC1967UpgradeUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ERC1967Upgrade contract.
type ERC1967UpgradeUpgradedIterator struct {
	Event *ERC1967UpgradeUpgraded // Event containing the contract specifics and raw log

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
func (it *ERC1967UpgradeUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC1967UpgradeUpgraded)
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
		it.Event = new(ERC1967UpgradeUpgraded)
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
func (it *ERC1967UpgradeUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC1967UpgradeUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC1967UpgradeUpgraded represents a Upgraded event raised by the ERC1967Upgrade contract.
type ERC1967UpgradeUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ERC1967UpgradeUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Upgrade.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ERC1967UpgradeUpgradedIterator{contract: _ERC1967Upgrade.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ERC1967UpgradeUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ERC1967Upgrade.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC1967UpgradeUpgraded)
				if err := _ERC1967Upgrade.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ERC1967Upgrade *ERC1967UpgradeFilterer) ParseUpgraded(log types.Log) (*ERC1967UpgradeUpgraded, error) {
	event := new(ERC1967UpgradeUpgraded)
	if err := _ERC1967Upgrade.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	return event, nil
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

// IBeaconFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IBeaconMetaData.Sigs instead.
var IBeaconFuncSigs = IBeaconMetaData.Sigs

// IBeacon is an auto generated Go binding around a Klaytn contract.
type IBeacon struct {
	IBeaconCaller     // Read-only binding to the contract
	IBeaconTransactor // Write-only binding to the contract
	IBeaconFilterer   // Log filterer for contract events
}

// IBeaconCaller is an auto generated read-only Go binding around a Klaytn contract.
type IBeaconCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IBeaconTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IBeaconFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IBeaconSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IBeaconSession struct {
	Contract     *IBeacon          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IBeaconCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IBeaconCallerSession struct {
	Contract *IBeaconCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// IBeaconTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IBeaconTransactorSession struct {
	Contract     *IBeaconTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IBeaconRaw is an auto generated low-level Go binding around a Klaytn contract.
type IBeaconRaw struct {
	Contract *IBeacon // Generic contract binding to access the raw methods on
}

// IBeaconCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IBeaconCallerRaw struct {
	Contract *IBeaconCaller // Generic read-only contract binding to access the raw methods on
}

// IBeaconTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
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

// IERC1822ProxiableMetaData contains all meta data concerning the IERC1822Proxiable contract.
var IERC1822ProxiableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"52d1902d": "proxiableUUID()",
	},
}

// IERC1822ProxiableABI is the input ABI used to generate the binding from.
// Deprecated: Use IERC1822ProxiableMetaData.ABI instead.
var IERC1822ProxiableABI = IERC1822ProxiableMetaData.ABI

// IERC1822ProxiableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IERC1822ProxiableBinRuntime = ``

// IERC1822ProxiableFuncSigs maps the 4-byte function signature to its string representation.
// Deprecated: Use IERC1822ProxiableMetaData.Sigs instead.
var IERC1822ProxiableFuncSigs = IERC1822ProxiableMetaData.Sigs

// IERC1822Proxiable is an auto generated Go binding around a Klaytn contract.
type IERC1822Proxiable struct {
	IERC1822ProxiableCaller     // Read-only binding to the contract
	IERC1822ProxiableTransactor // Write-only binding to the contract
	IERC1822ProxiableFilterer   // Log filterer for contract events
}

// IERC1822ProxiableCaller is an auto generated read-only Go binding around a Klaytn contract.
type IERC1822ProxiableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableTransactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC1822ProxiableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC1822ProxiableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1822ProxiableSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC1822ProxiableSession struct {
	Contract     *IERC1822Proxiable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IERC1822ProxiableCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC1822ProxiableCallerSession struct {
	Contract *IERC1822ProxiableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IERC1822ProxiableTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC1822ProxiableTransactorSession struct {
	Contract     *IERC1822ProxiableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IERC1822ProxiableRaw is an auto generated low-level Go binding around a Klaytn contract.
type IERC1822ProxiableRaw struct {
	Contract *IERC1822Proxiable // Generic contract binding to access the raw methods on
}

// IERC1822ProxiableCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC1822ProxiableCallerRaw struct {
	Contract *IERC1822ProxiableCaller // Generic read-only contract binding to access the raw methods on
}

// IERC1822ProxiableTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
type IERC1822ProxiableTransactorRaw struct {
	Contract *IERC1822ProxiableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIERC1822Proxiable creates a new instance of IERC1822Proxiable, bound to a specific deployed contract.
func NewIERC1822Proxiable(address common.Address, backend bind.ContractBackend) (*IERC1822Proxiable, error) {
	contract, err := bindIERC1822Proxiable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IERC1822Proxiable{IERC1822ProxiableCaller: IERC1822ProxiableCaller{contract: contract}, IERC1822ProxiableTransactor: IERC1822ProxiableTransactor{contract: contract}, IERC1822ProxiableFilterer: IERC1822ProxiableFilterer{contract: contract}}, nil
}

// NewIERC1822ProxiableCaller creates a new read-only instance of IERC1822Proxiable, bound to a specific deployed contract.
func NewIERC1822ProxiableCaller(address common.Address, caller bind.ContractCaller) (*IERC1822ProxiableCaller, error) {
	contract, err := bindIERC1822Proxiable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableCaller{contract: contract}, nil
}

// NewIERC1822ProxiableTransactor creates a new write-only instance of IERC1822Proxiable, bound to a specific deployed contract.
func NewIERC1822ProxiableTransactor(address common.Address, transactor bind.ContractTransactor) (*IERC1822ProxiableTransactor, error) {
	contract, err := bindIERC1822Proxiable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableTransactor{contract: contract}, nil
}

// NewIERC1822ProxiableFilterer creates a new log filterer instance of IERC1822Proxiable, bound to a specific deployed contract.
func NewIERC1822ProxiableFilterer(address common.Address, filterer bind.ContractFilterer) (*IERC1822ProxiableFilterer, error) {
	contract, err := bindIERC1822Proxiable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IERC1822ProxiableFilterer{contract: contract}, nil
}

// bindIERC1822Proxiable binds a generic wrapper to an already deployed contract.
func bindIERC1822Proxiable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IERC1822ProxiableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1822Proxiable *IERC1822ProxiableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1822Proxiable.Contract.IERC1822ProxiableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1822Proxiable *IERC1822ProxiableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1822Proxiable.Contract.IERC1822ProxiableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1822Proxiable *IERC1822ProxiableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1822Proxiable.Contract.IERC1822ProxiableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IERC1822Proxiable *IERC1822ProxiableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IERC1822Proxiable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IERC1822Proxiable *IERC1822ProxiableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IERC1822Proxiable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IERC1822Proxiable *IERC1822ProxiableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IERC1822Proxiable.Contract.contract.Transact(opts, method, params...)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822Proxiable *IERC1822ProxiableCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _IERC1822Proxiable.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822Proxiable *IERC1822ProxiableSession) ProxiableUUID() ([32]byte, error) {
	return _IERC1822Proxiable.Contract.ProxiableUUID(&_IERC1822Proxiable.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_IERC1822Proxiable *IERC1822ProxiableCallerSession) ProxiableUUID() ([32]byte, error) {
	return _IERC1822Proxiable.Contract.ProxiableUUID(&_IERC1822Proxiable.CallOpts)
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

// IERC1967 is an auto generated Go binding around a Klaytn contract.
type IERC1967 struct {
	IERC1967Caller     // Read-only binding to the contract
	IERC1967Transactor // Write-only binding to the contract
	IERC1967Filterer   // Log filterer for contract events
}

// IERC1967Caller is an auto generated read-only Go binding around a Klaytn contract.
type IERC1967Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Transactor is an auto generated write-only Go binding around a Klaytn contract.
type IERC1967Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Filterer is an auto generated log filtering Go binding around a Klaytn contract events.
type IERC1967Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IERC1967Session is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type IERC1967Session struct {
	Contract     *IERC1967         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IERC1967CallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type IERC1967CallerSession struct {
	Contract *IERC1967Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// IERC1967TransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type IERC1967TransactorSession struct {
	Contract     *IERC1967Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IERC1967Raw is an auto generated low-level Go binding around a Klaytn contract.
type IERC1967Raw struct {
	Contract *IERC1967 // Generic contract binding to access the raw methods on
}

// IERC1967CallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type IERC1967CallerRaw struct {
	Contract *IERC1967Caller // Generic read-only contract binding to access the raw methods on
}

// IERC1967TransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
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

	logs chan types.Log      // Log channel receiving the found contract events
	sub  klaytn.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
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
	return event, nil
}

// IERC1967BeaconUpgradedIterator is returned from FilterBeaconUpgraded and is used to iterate over the raw logs and unpacked data for BeaconUpgraded events raised by the IERC1967 contract.
type IERC1967BeaconUpgradedIterator struct {
	Event *IERC1967BeaconUpgraded // Event containing the contract specifics and raw log

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
	return event, nil
}

// IERC1967UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the IERC1967 contract.
type IERC1967UpgradedIterator struct {
	Event *IERC1967Upgraded // Event containing the contract specifics and raw log

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
	return event, nil
}

// ProxyMetaData contains all meta data concerning the Proxy contract.
var ProxyMetaData = &bind.MetaData{
	ABI: "[{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// ProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use ProxyMetaData.ABI instead.
var ProxyABI = ProxyMetaData.ABI

// ProxyBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ProxyBinRuntime = ``

// Proxy is an auto generated Go binding around a Klaytn contract.
type Proxy struct {
	ProxyCaller     // Read-only binding to the contract
	ProxyTransactor // Write-only binding to the contract
	ProxyFilterer   // Log filterer for contract events
}

// ProxyCaller is an auto generated read-only Go binding around a Klaytn contract.
type ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyTransactor is an auto generated write-only Go binding around a Klaytn contract.
type ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxyFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProxySession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type ProxySession struct {
	Contract     *Proxy            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type ProxyCallerSession struct {
	Contract *ProxyCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ProxyTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type ProxyTransactorSession struct {
	Contract     *ProxyTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProxyRaw is an auto generated low-level Go binding around a Klaytn contract.
type ProxyRaw struct {
	Contract *Proxy // Generic contract binding to access the raw methods on
}

// ProxyCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type ProxyCallerRaw struct {
	Contract *ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// ProxyTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
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

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Proxy *ProxyTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proxy.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Proxy *ProxySession) Receive() (*types.Transaction, error) {
	return _Proxy.Contract.Receive(&_Proxy.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Proxy *ProxyTransactorSession) Receive() (*types.Transaction, error) {
	return _Proxy.Contract.Receive(&_Proxy.TransactOpts)
}

// StorageSlotMetaData contains all meta data concerning the StorageSlot contract.
var StorageSlotMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60566037600b82828239805160001a607314602a57634e487b7160e01b600052600060045260246000fd5b30600052607381538281f3fe73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212200d7c79cc0d2dc5a13cd1f9ec7db59a64db211d68e514f0971fff90f2f3fff0a564736f6c63430008130033",
}

// StorageSlotABI is the input ABI used to generate the binding from.
// Deprecated: Use StorageSlotMetaData.ABI instead.
var StorageSlotABI = StorageSlotMetaData.ABI

// StorageSlotBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const StorageSlotBinRuntime = `73000000000000000000000000000000000000000030146080604052600080fdfea26469706673582212200d7c79cc0d2dc5a13cd1f9ec7db59a64db211d68e514f0971fff90f2f3fff0a564736f6c63430008130033`

// StorageSlotBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StorageSlotMetaData.Bin instead.
var StorageSlotBin = StorageSlotMetaData.Bin

// DeployStorageSlot deploys a new Klaytn contract, binding an instance of StorageSlot to it.
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

// StorageSlot is an auto generated Go binding around a Klaytn contract.
type StorageSlot struct {
	StorageSlotCaller     // Read-only binding to the contract
	StorageSlotTransactor // Write-only binding to the contract
	StorageSlotFilterer   // Log filterer for contract events
}

// StorageSlotCaller is an auto generated read-only Go binding around a Klaytn contract.
type StorageSlotCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotTransactor is an auto generated write-only Go binding around a Klaytn contract.
type StorageSlotTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotFilterer is an auto generated log filtering Go binding around a Klaytn contract events.
type StorageSlotFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StorageSlotSession is an auto generated Go binding around a Klaytn contract,
// with pre-set call and transact options.
type StorageSlotSession struct {
	Contract     *StorageSlot      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StorageSlotCallerSession is an auto generated read-only Go binding around a Klaytn contract,
// with pre-set call options.
type StorageSlotCallerSession struct {
	Contract *StorageSlotCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StorageSlotTransactorSession is an auto generated write-only Go binding around a Klaytn contract,
// with pre-set transact options.
type StorageSlotTransactorSession struct {
	Contract     *StorageSlotTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StorageSlotRaw is an auto generated low-level Go binding around a Klaytn contract.
type StorageSlotRaw struct {
	Contract *StorageSlot // Generic contract binding to access the raw methods on
}

// StorageSlotCallerRaw is an auto generated low-level read-only Go binding around a Klaytn contract.
type StorageSlotCallerRaw struct {
	Contract *StorageSlotCaller // Generic read-only contract binding to access the raw methods on
}

// StorageSlotTransactorRaw is an auto generated low-level write-only Go binding around a Klaytn contract.
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
