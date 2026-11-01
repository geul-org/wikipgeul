//ff:func feature=load type=model control=sequence
//ff:what Rows.AddEntity — every row one dump entity gives: raw line, head, terms, sitelinks, statements with qualifiers and references, lexeme parts
package load

import "github.com/geul-org/wikipgeul/internal/dump"

// AddEntity adds the rows of entity e from dump line line (its trimmed JSON
// raw): the raw line when c.Raw (entity_raw), the head (entity), terms in the
// kept languages (addTerms), sitelinks (addSitelinks), statements
// (addStatements) and, for a lexeme, its lemmas, forms and senses
// (addLexeme). An ID the loader cannot split is recorded as a load_error
// and nothing else is added.
func (r *Rows) AddEntity(c Config, line int64, raw []byte, e dump.Entity, refs *RefSet) {
	kind, num, ok := dump.SplitID(e.ID)
	if !ok {
		r.Error(c.DumpID, line, e.ID, "entity", "entity ID does not split into a letter and a number", raw)
		return
	}
	if c.Raw {
		text, clean := cleanText(string(raw))
		if !clean {
			r.Error(c.DumpID, line, e.ID, "raw", "NUL byte removed from the raw line", nil)
		}
		r.add("entity_raw", e.ID, c.DumpID, line, text)
	}
	var lastrev any
	if e.LastRevID != 0 {
		lastrev = e.LastRevID
	}
	r.add("entity", e.ID, line, string(kind), num, e.Type, nullable(e.Datatype), lastrev, nullable(e.Modified), c.DumpID)
	r.addTerms(c, line, e)
	r.addSitelinks(c, line, e)
	if next, ok := r.addStatements(c, line, e.ID, e.Claims, 0, refs); ok {
		r.addLexeme(c, line, e, num, next, refs)
	}
}
