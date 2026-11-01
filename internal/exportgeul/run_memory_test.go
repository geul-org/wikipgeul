//ff:func feature=convert type=command control=sequence
//ff:what TestRunMemory — pack mode runs under the memory guard: memory that cannot be read fails; a run out of memory stops with a StopError and -resume gives the set of a run never stopped
package exportgeul

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/geul-org/wikigeul/pkg/memguard"
)

func TestRunMemory(t *testing.T) {
	_, schema := fixtureDB(t)
	reg := t.TempDir()
	fixtureRegistry(t, reg)
	mi := filepath.Join(t.TempDir(), "meminfo")
	free := func(pct int) { os.WriteFile(mi, []byte(fmt.Sprintf("MemTotal: 1000000 kB\nMemAvailable: %d kB\n", pct*10000)), 0o644) }
	a := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	args := func(root string, extra ...string) []string {
		return append([]string{"-root", root, "-langreg-root", reg, "-from", "1", "-to", "3", "-commit", a, "-lib-commit", a, "-schema", schema, "-workers", "2",
			"-recorded-at", "1800000000", "-samples", "2", "-checkpoint", "1", "-mem-meminfo", mi}, extra...)
	}
	if err := Run(args(t.TempDir())); err == nil {
		t.Fatal("no meminfo: want an error")
	}
	free(5)
	root := t.TempDir()
	var stop memguard.StopError
	if err := Run(args(root)); !errors.As(err, &stop) {
		t.Fatalf("out of memory: %v", err)
	}
	free(60)
	if err := Run(args(root, "-resume")); err != nil {
		t.Fatalf("resume: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".well-known/geul/0x00010000/latest.txt")); err != nil {
		t.Fatal(err)
	}
}
