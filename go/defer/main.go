package main

import "fmt"

func main() {
	earlyReturn()

	getFunc(func() {
		println("callback func")

		defer func() {
			println("callback defer")
		}()
	})

	println("main defer")
	println("=======================")

	for i := 0; i < 3; i++ {
		j := i
		defer func() {
			fmt.Printf("loop defer: %d\n", j)
		}()
	}

	println("loop defer")
}

func getFunc(callback func()) {
	defer func() {
		println("getFunc defer")
	}()

	callback()
}

func earlyReturn() {
	if true {
		return
	}

	defer func() {
		println("earlyReturn defer")
	}()
}
