//ff:func feature=convert type=loader control=sequence
//ff:what TestRefOrderSQL — with reference_snak_order a LATERAL primary-key lookup and the place in perm (ordinal when none), without it nothing joined and the ordinal
package pgsource

import (
	"strings"
	"testing"
)

func TestRefOrderSQL(t *testing.T) {
	join, order := refOrderSQL(Store{Schema: "wd", Dump: 1, RefOrder: true})
	if !strings.HasPrefix(join, " LEFT JOIN LATERAL (SELECT perm FROM wd.reference_snak_order WHERE statement_id = r.statement_id AND ordinal = r.ordinal OFFSET 0) o") ||
		order != "coalesce(array_position(o.perm, s.ordinal) - 1, s.ordinal)" {
		t.Fatalf("with: %q %q", join, order)
	}
	if join, order := refOrderSQL(Store{Schema: "wd", Dump: 1}); join != "" || order != "s.ordinal" {
		t.Fatalf("without: %q %q", join, order)
	}
}
