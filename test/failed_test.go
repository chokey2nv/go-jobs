package test

import (
	"context"
	"testing"
	"time"

	"github.com/chokey2nv/go-jobs/job"
	"github.com/chokey2nv/go-jobs/job/stores"
	"github.com/chokey2nv/go-jobs/job/types"
	"github.com/stretchr/testify/require"
)

// waitForJob waits for the job to finish (Succeeded or Failed) or times out.
func waitForJob(ctx context.Context, svc *job.JobService, id string) (*types.Job, error) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			j, err := svc.Get(context.Background(), id)
			if err != nil {
				return nil, err
			}
			if j != nil && (j.Status == types.Succeeded || j.Status == types.Failed) {
				return j, nil
			}
		}
	}
}

func TestWorkerPool_Success(t *testing.T) {
	store := stores.NewMemoryStore()
	jobService := job.New(store)

	// Start a short async job
	job, err := jobService.StartAsync(context.Background(), "job-1",
		func(ctx context.Context, report types.ProgressReporter) (any, error) {
			report.Progress(50, "Halfway")
			time.Sleep(5 * time.Millisecond) // simulate work
			report.Progress(100, "Completed")
			return "ok", nil
		},
	)

	require.NoError(t, err)
	require.NotNil(t, job)

	// Wait for the job to complete with a 2-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	j, err := waitForJob(ctx, jobService, job.ID)
	require.NoError(t, err)

	require.Equal(t, types.Succeeded, j.Status)
	require.Equal(t, "ok", j.Result)
	require.Equal(t, 100, j.Progress)
	require.Equal(t, "Completed", j.Message) // If you set message in ProgressReporter
}

func TestWorkerPool_Failure(t *testing.T) {
	store := stores.NewMemoryStore()
	jobService := job.New(store)

	job, err := jobService.StartAsync(context.Background(), "job-fail",
		func(ctx context.Context, report types.ProgressReporter) (any, error) {
			report.Progress(10, "Starting")
			return nil, context.DeadlineExceeded
		},
	)
	require.NoError(t, err)
	require.NotNil(t, job)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	j, err := waitForJob(ctx, jobService, job.ID)
	require.NoError(t, err)

	require.Equal(t, types.Failed, j.Status)
	require.NotEmpty(t, j.Error)
}
