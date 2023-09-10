package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type Listener struct{}

func (l *Listener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	time.Sleep(time.Second * 3)
	fmt.Println("执行本地逻辑失败")
	// 本地执行逻辑无缘无故失败，比如代码异常，宕机等
	return primitive.UnknowState
}

func (l *Listener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("开始回查")
	time.Sleep(time.Second * 15)

	return primitive.CommitMessageState
}

func main() {
	p, err := rocketmq.NewTransactionProducer(
		&Listener{},
		producer.WithNameServer([]string{"http://go-rmqnamesrv:9876"}),
	)
	if err != nil {
		panic("生成 producer 失败")
	}

	if err = p.Start(); err != nil {
		panic("启动 producer 失败")
	}

	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("TransTopic", []byte("this is transaction message 回查")))

	if err != nil {
		fmt.Printf("发送消息失败: %s\n", err)
		return
	}
	fmt.Printf("发送消息成功: %s\n", res.String())

	time.Sleep(5 * time.Minute)
	if err = p.Shutdown(); err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}
