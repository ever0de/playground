package main

import (
	"context"
	"fmt"
	"io"

	example "github.com/ever0de/playground/proto-quic/proto"
	"github.com/golang/protobuf/proto"
	quic "github.com/quic-go/quic-go"
)

func NewServer(addr string) {
	// Start the QUIC server
	listener, err := quic.ListenAddr(addr, GenerateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}

	// Wait for a client to connect
	conn, err := listener.Accept(context.Background())
	println("server)accept")
	if err != nil {
		panic(err)
	}

	// Open a stream for sending and receiving messages
	stream, err := conn.AcceptStream(context.Background())
	println("server)accept stream")
	if err != nil {
		panic(err)
	}

	// Send a protobuf message
	message := &example.Message{
		Body: "Hello from server",
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
	fmt.Println("Received message from client:", receivedMessage.GetBody())
}
