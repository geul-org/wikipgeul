//ff:func feature=repair type=helper control=iteration dimension=1
//ff:what compareStatementRefs — statement_reference rows that differ between the dump and the database, and the lines they are on
package repair

import "github.com/geul-org/wikipgeul/internal/load"

// compareStatementRefs counts the rows of want (from the dump) and have
// (stored) that do not match: a key in one only, or the same key with
// another hash. It returns that count and the dump lines of those statements
// (key >> 20 within the dump's 35 line bits).
func compareStatementRefs(want, have []stmtRef) (int64, map[int64]bool) {
	stored := map[load.Occurrence]string{}
	for _, r := range have {
		stored[load.Occurrence{Statement: r.Statement, Ordinal: r.Ordinal}] = r.Hash
	}
	var diff int64
	lines := map[int64]bool{}
	for _, r := range want {
		k := load.Occurrence{Statement: r.Statement, Ordinal: r.Ordinal}
		if h, ok := stored[k]; !ok || h != r.Hash {
			diff++
			lines[r.Statement>>20&(1<<35-1)] = true
		}
		delete(stored, k)
	}
	for k := range stored {
		diff++
		lines[k.Statement>>20&(1<<35-1)] = true
	}
	return diff, lines
}
