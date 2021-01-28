package main

import (
	"fmt"
	"sync"
)

type debug struct {
	sync.Mutex
	pool         int
	queued       int
	delivered    int
	requests     int
	failed       int
	notdelivered int
	closed       int
}

// Just here while testing
var d = debug{
	queued:       0,
	delivered:    0,
	requests:     0,
	failed:       0,
	notdelivered: 0,
	closed:       0,
}

func (d *debug) Report() {
	fmt.Println("pool", d.pool)
	fmt.Println("queued", d.queued)
	fmt.Println("delivered", d.delivered)
	fmt.Println("requests", d.requests)
	fmt.Println("failed", d.failed)
	fmt.Println("notdelivered", d.notdelivered)
	fmt.Println("closed", d.closed)
}

//func DoRequest(cp *connPool, wg *sync.WaitGroup) {
//	cc := cp.AcquireConn()
//	conn := cc.conn
//
//	//d.Lock()
//	//d.requests++
//	//d.Unlock()
//
//	buff := make([]byte, 12)
//	_, write_err := conn.Write([]byte("HEAD / HTTP/1.1\r\nAccept-Encoding: gzip\r\n\r\n"))
//	if write_err != nil {
//		cp.ReleaseAndRedialConn(cc)
//		return
//	}
//
//	_, read_err := conn.Read(buff)
//	if read_err != nil {
//		cp.ReleaseAndRedialConn(cc)
//		return
//	}
//
//	//fmt.Println(fmt.Sprintf("[%s]", string(buff[9:12])))
//
//	wg.Done()
//	cp.ReleaseConn(cc)
//}

//func main() {
//	var wg sync.WaitGroup
//
//	cp := NewConnPool(100)
//	cp.Prepare()
//
//	d.pool = 100
//
//	for i := 0; i < 10000; i++ {
//		wg.Add(1)
//		go DoRequest(&cp, &wg)
//	}
//	wg.Wait()
//	d.Report()
//}

var slicePool *Pool
var connectionPool connPool

func main() {
	address := "127.0.0.1:3000"

	slicePool = NewFixedPool(225, func() (interface{}, error) {
		worker := new(Worker)
		worker.tmp = make([]byte, 12)
		worker.requestSlices = PreallocateRequestByteSlices()
		return worker, nil
	})

	connectionPool = NewConnPool(100, address)
	connectionPool.Prepare()

	runner := Runner{
		network: "tcp",
		address: address,
	}

	runner.Process("server/big.txt")
}

var headSliceLength = 6

func PreallocateRequestByteSlices() [][]byte {
	headSlice := []byte{72, 69, 65, 68, 32, 47} // HEAD /
	httpSlice := []byte(" HTTP/1.1\r\nHost: localhost\r\nAccept-Encoding: gzip\r\n\r\n")
	//httpSlice := []byte{32, 72, 84, 84, 80, 47, 49, 46, 48, 13, 10, 13, 10} // HTTP/1.0\r\n\r\n
	//httpSlice := []byte{32, 72, 84, 84, 80, 47, 49, 46, 48, 13, 10, 65, 99, 99, 101, 112, 116, 45, 69, 110, 99, 111, 100, 105, 110, 103, 58, 32, 103, 122, 105, 112, 13, 10, 13, 10} // HEAD / HTTP/1.1\r\nAccept-Encoding: gzip\r\n\r\n
	httpSliceLength := len(httpSlice)

	var slices [][]byte

	for i := 0; i <= 28; i++ {
		slice := make([]byte, headSliceLength+httpSliceLength+i) // Create slice of length i + other slice lengths above
		copy(slice[0:headSliceLength], headSlice)                // Copy in the headSlice at the start of our generated slice
		copy(slice[headSliceLength+i:], httpSlice)               // Copy the httpSlice at the end of our generated slice

		slices = append(slices, slice)
	}

	return slices
}
