// Copyright 2019 The klaytn Authors
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

pragma solidity ^0.4.24;

/**
 * @title AddressBookMock
 */

contract AddressBookMock {
    event ReviseRewardAddress(address cnNodeId, address prevRewardAddress, address curRewardAddress);

    /*
     *  Constants
     */
    uint256 public constant MAX_ADMIN = 50;
    uint256 public constant MAX_PENDING_REQUEST = 100;
    string public constant CONTRACT_TYPE = "AddressBookMock";
    uint8 public constant CN_NODE_ID_TYPE = 0;
    uint8 public constant CN_STAKING_ADDRESS_TYPE = 1;
    uint8 public constant CN_REWARD_ADDRESS_TYPE = 2;
    uint8 public constant POC_CONTRACT_TYPE = 3;
    uint8 public constant KIR_CONTRACT_TYPE = 4;
    uint256 public constant ONE_WEEK = 1 weeks;
    uint256 public constant TWO_WEEKS = 2 weeks;
    uint256 public constant VERSION = 0;

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

    /*
     *  Storage
     */
    address[] public adminList;
    uint256 public requirement;

    address public pocContractAddress;
    address public kirContractAddress;
    address public spareContractAddress;

    bool public isActivated;
    bool public isConstructed;

    mapping(address => uint256) public cnIndexMap;
    address[] public cnNodeIdList;
    address[] public cnStakingContractList;
    address[] public cnRewardAddressList;

    /*
     *  setter functions
     */
    function constructContract(address[] _adminList, uint256 _requirement) external {
        adminList = _adminList;
        requirement = _requirement;
        isConstructed = true;
    }

    function updateRequirement(uint256 _requirement) public {
        requirement = _requirement;
    }

    function activateAddressBook() public {
        require(isActivated == false, "Already activated.");
        require(isConstructed, "AddressBookMock not constructed");
        require(pocContractAddress != address(0), "PoC contract is not registered.");
        require(kirContractAddress != address(0), "KIR contract is not registered.");
        require(cnNodeIdList.length != 0, "No node ID is listed.");
        require(
            cnNodeIdList.length == cnStakingContractList.length,
            "Invalid length between node IDs and staking contracts."
        );
        require(
            cnStakingContractList.length == cnRewardAddressList.length,
            "Invalid length between staking contracts and reward addresses."
        );
        isActivated = true;
    }

    function updatePocContract(address _pocContractAddress, uint256 /* _version */) public {
        require(isConstructed, "AddressBookMock not constructed");
        pocContractAddress = _pocContractAddress;
    }

    function updateKirContract(address _kirContractAddress, uint256 /* _version */) public {
        require(isConstructed, "AddressBookMock not constructed");
        kirContractAddress = _kirContractAddress;
    }

    function updateSpareContract(address _spareContractAddress) public {
        require(isConstructed, "AddressBookMock not constructed");
        spareContractAddress = _spareContractAddress;
    }

    function registerCnStakingContract(
        address _cnNodeId,
        address _cnStakingContractAddress,
        address _cnRewardAddress
    ) public {
        require(isConstructed, "AddressBookMock not constructed");
        if (cnNodeIdList.length > 0) {
            require(cnNodeIdList[cnIndexMap[_cnNodeId]] != _cnNodeId, "CN node ID already exist.");
        }

        uint256 index = cnNodeIdList.length;
        cnIndexMap[_cnNodeId] = index;
        cnNodeIdList.push(_cnNodeId);
        cnStakingContractList.push(_cnStakingContractAddress);
        cnRewardAddressList.push(_cnRewardAddress);
    }

    function mockRegisterCnStakingContracts(
        address[] _cnNodeIdList,
        address[] _cnStakingContractAddressList,
        address[] _cnRewardAddressList
    ) external {
        require(
            _cnNodeIdList.length == _cnStakingContractAddressList.length,
            "Different cnNodeId and cnStaking lengths"
        );
        require(_cnNodeIdList.length == _cnRewardAddressList.length, "Different cnNodeId and reward lengths");
        for (uint i = 0; i < _cnNodeIdList.length; i++) {
            // skip duplicate
            address _cnNodeId = _cnNodeIdList[i];
            if (cnNodeIdList.length > 0 && cnNodeIdList[cnIndexMap[_cnNodeId]] == _cnNodeId) {
                continue;
            }
            registerCnStakingContract(_cnNodeIdList[i], _cnStakingContractAddressList[i], _cnRewardAddressList[i]);
        }
    }

    function unregisterCnStakingContract(address _cnNodeId) public {
        require(isConstructed, "AddressBookMock not constructed");

        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        require(cnNodeIdList.length > 1, "CN should be more than one.");

        if (index < cnNodeIdList.length - 1) {
            cnNodeIdList[index] = cnNodeIdList[cnNodeIdList.length - 1];
            cnStakingContractList[index] = cnStakingContractList[cnNodeIdList.length - 1];
            cnRewardAddressList[index] = cnRewardAddressList[cnNodeIdList.length - 1];

            cnIndexMap[cnNodeIdList[cnNodeIdList.length - 1]] = index;
        }

        delete cnIndexMap[_cnNodeId];
        delete cnNodeIdList[cnNodeIdList.length - 1];
        cnNodeIdList.length = cnNodeIdList.length - 1;
        delete cnStakingContractList[cnStakingContractList.length - 1];
        cnStakingContractList.length = cnStakingContractList.length - 1;
        delete cnRewardAddressList[cnRewardAddressList.length - 1];
        cnRewardAddressList.length = cnRewardAddressList.length - 1;
    }

    function mockUnregisterCnStakingContracts(address[] _cnNodeIdList) external {
        for (uint i = 0; i < _cnNodeIdList.length; i++) {
            // skip duplicate
            address _cnNodeId = _cnNodeIdList[i];
            uint256 index = cnIndexMap[_cnNodeId];
            if (cnNodeIdList[index] != _cnNodeId) {
                continue;
            }
            unregisterCnStakingContract(_cnNodeIdList[i]);
        }
    }

    function mockUnregisterCnStakingContract(address _cnNodeId) public {
        require(isConstructed, "AddressBookMock not constructed");

        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        // require(cnNodeIdList.length > 1, "CN should be more than one.");

        if (index < cnNodeIdList.length - 1) {
            cnNodeIdList[index] = cnNodeIdList[cnNodeIdList.length - 1];
            cnStakingContractList[index] = cnStakingContractList[cnNodeIdList.length - 1];
            cnRewardAddressList[index] = cnRewardAddressList[cnNodeIdList.length - 1];

            cnIndexMap[cnNodeIdList[cnNodeIdList.length - 1]] = index;
        }

        delete cnIndexMap[_cnNodeId];
        delete cnNodeIdList[cnNodeIdList.length - 1];
        cnNodeIdList.length = cnNodeIdList.length - 1;
        delete cnStakingContractList[cnStakingContractList.length - 1];
        cnStakingContractList.length = cnStakingContractList.length - 1;
        delete cnRewardAddressList[cnRewardAddressList.length - 1];
        cnRewardAddressList.length = cnRewardAddressList.length - 1;
    }

    function revokeRequest(Functions, bytes32, bytes32, bytes32) external {}

    /*
     * submit functions redirected to functions
     */
    function submitAddAdmin(address /* _admin */) external {}

    function submitDeleteAdmin(address /* _admin */) external {}

    function submitUpdateRequirement(uint256 _requirement) external {
        updateRequirement(_requirement);
    }

    function submitClearRequest() external {}

    function submitActivateAddressBook() external {
        activateAddressBook();
    }

    function submitUpdatePocContract(address _pocContractAddress, uint256 _version) external {
        updatePocContract(_pocContractAddress, _version);
    }

    function submitUpdateKirContract(address _kirContractAddress, uint256 _version) external {
        updateKirContract(_kirContractAddress, _version);
    }

    function submitUpdateSpareContract(address _spareContractAddress) external {
        updateSpareContract(_spareContractAddress);
    }

    function submitRegisterCnStakingContract(
        address _cnNodeId,
        address _cnStakingContractAddress,
        address _cnRewardAddress
    ) external {
        registerCnStakingContract(_cnNodeId, _cnStakingContractAddress, _cnRewardAddress);
    }

    function submitUnregisterCnStakingContract(address _cnNodeId) external {
        unregisterCnStakingContract(_cnNodeId);
    }

    function reviseRewardAddress(address _rewardAddress) external {
        bool foundIt = false;
        uint256 index = 0;
        uint256 cnStakingContractListCnt = cnStakingContractList.length;
        for (uint256 i = 0; i < cnStakingContractListCnt; i++) {
            if (cnStakingContractList[i] == msg.sender) {
                foundIt = true;
                index = i;
                break;
            }
        }

        address prevAddress = cnRewardAddressList[index];
        cnRewardAddressList[index] = _rewardAddress;
        emit ReviseRewardAddress(cnNodeIdList[index], prevAddress, cnRewardAddressList[index]);
    }

    /*
     * Getter functions
     */
    function getState() external view returns (address[], uint256) {
        return (adminList, 0);
    }

    function getAllAddress() external view returns (uint8[], address[]) {
        uint8[] memory typeList;
        address[] memory addressList;
        if (isActivated == false) {
            typeList = new uint8[](0);
            addressList = new address[](0);
        } else {
            typeList = new uint8[](cnNodeIdList.length * 3 + 2);
            addressList = new address[](cnNodeIdList.length * 3 + 2);
            uint256 cnNodeCnt = cnNodeIdList.length;
            for (uint256 i = 0; i < cnNodeCnt; i++) {
                //add node id and its type number to array
                typeList[i * 3] = uint8(CN_NODE_ID_TYPE);
                addressList[i * 3] = address(cnNodeIdList[i]);
                //add staking address and its type number to array
                typeList[i * 3 + 1] = uint8(CN_STAKING_ADDRESS_TYPE);
                addressList[i * 3 + 1] = address(cnStakingContractList[i]);
                //add reward address and its type number to array
                typeList[i * 3 + 2] = uint8(CN_REWARD_ADDRESS_TYPE);
                addressList[i * 3 + 2] = address(cnRewardAddressList[i]);
            }
            typeList[cnNodeCnt * 3] = uint8(POC_CONTRACT_TYPE);
            addressList[cnNodeCnt * 3] = address(pocContractAddress);
            typeList[cnNodeCnt * 3 + 1] = uint8(KIR_CONTRACT_TYPE);
            addressList[cnNodeCnt * 3 + 1] = address(kirContractAddress);
        }
        return (typeList, addressList);
    }

    function getAllAddressInfo() external view returns (address[], address[], address[], address, address) {
        return (cnNodeIdList, cnStakingContractList, cnRewardAddressList, pocContractAddress, kirContractAddress);
    }

    function getCnInfo(address _cnNodeId) external view returns (address, address, address) {
        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        return (cnNodeIdList[index], cnStakingContractList[index], cnRewardAddressList[index]);
    }
}

contract AddressBookMockWrong {
    function getAllAddressInfo() external view returns (address[], address[], address[], address, address) {
        return (
            // unequal array lengths
            new address[](2),
            new address[](3),
            new address[](1),
            address(0),
            address(0)
        );
    }
}

contract AddressBookMockThreeCN {
    address public constant dummy = 0x0000000000000000000000000000000000000000;
    // addresses derived from the mnemonic test-junk
    // addresses must be aligned with fixtures.ts:getActors()
    address public constant abookAdmin = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    address public constant cn0 = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;
    address public constant cn1 = 0x90F79bf6EB2c4f870365E785982E1f101E93b906;
    address public constant cn2 = 0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65;

    function getState() external pure returns (address[] memory, uint256) {
        address[] memory adminList = new address[](1);
        adminList[0] = abookAdmin;
        uint256 requirement = 1;
        return (adminList, requirement);
    }

    function getCnInfo(address _cnNodeId) external pure returns (address, address, address) {
        address[3] memory cnList = [cn0, cn1, cn2];

        for (uint256 i = 0; i < cnList.length; i++) {
            if (_cnNodeId == cnList[i]) {
                return (cnList[i], dummy, dummy);
            }
        }

        revert("Invalid CN node ID.");
    }
}

contract AddressBookMockOneCN {
    address public constant dummy = 0x0000000000000000000000000000000000000000;
    // addresses derived from the mnemonic test-junk
    // addresses must be aligned with fixtures.ts:getActors()
    address public constant abookAdmin = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
    address public constant cn0 = 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC;

    function getState() external pure returns (address[] memory, uint256) {
        address[] memory adminList = new address[](1);
        adminList[0] = abookAdmin;
        uint256 requirement = 1;
        return (adminList, requirement);
    }

    function getCnInfo(address _cnNodeId) external pure returns (address, address, address) {
        address[1] memory cnList = [cn0];

        for (uint256 i = 0; i < cnList.length; i++) {
            if (_cnNodeId == cnList[i]) {
                return (cnList[i], dummy, dummy);
            }
        }

        revert("Invalid CN node ID.");
    }
}

contract MockValues {
    address public constant nodeId0 = 0x0000000000000000000000000000000000000F00;
    address public constant nodeId1 = 0x0000000000000000000000000000000000000F03;
    address public constant nodeId2 = 0x0000000000000000000000000000000000000F06;
}

/**
 * @title AddressBookMockTwoCN
 */

contract AddressBookMockTwoCN is MockValues {

    function getAllAddress() external view returns (uint8[] memory typeList, address[] memory addressList) {
        typeList = new uint8[](8);
        addressList = new address[](8);

        typeList[0] = 0; // Node address
        typeList[1] = 1; // Staking address
        typeList[2] = 2; // Reward address
        typeList[3] = 0; // Node address
        typeList[4] = 1; // Staking address
        typeList[5] = 2; // Reward address
        typeList[6] = 3; // POC address
        typeList[7] = 4; // KIR address

        addressList[0] = nodeId0;
        addressList[1] = 0x0000000000000000000000000000000000000F01;
        addressList[2] = 0x0000000000000000000000000000000000000f02;
        addressList[3] = nodeId1;
        addressList[4] = 0x0000000000000000000000000000000000000f04;
        addressList[5] = 0x0000000000000000000000000000000000000f05;
        addressList[6] = 0x0000000000000000000000000000000000000F06;
        addressList[7] = 0x0000000000000000000000000000000000000f07;
    }
}
