package go_jobpool

import (
	"sync"
	"sync/atomic"
	"runtime"
)

type Worker struct {
	//the pool current worker owns to
	pool  *Pool
	// personal job queue of current worker
	jQueue  chan Job
}

func NewWorker(parent *Pool) *Worker {
	return &Worker{
		pool: 	parent,
		jQueue: make(chan Job),
	}
}

func (w *Worker)run() {
	go func() {
		select {
			case job := <- w.jQueue:
				if job == nil {
					atomic.AddInt32(&w.pool.running, -1)
					return
				}
				job.Before()
				job.Run()
				job.After()
				w.pool.putWorker(w)
			case <- w.pool.releaseSignal:
				return
			}
	}()
}

func (w *Worker)reciveJob(job Job) {
	w.jQueue <- job
}

func (w *Worker)close() {
	w.reciveJob(nil)
}



type Pool struct {
	capacity 		int32
	running			int32
	workers 		[]*Worker
	freeSignal		chan struct{}
	releaseSignal	chan struct{}
	syncLock 		sync.RWMutex
}

func NewDefaultPool() *Pool {
	return &Pool{
		capacity:		int32(runtime.NumCPU())+1,
		workers:		make([]*Worker, 0, int32(runtime.NumCPU())+1),
		freeSignal: 	make(chan struct{}),
		releaseSignal: 	make(chan struct{}),
		syncLock: 		sync.RWMutex{},
	}
}


func NewPool(size int32) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		capacity: 		size,
		workers: 		make([]*Worker, 0, size),
		freeSignal: 	make(chan struct{}),
		releaseSignal: 	make(chan struct{}),
		syncLock: 		sync.RWMutex{},
	}
}



func (p *Pool)Submit(job Job) error {
	if _, opened := <- p.releaseSignal;!opened {
		return ErrPoolClosed{}
	}

	worker := p.getWorker()
	worker.reciveJob(job)
	return nil
}


func (p *Pool)getWorker() *Worker {
	var worker *Worker
	wait := false

	p.syncLock.Lock()
	n := len(p.workers) - 1
	if n < 0 {
		if p.running >= p.capacity {
			wait = true
		}else{
			p.running ++
		}
	}else{
		worker = p.workers[n]
		p.workers = p.workers[:n]
	}
	p.syncLock.Unlock()

	if wait {
		// block until free worker
		<- p.freeSignal
		p.syncLock.Lock()
		n := len(p.workers) - 1
		worker = p.workers[n]
		p.workers = p.workers[:n]
		p.syncLock.Unlock()
	}else if worker == nil {
		worker = NewWorker(p)
		worker.run()
	}
	return worker
}


func (p *Pool) putWorker(worker *Worker) {
	p.syncLock.Lock()
	p.workers = append(p.workers, worker)
	p.syncLock.Unlock()
	p.freeSignal <- struct{}{}
}


func (p *Pool) Release() {
	close(p.releaseSignal)
}




