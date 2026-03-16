# Permissionless Unit Test Plan

## Staking (`kaiax/staking/impl/`) — Skipped

Already covered by existing multicall mock tests (ABI binding + return value verification). `parsePermissionlessCallResult` is simple mapping logic with no additional test contribution.

## State Transition — Done

File: `kaiax/valset/impl/state_transition_getter_test.go`

| Test | Status |
|------|--------|
| `TestGetEpochTransition` | Done |
| `TestGetTimeoutTransition` | Done |
| `TestGetViolationTransition` | Done |
| `TestGetPermlessCouncil` | Done |
| `TestApplyAllTransitions` | Done |

## Voting — Done

File: `kaiax/valset/impl/execution_test.go`

| Test | Status |
|------|--------|
| `TestPostInsertBlock_PermissionlessIgnoresVote` | Done |

## Multicall — Done

File: `blockchain/system/multicall_test.go`

| Test | Status |
|------|--------|
| `TestContractCallerForMultiCall` | Done |
| `TestMultiCallStakingInfoPermissionless` | Done |

## Genesis — Done

File: `blockchain/system/permissionless_test.go`

| Test | Status |
|------|--------|
| `TestAllocPermissionless` | Done |
| `TestAllocPermissionless_SingleValidator` | Done |
| `TestAllocPermissionless_MismatchedLengths` | Done |

## Permissionless Activation — Done

Files: `blockchain/system/addressbook_v2_test.go`, `kaiax/valset/impl/state_transition_test.go`

| Test | File | Status | Description |
|------|------|--------|-------------|
| `TestEncodeWriteNodes` | `addressbook_v2_test.go` | Done | Validator state → ABI encoding (state uint8, timeout encoding) |
| `TestInstallAddressBookV2` | `addressbook_v2_test.go` | Done | ABv2 proxy installed at 0x400 (code + impl slot verification) |
| `TestInstallAndInitializeABv2` | `state_transition_test.go` | Done | Runtime HF scenario: `v.InstallABv2()` core path |
| `TestReadGetAllValidators` | `state_transition_test.go` | Done | Pre/post transition ReadGetAllValidators verification (`writeNodesToContract` core path) |
| `TestGetCouncilPermissionless` | `state_transition_getter_test.go` | Done | Council filtering: Active/Paused/Ready included, others excluded |

## ABv2 Contract Reading — Done

| Test | File | Status | Description |
|------|------|--------|-------------|
| `TestReadAddressBookV2BlsAll` | `addressbook_v2_test.go` | Done | Read all BLS keys from ABv2 |
| `TestGetNodeByState` | `state_transition_getter_test.go` | Done | Permissionless state filtering public API |

## Code Changes

- Removed `getCouncilGenesisPermissionless` (dead code: ABv2 is directly readable at permissionless block 0)
- `getAllStateNodes`: removed `num == 0` branch, unified to `getOrComputeNodeStates`
- `NodeStateMap.Addresses()`: added `bytes.Compare` sorting (eliminates map iteration non-determinism)
- `ValidatorList`: removed `permlessMu` (not shared across goroutines), added defensive copy in `newValidatorList`