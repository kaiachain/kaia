// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {ReentrancyGuardTransient} from "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";
import {ICnStaking} from "./interfaces/ICnStaking.sol";
import {ICnStakingV4Factory} from "../CnStakingV4Factory/interfaces/ICnStakingV4Factory.sol";
import {IRegistry} from "../../system/IRegistry.sol";
import {IStakingTracker} from "../../system/IStakingTracker.sol";
import {IKIP163, IPublicDelegation} from "../../PublicDelegation/interfaces/IPublicDelegation.sol";

/// @title CnStaking
/// @notice Simplified staking contract for Kaia consensus nodes.
///
/// Key features:
/// 1. Beacon proxy pattern — all instances share one implementation via UpgradeableBeacon.
/// 2. No initial lockup — only delegated staking.
/// 3. No custom multisig — uses OwnableUpgradeable. External multisigs (e.g. Gnosis Safe) can be the owner.
/// 4. No redundant identity fields — no gcId, rewardAddress, voterAddress, stakingTracker storage.
///    StakingTracker is resolved dynamically from the system registry at 0x0...401.
/// 5. PublicDelegation is optional — factory deploys with or without PD.
///    When PD is present, redelegation is always enabled. When PD is absent, redelegation is disabled.
///
/// Access control:
/// - Staking (delegate/receive): PD-only when PD is set; anyone when PD is not set.
/// - Unstaking (approve/cancel/withdraw): PD-only when PD is set; owner-only when PD is not set.
/// - Redelegation: PD-only (only available when PD is set).
/// - handleRedelegation: validates caller is a factory-deployed CnStaking contract.
///
/// @dev delegate() and receive() intentionally omit nonReentrant because handleRedelegation
///      calls PD.stakeFor which calls back to this contract's delegate(). Adding nonReentrant
///      to delegate() would break this cross-contract call chain.
contract CnStakingV4 is ICnStaking, Initializable, OwnableUpgradeable, ReentrancyGuardTransient {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /* ========== CONSTANTS ========== */

    string public constant CONTRACT_TYPE = "CnStakingContract";
    uint256 public constant VERSION = 4;
    uint256 public constant STAKE_LOCKUP = 1 weeks;

    address internal constant REGISTRY_ADDRESS = 0x0000000000000000000000000000000000000401;

    /* ========== ERC-7201 NAMESPACED STORAGE ========== */

    /// @custom:storage-location erc7201:cnstakingv4.storage.StakingStorage
    struct StakingStorage {
        /// @notice The PublicDelegation contract (set once in initialize, address(0) if PD disabled).
        address publicDelegation;
        /// @notice Total delegated stake held by this contract.
        uint256 staking;
        /// @notice Total amount pending in approved withdrawal requests.
        uint256 unstaking;
        /// @notice Monotonically increasing withdrawal request counter.
        uint256 withdrawalRequestCount;
        /// @notice Withdrawal request details by ID.
        mapping(uint256 => WithdrawalRequest) withdrawalRequestMap;
        /// @notice Per-user last redelegation timestamp (anti-hopping).
        mapping(address => uint256) lastRedelegation;
    }

    // keccak256(abi.encode(uint256(keccak256("cnstakingv4.storage.StakingStorage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant STAKING_STORAGE_LOCATION =
        0x91cd8fd1bd76934695389d6e4c0e33fd1c626573ce745c07eba2138ce41b2f00;

    function _getStakingStorage() private pure returns (StakingStorage storage $) {
        assembly {
            $.slot := STAKING_STORAGE_LOCATION
        }
    }

    /* ========== MODIFIERS ========== */

    /// @dev Staking access: PD-only when PD is set, anyone otherwise.
    modifier onlyStaker() {
        _onlyStaker();
        _;
    }

    function _onlyStaker() private view {
        address pd = _getStakingStorage().publicDelegation;
        if (pd != address(0)) {
            if (msg.sender != pd) revert NotStaker();
        }
    }

    /// @dev Unstaking access: PD-only when PD is set, owner-only otherwise.
    modifier onlyUnstakingManager() {
        _onlyUnstakingManager();
        _;
    }

    function _onlyUnstakingManager() private view {
        address pd = _getStakingStorage().publicDelegation;
        if (pd != address(0)) {
            if (msg.sender != pd) revert NotUnstakingManager();
        } else {
            if (msg.sender != owner()) revert NotUnstakingManager();
        }
    }

    /// @dev Ensures address is non-zero.
    modifier notNull(address _address) {
        _notNull(_address);
        _;
    }

    function _notNull(address _address) private pure {
        if (_address == address(0)) revert ZeroAddress();
    }

    /* ========== INITIALIZE ========== */

    /// @inheritdoc ICnStaking
    /// @dev Called once by the factory after deploying the BeaconProxy.
    function initialize(address _owner) external initializer notNull(_owner) {
        __Ownable_init(_owner);

        emit DeployCnStakingV4(CONTRACT_TYPE, address(0));
    }

    /// @inheritdoc ICnStaking
    /// @dev Called once by the factory after deploying both CnStaking and PD beacon proxies.
    ///      Verifies the PD's baseCnStaking points back to this contract.
    function initializeWithPD(
        address _owner,
        address _publicDelegation
    ) external initializer notNull(_owner) notNull(_publicDelegation) {
        if (IPublicDelegation(payable(_publicDelegation)).baseCnStaking() != address(this)) {
            revert BaseCnStakingMismatch();
        }

        __Ownable_init(_owner);
        _getStakingStorage().publicDelegation = _publicDelegation;

        emit DeployCnStakingV4(CONTRACT_TYPE, _publicDelegation);
    }

    /* ========== STAKING FUNCTIONS ========== */

    /// @inheritdoc ICnStaking
    /// @dev When PD is enabled, only the PD contract can call this.
    ///      When PD is disabled, anyone can call this.
    ///      No nonReentrant — called back by PD during handleRedelegation flow.
    function delegate() external payable override onlyStaker {
        _delegateKaia();
    }

    /// @inheritdoc ICnStaking
    /// @dev Same access control as delegate(). No nonReentrant for same reason.
    receive() external payable override onlyStaker {
        _delegateKaia();
    }

    /* ========== UNSTAKING FUNCTIONS ========== */

    /// @inheritdoc ICnStaking
    function approveStakingWithdrawal(
        address _to,
        uint256 _value
    ) external override notNull(_to) onlyUnstakingManager nonReentrant returns (uint256 id) {
        StakingStorage storage $ = _getStakingStorage();
        if (_value == 0 || $.unstaking + _value > $.staking) revert InvalidValue();

        id = $.withdrawalRequestCount;

        uint256 time;
        unchecked {
            $.withdrawalRequestCount++;
            time = block.timestamp + STAKE_LOCKUP;
            $.unstaking += _value;
        }

        $.withdrawalRequestMap[id] = WithdrawalRequest({
            to: _to,
            value: _value,
            withdrawableFrom: time,
            state: WithdrawalStakingState.Unknown
        });

        _refreshStake();
        emit ApproveStakingWithdrawal(id, _to, _value, time);
    }

    /// @inheritdoc ICnStaking
    function cancelApprovedStakingWithdrawal(uint256 _id) external override onlyUnstakingManager nonReentrant {
        (StakingStorage storage $, WithdrawalRequest storage request) = _getActiveRequest(_id);

        request.state = WithdrawalStakingState.Canceled;
        unchecked {
            $.unstaking -= request.value;
        }

        _refreshStake();
        emit CancelApprovedStakingWithdrawal(_id, request.to, request.value);
    }

    /// @inheritdoc ICnStaking
    /// @dev If STAKE_LOCKUP has passed: execute the withdrawal.
    ///      If 2*STAKE_LOCKUP has passed: auto-cancel the withdrawal.
    function withdrawApprovedStaking(uint256 _id) external override onlyUnstakingManager nonReentrant {
        (StakingStorage storage $, WithdrawalRequest storage request) = _getActiveRequest(_id);
        if (request.value > $.staking) revert InvalidValue();
        if (request.withdrawableFrom > block.timestamp) revert NotWithdrawableYet();

        uint256 withdrawableUntil;
        unchecked {
            withdrawableUntil = request.withdrawableFrom + STAKE_LOCKUP;
        }

        if (withdrawableUntil <= block.timestamp) {
            // Auto-cancel: 2*STAKE_LOCKUP has passed
            request.state = WithdrawalStakingState.Canceled;
            unchecked {
                $.unstaking -= request.value;
            }

            _refreshStake();
            emit CancelApprovedStakingWithdrawal(_id, request.to, request.value);
        } else {
            // Execute withdrawal
            address to = request.to;
            uint256 value = request.value;

            request.state = WithdrawalStakingState.Transferred;
            unchecked {
                $.staking -= value;
                $.unstaking -= value;
            }

            _refreshStake();
            emit WithdrawApprovedStaking(_id, to, value);

            (bool success, ) = to.call{value: value}("");
            if (!success) revert TransferFailed();
        }
    }

    /* ========== REDELEGATION FUNCTIONS ========== */

    /// @inheritdoc ICnStaking
    /// @dev Only callable by the PD contract. Redelegation is always enabled when PD exists.
    ///      Anti-hopping: user cannot redelegate again within STAKE_LOCKUP of their last redelegation.
    function redelegate(
        address _user,
        address _targetCnStaking,
        uint256 _value
    ) external override notNull(_user) nonReentrant {
        StakingStorage storage $ = _getStakingStorage();
        address pdAddr = $.publicDelegation;
        if (pdAddr == address(0) || msg.sender != pdAddr) revert RedelegationDisabled();
        if (_targetCnStaking == address(this)) revert InvalidTarget();
        if (!_isFactoryDeployed(_targetCnStaking)) revert InvalidTarget();
        if (_value == 0 || $.unstaking + _value > $.staking) revert InvalidValue();
        if ($.lastRedelegation[_user] != 0 && $.lastRedelegation[_user] + STAKE_LOCKUP > block.timestamp) {
            revert RedelegationCooldown();
        }

        unchecked {
            $.staking -= _value;
        }

        _refreshStake();
        emit Redelegation(_user, _targetCnStaking, _value);

        ICnStaking(payable(_targetCnStaking)).handleRedelegation{value: _value}(_user);
    }

    /// @inheritdoc ICnStaking
    /// @dev Validates the caller is a factory-deployed CnStaking contract.
    ///      Stakes KAIA on behalf of the user via PD's stakeFor().
    function handleRedelegation(address _user) external payable override notNull(_user) nonReentrant {
        StakingStorage storage $ = _getStakingStorage();
        address pdAddr = $.publicDelegation;
        if (pdAddr == address(0)) revert RedelegationDisabled();
        if (!_isFactoryDeployed(msg.sender)) revert InvalidTarget();

        $.lastRedelegation[_user] = block.timestamp;

        IKIP163 pd = IKIP163(payable(pdAddr));

        uint256 expected;
        unchecked {
            expected = address(this).balance + pd.reward();
        }

        emit HandleRedelegation(_user, msg.sender, address(this), msg.value);

        pd.stakeFor{value: msg.value}(_user);
        if (expected != address(this).balance) revert InvalidStakeFor();
    }

    /* ========== INTERNAL FUNCTIONS ========== */

    /// @dev Add delegated stake and refresh tracker.
    function _delegateKaia() private {
        if (msg.value == 0) revert ZeroValue();

        unchecked {
            _getStakingStorage().staking += msg.value;
        }

        _refreshStake();
        emit DelegateKaia(msg.sender, msg.value);
    }

    /// @dev Load and validate a withdrawal request is active (exists and in Unknown state).
    function _getActiveRequest(
        uint256 _id
    ) private view returns (StakingStorage storage $, WithdrawalRequest storage request) {
        $ = _getStakingStorage();
        request = $.withdrawalRequestMap[_id];
        if (request.to == address(0)) revert WithdrawalNotFound();
        if (request.state != WithdrawalStakingState.Unknown) revert InvalidWithdrawalState();
    }

    /// @dev Resolve StakingTracker from registry and refresh stake.
    function _refreshStake() private {
        address tracker = IRegistry(REGISTRY_ADDRESS).getActiveAddr("StakingTracker");
        if (tracker != address(0)) {
            IStakingTracker(tracker).refreshStake(address(this));
        }
    }

    /// @dev Resolve factory from registry and check if target is factory-deployed CnStaking.
    function _isFactoryDeployed(address _target) private view returns (bool) {
        address factoryAddr = IRegistry(REGISTRY_ADDRESS).getActiveAddr("CnStakingFactory");
        if (factoryAddr == address(0)) return false;
        return ICnStakingV4Factory(factoryAddr).isDeployedCnStaking(_target);
    }

    /* ========== PUBLIC GETTERS ========== */

    /// @inheritdoc ICnStaking
    function publicDelegation() external view override returns (address) {
        return _getStakingStorage().publicDelegation;
    }

    /// @inheritdoc ICnStaking
    function staking() external view override returns (uint256) {
        return _getStakingStorage().staking;
    }

    /// @inheritdoc ICnStaking
    function unstaking() external view override returns (uint256) {
        return _getStakingStorage().unstaking;
    }

    /// @inheritdoc ICnStaking
    function withdrawalRequestCount() external view override returns (uint256) {
        return _getStakingStorage().withdrawalRequestCount;
    }

    /// @inheritdoc ICnStaking
    function lastRedelegation(address _account) external view override returns (uint256) {
        return _getStakingStorage().lastRedelegation[_account];
    }

    /// @inheritdoc ICnStaking
    function getApprovedStakingWithdrawalIds(
        uint256 _from,
        uint256 _to,
        WithdrawalStakingState _state
    ) external view override returns (uint256[] memory ids) {
        StakingStorage storage $ = _getStakingStorage();
        uint256 count = $.withdrawalRequestCount;
        uint256 end = (_to == 0 || _to >= count) ? count : _to;

        if (_from >= end) return new uint256[](0);

        ids = new uint256[](end - _from);
        uint256 cnt = 0;

        for (uint256 i = _from; i < end; ++i) {
            if ($.withdrawalRequestMap[i].state == _state) {
                unchecked {
                    ids[cnt++] = i;
                }
            }
        }

        assembly {
            mstore(ids, cnt)
        }
    }

    /// @inheritdoc ICnStaking
    function getApprovedStakingWithdrawalInfo(
        uint256 _index
    )
        external
        view
        override
        returns (address to, uint256 value, uint256 withdrawableFrom, WithdrawalStakingState state)
    {
        WithdrawalRequest storage request = _getStakingStorage().withdrawalRequestMap[_index];
        return (request.to, request.value, request.withdrawableFrom, request.state);
    }
}
