//ff:func feature=load type=writer control=sequence
//ff:what sampleDump — a DumpInfo for dump 1 of 2026-10-19 with every field set
package load

import "time"

// sampleDump returns the dump row of a test load: dump 1, 2026-10-19, with
// file, URL, size, hashes, times and loader filled in.
func sampleDump() DumpInfo {
	day := time.Date(2026, 10, 19, 0, 0, 0, 0, time.UTC)
	return DumpInfo{ID: 1, Date: day, FileName: "wikidata-20261019-all.json.bz2", URL: "https://dumps.wikimedia.org/x", Size: 103488697311,
		MD5: "cfec", SHA1: "a05b", SHA256: "abcd", PublishedAt: day.Add(50 * time.Hour), DownloadedAt: day.Add(7 * 24 * time.Hour),
		Loader: "github.com/geul-org/wikipgeul/cmd/wikipg", LoaderCommit: "test"}
}
