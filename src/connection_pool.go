package gogasm

import "net"

type ConnectionPool struct {
	Pool
}

func NewConnectionPool(size int, generator Generator) *ConnectionPool {
	if generator == nil {
		panic("Need generator function")
	}

	result := &ConnectionPool{
		Pool {
			queue:     make(chan interface{}, size),
			generator: generator,
		},
	}
	result.PreGenerate()
	return result
}

func (cp *ConnectionPool) CloseAndRedial(conn *net.TCPConn) {
	conn.Close()
	conn = cp.generator().(*net.TCPConn)
	cp.queue <- conn
}
