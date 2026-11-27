package blockchain

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/prometheus/client_golang/prometheus"
)

func TestBlockImportOperationMetrics(t *testing.T) {
	db, bc, err := newCanonical(faker.NewFaker(), 0, true)
	if err != nil {
		t.Fatal(err)
	}
	defer bc.Stop()
	block := MakeBlockChain(bc.CurrentBlock(), 1, faker.NewFaker(), db, 10)[0]
	defer poolobserve.EnableForTesting()()
	reg := prometheus.NewRegistry()
	poolobserve.Register(reg)
	// Empty, rejected, and successful calls all finish the operation. This metric
	// counts InsertChain calls, not imported blocks or pure execution time.
	if _, err = bc.InsertChain(nil); err != nil {
		t.Fatal(err)
	}
	forged := block.WithBody(types.Transactions{types.NewTransaction(0, common.Address{}, big.NewInt(0), 21000, big.NewInt(1), nil)})
	if _, err = bc.InsertChain(types.Blocks{forged}); err == nil {
		t.Fatal("expected body mismatch")
	}
	if _, err = bc.InsertChain(types.Blocks{block}); err != nil {
		t.Fatal(err)
	}
	fs, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, f := range fs {
		for _, m := range f.Metric {
			if len(m.Label) != 1 || m.Label[0].GetValue() != "block_import" {
				continue
			}
			switch f.GetName() {
			case "kaia_poolobs_operations_started_total":
				checked++
				if m.Counter.GetValue() != 3 {
					t.Fatalf("%s=%v", f.GetName(), m.Counter.GetValue())
				}
			case "kaia_poolobs_operations_inflight":
				checked++
				if m.Gauge.GetValue() != 0 {
					t.Fatal("inflight leaked")
				}
			case "kaia_poolobs_duration_seconds":
				checked++
				if m.Histogram.GetSampleCount() != 3 {
					t.Fatal("duration missing")
				}
			}
		}
	}
	if checked != 3 {
		t.Fatalf("checked %d metric families", checked)
	}
}
