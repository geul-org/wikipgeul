//ff:type feature=convert type=model
//ff:what Store — where a loaded dump lives: the schema of sql/wikidata.sql, the dump_id within it and whether reference_snak_order is there to read
package pgsource

// Store names one loaded dump (sql/wikidata.sql): the schema that holds the
// tables (wikipgeul -schema) and the dump_id of the dump within it (wikipgeul
// -dump-id). Every query reads that dump only. RefOrder says the schema has
// reference_snak_order (Dump.RefOrder): the snaks of a reference occurrence
// are then read in the order it gives for that occurrence, and in
// reference_snak's order (the hash's first occurrence in the dump) when it
// has no row; without it always in reference_snak's order.
type Store struct {
	Schema   string
	Dump     int16
	RefOrder bool
}
