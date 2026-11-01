//ff:func feature=cli type=command control=selection
//ff:what exitCode — the process exit code of a failed run: 3 when a strict export-geul run stopped at an anomaly, else 1
package main

import (
	"errors"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

// exitCode says why wikipgeul failed: export-geul -strict stopped at an anomaly
// (3, as wikigeul exits), any other error (1).
func exitCode(err error) int {
	switch {
	case errors.As(err, &pipeline.AnomalyError{}):
		return 3
	default:
		return 1
	}
}
