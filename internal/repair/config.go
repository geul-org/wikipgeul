//ff:type feature=repair type=model
//ff:what Config — what a reference_snak repair or order build reads and does: schema, dump, line or entity range, order mode, dry run, batch, workers, pace, log
package repair

import (
	"io"
	"time"
)

// Config is one repair run: the schema and dump, the range (lines From..To,
// To 0 meaning the dump's last line; or, when Kind is set, the entities
// Kind Lo..Hi, e.g. Q 1..10000), whether it builds reference_snak_order
// instead of repairing reference_snak (Order: wikipg build-reference-order;
// reference_snak is only read), whether it only compares (DryRun), the dump
// lines per batch, the parallel workers (each one database connection reading
// and parsing a batch), a pause after each batch per worker (pace), how often
// a status line is written to Log and the bound of the in-memory hash set
// (load.RefSet).
type Config struct {
	Schema   string
	DumpID   int16
	From     int64
	To       int64
	Kind     string
	Lo       int64
	Hi       int64
	Order    bool
	DryRun   bool
	Batch    int
	Workers  int
	Pause    time.Duration
	Status   time.Duration
	RefLimit int
	Log      io.Writer
}
