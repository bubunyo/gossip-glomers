package main

import "sync"

type broadcast struct {
	m sync.RWMutex
	// c chan int
	d []float64
}

func newBroadcast() *broadcast {
	b := &broadcast{
		// c: make(chan int, 100),
		d: make([]float64, 10),
	}
	return b
}

// func (b *broadcast) init() {
// 	go func() {
// 		for _, v := range b.c {

// 		}
// 	}()
// }

func (b *broadcast) save(m float64) {
	b.m.Lock()
	b.d = append(b.d, m)
	b.m.Unlock()
}

func (b *broadcast) read() []float64 {
	b.m.RLock()
	l := b.d
	b.m.RUnlock()
	return l
}

// func (b *broadcast) close() {
// 	close(b.c)
// }
