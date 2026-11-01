//ff:func feature=cli type=command control=sequence
//ff:what TestExitCode — an anomaly stop, wrapped too, is 3; a memory stop is 4 with the resume command told; any other error is 1
package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/geul-org/wikigeul/pkg/memguard"
	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestExitCode(t *testing.T) {
	var b bytes.Buffer
	if got := exitCode(fmt.Errorf("run: %w", pipeline.AnomalyError{}), nil, &b); got != 3 {
		t.Fatalf("anomaly: %d", got)
	}
	args := []string{"wikipg", "export-geul", "-root", "a"}
	if got := exitCode(fmt.Errorf("run: %w", memguard.StopError{}), args, &b); got != 4 || !strings.Contains(b.String(), "wikipg export-geul -root a -resume") {
		t.Fatalf("memory: %d %q", got, b.String())
	}
	if got := exitCode(errors.New("x"), nil, &b); got != 1 {
		t.Fatalf("other: %d", got)
	}
}
