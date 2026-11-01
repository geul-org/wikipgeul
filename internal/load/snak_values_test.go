//ff:func feature=load type=helper control=sequence
//ff:what TestSnakValues — a value snak gives type, JSON and entity; novalue gives NULLs; a \u0000 value goes to value_json
package load

import (
	"encoding/json"
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
)

func TestSnakValues(t *testing.T) {
	v := snakValues(dump.Snak{SnakType: "value", Property: "P31", Datatype: "wikibase-item", DataValue: &dump.DataValue{Type: "wikibase-entityid", Value: json.RawMessage(`{"id":"Q5"}`)}})
	if len(v) != 6 || v[0] != "value" || v[1] != "wikibase-item" || v[2] != "wikibase-entityid" || v[3] != `{"id":"Q5"}` || v[4] != "Q5" || v[5] != nil {
		t.Fatalf("value: %v", v)
	}
	n := snakValues(dump.Snak{SnakType: "novalue", Property: "P40"})
	if n[0] != "novalue" || n[1] != nil || n[2] != nil || n[3] != nil || n[4] != nil || n[5] != nil {
		t.Fatalf("novalue: %v", n)
	}
	z := snakValues(dump.Snak{SnakType: "value", Property: "P1", DataValue: &dump.DataValue{Type: "string", Value: json.RawMessage(`"a\u0000"`)}})
	if z[2] != "string" || z[3] != nil || z[4] != nil || z[5] != `"a\u0000"` {
		t.Fatalf("nul: %v", z)
	}
}
