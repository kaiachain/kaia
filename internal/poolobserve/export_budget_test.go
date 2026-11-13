// Copyright 2026 The Kaia Authors
// This file is part of the kaia library.
package poolobserve

import (
	"bytes"
	"os"
	"testing"

	"github.com/prometheus/common/expfmt"
)

func TestExportBudget(t *testing.T) {
	r := newRecorder()
	fs, e := r.registry.Gather()
	if e != nil {
		t.Fatal(e)
	}
	var b bytes.Buffer
	for _, f := range fs {
		if _, e := expfmt.MetricFamilyToText(&b, f); e != nil {
			t.Fatal(e)
		}
	}
	if b.Len() > 60000 {
		t.Fatalf("fixed scrape budget exceeded: %d bytes", b.Len())
	}
	t.Logf("schema5 idle scrape bytes=%d; includes HELP and TYPE; not compressed network or storage", b.Len())
	if p := os.Getenv("POOLOBS_EXPORT_TEST_PATH"); p != "" {
		if e := os.WriteFile(p, b.Bytes(), 0o600); e != nil {
			t.Fatal(e)
		}
	}
}
