package main

import (
	"crypto/tls"
	"fmt"
	"io"

	example "github.com/ever0de/playground/proto-quic/proto"
	"github.com/golang/protobuf/proto"
	quic "github.com/quic-go/quic-go"
)

func NewClient(addr string) {
	// Connect to the QUIC server
	conn, err := quic.DialAddr(addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}, nil)
	println("client)connect session")
	if err != nil {
		panic(err)
	}

	// Open a stream for sending and receiving messages
	stream, err := conn.OpenStream()
	println("client)open stream")
	if err != nil {
		panic(err)
	}

	// Send a protobuf message
	message := &example.Message{
		Body: "Hello from client",
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}
	_, err = stream.Write(data)
	if err != nil {
		panic(err)
	}
	err = stream.Close()
	if err != nil {
		panic(err)
	}

	// Receive a protobuf message
	buf, err := io.ReadAll(stream)
	if err != nil {
		panic(err)
	}
	receivedMessage := &example.Message{}
	err = proto.Unmarshal(buf, receivedMessage)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received message from server:", receivedMessage.GetBody())
}
