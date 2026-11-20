package blockchain

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/internal/poolobserve"
	kaiamock "github.com/kaiachain/kaia/kaiax/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestPromotionStopDistinguishesIdenticalPoolSnapshots(t *testing.T) {
	for _, name := range []string{"nonce_gap", "module_hold"} {
		t.Run(name, func(t *testing.T) {
			defer poolobserve.EnableForTesting()()
			cfg := *params.TestChainConfig
			cfg.MagmaCompatibleBlock = big.NewInt(0)
			pool, key := setupTxPoolWithConfig(&cfg)
			pool.Stop()
			r := prometheus.NewRegistry()
			poolobserve.Register(r)
			mod := kaiamock.NewMockTxPoolModule(gomock.NewController(t))
			mod.EXPECT().IsModuleTx(gomock.Any()).DoAndReturn(func(tx *types.Transaction) bool { return name == "module_hold" && tx.Nonce() == 1 }).AnyTimes()
			if name == "module_hold" {
				mod.EXPECT().GetCheckBalance().Return(nil).AnyTimes()
				mod.EXPECT().PreAddTx(gomock.Any(), true).Return(nil)
				mod.EXPECT().IsReady(gomock.Any(), uint64(1), gomock.Any()).Return(false).MinTimes(1)
			}
			pool.RegisterTxPoolModule(mod)
			nonce := uint64(2)
			if name == "module_hold" {
				nonce = 1
			}
			for _, n := range []uint64{0, nonce} {
				if err := pool.AddLocal(pricedTransaction(n, 21000, new(big.Int).Set(pool.gasPrice), key)); err != nil {
					t.Fatal(err)
				}
			}
			p, q := pool.Stats()
			snapshot, err := pool.Pending()
			if err != nil {
				t.Fatal(err)
			}
			size := 0
			for _, txs := range snapshot {
				size += len(txs)
			}
			a := propagationEvent(t, r, "admitted")
			m := propagationEvent(t, r, "promoted")
			if p != 1 || q != 1 || size != 1 || a != 2 || m != 1 {
				t.Fatalf("unexpected %d %d %d %g %g", p, q, size, a, m)
			}
			counts := promotionStops(t, r)
			if name == "nonce_gap" {
				require.Positive(t, counts["nonce_gap"])
				require.Zero(t, counts["module_not_ready"])
			} else {
				require.Positive(t, counts["module_not_ready"])
				require.Zero(t, counts["nonce_gap"])
			}
			require.Zero(t, counts["below_base_fee"])
		})
	}
}

func promotionStops(t *testing.T, reg *prometheus.Registry) map[string]float64 {
	t.Helper()
	families, err := reg.Gather()
	require.NoError(t, err)
	for _, family := range families {
		if family.GetName() == "kaia_poolobs_promotion_stops_total" {
			counts := make(map[string]float64)
			for _, sample := range family.Metric {
				counts[sample.Label[0].GetValue()] = sample.Counter.GetValue()
			}
			require.Len(t, counts, 3)
			return counts
		}
	}
	t.Fatal("promotion stop counters are missing")
	return nil
}

func TestPromotionStopAfterReadyPrefix(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	for _, tc := range []struct {
		name   string
		nonces []uint64
		prices []int64
		ready  int
		reason string
	}{
		{"empty", nil, nil, 0, ""},
		{"drained", []uint64{0, 1}, []int64{2, 2}, 2, ""},
		{"head gap", []uint64{1}, []int64{2}, 0, "nonce_gap"},
		{"gap after prefix", []uint64{0, 2}, []int64{2, 2}, 1, "nonce_gap"},
		{"price after prefix", []uint64{0, 1}, []int64{2, 1}, 1, "below_base_fee"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(poolobserve.EnableForTesting())
			reg := prometheus.NewRegistry()
			poolobserve.Register(reg)
			m := newTxSortedMap()
			for i, n := range tc.nonces {
				m.Put(pricedTransaction(n, 21000, big.NewInt(tc.prices[i]), key))
			}
			require.Len(t, m.ReadyWithGasPrice(0, big.NewInt(2), nil), tc.ready)
			require.Equal(t, len(tc.nonces)-tc.ready, m.Len())
			want := map[string]float64{"nonce_gap": 0, "below_base_fee": 0, "module_not_ready": 0}
			if tc.reason != "" {
				want[tc.reason] = 1
			}
			require.Equal(t, want, promotionStops(t, reg))
		})
	}
}
