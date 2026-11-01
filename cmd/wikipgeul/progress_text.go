//ff:func feature=cli type=helper control=sequence
//ff:what progress.Text — one status line: checkpoint line, entities, errors, entities per second
package main

import (
	"fmt"
	"time"
)

// Text is one line of status at now: the checkpoint line, the entities and
// errors written in this run's totals and the entity rate since Start.
func (p progress) Text(now time.Time, base int64) string {
	secs := now.Sub(p.Start).Seconds()
	rate := 0.0
	if secs > 0 {
		rate = float64(p.Entities-base) / secs
	}
	return fmt.Sprintf("line %d  entities %d  errors %d  %.0f entities/s", p.Line, p.Entities, p.Errors, rate)
}
