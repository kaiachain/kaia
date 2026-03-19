# PR Split Plan

Current working PR: kaiachain/kaia#748 (state-transition → permissionless)

## Dependency Order

```
PR 1 (Contracts) → PR 2 (Go Bindings) → PR 3 (dev sync) → PR 4 (ABv2 + State Transition) → PR 5 (Genesis Alloc + Integration)
```

All PRs target `permissionless` branch. Must be merged in order.

---

## PR 1: Solidity Contracts — `contracts_permissionless/`

**PR**: kaiachain/kaia#792
**Branch**: `pr1-contracts` → `permissionless`
**Status**: Merged
**Scope:** Separate permissionless contracts into independent `contracts_permissionless/` project. No Go binding, no `contracts/` changes.

| Category | Files |
|----------|-------|
| New project | `contracts_permissionless/` — sol, TS tests, hardhat config, `generate.go` |
| ABv2 | `AddressBookV2/`, `ABv2DataContract`, `AddressBookLegacy`, interfaces |
| CnStakingV4 | `CnStaking/CnStakingV4/`, `CnStakingV4Factory/` |
| PublicDelegation | `PublicDelegation/` |
| Proxy (OZ v5) | `Proxy/ERC1967Proxy.sol`, `UpgradeableBeacon.sol` |
| Libraries | `ABv2ConfigLib`, `NodeVerifier`, `SlotMath` |
| MultiCall | `multicall/MultiCallContract.sol` |
| System interfaces | `system/IRegistry.sol`, `IStakingTracker.sol`, `SystemCallable.sol` |
| Testing mocks | `testing/AddressBookMock.sol`, `AddressBookV2Mock.sol`, `MultiCallContractMock.sol` |
| Types | `types/Node.sol` |
| KIP interfaces | `kip113/IAddressBook.sol`, `kip149/IRegistry.sol`, `Registry.sol` |

**Review focus:** OZ v5 contract logic, ABv2 interface, new contract structure

---

## PR 2: Go Bindings

**PR**: kaiachain/kaia#796
**Branch**: `pr2-go-system` → `permissionless`
**Status**: Merged
**Scope:** `go generate` bindings only

| Category | Files |
|----------|-------|
| Go bindings (`go generate`) | `contracts_permissionless/contracts/**/*.go` (8 files) |

**Review focus:** Auto-generated code (~77K lines). Review focus is on the `generate.go` directives in PR1.

---

## PR 3: dev branch sync + refactoring

**Branch**: `pr3-dev-sync` → `permissionless`
**Status**: Not started
**Scope:** Merge upstream `dev` (up to `3347c86a5`) + adapt permissionless code to dev changes

| Category | Changes |
|----------|---------|
| dev merge | Merge `mainstream/dev` into `permissionless` — includes PR #768 (flex reward), PR #790 (fork override fix), Osaka mainnet fork number, etc. |
| MultiCall SpareAddress | `contracts_permissionless/multicall/MultiCallContract.sol` + regenerated binding — adds `SpareAddress` (KPF) to `multiCallStakingInfo` return |
| ConsolidatedNodes refactor | `kaiax/staking/staking_info.go` — revert to original `consolidateNodes()` (no council param), add `consolidateNodesPermissionless()`, `ConsolidatedNodes(rules)` signature |
| Staking getter | `kaiax/staking/impl/getter.go` — two multicall imports (contracts/ + contracts_permissionless/), `parsePermissionlessCallResult` |
| collectStakingAmounts | `kaiax/valset/impl/getter_demote.go` — iterate `cn.NodeIds` (plural, original) |
| Tests | `kaiax/staking/staking_info_test.go` — updated for `NodeIds`/`StakingContracts` struct |

**Review focus:** dev merge conflict resolution, ConsolidatedNodes cache safety, multicall SpareAddress integration

---

## PR 4: ABv2 Go Layer + State Transition (Core Logic)

**Branch**: `pr4-state-transition` → `permissionless`
**Status**: Not started
**Scope:** ABv2 read/write functions + validator state machine

| Category | Files |
|----------|-------|
| ABv2 functions | `blockchain/system/addressbook_v2.go` |
| Constants | `blockchain/system/constant.go` (ERC1967ProxyV5Code, mock addr, multicall import) |
| Constants test | `blockchain/system/constant_test.go` |
| KIP113 changes | `blockchain/system/kip113.go`, `kip113_test.go` |
| Multicall | `blockchain/system/multicall.go`, `multicall_test.go` |
| Contract backend | `accounts/abi/bind/backends/blockchain_state.go` |
| State transition | `kaiax/valset/impl/state_transition.go` |
| Transition getters | `kaiax/valset/impl/state_transition_getter.go` |
| Valset getter | `kaiax/valset/impl/getter.go` |
| Types | `kaiax/valset/types.go`, `kaiax/valset/address_set.go` |
| Interface | `kaiax/valset/interface.go`, `kaiax/valset/mock/module.go` |
| Staking (permissionless) | `kaiax/staking/staking_info.go`, `p2p_staking_info.go` |
| Staking getter | `kaiax/staking/impl/getter.go`, `kaiax/staking/interface.go`, `kaiax/staking/mock/` |
| Engine hook | `consensus/istanbul/backend/engine.go`, `consensus/faker/faker.go`, `consensus/mocks/engine_mock.go` |
| Worker | `work/worker.go` |
| Unit tests | `blockchain/system/addressbook_v2_test.go`, `multicall_test.go` |
| Unit tests | `kaiax/valset/impl/state_transition_test.go`, `state_transition_getter_test.go`, `execution_test.go` |
| Unit tests | `kaiax/valset/impl/getter_council_test.go`, `getter_context_test.go`, `getter_demote_test.go`, `getter_proposers_test.go` |
| Unit tests | `kaiax/staking/staking_info_test.go` |

**Review focus:** ABv2 read/write helpers, ERC1967Proxy v4/v5 split, state machine transitions, epoch/timeout/violation logic, getOrComputeNodeStates caching

---

## PR 5: Genesis Alloc + Integration

**Branch**: `pr5-integration` → `permissionless`
**Status**: Not started
**Scope:** Genesis allocation, config wiring, API, cleanup

| Category | Files |
|----------|-------|
| Genesis alloc | `blockchain/system/permissionless.go` |
| Test helpers | `blockchain/system/testing.go` |
| Chain config | `params/config.go` (PermissionlessCompatibleBlock) |
| Genesis CLI | `cmd/homi/setup/cmd.go` |
| Admin API | `node/cn/api_admin_chain.go`, `console/web3ext/web3ext.go` |
| P2P | `node/cn/peer_set.go`, `node/cn/handler.go`, `node/cn/handler_msg_test.go` |
| Valset init/schema | `kaiax/valset/impl/init.go`, `init_test.go`, `schema.go`, `schema_test.go`, `api.go`, `error.go` |
| Valset execution/council | `kaiax/valset/impl/execution.go`, `getter_council.go` |
| Staking init/api | `kaiax/staking/impl/api.go`, `init.go` |
| Reward/Randao | `kaiax/reward/impl/getter.go`, `kaiax/randao/impl/getter.go` |
| Misc tests | `tests/randao_fork_test.go`, `tests/block_test.go`, `tests/block_test_util.go`, `tests/state_test.go` |
| Backend | `consensus/istanbul/backend/backend.go`, `api/api_eth_test.go` |
| Unit tests | `blockchain/system/permissionless_test.go` |
| Docs | `permissionless/TEST_PLAN.md`, `TODO.md` |

**Review focus:** Genesis allocation logic, config wiring, API correctness, P2P changes

---

## Notes

- `contracts/` is untouched — existing bytecodes and Go bindings remain identical to base
- `contracts_permissionless/` is the new independent project for permissionless contracts (OZ v5, solc 0.8.25)
- Each PR should compile and pass tests independently
- PR 2 is largest by line count (Go bindings ~77K) but mostly generated
- PR 4 is the most logic-dense and needs the most careful review
- dev merge (PR 3) syncs up to `3347c86a5` (PR #790)