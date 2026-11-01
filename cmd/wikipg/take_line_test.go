//ff:func feature=cli type=command control=sequence
//ff:what TestTakeLine — an entity line at or after the start is taken; brackets and earlier lines are not
package main

import "testing"

func TestTakeLine(t *testing.T) {
	var b batch
	if takeLine(&b, []byte("{\"id\":\"Q1\"},\n"), 5, 3) != 1 || b.Last != 5 || b.Bytes != 11 || string(b.Lines[0]) != `{"id":"Q1"}` {
		t.Fatalf("take: %+v", b)
	}
	if takeLine(&b, []byte("[\n"), 6, 3) != 0 || takeLine(&b, []byte("{}\n"), 2, 3) != 0 || len(b.Lines) != 1 {
		t.Fatalf("skip: %+v", b)
	}
}
