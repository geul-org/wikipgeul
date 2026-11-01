//ff:func feature=convert type=helper control=sequence
//ff:what provenanceOf — the provenance every document of this run records: the loaded dump's URL, SHA-256 and time, and the converter and library commits
package exportgeul

import (
	"fmt"

	"github.com/geul-org/wikigeul/pkg/pipeline"
	"github.com/geul-org/wikigeul/pkg/port"
	"github.com/geul-org/wikipgeul/internal/pgsource"
)

// wikipgeulPath is the loader that wikipgeul (github.com/geul-org/wikipgeul) records in
// the dump row. The legacy paths are what it recorded under its earlier names
// (repository wikiport, command wikiload; the module was .../wikiport/go until
// 2026-10-27; repository and command wikipg until 2026-11-01); dumps loaded
// then keep them.
const (
	wikipgeulPath      = "github.com/geul-org/wikipgeul/cmd/wikipgeul"
	legacyWikiloadPath = "github.com/geul-org/wikiport/cmd/wikiload"
	legacyWikiportPath = "github.com/geul-org/wikiport/go/cmd/wikiload"
	legacyWikipgPath   = "github.com/geul-org/wikipg/cmd/wikipg"
)

// provenanceOf is the provenance of dump d (README 출처 기록): its URL,
// file SHA-256 and download time, then the tools as geul.org term definition
// documents (GEUL-ORG) — wikigeul (document 1, this repository) at commit and the
// GEUL library (document 2, gogeul) at libCommit. The loader is not a tool of
// the document: wikipgeul stores the dump losslessly (every line in
// entity_raw, values text or jsonb cannot hold in *_json), so the document
// follows from the dump and these two; its commit stays on the dump row. A
// loader other than wikipgeul is an error: its store is not known to be
// lossless.
func provenanceOf(d pgsource.Dump, commit, libCommit string) (port.Provenance, error) {
	if d.Loader != wikipgeulPath && d.Loader != legacyWikipgPath && d.Loader != legacyWikiloadPath && d.Loader != legacyWikiportPath {
		return port.Provenance{}, fmt.Errorf("wikigeul: dump loaded by %q, not %s", d.Loader, wikipgeulPath)
	}
	return port.Provenance{
		DumpURI:    d.URL,
		DumpSHA256: d.SHA256,
		AcquiredAt: d.DownloadedAt.Unix(),
		Tools:      pipeline.Tools(commit, libCommit),
	}, nil
}
