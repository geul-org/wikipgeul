//ff:func feature=convert type=command control=sequence
//ff:what packTestProv — the provenance of the fixture dump with the config's commits
package exportgeul

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/port"
)

// packTestProv builds c's provenance over fixtureDump.
func packTestProv(t *testing.T, c config) port.Provenance {
	p, err := provenanceOf(fixtureDump(), c.Commit, c.LibCommit)
	if err != nil {
		t.Fatal(err)
	}
	return p
}
