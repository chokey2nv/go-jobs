package job

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/chokey2nv/go-jobs/v1/stores"
	"github.com/chokey2nv/go-jobs/v1/types"
	"github.com/google/uuid"
)

// NewJob creates a new job instance with default Pending status.
func NewJob(name string) *types.Job {
	now := time.Now()
	return &types.Job{
		ID:        uuid.NewString(),
		Type:      name,
		Status:    types.Pending,
		Progress:  0,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// jobRuntime holds the cancel function for a running job.
type jobRuntime struct {
	cancel context.CancelFunc
}

// JobService manages job execution, tracking, and storage.
type JobService struct {
	store stores.Store

	mu       sync.Mutex
	runtimes map[string]*jobRuntime

	pool types.WorkerPool
}

// New creates a JobService with a default pool.
func New(store stores.Store) *JobService {
	p := NewWorkerPool(10, 100)
	p.Start()
	return &JobService{
		store:    store,
		pool:     p,
		runtimes: make(map[string]*jobRuntime),
	}
}

// StartAsync schedules a job to run asynchronously in the worker pool.
func (s *JobService) StartAsync(
	ctx context.Context,
	name string,
	exec types.Executor,
) (*types.Job, error) {

	job := NewJob(name)
	if err := s.store.Create(ctx, job); err != nil {
		return nil, err
	}

	jobCtx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.runtimes[job.ID] = &jobRuntime{cancel: cancel}
	s.mu.Unlock()

	err := s.pool.Submit(func() {
		s.run(job, jobCtx, exec)
	})
	if err != nil {
		// cleanup if pool queue is full
		s.mu.Lock()
		delete(s.runtimes, job.ID)
		s.mu.Unlock()
		return nil, err
	}

	return job, nil
}

// Start runs a job immediately in a separate goroutine (fire-and-forget).
func (s *JobService) Start(
	ctx context.Context,
	name string,
	exec types.Executor,
) (*types.Job, error) {

	job := NewJob(name)
	if err := s.store.Create(ctx, job); err != nil {
		return nil, err
	}

	jobCtx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	s.runtimes[job.ID] = &jobRuntime{cancel: cancel}
	s.mu.Unlock()

	go s.run(job, jobCtx, exec)

	return job, nil
}

// run executes the job and updates progress/status to the store.
func (s *JobService) run(
	job *types.Job,
	ctx context.Context,
	exec types.Executor,
) {
	defer func() {
		s.mu.Lock()
		delete(s.runtimes, job.ID)
		s.mu.Unlock()
	}()

	job.Status = types.Running
	_ = s.store.Update(context.Background(), job)

	report := NewProgressReporter(s.store, job)

	result, err := exec(ctx, report)

	switch {
	case errors.Is(ctx.Err(), context.Canceled):
		job.Status = types.Cancelled
		job.Message = "Job cancelled"
	case err != nil:
		job.Status = types.Failed
		job.Error = err.Error()
	default:
		job.Status = types.Succeeded
		job.Result = result
		job.Progress = 100
	}

	job.UpdatedAt = time.Now()
	_ = s.store.Update(context.Background(), job)
}

// Cancel stops a running job by ID.
func (s *JobService) Cancel(ctx context.Context, jobID string) error {
	s.mu.Lock()
	runtime, ok := s.runtimes[jobID]
	s.mu.Unlock()

	if !ok {
		return errors.New("job not running or already finished")
	}

	runtime.cancel()
	return nil
}

// Get retrieves a job by ID.
func (s *JobService) Get(ctx context.Context, id string) (*types.Job, error) {
	return s.store.Get(ctx, id)
}

// Remove a job by ID
func (s *JobService) Remove(ctx context.Context, id string) (string, error) {
	return s.store.Remove(ctx, id)
}

// List retrieves jobs by filter.
func (s *JobService) List(ctx context.Context, f stores.Filter) ([]*types.Job, error) {
	return s.store.List(ctx, f)
}
