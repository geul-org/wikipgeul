//ff:func feature=cli type=command control=sequence
//ff:what main — wikipg entry: loads a Wikidata JSON dump into the wikidata schema of geuldev (sql/wikidata.sql), or runs a subcommand
package main

import (
	"fmt"
	"os"
)

// main runs wikipg with the command-line arguments (runCommand) and exits
// with exitCode on error.
func main() {
	if err := runCommand(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "wikipg:", err)
		os.Exit(exitCode(err, os.Args, os.Stderr))
	}
}
