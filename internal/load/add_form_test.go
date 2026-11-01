//ff:func feature=load type=model control=sequence
//ff:what TestRowsAddForm — entity row of type form, form row with features (empty array when none), representations and statements; key range ends with false
package load

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/dump"
)

func TestRowsAddForm(t *testing.T) {
	r := NewRows()
	f := sampleLexeme().Forms[0]
	next, ok := r.addForm(Config{DumpID: 1}, 3, "L4", 4, 0, f, 5, NewRefSet(64))
	if !ok || next != 6 || r.Count("entity") != 1 || r.Count("form") != 1 || r.Count("form_representation") != 1 || r.Count("statement") != 1 {
		t.Fatalf("got %d %v %+v", next, ok, r.byTable)
	}
	if _, ok := r.addForm(Config{DumpID: 1}, 3, "L4", 4, 1, dump.Form{ID: "L4-F2"}, 6, NewRefSet(64)); !ok || r.Count("form") != 2 {
		t.Fatal("no features: want a row")
	}
	if _, ok := NewRows().addForm(Config{DumpID: 0}, 3, "L4", 4, 0, f, 0, NewRefSet(64)); ok {
		t.Fatal("bad dump: want false")
	}
}
