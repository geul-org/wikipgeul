//ff:func feature=convert type=command control=iteration dimension=1
//ff:what TestRunPack — pack mode writes the Issuer's set (every property, then the Q range) readable by documents and h/, the report; bad ID list, resume without a stream and an unopenable report fail
package exportgeul

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/pack"
)

func TestRunPack(t *testing.T) {
	_, schema := fixtureDB(t)
	root, reg := t.TempDir(), t.TempDir()
	fixtureRegistry(t, reg)
	a := strings.Repeat("a", 40)
	args := func(extra ...string) []string {
		return append([]string{"-root", root, "-langreg-root", reg, "-from", "1", "-to", "100", "-commit", a, "-lib-commit", a, "-schema", schema, "-workers", "2", "-checkpoint", "1", "-recorded-at", "1800000000"}, extra...)
	}
	rep := filepath.Join(root, "r.tsv")
	if err := Run(args("-report", rep)); err != nil {
		t.Fatal(err)
	}
	r, err := pack.Open(context.Background(), pack.FileStore{Root: filepath.Join(root, ".well-known/geul/0x00010000")}, pack.Options{})
	if err != nil {
		t.Fatal(err)
	}
	doc, err := r.Fetch(context.Background(), "0x00010000/0x510000000000002a/0x00000001.geul")
	ws, _ := core.ReadWords(doc)
	h, _ := core.ContentHash(ws)
	viaH, err2 := r.Fetch(context.Background(), "h/"+fmtHash(h)+".geul")
	if err != nil || err2 != nil || string(viaH) != string(doc) {
		t.Fatalf("%v %v", err, err2)
	}
	if b, _ := os.ReadFile(rep); strings.Count(string(b), "\tok\t") != 8 {
		t.Fatalf("report:\n%s", b)
	}
	broken := t.TempDir()
	os.MkdirAll(filepath.Join(broken, ".well-known/geul/0x0000"), 0o755)
	os.WriteFile(filepath.Join(broken, ".well-known/geul/0x0000/latest.txt"), []byte("1 00\n"), 0o644)
	for name, extra := range map[string][]string{
		"bad ID list":         {"-root", t.TempDir(), "-ids", filepath.Join(root, "none")},
		"resume, no stream":   {"-resume", "-root", t.TempDir()},
		"unopenable report":   {"-root", t.TempDir(), "-report", root},
		"existing set (new)":  {},
		"broken registry set": {"-root", t.TempDir(), "-langreg-root", broken},
	} {
		if err := Run(args(extra...)); err == nil {
			t.Errorf("%s: want error", name)
		}
	}
	if err := Run(append(args("-root", t.TempDir()), "-recorded-at", "0")); err != nil {
		t.Fatalf("RECORDED-AT from the clock: %v", err)
	}
}
