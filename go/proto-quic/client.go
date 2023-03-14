package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"sync"

	example "github.com/ever0de/playground/proto-quic/proto"
	quic "github.com/quic-go/quic-go"
	"google.golang.org/protobuf/proto"
)

func NewClientStream(addr string) {
	// Connect to the QUIC server
	conn, err := quic.DialAddr(addr, &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}, nil)
	println("client)connect session")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Open a stream for sending and receiving messages
	stream, err := conn.OpenStream()
	println("client)open stream")
	if err != nil {
		panic(err)
	}

	stream2, err := conn.OpenStream()
	println("client)open stream2")
	if err != nil {
		panic(err)
	}

	go func() {
		message := &example.Operation{
			Type: example.OpType_OpGet,
			Payload: &example.Operation_Get{
				Get: &example.Get{Key: []byte("Hello from client")},
			},
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
		buf, err := io.ReadAll(stream)
		if err != nil {
			panic(err)
		}
		receivedMessage := &example.GetResponse{}
		err = proto.Unmarshal(buf, receivedMessage)
		if err != nil {
			panic(err)
		}
		fmt.Println("Received message from server(stream1):", receivedMessage)

		wg.Done()
	}()

	go func() {
		message := &example.Operation{
			Type: example.OpType_OpGet,
			Payload: &example.Operation_Get{
				Get: &example.Get{Key: []byte("Hello from client2")},
			},
		}
		data, err := proto.Marshal(message)
		if err != nil {
			panic(err)
		}
		_, err = stream2.Write(data)
		if err != nil {
			panic(err)
		}

		err = stream2.Close()
		if err != nil {
			panic(err)
		}

		buf2, err := io.ReadAll(stream2)
		if err != nil {
			panic(err)
		}
		receivedMessage := &example.GetResponse{}
		err = proto.Unmarshal(buf2, receivedMessage)
		if err != nil {
			panic(err)
		}
		fmt.Println("Received message from server(stream2):", receivedMessage)

		wg.Done()
	}()

	wg.Wait()
}
