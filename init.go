package async_task

import (
	"sync"

	"github.com/Benjaminlii/async_task/config"
	"github.com/Benjaminlii/async_task/service/creator"
	"github.com/Benjaminlii/async_task/service/runner"
	"github.com/Benjaminlii/async_task/service/syncer"
)

type InitFun func(options *config.Options) error

var (
	client         *TaskCenterClient
	clientInitOnce sync.Once
	hasInit        bool
	initFunChain   = []InitFun{
		creator.Init,
		runner.Init,
		syncer.Init,
	}
)

func Init(optionFuncs ...config.OptionFunc) (*TaskCenterClient, error) {
	var err error
	clientInitOnce.Do(func() {
		client = &TaskCenterClient{
			Options: config.ExplainOption(optionFuncs...),
		}
		for _, initF := range initFunChain {
			err = initF(client.Options)
			if err != nil {
				return
			}
		}
		config.SetConfig(*client.Options)
		hasInit = true
	})
	return client, err
}

func checkInit() bool {
	return hasInit
}
