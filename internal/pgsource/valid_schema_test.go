//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what TestValidSchema — plain lower-case identifiers pass; empty, digit-first, quoted, upper-case and too long names fail
package pgsource

import (
	"strings"
	"testing"
)

func TestValidSchema(t *testing.T) {
	for _, s := range []string{"wikidata", "wikidata_test_1", "_x"} {
		if !validSchema(s) {
			t.Fatalf("%q: want valid", s)
		}
	}
	for _, s := range []string{"", "1x", "wiki.data", "x;drop", "Wikidata", strings.Repeat("a", 64)} {
		if validSchema(s) {
			t.Fatalf("%q: want invalid", s)
		}
	}
}
