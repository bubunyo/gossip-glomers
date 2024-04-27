package main

import "sync"

type store[T any] struct {
	m sync.RWMutex
	d []T
}

func newStore[T any]() *store[T] {
	s := &store[T]{
		d: make([]T, 0),
	}
	return s
}

func (s *store[T]) save(m T) {
	s.m.Lock()
	s.d = append(s.d, m)
	s.m.Unlock()
}

func (s *store[T]) read() []T {
	s.m.RLock()
	l := s.d
	s.m.RUnlock()
	return l
}
