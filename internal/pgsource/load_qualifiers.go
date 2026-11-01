//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what loadQualifiers — the qualifiers of an entity's claims a query returns as (statement ID, property, datatype, value) rows
package pgsource

import (
	"database/sql"

	"github.com/geul-org/wikigeul/pkg/source"
)

// loadQualifiers runs query (which selects statement ID, property, value type, datatype,
// value, value_json and snak type of the qualifiers of $1 = id's claims) and
// returns the rows in order, each value exact (exactValue).
func loadQualifiers(db *sql.DB, query, id string) ([]source.Qualifier, error) {
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []source.Qualifier
	for rows.Next() {
		var q source.Qualifier
		var exact string
		if err := rows.Scan(&q.StatementID, &q.Property, &q.Type, &q.Datatype, &q.Value, &exact, &q.SnakType); err != nil {
			return nil, err
		}
		q.Value = source.ExactValue(q.Value, exact)
		out = append(out, q)
	}
	return out, rows.Err()
}
