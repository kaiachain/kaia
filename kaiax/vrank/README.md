# kaiax/vrank

This module computes and exposes VRank scores (PFS and CFS) for validators and candidates, and implements the supporting `header.VRank` flow described in [KIP-227](https://github.com/kaiachain/kips/blob/main/KIPs/kip-227.md).

## Concepts

### Overview

VRank measures the reliability of validators and candidates over a rolling epoch of `Epoch = 86400` blocks. Two scores are maintained:

- **PFS (Proposal Failure Score)**: Counts how many times a validator failed to finalize a block as proposer within the current epoch. Each round change at block N adds one failure for the proposer at that round.
- **CFS (Candidate Failure Score)**: Counts how many times a candidate failed to respond to a `VRankPreprepare` message in time within the current epoch. Responses are collected by committee members during block `N`, and the resulting report is committed into `header(N+1).VRank` by the next block proposer.

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

For each committed header `x` in the range, `header(x).VRank` stores the `cfReport` written by the proposer of block `x`. That report is tallied from candidate responses observed during the consensus of block `x-1`: candidates who did not respond with the correct block hash before the timeout are included.

The raw per-block cfReports are first aggregated into a candidate-proposer (CP) matrix:

```
for each x in [epochStart(N), N]:
    reporter(x) = Proposer(x, header(x).Round)
    for each candidate in cfReport(x):
        cpMatrix_N[candidate][reporter(x)] += 1
```

The CP matrix is then passed through a **byzantine filter** before producing the final CFS. The filter discards the top-F reporter totals per candidate, where `F = floor((epochVACount - 1) / 3)`. `epochVACount` is supplied by the caller; the intended source is the epoch snapshot of `ValActive` from AddressBookV2, not the live committee length. This protects against up to F malicious proposers falsely accusing candidates.

```
scores_N(candidate) = sorted list of cpMatrix_N[candidate][reporter] over all reporters
filteredScores_N(candidate) = scores_N(candidate) with the largest F entries removed
CFS(N)[candidate] = sum(filteredScores_N(candidate))
```

On a cold start, the CP matrix is seeded from `GetCandTesting(N)` for the queried block `N`, so those candidates appear in the output even with a score of zero. Cached or checkpointed CP matrices are replayed forward from that seed.

### Scoring epoch

Both scores are epoch-local. `epochStart(N) = N - (N % epoch)`. The epoch-start block itself (`N % epoch == 0`) always has an empty `header.VRank` and does not contribute to CFS.

### Checkpoints

To avoid replaying the entire epoch on every query, PFS and the CP matrix are periodically persisted to the database every `scoreCheckpointInterval = Epoch / 8 = 10800` blocks. A `lastCheckpoint` pointer is also maintained for rewind bookkeeping. On startup, `Init` calls `GetPFS(head)` and `getCPMatrix(head)` to warm the caches, which internally replay only the blocks since the latest surviving checkpoint.

## Persistent schema

All keys are stored in the chain key-value store (`ChainKv`).

- `vrankCheckpoint || Uint64BE(blockNum)`: PFS map and CP matrix at a checkpoint block. Written every `scoreCheckpointInterval` blocks.
  ```
  "vrankCheckpoint" || Uint64BE(N) => RLP({ PFS: [{Addr, Score}], CPMatrix: [{Candidate, Reporter, Score}] })
  ```
- `vrankLastCheckpoint`: Block number of the most recently written checkpoint, used by rewind logic.
  ```
  "vrankLastCheckpoint" => Uint64BE(N)
  ```

## Module lifecycle

### Init

Calls `GetPFS(head)` and `getCPMatrix(head)` to warm both in-memory caches (PFS and CP matrix). Each internally loads the latest DB checkpoint and replays only the blocks since that checkpoint up to the current chain head.

- Dependencies:
  - `ChainKv`: Raw key-value database for checkpoint persistence.
  - `ChainConfig`: Holds the `PermissionlessCompatibleBlock` fork point.
  - `Chain`: Provides block headers.
  - `Valset`: Provides `GetProposer`, `GetCandTesting`, `GetCommittee`.
  - `Randao`: Provides registered BLS public keys for verifying `VRankCandidate` messages.
  - `RoundReader`: Reads the consensus round from committed headers.
  - `NodeKey`: ECDSA private key for signing `VRankPreprepare` and `VRankCandidate` messages.
  - `BlsKey`: BLS secret key for signing `VRankCandidate` messages.

#### Valset dependency

- getters
  - For all getters except `TallyCfReport`, `N` must exist in the header DB.
  - `GetPfReport` and `GetPFS` call `GetProposer` for each block `x` in `[epochStart(N), N]`. In the valset module, `GetProposer` fast-paths committed historical blocks by returning `Chain.Sealer().Author(header)` when `header(x)` exists and its round matches the requested round; otherwise it falls back to block-context lookup.
  - `GetCFS(N, epochVACount)` calls `GetCandTesting(N)` to seed a fresh CP matrix when there is no cache or checkpoint seed.
  - `TallyCfReport(N, round)` is called by the proposer of block `N+1` to fill `header(N+1).VRank`; it queries `GetCandTesting(N)` because it validates messages collected for block `N`.

| Function                             | Valset call                                                               | Notes                                                         |
| ------------------------------------ | ------------------------------------------------------------------------- | ------------------------------------------------------------- |
| `GetPfReport(N)`                     | `GetProposer(N, r)`; `r ∈ [0, header(N).Round)`                           | valset fast-paths committed headers via `Chain.Sealer().Author` |
| `GetCfReport(N)`                     | —                                                                         | —                                                             |
| `TallyCfReport(N, round)`            | `GetCandTesting(N)`                                                       | `N` is the reported block whose responses are being tallied   |
| `GetPFS(N)`                          | `GetProposer(x, r)`; `x ∈ [epochStart(N), N]`, `r ∈ [0, header(x).Round)` | valset fast-paths committed headers via `Chain.Sealer().Author` |
| `GetCFS(N, epochVACount)`            | `GetCandTesting(N)`                                                       | seeds `newCPMatrix` when no cache/checkpoint seed exists      |
|                                      | `GetProposer(x, header(x).Round)`; `x ∈ [epochStart(N), N]`               | reporter lookup per committed header                          |

- handlers
  - During consensus of block N, block N is not yet committed, but its validator/candidate/proposer set is already determined — so `Get*(N)` is safe for proposer / candidate / committee identity checks tied to that live view.
  - `HandleVRankCandidate` deliberately does not perform canonical committee-membership or candidate-membership validation at receive time. It acts as a bounded inbox keyed by `(BlockNumber, Round)`: it requires a current preprepared view, drops stale messages, rejects overly-future messages, verifies ECDSA and BLS signatures, and ignores duplicate senders per view.
  - Canonical semantic validation for candidate messages happens later in `TallyCfReport`, which checks the exact view's preprepare time, expected block hash, the candidate set from `GetCandTesting(N)`, and the timeout.

| Function                                  | Valset call                  | Notes                                                       |
| ----------------------------------------- | ---------------------------- | ----------------------------------------------------------- |
| `HandleIstanbulPreprepare(block N)`       | `GetCommittee(N, 0)`         | committee members collect `VRankCandidate` messages         |
|                                           | `GetProposer(N, view.Round)` | only proposer broadcasts `VRankPreprepare`                  |
|                                           | `GetCandTesting(N)`          | proposer broadcasts `VRankPreprepare` to candidates         |
| `HandleVRankPreprepare(block N)`          | `GetCandTesting(N)`          | only candidates handle `VRankPreprepare`                    |
|                                           | `GetProposer(N, view.Round)` | verifies `VRankPreprepare` sender is the proposer           |
|                                           | `GetCommittee(N, view.Round)` | candidates broadcast `VRankCandidate` to the round committee |
| `HandleVRankCandidate(msg.BlockNumber=M)` | —                             | receive-time path does not query valset                     |

### Start and stop

This module maintains one background goroutine that drains the broadcast channel and forwards VRankPreprepare and VRankCandidate messages to peers.

## Block processing

### Consensus

#### HandleIstanbulPreprepare

Called when the Istanbul core receives a PREPREPARE for block N. If this node is a committee member for block N, it records the preprepare timestamp and block hash for the collector. If this node is also the proposer, it broadcasts a `VRankPreprepare` to all candidates.

#### HandleVRankPreprepare

Called when a `VRankPreprepare` is received. If this node is a candidate for block N, it verifies the proposer's signature, suppresses exact replays, ignores same-view conflicting hashes, and then signs and broadcasts a `VRankCandidate` message to the round-specific committee.

#### HandleVRankCandidate

Called when a `VRankCandidate` message is received. The handler does not try to decide canonical candidate or committee membership for that message's view. Instead, it:

- requires a current preprepared view on the receiving node
- rejects messages that are too far in the future
- drops stale messages for already-passed views
- verifies the candidate's ECDSA and BLS signatures
- ignores duplicate senders for the same `(blockNum, round)`
- records the message in the collector for later tallying

`TallyCfReport` is the stage that applies canonical checks such as expected block hash, candidate membership in `GetCandTesting(N)`, and the timeout rule.

#### VerifyHeader

`VerifyHeader` checks the `VRank` field in committed headers:

- before the permissionless fork, `header.VRank` must be empty
- at epoch-start blocks, `header.VRank` must be empty
- otherwise, `header(N).VRank` must decode to a sorted, deduplicated candidate list whose addresses are all in `GetCandTesting(N-1)`, because `header(N).VRank` reports failures observed while building block `N-1`

### Execution

#### PostInsertBlock

After each canonical block is inserted, proactively calls `GetPFS` and `getCPMatrix` to keep caches warm. At every `scoreCheckpointInterval` boundary, writes the current PFS and CP matrix to the DB and updates the `lastCheckpoint` pointer.

### Rewind

- `RewindTo`: Purges both in-memory caches (PFS and CP matrix).
- `RewindDelete`: Handles exactly one deleted block. If that block is a checkpoint, its DB checkpoint is deleted; if it was also the `lastCheckpoint`, the pointer rolls back to the highest surviving earlier checkpoint.

## Getters

All getters return `ErrNotPermissionless` if the queried block is before the permissionless fork.

- `GetPfReport(N)`: Returns the list of proposers who failed at block N (one per failed round).
- `GetCfReport(N)`: Returns the raw cfReport committed in `header(N).VRank`.
- `TallyCfReport(N, round)`: Computes the cfReport for block N at the given round from the in-memory collector, for use by the proposer of block N+1.
- `GetPFS(N)`: Returns the cumulative PFS map over `[epochStart(N), N]`.
- `GetCFS(N, epochVACount)`: Returns the cumulative CFS map over `[epochStart(N), N]`.

The JSON-RPC wrapper currently resolves `epochVACount` to `1` as a temporary placeholder until the real epoch snapshot source is wired.
