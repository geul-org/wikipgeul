//ff:func feature=convert type=loader control=sequence
//ff:what TestEntityHead — Q42 exists with no datatype and its revision, P31 with wikibase-item and no revision, a never-created ID does not, a closed handle errors
package pgsource

import "testing"

func TestEntityHead(t *testing.T) {
	db, s := fixtureDB(t)
	q := queriesFor(s).Exists
	if e, ok, err := entityHead(db, q, "Q42"); err != nil || !ok || e.ID != "Q42" || e.Datatype != "" || e.LastRevID != 2557857485 || e.Modified != 1792210885 {
		t.Fatalf("Q42: %+v %v %v", e, ok, err)
	}
	if e, ok, err := entityHead(db, q, "P31"); err != nil || !ok || e.Datatype != "wikibase-item" || e.LastRevID != 0 || e.Modified != 0 {
		t.Fatalf("P31: %+v %v %v", e, ok, err)
	}
	if _, ok, err := entityHead(db, q, "Q6"); err != nil || ok {
		t.Fatalf("Q6: %v %v", ok, err)
	}
	if _, _, err := entityHead(closedDB(t), q, "Q42"); err == nil {
		t.Fatal("closed: want error")
	}
}
