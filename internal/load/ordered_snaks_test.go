//ff:func feature=load type=helper control=iteration dimension=1
//ff:what TestOrderedSnaks — the order list first with value ordinals per property, missing properties after in sorted order
package load

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
)

func TestOrderedSnaks(t *testing.T) {
	m := map[string][]dump.Snak{"P2": {{Hash: "a"}, {Hash: "b"}}, "P1": {{Hash: "c"}}, "P9": {{Hash: "z"}}, "P5": {{Hash: "y"}}}
	got := orderedSnaks([]string{"P2", "P1"}, m)
	want := []string{"a", "b", "c", "y", "z"}
	if len(got) != 5 {
		t.Fatalf("got %d", len(got))
	}
	for i, s := range got {
		if s.Snak.Hash != want[i] || s.Ordinal != i {
			t.Fatalf("%d: %+v", i, s)
		}
	}
	if got[1].ValueOrdinal != 1 || got[2].ValueOrdinal != 0 {
		t.Fatal("value ordinals")
	}
	if orderedSnaks(nil, nil) != nil {
		t.Fatal("empty: want nil")
	}
}
