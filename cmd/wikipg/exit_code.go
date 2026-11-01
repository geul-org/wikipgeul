//ff:func feature=cli type=command control=selection
//ff:what exitCode — the process exit code of a failed run: 3 when a strict export-geul run stopped at an anomaly, 4 when the memory guard stopped it (telling the command that resumes), else 1
package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/geul-org/wikigeul/pkg/memguard"
	"github.com/geul-org/wikigeul/pkg/pipeline"
)

// exitCode says why wikipg failed: export-geul -strict stopped at an anomaly
// (3, as wikigeul exits), the memory guard stopped export-geul (4, the resume
// command from args written to w), any other error (1).
func exitCode(err error, args []string, w io.Writer) int {
	switch {
	case errors.As(err, &pipeline.AnomalyError{}):
		return 3
	case errors.As(err, &memguard.StopError{}):
		fmt.Fprintln(w, "resume:", memguard.ResumeCommand(args))
		return memguard.ExitCode
	default:
		return 1
	}
}
