# kaiax/supply

This module is responsible for tracking the total supply of the native token (KAIA).

## Concepts

The total supply of the token can be theoretically calculated by traversing every account in the world state trie. However, this is extremely inefficient to traverse just once. Instead, this module indirectly calculates it by summing up the increasing factor and decreasing factors.

- AccReward (accumulated reward)
  - (+) Genesis supply
  - (+) Minted tokens in blocks 1 through N
  - (-) Burnt fees in blocks 1 through N
  - Genesis supply is calculated from the genesis state trie.
  - Since then, calculated for every block according to reward rules.
- CanonicalBurn
  - (-) Balances of the canonical burn addresses at block N
  - The tokens in the canonical burn addresses (0x0 and 0xdead) are considered burnt. Note that those balances are excluded from the total supply but still appears in the state trie.
  - Simple GetBalance query.
- RebalanceBurn
  - (-) Rebalance net burnt amount up to block N
    - KIP-103 burn amount if N >= Kip103ForkNum
    - KIP-160 burn amount if N >= Kip160ForkNum
    - Read from the TreasuryRebalance contract's `memo` variable. Note that Kairos (chainid=1001) data is overridden with fixed constant as its KIP-160 memo was irreversibly malformed.
    - Because memo has to be manually set after the fork, the data may not be available during the short period after the fork.

Because AccReward has to be accumulated over the whole blocks, intermediate results are periodically committed to the database as supply checkpoints. If total supply at a block is needed, the module will read the nearest supply checkpoint and re-accumulate the rewards from there.

## Persistent schema

- `AccReward(num)`: Stores the AccReward up to block `num`. Periodically committed to the database.
  ```
  "supplyCheckpoint" || Uint64BE(num) => RLP([Minted.Bytes(), BurntFee.Bytes()])
  ```
- `LastAccRewardNumber()`: The highest block number where the AccReward is stored.
  ```
  "lastSupplyCheckpointNumber" => Uint64BE(num)
  ```
- Note that the schema keys have "supplyCheckpoint" for backward compatibility.

## In-memory structures

### AccReward

AccReward represents the accumulated minted tokens and burnt fees up to a specific block number.

```go
type AccReward struct {
	TotalMinted   *big.Int // Genesis + Minted[1..n]
	BurntFee *big.Int // BurntFee[1..n]
}
```

### TotalSupply

TotalSupply represents the native token's total supply and its breakdown at a specific block number.

```go
type TotalSupply struct {
	TotalSupply *big.Int // TotalMinted - TotalBurnt

	// Because there is only minting source (block reward), TotalMinted equals to AccReward.Minted.
	TotalMinted *big.Int // Sum of all minted amounts: Genesis + Minted[1..n]

	// Tokens are burnt by various mechanisms.
	TotalBurnt  *big.Int // Sum of all burnt amounts: BurntFee[1..n] + CanonicalBurn[n] + RebalanceBurn[n]
	BurntFee    *big.Int // BurntFee[1..n]
	ZeroBurn    *big.Int // CanonicalBurn[n] at 0x0
	DeadBurn    *big.Int // CanonicalBurn[n] at 0xdead
	Kip103Burn  *big.Int // RebalanceBurn[n] by KIP-103
	Kip160Burn  *big.Int // RebalanceBurn[n] by KIP-160
}
```

### TotalSupplyResponse

TotalSupplyResponse is the response type for the `kaia_getTotalSupply` API. In addition to TotalSupply fields, it includes the block number and error string (if showPartial=true and some information is missing).

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
  - `showPartial`: If specified and equals true, the API will return a best-effort partial information even if block states or TreasuryRebalance memos are missing. Otherwise, the API will fail when some information is unavailable. The showPartial=true mode is primarily for debugging purposes.
- Returns
  - `TotalSupplyResponse`
- Example1 (full information)
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
- Example2 (partial information, showPartial=false)
  ```sh
  curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
    {"jsonrpc":"2.0","id":1,"method":"kaia_getTotalSupply","params":[
      "0x1000"
    ]}' | jq
  ```
  ```json
  {
    "jsonrpc": "2.0",
    "id": 1,
    "error": {
      "code": -32000,
      "message": "cannot determine canonical (0x0, 0xdead) burn amount: missing trie node cebc1a5911a6bda6ca34d46240fea8da49673e20afcdce46c3dd91dd2b6f41cc (path )"
    }
  }
  ```
- Example3 (partial information, showPartial=true)
  ```sh
  curl "http://localhost:8551" -X POST -H 'Content-Type: application/json' --data '
    {"jsonrpc":"2.0","id":1,"method":"kaia_getTotalSupply","params":[
      "0x1000", true
    ]}' | jq
  ```
  ```json
  {
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
      "number": "0x1000",
      "error": "cannot determine canonical (0x0, 0xdead) burn amount: missing trie node cebc1a5911a6bda6ca34d46240fea8da49673e20afcdce46c3dd91dd2b6f41cc (path )",
      "totalSupply": null,
      "totalMinted": "0x446c3b15f9926687d2c4053d515636313c00000000",
      "totalBurnt": null,
      "burntFee": "0x0",
      "zeroBurn": null,
      "deadBurn": null,
      "kip103Burn": "0x0",
      "kip160Burn": "0x0"
    }
  }
  ```


## Getters

- GetTotalSupply: Returns the TotalSupply at the given block. Because this function returns best-effort partial information, its return value has three cases:
  - If all information is available, it returns `(ts, nil)`
  - If some information is missing, it returns `(ts, err)`
  - If essential information is missing or something is wrong, it returns `(nil, err)`
  ```
  GetTotalSupply(num) -> (TotalSupply, error)
  ```
