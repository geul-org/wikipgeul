//ff:func feature=convert type=command control=sequence
//ff:what TestPortOne — Q42 is published under its number; Q30's somevalue/novalue snaks become slot markers; a missing ID is false and reported Missing; a malformed digest, a broken published version and an unwritable tree fail
package exportgeul

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/resolve"
	"github.com/geul-org/wikigeul/pkg/pipeline"

	"github.com/geul-org/wikipgeul/internal/pgsource"
)

func TestPortOne(t *testing.T) {
	db, schema := fixtureDB(t)
	root := t.TempDir()
	c := config{Config: pipeline.Config{Commit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", LibCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Schema: schema, Dump: 1}}
	prov, _ := provenanceOf(pgsource.Dump{URL: "u", Loader: wikipgeulPath, LoaderCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, c.Commit, c.LibCommit)
	rep, ok, err := portOne(db, resolve.Dir{Root: root}, "Q42", []core.Word{1, 0}, c, prov)
	if err != nil || !ok || rep.Missing || rep.Triples == 0 || rep.Qualifiers == 0 {
		t.Fatalf("Q42: %+v %v %v", rep, ok, err)
	}
	if _, err := os.Stat(filepath.Join(root, ".well-known/geul/0x00010000/0x510000000000002a/0x00000001.geul")); err != nil {
		t.Fatal(err)
	}
	q30, ok, err := portOne(db, resolve.Dir{Root: root}, "Q30", []core.Word{1, 0}, c, prov)
	if err != nil || !ok || q30.Triples != 3 || q30.Qualifiers != 3 || len(q30.Skipped) != 0 ||
		q30.Markers["NO-VALUE"] != 1 || q30.Markers["UNKNOWN-VALUE"] != 1 || q30.Markers["qualifier UNKNOWN-VALUE"] != 1 {
		t.Fatalf("Q30: %+v %v %v", q30, ok, err)
	}
	if q6, ok, err := portOne(db, resolve.Dir{Root: root}, "Q6", []core.Word{1, 0}, c, prov); ok || err != nil || !q6.Missing || q6.ID != "Q6" {
		t.Fatalf("Q6: %+v %v %v", q6, ok, err)
	}
	bad := prov
	bad.DumpSHA256 = "zz"
	if _, _, err := portOne(db, resolve.Dir{Root: t.TempDir()}, "Q42", []core.Word{1, 0}, c, bad); err == nil {
		t.Fatal("malformed dump digest: want error")
	}
	broken := t.TempDir()
	vdir := filepath.Join(broken, ".well-known/geul/0x00010000/0x510000000000002a")
	os.MkdirAll(vdir, 0o755)
	os.WriteFile(filepath.Join(vdir, "versions.geul"), []byte{1, 2, 3}, 0o644)
	if _, _, err := portOne(db, resolve.Dir{Root: broken}, "Q42", []core.Word{1, 0}, c, prov); err == nil {
		t.Fatal("broken published version list: want error")
	}
	noH := t.TempDir()
	os.MkdirAll(filepath.Join(noH, ".well-known/geul"), 0o755)
	os.WriteFile(filepath.Join(noH, ".well-known/geul/h"), nil, 0o644)
	if _, _, err := portOne(db, resolve.Dir{Root: noH}, "Q42", []core.Word{1, 0}, c, prov); err == nil {
		t.Fatal("h/ not a directory: want error")
	}
	blocked := filepath.Join(root, "blocked")
	os.WriteFile(blocked, nil, 0o644)
	if _, _, err := portOne(db, resolve.Dir{Root: blocked}, "Q42", []core.Word{1, 0}, c, prov); err == nil {
		t.Fatal("unwritable: want error")
	}
}
