//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestRunOrder — a dry run over a loaded schema counts the occurrence in another order and writes nothing; a write records it and leaves reference_snak alone; bad flags fail
package main

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
	"github.com/geul-org/wikipgeul/internal/load"
)

func TestRunOrder(t *testing.T) {
	db, schema := testSchema(t)
	r, refs := load.NewRows(), load.NewRefSet(64)
	for i, raw := range []string{
		`{"id":"Q1","type":"item","claims":{"P31":[{"id":"Q1$1","rank":"normal","mainsnak":{"snaktype":"novalue","property":"P31"},"references":[{"hash":"r","snaks":{"P143":[{"snaktype":"novalue","property":"P143"}],"P813":[{"snaktype":"somevalue","property":"P813"}]},"snaks-order":["P143","P813"]}]}]}}`,
		`{"id":"Q2","type":"item","claims":{"P31":[{"id":"Q2$1","rank":"normal","mainsnak":{"snaktype":"novalue","property":"P31"},"references":[{"hash":"r","snaks":{"P143":[{"snaktype":"novalue","property":"P143"}],"P813":[{"snaktype":"somevalue","property":"P813"}]},"snaks-order":["P813","P143"]}]}]}}`,
	} {
		e, err := dump.DecodeEntity([]byte(raw))
		if err != nil {
			t.Fatal(err)
		}
		r.AddEntity(load.Config{DumpID: 1, Raw: true}, int64(i+2), []byte(raw), e, refs)
	}
	if err := load.WriteRows(db, schema, r); err != nil {
		t.Fatal(err)
	}
	db.Exec("ALTER TABLE " + schema + ".reference_snak DROP COLUMN occ_statement_id, DROP COLUMN occ_ordinal")
	if err := runOrder([]string{"-schema", schema, "-dry-run", "-workers", "1"}); err != nil {
		t.Fatal(err)
	}
	var n int
	if db.QueryRow("SELECT count(*) FROM information_schema.tables WHERE table_schema = $1 AND table_name = 'reference_snak_order'", schema).Scan(&n); n != 0 {
		t.Fatal("dry run wrote")
	}
	if err := runOrder([]string{"-schema", schema, "-workers", "1"}); err != nil {
		t.Fatal(err)
	}
	var perm string
	if err := db.QueryRow("SELECT perm::text FROM " + schema + ".reference_snak_order WHERE statement_id = (1::bigint << 55) | (3::bigint << 20) AND ordinal = 0").Scan(&perm); err != nil || perm != "{1,0}" {
		t.Fatalf("order row %q %v", perm, err)
	}
	if db.QueryRow("SELECT count(*) FROM " + schema + ".reference_snak WHERE hash = 'r' AND (ordinal, property) IN ((0, 'P143'), (1, 'P813'))").Scan(&n); n != 2 {
		t.Fatalf("reference_snak changed: %d", n)
	}
	if err := runOrder([]string{"-batch", "0"}); err == nil {
		t.Fatal("bad flags: want error")
	}
}
