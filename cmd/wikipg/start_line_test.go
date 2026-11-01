//ff:func feature=cli type=command control=sequence
//ff:what TestStartLine — -from-line wins; a new dump starts at line 1; a checkpoint resumes after it with later rows cleared and counts kept; an unknown dump fails
package main

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestStartLine(t *testing.T) {
	db, schema := testSchema(t)
	c := loadConfig{Schema: schema, Dump: load.DumpInfo{ID: 1, FileName: "f", URL: "u", SHA256: "s", Loader: "l", LoaderCommit: "c"}}
	if from, _, _, err := startLine(db, loadConfig{From: 9}); err != nil || from != 9 {
		t.Fatalf("from-line: %d %v", from, err)
	}
	if _, _, _, err := startLine(db, c); err == nil {
		t.Fatal("no dump row: want error")
	}
	load.InsertDump(db, schema, c.Dump)
	if from, e, x, err := startLine(db, c); err != nil || from != 1 || e != 0 || x != 0 {
		t.Fatalf("new: %d %d %d %v", from, e, x, err)
	}
	r := load.NewRows()
	r.AddEntity(load.Config{DumpID: 1}, 5, nil, sampleOne(), load.NewRefSet(64))
	load.WriteRows(db, schema, r)
	load.SaveCheckpoint(db, schema, 1, 4, 3, 1)
	from, e, x, err := startLine(db, c)
	var left int
	db.QueryRow("SELECT count(*) FROM " + schema + ".entity").Scan(&left)
	if err != nil || from != 5 || e != 3 || x != 1 || left != 0 {
		t.Fatalf("resume: %d %d %d %v left %d", from, e, x, err, left)
	}
	db.Exec("DROP TABLE " + schema + ".label CASCADE")
	if _, _, _, err := startLine(db, c); err == nil {
		t.Fatal("clear failure: want error")
	}
}
