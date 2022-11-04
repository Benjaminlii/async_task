package async_task

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Benjaminlii/async_task/biz/common"
	"github.com/Benjaminlii/async_task/biz/config"
	"github.com/Benjaminlii/async_task/biz/service/creator"
	"github.com/Benjaminlii/async_task/biz/service/syncer"
)

type TaskCenterClient struct {
	Options *config.Options
}

func (c *TaskCenterClient) Create(ctx context.Context, taskType common.TaskType, bizRequest *interface{}, option *common.TaskAdditionalOption) (string, error) {
	if !CheckInit() {
		return "", errors.New("[Create] client need init")
	}
	taskID, err := creator.NewCreatorService().CreateTask(ctx, taskType, bizRequest, option)
	if err != nil {
		return "", errors.New("[Create] CreateTask error")
	}
	return taskID, nil
}

func (c *TaskCenterClient) AppendProgress(ctx context.Context, taskType common.TaskType, taskID string, appendProgress float64) error {
	if !CheckInit() {
		return errors.New("[AppendProgress] client need init")
	}
	if err := syncer.NewSyncerService().AppendProgress(ctx, taskType, taskID, appendProgress); err != nil {
		return errors.New("[AppendProgress] AppendProgress error")
	}
	return nil
}

func (c *TaskCenterClient) GetResult(ctx context.Context, taskType common.TaskType, taskID string, bizResponse interface{}) (*common.TaskStateInfo, error) {
	if !CheckInit() {
		return nil, errors.New("[GetResult] client need init")
	}

	// 获取任务执行状态
	info, err := syncer.NewSyncerService().GetTaskStateInfo(ctx, taskType, taskID)
	if err != nil {
		return nil, errors.New("[GetResult] GetTaskStateInf error")
	}

	// 若未进行完，直接返回
	if !common.TaskStateIsTermination(info.State) {
		return info, nil
	}

	// 否则查询业务结果
	if err := syncer.NewSyncerService().GetBizResponse(ctx, taskType, taskID, &bizResponse); err != nil {
		return nil, errors.New("[GetResult] GetBizResponse error")
	}
	return info, nil
}
