package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	requestTimeout = 10 * time.Second
	separatorWidth = 42
)

// the squad of urls we're gonna poke. one of them is fake on purpose, just to see the error flow
var targets = []string{
	"https://go.dev",
	"https://github.com",
	"https://example.com",
	"https://httpbin.org/status/404",
	"https://invalid-fake-host-test.org",
}

func main() {
	var checker Checker = NewHTTPChecker(requestTimeout)

	fmt.Println("== url health checker==")

	fmt.Println("\n[sequential check]")
	_, sequentialTime := runSequential(checker, targets)

	fmt.Println("\n[concurrent check]")
	summary, concurrentTime := runConcurrent(checker, targets)

	printReport(summary, sequentialTime, concurrentTime)
}
func runSequential(c Checker, urls []string) (Summary, time.Duration) {
	start := time.Now()

	var summary Summary
	for _, u := range urls {
		res := c.Check(u)
		fmt.Println(res)
		summary.Record(res)
	}

	return summary, time.Since(start)
}
func runConcurrent(c Checker, urls []string) (Summary, time.Duration) {
	start := time.Now()
	var (
		wg    sync.WaitGroup
		stats SafeStats
	)
	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			res := c.Check(url)
			fmt.Println(res)
			stats.Record(res)
		}(u)
	}
	wg.Wait()
	return stats.Summary(), time.Since(start)
}
func printReport(s Summary, sequential, concurrent time.Duration) {
	separator := strings.Repeat("=", separatorWidth)
	fmt.Println(separator)
	fmt.Println("final monitoring report:")
	fmt.Printf("- total urls checked: %d\n", s.Total)
	fmt.Printf("- all good (200 ok): %d\n", s.Success)
	fmt.Printf("- other http statuses (404 and friends): %d\n", s.Other)
	fmt.Printf("- network errors (rip): %d\n", s.Errors)
	fmt.Printf("- total time: %v (vs %v in sequential mode)\n",
		concurrent.Round(time.Millisecond), sequential.Round(time.Millisecond))
	fmt.Printf("- speedup: x%.1f\n", float64(sequential)/float64(concurrent))
	fmt.Println(separator)
}
