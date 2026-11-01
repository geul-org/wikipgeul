//ff:type feature=cli type=model
//ff:what failReader — a reader whose every read fails, for read error paths
package main

// failReader fails every Read with a "disk" error.
type failReader struct{}
