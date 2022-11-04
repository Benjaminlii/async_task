package syncer

import (
	"code/benjamin/async_task/biz/config"
	"code/benjamin/async_task/biz/driver/redis"
)

func Init(config *config.Options) (err error) {
	if err = redis.InitRedis(config); err != nil {
		return err
	}
	return nil
}

func NewSyncerService() SyncerService {
	return &SyncerServiceDefaultImpl{}
}