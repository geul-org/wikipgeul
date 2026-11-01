//ff:func feature=load type=helper control=iteration dimension=2
//ff:what orderedSnaks — the snaks of a property map laid out in the dump's property order, then any properties it misses
package load

import (
	"slices"
	"sort"

	"github.com/geul-org/wikipgeul/internal/dump"
)

// orderedSnaks lays out the snaks of m (by property) as the dump orders
// them: properties in order (qualifiers-order, snaks-order), each property's
// snaks in their array order. Properties of m missing from order follow in
// sorted order, so nothing is lost when the order list is incomplete.
func orderedSnaks(order []string, m map[string][]dump.Snak) []snakAt {
	props := slices.Clone(order)
	var rest []string
	for p := range m {
		if !slices.Contains(order, p) {
			rest = append(rest, p)
		}
	}
	sort.Strings(rest)
	var out []snakAt
	for _, p := range append(props, rest...) {
		for i, s := range m[p] {
			out = append(out, snakAt{Snak: s, Ordinal: len(out), ValueOrdinal: i})
		}
	}
	return out
}
