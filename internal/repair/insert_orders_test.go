//ff:func feature=repair type=writer control=sequence
//ff:what TestInsertOrders — stores permutations, keeps an equal row, replaces a different one; no orders is a no-op; a missing permutation or table fails
package repair

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestInsertOrders(t *testing.T) {
	db, schema := testSchema(t)
	c := Config{Schema: schema, Order: true}
	if err := ensureTables(db, c); err != nil {
		t.Fatal(err)
	}
	at := func(line int64, i int64, ord int) load.Occurrence {
		return load.Occurrence{Statement: 1<<55 | line<<20 | i, Ordinal: ord}
	}
	write := func(orders []OccOrder) error {
		tx, err := db.Begin()
		if err != nil {
			t.Fatal(err)
		}
		defer tx.Rollback()
		if err := insertOrders(tx, c, orders); err != nil {
			return err
		}
		return tx.Commit()
	}
	if err := write([]OccOrder{{At: at(3, 1, 0), Perm: []int{2, 3, 1, 0}}, {At: at(4, 0, 1), Perm: []int{1, 0}}}); err != nil {
		t.Fatal(err)
	}
	if got := orderList(t, db, c); got != "3.1/0={2,3,1,0} 4.0/1={1,0}" {
		t.Fatalf("stored %q", got)
	}
	var xmin1, xmin2 string
	db.QueryRow("SELECT xmin::text FROM " + schema + ".reference_snak_order WHERE ordinal = 0").Scan(&xmin1)
	if err := write([]OccOrder{{At: at(3, 1, 0), Perm: []int{2, 3, 1, 0}}, {At: at(4, 0, 1), Perm: []int{0, 1}}}); err != nil {
		t.Fatal(err)
	}
	db.QueryRow("SELECT xmin::text FROM " + schema + ".reference_snak_order WHERE ordinal = 0").Scan(&xmin2)
	if got := orderList(t, db, c); got != "3.1/0={2,3,1,0} 4.0/1={0,1}" || xmin1 != xmin2 || xmin1 == "" {
		t.Fatalf("again %q (row versions %s %s)", got, xmin1, xmin2)
	}
	if err := write(nil); err != nil {
		t.Fatal("none:", err)
	}
	if err := write([]OccOrder{{Hash: "h", At: at(5, 0, 0)}}); err == nil {
		t.Fatal("no permutation: want error")
	}
	c.Schema += "_none"
	if err := write([]OccOrder{{At: at(5, 0, 0), Perm: []int{1, 0}}}); err == nil {
		t.Fatal("missing table: want error")
	}
}
