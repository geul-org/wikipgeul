//ff:func feature=cli type=command control=iteration dimension=1
//ff:what parser — turns batches into rows (CPU work) and hands them to the writers
package main

import "github.com/geul-org/wikipgeul/internal/load"

// parser builds the rows of every batch from in (buildRows: JSON decoding and
// row building, the CPU-heavy part) and sends them to out for a writer.
func parser(c load.Config, refs *load.RefSet, in <-chan batch, out chan<- parsed) {
	for b := range in {
		rows, n := buildRows(c, b, refs)
		out <- parsed{Seq: b.Seq, Last: b.Last, Rows: rows, Entities: n}
	}
}
