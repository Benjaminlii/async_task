package creator

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"code/benjamin/async_task/biz/common"
	"code/benjamin/async_task/biz/config"
	"code/benjamin/async_task/biz/driver/logger"
	"code/benjamin/async_task/biz/driver/rocketmq"
	"code/benjamin/async_task/biz/service/syncer"
	"code/benjamin/async_task/biz/utils"
)

type CreatorServiceDefaultImpl struct{}

func (cs *CreatorServiceDefaultImpl) CreateTask(ctx context.Context, taskType common.TaskType, bizRequest *interface{}, option *common.TaskAdditionalOption) (string, error) {
	// 校验及补充参数
	if _, isOk := config.GetConfig().HandlerMapping[taskType]; !isOk {
		return "", errors.New("[CreateTask] req taskType is illegal")
	}
	if option == nil {
		option = common.NewTaskAdditionalOption()
	}
	taskID := option.CustomTaskID
	if taskID == nil {
		taskID = utils.NewString(uuid.NewString())
	}
	// 初始化task Info
	if err := syncer.NewSyncerService().InitTaskInfo(ctx, taskType, *taskID, option); err != nil {
		return "", errors.New("[CreateTask] InitTaskInfo error")
	}
	// 存储req信息
	if err := syncer.NewSyncerService().SetBizRequest(ctx, taskType, *taskID, *bizRequest); err != nil {
		return "", errors.New("[CreateTask] SetBizRequest error")
	}
	// 发送mq消息
	msgID, err := rocketmq.SendMessage(ctx, *taskID, []string{string(taskType)})
	if err != nil {
		return "", errors.New("[CreateTask] SetBizRequest error")
	}
	err = syncer.NewSyncerService().UpdateTaskStateInfo(ctx, taskType, *taskID, nil, nil, map[string]string{
		common.ASYNC_TASK_MSG_ID: msgID,
	})
	if err != nil {
		logger.Warnf(ctx, "[CreateTask] set msg id failed, taskType:%v, taskID:%v, msgID:%v", taskType, taskID, msgID)
	}
	return *taskID, nil
}
