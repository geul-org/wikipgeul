//ff:func feature=cli type=helper control=sequence
//ff:what parseEntities — an entity range like Q1-Q10000 as its kind letter and number bounds
package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// entityRange matches Q1-Q10000 (one kind on both ends).
var entityRange = regexp.MustCompile(`^([A-Z])(\d+)-([A-Z])(\d+)$`)

// parseEntities reads s, an entity range such as Q1-Q10000 (the same letter
// on both ends, lower bound not above the upper), as kind, lo and hi.
func parseEntities(s string) (string, int64, int64, error) {
	m := entityRange.FindStringSubmatch(s)
	if m == nil || m[1] != m[3] {
		return "", 0, 0, fmt.Errorf("wikipg: -entities %q is not a range like Q1-Q10000", s)
	}
	lo, _ := strconv.ParseInt(m[2], 10, 64)
	hi, _ := strconv.ParseInt(m[4], 10, 64)
	if lo > hi {
		return "", 0, 0, fmt.Errorf("wikipg: -entities %q: lower bound above upper", s)
	}
	return m[1], lo, hi, nil
}
