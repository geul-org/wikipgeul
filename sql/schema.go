// Package sql gives the wikidata schema (wikidata.sql) to other modules, such
// as wikigeul's tests, without copying it.
package sql

import _ "embed"

// Schema is wikidata.sql: sections 1–3 create the schema wikidata and its
// tables; section 4 (cleanup, keys and indexes after a load) is commented out.
//
//go:embed wikidata.sql
var Schema string

// RepairReferenceSnak is repair_reference_snak.sql: the work tables of
// `wikipg repair-reference-snak` (hashes done, checkpoints), created in the
// search_path schema.
//
//go:embed repair_reference_snak.sql
var RepairReferenceSnak string

// ReferenceSnakOrder is reference_snak_order.sql: the table of
// `wikipg build-reference-order` (each reference occurrence whose snak order
// differs from reference_snak, as a permutation of its ordinals) and its
// checkpoints, created in the search_path schema.
//
//go:embed reference_snak_order.sql
var ReferenceSnakOrder string

// MigrateCheckpointKey is migrate_checkpoint_key.sql: moves the checkpoint
// tables of the two commands above, made before 2026-10-29 with the key run
// alone, to the key (dump_id, run) in the search_path schema (tables missing
// or already moved are left alone).
//
//go:embed migrate_checkpoint_key.sql
var MigrateCheckpointKey string
