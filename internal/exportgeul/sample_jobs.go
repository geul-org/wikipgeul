//ff:func feature=convert type=loader control=sequence
//ff:what sampleJobs — the entities a pack-mode run trains its zstd dictionary on: the first Q entities of the ID list, else of the dump from Q<From>
package exportgeul

import (
	"database/sql"

	"github.com/geul-org/wikigeul/pkg/pipeline"

	"github.com/geul-org/wikipgeul/internal/pgsource"
)

// sampleJobs lists the first c.Samples Q entities of the run: the ID list's
// (pipeline.FirstQ) when there is one, else the dump's from Q<From>.
func sampleJobs(db *sql.DB, c config, ids []pipeline.PackJob) ([]pipeline.PackJob, error) {
	if ids != nil {
		return pipeline.FirstQ(ids, c.Samples), nil
	}
	nums, err := pgsource.EntityNumbers(db, pgsource.Store{Schema: c.Schema, Dump: c.Dump}, 'Q', c.From-1, c.To, c.Samples)
	if err != nil {
		return nil, err
	}
	return jobsOf('Q', nums), nil
}
