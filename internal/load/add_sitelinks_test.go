//ff:func feature=load type=model control=sequence
//ff:what TestRowsAddSitelinks — one row per site with badges; no badges is an empty array; a NUL title is cleaned and kept exact in title_json
package load

import (
	"database/sql/driver"
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
)

func TestRowsAddSitelinks(t *testing.T) {
	r := NewRows()
	e := dump.Entity{ID: "Q1", Sitelinks: map[string]dump.Sitelink{"enwiki": {Site: "enwiki", Title: "A\x00"}}}
	r.addSitelinks(Config{DumpID: 1}, 3, e)
	row := r.byTable["sitelink"][0]
	if r.Count("sitelink") != 1 || row[3] != "A" || row[5] != `"A\u0000"` || r.Errors != 0 {
		t.Fatalf("got %v", row)
	}
	if v, err := row[4].(driver.Valuer).Value(); err != nil || v != "{}" {
		t.Fatalf("badges %v %v", v, err)
	}
}
