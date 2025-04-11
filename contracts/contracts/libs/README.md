# Libraries

These contracts are external dependencies. Some old libraries are kept to support other legacy contracts.

> [!NOTE]
> Uniswap's gobind edits compiled binaries because the contract expects a specific binary.

- `kip13/InterfaceIdentifier.sol`: The ERC-165 & KIP-13 supportsInterface.
- `openzeppelin-contracts-v2`: https://github.com/OpenZeppelin/openzeppelin-contracts/releases/tag/v2.3.0
- `uniswap/factory`: gobind of https://github.com/Uniswap/v2-core/blob/v1.0.1/contracts/UniswapV2Factory.sol
- `uniswap/factory`: gobind of https://github.com/Uniswap/v2-periphery/blob/master/contracts/UniswapV2Router02.sol
- `Ownable.sol`: Ownable from https://github.com/klaytn/klaytn-contracts/blob/main/contracts/access/Ownable.sol.
- `SafeMath.sol`: SafeMath for older solidity versions.
