//ff:func feature=convert type=command control=selection
//ff:what produceJobs — send the run's entities after the checkpoint in output order: the ID list, or every property then Q<From>…Q<To> the dump holds
package exportgeul

import (
	"database/sql"

	"github.com/geul-org/wikigeul/pkg/pipeline"
)

// produceJobs numbers and sends the entities still to convert on jobs and
// closes it: the ID list ids when given, else every property of the dump
// then the dump's Q entities in [From, To] — properties first because their
// documents (0x50…) sort before items' (0x51…). Entities up to the
// checkpoint's (p.Kind, p.Num) are skipped. An error is left in *errp.
func produceJobs(db *sql.DB, c config, p pipeline.PackProgress, ids []pipeline.PackJob, jobs chan<- pipeline.PackJob, errp *error) {
	defer close(jobs)
	var seq int64
	send := func(j pipeline.PackJob) {
		if pipeline.JobAfter(j.Kind, j.Num, p.Kind, p.Num) {
			j.Seq, seq = seq, seq+1
			jobs <- j
		}
	}
	switch {
	case ids != nil:
		sendAll(ids, send)
	case p.Kind == 'Q':
		*errp = pageKind(db, c, 'Q', max(p.Num, c.From-1), c.To, send)
	default:
		*errp = pageKind(db, c, 'P', p.Num, ^uint64(0)>>8, send)
		if *errp == nil {
			*errp = pageKind(db, c, 'Q', c.From-1, c.To, send)
		}
	}
}
