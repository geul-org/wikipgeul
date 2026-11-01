//ff:func feature=cli type=command control=iteration dimension=1
//ff:what readBatches — reads dump lines from a line on, up to a limit, into batches until the input ends or the context stops
package main

import (
	"bufio"
	"context"
	"errors"
	"io"
)

// readBatches reads r line by line (line 1 is the dump's "["), skips lines
// before from, and sends batches of entity lines (takeLine) to out, each up
// to size lines or maxBytes bytes (whichever comes first: entity sizes vary
// from bytes to megabytes, so bytes bound the memory and the write chunk),
// numbered from 0, until the input ends, limit entity lines were taken (0: no
// limit) or ctx is done. A batch is also sent at the limit and at the
// end or a read error. It closes out and returns the last line read. Lines of
// any length are read whole (an entity line can be megabytes).
func readBatches(ctx context.Context, r io.Reader, from, limit int64, size int, maxBytes int64, out chan<- batch) (int64, error) {
	defer close(out)
	br := bufio.NewReaderSize(r, 1<<20)
	var line, taken int64
	b := batch{}
	for ctx.Err() == nil && (limit == 0 || taken < limit) {
		raw, err := br.ReadBytes('\n')
		if len(raw) > 0 {
			line++
			taken += takeLine(&b, raw, line, from)
		}
		atLimit := limit > 0 && taken == limit
		full := len(b.Lines) == size || b.Bytes >= maxBytes
		if len(b.Lines) > 0 && (full || err != nil || atLimit) {
			out <- b
			b = batch{Seq: b.Seq + 1}
		}
		if errors.Is(err, io.EOF) {
			return line, nil
		}
		if err != nil {
			return line, err
		}
	}
	return line, ctx.Err()
}
