//ff:func feature=cli type=command control=selection
//ff:what runCommand — dispatches wikipgeul's arguments: the repair-reference-snak, build-reference-order or export-geul subcommand, or a load (no subcommand)
package main

import "github.com/geul-org/wikipgeul/internal/exportgeul"

// runCommand runs the subcommand args[0] names with the rest of args
// (repair-reference-snak: runRepair, build-reference-order: runOrder,
// export-geul: exportgeul.Run); without a subcommand args are the load
// flags (runLoad).
func runCommand(args []string) error {
	switch {
	case len(args) > 0 && args[0] == "repair-reference-snak":
		return runRepair(args[1:])
	case len(args) > 0 && args[0] == "build-reference-order":
		return runOrder(args[1:])
	case len(args) > 0 && args[0] == "export-geul":
		return exportgeul.Run(args[1:])
	default:
		return runLoad(args)
	}
}
