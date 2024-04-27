package main

import "sync"

type store[T comparable] struct {
	m sync.RWMutex
	d map[T]struct{}
}

func newStore[T comparable]() *store[T] {
	s := &store[T]{
		d: map[T]struct{}{},
	}
	return s
}

func (s *store[T]) save(m ...T) {
	s.m.Lock()
	for _, v := range m {
		s.d[v] = struct{}{}
	}
	s.m.Unlock()
}

func (s *store[T]) read() []T {
	i := 0
	s.m.RLock()
	l := make([]T, len(s.d))
	for v := range s.d {
		l[i] = v
		i += 1
	}
	s.m.RUnlock()
	return l
}
