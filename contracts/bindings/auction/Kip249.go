// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package auction

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

// IAuctionEntryPointAuctionTx is an auto generated low-level Go binding around an user-defined struct.
type IAuctionEntryPointAuctionTx struct {
	TargetTxHash  [32]byte
	BlockNumber   *big.Int
	Sender        common.Address
	To            common.Address
	Nonce         *big.Int
	Bid           *big.Int
	CallGasLimit  *big.Int
	Data          []byte
	SearcherSig   []byte
	AuctioneerSig []byte
}

// IAuctionEntryPointMetaData contains all meta data concerning the IAuctionEntryPoint contract.
var IAuctionEntryPointMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Call\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"CallFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAuctioneer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAuctioneer\",\"type\":\"address\"}],\"name\":\"ChangeAuctioneer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldDepositVault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newDepositVault\",\"type\":\"address\"}],\"name\":\"ChangeDepositVault\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPerByteIntrinsic\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPerByteEip7623\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasContractExecution\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasBufferEstimate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasBufferUnmeasured\",\"type\":\"uint256\"}],\"name\":\"ChangeGasParameters\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"UseNonce\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"auctioneer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"targetTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"searcherSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"auctioneerSig\",\"type\":\"bytes\"}],\"internalType\":\"structIAuctionEntryPoint.AuctionTx\",\"name\":\"auctionTx\",\"type\":\"tuple\"}],\"name\":\"call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_auctioneer\",\"type\":\"address\"}],\"name\":\"changeAuctioneer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_depositVault\",\"type\":\"address\"}],\"name\":\"changeDepositVault\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_gasPerByteIntrinsic\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPerByteEip7623\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasContractExecution\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasBufferEstimate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasBufferUnmeasured\",\"type\":\"uint256\"}],\"name\":\"changeGasParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositVault\",\"outputs\":[{\"internalType\":\"contractIAuctionDepositVault\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasBufferEstimate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasBufferUnmeasured\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasContractExecution\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasPerByteEip7623\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasPerByteIntrinsic\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"targetTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"searcherSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"auctioneerSig\",\"type\":\"bytes\"}],\"internalType\":\"structIAuctionEntryPoint.AuctionTx\",\"name\":\"auctionTx\",\"type\":\"tuple\"}],\"name\":\"getAuctionTxHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"searchers\",\"type\":\"address[]\"}],\"name\":\"getNoncesAndDeposits\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nonces_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"deposits_\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IAuctionEntryPointABI is the input ABI used to generate the binding from.
// Deprecated: Use IAuctionEntryPointMetaData.ABI instead.
var IAuctionEntryPointABI = IAuctionEntryPointMetaData.ABI

// IAuctionEntryPointBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IAuctionEntryPointBinRuntime = ``

// IAuctionEntryPoint is an auto generated Go binding around a Kaia contract.
type IAuctionEntryPoint struct {
	IAuctionEntryPointCaller     // Read-only binding to the contract
	IAuctionEntryPointTransactor // Write-only binding to the contract
	IAuctionEntryPointFilterer   // Log filterer for contract events
}

// IAuctionEntryPointCaller is an auto generated read-only Go binding around a Kaia contract.
type IAuctionEntryPointCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionEntryPointTransactor is an auto generated write-only Go binding around a Kaia contract.
type IAuctionEntryPointTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionEntryPointFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IAuctionEntryPointFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionEntryPointSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IAuctionEntryPointSession struct {
	Contract     *IAuctionEntryPoint // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// IAuctionEntryPointCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IAuctionEntryPointCallerSession struct {
	Contract *IAuctionEntryPointCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// IAuctionEntryPointTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IAuctionEntryPointTransactorSession struct {
	Contract     *IAuctionEntryPointTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// IAuctionEntryPointRaw is an auto generated low-level Go binding around a Kaia contract.
type IAuctionEntryPointRaw struct {
	Contract *IAuctionEntryPoint // Generic contract binding to access the raw methods on
}

// IAuctionEntryPointCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IAuctionEntryPointCallerRaw struct {
	Contract *IAuctionEntryPointCaller // Generic read-only contract binding to access the raw methods on
}

// IAuctionEntryPointTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IAuctionEntryPointTransactorRaw struct {
	Contract *IAuctionEntryPointTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAuctionEntryPoint creates a new instance of IAuctionEntryPoint, bound to a specific deployed contract.
func NewIAuctionEntryPoint(address common.Address, backend bind.ContractBackend) (*IAuctionEntryPoint, error) {
	contract, err := bindIAuctionEntryPoint(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPoint{IAuctionEntryPointCaller: IAuctionEntryPointCaller{contract: contract}, IAuctionEntryPointTransactor: IAuctionEntryPointTransactor{contract: contract}, IAuctionEntryPointFilterer: IAuctionEntryPointFilterer{contract: contract}}, nil
}

// NewIAuctionEntryPointCaller creates a new read-only instance of IAuctionEntryPoint, bound to a specific deployed contract.
func NewIAuctionEntryPointCaller(address common.Address, caller bind.ContractCaller) (*IAuctionEntryPointCaller, error) {
	contract, err := bindIAuctionEntryPoint(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointCaller{contract: contract}, nil
}

// NewIAuctionEntryPointTransactor creates a new write-only instance of IAuctionEntryPoint, bound to a specific deployed contract.
func NewIAuctionEntryPointTransactor(address common.Address, transactor bind.ContractTransactor) (*IAuctionEntryPointTransactor, error) {
	contract, err := bindIAuctionEntryPoint(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointTransactor{contract: contract}, nil
}

// NewIAuctionEntryPointFilterer creates a new log filterer instance of IAuctionEntryPoint, bound to a specific deployed contract.
func NewIAuctionEntryPointFilterer(address common.Address, filterer bind.ContractFilterer) (*IAuctionEntryPointFilterer, error) {
	contract, err := bindIAuctionEntryPoint(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointFilterer{contract: contract}, nil
}

// bindIAuctionEntryPoint binds a generic wrapper to an already deployed contract.
func bindIAuctionEntryPoint(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAuctionEntryPointMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAuctionEntryPoint *IAuctionEntryPointRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAuctionEntryPoint.Contract.IAuctionEntryPointCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAuctionEntryPoint *IAuctionEntryPointRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.IAuctionEntryPointTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAuctionEntryPoint *IAuctionEntryPointRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.IAuctionEntryPointTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAuctionEntryPoint *IAuctionEntryPointCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAuctionEntryPoint.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAuctionEntryPoint *IAuctionEntryPointTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAuctionEntryPoint *IAuctionEntryPointTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.contract.Transact(opts, method, params...)
}

// Auctioneer is a free data retrieval call binding the contract method 0x5ec2c7bf.
//
// Solidity: function auctioneer() view returns(address)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) Auctioneer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "auctioneer")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Auctioneer is a free data retrieval call binding the contract method 0x5ec2c7bf.
//
// Solidity: function auctioneer() view returns(address)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) Auctioneer() (common.Address, error) {
	return _IAuctionEntryPoint.Contract.Auctioneer(&_IAuctionEntryPoint.CallOpts)
}

// Auctioneer is a free data retrieval call binding the contract method 0x5ec2c7bf.
//
// Solidity: function auctioneer() view returns(address)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) Auctioneer() (common.Address, error) {
	return _IAuctionEntryPoint.Contract.Auctioneer(&_IAuctionEntryPoint.CallOpts)
}

// DepositVault is a free data retrieval call binding the contract method 0xd7cd3949.
//
// Solidity: function depositVault() view returns(address)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) DepositVault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "depositVault")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// DepositVault is a free data retrieval call binding the contract method 0xd7cd3949.
//
// Solidity: function depositVault() view returns(address)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) DepositVault() (common.Address, error) {
	return _IAuctionEntryPoint.Contract.DepositVault(&_IAuctionEntryPoint.CallOpts)
}

// DepositVault is a free data retrieval call binding the contract method 0xd7cd3949.
//
// Solidity: function depositVault() view returns(address)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) DepositVault() (common.Address, error) {
	return _IAuctionEntryPoint.Contract.DepositVault(&_IAuctionEntryPoint.CallOpts)
}

// GasBufferEstimate is a free data retrieval call binding the contract method 0xa5b2ab40.
//
// Solidity: function gasBufferEstimate() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) GasBufferEstimate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "gasBufferEstimate")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GasBufferEstimate is a free data retrieval call binding the contract method 0xa5b2ab40.
//
// Solidity: function gasBufferEstimate() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) GasBufferEstimate() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasBufferEstimate(&_IAuctionEntryPoint.CallOpts)
}

// GasBufferEstimate is a free data retrieval call binding the contract method 0xa5b2ab40.
//
// Solidity: function gasBufferEstimate() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) GasBufferEstimate() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasBufferEstimate(&_IAuctionEntryPoint.CallOpts)
}

// GasBufferUnmeasured is a free data retrieval call binding the contract method 0x145aa6c7.
//
// Solidity: function gasBufferUnmeasured() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) GasBufferUnmeasured(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "gasBufferUnmeasured")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GasBufferUnmeasured is a free data retrieval call binding the contract method 0x145aa6c7.
//
// Solidity: function gasBufferUnmeasured() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) GasBufferUnmeasured() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasBufferUnmeasured(&_IAuctionEntryPoint.CallOpts)
}

// GasBufferUnmeasured is a free data retrieval call binding the contract method 0x145aa6c7.
//
// Solidity: function gasBufferUnmeasured() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) GasBufferUnmeasured() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasBufferUnmeasured(&_IAuctionEntryPoint.CallOpts)
}

// GasContractExecution is a free data retrieval call binding the contract method 0xc62a1105.
//
// Solidity: function gasContractExecution() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) GasContractExecution(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "gasContractExecution")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GasContractExecution is a free data retrieval call binding the contract method 0xc62a1105.
//
// Solidity: function gasContractExecution() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) GasContractExecution() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasContractExecution(&_IAuctionEntryPoint.CallOpts)
}

// GasContractExecution is a free data retrieval call binding the contract method 0xc62a1105.
//
// Solidity: function gasContractExecution() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) GasContractExecution() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasContractExecution(&_IAuctionEntryPoint.CallOpts)
}

// GasPerByteEip7623 is a free data retrieval call binding the contract method 0x5e155d7a.
//
// Solidity: function gasPerByteEip7623() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) GasPerByteEip7623(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "gasPerByteEip7623")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GasPerByteEip7623 is a free data retrieval call binding the contract method 0x5e155d7a.
//
// Solidity: function gasPerByteEip7623() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) GasPerByteEip7623() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasPerByteEip7623(&_IAuctionEntryPoint.CallOpts)
}

// GasPerByteEip7623 is a free data retrieval call binding the contract method 0x5e155d7a.
//
// Solidity: function gasPerByteEip7623() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) GasPerByteEip7623() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasPerByteEip7623(&_IAuctionEntryPoint.CallOpts)
}

// GasPerByteIntrinsic is a free data retrieval call binding the contract method 0xea147667.
//
// Solidity: function gasPerByteIntrinsic() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) GasPerByteIntrinsic(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "gasPerByteIntrinsic")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GasPerByteIntrinsic is a free data retrieval call binding the contract method 0xea147667.
//
// Solidity: function gasPerByteIntrinsic() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) GasPerByteIntrinsic() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasPerByteIntrinsic(&_IAuctionEntryPoint.CallOpts)
}

// GasPerByteIntrinsic is a free data retrieval call binding the contract method 0xea147667.
//
// Solidity: function gasPerByteIntrinsic() view returns(uint256)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) GasPerByteIntrinsic() (*big.Int, error) {
	return _IAuctionEntryPoint.Contract.GasPerByteIntrinsic(&_IAuctionEntryPoint.CallOpts)
}

// GetAuctionTxHash is a free data retrieval call binding the contract method 0xa8aa9450.
//
// Solidity: function getAuctionTxHash((bytes32,uint256,address,address,uint256,uint256,uint256,bytes,bytes,bytes) auctionTx) view returns(bytes32)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) GetAuctionTxHash(opts *bind.CallOpts, auctionTx IAuctionEntryPointAuctionTx) ([32]byte, error) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "getAuctionTxHash", auctionTx)
	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err
}

// GetAuctionTxHash is a free data retrieval call binding the contract method 0xa8aa9450.
//
// Solidity: function getAuctionTxHash((bytes32,uint256,address,address,uint256,uint256,uint256,bytes,bytes,bytes) auctionTx) view returns(bytes32)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) GetAuctionTxHash(auctionTx IAuctionEntryPointAuctionTx) ([32]byte, error) {
	return _IAuctionEntryPoint.Contract.GetAuctionTxHash(&_IAuctionEntryPoint.CallOpts, auctionTx)
}

// GetAuctionTxHash is a free data retrieval call binding the contract method 0xa8aa9450.
//
// Solidity: function getAuctionTxHash((bytes32,uint256,address,address,uint256,uint256,uint256,bytes,bytes,bytes) auctionTx) view returns(bytes32)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) GetAuctionTxHash(auctionTx IAuctionEntryPointAuctionTx) ([32]byte, error) {
	return _IAuctionEntryPoint.Contract.GetAuctionTxHash(&_IAuctionEntryPoint.CallOpts, auctionTx)
}

// GetNoncesAndDeposits is a free data retrieval call binding the contract method 0x0339ed37.
//
// Solidity: function getNoncesAndDeposits(address[] searchers) view returns(uint256[] nonces_, uint256[] deposits_)
func (_IAuctionEntryPoint *IAuctionEntryPointCaller) GetNoncesAndDeposits(opts *bind.CallOpts, searchers []common.Address) (struct {
	Nonces   []*big.Int
	Deposits []*big.Int
}, error,
) {
	var out []interface{}
	err := _IAuctionEntryPoint.contract.Call(opts, &out, "getNoncesAndDeposits", searchers)

	outstruct := new(struct {
		Nonces   []*big.Int
		Deposits []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Nonces = *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	outstruct.Deposits = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err
}

// GetNoncesAndDeposits is a free data retrieval call binding the contract method 0x0339ed37.
//
// Solidity: function getNoncesAndDeposits(address[] searchers) view returns(uint256[] nonces_, uint256[] deposits_)
func (_IAuctionEntryPoint *IAuctionEntryPointSession) GetNoncesAndDeposits(searchers []common.Address) (struct {
	Nonces   []*big.Int
	Deposits []*big.Int
}, error,
) {
	return _IAuctionEntryPoint.Contract.GetNoncesAndDeposits(&_IAuctionEntryPoint.CallOpts, searchers)
}

// GetNoncesAndDeposits is a free data retrieval call binding the contract method 0x0339ed37.
//
// Solidity: function getNoncesAndDeposits(address[] searchers) view returns(uint256[] nonces_, uint256[] deposits_)
func (_IAuctionEntryPoint *IAuctionEntryPointCallerSession) GetNoncesAndDeposits(searchers []common.Address) (struct {
	Nonces   []*big.Int
	Deposits []*big.Int
}, error,
) {
	return _IAuctionEntryPoint.Contract.GetNoncesAndDeposits(&_IAuctionEntryPoint.CallOpts, searchers)
}

// Call is a paid mutator transaction binding the contract method 0xca157554.
//
// Solidity: function call((bytes32,uint256,address,address,uint256,uint256,uint256,bytes,bytes,bytes) auctionTx) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointTransactor) Call(opts *bind.TransactOpts, auctionTx IAuctionEntryPointAuctionTx) (*types.Transaction, error) {
	return _IAuctionEntryPoint.contract.Transact(opts, "call", auctionTx)
}

// Call is a paid mutator transaction binding the contract method 0xca157554.
//
// Solidity: function call((bytes32,uint256,address,address,uint256,uint256,uint256,bytes,bytes,bytes) auctionTx) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointSession) Call(auctionTx IAuctionEntryPointAuctionTx) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.Call(&_IAuctionEntryPoint.TransactOpts, auctionTx)
}

// Call is a paid mutator transaction binding the contract method 0xca157554.
//
// Solidity: function call((bytes32,uint256,address,address,uint256,uint256,uint256,bytes,bytes,bytes) auctionTx) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointTransactorSession) Call(auctionTx IAuctionEntryPointAuctionTx) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.Call(&_IAuctionEntryPoint.TransactOpts, auctionTx)
}

// ChangeAuctioneer is a paid mutator transaction binding the contract method 0x774f45ec.
//
// Solidity: function changeAuctioneer(address _auctioneer) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointTransactor) ChangeAuctioneer(opts *bind.TransactOpts, _auctioneer common.Address) (*types.Transaction, error) {
	return _IAuctionEntryPoint.contract.Transact(opts, "changeAuctioneer", _auctioneer)
}

// ChangeAuctioneer is a paid mutator transaction binding the contract method 0x774f45ec.
//
// Solidity: function changeAuctioneer(address _auctioneer) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointSession) ChangeAuctioneer(_auctioneer common.Address) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.ChangeAuctioneer(&_IAuctionEntryPoint.TransactOpts, _auctioneer)
}

// ChangeAuctioneer is a paid mutator transaction binding the contract method 0x774f45ec.
//
// Solidity: function changeAuctioneer(address _auctioneer) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointTransactorSession) ChangeAuctioneer(_auctioneer common.Address) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.ChangeAuctioneer(&_IAuctionEntryPoint.TransactOpts, _auctioneer)
}

// ChangeDepositVault is a paid mutator transaction binding the contract method 0x9d59928b.
//
// Solidity: function changeDepositVault(address _depositVault) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointTransactor) ChangeDepositVault(opts *bind.TransactOpts, _depositVault common.Address) (*types.Transaction, error) {
	return _IAuctionEntryPoint.contract.Transact(opts, "changeDepositVault", _depositVault)
}

// ChangeDepositVault is a paid mutator transaction binding the contract method 0x9d59928b.
//
// Solidity: function changeDepositVault(address _depositVault) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointSession) ChangeDepositVault(_depositVault common.Address) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.ChangeDepositVault(&_IAuctionEntryPoint.TransactOpts, _depositVault)
}

// ChangeDepositVault is a paid mutator transaction binding the contract method 0x9d59928b.
//
// Solidity: function changeDepositVault(address _depositVault) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointTransactorSession) ChangeDepositVault(_depositVault common.Address) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.ChangeDepositVault(&_IAuctionEntryPoint.TransactOpts, _depositVault)
}

// ChangeGasParameters is a paid mutator transaction binding the contract method 0x2a215610.
//
// Solidity: function changeGasParameters(uint256 _gasPerByteIntrinsic, uint256 _gasPerByteEip7623, uint256 _gasContractExecution, uint256 _gasBufferEstimate, uint256 _gasBufferUnmeasured) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointTransactor) ChangeGasParameters(opts *bind.TransactOpts, _gasPerByteIntrinsic *big.Int, _gasPerByteEip7623 *big.Int, _gasContractExecution *big.Int, _gasBufferEstimate *big.Int, _gasBufferUnmeasured *big.Int) (*types.Transaction, error) {
	return _IAuctionEntryPoint.contract.Transact(opts, "changeGasParameters", _gasPerByteIntrinsic, _gasPerByteEip7623, _gasContractExecution, _gasBufferEstimate, _gasBufferUnmeasured)
}

// ChangeGasParameters is a paid mutator transaction binding the contract method 0x2a215610.
//
// Solidity: function changeGasParameters(uint256 _gasPerByteIntrinsic, uint256 _gasPerByteEip7623, uint256 _gasContractExecution, uint256 _gasBufferEstimate, uint256 _gasBufferUnmeasured) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointSession) ChangeGasParameters(_gasPerByteIntrinsic *big.Int, _gasPerByteEip7623 *big.Int, _gasContractExecution *big.Int, _gasBufferEstimate *big.Int, _gasBufferUnmeasured *big.Int) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.ChangeGasParameters(&_IAuctionEntryPoint.TransactOpts, _gasPerByteIntrinsic, _gasPerByteEip7623, _gasContractExecution, _gasBufferEstimate, _gasBufferUnmeasured)
}

// ChangeGasParameters is a paid mutator transaction binding the contract method 0x2a215610.
//
// Solidity: function changeGasParameters(uint256 _gasPerByteIntrinsic, uint256 _gasPerByteEip7623, uint256 _gasContractExecution, uint256 _gasBufferEstimate, uint256 _gasBufferUnmeasured) returns()
func (_IAuctionEntryPoint *IAuctionEntryPointTransactorSession) ChangeGasParameters(_gasPerByteIntrinsic *big.Int, _gasPerByteEip7623 *big.Int, _gasContractExecution *big.Int, _gasBufferEstimate *big.Int, _gasBufferUnmeasured *big.Int) (*types.Transaction, error) {
	return _IAuctionEntryPoint.Contract.ChangeGasParameters(&_IAuctionEntryPoint.TransactOpts, _gasPerByteIntrinsic, _gasPerByteEip7623, _gasContractExecution, _gasBufferEstimate, _gasBufferUnmeasured)
}

// IAuctionEntryPointCallIterator is returned from FilterCall and is used to iterate over the raw logs and unpacked data for Call events raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointCallIterator struct {
	Event *IAuctionEntryPointCall // Event containing the contract specifics and raw log

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
func (it *IAuctionEntryPointCallIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionEntryPointCall)
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
		it.Event = new(IAuctionEntryPointCall)
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
func (it *IAuctionEntryPointCallIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionEntryPointCallIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionEntryPointCall represents a Call event raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointCall struct {
	Sender common.Address
	Nonce  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCall is a free log retrieval operation binding the contract event 0x9e4369a013b5e45a103a959d8eb70f15c55dc496e8335f245690393dfb4b71d4.
//
// Solidity: event Call(address sender, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) FilterCall(opts *bind.FilterOpts) (*IAuctionEntryPointCallIterator, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.FilterLogs(opts, "Call")
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointCallIterator{contract: _IAuctionEntryPoint.contract, event: "Call", logs: logs, sub: sub}, nil
}

// WatchCall is a free log subscription operation binding the contract event 0x9e4369a013b5e45a103a959d8eb70f15c55dc496e8335f245690393dfb4b71d4.
//
// Solidity: event Call(address sender, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) WatchCall(opts *bind.WatchOpts, sink chan<- *IAuctionEntryPointCall) (event.Subscription, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.WatchLogs(opts, "Call")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionEntryPointCall)
				if err := _IAuctionEntryPoint.contract.UnpackLog(event, "Call", log); err != nil {
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

// ParseCall is a log parse operation binding the contract event 0x9e4369a013b5e45a103a959d8eb70f15c55dc496e8335f245690393dfb4b71d4.
//
// Solidity: event Call(address sender, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) ParseCall(log types.Log) (*IAuctionEntryPointCall, error) {
	event := new(IAuctionEntryPointCall)
	if err := _IAuctionEntryPoint.contract.UnpackLog(event, "Call", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionEntryPointCallFailedIterator is returned from FilterCallFailed and is used to iterate over the raw logs and unpacked data for CallFailed events raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointCallFailedIterator struct {
	Event *IAuctionEntryPointCallFailed // Event containing the contract specifics and raw log

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
func (it *IAuctionEntryPointCallFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionEntryPointCallFailed)
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
		it.Event = new(IAuctionEntryPointCallFailed)
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
func (it *IAuctionEntryPointCallFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionEntryPointCallFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionEntryPointCallFailed represents a CallFailed event raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointCallFailed struct {
	Sender common.Address
	Nonce  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterCallFailed is a free log retrieval operation binding the contract event 0xb9eaeae386d339f8115782f297a9e5f0e13fb587cd6b0d502f113cb8dd4d6cb0.
//
// Solidity: event CallFailed(address sender, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) FilterCallFailed(opts *bind.FilterOpts) (*IAuctionEntryPointCallFailedIterator, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.FilterLogs(opts, "CallFailed")
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointCallFailedIterator{contract: _IAuctionEntryPoint.contract, event: "CallFailed", logs: logs, sub: sub}, nil
}

// WatchCallFailed is a free log subscription operation binding the contract event 0xb9eaeae386d339f8115782f297a9e5f0e13fb587cd6b0d502f113cb8dd4d6cb0.
//
// Solidity: event CallFailed(address sender, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) WatchCallFailed(opts *bind.WatchOpts, sink chan<- *IAuctionEntryPointCallFailed) (event.Subscription, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.WatchLogs(opts, "CallFailed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionEntryPointCallFailed)
				if err := _IAuctionEntryPoint.contract.UnpackLog(event, "CallFailed", log); err != nil {
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

// ParseCallFailed is a log parse operation binding the contract event 0xb9eaeae386d339f8115782f297a9e5f0e13fb587cd6b0d502f113cb8dd4d6cb0.
//
// Solidity: event CallFailed(address sender, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) ParseCallFailed(log types.Log) (*IAuctionEntryPointCallFailed, error) {
	event := new(IAuctionEntryPointCallFailed)
	if err := _IAuctionEntryPoint.contract.UnpackLog(event, "CallFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionEntryPointChangeAuctioneerIterator is returned from FilterChangeAuctioneer and is used to iterate over the raw logs and unpacked data for ChangeAuctioneer events raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointChangeAuctioneerIterator struct {
	Event *IAuctionEntryPointChangeAuctioneer // Event containing the contract specifics and raw log

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
func (it *IAuctionEntryPointChangeAuctioneerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionEntryPointChangeAuctioneer)
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
		it.Event = new(IAuctionEntryPointChangeAuctioneer)
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
func (it *IAuctionEntryPointChangeAuctioneerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionEntryPointChangeAuctioneerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionEntryPointChangeAuctioneer represents a ChangeAuctioneer event raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointChangeAuctioneer struct {
	OldAuctioneer common.Address
	NewAuctioneer common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterChangeAuctioneer is a free log retrieval operation binding the contract event 0xc8a0473779a405800019e7658500d0930d92f78e776b553d286baeefb9c9a0f1.
//
// Solidity: event ChangeAuctioneer(address oldAuctioneer, address newAuctioneer)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) FilterChangeAuctioneer(opts *bind.FilterOpts) (*IAuctionEntryPointChangeAuctioneerIterator, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.FilterLogs(opts, "ChangeAuctioneer")
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointChangeAuctioneerIterator{contract: _IAuctionEntryPoint.contract, event: "ChangeAuctioneer", logs: logs, sub: sub}, nil
}

// WatchChangeAuctioneer is a free log subscription operation binding the contract event 0xc8a0473779a405800019e7658500d0930d92f78e776b553d286baeefb9c9a0f1.
//
// Solidity: event ChangeAuctioneer(address oldAuctioneer, address newAuctioneer)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) WatchChangeAuctioneer(opts *bind.WatchOpts, sink chan<- *IAuctionEntryPointChangeAuctioneer) (event.Subscription, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.WatchLogs(opts, "ChangeAuctioneer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionEntryPointChangeAuctioneer)
				if err := _IAuctionEntryPoint.contract.UnpackLog(event, "ChangeAuctioneer", log); err != nil {
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

// ParseChangeAuctioneer is a log parse operation binding the contract event 0xc8a0473779a405800019e7658500d0930d92f78e776b553d286baeefb9c9a0f1.
//
// Solidity: event ChangeAuctioneer(address oldAuctioneer, address newAuctioneer)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) ParseChangeAuctioneer(log types.Log) (*IAuctionEntryPointChangeAuctioneer, error) {
	event := new(IAuctionEntryPointChangeAuctioneer)
	if err := _IAuctionEntryPoint.contract.UnpackLog(event, "ChangeAuctioneer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionEntryPointChangeDepositVaultIterator is returned from FilterChangeDepositVault and is used to iterate over the raw logs and unpacked data for ChangeDepositVault events raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointChangeDepositVaultIterator struct {
	Event *IAuctionEntryPointChangeDepositVault // Event containing the contract specifics and raw log

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
func (it *IAuctionEntryPointChangeDepositVaultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionEntryPointChangeDepositVault)
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
		it.Event = new(IAuctionEntryPointChangeDepositVault)
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
func (it *IAuctionEntryPointChangeDepositVaultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionEntryPointChangeDepositVaultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionEntryPointChangeDepositVault represents a ChangeDepositVault event raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointChangeDepositVault struct {
	OldDepositVault common.Address
	NewDepositVault common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterChangeDepositVault is a free log retrieval operation binding the contract event 0x32718750012b831b1e21ea05c1ed833c9365fc94e7316126303f6d09f41beb5d.
//
// Solidity: event ChangeDepositVault(address oldDepositVault, address newDepositVault)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) FilterChangeDepositVault(opts *bind.FilterOpts) (*IAuctionEntryPointChangeDepositVaultIterator, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.FilterLogs(opts, "ChangeDepositVault")
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointChangeDepositVaultIterator{contract: _IAuctionEntryPoint.contract, event: "ChangeDepositVault", logs: logs, sub: sub}, nil
}

// WatchChangeDepositVault is a free log subscription operation binding the contract event 0x32718750012b831b1e21ea05c1ed833c9365fc94e7316126303f6d09f41beb5d.
//
// Solidity: event ChangeDepositVault(address oldDepositVault, address newDepositVault)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) WatchChangeDepositVault(opts *bind.WatchOpts, sink chan<- *IAuctionEntryPointChangeDepositVault) (event.Subscription, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.WatchLogs(opts, "ChangeDepositVault")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionEntryPointChangeDepositVault)
				if err := _IAuctionEntryPoint.contract.UnpackLog(event, "ChangeDepositVault", log); err != nil {
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

// ParseChangeDepositVault is a log parse operation binding the contract event 0x32718750012b831b1e21ea05c1ed833c9365fc94e7316126303f6d09f41beb5d.
//
// Solidity: event ChangeDepositVault(address oldDepositVault, address newDepositVault)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) ParseChangeDepositVault(log types.Log) (*IAuctionEntryPointChangeDepositVault, error) {
	event := new(IAuctionEntryPointChangeDepositVault)
	if err := _IAuctionEntryPoint.contract.UnpackLog(event, "ChangeDepositVault", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionEntryPointChangeGasParametersIterator is returned from FilterChangeGasParameters and is used to iterate over the raw logs and unpacked data for ChangeGasParameters events raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointChangeGasParametersIterator struct {
	Event *IAuctionEntryPointChangeGasParameters // Event containing the contract specifics and raw log

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
func (it *IAuctionEntryPointChangeGasParametersIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionEntryPointChangeGasParameters)
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
		it.Event = new(IAuctionEntryPointChangeGasParameters)
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
func (it *IAuctionEntryPointChangeGasParametersIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionEntryPointChangeGasParametersIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionEntryPointChangeGasParameters represents a ChangeGasParameters event raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointChangeGasParameters struct {
	GasPerByteIntrinsic  *big.Int
	GasPerByteEip7623    *big.Int
	GasContractExecution *big.Int
	GasBufferEstimate    *big.Int
	GasBufferUnmeasured  *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterChangeGasParameters is a free log retrieval operation binding the contract event 0x88fa9ca37e3311ae53f6c07c267edc3aae4fc605df62318b5215d44665eb0308.
//
// Solidity: event ChangeGasParameters(uint256 gasPerByteIntrinsic, uint256 gasPerByteEip7623, uint256 gasContractExecution, uint256 gasBufferEstimate, uint256 gasBufferUnmeasured)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) FilterChangeGasParameters(opts *bind.FilterOpts) (*IAuctionEntryPointChangeGasParametersIterator, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.FilterLogs(opts, "ChangeGasParameters")
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointChangeGasParametersIterator{contract: _IAuctionEntryPoint.contract, event: "ChangeGasParameters", logs: logs, sub: sub}, nil
}

// WatchChangeGasParameters is a free log subscription operation binding the contract event 0x88fa9ca37e3311ae53f6c07c267edc3aae4fc605df62318b5215d44665eb0308.
//
// Solidity: event ChangeGasParameters(uint256 gasPerByteIntrinsic, uint256 gasPerByteEip7623, uint256 gasContractExecution, uint256 gasBufferEstimate, uint256 gasBufferUnmeasured)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) WatchChangeGasParameters(opts *bind.WatchOpts, sink chan<- *IAuctionEntryPointChangeGasParameters) (event.Subscription, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.WatchLogs(opts, "ChangeGasParameters")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionEntryPointChangeGasParameters)
				if err := _IAuctionEntryPoint.contract.UnpackLog(event, "ChangeGasParameters", log); err != nil {
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

// ParseChangeGasParameters is a log parse operation binding the contract event 0x88fa9ca37e3311ae53f6c07c267edc3aae4fc605df62318b5215d44665eb0308.
//
// Solidity: event ChangeGasParameters(uint256 gasPerByteIntrinsic, uint256 gasPerByteEip7623, uint256 gasContractExecution, uint256 gasBufferEstimate, uint256 gasBufferUnmeasured)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) ParseChangeGasParameters(log types.Log) (*IAuctionEntryPointChangeGasParameters, error) {
	event := new(IAuctionEntryPointChangeGasParameters)
	if err := _IAuctionEntryPoint.contract.UnpackLog(event, "ChangeGasParameters", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionEntryPointUseNonceIterator is returned from FilterUseNonce and is used to iterate over the raw logs and unpacked data for UseNonce events raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointUseNonceIterator struct {
	Event *IAuctionEntryPointUseNonce // Event containing the contract specifics and raw log

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
func (it *IAuctionEntryPointUseNonceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionEntryPointUseNonce)
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
		it.Event = new(IAuctionEntryPointUseNonce)
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
func (it *IAuctionEntryPointUseNonceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionEntryPointUseNonceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionEntryPointUseNonce represents a UseNonce event raised by the IAuctionEntryPoint contract.
type IAuctionEntryPointUseNonce struct {
	Searcher common.Address
	Nonce    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUseNonce is a free log retrieval operation binding the contract event 0x7243aab3c75cebfd96dc41ece7762eb023fc8c81e9ba2ce1d487876419b918b5.
//
// Solidity: event UseNonce(address searcher, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) FilterUseNonce(opts *bind.FilterOpts) (*IAuctionEntryPointUseNonceIterator, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.FilterLogs(opts, "UseNonce")
	if err != nil {
		return nil, err
	}
	return &IAuctionEntryPointUseNonceIterator{contract: _IAuctionEntryPoint.contract, event: "UseNonce", logs: logs, sub: sub}, nil
}

// WatchUseNonce is a free log subscription operation binding the contract event 0x7243aab3c75cebfd96dc41ece7762eb023fc8c81e9ba2ce1d487876419b918b5.
//
// Solidity: event UseNonce(address searcher, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) WatchUseNonce(opts *bind.WatchOpts, sink chan<- *IAuctionEntryPointUseNonce) (event.Subscription, error) {
	logs, sub, err := _IAuctionEntryPoint.contract.WatchLogs(opts, "UseNonce")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionEntryPointUseNonce)
				if err := _IAuctionEntryPoint.contract.UnpackLog(event, "UseNonce", log); err != nil {
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

// ParseUseNonce is a log parse operation binding the contract event 0x7243aab3c75cebfd96dc41ece7762eb023fc8c81e9ba2ce1d487876419b918b5.
//
// Solidity: event UseNonce(address searcher, uint256 nonce)
func (_IAuctionEntryPoint *IAuctionEntryPointFilterer) ParseUseNonce(log types.Log) (*IAuctionEntryPointUseNonce, error) {
	event := new(IAuctionEntryPointUseNonce)
	if err := _IAuctionEntryPoint.contract.UnpackLog(event, "UseNonce", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
