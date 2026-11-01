//ff:func feature=convert type=loader control=sequence
//ff:what TestSampleJobs — the first Samples Q entities of the range or of the ID list; a bad schema fails
package exportgeul

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestSampleJobs(t *testing.T) {
	db, schema := fixtureDB(t)
	c := packTestConfig(t, schema)
	c.Samples = 3
	if js, err := sampleJobs(db, c, nil); err != nil || len(js) != 3 || js[0].Kind != 'Q' || js[0].Num != 1 {
		t.Fatalf("%+v %v", js, err)
	}
	ids := []pipeline.PackJob{{Kind: 'P', Num: 31}, {Kind: 'Q', Num: 42}}
	if js, err := sampleJobs(db, c, ids); err != nil || len(js) != 1 || js[0].Num != 42 {
		t.Fatalf("ID list: %+v %v", js, err)
	}
	if _, err := sampleJobs(db, packTestConfig(t, "1x"), nil); err == nil {
		t.Fatal("bad schema: want error")
	}
}
