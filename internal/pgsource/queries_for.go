//ff:func feature=convert type=loader control=sequence
//ff:what queriesFor — the SQL for one entity of a loaded dump: existence, texts in every language, claims, qualifiers and references in dump order, sitelinks by site
package pgsource

import "fmt"

// queriesFor reads one entity of dump s.Dump from schema s.Schema (wiki/
// wikidata.sql; the name must pass validSchema): whether it exists with its
// datatype (properties), revision ID and save time (POSIX seconds, 0 when
// none), labels and descriptions in every language by
// language code (byte order), aliases by language then dump order, statements
// by statement_id (wikipgeul numbers them by property ID as text between
// properties and in dump order within a property; the order between
// properties means nothing in Wikidata) and qualifiers by statement then their dump order
// (qualifiers-order), references by statement, reference order and
// snaks-order (statement_reference and reference_snak by their primary keys;
// the hash is compared as stored, never wrapped in a function; each
// occurrence's own snak order from reference_snak_order when s.RefOrder,
// refOrderSQL), sitelinks by site ID in byte order (the dump's key
// order) with their badges. A value is the bare entity ID for an entity
// (value_entity), the unquoted string for a JSON string, else the datavalue
// JSON as text (jsonb's spacing and key order); novalue and somevalue give an
// empty value, told apart by the snak type column (snaktype). A value wikipgeul kept in value_json (\u0000, which text and
// jsonb cannot hold) comes as a further column (exactValue). The statement IDs are passed as an array rather than joined,
// which PostgreSQL would plan as a slow JIT-compiled join.
func queriesFor(s Store) entityQueries {
	value := "coalesce(value_entity, value #>> '{}', '')"
	join, order := refOrderSQL(s)
	return entityQueries{
		Exists: fmt.Sprintf("SELECT count(*), coalesce(max(datatype), ''), coalesce(max(lastrevid), 0), coalesce(floor(extract(epoch FROM max(modified)))::bigint, 0) "+
			"FROM %s.entity WHERE entity_id = $1 AND dump_id = %d", s.Schema, s.Dump),
		Labels:       fmt.Sprintf("SELECT language, value, coalesce(value_json, '') FROM %s.label WHERE entity_id = $1 AND dump_id = %d ORDER BY language COLLATE \"C\"", s.Schema, s.Dump),
		Descriptions: fmt.Sprintf("SELECT language, value, coalesce(value_json, '') FROM %s.description WHERE entity_id = $1 AND dump_id = %d ORDER BY language COLLATE \"C\"", s.Schema, s.Dump),
		Aliases:      fmt.Sprintf("SELECT language, value, coalesce(value_json, '') FROM %s.alias WHERE entity_id = $1 AND dump_id = %d ORDER BY language COLLATE \"C\", ordinal", s.Schema, s.Dump),
		Claims: fmt.Sprintf("SELECT statement_id, property, coalesce(value_type, ''), coalesce(datatype, ''), %s, coalesce(value_json, ''), rank::text, snaktype::text FROM %s.statement "+
			"WHERE entity_id = $1 AND dump_id = %d ORDER BY statement_id", value, s.Schema, s.Dump),
		Qualifiers: fmt.Sprintf("SELECT statement_id, property, coalesce(value_type, ''), coalesce(datatype, ''), %s, coalesce(value_json, ''), snaktype::text FROM %s.qualifier "+
			"WHERE statement_id = ANY(ARRAY(SELECT statement_id FROM %s.statement WHERE entity_id = $1 AND dump_id = %d)) ORDER BY statement_id, ordinal",
			value, s.Schema, s.Schema, s.Dump),
		References: fmt.Sprintf("SELECT r.statement_id, r.ordinal, r.hash, s.property, coalesce(s.value_type, ''), coalesce(s.datatype, ''), %s, coalesce(s.value_json, ''), s.snaktype::text "+
			"FROM %s.statement_reference r JOIN %s.reference_snak s ON s.hash = r.hash AND s.dump_id = %d%s "+
			"WHERE r.statement_id = ANY(ARRAY(SELECT statement_id FROM %s.statement WHERE entity_id = $1 AND dump_id = %d)) ORDER BY r.statement_id, r.ordinal, %s",
			"coalesce(s.value_entity, s.value #>> '{}', '')", s.Schema, s.Schema, s.Dump, join, s.Schema, s.Dump, order),
		Sitelinks: fmt.Sprintf("SELECT site, title, coalesce(title_json, ''), badges FROM %s.sitelink WHERE entity_id = $1 AND dump_id = %d ORDER BY site COLLATE \"C\"", s.Schema, s.Dump),
	}
}
