package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"time"
)

func main() {
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"139.198.21.42:9876"}),
		consumer.WithGroupName("group1"),
		)
	if err := c.Subscribe("hellomq",  consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error){
		for i := range msgs {
			fmt.Printf("获取到值: %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		fmt.Println("读取消息失败")
	}

	_ = c.Start()
	//不能让主goroutine 推出
	time.Sleep(time.Hour)
	_ = c.Shutdown()
}

