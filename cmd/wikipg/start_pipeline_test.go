//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestStartPipeline — batches flow through parsers and writers into one result each, and the result channel closes
package main

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestStartPipeline(t *testing.T) {
	db, schema := testSchema(t)
	batches := make(chan batch, 4)
	for i := int64(0); i < 3; i++ {
		batches <- batch{Seq: i, Last: i + 2, Lines: [][]byte{[]byte(`{"type":"item","id":"Q` + string(rune('1'+i)) + `"}`)}, Nums: []int64{i + 2}}
	}
	close(batches)
	c := loadConfig{Schema: schema, Parsers: 2, Writers: 2}
	results := startPipeline(db, c, load.Config{DumpID: 1}, load.NewRefSet(64), batches)
	n := 0
	for r := range results {
		if r.Err != nil {
			t.Fatal(r.Err)
		}
		n++
	}
	var rows int
	db.QueryRow("SELECT count(*) FROM " + schema + ".entity").Scan(&rows)
	if n != 3 || rows != 3 {
		t.Fatalf("results %d rows %d", n, rows)
	}
}
