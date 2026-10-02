# kaiax/vrank

This module computes and exposes VRank scores (PFS and CFS) for validators and candidates, and implements the supporting `header.VRank` flow described in [KIP-227](https://github.com/kaiachain/kips/blob/main/KIPs/kip-227.md).

## Concepts

### Overview

VRank measures the reliability of validators and candidates over an epoch of `ChainConfig.VRankEpoch` blocks (default `86400`). Two scores are maintained:

- **PFS (Proposal Failure Score)**: Counts how many times a validator failed to finalize a block as proposer within the current epoch. Each round that was started but did not finalize block `N` adds one failure for the proposer of that round. The failed rounds of block `N` are recorded in `header(N+1).VRank` together with a committee certificate that proves them.
- **CFS (Candidate Failure Score)**: Counts how many times a candidate failed to respond to a `VRankPreprepare` message in time within the current epoch. Responses to the proposer of block `N` are collected by that proposer alone, and the resulting report is committed into the `header.VRank` of the next block that the same proposer builds within the epoch.

Both scores reset to zero at the start of each epoch.

### header.VRank

`header(N).VRank` is the RLP encoding of `VRankPayload`, written by the proposer of block `N` in `PrepareHeader`:

```
VRankPayload {
    Report              []Address  // CandTesting(N) at epoch-start blocks, otherwise a cfReport
    ParentRound         uint8      // a round at which a quorum committed block N-1; 0 when there was no proposal failure,
                                   // at epoch-start blocks, and when N-1 is before the permissionless fork
    ParentCommittedSeal [][]byte   // committed seals of block N-1 for ParentRound; empty when ParentRound == 0
}
```

A payload with an empty report, `ParentRound == 0` and no seals is encoded as an empty field, and an empty field decodes to that payload.

```
during consensus of block N
  proposer(N, round) ──[VRankPreprepare]──▶ CandTesting(N)
  candidates         ──[VRankCandidate]───▶ proposer(N, round)
                                            (the proposer records the preprepare time and each reply time)

when the same proposer builds its next block M (> N) in the epoch
  proposer(M) ── EvaluateCandidates(N, round) ──▶ cfReport ──▶ header(M).VRank.Report
  proposer(M) ── round and committed seals of M-1 ──▶ header(M).VRank.ParentRound / ParentCommittedSeal
```

- `VRankPreprepare{Block, View, Sig}`: signed by `proposer(N, view.Round)`; tells candidates a response is expected for `(N, view.Round)`.
- `VRankCandidate{BlockNumber, Round, BlockHash, Sig, BlsSig}`: signed by each candidate with both its node key and its BLS key; carries the block hash it observed.
- Both signatures cover `keccak256(domain || chainId(uint64 BE) || blockNumber(uint64 BE) || round(uint8) || blockHash)` with the domains `VRANK_PREPREPARE_V1` and `VRANK_CANDIDATE_V1`.
- `EvaluateCandidates(N, round)` lists every address in `GetCandTesting(N)` whose `VRankCandidate` did not arrive with the correct `BlockHash` within `candidateMsgTimeoutMs = 500` of the preprepare time.
- Rounds above `MaxRound = 10` are rejected on receipt and evaluate to an empty report.

### PFS

`PFS(N)` is the cumulative proposal failure score at block N, covering blocks `[epochStart(N), N]`.

For each block `x` in the range, `header(x).VRank.ParentRound` is a round at which a quorum committed block `x-1`. The proposers of block `x-1` at rounds `[0, ParentRound)` each receive one failure — those are the proposers who started a round but did not finalize the block. The proposer at `ParentRound` does not receive a failure.

```
epochStart(N) = N - (N % VRankEpoch)
PFS(N)[addr] = len({x in [epochStart(N), N] : addr in pfReport(x)})
pfReport(x)  = [Proposer(x-1, 0), Proposer(x-1, 1), ..., Proposer(x-1, ParentRound(x) - 1)]
```

`ParentRound` is backed by a certificate: `VerifyHeader` recovers the signers of `ParentCommittedSeal` against `header(x).ParentHash` and the claimed round, requires every signer to be a member of `GetCommittee(x-1, ParentRound)`, and requires at least a quorum of them. The certificate proves that a quorum committed `x-1` at that round, so a proposer cannot claim a higher round than any quorum reached. It does not prove that the round was the last one; the verifier does not compare it with the round stored in its own copy of `x-1`, since honest nodes may hold seals from different rounds for the same block. The epoch-start block and the first block after the fork carry no certificate (`ParentRound == 0`), so failures of the previous block do not count there.

### CFS

`CFS(N)` is the cumulative candidate failure score at block N, covering blocks `[epochStart(N), N]`.

For each committed header `x` in the range, `header(x).VRank.Report` stores the `cfReport` written by the author of block `x`. That report is tallied from the candidate responses the author observed while proposing its previous block in the epoch: candidates who did not respond with the correct block hash before the timeout are included.

The raw per-block cfReports are first aggregated into a candidate-proposer (CP) matrix. Every block contributes its author to the matrix even when its cfReport is empty, so the matrix's reporter set covers all distinct authors seen in the epoch. The author is read with `Sealer.Author(header(x))`: a hash-locked re-proposal carries the original `header.VRank` verbatim, so the proposer of the committing round is not necessarily its reporter.

```
for each x in [epochStart(N), N]:
    reporter(x) = Sealer.Author(header(x))
    cpMatrix_N.AddProposer(reporter(x))
    for each candidate in cfReport(x):
        cpMatrix_N[candidate][reporter(x)] += 1
```

A reported address that is not a row of the matrix is skipped with a warning.

The CP matrix is then passed through a **byzantine filter** before producing the final CFS. The filter discards the top-F reporter totals per candidate, where `F = floor(P/3)` and `P = cpMatrix.ProposerCount()` (the number of distinct authors in the epoch so far). This protects against up to F malicious proposers falsely accusing candidates. Early in an epoch `P` is small and the filter is correspondingly weak, but CFS is only consulted at epoch-end state transitions.

```
scores_N(candidate) = sorted list of cpMatrix_N[candidate][reporter] over all reporters
filteredScores_N(candidate) = scores_N(candidate) with the largest F entries removed
CFS(N)[candidate] = sum(filteredScores_N(candidate))
```

On a cold start, the CP matrix is seeded with the epoch's candidates, so those candidates appear in the output even with a score of zero. The candidates come from `GetCandTesting(N)`; when the state for `N` is unavailable, they are read from `header(epochStart(N)).VRank.Report` instead. Cached or checkpointed CP matrices are replayed forward from that seed.

### Scoring epoch

Both scores are epoch-local. `epochStart(N) = N - (N % VRankEpoch)`. The epoch-start block itself (`N % VRankEpoch == 0`) carries `CandTesting(N)` as its `Report` with `ParentRound == 0` and no seals. It contributes neither a cfReport nor a pfReport, but its author is still added to the CP matrix.

### Checkpoints

To avoid replaying the entire epoch on every query, PFS and the CP matrix are periodically persisted to the database every `scoreCheckpointInterval = VRankEpoch / 8` blocks (`10800` with the default epoch). A `lastCheckpoint` pointer is also maintained for rewind bookkeeping. A lookup first probes the in-memory caches for any of the previous 64 blocks of the epoch, then the checkpoint at `N - (N % scoreCheckpointInterval)`, and otherwise starts from the epoch start.

## Persistent schema

All keys are stored in the chain key-value store (`ChainKv`).

- `vrankCheckpoint || Uint64BE(blockNum)`: PFS map and CP matrix at a checkpoint block. Written every `scoreCheckpointInterval` blocks.
  ```
  "vrankCheckpoint" || Uint64BE(N) => RLP({
      PFS:      { Addrs: [Address], Scores: [uint64] },                                 // same index order
      CPMatrix: { Candidates: [Address], Reporters: [Address], Matrix: [[uint64]] }     // Matrix[candidate][reporter]
  })
  ```
  Addresses are sorted in each list.
- `vrankLastCheckpoint`: Block number of the most recently written checkpoint, used by rewind logic.
  ```
  "vrankLastCheckpoint" => Uint64BE(N)
  ```

## Module lifecycle

### Init

Requires every dependency below, including `ChainConfig.ChainID`. If the chain head is past the permissionless fork, `Init` calls `GetPFS(head)` and `getCPMatrix(head)` to warm both in-memory caches (PFS and CP matrix). Each internally loads the DB checkpoint at `head - (head % scoreCheckpointInterval)` and replays only the blocks since that checkpoint up to the current chain head. A warming failure is logged and the module starts cold.

- Dependencies:
  - `ChainKv`: Raw key-value database for checkpoint persistence.
  - `ChainConfig`: Holds the `PermissionlessCompatibleBlock` fork point, `VRankEpoch`, and `ChainID` (part of the signed message hashes).
  - `Chain`: Provides block headers and the current head.
  - `Valset`: Provides `GetProposer`, `GetCandTesting`, `GetCommittee`.
  - `Randao`: Provides registered BLS public keys for verifying `VRankCandidate` messages.
  - `Sealer`: Reads the committed round, author, and raw committed seals of a header, recovers committers from seals, and gives the quorum size.
  - `NodeKey`: ECDSA private key for signing `VRankPreprepare` and `VRankCandidate` messages.
  - `BlsKey`: BLS secret key for signing `VRankCandidate` messages.

#### Valset dependency

- getters
  - For all getters except `EvaluateCandidates`, `N` must exist in the header DB and must not be beyond the chain head.
  - `GetPfReport(N)` and `GetPFS` call `GetProposer(x-1, r)` for each block `x` in range and each failed round `r` read from `header(x).VRank.ParentRound`.
  - `GetCFS(N)` calls `GetCandTesting(N)` to seed a fresh CP matrix when there is no cache or checkpoint seed, falling back to the epoch-start header when the state is unavailable. Reporters are read from the header author, not from valset.
  - `EvaluateCandidates(N, round)` is called by the node that proposed block `N` at its committed round, when that node builds its next block; it queries `GetCandTesting(N)` because it validates messages collected for block `N`.

| Function                       | Valset call                                                                | Notes                                                         |
| ------------------------------ | -------------------------------------------------------------------------- | ------------------------------------------------------------- |
| `GetPfReport(N)`               | `GetProposer(N-1, r)`; `r ∈ [0, ParentRound(N))`                           | `ParentRound(N)` is read from `header(N).VRank`               |
| `EvaluateCandidates(N, round)` | `GetCandTesting(N)`                                                        | `N` is the reported block whose responses are being evaluated |
| `GetPFS(N)`                    | `GetProposer(x-1, r)`; `x ∈ [epochStart(N), N]`, `r ∈ [0, ParentRound(x))` | one lookup per failed round                                   |
| `GetCFS(N)`                    | `GetCandTesting(N)`                                                        | seeds `NewCPMatrix` when no cache/checkpoint seed exists      |

- handlers
  - During consensus of block N, block N is not yet committed, but its validator/candidate/proposer set is already determined — so `Get*(N)` is safe for proposer / candidate identity checks tied to that live view.
  - `HandleVRankCandidate` deliberately does not perform candidate-membership validation at receive time. It acts as a bounded inbox keyed by `(BlockNumber, Round)`: it accepts only views this node preprepared as proposer, verifies ECDSA and BLS signatures, and ignores duplicate senders per view.
  - Canonical semantic validation for candidate messages happens later in `EvaluateCandidates`, which checks the exact view's preprepare time, expected block hash, the candidate set from `GetCandTesting(N)`, and the timeout.
  - `VerifyHeader(N)` calls `GetCandTesting(N)` to confirm every address in `header(N).VRank.Report` is a candidate of the epoch (`CandTesting` is stable within an epoch), blocking malicious proposers from injecting arbitrary addresses to manipulate CFS scores. It calls `GetCommittee(N-1, ParentRound)` to verify the parent round certificate.

| Function                                  | Valset call                      | Notes                                                               |
| ----------------------------------------- | -------------------------------- | ------------------------------------------------------------------- |
| `HandleIstanbulPreprepare(block N)`       | `GetProposer(N, view.Round)`     | only the proposer records the view and broadcasts `VRankPreprepare` |
|                                           | `GetCandTesting(N)`              | proposer broadcasts `VRankPreprepare` to candidates                 |
| `HandleVRankPreprepare(block N)`          | `GetCandTesting(N)`              | only candidates handle `VRankPreprepare`                            |
|                                           | `GetProposer(N, view.Round)`     | verifies the `VRankPreprepare` sender is the proposer               |
| `HandleVRankCandidate(msg.BlockNumber=M)` | —                                | receive-time path does not query valset                             |
| `VerifyHeader(N)`                         | `GetCandTesting(N)`              | validates `header(N).VRank.Report` addresses are candidates         |
|                                           | `GetCommittee(N-1, ParentRound)` | validates `ParentCommittedSeal` signers and quorum                  |

### Start and stop

This module maintains one background goroutine that drains the broadcast channel into the `SubscribeVRank` feed; the subscriber (`node/cn`) delivers the VRankPreprepare and VRankCandidate messages to peers.

## Block processing

### Consensus

#### HandleIstanbulPreprepare

Called when the Istanbul core receives a PREPREPARE for block N. Blocks before the permissionless fork are ignored. If this node is the proposer for `(N, view.Round)`, it records the preprepare timestamp and block hash in the collector, drops collector views from earlier epochs, and broadcasts a `VRankPreprepare` to all candidates of block N. Other nodes keep no VRank state for the view in this handler.

#### HandleVRankPreprepare

Called when a `VRankPreprepare` is received. Blocks before the permissionless fork are ignored. The view must match the block (`view.Sequence == block.Number`) and `view.Round` must not exceed `MaxRound`. The block must be current: either the next block after the chain head with `ParentHash` equal to the head hash, or the head itself with the same hash and committed round. If this node is a candidate for block N, it verifies that the signer is `GetProposer(N, view.Round)`, rechecks that the block is still current, suppresses exact replays and same-view conflicting hashes within a window of the last 10 blocks, and then signs and sends a `VRankCandidate` message to the proposer.

#### HandleVRankCandidate

Called when a `VRankCandidate` message is received. The handler does not try to decide canonical candidate membership for that message's view. Instead, it:

- ignores messages for blocks before the permissionless fork
- rejects rounds above `MaxRound`
- accepts only views this node preprepared as proposer (others are dropped before signature recovery)
- recovers the sender from the ECDSA signature
- ignores duplicate senders for the same `(blockNum, round)`
- verifies the BLS signature against the sender's BLS public key registered as of `head + 1`
- records the message and its arrival time in the collector for later evaluation

`EvaluateCandidates` is the stage that applies canonical checks such as expected block hash, candidate membership in `GetCandTesting(N)`, and the timeout rule.

#### PrepareHeader

Called when this node builds block N. Before the permissionless fork `header.VRank` is left empty. Otherwise:

- at epoch-start blocks, `Report` is `GetCandTesting(N)`
- at other blocks, `Report` is `EvaluateCandidates(n, round)` for the most recent block `n < N` in the current epoch that this node proposed at its committed round and whose view is still in the collector. Views older than `n` are then dropped from the collector. If there is no such block (first proposal of the epoch, or a restart emptied the collector), the report is empty.
- `ParentRound` and `ParentCommittedSeal` are the round and committed seals read from this node's stored header of block `N-1` when that round is not 0. They are omitted at epoch-start blocks and when `N-1` is before the fork.

#### VerifyHeader

`VerifyHeader` checks the `VRank` field in committed headers:

- before the permissionless fork, `header.VRank` must be absent
- after the fork, `header.VRank` must decode to a `VRankPayload`
- `ParentRound` must be 0 with no seals at epoch-start blocks and when `N-1` is before the fork; a non-zero `ParentRound` must come with at most `len(GetCommittee(N-1, ParentRound))` seals that recover to distinct members of that committee, at least a quorum of them, over `header.ParentHash` and the claimed round
- at epoch-start blocks, `Report` must equal `GetCandTesting(N)`
- otherwise, `Report` must be empty or a sorted, deduplicated list whose addresses are all in `GetCandTesting(N)`

### Execution

#### PostInsertBlock

After each block is inserted, proactively calls `GetPFS` and `getCPMatrix` to keep caches warm. At every `scoreCheckpointInterval` boundary, writes the current PFS and CP matrix to the DB and updates the `lastCheckpoint` pointer. Warming may fail, but errors are logged and ignored; a skipped checkpoint is rebuilt from the epoch start on a later cold lookup.

### Rewind

- `RewindTo`: Purges both in-memory caches (PFS and CP matrix).
- `RewindDelete`: Handles exactly one deleted block. If that block is a checkpoint, its DB checkpoint is deleted; if it was also the `lastCheckpoint`, the pointer rolls back to the highest surviving earlier checkpoint, or is deleted when none survives.

## Getters

`GetPfReport`, `GetPFS` and `GetCFS` return `ErrNotPermissionless` if the queried block is before the permissionless fork; `EvaluateCandidates(N, round)` returns it if `N+1` is before the fork, since its report is written into block `N+1` or later. `GetPFS` and `GetCFS` also return `ErrFutureBlock` if the queried block is beyond the chain head.

- `GetPfReport(N)`: Returns the proposers of block `N-1` that failed at rounds `[0, ParentRound(N))`, as recorded in `header(N).VRank`.
- `EvaluateCandidates(N, round)`: Computes the cfReport for block N at the given round from the in-memory collector, for use by the node that proposed block N at its committed round when it builds its next block.
- `GetPFS(N)`: Returns the cumulative PFS map over `[epochStart(N), N]`.
- `GetCFS(N)`: Returns the cumulative CFS map over `[epochStart(N), N]`. The byzantine-filter threshold is derived internally from `cpMatrix.ProposerCount()`, the number of distinct authors seen so far in the epoch.

## APIs

Exposed under the `governance` namespace. A missing, `latest`, or `pending` block number resolves to the chain head.

- `governance_getPFS(blockNumber)`: `GetPFS` for the given block.
- `governance_getCFS(blockNumber)`: `GetCFS` for the given block.
