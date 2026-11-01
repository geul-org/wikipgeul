//ff:func feature=convert type=encoder control=sequence
//ff:what TestDocName — Q42 and P31 document numbers name q42 and p31; a number without a letter fails
package publish

import "testing"

func TestDocName(t *testing.T) {
	if s, err := docName(0x5100_0000_0000_002A); err != nil || s != "q42" {
		t.Fatalf("q42: %q %v", s, err)
	}
	if s, err := docName(0x5000_0000_0000_001F); err != nil || s != "p31" {
		t.Fatalf("p31: %q %v", s, err)
	}
	if _, err := docName(42); err == nil {
		t.Fatal("no letter: want error")
	}
}
