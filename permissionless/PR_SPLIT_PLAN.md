# PR Split Plan

Current working PR: kaiachain/kaia#748 (state-transition → permissionless)
- 171 files, ~87K lines changed — too large for a single review

## Dependency Order

```
PR 1 (Contracts) → PR 2 (ABv2 Go) → PR 3 (State Transition) → PR 4 (Integration)
```

Each PR targets `permissionless` branch. Must be merged in order.

---

## PR 1: Solidity Contracts (OZ v5 Migration + New Contracts)

**Scope:** All Solidity, TypeScript tests, Go bindings, build config

| Category | Files |
|----------|-------|
| ABv2 contracts | `contracts/.../AddressBookV2/*.sol`, `interfaces/`, `mocks/` |
| ABv2 Go bindings | `contracts/.../AddressBookV2/AddressBookV2.go`, `abv2data/ABv2DataContract.go` |
| CnStakingV4 | `contracts/.../CnStakingV4/`, `CnStakingV4Factory/` |
| PublicDelegation | `contracts/.../PublicDelegation/` |
| Proxy (OZ v5) | `contracts/.../Proxy/` (new), delete `proxy/` (old) |
| OZ v5 migration | All `.sol` changes across `auction/`, `consensus/`, `kaiabridge/`, `kip113/`, `kip247/`, `gov/`, `libraries/`, `multicall/` |
| Go bindings | All `.go` under `contracts/` (regenerated) |
| Build config | `contracts/package.json`, `package-lock.json`, `hardhat.config.ts`, `generate.go`, `abigenw` |
| TS tests | `contracts/test/**/*.ts`, `contracts/test/materials/fixtures/*.ts` |
| Constant update | `blockchain/system/constant.go` (bytecode hashes) |
| Test fix | `blockchain/system/constant_test.go`, `common/compiler/solidity_test.go` |

**Review focus:** OZ v4→v5 migration correctness, ABv2 contract logic, new contract interfaces

---

## PR 2: ABv2 Go Layer (Read/Write/Install)

**Scope:** Go code that interacts with ABv2 contract

| Category | Files |
|----------|-------|
| ABv2 functions | `blockchain/system/addressbook_v2.go` |
| Genesis alloc | `blockchain/system/permissionless.go` |
| Test helpers | `blockchain/system/testing.go` |
| KIP113 changes | `blockchain/system/kip113.go`, `kip113_test.go` (Proxy import) |
| Multicall | `blockchain/system/multicall_test.go` |
| Unit tests | `blockchain/system/addressbook_v2_test.go`, `permissionless_test.go` |

**Review focus:** ABI encoding/decoding, contract read helpers (ReadGetAllValidators, ReadABv2Timeouts, etc.), genesis allocation, InstallAddressBookV2

---

## PR 3: Permissionless State Transition (Core Logic)

**Scope:** Validator state machine, epoch/timeout/violation transitions

| Category | Files |
|----------|-------|
| State transition | `kaiax/valset/impl/state_transition.go` |
| Transition getters | `kaiax/valset/impl/state_transition_getter.go` |
| Valset getter | `kaiax/valset/impl/getter.go` |
| Types | `kaiax/valset/types.go`, `kaiax/valset/address_set.go` |
| Interface | `kaiax/valset/interface.go`, `kaiax/valset/mock/module.go` |
| Staking (permissionless path) | `kaiax/staking/staking_info.go` (ConsolidatedNodes), `kaiax/staking/p2p_staking_info.go` |
| Staking getter | `kaiax/staking/impl/getter.go`, `kaiax/staking/interface.go`, `kaiax/staking/mock/` |
| Engine hook | `consensus/consensus.go`, `consensus/istanbul/backend/engine.go`, `consensus/faker/faker.go`, `consensus/mocks/engine_mock.go` |
| State processor | `blockchain/state_processor.go`, `blockchain/state_transition.go` |
| Worker | `work/worker.go` |
| Unit tests | `kaiax/valset/impl/state_transition_test.go`, `state_transition_getter_test.go`, `execution_test.go` |
| Unit tests | `kaiax/valset/impl/getter_council_test.go`, `getter_context_test.go`, `getter_demote_test.go`, `getter_proposers_test.go` |
| Unit tests | `kaiax/staking/staking_info_test.go` |

**Review focus:** State machine (ValActive/ValReady/ValPaused/ValInactive/ValExiting transitions), epoch transition, timeout transition, violation transition, getOrComputeNodeStates caching

---

## PR 4: Integration (Config, API, Misc)

**Scope:** Remaining wiring, config, API, cleanup

| Category | Files |
|----------|-------|
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
| Docs | `permissionless/TEST_PLAN.md`, `TODO.md` |

**Review focus:** Config wiring, API correctness, P2P changes, integration test scenarios

---

## Notes

- Each PR should compile and pass tests independently (may need stub/no-op where downstream code isn't merged yet)
- PR 1 is the largest by line count (~70K) but mostly generated Go bindings — actual review is Solidity + build config
- PR 3 is the most logic-dense and needs the most careful review
- PR 4 contains cleanup items listed in `permissionless/TODO.md`