package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type Runner struct {
	network string
	address string
}

func (r Runner) Process(filename string) {
	data, _ := ioutil.ReadFile(filename)
	lines := strings.Split(string(data), "\n")

	var wg sync.WaitGroup

	startTime := time.Now().UnixNano()
	for _, line := range lines {
		go r.NewRequest(line, &wg)
	}

	wg.Wait()
	endTime := time.Now().UnixNano()

	difference := endTime - startTime
	differenceSeconds := float64(difference) / 1000000000.0

	fmt.Println(fmt.Sprintf("Took %f seconds to process %d line wordlist.\n", differenceSeconds, len(lines)))
	fmt.Println(fmt.Sprintf("Average of %f requests per second", float64(len(lines))/differenceSeconds))
}

func (r Runner) NewRequest(path string, wg *sync.WaitGroup) {
	wg.Add(1)
	worker, _ := slicePool.Borrow().(*Worker)
	worker.SocketRequest(path)
	go slicePool.Return(worker)
	wg.Done()
}
