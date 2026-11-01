//ff:func feature=repair type=writer control=sequence
//ff:what TestMarkDone — records hashes with occurrence and verdict, a later call overwrites; none is a no-op; a missing table fails
package repair

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestMarkDone(t *testing.T) {
	db, c := loadFixture(t)
	ensureTables(db, c)
	tx, _ := db.Begin()
	if err := markDone(tx, c, []HashCheck{{Hash: "a", At: load.Occurrence{Statement: 9, Ordinal: 1}, Verdict: Order}, {Hash: "b", Verdict: Same}}); err != nil {
		t.Fatal(err)
	}
	if err := markDone(tx, c, []HashCheck{{Hash: "a", At: load.Occurrence{Statement: 7, Ordinal: 2}, Verdict: Content}}); err != nil {
		t.Fatal(err)
	}
	if err := markDone(tx, c, nil); err != nil {
		t.Fatal(err)
	}
	tx.Commit()
	var sid, ord int64
	var v string
	db.QueryRow("SELECT occ_statement_id, occ_ordinal, verdict FROM "+c.Schema+".repair_reference_snak_done WHERE hash = 'a'").Scan(&sid, &ord, &v)
	if sid != 7 || ord != 2 || v != "content" {
		t.Fatalf("%d %d %s", sid, ord, v)
	}
	c.Schema += "_none"
	tx, _ = db.Begin()
	defer tx.Rollback()
	if err := markDone(tx, c, []HashCheck{{Hash: "a"}}); err == nil {
		t.Fatal("missing table: want error")
	}
}
