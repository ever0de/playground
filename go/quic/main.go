package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/quic-go/quic-go"
)

const (
	timeoutDuration = 1 * time.Second

	c = "client"
	s = "server"
)

var (
	serverUDPAddr = &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 4242,
	}

	clientUDPAddr = &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 4241,
	}
)

func main() {
	closeChan := make(chan struct{})
	var wg sync.WaitGroup

	log("global", "udp")
	{
		wg.Add(1)
		go func() { serverUDP(closeChan); wg.Done() }()
		clientUDP(closeChan)
		wg.Wait()
	}

	println("\n\n\n\n")
	log("global", "direct quic")
	{
		go serverQuic(closeChan)
		clientQuic(closeChan)
	}
}

func server(listener quic.EarlyListener) {
	for {
		log(s, "!!!!!!!! waiting for connection !!!!!!!!")
		conn, err := listener.Accept(context.Background())
		if err != nil {
			if errorIsServerClose(err) {
				log(s, "server closed")
				return
			}

			log(s, fmt.Sprintf("accept connection failed: %v", err))
			continue
		}

		log(s, fmt.Sprintf("accept connection from %s", conn.RemoteAddr().String()))
		go listen(conn)
	}
}

var (
	serverQuicConfig = &quic.Config{
		// MaxIdleTimeout: timeoutDuration,
	}
	clientQuicConfig = &quic.Config{
		MaxIdleTimeout: timeoutDuration,
	}
)

func serverUDP(closeChan chan struct{}) {
	conn := getUDPConn(serverUDPAddr)
	quicListener, err := quic.ListenEarly(
		conn,
		GenerateServerTLSConfig(),
		serverQuicConfig,
	)
	if err != nil {
		panic(err)
	}
	log(s, "create quic server")

	go server(quicListener)

	<-closeChan
	if err := quicListener.Close(); err != nil {
		panic(err)
	}
	if err := conn.Close(); err != nil {
		panic(err)
	}
}

func serverQuic(closeChan chan struct{}) {
	quicListener, err := quic.ListenAddrEarly(
		serverUDPAddr.String(),
		GenerateServerTLSConfig(),
		serverQuicConfig,
	)
	if err != nil {
		panic(err)
	}
	log(s, "create quic server")

	go server(quicListener)
	<-closeChan
	if err := quicListener.Close(); err != nil {
		panic(err)
	}
}

func clientUDP(closeChan chan struct{}) {
	type UDPConn struct {
		conn *net.UDPConn
	}

	conn := getUDPConn(clientUDPAddr)
	uc := &UDPConn{conn}
	now := time.Now()
	// timeout
	{
		quicConn := quicDial(uc.conn)
		log(c, "create quic client")
		time.Sleep(timeoutDuration * 2)

		ctx := context.Background()
		_, err := quicConn.OpenStreamSync(ctx)
		if err == nil {
			panic("should not open stream")
		}
		log(c, fmt.Sprintf("open stream failed: %v\ttimeout error: %v", err, errorIsIdleTimeout(err)))
		fmt.Printf("time cost: %v\n", time.Since(now))
		if err := quicConn.CloseWithError(99999, "quic client close"); err != nil {
			panic(err)
		}
		if err := uc.conn.Close(); err != nil {
			panic(err)
		}
	}

	// reconnect
	{
		conn := getUDPConn(clientUDPAddr)
		uc.conn = conn
		quicConn := quicDial(uc.conn)

		ctx := context.Background()
		stream, err := quicConn.OpenStreamSync(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("time cost: %v\n", time.Since(now))

		_, err = stream.Write([]byte("hello"))
		if err != nil {
			panic(err)
		}

		err = stream.Close()
		if err != nil {
			panic(err)
		}

		resp, err := io.ReadAll(stream)
		if err != nil {
			panic(err)
		}

		log(c, fmt.Sprintf("response: %s", resp))

		err = quicConn.CloseWithError(99999, "quic client close")
		if err != nil {
			panic(err)
		}
		if err := quicConn.CloseWithError(99999, "quic client close"); err != nil {
			panic(err)
		}

		if err := uc.conn.Close(); err != nil {
			panic(err)
		}
	}

	closeChan <- struct{}{}
}

func clientQuic(closeChan chan struct{}) {
	dial := func() quic.Connection {
		conn, err := quic.DialAddrEarly(
			serverUDPAddr.String(),
			GenerateClientTlsConfig(),
			clientQuicConfig,
		)
		if err != nil {
			panic(err)
		}
		return conn
	}

	now := time.Now()
	// timeout
	{
		conn := dial()
		log(c, "create quic client")
		time.Sleep(timeoutDuration * 2)

		ctx := context.Background()
		_, err := conn.OpenStreamSync(ctx)
		if err == nil {
			panic("should not open stream")
		}
		log(c, fmt.Sprintf("open stream failed: %v\ttimeout error: %v", err, errorIsIdleTimeout(err)))
		fmt.Printf("time cost: %v\n", time.Since(now))
		if err := conn.CloseWithError(99999, "quic client close"); err != nil {
			panic(err)
		}
	}

	// reconnect
	{
		conn := dial()
		ctx := context.Background()
		stream, err := conn.OpenStreamSync(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("time cost: %v\n", time.Since(now))

		_, err = stream.Write([]byte("hello"))
		if err != nil {
			panic(err)
		}

		err = stream.Close()
		if err != nil {
			panic(err)
		}

		resp, err := io.ReadAll(stream)
		if err != nil {
			panic(err)
		}

		log(c, fmt.Sprintf("response: %s", resp))

		err = conn.CloseWithError(99999, "quic client close")
		if err != nil {
			panic(err)
		}
		if err := conn.CloseWithError(99999, "quic client close"); err != nil {
			panic(err)
		}
	}

	closeChan <- struct{}{}
}

func listen(conn quic.Connection) {
	ctx := context.Background()

	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			if errorIsServerClose(err) {
				log(s, "server closed")
				return
			}

			log(s, fmt.Sprintf("accept stream failed: %v", err))
			return
		}
		log(s, "accept stream")

		request, err := io.ReadAll(stream)
		if err != nil {
			log(s, fmt.Sprintf("read stream failed: %v", err))
			continue
		}

		_, err = stream.Write(request)
		if err != nil {
			log(s, fmt.Sprintf("write stream failed: %v", err))
			continue
		}

		err = stream.Close()
		if err != nil {
			log(s, fmt.Sprintf("close stream failed: %v", err))
		}
	}
}

func getUDPConn(addr *net.UDPAddr) *net.UDPConn {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	log("global", fmt.Sprintf("udp connect port: %d", addr.Port))
	return conn
}

func quicDial(conn *net.UDPConn) quic.Connection {
	quicConn, err := quic.DialEarly(
		conn,
		serverUDPAddr,
		clientUDPAddr.String(),
		GenerateClientTlsConfig(),
		clientQuicConfig,
	)
	if err != nil {
		panic(err)
	}
	log("client", "create quic client")
	return quicConn
}

func log(prefix string, msg any) {
	fmt.Printf("[%s]: [%s]: %s\n", time.Now().Format(time.RFC3339Nano), prefix, msg)
}
