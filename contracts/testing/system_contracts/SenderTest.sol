// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.0;

/**
 * Test contract to represent KGF contract implementing getState()
 */
contract SenderTest1 {
    address[] _adminList;
    uint256 public minReq = 1;

    constructor() {
        _adminList.push(msg.sender);
    }

    /*
     * Getter functions
     */
    function getState() external view returns (address[] memory, uint256) {
        return (_adminList, minReq);
    }

    function emptyAdminList() public {
        _adminList.pop();
    }

    function changeMinReq(uint256 req) public {
        minReq = req;
    }

    function addAdmin(address admin) public {
        _adminList.push(admin);
    }

    /// @dev Add dummy admin addresses to test large adminList.
    function addDummyAdmins(uint256 count) public {
        address admin = 0xF000000000000000000000000000000000000001;
        for (uint256 i = 0; i < count; i++) {
            admin = address(uint160(admin) + 1);
            _adminList.push(admin);
        }
    }

    /*
     *  Deposit function
     */
    /// @dev Fallback function that allows to deposit KLAY
    fallback() external payable {
        require(msg.value > 0, "Invalid value.");
    }
}

/**
 * Test contract to represent KIR contract implementing getState()
 */
contract SenderTest2 {
    address[] _adminList;

    constructor() {
        _adminList.push(msg.sender);
    }

    /*
     * Getter functions
     */
    function getState() external view returns (address[] memory, uint256) {
        return (_adminList, 1);
    }

    /*
     *  Deposit function
     */
    /// @dev Fallback function that allows to deposit KLAY
    fallback() external payable {
        require(msg.value > 0, "Invalid value.");
    }
}
