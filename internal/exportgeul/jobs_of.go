//ff:func feature=convert type=helper control=iteration dimension=1
//ff:what jobsOf — jobs for entity numbers of one kind
package exportgeul

import "github.com/geul-org/wikigeul/pkg/pipeline"

// jobsOf makes a job per number of kind.
func jobsOf(kind byte, nums []uint64) []pipeline.PackJob {
	out := make([]pipeline.PackJob, len(nums))
	for i, n := range nums {
		out[i] = pipeline.PackJob{Kind: kind, Num: n}
	}
	return out
}
