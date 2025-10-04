[![GitHub Workflow Status (branch)](https://img.shields.io/github/actions/workflow/status/go-xlan/go-rocketmq/release.yml?branch=main&label=BUILD)](https://github.com/go-xlan/go-rocketmq/actions/workflows/release.yml?query=branch%3Amain)
[![GoDoc](https://pkg.go.dev/badge/github.com/go-xlan/go-rocketmq)](https://pkg.go.dev/github.com/go-xlan/go-rocketmq)
[![Coverage Status](https://img.shields.io/coveralls/github/go-xlan/go-rocketmq/main.svg)](https://coveralls.io/github/go-xlan/go-rocketmq?branch=main)
[![Supported Go Versions](https://img.shields.io/badge/Go-1.23+-lightgrey.svg)](https://go.dev/)
[![GitHub Release](https://img.shields.io/github/release/go-xlan/go-rocketmq.svg)](https://github.com/go-xlan/go-rocketmq/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-xlan/go-rocketmq)](https://goreportcard.com/report/github.com/go-xlan/go-rocketmq)

# go-rocketmq

Apache RocketMQ 客户端 Go 语言封装，提供简化的生产者和消费者接口。

---

<!-- TEMPLATE (ZH) BEGIN: LANGUAGE NAVIGATION -->
## 英文文档

[ENGLISH README](README.md)
<!-- TEMPLATE (ZH) END: LANGUAGE NAVIGATION -->

## 核心特性

🚀 **简化 API**: 为 RocketMQ 消息队列操作提供简洁的生产者和消费者 API
🔌 **自动解析**: 自动主机名解析，支持 IPv4 地址
🛡️ **无缝错误**: 与 mate 生态系统集成，提供优雅的错误处理
📝 **Go 实践**: 遵循 Go 语言最佳实践的接口设计
🔄 **自动拉取**: 基于回调的自动消息拉取和处理
⚙️ **可配置**: 消息发送操作的超时和重试设置

## 安装

```bash
go get github.com/go-xlan/go-rocketmq
```

## 前置要求

- Go 1.23.0 或更高版本
- Apache RocketMQ 服务器实例（可访问 NameServer）

## 使用方法

### 完整的生产者和消费者示例

此示例展示如何在同一个程序中创建生产者和消费者。

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

⬆️ **源码:** [源码](internal/demos/demo1x/main.go)

## API 参考

### 配置

- `Config` - 生产者和消费者的客户端配置
- `ProducerOptions` - 生产者特定配置选项

### 生产者

- `NewProducer(config *Config) (*Producer, error)` - 创建并启动新的生产者实例
- `SendMessage(ctx context.Context, topic string, payload []byte) error` - 发送消息到指定主题
- `Close() error` - 关闭生产者并释放资源

### 消费者

- `NewConsumer(config *Config) (*Consumer, error)` - 创建新的消费者实例
- `StartSubscribe(topic string, handler func(*primitive.MessageExt) (consumer.ConsumeResult, error)) error` - 启动消费者并订阅主题
- `Close() error` - 关闭消费者并释放资源

### 工具函数

- `ResolveNameServer(nameServer string) ([]string, error)` - 将 NameServer 主机名解析为 IPv4 地址

## 示例

### 创建生产者

**基本生产者设置：**
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

**发送消息：**
```go
payload := []byte(`{"name":"demo","value":123}`)
must.Done(producerCli.SendMessage(context.Background(), "MyTopic", payload))
```

### 创建消费者

**基本消费者设置：**
```go
config := &rocketmq.Config{
    NameServerAddress: "127.0.0.1:9876",
    GroupName:         "ConsumerGroup",
}
consumerCli := rese.P1(rocketmq.NewConsumer(config))
defer rese.F0(consumerCli.Close)
```

**订阅主题：**
```go
must.Done(consumerCli.StartSubscribe("MyTopic", func(message *primitive.MessageExt) (consumer.ConsumeResult, error) {
    fmt.Printf("Message: %s\n", string(message.Body))
    return consumer.ConsumeSuccess, nil
}))
```

### 主机名解析

**将 NameServer 主机名解析为 IPv4 地址：**
```go
addresses := rese.V1(rocketmq.ResolveNameServer("mq.example.com:9876"))
fmt.Printf("Resolved addresses: %v\n", addresses)
```

## 依赖

本项目使用 mate 生态系统进行优雅的错误处理和资源管理：

- `github.com/apache/rocketmq-client-go/v2` - Apache RocketMQ 官方 Go 客户端
- `github.com/yyle88/erero` - 错误处理
- `github.com/yyle88/must` - 基于 panic 的错误断言
- `github.com/yyle88/rese` - 带 panic 处理的资源管理
- `github.com/yyle88/zaplog` - 结构化日志
- `go.uber.org/zap` - 日志基础库

<!-- TEMPLATE (ZH) BEGIN: STANDARD PROJECT FOOTER -->
<!-- VERSION 2025-09-26 07:39:27.188023 +0000 UTC -->

## 📄 许可证类型

MIT 许可证。详见 [LICENSE](LICENSE)。

---

## 🤝 项目贡献

非常欢迎贡献代码！报告 BUG、建议功能、贡献代码：

- 🐛 **发现问题？** 在 GitHub 上提交问题并附上重现步骤
- 💡 **功能建议？** 创建 issue 讨论您的想法
- 📖 **文档疑惑？** 报告问题，帮助我们改进文档
- 🚀 **需要功能？** 分享使用场景，帮助理解需求
- ⚡ **性能瓶颈？** 报告慢操作，帮助我们优化性能
- 🔧 **配置困扰？** 询问复杂设置的相关问题
- 📢 **关注进展？** 关注仓库以获取新版本和功能
- 🌟 **成功案例？** 分享这个包如何改善工作流程
- 💬 **反馈意见？** 欢迎提出建议和意见

---

## 🔧 代码贡献

新代码贡献，请遵循此流程：

1. **Fork**：在 GitHub 上 Fork 仓库（使用网页界面）
2. **克隆**：克隆 Fork 的项目（`git clone https://github.com/yourname/repo-name.git`）
3. **导航**：进入克隆的项目（`cd repo-name`）
4. **分支**：创建功能分支（`git checkout -b feature/xxx`）
5. **编码**：实现您的更改并编写全面的测试
6. **测试**：（Golang 项目）确保测试通过（`go test ./...`）并遵循 Go 代码风格约定
7. **文档**：为面向用户的更改更新文档，并使用有意义的提交消息
8. **暂存**：暂存更改（`git add .`）
9. **提交**：提交更改（`git commit -m "Add feature xxx"`）确保向后兼容的代码
10. **推送**：推送到分支（`git push origin feature/xxx`）
11. **PR**：在 GitHub 上打开 Merge Request（在 GitHub 网页上）并提供详细描述

请确保测试通过并包含相关的文档更新。

---

## 🌟 项目支持

非常欢迎通过提交 Merge Request 和报告问题来为此项目做出贡献。

**项目支持：**

- ⭐ **给予星标**如果项目对您有帮助
- 🤝 **分享项目**给团队成员和（golang）编程朋友
- 📝 **撰写博客**关于开发工具和工作流程 - 我们提供写作支持
- 🌟 **加入生态** - 致力于支持开源和（golang）开发场景

**祝你用这个包编程愉快！** 🎉🎉🎉

<!-- TEMPLATE (ZH) END: STANDARD PROJECT FOOTER -->

---

## GitHub 标星点赞

[![Stargazers](https://starchart.cc/go-xlan/go-rocketmq.svg?variant=adaptive)](https://starchart.cc/go-xlan/go-rocketmq)
