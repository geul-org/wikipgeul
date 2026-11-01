//ff:func feature=convert type=command control=iteration dimension=1
//ff:what TestRun — a small range writes both Issuers' documents, indexes and the report with the language registry read from the root (or -langreg-root); bad flags, no registry, a broken registry pack set, an unparsable DSN, an unwritable root and unwritable Issuer directories fail
package exportgeul

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	_, schema := fixtureDB(t)
	root := t.TempDir()
	fixtureRegistry(t, root)
	report := filepath.Join(root, "r.tsv")
	if err := Run([]string{"-out", "dir", "-root", root, "-from", "1", "-to", "3", "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema, "-report", report}); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{".well-known/geul/0x00010000/documents.geul", ".well-known/geul/0x00010000/0x500000000000001f", "r.tsv"} {
		if _, err := os.Stat(filepath.Join(root, p)); err != nil {
			t.Fatal(err)
		}
	}
	if err := Run([]string{"-out", "dir", "-root", root, "-from", "1", "-to", "1", "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema}); err != nil {
		t.Fatalf("no report: %v", err)
	}
	broken := t.TempDir()
	os.MkdirAll(filepath.Join(broken, ".well-known/geul/0x0000"), 0o755)
	os.WriteFile(filepath.Join(broken, ".well-known/geul/0x0000/latest.txt"), []byte("1 00\n"), 0o644)
	if err := Run([]string{"-out", "dir", "-root", root, "-langreg-root", broken, "-from", "1", "-to", "1", "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema}); err == nil {
		t.Fatal("broken registry pack set: want error")
	}
	if err := Run(nil); err == nil {
		t.Fatal("no commit: want error")
	}
	for _, dump := range []string{"9", "2"} {
		if err := Run([]string{"-out", "dir", "-root", root, "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema, "-dump-id", dump}); err == nil {
			t.Fatalf("dump %s (missing, other loader): want error", dump)
		}
	}
	if err := Run([]string{"-out", "dir", "-root", root, "-dsn", "host=/nonexistent", "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema}); err == nil {
		t.Fatal("unreachable database: want error")
	}
	if err := Run([]string{"-out", "dir", "-root", root, "-dsn", "host='x", "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema}); err == nil {
		t.Fatal("unparsable DSN: want error")
	}
	if err := Run([]string{"-out", "dir", "-root", t.TempDir(), "-from", "1", "-to", "1", "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema}); err == nil {
		t.Fatal("no language registry: want error")
	}
	blocked := filepath.Join(root, "blocked")
	os.WriteFile(blocked, nil, 0o644)
	if err := Run([]string{"-out", "dir", "-root", blocked, "-langreg-root", root, "-from", "1", "-to", "1", "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema}); err == nil {
		t.Fatal("unwritable: want error")
	}
	// A path that is a file: P31's document folder (the property Q1 uses
	// cannot be published), the Issuer directory (Q1 cannot be published) and,
	// with no documents (Q6 was never created), the index.
	for _, c := range []struct{ file, from string }{
		{"0x00010000/0x500000000000001f", "1"},
		{"0x00010000", "1"},
		{"0x00010000", "6"},
	} {
		r := t.TempDir()
		os.MkdirAll(filepath.Dir(filepath.Join(r, ".well-known/geul", c.file)), 0o755)
		os.WriteFile(filepath.Join(r, ".well-known/geul", c.file), nil, 0o644)
		if err := Run([]string{"-out", "dir", "-root", r, "-langreg-root", root, "-from", c.from, "-to", c.from, "-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-lib-commit", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "-schema", schema}); err == nil {
			t.Fatalf("%s file, Q%s: want error", c.file, c.from)
		}
	}
}
