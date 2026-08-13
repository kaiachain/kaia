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

// CnStakingV4FactoryMetaData contains all meta data concerning the CnStakingV4Factory contract.
var CnStakingV4FactoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_cnStakingBeacon\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pdBeacon\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEAD_ADDRESS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"INITIAL_LOCKUP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cnStakingBeacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"deployCnStaking\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"cnStakingProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deployCnStakingWithPD\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_pdArgs\",\"type\":\"tuple\",\"internalType\":\"structIPublicDelegation.PDConstructorArgs\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commissionTo\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"commissionRate\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gcName\",\"type\":\"string\",\"internalType\":\"string\"}]}],\"outputs\":[{\"name\":\"cnStakingProxy\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"publicDelegationProxy\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"getDeployer\",\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isDeployedCnStaking\",\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isDeployedPublicDelegation\",\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pdBeacon\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DeployCnStakingV4\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"DeployCnStakingV4WithPD\",\"inputs\":[{\"name\":\"proxy\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"publicDelegation\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"InsufficientInitialStake\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
	Bin: "0x60c060405234801561000f575f80fd5b50604051610f59380380610f5983398101604081905261002e916100ae565b6001600160a01b0382166100555760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b03811661007c5760405163d92e233d60e01b815260040160405180910390fd5b6001600160a01b039182166080521660a0526100df565b80516001600160a01b03811681146100a9575f80fd5b919050565b5f80604083850312156100bf575f80fd5b6100c883610093565b91506100d660208401610093565b90509250929050565b60805160a051610e446101155f395f81816101fe01526102fc01525f818161014d015281816102a401526105480152610e445ff3fe608060405260043610610075575f3560e01c80630cf74c5c1461007957806333f5db27146100a35780634e6fd6c4146100c45780634ed7a764146100e6578063669d8d45146101055780637b0e7fdd1461013c5780637dfe297c1461016f5780639a429925146101b6578063aa777f4c146101ed575b5f80fd5b348015610084575f80fd5b50610090633b9aca0081565b6040519081526020015b60405180910390f35b6100b66100b13660046106fa565b610220565b60405161009a9291906107f2565b3480156100cf575f80fd5b506100d961dead81565b60405161009a919061080c565b3480156100f1575f80fd5b506100d9610100366004610820565b610518565b348015610110575f80fd5b506100d961011f366004610820565b6001600160a01b039081165f908152600260205260409020541690565b348015610147575f80fd5b506100d97f000000000000000000000000000000000000000000000000000000000000000081565b34801561017a575f80fd5b506101a6610189366004610820565b6001600160a01b03165f9081526001602052604090205460ff1690565b604051901515815260200161009a565b3480156101c1575f80fd5b506101a66101d0366004610820565b6001600160a01b03165f9081526020819052604090205460ff1690565b3480156101f8575f80fd5b506100d97f000000000000000000000000000000000000000000000000000000000000000081565b5f80633b9aca003410156102475760405163176c085f60e31b815260040160405180910390fd5b5f338560405160200161025b9291906107f2565b6040516020818303038152906040528051906020012090505f338686604051602001610289939291906108a2565b604051602081830303815290604052805190602001209050817f00000000000000000000000000000000000000000000000000000000000000006040516102cf90610664565b6102d991906108d6565b8190604051809103905ff59050801580156102f6573d5f803e3d5ffd5b509350807f000000000000000000000000000000000000000000000000000000000000000060405161032790610664565b61033191906108d6565b8190604051809103905ff590508015801561034e573d5f803e3d5ffd5b5060405163136793bd60e11b81529093506001600160a01b038416906326cf277a9061038090879089906004016108f8565b5f604051808303815f87803b158015610397575f80fd5b505af11580156103a9573d5f803e3d5ffd5b5050604051633e8f533560e01b81526001600160a01b0387169250633e8f533591506103db90899087906004016107f2565b5f604051808303815f87803b1580156103f2575f80fd5b505af1158015610404573d5f803e3d5ffd5b50506040516325fb490360e11b81526001600160a01b0386169250634bf69206915034906104389061dead9060040161080c565b5f604051808303818588803b15801561044f575f80fd5b505af1158015610461573d5f803e3d5ffd5b505050506001600160a01b038581165f818152602081815260408083208054600160ff1991821681179092558a87168086528285528386208054909216909217905584845260029092528083208054336001600160a01b031991821681179092559284529281902080549092169092179055519189169250907f87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f9061050790879061080c565b60405180910390a350509250929050565b5f80338360405160200161052d9291906107f2565b604051602081830303815290604052805190602001209050807f000000000000000000000000000000000000000000000000000000000000000060405161057390610664565b61057d91906108d6565b8190604051809103905ff590508015801561059a573d5f803e3d5ffd5b5060405163189acdbd60e31b81529092506001600160a01b0383169063c4d66de8906105ca90869060040161080c565b5f604051808303815f87803b1580156105e1575f80fd5b505af11580156105f3573d5f803e3d5ffd5b5050506001600160a01b038084165f81815260208181526040808320805460ff19166001179055600290915280822080546001600160a01b0319163317905551928716935090917f1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef9190a350919050565b6105148061092483390190565b80356001600160a01b0381168114610687575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b6040516080810167ffffffffffffffff811182821017156106c3576106c361068c565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156106f2576106f261068c565b604052919050565b5f806040838503121561070b575f80fd5b61071483610671565b915060208084013567ffffffffffffffff80821115610731575f80fd5b9085019060808288031215610744575f80fd5b61074c6106a0565b61075583610671565b8152610762848401610671565b8482015260408301356040820152606083013582811115610781575f80fd5b80840193505087601f840112610795575f80fd5b8235828111156107a7576107a761068c565b6107b9601f8201601f191686016106c9565b925080835288858286010111156107ce575f80fd5b80858501868501375f85828501015250816060820152809450505050509250929050565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b0391909116815260200190565b5f60208284031215610830575f80fd5b61083982610671565b9392505050565b5f60018060a01b0380835116845280602084015116602085015250604082015160408401526060820151608060608501528051806080860152806020830160a087015e5f60a0828701015260a0601f19601f8301168601019250505092915050565b6001600160a01b038481168252831660208201526060604082018190525f906108cd90830184610840565b95945050505050565b6001600160a01b039190911681526040602082018190525f9082015260600190565b6001600160a01b03831681526040602082018190525f9061091b90830184610840565b94935050505056fe60a060405260405161051438038061051483398101604081905261002291610331565b61002c828261003e565b506001600160a01b031660805261040d565b610047826100fb565b6040516001600160a01b038316907f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e905f90a28051156100ef576100ea826001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156100c0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906100e491906103ed565b82610209565b505050565b6100f76102aa565b5050565b806001600160a01b03163b5f0361013557604051631933b43b60e21b81526001600160a01b03821660048201526024015b60405180910390fd5b807fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b0392831617905560408051635c60da1b60e01b815290515f92841691635c60da1b9160048083019260209291908290030181865afa1580156101ae573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101d291906103ed565b9050806001600160a01b03163b5f036100f757604051634c9c8ce360e01b81526001600160a01b038216600482015260240161012c565b60605f61021684846102cb565b905080801561023757505f3d118061023757505f846001600160a01b03163b115b1561024c576102446102de565b9150506102a4565b801561027657604051639996b31560e01b81526001600160a01b038516600482015260240161012c565b3d15610289576102846102f7565b6102a2565b60405163d6bda27560e01b815260040160405180910390fd5b505b92915050565b34156102c95760405163b398979f60e01b815260040160405180910390fd5b565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b80516001600160a01b0381168114610318575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f8060408385031215610342575f80fd5b61034b83610302565b60208401519092506001600160401b0380821115610367575f80fd5b818501915085601f83011261037a575f80fd5b81518181111561038c5761038c61031d565b604051601f8201601f19908116603f011681019083821181831017156103b4576103b461031d565b816040528281528860208487010111156103cc575f80fd5b8260208601602083015e5f6020848301015280955050505050509250929050565b5f602082840312156103fd575f80fd5b61040682610302565b9392505050565b60805160f26104225f395f601d015260f25ff3fe6080604052600a600c565b005b60186014601a565b609d565b565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156076573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906098919060ba565b905090565b365f80375f80365f845af43d5f803e80801560b6573d5ff35b3d5ffd5b5f6020828403121560c9575f80fd5b81516001600160a01b038116811460de575f80fd5b939250505056fea164736f6c6343000819000aa164736f6c6343000819000a",
}

// CnStakingV4FactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use CnStakingV4FactoryMetaData.ABI instead.
var CnStakingV4FactoryABI = CnStakingV4FactoryMetaData.ABI

// CnStakingV4FactoryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const CnStakingV4FactoryBinRuntime = `608060405260043610610075575f3560e01c80630cf74c5c1461007957806333f5db27146100a35780634e6fd6c4146100c45780634ed7a764146100e6578063669d8d45146101055780637b0e7fdd1461013c5780637dfe297c1461016f5780639a429925146101b6578063aa777f4c146101ed575b5f80fd5b348015610084575f80fd5b50610090633b9aca0081565b6040519081526020015b60405180910390f35b6100b66100b13660046106fa565b610220565b60405161009a9291906107f2565b3480156100cf575f80fd5b506100d961dead81565b60405161009a919061080c565b3480156100f1575f80fd5b506100d9610100366004610820565b610518565b348015610110575f80fd5b506100d961011f366004610820565b6001600160a01b039081165f908152600260205260409020541690565b348015610147575f80fd5b506100d97f000000000000000000000000000000000000000000000000000000000000000081565b34801561017a575f80fd5b506101a6610189366004610820565b6001600160a01b03165f9081526001602052604090205460ff1690565b604051901515815260200161009a565b3480156101c1575f80fd5b506101a66101d0366004610820565b6001600160a01b03165f9081526020819052604090205460ff1690565b3480156101f8575f80fd5b506100d97f000000000000000000000000000000000000000000000000000000000000000081565b5f80633b9aca003410156102475760405163176c085f60e31b815260040160405180910390fd5b5f338560405160200161025b9291906107f2565b6040516020818303038152906040528051906020012090505f338686604051602001610289939291906108a2565b604051602081830303815290604052805190602001209050817f00000000000000000000000000000000000000000000000000000000000000006040516102cf90610664565b6102d991906108d6565b8190604051809103905ff59050801580156102f6573d5f803e3d5ffd5b509350807f000000000000000000000000000000000000000000000000000000000000000060405161032790610664565b61033191906108d6565b8190604051809103905ff590508015801561034e573d5f803e3d5ffd5b5060405163136793bd60e11b81529093506001600160a01b038416906326cf277a9061038090879089906004016108f8565b5f604051808303815f87803b158015610397575f80fd5b505af11580156103a9573d5f803e3d5ffd5b5050604051633e8f533560e01b81526001600160a01b0387169250633e8f533591506103db90899087906004016107f2565b5f604051808303815f87803b1580156103f2575f80fd5b505af1158015610404573d5f803e3d5ffd5b50506040516325fb490360e11b81526001600160a01b0386169250634bf69206915034906104389061dead9060040161080c565b5f604051808303818588803b15801561044f575f80fd5b505af1158015610461573d5f803e3d5ffd5b505050506001600160a01b038581165f818152602081815260408083208054600160ff1991821681179092558a87168086528285528386208054909216909217905584845260029092528083208054336001600160a01b031991821681179092559284529281902080549092169092179055519189169250907f87e576388f908fa28dd98074396103e586efe1e1a1c523831df5a717fc26359f9061050790879061080c565b60405180910390a350509250929050565b5f80338360405160200161052d9291906107f2565b604051602081830303815290604052805190602001209050807f000000000000000000000000000000000000000000000000000000000000000060405161057390610664565b61057d91906108d6565b8190604051809103905ff590508015801561059a573d5f803e3d5ffd5b5060405163189acdbd60e31b81529092506001600160a01b0383169063c4d66de8906105ca90869060040161080c565b5f604051808303815f87803b1580156105e1575f80fd5b505af11580156105f3573d5f803e3d5ffd5b5050506001600160a01b038084165f81815260208181526040808320805460ff19166001179055600290915280822080546001600160a01b0319163317905551928716935090917f1194ef4b0b09761863fcb550bda6ac525d646797cae7b851026bb5b4847fa6ef9190a350919050565b6105148061092483390190565b80356001600160a01b0381168114610687575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b6040516080810167ffffffffffffffff811182821017156106c3576106c361068c565b60405290565b604051601f8201601f1916810167ffffffffffffffff811182821017156106f2576106f261068c565b604052919050565b5f806040838503121561070b575f80fd5b61071483610671565b915060208084013567ffffffffffffffff80821115610731575f80fd5b9085019060808288031215610744575f80fd5b61074c6106a0565b61075583610671565b8152610762848401610671565b8482015260408301356040820152606083013582811115610781575f80fd5b80840193505087601f840112610795575f80fd5b8235828111156107a7576107a761068c565b6107b9601f8201601f191686016106c9565b925080835288858286010111156107ce575f80fd5b80858501868501375f85828501015250816060820152809450505050509250929050565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b0391909116815260200190565b5f60208284031215610830575f80fd5b61083982610671565b9392505050565b5f60018060a01b0380835116845280602084015116602085015250604082015160408401526060820151608060608501528051806080860152806020830160a087015e5f60a0828701015260a0601f19601f8301168601019250505092915050565b6001600160a01b038481168252831660208201526060604082018190525f906108cd90830184610840565b95945050505050565b6001600160a01b039190911681526040602082018190525f9082015260600190565b6001600160a01b03831681526040602082018190525f9061091b90830184610840565b94935050505056fe60a060405260405161051438038061051483398101604081905261002291610331565b61002c828261003e565b506001600160a01b031660805261040d565b610047826100fb565b6040516001600160a01b038316907f1cf3b03a6cf19fa2baba4df148e9dcabedea7f8a5c07840e207e5c089be95d3e905f90a28051156100ef576100ea826001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156100c0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906100e491906103ed565b82610209565b505050565b6100f76102aa565b5050565b806001600160a01b03163b5f0361013557604051631933b43b60e21b81526001600160a01b03821660048201526024015b60405180910390fd5b807fa3f0ad74e5423aebfd80d3ef4346578335a9a72aeaee59ff6cb3582b35133d5080546001600160a01b0319166001600160a01b0392831617905560408051635c60da1b60e01b815290515f92841691635c60da1b9160048083019260209291908290030181865afa1580156101ae573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906101d291906103ed565b9050806001600160a01b03163b5f036100f757604051634c9c8ce360e01b81526001600160a01b038216600482015260240161012c565b60605f61021684846102cb565b905080801561023757505f3d118061023757505f846001600160a01b03163b115b1561024c576102446102de565b9150506102a4565b801561027657604051639996b31560e01b81526001600160a01b038516600482015260240161012c565b3d15610289576102846102f7565b6102a2565b60405163d6bda27560e01b815260040160405180910390fd5b505b92915050565b34156102c95760405163b398979f60e01b815260040160405180910390fd5b565b5f805f835160208501865af49392505050565b6040513d81523d5f602083013e3d602001810160405290565b6040513d5f823e3d81fd5b80516001600160a01b0381168114610318575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f8060408385031215610342575f80fd5b61034b83610302565b60208401519092506001600160401b0380821115610367575f80fd5b818501915085601f83011261037a575f80fd5b81518181111561038c5761038c61031d565b604051601f8201601f19908116603f011681019083821181831017156103b4576103b461031d565b816040528281528860208487010111156103cc575f80fd5b8260208601602083015e5f6020848301015280955050505050509250929050565b5f602082840312156103fd575f80fd5b61040682610302565b9392505050565b60805160f26104225f395f601d015260f25ff3fe6080604052600a600c565b005b60186014601a565b609d565b565b5f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316635c60da1b6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156076573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906098919060ba565b905090565b365f80375f80365f845af43d5f803e80801560b6573d5ff35b3d5ffd5b5f6020828403121560c9575f80fd5b81516001600160a01b038116811460de575f80fd5b939250505056fea164736f6c6343000819000aa164736f6c6343000819000a`

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
// Solidity: function deployCnStaking(address _owner) returns(address cnStakingProxy)
func (_CnStakingV4Factory *CnStakingV4FactoryTransactor) DeployCnStaking(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _CnStakingV4Factory.contract.Transact(opts, "deployCnStaking", _owner)
}

// DeployCnStaking is a paid mutator transaction binding the contract method 0x4ed7a764.
//
// Solidity: function deployCnStaking(address _owner) returns(address cnStakingProxy)
func (_CnStakingV4Factory *CnStakingV4FactorySession) DeployCnStaking(_owner common.Address) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.DeployCnStaking(&_CnStakingV4Factory.TransactOpts, _owner)
}

// DeployCnStaking is a paid mutator transaction binding the contract method 0x4ed7a764.
//
// Solidity: function deployCnStaking(address _owner) returns(address cnStakingProxy)
func (_CnStakingV4Factory *CnStakingV4FactoryTransactorSession) DeployCnStaking(_owner common.Address) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.DeployCnStaking(&_CnStakingV4Factory.TransactOpts, _owner)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address cnStakingProxy, address publicDelegationProxy)
func (_CnStakingV4Factory *CnStakingV4FactoryTransactor) DeployCnStakingWithPD(opts *bind.TransactOpts, _owner common.Address, _pdArgs IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _CnStakingV4Factory.contract.Transact(opts, "deployCnStakingWithPD", _owner, _pdArgs)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address cnStakingProxy, address publicDelegationProxy)
func (_CnStakingV4Factory *CnStakingV4FactorySession) DeployCnStakingWithPD(_owner common.Address, _pdArgs IPublicDelegationPDConstructorArgs) (*types.Transaction, error) {
	return _CnStakingV4Factory.Contract.DeployCnStakingWithPD(&_CnStakingV4Factory.TransactOpts, _owner, _pdArgs)
}

// DeployCnStakingWithPD is a paid mutator transaction binding the contract method 0x33f5db27.
//
// Solidity: function deployCnStakingWithPD(address _owner, (address,address,uint256,string) _pdArgs) payable returns(address cnStakingProxy, address publicDelegationProxy)
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
