package main

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// go get golang.org/x/sync/errgroup
func main() {
	var group errgroup.Group

	group.Go(func() error {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("hello")
		return nil
	})

	group.Go(func() error {
		time.Sleep(80 * time.Millisecond)
		fmt.Println("hello")
		return errors.New("sorry")
	})

	group.Go(func() error {
		time.Sleep(50 * time.Millisecond)
		fmt.Println("hello")
		return errors.New("ops")
	})

	err := group.Wait()
	fmt.Println("wait result ", err)
}

// cd 04_errgroup/
// go run main.go
// Output:
// hello
// hello
// hello
// wait result  ops
