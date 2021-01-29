package main

import "net"

type ConnectionGenerator func() *net.TCPConn

type ConnectionPool struct {
	queue     chan *net.TCPConn
	generator ConnectionGenerator
}

func NewConnectionPool(size int, generator ConnectionGenerator) *ConnectionPool {
	if generator == nil {
		panic("Need generator function")
	}

	result := &ConnectionPool{
		queue:     make(chan *net.TCPConn, size),
		generator: generator,
	}
	result.PreGenerate()
	return result
}

func (cp *ConnectionPool) PreGenerate() {
	for cap(cp.queue) > 0 {
		select {
		case cp.queue <- cp.generator():
			continue
		default:
			return
		}
	}
}

func (cp *ConnectionPool) Get() *net.TCPConn {
	return <-cp.queue
}

func (cp *ConnectionPool) CloseAndRedial(conn *net.TCPConn) {
	conn.Close()
	conn = cp.generator()
	cp.queue <- conn
}
