package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"http://go-rmqnamesrv:9876"}))
	if err != nil {
		panic("生成 producer 失败")
	}

	if err = p.Start(); err != nil {
		panic("启动 producer 失败")
	}

	msg := &primitive.Message{
		Topic: "uccs",
		Body:  []byte("this is a delay message 2"),
	}
	msg.WithDelayTimeLevel(4)
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("发送消息失败: %s\n", err)
		return
	}
	fmt.Printf("发送消息成功: %s\n", res.String())
	if err := p.Shutdown(); err != nil {
		panic("关闭 producer 失败")
	}
}
