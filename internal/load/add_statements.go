//ff:func feature=load type=model control=iteration dimension=2
//ff:what Rows.addStatements — the statements of an entity from key number start, properties in sorted order, each with its key (statementID); the next number
package load

import (
	"slices"

	"github.com/geul-org/wikipgeul/internal/dump"
)

// addStatements adds the statements claims of id (an entity, or a lexeme's
// form or sense): properties in sorted order (so keys are the same on every
// run), each property's statements in array order with that position as the
// ordinal. Keys continue from start, the i-th statement of the dump line
// getting statementID(dump, line, i), so a lexeme and its forms and senses
// share one key sequence. It returns the next number and false when the key
// range ran out, which is recorded as one load_error.
func (r *Rows) addStatements(c Config, line int64, id string, claims map[string][]dump.Statement, start int, refs *RefSet) (int, bool) {
	props := make([]string, 0, len(claims))
	for p := range claims {
		props = append(props, p)
	}
	slices.Sort(props)
	i := start
	for _, p := range props {
		for ord, s := range claims[p] {
			sid, err := statementID(c.DumpID, line, i)
			if err != nil {
				r.Error(c.DumpID, line, id, "statement", err.Error(), nil)
				return i, false
			}
			r.addStatement(c, line, id, sid, ord, s, refs)
			i++
		}
	}
	return i, true
}
