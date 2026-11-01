//ff:func feature=convert type=helper control=iteration dimension=1
//ff:what sendAll — call send for each job in order
package exportgeul

import "github.com/geul-org/wikigeul/pkg/pipeline"

// sendAll passes every job of js to send, in order.
func sendAll(js []pipeline.PackJob, send func(pipeline.PackJob)) {
	for _, j := range js {
		send(j)
	}
}
