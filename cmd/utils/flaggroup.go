// Modifications Copyright 2024 The Kaia Authors
// Copyright 2020 The klaytn Authors
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

package utils

import (
	"sort"

	"github.com/kaiachain/kaia/api/debug"
	"github.com/urfave/cli/v2"
)

const uncategorized = "MISC" // Uncategorized flags will belong to this group

// FlagGroup is a collection of flags belonging to a single topic.
type FlagGroup struct {
	Name  string
	Flags []cli.Flag
}

// TODO-Kaia: consider changing the type of FlagGroups to map
// FlagGroups categorizes flags into groups to print structured help.
var FlagGroups = []FlagGroup{
	{
		Name: "KLAY",
		Flags: []cli.Flag{
			NtpDisableFlag,
			NtpServerFlag,
			DbTypeFlag,
			DataDirFlag,
			ChainDataDirFlag,
			IdentityFlag,
			SyncModeFlag,
			GCModeFlag,
			SrvTypeFlag,
			ExtraDataFlag,
			ConfigFileFlag,
			OverwriteGenesisFlag,
			StartBlockNumberFlag,
			BlockGenerationIntervalFlag,
			BlockGenerationTimeLimitFlag,
			OpcodeComputationCostLimitFlag,
		},
	},
	{
		Name: "ACCOUNT",
		Flags: []cli.Flag{
			UnlockedAccountFlag,
			PasswordFileFlag,
			LightKDFFlag,
			KeyStoreDirFlag,
		},
	},
	{
		Name: "TXPOOL",
		Flags: []cli.Flag{
			TxPoolNoLocalsFlag,
			TxPoolAllowLocalAnchorTxFlag,
			TxPoolDenyRemoteTxFlag,
			TxPoolJournalFlag,
			TxPoolJournalIntervalFlag,
			TxPoolPriceLimitFlag,
			TxPoolPriceBumpFlag,
			TxPoolExecSlotsAccountFlag,
			TxPoolExecSlotsAllFlag,
			TxPoolNonExecSlotsAccountFlag,
			TxPoolNonExecSlotsAllFlag,
			TxPoolLifetimeFlag,
			TxPoolKeepLocalsFlag,
			TxResendIntervalFlag,
			TxResendCountFlag,
			TxResendUseLegacyFlag,
		},
	},
	{
		Name: "DATABASE",
		Flags: []cli.Flag{
			LevelDBCacheSizeFlag,
			PebbleDBCacheSizeFlag,
			SingleDBFlag,
			NumStateTrieShardsFlag,
			LevelDBCompressionTypeFlag,
			LevelDBNoBufferPoolFlag,
			RocksDBSecondaryFlag,
			RocksDBCacheSizeFlag,
			RocksDBDumpMallocStatFlag,
			RocksDBCompressionTypeFlag,
			RocksDBBottommostCompressionTypeFlag,
			RocksDBFilterPolicyFlag,
			RocksDBDisableMetricsFlag,
			RocksDBMaxOpenFilesFlag,
			RocksDBCacheIndexAndFilterFlag,
			DynamoDBTableNameFlag,
			DynamoDBRegionFlag,
			DynamoDBIsProvisionedFlag,
			DynamoDBReadCapacityFlag,
			DynamoDBWriteCapacityFlag,
			DynamoDBReadOnlyFlag,
			NoParallelDBWriteFlag,
			SenderTxHashIndexingFlag,
			DBNoPerformanceMetricsFlag,
			TxPruningFlag,
			TxPruningRetentionFlag,
			ReceiptPruningFlag,
			ReceiptPruningRetentionFlag,
		},
	},
	{
		Name: "DATABASE SYNCER",
		Flags: []cli.Flag{
			EnableDBSyncerFlag,
			DBHostFlag,
			DBPortFlag,
			DBNameFlag,
			DBUserFlag,
			DBPasswordFlag,
			EnabledLogModeFlag,
			MaxIdleConnsFlag,
			MaxOpenConnsFlag,
			ConnMaxLifeTimeFlag,
			BlockSyncChannelSizeFlag,
			DBSyncerModeFlag,
			GenQueryThreadFlag,
			InsertThreadFlag,
			BulkInsertSizeFlag,
			EventModeFlag,
			MaxBlockDiffFlag,
		},
	},
	{
		Name: "CHAINDATAFETCHER",
		Flags: []cli.Flag{
			EnableChainDataFetcherFlag,
			ChainDataFetcherMode,
			ChainDataFetcherNoDefault,
			ChainDataFetcherNumHandlers,
			ChainDataFetcherJobChannelSize,
			ChainDataFetcherChainEventSizeFlag,
			ChainDataFetcherMaxProcessingDataSize,
			ChainDataFetcherKASDBHostFlag,
			ChainDataFetcherKASDBPortFlag,
			ChainDataFetcherKASDBNameFlag,
			ChainDataFetcherKASDBUserFlag,
			ChainDataFetcherKASDBPasswordFlag,
			ChainDataFetcherKASCacheUse,
			ChainDataFetcherKASCacheURLFlag,
			ChainDataFetcherKASXChainIdFlag,
			ChainDataFetcherKASBasicAuthParamFlag,
			ChainDataFetcherKafkaBrokersFlag,
			ChainDataFetcherKafkaTopicEnvironmentFlag,
			ChainDataFetcherKafkaTopicResourceFlag,
			ChainDataFetcherKafkaReplicasFlag,
			ChainDataFetcherKafkaPartitionsFlag,
			ChainDataFetcherKafkaMaxMessageBytesFlag,
			ChainDataFetcherKafkaSegmentSizeBytesFlag,
			ChainDataFetcherKafkaRequiredAcksFlag,
			ChainDataFetcherKafkaMessageVersionFlag,
			ChainDataFetcherKafkaProducerIdFlag,
		},
	},
	{
		Name: "DATABASE MIGRATION",
		Flags: []cli.Flag{
			DstDbTypeFlag,
			DstDataDirFlag,
			DstSingleDBFlag,
			DstLevelDBCompressionTypeFlag,
			DstLevelDBCacheSizeFlag,
			DstNumStateTrieShardsFlag,
			DstDynamoDBTableNameFlag,
			DstDynamoDBRegionFlag,
			DstDynamoDBIsProvisionedFlag,
			DstDynamoDBReadCapacityFlag,
			DstDynamoDBWriteCapacityFlag,
			DstRocksDBSecondaryFlag,
			DstRocksDBCacheSizeFlag,
			DstRocksDBDumpMallocStatFlag,
			DstRocksDBCompressionTypeFlag,
			DstRocksDBBottommostCompressionTypeFlag,
			DstRocksDBFilterPolicyFlag,
			DstRocksDBDisableMetricsFlag,
			DstRocksDBMaxOpenFilesFlag,
			DstRocksDBCacheIndexAndFilterFlag,
		},
	},
	{
		Name: "STATE",
		Flags: []cli.Flag{
			TrieMemoryCacheSizeFlag,
			TrieBlockIntervalFlag,
			TriesInMemoryFlag,
			LivePruningFlag,
			LivePruningRetentionFlag,
		},
	},
	{
		Name: "CACHE",
		Flags: []cli.Flag{
			CacheTypeFlag,
			CacheScaleFlag,
			CacheUsageLevelFlag,
			MemorySizeFlag,
			TrieNodeCacheTypeFlag,
			NumFetcherPrefetchWorkerFlag,
			UseSnapshotForPrefetchFlag,
			TrieNodeCacheLimitFlag,
			TrieNodeCacheSavePeriodFlag,
			TrieNodeCacheRedisEndpointsFlag,
			TrieNodeCacheRedisClusterFlag,
			TrieNodeCacheRedisPublishBlockFlag,
			TrieNodeCacheRedisSubscribeBlockFlag,
		},
	},
	{
		Name: "CONSENSUS",
		Flags: []cli.Flag{
			ServiceChainSignerFlag,
			RewardbaseFlag,
		},
	},
	{
		Name: "NETWORKING",
		Flags: []cli.Flag{
			BootnodesFlag,
			ListenPortFlag,
			SubListenPortFlag,
			MultiChannelUseFlag,
			MaxConnectionsFlag,
			MaxPendingPeersFlag,
			TargetGasLimitFlag,
			NATFlag,
			NoDiscoverFlag,
			RWTimerWaitTimeFlag,
			RWTimerIntervalFlag,
			NetrestrictFlag,
			NodeKeyFileFlag,
			NodeKeyHexFlag,
			NetworkIdFlag,
			KairosFlag,
			MainnetFlag,
		},
	},
	{
		Name: "METRICS",
		Flags: []cli.Flag{
			MetricsEnabledFlag,
			PrometheusExporterFlag,
			PrometheusExporterPortFlag,
		},
	},
	{
		Name: "VIRTUAL MACHINE",
		Flags: []cli.Flag{
			VMEnableDebugFlag,
			VMLogTargetFlag,
			VMTraceInternalTxFlag,
			VMOpDebugFlag,
		},
	},
	{
		Name: "API AND CONSOLE",
		Flags: []cli.Flag{
			RPCEnabledFlag,
			HeavyDebugRequestLimitFlag,
			StateRegenerationTimeLimitFlag,
			RPCListenAddrFlag,
			RPCPortFlag,
			RPCCORSDomainFlag,
			RPCVirtualHostsFlag,
			RPCApiFlag,
			RPCGlobalGasCap,
			RPCGlobalEVMTimeoutFlag,
			RPCGlobalEthTxFeeCapFlag,
			RPCConcurrencyLimit,
			RPCNonEthCompatibleFlag,
			RPCExecutionTimeoutFlag,
			RPCIdleTimeoutFlag,
			RPCReadTimeout,
			RPCWriteTimeoutFlag,
			RPCUpstreamArchiveENFlag,
			UnsafeDebugDisableFlag,
			IPCDisabledFlag,
			IPCPathFlag,
			WSEnabledFlag,
			WSListenAddrFlag,
			WSPortFlag,
			WSApiFlag,
			WSAllowedOriginsFlag,
			WSMaxConnections,
			WSMaxSubscriptionPerConn,
			WSReadDeadLine,
			WSWriteDeadLine,
			GRPCEnabledFlag,
			GRPCListenAddrFlag,
			GRPCPortFlag,
			JSpathFlag,
			ExecFlag,
			PreloadJSFlag,
			MaxRequestContentLengthFlag,
			APIFilterGetLogsDeadlineFlag,
			APIFilterGetLogsMaxItemsFlag,
		},
	},
	{
		Name:  "LOGGING AND DEBUGGING",
		Flags: debug.Flags,
	},
	{
		Name: "SERVICECHAIN",
		Flags: []cli.Flag{
			ChildChainIndexingFlag,
			MainBridgeFlag,
			MainBridgeListenPortFlag,
			SubBridgeFlag,
			SubBridgeListenPortFlag,
			AnchoringPeriodFlag,
			SentChainTxsLimit,
			ParentChainIDFlag,
			VTRecoveryFlag,
			VTRecoveryIntervalFlag,
			ServiceChainAnchoringFlag,
			ServiceChainNewAccountFlag,
			ServiceChainParentOperatorTxGasLimitFlag,
			ServiceChainChildOperatorTxGasLimitFlag,
			KASServiceChainAnchorFlag,
			KASServiceChainAnchorPeriodFlag,
			KASServiceChainAnchorUrlFlag,
			KASServiceChainAnchorOperatorFlag,
			KASServiceChainAccessKeyFlag,
			KASServiceChainSecretKeyFlag,
			KASServiceChainXChainIdFlag,
			KASServiceChainAnchorRequestTimeoutFlag,
		},
	},
	{
		Name: "GAS PRICE ORACLE",
		Flags: []cli.Flag{
			GpoBlocksFlag,
			GpoPercentileFlag,
			GpoMaxGasPriceFlag,
		},
	},
	{
		Name: "MISC",
		Flags: []cli.Flag{
			GenKeyFlag,
			WriteAddressFlag,
			AutoRestartFlag,
			RestartTimeOutFlag,
			DaemonPathFlag,
			KESNodeTypeServiceFlag,
			SnapshotFlag,
			SnapshotCacheSizeFlag,
			SnapshotAsyncGen,
			DocRootFlag,
		},
	},
}

// CategorizeFlags classifies each flag into pre-defined flagGroups.
func CategorizeFlags(flags []cli.Flag) []FlagGroup {
	flagGroupsMap := make(map[string][]cli.Flag)
	isFlagAdded := make(map[string]bool) // Check duplicated flags

	// Find its group for each flag
	for _, flag := range flags {
		if isFlagAdded[flag.Names()[0]] {
			logger.Debug("a flag is added in the help description more than one time", "flag", flag.Names()[0])
			continue
		}

		// Find a group of the flag. If a flag doesn't belong to any groups, categorize it as a MISC flag
		group := flagCategory(flag, FlagGroups)
		flagGroupsMap[group] = append(flagGroupsMap[group], flag)
		isFlagAdded[flag.Names()[0]] = true
	}

	// Convert flagGroupsMap to a slice of FlagGroup
	flagGroups := []FlagGroup{}
	for group, flags := range flagGroupsMap {
		flagGroups = append(flagGroups, FlagGroup{Name: group, Flags: flags})
	}

	// Sort flagGroups in ascending order of name
	sortFlagGroup(flagGroups, uncategorized)

	return flagGroups
}

// sortFlagGroup sorts a slice of FlagGroup in ascending order of name,
// but an uncategorized group is exceptionally placed at the end.
func sortFlagGroup(flagGroups []FlagGroup, uncategorized string) []FlagGroup {
	sort.Slice(flagGroups, func(i, j int) bool {
		if flagGroups[i].Name == uncategorized {
			return false
		}
		if flagGroups[j].Name == uncategorized {
			return true
		}
		return flagGroups[i].Name < flagGroups[j].Name
	})

	// Sort flags in each group i ascending order of flag name.
	for _, group := range flagGroups {
		sort.Slice(group.Flags, func(i, j int) bool {
			return group.Flags[i].Names()[0] < group.Flags[j].Names()[0]
		})
	}

	return flagGroups
}

// flagCategory returns belonged group of the given flag.
func flagCategory(flag cli.Flag, fg []FlagGroup) string {
	for _, category := range fg {
		for _, flg := range category.Flags {
			if flg.Names()[0] == flag.Names()[0] {
				return category.Name
			}
		}
	}
	return uncategorized
}
