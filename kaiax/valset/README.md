# kaiax/valset

This module is responsible for getting council and calculating committee or proposer.

## Concepts
| Components               | [Block 0]                          | [Block N-1] <br/>"AddValidator" <br/>"RemoveValidator"                                      | [Block N]                                        |
|--------------------------|------------------------------------|---------------------------------------------------------------------------------------------|--------------------------------------------------|
| demoted(Block)           | None                               | Council(N-2).demotedByQualification(state(N-2))                                             | Council(N-1).demotedByQualification(*state(N-1)) | 
| qualified(Block)         | GenesisCouncil                     | Council(N-2) - demoted(N-1)                                                                 | Council(N-1) - demoted(N)                        |
| └ Committee(Block,Round) | GenesisCouncil                     | qualified(N-1).sublist(Round)                                                               | qualified(N).sublist(Round')                     |
| └└ proposer(Block,Round) | GenesisCouncil's<br/>first element | Committee(N-1,Round).proposer(Round)                                                        | *Committee(N,Round').proposer(Round')            |
| Council(Block)           | GenesisCouncil                     | Council(N-2) <br/>+ validator(i), i is not in Council <br/>- validator(j), j is in Council  | ...                                              |
*Committee(N,Round').proposer(Round'): Actually, before Randao, proposer is not picked from Committee. Nevertheless, proposer is a member of the Committee, so it is representeas as belonging to the Committee.<br/>
*state(N-1): It's a result of processing block N-1 including state transition. The results are govParamSet, stakingInfo, and author.

### council: a set of registered CN
The genesis council is restored via genesis block(block 0)'s extraData. 
Then, the member of council is added/removed by "governance.addvalidator"/"governance.removevalidator" vote.
The council of Block N is finalized after the block(N).vote is applied, and it will be used to calculate the next committee or proposer.

The council is classified into two groups: qualified validators and demoted validators. 
* Qualified validators: the validators who are qualified to be a committee member. 
* Demoted validators: the remainder after removing the qualified validators.

The qualifications differ depending on the proposer policy.
* clique consensus: it implements default proposer policy (round-robin). 
  * round-robin: all council members are qualified
* istanbul consensus: it implements three proposer policies(round-robin, sticky, weighted-random).
  * round-robin: all council members are qualified
  * sticky: all council members are qualified
  * istanbul: Validators that meet the minimum staking requirement are qualified. 
    * If no node meet this requirement, all nodes are deemed qualified. In single mode, governingnode is qualified without any conditions.

### committee: a subset of qualified council members
A committee of Block N, Round R is calculated based on previous block's council and results(header.author, stakingInfo, govParamSet). However, the committee of genesis block doesn't have previous block, so it is copied from the genesis council.

- committeesize - The committee size can be updated via "istanbul.committeeSize". It decides the size of committee.
- committee shuffle seed - The seed is calculated using previous block's information. The copied qualified council is shuffled with the calculated seed to get the committee.

Committee selection logic is different before/after Randao Hardfork when it's proposer policy is weightedCouncil. So the condition to activate RandaoCommittee is `policy.IsDefaultSet() || (policy.IsWeightedCouncil() && !rules.IsRandao)`
- committee shuffle seed calculation logic
  - before Randao: the seed is calculated using prevHash. `seed = int64(binary.BigEndian.Uint64(prevHash.Bytes()[:8]))`. 
  - after Randao: the seed is calculated using mixHash. `seed = int64(binary.BigEndian.Uint64(mixHash[:8]))`
- qualified council shuffle
  - before Randao: extract (proposer, next proposer which is differnt from proposer) and shuffle it. Attach the proposers again and slice the it.
  - after Randao: shuffle the qualified council and slice the committee.

Example of BeforeRandaoCommittee
```
Condition: proposerIdx = 3, nextProposerIdx = 7, committeesize = 6, council = [0,1,2,3,4,5,6,7,8,9]
Step1. extract proposers to committee: proposers = [3,7], council = [0,1,2,4,5,6,8,9]
Step2. shuffle the council. proposers = [3,7], council = [4,5,6,8,9,0,1,2]
Step3. merge. council = [3,7,4,5,6,8,9,0,1,2]
Step4. slice the council by committee size. committee = [3,7,4,5,6,8]
```
### proposer: a member of committee who proposes the block
A proposer means the member of committee who proposes the block. We call as author after the block is created. Proposer is selected based on the proposer policy the network chosen. Also, each proposer selection logic has been updated per HF.
#### ProposerPolicy
- RoundRobin: `council(N-1)[(prevAuthorIdx+round)%len(council(N-1))]`
- Sticky: `council(N-1)[(prevAuthorIdx+round+1)%len(council(N-1))]`
- WeightedRandom
  - Before Randao: proposer of (block N, Round R) is picked from "proposers" array which is created at proposerupdateinterval block with staking amount and gini. 
  - After Randao: proposer of (block N, Round R) is the n'th element of the committee(N)

The selection of proposer policy is limited by consensus algorithm.
  - clique - 0: RoundRobin [default: 0]
  - istanbul - 0: RoundRobin, 1: Sticky, 2: WeightedRandom [default: 0]

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
  
  qualifiedValidators subsetCouncilSlice // qualified is a subset of prev block's council list
  demotedValidators   subsetCouncilSlice // demoted is a subset of prev block's council who are demoted as a member of committee
  
  councilAddressList subsetCouncilSlice // total council node address list. the order is reserved.
}
```
- Validators: Sorting purpose
```go
type subsetCouncilSlice []common.Address
```

### ValsetContext
- blockResult: It's a result(state) of previous block N. the committee/proposer of next block
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

