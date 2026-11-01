//ff:func feature=cli type=command control=sequence
//ff:what TestRunLoad — a small dump loads end to end (plain and .bz2), \u0000 values kept in *_json, a limit stops early, a rerun resumes after the checkpoint without duplicates; bad flags, an unparsable DSN, an unreachable database, a broken schema on resume and a missing decompressor fail
package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestRunLoad(t *testing.T) {
	db, schema := testSchema(t)
	p := writeDump(t, 7, 3)
	base := []string{"-schema", schema, "-commit", "test", "-writers", "2", "-parsers", "2", "-decoders", "1", "-batch", "2", "-langs", "all"}
	count := func(table string) (n int) {
		db.QueryRow("SELECT count(*) FROM " + schema + "." + table).Scan(&n)
		return n
	}
	if err := runLoad(append(base, "-in", p, "-limit", "3")); err != nil {
		t.Fatal(err)
	}
	if count("entity") != 2 || count("load_error") != 1 {
		t.Fatalf("limit: entities %d errors %d", count("entity"), count("load_error"))
	}
	if err := runLoad(append(base, "-in", p)); err != nil {
		t.Fatal(err)
	}
	var line, ents, errs int64
	var done bool
	db.QueryRow("SELECT checkpoint_line, entity_count, error_count, load_finished IS NOT NULL FROM "+schema+".dump").Scan(&line, &ents, &errs, &done)
	// The resumed run starts with an empty RefSet, and parallel parsers may
	// meet an earlier occurrence after a later one, so the shared source r1 is
	// written at least once per run (wikidata.sql 4.2 removes the duplicates),
	// and always at its first occurrence (Q1 on line 2, statement 0, reference 0).
	var firstOcc int64
	db.QueryRow("SELECT min(occ_statement_id) FROM " + schema + ".reference_snak WHERE hash = 'r1' AND occ_ordinal = 0").Scan(&firstOcc)
	if count("entity") != 6 || count("statement") != 6 || count("qualifier") != 6 || count("reference") < 2 || firstOcc != 1<<55|2<<20 || line != 8 || ents != 6 || errs != 1 || !done {
		t.Fatalf("resumed: entity %d statement %d qualifier %d reference %d line %d ents %d errs %d done %v",
			count("entity"), count("statement"), count("qualifier"), count("reference"), line, ents, errs, done)
	}
	var labelJSON, valueJSON string
	db.QueryRow("SELECT value_json FROM " + schema + ".label WHERE entity_id = 'Q2'").Scan(&labelJSON)
	db.QueryRow("SELECT value_json FROM " + schema + ".statement WHERE entity_id = 'Q2'").Scan(&valueJSON)
	if labelJSON != `"item\u00002"` || valueJSON != `"a\u0000b"` {
		t.Fatalf("\\u0000 kept: label %q value %q", labelJSON, valueJSON)
	}
	if err := exec.Command("bzip2", "-k", p).Run(); err == nil {
		os.WriteFile(p+".bz2.sha256", []byte("def456  wikidata-20261019-all.json.bz2\n"), 0o644)
		db2, schema2 := testSchema(t)
		if err := runLoad([]string{"-schema", schema2, "-commit", "test", "-in", p + ".bz2", "-decompressor", "bzip2"}); err != nil {
			t.Fatal(err)
		}
		var n int
		db2.QueryRow("SELECT count(*) FROM " + schema2 + ".entity").Scan(&n)
		if n != 6 {
			t.Fatalf("bz2: %d", n)
		}
		if err := runLoad([]string{"-schema", schema2, "-commit", "test", "-dump-id", "2", "-in", p + ".bz2", "-decompressor", "no-such-program"}); err == nil {
			t.Fatal("missing decompressor: want error")
		}
	}
	db.Exec("DROP TABLE " + schema + ".label CASCADE")
	if err := runLoad(append(base, "-in", p)); err == nil {
		t.Fatal("resume with a broken schema: want error")
	}
	if err := runLoad(nil); err == nil {
		t.Fatal("no flags: want error")
	}
	if err := runLoad([]string{"-dsn", "host='x", "-commit", "test", "-in", p}); err == nil {
		t.Fatal("unparsable DSN: want error")
	}
	if err := runLoad([]string{"-dsn", "host=/nonexistent", "-commit", "test", "-in", p}); err == nil {
		t.Fatal("unreachable database: want error")
	}
	if err := runLoad(append(base, "-in", p, "-dump-id", "9", "-schema", schema+"_none")); err == nil {
		t.Fatal("missing schema: want error")
	}
}
