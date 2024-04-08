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

package contracts

/*
You can use solc-select or svm-rs to quickly switch solc versions.
Example:
	solc-select use 0.4.24
	go generate --run 0.4.24
	solc-select use 0.5.6
	go generate --run 0.5.6
	solc-select use 0.8.19
	go generate --run 0.8.19
*/

// These files were compiled with solidity 0.4.24.

//go:generate abigen --pkg misc --sol ./contracts/system_contracts/misc/credit.sol --out ./contracts/system_contracts/misc/credit.go # 0.4.24
//go:generate abigen --pkg consensus --sol ./contracts/system_contracts/consensus/consensus.sol --out ./contracts/system_contracts/consensus/consensus.go # 0.4.24
//go:generate abigen --pkg reward --sol ./contracts/testing/reward/all.sol --out ./contracts/testing/reward/all.go # 0.4.24

// These files were compiled with solidity 0.5.6.

//go:generate abigen --pkg bridge --sol ./contracts/service_chain/bridge/Bridge.sol --out ./contracts/service_chain/bridge/bridge.go # 0.5.6
//go:generate abigen --pkg sc_erc20 --sol ./contracts/testing/sc_erc20/sc_token.sol --out ./contracts/testing/sc_erc20/sc_token.go # 0.5.6
//go:generate abigen --pkg sc_erc721 --sol ./contracts/testing/sc_erc721/sc_nft.sol --out ./contracts/testing/sc_erc721/sc_nft.go # 0.5.6
//go:generate abigen --pkg sc_erc721_no_uri --sol ./contracts/testing/sc_erc721_no_uri/sc_nft_no_uri.sol --out ./contracts/testing/sc_erc721_no_uri/sc_nft_no_uri.go # 0.5.6
//go:generate abigen --pkg extbridge --sol ./contracts/testing/extbridge/ext_bridge.sol --out ./contracts/testing/extbridge/ext_bridge.go # 0.5.6

// These files were compiled with solidity 0.8.19.

//go:generate abigen --pkg gov --sol ./contracts/system_contracts/gov/GovParam.sol --out ./contracts/system_contracts/gov/GovParam.go # 0.8.19
//go:generate abigen --pkg kip103 --sol ./contracts/system_contracts/kip103/TreasuryRebalance.sol --out ./contracts/system_contracts/kip103/TreasuryRebalance.go # 0.8.19
//go:generate abigen --pkg kip113 --sol ./contracts/system_contracts/kip113/SimpleBlsRegistry.sol --out ./contracts/system_contracts/kip113/SimpleBlsRegistry.go # 0.8.19
//go:generate abigen --pkg kip149 --sol ./contracts/system_contracts/kip149/Registry.sol --out ./contracts/system_contracts/kip149/Registry.go # 0.8.19
//go:generate abigen --pkg proxy --sol ./contracts/system_contracts/proxy/proxy.sol --out ./contracts/system_contracts/proxy/proxy.go # 0.8.19
//go:generate abigen --pkg system_contracts --sol ./contracts/testing/system_contracts/all.sol --out ./contracts/testing/system_contracts/all.go # 0.8.19
