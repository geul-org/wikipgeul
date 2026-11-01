//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what loadSitelinks — the sitelinks of an entity a query returns as (site, title, title_json, badges) rows
package pgsource

import (
	"database/sql"

	"github.com/geul-org/wikigeul/pkg/source"
	"github.com/lib/pq"
)

// loadSitelinks runs query (which selects site, title, title_json and badges
// for $1 = id) and returns the rows in order, each title exact (exactValue).
func loadSitelinks(db *sql.DB, query, id string) ([]source.Sitelink, error) {
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []source.Sitelink
	for rows.Next() {
		var s source.Sitelink
		var exact string
		if err := rows.Scan(&s.Site, &s.Title, &exact, pq.Array(&s.Badges)); err != nil {
			return nil, err
		}
		s.Title = source.ExactValue(s.Title, exact)
		out = append(out, s)
	}
	return out, rows.Err()
}
