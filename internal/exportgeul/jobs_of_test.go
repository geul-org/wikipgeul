//ff:func feature=convert type=helper control=sequence
//ff:what TestJobsOf — one job of the kind per number
package exportgeul

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestJobsOf(t *testing.T) {
	js := jobsOf('Q', []uint64{3, 9})
	if len(js) != 2 || js[1] != (pipeline.PackJob{Kind: 'Q', Num: 9}) {
		t.Fatal(js)
	}
}
