// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

/// @title ICnStaking
/// @notice Version-agnostic interface for CnStaking contracts (V4+).
interface ICnStaking {
    /* ========== STRUCTS ========== */

    struct WithdrawalRequest {
        address to;
        uint256 value;
        uint256 withdrawableFrom;
        WithdrawalStakingState state;
    }

    /* ========== ENUMS ========== */

    enum WithdrawalStakingState {
        Unknown,
        Transferred,
        Canceled
    }

    /* ========== EVENTS ========== */

    event DeployCnStakingV4(string contractType, address indexed publicDelegation);

    event DelegateKaia(address indexed from, uint256 value);

    event ApproveStakingWithdrawal(
        uint256 indexed approvedWithdrawalId,
        address indexed to,
        uint256 value,
        uint256 withdrawableFrom
    );

    event CancelApprovedStakingWithdrawal(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value);

    event WithdrawApprovedStaking(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value);

    event Redelegation(address indexed user, address indexed targetCnStaking, uint256 value);

    event HandleRedelegation(
        address indexed user,
        address indexed prevCnStaking,
        address indexed targetCnStaking,
        uint256 value
    );

    /* ========== ERRORS ========== */

    error ZeroAddress();
    error ZeroValue();
    error NotStaker();
    error NotUnstakingManager();
    error RedelegationDisabled();
    error InvalidValue();
    error InvalidTarget();
    error RedelegationCooldown();
    error WithdrawalNotFound();
    error InvalidWithdrawalState();
    error NotWithdrawableYet();
    error TransferFailed();
    error InvalidStakeFor();
    error BaseCnStakingMismatch();

    /* ========== INITIALIZATION FUNCTIONS ========== */

    /// @notice Initialize the CnStaking contract without PublicDelegation.
    /// @param _owner   The owner of this CnStaking (can be an EOA or multisig).
    function initialize(address _owner) external;

    /// @notice Initialize the CnStaking contract with PublicDelegation.
    /// @param _owner             The owner of this CnStaking (can be an EOA or multisig).
    /// @param _publicDelegation  The PublicDelegation proxy address (already deployed by factory).
    function initializeWithPD(address _owner, address _publicDelegation) external;

    /* ========== STAKING FUNCTIONS ========== */

    /// @notice Delegate KAIA to this contract.
    function delegate() external payable;

    /// @notice Receive KAIA and add to delegated stake.
    receive() external payable;

    /* ========== UNSTAKING FUNCTIONS ========== */

    /// @notice Schedule a withdrawal of delegated stakes.
    /// @param _to    Recipient address.
    /// @param _value Amount of KAIA to withdraw.
    /// @return id    The withdrawal request ID.
    function approveStakingWithdrawal(address _to, uint256 _value) external returns (uint256 id);

    /// @notice Cancel an approved withdrawal request.
    /// @param _id The withdrawal request ID.
    function cancelApprovedStakingWithdrawal(uint256 _id) external;

    /// @notice Execute or auto-cancel an approved withdrawal.
    /// @param _id The withdrawal request ID.
    function withdrawApprovedStaking(uint256 _id) external;

    /* ========== REDELEGATION FUNCTIONS ========== */

    /// @notice Redelegate stakes to another CnStaking contract.
    /// @param _user             The user initiating the redelegation (managed by PD).
    /// @param _targetCnStaking  The target CnStaking contract.
    /// @param _value            The amount of KAIA to redelegate.
    function redelegate(address _user, address _targetCnStaking, uint256 _value) external;

    /// @notice Handle incoming re-delegation from another CnStaking contract.
    /// @param _user The user receiving the redelegated stake.
    function handleRedelegation(address _user) external payable;

    /* ========== GETTERS ========== */

    /// @notice The contract type identifier.
    function CONTRACT_TYPE() external pure returns (string memory);

    /// @notice The contract version.
    function VERSION() external pure returns (uint256);

    /// @notice The stake lockup duration.
    function STAKE_LOCKUP() external pure returns (uint256);

    /// @notice The PublicDelegation contract address (address(0) if PD disabled).
    function publicDelegation() external view returns (address);

    /// @notice Total delegated stake held by this contract.
    function staking() external view returns (uint256);

    /// @notice Total amount pending in approved withdrawal requests.
    function unstaking() external view returns (uint256);

    /// @notice Monotonically increasing withdrawal request counter.
    function withdrawalRequestCount() external view returns (uint256);

    /// @notice Per-user last redelegation timestamp.
    /// @param _account The user address.
    function lastRedelegation(address _account) external view returns (uint256);

    /// @notice Query withdrawal request IDs matching a given state.
    /// @param _from  Search begin index.
    /// @param _to    Search end index (0 or >= requestCount means search till end).
    /// @param _state Withdrawal state to filter by.
    /// @return ids   Withdrawal IDs satisfying the conditions.
    function getApprovedStakingWithdrawalIds(
        uint256 _from,
        uint256 _to,
        WithdrawalStakingState _state
    ) external view returns (uint256[] memory ids);

    /// @notice Query a withdrawal request's details.
    /// @param _index Withdrawal request ID.
    /// @return to                Recipient address.
    /// @return value             Withdrawal amount.
    /// @return withdrawableFrom  Earliest withdrawal timestamp.
    /// @return state             Request state.
    function getApprovedStakingWithdrawalInfo(
        uint256 _index
    ) external view returns (address to, uint256 value, uint256 withdrawableFrom, WithdrawalStakingState state);
}
