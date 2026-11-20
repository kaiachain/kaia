package cn

import (
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/internal/poolobserve"
	"github.com/prometheus/client_golang/prometheus"
)

func TestPoolObservabilityBroadcastQueueLifecycle(t *testing.T) {
	defer poolobserve.EnableForTesting()()
	ch := make(chan blockchain.NewTxsEvent, txChanSize)
	feed := new(event.Feed)
	sub := feed.Subscribe(ch)
	defer sub.Unsubscribe()
	pm := &ProtocolManager{txsCh: ch, txsSub: sub}
	registry := prometheus.NewRegistry()
	poolobserve.Register(registry)
	done := make(chan struct{})
	go func() { pm.txBroadcastLoop(); close(done) }()
	value := func(family string) float64 {
		fs, err := registry.Gather()
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range fs {
			if f.GetName() == family {
				for _, m := range f.Metric {
					if m.Label[0].GetValue() == "broadcast" {
						return m.Gauge.GetValue()
					}
				}
			}
		}
		t.Fatal("metric missing")
		return 0
	}
	deadline := time.Now().Add(3 * time.Second)
	for value("kaia_poolobs_queue_observed") != 1 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if value("kaia_poolobs_queue_observed") != 1 || value("kaia_poolobs_queue_capacity") != float64(txChanSize) {
		t.Fatal("broadcast channel not observed")
	}
	sub.Unsubscribe()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("broadcast consumer did not stop")
	}
	if value("kaia_poolobs_queue_observed") != 0 {
		t.Fatal("stopped consumer reported healthy queue")
	}
}
