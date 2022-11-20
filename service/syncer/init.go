package syncer

import (
	"github.com/Benjaminlii/async_task/config"
	"github.com/Benjaminlii/async_task/driver/redis"
	"github.com/pkg/errors"
)

func Init(config *config.Options) (err error) {
	if err = redis.InitRedis(config); err != nil {
		return errors.Wrap(err, "[Init] InitRedis error")
	}
	return nil
}

func NewSyncerService() SyncerService {
	return &SyncerServiceDefaultImpl{}
}
