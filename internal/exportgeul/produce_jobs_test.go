//ff:func feature=convert type=command control=sequence
//ff:what TestProduceJobs — every property then the Q range in order and numbered; from a checkpoint in P or in Q; the ID list after the checkpoint; a bad schema leaves an error
package exportgeul

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestProduceJobs(t *testing.T) {
	db, schema := fixtureDB(t)
	c := config{Config: pipeline.Config{Schema: schema, Dump: 1, From: 2, To: 40}}
	collect := func(c config, p pipeline.PackProgress, ids []pipeline.PackJob) ([]string, error) {
		jobs := make(chan pipeline.PackJob, 100)
		var err error
		produceJobs(db, c, p, ids, jobs, &err)
		var out []string
		for j := range jobs {
			out = append(out, string(j.Kind)+string(rune('0'+j.Seq)))
		}
		return out, err
	}
	all, err := collect(c, pipeline.PackProgress{}, nil)
	if err != nil || len(all) != 6 || all[0] != "P0" || all[2] != "Q2" || all[5] != "Q5" {
		t.Fatalf("%v %v", all, err)
	}
	if got, _ := collect(c, pipeline.PackProgress{Kind: 'P', Num: 31}, nil); len(got) != 5 || got[0] != "P0" {
		t.Fatalf("after P31: %v", got)
	}
	if got, _ := collect(c, pipeline.PackProgress{Kind: 'Q', Num: 3}, nil); len(got) != 2 {
		t.Fatalf("after Q3: %v", got)
	}
	if got, _ := collect(c, pipeline.PackProgress{Kind: 'Q', Num: 1}, []pipeline.PackJob{{Kind: 'P', Num: 1}, {Kind: 'Q', Num: 1}, {Kind: 'Q', Num: 9}}); len(got) != 1 {
		t.Fatalf("ID list: %v", got)
	}
	if _, err := collect(config{Config: pipeline.Config{Schema: "1x", Dump: 1}}, pipeline.PackProgress{}, nil); err == nil {
		t.Fatal("bad schema: want error")
	}
}
