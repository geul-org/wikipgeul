//ff:func feature=convert type=encoder control=sequence
//ff:what docName — the name of a wikidata.geul.org document from its number: the lower-case kind letter and the entity number (q42, p31)
package publish

import (
	"fmt"
	"strconv"
)

// docName names the document number n (port.DocNumber: the kind letter in
// the top byte, the entity number below) as the lower-case letter followed by
// the number: 0x5100_0000_0000_002A is q42. A top byte that is not an
// upper-case ASCII letter is an error.
func docName(n uint64) (string, error) {
	letter := n >> 56
	if letter < 'A' || letter > 'Z' {
		return "", fmt.Errorf("wikigeul: document 0x%016x has no kind letter", n)
	}
	return string(rune(letter+'a'-'A')) + strconv.FormatUint(n&(1<<56-1), 10), nil
}
