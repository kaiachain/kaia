# kaiax/valset

This module is responsible for getting council and calculating committee or proposer.

## Concepts
*Council(N): a validator list which is composed of qualified and demoted for committee, and is used for mining block N.<br/>
*Committee(N, R): a subset of qualified validators for mining block N. <br/>
*Proposer(N, R): a proposer is ensured to be a committee member in all the times regardless of hardfork or policy. <br/>

| Components                | [Block 0]                          | [Block N]                                                                                                                         |
|---------------------------|------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------|
| Council(N)                | GenesisCouncil                     | Council(N-1).applyVote(header(N-1)) <br/>+ "Add Validator(i)", i is not in Council <br/>- "Remove Validator (j)", j is in Council |
| └ demoted validators(N)   | None                               | Council(N).demotedByQualification()                                                                                               | 
| └ qualified validators(N) | GenesisCouncil                     | Council(N) - demoted(N)                                                                                                           |
| └└ Committee(N, R)        | GenesisCouncil                     | qualified(N).sublist(R)                                                                                                           |
| └└└ Proposer(N,R)         | GenesisCouncil's<br/>first element | Committee(N).pickProposer(R)                                                                                                      |

Several dependencies are required to calculate the committee or proposer of block N:
- Header(N-1) // e.g. mixHash, blockHash
- Council(N)
- StakingInfo(N) // will read StateAt(num-1)
- EffectiveParamSet(N) // governance parameter value

## council: a set of registered CN
A `Council(N)` refers to the council determined after the block `N-1` is executed. 

```
  Council(N) = apply vote in header[N-1].Vote at Council(N-1)
```
* GenesisCouncil: restored from the extraData of the genesis block(block 0).
* Add, Remove votes: "istanbul.addvalidator" or "istanbul.removevalidator" vote on block.Header
  * For more information about header votes, see [kaiax/gov](../gov/headergov/README.md)

## qualified, demoted: divided by qualifications
Council members can be divided into two groups by qualifications: Qualified and demoted.
To reach consensus on block N, `Qualified(N)` is used as a superset of Committee(N).
```
  Demoted(N) = filter out by qualifications of Council(N)
  Qualified(N) = remaining after excluding the demoted(N) from Council(N)
```
the qualifications vary depending on the proposer policy:
* round-robin: all council members are qualified.
* sticky: all council members are qualified.
* weighted-random: Validators who staked more than `EffectiveParamSet(N).minimumstake` are qualified.
  * If no nodes meet this requirement, all nodes are considered qualified. 
  * In single mode, governingnode is qualified without any conditions.

## committee: a subset of qualified validators
A `Committee(N, R)` is randomly selected as a subset of `Qualified(N)` for mining the block `N`.
The subset size of committee size is determined by `EffectiveParamSet(N).committeesize`.
```
  Committee(N, R) = a subset of Qualified(N)
```

### Committee selection logic
- If Genesis block(N is 0), all GenesisCouncil are committee.
- If Qualified size <= committee size, all qualified are committee.
- If policy is weighted-random and randao hardfork is activated, derive SelectRandaoCommittee.
- If policy is round-robin, or policy is sticky, or policy is weighted-robin but randao is not activated, derive SelectRandomCommittee.

Selecting Committee logic do the shuffling at Qualified(N) and slices it:
- Calculate the shuffling seed by using previous block's hash or mixHash.
- Shuffle the qualified validators with the calculated seed.
- After shuffling, the council is "sliced" up to `EffectiveParamSet(N).committeesize`.

### SelectRandomCommittee
SelectRandomCommittee is a function that selects a committee by shuffling the Qualified using a seed derived from prevHash.

Next describes the SelectRandomCommittee logics by steps:
```
Condition: proposerIdx = 3, nextProposerIdx = 7, committeesize = 6, qualified validators = [0,1,2,3,4,5,6,7,8,9]

SimpleRandomCommittee:
- Step1. calculate the shuffling seed by prevHash: seed = int64(binary.BigEndian.Uint64(prevHash.Bytes()[:8]))
- Step2. secure the proposer and the next distinct proposer: proposers = [3,7], council = [0,1,2,4,5,6,8,9]
- Step3. shuffle the council: proposers = [3,7], council = [4,5,6,8,9,0,1,2]
- Step4. merge. council = [3,7,4,5,6,8,9,0,1,2]
- Step5. slice by committeesize: committee = [3,7,4,5,6,8]
```

### SelectRandaoCommittee
SelectRandaoCommittee is a function that selects a committee by shuffling the Qualified using a seed derived from prevMixHash. 

Next describes the SelectRandaoCommittee logic by steps:
```
Condition: proposerIdx = 3, nextProposerIdx = 7, committeesize = 6, qualified validators = [0,1,2,3,4,5,6,7,8,9]

RandaoCommittee:
- Step1. calculate the shuffling seed by prevMixHash: seed = int64(binary.BigEndian.Uint64(prevMixHash[:8]))
- Step2. shuffle the qualfied validators: council = [4,7,5,6,8,3,9,0,1,2]
- Step3. slice by committeesize: committee = [4,7,5,6,8,3]
```

## proposer: a committee member who proposes the block
A `Proposer(N, R)` proposes the block N at round R. The proposer is called an "author" once the block is created. The Mainnet and Kairos uses `WeightedRandom` policy.

### ProposerPolicy
Proposer is selected based on the governance parameter "istanbul.proposerpolicy". The mechanism of the same policy works differently on the hardfork.

Here are the list of policies: (the number in the paranthesis indicates the value of "istanbul.proposerpolicy")
- `RoundRobin` (0): every council member takes each turn.
- `Sticky` (1): a proposer is changed only when the round is changed.
- `WeightedRandom` (2) [default proposer policy]:
  - Since Randao HF, the proposer is randomly selected among the committee.
  - Before Randao HF, the proposer is randomly selected where the probability is proportional to the gini-coeff-applied staking amount. `proposers` is refreshed at every "reward.stakingupdateinterval".

A proposer is picked from different sources depending on the proposer policy and randao hardfork activation.
* RoundRobin or Sticky: a proposer is picked from Council(N)
* weighted-random and randao is activated: a proposer is picked from Committee(N)
* weighted-random and randao is not activated: a proposer is picked from proposers(pUpdateBlock)

### Proposer selection logic
- If Genesis Block(N is 0), a proposer is the first element of the GenesisCouncil.
- If policy is round-robin, proposer is selected as the next validator in the council list, based on the index (prevAuthorIdx + round).
- If policy is sticky, proposer remains the same as long as the round changes. If round changes, proposer is selected as the next validator in the council list.
- If policy is weighted-random and randao hardfork is not activated, proposer is selected from the Proposers(proposerUpdateBlock), based on the index ((N + round - proposerUpdateBlock) % len(proposers)).
  - proposerUpdateBlock = blockNum - (blockNum % `EffectiveParamSet(N).proposerUpdateInterval`)
- If policy is weighted-random and randao hardfork is activated, proposer is selected as the Round'th validator in the RandaoCommittee.

Next are the examples of the proposer selection logic per policies:
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
12     |   0   |    2
13     |   0   |    2
14     |   0   |    2
```
</details>
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

## Persistent Schema
### CouncilAddressList(num)
- Retrieve the record from miscdb with the key, "councilAddress" + uint64(v) 
- v is the latest block number from `num` which contains the addvalidator/removevalidator vote.
- Returns the sorted []common.Address

### ValidatorVoteDataBlockNums
- Retrieve the record with the key, "valSetVoteBlockNums"
- Returns []uint64

### IstanbulSnapshot(blockhash)
>**Note:** Since the voting block concept is added later, the migration of the valSetVoteBlockNums is required.
> Thus, until migration finished, the **voting blocks** may not be available for the past blocks.
> It would temporarily generate the council list for the target block by iterating over the vote data from the nearest istanbul snapshot.
- Retrieve the record with the key, "snapshot" + blockhash.
- Only available at intervals defined by istanbulCheckpoint.
- Returns a sorted []common.Address, which is a unified list of snap[blockhash].demoted + snap[blockhash].qualified.

### lowestScannedCheckpointIntervalKey
- Retrieve the record with the key, "lowestScannedCheckpointIntervalKey".
- Returns uint64.
- Identifies the voteBlock migration completion point; blocks below this number still requires migration.

## In-memory Structures
### Council
- Council structure of block N: It categorizes the council(N) into qualified and demoted for block N. It's not for display purpose, but it's for calculating committee or proposer of block N.
```go
type council struct {
  blockNumber uint64  // id of Council
  rules       params.Rules // rules of this council
  
  qualifiedValidators []common.Address
  demotedValidators   []common.Address
  
  councilAddressList []common.Address // total council node address list. the order is reserved.
}
```
### AddressList
Sorting purpose
```go
type AddressList []common.Address
```
### CommitteeContext: TBD

## Module lifecycle

### Init
- Dependencies:
  - ChainKv: Read/Write Voting blks and Council address lists. The keys of council address list is the voting blks.
  - Chain: Get Header and config information from headerChain.
  - HeaderGov: Get govParam from headerGov.
  - StakingInfo: Get block's staking info from stakingInfo.

### Start and stop
When starting the valSet module, the following checks and actions are performed:
* valSet initialization at genesis block
  * if the validator vote data (voteBlks) is missing, it stores []uint64{0} at voteBlks and genesis council to councilAddressListDB
  * if the current block number is genesis block, the lowest scanned interval(lowestScannedNum) is initialized as 0.
* lowestScannedNum to figure out if migration from istanbul snapshot to valSet DB is needed or not
  * if no lowestScannedNum exists, it replays the votes between [lastScannedNum, currentBlock)
  * updates the lowestScannedNum as lastScannedNum
* check lowestScannedNum is 0. If not, it still needs migration. Start the migration thread in background.

This module does not have stop logic.
## Block processing

### Consensus

This module does not have any consensus-related block processing logic.

### Execution
- PostInsertBlock: the addvalidator/removevalidator votes are handled. If success, the voteBlk and councilAddressListDB will be updated.

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
- GetCouncilAddressList(N): Council(N) is either retrieved from councilAddressListDB or istanbul snapshot DB.
  - when migration is finished at N-1, finalized block, it retrieves from councilAddressListDB.
  - when migration still in progress, it retrieves from istanbul snapshot DB.
     - however, snapshot is stored every checkpointinterval(=1024)
     - so, to get snap[N-1], need to iterate over the non-snapshot interval blocks
- GetCommittee(N, round): it calculates the list of committee at (block N, round R)
- GetProposer(N, round): it calculates the proposer at (block N, round R)

