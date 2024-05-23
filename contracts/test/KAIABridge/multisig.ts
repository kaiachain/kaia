import { expect } from "chai";

describe("[Multisig Test]", function () {
  let bridge;
  let operator;
  let guardian;
  let judge;
  let minOperatorRequiredConfirm;
  let minGuardianRequiredConfirm;
  let operator1;
  let operator2;
  let operator3;
  let operator4;
  let guardian1;
  let guardian2;
  let guardian3;
  let guardian4;
  let operatorCandidate1;
  let guardianCandidate1;
  let judgeCandidate1;
  let judgeCandidate2;
  let judge1;
  let judge2;
  let txID;
  let guardianTxID;

  let sender;
  let receiver;
  let amount;
  let seq;
  let maxTryTransfer;

  beforeEach(async function () {
    upgrades.silenceWarnings()
    sender   = "0x0000000000000000000000000000000000000123";
    receiver = "0x0000000000000000000000000000000000000456";
    amount = 1;
    seq = 1;

    const [
      _operator1, _operator2, _operator3, _operator4, p1,
      _guardian1, _guardian2, _guardian3, _guardian4, p2,
      _judge1, _judge2, _judgeCandidate1, _judgeCandidate2
    ] = await ethers.getSigners();

    operator1 = _operator1;
    operator2 = _operator2;
    operator3 = _operator3;
    operator4 = _operator4;
    guardian1 = _guardian1;
    guardian2 = _guardian2;
    guardian3 = _guardian3;
    guardian4 = _guardian4;

    judge1 = _judge1;
    judge2 = _judge2;
    judgeCandidate1 = _judgeCandidate1;
    judgeCandidate2 = _judgeCandidate2;

    operatorCandidate1 = p1;
    guardianCandidate1 = p2;
    minOperatorRequiredConfirm = 3;
    minGuardianRequiredConfirm = 3;
    maxTryTransfer = 0;

    const accs = [operator1, operator2, operator3, operator4, guardian1, guardian2, guardian3, guardian4];
    for (let acc of accs) {
      await hre.network.provider.send("hardhat_setBalance", [
        acc.address,
        "0x1000000000000000000000000000000000000",
      ]);
    }

    const guardianFactory = await ethers.getContractFactory("Guardian");
    guardian = await upgrades.deployProxy(guardianFactory, [[
      guardian1.address,
      guardian2.address,
      guardian3.address,
      guardian4.address,
    ], minGuardianRequiredConfirm]);

    const operatorFactory = await ethers.getContractFactory("Operator");
    operator = await upgrades.deployProxy(operatorFactory, [[
      operator1.address,
      operator2.address,
      operator3.address,
      operator4.address,
    ], guardian.address, minOperatorRequiredConfirm]);

    const judgeFactory = await ethers.getContractFactory("Judge");
    judge = await upgrades.deployProxy(judgeFactory, [[
      judge1.address, judge2.address
    ], guardian.address, 2]);

    const bridgeFactory = await ethers.getContractFactory("KAIABridge");
    bridge = await upgrades.deployProxy(bridgeFactory, [
      operator.address, guardian.address, judge.address, maxTryTransfer
    ]);

    let rawTxData = (await operator.populateTransaction.changeBridge(bridge.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(1);
    await expect(guardian.connect(guardian3).confirmTransaction(1))
      .to.emit(operator, "ChangeBridge");
    guardianTxID = 2;

    txID = 1;
  })

  it("#add operator", async function () {
    const initialOperatorLen = 4;
    expect((await operator.getOperators()).length).to.equal(initialOperatorLen);

    const rawTxData = (await operator.populateTransaction.addOperator(operatorCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    expect((await operator.getOperators()).length).to.equal(initialOperatorLen + 1);
  });

  it("#add operator - authentication failure", async function () {
    const initialOperatorLen = 4;
    expect((await operator.getOperators()).length).to.equal(initialOperatorLen);

    const rawTxData = (await operator.populateTransaction.addOperator(operatorCandidate1.address)).data;
    await operator.connect(operator1).submitTransaction(operator.address, rawTxData, 0);
    await operator.connect(operator2).confirmTransaction(txID);
    await expect(operator.connect(operator3).confirmTransaction(txID))
      .to.revertedWith("KAIA::Operator: Sender is not guardian contract");
  });

  it("#remove operator", async function () {
    const initialOperatorLen = 4;
    expect((await operator.getOperators()).length).to.equal(initialOperatorLen);

    // try to delete unknown address
    let rawTxData = (await operator.populateTransaction.removeOperator(operatorCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await expect(guardian.connect(guardian3).confirmTransaction(guardianTxID))
      .to.be.revertedWith("KAIA::Operator: Not an operator");

    // now try to remove known validator
    const newTxID = guardianTxID + 1;
    rawTxData = (await operator.populateTransaction.removeOperator(operator3.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(newTxID);
    await expect(guardian.connect(guardian3).confirmTransaction(newTxID))
      .to.emit(guardian, "Execution").withArgs(newTxID);

    expect((await operator.getOperators()).length).to.equal(initialOperatorLen - 1);
  });

  it("#operator removal may change the number of required confirmation", async function () {
    let rawTxData = (await operator.populateTransaction.changeRequirement(4)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);
    expect(await operator.minOperatorRequiredConfirm()).to.equal(4);

    rawTxData = (await operator.populateTransaction.removeOperator(operator3.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID + 1);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID + 1);
    expect((await operator.getOperators()).length).to.equal(3);

    expect(await guardian.minGuardianRequiredConfirm()).to.equal(3);
  });

  it("#replace operator", async function () {
    const initialOperatorLen = 4;
    expect((await operator.getOperators()).length).to.equal(initialOperatorLen);

    const rawTxData = (await operator.populateTransaction.replaceOperator(operator3.address, operatorCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await expect(guardian.connect(guardian3).confirmTransaction(guardianTxID))
      .to.emit(guardian, "Execution").withArgs(guardianTxID);

    expect((await operator.getOperators())).to.contain(operatorCandidate1.address);
  });

  it("#change guardian requirement - guardian", async function () {
    const initialOperatorLen = 4;
    expect((await operator.getOperators()).length).to.equal(initialOperatorLen);

    let rawTxData = (await guardian.populateTransaction.changeRequirement(2)).data;
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    // Now two votes can be effective
    const newTxID = guardianTxID + 1;
    rawTxData = (await operator.populateTransaction.addOperator(operatorCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await expect(guardian.connect(guardian3).confirmTransaction(newTxID))
      .to.emit(guardian, "Execution").withArgs(newTxID);

    expect((await operator.getOperators()).length).to.equal(initialOperatorLen + 1);
  });

  it("#change operator requirement - operator", async function () {
    expect(await operator.minOperatorRequiredConfirm()).to.equal(3);

    let rawTxData = (await operator.populateTransaction.changeRequirement(2)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    expect(await operator.minOperatorRequiredConfirm()).to.equal(2);
  });

  it("#change judge requirement - judge", async function () {
    expect(await judge.minJudgeRequiredConfirm()).to.equal(2);

    let rawTxData = (await judge.populateTransaction.changeRequirement(1)).data;
    await guardian.connect(guardian1).submitTransaction(judge.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await expect(guardian.connect(guardian3).confirmTransaction(guardianTxID))
      .to.emit(judge, "RequirementChange")

    expect(await judge.minJudgeRequiredConfirm()).to.equal(1);
  });

  it("#add guardian", async function () {
    const initGuardianLen = 4;
    expect((await guardian.getGuardians()).length).to.equal(initGuardianLen);

    const rawTxData = (await guardian.populateTransaction.addGuardian(guardianCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    expect((await guardian.getGuardians()).length).to.equal(initGuardianLen + 1);
  });

  it("#remove guardian", async function () {
    const initGuardianLen = 4;
    expect((await guardian.getGuardians()).length).to.equal(initGuardianLen);

    const rawTxData = (await guardian.populateTransaction.removeGuardian(guardian3.address)).data;
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    expect((await guardian.getGuardians()).length).to.equal(initGuardianLen - 1);
  });

  it("#guardian removal may change the number of required confirmation", async function () {
    let rawTxData = (await guardian.populateTransaction.changeRequirement(4)).data;
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);
    expect(await guardian.minGuardianRequiredConfirm()).to.equal(4);

    rawTxData = (await guardian.populateTransaction.removeGuardian(guardian3.address)).data;
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID + 1);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID + 1);
    await guardian.connect(guardian4).confirmTransaction(guardianTxID + 1);
    expect((await guardian.getGuardians()).length).to.equal(3);
  });

  it("#replace guardian", async function () {
    const initGuardianLen = 4;
    expect((await guardian.getGuardians()).length).to.equal(initGuardianLen);

    const rawTxData = (await guardian.populateTransaction.replaceGuardian(guardian3.address, guardianCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    expect((await guardian.getGuardians())).to.contain(guardianCandidate1.address);
  });

  it("#add judge", async function () {
    const initJudgeLen = 2;
    expect((await judge.getJudges()).length).to.equal(initJudgeLen);

    let rawTxData = (await judge.populateTransaction.addJudge(judgeCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(judge.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);
    expect((await judge.getJudges()).length).to.equal(initJudgeLen + 1);

    rawTxData = (await judge.populateTransaction.addJudge(judgeCandidate2.address)).data;
    await guardian.connect(guardian1).submitTransaction(judge.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID + 1);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID + 1);
    expect((await judge.getJudges()).length).to.equal(initJudgeLen + 2);
  });

  it("#remove judge", async function () {
    const initJudgeLen = 2;
    expect((await judge.getJudges()).length).to.equal(initJudgeLen);

    let rawTxData = (await judge.populateTransaction.addJudge(judgeCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(judge.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);
    expect((await judge.getJudges()).length).to.equal(initJudgeLen + 1);

    rawTxData = (await judge.populateTransaction.removeJudge(judgeCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(judge.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID + 1);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID + 1);
    expect((await judge.getJudges()).length).to.equal(initJudgeLen);
  });

  it("#replace judge", async function () {
    expect(await judge.getJudges()).contains(judge1.address);

    let rawTxData = (await judge.populateTransaction.replaceJudge(judge1.address, judgeCandidate1.address)).data;
    await guardian.connect(guardian1).submitTransaction(judge.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);
    expect(await judge.getJudges()).contains(judgeCandidate1.address);
  });

  it("#change guardian in judge contract", async function () {
    expect(await judge.guardian()).to.equal(guardian.address);

    const newGuardian = "0x0000000000000000000000000000000000000123";
    let rawTxData = (await judge.populateTransaction.changeGuardian(newGuardian)).data;
    await guardian.connect(guardian1).submitTransaction(judge.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    expect(await judge.guardian()).to.equal(newGuardian);
  });

  it("#revoke transaction - guardian", async function () {
    let rawTxData = (await operator.populateTransaction.changeRequirement(2)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);

    // Owner3 did commit confirmation before
    await expect(guardian.connect(guardian3).revokeConfirmation(guardianTxID))
      .to.be.revertedWith("KAIA::Guardian: No confirmation was committed yet");

    expect((await guardian.getConfirmations(guardianTxID)).length).to.equal(2);
    await expect(guardian.connect(guardian2).revokeConfirmation(guardianTxID))
      .to.emit(guardian, "Revocation");
    expect((await guardian.getConfirmations(guardianTxID)).length).to.equal(1);
    await expect(guardian.connect(guardian1).revokeConfirmation(guardianTxID))
      .to.emit(guardian, "Revocation");
    expect((await guardian.getConfirmations(guardianTxID)).length).to.equal(0);
  });

  it("#revoke transaction - operator", async function () {
    let rawTxData = (await operator.populateTransaction.changeRequirement(2)).data;
    await operator.connect(operator1).submitTransaction(operator.address, rawTxData, 0);
    await operator.connect(operator2).confirmTransaction(txID);

    // Owner3 did commit confirmation before
    await expect(operator.connect(operator3).revokeConfirmation(txID))
      .to.be.revertedWith("KAIA::Operator: No confirmation was committed yet");
    expect((await operator.getConfirmations(txID)).length).to.equal(2);
    await expect(operator.connect(operator2).revokeConfirmation(txID))
      .to.emit(operator, "Revocation");
    expect((await operator.getConfirmations(txID)).length).to.equal(1);
    await expect(operator.connect(operator1).revokeConfirmation(txID))
      .to.emit(operator, "Revocation");
    expect((await operator.getConfirmations(txID)).length).to.equal(0);
  });
  
  it("#revoke transaction - judge", async function () {
    let rawTxData = (await judge.populateTransaction.changeRequirement(1)).data;
    await judge.connect(judge1).submitTransaction(judge.address, rawTxData, 0);

    await expect(judge.connect(judge2).revokeConfirmation(txID))
      .to.be.revertedWith("KAIA::Judge: No confirmation was committed yet");
    expect((await judge.getConfirmations(txID)).length).to.equal(1);
    await expect(judge.connect(judge1).revokeConfirmation(txID))
      .to.emit(judge, "Revocation");
    expect((await judge.getConfirmations(txID)).length).to.equal(0);
  });

  it("#change guardian in operator contract", async function () {
    expect(await operator.guardian()).to.equal(guardian.address);

    const newGuardian = "0x0000000000000000000000000000000000000123";
    let rawTxData = (await operator.populateTransaction.changeGuardian(newGuardian)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    expect(await operator.guardian()).to.equal(newGuardian);
  });

  it("#duplicated confirmation", async function () {
    let rawTxData = (await operator.populateTransaction.changeRequirement(2)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await expect(guardian.connect(guardian2).confirmTransaction(guardianTxID))
      .to.be.revertedWith("KAIA::Guardian: Transaction was already confirmed");
  });

  it("#get transaction count - guardian", async function () {
    let [nPending, nExecuted] = await guardian.getTransactionCount(true, true);
    expect(nPending).to.equal(0);
    expect(nExecuted).to.equal(1);

    let rawTxData = (await guardian.populateTransaction.changeRequirement(2)).data;

    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);

    [nPending, nExecuted] = await guardian.getTransactionCount(true, true);
    await expect(nPending).to.equal(2);
    await expect(nExecuted).to.equal(2);
  });

  it("#get transaction count - opertaor", async function () {
    let [nPending, nExecuted] = await operator.getTransactionCount(true, true);
    await expect(nPending).to.equal(0);
    await expect(nExecuted).to.equal(0);

    let provision = [seq, sender, receiver, amount];
    let rawTxData = (await bridge.populateTransaction.provision(provision)).data;

    await operator.connect(operator1).submitTransaction(bridge.address, rawTxData, 0);
    await operator.connect(operator2).confirmTransaction(txID);
    await operator.connect(operator3).confirmTransaction(txID);

    provision = [seq + 1, sender, receiver, amount];
    rawTxData = (await bridge.populateTransaction.provision(provision)).data;
    await operator.connect(operator1).submitTransaction(bridge.address, rawTxData, 0);
    provision = [seq + 2, sender, receiver, amount];
    rawTxData = (await bridge.populateTransaction.provision(provision)).data;
    await operator.connect(operator1).submitTransaction(bridge.address, rawTxData, 0);

    [nPending, nExecuted] = await operator.getTransactionCount(true, true);
    await expect(nPending).to.equal(2);
    await expect(nExecuted).to.equal(1);
  });

  it("#get transaction count - judge", async function () {
    let [nPending, nExecuted] = await judge.getTransactionCount(true, true);
    await expect(nPending).to.equal(0);
    await expect(nExecuted).to.equal(0);

    let rawTxData = (await bridge.populateTransaction.holdClaim(1)).data;
    await judge.connect(judge1).submitTransaction(bridge.address, rawTxData, 0);
    await judge.connect(judge2).confirmTransaction(txID);

    rawTxData = (await bridge.populateTransaction.holdClaim(2)).data;
    await judge.connect(judge1).submitTransaction(bridge.address, rawTxData, 0);
    rawTxData = (await bridge.populateTransaction.holdClaim(3)).data;
    await judge.connect(judge1).submitTransaction(bridge.address, rawTxData, 0);

    [nPending, nExecuted] = await judge.getTransactionCount(true, true);
    await expect(nPending).to.equal(2);
    await expect(nExecuted).to.equal(1);
  });

  it("#get confirmation count - guardian", async function () {
    expect(await guardian.getConfirmationCount(0)).to.equal(0);
    let rawTxData = (await guardian.populateTransaction.changeRequirement(2)).data;
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    expect(await guardian.getConfirmationCount(guardianTxID)).to.equal(2);
  });

  it("#get confirmation count - judge", async function () {
    expect(await judge.getConfirmationCount(txID)).to.equal(0);

    let rawTxData = (await bridge.populateTransaction.holdClaim(1)).data;
    await judge.connect(judge1).submitTransaction(bridge.address, rawTxData, 0);
    await judge.connect(judge2).confirmTransaction(txID);

    expect(await judge.getConfirmationCount(txID)).to.equal(2);
  });

  it("#get transaction IDs - guardian", async function () {
    let [nPending, nExecuted] = await guardian.getTransactionIds(0, 100, true, true);
    await expect(nPending.length).to.equal(0);
    await expect(nExecuted.length).to.equal(1);

    let rawTxData = (await guardian.populateTransaction.changeRequirement(2)).data;

    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID);

    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);
    await guardian.connect(guardian1).submitTransaction(guardian.address, rawTxData, 0);

    [nPending, nExecuted] = await guardian.getTransactionIds(0, 100, true, true);
    await expect(nPending.length).to.equal(2);
    await expect(nExecuted.length).to.equal(2);
  });

  it("#get transaction IDs - operator", async function () {
    let [nPending, nExecuted] = await operator.getTransactionIds(0, 100, true, true);
    await expect(nPending.length).to.equal(0);
    await expect(nExecuted.length).to.equal(0);

    let provision = [seq, sender, receiver, amount];
    let rawTxData = (await bridge.populateTransaction.provision(provision)).data;

    await operator.connect(operator1).submitTransaction(bridge.address, rawTxData, 0);
    await operator.connect(operator2).confirmTransaction(txID);
    await operator.connect(operator3).confirmTransaction(txID);

    provision = [seq + 1, sender, receiver, amount];
    rawTxData = (await bridge.populateTransaction.provision(provision)).data;
    await operator.connect(operator1).submitTransaction(bridge.address, rawTxData, 0);
    provision = [seq + 2, sender, receiver, amount];
    rawTxData = (await bridge.populateTransaction.provision(provision)).data;
    await operator.connect(operator1).submitTransaction(bridge.address, rawTxData, 0);

    [nPending, nExecuted] = await operator.getTransactionIds(0, 100, true, true);
    await expect(nPending.length).to.equal(2);
    await expect(nExecuted.length).to.equal(1);
  });

  it("#get transaction IDs - judge", async function () {
    let [nPending, nExecuted] = await judge.getTransactionIds(0, 100, true, true);
    await expect(nPending.length).to.equal(0);
    await expect(nExecuted.length).to.equal(0);

    let rawTxData = (await bridge.populateTransaction.holdClaim(1)).data;
    await judge.connect(judge1).submitTransaction(bridge.address, rawTxData, 0);
    await judge.connect(judge2).confirmTransaction(txID);

    rawTxData = (await bridge.populateTransaction.holdClaim(2)).data;
    await judge.connect(judge1).submitTransaction(bridge.address, rawTxData, 0);
    rawTxData = (await bridge.populateTransaction.holdClaim(3)).data;
    await judge.connect(judge1).submitTransaction(bridge.address, rawTxData, 0);

    [nPending, nExecuted] = await judge.getTransactionIds(0, 100, true, true);
    await expect(nPending.length).to.equal(2);
    await expect(nExecuted.length).to.equal(1);
  });

  it("#try failed execution again when it becomes successful", async function () {
    // Failed to execute `addOperator` execution
    let addOperatorTx = (await operator.populateTransaction.addOperator(operator4.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, addOperatorTx, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID);
    await expect(guardian.connect(guardian3).confirmTransaction(guardianTxID))
      .to.be.revertedWith("KAIA::Operator: The address must not be operator");

    let removeOperatorTx = (await operator.populateTransaction.removeOperator(operator4.address)).data;
    await guardian.connect(guardian1).submitTransaction(operator.address, removeOperatorTx, 0);
    await guardian.connect(guardian2).confirmTransaction(guardianTxID + 1);
    await guardian.connect(guardian3).confirmTransaction(guardianTxID + 1);

    await expect(guardian.connect(guardian3).confirmTransaction(guardianTxID))
      .to.emit(guardian, "Execution");
  });
});
