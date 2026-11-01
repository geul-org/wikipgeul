//ff:func feature=cli type=command control=selection
//ff:what applyResult — adds a batch result to the totals and saves the checkpoint when it advances; a failure cancels reading
package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/geul-org/wikipgeul/internal/load"
)

// applyResult applies result r and returns the load's failure so far. A
// failed batch cancels reading and, if it is the first failure, becomes it.
// Otherwise its counts are added and, while nothing has failed, the
// checkpoint advances over batches done in order and is saved
// (load.SaveCheckpoint); a failed save is a failure too.
func applyResult(db *sql.DB, c loadConfig, cp *load.Checkpoint, p *progress, r result, failed error, cancel context.CancelFunc) error {
	switch {
	case r.Err != nil:
		cancel()
		if failed == nil {
			return fmt.Errorf("batch %d (to line %d): %w", r.Seq, r.Last, r.Err)
		}
		return failed
	case failed != nil:
		return failed
	}
	p.Entities, p.Errors = p.Entities+r.Entities, p.Errors+r.Errors
	if !cp.Done(r.Seq, r.Last) {
		return nil
	}
	p.Line = cp.Line
	if err := load.SaveCheckpoint(db, c.Schema, c.Dump.ID, p.Line, p.Entities, p.Errors); err != nil {
		cancel()
		return err
	}
	return nil
}
