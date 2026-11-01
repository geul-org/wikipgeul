//ff:func feature=load type=model control=iteration dimension=1
//ff:what Rows.addStatement — one statement row, its qualifiers in dump order and its references
package load

import "github.com/geul-org/wikipgeul/internal/dump"

// addStatement adds statement s of entity id as key sid at ordinal ord: the
// statement row (main snak columns by snakValues), a qualifier row per
// qualifier snak in dump order (orderedSnaks) and its references
// (addReferences).
func (r *Rows) addStatement(c Config, line int64, id string, sid int64, ord int, s dump.Statement, refs *RefSet) {
	main := snakValues(s.Mainsnak)
	r.add("statement", append([]any{sid, s.ID, id, s.Mainsnak.Property, ord, s.Rank}, append(main, c.DumpID)...)...)
	for _, q := range orderedSnaks(s.QualifiersOrder, s.Qualifiers) {
		vals := snakValues(q.Snak)
		r.add("qualifier", append([]any{sid, q.Ordinal, q.Snak.Property, q.ValueOrdinal}, append(vals, nullable(q.Snak.Hash))...)...)
	}
	r.addReferences(c, line, id, sid, s.References, refs)
}
