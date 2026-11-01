//ff:func feature=convert type=helper control=iteration dimension=1
//ff:what publishedNumbers — the document numbers already published in an Issuer directory (0x + 16 hex digit folders)
package exportgeul

import (
	"errors"
	"io/fs"
	"os"
	"strconv"
)

// publishedNumbers lists the document folders of an Issuer directory, named
// 0x and 16 hex digits (the 4-word document ID), as numbers. A missing
// directory has none; other entries are ignored.
func publishedNumbers(dir string) ([]uint64, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []uint64
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() || len(name) != 18 || name[:2] != "0x" {
			continue
		}
		if n, err := strconv.ParseUint(name[2:], 16, 64); err == nil {
			out = append(out, n)
		}
	}
	return out, nil
}
