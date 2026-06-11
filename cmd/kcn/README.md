# kcn valops - Validator Operations CLI

Command-line interface for interacting with the AddressBookV2 contract (`0x400`) on-chain.

## Usage

```
kcn valops <command> [flags] [args]
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
kcn valops suspend-validator --node-id <address>
kcn valops unsuspend-validator --node-id <address>
```

| Flag | Description |
|------|-------------|
| `--node-id` | Address of the target validator node |

### Node operator role

These commands require `msg.sender == node-id` (the private key **is** the node key).

```
kcn valops ready-candidate
kcn valops unready-candidate
kcn valops ready-validator
kcn valops unready-validator
kcn valops pause
kcn valops resume
kcn valops exit
kcn valops offboard
```

No extra arguments. The node-id is derived from the private key.

## Examples

```sh
# Suspend a validator (as suspender)
kcn valops suspend-validator --node-id 0xABCD... --private-key 0x1234...

# Use default IPC endpoint and nodekey file
kcn valops ready-candidate

# Specify a custom endpoint
kcn valops pause --endpoint http://localhost:8551

# Use a custom private key and endpoint
kcn valops ready-validator --private-key 0x1234... --endpoint /tmp/klay.ipc
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
