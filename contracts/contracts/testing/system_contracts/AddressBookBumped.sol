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
pragma solidity ^0.8.0;

/**
 * @title SafeMath
 * @dev Unsigned math operations with safety checks that revert on error
 */
library SafeMath {
    /**
     * @dev Multiplies two unsigned integers, reverts on overflow.
     */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
        // benefit is lost if 'b' is also tested.
        // See: https://github.com/OpenZeppelin/openzeppelin-solidity/pull/522
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        require(c / a == b);

        return c;
    }

    /**
     * @dev Integer division of two unsigned integers truncating the quotient, reverts on division by zero.
     */
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        // Solidity only automatically asserts when dividing by 0
        require(b > 0);
        uint256 c = a / b;
        // assert(a == b * c + a % b); // There is no case in which this doesn't hold

        return c;
    }

    /**
     * @dev Subtracts two unsigned integers, reverts on overflow (i.e. if subtrahend is greater than minuend).
     */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b <= a);
        uint256 c = a - b;

        return c;
    }

    /**
     * @dev Adds two unsigned integers, reverts on overflow.
     */
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a);

        return c;
    }

    /**
     * @dev Divides two unsigned integers and returns the remainder (unsigned integer modulo),
     * reverts when dividing by zero.
     */
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b != 0);
        return a % b;
    }
}

/**
 * @title AddressBook
 */

contract AddressBookBumped {
    using SafeMath for uint256;
    /*
     *  Events
     */
    event DeployContract(string contractType, address[] adminList, uint256 requirement);
    event AddAdmin(address indexed admin);
    event DeleteAdmin(address indexed admin);
    event UpdateRequirement(uint256 requirement);
    event ClearRequest();
    event SubmitRequest(
        bytes32 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg,
        address[] confirmers
    );
    event ExpiredRequest(
        bytes32 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg,
        address[] confirmers
    );
    event RevokeRequest(
        bytes32 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg,
        address[] confirmers
    );
    event CancelRequest(
        bytes32 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg
    );
    event ExecuteRequestSuccess(
        bytes32 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg
    );
    event ExecuteRequestFailure(
        bytes32 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg
    );

    event ActivateAddressBook();
    event UpdatePocContract(
        address prevPocContractAddress,
        uint256 prevVersion,
        address curPocContractAddress,
        uint256 curVersion
    );
    event UpdateKirContract(
        address prevKirContractAddress,
        uint256 prevVersion,
        address curKirContractAddress,
        uint256 curVersion
    );
    event UpdateSpareContract(address spareContractAddress);
    event RegisterCnStakingContract(address cnNodeId, address cnStakingContractAddress, address cnRewardAddress);
    event UnregisterCnStakingContract(address cnNodeId);
    event ReviseRewardAddress(address cnNodeId, address prevRewardAddress, address curRewardAddress);

    /*
     *  Constants
     */
    uint256 public constant MAX_ADMIN = 50;
    uint256 public constant MAX_PENDING_REQUEST = 100; //Max값에 도달하면 clearRequest()이외에 다른 함수 호출 불가
    string public constant CONTRACT_TYPE = "AddressBook";
    uint8 public constant CN_NODE_ID_TYPE = 0;
    uint8 public constant CN_STAKING_ADDRESS_TYPE = 1;
    uint8 public constant CN_REWARD_ADDRESS_TYPE = 2;
    uint8 public constant POC_CONTRACT_TYPE = 3;
    uint8 public constant KIR_CONTRACT_TYPE = 4;
    uint256 public constant ONE_WEEK = 1 weeks;
    uint256 public constant TWO_WEEKS = 2 weeks;
    uint256 public constant VERSION = 1;

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

    struct Request {
        Functions functionId;
        bytes32 firstArg;
        bytes32 secondArg;
        bytes32 thirdArg;
        address[] confirmers;
        uint256 initialProposedTime;
        RequestState state;
    }

    /*
     *  Storage
     */
    address[] private adminList; //멀티시그 신청 권한을 가지고 있는 관리자 목록
    uint256 public requirement; //멀티시그 리퀘스트가 실행되기 위한 submitRequest의 수
    mapping(address => bool) private isAdmin;
    mapping(bytes32 => Request) private requestMap;
    bytes32[] private pendingRequestList;

    address public pocContractAddress;
    address public kirContractAddress;
    address public spareContractAddress;

    //cnIndex맵을 통해 리스트 순회할 필요 없이 바로 인덱스 접근
    mapping(address => uint256) private cnIndexMap;
    address[] private cnNodeIdList;
    address[] private cnStakingContractList;
    address[] private cnRewardAddressList;

    //모든 셋팅이 끝난 후에 true로 전환되는 flag
    bool public isActivated;
    bool public isConstructed;

    /*
     *  Modifiers
     */
    modifier onlyMultisigTx() {
        require(msg.sender == address(this), "Not a multisig-transaction.");
        _;
    }

    modifier onlyAdmin(address _admin) {
        require(isAdmin[_admin], "Address is not admin.");
        _;
    }

    modifier adminDoesNotExist(address _admin) {
        require(!isAdmin[_admin], "Admin already exits.");
        _;
    }

    modifier notNull(address _address) {
        require(_address != address(0), "Address is null.");
        _;
    }

    modifier validRequirement(uint256 _adminCount, uint256 _requirement) {
        require(
            _adminCount <= MAX_ADMIN && _requirement <= _adminCount && _requirement != 0 && _adminCount != 0,
            "Invalid requirement."
        );
        _;
    }

    /*
     *  Constructor
     */
    function constructContract(
        address[] memory _adminList,
        uint256 _requirement
    ) external validRequirement(_adminList.length, _requirement) {
        // require(msg.sender == 0x854CA8508C8BE2bb1f3C244045786410Cb7D5D0a, "Invalid sender."); //하드코딩된 특정 주소만 initialize 가능
        require(isConstructed == false, "Already constructed.");
        uint256 adminListCnt = _adminList.length;

        isActivated = false;
        for (uint256 i = 0; i < adminListCnt; i++) {
            require(!isAdmin[_adminList[i]] && _adminList[i] != address(0), "Address is null or not unique.");
            isAdmin[_adminList[i]] = true;
        }
        adminList = _adminList;
        requirement = _requirement;
        isConstructed = true;
        emit DeployContract(CONTRACT_TYPE, adminList, requirement);
    }

    /*
     *  Private functions
     */
    function deleteFromPendingRequestList(bytes32 id) private {
        uint256 pendingRequestListCnt = pendingRequestList.length;
        for (uint256 i = 0; i < pendingRequestListCnt; i++) {
            if (id == pendingRequestList[i]) {
                //제거하려는 아이템이 리스트의 마지막이 항목이 아닌 경우
                if (i != pendingRequestListCnt - 1) {
                    pendingRequestList[i] = pendingRequestList[pendingRequestListCnt - 1];
                }
                //맨 마지막 항목 제거
                pendingRequestList.pop();
                break;
            }
        }
    }

    function getId(
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) private pure returns (bytes32) {
        return keccak256(abi.encodePacked(_functionId, _firstArg, _secondArg, _thirdArg));
    }

    function submitRequest(Functions _functionId, bytes32 _firstArg, bytes32 _secondArg, bytes32 _thirdArg) private {
        bytes32 id = getId(_functionId, _firstArg, _secondArg, _thirdArg);

        // 처음 수행되는 Proposal이 아닌 경우
        if (requestMap[id].initialProposedTime != 0) {
            // 최초 Proposal 로부터 2주 경과 후 동일한 request가 있을 시 새로운 request로 간주하고 기존 내용은 초기화
            if (requestMap[id].initialProposedTime + TWO_WEEKS < block.timestamp) {
                deleteFromPendingRequestList(id);
                delete requestMap[id];
            }
            // 최초 Proposal 로부터 1주 이상 경과 된 경우 request가 expired 된 상태
            else if (requestMap[id].initialProposedTime + ONE_WEEK < block.timestamp) {
                if (requestMap[id].state != RequestState.Expired) {
                    requestMap[id].state = RequestState.Expired;
                }
                emit ExpiredRequest(
                    id,
                    msg.sender,
                    _functionId,
                    _firstArg,
                    _secondArg,
                    _thirdArg,
                    requestMap[id].confirmers
                );
            }
            // Confirm
            else if (requestMap[id].initialProposedTime <= block.timestamp) {
                uint256 confirmersCnt = requestMap[id].confirmers.length;
                for (uint256 i = 0; i < confirmersCnt; i++) {
                    require(msg.sender != requestMap[id].confirmers[i], "Msg.sender already requested."); //이미 컨펌한 이력이 없는지 체크
                }
                requestMap[id].confirmers.push(msg.sender);
                emit SubmitRequest(
                    id,
                    msg.sender,
                    _functionId,
                    _firstArg,
                    _secondArg,
                    _thirdArg,
                    requestMap[id].confirmers
                );
            }
        }

        // 처음 수행되는 Proposal인 경우
        if (requestMap[id].initialProposedTime == 0) {
            if (pendingRequestList.length >= MAX_PENDING_REQUEST) {
                require(_functionId == Functions.ClearRequest, "Request list is full.");
            }
            requestMap[id] = Request({
                functionId: _functionId,
                firstArg: _firstArg,
                secondArg: _secondArg,
                thirdArg: _thirdArg,
                initialProposedTime: block.timestamp,
                confirmers: new address[](0),
                state: RequestState.NotConfirmed
            });
            requestMap[id].confirmers.push(msg.sender);
            pendingRequestList.push(id);
            emit SubmitRequest(
                id,
                msg.sender,
                _functionId,
                _firstArg,
                _secondArg,
                _thirdArg,
                requestMap[id].confirmers
            );
        }
    }

    function executeRequest(bytes32 _id) private {
        bool executed = false;
        Request memory _executeRequest = requestMap[_id];

        if (_executeRequest.functionId == Functions.AddAdmin) {
            //bytes4(keccak256("addAdmin(address)")) => 0x70480275
            (executed, ) = address(this).call(abi.encodeWithSelector(this.addAdmin.selector, _executeRequest.firstArg));
        } else if (_executeRequest.functionId == Functions.DeleteAdmin) {
            //bytes4(keccak256("deleteAdmin(address)")) => 0x27e1f7df
            (executed, ) = address(this).call(
                abi.encodeWithSelector(this.deleteAdmin.selector, _executeRequest.firstArg)
            );
        } else if (_executeRequest.functionId == Functions.UpdateRequirement) {
            //bytes4(keccak256("updateRequirement(uint256)")) => 0xc47afb3a
            (executed, ) = address(this).call(
                abi.encodeWithSelector(this.updateRequirement.selector, uint256(_executeRequest.firstArg))
            );
        } else if (_executeRequest.functionId == Functions.ClearRequest) {
            //bytes4(keccak256("clearRequest()")) => 0x4f97638f
            (executed, ) = address(this).call(abi.encodeWithSelector(this.clearRequest.selector));
        } else if (_executeRequest.functionId == Functions.ActivateAddressBook) {
            //bytes4(keccak256("activateAddressBook()")) => 0xcec92466
            (executed, ) = address(this).call(abi.encodeWithSelector(this.activateAddressBook.selector));
        } else if (_executeRequest.functionId == Functions.UpdatePocContract) {
            //bytes4(keccak256("updatePocContract(address,uint256)")) => 0xc7e9de75
            (executed, ) = address(this).call(
                abi.encodeWithSelector(
                    this.updatePocContract.selector,
                    _executeRequest.firstArg,
                    _executeRequest.secondArg
                )
            );
        } else if (_executeRequest.functionId == Functions.UpdateKirContract) {
            //bytes4(keccak256("updateKirContract(address,uint256)")) => 0x4c5d435c
            (executed, ) = address(this).call(
                abi.encodeWithSelector(
                    this.updateKirContract.selector,
                    _executeRequest.firstArg,
                    _executeRequest.secondArg
                )
            );
        } else if (_executeRequest.functionId == Functions.RegisterCnStakingContract) {
            //bytes4(keccak256("registerCnStakingContract(address,address,address)")) => 0x298b3c61
            (executed, ) = address(this).call(
                abi.encodeWithSelector(
                    this.registerCnStakingContract.selector,
                    _executeRequest.firstArg,
                    _executeRequest.secondArg,
                    _executeRequest.thirdArg
                )
            );
        } else if (_executeRequest.functionId == Functions.UnregisterCnStakingContract) {
            //bytes4(keccak256("unregisterCnStakingContract(address)")) => 0x579740db
            (executed, ) = address(this).call(
                abi.encodeWithSelector(this.unregisterCnStakingContract.selector, _executeRequest.firstArg)
            );
        } else if (_executeRequest.functionId == Functions.UpdateSpareContract) {
            //bytes4(keccak256("updateSpareContract(address)")) => 0xafaaf330
            (executed, ) = address(this).call(
                abi.encodeWithSelector(this.updateSpareContract.selector, _executeRequest.firstArg)
            );
        }

        deleteFromPendingRequestList(_id);
        if (executed) {
            if (requestMap[_id].initialProposedTime != 0) {
                // clear 되지 않고 유효한 request인 경우
                requestMap[_id].state = RequestState.Executed;
            }
            emit ExecuteRequestSuccess(
                _id,
                msg.sender,
                _executeRequest.functionId,
                _executeRequest.firstArg,
                _executeRequest.secondArg,
                _executeRequest.thirdArg
            );
        } else {
            if (requestMap[_id].initialProposedTime != 0) {
                // clear 되지 않고 유효한 request인 경우
                requestMap[_id].state = RequestState.ExecutionFailed;
            }
            emit ExecuteRequestFailure(
                _id,
                msg.sender,
                _executeRequest.functionId,
                _executeRequest.firstArg,
                _executeRequest.secondArg,
                _executeRequest.thirdArg
            );
        }
    }

    function checkQuorum(bytes32 _id) private view returns (bool) {
        //컨펌 조건이 채워졌을 경우 실행
        return (requestMap[_id].confirmers.length >= requirement);
    }

    /*
     *  external functions
     */
    function revokeRequest(
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) external onlyAdmin(msg.sender) {
        bytes32 id = getId(_functionId, _firstArg, _secondArg, _thirdArg);

        //아직 요청되지 않은 리퀘스트인지 확인
        require(requestMap[id].initialProposedTime != 0, "Invalid request.");
        require(requestMap[id].state == RequestState.NotConfirmed, "Must be at not-confirmed state.");
        bool foundIt = false;
        uint256 confirmerCnt = requestMap[id].confirmers.length;

        for (uint256 i = 0; i < confirmerCnt; i++) {
            if (msg.sender == requestMap[id].confirmers[i]) {
                foundIt = true;

                // 최초 Proposal 로부터 1주 이상 경과 된 경우 request가 expired 된 상태
                if (requestMap[id].initialProposedTime + ONE_WEEK < block.timestamp) {
                    // 최초 Proposal 로부터 2주 경과 후 동일한 request가 있을 시 해당 request는 delete
                    if (requestMap[id].initialProposedTime + TWO_WEEKS < block.timestamp) {
                        deleteFromPendingRequestList(id);
                        delete requestMap[id];
                    } else {
                        requestMap[id].state = RequestState.Expired;
                    }

                    emit ExpiredRequest(
                        id,
                        msg.sender,
                        _functionId,
                        _firstArg,
                        _secondArg,
                        _thirdArg,
                        requestMap[id].confirmers
                    );
                } else {
                    //제거하려는 아이템이 리스트의 마지막이 항목이 아닌 경우
                    if (i != confirmerCnt - 1) {
                        requestMap[id].confirmers[i] = requestMap[id].confirmers[confirmerCnt - 1];
                    }
                    requestMap[id].confirmers.pop();

                    emit RevokeRequest(
                        id,
                        msg.sender,
                        requestMap[id].functionId,
                        requestMap[id].firstArg,
                        requestMap[id].secondArg,
                        requestMap[id].thirdArg,
                        requestMap[id].confirmers
                    );

                    if (requestMap[id].confirmers.length == 0) {
                        deleteFromPendingRequestList(id);
                        delete requestMap[id];
                        emit CancelRequest(
                            id,
                            msg.sender,
                            requestMap[id].functionId,
                            requestMap[id].firstArg,
                            requestMap[id].secondArg,
                            requestMap[id].thirdArg
                        );
                    }
                }
                break;
            }
        }
        //revoke가 수행되지 않은 경우 revert
        require(foundIt, "Msg.sender has not requested.");
    }

    /*
     *  submit request functions
     */
    function submitAddAdmin(
        address _admin
    )
        external
        onlyAdmin(msg.sender)
        adminDoesNotExist(_admin)
        notNull(_admin)
        validRequirement(adminList.length.add(1), requirement)
    {
        bytes32 id = getId(Functions.AddAdmin, bytes32(uint256(uint160(_admin))), 0, 0);

        submitRequest(Functions.AddAdmin, bytes32(uint256(uint160(_admin))), 0, 0);
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitDeleteAdmin(
        address _admin
    )
        external
        onlyAdmin(_admin)
        onlyAdmin(msg.sender)
        notNull(_admin)
        validRequirement(adminList.length.sub(1), requirement)
    {
        bytes32 id = getId(Functions.DeleteAdmin, bytes32(uint256(uint160(_admin))), 0, 0);

        submitRequest(Functions.DeleteAdmin, bytes32(uint256(uint160(_admin))), 0, 0);
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUpdateRequirement(
        uint256 _requirement
    ) external onlyAdmin(msg.sender) validRequirement(adminList.length, _requirement) {
        require(requirement != _requirement, "Same requirement.");
        bytes32 id = getId(Functions.UpdateRequirement, bytes32(_requirement), 0, 0);

        submitRequest(Functions.UpdateRequirement, bytes32(_requirement), 0, 0);
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitClearRequest() external onlyAdmin(msg.sender) {
        bytes32 id = getId(Functions.ClearRequest, 0, 0, 0);

        submitRequest(Functions.ClearRequest, 0, 0, 0);
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitActivateAddressBook() external onlyAdmin(msg.sender) {
        require(isActivated == false, "Already activated.");
        require(adminList.length != 0, "No admin is listed.");
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

        bytes32 id = getId(Functions.ActivateAddressBook, 0, 0, 0);

        submitRequest(Functions.ActivateAddressBook, 0, 0, 0);
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUpdatePocContract(
        address _pocContractAddress,
        uint256 _version
    ) external notNull(_pocContractAddress) onlyAdmin(msg.sender) {
        //pocContract의 버전체크 함수를 통해 등록하려는 컨트랙트가 pocContract인지, 그리고 등록하려는 버전과 맞는지 확인
        require(PocContractInterface(_pocContractAddress).getPocVersion() == _version, "Invalid PoC version.");

        bytes32 id = getId(
            Functions.UpdatePocContract,
            bytes32(uint256(uint160(_pocContractAddress))),
            bytes32(_version),
            0
        );

        submitRequest(
            Functions.UpdatePocContract,
            bytes32(uint256(uint160(_pocContractAddress))),
            bytes32(_version),
            0
        );
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUpdateKirContract(
        address _kirContractAddress,
        uint256 _version
    ) external notNull(_kirContractAddress) onlyAdmin(msg.sender) {
        //kirContract의 버전체크 함수를 통해 등록하려는 컨트랙트가 kirContract인지, 그리고 등록하려는 버전과 맞는지 확인
        require(KirContractInterface(_kirContractAddress).getKirVersion() == _version, "Invalid KIR version.");

        bytes32 id = getId(
            Functions.UpdateKirContract,
            bytes32(uint256(uint160(_kirContractAddress))),
            bytes32(_version),
            0
        );

        submitRequest(
            Functions.UpdateKirContract,
            bytes32(uint256(uint160(_kirContractAddress))),
            bytes32(_version),
            0
        );
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUpdateSpareContract(address _spareContractAddress) external onlyAdmin(msg.sender) {
        bytes32 id = getId(Functions.UpdateSpareContract, bytes32(uint256(uint160(_spareContractAddress))), 0, 0);

        submitRequest(Functions.UpdateSpareContract, bytes32(uint256(uint160(_spareContractAddress))), 0, 0);
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitRegisterCnStakingContract(
        address _cnNodeId,
        address _cnStakingContractAddress,
        address _cnRewardAddress
    ) external notNull(_cnNodeId) notNull(_cnStakingContractAddress) notNull(_cnRewardAddress) onlyAdmin(msg.sender) {
        if (cnNodeIdList.length > 0) {
            require(cnNodeIdList[cnIndexMap[_cnNodeId]] != _cnNodeId, "CN node ID already exist."); //중복 등록 방지
        }
        //cnStakingContract에 접근하여 nodeId와 rewardAddress가 유효한지 검증
        require(CnStakingContractInterface(_cnStakingContractAddress).nodeId() == _cnNodeId, "Invalid CN node ID.");
        require(
            CnStakingContractInterface(_cnStakingContractAddress).rewardAddress() == _cnRewardAddress,
            "Invalid CN reward address."
        );
        //initialize 이후에만 등록 가능
        require(
            CnStakingContractInterface(_cnStakingContractAddress).isInitialized() == true,
            "CN contract is not initialized."
        );

        bytes32 id = getId(
            Functions.RegisterCnStakingContract,
            bytes32(uint256(uint160(_cnNodeId))),
            bytes32(uint256(uint160(_cnStakingContractAddress))),
            bytes32(uint256(uint160(_cnRewardAddress)))
        );

        submitRequest(
            Functions.RegisterCnStakingContract,
            bytes32(uint256(uint160(_cnNodeId))),
            bytes32(uint256(uint160(_cnStakingContractAddress))),
            bytes32(uint256(uint160(_cnRewardAddress)))
        );
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    function submitUnregisterCnStakingContract(address _cnNodeId) external notNull(_cnNodeId) onlyAdmin(msg.sender) {
        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        require(cnNodeIdList.length > 1, "CN should be more than one.");

        bytes32 id = getId(Functions.UnregisterCnStakingContract, bytes32(uint256(uint160(_cnNodeId))), 0, 0);

        submitRequest(Functions.UnregisterCnStakingContract, bytes32(uint256(uint160(_cnNodeId))), 0, 0);
        if (checkQuorum(id) && requestMap[id].state == RequestState.NotConfirmed) {
            executeRequest(id);
        }
    }

    /*
     *  Multisig functions
     */
    function addAdmin(
        address _admin
    ) external onlyMultisigTx adminDoesNotExist(_admin) validRequirement(adminList.length.add(1), requirement) {
        isAdmin[_admin] = true;
        adminList.push(_admin);
        clearRequest();
        emit AddAdmin(_admin);
    }

    function deleteAdmin(
        address _admin
    ) external onlyMultisigTx onlyAdmin(_admin) validRequirement(adminList.length.sub(1), requirement) {
        uint256 adminCnt = adminList.length;
        isAdmin[_admin] = false;

        for (uint256 i = 0; i < adminCnt - 1; i++) {
            if (adminList[i] == _admin) {
                adminList[i] = adminList[adminCnt - 1];
                break;
            }
        }

        adminList.pop();
        clearRequest();
        emit DeleteAdmin(_admin);
    }

    function updateRequirement(
        uint256 _requirement
    ) external onlyMultisigTx validRequirement(adminList.length, _requirement) {
        require(requirement != _requirement, "Same requirement.");
        requirement = _requirement;
        clearRequest();
        emit UpdateRequirement(_requirement);
    }

    function clearRequest() public onlyMultisigTx {
        uint256 pendingRequestCnt = pendingRequestList.length;

        for (uint256 i = 0; i < pendingRequestCnt; i++) {
            delete requestMap[pendingRequestList[i]];
        }
        delete pendingRequestList;
        emit ClearRequest();
    }

    function activateAddressBook() external onlyMultisigTx {
        require(isActivated == false, "Already activated.");
        require(adminList.length != 0, "No admin is listed.");
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

        emit ActivateAddressBook();
    }

    function updatePocContract(address _pocContractAddress, uint256 _version) external onlyMultisigTx {
        //pocContract의 버전체크 함수를 통해 등록하려는 컨트랙트가 pocContract인지, 그리고 등록하려는 버전과 맞는지 확인
        require(PocContractInterface(_pocContractAddress).getPocVersion() == _version, "Invalid PoC version.");

        address prevPocContractAddress = pocContractAddress;
        pocContractAddress = _pocContractAddress;
        uint256 prevVersion = 0;

        if (prevPocContractAddress != address(0)) {
            prevVersion = PocContractInterface(prevPocContractAddress).getPocVersion();
        }
        emit UpdatePocContract(prevPocContractAddress, prevVersion, _pocContractAddress, _version);
    }

    function updateKirContract(address _kirContractAddress, uint256 _version) external onlyMultisigTx {
        //kirContract의 버전체크 함수를 통해 등록하려는 컨트랙트가 kirContract인지, 그리고 등록하려는 버전과 맞는지 확인
        require(KirContractInterface(_kirContractAddress).getKirVersion() == _version, "Invalid KIR version.");

        address prevKirContractAddress = kirContractAddress;
        kirContractAddress = _kirContractAddress;
        uint256 prevVersion = 0;

        if (prevKirContractAddress != address(0)) {
            prevVersion = KirContractInterface(prevKirContractAddress).getKirVersion();
        }
        emit UpdateKirContract(prevKirContractAddress, prevVersion, _kirContractAddress, _version);
    }

    function updateSpareContract(address _spareContractAddress) external onlyMultisigTx {
        spareContractAddress = _spareContractAddress;
        emit UpdateSpareContract(spareContractAddress);
    }

    function registerCnStakingContract(
        address _cnNodeId,
        address _cnStakingContractAddress,
        address _cnRewardAddress
    ) external onlyMultisigTx {
        if (cnNodeIdList.length > 0) {
            require(cnNodeIdList[cnIndexMap[_cnNodeId]] != _cnNodeId, "CN node ID already exist."); //중복 등록 방지
        }
        //cnStakingContract에 접근하여 nodeId와 rewardAddress가 유효한지 검증
        require(CnStakingContractInterface(_cnStakingContractAddress).nodeId() == _cnNodeId, "Invalid CN node ID.");
        require(
            CnStakingContractInterface(_cnStakingContractAddress).rewardAddress() == _cnRewardAddress,
            "Invalid CN reward address."
        );
        //initialize 이후에만 등록 가능
        require(
            CnStakingContractInterface(_cnStakingContractAddress).isInitialized() == true,
            "CN contract is not initialized."
        );

        uint256 index = cnNodeIdList.length;
        cnIndexMap[_cnNodeId] = index;
        cnNodeIdList.push(_cnNodeId);
        cnStakingContractList.push(_cnStakingContractAddress);
        cnRewardAddressList.push(_cnRewardAddress);

        emit RegisterCnStakingContract(_cnNodeId, _cnStakingContractAddress, _cnRewardAddress);
    }

    function unregisterCnStakingContract(address _cnNodeId) external onlyMultisigTx {
        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        require(cnNodeIdList.length > 1, "CN should be more than one.");

        //제거하려는 CN이 리스트 마지막에 존재하지 않을 경우 마지막 id를 지우려는 id 인덱스에 복사한 후 마지막 아이템을 제거한다
        if (index < cnNodeIdList.length - 1) {
            cnNodeIdList[index] = cnNodeIdList[cnNodeIdList.length - 1];
            cnStakingContractList[index] = cnStakingContractList[cnNodeIdList.length - 1];
            cnRewardAddressList[index] = cnRewardAddressList[cnNodeIdList.length - 1];

            //바뀐 위치에 있는 nodeId에 새로운 인덱스 부여
            cnIndexMap[cnNodeIdList[cnNodeIdList.length - 1]] = index;
        }

        delete cnIndexMap[_cnNodeId];
        cnNodeIdList.pop();
        cnStakingContractList.pop();
        cnRewardAddressList.pop();

        emit UnregisterCnStakingContract(_cnNodeId);
    }

    /*
     * External function
     */
    function reviseRewardAddress(address _rewardAddress) external notNull(_rewardAddress) {
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
        //요청이 들어온 주소가 등록된 cnStakingContract가 아닐 시 revert
        require(foundIt, "Msg.sender is not CN contract.");
        address prevAddress = cnRewardAddressList[index];
        cnRewardAddressList[index] = _rewardAddress;

        emit ReviseRewardAddress(cnNodeIdList[index], prevAddress, cnRewardAddressList[index]);
    }

    /*
     * Getter functions
     */
    function getState() external view returns (address[] memory, uint256) {
        return (adminList, requirement);
    }

    function getPendingRequestList() external view returns (bytes32[] memory) {
        return pendingRequestList;
    }

    function getRequestInfo(
        bytes32 _id
    ) external view returns (Functions, bytes32, bytes32, bytes32, address[] memory, uint256, RequestState) {
        return (
            requestMap[_id].functionId,
            requestMap[_id].firstArg,
            requestMap[_id].secondArg,
            requestMap[_id].thirdArg,
            requestMap[_id].confirmers,
            requestMap[_id].initialProposedTime,
            requestMap[_id].state
        );
    }

    function getRequestInfoByArgs(
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) external view returns (bytes32, address[] memory, uint256, RequestState) {
        bytes32 _id = getId(_functionId, _firstArg, _secondArg, _thirdArg);
        return (_id, requestMap[_id].confirmers, requestMap[_id].initialProposedTime, requestMap[_id].state);
    }

    function getAllAddress() external view returns (uint8[] memory, address[] memory) {
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

    function getAllAddressInfo()
        external
        view
        returns (address[] memory, address[] memory, address[] memory, address, address)
    {
        return (cnNodeIdList, cnStakingContractList, cnRewardAddressList, pocContractAddress, kirContractAddress);
    }

    function getCnInfo(address _cnNodeId) external view notNull(_cnNodeId) returns (address, address, address) {
        uint256 index = cnIndexMap[_cnNodeId];
        require(cnNodeIdList[index] == _cnNodeId, "Invalid CN node ID.");
        return (cnNodeIdList[index], cnStakingContractList[index], cnRewardAddressList[index]);
    }
}

interface CnStakingContractInterface {
    function nodeId() external view returns (address);

    function rewardAddress() external view returns (address);

    function isInitialized() external view returns (bool);
}

interface PocContractInterface {
    function getPocVersion() external pure returns (uint256);
}

interface KirContractInterface {
    function getKirVersion() external pure returns (uint256);
}
