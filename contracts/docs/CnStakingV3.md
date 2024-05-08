# CnStakingV3

A member of Kaia Governance Council (GC) must own at least one CnStaking contract. There are three versions of CnStaking contract so far.

- [Deprecated] legacy/CnStakingContract (CnSV1): The original CnStaking contract that is being used since the start of Kaia mainnet.
- CnStakingV2 (CnSV2): An upgraded version of CnSV1, having all interfaces of V1 as well as new governance-related features.
- CnStakingV3MultiSig (CnSV3): An upgraded version of CnSV2, having all features of V2 as well as new public delegation features.

This document will describe the V3 contract.

## Code

- [ICnStakingV3.sol](../contracts/consensus/CnV3/ICnStakingV3.sol)
- [ICnStakingV3MultiSig.sol](../contracts/consensus/CnV3/ICnStakingV3MultiSig.sol)
- [CnStakingV3Storage.sol](../contracts/consensus/CnV3/CnStakingV3Storage.sol)
- [CnStakingV3MultiSigStorage.sol](../contracts/consensus/CnV3/CnStakingV3MultiSigStorage.sol)
- [CnStakingV3.sol](../contracts/consensus/CnV3/CnStakingV3.sol)
- [CnStakingV3MultiSig.sol](../contracts/consensus/CnV3/CnStakingV3MultiSig.sol)

## CnStakingV3MultiSig and CnStakingV3

The CnStakingV3MultiSig contract inherits the CnStakingV3 contract to add multi-sig features. We intentionally separated the multi-sig features from the main CnStaking contract to make the contract more modular and easier to maintain.

## Operation Modes

The CnSV3 contract has two operation modes:

1. **Disable Public Delegation**: Do not enable public delegation. In this mode, the contract behaves same as the CnSV2 contract.
2. **Enable Public Delegation**: Enable public delegation. In this mode, the general users can delegate their KAIA to the CnSV3 contract.

The public delegation can be enabled by passing zero reward address and initial lockup conditions to the constructor. The constructor will set the `isPublicDelegationEnabled` flag to true if the reward address is zero. The reward address must be set to public delegation contract address in `setPublicDelegation()` function. Since the reward address must not be updated, the `multisig UpdateRewardAddress` function won't be available.

```solidity
    // In constructor
    // Check if the initial conditions are valid.
    // If `_rewardaddress` is zero, `_unlockTime` and `_unlockAmount` must be empty.
    _validInitialConditions(_rewardAddress, _unlockTime, _unlockAmount);

    isPublicDelegationEnabled = _rewardAddress == address(0);
```

The initial lockup conditions must be empty since all KAIA must be delegated to CnSV3 through the public delegation contract if enabled.

## Access Control

The all access control is managed by the `openzeppelin/AccessControlEnumerable.sol` library. The CnSV3 contract has the following roles:

- `OPERATOR_ROLE`: The role that can call the multi-sig execution function. Usually set to contract itself.
- `ADMIN_ROLE`: The role that can submit a new multi-sig request.
- `STAKER_ROLE`: The role that can delegate KAIA to CnSV3 contract.
- `UNSTAKING_APPROVER_ROLE`: The role that can approve the unstaking request.
- `UNSTAKING_CLAIMER_ROLE`: The role that can claim the claimable unstaking request.
- `contractValidator, or CV`: The CV is a temporary admin who is involved in contract initialization. In the current GC onboarding process, the Kaia team takes the role. The CV is temporarily assigned to the `ADMIN_ROLE` and removed in the initialization process.

In the **Enable Public Delegation** mode, the staking functions are controlled by `Public Delegation (PD)` contract as follows:

| Functions                               | PD disabled    | PD enabled |
| --------------------------------------- | -------------- | ---------- |
| `delegate/receive`                      | No condition   | Only PD    |
| `submitApproveStakingWithdrawal`        | Only Admin     | N.A.       |
| `submitCancelApprovedStakingWithdrawal` | Only Admin     | N.A.       |
| `approveStakingWithdrawal`              | Only Admin     | Only PD    |
| `cancelApprovedStakingWithdrawal`       | Only Multi-sig | Only PD    |
| `withdrawApprovedStaking`               | Only Multi-sig | Only PD    |

## Contract Initialization

Initializing a CnSV3 takes multiple steps involving many entities.

1. `constructor()`, usually by the contractValidator, or CV.
2. `setGCId()`, by the CV, to specify the GC ID.
3. `setStakingTracker()` by the CV or any admin, to specify the StakingTracker address.
4. `setPublicDelegation()` (when PD enabled) by the CV or any admin, to deploy the PD and set its address as the reward address.
5. `reviewInitialConditions()` by every admins and the CV, to signify that they all agree with initial lockup conditions.
6. `depositLockupStakingAndInit()` to deposit the initial lockup KAIAs if there's any lockup and finish the initialization process.

## Multisig facility

Each multisig function call is stored in a Request struct. Upto three arguments are encoded as bytes32 type each.

```solidity
struct Request {
  Functions functionId;
  bytes32 firstArg;
  bytes32 secondArg;
  bytes32 thirdArg;
  address requestProposer;
  address[] confirmers;
  RequestState state;
}
```

A request can be in following `RequestState`

- `NotConfirmed`: before execution.
- `Executed`: successfully executed.
- `ExecutionFailed`: did call the function but failed.
- `Canceled`: canceled before execution.

A multisig function "X" is executed through following steps.

An admin calls `submitX()` function, creating a `Request`. Any preconditions (e.g. cannot add already admin) are checked at this time.
The `Request` content is available from either `SubmitRequest` event or `getRequestIds()`, `getRequestInfo()` functions.
For zero or more times, other admins call confirmRequest()
The function `x()` is executed when the requirement number of admins have called either `submitX()` or `confirmRequest()`. Any preconditions are checked again at this time because states may have changed since request creation.
The function call to `x()` never reverts as `.call()` is used. However, depending on the call result either `ExecuteRequestSuccess` or `ExecuteRequestFailure` event is emitted.

For simplicity, calling "multisig X()" refers to all above steps.

## Features

CnSV3 is a multi-purpose contract. Its features can be categorized into 4 groups.

- Admin management
- Initial lockup stakes
- Delegation
- Re-delegation
- External accounts management

### Admin management

The contract admins themselves are managed by multisig functions.

Related functions:

- `multisig AddAdmin(_admin)`: Add an account to `ADMIN_ROLE`.
- `multisig DeleteAdmin(_admin)`: Delete an account from `ADMIN_ROLE`.
- `multisig UpdateRequirement(_requirement)`: Change the requirement number.
- `multisig ClearRequest()`: Cancels all outstanding requests.

As a precaution, all outstanding (i.e. `RequestState.NotConfirmed`) requests are canceled when either `ADMIN_ROLE` or `requirement` changes.

### Initial lockup stakes

Initial lockup is a series of long-term fixed lockups. The lockup condition is stored in a singleton LockupConditions struct. Total amounts are tracked by initialLockupStaking and remainingLockupStaking.

```solidity
LockupConditions public lockupConditions;
uint256 public initialLockupStaking;
uint256 public remainingLockupStaking;
struct LockupConditions {
    uint256[] unlockTime;
    uint256[] unlockAmount;
    bool allReviewed;
    uint256 reviewedCount;
    mapping(address => bool) reviewedAdmin;
}
```

Related functions:

- `reviewInitialConditions()`: Agree to initial lockup conditions during contract init.
- `depositLockupStakingAndInit()`: Deposit the sum of initial lockup amounts during contract init.
- `multisig WithdrawLockupStaking(to, value)`: Withdraw all or part of lockup stakes.

### Delegation

The additioanl stakes can be added regardless of the operation mode, while there're different conditions to call as we already mentioned.

There's a 7 days of lockup period for the delegation.

```solidity
uint256 public staking;
uint256 public unstaking;
uint256 public withdrawalRequestCount;
mapping(uint256 => WithdrawalRequest) private withdrawalRequestMap;
struct WithdrawalRequest {
    address to;
    uint256 value;
    uint256 withdrawableFrom;
    WithdrawalStakingState state;
}
```

A withdrawal request can be in following `WithdrawalStakingState`

- `Unknown`: Before transfer or cancellation. Unlike `RequestState.Unknown`, `WithdrawalStakingState.Unknown` is a meaningful enum.
- `Transferred`: Successfully executed.
- `Canceled`: Canceled due to timeout or an explicit cancellation.

A withdrawal request becomes withdrawable from 7 days since creation (D+7). However it expires after another 7 days (D+14).

At any moment, sum of unstaking (pending withdrawal) amount cannot exceed the staking amount. This rule prevents an abusing scenario where a GC member "spray" withdrawal requests across different days and effectively bypass the 7-day delay.

Related functions:

- `delegate()`: Delegates KAIA to the contract.
- `receive()`: Equivalent to delegate().
- `multisig or PD ApproveStakingWithdrawal(to, value)`: Create a `WithdrawalRequest` to take out all or part of non-lockup stakes.
- `multisig or PD CancelApprovedStakingWithdrawal(id)`: Explicitly cancel a `WithdrawalRequest`.
- `withdrawApprovedStaking(id)`: Execute a `WithdrawalRequest` after 7 days of lockup.

### Re-delegation

The redelegation is a togglable feature, and both current CnSV3 and target CnSV3 must enable PD and redelegation to use it. Users can redelegate their KAIA to another CnSV3 without waiting for lockup period. But to prevent redelegation hopping, which makes KAIA's consensus be unstable, users who have been redelegated must wait a lockup period before they can redelegate again. For example, if a user redelegates from [`A` → `B`], then user must wait for a lockup period to redelegate from `B` to another CnSV3. The last redelegation records will be stored in each CnSV3s. As same as delegation, the re-delegation needs to be initiated by the PD contract.

```solidity
    bool public isRedelegationEnabled;
    // Record last delegation time
    mapping(address => uint256) public lastRedelegation;
```

Related functions:

- `redelegate(user, targetCnV3, value)`: Redelegate the `pdKAIA` to target CnSV3. Call `handleRedelegation` of the target CnSV3.
- `handleRedelegation(user)`: Handle the re-delegation request from the previous CnSV3.

![Redelegation flow](./assets/Redelegation.png)

To process the re-delegation in one transaction, the target CnsV3 will delegate KAIA on behalf of `user` by calling `stakeFor(user)` of the PD contract. The PD will delegate the KAIA back to the target CnSV3. To prevent the loss of KAIA, the target CnSV3 calculates and checks the expected balance of `KAIA` in `handleRedelegation`.

### External accounts management

A KAIA GC member uses several accounts for different purposes. Some of them are appointed via CnSV3 contract.

Related functions:

- `multisig UpdateRewardAddress(addr)`: Update the `rewardAddress` and also update to the `AddressBook` contract. Disabled if public delegation is enabled.
- `multisig UpdateStakingTracker(addr)`: Update the `stakingTracker` address. CnSV3 contract notifies staking balance change and voter address change to the `StakingTracker`.
- `multisig UpdateVoterAddress(addr)`: Update the `voterAddress` and also update to the `StakingTracker` contract.

## Interaction with other contracts

### AddressBook

The AddressBook contract (0x000..400) is the central directory of Kaia consensus nodes (CNs). The CnSV3 can change its reward address via `AddressBook.reviseRewardAddress()` function.

### StakingTracker

CnSV3 notifies its balance change and voter address change to the `StakingTracker` contract.

Whenever its balance or unstaking amount changes, `StakingTracker` must re-evaluate the voting powers. Therefore CnSV3 calls `StakingTracker.refreshStake` in staking-related functions.
When a GC member wishes to change its voter address, it must update the `StakingTracker` as well. Therefore CnSV3 calls `StakingTracker.refreshVoter` in multisig `UpdateVoterAddress()`.
