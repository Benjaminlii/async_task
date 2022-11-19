package config

import "github.com/Benjaminlii/async_task/common"

var config Options

type Options struct {
	RocketMQConfig *RocketMQConfig
	RedisConfig    *RedisConfig
	HandlerMapping map[common.TaskType]*common.AsyncTaskHandler
}

func SetConfig(o Options) {
	config = o
}

func GetConfig() Options {
	return config
}

func ExplainOption(options ...OptionFunc) *Options {
	ops := &Options{}
	for _, do := range options {
		do.F(ops)
	}
	return ops
}

type OptionFunc struct {
	F func(*Options)
}

func WithRocketMQ(rocketMQConfig *RocketMQConfig) OptionFunc {
	return OptionFunc{func(op *Options) {
		op.RocketMQConfig = rocketMQConfig
	}}
}

func WithRedis(redisConfig *RedisConfig) OptionFunc {
	return OptionFunc{func(op *Options) {
		op.RedisConfig = redisConfig
	}}
}

func WithHandler(taskType common.TaskType, handler *common.AsyncTaskHandler) OptionFunc {
	return OptionFunc{func(op *Options) {
		if len(op.HandlerMapping) == 0 {
			op.HandlerMapping = make(map[common.TaskType]*common.AsyncTaskHandler, 4)
		}
		op.HandlerMapping[taskType] = handler
	}}
}
