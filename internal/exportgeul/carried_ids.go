//ff:func feature=convert type=helper control=sequence
//ff:what carriedIDs — the LOCAL ID assignment for an entity's document: carried over from its published version and version list
package exportgeul

import (
	"strconv"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/delta"
	"github.com/geul-org/gogeul/pkg/resolve"

	"github.com/geul-org/wikigeul/pkg/port"
)

// carriedIDs reads the current version of id's document under issuer and its
// NEXT-LOCAL-ID (resolve.Dir.Current) and returns the assignment that keeps
// the IDs of unchanged statements (delta.NewLocalIDs). A document never
// published starts from 1.
func carriedIDs(dir resolve.Dir, issuer []core.Word, id string) (*delta.LocalIDs, error) {
	n, err := strconv.ParseUint(id[1:], 10, 64)
	if err != nil {
		return nil, err
	}
	num, err := port.DocNumber(id[0], n)
	if err != nil {
		return nil, err
	}
	prev, next, err := dir.Current(issuer, []core.Word{core.Word(num >> 48), core.Word(num >> 32), core.Word(num >> 16), core.Word(num)})
	if err != nil {
		return nil, err
	}
	return delta.NewLocalIDs(prev, next)
}
