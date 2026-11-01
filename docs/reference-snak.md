# 출처 snak 저장 규칙과 보조 도구

## 문제

위키베이스 출처 hash는 snak 순서(snaks-order, 같은 속성 안 값 순서)와 무관하다. 그래서 같은 hash가 개체마다 다른 순서로 나온다. `reference_snak`은 hash를 키로 하므로 순서 하나만 담을 수 있다.

## 출현 규칙 (결정적)

`reference_snak`은 hash마다 **한 출현의 snak 목록을 통째로** 담는다. 번호마다 다른 출현에서 고르면 목록이 섞인다.

그 출현은 덤프 순서상 첫 출현, 곧 `statement_reference`의 `(statement_id, ordinal)`이 가장 작은 출현이다.

1. `entity_raw.line_no`가 가장 작은 개체,
2. 그 안에서 적재기의 statement 순서(속성 ID 문자열 순, 같은 속성 안은 claims 배열 순; Lexeme → Form → Sense),
3. 그 statement의 references 배열 순서상 첫 번째.

이 규칙은 DB만으로 다시 계산된다.

- 적재기는 `reference_snak`에 스테이징 열 `occ_statement_id`·`occ_ordinal`(그 행을 쓴 출현)을 함께 적재한다.
- `RefSet`은 hash마다 기록한 가장 이른 출현을 기억해, 병렬 파서가 더 이른 출현을 늦게 만나도 다시 쓴다(첫 출현은 반드시 기록한다. 메모리는 `-ref-limit`로 상한).
- 4.2 정리가 hash마다 가장 작은 출현의 행만 남기고 두 열을 뺀다.
- 재개 시 지운 줄의 출현 행도 지운다.

## 출현별 순서: `reference_snak_order`

같은 hash를 다른 순서로 쓴 출현의 원문 순서는 `reference_snak_order`([`sql/reference_snak_order.sql`](../sql/reference_snak_order.sql))가 둔다. `wikipg build-reference-order`가 `entity_raw`에서 만든다. 덤프 파일은 다시 읽지 않는다.

- **표** `reference_snak_order(statement_id, ordinal, perm smallint[])`, PK `(statement_id, ordinal)`. 키는 `statement_reference`와 같다. 출현 하나 = 그 표의 한 행이다.
- `dump_id` 열이 없다. `statement_id`에 판이 들어 있다.
- `perm[i]` = 그 출현의 i번째 snak이 `reference_snak`에서 갖는 `ordinal`. 곧 `reference_snak`의 ordinal을 그 출현 순서로 늘어놓은 순열이다. 같은 내용의 snak이 둘이면 저장 순서대로 짝짓는다.
- 행은 내용이 같고 순서만 다른 출현에만 둔다. 행이 없으면 `reference_snak` 순서 그대로다. 내용이 다른 출현은 순열이 없으므로 쓰지 않고 보고한다.
- 읽기: `ORDER BY coalesce(array_position(o.perm, s.ordinal) - 1, s.ordinal)`. `o`는 `(statement_id, ordinal)` PK 조회다(예: `LEFT JOIN LATERAL (… OFFSET 0)`).
- FK는 두지 않는다. 수십억 행대 `statement_reference`를 참조하는 FK를 걸면 그 표에 RI 트리거와 잠금이 생긴다. 대신 도구가 같은 배치에서 그 줄들의 `statement_reference`를 DB와 대조한다. 확인 질의는 SQL 주석에 있다.
- 적재기 안에서 만들지 않는 이유: 어느 출현이 첫 출현인지는 4.2를 마쳐야 정해지고, 적재 중에 모든 출현의 순서를 남기려면 `statement_reference` 규모의 스테이징이 든다. 도구는 적재기의 행 만들기(`load.Rows.AddEntity`)를 그대로 쓰므로 순서 규칙이 같다.

```sh
GOMEMLIMIT=2GiB wikipg build-reference-order -dry-run -entities Q1-Q10000 -workers 4   # 세기만
GOMEMLIMIT=2GiB wikipg build-reference-order -workers 8 -status 60s                   # 전 구간 쓰기, 같은 명령으로 재개
```

- `entity_raw`를 줄 순서로 배치마다 읽는다(작업자 병렬, 적용은 줄 순서). 출현마다 적재기와 같은 행을 만들어 `reference_snak`과 비교하고(same·order·content·missing), order인 출현의 순열을 쓴다.
- `reference_snak`과 다른 표는 읽기만 한다.
- 배치의 행 쓰기와 체크포인트(`reference_snak_order_checkpoint`)는 한 트랜잭션이라 중단(SIGINT) 뒤 같은 명령으로 정확히 이어진다. 같은 순열이면 행을 다시 쓰지 않는다.

## 복구 도구: `repair-reference-snak`

출현 규칙이 생기기 전에 4.2를 거친 DB(번호마다 따로 고른 `DISTINCT ON (hash, dump_id, ordinal)`로 출현이 섞였을 수 있는 DB)를 고친다.

```sh
wikipg repair-reference-snak -dry-run -entities Q1-Q10000      # 비교만(읽기 전용 세션), 보고
wikipg repair-reference-snak                                    # 전 구간 쓰기, 중단 뒤 같은 명령으로 재개
wikipg repair-reference-snak -from-line 2 -to-line 1000000 -batch 1000 -workers 9 -pause 0 -status 30s
```

- `entity_raw`를 줄 순서로 배치마다 읽어 출현마다 적재기와 같은 행을 만들고, hash마다 범위 안 첫 출현의 목록을 DB와 비교해 다르면 그 hash의 행을 트랜잭션 안에서 바꾼다. 같은 패스에서 `statement_reference`를 원본과 비교해 차이를 보고만 한다.
- 처리한 hash와 출현은 `repair_reference_snak_done`, 재개 지점은 `repair_reference_snak_checkpoint`([`sql/repair_reference_snak.sql`](../sql/repair_reference_snak.sql), 쓰기 실행 때 만든다).
- 구간은 어떤 순서로 돌려도 전 구간을 마치면 hash마다 덤프 첫 출현이 된다. 줄 순서대로 전 구간을 돌리면 hash마다 정확히 한 번 처리한다.
- 행 교체(DELETE + INSERT)는 같은 트랜잭션에서 색인을 함께 갱신하므로 재색인이 필요 없다. 끝난 뒤 `VACUUM (ANALYZE) wikidata.reference_snak`을 돌리고, 작업 표 두 개는 결과를 확인한 뒤 지운다.

## 두 도구의 공통 사항

- 범위: `-from-line`/`-to-line`(기본 전체) 또는 `-entities Q1-Q10000`. 체크포인트는 판과 범위 이름(`lines:0-0`, `entities:Q1-Q10000`)마다 따로다(키 `(dump_id, run)`). 시작할 때 그 판·범위의 체크포인트를 stderr에 적는다.
- `-dry-run`은 `default_transaction_read_only=on` 세션으로 아무것도 쓰지 않는다(처리 표 대신 메모리, `-ref-limit`). 쓰기 실행은 스키마별 advisory lock 하나를 함께 써서 두 도구가 겹치지 않는다.
- 세션은 `jit=off`, `enable_seqscan=off`, `max_parallel_workers_per_gather=0`이다. 질의는 모두 색인 조회이고, 목록은 `unnest(배열) CROSS JOIN LATERAL (… OFFSET 0)`로 준다. 수만 개 hash를 `= ANY(배열)`로 주면 플래너가 행 수를 크게 잡아 거대한 병렬 순차 탐색을 고르기 때문이다.
- 기본 작업자 = CPU의 70%의 절반, 배치 1000줄. 다른 큰 작업(색인 생성 등)과 겹치면 `-workers`를 줄인다.
- 보고: hash별 판정(same·order·content·missing), 출현별 비교, statement_reference 차이.

## 옛 체크포인트 표 옮기기

체크포인트 표의 키는 `(dump_id, run)`이다. 키가 `run`뿐인 옛 표는 다른 판을 같은 구간 이름으로 쓸 때 덮어쓴다. 도구는 시작할 때 키를 확인하고(`checkKey`, dry-run 포함), 옛 키면 이전 SQL을 안내하며 멈춘다.

```sh
PGOPTIONS='-c search_path=wikidata' psql -w -v ON_ERROR_STOP=1 -f sql/migrate_checkpoint_key.sql
```

[`sql/migrate_checkpoint_key.sql`](../sql/migrate_checkpoint_key.sql)은 있는 표만 옮기고(다시 돌려도 된다), 행의 `dump_id`가 `dump`에 있고 `line_no`가 그 판의 마지막 줄을 넘지 않는지 확인한 뒤 PK만 바꾼다.
