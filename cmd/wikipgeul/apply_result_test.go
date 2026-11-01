//ff:func feature=cli type=command control=sequence
//ff:what TestApplyResult — counts add and the checkpoint saves in order; a failure cancels and is kept first; after a failure results are ignored; a failed save fails
package main

import (
	"errors"
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestApplyResult(t *testing.T) {
	db, schema := testSchema(t)
	d := load.DumpInfo{ID: 1, FileName: "f", URL: "u", SHA256: "s", Loader: "l", LoaderCommit: "c"}
	load.InsertDump(db, schema, d)
	c := loadConfig{Schema: schema, Dump: d}
	cp, p := load.NewCheckpoint(1), &progress{Line: 1}
	canceled := 0
	cancel := func() { canceled++ }
	if err := applyResult(db, c, cp, p, result{Seq: 1, Last: 5, Entities: 2}, nil, cancel); err != nil || p.Line != 1 || p.Entities != 2 {
		t.Fatalf("out of order: %v %+v", err, p)
	}
	if err := applyResult(db, c, cp, p, result{Seq: 0, Last: 3, Entities: 1, Errors: 1}, nil, cancel); err != nil || p.Line != 5 || p.Entities != 3 || p.Errors != 1 {
		t.Fatalf("in order: %v %+v", err, p)
	}
	first := applyResult(db, c, cp, p, result{Seq: 2, Err: errors.New("copy")}, nil, cancel)
	if first == nil || canceled != 1 {
		t.Fatal("failure: want error and cancel")
	}
	if err := applyResult(db, c, cp, p, result{Seq: 3, Err: errors.New("again")}, first, cancel); err != first {
		t.Fatal("second failure must keep the first")
	}
	if err := applyResult(db, c, cp, p, result{Seq: 2, Last: 9, Entities: 5}, first, cancel); err != first || p.Entities != 3 {
		t.Fatal("after failure: ignored")
	}
	gone := c
	gone.Schema = schema + "_none"
	if err := applyResult(db, gone, load.NewCheckpoint(0), &progress{}, result{Seq: 0, Last: 2}, nil, cancel); err == nil {
		t.Fatal("failed save: want error")
	}
}
