//ff:func feature=load type=writer control=iteration dimension=1
//ff:what TestReferenceDedup — a hash met in two orders, loaded out of dump order and twice, keeps after wikidata.sql 4.2 exactly its first occurrence's snak list
package load

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
)

func TestReferenceDedup(t *testing.T) {
	db, schema := testSchema(t)
	InsertDump(db, schema, sampleDump())
	item := func(p, q string) dump.Snak {
		return dump.Snak{SnakType: "value", Property: p, Datatype: "wikibase-item", DataValue: &dump.DataValue{Type: "wikibase-entityid", Value: json.RawMessage(`{"id":"` + q + `"}`)}}
	}
	snaks := map[string][]dump.Snak{"P248": {item("P248", "Q1")}, "P123": {item("P123", "Q7809"), item("P123", "Q2")}, "P813": {{SnakType: "novalue", Property: "P813"}}}
	entity := func(id string, order ...string) dump.Entity {
		ref := dump.Reference{Hash: "hx", Snaks: snaks, SnaksOrder: order}
		return dump.Entity{ID: id, Type: "item", Claims: map[string][]dump.Statement{"P31": {{ID: id + "$1", Rank: "normal",
			Mainsnak: item("P31", "Q5"), References: []dump.Reference{ref}}}}}
	}
	refs := NewRefSet(64)
	later, first := NewRows(), NewRows()
	later.AddEntity(Config{DumpID: 1}, 20, nil, entity("Q20", "P813", "P248", "P123"), refs)
	first.AddEntity(Config{DumpID: 1}, 10, nil, entity("Q10", "P123", "P813", "P248"), refs)
	again := NewRows()
	again.AddEntity(Config{DumpID: 1}, 10, nil, entity("Q10", "P123", "P813", "P248"), nil)
	for _, r := range []*Rows{later, first, again} {
		if err := WriteRows(db, schema, r); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := db.Exec(sectionSQL(t, "-- 4.2 ", schema)); err != nil {
		t.Fatal(err)
	}
	var got string
	if err := db.QueryRow("SELECT string_agg(property || '=' || coalesce(value->>'id', snaktype::text) || '/' || value_ordinal, ' ' ORDER BY ordinal) FROM " + schema + ".reference_snak WHERE hash = 'hx'").Scan(&got); err != nil {
		t.Fatal(err)
	}
	if want := "P123=Q7809/0 P123=Q2/1 P813=novalue/0 P248=Q1/0"; got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	var refRows, cols int
	db.QueryRow("SELECT count(*) FROM " + schema + ".reference").Scan(&refRows)
	db.QueryRow("SELECT count(*) FROM information_schema.columns WHERE table_schema = $1 AND table_name = 'reference_snak' AND column_name LIKE 'occ%'", strings.ToLower(schema)).Scan(&cols)
	if refRows != 1 || cols != 0 {
		t.Fatalf("reference %d rows, staging columns %d", refRows, cols)
	}
}
