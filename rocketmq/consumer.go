package rocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

type ConsumerClient struct {
	mqConsumer rocketmq.PushConsumer
}

func NewConsumerClient(config *Config) (*ConsumerClient, error) {
	zaplog.LOG.Debug("rocketmq-name-server", zap.String("nameServer", config.NameServer))
	nameSrvAddress := rese.A1(ResolveNameServer(config.NameServer))
	zaplog.LOG.Debug("rocketmq-name-server", zap.Strings("srvAddress", nameSrvAddress))

	mqConsumer, err := consumer.NewPushConsumer(
		consumer.WithNameServer(nameSrvAddress),
		consumer.WithGroupName(config.GroupName),
		// consumer.WithConsumerModel(consumer.Clustering), //这个是默认值没必要再设置
		consumer.WithMaxReconsumeTimes(3),
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return &ConsumerClient{
		mqConsumer: mqConsumer,
	}, nil
}

func (c *ConsumerClient) Close() {
	must.Done(c.mqConsumer.Shutdown())
}

func (c *ConsumerClient) StartConsume(topic string, consumeFunc func(msg *primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	if err := c.mqConsumer.Start(); err != nil {
		return erero.Wro(err)
	}
	zaplog.LOG.Debug("Consumer start success, waiting messages...")

	if err := c.mqConsumer.Subscribe(topic, consumer.MessageSelector{
		// Type:       consumer.TAG,
		// Expression: "tagA", //根据自己的需求选择是否启用标签功能
	}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, msg := range msgs {
			zaplog.LOG.Debug("consume message", zap.String("msgID", msg.MsgId))
			if result, err := consumeFunc(msg); err != nil {
				return result, erero.Wro(err)
			}
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		return erero.Wro(err)
	}
	return nil
}
