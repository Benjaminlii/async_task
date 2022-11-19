package rocketmq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/pkg/errors"

	"github.com/Benjaminlii/async_task/config"
	"github.com/Benjaminlii/async_task/driver/logger"
)

const (
	RetryTime = 3
	SleepTime = time.Millisecond * 10
)

var (
	Producer rocketmq.Producer
)

func InitProducer(mqConfig *config.RocketMQConfig) error {
	var err error
	if Producer, err = rocketmq.NewProducer(producer.WithNameServer(mqConfig.NameServers)); err != nil {
		return errors.Wrap(err, "[InitRocketMQ] init rocket mq producer error")
	}
	if err = Producer.Start(); err != nil {
		return errors.Wrap(err, "[InitRocketMQ] producer mq start error")
	}
	return nil
}

func SendMessage(ctx context.Context, body string, tags []string) (string, error) {
	msg := primitive.NewMessage(ASyncTaskTopic, []byte(body))
	if len(tags) > 0 {
		msg = msg.WithTag(strings.Join(tags, "||"))
	}
	for i := 0; i < RetryTime; i++ {
		res, err := Producer.SendSync(ctx, msg)
		if err != nil {
			time.Sleep(SleepTime)
			continue
		}
		logger.Infof(ctx, "[RocketMQ][SendMessage] send success, msgID:%s", res.MsgID)
		return res.MsgID, nil
	}
	return "", errors.New(fmt.Sprintf("[SendMessage] send failed, tags:%v, body:%v", tags, body))
}
