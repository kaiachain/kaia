// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abv2data

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

// BlsPublicKeyInfo is an auto generated low-level Go binding around an user-defined struct.
type BlsPublicKeyInfo struct {
	PublicKey []byte
	Pop       []byte
}

// IABv2DataContractInitData is an auto generated low-level Go binding around an user-defined struct.
type IABv2DataContractInitData struct {
	InitialOwner           common.Address
	ExitThreshold          *big.Int
	PauseTimeout           *big.Int
	IdleTimeout            *big.Int
	MaxValidatorCount      *big.Int
	MaxReadyCandidateCount *big.Int
	KefAddress             common.Address
	KifAddress             common.Address
	KpfAddress             common.Address
	NodeIds                []common.Address
	Infos                  []NodeInfo
}

// IPublicDelegationPDConstructorArgs is an auto generated low-level Go binding around an user-defined struct.
type IPublicDelegationPDConstructorArgs struct {
	Owner          common.Address
	CommissionTo   common.Address
	CommissionRate *big.Int
	GcName         string
}

// NodeInfo is an auto generated low-level Go binding around an user-defined struct.
type NodeInfo struct {
	Manager         common.Address
	StakingContract common.Address
	RewardAddress   common.Address
	VoterAddress    common.Address
	TimeoutAt       *big.Int
	GcId            *big.Int
	BlsInfo         BlsPublicKeyInfo
	Metadata        string
	State           uint8
}

// ABv2DataContractMetaData contains all meta data concerning the ABv2DataContract contract.
var ABv2DataContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_implementation\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"exitThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pauseTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idleTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxValidatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxReadyCandidateCount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"kefAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kifAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kpfAddress\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"nodeIds\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"voterAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gcId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structBlsPublicKeyInfo\",\"name\":\"blsInfo\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structNodeInfo[]\",\"name\":\"infos\",\"type\":\"tuple[]\"}],\"internalType\":\"structIABv2DataContract.InitData\",\"name\":\"data\",\"type\":\"tuple\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AddressAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInput\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInput\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"exitThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getInitData\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"exitThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pauseTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idleTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxValidatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxReadyCandidateCount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"kefAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kifAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kpfAddress\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"nodeIds\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"voterAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gcId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structBlsPublicKeyInfo\",\"name\":\"blsInfo\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structNodeInfo[]\",\"name\":\"infos\",\"type\":\"tuple[]\"}],\"internalType\":\"structIABv2DataContract.InitData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"idleTimeout\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kefAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kifAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"kpfAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxReadyCandidateCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"maxValidatorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseTimeout\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"9381ece0": "exitThreshold()",
		"ebe58ed7": "getInitData()",
		"195edace": "idleTimeout()",
		"5c60da1b": "implementation()",
		"29ba7bb2": "initialOwner()",
		"5ed5886f": "kefAddress()",
		"f95fa230": "kifAddress()",
		"c5cc6993": "kpfAddress()",
		"098f8520": "maxReadyCandidateCount()",
		"2ceec630": "maxValidatorCount()",
		"0190c5d0": "pauseTimeout()",
	},
	Bin: "0x6101c0604052348015610010575f80fd5b5060405161198e38038061198e83398101604081905261002f91610bc0565b6001600160a01b0382166100565760405163b4fa3fb360e01b815260040160405180910390fd5b80516001600160a01b031661007e5760405163b4fa3fb360e01b815260040160405180910390fd5b80602001515f036100a25760405163b4fa3fb360e01b815260040160405180910390fd5b80604001515f036100c65760405163b4fa3fb360e01b815260040160405180910390fd5b80606001515f036100ea5760405163b4fa3fb360e01b815260040160405180910390fd5b80608001515f0361010e5760405163b4fa3fb360e01b815260040160405180910390fd5b8060a001515f036101325760405163b4fa3fb360e01b815260040160405180910390fd5b60c08101516001600160a01b031661015d5760405163b4fa3fb360e01b815260040160405180910390fd5b60e08101516001600160a01b03166101885760405163b4fa3fb360e01b815260040160405180910390fd5b6101008101516001600160a01b03166101b45760405163b4fa3fb360e01b815260040160405180910390fd5b610120810151516101408201515181146101e15760405163b4fa3fb360e01b815260040160405180910390fd5b81608001518111156102065760405163b4fa3fb360e01b815260040160405180910390fd5b6001600160a01b0380841660809081528351821660a0908152602085015160c0908152604086015160e0908152606087015161010090815293870151610120529186015161014052850151831661016052840151821661018052830151166101a0525f5b8181101561051d575f836101200151828151811061028a5761028a610cda565b602002602001015190505f84610140015183815181106102ac576102ac610cda565b602090810291909101015180519091506001600160a01b03166102e25760405163b4fa3fb360e01b815260040160405180910390fd5b6103015f83836020015184604001518560c0015161052660201b60201c565b6001805480820182557fb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf60180546001600160a01b03199081166001600160a01b03868116919091179092556002805493840181555f5283517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace600a9094029384018054831691841691909117815560208501517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5acf85018054841691851691909117905560408501517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad085018054841691851691909117905560608501517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad185018054909316931692909217905560808301517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad283015560a08301517f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad383015560c0830151805184937f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ad4019081906104b89082610d6a565b50602082015160018201906104cd9082610d6a565b50505060e082015160088201906104e49082610d6a565b5061010082015160098201805460ff1916600183600881111561050957610509610e29565b02179055505050505080600101905061026a565b50505050610e5d565b6105328484848461054f565b61053c83836106c1565b61054885858585610771565b5050505050565b6001600160a01b038416158061056c57506001600160a01b038316155b8061057e57506001600160a01b038216155b1561059c5760405163b4fa3fb360e01b815260040160405180910390fd5b826001600160a01b0316846001600160a01b031614806105cd5750816001600160a01b0316846001600160a01b0316145b806105e95750816001600160a01b0316836001600160a01b0316145b156106075760405163b4fa3fb360e01b815260040160405180910390fd5b805151603014158061061f5750806020015151606014155b1561063d5760405163b4fa3fb360e01b815260040160405180910390fd5b805180516020909101207fc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd293148061069d575060208082015180519101207f46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21145b156106bb5760405163b4fa3fb360e01b815260040160405180910390fd5b50505050565b5f826001600160a01b031663e1a12d356040518163ffffffff1660e01b8152600401602060405180830381865afa1580156106fe573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906107229190610e3d565b90506001600160a01b0381161580159061074e5750816001600160a01b0316816001600160a01b031614155b1561076c5760405163b4fa3fb360e01b815260040160405180910390fd5b505050565b6001600160a01b0383165f9081526020859052604090205460ff16806107ae57506001600160a01b0382165f9081526020859052604090205460ff165b806107d057506001600160a01b0381165f9081526020859052604090205460ff165b156107ee576040516316a163b960e11b815260040160405180910390fd5b6001600160a01b039283165f90815260209490945260408085208054600160ff19918216811790925593851686528186208054851682179055919093168452919092208054909216179055565b80516001600160a01b0381168114610851575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b60405161012081016001600160401b038111828210171561088d5761088d610856565b60405290565b60405161016081016001600160401b038111828210171561088d5761088d610856565b604051601f8201601f191681016001600160401b03811182821017156108de576108de610856565b604052919050565b5f6001600160401b038211156108fe576108fe610856565b5060051b60200190565b5f82601f830112610917575f80fd5b8151602061092c610927836108e6565b6108b6565b8083825260208201915060208460051b87010193508684111561094d575f80fd5b602086015b84811015610970576109638161083b565b8352918301918301610952565b509695505050505050565b5f82601f83011261098a575f80fd5b81516001600160401b038111156109a3576109a3610856565b6109b6601f8201601f19166020016108b6565b8181528460208386010111156109ca575f80fd5b8160208501602083015e5f918101602001919091529392505050565b5f604082840312156109f6575f80fd5b604080519081016001600160401b038082118383101715610a1957610a19610856565b816040528293508451915080821115610a30575f80fd5b610a3c8683870161097b565b83526020850151915080821115610a51575f80fd5b50610a5e8582860161097b565b6020830152505092915050565b805160098110610851575f80fd5b5f82601f830112610a88575f80fd5b81516020610a98610927836108e6565b82815260059290921b84018101918181019086841115610ab6575f80fd5b8286015b848110156109705780516001600160401b0380821115610ad8575f80fd5b90880190610120828b03601f1901811315610af1575f80fd5b610af961086a565b610b0488850161083b565b81526040610b1381860161083b565b898301526060610b2481870161083b565b8284015260809150610b3782870161083b565b818401525060a0808601518284015260c0915081860151818401525060e08086015185811115610b65575f80fd5b610b738f8c838a01016109e6565b838501525061010091508186015185811115610b8d575f80fd5b610b9b8f8c838a010161097b565b828501525050610bac838601610a6b565b908201528652505050918301918301610aba565b5f8060408385031215610bd1575f80fd5b610bda8361083b565b60208401519092506001600160401b0380821115610bf6575f80fd5b908401906101608287031215610c0a575f80fd5b610c12610893565b610c1b8361083b565b81526020830151602082015260408301516040820152606083015160608201526080830151608082015260a083015160a0820152610c5b60c0840161083b565b60c0820152610c6c60e0840161083b565b60e0820152610100610c7f81850161083b565b908201526101208381015183811115610c96575f80fd5b610ca289828701610908565b8284015250506101408084015183811115610cbb575f80fd5b610cc789828701610a79565b8284015250508093505050509250929050565b634e487b7160e01b5f52603260045260245ffd5b600181811c90821680610d0257607f821691505b602082108103610d2057634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561076c57805f5260205f20601f840160051c81016020851015610d4b5750805b601f840160051c820191505b81811015610548575f8155600101610d57565b81516001600160401b03811115610d8357610d83610856565b610d9781610d918454610cee565b84610d26565b602080601f831160018114610dca575f8415610db35750858301515b5f19600386901b1c1916600185901b178555610e21565b5f85815260208120601f198616915b82811015610df857888601518255948401946001909101908401610dd9565b5085821015610e1557878501515f19600388901b60f8161c191681555b505060018460011b0185555b505050505050565b634e487b7160e01b5f52602160045260245ffd5b5f60208284031215610e4d575f80fd5b610e568261083b565b9392505050565b60805160a05160c05160e05161010051610120516101405161016051610180516101a051610a85610f095f395f8181610212015261044101525f818161024e015261041201525f81816101c401526103e301525f818160e901526103bd01525f8181610176015261039701525f8181610110015261037101525f818160af015261034b01525f81816101eb015261032501525f818161013701526102f601525f61019d0152610a855ff3fe608060405234801561000f575f80fd5b50600436106100a6575f3560e01c80635c60da1b1161006e5780635c60da1b146101985780635ed5886f146101bf5780639381ece0146101e6578063c5cc69931461020d578063ebe58ed714610234578063f95fa23014610249575f80fd5b80630190c5d0146100aa578063098f8520146100e4578063195edace1461010b57806329ba7bb2146101325780632ceec63014610171575b5f80fd5b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020015b60405180910390f35b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016100db565b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b61023c610270565b6040516100db9190610931565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b6102e86040518061016001604052805f6001600160a01b031681526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b0316815260200160608152602001606081525090565b6040518061016001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815260200160018054806020026020016040519081016040528092919081815260200182805480156104c257602002820191905f5260205f20905b81546001600160a01b031681526001909101906020018083116104a4575b505050505081526020016002805480602002602001604051908101604052809291908181526020015f905b82821015610765575f8481526020908190206040805161012081018252600a860290920180546001600160a01b039081168452600182015481169484019490945260028101548416838301526003810154909316606083015260048301546080830152600583015460a0830152805180820190915260068301805492939260c085019291908290829061057f90610a17565b80601f01602080910402602001604051908101604052809291908181526020018280546105ab90610a17565b80156105f65780601f106105cd576101008083540402835291602001916105f6565b820191905f5260205f20905b8154815290600101906020018083116105d957829003601f168201915b5050505050815260200160018201805461060f90610a17565b80601f016020809104026020016040519081016040528092919081815260200182805461063b90610a17565b80156106865780601f1061065d57610100808354040283529160200191610686565b820191905f5260205f20905b81548152906001019060200180831161066957829003601f168201915b50505050508152505081526020016008820180546106a390610a17565b80601f01602080910402602001604051908101604052809291908181526020018280546106cf90610a17565b801561071a5780601f106106f15761010080835404028352916020019161071a565b820191905f5260205f20905b8154815290600101906020018083116106fd57829003601f168201915b5050509183525050600982015460209091019060ff16600881111561074157610741610817565b600881111561075257610752610817565b81525050815260200190600101906104ed565b505050915250919050565b5f815180845260208085019450602084015f5b838110156107a85781516001600160a01b031687529582019590820190600101610783565b509495945050505050565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b5f8151604084526107f560408501826107b3565b90506020830151848203602086015261080e82826107b3565b95945050505050565b634e487b7160e01b5f52602160045260245ffd5b6009811061084757634e487b7160e01b5f52602160045260245ffd5b9052565b5f82825180855260208086019550808260051b8401018186015f5b8481101561092457601f19868403018952815180516001600160a01b03908116855285820151811686860152604080830151821690860152606080830151909116908501526080808201519085015260a0808201519085015260c08082015161012082870181905291906108dc838801826107e1565b9250505060e080830151868303828801526108f783826107b3565b925050506101008083015192506109108187018461082b565b509985019993505090830190600101610866565b5090979650505050505050565b6020815261094b6020820183516001600160a01b03169052565b602082015160408201526040820151606082015260608201516080820152608082015160a082015260a082015160c08201525f60c083015161099860e08401826001600160a01b03169052565b5060e08301516101006109b5818501836001600160a01b03169052565b84015190506101206109d1848201836001600160a01b03169052565b8085015191505061016061014081818601526109f1610180860184610770565b90860151858203601f190183870152909250610a0d838261084b565b9695505050505050565b600181811c90821680610a2b57607f821691505b602082108103610a4957634e487b7160e01b5f52602260045260245ffd5b5091905056fea2646970667358221220cad7c153543400f5ee5b8a6560c8fae7f6e2cc26842630820871ca539cd1202164736f6c63430008190033",
}

// ABv2DataContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ABv2DataContractMetaData.ABI instead.
var ABv2DataContractABI = ABv2DataContractMetaData.ABI

// ABv2DataContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const ABv2DataContractBinRuntime = `608060405234801561000f575f80fd5b50600436106100a6575f3560e01c80635c60da1b1161006e5780635c60da1b146101985780635ed5886f146101bf5780639381ece0146101e6578063c5cc69931461020d578063ebe58ed714610234578063f95fa23014610249575f80fd5b80630190c5d0146100aa578063098f8520146100e4578063195edace1461010b57806329ba7bb2146101325780632ceec63014610171575b5f80fd5b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6040519081526020015b60405180910390f35b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b0390911681526020016100db565b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b6100d17f000000000000000000000000000000000000000000000000000000000000000081565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b61023c610270565b6040516100db9190610931565b6101597f000000000000000000000000000000000000000000000000000000000000000081565b6102e86040518061016001604052805f6001600160a01b031681526020015f81526020015f81526020015f81526020015f81526020015f81526020015f6001600160a01b031681526020015f6001600160a01b031681526020015f6001600160a01b0316815260200160608152602001606081525090565b6040518061016001604052807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f000000000000000000000000000000000000000000000000000000000000000081526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031681526020017f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316815260200160018054806020026020016040519081016040528092919081815260200182805480156104c257602002820191905f5260205f20905b81546001600160a01b031681526001909101906020018083116104a4575b505050505081526020016002805480602002602001604051908101604052809291908181526020015f905b82821015610765575f8481526020908190206040805161012081018252600a860290920180546001600160a01b039081168452600182015481169484019490945260028101548416838301526003810154909316606083015260048301546080830152600583015460a0830152805180820190915260068301805492939260c085019291908290829061057f90610a17565b80601f01602080910402602001604051908101604052809291908181526020018280546105ab90610a17565b80156105f65780601f106105cd576101008083540402835291602001916105f6565b820191905f5260205f20905b8154815290600101906020018083116105d957829003601f168201915b5050505050815260200160018201805461060f90610a17565b80601f016020809104026020016040519081016040528092919081815260200182805461063b90610a17565b80156106865780601f1061065d57610100808354040283529160200191610686565b820191905f5260205f20905b81548152906001019060200180831161066957829003601f168201915b50505050508152505081526020016008820180546106a390610a17565b80601f01602080910402602001604051908101604052809291908181526020018280546106cf90610a17565b801561071a5780601f106106f15761010080835404028352916020019161071a565b820191905f5260205f20905b8154815290600101906020018083116106fd57829003601f168201915b5050509183525050600982015460209091019060ff16600881111561074157610741610817565b600881111561075257610752610817565b81525050815260200190600101906104ed565b505050915250919050565b5f815180845260208085019450602084015f5b838110156107a85781516001600160a01b031687529582019590820190600101610783565b509495945050505050565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b5f8151604084526107f560408501826107b3565b90506020830151848203602086015261080e82826107b3565b95945050505050565b634e487b7160e01b5f52602160045260245ffd5b6009811061084757634e487b7160e01b5f52602160045260245ffd5b9052565b5f82825180855260208086019550808260051b8401018186015f5b8481101561092457601f19868403018952815180516001600160a01b03908116855285820151811686860152604080830151821690860152606080830151909116908501526080808201519085015260a0808201519085015260c08082015161012082870181905291906108dc838801826107e1565b9250505060e080830151868303828801526108f783826107b3565b925050506101008083015192506109108187018461082b565b509985019993505090830190600101610866565b5090979650505050505050565b6020815261094b6020820183516001600160a01b03169052565b602082015160408201526040820151606082015260608201516080820152608082015160a082015260a082015160c08201525f60c083015161099860e08401826001600160a01b03169052565b5060e08301516101006109b5818501836001600160a01b03169052565b84015190506101206109d1848201836001600160a01b03169052565b8085015191505061016061014081818601526109f1610180860184610770565b90860151858203601f190183870152909250610a0d838261084b565b9695505050505050565b600181811c90821680610a2b57607f821691505b602082108103610a4957634e487b7160e01b5f52602260045260245ffd5b5091905056fea2646970667358221220cad7c153543400f5ee5b8a6560c8fae7f6e2cc26842630820871ca539cd1202164736f6c63430008190033`

// Deprecated: Use ABv2DataContractMetaData.Sigs instead.
// ABv2DataContractFuncSigs maps the 4-byte function signature to its string representation.
var ABv2DataContractFuncSigs = ABv2DataContractMetaData.Sigs

// ABv2DataContractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ABv2DataContractMetaData.Bin instead.
var ABv2DataContractBin = ABv2DataContractMetaData.Bin

// DeployABv2DataContract deploys a new Kaia contract, binding an instance of ABv2DataContract to it.
func DeployABv2DataContract(auth *bind.TransactOpts, backend bind.ContractBackend, _implementation common.Address, data IABv2DataContractInitData) (common.Address, *types.Transaction, *ABv2DataContract, error) {
	parsed, err := ABv2DataContractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ABv2DataContractBin), backend, _implementation, data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ABv2DataContract{ABv2DataContractCaller: ABv2DataContractCaller{contract: contract}, ABv2DataContractTransactor: ABv2DataContractTransactor{contract: contract}, ABv2DataContractFilterer: ABv2DataContractFilterer{contract: contract}}, nil
}

// ABv2DataContract is an auto generated Go binding around a Kaia contract.
type ABv2DataContract struct {
	ABv2DataContractCaller     // Read-only binding to the contract
	ABv2DataContractTransactor // Write-only binding to the contract
	ABv2DataContractFilterer   // Log filterer for contract events
}

// ABv2DataContractCaller is an auto generated read-only Go binding around a Kaia contract.
type ABv2DataContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ABv2DataContractTransactor is an auto generated write-only Go binding around a Kaia contract.
type ABv2DataContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ABv2DataContractFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type ABv2DataContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ABv2DataContractSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type ABv2DataContractSession struct {
	Contract     *ABv2DataContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ABv2DataContractCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type ABv2DataContractCallerSession struct {
	Contract *ABv2DataContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ABv2DataContractTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type ABv2DataContractTransactorSession struct {
	Contract     *ABv2DataContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ABv2DataContractRaw is an auto generated low-level Go binding around a Kaia contract.
type ABv2DataContractRaw struct {
	Contract *ABv2DataContract // Generic contract binding to access the raw methods on
}

// ABv2DataContractCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type ABv2DataContractCallerRaw struct {
	Contract *ABv2DataContractCaller // Generic read-only contract binding to access the raw methods on
}

// ABv2DataContractTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type ABv2DataContractTransactorRaw struct {
	Contract *ABv2DataContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewABv2DataContract creates a new instance of ABv2DataContract, bound to a specific deployed contract.
func NewABv2DataContract(address common.Address, backend bind.ContractBackend) (*ABv2DataContract, error) {
	contract, err := bindABv2DataContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ABv2DataContract{ABv2DataContractCaller: ABv2DataContractCaller{contract: contract}, ABv2DataContractTransactor: ABv2DataContractTransactor{contract: contract}, ABv2DataContractFilterer: ABv2DataContractFilterer{contract: contract}}, nil
}

// NewABv2DataContractCaller creates a new read-only instance of ABv2DataContract, bound to a specific deployed contract.
func NewABv2DataContractCaller(address common.Address, caller bind.ContractCaller) (*ABv2DataContractCaller, error) {
	contract, err := bindABv2DataContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ABv2DataContractCaller{contract: contract}, nil
}

// NewABv2DataContractTransactor creates a new write-only instance of ABv2DataContract, bound to a specific deployed contract.
func NewABv2DataContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ABv2DataContractTransactor, error) {
	contract, err := bindABv2DataContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ABv2DataContractTransactor{contract: contract}, nil
}

// NewABv2DataContractFilterer creates a new log filterer instance of ABv2DataContract, bound to a specific deployed contract.
func NewABv2DataContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ABv2DataContractFilterer, error) {
	contract, err := bindABv2DataContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ABv2DataContractFilterer{contract: contract}, nil
}

// bindABv2DataContract binds a generic wrapper to an already deployed contract.
func bindABv2DataContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ABv2DataContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ABv2DataContract *ABv2DataContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ABv2DataContract.Contract.ABv2DataContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ABv2DataContract *ABv2DataContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ABv2DataContract.Contract.ABv2DataContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ABv2DataContract *ABv2DataContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ABv2DataContract.Contract.ABv2DataContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ABv2DataContract *ABv2DataContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ABv2DataContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ABv2DataContract *ABv2DataContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ABv2DataContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ABv2DataContract *ABv2DataContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ABv2DataContract.Contract.contract.Transact(opts, method, params...)
}

// ExitThreshold is a free data retrieval call binding the contract method 0x9381ece0.
//
// Solidity: function exitThreshold() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCaller) ExitThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "exitThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExitThreshold is a free data retrieval call binding the contract method 0x9381ece0.
//
// Solidity: function exitThreshold() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractSession) ExitThreshold() (*big.Int, error) {
	return _ABv2DataContract.Contract.ExitThreshold(&_ABv2DataContract.CallOpts)
}

// ExitThreshold is a free data retrieval call binding the contract method 0x9381ece0.
//
// Solidity: function exitThreshold() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCallerSession) ExitThreshold() (*big.Int, error) {
	return _ABv2DataContract.Contract.ExitThreshold(&_ABv2DataContract.CallOpts)
}

// GetInitData is a free data retrieval call binding the contract method 0xebe58ed7.
//
// Solidity: function getInitData() view returns((address,uint256,uint256,uint256,uint256,uint256,address,address,address,address[],(address,address,address,address,uint256,uint256,(bytes,bytes),string,uint8)[]))
func (_ABv2DataContract *ABv2DataContractCaller) GetInitData(opts *bind.CallOpts) (IABv2DataContractInitData, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "getInitData")

	if err != nil {
		return *new(IABv2DataContractInitData), err
	}

	out0 := *abi.ConvertType(out[0], new(IABv2DataContractInitData)).(*IABv2DataContractInitData)

	return out0, err

}

// GetInitData is a free data retrieval call binding the contract method 0xebe58ed7.
//
// Solidity: function getInitData() view returns((address,uint256,uint256,uint256,uint256,uint256,address,address,address,address[],(address,address,address,address,uint256,uint256,(bytes,bytes),string,uint8)[]))
func (_ABv2DataContract *ABv2DataContractSession) GetInitData() (IABv2DataContractInitData, error) {
	return _ABv2DataContract.Contract.GetInitData(&_ABv2DataContract.CallOpts)
}

// GetInitData is a free data retrieval call binding the contract method 0xebe58ed7.
//
// Solidity: function getInitData() view returns((address,uint256,uint256,uint256,uint256,uint256,address,address,address,address[],(address,address,address,address,uint256,uint256,(bytes,bytes),string,uint8)[]))
func (_ABv2DataContract *ABv2DataContractCallerSession) GetInitData() (IABv2DataContractInitData, error) {
	return _ABv2DataContract.Contract.GetInitData(&_ABv2DataContract.CallOpts)
}

// IdleTimeout is a free data retrieval call binding the contract method 0x195edace.
//
// Solidity: function idleTimeout() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCaller) IdleTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "idleTimeout")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// IdleTimeout is a free data retrieval call binding the contract method 0x195edace.
//
// Solidity: function idleTimeout() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractSession) IdleTimeout() (*big.Int, error) {
	return _ABv2DataContract.Contract.IdleTimeout(&_ABv2DataContract.CallOpts)
}

// IdleTimeout is a free data retrieval call binding the contract method 0x195edace.
//
// Solidity: function idleTimeout() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCallerSession) IdleTimeout() (*big.Int, error) {
	return _ABv2DataContract.Contract.IdleTimeout(&_ABv2DataContract.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_ABv2DataContract *ABv2DataContractCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_ABv2DataContract *ABv2DataContractSession) Implementation() (common.Address, error) {
	return _ABv2DataContract.Contract.Implementation(&_ABv2DataContract.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_ABv2DataContract *ABv2DataContractCallerSession) Implementation() (common.Address, error) {
	return _ABv2DataContract.Contract.Implementation(&_ABv2DataContract.CallOpts)
}

// InitialOwner is a free data retrieval call binding the contract method 0x29ba7bb2.
//
// Solidity: function initialOwner() view returns(address)
func (_ABv2DataContract *ABv2DataContractCaller) InitialOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "initialOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// InitialOwner is a free data retrieval call binding the contract method 0x29ba7bb2.
//
// Solidity: function initialOwner() view returns(address)
func (_ABv2DataContract *ABv2DataContractSession) InitialOwner() (common.Address, error) {
	return _ABv2DataContract.Contract.InitialOwner(&_ABv2DataContract.CallOpts)
}

// InitialOwner is a free data retrieval call binding the contract method 0x29ba7bb2.
//
// Solidity: function initialOwner() view returns(address)
func (_ABv2DataContract *ABv2DataContractCallerSession) InitialOwner() (common.Address, error) {
	return _ABv2DataContract.Contract.InitialOwner(&_ABv2DataContract.CallOpts)
}

// KefAddress is a free data retrieval call binding the contract method 0x5ed5886f.
//
// Solidity: function kefAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractCaller) KefAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "kefAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KefAddress is a free data retrieval call binding the contract method 0x5ed5886f.
//
// Solidity: function kefAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractSession) KefAddress() (common.Address, error) {
	return _ABv2DataContract.Contract.KefAddress(&_ABv2DataContract.CallOpts)
}

// KefAddress is a free data retrieval call binding the contract method 0x5ed5886f.
//
// Solidity: function kefAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractCallerSession) KefAddress() (common.Address, error) {
	return _ABv2DataContract.Contract.KefAddress(&_ABv2DataContract.CallOpts)
}

// KifAddress is a free data retrieval call binding the contract method 0xf95fa230.
//
// Solidity: function kifAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractCaller) KifAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "kifAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KifAddress is a free data retrieval call binding the contract method 0xf95fa230.
//
// Solidity: function kifAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractSession) KifAddress() (common.Address, error) {
	return _ABv2DataContract.Contract.KifAddress(&_ABv2DataContract.CallOpts)
}

// KifAddress is a free data retrieval call binding the contract method 0xf95fa230.
//
// Solidity: function kifAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractCallerSession) KifAddress() (common.Address, error) {
	return _ABv2DataContract.Contract.KifAddress(&_ABv2DataContract.CallOpts)
}

// KpfAddress is a free data retrieval call binding the contract method 0xc5cc6993.
//
// Solidity: function kpfAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractCaller) KpfAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "kpfAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KpfAddress is a free data retrieval call binding the contract method 0xc5cc6993.
//
// Solidity: function kpfAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractSession) KpfAddress() (common.Address, error) {
	return _ABv2DataContract.Contract.KpfAddress(&_ABv2DataContract.CallOpts)
}

// KpfAddress is a free data retrieval call binding the contract method 0xc5cc6993.
//
// Solidity: function kpfAddress() view returns(address)
func (_ABv2DataContract *ABv2DataContractCallerSession) KpfAddress() (common.Address, error) {
	return _ABv2DataContract.Contract.KpfAddress(&_ABv2DataContract.CallOpts)
}

// MaxReadyCandidateCount is a free data retrieval call binding the contract method 0x098f8520.
//
// Solidity: function maxReadyCandidateCount() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCaller) MaxReadyCandidateCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "maxReadyCandidateCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxReadyCandidateCount is a free data retrieval call binding the contract method 0x098f8520.
//
// Solidity: function maxReadyCandidateCount() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractSession) MaxReadyCandidateCount() (*big.Int, error) {
	return _ABv2DataContract.Contract.MaxReadyCandidateCount(&_ABv2DataContract.CallOpts)
}

// MaxReadyCandidateCount is a free data retrieval call binding the contract method 0x098f8520.
//
// Solidity: function maxReadyCandidateCount() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCallerSession) MaxReadyCandidateCount() (*big.Int, error) {
	return _ABv2DataContract.Contract.MaxReadyCandidateCount(&_ABv2DataContract.CallOpts)
}

// MaxValidatorCount is a free data retrieval call binding the contract method 0x2ceec630.
//
// Solidity: function maxValidatorCount() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCaller) MaxValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "maxValidatorCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaxValidatorCount is a free data retrieval call binding the contract method 0x2ceec630.
//
// Solidity: function maxValidatorCount() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractSession) MaxValidatorCount() (*big.Int, error) {
	return _ABv2DataContract.Contract.MaxValidatorCount(&_ABv2DataContract.CallOpts)
}

// MaxValidatorCount is a free data retrieval call binding the contract method 0x2ceec630.
//
// Solidity: function maxValidatorCount() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCallerSession) MaxValidatorCount() (*big.Int, error) {
	return _ABv2DataContract.Contract.MaxValidatorCount(&_ABv2DataContract.CallOpts)
}

// PauseTimeout is a free data retrieval call binding the contract method 0x0190c5d0.
//
// Solidity: function pauseTimeout() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCaller) PauseTimeout(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ABv2DataContract.contract.Call(opts, &out, "pauseTimeout")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PauseTimeout is a free data retrieval call binding the contract method 0x0190c5d0.
//
// Solidity: function pauseTimeout() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractSession) PauseTimeout() (*big.Int, error) {
	return _ABv2DataContract.Contract.PauseTimeout(&_ABv2DataContract.CallOpts)
}

// PauseTimeout is a free data retrieval call binding the contract method 0x0190c5d0.
//
// Solidity: function pauseTimeout() view returns(uint256)
func (_ABv2DataContract *ABv2DataContractCallerSession) PauseTimeout() (*big.Int, error) {
	return _ABv2DataContract.Contract.PauseTimeout(&_ABv2DataContract.CallOpts)
}

// IABv2DataContractMetaData contains all meta data concerning the IABv2DataContract contract.
var IABv2DataContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getInitData\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"exitThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pauseTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"idleTimeout\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxValidatorCount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxReadyCandidateCount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"kefAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kifAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"kpfAddress\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"nodeIds\",\"type\":\"address[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stakingContract\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"rewardAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"voterAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeoutAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"gcId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"publicKey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"pop\",\"type\":\"bytes\"}],\"internalType\":\"structBlsPublicKeyInfo\",\"name\":\"blsInfo\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"metadata\",\"type\":\"string\"},{\"internalType\":\"enumState\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structNodeInfo[]\",\"name\":\"infos\",\"type\":\"tuple[]\"}],\"internalType\":\"structIABv2DataContract.InitData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"implementation\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"ebe58ed7": "getInitData()",
		"5c60da1b": "implementation()",
	},
}

// IABv2DataContractABI is the input ABI used to generate the binding from.
// Deprecated: Use IABv2DataContractMetaData.ABI instead.
var IABv2DataContractABI = IABv2DataContractMetaData.ABI

// IABv2DataContractBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IABv2DataContractBinRuntime = ``

// Deprecated: Use IABv2DataContractMetaData.Sigs instead.
// IABv2DataContractFuncSigs maps the 4-byte function signature to its string representation.
var IABv2DataContractFuncSigs = IABv2DataContractMetaData.Sigs

// IABv2DataContract is an auto generated Go binding around a Kaia contract.
type IABv2DataContract struct {
	IABv2DataContractCaller     // Read-only binding to the contract
	IABv2DataContractTransactor // Write-only binding to the contract
	IABv2DataContractFilterer   // Log filterer for contract events
}

// IABv2DataContractCaller is an auto generated read-only Go binding around a Kaia contract.
type IABv2DataContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IABv2DataContractTransactor is an auto generated write-only Go binding around a Kaia contract.
type IABv2DataContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IABv2DataContractFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IABv2DataContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IABv2DataContractSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IABv2DataContractSession struct {
	Contract     *IABv2DataContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// IABv2DataContractCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IABv2DataContractCallerSession struct {
	Contract *IABv2DataContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// IABv2DataContractTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IABv2DataContractTransactorSession struct {
	Contract     *IABv2DataContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// IABv2DataContractRaw is an auto generated low-level Go binding around a Kaia contract.
type IABv2DataContractRaw struct {
	Contract *IABv2DataContract // Generic contract binding to access the raw methods on
}

// IABv2DataContractCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IABv2DataContractCallerRaw struct {
	Contract *IABv2DataContractCaller // Generic read-only contract binding to access the raw methods on
}

// IABv2DataContractTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IABv2DataContractTransactorRaw struct {
	Contract *IABv2DataContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIABv2DataContract creates a new instance of IABv2DataContract, bound to a specific deployed contract.
func NewIABv2DataContract(address common.Address, backend bind.ContractBackend) (*IABv2DataContract, error) {
	contract, err := bindIABv2DataContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IABv2DataContract{IABv2DataContractCaller: IABv2DataContractCaller{contract: contract}, IABv2DataContractTransactor: IABv2DataContractTransactor{contract: contract}, IABv2DataContractFilterer: IABv2DataContractFilterer{contract: contract}}, nil
}

// NewIABv2DataContractCaller creates a new read-only instance of IABv2DataContract, bound to a specific deployed contract.
func NewIABv2DataContractCaller(address common.Address, caller bind.ContractCaller) (*IABv2DataContractCaller, error) {
	contract, err := bindIABv2DataContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IABv2DataContractCaller{contract: contract}, nil
}

// NewIABv2DataContractTransactor creates a new write-only instance of IABv2DataContract, bound to a specific deployed contract.
func NewIABv2DataContractTransactor(address common.Address, transactor bind.ContractTransactor) (*IABv2DataContractTransactor, error) {
	contract, err := bindIABv2DataContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IABv2DataContractTransactor{contract: contract}, nil
}

// NewIABv2DataContractFilterer creates a new log filterer instance of IABv2DataContract, bound to a specific deployed contract.
func NewIABv2DataContractFilterer(address common.Address, filterer bind.ContractFilterer) (*IABv2DataContractFilterer, error) {
	contract, err := bindIABv2DataContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IABv2DataContractFilterer{contract: contract}, nil
}

// bindIABv2DataContract binds a generic wrapper to an already deployed contract.
func bindIABv2DataContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IABv2DataContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IABv2DataContract *IABv2DataContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IABv2DataContract.Contract.IABv2DataContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IABv2DataContract *IABv2DataContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IABv2DataContract.Contract.IABv2DataContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IABv2DataContract *IABv2DataContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IABv2DataContract.Contract.IABv2DataContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IABv2DataContract *IABv2DataContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IABv2DataContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IABv2DataContract *IABv2DataContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IABv2DataContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IABv2DataContract *IABv2DataContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IABv2DataContract.Contract.contract.Transact(opts, method, params...)
}

// GetInitData is a free data retrieval call binding the contract method 0xebe58ed7.
//
// Solidity: function getInitData() view returns((address,uint256,uint256,uint256,uint256,uint256,address,address,address,address[],(address,address,address,address,uint256,uint256,(bytes,bytes),string,uint8)[]))
func (_IABv2DataContract *IABv2DataContractCaller) GetInitData(opts *bind.CallOpts) (IABv2DataContractInitData, error) {
	var out []interface{}
	err := _IABv2DataContract.contract.Call(opts, &out, "getInitData")

	if err != nil {
		return *new(IABv2DataContractInitData), err
	}

	out0 := *abi.ConvertType(out[0], new(IABv2DataContractInitData)).(*IABv2DataContractInitData)

	return out0, err

}

// GetInitData is a free data retrieval call binding the contract method 0xebe58ed7.
//
// Solidity: function getInitData() view returns((address,uint256,uint256,uint256,uint256,uint256,address,address,address,address[],(address,address,address,address,uint256,uint256,(bytes,bytes),string,uint8)[]))
func (_IABv2DataContract *IABv2DataContractSession) GetInitData() (IABv2DataContractInitData, error) {
	return _IABv2DataContract.Contract.GetInitData(&_IABv2DataContract.CallOpts)
}

// GetInitData is a free data retrieval call binding the contract method 0xebe58ed7.
//
// Solidity: function getInitData() view returns((address,uint256,uint256,uint256,uint256,uint256,address,address,address,address[],(address,address,address,address,uint256,uint256,(bytes,bytes),string,uint8)[]))
func (_IABv2DataContract *IABv2DataContractCallerSession) GetInitData() (IABv2DataContractInitData, error) {
	return _IABv2DataContract.Contract.GetInitData(&_IABv2DataContract.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IABv2DataContract *IABv2DataContractCaller) Implementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _IABv2DataContract.contract.Call(opts, &out, "implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IABv2DataContract *IABv2DataContractSession) Implementation() (common.Address, error) {
	return _IABv2DataContract.Contract.Implementation(&_IABv2DataContract.CallOpts)
}

// Implementation is a free data retrieval call binding the contract method 0x5c60da1b.
//
// Solidity: function implementation() view returns(address)
func (_IABv2DataContract *IABv2DataContractCallerSession) Implementation() (common.Address, error) {
	return _IABv2DataContract.Contract.Implementation(&_IABv2DataContract.CallOpts)
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

// IRegistryMetaData contains all meta data concerning the IRegistry contract.
var IRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"}],\"name\":\"getActiveAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Sigs: map[string]string{
		"e2693e3f": "getActiveAddr(string)",
	},
}

// IRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use IRegistryMetaData.ABI instead.
var IRegistryABI = IRegistryMetaData.ABI

// IRegistryBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const IRegistryBinRuntime = ``

// Deprecated: Use IRegistryMetaData.Sigs instead.
// IRegistryFuncSigs maps the 4-byte function signature to its string representation.
var IRegistryFuncSigs = IRegistryMetaData.Sigs

// IRegistry is an auto generated Go binding around a Kaia contract.
type IRegistry struct {
	IRegistryCaller     // Read-only binding to the contract
	IRegistryTransactor // Write-only binding to the contract
	IRegistryFilterer   // Log filterer for contract events
}

// IRegistryCaller is an auto generated read-only Go binding around a Kaia contract.
type IRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistryTransactor is an auto generated write-only Go binding around a Kaia contract.
type IRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistryFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type IRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IRegistrySession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type IRegistrySession struct {
	Contract     *IRegistry        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// IRegistryCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type IRegistryCallerSession struct {
	Contract *IRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// IRegistryTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type IRegistryTransactorSession struct {
	Contract     *IRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// IRegistryRaw is an auto generated low-level Go binding around a Kaia contract.
type IRegistryRaw struct {
	Contract *IRegistry // Generic contract binding to access the raw methods on
}

// IRegistryCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type IRegistryCallerRaw struct {
	Contract *IRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// IRegistryTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type IRegistryTransactorRaw struct {
	Contract *IRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIRegistry creates a new instance of IRegistry, bound to a specific deployed contract.
func NewIRegistry(address common.Address, backend bind.ContractBackend) (*IRegistry, error) {
	contract, err := bindIRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IRegistry{IRegistryCaller: IRegistryCaller{contract: contract}, IRegistryTransactor: IRegistryTransactor{contract: contract}, IRegistryFilterer: IRegistryFilterer{contract: contract}}, nil
}

// NewIRegistryCaller creates a new read-only instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryCaller(address common.Address, caller bind.ContractCaller) (*IRegistryCaller, error) {
	contract, err := bindIRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IRegistryCaller{contract: contract}, nil
}

// NewIRegistryTransactor creates a new write-only instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*IRegistryTransactor, error) {
	contract, err := bindIRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IRegistryTransactor{contract: contract}, nil
}

// NewIRegistryFilterer creates a new log filterer instance of IRegistry, bound to a specific deployed contract.
func NewIRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*IRegistryFilterer, error) {
	contract, err := bindIRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IRegistryFilterer{contract: contract}, nil
}

// bindIRegistry binds a generic wrapper to an already deployed contract.
func bindIRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IRegistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRegistry *IRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IRegistry.Contract.IRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRegistry *IRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRegistry.Contract.IRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRegistry *IRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRegistry.Contract.IRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IRegistry *IRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IRegistry *IRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IRegistry *IRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistryCaller) GetActiveAddr(opts *bind.CallOpts, name string) (common.Address, error) {
	var out []interface{}
	err := _IRegistry.contract.Call(opts, &out, "getActiveAddr", name)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistrySession) GetActiveAddr(name string) (common.Address, error) {
	return _IRegistry.Contract.GetActiveAddr(&_IRegistry.CallOpts, name)
}

// GetActiveAddr is a free data retrieval call binding the contract method 0xe2693e3f.
//
// Solidity: function getActiveAddr(string name) view returns(address)
func (_IRegistry *IRegistryCallerSession) GetActiveAddr(name string) (common.Address, error) {
	return _IRegistry.Contract.GetActiveAddr(&_IRegistry.CallOpts, name)
}

// NodeVerifierMetaData contains all meta data concerning the NodeVerifier contract.
var NodeVerifierMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"AddressAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FactoryNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInput\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StakingDeployerMismatch\",\"type\":\"error\"}]",
	Bin: "0x60556032600b8282823980515f1a607314602657634e487b7160e01b5f525f60045260245ffd5b305f52607381538281f3fe730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220d46e172e68b60b46429329f39b98fd027d8a510f577c41c1c86ef1dd91485a0164736f6c63430008190033",
}

// NodeVerifierABI is the input ABI used to generate the binding from.
// Deprecated: Use NodeVerifierMetaData.ABI instead.
var NodeVerifierABI = NodeVerifierMetaData.ABI

// NodeVerifierBinRuntime is the compiled bytecode used for adding genesis block without deploying code.
const NodeVerifierBinRuntime = `730000000000000000000000000000000000000000301460806040525f80fdfea2646970667358221220d46e172e68b60b46429329f39b98fd027d8a510f577c41c1c86ef1dd91485a0164736f6c63430008190033`

// NodeVerifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NodeVerifierMetaData.Bin instead.
var NodeVerifierBin = NodeVerifierMetaData.Bin

// DeployNodeVerifier deploys a new Kaia contract, binding an instance of NodeVerifier to it.
func DeployNodeVerifier(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NodeVerifier, error) {
	parsed, err := NodeVerifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NodeVerifierBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NodeVerifier{NodeVerifierCaller: NodeVerifierCaller{contract: contract}, NodeVerifierTransactor: NodeVerifierTransactor{contract: contract}, NodeVerifierFilterer: NodeVerifierFilterer{contract: contract}}, nil
}

// NodeVerifier is an auto generated Go binding around a Kaia contract.
type NodeVerifier struct {
	NodeVerifierCaller     // Read-only binding to the contract
	NodeVerifierTransactor // Write-only binding to the contract
	NodeVerifierFilterer   // Log filterer for contract events
}

// NodeVerifierCaller is an auto generated read-only Go binding around a Kaia contract.
type NodeVerifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeVerifierTransactor is an auto generated write-only Go binding around a Kaia contract.
type NodeVerifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeVerifierFilterer is an auto generated log filtering Go binding around a Kaia contract events.
type NodeVerifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NodeVerifierSession is an auto generated Go binding around a Kaia contract,
// with pre-set call and transact options.
type NodeVerifierSession struct {
	Contract     *NodeVerifier     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NodeVerifierCallerSession is an auto generated read-only Go binding around a Kaia contract,
// with pre-set call options.
type NodeVerifierCallerSession struct {
	Contract *NodeVerifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// NodeVerifierTransactorSession is an auto generated write-only Go binding around a Kaia contract,
// with pre-set transact options.
type NodeVerifierTransactorSession struct {
	Contract     *NodeVerifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// NodeVerifierRaw is an auto generated low-level Go binding around a Kaia contract.
type NodeVerifierRaw struct {
	Contract *NodeVerifier // Generic contract binding to access the raw methods on
}

// NodeVerifierCallerRaw is an auto generated low-level read-only Go binding around a Kaia contract.
type NodeVerifierCallerRaw struct {
	Contract *NodeVerifierCaller // Generic read-only contract binding to access the raw methods on
}

// NodeVerifierTransactorRaw is an auto generated low-level write-only Go binding around a Kaia contract.
type NodeVerifierTransactorRaw struct {
	Contract *NodeVerifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNodeVerifier creates a new instance of NodeVerifier, bound to a specific deployed contract.
func NewNodeVerifier(address common.Address, backend bind.ContractBackend) (*NodeVerifier, error) {
	contract, err := bindNodeVerifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NodeVerifier{NodeVerifierCaller: NodeVerifierCaller{contract: contract}, NodeVerifierTransactor: NodeVerifierTransactor{contract: contract}, NodeVerifierFilterer: NodeVerifierFilterer{contract: contract}}, nil
}

// NewNodeVerifierCaller creates a new read-only instance of NodeVerifier, bound to a specific deployed contract.
func NewNodeVerifierCaller(address common.Address, caller bind.ContractCaller) (*NodeVerifierCaller, error) {
	contract, err := bindNodeVerifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NodeVerifierCaller{contract: contract}, nil
}

// NewNodeVerifierTransactor creates a new write-only instance of NodeVerifier, bound to a specific deployed contract.
func NewNodeVerifierTransactor(address common.Address, transactor bind.ContractTransactor) (*NodeVerifierTransactor, error) {
	contract, err := bindNodeVerifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NodeVerifierTransactor{contract: contract}, nil
}

// NewNodeVerifierFilterer creates a new log filterer instance of NodeVerifier, bound to a specific deployed contract.
func NewNodeVerifierFilterer(address common.Address, filterer bind.ContractFilterer) (*NodeVerifierFilterer, error) {
	contract, err := bindNodeVerifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NodeVerifierFilterer{contract: contract}, nil
}

// bindNodeVerifier binds a generic wrapper to an already deployed contract.
func bindNodeVerifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := NodeVerifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeVerifier *NodeVerifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeVerifier.Contract.NodeVerifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeVerifier *NodeVerifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeVerifier.Contract.NodeVerifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeVerifier *NodeVerifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeVerifier.Contract.NodeVerifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NodeVerifier *NodeVerifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NodeVerifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NodeVerifier *NodeVerifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NodeVerifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NodeVerifier *NodeVerifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NodeVerifier.Contract.contract.Transact(opts, method, params...)
}
