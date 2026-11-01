//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what attachReferences — puts each reference on the claim with its statement ID, keeping their order
package pgsource

import "github.com/geul-org/wikigeul/pkg/source"

// attachReferences appends each reference in rs to the claim in claims whose
// ID is its statement ID, in rs order. A reference of no listed claim is dropped.
func attachReferences(claims []source.Claim, rs []source.Reference) {
	at := map[int64]int{}
	for i, c := range claims {
		at[c.ID] = i
	}
	for _, r := range rs {
		if i, ok := at[r.StatementID]; ok {
			claims[i].References = append(claims[i].References, r)
		}
	}
}
