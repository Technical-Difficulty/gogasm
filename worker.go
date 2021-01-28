package main

import (
	"fmt"
)

type Worker struct {
	tmp           []byte
	requestSlices [][]byte
}

func (w Worker) CreateRequestByteSlice(line string) []byte {
	copy(w.requestSlices[len(line)][headSliceLength:], line)
	return w.requestSlices[len(line)]
}

func (w Worker) SocketRequest(path string) {
	requestSlice := w.CreateRequestByteSlice(path)

	cc := connectionPool.AcquireConn()
	conn := cc.conn

	_, err := conn.Write(requestSlice)
	if err != nil {
		connectionPool.ReleaseAndRedialConn(cc)
		return
	}

	_, err = conn.Read(w.tmp)
	if err != nil {
		connectionPool.ReleaseAndRedialConn(cc)
		return
	}

	go w.ReadStatusCode(path)
	go w.CloseAndRelease(cc)
}

func (w Worker) ReadStatusCode(path string) {
	if w.tmp[9] != 52 || w.tmp[11] != 52 { // 52 = 4 in decimal, we are checking for 404.
		fmt.Println(fmt.Sprintf("[%s] %s", string(w.tmp[9:12]), path))
	}
}

func (w Worker) CloseAndRelease(conn *clientConn) {
	connectionPool.ReleaseConn(conn)
}
