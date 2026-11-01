//ff:func feature=convert type=helper control=iteration dimension=1
//ff:what TestProvenanceOf — the dump's URL and download time, the converter and the library at their commits (not the loader); wikipgeul under its old names passes; another loader fails
package exportgeul

import (
	"testing"
	"time"

	"github.com/geul-org/wikigeul/pkg/port"
	"github.com/geul-org/wikipgeul/internal/pgsource"
)

func TestProvenanceOf(t *testing.T) {
	d := pgsource.Dump{URL: "u", SHA256: "s", DownloadedAt: time.Unix(5, 0), Loader: wikipgeulPath, LoaderCommit: "l"}
	p, err := provenanceOf(d, "abc", "def")
	if err != nil || p.DumpURI != "u" || p.DumpSHA256 != "s" || p.AcquiredAt != 5 || len(p.Tools) != 2 ||
		p.Tools[0] != (port.Tool{Doc: 1, Commit: "abc"}) || p.Tools[1] != (port.Tool{Doc: 2, Commit: "def"}) {
		t.Fatalf("got %+v %v", p, err)
	}
	for _, old := range []string{legacyWikipgPath, legacyWikiloadPath, legacyWikiportPath} {
		d.Loader = old
		if _, err := provenanceOf(d, "abc", "def"); err != nil {
			t.Fatalf("wikipgeul under its old name %s: %v", old, err)
		}
	}
	for _, other := range []string{"github.com/geul-org/geul-entity/artifacts/wikidata", "github.com/geul-org/wikipgeul/cmd/wikipgeul"} {
		d.Loader = other
		if _, err := provenanceOf(d, "abc", "def"); err == nil {
			t.Fatalf("loader %s: want error", other)
		}
	}
}
