package job

import (
	"context"
	"time"

	"github.com/chokey2nv/go-jobs/v1/stores"
	"github.com/chokey2nv/go-jobs/v1/types"
)

type progressReporter struct {
	store stores.Store
	job   *types.Job
}

func NewProgressReporter(store stores.Store, job *types.Job) types.ProgressReporter {
	return &progressReporter{
		store: store,
		job:   job,
	}
}
func (r *progressReporter) Progress(p int, msg string) {
	if p < 0 {
		p = 0
	}
	if p > 100 {
		p = 100
	}
	r.job.Progress = p
	r.job.Message = msg
	r.job.UpdatedAt = time.Now()
	_ = r.store.Update(context.Background(), r.job)
}
