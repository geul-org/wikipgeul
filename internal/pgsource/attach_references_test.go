//ff:func feature=convert type=loader control=sequence
//ff:what TestAttachReferences — references go to their claim in order; one of an unlisted claim is dropped
package pgsource

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/source"
)

func TestAttachReferences(t *testing.T) {
	cs := []source.Claim{{ID: 1}, {ID: 2}}
	attachReferences(cs, []source.Reference{{StatementID: 2, Hash: "a"}, {StatementID: 9, Hash: "x"}, {StatementID: 2, Hash: "b"}})
	if len(cs[0].References) != 0 || len(cs[1].References) != 2 || cs[1].References[0].Hash != "a" || cs[1].References[1].Hash != "b" {
		t.Fatalf("got %+v", cs)
	}
}
