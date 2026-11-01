//ff:func feature=cli type=command control=sequence
//ff:what TestExecRepair — runs a parsed repair over a loaded schema and reports; a DSN lib/pq cannot parse and a missing schema fail
package main

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/repair"
)

func TestExecRepair(t *testing.T) {
	_, schema := testSchema(t)
	f, err := parseRepairFlags("wikipgeul build-reference-order", []string{"-schema", schema, "-workers", "1", "-dry-run"})
	if err != nil {
		t.Fatal(err)
	}
	f.Config.Order = true
	if err := execRepair(f); err != nil {
		t.Fatal(err)
	}
	f.Config.DryRun = false
	if err := execRepair(f); err != nil {
		t.Fatal("empty write:", err)
	}
	bad := f
	bad.DSN = "host='x"
	if err := execRepair(bad); err == nil {
		t.Fatal("unparsable DSN: want error")
	}
	bad = repairFlags{Config: repair.Config{Schema: schema + "_none", DumpID: 1, Batch: 1, Workers: 1, Status: 1, Order: true}}
	if err := execRepair(bad); err == nil {
		t.Fatal("missing schema: want error")
	}
}
