package creator

import (
	"context"

	"code/benjamin/async_task/biz/common"
)

type CreatorService interface {
	CreateTask(ctx context.Context, taskType common.TaskType, bizRequest *interface{}, option *common.TaskAdditionalOption) (taskID string, err error)
}
