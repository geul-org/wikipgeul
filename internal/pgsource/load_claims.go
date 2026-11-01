//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what loadClaims — the statements a query returns as (statement ID, property, type, value, rank) rows
package pgsource

import (
	"database/sql"

	"github.com/geul-org/wikigeul/pkg/source"
)

// loadClaims runs query (which selects statement ID, property, value type, datatype, value,
// value_json, rank and snak type for $1 = id) and returns the rows in order,
// each value exact (exactValue).
func loadClaims(db *sql.DB, query, id string) ([]source.Claim, error) {
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []source.Claim
	for rows.Next() {
		var c source.Claim
		var exact string
		if err := rows.Scan(&c.ID, &c.Property, &c.Type, &c.Datatype, &c.Value, &exact, &c.Rank, &c.SnakType); err != nil {
			return nil, err
		}
		c.Value = source.ExactValue(c.Value, exact)
		out = append(out, c)
	}
	return out, rows.Err()
}
