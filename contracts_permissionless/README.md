# contracts_permissionless

Permissionless (KIP-286/287) Solidity contracts, separated from `contracts/` to avoid OZ v4/v5 dependency conflicts.

## Why Separate?

- `contracts/` uses OZ v4 (`@openzeppelin/contracts@^4.9.6`) with Hardhat 2.19 (evmVersion: paris)
- Permissionless contracts use OZ v5 (`@openzeppelin/contracts@^5.5.0`) which requires evmVersion: cancun (tload/tstore/mcopy)
- Mixing both in one Hardhat project changes existing contract bytecode — not allowed

## Contracts

| Directory | Description |
|---|---|
| `AddressBookV2/` | ABv2 core — state management, node lifecycle, UUPS proxy |
| `CnStaking/CnStakingV4/` | CnStakingV4 — validator staking with ReentrancyGuardTransient |
| `CnStaking/CnStakingV4Factory/` | Factory for deploying CnStakingV4 via beacon proxy |
| `PublicDelegation/` | Public delegation contract |
| `Proxy/` | ERC1967Proxy, UpgradeableBeacon (permissionless-specific) |
| `multicall/` | MultiCallContract with `multiCallStakingInfoPermissionless` |
| `system/` | SystemCallable, IRegistry, IStakingTracker |
| `libraries/` | ABv2ConfigLib, NodeVerifier, SlotMath |
| `types/` | Node.sol (Profile, State) |

## Testing Contracts (from `contracts/`)

| File | Source | Reason |
|---|---|---|
| `kip149/Registry.sol` | `contracts/system_contracts/kip149/` | Test helper deploys Registry at 0x401 |
| `kip113/IAddressBook.sol` | `contracts/system_contracts/kip113/` | MultiCallContractMock dependency |
| `testing/AddressBookMock.sol` | `contracts/testing/reward/` | Multicall permissioned mock (legacy) |
| `testing/AddressBookV2Mock.sol` | New | ABv2 mock for permissionless test |
| `testing/MultiCallContractMock.sol` | New | MultiCall mock with permissionless support |

## Setup

```bash
npm install
npx hardhat compile
npx hardhat test
```