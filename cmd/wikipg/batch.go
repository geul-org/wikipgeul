//ff:type feature=cli type=model
//ff:what batch — consecutive dump lines given to one parser: its number, last line, the lines, their numbers and size
package main

// batch is a run of dump lines one worker turns into one transaction: Seq
// numbers batches from 0 in reading order, Last is the number of the last
// line read into it (the checkpoint moves to it when the batch is done), and
// Lines with Nums are the entity lines and their line numbers, Bytes their
// total size.
type batch struct {
	Seq   int64
	Last  int64
	Lines [][]byte
	Nums  []int64
	Bytes int64
}
