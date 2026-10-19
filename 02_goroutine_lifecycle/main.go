package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	SubRoutine()
}

// go run 02_goroutine_lifecycle/main.go
// Output:
// child begin
// grandchild begin
// Number of Goroutine (Before): 3
// child finish
// Number of Goroutine (After): 2
// grandchild finish

// main goroutine will use this SubRoutine function
func SubRoutine() {
	go child()
	time.Sleep(5 * time.Second)
}

func child() {
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
	fmt.Println("grandChild begin")
	fmt.Printf("Number of Goroutine (Before): %d\n", runtime.NumGoroutine()) // 3
	time.Sleep(3 * time.Second)
	fmt.Printf("Number of Goroutine (After): %d\n", runtime.NumGoroutine()) // 2
	fmt.Println("grandChild finish")
}
