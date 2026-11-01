//ff:func feature=cli type=command control=sequence
//ff:what TestWriter — a parsed batch is committed and reported; a failing write is reported with its error
package main

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestWriter(t *testing.T) {
	db, schema := testSchema(t)
	r := load.NewRows()
	r.AddEntity(load.Config{DumpID: 1}, 2, nil, sampleOne(), load.NewRefSet(64))
	r.Error(1, 3, "", "parse", "x", nil)
	in, out := make(chan parsed, 2), make(chan result, 2)
	in <- parsed{Seq: 0, Last: 3, Rows: r, Entities: 1}
	in <- parsed{Seq: 1, Last: 4, Rows: r, Entities: 1}
	close(in)
	db.Exec("INSERT INTO " + schema + ".dump (dump_id, dump_date, file_name, url, size_bytes, sha256, downloaded_at, loader, loader_commit) VALUES (1, '2026-10-19', 'f', 'u', 1, 's', now(), 'l', 'c')")
	db.Exec("ALTER TABLE " + schema + ".entity ADD PRIMARY KEY (entity_id, dump_id)")
	writer(db, schema, in, out)
	a, b := <-out, <-out
	if a.Err != nil || a.Entities != 1 || a.Errors != 1 || a.Last != 3 {
		t.Fatalf("first %+v", a)
	}
	if b.Err == nil {
		t.Fatal("duplicate batch: want error")
	}
}
