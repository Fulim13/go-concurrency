package main

import (
	"fmt"
	"runtime"
	"time"
)

func Add(x, y int) int {
	fmt.Println("Add")
	return x + y
}

var product = func(x, y int) int {
	fmt.Println("Product")
	return x * y
}

func SimpleGoroutine() {
	fmt.Printf("Number of logical processors (CPU Cores): %d\n", runtime.NumCPU())      // 15
	runtime.GOMAXPROCS(runtime.NumCPU() / 2)                                            // set the max number of CPU Cores that current process can use
	fmt.Printf("Number of existing goroutines (Initial): %d\n", runtime.NumGoroutine()) // 1 (Main function)

	// call function
	go Add(1, 2)

	// Anonymous go function
	go func(x, y int) int {
		fmt.Println("Subtract")
		return x - y
	}(3, 1)

	// call function definition
	go product(2, 4)

	// create more goroutines so that goroutines > CPU cores (still can run, because all those goroutines will run on these 15 CPU Cores)
	for i := 0; i < 20; i++ {
		go Add(1, 2)
	}

	fmt.Printf("Number of existing goroutines: %d\n", runtime.NumGoroutine()) // 24
	time.Sleep(1 * time.Millisecond)
}

func main() {
	SimpleGoroutine()
}
// go run 01_simple_goroutine/main.go
// Output:
// Number of logical processors (CPU Cores): 15
// Number of existing goroutines (Initial): 1
// Number of existing goroutines: 24
// Add ... Product  Subtract (Order might different)
