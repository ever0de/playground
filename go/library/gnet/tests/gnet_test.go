package gnet_test

import (
	"encoding/binary"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gnet "github.com/ever0de/playground/library/gnet/tests"
)

func TestGnetWithPebble(t *testing.T) {
	server := gnet.NewServer("localhost:8080")

	go server.Start()

	client := gnet.NewClient("localhost:8080")

	now := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(150)

	go func() {
		cn := &atomic.Uint32{}
		cn.Store(0)
		for {
			time.Sleep(50 * time.Millisecond)

			go func() {
				now := time.Now()
				fmt.Println("client: send")
				val := make([]byte, 10*1024*1024)

				key := make([]byte, 4)
				binary.BigEndian.PutUint32(key, cn.Add(1))
				if err := client.Send(key, val); err != nil {
					panic(err)
				}
				wg.Done()
				fmt.Printf("client: send time: %v\n", time.Since(now))
			}()
		}
	}()

	wg.Wait()
	fmt.Printf("time: %v\n", time.Since(now))
}
