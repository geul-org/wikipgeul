//ff:func feature=cli type=helper control=sequence
//ff:what failReader.Read — always fails with a disk error
package main

import "errors"

// Read fails with a "disk" error and reads nothing.
func (failReader) Read([]byte) (int, error) { return 0, errors.New("disk") }
