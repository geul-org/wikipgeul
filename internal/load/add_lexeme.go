//ff:func feature=load type=model control=iteration dimension=1
//ff:what Rows.addLexeme — a lexeme's own rows: lexeme (language, category), lemmas, then each form and sense with their statements continuing the key sequence
package load

import "github.com/geul-org/wikipgeul/internal/dump"

// addLexeme adds what only a lexeme (lexeme dump) has: the lexeme row with
// its language and lexical category (item IDs), its lemmas (every language
// variant: a lemma is the word itself, not a translation, so -langs does not
// filter it), then each form (addForm) and sense (addSense) in order. Their
// statements take the keys after the lexeme's own, from next. Another kind
// of entity adds nothing.
func (r *Rows) addLexeme(c Config, line int64, e dump.Entity, num int64, next int, refs *RefSet) {
	if e.Type != "lexeme" {
		return
	}
	r.add("lexeme", e.ID, c.DumpID, nullable(e.Language), nullable(e.LexicalCategory))
	for lang, t := range e.Lemmas {
		r.addText("lemma", t.Value, e.ID, c.DumpID, lang)
	}
	ok := true
	for ord, f := range e.Forms {
		if ok {
			next, ok = r.addForm(c, line, e.ID, num, ord, f, next, refs)
		}
	}
	for ord, s := range e.Senses {
		if ok {
			next, ok = r.addSense(c, line, e.ID, num, ord, s, next, refs)
		}
	}
}
