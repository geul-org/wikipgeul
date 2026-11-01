//ff:func feature=cli type=command control=sequence
//ff:what takeLine — adds a dump line to the batch when it is an entity line at or after the start; 1 when taken
package main

import "github.com/geul-org/wikipgeul/internal/dump"

// takeLine adds raw, dump line number line, to b when it holds an entity
// (dump.TrimLine) and line is at or after from, and returns 1 then, else 0.
func takeLine(b *batch, raw []byte, line, from int64) int64 {
	e, ok := dump.TrimLine(raw)
	if !ok || line < from {
		return 0
	}
	b.Lines, b.Nums, b.Last, b.Bytes = append(b.Lines, e), append(b.Nums, line), line, b.Bytes+int64(len(e))
	return 1
}
