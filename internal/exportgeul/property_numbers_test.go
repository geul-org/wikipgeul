//ff:func feature=convert type=helper control=sequence
//ff:what TestPropertyNumbers — distinct, ascending, malformed IDs dropped
package exportgeul

import (
	"slices"
	"testing"

	"github.com/geul-org/wikigeul/pkg/port"
)

func TestPropertyNumbers(t *testing.T) {
	got := propertyNumbers([]port.Report{{Properties: []string{"P31", "P5", "Px"}}, {Properties: []string{"P31"}}})
	if !slices.Equal(got, []uint64{5, 31}) {
		t.Fatalf("got %v", got)
	}
}
