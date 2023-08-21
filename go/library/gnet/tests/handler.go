package gnet

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/panjf2000/gnet/v2"
)

type (
	customCtxKey int

	ProtocolHandler func(ctx context.Context, conn gnet.Conn, routingId uint32, bz []byte) gnet.Action

	ChunkingHandler struct {
		handler         ProtocolHandler
		onBootCallback  func(gnet.Engine)
		onOpenCallback  func(gnet.Conn)
		onCloseCallback func(gnet.Conn, error) gnet.Action
	}
)

var _ gnet.EventHandler = (*ChunkingHandler)(nil)

const (
	ctxRoutingIdKey customCtxKey = iota
)
const (
	ProtocolHeaderSize int = 4 + 4
)

func NewChunkingHandler(
	handler ProtocolHandler,
	onBootCallback func(gnet.Engine),
	onOpenCallback func(gnet.Conn),
	onCloseCallback func(gnet.Conn, error) gnet.Action,
) *ChunkingHandler {
	return &ChunkingHandler{
		handler:         handler,
		onBootCallback:  onBootCallback,
		onOpenCallback:  onOpenCallback,
		onCloseCallback: onCloseCallback,
	}
}

func (h *ChunkingHandler) OnBoot(eng gnet.Engine) (action gnet.Action) {
	h.onBootCallback(eng)

	return gnet.None
}

func (h *ChunkingHandler) OnShutdown(eng gnet.Engine) {
}

func (h *ChunkingHandler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	h.onOpenCallback(c)

	return nil, gnet.None
}

func (h *ChunkingHandler) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	h.onCloseCallback(c, err)

	return gnet.None
}

func (h *ChunkingHandler) OnTraffic(c gnet.Conn) (action gnet.Action) {
	if c.InboundBuffered() < ProtocolHeaderSize {
		return gnet.None
	}

	for c.InboundBuffered() >= ProtocolHeaderSize {
		// [routingId, payloadSize]
		headerBuffer, err := c.Peek(ProtocolHeaderSize)
		if err != nil {
			fmt.Printf("peek error: %v(from %s)\n", err, c.RemoteAddr().String())
			return gnet.Close
		}

		routingId := binary.BigEndian.Uint32(headerBuffer[0:4])
		payloadSize := binary.BigEndian.Uint32(headerBuffer[4:8])

		if c.InboundBuffered() < ProtocolHeaderSize+int(payloadSize) {
			return gnet.None
		}

		_, err = c.Discard(ProtocolHeaderSize)
		if err != nil {
			fmt.Printf("discard error: %v(from %s)\n", err, c.RemoteAddr().String())
			return gnet.Close
		}

		bz, err := c.Next(int(payloadSize))

		if err != nil {
			fmt.Printf("next error: %v(from %s)\n", err, c.RemoteAddr().String())
			return gnet.Close
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, ctxRoutingIdKey, routingId)
		action = h.handler(ctx, c, routingId, bz)
		if action != gnet.None {
			return action
		}
	}

	return gnet.None
}

func (g *ChunkingHandler) OnTick() (delay time.Duration, action gnet.Action) { return }
