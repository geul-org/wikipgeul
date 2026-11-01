//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what TestLoadSitelinks — Q42's sitelinks of its dump in byte order of site with badges; a NUL title exact; a closed handle and a query with the wrong columns error
package pgsource

import (
	"slices"
	"testing"
)

func TestLoadSitelinks(t *testing.T) {
	db, s := fixtureDB(t)
	q := queriesFor(s).Sitelinks
	ls, err := loadSitelinks(db, q, "Q42")
	var sites []string
	for _, l := range ls {
		sites = append(sites, l.Site)
	}
	if err != nil || !slices.Equal(sites, []string{"be_x_oldwiki", "commonswiki", "enwiki", "enwikiquote", "kowiki", "xxsite", "zh_min_nanwiki"}) {
		t.Fatalf("sites %v %v", sites, err)
	}
	if ls[2].Title != "Douglas Adams" || !slices.Equal(ls[2].Badges, []string{"Q17437796"}) || len(ls[0].Badges) != 0 || !slices.Equal(ls[3].Badges, []string{"Q17437796", "P1"}) {
		t.Fatalf("enwiki %+v, be_x_oldwiki %+v, enwikiquote %+v", ls[2], ls[0], ls[3])
	}
	if q7, err := loadSitelinks(db, q, "Q7"); err != nil || len(q7) != 1 || q7[0].Title != "a\x00b" {
		t.Fatalf("Q7 (\\u0000 in title_json): %+v %v", q7, err)
	}
	if _, err := loadSitelinks(closedDB(t), q, "Q42"); err == nil {
		t.Fatal("closed: want error")
	}
	if _, err := loadSitelinks(db, "SELECT site FROM "+s.Schema+".sitelink WHERE entity_id = $1", "Q42"); err == nil {
		t.Fatal("one column: want scan error")
	}
}
