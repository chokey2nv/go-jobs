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
	ID        string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" gql:"required"`
	Type      string    `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty" gql:"required"`
	Status    Status    `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty" gql:"required"`
	Progress  int       `protobuf:"bytes,1,opt,name=progress,proto3" json:"progress,omitempty" gql:"required"`
	Message   string    `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty" gql:"required"`
	Result    any       `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty" gql:"required"`
	Error     string    `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty" gql:"required"`
	CreatedAt time.Time `protobuf:"bytes,1,opt,name=createdAt,proto3" json:"createdAt,omitempty" gql:"required"`
	UpdatedAt time.Time `protobuf:"bytes,1,opt,name=updatedAt,proto3" json:"updatedAt,omitempty" gql:"required"`
}

type JobExecutor func(
	ctx context.Context,
	report func(message string, progress int),
) (any, error)
