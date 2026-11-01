//ff:func feature=convert type=loader control=sequence
//ff:what LoadEntity — an entity's labels, descriptions, aliases and claims from a loaded dump; false when the ID is not in it
package pgsource

import (
	"database/sql"
	"fmt"

	"github.com/geul-org/wikigeul/pkg/source"
)

// LoadEntity reads the entity id from the dump s (queriesFor): its labels,
// descriptions and aliases in every language, all its claims and its sitelinks. It is false, with no
// error, when the dump has no such ID. An invalid schema name is an error.
func LoadEntity(db *sql.DB, s Store, id string) (source.Entity, bool, error) {
	if !validSchema(s.Schema) {
		return source.Entity{}, false, fmt.Errorf("wikigeul: schema %q is not a plain identifier", s.Schema)
	}
	return loadEntityWith(db, queriesFor(s), id)
}
