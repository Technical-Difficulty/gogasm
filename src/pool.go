package gogasm

type Generator func() interface{}

type Pool struct {
	queue     chan interface{}
	generator Generator
}

func (p *Pool) PreGenerate() {
	for cap(p.queue) > 0 {
		select {
		case p.queue <- p.generator():
			continue
		default:
			return
		}
	}
}

func (p *Pool) Get() interface{} {
	return <-p.queue
}

func (p *Pool) Put(worker interface{}) {
	p.queue <- worker
}
