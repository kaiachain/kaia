# kaiax/valset

This module is responsible for getting council and calculating committee or proposer.

## Concepts
| Components                  | [Block 0]                          | [Block N]  <br/>"Add Validator i" ,"Remove Validator j"                                    |
|-----------------------------|------------------------------------|--------------------------------------------------------------------------------------------|
| demoted validators(Block)   | None                               | Council(N-1).demotedByQualification(*state(N-1))                                           | 
| qualified validators(Block) | GenesisCouncil                     | Council(N-1) - demoted(N)                                                                  |
| └ Committee(Block,Round)    | GenesisCouncil                     | qualified(N).sublist(Round')                                                               |
| └└ proposer(Block,Round)    | GenesisCouncil's<br/>first element | *Committee(N,Round').proposer(Round')                                                      |
| Council(Block)              | GenesisCouncil                     | Council(N-1) <br/>+ validator(i), i is not in Council <br/>- validator(j), j is in Council |
*Committee(N,Round').proposer(Round'): A Proposer is calculated first than the committee unless the Randao hardfork is activated. 
Nevertheless, proposer is ensured to be a committee member, so it is represented that the proposer is a member of committee.<br/>
*state(N-1): It's a result of processing block N-1 including state transition. The results are govParamSet, stakingInfo, and author.

### council: a set of registered CN
A `Council(N)` refers to the council determined after the block `N` is executed, based on `header(N).Vote` and `gov.EffectiveParamsAt(N)`.
In other words, to reach consensus on block N, `Council(N-1)` is used as a superset of Committee(N) or proposer(N).
Additionally, the genesis council is directly restored from the extraData of the genesis block(block 0).

### qualified validators: a subset of council members who are qualified to be a committee
To determine a committee or proposer at block N, the council members must be filtered based on their qualification.
* Qualified validators: Validators who meet the criteria to be a committee member.
* Demoted validators: The remaining council members after excluding the qualified validators.

Committee members are selected only from the qualified validators, and the qualifications vary depending on the proposer policy:
* clique consensus: Uses default proposer policy (`round-robin`).
  * round-robin: all council members are qualified
* istanbul consensus: Supports three proposer policies(`round-robin`, `sticky`, `weighted-random`).
  * round-robin: all council members are qualified
  * sticky: all council members are qualified
  * weighted-random: Validators who staked more than **[minimum staking amount]** are qualified.
    * If no nodes meet this requirement, all nodes are considered qualified. In single mode, governingnode is qualified without any conditions.

### committee: a subset of qualified validators
A `Committee(N, R)` is randomly selected as a subset of `Council(N-1)` for mining the block `N` at the round `R`.

The `Committee(N, R)` is chosen based on the followings: `Header(N-1)`, `StakingInfo(N-1)`, `gov.EffectiveParamsAt(N)`, and `Council(N-1)`.
The committee size is determined by the governance parameter value, "istanbul.committeeSize", which indicates the maximum size of committee.
Additionally, the genesis block doesn't have previous block, so genesis committee is copied from the genesis council.

#### Committee selection logic
- Selection logic involves a shuffling algorithms: `SimpleRandomCommittee` and `RandaoRandomCommittee` which differs depending on whether it's before/after Randao hardfork and the proposer policy.
  - More specifically, `RandaoCommittee` logic is used when `policy.IsDefaultSet() || (policy.IsWeightedCouncil() && !rules.IsRandao)`, which corresponds to Mainnet and Kairos.
- Calculate the shuffling seed
  - `SimpleRandomCommittee`: Calculate the seed using prevHash. `seed = int64(binary.BigEndian.Uint64(prevHash.Bytes()[:8]))`.
  - `RandaoRandomCommittee`: Calculate the seed using mixHash. `seed = int64(binary.BigEndian.Uint64(mixHash[:8]))`
- Shuffle the qualified validators
  - `SimpleRandomCommittee`: To secure the proposers, extract the current round's proposer and next distinct proposer which has different address. Shuffle the rest and attach the two proposers.
  - `RandaoRandomCommittee`: Shuffle the qualified validators.
- After shuffling, the council is "sliced" up to "istanbul.committeesize".
```
Condition: proposerIdx = 3, nextProposerIdx = 7, committeesize = 6, qualified validators = [0,1,2,3,4,5,6,7,8,9]

SimpleRandomCommittee:
- Step1. secure the proposer and the next distinct proposer: proposers = [3,7], council = [0,1,2,4,5,6,8,9]
- Step2. shuffle the council: proposers = [3,7], council = [4,5,6,8,9,0,1,2]
- Step3. merge. council = [3,7,4,5,6,8,9,0,1,2]
- Step4. slice by committeesize: committee = [3,7,4,5,6,8]

RandaoCommittee:
- Step1. shuffle the council: council = [4,7,5,6,8,3,9,0,1,2]
- Step2. slice by committeesize: committee = [4,7,5,6,8,3]
```

### proposer: a member of committee who proposes the block
A `Proposer(N, R)` proposes the block N at round R who are selected from the committee.
Note that proposer is called an "author" once the block is created.

#### ProposerPolicy
Proposer is selected based on the governance parameter "istanbul.proposerpolicy".
The mechanism of the same policy works differently on the hardfork.

Here are the list of policies: (the number in the paranthesis indicates the value of "istanbul.proposerpolicy")
- `RoundRobin` (0): every council member takes each turn. That is, `council(N-1)[(prevAuthorIdx+round)%len(council(N-1))]`
  <details>
    <summary>RoundRobin Example</summary>
  
  ```
  proposers=[0,1,2,3]
  
  block  | round | proposer
  -------------------------
   1     |   0   |    0
   2     |   0   |    1
   3     |   0   |    2
   4     |   0   |    3
   5     |   0   |    0
   6     |   0   |    1
   7     |   0   |    2
   8     |   0   |    3
   9     |   0   |    0
  10     |   0   |    1
  11     |   0   |    2        << round change
  11     |   1   |    3        << round change
  11     |   2   |    0
  12     |   0   |    1
  13     |   0   |    2
  14     |   0   |    3
  ```
  </details>

- `Sticky` (1): every council member takes each turn. That is, `council(N-1)[(prevAuthorIdx+round+1)%len(council(N-1))]`
  <details>
    <summary>Sticky Example</summary>

  ```
  proposers=[0,1,2,3]
  
  block  | round | proposer
  -------------------------
   1     |   0   |    0
   2     |   0   |    0
   3     |   0   |    0
   4     |   0   |    0
   5     |   0   |    0
   6     |   0   |    0
   7     |   0   |    0
   8     |   0   |    0
   9     |   0   |    0
  10     |   0   |    0
  11     |   0   |    0        << round change
  11     |   1   |    1        << round change
  11     |   2   |    2
  12     |   0   |    0
  13     |   0   |    0
  14     |   0   |    0
  ```
  </details>
- `WeightedRandom` (2):
  - Since Randao HF, the proposer is randomly selected among the committee. That is, `Proposer(N, R) = Committee(N)[R]`.
      <details>
        <summary>Weightedrandom Randao Example</summary>

      ```
      proposers=[0,1,2,3]
          
      block  | round | proposer
      -------------------------
         1   |   0   |    3
         2   |   0   |    1
         3   |   0   |    2
         4   |   0   |    3
         5   |   0   |    3
         6   |   0   |    0
         7   |   0   |    1
         8   |   0   |    0
         9   |   0   |    2
        10   |   0   |    3
        11   |   0   |    1        << round change
        11   |   1   |    2        << round change
        11   |   2   |    2
        12   |   0   |    3
        13   |   0   |    0
        14   |   0   |    1
      ```
      </details>
  - Before Randao HF, the proposer is randomly selected where the probability is proportional to the gini-coeff-applied staking amount.
    That is, `Proposer(N, R) = proposers[(N + R) % len(proposers)]` where `proposers` is refreshed at every "reward.stakingupdateinterval".

      <details>
        <summary>Weightedrandom Default Example</summary>
  
      ```
      proposers=[0,1,0,0,2,0,3,0,0,2,2,3,3,0,2,3]
  
      block  | round | proposer
      -------------------------
         1   |   0   |    0
         2   |   0   |    1
         3   |   0   |    0
         4   |   0   |    0
         5   |   0   |    2
         6   |   0   |    0
         7   |   0   |    3
         8   |   0   |    0
         9   |   0   |    0
        10   |   0   |    2
        11   |   0   |    2        << round change
        11   |   1   |    3        << round change
        11   |   2   |    3
        12   |   0   |    2        << next proposer of block 11
        13   |   0   |    3
        14   |   0   |    3
      ```
      </details>

The Mainnet and Kairos uses `WeightedRandom` policy.

The selection of proposer policy is limited by consensus algorithm:
- clique: RoundRobin [default: RoundRobin]
- istanbul - RoundRobin, Sticky, WeightedRandom [default: WeightedRandom]

## Persistent Schema
The voting blocks and the council addressList is stored at miscDB.

### Voting Blks
- valSetVoteBlockNums: The block numbers whose header contains addvalidator/removevalidator vote data.

### ValSetSnapshot
>**Note:** Since the voting block is a concept added later, migration is required. 
> Thus, until migration finished, the voting blocks may not be available for the past blocks. 
> To support this, we can still generate the council list for the target block by iterating over the vote data from the nearest ValSetSnapshot.
- ReadValSetSnapshot - it reads the closest ValSetSnapshot of block N. 
- StoreValSetSnapshot - deprecated.

### CouncilAddressList: The validator list of the council.
- councilAddressList(n) - It founds the closest addvalidator/removevalidator voting block 'v' of block number 'n'. Then, it reads the council(v)

## In-memory Structures
###  Council
- Council structure of block N: It categorizes the council(N-1) into qualified and demoted for block N. It's not for display purpose, but it's for calculating committee or proposer of block N.
```go
type council struct {
  blockNumber uint64  // id of Council
  
  qualifiedValidators []common.Address // qualified is a subset of prev block's council list
  demotedValidators   []common.Address // demoted is a subset of prev block's council who are demoted as a member of committee
  
  councilAddressList []common.Address // total council node address list. the order is reserved.
}
```
- Validators: Sorting purpose
```go
type subsetCouncilSlice []common.Address
```

### ValsetContext
- blockResult: It's a result(state) of the block. 
  - valSet(N) = valSet(N-1).handleVote(Block(N-1).vote)
  - pSet(N) = pSet(N-1).handleGov(Block(N-1).gov) (only if N is an epoch)
  - stakingInfo(N) = Block(N-1).stake_change (only if N is a stakingUpdateInterval)
  - header(N) = Block(N).header
  - author(N) = Author(Block(N).header)
```go
type blockResult struct {
  staking         *staking.StakingInfo
  header          *types.Header
  author          common.Address
  pSet            gov.ParamSet
}
```

## Module lifecycle

### Init
- Dependencies:
  - ChainKv: Read/Write Voting blks and Council address lists. The keys of council address list is the voting blks.
  - Chain: Get Header and config information from headerChain.
  - HeaderGov: Get govParam from headerGov.
  - StakingInfo: Get block's staking info from stakingInfo.

### Start and stop

This module does not have any background threads.

## Block processing

### Consensus

This module does not have any consensus-related block processing logic.

### Execution
- PostInsertBlock: At the end of the block execution, the addvalidator/removevalidator votes are handled and they are remove from MyVotes. If succeed, the voteBlk and councilAddressList db will be updated.

### Rewind

This module does not have any consensus-related block processing logic.

## APIs
- GetCouncil(number *rpc.BlockNumber) ([]common.Address, error)
- GetCouncilSize(number *rpc.BlockNumber) (int, error)
- GetCommittee(number *rpc.BlockNumber) ([]common.Address, error)
- GetCommitteeSize(number *rpc.BlockNumber) (int, error)
- GetValidators(number *rpc.BlockNumber) ([]common.Address, error)
- GetValidatorsAtHash(hash common.Hash) ([]common.Address, error)
- GetDemotedValidators(number *rpc.BlockNumber) ([]common.Address, error)
- GetDemotedValidatorsAtHash(hash common.Hash) ([]common.Address, error)

## Getters
- GetCouncilAddressList(N): it returns the list of council after executing 'PostInsertBlock'.
- GetCommitteeAddressList(N, round): it calculates the list of committee at the view (block N, round R)
- GetProposer(N, round): it calculates the proposer at the view (block N, round R)

