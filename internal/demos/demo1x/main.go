package main

import (
	"context"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-xlan/go-rocketmq/rocketmq"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
	"github.com/yyle88/zaplog"
)

func main() {
	config := &rocketmq.Config{
		NameServerAddress: "127.0.0.1:9876",
		GroupName:         "TestGroup",
		ProducerOptions: &rocketmq.ProducerOptions{
			SendMessageTimeout: 3 * time.Second,
			RetryAttempts:      3,
		},
	}

	const topic = "TestTopic"

	consumerCli := rese.P1(rocketmq.NewConsumer(config))
	defer rese.F0(consumerCli.Close)
	must.Done(consumerCli.StartSubscribe(topic, func(message *primitive.MessageExt) (consumer.ConsumeResult, error) {
		zaplog.SUG.Debugln("consume message body:", string(message.Body))
		return consumer.ConsumeSuccess, nil
	}))

	producerCli := rese.P1(rocketmq.NewProducer(config))
	defer rese.F0(producerCli.Close)
	for idx := 0; idx < 10000; idx++ {
		type MessagePayload struct {
			Name       string `json:"name"`
			SequenceNo int64  `json:"sequenceNo"`
		}

		payload := neatjsons.S(&MessagePayload{
			Name:       "demo",
			SequenceNo: int64(idx),
		})
		must.Done(producerCli.SendMessage(context.Background(), topic, []byte(payload)))
		time.Sleep(time.Second)
	}

	select {}
}
