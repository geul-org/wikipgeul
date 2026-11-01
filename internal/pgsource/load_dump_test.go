//ff:func feature=convert type=loader control=sequence
//ff:what TestLoadDump — the fixture dump's URL, download time, loader and commit, and whether reference_snak_order is there; a missing dump and a bad schema error
package pgsource

import (
	"strings"
	"testing"
)

func TestLoadDump(t *testing.T) {
	db, s := fixtureDB(t)
	d, err := LoadDump(db, s)
	if err != nil || !strings.HasSuffix(d.URL, "wikidata-20261019-all.json.bz2") || !strings.HasPrefix(d.SHA256, "a249a19d") || d.DownloadedAt.Unix() != 1792996483 || d.Loader != "github.com/geul-org/wikipgeul/cmd/wikipg" || len(d.LoaderCommit) != 40 || d.RefOrder {
		t.Fatalf("got %+v %v", d, err)
	}
	db.Exec("CREATE TABLE " + s.Schema + ".reference_snak_order (statement_id bigint, ordinal smallint, perm smallint[], PRIMARY KEY (statement_id, ordinal))")
	if d, err := LoadDump(db, s); err != nil || !d.RefOrder {
		t.Fatalf("with reference_snak_order: %+v %v", d, err)
	}
	if _, err := LoadDump(db, Store{Schema: s.Schema, Dump: 9}); err == nil {
		t.Fatal("missing dump: want error")
	}
	if _, err := LoadDump(db, Store{Schema: "Bad", Dump: 1}); err == nil {
		t.Fatal("bad schema: want error")
	}
}
