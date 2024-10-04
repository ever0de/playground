package gnet

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/panjf2000/gnet/v2"
)

type (
	client struct {
		conn gnet.Conn
		cfg  ClientConfig

		router *Router

		closed *atomic.Bool

		closeCh chan struct{}
	}

	ClientConfig struct {
		// dial target address
		ServerAddr string
		// Timeout option limits the response time of Send by recording a deadline
		// using `context.WithTimeout`.
		Timeout time.Duration

		// == eventrouter ==
		// EventRouterSize is the size of the eventrouter.
		//
		// Default: `16384`
		EventRouterSize uint32
		// == GNET ==
		// TCPKeepAlive sets up a duration for (SO_KEEPALIVE) socket option.
		//
		// Default: `2 minutes`
		KeepAlive time.Duration
		// NumEventLoop is set up to start the given number of event-loop goroutine.
		//
		// Default: `runtime.NumCPU()`
		NumEventLoop int
	}
)

const (
	DefaultEventRouterSize                        = 16384
	DefaultPersistedEventRouterSize               = 32
	DefaultClientKeepAliveDuration  time.Duration = 2 * time.Minute
)

func setupDefaultConfig(cfg ClientConfig) ClientConfig {
	if cfg.EventRouterSize == 0 {
		cfg.EventRouterSize = DefaultEventRouterSize
	}

	if cfg.KeepAlive == 0 {
		cfg.KeepAlive = DefaultClientKeepAliveDuration
	}
	if cfg.NumEventLoop == 0 {
		cfg.NumEventLoop = runtime.NumCPU()
	}

	return cfg
}

func clientProtocolHandler(router *Router) protocolHandler {
	return func(ctx context.Context, peer *peer, bz []byte, closer func()) {
		defer closer()
		routingID := RoutingID(peer.sequence)

		switch routingID.RoutingIDType() {
		case RoutingTypeSingular:
			router.ExecuteRoutingEntry(ctx, routingID, bz)
		default:
			panic("gnet/client: get unknown routingID type(maybe broken connection, get other version of tetrapod)")
		}
	}
}

func NewClient(cfg ClientConfig) (*client, error) {
	cfg = setupDefaultConfig(cfg)
	router := NewRouter(cfg.EventRouterSize)
	closeCh := make(chan struct{}, 1)
	client := &client{
		cfg: cfg,

		router: router,

		closed: &atomic.Bool{},

		closeCh: closeCh,
	}

	openCh := make(chan struct{}, 1)
	handler := newChunkingHandler(clientProtocolHandler(router), false).
		WithOpen(func(c gnet.Conn) {
			fmt.Println("gnet/client: open connection",
				"remote-addr:", c.RemoteAddr(),
				"local-addr:", c.LocalAddr(),
			)

			openCh <- struct{}{}
		}).
		WithClose(func(c gnet.Conn, err error) gnet.Action {
			fmt.Println("gnet/client: closed connection",
				"error:", err,
				"local-addr:", c.LocalAddr(),
			)

			close(closeCh)
			return gnet.Shutdown
		})

	c, err := gnet.NewClient(
		handler,

		gnet.WithNumEventLoop(cfg.NumEventLoop),
		gnet.WithLoadBalancing(gnet.SourceAddrHash),

		gnet.WithTCPKeepAlive(cfg.KeepAlive),
	)
	if err != nil {
		return nil, err
	}

	if err := c.Start(); err != nil {
		return nil, err
	}

	conn, err := c.Dial("tcp", cfg.ServerAddr)
	if err != nil {
		return nil, err
	}

	client.conn = conn

	<-openCh
	return client, nil
}

func (c *client) Send(ctx context.Context, req []byte) ([]byte, error) {
	if c.closed.Load() {
		return nil, fmt.Errorf("gnet/client: connection closed")
	}

	resch := make(chan []byte, 1)
	routingID, err := c.router.NewRoutingEntry(
		func(_ context.Context, routingID RoutingID, response []byte) {
			resch <- response
		},
	)
	if err != nil {
		return nil, err
	}

	headerBuffer := NewProtocolHeader(
		uint32(len(req)),
		uint32(routingID),
	).ToBytes()

	errCh := make(chan error, 1)
	err = c.conn.AsyncWritev(
		[][]byte{headerBuffer, req},
		func(_ gnet.Conn, err error) error {
			if err != nil {
				errCh <- err
			}
			return err
		},
	)
	if err != nil {
		return nil, err
	}

	select {
	case <-c.closeCh:
		return nil, fmt.Errorf("gnet/client: connection closed")
	case response := <-resch:
		return response, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		err := ctx.Err()
		return nil, err
	}
}

func (c *client) Ident() string {
	return c.conn.LocalAddr().String()
}

func (c *client) Close() error {
	return c.conn.Close()
}
