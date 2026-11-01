//ff:func feature=cli type=helper control=sequence
//ff:what TestProgressText — line, totals and the rate of this run; no elapsed time gives a zero rate
package main

import (
	"strings"
	"testing"
	"time"
)

func TestProgressText(t *testing.T) {
	start := time.Unix(100, 0)
	p := progress{Line: 7, Entities: 30, Errors: 1, Start: start}
	if got := p.Text(start.Add(10*time.Second), 10); got != "line 7  entities 30  errors 1  2 entities/s" {
		t.Fatalf("got %q", got)
	}
	if got := p.Text(start, 0); !strings.HasSuffix(got, " 0 entities/s") {
		t.Fatalf("zero: %q", got)
	}
}
