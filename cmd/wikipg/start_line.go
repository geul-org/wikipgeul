//ff:func feature=cli type=command control=sequence
//ff:what startLine — where a load starts: the given line, or after the checkpoint with later partial rows cleared
package main

import (
	"database/sql"

	"github.com/geul-org/wikipgeul/internal/load"
)

// startLine decides the first line to read and the counts to continue from.
// With -from-line it is that line and counting starts at zero. Otherwise it is
// the line after the dump's checkpoint (1 for a new dump): rows of later lines
// a stopped run committed in part are cleared first (load.ClearAfter), and
// the saved counts carry on.
func startLine(db *sql.DB, c loadConfig) (from, entities, errs int64, err error) {
	if c.From > 0 {
		return c.From, 0, 0, nil
	}
	line, entities, errs, err := load.ReadCheckpoint(db, c.Schema, c.Dump.ID)
	if err != nil {
		return 0, 0, 0, err
	}
	if line > 0 {
		if err := load.ClearAfter(db, c.Schema, c.Dump.ID, line); err != nil {
			return 0, 0, 0, err
		}
	}
	return line + 1, entities, errs, nil
}
