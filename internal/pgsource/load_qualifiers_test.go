//ff:func feature=convert type=loader control=sequence
//ff:what TestLoadQualifiers — Q42's qualifiers in qualifiers-order (P580, P642 as a bare ID, novalue P582 with its snak type); a closed handle and wrong columns error
package pgsource

import (
	"strings"
	"testing"
)

func TestLoadQualifiers(t *testing.T) {
	db, s := fixtureDB(t)
	q := queriesFor(s).Qualifiers
	qs, err := loadQualifiers(db, q, "Q42")
	if err != nil || len(qs) != 3 || qs[0].StatementID != 4201 || qs[0].Property != "P580" || qs[0].Type != "time" || qs[0].Datatype != "time" || !strings.Contains(qs[0].Value, "+1979-") ||
		qs[1].Value != "Q1" || qs[2].Property != "P582" || qs[2].Value != "" || qs[2].SnakType != "novalue" || qs[1].SnakType != "value" {
		t.Fatalf("got %+v %v", qs, err)
	}
	if _, err := loadQualifiers(closedDB(t), q, "Q42"); err == nil {
		t.Fatal("closed: want error")
	}
	if _, err := loadQualifiers(db, "SELECT 1, 2 FROM "+s.Schema+".statement WHERE entity_id = $1", "Q42"); err == nil {
		t.Fatal("columns: want error")
	}
}
