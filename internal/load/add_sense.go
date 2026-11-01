//ff:func feature=load type=model control=iteration dimension=1
//ff:what Rows.addSense — one lexeme sense: its entity row (type sense), sense row, glosses in the kept languages and statements; the next key number
package load

import "github.com/geul-org/wikipgeul/internal/dump"

// addSense adds sense s, the ord-th sense of lexeme (number num): an entity
// row of type sense, the sense row, a gloss per language that c keeps
// (glosses are translations of the meaning, filtered like descriptions) and
// its statements from key number next (addStatements). It returns the next
// number and whether keys remain.
func (r *Rows) addSense(c Config, line int64, lexeme string, num int64, ord int, s dump.Sense, next int, refs *RefSet) (int, bool) {
	r.add("entity", s.ID, line, "L", num, "sense", nil, nil, nil, c.DumpID)
	r.add("sense", s.ID, lexeme, c.DumpID, ord)
	for lang, t := range s.Glosses {
		if c.keepLang(lang) {
			r.addText("sense_gloss", t.Value, s.ID, c.DumpID, lang)
		}
	}
	return r.addStatements(c, line, s.ID, s.Claims, next, refs)
}
