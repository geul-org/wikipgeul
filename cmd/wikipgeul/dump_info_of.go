//ff:func feature=cli type=command control=sequence
//ff:what dumpInfoOf — fills the dump record from the file: name, date, URL, size, modification time, sidecar hashes
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

// dumpInfoOf completes c.Dump: the file name, the date from -dump-date or the
// name (wikidata-YYYYMMDD-…), the dumps.wikimedia.org URL when -url is empty,
// the size and download time from the file, the hashes from the sidecar files
// next to it (readSum: <name>.sha256 and wikidata-YYYYMMDD-md5sums.txt,
// -sha1sums.txt) and the loader. A dump with no SHA-256 is refused: the dump
// row must say which bytes were loaded.
func dumpInfoOf(c *loadConfig, date, published string) error {
	d := &c.Dump
	d.Loader = "github.com/geul-org/wikipgeul/cmd/wikipgeul"
	d.FileName = filepath.Base(c.In)
	if date == "" {
		if m := regexp.MustCompile(`wikidata-(\d{8})-`).FindStringSubmatch(d.FileName); m != nil {
			date = m[1]
		}
	}
	day, err := time.Parse("20060102", date)
	if err != nil {
		return fmt.Errorf("wikipgeul: dump date %q: need -dump-date YYYYMMDD", date)
	}
	d.Date = day
	if d.URL == "" {
		d.URL = "https://dumps.wikimedia.org/wikidatawiki/entities/" + date + "/" + d.FileName
	}
	if published != "" {
		if d.PublishedAt, err = time.Parse(time.RFC3339, published); err != nil {
			return fmt.Errorf("wikipgeul: -published: %w", err)
		}
	}
	st, err := os.Stat(c.In)
	if err != nil {
		return err
	}
	d.Size, d.DownloadedAt = st.Size(), st.ModTime()
	dir := filepath.Dir(c.In)
	d.SHA256 = readSum(c.In+".sha256", d.FileName)
	d.MD5 = readSum(filepath.Join(dir, "wikidata-"+date+"-md5sums.txt"), d.FileName)
	d.SHA1 = readSum(filepath.Join(dir, "wikidata-"+date+"-sha1sums.txt"), d.FileName)
	if d.SHA256 == "" {
		return fmt.Errorf("wikipgeul: no SHA-256 for %s (write it to %s.sha256 with sha256sum)", d.FileName, c.In)
	}
	return nil
}
