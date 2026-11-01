//ff:func feature=cli type=helper control=sequence
//ff:what TestCmdReaderClose — closing waits for the process and returns its exit status
package main

import (
	"os/exec"
	"testing"
)

func TestCmdReaderClose(t *testing.T) {
	ok := exec.Command("true")
	out, _ := ok.StdoutPipe()
	ok.Start()
	if err := (cmdReader{ReadCloser: out, Cmd: ok}).Close(); err != nil {
		t.Fatal(err)
	}
	fail := exec.Command("false")
	out, _ = fail.StdoutPipe()
	fail.Start()
	if err := (cmdReader{ReadCloser: out, Cmd: fail}).Close(); err == nil {
		t.Fatal("false: want error")
	}
}
