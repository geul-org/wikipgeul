//ff:func feature=cli type=command control=sequence
//ff:what runLoad — loads a dump through a parallel pipeline: decompressor, reader, parsers, writers, collector; stops cleanly on interrupt
package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lib/pq"

	"github.com/geul-org/wikipgeul/internal/load"
)

// runLoad loads the dump c.In into c.Schema: it writes the dump row
// (load.InsertDump), finds the start line (startLine) and runs the pipeline:
// the decompressor process (openInput, c.Decoders threads), one reader
// splitting lines into batches (readBatches), c.Parsers parsers and c.Writers
// writers (startPipeline), and the collector recording checkpoints (collect). SIGINT or SIGTERM
// stops reading; batches already read are still written, so the checkpoint
// stays exact.
func runLoad(args []string) error {
	c, err := parseLoadFlags(args)
	if err != nil {
		return err
	}
	conn, err := pq.NewConnector(c.DSN)
	if err != nil {
		return err
	}
	db := sql.OpenDB(conn)
	defer db.Close()
	db.SetMaxOpenConns(c.Writers + 2)
	if err := load.InsertDump(db, c.Schema, c.Dump); err != nil {
		return err
	}
	from, entities, errs, err := startLine(db, c)
	if err != nil {
		return err
	}
	in, err := openInput(c.In, c.Decompressor, c.Decoders)
	if err != nil {
		return err
	}
	defer in.Close()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	lc := load.Config{DumpID: c.Dump.ID, Langs: langSet(c.Langs), Raw: c.Raw}
	batches := make(chan batch, c.Parsers)
	read := make(chan error, 1)
	go func() {
		_, err := readBatches(ctx, in, from, c.Limit, c.Batch, c.BatchBytes, batches)
		read <- err
	}()
	results := startPipeline(db, c, lc, load.NewRefSet(c.RefLimit), batches)
	p := progress{Line: from - 1, Entities: entities, Errors: errs, Start: time.Now()}
	return collect(ctx, cancel, db, c, p, results, read)
}
