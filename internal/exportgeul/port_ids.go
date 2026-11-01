//ff:func feature=convert type=command control=iteration dimension=1
//ff:what portIDs — port each entity ID in order; the reports and document numbers of the documents written
package exportgeul

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/resolve"

	"github.com/geul-org/wikigeul/pkg/port"
)

// portIDs ports prefix+n for each n under issuer with prov (portOne). It returns the reports of
// all IDs (missing ones too, with no counts) and the document numbers
// (port.DocNumber) of the documents written. The first error stops it.
func portIDs(db *sql.DB, dir resolve.Dir, prefix string, nums []uint64, issuer []core.Word, c config, prov port.Provenance) ([]port.Report, []uint64, error) {
	var reps []port.Report
	var written []uint64
	for _, n := range nums {
		id := prefix + strconv.FormatUint(n, 10)
		rep, ok, err := portOne(db, dir, id, issuer, c, prov)
		if err != nil {
			return reps, written, fmt.Errorf("%s: %w", id, err)
		}
		reps = append(reps, rep)
		if ok {
			doc, _ := port.DocNumber(prefix[0], n) // BuildDoc already accepted this ID's number
			written = append(written, doc)
		}
	}
	return reps, written, nil
}
