// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {ERC20Upgradeable} from "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import {Math} from "@openzeppelin/contracts/utils/math/Math.sol";
import {Address} from "@openzeppelin/contracts/utils/Address.sol";
import {ReentrancyGuardTransient} from "@openzeppelin/contracts/utils/ReentrancyGuardTransient.sol";
import {IPublicDelegation, IKIP163} from "./interfaces/IPublicDelegation.sol";
import {ICnStaking} from "../CnStaking/CnStakingV4/interfaces/ICnStaking.sol";

/// @title PublicDelegation (V2)
/// @notice Tokenized staking vault for Kaia consensus nodes.
///
/// PublicDelegation is an interest-bearing token (pdKAIA) based on staked KAIA in a CnStaking contract.
/// Its math is based on ERC4626 — exchange rate = totalAssets / totalSupply.
/// As block rewards accumulate, each pdKAIA share becomes worth more KAIA.
///
/// Key changes from V1:
/// 1. Beacon proxy upgradeable (Initializable + ERC20Upgradeable + OwnableUpgradeable).
/// 2. pdKAIA is fully transferable (standard ERC20, no _update override).
/// 3. Redelegation is always enabled (no toggle check).
/// 4. ERC-7201 namespaced storage.
/// 5. Target validation delegated to CnStaking (no ValidContract checks in PD).
contract PublicDelegation is
    IPublicDelegation,
    Initializable,
    ERC20Upgradeable,
    OwnableUpgradeable,
    ReentrancyGuardTransient
{
    using Math for uint256;
    using Address for address payable;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /* ========== CONSTANTS ========== */

    string public constant CONTRACT_TYPE = "PublicDelegation";
    uint256 public constant VERSION = 2;
    uint256 public constant MAX_COMMISSION_RATE = 1e4;
    uint256 public constant COMMISSION_DENOMINATOR = 1e4;

    /* ========== ERC-7201 NAMESPACED STORAGE ========== */

    /// @custom:storage-location erc7201:publicdelegation.storage.PDStorage
    struct PDStorage {
        /// @notice The base CnStaking contract this PD is linked to.
        ICnStaking baseCnStaking;
        /// @notice Commission recipient address.
        address commissionTo;
        /// @notice Commission rate (1e4 = 100%).
        uint256 commissionRate;
        /// @notice User -> list of withdrawal request IDs.
        mapping(address => uint256[]) userRequestIds;
        /// @notice Request ID -> owner address.
        mapping(uint256 => address) requestIdToOwner;
    }

    // keccak256(abi.encode(uint256(keccak256("publicdelegation.storage.PDStorage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant PD_STORAGE_LOCATION = 0x5bce486aa234ee45b8f26e442e90eb587481d7fb346f5ce504cb8d8242c62700;

    function _getPDStorage() private pure returns (PDStorage storage $) {
        assembly {
            $.slot := PD_STORAGE_LOCATION
        }
    }

    /* ========== MODIFIERS ========== */

    modifier notNull(address _address) {
        _notNull(_address);
        _;
    }

    function _notNull(address _address) private pure {
        if (_address == address(0)) revert ZeroAddress();
    }

    /* ========== INITIALIZE ========== */

    /// @inheritdoc IPublicDelegation
    /// @dev Called once by the factory after deploying the BeaconProxy.
    function initialize(
        address _baseCnStaking,
        PDConstructorArgs memory _args
    ) external initializer notNull(_baseCnStaking) notNull(_args.owner) {
        if (_args.commissionRate > MAX_COMMISSION_RATE) revert CommissionRateTooHigh();

        __ERC20_init(
            string(abi.encodePacked(_args.gcName, " Public Delegated KAIA")),
            string(abi.encodePacked(_args.gcName, "-pdKAIA"))
        );
        __Ownable_init(_args.owner);

        PDStorage storage $ = _getPDStorage();
        $.baseCnStaking = ICnStaking(payable(_baseCnStaking));
        $.commissionTo = _args.commissionTo;
        $.commissionRate = _args.commissionRate;

        emit DeployPublicDelegation(CONTRACT_TYPE, _baseCnStaking, _args);
    }

    /* ========== OWNER FUNCTIONS ========== */

    /// @inheritdoc IPublicDelegation
    function updateCommissionTo(address _commissionTo) external override onlyOwner notNull(_commissionTo) nonReentrant {
        _sweepAndStake(address(0), 0);

        PDStorage storage $ = _getPDStorage();
        address prev = $.commissionTo;
        $.commissionTo = _commissionTo;

        emit UpdateCommissionTo(prev, _commissionTo);
    }

    /// @inheritdoc IPublicDelegation
    function updateCommissionRate(uint256 _commissionRate) external override onlyOwner nonReentrant {
        if (_commissionRate > MAX_COMMISSION_RATE) revert CommissionRateTooHigh();

        _sweepAndStake(address(0), 0);

        PDStorage storage $ = _getPDStorage();
        uint256 prev = $.commissionRate;
        $.commissionRate = _commissionRate;

        emit UpdateCommissionRate(prev, _commissionRate);
    }

    /* ========== STAKING FUNCTIONS ========== */

    /// @inheritdoc IPublicDelegation
    function stake() external payable override nonReentrant {
        _sweepAndStake(msg.sender, msg.value);
    }

    /// @inheritdoc IKIP163
    function stakeFor(address _recipient) external payable override notNull(_recipient) nonReentrant {
        _sweepAndStake(_recipient, msg.value);
    }

    /// @inheritdoc IPublicDelegation
    receive() external payable override nonReentrant {
        _sweepAndStake(msg.sender, msg.value);
    }

    /* ========== WITHDRAWAL FUNCTIONS ========== */

    /// @inheritdoc IPublicDelegation
    function withdraw(address _recipient, uint256 _assets) external override notNull(_recipient) nonReentrant {
        _sweepAndStake(address(0), 0);

        uint256 _shares = previewWithdraw(_assets);
        _burn(msg.sender, _shares);
        _requestWithdrawal(msg.sender, _recipient, _assets);

        emit Redeemed(msg.sender, _recipient, _assets, _shares);
    }

    /// @inheritdoc IPublicDelegation
    function redeem(address _recipient, uint256 _shares) external override notNull(_recipient) nonReentrant {
        _sweepAndStake(address(0), 0);

        uint256 _assets = previewRedeem(_shares);
        _burn(msg.sender, _shares);
        _requestWithdrawal(msg.sender, _recipient, _assets);

        emit Redeemed(msg.sender, _recipient, _assets, _shares);
    }

    /// @inheritdoc IPublicDelegation
    function cancelApprovedStakingWithdrawal(uint256 _requestId) external override nonReentrant {
        _sweepAndStake(address(0), 0);
        _onlyRequestOwner(_requestId);

        ICnStaking cn = _cnStaking();

        (, uint256 _assets, , ) = cn.getApprovedStakingWithdrawalInfo(_requestId);
        uint256 _shares = previewDeposit(_assets);
        _mint(msg.sender, _shares);

        cn.cancelApprovedStakingWithdrawal(_requestId);

        emit RequestCancelWithdrawal(msg.sender, _requestId);
    }

    /// @inheritdoc IPublicDelegation
    function claim(uint256 _requestId) external override nonReentrant {
        _sweepAndStake(address(0), 0);
        _onlyRequestOwner(_requestId);

        ICnStaking cn = _cnStaking();
        cn.withdrawApprovedStaking(_requestId);

        (, uint256 _assets, , ICnStaking.WithdrawalStakingState _state) = cn.getApprovedStakingWithdrawalInfo(
            _requestId
        );

        // If auto-canceled (2*STAKE_LOCKUP passed), re-mint shares
        if (_state == ICnStaking.WithdrawalStakingState.Canceled) {
            uint256 _shares = _convertToShares(_assets, totalSupply(), _totalAssets() - _assets);
            _mint(msg.sender, _shares);
            return;
        }

        emit Claimed(msg.sender, _requestId);
    }

    /* ========== REDELEGATION FUNCTIONS ========== */

    /// @inheritdoc IPublicDelegation
    function redelegateByAssets(address _targetCnStaking, uint256 _assets) external override nonReentrant {
        _sweepAndStake(address(0), 0);

        uint256 _shares = previewWithdraw(_assets);
        _burn(msg.sender, _shares);
        _redelegate(_targetCnStaking, _assets);
    }

    /// @inheritdoc IPublicDelegation
    function redelegateByShares(address _targetCnStaking, uint256 _shares) external override nonReentrant {
        _sweepAndStake(address(0), 0);

        uint256 _assets = previewRedeem(_shares);
        _burn(msg.sender, _shares);
        _redelegate(_targetCnStaking, _assets);
    }

    /* ========== SWEEP ========== */

    /// @inheritdoc IPublicDelegation
    function sweep() external override nonReentrant {
        _sweepAndStake(address(0), 0);
    }

    /* ========== PRIVATE FUNCTIONS ========== */

    /// @dev Reverts if caller is not the owner of the withdrawal request.
    function _onlyRequestOwner(uint256 _requestId) private view {
        if (_getPDStorage().requestIdToOwner[_requestId] != msg.sender) revert NotRequestOwner();
    }

    /// @dev Core sweep-and-stake logic (CEI order).
    ///      1. Calculate accumulated reward (address(this).balance - _assets).
    ///      2. Deduct commission from reward.
    ///      3. Mint shares if _recipient is provided (Effects).
    ///      4. Delegate to CnStaking, then send commission (Interactions).
    function _sweepAndStake(address _recipient, uint256 _assets) private {
        PDStorage storage $ = _getPDStorage();
        ICnStaking cn = _cnStaking();

        unchecked {
            uint256 _reward = address(this).balance - _assets;
            uint256 _commission = _calcCommission(_reward, $.commissionRate);

            // Effects: mint shares before any external calls
            if (_recipient != address(0)) {
                uint256 _baseAssets = cn.staking() - cn.unstaking() + _reward - _commission;
                uint256 _shares = _convertToShares(_assets, totalSupply(), _baseAssets);
                if (_shares == 0) revert StakeAmountTooLow();
                _mint(_recipient, _shares);
                emit Staked(_recipient, _assets, _shares);
            }

            // Interactions: external calls after all state changes
            uint256 _toStake = _reward + _assets - _commission;
            if (_toStake > 0) {
                cn.delegate{value: _toStake}();
            }

            if (_commission > 0) {
                payable($.commissionTo).sendValue(_commission);
                emit SendCommission($.commissionTo, _commission);
            }
        }
    }

    /// @dev Request withdrawal from CnStaking.
    function _requestWithdrawal(address _owner, address _recipient, uint256 _assets) private {
        if (_assets == 0) revert WithdrawalAmountTooLow();

        PDStorage storage $ = _getPDStorage();
        uint256 _id = _cnStaking().approveStakingWithdrawal(_recipient, _assets);

        $.userRequestIds[_owner].push(_id);
        $.requestIdToOwner[_id] = _owner;

        emit RequestWithdrawal(_owner, _recipient, _id, _assets);
    }

    /// @dev Execute redelegation via CnStaking.
    function _redelegate(address _targetCnStaking, uint256 _assets) private {
        if (_assets == 0) revert RedelegateAmountTooLow();

        _cnStaking().redelegate(msg.sender, _targetCnStaking, _assets);

        emit Redelegated(msg.sender, _targetCnStaking, _assets);
    }

    /// @dev Returns the base CnStaking contract this PD is linked to.
    function _cnStaking() private view returns (ICnStaking) {
        return _getPDStorage().baseCnStaking;
    }

    /// @dev Calculate commission amount.
    function _calcCommission(uint256 _amount, uint256 _rate) private pure returns (uint256) {
        return _amount.mulDiv(_rate, COMMISSION_DENOMINATOR, Math.Rounding.Floor);
    }

    /// @dev Convert assets to shares with custom supply and total assets.
    function _convertToShares(
        uint256 _assets,
        uint256 _customSupply,
        uint256 _customAssets
    ) private pure returns (uint256) {
        return _customSupply == 0 ? _assets : _assets.mulDiv(_customSupply, _customAssets, Math.Rounding.Floor);
    }

    /// @dev Total staking minus unstaking.
    function _totalStaking() private view returns (uint256) {
        ICnStaking cn = _cnStaking();
        unchecked {
            return cn.staking() - cn.unstaking();
        }
    }

    /// @dev Total reward (balance) minus commission.
    function _pureReward() private view returns (uint256) {
        uint256 _reward = address(this).balance;
        uint256 _commission = _calcCommission(_reward, _getPDStorage().commissionRate);
        unchecked {
            return _reward - _commission;
        }
    }

    /// @dev Total assets = staking net + pure reward.
    function _totalAssets() private view returns (uint256) {
        unchecked {
            return _totalStaking() + _pureReward();
        }
    }

    /* ========== PUBLIC GETTERS ========== */

    /// @inheritdoc IPublicDelegation
    function baseCnStaking() external view override returns (address) {
        return address(_getPDStorage().baseCnStaking);
    }

    /// @inheritdoc IPublicDelegation
    function commissionTo() external view override returns (address) {
        return _getPDStorage().commissionTo;
    }

    /// @inheritdoc IPublicDelegation
    function commissionRate() external view override returns (uint256) {
        return _getPDStorage().commissionRate;
    }

    /// @inheritdoc IPublicDelegation
    function totalAssets() public view override returns (uint256) {
        return _totalAssets();
    }

    /// @inheritdoc IPublicDelegation
    function convertToShares(uint256 _assets) public view override returns (uint256) {
        uint256 supply = totalSupply();
        return supply == 0 ? _assets : _assets.mulDiv(supply, totalAssets(), Math.Rounding.Floor);
    }

    /// @inheritdoc IPublicDelegation
    function convertToAssets(uint256 _shares) public view override returns (uint256) {
        uint256 supply = totalSupply();
        return supply == 0 ? _shares : _shares.mulDiv(totalAssets(), supply, Math.Rounding.Floor);
    }

    /// @inheritdoc IPublicDelegation
    function previewDeposit(uint256 _assets) public view override returns (uint256) {
        return convertToShares(_assets);
    }

    /// @inheritdoc IPublicDelegation
    function previewWithdraw(uint256 _assets) public view override returns (uint256) {
        uint256 supply = totalSupply();
        return supply == 0 ? _assets : _assets.mulDiv(supply, totalAssets(), Math.Rounding.Ceil);
    }

    /// @inheritdoc IPublicDelegation
    function previewRedeem(uint256 _shares) public view override returns (uint256) {
        return convertToAssets(_shares);
    }

    /// @inheritdoc IPublicDelegation
    function maxRedeem(address _owner) public view override returns (uint256) {
        return balanceOf(_owner);
    }

    /// @inheritdoc IPublicDelegation
    function maxWithdraw(address _owner) public view override returns (uint256) {
        return previewRedeem(balanceOf(_owner));
    }

    /// @inheritdoc IKIP163
    function reward() public view override returns (uint256) {
        return _pureReward();
    }

    /// @inheritdoc IPublicDelegation
    function userRequestIds(address _owner, uint256 _index) external view override returns (uint256) {
        return _getPDStorage().userRequestIds[_owner][_index];
    }

    /// @inheritdoc IPublicDelegation
    function requestIdToOwner(uint256 _requestId) external view override returns (address) {
        return _getPDStorage().requestIdToOwner[_requestId];
    }

    /// @inheritdoc IPublicDelegation
    function getCurrentWithdrawalRequestState(
        uint256 _requestId
    ) public view override returns (WithdrawalRequestState) {
        ICnStaking cn = _cnStaking();
        (, , uint256 _withdrawableFrom, ICnStaking.WithdrawalStakingState _state) = cn.getApprovedStakingWithdrawalInfo(
            _requestId
        );

        if (_withdrawableFrom == 0) return WithdrawalRequestState.Undefined;
        if (_state == ICnStaking.WithdrawalStakingState.Canceled) return WithdrawalRequestState.Canceled;
        if (_state == ICnStaking.WithdrawalStakingState.Transferred) return WithdrawalRequestState.Withdrawn;

        uint256 _withdrawableUntil;
        unchecked {
            _withdrawableUntil = _withdrawableFrom + cn.STAKE_LOCKUP();
        }

        if (block.timestamp < _withdrawableFrom) return WithdrawalRequestState.Requested;
        if (block.timestamp < _withdrawableUntil) return WithdrawalRequestState.Withdrawable;

        return WithdrawalRequestState.PendingCancel;
    }

    /// @inheritdoc IPublicDelegation
    function getUserRequestCount(address _owner) public view override returns (uint256) {
        return _getPDStorage().userRequestIds[_owner].length;
    }

    /// @inheritdoc IPublicDelegation
    function getUserRequestIdsWithState(
        address _owner,
        WithdrawalRequestState _state
    ) public view override returns (uint256[] memory) {
        PDStorage storage $ = _getPDStorage();
        uint256[] storage ids = $.userRequestIds[_owner];
        uint256 len = ids.length;

        uint256[] memory result = new uint256[](len);
        uint256 cnt = 0;

        for (uint256 i = 0; i < len; ++i) {
            if (getCurrentWithdrawalRequestState(ids[i]) == _state) {
                unchecked {
                    result[cnt++] = ids[i];
                }
            }
        }

        assembly {
            mstore(result, cnt)
        }
        return result;
    }

    /// @inheritdoc IPublicDelegation
    function getUserRequestIds(address _owner) public view override returns (uint256[] memory) {
        return _getPDStorage().userRequestIds[_owner];
    }
}
