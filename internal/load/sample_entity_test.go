//ff:func feature=load type=model control=sequence
//ff:what sampleEntity — a decoded item with terms in two languages, a sitelink and claims with qualifiers, a reference, a quantity and novalue
package load

import (
	"encoding/json"

	"github.com/geul-org/wikipgeul/internal/dump"
)

// sampleEntity returns item Q42 as DecodeEntity gives it: labels and
// descriptions in en and ko, two English aliases, an enwiki sitelink, a P69
// claim with qualifiers P580 and P582 (qualifiers-order P580, P582) and the
// reference r1, a quantity P2048 with unit and bounds, and a novalue P40.
func sampleEntity() dump.Entity {
	dv := func(typ, v string) *dump.DataValue { return &dump.DataValue{Type: typ, Value: json.RawMessage(v)} }
	time := func(p, y string) dump.Snak {
		return dump.Snak{SnakType: "value", Property: p, Hash: "h" + p, Datatype: "time", DataValue: dv("time", `{"time":"+`+y+`-00-00T00:00:00Z","precision":9}`)}
	}
	return dump.Entity{
		ID: "Q42", Type: "item", LastRevID: 2401234567, Modified: "2026-10-18T09:00:00Z",
		Labels:       map[string]dump.Term{"en": {Language: "en", Value: "Douglas Adams"}, "ko": {Language: "ko", Value: "더글러스 애덤스"}},
		Descriptions: map[string]dump.Term{"en": {Language: "en", Value: "English writer"}, "ko": {Language: "ko", Value: "영국 작가"}},
		Aliases:      map[string][]dump.Term{"en": {{Language: "en", Value: "DNA"}, {Language: "en", Value: "Douglas Noël Adams"}}},
		Sitelinks:    map[string]dump.Sitelink{"enwiki": {Site: "enwiki", Title: "Douglas Adams", Badges: []string{"Q17437796"}}},
		Claims: map[string][]dump.Statement{
			"P69": {{ID: "Q42$A", Rank: "normal",
				Mainsnak:        dump.Snak{SnakType: "value", Property: "P69", Datatype: "wikibase-item", DataValue: dv("wikibase-entityid", `{"entity-type":"item","id":"Q691283"}`)},
				Qualifiers:      map[string][]dump.Snak{"P582": {time("P582", "1974")}, "P580": {time("P580", "1971")}},
				QualifiersOrder: []string{"P580", "P582"},
				References: []dump.Reference{{Hash: "r1", SnaksOrder: []string{"P248"},
					Snaks: map[string][]dump.Snak{"P248": {{SnakType: "value", Property: "P248", Datatype: "wikibase-item", DataValue: dv("wikibase-entityid", `{"entity-type":"item","id":"Q5375741"}`)}}}}}}},
			"P2048": {{ID: "Q42$B", Rank: "preferred", Mainsnak: dump.Snak{SnakType: "value", Property: "P2048", Datatype: "quantity",
				DataValue: dv("quantity", `{"amount":"+1.96","unit":"http://www.wikidata.org/entity/Q11573","upperBound":"+1.97","lowerBound":"+1.95"}`)}}},
			"P40": {{ID: "Q42$C", Rank: "normal", Mainsnak: dump.Snak{SnakType: "novalue", Property: "P40", Datatype: "wikibase-item"}}},
		},
	}
}
