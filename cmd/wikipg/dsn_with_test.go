//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what TestDsnWith — a parameter joins key=value, empty and URL DSNs; a session opened read-only cannot write
package main

import (
	"database/sql"
	"testing"
)

func TestDsnWith(t *testing.T) {
	const ro = "default_transaction_read_only=on"
	for in, want := range map[string]string{
		"":                            ro,
		"host=x":                      "host=x " + ro,
		"postgres://h/db":             "postgres://h/db?" + ro,
		"postgresql://h/db?sslmode=x": "postgresql://h/db?sslmode=x&" + ro,
	} {
		if got := dsnWith(in, ro); got != want {
			t.Fatalf("%q: %q want %q", in, got, want)
		}
	}
	db, err := sql.Open("postgres", dsnWith(dsnWith("", "jit=off"), ro))
	if err == nil {
		err = db.Ping()
	}
	if err != nil {
		t.Skipf("database not reachable: %v", err)
	}
	defer db.Close()
	var jit string
	if db.QueryRow("SHOW jit").Scan(&jit); jit != "off" {
		t.Fatalf("jit %q", jit)
	}
	if _, err := db.Exec("CREATE TEMP TABLE read_only_probe (x int)"); err == nil {
		t.Fatal("a read-only session created a table")
	}
}
