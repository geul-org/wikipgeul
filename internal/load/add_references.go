//ff:func feature=load type=model control=iteration dimension=2
//ff:what Rows.addReferences — a statement's reference links in order, and each source's snaks with their occurrence when it must be written (RefSet.First)
package load

import "github.com/geul-org/wikipgeul/internal/dump"

// addReferences links statement sid to each of its references in order
// (statement_reference). A source whose occurrence refs says to write
// (RefSet.First: a new hash, or an occurrence earlier than the one written)
// gets its reference row and the snaks of this occurrence (ReferenceSnaks)
// with the occurrence (occ_statement_id = sid, occ_ordinal = the position in
// the references), so wikidata.sql 4.2 keeps each hash's first occurrence
// whole. Other occurrences are only linked.
func (r *Rows) addReferences(c Config, line int64, id string, sid int64, rs []dump.Reference, refs *RefSet) {
	for ord, ref := range rs {
		r.add("statement_reference", sid, ord, ref.Hash)
		if !refs.First(ref.Hash, Occurrence{Statement: sid, Ordinal: ord}) {
			continue
		}
		r.add("reference", ref.Hash, c.DumpID)
		for _, row := range ReferenceSnaks(c.DumpID, ref) {
			r.add("reference_snak", append(row, sid, ord)...)
		}
	}
}
