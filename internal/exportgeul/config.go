//ff:type feature=convert type=model
//ff:what config — an export-geul run: what every conversion run shares, the database, whether the dump keeps each reference occurrence's snak order, and the output form
package exportgeul

import "github.com/geul-org/wikigeul/pkg/pipeline"

// config is one wikipg export-geul run: the shared run settings (Schema and
// Dump name the loaded dump, sql/wikidata.sql), the PostgreSQL DSN (empty:
// lib/pq reads the PG* environment variables) and Out, "pack" (default: the
// Issuer's pack set under Root written as one stream, runPack) or "dir" (the
// directory tree, for tests). RefOrder is not a flag: Run sets it from the
// loaded dump (pgsource.Dump.RefOrder, reference_snak_order present).
type config struct {
	pipeline.Config
	DSN      string
	RefOrder bool
	Out      string
}
