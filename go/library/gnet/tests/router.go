package gnet

import (
	"context"
	"fmt"
	"sync/atomic"
)

type (
	Router struct {
		size           uint32
		routingEntries []RoutingCallback
		lastRoutingId  *atomic.Uint32
	}

	RoutingCallback func(ctx context.Context, id uint32, bz []byte)
)

func NewRouter(size uint32) *Router {
	return &Router{
		size:           size,
		routingEntries: make([]RoutingCallback, size),
		lastRoutingId:  &atomic.Uint32{},
	}
}

func (r *Router) nextRoutingID() uint32 {
	return r.lastRoutingId.Add(1) % r.size
}

func (r *Router) NewRoutingEntry(cb RoutingCallback) uint32 {
	nextId := r.nextRoutingID()

	if r.routingEntries[nextId] != nil {
		panic(fmt.Sprintf("routing entry already exists: %d", nextId))
	}

	r.routingEntries[nextId] = cb
	return nextId
}

func (r *Router) ExecuteRoutingEntry(ctx context.Context, id uint32, bz []byte) {
	r.routingEntries[id](ctx, id, bz)
}
