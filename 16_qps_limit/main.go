package main

import (
	"fmt"
	"sync"
	"time"
)

// Why we want to rate list the request
// 1. Too many concurrent request may cause CPU high load
// 2. Inside the handler, they might initialize a large map/slice that will cause Memory high load
// 3. Inside the handler, it calls Redis/Postgres, if too many concurrent request, Redis/Postgres CPU will high load

var qps = make(chan struct{}, 100)

func handler() {
	qps <- struct{}{}
	fmt.Println("Received Request")
	defer func() {
		<-qps
		fmt.Println("Blocking UnLock")
	}()
	time.Sleep(3 * time.Second)
}

func main() {
	const P = 1000
	wg := sync.WaitGroup{}
	for i := 0; i < P; i++ {
		wg.Go(func() {
			handler()
		})
	}
	wg.Wait()
}
