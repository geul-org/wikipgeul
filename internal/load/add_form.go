//ff:func feature=load type=model control=iteration dimension=1
//ff:what Rows.addForm — one lexeme form: its entity row (type form), form row with grammatical features, representations and statements; the next key number
package load

import (
	"github.com/lib/pq"

	"github.com/geul-org/wikipgeul/internal/dump"
)

// addForm adds form f, the ord-th form of lexeme (number num): an entity row
// of type form (kind L, the lexeme's number) so its statements have an
// entity to point at, the form row with its grammatical features, a
// representation per language variant (not filtered, like lemmas) and its
// statements from key number next (addStatements). It returns the next
// number and whether keys remain.
func (r *Rows) addForm(c Config, line int64, lexeme string, num int64, ord int, f dump.Form, next int, refs *RefSet) (int, bool) {
	features := f.GrammaticalFeatures
	if features == nil {
		features = []string{}
	}
	r.add("entity", f.ID, line, "L", num, "form", nil, nil, nil, c.DumpID)
	r.add("form", f.ID, lexeme, c.DumpID, ord, pq.Array(features))
	for lang, t := range f.Representations {
		r.addText("form_representation", t.Value, f.ID, c.DumpID, lang)
	}
	return r.addStatements(c, line, f.ID, f.Claims, next, refs)
}
