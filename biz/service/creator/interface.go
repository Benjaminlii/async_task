package creator

import (
	"context"

	"github.com/Benjaminlii/async_task/biz/common"
)

type CreatorService interface {
	CreateTask(ctx context.Context, taskType common.TaskType, bizRequest *interface{}, option *common.TaskAdditionalOption) (taskID string, err error)
}
