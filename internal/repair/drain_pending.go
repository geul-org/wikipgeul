//ff:func feature=repair type=writer control=iteration dimension=1
//ff:what drainPending — applies the waiting batches that are next in order, releasing a batch slot for each
package repair

import (
	"database/sql"
	"fmt"

	"github.com/geul-org/wikipgeul/internal/load"
)

// drainPending applies (applyFound) and removes the batches of pending from
// *next on while they are there, so batches apply in line order however the
// workers finish; each applied batch frees one slot of tokens for the feeder.
func drainPending(db *sql.DB, c Config, pending map[int64]Found, next *int64, tokens <-chan struct{}, refs *load.RefSet, st *Stats) error {
	for f, ok := pending[*next]; ok; f, ok = pending[*next] {
		if err := applyFound(db, c, f, refs, st); err != nil {
			return fmt.Errorf("batch %d (line %d): %w", f.Seq, f.Last, err)
		}
		delete(pending, *next)
		*next++
		<-tokens
	}
	return nil
}
