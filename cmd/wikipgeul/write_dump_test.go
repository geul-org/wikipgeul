//ff:func feature=cli type=command control=iteration dimension=1
//ff:what writeDump — a small dump file (JSON array, one entity per line) named like a Wikidata dump, with its .sha256 sidecar
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeDump writes, in a fresh directory, wikidata-20261019-all.json holding
// "[", n entity lines Q1…Qn (each with an English label, a claim with a
// qualifier and a reference) and "]", plus its .sha256 sidecar, and returns
// the file path. Entity bad (1-based, 0 for none) is written as broken JSON.
// Entity 2's label and claim value hold the escape \u0000, which text and
// jsonb cannot hold.
func writeDump(t *testing.T, n, bad int) string {
	t.Helper()
	var b strings.Builder
	b.WriteString("[\n")
	for i := 1; i <= n; i++ {
		line := fmt.Sprintf(`{"type":"item","id":"Q%d","lastrevid":%d,"modified":"2026-10-18T00:00:00Z","labels":{"en":{"language":"en","value":"item %d"}},`+
			`"claims":{"P31":[{"mainsnak":{"snaktype":"value","property":"P31","datavalue":{"value":{"entity-type":"item","id":"Q5"},"type":"wikibase-entityid"},"datatype":"wikibase-item"},`+
			`"id":"Q%d$a","rank":"normal","qualifiers":{"P580":[{"snaktype":"novalue","property":"P580","hash":"q"}]},"qualifiers-order":["P580"],`+
			`"references":[{"hash":"r1","snaks":{"P143":[{"snaktype":"value","property":"P143","datavalue":{"value":{"entity-type":"item","id":"Q328"},"type":"wikibase-entityid"},"datatype":"wikibase-item"}]},"snaks-order":["P143"]}]}]}}`, i, 100+i, i, i)
		if i == 2 {
			line = strings.Replace(line, `"item 2"`, `"item\u00002"`, 1)
			line = strings.Replace(line, `"datavalue":{"value":{"entity-type":"item","id":"Q5"},"type":"wikibase-entityid"},"datatype":"wikibase-item"`, `"datavalue":{"value":"a\u0000b","type":"string"},"datatype":"string"`, 1)
		}
		if i == bad {
			line = `{"type":"item","id":`
		}
		sep := ","
		if i == n {
			sep = ""
		}
		b.WriteString(line + sep + "\n")
	}
	b.WriteString("]\n")
	dir := t.TempDir()
	p := filepath.Join(dir, "wikidata-20261019-all.json")
	if err := os.WriteFile(p, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(p+".sha256", []byte("abc123  wikidata-20261019-all.json\n"), 0o644)
	return p
}
