//ff:func feature=repair type=helper control=sequence
//ff:what TestFilterDone — not done and done later are processed, done at the same or an earlier occurrence skipped
package repair

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestFilterDone(t *testing.T) {
	at := func(s int64) load.Occurrence { return load.Occurrence{Statement: s} }
	checks := []HashCheck{{Hash: "new", At: at(5)}, {Hash: "later", At: at(5)}, {Hash: "same", At: at(5)}, {Hash: "earlier", At: at(5)}}
	got, skipped := filterDone(checks, map[string]load.Occurrence{"later": at(9), "same": at(5), "earlier": at(1)})
	if skipped != 2 || len(got) != 2 || got[0].Hash != "new" || got[1].Hash != "later" {
		t.Fatalf("%v %d", got, skipped)
	}
}
