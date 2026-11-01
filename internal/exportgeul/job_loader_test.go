//ff:func feature=convert type=loader control=sequence
//ff:what TestJobLoader — an entity of the dump loads by kind and number; a missing one is false; a bad schema fails
package exportgeul

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestJobLoader(t *testing.T) {
	db, schema := fixtureDB(t)
	load := jobLoader(db, packTestConfig(t, schema))
	if e, ok, err := load(pipeline.PackJob{Kind: 'Q', Num: 42}); err != nil || !ok || e.ID != "Q42" || len(e.Labels) == 0 {
		t.Fatalf("%+v %v %v", e, ok, err)
	}
	if _, ok, err := load(pipeline.PackJob{Kind: 'Q', Num: 6}); ok || err != nil {
		t.Fatalf("missing: %v %v", ok, err)
	}
	if _, _, err := jobLoader(db, packTestConfig(t, "1x"))(pipeline.PackJob{Kind: 'Q', Num: 42}); err == nil {
		t.Fatal("bad schema: want error")
	}
}
