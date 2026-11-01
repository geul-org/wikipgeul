//ff:func feature=load type=model control=sequence
//ff:what TestReferenceSnaks — rows follow this occurrence's snaks-order with overall and per-property ordinals and the snak columns
package load

import (
	"encoding/json"
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
)

func TestReferenceSnaks(t *testing.T) {
	item := func(p, q string) dump.Snak {
		return dump.Snak{SnakType: "value", Property: p, Datatype: "wikibase-item", DataValue: &dump.DataValue{Type: "wikibase-entityid", Value: json.RawMessage(`{"id":"` + q + `"}`)}}
	}
	ref := dump.Reference{Hash: "h", SnaksOrder: []string{"P854", "P248"},
		Snaks: map[string][]dump.Snak{"P248": {item("P248", "Q1"), item("P248", "Q2")}, "P854": {{SnakType: "somevalue", Property: "P854"}}}}
	rows := ReferenceSnaks(3, ref)
	if len(rows) != 3 || len(rows[0]) != 11 {
		t.Fatalf("rows %v", rows)
	}
	if r := rows[0]; r[0] != "h" || r[1] != int16(3) || r[2] != 0 || r[3] != "P854" || r[4] != 0 || r[5] != "somevalue" || r[8] != nil {
		t.Fatalf("first %v", r)
	}
	if r := rows[2]; r[2] != 2 || r[3] != "P248" || r[4] != 1 || r[9] != "Q2" {
		t.Fatalf("last %v", r)
	}
	if ReferenceSnaks(1, dump.Reference{Hash: "e"}) != nil {
		t.Fatal("no snaks: no rows")
	}
}
