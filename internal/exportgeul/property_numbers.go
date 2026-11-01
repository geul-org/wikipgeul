//ff:func feature=convert type=helper control=iteration dimension=2
//ff:what propertyNumbers — the distinct property numbers the reports used, ascending
package exportgeul

import (
	"slices"
	"strconv"

	"github.com/geul-org/wikigeul/pkg/port"
)

// propertyNumbers collects the numbers of the properties (P…) the reports
// used as predicates or values, without duplicates, ascending.
func propertyNumbers(reps []port.Report) []uint64 {
	var out []uint64
	for _, r := range reps {
		for _, p := range r.Properties {
			if n, err := strconv.ParseUint(p[1:], 10, 64); err == nil {
				out = append(out, n)
			}
		}
	}
	slices.Sort(out)
	return slices.Compact(out)
}
