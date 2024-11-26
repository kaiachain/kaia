# kaiax/gov/headergov

This module is responsible for providing the governance parameter set from **header governance** at a given block number.

## Concepts

Please read [gov module](../README.md) first.

### Key Concepts

- _vote_: a vote for a parameter change.
- _ratification_: votes in an epoch being declared as ratified (accepted).
- _epoch_: a fixed period for header governance ratification. Defined in the genesis block. In Mainnet and Kairos, epoch is 604800 blocks (1 week).
- _epoch index_: the index of the epoch, starting from 0. Given a block number `N`, its epoch index is `floor(N / epoch)`. In other words, all blocks in `[k*epoch, (k+1)*epoch - 1]` belong to the `k`-th epoch.
- _effective parameter set at blockNum_: the governance parameter set that are effective when mining the given block.

### Header governance

Header governance is the process of changing the governance parameters among members of the GC via block header.
This module writes and interprets block header's `Vote` and `Governance` fields. Note that this module writes validator addition/removal votes to `header.Vote`, but does not read them from `header.Vote` which is left to the kaiax/valset module.

```
      vote         ratified    vote   vote   ratified
       V                V       V       V    V
   |---+----------------|-------+-------+----|

   *.....0th epoch......O
                        *......1st epoch.....O
                                             *....2nd epoch....
```

A vote is initiated by a governing member of the GC casting a vote.
The vote is cast by the `governance_vote` API, which is inscribed in the block header `header.Vote` when the node that received the API call becomes the proposer.
The voter casts a vote as a tuple `(parameter name, new value)`.
The vote is collected for the entire epoch.

At every epoch block (i.e., `k*epoch` blocks), the node that becomes the proposer will check if the vote has been ratified.
If the vote is ratified, the node will announce the ratification in the block header `header.Governance`.
In other words, `header.Governance` can contain data only at epoch blocks.
If there are no votes in an epoch, the next starting block of the epoch will have an empty `header.Governance`.
It contains a JSON object of `{name: value}` for each ratified parameter.

The ratification condition is determined by the `governance.governancemode` parameter. Mainnet and Kairos both operate in `single` mode. There are two governance modes:

- `none` mode: all members of the GC can vote. For each governance parameter, the last vote in the epoch will be ratified.
- `single` mode: only one member of the GC, stipulated in the parameter `governance.governingnode`, can vote. The vote will be ratified if it is the only vote in the epoch.

Parameter change ratified at `k*epoch` block takes effect starting from `(k+1)*epoch` block.
It is worth noting that the effective time of the ratification is `(k+1)*epoch + 1` before Kore.

### Reading a parameter set

The effective parameter set at block `N` (in `k`-th epoch) is determined as follows:

- Collect all the ratified parameters from 0-th to `k-1`-th epoch. In case of duplication, recent ratification is prioritized.
  - `k-1` is calculated by [PrevEpochStart](./impl/getter.go#L41).
- For each parameter, take the last ratified value. If a parameter has never been ratified, use the default value as a fallback.

This is the description of `EffectiveParamSet(N)`, which is implemented [here](./impl/getter.go#L9).

For example, given `epoch=1000`, assume that `header` is as follows:

```
num  |  header
--------------
0    |  Governance: {"governance.unitprice": 25 kei, "reward.mintingamount": 5 KAIA}
400  |  Vote: ("governance.unitprice", 50 kei)
500  |  Vote: ("governance.unitprice", 100 kei)
600  |  Vote: ("reward.mintingamount", 10 KAIA)
1000 |  Governance: {"governance.unitprice": 100 kei, "reward.mintingamount": 10 KAIA}
```

Then, this module will return the effective parameter set as follows:

```
num  |  effective parameter set at num
------------------------------------
1999 |  {"governance.unitprice": 25 kei, "reward.mintingamount": 5 KAIA}
2000 |  {"governance.unitprice": 100 kei, "reward.mintingamount": 10 KAIA}
2001 |  same as above
...  |  same as above
```

## Persistent Schema

- `governanceVoteDataBlockNums`: The block numbers whose header contains the vote data.
  ```
  "governanceVoteDataBlockNums" => JSON.Marshal([num1, num2, ...])
  ```

- `governanceDataBlockNums`: The block numbers whose header contains the governance data.
  ```
  "governanceDataBlockNums" => JSON.Marshal([num1, num2, ...])
  ```
  
- `lowestVoteScannedBlockNumKey`: The lowest block number whose vote data is scanned. That is, only vote block numbers which are greater than or equal to this value are stored in `governanceVoteDataBlockNums`. Grows downwards by `epoch` blocks.
  ```
  "governanceLowestVoteScannedBlockNum" => Uint64BE(num)
  ```

- `legacyIdxHistoryKey`: legacy governance data block numbers.
  ```
  "governanceIdxHistory" => JSON.Marshal([num1, num2, ...])
  ```

## In-memory Structures

### VoteData

`VoteData` is used for storing `header.Vote` in memory.
All `VoteData` values are canonicalized and format-checked.

See [vote.go](./vote.go).

- `ToVoteBytes()` returns the serialized bytes which is written in `header.Vote`.

### GovData

`GovData` is used for storing header's `Governance` in memory.
All `GovData` values are canonicalized and format-checked.
Unlike `VoteData`, vote-forbidden parameters are allowed for parsing the genesis block.

See [gov.go](./gov.go).

- `ToGovBytes()` returns the serialized bytes of the governance which is written in `header.Governance`.

### History

`History` is used for obtaining the parameter set at a given block number.

See [history.go](./history.go).

### HeaderCache

`HeaderCache` is used for caching DB data in memory.
Cache is always fully synced with the DB, so there's no need to write from DB.
In that sense, writing to the cache will write to DB as well.

See [cache.go](./cache.go).

### VotesResponse

The response type for `governance_votes`.

See [impl/api.go](./impl/api.go).

### MyVotesResponse

The response type for `governance_myVotes`. `MyVotes` contains votes that the current node has casted in this epoch, and votes waiting to be casted when the node becomes a proposer.

See [impl/api.go](./impl/api.go).

### StatusResponse

The response type for `governance_status`.

See [impl/api.go](./impl/api.go).

## Module lifecycle

### Init

- Dependencies:
  - ChainDB: Raw key-value database to access this module's persistent schema.
  - ChainConfig: Holds the epoch value at genesis.
  - Chain: Provides the blocks and states.
  - NodeAddress: Provides `governance_nodeAddress` API. Facilitates checks in `governance_vote`.

### Start and stop

This module does not have any background threads.

## Block processing

### Consensus

#### PrepareHeader

This module writes `header.Vote` and `header.Governance` during the block processing.
Specifically, it writes `header.Vote` if `governance_vote` API is called on this node.
It writes `header.Governance` at an epoch block if there are any ratified votes in the previous epoch.

#### VerifyHeader

This module verifies `header.Vote` and `header.Governance` during the block processing.
Specifically, it checks the following for `header.Vote` if it exists:

- The voter is the block proposer.
- The voter has the right to vote.
- The voted parameter does not break the consistency.

It checks the following for `header.Governance` if it exists:

- The block is an epoch block.
- The ratification is built based on the votes in the previous epoch.

#### FinalizeHeader

This module does not have any block processing logic at `FinalizeHeader`.

### Execution

This module updates cache and DB based on `header.Vote` and `header.Governance`.

### Rewind

Upon rewind, this module deletes the related persistent data and flushes the in-memory cache.

## APIs

### governance_vote

Cast a vote for a parameter.

- Parameters
  - `name`: name of the parameter
  - `value`: new value of the parameter
- Returns
  - `string`: confirmation message
  - `error`: error
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_vote","params":[
    "governance.unitprice",
    100
  ]}' | jq '.result'
"(kaiax) Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate."
```

### governance_idxCache

Returns all vote block numbers from cache. The API name is retained for legacy compatibility.

- Parameters: none
- Returns
  - `[]uint64`: block numbers
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_idxCache","params":[]}' | jq '.result'
[100, 200, 300]
```

### governance_votes

Returns all votes in the epoch that the given block number belongs to.

- Parameters:
  - `num`: block number
- Returns
  - `VotesResponse`: votes
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_votes","params":[
    100
  ]}' | jq '.result'
[
  {
    "BlockNum": 102,
    "Key": "governance.unitprice",
    "Value": 100
  }
]
```

### governance_myVotes

Returns `MyVotes` which contains votes that the current node has casted in this epoch, and votes waiting to be casted when the node becomes a proposer.

- Parameters: none
- Returns
  - `MyVotesResponse`: votes with casted flag
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_myVotes","params":[]}' | jq '.result'
[
  {
    "BlockNum": 102,
    "Key": "governance.unitprice",
    "Value": 100,
    "Casted": true
  }
]

```

### governance_status

Returns in-memory data of this module.

- Parameters: none
- Returns
  - `StatusResponse`: status
- Example

```
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"governance_status","params":[]}' | jq '.result'
{
  "groupedVotes": {
    "3": {
      "102": {
        "Voter": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
        "Name": "governance.unitprice",
        "Value": 100
      }
    },
    "4": {
      "143": {
        "Voter": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
        "Name": "governance.unitprice",
        "Value": 400
      }
    },
    "7": {
      "215": {
        "Voter": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
        "Name": "governance.unitprice",
        "Value": 500
      }
    }
  },
  "governances": {
    "0": {
      "governance.deriveshaimpl": 2,
      "governance.governancemode": "single",
      "governance.governingnode": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
      "governance.unitprice": 25000000000,
      "istanbul.committeesize": 13,
      "istanbul.epoch": 30,
      "istanbul.policy": 2,
      "reward.deferredtxfee": true,
      "reward.minimumstake": "5000000",
      "reward.mintingamount": "9600000000000000000",
      "reward.proposerupdateinterval": 30,
      "reward.ratio": "34/54/12",
      "reward.stakingupdateinterval": 60,
      "reward.useginicoeff": false
    },
    "120": {
      "governance.unitprice": 100
    },
    "150": {
      "governance.unitprice": 400
    },
    "240": {
      "governance.unitprice": 500
    }
  },
  "govHistory": {
    "0": {
      "GovernanceMode": "single",
      "GoverningNode": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
      "GovParamContract": "0x0000000000000000000000000000000000000000",
      "CommitteeSize": 13,
      "ProposerPolicy": 2,
      "Epoch": 30,
      "Ratio": "34/54/12",
      "Kip82Ratio": "20/80",
      "StakingUpdateInterval": 60,
      "ProposerUpdateInterval": 30,
      "MintingAmount": 9600000000000000000,
      "MinimumStake": 5000000,
      "UseGiniCoeff": false,
      "DeferredTxFee": true,
      "LowerBoundBaseFee": 25000000000,
      "UpperBoundBaseFee": 750000000000,
      "GasTarget": 30000000,
      "MaxBlockGasUsedForBaseFee": 60000000,
      "BaseFeeDenominator": 20,
      "DeriveShaImpl": 2,
      "UnitPrice": 25000000000
    },
    "120": {
      "GovernanceMode": "single",
      "GoverningNode": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
      "GovParamContract": "0x0000000000000000000000000000000000000000",
      "CommitteeSize": 13,
      "ProposerPolicy": 2,
      "Epoch": 30,
      "Ratio": "34/54/12",
      "Kip82Ratio": "20/80",
      "StakingUpdateInterval": 60,
      "ProposerUpdateInterval": 30,
      "MintingAmount": 9600000000000000000,
      "MinimumStake": 5000000,
      "UseGiniCoeff": false,
      "DeferredTxFee": true,
      "LowerBoundBaseFee": 25000000000,
      "UpperBoundBaseFee": 750000000000,
      "GasTarget": 30000000,
      "MaxBlockGasUsedForBaseFee": 60000000,
      "BaseFeeDenominator": 20,
      "DeriveShaImpl": 2,
      "UnitPrice": 100
    },
    "150": {
      "GovernanceMode": "single",
      "GoverningNode": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
      "GovParamContract": "0x0000000000000000000000000000000000000000",
      "CommitteeSize": 13,
      "ProposerPolicy": 2,
      "Epoch": 30,
      "Ratio": "34/54/12",
      "Kip82Ratio": "20/80",
      "StakingUpdateInterval": 60,
      "ProposerUpdateInterval": 30,
      "MintingAmount": 9600000000000000000,
      "MinimumStake": 5000000,
      "UseGiniCoeff": false,
      "DeferredTxFee": true,
      "LowerBoundBaseFee": 25000000000,
      "UpperBoundBaseFee": 750000000000,
      "GasTarget": 30000000,
      "MaxBlockGasUsedForBaseFee": 60000000,
      "BaseFeeDenominator": 20,
      "DeriveShaImpl": 2,
      "UnitPrice": 400
    },
    "240": {
      "GovernanceMode": "single",
      "GoverningNode": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
      "GovParamContract": "0x0000000000000000000000000000000000000000",
      "CommitteeSize": 13,
      "ProposerPolicy": 2,
      "Epoch": 30,
      "Ratio": "34/54/12",
      "Kip82Ratio": "20/80",
      "StakingUpdateInterval": 60,
      "ProposerUpdateInterval": 30,
      "MintingAmount": 9600000000000000000,
      "MinimumStake": 5000000,
      "UseGiniCoeff": false,
      "DeferredTxFee": true,
      "LowerBoundBaseFee": 25000000000,
      "UpperBoundBaseFee": 750000000000,
      "GasTarget": 30000000,
      "MaxBlockGasUsedForBaseFee": 60000000,
      "BaseFeeDenominator": 20,
      "DeriveShaImpl": 2,
      "UnitPrice": 500
    }
  },
  "nodeAddress": "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
  "myVotes": []
}
```

## Getters

- `EffectiveParamSet(num)`: Returns the effective parameter set at the block `num`.

  ```
  EffectiveParamSet(num) -> ParamSet
  ```

- `EffectiveParamsPartial(num)`: Returns only the parameters effective by header governance, which is the union of `header.governance` from block 0 to `num`. It is used for assembling parameters in a gov module.
  ```
  EffectiveParamsPartial(num) -> PartialParamSet
  ```
