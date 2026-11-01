//ff:func feature=convert type=loader control=sequence
//ff:what TestLoadTexts — Q42's labels in every language (byte order of the code); a closed handle and a query with the wrong columns error
package pgsource

import (
	"testing"
)

func TestLoadTexts(t *testing.T) {
	db, s := fixtureDB(t)
	q := queriesFor(s).Labels
	ts, err := loadTexts(db, q, "Q42")
	if err != nil || len(ts) != 3 || ts[0].Lang != "en" || ts[0].Value != "Douglas Adams" || ts[1].Lang != "ko" || ts[2].Lang != "mul" {
		t.Fatalf("got %+v %v", ts, err)
	}
	if _, err := loadTexts(closedDB(t), q, "Q42"); err == nil {
		t.Fatal("closed: want error")
	}
	bad := "SELECT language, value FROM " + s.Schema + ".label WHERE entity_id = $1"
	if _, err := loadTexts(db, bad, "Q42"); err == nil {
		t.Fatal("two columns: want scan error")
	}
}
