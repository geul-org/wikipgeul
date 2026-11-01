//ff:func feature=convert type=helper control=sequence
//ff:what fmtHash — a 32-byte hash as 0x and 64 hex digits
package exportgeul

import "fmt"

// fmtHash writes h as the h/ path names it.
func fmtHash(h [32]byte) string {
	return fmt.Sprintf("0x%x", h)
}
