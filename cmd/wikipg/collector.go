//ff:type feature=cli type=model
//ff:what collector — the collector's state: database, config, checkpoint, totals, first failure, how to stop reading
package main

import (
	"context"
	"database/sql"

	"github.com/geul-org/wikipgeul/internal/load"
)

// collector holds what collect keeps between results: the database and
// config, the checkpoint, the running totals (Base: the entities before this
// run), the first failure and the function that stops reading.
type collector struct {
	DB     *sql.DB
	C      loadConfig
	CP     *load.Checkpoint
	P      progress
	Base   int64
	Failed error
	Cancel context.CancelFunc
}
