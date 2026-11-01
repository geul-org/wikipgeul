//ff:func feature=convert type=command control=sequence
//ff:what openDB — one database connection with JIT off, as run uses
package exportgeul

import (
	"database/sql"
	"database/sql/driver"
)

// openDB opens a pool of one connection over conn with JIT off (run: the
// planner JIT-compiles the reference join, about 160 ms per entity against
// 0.5 ms without). Pack mode gives each worker its own.
func openDB(conn driver.Connector) *sql.DB {
	db := sql.OpenDB(conn)
	db.SetMaxOpenConns(1)
	db.Exec("SET jit = off") // speed only: an error leaves JIT on
	return db
}
