//ff:func feature=cli type=command control=sequence
//ff:what execRepair — opens the database for a repair or order build run, runs it until done or interrupted and writes the report
package main

import (
	"context"
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/lib/pq"

	"github.com/geul-org/wikipgeul/internal/repair"
)

// execRepair opens the database of f (sessions with jit off, as the
// per-line and per-hash lookups are cheap but estimated large; sequential
// scans off, as every query of the run goes by an index and a scan of these
// tables takes hours; no parallel workers, so the queries at once are at
// most the workers; read-only for -dry-run: dsnWith; a DSN lib/pq cannot
// parse fails here) with a connection per worker and two more, runs
// repair.Run until done or SIGINT/SIGTERM (batches already read are still
// applied, so the checkpoint stays exact) and writes the report to stdout.
func execRepair(f repairFlags) error {
	dsn := dsnWith(dsnWith(dsnWith(f.DSN, "jit=off"), "enable_seqscan=off"), "max_parallel_workers_per_gather=0")
	if f.Config.DryRun {
		dsn = dsnWith(dsn, "default_transaction_read_only=on")
	}
	conn, err := pq.NewConnector(dsn)
	if err != nil {
		return err
	}
	db := sql.OpenDB(conn)
	defer db.Close()
	db.SetMaxOpenConns(f.Config.Workers + 2)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	st, err := repair.Run(ctx, db, f.Config)
	writeReport(os.Stdout, f.Config, st)
	return err
}
