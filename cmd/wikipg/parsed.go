//ff:type feature=cli type=model
//ff:what parsed — a batch turned into rows by a parser, waiting for a writer
package main

import "github.com/geul-org/wikipgeul/internal/load"

// parsed is batch Seq (ending at line Last) after a parser built its rows:
// the rows and the number of entities, for a writer to commit.
type parsed struct {
	Seq      int64
	Last     int64
	Rows     *load.Rows
	Entities int64
}
