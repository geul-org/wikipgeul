//ff:func feature=convert type=loader control=iteration dimension=1
//ff:what validSchema — whether a schema name is a plain lower-case SQL identifier safe to put in a query
package pgsource

// validSchema reports whether s is a plain identifier — a lower-case letter
// or _ followed by lower-case letters, digits or _, at most 63 bytes — so it
// can be written into the SQL unquoted (queriesFor, LoadDump).
func validSchema(s string) bool {
	if s == "" || len(s) > 63 || (s[0] >= '0' && s[0] <= '9') {
		return false
	}
	for _, c := range s {
		if !(c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '_') {
			return false
		}
	}
	return true
}
