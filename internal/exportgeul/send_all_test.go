//ff:func feature=convert type=helper control=sequence
//ff:what TestSendAll — every job in order
package exportgeul

import (
	"testing"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

func TestSendAll(t *testing.T) {
	var got []uint64
	sendAll([]pipeline.PackJob{{Num: 1}, {Num: 2}}, func(j pipeline.PackJob) { got = append(got, j.Num) })
	if len(got) != 2 || got[1] != 2 {
		t.Fatal(got)
	}
}
