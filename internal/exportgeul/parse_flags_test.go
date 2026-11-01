//ff:func feature=convert type=command control=iteration dimension=1
//ff:what TestParseFlags — database defaults and overrides on top of the shared flags; a -dump-id out of range, an unknown -out, the dump file's flags and bad shared flags fail
package exportgeul

import (
	"strings"
	"testing"
)

func TestParseFlags(t *testing.T) {
	a, b := strings.Repeat("a", 40), strings.Repeat("b", 40)
	need := []string{"-commit", a, "-lib-commit", b}
	c, err := parseFlags(need)
	if err != nil || c.Root != "." || c.From != 1 || c.To != 1000 || c.Commit != a || c.LibCommit != b || c.Schema != "wikidata" || c.Dump != 1 || c.DSN != "" || c.Out != "pack" || c.Workers != defaultWorkers() {
		t.Fatalf("defaults: %+v %v", c, err)
	}
	c, err = parseFlags(append(need, "-schema", "wd", "-dump-id", "3", "-dsn", "host=x", "-out", "dir", "-workers", "3"))
	if err != nil || c.Schema != "wd" || c.Dump != 3 || c.DSN != "host=x" || c.Out != "dir" || c.Workers != 3 {
		t.Fatalf("overrides: %+v %v", c, err)
	}
	for _, args := range [][]string{{}, {"-commit", a}, {"-nope"}, append(need, "-dump-id", "0"), append(need, "-dump-id", "256"), append(need, "-out", "zip"), append(need, "-dump", "d.json"), append(need, "-langs", "en")} {
		if _, err := parseFlags(args); err == nil {
			t.Fatalf("%v: want error", args)
		}
	}
}
