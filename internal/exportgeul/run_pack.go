//ff:func feature=convert type=command control=iteration dimension=1
//ff:what runPack — pack mode: convert the run's entities with parallel workers, each with its own connection, into the Issuer's pack set (zstd), checkpointing and resuming
package exportgeul

import (
	"database/sql"
	"database/sql/driver"
	"path/filepath"
	"time"

	"github.com/geul-org/wikigeul/pkg/pipeline"
	"github.com/geul-org/wikigeul/pkg/port"
)

// runPack writes the pack set of wikidata.geul.org under
// <root>/.well-known/geul/0x00010000 as one pack stream (pipeline.WritePack):
// the entities of -ids, or every property then Q<From>…Q<To> the dump holds
// (produceJobs), each converted by one of c.Workers goroutines with its own
// connection. No documents.geul is written: at Wikidata's size the name index
// would be one 4 GB document, and a document's name follows from its ID
// (port.DocNumber). db (JIT off) lists the entities and builds the dictionary
// samples.
func runPack(c config, conn driver.Connector, db *sql.DB, prov port.Provenance) error {
	if c.RecordedAt == 0 && !c.Resume {
		c.RecordedAt = time.Now().Unix()
	}
	ids, idsHash, err := pipeline.ReadIDs(c.IDs)
	if err != nil {
		return err
	}
	set := filepath.Join(c.Root, ".well-known", "geul", "0x00010000")
	samples := func() ([]pipeline.PackJob, error) { return sampleJobs(db, c, ids) }
	s, p, err := pipeline.OpenPackStream(set, c.Config, prov, jobLoader(db, c), samples, idsHash)
	if err != nil {
		return err
	}
	loaders := make([]pipeline.Loader, c.Workers)
	for i := range loaders {
		loaders[i] = jobLoader(openDB(conn), c)
	}
	return pipeline.WritePack(s, c.Config, p, prov, loaders, func(jobs chan<- pipeline.PackJob, errp *error) {
		produceJobs(db, c, p, ids, jobs, errp)
	})
}
