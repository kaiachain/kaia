// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package valset

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

// IValidatorStateSystemValidatorUpdateRequest is an auto generated low-level Go binding around an user-defined struct.
type IValidatorStateSystemValidatorUpdateRequest struct {
	Addr  common.Address
	State uint8
}

// IValidatorStateValidatorState is an auto generated low-level Go binding around an user-defined struct.
type IValidatorStateValidatorState struct {
	Addr          common.Address
	State         uint8
	IdleTimeout   *big.Int
	PausedTimeout *big.Int
}

// EnumerableSetMetaData contains all meta data concerning the EnumerableSet contract.
var EnumerableSetMetaData = &bind.MetaData{
	ABI: "[]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220ab7ec786af90a9f4a47555f20c6e9ec9c8cbd49047b90c75fa05c2991f6f865664736f6c63430008190033",
}

// EnumerableSetABI is the input ABI used to generate the binding from.
// Deprecated: Use EnumerableSetMetaData.ABI instead.
var EnumerableSetABI = EnumerableSetMetaData.ABI

// EnumerableSetBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const EnumerableSetBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220ab7ec786af90a9f4a47555f20c6e9ec9c8cbd49047b90c75fa05c2991f6f865664736f6c63430008190033`

// EnumerableSetBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EnumerableSetMetaData.Bin instead.
var EnumerableSetBin = EnumerableSetMetaData.Bin

// DeployEnumerableSet deploys a new Kaia contract, binding an instance of EnumerableSet to it.
func DeployEnumerableSet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *EnumerableSet, error) {
	parsed, err := EnumerableSetMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EnumerableSetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &EnumerableSet{EnumerableSetCaller: EnumerableSetCaller{contract: contract}, EnumerableSetTransactor: EnumerableSetTransactor{contract: contract}, EnumerableSetFilterer: EnumerableSetFilterer{contract: contract}}, nil
}

// EnumerableSet is an auto generated Go binding around a Kaia contract.
type EnumerableSet struct {
	EnumerableSetCaller     // Read-only binding to the contract
	EnumerableSetTransactor // Write-only binding to the contract
	EnumerableSetFilterer   // Log filterer for contract events
}

// EnumerableSetCaller is an auto generated read-only Go binding around a Kaia contract.
type EnumerableSetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetTransactor is an auto generated write-only Go binding around a Kaia contract.
type EnumerableSetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type EnumerableSetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EnumerableSetSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type EnumerableSetSession struct {
	Contract     *EnumerableSet    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EnumerableSetCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type EnumerableSetCallerSession struct {
	Contract *EnumerableSetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// EnumerableSetTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type EnumerableSetTransactorSession struct {
	Contract     *EnumerableSetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// EnumerableSetRaw is an auto generated low-level Go binding around a Kaia contract.
type EnumerableSetRaw struct {
	Contract *EnumerableSet // Generic contract binding to access the raw methods on
}

// EnumerableSetCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type EnumerableSetCallerRaw struct {
	Contract *EnumerableSetCaller // Generic read-only contract binding to access the raw methods on
}

// EnumerableSetTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type EnumerableSetTransactorRaw struct {
	Contract *EnumerableSetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEnumerableSet creates a new instance of EnumerableSet, bound to a specific deployed contract.
func NewEnumerableSet(address common.Address, backend bind.ContractBackend) (*EnumerableSet, error) {
	contract, err := bindEnumerableSet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EnumerableSet{EnumerableSetCaller: EnumerableSetCaller{contract: contract}, EnumerableSetTransactor: EnumerableSetTransactor{contract: contract}, EnumerableSetFilterer: EnumerableSetFilterer{contract: contract}}, nil
}

// NewEnumerableSetCaller creates a new read-only instance of EnumerableSet, bound to a specific deployed contract.
func NewEnumerableSetCaller(address common.Address, caller bind.ContractCaller) (*EnumerableSetCaller, error) {
	contract, err := bindEnumerableSet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetCaller{contract: contract}, nil
}

// NewEnumerableSetTransactor creates a new write-only instance of EnumerableSet, bound to a specific deployed contract.
func NewEnumerableSetTransactor(address common.Address, transactor bind.ContractTransactor) (*EnumerableSetTransactor, error) {
	contract, err := bindEnumerableSet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetTransactor{contract: contract}, nil
}

// NewEnumerableSetFilterer creates a new log filterer instance of EnumerableSet, bound to a specific deployed contract.
func NewEnumerableSetFilterer(address common.Address, filterer bind.ContractFilterer) (*EnumerableSetFilterer, error) {
	contract, err := bindEnumerableSet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EnumerableSetFilterer{contract: contract}, nil
}

// bindEnumerableSet binds a generic wrapper to an already deployed contract.
func bindEnumerableSet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EnumerableSetMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnumerableSet *EnumerableSetRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnumerableSet.Contract.EnumerableSetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnumerableSet *EnumerableSetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnumerableSet.Contract.EnumerableSetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnumerableSet *EnumerableSetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnumerableSet.Contract.EnumerableSetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EnumerableSet *EnumerableSetCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EnumerableSet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EnumerableSet *EnumerableSetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EnumerableSet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EnumerableSet *EnumerableSetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EnumerableSet.Contract.contract.Transact(opts, method, params...)
}

// IAddressBookMetaData contains all meta data concerning the IAddressBook contract.
var IAddressBookMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_cnNodeId\",\"type\":\"address\"}],\"name\":\"getCnInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"cnNodeId\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"cnStakingcontract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"cnRewardAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"15575d5a": "getCnInfo(address)",
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

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookCaller) GetCnInfo(opts *bind.CallOpts, _cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error) {
	var out []interface{}
	err := _IAddressBook.contract.Call(opts, &out, "getCnInfo", _cnNodeId)

	outstruct := new(struct {
		CnNodeId          common.Address
		CnStakingcontract common.Address
		CnRewardAddress   common.Address
	})
	if err != nil {
		return *outstruct, err
	}

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
}, error) {
	return _IAddressBook.Contract.GetCnInfo(&_IAddressBook.CallOpts, _cnNodeId)
}

// GetCnInfo is a free data retrieval call binding the contract method 0x15575d5a.
//
// Solidity: function getCnInfo(address _cnNodeId) view returns(address cnNodeId, address cnStakingcontract, address cnRewardAddress)
func (_IAddressBook *IAddressBookCallerSession) GetCnInfo(_cnNodeId common.Address) (struct {
	CnNodeId          common.Address
	CnStakingcontract common.Address
	CnRewardAddress   common.Address
}, error) {
	return _IAddressBook.Contract.GetCnInfo(&_IAddressBook.CallOpts, _cnNodeId)
}

// ICnStakingMetaData contains all meta data concerning the ICnStaking contract.
var ICnStakingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"staking\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unstaking\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
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

// IValidatorStateMetaData contains all meta data concerning the IValidatorState contract.
var IValidatorStateMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"SetIdleTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"SetPausedTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumIValidatorState.State\",\"name\":\"oldStatw\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumIValidatorState.State\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"StateTranstiion\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getAllValidators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"enumIValidatorState.State\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"idleTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pausedTimeout\",\"type\":\"uint256\"}],\"internalType\":\"structIValidatorState.ValidatorState[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"enumIValidatorState.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structIValidatorState.SystemValidatorUpdateRequest[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"name\":\"setValidatorStates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"f3513a37": "getAllValidators()",
		"513e846c": "setValidatorStates((address,uint8)[])",
	},
}

// IValidatorStateABI is the input ABI used to generate the binding from.
// Deprecated: Use IValidatorStateMetaData.ABI instead.
var IValidatorStateABI = IValidatorStateMetaData.ABI

// IValidatorStateBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IValidatorStateBinRuntime = ``

// Deprecated: Use IValidatorStateMetaData.Sigs instead.
// IValidatorStateFuncSigs maps the 4-byte function signature to its string representation.
var IValidatorStateFuncSigs = IValidatorStateMetaData.Sigs

// IValidatorState is an auto generated Go binding around a Kaia contract.
type IValidatorState struct {
	IValidatorStateCaller     // Read-only binding to the contract
	IValidatorStateTransactor // Write-only binding to the contract
	IValidatorStateFilterer   // Log filterer for contract events
}

// IValidatorStateCaller is an auto generated read-only Go binding around a Kaia contract.
type IValidatorStateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IValidatorStateTransactor is an auto generated write-only Go binding around a Kaia contract.
type IValidatorStateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IValidatorStateFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IValidatorStateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IValidatorStateSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IValidatorStateSession struct {
	Contract     *IValidatorState  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IValidatorStateCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IValidatorStateCallerSession struct {
	Contract *IValidatorStateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// IValidatorStateTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IValidatorStateTransactorSession struct {
	Contract     *IValidatorStateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// IValidatorStateRaw is an auto generated low-level Go binding around a Kaia contract.
type IValidatorStateRaw struct {
	Contract *IValidatorState // Generic contract binding to access the raw methods on
}

// IValidatorStateCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IValidatorStateCallerRaw struct {
	Contract *IValidatorStateCaller // Generic read-only contract binding to access the raw methods on
}

// IValidatorStateTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IValidatorStateTransactorRaw struct {
	Contract *IValidatorStateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIValidatorState creates a new instance of IValidatorState, bound to a specific deployed contract.
func NewIValidatorState(address common.Address, backend bind.ContractBackend) (*IValidatorState, error) {
	contract, err := bindIValidatorState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IValidatorState{IValidatorStateCaller: IValidatorStateCaller{contract: contract}, IValidatorStateTransactor: IValidatorStateTransactor{contract: contract}, IValidatorStateFilterer: IValidatorStateFilterer{contract: contract}}, nil
}

// NewIValidatorStateCaller creates a new read-only instance of IValidatorState, bound to a specific deployed contract.
func NewIValidatorStateCaller(address common.Address, caller bind.ContractCaller) (*IValidatorStateCaller, error) {
	contract, err := bindIValidatorState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IValidatorStateCaller{contract: contract}, nil
}

// NewIValidatorStateTransactor creates a new write-only instance of IValidatorState, bound to a specific deployed contract.
func NewIValidatorStateTransactor(address common.Address, transactor bind.ContractTransactor) (*IValidatorStateTransactor, error) {
	contract, err := bindIValidatorState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IValidatorStateTransactor{contract: contract}, nil
}

// NewIValidatorStateFilterer creates a new log filterer instance of IValidatorState, bound to a specific deployed contract.
func NewIValidatorStateFilterer(address common.Address, filterer bind.ContractFilterer) (*IValidatorStateFilterer, error) {
	contract, err := bindIValidatorState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IValidatorStateFilterer{contract: contract}, nil
}

// bindIValidatorState binds a generic wrapper to an already deployed contract.
func bindIValidatorState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IValidatorStateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IValidatorState *IValidatorStateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IValidatorState.Contract.IValidatorStateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IValidatorState *IValidatorStateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IValidatorState.Contract.IValidatorStateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IValidatorState *IValidatorStateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IValidatorState.Contract.IValidatorStateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IValidatorState *IValidatorStateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IValidatorState.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IValidatorState *IValidatorStateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IValidatorState.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IValidatorState *IValidatorStateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IValidatorState.Contract.contract.Transact(opts, method, params...)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((address,uint8,uint256,uint256)[])
func (_IValidatorState *IValidatorStateCaller) GetAllValidators(opts *bind.CallOpts) ([]IValidatorStateValidatorState, error) {
	var out []interface{}
	err := _IValidatorState.contract.Call(opts, &out, "getAllValidators")

	if err != nil {
		return *new([]IValidatorStateValidatorState), err
	}

	out0 := *abi.ConvertType(out[0], new([]IValidatorStateValidatorState)).(*[]IValidatorStateValidatorState)

	return out0, err

}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((address,uint8,uint256,uint256)[])
func (_IValidatorState *IValidatorStateSession) GetAllValidators() ([]IValidatorStateValidatorState, error) {
	return _IValidatorState.Contract.GetAllValidators(&_IValidatorState.CallOpts)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((address,uint8,uint256,uint256)[])
func (_IValidatorState *IValidatorStateCallerSession) GetAllValidators() ([]IValidatorStateValidatorState, error) {
	return _IValidatorState.Contract.GetAllValidators(&_IValidatorState.CallOpts)
}

// SetValidatorStates is a paid mutator transaction binding the contract method 0x513e846c.
//
// Solidity: function setValidatorStates((address,uint8)[] ) returns()
func (_IValidatorState *IValidatorStateTransactor) SetValidatorStates(opts *bind.TransactOpts, arg0 []IValidatorStateSystemValidatorUpdateRequest) (*types.Transaction, error) {
	return _IValidatorState.contract.Transact(opts, "setValidatorStates", arg0)
}

// SetValidatorStates is a paid mutator transaction binding the contract method 0x513e846c.
//
// Solidity: function setValidatorStates((address,uint8)[] ) returns()
func (_IValidatorState *IValidatorStateSession) SetValidatorStates(arg0 []IValidatorStateSystemValidatorUpdateRequest) (*types.Transaction, error) {
	return _IValidatorState.Contract.SetValidatorStates(&_IValidatorState.TransactOpts, arg0)
}

// SetValidatorStates is a paid mutator transaction binding the contract method 0x513e846c.
//
// Solidity: function setValidatorStates((address,uint8)[] ) returns()
func (_IValidatorState *IValidatorStateTransactorSession) SetValidatorStates(arg0 []IValidatorStateSystemValidatorUpdateRequest) (*types.Transaction, error) {
	return _IValidatorState.Contract.SetValidatorStates(&_IValidatorState.TransactOpts, arg0)
}

// IValidatorStateSetIdleTimeoutIterator is returned from FilterSetIdleTimeout and is used to iterate over the raw logs and unpacked data for SetIdleTimeout events raised by the IValidatorState contract.
type IValidatorStateSetIdleTimeoutIterator struct {
	Event *IValidatorStateSetIdleTimeout // Event containing the contract specifics and raw log

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
func (it *IValidatorStateSetIdleTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorStateSetIdleTimeout)
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
		it.Event = new(IValidatorStateSetIdleTimeout)
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
func (it *IValidatorStateSetIdleTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorStateSetIdleTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorStateSetIdleTimeout represents a SetIdleTimeout event raised by the IValidatorState contract.
type IValidatorStateSetIdleTimeout struct {
	Addr      common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetIdleTimeout is a free log retrieval operation binding the contract event 0xc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b818.
//
// Solidity: event SetIdleTimeout(address addr, uint256 timestamp)
func (_IValidatorState *IValidatorStateFilterer) FilterSetIdleTimeout(opts *bind.FilterOpts) (*IValidatorStateSetIdleTimeoutIterator, error) {

	logs, sub, err := _IValidatorState.contract.FilterLogs(opts, "SetIdleTimeout")
	if err != nil {
		return nil, err
	}
	return &IValidatorStateSetIdleTimeoutIterator{contract: _IValidatorState.contract, event: "SetIdleTimeout", logs: logs, sub: sub}, nil
}

// WatchSetIdleTimeout is a free log subscription operation binding the contract event 0xc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b818.
//
// Solidity: event SetIdleTimeout(address addr, uint256 timestamp)
func (_IValidatorState *IValidatorStateFilterer) WatchSetIdleTimeout(opts *bind.WatchOpts, sink chan<- *IValidatorStateSetIdleTimeout) (event.Subscription, error) {

	logs, sub, err := _IValidatorState.contract.WatchLogs(opts, "SetIdleTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorStateSetIdleTimeout)
				if err := _IValidatorState.contract.UnpackLog(event, "SetIdleTimeout", log); err != nil {
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

// ParseSetIdleTimeout is a log parse operation binding the contract event 0xc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b818.
//
// Solidity: event SetIdleTimeout(address addr, uint256 timestamp)
func (_IValidatorState *IValidatorStateFilterer) ParseSetIdleTimeout(log types.Log) (*IValidatorStateSetIdleTimeout, error) {
	event := new(IValidatorStateSetIdleTimeout)
	if err := _IValidatorState.contract.UnpackLog(event, "SetIdleTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorStateSetPausedTimeoutIterator is returned from FilterSetPausedTimeout and is used to iterate over the raw logs and unpacked data for SetPausedTimeout events raised by the IValidatorState contract.
type IValidatorStateSetPausedTimeoutIterator struct {
	Event *IValidatorStateSetPausedTimeout // Event containing the contract specifics and raw log

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
func (it *IValidatorStateSetPausedTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorStateSetPausedTimeout)
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
		it.Event = new(IValidatorStateSetPausedTimeout)
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
func (it *IValidatorStateSetPausedTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorStateSetPausedTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorStateSetPausedTimeout represents a SetPausedTimeout event raised by the IValidatorState contract.
type IValidatorStateSetPausedTimeout struct {
	Addr      common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetPausedTimeout is a free log retrieval operation binding the contract event 0xbc1a6d4bb677fb39bdb63d1713dc004bd314de87091c5dfbba6a8d78408d0dc9.
//
// Solidity: event SetPausedTimeout(address addr, uint256 timestamp)
func (_IValidatorState *IValidatorStateFilterer) FilterSetPausedTimeout(opts *bind.FilterOpts) (*IValidatorStateSetPausedTimeoutIterator, error) {

	logs, sub, err := _IValidatorState.contract.FilterLogs(opts, "SetPausedTimeout")
	if err != nil {
		return nil, err
	}
	return &IValidatorStateSetPausedTimeoutIterator{contract: _IValidatorState.contract, event: "SetPausedTimeout", logs: logs, sub: sub}, nil
}

// WatchSetPausedTimeout is a free log subscription operation binding the contract event 0xbc1a6d4bb677fb39bdb63d1713dc004bd314de87091c5dfbba6a8d78408d0dc9.
//
// Solidity: event SetPausedTimeout(address addr, uint256 timestamp)
func (_IValidatorState *IValidatorStateFilterer) WatchSetPausedTimeout(opts *bind.WatchOpts, sink chan<- *IValidatorStateSetPausedTimeout) (event.Subscription, error) {

	logs, sub, err := _IValidatorState.contract.WatchLogs(opts, "SetPausedTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorStateSetPausedTimeout)
				if err := _IValidatorState.contract.UnpackLog(event, "SetPausedTimeout", log); err != nil {
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

// ParseSetPausedTimeout is a log parse operation binding the contract event 0xbc1a6d4bb677fb39bdb63d1713dc004bd314de87091c5dfbba6a8d78408d0dc9.
//
// Solidity: event SetPausedTimeout(address addr, uint256 timestamp)
func (_IValidatorState *IValidatorStateFilterer) ParseSetPausedTimeout(log types.Log) (*IValidatorStateSetPausedTimeout, error) {
	event := new(IValidatorStateSetPausedTimeout)
	if err := _IValidatorState.contract.UnpackLog(event, "SetPausedTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// IValidatorStateStateTranstiionIterator is returned from FilterStateTranstiion and is used to iterate over the raw logs and unpacked data for StateTranstiion events raised by the IValidatorState contract.
type IValidatorStateStateTranstiionIterator struct {
	Event *IValidatorStateStateTranstiion // Event containing the contract specifics and raw log

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
func (it *IValidatorStateStateTranstiionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IValidatorStateStateTranstiion)
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
		it.Event = new(IValidatorStateStateTranstiion)
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
func (it *IValidatorStateStateTranstiionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IValidatorStateStateTranstiionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IValidatorStateStateTranstiion represents a StateTranstiion event raised by the IValidatorState contract.
type IValidatorStateStateTranstiion struct {
	OldStatw uint8
	NewState uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStateTranstiion is a free log retrieval operation binding the contract event 0x8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee36091.
//
// Solidity: event StateTranstiion(uint8 oldStatw, uint8 newState)
func (_IValidatorState *IValidatorStateFilterer) FilterStateTranstiion(opts *bind.FilterOpts) (*IValidatorStateStateTranstiionIterator, error) {

	logs, sub, err := _IValidatorState.contract.FilterLogs(opts, "StateTranstiion")
	if err != nil {
		return nil, err
	}
	return &IValidatorStateStateTranstiionIterator{contract: _IValidatorState.contract, event: "StateTranstiion", logs: logs, sub: sub}, nil
}

// WatchStateTranstiion is a free log subscription operation binding the contract event 0x8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee36091.
//
// Solidity: event StateTranstiion(uint8 oldStatw, uint8 newState)
func (_IValidatorState *IValidatorStateFilterer) WatchStateTranstiion(opts *bind.WatchOpts, sink chan<- *IValidatorStateStateTranstiion) (event.Subscription, error) {

	logs, sub, err := _IValidatorState.contract.WatchLogs(opts, "StateTranstiion")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IValidatorStateStateTranstiion)
				if err := _IValidatorState.contract.UnpackLog(event, "StateTranstiion", log); err != nil {
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

// ParseStateTranstiion is a log parse operation binding the contract event 0x8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee36091.
//
// Solidity: event StateTranstiion(uint8 oldStatw, uint8 newState)
func (_IValidatorState *IValidatorStateFilterer) ParseStateTranstiion(log types.Log) (*IValidatorStateStateTranstiion, error) {
	event := new(IValidatorStateStateTranstiion)
	if err := _IValidatorState.contract.UnpackLog(event, "StateTranstiion", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorStateMetaData contains all meta data concerning the ValidatorState contract.
var ValidatorStateMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"enumIValidatorState.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structIValidatorState.SystemValidatorUpdateRequest[]\",\"name\":\"valStates\",\"type\":\"tuple[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"SetIdleTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"SetPausedTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"enumIValidatorState.State\",\"name\":\"oldStatw\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"enumIValidatorState.State\",\"name\":\"newState\",\"type\":\"uint8\"}],\"name\":\"StateTranstiion\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ABOOK\",\"outputs\":[{\"internalType\":\"contractIAddressBook\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ValIdleTimeout\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ValPausedTimeout\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllValidators\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"enumIValidatorState.State\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"idleTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pausedTimeout\",\"type\":\"uint256\"}],\"internalType\":\"structIValidatorState.ValidatorState[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setCandInactive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setCandInactiveSuper\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setCandReady\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"duration\",\"type\":\"uint256\"}],\"name\":\"setIdleTimeoutSuper\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setValActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setValExiting\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setValPaused\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setValReady\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"enumIValidatorState.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structIValidatorState.SystemValidatorUpdateRequest[]\",\"name\":\"valStates\",\"type\":\"tuple[]\"}],\"name\":\"setValidatorStates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorStates\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"enumIValidatorState.State\",\"name\":\"state\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"idleTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pausedTimeout\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e675a75a": "ABOOK()",
		"26089a4c": "ValIdleTimeout()",
		"3848ebf1": "ValPausedTimeout()",
		"f3513a37": "getAllValidators()",
		"8da5cb5b": "owner()",
		"8f7a2984": "setCandInactive()",
		"30488def": "setCandInactiveSuper()",
		"8c587d4d": "setCandReady()",
		"497da016": "setIdleTimeoutSuper(uint256)",
		"b45fb92f": "setValActive()",
		"6d837b28": "setValExiting()",
		"e2dda149": "setValPaused()",
		"25c4639e": "setValReady()",
		"513e846c": "setValidatorStates((address,uint8)[])",
		"3ba9925a": "validatorStates(address)",
	},
	Bin: "0x6080604052600180546001600160a01b0319166002600160a01b03179055348015610028575f80fd5b5060405161171438038061171483398101604081905261004791610315565b5f80546001600160a01b0319163317905561006181610067565b50610484565b5f5b81518110156100fb575f828281518110610085576100856103fb565b60200260200101515f015190505f8383815181106100a5576100a56103fb565b60200260200101516020015190506100c382826100ff60201b60201c565b6100ce60038361019f565b5060048160088111156100e3576100e361040f565b036100f1576100f1826101bc565b5050600101610069565b5050565b6001600160a01b0382165f90815260026020526040908190205490517f8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee360919161015491600160a01b90910460ff16908490610443565b60405180910390a16001600160a01b0382165f908152600260205260409020805482919060ff60a01b1916600160a01b8360088111156101965761019661040f565b02179055505050565b5f6101b3836001600160a01b03841661025d565b90505b92915050565b6001600160a01b0381165f90815260026020526040902060010154421115610207576101eb62278d0042610465565b6001600160a01b0382165f908152600260205260409020600101555b6001600160a01b0381165f81815260026020908152604091829020600101548251938452908301527fc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b818910160405180910390a150565b5f8181526001830160205260408120546102a257508154600181810184555f8481526020808220909301849055845484825282860190935260409020919091556101b6565b505f6101b6565b634e487b7160e01b5f52604160045260245ffd5b604080519081016001600160401b03811182821017156102df576102df6102a9565b60405290565b604051601f8201601f191681016001600160401b038111828210171561030d5761030d6102a9565b604052919050565b5f6020808385031215610326575f80fd5b82516001600160401b038082111561033c575f80fd5b818501915085601f83011261034f575f80fd5b815181811115610361576103616102a9565b61036f848260051b016102e5565b818152848101925060069190911b83018401908782111561038e575f80fd5b928401925b818410156103f057604084890312156103aa575f80fd5b6103b26102bd565b84516001600160a01b03811681146103c8575f80fd5b815284860151600981106103da575f80fd5b8187015283526040939093019291840191610393565b979650505050505050565b634e487b7160e01b5f52603260045260245ffd5b634e487b7160e01b5f52602160045260245ffd5b6009811061043f57634e487b7160e01b5f52602160045260245ffd5b9052565b604081016104518285610423565b61045e6020830184610423565b9392505050565b808201808211156101b657634e487b7160e01b5f52601160045260245ffd5b611283806104915f395ff3fe608060405234801561000f575f80fd5b50600436106100f0575f3560e01c80636d837b2811610093578063b45fb92f11610063578063b45fb92f146101e8578063e2dda149146101f0578063e675a75a146101f8578063f3513a3714610201575f80fd5b80636d837b28146101a65780638c587d4d146101ae5780638da5cb5b146101b65780638f7a2984146101e0575f80fd5b80633848ebf1116100ce5780633848ebf1146101235780633ba9925a1461012c578063497da01614610180578063513e846c14610193575f80fd5b806325c4639e146100f457806326089a4c146100fe57806330488def1461011b575b5f80fd5b6100fc610216565b005b61010862278d0081565b6040519081526020015b60405180910390f35b6100fc610417565b61010861708081565b61017061013a366004610e31565b600260208190525f91825260409091208054600182015491909201546001600160a01b03831692600160a01b900460ff16919084565b6040516101129493929190610e87565b6100fc61018e366004610eb6565b610464565b6100fc6101a1366004610f3b565b6104ea565b6100fc610531565b6100fc61060f565b5f546101c8906001600160a01b031681565b6040516001600160a01b039091168152602001610112565b6100fc610813565b6100fc6108c8565b6100fc610950565b6101c861040081565b6102096109e1565b6040516101129190611017565b33610222600382610b29565b6102475760405162461bcd60e51b815260040161023e9061108c565b60405180910390fd5b336004806001600160a01b0383165f90815260026020526040902054600160a01b900460ff16600881111561027e5761027e610e53565b1461029b5760405162461bcd60e51b815260040161023e906110d3565b604051630aabaead60e11b81523360048201819052905f90610400906315575d5a90602401606060405180830381865afa1580156102db573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906102ff919061110a565b509150505f816001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa158015610340573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103649190611154565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103a0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103c49190611154565b6103ce919061117f565b90506a0422ca8b0a00a4250000008110156103fb5760405162461bcd60e51b815260040161023e90611192565b610406836007610b4f565b61040f83610bef565b505050505050565b5f546001600160a01b031633148061043957506001546001600160a01b031633145b6104555760405162461bcd60e51b815260040161023e906111d4565b33610461816001610b4f565b50565b33610470600382610b29565b61048c5760405162461bcd60e51b815260040161023e9061108c565b5f546001600160a01b03163314806104ae57506001546001600160a01b031633145b6104ca5760405162461bcd60e51b815260040161023e906111d4565b6104d4824261120b565b335f908152600260205260409020600101555050565b5f546001600160a01b031633148061050c57506001546001600160a01b031633145b6105285760405162461bcd60e51b815260040161023e906111d4565b61046181610c91565b3361053d600382610b29565b6105595760405162461bcd60e51b815260040161023e9061108c565b3360086005816001600160a01b0384165f90815260026020526040902054600160a01b900460ff16600881111561059257610592610e53565b14806105e057508060088111156105ab576105ab610e53565b6001600160a01b0384165f90815260026020526040902054600160a01b900460ff1660088111156105de576105de610e53565b145b6105fc5760405162461bcd60e51b815260040161023e906110d3565b33610608816006610b4f565b5050505050565b335f818152600260205260408120546001908290600160a01b900460ff16600881111561063e5761063e610e53565b148061068c575080600881111561065757610657610e53565b6001600160a01b0384165f90815260026020526040902054600160a01b900460ff16600881111561068a5761068a610e53565b145b6106a85760405162461bcd60e51b815260040161023e906110d3565b604051630aabaead60e11b81523360048201819052905f90610400906315575d5a90602401606060405180830381865afa1580156106e8573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061070c919061110a565b509150505f816001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa15801561074d573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107719190611154565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa1580156107ad573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107d19190611154565b6107db919061117f565b90506a0422ca8b0a00a4250000008110156108085760405162461bcd60e51b815260040161023e90611192565b61040f836002610b4f565b3360026004816001600160a01b0384165f90815260026020526040902054600160a01b900460ff16600881111561084c5761084c610e53565b148061089a575080600881111561086557610865610e53565b6001600160a01b0384165f90815260026020526040902054600160a01b900460ff16600881111561089857610898610e53565b145b6108b65760405162461bcd60e51b815260040161023e906110d3565b336108c2816001610b4f565b50505050565b336108d4600382610b29565b6108f05760405162461bcd60e51b815260040161023e9061108c565b336005806001600160a01b0383165f90815260026020526040902054600160a01b900460ff16600881111561092757610927610e53565b146109445760405162461bcd60e51b815260040161023e906110d3565b336108c2816008610b4f565b3361095c600382610b29565b6109785760405162461bcd60e51b815260040161023e9061108c565b336008806001600160a01b0383165f90815260026020526040902054600160a01b900460ff1660088111156109af576109af610e53565b146109cc5760405162461bcd60e51b815260040161023e906110d3565b336109d8816005610b4f565b6108c281610d23565b60605f6109ee6003610d83565b90505f8167ffffffffffffffff811115610a0a57610a0a610ecd565b604051908082528060200260200182016040528015610a6457816020015b610a5160408051608081019091525f808252602082019081526020015f81526020015f81525090565b815260200190600190039081610a285790505b5090505f5b82811015610b22575f610a7d600383610d8c565b604080516080810182526001600160a01b0383168082525f90815260026020908152929020549293509190820190600160a01b900460ff166008811115610ac657610ac6610e53565b81526001600160a01b0383165f8181526002602081815260408084206001810154838801529490935281905291909101549101528351849084908110610b0e57610b0e61121e565b602090810291909101015250600101610a69565b5092915050565b6001600160a01b0381165f90815260018301602052604081205415155b90505b92915050565b6001600160a01b0382165f90815260026020526040908190205490517f8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee3609191610ba491600160a01b90910460ff16908490611232565b60405180910390a16001600160a01b0382165f908152600260205260409020805482919060ff60a01b1916600160a01b836008811115610be657610be6610e53565b02179055505050565b6001600160a01b0381165f90815260026020526040902060010154421115610c3a57610c1e62278d004261120b565b6001600160a01b0382165f908152600260205260409020600101555b6001600160a01b0381165f81815260026020908152604091829020600101548251938452908301527fc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b81891015b60405180910390a150565b5f5b8151811015610d1f575f828281518110610caf57610caf61121e565b60200260200101515f015190505f838381518110610ccf57610ccf61121e565b6020026020010151602001519050610ce78282610b4f565b610cf2600383610d97565b506004816008811115610d0757610d07610e53565b03610d1557610d1582610bef565b5050600101610c93565b5050565b610d2f6170804261120b565b6001600160a01b0382165f8181526002602081815260409283902090910184905581519283528201929092527fbc1a6d4bb677fb39bdb63d1713dc004bd314de87091c5dfbba6a8d78408d0dc99101610c86565b5f610b49825490565b5f610b468383610dab565b5f610b46836001600160a01b038416610dd1565b5f825f018281548110610dc057610dc061121e565b905f5260205f200154905092915050565b5f818152600183016020526040812054610e1657508154600181810184555f848152602080822090930184905584548482528286019093526040902091909155610b49565b505f610b49565b6001600160a01b0381168114610461575f80fd5b5f60208284031215610e41575f80fd5b8135610e4c81610e1d565b9392505050565b634e487b7160e01b5f52602160045260245ffd5b60098110610e8357634e487b7160e01b5f52602160045260245ffd5b9052565b6001600160a01b038516815260808101610ea46020830186610e67565b60408201939093526060015292915050565b5f60208284031215610ec6575f80fd5b5035919050565b634e487b7160e01b5f52604160045260245ffd5b6040805190810167ffffffffffffffff81118282101715610f0457610f04610ecd565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610f3357610f33610ecd565b604052919050565b5f6020808385031215610f4c575f80fd5b823567ffffffffffffffff80821115610f63575f80fd5b818501915085601f830112610f76575f80fd5b813581811115610f8857610f88610ecd565b610f96848260051b01610f0a565b818152848101925060069190911b830184019087821115610fb5575f80fd5b928401925b8184101561100c5760408489031215610fd1575f80fd5b610fd9610ee1565b8435610fe481610e1d565b81528486013560098110610ff6575f80fd5b8187015283526040939093019291840191610fba565b979650505050505050565b602080825282518282018190525f919060409081850190868401855b8281101561107f57815180516001600160a01b031685528681015161105a88870182610e67565b5080860151858701526060908101519085015260809093019290850190600101611033565b5091979650505050505050565b60208082526027908201527f4572726f723a20676976656e2061646472657373206973206e6f74206120766160408201526634b63230ba37b960c91b606082015260800190565b6020808252601c908201527f4572726f723a204e6f7420616e20657870656374656420737461746500000000604082015260600190565b5f805f6060848603121561111c575f80fd5b835161112781610e1d565b602085015190935061113881610e1d565b604085015190925061114981610e1d565b809150509250925092565b5f60208284031215611164575f80fd5b5051919050565b634e487b7160e01b5f52601160045260245ffd5b81810381811115610b4957610b4961116b565b60208082526022908201527f4572726f723a20696e73756666696369656e74207374616b696e6720616d6f756040820152611b9d60f21b606082015260800190565b6020808252601e908201527f4572726f723a2043616c6c6572206973206e6f7420746865206f776e65720000604082015260600190565b80820180821115610b4957610b4961116b565b634e487b7160e01b5f52603260045260245ffd5b604081016112408285610e67565b610e4c6020830184610e6756fea2646970667358221220acff3646ec90a918dcc2c55c66205a0bd986c42c27f7cd7911674397303f217664736f6c63430008190033",
}

// ValidatorStateABI is the input ABI used to generate the binding from.
// Deprecated: Use ValidatorStateMetaData.ABI instead.
var ValidatorStateABI = ValidatorStateMetaData.ABI

// ValidatorStateBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ValidatorStateBinRuntime = `608060405234801561000f575f80fd5b50600436106100f0575f3560e01c80636d837b2811610093578063b45fb92f11610063578063b45fb92f146101e8578063e2dda149146101f0578063e675a75a146101f8578063f3513a3714610201575f80fd5b80636d837b28146101a65780638c587d4d146101ae5780638da5cb5b146101b65780638f7a2984146101e0575f80fd5b80633848ebf1116100ce5780633848ebf1146101235780633ba9925a1461012c578063497da01614610180578063513e846c14610193575f80fd5b806325c4639e146100f457806326089a4c146100fe57806330488def1461011b575b5f80fd5b6100fc610216565b005b61010862278d0081565b6040519081526020015b60405180910390f35b6100fc610417565b61010861708081565b61017061013a366004610e31565b600260208190525f91825260409091208054600182015491909201546001600160a01b03831692600160a01b900460ff16919084565b6040516101129493929190610e87565b6100fc61018e366004610eb6565b610464565b6100fc6101a1366004610f3b565b6104ea565b6100fc610531565b6100fc61060f565b5f546101c8906001600160a01b031681565b6040516001600160a01b039091168152602001610112565b6100fc610813565b6100fc6108c8565b6100fc610950565b6101c861040081565b6102096109e1565b6040516101129190611017565b33610222600382610b29565b6102475760405162461bcd60e51b815260040161023e9061108c565b60405180910390fd5b336004806001600160a01b0383165f90815260026020526040902054600160a01b900460ff16600881111561027e5761027e610e53565b1461029b5760405162461bcd60e51b815260040161023e906110d3565b604051630aabaead60e11b81523360048201819052905f90610400906315575d5a90602401606060405180830381865afa1580156102db573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906102ff919061110a565b509150505f816001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa158015610340573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103649190611154565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa1580156103a0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906103c49190611154565b6103ce919061117f565b90506a0422ca8b0a00a4250000008110156103fb5760405162461bcd60e51b815260040161023e90611192565b610406836007610b4f565b61040f83610bef565b505050505050565b5f546001600160a01b031633148061043957506001546001600160a01b031633145b6104555760405162461bcd60e51b815260040161023e906111d4565b33610461816001610b4f565b50565b33610470600382610b29565b61048c5760405162461bcd60e51b815260040161023e9061108c565b5f546001600160a01b03163314806104ae57506001546001600160a01b031633145b6104ca5760405162461bcd60e51b815260040161023e906111d4565b6104d4824261120b565b335f908152600260205260409020600101555050565b5f546001600160a01b031633148061050c57506001546001600160a01b031633145b6105285760405162461bcd60e51b815260040161023e906111d4565b61046181610c91565b3361053d600382610b29565b6105595760405162461bcd60e51b815260040161023e9061108c565b3360086005816001600160a01b0384165f90815260026020526040902054600160a01b900460ff16600881111561059257610592610e53565b14806105e057508060088111156105ab576105ab610e53565b6001600160a01b0384165f90815260026020526040902054600160a01b900460ff1660088111156105de576105de610e53565b145b6105fc5760405162461bcd60e51b815260040161023e906110d3565b33610608816006610b4f565b5050505050565b335f818152600260205260408120546001908290600160a01b900460ff16600881111561063e5761063e610e53565b148061068c575080600881111561065757610657610e53565b6001600160a01b0384165f90815260026020526040902054600160a01b900460ff16600881111561068a5761068a610e53565b145b6106a85760405162461bcd60e51b815260040161023e906110d3565b604051630aabaead60e11b81523360048201819052905f90610400906315575d5a90602401606060405180830381865afa1580156106e8573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061070c919061110a565b509150505f816001600160a01b031663630b11466040518163ffffffff1660e01b8152600401602060405180830381865afa15801561074d573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107719190611154565b826001600160a01b0316634cf088d96040518163ffffffff1660e01b8152600401602060405180830381865afa1580156107ad573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107d19190611154565b6107db919061117f565b90506a0422ca8b0a00a4250000008110156108085760405162461bcd60e51b815260040161023e90611192565b61040f836002610b4f565b3360026004816001600160a01b0384165f90815260026020526040902054600160a01b900460ff16600881111561084c5761084c610e53565b148061089a575080600881111561086557610865610e53565b6001600160a01b0384165f90815260026020526040902054600160a01b900460ff16600881111561089857610898610e53565b145b6108b65760405162461bcd60e51b815260040161023e906110d3565b336108c2816001610b4f565b50505050565b336108d4600382610b29565b6108f05760405162461bcd60e51b815260040161023e9061108c565b336005806001600160a01b0383165f90815260026020526040902054600160a01b900460ff16600881111561092757610927610e53565b146109445760405162461bcd60e51b815260040161023e906110d3565b336108c2816008610b4f565b3361095c600382610b29565b6109785760405162461bcd60e51b815260040161023e9061108c565b336008806001600160a01b0383165f90815260026020526040902054600160a01b900460ff1660088111156109af576109af610e53565b146109cc5760405162461bcd60e51b815260040161023e906110d3565b336109d8816005610b4f565b6108c281610d23565b60605f6109ee6003610d83565b90505f8167ffffffffffffffff811115610a0a57610a0a610ecd565b604051908082528060200260200182016040528015610a6457816020015b610a5160408051608081019091525f808252602082019081526020015f81526020015f81525090565b815260200190600190039081610a285790505b5090505f5b82811015610b22575f610a7d600383610d8c565b604080516080810182526001600160a01b0383168082525f90815260026020908152929020549293509190820190600160a01b900460ff166008811115610ac657610ac6610e53565b81526001600160a01b0383165f8181526002602081815260408084206001810154838801529490935281905291909101549101528351849084908110610b0e57610b0e61121e565b602090810291909101015250600101610a69565b5092915050565b6001600160a01b0381165f90815260018301602052604081205415155b90505b92915050565b6001600160a01b0382165f90815260026020526040908190205490517f8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee3609191610ba491600160a01b90910460ff16908490611232565b60405180910390a16001600160a01b0382165f908152600260205260409020805482919060ff60a01b1916600160a01b836008811115610be657610be6610e53565b02179055505050565b6001600160a01b0381165f90815260026020526040902060010154421115610c3a57610c1e62278d004261120b565b6001600160a01b0382165f908152600260205260409020600101555b6001600160a01b0381165f81815260026020908152604091829020600101548251938452908301527fc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b81891015b60405180910390a150565b5f5b8151811015610d1f575f828281518110610caf57610caf61121e565b60200260200101515f015190505f838381518110610ccf57610ccf61121e565b6020026020010151602001519050610ce78282610b4f565b610cf2600383610d97565b506004816008811115610d0757610d07610e53565b03610d1557610d1582610bef565b5050600101610c93565b5050565b610d2f6170804261120b565b6001600160a01b0382165f8181526002602081815260409283902090910184905581519283528201929092527fbc1a6d4bb677fb39bdb63d1713dc004bd314de87091c5dfbba6a8d78408d0dc99101610c86565b5f610b49825490565b5f610b468383610dab565b5f610b46836001600160a01b038416610dd1565b5f825f018281548110610dc057610dc061121e565b905f5260205f200154905092915050565b5f818152600183016020526040812054610e1657508154600181810184555f848152602080822090930184905584548482528286019093526040902091909155610b49565b505f610b49565b6001600160a01b0381168114610461575f80fd5b5f60208284031215610e41575f80fd5b8135610e4c81610e1d565b9392505050565b634e487b7160e01b5f52602160045260245ffd5b60098110610e8357634e487b7160e01b5f52602160045260245ffd5b9052565b6001600160a01b038516815260808101610ea46020830186610e67565b60408201939093526060015292915050565b5f60208284031215610ec6575f80fd5b5035919050565b634e487b7160e01b5f52604160045260245ffd5b6040805190810167ffffffffffffffff81118282101715610f0457610f04610ecd565b60405290565b604051601f8201601f1916810167ffffffffffffffff81118282101715610f3357610f33610ecd565b604052919050565b5f6020808385031215610f4c575f80fd5b823567ffffffffffffffff80821115610f63575f80fd5b818501915085601f830112610f76575f80fd5b813581811115610f8857610f88610ecd565b610f96848260051b01610f0a565b818152848101925060069190911b830184019087821115610fb5575f80fd5b928401925b8184101561100c5760408489031215610fd1575f80fd5b610fd9610ee1565b8435610fe481610e1d565b81528486013560098110610ff6575f80fd5b8187015283526040939093019291840191610fba565b979650505050505050565b602080825282518282018190525f919060409081850190868401855b8281101561107f57815180516001600160a01b031685528681015161105a88870182610e67565b5080860151858701526060908101519085015260809093019290850190600101611033565b5091979650505050505050565b60208082526027908201527f4572726f723a20676976656e2061646472657373206973206e6f74206120766160408201526634b63230ba37b960c91b606082015260800190565b6020808252601c908201527f4572726f723a204e6f7420616e20657870656374656420737461746500000000604082015260600190565b5f805f6060848603121561111c575f80fd5b835161112781610e1d565b602085015190935061113881610e1d565b604085015190925061114981610e1d565b809150509250925092565b5f60208284031215611164575f80fd5b5051919050565b634e487b7160e01b5f52601160045260245ffd5b81810381811115610b4957610b4961116b565b60208082526022908201527f4572726f723a20696e73756666696369656e74207374616b696e6720616d6f756040820152611b9d60f21b606082015260800190565b6020808252601e908201527f4572726f723a2043616c6c6572206973206e6f7420746865206f776e65720000604082015260600190565b80820180821115610b4957610b4961116b565b634e487b7160e01b5f52603260045260245ffd5b604081016112408285610e67565b610e4c6020830184610e6756fea2646970667358221220acff3646ec90a918dcc2c55c66205a0bd986c42c27f7cd7911674397303f217664736f6c63430008190033`

// Deprecated: Use ValidatorStateMetaData.Sigs instead.
// ValidatorStateFuncSigs maps the 4-byte function signature to its string representation.
var ValidatorStateFuncSigs = ValidatorStateMetaData.Sigs

// ValidatorStateBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ValidatorStateMetaData.Bin instead.
var ValidatorStateBin = ValidatorStateMetaData.Bin

// DeployValidatorState deploys a new Kaia contract, binding an instance of ValidatorState to it.
func DeployValidatorState(auth *bind.TransactOpts, backend bind.ContractBackend, valStates []IValidatorStateSystemValidatorUpdateRequest) (common.Address, *types.Transaction, *ValidatorState, error) {
	parsed, err := ValidatorStateMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ValidatorStateBin), backend, valStates)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ValidatorState{ValidatorStateCaller: ValidatorStateCaller{contract: contract}, ValidatorStateTransactor: ValidatorStateTransactor{contract: contract}, ValidatorStateFilterer: ValidatorStateFilterer{contract: contract}}, nil
}

// ValidatorState is an auto generated Go binding around a Kaia contract.
type ValidatorState struct {
	ValidatorStateCaller     // Read-only binding to the contract
	ValidatorStateTransactor // Write-only binding to the contract
	ValidatorStateFilterer   // Log filterer for contract events
}

// ValidatorStateCaller is an auto generated read-only Go binding around a Kaia contract.
type ValidatorStateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorStateTransactor is an auto generated write-only Go binding around a Kaia contract.
type ValidatorStateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorStateFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ValidatorStateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorStateSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ValidatorStateSession struct {
	Contract     *ValidatorState   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValidatorStateCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ValidatorStateCallerSession struct {
	Contract *ValidatorStateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ValidatorStateTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ValidatorStateTransactorSession struct {
	Contract     *ValidatorStateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ValidatorStateRaw is an auto generated low-level Go binding around a Kaia contract.
type ValidatorStateRaw struct {
	Contract *ValidatorState // Generic contract binding to access the raw methods on
}

// ValidatorStateCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ValidatorStateCallerRaw struct {
	Contract *ValidatorStateCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatorStateTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ValidatorStateTransactorRaw struct {
	Contract *ValidatorStateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatorState creates a new instance of ValidatorState, bound to a specific deployed contract.
func NewValidatorState(address common.Address, backend bind.ContractBackend) (*ValidatorState, error) {
	contract, err := bindValidatorState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ValidatorState{ValidatorStateCaller: ValidatorStateCaller{contract: contract}, ValidatorStateTransactor: ValidatorStateTransactor{contract: contract}, ValidatorStateFilterer: ValidatorStateFilterer{contract: contract}}, nil
}

// NewValidatorStateCaller creates a new read-only instance of ValidatorState, bound to a specific deployed contract.
func NewValidatorStateCaller(address common.Address, caller bind.ContractCaller) (*ValidatorStateCaller, error) {
	contract, err := bindValidatorState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorStateCaller{contract: contract}, nil
}

// NewValidatorStateTransactor creates a new write-only instance of ValidatorState, bound to a specific deployed contract.
func NewValidatorStateTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatorStateTransactor, error) {
	contract, err := bindValidatorState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorStateTransactor{contract: contract}, nil
}

// NewValidatorStateFilterer creates a new log filterer instance of ValidatorState, bound to a specific deployed contract.
func NewValidatorStateFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatorStateFilterer, error) {
	contract, err := bindValidatorState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatorStateFilterer{contract: contract}, nil
}

// bindValidatorState binds a generic wrapper to an already deployed contract.
func bindValidatorState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ValidatorStateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorState *ValidatorStateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ValidatorState.Contract.ValidatorStateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorState *ValidatorStateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.Contract.ValidatorStateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorState *ValidatorStateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidatorState.Contract.ValidatorStateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorState *ValidatorStateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ValidatorState.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorState *ValidatorStateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorState *ValidatorStateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidatorState.Contract.contract.Transact(opts, method, params...)
}

// ABOOK is a free data retrieval call binding the contract method 0xe675a75a.
//
// Solidity: function ABOOK() view returns(address)
func (_ValidatorState *ValidatorStateCaller) ABOOK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ValidatorState.contract.Call(opts, &out, "ABOOK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ABOOK is a free data retrieval call binding the contract method 0xe675a75a.
//
// Solidity: function ABOOK() view returns(address)
func (_ValidatorState *ValidatorStateSession) ABOOK() (common.Address, error) {
	return _ValidatorState.Contract.ABOOK(&_ValidatorState.CallOpts)
}

// ABOOK is a free data retrieval call binding the contract method 0xe675a75a.
//
// Solidity: function ABOOK() view returns(address)
func (_ValidatorState *ValidatorStateCallerSession) ABOOK() (common.Address, error) {
	return _ValidatorState.Contract.ABOOK(&_ValidatorState.CallOpts)
}

// ValIdleTimeout is a free data retrieval call binding the contract method 0x26089a4c.
//
// Solidity: function ValIdleTimeout() view returns(uint256)
func (_ValidatorState *ValidatorStateCaller) ValIdleTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorState.contract.Call(opts, &out, "ValIdleTimeout")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValIdleTimeout is a free data retrieval call binding the contract method 0x26089a4c.
//
// Solidity: function ValIdleTimeout() view returns(uint256)
func (_ValidatorState *ValidatorStateSession) ValIdleTimeout() (*big.Int, error) {
	return _ValidatorState.Contract.ValIdleTimeout(&_ValidatorState.CallOpts)
}

// ValIdleTimeout is a free data retrieval call binding the contract method 0x26089a4c.
//
// Solidity: function ValIdleTimeout() view returns(uint256)
func (_ValidatorState *ValidatorStateCallerSession) ValIdleTimeout() (*big.Int, error) {
	return _ValidatorState.Contract.ValIdleTimeout(&_ValidatorState.CallOpts)
}

// ValPausedTimeout is a free data retrieval call binding the contract method 0x3848ebf1.
//
// Solidity: function ValPausedTimeout() view returns(uint256)
func (_ValidatorState *ValidatorStateCaller) ValPausedTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ValidatorState.contract.Call(opts, &out, "ValPausedTimeout")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ValPausedTimeout is a free data retrieval call binding the contract method 0x3848ebf1.
//
// Solidity: function ValPausedTimeout() view returns(uint256)
func (_ValidatorState *ValidatorStateSession) ValPausedTimeout() (*big.Int, error) {
	return _ValidatorState.Contract.ValPausedTimeout(&_ValidatorState.CallOpts)
}

// ValPausedTimeout is a free data retrieval call binding the contract method 0x3848ebf1.
//
// Solidity: function ValPausedTimeout() view returns(uint256)
func (_ValidatorState *ValidatorStateCallerSession) ValPausedTimeout() (*big.Int, error) {
	return _ValidatorState.Contract.ValPausedTimeout(&_ValidatorState.CallOpts)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((address,uint8,uint256,uint256)[])
func (_ValidatorState *ValidatorStateCaller) GetAllValidators(opts *bind.CallOpts) ([]IValidatorStateValidatorState, error) {
	var out []interface{}
	err := _ValidatorState.contract.Call(opts, &out, "getAllValidators")

	if err != nil {
		return *new([]IValidatorStateValidatorState), err
	}

	out0 := *abi.ConvertType(out[0], new([]IValidatorStateValidatorState)).(*[]IValidatorStateValidatorState)

	return out0, err

}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((address,uint8,uint256,uint256)[])
func (_ValidatorState *ValidatorStateSession) GetAllValidators() ([]IValidatorStateValidatorState, error) {
	return _ValidatorState.Contract.GetAllValidators(&_ValidatorState.CallOpts)
}

// GetAllValidators is a free data retrieval call binding the contract method 0xf3513a37.
//
// Solidity: function getAllValidators() view returns((address,uint8,uint256,uint256)[])
func (_ValidatorState *ValidatorStateCallerSession) GetAllValidators() ([]IValidatorStateValidatorState, error) {
	return _ValidatorState.Contract.GetAllValidators(&_ValidatorState.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidatorState *ValidatorStateCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ValidatorState.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidatorState *ValidatorStateSession) Owner() (common.Address, error) {
	return _ValidatorState.Contract.Owner(&_ValidatorState.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ValidatorState *ValidatorStateCallerSession) Owner() (common.Address, error) {
	return _ValidatorState.Contract.Owner(&_ValidatorState.CallOpts)
}

// ValidatorStates is a free data retrieval call binding the contract method 0x3ba9925a.
//
// Solidity: function validatorStates(address ) view returns(address addr, uint8 state, uint256 idleTimeout, uint256 pausedTimeout)
func (_ValidatorState *ValidatorStateCaller) ValidatorStates(opts *bind.CallOpts, arg0 common.Address) (struct {
	Addr          common.Address
	State         uint8
	IdleTimeout   *big.Int
	PausedTimeout *big.Int
}, error) {
	var out []interface{}
	err := _ValidatorState.contract.Call(opts, &out, "validatorStates", arg0)

	outstruct := new(struct {
		Addr          common.Address
		State         uint8
		IdleTimeout   *big.Int
		PausedTimeout *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Addr = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.State = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.IdleTimeout = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PausedTimeout = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ValidatorStates is a free data retrieval call binding the contract method 0x3ba9925a.
//
// Solidity: function validatorStates(address ) view returns(address addr, uint8 state, uint256 idleTimeout, uint256 pausedTimeout)
func (_ValidatorState *ValidatorStateSession) ValidatorStates(arg0 common.Address) (struct {
	Addr          common.Address
	State         uint8
	IdleTimeout   *big.Int
	PausedTimeout *big.Int
}, error) {
	return _ValidatorState.Contract.ValidatorStates(&_ValidatorState.CallOpts, arg0)
}

// ValidatorStates is a free data retrieval call binding the contract method 0x3ba9925a.
//
// Solidity: function validatorStates(address ) view returns(address addr, uint8 state, uint256 idleTimeout, uint256 pausedTimeout)
func (_ValidatorState *ValidatorStateCallerSession) ValidatorStates(arg0 common.Address) (struct {
	Addr          common.Address
	State         uint8
	IdleTimeout   *big.Int
	PausedTimeout *big.Int
}, error) {
	return _ValidatorState.Contract.ValidatorStates(&_ValidatorState.CallOpts, arg0)
}

// SetCandInactive is a paid mutator transaction binding the contract method 0x8f7a2984.
//
// Solidity: function setCandInactive() returns()
func (_ValidatorState *ValidatorStateTransactor) SetCandInactive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setCandInactive")
}

// SetCandInactive is a paid mutator transaction binding the contract method 0x8f7a2984.
//
// Solidity: function setCandInactive() returns()
func (_ValidatorState *ValidatorStateSession) SetCandInactive() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetCandInactive(&_ValidatorState.TransactOpts)
}

// SetCandInactive is a paid mutator transaction binding the contract method 0x8f7a2984.
//
// Solidity: function setCandInactive() returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetCandInactive() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetCandInactive(&_ValidatorState.TransactOpts)
}

// SetCandInactiveSuper is a paid mutator transaction binding the contract method 0x30488def.
//
// Solidity: function setCandInactiveSuper() returns()
func (_ValidatorState *ValidatorStateTransactor) SetCandInactiveSuper(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setCandInactiveSuper")
}

// SetCandInactiveSuper is a paid mutator transaction binding the contract method 0x30488def.
//
// Solidity: function setCandInactiveSuper() returns()
func (_ValidatorState *ValidatorStateSession) SetCandInactiveSuper() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetCandInactiveSuper(&_ValidatorState.TransactOpts)
}

// SetCandInactiveSuper is a paid mutator transaction binding the contract method 0x30488def.
//
// Solidity: function setCandInactiveSuper() returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetCandInactiveSuper() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetCandInactiveSuper(&_ValidatorState.TransactOpts)
}

// SetCandReady is a paid mutator transaction binding the contract method 0x8c587d4d.
//
// Solidity: function setCandReady() returns()
func (_ValidatorState *ValidatorStateTransactor) SetCandReady(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setCandReady")
}

// SetCandReady is a paid mutator transaction binding the contract method 0x8c587d4d.
//
// Solidity: function setCandReady() returns()
func (_ValidatorState *ValidatorStateSession) SetCandReady() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetCandReady(&_ValidatorState.TransactOpts)
}

// SetCandReady is a paid mutator transaction binding the contract method 0x8c587d4d.
//
// Solidity: function setCandReady() returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetCandReady() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetCandReady(&_ValidatorState.TransactOpts)
}

// SetIdleTimeoutSuper is a paid mutator transaction binding the contract method 0x497da016.
//
// Solidity: function setIdleTimeoutSuper(uint256 duration) returns()
func (_ValidatorState *ValidatorStateTransactor) SetIdleTimeoutSuper(opts *bind.TransactOpts, duration *big.Int) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setIdleTimeoutSuper", duration)
}

// SetIdleTimeoutSuper is a paid mutator transaction binding the contract method 0x497da016.
//
// Solidity: function setIdleTimeoutSuper(uint256 duration) returns()
func (_ValidatorState *ValidatorStateSession) SetIdleTimeoutSuper(duration *big.Int) (*types.Transaction, error) {
	return _ValidatorState.Contract.SetIdleTimeoutSuper(&_ValidatorState.TransactOpts, duration)
}

// SetIdleTimeoutSuper is a paid mutator transaction binding the contract method 0x497da016.
//
// Solidity: function setIdleTimeoutSuper(uint256 duration) returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetIdleTimeoutSuper(duration *big.Int) (*types.Transaction, error) {
	return _ValidatorState.Contract.SetIdleTimeoutSuper(&_ValidatorState.TransactOpts, duration)
}

// SetValActive is a paid mutator transaction binding the contract method 0xb45fb92f.
//
// Solidity: function setValActive() returns()
func (_ValidatorState *ValidatorStateTransactor) SetValActive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setValActive")
}

// SetValActive is a paid mutator transaction binding the contract method 0xb45fb92f.
//
// Solidity: function setValActive() returns()
func (_ValidatorState *ValidatorStateSession) SetValActive() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValActive(&_ValidatorState.TransactOpts)
}

// SetValActive is a paid mutator transaction binding the contract method 0xb45fb92f.
//
// Solidity: function setValActive() returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetValActive() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValActive(&_ValidatorState.TransactOpts)
}

// SetValExiting is a paid mutator transaction binding the contract method 0x6d837b28.
//
// Solidity: function setValExiting() returns()
func (_ValidatorState *ValidatorStateTransactor) SetValExiting(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setValExiting")
}

// SetValExiting is a paid mutator transaction binding the contract method 0x6d837b28.
//
// Solidity: function setValExiting() returns()
func (_ValidatorState *ValidatorStateSession) SetValExiting() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValExiting(&_ValidatorState.TransactOpts)
}

// SetValExiting is a paid mutator transaction binding the contract method 0x6d837b28.
//
// Solidity: function setValExiting() returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetValExiting() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValExiting(&_ValidatorState.TransactOpts)
}

// SetValPaused is a paid mutator transaction binding the contract method 0xe2dda149.
//
// Solidity: function setValPaused() returns()
func (_ValidatorState *ValidatorStateTransactor) SetValPaused(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setValPaused")
}

// SetValPaused is a paid mutator transaction binding the contract method 0xe2dda149.
//
// Solidity: function setValPaused() returns()
func (_ValidatorState *ValidatorStateSession) SetValPaused() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValPaused(&_ValidatorState.TransactOpts)
}

// SetValPaused is a paid mutator transaction binding the contract method 0xe2dda149.
//
// Solidity: function setValPaused() returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetValPaused() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValPaused(&_ValidatorState.TransactOpts)
}

// SetValReady is a paid mutator transaction binding the contract method 0x25c4639e.
//
// Solidity: function setValReady() returns()
func (_ValidatorState *ValidatorStateTransactor) SetValReady(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setValReady")
}

// SetValReady is a paid mutator transaction binding the contract method 0x25c4639e.
//
// Solidity: function setValReady() returns()
func (_ValidatorState *ValidatorStateSession) SetValReady() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValReady(&_ValidatorState.TransactOpts)
}

// SetValReady is a paid mutator transaction binding the contract method 0x25c4639e.
//
// Solidity: function setValReady() returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetValReady() (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValReady(&_ValidatorState.TransactOpts)
}

// SetValidatorStates is a paid mutator transaction binding the contract method 0x513e846c.
//
// Solidity: function setValidatorStates((address,uint8)[] valStates) returns()
func (_ValidatorState *ValidatorStateTransactor) SetValidatorStates(opts *bind.TransactOpts, valStates []IValidatorStateSystemValidatorUpdateRequest) (*types.Transaction, error) {
	return _ValidatorState.contract.Transact(opts, "setValidatorStates", valStates)
}

// SetValidatorStates is a paid mutator transaction binding the contract method 0x513e846c.
//
// Solidity: function setValidatorStates((address,uint8)[] valStates) returns()
func (_ValidatorState *ValidatorStateSession) SetValidatorStates(valStates []IValidatorStateSystemValidatorUpdateRequest) (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValidatorStates(&_ValidatorState.TransactOpts, valStates)
}

// SetValidatorStates is a paid mutator transaction binding the contract method 0x513e846c.
//
// Solidity: function setValidatorStates((address,uint8)[] valStates) returns()
func (_ValidatorState *ValidatorStateTransactorSession) SetValidatorStates(valStates []IValidatorStateSystemValidatorUpdateRequest) (*types.Transaction, error) {
	return _ValidatorState.Contract.SetValidatorStates(&_ValidatorState.TransactOpts, valStates)
}

// ValidatorStateSetIdleTimeoutIterator is returned from FilterSetIdleTimeout and is used to iterate over the raw logs and unpacked data for SetIdleTimeout events raised by the ValidatorState contract.
type ValidatorStateSetIdleTimeoutIterator struct {
	Event *ValidatorStateSetIdleTimeout // Event containing the contract specifics and raw log

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
func (it *ValidatorStateSetIdleTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorStateSetIdleTimeout)
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
		it.Event = new(ValidatorStateSetIdleTimeout)
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
func (it *ValidatorStateSetIdleTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorStateSetIdleTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorStateSetIdleTimeout represents a SetIdleTimeout event raised by the ValidatorState contract.
type ValidatorStateSetIdleTimeout struct {
	Addr      common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetIdleTimeout is a free log retrieval operation binding the contract event 0xc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b818.
//
// Solidity: event SetIdleTimeout(address addr, uint256 timestamp)
func (_ValidatorState *ValidatorStateFilterer) FilterSetIdleTimeout(opts *bind.FilterOpts) (*ValidatorStateSetIdleTimeoutIterator, error) {

	logs, sub, err := _ValidatorState.contract.FilterLogs(opts, "SetIdleTimeout")
	if err != nil {
		return nil, err
	}
	return &ValidatorStateSetIdleTimeoutIterator{contract: _ValidatorState.contract, event: "SetIdleTimeout", logs: logs, sub: sub}, nil
}

// WatchSetIdleTimeout is a free log subscription operation binding the contract event 0xc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b818.
//
// Solidity: event SetIdleTimeout(address addr, uint256 timestamp)
func (_ValidatorState *ValidatorStateFilterer) WatchSetIdleTimeout(opts *bind.WatchOpts, sink chan<- *ValidatorStateSetIdleTimeout) (event.Subscription, error) {

	logs, sub, err := _ValidatorState.contract.WatchLogs(opts, "SetIdleTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorStateSetIdleTimeout)
				if err := _ValidatorState.contract.UnpackLog(event, "SetIdleTimeout", log); err != nil {
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

// ParseSetIdleTimeout is a log parse operation binding the contract event 0xc530fcdc46aac627fc4102dfb818925cf48a27607235fa293bec698ebae8b818.
//
// Solidity: event SetIdleTimeout(address addr, uint256 timestamp)
func (_ValidatorState *ValidatorStateFilterer) ParseSetIdleTimeout(log types.Log) (*ValidatorStateSetIdleTimeout, error) {
	event := new(ValidatorStateSetIdleTimeout)
	if err := _ValidatorState.contract.UnpackLog(event, "SetIdleTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorStateSetPausedTimeoutIterator is returned from FilterSetPausedTimeout and is used to iterate over the raw logs and unpacked data for SetPausedTimeout events raised by the ValidatorState contract.
type ValidatorStateSetPausedTimeoutIterator struct {
	Event *ValidatorStateSetPausedTimeout // Event containing the contract specifics and raw log

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
func (it *ValidatorStateSetPausedTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorStateSetPausedTimeout)
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
		it.Event = new(ValidatorStateSetPausedTimeout)
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
func (it *ValidatorStateSetPausedTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorStateSetPausedTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorStateSetPausedTimeout represents a SetPausedTimeout event raised by the ValidatorState contract.
type ValidatorStateSetPausedTimeout struct {
	Addr      common.Address
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetPausedTimeout is a free log retrieval operation binding the contract event 0xbc1a6d4bb677fb39bdb63d1713dc004bd314de87091c5dfbba6a8d78408d0dc9.
//
// Solidity: event SetPausedTimeout(address addr, uint256 timestamp)
func (_ValidatorState *ValidatorStateFilterer) FilterSetPausedTimeout(opts *bind.FilterOpts) (*ValidatorStateSetPausedTimeoutIterator, error) {

	logs, sub, err := _ValidatorState.contract.FilterLogs(opts, "SetPausedTimeout")
	if err != nil {
		return nil, err
	}
	return &ValidatorStateSetPausedTimeoutIterator{contract: _ValidatorState.contract, event: "SetPausedTimeout", logs: logs, sub: sub}, nil
}

// WatchSetPausedTimeout is a free log subscription operation binding the contract event 0xbc1a6d4bb677fb39bdb63d1713dc004bd314de87091c5dfbba6a8d78408d0dc9.
//
// Solidity: event SetPausedTimeout(address addr, uint256 timestamp)
func (_ValidatorState *ValidatorStateFilterer) WatchSetPausedTimeout(opts *bind.WatchOpts, sink chan<- *ValidatorStateSetPausedTimeout) (event.Subscription, error) {

	logs, sub, err := _ValidatorState.contract.WatchLogs(opts, "SetPausedTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorStateSetPausedTimeout)
				if err := _ValidatorState.contract.UnpackLog(event, "SetPausedTimeout", log); err != nil {
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

// ParseSetPausedTimeout is a log parse operation binding the contract event 0xbc1a6d4bb677fb39bdb63d1713dc004bd314de87091c5dfbba6a8d78408d0dc9.
//
// Solidity: event SetPausedTimeout(address addr, uint256 timestamp)
func (_ValidatorState *ValidatorStateFilterer) ParseSetPausedTimeout(log types.Log) (*ValidatorStateSetPausedTimeout, error) {
	event := new(ValidatorStateSetPausedTimeout)
	if err := _ValidatorState.contract.UnpackLog(event, "SetPausedTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ValidatorStateStateTranstiionIterator is returned from FilterStateTranstiion and is used to iterate over the raw logs and unpacked data for StateTranstiion events raised by the ValidatorState contract.
type ValidatorStateStateTranstiionIterator struct {
	Event *ValidatorStateStateTranstiion // Event containing the contract specifics and raw log

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
func (it *ValidatorStateStateTranstiionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorStateStateTranstiion)
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
		it.Event = new(ValidatorStateStateTranstiion)
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
func (it *ValidatorStateStateTranstiionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorStateStateTranstiionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorStateStateTranstiion represents a StateTranstiion event raised by the ValidatorState contract.
type ValidatorStateStateTranstiion struct {
	OldStatw uint8
	NewState uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStateTranstiion is a free log retrieval operation binding the contract event 0x8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee36091.
//
// Solidity: event StateTranstiion(uint8 oldStatw, uint8 newState)
func (_ValidatorState *ValidatorStateFilterer) FilterStateTranstiion(opts *bind.FilterOpts) (*ValidatorStateStateTranstiionIterator, error) {

	logs, sub, err := _ValidatorState.contract.FilterLogs(opts, "StateTranstiion")
	if err != nil {
		return nil, err
	}
	return &ValidatorStateStateTranstiionIterator{contract: _ValidatorState.contract, event: "StateTranstiion", logs: logs, sub: sub}, nil
}

// WatchStateTranstiion is a free log subscription operation binding the contract event 0x8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee36091.
//
// Solidity: event StateTranstiion(uint8 oldStatw, uint8 newState)
func (_ValidatorState *ValidatorStateFilterer) WatchStateTranstiion(opts *bind.WatchOpts, sink chan<- *ValidatorStateStateTranstiion) (event.Subscription, error) {

	logs, sub, err := _ValidatorState.contract.WatchLogs(opts, "StateTranstiion")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorStateStateTranstiion)
				if err := _ValidatorState.contract.UnpackLog(event, "StateTranstiion", log); err != nil {
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

// ParseStateTranstiion is a log parse operation binding the contract event 0x8cfe2ef835f5d8ee3fe3dd2e3da396992f244536d16d4ca2fb814ae56ee36091.
//
// Solidity: event StateTranstiion(uint8 oldStatw, uint8 newState)
func (_ValidatorState *ValidatorStateFilterer) ParseStateTranstiion(log types.Log) (*ValidatorStateStateTranstiion, error) {
	event := new(ValidatorStateStateTranstiion)
	if err := _ValidatorState.contract.UnpackLog(event, "StateTranstiion", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
