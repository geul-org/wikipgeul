//ff:func feature=repair type=helper control=sequence
//ff:what TestStatsAddOrder — an order batch adds the batch counts and its orders, written only when stored, with examples up to 20 and no hash verdicts
package repair

import (
	"strings"
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestStatsAddOrder(t *testing.T) {
	var st Stats
	f := Found{Last: 4, Entities: 3, Occurrences: Tally{Same: 5, Order: 2}, Orders: []OccOrder{{Hash: "hx", Entity: "Q2", Perm: []int{2, 3, 1, 0}, At: load.Occurrence{Statement: 7}}, {Hash: "hy", Entity: "Q3", Perm: []int{1, 0}}}}
	st.addOrder(f, false)
	st.addOrder(f, true)
	if st.Spans != 2 || st.Last != 4 || st.Entities != 6 || st.Occurrences != (Tally{Same: 10, Order: 4}) || st.Orders != 4 || st.Written != 2 || st.Hashes != (Tally{}) || st.Replaced != 0 {
		t.Fatalf("%+v", st)
	}
	if len(st.Examples) != 4 || st.Examples[0] != "hx Q2 order [2 3 1 0] (statement 7 reference 0)" {
		t.Fatalf("examples %q", st.Examples)
	}
	full := Stats{Examples: make([]string, 19)}
	full.addOrder(f, true)
	if len(full.Examples) != 20 || !strings.HasPrefix(full.Examples[19], "hx") {
		t.Fatalf("bound %q", full.Examples)
	}
}
