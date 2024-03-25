// Copyright 2022 The klaytn Authors
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

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./IGovParam.sol";

/// @dev Contract to store and update governance parameters
/// This contract can be called by node to read the param values in the current block
/// Also, the governance contract can change the parameter values.
contract GovParam is Ownable, IGovParam {
    /// @dev Returns all parameter names that ever existed
    string[] public override paramNames;

    mapping(string => Param[]) private _checkpoints;

    /// @dev Returns all parameter names that ever existed, including those that are currently non-existing
    function getAllParamNames() external view override returns (string[] memory) {
        return paramNames;
    }

    /// @dev Returns all checkpoints of the parameter
    /// @param name The parameter name
    function checkpoints(string calldata name) public view override returns (Param[] memory) {
        return _checkpoints[name];
    }

    /// @dev Returns the last checkpoint whose activation block has passed.
    ///      WARNING: Before calling this function, you must ensure that
    ///               _checkpoints[name].length > 0
    function _param(string memory name) private view returns (Param storage) {
        Param[] storage ckpts = _checkpoints[name];
        uint256 len = ckpts.length;

        // there can be up to one checkpoint whose activation block has not passed yet
        // because setParam() will overwrite if there already exists such a checkpoint
        // thus, if the last checkpoint's activation is in the future,
        // it is guaranteed that the next-to-last is activated
        if (ckpts[len - 1].activation <= block.number) {
            return ckpts[len - 1];
        } else {
            return ckpts[len - 2];
        }
    }

    /// @dev Returns the parameter viewed by the current block
    /// @param name The parameter name
    /// @return (1) Whether the parameter exists, and if the parameter exists, (2) its value
    function getParam(string calldata name) external view override returns (bool, bytes memory) {
        if (_checkpoints[name].length == 0) {
            return (false, "");
        }

        Param memory p = _param(name);
        return (p.exists, p.val);
    }

    /// @dev Average of two integers without overflow
    /// https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.7.3/contracts/utils/math/Math.sol#L34
    function average(uint256 a, uint256 b) internal pure returns (uint256) {
        // (a + b) / 2 can overflow.
        return (a & b) + (a ^ b) / 2;
    }

    /// @dev Returns the parameters used for generating the "blockNumber" block
    ///      WARNING: for future blocks, the result may change
    function getParamAt(string memory name, uint256 blockNumber) public view override returns (bool, bytes memory) {
        uint256 len = _checkpoints[name].length;
        if (len == 0) {
            return (false, "");
        }

        // See https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/extensions/ERC20Votes.sol#L99
        // We run a binary search to look for the earliest checkpoint taken after `blockNumber`.
        // During the loop, the index of the wanted checkpoint remains in the range [low-1, high).
        // With each iteration, either `low` or `high` is moved towards the middle of the range to maintain the invariant.
        // - If the middle checkpoint is after `blockNumber`, we look in [low, mid)
        // - If the middle checkpoint is before or equal to `blockNumber`, we look in [mid+1, high)
        // Once we reach a single value (when low == high), we've found the right checkpoint at the index high-1, if not
        // out of bounds (in which case we're looking too far in the past and the result is 0).
        // Note that if the latest checkpoint available is exactly for `blockNumber`, we end up with an index that is
        // past the end of the array, so we technically don't find a checkpoint after `blockNumber`, but it works out
        // the same.
        uint256 low = 0;
        uint256 high = len;

        Param[] storage ckpts = _checkpoints[name];

        while (low < high) {
            uint256 mid = average(low, high);
            if (ckpts[mid].activation > blockNumber) {
                high = mid;
            } else {
                low = mid + 1;
            }
        }

        // high can't be zero. For high to be zero, The "high = mid" line should be executed when mid is zero.
        // When mid = 0, ckpts[mid].activation is always 0 due to the sentinel checkpoint.
        // Therefore, ckpts[mid].activation <= blockNumber,
        // and the "high = mid" line is never executed.
        return (ckpts[high - 1].exists, ckpts[high - 1].val);
    }

    /// @dev Returns existing parameters viewed by the current block
    function getAllParams() external view override returns (string[] memory, bytes[] memory) {
        // solidity doesn't allow memory arrays to be resized
        // so we calculate the size in advance (existCount)
        // See https://docs.soliditylang.org/en/latest/types.html#allocating-memory-arrays
        uint256 existCount = 0;
        for (uint256 i = 0; i < paramNames.length; i++) {
            Param storage tmp = _param(paramNames[i]);
            if (tmp.exists) {
                existCount++;
            }
        }

        string[] memory names = new string[](existCount);
        bytes[] memory vals = new bytes[](existCount);

        uint256 idx = 0;
        for (uint256 i = 0; i < paramNames.length; i++) {
            Param storage tmp = _param(paramNames[i]);
            if (tmp.exists) {
                names[idx] = paramNames[i];
                vals[idx] = tmp.val;
                idx++;
            }
        }
        return (names, vals);
    }

    /// @dev Returns parameters used for generating the "blockNumber" block
    ///      WARNING: for future blocks, the result may change
    function getAllParamsAt(uint256 blockNumber) external view override returns (string[] memory, bytes[] memory) {
        // solidity doesn't allow memory arrays to be resized
        // so we calculate the size in advance (existCount)
        // See https://docs.soliditylang.org/en/latest/types.html#allocating-memory-arrays
        uint256 existCount = 0;
        for (uint256 i = 0; i < paramNames.length; i++) {
            (bool exists, ) = getParamAt(paramNames[i], blockNumber);
            if (exists) {
                existCount++;
            }
        }

        string[] memory names = new string[](existCount);
        bytes[] memory vals = new bytes[](existCount);

        uint256 idx = 0;
        for (uint256 i = 0; i < paramNames.length; i++) {
            (bool exists, bytes memory val) = getParamAt(paramNames[i], blockNumber);
            if (exists) {
                names[idx] = paramNames[i];
                vals[idx] = val;
                idx++;
            }
        }

        return (names, vals);
    }

    /// @dev Returns all parameters as stored in the contract
    function getAllCheckpoints() external view override returns (string[] memory, Param[][] memory) {
        Param[][] memory ckptsArr = new Param[][](paramNames.length);
        for (uint256 i = 0; i < paramNames.length; i++) {
            ckptsArr[i] = _checkpoints[paramNames[i]];
        }
        return (paramNames, ckptsArr);
    }

    /// @dev Returns all parameters as stored in the contract
    function setParam(string calldata name, bool exists, bytes calldata val, uint256 activation)
        public
        override
        onlyOwner
    {
        require(bytes(name).length > 0, "GovParam: name cannot be empty");
        require(
            activation > block.number,
            "GovParam: activation must be in the future"
        );
        require(
            !exists || val.length > 0,
            "GovParam: val must not be empty if exists=true"
        );
        require(
            exists || val.length == 0,
            "GovParam: val must be empty if exists=false"
        );

        Param memory newParam = Param(activation, exists, val);
        Param[] storage ckpts = _checkpoints[name];

        // for a new parameter, push occurs twice
        // (1) sentinel checkpoint
        // (2) newParam
        // this ensures that if name is in paramNames, then ckpts.length >= 2
        if (ckpts.length == 0) {
            paramNames.push(name);

            // insert a sentinel checkpoint
            ckpts.push(Param(0, false, ""));
        }

        uint256 lastPos = ckpts.length - 1;
        // if the last checkpoint's activation is in the past, push newParam
        // otherwise, overwrite the last checkpoint with newParam
        if (ckpts[lastPos].activation <= block.number) {
            ckpts.push(newParam);
        } else {
            ckpts[lastPos] = newParam;
        }

        emit SetParam(name, exists, val, activation);
    }

    /// @dev Updates the parameter to the given state at the relative activation block
    function setParamIn(string calldata name, bool exists, bytes calldata val, uint256 relativeActivation)
        external
        override
        onlyOwner
    {
        uint256 activation = block.number + relativeActivation;
        setParam(name, exists, val, activation);
    }
}
