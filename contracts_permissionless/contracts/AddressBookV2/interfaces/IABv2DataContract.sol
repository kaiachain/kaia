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
        address initialOwner;
        uint256 exitThreshold;
        uint256 pauseTimeout;
        uint256 idleTimeout;
        uint256 maxValidatorCount;
        uint256 maxReadyCandidateCount;
        address kefAddress;
        address kifAddress;
        address kpfAddress;
        /// @dev Node addresses for initial validators
        address[] nodeIds;
        /// @dev NodeInfo for each initial validator (parallel to nodeIds); infos[i].manager holds the manager
        NodeInfo[] infos;
    }

    /// @notice Returns the ABv2 implementation (logic) contract address.
    /// @dev Used by the core client to set up the UUPS proxy before calling initialize().
    function implementation() external view returns (address);

    /// @notice Returns all genesis initialization data in a single call.
    /// @return The complete InitData struct
    function getInitData() external view returns (InitData memory);
}
