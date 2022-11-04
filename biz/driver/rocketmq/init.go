package rocketmq

import (
	"code/benjamin/async_task/biz/config"
	"sync"

	"github.com/pkg/errors"
)

var (
	rocketMQOnce         sync.Once
	ASyncTaskTopic       = "default_async_task_topic"
	ASyncTaskNameServers = []string{}
)

func InitRocketMQ(config *config.Options) error {
	var err error
	rocketMQOnce.Do(func() {
		if config == nil {
			err = errors.New("[InitRocketMQ] config is nil")
			return
		}
		mqConfig := config.RocketMQConfig
		if config == nil {
			err = errors.New("[InitRocketMQ] mqConfig is nil")
			return
		}
		if mqConfig.Topic != "" {
			ASyncTaskTopic = mqConfig.Topic
		}
		if len(mqConfig.NameServers) > 0 {
			ASyncTaskNameServers = mqConfig.NameServers
		}
		if err = InitProducer(mqConfig); err != nil {
			err = errors.Wrap(err, "[InitRocketMQ] InitProducer error")
			return
		}
	})
	return err
}
