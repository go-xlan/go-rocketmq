package main

import (
	"context"
	"time"

	"github.com/go-xlan/go-rocketmq/rocketmq"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
	"github.com/yyle88/rese"
)

func main() {
	config := &rocketmq.Config{
		NameServer: "127.0.0.1:9876",
		GroupName:  "TestGroup",
		Producer: &rocketmq.ProducerConfig{
			SendMsgTimeout: 3000,
			RetryTimes:     3,
		},
	}

	const topic = "TestTopic"

	consumerClient := rese.P1(rocketmq.NewConsumerClient(config))
	defer consumerClient.Close()
	must.Done(consumerClient.StartReceiveMessage(topic))

	producerClient := rese.P1(rocketmq.NewProducerClient(config))
	defer producerClient.Close()
	for i := 0; i < 10000; i++ {
		type MessageType struct {
			Name       string `json:"name"`
			SequenceNo int64  `json:"sequenceNo"`
		}

		message := neatjsons.S(&MessageType{
			Name:       "demo",
			SequenceNo: int64(i),
		})
		must.Done(producerClient.SendMsg(context.Background(), topic, []byte(message)))
		time.Sleep(time.Second)
	}

	select {}
}
