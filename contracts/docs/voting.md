# Voting contract

## Klaytn governance structure

Klaytn as it is today is run by a permissioned validators called the Governance Council (GC). Each GC member run one Consensus Node (CN) to produce and validate blocks. The GC members also have authority to change Klaytn's network parameters such as gas price and block reward. Our governance contracts are trying to facilitate a transparent and stake-based voting by the GC members.

### Voting power

- Each GC member can have an integer amount of votes.
- A GC member obtains 1 vote for every 5,000,000 KLAYs staked.
- But one GC member cannot have more than [numGCs - 1] votes.
  ```
  votes = min( floor(stakeAmount / 5M), numGC - 1 )
  ```
- The vote cap of [numGC - 1] prevents single voter having more than 50% of the power. For instance, if there are 50 GC members with the most inequal situation, the top voter would have exactly 50%.
  ```
  49, 1, 1, 1, 1, 1, ..., 1, 1, 1, 1, 1
  |   <-- 49 members with min votes -->
  one member with max votes
  ```

### Quorum

A proposal passes when it satisfies a combination of the following conditions.

- `CountQuorum` = At least 1/3 of all eligible voters cast votes
- `PowerQuorum` = At least 1/3 of all eligible voting powers cast votes
- `MajorityYes` = Yes votes are more than half of total cast votes
- `PassCondition` = (CountQuorum or PowerQuorum) and MajorityYes

## Code

See [IVoting.sol](../contracts/IVoting.sol) and [Voting.sol](../contracts/Voting.sol)

## Difference from KIP-81

The [KIP-81 as of 2022-11-09](https://github.com/klaytn/kips/blob/a1d99a58a60d0e3743774aca00beedb49a3c89a8/KIPs/kip-81.md) has some outdated content. The differences are mainly access control related. Several business logic had to be added in the actual Voting contract.

## Access Control

Voting contract has two categories of users, secretary and voters.

- The secretary is a single account stored in the `secretary` variable. This account is intended to be controlled by the Klaytn Foundation. It will serve the GC by assisting administrative works such as submitting and executing proposals.
- Voters are identified by their `NodeID`. The list of voters differs per proposal, depending on the list of GC members registered in AddressBook and their staking amounts at the time of proposal submission.

Each function has different access control rule.
- In the long term, the Foundation is wants GC to run governance by themselves. This is to push Klaytn towards decentralization.
- There exists some feature that allows secretary to share or yield its rights to voters (i.e. GC). The Foundation, represented by the secretary, will at some point invoke those features to change Voting configuration. Once a change is made, the secretary should be impossible to revert by itself, though possible by voters. 

| Function | Allowed to | Reasoning |
|-|-|-|
| `propose` | secretary or voters or both. Depends on `accessRule` | for GC autonomy |
| `cancel` | proposer of the given proposal | |
| `castVote` | to voters | |
| `queue` and `execute` | secretary or voters or both. Depends on `accessRule` | for GC autonomy |
| `updateStakingTracker` | governance voting | StakingTracker is very important because it determines how voting powers are calculated. |
| `updateSecretary` | governance voting | Giving power outside GC must be careful. The secretary account must be a trusted 3rd party (e.g. Foundation). |
| `updateAccessRule` | governance voting or secretary | allowed to governance voting: for GC autonomy, allowed to secretary: for quick transition into GC autonomy |
| `updateTimingRule` | governance voting or secretary | allowed to governance voting: for GC autonomy, allowed to secretary: for emergency agenda |

## Flexible timing

In Klaytn there exists a on-chain governance system leveraging block headers ([docs](https://docs.klaytn.foundation/content/dapp/json-rpc/api-references/governance), [explainer in Korean](https://www.youtube.com/watch?v=UPyf7B0YvI0)). Major complaint about the existing system was that it takes 1 to 2 weeks for a parameter change to take effect ([example](https://medium.com/klaytn/klaytn-gas-price-reduction-schedule-2ba158e3630d))

Therefore the Voting contract will allow each proposal to have custom timing. The `propose()` function accepts two timing parameters `votingDelay` and `votingPeriod`, given that they fall within the `timingRule` ranges.

In the default setting, both `votingDelay` and `votingPeriod` are bound between 1 and 28 days. This makes the quickest proposal execution down to 4 days (1 day Pending, 1 day Active, 2 days Queued).

