# kaiax modules

A kaiax module is a collection of modfications to the base blockchain that implement a related set of features.

## Module interface

A module must have following functions:

- `NewAbcModule()`: to declare modules first, to be interlinked by `Init` functions later.
- `Init(...)`: to link kaiax modules and components (BlockChain, ChainDB, etc.) together.

Having both Constructor and Init enables circular dependencies. e.g.

```go
mValset = NewValsetModule()
mGov = NewGovModule()

mValset.Init(InitOpts{... mGov ...})
mGov.Init(InitOpts{... mValset ...})
```

A module also must implement the `BaseModule` interface.

- `Start()`: starts any background goroutines after every module's `Init` is complete. Used in two places:
  - At `CN.Start()` after other components are started
  - At `CNAPIBackend.SetHead()` where modules are restarted after a chain rewind.
- `Stop()`: stops any background goroutines. Similar to `Start()`.

A module may implement other `Module` interfaces such as `JsonRpcModule` and `ConsensusModule`.

A non-kaiax component (such as BlockChain, CN, consensus.Engine) may implement `ModuleHost` interface to interact with kaiax modules.

```go
blockChain.RegisterExecutionModule(mGov)
```

