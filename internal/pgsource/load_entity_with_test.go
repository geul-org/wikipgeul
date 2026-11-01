//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what TestLoadEntityWith — Q42 loads; a missing ID is false; a failing existence, label, description, alias, claim, qualifier, reference or sitelink query errors
package pgsource

import "testing"

func TestLoadEntityWith(t *testing.T) {
	db, st := fixtureDB(t)
	q := queriesFor(st)
	if e, ok, err := loadEntityWith(db, q, "Q42"); err != nil || !ok || len(e.Labels) != 3 || len(e.Descriptions) != 2 || len(e.Sitelinks) != 7 {
		t.Fatalf("Q42: %v %v", ok, err)
	}
	if _, ok, err := loadEntityWith(db, q, "Q6"); ok || err != nil {
		t.Fatalf("Q6: %v %v", ok, err)
	}
	bad := "SELECT nope FROM nowhere"
	for i, b := range []entityQueries{
		{Exists: bad, Labels: q.Labels, Descriptions: q.Descriptions, Aliases: q.Aliases, Claims: q.Claims, Qualifiers: q.Qualifiers, References: q.References, Sitelinks: q.Sitelinks},
		{Exists: q.Exists, Labels: bad, Descriptions: q.Descriptions, Aliases: q.Aliases, Claims: q.Claims, Qualifiers: q.Qualifiers, References: q.References, Sitelinks: q.Sitelinks},
		{Exists: q.Exists, Labels: q.Labels, Descriptions: bad, Aliases: q.Aliases, Claims: q.Claims, Qualifiers: q.Qualifiers, References: q.References, Sitelinks: q.Sitelinks},
		{Exists: q.Exists, Labels: q.Labels, Descriptions: q.Descriptions, Aliases: bad, Claims: q.Claims, Qualifiers: q.Qualifiers, References: q.References, Sitelinks: q.Sitelinks},
		{Exists: q.Exists, Labels: q.Labels, Descriptions: q.Descriptions, Aliases: q.Aliases, Claims: bad, Qualifiers: q.Qualifiers, References: q.References, Sitelinks: q.Sitelinks},
		{Exists: q.Exists, Labels: q.Labels, Descriptions: q.Descriptions, Aliases: q.Aliases, Claims: q.Claims, Qualifiers: bad, References: q.References, Sitelinks: q.Sitelinks},
		{Exists: q.Exists, Labels: q.Labels, Descriptions: q.Descriptions, Aliases: q.Aliases, Claims: q.Claims, Qualifiers: q.Qualifiers, References: bad, Sitelinks: q.Sitelinks},
		{Exists: q.Exists, Labels: q.Labels, Descriptions: q.Descriptions, Aliases: q.Aliases, Claims: q.Claims, Qualifiers: q.Qualifiers, References: q.References, Sitelinks: bad},
	} {
		if _, _, err := loadEntityWith(db, b, "Q42"); err == nil {
			t.Fatalf("bad %d: want error", i)
		}
	}
}
