//ff:func feature=convert type=helper control=sequence
//ff:what TestPageKind — every entity of a kind in the bound, in order; a bad schema fails
package exportgeul

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestPageKind(t *testing.T) {
	db, schema := fixtureDB(t)
	c := config{Config: pipeline.Config{Schema: schema, Dump: 1}}
	var got []uint64
	if err := pageKind(db, c, 'Q', 2, 42, func(j pipeline.PackJob) { got = append(got, j.Num) }); err != nil || len(got) != 4 || got[0] != 3 || got[3] != 42 {
		t.Fatalf("%v %v", got, err)
	}
	if err := pageKind(db, config{Config: pipeline.Config{Schema: "1x", Dump: 1}}, 'Q', 0, 9, func(pipeline.PackJob) {}); err == nil {
		t.Fatal("bad schema: want error")
	}
}
