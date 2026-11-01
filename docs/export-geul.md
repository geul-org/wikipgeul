# export-geul: 적재본 → GEUL

`wikipgeul export-geul`은 wikipgeul가 적재한 스키마(기본 `wikidata`)를 읽어 위키데이터 개체를 GEUL 문서로 옮긴다. DB는 읽기만 한다.

```sh
wikipgeul export-geul -root ./archive -from 1 -to 200000000 \
  -commit <wikigeul 커밋> -lib-commit <gogeul 커밋> -report report.tsv
```

- **변환 코드는 [wikigeul](https://github.com/geul-org/wikigeul) 한 벌이다.** wikipgeul는 DB 읽기(`internal/pgsource`)와 작업 생성(`internal/exportgeul`)만 갖고, 변환(`wikigeul/pkg/port`)과 pack 쓰기(`wikigeul/pkg/pipeline`)는 import한다. 덤프 파일을 바로 읽는 `wikigeul`과 입력만 다르다(`pipeline.Loader`). 같은 덤프·커밋·레지스트리 판·`-recorded-at`이면 문서가 바이트 단위로 같다. 예외와 대응표는 wikigeul [docs/dump.md](https://github.com/geul-org/wikigeul/blob/main/docs/dump.md)에 있다.
- **`-commit`은 wikigeul 커밋이다.** 문서의 TRANSFORMED-BY는 변환 도구 wikigeul(`0x0001000000000001`)과 gogeul을 적는다. wikipgeul 커밋은 적지 않는다(적재 커밋은 `dump.loader_commit`). wikipgeul가 빌드에 쓴 wikigeul·gogeul 커밋을 준다.
- **이름은 옮기기 전과 같다.** 보고서 첫 줄 `# wikigeul -from … -schema … -dump-id …`, 기본 이상점 보고서 `wikigeul-anomaly.tsv`, 공유 코드 오류 문구의 `wikigeul:` 접두, 체크포인트 형식은 그대로다. 이전 체크포인트를 그대로 재개할 수 있다.
- **플래그:** 공통 플래그는 wikigeul과 같다(`-root`, `-from`, `-to`, `-ids`, `-commit`, `-lib-commit`, `-report`, `-workers`, `-checkpoint`, `-resume`, `-recorded-at`, `-samples`, `-strict`, `-anomaly`, `-allow`, `-langreg-root`, `-langreg-version`, 뜻은 wikigeul README). 여기에만 있는 것은 아래 셋과 `-out`이다.

| 플래그 | 뜻 | 기본값 |
|---|---|---|
| `-dsn` | PostgreSQL DSN | 비면 `PG*` 환경 변수 |
| `-schema` | 적재된 덤프의 스키마(적재 때 `-schema`) | `wikidata` |
| `-dump-id` | 적재된 덤프 번호(적재 때 `-dump-id`, 1–255) | `1` |
| `-out` | `pack`(Issuer pack 세트, 스트림) 또는 `dir`(디렉터리 트리, 시험용) | `pack` |

- `-workers` 기본값은 CPU 70%의 절반이다. 작업자마다 PostgreSQL 백엔드 하나가 거의 같은 만큼 바쁘다.
- 엄격 모드가 이상점에서 멈추면 종료 코드 3이다.

## 디렉터리 모드 (`-out dir`)

시험용이다. Q1–Q1000을 옮기고 그 문서들이 쓴 속성을 이어서 옮겨 `<root>/.well-known/geul/` 디렉터리 트리에 쓴 뒤 Issuer의 이름 색인 `documents.geul`(`q42`, `p31`, …)을 쓴다. 색인은 이미 공표된 문서를 유지하므로 범위를 나눠 실행할 수 있다. 실행마다 문서의 버전 1을 덮어쓰고, 내용이 같으면 파일을 건드리지 않으며, Edge ID는 이전 공표본에서 이어받아 바뀌지 않은 statement의 번호를 지킨다. 같은 입력·커밋이면 pack 모드와 문서 바이트가 같다.

## 읽는 DB 테이블

DB 없이 덤프 파일을 바로 읽는 경로에서 각 열을 덤프 JSON의 어디에서 얻는지는 wikigeul [docs/dump.md](https://github.com/geul-org/wikigeul/blob/main/docs/dump.md)에 있다. 테이블 정의는 `sql/wikidata.sql`이다.

| 테이블 | 쓰는 열 |
|---|---|
| `dump` | `url`, `sha256`, `downloaded_at`, `loader`, `loader_commit` (출처 기록) |
| `entity` | `entity_id`(있는지), `lastrevid`·`modified`(개체 판, 출처 기록), `datatype`(속성) |
| `label`, `description`, `alias` | `language`, `value` (모든 언어, 언어 코드 바이트 순. 별칭은 그 안에서 `ordinal` 순) |
| `statement` | `statement_id`, `property`, `rank`, `snaktype`, `value_type`, `value`(jsonb), `value_entity` |
| `qualifier` | `statement_id`, `ordinal`(qualifiers-order), `property`, `snaktype`, `value_type`, `value`, `value_entity` |
| `statement_reference`, `reference_snak` | 출처: statement의 출처 순서·hash, 출처 항목(속성, 값 열은 qualifier와 같다) |
| `reference_snak_order` (있으면) | 출현별 출처 항목 순서: `(statement_id, ordinal)` → `perm`(wikipgeul `build-reference-order`). 표가 없거나 그 출현에 행이 없으면 `reference_snak.ordinal` 순서 |
| `sitelink` | `site`, `title`(`title_json`), `badges` (사이트 ID 바이트 순 = 덤프 키 순) |

## 읽기 규칙

- 모든 질의는 `dump_id`로 한 판만 읽는다(`-dump-id`).
- `statement_id`는 wikipgeul가 매긴다. 속성 사이는 속성 ID 문자열 순, 같은 속성 안은 덤프 순이다. 위키데이터에서 속성 사이 순서는 뜻이 없다.
- 값은 개체면 맨 ID(`value_entity`), 문자열이면 따옴표 없는 문자열, 그 밖은 datavalue JSON(jsonb의 공백·키 순서)이다.
- `\u0000`이 든 값은 wikipgeul가 `value_json`(레이블·별칭·설명도 `value_json`)에 원문 JSON으로 두므로 그것을 읽는다(`source.ExactValue`).
- novalue·somevalue는 값이 비어 있고 `snaktype` 열로 구별한다(자리 표지로 옮긴다, [변환 규칙](conversion.md)). 값 snak인데 값이 빈 것만 건너뛴다(보고서 `빈 값`).
- 본 값과 qualifier 값은 같은 형태다.

## 적재 도구 확인

`dump.loader`가 wikipgeul(`github.com/geul-org/wikipgeul/cmd/wikipgeul`, 또는 wikipgeul의 옛 이름으로 기록된 경로)가 아니면 실행을 거부한다. wikipgeul는 덤프를 무손실로 저장하므로(덤프 줄 원문 `entity_raw`, 구조화 표가 못 담는 값은 `*_json` 열) 문서 내용은 덤프·wikigeul·GEUL 라이브러리로 정해진다. 다른 적재 도구는 무손실을 알 수 없다.

## 출현별 출처 항목 순서

- `reference_snak`은 hash마다 덤프상 첫 출현의 순서 하나만 둔다. 위키베이스 출처 hash는 snak 순서와 무관해서, 같은 hash의 다른 출현은 순서가 다를 수 있다.
- 그런 출현은 wikipgeul의 `reference_snak_order`가 순열 `perm`(reference_snak ordinal을 그 출현 순서로)을 둔다. `LEFT JOIN LATERAL`로 PK를 한 번 조회해 `coalesce(array_position(perm, ordinal) - 1, ordinal)` 순으로 읽는다(`pgsource.refOrderSQL`).
- 행이 없는 출현은 `reference_snak.ordinal` 순서다. 표가 없는 스키마(옛 적재)도 그렇다. `pgsource.LoadDump`가 `to_regclass`로 표가 있는지 보고, 실행 첫 줄(stdout)에 `reference snak order: per occurrence (reference_snak_order) true|false`로 알린다.

## 세션

DB는 연결 하나에 `SET jit = off`로 쓴다. 출처 질의(statement_reference ⋈ reference_snak)의 행 수를 플래너가 크게 잡아 JIT 컴파일이 개체마다 수백 ms 걸리기 때문이다(끄면 1ms 안쪽).
