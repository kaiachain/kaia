# kcn valops - Validator Operations CLI

Command-line interface for interacting with the AddressBookV2 contract (`0x400`) on-chain, plus an offline `generate-keys` helper for creating validator onboarding keys.

## Usage

```
kcn valops <command> [flags] [args]
```

## Flags (on-chain commands)

These apply to every command except the offline `generate-keys` and `sign-create-node` (see [Key generation](#key-generation-offline) and [Registration signature](#registration-signature-offline)).

| Flag | Default | Description |
|------|---------|-------------|
| `--datadir` | `/var/kcnd/data` | Node data directory. rpm installs use the default; for a tar install, pass `DATA_DIR` from `conf/kcnd.conf` |
| `--endpoint` | `<datadir>/klay.ipc` | RPC endpoint (HTTP/WS/IPC) |
| `--private-key` | | Hex-encoded private key, used instead of the node key file |
| `--timeout` | `2m` | How long to wait for the transaction receipt |

Unless `--private-key` is set, the key is loaded from the file the node uses: `<datadir>/nodekey` if it exists, else `<datadir>/klay/nodekey`. Prefer the node key file: a key given with `--private-key` stays in the shell history and is visible in process listings while the command runs.

Chain ID is fetched automatically from the endpoint. The gas limit is the estimate plus 50%, and the gas price is the endpoint's suggested price. The sender must hold at least gas limit × gas price; the command checks this before sending.

## Commands

### Suspender role

These commands require the caller to hold the suspender role in AddressBookV2.

```
kcn valops suspend-validator --node-id <address> --private-key <key> --endpoint <endpoint>
kcn valops unsuspend-validator --node-id <address> --private-key <key> --endpoint <endpoint>
```

| Flag | Description |
|------|-------------|
| `--node-id` | Address of the target validator node |

### Node operator role

These commands require `msg.sender == node-id` (the private key **is** the node key).

```
kcn valops ready-candidate --datadir <path>
kcn valops unready-candidate --datadir <path>
kcn valops ready-validator --datadir <path>
kcn valops unready-validator --datadir <path>
kcn valops pause --datadir <path>
kcn valops resume --datadir <path>
kcn valops exit --datadir <path>
kcn valops offboard --datadir <path>
```

The node-id is derived from the key.

`ready-candidate` and `ready-validator` print a warning, and still send the transaction, when:

- less than 5 minutes are left before the next epoch block (one block per second). If the transaction is included after that block, the transition takes effect one epoch later.
- (`ready-validator` only) the node's idle timeout comes before the epoch block at which the transition would take effect. The node then becomes Registered, not ValActive, and must go through `ready-candidate` again.

The node key address pays the gas, so it needs a small KAIA balance.

### Registration signature (offline)

`sign-create-node` signs the node-ID ownership proof that `AddressBookV2.createNode` requires. It is offline and takes the key the same way as the on-chain commands (`--datadir` or `--private-key`).

```
kcn valops sign-create-node --datadir <path> --chain-id <id> --manager <address> --staking-contract <address> [--node-id <address>]
```

| Flag | Description |
|------|-------------|
| `--chain-id` | Chain ID the node runs on (8217 mainnet, 1001 kairos) |
| `--manager` | Address that will send `createNode` (becomes the node's manager) |
| `--staking-contract` | CnStaking contract address bound into the proof |
| `--node-id` | Optional. Expected node ID; the command fails if it does not match the signing key |

The signature is the only output on stdout; pass it to `createNode` as the node-ID signature. The inputs it signed are printed on stderr. `createNode` must be sent from `--manager` with the same staking contract, or it reverts.

### Key generation (offline)

`generate-keys` creates the full validator onboarding key set. It is offline (no chain interaction) and uses `--datadir` instead of `--endpoint`/`--private-key`.

```
kcn valops generate-keys --datadir <path>
```

| Flag | Description |
|------|-------------|
| `--datadir` | Node data directory (default `/var/kcnd/data`, the rpm install path; for a tar install, pass `DATA_DIR` from `conf/kcnd.conf`); keys are written under `<datadir>/klay` |

Existing files are never overwritten: every file is created only if it does not exist at the moment of writing. The command reuses existing key material and creates only the missing files:

- An existing raw `nodekey` or `bls-nodekey` is loaded and used (`reused`). Each is looked up as the node does: `<datadir>/<file>` if it exists, else `<datadir>/klay/<file>`. This covers a node that was started before `generate-keys`. If only `nodekey` exists, the BLS key is derived from it the same way the node derives it.
- An existing keystore is decrypted with its `.pass` file. It is kept (`skipped`) only if it decrypts, its `address` (or EIP-2335 `pubkey`) field belongs to the decrypted key, and, for `nodekey.json` and `bls-nodekey.json`, the key equals the raw key. Operator keystores print the address of the decrypted key.
- An existing `bls-pub` is kept if it equals the BLS public key, and an existing `bls-pop` if it is a valid proof-of-possession for it.
- Missing files are generated (`created`).
- When a new `nodekey` is created, the command prints a warning, because it means a new node ID. On a node that has already run, check `--datadir`.

All checks run before any write. If a file does not match its raw key, or a `.json` exists without its `.pass` (or the reverse), the command writes nothing and names the file. It writes:

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
# rpm install: nodekey and IPC endpoint under /var/kcnd/data
kcn valops ready-candidate

# tar install: pass DATA_DIR from conf/kcnd.conf; the IPC endpoint is <datadir>/klay.ipc
kcn valops ready-candidate --datadir /data/kcnd

# Specify a custom endpoint
kcn valops pause --datadir /var/kcnd/data --endpoint http://localhost:8551

# Suspend a validator from an rpm node whose node key holds the suspender role
kcn valops suspend-validator --node-id 0xABCD...

# Suspend a validator (as suspender); read the key from a file instead of typing it
kcn valops suspend-validator --node-id 0xABCD... --private-key "$(cat /path/to/suspender.key)" --endpoint http://localhost:8551

# Sign the createNode proof (offline)
kcn valops sign-create-node --datadir /var/kcnd/data --chain-id 8217 --manager 0x1111... --staking-contract 0x2222...

# Generate the validator key set, reusing keys of an already started node (offline)
kcn valops generate-keys --datadir /var/kcnd/data
```

## Output

Each on-chain command first prints what it is about to send, then the transaction hash, then waits for the receipt (`generate-keys` instead prints the key summary with a created/reused/skipped status per key):

```
method:   readyCandidate
datadir:  /var/kcnd/data (default)
endpoint: /var/kcnd/data/klay.ipc
chain ID: 8217
sender:   0x9a3f...
node ID:  0x9a3f...
epoch:    next epoch block 233366400, about 3h12m0s left
tx: 0xabc123...
waiting for receipt...
status: success
state:    CandReady
suspended: false
```

The `epoch` line and any `WARNING:` lines appear only for `ready-candidate` and `ready-validator`; `ready-validator` also prints a `timeout:` line when the node has an idle timeout. `state` and `suspended` are read from AddressBookV2 after the receipt.

On failure:

```
pre-flight check reverted: <revert reason> (no transaction sent)
sender 0x9a3f... needs <amount> KAIA for gas but holds <balance> KAIA; send KAIA to this address (no transaction sent)
transaction failed: <revert reason>
no receipt for 0xabc123... after 2m0s; the outcome is unknown. Check the transaction hash before sending again
```
