//ff:func feature=cli type=command control=sequence
//ff:what TestCollectorStep — a result is applied, a tick prints and goes on, a closed channel ends
package main

import (
	"testing"
	"time"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestCollectorStep(t *testing.T) {
	db, schema := testSchema(t)
	d := load.DumpInfo{ID: 1, FileName: "f", URL: "u", SHA256: "s", Loader: "l", LoaderCommit: "c"}
	load.InsertDump(db, schema, d)
	k := &collector{DB: db, C: loadConfig{Schema: schema, Dump: d}, CP: load.NewCheckpoint(1), Cancel: func() {}}
	results, tick := make(chan result, 1), make(chan time.Time, 1)
	results <- result{Seq: 0, Last: 4, Entities: 2}
	if !k.step(results, tick) || k.P.Entities != 2 || k.P.Line != 4 {
		t.Fatalf("result: %+v", k.P)
	}
	tick <- time.Now()
	if !k.step(results, tick) {
		t.Fatal("tick: want true")
	}
	close(results)
	if k.step(results, tick) {
		t.Fatal("closed: want false")
	}
}
