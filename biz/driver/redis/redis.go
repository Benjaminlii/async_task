package redis

import (
	"code/benjamin/async_task/biz/common"
	"code/benjamin/async_task/biz/config"
	"fmt"
	"sync"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

var (
	redisOnce   sync.Once
	redisClient *redis.Client
)

func Redis() *redis.Client {
	return redisClient
}

func InitRedis(config *config.Options) error {
	var err error
	redisOnce.Do(func() {
		if config == nil {
			err = errors.New("[InitRedis] config is nil")
			return
		}
		redisConfig := config.RedisConfig
		if redisConfig == nil {
			err = errors.New("[InitRedis] redisConfig is nil")
			return
		}
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisConfig.Address,
			Password: redisConfig.Password,
		})

		if _, curErr := redisClient.Ping().Result(); curErr != nil {
			err = errors.Wrap(curErr, "[InitRedis] Ping failed")
			return
		}
	})
	return err
}

// 任务执行入参信息
func GetTaskID2TaskRequestKey(taskType common.TaskType, taskID string) string {
	return fmt.Sprintf("taskID_to_request_%s_%s", taskType, taskID)
}

// 任务执行状态
func GetTaskID2TaskStateKey(taskType common.TaskType, taskID string) string {
	return fmt.Sprintf("taskID_to_state_%s_%s", taskType, taskID)
}

// 任务附加选项
func GetTaskID2TaskOptionKey(taskType common.TaskType, taskID string) string {
	return fmt.Sprintf("taskID_to_option_%s_%s", taskType, taskID)
}

// 任务执行结果
func GetTaskID2TaskRespKey(taskType common.TaskType, taskID string) string {
	return fmt.Sprintf("taskID_to_resp_%s_%s", taskType, taskID)
}
