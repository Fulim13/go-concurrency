package main

import (
	"fmt"
	"sync"
	"time"
)

var n int32 = 0
var lock sync.Mutex // sync.Mutex is a Write Lock

func inc() {
	lock.Lock()   // Add Write Lock
	n++           // Anytime, only one goroutine can go to this critical region
	lock.Unlock() // Release Write Lock
}

func Lock() {
	const P = 1000 // typically need around 1,000 concurrent goroutines to test the lost update
	wg := sync.WaitGroup{}
	for i := 0; i < P; i++ {
		wg.Go(func() {
			inc()
		})
	}
	wg.Wait()
	fmt.Printf("Lock Finally n=%d\n", n) // Lock Finally n=1000
}

var (
	mu = sync.RWMutex{}
)

// Read locks are reentrant: the same goroutine can acquire the same read lock multiple times (without releasing the previous one)
// because is the same goroutine access the RLock
func ReentranceRLock(n int) {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println(n)
	if n > 0 {
		ReentranceRLock(n - 1)
	}
	time.Sleep(1 * time.Second)
}

// Write locks are not reentrant
// deadlock
func ReentranceWLock(n int) {
	mu.Lock()
	defer mu.Unlock()
	fmt.Println(n)
	if n > 0 {
		ReentranceWLock(n - 1)
	}
	time.Sleep(1 * time.Second)
}

func RLockExclusion() {
	mu.RLock() // Get Read Lock,
	// other go routine that get read lock can get the read lock,
	// but other goroutine get write lock, cannot
	defer mu.RUnlock()

	go func() {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("child goroutine get the write lock")
	}()

	time.Sleep(1 * time.Second)
}

// Nothing will print out

func RLockExclusion1() {
	mu.RLock() // Get Read Lock,
	// other go routine that get read lock can get the read lock,
	// but other goroutine get write lock, cannot
	defer mu.RUnlock()

	go func() {
		time.Sleep(100 * time.Millisecond)
		mu.RLock()
		defer mu.RUnlock()
		fmt.Println("child goroutine get the read lock")
	}()

	time.Sleep(1 * time.Second)
}

// child goroutine get the read lock

func RLockExclusion2() {
	mu.RLock() // Get Read Lock,
	// other go routine that get read lock can get the read lock,
	// but other goroutine get write lock, cannot
	defer mu.RUnlock()

	go func() {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("child goroutine get the write lock")
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		mu.RLock()
		defer mu.RUnlock()
		fmt.Println("child goroutine get the read lock")
	}()

	time.Sleep(1 * time.Second)
}

// Nothing will print out, because the write lock cause blocking, so read lock also cannot get

func RLockExclusion3() {
	mu.RLock() // Get Read Lock,
	// other go routine that get read lock can get the read lock,
	// but other goroutine get write lock, cannot

	go func() {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("child goroutine get the write lock")
	}()

	go func() {
		time.Sleep(100 * time.Millisecond)
		mu.RLock()
		defer mu.RUnlock()
		fmt.Println("child goroutine get the read lock")
	}()

	time.Sleep(1 * time.Second)
	mu.RUnlock()
	time.Sleep(1 * time.Second)
}

// child goroutine get the write lock
// child goroutine get the read lock

func WLockExclusion() {
	mu.Lock() // Get Write Lock,
	// other goroutines can't acquire a read lock or a write lock
	defer mu.Unlock()

	go func() {
		mu.RLock()
		defer mu.RUnlock()
		fmt.Println("child goroutine get the read lock")
	}()

	go func() {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("child goroutine get the write lock")
	}()

	time.Sleep(1 * time.Second)
}

// Nothing will print out

func main() {
	// Lock()
	// Lock Finally n=1000

	// ReentranceRLock(3)
	// 3 2 1 0

	// ReentranceWLock(3)
	// 3
	// fatal error: all goroutines are asleep - deadlock!
	// panic can be recoverable, but fatal error cannot

	// RLockExclusion() // nothing

	// RLockExclusion1() // child goroutine get the read lock

	// RLockExclusion2() // nothing

	// RLockExclusion3()
	// child goroutine get the write lock
	// child goroutine get the read lock
	// write get first, because read wait 100ms

	WLockExclusion() // nothing
}

// go run 06_rwlock/main.go
