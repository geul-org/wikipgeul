//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what writeReport — the end report of a repair or order build run: range, totals, hash and occurrence verdicts, occurrence orders, statement_reference check, errors and examples
package main

import (
	"fmt"
	"io"

	"github.com/geul-org/wikipgeul/internal/repair"
)

// writeReport writes the totals st of run c to w: the range and mode, the
// verdict of each hash at its first occurrence in the range (to replace in a
// dry run, replaced when writing), every occurrence compared with the rows
// read and the entities affected, for an order build (c.Order) the
// occurrences in another order found and written, the statement_reference
// check, then the parse errors and examples (up to 20 each).
func writeReport(w io.Writer, c repair.Config, st repair.Stats) {
	mode, name := "write", "repair-reference-snak"
	if c.DryRun {
		mode = "dry run"
	}
	if c.Order {
		name = "build-reference-order"
	}
	fmt.Fprintf(w, "%s %s (%s), dump %d: %d/%d batches through line %d, %d entities, %.1f MB raw, %s\n",
		name, c.Run(), mode, c.DumpID, st.Spans, st.Planned, st.Last, st.Entities, float64(st.Bytes)/1e6, st.Elapsed.Round(1e9))
	fmt.Fprintf(w, "hashes (first occurrence in range): same %d, order %d, content %d, missing %d; replaced %d, skipped as done %d\n",
		st.Hashes.Same, st.Hashes.Order, st.Hashes.Content, st.Hashes.Missing, st.Replaced, st.Skipped)
	fmt.Fprintf(w, "occurrences vs stored rows: same %d, order %d (%d entities), content %d (%d entities), missing %d\n",
		st.Occurrences.Same, st.Occurrences.Order, st.OrderEntities, st.Occurrences.Content, st.ContentEntities, st.Occurrences.Missing)
	if c.Order {
		fmt.Fprintf(w, "occurrence orders (reference_snak_order): found %d, written %d\n", st.Orders, st.Written)
	}
	fmt.Fprintf(w, "statement_reference: %d rows from the dump, %d differ %v\n", st.StatementRefs, st.RefDiffs, st.RefDiffEntities)
	for _, e := range append(st.Errors, st.Examples...) {
		fmt.Fprintln(w, "  ", e)
	}
}
