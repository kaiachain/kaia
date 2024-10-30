# kaiax/supply

This module is responsible for tracking the total supply of the native token (KAIA).

## Concepts

The total supply of the token can be theoretically calculated by traversing every account in the world state trie. However, this is extremely inefficient to traverse just once. Instead, this module indirectly calculates it by summing up the increasing factor and decreasing factors at each block:

- The initial total supply at Genesis block is calculated from the genesis state trie.
- The total supply increases every block by the minting amount.
- The total supply decreases every block by the burnt fees.
- The total supply may increase or decrease after the TreasuryRebalance KIP-103 and KIP-160.
- The total supply may considered decrease if tokens are burnt by sending to the canonical burn addresses (0x0 and 0xdead). Note that in this case, the sum of balances in the state trie does not change but the total supply does.

Therefore, total supply at a given block N is calculated by:

- TotalMinted
  - (+) Genesis supply
  - (+) Minted tokens in blocks 1 through N
  - Must be accumulated for every block.
- BurntFee
  - (-) Burnt fees in blocks 1 through N
  - Must be accumulated for every block.
- CanonicalBurn
  - (-) Balances of the canonical burn addresses at block N
  - Simple GetBalance query.
- RebalanceBurn
  - (-) Rebalance net burnt amount up to block N
    - KIP-103 burn amount if N >= Kip103ForkNum
    - KIP-160 burn amount if N >= Kip160ForkNum
    - Read the TreasuryRebalance contract's `memo` variable. Note 


## Persistent schema

- `SupplyCheckpoint(num)`: Accumulated TotalMinted and BurntFee up to block `num`.
  ```
  "supplyCheckpoint" || Uint64BE(num) => RLP(TotalMinted.Bytes(), BurntFee.Bytes())
  ```
- `lastSupplyCheckpointNumber`: The largest block number where the SupplyCheckpoint is stored.
  ```
  "lastSupplyCheckpointNumber" => Uint64BE(num)
  ```

## In-memory structures

### SupplyCheckpoint


### TotalSupply


### TotalSupplyResponse

## Module lifecycle

### Init

- Dependencies:
  - ChainDB: Raw key-value database to access this module's persistent schema.
  - Chain: Provides the blocks and states.
  - kaiax/reward: Query the RewardSummary.

### Start and stop

This module operates one background thread to iterate over the blocks and accumulate the SupplyCheckpoint up to the latest block.

## Block processing

### Consensus

This module does not have any consensus-related block processing logic.

### Execution

After a new block is inserted, this module updates the SupplyCheckpoint.

### Rewind

Upon rewind, this module deletes the related persistent data and flushes the in-memory cache.

## APIs

### kaia_getTotalSupply

Query the total supply at the given block.

- Parameters
  - `num`: block number
  - `showPartial`: If specified and equals true, the API will return a best-effort partial information even if block states or TreasuryRebalance memos are missing. Otherwise, the API will fail when some information is unavailable.
- Returns
  - `TotalSupplyResponse`
- Example
```sh
curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
  {"jsonrpc":"2.0","id":1,"method":"kaia_getTotalSupply","params":[
    "latest"
  ]}' | jq .result
```
```json
{
  "number": "0xa0ded60",
  "totalSupply": "0x446c3b15f9926687d2d44701ad1de466362b643749",
  "totalMinted": "0x446c3b15f9926687d2c8b9b2855a06795ef6400000",
  "totalBurnt": "-0xb8d4f27c3ddecd735243749",
  "burntFee": "0xba97a444dd853f01bbfd",
  "zeroBurn": "0x65e15d8ed6b50c008ac",
  "deadBurn": "0x3a1a8c88eb4fded2fb",
  "kip103Burn": "0x6a15eef7e3a354cf657913",
  "kip160Burn": "-0xbf82646906be407e42a4800"
}
```

## Getters

- GetTotalSupply: Returns the TotalSupply at the given block.
  ```
  GetTotalSupply(num) -> TotalSupply
  ```
