-- reference_snak_order.sql — 출현별 출처 snak 순서 (`wikipg build-reference-order`가 만든다)
--
-- 2026-10-29. reference_snak은 hash마다 덤프상 첫 출현의 snak 목록 하나만 담는다(wikidata.sql
-- reference_snak 주석의 출현 규칙). 위키베이스 출처 hash는 snaks-order와 같은 속성 안 값 순서에
-- 무관해서, 같은 hash의 다른 출현은 순서가 다를 수 있다(2026-10-29 재변환 검증 S10: Q1–Q10000 출현
-- 435,343개 중 900개). 이 표는 그 출현들의 순서를 남긴다. 나머지 출현(대부분)은 행이 없고, 순서는
-- reference_snak.ordinal 그대로다. 원천은 entity_raw다(덤프 파일을 다시 읽지 않는다).
--
-- 도구가 search_path를 대상 스키마로 두고 실행한다(CREATE TABLE IF NOT EXISTS). 4.3을 마친 뒤
-- (entity_raw_dump_line, statement_reference·reference_snak PK가 있어야 한다) 4.4로 돌린다.

-- 출현 하나 = statement_reference의 (statement_id, ordinal). 키를 statement_reference와 같게 둔다.
--   - dump_id 열이 없다: statement_reference처럼 statement_id(dump << 55 | 줄 << 20 | 순번)에 판이
--     들어 있다. 판별 조회는 statement_id 범위다(reference_snak처럼 hash로 찾는 표가 아니다).
--   - perm[i] = 이 출현의 i번째(0부터) snak이 reference_snak에서 갖는 ordinal. 곧 reference_snak의
--     ordinal을 이 출현의 순서로 늘어놓은 순열이다. 이 출현의 snak 목록은
--     reference_snak(hash, dump_id, ordinal = perm[0]), (…, perm[1]), … 이다. SQL에서는
--     ORDER BY coalesce(array_position(o.perm, s.ordinal), s.ordinal + 1)로 읽는다.
--     value_ordinal(같은 속성 안 순서)은 reference_snak 행의 값(첫 출현 기준)이고, 이 출현의 값 순서는
--     perm이 준다.
--   - 행은 순서가 reference_snak과 다른 출현만 둔다(내용이 같고 순서만 다른 출현). 내용이 다른 출현은
--     순열이 없으므로 쓰지 않고 도구가 보고한다.
--   - FK(statement_id, ordinal) → statement_reference는 두지 않는다. FK를 걸면 참조되는 statement_reference
--     (13억 행대)에 RI 트리거와 잠금이 생긴다. 이미 적재된 표의 정의를 바꾸지 않으려는 것이다. 행은
--     도구가 entity_raw 줄에서 만든 출현에서만 나오고, 같은 배치에서 그 줄의 statement_reference가 DB와
--     같은지 대조한다(차이는 보고). 확인 질의(색인 조회):
--       SELECT count(*) FROM reference_snak_order o
--       WHERE NOT EXISTS (SELECT 1 FROM statement_reference r WHERE r.statement_id = o.statement_id AND r.ordinal = o.ordinal);
CREATE TABLE IF NOT EXISTS reference_snak_order (
    statement_id  bigint     NOT NULL,                 -- statement_reference.statement_id
    ordinal       smallint   NOT NULL,                 -- statement_reference.ordinal (그 statement의 references 안 순서)
    perm          smallint[] NOT NULL,                 -- 이 출현 순서대로 reference_snak.ordinal
    PRIMARY KEY (statement_id, ordinal)
);

-- 구간별 재개 지점(repair_reference_snak_checkpoint와 같은 모양). 키는 (dump_id, run)이다. run은 구간
-- 이름(lines:0-0, entities:Q1-Q10000)이라 판마다 같을 수 있다. line_no는 이 줄까지 처리를 마쳤다는
-- 뜻이다. 배치의 행 쓰기와 체크포인트는 한 트랜잭션에 커밋된다. 출현마다 결과가 정해져 있어(같은 순열)
-- 구간이 겹쳐도 같은 행이 된다. 2026-10-29 이전 표(키 run뿐)는 migrate_checkpoint_key.sql로 옮긴다.
-- 작업을 마치고 결과를 확인한 뒤 이 표는 지워도 된다(reference_snak_order는 남긴다).
CREATE TABLE IF NOT EXISTS reference_snak_order_checkpoint (
    dump_id  smallint    NOT NULL,
    run      text        NOT NULL,
    line_no  bigint      NOT NULL,
    updated  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (dump_id, run)
);
