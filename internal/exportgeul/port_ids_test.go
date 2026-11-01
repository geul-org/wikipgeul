//ff:func feature=convert type=command control=sequence
//ff:what TestPortIDs — present IDs are written and missing ones only reported; the first error stops the run
package exportgeul

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/resolve"
	"github.com/geul-org/wikigeul/pkg/pipeline"

	"github.com/geul-org/wikipgeul/internal/pgsource"
)

func TestPortIDs(t *testing.T) {
	db, schema := fixtureDB(t)
	c := config{Config: pipeline.Config{Commit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", LibCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Schema: schema, Dump: 1}}
	prov, _ := provenanceOf(pgsource.Dump{URL: "u", Loader: wikipgeulPath, LoaderCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, c.Commit, c.LibCommit)
	reps, written, err := portIDs(db, resolve.Dir{Root: t.TempDir()}, "Q", []uint64{1, 6}, []core.Word{1, 0}, c, prov)
	if err != nil || len(reps) != 2 || !slices.Equal(written, []uint64{0x5100_0000_0000_0001}) {
		t.Fatalf("got %d %v %v", len(reps), written, err)
	}
	blocked := filepath.Join(t.TempDir(), "blocked")
	os.WriteFile(blocked, nil, 0o644)
	if _, _, err := portIDs(db, resolve.Dir{Root: blocked}, "Q", []uint64{1}, []core.Word{1, 0}, c, prov); err == nil {
		t.Fatal("unwritable: want error")
	}
}
