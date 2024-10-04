package gnet

import (
	"bytes"
	"context"
	"sync"
	"time"

	"github.com/panjf2000/gnet/v2"
)

var _ gnet.EventHandler = (*chunkingHandler)(nil)

type (
	protocolHandler func(ctx context.Context, peer *peer, bz []byte, closer func())

	chunkingHandler struct {
		handler protocolHandler

		poolVersion     bool
		bytesBufferPool *sync.Pool

		onBoot  func(gnet.Engine)
		onOpen  func(gnet.Conn)
		onClose func(gnet.Conn, error) gnet.Action
	}
)

func newChunkingHandler(handler protocolHandler, poolVersion bool) *chunkingHandler {
	return &chunkingHandler{
		handler: handler,

		poolVersion: poolVersion,
		bytesBufferPool: &sync.Pool{
			New: func() any {
				initBuf := make([]byte, 0, 1024)
				return bytes.NewBuffer(initBuf)
			},
		},

		onBoot:  func(gnet.Engine) {},
		onOpen:  func(gnet.Conn) {},
		onClose: func(gnet.Conn, error) gnet.Action { return gnet.Close },
	}
}
func (g *chunkingHandler) WithBoot(callback func(gnet.Engine)) *chunkingHandler {
	g.onBoot = callback
	return g
}
func (g *chunkingHandler) WithOpen(callback func(gnet.Conn)) *chunkingHandler {
	g.onOpen = callback
	return g
}
func (g *chunkingHandler) WithClose(callback func(gnet.Conn, error) gnet.Action) *chunkingHandler {
	g.onClose = callback
	return g
}

// ======== GNET engine implementation ========

func (g *chunkingHandler) OnBoot(eng gnet.Engine) (action gnet.Action) {
	g.onBoot(eng)

	return gnet.None
}

func (g *chunkingHandler) OnShutdown(eng gnet.Engine) {
}

func (g *chunkingHandler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	g.onOpen(c)

	return nil, gnet.None
}

func (g *chunkingHandler) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	return g.onClose(c, err)
}

func (g *chunkingHandler) OnTraffic(c gnet.Conn) (action gnet.Action) {
	if c.InboundBuffered() < ProtocolHeaderSize {
		return gnet.None
	}

	for c.InboundBuffered() >= ProtocolHeaderSize {
		pheader, err := ProtocolHeader{}.
			FromPeek(c.Peek(ProtocolHeaderSize))
		if err != nil {
			println("gnet: error while reading protocol header, try closing connection",
				"error", err,
				"remote-addr", c.RemoteAddr(),
			)
			return gnet.Close
		}

		if c.InboundBuffered() < int(ProtocolHeaderSize+pheader.PacketLength) {
			return gnet.None
		}

		_, err = c.Discard(ProtocolHeaderSize)
		if err != nil {
			println("gnet: error while discarding protocol header, try closing connection",
				"error", err,
				"remote-addr", c.RemoteAddr(),
			)
			return gnet.Close
		}

		bz, err := c.Next(int(pheader.PacketLength))
		if err != nil {
			println("gnet: error while reading packet, try closing connection",
				"error", err,
				"remote-addr", c.RemoteAddr(),
			)
			return gnet.Close
		}

		length := len(bz)
		if g.poolVersion { // use sync.Pool -> error(corrupted data) in concurrent_test.go
			copyBuf := g.bytesBufferPool.Get().(*bytes.Buffer)
			copyBuf.Reset()
			copyBuf.Grow(length)
			copyBuf.Write(bz)

			g.handler(
				context.Background(),
				newPeer(
					c,
					pheader.Sequence,
				),
				copyBuf.Bytes(),
				func() {
					copyBuf.Reset()
					g.bytesBufferPool.Put(copyBuf)
				},
			)
		} else { // use copy -> no error in concurrent_test.go
			copyBuf := make([]byte, length)
			copy(copyBuf, bz)

			g.handler(
				context.Background(),
				newPeer(
					c,
					pheader.Sequence,
				),
				copyBuf,
				func() {},
			)
		}
	}

	return gnet.None
}

func (g *chunkingHandler) OnTick() (delay time.Duration, action gnet.Action) { return }
