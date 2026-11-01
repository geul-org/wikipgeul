//ff:func feature=load type=model control=iteration dimension=2
//ff:what Rows.addTerms — label, description and alias rows of an entity in the kept languages, aliases in order
package load

import "github.com/geul-org/wikipgeul/internal/dump"

// addTerms adds the labels, descriptions and aliases of e in the languages c
// keeps; each alias with its position in its language's list.
func (r *Rows) addTerms(c Config, line int64, e dump.Entity) {
	for lang, t := range e.Labels {
		if c.keepLang(lang) {
			r.addText("label", t.Value, e.ID, c.DumpID, lang)
		}
	}
	for lang, t := range e.Descriptions {
		if c.keepLang(lang) {
			r.addText("description", t.Value, e.ID, c.DumpID, lang)
		}
	}
	for lang, ts := range e.Aliases {
		if !c.keepLang(lang) {
			continue
		}
		for i, t := range ts {
			r.addText("alias", t.Value, e.ID, c.DumpID, lang, i)
		}
	}
}
