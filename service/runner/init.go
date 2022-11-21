package runner

import (
	"sync"

	"github.com/Benjaminlii/async_task/config"
	"github.com/Benjaminlii/async_task/driver/redis"
	"github.com/Benjaminlii/async_task/driver/rocketmq"
	"github.com/pkg/errors"
)

var (
	registerHandlerOnce sync.Once
	consumerCount       int
	shutdownChan        chan struct{}
)

func Init(config *config.Options) error {
	shutdownChan = make(chan struct{})
	if err := rocketmq.InitRocketMQ(config); err != nil {
		return errors.Wrap(err, "[Init] InitRocketMQ error")
	}
	if err := redis.InitRedis(config); err != nil {
		return errors.Wrap(err, "[Init] InitRedis error")
	}
	if err := NewRunnerService().RegisterHandler(config.HandlerMapping); err != nil {
		return errors.Wrap(err, "[Init] RegisterHandler error")
	}
	return nil
}

func Shutdown() {
	for i := 0; i < consumerCount; i++ {
		shutdownChan <- struct{}{}
	}
}

func NewRunnerService() RunnerService {
	return &RunnerServiceDefaultImpl{}
}
