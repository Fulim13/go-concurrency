package main

import (
	"errors"
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

// New Version
type RateLimiterServer struct {
	sem chan struct{}
}

func NewServer(maxConcurrent int) *RateLimiterServer {
	return &RateLimiterServer{
		sem: make(chan struct{}, maxConcurrent),
	}
}

func (s *RateLimiterServer) processTask() error {
	time.Sleep(1 * time.Second)
	return nil
}

func (s *RateLimiterServer) HandlerRequest() error {
	select {
	// Try to acquire a slot without blocking, if already reach maxConcurrent, it will go to default
	case s.sem <- struct{}{}:
		defer func() {
			// release the slot after done
			<-s.sem
		}()
		return s.processTask()
	default:
		return errors.New("server is busy, please try again later")
	}
}
