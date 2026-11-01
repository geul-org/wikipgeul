//ff:type feature=cli type=model
//ff:what cmdReader — the output of a decompressor process, closed by waiting for the process
package main

import (
	"io"
	"os/exec"
)

// cmdReader reads the standard output of Cmd (a decompressor). Closing it
// closes the pipe and waits for the process.
type cmdReader struct {
	io.ReadCloser
	Cmd *exec.Cmd
}
