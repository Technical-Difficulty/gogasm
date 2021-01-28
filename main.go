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

func DoRequest(cp *connPool, wg *sync.WaitGroup) {
	cc := cp.AcquireConn()
	conn := cc.conn

	d.Lock()
	d.requests++
	d.Unlock()

	buff := make([]byte, 12)
	_, write_err := conn.Write([]byte("HEAD / HTTP/1.1\r\nAccept-Encoding: gzip\r\n\r\n"))
	if write_err != nil {
		cp.ReleaseAndRedialConn(cc)
		return
	}

	_, read_err := conn.Read(buff)
	if read_err != nil {
		cp.ReleaseAndRedialConn(cc)
		return
	}

	//fmt.Println(fmt.Sprintf("[%s]", string(buff[9:12])))

	wg.Done()
	cp.ReleaseConn(cc)
}

func main() {
	var wg sync.WaitGroup

	cp := NewConnPool(100)
	cp.Prepare()

	d.pool = 100

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go DoRequest(&cp, &wg)
	}
	wg.Wait()
	d.Report()
}
