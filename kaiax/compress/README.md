# kaiax/compression

This module is responsible for reducing the disk storage size.


## Concepts
The module compresses preodically header, body, and receipts storage types.

The preodic parameters are determined by (1) number of blocks (chunk) and (2) number of byte size (cap).
If the chunk size resides in cap size, a chunk is compressed together.
If not, number of blocks to be compressed would be less than predefeind chunk value.

Once the module compressed storage data by chunk level, origin storage data is remomoved.
Data retreival process is expanded by two phase.
1. Read the data (i.e., `ReadHeader`, `ReadBody`, `ReadRecceipts`)
2. If `1` is failed, find a chunk from compressed storage and decompress it.

Eventually, origin storage disk will be get rid off progressively and compressed storage occupies most


This module is designed to periodically compress storage types, including headers, bodies, and receipts. The compression process is based on configurable periodic parameters.

### Compression Parameters
The module uses two primary parameters to determine compression:

1. Number of Blocks (Chunk)
- Defines the number of blocks grouped together for compression.

2. Upper Chunk Size (Cap)
- Specifies the maximum byte size of a chunk.

If the chunk size fits within the cap, the entire chunk is compressed.
If the chunk exceeds the cap size, only a portion of the chunk (fewer blocks) is compressed, adhering to the size limit.
Once compression is complete, the original uncompressed storage data is removed, reducing storage usage.

### Data Retrieval
The data retrieval process involves two phases:

1. Direct Read
- Attempt to read data directly using one of the following functions (`ReadHeader`, `ReadBody`, `ReadReceipts`)
2. Decompression (Fallback)
- If the direct read fails, the module locates the corresponding compressed chunk in storage and decompresses it to retrieve the requested data.

## Persistent schema
1. `Compressed` storage
Three types of keys are newly introduced:
- `CompressedHeader-` + block range (`to-from`)
- `CompressedBody-` + block range (`to-from`)
- `CompressedReceipts-` + block range (`to-from`)

These key-value is stored into directories `compressed_header`, `compressed_body`, and `compressed_receipts`, respectively.

2. `Misc` storage
Three types key has been added to `Misc` database.
- `CompressedHeader-` + next compression block number
- `CompressedBody-` + next compression block number
- `CompressedReceipts-` + next compression block number

These key-value is stored into `misc` directory and represents where compression should be started continuously.

### Start and stop

This module operates one background thread to iterate over the blocks and compress three kind of storage types: header, bdoy, and receipts.

### Rewind

Upon rewind, this module deletes the related persistent compressed data and flushes the in-memory cache.

## APIs
No additional APIs has been added within this package.

## Getters
No export getters
