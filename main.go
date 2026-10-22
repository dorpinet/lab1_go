package main

import (
	"fmt"
	"time"
)

const requestTimeout = 10 * time.Second

var targets = []string{
	"https://go.dev",
	"https://github.com",
	"https://example.com",
	"https://httpbin.org/status/404",
	"https://invalid-fake-host-test.org",
}

func main() {
	var checker Checker = NewHTTPChecker(requestTimeout)

	fmt.Println("== url health checker ==")
	runSequential(checker, targets)
}
func runSequential(c Checker, urls []string) {
	for _, u := range urls {
		fmt.Println(c.Check(u))
	}
}
