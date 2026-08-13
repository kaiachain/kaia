# kaiax/system

This module is responsible for system state transitions that are
triggered at specific hardfork boundaries.

## Concepts

This module groups system-wide state mutations that do not naturally belong to
feature-specific modules.
Examples include:

- Installing the Randao registry at the fork block
- Executing treasury rebalance logic at KIP-103 and KIP-160 fork blocks
- Replacing or restoring the Mainnet credit contract code at Kaia and Osaka fork boundaries

## Persistent schema

This module does not persist any data.

## In-memory structures

This module does not maintain in-memory caches.

## Module lifecycle

### Init

- Dependencies:
  - Chain: Provides chain config and block context for fork checks.

### Start and stop

This module does not have any background threads.

## Block processing

### BlockStateModule

#### InitializeState

This module does not modify state before transaction execution begins.

#### FinalizeState

After all transactions are processed, this module applies hardfork-triggered
system state transitions in-place. This includes randao registry installation,
treasury rebalance execution, and Mainnet credit contract migration steps when
their corresponding fork conditions are met.

## APIs

This module does not expose APIs.
