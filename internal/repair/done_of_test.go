//ff:func feature=repair type=reader control=sequence
//ff:what TestDoneOf — the occurrence of the done hashes asked for; a NULL occurrence and a missing table fail
package repair

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestDoneOf(t *testing.T) {
	db, c := loadFixture(t)
	ensureTables(db, c)
	tx, _ := db.Begin()
	defer tx.Rollback()
	markDone(tx, c, []HashCheck{{Hash: "a", At: load.Occurrence{Statement: 9, Ordinal: 1}}, {Hash: "b"}})
	got, err := doneOf(tx, c, []string{"a", "z"})
	if err != nil || len(got) != 1 || got["a"] != (load.Occurrence{Statement: 9, Ordinal: 1}) {
		t.Fatalf("%v %v", got, err)
	}
	tx.Exec("ALTER TABLE " + c.Schema + ".repair_reference_snak_done ALTER occ_ordinal DROP NOT NULL")
	tx.Exec("INSERT INTO " + c.Schema + ".repair_reference_snak_done VALUES ('n', 1, 1, NULL, 'same')")
	if _, err := doneOf(tx, c, []string{"n"}); err == nil {
		t.Fatal("NULL ordinal: want error")
	}
	tx2, _ := db.Begin()
	defer tx2.Rollback()
	c.Schema += "_none"
	if _, err := doneOf(tx2, c, []string{"a"}); err == nil {
		t.Fatal("missing table: want error")
	}
}
