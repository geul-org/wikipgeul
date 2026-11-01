//ff:func feature=convert type=helper control=sequence
//ff:what TestWriteIndex — documents.geul under the Issuer directory names this run's and earlier published documents; a bad Issuer and an unwritable root fail
package exportgeul

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

func TestWriteIndex(t *testing.T) {
	root := t.TempDir()
	if err := writeIndex(root, []core.Word{1, 0}, []uint64{0x5000_0000_0000_001F}); err != nil {
		t.Fatal(err)
	}
	os.Mkdir(filepath.Join(root, ".well-known/geul/0x00010000/0x5100000000000007"), 0o755)
	if err := writeIndex(root, []core.Word{1, 0}, []uint64{0x5000_0000_0000_001F}); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(root, ".well-known/geul/0x00010000/documents.geul"))
	if err != nil || !bytes.Contains(b, []byte("q7")) || !bytes.Contains(b, []byte("p31")) {
		t.Fatalf("names: %v", err)
	}
	if err := writeIndex(root, []core.Word{0, 1}, []uint64{0x5000_0000_0000_001F}); err == nil {
		t.Fatal("bad Issuer: want error")
	}
	blocked := filepath.Join(root, "blocked")
	os.WriteFile(blocked, nil, 0o644)
	os.MkdirAll(filepath.Join(root, "unreadable/.well-known/geul"), 0o755)
	os.WriteFile(filepath.Join(root, "unreadable/.well-known/geul/0x00010000"), nil, 0o644)
	if err := writeIndex(filepath.Join(root, "unreadable"), []core.Word{1, 0}, nil); err == nil {
		t.Fatal("Issuer path is a file: want error")
	}
	ro := filepath.Join(root, "ro/.well-known/geul")
	os.MkdirAll(ro, 0o755)
	os.Chmod(ro, 0o555)
	defer os.Chmod(ro, 0o755)
	if err := writeIndex(filepath.Join(root, "ro"), []core.Word{1, 0}, nil); err == nil {
		t.Fatal("read-only parent: want error")
	}
	if err := writeIndex(blocked, []core.Word{1, 0}, []uint64{0x5000_0000_0000_001F}); err == nil {
		t.Fatal("unwritable: want error")
	}
}
