package common

const (
	ASYNC_TASK_MSG_ID = "MSG_ID"
)

const (
	TASK_STATE_NOT_START  TaskState = iota // 进行中
	TASK_STATE_PROCESSING                  // 进行中

	TASK_STATE_SUCCESS   // 完成并成功
	TASK_STATE_FAILED    // 完成但失败
	TASK_STATE_SYS_ERROR // 系统异常
)
