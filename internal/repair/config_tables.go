//ff:func feature=repair type=helper control=selection
//ff:what Config.tables — the work tables of the run's mode: their SQL (sql/*.sql) and the checkpoint table's name
package repair

import schema "github.com/geul-org/wikipgeul/sql"

// tables returns the SQL creating the work tables of c's mode and the name
// of its checkpoint table: reference_snak_order.sql and
// reference_snak_order_checkpoint for an order build, else
// repair_reference_snak.sql and repair_reference_snak_checkpoint.
func (c Config) tables() (string, string) {
	switch {
	case c.Order:
		return schema.ReferenceSnakOrder, "reference_snak_order_checkpoint"
	default:
		return schema.RepairReferenceSnak, "repair_reference_snak_checkpoint"
	}
}
