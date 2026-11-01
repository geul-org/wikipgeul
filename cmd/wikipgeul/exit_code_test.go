//ff:func feature=cli type=command control=sequence
//ff:what TestExitCode — an anomaly stop, wrapped too, is 3; any other error is 1
package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestExitCode(t *testing.T) {
	if got := exitCode(fmt.Errorf("run: %w", pipeline.AnomalyError{})); got != 3 {
		t.Fatalf("anomaly: %d", got)
	}
	if got := exitCode(errors.New("x")); got != 1 {
		t.Fatalf("other: %d", got)
	}
}
