package gnet

import (
	"github.com/panjf2000/gnet/v2"
)

type peer struct {
	conn     gnet.Conn
	ident    string
	sequence uint32
}

func newPeer(
	conn gnet.Conn,
	sequence uint32,
) *peer {
	var ident string
	if conn.RemoteAddr() != nil {
		ident = conn.RemoteAddr().String()
	} else {
		ident = "unknown, already closed connection"
	}

	return &peer{
		conn:     conn,
		ident:    ident,
		sequence: sequence,
	}
}

func (c *peer) Send(bz []byte, closer func()) error {
	pHeaderBz := NewProtocolHeader(uint32(len(bz)), c.sequence).
		ToBytes()

	errCh := make(chan error, 1)

	err := c.conn.AsyncWritev(
		[][]byte{pHeaderBz, bz},
		func(_ gnet.Conn, err error) error {
			closer()

			errCh <- err
			if err != nil {
				println("listener/peer.Send: error sending response",
					"error", err,
				)
			}

			return err
		},
	)
	if err != nil {
		return err
	}

	return <-errCh
}

func (c *peer) Ident() string {
	return c.ident
}

func (c *peer) Close() error {
	return c.conn.Close()
}
