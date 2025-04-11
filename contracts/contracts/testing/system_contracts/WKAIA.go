// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package system_contracts

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

// WKAIAMetaData contains all meta data concerning the WKAIA contract.
var WKAIAMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"guy\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"src\",\"type\":\"address\"},{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"dst\",\"type\":\"address\"},{\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"guy\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"dst\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"src\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"wad\",\"type\":\"uint256\"}],\"name\":\"Withdrawal\",\"type\":\"event\"}]",
	Sigs: map[string]string{
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"313ce567": "decimals()",
		"d0e30db0": "deposit()",
		"06fdde03": "name()",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"2e1a7d4d": "withdraw(uint256)",
	},
	Bin: "0x60c0604052600c60808190527f57726170706564204b616961000000000000000000000000000000000000000060a090815261003e91600091906100a3565b506040805180820190915260058082527f574b4149410000000000000000000000000000000000000000000000000000006020909201918252610083916001916100a3565b506002805460ff1916601217905534801561009d57600080fd5b5061013e565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100e457805160ff1916838001178555610111565b82800160010185558215610111579182015b828111156101115782518255916020019190600101906100f6565b5061011d929150610121565b5090565b61013b91905b8082111561011d5760008155600101610127565b90565b6106e48061014d6000396000f3fe60806040526004361061009c5760003560e01c8063313ce56711610064578063313ce5671461021157806370a082311461023c57806395d89b411461026f578063a9059cbb14610284578063d0e30db01461009c578063dd62ed3e146102bd5761009c565b806306fdde03146100a6578063095ea7b31461013057806318160ddd1461017d57806323b872dd146101a45780632e1a7d4d146101e7575b6100a46102f8565b005b3480156100b257600080fd5b506100bb610347565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100f55781810151838201526020016100dd565b50505050905090810190601f1680156101225780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561013c57600080fd5b506101696004803603604081101561015357600080fd5b506001600160a01b0381351690602001356103d5565b604080519115158252519081900360200190f35b34801561018957600080fd5b5061019261043b565b60408051918252519081900360200190f35b3480156101b057600080fd5b50610169600480360360608110156101c757600080fd5b506001600160a01b03813581169160208101359091169060400135610440565b3480156101f357600080fd5b506100a46004803603602081101561020a57600080fd5b5035610574565b34801561021d57600080fd5b50610226610609565b6040805160ff9092168252519081900360200190f35b34801561024857600080fd5b506101926004803603602081101561025f57600080fd5b50356001600160a01b0316610612565b34801561027b57600080fd5b506100bb610624565b34801561029057600080fd5b50610169600480360360408110156102a757600080fd5b506001600160a01b03813516906020013561067e565b3480156102c957600080fd5b50610192600480360360408110156102e057600080fd5b506001600160a01b0381358116916020013516610692565b33600081815260036020908152604091829020805434908101909155825190815291517fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9281900390910190a2565b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156103cd5780601f106103a2576101008083540402835291602001916103cd565b820191906000526020600020905b8154815290600101906020018083116103b057829003601f168201915b505050505081565b3360008181526004602090815260408083206001600160a01b038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b303190565b6001600160a01b03831660009081526003602052604081205482111561046557600080fd5b6001600160a01b03841633148015906104a357506001600160a01b038416600090815260046020908152604080832033845290915290205460001914155b15610503576001600160a01b03841660009081526004602090815260408083203384529091529020548211156104d857600080fd5b6001600160a01b03841660009081526004602090815260408083203384529091529020805483900390555b6001600160a01b03808516600081815260036020908152604080832080548890039055938716808352918490208054870190558351868152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35060019392505050565b3360009081526003602052604090205481111561059057600080fd5b33600081815260036020526040808220805485900390555183156108fc0291849190818181858888f193505050501580156105cf573d6000803e3d6000fd5b5060408051828152905133917f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65919081900360200190a250565b60025460ff1681565b60036020526000908152604090205481565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156103cd5780601f106103a2576101008083540402835291602001916103cd565b600061068b338484610440565b9392505050565b60046020908152600092835260408084209091529082529020548156fea265627a7a72305820a429f1117c3817e78ee7f826a6aa688db64890dbdc661379250d267160a88a7264736f6c63430005090032",
}

// WKAIAABI is the input ABI used to generate the binding from.
// Deprecated: Use WKAIAMetaData.ABI instead.
var WKAIAABI = WKAIAMetaData.ABI

// WKAIABinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const WKAIABinRuntime = `60806040526004361061009c5760003560e01c8063313ce56711610064578063313ce5671461021157806370a082311461023c57806395d89b411461026f578063a9059cbb14610284578063d0e30db01461009c578063dd62ed3e146102bd5761009c565b806306fdde03146100a6578063095ea7b31461013057806318160ddd1461017d57806323b872dd146101a45780632e1a7d4d146101e7575b6100a46102f8565b005b3480156100b257600080fd5b506100bb610347565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100f55781810151838201526020016100dd565b50505050905090810190601f1680156101225780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561013c57600080fd5b506101696004803603604081101561015357600080fd5b506001600160a01b0381351690602001356103d5565b604080519115158252519081900360200190f35b34801561018957600080fd5b5061019261043b565b60408051918252519081900360200190f35b3480156101b057600080fd5b50610169600480360360608110156101c757600080fd5b506001600160a01b03813581169160208101359091169060400135610440565b3480156101f357600080fd5b506100a46004803603602081101561020a57600080fd5b5035610574565b34801561021d57600080fd5b50610226610609565b6040805160ff9092168252519081900360200190f35b34801561024857600080fd5b506101926004803603602081101561025f57600080fd5b50356001600160a01b0316610612565b34801561027b57600080fd5b506100bb610624565b34801561029057600080fd5b50610169600480360360408110156102a757600080fd5b506001600160a01b03813516906020013561067e565b3480156102c957600080fd5b50610192600480360360408110156102e057600080fd5b506001600160a01b0381358116916020013516610692565b33600081815260036020908152604091829020805434908101909155825190815291517fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9281900390910190a2565b6000805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156103cd5780601f106103a2576101008083540402835291602001916103cd565b820191906000526020600020905b8154815290600101906020018083116103b057829003601f168201915b505050505081565b3360008181526004602090815260408083206001600160a01b038716808552908352818420869055815186815291519394909390927f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925928290030190a350600192915050565b303190565b6001600160a01b03831660009081526003602052604081205482111561046557600080fd5b6001600160a01b03841633148015906104a357506001600160a01b038416600090815260046020908152604080832033845290915290205460001914155b15610503576001600160a01b03841660009081526004602090815260408083203384529091529020548211156104d857600080fd5b6001600160a01b03841660009081526004602090815260408083203384529091529020805483900390555b6001600160a01b03808516600081815260036020908152604080832080548890039055938716808352918490208054870190558351868152935191937fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef929081900390910190a35060019392505050565b3360009081526003602052604090205481111561059057600080fd5b33600081815260036020526040808220805485900390555183156108fc0291849190818181858888f193505050501580156105cf573d6000803e3d6000fd5b5060408051828152905133917f7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65919081900360200190a250565b60025460ff1681565b60036020526000908152604090205481565b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156103cd5780601f106103a2576101008083540402835291602001916103cd565b600061068b338484610440565b9392505050565b60046020908152600092835260408084209091529082529020548156fea265627a7a72305820a429f1117c3817e78ee7f826a6aa688db64890dbdc661379250d267160a88a7264736f6c63430005090032`

// Deprecated: Use WKAIAMetaData.Sigs instead.
// WKAIAFuncSigs maps the 4-byte function signature to its string representation.
var WKAIAFuncSigs = WKAIAMetaData.Sigs

// WKAIABin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use WKAIAMetaData.Bin instead.
var WKAIABin = WKAIAMetaData.Bin

// DeployWKAIA deploys a new Kaia contract, binding an instance of WKAIA to it.
func DeployWKAIA(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *WKAIA, error) {
	parsed, err := WKAIAMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(WKAIABin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &WKAIA{WKAIACaller: WKAIACaller{contract: contract}, WKAIATransactor: WKAIATransactor{contract: contract}, WKAIAFilterer: WKAIAFilterer{contract: contract}}, nil
}

// WKAIA is an auto generated Go binding around a Kaia contract.
type WKAIA struct {
	WKAIACaller     // Read-only binding to the contract
	WKAIATransactor // Write-only binding to the contract
	WKAIAFilterer   // Log filterer for contract events
}

// WKAIACaller is an auto generated read-only Go binding around a Kaia contract.
type WKAIACaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WKAIATransactor is an auto generated write-only Go binding around a Kaia contract.
type WKAIATransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WKAIAFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type WKAIAFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WKAIASession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type WKAIASession struct {
	Contract     *WKAIA            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WKAIACallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type WKAIACallerSession struct {
	Contract *WKAIACaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// WKAIATransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type WKAIATransactorSession struct {
	Contract     *WKAIATransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WKAIARaw is an auto generated low-level Go binding around a Kaia contract.
type WKAIARaw struct {
	Contract *WKAIA // Generic contract binding to access the raw methods on
}

// WKAIACallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type WKAIACallerRaw struct {
	Contract *WKAIACaller // Generic read-only contract binding to access the raw methods on
}

// WKAIATransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type WKAIATransactorRaw struct {
	Contract *WKAIATransactor // Generic write-only contract binding to access the raw methods on
}

// NewWKAIA creates a new instance of WKAIA, bound to a specific deployed contract.
func NewWKAIA(address common.Address, backend bind.ContractBackend) (*WKAIA, error) {
	contract, err := bindWKAIA(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &WKAIA{WKAIACaller: WKAIACaller{contract: contract}, WKAIATransactor: WKAIATransactor{contract: contract}, WKAIAFilterer: WKAIAFilterer{contract: contract}}, nil
}

// NewWKAIACaller creates a new read-only instance of WKAIA, bound to a specific deployed contract.
func NewWKAIACaller(address common.Address, caller bind.ContractCaller) (*WKAIACaller, error) {
	contract, err := bindWKAIA(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WKAIACaller{contract: contract}, nil
}

// NewWKAIATransactor creates a new write-only instance of WKAIA, bound to a specific deployed contract.
func NewWKAIATransactor(address common.Address, transactor bind.ContractTransactor) (*WKAIATransactor, error) {
	contract, err := bindWKAIA(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WKAIATransactor{contract: contract}, nil
}

// NewWKAIAFilterer creates a new log filterer instance of WKAIA, bound to a specific deployed contract.
func NewWKAIAFilterer(address common.Address, filterer bind.ContractFilterer) (*WKAIAFilterer, error) {
	contract, err := bindWKAIA(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WKAIAFilterer{contract: contract}, nil
}

// bindWKAIA binds a generic wrapper to an already deployed contract.
func bindWKAIA(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := WKAIAMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WKAIA *WKAIARaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WKAIA.Contract.WKAIACaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WKAIA *WKAIARaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WKAIA.Contract.WKAIATransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WKAIA *WKAIARaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WKAIA.Contract.WKAIATransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_WKAIA *WKAIACallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _WKAIA.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_WKAIA *WKAIATransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WKAIA.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_WKAIA *WKAIATransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _WKAIA.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WKAIA *WKAIACaller) Allowance(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WKAIA.contract.Call(opts, &out, "allowance", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WKAIA *WKAIASession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WKAIA.Contract.Allowance(&_WKAIA.CallOpts, arg0, arg1)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address , address ) view returns(uint256)
func (_WKAIA *WKAIACallerSession) Allowance(arg0 common.Address, arg1 common.Address) (*big.Int, error) {
	return _WKAIA.Contract.Allowance(&_WKAIA.CallOpts, arg0, arg1)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WKAIA *WKAIACaller) BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _WKAIA.contract.Call(opts, &out, "balanceOf", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WKAIA *WKAIASession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WKAIA.Contract.BalanceOf(&_WKAIA.CallOpts, arg0)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address ) view returns(uint256)
func (_WKAIA *WKAIACallerSession) BalanceOf(arg0 common.Address) (*big.Int, error) {
	return _WKAIA.Contract.BalanceOf(&_WKAIA.CallOpts, arg0)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WKAIA *WKAIACaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _WKAIA.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WKAIA *WKAIASession) Decimals() (uint8, error) {
	return _WKAIA.Contract.Decimals(&_WKAIA.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_WKAIA *WKAIACallerSession) Decimals() (uint8, error) {
	return _WKAIA.Contract.Decimals(&_WKAIA.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WKAIA *WKAIACaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WKAIA.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WKAIA *WKAIASession) Name() (string, error) {
	return _WKAIA.Contract.Name(&_WKAIA.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_WKAIA *WKAIACallerSession) Name() (string, error) {
	return _WKAIA.Contract.Name(&_WKAIA.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WKAIA *WKAIACaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _WKAIA.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WKAIA *WKAIASession) Symbol() (string, error) {
	return _WKAIA.Contract.Symbol(&_WKAIA.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_WKAIA *WKAIACallerSession) Symbol() (string, error) {
	return _WKAIA.Contract.Symbol(&_WKAIA.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WKAIA *WKAIACaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _WKAIA.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WKAIA *WKAIASession) TotalSupply() (*big.Int, error) {
	return _WKAIA.Contract.TotalSupply(&_WKAIA.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_WKAIA *WKAIACallerSession) TotalSupply() (*big.Int, error) {
	return _WKAIA.Contract.TotalSupply(&_WKAIA.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WKAIA *WKAIATransactor) Approve(opts *bind.TransactOpts, guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.contract.Transact(opts, "approve", guy, wad)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WKAIA *WKAIASession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.Contract.Approve(&_WKAIA.TransactOpts, guy, wad)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address guy, uint256 wad) returns(bool)
func (_WKAIA *WKAIATransactorSession) Approve(guy common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.Contract.Approve(&_WKAIA.TransactOpts, guy, wad)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WKAIA *WKAIATransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _WKAIA.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WKAIA *WKAIASession) Deposit() (*types.Transaction, error) {
	return _WKAIA.Contract.Deposit(&_WKAIA.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_WKAIA *WKAIATransactorSession) Deposit() (*types.Transaction, error) {
	return _WKAIA.Contract.Deposit(&_WKAIA.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WKAIA *WKAIATransactor) Transfer(opts *bind.TransactOpts, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.contract.Transact(opts, "transfer", dst, wad)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WKAIA *WKAIASession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.Contract.Transfer(&_WKAIA.TransactOpts, dst, wad)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address dst, uint256 wad) returns(bool)
func (_WKAIA *WKAIATransactorSession) Transfer(dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.Contract.Transfer(&_WKAIA.TransactOpts, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WKAIA *WKAIATransactor) TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.contract.Transact(opts, "transferFrom", src, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WKAIA *WKAIASession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.Contract.TransferFrom(&_WKAIA.TransactOpts, src, dst, wad)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address src, address dst, uint256 wad) returns(bool)
func (_WKAIA *WKAIATransactorSession) TransferFrom(src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.Contract.TransferFrom(&_WKAIA.TransactOpts, src, dst, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WKAIA *WKAIATransactor) Withdraw(opts *bind.TransactOpts, wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.contract.Transact(opts, "withdraw", wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WKAIA *WKAIASession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.Contract.Withdraw(&_WKAIA.TransactOpts, wad)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 wad) returns()
func (_WKAIA *WKAIATransactorSession) Withdraw(wad *big.Int) (*types.Transaction, error) {
	return _WKAIA.Contract.Withdraw(&_WKAIA.TransactOpts, wad)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WKAIA *WKAIATransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _WKAIA.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WKAIA *WKAIASession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WKAIA.Contract.Fallback(&_WKAIA.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_WKAIA *WKAIATransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _WKAIA.Contract.Fallback(&_WKAIA.TransactOpts, calldata)
}

// WKAIAApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the WKAIA contract.
type WKAIAApprovalIterator struct {
	Event *WKAIAApproval // Event containing the contract specifics and raw log

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
func (it *WKAIAApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WKAIAApproval)
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
		it.Event = new(WKAIAApproval)
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
func (it *WKAIAApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WKAIAApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WKAIAApproval represents a Approval event raised by the WKAIA contract.
type WKAIAApproval struct {
	Src common.Address
	Guy common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WKAIA *WKAIAFilterer) FilterApproval(opts *bind.FilterOpts, src []common.Address, guy []common.Address) (*WKAIAApprovalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WKAIA.contract.FilterLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return &WKAIAApprovalIterator{contract: _WKAIA.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WKAIA *WKAIAFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *WKAIAApproval, src []common.Address, guy []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var guyRule []interface{}
	for _, guyItem := range guy {
		guyRule = append(guyRule, guyItem)
	}

	logs, sub, err := _WKAIA.contract.WatchLogs(opts, "Approval", srcRule, guyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WKAIAApproval)
				if err := _WKAIA.contract.UnpackLog(event, "Approval", log); err != nil {
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
// Solidity: event Approval(address indexed src, address indexed guy, uint256 wad)
func (_WKAIA *WKAIAFilterer) ParseApproval(log types.Log) (*WKAIAApproval, error) {
	event := new(WKAIAApproval)
	if err := _WKAIA.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WKAIADepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the WKAIA contract.
type WKAIADepositIterator struct {
	Event *WKAIADeposit // Event containing the contract specifics and raw log

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
func (it *WKAIADepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WKAIADeposit)
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
		it.Event = new(WKAIADeposit)
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
func (it *WKAIADepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WKAIADepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WKAIADeposit represents a Deposit event raised by the WKAIA contract.
type WKAIADeposit struct {
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WKAIA *WKAIAFilterer) FilterDeposit(opts *bind.FilterOpts, dst []common.Address) (*WKAIADepositIterator, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WKAIA.contract.FilterLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return &WKAIADepositIterator{contract: _WKAIA.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WKAIA *WKAIAFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *WKAIADeposit, dst []common.Address) (event.Subscription, error) {

	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WKAIA.contract.WatchLogs(opts, "Deposit", dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WKAIADeposit)
				if err := _WKAIA.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c.
//
// Solidity: event Deposit(address indexed dst, uint256 wad)
func (_WKAIA *WKAIAFilterer) ParseDeposit(log types.Log) (*WKAIADeposit, error) {
	event := new(WKAIADeposit)
	if err := _WKAIA.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WKAIATransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the WKAIA contract.
type WKAIATransferIterator struct {
	Event *WKAIATransfer // Event containing the contract specifics and raw log

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
func (it *WKAIATransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WKAIATransfer)
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
		it.Event = new(WKAIATransfer)
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
func (it *WKAIATransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WKAIATransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WKAIATransfer represents a Transfer event raised by the WKAIA contract.
type WKAIATransfer struct {
	Src common.Address
	Dst common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WKAIA *WKAIAFilterer) FilterTransfer(opts *bind.FilterOpts, src []common.Address, dst []common.Address) (*WKAIATransferIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WKAIA.contract.FilterLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return &WKAIATransferIterator{contract: _WKAIA.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WKAIA *WKAIAFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *WKAIATransfer, src []common.Address, dst []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}
	var dstRule []interface{}
	for _, dstItem := range dst {
		dstRule = append(dstRule, dstItem)
	}

	logs, sub, err := _WKAIA.contract.WatchLogs(opts, "Transfer", srcRule, dstRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WKAIATransfer)
				if err := _WKAIA.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(address indexed src, address indexed dst, uint256 wad)
func (_WKAIA *WKAIAFilterer) ParseTransfer(log types.Log) (*WKAIATransfer, error) {
	event := new(WKAIATransfer)
	if err := _WKAIA.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// WKAIAWithdrawalIterator is returned from FilterWithdrawal and is used to iterate over the raw logs and unpacked data for Withdrawal events raised by the WKAIA contract.
type WKAIAWithdrawalIterator struct {
	Event *WKAIAWithdrawal // Event containing the contract specifics and raw log

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
func (it *WKAIAWithdrawalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WKAIAWithdrawal)
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
		it.Event = new(WKAIAWithdrawal)
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
func (it *WKAIAWithdrawalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WKAIAWithdrawalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WKAIAWithdrawal represents a Withdrawal event raised by the WKAIA contract.
type WKAIAWithdrawal struct {
	Src common.Address
	Wad *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterWithdrawal is a free log retrieval operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WKAIA *WKAIAFilterer) FilterWithdrawal(opts *bind.FilterOpts, src []common.Address) (*WKAIAWithdrawalIterator, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WKAIA.contract.FilterLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return &WKAIAWithdrawalIterator{contract: _WKAIA.contract, event: "Withdrawal", logs: logs, sub: sub}, nil
}

// WatchWithdrawal is a free log subscription operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WKAIA *WKAIAFilterer) WatchWithdrawal(opts *bind.WatchOpts, sink chan<- *WKAIAWithdrawal, src []common.Address) (event.Subscription, error) {

	var srcRule []interface{}
	for _, srcItem := range src {
		srcRule = append(srcRule, srcItem)
	}

	logs, sub, err := _WKAIA.contract.WatchLogs(opts, "Withdrawal", srcRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WKAIAWithdrawal)
				if err := _WKAIA.contract.UnpackLog(event, "Withdrawal", log); err != nil {
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

// ParseWithdrawal is a log parse operation binding the contract event 0x7fcf532c15f0a6db0bd6d0e038bea71d30d808c7d98cb3bf7268a95bf5081b65.
//
// Solidity: event Withdrawal(address indexed src, uint256 wad)
func (_WKAIA *WKAIAFilterer) ParseWithdrawal(log types.Log) (*WKAIAWithdrawal, error) {
	event := new(WKAIAWithdrawal)
	if err := _WKAIA.contract.UnpackLog(event, "Withdrawal", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
