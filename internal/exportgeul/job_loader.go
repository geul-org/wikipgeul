//ff:func feature=convert type=loader control=sequence
//ff:what jobLoader — the database's pipeline.Loader over one connection: a job's entity from the loaded dump
package exportgeul

import (
	"database/sql"
	"strconv"

	"github.com/geul-org/wikigeul/pkg/pipeline"
	"github.com/geul-org/wikigeul/pkg/source"

	"github.com/geul-org/wikipgeul/internal/pgsource"
)

// jobLoader reads a job's entity through db (pgsource.LoadEntity, each
// reference occurrence in its own snak order when c.RefOrder).
func jobLoader(db *sql.DB, c config) pipeline.Loader {
	store := pgsource.Store{Schema: c.Schema, Dump: c.Dump, RefOrder: c.RefOrder}
	return func(job pipeline.PackJob) (source.Entity, bool, error) {
		return pgsource.LoadEntity(db, store, string(job.Kind)+strconv.FormatUint(job.Num, 10))
	}
}
