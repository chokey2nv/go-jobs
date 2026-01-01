package types

import (
	"context"
)

type ProgressReporter interface {
	Progress(p int, msg string)
}

type Executor func(
	ctx context.Context,
	report ProgressReporter,
) (result any, err error)

// type Executor func(
// 	ctx context.Context,
// 	report func(message string, progress int),
// ) (any, error)