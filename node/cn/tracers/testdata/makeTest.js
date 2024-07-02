// Create a testdata JSON out of a transaction hash (tx)
var makeTest = function(tx, tracerConfig) {
  var block  = eth.getBlock(eth.getTransaction(tx).blockHash);

  // Generate the genesis block from the block, transaction and prestate data
  // as if the blockchain was born just before the transaction.
  // Note that only config and alloc are used in the tracers_test.
  var alloc = debug.traceTransaction(tx, {tracer: "prestateTracer"});
  for (var key in alloc) {
    alloc[key].nonce = alloc[key].nonce.toString();
  }
  var config = klay.getChainConfig(block.number);
  var genesis = {
    alloc,
    config,
  };

  // Fill the prestate for feePayer account, if the tx is fee delegated.
  // You must use an archive node for fee delegated tx.
  var klayTx = klay.getTransaction(tx);
  if (klayTx.feePayer) {
    var account = klay.getAccount(klayTx.feePayer, klayTx.blockNumber - 1).account;
    genesis.alloc[klayTx.feePayer] = {
      balance: account.balance,
      nonce: account.nonce,
      // feePayer's code and storage, even if exists, are irrelevant executing current tx.
      code: "0x",
      storage: {},
    };
  }

  // Collect the necessary block context.
  var context = {
    mixHash: block.mixHash,
    number: block.number.toString(),
    timestamp: block.timestamp.toString(),
    blockScore: block.difficulty.toString(),
  }
  if (block.baseFeePerGas) {
    context.baseFeePerGas = block.baseFeePerGas.toString();
  }

  // Generate the correct call trace
  var result = debug.traceTransaction(tx, tracerConfig);

  console.log(JSON.stringify({
    _comment: `chainId ${parseInt(eth.chainId())} txHash ${tx}`, // for the record
    genesis: genesis,
    context: context,
    // Must use klay.getRawTransaction(tx) instead of eth.getRawTransaction(tx) to preserve eth envelope 0x78.
    input:  klay.getRawTransaction(tx),
    result: result,
  }, null, 2));
}
