package types

import (
	"context"
	"time"
)

type Status string

const (
	Pending   Status = "PENDING"
	Running   Status = "RUNNING"
	Succeeded Status = "SUCCEEDED"
	Failed    Status = "FAILED"
	Cancelled Status = "CANCELLED"
)

type Job struct {
	ID        string
	Type      string
	Status    Status
	Progress  int
	Message   string
	Result    any
	Error     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type JobExecutor func(
	ctx context.Context,
	report func(message string, progress int),
) (any, error)
