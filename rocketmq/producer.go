// Package rocketmq: Apache RocketMQ client wrapper with simplified producer and consumer interfaces
// Provides elegant Go integration for RocketMQ message queue operations with mate ecosystem patterns
// Features automatic hostname resolution, IPv4 address support, and seamless error handling
//
// rocketmq: Apache RocketMQ 客户端封装，提供简化的生产者和消费者接口
// 为 RocketMQ 消息队列操作提供优雅的 Go 集成，采用 mate 生态模式
// 具有自动主机名解析、IPv4 地址支持和无缝错误处理功能
package rocketmq

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/yyle88/erero"
	"github.com/yyle88/must"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
	"go.uber.org/zap"
)

// Config represents the RocketMQ client configuration for both producer and consumer
// Contains essential connection parameters and producer-specific options
//
// Config 代表 RocketMQ 客户端配置，用于生产者和消费者
// 包含基本连接参数和生产者特定选项
type Config struct {
	NameServerAddress string           // NameServer address in format host:port // NameServer 地址，格式为 host:port
	GroupName         string           // Consumer or producer group name // 消费者或生产者组名
	ProducerOptions   *ProducerOptions // Producer specific configuration options // 生产者特定配置选项
}

// ProducerOptions contains configuration options specific to message producer
// Defines timeout and retry behavior for message sending operations
//
// ProducerOptions 包含消息生产者特定的配置选项
// 定义消息发送操作的超时和重试行为
type ProducerOptions struct {
	SendMessageTimeout time.Duration // Message sending timeout duration // 消息发送超时时长
	RetryAttempts      int           // Number of retry attempts on send failure // 发送失败时的重试次数
}

// Producer represents a RocketMQ message producer client
// Encapsulates the underlying RocketMQ producer instance with simplified interface
//
// Producer 代表 RocketMQ 消息生产者客户端
// 封装底层 RocketMQ 生产者实例并提供简化接口
type Producer struct {
	instance rocketmq.Producer // Underlying RocketMQ producer instance // 底层 RocketMQ 生产者实例
}

// NewProducer creates and initializes a new RocketMQ producer instance
// Resolves NameServer address, configures producer with retry and timeout settings, and starts the producer
// Returns configured and running producer instance or error on configuration/connection failure
//
// NewProducer 创建并初始化新的 RocketMQ 生产者实例
// 解析 NameServer 地址，配置生产者的重试和超时设置，并启动生产者
// 返回已配置并运行的生产者实例，或在配置/连接失败时返回错误
func NewProducer(config *Config) (*Producer, error) {
	zaplog.LOG.Debug("resolve name server address", zap.String("nameServer", config.NameServerAddress))
	addresses := rese.A1(ResolveNameServer(config.NameServerAddress))
	zaplog.LOG.Debug("name server resolved", zap.Strings("addresses", addresses))

	instance, err := producer.NewDefaultProducer(
		producer.WithNameServer(addresses),
		producer.WithGroupName(must.Nice(config.GroupName)),
		producer.WithRetry(must.Nice(config.ProducerOptions.RetryAttempts)),
		producer.WithSendMsgTimeout(must.Nice(config.ProducerOptions.SendMessageTimeout)),
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	if err = instance.Start(); err != nil {
		return nil, erero.Wro(err)
	}
	return &Producer{instance: instance}, nil
}

// Close shuts down the producer and releases all resources
// Returns error if shutdown process encounters issues
//
// Close 关闭生产者并释放所有资源
// 如果关闭过程遇到问题则返回错误
func (p *Producer) Close() error {
	return p.instance.Shutdown()
}

// SendMessage sends a message with given payload to the specified topic
// Uses synchronous sending mode with context for timeout control
// Returns error if message sending fails or times out
//
// SendMessage 将给定载荷的消息发送到指定主题
// 使用同步发送模式，通过 context 控制超时
// 如果消息发送失败或超时则返回错误
func (p *Producer) SendMessage(ctx context.Context, topic string, payload []byte) error {
	zaplog.LOG.Debug("send message to topic", zap.String("topic", topic), zap.Int("size", len(payload)))

	message := primitive.NewMessage(topic, payload)
	// message.WithTag("tagA") // Enable tag feature based on your needs // 根据需求启用标签功能

	result, err := p.instance.SendSync(ctx, message)
	if err != nil {
		return erero.Wro(err)
	}

	zaplog.LOG.Debug("message sent", zap.String("msgID", result.MsgID), zap.Any("status", result.Status))
	return nil
}

// ResolveNameServer resolves NameServer hostname to IPv4 addresses with port
// Accepts input in format "hostname:port" and performs DNS resolution
// Returns slice of IPv4 addresses with port or error on invalid format or resolution failure
// Filters out IPv6 addresses to ensure compatibility with RocketMQ client requirements
//
// ResolveNameServer 将 NameServer 主机名解析为带端口的 IPv4 地址
// 接受 "hostname:port" 格式的输入并执行 DNS 解析
// 返回带端口的 IPv4 地址切片，或在格式无效或解析失败时返回错误
// 过滤掉 IPv6 地址以确保与 RocketMQ 客户端要求兼容
func ResolveNameServer(nameServer string) ([]string, error) {
	parts := strings.Split(nameServer, ":")
	if len(parts) != 2 {
		return nil, erero.Errorf("invalid name server format: %s (expected host:port)", nameServer)
	}
	hostname, port := parts[0], parts[1]
	ipAddresses, err := net.LookupIP(hostname)
	if err != nil {
		return nil, erero.Errorf("resolve hostname %s: %v", hostname, err)
	}
	var addresses []string
	for _, ip := range ipAddresses {
		if ip.To4() != nil { // Use IPv4 addresses only // 仅使用 IPv4 地址
			addresses = append(addresses, ip.String()+":"+port)
		}
	}
	return addresses, nil
}
