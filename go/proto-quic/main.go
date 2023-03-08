package main

func main() {
	go func() {
		NewServer()
	}()

	NewClient()
}
