package main

import (
	"context"
	"fmt"
	"io"

	example "github.com/ever0de/playground/proto-quic/proto"
	quic "github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

func NewServerStream(addr string) {
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

	go func() {
		for {
			// Open a stream for sending and receiving messages
			stream, err := conn.AcceptStream(context.Background())
			println("server)accept stream")
			if err != nil {
				panic(err)
			}

			go func() {
				// Receive a protobuf message
				buf, err := io.ReadAll(stream)
				if err != nil {
					panic(err)
				}

				receivedMessage := &example.Operation{}
				err = proto.Unmarshal(buf, receivedMessage)
				if err != nil {
					panic(err)
				}
				fmt.Println("Received message from client:", receivedMessage)

				switch receivedMessage.Type {
				case example.OpType_OpGet:
					message := &example.GetResponse{
						Key:   append(receivedMessage.GetGet().Key, []byte("/key")...),
						Value: []byte("Hello from server/value"),
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

				default:
					panic("unknown message type")
				}
			}()
		}
	}()
}
