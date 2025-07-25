// Modifications Copyright 2024 The Kaia Authors
// Copyright 2018 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package log

import (
	"strconv"
	"strings"
	"time"
)

// StatsReportLimit is the time limit during working after which we always print
// out progress. This avoids the user wondering what's going on.
const StatsReportLimit = 10 * time.Second

type ModuleID int

// When printID is true, log prints ModuleID instead of ModuleName.
// TODO-Kaia This can be controlled by runtime configuration.
var printID = true

func GetModuleName(mi ModuleID) string {
	return moduleNames[int(mi)]
}

func GetModuleID(moduleName string) ModuleID {
	moduleName = strings.ToLower(moduleName)
	for i := 0; i < len(moduleNames); i++ {
		if moduleName == moduleNames[i] {
			return ModuleID(i)
		}
	}
	return ModuleNameLen
}

func (mi ModuleID) String() string {
	if printID {
		return strconv.Itoa(int(mi))
	}
	return moduleNames[mi]
}

// NOTE-Kaia-Log Please add module in lexicographical order.
const (
	// 0
	BaseLogger ModuleID = iota

	// 1~10
	AccountsAbiBind
	AccountsKeystore
	API
	APIDebug
	Blockchain
	BlockchainState
	BlockchainTypes
	BlockchainTypesAccount
	BlockchainTypesAccountKey
	CMDIstanbul

	// 11~20
	CMDKBN
	CMDKCN
	CMDKEN
	CMDKGEN
	CMDKaia
	CMDKPN
	CMDKSCN
	CMDUtils
	CMDUtilsNodeCMD
	Common

	// 21~30
	ConsensusClique
	ConsensusGxhash
	ConsensusIstanbul
	ConsensusIstanbulBackend
	ConsensusIstanbulCore
	ConsensusIstanbulValidator
	Console
	DatasyncDownloader
	DatasyncFetcher
	Governance

	// 31~40
	Metrics
	NetworksGRPC
	NetworksP2P
	NetworksP2PDiscover
	NetworksP2PNat
	NetworksP2PSimulations
	NetworksP2PSimulationsAdapters
	NetworksP2PSimulationsCnism
	NetworksRPC
	Node

	// 41~50
	NodeCN
	NodeCNFilters
	NodeCNTracers
	Reward
	ServiceChain
	Snapshot
	SnapshotSync
	StorageDatabase
	StorageStateDB
	VM

	// 51~60
	Work
	CMDKSPN
	CMDKSEN
	ChainDataFetcher
	KAS
	FORK
	NodeCnGasPrice
	KaiaxStaking
	KaiaxReward
	KaiaxSupply

	// 61~70
	KaiaxGov
	KaiaxValset
	KaiaxRandao
	KaiaxGasless
	Builder
	KaiaxAuction

	// ModuleNameLen should be placed at the end of the list.
	ModuleNameLen
)

var moduleNames = [ModuleNameLen]string{
	// 0
	"defaultLogger",

	// 1~10
	"accounts/abi/bind",
	"accounts/keystore",
	"api",
	"api/debug",
	"blockchain",
	"blockchain/state",
	"blockchain/types",
	"blockchain/types/account",
	"blockchain/types/accountkey",
	"cmd/istanbul",

	// 11~20
	"cmd/kbn",
	"cmd/kcn",
	"cmd/ken",
	"cmd/kgen",
	"cmd/klay",
	"cmd/kpn",
	"cmd/kscn",
	"cmd/utils",
	"cmd/utils/nodecmd",
	"common",

	// 21~30
	"consensus/clique",
	"consensus/gxhash",
	"consensus/istanbul",
	"consensus/istanbul/backend",
	"consensus/istanbul/core",
	"consensus/istanbul/validator",
	"console",
	"datasync/downloader",
	"datasync/fetcher",
	"governance/governance",

	// 31~40
	"metrics",
	"networks/grpc",
	"networks/p2p",
	"networks/p2p/discover",
	"networks/p2p/nat",
	"networks/p2p/simulations",
	"networks/p2p/simulations/adapters",
	"networks/p2p/simulations/cnism",
	"networks/rpc",
	"node",

	// 41~50
	"node/cn",
	"node/cn/filters",
	"node/cn/tracers",
	"contracts/reward",
	"servicechain",
	"snapshot",
	"node/cn/snap",
	"storage/database",
	"storage/statedb",
	"vm",

	// 51~60
	"work",
	"cmd/kspn",
	"cmd/ksen",
	"datasync/chaindatafetcher",
	"kas",
	"fork",
	"node/cn/gasprice",
	"kaiax/staking",
	"kaiax/reward",
	"kaiax/supply",

	// 61~70
	"kaiax/gov",
	"kaiax/valset",
	"kaiax/randao",
	"kaiax/gasless",
	"builder",
	"kaiax/auction",
}
