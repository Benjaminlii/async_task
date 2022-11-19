package common

import (
	"errors"
	"time"
)

// 任务类型，可扩展 为支持自定义taskID 保证taskID不重复 tasktype不应为纯数字
type TaskType string

// 异步任务附加参数
type TaskAdditionalOption struct {
	TaskStateInfoTimeout time.Duration // 任务状态信息的超时时间，默认12小时
	TaskResultTimeout    time.Duration // 任务执行结果信息的超时时间，默认12小时
	CustomTaskID         *string       // 自定义TaskID, 若启用调用方需自己保证自己id在type内t唯一
}

func NewTaskAdditionalOption() *TaskAdditionalOption {
	return &TaskAdditionalOption{
		TaskStateInfoTimeout: time.Hour * 12,
		TaskResultTimeout:    time.Hour * 12,
		CustomTaskID:         nil,
	}
}

// 异步任务状态
type TaskState int64

func TaskStatePtr(ts TaskState) *TaskState {
	return &ts
}

var taskStateIsTerminationMap = map[TaskState]bool{
	TASK_STATE_SUCCESS:   true,
	TASK_STATE_FAILED:    true,
	TASK_STATE_SYS_ERROR: true,
}

func TaskStateIsTermination(ts TaskState) bool {
	return taskStateIsTerminationMap[ts]
}

var calculaterMap = map[TaskState]map[TaskState]bool{
	TASK_STATE_NOT_START: {
		TASK_STATE_PROCESSING: true,
		TASK_STATE_SYS_ERROR:  true,
	},
	TASK_STATE_PROCESSING: {
		TASK_STATE_PROCESSING: true,
		TASK_STATE_SUCCESS:    true,
		TASK_STATE_FAILED:     true,
		TASK_STATE_SYS_ERROR:  true,
	},
}

// 任务计算状态机
func CalculateTaskState(beforeState TaskState, afterState TaskState) (TaskState, error) {
	if secondMap, isOk := calculaterMap[beforeState]; isOk {
		if result := secondMap[afterState]; result {
			return afterState, nil
		}
	}
	return beforeState, errors.New("[CalculateTaskState] calculate failed")
}

// 任务执行状态信息
type TaskStateInfo struct {
	State      TaskState         // 处理状态
	Progress   float64           // 处理进度：0~1的float数字（1表示100%）
	ResultInfo map[string]string // 业务自定义结果信息
}
