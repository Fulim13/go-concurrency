package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var wt = sync.WaitGroup{}

func main() {
	SubRoutine()
	fmt.Println("====================")
	time.Sleep(5 * time.Second)
	WaitGroup()
	fmt.Println("====================")
	time.Sleep(5 * time.Second)
	WaitGroup2()
}

// go run  03_waitgroup/main.go
// Output:
// child begin
// grandchild begin
// Number of Goroutine (Before): 3
// child finish
// Number of Goroutine (After): 2
// grandchild finish
//==================== (wait 5 sec)
// Number of Goroutine (Before): 11
// Number of Goroutine (Before): 1
//==================== (wait 5 sec)
// Number of Goroutine (Before): 11
// Number of Goroutine (Before): 1

// main goroutine will use this SubRoutine function
func SubRoutine() {
	wt.Add(2)
	go child()
	wt.Wait()
}

func child() {
	defer wt.Done()
	fmt.Println("child begin")
	go grandChild()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("child finish")
	// After the child goroutine exits,
	// the grandchild goroutine keeps running.
	// All goroutines are equal and independent.
	// One goroutine's lifetime is not controlled by another goroutine (except the main goroutine).
}

func grandChild() {
	defer wt.Done()
	fmt.Println("grandChild begin")
	fmt.Printf("Number of Goroutine (Before): %d\n", runtime.NumGoroutine()) // 3
	time.Sleep(3 * time.Second)
	fmt.Printf("Number of Goroutine (After): %d\n", runtime.NumGoroutine()) // 2
	fmt.Println("grandChild finish")
}

func WaitGroup() {
	const N = 10
	wg := sync.WaitGroup{}
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(x, y int) { // Open N goroutines
			defer wg.Done() // Minus 1 goroutine
			time.Sleep(10 * time.Millisecond)
			_ = x + y
		}(i, i+1)
	}
	fmt.Printf("Number of Goroutine (Before): %d\n", runtime.NumGoroutine()) // N + 1
	wg.Wait()                                                                // Wait all goroutines to be 0
	fmt.Printf("Number of Goroutine (After): %d\n", runtime.NumGoroutine())  // 1
}

func WaitGroup2() {
	const N = 10
	wg := sync.WaitGroup{}
	for i := 0; i < N; i++ {
		wg.Go(func() {
			time.Sleep(10 * time.Millisecond)
		})
	}
	fmt.Printf("Number of Goroutine (Before): %d\n", runtime.NumGoroutine()) // N + 1
	wg.Wait()                                                                // Wait all goroutines to be 0
	fmt.Printf("Number of Goroutine (After): %d\n", runtime.NumGoroutine())  // 1
}

// go 1.25
// wg.Go(func() {
//
// })
// wg.Wait()
// same as

// wg.Add(1)
// go func() {
// 	defer wg.Done()
//	f()
//}()
// wg.Wait()
