//ff:type feature=cli type=model
//ff:what repairFlags — the parsed command line of repair-reference-snak: the database DSN and the repair run
package main

import "github.com/geul-org/wikipgeul/internal/repair"

// repairFlags is the command line of wikipg repair-reference-snak: the
// PostgreSQL DSN (empty: PG* variables) and the run (repair.Config).
type repairFlags struct {
	DSN    string
	Config repair.Config
}
