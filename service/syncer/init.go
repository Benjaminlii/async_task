package syncer

import (
	"github.com/Benjaminlii/async_task/config"
	"github.com/Benjaminlii/async_task/driver/redis"
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