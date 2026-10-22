package main

import (
	"net/http"
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
