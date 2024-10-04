package gnet

import (
	"context"
	"fmt"
)

const (
	MaxRoutingEntries = uint32(1 << 31)

	RoutingTypeSingular = iota
	RoutingTypePersisted
	RoutingTypeUnknown
)

type Router struct {
	routerSize     uint32
	routingEntries []RoutingEntryCallback
	routingQueue   chan uint32
}

type RoutingEntryCallback func(ctx context.Context, routingID RoutingID, bz []byte)

type RoutingID uint32

func (r RoutingID) Singular() RoutingID {
	return r & 0x7fff_ffff
}

func (r RoutingID) Persisted() RoutingID {
	return r | 0x8000_0000
}

func (r RoutingID) String() string {
	switch r.RoutingIDType() {
	case RoutingTypeSingular:
		return fmt.Sprintf("RoutingID.Singular(%d)", r.Singular())
	case RoutingTypePersisted:
		return fmt.Sprintf("RoutingID.Persisted(%d)", r.Persisted())
	}

	return fmt.Sprintf("RoutingID.Unknown(%d)", r)
}

func (r RoutingID) RoutingIDType() int {
	switch r >> 31 {
	case 0x0: // no MSB set
		return RoutingTypeSingular
	case 0x1: // only response MSB is set
		return RoutingTypePersisted
	default:
		return RoutingTypeUnknown
	}
}

func (r RoutingID) Mod(m uint32) RoutingID {
	return RoutingID(uint32(r) % m)
}

func NewRouter(routerSize uint32) *Router {
	if routerSize == 0 {
		panic("routerSize must be greater than 0")
	}
	if routerSize > MaxRoutingEntries {
		panic(fmt.Sprintf("routerSize is too large: %d, maximum: %d", routerSize, MaxRoutingEntries))
	}

	routingQueue := make(chan uint32, routerSize)
	for i := uint32(0); i < routerSize; i++ {
		routingQueue <- i
	}

	return &Router{
		routerSize:     routerSize,
		routingEntries: make([]RoutingEntryCallback, routerSize),
		routingQueue:   routingQueue,
	}
}

func (r *Router) nextRoutingID() RoutingID {
	id := <-r.routingQueue

	return RoutingID(id).
		Mod(r.routerSize).
		Singular()
}

func (r *Router) NewRoutingEntry(callback RoutingEntryCallback) (RoutingID, error) {
	nextRoutingId := r.nextRoutingID()

	if r.routingEntries[nextRoutingId] != nil {
		return 0, fmt.Errorf(
			"routing entry already exists for routingID: %d",
			nextRoutingId,
		)
	}

	r.routingEntries[nextRoutingId] = callback
	return nextRoutingId, nil
}

func (r *Router) ExecuteRoutingEntry(ctx context.Context, routingID RoutingID, bz []byte) {
	r.routingEntries[routingID](ctx, routingID, bz)
	r.RemoveRoutingEntry(routingID)
}

func (r *Router) RemoveRoutingEntry(routingID RoutingID) {
	r.routingEntries[routingID] = nil
	r.routingQueue <- uint32(routingID)
}
