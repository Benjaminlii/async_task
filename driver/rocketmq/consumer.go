package rocketmq

import (
	"context"
	"strings"
	"sync"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/pkg/errors"
)

func RunConsumer(tags []string, bizFunc func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error), consumerShutdownCond *sync.Cond) {
	go func() {
		PushConsumer, err := rocketmq.NewPushConsumer(consumer.WithNameServer(ASyncTaskNameServers))
		if err != nil {
			panic(errors.Wrap(err, "[registerConsumer] init rocket mq push consumer error"))
		}
		var messageSelector consumer.MessageSelector
		if len(tags) > 0 {
			messageSelector = consumer.MessageSelector{
				Type:       consumer.TAG,
				Expression: strings.Join(tags, "||"),
			}
		}
		// 订阅topic
		err = PushConsumer.Subscribe(ASyncTaskTopic, messageSelector, bizFunc)
		if err != nil {
			panic(errors.Wrap(err, "[registerConsumer] Subscribe error"))
		}
		// 启动consumer
		if err = PushConsumer.Start(); err != nil {
			panic(errors.Wrap(err, "[registerConsumer] consumer start error"))
		}
		consumerShutdownCond.Wait()
		if err = PushConsumer.Shutdown(); err != nil {
			panic(errors.Wrap(err, "[registerConsumer] Subscribe error"))
		}
	}()
}
