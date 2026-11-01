//ff:func feature=convert type=loader control=sequence
//ff:what closedDB — test fixture: a database handle that is already closed, so every query fails
package pgsource

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
)

// closedDB returns a lib/pq handle closed before use: every query on it
// fails with sql.ErrConnDone-like errors, which exercises the error paths.
func closedDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", "")
	if err != nil {
		t.Fatal(err)
	}
	db.Close()
	return db
}
