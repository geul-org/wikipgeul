//ff:func feature=cli type=helper control=sequence
//ff:what cmdReader.Close — closes the pipe and waits for the decompressor
package main

// Close closes the output pipe and waits for the process, returning its exit
// error. A process stopped early by the closed pipe reports that error too.
func (r cmdReader) Close() error {
	r.ReadCloser.Close()
	return r.Cmd.Wait()
}
