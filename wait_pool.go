package main

import "sync"

type waitPool struct {
	sync.Mutex
	waiting []*wantConn
}

type wantConn struct {
	ready chan *clientConn
}

func (wp *waitPool) queueIdle() *wantConn {
	d.Lock()
	d.queued++
	d.Unlock()

	wp.Lock()
	defer wp.Unlock()

	wc := &wantConn{
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
