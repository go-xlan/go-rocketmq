package rocketmq

import (
	"context"
	"fmt"
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

type Config struct {
	NameServer string
	GroupName  string
	Producer   *ProducerConfig
}

type ProducerConfig struct {
	SendMsgTimeout int32
	RetryTimes     int32
}

type ProducerClient struct {
	mqProducer rocketmq.Producer
}

func NewProducerClient(config *Config) (*ProducerClient, error) {
	zaplog.LOG.Debug("rocketmq-name-server", zap.String("nameServer", config.NameServer))
	nameSrvAddress := rese.A1(ResolveNameServer(config.NameServer))
	zaplog.LOG.Debug("rocketmq-name-server", zap.Strings("srvAddress", nameSrvAddress))

	mqProducer, err := producer.NewDefaultProducer(
		producer.WithNameServer(nameSrvAddress),
		producer.WithGroupName(must.Nice(config.GroupName)),
		producer.WithRetry(int(must.Nice(config.Producer.RetryTimes))),
		producer.WithSendMsgTimeout(time.Duration(config.Producer.SendMsgTimeout)*time.Second),
	)
	if err != nil {
		return nil, erero.Wro(err)
	}
	if err = mqProducer.Start(); err != nil {
		return nil, erero.Wro(err)
	}
	res := &ProducerClient{
		mqProducer: mqProducer,
	}
	return res, nil
}

func (c *ProducerClient) Close() {
	must.Done(c.mqProducer.Shutdown())
}

func (c *ProducerClient) SendMsg(ctx context.Context, topic string, msgBytes []byte) error {
	zaplog.LOG.Debug("sending rocket-mq-message")

	msg := primitive.NewMessage(topic, msgBytes)
	// msg.WithTag("tagA") // 根据自己的需求选择是否启用标签功能

	result, err := c.mqProducer.SendSync(ctx, msg)
	if err != nil {
		return erero.Wro(err)
	}

	zaplog.LOG.Debug("rocket-mq-message sent success", zap.Any("result", result))
	return nil
}

func ResolveNameServer(nameServer string) ([]string, error) {
	args := strings.Split(nameServer, ":")
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong name_server format: %s", nameServer)
	}
	host, port := args[0], args[1]
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("parse %s wrong: %v", host, err)
	}
	var results []string
	for _, ip := range ips {
		if ip.To4() != nil { // 使用 IPv4 地址
			results = append(results, ip.String()+":"+port)
		}
	}
	return results, nil
}
