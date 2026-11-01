//ff:func feature=cli type=command control=sequence
//ff:what TestDumpInfoOf — name, date, URL, size, download time and sidecar hashes fill in; a bad date or publication time, a missing file and a missing SHA-256 fail
package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDumpInfoOf(t *testing.T) {
	p := writeDump(t, 2, 0)
	dir := filepath.Dir(p)
	os.WriteFile(filepath.Join(dir, "wikidata-20261019-md5sums.txt"), []byte("m5  wikidata-20261019-all.json\nzz  other\n"), 0o644)
	c := loadConfig{In: p}
	if err := dumpInfoOf(&c, "", ""); err != nil {
		t.Fatal(err)
	}
	d := c.Dump
	if d.FileName != "wikidata-20261019-all.json" || d.Date.Format("20060102") != "20261019" || d.URL != "https://dumps.wikimedia.org/wikidatawiki/entities/20261019/wikidata-20261019-all.json" ||
		d.Size == 0 || d.DownloadedAt.IsZero() || d.MD5 != "m5" || d.SHA1 != "" || d.SHA256 != "abc123" || d.Loader == "" {
		t.Fatalf("%+v", d)
	}
	if err := dumpInfoOf(&loadConfig{In: p}, "", "2026-10-21T02:33:31Z"); err != nil {
		t.Fatal(err)
	}
	if err := dumpInfoOf(&loadConfig{In: p}, "", "yesterday"); err == nil {
		t.Fatal("bad published: want error")
	}
	other := filepath.Join(dir, "dump.json")
	os.WriteFile(other, nil, 0o644)
	if err := dumpInfoOf(&loadConfig{In: other}, "", ""); err == nil {
		t.Fatal("no date: want error")
	}
	if err := dumpInfoOf(&loadConfig{In: other}, "20261019", ""); err == nil {
		t.Fatal("no sha256: want error")
	}
	if err := dumpInfoOf(&loadConfig{In: filepath.Join(dir, "none.json")}, "20261019", ""); err == nil {
		t.Fatal("missing file: want error")
	}
}
