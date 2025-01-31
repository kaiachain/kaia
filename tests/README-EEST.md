# Kaia's EEST (Ethereum Compatibility Test)

This test is using EEST for the Ethereum compatibility. It is performed in the following:

- `TestExecutionSpecBlockTestSuite` in `tests/block_test.go`
- `TestExecutionSpecStateTestSuite` in `tests/state_test.go`

Some EIP tests skip because Kaia don't support them.

## How it works

Vanilla Kaia can't execute EEST because of the different process from Ethereum such as a rewarding system. Specifically, the following items are different:

- Rule of precompiled contract address
- Intrinsic gas
- Gas price
- Op code constant gas
- Reward amount
- State root

This test can execute EEST overwriting these into the same one as Ethereum.

Regarding rule of precompiled contract address, this rule is used everywhere so this test changes the rule before launching blockchain.

Regarding intrinsic gas, gas price, and op code, these values are used in `ApplyMessage` so Kaia overwrites a transaction and evm before executing `ApplyMessage` to use them.

Regarding reward amount and state root, these
values are used after executing all transactions so this test simply overwrite the values after retrieving executed state DB.
