//ff:func feature=cli type=helper control=iteration dimension=1
//ff:what langSet — the chosen term languages as a set; nil for all languages
package main

// langSet turns the -langs list into the set load.Config keeps; no list (all
// languages) is nil.
func langSet(langs []string) map[string]bool {
	if len(langs) == 0 {
		return nil
	}
	m := map[string]bool{}
	for _, l := range langs {
		m[l] = true
	}
	return m
}
