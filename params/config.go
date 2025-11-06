// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from params/config.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package params

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash      = common.HexToHash("0xc72e5293c3c3ba38ed8ae910f780e4caaa9fb95e79784f7ab74c3c262ea7137e") // mainnet genesis hash to enforce below configs on
	KairosGenesisHash       = common.HexToHash("0xe33ff05ceec2581ca9496f38a2bf9baad5d4eed629e896ccb33d1dc991bc4b4a") // Kairos genesis hash to enforce below configs on
	AuthorAddressForTesting = common.HexToAddress("0xc0ea08a2d404d3172d2add29a45be56da40e2949")
	mintingAmount, _        = new(big.Int).SetString("9600000000000000000", 10)
	logger                  = log.NewModuleLogger(log.Governance)
)

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:                  big.NewInt(int64(MainnetNetworkId)),
		IstanbulCompatibleBlock:  big.NewInt(86816005),
		LondonCompatibleBlock:    big.NewInt(86816005),
		EthTxTypeCompatibleBlock: big.NewInt(86816005),
		MagmaCompatibleBlock:     big.NewInt(99841497),
		KoreCompatibleBlock:      big.NewInt(119750400),
		ShanghaiCompatibleBlock:  big.NewInt(135456000),
		CancunCompatibleBlock:    big.NewInt(147534000),
		KaiaCompatibleBlock:      big.NewInt(162900480),
		RandaoCompatibleBlock:    big.NewInt(147534000),
		RandaoRegistry: &RegistryConfig{
			Records: map[string]common.Address{
				"KIP113": common.HexToAddress("0x3e80e75975bdb8e04B800485DD28BebeC6d97679"),
			},
			Owner: common.HexToAddress("0x04992a2B7E7CE809d409adE32185D49A96AAa32d"),
		},
		PragueCompatibleBlock:    big.NewInt(190670000),
		OsakaCompatibleBlock:     nil, // TODO-kaia-osaka: set Mainnet's OsakaCompatibleBlock
		BPO1CompatibleBlock:      nil,
		BPO2CompatibleBlock:      nil,
		BPO3CompatibleBlock:      nil,
		BPO4CompatibleBlock:      nil,
		BPO5CompatibleBlock:      nil,
		AmsterdamCompatibleBlock: nil,
		VerkleCompatibleBlock:    nil,
		Kip103CompatibleBlock:    big.NewInt(119750400),
		Kip103ContractAddress:    common.HexToAddress("0xD5ad6D61Dd87EdabE2332607C328f5cc96aeCB95"),
		Kip160CompatibleBlock:    big.NewInt(162900480),
		Kip160ContractAddress:    common.HexToAddress("0xa4df15717Da40077C0aD528296AdBBd046579Ee9"),
		DeriveShaImpl:            2,
		Governance: &GovernanceConfig{
			GoverningNode:  common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"),
			GovernanceMode: "single",
			Reward: &RewardConfig{
				MintingAmount:          mintingAmount,
				Ratio:                  "34/54/12",
				UseGiniCoeff:           true,
				DeferredTxFee:          true,
				StakingUpdateInterval:  86400,
				ProposerUpdateInterval: 3600,
				MinimumStake:           big.NewInt(5000000),
			},
		},
		Istanbul: &IstanbulConfig{
			Epoch:          604800,
			ProposerPolicy: 2,
			SubGroupSize:   22,
		},
		BlobScheduleConfig: &BlobScheduleConfig{
			Osaka: DefaultOsakaBlobConfig,
		},
		UnitPrice: 25000000000,
	}

	// KairosChainConfig contains the chain parameters to run a node on the Kairos.
	KairosChainConfig = &ChainConfig{
		ChainID:                  big.NewInt(int64(KairosNetworkId)),
		IstanbulCompatibleBlock:  big.NewInt(75373312),
		LondonCompatibleBlock:    big.NewInt(80295291),
		EthTxTypeCompatibleBlock: big.NewInt(86513895),
		MagmaCompatibleBlock:     big.NewInt(98347376),
		KoreCompatibleBlock:      big.NewInt(111736800),
		ShanghaiCompatibleBlock:  big.NewInt(131608000),
		CancunCompatibleBlock:    big.NewInt(141367000),
		KaiaCompatibleBlock:      big.NewInt(156660000),
		// Optional forks
		RandaoCompatibleBlock: big.NewInt(141367000),
		RandaoRegistry: &RegistryConfig{
			Records: map[string]common.Address{
				"KIP113": common.HexToAddress("0x4BEed0651C46aE5a7CB3b7737345d2ee733789e6"),
			},
			Owner: common.HexToAddress("0x04992a2B7E7CE809d409adE32185D49A96AAa32d"),
		},
		PragueCompatibleBlock:    big.NewInt(187930000),
		OsakaCompatibleBlock:     nil, // TODO-kaia-osaka: set Kairos' OsakaCompatibleBlock
		BPO1CompatibleBlock:      nil,
		BPO2CompatibleBlock:      nil,
		BPO3CompatibleBlock:      nil,
		BPO4CompatibleBlock:      nil,
		BPO5CompatibleBlock:      nil,
		AmsterdamCompatibleBlock: nil,
		VerkleCompatibleBlock:    nil,
		Kip103CompatibleBlock:    big.NewInt(119145600),
		Kip103ContractAddress:    common.HexToAddress("0xD5ad6D61Dd87EdabE2332607C328f5cc96aeCB95"),
		Kip160CompatibleBlock:    big.NewInt(156660000),
		Kip160ContractAddress:    common.HexToAddress("0x3D478E73c9dBebB72332712D7265961B1868d193"),
		// Genesis governance parameters
		DeriveShaImpl: 2,
		Governance: &GovernanceConfig{
			GoverningNode:  common.HexToAddress("0x99fb17d324fa0e07f23b49d09028ac0919414db6"),
			GovernanceMode: "single",
			Reward: &RewardConfig{
				MintingAmount:          mintingAmount,
				Ratio:                  "34/54/12",
				UseGiniCoeff:           true,
				DeferredTxFee:          true,
				StakingUpdateInterval:  86400,
				ProposerUpdateInterval: 3600,
				MinimumStake:           big.NewInt(5000000),
			},
		},
		Istanbul: &IstanbulConfig{
			Epoch:          604800,
			ProposerPolicy: 2,
			SubGroupSize:   22,
		},
		BlobScheduleConfig: &BlobScheduleConfig{
			Osaka: DefaultOsakaBlobConfig,
		},
		UnitPrice: 25000000000,
	}

	TestChainConfig = TestKaiaConfig("kaia")
)

func TestKaiaConfig(maxHardfork string) *ChainConfig {
	// Create a custom governance config with lower/upper bound base fee as 1 for testing
	testGovConfig := GetDefaultGovernanceConfig()
	testGovConfig.KIP71.LowerBoundBaseFee = 0
	testGovConfig.KIP71.UpperBoundBaseFee = 1

	chainConfig := &ChainConfig{
		ChainID:                 big.NewInt(1),
		IstanbulCompatibleBlock: big.NewInt(0),
		Istanbul: &IstanbulConfig{
			Epoch:          DefaultEpoch,
			ProposerPolicy: WeightedRandom,
			SubGroupSize:   DefaultSubGroupSize,
		},
		UnitPrice:     1, // NOTE-Kaia Use 1 for testing
		DeriveShaImpl: 0,
		Governance:    testGovConfig,
	}
	if maxHardfork == "istanbul" {
		return chainConfig
	}
	chainConfig.LondonCompatibleBlock = big.NewInt(0)
	if maxHardfork == "london" {
		return chainConfig
	}
	chainConfig.EthTxTypeCompatibleBlock = big.NewInt(0)
	if maxHardfork == "ethTxType" {
		return chainConfig
	}
	chainConfig.MagmaCompatibleBlock = big.NewInt(0)
	if maxHardfork == "magma" {
		return chainConfig
	}
	chainConfig.KoreCompatibleBlock = big.NewInt(0)
	if maxHardfork == "kore" {
		return chainConfig
	}
	chainConfig.ShanghaiCompatibleBlock = big.NewInt(0)
	if maxHardfork == "shanghai" {
		return chainConfig
	}
	chainConfig.CancunCompatibleBlock = big.NewInt(0)
	if maxHardfork == "cancun" {
		return chainConfig
	}
	chainConfig.KaiaCompatibleBlock = big.NewInt(0)
	if maxHardfork == "kaia" {
		return chainConfig
	}
	chainConfig.PragueCompatibleBlock = big.NewInt(0)
	if maxHardfork == "prague" {
		return chainConfig
	}
	chainConfig.RandaoCompatibleBlock = big.NewInt(0)
	if maxHardfork == "randao" {
		return chainConfig
	}

	return chainConfig
}

// VMLogTarget sets the output target of vmlog.
// The values below can be OR'ed.
//   - 0x0: no output (default)
//   - 0x1: file (DATADIR/logs/vm.log)
//   - 0x2: stdout (like logger.DEBUG)
var VMLogTarget = 0x0

const (
	VMLogToFile   = 0x1
	VMLogToStdout = 0x2
	VMLogToAll    = VMLogToFile | VMLogToStdout

	UpperGasLimit = uint64(999999999999)

	// Default max price for gas price oracle
	DefaultGPOMaxPrice = 500 * Gkei
)

const (
	PasswordLength = 16
)

var (
	// DefaultOsakaBlobConfig is the default blob configuration for the Osaka fork.
	DefaultOsakaBlobConfig = &BlobConfig{
		Target:         6,
		Max:            9,
		UpdateFraction: 5007716,
	}
	// DefaultBPO1BlobConfig is the default blob configuration for the BPO1 fork.
	DefaultBPO1BlobConfig = &BlobConfig{
		Target:         10,
		Max:            15,
		UpdateFraction: 8346193,
	}
	// DefaultBPO2BlobConfig is the default blob configuration for the BPO2 fork.
	DefaultBPO2BlobConfig = &BlobConfig{
		Target:         14,
		Max:            21,
		UpdateFraction: 11684671,
	}
	// DefaultBPO3BlobConfig is the default blob configuration for the BPO3 fork.
	DefaultBPO3BlobConfig = &BlobConfig{
		Target:         21,
		Max:            32,
		UpdateFraction: 20609697,
	}
	// DefaultBPO4BlobConfig is the default blob configuration for the BPO4 fork.
	DefaultBPO4BlobConfig = &BlobConfig{
		Target:         14,
		Max:            21,
		UpdateFraction: 13739630,
	}
	// DefaultBlobSchedule is the latest configured blob schedule for Ethereum mainnet.
	DefaultBlobSchedule = &BlobScheduleConfig{
		Osaka: DefaultOsakaBlobConfig,
	}
)

// ChainConfig is the blockchain config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	// "Compatible" means that it is EVM compatible(the opcode and precompiled contracts are the same as Ethereum EVM).
	// In other words, not all the hard fork items are included.
	IstanbulCompatibleBlock  *big.Int `json:"istanbulCompatibleBlock,omitempty"`  // IstanbulCompatibleBlock switch block (nil = no fork, 0 = already on istanbul)
	LondonCompatibleBlock    *big.Int `json:"londonCompatibleBlock,omitempty"`    // LondonCompatibleBlock switch block (nil = no fork, 0 = already on london)
	EthTxTypeCompatibleBlock *big.Int `json:"ethTxTypeCompatibleBlock,omitempty"` // EthTxTypeCompatibleBlock switch block (nil = no fork, 0 = already on ethTxType)
	MagmaCompatibleBlock     *big.Int `json:"magmaCompatibleBlock,omitempty"`     // MagmaCompatible switch block (nil = no fork, 0 already on Magma)
	KoreCompatibleBlock      *big.Int `json:"koreCompatibleBlock,omitempty"`      // KoreCompatible switch block (nil = no fork, 0 already on Kore)
	ShanghaiCompatibleBlock  *big.Int `json:"shanghaiCompatibleBlock,omitempty"`  // ShanghaiCompatible switch block (nil = no fork, 0 already on shanghai)
	CancunCompatibleBlock    *big.Int `json:"cancunCompatibleBlock,omitempty"`    // CancunCompatible switch block (nil = no fork, 0 already on Cancun)
	KaiaCompatibleBlock      *big.Int `json:"kaiaCompatibleBlock,omitempty"`      // KaiaCompatible switch block (nil = no fork, 0 already on Kaia)
	PragueCompatibleBlock    *big.Int `json:"pragueCompatibleBlock,omitempty"`    // PragueCompatible switch block (nil = no fork)
	OsakaCompatibleBlock     *big.Int `json:"osakaCompatibleBlock,omitempty"`     // OsakaCompatible switch block (nil = no fork)
	BPO1CompatibleBlock      *big.Int `json:"bpo1CompatibleBlock,omitempty"`      // BPO1Compatible switch block (nil = no fork, 0 = already on bpo1)
	BPO2CompatibleBlock      *big.Int `json:"bpo2CompatibleBlock,omitempty"`      // BPO2Compatible switch block (nil = no fork, 0 = already on bpo2)
	BPO3CompatibleBlock      *big.Int `json:"bpo3CompatibleBlock,omitempty"`      // BPO3Compatible switch block (nil = no fork, 0 = already on bpo3)
	BPO4CompatibleBlock      *big.Int `json:"bpo4CompatibleBlock,omitempty"`      // BPO4Compatible switch block (nil = no fork, 0 = already on bpo4)
	BPO5CompatibleBlock      *big.Int `json:"bpo5CompatibleBlock,omitempty"`      // BPO5Compatible switch block (nil = no fork, 0 = already on bpo5)
	AmsterdamCompatibleBlock *big.Int `json:"amsterdamCompatibleBlock,omitempty"` // AmsterdamCompatible switch block (nil = no fork, 0 = already on amsterdam)
	VerkleCompatibleBlock    *big.Int `json:"verkleCompatibleBlock,omitempty"`    // VerkleCompatible switch block (nil = no fork, 0 = already on verkle)

	// Kip103 is a special purpose hardfork feature that can be executed only once
	// Both Kip103CompatibleBlock and Kip103ContractAddress should be specified to enable KIP103
	Kip103CompatibleBlock *big.Int       `json:"kip103CompatibleBlock,omitempty"` // Kip103Compatible activate block (nil = no fork)
	Kip103ContractAddress common.Address `json:"kip103ContractAddress,omitempty"` // Kip103 contract address already deployed on the network

	// Kip160 is an optional hardfork
	// Both Kip160CompatibleBlock and Kip160ContractAddress should be specified to enable KIP160
	Kip160CompatibleBlock *big.Int       `json:"kip160CompatibleBlock,omitempty"` // Kip160Compatible activate block (nil = no fork)
	Kip160ContractAddress common.Address `json:"kip160ContractAddress,omitempty"` // Kip160 contract address already deployed on the network

	// Randao is an optional hardfork
	// RandaoCompatibleBlock, RandaoRegistryRecords and RandaoRegistryOwner all must be specified to enable Randao
	RandaoCompatibleBlock *big.Int        `json:"randaoCompatibleBlock,omitempty"` // RandaoCompatible activate block (nil = no fork)
	RandaoRegistry        *RegistryConfig `json:"randaoRegistry,omitempty"`        // Registry initial states

	// Various consensus engines
	Istanbul           *IstanbulConfig     `json:"istanbul,omitempty"`
	BlobScheduleConfig *BlobScheduleConfig `json:"blobSchedule,omitempty"`

	UnitPrice     uint64            `json:"unitPrice"`
	DeriveShaImpl int               `json:"deriveShaImpl"`
	Governance    *GovernanceConfig `json:"governance"`
}

// UnmarshalJSON implements json.Unmarshaler to detect deprecated consensus engine configurations
func (c *ChainConfig) UnmarshalJSON(data []byte) error {
	// First, unmarshal into a temporary struct that includes deprecated fields
	type Alias ChainConfig
	aux := &struct {
		*Alias
		Clique json.RawMessage `json:"clique,omitempty"`
		Gxhash json.RawMessage `json:"gxhash,omitempty"`
	}{
		Alias: (*Alias)(c),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Check for deprecated fields and return error
	if aux.Clique != nil {
		return errors.New("DEPRECATED: 'clique' consensus engine configuration is deprecated and has been removed. Please use 'istanbul' instead")
	}
	if aux.Gxhash != nil {
		return errors.New("DEPRECATED: 'gxhash' consensus engine configuration is deprecated and has been removed. Please use 'istanbul' instead")
	}

	return nil
}

// GovernanceConfig stores governance information for a network
type GovernanceConfig struct {
	GoverningNode    common.Address `json:"governingNode"`
	GovernanceMode   string         `json:"governanceMode"`
	GovParamContract common.Address `json:"govParamContract"`
	Reward           *RewardConfig  `json:"reward,omitempty"`
	KIP71            *KIP71Config   `json:"kip71,omitempty"`
}

func (g *GovernanceConfig) DeferredTxFee() bool {
	return g.Reward.DeferredTxFee
}

// RewardConfig stores information about the network's token economy
type RewardConfig struct {
	MintingAmount          *big.Int `json:"mintingAmount"`
	Ratio                  string   `json:"ratio"`                  // Define how much portion of reward be distributed to CN/KIF/KEF
	Kip82Ratio             string   `json:"kip82ratio,omitempty"`   // Define how much portion of reward be distributed to proposer/stakers
	UseGiniCoeff           bool     `json:"useGiniCoeff"`           // Decide if Gini Coefficient will be used or not
	DeferredTxFee          bool     `json:"deferredTxFee"`          // Decide if TX fee will be handled instantly or handled later at block finalization
	StakingUpdateInterval  uint64   `json:"stakingUpdateInterval"`  // Interval when staking information is updated
	ProposerUpdateInterval uint64   `json:"proposerUpdateInterval"` // Interval when proposer information is updated
	MinimumStake           *big.Int `json:"minimumStake"`           // Minimum amount of kei to join CCO
}

// Magma governance parameters
type KIP71Config struct {
	LowerBoundBaseFee         uint64 `json:"lowerboundbasefee"`         // Minimum base fee for dynamic gas price
	UpperBoundBaseFee         uint64 `json:"upperboundbasefee"`         // Maximum base fee for dynamic gas price
	GasTarget                 uint64 `json:"gastarget"`                 // Gauge parameter increasing or decreasing gas price
	MaxBlockGasUsedForBaseFee uint64 `json:"maxblockgasusedforbasefee"` // Maximum network and process capacity to allow in a block
	BaseFeeDenominator        uint64 `json:"basefeedenominator"`        // For normalizing effect of the rapid change like impulse gas used
}

// IstanbulConfig is the consensus engine configs for Istanbul based sealing.
type IstanbulConfig struct {
	Epoch          uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
	ProposerPolicy uint64 `json:"policy"` // The policy for proposer selection; 0: Round Robin, 1: Sticky, 2: Weighted Random
	SubGroupSize   uint64 `json:"sub"`
}

// RegistryConfig is the initial KIP-149 system contract registry states.
// It is installed at block (RandaoCompatibleBlock - 1). Initial states are not applied if RandaoCompatibleBlock is nil or 0.
// To install the initial states from the block 0, use the AllocRegistry to generate GenesisAlloc.
//
// This struct only represents a special case of Registry state where:
// - there is only one record for each name
// - the activation of all records is zero
// - the names array is lexicographically sorted
type RegistryConfig struct {
	Records map[string]common.Address `json:"records"`
	Owner   common.Address            `json:"owner"`
}

// String implements the stringer interface, returning the consensus engine details.
func (c *IstanbulConfig) String() string {
	return "istanbul"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Istanbul != nil:
		engine = c.Istanbul
	default:
		engine = "unknown"
	}

	kip103 := fmt.Sprintf(" KIP103CompatibleBlock: %v KIP103ContractAddress %s", c.Kip103CompatibleBlock, c.Kip103ContractAddress.String())
	kip160 := fmt.Sprintf(" KIP160CompatibleBlock: %v KIP160ContractAddress %s", c.Kip160CompatibleBlock, c.Kip160ContractAddress.String())
	var subGroupSize string = ""
	if c.Istanbul != nil {
		subGroupSize = fmt.Sprintf(" SubGroupSize: %d", c.Istanbul.SubGroupSize)
	}

	return fmt.Sprintf("{ChainID: %v"+
		" IstanbulCompatibleBlock: %v"+
		" LondonCompatibleBlock: %v"+
		" EthTxTypeCompatibleBlock: %v"+
		" MagmaCompatibleBlock: %v"+
		" KoreCompatibleBlock: %v"+
		" ShanghaiCompatibleBlock: %v"+
		" CancunCompatibleBlock: %v"+
		" KaiaCompatibleBlock: %v"+
		" RandaoCompatibleBlock: %v"+
		" PragueCompatibleBlock: %v"+
		" OsakaCompatibleBlock: %v"+
		"%s%s%s"+
		" UnitPrice: %d"+
		" DeriveShaImpl: %d"+
		" Engine: %v}",
		c.ChainID,
		c.IstanbulCompatibleBlock,
		c.LondonCompatibleBlock,
		c.EthTxTypeCompatibleBlock,
		c.MagmaCompatibleBlock,
		c.KoreCompatibleBlock,
		c.ShanghaiCompatibleBlock,
		c.CancunCompatibleBlock,
		c.KaiaCompatibleBlock,
		c.RandaoCompatibleBlock,
		c.PragueCompatibleBlock,
		c.OsakaCompatibleBlock,
		kip103, kip160, subGroupSize,
		c.UnitPrice,
		c.DeriveShaImpl,
		engine,
	)
}

func (c *ChainConfig) Copy() *ChainConfig {
	r := &ChainConfig{}
	j, _ := json.Marshal(c)
	json.Unmarshal(j, r)
	return r
}

// BlobConfig specifies the target and max blobs per block for the associated fork.
type BlobConfig struct {
	Target         int    `json:"target"`
	Max            int    `json:"max"`
	UpdateFraction uint64 `json:"baseFeeUpdateFraction"`
}

// String implement fmt.Stringer, returning string format blob config.
func (bc *BlobConfig) String() string {
	if bc == nil {
		return "nil"
	}
	return fmt.Sprintf("target: %d, max: %d, fraction: %d", bc.Target, bc.Max, bc.UpdateFraction)
}

// BlobScheduleConfig determines target and max number of blobs allow per fork.
type BlobScheduleConfig struct {
	Osaka     *BlobConfig `json:"osaka,omitempty"`
	Verkle    *BlobConfig `json:"verkle,omitempty"`
	BPO1      *BlobConfig `json:"bpo1,omitempty"`
	BPO2      *BlobConfig `json:"bpo2,omitempty"`
	BPO3      *BlobConfig `json:"bpo3,omitempty"`
	BPO4      *BlobConfig `json:"bpo4,omitempty"`
	BPO5      *BlobConfig `json:"bpo5,omitempty"`
	Amsterdam *BlobConfig `json:"amsterdam,omitempty"`
}

func (c *ChainConfig) BlobConfig(head *big.Int) *BlobConfig {
	if c.IsBPO5ForkEnabled(head) {
		return c.BlobScheduleConfig.BPO5
	}
	if c.IsBPO4ForkEnabled(head) {
		return c.BlobScheduleConfig.BPO4
	}
	if c.IsBPO3ForkEnabled(head) {
		return c.BlobScheduleConfig.BPO3
	}
	if c.IsBPO2ForkEnabled(head) {
		return c.BlobScheduleConfig.BPO2
	}
	if c.IsBPO1ForkEnabled(head) {
		return c.BlobScheduleConfig.BPO1
	}
	if c.IsOsakaForkEnabled(head) {
		return c.BlobScheduleConfig.Osaka
	}
	return nil
}

// IsIstanbulForkEnabled returns whether num is either equal to the istanbul block or greater.
func (c *ChainConfig) IsIstanbulForkEnabled(num *big.Int) bool {
	return isForked(c.IstanbulCompatibleBlock, num)
}

// IsLondonForkEnabled returns whether num is either equal to the london block or greater.
func (c *ChainConfig) IsLondonForkEnabled(num *big.Int) bool {
	return isForked(c.LondonCompatibleBlock, num)
}

// IsEthTxTypeForkEnabled returns whether num is either equal to the ethTxType block or greater.
func (c *ChainConfig) IsEthTxTypeForkEnabled(num *big.Int) bool {
	return isForked(c.EthTxTypeCompatibleBlock, num)
}

// IsMagmaForkEnabled returns whether num is either equal to the magma block or greater.
func (c *ChainConfig) IsMagmaForkEnabled(num *big.Int) bool {
	return isForked(c.MagmaCompatibleBlock, num)
}

// IsKoreForkEnabled returns whether num is either equal to the kore block or greater.
func (c *ChainConfig) IsKoreForkEnabled(num *big.Int) bool {
	return isForked(c.KoreCompatibleBlock, num)
}

// IsShanghaiForkEnabled returns whether num is either equal to the shanghai block or greater.
func (c *ChainConfig) IsShanghaiForkEnabled(num *big.Int) bool {
	return isForked(c.ShanghaiCompatibleBlock, num)
}

// IsCancunForkEnabled returns whether num is either equal to the cancun block or greater.
func (c *ChainConfig) IsCancunForkEnabled(num *big.Int) bool {
	return isForked(c.CancunCompatibleBlock, num)
}

// IsKaiaForkEnabled returns whether num is either equal to the kaia block or greater.
func (c *ChainConfig) IsKaiaForkEnabled(num *big.Int) bool {
	return isForked(c.KaiaCompatibleBlock, num)
}

// IsKip103ForkEnabled returns whether num is either equal to the kip103 block or greater.
func (c *ChainConfig) IsKip103ForkEnabled(num *big.Int) bool {
	return isForked(c.Kip103CompatibleBlock, num)
}

// IsKip160ForkEnabled returns whether num is either equal to the kip160 block or greater.
func (c *ChainConfig) IsKip160ForkEnabled(num *big.Int) bool {
	return isForked(c.Kip160CompatibleBlock, num)
}

// IsRandaoForkEnabled returns whether num is either equal to the randao block or greater.
func (c *ChainConfig) IsRandaoForkEnabled(num *big.Int) bool {
	return isForked(c.RandaoCompatibleBlock, num)
}

// IsPragueForkEnabled returns whether num is either equal to the prague block or greater.
func (c *ChainConfig) IsPragueForkEnabled(num *big.Int) bool {
	return isForked(c.PragueCompatibleBlock, num)
}

// IsOsakaForkEnabled returns whether num is either equal to the osaka block or greater.
func (c *ChainConfig) IsOsakaForkEnabled(num *big.Int) bool {
	return isForked(c.OsakaCompatibleBlock, num)
}

// IsBPO1ForkEnabled returns whether num is either equal to the bpo1 block or greater.
func (c *ChainConfig) IsBPO1ForkEnabled(num *big.Int) bool {
	return isForked(c.BPO1CompatibleBlock, num)
}

// IsBPO2ForkEnabled returns whether num is either equal to the bpo2 block or greater.
func (c *ChainConfig) IsBPO2ForkEnabled(num *big.Int) bool {
	return isForked(c.BPO2CompatibleBlock, num)
}

// IsBPO3ForkEnabled returns whether num is either equal to the bpo3 block or greater.
func (c *ChainConfig) IsBPO3ForkEnabled(num *big.Int) bool {
	return isForked(c.BPO3CompatibleBlock, num)
}

// IsBPO4ForkEnabled returns whether num is either equal to the bpo4 block or greater.
func (c *ChainConfig) IsBPO4ForkEnabled(num *big.Int) bool {
	return isForked(c.BPO4CompatibleBlock, num)
}

// IsBPO5ForkEnabled returns whether num is either equal to the bpo5 block or greater.
func (c *ChainConfig) IsBPO5ForkEnabled(num *big.Int) bool {
	return isForked(c.BPO5CompatibleBlock, num)
}

// IsKIP103ForkBlock returns whether num is equal to the kip103 block.
func (c *ChainConfig) IsKIP103ForkBlock(num *big.Int) bool {
	return isForkBlock(c.Kip103CompatibleBlock, num)
}

// IsKIP160ForkBlock returns whether num is equal to the kip160 block.
func (c *ChainConfig) IsKIP160ForkBlock(num *big.Int) bool {
	return isForkBlock(c.Kip160CompatibleBlock, num)
}

// IsRandaoForkBlockParent returns whether num is one block before the randao block.
func (c *ChainConfig) IsRandaoForkBlockParent(num *big.Int) bool {
	return isForkBlockParent(c.RandaoCompatibleBlock, num)
}

// IsRandaoForkBlock returns whether num is equal to the randao block.
func (c *ChainConfig) IsRandaoForkBlock(num *big.Int) bool {
	return isForkBlock(c.RandaoCompatibleBlock, num)
}

// IsKaiaForkBlockParent returns whether num is equal to the kaia block.
func (c *ChainConfig) IsKaiaForkBlockParent(num *big.Int) bool {
	return isForkBlockParent(c.KaiaCompatibleBlock, num)
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name     string
		block    *big.Int
		optional bool // if true, the fork may be nil and next fork is still allowed
	}
	var lastFork fork
	for _, cur := range []fork{
		{name: "istanbulBlock", block: c.IstanbulCompatibleBlock},
		{name: "londonBlock", block: c.LondonCompatibleBlock},
		{name: "ethTxTypeBlock", block: c.EthTxTypeCompatibleBlock},
		{name: "magmaBlock", block: c.MagmaCompatibleBlock},
		{name: "koreBlock", block: c.KoreCompatibleBlock},
		{name: "shanghaiBlock", block: c.ShanghaiCompatibleBlock},
		{name: "cancunBlock", block: c.CancunCompatibleBlock},
		{name: "randaoBlock", block: c.RandaoCompatibleBlock, optional: true},
		{name: "kaiaBlock", block: c.KaiaCompatibleBlock},
		{name: "pragueBlock", block: c.PragueCompatibleBlock},
		{name: "osakaBlock", block: c.OsakaCompatibleBlock},
	} {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || cur.block != nil {
			lastFork = cur
		}
	}
	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.IstanbulCompatibleBlock, newcfg.IstanbulCompatibleBlock, head) {
		return newCompatError("Istanbul Block", c.IstanbulCompatibleBlock, newcfg.IstanbulCompatibleBlock)
	}
	if isForkIncompatible(c.LondonCompatibleBlock, newcfg.LondonCompatibleBlock, head) {
		return newCompatError("London Block", c.LondonCompatibleBlock, newcfg.LondonCompatibleBlock)
	}
	if isForkIncompatible(c.EthTxTypeCompatibleBlock, newcfg.EthTxTypeCompatibleBlock, head) {
		return newCompatError("EthTxType Block", c.EthTxTypeCompatibleBlock, newcfg.EthTxTypeCompatibleBlock)
	}
	if isForkIncompatible(c.MagmaCompatibleBlock, newcfg.MagmaCompatibleBlock, head) {
		return newCompatError("Magma Block", c.MagmaCompatibleBlock, newcfg.MagmaCompatibleBlock)
	}
	if isForkIncompatible(c.KoreCompatibleBlock, newcfg.KoreCompatibleBlock, head) {
		return newCompatError("Kore Block", c.KoreCompatibleBlock, newcfg.KoreCompatibleBlock)
	}
	// We have intentionally skipped kip103Block in the fork ordering check since kip103 is designed
	// as an optional hardfork and there are no dependency with other forks.
	if isForkIncompatible(c.ShanghaiCompatibleBlock, newcfg.ShanghaiCompatibleBlock, head) {
		return newCompatError("Shanghai Block", c.ShanghaiCompatibleBlock, newcfg.ShanghaiCompatibleBlock)
	}
	if isForkIncompatible(c.CancunCompatibleBlock, newcfg.CancunCompatibleBlock, head) {
		return newCompatError("Cancun Block", c.CancunCompatibleBlock, newcfg.CancunCompatibleBlock)
	}
	if isForkIncompatible(c.KaiaCompatibleBlock, newcfg.KaiaCompatibleBlock, head) {
		return newCompatError("Kaia Block", c.KaiaCompatibleBlock, newcfg.KaiaCompatibleBlock)
	}
	if isForkIncompatible(c.RandaoCompatibleBlock, newcfg.RandaoCompatibleBlock, head) {
		return newCompatError("Randao Block", c.RandaoCompatibleBlock, newcfg.RandaoCompatibleBlock)
	}
	if isForkIncompatible(c.PragueCompatibleBlock, newcfg.PragueCompatibleBlock, head) {
		return newCompatError("Prague Block", c.PragueCompatibleBlock, newcfg.PragueCompatibleBlock)
	}
	if isForkIncompatible(c.OsakaCompatibleBlock, newcfg.OsakaCompatibleBlock, head) {
		return newCompatError("Osaka Block", c.OsakaCompatibleBlock, newcfg.OsakaCompatibleBlock)
	}
	return nil
}

// SetDefaultsForGenesis fills undefined chain config with default values.
// Only used for generating genesis.
// Empty values from genesis.json will be left out from genesis.
func (c *ChainConfig) SetDefaultsForGenesis() {
	if c.Istanbul == nil {
		c.Istanbul = GetDefaultIstanbulConfig()
		logger.Warn("Override the default Istanbul config to the chain config")
	}

	if c.Governance == nil {
		c.Governance = GetDefaultGovernanceConfigForGenesis()
		logger.Warn("Override the default governance config to the chain config")
	}

	if c.Governance.Reward == nil {
		c.Governance.Reward = GetDefaultRewardConfigForGenesis()
		logger.Warn("Override the default governance reward config to the chain config", "reward",
			c.Governance.Reward)
	}

	// StakingUpdateInterval must be nonzero because it is used as denominator
	if c.Governance.Reward.StakingUpdateInterval == 0 {
		c.Governance.Reward.StakingUpdateInterval = StakingUpdateInterval()
		logger.Warn("Override the default staking update interval to the chain config", "interval",
			c.Governance.Reward.StakingUpdateInterval)
	}

	// ProposerUpdateInterval must be nonzero because it is used as denominator
	if c.Governance.Reward.ProposerUpdateInterval == 0 {
		c.Governance.Reward.ProposerUpdateInterval = ProposerUpdateInterval()
		logger.Warn("Override the default proposer update interval to the chain config", "interval",
			c.Governance.Reward.ProposerUpdateInterval)
	}
}

// SetDefaults fills undefined chain config with default values
// so that nil pointer does not exist in the chain config
func (c *ChainConfig) SetDefaults() {
	c.SetDefaultsForGenesis()

	if c.Governance.KIP71 == nil {
		c.Governance.KIP71 = GetDefaultKIP71Config()
	}
	if c.Governance.Reward.Kip82Ratio == "" {
		c.Governance.Reward.Kip82Ratio = DefaultKip82Ratio
	}
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

// isForkBlock returns whether given head block is exactly the fork block s.
func isForkBlock(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) == 0
}

func isForkBlockParent(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	nextNum := new(big.Int).Add(head, common.Big1)
	return s.Cmp(nextNum) == 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID     *big.Int
	IsIstanbul  bool
	IsLondon    bool
	IsEthTxType bool
	IsMagma     bool
	IsKore      bool
	IsShanghai  bool
	IsCancun    bool
	IsKaia      bool
	IsRandao    bool
	IsPrague    bool
	IsOsaka     bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:     new(big.Int).Set(chainID),
		IsIstanbul:  c.IsIstanbulForkEnabled(num),
		IsLondon:    c.IsLondonForkEnabled(num),
		IsEthTxType: c.IsEthTxTypeForkEnabled(num),
		IsMagma:     c.IsMagmaForkEnabled(num),
		IsKore:      c.IsKoreForkEnabled(num),
		IsShanghai:  c.IsShanghaiForkEnabled(num),
		IsCancun:    c.IsCancunForkEnabled(num),
		IsKaia:      c.IsKaiaForkEnabled(num),
		IsRandao:    c.IsRandaoForkEnabled(num),
		IsPrague:    c.IsPragueForkEnabled(num),
		IsOsaka:     c.IsOsakaForkEnabled(num),
	}
}

// Mainnet genesis config
func GetDefaultGovernanceConfigForGenesis() *GovernanceConfig {
	gov := &GovernanceConfig{
		GovernanceMode: DefaultGovernanceMode,
		GoverningNode:  common.HexToAddress(DefaultGoverningNode),
		Reward:         GetDefaultRewardConfigForGenesis(),
	}
	return gov
}

func GetDefaultGovernanceConfig() *GovernanceConfig {
	gov := &GovernanceConfig{
		GovernanceMode:   DefaultGovernanceMode,
		GoverningNode:    common.HexToAddress(DefaultGoverningNode),
		GovParamContract: common.HexToAddress(DefaultGovParamContract),
		Reward:           GetDefaultRewardConfig(),
		KIP71:            GetDefaultKIP71Config(),
	}
	return gov
}

func GetDefaultIstanbulConfig() *IstanbulConfig {
	return &IstanbulConfig{
		Epoch:          DefaultEpoch,
		ProposerPolicy: DefaultProposerPolicy,
		SubGroupSize:   DefaultSubGroupSize,
	}
}

func GetDefaultRewardConfigForGenesis() *RewardConfig {
	return &RewardConfig{
		MintingAmount:          DefaultMintingAmount,
		Ratio:                  DefaultRatio,
		UseGiniCoeff:           DefaultUseGiniCoeff,
		DeferredTxFee:          DefaultDeferredTxFee,
		StakingUpdateInterval:  DefaultStakeUpdateInterval,
		ProposerUpdateInterval: DefaultProposerRefreshInterval,
		MinimumStake:           DefaultMinimumStake,
	}
}

func GetDefaultRewardConfig() *RewardConfig {
	return &RewardConfig{
		MintingAmount:          DefaultMintingAmount,
		Ratio:                  DefaultRatio,
		Kip82Ratio:             DefaultKip82Ratio,
		UseGiniCoeff:           DefaultUseGiniCoeff,
		DeferredTxFee:          DefaultDeferredTxFee,
		StakingUpdateInterval:  DefaultStakeUpdateInterval,
		ProposerUpdateInterval: DefaultProposerRefreshInterval,
		MinimumStake:           DefaultMinimumStake,
	}
}

func GetDefaultKIP71Config() *KIP71Config {
	return &KIP71Config{
		LowerBoundBaseFee:         DefaultLowerBoundBaseFee,
		UpperBoundBaseFee:         DefaultUpperBoundBaseFee,
		GasTarget:                 DefaultGasTarget,
		MaxBlockGasUsedForBaseFee: DefaultMaxBlockGasUsedForBaseFee,
		BaseFeeDenominator:        DefaultBaseFeeDenominator,
	}
}
