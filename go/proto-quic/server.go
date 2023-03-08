package main

import (
	"context"
	"fmt"

	example "github.com/ever0de/playground/proto-quic/proto"
	"github.com/golang/protobuf/proto"
	quic "github.com/quic-go/quic-go"
)

func NewServer() {
	// Start the QUIC server
	listener, err := quic.ListenAddr("localhost:4242", GenerateTLSConfig(), nil)
	if err != nil {
		panic(err)
	}

	// Wait for a client to connect
	session, err := listener.Accept(context.Background())
	println("server)accept")
	if err != nil {
		panic(err)
	}

	// Open a stream for sending and receiving messages
	stream, err := session.AcceptStream(context.Background())
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

	// Receive a protobuf message
	buffer := make([]byte, 1024)
	n, err := stream.Read(buffer)
	if err != nil {
		panic(err)
	}
	receivedMessage := &example.Message{}
	err = proto.Unmarshal(buffer[:n], receivedMessage)
	if err != nil {
		panic(err)
	}
	fmt.Println("Received message from client:", receivedMessage.GetBody())
}
