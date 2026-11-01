//ff:func feature=convert type=loader control=sequence
//ff:what TestAttachQualifiers — qualifiers go to their claim in order; one of an unlisted claim is dropped
package pgsource

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/source"
)

func TestAttachQualifiers(t *testing.T) {
	cs := []source.Claim{{ID: 1}, {ID: 2}}
	attachQualifiers(cs, []source.Qualifier{{StatementID: 2, Property: "P580"}, {StatementID: 9, Property: "P1"}, {StatementID: 2, Property: "P582"}})
	if len(cs[0].Qualifiers) != 0 || len(cs[1].Qualifiers) != 2 || cs[1].Qualifiers[0].Property != "P580" || cs[1].Qualifiers[1].Property != "P582" {
		t.Fatalf("got %+v", cs)
	}
}
