//ff:func feature=cli type=command control=iteration dimension=1
//ff:what TestParseRepairFlags — defaults, a line range, an entity range; bad values, both ranges, a bad entity range and an unknown flag fail
package main

import (
	"testing"
	"time"
)

func TestParseRepairFlags(t *testing.T) {
	f, err := parseRepairFlags("wikipg repair-reference-snak", nil)
	if err != nil || f.Config.Schema != "wikidata" || f.Config.DumpID != 1 || f.Config.DryRun || f.Config.Batch != 1000 || f.Config.Workers < 1 || f.Config.Status != 30*time.Second || f.Config.Log == nil {
		t.Fatalf("defaults: %+v %v", f, err)
	}
	f, err = parseRepairFlags("wikipg repair-reference-snak", []string{"-dry-run", "-from-line", "2", "-to-line", "9", "-dsn", "host=x", "-pause", "1s"})
	if err != nil || !f.Config.DryRun || f.Config.From != 2 || f.Config.To != 9 || f.DSN != "host=x" || f.Config.Pause != time.Second {
		t.Fatalf("lines: %+v %v", f, err)
	}
	f, err = parseRepairFlags("wikipg repair-reference-snak", []string{"-entities", "Q1-Q10000", "-dump-id", "2"})
	if err != nil || f.Config.Kind != "Q" || f.Config.Hi != 10000 || f.Config.DumpID != 2 {
		t.Fatalf("entities: %+v %v", f, err)
	}
	for _, bad := range [][]string{{"-dump-id", "0"}, {"-workers", "0"}, {"-from-line", "5", "-to-line", "4"}, {"-entities", "Q1-Q2", "-from-line", "3"}, {"-entities", "x"}, {"-nope"}} {
		if _, err := parseRepairFlags("wikipg repair-reference-snak", bad); err == nil {
			t.Fatalf("%v: want error", bad)
		}
	}
}
