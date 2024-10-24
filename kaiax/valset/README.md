# kaiax/valset

This module is responsible for getting council and calculating committee or proposer.

## Concepts
| Components          | [Block 0]                          | [Block N]<br/>"AddValidator" <br/>"RemoveValidator"                                          | [Block N+1]                             |
|---------------------|------------------------------------|----------------------------------------------------------------------------------------------|-----------------------------------------|
| Council(N)          | GenesisCouncil                     | GenesisCouncil <br/>+ validator(i), i is not in Council <br/>- validator(j), j is in Council | ---                                     |
| ㄴ qualified(N)      | GenesisCouncil                     | Council(N-1).filter()                                                                        | Council(N).filter()                     |
| ㄴ demoted(N)        | None                               | Council(N-1)<br/>- Council(N-1).qualified()                                                  | Council(N)<br/>- Council(N).qualified() | 
| Committee(N,Round)  | GenesisCouncil                     | qualified(N).sublist(Round)                                                                  | qualified(N).sublist(R)                 |
| ㄴ proposer(N,Round) | GenesisCouncil's<br/>first element | Committee(N,R).proposer(Round)                                                               | Committee(N,R).proposer(Round)          | 

### council: a list of registered CN
A member of council is added/removed by "governance.addvalidator" vote. The genesis council is restored via genesis extraData.
The council(N) is decided after the block(N).vote and govParam(N) is applied, and it will be used to calculate next committee or proposer.

### committee: a subset of council who participates on consensus
A committee(N, R) is calculated based on previous block's results(council, prevHeader, stakingInfo of council) and sorted. However, the genesis committee is copied from the genesis council.
- minimum staking amount - The council members who have less than minimum staking amount is demoted, so it cannot be a member of committee. 
- committeesize - The committee size can be updated via "istanbul.committeeSize". It decides the size of committee.
- committee shuffle seed - The seed is calculated using previous block's information. The copied council is shuffled with the calculated seed to get the committee.

Committee selection logic is different before/after Randao Hardfork when it's proposer policy is weightedCouncil. So the condition to activate RandaoCommittee is `policy.IsDefaultSet() || (policy.IsWeightedCouncil() && !rules.IsRandao)`
- committee shuffle seed calculation logic
  - before Randao: the seed is calculated using prevHash. `seed = int64(binary.BigEndian.Uint64(prevHash.Bytes()[:8]))`. 
  - after Randao: the seed is calculated using mixHash. `seed = int64(binary.BigEndian.Uint64(mixHash[:8]))`
- council shuffle
  - before Randao: extract (proposer, next proposer which is differnt from proposer) and shuffle it. Attach the proposers again and slice the council.
  - after Randao: shuffle the council and slice the committee.

Example of BeforeRandaoCommittee
```
Condition: proposerIdx = 3, nextProposerIdx = 7, committeesize = 6, council = [0,1,2,3,4,5,6,7,8,9]
Step1. extract proposers to committee: proposers = [3,7], council = [0,1,2,4,5,6,8,9]
Step2. shuffle the council. proposers = [3,7], council = [4,5,6,8,9,0,1,2]
Step3. merge. council = [3,7,4,5,6,8,9,0,1,2]
Step4. slice the council by committee size. committee = [3,7,4,5,6,8]
```
### proposer: a member of committee who proposes the block
A proposer proposes the block. We call as author after the block is created.
#### ProposerPolicy
Proposer is selected based on the proposer policy the network chosen. Also, each proposer selection logic has been updated per HF.
- RoundRobin: `council(N-1)[(prevAuthorIdx+round)%len(council(N-1))]`
- Sticky: `council(N-1)[(prevAuthorIdx+round+1)%len(council(N-1))]`
- WeightedRandom
  - the council is splitted into qualified, demoted addresses by minimum staking amount since istanbul HF.
  - proposer is picked from proposers array which is created at proposerupdateinterval block with staking amount and gini. However, it is deprecated since Randao HF.

The selection of proposer policy is limited by consensus algorithm.
  - clique - 0: RoundRobin [default: 0]
  - istanbul - 0: RoundRobin, 1: Sticky, 2: WeightedRandom [default: 0]

## Persistent Schema
The voting blks and the council addressList is stored at miscDB
### Voting Blks
- ReadVoteBlks - it reads whole addvalidator/removevalidator voting blks 
- StoreVoteBlks - it stores/updates whole addvalidator/removevalidator voting blks

### ValSetSnapshot
- ReadValSetSnapshot - it reads the closest ValSetSnapshot of block N. 
- StoreValSetSnapshot - deprecated. kaiax no longer stores ValSetSnapshot.

### Valset
- ReadCouncilAddressListFromDb(n) - n is the addvalidator/removevalidator voting blks
- WriteCouncilAddressListToDb(n, council) - n is the addvalidator/removevalidator voting blks

## In-memory Structures
###  Council
- Council: it categorizes the previous council list for block N. it's not for display purpose, but it's for calculating committee or proposer of block N.
```go
type council struct {
  blockNumber uint64  // id of Council
  
  qualifiedValidators subsetCouncilSlice // qualified is a subset of prev block's council list
  demotedValidators   subsetCouncilSlice // demoted is a subset of prev block's council who are demoted as a member of committee
  
  councilAddressList subsetCouncilSlice // total council node address list. the order is reserved.
}
```
- Validators: sorting purpose
```go
type subsetCouncilSlice []common.Address
```

### ValsetContext
- blockResult: it's a result(state) of previous block N. the committee/proposer of next block
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
- PostInsertBlock: At the end of the block execution, the addvalidator or removevalidator votes are handled and they are remove from MyVotes. If succeed, the voteBlk and councilAddressList db will be updated.

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

