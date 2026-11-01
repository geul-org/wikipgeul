//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what readSum — the hash a checksum file (sha256sum / md5sums format) gives for a file name; empty when absent
package main

import (
	"os"
	"path/filepath"
	"strings"
)

// readSum reads a checksum file in the "<hash>  <name>" format (sha256sum,
// Wikimedia's md5sums.txt) and returns the hash of the line whose name ends
// with name. A missing file or line is "".
func readSum(path, name string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(b), "\n") {
		f := strings.Fields(line)
		if len(f) == 2 && filepath.Base(strings.TrimPrefix(f[1], "*")) == name {
			return f[0]
		}
	}
	return ""
}
