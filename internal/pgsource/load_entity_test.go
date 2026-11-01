//ff:func feature=convert type=loader control=sequence
//ff:what TestLoadEntity — Q42 with labels, descriptions, aliases (every language), claims and sitelinks of its dump only, P31 with its datatype; Q30 with novalue/somevalue snak types; a missing ID is false; a bad schema and a closed handle error
package pgsource

import "testing"

func TestLoadEntity(t *testing.T) {
	db, s := fixtureDB(t)
	e, ok, err := LoadEntity(db, s, "Q42")
	if err != nil || !ok || e.ID != "Q42" || len(e.Labels) != 3 || len(e.Sitelinks) != 7 || len(e.Descriptions) != 2 || len(e.Aliases) != 2 || len(e.Claims) != 8 || len(e.Claims[0].Qualifiers) != 3 || len(e.Claims[0].References) != 2 || len(e.Claims[1].References) != 0 {
		t.Fatalf("Q42: %v %v %+v", ok, err, e)
	}
	if p, ok, err := LoadEntity(db, s, "P31"); err != nil || !ok || p.Datatype != "wikibase-item" {
		t.Fatalf("P31: %+v %v %v", p.Datatype, ok, err)
	}
	if q7, ok, err := LoadEntity(db, s, "Q7"); err != nil || !ok || len(q7.Labels) != 1 || q7.Labels[0].Value != "a\x00b" || len(q7.Sitelinks) != 1 || q7.Sitelinks[0].Title != "a\x00b" ||
		len(q7.Claims) != 1 || q7.Claims[0].Value != "x\x00" || len(q7.Claims[0].Qualifiers) != 1 || q7.Claims[0].Qualifiers[0].Value != "q\x00" {
		t.Fatalf("Q7 (\\u0000 in value_json): %+v %v %v", q7, ok, err)
	}
	if q30, ok, err := LoadEntity(db, s, "Q30"); err != nil || !ok || len(q30.Claims) != 3 || q30.Claims[0].SnakType != "novalue" || len(q30.Claims[0].Qualifiers) != 2 ||
		q30.Claims[1].SnakType != "somevalue" || len(q30.Claims[2].Qualifiers) != 1 || q30.Claims[2].Qualifiers[0].SnakType != "somevalue" || q30.Claims[2].Qualifiers[0].Value != "" {
		t.Fatalf("Q30 (special snaks): %+v %v %v", q30, ok, err)
	}
	if _, ok, err := LoadEntity(db, s, "Q6"); err != nil || ok {
		t.Fatalf("Q6: %v %v", ok, err)
	}
	if _, ok, err := LoadEntity(db, s, "Q99"); err != nil || ok {
		t.Fatalf("Q99 of another dump: %v %v", ok, err)
	}
	if _, _, err := LoadEntity(db, Store{Schema: "x;drop", Dump: 1}, "Q42"); err == nil {
		t.Fatal("bad schema: want error")
	}
	if _, _, err := LoadEntity(closedDB(t), s, "Q42"); err == nil {
		t.Fatal("closed: want error")
	}
}
