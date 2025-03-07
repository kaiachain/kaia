# kaiax/compress

This module is responsible for compressing parts of the blockchain database to reduce the storage footprint.

## Concepts

This module compresses header, (transaction) body, and receipts databases. These data has repetitive fields (e.g. extraData, calldata, event logs) that are worth compressing. But compressing each item individually is not efficient, so this module compresses the data in chunks. Find more details in [KIP-237](https://github.com/kaiachain/kips/blob/main/KIPs/kip-237.md).

This module will
1. Scan from the genesis to current block, collect items into chunks and compress them, delete the uncompressed data.
2. When new blocks are inserted, collect the new items into chunks and compress them, delete the uncompressed data.
3. When blockchain is rewound, decompress the chunks into original (uncompressed) databases.
4. Assist the DBManager to read from the compressed databases.

### Chunk size and retention

By default, a chunk is finalized (i.e. compressed) when either of the following conditions is met:
- `ChunkItemCap`: The chunk contains 10,000 items.
- `ChunkByteCap`: The chunk's uncompressed size exceeds 1MB.

A compression chunk size was determined by the following considerations:
- The chunk size should be small enough to be compressed and decompressed efficiently.
- The chunk size should be large enough to make the compression ratio high.
- The compression thread must not indifinetly wait for the chunk to be large enough.

This module will preserve the uncompressed data for at least `Retention` blocks for fast retrieval without decompression overhead. The default retention is 172,800 blocks (2 days).

### Compressed schema

In the [DBManager](../../storage/database/db_manager.go), each of these databases are stored in separate key-value databases at dedicated directories.

For efficiency and observability, compressed data will be stored other separate directories under `DATA_DIR/klay/chaindata/`

| Data type | Original (uncompressed) dir | Compressed dir |
|-|-|-|
| Headers   | `header`   | `header_compressed`  |
| Bodies    | `body`     | `body_compressed`    |
| Receipts  | `receipts` | `receipts_compressed`|

### Data retrieval

To read header, body, receipt, the caller shall:

- First attempt to read from the original (uncompressed) database.
- If not found, find the chunk that contains the item. Decompress the chunk and extract the item. The decompressed chunk may be cached.

## Persistent schema

This module uses schemas in multiple databases.

- In `header_compressed` database,
  - `CompressedHeader(from, to)`: Compressed header chunk for block numbers [from, to]. Note that `to` comes before `from` to enable the one-shot retrieval.
    ```
    "Compressed-h" || Uint64BE(to) || Uint64BE(from) => RLP.Encode([]{num, hash, headerRLP})
    ```
  - `NextCompressingHeaderNum`: The number after the last compressed chunk, i.e. the next block number to be compressed.
    ```
    "NextCompressingNum-h" => Uint64BE(num)
    ```
- In `body_compressed` database,
  - `CompressedBody(from, to)`: Compressed body chunk for block numbers [from, to].
    ```
    "Compressed-b" || Uint64BE(to) || Uint64BE(from) => RLP.Encode([]{num, hash, bodyRLP})
    ```
  - `NextCompressingBodyNum`: The number after the last compressed chunk, i.e. the next block number to be compressed.
    ```
    "NextCompressingNum-b" => Uint64BE(num)
    ```
- In `receipts_compressed` database,
  - `CompressedReceipt(from, to)`: Compressed receipt chunk for block numbers [from, to].
    ```
    "Compressed-r" || Uint64BE(to) || Uint64BE(from) => RLP.Encode([]{num, hash, receiptsRLP})
    ```
  - `NextCompressingReceiptsNum`: The number after the last compressed chunk, i.e. the next block number to be compressed.
    ```
    "NextCompressingNum-r" => Uint64BE(num)
    ```

## In-memory Structures

### ItemSchema

ItemSchema is an abstract interface that implements the type-specific operations for database interaction and compression. Data items such as headerRLP are handled as opaque bytes.

```go
type ItemSchema interface {
	name() string                                        // name that appears in logs
	uncompressedDb() database.Database                   // uncompressed database handle (e.g. HeaderDB)
	compressedDb() database.Database                     // compressed database handle (e.g. CompressedHeaderDB)
	uncompressedKey(num uint64, hash common.Hash) []byte // key for uncompressed item
	compressedKey(from, to uint64) []byte                // key for compressed chunk
	nextNumKey() []byte                                  // key for next number to be compressed
}
```

### Codec

Codec is an abstract interface that implements the compression and decompression logic. Currently this module supports zstd (https://github.com/klauspost/compress/tree/master/zstd#zstd).

```go
type Codec interface {
	compress(src []byte) ([]byte, error)
	decompress(src []byte) ([]byte, error)
}
```

## Module lifecycle

### Init

- Dependencies:
  - Chain: Provides the current block number.
  - DBManager: Provides handles to the databases.

### Start and stop

This module runs a background threads that compress data from the genesis to current block number minus retention.

## Block processing

### Rewind

Upon rewind, this module decompresses the chunks and restores the original databases.

## APIs

This module does not expose any APIs.

## Getters

- FindHeaderRLP: Find the header RLP from the compressed database.
  ```
  FindHeaderRLP(num, hash) -> (headerRLP, ok)
  ```
- FindBodyRLP: Find the body RLP from the compressed database.
  ```
  FindBodyRLP(num, hash) -> (bodyRLP, ok)
  ```
- FindReceiptsRLP: Find the receipts RLP from the compressed database.
  ```
  FindReceiptsRLP(num, hash) -> (receiptsRLP, ok)
  ```
