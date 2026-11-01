//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what TestParseEntities — Q1-Q10000 splits into kind and bounds; other forms, mixed kinds and reversed bounds fail
package main

import "testing"

func TestParseEntities(t *testing.T) {
	if k, lo, hi, err := parseEntities("Q1-Q10000"); err != nil || k != "Q" || lo != 1 || hi != 10000 {
		t.Fatalf("%s %d %d %v", k, lo, hi, err)
	}
	for _, bad := range []string{"Q1", "Q1-P5", "q1-q2", "Q9-Q1"} {
		if _, _, _, err := parseEntities(bad); err == nil {
			t.Fatalf("%s: want error", bad)
		}
	}
}
