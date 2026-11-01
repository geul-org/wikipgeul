//ff:func feature=load type=helper control=selection
//ff:what snakValues — the snak columns of a snak: type, datatype, value type, value JSON, entity ID, JSON text jsonb cannot hold
package load

import (
	"bytes"

	"github.com/geul-org/wikipgeul/internal/dump"
)

// snakValues returns the snak columns (snaktype, datatype, value_type,
// value, value_entity, value_json) of s, NULL where absent. The value is the
// datavalue JSON as written. PostgreSQL jsonb rejects the escape \u0000, so
// such a value goes to value_json as JSON text instead (text holds the
// six-character escape): no value is dropped.
func snakValues(s dump.Snak) []any {
	vals := []any{s.SnakType, nullable(s.Datatype), nil, nil, nil, nil}
	dv := s.DataValue
	switch {
	case dv == nil:
	case bytes.Contains(dv.Value, []byte(`\u0000`)):
		vals[2], vals[4], vals[5] = dv.Type, nullable(dump.ValueEntity(dv)), string(dv.Value)
	default:
		vals[2], vals[3], vals[4] = dv.Type, string(dv.Value), nullable(dump.ValueEntity(dv))
	}
	return vals
}
