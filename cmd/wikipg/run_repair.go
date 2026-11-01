//ff:func feature=cli type=command control=sequence
//ff:what runRepair — wikipg repair-reference-snak: repairs (or with -dry-run compares) reference_snak against entity_raw and reports; stops cleanly on interrupt
package main

// runRepair parses args (parseRepairFlags) and runs the repair (execRepair).
func runRepair(args []string) error {
	f, err := parseRepairFlags("wikipg repair-reference-snak", args)
	if err != nil {
		return err
	}
	return execRepair(f)
}
