-- fixture.sql — a few entities in the wikidata schema (wikipgeul sql/wikidata.sql) for export-geul's tests.
-- Run after sections 1–3 with search_path set to the test schema. Q6 is left out (never created).
-- Q42 has labels in three languages (mul included) and sitelinks of several sites, one unknown.
-- Q7 holds \u0000 values that wikipgeul keeps in the *_json columns (value NULL or cleaned).
-- Q30 holds special snaks: a border whose end (P582) is somevalue, a novalue central bank (P1304) with time qualifiers, a somevalue main snak.
INSERT INTO dump (dump_id, dump_date, file_name, url, size_bytes, sha256, downloaded_at, loader, loader_commit) VALUES
  (1, '2026-10-19', 'wikidata-20261019-all.json.bz2', 'https://dumps.wikimedia.org/wikidatawiki/entities/20261019/wikidata-20261019-all.json.bz2',
   103488697311, 'a249a19d12b1b0e8443248da3583090682be3f20c6bdd6ed1f7e1ae05b6b0774', '2026-10-26 06:34:43+00',
   'github.com/geul-org/wikipgeul/cmd/wikipgeul', 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb'),
  (2, '2026-10-19', 'other.json.bz2', 'https://example.org/other', 1, '00', '2026-10-26 06:34:43+00', 'elsewhere/loader', 'bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb');

INSERT INTO entity (entity_id, line_no, kind, num, type, datatype, lastrevid, modified, dump_id) VALUES
  ('Q42', 5, 'Q', 42, 'item', NULL, 2557857485, '2026-10-17 04:21:25+00', 1);
INSERT INTO entity (entity_id, line_no, kind, num, type, datatype, dump_id) VALUES
  ('Q1', 2, 'Q', 1, 'item', NULL, 1), ('Q2', 3, 'Q', 2, 'item', NULL, 1), ('Q3', 4, 'Q', 3, 'item', NULL, 1),
  ('P31', 6, 'P', 31, 'property', 'wikibase-item', 1),
  ('P580', 7, 'P', 580, 'property', 'time', 1), ('Q99', 8, 'Q', 99, 'item', NULL, 2), ('Q7', 9, 'Q', 7, 'item', NULL, 1),
  ('Q30', 10, 'Q', 30, 'item', NULL, 1);

INSERT INTO label (entity_id, dump_id, language, value) VALUES
  ('Q1', 1, 'en', 'Universe'), ('Q2', 1, 'en', 'Earth'), ('Q42', 1, 'en', 'Douglas Adams'), ('Q42', 1, 'ko', '더글러스 애덤스'), ('Q42', 1, 'mul', 'Douglas Adams'),
  ('P31', 1, 'en', 'instance of'), ('Q42', 2, 'en', 'other dump');
INSERT INTO label (entity_id, dump_id, language, value, value_json) VALUES
  ('Q7', 1, 'en', 'ab', '"a\u0000b"');
INSERT INTO description (entity_id, dump_id, language, value) VALUES
  ('Q42', 1, 'en', 'English writer and humorist'), ('Q42', 1, 'ko', '영국의 작가');
INSERT INTO alias (entity_id, dump_id, language, ordinal, value) VALUES
  ('Q42', 1, 'en', 1, 'Douglas Noël Adams'), ('Q42', 1, 'en', 0, 'DNA');

INSERT INTO statement (statement_id, guid, entity_id, property, ordinal, rank, snaktype, datatype, value_type, value, value_entity, dump_id) VALUES
  (101, 'Q1$a', 'Q1', 'P31', 0, 'normal', 'value', 'wikibase-item', 'wikibase-entityid', '{"id": "Q36906466", "numeric-id": 36906466, "entity-type": "item"}', 'Q36906466', 1),
  (201, 'Q2$a', 'Q2', 'P31', 0, 'normal', 'value', 'wikibase-item', 'wikibase-entityid', '{"id": "Q3504248", "numeric-id": 3504248, "entity-type": "item"}', 'Q3504248', 1),
  (4201, 'Q42$a', 'Q42', 'P31', 0, 'preferred', 'value', 'wikibase-item', 'wikibase-entityid', '{"id": "Q5", "numeric-id": 5, "entity-type": "item"}', 'Q5', 1),
  (4202, 'Q42$b', 'Q42', 'P569', 0, 'normal', 'value', 'time', 'time', '{"time": "+1952-03-11T00:00:00Z", "after": 0, "before": 0, "timezone": 0, "precision": 11, "calendarmodel": "http://www.wikidata.org/entity/Q1985727"}', NULL, 1),
  (4203, 'Q42$c', 'Q42', 'P1559', 0, 'normal', 'value', 'monolingualtext', 'monolingualtext', '{"text": "Douglas Adams", "language": "en"}', NULL, 1),
  (4204, 'Q42$d', 'Q42', 'P214', 0, 'normal', 'value', 'external-id', 'string', '"113230702"', NULL, 1),
  (4205, 'Q42$e', 'Q42', 'P40', 0, 'normal', 'novalue', 'wikibase-item', NULL, NULL, NULL, 1),
  (4206, 'Q42$f', 'Q42', 'P2067', 0, 'normal', 'value', 'quantity', 'quantity', '{"amount": "+73", "unit": "http://www.wikidata.org/entity/Q11570"}', NULL, 1),
  (4207, 'Q42$g', 'Q42', 'P856', 0, 'normal', 'value', 'url', 'string', '"https://douglasadams.com/더글러스"', NULL, 1),
  (4208, 'Q42$h', 'Q42', 'P18', 0, 'normal', 'value', 'commonsMedia', 'string', '"Douglas adams portrait.jpg"', NULL, 1),
  (9901, 'Q99$a', 'Q99', 'P31', 0, 'normal', 'value', 'wikibase-item', 'wikibase-entityid', '{"id": "Q5", "numeric-id": 5, "entity-type": "item"}', 'Q5', 2);

INSERT INTO qualifier (statement_id, ordinal, property, value_ordinal, snaktype, datatype, value_type, value, value_entity) VALUES
  (4201, 1, 'P642', 0, 'value', 'wikibase-item', 'wikibase-entityid', '{"id": "Q1", "numeric-id": 1, "entity-type": "item"}', 'Q1'),
  (4201, 0, 'P580', 0, 'value', 'time', 'time', '{"time": "+1979-00-00T00:00:00Z", "after": 0, "before": 0, "timezone": 0, "precision": 9, "calendarmodel": "http://www.wikidata.org/entity/Q1985727"}', NULL),
  (4201, 2, 'P582', 0, 'novalue', 'time', NULL, NULL, NULL);
INSERT INTO statement (statement_id, guid, entity_id, property, ordinal, rank, snaktype, datatype, value_type, value, value_entity, value_json, dump_id) VALUES
  (701, 'Q7$a', 'Q7', 'P14457', 0, 'normal', 'value', 'string', 'string', NULL, NULL, '"x\u0000"', 1);
INSERT INTO qualifier (statement_id, ordinal, property, value_ordinal, snaktype, datatype, value_type, value, value_entity, value_json) VALUES
  (701, 0, 'P1810', 0, 'value', 'string', 'string', NULL, NULL, '"q\u0000"');

INSERT INTO statement (statement_id, guid, entity_id, property, ordinal, rank, snaktype, datatype, value_type, value, value_entity, dump_id) VALUES
  (3001, 'Q30$a', 'Q30', 'P1304', 0, 'normal', 'novalue', 'wikibase-item', NULL, NULL, NULL, 1),
  (3002, 'Q30$b', 'Q30', 'P35', 0, 'deprecated', 'somevalue', 'wikibase-item', NULL, NULL, NULL, 1),
  (3003, 'Q30$c', 'Q30', 'P47', 0, 'normal', 'value', 'wikibase-item', 'wikibase-entityid', '{"id": "Q34266", "numeric-id": 34266, "entity-type": "item"}', 'Q34266', 1);
INSERT INTO qualifier (statement_id, ordinal, property, value_ordinal, snaktype, datatype, value_type, value, value_entity) VALUES
  (3001, 0, 'P580', 0, 'value', 'time', 'time', '{"time": "+1811-00-00T00:00:00Z", "after": 0, "before": 0, "timezone": 0, "precision": 9, "calendarmodel": "http://www.wikidata.org/entity/Q1985727"}', NULL),
  (3001, 1, 'P582', 0, 'value', 'time', 'time', '{"time": "+1816-00-00T00:00:00Z", "after": 0, "before": 0, "timezone": 0, "precision": 9, "calendarmodel": "http://www.wikidata.org/entity/Q1985727"}', NULL),
  (3003, 0, 'P582', 0, 'somevalue', 'time', NULL, NULL, NULL);

INSERT INTO sitelink (entity_id, dump_id, site, title, badges) VALUES
  ('Q42', 1, 'enwiki', 'Douglas Adams', '{Q17437796}'), ('Q42', 1, 'kowiki', '더글러스 애덤스', '{}'),
  ('Q42', 1, 'commonswiki', 'Category:Douglas Adams', '{}'), ('Q42', 1, 'be_x_oldwiki', 'Дуглас Адамс', '{}'),
  ('Q42', 1, 'enwikiquote', 'What?/A & B''s 100%', '{Q17437796,P1}'), ('Q42', 1, 'xxsite', 'x', '{}'),
  ('Q42', 1, 'zh_min_nanwiki', 'Douglas Adams', '{}'), ('Q42', 2, 'enwiki', 'other dump', '{}');
INSERT INTO sitelink (entity_id, dump_id, site, title, title_json) VALUES
  ('Q7', 1, 'enwiki', 'ab', '"a\u0000b"');

-- Q42 P31 (4201) has two references: Britannica (P248 + its ID P1417) and a URL with its retrieval date.
-- Hash h3 exists only in dump 2 (must not be read with dump 1).
INSERT INTO reference (hash, dump_id) VALUES ('h1', 1), ('h2', 1), ('h3', 2);
INSERT INTO statement_reference (statement_id, ordinal, hash) VALUES (4201, 1, 'h2'), (4201, 0, 'h1'), (4202, 0, 'h3');
INSERT INTO reference_snak (hash, dump_id, ordinal, property, value_ordinal, snaktype, datatype, value_type, value, value_entity) VALUES
  ('h1', 1, 1, 'P1417', 0, 'value', 'external-id', 'string', '"biography/Douglas-Adams"', NULL),
  ('h1', 1, 0, 'P248', 0, 'value', 'wikibase-item', 'wikibase-entityid', '{"id": "Q5375741", "numeric-id": 5375741, "entity-type": "item"}', 'Q5375741'),
  ('h2', 1, 0, 'P854', 0, 'value', 'url', 'string', '"http://www.eskimo.com/~rkj/weekly/aa051701a.htm"', NULL),
  ('h2', 1, 1, 'P813', 0, 'somevalue', 'time', NULL, NULL, NULL),
  ('h3', 2, 0, 'P143', 0, 'value', 'wikibase-item', 'wikibase-entityid', '{"id": "Q328", "numeric-id": 328, "entity-type": "item"}', 'Q328');
