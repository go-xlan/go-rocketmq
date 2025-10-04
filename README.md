[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-xlan/go-rocketmq/release.yml?branch=main&label=BUILD)](https://github.com/go-xlan/go-rocketmq/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-xlan/go-rocketmq)](https://pkg.go.dev/github.com/go-xlan/go-rocketmq)
[![Coverage Status](https://img.shields.io/coveralls/github/go-xlan/go-rocketmq/main.svg)](https://coveralls.io/github/go-xlan/go-rocketmq?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.23+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-xlan/go-rocketmq.svg)](https://github.com/go-xlan/go-rocketmq/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-xlan/go-rocketmq)](https://goreportcard.com/report/github.com/go-xlan/go-rocketmq)

# go-rocketmq

Apache RocketMQ client wrapper with simplified producer and consumer interfaces for Go.

---

<!-- TEMPLATE (EN) BEGIN: LANGUAGE NAVIGATION -->
## CHINESE README

[中文说明](README.zh.md)
<!-- TEMPLATE (EN) END: LANGUAGE NAVIGATION -->

## Main Features

🚀 **Simplified API**: Clean producer and consumer API for RocketMQ message queue operations
🔌 **Auto Resolution**: Automatic hostname resolution with IPv4 address support
🛡️ **Seamless Errors**: Elegant error handling with mate ecosystem integration
📝 **Go Practices**: Interface design following Go best practices
🔄 **Auto Pulling**: Automatic message pulling with callback-based message handling
⚙️ **Configurable**: Timeout and retry settings for message sending operations

## Installation

```bash
go get github.com/go-xlan/go-rocketmq
```

## Prerequisites

- Go 1.23.0 or higher
- Apache RocketMQ server instance (NameServer accessible)

## Usage

### Complete Producer and Consumer Example

This example shows how to create both a producer and consumer in the same program.

```go
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
```

⬆️ **Source:** [Source](internal/demos/demo1x/main.go)

## API Reference

### Configuration

- `Config` - Client configuration for both producer and consumer
- `ProducerOptions` - Producer specific configuration options

### Producer

- `NewProducer(config *Config) (*Producer, error)` - Creates and starts a new producer instance
- `SendMessage(ctx context.Context, topic string, payload []byte) error` - Sends message to specified topic
- `Close() error` - Shuts down producer and releases resources

### Consumer

- `NewConsumer(config *Config) (*Consumer, error)` - Creates a new consumer instance
- `StartSubscribe(topic string, handler func(*primitive.MessageExt) (consumer.ConsumeResult, error)) error` - Starts consumer and subscribes to topic
- `Close() error` - Shuts down consumer and releases resources

### Utilities

- `ResolveNameServer(nameServer string) ([]string, error)` - Resolves NameServer hostname to IPv4 addresses

## Examples

### Creating Producer

**Basic producer setup:**
```go
config := &rocketmq.Config{
    NameServerAddress: "127.0.0.1:9876",
    GroupName:         "ProducerGroup",
    ProducerOptions: &rocketmq.ProducerOptions{
        SendMessageTimeout: 3 * time.Second,
        RetryAttempts:      3,
    },
}
producerCli := rese.P1(rocketmq.NewProducer(config))
defer rese.F0(producerCli.Close)
```

**Sending messages:**
```go
payload := []byte(`{"name":"demo","value":123}`)
must.Done(producerCli.SendMessage(context.Background(), "MyTopic", payload))
```

### Creating Consumer

**Basic consumer setup:**
```go
config := &rocketmq.Config{
    NameServerAddress: "127.0.0.1:9876",
    GroupName:         "ConsumerGroup",
}
consumerCli := rese.P1(rocketmq.NewConsumer(config))
defer rese.F0(consumerCli.Close)
```

**Subscribing to topics:**
```go
must.Done(consumerCli.StartSubscribe("MyTopic", func(message *primitive.MessageExt) (consumer.ConsumeResult, error) {
    fmt.Printf("Message: %s\n", string(message.Body))
    return consumer.ConsumeSuccess, nil
}))
```

### Hostname Resolution

**Resolve NameServer hostname to IPv4 addresses:**
```go
addresses := rese.V1(rocketmq.ResolveNameServer("mq.example.com:9876"))
fmt.Printf("Resolved addresses: %v\n", addresses)
```

## Dependencies

This project uses the mate ecosystem for elegant error handling and resource management:

- `github.com/apache/rocketmq-client-go/v2` - Official Apache RocketMQ Go client
- `github.com/yyle88/erero` - Error handling
- `github.com/yyle88/must` - Panic-based assertion for errors
- `github.com/yyle88/rese` - Resource management with panic handling
- `github.com/yyle88/zaplog` - Structured logging
- `go.uber.org/zap` - Logging foundation

<!-- TEMPLATE (EN) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 License

MIT License. See [LICENSE](LICENSE).

---

## 🤝 Contributing

Contributions are welcome! Report bugs, suggest features, and contribute code:

- 🐛 **Found a mistake?** Open an issue on GitHub with reproduction steps
- 💡 **Have a feature idea?** Create an issue to discuss the suggestion
- 📖 **Documentation confusing?** Report it so we can improve
- 🚀 **Need new features?** Share the use cases to help us understand requirements
- ⚡ **Performance issue?** Help us optimize through reporting slow operations
- 🔧 **Configuration problem?** Ask questions about complex setups
- 📢 **Follow project progress?** Watch the repo to get new releases and features
- 🌟 **Success stories?** Share how this package improved the workflow
- 💬 **Feedback?** We welcome suggestions and comments

---

## 🔧 Development

New code contributions, follow this process:

1. **Fork**: Fork the repo on GitHub (using the webpage UI).
2. **Clone**: Clone the forked project (`git clone https://github.com/yourname/repo-name.git`).
3. **Navigate**: Navigate to the cloned project (`cd repo-name`)
4. **Branch**: Create a feature branch (`git checkout -b feature/xxx`).
5. **Code**: Implement the changes with comprehensive tests
6. **Testing**: (Golang project) Ensure tests pass (`go test ./...`) and follow Go code style conventions
7. **Documentation**: Update documentation to support client-facing changes and use significant commit messages
8. **Stage**: Stage changes (`git add .`)
9. **Commit**: Commit changes (`git commit -m "Add feature xxx"`) ensuring backward compatible code
10. **Push**: Push to the branch (`git push origin feature/xxx`).
11. **PR**: Open a merge request on GitHub (on the GitHub webpage) with detailed description.

Please ensure tests pass and include relevant documentation updates.

---

## 🌟 Support

Welcome to contribute to this project via submitting merge requests and reporting issues.

**Project Support:**

- ⭐ **Give GitHub stars** if this project helps you
- 🤝 **Share with teammates** and (golang) programming friends
- 📝 **Write tech blogs** about development tools and workflows - we provide content writing support
- 🌟 **Join the ecosystem** - committed to supporting open source and the (golang) development scene

**Have Fun Coding with this package!** 🎉🎉🎉

<!-- TEMPLATE (EN) END: STANDARD PROJECT FOOTER -->

---

## GitHub Stars

[![Stargazers](https://starchart.cc/go-xlan/go-rocketmq.svg?variant=adaptive)](https://starchart.cc/go-xlan/go-rocketmq)
