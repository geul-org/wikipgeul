//ff:func feature=convert type=helper control=iteration dimension=1
//ff:what pageKind — send every entity of one kind the dump holds after a number and up to a bound, page by page
package exportgeul

import (
	"database/sql"

	"github.com/geul-org/wikigeul/pkg/pipeline"
	"github.com/geul-org/wikipgeul/internal/pgsource"
)

// pageKind walks the dump's entities of kind in (after, to] in pages of
// 10,000 (pgsource.EntityNumbers) and calls send for each, in order.
func pageKind(db *sql.DB, c config, kind byte, after, to uint64, send func(pipeline.PackJob)) error {
	for {
		nums, err := pgsource.EntityNumbers(db, pgsource.Store{Schema: c.Schema, Dump: c.Dump}, kind, after, to, 10000)
		if err != nil || len(nums) == 0 {
			return err
		}
		for _, n := range nums {
			send(pipeline.PackJob{Kind: kind, Num: n})
		}
		after = nums[len(nums)-1]
	}
}
