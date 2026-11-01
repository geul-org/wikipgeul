//ff:func feature=cli type=command control=sequence
//ff:what parseRepairFlags — the command line of repair-reference-snak and build-reference-order: range by lines or entities, dry run, batch, workers, pace; defaults keep 30% of the machine free
package main

import (
	"errors"
	"flag"
	"os"
	"runtime"
	"time"
)

// parseRepairFlags reads the arguments of wikipg repair-reference-snak or
// build-reference-order (name, the command shown in usage and errors). The
// range is -from-line/-to-line (default: the whole dump) or -entities
// (Q1-Q10000), not both. The defaults keep the machine about 30% free: 70%
// of the CPUs, half of it for the workers (each parses in Go while its
// database connection decompresses the lz4 raw lines) and half left to the
// server, batches of 1000 lines.
func parseRepairFlags(name string, args []string) (repairFlags, error) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	var f repairFlags
	c := &f.Config
	var id int
	var entities string
	fs.StringVar(&f.DSN, "dsn", "", "PostgreSQL DSN (empty: PG* environment variables)")
	fs.StringVar(&c.Schema, "schema", "wikidata", "schema made by sql/wikidata.sql")
	fs.IntVar(&id, "dump-id", 1, "dump_id to repair (1..255)")
	fs.BoolVar(&c.DryRun, "dry-run", false, "only compare and report; read-only session, nothing written")
	fs.Int64Var(&c.From, "from-line", 0, "first dump line of the range")
	fs.Int64Var(&c.To, "to-line", 0, "last dump line of the range (0: the dump's last)")
	fs.StringVar(&entities, "entities", "", "entity range instead of lines, e.g. Q1-Q10000")
	fs.IntVar(&c.Batch, "batch", 1000, "dump lines per batch (one transaction when writing)")
	fs.IntVar(&c.Workers, "workers", max(runtime.NumCPU()*7/10/2, 1), "parallel readers (database connections)")
	fs.DurationVar(&c.Pause, "pause", 0, "pause after each batch per worker (pace limit)")
	fs.DurationVar(&c.Status, "status", 30*time.Second, "status line interval")
	fs.IntVar(&c.RefLimit, "ref-limit", 20000000, "hashes remembered in memory to skip repeats (beyond it the done table decides)")
	if err := fs.Parse(args); err != nil {
		return repairFlags{}, err
	}
	if id < 1 || id > 255 || c.Batch < 1 || c.Workers < 1 || c.Status <= 0 || c.Pause < 0 || c.From < 0 || c.To < 0 || c.To != 0 && c.To < c.From {
		return repairFlags{}, errors.New("wikipg: need -dump-id 1..255, -batch, -workers ≥ 1, -status > 0, -pause ≥ 0, 0 ≤ -from-line ≤ -to-line")
	}
	c.DumpID, c.Log = int16(id), os.Stderr
	if entities == "" {
		return f, nil
	}
	if c.From != 0 || c.To != 0 {
		return repairFlags{}, errors.New("wikipg: -entities and -from-line/-to-line exclude each other")
	}
	var err error
	c.Kind, c.Lo, c.Hi, err = parseEntities(entities)
	return f, err
}
