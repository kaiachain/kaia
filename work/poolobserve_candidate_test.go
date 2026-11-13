package work

import (
	"errors"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/kzg4844"
	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/kaiachain/kaia/kaiax/gov"
	govmock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work/builder"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/prometheus/client_golang/prometheus"
)

func candidateCounter(t *testing.T, r *prometheus.Registry, name, label, value string) float64 {
	t.Helper()
	fs, err := r.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fs {
		if f.GetName() != name {
			continue
		}
		for _, m := range f.Metric {
			if label == "" && len(m.Label) == 0 {
				return m.GetCounter().GetValue()
			}
			for _, l := range m.Label {
				if l.GetName() == label && l.GetValue() == value {
					if m.Histogram != nil {
						return float64(m.Histogram.GetSampleCount())
					}
					return m.GetCounter().GetValue()
				}
			}
		}
	}
	t.Fatalf("missing %s{%s=%s}", name, label, value)
	return 0
}

type candidateFilter struct {
	builder.TxBundlingModule
	remove common.Address
}

func (f candidateFilter) FilterTxs(txs map[common.Address]types.Transactions) { delete(txs, f.remove) }

func TestCandidatePreparationMetrics(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	reg := prometheus.NewRegistry()
	poolobserve.Register(reg)
	ctrl := gomock.NewController(t)
	chain := mocks.NewMockBlockChain(ctrl)
	chain.EXPECT().CurrentBlock().Return(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(0), Time: big.NewInt(0)})).AnyTimes()
	chain.EXPECT().Validator().Return(&fakeValidator{}).AnyTimes()
	chain.EXPECT().PrepareHeader(gomock.Any()).Return(errors.New("stop after preparation"))
	pool := mocks.NewMockTxPool(ctrl)
	low, high := common.Address{1}, common.Address{2}
	pool.EXPECT().Pending().Return(map[common.Address]types.Transactions{
		low:  {types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1), nil)},
		high: {types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(100000000000), nil)},
	}, nil)
	g := govmock.NewMockGovModule(ctrl)
	g.EXPECT().GetParamSet(uint64(1)).Return(*gov.GetDefaultGovernanceParamSet())
	w := &worker{config: &params.ChainConfig{MagmaCompatibleBlock: big.NewInt(1)}, backend: &testBackend{chain: chain, txPool: pool}, chain: chain, govModule: g, txBundlingModules: []builder.TxBundlingModule{candidateFilter{remove: high}}}
	w.commitNewWork()
	for phase, count := range map[string]float64{"pool_input": 2, "after_fee": 1, "after_modules": 0} {
		if got := candidateCounter(t, reg, "kaia_poolobs_candidate_transactions_total", "phase", phase); got != count {
			t.Fatalf("%s: got %v want %v", phase, got, count)
		}
		if got := candidateCounter(t, reg, "kaia_poolobs_candidate_preparations_total", "", ""); got != 1 {
			t.Fatalf("%s observations=%v", phase, got)
		}
	}
	if candidateCounter(t, reg, "kaia_poolobs_candidate_empty_total", "phase", "after_modules") != 1 {
		t.Fatal("module-produced empty candidate not observed")
	}
	if candidateCounter(t, reg, "kaia_poolobs_duration_seconds", "stage", "candidate_pending") != 1 {
		t.Fatal("pending operation not completed")
	}
}

// Inject execution outcomes at the real worker boundary, without simulating EVM behavior.
type candidateMetricsChain struct {
	BlockChain
	firstError         error
	firstFailedReceipt bool
	attempts           int
}

func (c *candidateMetricsChain) ApplyTransaction(_ *params.ChainConfig, _ *common.Address, _ *state.StateDB, _ *types.Header, _ *types.Transaction, _ *uint64, _ *vm.Config) (*types.Receipt, *vm.InternalTxTrace, error) {
	c.attempts++
	if c.attempts == 1 && c.firstError != nil {
		return nil, nil, c.firstError
	}
	if c.attempts == 1 && c.firstFailedReceipt {
		return &types.Receipt{Status: types.ReceiptStatusFailed}, nil, nil
	}
	return &types.Receipt{Status: types.ReceiptStatusSuccessful}, nil, nil
}

type failingGeneratedBundle struct{}

func (failingGeneratedBundle) FilterTxs(map[common.Address]types.Transactions) {}
func (failingGeneratedBundle) ExtractTxBundles([]*types.Transaction, []*builder.Bundle) []*builder.Bundle {
	entry := builder.NewTxOrGenFromGen(func(uint64) (*types.Transaction, error) { return nil, errors.New("generated transaction unavailable") }, common.Hash{99})
	return []*builder.Bundle{builder.NewBundle([]*builder.TxOrGen{entry}, common.Hash{}, false)}
}

func TestCandidateExecutionMetrics(t *testing.T) {
	for _, tc := range []struct {
		name     string
		input    int
		err      error
		included int
	}{{"empty", 0, nil, 0}, {"included", 2, nil, 2}, {"nonce_dependency", 2, blockchain.ErrNonceTooHigh, 0}, {"old_nonce", 2, blockchain.ErrNonceTooLow, 1}, {"generated_entry_not_input", 2, nil, 2}} {
		t.Run(tc.name, func(t *testing.T) {
			defer poolobserve.EnableForTesting()()
			reg := prometheus.NewRegistry()
			poolobserve.Register(reg)
			chainID := big.NewInt(1)
			signer := types.LatestSignerForChainID(chainID)
			key, err := crypto.GenerateKey()
			if err != nil {
				t.Fatal(err)
			}
			input := make(types.Transactions, tc.input)
			for i := range input {
				input[i], err = types.SignTx(types.NewTransaction(uint64(i), common.Address{}, big.NewInt(1), 21000, big.NewInt(25000000000), nil), signer, key)
				if err != nil {
					t.Fatal(err)
				}
			}
			db, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			task := &Task{config: &params.ChainConfig{ChainID: chainID}, signer: signer, state: db, header: &types.Header{Number: big.NewInt(1), Time: big.NewInt(1)}}
			pending := make(map[common.Address]types.Transactions)
			if len(input) > 0 {
				pending[crypto.PubkeyToAddress(key.PublicKey)] = input
			}
			txs := types.NewTransactionsByPriceAndNonce(signer, pending, nil)
			var modules []builder.TxBundlingModule
			if tc.name == "generated_entry_not_input" {
				modules = []builder.TxBundlingModule{failingGeneratedBundle{}}
			}
			task.ApplyTransactions(txs, &candidateMetricsChain{firstError: tc.err}, common.Address{}, modules)
			for phase, want := range map[string]int{"execution_input": tc.input, "included": tc.included} {
				if got := candidateCounter(t, reg, "kaia_poolobs_candidate_transactions_total", "phase", phase); got != float64(want) {
					t.Fatalf("%s got %v want %d", phase, got, want)
				}
			}
			for reason, executionErr := range map[string]error{"nonce_low": blockchain.ErrNonceTooLow, "nonce_high": blockchain.ErrNonceTooHigh} {
				want := float64(0)
				if tc.err == executionErr {
					want = 1
				}
				if got := candidateCounter(t, reg, "kaia_poolobs_candidate_outcomes_total", "reason", reason); got != want {
					t.Fatalf("%s got %v want %v", reason, got, want)
				}
			}
			if len(task.txs) != tc.included {
				t.Fatal("unexpected assembly")
			}
			if candidateCounter(t, reg, "kaia_poolobs_duration_seconds", "stage", "candidate_build") != 1 {
				t.Fatal("candidate operation not completed")
			}
		})
	}
}

type outcomeBundle struct{ mode string }

func (outcomeBundle) FilterTxs(map[common.Address]types.Transactions) {}
func (m outcomeBundle) ExtractTxBundles(txs []*types.Transaction, _ []*builder.Bundle) []*builder.Bundle {
	if m.mode == "target" {
		return []*builder.Bundle{builder.NewBundle(builder.NewTxOrGenList(txs[1]), txs[0].Hash(), true)}
	}
	if m.mode == "generation" {
		entry := builder.NewTxOrGenFromGen(func(uint64) (*types.Transaction, error) { return nil, errors.New("cannot generate") }, common.Hash{99})
		return []*builder.Bundle{builder.NewBundle(builder.NewTxOrGenList(txs[0], entry), common.Hash{}, false)}
	}
	return []*builder.Bundle{builder.NewBundle(builder.NewTxOrGenList(txs[0], txs[1]), common.Hash{}, false)}
}

func TestCandidateOutcomeMetrics(t *testing.T) {
	for _, tc := range []struct {
		name, reason, module      string
		failedReceipt, size, blob bool
		executionErr              error
	}{
		{name: "target", reason: "bundle_target_rejected", module: "target", failedReceipt: true},
		{name: "generation_before_execution", reason: "generator_failed", module: "initial_generation"},
		{name: "generation_in_bundle", reason: "generator_failed", module: "generation"},
		{name: "size_tx", reason: "size_limit", size: true},
		{name: "size_bundle", reason: "size_limit", module: "bundle", size: true},
		{name: "missing_blob_sidecar", reason: "blob_space_rejected", blob: true},
		{name: "blob_capacity", reason: "blob_space_rejected", blob: true},
		{name: "unsupported", reason: "unsupported_type", executionErr: blockchain.ErrTxTypeNotSupported},
		{name: "reverted_bundle", reason: "bundle_reverted", module: "bundle", failedReceipt: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer poolobserve.EnableForTesting()()
			reg := prometheus.NewRegistry()
			poolobserve.Register(reg)
			signer := types.LatestSignerForChainID(big.NewInt(1))
			key, err := crypto.GenerateKey()
			if err != nil {
				t.Fatal(err)
			}
			input := make(types.Transactions, 2)
			for i := range input {
				raw := types.NewTransaction(uint64(i), common.Address{}, big.NewInt(1), 21000, big.NewInt(25000000000), nil)
				if tc.blob {
					raw = types.NewTx(&types.TxInternalDataEthereumBlob{ChainID: uint256.NewInt(1), AccountNonce: uint64(i), Amount: uint256.NewInt(1), GasFeeCap: uint256.NewInt(25000000000), GasTipCap: uint256.NewInt(1), BlobFeeCap: uint256.NewInt(1), GasLimit: 21000, BlobHashes: []common.Hash{{1}}})
				}
				if tc.name == "blob_capacity" {
					raw = raw.WithBlobTxSidecar(types.NewBlobTxSidecar(0, []kzg4844.Blob{{}}, []kzg4844.Commitment{{}}, []kzg4844.Proof{{}}))
				}
				input[i], err = types.SignTx(raw, signer, key)
				if err != nil {
					t.Fatal(err)
				}
			}
			db, err := state.New(common.Hash{}, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			config := &params.ChainConfig{ChainID: big.NewInt(1)}
			if tc.blob {
				config.OsakaCompatibleBlock = big.NewInt(0)
				config.BlobScheduleConfig = &params.BlobScheduleConfig{Osaka: params.DefaultOsakaBlobConfig}
			}
			task := &Task{config: config, signer: signer, state: db, header: &types.Header{Number: big.NewInt(1), Time: big.NewInt(1)}}
			if tc.name == "blob_capacity" {
				task.blobs = config.LatestBlobConfigMax(task.header.Number)
			}
			if tc.size {
				task.size = params.MaxBlockSize - maxBlockSizeBufferZone
			}
			var modules []builder.TxBundlingModule
			if tc.module == "initial_generation" {
				modules = []builder.TxBundlingModule{failingGeneratedBundle{}}
			} else if tc.module != "" {
				modules = []builder.TxBundlingModule{outcomeBundle{tc.module}}
			}
			txs := types.NewTransactionsByPriceAndNonce(signer, map[common.Address]types.Transactions{crypto.PubkeyToAddress(key.PublicKey): input}, nil)
			bc := &candidateMetricsChain{firstError: tc.executionErr, firstFailedReceipt: tc.failedReceipt}
			task.ApplyTransactions(txs, bc, common.Address{}, modules)
			for _, reason := range []string{"bundle_target_rejected", "generator_failed", "size_limit", "blob_space_rejected", "unsupported_type", "bundle_reverted"} {
				want := float64(0)
				if reason == tc.reason {
					want = 1
				}
				if got := candidateCounter(t, reg, "kaia_poolobs_candidate_outcomes_total", "reason", reason); got != want {
					t.Fatalf("%s got %v want %v", reason, got, want)
				}
			}
		})
	}
}
