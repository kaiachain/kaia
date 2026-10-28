# consensus/bft

This package is the public entry point for Kaia's BFT consensus stack.

It does not implement the consensus runtime itself. Instead, it exposes stable constructor-style APIs over the current concrete implementation in [consensus/istanbul](../istanbul/README.md).

## Overview

The rest of the codebase depends on the generic `consensus.Engine` and `consensus.Sealer` interfaces.

This package exists so callers can:

- construct the live BFT engine without importing the concrete Istanbul backend directly
- construct a consensus-aware sealer for header parsing, header writing, and hashing

In other words, `consensus/bft` is a facade package, not the place where the runtime protocol is implemented.

## Interfaces

### `consensus.Engine`

`Engine` is the live consensus-runtime interface.

It is responsible for:

- starting and stopping the active consensus runtime
- accepting the current execution result through `SubmitTransactions(...)`
- exposing consensus RPC APIs
- notifying the worker when a new sequence begins

This interface is about the current head of the chain and the node's active consensus role. A node runs one concrete engine for the current chain, so `NewEngine(...)` simply chooses the concrete runtime implementation and returns it.

If a future hardfork requires changing the live consensus runtime, this is the layer where engine-selection or engine-swap logic can be added while the rest of the codebase continues to depend on the same `consensus.Engine` interface.

#### Engine methods

- `RegisterKaiaxModules(mGov, mValset)`
  - Wires governance and validator-set modules into the live engine.
- `Start(chain, currentBlock, hasBadBlock, executor)`
  - Starts the active consensus runtime with the chain reader, head callbacks, and execution bridge.
- `Stop()`
  - Stops the active consensus runtime.
- `SubmitTransactions(txs, state, header, mux, onPrepared)`
  - Hands the current execution work to the consensus engine and returns a channel for the finalized result.
- `APIs(chain)`
  - Returns the RPC APIs exposed by the concrete engine.
- `PurgeCache()`
  - Clears engine-specific transient caches when needed.
- `SubscribeNewSequence()`
  - Subscribes to notifications that a new consensus sequence has started.

### `consensus.Sealer`

`Sealer` is the header and signature logic interface, independent from the live runtime.

It is responsible for:

- parsing consensus metadata from `header.Extra`
- writing validators, round, proposer seal, and committed seals into headers
- computing consensus-aware header hashes
- creating and verifying seal-related data

Unlike `Engine`, `Sealer` is used not only for the next block being produced, but also for arbitrary historical blocks. Blockchain code, APIs, tests, and utilities may ask questions such as:

- who authored this old block?
- who committed this block?
- what validators were encoded in this header?
- what hash should this historical header have under its consensus rules?

That is why `Sealer` has to be designed with per-block compatibility in mind.

#### Sealer methods

- `Author(header)`
  - Returns the proposer address recovered from the proposer seal in the header.
- `Committers(header)`
  - Returns the validator addresses recovered from the committed seals in the header.
- `Vanity(header)`
  - Returns the vanity bytes stored in `header.Extra`.
- `RawSeals(header)`
  - Returns the raw proposer seal and committed seals from the header.
- `Round(header)`
  - Returns the consensus round encoded in the header.
- `Validators(header)`
  - Returns the validator list encoded in the header.
- `WriteAuthorSeal(header, seal)`
  - Writes the proposer seal into the header.
- `WriteCommittedSeals(header, committedSeals)`
  - Writes the committed seals into the header.
- `WriteRound(header, round)`
  - Writes the round number into the header.
- `WriteValidators(header, validators)`
  - Writes the validator list into the header.
- `MakeAuthorSeal(header)`
  - Creates the local proposer seal bytes for the header.
- `MakeCommittedSeal(header)`
  - Creates the local committed seal bytes for the header.
- `HeaderHash(header)`
  - Returns the consensus-aware block hash for the header.
- `SigHash(header)`
  - Returns the proposer-signing hash for the header.
- `F(blockNum, qualifiedLen, committeeSize)`
  - Returns the number of faulty validators tolerated at the given block context.
- `Quorum(blockNum, qualifiedLen, committeeSize)`
  - Returns the quorum threshold required at the given block context.

## Exposed constructors

### `NewEngine(opts)`

`NewEngine(opts) consensus.Engine` constructs the live consensus engine used by CN.

It is the entry point for:

- live consensus runtime
- kaiax module registration
- p2p consensus protocol handling
- Istanbul RPC API exposure

### `NewSealer(chainConfig, privateKey)`

`NewSealer(chainConfig, privateKey) consensus.Sealer` constructs the consensus-aware header helper used by:

- blockchain code
- tests
- genesis helpers
- CLI tooling

This sealer is responsible for:

- reading validators, author, committers, and round from headers
- writing validators and seals into `header.Extra`
- computing consensus-aware header hashes
- creating proposer and committed signatures

## Current implementation

Although the package name is generic, the current implementation is Istanbul-specific.

- `NewEngine` returns `consensus/istanbul/backend`
- `NewSealer` is backed by `istanbul.NewSealerImpl`

This arrangement keeps the concrete implementation behind a narrow public surface. If Kaia later changes the concrete BFT engine or introduces fork-dependent selection logic, the rest of the codebase does not need to change its import surface.

## Important behavior

### Node type forwarding

`NewEngine(opts)` passes `opts.NodeType` into the concrete engine.

This matters because live consensus behavior depends on node role.

- `CONSENSUSNODE`
  - actively participates in validator-side runtime behavior such as timeout handling and round-change sending
- other node types
  - may host the engine and consume consensus data, but do not drive the live consensus flow in the same way

### Header hash override

`NewSealer(...)` installs `types.SetHeaderHashFn(s.headerHash)`.

As a result, the following use the consensus-aware hash instead of the default plain RLP header hash:

- `types.Header.Hash()`
- `types.Block.Hash()`

### Height dispatch point

`dynamicSealer` has an `implAt(number)` dispatch point.

Today it always returns the Istanbul sealer implementation, but the indirection exists for an important reason.

`Engine` only needs to represent the currently active runtime. `Sealer`, on the other hand, may be asked to interpret or hash headers from many different block heights.

This matters because consensus rules can be height-dependent:

- pre-Istanbul blocks may use plain RLP header hashing
- Istanbul blocks use Istanbul-aware `Extra` parsing and header hashing
- future forks could introduce new header formats or seal rules while old blocks must remain readable

So the dispatch point is there to let `Sealer` choose the correct implementation for the block number being queried, even though today's implementation still routes every Istanbul-era block to the same sealer.

## What this package does not do

This package does not:

- implement the consensus protocol itself
- expose its own RPC namespace
- define its own persistent schema
- manage its own long-lived runtime beyond constructing the concrete engine

For runtime details, see [consensus/istanbul](../istanbul/README.md).
