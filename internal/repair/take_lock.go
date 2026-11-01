//ff:func feature=repair type=writer control=sequence
//ff:what takeLock — holds the repair's advisory lock on a connection of its own, so two writing runs never overlap
package repair

import (
	"context"
	"database/sql"
)

// lockKey names the advisory lock of writing runs (with the schema). Repairs
// and order builds share it: an order build reads the reference_snak rows a
// repair replaces, so the two never write at once.
const lockKey = "wikipg repair-reference-snak"

// takeLock takes the session advisory lock of the repair of c.Schema on a
// connection it keeps (released by releaseLock, or when the process ends). Another
// writing run holding it is an error (lockError).
func takeLock(ctx context.Context, db *sql.DB, c Config) (*sql.Conn, error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	var ok bool
	err = conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock(hashtext($1))", lockKey+" "+c.Schema).Scan(&ok)
	if err != nil || !ok {
		conn.Close()
		return nil, lockError(err)
	}
	return conn, nil
}
