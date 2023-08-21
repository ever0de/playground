package gnet

import (
	"context"
	"encoding/binary"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/panjf2000/gnet/v2"
)

type Client struct {
	conn gnet.Conn

	router *Router
}

func NewClient(addr string) *Client {
	openCh := make(chan struct{}, 1)

	router := NewRouter(1024)
	handler := NewChunkingHandler(
		func(ctx context.Context, conn gnet.Conn, routingId uint32, bz []byte) gnet.Action {
			router.ExecuteRoutingEntry(ctx, routingId, bz)
			return gnet.None
		},
		func(engine gnet.Engine) {
			fmt.Println("client boot event")
		},
		func(c gnet.Conn) {
			fmt.Println("client open event")
			openCh <- struct{}{}
		},
		func(c gnet.Conn, err error) gnet.Action {
			fmt.Printf("client close event: %v\n", err)
			return gnet.Shutdown
		},
	)

	c, err := gnet.NewClient(
		handler,
		gnet.WithMulticore(true),
		gnet.WithTCPKeepAlive(time.Minute),
	)
	if err != nil {
		panic(err)
	}

	if err := c.Start(); err != nil {
		panic(err)
	}

	conn, err := c.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("client: waiting for open")
	<-openCh
	fmt.Println("client: open")
	return &Client{
		conn:   conn,
		router: router,
	}
}

func (c *Client) Send(key, value []byte) error {
	respCh := make(chan []byte, 1)
	closeFlag := &atomic.Bool{}
	closeFlag.Store(false)
	defer func() {
		closeFlag.Store(true)
		close(respCh)
	}()

	routingID := c.router.NewRoutingEntry(func(ctx context.Context, id uint32, bz []byte) {
		if closeFlag.Load() {
			return
		}

		respCh <- bz
	})
	fmt.Printf("client send routingID: %d\n", routingID)

	headerBuf := make([]byte, ProtocolHeaderSize)
	binary.BigEndian.PutUint32(headerBuf[0:4], routingID)
	binary.BigEndian.PutUint32(headerBuf[4:8], uint32(len(key)+len(value)))

	errCh := make(chan error, 1)
	defer close(errCh)

	err := c.conn.AsyncWritev([][]byte{
		headerBuf,
		append(key, value...),
	}, func(c gnet.Conn, err error) error {
		if err != nil {
			errCh <- err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	select {
	case err := <-errCh:
		return err
	case resp := <-respCh:
		fmt.Printf("client get response: %s\n", resp)
		return nil
	}
}
