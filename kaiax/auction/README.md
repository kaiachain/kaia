# kaiax/auction

This module is responsible for processing the auction specified by [KIP-249](https://kips.kaia.io/KIPs/kip-249).

## Concepts

The bid is a data that contains the information to generate a transaction to be executed right after the target transaction is executed.

Note that the bid comes from the `Auctioneer`, which is an independent service that is responsible for processing auction and submit winner's bid to the Kaia client (CN).

## Bid pool validation rules

A bid pool is responsible for managing the valid bids from the auctioneer. The bid must satisfy the following rules:

1. The `bid.BlockNumber` must be in range of `[currentBlockNumber + 1, currentBlockNumber + allowFutureBlock]`.
2. The `bid.Bid` must be greater than 0.
3. The gas limit must be less than the maximum limit.
4. The `bid.Sender` must not be in the winner list of the same block number if the new bid isn't equal to the previous bid.
5. The `bid.SearcherSig` and `bid.AuctioneerSig` must be valid.

## Block building rules

Upon detection of target transaction from the tx pool, the following logics are executed:

- The corresponding bid is retrieved from the bid pool.
- If the bid is found, a new bundle is generated which contain `[BidTx]`.
- If the target transaction is not found in the bid pool, the bid will be ignored.

## Persistent schema

This module does not persist any data.

## Module lifecycle

### Init

- Dependencies:
  - ChainConfig: to generate the latest signer.
  - NodeKey: for BidTxGenerator.
- Notable dependents:
  - worker: to extract bundles.

### Start and stop

It starts/stops the bid pool and corresponding background threads.

## Block processing

### Execution

This module reads `SystemRegistry` and `AuctionEntryPoint` to detect any changes in `AuctionEntryPoint` and `Auctioneer` address. If the one of them is changed, the module will clear the existing bids in the bid pool.

## APIs

### auction_submitBid

Send a bid and get bid hash if successful, otherwise empty hash with error.

```sh
curl -H "Content-Type: application/json" \
    --data '{
        "jsonrpc": "2.0",
        "method": "auction_submitBid",
        "params": [
            {
                "targetTxRaw": "0xf8674785066720b30083015f909496bd8e216c0d894c0486341288bf486d5686c5b601808207f4a0a97fa83b989a6d66acc942d1cbd70f548c21e24eefea12e72f8c27ba4369a434a01900811315ba3c64055e9778470f438128b54a46712cc032f25a1487e2144578",
                "targetTxHash": "0xc7f1b27b0c69006738b17567a7127c4d163fac7b575d046c6cbc90e62e6355e8",
                "blockNumber": 1,
                "sender": "0x14791697260E4c9A71f18484C9f997B308e59325",
                "to": "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707",
                "nonce": 4,
                "bid": 3,
                "callGasLimit": 2,
                "data": "0x1234",
                "searchersig": "0x2cd97ec3eb8230a8cac9169146ea6ca406d908edd488e5fda30811ebf56647d94740d582c592e3476481b3fbab38a100623d2f4b0615da8b8dfd0f99128879901b",
                "auctioneerSig": "0xd87718806c267dd6de19f4ed1111742750ee8040fdb3d18b1bd0dc1020ad8ca84262dfb4a3449f53b2cef8e2142796a96cca9ff8d08302f07db1d53a7b792e8d1c"
            }
        ],
        "id": 1
    }' http://localhost:8551
```

```
json
{
  "bidHash": 0x...
  "err": "..."
}
```

Go client can use `SendAuctionTx(Context, BidInput) (map[string]any, error)`, which is the same format with JSON RPC.
