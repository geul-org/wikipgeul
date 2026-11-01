//ff:type feature=repair type=model
//ff:what HashCheck — a hash's first occurrence in a batch with its verdict against the stored rows and the rows to store
package repair

import "github.com/geul-org/wikipgeul/internal/load"

// HashCheck is the first occurrence of a hash in a batch: the hash, where it
// occurs, the entity, how the stored rows compare with it and its rows
// (load.ReferenceSnaks), which replace the stored ones when they differ.
type HashCheck struct {
	Hash    string
	At      load.Occurrence
	Entity  string
	Verdict Verdict
	Rows    [][]any
}
