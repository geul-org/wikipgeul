//ff:type feature=cli type=model
//ff:what loadConfig — what wikipg runs: database, schema, input, line range, parallelism per stage, kept languages, raw, the dump record
package main

import (
	"time"

	"github.com/geul-org/wikipgeul/internal/load"
)

// loadConfig is the parsed command line of wikipg: the PostgreSQL DSN
// (empty: PG* variables) and schema, the dump file,
// the first line to read and how many entity lines at most (0: all), the
// parallelism of each stage (decompressor threads, parsers, writers) and the
// decompressor program, the entities and bytes per batch at most, the term languages (nil:
// all), whether raw lines are kept, the reference hash memory bound, how
// often a status line is printed and the dump record.
type loadConfig struct {
	DSN          string
	Schema       string
	In           string
	From         int64
	Limit        int64
	Decoders     int
	Parsers      int
	Writers      int
	Decompressor string
	Batch        int
	BatchBytes   int64
	Langs        []string
	Raw          bool
	RefLimit     int
	Status       time.Duration
	Dump         load.DumpInfo
}
