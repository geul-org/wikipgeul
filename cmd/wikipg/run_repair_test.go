//ff:func feature=cli type=command control=sequence
//ff:what TestRunRepair — a dry run over a loaded schema reports and writes nothing; a write repairs; bad flags and an unreachable database fail
package main

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
	"github.com/geul-org/wikipgeul/internal/load"
)

func TestRunRepair(t *testing.T) {
	db, schema := testSchema(t)
	raw := []byte(`{"id":"Q1","type":"item","claims":{"P31":[{"id":"Q1$1","rank":"normal","mainsnak":{"snaktype":"novalue","property":"P31"},` +
		`"references":[{"hash":"r","snaks":{"P143":[{"snaktype":"novalue","property":"P143"}]},"snaks-order":["P143"]}]}]}}`)
	e, err := dump.DecodeEntity(raw)
	if err != nil {
		t.Fatal(err)
	}
	r := load.NewRows()
	r.AddEntity(load.Config{DumpID: 1, Raw: true}, 2, raw, e, nil)
	if err := load.WriteRows(db, schema, r); err != nil {
		t.Fatal(err)
	}
	db.Exec("ALTER TABLE " + schema + ".reference_snak DROP COLUMN occ_statement_id, DROP COLUMN occ_ordinal")
	db.Exec("DELETE FROM " + schema + ".reference_snak")
	if err := runRepair([]string{"-schema", schema, "-dry-run", "-workers", "1"}); err != nil {
		t.Fatal(err)
	}
	var n int
	if db.QueryRow("SELECT count(*) FROM " + schema + ".reference_snak").Scan(&n); n != 0 {
		t.Fatal("dry run wrote")
	}
	if err := runRepair([]string{"-schema", schema, "-workers", "1"}); err != nil {
		t.Fatal(err)
	}
	if db.QueryRow("SELECT count(*) FROM " + schema + ".reference_snak").Scan(&n); n != 1 {
		t.Fatalf("repaired rows %d", n)
	}
	if err := runRepair([]string{"-batch", "0"}); err == nil {
		t.Fatal("bad flags: want error")
	}
	if err := runRepair([]string{"-dsn", "postgres://%zz", "-schema", schema}); err == nil {
		t.Fatal("bad DSN URL: want error")
	}
	if err := runRepair([]string{"-dsn", "host='x", "-schema", schema}); err == nil {
		t.Fatal("unparsable DSN: want error")
	}
}
