//ff:func feature=convert type=helper control=iteration dimension=1
//ff:what rangeNumbers — the numbers from through to, inclusive
package exportgeul

// rangeNumbers returns from, from+1, …, to.
func rangeNumbers(from, to uint64) []uint64 {
	var out []uint64
	for n := from; n <= to; n++ {
		out = append(out, n)
	}
	return out
}
