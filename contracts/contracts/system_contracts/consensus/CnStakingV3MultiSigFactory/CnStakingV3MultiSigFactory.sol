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

import {ICnStakingV3MultiSigChunk} from "./ICnStakingV3MultiSigChunk.sol";
import {ICnStakingV3MultiSigFactory} from "./ICnStakingV3MultiSigFactory.sol";

contract CnStakingV3MultiSigFactory is ICnStakingV3MultiSigFactory {
    bytes32 public constant CNSTAKINGV3_MULTISIG_BYTECODE_HASH =
        0x2f15da75ebe93ab05e8a43462d97f4e7d2f5df0c61063df35c04efd4a58eb005;

    ICnStakingV3MultiSigChunk public immutable chunkOne;
    ICnStakingV3MultiSigChunk public immutable chunkTwo;

    mapping(address => address[]) private _deployedCnStakingV3MultiSigList;

    constructor(address _chunkOne, address _chunkTwo) {
        chunkOne = ICnStakingV3MultiSigChunk(_chunkOne);
        chunkTwo = ICnStakingV3MultiSigChunk(_chunkTwo);

        require(chunkOne.CHUNK_ID() == 1 && chunkTwo.CHUNK_ID() == 2, "Invalid chunk ordering");
        require(
            keccak256(_assembleBytecode()) == CNSTAKINGV3_MULTISIG_BYTECODE_HASH,
            "Invalid CnStakingV3MultiSig bytecode hash"
        );
    }

    function deployCnStakingV3MultiSig(
        address _contractValidator,
        address _nodeId,
        address _rewardAddress,
        address[] memory _cnAdminlist,
        uint256 _requirement,
        uint256[] memory _unlockTime,
        uint256[] memory _unlockAmount
    ) public returns (address cnStakingV3MultiSig) {
        bytes memory constructorArgs = abi.encode(
            _contractValidator,
            _nodeId,
            _rewardAddress,
            _cnAdminlist,
            _requirement,
            _unlockTime,
            _unlockAmount
        );
        bytes memory initCode = abi.encodePacked(_assembleBytecode(), constructorArgs);
        bytes32 salt = keccak256(constructorArgs);

        assembly {
            cnStakingV3MultiSig := create2(0, add(initCode, 0x20), mload(initCode), salt)
        }
        if (cnStakingV3MultiSig == address(0)) revert DeploymentFailed();
        _deployedCnStakingV3MultiSigList[msg.sender].push(cnStakingV3MultiSig);

        emit CnStakingV3MultiSigDeployed(msg.sender, cnStakingV3MultiSig);

        return cnStakingV3MultiSig;
    }

    function _assembleBytecode() private view returns (bytes memory) {
        return abi.encodePacked(chunkOne.CHUNK_BYTECODE(), chunkTwo.CHUNK_BYTECODE());
    }

    function getDeployedCnStakingV3MultiSigList(address _deployer) public view returns (address[] memory) {
        return _deployedCnStakingV3MultiSigList[_deployer];
    }

    function isDeployedBy(address _deployer, address _cnStakingV3MultiSig) public view returns (bool) {
        for (uint256 i = 0; i < _deployedCnStakingV3MultiSigList[_deployer].length; i++) {
            if (_deployedCnStakingV3MultiSigList[_deployer][i] == _cnStakingV3MultiSig) {
                return true;
            }
        }
        return false;
    }
}
