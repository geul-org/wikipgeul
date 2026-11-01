//ff:func feature=load type=writer control=iteration dimension=1
//ff:what TestClearAfter — rows of lines after the checkpoint go from every entity table and their reference_snak occurrences, earlier lines stay; a closed handle and a missing schema fail
package load

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
)

func TestClearAfter(t *testing.T) {
	db, schema := testSchema(t)
	InsertDump(db, schema, sampleDump())
	r := NewRows()
	refs := NewRefSet(64)
	r.AddEntity(Config{DumpID: 1, Raw: true}, 5, []byte(`{}`), sampleEntity(), refs)
	later := sampleEntity()
	later.ID = "Q43"
	r.AddEntity(Config{DumpID: 1, Raw: true}, 6, []byte(`{}`), later, nil)
	r.Error(1, 6, "Q43", "x", "y", nil)
	r.AddEntity(Config{DumpID: 1}, 7, nil, dump.Entity{ID: "Q44"}, refs)
	if err := WriteRows(db, schema, r); err != nil {
		t.Fatal(err)
	}
	if err := ClearAfter(db, schema, 1, 5); err != nil {
		t.Fatal(err)
	}
	for table, want := range map[string]int{"entity_raw": 1, "entity": 1, "label": 2, "description": 2, "alias": 2, "sitelink": 1, "statement": 3, "qualifier": 2, "statement_reference": 1, "reference_snak": 1, "load_error": 0} {
		var n int
		db.QueryRow("SELECT count(*) FROM " + schema + "." + table).Scan(&n)
		if n != want {
			t.Fatalf("%s: %d want %d", table, n, want)
		}
	}
	if err := ClearAfter(closedDB(t), schema, 1, 5); err == nil {
		t.Fatal("closed: want error")
	}
	if err := ClearAfter(db, schema+"_none", 1, 5); err == nil {
		t.Fatal("missing schema: want error")
	}
}
