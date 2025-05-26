import { HardhatUserConfig, subtask } from "hardhat/config";
import { TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS } from "hardhat/builtin-tasks/task-names";
import "@nomicfoundation/hardhat-toolbox";
import "@openzeppelin/hardhat-upgrades";

const config: HardhatUserConfig = {
  solidity: {
    compilers: [
      {
        version: "0.4.24",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
      {
        version: "0.5.6",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
      {
        version: "0.5.9",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
      {
        version: "0.8.19",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
      {
        version: "0.8.24",
        settings: { optimizer: { enabled: true, runs: 1000 }, viaIR: true },
      },
      {
        version: "0.8.25",
        settings: { optimizer: { enabled: true, runs: 200 } },
      },
    ],
  },
};

// Remove system_contracts/misc from the source paths
subtask(TASK_COMPILE_SOLIDITY_GET_SOURCE_PATHS).setAction(async (_, __, runSuper) => {
  const paths = await runSuper();
  return paths.filter((p: string) => !p.includes("system_contracts/misc"));
});

export default config;
