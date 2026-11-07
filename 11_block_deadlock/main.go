package main

import (
	"fmt"
	"sync"
	"time"
)

// 6 Ways to cause blocking

// Method 1: time.Sleep (Will unblock by time)
func Block() {
	go func() {
		time.Sleep(5 * time.Hour)
	}()

	time.Sleep(time.Second)
	fmt.Println("Sleep Over")
}

// Output: Sleep Over (After 1 sec)

// Method 2: WaitGroup
func Block2() {
	go func() {
		time.Sleep(5 * time.Hour)
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	fmt.Println("Wait Over")
}

// Output: Wait Over (won't print)

func Block3() {
	go func() {
		time.Sleep(5 * time.Hour)
	}()

	ch := make(chan int, 10)
	<-ch
	fmt.Println("Channel Over")
}

// Output: Chanel Over (won't print)

func Block4() {
	go func() {
		time.Sleep(5 * time.Hour)
	}()

	ch := make(chan int, 10)
	for ele := range ch {
		fmt.Println(ele)
	}
	fmt.Println("Channel Range Over")
}

// Output: Chanel Range Over (won't print)

func Block5() {
	go func() {
		time.Sleep(5 * time.Hour)
	}()

	mu := sync.Mutex{}
	mu.Lock()
	mu.Lock()
	fmt.Println("Lock Over")
}

// Output: Lock Over (won't print)

func Block6() {
	go func() {
		time.Sleep(5 * time.Hour)
	}()

	select {}
	fmt.Println("Select Over")
}

// Output: Select Over (won't print)

// Deadlock happen when all the current live goroutine all in blocking (except time.Sleep case)
func DeadLock() {
	ch := make(chan int, 10)
	<-ch
	fmt.Println("Channel Over")
}

// fatal error: deadlock

func DeadLock2() {
	go func() {
		time.Sleep(5 * time.Second) // After 5 sec, this go routine ends, main routine is in the blocking cause deadlock
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	fmt.Println("Wait Over")
}

// fatal error: deadlock (After 5 sec)

func DeadLock3() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	fmt.Println("Wait Over")

	// won't reach here, because above already cause blocking
	go func() {
		time.Sleep(5 * time.Second)
	}()
}

// fatal error: deadlock

func main() {
	// Block()
	// Block2()
	// Block3()
	// Block4()
	// Block5()
	// Block6()
	// DeadLock()
	// DeadLock2()
	// DeadLock3()
}
