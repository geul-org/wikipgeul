//ff:func feature=convert type=loader control=sequence
//ff:what LoadDump — the provenance row of a loaded dump and whether its schema has reference_snak_order; an error for a bad schema name or a missing dump
package pgsource

import (
	"database/sql"
	"fmt"
)

// LoadDump reads the dump row of s (URL, file SHA-256, download time, loader
// and its commit) and whether the schema has the table reference_snak_order
// (to_regclass, in the same query). An invalid schema name or a dump_id with
// no row is an error.
func LoadDump(db *sql.DB, s Store) (Dump, error) {
	if !validSchema(s.Schema) {
		return Dump{}, fmt.Errorf("wikigeul: schema %q is not a plain identifier", s.Schema)
	}
	var d Dump
	err := db.QueryRow(fmt.Sprintf("SELECT url, sha256, downloaded_at, loader, loader_commit, to_regclass('%s.reference_snak_order') IS NOT NULL FROM %s.dump WHERE dump_id = $1", s.Schema, s.Schema), s.Dump).
		Scan(&d.URL, &d.SHA256, &d.DownloadedAt, &d.Loader, &d.LoaderCommit, &d.RefOrder)
	return d, err
}
