//ff:func feature=cli type=helper control=sequence
//ff:what openInput — the dump's JSON lines: a .bz2 through the decompressor program (lbzip2 with threads, or bzip2), anything else as a plain file
package main

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// openInput opens the dump at path for reading its JSON lines. A .bz2 file is
// read through the decompressor program: lbzip2 runs with threads threads
// (lbzip2 -dc -n), which parallelizes even a single-stream file; anything
// else (bzip2) runs as prog -dc. Other files are read as they are.
func openInput(path, prog string, threads int) (io.ReadCloser, error) {
	if !strings.HasSuffix(path, ".bz2") {
		return os.Open(path)
	}
	args := []string{"-dc"}
	if strings.HasPrefix(filepath.Base(prog), "lbzip2") {
		args = append(args, "-n", strconv.Itoa(threads))
	}
	cmd := exec.Command(prog, append(args, path)...)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmdReader{ReadCloser: out, Cmd: cmd}, nil
}
