package main

import (
	"net/http"
	"sync"
)

type Summary struct {
	Total   int
	Success int
	Other   int
	Errors  int
}

func (s *Summary) Record(r Result) {
	s.Total++
	switch {
	case r.Err != nil:
		s.Errors++
	case r.Status == http.StatusOK:
		s.Success++
	default:
		s.Other++
	}
}

type SafeStats struct {
	mu      sync.Mutex
	summary Summary
}

func (s *SafeStats) Record(r Result) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.summary.Record(r)
}
func (s *SafeStats) Summary() Summary {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.summary
}
