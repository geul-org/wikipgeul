//ff:func feature=convert type=command control=sequence
//ff:what Run — wikipgeul export-geul: convert the loaded dump into GEUL documents, as a pack set or as a directory tree with the name index and the report
package exportgeul

import (
	"database/sql"
	"fmt"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/resolve"
	"github.com/geul-org/wikigeul/pkg/pipeline"
	"github.com/lib/pq"

	"github.com/geul-org/wikipgeul/internal/pgsource"
)

// Run reads the loaded dump's provenance (pgsource.LoadDump, provenanceOf)
// and, in directory mode, ports Q<from>…Q<to>, then every property those
// documents use, all under the one Issuer wikidata.geul.org (0x00010000;
// document IDs by port.DocNumber), writes its documents.geul and, when asked,
// the report; in pack mode runPack streams the pack set. A DSN lib/pq cannot
// parse fails before any work (pq.NewConnector). It keeps one connection with
// JIT off: the planner overestimates the reference join (statement_reference
// ⋈ reference_snak) and JIT-compiles it, about 160 ms per entity against 0.5
// ms without (2026-10-27). When the schema has reference_snak_order
// (pgsource.Dump.RefOrder) every reference occurrence keeps its own snak
// order; the run says which on stdout.
func Run(args []string) error {
	c, err := parseFlags(args)
	if err != nil {
		return err
	}
	conn, err := pq.NewConnector(c.DSN)
	if err != nil {
		return err
	}
	db := sql.OpenDB(conn)
	defer db.Close()
	db.SetMaxOpenConns(1)
	db.Exec("SET jit = off") // speed only: an error leaves JIT on
	d, err := pgsource.LoadDump(db, pgsource.Store{Schema: c.Schema, Dump: c.Dump})
	if err != nil {
		return fmt.Errorf("dump %s/%d: %w", c.Schema, c.Dump, err)
	}
	c.RefOrder = d.RefOrder
	if err := pipeline.LoadLanguages(&c.Config); err != nil {
		return err
	}
	fmt.Printf("language registry: %s, %d tags\n", c.LangDoc, len(c.Langs))
	fmt.Printf("reference snak order: per occurrence (reference_snak_order) %v\n", c.RefOrder)
	prov, err := provenanceOf(d, c.Commit, c.LibCommit)
	if err != nil {
		return err
	}
	if c.Out == "pack" {
		return runPack(c, conn, db, prov)
	}
	dir := resolve.Dir{Root: c.Root}
	issuer := []core.Word{1, 0}
	qReps, qDocs, err := portIDs(db, dir, "Q", rangeNumbers(c.From, c.To), issuer, c, prov)
	if err != nil {
		return err
	}
	pReps, pDocs, err := portIDs(db, dir, "P", propertyNumbers(qReps), issuer, c, prov)
	if err != nil {
		return err
	}
	if err := writeIndex(c.Root, issuer, append(qDocs, pDocs...)); err != nil {
		return err
	}
	fmt.Printf("Q: %d documents of %d IDs, P: %d documents of %d IDs\n", len(qDocs), len(qReps), len(pDocs), len(pReps))
	if c.Report == "" {
		return nil
	}
	return pipeline.WriteReport(c.Report, pipeline.RunLine(c.Config), append(qReps, pReps...))
}
