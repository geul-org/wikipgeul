//ff:func feature=cli type=command control=iteration dimension=1
//ff:what collect — gathers batch results until the writers are done, printing status, then ends the load (endLoad)
package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/geul-org/wikipgeul/internal/load"
)

// collect reads results until every writer is done (collector.step: counts,
// checkpoint, failure, a status line every c.Status), prints a last status
// line and lets endLoad decide the outcome with the reader's error.
func collect(ctx context.Context, cancel context.CancelFunc, db *sql.DB, c loadConfig, p progress, results <-chan result, read <-chan error) error {
	k := &collector{DB: db, C: c, CP: load.NewCheckpoint(p.Line), P: p, Base: p.Entities, Cancel: cancel}
	tick := time.NewTicker(c.Status)
	defer tick.Stop()
	for k.step(results, tick.C) {
	}
	fmt.Fprintln(os.Stderr, k.P.Text(time.Now(), k.Base))
	return endLoad(ctx, db, c, k.P, k.Failed, <-read)
}
