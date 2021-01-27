package gogasm

type WorkerGenerator func() Worker

type WorkerPool struct {
	queue     chan Worker
	generator WorkerGenerator
}

func NewWorkerPool(size int, generator WorkerGenerator) *WorkerPool {
	if generator == nil {
		panic("Need generator function")
	}
	result := &WorkerPool{
		queue:     make(chan Worker, size),
		generator: generator,
	}
	result.PreGenerate()
	return result
}

func (wp *WorkerPool) PreGenerate() {
	for cap(wp.queue) > 0 {
		select {
		case wp.queue <- wp.generator():
			continue
		default:
			return
		}
	}
}

func (wp *WorkerPool) Get() Worker {
	return <-wp.queue
}

func (wp *WorkerPool) Put(worker Worker) {
	wp.queue <- worker
}
