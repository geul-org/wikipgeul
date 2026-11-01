//ff:func feature=cli type=command control=selection
//ff:what endLoad — the outcome of a load: its failure, the reader's error with the resume line, an early end, or the dump marked finished
package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/geul-org/wikipgeul/internal/load"
)

// endLoad returns how the load ended: the batch or checkpoint failure, the
// reader's error with the line to resume after, nothing for a limited or
// stopped run, or, when the whole input was read, the dump marked finished
// (load.FinishDump).
func endLoad(ctx context.Context, db *sql.DB, c loadConfig, p progress, failed, readErr error) error {
	switch {
	case failed != nil:
		return failed
	case readErr != nil:
		return fmt.Errorf("stopped after line %d (run the same command to resume): %w", p.Line, readErr)
	case c.Limit > 0 || ctx.Err() != nil:
		return nil
	default:
		return load.FinishDump(db, c.Schema, c.Dump.ID)
	}
}
