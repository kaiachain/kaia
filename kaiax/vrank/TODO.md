# Permissionless TODO

## Go

| # | File | Line | Description | Priority |
|---|------|------|-------------|----------|
| 1 | `kaiax/valset/impl/state_transition.go` | 102 | `isPassVrankTest()` — replace with KIP-227 CFS implementation (currently hardcoded `return true`). Blocked on cfsThreshold contract addition. | High |

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
| 8 | `consensus/istanbul/backend/backend.go` | Remove `SetSkipPropose` and `skipPropose` flag |
| 9 | `consensus/istanbul/backend/engine.go` | Remove `skipPropose` check in Seal |
| 10 | `node/cn/api_admin_chain.go` | Remove `SetSkipPropose` API |
| 11 | `console/web3ext/web3ext.go` | Remove `admin.setSkipPropose` web3 method |
