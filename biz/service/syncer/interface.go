package syncer

import (
	"context"

	"code/benjamin/async_task/biz/common"
)

type SyncerService interface {
	// 初始化任务状态
	InitTaskInfo(ctx context.Context, taskType common.TaskType, taskID string, taskAdditionalOption *common.TaskAdditionalOption) error
	// 查询任务状态
	GetTaskStateInfo(ctx context.Context, taskType common.TaskType, taskID string) (*common.TaskStateInfo, error)
	// 更新任务状态 taskState：任务状态（不更新可传nil）;progress：处理进度（0.0~1.0，不更新可传nil);resultInfo：任务自定义结果信息（可传nil)
	UpdateTaskStateInfo(ctx context.Context, taskType common.TaskType, taskID string, taskState *common.TaskState, progress *float64, resultInfo map[string]string) error

	// 追加进度
	AppendProgress(ctx context.Context, taskType common.TaskType, taskID string, appendProgress float64) error
	// 开始任务
	TaskStart(ctx context.Context, taskType common.TaskType, taskID string) error
	// 同步成功完成任务
	TaskFinishSuccess(ctx context.Context, taskType common.TaskType, taskID string) error
	// 同步失败完成任务
	TaskFinishFailed(ctx context.Context, taskType common.TaskType, taskID string) error

	// 设置任务附加选项
	SetTaskAdditionalOption(ctx context.Context, taskType common.TaskType, taskID string, taskAdditionalOption *common.TaskAdditionalOption) error
	// 获取任务附加选项
	getTaskAdditionalOption(ctx context.Context, taskType common.TaskType, taskID string) (*common.TaskAdditionalOption, error)

	// 设置任务执行业务参数
	SetBizRequest(ctx context.Context, taskType common.TaskType, taskID string, bizRequest interface{}) error
	// 获取任务执行业务参数
	GetBizRequest(ctx context.Context, taskType common.TaskType, taskID string) (string, error)

	// 设置业务结果
	SetBizResponse(ctx context.Context, taskType common.TaskType, taskID string, bizResponse interface{}) error
	// 获取业务结果
	GetBizResponse(ctx context.Context, taskType common.TaskType, taskID string, bizResponse *interface{}) error
}
