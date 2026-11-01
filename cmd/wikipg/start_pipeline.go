//ff:func feature=cli type=command control=iteration dimension=1
//ff:what startPipeline — starts the parsers and writers between the batch channel and the result channel, each stage closing the next when done
package main

import (
	"database/sql"
	"sync"

	"github.com/geul-org/wikipgeul/internal/load"
)

// startPipeline starts c.Parsers parsers on batches and c.Writers writers on
// their output, and returns the results channel. The parsed channel closes
// once every parser is done, the results channel once every writer is. The
// queues hold one batch per parser and per writer, so at most about
// 2 × (parsers + writers) batches are in memory at once.
func startPipeline(db *sql.DB, c loadConfig, lc load.Config, refs *load.RefSet, batches <-chan batch) <-chan result {
	rows := make(chan parsed, c.Writers)
	results := make(chan result, c.Writers*2)
	var ps, ws sync.WaitGroup
	for i := 0; i < c.Parsers; i++ {
		ps.Add(1)
		go func() { defer ps.Done(); parser(lc, refs, batches, rows) }()
	}
	for i := 0; i < c.Writers; i++ {
		ws.Add(1)
		go func() { defer ws.Done(); writer(db, c.Schema, rows, results) }()
	}
	go func() { ps.Wait(); close(rows) }()
	go func() { ws.Wait(); close(results) }()
	return results
}
