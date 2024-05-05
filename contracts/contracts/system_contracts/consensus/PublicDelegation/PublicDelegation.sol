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
pragma solidity 0.8.25;

import "../../../libs/ValidContract.sol";
import "../CnV3/ICnStakingV3.sol";
import "./IPublicDelegation.sol";
import "./PublicDelegationStorage.sol";
import "openzeppelin-contracts-5.0/access/Ownable.sol";
import "openzeppelin-contracts-5.0/utils/Address.sol";
import "openzeppelin-contracts-5.0/utils/math/Math.sol";
import "openzeppelin-contracts-5.0/token/ERC20/ERC20.sol";

// Implementation of public delegation of Kaia.
//
// PublicDelegation is an non-transferable interest-bearing token based on the user's staked KAIA in baseCnStakingV3.
// It must be deployed and set as reward address of CnStakingV3.
//
// The balance of PublicDelegation will be increased by the block reward, not by transaction.
// The exchange ratio between pdKAIA:KAIA keeps changing over time by the block reward.
//
// Its math is mainly based on ERC4626, which is a standard for tokenized vaults.
// See https://eips.ethereum.org/EIPS/eip-4626 for more details.

// Code organization
// - Constants
// - State variables (PublicDelegationStorage)
// - Modifiers
// - Mutators
//   - Initialize functions
//   - Operation functions
//   - Public functions
//   - Private functions
// - Private getters
// - Public getters

contract PublicDelegation is IPublicDelegation, PublicDelegationStorage, ERC20, Ownable {
    using Math for uint256;
    using Address for address payable;
    using ValidContract for address;

    /* ========== MODIFIERS ========== */

    modifier notNull(address _address) {
        require(_address != address(0), "Address is null.");
        _;
    }

    /* ========== INITIALIZE FUNCTIONS ========== */

    /// @dev Fill in initial values for the contract.
    /// @param _baseCnStakingV3 The address of the base CnStakingV3.
    /// @param _pdArgs          The constructor arguments.
    /// Emits a DeployContract event.
    constructor(
        address _baseCnStakingV3,
        PDConstructorArgs memory _pdArgs
    )
        Ownable(_pdArgs.owner)
        ERC20(
            string(abi.encodePacked(_pdArgs.gcName, " Public Delegated KAIA")),
            string(abi.encodePacked(_pdArgs.gcName, "-pdKAIA"))
        )
    {
        require(_pdArgs.commissionRate <= MAX_COMMISSION_RATE, "Commission rate is too high.");

        /// @dev Cannot check the validity of the CnStakingV3 here
        baseCnStakingV3 = ICnStakingV3(payable(_baseCnStakingV3));

        commissionRate = _pdArgs.commissionRate;
        commissionTo = _pdArgs.commissionTo;

        emit DeployContract(CONTRACT_TYPE, _baseCnStakingV3, _pdArgs);
    }

    /* ========== OPERATION FUNCTIONS ========== */

    /// @dev Update the commission receiver address.
    /// If `_commissionTo == 0x0000...0000dead`, the commission will be burned.
    /// Clear previous rewards.
    function updateCommissionTo(address _commissionTo) external override onlyOwner {
        _sweepAndStake(address(0), 0);

        address _prevCommissionTo = commissionTo;
        commissionTo = _commissionTo;

        emit UpdateCommissionTo(_prevCommissionTo, _commissionTo);
    }

    /// @dev Update the commission rate.
    /// Clear previous rewards.
    function updateCommissionRate(uint256 _commissionRate) external override onlyOwner {
        require(_commissionRate <= MAX_COMMISSION_RATE, "Commission rate is too high.");
        _sweepAndStake(address(0), 0);

        uint256 _prevCommissionRate = commissionRate;
        commissionRate = _commissionRate;

        emit UpdateCommissionRate(_prevCommissionRate, _commissionRate);
    }

    /* ========== STAKE FUNCTIONS ========== */

    /// @dev Stake KAIA to CnStakingV3.
    /// Emits a Staked event.
    function stake() external payable override {
        _sweepAndStake(_msgSender(), msg.value);
    }

    /// @dev Stake KAIA to CnStakingV3 for a specific address.
    /// @param _recipient The address to which the shares are minted.
    /// Emits a Staked event.
    function stakeFor(address _recipient) external payable override notNull(_recipient) {
        _sweepAndStake(_recipient, msg.value);
    }

    /// @dev Fallback function to receive KAIA.
    receive() external payable override {
        _sweepAndStake(_msgSender(), msg.value);
    }

    /* ========== REDELEGATE FUNCTIONS ========== */

    /// @dev Redelegate staked KAIA to another CnStakingV3 by _assets.
    /// @param _targetCnV3  The address of the target CnStakingV3.
    /// @param _assets      The amount of KAIA to redelegate.
    ///
    /// BaseCnStakingV3 will call `stakeFor` of public delegation for _targetCnV3.
    /// New pdKAIA of _targetCnV3 will be minted to same sender if _targetCnV3 uses PublicDelegation.
    ///
    /// Use `redelegateByShares` to redelegate all the KAIA.
    /// Read `redelegateByShares` for details.
    ///
    /// Emits a Redelegated event.
    function redelegateByAssets(address _targetCnV3, uint256 _assets) external override {
        require(_isRedelegationEnabled(), "Redelegation disabled.");
        require(_targetCnV3._validCnStaking(3), "Invalid CnStakingV3.");

        _sweepAndStake(address(0), 0);

        uint256 _shares = previewWithdraw(_assets);
        _burn(_msgSender(), _shares);

        _redelegate(_targetCnV3, _assets);
    }

    /// @dev Redelegate staked KAIA to another CnStakingV3 by _shares.
    /// @param _targetCnV3  The address of the target CnStakingV3.
    /// @param _shares      The amount of shares to redelegate.
    ///
    /// BaseCnStakingV3 will call `stakeFor` of public delegation for _targetCnV3.
    /// New pdKAIA of _targetCnV3 will be minted to same sender if _targetCnV3 uses PublicDelegation.
    ///
    /// Users must use `byShares` to redelegate all the KAIA
    /// since `byAssets` might leave a dust of shares due to block reward before tx arrived.
    ///
    /// Emits a Redelegated event.
    function redelegateByShares(address _targetCnV3, uint256 _shares) external override {
        require(_isRedelegationEnabled(), "Redelegation disabled.");
        require(_targetCnV3._validCnStaking(3), "Invalid CnStakingV3.");

        _sweepAndStake(address(0), 0);

        uint256 _assets = previewRedeem(_shares);
        _burn(_msgSender(), _shares);

        _redelegate(_targetCnV3, _assets);
    }

    /* ========== WITHDRAW FUNCTIONS ========== */

    /// @dev Request withdrawal of staked KAIA by assets.
    /// @param _recipient The address to which the KAIA is sent.
    /// @param _assets    The amount of KAIA to redeem.
    ///
    /// Since unstaking KAIA is not eligible for rewards, shares need to be burned before actual withdrawal.
    /// Users can withdraw KAIA after the 7-day withdrawal period.
    ///
    /// Use `redeem` to withdraw all the KAIA.
    /// Read `redeem` for details.
    ///
    /// Emits a RequestWithdrawal event.
    /// Emits a Redeemed event.
    function withdraw(address _recipient, uint256 _assets) external override notNull(_recipient) {
        _sweepAndStake(address(0), 0);

        uint256 _shares = previewWithdraw(_assets);

        _burn(_msgSender(), _shares);
        _requestWithdrawal(_msgSender(), _recipient, _assets);

        emit Redeemed(_msgSender(), _recipient, _assets, _shares);
    }

    /// @dev Request withdrawal of staked KAIA by shares.
    /// @param _recipient The address to which the KAIA is sent.
    /// @param _shares    The amount of shares to redeem considering PRECISION.
    ///
    /// Since unstaking KAIA is not eligible for rewards, shares need to be burned before actual withdrawal.
    /// Users can withdraw KAIA after the 7-day withdrawal period.
    ///
    /// Users must use `redeem` to withdraw all the KAIA
    /// since `withdraw` might leave a dust of shares due to block reward before tx arrived.
    ///
    /// Emits a RequestWithdrawal event.
    /// Emits a Redeemed event.
    function redeem(address _recipient, uint256 _shares) external override notNull(_recipient) {
        _sweepAndStake(address(0), 0);

        uint256 _assets = previewRedeem(_shares);

        _burn(_msgSender(), _shares);
        _requestWithdrawal(_msgSender(), _recipient, _assets);

        emit Redeemed(_msgSender(), _recipient, _assets, _shares);
    }

    /* ========== CANCEL/CLAIM WITHDRAWAL FUNCTIONS ========== */

    /// @dev Cancel the approved staking withdrawal.
    /// @param _requestId The ID of the approved staking withdrawal.
    ///
    /// Revive the shares and cancel the approved staking withdrawal.
    ///
    /// Emits a CancelApprovedStakingWithdrawal event.
    function cancelApprovedStakingWithdrawal(uint256 _requestId) external override {
        _sweepAndStake(address(0), 0);

        _cancelApprovedStakingWithdrawal(_requestId);
    }

    /// @dev Claim the approved staking withdrawal if the withdrawal period is over.
    /// @param _requestId The ID of the approved staking withdrawal.
    ///
    /// If the withdrawal was canceled, revive the shares.
    ///
    /// Emits a CancelApprovedStakingWithdrawal event if the withdrawal was canceled.
    /// Emits a Claimed event.
    function claim(uint256 _requestId) external override {
        _sweepAndStake(address(0), 0);

        _claim(_requestId);
    }

    /* ========== TRANSFER FUNCTIONS ========== */

    /// @dev Transfer is not allowed.
    function _update(address _from, address _to, uint256 _value) internal override {
        /// @dev Since _transfer already filters the below condition,
        /// it can prevent the transfer between the addresses except for the mint and burn.
        require(_from == address(0) || _to == address(0), "Transfer not allowed.");

        super._update(_from, _to, _value);
    }

    /* ========== UTIL FUNCTION ========== */

    /// @dev Manually sweep the rewards.
    function sweep() external override {
        _sweepAndStake(address(0), 0);
    }

    /* ========== PRIVATE CALLBACKS ========== */

    /// @dev Sweep rewards and stake KAIA to CnStakingV3.
    /// @param _recipient The address to which the shares are minted.
    ///
    /// The execution order is matter for the security of the contract.
    /// For example, in `previewRedeem`, the withdrawable KAIA reflects the cumulated rewards until the last block.
    /// But if the rewards haven't been swept to CnStakingV3 before requesting withdrawal, it might be reverted due the lack of KAIA.
    ///
    /// Emits a SendCommission event.
    /// Emits a Staked event if it's not just sweeping.
    function _sweepAndStake(address _recipient, uint256 _assets) private {
        unchecked {
            uint256 _reward = address(this).balance - _assets;
            uint256 _commission = _calcCommission(_reward);
            if (_commission > 0) {
                _sendCommission(_commission);
            }

            if (_recipient != address(0)) {
                /// @dev Can't use _totalAssets() here since it contains the user's staking amount.
                uint256 _baseAssets = _totalStaking() + _reward - _commission;
                uint256 _shares = _convertToShares(_assets, totalSupply(), _baseAssets);
                require(_shares > 0, "Stake amount is too low.");
                _mint(_recipient, _shares);
                emit Staked(_recipient, _assets, _shares);
            }

            uint256 _toStake = _reward + _assets - _commission;
            if (_toStake > 0) {
                _delegate(_toStake);
            }
        }
    }

    /// @dev Delegate staked KAIA to CnStakingV3.
    function _delegate(uint256 _amount) private {
        baseCnStakingV3.delegate{value: _amount}();
    }

    /// @dev Redelegate staked KAIA to another CnStakingV3.
    function _redelegate(address _targetCnV3, uint256 _assets) private {
        require(_assets > 0, "Redelegate amount is too low.");

        baseCnStakingV3.redelegate(_msgSender(), _targetCnV3, _assets);

        emit Redelegated(_msgSender(), _targetCnV3, _assets);
    }

    /// @dev Send the approved staking withdrawal.
    function _requestWithdrawal(address _owner, address _recipient, uint256 _assets) private {
        require(_assets > 0, "Withdrawal amount is too low.");

        uint256 _id = baseCnStakingV3.approveStakingWithdrawal(_recipient, _assets);
        userRequestIds[_owner].push(_id);
        requestIdToOwner[_id] = _owner;

        emit RequestWithdrawal(_owner, _recipient, _id, _assets);
    }

    /// @dev Cancel the approved staking withdrawal.
    function _cancelApprovedStakingWithdrawal(uint256 _id) private {
        require(requestIdToOwner[_id] == _msgSender(), "Not the owner of the request.");

        // Revive the shares
        (, uint256 _assets, , ) = _withdrawalRequest(_id);
        uint256 _shares = previewDeposit(_assets);
        _mint(_msgSender(), _shares);

        baseCnStakingV3.cancelApprovedStakingWithdrawal(_id);

        emit RequestCancelWithdrawal(_msgSender(), _id);
    }

    /// @dev Claim the approved staking withdrawal.
    function _claim(uint256 _id) private {
        require(requestIdToOwner[_id] == _msgSender(), "Not the owner of the request.");

        baseCnStakingV3.withdrawApprovedStaking(_id);

        (, uint256 _asset, , ICnStakingV3.WithdrawalStakingState _state) = _withdrawalRequest(_id);

        // If the withdrawal was canceled, revive the shares.
        // Since the `_totalAsset()` was already increased by _asset, it needs to be subtracted.
        if (_state == ICnStakingV3.WithdrawalStakingState.Canceled) {
            uint256 _shares = _convertToShares(_asset, totalSupply(), _totalAssets() - _asset);
            _mint(_msgSender(), _shares);
            return;
        }

        emit Claimed(_msgSender(), _id);
    }

    /// @dev Send the commission to the commission address.
    function _sendCommission(uint256 _commission) private {
        payable(commissionTo).sendValue(_commission);

        emit SendCommission(commissionTo, _commission);
    }

    /* ========== PRIVATE GETTERS ========== */

    function _isRedelegationEnabled() private view returns (bool) {
        return baseCnStakingV3.isRedelegationEnabled();
    }

    /// @dev Get the approved staking withdrawal info.
    function _withdrawalRequest(
        uint256 _id
    ) private view returns (address, uint256, uint256, ICnStakingV3.WithdrawalStakingState) {
        return baseCnStakingV3.getApprovedStakingWithdrawalInfo(_id);
    }

    /// @dev Get the total asset of the contract.
    /// Note that it can't be used to calculate user's share in any payable function
    /// since _pureReward() will temporarily contain the user's staking amount.
    function _totalAssets() private view returns (uint256) {
        /// @dev Total asset is the sum of staking and reward except commission.
        unchecked {
            return _totalStaking() + _pureReward();
        }
    }

    /// @dev Get the total staking amount.
    function _totalStaking() private view returns (uint256) {
        unchecked {
            return baseCnStakingV3.staking() - baseCnStakingV3.unstaking();
        }
    }

    /// @dev Calculate the commission.
    /// @param _amount   The amount to calculate the commission.
    function _calcCommission(uint256 _amount) private view returns (uint256) {
        return _amount.mulDiv(commissionRate, COMMISSION_DENOMINATOR, Math.Rounding.Floor);
    }

    /// @dev Get the current reward amount.
    function _totalReward() private view returns (uint256) {
        return address(this).balance;
    }

    /// @dev Get the pure reward amount.
    function _pureReward() private view returns (uint256) {
        uint256 _reward = _totalReward();
        uint256 _commission = _calcCommission(_reward);
        unchecked {
            return _reward - _commission;
        }
    }

    /// @dev Convert assets to shares with rounding down with custom totalSupply and totalAssets.
    function _convertToShares(
        uint256 _assets,
        uint256 _customSupply,
        uint256 _customAssets
    ) private pure returns (uint256) {
        return _customSupply == 0 ? _assets : _assets.mulDiv(_customSupply, _customAssets, Math.Rounding.Floor);
    }

    /* ========== PUBLIC GETTERS ========== */

    /// @dev Get the current state of withdrawal request.
    /// It will be useful since CnStakingV3's `getApprovedStakingWithdrawalInfo` only returns static info.
    function getCurrentWithdrawalRequestState(uint256 _id) public view override returns (WithdrawalRequestState) {
        (, , uint256 _withdrawableFrom, ICnStakingV3.WithdrawalStakingState _state) = _withdrawalRequest(_id);

        if (_withdrawableFrom == 0) {
            return WithdrawalRequestState.Undefined;
        }

        if (_state == ICnStakingV3.WithdrawalStakingState.Canceled) {
            return WithdrawalRequestState.Canceled;
        }

        if (_state == ICnStakingV3.WithdrawalStakingState.Transferred) {
            return WithdrawalRequestState.Withdrawn;
        }

        uint256 _withdrawableUntil;
        unchecked {
            _withdrawableUntil = _withdrawableFrom + baseCnStakingV3.STAKE_LOCKUP();
        }
        if (block.timestamp < _withdrawableFrom) {
            return WithdrawalRequestState.Requested;
        } else if (block.timestamp < _withdrawableUntil) {
            return WithdrawalRequestState.Withdrawable;
        }

        return WithdrawalRequestState.PendingCancel;
    }

    /// @dev Get the user's all request count.
    function getUserRequestCount(address _owner) public view override returns (uint256) {
        return userRequestIds[_owner].length;
    }

    /// @dev Get the user's request IDs with the specific state.
    function getUserRequestIdsWithState(
        address _owner,
        WithdrawalRequestState _state
    ) public view override returns (uint256[] memory _requestIds) {
        uint256 _count = userRequestIds[_owner].length;
        _requestIds = new uint256[](_count);
        uint256 _index = 0;

        for (uint256 i = 0; i < _count; i++) {
            uint256 _id = userRequestIds[_owner][i];
            if (getCurrentWithdrawalRequestState(_id) == _state) {
                _requestIds[_index++] = _id;
            }
        }

        assembly {
            mstore(_requestIds, _index)
        }

        return _requestIds;
    }

    /// @dev Get the user's all request IDs.
    function getUserRequestIds(address _owner) public view override returns (uint256[] memory) {
        return userRequestIds[_owner];
    }

    /// @dev Get the total asset of the contract.
    function totalAssets() public view override returns (uint256) {
        return _totalAssets();
    }

    /// @dev Get the maximum redeemable shares of the owner.
    function maxRedeem(address _owner) public view override returns (uint256) {
        return super.balanceOf(_owner);
    }

    /// @dev Get the maximum withdrawable KAIA of the owner (= total assets owned by the owner).
    function maxWithdraw(address _owner) public view override returns (uint256) {
        return previewRedeem(balanceOf(_owner));
    }

    /// @dev Get the reward amount except commission.
    function reward() public view override returns (uint256) {
        return _pureReward();
    }

    /// @dev Convert assets to shares with rounding down.
    function convertToShares(uint256 _assets) public view returns (uint256) {
        uint256 supply = totalSupply();

        return supply == 0 ? _assets : _assets.mulDiv(supply, totalAssets(), Math.Rounding.Floor);
    }

    /// @dev Convert shares to assets with rounding down.
    function convertToAssets(uint256 _shares) public view returns (uint256) {
        uint256 supply = totalSupply();

        return supply == 0 ? _shares : _shares.mulDiv(totalAssets(), supply, Math.Rounding.Floor);
    }

    /// @dev Preview the expected shares to mint for the deposit amount.
    function previewDeposit(uint256 _assets) public view returns (uint256) {
        return convertToShares(_assets);
    }

    /// @dev Preview the expected shares to redeem for the withdrawal amount.
    function previewWithdraw(uint256 _assets) public view returns (uint256) {
        uint256 supply = totalSupply();

        return supply == 0 ? _assets : _assets.mulDiv(supply, totalAssets(), Math.Rounding.Ceil);
    }

    /// @dev Preview the expected assets to redeem for the shares.
    function previewRedeem(uint256 _shares) public view returns (uint256) {
        return convertToAssets(_shares);
    }
}
