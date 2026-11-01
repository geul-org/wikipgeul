//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what EntityNumbers — the next entity numbers of one kind in the dump, ascending, after a number and up to a bound, at most limit
package pgsource

import (
	"database/sql"
	"fmt"
)

// EntityNumbers returns up to limit numbers of the entities of kind ('Q',
// 'P', 'L') that dump s holds, in ascending order, greater than after and at
// most to. It reads the (dump_id, kind, num) index, so a whole kind is walked
// page by page without holding it. Entities the dump does not have (deleted,
// redirected) never appear. An invalid schema name is an error.
func EntityNumbers(db *sql.DB, s Store, kind byte, after, to uint64, limit int) ([]uint64, error) {
	if !validSchema(s.Schema) {
		return nil, fmt.Errorf("wikigeul: schema %q is not a plain identifier", s.Schema)
	}
	q := fmt.Sprintf("SELECT num FROM %s.entity WHERE dump_id = %d AND kind = $1 AND num > $2 AND num <= $3 ORDER BY num LIMIT $4", s.Schema, s.Dump)
	rows, err := db.Query(q, string(kind), int64(after), int64(min(to, 1<<62)), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []uint64
	for rows.Next() {
		var n int64
		if err := rows.Scan(&n); err != nil {
			return nil, err
		}
		out = append(out, uint64(n))
	}
	return out, rows.Err()
}
