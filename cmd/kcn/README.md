# kcn abv2 — AddressBookV2 CLI

Command-line interface for interacting with the AddressBookV2 contract (`0x400`) on-chain.

## Usage

```
kcn abv2 <command> [flags] [args]
```

## Flags (all commands)

| Flag | Default | Description |
|------|---------|-------------|
| `--endpoint` | `/var/kcnd/data/klay.ipc` | RPC endpoint (HTTP/WS/IPC) |
| `--private-key` | *(load from nodekey file)* | Hex-encoded private key |

When `--private-key` is omitted, the key is loaded from `/var/kcnd/data/klay/nodekey`.

Chain ID is fetched automatically from the endpoint.

## Commands

### Suspender role

These commands require the caller to hold the suspender role in AddressBookV2.

```
kcn abv2 suspend-validator --node-id <address>
kcn abv2 unsuspend-validator --node-id <address>
```

| Flag | Description |
|------|-------------|
| `--node-id` | Address of the target validator node |

### Node operator role

These commands require `msg.sender == node-id` (the private key **is** the node key).

```
kcn abv2 ready-candidate
kcn abv2 unready-candidate
kcn abv2 ready-validator
kcn abv2 unready-validator
kcn abv2 pause
kcn abv2 resume
kcn abv2 exit
kcn abv2 offboard
```

No extra arguments. The node-id is derived from the private key.

## Examples

```sh
# Suspend a validator (as suspender)
kcn abv2 suspend-validator --node-id 0xABCD... --private-key 0x1234...

# Use default IPC endpoint and nodekey file
kcn abv2 ready-candidate

# Specify a custom endpoint
kcn abv2 pause --endpoint http://localhost:8551

# Use a custom private key and endpoint
kcn abv2 ready-validator --private-key 0x1234... --endpoint /tmp/klay.ipc
```

## Output

Each command prints the transaction hash immediately upon submission, then waits for the receipt:

```
tx: 0xabc123...
status: success
```

On failure:

```
tx: 0xabc123...
transaction failed: <revert reason>
```
