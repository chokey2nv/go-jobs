package types

import (
	"context"
)

type Status string

const (
	Pending   Status = "PENDING"
	Running   Status = "RUNNING"
	Succeeded Status = "SUCCEEDED"
	Failed    Status = "FAILED"
	Cancelled Status = "CANCELLED"
)

// type Job struct {
// 	ID        string
// 	Type      string
// 	Status    Status
// 	Progress  int
// 	Message   string
// 	Result    any
// 	Error     string
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

type Job struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type,omitempty"`
	Status    Status `json:"status,omitempty"`
	Progress  int    `json:"progress,omitempty"`
	Message   string `json:"message,omitempty"`
	Result    any    `json:"result,omitempty"`
	Error     string `json:"error,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

type JobExecutor func(
	ctx context.Context,
	report func(message string, progress int),
) (any, error)
