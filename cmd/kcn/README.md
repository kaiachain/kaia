# kcn valops - Validator Operations CLI

Command-line interface for interacting with the AddressBookV2 contract (`0x400`) on-chain, plus an offline `generate-keys` helper for creating validator onboarding keys.

## Usage

```
kcn valops <command> [flags] [args]
```

## Flags (on-chain commands)

These apply to every command except the offline `generate-keys` (see [Key generation](#key-generation-offline)).

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

### Key generation (offline)

`generate-keys` creates the full validator onboarding key set. It is offline (no chain interaction) and uses `--datadir` instead of `--endpoint`/`--private-key`.

```
kcn valops generate-keys --datadir <path>
```

| Flag | Description |
|------|-------------|
| `--datadir` | Node data directory; keys are written under `<datadir>/klay` |

Existing files are never overwritten — the command aborts if any target already exists. It writes:

```
<datadir>/klay/
├── nodekey                          raw hex   (node loads this)
├── bls-nodekey                      raw hex   (node loads this)
└── operator-keys/
    ├── nodekey.json + .pass         v3 keystore
    ├── manager.json + .pass         v3 keystore
    ├── voter.json + .pass           v3 keystore
    ├── reward.json + .pass          v3 keystore
    ├── cnstaking-owner.json + .pass v3 keystore
    ├── mev-reward.json + .pass      v3 keystore
    ├── bls-nodekey.json + .pass     EIP-2335 keystore
    └── bls-pub / bls-pop            raw hex   (public; createNode blsInfo)
```

Key roles:

| Key | Role |
|-----|------|
| `nodekey` | node identity (p2p) + createNode `nodeId`; signs state transitions (onlyNodeId) |
| `bls-nodekey` | consensus BLS (randao/vrank) + createNode `blsInfo` (pub/pop) |
| `manager` | deploys CnStaking, stakes, and sends createNode (becomes NodeInfo.manager); holds ≥ 5M KAIA |
| `cnstaking-owner` | owner of the deployed CnStaking |
| `voter` | createNode `voterAddress` (on-chain governance) |
| `reward` | createNode `rewardAddress` (unused when PublicDelegation is enabled) |
| `mev-reward` | reward recipient for the auction (MEV) contract; not a createNode argument |

ECDSA keys are encrypted as Web3 Secret Storage **v3** keystores; the BLS key as an **EIP-2335** keystore. Each keystore's password is a random value written to the adjacent `.pass` file. `nodekey`/`bls-nodekey` are also kept as raw hex because the node loads them directly; `bls-pub`/`bls-pop` are public values for createNode.

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

# Generate a fresh validator key set (offline)
kcn valops generate-keys --datadir /var/kcnd/data
```

## Output

Each on-chain command prints the transaction hash immediately upon submission, then waits for the receipt (`generate-keys` instead prints the generated key summary):

```
tx: 0xabc123...
status: success
```

On failure:

```
tx: 0xabc123...
transaction failed: <revert reason>
```
