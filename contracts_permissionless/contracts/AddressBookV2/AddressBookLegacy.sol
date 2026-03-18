// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

/// @title AddressBookLegacy
/// @notice Backward-compatible storage layout and legacy getters for the original AddressBook (v1).
/// @dev This abstract contract preserves the storage layout of the original AddressBook at 0x400
///      (Solidity 0.4.24, sequential slots 0-12). After the hard fork upgrade to ABv2 + UUPS proxy,
///      the old storage persists and is NOT cleared. This contract provides:
///      - Storage variables matching the original slot layout (slots 0-12)
///      - Legacy view functions returning real data from ABv2 storage (CN-related) or frozen data (request-related)
///      - Legacy mutating functions that revert with LegacyFunctionDeprecated()
///      ABv2 inherits this to maintain full backward compatibility with external tools.
///
///      Storage safety: Both ABv2 and dependent contracts use ERC-7201 namespaced storage.
abstract contract AddressBookLegacy {
    /* ========== ERRORS ========== */

    /// @notice Thrown when a deprecated legacy function is called
    error LegacyFunctionDeprecated();

    /* ========== LEGACY ENUMS (matching original AddressBook values) ========== */

    enum RequestState {
        Unknown, // 0
        NotConfirmed, // 1
        Executed, // 2
        ExecutionFailed, // 3
        Expired // 4
    }

    enum Functions {
        Unknown, // 0
        AddAdmin, // 1
        DeleteAdmin, // 2
        UpdateRequirement, // 3
        ClearRequest, // 4
        ActivateAddressBook, // 5
        UpdatePocContract, // 6
        UpdateKirContract, // 7
        RegisterCnStakingContract, // 8
        UnregisterCnStakingContract, // 9
        UpdateSpareContract // 10
    }

    /* ========== LEGACY STRUCTS ========== */

    /// @dev Must match exact field order and types of the original Request struct
    ///      for storage-compatible reading from frozen requestMap (slot 3).
    struct Request {
        Functions functionId;
        bytes32 firstArg;
        bytes32 secondArg;
        bytes32 thirdArg;
        address[] confirmers;
        uint256 initialProposedTime;
        RequestState state;
    }

    /* ========== LEGACY CONSTANTS ========== */
    // Note: CONTRACT_TYPE and VERSION are NOT redeclared here — ABv2 defines them (VERSION=2).

    uint256 public constant MAX_ADMIN = 50;
    uint256 public constant MAX_PENDING_REQUEST = 100;
    uint8 public constant CN_NODE_ID_TYPE = 0;
    uint8 public constant CN_STAKING_ADDRESS_TYPE = 1;
    uint8 public constant CN_REWARD_ADDRESS_TYPE = 2;
    uint8 public constant POC_CONTRACT_TYPE = 3;
    uint8 public constant KIR_CONTRACT_TYPE = 4;
    uint256 public constant ONE_WEEK = 1 weeks;
    uint256 public constant TWO_WEEKS = 2 weeks;

    /* ========== LEGACY STORAGE (slots 0-12, matching original AddressBook) ========== */
    // IMPORTANT: Declaration order determines slot positions. Do NOT reorder.

    /// @dev Slot 0: admin list for the internal multisig
    address[] private adminList;

    /// @dev Slot 1: multisig requirement — renamed from `requirement` (was public) to allow
    ///      manual getter override. Storage slot position is unchanged (determined by order, not name).
    uint256 internal _requirement;

    /// @dev Slot 2: admin check mapping
    mapping(address => bool) private isAdmin;

    /// @dev Slot 3: multisig request mapping (Request struct layout must match original)
    mapping(bytes32 => Request) private requestMap;

    /// @dev Slot 4: pending request list
    bytes32[] private pendingRequestList;

    /// @dev Slot 5: PoC contract address (public — auto-generated getter reads frozen value)
    address public pocContractAddress;

    /// @dev Slot 6: KIR contract address (public — auto-generated getter reads frozen value)
    address public kirContractAddress;

    /// @dev Slot 7: spare contract address (public — auto-generated getter reads frozen value)
    address public spareContractAddress;

    /// @dev Slot 8: CN index mapping for fast lookup (frozen, unused by new getters)
    mapping(address => uint256) private cnIndexMap;

    /// @dev Slot 9: CN node ID list (frozen, unused by new getters)
    address[] private cnNodeIdList;

    /// @dev Slot 10: CN staking contract list (frozen, unused by new getters)
    address[] private cnStakingContractList;

    /// @dev Slot 11: CN reward address list (frozen, unused by new getters)
    address[] private cnRewardAddressList;

    /// @dev Slot 12: activation and construction flags (packed bools, 1 byte each)
    bool public isActivated;
    bool public isConstructed;

    /* ========== ABSTRACT: Data bridge to ABv2 storage ========== */
    // ABv2 implements these to provide live node data from ERC-7201 storage.

    /// @notice Returns all registered node data and fund addresses
    /// @dev Implemented by ABv2 — iterates both sets and reads nodeInfo + fund addresses.
    ///      KIF maps to legacy PoC, KEF maps to legacy KIR.
    function _getAllAddressData()
        internal
        view
        virtual
        returns (
            address[] memory nodeIds,
            address[] memory stakingContracts,
            address[] memory rewardAddresses,
            address kifAddress,
            address kefAddress
        );

    /// @notice Returns data for a specific registered node
    /// @dev Implemented by ABv2 — reads from nodeInfo mapping.
    /// @return stakingContract The node's staking contract address
    /// @return rewardAddress The node's reward address
    /// @return exists True if the node is registered (state != Unknown)
    function _getNodeData(
        address nodeId
    ) internal view virtual returns (address stakingContract, address rewardAddress, bool exists);

    /* ========== LEGACY VIEW: requirement() override ========== */

    /// @notice Returns the multisig requirement.
    /// @dev Overrides the old `public requirement` auto-getter. Returns 1 to reflect
    ///      the new single-owner governance model (consistent with getState() returning ([owner], 1)).
    function requirement() external pure returns (uint256) {
        return 1;
    }

    /* ========== LEGACY VIEWS: CN-related (live data from ABv2 + frozen poc/kir) ========== */

    /// @notice Returns all addresses in the legacy format
    /// @dev CN data (nodeId, stakingContract, rewardAddress) and fund addresses (KEF, KIF)
    ///      are read from ABv2's live ERC-7201 storage.
    ///      The isActivated check is preserved for backward compatibility.
    function getAllAddress() external view returns (uint8[] memory typeList, address[] memory addressList) {
        if (!isActivated) {
            typeList = new uint8[](0);
            addressList = new address[](0);
        } else {
            (
                address[] memory nodeIds,
                address[] memory staking,
                address[] memory rewards,
                address kif,
                address kef
            ) = _getAllAddressData();
            uint256 cnNodeCnt = nodeIds.length;
            typeList = new uint8[](cnNodeCnt * 3 + 2);
            addressList = new address[](cnNodeCnt * 3 + 2);
            for (uint256 i = 0; i < cnNodeCnt; ++i) {
                typeList[i * 3] = CN_NODE_ID_TYPE;
                addressList[i * 3] = nodeIds[i];
                typeList[i * 3 + 1] = CN_STAKING_ADDRESS_TYPE;
                addressList[i * 3 + 1] = staking[i];
                typeList[i * 3 + 2] = CN_REWARD_ADDRESS_TYPE;
                addressList[i * 3 + 2] = rewards[i];
            }
            typeList[cnNodeCnt * 3] = POC_CONTRACT_TYPE;
            addressList[cnNodeCnt * 3] = kif;
            typeList[cnNodeCnt * 3 + 1] = KIR_CONTRACT_TYPE;
            addressList[cnNodeCnt * 3 + 1] = kef;
        }
    }

    /// @notice Returns all address info in the legacy format
    /// @dev All data from ABv2's live ERC-7201 storage. KIF maps to PoC, KEF maps to KIR.
    function getAllAddressInfo()
        external
        view
        returns (address[] memory, address[] memory, address[] memory, address, address)
    {
        (
            address[] memory nodeIds,
            address[] memory staking,
            address[] memory rewards,
            address kif,
            address kef
        ) = _getAllAddressData();
        return (nodeIds, staking, rewards, kif, kef);
    }

    /// @notice Returns CN info for a specific node ID
    /// @dev Reads from ABv2's live nodeInfo storage.
    function getCnInfo(address _cnNodeId) external view returns (address, address, address) {
        (address staking, address reward, bool exists) = _getNodeData(_cnNodeId);
        require(exists, "Invalid CN node ID.");
        return (_cnNodeId, staking, reward);
    }

    /* ========== LEGACY MUTATORS: All deprecated ========== */
    // All mutating functions from the original AddressBook are deprecated.
    // Any call with an unrecognized selector (including all legacy submit*, multisig,
    // and admin functions) will hit this fallback and revert.

    fallback() external {
        revert LegacyFunctionDeprecated();
    }
}
