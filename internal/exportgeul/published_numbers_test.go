//ff:func feature=convert type=helper control=iteration dimension=1
//ff:what TestPublishedNumbers — document folders become numbers; files, other names, bad hex and a missing directory are ignored; an unreadable path fails
package exportgeul

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestPublishedNumbers(t *testing.T) {
	dir := t.TempDir()
	for _, d := range []string{"0x000000000000002a", "0x0000000000000001", "h", "0xzz00000000000000"} {
		os.Mkdir(filepath.Join(dir, d), 0o755)
	}
	os.WriteFile(filepath.Join(dir, "0x0000000000000005"), nil, 0o644)
	got, err := publishedNumbers(dir)
	slices.Sort(got)
	if err != nil || !slices.Equal(got, []uint64{1, 42}) {
		t.Fatalf("%v %v", got, err)
	}
	if got, err := publishedNumbers(filepath.Join(dir, "none")); err != nil || got != nil {
		t.Fatalf("missing: %v %v", got, err)
	}
	if _, err := publishedNumbers(filepath.Join(dir, "0x0000000000000005")); err == nil {
		t.Fatal("file: want error")
	}
}
