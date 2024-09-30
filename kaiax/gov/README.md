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
###  Init

- Dependencies:
  - headergov: To retrieve header governance parameters.
  - contractgov: To retrieve contract governance parameters.
- Notable dependents:
  - kaiax/valset: Provides committee size.
  - kaiax/reward: Provides parameters related to rewards.
  - kaiax/staking: Provides the useGini and minStake for the API.


###  Start and stop
This module does not have any background threads.

## Block processing

See [headergov](./headergov/README.md#Block-processing).

## APIs

Other APIs can be found in [headergov](./headergov/README.md#APIs).

### governance_getRewardsAccumulated

### kaia_getChainConfig

### kaia_getParams

### kaia_getRewards

### kaia_nodeAddress

## Getters

- `EffectiveParamSet(num)`: Returns the effective parameter set at the block `num`.
  ```
  EffectiveParamSet(num) -> ParamSet
  ```
