package runner

import (
	"sync"

	"github.com/Benjaminlii/async_task/config"
	"github.com/Benjaminlii/async_task/driver/redis"
	"github.com/Benjaminlii/async_task/driver/rocketmq"
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
