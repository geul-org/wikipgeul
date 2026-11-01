//ff:func feature=convert type=command control=sequence
//ff:what portOne — load, build and store one entity's document as version 1; false when the entity was never loaded
package exportgeul

import (
	"database/sql"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/resolve"

	"github.com/geul-org/wikigeul/pkg/port"
	"github.com/geul-org/wikipgeul/internal/pgsource"
)

// portOne ports entity id: loads it from c's dump (pgsource.LoadEntity, each
// reference occurrence in its own snak order when c.RefOrder), builds
// its document under issuer with prov and the LOCAL IDs carried over from
// its published version (carriedIDs, port.BuildDoc, languages by the run's
// registry c.Langs) and stores it in dir as version 1 over any
// earlier one (resolve.Dir.Overwrite: development, no version per run). It is
// false when the entity is not in the dump, and its report then says so
// (Missing).
func portOne(db *sql.DB, dir resolve.Dir, id string, issuer []core.Word, c config, prov port.Provenance) (port.Report, bool, error) {
	e, ok, err := pgsource.LoadEntity(db, pgsource.Store{Schema: c.Schema, Dump: c.Dump, RefOrder: c.RefOrder}, id)
	if err != nil || !ok {
		return port.Report{ID: id, Missing: err == nil}, false, err
	}
	ids, err := carriedIDs(dir, issuer, id)
	if err != nil {
		return port.Report{ID: id}, false, err
	}
	b, rep, err := port.BuildDoc(e, issuer, prov, ids, c.Langs)
	if err != nil {
		return rep, false, err
	}
	if _, err := dir.Overwrite(b); err != nil {
		return rep, false, err
	}
	return rep, true, nil
}
