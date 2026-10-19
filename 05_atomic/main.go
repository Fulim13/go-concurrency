package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Data Races and Atomic Operations

// When multiple goroutines concurrently access and modify the same shared memory, it causes a data race.

// n++ is not atomic.
// It’s a read-modify-write operation made of 3 steps:
// load n, increment it, and store the result back.
// When goroutines interleave these steps, increments get overwritten, causing lost updates.
// The race is intermittent, so you typically need around 1,000 concurrent goroutines to reproduce it.

// Tip: run your program with go run -race or go test -race
// to make the race detector report data races,
// rather than relying on 1,000 goroutines to show the bug.

// n is global variable, a and b are not
// Thread 1		Thread 2
// a=n		a=0(1)	a=n		a=0(2)
// b=a+1	b=1(3)	b=a+1		b=1(5)
// n=b		n=1(4)	n=b		n=1(6)
// end up, n = 1, although 2 n++

// Using the sync/atomic package makes the increment atomic, which eliminates the data race and prevents lost updates:
// atomic.AddInt32(addr *int32, delta int32) (new int32): atomically adds delta and returns the new value
// atomic.LoadInt32(addr *int32) (val int32): atomically reads the current value

var n int32 = 0

func inc1() {
	n++ // Not an atomic operation, 3 step, get n value, plus 1, result to n
}

func inc2() {
	atomic.AddInt32(&n, 1) // atomic operation
}

func NotAtomic() {
	const P = 1000 // typically need around 1,000 concurrent goroutines to test the lost update
	wg := sync.WaitGroup{}
	for i := 0; i < P; i++ {
		wg.Go(func() {
			inc1()
		})
	}
	wg.Wait()
	fmt.Printf("NotAtomic Finally n=%d\n", n)
}

func Atomic() {
	const P = 1000 // typically need around 1,000 concurrent goroutines to test the lost update
	wg := sync.WaitGroup{}
	for i := 0; i < P; i++ {
		wg.Go(func() {
			inc2()
			// _ = n + 1
			// _ = atomic.LoadInt32(&n) + 1 // Even under high concurrency, there's no need to use atomic.Load to read a variable
		})
	}
	wg.Wait()
	fmt.Printf("Atomic Finally n=%d\n", n)
}

func main() {
	NotAtomic()
	fmt.Println("===========")
	n = 0 // Reset n
	Atomic()
}

// go run 05_atomic/main.go
// NotAtomic Finally n=959 (Vary)
// ===========
// Atomic Finally n=1000

// go run -race 05_atomic/main.go
// Found 2 data race(s)
