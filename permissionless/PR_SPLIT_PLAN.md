# PR Split Plan

Current working PR: kaiachain/kaia#748 (state-transition → permissionless)

## Dependency Order

```
PR 1 (Contracts) → PR 2 (Go Bindings) → PR 3 (dev sync) → PR 4 (Validator Lifecycle)
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

## PR 3: dev branch sync

**Branch**: `pr3-dev-sync` → `permissionless`
**Status**: Not started
**Scope:** Merge upstream `dev` (up to `3347c86a5`) only. No permissionless-specific refactoring.

| Category | Changes |
|----------|---------|
| dev merge | Merge `mainstream/dev` — includes PR #768 (flex reward: KPF, 4-part ratio, SpareAddress), PR #790 (fork override fix), Osaka mainnet/Kairos fork numbers |
| MultiCall SpareAddress | `contracts_permissionless/multicall/MultiCallContract.sol` + regenerated binding — adds `SpareAddress` (KPF) to `multiCallStakingInfo` return |

**Review focus:** dev merge conflict resolution, MultiCall SpareAddress integration

---

## PR 4: Validator Lifecycle (ABv2 + State Transition + Genesis + Integration)

**Branch**: `pr4-validator-lifecycle` → `permissionless`
**Status**: In progress
**Scope:** All Go-side permissionless logic — ABv2, state transition, genesis alloc, integration

| Category | Files |
|----------|-------|
| ABv2 functions | `blockchain/system/addressbook_v2.go`, `addressbook_v2_test.go` |
| Genesis alloc | `blockchain/system/permissionless.go`, `permissionless_test.go` |
| Test helpers | `blockchain/system/testing.go` |
| Constants | `blockchain/system/constant.go` (ERC1967ProxyV5Code, mock addr) |
| KIP113 changes | `blockchain/system/kip113.go`, `kip113_test.go` |
| Multicall | `blockchain/system/multicall.go`, `multicall_test.go` |
| MultiCall SpareAddress | `contracts_permissionless/multicall/MultiCallContract.sol` + binding |
| ConsolidatedNodes refactor | `kaiax/staking/staking_info.go` — `ConsolidatedNodes(rules)`, `consolidateNodesPermissionless()` |
| Staking | `kaiax/staking/impl/getter.go`, `interface.go`, `mock/`, `p2p_staking_info.go`, `api.go` |
| State transition | `kaiax/valset/impl/state_transition.go`, `state_transition_getter.go` |
| Valset | `kaiax/valset/types.go`, `address_set.go`, `interface.go`, `mock/module.go` |
| Valset impl | `kaiax/valset/impl/getter.go`, `getter_council.go`, `getter_demote.go`, `execution.go`, `init.go`, `schema.go`, `api.go`, `error.go` |
| Engine hook | `consensus/istanbul/backend/engine.go` |
| Chain config | `params/config.go` (PermissionlessCompatibleBlock) |
| Genesis CLI | `cmd/homi/setup/cmd.go` |
| Admin API | `node/cn/api_admin_chain.go`, `console/web3ext/web3ext.go` |
| P2P | `node/cn/peer_set.go`, `node/cn/handler.go`, `node/cn/handler_msg_test.go` |
| Reward/Randao | `kaiax/reward/impl/getter.go`, `kaiax/randao/impl/getter.go` |
| Misc tests | `tests/randao_fork_test.go` |
| Unit tests | All `*_test.go` in above packages |

**Review focus:** ABv2 read/write, ERC1967Proxy v4/v5 split, ConsolidatedNodes cache safety, state machine transitions, genesis allocation, config wiring

---

## Notes

- `contracts/` is untouched — existing bytecodes and Go bindings remain identical to base
- `contracts_permissionless/` is the new independent project for permissionless contracts (OZ v5, solc 0.8.25)
- PR 2 is largest by line count (Go bindings ~77K) but mostly generated
- PR 3 is pure dev sync — no permissionless-specific code changes
- PR 4 is the most logic-dense — combines ABv2, state transition, genesis alloc, and integration
- dev merge (PR 3) syncs up to `3347c86a5` (PR #790)