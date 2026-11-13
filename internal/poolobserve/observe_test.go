package poolobserve

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestRecorderConcurrentCompletion(t *testing.T) {
	defer EnableForTesting()()
	r := current.Load()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); end := Begin(AdmissionWriteSingle); end(); end() }()
	}
	wg.Wait()
	fs, err := r.registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range fs {
		for _, m := range f.Metric {
			if len(m.Label) == 0 || m.Label[0].GetValue() != names[AdmissionWriteSingle] {
				continue
			}
			switch f.GetName() {
			case "kaia_poolobs_operations_started_total":
				if m.Counter.GetValue() != 100 {
					t.Fatal(f.GetName(), m.Counter.GetValue())
				}
			case "kaia_poolobs_operations_inflight":
				if m.Gauge.GetValue() != 0 {
					t.Fatal("inflight leak")
				}
			case "kaia_poolobs_duration_seconds":
				if m.Histogram.GetSampleCount() != 100 {
					t.Fatal("duplicate histogram completion")
				}
			}
		}
	}
}

func TestRecorderDisabledAndReplacement(t *testing.T) {
	previous := current.Swap(nil)
	defer current.Store(previous)
	reg := prometheus.NewRegistry()
	Register(reg)
	Begin(AdmissionWriteSingle)()
	fs, err := reg.Gather()
	if err != nil || len(fs) != 0 {
		t.Fatal("disabled collector exported metrics", err)
	}
	restore := EnableForTesting()
	defer restore()
	end := Begin(ResetHold)
	old := current.Load()
	next := newRecorder()
	current.Store(next)
	end()
	end()
	for _, r := range []*recorder{old, next} {
		fs, err := r.registry.Gather()
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range fs {
			if f.GetName() == "kaia_poolobs_operations_inflight" {
				for _, m := range f.Metric {
					if m.Gauge.GetValue() != 0 {
						t.Fatal("late completion corrupted recorder")
					}
				}
			}
		}
	}
	// No active-operation map exists: retaining arbitrarily many operations never
	// allocates recorder slots or requires an overflow policy.
}

func TestRecorderSchema(t *testing.T) {
	r := newRecorder()
	fs, err := r.registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	series := 0
	for _, f := range fs {
		if f.GetHelp() == "" {
			t.Fatalf("missing HELP %s", f.GetName())
		}
		for _, m := range f.Metric {
			if m.Histogram != nil {
				series += len(m.Histogram.Bucket) + 3
				want := []float64{.1, .15, .25, 1, 5, 10, 20}
				if len(m.Histogram.Bucket) != len(want) {
					t.Fatalf("%s bucket count", f.GetName())
				}
				for i, bucket := range m.Histogram.Bucket {
					if bucket.GetUpperBound() != want[i] {
						t.Fatalf("%s bucket %d: got %g want %g", f.GetName(), i, bucket.GetUpperBound(), want[i])
					}
				}
			} else {
				series++
			}
			for _, l := range m.Label {
				if l.GetName() != "stage" && l.GetName() != "queue" && l.GetName() != "event" && l.GetName() != "phase" && l.GetName() != "reason" {
					t.Fatal("unbounded label")
				}
				if strings.Contains(l.GetValue(), "capacity_read") {
					t.Fatal("removed lock exported as zero")
				}
			}
			switch f.GetName() {
			case "kaia_poolobs_schema_version":
				if m.Gauge.GetValue() != 5 {
					t.Fatal("wrong schema")
				}
			}
		}
	}
	if series != 268 {
		t.Fatalf("series=%d want268", series)
	}
}

func BenchmarkObservation(b *testing.B) {
	previous := current.Load()
	defer current.Store(previous)
	for _, enabled := range []bool{false, true} {
		name := "disabled"
		current.Store(nil)
		if enabled {
			name = "enabled"
			current.Store(newRecorder())
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				Begin(AdmissionWriteSingle)()
			}
		})
	}
}

func queueMetric(t *testing.T, r *recorder, family, queue string) float64 {
	t.Helper()
	fs, e := r.registry.Gather()
	if e != nil {
		t.Fatal(e)
	}
	for _, f := range fs {
		if f.GetName() == family {
			for _, m := range f.Metric {
				if m.Label[0].GetValue() == queue {
					return m.Gauge.GetValue()
				}
			}
		}
	}
	t.Fatal("queue metric missing")
	return 0
}

func TestQueueLifecycle(t *testing.T) {
	defer EnableForTesting()()
	r := current.Load()
	if queueMetric(t, r, "kaia_poolobs_queue_observed", "feed") != 0 {
		t.Fatal("unbound queue looked observed")
	}
	ch := make(chan int, 2)
	detach := TrackQueue(FeedQueue, func() int { return len(ch) }, cap(ch))
	ch <- 1
	if queueMetric(t, r, "kaia_poolobs_queue_slots", "feed") != 1 || queueMetric(t, r, "kaia_poolobs_queue_capacity", "feed") != 2 {
		t.Fatal("queue units differ")
	}
	newer := make(chan int, 3)
	detachNew := TrackQueue(FeedQueue, func() int { return len(newer) }, cap(newer))
	detach()
	if queueMetric(t, r, "kaia_poolobs_queue_observed", "feed") != 1 || queueMetric(t, r, "kaia_poolobs_queue_capacity", "feed") != 3 {
		t.Fatal("old cleanup detached replacement")
	}
	detachNew()
	if queueMetric(t, r, "kaia_poolobs_queue_observed", "feed") != 0 {
		t.Fatal("stopped queue still observed")
	}
}

// Each child starts a new OS process so identity and disabled export are tested
// across real process lifetimes, not just swaps of the in-memory recorder.
func TestProcessIdentity(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	run := func(enabled string) map[string]float64 {
		cmd := exec.CommandContext(t.Context(), executable, "-test.run=^TestProcessIdentityChild$")
		for _, e := range os.Environ() {
			if !strings.HasPrefix(e, "KAIA_POOLOBSERVE=") && !strings.HasPrefix(e, "POOLOBS_CHILD=") {
				cmd.Env = append(cmd.Env, e)
			}
		}
		cmd.Env = append(cmd.Env, "KAIA_POOLOBSERVE="+enabled, "POOLOBS_CHILD=1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("child failed: %v %s", err, output)
		}
		for _, line := range strings.Split(string(output), "\n") {
			if strings.HasPrefix(line, "IDENTITY ") {
				m := map[string]float64{}
				if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "IDENTITY ")), &m); err != nil {
					t.Fatal(err)
				}
				return m
			}
		}
		t.Fatalf("child identity missing: %s", output)
		return nil
	}
	a, b := run("1"), run("1")
	if a["start"] == b["start"] {
		t.Fatal("restarted process reused observation start timestamp")
	}
	if a["start"] <= 0 || b["start"] <= 0 {
		t.Fatal("invalid process metadata")
	}
	if off := run("0"); len(off) != 0 {
		t.Fatal("disabled process exported collector")
	}
}

func TestProcessIdentityChild(t *testing.T) {
	if os.Getenv("POOLOBS_CHILD") != "1" {
		t.Skip("subprocess helper")
	}
	registry := prometheus.NewRegistry()
	Register(registry)
	fs, err := registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]float64{}
	for _, f := range fs {
		for _, x := range f.Metric {
			switch f.GetName() {
			case "kaia_poolobs_process_start_time_seconds":
				m["start"] = x.Gauge.GetValue()
			}
		}
	}
	data, _ := json.Marshal(m)
	fmt.Printf("IDENTITY %s\n", data)
}

// A gauge can be zero at both scrapes even though an operation happened between
// them. The native cumulative histogram retains the completed observation.
func TestCompletedOperationBetweenScrapes(t *testing.T) {
	defer EnableForTesting()()
	r := current.Load()
	before, err := r.registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	if len(before) == 0 {
		t.Fatal("no initial sample")
	}
	done := Begin(FeedEnqueue)
	done()
	after, err := r.registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	seen := false
	for _, family := range after {
		if family.GetName() != "kaia_poolobs_duration_seconds" {
			continue
		}
		for _, m := range family.Metric {
			for _, label := range m.Label {
				if label.GetName() == "stage" && label.GetValue() == "feed_enqueue" {
					if m.Histogram.GetSampleCount() != 1 {
						t.Fatal("completed sub-scrape operation lost")
					}
					seen = true
				}
			}
		}
	}
	if !seen || metricValue(r.stages[FeedEnqueue].inflight) != 0 {
		t.Fatal("wrong cumulative or live state")
	}
}
