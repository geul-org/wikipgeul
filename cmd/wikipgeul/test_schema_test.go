//ff:func feature=cli type=command control=sequence
//ff:what testSchema — a throwaway schema of sql/wikidata.sql, named tsma_test_*, dropped after the test; skips when the database is not reachable
package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// testSchema opens the database of lib/pq's PG* environment variables, runs
// sections 1–3 of sql/wikidata.sql under a fresh schema (the schema wikidata
// is replaced) and drops that schema when the test ends. It skips the test when
// the database is not reachable.
func testSchema(t *testing.T) (*sql.DB, string) {
	t.Helper()
	db, err := sql.Open("postgres", "")
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		t.Skipf("database not reachable: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	b, err := os.ReadFile("../../sql/wikidata.sql")
	if err != nil {
		t.Fatal(err)
	}
	name := fmt.Sprintf("tsma_test_%d", time.Now().UnixNano())
	head := "CREATE SCHEMA IF NOT EXISTS wikidata;\nSET search_path = wikidata;"
	if !strings.Contains(string(b), head) {
		t.Fatal("sql/wikidata.sql: schema header not found; refusing to run the DDL outside a test schema")
	}
	ddl := strings.Replace(string(b), head, "CREATE SCHEMA "+name+";\nSET search_path = "+name+";", 1)
	if _, err := db.Exec(ddl); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Exec("DROP SCHEMA " + name + " CASCADE") })
	return db, name
}
