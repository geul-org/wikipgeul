//ff:func feature=cli type=command control=sequence
//ff:what TestBuildRows — decodable lines become entities, a broken line a parse error
package main

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestBuildRows(t *testing.T) {
	b := batch{Lines: [][]byte{[]byte(`{"type":"item","id":"Q1"}`), []byte(`{"id":`)}, Nums: []int64{2, 3}}
	r, n := buildRows(load.Config{DumpID: 1, Raw: true}, b, load.NewRefSet(64))
	if n != 1 || r.Errors != 1 || r.Count("entity") != 1 || r.Count("load_error") != 1 {
		t.Fatalf("n %d errors %d", n, r.Errors)
	}
}
