package types

import (
	"context"
)

// type WorkerPool interface {
// 	Submit(job *Job, fn func()) error
// 	Start()
// 	Stop(ctx context.Context) error
// }

type WorkerPool interface {
	Start()
	Submit(task func()) error
	Stop(ctx context.Context) error
}
