//ff:func feature=repair type=decoder control=iteration dimension=1
//ff:what occurrencesOf — one entity_raw line as the loader reads it: entity ID, every reference occurrence with its snak rows, and its statement_reference rows
package repair

import (
	"errors"
	"fmt"

	"github.com/geul-org/wikipgeul/internal/dump"
	"github.com/geul-org/wikipgeul/internal/load"
)

// occurrencesOf decodes raw (dump line line) and builds its rows as the loader
// does (load.Rows.AddEntity, same statement keys) with no RefSet, so every
// reference occurrence keeps its own snak list. It returns the entity ID, the
// occurrences in dump order (statement key, then position; a reference with
// no snaks has no rows) and the statement_reference rows. A line that does
// not decode, or that the loader records a load_error for, is an error.
func occurrencesOf(dumpID int16, line int64, raw string) (string, []occ, []stmtRef, error) {
	e, err := dump.DecodeEntity([]byte(raw))
	if err != nil {
		return "", nil, nil, err
	}
	r := load.NewRows()
	r.AddEntity(load.Config{DumpID: dumpID}, line, nil, e, nil)
	if r.Errors > 0 {
		return e.ID, nil, nil, errors.New(fmt.Sprint(r.Table("load_error")[0][4]))
	}
	var occs []occ
	var refs []stmtRef
	at := map[load.Occurrence]int{}
	for _, row := range r.Table("statement_reference") {
		ref := stmtRef{Statement: row[0].(int64), Ordinal: row[1].(int), Hash: row[2].(string)}
		refs = append(refs, ref)
		at[load.Occurrence{Statement: ref.Statement, Ordinal: ref.Ordinal}] = len(occs)
		occs = append(occs, occ{Hash: ref.Hash, At: load.Occurrence{Statement: ref.Statement, Ordinal: ref.Ordinal}, Entity: e.ID})
	}
	for _, row := range r.Table("reference_snak") {
		i := at[load.Occurrence{Statement: row[11].(int64), Ordinal: row[12].(int)}]
		occs[i].Rows = append(occs[i].Rows, row[:11])
	}
	return e.ID, occs, refs, nil
}
