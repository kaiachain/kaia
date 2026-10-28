import * as path from "path";
import { HardhatUserConfig, subtask } from "hardhat/config";
import { TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS } from "hardhat/builtin-tasks/task-names";
import "@nomicfoundation/hardhat-toolbox";

import * as glob from "glob";

subtask(TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS, async () =>
  glob.sync(path.join(__dirname, "{libs,service_chain,testing}/**/*.sol"))
);

const config: HardhatUserConfig = {
  solidity: {
    compilers: [
      {
        // libs/openzeppelin-contracts-v2 (^0.4.24), testing/compiler/UnsafeMultiply_0.4.24.sol
        version: "0.4.24",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
      {
        // service_chain/bridge (0.5.6), testing/sc_erc20, testing/sc_erc721 (^0.5.0, ^0.5.6)
        version: "0.5.6",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
      {
        // testing/system_contracts/WKAIA.sol (>=0.5.9 <0.6.0); not satisfied by 0.5.6
        version: "0.5.9",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
      {
        // testing/system_contracts (^0.8.x)
        version: "0.8.19",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
      {
        // testing/system_contracts/MockCnStakingOverV2.sol (0.8.25)
        version: "0.8.25",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
    ],
  },
};

export default config;
