//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what loadReferences — the references of an entity's claims a query returns as one row per snak, grouped into references in order
package pgsource

import (
	"database/sql"

	"github.com/geul-org/wikigeul/pkg/source"
)

// loadReferences runs query (which selects statement ID, reference ordinal,
// hash, then property, value type, datatype, value, value_json and snak type
// of each snak of the references of $1 = id's claims, ordered by statement,
// reference ordinal and snak ordinal) and groups consecutive rows of the same
// statement and ordinal into one Reference, each value exact (exactValue).
func loadReferences(db *sql.DB, query, id string) ([]source.Reference, error) {
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []source.Reference
	for rows.Next() {
		var r source.Reference
		var q source.Qualifier
		var exact string
		if err := rows.Scan(&r.StatementID, &r.Ordinal, &r.Hash, &q.Property, &q.Type, &q.Datatype, &q.Value, &exact, &q.SnakType); err != nil {
			return nil, err
		}
		q.StatementID, q.Value = r.StatementID, source.ExactValue(q.Value, exact)
		if n := len(out); n == 0 || out[n-1].StatementID != r.StatementID || out[n-1].Ordinal != r.Ordinal {
			out = append(out, r)
		}
		out[len(out)-1].Snaks = append(out[len(out)-1].Snaks, q)
	}
	return out, rows.Err()
}
