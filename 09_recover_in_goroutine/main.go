package main

import (
	"time"
)

var table = []int{1, 2, 3, 4, 5, 6, 7}

// index out of bound if provide index > 7
func GetHandler(index int) int {
	// recover only catch the panic from curent go routine
	// So now it catch the panic of index out of bound
	defer func() {
		recover()
	}()
	return table[index]
}

func main() {
	// The panic cause by GetHandler (another goroutine) and SetHandler(another goroutine) won't catch by this main goroutine
	// recover only catch the panic from curent go routine
	defer func() {
		recover()
	}()
	go GetHandler(10)
	time.Sleep(time.Second)
}

// go run 09_recover_in_goroutine/main.go
