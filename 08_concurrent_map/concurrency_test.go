package concurrentmap

import (
	"sync"
	"testing"
)

func TestConcurrency(t *testing.T) {
	cm := NewConcurrentMap[int, bool](50)
	const P = 10
	wg := sync.WaitGroup{}
	for i := 0; i < P; i++ {
		wg.Go(func() {
			for j := 0; j < 1000; j++ {
				cm.Store(j, true)
			}
		})
	}
	wg.Wait()
}

// won't have fatal error under high concurrency write to the map
// cd 08_concurrent_map
// go test -v . -run '^TestConcurrency$'
