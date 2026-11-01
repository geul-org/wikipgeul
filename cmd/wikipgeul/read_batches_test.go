//ff:func feature=cli type=command control=sequence
//ff:what TestReadBatches — batches of the given size or byte bound with line numbers; from, limit, a stopped context and a read error behave
package main

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

func TestReadBatches(t *testing.T) {
	f, _ := os.Open(writeDump(t, 5, 0))
	defer f.Close()
	ch := make(chan batch, 10)
	last, err := readBatches(context.Background(), f, 0, 0, 2, 1<<30, ch)
	bs := drain(ch)
	if err != nil || last != 7 || len(bs) != 3 || bs[0].Nums[0] != 2 || bs[2].Last != 6 || bs[2].Seq != 2 || len(bs[2].Lines) != 1 {
		t.Fatalf("all: %d %v %+v", last, err, bs)
	}
	f.Seek(0, io.SeekStart)
	ch = make(chan batch, 10)
	readBatches(context.Background(), f, 4, 2, 10, 1<<30, ch)
	bs = drain(ch)
	if len(bs) != 1 || bs[0].Nums[0] != 4 || bs[0].Last != 5 {
		t.Fatalf("from+limit: %+v", bs)
	}
	f.Seek(0, io.SeekStart)
	ch = make(chan batch, 10)
	readBatches(context.Background(), f, 0, 0, 100, 1, ch)
	if bs = drain(ch); len(bs) != 5 {
		t.Fatalf("by bytes: %d batches", len(bs))
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ch = make(chan batch, 10)
	if _, err := readBatches(ctx, strings.NewReader("[\n{}\n"), 0, 0, 10, 1<<30, ch); !errors.Is(err, context.Canceled) || len(drain(ch)) != 0 {
		t.Fatalf("canceled: %v", err)
	}
	ch = make(chan batch, 10)
	if _, err := readBatches(context.Background(), failReader{}, 0, 0, 10, 1<<30, ch); err == nil || len(drain(ch)) != 0 {
		t.Fatal("read error: want error")
	}
	ch = make(chan batch, 10)
	readBatches(context.Background(), io.MultiReader(strings.NewReader("{\"id\":\"Q1\"}\n"), failReader{}), 0, 0, 10, 1<<30, ch)
	if bs := drain(ch); len(bs) != 1 {
		t.Fatalf("partial batch before error: %+v", bs)
	}
}
