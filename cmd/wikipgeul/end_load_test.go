//ff:func feature=cli type=command control=sequence
//ff:what TestEndLoad — failure first, then the read error with the resume line, then limited or stopped runs end quietly, else the dump is finished
package main

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestEndLoad(t *testing.T) {
	db, schema := testSchema(t)
	d := load.DumpInfo{ID: 1, FileName: "f", URL: "u", SHA256: "s", Loader: "l", LoaderCommit: "c"}
	load.InsertDump(db, schema, d)
	c := loadConfig{Schema: schema, Dump: d}
	bg := context.Background()
	fail := errors.New("batch")
	if err := endLoad(bg, db, c, progress{}, fail, errors.New("read")); err != fail {
		t.Fatal("failure first")
	}
	if err := endLoad(bg, db, c, progress{Line: 7}, nil, errors.New("read")); err == nil || !strings.Contains(err.Error(), "after line 7") {
		t.Fatalf("read error: %v", err)
	}
	limited := c
	limited.Limit = 3
	stopped, stop := context.WithCancel(bg)
	stop()
	if endLoad(bg, db, limited, progress{}, nil, nil) != nil || endLoad(stopped, db, c, progress{}, nil, nil) != nil {
		t.Fatal("early end: want nil")
	}
	if err := endLoad(bg, db, c, progress{}, nil, nil); err != nil {
		t.Fatal(err)
	}
	var done bool
	db.QueryRow("SELECT load_finished IS NOT NULL FROM " + schema + ".dump").Scan(&done)
	if !done {
		t.Fatal("not finished")
	}
}
