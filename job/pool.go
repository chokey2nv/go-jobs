package job

import (
	"context"
	"errors"
	"sync"

	types "github.com/chokey2nv/go-jobs/job/types"
)

// WorkerPool implements a fixed-size pool of goroutines
// executing submitted tasks.
type pool struct {
	workers int
	queue   chan func()
	wg      sync.WaitGroup

	stopOnce sync.Once
	quit     chan struct{}
}

// NewWorkerPool creates a new worker pool.
func NewWorkerPool(workers, queueSize int) types.WorkerPool {
	return &pool{
		workers: workers,
		queue:   make(chan func(), queueSize),
		quit:    make(chan struct{}),
	}
}

// Start launches the worker goroutines.
func (p *pool) Start() {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

// worker runs in a goroutine, processing tasks from the queue.
func (p *pool) worker() {
	defer p.wg.Done()
	for {
		select {
		case task := <-p.queue:
			task()
		case <-p.quit:
			return
		}
	}
}

// Submit adds a task to the pool queue.
func (p *pool) Submit(task func()) error {
	select {
	case p.queue <- task:
		return nil
	default:
		return errors.New("worker pool queue is full")
	}
}

// Stop waits for all workers to finish or returns on ctx cancel.
func (p *pool) Stop(ctx context.Context) error {
	p.stopOnce.Do(func() {
		close(p.quit)
	})

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
