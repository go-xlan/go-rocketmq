package rocketmq

import (
	"context"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/yyle88/erero"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// Consumer represents a RocketMQ push consumer client
// Encapsulates the underlying push consumer instance with simplified subscription interface
// Provides automatic message pulling and callback-based message handling
//
// Consumer 代表 RocketMQ 推送消费者客户端
// 封装底层推送消费者实例并提供简化的订阅接口
// 提供自动消息拉取和基于回调的消息处理
type Consumer struct {
	instance rocketmq.PushConsumer // Underlying RocketMQ push consumer instance // 底层 RocketMQ 推送消费者实例
}

// NewConsumer creates and initializes a new RocketMQ consumer instance
// Resolves NameServer address and configures consumer with group name and retry settings
// Returns configured consumer instance (not yet started) or error on configuration failure
//
// NewConsumer 创建并初始化新的 RocketMQ 消费者实例
// 解析 NameServer 地址并配置消费者的组名和重试设置
// 返回已配置的消费者实例（尚未启动）或在配置失败时返回错误
func NewConsumer(config *Config) (*Consumer, error) {
	zaplog.LOG.Debug("resolve name server address", zap.String("nameServer", config.NameServerAddress))
	addresses := rese.A1(ResolveNameServer(config.NameServerAddress))
	zaplog.LOG.Debug("name server resolved", zap.Strings("addresses", addresses))

	instance, err := consumer.NewPushConsumer(
		consumer.WithNameServer(addresses),
		consumer.WithGroupName(config.GroupName),
		// consumer.WithConsumerModel(consumer.Clustering), // Default value, no need to set // 默认值，无需设置
		consumer.WithMaxReconsumeTimes(3),
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	return &Consumer{instance: instance}, nil
}

// Close shuts down the consumer and releases all resources
// Returns error if shutdown process encounters issues
//
// Close 关闭消费者并释放所有资源
// 如果关闭过程遇到问题则返回错误
func (c *Consumer) Close() error {
	return c.instance.Shutdown()
}

// StartSubscribe starts the consumer and subscribes to the specified topic with message handler
// First starts the consumer instance, then subscribes to topic with callback function
// The handler function is invoked for each received message and should return consume result
// Returns error if consumer start fails or subscription setup encounters issues
//
// StartSubscribe 启动消费者并使用消息处理器订阅指定主题
// 首先启动消费者实例，然后使用回调函数订阅主题
// 处理器函数会对每个接收到的消息调用，应返回消费结果
// 如果消费者启动失败或订阅设置遇到问题则返回错误
func (c *Consumer) StartSubscribe(topic string, handler func(message *primitive.MessageExt) (consumer.ConsumeResult, error)) error {
	zaplog.LOG.Debug("start consumer instance", zap.String("topic", topic))
	if err := c.instance.Start(); err != nil {
		return erero.Wro(err)
	}
	zaplog.LOG.Debug("consumer is running", zap.String("topic", topic))

	zaplog.LOG.Debug("subscribe to topic", zap.String("topic", topic))
	if err := c.instance.Subscribe(topic, consumer.MessageSelector{
		// Type:       consumer.TAG,
		// Expression: "tagA", // Enable tag feature based on your needs // 根据需求启用标签功能
	}, func(ctx context.Context, messages ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, message := range messages {
			zaplog.LOG.Debug("receive message", zap.String("msgID", message.MsgId), zap.String("topic", message.Topic))
			if result, err := handler(message); err != nil {
				return result, erero.Wro(err)
			}
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		return erero.Wro(err)
	}
	zaplog.LOG.Debug("subscription complete", zap.String("topic", topic))
	return nil
}
