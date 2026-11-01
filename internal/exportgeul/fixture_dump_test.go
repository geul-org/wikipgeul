//ff:func feature=convert type=helper control=sequence
//ff:what fixtureDump — test fixture: the dump row of testdata/fixture.sql as provenanceOf reads it
package exportgeul

import "github.com/geul-org/wikipgeul/internal/pgsource"

// fixtureDump is a dump loaded by wikipgeul, for building a provenance.
func fixtureDump() pgsource.Dump {
	return pgsource.Dump{URL: "u", Loader: wikipgeulPath, LoaderCommit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
}
