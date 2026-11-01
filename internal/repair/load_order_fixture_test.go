//ff:func feature=repair type=writer control=iteration dimension=1
//ff:what loadOrderFixture — a throwaway schema loaded from testdata/fixture.jsonl as wikidata.sql 4.2 leaves it, undamaged: each hash keeps its first occurrence, later ones may differ in order
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

// loadOrderFixture loads testdata/fixture.jsonl (Q1, Q2, Q3 on lines 2–4)
// with the loader in line order, so each hash keeps its first occurrence, and
// drops the staging columns as wikidata.sql 4.2 does. Two later occurrences
// have the stored snaks in another order: hx on Q2 (P123 ×2, P813, P248
// against P248, P813, P123 ×2) and hy on Q3 (P248, P854 against P854, P248).
// It returns the database and a dry-run order-build Config over all lines in
// batches of 2 with 2 workers.
func loadOrderFixture(t *testing.T) (*sql.DB, Config) {
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
	if _, err := db.Exec("ALTER TABLE " + schema + ".reference_snak DROP COLUMN occ_statement_id, DROP COLUMN occ_ordinal"); err != nil {
		t.Fatal(err)
	}
	return db, Config{Schema: schema, DumpID: 1, Order: true, DryRun: true, Batch: 2, Workers: 2, Status: time.Hour, RefLimit: 64, Log: io.Discard}
}
