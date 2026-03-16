# Permissionless TODO

## Go

| # | File | Line | Description | Priority |
|---|------|------|-------------|----------|
| 1 | `kaiax/valset/impl/state_transition.go` | 242 | `isPassVrankTest()` — replace with KIP-227 VRank implementation (currently hardcoded `return true`) | High |
| 2 | `kaiax/valset/impl/state_transition_getter.go` | 43 | Remove devLog debug logs (`TODO-Permissionless: Remove this log`) | Medium — remove before merge |
| 3 | `kaiax/valset/impl/state_transition_getter.go` | 180 | Implement `rule2: RC` (inside `getViolationTransition`) | High |

## Config

| # | File | Line | Description | Priority |
|---|------|------|-------------|----------|
| 4 | `params/config.go` | 65 | Set mainnet `PermissionlessCompatibleBlock` | Deploy-time |
| 5 | `params/config.go` | 117 | Set Kairos `PermissionlessCompatibleBlock` | Deploy-time |

## Cleanup (Remove before production)

| # | File | Description |
|---|------|-------------|
| 6 | `node/cn/api_admin_chain.go` | Remove `SetPermissionlessForkBlock` API |
| 7 | `console/web3ext/web3ext.go` | Remove `admin.setPermissionlessForkBlock` web3 method |

## Refactoring

| # | Files | Description | Priority |
|---|-------|-------------|----------|
| 8 | `kaiax/valset/types.go`, `kaiax/valset/impl/state_transition.go`, `kaiax/valset/impl/getter.go` | Remove `CommonAddressSet` interface — NO-OP methods (`Add`/`Remove`), semantic mismatch (`Council()`), interface pollution (`EqualState`). Each path should return `[]common.Address`/`NodeStateMap` directly, eliminating the `ValidatorList` wrapper. Touches permissioned path — separate task. | Low |
