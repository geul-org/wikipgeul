//ff:type feature=repair type=model
//ff:what occ — one occurrence of a reference on a dump line: hash, place, entity and its reference_snak rows as the loader builds them
package repair

import "github.com/geul-org/wikipgeul/internal/load"

// occ is one occurrence of a reference read from entity_raw: its hash, where
// it occurs (statement key, position in the references), the entity whose
// line it is on and the reference_snak rows of this occurrence (hash, dump_id,
// ordinal, property, value_ordinal and the snak columns: load.ReferenceSnaks).
type occ struct {
	Hash   string
	At     load.Occurrence
	Entity string
	Rows   [][]any
}
