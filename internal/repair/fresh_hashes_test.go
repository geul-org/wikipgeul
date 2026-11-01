//ff:func feature=repair type=helper control=sequence
//ff:what TestFreshHashes — a hash already handled at an earlier occurrence is left out, an earlier occurrence is kept
package repair

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestFreshHashes(t *testing.T) {
	refs := load.NewRefSet(64)
	at := func(s int64) load.Occurrence { return load.Occurrence{Statement: s} }
	if got := freshHashes(Found{Hashes: []HashCheck{{Hash: "a", At: at(5)}, {Hash: "b", At: at(5)}}}, refs); len(got) != 2 {
		t.Fatalf("first: %v", got)
	}
	got := freshHashes(Found{Hashes: []HashCheck{{Hash: "a", At: at(9)}, {Hash: "b", At: at(1)}, {Hash: "c", At: at(9)}}}, refs)
	if len(got) != 2 || got[0].Hash != "b" || got[1].Hash != "c" {
		t.Fatalf("second: %v", got)
	}
}
