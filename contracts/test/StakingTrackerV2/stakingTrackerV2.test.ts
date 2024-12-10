import { loadFixture, setBalance } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import { ABOOK_ADDRESS, REGISTRY_ADDRESS, MILLIONS } from "../common/helper";
import { BigNumber } from "ethers";
import { registerContract, stakingTrackerV2TestFixture } from "../materials";
import { FakeContract, smock } from "@defi-wonderland/smock";
import {
  ICLPool,
  ICLPool__factory,
  ICLRegistry,
  ICLRegistry__factory,
  IERC20,
  IERC20__factory,
  StakingTrackerV2,
} from "../../typechain-types";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

/**
 * @dev This unit & scenario test is for StakingTrackerV2.sol.
 */
describe("StakingTrackerV2.sol", function () {
  async function verifyTrackerState(
    stakingTracker: StakingTrackerV2,
    trackerId: number,
    trackStart: number,
    trackEnd: number,
    gcIds: number[],
    cnStakingBalances: BigNumber[],
    gcBalances: BigNumber[],
    votes: number[]
  ) {
    const totalVotes = votes.reduce((a, b) => a + b, 0);
    const lenGCs = gcIds.length;
    const numEligible = cnStakingBalances.filter((x) => x.gt(MILLIONS.FIVE)).length;

    const summary = await stakingTracker.getTrackerSummary(trackerId);
    expect(summary).to.deep.equal([trackStart, trackEnd, lenGCs, totalVotes, numEligible]);

    const trackedGCs = await stakingTracker.getAllTrackedGCs(trackerId);
    expect(trackedGCs).to.deep.equal([gcIds, gcBalances, votes]);

    for (let i = 0; i < lenGCs; i++) {
      const gcBalance = await stakingTracker.getTrackedGCBalance(trackerId, gcIds[i]);
      expect(gcBalance).to.deep.equal([cnStakingBalances[i], gcBalances[i]]);
    }
  }

  let fixture: UnPromisify<ReturnType<typeof stakingTrackerV2TestFixture>>;
  let fakeCLRegistry: FakeContract<ICLRegistry>;
  let fakeWKaia: FakeContract<IERC20>;
  let clPoolA: FakeContract<ICLPool>;
  let clPoolB: FakeContract<ICLPool>;
  let clPoolC: FakeContract<ICLPool>;
  beforeEach(async function () {
    fixture = await loadFixture(stakingTrackerV2TestFixture);

    fakeCLRegistry = await smock.fake<ICLRegistry>(ICLRegistry__factory.abi);
    fakeWKaia = await smock.fake<IERC20>(IERC20__factory.abi);

    await registerContract(fixture.registry, "CLRegistry", fakeCLRegistry.address);
    await registerContract(fixture.registry, "WrappedKaia", fakeWKaia.address);

    clPoolA = await smock.fake<ICLPool>(ICLPool__factory.abi);
    clPoolB = await smock.fake<ICLPool>(ICLPool__factory.abi);
    clPoolC = await smock.fake<ICLPool>(ICLPool__factory.abi);

    clPoolA.stakingTracker.returns(fixture.stakingTracker.address);
    clPoolB.stakingTracker.returns(fixture.stakingTracker.address);
    clPoolC.stakingTracker.returns(fixture.stakingTracker.address);

    // Ignore nodeIds and CLStakings.
    fakeCLRegistry.getAllCLs.returns([[], [699, 700, 701], [clPoolA.address, clPoolB.address, clPoolC.address], []]);
  });

  describe("StakingTrackerV2 Initialize", function () {
    it("Check staking tracker constants", async function () {
      const { stakingTracker, contractValidator } = fixture;

      expect(await stakingTracker.CONTRACT_TYPE()).to.equal("StakingTracker");
      expect(await stakingTracker.VERSION()).to.equal(1);
      expect(await stakingTracker.ADDRESS_BOOK_ADDRESS()).to.equal(ABOOK_ADDRESS);
      expect(await stakingTracker.REGISTRY_ADDRESS()).to.equal(REGISTRY_ADDRESS);
      expect(await stakingTracker.MIN_STAKE()).to.equal(MILLIONS.FIVE);
      expect(await stakingTracker.owner()).to.equal(contractValidator.address);
    });
  });
  describe("Create a tracker", function () {
    it("#createTracker: no CLRegistry", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { SIX } = MILLIONS;

      await registerContract(fixture.registry, "CLRegistry", hre.ethers.constants.AddressZero);

      await expect(stakingTracker.createTracker(trackStart, trackEnd))
        .to.emit(stakingTracker, "CreateTracker")
        .withArgs(1, trackStart, trackEnd, [700, 701, 702]);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [SIX, SIX, SIX],
        [1, 1, 1]
      );
    });
    it("#createTracker: no wKaia so do not track wKaia balance", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { SIX, TWO, FIVE } = MILLIONS;

      await registerContract(fixture.registry, "WrappedKaia", hre.ethers.constants.AddressZero);

      fakeWKaia.balanceOf.whenCalledWith(clPoolB.address).returns(TWO);
      fakeWKaia.balanceOf.whenCalledWith(clPoolC.address).returns(FIVE);

      await expect(stakingTracker.createTracker(trackStart, trackEnd))
        .to.emit(stakingTracker, "CreateTracker")
        .withArgs(1, trackStart, trackEnd, [700, 701, 702]);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [SIX, SIX, SIX], // No wKaia balance reflected since wKaia is not active yet.
        [1, 1, 1]
      );
    });
    it("#createTracker: no CLPool", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { SIX } = MILLIONS;

      await expect(stakingTracker.createTracker(trackStart, trackEnd))
        .to.emit(stakingTracker, "CreateTracker")
        .withArgs(1, trackStart, trackEnd, [700, 701, 702]);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [SIX, SIX, SIX],
        [1, 1, 1]
      );
    });
    it("#createTracker: do not track CLPool with invalid StakingTracker address", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { SIX } = MILLIONS;

      clPoolB.stakingTracker.returns(hre.ethers.constants.AddressZero);

      await expect(stakingTracker.createTracker(trackStart, trackEnd))
        .to.emit(stakingTracker, "CreateTracker")
        .withArgs(1, trackStart, trackEnd, [700, 701, 702]);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [SIX, SIX, SIX],
        [1, 1, 1]
      );

      expect(await stakingTracker.isCLPool(1, clPoolB.address)).to.equal(false);
      expect(await stakingTracker.isCLPool(1, clPoolC.address)).to.equal(true);
      expect(await stakingTracker.stakingToGCId(1, clPoolB.address)).to.equal(0);
      expect(await stakingTracker.stakingToGCId(1, clPoolC.address)).to.equal(701);
    });
    it("#createTracker: do not track CLPool with non-existent gcId", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { SIX } = MILLIONS;

      // clPoolB is assigned gcId 800, which is not a valid gcId
      fakeCLRegistry.getAllCLs.returns([[], [699, 800, 701], [clPoolA.address, clPoolB.address, clPoolC.address], []]);

      await expect(stakingTracker.createTracker(trackStart, trackEnd))
        .to.emit(stakingTracker, "CreateTracker")
        .withArgs(1, trackStart, trackEnd, [700, 701, 702]);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [SIX, SIX, SIX],
        [1, 1, 1]
      );

      expect(await stakingTracker.isCLPool(1, clPoolB.address)).to.equal(false);
      expect(await stakingTracker.isCLPool(1, clPoolC.address)).to.equal(true);
      expect(await stakingTracker.stakingToGCId(1, clPoolB.address)).to.equal(0);
      expect(await stakingTracker.stakingToGCId(1, clPoolC.address)).to.equal(701);
    });
    it("#createTracker: two CLPools", async function () {
      const { stakingTracker, cnStakingV2B, cnStakingV2C, cnStakingV2D, trackStart, trackEnd } = fixture;
      const { TWO, FIVE, SIX, EIGHT, ELEVEN } = MILLIONS;

      fakeWKaia.balanceOf.whenCalledWith(clPoolA.address).returns(TWO);
      fakeWKaia.balanceOf.whenCalledWith(clPoolB.address).returns(TWO);
      fakeWKaia.balanceOf.whenCalledWith(clPoolC.address).returns(FIVE);

      await expect(stakingTracker.createTracker(trackStart, trackEnd))
        .to.emit(stakingTracker, "CreateTracker")
        .withArgs(1, trackStart, trackEnd, [700, 701, 702]);

      const gcBalances = await stakingTracker.getTrackedGCBalance(1, 699);
      // 699 is not tracked
      expect(gcBalances).to.deep.equal([0, 0]);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [EIGHT, ELEVEN, SIX],
        [1, 2, 1]
      );

      expect(await stakingTracker.isCLPool(1, clPoolA.address)).to.equal(false);
      expect(await stakingTracker.isCLPool(1, clPoolB.address)).to.equal(true);
      expect(await stakingTracker.isCLPool(1, clPoolC.address)).to.equal(true);
      expect(await stakingTracker.isCLPool(1, cnStakingV2B.address)).to.equal(false);
      expect(await stakingTracker.isCLPool(1, cnStakingV2C.address)).to.equal(false);
      expect(await stakingTracker.isCLPool(1, cnStakingV2D.address)).to.equal(false);
      expect(await stakingTracker.stakingToGCId(1, cnStakingV2B.address)).to.equal(700);
      expect(await stakingTracker.stakingToGCId(1, cnStakingV2C.address)).to.equal(701);
      expect(await stakingTracker.stakingToGCId(1, cnStakingV2D.address)).to.equal(702);
      expect(await stakingTracker.stakingToGCId(1, clPoolB.address)).to.equal(700);
      expect(await stakingTracker.stakingToGCId(1, clPoolC.address)).to.equal(701);
    });
    it("#createTracker: CnStaking has less than 5M", async function () {
      const { stakingTracker, cnStakingV2C, trackStart, trackEnd } = fixture;
      const { TWO, SIX, FOUR, FIVE, EIGHT, NINE, TEN } = MILLIONS;

      fakeWKaia.balanceOf.whenCalledWith(clPoolB.address).returns(TWO);
      fakeWKaia.balanceOf.whenCalledWith(clPoolC.address).returns(FIVE);

      // Force set the balance of CnStakingV2C to 4M
      await setBalance(cnStakingV2C.address, FOUR);

      await expect(stakingTracker.createTracker(trackStart, trackEnd))
        .to.emit(stakingTracker, "CreateTracker")
        .withArgs(1, trackStart, trackEnd, [700, 701, 702]);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, FOUR, SIX],
        [EIGHT, NINE, SIX],
        [1, 0, 1] // CnStakingV2C is not eligible since its CnStaking balance is less than 5M
      );

      // The addition of CLPoolC's balance should not affect the votes of CnStakingV2C
      fakeWKaia.balanceOf.whenCalledWith(clPoolC.address).returns(TEN);
      await expect(stakingTracker.refreshStake(clPoolC.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 701, clPoolC.address, TEN, TEN.add(FOUR), 0, 2);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, FOUR, SIX],
        [EIGHT, TEN.add(FOUR), SIX],
        [1, 0, 1]
      );
    });
  });
  describe("Refresh stake", function () {
    beforeEach(async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { TWO, FIVE } = MILLIONS;

      fakeWKaia.balanceOf.whenCalledWith(clPoolB.address).returns(TWO);
      fakeWKaia.balanceOf.whenCalledWith(clPoolC.address).returns(FIVE);

      // We have 1, 2, 1 votes for B, C, D respectively
      await stakingTracker.createTracker(trackStart, trackEnd);
    });
    it("#refreshStake: From CnStaking, eligible -> eligible", async function () {
      const { stakingTracker, cnStakingV2B, trackStart, trackEnd } = fixture;
      const { SIX, EIGHT, TEN, ELEVEN } = MILLIONS;

      await setBalance(cnStakingV2B.address, EIGHT); // 6M -> 8M, votes: 1 -> 2

      await expect(stakingTracker.refreshStake(cnStakingV2B.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, cnStakingV2B.address, EIGHT, TEN, 2, 5);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [EIGHT, SIX, SIX],
        [TEN, ELEVEN, SIX],
        [2, 2, 1]
      );
    });
    it("#refreshStake: From CnStaking, eligible -> ineligible", async function () {
      const { stakingTracker, cnStakingV2B, trackStart, trackEnd } = fixture;
      const { FOUR, SIX, ELEVEN } = MILLIONS;
      await setBalance(cnStakingV2B.address, FOUR); // 6M -> 4M, votes: 1 -> 0

      await expect(stakingTracker.refreshStake(cnStakingV2B.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, cnStakingV2B.address, FOUR, SIX, 0, 2);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [FOUR, SIX, SIX],
        [SIX, ELEVEN, SIX],
        [0, 1, 1] // the votes of CnStakingV2C are also capped by 1
      );
    });
    it("#refreshStake: From CnStaking, ineligible -> eligible", async function () {
      const { stakingTracker, cnStakingV2B, trackStart, trackEnd } = fixture;
      const { FOUR, SIX, EIGHT, TEN, ELEVEN } = MILLIONS;

      await setBalance(cnStakingV2B.address, FOUR); // 6M -> 4M, votes: 1 -> 0
      await expect(stakingTracker.refreshStake(cnStakingV2B.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, cnStakingV2B.address, FOUR, SIX, 0, 2);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [FOUR, SIX, SIX],
        [SIX, ELEVEN, SIX],
        [0, 1, 1]
      );

      await setBalance(cnStakingV2B.address, EIGHT); // 4M -> 8M, votes: 0 -> 2
      await expect(stakingTracker.refreshStake(cnStakingV2B.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, cnStakingV2B.address, EIGHT, TEN, 2, 5);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [EIGHT, SIX, SIX],
        [TEN, ELEVEN, SIX],
        [2, 2, 1]
      );
    });
    it("#refreshStake: From CnStaking, ineligible -> ineligible", async function () {
      const { stakingTracker, cnStakingV2B, trackStart, trackEnd } = fixture;
      const { TWO, FOUR, SIX, ELEVEN } = MILLIONS;
      await setBalance(cnStakingV2B.address, FOUR); // 6M -> 4M, votes: 1 -> 0

      await expect(stakingTracker.refreshStake(cnStakingV2B.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, cnStakingV2B.address, FOUR, SIX, 0, 2);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [FOUR, SIX, SIX],
        [SIX, ELEVEN, SIX],
        [0, 1, 1]
      );

      await setBalance(cnStakingV2B.address, TWO); // 4M -> 2M, votes: 0 -> 0
      await expect(stakingTracker.refreshStake(cnStakingV2B.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, cnStakingV2B.address, TWO, FOUR, 0, 2);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [TWO, SIX, SIX],
        [FOUR, ELEVEN, SIX],
        [0, 1, 1]
      );
    });
    it("#refreshStake: From CLPool, votes unchanged", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { THREE, SIX, NINE, ELEVEN } = MILLIONS;

      fakeWKaia.balanceOf.whenCalledWith(clPoolB.address).returns(THREE); // 2M -> 3M, votes: 1 -> 1
      await expect(stakingTracker.refreshStake(clPoolB.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, clPoolB.address, THREE, NINE, 1, 4);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [NINE, ELEVEN, SIX],
        [1, 2, 1]
      );
    });
    it("#refreshStake: From CLPool, votes increased", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { FOUR, SIX, TEN, ELEVEN } = MILLIONS;

      fakeWKaia.balanceOf.whenCalledWith(clPoolB.address).returns(FOUR); // 2M -> 4M, votes: 1 -> 2
      await expect(stakingTracker.refreshStake(clPoolB.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, clPoolB.address, FOUR, TEN, 2, 5);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [TEN, ELEVEN, SIX],
        [2, 2, 1]
      );
    });
    it("#refreshStake: From CLPool, votes decreased", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;
      const { TWO, FOUR, SIX, EIGHT, TEN, ELEVEN } = MILLIONS;

      fakeWKaia.balanceOf.whenCalledWith(clPoolB.address).returns(FOUR); // 2M -> 4M, votes: 1 -> 2
      await expect(stakingTracker.refreshStake(clPoolB.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, clPoolB.address, FOUR, TEN, 2, 5);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [TEN, ELEVEN, SIX],
        [2, 2, 1]
      );

      fakeWKaia.balanceOf.whenCalledWith(clPoolB.address).returns(TWO); // 4M -> 2M, votes: 2 -> 1
      await expect(stakingTracker.refreshStake(clPoolB.address))
        .to.emit(stakingTracker, "RefreshStake")
        .withArgs(1, 700, clPoolB.address, TWO, EIGHT, 1, 4);

      await verifyTrackerState(
        stakingTracker,
        1,
        trackStart,
        trackEnd,
        [700, 701, 702],
        [SIX, SIX, SIX],
        [EIGHT, ELEVEN, SIX],
        [1, 2, 1]
      );
    });
  });
});
