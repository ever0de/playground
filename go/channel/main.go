package main

import "sync"

type Peer struct {
	channel chan []byte
}

func (p *Peer) write(data []byte) int {
	p.channel <- data
	return len(data)
}

func (p *Peer) read() []byte {
	return <-p.channel
}

func main() {
	channel := make(chan []byte)
	peer := Peer{channel: channel}

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go func() {
		peer.write([]byte("hello"))
		wg.Done()
	}()

	go func() {
		peer.write([]byte("world"))
		wg.Done()
	}()

	go func() {
		for i := 0; i < 2; i++ {
			println(string(peer.read()))
		}
		wg.Done()
	}()

	wg.Wait()
}
