//ff:func feature=load type=model control=sequence
//ff:what sampleLexeme — test fixture: lexeme L4 "windsurf" with two lemmas, a form with features and a claim, a sense with two glosses and a claim
package load

import (
	"encoding/json"

	"github.com/geul-org/wikipgeul/internal/dump"
)

// sampleLexeme is a small lexeme for tests: L4 (English verb) with a lexeme
// claim, form L4-F1 "windsurfing" with one feature and one claim, and sense
// L4-S1 glossed in en and fr with a P5137 claim.
func sampleLexeme() dump.Entity {
	item := func(id string) *dump.DataValue {
		return &dump.DataValue{Type: "wikibase-entityid", Value: json.RawMessage(`{"entity-type":"item","id":"` + id + `"}`)}
	}
	claim := func(id, p, v string) map[string][]dump.Statement {
		return map[string][]dump.Statement{p: {{ID: id, Rank: "normal", Mainsnak: dump.Snak{SnakType: "value", Property: p, Datatype: "wikibase-item", DataValue: item(v)}}}}
	}
	return dump.Entity{
		ID: "L4", Type: "lexeme", LastRevID: 7, Language: "Q1860", LexicalCategory: "Q24905",
		Lemmas: map[string]dump.Term{"en": {Language: "en", Value: "windsurf"}, "en-gb": {Language: "en-gb", Value: "windsurf"}},
		Claims: claim("L4$A", "P5185", "Q499327"),
		Forms: []dump.Form{{ID: "L4-F1", Representations: map[string]dump.Term{"en": {Language: "en", Value: "windsurfing"}},
			GrammaticalFeatures: []string{"Q10345583"}, Claims: claim("L4-F1$A", "P898", "Q1")}},
		Senses: []dump.Sense{{ID: "L4-S1", Glosses: map[string]dump.Term{"en": {Language: "en", Value: "to sail on a sailboard"}, "fr": {Language: "fr", Value: "faire de la planche à voile"}},
			Claims: claim("L4-S1$A", "P5137", "Q217636")}},
	}
}
