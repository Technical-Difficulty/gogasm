package gogasm

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var HeadSlice = []byte("HEAD /")
var HeadSliceLength = len(HeadSlice)

type Runner struct {
	network        string
	address        string
	wordlistPath   string
	slicePool      *WorkerPool
	connectionPool *ConnectionPool
	requestSlice   RunnerHTTPSlice
}

type RunnerHTTPSlice struct {
	bytes  []byte
	length int
}

func NewRunner(network string, address string, wordlistPath string) Runner {
	requestAddress, _ := net.ResolveTCPAddr(network, address)
	requestSlice := []byte(" HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n")

	runner := Runner{
		network:      network,
		address:      address,
		wordlistPath: wordlistPath,
		requestSlice: RunnerHTTPSlice{bytes: requestSlice, length: len(requestSlice)},
		connectionPool: NewConnectionPool(300, func() interface{} {
			conn, _ := net.DialTCP(network, nil, requestAddress)
			conn.SetLinger(0)

			return conn
		}),
	}

	runner.slicePool = NewWorkerPool(50, func() interface{} {
		worker := Worker{
			tmp:            make([]byte, 12),
			requestSlices:  runner.PreallocateRequestByteSlices(),
			connectionPool: runner.connectionPool,
		}
		return worker
	})

	return runner
}

func (r Runner) Start() {
	data, _ := ioutil.ReadFile(r.wordlistPath)
	lines := strings.Split(string(data), "\n")

	var wg sync.WaitGroup

	startTime := time.Now().UnixNano()
	for i := range lines {
		if line := lines[i]; line != "" {
			go r.DoRequest(line, &wg)
		}
	}

	wg.Wait()
	endTime := time.Now().UnixNano()

	difference := endTime - startTime
	differenceSeconds := float64(difference) / 1000000000.0

	fmt.Println(fmt.Sprintf("%s\nTook %f seconds to process %d line wordlist.", TERMINAL_CLEAR_LINE, differenceSeconds, len(lines)))
	fmt.Println(fmt.Sprintf("Average of %f requests per second\n", float64(len(lines))/differenceSeconds))
}

func (r Runner) DoRequest(path string, wg *sync.WaitGroup) {
	wg.Add(1)
	worker := r.slicePool.Get().(Worker)
	worker.SocketRequest(path)
	wg.Done()
	r.slicePool.Put(worker)
	go func() {
		fmt.Fprintf(os.Stderr, "%s Processing :: /%s", TERMINAL_CLEAR_LINE, path)
	}()
}

func (r Runner) PreallocateRequestByteSlices() [][]byte {
	var slices [][]byte

	for i := 0; i <= 28; i++ {
		slice := make([]byte, HeadSliceLength+r.requestSlice.length+i) // Create slice of length i + other slice lengths above
		copy(slice[0:HeadSliceLength], HeadSlice)                      // Copy in the headSlice at the start of our generated slice
		copy(slice[HeadSliceLength+i:], r.requestSlice.bytes)          // Copy the httpSlice at the end of our generated slice

		slices = append(slices, slice)
	}

	return slices
}
