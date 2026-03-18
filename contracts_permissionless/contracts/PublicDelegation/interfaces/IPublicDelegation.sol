// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

/// @title IKIP163
/// @notice Standard staking interface for redelegation handoff.
interface IKIP163 {
    /// @notice Stake KAIA on behalf of a recipient.
    /// @param _recipient The address to receive pdKAIA shares.
    function stakeFor(address _recipient) external payable;

    /// @notice Get the pure reward amount (balance minus commission).
    function reward() external view returns (uint256);
}

/// @title IPublicDelegation
/// @notice Interface for PublicDelegation V2 — a tokenized staking vault (pdKAIA).
interface IPublicDelegation is IKIP163 {
    /* ========== STRUCTS ========== */

    struct PDConstructorArgs {
        address owner;
        address commissionTo;
        uint256 commissionRate;
        string gcName;
    }

    /* ========== ENUMS ========== */

    enum WithdrawalRequestState {
        Undefined,
        Requested,
        Withdrawable,
        Withdrawn,
        PendingCancel,
        Canceled
    }

    /* ========== EVENTS ========== */

    event DeployPublicDelegation(string contractType, address indexed baseCnStaking, PDConstructorArgs pdArgs);
    event UpdateCommissionTo(address indexed prevCommissionTo, address indexed commissionTo);
    event UpdateCommissionRate(uint256 indexed prevCommissionRate, uint256 indexed commissionRate);
    event SendCommission(address indexed commissionTo, uint256 commission);
    event Staked(address indexed user, uint256 assets, uint256 shares);
    event Redeemed(address indexed user, address indexed recipient, uint256 assets, uint256 shares);
    event Redelegated(address indexed user, address indexed targetCnStaking, uint256 assets);
    event RequestWithdrawal(address indexed user, address indexed recipient, uint256 indexed requestId, uint256 assets);
    event RequestCancelWithdrawal(address indexed user, uint256 indexed requestId);
    event Claimed(address indexed user, uint256 indexed requestId);

    /* ========== ERRORS ========== */

    error ZeroAddress();
    error CommissionRateTooHigh();
    error StakeAmountTooLow();
    error WithdrawalAmountTooLow();
    error RedelegateAmountTooLow();
    error NotRequestOwner();

    /* ========== INITIALIZATION ========== */

    /// @notice Initialize the PublicDelegation contract.
    /// @param _baseCnStaking The CnStaking proxy address this PD is linked to.
    /// @param _args          Constructor arguments (owner, commissionTo, commissionRate, gcName).
    function initialize(address _baseCnStaking, PDConstructorArgs memory _args) external;

    /* ========== CONSTANTS ========== */

    /// @notice The contract type identifier.
    function CONTRACT_TYPE() external pure returns (string memory);

    /// @notice The contract version.
    function VERSION() external pure returns (uint256);

    /// @notice Maximum commission rate (1e4 = 100%).
    function MAX_COMMISSION_RATE() external pure returns (uint256);

    /// @notice Commission rate denominator.
    function COMMISSION_DENOMINATOR() external pure returns (uint256);

    /* ========== OWNER FUNCTIONS ========== */

    /// @notice Update the commission recipient. Sweeps rewards with old settings first.
    /// @param _commissionTo The new commission recipient address.
    function updateCommissionTo(address _commissionTo) external;

    /// @notice Update the commission rate. Sweeps rewards with old rate first.
    /// @param _commissionRate The new commission rate.
    function updateCommissionRate(uint256 _commissionRate) external;

    /* ========== STAKING FUNCTIONS ========== */

    /// @notice Stake KAIA and receive pdKAIA shares.
    function stake() external payable;
    // stakeFor inherited from IKIP163

    /// @notice Receive KAIA and stake.
    receive() external payable;

    /* ========== WITHDRAWAL FUNCTIONS ========== */

    /// @notice Withdraw exact KAIA amount. Burns ceiling shares.
    /// @param _recipient The address to receive the withdrawn KAIA.
    /// @param _assets    The exact amount of KAIA to withdraw.
    function withdraw(address _recipient, uint256 _assets) external;

    /// @notice Redeem exact shares for KAIA. Receives floor assets.
    /// @param _recipient The address to receive the redeemed KAIA.
    /// @param _shares    The exact number of shares to redeem.
    function redeem(address _recipient, uint256 _shares) external;

    /// @notice Cancel a withdrawal request. Re-mints shares at current rate.
    /// @param _requestId The withdrawal request ID to cancel.
    function cancelApprovedStakingWithdrawal(uint256 _requestId) external;

    /// @notice Claim a withdrawable request, or auto-cancel if expired.
    /// @param _requestId The withdrawal request ID to claim.
    function claim(uint256 _requestId) external;

    /* ========== REDELEGATION FUNCTIONS ========== */

    /// @notice Redelegate by exact KAIA amount. Burns ceiling shares.
    /// @param _targetCnStaking The target CnStaking contract to redelegate to.
    /// @param _assets          The exact amount of KAIA to redelegate.
    function redelegateByAssets(address _targetCnStaking, uint256 _assets) external;

    /// @notice Redelegate by exact shares. Sends floor KAIA.
    /// @param _targetCnStaking The target CnStaking contract to redelegate to.
    /// @param _shares          The exact number of shares to redelegate.
    function redelegateByShares(address _targetCnStaking, uint256 _shares) external;

    /* ========== SWEEP ========== */

    /// @notice Manually trigger reward sweep and stake.
    function sweep() external;

    /* ========== GETTERS ========== */

    /// @notice The base CnStaking contract this PD is linked to.
    function baseCnStaking() external view returns (address);

    /// @notice The commission recipient address.
    function commissionTo() external view returns (address);

    /// @notice The commission rate.
    function commissionRate() external view returns (uint256);

    /// @notice Total assets under management (staking net + pure reward).
    function totalAssets() external view returns (uint256);

    /// @notice Convert KAIA assets to pdKAIA shares (floor).
    /// @param _assets The amount of KAIA to convert.
    function convertToShares(uint256 _assets) external view returns (uint256);

    /// @notice Convert pdKAIA shares to KAIA assets (floor).
    /// @param _shares The number of shares to convert.
    function convertToAssets(uint256 _shares) external view returns (uint256);

    /// @notice Preview the number of shares minted for a deposit.
    /// @param _assets The amount of KAIA to deposit.
    function previewDeposit(uint256 _assets) external view returns (uint256);

    /// @notice Preview the number of shares burned for a withdrawal (ceiling).
    /// @param _assets The amount of KAIA to withdraw.
    function previewWithdraw(uint256 _assets) external view returns (uint256);

    /// @notice Preview the amount of KAIA received for redeeming shares (floor).
    /// @param _shares The number of shares to redeem.
    function previewRedeem(uint256 _shares) external view returns (uint256);

    /// @notice Maximum number of shares redeemable by an owner.
    /// @param _owner The owner address.
    function maxRedeem(address _owner) external view returns (uint256);

    /// @notice Maximum KAIA withdrawable by an owner.
    /// @param _owner The owner address.
    function maxWithdraw(address _owner) external view returns (uint256);

    /// @notice Get a user's withdrawal request ID at a given index.
    /// @param _owner  The user address.
    /// @param _index  The index in the user's request list.
    function userRequestIds(address _owner, uint256 _index) external view returns (uint256);

    /// @notice Get the owner of a withdrawal request.
    /// @param _requestId The withdrawal request ID.
    function requestIdToOwner(uint256 _requestId) external view returns (address);

    /// @notice Get the current state of a withdrawal request.
    /// @param _requestId The withdrawal request ID.
    function getCurrentWithdrawalRequestState(uint256 _requestId) external view returns (WithdrawalRequestState);

    /// @notice Get the total number of withdrawal requests for a user.
    /// @param _owner The user address.
    function getUserRequestCount(address _owner) external view returns (uint256);

    /// @notice Get a user's withdrawal request IDs filtered by state.
    /// @param _owner  The user address.
    /// @param _state  The withdrawal request state to filter by.
    function getUserRequestIdsWithState(
        address _owner,
        WithdrawalRequestState _state
    ) external view returns (uint256[] memory);

    /// @notice Get all withdrawal request IDs for a user.
    /// @param _owner The user address.
    function getUserRequestIds(address _owner) external view returns (uint256[] memory);
}
