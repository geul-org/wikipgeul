//ff:func feature=cli type=command control=iteration dimension=1
//ff:what drain — every batch a closed channel still holds
package main

// drain reads ch until it is closed and returns the batches in order.
func drain(ch <-chan batch) []batch {
	var bs []batch
	for b := range ch {
		bs = append(bs, b)
	}
	return bs
}
