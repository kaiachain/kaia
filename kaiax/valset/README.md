# kaiax/valset

This module is responsible for tracking the set of block validators.

## Concepts

### Overview

The validators are subdivided as below subsets.

- `Council(N)`: A set of validator addresses where the proposers and committee for the block N are selected from.
  - `DemotedValidators(N)`: Subset of the council not eligible to be committee and proposer.
  - `QualifiedValidators(N)`: Subset of the council eligible to be committee and proposer.
    - `Committee(N,R)`: Subset of qualified validators that can validate the block N at round R.
    - `Proposer(N,R)`: One of the qualified validators that will finalize and propose the block N at round R. Note that Proposer(N,R) is must be included in Committee(N,R).
  - The demoted and qualified validators are mutually exclusive and collecively form the council.
  - The qualified validators are often simply called "validators" as in `IstanbulExtra.Validators`.

The rules determining the subsets depend on the proposer policy. Proposer policy is an integer governance parameter (`istanbul.policy`) that is first defined in the ChainConfig (`Istanbul.ProposerPolicy`) at genesis and never changes afterwards.

- RoundRobin (0): The simplest rule where proposers are selected in round-robin fashion.
- Sticky (1): Another simple rule where the proposer is changed only upon round change.
- WeightedRandom (2): The full-featured rule utilizing staking and randomness.
  - Before Kore hardfork, proposers are selected at random with their staking amounts (see [kaiax/staking](../staking/README.md)) affecting the probability.
  - After Kore hardfork, proposers are selected at random with uniform probabiliy.
  - After Randao hardfork, proposers are selected at random with uniform probabiliy. In addition, the selection becomes unpredictable.

At the genesis block (N=0), it is meaningless to define DemotedValidators, QualifiedValidators, Committee or Proposer because genesis block is not subject to consensus. But we define them nonetheless for completeness.

### Council

Council is the largest set that all other subsets are drawn from. Council is the set of registered validators.

- Council is initially defined at genesis block's `extraData` field in the `IstanbulExtra` format.
- Validators can enter to or exit from the council by the header governance votes (read more at [kaiax/gov](../gov/headergov/README.md))
  - "governance.addvalidator" vote contains one or more addresses to be added to the council of the next block and onwards.
  - "governance.removevalidator" vote contains one or more addresses to be removed from the concil of the next block and onwards.
  - Because council is a set of unique addresses, already-existing address cannot be added and non-existing address cannot be removed.
- The votes are incremental. i.e.
  ```
  Council(0) = parse(Genesis.ExtraData)
  Council(N) = Council(N-1) + AddValidatorVotes(N-1) - RemoveValidatorVotes(N-1)
  ```

### Qualified and Demoted validators

Each proposer policy require different qualifications to council members.

- RoundRobin: All council members are qualified.
- Sticky: All council members are qualified.
- WeightedRandom: Validators are qualified if it stakes no less than the minimum staking amount.
  - For genesis block, all council members are qualified.
  - For blocks before Istanbul hardfork, all council members are qualified.
  - If no validator meets the minimum staking requirement, all council members are qualified.
  - If the governance mode is "single", the governing node is unconditionally qualified to allow governance parameter change under any circumstances.
  - The minimum staking amount refers to the `reward.minstake` parameter, the governance moce is the `governance.governancemode` parameter and governing node is the `governance.governingnode` parameter.
- In the implementation, the demoted validators are first picked from the council. Then qualified validators are calculated as the rest. i.e.
  ```
  DemotedValidators(N) = (various rules)
  QualifiedValidators(N) = Council(N) - DemotedValidators(N)
  ```

### Use of PRNG

Some algorithms in this module uses PRNG to generate deterministic yet statistically well-distributed results.

Those algorithms calculates a deterministic random seed that depends on the block number N and the round number R. The random seed is fed into the golang's PRNG [`math/rand`](https://pkg.go.dev/math/rand) that implements a linear congruential generator with a=48271 and c=0 ([wiki](https://en.wikipedia.org/wiki/Lehmer_random_number_generator), [source](https://cs.opensource.google/go/go/+/refs/tags/go1.23.4:src/math/rand/rng.go)).

The PRNG is then used to shuffle a lexicographically sorted array of addresses. There are two shuffling methods in use.
- Legacy-shuffle: use the following algorithm that depends on a stream of integers modulo array length.
    ```go
    for i = 0; i < len(a); i++ {
      j = prng.Intn(len(a)) // random() % len(a)
      a[i], a[j] = a[j], a[i]
    }
    ```
- Builtin-shuffle: use the PRNG's `Shuffle()` function that implements the Fisher–Yates shuffle ([wiki](https://en.wikipedia.org/wiki/Fisher%E2%80%93Yates_shuffle), [source](https://cs.opensource.google/go/go/+/refs/tags/go1.23.4:src/math/rand/rand.go)).
    ```go
    prng.Shuffle(len(a), func(i, j int) {
      a[i], a[j] = a[j], a[i]
    })
    ```

### Committee

A Committee is a randomly sampled subset of the qualified validators. To make a block valid, a quorum number of committee members must sign the block hash.

- The committee of the genesis block is all qualified validators, and by extension all council members.
- If qualified validator size is less than or equal to the committee size, the committee is all qualified validators.
    - The committee size refers to the `istanbul.committeesize` paramter.
- If committee size is 1, return an 1-element set `{Proposer(N,R)}`.
    - From here, we can assume `1 < committee size < qualified validator size`.
- Otherwise, randomly sample from the qualified validators. The sampling algorithm depends on the proposer policy.
  - RoundRobin: SelectRandomCommittee
  - Sticky: SelectRandomCommittee
  - WeightedRandom before Kore hardfork: SelectRandomCommittee
  - WeightedRandom after Kore hardfork: SelectRandomCommittee
  - WeightedRandom after Randao hardfork: SelectRandaoCommittee

<details>
<summary><b>SelectRandomCommittee</b></summary>

The "SelectRandomCommittee" algorithm uses previous block hash as the PRNG seed. The resulting committee must include the current round's proposer and the next distinct proposer. Next distinct proposer is first of the following round proposers that is not equal to the current proposer.

1. Find the next distinct proposer.
    - Note that at this point, there are 2 or more qualified validators.
    - RoundRobin: Next round's proposer. It must be different from the current proposer.
    - Sticky: Next round's proposer. It must be different from the current proposer.
    - WeightedRandom before Kore hardfork: Traverse through the proposers list (i.e. `Proposer(N,R+x), x=1,2,..`) until finding a distinct proposer.
    - WeightedRandom after Kore hardfork: Traverse through the proposers list (i.e. `Proposer(N,R+x), x=1,2,..`) until finding a distinct proposer.
    - WeightedRandom after Randao hardfork: Not applicable.
1. Current proposer and next distinct proposer are temporarily removed from the qualified validators.
    ```go
    // Suppose QualifiedValidators has K elements (K >= 2)
    ShuffledValidators = QualifiedValidators - CurrentProposer - NextDistinctProposer
    // ShuffledValidators has K-2 elements
    ```
1. PRNG seed is the first 15 hex digits (nibbles) of previous block hash (prevHash)
    ```go
    hexstring = strings.TrimPrefix(prevHash.Hex(), "0x")[:15]
    seed = strconv.ParseInt(hashstring, 16, 64)
    ```
1. Lexicographically sort and legacy-shuffle ShuffledValidators.
1. Return current proposer, next distinct proposer and the first (CommitteeSize - 2) elements from the shuffled array.
    ```
    {CurrentProposer, NextDistinctProposer, ShuffledValidators[:CommitteeSize-2]}
    ```

For example,
```
Given:
- PrevHash = 0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84
- Sorted QualifiedValidators = [0 1 2 3 4 5 6 7 8 9]
- CommitteeSize = 6
- CurrentProposer = 3
- NextDistinctProposer = 7

Calculate:
- Seed = 0x112233445566778 = 77162851027281784
- ShuffledValidators = [0 1 2 4 5 6 8 9]
- ShuffledValidators = [6 8 1 2 4 0 9 5]
- Committee = [3 7 6 8 1 2]
```

</details>

<details>
<summary><b>SelectRandaoCommittee</b></summary>

The "SelectRandaoCommittee" algorithm uses the previous block RANDAO mixHash (prevMixHash) as the PRNG seed. See [KIP-146](https://github.com/kaiachain/kips/blob/main/KIPs/kip-146.md) for more details.

1. PRNG seed is the first 8 bytes of the prevMixHash
    ```
    seed = int64(Uint64BE(prevMixHash))
    ```
1. Lexicographically sort and builtin-shuffle the QualifiedValidators
1. Return the first CommitteeSize elements
    ```
    QualifiedValidators[:CommitteeSize]
    ```

For example,
```
Given:
- PrevHash = 0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84
- Sorted QualifiedValidators = [0 1 2 3 4 5 6 7 8 9]
- CommitteeSize = 6

Calculate:
- Seed = 0x1122334455667788 = 1234605616436508552
- QualifiedValidators = [8 3 5 1 0 9 2 6 7 4]
- Committee = [8 3 5 1 0 9]
```

</details>

### Proposer

A proposer is a member of the committee that decides the transaction execution order of the block and broadcasts (to "propose") to other validators for approval. The proposer of a block is also called the "author" of the block.

- Proposer of the genesis block is "0x0000000000000000000000000000000000000000".
- Proposer of other blocks depend on the proposer policy, previous block's proposer, current block number and round number.
  - RoundRobin: SelectRoundRobinProposer
  - Sticky: SelectStickyProposer
  - WeightedRandom before Kore hardfork: SelectWeightedRandomProposer
  - WeightedRandom after Kore hardfork: SelectUniformRandomProposer
  - WeightedRandom after Randao hardfork: SelectRandaoProposer

<details>
<summary><b>SelectRoundRobinProposer</b></summary>

In "SelectRoundRobinProposer" algorithm, each validator takes turn.

1. Lexicographically sort the QualifiedValidators array.
2. Let `prevIdx` be the index of the previous block's proposer. If the previous proposer is not in the array (e.g. genesis block), let `prevIdx=0`.
3. Return `[(prevIdx + round + 1) % len]`-th element of the array.

For example,
```
Given:
- QualifiedValidators = [0 1 2 3]

Calculate:
block     0  1  2  3  4  5  6  6  6  7  8
round     0  0  0  0  0  0  0  1  2  0  0
                               RC RC
proposer  x  0  1  2  3  0  1  2  3  0  1

*RC = round change
*At block=6, round=0 => prevIdx=0 => proposerIdx=(0+0+1)%4=1
*At block=6, round=1 => prevIdx=0 => proposerIdx=(0+1+1)%4=2
*At block=6, round=2 => prevIdx=0 => proposerIdx=(0+2+1)%4=3
*At block=7, round=0 => prevIdx=3 => proposerIdx=(3+0+1)%4=0
```

</details>
<details>
<summary><b>SelectStickyProposer</b></summary>

In "SelectStickyProposer" algorithm, proposer remains the same unless a round change occurs.

1. Lexicographically sort the QualifiedValidators array.
2. Let `prevIdx` be the index of the previous block's proposer. If the previous proposer is not in the array (e.g. genesis block), let `prevIdx=0`.
3. Return `[(prevIdx + round) % len]`-th element of the array.


For example,
```
Given:
- QualifiedValidators = [0 1 2 3]

Calculate:
block     0  1  2  3  4  5  6  6  6  7  8
round     0  0  0  0  0  0  0  1  2  0  0
                               RC RC
proposer  x  0  0  0  0  0  0  1  2  2  2

*RC = round change
*At block=6, round=0 => prevIdx=0 => proposerIdx=(0+0)%4=0
*At block=6, round=1 => prevIdx=0 => proposerIdx=(0+1)%4=1
*At block=6, round=2 => prevIdx=0 => proposerIdx=(0+2)%4=2
*At block=7, round=0 => prevIdx=3 => proposerIdx=(2+0)%4=2
```

</details>
<details>
<summary><b>SelectWeightedRandomProposer</b></summary>

In "SelectWeightedRandomProposer", a proposer list is pre-generated every proposer interval block. A single validator can appear multiple times in the list. Higher the staking amount more frequently the validator appears in the list. The list is randomized using a deterministic PRNG.

**proposer list generation**

Given "proposer update block" in which all the information below (qualified validators, staking info, block hash) is based on,
1. Adjust the staking amounts of the qualified validators.
    - $S_i = \text{Round}({S_i}^{1/G})$ if Gini option is enabled.
    - $S_i = S_i$ (unchanged) if Gini option is disabled.
    - Gini option is controlled by the `reward.useginicoeff` parameter.
    - G is the Gini coefficient of the staking amounts that are no less than the minimum staking amount (i.e. `G = stakingInfo.Gini(minStake)`).
1. Calculate the total staking amount of the qualified validators. Use the Gini-adjusted amounts, if enabled.
    - $TS = \sum{S_i}$
1. Calculate the percentile weight (integers 1 to 100) of the qualified validators. Every validator has a nonzero chance.
    - $W_i = \text{max}\left(1, \text{round}\left(\frac{100S_i}{TS}\right)\right)$ if TS > 0
    - $W_i = 0$ if TS = 0
1. Construct the "proposer list" where each validator appears $W_i$ (could be 0) times. If the proposer list is empty, put each validator once.
    - $PL = \text{repeat}(0, W_0) || ... || \text{repeat}(k, W_k)$ if $\sum{W_i} > 0$
    - $PL = 0, 1, ..., k$ if $\sum{W_i} = 0$
    - where $\text{repeat}(v,n) = v, v, ..., v\ \ (n\text{ times})$
1. PRNG seed is the first 15 hex digits (nibbles) of proposer update block's hash (updateHash)
    ```go
    hexstring = strings.TrimPrefix(updateHash.Hex(), "0x")[:15]
    seed = strconv.ParseInt(hashstring, 16, 64)
    ```
1. Lexicographically sort and legacy-shuffle the proposer list.

For example,
```
Given:
- UpdateHash = 0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84
- Sorted QualifiedValidators = [0 1 2 3]
- StakingAmounts = [5m 10m 15m 20m]
- ProposerInterval = 10
- MinStake = 5m
- UseGini = false

Calculate:
- Adjusted staking amounts S_i = [5m 10m 15m 20m]
- Total staking amount TS = 50m
- Weights W_i = [10 20 30 40]
- Seed = 0x112233445566778 = 77162851027281784
- Proposer List PL = [0 0 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 2 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3 3]
- Proposer List PL = [1 1 3 2 0 3 2 3 1 1 3 1 3 2 3 1 2 2 0 3 3 2 3 3 2 1 1 1 3 3 2 3 1 2 3 1 2 3 2 2 0 3 2 2 1 3 1 0 2 2 2 1 1 3 2 3 1 2 3 3 0 3 3 3 2 3 3 3 2 3 3 3 0 3 2 0 1 3 2 2 1 2 3 1 3 3 0 0 3 1 2 3 3 2 2 3 2 0 3 2]
```

**proposer retrieval**

Given block number N and round R,

1. Calculate the "proposer update block number" U.
    ```
    U = RoundDown(N-1, ProposerInterval)
    RoundDown(n, p) = n - (n % p)
    ```
    - ProposerInterval refers to the `reward.proposerupdateinterval` that is first defined in the ChainConfig (`Reward.ProposerUpdateInterval`) at genesis and never changes afterwards.
1. Generate the proposer list `ProposerList(U)` from the information at block S.
    - Because multiple block numbers (N) maps to the same proposer update block number (U), it is worth caching the list.
1. Return the `[(N + R - U - 1) % len]`-th element of the proposer list.

For example,

```
Given:
- ProposerInterval = 10
- ProposerList(0) = Council(0) = [0 1 2 3]
- ProposerList(10) = [1 1 3 2 0 3 2 3 1 1 3 1 3 2 3 1 2 2 0 3 3 2 3 3 2 1 1 1 3 3 2 3 1 2 3 1 2 3 2 2 0 3 2 2 1 3 1 0 2 2 2 1 1 3 2 3 1 2 3 3 0 3 3 3 2 3 3 3 2 3 3 3 0 3 2 0 1 3 2 2 1 2 3 1 3 3 0 0 3 1 2 3 3 2 2 3 2 0 3 2]
- ProposerList(20) = [2 1 0 3 1 2 0 3 3 2 3 3 3 3 1 0 2 3 0 2 1 1 3 3 1 1 3 1 1 2 0 2 2 2 3 3 2 2 2 0 3 0 2 3 1 3 1 3 3 3 3 3 2 3 2 3 2 3 3 3 2 2 2 3 1 3 3 2 2 2 3 2 1 2 3 3 3 3 3 2 0 2 3 2 2 1 0 1 1 1 3 3 3 1 2 1 0 2 1 3]

Calculate:
block     0  1  2  3  4  5  6  6  6  7  8  9 10 11 12 13 14 15 16 17 18 19 20 21 22
round     0  0  0  0  0  0  0  1  2  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0  0
                               RC RC
updateNum x  0  0  0  0  0  0  0  0  0  0  0  0 10 10 10 10 10 10 10 10 10 10 20 20
proposer  x  0  1  2  3  0  1  2  3  2  3  0  1  1  1  3  2  0  3  2  3  1  1  2  1

*RC = round change
*At block=6, round=0 => idx=(6+0-0-1)%4=1 => proposer=1
*At block=6, round=1 => idx=(6+1-0-1)%4=2 => proposer=2
*At block=6, round=2 => idx=(6+2-0-1)%4=3 => proposer=3
*At block=7, round=0 => idx=(7+0-0-1)%4=2 => proposer=2
*At block=11, round=0 => idx=(11+0-10-1)%100=0 => proposer=1
```

</details>
<details>
<summary><b>SelectUniformRandomProposer</b></summary>

The "SelectUniformRandomProposer" is same as the "SelectWeightedRandomProposer" algorithm except that the weights $W_i$ are always zero.

**proposer list generation**

Given "proposer update block" in which all the information below (qualified validators, staking info, block hash) is based on,
1. Construct the "proposer list" where each validator appears once.
    - $PL = 0, 1, ..., k$
1. PRNG seed is the first 15 hex digits (nibbles) of proposer update block's hash (updateHash)
    ```go
    hexstring = strings.TrimPrefix(updateHash.Hex(), "0x")[:15]
    seed = strconv.ParseInt(hashstring, 16, 64)
    ```
1. Lexicographically sort and legacy-shuffle the proposer list.

For example,
```
Given:
- UpdateHash = 0x1122334455667788af897911c946935ca28f37cb3b1bf9a30f17c84084276a84
- Sorted QualifiedValidators = [0 1 2 3]

Calculate:
- Seed = 0x112233445566778 = 77162851027281784
- Proposer List PL = [0 1 2 3]
- Proposer List PL = [1 3 0 2]
```

**proposer retrieval**

Identical to SelectWeightedRandomProposer.

</details>
<details>
<summary><b>SelectRandaoProposer</b></summary>

THe "SelectRandaoProposer" selects the proposer from the committee calculated by SelectRandaoCommittee. See [KIP-146](https://github.com/kaiachain/kips/blob/main/KIPs/kip-146.md) for more details.

1. Calculate the committee.
2. Return `[round % len]`-th element of the committee.

For example,

```
Given:
- QualifiedValidators = [0 1 2 3]
- Committee(0) = [0 1 2 3]
- Committee(1) = [1 3 0 2]
- Committee(2) = [3 0 2 1]
- Committee(3) = [2 1 3 0]
- Committee(4) = [0 2 1 3]
- Committee(5) = [3 0 2 1]
- Committee(6) = [2 3 1 0]
- Committee(7) = [2 1 0 3]
- Committee(8) = [3 2 0 1]

Calculate:
block     0  1  2  3  4  5  6  6  6  7  8
round     0  0  0  0  0  0  0  1  2  0  0
                               RC RC
proposer  x  1  3  2  0  3  2  3  1  2  3

*RC = round change
*At block=6, round=0 => idx=0 => proposer=2
*At block=6, round=1 => idx=1 => proposer=3
*At block=6, round=2 => idx=2 => proposer=1
*At block=7, round=0 => idx=0 => proposer=2
```

</details>

## Persistent Schema

- `validatorVoteBlockNums`: The block numbers whose header contains the validator vote data. It must at least include the number 0.
  ```
  "validatorVoteBlockNums" => JSON.Marshal([num1, num2, ...])
  ```
- `council`: The council at the given block number. Stored at the block numbers whose header contains the validator vote data. The addresses are sorted in lexicographic order. At least one entry at number 0 for genesis council must be present.
  ```
  "council" || Uint64BE(num) => JSON.Marshal([addr1, addr2, ...])
  ```
- `lowestCheckpointScannedBlockNum`: The lowest block number whose vote data and council is calculated from the legacy istanbul snapshot schema. That is, only vote block numbers are greater than or equal to this value are stored in `validatorVoteBlockNums`. It grows downwards by `istanbulCheckpointInterval` blocks.
  ```
  "lowestCheckpointScannedBlockNum" => Uint64BE(num)
  ```
- `istanbulSnapshot`: The legacy schema that periodically (every `istanbulCheckpointInterval` block) commits the council and other fields. Council at an arbitrary block can be recovered by accumulating the validator votes from the nearest istanbul snapshot.
  ```
  "snapshot" || blockHash[:] => JSON.Marshal(IstanbulSnapshot)
  ```

## In-memory Structures

### sortableAddressList

```go
type sortableAddressList []common.Address
```

Utility type for lexicographically sorting the addresses. Note that the sort order is based on EIP-155 mixed-case string format, not the byte array format.

### AddressSet

```go
type AddressSet struct {
  list sortableAddressList
  mu sync.RWMutex
}
```

An ordered set of addresses. Used to represent a set of validators such as council and committee.

### CommitteeContext

A committeeContext is a context for calculating the committee or proposer of block N. The pUpdateBlock and proposers is no longer used since Randao Hardfork. It provides `getCommittee` and `getProposer` methods.

```go
type committeeContext struct {
	qualified valset.AddressList
	num       uint64
	rules     params.Rules

	// pSet
	committeeSize          uint64
	proposerPolicy         ProposerPolicy
	proposerUpdateInterval uint64

	// num-1
	prevHeader *types.Header
	prevAuthor common.Address

	// num - (num % proposerUpdateInterval)
	pUpdateBlock uint64
	proposers    []common.Address
}
```

## Module lifecycle

### Init

- Dependencies:
  - ChainKv: Raw key-value database to access this module's persistent schema.
  - ChainConfig: Holds the genesis parameters such as ProposerUpdateInterval and ProposerPolicy.
  - Chain: Provides block headers.
  - GovModule: Provides governance parameters.
  - StakingModule: Provides staking info.

### Start and stop

This module maintains a background thread that migrates `istanbulCheckpoint` schema into the new `valsetVoteBlockNums` and `council` schemas.

## Block processing

### Consensus

This module does not have any consensus-related block processing logic.

### Execution

This module updates persistent schema based on `header.Vote`.

### Rewind

Upon rewind, this module deletes the related persistent data and flushes in-memory caches.

## APIs

This module does not expose APIs.

## Getters

- `GetCouncil(num)`: Returns the council at the block `num`.
  ```
  GetCouncil(num) -> []common.Address
  ```
- `GetValidators(num)`: Returns the demoted validators at block `num`. Note that you can calculate the qualified validators as Council.Subtract(DemotedValidators).
  ```
  GetDemotedValidators(num) -> []common.Address
  ```
- `GetCommittee(num, round)`: Returns the committee at the block `num` and round `round`.
  ```
  GetCommittee(num, round) -> []common.Address
  ```
- `GetProposer(num, round)`: Returns the proposer at the block `num` and round `round`.
  ```
  GetProposer(num, round) -> common.Address
  ```

### Implementation outline

```
# getters
GetCouncil
  getCouncilGenesis
  getCouncilLegacyDB
  getCouncilDB
GetValidators
  GetCouncil
  getStakingDemoted
GetCommittee
  GetValidators
  selectRandomCommittee
  selectRandaoCommittee
GetProposer
  GetValidators
  selectRoundRobinProposer
  selectStickyProposer
  selectWeightedRandomProposer
  selectUniformRandomProposer
  selectRandaoProposer

# types
AddressSet
  Add
  Remove
  Has
  ToSortedList

AddressList
  At
  IndexOf
  Swap
  Sort
  ShuffleLegacy
  ShuffleRandao
```