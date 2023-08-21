package gnet

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/panjf2000/gnet/v2"
)

type Server struct {
	addr   string
	engine gnet.Engine

	db *pebble.DB
}

func NewServer(addr string) *Server {
	db, err := pebble.Open("./tmp", &pebble.Options{})
	if err != nil {
		panic(err)
	}

	return &Server{
		addr: addr,

		db: db,
	}
}

func (s *Server) Start() {
	handler := NewChunkingHandler(
		func(ctx context.Context, conn gnet.Conn, routingId uint32, bz []byte) gnet.Action {
			key := binary.BigEndian.Uint32(bz[0:4])
			fmt.Printf("server get(%d) key: %d, request length: %d\n", routingId, key, len(bz))

			go func(c gnet.Conn, bz []byte) {
				batch := s.db.NewBatch()
				err := batch.Set(bz[0:4], bz[4:], pebble.Sync)
				if err != nil {
					panic(err)
				}
				if err := batch.Commit(pebble.Sync); err != nil {
					panic(err)
				}

				headerBuf := make([]byte, ProtocolHeaderSize)
				binary.BigEndian.PutUint32(headerBuf[0:4], routingId)
				binary.BigEndian.PutUint32(headerBuf[4:8], uint32(len("pong")))

				if err := conn.AsyncWritev([][]byte{
					headerBuf,
					[]byte("pong"),
				}, nil); err != nil {
					panic(err)
				}
			}(conn, bz)

			return gnet.None
		},
		func(engine gnet.Engine) {
			fmt.Println("server boot")
			s.engine = engine
		},
		func(c gnet.Conn) {
			fmt.Printf("server open connection, remote: %s\n", c.RemoteAddr().String())
		},
		func(c gnet.Conn, err error) gnet.Action {
			fmt.Printf("server close connection: err: %v\n", err)
			return gnet.Close
		},
	)

	err := gnet.Run(
		handler,
		s.addr,
		gnet.WithMulticore(true),
		gnet.WithReuseAddr(true),
		gnet.WithTCPKeepAlive(time.Minute),
	)
	if err != nil {
		panic(err)
	}
}
