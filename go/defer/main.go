package main

func main() {
	earlyReturn()

	getFunc(func() {
		println("callback func")

		defer func() {
			println("callback defer")
		}()
	})

	println("main defer")
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
