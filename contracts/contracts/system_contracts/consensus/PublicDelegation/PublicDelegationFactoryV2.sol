// Copyright 2025 The kaia Authors
// This file is part of the kaia library.
//
// The kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the kaia library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {PublicDelegation} from "./PublicDelegation.sol";
import {IPublicDelegationFactoryV2, IPublicDelegation} from "./IPublicDelegationFactoryV2.sol";

// Simple factory contract to deploy PublicDelegation contracts.
// The caller of `deployPublicDelegation` will be the corresponding CnStakingV3MultiSig.
contract PublicDelegationFactoryV2 is IPublicDelegationFactoryV2 {
    /* ========== CONSTANTS ========== */

    uint256 public constant VERSION = 2;

    string public constant CONTRACT_TYPE = "PublicDelegationFactory";

    mapping(address => address) private _deployedPublicDelegation; // CnStakingV3MultiSig => PublicDelegation

    /* ========== DEPLOYMENT ========== */

    function deployPublicDelegation(
        IPublicDelegation.PDConstructorArgs memory _pdArgs
    ) external override returns (IPublicDelegation publicDelegation) {
        bytes memory constructorArgs = abi.encode(msg.sender, _pdArgs);
        bytes memory initCode = abi.encodePacked(type(PublicDelegation).creationCode, constructorArgs);
        bytes32 salt = keccak256(constructorArgs);

        assembly {
            publicDelegation := create2(0, add(initCode, 0x20), mload(initCode), salt)
        }
        if (address(publicDelegation) == address(0)) revert DeploymentFailed();

        _deployedPublicDelegation[msg.sender] = address(publicDelegation);

        emit PublicDelegationDeployed(msg.sender, address(publicDelegation));

        return publicDelegation;
    }

    function isDeployedBy(address _cnStakingV3MultiSig, address _publicDelegation) external view returns (bool) {
        return _deployedPublicDelegation[_cnStakingV3MultiSig] == _publicDelegation;
    }
}
