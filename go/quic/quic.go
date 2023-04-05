package main

import (
	"github.com/quic-go/quic-go"
)

var (
	timeout = &quic.IdleTimeoutError{}
)

func errorIsIdleTimeout(err error) bool {
	return err == timeout
}

func errorIsServerClose(err error) bool {
	return err == quic.ErrServerClosed
}
