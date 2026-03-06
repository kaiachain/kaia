// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

import {NodeInfo} from "../../types/Node.sol";

/// @title IABv2DataContract
/// @notice Interface for the read-only data contract that provides genesis configuration
///         for AddressBookV2 initialization. Registered in the system registry (0x401)
///         as "ABv2DataContract".
interface IABv2DataContract {
    /// @notice Aggregated genesis data for ABv2 initialization
    struct InitData {
        /// @dev Owner of AddressBookV2 (can upgrade, suspend, update config)
        address initialOwner;
        /// @dev Proposal failure count that triggers forced exit
        uint256 exitThreshold;
        /// @dev Duration (seconds) for ValPaused timeout
        uint256 pauseTimeout;
        /// @dev Duration (seconds) for ValInactive/ValReady timeout
        uint256 idleTimeout;
        /// @dev Maximum validators in the pipeline (activeSet)
        uint256 maxValidatorCount;
        /// @dev Maximum candidates in CandReady state
        uint256 maxReadyCandidateCount;
        /// @dev Kaia Ecosystem Fund address
        address kefAddress;
        /// @dev Kaia Infrastructure Fund address
        address kifAddress;
        /// @dev Kaia Protocol Fund address
        address kpfAddress;
        /// @dev Node addresses for initial validators
        address[] nodeIds;
        /// @dev NodeInfo for each initial validator (parallel to nodeIds); infos[i].manager holds the manager
        NodeInfo[] infos;
    }

    /// @notice Returns all genesis data in a single call
    /// @return The complete InitData struct
    function getData() external view returns (InitData memory);
}
