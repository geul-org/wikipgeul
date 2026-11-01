//ff:func feature=convert type=loader control=sequence
//ff:what loadEntityWith — an entity's revision, labels, descriptions, aliases, claims with their qualifiers and references by the given queries; false when the ID was never loaded
package pgsource

import (
	"database/sql"

	"github.com/geul-org/wikigeul/pkg/source"
)

// loadEntityWith reads the entity id by q: existence, datatype and revision, then labels, descriptions and
// aliases in every language, then claims with their qualifiers (attachQualifiers) and
// references (attachReferences), then sitelinks. It is false, with no error, when the entity does not
// exist.
func loadEntityWith(db *sql.DB, q entityQueries, id string) (source.Entity, bool, error) {
	e, ok, err := entityHead(db, q.Exists, id)
	if err != nil || !ok {
		return source.Entity{}, false, err
	}
	if e.Labels, err = loadTexts(db, q.Labels, id); err != nil {
		return source.Entity{}, false, err
	}
	if e.Descriptions, err = loadTexts(db, q.Descriptions, id); err != nil {
		return source.Entity{}, false, err
	}
	if e.Aliases, err = loadTexts(db, q.Aliases, id); err != nil {
		return source.Entity{}, false, err
	}
	if e.Claims, err = loadClaims(db, q.Claims, id); err != nil {
		return source.Entity{}, false, err
	}
	qs, err := loadQualifiers(db, q.Qualifiers, id)
	if err != nil {
		return source.Entity{}, false, err
	}
	attachQualifiers(e.Claims, qs)
	rs, err := loadReferences(db, q.References, id)
	if err != nil {
		return source.Entity{}, false, err
	}
	attachReferences(e.Claims, rs)
	if e.Sitelinks, err = loadSitelinks(db, q.Sitelinks, id); err != nil {
		return source.Entity{}, false, err
	}
	return e, true, nil
}
