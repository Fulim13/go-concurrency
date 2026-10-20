package main

import (
	"fmt"
	"sync"
)

// Arrays, slices, and structs can be modified concurrently without crashing,
// though this may produce dirty writes (a data race).
// Concurrent writes to a map, however, can (not always) trigger a fatal error (fatal error: concurrent map writes).
// recover() can only catch a panic. It cannot catch a fatal error, which terminates the process immediately.
// If you need to modify a map concurrently, use sync.Map, or guard a regular map with a sync.Mutex / sync.RWMutex.

type Student struct {
	Name string
	Age  int32
}

var arr = [10]int{}
var mp = make(map[int]bool, 10)
var m = sync.Map{}

func CollectionSafety() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Go(func() { // Write Event Index
		for i := 0; i < len(arr); i += 2 {
			arr[i] = 0
		}
	})

	wg.Go(func() { // Write Odd Index
		for i := 1; i < len(arr); i += 2 {
			arr[i] = 1
		}
	})

	wg.Wait()
	fmt.Println(arr) // [0 1 0 1 0 1 0 1 0 1]
	fmt.Println("========================")

	var student Student
	wg.Go(func() {
		student.Name = "Fu Lim"
	})

	wg.Go(func() {
		student.Age = 24
	})

	wg.Wait()
	fmt.Printf("%v %v\n", student.Name, student.Age) // Fu Lim 24
	fmt.Println("========================")

	// fatal error: concurrent map writes
	// wg.Go(func() {
	// 	for i := 0; i < 10000; i += 2 {
	// 		mp[i] = true
	// 	}
	// })

	// wg.Go(func() {
	// 	for i := 1; i < 10000; i += 2 {
	// 		mp[i] = false
	// 	}
	// })

	// wg.Wait()
	// fmt.Printf("%v\n", mp[1])
	// fmt.Println("========================")

	wg.Go(func() {
		m.Store("k1", "v1")
	})
	wg.Go(func() {
		m.Store("k1", "v2")
	})

	wg.Wait()
	if v, exists := m.Load("k1"); exists {
		fmt.Println(v)
	}
}

func main() {
	CollectionSafety()
}
