//ff:type feature=cli type=model
//ff:what result — what a worker reports for a batch: which batch, its last line, entities and errors written, or the failure
package main

// result is a worker's report on batch Seq ending at line Last: the entities
// and load errors it wrote, or Err when the batch could not be written.
type result struct {
	Seq      int64
	Last     int64
	Entities int64
	Errors   int64
	Err      error
}
