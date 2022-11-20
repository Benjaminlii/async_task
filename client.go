package async_task

import (
	"context"

	"github.com/pkg/errors"

	"github.com/Benjaminlii/async_task/common"
	"github.com/Benjaminlii/async_task/config"
	"github.com/Benjaminlii/async_task/service/creator"
	"github.com/Benjaminlii/async_task/service/syncer"
)

type TaskCenterClient struct {
	Options *config.Options
}

func (c *TaskCenterClient) Create(ctx context.Context, taskType common.TaskType, bizRequest *interface{}, option *common.TaskAdditionalOption) (string, error) {
	if !checkInit() {
		return "", errors.New("[Create] client need init")
	}
	taskID, err := creator.NewCreatorService().CreateTask(ctx, taskType, bizRequest, option)
	if err != nil {
		return "", errors.Wrap(err, "[Create] CreateTask error")
	}
	return taskID, nil
}

func (c *TaskCenterClient) AppendProgress(ctx context.Context, taskType common.TaskType, taskID string, appendProgress float64) error {
	if !checkInit() {
		return errors.New("[AppendProgress] client need init")
	}
	if err := syncer.NewSyncerService().AppendProgress(ctx, taskType, taskID, appendProgress); err != nil {
		return errors.New("[AppendProgress] AppendProgress error")
	}
	return nil
}

func (c *TaskCenterClient) GetResult(ctx context.Context, taskType common.TaskType, taskID string, bizResponse interface{}) (*common.TaskStateInfo, error) {
	if !checkInit() {
		return nil, errors.New("[GetResult] client need init")
	}

	// 获取任务执行状态
	info, err := syncer.NewSyncerService().GetTaskStateInfo(ctx, taskType, taskID)
	if err != nil {
		return nil, errors.Wrap(err, "[GetResult] GetTaskStateInf error")
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
