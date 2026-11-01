//ff:func feature=cli type=command control=selection
//ff:what collector.step — waits for the next result or status tick; false once every writer is done
package main

import (
	"fmt"
	"os"
	"time"
)

// step waits for one result, which it applies (applyResult), or one tick, on
// which it prints a status line. It is false when the results channel is
// closed.
func (k *collector) step(results <-chan result, tick <-chan time.Time) bool {
	select {
	case r, ok := <-results:
		if ok {
			k.Failed = applyResult(k.DB, k.C, k.CP, &k.P, r, k.Failed, k.Cancel)
		}
		return ok
	case now := <-tick:
		fmt.Fprintln(os.Stderr, k.P.Text(now, k.Base))
		return true
	}
}
