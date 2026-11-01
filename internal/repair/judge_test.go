//ff:func feature=repair type=helper control=sequence
//ff:what TestJudge — every occurrence is counted, entities listed once per kind of difference, each hash checked at its first occurrence, occurrences in another order listed with their permutation
package repair

import (
	"fmt"
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestJudge(t *testing.T) {
	row := func(ord int, p string) []any {
		return []any{"", int16(1), ord, p, 0, "novalue", nil, nil, nil, nil, nil}
	}
	ab, ba := [][]any{row(0, "A"), row(1, "B")}, [][]any{row(0, "B"), row(1, "A")}
	stored := map[string][]Snak{"h": {snakOfRow(row(0, "A"), nil), snakOfRow(row(1, "B"), nil)}, "k": {snakOfRow(row(0, "A"), nil)}}
	d := spanData{Occs: []occ{
		{Hash: "h", Entity: "Q1", Rows: ba, At: load.Occurrence{Statement: 1}},
		{Hash: "h", Entity: "Q1", Rows: ba, At: load.Occurrence{Statement: 2}},
		{Hash: "k", Entity: "Q1", Rows: ab}, {Hash: "k", Entity: "Q1", Rows: ab},
		{Hash: "h", Entity: "Q2", Rows: ab},
	}}
	f := judge(d, stored, nil)
	if f.Occurrences != (Tally{Same: 1, Order: 2, Content: 2}) || len(f.ContentEntities) != 1 || len(f.OrderEntities) != 1 {
		t.Fatalf("%+v", f)
	}
	if len(f.Hashes) != 2 || f.Hashes[0].Verdict != Order || f.Hashes[0].At.Statement != 1 || f.Hashes[1].Verdict != Content {
		t.Fatalf("hashes %+v", f.Hashes)
	}
	if len(f.Orders) != 2 || f.Orders[0].At.Statement != 1 || f.Orders[1].At.Statement != 2 || f.Orders[1].Hash != "h" || f.Orders[1].Entity != "Q1" || fmt.Sprint(f.Orders[0].Perm) != "[1 0]" {
		t.Fatalf("orders %+v", f.Orders)
	}
}
