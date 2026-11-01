//ff:func feature=cli type=helper control=sequence
//ff:what TestOpenInput — a plain file and a .bz2 (bzip2, and an lbzip2 given -n threads) read the same bytes; a missing program or file fails, no free file descriptor for the pipe fails, a bad .bz2 fails on close
package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
)

func TestOpenInput(t *testing.T) {
	p := writeDump(t, 3, 0)
	want, _ := os.ReadFile(p)
	r, err := openInput(p, "bzip2", 1)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(r)
	r.Close()
	if string(got) != string(want) {
		t.Fatal("plain differs")
	}
	if err := exec.Command("bzip2", "-k", p).Run(); err != nil {
		t.Skipf("bzip2: %v", err)
	}
	r, err = openInput(p+".bz2", "bzip2", 1)
	if err != nil {
		t.Fatal(err)
	}
	got, _ = io.ReadAll(r)
	if err := r.Close(); err != nil || string(got) != string(want) {
		t.Fatalf("bz2: %v", err)
	}
	fake := filepath.Join(t.TempDir(), "lbzip2")
	os.WriteFile(fake, []byte("#!/bin/sh\n[ \"$2\" = -n ] || exit 3\nexec bzip2 -dc \"$4\"\n"), 0o755)
	r, _ = openInput(p+".bz2", fake, 4)
	if got, _ = io.ReadAll(r); r.Close() != nil || string(got) != string(want) {
		t.Fatal("lbzip2 with threads differs")
	}
	if _, err := openInput(p+".bz2", "no-such-decompressor", 1); err == nil {
		t.Fatal("missing program: want error")
	}
	// The stdout pipe needs two new descriptors: with the soft limit at the
	// lowest free descriptor number, the kernel has none to give (EMFILE).
	var lim syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
		t.Fatal(err)
	}
	probe, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	low := uint64(probe.Fd())
	probe.Close()
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: low, Max: lim.Max}); err != nil {
		t.Fatal(err)
	}
	r, err = openInput(p+".bz2", "bzip2", 1)
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &lim); err != nil {
		t.Fatal(err)
	}
	if r != nil || !errors.Is(err, syscall.EMFILE) {
		t.Fatalf("no free descriptor: %v %v", r, err)
	}
	if _, err := openInput(filepath.Join(t.TempDir(), "none.json"), "bzip2", 1); err == nil {
		t.Fatal("missing: want error")
	}
	bad := filepath.Join(t.TempDir(), "bad.json.bz2")
	os.WriteFile(bad, []byte("not bzip2"), 0o644)
	r, err = openInput(bad, "bzip2", 1)
	if err != nil {
		t.Fatal(err)
	}
	io.ReadAll(r)
	if err := r.Close(); err == nil {
		t.Fatal("bad bz2: want error on close")
	}
}
