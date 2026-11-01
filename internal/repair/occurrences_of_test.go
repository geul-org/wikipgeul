//ff:func feature=repair type=decoder control=sequence
//ff:what TestOccurrencesOf — every occurrence keeps its own snak order and key, a reference without snaks still counts, bad JSON and a loader error fail
package repair

import (
	"os"
	"strings"
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestOccurrencesOf(t *testing.T) {
	b, err := os.ReadFile("testdata/fixture.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(string(b), "\n")
	id, occs, refs, err := occurrencesOf(1, 3, lines[1])
	if err != nil || id != "Q2" || len(occs) != 3 || len(refs) != 3 {
		t.Fatalf("Q2: %s %v %v %v", id, occs, refs, err)
	}
	// Q2: P17 (hz, hw) sorts before P31 (hx), keys continue on line 3.
	if o := occs[2]; o.Hash != "hx" || o.At != (load.Occurrence{Statement: 1<<55 | 3<<20 | 1, Ordinal: 0}) || o.Entity != "Q2" || len(o.Rows) != 4 || o.Rows[0][3] != "P123" {
		t.Fatalf("hx: %+v", o)
	}
	if refs[1] != (stmtRef{Statement: 1<<55 | 3<<20, Ordinal: 1, Hash: "hw"}) {
		t.Fatalf("refs %v", refs)
	}
	_, occs, _, err = occurrencesOf(1, 9, `{"id":"Q9","claims":{"P1":[{"id":"Q9$1","rank":"normal","mainsnak":{"snaktype":"novalue","property":"P1"},"references":[{"hash":"e"}]}]}}`)
	if err != nil || len(occs) != 1 || occs[0].Rows != nil {
		t.Fatalf("empty reference: %v %v", occs, err)
	}
	if _, _, _, err := occurrencesOf(1, 9, "{bad"); err == nil {
		t.Fatal("bad JSON: want error")
	}
	if id, _, _, err := occurrencesOf(1, 9, `{"id":"bad"}`); err == nil || id != "bad" {
		t.Fatalf("unsplittable ID: %q %v", id, err)
	}
}
