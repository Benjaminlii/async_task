package runner

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/pkg/errors"

	"github.com/Benjaminlii/async_task/common"
	"github.com/Benjaminlii/async_task/driver/logger"
	"github.com/Benjaminlii/async_task/driver/rocketmq"
	"github.com/Benjaminlii/async_task/service/syncer"
)

type RunnerServiceDefaultImpl struct{}

func (impl *RunnerServiceDefaultImpl) RegisterHandler(handlerMapping map[common.TaskType]common.AsyncTaskHandler) error {
	var err error
	registerHandlerOnce.Do(func() {
		if len(handlerMapping) == 0 {
			err = errors.New("[registerHandler] handlerMapping length is 0")
			return
		}
		for taskType := range handlerMapping {
			handler := handlerMapping[taskType]
			// 按照tag消费，未注册的taskType和handler不会执行
			tags := []string{string(taskType)}
			f := func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
				for i := range msgs {
					taskID := msgs[i].Body
					logger.Infof(ctx, "[consumer running] taskID:%s", string(taskID))
					if err := impl.runTask(ctx, taskType, string(taskID), handler); err != nil {
						return consumer.ConsumeRetryLater, errors.Wrap(err, "[consumer running] runTask error")
					}
				}
				return consumer.ConsumeSuccess, nil
			}
			rocketmq.RunConsumer(tags, f, shutdownChan)
			consumerCount++
		}
	})

	return err
}

func (impl *RunnerServiceDefaultImpl) runTask(ctx context.Context, taskType common.TaskType, taskID string, handler common.AsyncTaskHandler) error {
	handleReq, err := syncer.NewSyncerService().GetBizRequest(ctx, taskType, taskID)
	if err != nil {
		return errors.Wrap(err, "[runTask] GetBizRequest error")
	}

	// 更新taskState为进行中
	if err = syncer.NewSyncerService().TaskStart(ctx, taskType, taskID); err != nil {
		return errors.Wrap(err, "[runTask] UpdateTaskStateInfo error")
	}
	handleResp, err := handler.HandleMessage(ctx, handleReq)
	if err != nil {
		// 更新taskState为finish_Failed
		if err = syncer.NewSyncerService().TaskFinishFailed(ctx, taskType, taskID); err != nil {
			return errors.Wrap(err, "[runTask] TaskFinishFailed error")
		}
		return errors.Wrap(err, "[runTask] run handleMessage error")
	}
	// 更新taskState为finish_success
	if err = syncer.NewSyncerService().TaskFinishSuccess(ctx, taskType, taskID); err != nil {
		return errors.Wrap(err, "[runTask] TaskFinishSuccess error")
	}

	// 持久化resp
	if err = syncer.NewSyncerService().SetBizResponse(ctx, taskType, taskID, handleResp); err != nil {
		return errors.Wrap(err, "[runTask] SetBizResult error")
	}
	return nil
}
