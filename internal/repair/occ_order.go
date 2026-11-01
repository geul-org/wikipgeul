//ff:type feature=repair type=model
//ff:what OccOrder — a reference occurrence whose snaks are those stored for its hash in another order: where it is and the permutation of the stored ordinals
package repair

import "github.com/geul-org/wikipgeul/internal/load"

// OccOrder is one reference occurrence whose snak list has the snaks of the
// stored reference_snak rows of its hash in another order (verdict Order):
// its hash, where it occurs (the key of reference_snak_order), the entity and
// Perm, the stored ordinal of each of its snaks in its own order.
type OccOrder struct {
	Hash   string
	At     load.Occurrence
	Entity string
	Perm   []int
}
