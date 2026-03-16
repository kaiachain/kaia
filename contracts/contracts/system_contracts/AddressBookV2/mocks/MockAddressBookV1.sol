// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.25;

/// @title MockAddressBookV1
/// @notice Mock of the original AddressBook with the exact same storage layout (slots 0-12).
///         Provides simplified setter functions for testing the upgrade to ABv2.
///         Deploy at 0x400 → populate data → overlay proxy → verify ABv2 reads it correctly.
/// @dev Storage layout must EXACTLY match the original AddressBook (Solidity 0.4.24).
///      See AddressBookLegacy.sol for slot documentation.
contract MockAddressBookV1 {
    /* ========== ENUMS (matching original) ========== */

    enum RequestState {
        Unknown,
        NotConfirmed,
        Executed,
        ExecutionFailed,
        Expired
    }

    enum Functions {
        Unknown,
        AddAdmin,
        DeleteAdmin,
        UpdateRequirement,
        ClearRequest,
        ActivateAddressBook,
        UpdatePocContract,
        UpdateKirContract,
        RegisterCnStakingContract,
        UnregisterCnStakingContract,
        UpdateSpareContract
    }

    /* ========== STRUCTS (matching original) ========== */

    struct Request {
        Functions functionId;
        bytes32 firstArg;
        bytes32 secondArg;
        bytes32 thirdArg;
        address[] confirmers;
        uint256 initialProposedTime;
        RequestState state;
    }

    /* ========== STORAGE (slots 0-12, matching original AddressBook) ========== */

    address[] private adminList; // slot 0
    uint256 public requirement; // slot 1
    mapping(address => bool) private isAdmin; // slot 2
    mapping(bytes32 => Request) private requestMap; // slot 3
    bytes32[] private pendingRequestList; // slot 4
    address public pocContractAddress; // slot 5
    address public kirContractAddress; // slot 6
    address public spareContractAddress; // slot 7
    mapping(address => uint256) private cnIndexMap; // slot 8
    address[] private cnNodeIdList; // slot 9
    address[] private cnStakingContractList; // slot 10
    address[] private cnRewardAddressList; // slot 11
    bool public isActivated; // slot 12 (byte 0)
    bool public isConstructed; // slot 12 (byte 1)

    /* ========== MOCK SETTERS ========== */

    /// @notice Sets admin list and requirement
    function mockSetAdmins(address[] calldata admins, uint256 req) external {
        delete adminList;
        for (uint256 i = 0; i < admins.length; i++) {
            adminList.push(admins[i]);
            isAdmin[admins[i]] = true;
        }
        requirement = req;
    }

    /// @notice Sets PoC, KIR, and spare contract addresses
    function mockSetContracts(address poc, address kir, address spare) external {
        pocContractAddress = poc;
        kirContractAddress = kir;
        spareContractAddress = spare;
    }

    /// @notice Registers a CN with node ID, staking contract, and reward address
    function mockRegisterCn(address nodeId, address staking, address reward) external {
        uint256 index = cnNodeIdList.length;
        cnIndexMap[nodeId] = index;
        cnNodeIdList.push(nodeId);
        cnStakingContractList.push(staking);
        cnRewardAddressList.push(reward);
    }

    /// @notice Activates the address book (sets both flags to true)
    function mockActivate() external {
        isActivated = true;
        isConstructed = true;
    }

    /// @notice Creates a mock pending request
    function mockAddRequest(
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg,
        address confirmer
    ) external {
        bytes32 id = keccak256(abi.encodePacked(functionId, firstArg, secondArg, thirdArg));
        requestMap[id].functionId = functionId;
        requestMap[id].firstArg = firstArg;
        requestMap[id].secondArg = secondArg;
        requestMap[id].thirdArg = thirdArg;
        requestMap[id].confirmers.push(confirmer);
        requestMap[id].initialProposedTime = block.timestamp;
        requestMap[id].state = RequestState.NotConfirmed;
        pendingRequestList.push(id);
    }

    /* ========== MOCK GETTERS (for verifying data before upgrade) ========== */

    function mockGetCnNodeIdList() external view returns (address[] memory) {
        return cnNodeIdList;
    }

    function mockGetAdminList() external view returns (address[] memory) {
        return adminList;
    }

    function mockGetPendingRequestList() external view returns (bytes32[] memory) {
        return pendingRequestList;
    }
}
