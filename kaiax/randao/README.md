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

#### PostInsertBlock

This module caches the BLS public key for the next block.

During synchronization, when blocks are processed rapidly in succession, the module skips the caching of future block BLS keys to avoid memory bloat. (And not effective as there's no enough interval between `PostInsertBlock`.) The cache is only filled once the node catches up to the network.

### Rewind

Upon rewind, this module purges the in-memory cache.

## APIs

### kaia_getBlsInfos

Returns the BLS public key and PoP of a given block.

```sh
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"kaia_getBlsInfos","params":[
    "latest"
  ]}' | jq .result
```

```json
{
  "0x136f6AB7fC073EAa7f943bf8863983BF9Ea4d899": {
    "pop": "b748eab431cd14bb6f29fe558fd86f8c719150f5cb11c8cf8c94b99243edc85554d7ff9e0b007eb56c4295ddb09929600189d53241a8cabd1a46f30f36147afa543980ba41097a5a347b84c7f7398e6faead0c770254aa1aad4ae867bce6f2c8",
    "publicKey": "9487860295d9d3f8137a57d85a41d4f31a89cf00622d425fe96e26f3c2848a3133158e2064d1793ce2bf9db443cec526",
    "verifyErr": null
  },
  "0x2C766CB7B2B8F21C0C82F23A5284E8bdC9b988e9": {
    "pop": "94e793e5089c08d6b2735c3236eee789386ebf4f640872772709e74ac4185644f7f405867a891d2d153ad26326ced2d00f459bdd834641b05aad4dfa14679f1b9305f243c53ec4315d31d98a20b471910bfcd19b7e490e3013839aaf14e08a9a",
    "publicKey": "aefafce19d30d92e1dd5156f88ccd6f99d29f4698af9645bdbed27d55c306f1f73be317fa3ac61bd3485765ff388eff2",
    "verifyErr": null
  },
  ...
}
```
