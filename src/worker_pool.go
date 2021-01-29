package gogasm

type WorkerPool struct {
	Pool
}

func NewWorkerPool(size int, generator Generator) *WorkerPool {
	if generator == nil {
		panic("Need generator function")
	}
	result := &WorkerPool{
		Pool {
			queue:     make(chan interface{}, size),
			generator: generator,
		},
	}
	result.PreGenerate()
	return result
}