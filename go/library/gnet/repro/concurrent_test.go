package gnet

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newClient(t *testing.T, serverAddr string) *client {
	client, err := NewClient(ClientConfig{
		ServerAddr:      serverAddr,
		EventRouterSize: 1024,
	})
	assert.NoError(t, err)
	return client
}

func TestConcurrentSending(t *testing.T) {
	serverAddr := "0.0.0.0:8080"

	listener, err := NewListener(ListenerConfig{
		Addr: serverAddr,
	})
	assert.NoError(t, err)

	err = listener.AcceptAll(
		func(ctx context.Context, lpeer *peer, bz []byte) ([]byte, error) {
			return bz, nil
		},
	)
	assert.NoError(t, err)
	defer func() { assert.NoError(t, listener.Close()) }()

	client := newClient(t, serverAddr)
	defer func() { assert.NoError(t, client.Close()) }()

	resp, err := client.Send(context.Background(), []byte("first connect\n"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("first connect\n"), resp)

	cnt := 1_000
	list := make([]time.Duration, cnt)
	wg := new(sync.WaitGroup)
	ctx := context.Background()

	now := time.Now()
	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func(i int) {
			req := []byte(fmt.Sprintf("connect %d\n", i))

			now := time.Now()
			var resp []byte
			var err error
			{
				resp, err = client.Send(ctx, req)
			}
			since := time.Since(now)
			list[i] = since

			assert.NoError(t, err)
			assert.Equal(t, req, resp)

			wg.Done()
		}(i)
	}
	wg.Wait()

	total := time.Since(now)
	min := list[0]
	max := list[0]
	avg := time.Duration(0)
	for _, v := range list {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		avg += v
	}

	t.Log("total", total)
	t.Log("min", min)
	t.Log("max", max)
	t.Log("avg", avg/time.Duration(len(list)))
}
