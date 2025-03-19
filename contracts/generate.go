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

	solc-select install 0.4.24 0.5.6 0.8.19 0.8.25
	go generate

Othewise, you can manually switch solc versions and run go generate for each solc version.

	go generate --run 0.4.24
	go generate --run 0.5.6
	go generate --run 0.8.19
	go generate --run 0.8.25
*/

// These files were compiled with solidity 0.4.24.

//go:generate ./abigenw --pkg misc --sol ./contracts/system_contracts/misc/credit.sol --out ./contracts/system_contracts/misc/credit.go --ver 0.4.24
//go:generate ./abigenw --pkg consensus --sol ./contracts/system_contracts/consensus/consensus.sol --out ./contracts/system_contracts/consensus/consensus.go --ver 0.4.24
//go:generate ./abigenw --pkg reward --sol ./contracts/testing/reward/all.sol --out ./contracts/testing/reward/all.go --ver 0.4.24

// These files were compiled with solidity 0.5.6.

//go:generate ./abigenw --pkg bridge --sol ./contracts/service_chain/bridge/Bridge.sol --out ./contracts/service_chain/bridge/bridge.go --ver 0.5.6
//go:generate ./abigenw --pkg sc_erc20 --sol ./contracts/testing/sc_erc20/sc_token.sol --out ./contracts/testing/sc_erc20/sc_token.go --ver 0.5.6
//go:generate ./abigenw --pkg sc_erc721 --sol ./contracts/testing/sc_erc721/sc_nft.sol --out ./contracts/testing/sc_erc721/sc_nft.go --ver 0.5.6
//go:generate ./abigenw --pkg sc_erc721_no_uri --sol ./contracts/testing/sc_erc721_no_uri/sc_nft_no_uri.sol --out ./contracts/testing/sc_erc721_no_uri/sc_nft_no_uri.go --ver 0.5.6
//go:generate ./abigenw --pkg extbridge --sol ./contracts/testing/extbridge/ext_bridge.sol --out ./contracts/testing/extbridge/ext_bridge.go --ver 0.5.6

// These files were compiled with solidity 0.8.19.

//go:generate ./abigenw --pkg gov --sol ./contracts/system_contracts/gov/GovParam.sol --out ./contracts/system_contracts/gov/GovParam.go --ver 0.8.19
//go:generate ./abigenw --pkg rebalance --sol ./contracts/system_contracts/rebalance/all.sol --out ./contracts/system_contracts/rebalance/all.go --ver 0.8.19
//go:generate ./abigenw --pkg kip113 --sol ./contracts/system_contracts/kip113/SimpleBlsRegistry.sol --out ./contracts/system_contracts/kip113/SimpleBlsRegistry.go --ver 0.8.19
//go:generate ./abigenw --pkg kip149 --sol ./contracts/system_contracts/kip149/Registry.sol --out ./contracts/system_contracts/kip149/Registry.go --ver 0.8.19
//go:generate ./abigenw --pkg proxy --sol ./contracts/system_contracts/proxy/proxy.sol --out ./contracts/system_contracts/proxy/proxy.go --ver 0.8.19
//go:generate ./abigenw --pkg system_contracts --sol ./contracts/testing/system_contracts/all.sol --out ./contracts/testing/system_contracts/all.go --ver 0.8.19
//go:generate ./abigenw --pkg multicall --sol ./contracts/system_contracts/multicall/MultiCallContract.sol --out ./contracts/system_contracts/multicall/MultiCallContract.go --ver 0.8.19

// These files were compiled with solidity 0.8.25.

//go:generate ./abigenw --pkg consensus --sol ./contracts/system_contracts/consensus/Kip163.sol --out ./contracts/system_contracts/consensus/Kip163.go --ver 0.8.25
