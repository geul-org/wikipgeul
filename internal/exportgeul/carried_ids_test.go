//ff:func feature=convert type=helper control=iteration dimension=1
//ff:what TestCarriedIDs — an unpublished entity starts at 1; a published one continues above its IDs; a malformed ID, a number past 56 bits and a broken version list fail
package exportgeul

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/resolve"
	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestCarriedIDs(t *testing.T) {
	dir := resolve.Dir{Root: t.TempDir()}
	issuer := []core.Word{1, 0}
	if ids, err := carriedIDs(dir, issuer, "Q42"); err != nil || ids.Next() != 1 {
		t.Fatalf("unpublished: %v", err)
	}
	db, schema := fixtureDB(t)
	c := config{Config: pipeline.Config{Commit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", LibCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", Schema: schema, Dump: 1}}
	prov, _ := provenanceOf(fixtureDump(), c.Commit, c.LibCommit)
	if _, ok, err := portOne(db, dir, "Q42", issuer, c, prov); err != nil || !ok {
		t.Fatalf("port: %v", err)
	}
	if ids, err := carriedIDs(dir, issuer, "Q42"); err != nil || ids.Next() < 2 {
		t.Fatalf("published: %v", err)
	}
	for _, id := range []string{"Qx", "Q72057594037927936"} {
		if _, err := carriedIDs(dir, issuer, id); err == nil {
			t.Fatalf("%s: want error", id)
		}
	}
	root := t.TempDir()
	list := filepath.Join(root, ".well-known/geul/0x00010000/0x5100000000000001")
	os.MkdirAll(list, 0o755)
	os.WriteFile(filepath.Join(list, "versions.geul"), []byte{1}, 0o644)
	if _, err := carriedIDs(resolve.Dir{Root: root}, issuer, "Q1"); err == nil {
		t.Fatal("broken version list: want error")
	}
}
