//ff:func feature=cli type=command control=sequence
//ff:what TestParser — every batch comes out parsed with its number, last line and entity count
package main

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestParser(t *testing.T) {
	in, out := make(chan batch, 2), make(chan parsed, 2)
	in <- batch{Seq: 4, Last: 9, Lines: [][]byte{[]byte(`{"id":"Q1"}`)}, Nums: []int64{9}}
	close(in)
	parser(load.Config{DumpID: 1}, load.NewRefSet(64), in, out)
	p := <-out
	if p.Seq != 4 || p.Last != 9 || p.Entities != 1 || p.Rows.Count("entity") != 1 {
		t.Fatalf("got %+v", p)
	}
}
