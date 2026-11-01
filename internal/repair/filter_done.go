//ff:func feature=repair type=helper control=iteration dimension=1
//ff:what filterDone — the checks to process: those not done, or done only at a later occurrence; and how many were skipped
package repair

import "github.com/geul-org/wikipgeul/internal/load"

// filterDone returns the checks whose hash is not in done or was done at a
// later occurrence than the check's (an earlier range run after a later
// one), and the number skipped because done at the same or an earlier one.
func filterDone(checks []HashCheck, done map[string]load.Occurrence) ([]HashCheck, int64) {
	var out []HashCheck
	var skipped int64
	for _, h := range checks {
		if d, ok := done[h.Hash]; ok && !h.At.Before(d) {
			skipped++
			continue
		}
		out = append(out, h)
	}
	return out, skipped
}
