package syncer

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"

	"github.com/Benjaminlii/async_task/common"
	"github.com/Benjaminlii/async_task/driver/redis"
)

type SyncerServiceDefaultImpl struct{}

func (impl *SyncerServiceDefaultImpl) InitTaskInfo(ctx context.Context, taskType common.TaskType, taskID string, option *common.TaskAdditionalOption) error {
	taskStateInfo := &common.TaskStateInfo{
		State:      common.TASK_STATE_NOT_START,
		Progress:   0,
		ResultInfo: map[string]string{},
	}
	err := impl.SetTaskAdditionalOption(ctx, taskType, taskID, option)
	if err != nil {
		return errors.Wrap(err, "[InitTaskStateInfo] GetTaskAdditionalOption error")
	}

	redisKey := redis.GetTaskID2TaskStateKey(taskType, taskID)
	_, err = redis.Redis().Set(redisKey, *taskStateInfo, option.TaskStateInfoTimeout).Result()
	if err != nil {
		return errors.Wrap(err, "[InitTaskStateInfo] redis set error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) GetTaskStateInfo(ctx context.Context, taskType common.TaskType, taskID string) (*common.TaskStateInfo, error) {
	redisKey := redis.GetTaskID2TaskStateKey(taskType, taskID)
	dataBytes, err := redis.Redis().Get(redisKey).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "[GetTaskStateInfo] redis get error")
	}

	taskStateInfo := &common.TaskStateInfo{}
	if err = json.Unmarshal(dataBytes, taskStateInfo); err != nil {
		return nil, errors.Wrap(err, "[GetTaskStateInfo] json unmarshal error")
	}
	return taskStateInfo, nil
}

func (impl *SyncerServiceDefaultImpl) UpdateTaskStateInfo(ctx context.Context, taskType common.TaskType, taskID string, taskState *common.TaskState, appendProgress *float64, resultInfo map[string]string) error {
	if appendProgress != nil && *appendProgress < 0 {
		return errors.New("[UpdateTaskStateInfo] req appendProgress illegal")
	}

	taskStateInfo, err := impl.GetTaskStateInfo(ctx, taskType, taskID)
	if err != nil {
		return errors.Wrap(err, "[UpdateTaskStateInfo] GetTaskStateInfo error")
	}

	// verify state change
	if taskState != nil {
		if afterState, err := common.CalculateTaskState(taskStateInfo.State, *taskState); err != nil {
			return errors.Wrap(err, "[UpdateTaskStateInfo] CalculateTaskState error")
		} else {
			taskStateInfo.State = afterState
		}
	}

	// verify progress
	if appendProgress != nil && *appendProgress > 0 {
		taskStateInfo.Progress += *appendProgress
		if taskStateInfo.Progress > 1 {
			taskStateInfo.Progress = 1
		}
	}

	// append resultInfo
	for key := range resultInfo {
		taskStateInfo.ResultInfo[key] = resultInfo[key]
	}

	redisKey := redis.GetTaskID2TaskStateKey(taskType, taskID)

	var mutex sync.Mutex
	mutex.Lock()
	_, err = redis.Redis().Set(redisKey, *taskStateInfo, redis.Redis().TTL(redisKey).Val()).Result()
	mutex.Unlock()

	if err != nil {
		return errors.Wrap(err, "[UpdateTaskStateInfo] redis set error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) AppendProgress(ctx context.Context, taskType common.TaskType, taskID string, appendProgress float64) error {
	if err := impl.UpdateTaskStateInfo(ctx, taskType, taskID, nil, &appendProgress, nil); err != nil {
		return errors.Wrap(err, "[AppendProgress] updateTaskStateInfo error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) TaskStart(ctx context.Context, taskType common.TaskType, taskID string) error {
	if err := impl.UpdateTaskStateInfo(ctx, taskType, taskID, common.TaskStatePtr(common.TASK_STATE_PROCESSING), nil, nil); err != nil {
		return errors.Wrap(err, "[TaskStart] updateTaskStateInfo error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) TaskFinishSuccess(ctx context.Context, taskType common.TaskType, taskID string) error {
	if err := impl.UpdateTaskStateInfo(ctx, taskType, taskID, common.TaskStatePtr(common.TASK_STATE_SUCCESS), nil, nil); err != nil {
		return errors.Wrap(err, "[TaskFinishSuccess] updateTaskStateInfo error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) TaskFinishFailed(ctx context.Context, taskType common.TaskType, taskID string) error {
	if err := impl.UpdateTaskStateInfo(ctx, taskType, taskID, common.TaskStatePtr(common.TASK_STATE_FAILED), nil, nil); err != nil {
		return errors.Wrap(err, "[TaskFinishFailed] updateTaskStateInfo error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) SetBizRequest(ctx context.Context, taskType common.TaskType, taskID string, bizRequest interface{}) error {
	option, err := impl.getTaskAdditionalOption(ctx, taskType, taskID)
	if err != nil {
		return errors.Wrap(err, "[SetBizRequest] getTaskAdditionalOption error")
	}

	redisKey := redis.GetTaskID2TaskRequestKey(taskType, taskID)
	reqBytes, err := json.Marshal(bizRequest)
	if err != nil {
		return errors.Wrap(err, "[SetBizRequest] Marshal error")
	}
	_, err = redis.Redis().Set(redisKey, reqBytes, option.TaskStateInfoTimeout).Result()
	if err != nil {
		return errors.Wrap(err, "[SetBizRequest] redis set error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) GetBizRequest(ctx context.Context, taskType common.TaskType, taskID string) (string, error) {
	redisKey := redis.GetTaskID2TaskRespKey(taskType, taskID)
	dataStr, err := redis.Redis().Get(redisKey).Result()
	if err != nil {
		return "", errors.Wrap(err, "[GetBizRequest] redis get error")
	}
	return dataStr, nil
}

func (impl *SyncerServiceDefaultImpl) SetBizResponse(ctx context.Context, taskType common.TaskType, taskID string, bizResponse interface{}) error {
	option, err := impl.getTaskAdditionalOption(ctx, taskType, taskID)
	if err != nil {
		return errors.Wrap(err, "[SetBizResponse] GetTaskAdditionalOption error")
	}

	redisKey := redis.GetTaskID2TaskRespKey(taskType, taskID)

	respBytes, err := json.Marshal(bizResponse)
	if err != nil {
		return errors.Wrap(err, "[SetBizResponse] Marshal error")
	}

	_, err = redis.Redis().Set(redisKey, respBytes, option.TaskResultTimeout).Result()
	if err != nil {
		return errors.Wrap(err, "[SetBizResponse] redis set error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) GetBizResponse(ctx context.Context, taskType common.TaskType, taskID string, bizResponse *interface{}) error {
	redisKey := redis.GetTaskID2TaskRespKey(taskType, taskID)
	dataBytes, err := redis.Redis().Get(redisKey).Bytes()
	if err != nil {
		return errors.Wrap(err, "[GetBizResponse] redis get error")
	}
	if err = json.Unmarshal(dataBytes, bizResponse); err != nil {
		return errors.Wrap(err, "[GetBizResponse] json unmarshal error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) SetTaskAdditionalOption(ctx context.Context, taskType common.TaskType, taskID string, taskAdditionalOption *common.TaskAdditionalOption) error {
	redisKey := redis.GetTaskID2TaskOptionKey(taskType, taskID)
	if _, err := redis.Redis().Set(redisKey, *taskAdditionalOption, taskAdditionalOption.TaskStateInfoTimeout).Result(); err != nil {
		return errors.Wrap(err, "[SetTaskAdditionalOption] redis set error")
	}
	return nil
}

func (impl *SyncerServiceDefaultImpl) getTaskAdditionalOption(ctx context.Context, taskType common.TaskType, taskID string) (*common.TaskAdditionalOption, error) {
	redisKey := redis.GetTaskID2TaskOptionKey(taskType, taskID)
	dataBytes, err := redis.Redis().Get(redisKey).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "[GetTaskAdditionalOption] redis get error")
	}

	res := common.NewTaskAdditionalOption()

	if err = json.Unmarshal(dataBytes, res); err != nil {
		return nil, errors.Wrap(err, "[GetTaskAdditionalOption] json unmarshal error")
	}
	return res, nil
}
