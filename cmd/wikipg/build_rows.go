//ff:func feature=cli type=command control=iteration dimension=1
//ff:what buildRows — the rows of one batch: each line decoded and added, or recorded as a parse error
package main

import (
	"github.com/geul-org/wikipgeul/internal/dump"
	"github.com/geul-org/wikipgeul/internal/load"
)

// buildRows decodes each line of b (dump.DecodeEntity) and adds its rows
// (load.Rows.AddEntity). A line that does not decode is a load_error of stage
// parse and the batch goes on. It returns the rows and the number of entities.
func buildRows(c load.Config, b batch, refs *load.RefSet) (*load.Rows, int64) {
	r := load.NewRows()
	var n int64
	for i, raw := range b.Lines {
		e, err := dump.DecodeEntity(raw)
		if err != nil {
			r.Error(c.DumpID, b.Nums[i], "", "parse", err.Error(), raw)
			continue
		}
		r.AddEntity(c, b.Nums[i], raw, e, refs)
		n++
	}
	return r, n
}
