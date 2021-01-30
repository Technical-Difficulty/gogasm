package gogasm

import (
	"fmt"
	"net"
)

type Worker struct {
	tmp            []byte
	requestSlices  [][]byte
	connectionPool *ConnectionPool
}

func (w *Worker) CreateRequestByteSlice(line string, length int) []byte {
	copy(w.requestSlices[length][6:], line) // 6 = head slice length
	return w.requestSlices[length]
}

func (w *Worker) SocketRequest(path string) {
	requestSlice := w.CreateRequestByteSlice(path, len(path))

	conn := w.connectionPool.Get().(*net.TCPConn)
	conn.Write(requestSlice)
	conn.Read(w.tmp)

	go w.connectionPool.CloseAndRedial(conn)

	w.ReadStatusCode(path)
}

func (w *Worker) ReadStatusCode(path string) {
	if w.tmp[9] != 52 || (w.tmp[11] != 52 && w.tmp[11] != 48) { // 52 = 4 in decimal, we are checking for 404.
		fmt.Println(fmt.Sprintf("%s [%s] /%s", TERMINAL_CLEAR_LINE, string(w.tmp[9:12]), path))
	}
}
