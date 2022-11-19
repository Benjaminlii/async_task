package creator

import (
	"github.com/Benjaminlii/async_task/config"
	"github.com/Benjaminlii/async_task/driver/redis"
	"github.com/Benjaminlii/async_task/driver/rocketmq"
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
