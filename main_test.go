package main

import (
	gg "gogasm/src"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
)

func BenchmarkRunner(b *testing.B) {
	r := gg.NewRunner("tcp", "127.0.0.1:3000", "server/big.txt")

	data, _ := ioutil.ReadFile("server/big.txt")
	lines := strings.Split(string(data), "\n")

	var wg sync.WaitGroup
	for n := 0; n <= b.N; n++ {
		go r.DoRequest(lines[n], &wg)
	}
	wg.Wait()
}
