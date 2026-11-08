package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Select to Listen Multiple Channel
func ListenMultiWay() {
	ch1 := make(chan int, 1000)
	ch2 := make(chan byte, 1000)

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			ch1 <- rand.Int()
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
			ch2 <- byte(rand.Int())
		}
	}()

loop:
	for {
		time.Sleep(time.Second)
		select { // Listen Multiple Channels, channel got data will get executed first, if both channels have data, randomly choose one
		case v1 := <-ch1:
			fmt.Printf("v1=%v\n", v1)
		case v2 := <-ch2:
			fmt.Printf("v2=%v\n", v2)
			if v2 < 40 {
				break loop
			}
		default:
			fmt.Println("default")
		}
	}

	// Final fallback: if a new element was written to ch1 right at the moment we hit `break AB`,
	// the select below will catch it and print it. Why put it in a select? Because if ch1 is empty,
	// I don't want to block. I want the main goroutine to exit immediately.
	select { // non-blocking way to read from a channel
	case v1 := <-ch1:
		fmt.Printf("at last v1=%d\n", v1)
	default: // Non-blocking way to read channel
	}
}

func SelectBlock() {
	// work goroutine
	go func() {
		time.Sleep(time.Second)
		fmt.Println("I still providing services for outside")
	}()

	// Method 1
	select {} // Block main goroutine infinitely, to avoid the work goroutine get terminated

	// Method 2
	// ch := make(chan struct{})
	// ch <- struct{}{} // Block main goroutine infinitely
}

func main() {
	ListenMultiWay()
	SelectBlock()
}

// Output:
// v2=53
// v1=4958900944720906762
// v2=52
// v2=57
// v1=6815316736877713025
// v2=20
// at last v1=7680512032676148418
// I still providing services for outside
