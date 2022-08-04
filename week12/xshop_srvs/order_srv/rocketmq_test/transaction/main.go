package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type OrderListener struct {}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState{
	//return primitive.CommitMessageState
	fmt.Println("开始逻辑")
	time.Sleep(time.Second * 3)
	fmt.Println("逻辑成功")
	return primitive.CommitMessageState
	//return primitive.RollbackMessageState
}


func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState{
	return primitive.CommitMessageState
}


func main() {
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"139.198.21.42:9876"}),
		)
	if err != nil {
		panic("生成producer失败")
	}

	if err = p.Start(); err != nil {
		panic("启动Start失败")
	}

	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("hellomq", []byte("this is a transaction message")))
	if err != nil {
		fmt.Printf("发送失败:%s\n",err)
	}else{
		fmt.Printf("发送成功:%s\n",res.String())
	}

	time.Sleep(time.Hour)
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}
}