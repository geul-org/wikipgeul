//ff:func feature=load type=model control=iteration dimension=1
//ff:what Rows.addSitelinks — one sitelink row per site of an entity, badges as a text array
package load

import (
	"github.com/lib/pq"

	"github.com/geul-org/wikipgeul/internal/dump"
)

// addSitelinks adds a sitelink row for each site of e: site, page title
// (without NUL bytes), badges (an empty array when there are none) and
// title_json, the exact title when it had NUL bytes (nulJSON).
func (r *Rows) addSitelinks(c Config, line int64, e dump.Entity) {
	for site, s := range e.Sitelinks {
		badges := s.Badges
		if badges == nil {
			badges = []string{}
		}
		title, _ := cleanText(s.Title)
		r.add("sitelink", e.ID, c.DumpID, site, title, pq.Array(badges), nulJSON(s.Title))
	}
}
