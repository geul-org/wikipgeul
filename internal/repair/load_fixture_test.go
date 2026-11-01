//ff:func feature=repair type=writer control=iteration dimension=1
//ff:what loadFixture — a throwaway schema loaded from testdata/fixture.jsonl and damaged like the 2026-10-27 cleanup: mixed, reordered and missing snak lists, a lost statement_reference row, a broken line
package repair

import (
	"database/sql"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/geul-org/wikipgeul/internal/dump"
	"github.com/geul-org/wikipgeul/internal/load"
)

// loadFixture loads testdata/fixture.jsonl (Q1, Q2, Q3 on lines 2–4) with
// the loader in line order, so each hash keeps its first occurrence, drops
// the staging columns as wikidata.sql 4.2 does, and damages the result:
// hx (Q1: P248, P813, P123 ×2) gets P813 twice and loses P123 Q7809 (content),
// hy takes Q3's order P248, P854 (order), hz loses its row (missing), hw stays
// (same), Q3's statement_reference row for hw is gone, and line 5 (Q4) is
// broken JSON. It returns the database and a dry-run Config over all lines
// in batches of 2 with 2 workers.
func loadFixture(t *testing.T) (*sql.DB, Config) {
	t.Helper()
	db, schema := testSchema(t)
	b, err := os.ReadFile("testdata/fixture.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	r, refs := load.NewRows(), load.NewRefSet(64)
	for i, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		e, err := dump.DecodeEntity([]byte(line))
		if err != nil {
			t.Fatal(err)
		}
		r.AddEntity(load.Config{DumpID: 1, Raw: true}, int64(i+2), []byte(line), e, refs)
	}
	if err := load.WriteRows(db, schema, r); err != nil {
		t.Fatal(err)
	}
	s := schema + "."
	for _, q := range []string{
		"ALTER TABLE " + s + "reference_snak DROP COLUMN occ_statement_id, DROP COLUMN occ_ordinal",
		"INSERT INTO " + s + "entity_raw VALUES ('Q4', 1, 5, '{bad')",
		"INSERT INTO " + s + "entity VALUES ('Q4', 5, 'Q', 4, 'item', NULL, NULL, NULL, 1)",
		"UPDATE " + s + "reference_snak d SET property = o.property, value_ordinal = o.value_ordinal, datatype = o.datatype, value_type = o.value_type, value = o.value " +
			"FROM " + s + "reference_snak o WHERE d.hash = 'hx' AND d.ordinal = 2 AND o.hash = 'hx' AND o.ordinal = 1",
		"UPDATE " + s + "reference_snak SET ordinal = 1 - ordinal, value_ordinal = 0 WHERE hash = 'hy'",
		"DELETE FROM " + s + "reference_snak WHERE hash = 'hz'",
		"DELETE FROM " + s + "statement_reference WHERE hash = 'hw' AND statement_id = ((1::bigint << 55) | (4::bigint << 20))",
	} {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(q, err)
		}
	}
	return db, Config{Schema: schema, DumpID: 1, DryRun: true, Batch: 2, Workers: 2, Status: time.Hour, RefLimit: 64, Log: io.Discard}
}
