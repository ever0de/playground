package gnet

import (
	"context"
	"fmt"
	"net"
	"runtime"
	"strings"
	"time"

	"github.com/panjf2000/gnet/v2"
)

type (
	listener struct {
		config ListenerConfig

		addr   net.Addr
		engine gnet.Engine
	}

	ListenerConfig struct {
		Addr string

		KeepAlive    time.Duration
		NumEventLoop int
	}
)

const (
	DefaultListenerKeepAliveDuration time.Duration = time.Minute
)

func NewListener(config ListenerConfig) (*listener, error) {
	if config.KeepAlive == 0 {
		config.KeepAlive = DefaultListenerKeepAliveDuration
	}
	if config.NumEventLoop == 0 {
		config.NumEventLoop = runtime.NumCPU()
	}

	addr, err := net.ResolveTCPAddr("tcp", config.Addr)
	if err != nil {
		return nil, err
	}

	return &listener{
		config: config,
		addr:   addr,
		engine: gnet.Engine{},
	}, nil
}

func listenerProtocolHandler(handler func(context.Context, *peer, []byte) ([]byte, error)) protocolHandler {
	return func(ctx context.Context, p *peer, bz []byte, closer func()) {
		go func(ctx context.Context, peer *peer, bz []byte) {
			response, err := handler(ctx, peer, bz)
			if err != nil {
				println("gnet/listener: error handling request, try closing peer",
					"error", err,
					"ident", peer.Ident(),
				)

				if err := peer.Close(); err != nil {
					println("gnet/listener: error closing peer",
						"error", err,
					)
				}
				return
			}

			if len(response) == 0 {
				return
			}

			if err := peer.Send(response, closer); err != nil {
				println("error sending response",
					"error", err,
				)

				if err := peer.Close(); err != nil {
					println("error closing peer",
						"error", err)
				}
				return
			}
		}(ctx, p, bz)
	}
}

func (l *listener) AcceptAll(
	h func(context.Context, *peer, []byte) ([]byte, error),
) error {
	listeningStarted := make(chan error)

	handler := newChunkingHandler(listenerProtocolHandler(h), true).
		WithBoot(func(e gnet.Engine) {
			fmt.Println("gnet/listener: engine booted")

			l.engine = e
			listeningStarted <- nil
		}).
		WithOpen(func(conn gnet.Conn) {
			fmt.Println("gnet/listener: open connection",
				"remoteAddr:", conn.RemoteAddr(),
			)
		}).
		WithClose(func(conn gnet.Conn, err error) gnet.Action {
			fmt.Println(
				"gnet/listener: connection closed",
				"remoteAddr:", conn.RemoteAddr(),
				"error:", err,
			)

			return gnet.Close
		})

	go func() {
		addr := l.config.Addr
		if strings.Count(addr, "://") != 1 {
			addr = "tcp://" + addr
		}

		err := gnet.Run(
			handler,
			addr,

			gnet.WithNumEventLoop(l.config.NumEventLoop),
			gnet.WithLoadBalancing(gnet.SourceAddrHash),

			gnet.WithReuseAddr(true),
			gnet.WithTCPKeepAlive(l.config.KeepAlive),
		)

		if err != nil {
			listeningStarted <- err
		}
	}()

	return <-listeningStarted
}

func (l *listener) LocalAddr() net.Addr {
	return l.addr
}

func (l *listener) Close() error {
	return l.engine.Stop(context.Background())
}
