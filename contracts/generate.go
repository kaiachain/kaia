// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package contracts

/*
Recommended to install solc-select or svm-rs to switch solc versions within abigenw.

	solc-select install 0.4.24 0.5.6 0.5.9 0.5.16 0.6.6 0.8.19 0.8.25
	go generate

Othewise, you can manually switch solc versions and run go generate for each solc version.

	go generate --run 0.4.24
	go generate --run 0.5.6
	go generate --run 0.5.9
	go generate --run 0.5.16
	go generate --run 0.6.6
	go generate --run 0.8.19
	go generate --run 0.8.25
*/

// These files were compiled with solidity 0.4.24.

// WARNING: Do not regenerate credit.go — the BinRuntime must match the mainnet-deployed bytecode.
// //go:generate ./abigenw --pkg misc --sol ./contracts/system_contracts/misc/credit.sol --out ./bindings/misc/credit.go --ver 0.4.24
//go:generate ./abigenw --pkg reward --sol ./testing/reward/all.sol --out ./bindings/testing/reward/all.go --ver 0.4.24

// These files were compiled with solidity 0.5.6.

//go:generate ./abigenw --pkg bridge --sol ./service_chain/bridge/Bridge.sol --out ./bindings/bridge/bridge.go --ver 0.5.6
//go:generate ./abigenw --pkg sc_erc20 --sol ./testing/sc_erc20/sc_token.sol --out ./bindings/testing/sc_erc20/sc_token.go --ver 0.5.6
//go:generate ./abigenw --pkg sc_erc721 --sol ./testing/sc_erc721/sc_nft.sol --out ./bindings/testing/sc_erc721/sc_nft.go --ver 0.5.6
//go:generate ./abigenw --pkg sc_erc721_no_uri --sol ./testing/sc_erc721_no_uri/sc_nft_no_uri.sol --out ./bindings/testing/sc_erc721_no_uri/sc_nft_no_uri.go --ver 0.5.6
//go:generate ./abigenw --pkg extbridge --sol ./testing/extbridge/ext_bridge.sol --out ./bindings/testing/extbridge/ext_bridge.go --ver 0.5.6

// These files were compiled with solidity 0.5.9.

//go:generate ./abigenw --pkg system_contracts --sol ./testing/system_contracts/WKAIA.sol --out ./bindings/testing/system_contracts/WKAIA.go --ver 0.5.9

// These files were compiled with solidity 0.5.16.

//go:generate ./abigenw --pkg factory --sol ./node_modules/@uniswap/v2-core/contracts/UniswapV2Factory.sol --out ./bindings/uniswap/factory/UniswapV2Factory.go --ver 0.5.16

// These files were compiled with solidity 0.6.6.

//go:generate ./abigenw --pkg router --sol ./node_modules/@uniswap/v2-periphery/contracts/UniswapV2Router02.sol --out ./bindings/uniswap/router/UniswapV2Router02.go --ver 0.6.6

// These files were compiled with solidity 0.8.19.

//go:generate ./abigenw --pkg gov --type GovParam --artifact ./kaia-system-contract/contracts-klaytn-v1.10/artifacts/contracts/GovParam.sol/GovParam.json --out ./bindings/gov/GovParam.go
//go:generate ./abigenw --pkg rebalance --type TreasuryRebalance --artifact ./kaia-system-contract/contracts-klaytn-v1.10/artifacts/contracts/TreasuryRebalance.sol/TreasuryRebalance.json --out ./bindings/rebalance/TreasuryRebalance.go
//go:generate ./abigenw --pkg rebalance --type TreasuryRebalanceV2 --artifact ./kaia-system-contract/contracts-v1.0/artifacts/contracts/TreasuryRebalanceV2.sol/TreasuryRebalanceV2.json --out ./bindings/rebalance/TreasuryRebalanceV2.go
// WARNING: Do not regenerate SimpleBlsRegistry.go — the BinRuntime must match the mainnet-deployed bytecode.
// //go:generate ./abigenw --pkg kip113 --type SimpleBlsRegistry --artifact ./kaia-system-contract/contracts-klaytn-v1.12/artifacts/contracts/SimpleBlsRegistry.sol/SimpleBlsRegistry.json --out ./bindings/kip113/SimpleBlsRegistry.go
// WARNING: Do not regenerate Registry.go — the BinRuntime must match the mainnet-deployed bytecode.
// //go:generate ./abigenw --pkg kip149 --type Registry --artifact ./kaia-system-contract/contracts-klaytn-v1.12/artifacts/contracts/Registry.sol/Registry.json --out ./bindings/kip149/Registry.go
// WARNING: Do not regenerate proxyv4/Proxy.go — the BinRuntime must match the mainnet-deployed bytecode.
// //go:generate ./abigenw --pkg proxyv4 --type ERC1967Proxy --sol ./contracts/system_contracts/proxy/proxy.sol --out ./bindings/proxyv4/Proxy.go --ver 0.8.19
//go:generate ./abigenw --pkg system_contracts --sol ./testing/system_contracts/all.sol --out ./bindings/testing/system_contracts/all.go --ver 0.8.19

// These files were compiled with solidity 0.8.25.

//go:generate ./abigenw --pkg kip247 --type GaslessSwapRouter --artifact ./kaia-system-contract/contracts-v2.0/artifacts/contracts/gasless/GaslessSwapRouter.sol/GaslessSwapRouter.json --out ./bindings/kip247/GaslessSwapRouter.go
//go:generate ./abigenw --pkg gasless --type TestToken --artifact ./kaia-system-contract/contracts-v2.0/artifacts/contracts/gasless/TestToken.sol/TestToken.json --out ./bindings/testing/system_contracts/gasless/TestToken.go

//go:generate ./abigenw --pkg auction --type IAuctionEntryPoint --artifact ./kaia-system-contract/contracts-v2.1/artifacts/contracts/auction/interfaces/IAuctionEntryPoint.sol/IAuctionEntryPoint.json --out ./bindings/auction/Kip249.go

// Permissionless contracts (v3.0) — generated from forge artifacts in kaia-system-contract submodule.
// MultiCallContractMock is excluded (not yet in kaia-system-contract; kept as copied binding).

//go:generate ./abigenw --pkg addressbookv2 --type AddressBookV2 --artifact ./kaia-system-contract/contracts-v3.0/out/AddressBookV2.sol/AddressBookV2.json --out ./bindings/addressbookv2/AddressBookV2.go
//go:generate ./abigenw --pkg abv2data --type ABv2DataContract --artifact ./kaia-system-contract/contracts-v3.0/out/ABv2DataContract.sol/ABv2DataContract.json --out ./bindings/abv2data/ABv2DataContract.go
//go:generate ./abigenw --pkg cnstakingv4 --type CnStakingV4 --artifact ./kaia-system-contract/contracts-v3.0/out/CnStakingV4.sol/CnStakingV4.json --out ./bindings/cnstakingv4/CnStakingV4.go
//go:generate ./abigenw --pkg cnstakingv4factory --type CnStakingV4Factory --artifact ./kaia-system-contract/contracts-v3.0/out/CnStakingV4Factory.sol/CnStakingV4Factory.json --out ./bindings/cnstakingv4factory/CnStakingV4Factory.go
//go:generate ./abigenw --pkg multicall --type MultiCallContract --artifact ./kaia-system-contract/contracts-v3.0/out/MultiCallContract.sol/MultiCallContract.json --out ./bindings/multicall/MultiCallContract.go
//go:generate ./abigenw --pkg proxyv5 --type ERC1967Proxy --artifact ./kaia-system-contract/contracts-v3.0/out/ERC1967Proxy.sol/ERC1967Proxy.json --out ./bindings/proxyv5/Proxy.go
//go:generate ./abigenw --pkg beacon --type UpgradeableBeacon --artifact ./kaia-system-contract/contracts-v3.0/out/UpgradeableBeacon.sol/UpgradeableBeacon.json --out ./bindings/beacon/UpgradeableBeacon.go
//go:generate ./abigenw --pkg publicdelegation --type PublicDelegation --artifact ./kaia-system-contract/contracts-v3.0/out/PublicDelegation.sol/PublicDelegation.json --out ./bindings/publicdelegation/PublicDelegation.go
