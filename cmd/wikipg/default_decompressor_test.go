//ff:func feature=cli type=helper control=sequence
//ff:what TestDefaultDecompressor — lbzip2 when on PATH, bzip2 otherwise
package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDefaultDecompressor(t *testing.T) {
	want := "bzip2"
	if _, err := exec.LookPath("lbzip2"); err == nil {
		want = "lbzip2"
	}
	if got := defaultDecompressor(); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "lbzip2"), []byte("#!/bin/sh\n"), 0o755)
	t.Setenv("PATH", dir)
	if got := defaultDecompressor(); got != "lbzip2" {
		t.Fatalf("lbzip2 on PATH: %q", got)
	}
	t.Setenv("PATH", t.TempDir())
	if got := defaultDecompressor(); got != "bzip2" {
		t.Fatalf("no lbzip2: %q", got)
	}
}
