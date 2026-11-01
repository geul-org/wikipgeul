//ff:func feature=cli type=command control=sequence
//ff:what parseLoadFlags — wikipg's command line as a loadConfig; -in and -commit are required, the dump record is filled from the file
package main

import (
	"errors"
	"flag"
	"runtime"
	"strings"
	"time"
)

// parseLoadFlags reads args. -in (the dump file) and -commit (the loader's git commit, recorded on the dump row)
// are required. -langs is a comma list or "all". The defaults keep the
// machine about 30% free: of 70% of the CPUs, a fifth decompress and two
// fifths parse, with 4 writers (database connections, whose server processes
// take the rest), and batches of 8 MiB keep memory near 1 GB (measured
// 2026-10-26: 4 writers and 8 MiB were fastest; more writers contend on table
// extension, bigger batches only add memory). The rest of the dump record
// comes from the file (dumpInfoOf).
func parseLoadFlags(args []string) (loadConfig, error) {
	fs := flag.NewFlagSet("wikipg", flag.ContinueOnError)
	var c loadConfig
	var langs, date, published string
	var id int
	fs.StringVar(&c.DSN, "dsn", "", "PostgreSQL DSN (empty: PG* environment variables)")
	fs.StringVar(&c.Schema, "schema", "wikidata", "schema made by sql/wikidata.sql")
	fs.StringVar(&c.In, "in", "", "dump file, .json or .json.bz2 (required)")
	fs.Int64Var(&c.From, "from-line", 0, "first dump line to read (0: resume after the checkpoint)")
	fs.Int64Var(&c.Limit, "limit", 0, "entity lines to load at most (0: all)")
	budget := max(runtime.NumCPU()*7/10, 3)
	fs.IntVar(&c.Decoders, "decoders", max(budget/5, 1), "decompressor threads (lbzip2 -n)")
	fs.IntVar(&c.Parsers, "parsers", max(budget*2/5, 1), "parallel JSON parsers")
	fs.IntVar(&c.Writers, "writers", 4, "parallel database writers (connections)")
	fs.StringVar(&c.Decompressor, "decompressor", defaultDecompressor(), "bzip2 program: lbzip2 (parallel) or bzip2")
	fs.IntVar(&c.Batch, "batch", 100000, "entities per batch at most (one transaction)")
	var batchMiB int64
	fs.Int64Var(&batchMiB, "batch-mib", 8, "raw MiB per batch at most (one transaction; bounds memory and the write chunk)")
	fs.StringVar(&langs, "langs", "all", `label, description, alias and gloss languages, a comma list or "all" (default: every language, a lossless load)`)
	fs.BoolVar(&c.Raw, "raw", true, "keep each dump line in entity_raw")
	fs.IntVar(&c.RefLimit, "ref-limit", 4000000, "reference hashes remembered to skip repeats")
	fs.DurationVar(&c.Status, "status", 30*time.Second, "status line interval")
	fs.IntVar(&id, "dump-id", 1, "dump_id of this dump (1..255)")
	fs.StringVar(&date, "dump-date", "", "dump date YYYYMMDD (default: from the file name)")
	fs.StringVar(&c.Dump.URL, "url", "", "where the dump was downloaded (default: dumps.wikimedia.org path of the file)")
	fs.StringVar(&published, "published", "", "publication time, RFC 3339 (optional)")
	fs.StringVar(&c.Dump.LoaderCommit, "commit", "", "loader git commit (required)")
	if err := fs.Parse(args); err != nil {
		return loadConfig{}, err
	}
	if c.In == "" || c.Dump.LoaderCommit == "" {
		return loadConfig{}, errors.New("wikipg: -in and -commit are required")
	}
	if id < 1 || id > 255 || c.Status <= 0 || c.Decoders < 1 || c.Parsers < 1 || c.Writers < 1 || c.Batch < 1 || batchMiB < 1 || c.From < 0 || c.Limit < 0 {
		return loadConfig{}, errors.New("wikipg: need -dump-id 1..255, -decoders, -parsers, -writers, -batch and -batch-mib ≥ 1, -from-line and -limit ≥ 0")
	}
	if langs != "all" {
		c.Langs = strings.Split(strings.ReplaceAll(langs, " ", ""), ",")
	}
	c.Dump.ID, c.BatchBytes = int16(id), batchMiB<<20
	return c, dumpInfoOf(&c, date, published)
}
