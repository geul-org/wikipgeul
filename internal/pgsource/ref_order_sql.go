//ff:func feature=convert type=loader control=selection
//ff:what refOrderSQL — the references query's join and snak ordering: each occurrence's own order from reference_snak_order when the schema has it, else reference_snak's
package pgsource

import "fmt"

// refOrderSQL returns what the references query adds after its join of
// reference_snak (alias s) to statement_reference (alias r) and the
// expression ordering the snaks of one occurrence. With s.RefOrder the
// occurrence's row of reference_snak_order (perm: reference_snak ordinals in
// the occurrence's order) is looked up by primary key (LATERAL with OFFSET 0,
// one lookup per occurrence) and a snak goes at its place in perm, at its
// ordinal when the occurrence has no row (its order is reference_snak's).
// Without it the snaks go by ordinal.
func refOrderSQL(s Store) (string, string) {
	switch {
	case s.RefOrder:
		return fmt.Sprintf(" LEFT JOIN LATERAL (SELECT perm FROM %s.reference_snak_order WHERE statement_id = r.statement_id AND ordinal = r.ordinal OFFSET 0) o ON true", s.Schema),
			"coalesce(array_position(o.perm, s.ordinal) - 1, s.ordinal)"
	default:
		return "", "s.ordinal"
	}
}
