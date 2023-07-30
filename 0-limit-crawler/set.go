package main

import "sync"

type Set struct {
	m map[string]bool
	sync.Mutex
}

func New() *Set {
	return &Set{
		m: make(map[string]bool),
	}
}

func (s *Set) Add(item string) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) Has(item string) bool {
	s.Lock()
	defer s.Unlock()
	_, ok := s.m[item]
	return ok
}
