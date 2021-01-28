package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type connPool struct {
	sync.Mutex
	size   int
	dialer net.Dialer
	addr   string
	conns  []*clientConn
	wp     waitPool
}

type clientConn struct {
	conn net.Conn
}

func NewConnPool(size int, address string) connPool {
	cp := connPool{
		size:   size,
		addr:   address,
		conns:  []*clientConn{},
		dialer: net.Dialer{Timeout: 10 * time.Second},
	}
	return cp
}

func (cp *connPool) Prepare() {
	cp.Lock()
	defer cp.Unlock()

	for i := 0; i < cp.size; i++ {
		conn, err := cp.dialer.Dial("tcp", cp.addr)
		if err != nil {
			panic(err)
		}

		cc := clientConn{conn: conn}
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
		case cc := <-wc.ready:
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
