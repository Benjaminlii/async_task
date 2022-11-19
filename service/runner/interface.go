package runner

import (
	"github.com/Benjaminlii/async_task/common"
)

type RunnerService interface {
	RegisterHandler(handlerMapping map[common.TaskType]*common.AsyncTaskHandler) error
}
