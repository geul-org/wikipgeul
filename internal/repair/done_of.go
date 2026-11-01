//ff:func feature=repair type=reader control=iteration dimension=1
//ff:what doneOf — the occurrence each of some hashes was processed at in an earlier batch or run (repair_reference_snak_done), one key lookup per hash
package repair

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/geul-org/wikipgeul/internal/load"
)

// doneOf returns, for those of hashes already processed, the occurrence they
// were processed at, read in tx: one primary-key lookup per hash (LATERAL
// with OFFSET 0, as in fetchStored; the done table grows to every hash).
func doneOf(tx *sql.Tx, c Config, hashes []string) (map[string]load.Occurrence, error) {
	rows, err := tx.Query("SELECT d.hash, d.occ_statement_id, d.occ_ordinal FROM unnest($2::text[]) AS h(hash) CROSS JOIN LATERAL "+
		"(SELECT * FROM "+pq.QuoteIdentifier(c.Schema)+".repair_reference_snak_done WHERE hash = h.hash AND dump_id = $1 OFFSET 0) d",
		c.DumpID, pq.Array(hashes))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]load.Occurrence{}
	for rows.Next() {
		var h string
		var o load.Occurrence
		if err := rows.Scan(&h, &o.Statement, &o.Ordinal); err != nil {
			return nil, err
		}
		out[h] = o
	}
	return out, rows.Err()
}
