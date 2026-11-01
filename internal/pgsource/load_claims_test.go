//ff:func feature=convert type=loader control=sequence
//ff:what TestLoadClaims — Q42's statements in dump order with bare entity IDs, unquoted strings and empty novalue with its snak type; a closed handle and wrong columns error
package pgsource

import "testing"

func TestLoadClaims(t *testing.T) {
	db, s := fixtureDB(t)
	q := queriesFor(s).Claims
	cs, err := loadClaims(db, q, "Q42")
	if err != nil || len(cs) != 8 || cs[6].Datatype != "url" || cs[7].Datatype != "commonsMedia" || cs[0].ID != 4201 || cs[0].Value != "Q5" || cs[0].Type != "wikibase-entityid" || cs[0].Rank != "preferred" ||
		cs[3].Value != "113230702" || cs[4].Value != "" || cs[4].Type != "" || cs[4].SnakType != "novalue" || cs[0].SnakType != "value" || cs[1].Value[0] != '{' {
		t.Fatalf("got %+v %v", cs, err)
	}
	if _, err := loadClaims(closedDB(t), q, "Q42"); err == nil {
		t.Fatal("closed: want error")
	}
	if _, err := loadClaims(db, "SELECT statement_id FROM "+s.Schema+".statement WHERE entity_id = $1", "Q42"); err == nil {
		t.Fatal("one column: want scan error")
	}
}
