-- repair_reference_snak.sql — `wikipg repair-reference-snak`의 작업 표 (쓰기 실행 때 도구가 만든다)
--
-- 2026-10-28. 2026-10-27의 wikidata.sql 4.2 정리가 reference_snak을 번호마다 다른 출현에서 골라
-- 출처 snak 목록이 섞였다(표본 검증 E1). 복구 도구는 entity_raw를 줄 순서로 읽어 hash마다 첫 출현
-- (wikidata.sql reference_snak 주석의 출현 규칙)의 snak 목록을 만들고, DB와 다르면 그 hash의 행을
-- 트랜잭션 안에서 바꾼다. 이 두 표는 그 작업의 재개 기록이다. search_path는 도구가 대상 스키마로 둔다.
-- 복구를 끝내고 결과를 확인한 뒤 두 표는 지워도 된다(DROP TABLE).

-- 처리한 hash. 한 hash는 한 번만 처리한다: 여기 있는 출현보다 늦은 출현은 건너뛰고, 이른 출현
-- (앞 구간을 나중에 돌린 경우)이 오면 다시 처리해 덮어쓴다. 그래서 구간을 어떤 순서로 돌려도
-- 전 구간을 마치면 결과는 hash마다 덤프 첫 출현이다. 줄 순서대로 돌리면 hash마다 정확히 한 번이다.
CREATE TABLE IF NOT EXISTS repair_reference_snak_done (
    hash              text     NOT NULL,
    dump_id           smallint NOT NULL,
    occ_statement_id  bigint   NOT NULL,                 -- 처리한 출현 (statement_reference.statement_id)
    occ_ordinal       smallint NOT NULL,                 -- 그 statement의 references 안 순서
    verdict           text     NOT NULL,                 -- same(그대로), order(순서만 달라 교체), content(내용이 달라 교체), missing(행이 없어 넣음)
    PRIMARY KEY (hash, dump_id)
);

-- 구간별 재개 지점. 키는 (dump_id, run)이다. run은 구간 이름(lines:2-121971282, entities:Q1-Q10000)이고
-- 판과 무관하게 같아질 수 있어(전체 구간은 어느 판이든 lines:0-0) 판 번호와 함께 키가 된다. line_no는
-- 이 줄까지 처리를 마쳤다는 뜻이다. 배치의 교체·처리 표·체크포인트는 한 트랜잭션에 커밋된다.
-- 2026-10-29 이전 표는 키가 run뿐이었다. 그런 표는 migrate_checkpoint_key.sql로 옮긴다(도구는 시작할 때
-- 키를 확인하고, 옛 키면 멈춘다).
CREATE TABLE IF NOT EXISTS repair_reference_snak_checkpoint (
    dump_id  smallint    NOT NULL,
    run      text        NOT NULL,
    line_no  bigint      NOT NULL,
    updated  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (dump_id, run)
);
