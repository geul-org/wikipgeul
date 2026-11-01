//ff:func feature=convert type=loader control=sequence
//ff:what TestEntityNumbers — numbers of one kind in the dump, ascending, after and up to bounds, limited; other dumps and kinds excluded; a bad schema and a closed database fail
package pgsource

import (
	"slices"
	"testing"
)

func TestEntityNumbers(t *testing.T) {
	db, s := fixtureDB(t)
	if got, err := EntityNumbers(db, s, 'Q', 0, 1000, 100); err != nil || !slices.Equal(got, []uint64{1, 2, 3, 7, 30, 42}) {
		t.Fatalf("Q: %v %v", got, err)
	}
	if got, _ := EntityNumbers(db, s, 'Q', 2, 30, 2); !slices.Equal(got, []uint64{3, 7}) {
		t.Fatalf("after 2, limit 2: %v", got)
	}
	if got, _ := EntityNumbers(db, s, 'P', 0, ^uint64(0), 10); !slices.Equal(got, []uint64{31, 580}) {
		t.Fatalf("P: %v", got)
	}
	if _, err := EntityNumbers(db, Store{Schema: "1bad", Dump: 1}, 'Q', 0, 10, 10); err == nil {
		t.Fatal("bad schema: want error")
	}
	if _, err := EntityNumbers(db, Store{Schema: "no_such_schema", Dump: 1}, 'Q', 0, 10, 10); err == nil {
		t.Fatal("missing schema: want error")
	}
	bad := s.Schema + "_n"
	if _, err := db.Exec("CREATE SCHEMA " + bad + "; CREATE TABLE " + bad + ".entity (dump_id int, kind text, num numeric); INSERT INTO " + bad + ".entity VALUES (1, 'Q', 1.5)"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Exec("DROP SCHEMA " + bad + " CASCADE") })
	if _, err := EntityNumbers(db, Store{Schema: bad, Dump: 1}, 'Q', 0, 10, 10); err == nil {
		t.Fatal("a number that is no integer: want error")
	}
}
