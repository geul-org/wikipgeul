//ff:type feature=load type=model
//ff:what snakAt — a qualifier or reference snak with its place: overall ordinal and ordinal within its property
package load

import "github.com/geul-org/wikipgeul/internal/dump"

// snakAt is one snak of a qualifier or reference set with its overall
// position (Ordinal) and its position among snaks of the same property
// (ValueOrdinal), both 0-based.
type snakAt struct {
	Snak         dump.Snak
	Ordinal      int
	ValueOrdinal int
}
