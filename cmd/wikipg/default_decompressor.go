//ff:func feature=cli type=helper control=sequence
//ff:what defaultDecompressor — lbzip2 when it is on PATH (parallel on a single-stream dump), else bzip2
package main

import "os/exec"

// defaultDecompressor is the bzip2 program used when -decompressor is not
// given: lbzip2 when it is on PATH, since it decompresses even the
// single-stream Wikidata dump in parallel (about 9 times bzip2 with 16
// threads), else bzip2.
func defaultDecompressor() string {
	if _, err := exec.LookPath("lbzip2"); err == nil {
		return "lbzip2"
	}
	return "bzip2"
}
