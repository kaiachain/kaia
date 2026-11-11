# consensus/istanbul

This package is responsible for implementing Kaia's Istanbul BFT consensus runtime.

## Concepts

### Overview

The implementation is split into three layers.

- Top-level `consensus/istanbul`
  - Shared protocol types such as `View`, `Preprepare`, `Subject`, and `ConsensusMsg`
  - Shared events such as `RequestEvent`, `MessageEvent`, `ChainHeadEvent`, and `NewSequenceEvent`
  - Shared `IstanbulSealer` implementation for header encoding, hashing, and signatures
- `consensus/istanbul/core`
  - Consensus state machine for `PREPREPARE`, `PREPARE`, `COMMIT`, and round change
- `consensus/istanbul/backend`
  - Integration layer between the state machine and the node, blockchain, executor, and p2p stack

The package is normally constructed through [consensus/bft](../bft/README.md), which exposes the generic `consensus.Engine` interface used by the rest of the node.

### View

Consensus progresses on a `View`, which is a pair of:

- `Sequence`: the block number being finalized
- `Round`: the retry number for that block number

For each new sequence, consensus starts from round `0`. When the proposer fails to finalize a block in time or the proposal becomes invalid, validators move to a higher round.

### Consensus phases

Istanbul uses three voting phases.

1. `PREPREPARE`
   - The proposer broadcasts a block proposal for the current `(sequence, round)`.
2. `PREPARE`
   - Committee members vote for the proposal digest and parent hash.
3. `COMMIT`
   - Committee members attach committed seals for the proposal hash.

When enough `COMMIT` messages are collected, the core asks the backend to commit the proposal. If progress stalls, the core sends `ROUND CHANGE` messages and starts a higher round.

### Quorum and fault tolerance

For a given sequence and round, the core derives:

- `qualifiedLen`: the number of qualified validators
- `committeeSize`: the effective committee size from governance
- `f`: the maximum number of faulty validators tolerated by the round
- `requiredMsgCnt`: the quorum size needed to move consensus forward

The effective committee size is `min(qualifiedLen, committeeSize)`.

- If the effective size is less than `4`, quorum is the effective size itself.
- Otherwise, quorum is `ceil(2 * effectiveSize / 3)`.
- Fault tolerance is `ceil(effectiveSize / 3) - 1`.

This means the round state is not driven by the whole council directly. It is driven by the qualified validators and the governance-selected committee size for that block.

### Header extra-data

Consensus metadata is stored in `header.Extra` in Istanbul format.

- First 32 bytes: vanity area
- Last byte of the vanity area: round number
- Remaining RLP payload:
  - `Validators []common.Address`
  - `Seal []byte`
  - `CommittedSeal [][]byte`

`IstanbulSealer` is responsible for:

- parsing and caching `IstanbulExtra`
- reading validators, author, committers, and round from a header
- writing validators, proposer seal, and committed seals into a header
- computing `SigHash` and consensus-aware `HeaderHash`
- creating proposer and committed signatures

### Signature and hash basis

Istanbul uses different hash inputs for different signatures.

- Proposer seal
  - Signed over `SigHash(header)`
  - The hash is calculated from the header with proposer seal removed and committed seals removed
- Committed seal
  - Signed over `Keccak256(prepareCommittedSeal(HeaderHash(header)))`
  - `prepareCommittedSeal` appends the commit message code marker to the proposal hash
- Block hash
  - `HeaderHash(header)` is the consensus-aware block hash
  - It keeps the proposer seal but removes committed seals before hashing

This distinction is important because:

- proposer identity is recovered from the proposer seal over `SigHash`
- committers are recovered from committed seals over `HeaderHash`
- the block hash used by the chain is not the same thing as the proposer-signing preimage

### kaiax integration

Istanbul depends on kaiax modules for validator and governance logic.

- `kaiax/valset`
  - Provides council, qualified validators, committee, and proposer for a given block and round
- `kaiax/gov`
  - Provides committee size and other governance-driven parameters used by consensus

`backend.RegisterKaiaxModules(mGov, mValset)` wires those modules into both the backend and the core. Without them, proposal preparation and validation cannot proceed.

### Networking

The backend exposes an Istanbul subprotocol.

- Protocol name: `istanbul`
- Consensus message code: `0x11`

Network messages are carried as `ConsensusMsg{PrevHash, Payload}` where `Payload` is an RLP-encoded core message.

To reduce redundant gossip, the backend keeps:

- `recentMessages`: messages recently seen from each peer
- `knownMessages`: messages already seen locally

`GossipSubPeer` narrows broadcast targets to the union of the current-round and next-round committees.

## Persistent schema

This package does not define its own persistent schema.

Blocks, headers, receipts, and governance state are persisted by other components.

## In-memory structures

### `IstanbulSealer`

`IstanbulSealer` stores:

- the local private key used to create proposer and committed seals
- an LRU cache of parsed `IstanbulExtra`

### `backend.backend`

The backend keeps the node-facing runtime state:

- `istanbulEventMux`: shared event bus between worker, backend, and core
- `commitCh`, `proposedBlockHash`, `sealSkippedNum`: coordination state for the blocking `Seal()` path
- `recentMessages`, `knownMessages`: duplicate-suppression caches
- `currentView`: latest view published by the core
- `candidates`: local validator authorization proposals exposed through RPC
- references to kaiax modules, broadcaster, chain accessors, and executor

### `core.core`

The core keeps the consensus-facing runtime state:

- current `roundState`
- pending requests and backlog queues for future messages
- round-change tracking and timers
- committee, quorum, and proposer metadata for the active view
- observability metrics and consensus timestamps

## Lifecycle

### Construction

`backend.New(opts)` creates:

- the local sealer and node address
- the event mux
- the message caches
- the Istanbul core instance

### Start and stop

`backend.Start(...)`:

- injects chain access, current-block callback, bad-block callback, and executor
- resets transient seal state
- starts the core event loop

`backend.Stop()`:

- closes any pending `commitCh`
- stops the core event loop

`core.Start()` subscribes to the event mux, starts from the latest sequence plus one, and launches the handler goroutine.

`core.Stop()` unsubscribes from all event streams, stops timers, and waits for the handler goroutine to exit.

## Block processing

### Consensus

#### SubmitTransactions

`backend.SubmitTransactions(...)` bridges execution and consensus.

Its flow is:

1. Read the qualified validators for the new block from `kaiax/valset`.
2. Write validators into the header using the sealer.
3. Reset the executor with the current state and header.
4. Execute transactions.
5. Finalize block state.
6. Notify the caller that the block is prepared.
7. Call `Seal()` to run Istanbul consensus.

If the local node is not the winning proposer, `Seal()` may return `nil` and the block will instead arrive through normal block import.

#### Seal

`backend.Seal(...)`:

- writes the proposer seal into the block header
- initializes the seal coordination state
- posts a `RequestEvent` into the event mux
- waits for `backend.Commit(...)` to deliver the finalized block through `commitCh`

#### Verify

`backend.Verify(...)` validates:

- block proposal type
- bad-block blacklist state
- transaction root and blob sidecar consistency
- header validity through `chain.ValidateHeader`

#### Commit

When the core gathers quorum `COMMIT` messages:

- it collects committed seals
- it calls `backend.Commit(proposal, seals)`
- the backend writes round and committed seals into the header
- if the local node proposed the block, the sealed block is sent back to `Seal()`
- otherwise the block is enqueued for normal import

#### New sequence

When a new sequence starts, the core posts `NewSequenceEvent` so the worker can prepare the next block.

## APIs

The backend exposes the public `istanbul` RPC namespace.

### `istanbul_getValidators`

Returns the qualified validators for the given block number.

- Parameters
  - `number`: block number, `latest`, or `pending`
- Returns
  - `[]common.Address`

### `istanbul_getDemotedValidators`

Returns the demoted validators for the given block number.

- Parameters
  - `number`: block number, `latest`, or `pending`
- Returns
  - `[]common.Address`

### `istanbul_getValidatorsAtHash`

Returns the qualified validators for the given block hash.

- Parameters
  - `hash`: block hash
- Returns
  - `[]common.Address`

### `istanbul_getDemotedValidatorsAtHash`

Returns the demoted validators for the given block hash.

- Parameters
  - `hash`: block hash
- Returns
  - `[]common.Address`

### `istanbul_candidates`

Returns the local candidate map currently tracked by the backend.

- Parameters
  - none
- Returns
  - `map[common.Address]bool`

### `istanbul_propose`

Adds or updates a local authorization candidate.

- Parameters
  - `address`: candidate validator address
  - `auth`: whether the node votes to authorize or deauthorize the address
- Returns
  - none

### `istanbul_discard`

Removes a local authorization candidate.

- Parameters
  - `address`: candidate validator address
- Returns
  - none

### `istanbul_getTimeout`

Returns the default Istanbul timeout value exposed by the API implementation.

- Parameters
  - none
- Returns
  - `uint64`
