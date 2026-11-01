//ff:func feature=convert type=command control=sequence
//ff:what parseFlags — the run config from export-geul's arguments; full wikigeul and library commits are required
package exportgeul

import (
	"errors"
	"flag"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

// parseFlags reads args into a config: the flags every conversion run shares
// (pipeline.BindFlags) and the database's. -schema and -dump-id name the
// loaded dump (default wikidata, 1).
func parseFlags(args []string) (config, error) {
	fs := flag.NewFlagSet("wikipg export-geul", flag.ContinueOnError)
	var c config
	var dump int
	check := pipeline.BindFlags(fs, &c.Config, defaultWorkers(), "pack mode: conversion goroutines, each with its own database connection")
	fs.StringVar(&c.DSN, "dsn", "", "PostgreSQL DSN (empty: PG* environment variables)")
	fs.StringVar(&c.Schema, "schema", "wikidata", "schema of the loaded dump (wikipg -schema)")
	fs.IntVar(&dump, "dump-id", 1, "dump_id of the loaded dump (1..255)")
	fs.StringVar(&c.Out, "out", "pack", "output: pack (the Issuer's pack set, zstd, streamed) or dir (directory tree, for tests)")
	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	if err := check(); err != nil {
		return config{}, err
	}
	if dump < 1 || dump > 255 {
		return config{}, errors.New("wikipg export-geul: need 1 <= -dump-id <= 255")
	}
	c.Dump = int16(dump)
	if c.Out != "pack" && c.Out != "dir" {
		return config{}, errors.New("wikipg export-geul: -out is pack or dir")
	}
	return c, nil
}
