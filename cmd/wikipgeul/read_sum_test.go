//ff:func feature=cli type=helper control=sequence
//ff:what TestReadSum — the hash of the named file, binary-mode names included; other names and missing files are empty
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadSum(t *testing.T) {
	p := filepath.Join(t.TempDir(), "sums")
	os.WriteFile(p, []byte("aa  /x/a.bz2\nbb *b.bz2\nbad line here\n"), 0o644)
	if readSum(p, "a.bz2") != "aa" || readSum(p, "b.bz2") != "bb" || readSum(p, "c.bz2") != "" || readSum(p+"x", "a.bz2") != "" {
		t.Fatal("wrong sums")
	}
}
