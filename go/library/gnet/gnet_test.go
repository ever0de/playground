package gnet_test

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/panjf2000/gnet/v2"
	"github.com/stretchr/testify/assert"
)

var _ gnet.EventHandler = &Handler{}

type Handler struct {
	eng gnet.Engine

	// buf bytes.Buffer

	onBoot     func(eng gnet.Engine) gnet.Action
	onShutdown func(eng gnet.Engine)
	onOpen     func(c gnet.Conn) ([]byte, gnet.Action)
	onClose    func(c gnet.Conn, err error) gnet.Action
	onTraffic  func(c gnet.Conn) gnet.Action
}

func (h *Handler) OnBoot(eng gnet.Engine) (action gnet.Action) {
	h.eng = eng
	if h.onBoot != nil {
		return h.onBoot(eng)
	}

	return gnet.None
}

func (h *Handler) OnShutdown(eng gnet.Engine) {
	if h.onShutdown != nil {
		h.onShutdown(eng)
		return
	}

}

func (h *Handler) OnOpen(c gnet.Conn) (out []byte, action gnet.Action) {
	if h.onOpen != nil {
		return h.onOpen(c)
	}

	return nil, gnet.None
}

func (h *Handler) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if h.onClose != nil {
		return h.onClose(c, err)
	}

	return gnet.None
}

func (h *Handler) OnTraffic(c gnet.Conn) (action gnet.Action) {
	if h.onTraffic != nil {
		return h.onTraffic(c)
	}

	return gnet.None
}

func (h *Handler) OnTick() (delay time.Duration, action gnet.Action) { return }

func TestGnet(t *testing.T) {
	addr := "localhost:8081"
	serverHandler := new(Handler)
	go func() {
		i := 0
		serverHandler.onOpen = func(c gnet.Conn) ([]byte, gnet.Action) {
			fmt.Printf("Open connection: %s\n", c.RemoteAddr().String())
			return nil, gnet.None
		}
		serverHandler.onTraffic = func(c gnet.Conn) gnet.Action {
			req, _ := io.ReadAll(c)
			fmt.Printf("Receive data: %s\n", string(req))

			i++
			if i%2 == 0 {
				fmt.Printf("Close connection: %s\n", c.RemoteAddr().String())
				if err := c.Close(); err != nil {
					panic(err)
				}
			}

			return gnet.None
		}

		err := gnet.Run(serverHandler, addr)
		if err != nil {
			panic(err)
		}
	}()

	h := new(Handler)

	closeCh := make(chan struct{})
	h.onClose = func(c gnet.Conn, err error) gnet.Action {
		closeCh <- struct{}{}
		return gnet.None
	}

	client, err := gnet.NewClient(h)
	assert.NoError(t, err)
	assert.NoError(t, client.Start())

	write := func(conn gnet.Conn, data string) error {
		_, err := conn.Write([]byte(data))
		return err
	}

	{
		fmt.Println("first connection")
		conn, err := client.Dial("tcp", addr)
		assert.NoError(t, err)

		assert.NoError(t, client.Start())
		assert.NoError(t, write(conn, "hello1\t"))
		assert.NoError(t, write(conn, "hello2\t"))
		assert.NoError(t, write(conn, "hello3\t"))

		assert.NoError(t, conn.Close())
		<-closeCh
		err = write(conn, "errrrrrrrrror\t")
		fmt.Printf("after close, error: %v\n", err)
		assert.Error(t, err)
	}

	{
		fmt.Println("second connection")
		conn, err := client.Dial("tcp", addr)
		assert.NoError(t, err)

		assert.NoError(t, write(conn, "hello1\t"))
		assert.NoError(t, write(conn, "hello2\t"))
		assert.NoError(t, conn.Close())
	}
}
