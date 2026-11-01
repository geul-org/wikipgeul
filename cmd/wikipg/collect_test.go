//ff:func feature=cli type=command control=sequence
//ff:what TestCollect — in-order results save the checkpoint and finish the dump; a limit or a stop does not finish; a failed batch, a read error and a failed save are returned
package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestCollect(t *testing.T) {
	db, schema := testSchema(t)
	d := load.DumpInfo{ID: 1, FileName: "f", URL: "u", SHA256: "s", Loader: "l", LoaderCommit: "c"}
	load.InsertDump(db, schema, d)
	c := loadConfig{Schema: schema, Dump: d, Status: time.Millisecond}
	run := func(c loadConfig, ctx context.Context, rs []result, readErr error) error {
		ch, read := make(chan result, len(rs)), make(chan error, 1)
		for _, r := range rs {
			ch <- r
		}
		close(ch)
		read <- readErr
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		time.Sleep(2 * time.Millisecond)
		return collect(ctx, cancel, db, c, progress{Start: time.Now()}, ch, read)
	}
	ok := []result{{Seq: 1, Last: 4, Entities: 2}, {Seq: 0, Last: 2, Entities: 1, Errors: 1}}
	if err := run(c, context.Background(), ok, nil); err != nil {
		t.Fatal(err)
	}
	var line, ents int64
	var done bool
	db.QueryRow("SELECT checkpoint_line, entity_count, load_finished IS NOT NULL FROM "+schema+".dump").Scan(&line, &ents, &done)
	if line != 4 || ents != 3 || !done {
		t.Fatalf("checkpoint %d entities %d finished %v", line, ents, done)
	}
	limited := c
	limited.Limit = 5
	if err := run(limited, context.Background(), nil, nil); err != nil {
		t.Fatalf("limit: %v", err)
	}
	stopped, stop := context.WithCancel(context.Background())
	stop()
	if err := run(c, stopped, nil, nil); err != nil {
		t.Fatalf("stopped: %v", err)
	}
	if err := run(c, context.Background(), []result{{Seq: 0, Err: errors.New("copy")}, {Seq: 1, Err: errors.New("again")}}, nil); err == nil {
		t.Fatal("failed batch: want error")
	}
	if err := run(c, context.Background(), nil, errors.New("disk")); err == nil {
		t.Fatal("read error: want error")
	}
	gone := c
	gone.Schema = schema + "_none"
	if err := run(gone, context.Background(), []result{{Seq: 0, Last: 2}}, nil); err == nil {
		t.Fatal("failed save: want error")
	}
}
