# kaiax/vrank

This module is responsible for computing and exposing VRank scores (PFS and CFS) for validators and candidates, as specified in [KIP-227](https://github.com/kaiachain/kips/blob/main/KIPs/kip-227.md).

## Concepts

### Overview

VRank measures the reliability of validators and candidates over a rolling epoch of `Epoch = 86400` blocks. Two scores are maintained:

- **PFS (Proposal Failure Score)**: Counts how many times a validator failed to finalize a block as proposer within the current epoch. Each round change at block N adds one failure for the proposer at that round.
- **CFS (Candidate Failure Score)**: Counts how many times a candidate failed to respond to a VRankPreprepare message in time within the current epoch. Responses are collected by validators and committed into `header.VRank` by the block proposer.

Both scores reset to zero at the start of each epoch.

### PFS

`PFS(N)` is the cumulative proposal failure score at block N, covering blocks `[epochStart(N), N]`.

For each block `x` in the range, the proposers at rounds `[0, header(x).Round)` each receive one failure — those are the proposers who started a round but did not finalize the block. The proposer at `header(x).Round` finalized the block and does not receive a failure.

```
epochStart(N) = N - (N % Epoch)
PFS(N)[addr] = len({x in [epochStart(N), N] : addr in pfReport(x)})
pfReport(x)  = [Proposer(x, 0), Proposer(x, 1), ..., Proposer(x, header(x).Round - 1)]
```

`pfReport` is derived entirely from `header(x).Extra` (the round byte) and the valset module. It cannot be forged because the round byte is part of the header hash that 2F+1 validators sign with their committed seals.

### CFS

`CFS(N)` is the cumulative candidate failure score at block N, covering blocks `[epochStart(N), N]`.

For each block `x` in the range, `header(x).VRank` stores `cfReport(x)`, written by the proposer of block `x`. That report is tallied from candidate responses observed during the consensus of block `x-1`: candidates who did not respond with the correct block hash before the timeout are included in `cfReport(x)`.

The raw per-block cfReports are first aggregated into a candidate-proposer (CP) matrix:

```
for each x in [epochStart(N), N]:
    reporter(x) = Proposer(x, header(x).Round)
    for each candidate in cfReport(x):
        cpMatrix_N[candidate][reporter(x)] += 1
```

The CP matrix is then passed through a **byzantine filter** before producing the final CFS. The filter discards the top-F reporter totals per candidate, where `F = floor((committeeSize - 1) / 3)`. This protects against up to F malicious proposers falsely accusing candidates.

```
scores_N(candidate) = sorted list of cpMatrix_N[candidate][reporter] over all reporters
filteredScores_N(candidate) = scores_N(candidate) with the largest F entries removed
CFS(N)[candidate] = sum(filteredScores_N(candidate))
```

Candidates are seeded from `GetCandidates(epochStart(N))`, and thus every candidate appears in the output even with a score of zero.

### Scoring epoch

Both scores are epoch-local. `epochStart(N) = N - (N % epoch)`. The epoch-start block itself (`N % epoch == 0`) always has an empty `header.VRank` and does not contribute to CFS.

### Checkpoints

To avoid replaying the entire epoch on every query, PFS and the CP matrix are periodically persisted to the database every `scoreCheckpointInterval = Epoch / 8 = 10800` blocks. A `lastCheckpoint` pointer enables O(1) lookup of the most recent checkpoint. On startup, `catchUp` replays only the blocks since the latest surviving checkpoint.

## Persistent schema

All keys are stored in the chain key-value store (`ChainKv`).

- `vrankCheckpoint || Uint64BE(blockNum)`: PFS map and CP matrix at a checkpoint block. Written every `scoreCheckpointInterval` blocks.
  ```
  "vrankCheckpoint" || Uint64BE(N) => RLP({ PFS: [{Addr, Score}], CPMatrix: [{Candidate, Reporter, Score}] })
  ```
- `vrankLastCheckpoint`: Block number of the most recently written checkpoint.
  ```
  "vrankLastCheckpoint" => Uint64BE(N)
  ```

## Module lifecycle

### Init

Loads the latest DB checkpoint and replays blocks up to the current chain head to warm both in-memory caches (PFS and CP matrix).

- Dependencies:
  - `ChainKv`: Raw key-value database for checkpoint persistence.
  - `ChainConfig`: Holds the `PermissionlessCompatibleBlock` fork point.
  - `Chain`: Provides block headers.
  - `Valset`: Provides `GetProposer`, `GetCandidates`, `GetCommittee`.
  - `NodeKey`: ECDSA private key for signing VRankCandidate messages.

#### Valset dependency

- getters
  - For all getters except `TallyCfReport`, `N` must exist in the header DB.
  - `GetPfReport` and `GetPFS` call `GetProposer` at the block's own number — correct because proposer rotation is per-block, not per-epoch.
  - `GetCFS` queries `GetCandidates` and `GetCommittee` at `epochStart` — the candidate set and committee size are fixed for the whole epoch.
  - `TallyCfReport(N, round)` is called by the proposer of block `N+1` to fill `header(N+1).VRank`; it uses `calcEpochStart(N+1)` for candidates because the cfReport belongs to block `N+1`'s epoch.

| Function                  | Valset call                         | Notes                                                        |
| ------------------------- | ----------------------------------- | ------------------------------------------------------------ |
| `GetPfReport(N)`          | `GetProposer(N, x)`                 | each `x` in `[0, header(N).Round)`                           |
| `GetCfReport(N)`          | —                                   | —                                                            |
| `TallyCfReport(N, round)` | `GetCandidates(N)`                  | `N+1` because `header(N+1).VRank = TallyCfReport(N, R)`      |
| `GetPFS(N)`               | `GetProposer(x, r)`                 | each `x` in `[start, N]`, each `r` in `[0, header(x).Round)` |
| `GetCFS(N)`               | `GetCandidates(N)`                  | newCPMatrix                                                  |
|                           | `GetProposer(x, header(x).Round())` | each `x` in `[start, N]`                                     |
|                           | `GetCommittee(N, 0)`                | for byzantine filtering `F` (TODO: remove me)                |

- handlers
  - During consensus of block N, block N is not yet committed, but its validator/candidate/proposer set is already determined — so `Get*(N)` is safe for proposer / candidate / committee identity checks tied to that live view.
  - `HandleVRankCandidate` deliberately does not perform committee-membership or candidate-membership validation at receive time. It acts as a bounded inbox keyed by `(BlockNumber, Round)`: stale messages are dropped, overly-future messages are rejected, signatures are checked, and duplicate senders per view are ignored.
  - Canonical semantic validation for candidate messages happens later in `TallyCfReport`, which checks the exact view's preprepare time, expected block hash, canonical epoch-start candidate set, and timeout.

| Function                                  | Valset call                  | Notes                                                 |
| ----------------------------------------- | ---------------------------- | ----------------------------------------------------- |
| `HandleIstanbulPreprepare(block N)`       | `GetCommittee(N, 0)`         | only committee handles Preprepare                     |
|                                           | `GetProposer(N, view.Round)` | only proposer broadcasts VRankPreprepare              |
|                                           | `GetCandidates(N)`           | only proposer broadcasts VRankPreprepare              |
| `HandleVRankPreprepare(block N)`          | `GetCandidates(N)`           | only candidates handle VRankPreprepare                |
|                                           | `GetProposer(N, view.Round)` | verify VRankPreprepare sender is the proposer         |
|                                           | `GetCommittee(N, 0)`         | only candidates broadcast VRankCandidate to committee |
| `HandleVRankCandidate(msg.BlockNumber=M)` | —                            | receive-time path does not query valset               |

### Start and stop

This module maintains one background goroutine that drains the broadcast channel and forwards VRankPreprepare and VRankCandidate messages to peers.

## Block processing

### Consensus

#### HandleIstanbulPreprepare

Called when the Istanbul core receives a PREPREPARE for block N. If this node is a validator for block N, it records the preprepare timestamp and block hash for the collector. If this node is also the proposer, it broadcasts a VRankPreprepare to all candidates.

#### HandleVRankPreprepare

Called when a VRankPreprepare is received. If this node is a candidate for block N, it signs and broadcasts a VRankCandidate message to all validators.

#### HandleVRankCandidate

Called when a VRankCandidate message is received. The handler does not try to decide canonical candidate or committee membership for that message's view. Instead, it:

- rejects messages that are too far in the future
- drops stale messages for already-passed views
- verifies the candidate signature
- ignores duplicate senders for the same `(blockNum, round)`
- records the message in the collector for later tallying

`TallyCfReport` is the stage that applies canonical checks such as expected block hash, candidate membership in `GetCandidates(calcEpochStart(N+1))`, and the timeout rule.

### Execution

#### PostInsertBlock

After each canonical block is inserted, proactively calls `GetPFS` and `GetCFS` to keep caches warm. At every `scoreCheckpointInterval` boundary, writes the current PFS and CP matrix to the DB and updates the `lastCheckpoint` pointer.

### Rewind

- `RewindTo`: Purges both in-memory caches (PFS and CP matrix).
- `RewindDelete`: Handles exactly one deleted block. If that block is a checkpoint, its DB checkpoint is deleted; if it was also the `lastCheckpoint`, the pointer rolls back to the highest surviving earlier checkpoint.

## Getters

All getters return `ErrNotPermissionless` if the queried block is before the permissionless fork.

- `GetPfReport(N)`: Returns the list of proposers who failed at block N (one per failed round).
- `GetCfReport(N)`: Returns the raw cfReport committed in `header(N).VRank`.
- `TallyCfReport(N, round)`: Computes the cfReport for block N at the given round from the in-memory collector, for use by the proposer of block N+1.
- `GetPFS(N)`: Returns the cumulative PFS map over `[epochStart(N), N]`.
- `GetCFS(N)`: Returns the cumulative CFS map over `[epochStart(N), N]`.
