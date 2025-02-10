# kaiax/randao

This module is responsible for providing the BLS public key at a given block number.

## Concepts

The BLS key serves as a cryptographic source of randomness for the Kaia network since the Randao hardfork. It has two primary use cases:

1. **Proposer Selection**: Used to deterministically but unpredictably select the block proposer for each block. [KIP-146](https://kips.kaia.io/KIPs/kip-146)

2. **PREVRANDAO**: Provides randomness to EVM through the RANDOM opcode. [KIP-114](https://kips.kaia.io/KIPs/kip-114)

The BLS public key is used to verify the BLS signature signed by the validators.

## Persistent schema

This module does not persist any data.

## Module lifecycle

### Init

- Dependencies:
  - ChainConfig: Holds the RandaomCompatibleBlock value at genesis.
  - Chain: Provides the blocks and states.

### Start and stop

This module does not have any background threads.

## Block processing

### Execution

This module caches the BLS public key for the next block.

During synchronization, when blocks are processed rapidly in succession, the module skips the caching of future block BLS keys to avoid memory bloat. (And not effective as there's no enough interval between `PostInsertBlock`.) The cache is only filled once the node catches up to the network.

### Rewind

Upon rewind, this module purges the in-memory cache.

## APIs

This module does not expose any APIs.
