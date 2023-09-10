package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"http://go-rmqnamesrv:9876"}),
		consumer.WithGroupName("mxshop"),
	)
	if err := c.Subscribe("uccs", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("收到消息: %s\n", msgs[i].Body)
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		panic(err)
	}
	c.Start()

	time.Sleep(time.Hour)
	_ = c.Shutdown()

}
