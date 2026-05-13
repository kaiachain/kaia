# Contracts

## Install dependencies

- Install node.js 20+
- Install solidity compiler (for Go binding regeneration):
  - [solc-select](https://github.com/crytic/solc-select)
  - or [svm-rs](https://github.com/alloy-rs/svm-rs)

## Generate Go bindings via `abigenw`

Regenerate after upstream system contracts change.

```
cd contracts
go generate
```

Bindings are generated from [kaia-system-contracts](./kaia-system-contracts/) submodule artifacts.
After updating the submodule, run `go generate` to regenerate `bindings/`.

## Unit test

```
cd contracts
npm install
npx hardhat compile
npx hardhat test

# Run a subset
npx hardhat test --grep Bridge
npx hardhat coverage --testfiles test/Bridge/bridge.test.ts
```

## Directory structure

```
contracts/
├── abigenw                   # Go binding generator script
├── generate.go               # go generate directives
├── bindings/                 # Generated Go bindings (do not edit manually)
├── kaia-system-contracts/    # Git submodule — single source of truth for system contract .sol
├── libs/                     # Vendored Solidity libraries
├── service_chain/            # Service chain bridge contracts
└── testing/                  # Testing and mock contracts
```

### `bindings/`

Generated Go wrappers for all system contracts. Regenerated via `go generate`.

### `kaia-system-contracts/`

Git submodule pointing to [kaia-system-contract](https://github.com/kaiachain/kaia-system-contract).
Contains compiled artifacts used by `go generate`.

### `libs/`

Vendored Solidity dependencies.

- `kip13/`: ERC-165 supportsInterface.
- `openzeppelin-contracts-v2/`: OpenZeppelin contracts v2 (for legacy 0.4.x contracts).
- `uniswap/`: Uniswap V2 Go bindings (generated from `node_modules/@uniswap/`).

### `service_chain/`

Service chain token bridge contracts.

- `bridge/`: Token bridge implementation (deployed by `subbridge_deployBridge` API).
- `sc_erc20/`, `sc_erc721/`, `sc_erc721_no_uri/`: Bridge receiver interfaces.

### `testing/`

Mock and testing contracts used in Go unit tests and hardhat tests. Not deployed on mainnet.

- `system_contracts/`: Mock versions of system contracts (e.g. `RegistryMock`, `KIP113Mock`).
- `reward/`: Mock AddressBook and reward contracts for genesis testing.
- `sc_erc20/`, `sc_erc721/`, `extbridge/`: Service chain testing helpers.
