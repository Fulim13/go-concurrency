package main

import (
	"context"
	"fmt"
	"time"
)

/*
type Context interface {
    // Deadline() records the expiration time, and whether one is set
    Deadline() (deadline time.Time, ok bool)

    // Done() returns a receive-only channel. No elements are ever sent on it (struct{}),
    // so the channel exists purely to block until it's closed
    Done() <-chan struct{}

    // Err() records why the Done() channel was closed, e.g. because of a timeout,
    // or because it was forcibly cancelled
    Err() error

    // Value() returns the value for a given key. You can think of the Context
    // as maintaining a map internally
    Value(key any) any
}

// emptyCtx can't be accessed directly from outside the package; you have to go through
// Background() or TODO()
type emptyCtx int // implements the Context interface

var (
    background = new(emptyCtx)
    todo       = new(emptyCtx)
)

func Background() Context {
    return background
}

func TODO() Context {
    return todo
}
*/

func Timeout1() {
	// when execute cancel, it will close the read channel (Done() <-chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})

	go func() {
		fmt.Println("Business Logic Execution")
		time.Sleep(100 * time.Millisecond)
		done <- struct{}{}
	}()

	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel() // call cancel, activate Done
	}()

	select {
	case <-done:
		fmt.Println("Business Logic not expired")
		//ctx.Done() is a channel, when called cancel() will close this channel,
		// any read channel operation will immediately return
	case <-ctx.Done():
		err := ctx.Err()                               // if happen Done (Channel get close), Err returned Reason of Done, could be Cancel, could be timeout
		fmt.Println("Business Logic not expired", err) // Output: context cancel
	}
}

func Timeout2() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50) // timeout will auto call context Deadline, and Deadline will call Done
	defer cancel()
	done := make(chan struct{})

	go func() {
		fmt.Println("Business Logic Execution")
		time.Sleep(100 * time.Millisecond)
		done <- struct{}{}
	}()

	select {
	case <-done:
		fmt.Println("Business Logic not expired")
		//ctx.Done() is a channel, when called cancel() will close this channel,
		// any read channel operation will immediately return
	case <-ctx.Done():
		err := ctx.Err()                               // if happen Done (Channel get close), Err returned Reason of Done, could be Cancel, could be timeout
		fmt.Println("Business Logic not expired", err) // Output: context cancel
	}
}

// The Context created by context.WithTimout, its life won't over its parent Context
func InheritTimeout() {
	parent, cancel1 := context.WithTimeout(context.Background(), time.Millisecond*1000) // parent set 1000ms timeout
	t0 := time.Now()
	defer cancel1()

	time.Sleep(500 * time.Millisecond) // Used 500ms

	child, cancel2 := context.WithTimeout(parent, time.Millisecond*1000) // parent left 500ms, child set 1000ms, child.Done() channel will close based on the lesser value, which is 500ms
	// child, cancel2 := context.WithTimeout(parent, time.Millisecond * 100) // parent left 500ms, child set 100ms, child.Done() channel will close based on the lesser value, which is 100ms
	t1 := time.Now()
	defer cancel2()

	<-child.Done()
	fmt.Println(time.Since(t0).Milliseconds(), time.Since(t1).Milliseconds())
	fmt.Println(child.Err()) // context deadline exceeded
}

type StringKey string
type StringValue string
type IntValue int

func RoutineID() {
	for i := 0; i < 3; i++ {
		ctx := context.WithValue(context.Background(), StringKey("gid"), IntValue(1))
		ctx = context.WithValue(ctx, StringKey("owner"), StringValue("Fu Lim"))
		go func(ctx context.Context) {
			if gid, ok := ctx.Value(StringKey("gid")).(IntValue); ok {
				fmt.Printf("Current Goroutine %d\n", gid)
			}
			if owner, ok := ctx.Value(StringKey("owner")).(StringValue); ok {
				fmt.Printf("Current Goroutine %s\n", owner)
			}
		}(ctx)
	}
	time.Sleep(time.Second)
}

// Use context values only for request-scoped data that transits processes and APIs,
// not for passing optional parameters to functions.
