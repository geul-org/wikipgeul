//ff:type feature=convert type=model
//ff:what Dump — a loaded dump's provenance from the dump table: URL, file SHA-256, when it was downloaded, the loader and its commit; and whether its schema keeps each reference occurrence's snak order
package pgsource

import "time"

// Dump is the provenance of a loaded dump (sql/wikidata.sql dump table):
// where it was downloaded from, the SHA-256 of the file in hex, when the file
// was received, the loader that filled the tables (its Go package path) and
// the loader's git commit. RefOrder tells whether the schema has
// reference_snak_order (wikipgeul build-reference-order: the snak order of each
// reference occurrence that differs from reference_snak); then references
// are read in each occurrence's own order (Store.RefOrder).
type Dump struct {
	URL          string
	SHA256       string
	DownloadedAt time.Time
	Loader       string
	LoaderCommit string
	RefOrder     bool
}
