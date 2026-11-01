//ff:func feature=repair type=writer control=sequence
//ff:what TestDrainPending — applies waiting batches only in order, frees a slot each, and stops at a failed apply
package repair

import (
	"testing"

	"github.com/geul-org/wikipgeul/internal/load"
)

func TestDrainPending(t *testing.T) {
	db, c := loadFixture(t)
	tokens := make(chan struct{}, 3)
	tokens <- struct{}{}
	tokens <- struct{}{}
	tokens <- struct{}{}
	pending := map[int64]Found{0: {Seq: 0, Last: 2}, 1: {Seq: 1, Last: 3}, 3: {Seq: 3, Last: 9}}
	var next int64
	var st Stats
	if err := drainPending(db, c, pending, &next, tokens, load.NewRefSet(64), &st); err != nil || next != 2 || len(pending) != 1 || len(tokens) != 1 || st.Last != 3 {
		t.Fatalf("%d %v %d %+v %v", next, pending, len(tokens), st, err)
	}
	c.DryRun = false
	pending[2] = Found{Seq: 2, Last: 5}
	if err := drainPending(db, c, pending, &next, tokens, load.NewRefSet(64), &st); err == nil || next != 2 {
		t.Fatalf("no work tables: %v", err)
	}
}
