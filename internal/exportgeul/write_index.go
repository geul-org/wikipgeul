//ff:func feature=convert type=helper control=sequence
//ff:what writeIndex — the Issuer's documents.geul for the documents written
package exportgeul

import (
	"os"
	"path/filepath"

	"github.com/geul-org/gogeul/pkg/core"

	"github.com/geul-org/wikipgeul/internal/publish"
)

// writeIndex writes `<root>/.well-known/geul/<issuer>/documents.geul` naming
// each written document by its kind letter and number (q42, p31) and every document already
// published under the Issuer (publishedNumbers), so a run over another range
// does not drop the names of earlier runs (publish.NameIndex).
func writeIndex(root string, issuer []core.Word, nums []uint64) error {
	dir := filepath.Join(root, ".well-known", "geul", core.FormatID(issuer))
	published, err := publishedNumbers(dir)
	if err != nil {
		return err
	}
	b, err := publish.NameIndex(issuer, append(published, nums...))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "documents.geul"), b, 0o644)
}
