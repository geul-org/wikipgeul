//ff:func feature=cli type=helper control=sequence
//ff:what TestLangSet — a list becomes a set; no list is nil (all languages)
package main

import "testing"

func TestLangSet(t *testing.T) {
	if m := langSet([]string{"ko", "en"}); !m["ko"] || !m["en"] || len(m) != 2 {
		t.Fatalf("got %v", m)
	}
	if langSet(nil) != nil {
		t.Fatal("all: want nil")
	}
}
