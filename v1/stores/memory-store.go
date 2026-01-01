package stores

import (
	"context"
	"errors"
	"sync"

	"github.com/chokey2nv/go-jobs/v1/types"
)

type MemoryStore struct {
	mu   sync.RWMutex
	jobs map[string]*types.Job
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		jobs: make(map[string]*types.Job),
	}
}

func (m *MemoryStore) Create(ctx context.Context, job *types.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs[job.ID] = clone(job)
	return nil
}

func (m *MemoryStore) Update(ctx context.Context, job *types.Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.jobs[job.ID] = clone(job)
	return nil
}

func (m *MemoryStore) Get(ctx context.Context, id string) (*types.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	j, ok := m.jobs[id]
	if !ok {
		return nil, errors.New("job not found")
	}
	return clone(j), nil
}
func (m *MemoryStore) Remove(ctx context.Context, id string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	delete(m.jobs, id)
	return id, nil
}

func (m *MemoryStore) List(ctx context.Context, f Filter) ([]*types.Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var res []*types.Job
	for _, j := range m.jobs {
		if f.Type != "" && j.Type != f.Type {
			continue
		}
		if f.Status != "" && j.Status != f.Status {
			continue
		}
		res = append(res, clone(j))
		if f.Limit > 0 && len(res) >= f.Limit {
			break
		}
	}
	return res, nil
}

func clone(j *types.Job) *types.Job {
	cp := *j
	return &cp
}
