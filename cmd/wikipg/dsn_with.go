//ff:func feature=cli type=helper control=selection
//ff:what dsnWith — a DSN with one more session parameter (key=value), for a postgres:// URL or a key=value DSN
package main

import "strings"

// dsnWith returns dsn with the run-time parameter param (key=value), which
// lib/pq sends at connection start (e.g. jit=off,
// default_transaction_read_only=on): a query parameter of a postgres:// URL,
// a key=value pair otherwise (an empty DSN keeps the PG* variables).
func dsnWith(dsn, param string) string {
	switch {
	case strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://"):
		if strings.Contains(dsn, "?") {
			return dsn + "&" + param
		}
		return dsn + "?" + param
	default:
		return strings.TrimSpace(dsn + " " + param)
	}
}
