//ff:func feature=repair type=reader control=sequence
//ff:what TestCheckKey — no table, a table with the key (dump_id, run), a table of before 2026-10-29 (key run) failing until sql/migrate_checkpoint_key.sql moves it with its rows; a failed query
package repair

import (
	"strings"
	"testing"

	schemasql "github.com/geul-org/wikipgeul/sql"
)

func TestCheckKey(t *testing.T) {
	db, schema := testSchema(t)
	c := Config{Schema: schema, DumpID: 2, Order: true}
	if ok, err := checkKey(db, c); ok || err != nil {
		t.Fatalf("no table: %v %v", ok, err)
	}
	s := schema + "."
	old := "CREATE TABLE " + s + "reference_snak_order_checkpoint (run text PRIMARY KEY, dump_id smallint NOT NULL, line_no bigint NOT NULL, updated timestamptz NOT NULL DEFAULT now());\n" +
		"INSERT INTO " + s + "reference_snak_order_checkpoint (run, dump_id, line_no) VALUES ('lines:0-0', 1, 121971282), ('lines:2-0', 2, 1707633)"
	if _, err := db.Exec(old); err != nil {
		t.Fatal(err)
	}
	if ok, err := checkKey(db, c); !ok || err == nil || !strings.Contains(err.Error(), "(run)") || !strings.Contains(err.Error(), "migrate_checkpoint_key.sql") {
		t.Fatalf("old key: %v %v", ok, err)
	}
	migrate := "SET LOCAL search_path = " + schema + ";\n" + schemasql.MigrateCheckpointKey
	if _, err := db.Exec(migrate); err == nil || !strings.Contains(err.Error(), "dump 1") {
		t.Fatalf("rows of dumps not in the dump table: want error, got %v", err)
	}
	if ok, err := checkKey(db, c); !ok || err == nil {
		t.Fatalf("failed migration left the key: %v %v", ok, err)
	}
	if _, err := db.Exec("INSERT INTO " + s + "dump (dump_id, dump_date, file_name, url, size_bytes, sha256, downloaded_at, loader, loader_commit, checkpoint_line) VALUES " +
		"(1, '2026-10-19', 'a', 'u', 1, 'h', now(), 'l', 'c', 121971282), (2, '2026-10-21', 'b', 'u', 1, 'h', now(), 'l', 'c', 1707633)"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec(migrate); err != nil {
		t.Fatal("migrate:", err)
	}
	if _, err := db.Exec(migrate); err != nil {
		t.Fatal("migrate again (a no-op):", err)
	}
	if ok, err := checkKey(db, c); !ok || err != nil {
		t.Fatalf("migrated: %v %v", ok, err)
	}
	if n, err := readCheckpoint(db, c); err != nil || n != 0 {
		t.Fatalf("dump 2 lines:0-0 after migration: %d %v", n, err)
	}
	tx, _ := db.Begin()
	saveCheckpoint(tx, c, 99)
	tx.Commit()
	c.DumpID = 1
	if n, err := readCheckpoint(db, c); err != nil || n != 121971282 {
		t.Fatalf("dump 1 kept its line: %d %v", n, err)
	}
	c.DumpID, c.From = 2, 2
	if n, err := readCheckpoint(db, c); err != nil || n != 1707633 {
		t.Fatalf("dump 2 lines:2-0 kept: %d %v", n, err)
	}
	if _, err := checkKey(closedDB(t), c); err == nil {
		t.Fatal("closed database: want error")
	}
}
