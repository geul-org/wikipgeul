# 적재 스키마와 절차

정의는 [`sql/wikidata.sql`](../sql/wikidata.sql)이다. 원칙과 열 설명은 그 파일 주석에 있다. 이 문서는 요점과 적재 절차를 적는다.

## 원칙

1. **원본을 먼저 보존한다.** 개체마다 덤프의 JSON 한 줄을 그대로 `entity_raw`에 둔다. 나머지 테이블은 그 원본에서 파생한 질의용 색인이다. 해석 규칙에 버그가 있어도 재적재 없이 원본에서 다시 파생한다.
2. **적재 단계에서 아무것도 버리지 않는다.** novalue·somevalue(`snaktype`), qualifier 순서(`ordinal`, qualifiers-order), 출처(`reference*`)를 남긴다. 실패한 줄은 `load_error`에 남긴다.
3. **판본을 추정하지 않는다.** 덤프 파일·URL·크기·MD5·SHA-1·SHA-256·적재 도구와 커밋을 `dump` 행에 기록하고, 개체 단위 테이블은 모두 `dump_id`를 가진다.

- jsonb·text가 못 담는 값(`\u0000`)은 JSON 원문을 `value_json`(statement·qualifier·reference_snak·label·description·alias·lemma·form_representation·sense_gloss)과 `title_json`(sitelink)에 둔다.
- 판을 둘 이상 넣어 비교할 수 있다. 키는 (개체, dump_id, …) 순이라 한 판 안의 개체 조회가 색인 하나로 끝난다.

## 테이블

`dump`, `entity_raw`(덤프 줄 원문), `load_error`, `entity`, `label`, `description`, `alias`, `sitelink`, `statement`, `qualifier`, `reference`, `statement_reference`, `reference_snak`, Lexeme용 `lexeme`, `lemma`, `form`, `form_representation`, `sense`, `sense_gloss`.

- `statement_id` = `dump << 55 | 줄 번호 << 20 | 개체 안 순번`. 개체 안 순서는 속성 ID 문자열 순, 같은 속성 안은 claims 배열 순이다(Lexeme은 Lexeme → Form → Sense).
- `statement_guid` 색인은 UNIQUE가 아니다. 위키데이터 원본에 같은 GUID가 두 statement에 쓰인 경우가 있다(예: Q34433114 P953). statement의 키는 `statement_id`이고 GUID는 찾기용이다. `sitelink_site_title`도 같은 이유로 UNIQUE가 아니다.
- 출처 hash와 snak 순서의 관계는 [reference-snak.md](reference-snak.md).

## Lexeme

Lexeme은 별도 덤프(`wikidata-YYYYMMDD-lexemes.json.bz2`)라 다른 판 번호로 적재한다(`-dump-id 2` 등). Form·Sense도 `entity`에 행(type `form`·`sense`, kind L)을 두어 그 statement(`entity_id` = `L4-F1`)가 가리킬 개체가 있다. 표제어·어형 표기는 낱말 자체라 `-langs`로 거르지 않고, 뜻풀이(gloss)는 거른다. Lexeme·Form·Sense의 statement 키는 한 줄 안에서 이어진다.

## 적재 절차

1. `sql/wikidata.sql`을 실행한다. 1~3절(타입·테이블)만 실행되고 4절은 주석이다. 스키마 이름은 파일 안 `CREATE SCHEMA`·`search_path` 두 줄이 정한다(기본 `wikidata`). 다른 스키마에 넣으려면 그 두 줄을 바꾸고 도구에 `-schema`를 준다.
2. `wikipg`로 덤프를 적재한다. 대량 적재 중에는 색인이 없어 질의가 순차 스캔이다. 적재 중 시험 질의는 하지 않는다.
3. 4절을 주석을 풀어 실행한다.
   - 4.0 무손실 열: `value_json`·`title_json` 열이 없던 옛 스키마에만 필요하다.
   - 4.1 재개: 중단 뒤 재개했다면 체크포인트 이후 줄의 행이 일부 커밋됐을 수 있다. `-from-line` 없이 재개하면 도구가 먼저 지운다.
   - 4.2 출처 중복 정리: 적재기는 최근에 본 hash만 기억해 거르므로(`-ref-limit`) 남은 중복을 PK 전에 지운다. `reference_snak`은 hash마다 첫 출현의 행을 통째로 남긴다([reference-snak.md](reference-snak.md)).
   - 4.3 제약과 색인(PK, 색인, FK NOT VALID → VALIDATE, ANALYZE).
4. 4.4: `wikipg build-reference-order`로 출현별 출처 snak 순서를 만든다([reference-snak.md](reference-snak.md)).

## 기본 자원 사용

기본값은 기계를 약 30% 남긴다. CPU의 70% 중 1/5을 압축 해제(lbzip2), 2/5를 JSON 파싱에 쓰고, 작성자(DB 연결) 4개, 배치 8 MiB(메모리 약 1 GB)다. 작성자를 늘리면 테이블 확장에서 경합하고, 배치를 키우면 메모리만 는다.
