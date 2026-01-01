package stores

import (
	"context"

	"github.com/chokey2nv/go-jobs/v1/types"
)

type Filter struct {
	Type   string
	Status types.Status
	Limit  int
}

type Store interface {
	Create(ctx context.Context, job *types.Job) error
	Update(ctx context.Context, job *types.Job) error
	Get(ctx context.Context, id string) (*types.Job, error)
	List(ctx context.Context, f Filter) ([]*types.Job, error)
}
