//ff:func feature=repair type=command control=iteration dimension=1
//ff:what applyInOrder — takes the workers' outcomes as they come and applies them in batch order, writing a status line now and then; stops all on the first error
package repair

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/geul-org/wikipgeul/internal/load"
)

// applyInOrder reads outcomes until the channel closes, keeps each batch
// until the ones before it are applied (drainPending) and writes statusLine
// to c.Log every c.Status, after the local time. The first failed batch or apply cancels the run
// (stop) and is returned; later outcomes are only drained.
func applyInOrder(db *sql.DB, c Config, results <-chan outcome, tokens <-chan struct{}, stop func(), st *Stats) error {
	pending := map[int64]Found{}
	refs := load.NewRefSet(c.RefLimit)
	var next int64
	var first error
	start, tick := time.Now(), time.Now()
	for out := range results {
		if out.Err != nil && first == nil {
			first = fmt.Errorf("batch %d (line %d): %w", out.Found.Seq, out.Found.Last, out.Err)
			stop()
		}
		if first != nil {
			continue
		}
		pending[out.Found.Seq] = out.Found
		if first = drainPending(db, c, pending, &next, tokens, refs, st); first != nil {
			stop()
		}
		if time.Since(tick) >= c.Status {
			tick = time.Now()
			fmt.Fprintln(c.Log, time.Now().Format("2006-01-02 15:04:05"), statusLine(*st, time.Since(start)))
		}
	}
	st.Elapsed = time.Since(start)
	return first
}
