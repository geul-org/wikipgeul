//ff:func feature=convert type=command control=sequence
//ff:what TestOpenDB — one connection with JIT off
package exportgeul

import (
	"testing"

	"github.com/lib/pq"
)

func TestOpenDB(t *testing.T) {
	fixtureDB(t) // skips without a database
	conn, _ := pq.NewConnector("")
	db := openDB(conn)
	defer db.Close()
	var jit string
	if err := db.QueryRow("SHOW jit").Scan(&jit); err != nil || jit != "off" || db.Stats().MaxOpenConnections != 1 {
		t.Fatalf("%q %v", jit, err)
	}
}
