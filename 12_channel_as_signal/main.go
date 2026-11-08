package main

import (
	"fmt"
	"time"
)

func Broadcast() {
	ch := make(chan struct{})

	const P = 3
	for i := 0; i < P; i++ {
		go func() {
			<-ch // Blocking
			fmt.Printf("%d Start\n", i)
		}()
	}

	time.Sleep(2 * time.Second)
	fmt.Println("Guys we can start")
	close(ch) // broadcast
	time.Sleep(time.Second)
}

// Output:
// Guys we can start (After 2 seconds)
// 0 Start
// 1 Start
// 2 Start

func CountDownLatch() {
	ch := make(chan struct{})

	const P = 3
	for i := 0; i < P; i++ {
		go func() {
			time.Sleep(time.Duration(i) * time.Second)
			fmt.Printf("%d completed its job\n", i)
			ch <- struct{}{}
		}()
	}

	// Act like WaitGroup
	for i := 0; i < P; i++ {
		<-ch
	}

	fmt.Println("All Completed, I will start now")
}

// 0 completed its job (After 1 sec)
// 1 completed its job (After 2 sec)
// 2 completed its job (After 3 sec)
// All Completed, I will start now (After 3 sec)

func main() {
	// Broadcast()
	CountDownLatch()
}
