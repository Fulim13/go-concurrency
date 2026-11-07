package main

import (
	"fmt"
	"time"
)

var ch = make(chan int, 2) // Buffered Channel

func ChannelBlock() {
	go func() {
		time.Sleep(time.Second)
		if v, ok := <-ch; ok {
			fmt.Println("Take ", v)
		}
		time.Sleep(time.Second)
		if v, ok := <-ch; ok {
			fmt.Println("Take ", v)
		}
	}()

	ch <- 1
	fmt.Println("Send 1")
	ch <- 2
	fmt.Println("Send 2")
	ch <- 3 // Blocking, because buffered channel is full
	fmt.Println("Send 3")
	time.Sleep(2 * time.Second)
}

// Output:
// Send 1
// Send 2
// Take 1 (After 1 Sec)
// Send 3 (After 1 Sec)
// Take 2 (After 2 Sec)

func CloseChannel() {
	ch <- 1
	go func() {
		time.Sleep(3 * time.Second)
		ch <- 2
	}()
	// close(ch) // When close channel, read operation will return immediately with 0 and false
	v, ok := <-ch
	fmt.Println(v, ok)
	v, ok = <-ch // Blocking for 3 sec, if the close channel line is comment
	fmt.Println(v, ok)
	ch <- 3 // After closing the channel, cannot write to channel (panic: send on closed channel)
}

// Output:
// 1 true
// 0 false (If close channel)
// 2 true (If does not close channel, After 3 sec)

func TraverseChannel() {
	ch <- 1
	ch <- 2
	// close(ch) // if we don't close the channel while the for loop keep reading the channel it will cause deadlock

	// if we did not close the channel, but every sec, there is 1 value put in the channel, it won't cause dead lokc
	// go func() {
	// 	for {
	// 		ch <- 10
	// 		time.Sleep(time.Second)
	// 	}
	// }()

	// Method 1
	for ele := range ch {
		fmt.Println(ele)
	}
	// Output:
	// 1
	// 2

	// Method 2
	// for {
	// 	if v, ok := <-ch; ok {
	// 		fmt.Println(v)
	// 	} else {
	// 		break
	// 	}
	// }
	// Output:
	// 1
	// 2
}

func main() {
	// ChannelBlock()
	// CloseChannel()
	TraverseChannel()
}
