//ff:func feature=load type=model control=iteration dimension=1
//ff:what TestRowsAddEntity — the sample item gives raw, head, terms, sitelink, statements, qualifiers and references; raw off and a bad ID behave
package load

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
)

func TestRowsAddEntity(t *testing.T) {
	r := NewRows()
	r.AddEntity(Config{DumpID: 1, Raw: true}, 5, []byte(`{"id":"Q42"}`), sampleEntity(), NewRefSet(64))
	for table, n := range map[string]int{"entity_raw": 1, "entity": 1, "label": 2, "description": 2, "alias": 2, "sitelink": 1, "statement": 3, "qualifier": 2, "statement_reference": 1, "reference": 1, "reference_snak": 1} {
		if r.Count(table) != n {
			t.Fatalf("%s: %d want %d", table, r.Count(table), n)
		}
	}
	head := r.byTable["entity"][0]
	if head[1] != int64(5) || head[2] != "Q" || head[3] != int64(42) || head[5] != nil || head[6] != int64(2401234567) {
		t.Fatalf("head %v", head)
	}
	off := NewRows()
	off.AddEntity(Config{DumpID: 1}, 5, nil, dump.Entity{ID: "Q1"}, NewRefSet(64))
	if off.Count("entity_raw") != 0 || off.Count("entity") != 1 || off.byTable["entity"][0][6] != nil {
		t.Fatalf("raw off %+v", off.byTable)
	}
	bad := NewRows()
	bad.AddEntity(Config{DumpID: 1, Raw: true}, 6, []byte("x"), dump.Entity{ID: "Qx"}, NewRefSet(64))
	if bad.Count("entity") != 0 || bad.Errors != 1 {
		t.Fatalf("bad ID %+v", bad.byTable)
	}
	nul := NewRows()
	nul.AddEntity(Config{DumpID: 1, Raw: true}, 7, []byte("{\x00}"), dump.Entity{ID: "Q7"}, NewRefSet(64))
	if nul.Errors != 1 || nul.byTable["entity_raw"][0][3] != "{}" {
		t.Fatalf("NUL raw %+v", nul.byTable)
	}
}
