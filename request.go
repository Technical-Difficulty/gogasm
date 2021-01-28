package main

import (
	"net"
	"time"
	"sync"
	"fmt"
)

type debug struct {
	sync.Mutex
	pool int
	queued int
	delivered int
	requests int
	failed int
	notdelivered int
	closed int
}

// Just here while testing
var d debug = debug { 
	queued: 0, 
	delivered: 0, 
	requests: 0, 
	failed: 0, 
	notdelivered: 0, 
	closed: 0,
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

type clientConn struct {
	conn net.Conn
}

type connPool struct {
	sync.Mutex
	size int
	dialer net.Dialer
	addr string
	conns []*clientConn
	wp waitPool
}

type wantConn struct {
	ready chan *clientConn
}

type waitPool struct {
	sync.Mutex
	waiting []*wantConn
}

func (wp *waitPool) queueIdle() *wantConn {
	d.Lock()
	d.queued++
	d.Unlock()

	wp.Lock()
	defer wp.Unlock()

	wc := &wantConn {
		ready: make(chan *clientConn, 1),
	}

	wp.waiting = append(wp.waiting, wc)

	return wc
}

func (wp *waitPool) TryDeliverConn(cc *clientConn) bool {
	for wp.len() > 0 {
		wc := wp.Shift()
		select {
		case wc.ready <- cc:
			d.Lock()
			d.delivered++
			d.Unlock()
			return true
		default:
			d.Lock()
			d.notdelivered++
			d.Unlock()
			return false
		}
	}
	return false
}

func (wp *waitPool) len() int {
	return len(wp.waiting)
}

func (wp *waitPool) Shift() *wantConn {
	wp.Lock()
	defer wp.Unlock()
	wc := wp.waiting[0]
	wp.waiting = wp.waiting[1:]
	return wc
}

func (cp *connPool) Prepare() {
	cp.Lock()
	defer cp.Unlock()

	for i:=0; i<cp.size; i++ {
		conn, err := cp.dialer.Dial("tcp", cp.addr)
		if err != nil {
			panic(err)
		}

		cc := clientConn { conn: conn }
		cp.conns = append(cp.conns, &cc)
	}
}

func (cp *connPool) Shift() *clientConn {
	cp.Lock()
	defer cp.Unlock()
	conn := cp.conns[0]
	cp.conns = cp.conns[1:]
	return conn
}

func (cp *connPool) Push(cc *clientConn) {
	cp.Lock()
	defer cp.Unlock()
	cp.conns = append(cp.conns, cc)
}

func (cp *connPool) AcquireConn() *clientConn {
	if len(cp.conns) <= 0 {
		wc := cp.wp.queueIdle()
		select {
		case cc := <- wc.ready:
			return cc
		// TODO case <- timer.done:
		}
	}

	clientConn := cp.Shift()
	return clientConn
}

func (cp *connPool) ReleaseConn(cc *clientConn) {
	if len(cp.conns) >= cp.size {
		d.Lock()
		d.closed++
		d.Unlock()
		cc.conn.Close()
	}

	delivered := false
	if cp.wp.len() > 0 {
		delivered = cp.wp.TryDeliverConn(cc)
	}

	if !delivered {
		cp.Push(cc)
	}
}

func (cp *connPool) ReleaseAndRedialConn(cc *clientConn) {
	fmt.Println("Redial")
	d.Lock()
	d.failed++
	d.Unlock()
}

func NewConnPool(size int) connPool {
	cp := connPool { 
		size: size, 
		addr: "localhost:3000",
		conns: []*clientConn{},
		dialer: net.Dialer{Timeout: 10 * time.Second},
	}
	return cp
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

	for i:=0; i<1000; i++ {
		wg.Add(1)
		go DoRequest(&cp, &wg)
	}
	wg.Wait()
	d.Report()
}
