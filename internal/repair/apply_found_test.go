//ff:func feature=repair type=writer control=sequence
//ff:what TestApplyFound — a dry run counts without writing, a write repairs and counts replacements, a failed write is returned; an order build counts orders and leaves reference_snak alone
package repair

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestApplyFound(t *testing.T) {
	db, c := loadFixture(t)
	f, err := checkSpan(db, c, Span{From: 2, To: 3, Last: 3})
	if err != nil {
		t.Fatal(err)
	}
	var st Stats
	if err := applyFound(db, c, f, load.NewRefSet(64), &st); err != nil || st.Hashes.Content != 1 || st.Replaced != 0 || snakList(t, db, c, "hz") != "" {
		t.Fatalf("dry: %+v %v", st, err)
	}
	c.DryRun = false
	if err := applyFound(db, c, f, load.NewRefSet(64), &st); err == nil {
		t.Fatal("no work tables: want error")
	}
	ensureTables(db, c)
	if err := applyFound(db, c, f, load.NewRefSet(64), &st); err != nil || st.Replaced != 3 || snakList(t, db, c, "hz") == "" {
		t.Fatalf("write: %+v %v", st, err)
	}
	c.Order, c.DryRun = true, true
	var ost Stats
	if err := applyFound(db, c, f, load.NewRefSet(64), &ost); err != nil || ost.Orders != 1 || ost.Replaced != 0 || ost.Hashes != (Tally{}) || ost.Spans != 1 {
		t.Fatalf("order: %+v %v", ost, err)
	}
}
