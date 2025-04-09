package pool

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Task func()

type Pool struct {
	capacity int32
	running  int32
	workers  chan Task
	wg       sync.WaitGroup
	closed   int32
}

func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, errors.New("不能小于0")
	}
	p := &Pool{
		capacity: int32(size),
		workers:  make(chan Task, size),
	}
	for i := 0; i < size; i++ {
		p.wg.Add(1)
		go p.worker()
	}
	return p, nil
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.workers {
		if task == nil {
			return
		}
		atomic.AddInt32(&p.running, 1)
		task()
		atomic.AddInt32(&p.running, -1)
	}
}

func (p *Pool) Submit(task Task) error {
	if atomic.LoadInt32(&p.closed) == 1 {
		return errors.New("池已关闭")
	}
	p.workers <- task
	return nil
}

func (p *Pool) Running() int {
	return int(atomic.LoadInt32((&p.running)))
}

func (p *Pool) Cap() int {
	return int(p.capacity)
}

func (p *Pool) Close() {
	if atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		close(p.workers)
		p.wg.Wait()
	}
}
