// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {IABv2DataContract} from "./interfaces/IABv2DataContract.sol";
import {NodeInfo, BlsPublicKeyInfo} from "../types/Node.sol";
import {NodeVerifier} from "../libraries/NodeVerifier.sol";

/// @title ABv2DataContract
/// @notice Read-only data contract holding genesis configuration for AddressBookV2.
///         Deployed before hardfork, registered in Registry (0x401) as "ABv2DataContract".
///         All validation happens in the constructor — ABv2.initialize() trusts this data.
contract ABv2DataContract is IABv2DataContract {
    /* ========== ERRORS ========== */

    error InvalidInput();

    /* ========== IMMUTABLES (scalars) ========== */

    /// @inheritdoc IABv2DataContract
    address public immutable override implementation;

    address public immutable initialOwner;
    address public immutable initialSuspender;
    address public immutable initialConfigurator;
    uint256 public immutable pfsThreshold;
    uint256 public immutable cfsThreshold;
    uint256 public immutable pauseTimeout;
    uint256 public immutable idleTimeout;
    uint256 public immutable maxNodeCount;
    uint256 public immutable maxValActivePausedCount;
    uint256 public immutable maxCandReadyCount;
    address public immutable kefAddress;
    address public immutable kifAddress;
    address public immutable kpfAddress;

    /* ========== STORAGE (arrays) ========== */

    /// @dev Address uniqueness registry (storage mapping as scratch space for constructor validation)
    mapping(address => bool) private _usedAddresses;

    address[] private _nodeIds;
    NodeInfo[] private _infos;

    /* ========== CONSTRUCTOR ========== */

    /// @notice Validates and stores all genesis data. Reverts on any invalid input.
    /// @param _implementation The ABv2 logic contract address for proxy deployment
    /// @param data The complete genesis configuration
    constructor(address _implementation, InitData memory data) {
        // Validate implementation address
        if (_implementation == address(0)) revert InvalidInput();

        // Validate config
        if (data.initialOwner == address(0)) revert InvalidInput();
        if (data.initialSuspender == address(0)) revert InvalidInput();
        if (data.initialConfigurator == address(0)) revert InvalidInput();
        if (data.pfsThreshold == 0) revert InvalidInput();
        if (data.cfsThreshold == 0) revert InvalidInput();
        if (data.pauseTimeout == 0) revert InvalidInput();
        if (data.idleTimeout == 0) revert InvalidInput();
        if (data.maxNodeCount == 0) revert InvalidInput();
        if (data.maxValActivePausedCount == 0) revert InvalidInput();
        if (data.maxCandReadyCount == 0) revert InvalidInput();
        if (data.kefAddress == address(0)) revert InvalidInput();
        if (data.kifAddress == address(0)) revert InvalidInput();
        if (data.kpfAddress == address(0)) revert InvalidInput();

        // Validate array length consistency
        uint256 len = data.nodeIds.length;
        if (len != data.infos.length) revert InvalidInput();
        if (len > data.maxNodeCount) revert InvalidInput();

        // Store immutables
        implementation = _implementation;
        initialOwner = data.initialOwner;
        initialSuspender = data.initialSuspender;
        initialConfigurator = data.initialConfigurator;
        pfsThreshold = data.pfsThreshold;
        cfsThreshold = data.cfsThreshold;
        pauseTimeout = data.pauseTimeout;
        idleTimeout = data.idleTimeout;
        maxNodeCount = data.maxNodeCount;
        maxValActivePausedCount = data.maxValActivePausedCount;
        maxCandReadyCount = data.maxCandReadyCount;
        kefAddress = data.kefAddress;
        kifAddress = data.kifAddress;
        kpfAddress = data.kpfAddress;

        // Validate and store nodes
        for (uint256 i; i < len; ++i) {
            address nodeId = data.nodeIds[i];
            NodeInfo memory info = data.infos[i];

            if (info.manager == address(0)) revert InvalidInput();
            NodeVerifier.registerNodeGenesis(
                _usedAddresses,
                nodeId,
                info.stakingContract,
                info.rewardAddress,
                info.blsInfo
            );

            _nodeIds.push(nodeId);
            _infos.push(info);
        }
    }

    /* ========== GETTERS ========== */

    /// @inheritdoc IABv2DataContract
    function getInitData() external view override returns (InitData memory) {
        return
            InitData({
                initialOwner: initialOwner,
                initialSuspender: initialSuspender,
                initialConfigurator: initialConfigurator,
                pfsThreshold: pfsThreshold,
                cfsThreshold: cfsThreshold,
                pauseTimeout: pauseTimeout,
                idleTimeout: idleTimeout,
                maxNodeCount: maxNodeCount,
                maxValActivePausedCount: maxValActivePausedCount,
                maxCandReadyCount: maxCandReadyCount,
                kefAddress: kefAddress,
                kifAddress: kifAddress,
                kpfAddress: kpfAddress,
                nodeIds: _nodeIds,
                infos: _infos
            });
    }
}
