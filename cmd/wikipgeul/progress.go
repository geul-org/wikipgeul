//ff:type feature=cli type=model
//ff:what progress — the running totals of a load: checkpoint line, entities and errors, and when it started
package main

import "time"

// progress is what a load has committed so far: the checkpoint line, the
// entities and load errors written, and the start time for the rate.
type progress struct {
	Line     int64
	Entities int64
	Errors   int64
	Start    time.Time
}
