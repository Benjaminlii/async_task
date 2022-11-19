package common

import "context"

type AsyncTaskHandler interface {
	HandleMessage(ctx context.Context, taskRequest string) (taskResponse interface{}, err error)
}
