# API Documentation

This directory contains the API implementations for the Kaia blockchain node.

## API Files and Structs

### api/api_debug.go (DebugAPI)

Debug information and data inspection

- `debug_getBlockRlp`

### api/api_debug_util.go (DebugUtilAPI)

Debug utilities and system state management

- `debug_chaindbProperty`
- `debug_chaindbCompact`
- `debug_setHead`
- `debug_printBlock`

### api/api_eth.go (EthAPI)

Ethereum-compatible RPC methods including:

- `eth_etherbase`
- `eth_coinbase`
- `eth_hashrate`
- `eth_mining`
- `eth_getWork`
- `eth_submitWork`
- `eth_submitHashrate`
- `eth_getHashrate`
- `eth_newPendingTransactionFilter`
- `eth_newPendingTransactions`
- `eth_newBlockFilter`
- `eth_newHeads`
- `eth_logs`
- `eth_newFilter`
- `eth_getLogs`
- `eth_uninstallFilter`
- `eth_getFilterLogs`
- `eth_getFilterChanges`
- `eth_gasPrice`
- `eth_upperBoundGasPrice`
- `eth_lowerBoundGasPrice`
- `eth_maxPriorityFeePerGas`
- `eth_feeHistory`
- `eth_syncing`
- `eth_chainId`
- `eth_blockNumber`
- `eth_getBalance`
- `eth_getProof`
- `eth_getHeaderByNumber`
- `eth_getHeaderByHash`
- `eth_getBlockByNumber`
- `eth_getBlockByHash`
- `eth_getUncleByBlockNumberAndIndex`
- `eth_getUncleByBlockHashAndIndex`
- `eth_getUncleCountByBlockNumber`
- `eth_getUncleCountByBlockHash`
- `eth_getCode`
- `eth_getStorageAt`
- `eth_call`
- `eth_estimateGas`
- `eth_getBlockTransactionCountByNumber`
- `eth_getBlockTransactionCountByHash`
- `eth_getTransactionByBlockNumberAndIndex`
- `eth_getTransactionByBlockHashAndIndex`
- `eth_getRawTransactionByBlockNumberAndIndex`
- `eth_getRawTransactionByBlockHashAndIndex`
- `eth_getTransactionCount`
- `eth_getTransactionByHash`
- `eth_getRawTransactionByHash`
- `eth_getTransactionReceipt`
- `eth_getBlockReceipts`
- `eth_sendTransaction`
- `eth_fillTransaction`
- `eth_sendRawTransaction`
- `eth_sign`
- `eth_signTransaction`
- `eth_pendingTransactions`
- `eth_resend`
- `eth_accounts`
- `eth_createAccessList`

### api/api_kaia.go (KaiaAPI)

Kaia-specific RPC methods

- `kaia_gasPrice`
- `kaia_upperBoundGasPrice`
- `kaia_lowerBoundGasPrice`
- `kaia_protocolVersion`
- `kaia_maxPriorityFeePerGas`
- `kaia_feeHistory`
- `kaia_syncing`
- `kaia_encodeAccountKey`
- `kaia_decodeAccountKey`

### api/api_kaia_blockchain.go (KaiaBlockChainAPI)

Blockchain operations and queries

- `kaia_blockNumber`
- `kaia_chainId`
- `kaia_isContractAccount`
- `kaia_getBlockReceipts`
- `kaia_getBalance`
- `kaia_accountCreated`
- `kaia_getAccount`
- `kaia_forkStatus`
- `kaia_getHeaderByNumber`
- `kaia_getHeaderByHash`
- `kaia_getBlockByNumber`
- `kaia_getBlockByHash`
- `kaia_getCode`
- `kaia_getStorageAt`
- `kaia_getAccountKey`
- `kaia_isParallelDBWrite`
- `kaia_isSenderTxHashIndexingEnabled`
- `kaia_isConsoleLogEnabled`
- `kaia_call`
- `kaia_estimateComputationCost`
- `kaia_estimateGas`
- `kaia_createAccessList`
- `kaia_getProof`
- `kaia_getCypressCredit`

### api/api_kaia_transaction.go (KaiaTransactionAPI)

Transaction-related operations

- `kaia_getBlockTransactionCountByNumber`
- `kaia_getBlockTransactionCountByHash`
- `kaia_getTransactionByBlockNumberAndIndex`
- `kaia_getTransactionByBlockHashAndIndex`
- `kaia_getRawTransactionByBlockNumberAndIndex`
- `kaia_getRawTransactionByBlockHashAndIndex`
- `kaia_getTransactionCount`
- `kaia_getTransactionBySenderTxHash`
- `kaia_getTransactionByHash`
- `kaia_getDecodedAnchoringTransactionByHash`
- `kaia_getRawTransactionByHash`
- `kaia_getTransactionReceiptBySenderTxHash`
- `kaia_getTransactionReceipt`
- `kaia_getTransactionReceiptInCache`
- `kaia_sendTransaction`
- `kaia_sendTransactionAsFeePayer`
- `kaia_sendRawTransaction`
- `kaia_sign`
- `kaia_signTransaction`
- `kaia_signTransactionAsFeePayer`
- `kaia_pendingTransactions`
- `kaia_resend`
- `kaia_recoverFromTransaction`
- `kaia_recoverFromMessage`

### api/api_net.go (NetAPI)

Network-related RPC methods

- `net_listening`
- `net_peerCount`
- `net_peerCountByType`
- `net_version`
- `net_networkID`

### api/api_personal.go (PersonalAPI)

Personal account management methods

- `personal_listAccounts`
- `personal_listWallets`
- `personal_openWallet`
- `personal_deriveAccount`
- `personal_newAccount`
- `personal_replaceRawKey`
- `personal_importRawKey`
- `personal_unlockAccount`
- `personal_lockAccount`
- `personal_sendTransaction`
- `personal_sendTransactionAsFeePayer`
- `personal_sendAccountUpdate`
- `personal_sendValueTransfer`
- `personal_signTransaction`
- `personal_signTransactionAsFeePayer`
- `personal_sign`
- `personal_ecRecover`
- `personal_signAndSendTransaction`

### api/api_txpool.go (TxPoolAPI)

Transaction pool operations

- `txpool_content`
- `txpool_status`
- `txpool_inspect`

## Node APIs

### node/api_admin.go (AdminNodeAPI)

Node administration and management

- `admin_datadir`
- `admin_nodeInfo`
- `admin_peers`

### node/api_admin_network.go (AdminNetworkNodeAPI)

Network administration and peer management

- `admin_addPeer`
- `admin_peerEvents`
- `admin_removePeer`
- `admin_setMaxSubscriptionPerWSConn`
- `admin_startHTTP`
- `admin_startRPC`
- `admin_startWS`
- `admin_stopHTTP`
- `admin_stopRPC`
- `admin_stopWS`

### node/api_debug.go (DebugNodeAPI)

Node-level debug operations

- `debug_metrics`

### node/api_kaia.go (KaiaNodeAPI)

Node-level Kaia operations

- `kaia_clientVersion`
- `kaia_sha3`

## CN (Consensus Node) APIs

### node/cn/api_admin_chain.go (AdminChainCNAPI)

CN chain administration and management

- `admin_exportChain`
- `admin_getSpamThrottlerCandidateList`
- `admin_getSpamThrottlerThrottleList`
- `admin_getSpamThrottlerWhiteList`
- `admin_importChain`
- `admin_importChainFromString`
- `admin_nodeConfig`
- `admin_saveTrieNodeCacheToDisk`
- `admin_setSpamThrottlerWhiteList`
- `admin_spamThrottlerConfig`
- `admin_startSpamThrottler`
- `admin_startStateMigration`
- `admin_stateMigrationStatus`
- `admin_stopSpamThrottler`
- `admin_stopStateMigration`

### node/cn/api_debug.go (DebugCNAPI)

CN debug operations and data inspection

- `debug_dumpBlock`
- `debug_dumpStateTrie`
- `debug_getBadBlocks`
- `debug_getModifiedAccountsByHash`
- `debug_getModifiedAccountsByNumber`
- `debug_getModifiedStorageNodesByNumber`

### node/cn/api_debug_storage.go (DebugStorageCNAPI)

CN storage debug operations

- `debug_preimage`
- `debug_startCollectingTrieStats`
- `debug_startContractWarmUp`
- `debug_startWarmUp`
- `debug_stopWarmUp`
- `debug_storageRangeAt`

### node/cn/api_kaia.go (KaiaCNAPI)

CN-specific Kaia operations

- `kaia_rewardbase`

### node/cn/filters/api_kaia_filter.go (KaiaFilterAPI)

CN filter operations for events and logs

- `kaia_events`
- `kaia_getFilterChanges`
- `kaia_getFilterLogs`
- `kaia_getLogs`
- `kaia_logs`
- `kaia_newBlockFilter`
- `kaia_newFilter`
- `kaia_newHeads`
- `kaia_newPendingTransactionFilter`
- `kaia_newPendingTransactions`
- `kaia_uninstallFilter`

## Bootnode APIs

### cmd/kbn/api_bootnode.go (BootnodeAPI)

Bootnode operations

- `bootnode_getAuthorizedNodes`

### cmd/kbn/api_bootnode_registry.go (BootnodeRegistryAPI)

Bootnode registry operations

- `bootnode_createUpdateNodeOnDB`
- `bootnode_createUpdateNodeOnTable`
- `bootnode_deleteAuthorizedNodes`
- `bootnode_deleteNodeFromDB`
- `bootnode_deleteNodeFromTable`
- `bootnode_getNodeFromDB`
- `bootnode_getTableEntries`
- `bootnode_getTableReplacements`
- `bootnode_lookup`
- `bootnode_name`
- `bootnode_putAuthorizedNodes`
- `bootnode_readRandomNodes`
- `bootnode_resolve`

## DataSync APIs

### datasync/chaindatafetcher/api_chaindatafetcher.go (ChainDataFetcherAPI)

Chain data fetching operations

- `chaindatafetcher_getConfig`
- `chaindatafetcher_readCheckpoint`
- `chaindatafetcher_startFetching`
- `chaindatafetcher_startRangeFetching`
- `chaindatafetcher_status`
- `chaindatafetcher_stopFetching`
- `chaindatafetcher_stopRangeFetching`
- `chaindatafetcher_writeCheckpoint`

### datasync/downloader/api_kaia_downloader.go (KaiaDownloaderAPI)

Kaia downloader operations

- `kaia_subscribeSyncStatus`
- `kaia_syncing`

### datasync/downloader/api_kaia_downloader_sync.go (KaiaDownloaderSyncAPI)

Kaia downloader sync operations

- `kaia_syncStakingInfo`
- `kaia_syncStakingInfoStatus`
