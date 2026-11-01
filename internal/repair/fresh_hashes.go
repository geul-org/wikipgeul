//ff:func feature=repair type=helper control=iteration dimension=1
//ff:what freshHashes — the hash checks of a batch not already handled in this run at an earlier occurrence (and remembers them)
package repair

import "github.com/geul-org/wikipgeul/internal/load"

// freshHashes returns the checks of f whose hash this run has not handled at
// an earlier or the same occurrence (load.RefSet.First), remembering them.
// Batches are applied in line order, so the first time a hash comes is its
// first occurrence in the run's range.
func freshHashes(f Found, refs *load.RefSet) []HashCheck {
	var out []HashCheck
	for _, h := range f.Hashes {
		if refs.First(h.Hash, h.At) {
			out = append(out, h)
		}
	}
	return out
}
