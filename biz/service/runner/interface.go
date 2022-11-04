package runner

import (
	"code/benjamin/async_task/biz/common"
)

type RunnerService interface {
	RegisterHandler(handlerMapping map[common.TaskType]*common.AsyncTaskHandler) error
}
