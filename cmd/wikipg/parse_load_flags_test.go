//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestParseLoadFlags — defaults and overrides; missing -in or -commit, bad ranges, bad -published and an unknown flag fail
package main

import (
	"slices"
	"testing"
)

func TestParseLoadFlags(t *testing.T) {
	p := writeDump(t, 1, 0)
	c, err := parseLoadFlags([]string{"-in", p, "-commit", "abc"})
	if err != nil || c.Schema != "wikidata" || c.Writers != 4 || c.Parsers < 1 || c.Decoders < 1 || c.Decompressor == "" || c.Batch != 100000 || c.BatchBytes != 8<<20 || c.Langs != nil || !c.Raw || c.Dump.ID != 1 || c.Dump.SHA256 != "abc123" {
		t.Fatalf("defaults: %+v %v", c, err)
	}
	c, err = parseLoadFlags([]string{"-in", p, "-commit", "abc", "-langs", "ko,en", "-raw=false", "-writers", "2", "-parsers", "3", "-decoders", "4", "-decompressor", "lbzip2", "-batch", "10", "-dump-id", "3", "-published", "2026-10-21T02:33:31Z"})
	if err != nil || !slices.Equal(c.Langs, []string{"ko", "en"}) || c.Raw || c.Writers != 2 || c.Parsers != 3 || c.Decoders != 4 || c.Decompressor != "lbzip2" || c.Dump.ID != 3 || c.Dump.PublishedAt.IsZero() {
		t.Fatalf("overrides: %+v %v", c, err)
	}
	for _, args := range [][]string{
		{"-commit", "a"}, {"-in", p}, {"-in", p, "-commit", "a", "-dump-id", "0"}, {"-in", p, "-commit", "a", "-writers", "0"}, {"-in", p, "-commit", "a", "-parsers", "0"}, {"-in", p, "-commit", "a", "-decoders", "0"}, {"-in", p, "-commit", "a", "-batch-mib", "0"},
		{"-in", p, "-commit", "a", "-published", "x"}, {"-nope"},
	} {
		if _, err := parseLoadFlags(args); err == nil {
			t.Fatalf("%v: want error", args)
		}
	}
}
