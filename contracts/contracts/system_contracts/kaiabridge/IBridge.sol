// Copyright 2024 The klaytn Authors
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
pragma solidity 0.8.24;

import "./EnumerableSetUint64.sol";

abstract contract IBridge {
    //////////////////// Struct ////////////////////
    /// @dev ProvisionData struct. This struct is identical with the wasm event struct
    struct ProvisionData {
        uint64 seq;
        string sender; // FNSA address
        address receiver;
        uint256 amount;
    }

    struct ProvisionIndividualEvent {
        uint64 seq;
        string sender; // FNSA address
        address receiver;
        uint256 amount;
        uint64 txID;
        address operator;
    }

    struct ProvisionConfirmedEvent {
        uint64 seq;
        string sender; // FNSA address
        address receiver;
        uint256 amount;
    }

    struct SwapRequest {
        uint64 seq;
        address sender;
        string receiver;
        uint256 amount;
    }

    //////////////////// Modifier ////////////////////
    modifier onlyOperator {
        require(operator == msg.sender, "KAIA::Bridge: Not an operator");
        _;
    }

    modifier onlyGuardian {
        require(guardian == msg.sender, "KAIA::Bridge: Not an guardian");
        _;
    }

    modifier onlyJudge {
        require(judge == msg.sender, "KAIA::Bridge: Not an judge");
        _;
    }

    modifier notPause {
        require(!pause, "KAIA::Bridge: Bridge has been paused");
        _;
    }

    modifier inPause {
        require(pause, "KAIA::Bridge: Bridge has not been paused");
        _;
    }

    modifier transferEnable {
        require(transferFromKaiaOn, "KAIA::Bridge: Swap reqeust from KAIA is disabled");
        _;
    }

    //////////////////// Event ////////////////////
    /// @dev Emitted when a `transferFromKaiaOn` is changed
    event TransferFromKaiaOnOffChanged(bool indexed transferFromKaiaOn, bool indexed set);

    /// @dev Emitted when a provision is submitted
    event ProvisionConfirm(ProvisionConfirmedEvent provision);

    /// @dev Emitted when a provision is confirmed
    event Provision(ProvisionIndividualEvent provision);

    /// @dev Emitted when a provision is removed
    event RemoveProvision(ProvisionData provision);

    /// @dev Emitted when KAIA is charged
    event KAIACharged(address sender, uint256 amount);

    /// @dev Emitted when transfer (swap request) is done
    event Transfer(SwapRequest lockInfo);

    /// @dev Emitted when `minLockableKAIA` is changed
    event MinLockableKAIAChange(uint256 indexed beforeMinLock, uint256 indexed newMinLock);

    /// @dev Emitted when `maxLockableKAIA` is changed
    event MaxLockableKAIAChange(uint256 indexed beforeMaxLock, uint256 indexed newMaxLock);

    /// @dev Emitted when claim is done
    event Claim(ProvisionData indexed provision);

    /// @dev Emitted when provision receiver is changed
    event ProvisionReceiverChanged(address indexed beforeReceiver, address indexed newReceiver);

    /// @dev Emitted when `maxTryTransfer` is changed
    event MaxTryTransferChange(uint256 indexed maxTryTransfer, uint256 indexed newMaxTryTransfer);

    /// @dev Emitted when the operator contract address is changed
    event ChangeOperator(address indexed beforeOperator, address indexed newOperator);

    /// @dev Emitted when the guardian contract address is changed
    event ChangeGuardian(address indexed beforeGuardian, address indexed newGuardian);

    /// @dev Emitted when the judge contract address is changed
    event ChangeJudge(address indexed beforeJudge, address indexed newJudge);

    /// @dev Emitted when `maxTryTransfer` is changed
    event ChangeAddrValidation(bool indexed before, bool now);

    /// @dev Emitted when mintlock duration is updated
    event ChangeTransferTimeLock(uint256 indexed time);

    /// @dev Emitted when claim is blocked
    event HoldClaim(uint256 indexed seq, uint256 indexed time);

    /// @dev Emitted when claim is released
    event ReleaseClaim(uint256 indexed seq, uint256 indexed time);

    /// @dev Emitted when bridge is paused
    event BridgePause(string indexed msg);

    /// @dev Emitted when bridge is resumed
    event BridgeResume(string indexed msg);

    /// @dev Emitted when the bridge service period is changed
    event ChangeBridgeServicePeriod(uint256 indexed bridgeServicePeriod, uint256 indexed newPeriod);

    /// @dev Emitted when the bridge balance is burned
    event BridgeBalanceBurned(uint256 indexed bridgeBalance);

    //////////////////// Exported functions ////////////////////
    /// @dev Set the transferable option from KAIA to FNSA
    /// @param set Enable transfer function if it is true
    function changeTransferEnable(bool set) external virtual;

    /// @dev A gateway function that triggers token swap between FNSA and KAIA
    /// @param prov Burning provision
    function provision(ProvisionData calldata prov) external virtual;

    /// @dev Return true if provisioned, return 0 otherwise.
    /// @param seq sequence number
    function isProvisioned(uint64 seq) external virtual view returns (bool);

    /// @dev A list version of the `isProvisioned`
    /// @param from start seqeunce number
    /// @param to end seqeunce number
    function isProvisionedRange(uint64 from, uint64 to) external virtual view returns (uint64[] memory);

    /// @dev Try to mint for a provision submitted before
    /// @param seq ProvisionData sequence
    function requestClaim(uint64 seq) external virtual returns (bool);

    /// @dev A version of multiple {IBridge-tryClaim} in one transaction
    /// @param range Unclaimed provision set iteration range
    function requestBatchClaim(uint64 range) external virtual;

    /// @dev Remove a confirmed provision
    /// @param seq sequence number
    function removeProvision(uint64 seq) external virtual;

    /// @dev Unclaimable provision is resolved by guardian group. Receiver address must not be a contract address
    /// @param seq sequence number
    /// @param newReceiver new reciever address
    function resolveUnclaimable(uint64 seq, address newReceiver) external virtual;

    /// @dev Change the operator contract address
    /// @param newOperator Operator contract address
    function changeOperator(address newOperator) external virtual;

    /// @dev Change the guardian contract address
    /// @param newGuardian guardian contract address
    function changeGuardian(address newGuardian) external virtual;

    /// @dev Change the judge contract address
    /// @param newJudge judge contract address
    function changeJudge(address newJudge) external virtual;

    /// @dev Convert raw tx data to `Claim` struct. `try/catch` pattern may be reqruied to bypass runtime exception
    /// @param data raw tx data
    function bytes2Provision(bytes calldata data) external virtual pure returns (ProvisionData memory);

    /// @dev Update mintLockableKAIA with a new value
    /// @param newMinLockableKAIA new value of mintLockableKAIA
    function changeMinLockableKAIA(uint256 newMinLockableKAIA) external virtual;

    /// @dev Update maxtLockableKAIA with a new value
    /// @param newMaxLockableKAIA new value of maxtLockableKAIA
    function changeMaxLockableKAIA(uint256 newMaxLockableKAIA) external virtual;

    /// @dev Update bridge address
    /// @param newMaxTryTransfer new value of maxTryTransfer
    function changeMaxTryTransfer(uint256 newMaxTryTransfer) external virtual;

    /// @dev Transfer the amount of to be sent and record its locking information
    /// @param receiver receiver addresss
    function transfer(string calldata receiver) external payable virtual;

    /// @dev Set address validation on and off
    /// @param onOff Enable address validation if it is true, disable otherwise
    function setAddrValidation(bool onOff) external virtual;

    /// @dev Change mintlock duration
    /// @param duration Transferlock duration
    function changeTransferTimeLock(uint256 duration) external virtual;

    /// @dev Test if mintlock duration passed over
    /// @param seq Claim sequence
    function isPassedTimeLockDuration(uint256 seq) external virtual returns (bool);

    /// @dev The time of minting is extended to indefinitely
    /// @param seq Claim sequence
    function holdClaim(uint256 seq) external virtual;

    /// @dev The time of minting is set to `TRANSFERLOCK`
    /// @param seq Claim sequence
    function releaseClaim(uint256 seq) external virtual;

    /// @dev Pause the bridge and emits the pause message
    /// @param pauseMsg pause message
    function pauseBridge(string calldata pauseMsg) external virtual;

    /// @dev Resume the bridge and emits the resume message
    /// @param resumeMsg resume message
    function resumeBridge(string calldata resumeMsg) external virtual;

    /// @dev Change the bridge service period
    /// @param newPeriod New period to be replaced
    function changeBridgeServicePeriod(uint256 newPeriod) external virtual;

    /// @dev Burn the bridge balance (sending to 0xDEAD address)
    function burnBridgeBalance() external virtual;

    /// @dev Get all locked
    function getAllSwapRequests() external virtual view returns (SwapRequest[] memory);

    /// @dev Get claim ranges
    /// @param from Range start number
    /// @param to Range end number
    function getSwapRequests(uint256 from, uint256 to) external virtual view returns (SwapRequest[] memory);

    /// @dev Return `claimCandidates` set
    function getClaimCandidates() external virtual view returns (uint64[] memory);

    /// @dev Return `claimCandidates` set with the given range
    /// @param range Range of set
    function getClaimCandidatesRange(uint64 range) external virtual view returns (uint64[] memory);

    /// @dev Return `claimFailures` set
    function getClaimFailures() external virtual view returns (uint64[] memory);

    /// @dev Return `claimFailures` set
    /// @param range Range of set
    function getClaimFailuresRange(uint64 range) external virtual view returns (uint64[] memory);

    //////////////////// Constant ////////////////////
    uint256 constant public KAIA_UNIT = 10e18;
    uint256 constant INFINITE = type(uint64).max;
    address constant BURN_TARGET = 0x000000000000000000000000000000000000dEaD;

    //////////////////// Storage variables ////////////////////
    uint256 public bridgeServiceStarted;
    uint256 public bridgeServicePeriod;

    bool public transferFromKaiaOn;
    address public operator;
    address public guardian;
    address public judge;
    uint256 public greatestConfirmedSeq; // Gratest seq, which are known on this contract
    uint256 public nProvisioned;
    uint256 public nClaimed;
    uint256 public minLockableKAIA;
    uint256 public maxLockableKAIA;
    uint64 public seq;
    uint64 public nextProvisionSeq;
    uint256 public maxTryTransfer;
    bool public addrValidationOn;
    uint256 public TRANSFERLOCK;
    mapping (uint256 => uint256) public timelocks; // <sequence, residual lock duration>
    bool public pause;

    SwapRequest[] public locked;
    mapping (uint64 => ProvisionData) public provisions; // <sequence, provision info>
    mapping (uint64 => bool) public claimed; // <sequence, is claimed?>
    mapping (uint256 => uint256) public transferFail; // <seq, # of transfer failure>
    mapping (uint64 => uint256) public seq2BlockNum; // <seq, block number>
    EnumerableSet.UintSet claimFailures; // claim failures by contract receiver
    EnumerableSet.UintSet claimCandidates; // unclaimed provisions sequence numbers

    uint256[100] __gap;
}
