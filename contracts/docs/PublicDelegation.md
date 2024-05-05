# PublicDelegation

The public delegation (PD) is a non-transferable ERC-4626 based contract that allows general users to delegate and re-delegate their KAIA to a Kaia Governance Council (GC) who uses CnStakingV3MultiSig (CnSV3) and enable PD.

It mints the tokenized shares to the delegator, which is called `pdKAIA`. The `pdKAIA` is a non-transferable interest-bearing token that represents the delegator's share of the total KAIA delegated to the GC. As rewards are compunded, the exchange rate of `pdKAIA` to KAIA increases. The delegator can burn the `pdKAIA` to get the KAIA back. All the math comes from the ERC-4626 standard.

Unlike usual ERC-4626 vault contracts, the reward is directly distributed to PD contract by state modification at the consensus-level. The reward will be automatically compounded to the CnSV3 contract.

It is deployed during setup process of CnSV3 contract, namely `setPublicDelegation` function. The PD is only compatible with the CnSV3 contract.

## Code

- [IPublicDelegation.sol](../contracts/consensus/PublicDelegation/IPublicDelegation.sol)
- [PublicDelegationStorage.sol](../contracts/consensus/PublicDelegation/PublicDelegationStorage.sol)
- [PublicDelegation.sol](../contracts/consensus/PublicDelegation/PublicDelegation.sol)

## Access Control

The access control is managed by the `openzeppelin/Ownable.sol` library. The ownable functions in PD are as follows:

- `updateCommissionTo(addr)`: Update the commission receiver address.
- `updateCommissionRate(commissionRate)`: Update the commission rate. `MAX_COMMISSION_RATE` is 3,000, which is 30%.

## Features

The PD's features can be categorized into 4 groups:

- Commission management
- Delegation
- Withdrawal
- Re-delegation

### Commission Management

The PD contract can collect the commission from the rewards. The commission information can be updated by the owner.

Related functions:

- `updateCommissionTo(addr)`: Update the commission receiver address.
- `updateCommissionRate(commissionRate)`: Update the commission rate. `MAX_COMMISSION_RATE` is 3,000, which is 30%.

The commission is calculated and sent to the commission receiver whenever the rewards are compounded. The commission is calculated as follows:

$commission = ⌊rewards * commissionRate / MAX\_COMMISSION\_RATE⌋$

### Delegation

All general users can delegate their KAIA to the CnSV3 through the PD contract.

Related functions:

- `stake()`: Delegate the staked KAIA to the CnSV3 and mint the `pdKAIA` to the delegator.
- `stakeFor(_receipient)`: `stake` on behalf of the `_receipient`.

The delegators will receive the corresponding `pdKAIA` as shares. The shares are minted by the following formula when delegating `stakedKAIA`:

$pdKAIA = ⌊stakedKAIA * totalShares / totalStakedKAIA⌋$

Note that `totalStakedKAIA⌋` includes the current rewards, which are not yet compounded.

### Withdrawal

The withdrawal needs two steps: `withdraw` and `claim`. The withdrawal request will initiate the lockup period, and the delegator can claim the withdrawal request after the lockup period.

Related functions:

- `withdraw(recipient, KAIA)`: Burn the `previewWithdraw(KAIA)` and request withdrawal of the `KAIA` to the CnSV3.
- `redeem(recipient, pdKAIA)`: Burn the `pdKAIA` and request withdrawal of the `previewRedeem(KAIA)` to the CnSV3.

- `claim(requestId)`: Claim the withdrawal request after lockup period.

The expected burn amount of `pdKAIA` when withdraw `KAIA` is calculated as follows:

$pdKAIA = ⌈KAIA * totalShares / totalStakedKAIA⌉$

The `withdraw` function is used to withdraw the exact amount of `KAIA`, while the `redeem` function is used to withdraw the `KAIA` equivalent to the `pdKAIA`. The users must use the `redeem` to withdraw the all `KAIA` since the rewards are accumulated every block (every 1 second), which leaves the dust amount of `KAIA` in the contract if use `withdraw`.

### Re-delegation

Users can redelegate their KAIA to another CnSV3 without waiting for lockup period. The detailed process can be found in [Re-delegation](./CnStakingV3.md#re-delegation).

- `redelegateByAssets(targetCnV3, KAIA)`: Redelegate the `KAIA` to the `targetCnV3`.
- `redelegateByShares(targetCnV3, pdKAIA)`: Redelegate the `previewRedeem(pdKAIA)` to the `targetCnV3`.

As same as the withdrawal, the `redelegateByAssets` is used to redelegate the exact amount of `KAIA`, while the `redelegateByShares` is used to redelegate the `KAIA` equivalent to the `pdKAIA`. The users must use the `redelegateByShares` to redelegate the all `KAIA` as well.
