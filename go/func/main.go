package main

func main() {
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
