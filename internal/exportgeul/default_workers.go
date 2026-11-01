//ff:func feature=convert type=helper control=sequence
//ff:what defaultWorkers — pack-mode conversion goroutines that leave about 30% of the CPUs free, counting a database backend per worker
package exportgeul

import "runtime"

// defaultWorkers is 70% of the CPUs halved, at least 1: each worker keeps a
// PostgreSQL backend about as busy as itself (the queries of one entity),
// and the zstd encoder and the writer need a share too. Memory per worker is
// small (one entity and its document); the encoder holds the most
// (pack.Zstd.Workers).
func defaultWorkers() int {
	return max(1, runtime.NumCPU()*7/10/2)
}
