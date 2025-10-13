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

// ContextMetaData contains all meta data concerning the Context contract.
var ContextMetaData = &bind.MetaData{
	ABI: "[]",
}

// ContextABI is the input ABI used to generate the binding from.
// Deprecated: Use ContextMetaData.ABI instead.
var ContextABI = ContextMetaData.ABI

// ContextBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ContextBinRuntime = ``

// Context is an auto generated Go binding around a Kaia contract.
type Context struct {
	ContextCaller     // Read-only binding to the contract
	ContextTransactor // Write-only binding to the contract
	ContextFilterer   // Log filterer for contract events
}

// ContextCaller is an auto generated read-only Go binding around a Kaia contract.
type ContextCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextTransactor is an auto generated write-only Go binding around a Kaia contract.
type ContextTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ContextFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ContextSession struct {
	Contract     *Context          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContextCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ContextCallerSession struct {
	Contract *ContextCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ContextTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ContextTransactorSession struct {
	Contract     *ContextTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ContextRaw is an auto generated low-level Go binding around a Kaia contract.
type ContextRaw struct {
	Contract *Context // Generic contract binding to access the raw methods on
}

// ContextCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ContextCallerRaw struct {
	Contract *ContextCaller // Generic read-only contract binding to access the raw methods on
}

// ContextTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ContextTransactorRaw struct {
	Contract *ContextTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContext creates a new instance of Context, bound to a specific deployed contract.
func NewContext(address common.Address, backend bind.ContractBackend) (*Context, error) {
	contract, err := bindContext(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Context{ContextCaller: ContextCaller{contract: contract}, ContextTransactor: ContextTransactor{contract: contract}, ContextFilterer: ContextFilterer{contract: contract}}, nil
}

// NewContextCaller creates a new read-only instance of Context, bound to a specific deployed contract.
func NewContextCaller(address common.Address, caller bind.ContractCaller) (*ContextCaller, error) {
	contract, err := bindContext(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextCaller{contract: contract}, nil
}

// NewContextTransactor creates a new write-only instance of Context, bound to a specific deployed contract.
func NewContextTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextTransactor, error) {
	contract, err := bindContext(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextTransactor{contract: contract}, nil
}

// NewContextFilterer creates a new log filterer instance of Context, bound to a specific deployed contract.
func NewContextFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextFilterer, error) {
	contract, err := bindContext(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextFilterer{contract: contract}, nil
}

// bindContext binds a generic wrapper to an already deployed contract.
func bindContext(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContextMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Context.Contract.ContextCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Context.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.contract.Transact(opts, method, params...)
}

// IAuctionDepositVaultMetaData contains all meta data concerning the IAuctionDepositVault contract.
var IAuctionDepositVaultMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldFeeVault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeVault\",\"type\":\"address\"}],\"name\":\"ChangeAuctionFeeVault\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"ChangeMinDepositAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldLocktime\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newLocktime\",\"type\":\"uint256\"}],\"name\":\"ChangeMinWithdrawLocktime\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"InsufficientBalance\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TakenBid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TakenBidFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAmount\",\"type\":\"uint256\"}],\"name\":\"TakenGas\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasAmount\",\"type\":\"uint256\"}],\"name\":\"TakenGasFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"VaultDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"VaultReserveWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"VaultWithdraw\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAuctionFeeVault\",\"type\":\"address\"}],\"name\":\"changeAuctionFeeVault\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newMinAmount\",\"type\":\"uint256\"}],\"name\":\"changeMinDepositAmount\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newLocktime\",\"type\":\"uint256\"}],\"name\":\"changeMinWithdrawLocktime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"}],\"name\":\"depositBalances\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"}],\"name\":\"depositFor\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getAllAddrsOverMinDeposit\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"searchers\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"depositAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"nonces\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getDepositAddrs\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"searchers\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDepositAddrsLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"}],\"name\":\"isMinDepositOver\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reserveWithdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"takeBid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasUsed\",\"type\":\"uint256\"}],\"name\":\"takeGas\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"c2cf75ec": "changeAuctionFeeVault(address)",
		"8caad7b5": "changeMinDepositAmount(uint256)",
		"39253624": "changeMinWithdrawLocktime(uint256)",
		"d0e30db0": "deposit()",
		"1eb903cf": "depositBalances(address)",
		"aa67c919": "depositFor(address)",
		"5be1a55e": "getAllAddrsOverMinDeposit(uint256,uint256)",
		"e45076ac": "getDepositAddrs(uint256,uint256)",
		"84792e0b": "getDepositAddrsLength()",
		"48f928e8": "isMinDepositOver(address)",
		"bea39cab": "reserveWithdraw()",
		"08bf759d": "takeBid(address,uint256)",
		"b4fdfae5": "takeGas(address,uint256)",
		"3ccfd60b": "withdraw()",
	},
}

// IAuctionDepositVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use IAuctionDepositVaultMetaData.ABI instead.
var IAuctionDepositVaultABI = IAuctionDepositVaultMetaData.ABI

// IAuctionDepositVaultBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IAuctionDepositVaultBinRuntime = ``

// Deprecated: Use IAuctionDepositVaultMetaData.Sigs instead.
// IAuctionDepositVaultFuncSigs maps the 4-byte function signature to its string representation.
var IAuctionDepositVaultFuncSigs = IAuctionDepositVaultMetaData.Sigs

// IAuctionDepositVault is an auto generated Go binding around a Kaia contract.
type IAuctionDepositVault struct {
	IAuctionDepositVaultCaller     // Read-only binding to the contract
	IAuctionDepositVaultTransactor // Write-only binding to the contract
	IAuctionDepositVaultFilterer   // Log filterer for contract events
}

// IAuctionDepositVaultCaller is an auto generated read-only Go binding around a Kaia contract.
type IAuctionDepositVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionDepositVaultTransactor is an auto generated write-only Go binding around a Kaia contract.
type IAuctionDepositVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionDepositVaultFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IAuctionDepositVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionDepositVaultSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IAuctionDepositVaultSession struct {
	Contract     *IAuctionDepositVault // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// IAuctionDepositVaultCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IAuctionDepositVaultCallerSession struct {
	Contract *IAuctionDepositVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// IAuctionDepositVaultTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IAuctionDepositVaultTransactorSession struct {
	Contract     *IAuctionDepositVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// IAuctionDepositVaultRaw is an auto generated low-level Go binding around a Kaia contract.
type IAuctionDepositVaultRaw struct {
	Contract *IAuctionDepositVault // Generic contract binding to access the raw methods on
}

// IAuctionDepositVaultCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IAuctionDepositVaultCallerRaw struct {
	Contract *IAuctionDepositVaultCaller // Generic read-only contract binding to access the raw methods on
}

// IAuctionDepositVaultTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IAuctionDepositVaultTransactorRaw struct {
	Contract *IAuctionDepositVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAuctionDepositVault creates a new instance of IAuctionDepositVault, bound to a specific deployed contract.
func NewIAuctionDepositVault(address common.Address, backend bind.ContractBackend) (*IAuctionDepositVault, error) {
	contract, err := bindIAuctionDepositVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVault{IAuctionDepositVaultCaller: IAuctionDepositVaultCaller{contract: contract}, IAuctionDepositVaultTransactor: IAuctionDepositVaultTransactor{contract: contract}, IAuctionDepositVaultFilterer: IAuctionDepositVaultFilterer{contract: contract}}, nil
}

// NewIAuctionDepositVaultCaller creates a new read-only instance of IAuctionDepositVault, bound to a specific deployed contract.
func NewIAuctionDepositVaultCaller(address common.Address, caller bind.ContractCaller) (*IAuctionDepositVaultCaller, error) {
	contract, err := bindIAuctionDepositVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultCaller{contract: contract}, nil
}

// NewIAuctionDepositVaultTransactor creates a new write-only instance of IAuctionDepositVault, bound to a specific deployed contract.
func NewIAuctionDepositVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*IAuctionDepositVaultTransactor, error) {
	contract, err := bindIAuctionDepositVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultTransactor{contract: contract}, nil
}

// NewIAuctionDepositVaultFilterer creates a new log filterer instance of IAuctionDepositVault, bound to a specific deployed contract.
func NewIAuctionDepositVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*IAuctionDepositVaultFilterer, error) {
	contract, err := bindIAuctionDepositVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultFilterer{contract: contract}, nil
}

// bindIAuctionDepositVault binds a generic wrapper to an already deployed contract.
func bindIAuctionDepositVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAuctionDepositVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAuctionDepositVault *IAuctionDepositVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAuctionDepositVault.Contract.IAuctionDepositVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAuctionDepositVault *IAuctionDepositVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.IAuctionDepositVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAuctionDepositVault *IAuctionDepositVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.IAuctionDepositVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAuctionDepositVault *IAuctionDepositVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAuctionDepositVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.contract.Transact(opts, method, params...)
}

// DepositBalances is a free data retrieval call binding the contract method 0x1eb903cf.
//
// Solidity: function depositBalances(address searcher) view returns(uint256)
func (_IAuctionDepositVault *IAuctionDepositVaultCaller) DepositBalances(opts *bind.CallOpts, searcher common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IAuctionDepositVault.contract.Call(opts, &out, "depositBalances", searcher)
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// DepositBalances is a free data retrieval call binding the contract method 0x1eb903cf.
//
// Solidity: function depositBalances(address searcher) view returns(uint256)
func (_IAuctionDepositVault *IAuctionDepositVaultSession) DepositBalances(searcher common.Address) (*big.Int, error) {
	return _IAuctionDepositVault.Contract.DepositBalances(&_IAuctionDepositVault.CallOpts, searcher)
}

// DepositBalances is a free data retrieval call binding the contract method 0x1eb903cf.
//
// Solidity: function depositBalances(address searcher) view returns(uint256)
func (_IAuctionDepositVault *IAuctionDepositVaultCallerSession) DepositBalances(searcher common.Address) (*big.Int, error) {
	return _IAuctionDepositVault.Contract.DepositBalances(&_IAuctionDepositVault.CallOpts, searcher)
}

// GetAllAddrsOverMinDeposit is a free data retrieval call binding the contract method 0x5be1a55e.
//
// Solidity: function getAllAddrsOverMinDeposit(uint256 start, uint256 limit) view returns(address[] searchers, uint256[] depositAmounts, uint256[] nonces)
func (_IAuctionDepositVault *IAuctionDepositVaultCaller) GetAllAddrsOverMinDeposit(opts *bind.CallOpts, start *big.Int, limit *big.Int) (struct {
	Searchers      []common.Address
	DepositAmounts []*big.Int
	Nonces         []*big.Int
}, error,
) {
	var out []interface{}
	err := _IAuctionDepositVault.contract.Call(opts, &out, "getAllAddrsOverMinDeposit", start, limit)

	outstruct := new(struct {
		Searchers      []common.Address
		DepositAmounts []*big.Int
		Nonces         []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Searchers = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.DepositAmounts = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.Nonces = *abi.ConvertType(out[2], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err
}

// GetAllAddrsOverMinDeposit is a free data retrieval call binding the contract method 0x5be1a55e.
//
// Solidity: function getAllAddrsOverMinDeposit(uint256 start, uint256 limit) view returns(address[] searchers, uint256[] depositAmounts, uint256[] nonces)
func (_IAuctionDepositVault *IAuctionDepositVaultSession) GetAllAddrsOverMinDeposit(start *big.Int, limit *big.Int) (struct {
	Searchers      []common.Address
	DepositAmounts []*big.Int
	Nonces         []*big.Int
}, error,
) {
	return _IAuctionDepositVault.Contract.GetAllAddrsOverMinDeposit(&_IAuctionDepositVault.CallOpts, start, limit)
}

// GetAllAddrsOverMinDeposit is a free data retrieval call binding the contract method 0x5be1a55e.
//
// Solidity: function getAllAddrsOverMinDeposit(uint256 start, uint256 limit) view returns(address[] searchers, uint256[] depositAmounts, uint256[] nonces)
func (_IAuctionDepositVault *IAuctionDepositVaultCallerSession) GetAllAddrsOverMinDeposit(start *big.Int, limit *big.Int) (struct {
	Searchers      []common.Address
	DepositAmounts []*big.Int
	Nonces         []*big.Int
}, error,
) {
	return _IAuctionDepositVault.Contract.GetAllAddrsOverMinDeposit(&_IAuctionDepositVault.CallOpts, start, limit)
}

// GetDepositAddrs is a free data retrieval call binding the contract method 0xe45076ac.
//
// Solidity: function getDepositAddrs(uint256 start, uint256 limit) view returns(address[] searchers)
func (_IAuctionDepositVault *IAuctionDepositVaultCaller) GetDepositAddrs(opts *bind.CallOpts, start *big.Int, limit *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _IAuctionDepositVault.contract.Call(opts, &out, "getDepositAddrs", start, limit)
	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err
}

// GetDepositAddrs is a free data retrieval call binding the contract method 0xe45076ac.
//
// Solidity: function getDepositAddrs(uint256 start, uint256 limit) view returns(address[] searchers)
func (_IAuctionDepositVault *IAuctionDepositVaultSession) GetDepositAddrs(start *big.Int, limit *big.Int) ([]common.Address, error) {
	return _IAuctionDepositVault.Contract.GetDepositAddrs(&_IAuctionDepositVault.CallOpts, start, limit)
}

// GetDepositAddrs is a free data retrieval call binding the contract method 0xe45076ac.
//
// Solidity: function getDepositAddrs(uint256 start, uint256 limit) view returns(address[] searchers)
func (_IAuctionDepositVault *IAuctionDepositVaultCallerSession) GetDepositAddrs(start *big.Int, limit *big.Int) ([]common.Address, error) {
	return _IAuctionDepositVault.Contract.GetDepositAddrs(&_IAuctionDepositVault.CallOpts, start, limit)
}

// GetDepositAddrsLength is a free data retrieval call binding the contract method 0x84792e0b.
//
// Solidity: function getDepositAddrsLength() view returns(uint256)
func (_IAuctionDepositVault *IAuctionDepositVaultCaller) GetDepositAddrsLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _IAuctionDepositVault.contract.Call(opts, &out, "getDepositAddrsLength")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// GetDepositAddrsLength is a free data retrieval call binding the contract method 0x84792e0b.
//
// Solidity: function getDepositAddrsLength() view returns(uint256)
func (_IAuctionDepositVault *IAuctionDepositVaultSession) GetDepositAddrsLength() (*big.Int, error) {
	return _IAuctionDepositVault.Contract.GetDepositAddrsLength(&_IAuctionDepositVault.CallOpts)
}

// GetDepositAddrsLength is a free data retrieval call binding the contract method 0x84792e0b.
//
// Solidity: function getDepositAddrsLength() view returns(uint256)
func (_IAuctionDepositVault *IAuctionDepositVaultCallerSession) GetDepositAddrsLength() (*big.Int, error) {
	return _IAuctionDepositVault.Contract.GetDepositAddrsLength(&_IAuctionDepositVault.CallOpts)
}

// IsMinDepositOver is a free data retrieval call binding the contract method 0x48f928e8.
//
// Solidity: function isMinDepositOver(address searcher) view returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultCaller) IsMinDepositOver(opts *bind.CallOpts, searcher common.Address) (bool, error) {
	var out []interface{}
	err := _IAuctionDepositVault.contract.Call(opts, &out, "isMinDepositOver", searcher)
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// IsMinDepositOver is a free data retrieval call binding the contract method 0x48f928e8.
//
// Solidity: function isMinDepositOver(address searcher) view returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultSession) IsMinDepositOver(searcher common.Address) (bool, error) {
	return _IAuctionDepositVault.Contract.IsMinDepositOver(&_IAuctionDepositVault.CallOpts, searcher)
}

// IsMinDepositOver is a free data retrieval call binding the contract method 0x48f928e8.
//
// Solidity: function isMinDepositOver(address searcher) view returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultCallerSession) IsMinDepositOver(searcher common.Address) (bool, error) {
	return _IAuctionDepositVault.Contract.IsMinDepositOver(&_IAuctionDepositVault.CallOpts, searcher)
}

// ChangeAuctionFeeVault is a paid mutator transaction binding the contract method 0xc2cf75ec.
//
// Solidity: function changeAuctionFeeVault(address newAuctionFeeVault) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) ChangeAuctionFeeVault(opts *bind.TransactOpts, newAuctionFeeVault common.Address) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "changeAuctionFeeVault", newAuctionFeeVault)
}

// ChangeAuctionFeeVault is a paid mutator transaction binding the contract method 0xc2cf75ec.
//
// Solidity: function changeAuctionFeeVault(address newAuctionFeeVault) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultSession) ChangeAuctionFeeVault(newAuctionFeeVault common.Address) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.ChangeAuctionFeeVault(&_IAuctionDepositVault.TransactOpts, newAuctionFeeVault)
}

// ChangeAuctionFeeVault is a paid mutator transaction binding the contract method 0xc2cf75ec.
//
// Solidity: function changeAuctionFeeVault(address newAuctionFeeVault) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) ChangeAuctionFeeVault(newAuctionFeeVault common.Address) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.ChangeAuctionFeeVault(&_IAuctionDepositVault.TransactOpts, newAuctionFeeVault)
}

// ChangeMinDepositAmount is a paid mutator transaction binding the contract method 0x8caad7b5.
//
// Solidity: function changeMinDepositAmount(uint256 newMinAmount) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) ChangeMinDepositAmount(opts *bind.TransactOpts, newMinAmount *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "changeMinDepositAmount", newMinAmount)
}

// ChangeMinDepositAmount is a paid mutator transaction binding the contract method 0x8caad7b5.
//
// Solidity: function changeMinDepositAmount(uint256 newMinAmount) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultSession) ChangeMinDepositAmount(newMinAmount *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.ChangeMinDepositAmount(&_IAuctionDepositVault.TransactOpts, newMinAmount)
}

// ChangeMinDepositAmount is a paid mutator transaction binding the contract method 0x8caad7b5.
//
// Solidity: function changeMinDepositAmount(uint256 newMinAmount) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) ChangeMinDepositAmount(newMinAmount *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.ChangeMinDepositAmount(&_IAuctionDepositVault.TransactOpts, newMinAmount)
}

// ChangeMinWithdrawLocktime is a paid mutator transaction binding the contract method 0x39253624.
//
// Solidity: function changeMinWithdrawLocktime(uint256 newLocktime) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) ChangeMinWithdrawLocktime(opts *bind.TransactOpts, newLocktime *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "changeMinWithdrawLocktime", newLocktime)
}

// ChangeMinWithdrawLocktime is a paid mutator transaction binding the contract method 0x39253624.
//
// Solidity: function changeMinWithdrawLocktime(uint256 newLocktime) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultSession) ChangeMinWithdrawLocktime(newLocktime *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.ChangeMinWithdrawLocktime(&_IAuctionDepositVault.TransactOpts, newLocktime)
}

// ChangeMinWithdrawLocktime is a paid mutator transaction binding the contract method 0x39253624.
//
// Solidity: function changeMinWithdrawLocktime(uint256 newLocktime) returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) ChangeMinWithdrawLocktime(newLocktime *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.ChangeMinWithdrawLocktime(&_IAuctionDepositVault.TransactOpts, newLocktime)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_IAuctionDepositVault *IAuctionDepositVaultSession) Deposit() (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.Deposit(&_IAuctionDepositVault.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) Deposit() (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.Deposit(&_IAuctionDepositVault.TransactOpts)
}

// DepositFor is a paid mutator transaction binding the contract method 0xaa67c919.
//
// Solidity: function depositFor(address searcher) payable returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) DepositFor(opts *bind.TransactOpts, searcher common.Address) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "depositFor", searcher)
}

// DepositFor is a paid mutator transaction binding the contract method 0xaa67c919.
//
// Solidity: function depositFor(address searcher) payable returns()
func (_IAuctionDepositVault *IAuctionDepositVaultSession) DepositFor(searcher common.Address) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.DepositFor(&_IAuctionDepositVault.TransactOpts, searcher)
}

// DepositFor is a paid mutator transaction binding the contract method 0xaa67c919.
//
// Solidity: function depositFor(address searcher) payable returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) DepositFor(searcher common.Address) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.DepositFor(&_IAuctionDepositVault.TransactOpts, searcher)
}

// ReserveWithdraw is a paid mutator transaction binding the contract method 0xbea39cab.
//
// Solidity: function reserveWithdraw() returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) ReserveWithdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "reserveWithdraw")
}

// ReserveWithdraw is a paid mutator transaction binding the contract method 0xbea39cab.
//
// Solidity: function reserveWithdraw() returns()
func (_IAuctionDepositVault *IAuctionDepositVaultSession) ReserveWithdraw() (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.ReserveWithdraw(&_IAuctionDepositVault.TransactOpts)
}

// ReserveWithdraw is a paid mutator transaction binding the contract method 0xbea39cab.
//
// Solidity: function reserveWithdraw() returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) ReserveWithdraw() (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.ReserveWithdraw(&_IAuctionDepositVault.TransactOpts)
}

// TakeBid is a paid mutator transaction binding the contract method 0x08bf759d.
//
// Solidity: function takeBid(address searcher, uint256 amount) returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) TakeBid(opts *bind.TransactOpts, searcher common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "takeBid", searcher, amount)
}

// TakeBid is a paid mutator transaction binding the contract method 0x08bf759d.
//
// Solidity: function takeBid(address searcher, uint256 amount) returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultSession) TakeBid(searcher common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.TakeBid(&_IAuctionDepositVault.TransactOpts, searcher, amount)
}

// TakeBid is a paid mutator transaction binding the contract method 0x08bf759d.
//
// Solidity: function takeBid(address searcher, uint256 amount) returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) TakeBid(searcher common.Address, amount *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.TakeBid(&_IAuctionDepositVault.TransactOpts, searcher, amount)
}

// TakeGas is a paid mutator transaction binding the contract method 0xb4fdfae5.
//
// Solidity: function takeGas(address searcher, uint256 gasUsed) returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) TakeGas(opts *bind.TransactOpts, searcher common.Address, gasUsed *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "takeGas", searcher, gasUsed)
}

// TakeGas is a paid mutator transaction binding the contract method 0xb4fdfae5.
//
// Solidity: function takeGas(address searcher, uint256 gasUsed) returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultSession) TakeGas(searcher common.Address, gasUsed *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.TakeGas(&_IAuctionDepositVault.TransactOpts, searcher, gasUsed)
}

// TakeGas is a paid mutator transaction binding the contract method 0xb4fdfae5.
//
// Solidity: function takeGas(address searcher, uint256 gasUsed) returns(bool)
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) TakeGas(searcher common.Address, gasUsed *big.Int) (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.TakeGas(&_IAuctionDepositVault.TransactOpts, searcher, gasUsed)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionDepositVault.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_IAuctionDepositVault *IAuctionDepositVaultSession) Withdraw() (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.Withdraw(&_IAuctionDepositVault.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_IAuctionDepositVault *IAuctionDepositVaultTransactorSession) Withdraw() (*types.Transaction, error) {
	return _IAuctionDepositVault.Contract.Withdraw(&_IAuctionDepositVault.TransactOpts)
}

// IAuctionDepositVaultChangeAuctionFeeVaultIterator is returned from FilterChangeAuctionFeeVault and is used to iterate over the raw logs and unpacked data for ChangeAuctionFeeVault events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultChangeAuctionFeeVaultIterator struct {
	Event *IAuctionDepositVaultChangeAuctionFeeVault // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultChangeAuctionFeeVaultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultChangeAuctionFeeVault)
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
		it.Event = new(IAuctionDepositVaultChangeAuctionFeeVault)
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
func (it *IAuctionDepositVaultChangeAuctionFeeVaultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultChangeAuctionFeeVaultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultChangeAuctionFeeVault represents a ChangeAuctionFeeVault event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultChangeAuctionFeeVault struct {
	OldFeeVault common.Address
	NewFeeVault common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterChangeAuctionFeeVault is a free log retrieval operation binding the contract event 0x384b2820c6944bcf938c3d8c170a4421377f8e7a397ffb5ff321b8dce59bb97d.
//
// Solidity: event ChangeAuctionFeeVault(address oldFeeVault, address newFeeVault)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterChangeAuctionFeeVault(opts *bind.FilterOpts) (*IAuctionDepositVaultChangeAuctionFeeVaultIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "ChangeAuctionFeeVault")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultChangeAuctionFeeVaultIterator{contract: _IAuctionDepositVault.contract, event: "ChangeAuctionFeeVault", logs: logs, sub: sub}, nil
}

// WatchChangeAuctionFeeVault is a free log subscription operation binding the contract event 0x384b2820c6944bcf938c3d8c170a4421377f8e7a397ffb5ff321b8dce59bb97d.
//
// Solidity: event ChangeAuctionFeeVault(address oldFeeVault, address newFeeVault)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchChangeAuctionFeeVault(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultChangeAuctionFeeVault) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "ChangeAuctionFeeVault")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultChangeAuctionFeeVault)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "ChangeAuctionFeeVault", log); err != nil {
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

// ParseChangeAuctionFeeVault is a log parse operation binding the contract event 0x384b2820c6944bcf938c3d8c170a4421377f8e7a397ffb5ff321b8dce59bb97d.
//
// Solidity: event ChangeAuctionFeeVault(address oldFeeVault, address newFeeVault)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseChangeAuctionFeeVault(log types.Log) (*IAuctionDepositVaultChangeAuctionFeeVault, error) {
	event := new(IAuctionDepositVaultChangeAuctionFeeVault)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "ChangeAuctionFeeVault", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultChangeMinDepositAmountIterator is returned from FilterChangeMinDepositAmount and is used to iterate over the raw logs and unpacked data for ChangeMinDepositAmount events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultChangeMinDepositAmountIterator struct {
	Event *IAuctionDepositVaultChangeMinDepositAmount // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultChangeMinDepositAmountIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultChangeMinDepositAmount)
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
		it.Event = new(IAuctionDepositVaultChangeMinDepositAmount)
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
func (it *IAuctionDepositVaultChangeMinDepositAmountIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultChangeMinDepositAmountIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultChangeMinDepositAmount represents a ChangeMinDepositAmount event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultChangeMinDepositAmount struct {
	OldAmount *big.Int
	NewAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterChangeMinDepositAmount is a free log retrieval operation binding the contract event 0x404f17a8e681362e0d411e6b9b82111f6f7657596601fbe44eb9e2c960438f54.
//
// Solidity: event ChangeMinDepositAmount(uint256 oldAmount, uint256 newAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterChangeMinDepositAmount(opts *bind.FilterOpts) (*IAuctionDepositVaultChangeMinDepositAmountIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "ChangeMinDepositAmount")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultChangeMinDepositAmountIterator{contract: _IAuctionDepositVault.contract, event: "ChangeMinDepositAmount", logs: logs, sub: sub}, nil
}

// WatchChangeMinDepositAmount is a free log subscription operation binding the contract event 0x404f17a8e681362e0d411e6b9b82111f6f7657596601fbe44eb9e2c960438f54.
//
// Solidity: event ChangeMinDepositAmount(uint256 oldAmount, uint256 newAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchChangeMinDepositAmount(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultChangeMinDepositAmount) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "ChangeMinDepositAmount")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultChangeMinDepositAmount)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "ChangeMinDepositAmount", log); err != nil {
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

// ParseChangeMinDepositAmount is a log parse operation binding the contract event 0x404f17a8e681362e0d411e6b9b82111f6f7657596601fbe44eb9e2c960438f54.
//
// Solidity: event ChangeMinDepositAmount(uint256 oldAmount, uint256 newAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseChangeMinDepositAmount(log types.Log) (*IAuctionDepositVaultChangeMinDepositAmount, error) {
	event := new(IAuctionDepositVaultChangeMinDepositAmount)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "ChangeMinDepositAmount", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultChangeMinWithdrawLocktimeIterator is returned from FilterChangeMinWithdrawLocktime and is used to iterate over the raw logs and unpacked data for ChangeMinWithdrawLocktime events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultChangeMinWithdrawLocktimeIterator struct {
	Event *IAuctionDepositVaultChangeMinWithdrawLocktime // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultChangeMinWithdrawLocktimeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultChangeMinWithdrawLocktime)
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
		it.Event = new(IAuctionDepositVaultChangeMinWithdrawLocktime)
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
func (it *IAuctionDepositVaultChangeMinWithdrawLocktimeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultChangeMinWithdrawLocktimeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultChangeMinWithdrawLocktime represents a ChangeMinWithdrawLocktime event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultChangeMinWithdrawLocktime struct {
	OldLocktime *big.Int
	NewLocktime *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterChangeMinWithdrawLocktime is a free log retrieval operation binding the contract event 0x5e273c905caf5b3f4d3a98caf8814cdccb3860b2b334a58140b692fcf28b286d.
//
// Solidity: event ChangeMinWithdrawLocktime(uint256 oldLocktime, uint256 newLocktime)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterChangeMinWithdrawLocktime(opts *bind.FilterOpts) (*IAuctionDepositVaultChangeMinWithdrawLocktimeIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "ChangeMinWithdrawLocktime")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultChangeMinWithdrawLocktimeIterator{contract: _IAuctionDepositVault.contract, event: "ChangeMinWithdrawLocktime", logs: logs, sub: sub}, nil
}

// WatchChangeMinWithdrawLocktime is a free log subscription operation binding the contract event 0x5e273c905caf5b3f4d3a98caf8814cdccb3860b2b334a58140b692fcf28b286d.
//
// Solidity: event ChangeMinWithdrawLocktime(uint256 oldLocktime, uint256 newLocktime)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchChangeMinWithdrawLocktime(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultChangeMinWithdrawLocktime) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "ChangeMinWithdrawLocktime")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultChangeMinWithdrawLocktime)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "ChangeMinWithdrawLocktime", log); err != nil {
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

// ParseChangeMinWithdrawLocktime is a log parse operation binding the contract event 0x5e273c905caf5b3f4d3a98caf8814cdccb3860b2b334a58140b692fcf28b286d.
//
// Solidity: event ChangeMinWithdrawLocktime(uint256 oldLocktime, uint256 newLocktime)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseChangeMinWithdrawLocktime(log types.Log) (*IAuctionDepositVaultChangeMinWithdrawLocktime, error) {
	event := new(IAuctionDepositVaultChangeMinWithdrawLocktime)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "ChangeMinWithdrawLocktime", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultInsufficientBalanceIterator is returned from FilterInsufficientBalance and is used to iterate over the raw logs and unpacked data for InsufficientBalance events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultInsufficientBalanceIterator struct {
	Event *IAuctionDepositVaultInsufficientBalance // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultInsufficientBalanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultInsufficientBalance)
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
		it.Event = new(IAuctionDepositVaultInsufficientBalance)
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
func (it *IAuctionDepositVaultInsufficientBalanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultInsufficientBalanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultInsufficientBalance represents a InsufficientBalance event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultInsufficientBalance struct {
	Searcher common.Address
	Balance  *big.Int
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterInsufficientBalance is a free log retrieval operation binding the contract event 0xdb42144d928cd19733b9dfeead8bc02ed403845c274e06a6eee0944fb25ca5c4.
//
// Solidity: event InsufficientBalance(address searcher, uint256 balance, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterInsufficientBalance(opts *bind.FilterOpts) (*IAuctionDepositVaultInsufficientBalanceIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "InsufficientBalance")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultInsufficientBalanceIterator{contract: _IAuctionDepositVault.contract, event: "InsufficientBalance", logs: logs, sub: sub}, nil
}

// WatchInsufficientBalance is a free log subscription operation binding the contract event 0xdb42144d928cd19733b9dfeead8bc02ed403845c274e06a6eee0944fb25ca5c4.
//
// Solidity: event InsufficientBalance(address searcher, uint256 balance, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchInsufficientBalance(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultInsufficientBalance) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "InsufficientBalance")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultInsufficientBalance)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "InsufficientBalance", log); err != nil {
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

// ParseInsufficientBalance is a log parse operation binding the contract event 0xdb42144d928cd19733b9dfeead8bc02ed403845c274e06a6eee0944fb25ca5c4.
//
// Solidity: event InsufficientBalance(address searcher, uint256 balance, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseInsufficientBalance(log types.Log) (*IAuctionDepositVaultInsufficientBalance, error) {
	event := new(IAuctionDepositVaultInsufficientBalance)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "InsufficientBalance", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultTakenBidIterator is returned from FilterTakenBid and is used to iterate over the raw logs and unpacked data for TakenBid events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultTakenBidIterator struct {
	Event *IAuctionDepositVaultTakenBid // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultTakenBidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultTakenBid)
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
		it.Event = new(IAuctionDepositVaultTakenBid)
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
func (it *IAuctionDepositVaultTakenBidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultTakenBidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultTakenBid represents a TakenBid event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultTakenBid struct {
	Searcher common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTakenBid is a free log retrieval operation binding the contract event 0xbf3f09205213cdb412585918db004ac3947fcd2db1f85fdbb20faaa950e08d2e.
//
// Solidity: event TakenBid(address searcher, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterTakenBid(opts *bind.FilterOpts) (*IAuctionDepositVaultTakenBidIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "TakenBid")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultTakenBidIterator{contract: _IAuctionDepositVault.contract, event: "TakenBid", logs: logs, sub: sub}, nil
}

// WatchTakenBid is a free log subscription operation binding the contract event 0xbf3f09205213cdb412585918db004ac3947fcd2db1f85fdbb20faaa950e08d2e.
//
// Solidity: event TakenBid(address searcher, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchTakenBid(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultTakenBid) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "TakenBid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultTakenBid)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "TakenBid", log); err != nil {
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

// ParseTakenBid is a log parse operation binding the contract event 0xbf3f09205213cdb412585918db004ac3947fcd2db1f85fdbb20faaa950e08d2e.
//
// Solidity: event TakenBid(address searcher, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseTakenBid(log types.Log) (*IAuctionDepositVaultTakenBid, error) {
	event := new(IAuctionDepositVaultTakenBid)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "TakenBid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultTakenBidFailedIterator is returned from FilterTakenBidFailed and is used to iterate over the raw logs and unpacked data for TakenBidFailed events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultTakenBidFailedIterator struct {
	Event *IAuctionDepositVaultTakenBidFailed // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultTakenBidFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultTakenBidFailed)
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
		it.Event = new(IAuctionDepositVaultTakenBidFailed)
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
func (it *IAuctionDepositVaultTakenBidFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultTakenBidFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultTakenBidFailed represents a TakenBidFailed event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultTakenBidFailed struct {
	Searcher common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTakenBidFailed is a free log retrieval operation binding the contract event 0xd26b9fcb04629bad3365ea66f0f4d563e8c71753e7ee1a0c601f1a02ceebed8b.
//
// Solidity: event TakenBidFailed(address searcher, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterTakenBidFailed(opts *bind.FilterOpts) (*IAuctionDepositVaultTakenBidFailedIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "TakenBidFailed")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultTakenBidFailedIterator{contract: _IAuctionDepositVault.contract, event: "TakenBidFailed", logs: logs, sub: sub}, nil
}

// WatchTakenBidFailed is a free log subscription operation binding the contract event 0xd26b9fcb04629bad3365ea66f0f4d563e8c71753e7ee1a0c601f1a02ceebed8b.
//
// Solidity: event TakenBidFailed(address searcher, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchTakenBidFailed(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultTakenBidFailed) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "TakenBidFailed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultTakenBidFailed)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "TakenBidFailed", log); err != nil {
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

// ParseTakenBidFailed is a log parse operation binding the contract event 0xd26b9fcb04629bad3365ea66f0f4d563e8c71753e7ee1a0c601f1a02ceebed8b.
//
// Solidity: event TakenBidFailed(address searcher, uint256 amount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseTakenBidFailed(log types.Log) (*IAuctionDepositVaultTakenBidFailed, error) {
	event := new(IAuctionDepositVaultTakenBidFailed)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "TakenBidFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultTakenGasIterator is returned from FilterTakenGas and is used to iterate over the raw logs and unpacked data for TakenGas events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultTakenGasIterator struct {
	Event *IAuctionDepositVaultTakenGas // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultTakenGasIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultTakenGas)
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
		it.Event = new(IAuctionDepositVaultTakenGas)
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
func (it *IAuctionDepositVaultTakenGasIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultTakenGasIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultTakenGas represents a TakenGas event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultTakenGas struct {
	Searcher  common.Address
	GasAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTakenGas is a free log retrieval operation binding the contract event 0x672970b97597bfb4460c604bb4adc66538c7b00a40bfb2cb9b6566f1692266fd.
//
// Solidity: event TakenGas(address searcher, uint256 gasAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterTakenGas(opts *bind.FilterOpts) (*IAuctionDepositVaultTakenGasIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "TakenGas")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultTakenGasIterator{contract: _IAuctionDepositVault.contract, event: "TakenGas", logs: logs, sub: sub}, nil
}

// WatchTakenGas is a free log subscription operation binding the contract event 0x672970b97597bfb4460c604bb4adc66538c7b00a40bfb2cb9b6566f1692266fd.
//
// Solidity: event TakenGas(address searcher, uint256 gasAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchTakenGas(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultTakenGas) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "TakenGas")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultTakenGas)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "TakenGas", log); err != nil {
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

// ParseTakenGas is a log parse operation binding the contract event 0x672970b97597bfb4460c604bb4adc66538c7b00a40bfb2cb9b6566f1692266fd.
//
// Solidity: event TakenGas(address searcher, uint256 gasAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseTakenGas(log types.Log) (*IAuctionDepositVaultTakenGas, error) {
	event := new(IAuctionDepositVaultTakenGas)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "TakenGas", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultTakenGasFailedIterator is returned from FilterTakenGasFailed and is used to iterate over the raw logs and unpacked data for TakenGasFailed events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultTakenGasFailedIterator struct {
	Event *IAuctionDepositVaultTakenGasFailed // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultTakenGasFailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultTakenGasFailed)
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
		it.Event = new(IAuctionDepositVaultTakenGasFailed)
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
func (it *IAuctionDepositVaultTakenGasFailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultTakenGasFailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultTakenGasFailed represents a TakenGasFailed event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultTakenGasFailed struct {
	Searcher  common.Address
	GasAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTakenGasFailed is a free log retrieval operation binding the contract event 0x53605ab3492a9ae1542e049df64874ee628a589811e81b2594f416e9d092b7ca.
//
// Solidity: event TakenGasFailed(address searcher, uint256 gasAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterTakenGasFailed(opts *bind.FilterOpts) (*IAuctionDepositVaultTakenGasFailedIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "TakenGasFailed")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultTakenGasFailedIterator{contract: _IAuctionDepositVault.contract, event: "TakenGasFailed", logs: logs, sub: sub}, nil
}

// WatchTakenGasFailed is a free log subscription operation binding the contract event 0x53605ab3492a9ae1542e049df64874ee628a589811e81b2594f416e9d092b7ca.
//
// Solidity: event TakenGasFailed(address searcher, uint256 gasAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchTakenGasFailed(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultTakenGasFailed) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "TakenGasFailed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultTakenGasFailed)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "TakenGasFailed", log); err != nil {
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

// ParseTakenGasFailed is a log parse operation binding the contract event 0x53605ab3492a9ae1542e049df64874ee628a589811e81b2594f416e9d092b7ca.
//
// Solidity: event TakenGasFailed(address searcher, uint256 gasAmount)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseTakenGasFailed(log types.Log) (*IAuctionDepositVaultTakenGasFailed, error) {
	event := new(IAuctionDepositVaultTakenGasFailed)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "TakenGasFailed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultVaultDepositIterator is returned from FilterVaultDeposit and is used to iterate over the raw logs and unpacked data for VaultDeposit events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultVaultDepositIterator struct {
	Event *IAuctionDepositVaultVaultDeposit // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultVaultDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultVaultDeposit)
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
		it.Event = new(IAuctionDepositVaultVaultDeposit)
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
func (it *IAuctionDepositVaultVaultDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultVaultDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultVaultDeposit represents a VaultDeposit event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultVaultDeposit struct {
	Searcher    common.Address
	Amount      *big.Int
	TotalAmount *big.Int
	Nonce       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVaultDeposit is a free log retrieval operation binding the contract event 0xb356f70fb7268cd9ebf6cd2c25e8628582435a3754305d2edade9b81c0961181.
//
// Solidity: event VaultDeposit(address searcher, uint256 amount, uint256 totalAmount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterVaultDeposit(opts *bind.FilterOpts) (*IAuctionDepositVaultVaultDepositIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "VaultDeposit")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultVaultDepositIterator{contract: _IAuctionDepositVault.contract, event: "VaultDeposit", logs: logs, sub: sub}, nil
}

// WatchVaultDeposit is a free log subscription operation binding the contract event 0xb356f70fb7268cd9ebf6cd2c25e8628582435a3754305d2edade9b81c0961181.
//
// Solidity: event VaultDeposit(address searcher, uint256 amount, uint256 totalAmount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchVaultDeposit(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultVaultDeposit) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "VaultDeposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultVaultDeposit)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "VaultDeposit", log); err != nil {
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

// ParseVaultDeposit is a log parse operation binding the contract event 0xb356f70fb7268cd9ebf6cd2c25e8628582435a3754305d2edade9b81c0961181.
//
// Solidity: event VaultDeposit(address searcher, uint256 amount, uint256 totalAmount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseVaultDeposit(log types.Log) (*IAuctionDepositVaultVaultDeposit, error) {
	event := new(IAuctionDepositVaultVaultDeposit)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "VaultDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultVaultReserveWithdrawIterator is returned from FilterVaultReserveWithdraw and is used to iterate over the raw logs and unpacked data for VaultReserveWithdraw events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultVaultReserveWithdrawIterator struct {
	Event *IAuctionDepositVaultVaultReserveWithdraw // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultVaultReserveWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultVaultReserveWithdraw)
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
		it.Event = new(IAuctionDepositVaultVaultReserveWithdraw)
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
func (it *IAuctionDepositVaultVaultReserveWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultVaultReserveWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultVaultReserveWithdraw represents a VaultReserveWithdraw event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultVaultReserveWithdraw struct {
	Searcher common.Address
	Amount   *big.Int
	Nonce    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterVaultReserveWithdraw is a free log retrieval operation binding the contract event 0x89312dfc642f95af3c0376f434dd92d2ce54d713297a8df3921828edbd2ef1a1.
//
// Solidity: event VaultReserveWithdraw(address searcher, uint256 amount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterVaultReserveWithdraw(opts *bind.FilterOpts) (*IAuctionDepositVaultVaultReserveWithdrawIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "VaultReserveWithdraw")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultVaultReserveWithdrawIterator{contract: _IAuctionDepositVault.contract, event: "VaultReserveWithdraw", logs: logs, sub: sub}, nil
}

// WatchVaultReserveWithdraw is a free log subscription operation binding the contract event 0x89312dfc642f95af3c0376f434dd92d2ce54d713297a8df3921828edbd2ef1a1.
//
// Solidity: event VaultReserveWithdraw(address searcher, uint256 amount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchVaultReserveWithdraw(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultVaultReserveWithdraw) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "VaultReserveWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultVaultReserveWithdraw)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "VaultReserveWithdraw", log); err != nil {
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

// ParseVaultReserveWithdraw is a log parse operation binding the contract event 0x89312dfc642f95af3c0376f434dd92d2ce54d713297a8df3921828edbd2ef1a1.
//
// Solidity: event VaultReserveWithdraw(address searcher, uint256 amount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseVaultReserveWithdraw(log types.Log) (*IAuctionDepositVaultVaultReserveWithdraw, error) {
	event := new(IAuctionDepositVaultVaultReserveWithdraw)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "VaultReserveWithdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionDepositVaultVaultWithdrawIterator is returned from FilterVaultWithdraw and is used to iterate over the raw logs and unpacked data for VaultWithdraw events raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultVaultWithdrawIterator struct {
	Event *IAuctionDepositVaultVaultWithdraw // Event containing the contract specifics and raw log

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
func (it *IAuctionDepositVaultVaultWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionDepositVaultVaultWithdraw)
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
		it.Event = new(IAuctionDepositVaultVaultWithdraw)
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
func (it *IAuctionDepositVaultVaultWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionDepositVaultVaultWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionDepositVaultVaultWithdraw represents a VaultWithdraw event raised by the IAuctionDepositVault contract.
type IAuctionDepositVaultVaultWithdraw struct {
	Searcher common.Address
	Amount   *big.Int
	Nonce    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterVaultWithdraw is a free log retrieval operation binding the contract event 0x97c7b889c23f0ecb7c22eac5c58266eae022989cfdbd9b6dcf63417338c787de.
//
// Solidity: event VaultWithdraw(address searcher, uint256 amount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) FilterVaultWithdraw(opts *bind.FilterOpts) (*IAuctionDepositVaultVaultWithdrawIterator, error) {
	logs, sub, err := _IAuctionDepositVault.contract.FilterLogs(opts, "VaultWithdraw")
	if err != nil {
		return nil, err
	}
	return &IAuctionDepositVaultVaultWithdrawIterator{contract: _IAuctionDepositVault.contract, event: "VaultWithdraw", logs: logs, sub: sub}, nil
}

// WatchVaultWithdraw is a free log subscription operation binding the contract event 0x97c7b889c23f0ecb7c22eac5c58266eae022989cfdbd9b6dcf63417338c787de.
//
// Solidity: event VaultWithdraw(address searcher, uint256 amount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) WatchVaultWithdraw(opts *bind.WatchOpts, sink chan<- *IAuctionDepositVaultVaultWithdraw) (event.Subscription, error) {
	logs, sub, err := _IAuctionDepositVault.contract.WatchLogs(opts, "VaultWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionDepositVaultVaultWithdraw)
				if err := _IAuctionDepositVault.contract.UnpackLog(event, "VaultWithdraw", log); err != nil {
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

// ParseVaultWithdraw is a log parse operation binding the contract event 0x97c7b889c23f0ecb7c22eac5c58266eae022989cfdbd9b6dcf63417338c787de.
//
// Solidity: event VaultWithdraw(address searcher, uint256 amount, uint256 nonce)
func (_IAuctionDepositVault *IAuctionDepositVaultFilterer) ParseVaultWithdraw(log types.Log) (*IAuctionDepositVaultVaultWithdraw, error) {
	event := new(IAuctionDepositVaultVaultWithdraw)
	if err := _IAuctionDepositVault.contract.UnpackLog(event, "VaultWithdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionEntryPointMetaData contains all meta data concerning the IAuctionEntryPoint contract.
var IAuctionEntryPointMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"Call\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"CallFailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldAuctioneer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAuctioneer\",\"type\":\"address\"}],\"name\":\"ChangeAuctioneer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"oldDepositVault\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newDepositVault\",\"type\":\"address\"}],\"name\":\"ChangeDepositVault\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPerByteIntrinsic\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasPerByteEip7623\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasContractExecution\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasBufferEstimate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"gasBufferUnmeasured\",\"type\":\"uint256\"}],\"name\":\"ChangeGasParameters\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"UseNonce\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"auctioneer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"targetTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"searcherSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"auctioneerSig\",\"type\":\"bytes\"}],\"internalType\":\"structIAuctionEntryPoint.AuctionTx\",\"name\":\"auctionTx\",\"type\":\"tuple\"}],\"name\":\"call\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_auctioneer\",\"type\":\"address\"}],\"name\":\"changeAuctioneer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_depositVault\",\"type\":\"address\"}],\"name\":\"changeDepositVault\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_gasPerByteIntrinsic\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasPerByteEip7623\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasContractExecution\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasBufferEstimate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_gasBufferUnmeasured\",\"type\":\"uint256\"}],\"name\":\"changeGasParameters\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositVault\",\"outputs\":[{\"internalType\":\"contractIAuctionDepositVault\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasBufferEstimate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"targetTxHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bid\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"callGasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"searcherSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"auctioneerSig\",\"type\":\"bytes\"}],\"internalType\":\"structIAuctionEntryPoint.AuctionTx\",\"name\":\"auctionTx\",\"type\":\"tuple\"}],\"name\":\"getAuctionTxHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"searchers\",\"type\":\"address[]\"}],\"name\":\"getNoncesAndDeposits\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"nonces_\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"deposits_\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"5ec2c7bf": "auctioneer()",
		"ca157554": "call((bytes32,uint256,address,address,uint256,uint256,uint256,bytes,bytes,bytes))",
		"774f45ec": "changeAuctioneer(address)",
		"9d59928b": "changeDepositVault(address)",
		"2a215610": "changeGasParameters(uint256,uint256,uint256,uint256,uint256)",
		"d7cd3949": "depositVault()",
		"a5b2ab40": "gasBufferEstimate()",
		"a8aa9450": "getAuctionTxHash((bytes32,uint256,address,address,uint256,uint256,uint256,bytes,bytes,bytes))",
		"0339ed37": "getNoncesAndDeposits(address[])",
	},
}

// IAuctionEntryPointABI is the input ABI used to generate the binding from.
// Deprecated: Use IAuctionEntryPointMetaData.ABI instead.
var IAuctionEntryPointABI = IAuctionEntryPointMetaData.ABI

// IAuctionEntryPointBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IAuctionEntryPointBinRuntime = ``

// Deprecated: Use IAuctionEntryPointMetaData.Sigs instead.
// IAuctionEntryPointFuncSigs maps the 4-byte function signature to its string representation.
var IAuctionEntryPointFuncSigs = IAuctionEntryPointMetaData.Sigs

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

// IAuctionFeeVaultMetaData contains all meta data concerning the IAuctionFeeVault contract.
var IAuctionFeeVaultMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"paybackAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"validatorPaybackAmount\",\"type\":\"uint256\"}],\"name\":\"FeeDeposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FeeWithdrawal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reward\",\"type\":\"address\"}],\"name\":\"RewardAddressRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"searcherPaybackRate\",\"type\":\"uint256\"}],\"name\":\"SearcherPaybackRateUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"validatorPaybackRate\",\"type\":\"uint256\"}],\"name\":\"ValidatorPaybackRateUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"}],\"name\":\"getRewardAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddr\",\"type\":\"address\"}],\"name\":\"registerRewardAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_searcherPaybackRate\",\"type\":\"uint256\"}],\"name\":\"setSearcherPaybackRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_validatorPaybackRate\",\"type\":\"uint256\"}],\"name\":\"setValidatorPaybackRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"searcher\",\"type\":\"address\"}],\"name\":\"takeBid\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"27a50f72": "getRewardAddr(address)",
		"363d5183": "registerRewardAddress(address,address)",
		"36cf2c63": "setSearcherPaybackRate(uint256)",
		"11062696": "setValidatorPaybackRate(uint256)",
		"8573e2ff": "takeBid(address)",
		"51cff8d9": "withdraw(address)",
	},
}

// IAuctionFeeVaultABI is the input ABI used to generate the binding from.
// Deprecated: Use IAuctionFeeVaultMetaData.ABI instead.
var IAuctionFeeVaultABI = IAuctionFeeVaultMetaData.ABI

// IAuctionFeeVaultBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IAuctionFeeVaultBinRuntime = ``

// Deprecated: Use IAuctionFeeVaultMetaData.Sigs instead.
// IAuctionFeeVaultFuncSigs maps the 4-byte function signature to its string representation.
var IAuctionFeeVaultFuncSigs = IAuctionFeeVaultMetaData.Sigs

// IAuctionFeeVault is an auto generated Go binding around a Kaia contract.
type IAuctionFeeVault struct {
	IAuctionFeeVaultCaller     // Read-only binding to the contract
	IAuctionFeeVaultTransactor // Write-only binding to the contract
	IAuctionFeeVaultFilterer   // Log filterer for contract events
}

// IAuctionFeeVaultCaller is an auto generated read-only Go binding around a Kaia contract.
type IAuctionFeeVaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionFeeVaultTransactor is an auto generated write-only Go binding around a Kaia contract.
type IAuctionFeeVaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionFeeVaultFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IAuctionFeeVaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IAuctionFeeVaultSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IAuctionFeeVaultSession struct {
	Contract     *IAuctionFeeVault // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IAuctionFeeVaultCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IAuctionFeeVaultCallerSession struct {
	Contract *IAuctionFeeVaultCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// IAuctionFeeVaultTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IAuctionFeeVaultTransactorSession struct {
	Contract     *IAuctionFeeVaultTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IAuctionFeeVaultRaw is an auto generated low-level Go binding around a Kaia contract.
type IAuctionFeeVaultRaw struct {
	Contract *IAuctionFeeVault // Generic contract binding to access the raw methods on
}

// IAuctionFeeVaultCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IAuctionFeeVaultCallerRaw struct {
	Contract *IAuctionFeeVaultCaller // Generic read-only contract binding to access the raw methods on
}

// IAuctionFeeVaultTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IAuctionFeeVaultTransactorRaw struct {
	Contract *IAuctionFeeVaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIAuctionFeeVault creates a new instance of IAuctionFeeVault, bound to a specific deployed contract.
func NewIAuctionFeeVault(address common.Address, backend bind.ContractBackend) (*IAuctionFeeVault, error) {
	contract, err := bindIAuctionFeeVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVault{IAuctionFeeVaultCaller: IAuctionFeeVaultCaller{contract: contract}, IAuctionFeeVaultTransactor: IAuctionFeeVaultTransactor{contract: contract}, IAuctionFeeVaultFilterer: IAuctionFeeVaultFilterer{contract: contract}}, nil
}

// NewIAuctionFeeVaultCaller creates a new read-only instance of IAuctionFeeVault, bound to a specific deployed contract.
func NewIAuctionFeeVaultCaller(address common.Address, caller bind.ContractCaller) (*IAuctionFeeVaultCaller, error) {
	contract, err := bindIAuctionFeeVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVaultCaller{contract: contract}, nil
}

// NewIAuctionFeeVaultTransactor creates a new write-only instance of IAuctionFeeVault, bound to a specific deployed contract.
func NewIAuctionFeeVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*IAuctionFeeVaultTransactor, error) {
	contract, err := bindIAuctionFeeVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVaultTransactor{contract: contract}, nil
}

// NewIAuctionFeeVaultFilterer creates a new log filterer instance of IAuctionFeeVault, bound to a specific deployed contract.
func NewIAuctionFeeVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*IAuctionFeeVaultFilterer, error) {
	contract, err := bindIAuctionFeeVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVaultFilterer{contract: contract}, nil
}

// bindIAuctionFeeVault binds a generic wrapper to an already deployed contract.
func bindIAuctionFeeVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IAuctionFeeVaultMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAuctionFeeVault *IAuctionFeeVaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAuctionFeeVault.Contract.IAuctionFeeVaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAuctionFeeVault *IAuctionFeeVaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.IAuctionFeeVaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAuctionFeeVault *IAuctionFeeVaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.IAuctionFeeVaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IAuctionFeeVault *IAuctionFeeVaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IAuctionFeeVault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IAuctionFeeVault *IAuctionFeeVaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IAuctionFeeVault *IAuctionFeeVaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.contract.Transact(opts, method, params...)
}

// GetRewardAddr is a free data retrieval call binding the contract method 0x27a50f72.
//
// Solidity: function getRewardAddr(address nodeId) view returns(address)
func (_IAuctionFeeVault *IAuctionFeeVaultCaller) GetRewardAddr(opts *bind.CallOpts, nodeId common.Address) (common.Address, error) {
	var out []interface{}
	err := _IAuctionFeeVault.contract.Call(opts, &out, "getRewardAddr", nodeId)
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// GetRewardAddr is a free data retrieval call binding the contract method 0x27a50f72.
//
// Solidity: function getRewardAddr(address nodeId) view returns(address)
func (_IAuctionFeeVault *IAuctionFeeVaultSession) GetRewardAddr(nodeId common.Address) (common.Address, error) {
	return _IAuctionFeeVault.Contract.GetRewardAddr(&_IAuctionFeeVault.CallOpts, nodeId)
}

// GetRewardAddr is a free data retrieval call binding the contract method 0x27a50f72.
//
// Solidity: function getRewardAddr(address nodeId) view returns(address)
func (_IAuctionFeeVault *IAuctionFeeVaultCallerSession) GetRewardAddr(nodeId common.Address) (common.Address, error) {
	return _IAuctionFeeVault.Contract.GetRewardAddr(&_IAuctionFeeVault.CallOpts, nodeId)
}

// RegisterRewardAddress is a paid mutator transaction binding the contract method 0x363d5183.
//
// Solidity: function registerRewardAddress(address nodeId, address rewardAddr) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactor) RegisterRewardAddress(opts *bind.TransactOpts, nodeId common.Address, rewardAddr common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.contract.Transact(opts, "registerRewardAddress", nodeId, rewardAddr)
}

// RegisterRewardAddress is a paid mutator transaction binding the contract method 0x363d5183.
//
// Solidity: function registerRewardAddress(address nodeId, address rewardAddr) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultSession) RegisterRewardAddress(nodeId common.Address, rewardAddr common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.RegisterRewardAddress(&_IAuctionFeeVault.TransactOpts, nodeId, rewardAddr)
}

// RegisterRewardAddress is a paid mutator transaction binding the contract method 0x363d5183.
//
// Solidity: function registerRewardAddress(address nodeId, address rewardAddr) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactorSession) RegisterRewardAddress(nodeId common.Address, rewardAddr common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.RegisterRewardAddress(&_IAuctionFeeVault.TransactOpts, nodeId, rewardAddr)
}

// SetSearcherPaybackRate is a paid mutator transaction binding the contract method 0x36cf2c63.
//
// Solidity: function setSearcherPaybackRate(uint256 _searcherPaybackRate) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactor) SetSearcherPaybackRate(opts *bind.TransactOpts, _searcherPaybackRate *big.Int) (*types.Transaction, error) {
	return _IAuctionFeeVault.contract.Transact(opts, "setSearcherPaybackRate", _searcherPaybackRate)
}

// SetSearcherPaybackRate is a paid mutator transaction binding the contract method 0x36cf2c63.
//
// Solidity: function setSearcherPaybackRate(uint256 _searcherPaybackRate) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultSession) SetSearcherPaybackRate(_searcherPaybackRate *big.Int) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.SetSearcherPaybackRate(&_IAuctionFeeVault.TransactOpts, _searcherPaybackRate)
}

// SetSearcherPaybackRate is a paid mutator transaction binding the contract method 0x36cf2c63.
//
// Solidity: function setSearcherPaybackRate(uint256 _searcherPaybackRate) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactorSession) SetSearcherPaybackRate(_searcherPaybackRate *big.Int) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.SetSearcherPaybackRate(&_IAuctionFeeVault.TransactOpts, _searcherPaybackRate)
}

// SetValidatorPaybackRate is a paid mutator transaction binding the contract method 0x11062696.
//
// Solidity: function setValidatorPaybackRate(uint256 _validatorPaybackRate) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactor) SetValidatorPaybackRate(opts *bind.TransactOpts, _validatorPaybackRate *big.Int) (*types.Transaction, error) {
	return _IAuctionFeeVault.contract.Transact(opts, "setValidatorPaybackRate", _validatorPaybackRate)
}

// SetValidatorPaybackRate is a paid mutator transaction binding the contract method 0x11062696.
//
// Solidity: function setValidatorPaybackRate(uint256 _validatorPaybackRate) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultSession) SetValidatorPaybackRate(_validatorPaybackRate *big.Int) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.SetValidatorPaybackRate(&_IAuctionFeeVault.TransactOpts, _validatorPaybackRate)
}

// SetValidatorPaybackRate is a paid mutator transaction binding the contract method 0x11062696.
//
// Solidity: function setValidatorPaybackRate(uint256 _validatorPaybackRate) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactorSession) SetValidatorPaybackRate(_validatorPaybackRate *big.Int) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.SetValidatorPaybackRate(&_IAuctionFeeVault.TransactOpts, _validatorPaybackRate)
}

// TakeBid is a paid mutator transaction binding the contract method 0x8573e2ff.
//
// Solidity: function takeBid(address searcher) payable returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactor) TakeBid(opts *bind.TransactOpts, searcher common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.contract.Transact(opts, "takeBid", searcher)
}

// TakeBid is a paid mutator transaction binding the contract method 0x8573e2ff.
//
// Solidity: function takeBid(address searcher) payable returns()
func (_IAuctionFeeVault *IAuctionFeeVaultSession) TakeBid(searcher common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.TakeBid(&_IAuctionFeeVault.TransactOpts, searcher)
}

// TakeBid is a paid mutator transaction binding the contract method 0x8573e2ff.
//
// Solidity: function takeBid(address searcher) payable returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactorSession) TakeBid(searcher common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.TakeBid(&_IAuctionFeeVault.TransactOpts, searcher)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address to) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactor) Withdraw(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.contract.Transact(opts, "withdraw", to)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address to) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultSession) Withdraw(to common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.Withdraw(&_IAuctionFeeVault.TransactOpts, to)
}

// Withdraw is a paid mutator transaction binding the contract method 0x51cff8d9.
//
// Solidity: function withdraw(address to) returns()
func (_IAuctionFeeVault *IAuctionFeeVaultTransactorSession) Withdraw(to common.Address) (*types.Transaction, error) {
	return _IAuctionFeeVault.Contract.Withdraw(&_IAuctionFeeVault.TransactOpts, to)
}

// IAuctionFeeVaultFeeDepositIterator is returned from FilterFeeDeposit and is used to iterate over the raw logs and unpacked data for FeeDeposit events raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultFeeDepositIterator struct {
	Event *IAuctionFeeVaultFeeDeposit // Event containing the contract specifics and raw log

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
func (it *IAuctionFeeVaultFeeDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionFeeVaultFeeDeposit)
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
		it.Event = new(IAuctionFeeVaultFeeDeposit)
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
func (it *IAuctionFeeVaultFeeDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionFeeVaultFeeDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionFeeVaultFeeDeposit represents a FeeDeposit event raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultFeeDeposit struct {
	Sender                 common.Address
	Amount                 *big.Int
	PaybackAmount          *big.Int
	ValidatorPaybackAmount *big.Int
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterFeeDeposit is a free log retrieval operation binding the contract event 0xa34c9ef6ada915fef21639b2d5c085580cf79046cca66c2c2e8b87e2f3bd8567.
//
// Solidity: event FeeDeposit(address indexed sender, uint256 amount, uint256 paybackAmount, uint256 validatorPaybackAmount)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) FilterFeeDeposit(opts *bind.FilterOpts, sender []common.Address) (*IAuctionFeeVaultFeeDepositIterator, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IAuctionFeeVault.contract.FilterLogs(opts, "FeeDeposit", senderRule)
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVaultFeeDepositIterator{contract: _IAuctionFeeVault.contract, event: "FeeDeposit", logs: logs, sub: sub}, nil
}

// WatchFeeDeposit is a free log subscription operation binding the contract event 0xa34c9ef6ada915fef21639b2d5c085580cf79046cca66c2c2e8b87e2f3bd8567.
//
// Solidity: event FeeDeposit(address indexed sender, uint256 amount, uint256 paybackAmount, uint256 validatorPaybackAmount)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) WatchFeeDeposit(opts *bind.WatchOpts, sink chan<- *IAuctionFeeVaultFeeDeposit, sender []common.Address) (event.Subscription, error) {
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _IAuctionFeeVault.contract.WatchLogs(opts, "FeeDeposit", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionFeeVaultFeeDeposit)
				if err := _IAuctionFeeVault.contract.UnpackLog(event, "FeeDeposit", log); err != nil {
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

// ParseFeeDeposit is a log parse operation binding the contract event 0xa34c9ef6ada915fef21639b2d5c085580cf79046cca66c2c2e8b87e2f3bd8567.
//
// Solidity: event FeeDeposit(address indexed sender, uint256 amount, uint256 paybackAmount, uint256 validatorPaybackAmount)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) ParseFeeDeposit(log types.Log) (*IAuctionFeeVaultFeeDeposit, error) {
	event := new(IAuctionFeeVaultFeeDeposit)
	if err := _IAuctionFeeVault.contract.UnpackLog(event, "FeeDeposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionFeeVaultFeeWithdrawalIterator is returned from FilterFeeWithdrawal and is used to iterate over the raw logs and unpacked data for FeeWithdrawal events raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultFeeWithdrawalIterator struct {
	Event *IAuctionFeeVaultFeeWithdrawal // Event containing the contract specifics and raw log

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
func (it *IAuctionFeeVaultFeeWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionFeeVaultFeeWithdrawal)
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
		it.Event = new(IAuctionFeeVaultFeeWithdrawal)
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
func (it *IAuctionFeeVaultFeeWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionFeeVaultFeeWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionFeeVaultFeeWithdrawal represents a FeeWithdrawal event raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultFeeWithdrawal struct {
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterFeeWithdrawal is a free log retrieval operation binding the contract event 0x706d7f48c702007c2fb0881cea5759732e64f52faee427d5ab030787cfb7d787.
//
// Solidity: event FeeWithdrawal(uint256 amount)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) FilterFeeWithdrawal(opts *bind.FilterOpts) (*IAuctionFeeVaultFeeWithdrawalIterator, error) {
	logs, sub, err := _IAuctionFeeVault.contract.FilterLogs(opts, "FeeWithdrawal")
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVaultFeeWithdrawalIterator{contract: _IAuctionFeeVault.contract, event: "FeeWithdrawal", logs: logs, sub: sub}, nil
}

// WatchFeeWithdrawal is a free log subscription operation binding the contract event 0x706d7f48c702007c2fb0881cea5759732e64f52faee427d5ab030787cfb7d787.
//
// Solidity: event FeeWithdrawal(uint256 amount)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) WatchFeeWithdrawal(opts *bind.WatchOpts, sink chan<- *IAuctionFeeVaultFeeWithdrawal) (event.Subscription, error) {
	logs, sub, err := _IAuctionFeeVault.contract.WatchLogs(opts, "FeeWithdrawal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionFeeVaultFeeWithdrawal)
				if err := _IAuctionFeeVault.contract.UnpackLog(event, "FeeWithdrawal", log); err != nil {
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

// ParseFeeWithdrawal is a log parse operation binding the contract event 0x706d7f48c702007c2fb0881cea5759732e64f52faee427d5ab030787cfb7d787.
//
// Solidity: event FeeWithdrawal(uint256 amount)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) ParseFeeWithdrawal(log types.Log) (*IAuctionFeeVaultFeeWithdrawal, error) {
	event := new(IAuctionFeeVaultFeeWithdrawal)
	if err := _IAuctionFeeVault.contract.UnpackLog(event, "FeeWithdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionFeeVaultRewardAddressRegisteredIterator is returned from FilterRewardAddressRegistered and is used to iterate over the raw logs and unpacked data for RewardAddressRegistered events raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultRewardAddressRegisteredIterator struct {
	Event *IAuctionFeeVaultRewardAddressRegistered // Event containing the contract specifics and raw log

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
func (it *IAuctionFeeVaultRewardAddressRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionFeeVaultRewardAddressRegistered)
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
		it.Event = new(IAuctionFeeVaultRewardAddressRegistered)
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
func (it *IAuctionFeeVaultRewardAddressRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionFeeVaultRewardAddressRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionFeeVaultRewardAddressRegistered represents a RewardAddressRegistered event raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultRewardAddressRegistered struct {
	NodeId common.Address
	Reward common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardAddressRegistered is a free log retrieval operation binding the contract event 0xe608476cb01b1d04f944f0fdb25841b1f483d26965d42c4a1fab67b8b1488b3b.
//
// Solidity: event RewardAddressRegistered(address indexed nodeId, address indexed reward)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) FilterRewardAddressRegistered(opts *bind.FilterOpts, nodeId []common.Address, reward []common.Address) (*IAuctionFeeVaultRewardAddressRegisteredIterator, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var rewardRule []interface{}
	for _, rewardItem := range reward {
		rewardRule = append(rewardRule, rewardItem)
	}

	logs, sub, err := _IAuctionFeeVault.contract.FilterLogs(opts, "RewardAddressRegistered", nodeIdRule, rewardRule)
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVaultRewardAddressRegisteredIterator{contract: _IAuctionFeeVault.contract, event: "RewardAddressRegistered", logs: logs, sub: sub}, nil
}

// WatchRewardAddressRegistered is a free log subscription operation binding the contract event 0xe608476cb01b1d04f944f0fdb25841b1f483d26965d42c4a1fab67b8b1488b3b.
//
// Solidity: event RewardAddressRegistered(address indexed nodeId, address indexed reward)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) WatchRewardAddressRegistered(opts *bind.WatchOpts, sink chan<- *IAuctionFeeVaultRewardAddressRegistered, nodeId []common.Address, reward []common.Address) (event.Subscription, error) {
	var nodeIdRule []interface{}
	for _, nodeIdItem := range nodeId {
		nodeIdRule = append(nodeIdRule, nodeIdItem)
	}
	var rewardRule []interface{}
	for _, rewardItem := range reward {
		rewardRule = append(rewardRule, rewardItem)
	}

	logs, sub, err := _IAuctionFeeVault.contract.WatchLogs(opts, "RewardAddressRegistered", nodeIdRule, rewardRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionFeeVaultRewardAddressRegistered)
				if err := _IAuctionFeeVault.contract.UnpackLog(event, "RewardAddressRegistered", log); err != nil {
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

// ParseRewardAddressRegistered is a log parse operation binding the contract event 0xe608476cb01b1d04f944f0fdb25841b1f483d26965d42c4a1fab67b8b1488b3b.
//
// Solidity: event RewardAddressRegistered(address indexed nodeId, address indexed reward)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) ParseRewardAddressRegistered(log types.Log) (*IAuctionFeeVaultRewardAddressRegistered, error) {
	event := new(IAuctionFeeVaultRewardAddressRegistered)
	if err := _IAuctionFeeVault.contract.UnpackLog(event, "RewardAddressRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionFeeVaultSearcherPaybackRateUpdatedIterator is returned from FilterSearcherPaybackRateUpdated and is used to iterate over the raw logs and unpacked data for SearcherPaybackRateUpdated events raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultSearcherPaybackRateUpdatedIterator struct {
	Event *IAuctionFeeVaultSearcherPaybackRateUpdated // Event containing the contract specifics and raw log

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
func (it *IAuctionFeeVaultSearcherPaybackRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionFeeVaultSearcherPaybackRateUpdated)
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
		it.Event = new(IAuctionFeeVaultSearcherPaybackRateUpdated)
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
func (it *IAuctionFeeVaultSearcherPaybackRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionFeeVaultSearcherPaybackRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionFeeVaultSearcherPaybackRateUpdated represents a SearcherPaybackRateUpdated event raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultSearcherPaybackRateUpdated struct {
	SearcherPaybackRate *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSearcherPaybackRateUpdated is a free log retrieval operation binding the contract event 0x71745430318b073bd776904f2432cb283ce3d2ded537bafe2640cf4d6e4bc64f.
//
// Solidity: event SearcherPaybackRateUpdated(uint256 searcherPaybackRate)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) FilterSearcherPaybackRateUpdated(opts *bind.FilterOpts) (*IAuctionFeeVaultSearcherPaybackRateUpdatedIterator, error) {
	logs, sub, err := _IAuctionFeeVault.contract.FilterLogs(opts, "SearcherPaybackRateUpdated")
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVaultSearcherPaybackRateUpdatedIterator{contract: _IAuctionFeeVault.contract, event: "SearcherPaybackRateUpdated", logs: logs, sub: sub}, nil
}

// WatchSearcherPaybackRateUpdated is a free log subscription operation binding the contract event 0x71745430318b073bd776904f2432cb283ce3d2ded537bafe2640cf4d6e4bc64f.
//
// Solidity: event SearcherPaybackRateUpdated(uint256 searcherPaybackRate)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) WatchSearcherPaybackRateUpdated(opts *bind.WatchOpts, sink chan<- *IAuctionFeeVaultSearcherPaybackRateUpdated) (event.Subscription, error) {
	logs, sub, err := _IAuctionFeeVault.contract.WatchLogs(opts, "SearcherPaybackRateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionFeeVaultSearcherPaybackRateUpdated)
				if err := _IAuctionFeeVault.contract.UnpackLog(event, "SearcherPaybackRateUpdated", log); err != nil {
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

// ParseSearcherPaybackRateUpdated is a log parse operation binding the contract event 0x71745430318b073bd776904f2432cb283ce3d2ded537bafe2640cf4d6e4bc64f.
//
// Solidity: event SearcherPaybackRateUpdated(uint256 searcherPaybackRate)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) ParseSearcherPaybackRateUpdated(log types.Log) (*IAuctionFeeVaultSearcherPaybackRateUpdated, error) {
	event := new(IAuctionFeeVaultSearcherPaybackRateUpdated)
	if err := _IAuctionFeeVault.contract.UnpackLog(event, "SearcherPaybackRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IAuctionFeeVaultValidatorPaybackRateUpdatedIterator is returned from FilterValidatorPaybackRateUpdated and is used to iterate over the raw logs and unpacked data for ValidatorPaybackRateUpdated events raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultValidatorPaybackRateUpdatedIterator struct {
	Event *IAuctionFeeVaultValidatorPaybackRateUpdated // Event containing the contract specifics and raw log

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
func (it *IAuctionFeeVaultValidatorPaybackRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IAuctionFeeVaultValidatorPaybackRateUpdated)
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
		it.Event = new(IAuctionFeeVaultValidatorPaybackRateUpdated)
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
func (it *IAuctionFeeVaultValidatorPaybackRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IAuctionFeeVaultValidatorPaybackRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IAuctionFeeVaultValidatorPaybackRateUpdated represents a ValidatorPaybackRateUpdated event raised by the IAuctionFeeVault contract.
type IAuctionFeeVaultValidatorPaybackRateUpdated struct {
	ValidatorPaybackRate *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterValidatorPaybackRateUpdated is a free log retrieval operation binding the contract event 0x5309d48fe743a67ce32d8f66af9e2388d65bfc8cc026a4e1fbed3a4612a0af98.
//
// Solidity: event ValidatorPaybackRateUpdated(uint256 validatorPaybackRate)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) FilterValidatorPaybackRateUpdated(opts *bind.FilterOpts) (*IAuctionFeeVaultValidatorPaybackRateUpdatedIterator, error) {
	logs, sub, err := _IAuctionFeeVault.contract.FilterLogs(opts, "ValidatorPaybackRateUpdated")
	if err != nil {
		return nil, err
	}
	return &IAuctionFeeVaultValidatorPaybackRateUpdatedIterator{contract: _IAuctionFeeVault.contract, event: "ValidatorPaybackRateUpdated", logs: logs, sub: sub}, nil
}

// WatchValidatorPaybackRateUpdated is a free log subscription operation binding the contract event 0x5309d48fe743a67ce32d8f66af9e2388d65bfc8cc026a4e1fbed3a4612a0af98.
//
// Solidity: event ValidatorPaybackRateUpdated(uint256 validatorPaybackRate)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) WatchValidatorPaybackRateUpdated(opts *bind.WatchOpts, sink chan<- *IAuctionFeeVaultValidatorPaybackRateUpdated) (event.Subscription, error) {
	logs, sub, err := _IAuctionFeeVault.contract.WatchLogs(opts, "ValidatorPaybackRateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IAuctionFeeVaultValidatorPaybackRateUpdated)
				if err := _IAuctionFeeVault.contract.UnpackLog(event, "ValidatorPaybackRateUpdated", log); err != nil {
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

// ParseValidatorPaybackRateUpdated is a log parse operation binding the contract event 0x5309d48fe743a67ce32d8f66af9e2388d65bfc8cc026a4e1fbed3a4612a0af98.
//
// Solidity: event ValidatorPaybackRateUpdated(uint256 validatorPaybackRate)
func (_IAuctionFeeVault *IAuctionFeeVaultFilterer) ParseValidatorPaybackRateUpdated(log types.Log) (*IAuctionFeeVaultValidatorPaybackRateUpdated, error) {
	event := new(IAuctionFeeVaultValidatorPaybackRateUpdated)
	if err := _IAuctionFeeVault.contract.UnpackLog(event, "ValidatorPaybackRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OwnableMetaData contains all meta data concerning the Ownable contract.
var OwnableMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"8da5cb5b": "owner()",
		"715018a6": "renounceOwnership()",
		"f2fde38b": "transferOwnership(address)",
	},
}

// OwnableABI is the input ABI used to generate the binding from.
// Deprecated: Use OwnableMetaData.ABI instead.
var OwnableABI = OwnableMetaData.ABI

// OwnableBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const OwnableBinRuntime = ``

// Deprecated: Use OwnableMetaData.Sigs instead.
// OwnableFuncSigs maps the 4-byte function signature to its string representation.
var OwnableFuncSigs = OwnableMetaData.Sigs

// Ownable is an auto generated Go binding around a Kaia contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around a Kaia contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around a Kaia contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around a Kaia contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
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

	logs chan types.Log    // Log channel receiving the found contract events
	sub  kaia.Subscription // Subscription for errors, completion and termination
	done bool              // Whether the subscription completed delivering logs
	fail error             // Occurred error to stop iteration
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
	event.Raw = log
	return event, nil
}
