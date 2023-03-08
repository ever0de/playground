package main

import (
	"context"
	"crypto/tls"
	"fmt"

	example "github.com/ever0de/playground/proto-quic/proto"
	quic "github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

func NewClient() {
	// Connect to the QUIC server
	session, err := quic.DialAddr("localhost:4242", &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}, nil)
	println("client)connect session")
	if err != nil {
		panic(err)
	}

	// Open a stream for sending and receiving messages
	stream, err := session.OpenStreamSync(context.Background())
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
	fmt.Println("Received message from server:", receivedMessage.GetBody())
}
