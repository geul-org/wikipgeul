//ff:func feature=repair type=writer control=iteration dimension=1
//ff:what TestApplyWrite — replaces differing hashes, records all and the checkpoint in one transaction, skips done ones; each failing step rolls the batch back
package repair

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestApplyWrite(t *testing.T) {
	db, c := loadFixture(t)
	c.DryRun = false
	ensureTables(db, c)
	f, err := checkSpan(db, c, Span{From: 2, To: 2, Last: 2})
	if err != nil {
		t.Fatal(err)
	}
	processed, skipped, err := applyWrite(db, c, 2, f.Hashes)
	if err != nil || len(processed) != 2 || skipped != 0 {
		t.Fatalf("%v %d %v", processed, skipped, err)
	}
	if got := snakList(t, db, c, "hy"); got != "P854=https://example.org/a/0 P248=Q1/0" {
		t.Fatalf("hy %q", got)
	}
	if n, _ := readCheckpoint(db, c); n != 2 {
		t.Fatalf("checkpoint %d", n)
	}
	if _, skipped, err := applyWrite(db, c, 2, f.Hashes); err != nil || skipped != 2 {
		t.Fatalf("again: %d %v", skipped, err)
	}
	late := []HashCheck{{Hash: "hq", At: load.Occurrence{Statement: 1}, Verdict: Missing, Rows: [][]any{{"hq", int16(1), 0, "P1", 0, "novalue", nil, nil, nil, nil, nil}}}}
	steps := []string{
		"ALTER TABLE " + c.Schema + ".repair_reference_snak_checkpoint RENAME TO cp_gone",
		"ALTER TABLE " + c.Schema + ".repair_reference_snak_done ADD CHECK (verdict <> 'missing')",
		"ALTER TABLE " + c.Schema + ".reference_snak RENAME TO rs_gone",
		"ALTER TABLE " + c.Schema + ".repair_reference_snak_done RENAME TO done_gone",
	}
	for _, q := range steps {
		if _, err := db.Exec(q); err != nil {
			t.Fatal(q, err)
		}
		if _, _, err := applyWrite(db, c, 3, late); err == nil {
			t.Fatalf("after %s: want error", q)
		}
	}
	if _, _, err := applyWrite(closedDB(t), c, 3, late); err == nil {
		t.Fatal("closed: want error")
	}
}
