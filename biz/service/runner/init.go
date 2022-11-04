package runner

import (
	"sync"

	"code/benjamin/async_task/biz/config"
	"code/benjamin/async_task/biz/driver/redis"
	"code/benjamin/async_task/biz/driver/rocketmq"
)

var (
	registerHandlerOnce sync.Once
)

func Init(config *config.Options) error {
	if err := rocketmq.InitRocketMQ(config); err != nil {
		return err
	}
	if err := redis.InitRedis(config); err != nil {
		return err
	}
	if err := NewRunnerService().RegisterHandler(config.HandlerMapping); err != nil {
		return err
	}
	return nil
}

func NewRunnerService() RunnerService {
	return &RunnerServiceDefaultImpl{}
}
