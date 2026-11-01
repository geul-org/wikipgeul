//ff:func feature=cli type=command control=sequence
//ff:what runOrder — wikipg build-reference-order: records (or with -dry-run counts) each reference occurrence whose snak order differs from reference_snak, from entity_raw; stops cleanly on interrupt
package main

// runOrder parses args (parseRepairFlags: the same range, batch and pace
// flags as repair-reference-snak) and runs an order build (execRepair with
// Config.Order): reference_snak_order gets the permutation of every
// reference occurrence whose snaks are those of reference_snak in another
// order; reference_snak and the other tables are only read.
func runOrder(args []string) error {
	f, err := parseRepairFlags("wikipg build-reference-order", args)
	if err != nil {
		return err
	}
	f.Config.Order = true
	return execRepair(f)
}
