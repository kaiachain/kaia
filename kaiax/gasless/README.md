# kaiax/gasless

This module is responsible for enabling gasless transactions.

## Concepts

### Transaction patterns

- LendTx
- GaslessApproveTx
- GaslessSwapTx

### Transaction pool rules

TBU

### Block building rules

TBU

## In-memory structures

## Module lifecycle

## APIs

## Getters

- IsGaslessApproveTx(tx *types.Transaction) bool
- IsGaslessSwapTx(tx *types.Transaction) bool
- IsGaslessPattern(approveTxOrNil, swapTx *types.Transaction) bool
- MakeLendTx(approveTxOrNil, swapTx *types.Transaction) (*types.Transaction, error)
