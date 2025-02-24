# kaiax/staking

This module is responsible for tracking validator staking amounts and their address configurations.

## Concepts

- StakingInfo is a struct representing Validator staking information at a certain block including staked amount, reward address, and node address. It is primarily used to determine validator set and rewards distribution.
- StakingInfo summarizes the current AddressBook contract state, and all the staking contracts registered in the AddressBook, and their native token balances.
  - Since the Prague hardfork, the StakingInfo will include the [consensus liquidity](https://kips.kaia.io/KIPs/kip-226) information from the CLRegistry.
- When processing a block at `num`, the StakingInfo from a historic block state is used and the historic block (i.e. source block) is determined by the following:

  - If `num` is before Kaia hardfork, then StakingInfo is drawn from the beginning of the previous staking interval. Note that if the `num` is a multiple of StakingInterval, the staking info is drawn from two epochs ahead (e.g. in the example below, the staking info at block 3000 is drawn from block 1000).
    ```go
    SourceNum(num) = RoundDown(num - 1, StakingInterval) - StakingInterval
    RoundDown(n, p) = n - (n % p)
    ```
  - If `num` is after Kaia hardfork, then StakingInfo is drawn from the previous block.
    ```go
    SourceNum(num) = num - 1
    ```
  - Example

    ```
    Given StakingInterval = 1000,
    Given Kaia hardfork at block 3456,

    num      source
    2001     1000 // RoundDown(2000 , 1000) - 1000 = 1000
    2002     1000
    2999     1000
    3000     1000 // RoundDown(2999 , 1000) - 1000 = 1000

    3001     2000 // RoundDown(3000 , 1000) - 1000 = 2000
    3002     2000
    3455     2000 // before Kaia HF, using the interval rule

    3456     3455 // after Kaia HF, using the previous block rule
    3457     3456
    ```

- StakingInterval a governance parameter (`reward.stakingupdateinterval`) that is first defined in the ChainConfig (`ChainConfig.Reward.StakingUpdateInterval`) at genesis and never changes afterwards.
- Kaia has two treasury fund addresses.
  - KEF (Kaia Ecosystem Fund). Stored in the AddressBook contract's `kirContractAddress` variable. This variable previously held the KCF, KIR addresses.
  - KIF (Kaia Infrastructure Fund). Stored in the AddressBook contract's `pocContractAddress` variable. This variable previously held the KFF, KGF, PoC addresses.

## Persistent Schema

- `StakingInfo(sourceNum)` The StakingInfo captured from the states at the block `sourceNum`. Persisted every StakingInterval before Kaia hardfork.
  ```
  "stakingInfo" || Uint64LE(num) => JSON.Marshal(StakingInfo)
  ```

## In-memory Structures

### StakingInfo

The staking info to be used for block processing.

```go
type StakingInfo struct {
  // The source block number where the staking info is captured.
  SourceBlockNum uint64 `json:"blockNum"`

  // The AddressBook triplets
  NodeIds          []common.Address `json:"councilNodeAddrs"`
  StakingContracts []common.Address `json:"councilStakingAddrs"`
  RewardAddrs      []common.Address `json:"councilRewardAddrs"`

  // Treasury fund addresses
  KEFAddr common.Address `json:"kefAddr"` // KEF contract address (or KCF, KIR)
  KIFAddr common.Address `json:"kifAddr"` // KIF contract address (or KFF, KGF, PoC)

  // Staking amounts of each staking contracts, in KAIA, rounded down.
  StakingAmounts []uint64 `json:"councilStakingAmounts"`

  // The consensus liquidity information
  CLStakingInfos *CLStakingInfos `json:"clStakingInfos"`
}

type CLStakingInfo struct {
	CLNodeId        common.Address `json:"clNodeId"`
	CLPoolAddr      common.Address `json:"clPoolAddr"`
	CLStakingAmount uint64         `json:"clStakingAmount"`
}

type CLStakingInfos []*CLStakingInfo
```

- `ConsolidatedNodes()` returns the nodes consolidated by the duplicating reward addresses. Note that the AddressBook entries with the same reward address are considered a single validator.
  - If `CLStakingInfo` exists for a validator after the Prague hardfork, the `ConsolidatedNodes()` will consolidate it. Note that `CLStakingInfo` has different reward address so the `CLNodeId` is used to consolidate.
- `Gini(minStake)` returns the gini coefficient of the staking amounts that are no less than `minStake`.

### StakingInfoResponse

The response type for `kaia_getStakingInfo` and `governance_getStakingInfo`. Adds additional fields for backward compatibility.

```go
type StakingInfoResponse struct {
  StakingInfo

  // Legacy treasury fund address fields for backward-compatibility
  KIRAddr common.Address `json:"KIRAddr"` // = KEFAddr
  PoCAddr common.Address `json:"PoCAddr"` // = KIFAddr
  KCFAddr common.Address `json:"kcfAddr"` // = KEFAddr
  KFFAddr common.Address `json:"kffAddr"` // = KIFAddr

  // Computed fields
  UseGini bool `json:"useGini"` // Whether the gini coefficient is used at the requested block number
  Gini float64 `json:"gini"` // The gini coefficient at the requested block number. Returned regardless of `UseGini` value.
}
```

## Module lifecycle

### Init

- Dependencies:
  - ChainDB: Raw key-value database to access this module's persistent schema.
  - ChainConfig: Holds the StakingInterval value at genesis.
  - Chain: Provides the blocks and states.

### Start and stop

This module does not have any background threads.

## Block processing

### Consensus

This module does not have any consensus-related block processing logic.

### Execution

This module makes sure that the corresponding StakingInfo is persisted, if applicable.

### Rewind

Upon rewind, this module deletes the related persistent data and flushes the in-memory cache.

## APIs

### kaia_getStakingInfo, governance_getStakingInfo

Query the StakingInfo to be used for the block `num`.

- Parameters
  - `num`: block number or hash
- Returns
  - `StakingInfoResponse`
- Example

```json
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"kaia_getStakingInfo","params":[
    "latest"
  ]}' | jq

{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "blockNum": 165145974,
    "councilNodeAddrs": [
      "0x99fb17d324fa0e07f23b49d09028ac0919414db6",
      "0x571e53df607be97431a5bbefca1dffe5aef56f4d",
      "0xb74ff9dea397fe9e231df545eb53fe2adf776cb2",
      "0x5cb1a7dccbd0dc446e3640898ede8820368554c8"
    ],
    "councilStakingAddrs": [
      "0x12fa1ab4c3e17c1c08c1b5a945c864c8e8bf707e",
      "0xfd56604f1a20268ff7a0eab2ab48e25ee1e0f653",
      "0x1e0f6aaa9baa6081dc4910a854eebf8854c262ab",
      "0x5e6988415ebe0f6b088f5a676003ba60f572875a"
    ],
    "councilRewardAddrs": [
      "0xb2bd3178affccd9f9f5189457f1cad7d17a01c9d",
      "0x6559a7b6248b342bc11fbcdf9343212bbc347edc",
      "0x82829a60c6eac4e3e9d6ed00891c69e88537fd4d",
      "0xf90675a56a03f836204d66c0f923e00500ddc90a"
    ],
    "useGini": false,
    "gini": 0.25,
    "councilStakingAmounts": [
      5000001,
      15000000,
      5000001,
      5000000
    ],
    "KIRAddr": "0x819d4b7245164e6a94341f4b5c2ae587372bb669",
    "PoCAddr": "0x8436e5bd1a6d622c278c946e2f8988a26136a16f",
    "kcfAddr": "0x819d4b7245164e6a94341f4b5c2ae587372bb669",
    "kffAddr": "0x8436e5bd1a6d622c278c946e2f8988a26136a16f",
    "kefAddr": "0x819d4b7245164e6a94341f4b5c2ae587372bb669",
    "kifAddr": "0x8436e5bd1a6d622c278c946e2f8988a26136a16f"
  }
}
```

## Getters

- GetStakingInfo: Returns the StakingInfo for the block `num`.
  ```
  GetStakingInfo(num) -> StakingInfo
  ```
