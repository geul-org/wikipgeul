//ff:func feature=convert type=helper control=sequence
//ff:what TestDefaultWorkers — at least one, at most a third of the CPUs
package exportgeul

import (
	"runtime"
	"testing"
)

func TestDefaultWorkers(t *testing.T) {
	if n := defaultWorkers(); n < 1 || n > max(1, runtime.NumCPU()/2) {
		t.Fatal(n)
	}
}
