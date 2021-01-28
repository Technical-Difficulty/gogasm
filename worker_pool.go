package main

type MakerFunc func() (interface{}, error)

type Pool struct {
	queue chan interface{}
	maker MakerFunc
}

// Create a new fixed pool. size is the max number of object to pool
// maker is the function that generates new object for the pool (when pool is empty)
func NewFixedPool(size int, maker MakerFunc) *Pool {
	if maker == nil {
		panic("Need maker function")
	}
	result := &Pool{
		queue: make(chan interface{}, size),
		maker: maker,
	}
	result.PreFill()
	return result
}

// Prepopulate the pool with full elements. This will call the maker repeately until it is full
// Failed maker will be discarded. If maker never return successful result, this may be in dead loop
func (v *Pool) PreFill() {
	for cap(v.queue) > 0 {
		elem, err := v.maker()
		if err != nil {
			continue
		}
		select {
		case v.queue <- elem:
			continue
		default:
			return
		}
	}
}

// Borrow a object from the pool, block until one is available.
// If an object failed test upon checkout because of tester func fails, a new object will be made and returned
// Maker will be tried 3 times, with 1 seconds delay in between
func (v *Pool) Borrow() interface{} {
	return <-v.queue
}

// Return an object to the pool, the object doesn't has to be borrowed
// Returns true if returned successfully
// Returns false if pool is full and object had been discarded
// (which is unlikely unless you returned something extra to the pool)
func (v *Pool) Return(c interface{}) bool {
	select {
	case v.queue <- c:
		return true
	default:
		return false
	}
}
