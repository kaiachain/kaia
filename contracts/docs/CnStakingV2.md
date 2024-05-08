# CnStakingV2

A member of Klaytn Governance Council (GC) must own at least one CnStaking contract. There are two versions of CnStaking contract so far.
- legacy/CnStakingContract (CnSV1): The original CnStaking contract that is being used since the start of Klaytn mainnet.
- CnStakingV2 (CnSV2): An upgraded version of CnSV1, having all interfaces of V1 as well as new governance-related features.

This document will describe the V2 contract.

## Code

See [ICnStakingV2.sol](../contracts/ICnStakingV2.sol) and [CnStakingV2.sol](../contracts/CnStakingV2.sol)

## Access control

CnStakingV2 has built-in multisig facility that requires certain number of approvals from the contract admin accounts, stored in following state variables. 
- `contractValidator` or CV:  The CV is a temporary admin who is involved in contract initialization. In the current GC onboarding process, the Klaytn team takes the role.
- `adminList`: the list of admins except `contractValidator`. Each GC member usually has 1, 2 or 3 admins.
- `requirement`: the required number of approvals.
- `isAdmin`: true for admins.

## Contract intialization

Initializing a CnStakingV2 takes multiple steps involving many entities.
1. `constructor()`, usually by the `contractValidator`, or CV.
2. `setStakingTracker()` by the CV or any admin, to specify the StakingTracker address.
3. `reviewInitialConditions()` by every admins and the CV, to signify that they all agree with initial lockup conditions.
4. `depositLockupStakingAndInit()` to deposit the initial lockup KLAYs and finish the initialization process.

## Multisig facility

Each multisig function call is stored in a `Request` struct. Upto three arguments are encoded as `bytes32` type each.

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
- NotConfirmed: before execution.
- Executed: successfully executed.
- ExecutionFailed: did call the function but failed.
- Canceled: canceled before execution.

A multisig function "X" is executed through following steps.
- An admin calls `submitX()` function, creating a `Request`. Any preconditions (e.g. cannot add already admin) are checked at this time.
- The `Request` content is available from either `SubmitRequest` event or `getRequestIds()`, `getRequestInfo()` functions.
- For zero or more times, other admins call `confirmRequest()`
- The function `x()` is executed when the `requirement` number of admins have called either `submitX()` or `confirmRequest()`. Any preconditions are checked again at this time because states may have changed since request creation.
- The function call to `x()` never reverts as `.call()` is used. However, depending on the call result either `ExecuteRequestSuccess` or `ExecuteRequestFailure` event is emitted.

For simplicity, calling "multisig X()" refers to all above steps.

## Features

CnStakingV2 is a multi-purpose contract. Its features can be categorized into 4 groups.
- Admin management
- Initial lockup stakes
- Free stakes
- External accounts management

### Admin management

The contract admins themselves are managed by multisig functions.

Related functions:
- `multisig AddAdmin(_admin)`: Add an account to admins.
- `multisig DeleteAdmin(_admin)`: Delete an account from admins.
- `multisig UpdateRequirement(_requirement)`: Change the requirement number.
- `multisig ClearRequest()`: Cancels all outstanding requests.

As a precaution, all outstanding (i.e. NotConfirmed) requests are canceled when either `adminList` or `requirement` changes.

### Initial lockup stakes

Initial lockup is a series of long-term fixed lockups. The lockup condition is stored in a singleton `LockupConditions` struct. Total amounts are tracked by `initialLockupStaking` and `remainingLockupStaking`.

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

### Non-lockup stakes

Admins can choose to deposit more stakes at their discretion. This feature is called "non-lockup" staking or "free" staking, though there is a 7 days of lockup before withdrawal.

The non-lockup staking amounts are tracked by `staking`, `unstaking`. Each withdrawal request is stored in `WithdrawalRequest` structs.

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
- Unknown: before transfer or cancellation. Unlike `RequestState.Unknown`, `WithdrawalStakingState.Unknwon` is a meaningful enum.
- Transferred: successfully executed.
- Canceled: canceled due to timeout or an explicit cancellation.

A withdrawal request becomes withdrawable from 7 days since creation (D+7). However it expires after another 7 days (D+14). 

At any moment, sum of unstaking (pending withdrawal) amount cannot exceed the staking amount. This rule prevents an abusing scenario where a GC member "spray" withdrawal requests across different days and effectively bypass the 7-day delay.

Related functions:
- `stakeKlay()`: Deposit non-lockup stakes.
- `receive()`: Equivalent to `stakeKlay()`.
- `multisig ApproveStakingWithdrawal(to, value)`: Create a WithdrawalRequest to take out all or part of non-lockup stakes.
- `multisig CancelApprovedStakingWithdrawal(id)`: Explicitly cancel a WithdrawalRequest.
- `withdrawApprovedStaking(id)`: Execute a WithdrawalRequest after 7 days of delay.

### External accounts

A Klaytn GC member uses several accounts for different purposes. Some of them are appointed via CnStakingV2 contract.

Related functions:
- `multisig UpdateRewardAddress(addr)`: Update the `rewardAddress` and also update to the AddressBook contract.
- `multisig UpdateStakingTracker(addr)`: Update the `stakingTracker` address. CnStakingV2 contract notifies staking balance change and voter address change to the StakingTracker.
- `multisig UpdateVoterAddress(addr)`: Update the `voterAddress` and also update to the StakingTracker contract.

## Interaction with other contracts

### AddressBook

The AddressBook contract (0x000..400) is the central directory of Klaytn consensus nodes (CNs). The CnStakingV2 can change its reward address via `AddressBook.reviseRewardAddress()` function.

### StakingTracker

CnStakingV2 notifies its balance change and voter address change to the StakingTracker contract.
- Whenever its balance or unstaking amount changes, StakingTracker must re-evaluate the voting powers. Therefore CnStakingV2 calls `StakingTracker.refreshStake` in staking-related functions.
- When a GC member wishes to change its voter address, it must update the StakingTracker as well. Therefore CnStakingV2 calls `StakingTracker.refreshVoter` in `multisig UpdateVoterAddress()` .


