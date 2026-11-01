//ff:func feature=cli type=command control=iteration dimension=1
//ff:what writer — commits each parsed batch in one transaction and reports it
package main

import (
	"database/sql"

	"github.com/geul-org/wikipgeul/internal/load"
)

// writer writes every parsed batch from in in one transaction
// (load.WriteRows) and reports a result for each to out, its error included.
// Writers wait on the database, parsers on the CPU; keeping them apart lets
// both run at once.
func writer(db *sql.DB, schema string, in <-chan parsed, out chan<- result) {
	for p := range in {
		err := load.WriteRows(db, schema, p.Rows)
		out <- result{Seq: p.Seq, Last: p.Last, Entities: p.Entities, Errors: int64(p.Rows.Errors), Err: err}
	}
}
