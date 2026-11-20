# Transaction processing metrics

This optional collector identifies where local transaction processing is waiting:
pool admission and reset, promotion notifications, peer writes, and candidate
execution. It does not change scheduling, pool policy, or propagation policy. It
adds no logs, transaction hashes, account addresses, or peer identities.

## Enable and collect

Start a node with `KAIA_POOLOBSERVE=1`, `--metrics`, and `--prometheus`.
Without the environment variable the hooks do not record observations and export no `poolobs`
series. The async queue envelope still reserves a 24-byte timestamp per batch
slot, about 3 KiB per peer at capacity 128, plus one observer pointer. Disabled
mode avoids the observer object, clocks and its lock, but not this envelope storage. Enabling the environment variable without the exporter still incurs the
recording cost and makes the observations unavailable remotely.

Schema 9 exports **228 fixed sample series per node in 29 families**, including histogram
buckets and metadata samples. Counts do not grow with transactions,
accounts, peers, candidate heights, or restarts. Native Prometheus collectors
are updated directly; they do not wait for the legacy metrics-to-Prometheus
conversion interval. The collector still has to scrape the endpoint. The central
pipeline must retain `stage`, `queue`, `event`, `phase`, `reason`, `result`, and `le`
labels and the cumulative counters/buckets. Validate that pipeline before rollout.

## Questions and measurements

### Has the local chain advanced while the pool still uses older state?

Compare `kaia_poolobs_chain_head_number`, `pool_state_head_number`, and
`pool_reset_head_number` (all with the `kaia_poolobs_` prefix). Require the
corresponding `*_observed` gauges to be one. State installation and full reset
completion are separate boundaries. A persistent gap narrows the investigation
to the local head/reset path. The gauges are not an atomic snapshot, so a single
one-block difference is not proof of a stall, and reorganizations need care.
`block_import` times the entire public `InsertChain` call, including synchronization
and event delivery, not just EVM execution. A slow import and a slow reset are
separate observations; do not infer either from the other.

### Are pool locks, admission, or promotion notification blocking progress?

Every stage exposes `kaia_poolobs_operations_started_total`,
`operations_inflight`, and `duration_seconds` (all with the same prefix).
The histogram's `kaia_poolobs_duration_seconds_count` is the completion counter;
there is no duplicate standalone completed counter. Durations use monotonic time;
completion includes unsuccessful operations. Inflight distinguishes work that
has not finished and therefore has not entered the duration histogram.

- `admission_write_single`, `admission_write_batch`: wait for the public admission
  pool lock. `admission_hold`: work while holding that lock.
- `head_reset_wait`: wait for `pool.mu` on a head update. `reset_body`: reset
  processing, including inner lock waits; `reset_hold`: time holding `pool.mu`
  around that processing. These overlap; do not add durations.
- `pending_wait`: wait for `pool.mu` in `Pending`, excluding the later `txMu` wait.
  `candidate_pending`: the complete candidate call to `Pending`.
- `remote_enqueue`: delivery into the remote-admission channel;
  `remote_admission`: consumer processing through `AddRemotes`.
- `feed_enqueue`: enqueue a promotion event batch; `feed_send`: deliver a batch
  to event subscribers. Neither counts individual transactions.

`kaia_poolobs_queue_{observed,slots,capacity}{queue="feed|broadcast|remote"}`
measures channel slots containing batches, not transactions. A zero observed
flag means the source is unavailable. Queue gauges can miss brief peaks between
scrapes; wait histograms and cumulative counts preserve completed short events.

`kaia_poolobs_events_total{event=...}` uses the following units:

- `admitted` / `admission_rejected`: completed public admission attempts, one per
  input transaction. Rejections include duplicates; success includes replacement.
  Internal re-injection is excluded. These are not RPC response timestamps.
- `promoted`: promotion events before notification and possible later eviction;
  not a durable pending-pool population.
- `remote_received`: transactions presented to the handler, before denial or
  throttling. Not successful pool admission.
- `remote_admitted` / `remote_rejected`: public remote admission results, a subset
  of the overall admitted/rejected totals. Local RPC admission is excluded.
- `capacity_rejected`: actual capacity rejection branches, including internal
  reinjection and batch overflow tails. It is not restricted to remote input.
- `remote_denied`: input discarded by the remote-denial policy.
- `remote_throttled`: transactions successfully handed into the throttle queue.
  Its existing drop counter still accounts for overflow.

Use existing pool size, underpriced, invalid/refused, replacement, and rate-limit
metrics to refine the reason. Do not subtract admitted and promoted to construct
a loss rate: their populations and retry behavior differ.

### Is local transmission blocked after transactions become ready?

- `broadcast`: the broadcast operation; on CN this includes dispatch into async
  queues, not completion of each peer write.
- `peer_send`: the common actual transaction send path, including CN async
  consumers, EN direct sends, retained PN compatibility paths, resend, and
  multichannel implementations. Current deployment targets are EN and CN.
- `rebroadcast`: the resend operation.
- Event counters `broadcast_transactions` and `rebroadcast_transactions` count
  operation input transactions; `peer_send_transactions` counts transaction
  copies in each attempted send batch. A failed early batch does not count its
  unattempted tail. `peer_send_errors` counts failed logical peer calls.
- `peer_send_dropped` counts copies rejected by a full CN async send queue.
- `broadcast_no_recipients` / `rebroadcast_no_recipients` count nonempty calls
  whose selected recipient maps are entirely empty. This can be normal known-tx
  suppression; it is not an error or a count of unsent individual transactions.

CN asynchronous enqueue-to-dequeue latency is recorded separately in
`kaia_poolobs_peer_queue_delay_seconds`. It includes the nonblocking enqueue
attempt and ends before the actual write. `kaia_poolobs_peer_queue_batches`
tracks unfinished batch observations across peers, with a provisional increment
just before enqueue and immediate rollback on a full queue. On peer close,
remaining observations are removed from this gauge and counted in
`kaia_poolobs_peer_queue_abandoned_batches_total`. Abandoned observations are
not proof of transaction loss. No peer label is exported. Each peer observer has
a short local mutex that never covers a channel operation or socket write.
A batch carries one timestamp; no per-transaction lookup table is added.

A successful write is not a remote pool acknowledgement. Per-peer copies,
retries, and resends can count the same transaction repeatedly. Aggregate local
send/receive counters from different nodes are not a delivery-loss equation.

### Did the local candidate lose transactions before or during execution?

`kaia_poolobs_candidate_{transactions,empty}_total{phase=...}`
records summed transaction counts and zero-transaction observations.
`kaia_poolobs_candidate_preparations_total` counts successful pending snapshots
once for the three preparation phases. The `candidate_build` histogram count
counts completed execution attempts, including empty attempts:

- `pool_input`: copied pending map before fee filtering.
- `after_fee`: map after fee filtering.
- `after_modules`: map after module filtering.
- `execution_input`: concrete transaction slice before bundle construction.
- `included`: actual transactions in the completed local candidate.

`candidate_build` measures `ApplyTransactions`, including bundle construction,
not the ideal-block-time wait or the whole consensus round. Candidates may be
repeated, abandoned, or never proposed. Generated transactions can make included
counts exceed execution input; phase differences are not canonical loss counts.
Execution input is observed before completion, so averages using the completion
count require a settled time window; they are not exact instantaneous ratios.

`kaia_poolobs_candidate_outcomes_total{reason=...}` adds eight bounded branch counts:

- `bundle_target_rejected`: target receipt does not allow the bundle.
- `generator_failed`: generating a transaction fails.
- `size_limit`: a transaction or bundle would exceed the block size limit.
- `blob_space_rejected`: missing sidecar or insufficient blob capacity.
- `unsupported_type`: unsupported transaction type.
- `bundle_reverted`: bundle reversion.
- `nonce_low`, `nonce_high`: nonce rejection decisions, preserved cumulatively.

These count branch decisions, not affected transactions. One rejected bundle can
contain several transactions. Existing `miner/timelimitreached`,
`miner/toolongtx`, `miner/usedalltxs`, and `miner/strangeerror/txs` cover execution
budget and other execution outcomes. Existing `miner/nonce/low/txs` and
`miner/nonce/high/txs` are per-call gauges, so their short-lived values can be
missed by scraping. The new nonce reason counters preserve the corresponding branch events even when the legacy gauge has already been overwritten.

### What happened in candidates this node actually proposed and adopted?

The existing `candidate_*` families cover all local preparation and execution
attempts. The new `kaia_poolobs_adopted_candidate_*` families cover only sealed
results consumed by the local worker and confirmed to be local-author canonical
blocks at publication time. An already-imported result is eligible when its hash
matches the local canonical header at that height.

`adopted_candidate_blocks_total{result=...}` distinguishes three populations:

- `observed`: preparation and execution observations belong to the matched
  attempt and recorder, with every phase recorded. All detail below uses this
  population.
- `unobserved`: the result is local-author canonical, but the matching complete
  observations are unavailable. Missing data must not look like an empty block.
- `not_adopted`: a consumed sealed result fails the canonical or local-author
  conditions, including write errors. Nil and nonproposer results are excluded.

For `observed` results, `adopted_candidate_transactions_total{phase}` sums
transaction counts at the same five phases defined above, and
`adopted_candidate_empty_total{phase}` counts zero-transaction phases.
`adopted_candidate_outcomes_total{reason}` counts the same eight execution branch
choices and two mutually exclusive deadline outcomes. When the loop exits on its timer flag or an execution time-limit error, `time_limit_with_remaining` means worklist entries remained and
`time_limit_without_remaining` means the list was drained. Remaining entries
can include bundle generators or ineligible transactions. These counters apply
only to matched adopted executions. Neither identifies an excluded hash. `adopted_candidate_execute_seconds` measures the executor call,
including its lock wait, transaction selection and application, but excluding
the preceding ResetWithState call, finalization and sealing. It is not an EVM-only duration or a direct count of
execution-budget timeouts. Its buckets match candidate_build so the two populations can be compared.

These statistics separate adopted local proposals from speculative candidates.
For example, compare `pool_input`, `after_fee`, and `included` within the adopted
population to distinguish a small input from filtering or execution exclusion.
The complete-observation counter is the denominator for per-block averages; do
not divide these phase totals by all candidate preparations.

The worker may replace a result channel before consuming its result. Such a
result is absent from these counters even when its block was imported through
another path. Therefore this is not an exhaustive canonical-block counter, and
`observed + unobserved` must not be treated as the total number of blocks authored
by the node. Compare with chain data when checking coverage. As with other
cumulative metrics, publication is not an atomic multi-family scrape. Use settled
windows; the metrics do not identify individual block or transaction hashes.

### Diagnosing low inclusion on mainnet

1. Count transactions on chain and match the actual proposer to the exporting node. Separate a block with zero transactions from one with zero test transactions. Check generator send times when attributing a supply gap. Node counters include other traffic.
2. Check the observed adopted-block count before interpreting phase deltas. For a settled interval containing exactly one observed adopted block on that node, phase deltas give that candidate's counts. A small pool input points upstream; a drop after fee or module filtering locates the reduction inside preparation. For multiple observed blocks, the same deltas are totals across candidates.
3. If input reached execution and inclusion was low, check the adopted deadline outcome with remaining work. A slow executor call alone cannot establish that the deadline stopped selection. Conversely, a short call on a sparse candidate makes the pool input and filters especially informative.
4. For a small candidate input, compare submission-node admission, readiness stops and send pressure with receiving-node handoff and admission. Inflight operations expose a stall before a duration sample exists. Sender and recipient counters can narrow the affected stage across the workload; they cannot pair an individual send with a receipt.

Choose internal duration bounds from the operation being timed. RPC response, header-based inclusion and first receipt observation have different endpoints. A long tail in one of those measurements does not by itself justify the same bucket in every node stage.

Finite duration boundaries are specific to the operation. They are diagnostic thresholds, not deadlines. The default 250ms execution budget applies to candidate execution. Compare matching stages across nodes. For comparisons across stages or schema versions, use shared thresholds; retain the stage label when computing quantiles.

| Measurement | Upper bounds | Diagnostic use |
|---|---|---|
| `admission_write_single` | 100ms, 250ms, 1s | Detect lock contention and count waits that consume a quarter or a full block interval. |
| `admission_write_batch` | 100ms, 1s | Identify batch lock contention on the remote admission path. |
| `head_reset_wait` | 10ms, 100ms, 1s | Separate delayed lock acquisition from time spent updating the pool. |
| `pending_wait` | 10ms, 100ms, 1s | Identify pool-lock delay inside the pending snapshot. |
| `reset_body` | 10ms, 100ms, 1s | Distinguish ordinary refresh work from a reset spanning a block interval. |
| `reset_hold` | 10ms, 100ms, 1s | Compare lock occupancy with the reset body at the same thresholds. |
| `feed_enqueue` | 100ms, 1s | Identify blocking while handing a promotion to the notification queue. |
| `feed_send` | 100ms, 1s | Distinguish subscriber delivery delay from queue submission. |
| `admission_hold` | 10ms, 100ms, 1s | Separate time holding the shared lock from time waiting for it. |
| `remote_enqueue` | 100ms, 1s | Identify backpressure before the remote admission consumer. |
| `remote_admission` | 10ms, 100ms, 1s | Separate slow batch processing from enqueue backpressure. |
| `broadcast` | 100ms, 1s | Identify a slow dispatch cycle and compare it with writes. |
| `peer_send` | 100ms, 1s | Compare write delays with async queue residence at both thresholds. |
| `rebroadcast` | 100ms, 1s | Detect slow resend calls; this timer excludes the wait between retries. |
| `block_import` | 10ms, 100ms, 250ms, 1s | Separate ordinary imports from increasingly expensive imports and head lag. |
| `candidate_build` | 10ms, 100ms, 250ms, 1s | Distinguish fast builds, substantial execution, and the default budget boundary. |
| `candidate_pending` | 10ms, 100ms, 1s | Compare total snapshot latency with its inner lock wait. |
| `peer_queue_delay` | 100ms, 1s | Identify async batch residence before the peer write. |
| `adopted_candidate_execute` | 10ms, 100ms, 250ms, 1s | Compare adopted executor calls with candidate builds; use outcome counters for deadline exits. |

Schema 7 added three `promotion_stops_total{reason}` counters and the two
adopted deadline outcomes. A fee-aware scan that stops early records its first
unmet condition: a missing next nonce (`nonce_gap`), a price below the base
fee supplied to that scan (`below_base_fee`), or a pool module withholding
readiness (`module_not_ready`). Repeated scans of one account count again.
Empty and fully drained scans record no stop. These counters cover the
Magma fee-aware readiness path after earlier pool checks. A nonce gap describes the
queue at scan time, even if an earlier balance check removed the missing entry.

Schema 7 also counts completed incoming transaction messages in
`tx_messages_total{result}`. Disabled reception, decode errors and invalid
blobs end before pool handoff. A truncated legacy list or skipped nil item is
partial when handed off. Forwarded means the full decoded list was handed off,
including an empty list. A pool policy may still deny every transaction. These
are per-message counts, unlike the per-copy `remote_received` event. No
message is counted while its handler is still blocked.

Schema 6 added five families and 31 sample series: three result counters, five
phase totals, five empty-phase counters, eight outcome counters, and ten
histogram samples. Existing family names and populations remain unchanged.

## Interpretation and limits

Start with completion rates, mean durations and the fraction exceeding a bucket
boundary, then inspect inflight work and head gaps. These examples retain both
common source-node labels; use the identifying labels retained by your collector.
Evaluate duration ratios only where the completion rate is greater than zero.

```promql
# Completed operations per second, by source node and stage.
sum by (job, instance, host, stage) (rate(kaia_poolobs_duration_seconds_count[5m]))

# Mean completed duration in seconds.
sum by (job, instance, host, stage) (rate(kaia_poolobs_duration_seconds_sum[5m]))
/
sum by (job, instance, host, stage) (rate(kaia_poolobs_duration_seconds_count[5m]))

# Fraction of completed operations taking more than 100ms.
1 - (
  sum by (job, instance, host, stage) (rate(kaia_poolobs_duration_seconds_bucket{le="0.1"}[5m]))
  /
  sum by (job, instance, host, stage) (rate(kaia_poolobs_duration_seconds_count[5m]))
)

increase(kaia_poolobs_events_total{event="peer_send_dropped"}[5m])
```

Preserve the source node's `instance` or `host` label in the scraper and central
pipeline. These target labels identify the node exporting a measurement; they
are separate from a remote peer identity. The new collector does not add peer
labels, so peer queue and write measurements aggregate all peers of that source
node. They locate a source node's transmission pressure but do not identify which
remote connection is responsible. Keep host identity for diagnosis rather than
summing unrelated nodes. Handle
process restart/counter reset. Use the observation process-start timestamp as a restart hint, not a globally
unique identity. Every histogram includes a +Inf bucket, sum and count. Use the per-stage table above for finite bounds. The mean and the fraction above a measured threshold are directly computable. Quantiles interpolate within buckets; in the overflow bucket there is no observed upper bound.

The metrics can locate aggregate stalls and show local candidate reasons. They
cannot retrospectively prove the path of a chosen hash, join candidate snapshots
to canonical block hashes, or tell what an unobserved proposer received. Do not
label an external node faulty from the absence of local evidence. Those are
explicit limits of this metrics-only change, not a reason to add unbounded labels.

At 10-second scraping, 228 series produce 1,969,920 samples per node per day;
at 30 seconds, 656,640. Text export bytes depend on current numeric values.
Measure compressed transport and storage separately. Start with owned nodes and
preserve histogram/counter semantics through the collector before widening scope.


## Resolution at each boundary

Existing pool-size gauges refresh in the node's eight-second statistics loop;
legacy registry metrics then pass through the three-second Prometheus bridge.
Those intervals are distinct from Telegraf/other scraper intervals, central
storage aggregation, and Grafana query step. This native collector removes only
the legacy bridge delay for the new metrics, not the downstream intervals.

A histogram retains completed duration bucket counts between scrapes, even if
an operation starts and ends between them. It does not retain the exact event
time. Retention assumes endpoint samples from the same process survive; restarts,
missing samples and inappropriate downsampling can destroy that evidence.
In-flight, channel-depth and head gauges remain instantaneous and can miss a
transient event. Started/histogram-count/inflight are also not atomically scraped as
one snapshot. Peer queue durations now cover completed async buffering delays;
feed/remote enqueue timers still measure enqueue blocking, not every buffered
batch's residence time.

For a bounded one-block investigation, collect these new families at one-second
intervals from selected owned nodes and preserve raw samples. Use a query step
no coarser than the retained samples; a finer step cannot create missing data.
Use a sufficiently long window for rate/quantile estimates, or inspect raw
counter deltas for a short interval. The one-second experimental rate is a
proposal, not a claim that the deployed central pipeline currently supports it.
At one second 228 series add 410,400 samples per node for 30 minutes. Start with
owned nodes, and choose a slower sustained rate based on retention/cost evidence.
Even one-second sampling cannot establish the order of subsecond events on
separate nodes or correlate an arbitrary transaction hash.


## Removed duplication (schema 3 to 5)

- Removed 17 standalone completion counters: histogram `_count` already counts
  exactly the same completed operations. Started remains a true arrival counter;
  inflight remains necessary for operations that never finish.
- Replaced eight bounds with seven in 18 histograms: 18 fewer series. The
  boundaries are 100ms, 150ms, 250ms, 1s, 5s, 10s and 20s.
- Replaced five per-phase observation counters with one preparation counter;
  completed execution counts come from the build histogram: four fewer series.
- Removed three random boot identity/validity series inherited from detailed
  capture design. The start timestamp and counter reset handling remain.

All 17 stage boundaries and all admission, transmission and candidate reasons
remain. Head/queue validity flags also remain: unavailable data must not be
interpreted as zero backlog. Total: 310 to 268 series (13.5% reduction).
