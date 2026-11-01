//ff:func feature=load type=model control=iteration dimension=1
//ff:what ReferenceSnaks — the reference_snak rows of one occurrence of a source: its snaks in the order this occurrence gives (snaks-order)
package load

import "github.com/geul-org/wikipgeul/internal/dump"

// ReferenceSnaks returns the reference_snak rows of ref as it occurs in one
// statement: hash, dump_id, ordinal, property, value_ordinal and the snak
// columns (snakValues), the snaks laid out by this occurrence's snaks-order
// (orderedSnaks). Another occurrence of the same hash may order them
// differently; which occurrence a hash keeps is decided by Occurrence.
func ReferenceSnaks(dumpID int16, ref dump.Reference) [][]any {
	var rows [][]any
	for _, sn := range orderedSnaks(ref.SnaksOrder, ref.Snaks) {
		rows = append(rows, append([]any{ref.Hash, dumpID, sn.Ordinal, sn.Snak.Property, sn.ValueOrdinal}, snakValues(sn.Snak)...))
	}
	return rows
}
