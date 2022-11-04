package creator

import (
	"code/benjamin/async_task/biz/config"
	"code/benjamin/async_task/biz/driver/redis"
	"code/benjamin/async_task/biz/driver/rocketmq"
)

func Init(config *config.Options) (err error) {
	if err = rocketmq.InitRocketMQ(config); err != nil {
		return err
	}
	if err = redis.InitRedis(config); err != nil {
		return err
	}
	return nil
}

func NewCreatorService() CreatorService {
	return &CreatorServiceDefaultImpl{}
}
