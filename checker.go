package main

import (
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	URL      string
	Status   int
	Duration time.Duration
	Err      error
}

func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("-> %s: error (%v)", r.URL, r.Err)
	}
	return fmt.Sprintf("-> %s: Status %d %s (time: %v)",
		r.URL, r.Status, http.StatusText(r.Status), r.Duration.Round(time.Millisecond))
}

type Checker interface {
	Check(url string) Result
}
type HTTPChecker struct {
	client *http.Client
}

func NewHTTPChecker(timeout time.Duration) *HTTPChecker {
	return &HTTPChecker{client: &http.Client{Timeout: timeout}}
}
func (c *HTTPChecker) Check(url string) Result {
	start := time.Now()
	resp, err := c.client.Get(url)
	elapsed := time.Since(start)

	if err != nil {
		return Result{URL: url, Duration: elapsed, Err: err}
	}
	defer resp.Body.Close()

	return Result{URL: url, Status: resp.StatusCode, Duration: elapsed}
}
