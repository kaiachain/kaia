// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package contracts_permissionless

/*
Go binding generation for permissionless contracts.
Uses abigenw from contracts/ directory. OZ v5 imports are resolved from this directory's node_modules.

	cd contracts_permissionless
	solc-select install 0.8.19 0.8.25
	go generate
*/

//go:generate ../contracts/abigenw --pkg addressbookv2 --sol ./contracts/AddressBookV2/AddressBookV2.sol --out ../bindings/addressbookv2/AddressBookV2.go --ver 0.8.25
//go:generate ../contracts/abigenw --pkg abv2data --sol ./contracts/AddressBookV2/ABv2DataContract.sol --out ../bindings/abv2data/ABv2DataContract.go --ver 0.8.25
//go:generate ../contracts/abigenw --pkg cnstakingv4 --sol ./contracts/CnStaking/CnStakingV4/CnStakingV4.sol --out ../bindings/cnstakingv4/CnStakingV4.go --ver 0.8.25
//go:generate ../contracts/abigenw --pkg cnstakingv4factory --sol ./contracts/CnStaking/CnStakingV4Factory/CnStakingV4Factory.sol --out ../bindings/cnstakingv4factory/CnStakingV4Factory.go --ver 0.8.25
//go:generate ../contracts/abigenw --pkg publicdelegation --sol ./contracts/PublicDelegation/PublicDelegation.sol --out ../bindings/publicdelegation/PublicDelegation.go --ver 0.8.25
//go:generate ../contracts/abigenw --pkg proxy --sol ./contracts/Proxy/ERC1967Proxy.sol --out ../bindings/proxy/Proxy.go --ver 0.8.25
//go:generate ../contracts/abigenw --pkg beacon --sol ./contracts/Proxy/UpgradeableBeacon.sol --out ../bindings/beacon/UpgradeableBeacon.go --ver 0.8.25
//go:generate ../contracts/abigenw --pkg multicall --sol ./contracts/Multicall/MultiCallContract.sol --out ../bindings/multicall/MultiCallContract.go --ver 0.8.25
//go:generate ../contracts/abigenw --pkg testing --sol ./contracts/testing/MultiCallContractMock.sol --out ../bindings/testing/MultiCallContractMock.go --ver 0.8.25
