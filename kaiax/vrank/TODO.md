# Permissionless TODO

## Go

No pending items.

## Config

| # | File | Line | Description | Priority |
|---|------|------|-------------|----------|
| 4 | `params/config.go` | 65 | Set mainnet `PermissionlessCompatibleBlock` | Deploy-time |
| 5 | `params/config.go` | 117 | Set Kairos `PermissionlessCompatibleBlock` | Deploy-time |

## Cleanup (Remove before production)

| # | File | Description |
|---|------|-------------|
| 8 | `consensus/istanbul/backend/backend.go` | Remove `SetSkipPropose` and `skipPropose` flag |
| 9 | `consensus/istanbul/backend/engine.go` | Remove `skipPropose` check in Seal |
| 10 | `node/cn/api_admin_chain.go` | Remove `SetSkipPropose` API |
| 11 | `console/web3ext/web3ext.go` | Remove `admin.setSkipPropose` web3 method |
| 12 | `kaiax/vrank/impl/init.go` | Remove `SetSkipCandidate` and `skipCandidate` flag |
| 13 | `node/cn/api_admin_chain.go` | Remove `SetSkipCandidate` API |
| 14 | `console/web3ext/web3ext.go` | Remove `admin.setSkipCandidate` web3 method |
