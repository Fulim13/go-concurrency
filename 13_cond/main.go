package main

import (
	"fmt"
	"sync"
	"time"
)

func CondSignal() {
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

	wg := sync.WaitGroup{}

	wg.Go(func() {
		time.Sleep(time.Second)
		cond.Signal() // Only 1 goroutine will receive the singal
		fmt.Println("Send Signal")

		time.Sleep(time.Second)
		cond.Signal() // Only 1 goroutine will receive the singal
		fmt.Println("Send Signal")
	})

	wg.Go(func() {
		cond.L.Lock() // Wait will execute Unlock, so before Wait we need to lock it because we cannot unlock a lock that is not lock
		cond.Wait()
		fmt.Println("Receive signal, execute something")
		cond.L.Unlock() // Need to unlock after wait, because next time another goroutine may need to lock this Lock

		// time.Sleep(2 * time.Second) // With this code, the second Send Signal not able to receive, because Wait have to setup first before able to send
		// the goroutine is Lock and Wait, and the main routine is also in Wait(), cause deadlock
		cond.L.Lock() // Wait will execute Unlock, so before Wait we need to lock it because we cannot unlock a lock that is not lock
		cond.Wait()
		fmt.Println("Receive signal, execute something")
		cond.L.Unlock()
	})

	wg.Wait()
	fmt.Println("===================================")
}

// Send Signal
// Receive signal, execute something
// Send Signal
// Receive signal, execute something
// ===================================

func ChannelSignal() {
	ch := make(chan struct{}, 100)
	// if you make buffered channel, the send won't blocking if the queue is not full, but the receive will block if there is no signal send to channel
	// if you use unbuffered channel or the channel is full, the send side will block, and the receive will block if there is no signal send to channel
	wg := sync.WaitGroup{}

	wg.Go(func() {
		time.Sleep(time.Second)
		ch <- struct{}{}
		fmt.Println("Send Signal")

		time.Sleep(time.Second)
		ch <- struct{}{}
		fmt.Println("Send Signal")
	})

	wg.Go(func() {
		<-ch
		fmt.Println("Receive signal, execute something")

		time.Sleep(2 * time.Second) // Able to Receive Signal After the signal being sent
		<-ch
		fmt.Println("Receive signal, execute something")
	})

	wg.Wait()
	fmt.Println("===================================")
}

// Send Signal
// Receive signal, execute something
// Send Signal
// Receive signal, execute something
// ===================================

func CondBroadcast() {
	mu := sync.Mutex{}
	cond := sync.NewCond(&mu)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		time.Sleep(time.Second)
		cond.Broadcast() // all go routine will receive the signal
		fmt.Println("Broadcast message")

		time.Sleep(time.Second)
		cond.Broadcast() // all go routine will receive the signal
		fmt.Println("Broadcast message")
	})

	for i := 0; i < 3; i++ {
		wg.Go(func() {
			cond.L.Lock()
			cond.Wait()
			fmt.Println("Receive signal, execute something")
			cond.L.Unlock()

			// time.Sleep(2 * time.Second) // Second) // With this code, the second Send Signal not able to receive, because Wait have to setup first before able to send
			// Three of the goroutine is Lock and Wait, and the main routine is also in Wait(), cause deadlock
			cond.L.Lock()
			cond.Wait()
			fmt.Println("Receive signal, execute something")
			cond.L.Unlock()
		})
	}
	wg.Wait()
	fmt.Println("===================================")
}

// Broadcast message
// Receive signal, execute something
// Receive signal, execute something
// Receive signal, execute something
// Broadcast message
// Receive signal, execute something
// Receive signal, execute something
// Receive signal, execute something
// ===================================

func ChannelBroadcast() {
	ch := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Go(func() {
		time.Sleep(time.Second)
		close(ch)
		fmt.Println("Broadcast message")
	})

	for i := 0; i < 3; i++ {
		wg.Go(func() {
			time.Sleep(2 * time.Second)
			<-ch
			fmt.Println("Receive signal, execute something")
		})
	}
	wg.Wait()
	fmt.Println("===================================")
}

// Broadcast message
// Receive signal, execute something
// Receive signal, execute something
// Receive signal, execute something
// ===================================

func main() {
	CondSignal()
	ChannelSignal()
	CondBroadcast()
	ChannelBroadcast()
}
