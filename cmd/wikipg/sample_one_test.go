//ff:func feature=cli type=command control=sequence
//ff:what sampleOne — a minimal decoded item Q1 with one English label
package main

import "github.com/geul-org/wikipgeul/internal/dump"

// sampleOne returns item Q1 with an English label, for writer tests.
func sampleOne() dump.Entity {
	return dump.Entity{ID: "Q1", Type: "item", Labels: map[string]dump.Term{"en": {Language: "en", Value: "universe"}}}
}
