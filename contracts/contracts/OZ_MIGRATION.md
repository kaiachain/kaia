# OZ v4/v5 Migration Issue

## Goal
Introduce permissionless contracts without changing existing contract bytecode/imports.

## Core Conflict
1. Permissionless contracts use OZ v5 upgradeable (`Initializable`, `OwnableUpgradeable`, `ReentrancyGuardTransient`, etc.)
2. OZ v5 upgradeable requires peer dependency `@openzeppelin/contracts@v5`
3. `@openzeppelin/contracts` must remain v4 (existing contract bytecode preservation)
4. npm does not recognize aliased packages (`openzeppelin-contracts-5.0`) as peer dependency resolution

## Additional Constraint
Existing contracts (kaiabridge, KIP113) directly import `@openzeppelin/contracts-upgradeable` v4 → cannot upgrade to v5.

## Conclusion
No clean way to install OZ v5 upgradeable via `npm install` while keeping `@openzeppelin/contracts` v4 and `@openzeppelin/contracts-upgradeable` v4.

## Approaches Tried
1. **Solc version split (0.8.28)**: OZ v5 node_modules have `^0.8.20` pragma, still compiled by 0.8.25 job → mcopy error
2. **Hardhat overrides**: OZ v5 files in node_modules not covered by per-file overrides → mcopy error
3. **@openzeppelin/contracts v4 + aliased v5 path**: `openzeppelin-contracts-upgradeable-5.0` peer requires `@openzeppelin/contracts@v5` → npm conflict
4. **npm overrides**: npm does not apply overrides to aliased package peer dependency resolution
5. **Reverse aliasing (v5 default, v4 alias)**: existing contracts import `@openzeppelin/contracts` (v4 API) → breaks if changed to v5
