package main

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

func main() {
	globalError()
	errorAs()
}

func globalError() {
	err := doSomething()

	if err == errors.New("not found") {
		panic("unreachable")
	}

	if errors.Is(err, ErrNotFound) {
		fmt.Printf("%#v\n", err)
		return
	}

	panic("unreachable")
}

type errClose struct{}

func (errClose) Error() string {
	return "close"
}

func errorAs() {
	var err error = &errClose{}
	if e := (&errClose{}); errors.As(err, &e) {
		fmt.Printf("%#v\n", e)
		return
	}

	panic("unreachable")
}

func doSomething() error {
	return ErrNotFound
}
