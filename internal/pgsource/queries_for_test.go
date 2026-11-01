//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what TestQueriesFor — every query reads the given schema and dump and orders its rows (dump order, byte order of language and site); no language filter; with reference_snak_order the references go in each occurrence's order
package pgsource

import (
	"strings"
	"testing"
)

func TestQueriesFor(t *testing.T) {
	q := queriesFor(Store{Schema: "wd", Dump: 3})
	for _, s := range []string{q.Exists, q.Labels, q.Descriptions, q.Aliases, q.Claims, q.Qualifiers, q.References, q.Sitelinks} {
		if !strings.Contains(s, "wd.") || !strings.Contains(s, "dump_id = 3") {
			t.Fatalf("schema or dump missing: %s", s)
		}
	}
	if !strings.Contains(q.Aliases, "ordinal") || !strings.Contains(q.Claims, "ORDER BY statement_id") || !strings.Contains(q.Qualifiers, "ORDER BY statement_id, ordinal") || !strings.Contains(q.References, "ORDER BY r.statement_id, r.ordinal, s.ordinal") || !strings.Contains(q.References, "s.hash = r.hash") ||
		!strings.Contains(q.Sitelinks, `ORDER BY site COLLATE "C"`) || !strings.Contains(q.Labels, `ORDER BY language COLLATE "C"`) || strings.Contains(q.Labels+q.Descriptions+q.Aliases, "$2") {
		t.Fatalf("order: %+v", q)
	}
	if strings.Contains(q.References, "reference_snak_order") {
		t.Fatalf("no RefOrder: %s", q.References)
	}
	o := queriesFor(Store{Schema: "wd", Dump: 3, RefOrder: true}).References
	if !strings.Contains(o, "LEFT JOIN LATERAL (SELECT perm FROM wd.reference_snak_order WHERE statement_id = r.statement_id AND ordinal = r.ordinal OFFSET 0) o ON true WHERE") ||
		!strings.HasSuffix(o, "ORDER BY r.statement_id, r.ordinal, coalesce(array_position(o.perm, s.ordinal) - 1, s.ordinal)") {
		t.Fatalf("RefOrder: %s", o)
	}
}
