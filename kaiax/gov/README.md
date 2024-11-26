# kaiax/gov

This submodule is responsible for providing the governance parameter set at a given block number.

## Concepts

A governance parameter is a specific configurable value that influences the behavior or rules of the blockchain network.
It's essentially a setting that can be adjusted to fine-tune how the network operates.
These parameters could control stuff like transaction fees, inflation rate, reward distribution, etc.
A governance parameter has a name and a value; see [./param.go](./param.go).

`EffectiveParams(blockNum)` returns the governance parameter set that are used for mining the given block.
For example, if `EffectiveParams(10000)` is `{UnitPrice: 25 kei, ...}`, it indicates that the unit price will be 25 kei when mining the 10000th block.

Here are the list of governance parameters:

```
<mutable parameters>
governance.deriveshaimpl
governance.governingnode
governance.govparamcontract
governance.unitprice
istanbul.committeesize
kip71.basefeedenominator
kip71.gastarget
kip71.lowerboundbasefee
kip71.maxblockgasusedforbasefee
kip71.upperboundbasefee
reward.kip82ratio

reward.mintingamount
reward.ratio

<immutable parameters - Mainnet configuration>
governance.governancemode: single
istanbul.epoch: 604800
istanbul.policy: 2
reward.deferredtxfee: true
reward.minimumstake: 5000000 (KAIA)
reward.proposerupdateinterval: 3600
reward.stakingupdateinterval: 86400
reward.useginicoeff: true
```

This module utilizes [header governance](./headergov/README.md) and [contract governance](./contractgov/README.md) underneath to fetch the parameter set and to handle governance parameter updates.

```
EffectiveParams(blockNum):
    ret := defaultParamSet()
    merge ret with HeaderGov.EffectiveParams(blockNum)
    if blockNum is post-Kore-HF:
        merge ret with ContractGov.EffectiveParams(blockNum)
    return ret
```

## Persistent Schema

See [headergov schema](./headergov/README.md#persistent-schema).

## In-memory Structures

## Module lifecycle

### Init

- Dependencies:
  - headergov: To retrieve header governance parameters.
  - contractgov: To retrieve contract governance parameters.
- Notable dependents:
  - kaiax/valset: Provides committee size.
  - kaiax/reward: Provides parameters related to rewards.
  - kaiax/staking: Provides the useGini and minStake for the API.

### Start and stop

This module does not have any background threads.

## Block processing

See [headergov](./headergov/README.md#Block-processing).

## APIs

Other APIs can be found in [headergov](./headergov/README.md#APIs).

### governance_getRewardsAccumulated

Returns the accumulated rewards for the given block range.

- Parameters:
  - `lower`: the starting block number
  - `upper`: the ending block number
- Returns
  - `AccumulatedRewardsResponse`: accumulated rewards
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_getRewardsAccumulated","params":["0x10", "0x20"]}' | jq '.result'
{
  "firstBlockTime": "2024-09-30 15:57:35 +0900 KST",
  "lastBlockTime": "2024-09-30 15:57:51 +0900 KST",
  "firstBlock": 16,
  "lastBlock": 32,
  "totalMinted": 163200000000000000000,
  "totalTxFee": 0,
  "totalBurntTxFee": 0,
  "totalProposerRewards": 163200000000000000000,
  "totalStakingRewards": 0,
  "totalKIFRewards": 0,
  "totalKEFRewards": 0,
  "rewards": {
    "0x70997970c51812dc3a010c7d01b50e0d17dc79c8": 105600000000000000000,
    "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266": 57600000000000000000
  }
}
```

### kaia_getChainConfig

Returns the chain config at the block `num`.

- Parameters:
  - `num`: block number
- Returns
  - `ChainConfig`: chain config
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_getChainConfig","params":[]}' | jq '.result'
=> TODO
```

### governance_getParams

Returns the effective parameter set at the block `num`.

- Parameters:
  - `num`: block number
- Returns
  - `map[ParamName]any`: parameter set
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_getParams","params":[
    100
  ]}' | jq '.result'
{
  "governance.deriveshaimpl": 2,
  "governance.governancemode": "single",
  "governance.governingnode": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
  "governance.govparamcontract": "0x0000000000000000000000000000000000000000",
  "governance.unitprice": 25000000000,
  "istanbul.committeesize": 13,
  "istanbul.epoch": 30,
  "istanbul.policy": 2,
  "kip71.basefeedenominator": 20,
  "kip71.gastarget": 30000000,
  "kip71.lowerboundbasefee": 25000000000,
  "kip71.maxblockgasusedforbasefee": 60000000,
  "kip71.upperboundbasefee": 750000000000,
  "reward.deferredtxfee": true,
  "reward.kip82ratio": "20/80",
  "reward.minimumstake": "5000000",
  "reward.mintingamount": "9600000000000000000",
  "reward.proposerupdateinterval": 1,
  "reward.ratio": "34/54/12",
  "reward.stakingupdateinterval": 1,
  "reward.useginicoeff": false
}
```

### kaia_getRewards

Returns the rewards at the block `num`.

- Parameters:
  - `num`: block number
- Returns
  - `RewardSpec`: rewards
- Example
```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"kaia_getRewards","params":["0x10"]}' | jq '.result'
{
  "minted": 9600000000000000000,
  "totalFee": 0,
  "burntFee": 0,
  "proposer": 9600000000000000000,
  "stakers": 0,
  "kif": 0,
  "kef": 0,
  "rewards": {
    "0x70997970c51812dc3a010c7d01b50e0d17dc79c8": 9600000000000000000
  }
}
```


### kaia_nodeAddress, governance_nodeAddress

Returns the node address.

- Parameters: none
- Returns
  - `address`: node address
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_nodeAddress","params":[]}' | jq '.result'
"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
```


## Getters

- `EffectiveParamSet(num)`: Returns the effective parameter set at the block `num`.
  ```
  EffectiveParamSet(num) -> ParamSet
  ```
