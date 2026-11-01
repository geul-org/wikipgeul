//ff:func feature=convert type=encoder control=sequence
//ff:what TestNameIndex — sorted, deduplicated BINDs that resolve.ParseDocumentNames reads back; a bad Issuer and a number without a kind letter fail
package publish

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/resolve"
)

func TestNameIndex(t *testing.T) {
	b, err := NameIndex([]core.Word{1, 0}, []uint64{0x5100_0000_0000_002A, 0x5000_0000_0000_001F, 0x5100_0000_0000_002A})
	if err != nil {
		t.Fatal(err)
	}
	names, err := resolve.ParseDocumentNames(b)
	if err != nil || len(names) != 2 || names["q42"] != (resolve.DocKey{Issuer: "0x00010000", Doc: "0x510000000000002a"}) || names["p31"].Doc != "0x500000000000001f" {
		t.Fatalf("got %v %v", names, err)
	}
	if _, err := NameIndex([]core.Word{0, 1}, []uint64{0x5100_0000_0000_0001}); err == nil {
		t.Fatal("bad Issuer: want error")
	}
	if _, err := NameIndex([]core.Word{1, 0}, []uint64{1}); err == nil {
		t.Fatal("no kind letter: want error")
	}
}
